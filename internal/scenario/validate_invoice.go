package scenario

import (
	"fmt"

	"pty-claude-test/internal/capability"
)

// ValidateInvoice composes Conductor → Procurement to check whether an
// incoming invoice is legitimate, and proposes proactive next actions
// based on the outcome (block + notify on impersonation; match-to-PO on
// trusted).
type ValidateInvoice struct{}

func (ValidateInvoice) Name() string { return "validate_invoice" }
func (ValidateInvoice) Description() string {
	return "Validate a vendor invoice against the vendor registry; flag impersonation, suggest actions."
}
func (ValidateInvoice) Example() map[string]any {
	return map[string]any{
		"sender_email": "ahmed.r@gmail.com",
		"vendor_id":    "vendor_xyz",
		"amount":       340000,
	}
}

func (ValidateInvoice) Run(ctx Context) (Result, error) {
	res := Result{Scenario: "validate_invoice", Scope: ctx.Scope}

	sender, _ := ctx.Input["sender_email"].(string)
	vendorID, _ := ctx.Input["vendor_id"].(string)
	amount := readFloat(ctx.Input, "amount")

	classify := Step{
		Agent: "conductor", Action: "decompose",
		Input: map[string]any{"sender_email": sender, "vendor_id": vendorID, "amount": amount},
		OK:    true,
		Note:  "routing to Procurement for vendor verification",
	}
	res.Steps = append(res.Steps, classify)
	HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, classify)

	proc, err := capability.NewProcurement(ctx.DataRoot, ctx.Scope)
	if err != nil {
		return res, err
	}
	check, err := proc.ValidateInvoice(sender, vendorID, amount)
	if err != nil {
		return res, err
	}
	pStep := Step{
		Agent: "procurement", Action: "validate_invoice",
		Input:  map[string]any{"sender_email": sender, "vendor_id": vendorID, "amount": amount},
		Output: map[string]any{"check": check},
		OK:     true,
	}
	switch {
	case check.ImpersonationRisk:
		pStep.Note = "IMPERSONATION risk detected — domain mismatch"
	case check.Trusted:
		pStep.Note = fmt.Sprintf("trusted sender for %s", check.VendorName)
	default:
		pStep.Note = "unverified — recommend manual review"
	}
	res.Steps = append(res.Steps, pStep)
	HelpEmitStep(ctx.Emitter, ctx.Scope, res.Scenario, pStep)

	// Resolve a display name for the claimed vendor — VendorName is empty when the
	// sender couldn't be tied to a registry entry (e.g. a gmail BEC sender), so
	// fall back to the claimed id / a generic phrase rather than printing "".
	claimed := check.VendorName
	if claimed == "" {
		claimed = vendorID
	}
	if claimed == "" {
		claimed = "a known vendor"
	}
	switch {
	case check.ImpersonationRisk:
		res.Summary = fmt.Sprintf("Impersonation risk: %s purporting to be %s.",
			sender, claimed)
		res.Suggestions = append(res.Suggestions, Suggestion{
			Title:  "Block this sender",
			Detail: "Add the sender to the blocklist and refuse any payment changes.",
			Action: "block_vendor",
			Payload: map[string]any{
				"sender_email": sender, "vendor_id": vendorID,
			},
		})
		res.Suggestions = append(res.Suggestions, Suggestion{
			Title:  "Notify finance & PM",
			Detail: "Heads up on a likely BEC attempt; freeze pending payments.",
			Action: "notify_finance",
			Payload: map[string]any{
				"reasons": check.Reasons,
			},
		})
		res.Suggestions = append(res.Suggestions, Suggestion{
			Title:  "Verify via known channel",
			Detail: fmt.Sprintf("Call %s on a number on file to confirm — never reply to the suspicious email.", claimed),
			Action: "verify_oob",
		})
	case check.Trusted && check.Recommendation == "pause_payment_until_contract_renewed":
		res.Summary = fmt.Sprintf("Trusted vendor (%s) but contract is INACTIVE — pause payment.", check.VendorName)
		res.Suggestions = append(res.Suggestions, Suggestion{
			Title:   "Pause payment until contract is reactivated",
			Detail:  "Sender is legitimate, but paying without an active contract violates procurement policy.",
			Action:  "pause_payment",
			Payload: map[string]any{"vendor_id": vendorID, "amount": amount},
		})
		res.Suggestions = append(res.Suggestions, Suggestion{
			Title:   "Reactivate the vendor contract",
			Detail:  "Open the contract record; renew or revise terms before next payment cycle.",
			Action:  "reactivate_contract",
			Payload: map[string]any{"vendor_id": vendorID},
		})
		res.Suggestions = append(res.Suggestions, Suggestion{
			Title:   "Notify Procurement to follow up",
			Detail:  "Procurement to confirm contract status with the vendor.",
			Action:  "notify_procurement",
			Payload: map[string]any{"vendor_id": vendorID},
		})
	case check.Trusted:
		res.Summary = fmt.Sprintf("Trusted vendor (%s).", check.VendorName)
		res.Suggestions = append(res.Suggestions, Suggestion{
			Title:   "Match invoice to a PO",
			Detail:  "Locate the matching purchase order; flag any mismatch.",
			Action:  "match_to_po",
			Payload: map[string]any{"vendor_id": vendorID, "amount": amount},
		})
		res.Suggestions = append(res.Suggestions, Suggestion{
			Title:   "Schedule payment after PO match",
			Detail:  "Queue payment subject to PO + receipt confirmation.",
			Action:  "schedule_payment",
			Payload: map[string]any{"vendor_id": vendorID, "amount": amount},
		})
	default:
		res.Summary = "Sender not in vendor registry."
		res.Suggestions = append(res.Suggestions, Suggestion{
			Title:  "Have Procurement onboard the vendor",
			Detail: "Run KYC; collect bank details OOB before any payment.",
			Action: "onboard_vendor",
			Payload: map[string]any{"sender_email": sender},
		})
	}
	res.OK = true
	return res, nil
}

func readFloat(in map[string]any, key string) float64 {
	if v, ok := in[key].(float64); ok {
		return v
	}
	if v, ok := in[key].(int); ok {
		return float64(v)
	}
	return 0
}
