package web

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	router   *brain.Router
	executor *brain.Executor
	convos   *brain.ConversationStore
	orch     *orchestrator.Orchestrator
	store    *message.Store
	projDir  string
	hub      *Hub
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

	writeJSON(w, map[string]interface{}{
		"response":    result.Response,
		"action":      result.Action,
		"project_id":  result.ProjectID,
		"success":     result.Success,
		"suggestions": decision.Suggestions,
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

// handleSendReply sends a reply via the email handlers.
func (bh *BrainHandler) handleSendReply(to, subject, body string) error {
	log.Printf("[BRAIN] Sending reply to %s: %s", to, subject)
	// Would use Resend here — for now just log it
	return nil
}

// handleBrainChat is the Server method that routes to BrainHandler.
func (s *Server) handleBrainChat(w http.ResponseWriter, r *http.Request) {
	if s.brainHandler == nil {
		http.Error(w, "Brain not initialized", http.StatusServiceUnavailable)
		return
	}
	s.brainHandler.HandleChat(w, r)
}
