package scenario

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"pty-claude-test/internal/capability"
)

// StaffProject is the story scenario that joins live project data with the
// live employee directory: "NCMS needs 2 welders — who can we send?"
//
//  1. Conductor parses the staffing request.
//  2. Project Manager pulls the target project's live context (headcount +
//     idle exposure) from the HRMS reports.
//  3. HR sweeps the full live directory for the requested trade.
//  4. HR screens eligibility — active, documents valid, not already on the
//     target project — and ranks idle workers first (they are already on
//     payroll doing nothing, so moving them erases idle cost).
//  5. Conductor recommends the picks.
//
// View-only end to end: the actual assignment happens in Worqplace by a
// human; the suggestion chip says exactly that.
type StaffProject struct{}

func (StaffProject) Name() string { return "staff_project" }
func (StaffProject) Description() string {
	return "Staff a project from the live HRMS: pull the project's headcount + idle exposure, sweep the directory for the requested trade, screen eligibility (active, valid documents, not already on the project) and rank idle workers first."
}
func (StaffProject) Example() map[string]any {
	return map[string]any{"role": "welder", "project": "NCMS", "count": "2"}
}

// staffCandidate is one screened employee with the human-readable reasons
// that made (or would break) the cut.
type staffCandidate struct {
	rec     map[string]any
	idle    int
	project string
}

func (StaffProject) Run(ctx Context) (Result, error) {
	res := Result{Scenario: "staff_project", Scope: ctx.Scope}
	h := capability.SharedHRMS()

	role := strings.ToLower(strFrom(ctx.Input, "role"))
	project := strFrom(ctx.Input, "project")
	count := 3
	if n, err := strconv.Atoi(strFrom(ctx.Input, "count")); err == nil && n > 0 {
		count = n
	}

	if role == "" {
		res.Summary = "Staffing request needs a role/trade (e.g. welder, electrician, fitter). Ask again with the role to fill."
		res.OK = false
		return res, nil
	}
	if !h.Configured() {
		res.Summary = "HRMS integration is not configured (set WQ_API, WQ_EMAIL, WQ_PASSWORD); cannot screen the live directory."
		res.OK = false
		return res, nil
	}

	// 1) Conductor frames the request.
	want := fmt.Sprintf("%d × %s", count, role)
	if project != "" {
		want += " for " + project
	}
	parse := Step{
		Agent: "conductor", Action: "parse_request", OK: true,
		Input: map[string]any{"role": role, "project": project, "count": count},
		Note:  "staffing request: " + want,
	}
	res.Steps = append(res.Steps, parse)
	HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, parse)

	// 2) PM: live project context from the reports.
	if project != "" {
		headcount := -1
		if rep, err := h.Report("headcount-by-project", nil); err == nil {
			for _, r := range rowsOf(rep) {
				name, _ := r["project"].(string)
				if matchProject(name, project) {
					headcount = intVal(r["headcount"])
					break
				}
			}
		}
		idleCount, idleCost := 0, ""
		if rep, err := h.Report("idle-by-project", nil); err == nil {
			for _, r := range rowsOf(rep) {
				name, _ := r["project_name"].(string)
				if matchProject(name, project) {
					idleCount = intVal(r["employee_count"])
					idleCost, _ = r["total_idle_cost"].(string)
					break
				}
			}
		}
		note := fmt.Sprintf("%s currently runs %d people", project, headcount)
		if headcount < 0 {
			note = fmt.Sprintf("%s not found in the live headcount report", project)
		}
		if idleCount > 0 {
			note += fmt.Sprintf(" · %d idle costing %s SAR", idleCount, idleCost)
		}
		pm := Step{
			Agent: "project_manager", Action: "report_project_context", OK: headcount >= 0,
			Output: map[string]any{"project": project, "headcount": headcount, "idle": idleCount, "idle_cost_sar": idleCost},
			Note:   note,
		}
		res.Steps = append(res.Steps, pm)
		HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, pm)
	}

	// 3) HR: sweep the whole live directory (designation is not a server-side
	// filter, so page through and match client-side).
	all, err := fetchAllEmployees(h, 10)
	if err != nil {
		fail := Step{Agent: "hr", Action: "scan_directory", OK: false, Note: "directory sweep failed: " + err.Error()}
		res.Steps = append(res.Steps, fail)
		HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, fail)
		res.Summary = "Could not sweep the live HRMS directory: " + err.Error()
		res.OK = false
		return res, nil
	}
	matches := []map[string]any{}
	for _, e := range all {
		desig, _ := e["designation"].(string)
		name, _ := e["full_name"].(string)
		// ZZTEST records are the client's DB fixtures — never staffable.
		if strings.HasPrefix(strings.ToLower(name), "zztest") ||
			strings.HasPrefix(strings.ToLower(desig), "zztest") {
			continue
		}
		if strings.Contains(strings.ToLower(desig), role) {
			matches = append(matches, e)
		}
	}
	scan := Step{
		Agent: "hr", Action: "scan_directory", OK: true,
		Output: map[string]any{"screened": len(all), "matches": len(matches)},
		Note:   fmt.Sprintf("screened %d live employee record(s) → %d %s(s) found", len(all), len(matches), role),
	}
	res.Steps = append(res.Steps, scan)
	HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, scan)

	// 4) HR: eligibility screen + ranking.
	eligible := []staffCandidate{}
	for _, e := range matches {
		status, _ := e["status"].(string)
		docs, _ := e["document_status"].(string)
		onProject, _ := e["project"].(string)
		if status == "terminated" || docs != "valid" {
			continue
		}
		if project != "" && onProject != "" && matchProject(onProject, project) {
			continue // already there
		}
		eligible = append(eligible, staffCandidate{rec: e, idle: intVal(e["idle_30d"]), project: onProject})
	}
	// Idle first (moving them erases idle cost), then unassigned, then name.
	sort.SliceStable(eligible, func(i, j int) bool {
		if eligible[i].idle != eligible[j].idle {
			return eligible[i].idle > eligible[j].idle
		}
		if (eligible[i].project == "") != (eligible[j].project == "") {
			return eligible[i].project == ""
		}
		ni, _ := eligible[i].rec["full_name"].(string)
		nj, _ := eligible[j].rec["full_name"].(string)
		return ni < nj
	})

	picks := []map[string]any{}
	for i, c := range eligible {
		if i >= count {
			break
		}
		name, _ := c.rec["full_name"].(string)
		code, _ := c.rec["emp_code"].(string)
		desig, _ := c.rec["designation"].(string)
		nat, _ := c.rec["nationality"].(string)
		why := []string{"documents valid"}
		if c.idle > 0 {
			why = append(why, fmt.Sprintf("idle %dd — already on payroll, zero mobilization cost", c.idle))
		}
		if c.project == "" {
			why = append(why, "currently unassigned")
		} else {
			why = append(why, "currently at "+c.project)
		}
		picks = append(picks, map[string]any{
			"emp_code": code, "full_name": name, "designation": desig,
			"nationality": nat, "idle_30d": c.idle, "current_project": c.project,
			"reasons": why,
		})
	}
	screen := Step{
		Agent: "hr", Action: "screen_eligibility", OK: true,
		Output: map[string]any{"eligible": len(eligible), "picks": picks},
		Note: fmt.Sprintf("%d of %d %s(s) eligible (active, documents valid, not already on the project) — proposing top %d",
			len(eligible), len(matches), role, len(picks)),
	}
	res.Steps = append(res.Steps, screen)
	HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, screen)

	// 5) Conductor recommends.
	names := []string{}
	for _, p := range picks {
		names = append(names, fmt.Sprintf("%s (%s%s)", p["full_name"], p["emp_code"],
			map[bool]string{true: fmt.Sprintf(", idle %dd", p["idle_30d"]), false: ""}[intVal(p["idle_30d"]) > 0]))
	}
	rec := Step{
		Agent: "conductor", Action: "recommend_staffing", OK: len(picks) > 0,
		Note: "recommendation ready: " + strings.Join(names, "; "),
	}
	if len(picks) == 0 {
		rec.Note = "no eligible candidates — consider hiring (process_applicants) or relaxing constraints"
	}
	res.Steps = append(res.Steps, rec)
	HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, rec)

	if len(picks) > 0 {
		// Markdown table — the Conductor chat renders GFM tables natively.
		var b strings.Builder
		fmt.Fprintf(&b, "Staffing request %s: screened %d live records → %d %s(s) → %d eligible.\n\n",
			want, len(all), len(matches), role, len(eligible))
		b.WriteString("| Name | Code | Designation | Idle 30d | Currently | Nationality |\n")
		b.WriteString("|---|---|---|---|---|---|\n")
		for _, p := range picks {
			cur, _ := p["current_project"].(string)
			if cur == "" {
				cur = "unassigned"
			}
			fmt.Fprintf(&b, "| %v | %v | %v | %vd | %s | %v |\n",
				p["full_name"], p["emp_code"], p["designation"], p["idle_30d"], cur, p["nationality"])
		}
		b.WriteString("\nIdle workers rank first — they are already on payroll, so mobilizing them erases idle cost.")
		res.Summary = b.String()
		res.Suggestions = append(res.Suggestions, Suggestion{
			Title:  fmt.Sprintf("Assign %d pick(s) in Worqplace", len(picks)),
			Detail: "Agent House is view-only on the HRMS — confirm the assignment in Worqplace itself.",
			Action: "open_hrms",
		})
		res.Suggestions = append(res.Suggestions, Suggestion{
			Title:  "Cross-check pick documents before mobilization",
			Detail: "Run the documents-expiring report for the picked employees' codes.",
			Action: "view_hrms_report",
			Payload: map[string]any{
				"slug": "documents-expiring",
			},
		})
	} else {
		res.Summary = fmt.Sprintf("Staffing request %s: screened %d live records but found no eligible %s(s) (matches: %d). Suggest opening a hiring request.",
			want, len(all), role, len(matches))
		res.Suggestions = append(res.Suggestions, Suggestion{
			Title:  "Open a hiring request for " + role,
			Detail: "Screen the applicant pipeline for this trade instead.",
			Action: "run_process_applicants",
		})
	}
	res.OK = true
	return res, nil
}

// fetchAllEmployees pages through the live directory (max maxPages × 100).
func fetchAllEmployees(h *capability.HRMS, maxPages int) ([]map[string]any, error) {
	out := []map[string]any{}
	for page := 1; page <= maxPages; page++ {
		items, _, err := h.Employees(map[string]string{
			"page": strconv.Itoa(page), "page_size": "100",
			"sort_by": "full_name", "sort_dir": "asc",
		})
		if err != nil {
			if page == 1 {
				return nil, err
			}
			break // partial sweep is still useful
		}
		out = append(out, items...)
		if len(items) < 100 {
			break
		}
	}
	return out, nil
}

// matchProject compares an HRMS project name against a spoken/typed one,
// surviving punctuation and "and"/"&" variants ("L & T" vs "L and T").
func matchProject(hrmsName, spoken string) bool {
	a, b := normalizeProject(hrmsName), normalizeProject(spoken)
	if a == "" || b == "" {
		return false
	}
	return strings.Contains(a, b) || strings.Contains(b, a)
}

func normalizeProject(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "&", " and ")
	var out strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			out.WriteRune(r)
		}
	}
	res := out.String()
	// "and" is pure noise between initials: "l and t" ≡ "l t" ≡ "l&t".
	return strings.ReplaceAll(res, "and", "")
}

// rowsOf extracts a report's rows tolerant of both ReportResult shapes.
func rowsOf(rep map[string]any) []map[string]any {
	rows, _ := rep["rows"].([]any)
	out := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		if m, ok := r.(map[string]any); ok {
			out = append(out, m)
		}
	}
	return out
}
