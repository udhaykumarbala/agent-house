package scenario

import (
	"fmt"
	"strings"

	"pty-claude-test/internal/capability"
)

// WorkforceSnapshot pulls a live, view-only workforce picture from the real
// Worqplace HRMS (capability.SharedHRMS) instead of the local JSON stores:
// headcount + active projects, expiring documents, saudization ratio, and
// per-project staffing. Every number in the output came over the wire from
// the HRMS seconds ago — nothing is seeded.
type WorkforceSnapshot struct{}

func (WorkforceSnapshot) Name() string { return "workforce_snapshot" }
func (WorkforceSnapshot) Description() string {
	return "Live workforce snapshot from the Worqplace HRMS (view-only): headcount, active projects, expiring documents, saudization ratio, per-project staffing."
}
func (WorkforceSnapshot) Example() map[string]any {
	return map[string]any{"q": "", "nationality": ""}
}

func (WorkforceSnapshot) Run(ctx Context) (Result, error) {
	res := Result{Scenario: "workforce_snapshot", Scope: ctx.Scope}
	h := capability.SharedHRMS()

	// Degrade, don't error: the demo must survive a missing VPN or env.
	if !h.Configured() {
		res.Steps = append(res.Steps, Step{
			Agent: "hr", Action: "connect_hrms", OK: false,
			Note: "HRMS not configured — set WQ_API, WQ_EMAIL, WQ_PASSWORD in .env",
		})
		res.Summary = "HRMS integration is not configured (set WQ_API, WQ_EMAIL, WQ_PASSWORD); no live data available."
		res.OK = false
		return res, nil
	}

	st := h.Status()
	connect := Step{
		Agent: "conductor", Action: "connect_hrms",
		OK:     st["healthy"] == true && st["authenticated"] == true,
		Output: map[string]any{"status": st},
		Note:   fmt.Sprintf("Worqplace HRMS at %v — healthy=%v authenticated=%v", st["base_url"], st["healthy"], st["authenticated"]),
	}
	res.Steps = append(res.Steps, connect)
	HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, connect)
	if !connect.OK {
		res.Summary = fmt.Sprintf("Could not reach the live HRMS: %v", st["error"])
		res.OK = false
		return res, nil
	}

	// A directory filter (q or nationality) turns the run into a focused
	// live search instead of the full snapshot — this is how the Conductor
	// answers "list our Nepalese employees" with actual people.
	q := strFrom(ctx.Input, "q")
	nationality := strFrom(ctx.Input, "nationality")
	if q != "" || nationality != "" {
		return runDirectorySearch(ctx, res, h, q, nationality)
	}

	// HR: headcount + org counters.
	people, projects := 0, 0
	docsExpiring, idleThisWeek := 0, 0
	if rep, err := h.Report("footer-stats", nil); err == nil {
		people = intFrom(rep, "summary", "people_count")
		projects = intFrom(rep, "summary", "active_projects")
	}
	counts, countsErr := h.Counts()
	if countsErr == nil {
		if emp, ok := counts["employees"].(map[string]any); ok {
			docsExpiring = intVal(emp["docs_expiring"])
			idleThisWeek = intVal(emp["idle_this_week"])
		}
	}
	hc := Step{
		Agent: "hr", Action: "report_headcount", OK: countsErr == nil,
		Output: map[string]any{"people": people, "active_projects": projects, "counts": counts},
		Note: fmt.Sprintf("%d people across %d active project(s) · %d doc(s) expiring · %d idle this week",
			people, projects, docsExpiring, idleThisWeek),
	}
	res.Steps = append(res.Steps, hc)
	HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, hc)

	// HR: document compliance from the expiring-documents report.
	totalExpiring, expiredCount, horizon := 0, 0, 0
	docRep, docErr := h.Report("documents-expiring", nil)
	if docErr == nil {
		totalExpiring = intFrom(docRep, "summary", "total_expiring")
		expiredCount = intFrom(docRep, "summary", "expired_count")
		horizon = intFrom(docRep, "summary", "horizon_days")
	}
	comp := Step{
		Agent: "hr", Action: "report_compliance", OK: docErr == nil,
		Output: map[string]any{"total_expiring": totalExpiring, "expired": expiredCount, "horizon_days": horizon},
		Note: fmt.Sprintf("%d document(s) expiring within %d day(s), %d already expired",
			totalExpiring, horizon, expiredCount),
	}
	if docErr != nil {
		comp.Note = "documents-expiring report failed: " + docErr.Error()
	}
	res.Steps = append(res.Steps, comp)
	HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, comp)

	// HR: saudization (Nitaqat) ratio.
	ratio := 0.0
	saudi, totalNat := 0, 0
	saudRep, saudErr := h.Report("saudization", nil)
	if saudErr == nil {
		if sum, ok := saudRep["summary"].(map[string]any); ok {
			if f, ok := sum["ratio"].(float64); ok {
				ratio = f
			}
			saudi = intVal(sum["saudi"])
			totalNat = intVal(sum["total"])
		}
	}
	saud := Step{
		Agent: "hr", Action: "report_saudization", OK: saudErr == nil,
		Output: map[string]any{"ratio": ratio, "saudi": saudi, "total": totalNat},
		Note:   fmt.Sprintf("saudization ratio %.2f%% (%d of %d)", ratio, saudi, totalNat),
	}
	if saudErr != nil {
		saud.Note = "saudization report failed: " + saudErr.Error()
	}
	res.Steps = append(res.Steps, saud)
	HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, saud)

	// PM: top projects by headcount (first rows of the report).
	prjRep, prjErr := h.Report("headcount-by-project", nil)
	topProjects := []map[string]any{}
	if prjErr == nil {
		if rows, ok := prjRep["rows"].([]any); ok {
			for i, r := range rows {
				if i >= 5 {
					break
				}
				if m, ok := r.(map[string]any); ok {
					topProjects = append(topProjects, m)
				}
			}
		}
	}
	pm := Step{
		Agent: "project_manager", Action: "report_project_staffing", OK: prjErr == nil,
		Output: map[string]any{"top_projects": topProjects},
		Note:   fmt.Sprintf("top %d project(s) by headcount pulled from live HRMS", len(topProjects)),
	}
	if prjErr != nil {
		pm.Note = "headcount-by-project report failed: " + prjErr.Error()
	}
	res.Steps = append(res.Steps, pm)
	HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, pm)

	res.Summary = fmt.Sprintf(
		"Live HRMS: %d people · %d active project(s) · %d document(s) expiring (%d expired) · saudization %.2f%% · %d idle this week",
		people, projects, totalExpiring, expiredCount, ratio, idleThisWeek)

	if totalExpiring > 0 {
		res.Suggestions = append(res.Suggestions, Suggestion{
			Title:  fmt.Sprintf("Review %d expiring document(s)", totalExpiring),
			Detail: fmt.Sprintf("Pull the documents-expiring report (%d already expired) and start renewals.", expiredCount),
			Action: "view_hrms_report",
			Payload: map[string]any{
				"slug": "documents-expiring",
			},
		})
	}
	if idleThisWeek > 0 {
		res.Suggestions = append(res.Suggestions, Suggestion{
			Title:  fmt.Sprintf("Investigate %d idle employee(s) this week", idleThisWeek),
			Detail: "Idle workers are unbilled cost — check idle-by-project for where they sit.",
			Action: "view_hrms_report",
			Payload: map[string]any{
				"slug": "idle-by-project",
			},
		})
	}
	if len(res.Suggestions) == 0 {
		res.Suggestions = append(res.Suggestions, Suggestion{
			Title: "Workforce healthy · nothing requires action", Detail: "All live HRMS indicators clear.",
			Action: "send_status_to_client",
		})
	}
	res.OK = true
	return res, nil
}

// runDirectorySearch is the focused variant: one live /employees query with
// the matches surfaced by name in the summary.
func runDirectorySearch(ctx Context, res Result, h *capability.HRMS, q, nationality string) (Result, error) {
	// Sort by name: the API default (created_at desc) leads with the
	// client's ZZTEST fixture employees, which reads terribly in a demo.
	params := map[string]string{"page_size": "10", "sort_by": "full_name", "sort_dir": "asc"}
	label := []string{}
	if q != "" {
		params["q"] = q
		label = append(label, fmt.Sprintf("q=%s", q))
	}
	if nationality != "" {
		params["nationality"] = nationality
		label = append(label, fmt.Sprintf("nationality=%s", nationality))
	}
	items, meta, err := h.Employees(params)
	step := Step{
		Agent: "hr", Action: "search_directory", OK: err == nil,
		Input:  map[string]any{"q": q, "nationality": nationality},
		Output: map[string]any{"employees": items, "meta": meta},
	}
	if err != nil {
		step.Note = "live directory search failed: " + err.Error()
		res.Steps = append(res.Steps, step)
		HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, step)
		res.Summary = "Could not search the live HRMS directory: " + err.Error()
		res.OK = false
		return res, nil
	}
	total := intVal(meta["total"])
	if total == 0 {
		total = len(items)
	}
	names := []string{}
	for i, e := range items {
		if i >= 5 {
			break
		}
		name, _ := e["full_name"].(string)
		code, _ := e["emp_code"].(string)
		if name != "" {
			names = append(names, fmt.Sprintf("%s (%s)", name, code))
		}
	}
	step.Note = fmt.Sprintf("%d live match(es) for %s", total, strings.Join(label, ", "))
	res.Steps = append(res.Steps, step)
	HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, step)

	res.Summary = fmt.Sprintf("Live HRMS directory: %d match(es) for %s — %s",
		total, strings.Join(label, ", "), strings.Join(names, ", "))
	if total > len(names) {
		res.Summary += fmt.Sprintf(" (+%d more)", total-len(names))
	}
	res.Suggestions = append(res.Suggestions, Suggestion{
		Title:  "Open the live directory with these filters",
		Detail: "Full list with pagination via the HRMS employees endpoint.",
		Action: "view_hrms_employees",
		Payload: map[string]any{
			"q": q, "nationality": nationality,
		},
	})
	res.OK = true
	return res, nil
}

func strFrom(input map[string]any, key string) string {
	if input == nil {
		return ""
	}
	if s, ok := input[key].(string); ok {
		return strings.TrimSpace(s)
	}
	return ""
}

// intFrom digs data[section][key] and coerces the JSON float64 to int.
func intFrom(m map[string]any, section, key string) int {
	if s, ok := m[section].(map[string]any); ok {
		return intVal(s[key])
	}
	return 0
}

func intVal(v any) int {
	if f, ok := v.(float64); ok {
		return int(f)
	}
	return 0
}
