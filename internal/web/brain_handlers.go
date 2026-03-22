package web

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"pty-claude-test/internal/agent"
	"pty-claude-test/internal/brain"
	"pty-claude-test/internal/message"
	"pty-claude-test/internal/orchestrator"
	"pty-claude-test/internal/session"
	"pty-claude-test/internal/task"
)

// BrainHandler manages the Brain chat interface.
type BrainHandler struct {
	router       *brain.Router
	executor     *brain.Executor
	convos       *brain.ConversationStore
	orch         *orchestrator.Orchestrator
	store        *message.Store
	projDir      string
	hub          *Hub
	emailEngine  interface{ MarkReplied(string); DeleteEmail(string) bool; UpdateApplicantStatus(string,string) bool }
}

// NewBrainHandler creates the Brain handler with all dependencies wired.
func NewBrainHandler(apiClient *session.APIClient, orch *orchestrator.Orchestrator, store *message.Store, hub *Hub, projDir string) *BrainHandler {
	convos := brain.NewConversationStore(projDir + "/.brain/conversations")
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

	// Route through Brain
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	decision, err := bh.router.Route(ctx, req.UserID, req.Message)
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
	result := bh.executor.Execute(decision)

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

	// Ensure suggestions is never nil
	suggestions := decision.Suggestions
	if len(suggestions) == 0 {
		suggestions = []string{"Show inbox", "List projects", "Check status"}
	}

	writeJSON(w, map[string]interface{}{
		"response":    result.Response,
		"action":      result.Action,
		"project_id":  result.ProjectID,
		"success":     result.Success,
		"suggestions": suggestions,
	})
}

// handleCreateProject is called by the executor when Brain decides to create a project.
func (bh *BrainHandler) handleCreateProject(projectID, taskStr string) error {
	log.Printf("[BRAIN] Creating project: %s — %s", projectID, taskStr[:min(len(taskStr), 80)])

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
