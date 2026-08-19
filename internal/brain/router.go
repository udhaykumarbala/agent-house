package brain

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"pty-claude-test/internal/session"
)

// Router handles Brain decision-making via oneshot API calls.
type Router struct {
	apiClient    *session.APIClient
	projectsDir  string
	conversations *ConversationStore
	systemPrompt string
}

// NewRouter creates a Brain router.
func NewRouter(apiClient *session.APIClient, projectsDir string, conversations *ConversationStore) *Router {
	return &Router{
		apiClient:     apiClient,
		projectsDir:   projectsDir,
		conversations: conversations,
		systemPrompt:  brainSystemPrompt,
	}
}

// Route processes a user message and returns a BrainDecision.
// This is the core Brain loop: assemble context → call API → parse decision.
// Supports re-calls for server-mediated actions (read_file, etc.)
func (r *Router) Route(ctx context.Context, userID, message, scope string) (*BrainDecision, error) {
	// Store user message
	r.conversations.Add(userID, "user", message)

	// Assemble workspace context
	workspaceCtx := AssembleContext(r.projectsDir)

	// Make the Brain company-aware. Agent House runs two companies; the scope
	// selects which one so the Brain routes intents correctly (and never offers
	// the other company's actions).
	switch scope {
	case "software":
		workspaceCtx = "COMPANY: Software Studio — an agent-driven software solutions company. You are its Conductor.\n" +
			"- \"Build / Create / Make <an app>\" → action create_project (pick a short, kebab-case project_id from the app's name).\n" +
			"- \"Change / Add / Fix / Improve <X> in <app>\" → action delegate to the right engineer on that project.\n" +
			"- The team runs a full autonomous SDLC: PRD → design → development → QA, with human checkpoints whose autonomy is governed by the project's run mode.\n" +
			"- This studio has NO email inbox, vendors, applicants, invoices, or construction-site data. If asked to process an inbox, validate an invoice, check vendors, or run a site scenario, briefly say that belongs to the construction company — NOT the software studio. NEVER invent inbox/email/vendor/applicant results here.\n" +
			"- Do NOT use the construction scenarios (process_inbox, validate_invoice, route_rfi, schedule_check, morning_briefing, workforce_snapshot) — those belong to the EPC company.\n\n" +
			workspaceCtx
	case "", "default":
		// Generic — no company banner.
	default:
		// EPC construction site/tenant.
		workspaceCtx = fmt.Sprintf("CURRENT EPC SITE / SCOPE: %s\n(All scenarios and capability data operate on this scope.)\n\n%s", scope, workspaceCtx)
	}

	// Get recent conversation history
	history := r.conversations.GetRecent(userID, 10)

	// Build the full prompt
	prompt := r.buildPrompt(workspaceCtx, history, message)

	// Call Brain (oneshot)
	decision, err := r.callBrain(ctx, prompt)
	if err != nil {
		return nil, err
	}

	// Handle server-mediated actions (re-call loop, max 3)
	for i := 0; i < 3; i++ {
		if decision.Action != ActionReadFile {
			break
		}

		// Read the requested file
		filePath := decision.Params["path"]
		project := decision.Params["project"]
		if project != "" && !strings.Contains(filePath, "/") {
			filePath = project + "/" + filePath
		}

		fullPath := r.projectsDir + "/" + filePath
		content, err := readFileContent(fullPath)
		if err != nil {
			content = fmt.Sprintf("Error reading file: %v", err)
		}

		// Re-call Brain with file content
		followUp := fmt.Sprintf("Here's the content of %s:\n\n```\n%s\n```\n\nNow respond to the user's original request.", filePath, truncateFile(content, 8000))

		r.conversations.Add(userID, "assistant", fmt.Sprintf("[Reading %s...]", filePath))
		prompt = r.buildPrompt(workspaceCtx, r.conversations.GetRecent(userID, 12), followUp)
		decision, err = r.callBrain(ctx, prompt)
		if err != nil {
			return nil, err
		}
	}

	// If still asking for files after 3 reads, escalate
	if decision.Action == ActionReadFile {
		decision.Action = ActionEscalate
		decision.Response = "This requires deeper analysis. Let me look into it more thoroughly."
	}

	// Safety: a reply only goes out AFTER the user has seen a draft and confirmed.
	// Confirmation is proven by CONVERSATION STATE (the prior assistant turn was a
	// draft) — NOT by keywords, because a chip label like "Send reply to Sarah"
	// contains "send" yet no draft was ever shown.
	if decision.Action == ActionSendReply {
		priorDraft := false
		for i := len(history) - 1; i >= 0; i-- {
			if history[i].Role == "assistant" {
				lc := strings.ToLower(history[i].Content)
				priorDraft = strings.Contains(lc, "confirm to send") ||
					strings.Contains(lc, "draft reply") ||
					strings.Contains(lc, "draft:")
				break
			}
		}
		if !priorDraft {
			log.Printf("[BRAIN] send_reply with no prior draft — showing the draft first")
			if body := decision.Params["body"]; body != "" {
				draft := fmt.Sprintf("**Draft Reply:**\n\n---\n\nTo: %s\nSubject: %s\n\n%s\n\n---\n\nConfirm to send?",
					decision.Params["to"], decision.Params["subject"], body)
				decision = &BrainDecision{
					Action:      ActionRespond,
					Response:    draft,
					Suggestions: []string{"Send this reply", "Edit draft", "Cancel"},
				}
			} else {
				decision = &BrainDecision{
					Action:      ActionRespond,
					Response:    "Let me draft that reply first so you can review it before it goes out.",
					Suggestions: []string{"Draft the reply", "Cancel"},
				}
			}
		}
	}

	// Store assistant response
	if decision.Response != "" {
		r.conversations.Add(userID, "assistant", decision.Response)
	}

	return decision, nil
}

func (r *Router) callBrain(ctx context.Context, prompt string) (*BrainDecision, error) {
	start := time.Now()

	resp, err := r.apiClient.SendMessage(ctx, r.systemPrompt, prompt)
	if err != nil {
		return nil, fmt.Errorf("brain API call failed: %w", err)
	}

	log.Printf("[BRAIN] API call: %dms, %d in / %d out tokens",
		time.Since(start).Milliseconds(), resp.Usage.InputTokens, resp.Usage.OutputTokens)

	// Parse structured JSON from response
	decision, err := parseDecision(resp.Text)
	if err != nil {
		// JSON parsing failed — infer action from text content
		log.Printf("[BRAIN] JSON parse failed, inferring action from text: %v", err)
		decision = inferActionFromText(resp.Text)
	}

	log.Printf("[BRAIN] Decision: action=%s", decision.Action)
	return decision, nil
}

func (r *Router) buildPrompt(workspaceCtx string, history []ChatMessage, currentMessage string) string {
	var sb strings.Builder

	sb.WriteString(workspaceCtx)
	sb.WriteString("\n")

	if len(history) > 1 {
		sb.WriteString("## Recent Conversation\n\n")
		// Show all but the last message (which is the current one)
		for _, msg := range history[:len(history)-1] {
			if msg.Role == "user" {
				sb.WriteString("**User:** " + msg.Content + "\n\n")
			} else {
				sb.WriteString("**Assistant:** " + msg.Content + "\n\n")
			}
		}
	}

	sb.WriteString("## Current Message\n\n")
	sb.WriteString(currentMessage)

	return sb.String()
}

func parseDecision(text string) (*BrainDecision, error) {
	text = strings.TrimSpace(text)

	// Strip markdown code block wrappers
	if strings.HasPrefix(text, "```") {
		lines := strings.Split(text, "\n")
		var clean []string
		for _, l := range lines {
			if strings.HasPrefix(strings.TrimSpace(l), "```") {
				continue
			}
			clean = append(clean, l)
		}
		text = strings.Join(clean, "\n")
	}

	var decision BrainDecision

	// Try 1: Direct unmarshal of full text
	if err := json.Unmarshal([]byte(text), &decision); err == nil && decision.Action != "" {
		goto parsed
	}

	// Try 2: Find {"action" and use json.Decoder (handles strings with special chars)
	if idx := strings.Index(text, `{"action"`); idx >= 0 {
		dec := json.NewDecoder(strings.NewReader(text[idx:]))
		if err := dec.Decode(&decision); err == nil && decision.Action != "" {
			goto parsed
		}
	}

	// Try 3: Find any { and use json.Decoder
	if idx := strings.Index(text, "{"); idx >= 0 {
		dec := json.NewDecoder(strings.NewReader(text[idx:]))
		if err := dec.Decode(&decision); err == nil && decision.Action != "" {
			goto parsed
		}
	}

	// Try 4: tolerant salvage — the text IS a decision envelope but strict JSON
	// parsing failed (e.g. the long markdown "response" body carries unescaped
	// newlines). Hand-extract the response so we answer instead of dumping raw
	// JSON at the user.
	if strings.Contains(text, `"action"`) && strings.Contains(text, `"response"`) {
		if resp := extractJSONStringField(text, "response"); resp != "" {
			decision.Action = ActionRespond
			decision.Response = resp
			goto parsed
		}
	}

	return nil, fmt.Errorf("no valid JSON action found in response (%.200s)", text)

parsed:

	// Safety-net: some models wrap the real decision inside a respond's response
	// field (a nested JSON object). Unwrap it so the intended action actually fires.
	if decision.Action == ActionRespond {
		inner := strings.TrimSpace(decision.Response)
		if strings.HasPrefix(inner, "{") && strings.Contains(inner, `"action"`) {
			var nested BrainDecision
			if json.Unmarshal([]byte(inner), &nested) == nil && nested.Action != "" && nested.Action != ActionRespond {
				decision = nested
			}
		}
	}

	// ALWAYS clean suggestion text from response (even if JSON suggestions exist)
	if decision.Response != "" {
		extracted := extractSuggestionsFromText(decision.Response)
		decision.Response = cleanSuggestionsFromText(decision.Response)
		// Use extracted suggestions if JSON field was empty
		if len(decision.Suggestions) == 0 && len(extracted) > 0 {
			decision.Suggestions = extracted
		}
	}

	// If STILL no suggestions, generate defaults based on context
	if len(decision.Suggestions) == 0 {
		decision.Suggestions = generateDefaultSuggestions(&decision)
	}

	return &decision, nil
}

// extractJSONStringField pulls a single string field value out of a JSON-ish
// blob even when the blob is NOT valid JSON (e.g. the value carries unescaped
// newlines). It locates "key": "…" and returns the unescaped value, treating
// the first quote that is followed (modulo whitespace) by , or } as the close.
func extractJSONStringField(text, key string) string {
	marker := `"` + key + `"`
	ki := strings.Index(text, marker)
	if ki < 0 {
		return ""
	}
	i := ki + len(marker)
	for i < len(text) && text[i] != '"' { // advance to the value's opening quote
		i++
	}
	if i >= len(text) {
		return ""
	}
	i++ // past the opening quote
	var b strings.Builder
	for ; i < len(text); i++ {
		c := text[i]
		if c == '\\' && i+1 < len(text) {
			switch n := text[i+1]; n {
			case 'n':
				b.WriteByte('\n')
			case 't':
				b.WriteByte('\t')
			case 'r':
				// drop carriage returns
			case '"', '\\', '/':
				b.WriteByte(n)
			default:
				b.WriteByte(n)
			}
			i++
			continue
		}
		if c == '"' {
			j := i + 1
			for j < len(text) && (text[j] == ' ' || text[j] == '\t' || text[j] == '\n' || text[j] == '\r') {
				j++
			}
			if j >= len(text) || text[j] == ',' || text[j] == '}' {
				break // real terminator
			}
			b.WriteByte('"') // literal quote inside the markdown body
			continue
		}
		b.WriteByte(c)
	}
	return strings.TrimSpace(b.String())
}

// extractSuggestionsFromText pulls suggestion-like lines from response text.
func extractSuggestionsFromText(text string) []string {
	// First try to extract JSON array: "suggestions": ["a", "b"]
	if idx := strings.Index(text, `"suggestions"`); idx >= 0 {
		rest := text[idx:]
		if start := strings.Index(rest, "["); start >= 0 {
			if end := strings.Index(rest[start:], "]"); end >= 0 {
				var arr []string
				if json.Unmarshal([]byte(rest[start:start+end+1]), &arr) == nil && len(arr) > 0 {
					return arr
				}
			}
		}
	}

	// Fallback: extract from bullet lines
	var suggestions []string
	lines := strings.Split(text, "\n")
	inSuggestions := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		lower := strings.ToLower(trimmed)

		// Detect suggestions section
		if strings.Contains(lower, "suggestion") || strings.Contains(lower, "you can") || strings.Contains(lower, "next steps") || strings.Contains(lower, "options:") {
			inSuggestions = true
			continue
		}

		if inSuggestions {
			// Extract quoted text from bullet lines
			if strings.HasPrefix(trimmed, "•") || strings.HasPrefix(trimmed, "-") || strings.HasPrefix(trimmed, "*") {
				s := strings.TrimLeft(trimmed, "•-* ")
				// Remove surrounding quotes
				s = strings.Trim(s, "\"'`")
				s = strings.TrimSpace(s)
				if len(s) > 3 && len(s) < 80 {
					suggestions = append(suggestions, s)
				}
			}
		}
	}

	// Cap at 4
	if len(suggestions) > 4 {
		suggestions = suggestions[:4]
	}
	return suggestions
}

// cleanSuggestionsFromText removes the suggestions section from response text.
func cleanSuggestionsFromText(text string) string {
	lines := strings.Split(text, "\n")
	var cleaned []string
	skip := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		lower := strings.ToLower(trimmed)

		// Start skipping at suggestion headers
		if strings.Contains(lower, "suggestion") || strings.Contains(lower, "next steps") ||
			strings.Contains(lower, "would you like") || strings.Contains(lower, "what would you like") ||
			strings.Contains(lower, "options:") || strings.Contains(lower, "actions available") {
			skip = true
			continue
		}

		// Skip bullet lines and JSON-style arrays
		if skip {
			if strings.HasPrefix(trimmed, "•") || strings.HasPrefix(trimmed, "-") ||
				strings.HasPrefix(trimmed, "*") || strings.HasPrefix(trimmed, "[") ||
				strings.HasPrefix(trimmed, "\"") || trimmed == "" {
				continue
			}
			// Non-bullet, non-empty line after suggestions → stop skipping
			skip = false
		}

		cleaned = append(cleaned, line)
	}

	result := strings.TrimSpace(strings.Join(cleaned, "\n"))
	result = strings.TrimRight(result, "-\n ")

	// Also strip inline JSON suggestion arrays like: "suggestions": ["a", "b"]
	if idx := strings.Index(result, `"suggestions"`); idx >= 0 {
		// Find the end of the array
		rest := result[idx:]
		if end := strings.Index(rest, "]"); end >= 0 {
			result = strings.TrimSpace(result[:idx] + result[idx+end+1:])
		}
	}

	// Strip trailing commas, quotes, braces from JSON remnants
	result = strings.TrimRight(result, " ,}\"")
	return strings.TrimSpace(result)
}

// inferActionFromText analyzes text response to determine the intended action.
// Called when the LLM doesn't return valid JSON.
func inferActionFromText(text string) *BrainDecision {
	lower := strings.ToLower(text)

	// Check for send/reply intent
	if (strings.Contains(lower, "reply sent") || strings.Contains(lower, "email sent") ||
		strings.Contains(lower, "message sent") || strings.Contains(lower, "sending reply")) &&
		!strings.Contains(lower, "draft") {
		// Extract "to" address from text
		to := extractEmailAddr(text)
		return &BrainDecision{
			Action:   ActionSendReply,
			Params:   map[string]string{"to": to, "subject": "Re:", "body": text},
			Response: text,
		}
	}

	// Check for delete intent
	if strings.Contains(lower, "deleted") || strings.Contains(lower, "removed") ||
		(strings.Contains(lower, "delete") && (strings.Contains(lower, "email") || strings.Contains(lower, "impersonation") || strings.Contains(lower, "scam"))) ||
		strings.Contains(lower, "delete_email") {
		emailID := extractIDFromText(text, "email_")
		return &BrainDecision{
			Action:   ActionDeleteEmail,
			Params:   map[string]string{"email_id": emailID},
			Response: text,
		}
	}

	// Check for shortlist intent
	if strings.Contains(lower, "shortlist") || strings.Contains(lower, "shortlisted") {
		appID := extractIDFromText(text, "app_")
		return &BrainDecision{
			Action:   ActionShortlistApplicant,
			Params:   map[string]string{"applicant_id": appID},
			Response: text,
		}
	}

	// Check for project creation
	if strings.Contains(lower, "creating project") || strings.Contains(lower, "project created") {
		return &BrainDecision{
			Action:   ActionCreateProject,
			Response: text,
		}
	}

	// Default: treat as text response
	return &BrainDecision{
		Action:   ActionRespond,
		Response: text,
	}
}

// extractEmailAddr finds an email address in text.
func extractEmailAddr(text string) string {
	words := strings.Fields(text)
	for _, w := range words {
		w = strings.Trim(w, "<>(),;\"'")
		if strings.Contains(w, "@") && strings.Contains(w, ".") {
			return w
		}
	}
	return ""
}

// extractIDFromText finds an ID with a given prefix in text.
func extractIDFromText(text, prefix string) string {
	idx := strings.Index(text, prefix)
	if idx < 0 {
		return ""
	}
	// Extract until non-alphanumeric/underscore
	end := idx + len(prefix)
	for end < len(text) && (text[end] >= '0' && text[end] <= '9' || text[end] == '_') {
		end++
	}
	if end > idx+len(prefix) {
		return text[idx:end]
	}
	return ""
}

// generateDefaultSuggestions creates contextual suggestions when the LLM didn't provide any.
func generateDefaultSuggestions(d *BrainDecision) []string {
	resp := strings.ToLower(d.Response)

	// Email-related responses
	if strings.Contains(resp, "email") || strings.Contains(resp, "inbox") || strings.Contains(resp, "mail") {
		return []string{"Draft reply to vendor", "Delete suspicious email", "Review next email", "Check contract terms"}
	}
	// Project-related
	if d.Action == "create_project" {
		return []string{"Check progress", "View live dashboard", "Add requirements"}
	}
	if d.Action == "project_status" || d.Action == "list_projects" {
		return []string{"Create new project", "Check emails", "Generate report"}
	}
	// Vendor-related
	if strings.Contains(resp, "vendor") || strings.Contains(resp, "contract") || strings.Contains(resp, "delivery") {
		return []string{"Draft firm reply", "Check penalty terms", "Escalate to management", "View all vendors"}
	}
	// Default
	return []string{"Show inbox", "List projects", "Create a project"}
}

func readFileContent(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func truncateFile(content string, maxLen int) string {
	if len(content) <= maxLen {
		return content
	}
	return content[:maxLen] + "\n... (truncated)"
}

// brainSystemPrompt is the core system prompt for the Brain agent.
var brainSystemPrompt = `You are the Brain of Agent House — an AI workspace manager.

You manage projects built by teams of AI agents (CEO, PM, UX, UI, Security, Architect, Senior Dev, Junior Dev).

## Critical Rules

1. ALWAYS respond with a single JSON object. No other text outside the JSON.
2. NEVER say "above", "as shown", "draft ready above", "review the draft above". Each response is SELF-CONTAINED — the user can ONLY see your response field.
3. When drafting an email, ALWAYS include the FULL draft text in your response field. Example:

{"action": "respond", "response": "Here is the draft reply:\n\n---\n\nSubject: Re: Steel Delivery\n\nDear Ahmed,\n\n[full email body here]\n\nRegards,\n[Name]\n\n---\n\nConfirm to send?", "suggestions": ["Send this reply", "Edit draft", "Cancel"]}

4. ONLY use send_reply action AFTER the user confirms (says "yes", "send", "confirm"). First show the draft with action "respond", then send on confirmation.
5. For send_reply, include email_id of the original email in params.
6. Put suggestions in the JSON "suggestions" array, NEVER as bullet text in the response.
7. NEVER invent data. Do not state counts, names, IDs, or facts that are not in the workspace context. If asked for something that isn't there (e.g. "list all 5 candidates" when only one applicant exists), state exactly what IS there and that there are no others. "Only one applicant (Raj Kumar) has applied" beats fabricating names or numbers — a wrong count destroys trust.
8. The only project/scope is the current EPC site (e.g. "atlas-site"). "Project Alpha / Beta / Gamma" are work items WITHIN that site, NOT separate Agent House projects — never call project_status or read_file with project_id "alpha"/"beta"/"gamma", and never suggest doing so. For their status use run_scenario "schedule_check" or "morning_briefing".
9. escalate is ONLY for an initial, open-ended "go deeply analyze X" request. If the user ASKS YOU TO PRODUCE a plan/summary/analysis (e.g. "give me the mitigation plan", "show the cascading impact and recovery options", or after approving an escalation), RESPOND directly with the full plan — do NOT escalate again.
10. Any reply-related suggestion you emit MUST say "Draft …" (never "Send …"). A reply is only ever sent after the user has seen the draft and explicitly confirms.

## Available Actions

{"action": "respond", "response": "Your answer here"}
  → Answer the user directly. Use for questions, status requests, greetings.

{"action": "create_project", "params": {"project_id": "my-project", "task": "Build a..."}, "response": "Creating project..."}
  → Start a new project. Choose a short, lowercase, hyphenated project_id.

{"action": "delegate", "params": {"project_id": "existing-project", "agent_role": "senior_dev", "task": "Add dark mode"}, "response": "Delegating to Senior Dev..."}
  → Send a task to a specific agent in an existing project.

{"action": "project_status", "params": {"project_id": "my-project"}, "response": "Let me check..."}
  → Get detailed status of a project (server will provide and re-call you).

{"action": "read_file", "params": {"project": "my-project", "path": "src/index.html"}, "response": "Reading file..."}
  → Read a file from a project. Server reads it and re-calls you with content.

{"action": "list_projects", "response": "Here are your projects:"}
  → List all projects with status.

{"action": "escalate", "params": {"project_id": "my-project", "task": "Analyze auth security"}, "response": "This needs deep analysis..."}
  → For complex multi-file analysis. Spawns a full Claude Code session.

{"action": "delete_email", "params": {"email_id": "email_123"}, "response": "Email deleted."}
  → Delete/archive a specific email from inbox.

{"action": "send_reply", "params": {"to": "email@example.com", "subject": "Re: ...", "body": "Dear...", "email_id": "original_email_id"}, "response": "Reply sent."}
  → Send an email reply. Include email_id of the original email. User must confirm before executing.

{"action": "shortlist_applicant", "params": {"applicant_id": "app_123", "notes": "Strong candidate"}, "response": "Applicant shortlisted."}
  → Mark a job applicant as shortlisted for a position.

{"action": "archive_email", "params": {"email_id": "email_123"}, "response": "Email archived."}
  → Archive an email (move out of active inbox).

{"action": "run_scenario", "params": {"scenario": "process_inbox"}, "response": "Sweeping the inbox and routing each item to the right specialist..."}
  → Run a multi-agent EPC operations scenario that FANS WORK OUT to specialist agents (procurement, hr, project_manager, site_engineer) — they pulse live in the workforce panel. Use this for EPC operations. Available scenarios (put the name in params.scenario):
    • "process_inbox"      → triage/sweep the whole inbox, route each email to a discipline
    • "morning_briefing"   → one-page brief: schedule slips + inbox state + HR pipeline + vendor health
    • "validate_invoice"   → verify a vendor invoice; flag impersonation/BEC (params: sender_email, vendor_id, amount)
    • "route_rfi"          → classify an incoming RFI by discipline + find a matching specialist (params: rfi_id, subject, body)
    • "schedule_check"     → scan milestones; flag slipping ones with mitigations
    • "process_applicants" → screen/rank job applicants for a role (params: request)
    • "workforce_snapshot" → LIVE workforce picture from the real Worqplace HRMS (view-only): headcount, active projects, expiring documents, saudization ratio, per-project staffing

## Decision Rules

1. Greeting or simple question → respond
2. "What projects exist?" / "What's running?" → list_projects
3. "How's project X?" / "Status of X" → project_status
4. "Build/Create/Make [something]" → create_project (pick a good project_id)
5. "Change/Add/Fix [something] in [project]" → delegate to appropriate agent
6. "Show me [file] from [project]" → read_file
7. Complex analysis across multiple files → escalate
8. If unsure, ask for clarification via respond
9. "Process / triage / sweep my inbox" → run_scenario (process_inbox)
10. "What needs my attention today / morning briefing" → run_scenario (morning_briefing)
11. "Is this invoice/vendor legit / verify invoice / is this email safe" → run_scenario (validate_invoice)
12. "Incoming RFI / who handles RFI #N" → run_scenario (route_rfi)
13. "Check the schedule / what's slipping / milestone status" → run_scenario (schedule_check)
14. "Screen / rank applicants / who can fill [role]" → run_scenario (process_applicants)
15. "Headcount / workforce status / how many employees / expiring documents (iqama, passport) / saudization / HRMS data" → run_scenario (workforce_snapshot) — this pulls LIVE data from the company's real HRMS

## Proactive Delegation (CRITICAL — do not just narrate)
When a task clearly belongs to a discipline or an operations scenario, DELEGATE or RUN_SCENARIO in the SAME turn — never merely state which agent *should* handle it. If your response names an agent or says you will route / triage / verify / check / screen something, your "action" MUST be "delegate" or "run_scenario", NOT "respond". Keep the descriptive prose in "response" for the human, but always carry the real action so work actually fans out.

## Confirmation for Destructive Actions

For these actions, FIRST show the user what will happen and ask for confirmation:
- delete_email: Show which email will be deleted, ask "Confirm delete?"
- send_reply: Show the full draft, ask "Send this reply?" with suggestions ["Send it", "Edit draft", "Cancel"]
- Only execute the action when the user explicitly confirms (says "yes", "send it", "confirm", "delete it")

When the user confirms, use the ACTUAL email_id or applicant_id from the context (e.g., email_1774159516188, not email_2).

## Vendor & Email Awareness

You have access to vendor data and email inbox in the workspace context.
When discussing vendor emails, ALWAYS reference:
- The vendor's contract terms (from vendors.json data in context)
- The vendor's trusted contacts
- Any trust alerts on the email

When drafting replies to vendor emails:
- Reference specific contract clauses (delays, penalties, payment terms)
- Use the vendor contact's name
- Be factual and reference the data you have
- If the email has impersonation risk, WARN the user prominently

## Suggestions (IMPORTANT)

ALWAYS include a "suggestions" array with 2-4 quick follow-up actions the user might want.
These become clickable buttons in the UI.

Example:
{"action": "respond", "response": "...", "suggestions": ["Draft firm reply", "Check contract terms", "Delete this email", "Escalate to management"]}

Make suggestions contextual:
- After reviewing an email → "Draft reply", "Delete email", "Flag as urgent"
- After showing project status → "View files", "Assign new task", "Generate report"
- After creating a project → "Check progress", "View live dashboard", "Add requirements"
- After listing emails → "Review email 1", "Check impersonation alert", "Draft replies for all"

## Project ID Rules
- Lowercase, hyphens only: "landing-page", "todo-app", "api-server"
- Short and descriptive
- No spaces or special characters

## Agent Roles for Delegation
- senior_dev: Code changes, new features, bug fixes
- junior_dev: Simple code tasks, file creation
- architect: Architecture decisions, structure changes
- pm: Requirements, documentation, specs
- security: Security audits, vulnerability fixes
- ux: UX improvements, flow changes
- ui: Visual design, styling changes
- ceo: Strategic decisions, reviews

## EPC / Construction Agent Roles (for delegation on EPC sites)
- procurement: vendors, invoices, POs, BEC / impersonation checks
- project_manager: schedule, milestones, slips, client progress updates
- site_engineer: RFIs, field/technical queries, drawings
- hse: safety, compliance, incident response
- qa_inspector: inspections, punch lists, quality
- hr: hiring, applicant screening`
