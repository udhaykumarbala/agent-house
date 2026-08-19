package scenario

import (
	"fmt"
	"strconv"
	"strings"

	"pty-claude-test/internal/capability"
)

// EstimateProject is the estimation story: a client asks for a new-project
// estimate (by email or in chat), the AI drafts the requirements including
// a manpower plan, HR checks live availability for every trade in the real
// HRMS and costs the plan from live data, and the Conductor composes the
// estimation table ready to send back.
//
//  1. Conductor reads the request — from params, or the newest estimation-
//     looking email in the inbox.
//  2. Project Manager drafts requirements: the manpower plan comes from the
//     Brain (AI-generated "trade:count" list passed as a param) or from a
//     deterministic template keyed on the request text.
//  3. HR sweeps the live directory once and reports, per trade: matches,
//     eligible (active + valid documents), available now (idle/unassigned),
//     and the hiring gap.
//  4. HR costs the plan with a blended day rate derived LIVE from the HRMS
//     idle-cost report (total idle cost / total idle days) — real payroll-
//     grade costing without exposing a single individual salary.
//  5. Conductor composes the estimation as a markdown table (the chat UI
//     renders GFM tables) with a grand total and next-step chips.
type EstimateProject struct{}

func (EstimateProject) Name() string { return "estimate_project" }
func (EstimateProject) Description() string {
	return "Turn a client estimation request (email or chat) into a costed manpower estimate: AI-drafted requirements, live HRMS availability per trade, and costing from a blended day rate derived live from the idle-cost report."
}
func (EstimateProject) Example() map[string]any {
	return map[string]any{
		"request":  "Warehouse construction in Dammam, 4,500 sqm",
		"manpower": "mason:4, steel fixer:6, electrician:2, labour:10",
		"months":   "6",
	}
}

// tradePlan is one line of the manpower requirement.
type tradePlan struct {
	Trade string
	Count int
}

const (
	workdaysPerMonth = 26
	overheadPct      = 0.10 // supervision + PPE + consumables
	fallbackDayRate  = 150.0
)

func (EstimateProject) Run(ctx Context) (Result, error) {
	res := Result{Scenario: "estimate_project", Scope: ctx.Scope}
	h := capability.SharedHRMS()

	// 1) Conductor: read the request (param first, inbox second).
	request := strFrom(ctx.Input, "request")
	source := "chat"
	if request == "" {
		if em := latestEstimationEmail(); em != nil {
			request = em.Subject
			if em.Body != "" {
				request += " — " + em.Body
			}
			source = fmt.Sprintf("inbox email from %s", em.From)
			read := Step{
				Agent: "conductor", Action: "read_request", OK: true,
				Input:  map[string]any{"email_id": em.ID},
				Output: map[string]any{"from": em.From, "subject": em.Subject},
				Note:   fmt.Sprintf("estimation request picked up from inbox: %q from %s", em.Subject, em.From),
			}
			res.Steps = append(res.Steps, read)
			HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, read)
		}
	} else {
		read := Step{
			Agent: "conductor", Action: "read_request", OK: true,
			Note: "estimation request: " + trimTo(request, 140),
		}
		res.Steps = append(res.Steps, read)
		HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, read)
	}
	if request == "" {
		res.Summary = "No estimation request found — pass one in chat or drop the client's email into the inbox first."
		res.OK = false
		return res, nil
	}
	if !h.Configured() {
		res.Summary = "HRMS integration is not configured (set WQ_API, WQ_EMAIL, WQ_PASSWORD); cannot check availability or cost the plan."
		res.OK = false
		return res, nil
	}

	months := 3
	if n, err := strconv.Atoi(strFrom(ctx.Input, "months")); err == nil && n > 0 {
		months = n
	}

	// 2) PM: requirements. AI-generated manpower param wins; else template.
	plan := parseManpower(strFrom(ctx.Input, "manpower"))
	planSource := "AI-drafted"
	if len(plan) == 0 {
		plan = defaultPlan(request)
		planSource = "template-drafted"
	}
	planOut := make([]map[string]any, 0, len(plan))
	planNote := make([]string, 0, len(plan))
	totalHead := 0
	for _, p := range plan {
		planOut = append(planOut, map[string]any{"trade": p.Trade, "count": p.Count})
		planNote = append(planNote, fmt.Sprintf("%d × %s", p.Count, p.Trade))
		totalHead += p.Count
	}
	reqStep := Step{
		Agent: "project_manager", Action: "draft_requirements", OK: true,
		Output: map[string]any{"manpower": planOut, "months": months, "source": planSource},
		Note:   fmt.Sprintf("%s requirements: %s over %d month(s)", planSource, strings.Join(planNote, ", "), months),
	}
	res.Steps = append(res.Steps, reqStep)
	HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, reqStep)

	// 3) HR: one live directory sweep, availability per trade.
	all, err := fetchAllEmployees(h, 10)
	if err != nil {
		fail := Step{Agent: "hr", Action: "check_availability", OK: false, Note: "directory sweep failed: " + err.Error()}
		res.Steps = append(res.Steps, fail)
		HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, fail)
		res.Summary = "Could not check live availability: " + err.Error()
		res.OK = false
		return res, nil
	}
	rows := make([]map[string]any, 0, len(plan))
	availNote := make([]string, 0, len(plan))
	totalGap := 0
	for _, p := range plan {
		matches, eligible, available := 0, 0, 0
		for _, e := range all {
			desig, _ := e["designation"].(string)
			name, _ := e["full_name"].(string)
			if strings.HasPrefix(strings.ToLower(name), "zztest") ||
				strings.HasPrefix(strings.ToLower(desig), "zztest") {
				continue
			}
			if !strings.Contains(strings.ToLower(desig), p.Trade) {
				continue
			}
			matches++
			status, _ := e["status"].(string)
			docs, _ := e["document_status"].(string)
			if status == "terminated" || docs != "valid" {
				continue
			}
			eligible++
			onProject, _ := e["project"].(string)
			if intVal(e["idle_30d"]) > 0 || onProject == "" {
				available++
			}
		}
		gap := p.Count - eligible
		if gap < 0 {
			gap = 0
		}
		totalGap += gap
		rows = append(rows, map[string]any{
			"trade": p.Trade, "required": p.Count, "matches": matches,
			"eligible": eligible, "available_now": available, "gap": gap,
		})
		availNote = append(availNote, fmt.Sprintf("%s %d/%d", p.Trade, eligible, p.Count))
	}
	availStep := Step{
		Agent: "hr", Action: "check_availability", OK: true,
		Output: map[string]any{"screened": len(all), "rows": rows, "total_gap": totalGap},
		Note: fmt.Sprintf("screened %d live records — eligible/required: %s · hiring gap %d",
			len(all), strings.Join(availNote, ", "), totalGap),
	}
	res.Steps = append(res.Steps, availStep)
	HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, availStep)

	// 4) HR: costing from the live blended day rate.
	rate, rateSource := blendedDayRate(h)
	laborTotal := float64(totalHead) * rate * workdaysPerMonth * float64(months)
	grandTotal := laborTotal * (1 + overheadPct)
	costStep := Step{
		Agent: "hr", Action: "cost_estimate", OK: true,
		Output: map[string]any{
			"day_rate_sar": rate, "rate_source": rateSource,
			"workdays_per_month": workdaysPerMonth, "months": months,
			"labor_total_sar": laborTotal, "overhead_pct": overheadPct,
			"grand_total_sar": grandTotal,
		},
		Note: fmt.Sprintf("blended day rate %.0f SAR (%s) × %d workday(s)/month × %d month(s) × %d worker(s) +%d%% overhead → %s SAR",
			rate, rateSource, workdaysPerMonth, months, totalHead, int(overheadPct*100), fmtSAR(grandTotal)),
	}
	res.Steps = append(res.Steps, costStep)
	HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, costStep)

	// 5) Conductor: compose the estimation table.
	var b strings.Builder
	fmt.Fprintf(&b, "Estimation — %s (%d month(s), %d worker(s), source: %s):\n\n",
		trimTo(request, 90), months, totalHead, source)
	b.WriteString("| Trade | Required | Eligible in-house | Available now | Hiring gap | Monthly (SAR) | Total (SAR) |\n")
	b.WriteString("|---|---|---|---|---|---|---|\n")
	for _, r := range rows {
		req := r["required"].(int)
		monthly := float64(req) * rate * workdaysPerMonth
		fmt.Fprintf(&b, "| %s | %d | %d | %d | %d | %s | %s |\n",
			r["trade"], req, r["eligible"], r["available_now"], r["gap"],
			fmtSAR(monthly), fmtSAR(monthly*float64(months)))
	}
	fmt.Fprintf(&b, "| **Total** | **%d** |  |  | **%d** | **%s** | **%s** |\n",
		totalHead, totalGap, fmtSAR(laborTotal/float64(months)), fmtSAR(laborTotal))
	fmt.Fprintf(&b, "\nLabor basis: blended day rate **%.0f SAR/day** (%s) · %d workdays/month · +%d%% supervision & consumables overhead.\n\n**Grand total: %s SAR**",
		rate, rateSource, workdaysPerMonth, int(overheadPct*100), fmtSAR(grandTotal))

	compose := Step{
		Agent: "conductor", Action: "compose_estimation", OK: true,
		Note: fmt.Sprintf("estimation ready: %d worker(s), %d month(s), grand total %s SAR (hiring gap %d)",
			totalHead, months, fmtSAR(grandTotal), totalGap),
	}
	res.Steps = append(res.Steps, compose)
	HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, compose)

	res.Summary = b.String()
	res.Suggestions = append(res.Suggestions, Suggestion{
		Title:  "Send the estimation reply to the client",
		Detail: "Draft the outbound email with the table above.",
		Action: "send_estimation_email",
		Payload: map[string]any{
			"request": request, "grand_total_sar": grandTotal, "months": months,
		},
	})
	if totalGap > 0 {
		res.Suggestions = append(res.Suggestions, Suggestion{
			Title:  fmt.Sprintf("Open hiring for the %d-worker shortfall", totalGap),
			Detail: "Screen the applicant pipeline (process_applicants) for the missing trades.",
			Action: "run_process_applicants",
		})
	}
	res.OK = true
	return res, nil
}

// parseManpower parses "mason:4, steel fixer:6" into a plan. Unparseable
// fragments are skipped; a missing count defaults to 1.
func parseManpower(s string) []tradePlan {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	out := []tradePlan{}
	for _, part := range strings.FieldsFunc(s, func(r rune) bool { return r == ',' || r == ';' }) {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		trade, count := part, 1
		if i := strings.LastIndex(part, ":"); i >= 0 {
			trade = strings.TrimSpace(part[:i])
			if n, err := strconv.Atoi(strings.TrimSpace(part[i+1:])); err == nil && n > 0 {
				count = n
			}
		}
		trade = strings.ToLower(strings.TrimSpace(trade))
		if trade != "" {
			out = append(out, tradePlan{Trade: trade, Count: count})
		}
	}
	return out
}

// defaultPlan is the deterministic fallback when no AI-drafted manpower
// arrives: a crew template keyed on what the request sounds like, using
// trades that actually exist in the live directory.
func defaultPlan(request string) []tradePlan {
	r := strings.ToLower(request)
	switch {
	case strings.Contains(r, "warehouse") || strings.Contains(r, "factory") || strings.Contains(r, "industrial"):
		return []tradePlan{
			{"steel fixer", 6}, {"shuttering carpenter", 4}, {"mason", 4},
			{"electrician", 2}, {"plumber", 1}, {"labour", 10},
		}
	case strings.Contains(r, "villa") || strings.Contains(r, "residential") || strings.Contains(r, "building"):
		return []tradePlan{
			{"mason", 6}, {"carpenter", 4}, {"electrician", 2},
			{"plumber", 2}, {"labour", 8},
		}
	case strings.Contains(r, "camp"):
		return []tradePlan{
			{"labour", 10}, {"electrician", 2}, {"plumber", 2}, {"carpenter", 2},
		}
	default:
		return []tradePlan{
			{"mason", 2}, {"carpenter", 2}, {"electrician", 1}, {"labour", 8},
		}
	}
}

// blendedDayRate derives an org-wide day rate from the live idle-cost
// report (total idle cost / total idle days) — real money data with zero
// individual salary exposure. Falls back to a fixed rate when the report
// is empty or unreachable.
func blendedDayRate(h *capability.HRMS) (float64, string) {
	rep, err := h.Report("idle-cost", nil)
	if err != nil {
		return fallbackDayRate, "fallback rate (idle-cost report unavailable)"
	}
	sum, _ := rep["summary"].(map[string]any)
	if sum == nil {
		return fallbackDayRate, "fallback rate (idle-cost report empty)"
	}
	days := float64(intVal(sum["total_idle_days"]))
	cost := 0.0
	switch v := sum["total_idle_cost"].(type) {
	case string:
		cost, _ = strconv.ParseFloat(v, 64)
	case float64:
		cost = v
	}
	if days <= 0 || cost <= 0 {
		return fallbackDayRate, "fallback rate (no idle data)"
	}
	return cost / days, "derived live from HRMS idle-cost"
}

// latestEstimationEmail scans the global inbox for the newest email that
// reads like an estimation/quotation request.
func latestEstimationEmail() *capability.Email {
	in := capability.NewInbox("")
	all, err := in.List() // newest first
	if err != nil {
		return nil
	}
	for i := range all {
		hay := strings.ToLower(all[i].Subject + " " + all[i].Body)
		for _, kw := range []string{"estimat", "quotation", "quote", "rfq", "proposal", "tender"} {
			if strings.Contains(hay, kw) {
				return &all[i]
			}
		}
	}
	return nil
}

// fmtSAR renders a SAR amount with thousands separators, no decimals.
func fmtSAR(f float64) string {
	n := int64(f + 0.5)
	s := strconv.FormatInt(n, 10)
	neg := strings.HasPrefix(s, "-")
	if neg {
		s = s[1:]
	}
	var parts []string
	for len(s) > 3 {
		parts = append([]string{s[len(s)-3:]}, parts...)
		s = s[:len(s)-3]
	}
	parts = append([]string{s}, parts...)
	out := strings.Join(parts, ",")
	if neg {
		out = "-" + out
	}
	return out
}

func trimTo(s string, n int) string {
	s = strings.TrimSpace(s)
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
