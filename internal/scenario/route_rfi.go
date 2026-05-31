package scenario

import (
	"fmt"
	"strings"

	"pty-claude-test/internal/capability"
)

// RouteRFI takes an incoming RFI (Request For Information) text, decides
// which discipline owns it, assigns the Site Engineer as the standing
// owner, and uses HR's applicant skill index to surface any matching
// specialists in the hiring pipeline. Suggests next steps based on whether
// a current specialist exists.
type RouteRFI struct{}

func (RouteRFI) Name() string { return "route_rfi" }
func (RouteRFI) Description() string {
	return "Classify an incoming RFI by discipline, assign to Site Engineer, check HR pipeline for matching specialists, propose next steps."
}
func (RouteRFI) Example() map[string]any {
	return map[string]any{
		"rfi_id":  "RFI-18",
		"subject": "Column tolerances on grid C-7",
		"body":    "Per drawing A-102 v3, max allowable plumb deviation on cast-in-place columns is ±10mm, but the structural spec calls out ±6mm. Please confirm the governing tolerance before pour.",
	}
}

// disciplineKeywords drives the routing decision. Each entry maps a
// discipline label to keywords that strongly suggest it. First match wins;
// "general" is the fallback.
var disciplineKeywords = []struct {
	label    string
	keywords []string
}{
	{"structural", []string{"column", "beam", "rebar", "slab", "load-bearing", "tolerance", "plumb", "concrete"}},
	{"electrical", []string{"electrical", "wiring", "cable", "transformer", "switchgear", "earthing", "voltage", "amperage"}},
	{"mechanical", []string{"hvac", "ducting", "chiller", "pump", "boiler", "piping", "mep", "mechanical"}},
	{"civil", []string{"foundation", "excavation", "drainage", "earthwork", "soil", "geotech", "highway", "bridge", "roadway"}},
	{"safety", []string{"safety", "ppe", "hazard", "accident", "incident", "fall protection", "scaffold"}},
	{"qa", []string{"inspection", "checklist", "punchlist", "qa", "qc", "test report", "non-conformance", "ncr"}},
}

func (RouteRFI) Run(ctx Context) (Result, error) {
	res := Result{Scenario: "route_rfi", Scope: ctx.Scope}

	subject, _ := ctx.Input["subject"].(string)
	body, _ := ctx.Input["body"].(string)
	rfiID, _ := ctx.Input["rfi_id"].(string)
	text := strings.ToLower(subject + " " + body)

	// 1. Conductor classifies discipline.
	discipline, hits := classifyDiscipline(text)
	classify := Step{
		Agent: "conductor", Action: "classify_discipline",
		Input: map[string]any{
			"rfi_id":  rfiID,
			"subject": subject,
		},
		Output: map[string]any{
			"discipline":     discipline,
			"matched_terms":  hits,
		},
		OK: true,
		Note: fmt.Sprintf("classified as %s discipline%s",
			discipline,
			ifNonEmpty(hits, " (matched: "+strings.Join(hits, ", ")+")")),
	}
	res.Steps = append(res.Steps, classify)
	HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, classify)

	// 2. Site Engineer accepts the routing.
	site := Step{
		Agent: "site_engineer", Action: "accept_routing", OK: true,
		Input: map[string]any{"discipline": discipline},
		Note: fmt.Sprintf("standing owner for %s discipline on this site",
			discipline),
	}
	res.Steps = append(res.Steps, site)
	HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, site)

	// 3. HR checks the applicant pipeline for matching specialists. The JD
	//    is built from the discipline + matched terms — same Match function
	//    process_applicants uses, so the score scale is consistent.
	hr, _ := capability.NewHR(ctx.DataRoot, ctx.Scope)
	jd := capability.JD{
		Title:    discipline + " specialist",
		MustHave: append([]string{discipline}, hits...),
	}
	matches, _ := hr.Match(jd, 3)
	hrStep := Step{
		Agent: "hr", Action: "find_specialist", OK: true,
		Input:  map[string]any{"jd": jd},
		Output: map[string]any{"matches": matches, "count": len(matches)},
	}
	if len(matches) > 0 {
		hrStep.Note = fmt.Sprintf("found %d specialist(s) in pipeline; top: %s (score %d)",
			len(matches), matches[0].Applicant.Name, matches[0].Score)
	} else {
		hrStep.Note = "no matching specialist in the applicant pipeline"
	}
	res.Steps = append(res.Steps, hrStep)
	HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, hrStep)

	// Summary + proactive suggestions.
	res.Summary = fmt.Sprintf("RFI %s routed to %s discipline (Site Engineer assigned). %s",
		rfiID, discipline,
		ifElse(len(matches) > 0,
			fmt.Sprintf("%d candidate specialist(s) in pipeline.", len(matches)),
			"No pipeline specialists matched."))

	res.Suggestions = append(res.Suggestions, Suggestion{
		Title:  "Assign to Site Engineer",
		Detail: "Site Engineer is the standing on-site owner for this discipline.",
		Action: "assign_to_site_engineer",
		Payload: map[string]any{
			"rfi_id":     rfiID,
			"discipline": discipline,
		},
	})
	if len(matches) > 0 {
		top := matches[0]
		res.Suggestions = append(res.Suggestions, Suggestion{
			Title:  fmt.Sprintf("Loop in %s as specialist consultant", top.Applicant.Name),
			Detail: fmt.Sprintf("HR pipeline match (score %d). Could provide subject-matter review.", top.Score),
			Action: "consult_pipeline_specialist",
			Payload: map[string]any{
				"applicant_id": top.Applicant.ID,
				"rfi_id":       rfiID,
			},
		})
	} else {
		res.Suggestions = append(res.Suggestions, Suggestion{
			Title:  fmt.Sprintf("Hire a %s specialist for the gap", discipline),
			Detail: "No matching pipeline candidate. Open a JD or escalate to an external consultant.",
			Action: "open_jd_for_specialist",
			Payload: map[string]any{
				"discipline": discipline,
				"rfi_id":     rfiID,
			},
		})
	}
	res.Suggestions = append(res.Suggestions, Suggestion{
		Title:   "Draft an RFI response for review",
		Detail:  "Site Engineer drafts; you approve before send.",
		Action:  "draft_rfi_response",
		Payload: map[string]any{"rfi_id": rfiID},
	})

	res.OK = true
	return res, nil
}

func classifyDiscipline(textLower string) (string, []string) {
	for _, d := range disciplineKeywords {
		hits := []string{}
		for _, kw := range d.keywords {
			if strings.Contains(textLower, kw) {
				hits = append(hits, kw)
			}
		}
		if len(hits) > 0 {
			return d.label, hits
		}
	}
	return "general", nil
}

func ifNonEmpty(s []string, suffix string) string {
	if len(s) == 0 {
		return ""
	}
	return suffix
}

func ifElse(cond bool, a, b string) string {
	if cond {
		return a
	}
	return b
}
