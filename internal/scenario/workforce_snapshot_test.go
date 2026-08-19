package scenario

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"pty-claude-test/internal/capability"
)

// fakeHRMS serves just enough of the Worqplace API for the snapshot flow.
func fakeHRMS(t *testing.T) *httptest.Server {
	t.Helper()
	env := func(data any) map[string]any {
		return map[string]any{"success": true, "data": data}
	}
	mux := http.NewServeMux()
	send := func(w http.ResponseWriter, v any) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(v)
	}
	mux.HandleFunc("/api/v1/healthz", func(w http.ResponseWriter, r *http.Request) {
		send(w, map[string]any{"status": "ok"})
	})
	mux.HandleFunc("/api/v1/auth/login", func(w http.ResponseWriter, r *http.Request) {
		send(w, env(map[string]any{"access_token": "tok", "refresh_token": "ref"}))
	})
	mux.HandleFunc("/api/v1/reports/footer-stats", func(w http.ResponseWriter, r *http.Request) {
		send(w, env(map[string]any{"summary": map[string]any{
			"people_count": 859, "active_projects": 48, "payroll_years": 1}}))
	})
	mux.HandleFunc("/api/v1/employees/tab-counts", func(w http.ResponseWriter, r *http.Request) {
		send(w, env(map[string]any{"everyone": 859, "on_probation": 0,
			"docs_expiring": 504, "idle_this_week": 42, "saudi_nationals": 132}))
	})
	mux.HandleFunc("/api/v1/compliance/documents/counts", func(w http.ResponseWriter, r *http.Request) {
		send(w, env(map[string]any{"all": 2100, "valid": 1385, "expiring": 504, "expired": 211}))
	})
	mux.HandleFunc("/api/v1/triage/counts", func(w http.ResponseWriter, r *http.Request) {
		send(w, env(map[string]any{"badge": 7}))
	})
	mux.HandleFunc("/api/v1/reports/documents-expiring", func(w http.ResponseWriter, r *http.Request) {
		send(w, env(map[string]any{"summary": map[string]any{
			"total_expiring": 514, "expired_count": 211, "horizon_days": 60}}))
	})
	mux.HandleFunc("/api/v1/reports/saudization", func(w http.ResponseWriter, r *http.Request) {
		send(w, env(map[string]any{"summary": map[string]any{
			"ratio": 15.58, "saudi": 132, "total": 847}}))
	})
	mux.HandleFunc("/api/v1/employees", func(w http.ResponseWriter, r *http.Request) {
		// The real DB holds ZZTEST fixture employees that are newest by
		// created_at (the API default sort) — searches must sort by name
		// so demo output leads with real people.
		if r.URL.Query().Get("sort_by") != "full_name" || r.URL.Query().Get("sort_dir") != "asc" {
			send(w, map[string]any{"success": false,
				"error": map[string]any{"code": "VALIDATION_FAILED", "message": "expected sort_by=full_name&sort_dir=asc"}})
			return
		}
		if r.URL.Query().Get("nationality") != "Nepalese" {
			send(w, env([]map[string]any{}))
			return
		}
		resp := env([]map[string]any{
			{"emp_code": "EMP-101", "full_name": "Hari Sharma", "designation": "Fitter", "nationality": "Nepalese", "project": "L & T", "status": "active"},
			{"emp_code": "EMP-102", "full_name": "Bikram Rai", "designation": "Welder", "nationality": "Nepalese", "project": "NCMS", "status": "active"},
		})
		resp["meta"] = map[string]any{"total": 2}
		send(w, resp)
	})
	mux.HandleFunc("/api/v1/reports/headcount-by-project", func(w http.ResponseWriter, r *http.Request) {
		send(w, env(map[string]any{"rows": []map[string]any{
			{"project": "NEOM Site A", "headcount": 120},
			{"project": "Red Sea Camp", "headcount": 95},
		}}))
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func TestWorkforceSnapshotUnconfigured(t *testing.T) {
	t.Setenv("WQ_API", "")
	t.Setenv("WQ_EMAIL", "")
	t.Setenv("WQ_PASSWORD", "")
	capability.ResetHRMSForTest()
	t.Cleanup(capability.ResetHRMSForTest)

	res, err := WorkforceSnapshot{}.Run(Context{Scope: "atlas-site"})
	if err != nil {
		t.Fatalf("unconfigured must degrade, not error: %v", err)
	}
	if res.OK {
		t.Fatal("expected OK=false when HRMS is unconfigured")
	}
	if !strings.Contains(res.Summary, "WQ_") {
		t.Fatalf("summary should tell the operator which env vars to set, got %q", res.Summary)
	}
}

func TestWorkforceSnapshotLiveData(t *testing.T) {
	srv := fakeHRMS(t)
	t.Setenv("WQ_API", srv.URL+"/api/v1")
	t.Setenv("WQ_EMAIL", "demo@x.com")
	t.Setenv("WQ_PASSWORD", "secret-pass")
	capability.ResetHRMSForTest()
	t.Cleanup(capability.ResetHRMSForTest)

	res, err := WorkforceSnapshot{}.Run(Context{Scope: "atlas-site"})
	if err != nil {
		t.Fatal(err)
	}
	if !res.OK {
		t.Fatalf("expected OK, got %#v", res)
	}
	if len(res.Steps) < 4 {
		t.Fatalf("expected at least 4 steps (connect, headcount, compliance, saudization), got %d", len(res.Steps))
	}
	agents := map[string]bool{}
	for _, s := range res.Steps {
		agents[s.Agent] = true
		if !s.OK {
			t.Fatalf("step %s/%s failed: %s", s.Agent, s.Action, s.Note)
		}
	}
	if !agents["hr"] {
		t.Fatal("snapshot must attribute work to the hr agent")
	}
	if !strings.Contains(res.Summary, "859") {
		t.Fatalf("summary should carry the live headcount, got %q", res.Summary)
	}
	if len(res.Suggestions) == 0 {
		t.Fatal("504 expiring documents must produce at least one suggestion")
	}
}

// A directory filter turns the snapshot into a focused live search: connect
// + search_directory only, with real names surfaced in the summary so the
// Conductor answers "list our Nepalese employees" with people, not totals.
func TestWorkforceSnapshotDirectorySearch(t *testing.T) {
	srv := fakeHRMS(t)
	t.Setenv("WQ_API", srv.URL+"/api/v1")
	t.Setenv("WQ_EMAIL", "demo@x.com")
	t.Setenv("WQ_PASSWORD", "secret-pass")
	capability.ResetHRMSForTest()
	t.Cleanup(capability.ResetHRMSForTest)

	res, err := WorkforceSnapshot{}.Run(Context{
		Scope: "atlas-site",
		Input: map[string]any{"nationality": "Nepalese"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !res.OK {
		t.Fatalf("expected OK, got %#v", res)
	}
	var search *Step
	for i := range res.Steps {
		if res.Steps[i].Action == "search_directory" {
			search = &res.Steps[i]
		}
		if res.Steps[i].Action == "report_saudization" {
			t.Fatal("focused search must not run the full snapshot reports")
		}
	}
	if search == nil || search.Agent != "hr" || !search.OK {
		t.Fatalf("expected an OK hr/search_directory step, got %#v", res.Steps)
	}
	if !strings.Contains(res.Summary, "Hari Sharma") || !strings.Contains(res.Summary, "2") {
		t.Fatalf("summary should name the live matches, got %q", res.Summary)
	}
}

func TestWorkforceSnapshotRegisteredShape(t *testing.T) {
	r := WorkforceSnapshot{}
	if r.Name() != "workforce_snapshot" {
		t.Fatalf("unexpected name %q", r.Name())
	}
	if r.Description() == "" || r.Example() == nil {
		t.Fatal("runner must be self-describing for the dynamic UI")
	}
}
