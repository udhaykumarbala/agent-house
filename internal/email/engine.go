package email

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Engine manages email inbox, trust verification, and applicant storage.
type Engine struct {
	mu         sync.RWMutex
	dataDir    string // project data directory
	inbox      []*Email
	vendors    []*Vendor
	applicants []*Applicant
	listeners  []func(*Email) // called when new email arrives
	blocklist  *Blocklist     // gateway-level deny list (exacts/domains/patterns) + BEC content rules
}

// NewEngine creates an email engine rooted at the given data directory.
func NewEngine(dataDir string) *Engine {
	e := &Engine{dataDir: dataDir}
	os.MkdirAll(filepath.Join(dataDir, "inbox"), 0755)
	os.MkdirAll(filepath.Join(dataDir, "outbox"), 0755)
	e.blocklist = NewBlocklist(dataDir)
	e.loadVendors()
	e.loadInbox()
	e.loadApplicants()
	return e
}

// OnEmail registers a callback for new incoming emails.
func (e *Engine) OnEmail(fn func(*Email)) {
	e.mu.Lock()
	e.listeners = append(e.listeners, fn)
	e.mu.Unlock()
}

// ReceiveEmail processes an incoming email: classify, trust-check, store.
func (e *Engine) ReceiveEmail(from, fromName, to, subject, body string) (*Email, *TrustResult) {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Gateway-level security check. Runs FIRST so a blocked sender
	// never reaches classify/trust/store/listeners. The blocklist is
	// also responsible for auditing the drop to blocked.jsonl.
	if block := e.evaluateSecurity(from, fromName, subject, body); block != nil {
		log.Printf("[EMAIL] BLOCKED at gateway: from=%s subject=%q mode=%s match=%q reason=%s",
			from, subject, block.BlockMode, block.BlockMatch, block.Reason)
		return nil, block
	}

	email := &Email{
		ID:        fmt.Sprintf("email_%d", time.Now().UnixMilli()),
		From:      from,
		FromName:  fromName,
		To:        to,
		Subject:   subject,
		Body:      body,
		Date:      time.Now(),
		Read:      false,
		Direction: "inbound",
	}

	// Classify the email
	email.Category = e.classify(subject, body, from)

	// Trust check for vendor-related emails
	trust := e.checkTrust(from, fromName, subject, body)
	email.TrustStatus = trust.Status
	email.TrustReason = trust.Reason
	email.VendorID = trust.VendorID

	// Store in inbox
	e.inbox = append(e.inbox, email)
	e.saveEmail(email)

	// If job application, parse and store applicant
	if email.Category == "job_application" {
		e.parseApplicant(email)
	}

	log.Printf("[EMAIL] Received: from=%s subject=%q category=%s trust=%s",
		from, subject, email.Category, trust.Status)

	// Notify listeners
	for _, fn := range e.listeners {
		go fn(email)
	}

	return email, trust
}

// evaluateSecurity is the gateway deny-list check. Returns a populated
// *TrustResult with Blocked=true when the email should be dropped, or
// nil when the message is allowed through to classify/trust/store.
//
// Order of checks (each one fires independently and short-circuits):
//  1. Blocklist exact address (e.g. ahmed.r@gmail.com)
//  2. Blocklist domain (e.g. any *@gmail.com — only if an operator
//     added that rule; the default list is exact-only to keep the
//     false-positive rate at zero on legitimate gmail senders)
//  3. Blocklist address-pattern
//  4. Built-in BEC content patterns (urgent-bank-change, etc.)
//  5. Display-name spoofing (public-domain sender claiming a senior
//     exec/finance role in the display name)
func (e *Engine) evaluateSecurity(from, fromName, subject, body string) *TrustResult {
	if e.blocklist == nil {
		return nil
	}
	if entry, ok := e.blocklist.MatchAddress(from); ok {
		reason := entry.Reason
		if reason == "" {
			reason = fmt.Sprintf("sender %q on blocklist (%s)", from, entry.Mode)
		}
		e.blocklist.Audit(BlockedEvent{
			From: from, FromName: fromName, Subject: subject,
			Reason: reason, Mode: entry.Mode, Match: entry.Value,
		})
		return &TrustResult{
			Status:     "blocked",
			Reason:     reason,
			Blocked:    true,
			BlockMode:  entry.Mode,
			BlockMatch: entry.Value,
		}
	}
	if name, snippet := e.blocklist.MatchContent(subject, body); name != "" {
		reason := fmt.Sprintf("BEC content filter (%s) fired: %q", name, snippet)
		e.blocklist.Audit(BlockedEvent{
			From: from, FromName: fromName, Subject: subject,
			Reason: reason, Mode: "bec_filter", Match: name,
		})
		return &TrustResult{
			Status:     "blocked",
			Reason:     reason,
			Blocked:    true,
			BlockMode:  "bec_filter",
			BlockMatch: name,
		}
	}
	if role := e.blocklist.MatchDisplayName(from, fromName); role != "" {
		reason := fmt.Sprintf("display-name spoof: public-mail sender %q claims role %q in display name", from, role)
		e.blocklist.Audit(BlockedEvent{
			From: from, FromName: fromName, Subject: subject,
			Reason: reason, Mode: "display_name_spoof", Match: role,
		})
		return &TrustResult{
			Status:     "blocked",
			Reason:     reason,
			Blocked:    true,
			BlockMode:  "display_name_spoof",
			BlockMatch: role,
		}
	}
	return nil
}

// GetInbox returns all inbox emails.
func (e *Engine) GetInbox() []*Email {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.inbox
}

// GetEmail returns a single email by ID.
func (e *Engine) GetEmail(id string) *Email {
	e.mu.RLock()
	defer e.mu.RUnlock()
	for _, email := range e.inbox {
		if email.ID == id {
			email.Read = true
			return email
		}
	}
	return nil
}

// GetApplicants returns all stored applicants.
func (e *Engine) GetApplicants() []*Applicant {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.applicants
}

// MarkRead marks an email as read.
func (e *Engine) MarkRead(id string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	for _, email := range e.inbox {
		if email.ID == id {
			email.Read = true
			e.saveEmail(email)
			return
		}
	}
}

// MarkReplied marks an email as replied and read.
func (e *Engine) MarkReplied(id string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	for _, email := range e.inbox {
		if email.ID == id {
			email.Read = true
			email.Replied = true
			email.RepliedAt = time.Now().Format(time.RFC3339)
			e.saveEmail(email)
			log.Printf("[EMAIL] Marked %s as replied", id)
			return
		}
	}
}

// DeleteEmail removes an email from inbox.
func (e *Engine) DeleteEmail(id string) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	for i, email := range e.inbox {
		if email.ID == id {
			e.inbox = append(e.inbox[:i], e.inbox[i+1:]...)
			path := filepath.Join(e.dataDir, "inbox", id+".json")
			os.Remove(path)
			log.Printf("[EMAIL] Deleted email %s", id)
			return true
		}
	}
	return false
}

// UpdateApplicantStatus updates an applicant's status.
func (e *Engine) UpdateApplicantStatus(applicantID, status string) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	for _, a := range e.applicants {
		if a.ID == applicantID {
			a.Status = status
			e.saveApplicants()
			log.Printf("[EMAIL] Applicant %s status → %s", applicantID, status)
			return true
		}
	}
	return false
}

// AddTrustedEmail adds an email to a vendor's trusted list.
func (e *Engine) AddTrustedEmail(vendorID, emailAddr string) bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, v := range e.vendors {
		if v.ID == vendorID {
			v.TrustedEmails = append(v.TrustedEmails, emailAddr)
			e.saveVendors()
			log.Printf("[EMAIL] Added %s to vendor %s trusted list", emailAddr, vendorID)
			return true
		}
	}
	return false
}

// BlockSender adds a sender / domain / pattern to the gateway blocklist.
func (e *Engine) BlockSender(mode, value, reason, addedBy string) error {
	return e.blocklist.Add(BlockedEntry{
		Mode:    mode,
		Value:   value,
		Reason:  reason,
		AddedBy: addedBy,
	})
}

// UnblockSender removes an entry from the blocklist.
func (e *Engine) UnblockSender(mode, value string) error {
	return e.blocklist.Remove(mode, value)
}

// ListBlocked returns a copy of the current blocklist.
func (e *Engine) ListBlocked() []BlockedEntry {
	return e.blocklist.List()
}

// classify determines the email category.
func (e *Engine) classify(subject, body, from string) string {
	lower := strings.ToLower(subject + " " + body)

	// Job application keywords
	if containsAny(lower, "resume", "application", "apply", "cv ", "curriculum vitae", "job opening", "position") {
		return "job_application"
	}

	// Invoice keywords
	if containsAny(lower, "invoice", "payment", "billing", "amount due") {
		return "vendor"
	}

	// Check if from a known vendor domain
	for _, v := range e.vendors {
		if strings.HasSuffix(from, "@"+v.Domain) {
			return "vendor"
		}
	}

	// Client keywords
	if containsAny(lower, "update", "progress", "status", "deliverable", "milestone") {
		return "client"
	}

	// Internal keywords
	if containsAny(lower, "internal", "team", "resource", "leave", "timesheet") {
		return "internal"
	}

	return "unknown"
}

// checkTrust verifies the sender against vendor trusted email lists.
func (e *Engine) checkTrust(from, fromName, subject, body string) *TrustResult {
	fromDomain := extractDomain(from)
	lower := strings.ToLower(subject + " " + body + " " + fromName)

	// Check each vendor
	for _, v := range e.vendors {
		// Direct match: sender is in trusted list
		for _, trusted := range v.TrustedEmails {
			if strings.EqualFold(from, trusted) {
				return &TrustResult{
					Status:     "trusted",
					Reason:     fmt.Sprintf("Sender matches %s's trusted email list", v.Name),
					VendorID:   v.ID,
					VendorName: v.Name,
				}
			}
		}

		// Domain match but not in trusted list
		if strings.EqualFold(fromDomain, v.Domain) {
			return &TrustResult{
				Status:     "new_contact",
				Reason:     fmt.Sprintf("New email from %s's domain (%s) — not yet in trusted list", v.Name, v.Domain),
				VendorID:   v.ID,
				VendorName: v.Name,
			}
		}

		// Impersonation: mentions vendor name but different domain
		vendorNameLower := strings.ToLower(v.Name)
		if strings.Contains(lower, vendorNameLower) && !strings.EqualFold(fromDomain, v.Domain) {
			risks := []string{
				fmt.Sprintf("Sender domain (%s) does not match vendor domain (%s)", fromDomain, v.Domain),
				fmt.Sprintf("Known contact: %s — sender: %s", v.ContactPerson, fromName),
			}
			if containsAny(lower, "payment", "bank", "wire", "transfer", "account details", "urgent") {
				risks = append(risks, "Email contains payment/banking language — possible BEC fraud")
			}
			return &TrustResult{
				Status:      "impersonation",
				Reason:      fmt.Sprintf("IMPERSONATION RISK — claims association with %s but sending from %s", v.Name, fromDomain),
				RiskFactors: risks,
				VendorID:    v.ID,
				VendorName:  v.Name,
			}
		}
	}

	// Not vendor-related
	return &TrustResult{
		Status: "trusted",
		Reason: "Non-vendor email — no trust check required",
	}
}

// parseApplicant extracts applicant info from a job application email.
func (e *Engine) parseApplicant(email *Email) {
	applicant := &Applicant{
		ID:            fmt.Sprintf("app_%d", time.Now().UnixMilli()),
		Name:          email.FromName,
		Email:         email.From,
		AppliedFor:    extractJobTitle(email.Subject, email.Body),
		ResumeEmailID: email.ID,
		ReceivedDate:  email.Date.Format("2006-01-02"),
		Status:        "new",
	}

	e.applicants = append(e.applicants, applicant)
	e.saveApplicants()

	log.Printf("[EMAIL] Stored applicant: %s for %s", applicant.Name, applicant.AppliedFor)
}

// File I/O
func (e *Engine) loadVendors() {
	path := filepath.Join(e.dataDir, "vendors.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	json.Unmarshal(data, &e.vendors)
	log.Printf("[EMAIL] Loaded %d vendors", len(e.vendors))
}

func (e *Engine) saveVendors() {
	data, _ := json.MarshalIndent(e.vendors, "", "  ")
	os.WriteFile(filepath.Join(e.dataDir, "vendors.json"), data, 0644)
}

func (e *Engine) loadInbox() {
	dir := filepath.Join(e.dataDir, "inbox")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			continue
		}
		var email Email
		if json.Unmarshal(data, &email) == nil {
			e.inbox = append(e.inbox, &email)
		}
	}
	log.Printf("[EMAIL] Loaded %d emails from inbox", len(e.inbox))
}

func (e *Engine) saveEmail(email *Email) {
	data, _ := json.MarshalIndent(email, "", "  ")
	path := filepath.Join(e.dataDir, "inbox", email.ID+".json")
	os.WriteFile(path, data, 0644)
}

func (e *Engine) loadApplicants() {
	path := filepath.Join(e.dataDir, "applicants.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	json.Unmarshal(data, &e.applicants)
	log.Printf("[EMAIL] Loaded %d applicants", len(e.applicants))
}

func (e *Engine) saveApplicants() {
	data, _ := json.MarshalIndent(e.applicants, "", "  ")
	os.WriteFile(filepath.Join(e.dataDir, "applicants.json"), data, 0644)
}

// Helpers
func extractDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

func containsAny(text string, keywords ...string) bool {
	for _, kw := range keywords {
		if strings.Contains(text, kw) {
			return true
		}
	}
	return false
}

func extractJobTitle(subject, body string) string {
	lower := strings.ToLower(subject)
	// Try to extract role from subject
	for _, prefix := range []string{"application for ", "applying for ", "application — ", "resume — ", "cv — "} {
		if idx := strings.Index(lower, prefix); idx >= 0 {
			title := subject[idx+len(prefix):]
			if nl := strings.IndexAny(title, "\n\r"); nl > 0 {
				title = title[:nl]
			}
			return strings.TrimSpace(title)
		}
	}
	return "General Application"
}
