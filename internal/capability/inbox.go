package capability

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Email mirrors the on-disk shape written by the email engine + the
// scripts/scenarios/load-epc-atlas.sh seed. The capability layer reads
// this directly instead of going through internal/email so it stays
// dependency-free (same shape, owned copy).
type Email struct {
	ID          string `json:"id"`
	From        string `json:"from"`
	FromName    string `json:"from_name,omitempty"`
	To          string `json:"to,omitempty"`
	Subject     string `json:"subject"`
	Body        string `json:"body,omitempty"`
	Date        string `json:"date,omitempty"`
	Read        bool   `json:"read,omitempty"`
	Replied     bool   `json:"replied,omitempty"`
	Category    string `json:"category"`     // vendor | applicant | client | internal
	VendorID    string `json:"vendor_id,omitempty"`
	TrustStatus string `json:"trust_status,omitempty"` // trusted | impersonation | new_contact
	TrustReason string `json:"trust_reason,omitempty"`
	Direction   string `json:"direction,omitempty"`
}

// Inbox is the global inbox reader. There is one inbox shared across scopes
// — that matches the production reality (one mailbox). It reads from
// projects/inbox/email_*.json the same files the email engine + the
// scenario seed scripts write.
type Inbox struct {
	path string
}

// NewInbox returns the global inbox reader. The on-disk path is fixed at
// projects/inbox/ to remain compatible with the existing email engine and
// the seed scripts; the dataRoot argument is unused but kept for parity
// with the other NewX constructors.
func NewInbox(_ string) *Inbox {
	return &Inbox{path: filepath.Join("projects", "inbox")}
}

func (i *Inbox) List() ([]Email, error) {
	entries, err := os.ReadDir(i.path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Email{}, nil
		}
		return nil, err
	}
	out := make([]Email, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		// Skip the applicants list and any non-email artefacts.
		if e.Name() == "applicants.json" || e.Name() == "vendors.json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(i.path, e.Name()))
		if err != nil {
			continue
		}
		var em Email
		if err := json.Unmarshal(data, &em); err != nil || em.ID == "" {
			continue
		}
		out = append(out, em)
	}
	sort.Slice(out, func(a, b int) bool { return out[a].Date > out[b].Date })
	return out, nil
}

// Inject writes a single email into the inbox directory as JSON. The
// caller (Lab UI / seed scripts) uses this to simulate inbound mail
// without going through a real SMTP/IMAP server. The file is named after
// the email's ID, so callers can pass a deterministic id (e.g.
// "mock_impersonation_<ts>") to keep the inbox directory readable.
func (i *Inbox) Inject(e Email) (Email, error) {
	if e.Subject == "" {
		return e, fmt.Errorf("subject required")
	}
	if e.From == "" {
		return e, fmt.Errorf("from required")
	}
	if e.ID == "" {
		e.ID = fmt.Sprintf("mock_%d", time.Now().UnixNano())
	}
	if e.Date == "" {
		e.Date = time.Now().UTC().Format(time.RFC3339)
	}
	if e.Direction == "" {
		e.Direction = "inbound"
	}
	if e.Category == "" {
		e.Category = "internal"
	}
	if err := os.MkdirAll(i.path, 0o755); err != nil {
		return e, err
	}
	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return e, err
	}
	path := filepath.Join(i.path, e.ID+".json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return e, err
	}
	return e, nil
}

func (i *Inbox) Unread() ([]Email, error) {
	all, err := i.List()
	if err != nil {
		return nil, err
	}
	out := []Email{}
	for _, e := range all {
		if !e.Read {
			out = append(out, e)
		}
	}
	return out, nil
}

// CategoryCounts groups unread emails by category and (for vendor) by trust.
type CategoryCounts struct {
	Total           int            `json:"total"`
	ByCategory      map[string]int `json:"by_category"`
	Impersonations  int            `json:"impersonations"`
	NewContacts     int            `json:"new_contacts"`
}

func (i *Inbox) Summary() (CategoryCounts, error) {
	out := CategoryCounts{ByCategory: map[string]int{}}
	all, err := i.List()
	if err != nil {
		return out, err
	}
	for _, e := range all {
		if e.Read {
			continue
		}
		out.Total++
		out.ByCategory[e.Category]++
		if e.TrustStatus == "impersonation" {
			out.Impersonations++
		}
		if e.TrustStatus == "new_contact" {
			out.NewContacts++
		}
	}
	return out, nil
}
