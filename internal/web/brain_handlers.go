package web

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"pty-claude-test/internal/agent"
	"pty-claude-test/internal/brain"
	"pty-claude-test/internal/checkpoint"
	"pty-claude-test/internal/message"
	"pty-claude-test/internal/orchestrator"
	"pty-claude-test/internal/session"
	"pty-claude-test/internal/task"
)

// globalSettingsDir holds the global default WorkflowSettings (run mode +
// decision timeout) that newly created projects inherit. It's a pseudo-project
// directory so it reuses the same checkpoint settings store + API.
const globalSettingsProject = "_global"

// BrainHandler manages the Brain chat interface.
type BrainHandler struct {
	router       *brain.Router
	executor     *brain.Executor
	convos       *brain.ConversationStore
	orch         *orchestrator.Orchestrator
	store        *message.Store
	projDir      string
	hub          *Hub
	routerClient *session.APIClient // fast tier-1 model (flash-lite) for instant intermediary feedback; nil disables it
	emailEngine  interface{ MarkReplied(string); DeleteEmail(string) bool; UpdateApplicantStatus(string,string) bool }
}

// NewBrainHandler creates the Brain handler with all dependencies wired.
func NewBrainHandler(apiClient *session.APIClient, orch *orchestrator.Orchestrator, store *message.Store, hub *Hub, projDir string) *BrainHandler {
	convos := brain.NewConversationStore(projDir + "/.brain/conversations")
	// Optional semantic search: if OPENAI_API_KEY is set, conversation
	// search upgrades from substring to cosine-ranked. Without the key the
	// store stays in lexical-only mode.
	if key := os.Getenv("OPENAI_API_KEY"); key != "" {
		model := os.Getenv("OPENAI_EMBEDDING_MODEL")
		convos.SetEmbedder(brain.NewOpenAIEmbedder(key, model))
	} else {
		log.Printf("[brain] OPENAI_API_KEY not set — conversation search will use lexical fallback")
	}
	router := brain.NewRouter(apiClient, projDir, convos)
	executor := brain.NewExecutor(projDir)

	bh := &BrainHandler{
		router:   router,
		executor: executor,
		convos:   convos,
		orch:     orch,
		store:    store,
		projDir:  projDir,
		hub:      hub,
	}

	// Wire executor callbacks
	executor.OnCreateProject = bh.handleCreateProject
	executor.OnDelegate = bh.handleDelegate
	executor.OnSendReply = bh.handleSendReply
	executor.OnDeleteEmail = bh.handleDeleteEmail
	executor.OnMarkRead = bh.handleMarkReplied
	executor.OnShortlist = bh.handleShortlist

	return bh
}

// HandleChat handles POST /api/chat
func (bh *BrainHandler) HandleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Return conversation history
		userID := r.URL.Query().Get("user")
		if userID == "" {
			userID = "default"
		}
		messages := bh.convos.GetAll(userID)
		writeJSON(w, map[string]interface{}{
			"messages": messages,
		})
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Message string `json:"message"`
		UserID  string `json:"user_id"`
		Scope   string `json:"scope"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	if req.Message == "" {
		http.Error(w, "message required", http.StatusBadRequest)
		return
	}
	if req.UserID == "" {
		req.UserID = "default"
	}
	// Scope selects the EPC site/tenant for scoped actions (run_scenario, etc.).
	// Until the two-doors landing passes it explicitly per company, default to
	// the seeded EPC demo site so the Conductor operates on Atlas data.
	if req.Scope == "" {
		req.Scope = "atlas-site"
	}

	// Tier 1: fire an instant intermediary via the fast flash-lite router so the
	// user sees the Conductor dispatching to a specialist while the authoritative
	// (tier-2) response is still generating. Best-effort, never blocks Route.
	if bh.routerClient != nil {
		go bh.emitIntermediary(req.Message, req.Scope)
	}

	// Route through Brain. Allow generous headroom for slower third-party model
	// hosts (e.g. Featherless cold-starts); a warm call should be far faster.
	ctx, cancel := context.WithTimeout(r.Context(), 120*time.Second)
	defer cancel()

	decision, err := bh.router.Route(ctx, req.UserID, req.Message, req.Scope)
	if err != nil {
		log.Printf("[BRAIN] Routing failed: %v", err)
		writeJSON(w, map[string]interface{}{
			"response": fmt.Sprintf("I'm having trouble processing that: %v", err),
			"action":   "error",
			"success":  false,
		})
		return
	}

	// Execute the decision
	result := bh.executor.Execute(decision, req.Scope)

	// Never emit a blank bubble: if the model returned an empty response (e.g. an
	// off-topic prompt it correctly declined to act on), substitute a polite,
	// on-domain redirect so the user never sees an empty message.
	if strings.TrimSpace(result.Response) == "" {
		result.Response = "I'm focused on running the Atlas site — projects, vendors, invoices, schedule, RFIs, the inbox, and hiring. Ask me about any of those and I'll get on it."
		if result.Action == "" || result.Action == brain.ActionRespond {
			result.Action = brain.ActionRespond
		}
	}

	// Broadcast brain activity to WebSocket
	brainEvent, _ := json.Marshal(map[string]interface{}{
		"type": "brain_event",
		"event": map[string]interface{}{
			"action":     result.Action,
			"response":   result.Response,
			"project_id": result.ProjectID,
			"success":    result.Success,
			"timestamp":  time.Now().UnixMilli(),
		},
	})
	bh.hub.BroadcastAgentSessionEvent(brainEvent)

	// Ensure suggestions is never nil. When the Brain doesn't emit any, fall
	// back to chips chosen from the action + the response shape — so a
	// confirmation prompt always gets "Yes / Cancel"-style chips rather
	// than generic tool names.
	suggestions := decision.Suggestions
	if len(suggestions) == 0 {
		suggestions = defaultSuggestionsFor(result.Action, result.Response)
	}

	writeJSON(w, map[string]interface{}{
		"response":    result.Response,
		"action":      result.Action,
		"project_id":  result.ProjectID,
		"success":     result.Success,
		"suggestions": suggestions,
	})
}

// emitIntermediary runs the fast tier-1 router (flash-lite) to pick the likely
// specialist plus a one-line status, then emits it as a message FROM that agent.
// That instantly pulses the agent in the workforce panel and shows a
// "dispatching…" line in the conductor — masking tier-2 latency. Best-effort:
// any failure is silently ignored so the authoritative response is unaffected.
func (bh *BrainHandler) emitIntermediary(userMsg, scope string) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	const sys = `You are a fast dispatcher for an EPC construction-site AI assistant. ` +
		`Given the director's message, reply with ONLY a compact JSON object: ` +
		`{"agent":"<one of: procurement, project_manager, site_engineer, hr, hse, qa_inspector, conductor>", ` +
		`"line":"<a 4-9 word present-tense status, e.g. Routing to Procurement to verify the invoice>"}`
	resp, err := bh.routerClient.SendMessage(ctx, sys, userMsg)
	if err != nil || resp == nil {
		return
	}
	var r struct {
		Agent string `json:"agent"`
		Line  string `json:"line"`
	}
	txt := strings.TrimSpace(resp.Text)
	if i := strings.Index(txt, "{"); i >= 0 {
		txt = txt[i:]
	}
	if json.Unmarshal([]byte(txt), &r) != nil || strings.TrimSpace(r.Line) == "" {
		return
	}
	if r.Agent == "" {
		r.Agent = "conductor"
	}
	m := message.NewMessage(message.TypeSystem, r.Agent, "conductor", r.Line)
	m.Metadata.ProjectID = scope
	m.Metadata.Tags = []string{"routing"} // NOT "scenario" — conductor suppresses scenario-tagged WS dupes
	bh.store.Add(m)
	bh.hub.Broadcast(m)
}

// defaultSuggestionsFor builds context-aware fallback chips when the Brain
// returned an empty `Suggestions` list. Two signals drive the choice:
//   1. action  — what the Brain decided to do (mutation? confirmation?)
//   2. response — does the body actually ask a yes/no question?
// When the response is asking confirmation, we prioritise direct accept /
// cancel chips so the user can resolve the decision in one click without
// retyping. Otherwise we pick chips that follow naturally from the action.
func defaultSuggestionsFor(action, response string) []string {
	asksConfirmation := isConfirmationPrompt(response)

	if asksConfirmation {
		switch action {
		case "delete_email":
			return []string{withCount("Yes, delete", action, response), "No, keep them", "Show inbox"}
		case "archive_email":
			return []string{withCount("Yes, archive", action, response), "Cancel", "Show inbox"}
		case "send_reply":
			return []string{withCount("Yes, send", action, response), "Edit draft first", "Cancel"}
		case "shortlist_applicant":
			return []string{withCount("Yes, shortlist", action, response), "Show their CV", "Cancel"}
		case "escalate":
			return []string{withCount("Yes, escalate", action, response), "Add context first", "Cancel"}
		case "create_project":
			return []string{"Yes, create it", "Adjust scope first", "Cancel"}
		}
		return []string{"Yes, proceed", "Cancel"}
	}

	// Non-confirmation: chips that flow naturally from what just happened.
	switch action {
	case "list_projects":
		return []string{"Show one in detail", "Create a new project", "Check inbox"}
	case "project_status":
		return []string{"List all projects", "Show recent activity", "Check inbox"}
	case "delete_email", "archive_email":
		return []string{"Show inbox", "Process new applicants", "Any other impersonation?"}
	case "send_reply":
		return []string{"Show inbox", "Any other drafts pending?", "Check status"}
	case "shortlist_applicant":
		return []string{"Schedule interview", "Show their CV", "Find more like them"}
	case "set_reminder":
		return []string{"List reminders", "Add another", "Check status"}
	case "respond":
		return []string{"Show inbox", "What needs my attention?", "Check status"}
	}
	return []string{"Show inbox", "List projects", "Check status"}
}

// withCount turns a bare verb phrase like "Yes, delete" into a count-aware
// chip when the response mentions a quantity near a matching noun:
//
//   ("Yes, delete",   "delete_email",        "Found 2 impersonation emails. Confirm?")
//       → "Yes, delete all 2 emails"
//   ("Yes, delete",   "delete_email",        "Confirm deletion of this email?")  (n=1)
//       → "Yes, delete the email"
//   ("Yes, delete",   "delete_email",        "Confirm deletion?")                (no count)
//       → "Yes, delete"
//
// The matched noun is preserved verbatim so the chip reads naturally
// regardless of which synonym the Brain used ("2 messages" → "messages",
// "2 slips" → "slips"). When no count is detected we return the verb
// unchanged.
func withCount(verb, action, response string) string {
	n, noun, ok := extractActionCount(action, response)
	if !ok {
		return verb
	}
	if n == 1 {
		return fmt.Sprintf("%s the %s", verb, noun)
	}
	return fmt.Sprintf("%s all %d %ss", verb, n, noun)
}

// extractActionCount finds a small integer followed (within a short
// window) by a noun that matches the action. Numbers in lines that contain
// "?" win — that's typically the confirmation sentence, where the count is
// most likely what the user is being asked to confirm. Falls back to
// scanning the whole response.
func extractActionCount(action, response string) (int, string, bool) {
	nouns := nounsForAction(action)
	if len(nouns) == 0 {
		return 0, "", false
	}

	// Prefer the confirmation line.
	for _, line := range strings.Split(response, "\n") {
		if !strings.Contains(line, "?") {
			continue
		}
		if n, noun, ok := countNearNoun(line, nouns); ok {
			return n, noun, true
		}
	}
	// Fall back to anywhere in the response.
	return countNearNoun(response, nouns)
}

// nounsForAction lists, in priority order, the nouns we look for near a
// digit for each confirmation action. First match wins, so put the most
// natural synonym first.
func nounsForAction(action string) []string {
	switch action {
	case "delete_email", "archive_email":
		return []string{"email", "message"}
	case "shortlist_applicant":
		return []string{"candidate", "applicant"}
	case "escalate":
		return []string{"slip", "issue", "item"}
	case "send_reply":
		return []string{"draft", "reply", "email"}
	}
	return nil
}

// countNumber matches a standalone 1–3 digit integer — wide enough for
// realistic counts, narrow enough to skip years and long IDs.
var countNumber = regexp.MustCompile(`\b(\d{1,3})\b`)

// countNearNoun returns the first standalone number that's followed within
// `lookahead` characters by any of the given nouns (case-insensitive), and
// the noun that matched (so the caller can echo it in the chip without
// forcing a canonical synonym).
func countNearNoun(text string, nouns []string) (int, string, bool) {
	const lookahead = 40
	for _, m := range countNumber.FindAllStringSubmatchIndex(text, -1) {
		digEnd := m[3]
		n, err := strconv.Atoi(text[m[2]:digEnd])
		if err != nil || n < 1 || n > 100 {
			continue
		}
		end := digEnd + lookahead
		if end > len(text) {
			end = len(text)
		}
		tail := strings.ToLower(text[digEnd:end])
		for _, noun := range nouns {
			if strings.Contains(tail, noun) {
				return n, noun, true
			}
		}
	}
	return 0, "", false
}

// isConfirmationPrompt is a cheap heuristic — the response ends in a
// question mark AND contains a verb that signals a yes/no decision. False
// negatives are fine (we just fall back to action-only chips); false
// positives would be more annoying.
func isConfirmationPrompt(s string) bool {
	if !strings.Contains(s, "?") {
		return false
	}
	low := strings.ToLower(s)
	for _, marker := range []string{
		"confirm", "are you sure", "shall i", "should i",
		"do you want", "proceed?", "go ahead?",
	} {
		if strings.Contains(low, marker) {
			return true
		}
	}
	return false
}

// handleCreateProject is called by the executor when Brain decides to create a project.
func (bh *BrainHandler) handleCreateProject(projectID, taskStr string) error {
	log.Printf("[BRAIN] Creating project: %s — %s", projectID, taskStr[:min(len(taskStr), 80)])

	// Inherit the global default run mode so the build runs at the autonomy
	// level the user picked (manual / semi_auto / full_auto / blitz). New
	// projects with no settings file otherwise default to semi_auto.
	if g, err := checkpoint.LoadSettings(bh.projDir + "/" + globalSettingsProject); err == nil && g.RunMode != "" {
		s := *g
		s.ProjectID = projectID
		s.ApplyRunModePreset(g.RunMode)
		if err := checkpoint.SaveSettings(bh.projDir+"/"+projectID, &s); err != nil {
			log.Printf("[BRAIN] could not seed run mode for %s: %v", projectID, err)
		} else {
			log.Printf("[BRAIN] project %s inherits run mode %q (timeout %dm)", projectID, s.RunMode, s.DecisionTimeoutMinutes)
		}
	}

	// Start task in background (same as handleTask in server.go)
	go func() {
		projDir := bh.projDir + "/" + projectID

		config := orchestrator.Config{
			ProjectDir:   projDir,
			MaxDepth:     5,
			MaxTurns:     25,
			EnableFiles:  true,
			Verbose:      true,
			EnablePhases: true,
		}

		sm := bh.orch.GetSessionManager()
		var newOrch *orchestrator.Orchestrator
		if sm != nil {
			newOrch = orchestrator.NewWithSessions(config, bh.store, sm)
		} else {
			newOrch = orchestrator.New(config, bh.store)
		}
		// Wire callbacks for WebSocket broadcast
		newOrch.OnMessage(func(msg *message.Message) {
			data, _ := json.Marshal(map[string]interface{}{
				"type": "message", "message": msg,
			})
			bh.hub.BroadcastAgentSessionEvent(data)
		})

		taskID := task.GenerateTaskID()
		log.Printf("[BRAIN] Starting orchestration: project=%s task=%s", projectID, taskID)
		newOrch.ProcessTaskWithID(projectID, taskStr, "", taskID)
		log.Printf("[BRAIN] Orchestration complete: project=%s", projectID)
	}()

	return nil
}

// handleDelegate is called by the executor when Brain decides to delegate to an agent.
func (bh *BrainHandler) handleDelegate(projectID, agentRole, taskStr string) error {
	injected := &orchestrator.InjectedTask{
		ID:        fmt.Sprintf("brain_%d", time.Now().UnixMilli()),
		Task:      taskStr,
		AgentRole: agent.Role(agentRole),
		ProjectID: projectID,
		Priority:  "normal",
		CreatedAt: time.Now(),
	}
	bh.orch.InjectTask(injected)
	return nil
}

// handleSendReply sends a reply and marks the original email as replied.
func (bh *BrainHandler) handleSendReply(to, subject, body, emailID string) error {
	log.Printf("[BRAIN] Sending reply to %s: %s (original: %s)", to, subject, emailID)

	// Mark original as replied — try email_id first, then find by sender
	if emailID != "" {
		bh.handleMarkReplied(emailID)
	} else if to != "" {
		// Find email by sender address and mark it
		emails := bh.getEmailsBySender(to)
		for _, eid := range emails {
			bh.handleMarkReplied(eid)
		}
	}
	return nil
}

// getEmailsBySender finds email IDs by sender address.
func (bh *BrainHandler) getEmailsBySender(sender string) []string {
	var ids []string
	inboxDir := bh.projDir + "/inbox"
	entries, _ := os.ReadDir(inboxDir)
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		data, err := os.ReadFile(inboxDir + "/" + entry.Name())
		if err != nil {
			continue
		}
		var email struct {
			ID   string `json:"id"`
			From string `json:"from"`
		}
		if json.Unmarshal(data, &email) == nil && strings.EqualFold(email.From, sender) {
			ids = append(ids, email.ID)
		}
	}
	return ids
}

// handleDeleteEmail deletes an email from the inbox.
func (bh *BrainHandler) handleDeleteEmail(emailID string) error {
	log.Printf("[BRAIN] Deleting email: %s", emailID)
	if bh.emailEngine != nil {
		bh.emailEngine.DeleteEmail(emailID)
	}
	return nil
}

// handleMarkReplied marks an email as replied via the engine (updates memory + disk).
func (bh *BrainHandler) handleMarkReplied(emailID string) {
	log.Printf("[BRAIN] Marking email %s as replied", emailID)
	if bh.emailEngine != nil {
		bh.emailEngine.MarkReplied(emailID)
	}
}

// handleShortlist updates applicant status via email engine.
func (bh *BrainHandler) handleShortlist(applicantID, status string) {
	if bh.emailEngine != nil {
		bh.emailEngine.UpdateApplicantStatus(applicantID, status)
	}
}

// handleBrainChat is the Server method that routes to BrainHandler.
func (s *Server) handleBrainChat(w http.ResponseWriter, r *http.Request) {
	if s.brainHandler == nil {
		http.Error(w, "Brain not initialized", http.StatusServiceUnavailable)
		return
	}
	s.brainHandler.HandleChat(w, r)
}
