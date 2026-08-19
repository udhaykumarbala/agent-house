package web

import (
	"net/http"
	"strings"

	"pty-claude-test/internal/capability"
)

// ── HRMS (live Worqplace, view-only) ──────────────────────────────────
//
// Unlike the file-backed capabilities these handlers proxy a real external
// HRMS through capability.SharedHRMS() (one process-wide client so the
// 15-minute token is cached; login is rate-limited upstream). The whole
// surface is GET — Agent House never mutates the client's HRMS.

// hrmsEmployeeFilters is the pass-through allowlist for /employees. Keeps
// the demo from becoming an open proxy for arbitrary query params.
var hrmsEmployeeFilters = []string{
	"tab", "q", "status", "nationality", "department_id", "project_id",
	"page", "page_size", "sort_by", "sort_dir",
}

func requireGET(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed (HRMS surface is view-only)", 405)
		return false
	}
	return true
}

// handleHRMSStatus reports connectivity + auth state. Never 500s: an
// unconfigured or unreachable HRMS is a state to display, not an error.
func (s *Server) handleHRMSStatus(w http.ResponseWriter, r *http.Request) {
	if !requireGET(w, r) {
		return
	}
	writeJSON(w, capability.SharedHRMS().Status())
}

// handleHRMSCounts returns the composite org-wide counters (employees /
// compliance / triage) in one call.
func (s *Server) handleHRMSCounts(w http.ResponseWriter, r *http.Request) {
	if !requireGET(w, r) {
		return
	}
	counts, err := capability.SharedHRMS().Counts()
	if err != nil {
		http.Error(w, err.Error(), 502)
		return
	}
	writeJSON(w, counts)
}

// handleHRMSEmployees proxies one page of the live directory, forwarding
// only allowlisted filters.
func (s *Server) handleHRMSEmployees(w http.ResponseWriter, r *http.Request) {
	if !requireGET(w, r) {
		return
	}
	params := map[string]string{}
	for _, k := range hrmsEmployeeFilters {
		if v := r.URL.Query().Get(k); v != "" {
			params[k] = v
		}
	}
	items, meta, err := capability.SharedHRMS().Employees(params)
	if err != nil {
		http.Error(w, err.Error(), 502)
		return
	}
	writeJSON(w, map[string]any{"employees": items, "meta": meta, "count": len(items)})
}

// handleHRMSReport serves GET /api/cap/hrms/reports/{slug}. The slug
// allowlist lives in the capability layer (aggregates only).
func (s *Server) handleHRMSReport(w http.ResponseWriter, r *http.Request) {
	if !requireGET(w, r) {
		return
	}
	slug := strings.TrimPrefix(r.URL.Path, "/api/cap/hrms/reports/")
	if slug == "" || strings.Contains(slug, "/") {
		http.Error(w, "report slug required: GET /api/cap/hrms/reports/{slug}", 400)
		return
	}
	params := map[string]string{}
	for _, k := range []string{"year", "days", "from", "to", "limit"} {
		if v := r.URL.Query().Get(k); v != "" {
			params[k] = v
		}
	}
	rep, err := capability.SharedHRMS().Report(slug, params)
	if err != nil {
		code := 502
		if strings.Contains(err.Error(), "allowlist") {
			code = 400
		}
		http.Error(w, err.Error(), code)
		return
	}
	writeJSON(w, map[string]any{"slug": slug, "report": rep})
}
