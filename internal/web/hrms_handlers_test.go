package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"pty-claude-test/internal/capability"
)

// fakeWQAPI is the minimal Worqplace fake for the handler layer.
func fakeWQAPI(t *testing.T) *httptest.Server {
	t.Helper()
	env := func(data any) map[string]any { return map[string]any{"success": true, "data": data} }
	send := func(w http.ResponseWriter, v any) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(v)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/auth/login", func(w http.ResponseWriter, r *http.Request) {
		send(w, env(map[string]any{"access_token": "tok", "refresh_token": "ref"}))
	})
	mux.HandleFunc("/api/v1/healthz", func(w http.ResponseWriter, r *http.Request) {
		send(w, map[string]any{"status": "ok"})
	})
	mux.HandleFunc("/api/v1/reports/footer-stats", func(w http.ResponseWriter, r *http.Request) {
		send(w, env(map[string]any{"summary": map[string]any{"people_count": 859}}))
	})
	mux.HandleFunc("/api/v1/employees", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("nationality") != "Nepalese" {
			send(w, env([]map[string]any{}))
			return
		}
		resp := env([]map[string]any{{"emp_code": "EMP-9", "full_name": "Gamma"}})
		resp["meta"] = map[string]any{"total": 1}
		send(w, resp)
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func pointHRMSAtFake(t *testing.T, base string) {
	t.Helper()
	t.Setenv("WQ_API", base+"/api/v1")
	t.Setenv("WQ_EMAIL", "demo@x.com")
	t.Setenv("WQ_PASSWORD", "secret-pass")
	capability.ResetHRMSForTest()
	t.Cleanup(capability.ResetHRMSForTest)
}

func TestHRMSStatusUnconfiguredIs200(t *testing.T) {
	t.Setenv("WQ_API", "")
	t.Setenv("WQ_EMAIL", "")
	t.Setenv("WQ_PASSWORD", "")
	capability.ResetHRMSForTest()
	t.Cleanup(capability.ResetHRMSForTest)

	s := &Server{}
	rec := httptest.NewRecorder()
	s.handleHRMSStatus(rec, httptest.NewRequest("GET", "/api/cap/hrms/status", nil))
	if rec.Code != 200 {
		t.Fatalf("status endpoint must never 500 on missing config, got %d", rec.Code)
	}
	var body map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &body)
	if body["configured"] != false {
		t.Fatalf("expected configured=false, got %#v", body)
	}
}

func TestHRMSHandlersAreGETOnly(t *testing.T) {
	s := &Server{}
	for path, h := range map[string]http.HandlerFunc{
		"/api/cap/hrms/status":                    s.handleHRMSStatus,
		"/api/cap/hrms/counts":                    s.handleHRMSCounts,
		"/api/cap/hrms/employees":                 s.handleHRMSEmployees,
		"/api/cap/hrms/reports/documents-expiring": s.handleHRMSReport,
	} {
		rec := httptest.NewRecorder()
		h(rec, httptest.NewRequest("POST", path, strings.NewReader("{}")))
		if rec.Code != 405 {
			t.Fatalf("POST %s: expected 405 (view-only surface), got %d", path, rec.Code)
		}
	}
}

func TestHRMSReportHandlerReturnsData(t *testing.T) {
	pointHRMSAtFake(t, fakeWQAPI(t).URL)
	s := &Server{}
	rec := httptest.NewRecorder()
	s.handleHRMSReport(rec, httptest.NewRequest("GET", "/api/cap/hrms/reports/footer-stats", nil))
	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var body map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &body)
	rep, _ := body["report"].(map[string]any)
	if rep == nil {
		t.Fatalf("expected report payload, got %s", rec.Body.String())
	}
}

func TestHRMSReportHandlerRejectsNonAllowlisted(t *testing.T) {
	pointHRMSAtFake(t, fakeWQAPI(t).URL)
	s := &Server{}
	rec := httptest.NewRecorder()
	s.handleHRMSReport(rec, httptest.NewRequest("GET", "/api/cap/hrms/reports/top-employees-by-cost", nil))
	if rec.Code != 400 {
		t.Fatalf("expected 400 for non-allowlisted report, got %d", rec.Code)
	}
}

func TestHRMSEmployeesHandlerForwardsFilters(t *testing.T) {
	pointHRMSAtFake(t, fakeWQAPI(t).URL)
	s := &Server{}
	rec := httptest.NewRecorder()
	s.handleHRMSEmployees(rec, httptest.NewRequest("GET", "/api/cap/hrms/employees?nationality=Nepalese", nil))
	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var body struct {
		Employees []map[string]any `json:"employees"`
		Meta      map[string]any   `json:"meta"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &body)
	if len(body.Employees) != 1 || body.Employees[0]["emp_code"] != "EMP-9" {
		t.Fatalf("nationality filter was not forwarded: %s", rec.Body.String())
	}
}
