package scenario

import (
	"fmt"
	"strings"

	"pty-claude-test/internal/capability"
)

// ProcessInbox is the most "multi-agent" scenario of the set. Conductor
// reads the global inbox, fans every unread email out to the agent that
// owns its category, asks each agent for a quick verdict + suggested
// next step, and returns a per-email triage.
type ProcessInbox struct{}

func (ProcessInbox) Name() string { return "process_inbox" }
func (ProcessInbox) Description() string {
	return "Sweep the inbox; fan vendor / applicant / client / internal mail out to the right agent; return per-email triage with proactive next steps."
}
func (ProcessInbox) Example() map[string]any {
	return map[string]any{"max": 8}
}

func (ProcessInbox) Run(ctx Context) (Result, error) {
	res := Result{Scenario: "process_inbox", Scope: ctx.Scope}

	max := 15 // triage the full demo inbox in one sweep; caller can override via "max"
	if v, ok := ctx.Input["max"].(float64); ok && v > 0 {
		max = int(v)
	}

	// 1. Conductor reads the global inbox.
	inbox := capability.NewInbox(ctx.DataRoot)
	unread, err := inbox.Unread()
	if err != nil {
		return res, err
	}
	if len(unread) > max {
		unread = unread[:max]
	}
	summary, _ := inbox.Summary()

	classify := Step{
		Agent: "conductor", Action: "sweep_inbox",
		Input:  map[string]any{"max": max},
		Output: map[string]any{"unread": len(unread), "summary": summary},
		OK:     true,
		Note:   fmt.Sprintf("%d unread · %d impersonations · %d new contacts", summary.Total, summary.Impersonations, summary.NewContacts),
	}
	res.Steps = append(res.Steps, classify)
	HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, classify)

	if len(unread) == 0 {
		res.Summary = "Inbox is clean — no unread mail."
		res.Suggestions = append(res.Suggestions, Suggestion{
			Title: "Audit silenced senders", Detail: "Confirm no senders are silently filtered.",
			Action: "audit_filters",
		})
		res.OK = true
		return res, nil
	}

	// 2. Route each email to the agent that owns its category. The agents
	// produce a per-email verdict; their messages stream into the hub so
	// each one pulses in the workforce panel.
	proc, _ := capability.NewProcurement(ctx.DataRoot, ctx.Scope)
	hr, _ := capability.NewHR(ctx.DataRoot, ctx.Scope)

	type Triage struct {
		ID, From, Subject, Category, Verdict, Suggestion string
		Severity                                         string // info|warn|critical
	}
	triages := []Triage{}
	stored := 0
	impersonations := 0
	clientFlags := 0

	for _, e := range unread {
		t := Triage{
			ID: e.ID, From: e.From, Subject: e.Subject, Category: e.Category,
			Severity: "info",
		}
		switch e.Category {
		case "vendor":
			// Trust the ingestion-time trust check first: the email engine already
			// flagged impersonation (domain mismatch / public-mail BEC) when the
			// message was received. This is authoritative and catches the fake or
			// unknown vendor_id that a registry-only ValidateInvoice check misses.
			if e.TrustStatus == "impersonation" {
				impersonations++
				t.Severity = "critical"
				reason := e.TrustReason
				if reason == "" {
					reason = "sender domain does not match the claimed vendor (possible BEC)"
				}
				if strings.HasPrefix(strings.ToUpper(reason), "IMPERSONATION") {
					t.Verdict = reason
				} else {
					t.Verdict = "IMPERSONATION RISK — " + reason
				}
				t.Suggestion = "block_vendor_and_notify_finance"
			} else if proc != nil {
				// Procurement does a domain check against the registry.
				check, _ := proc.ValidateInvoice(e.From, e.VendorID, 0)
				if check.ImpersonationRisk {
					impersonations++
					t.Severity = "critical"
					t.Verdict = fmt.Sprintf("IMPERSONATION RISK — %s", strings.Join(check.Reasons, "; "))
					t.Suggestion = "block_vendor_and_notify_finance"
				} else if check.Trusted {
					t.Verdict = fmt.Sprintf("trusted vendor (%s)", check.VendorName)
					t.Suggestion = "match_to_po_or_release_payment"
				} else {
					t.Severity = "warn"
					t.Verdict = "sender not in vendor registry"
					t.Suggestion = "onboard_vendor"
				}
			}
			HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, Step{
				Agent: "procurement", Action: "triage_email", OK: true,
				Input:  map[string]any{"email_id": e.ID, "from": e.From},
				Output: map[string]any{"verdict": t.Verdict, "severity": t.Severity},
				Note:   t.Verdict,
			})

		case "applicant":
			// HR stores the applicant if the email looks like an application.
			if hr != nil {
				name := e.FromName
				if name == "" {
					name = e.From
				}
				// Derive a DETERMINISTIC id from the source email so re-sweeping the
				// same inbox replaces the record instead of minting a new one each
				// time (an empty id => new app_<nanotime> on every run => runaway
				// duplicates). Idempotent: one application email => one applicant.
				applied := strings.TrimSpace(strings.TrimPrefix(e.Subject, "Application —"))
				applied = strings.TrimSpace(strings.TrimPrefix(applied, "Application -"))
				if applied == "" {
					applied = e.Subject
				}
				_, _ = hr.Store(capability.Applicant{
					ID:            "app_email_" + e.ID,
					Name:          name,
					Email:         e.From,
					AppliedFor:    applied,
					ResumeEmailID: e.ID,
				})
				stored++
				t.Verdict = fmt.Sprintf("added %s to pipeline", name)
				t.Suggestion = "shortlist_or_match_against_open_jd"
			}
			HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, Step{
				Agent: "hr", Action: "triage_email", OK: true,
				Input:  map[string]any{"email_id": e.ID, "from": e.From},
				Output: map[string]any{"verdict": t.Verdict},
				Note:   t.Verdict,
			})

		case "client":
			clientFlags++
			t.Severity = "warn"
			t.Verdict = "client-facing — needs response"
			t.Suggestion = "draft_reply_for_review"
			HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, Step{
				Agent: "project_manager", Action: "triage_email", OK: true,
				Input:  map[string]any{"email_id": e.ID, "from": e.From},
				Note:   "flagged for client response",
			})

		default:
			t.Verdict = "internal · no action"
			t.Suggestion = "archive"
		}
		triages = append(triages, t)
	}

	// 3. Roll-up step that summarises across agents.
	rollup := Step{
		Agent: "conductor", Action: "rollup",
		Output: map[string]any{
			"triaged":         len(triages),
			"impersonations":  impersonations,
			"applicants_stored": stored,
			"client_flags":    clientFlags,
		},
		OK:   true,
		Note: fmt.Sprintf("triaged %d · %d impersonations · %d new applicants · %d client flags", len(triages), impersonations, stored, clientFlags),
	}
	res.Steps = append(res.Steps, rollup)
	HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, rollup)

	// Expose the per-email table on the second step so the ScenarioCard's
	// shape-driven enhancer picks it up generically.
	res.Steps[len(res.Steps)-1].Output["triages"] = triages

	res.Summary = fmt.Sprintf("Triaged %d unread email(s). %d impersonation flag(s), %d new applicant(s), %d client flag(s).",
		len(triages), impersonations, stored, clientFlags)

	if impersonations > 0 {
		res.Suggestions = append(res.Suggestions, Suggestion{
			Title:  fmt.Sprintf("Block %d impersonation sender(s)", impersonations),
			Detail: "Add to vendor blocklist; alert finance.",
			Action: "block_impersonators",
		})
	}
	if clientFlags > 0 {
		res.Suggestions = append(res.Suggestions, Suggestion{
			Title:  fmt.Sprintf("Draft %d client reply(s) for review", clientFlags),
			Detail: "PM to draft progress/RFI responses; you approve before send.",
			Action: "draft_client_replies",
		})
	}
	if stored > 0 {
		res.Suggestions = append(res.Suggestions, Suggestion{
			Title:  fmt.Sprintf("Match %d new applicant(s) against open JDs", stored),
			Detail: "Auto-rank against any open roles in HR.",
			Action: "rank_new_applicants",
		})
	}
	if len(res.Suggestions) == 0 {
		res.Suggestions = append(res.Suggestions, Suggestion{
			Title: "All triaged · nothing urgent", Detail: "Inbox is processed.", Action: "no_op",
		})
	}
	res.OK = true
	return res, nil
}
