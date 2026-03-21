package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"pty-claude-test/internal/agent"
	"pty-claude-test/internal/orchestrator"
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

// handleInjectTask handles POST /api/inject — inject a task into the running pipeline
func (s *Server) handleInjectTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Task      string `json:"task"`
		AgentRole string `json:"agent_role"` // ceo, pm, senior_dev, etc. Empty = auto
		ProjectID string `json:"project_id"`
		Priority  string `json:"priority"`   // high, normal
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	if req.Task == "" {
		http.Error(w, "task is required", http.StatusBadRequest)
		return
	}
	if req.ProjectID == "" {
		req.ProjectID = "default"
	}
	if req.Priority == "" {
		req.Priority = "normal"
	}

	injected := &orchestrator.InjectedTask{
		ID:        fmt.Sprintf("inj_%d", time.Now().UnixMilli()),
		Task:      req.Task,
		AgentRole: agent.Role(req.AgentRole),
		ProjectID: req.ProjectID,
		Priority:  req.Priority,
		CreatedAt: time.Now(),
	}

	s.orchestrator.InjectTask(injected)

	writeJSON(w, map[string]interface{}{
		"success":   true,
		"inject_id": injected.ID,
		"message":   fmt.Sprintf("Task injected for %s — will execute between phases", req.AgentRole),
	})
}

// handleProjectDetail handles GET /api/projects/{projectID} — full project summary
func (s *Server) handleProjectDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	projectID := strings.TrimPrefix(r.URL.Path, "/api/projects/")
	if projectID == "" {
		http.Error(w, "Project ID required", http.StatusBadRequest)
		return
	}

	projectDir := s.projectDir + "/" + projectID

	// Collect project files (non-hidden)
	var files []map[string]interface{}
	filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			base := filepath.Base(path)
			if strings.HasPrefix(base, ".") && path != projectDir {
				return filepath.SkipDir
			}
			if base == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}
		rel, _ := filepath.Rel(projectDir, path)
		files = append(files, map[string]interface{}{
			"path": rel,
			"size": info.Size(),
		})
		return nil
	})

	// Get task history
	var tasks []map[string]interface{}
	historyPath := projectDir + "/.tasks/history.json"
	if data, err := os.ReadFile(historyPath); err == nil {
		json.Unmarshal(data, &tasks)
	}

	// Get session metrics
	var sessionMetrics map[string]interface{}
	sm := s.orchestrator.GetSessionManager()
	if sm != nil {
		sessionMetrics = sm.GetProjectMetrics(projectID)
	}

	// Get development plan
	var devPlan interface{}
	planPath := projectDir + "/.plans/development-plan.json"
	if data, err := os.ReadFile(planPath); err == nil {
		json.Unmarshal(data, &devPlan)
	}

	writeJSON(w, map[string]interface{}{
		"project_id": projectID,
		"files":      files,
		"file_count": len(files),
		"tasks":      tasks,
		"sessions":   sessionMetrics,
		"dev_plan":   devPlan,
	})
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
