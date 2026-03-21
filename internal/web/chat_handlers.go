package web

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"pty-claude-test/internal/agent"
	"pty-claude-test/internal/chat"
)

// handleAgentChat handles GET/POST /api/agents/{role}/chat
func (s *Server) handleAgentChat(w http.ResponseWriter, r *http.Request) {
	// Parse role from path: /api/agents/{role}/chat
	path := strings.TrimPrefix(r.URL.Path, "/api/agents/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[1] != "chat" {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	roleStr := parts[0]

	projectID := r.URL.Query().Get("project")
	if projectID == "" {
		projectID = "default"
	}
	projectDir := filepath.Join(s.projectDir, projectID)

	switch r.Method {
	case http.MethodGet:
		limit := 50
		msgs, err := chat.LoadMessages(projectDir, roleStr, limit)
		if err != nil {
			msgs = []chat.ChatMessage{}
		}
		writeJSON(w, map[string]interface{}{
			"messages": msgs,
			"count":    len(msgs),
		})

	case http.MethodPost:
		var req struct {
			Content string `json:"content"`
			Context string `json:"context"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if req.Content == "" {
			http.Error(w, "Content is required", http.StatusBadRequest)
			return
		}
		if req.Context == "" {
			req.Context = "idle_chat"
		}

		// Save user message
		userMsg := chat.ChatMessage{
			ID:        chat.GenerateMessageID(),
			ProjectID: projectID,
			AgentRole: roleStr,
			Direction: "user_to_agent",
			Content:   req.Content,
			Context:   req.Context,
			Timestamp: time.Now(),
		}
		if err := chat.AddMessage(projectDir, userMsg); err != nil {
			http.Error(w, "Failed to save message", http.StatusInternalServerError)
			return
		}

		// Broadcast user message via WebSocket
		s.hub.BroadcastCheckpointEvent(&CheckpointEvent{
			EventType:      "chat_message",
			CheckpointType: roleStr,
			Data: map[string]interface{}{
				"direction": "user_to_agent",
				"content":   req.Content,
				"agent":     roleStr,
				"id":        userMsg.ID,
			},
		})

		// Spawn agent response in background
		go s.processAgentChat(projectDir, projectID, roleStr, req.Content)

		writeJSON(w, map[string]interface{}{
			"success":    true,
			"message_id": userMsg.ID,
			"status":     "processing",
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// processAgentChat spawns an agent to respond to a chat message
func (s *Server) processAgentChat(projectDir, projectID, roleStr, userMessage string) {
	role := agent.Role(roleStr)

	// Build prompt with conversation context
	prompt, err := chat.BuildConversationContext(projectDir, roleStr, userMessage)
	if err != nil {
		log.Printf("[CHAT] Failed to build context for %s: %v", roleStr, err)
		return
	}

	// Create agent
	a, err := agent.NewAgent(role)
	if err != nil {
		log.Printf("[CHAT] Failed to create agent %s: %v", roleStr, err)
		return
	}

	log.Printf("[CHAT] %s processing chat message", roleStr)

	// Process — read-only mode (ProcessInDir doesn't write files unless explicitly asked)
	response, err := a.ProcessInDir(prompt, projectDir)
	if err != nil {
		log.Printf("[CHAT] %s failed: %v", roleStr, err)
		// Save error as response
		errMsg := chat.ChatMessage{
			ID:        chat.GenerateMessageID(),
			ProjectID: projectID,
			AgentRole: roleStr,
			Direction: "agent_to_user",
			Content:   "Sorry, I encountered an error processing your message.",
			Context:   "error",
			Timestamp: time.Now(),
		}
		chat.AddMessage(projectDir, errMsg)
		return
	}

	// Save agent response
	agentMsg := chat.ChatMessage{
		ID:        chat.GenerateMessageID(),
		ProjectID: projectID,
		AgentRole: roleStr,
		Direction: "agent_to_user",
		Content:   response.Content,
		Context:   "idle_chat",
		Timestamp: time.Now(),
	}
	if err := chat.AddMessage(projectDir, agentMsg); err != nil {
		log.Printf("[CHAT] Failed to save response: %v", err)
	}

	log.Printf("[CHAT] %s responded (%d bytes)", roleStr, len(response.Content))

	// Broadcast agent response via WebSocket
	s.hub.BroadcastCheckpointEvent(&CheckpointEvent{
		EventType:      "chat_response",
		CheckpointType: roleStr,
		Data: map[string]interface{}{
			"direction": "agent_to_user",
			"content":   response.Content,
			"agent":     roleStr,
			"id":        agentMsg.ID,
		},
	})
}
