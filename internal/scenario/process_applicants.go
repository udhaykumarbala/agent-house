package scenario

import (
	"fmt"
	"strings"

	"pty-claude-test/internal/capability"
)

// ProcessApplicants is the canonical HR demo: user asks for a hire,
// Conductor decomposes to a JD, HR matches against its applicant store,
// scenario proposes proactive next actions (schedule interviews, post
// listing). End-to-end without an LLM call — deterministic + testable.
type ProcessApplicants struct{}

func (ProcessApplicants) Name() string { return "process_applicants" }
func (ProcessApplicants) Description() string {
	return "Decompose a hiring request into a JD, match against HR applicant store, suggest next steps."
}
func (ProcessApplicants) Example() map[string]any {
	return map[string]any{
		"request": "Hire a Site Supervisor for the Highway Bridge Project",
	}
}

func (ProcessApplicants) Run(ctx Context) (Result, error) {
	res := Result{Scenario: "process_applicants", Scope: ctx.Scope}

	// 1. Conductor classifies the request into a JD object. Input is a
	// free-form request OR a fully-formed JD.
	jd, classifyNote := buildJD(ctx.Input)
	classify := Step{
		Agent: "conductor", Action: "decompose",
		Input:  map[string]any{"request": ctx.Input["request"], "jd_provided": ctx.Input["jd"] != nil},
		Output: map[string]any{"jd": jd},
		OK:     true, Note: classifyNote,
	}
	res.Steps = append(res.Steps, classify)
	HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, classify)

	// 2. HR runs Match against its store.
	hr, err := capability.NewHR(ctx.DataRoot, ctx.Scope)
	if err != nil {
		return res, err
	}
	matches, err := hr.Match(jd, 5)
	if err != nil {
		return res, err
	}
	hrStep := Step{
		Agent: "hr", Action: "match_against_jd",
		Input:  map[string]any{"jd": jd},
		Output: map[string]any{"matches": matches, "count": len(matches)},
		OK:     true,
	}
	if len(matches) > 0 {
		hrStep.Note = fmt.Sprintf("found %d candidate(s); top: %s (score %d)",
			len(matches), matches[0].Applicant.Name, matches[0].Score)
	} else {
		hrStep.Note = "no candidates matched the JD"
	}
	res.Steps = append(res.Steps, hrStep)
	HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, hrStep)

	// 3. Proactive suggestions based on the matches.
	if len(matches) > 0 {
		top := matches[0]
		res.Summary = fmt.Sprintf("HR found %d candidate(s) for %q. Top: %s (score %d, %d matched / %d missing).",
			len(matches), jd.Title, top.Applicant.Name, top.Score,
			len(top.Matched), len(top.Missing))

		ids := make([]string, 0, len(matches))
		for _, m := range matches {
			ids = append(ids, m.Applicant.ID)
		}
		res.Suggestions = append(res.Suggestions, Suggestion{
			Title:  "Schedule interviews with the top candidates",
			Detail: fmt.Sprintf("Send interview invites to %d applicant(s).", len(matches)),
			Action: "schedule_interviews",
			Payload: map[string]any{
				"applicant_ids": ids,
				"jd_title":      jd.Title,
			},
		})
		res.Suggestions = append(res.Suggestions, Suggestion{
			Title:  "Shortlist the top match",
			Detail: fmt.Sprintf("Mark %s as shortlisted.", top.Applicant.Name),
			Action: "update_applicant_status",
			Payload: map[string]any{
				"applicant_id": top.Applicant.ID,
				"status":       "shortlisted",
			},
		})
		if len(top.Missing) > 0 {
			res.Suggestions = append(res.Suggestions, Suggestion{
				Title:  "Verify missing requirements with top candidate",
				Detail: fmt.Sprintf("Top match is missing: %s", strings.Join(top.Missing, ", ")),
				Action: "send_followup_email",
				Payload: map[string]any{
					"applicant_id": top.Applicant.ID,
					"questions":    top.Missing,
				},
			})
		}
	} else {
		res.Summary = fmt.Sprintf("No applicants in the store match %q.", jd.Title)
		res.Suggestions = append(res.Suggestions, Suggestion{
			Title:   "Post the job listing externally",
			Detail:  "No internal pipeline candidates; post to job boards.",
			Action:  "post_job_listing",
			Payload: map[string]any{"jd": jd},
		})
	}

	res.OK = true
	return res, nil
}

func buildJD(input map[string]any) (capability.JD, string) {
	// Path 1: caller passed a structured JD object.
	if raw, ok := input["jd"].(map[string]any); ok {
		jd := capability.JD{}
		if t, _ := raw["title"].(string); t != "" {
			jd.Title = t
		}
		if mh, ok := raw["must_have"].([]any); ok {
			for _, v := range mh {
				if s, _ := v.(string); s != "" {
					jd.MustHave = append(jd.MustHave, s)
				}
			}
		}
		if nh, ok := raw["nice_to_have"].([]any); ok {
			for _, v := range nh {
				if s, _ := v.(string); s != "" {
					jd.NiceToHave = append(jd.NiceToHave, s)
				}
			}
		}
		if y, ok := raw["min_years"].(float64); ok {
			jd.MinYears = int(y)
		}
		return jd, "JD provided directly by caller"
	}
	// Path 2: derive a JD from a free-form request. Deterministic, dumb but
	// transparent: pulls a title after "for/as" and uses the whole text as a
	// keyword bag for must-haves.
	req, _ := input["request"].(string)
	lower := strings.ToLower(req)
	title := req
	for _, sep := range []string{" for ", " as ", " — ", ": "} {
		if i := strings.Index(lower, sep); i > 0 {
			title = strings.TrimSpace(req[i+len(sep):])
			break
		}
	}
	jd := capability.JD{Title: title}
	// Crude keyword extraction — keywords that commonly appear in EPC JDs.
	for _, kw := range []string{
		"civil", "structural", "electrical", "mechanical", "site supervisor",
		"bridge", "highway", "solar", "concrete", "rebar", "pe", "pmp",
		"safety", "qa", "qc", "hse", "rfi", "primavera", "autocad",
	} {
		if strings.Contains(lower, kw) {
			jd.MustHave = append(jd.MustHave, kw)
		}
	}
	if strings.Contains(lower, "senior") || strings.Contains(lower, "lead") {
		jd.MinYears = 8
	} else if strings.Contains(lower, "junior") {
		jd.MinYears = 2
	} else {
		jd.MinYears = 5
	}
	return jd, "JD derived from free-form request via keyword extraction"
}
