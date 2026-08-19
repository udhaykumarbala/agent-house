package scenario

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"pty-claude-test/internal/capability"
)

// fakeStaffHRMS serves a small directory with every eligibility edge case:
// idle vs assigned welders, expired documents, someone already on the
// target project, a ZZTEST fixture, and a non-matching designation.
func fakeStaffHRMS(t *testing.T) *httptest.Server {
	t.Helper()
	env := func(data any) map[string]any { return map[string]any{"success": true, "data": data} }
	send := func(w http.ResponseWriter, v any) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(v)
	}
	employees := []map[string]any{
		{"emp_code": "W1", "full_name": "WELDER IDLE", "designation": "Welder", "status": "active",
			"document_status": "valid", "idle_30d": 15, "project": "71 CAMP", "nationality": "Nepalese"},
		{"emp_code": "W2", "full_name": "WELDER FREE", "designation": "WELDER 6G", "status": "active",
			"document_status": "valid", "idle_30d": 0, "project": "", "nationality": "Indian"},
		{"emp_code": "W3", "full_name": "WELDER BUSY ON TARGET", "designation": "Welder", "status": "active",
			"document_status": "valid", "idle_30d": 0, "project": "NCMS", "nationality": "Indian"},
		{"emp_code": "W4", "full_name": "WELDER EXPIRED DOCS", "designation": "Welder", "status": "active",
			"document_status": "expired", "idle_30d": 20, "project": "", "nationality": "Nepalese"},
		{"emp_code": "W5", "full_name": "ZZTEST Welder", "designation": "ZZTEST Welder", "status": "active",
			"document_status": "valid", "idle_30d": 30, "project": "", "nationality": "Nepalese"},
		{"emp_code": "E1", "full_name": "SPARKY ELECTRICIAN", "designation": "Electrician", "status": "active",
			"document_status": "valid", "idle_30d": 9, "project": "", "nationality": "Indian"},
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/auth/login", func(w http.ResponseWriter, r *http.Request) {
		send(w, env(map[string]any{"access_token": "tok", "refresh_token": "ref"}))
	})
	mux.HandleFunc("/api/v1/employees", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("page") != "" && r.URL.Query().Get("page") != "1" {
			resp := env([]map[string]any{})
			resp["meta"] = map[string]any{"total": len(employees)}
			send(w, resp)
			return
		}
		resp := env(employees)
		resp["meta"] = map[string]any{"total": len(employees)}
		send(w, resp)
	})
	mux.HandleFunc("/api/v1/reports/headcount-by-project", func(w http.ResponseWriter, r *http.Request) {
		send(w, env(map[string]any{"rows": []map[string]any{
			{"project": "NCMS", "code": "NCMS", "headcount": 187},
			{"project": "L & T", "code": "LT", "headcount": 192},
		}}))
	})
	// Idle-engine report: bare array data (the real API's shape).
	mux.HandleFunc("/api/v1/reports/idle-by-project", func(w http.ResponseWriter, r *http.Request) {
		send(w, env([]map[string]any{
			{"project_name": "NCMS", "employee_count": 3, "total_idle_cost": "1503.67"},
		}))
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func TestStaffProjectEligiblePicking(t *testing.T) {
	srv := fakeStaffHRMS(t)
	t.Setenv("WQ_API", srv.URL+"/api/v1")
	t.Setenv("WQ_EMAIL", "demo@x.com")
	t.Setenv("WQ_PASSWORD", "secret-pass")
	capability.ResetHRMSForTest()
	t.Cleanup(capability.ResetHRMSForTest)

	res, err := StaffProject{}.Run(Context{
		Scope: "atlas-site",
		Input: map[string]any{"role": "welder", "project": "NCMS", "count": "2"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !res.OK {
		t.Fatalf("expected OK, got %#v", res)
	}

	agents := map[string]bool{}
	var screen *Step
	for i := range res.Steps {
		agents[res.Steps[i].Agent] = true
		if res.Steps[i].Action == "screen_eligibility" {
			screen = &res.Steps[i]
		}
	}
	for _, want := range []string{"conductor", "project_manager", "hr"} {
		if !agents[want] {
			t.Fatalf("story must involve %s, steps: %#v", want, res.Steps)
		}
	}
	if screen == nil || !screen.OK {
		t.Fatalf("expected an OK screen_eligibility step, got %#v", res.Steps)
	}
	picks, _ := screen.Output["picks"].([]map[string]any)
	if len(picks) != 2 {
		t.Fatalf("asked for 2 eligible picks, got %d: %#v", len(picks), picks)
	}
	// Idle welder ranks first (already paid, doing nothing), unassigned second.
	if picks[0]["emp_code"] != "W1" || picks[1]["emp_code"] != "W2" {
		t.Fatalf("wrong ranking: %#v", picks)
	}
	// Exclusions: expired docs, already on target project, ZZTEST fixtures.
	for _, p := range picks {
		if code := p["emp_code"]; code == "W3" || code == "W4" || code == "W5" {
			t.Fatalf("ineligible employee picked: %#v", p)
		}
	}
	if !strings.Contains(res.Summary, "NCMS") || !strings.Contains(res.Summary, "WELDER IDLE") {
		t.Fatalf("summary should tell the story with project + top pick, got %q", res.Summary)
	}
	// Project context came from the live reports.
	found := false
	for _, s := range res.Steps {
		if s.Action == "report_project_context" && strings.Contains(s.Note, "187") {
			found = true
		}
	}
	if !found {
		t.Fatalf("PM step should carry live project headcount, steps: %#v", res.Steps)
	}
}

// The Brain hands over spoken project names ("L and T"), the HRMS stores
// "L & T" — matching must survive punctuation and and/& variants.
func TestProjectNameFuzzyMatch(t *testing.T) {
	cases := []struct {
		hrms, spoken string
		want         bool
	}{
		{"L & T", "L and T", true},
		{"L & T", "l&t", true},
		{"L & T", "LT", true},
		{"NCMS", "ncms", true},
		{"HEAD OFFICE", "head office", true},
		{"NCMS", "L and T", false},
		{"71 CAMP IDLE", "NCMS", false},
	}
	for _, c := range cases {
		if got := matchProject(c.hrms, c.spoken); got != c.want {
			t.Errorf("matchProject(%q, %q) = %v, want %v", c.hrms, c.spoken, got, c.want)
		}
	}
}

func TestStaffProjectRequiresRole(t *testing.T) {
	srv := fakeStaffHRMS(t)
	t.Setenv("WQ_API", srv.URL+"/api/v1")
	t.Setenv("WQ_EMAIL", "demo@x.com")
	t.Setenv("WQ_PASSWORD", "secret-pass")
	capability.ResetHRMSForTest()
	t.Cleanup(capability.ResetHRMSForTest)

	res, err := StaffProject{}.Run(Context{Scope: "atlas-site", Input: map[string]any{}})
	if err != nil {
		t.Fatalf("missing role must degrade, not error: %v", err)
	}
	if res.OK {
		t.Fatal("expected OK=false without a role")
	}
	if !strings.Contains(strings.ToLower(res.Summary), "role") {
		t.Fatalf("summary should ask for the role, got %q", res.Summary)
	}
}

func TestStaffProjectUnconfigured(t *testing.T) {
	t.Setenv("WQ_API", "")
	t.Setenv("WQ_EMAIL", "")
	t.Setenv("WQ_PASSWORD", "")
	capability.ResetHRMSForTest()
	t.Cleanup(capability.ResetHRMSForTest)

	res, err := StaffProject{}.Run(Context{Scope: "atlas-site", Input: map[string]any{"role": "welder"}})
	if err != nil {
		t.Fatalf("unconfigured must degrade, not error: %v", err)
	}
	if res.OK || !strings.Contains(res.Summary, "WQ_") {
		t.Fatalf("expected graceful unconfigured result, got %#v", res)
	}
}
