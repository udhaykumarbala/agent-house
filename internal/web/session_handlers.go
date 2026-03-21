package web

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// handleSessionRouting routes /api/sessions/* requests.
func (s *Server) handleSessionRouting(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/sessions/")
	parts := strings.Split(path, "/")

	// POST actions: /api/sessions/{projectID}/{agentRole}/{action}
	if r.Method == http.MethodPost && len(parts) >= 3 {
		s.handleAgentSessionAction(w, r)
		return
	}

	// GET with projectID
	if r.Method == http.MethodGet {
		s.handleProjectSessions(w, r)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// handleSessions handles GET /api/sessions — list all agent sessions
func (s *Server) handleSessions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sm := s.orchestrator.GetSessionManager()
	log.Printf("[SESSIONS] GetSessionManager returned: %v", sm != nil)
	if sm == nil {
		writeJSON(w, map[string]interface{}{
			"sessions": []interface{}{},
			"message":  "Session mode not enabled",
		})
		return
	}

	sessions := sm.ListAll()
	writeJSON(w, map[string]interface{}{
		"sessions": sessions,
		"count":    len(sessions),
	})
}

// handleProjectSessions handles GET /api/sessions/{projectID} — list sessions for a project
func (s *Server) handleProjectSessions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sm := s.orchestrator.GetSessionManager()
	if sm == nil {
		writeJSON(w, map[string]interface{}{"sessions": []interface{}{}})
		return
	}

	// Extract projectID from path: /api/sessions/{projectID}
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/sessions/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		http.Error(w, "Project ID required", http.StatusBadRequest)
		return
	}
	projectID := parts[0]

	// If there's also an agent role: /api/sessions/{projectID}/{agentRole}
	if len(parts) >= 2 && parts[1] != "" {
		s.handleAgentSessionDetail(w, r, projectID, parts[1])
		return
	}

	sessions := sm.ListByProject(projectID)
	list := make([]interface{}, 0, len(sessions))
	for _, sess := range sessions {
		list = append(list, sess.GetInfo())
	}

	writeJSON(w, map[string]interface{}{
		"sessions":   list,
		"count":      len(list),
		"project_id": projectID,
	})
}

// handleAgentSessionDetail handles GET /api/sessions/{projectID}/{agentRole}
func (s *Server) handleAgentSessionDetail(w http.ResponseWriter, r *http.Request, projectID, agentRole string) {
	sm := s.orchestrator.GetSessionManager()
	if sm == nil {
		http.Error(w, "Session mode not enabled", http.StatusServiceUnavailable)
		return
	}

	sess := sm.Get(projectID, agentRole)
	if sess == nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	writeJSON(w, sess.GetInfo())
}

// handleAgentSessionAction handles POST actions on agent sessions
// POST /api/sessions/{projectID}/{agentRole}/abort
// POST /api/sessions/{projectID}/{agentRole}/approve
// POST /api/sessions/{projectID}/{agentRole}/deny
func (s *Server) handleAgentSessionAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sm := s.orchestrator.GetSessionManager()
	if sm == nil {
		http.Error(w, "Session mode not enabled", http.StatusServiceUnavailable)
		return
	}

	// Parse path: /api/sessions/{projectID}/{agentRole}/{action}
	path := strings.TrimPrefix(r.URL.Path, "/api/sessions/")
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	projectID := parts[0]
	agentRole := parts[1]
	action := parts[2]

	sess := sm.Get(projectID, agentRole)
	if sess == nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	switch action {
	case "abort":
		if err := sess.Abort(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, map[string]string{"status": "aborted"})

	case "approve":
		var body struct {
			ToolUseID string `json:"tool_use_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid body", http.StatusBadRequest)
			return
		}
		if err := sess.ApproveTool(body.ToolUseID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, map[string]string{"status": "approved"})

	case "deny":
		var body struct {
			ToolUseID string `json:"tool_use_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid body", http.StatusBadRequest)
			return
		}
		if err := sess.DenyTool(body.ToolUseID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, map[string]string{"status": "denied"})

	default:
		http.Error(w, "Unknown action: "+action, http.StatusBadRequest)
	}
}

// handleProjectReport handles GET /api/reports/{projectID}
func (s *Server) handleProjectReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sm := s.orchestrator.GetSessionManager()
	if sm == nil {
		http.Error(w, "Session mode not enabled", http.StatusServiceUnavailable)
		return
	}

	projectID := strings.TrimPrefix(r.URL.Path, "/api/reports/")
	if projectID == "" {
		http.Error(w, "Project ID required", http.StatusBadRequest)
		return
	}

	metrics := sm.GetProjectMetrics(projectID)
	writeJSON(w, metrics)
}
