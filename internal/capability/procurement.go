package capability

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Vendor — known supplier with a verified domain. The capability layer uses
// the domain mismatch test as the impersonation signal (matches the
// existing trust_reason wording in the email engine).
type Vendor struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Domain        string   `json:"domain"`         // canonical email domain, e.g. "vendorxyz.com"
	TrustedEmails []string `json:"trusted_emails"` // optional explicit allowlist
	ContactPerson string   `json:"contact_person,omitempty"`
	ContractActive bool    `json:"contract_active"`
}

// InvoiceCheck is the structured Procurement decision for a single invoice
// based on its sender vs the known vendor record.
type InvoiceCheck struct {
	VendorID          string   `json:"vendor_id"`
	VendorName        string   `json:"vendor_name,omitempty"`
	SenderEmail       string   `json:"sender_email"`
	Trusted           bool     `json:"trusted"`
	ImpersonationRisk bool     `json:"impersonation_risk"`
	Reasons           []string `json:"reasons"`
	Recommendation    string   `json:"recommendation"`
}

// Procurement is the per-scope vendor store + invoice validator.
type Procurement struct {
	mu   sync.RWMutex
	path string
}

func NewProcurement(dataRoot, scope string) (*Procurement, error) {
	dir := filepath.Join(dataRoot, scope)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return &Procurement{path: filepath.Join(dir, "vendors.json")}, nil
}

func (p *Procurement) load() ([]Vendor, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	data, err := os.ReadFile(p.path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Vendor{}, nil
		}
		return nil, err
	}
	var out []Vendor
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("parse %s: %w", p.path, err)
	}
	return out, nil
}

func (p *Procurement) save(in []Vendor) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	data, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p.path, data, 0o644)
}

func (p *Procurement) List() ([]Vendor, error) { return p.load() }

func (p *Procurement) Store(v Vendor) (Vendor, error) {
	if v.ID == "" {
		return v, fmt.Errorf("vendor id required")
	}
	cur, err := p.load()
	if err != nil {
		return v, err
	}
	replaced := false
	for i := range cur {
		if cur[i].ID == v.ID {
			cur[i] = v
			replaced = true
			break
		}
	}
	if !replaced {
		cur = append(cur, v)
	}
	if err := p.save(cur); err != nil {
		return v, err
	}
	return v, nil
}

func (p *Procurement) Get(id string) (*Vendor, error) {
	cur, err := p.load()
	if err != nil {
		return nil, err
	}
	for i := range cur {
		if cur[i].ID == id {
			return &cur[i], nil
		}
	}
	return nil, fmt.Errorf("vendor %q not found", id)
}

// ValidateInvoice compares the sender's email domain against the vendor's
// canonical domain + trusted email allowlist. This is the classic BEC
// pattern check the existing email engine flags via trust_status.
func (p *Procurement) ValidateInvoice(senderEmail, vendorID string, amount float64) (InvoiceCheck, error) {
	out := InvoiceCheck{
		VendorID:    vendorID,
		SenderEmail: senderEmail,
	}
	if vendorID == "" {
		out.Reasons = append(out.Reasons, "no vendor specified")
		out.Recommendation = "ask_procurement_to_classify"
		return out, nil
	}
	v, err := p.Get(vendorID)
	if err != nil {
		out.Reasons = append(out.Reasons, fmt.Sprintf("vendor %q not in registry", vendorID))
		out.Recommendation = "verify_vendor_identity"
		return out, nil
	}
	out.VendorName = v.Name

	senderDomain := domainOf(senderEmail)
	switch {
	case senderDomain == "":
		out.Reasons = append(out.Reasons, "sender email malformed")
		out.Recommendation = "reject_and_query"
	case stringSliceContainsFold(v.TrustedEmails, senderEmail):
		out.Trusted = true
		out.Reasons = append(out.Reasons, fmt.Sprintf("sender on %s allowlist", v.Name))
		out.Recommendation = "match_to_po"
	case strings.EqualFold(senderDomain, v.Domain):
		out.Trusted = true
		out.Reasons = append(out.Reasons,
			fmt.Sprintf("domain %s matches registered vendor domain", senderDomain))
		out.Recommendation = "match_to_po"
	default:
		out.ImpersonationRisk = true
		out.Reasons = append(out.Reasons,
			fmt.Sprintf("domain mismatch: sender @%s but %s is registered @%s",
				senderDomain, v.Name, v.Domain))
		if isPublicDomain(senderDomain) {
			out.Reasons = append(out.Reasons,
				fmt.Sprintf("sender uses public mail provider (%s) — classic BEC pattern", senderDomain))
		}
		if amount >= 50000 {
			out.Reasons = append(out.Reasons,
				fmt.Sprintf("amount $%.0f exceeds $50k human-review threshold", amount))
		}
		out.Recommendation = "block_vendor_and_notify_finance"
	}
	// An inactive contract overrides the recommendation even when the sender
	// is genuine: paying a trusted-but-out-of-contract vendor is itself a
	// procurement-process violation. Surface the contract status as the
	// gating issue rather than a passive note.
	if !v.ContractActive {
		out.Reasons = append(out.Reasons, "vendor contract not active")
		if out.Trusted {
			// Keep Trusted=true (sender IS legitimate) but reclassify the
			// recommendation so finance doesn't blindly schedule payment.
			out.Recommendation = "pause_payment_until_contract_renewed"
		}
	}
	return out, nil
}

func domainOf(email string) string {
	at := strings.LastIndex(email, "@")
	if at < 0 || at == len(email)-1 {
		return ""
	}
	return strings.ToLower(email[at+1:])
}

func isPublicDomain(d string) bool {
	switch strings.ToLower(d) {
	case "gmail.com", "yahoo.com", "outlook.com", "hotmail.com", "proton.me", "icloud.com":
		return true
	}
	return false
}

func stringSliceContainsFold(s []string, t string) bool {
	for _, v := range s {
		if strings.EqualFold(v, t) {
			return true
		}
	}
	return false
}
