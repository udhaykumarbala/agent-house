package scenario

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"pty-claude-test/internal/capability"
)

// fakeEstimateHRMS: a directory with enough of each trade to make the
// availability math interesting, plus the idle-cost report the blended
// day rate is derived from (42875.45 SAR / 397 days ≈ 108 SAR/day).
func fakeEstimateHRMS(t *testing.T) *httptest.Server {
	t.Helper()
	env := func(data any) map[string]any { return map[string]any{"success": true, "data": data} }
	send := func(w http.ResponseWriter, v any) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(v)
	}
	employees := []map[string]any{
		{"emp_code": "M1", "full_name": "MASON ONE", "designation": "Mason", "status": "active",
			"document_status": "valid", "idle_30d": 10, "project": ""},
		{"emp_code": "M2", "full_name": "MASON TWO", "designation": "MASON", "status": "active",
			"document_status": "valid", "idle_30d": 0, "project": "NCMS"},
		{"emp_code": "M3", "full_name": "MASON EXPIRED", "designation": "Mason", "status": "active",
			"document_status": "expired", "idle_30d": 0, "project": ""},
		{"emp_code": "S1", "full_name": "FIXER ONE", "designation": "Steel Fixer", "status": "active",
			"document_status": "valid", "idle_30d": 5, "project": ""},
		{"emp_code": "L1", "full_name": "LABOUR ONE", "designation": "Labour", "status": "active",
			"document_status": "valid", "idle_30d": 0, "project": "L & T"},
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/auth/login", func(w http.ResponseWriter, r *http.Request) {
		send(w, env(map[string]any{"access_token": "tok", "refresh_token": "ref"}))
	})
	mux.HandleFunc("/api/v1/employees", func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		if page != "" && page != "1" {
			resp := env([]map[string]any{})
			resp["meta"] = map[string]any{"total": len(employees)}
			send(w, resp)
			return
		}
		resp := env(employees)
		resp["meta"] = map[string]any{"total": len(employees)}
		send(w, resp)
	})
	mux.HandleFunc("/api/v1/reports/idle-cost", func(w http.ResponseWriter, r *http.Request) {
		send(w, env(map[string]any{"summary": map[string]any{
			"total_idle_cost": "42875.45", "total_idle_days": 397, "idle_employees": 25}}))
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func pointEstimateAtFake(t *testing.T) {
	t.Helper()
	srv := fakeEstimateHRMS(t)
	t.Setenv("WQ_API", srv.URL+"/api/v1")
	t.Setenv("WQ_EMAIL", "demo@x.com")
	t.Setenv("WQ_PASSWORD", "secret-pass")
	capability.ResetHRMSForTest()
	t.Cleanup(capability.ResetHRMSForTest)
}

func TestEstimateProjectManpowerCosting(t *testing.T) {
	pointEstimateAtFake(t)

	res, err := EstimateProject{}.Run(Context{
		Scope: "atlas-site",
		Input: map[string]any{
			"request":  "Warehouse construction in Dammam",
			"manpower": "mason:3, steel fixer:1",
			"months":   "2",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !res.OK {
		t.Fatalf("expected OK, got %#v", res)
	}

	var avail, cost *Step
	agents := map[string]bool{}
	for i := range res.Steps {
		agents[res.Steps[i].Agent] = true
		switch res.Steps[i].Action {
		case "check_availability":
			avail = &res.Steps[i]
		case "cost_estimate":
			cost = &res.Steps[i]
		}
	}
	for _, want := range []string{"conductor", "project_manager", "hr"} {
		if !agents[want] {
			t.Fatalf("story must involve %s, got %#v", want, res.Steps)
		}
	}
	if avail == nil || !avail.OK || cost == nil || !cost.OK {
		t.Fatalf("expected OK availability + costing steps, got %#v", res.Steps)
	}

	// Availability: masons — M1 eligible+available, M2 eligible (on NCMS),
	// M3 excluded (expired docs) → required 3, eligible 2, gap 1.
	rows, _ := avail.Output["rows"].([]map[string]any)
	if len(rows) != 2 {
		t.Fatalf("expected one availability row per trade, got %#v", rows)
	}
	var masonRow map[string]any
	for _, r := range rows {
		if r["trade"] == "mason" {
			masonRow = r
		}
	}
	if masonRow == nil || masonRow["required"] != 3 || masonRow["eligible"] != 2 || masonRow["gap"] != 1 {
		t.Fatalf("mason availability wrong: %#v", masonRow)
	}

	// Costing: blended live rate × 26 workdays × months × headcount, +10%.
	rate := 42875.45 / 397
	wantTotal := 4 * rate * 26 * 2 * 1.10
	got, _ := cost.Output["grand_total_sar"].(float64)
	if math.Abs(got-wantTotal) > 1.0 {
		t.Fatalf("grand total: got %.2f, want %.2f", got, wantTotal)
	}

	// The estimation must arrive as a markdown table with a hiring flag.
	if !strings.Contains(res.Summary, "| Trade") || !strings.Contains(res.Summary, "|---") {
		t.Fatalf("summary should carry a markdown estimation table, got %q", res.Summary)
	}
	if !strings.Contains(res.Summary, "SAR") || !strings.Contains(res.Summary, "mason") {
		t.Fatalf("summary should carry costs per trade, got %q", res.Summary)
	}
	hiring := false
	for _, s := range res.Suggestions {
		if strings.Contains(strings.ToLower(s.Title), "hir") {
			hiring = true
		}
	}
	if !hiring {
		t.Fatalf("a manpower gap must suggest hiring, got %#v", res.Suggestions)
	}
}

func TestEstimateProjectDefaultPlanFromRequest(t *testing.T) {
	pointEstimateAtFake(t)

	res, err := EstimateProject{}.Run(Context{
		Scope: "atlas-site",
		Input: map[string]any{"request": "New warehouse project, need an estimation", "months": "3"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !res.OK {
		t.Fatalf("expected OK with template-derived manpower, got %#v", res)
	}
	var req *Step
	for i := range res.Steps {
		if res.Steps[i].Action == "draft_requirements" {
			req = &res.Steps[i]
		}
	}
	if req == nil || !req.OK {
		t.Fatalf("expected a draft_requirements step, got %#v", res.Steps)
	}
	plan, _ := req.Output["manpower"].([]map[string]any)
	if len(plan) < 3 {
		t.Fatalf("default warehouse plan should carry several trades, got %#v", plan)
	}
}

func TestEstimateProjectPicksUpInboxRequest(t *testing.T) {
	pointEstimateAtFake(t)

	// The inbox path is fixed at ./projects/inbox — run in a temp cwd.
	old, _ := os.Getwd()
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })

	in := capability.NewInbox("")
	if _, err := in.Inject(capability.Email{
		ID: "mock_est_1", From: "client@newco.com", FromName: "NewCo Projects",
		Subject: "Request for Estimation — Warehouse Construction, Dammam",
		Body:    "Please share a manpower + cost estimation for a 4,500 sqm warehouse. Target duration 6 months.",
		Category: "client",
	}); err != nil {
		t.Fatal(err)
	}

	res, err := EstimateProject{}.Run(Context{Scope: "atlas-site", Input: map[string]any{}})
	if err != nil {
		t.Fatal(err)
	}
	if !res.OK {
		t.Fatalf("expected OK from inbox pickup, got %#v", res)
	}
	read := res.Steps[0]
	if read.Action != "read_request" || !strings.Contains(read.Note, "Warehouse Construction") {
		t.Fatalf("first step should read the estimation email from the inbox, got %#v", read)
	}
}

func TestEstimateProjectNoRequestAnywhere(t *testing.T) {
	pointEstimateAtFake(t)
	old, _ := os.Getwd()
	if err := os.Chdir(t.TempDir()); err != nil { // empty inbox
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })

	res, err := EstimateProject{}.Run(Context{Scope: "atlas-site", Input: map[string]any{}})
	if err != nil {
		t.Fatalf("must degrade, not error: %v", err)
	}
	if res.OK || !strings.Contains(strings.ToLower(res.Summary), "request") {
		t.Fatalf("expected graceful ask for the request, got %#v", res)
	}
}
