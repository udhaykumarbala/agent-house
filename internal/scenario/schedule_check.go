package scenario

import (
	"fmt"
	"strings"

	"pty-claude-test/internal/capability"
)

// ScheduleCheck composes Conductor → Project Manager → schedule capability.
// Identifies slipping milestones and emits per-slip mitigation suggestions.
type ScheduleCheck struct{}

func (ScheduleCheck) Name() string { return "schedule_check" }
func (ScheduleCheck) Description() string {
	return "Scan the milestone register; flag slipping milestones; suggest mitigations per severity."
}
func (ScheduleCheck) Example() map[string]any { return map[string]any{} }

func (ScheduleCheck) Run(ctx Context) (Result, error) {
	res := Result{Scenario: "schedule_check", Scope: ctx.Scope}

	classify := Step{
		Agent: "conductor", Action: "decompose",
		Input: map[string]any{"intent": "schedule_check"},
		OK:    true, Note: "routing to Project Manager for milestone scan",
	}
	res.Steps = append(res.Steps, classify)
	HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, classify)

	sc, err := capability.NewSchedule(ctx.DataRoot, ctx.Scope)
	if err != nil {
		return res, err
	}
	slips, err := sc.Slips()
	if err != nil {
		return res, err
	}
	pmStep := Step{
		Agent: "project_manager", Action: "detect_slips",
		Output: map[string]any{"slips": slips, "count": len(slips)},
		OK:     true,
	}
	if len(slips) == 0 {
		pmStep.Note = "no slips detected; all milestones on track"
	} else {
		titles := make([]string, 0, len(slips))
		for _, s := range slips {
			titles = append(titles, fmt.Sprintf("%s (+%dd)", s.Milestone.Title, s.DaysLate))
		}
		pmStep.Note = fmt.Sprintf("%d slip(s): %s", len(slips), strings.Join(titles, " · "))
	}
	res.Steps = append(res.Steps, pmStep)
	HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, pmStep)

	if len(slips) == 0 {
		res.Summary = "All milestones on track."
		res.Suggestions = append(res.Suggestions, Suggestion{
			Title:  "Send the on-track report to client",
			Detail: "Routine update — no escalation needed.",
			Action: "notify_client_on_track",
		})
	} else {
		critical, moderate, minor := bucket(slips)
		res.Summary = fmt.Sprintf("%d slip(s): %d critical / %d moderate / %d minor.",
			len(slips), critical, moderate, minor)
		for _, sl := range slips {
			res.Suggestions = append(res.Suggestions, slipSuggestion(sl))
		}
	}
	res.OK = true
	return res, nil
}

func bucket(slips []capability.Slip) (critical, moderate, minor int) {
	for _, s := range slips {
		switch s.Severity {
		case "critical":
			critical++
		case "moderate":
			moderate++
		default:
			minor++
		}
	}
	return
}

func slipSuggestion(s capability.Slip) Suggestion {
	var title, detail, action string
	switch s.Severity {
	case "critical":
		title = fmt.Sprintf("Escalate: %s is %d days late", s.Milestone.Title, s.DaysLate)
		detail = "Critical slip — escalate to CEO, propose crew increase or scope cut."
		action = "escalate_slip"
	case "moderate":
		title = fmt.Sprintf("Re-sequence: %s is %d days late", s.Milestone.Title, s.DaysLate)
		detail = "Moderate slip — propose re-sequencing or adding a parallel crew."
		action = "resequence_workstream"
	default:
		title = fmt.Sprintf("Recover: %s is %d days late", s.Milestone.Title, s.DaysLate)
		detail = "Minor slip — Site Engineer to confirm catch-up plan."
		action = "ask_for_catchup_plan"
	}
	return Suggestion{
		Title: title, Detail: detail, Action: action,
		Payload: map[string]any{
			"milestone_id":  s.Milestone.ID,
			"days_late":     s.DaysLate,
			"severity":      s.Severity,
		},
	}
}
