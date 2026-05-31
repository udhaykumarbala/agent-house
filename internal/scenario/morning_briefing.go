package scenario

import (
	"fmt"

	"pty-claude-test/internal/capability"
)

// MorningBriefing is the composite ritual: each agent reports a short
// status; Conductor assembles them into the morning brief. Pure
// composition over the existing capabilities — no new data sources.
type MorningBriefing struct{}

func (MorningBriefing) Name() string { return "morning_briefing" }
func (MorningBriefing) Description() string {
	return "Compose a one-page morning briefing: schedule slips (PM) + inbox state (Conductor) + HR pipeline + Procurement vendor health."
}
func (MorningBriefing) Example() map[string]any {
	return map[string]any{}
}

func (MorningBriefing) Run(ctx Context) (Result, error) {
	res := Result{Scenario: "morning_briefing", Scope: ctx.Scope}

	intent := Step{
		Agent: "conductor", Action: "compose_briefing",
		OK:   true,
		Note: "polling each discipline for a status update",
	}
	res.Steps = append(res.Steps, intent)
	HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, intent)

	// PM: schedule slips
	sc, _ := capability.NewSchedule(ctx.DataRoot, ctx.Scope)
	slips, _ := sc.Slips()
	ms, _ := sc.List()
	critical, moderate, minor := bucket(slips)
	pm := Step{
		Agent: "project_manager", Action: "report_schedule", OK: true,
		Output: map[string]any{
			"milestones":  len(ms),
			"slips":       slips,
			"critical":    critical,
			"moderate":    moderate,
			"minor":       minor,
		},
		Note: fmt.Sprintf("%d milestone(s) · %d slip(s) (%d critical / %d moderate / %d minor)",
			len(ms), len(slips), critical, moderate, minor),
	}
	res.Steps = append(res.Steps, pm)
	HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, pm)

	// HR: pipeline
	hr, _ := capability.NewHR(ctx.DataRoot, ctx.Scope)
	apps, _ := hr.List()
	byStatus := map[string]int{}
	for _, a := range apps {
		st := a.Status
		if st == "" {
			st = "new"
		}
		byStatus[st]++
	}
	hrStep := Step{
		Agent: "hr", Action: "report_pipeline", OK: true,
		Output: map[string]any{"total": len(apps), "by_status": byStatus},
		Note:   fmt.Sprintf("%d applicant(s) in pipeline", len(apps)),
	}
	res.Steps = append(res.Steps, hrStep)
	HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, hrStep)

	// Procurement: vendor health
	pr, _ := capability.NewProcurement(ctx.DataRoot, ctx.Scope)
	vendors, _ := pr.List()
	activeVendors := 0
	for _, v := range vendors {
		if v.ContractActive {
			activeVendors++
		}
	}
	prStep := Step{
		Agent: "procurement", Action: "report_vendor_health", OK: true,
		Output: map[string]any{"vendors": len(vendors), "active_contracts": activeVendors},
		Note:   fmt.Sprintf("%d vendor(s) · %d active contract(s)", len(vendors), activeVendors),
	}
	res.Steps = append(res.Steps, prStep)
	HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, prStep)

	// Conductor: inbox snapshot
	in := capability.NewInbox(ctx.DataRoot)
	mailSummary, _ := in.Summary()
	mailStep := Step{
		Agent: "conductor", Action: "report_inbox", OK: true,
		Output: map[string]any{"inbox": mailSummary},
		Note:   fmt.Sprintf("%d unread · %d impersonation flag(s)", mailSummary.Total, mailSummary.Impersonations),
	}
	res.Steps = append(res.Steps, mailStep)
	HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, mailStep)

	// Summary + suggestions.
	res.Summary = fmt.Sprintf("%d milestone(s) · %d slip(s) (%d critical) · %d applicant(s) · %d vendor(s) (%d active) · %d unread mail(s) (%d impersonation flag(s))",
		len(ms), len(slips), critical, len(apps), len(vendors), activeVendors,
		mailSummary.Total, mailSummary.Impersonations)

	if critical > 0 {
		res.Suggestions = append(res.Suggestions, Suggestion{
			Title:  fmt.Sprintf("Escalate %d critical schedule slip(s)", critical),
			Detail: "Run schedule_check for the full breakdown and mitigations.",
			Action: "run_schedule_check",
		})
	}
	if mailSummary.Impersonations > 0 {
		res.Suggestions = append(res.Suggestions, Suggestion{
			Title:  fmt.Sprintf("Investigate %d impersonation flag(s) in inbox", mailSummary.Impersonations),
			Detail: "Procurement to validate; finance to freeze related payments.",
			Action: "investigate_impersonations",
		})
	}
	if mailSummary.Total > 0 {
		res.Suggestions = append(res.Suggestions, Suggestion{
			Title:  fmt.Sprintf("Process %d unread email(s)", mailSummary.Total),
			Detail: "Triage the inbox; agents will handle by category.",
			Action: "run_process_inbox",
		})
	}
	if len(apps) > 0 {
		res.Suggestions = append(res.Suggestions, Suggestion{
			Title:  "Review HR pipeline",
			Detail: "Match new applicants against open JDs.",
			Action: "run_process_applicants",
		})
	}
	if len(res.Suggestions) == 0 {
		res.Suggestions = append(res.Suggestions, Suggestion{
			Title: "Quiet morning · nothing requires action", Detail: "All disciplines reporting clear.",
			Action: "send_status_to_client",
		})
	}
	res.OK = true
	return res, nil
}
