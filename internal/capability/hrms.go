package capability

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

// HRMS is a VIEW-ONLY client for the Worqplace HRMS REST API. Unlike the
// file-backed capabilities it is a process-wide singleton (SharedHRMS):
// Worqplace rate-limits login to 5/min/IP, so the 15-minute access token
// must be cached and rotated via /auth/refresh, never re-acquired per
// request. The client exposes no mutating methods by construction — the
// only POSTs it ever issues are /auth/login and /auth/refresh.
type HRMS struct {
	base     string // e.g. http://worqplace.alredaa/api/v1 (no trailing slash)
	email    string
	password string
	hc       *http.Client

	mu      sync.Mutex
	access  string
	refresh string
}

// ErrHRMSNotConfigured is returned when WQ_API / WQ_EMAIL / WQ_PASSWORD
// are absent, so callers can degrade gracefully instead of erroring the
// whole demo (same philosophy as the Resend email sender).
var ErrHRMSNotConfigured = errors.New("hrms not configured: set WQ_API, WQ_EMAIL, WQ_PASSWORD")

// hrmsReportAllow is the demo's view-only report surface. Deliberately
// absent: top-employees-by-cost (per-employee compensation) and anything
// payslip/salary-structure shaped — the demo shows aggregates only.
var hrmsReportAllow = map[string]bool{
	"footer-stats":              true,
	"headcount-by-project":      true,
	"payroll-cost-mom":          true,
	"saudization":               true,
	"documents-expiring":        true,
	"leave-approval-turnaround": true,
	"loan-exposure":             true,
	"document-renewal-cost":     true,
	"idle-cost":                 true,
	"idle-by-project":           true,
	"idle-exposure":             true,
}

// NewHRMSClient builds a client with explicit config (tests use this).
func NewHRMSClient(baseURL, email, password string) *HRMS {
	return &HRMS{
		base:     strings.TrimRight(baseURL, "/"),
		email:    email,
		password: password,
		hc:       &http.Client{Timeout: 15 * time.Second},
	}
}

var (
	hrmsMu   sync.Mutex
	hrmsInst *HRMS
)

// SharedHRMS returns the process-wide client, lazily built from env on
// first use. Always the same instance so the token cache survives across
// requests and scenario runs.
func SharedHRMS() *HRMS {
	hrmsMu.Lock()
	defer hrmsMu.Unlock()
	if hrmsInst == nil {
		hrmsInst = NewHRMSClient(os.Getenv("WQ_API"), os.Getenv("WQ_EMAIL"), os.Getenv("WQ_PASSWORD"))
	}
	return hrmsInst
}

// ResetHRMSForTest drops the singleton so tests can re-point WQ_* env vars.
func ResetHRMSForTest() {
	hrmsMu.Lock()
	hrmsInst = nil
	hrmsMu.Unlock()
}

// Configured reports whether the client has enough config to try talking
// to a real Worqplace instance.
func (h *HRMS) Configured() bool {
	return h.base != "" && h.email != "" && h.password != ""
}

// hrmsEnvelope is Worqplace's uniform response wrapper.
type hrmsEnvelope struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
	Meta    json.RawMessage `json:"meta"`
	Error   *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (h *HRMS) postJSON(path string, body any) (*hrmsEnvelope, error) {
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	resp, err := h.hc.Post(h.base+path, "application/json", bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var env hrmsEnvelope
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&env); err != nil {
		return nil, fmt.Errorf("hrms %s: bad response: %w", path, err)
	}
	if !env.Success {
		code := "UNKNOWN"
		msg := resp.Status
		if env.Error != nil {
			code, msg = env.Error.Code, env.Error.Message
		}
		return nil, fmt.Errorf("hrms %s: %s: %s", path, code, msg)
	}
	return &env, nil
}

// login acquires a fresh token pair. Caller must hold h.mu.
func (h *HRMS) login() error {
	env, err := h.postJSON("/auth/login", map[string]string{"email": h.email, "password": h.password})
	if err != nil {
		return err
	}
	return h.storeTokens(env.Data)
}

// refreshTokens rotates the pair; falls back to a full login when the
// refresh token itself is rejected. Caller must hold h.mu.
func (h *HRMS) refreshTokens() error {
	if h.refresh == "" {
		return h.login()
	}
	env, err := h.postJSON("/auth/refresh", map[string]string{"refresh_token": h.refresh})
	if err != nil {
		return h.login()
	}
	return h.storeTokens(env.Data)
}

func (h *HRMS) storeTokens(data json.RawMessage) error {
	var tok struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.Unmarshal(data, &tok); err != nil || tok.AccessToken == "" {
		return fmt.Errorf("hrms auth: token missing in response")
	}
	h.access, h.refresh = tok.AccessToken, tok.RefreshToken
	return nil
}

func (h *HRMS) token() (string, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.access == "" {
		if err := h.login(); err != nil {
			return "", err
		}
	}
	return h.access, nil
}

// get performs an authenticated GET, transparently refreshing (or
// re-logging-in) once on a 401 before giving up.
func (h *HRMS) get(path string, params map[string]string) (*hrmsEnvelope, error) {
	if !h.Configured() {
		return nil, ErrHRMSNotConfigured
	}
	tok, err := h.token()
	if err != nil {
		return nil, err
	}
	env, status, err := h.doGet(path, params, tok)
	if err != nil {
		return nil, err
	}
	if status == http.StatusUnauthorized {
		h.mu.Lock()
		refreshErr := h.refreshTokens()
		tok = h.access
		h.mu.Unlock()
		if refreshErr != nil {
			return nil, refreshErr
		}
		if env, status, err = h.doGet(path, params, tok); err != nil {
			return nil, err
		}
	}
	if status != http.StatusOK || !env.Success {
		code := fmt.Sprintf("HTTP %d", status)
		msg := ""
		if env.Error != nil {
			code, msg = env.Error.Code, env.Error.Message
		}
		return nil, fmt.Errorf("hrms GET %s: %s %s", path, code, msg)
	}
	return env, nil
}

func (h *HRMS) doGet(path string, params map[string]string, tok string) (*hrmsEnvelope, int, error) {
	u := h.base + path
	if len(params) > 0 {
		q := url.Values{}
		for k, v := range params {
			if v != "" {
				q.Set(k, v)
			}
		}
		u += "?" + q.Encode()
	}
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Authorization", "Bearer "+tok)
	resp, err := h.hc.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	var env hrmsEnvelope
	if err := json.NewDecoder(io.LimitReader(resp.Body, 8<<20)).Decode(&env); err != nil {
		return nil, resp.StatusCode, fmt.Errorf("hrms GET %s: bad response: %w", path, err)
	}
	return &env, resp.StatusCode, nil
}

// Report fetches one report by slug. Slugs are allowlisted so the demo
// surface stays aggregates-only regardless of what the account can read.
func (h *HRMS) Report(slug string, params map[string]string) (map[string]any, error) {
	if !hrmsReportAllow[slug] {
		return nil, fmt.Errorf("report %q is not in the view-only allowlist", slug)
	}
	env, err := h.get("/reports/"+slug, params)
	if err != nil {
		return nil, err
	}
	var out map[string]any
	if err := json.Unmarshal(env.Data, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Employees returns one page of the org-wide directory plus pagination meta.
func (h *HRMS) Employees(params map[string]string) ([]map[string]any, map[string]any, error) {
	env, err := h.get("/employees", params)
	if err != nil {
		return nil, nil, err
	}
	var items []map[string]any
	if err := json.Unmarshal(env.Data, &items); err != nil {
		return nil, nil, err
	}
	meta := map[string]any{}
	if len(env.Meta) > 0 {
		_ = json.Unmarshal(env.Meta, &meta)
	}
	return items, meta, nil
}

// Counts rolls the cheap org-wide counters into one snapshot. Sections
// degrade independently: a failing endpoint becomes {"error": ...} rather
// than sinking the whole snapshot.
func (h *HRMS) Counts() (map[string]any, error) {
	if !h.Configured() {
		return nil, ErrHRMSNotConfigured
	}
	out := map[string]any{}
	sections := []struct{ key, path string }{
		{"employees", "/employees/tab-counts"},
		{"compliance", "/compliance/documents/counts"},
		{"triage", "/triage/counts"},
	}
	var firstErr error
	for _, s := range sections {
		env, err := h.get(s.path, nil)
		if err != nil {
			out[s.key] = map[string]any{"error": err.Error()}
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		var v map[string]any
		if err := json.Unmarshal(env.Data, &v); err != nil {
			out[s.key] = map[string]any{"error": err.Error()}
			continue
		}
		out[s.key] = v
		firstErr = nil // at least one section worked; report partial data
	}
	if len(out) > 0 {
		return out, nil
	}
	return nil, firstErr
}

// Status answers "is the live HRMS reachable and are we logged in" without
// ever returning an error — it is the demo's connectivity indicator.
func (h *HRMS) Status() map[string]any {
	st := map[string]any{
		"configured":    h.Configured(),
		"healthy":       false,
		"authenticated": false,
	}
	if !h.Configured() {
		st["error"] = ErrHRMSNotConfigured.Error()
		return st
	}
	st["base_url"] = h.base

	// healthz is unauthenticated and NOT enveloped: {"status":"ok"}.
	resp, err := h.hc.Get(h.base + "/healthz")
	if err != nil {
		st["error"] = err.Error()
		return st
	}
	var hz struct {
		Status string `json:"status"`
	}
	err = json.NewDecoder(io.LimitReader(resp.Body, 4096)).Decode(&hz)
	resp.Body.Close()
	if err != nil || hz.Status != "ok" {
		st["error"] = fmt.Sprintf("healthz not ok (%v)", err)
		return st
	}
	st["healthy"] = true

	if _, err := h.token(); err != nil {
		st["error"] = err.Error()
		return st
	}
	st["authenticated"] = true
	return st
}
