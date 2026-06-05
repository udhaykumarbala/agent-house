package email

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

// BlockedEntry is a single address or domain that the gateway will refuse.
//
// Mode determines how the entry matches:
//   - "exact":   the lower-cased full email address (e.g. "ahmed.r@gmail.com")
//   - "domain":  the lower-cased domain (e.g. "gmail.com")
//   - "pattern": a Go regexp (RE2) applied to the lower-cased address
type BlockedEntry struct {
	Mode    string `json:"mode"`              // "exact" | "domain" | "pattern"
	Value   string `json:"value"`             // address, domain, or regex source
	Reason  string `json:"reason,omitempty"`  // why it was blocked (audit trail)
	AddedBy string `json:"added_by,omitempty"`
	AddedAt string `json:"added_at,omitempty"` // RFC3339
}

// BlockedEvent is a single audit record of an email the gateway dropped.
// Written to blocked.jsonl (append-only) so security can review what was
// caught without scrolling the inbox.
type BlockedEvent struct {
	Time     string `json:"time"`              // RFC3339
	From     string `json:"from"`
	FromName string `json:"from_name,omitempty"`
	Subject  string `json:"subject,omitempty"`
	Reason   string `json:"reason"`
	Mode     string `json:"mode,omitempty"` // "exact" | "domain" | "pattern" | "bec_filter" | "display_name_spoof"
	Match    string `json:"match,omitempty"` // the rule that fired
}

// becRule is one entry in the built-in anti-BEC content filter set.
// All rules are case-insensitive; the compiled regexes match against
// the lower-cased (subject + body) pair.
type becRule struct {
	name    string
	pattern *regexp.Regexp
}

// DefaultBlocklist returns a Blocklist pre-seeded with the high-confidence
// BEC entries every modern mail gateway uses. These are NOT customizable
// at runtime — they're the floor. Per-scope additions go in blocklist.json.
//
// Seeded with the sender currently being abused. Adding future variants
// of the same actor is a one-liner against this file (e.g.
// ahmed.rahman@gmail.com, ahmed.r@proton.me).
func DefaultBlocklist() []BlockedEntry {
	now := time.Now().UTC().Format(time.RFC3339)
	return []BlockedEntry{
		{
			Mode:    "exact",
			Value:   "ahmed.r@gmail.com",
			Reason:  "BEC: confirmed impersonation of XYZ Steel Corp (gmail, urgent bank-account change)",
			AddedBy: "security",
			AddedAt: now,
		},
	}
}

// DefaultBECPatterns is the gateway's built-in anti-phishing / anti-BEC
// heuristic set. A rule firing causes the email to be dropped at the
// gateway (same outcome as a blocklist hit) — but the audit log marks
// `mode = "bec_filter"` so security can distinguish address-based blocks
// from content-based blocks.
func DefaultBECPatterns() []becRule {
	mk := func(name, src string) becRule {
		return becRule{name: name, pattern: regexp.MustCompile("(?i)" + src)}
	}
	return []becRule{
		// Classic "wire / bank change" payload. Phrased loosely so it
		// catches common variants ("updated banking details", "new
		// account number", "wire instructions changed").
		mk("urgent-bank-change", `\b(urgent|immediate(ly)?|asap)\b.{0,40}\b(bank(ing)?|wire|routing|swift|account)\b.{0,40}\b(change|update|details|new|correct(ed|ion)?)\b`),

		// "Send money to a new account" — paired with invoice/PO
		// language, this is the dominant BEC payload.
		mk("invoice-bank-change", `\b(invoice|po|payment)\b.{0,80}\b(bank|wire|account|swift|routing)\b.{0,80}\b(change|update|new|correct|replace)\b`),

		// Gift-card / wire-to-personal-account scam — text lifted
		// almost verbatim from the FBI IC3 reports. The "do not
		// mention" caveat is itself a tell.
		mk("gift-card-scam", `\b(buy|purchase|send|get)\b.{0,40}\b(gift\s*card(s)?|itunes|google\s*play|steam|amazon)\b.{0,80}\b(do not (tell|mention|inform|notify)|keep (it|this) (confidential|secret|between))\b`),

		// CEO/exec impersonation — "are you at your desk", "need a
		// favor", "i need you to process something for me". The
		// FromName display-name spoof check is what makes this rule
		// high-confidence; both have to fire.
		mk("exec-impersonation-favor", `\b(are you (at|by) (your )?(desk|computer|office)|i need (a |an )?favor|do (you|me) a favor|need (you|this) handled (quietly|discreetly|before (i|the)))\b`),

		// "Click here to view document" + first-time sender is a
		// phishing marker. Kept conservative (requires "secure" /
		// "verify" language too) to avoid false-positives on
		// legitimate file shares.
		mk("phish-secure-link", `\b(secure|verify|validate|confirm)\b.{0,40}\b(document|file|invoice|account|password)\b.{0,80}\b(https?://\S+)\b`),
	}
}

// execRoleKeywords is the list of display-name tokens that, when present
// in a sender's name on a public mail provider, indicate a likely
// CEO/CFO/exec-impersonation BEC. The list is intentionally narrow —
// false positives are far more expensive than false negatives here.
var execRoleKeywords = []string{
	"ceo", "cfo", "coo", "cto", "cmo",
	"chief executive", "chief financial", "chief operating", "chief technology",
	"managing director", "president", "vp ", "vice president",
	"director of finance", "head of finance", "head of accounting",
	"controller", "treasurer",
}

// domainOf returns the lower-cased domain of an email address, or "" if
// the address is malformed. Mirrors capability.domainOf — duplicated
// here because the email package cannot import capability (capability
// already imports email).
func domainOf(addr string) string {
	at := strings.LastIndex(addr, "@")
	if at < 0 || at == len(addr)-1 {
		return ""
	}
	return strings.ToLower(addr[at+1:])
}

// isPublicDomain reports whether the domain is a free / consumer mail
// provider. Mirrors capability.isPublicDomain for the same reason as
// domainOf above.
func isPublicDomain(d string) bool {
	switch strings.ToLower(d) {
	case "gmail.com", "yahoo.com", "outlook.com", "hotmail.com",
		"proton.me", "protonmail.com", "icloud.com", "aol.com",
		"live.com", "msn.com", "yandex.com", "yandex.ru", "mail.ru":
		return true
	}
	return false
}

// Blocklist is the per-engine sender deny list plus the BEC pattern
// rules. Entries are loaded from dataDir/blocklist.json and audited
// drops are appended to dataDir/blocked.jsonl.
type Blocklist struct {
	mu       sync.RWMutex
	path     string // blocklist.json
	auditLog string // blocked.jsonl
	entries  []BlockedEntry
	// Compiled patterns are kept separate from the source so a
	// malformed regex never panics the gateway — the load step drops
	// bad patterns and the load still succeeds.
	patterns   []*regexp.Regexp
	patternSrc []string
	// becPatterns are gateway-wide heuristic rules that fire on any
	// (subject, body) pair regardless of the sender.
	becPatterns []becRule
}

// NewBlocklist returns a Blocklist that loads persisted entries from
// dataDir/blocklist.json and falls back to the default set when the
// file does not exist or fails to parse.
func NewBlocklist(dataDir string) *Blocklist {
	b := &Blocklist{
		path:     filepath.Join(dataDir, "blocklist.json"),
		auditLog: filepath.Join(dataDir, "blocked.jsonl"),
	}
	_ = os.MkdirAll(dataDir, 0o755)
	b.load()
	if len(b.entries) == 0 {
		// First boot — seed with defaults so the gateway is
		// protective even before an operator touches the file.
		b.entries = append(b.entries, DefaultBlocklist()...)
		_ = b.save()
	}
	b.becPatterns = DefaultBECPatterns()
	return b
}

func (b *Blocklist) load() {
	data, err := os.ReadFile(b.path)
	if err != nil {
		return
	}
	var entries []BlockedEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return
	}
	b.entries = entries
	b.compilePatterns()
}

func (b *Blocklist) save() error {
	data, err := json.MarshalIndent(b.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(b.path, data, 0o644)
}

func (b *Blocklist) compilePatterns() {
	b.patterns = b.patterns[:0]
	b.patternSrc = b.patternSrc[:0]
	for _, e := range b.entries {
		if e.Mode != "pattern" {
			continue
		}
		re, err := regexp.Compile(e.Value)
		if err != nil {
			// Drop a bad pattern rather than failing the whole
			// load — partial protection beats a panic.
			continue
		}
		b.patterns = append(b.patterns, re)
		b.patternSrc = append(b.patternSrc, e.Value)
	}
}

// Add inserts a new block entry. Persists the file on success.
func (b *Blocklist) Add(e BlockedEntry) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if e.Mode == "" {
		e.Mode = "exact"
	}
	e.Value = strings.ToLower(strings.TrimSpace(e.Value))
	if e.Value == "" {
		return fmt.Errorf("blocklist: empty value")
	}
	if e.Mode != "exact" && e.Mode != "domain" && e.Mode != "pattern" {
		return fmt.Errorf("blocklist: unknown mode %q", e.Mode)
	}
	if e.Mode == "pattern" {
		if _, err := regexp.Compile(e.Value); err != nil {
			return fmt.Errorf("blocklist: bad pattern: %w", err)
		}
	}
	if e.AddedAt == "" {
		e.AddedAt = time.Now().UTC().Format(time.RFC3339)
	}
	// Dedup: replace if same (mode, value) already present.
	for i := range b.entries {
		if b.entries[i].Mode == e.Mode && b.entries[i].Value == e.Value {
			b.entries[i] = e
			b.compilePatterns()
			return b.save()
		}
	}
	b.entries = append(b.entries, e)
	b.compilePatterns()
	return b.save()
}

// Remove deletes an entry by mode+value (case-insensitive on value).
func (b *Blocklist) Remove(mode, value string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	mode = strings.ToLower(strings.TrimSpace(mode))
	value = strings.ToLower(strings.TrimSpace(value))
	for i := range b.entries {
		if b.entries[i].Mode == mode && b.entries[i].Value == value {
			b.entries = append(b.entries[:i], b.entries[i+1:]...)
			b.compilePatterns()
			return b.save()
		}
	}
	return fmt.Errorf("blocklist: %s/%s not found", mode, value)
}

// List returns a copy of all current entries.
func (b *Blocklist) List() []BlockedEntry {
	b.mu.RLock()
	defer b.mu.RUnlock()
	out := make([]BlockedEntry, len(b.entries))
	copy(out, b.entries)
	return out
}

// MatchAddress checks the given address against all entries. Returns
// the first matching entry and a boolean, or (nil, false) when nothing
// matches.
func (b *Blocklist) MatchAddress(from string) (*BlockedEntry, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	from = strings.ToLower(strings.TrimSpace(from))
	// Exact + domain checks use the entries slice; pattern checks use
	// the parallel slices. Two passes keep the matching logic obvious.
	for i := range b.entries {
		e := &b.entries[i]
		switch e.Mode {
		case "exact":
			if e.Value == from {
				return e, true
			}
		case "domain":
			if d := domainOf(from); d != "" && d == e.Value {
				return e, true
			}
		}
	}
	for i, re := range b.patterns {
		if re.MatchString(from) {
			return &BlockedEntry{
				Mode:   "pattern",
				Value:  b.patternSrc[i],
				Reason: "address matches blocklist pattern",
			}, true
		}
	}
	return nil, false
}

// MatchContent runs the built-in anti-BEC content rules against the
// lower-cased (subject + body). Returns the rule name + the matched
// snippet, or empty strings when nothing fires.
func (b *Blocklist) MatchContent(subject, body string) (name, snippet string) {
	text := strings.ToLower(subject + "\n" + body)
	for _, r := range b.becPatterns {
		if loc := r.pattern.FindStringIndex(text); loc != nil {
			return r.name, text[loc[0]:loc[1]]
		}
	}
	return "", ""
}

// MatchDisplayName detects CEO-fraud style display-name spoofing: the
// sender is on a public mail provider (gmail/outlook/etc.) AND the
// display name claims a senior finance/exec role. The first signal
// alone is fine; the second alone is fine; both together is the
// textbook BEC payload.
//
// Returns the matched role keyword when both fire, or "" otherwise.
func (b *Blocklist) MatchDisplayName(from, fromName string) string {
	fromDomain := domainOf(from)
	if !isPublicDomain(fromDomain) {
		return ""
	}
	low := strings.ToLower(fromName)
	for _, role := range execRoleKeywords {
		if strings.Contains(low, role) {
			return role
		}
	}
	return ""
}

// Audit appends a drop event to blocked.jsonl. Best-effort: a write
// failure is logged but does not change the gateway's behavior.
func (b *Blocklist) Audit(ev BlockedEvent) {
	if b.auditLog == "" {
		return
	}
	if ev.Time == "" {
		ev.Time = time.Now().UTC().Format(time.RFC3339)
	}
	line, err := json.Marshal(ev)
	if err != nil {
		return
	}
	f, err := os.OpenFile(b.auditLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer f.Close()
	_, _ = f.Write(append(line, '\n'))
}
