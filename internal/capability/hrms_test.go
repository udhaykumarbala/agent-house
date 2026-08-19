package capability

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

// fakeWQ is a minimal in-memory Worqplace API: envelope responses, JWT-ish
// bearer auth with login/refresh counters, and a one-shot "expire the
// current token" switch to exercise the 401→refresh→retry path.
type fakeWQ struct {
	srv        *httptest.Server
	logins     int32
	refreshes  int32
	tokenGen   int32 // current valid token generation
	expireOnce int32 // when 1, the next authed request 401s once
}

func (f *fakeWQ) currentToken() string {
	return fmt.Sprintf("tok-%d", atomic.LoadInt32(&f.tokenGen))
}

func envelope(data any) map[string]any {
	return map[string]any{"success": true, "request_id": "test", "data": data}
}

func newFakeWQ(t *testing.T) *fakeWQ {
	t.Helper()
	f := &fakeWQ{tokenGen: 1}
	mux := http.NewServeMux()

	writeJSON := func(w http.ResponseWriter, code int, v any) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		_ = json.NewEncoder(w).Encode(v)
	}

	mux.HandleFunc("/api/v1/auth/login", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&f.logins, 1)
		var body struct{ Email, Password string }
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body.Email != "demo@x.com" || body.Password != "secret-pass" {
			writeJSON(w, 401, map[string]any{"success": false,
				"error": map[string]any{"code": "AUTH_INVALID_CREDENTIALS", "message": "bad creds"}})
			return
		}
		writeJSON(w, 200, envelope(map[string]any{
			"access_token": f.currentToken(), "refresh_token": "refresh-1",
			"user": map[string]any{"email": body.Email, "roles": []string{"viewer"}},
		}))
	})

	mux.HandleFunc("/api/v1/auth/refresh", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&f.refreshes, 1)
		atomic.AddInt32(&f.tokenGen, 1)
		writeJSON(w, 200, envelope(map[string]any{
			"access_token": f.currentToken(), "refresh_token": "refresh-2",
		}))
	})

	// healthz is NOT enveloped in the real API.
	mux.HandleFunc("/api/v1/healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{"status": "ok"})
	})

	authed := func(next func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			got := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			if atomic.CompareAndSwapInt32(&f.expireOnce, 1, 0) {
				atomic.AddInt32(&f.tokenGen, 1) // invalidate what the client holds
				writeJSON(w, 401, map[string]any{"success": false,
					"error": map[string]any{"code": "AUTH_TOKEN_EXPIRED", "message": "expired"}})
				return
			}
			if got != f.currentToken() {
				writeJSON(w, 401, map[string]any{"success": false,
					"error": map[string]any{"code": "AUTH_INVALID_TOKEN", "message": "invalid"}})
				return
			}
			next(w, r)
		}
	}

	mux.HandleFunc("/api/v1/reports/footer-stats", authed(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, envelope(map[string]any{
			"slug":    "footer-stats",
			"summary": map[string]any{"people_count": 859, "active_projects": 48},
		}))
	}))

	mux.HandleFunc("/api/v1/employees/tab-counts", authed(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, envelope(map[string]any{"everyone": 859, "docs_expiring": 504}))
	}))
	mux.HandleFunc("/api/v1/compliance/documents/counts", authed(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, envelope(map[string]any{"all": 900, "expired": 211}))
	}))
	mux.HandleFunc("/api/v1/triage/counts", authed(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, envelope(map[string]any{"badge": 7}))
	}))

	// The idle-engine reports break the uniform ReportResult shape: data is
	// a bare array, not an object. The client must tolerate both.
	mux.HandleFunc("/api/v1/reports/idle-by-project", authed(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, envelope([]map[string]any{
			{"project_name": "HEAD OFFICE", "employee_count": 7, "total_idle_cost": "54900.00"},
		}))
	}))

	mux.HandleFunc("/api/v1/employees", authed(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("page_size") != "2" {
			writeJSON(w, 400, map[string]any{"success": false,
				"error": map[string]any{"code": "VALIDATION_FAILED", "message": "expected page_size=2"}})
			return
		}
		resp := envelope([]map[string]any{
			{"emp_code": "EMP-1", "full_name": "Alpha"},
			{"emp_code": "EMP-2", "full_name": "Beta"},
		})
		resp["meta"] = map[string]any{"page": 1, "page_size": 2, "total": 859}
		writeJSON(w, 200, resp)
	}))

	f.srv = httptest.NewServer(mux)
	t.Cleanup(f.srv.Close)
	return f
}

func clientFor(f *fakeWQ) *HRMS {
	return NewHRMSClient(f.srv.URL+"/api/v1", "demo@x.com", "secret-pass")
}

func TestHRMSLoginOnceAcrossCalls(t *testing.T) {
	f := newFakeWQ(t)
	h := clientFor(f)
	for i := 0; i < 3; i++ {
		if _, err := h.Report("footer-stats", nil); err != nil {
			t.Fatalf("call %d: %v", i, err)
		}
	}
	if got := atomic.LoadInt32(&f.logins); got != 1 {
		t.Fatalf("expected exactly 1 login across 3 calls, got %d", got)
	}
}

func TestHRMSReportReturnsDecodedData(t *testing.T) {
	f := newFakeWQ(t)
	h := clientFor(f)
	rep, err := h.Report("footer-stats", nil)
	if err != nil {
		t.Fatal(err)
	}
	summary, _ := rep["summary"].(map[string]any)
	if summary == nil || summary["people_count"] != float64(859) {
		t.Fatalf("expected summary.people_count=859, got %#v", rep)
	}
}

func TestHRMSRefreshOnExpiredToken(t *testing.T) {
	f := newFakeWQ(t)
	h := clientFor(f)
	if _, err := h.Report("footer-stats", nil); err != nil {
		t.Fatalf("warmup: %v", err)
	}
	atomic.StoreInt32(&f.expireOnce, 1)
	if _, err := h.Report("footer-stats", nil); err != nil {
		t.Fatalf("expected transparent recovery after token expiry, got %v", err)
	}
	if atomic.LoadInt32(&f.refreshes) != 1 {
		t.Fatalf("expected 1 refresh, got %d", f.refreshes)
	}
	if atomic.LoadInt32(&f.logins) != 1 {
		t.Fatalf("expected no re-login (refresh should suffice), got %d logins", f.logins)
	}
}

func TestHRMSReportToleratesArrayData(t *testing.T) {
	f := newFakeWQ(t)
	h := clientFor(f)
	rep, err := h.Report("idle-by-project", nil)
	if err != nil {
		t.Fatalf("array-shaped report data must not error: %v", err)
	}
	rows, _ := rep["rows"].([]any)
	if len(rows) != 1 {
		t.Fatalf("expected array data wrapped as rows, got %#v", rep)
	}
	row, _ := rows[0].(map[string]any)
	if row["project_name"] != "HEAD OFFICE" {
		t.Fatalf("unexpected row: %#v", row)
	}
}

func TestHRMSNotConfigured(t *testing.T) {
	h := NewHRMSClient("", "", "")
	if h.Configured() {
		t.Fatal("empty client must report not configured")
	}
	if _, err := h.Report("footer-stats", nil); err == nil {
		t.Fatal("expected error from unconfigured client")
	}
}

func TestHRMSReportSlugAllowlist(t *testing.T) {
	f := newFakeWQ(t)
	h := clientFor(f)
	// Per-employee compensation is deliberately outside the demo surface.
	if _, err := h.Report("top-employees-by-cost", nil); err == nil ||
		!strings.Contains(err.Error(), "not in the view-only allowlist") {
		t.Fatalf("expected allowlist rejection, got %v", err)
	}
	// And unknown slugs never reach the network.
	if _, err := h.Report("../../admin/users", nil); err == nil {
		t.Fatal("expected rejection of arbitrary path as slug")
	}
}

func TestHRMSEmployeesPassesParamsAndMeta(t *testing.T) {
	f := newFakeWQ(t)
	h := clientFor(f)
	items, meta, err := h.Employees(map[string]string{"page_size": "2"})
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 || items[0]["emp_code"] != "EMP-1" {
		t.Fatalf("unexpected items: %#v", items)
	}
	if meta["total"] != float64(859) {
		t.Fatalf("expected meta.total=859, got %#v", meta)
	}
}

func TestHRMSCountsComposite(t *testing.T) {
	f := newFakeWQ(t)
	h := clientFor(f)
	c, err := h.Counts()
	if err != nil {
		t.Fatal(err)
	}
	emp, _ := c["employees"].(map[string]any)
	if emp == nil || emp["docs_expiring"] != float64(504) {
		t.Fatalf("expected employees.docs_expiring=504, got %#v", c)
	}
	if _, ok := c["compliance"]; !ok {
		t.Fatalf("expected compliance section, got %#v", c)
	}
	if _, ok := c["triage"]; !ok {
		t.Fatalf("expected triage section, got %#v", c)
	}
}

func TestHRMSStatusHealthyAndAuthed(t *testing.T) {
	f := newFakeWQ(t)
	h := clientFor(f)
	st := h.Status()
	if st["configured"] != true || st["healthy"] != true || st["authenticated"] != true {
		t.Fatalf("unexpected status: %#v", st)
	}
}

func TestSharedHRMSReadsEnvAndResets(t *testing.T) {
	f := newFakeWQ(t)
	t.Setenv("WQ_API", f.srv.URL+"/api/v1")
	t.Setenv("WQ_EMAIL", "demo@x.com")
	t.Setenv("WQ_PASSWORD", "secret-pass")
	ResetHRMSForTest()
	t.Cleanup(ResetHRMSForTest)

	h := SharedHRMS()
	if !h.Configured() {
		t.Fatal("SharedHRMS should pick up WQ_* env vars")
	}
	if h2 := SharedHRMS(); h2 != h {
		t.Fatal("SharedHRMS must return the same instance (token cache)")
	}
}
