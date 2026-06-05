package email

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestNewBlocklist_SeedsDefault is the regression test for the reported
// incident: ahmed.r@gmail.com must be blocked out of the box, with no
// operator action required.
func TestNewBlocklist_SeedsDefault(t *testing.T) {
	dir := t.TempDir()
	b := NewBlocklist(dir)

	entry, ok := b.MatchAddress("ahmed.r@gmail.com")
	if !ok {
		t.Fatal("ahmed.r@gmail.com should be blocked by the seeded default")
	}
	if entry.Mode != "exact" {
		t.Errorf("default block mode = %q, want exact", entry.Mode)
	}
	if !strings.Contains(strings.ToLower(entry.Reason), "bec") {
		t.Errorf("default block reason should mention BEC, got %q", entry.Reason)
	}
}

// TestNewBlocklist_PersistsToDisk verifies the blocklist survives a
// process restart. Drops go to blocked.jsonl and the entries themselves
// live in blocklist.json.
func TestNewBlocklist_PersistsToDisk(t *testing.T) {
	dir := t.TempDir()

	b1 := NewBlocklist(dir)
	if err := b1.Add(BlockedEntry{Mode: "domain", Value: "spamco.example", Reason: "test"}); err != nil {
		t.Fatalf("Add: %v", err)
	}

	// New engine / blocklist on the same dir should see the same entries.
	b2 := NewBlocklist(dir)
	_, ok := b2.MatchAddress("anything@spamco.example")
	if !ok {
		t.Fatal("domain block should persist across NewBlocklist calls")
	}
}

// TestBlocklist_AddRemove covers the management API: add → match →
// remove → no-match, and bad input is rejected.
func TestBlocklist_AddRemove(t *testing.T) {
	b := NewBlocklist(t.TempDir())

	if err := b.Add(BlockedEntry{Mode: "domain", Value: "evil.example", Reason: "test"}); err != nil {
		t.Fatalf("Add domain: %v", err)
	}
	if _, ok := b.MatchAddress("ceo@evil.example"); !ok {
		t.Fatal("expected match after add")
	}
	if err := b.Remove("domain", "evil.example"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if _, ok := b.MatchAddress("ceo@evil.example"); ok {
		t.Fatal("expected no match after remove")
	}

	// Dedup: second add with same (mode, value) replaces, not duplicates.
	if err := b.Add(BlockedEntry{Mode: "exact", Value: "x@y.com", Reason: "first"}); err != nil {
		t.Fatal(err)
	}
	if err := b.Add(BlockedEntry{Mode: "exact", Value: "x@y.com", Reason: "second"}); err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, e := range b.List() {
		if e.Mode == "exact" && e.Value == "x@y.com" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("dedup failed, got %d entries for x@y.com", count)
	}

	// Bad pattern is rejected, never reaches the compiled list.
	if err := b.Add(BlockedEntry{Mode: "pattern", Value: "([unclosed"}); err == nil {
		t.Fatal("expected bad pattern to be rejected")
	}
}

// TestBlocklist_MatchAddress covers exact + domain + pattern matching.
func TestBlocklist_MatchAddress(t *testing.T) {
	b := NewBlocklist(t.TempDir())
	mustAdd := func(e BlockedEntry) {
		t.Helper()
		if err := b.Add(e); err != nil {
			t.Fatalf("Add %v: %v", e, err)
		}
	}
	mustAdd(BlockedEntry{Mode: "exact", Value: "ceo@scammer.test"})
	mustAdd(BlockedEntry{Mode: "domain", Value: "phish.example"})
	mustAdd(BlockedEntry{Mode: "pattern", Value: `^finance-.*@.+\.ru$`})

	cases := []struct {
		addr     string
		wantHit  bool
		wantMode string
	}{
		{"ceo@scammer.test", true, "exact"},
		{"CEO@SCAMMER.TEST", true, "exact"}, // case-insensitive
		{"any@phish.example", true, "domain"},
		{"any@sub.phish.example", false, ""}, // domain match is exact, not suffix
		{"finance-anna@bank.ru", true, "pattern"},
		{"finance-anna@bank.com", false, ""},
		{"legit@vendorxyz.com", false, ""}, // control: legitimate sender
	}
	for _, tc := range cases {
		t.Run(tc.addr, func(t *testing.T) {
			entry, ok := b.MatchAddress(tc.addr)
			if ok != tc.wantHit {
				t.Fatalf("MatchAddress(%q) ok=%v, want %v", tc.addr, ok, tc.wantHit)
			}
			if ok && entry.Mode != tc.wantMode {
				t.Errorf("MatchAddress(%q) mode=%q, want %q", tc.addr, entry.Mode, tc.wantMode)
			}
		})
	}
}

// TestBlocklist_MatchContent covers all five built-in BEC content rules
// with representative payloads. Each row is the smallest payload that
// should fire that rule.
func TestBlocklist_MatchContent(t *testing.T) {
	b := NewBlocklist(t.TempDir())
	cases := []struct {
		name         string
		subject      string
		body         string
		wantRule     string
		wantContains string // substring expected in the matched snippet
	}{
		{
			name:         "urgent-bank-change",
			subject:      "URGENT — wire instructions updated",
			body:         "Please update our banking details effective immediately for the new account.",
			wantRule:     "urgent-bank-change",
			wantContains: "bank",
		},
		{
			name:         "invoice-bank-change",
			subject:      "Invoice #1042 — payment details",
			body:         "Please use the new bank account on the attached PO for this invoice.",
			wantRule:     "invoice-bank-change",
			wantContains: "invoice",
		},
		{
			name:         "gift-card-scam",
			subject:      "Quick favor",
			body:         "Please buy 10 Apple gift cards and email me the codes. Do not tell anyone in finance.",
			wantRule:     "gift-card-scam",
			wantContains: "gift",
		},
		{
			name:         "exec-impersonation-favor",
			subject:      "Are you at your desk?",
			body:         "I need a favor — please handle this quietly before I get back to the office.",
			wantRule:     "exec-impersonation-favor",
			wantContains: "favor",
		},
		{
			name:         "phish-secure-link",
			subject:      "Action required: secure document",
			body:         "Please verify your account here: https://login-secure-bogus.example/auth",
			wantRule:     "phish-secure-link",
			wantContains: "verify",
		},
		// Controls — these should NOT fire any rule.
		{
			name:    "legitimate-invoice",
			subject: "Invoice #1042 for May services",
			body:    "Please remit via the account on file. Thanks for your business.",
		},
		{
			name:    "legitimate-meeting",
			subject: "Are you free for a 1:1 on Thursday?",
			body:    "Want to discuss the Q3 plan. Nothing urgent.",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rule, snippet := b.MatchContent(tc.subject, tc.body)
			if tc.wantRule == "" {
				if rule != "" {
					t.Fatalf("control payload fired rule %q on %q", rule, snippet)
				}
				return
			}
			if rule != tc.wantRule {
				t.Fatalf("rule=%q, want %q (snippet=%q)", rule, tc.wantRule, snippet)
			}
			if !strings.Contains(snippet, tc.wantContains) {
				t.Errorf("matched snippet %q does not contain %q", snippet, tc.wantContains)
			}
		})
	}
}

// TestBlocklist_MatchDisplayName verifies the CEO-fraud display-name
// spoofing rule: it only fires when BOTH the address is on a public
// mail provider AND the display name claims a senior role.
func TestBlocklist_MatchDisplayName(t *testing.T) {
	b := NewBlocklist(t.TempDir())
	cases := []struct {
		from, fromName string
		want           string // empty means no match
	}{
		// Should fire — public mail + CEO in display name.
		{"attacker@gmail.com", "John Smith, CEO", "ceo"},
		{"ceo.impersonator@outlook.com", "Sarah Lee — CFO", "cfo"},
		{"finance@proton.me", "Mark (Managing Director)", "managing director"},
		// Should NOT fire — public mail but no exec role claim.
		{"legit@gmail.com", "Jane Vendor", ""},
		// Should NOT fire — exec role in display name, but sender is
		// on a corporate domain (legitimate scenario, no spoofing).
		{"ceo@vendorxyz.com", "Jane CEO", ""},
		// Should NOT fire — exec role misspelled / not in the keyword
		// list. Better to miss than to false-positive on "supervisor".
		{"someone@gmail.com", "Project Supervisor", ""},
	}
	for _, tc := range cases {
		t.Run(tc.from+"__"+tc.fromName, func(t *testing.T) {
			got := b.MatchDisplayName(tc.from, tc.fromName)
			if got != tc.want {
				t.Errorf("MatchDisplayName(%q, %q) = %q, want %q",
					tc.from, tc.fromName, got, tc.want)
			}
		})
	}
}

// TestEngine_ReceiveEmail_Blocks at the gateway is the end-to-end
// integration: confirms that a blocked email is NOT stored in the
// inbox, NOT classified, and NOT saved to disk.
func TestEngine_ReceiveEmail_Blocks(t *testing.T) {
	dir := t.TempDir()
	e := NewEngine(dir)

	// The default seed should block the reported attacker with no
	// extra setup. Subject is the actual BEC payload from the incident
	// sample (email_atlas_1_impersonation.json).
	email, trust := e.ReceiveEmail(
		"ahmed.r@gmail.com",
		"Ahmed Rahman — XYZ Steel",
		"ops@atlas-construction.com",
		"URGENT — Updated Bank Account Details for Payment",
		"Please update our payment details immediately. New bank: First National Bank, account 8834-2291-0057.",
	)

	if email != nil {
		t.Errorf("blocked email should not be returned, got %+v", email)
	}
	if trust == nil {
		t.Fatal("expected a non-nil trust result for a blocked message")
	}
	if !trust.Blocked {
		t.Error("trust.Blocked should be true")
	}
	if trust.Status != "blocked" {
		t.Errorf("trust.Status = %q, want blocked", trust.Status)
	}
	if trust.BlockMode != "exact" {
		t.Errorf("trust.BlockMode = %q, want exact", trust.BlockMode)
	}

	// Inbox should be empty — blocked emails must not be stored.
	if got := len(e.GetInbox()); got != 0 {
		t.Errorf("blocked email was stored in inbox: got %d emails, want 0", got)
	}

	// And no JSON file should be written to inbox/.
	entries, _ := os.ReadDir(filepath.Join(dir, "inbox"))
	for _, ent := range entries {
		if strings.HasSuffix(ent.Name(), ".json") {
			t.Errorf("blocked email wrote %q to inbox/", ent.Name())
		}
	}

	// Audit trail must record the drop.
	audit, err := os.ReadFile(filepath.Join(dir, "blocked.jsonl"))
	if err != nil {
		t.Fatalf("blocked.jsonl not written: %v", err)
	}
	if !strings.Contains(string(audit), `"ahmed.r@gmail.com"`) {
		t.Errorf("audit log missing the blocked sender, got:\n%s", audit)
	}
	if !strings.Contains(string(audit), `"mode":"exact"`) {
		t.Errorf("audit log missing block mode, got:\n%s", audit)
	}
}

// TestEngine_ReceiveEmail_BECContentFilter ensures the content rules
// fire even when the sender is not on the blocklist (catches the
// "newly-registered throwaway address" case).
func TestEngine_ReceiveEmail_BECContentFilter(t *testing.T) {
	dir := t.TempDir()
	e := NewEngine(dir)

	email, trust := e.ReceiveEmail(
		"new-throwaway-9f8s7@proton.me", // unknown sender, not on blocklist
		"Mark",
		"ops@atlas-construction.com",
		"URGENT — wire instructions updated",
		"Please update our banking details effective immediately for the new account.",
	)
	if email != nil {
		t.Errorf("BEC content payload should be dropped at gateway, got %+v", email)
	}
	if trust == nil || !trust.Blocked {
		t.Fatalf("expected blocked=true, got %+v", trust)
	}
	if trust.BlockMode != "bec_filter" {
		t.Errorf("BlockMode = %q, want bec_filter", trust.BlockMode)
	}
	if trust.BlockMatch != "urgent-bank-change" {
		t.Errorf("BlockMatch = %q, want urgent-bank-change", trust.BlockMatch)
	}
}

// TestEngine_ReceiveEmail_DisplayNameSpoof covers the CEO-fraud path
// where the address is unknown but the display name claims a senior
// role on a public mail provider.
func TestEngine_ReceiveEmail_DisplayNameSpoof(t *testing.T) {
	dir := t.TempDir()
	e := NewEngine(dir)

	email, trust := e.ReceiveEmail(
		"impostor@gmail.com",
		"Jane Smith, CEO",
		"finance@atlas-construction.com",
		"Are you at your desk?",
		"I need a favor handled quietly before I get back.",
	)
	if email != nil {
		t.Errorf("display-name spoof should be dropped, got %+v", email)
	}
	if trust == nil || !trust.Blocked {
		t.Fatalf("expected blocked=true, got %+v", trust)
	}
	if trust.BlockMode != "display_name_spoof" {
		t.Errorf("BlockMode = %q, want display_name_spoof", trust.BlockMode)
	}
}

// TestEngine_ReceiveEmail_AllowsLegitimate is the negative control: a
// normal-looking vendor email on a corporate domain is NOT blocked.
func TestEngine_ReceiveEmail_AllowsLegitimate(t *testing.T) {
	dir := t.TempDir()
	e := NewEngine(dir)

	email, trust := e.ReceiveEmail(
		"ap@vendorxyz.com",
		"Accounts Payable — XYZ Steel",
		"ops@atlas-construction.com",
		"Invoice #1042 for May services",
		"Please remit via the account on file. Thanks for your business.",
	)
	if email == nil {
		t.Fatalf("legitimate email was dropped at gateway: trust=%+v", trust)
	}
	if trust.Status == "blocked" {
		t.Errorf("legitimate email flagged as blocked: %+v", trust)
	}
}

// TestEngine_BlockSender_API exercises the engine-level management
// surface, used by future UI / API endpoints to add new entries.
func TestEngine_BlockSender_API(t *testing.T) {
	dir := t.TempDir()
	e := NewEngine(dir)

	// Block a future variant of the same actor — this is the workflow
	// an analyst would use after spotting a follow-up attempt.
	if err := e.BlockSender("exact", "ahmed.rahman@gmail.com", "BEC: same actor, new address", "security"); err != nil {
		t.Fatalf("BlockSender: %v", err)
	}
	email, trust := e.ReceiveEmail(
		"ahmed.rahman@gmail.com", "Ahmed Rahman", "ops@atlas-construction.com",
		"Quick follow-up", "Just checking in on the wire transfer.",
	)
	if email != nil || (trust != nil && trust.Status == "blocked") {
		t.Errorf("newly-blocked sender was not blocked: email=%+v trust=%+v", email, trust)
	}

	// Unblock and re-test — the management surface must support both
	// directions (false positives happen).
	if err := e.UnblockSender("exact", "ahmed.rahman@gmail.com"); err != nil {
		t.Fatalf("UnblockSender: %v", err)
	}
	email2, _ := e.ReceiveEmail(
		"ahmed.rahman@gmail.com", "Ahmed Rahman", "ops@atlas-construction.com",
		"Hello", "Just a test.",
	)
	if email2 == nil {
		t.Error("after unblock, the same sender was still rejected at the gateway")
	}

	// ListBlocked should now contain the default seed but not the removed entry.
	listed := e.ListBlocked()
	for _, e := range listed {
		if e.Mode == "exact" && e.Value == "ahmed.rahman@gmail.com" {
			t.Error("removed entry still appears in ListBlocked")
		}
	}
	if len(listed) == 0 {
		t.Error("ListBlocked should at least contain the default seed")
	}
}
