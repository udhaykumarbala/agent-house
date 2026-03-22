package web

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"pty-claude-test/internal/agent"
	"pty-claude-test/internal/agentask"
	"pty-claude-test/internal/message"
	"pty-claude-test/internal/orchestrator"
	"pty-claude-test/internal/session"
	"pty-claude-test/internal/task"
)

// Server handles HTTP requests for the dashboard
type Server struct {
	store           *message.Store
	orchestrator    *orchestrator.Orchestrator
	projectDir      string
	mu              sync.RWMutex
	projectTasks    map[string]bool // projectID -> isRunning (allows concurrent projects)
	hub             *Hub
	historyManager  *task.HistoryManager
	agentTaskMgr    *agentask.Manager // Agent task tracking
	cronScheduler   *CronScheduler     // Cron scheduler for recurring tasks
	brainHandler    *BrainHandler      // Brain chat handler
	emailHandlers   *EmailHandlers     // Email engine handlers
}

// Config holds server configuration
type Config struct {
	Port       int
	ProjectDir string
	Store      *message.Store
}

// NewServer creates a new web server
func NewServer(config Config) *Server {
	orchConfig := orchestrator.Config{
		ProjectDir:   config.ProjectDir,
		MaxDepth:     5,
		MaxTurns:     25, // Increased for 4-phase workflow (5+5+1+2 = 13 agents minimum)
		EnableFiles:  true,
		Verbose:      true,
		EnablePhases: true, // Enable multi-phase Research → Plan → Discuss → Develop workflow
	}

	hub := NewHub()

	// Create session manager for Claude Code sessions per agent
	sessionMgr := session.NewSessionManager()
	log.Printf("[SERVER] Session manager created: %v", sessionMgr != nil)

	server := &Server{
		store:          config.Store,
		orchestrator:   orchestrator.NewWithSessions(orchConfig, config.Store, sessionMgr),
		projectDir:     config.ProjectDir,
		projectTasks:   make(map[string]bool),
		hub:            hub,
		historyManager: task.NewHistoryManager(),
		agentTaskMgr:   agentask.NewManager(config.ProjectDir),
	}

	// Wire session manager events to WebSocket hub
	go func() {
		eventCh := sessionMgr.SubscribeAll()
		for ev := range eventCh {
			data, err := json.Marshal(map[string]interface{}{
				"type":  "agent_event",
				"event": ev,
			})
			if err != nil {
				continue
			}
			hub.BroadcastAgentSessionEvent(data)
		}
	}()

	// Load file-based agent registry
	agentRegistry := agent.NewRegistry("agents")
	server.orchestrator.SetAgentRegistry(agentRegistry)

	// Initialize email handlers
	server.emailHandlers = NewEmailHandlers(config.ProjectDir)

	// Initialize Brain handler
	apiClient := session.NewAPIClient()
	server.brainHandler = NewBrainHandler(apiClient, server.orchestrator, config.Store, hub, config.ProjectDir)

	// Start cron scheduler
	server.cronScheduler = NewCronScheduler(server.orchestrator)

	// Wire WebSocket session action handler (approve/deny/abort via WS)
	hub.SetSessionActionHandler(func(projectID, agentRole, action, toolUseID string) {
		sess := sessionMgr.Get(projectID, agentRole)
		if sess == nil {
			return
		}
		switch action {
		case "approve":
			sess.ApproveTool(toolUseID)
		case "deny":
			sess.DenyTool(toolUseID)
		case "abort":
			sess.Abort()
		}
	})

	// Register orchestrator message callback for WebSocket broadcast
	server.orchestrator.OnMessage(func(msg *message.Message) {
		hub.Broadcast(msg)
	})

	// Register checkpoint callback for WebSocket broadcast
	server.orchestrator.SetCheckpointNotify(func(cpID, cpType, summary string, phaseIndex int, resolved bool, decidedBy, action string) {
		eventType := "checkpoint_reached"
		if resolved {
			if decidedBy == "ceo_auto" {
				eventType = "checkpoint_auto_delegated"
			} else {
				eventType = "checkpoint_resolved"
			}
		}
		hub.BroadcastCheckpointEvent(&CheckpointEvent{
			EventType:      eventType,
			CheckpointID:   cpID,
			CheckpointType: cpType,
			PhaseIndex:     phaseIndex,
			Data: map[string]interface{}{
				"artifact_summary": summary,
				"decided_by":       decidedBy,
				"action":           action,
			},
		})
	})

	return server
}

// Start starts the HTTP server
func (s *Server) Start(port int) error {
	mux := http.NewServeMux()

	// Start WebSocket hub
	go s.hub.Run()

	// API routes
	mux.HandleFunc("/api/messages", s.handleMessages)
	mux.HandleFunc("/api/agents", s.handleAgents)
	mux.HandleFunc("/api/task", s.handleTask)
	mux.HandleFunc("/api/tasks", s.handleTasks)
	mux.HandleFunc("/api/task/", s.handleTaskByID)
	mux.HandleFunc("/api/files", s.handleFiles)
	mux.HandleFunc("/api/status", s.handleStatus)
	mux.HandleFunc("/api/phase", s.handlePhase)

	// Agent task endpoints
	mux.HandleFunc("/api/agent-tasks", s.handleAgentTasks)
	mux.HandleFunc("/api/agent-tasks/", s.handleAgentTaskByID)
	mux.HandleFunc("/api/agents/", s.handleAgentEndpoints)

	// Projects endpoint
	mux.HandleFunc("/api/projects", s.handleProjects)

	// File content endpoint
	mux.HandleFunc("/api/file-content", s.handleFileContent)

	// Kanban endpoints
	mux.HandleFunc("/api/kanban/tasks", s.handleKanbanTasks)
	mux.HandleFunc("/api/kanban/tasks/", s.handleKanbanTaskByID)
	mux.HandleFunc("/api/kanban/sync", s.handleKanbanSync)

	// Checkpoint endpoints
	mux.HandleFunc("/api/checkpoints", s.handleCheckpoints)
	mux.HandleFunc("/api/checkpoints/", s.handleCheckpointByID)

	// Workflow settings endpoints
	mux.HandleFunc("/api/settings/workflow", s.handleWorkflowSettings)

	// Development plan endpoints
	mux.HandleFunc("/api/development-plan", s.handleDevelopmentPlan)
	mux.HandleFunc("/api/subtasks", s.handleSubTasks)
	mux.HandleFunc("/api/qa-reviews", s.handleQAReviews)
	mux.HandleFunc("/api/phase-status", s.handlePhaseStatus)

	// Project detail endpoint
	mux.HandleFunc("/api/projects/", s.handleProjectDetail)

	// Email endpoints
	mux.HandleFunc("/api/email/", s.handleEmail)

	// Brain chat endpoint
	mux.HandleFunc("/api/chat", s.handleBrainChat)

	// Agent mode management
	mux.HandleFunc("/api/agents/modes", s.handleAgentModes)

	// Inject task endpoint
	mux.HandleFunc("/api/inject", s.handleInjectTask)

	// Webhook hooks — external systems trigger agents
	mux.HandleFunc("/api/hooks/", s.handleWebhook)

	// Cron scheduler — recurring agent tasks
	mux.HandleFunc("/api/cron", s.handleCronRouting)
	mux.HandleFunc("/api/cron/", s.handleCronRouting)

	// Agent session endpoints
	mux.HandleFunc("/api/sessions", s.handleSessions)
	mux.HandleFunc("/api/sessions/", s.handleSessionRouting)

	// Project report endpoints
	mux.HandleFunc("/api/reports/", s.handleProjectReport)

	// WebSocket endpoint
	mux.HandleFunc("/ws", s.hub.ServeWS)

	// Static files
	mux.HandleFunc("/", s.handleStatic)

	// CORS middleware
	handler := corsMiddleware(mux)

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting server on http://localhost%s", addr)
	log.Printf("WebSocket available at ws://localhost%s/ws", addr)

	return http.ListenAndServe(addr, handler)
}

// corsMiddleware adds CORS headers
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// handleMessages returns all messages
func (s *Server) handleMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check for query params
	projectID := r.URL.Query().Get("project")
	agentName := r.URL.Query().Get("agent")
	afterStr := r.URL.Query().Get("after")

	var messages []*message.Message

	if afterStr != "" {
		// Parse timestamp
		afterTime, err := time.Parse(time.RFC3339, afterStr)
		if err != nil {
			messages = s.store.GetAll()
		} else {
			messages = s.store.GetAfter(afterTime)
		}
	} else if projectID != "" {
		messages = s.store.GetByProject(projectID)
	} else if agentName != "" {
		messages = s.store.GetByAgent(agentName)
	} else {
		messages = s.store.GetAll()
	}

	writeJSON(w, map[string]interface{}{
		"messages": messages,
		"count":    len(messages),
	})
}

// handleAgents returns agent information
func (s *Server) handleAgents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get active agents from orchestrator
	activeAgents := s.orchestrator.GetAgents()

	// Build response with all roles
	roles := []agent.Role{
		agent.RoleCEO, agent.RolePM, agent.RoleUX, agent.RoleUI,
		agent.RoleSecurity, agent.RoleArchitect, agent.RoleSeniorDev, agent.RoleJuniorDev,
	}

	agentInfos := make([]map[string]interface{}, 0)
	for _, role := range roles {
		info := map[string]interface{}{
			"role":        string(role),
			"name":        getRoleName(role),
			"permissions": agent.GetPermission(role),
			"active":      false,
			"color":       getRoleColor(role),
		}

		if a, ok := activeAgents[role]; ok {
			info["active"] = true
			info["name"] = a.Name
		}

		agentInfos = append(agentInfos, info)
	}

	writeJSON(w, map[string]interface{}{
		"agents": agentInfos,
	})
}

// TaskRequest represents a task submission
type TaskRequest struct {
	Task         string `json:"task"`
	ProjectID    string `json:"project_id"`
	ContinueFrom string `json:"continue_from,omitempty"` // Previous task ID to continue from
}

// handleTask handles task submission
func (s *Server) handleTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request
	var req TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Task == "" {
		http.Error(w, "Task is required", http.StatusBadRequest)
		return
	}

	if req.ProjectID == "" {
		req.ProjectID = "default"
	}

	// Check if task is already running for this project
	s.mu.Lock()
	if s.projectTasks[req.ProjectID] {
		s.mu.Unlock()
		writeJSON(w, map[string]interface{}{
			"error":   fmt.Sprintf("Task already running for project %s", req.ProjectID),
			"success": false,
		})
		return
	}
	// Mark project as running
	s.projectTasks[req.ProjectID] = true
	s.mu.Unlock()

	// Generate task ID upfront so we can return it
	taskID := task.GenerateTaskID()

	log.Printf("[API] Task submitted: id=%s project=%s task=%q", taskID, req.ProjectID, req.Task)

	// Start task in background
	go func(taskID, projectID string) {
		defer func() {
			s.mu.Lock()
			delete(s.projectTasks, projectID)
			s.mu.Unlock()
		}()

		projectDir := filepath.Join(s.projectDir, req.ProjectID)
		os.MkdirAll(projectDir, 0755)

		// Build task with continuation context if provided
		taskWithContext := req.Task
		var parentTaskID string

		if req.ContinueFrom != "" {
			parentTaskID = req.ContinueFrom
			prevMeta, prevMessages, err := s.historyManager.GetTask(projectDir, req.ContinueFrom)
			if err == nil {
				taskWithContext = buildContinuationContext(prevMeta, prevMessages, req.Task)
			}
		}

		// Update orchestrator project dir, preserving session manager
		config := orchestrator.Config{
			ProjectDir:   projectDir,
			MaxDepth:     5,
			MaxTurns:     25, // Increased for 4-phase workflow
			EnableFiles:  true,
			Verbose:      true,
			EnablePhases: true, // Enable multi-phase workflow
		}
		sm := s.orchestrator.GetSessionManager()
		if sm != nil {
			s.orchestrator = orchestrator.NewWithSessions(config, s.store, sm)
		} else {
			s.orchestrator = orchestrator.New(config, s.store)
		}

		// Re-register WebSocket callbacks on new orchestrator
		s.orchestrator.OnMessage(func(msg *message.Message) {
			s.hub.Broadcast(msg)
		})
		s.orchestrator.SetCheckpointNotify(func(cpID, cpType, summary string, phaseIndex int, resolved bool, decidedBy, action string) {
			eventType := "checkpoint_reached"
			if resolved {
				if decidedBy == "ceo_auto" {
					eventType = "checkpoint_auto_delegated"
				} else {
					eventType = "checkpoint_resolved"
				}
			}
			s.hub.BroadcastCheckpointEvent(&CheckpointEvent{
				EventType:      eventType,
				CheckpointID:   cpID,
				CheckpointType: cpType,
				PhaseIndex:     phaseIndex,
				Data: map[string]interface{}{
					"artifact_summary": summary,
					"decided_by":       decidedBy,
					"action":           action,
				},
			})
		})

		log.Printf("[TASK] Starting orchestration: id=%s project=%s", taskID, projectID)
		s.orchestrator.ProcessTaskWithID(req.ProjectID, taskWithContext, parentTaskID, taskID)
		log.Printf("[TASK] Orchestration finished: id=%s project=%s", taskID, projectID)
	}(taskID, req.ProjectID)

	writeJSON(w, map[string]interface{}{
		"success":    true,
		"message":    "Task started",
		"project_id": req.ProjectID,
		"task_id":    taskID,
	})
}

// FileInfo represents a file in the project
type FileInfo struct {
	Path      string `json:"path"`
	Name      string `json:"name"`
	IsDir     bool   `json:"is_dir"`
	Size      int64  `json:"size"`
	Modified  string `json:"modified"`
	CreatedBy string `json:"created_by,omitempty"` // Agent that created this file
}

// handleFiles returns project files
func (s *Server) handleFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	projectID := r.URL.Query().Get("project")
	if projectID == "" {
		projectID = "default"
	}

	projectPath := filepath.Join(s.projectDir, projectID)

	// Ensure directory exists
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		writeJSON(w, map[string]interface{}{
			"files":   []FileInfo{},
			"project": projectID,
		})
		return
	}

	files := make([]FileInfo, 0)

	filepath.WalkDir(projectPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		// Skip hidden files and directories
		if strings.HasPrefix(d.Name(), ".") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		relPath, _ := filepath.Rel(projectPath, path)
		if relPath == "." {
			return nil
		}

		info, _ := d.Info()
		var size int64
		var modified string
		if info != nil {
			size = info.Size()
			modified = info.ModTime().Format(time.RFC3339)
		}

		files = append(files, FileInfo{
			Path:     relPath,
			Name:     d.Name(),
			IsDir:    d.IsDir(),
			Size:     size,
			Modified: modified,
		})

		return nil
	})

	writeJSON(w, map[string]interface{}{
		"files":   files,
		"project": projectID,
	})
}

// handleStatus returns server status
func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.mu.RLock()
	runningProjects := make([]string, 0, len(s.projectTasks))
	for projectID, running := range s.projectTasks {
		if running {
			runningProjects = append(runningProjects, projectID)
		}
	}
	anyRunning := len(runningProjects) > 0
	s.mu.RUnlock()

	wsClients := 0
	if s.hub != nil {
		wsClients = s.hub.ClientCount()
	}

	writeJSON(w, map[string]interface{}{
		"status":           "ok",
		"task_running":     anyRunning, // For backward compatibility
		"running_projects": runningProjects,
		"message_count":    s.store.Count(),
		"project_dir":      s.projectDir,
		"ws_clients":       wsClients,
	})
}

// PhaseInfo represents phase information for the API
type PhaseInfo struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"` // pending, active, completed
}

// handlePhase returns the current workflow phase
func (s *Server) handlePhase(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get current phase from orchestrator
	currentPhase := s.orchestrator.GetCurrentPhase()

	// Define all phases in order
	phaseOrder := []string{"template_selection", "research", "planning", "discussion", "development"}
	phaseNames := map[string]string{
		"template_selection": "Template Selection",
		"research":           "Research",
		"planning":           "Planning",
		"discussion":         "Discussion",
		"development":        "Development",
	}

	// Build phases with status
	phases := make([]PhaseInfo, len(phaseOrder))
	currentIdx := -1

	// Find current phase index
	for i, p := range phaseOrder {
		if p == currentPhase {
			currentIdx = i
			break
		}
	}

	for i, p := range phaseOrder {
		status := "pending"
		if currentIdx >= 0 {
			if i < currentIdx {
				status = "completed"
			} else if i == currentIdx {
				status = "active"
			}
		}
		phases[i] = PhaseInfo{
			ID:     p,
			Name:   phaseNames[p],
			Status: status,
		}
	}

	writeJSON(w, map[string]interface{}{
		"current_phase": currentPhase,
		"phases":        phases,
	})
}

// handleTasks returns the list of tasks for a project
func (s *Server) handleTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	projectID := r.URL.Query().Get("project")
	if projectID == "" {
		projectID = "default"
	}

	projectPath := filepath.Join(s.projectDir, projectID)

	history, err := s.historyManager.GetTaskList(projectPath)
	if err != nil {
		writeJSON(w, map[string]interface{}{
			"error":   err.Error(),
			"success": false,
		})
		return
	}

	writeJSON(w, map[string]interface{}{
		"project": projectID,
		"tasks":   history.Tasks,
		"count":   len(history.Tasks),
	})
}

// handleTaskByID returns a specific task's details and conversation
func (s *Server) handleTaskByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract task ID from URL: /api/task/{taskId}
	path := strings.TrimPrefix(r.URL.Path, "/api/task/")
	taskID := strings.TrimSuffix(path, "/")

	if taskID == "" {
		http.Error(w, "Task ID required", http.StatusBadRequest)
		return
	}

	projectID := r.URL.Query().Get("project")
	if projectID == "" {
		projectID = "default"
	}

	projectPath := filepath.Join(s.projectDir, projectID)

	meta, messages, err := s.historyManager.GetTask(projectPath, taskID)
	if err != nil {
		writeJSON(w, map[string]interface{}{
			"error":   err.Error(),
			"success": false,
		})
		return
	}

	writeJSON(w, map[string]interface{}{
		"task":     meta,
		"messages": messages,
		"success":  true,
	})
}

// handleStatic serves static files
func (s *Server) handleStatic(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" || r.URL.Path == "/index.html" {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(dashboardHTML))
		return
	}

	if r.URL.Path == "/office" {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(officeHTML))
		return
	}

	if r.URL.Path == "/mission" {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(missionHTML))
		return
	}

	if r.URL.Path == "/live" {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(agentActivityHTML))
		return
	}

	if strings.HasPrefix(r.URL.Path, "/project/") {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(projectDetailHTML))
		return
	}

	http.NotFound(w, r)
}

// buildContinuationContext creates context from a previous task for continuation
func buildContinuationContext(prevMeta *task.TaskMetadata, prevMessages []*message.Message, newTask string) string {
	// Build a summary of what was done before
	var filesCreated string
	if len(prevMeta.FilesCreated) > 0 {
		filesCreated = "\n- " + strings.Join(prevMeta.FilesCreated, "\n- ")
	} else {
		filesCreated = "(none)"
	}

	// Extract key decisions from previous messages (simplified - just get agent responses)
	var previousWork strings.Builder
	for _, msg := range prevMessages {
		if msg.Type == message.TypeResponse {
			// Truncate long responses
			content := msg.Content
			if len(content) > 500 {
				content = content[:500] + "..."
			}
			previousWork.WriteString(fmt.Sprintf("\n[%s]: %s\n", msg.From, content))
		}
	}

	return fmt.Sprintf(`CONTINUATION OF PREVIOUS TASK
=============================

PREVIOUS TASK: %s
STATUS: %s
FILES CREATED: %s

PREVIOUS WORK SUMMARY:
%s

=============================
NEW REQUEST: %s

Please continue from where we left off. The files listed above already exist in the project.
Modify or add to them as needed to fulfill the new request.`,
		prevMeta.Task,
		prevMeta.Status,
		filesCreated,
		previousWork.String(),
		newTask,
	)
}

// writeJSON writes JSON response
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// getRoleName returns a display name for the role
func getRoleName(role agent.Role) string {
	names := map[agent.Role]string{
		agent.RoleCEO:       "Chief Executive Officer",
		agent.RolePM:        "Product Manager",
		agent.RoleUX:        "UX Designer",
		agent.RoleUI:        "UI Designer",
		agent.RoleSecurity:  "Security Expert",
		agent.RoleArchitect: "Software Architect",
		agent.RoleSeniorDev: "Senior Developer",
		agent.RoleJuniorDev: "Junior Developer",
	}
	if name, ok := names[role]; ok {
		return name
	}
	return string(role)
}

// getRoleColor returns a color for the role
func getRoleColor(role agent.Role) string {
	colors := map[agent.Role]string{
		agent.RoleCEO:       "#4A90A4",
		agent.RolePM:        "#7B68EE",
		agent.RoleUX:        "#FF6B6B",
		agent.RoleUI:        "#4ECDC4",
		agent.RoleSecurity:  "#F39C12",
		agent.RoleArchitect: "#9B59B6",
		agent.RoleSeniorDev: "#2ECC71",
		agent.RoleJuniorDev: "#3498DB",
	}
	if color, ok := colors[role]; ok {
		return color
	}
	return "#888888"
}

// handleAgentTasks handles GET /api/agent-tasks
func (s *Server) handleAgentTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get query parameters
	projectID := r.URL.Query().Get("project")
	roleStr := r.URL.Query().Get("role")
	statusStr := r.URL.Query().Get("status")

	// Convert role string to agent.Role
	var role agent.Role
	if roleStr != "" {
		role = agent.Role(roleStr)
	}

	// Convert status string to AgentTaskStatus
	var status agentask.AgentTaskStatus
	if statusStr != "" {
		status = agentask.AgentTaskStatus(statusStr)
	}

	// Filter tasks
	tasks := s.agentTaskMgr.FilterTasks(projectID, role, status)

	writeJSON(w, map[string]interface{}{
		"tasks": tasks,
		"count": len(tasks),
	})
}

// handleAgentTaskByID handles GET /api/agent-tasks/{taskID}
func (s *Server) handleAgentTaskByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract task ID from path
	path := strings.TrimPrefix(r.URL.Path, "/api/agent-tasks/")
	taskID := strings.TrimSuffix(path, "/")

	if taskID == "" {
		http.Error(w, "Task ID required", http.StatusBadRequest)
		return
	}

	// Get task
	task, err := s.agentTaskMgr.GetTask(taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, task)
}

// handleAgentEndpoints handles /api/agents/{role}/* endpoints
func (s *Server) handleAgentEndpoints(w http.ResponseWriter, r *http.Request) {
	// Parse path: /api/agents/{role}/{endpoint}
	path := strings.TrimPrefix(r.URL.Path, "/api/agents/")
	parts := strings.Split(path, "/")

	if len(parts) < 2 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	roleStr := parts[0]
	endpoint := parts[1]

	role := agent.Role(roleStr)

	switch endpoint {
	case "tasks":
		s.handleAgentRoleTasks(w, r, role)
	case "status":
		s.handleAgentRoleStatus(w, r, role)
	case "outputs":
		s.handleAgentRoleOutputs(w, r, role)
	case "chat":
		s.handleAgentChat(w, r)
	default:
		http.Error(w, "Unknown endpoint", http.StatusNotFound)
	}
}

// handleAgentRoleTasks handles GET /api/agents/{role}/tasks
func (s *Server) handleAgentRoleTasks(w http.ResponseWriter, r *http.Request, role agent.Role) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tasks := s.agentTaskMgr.GetAgentTasks(role)

	writeJSON(w, map[string]interface{}{
		"role":  role,
		"tasks": tasks,
		"count": len(tasks),
	})
}

// handleAgentRoleStatus handles GET /api/agents/{role}/status
func (s *Server) handleAgentRoleStatus(w http.ResponseWriter, r *http.Request, role agent.Role) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	status, err := s.agentTaskMgr.GetAgentStatus(role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, status)
}

// handleAgentRoleOutputs handles GET /api/agents/{role}/outputs
func (s *Server) handleAgentRoleOutputs(w http.ResponseWriter, r *http.Request, role agent.Role) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get all tasks for the agent
	tasks := s.agentTaskMgr.GetAgentTasks(role)

	// Collect all outputs from all tasks
	var outputs []agentask.AgentOutput
	for _, task := range tasks {
		outputs = append(outputs, task.Outputs...)
	}

	writeJSON(w, map[string]interface{}{
		"role":    role,
		"outputs": outputs,
		"count":   len(outputs),
	})
}

// handleProjects handles GET /api/projects
func (s *Server) handleProjects(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read all directories in the projects folder
	entries, err := os.ReadDir(s.projectDir)
	if err != nil {
		http.Error(w, "Failed to read projects directory", http.StatusInternalServerError)
		return
	}

	type ProjectInfo struct {
		ID          string    `json:"id"`
		Name        string    `json:"name"`
		TaskCount   int       `json:"task_count"`
		LastUpdated time.Time `json:"last_updated,omitempty"`
	}

	projects := make([]ProjectInfo, 0)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		projectID := entry.Name()
		info := ProjectInfo{
			ID:   projectID,
			Name: projectID,
		}

		// Check if .tasks directory exists
		tasksDir := filepath.Join(s.projectDir, projectID, ".tasks")
		if tasksDirInfo, err := os.Stat(tasksDir); err == nil && tasksDirInfo.IsDir() {
			// Count tasks by reading history.json
			historyPath := filepath.Join(tasksDir, "history.json")
			if historyData, err := os.ReadFile(historyPath); err == nil {
				var history struct {
					Tasks []struct {
						CompletedAt *time.Time `json:"completedAt"`
					} `json:"tasks"`
				}
				if err := json.Unmarshal(historyData, &history); err == nil {
					info.TaskCount = len(history.Tasks)
					// Get last updated time from most recent task
					for _, task := range history.Tasks {
						if task.CompletedAt != nil && task.CompletedAt.After(info.LastUpdated) {
							info.LastUpdated = *task.CompletedAt
						}
					}
				}
			}

			// If no last updated from tasks, use directory mod time
			if info.LastUpdated.IsZero() {
				info.LastUpdated = tasksDirInfo.ModTime()
			}
		} else {
			// No .tasks directory yet - use project directory mod time
			if dirInfo, err := entry.Info(); err == nil {
				info.LastUpdated = dirInfo.ModTime()
			}
		}

		// Always add project, even if it has no tasks yet
		projects = append(projects, info)
	}

	// Sort by last updated (most recent first)
	sort.Slice(projects, func(i, j int) bool {
		return projects[i].LastUpdated.After(projects[j].LastUpdated)
	})

	writeJSON(w, map[string]interface{}{
		"projects": projects,
		"count":    len(projects),
	})
}

// handleFileContent serves the content of a file from the projects directory
func (s *Server) handleFileContent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get file path from query parameter
	filePath := r.URL.Query().Get("path")
	if filePath == "" {
		http.Error(w, "Missing path parameter", http.StatusBadRequest)
		return
	}

	// Security: ensure the path is within the projects directory
	absPath, err := filepath.Abs(filepath.Join(s.projectDir, filePath))
	if err != nil {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	absProjectDir, err := filepath.Abs(s.projectDir)
	if err != nil {
		http.Error(w, "Server configuration error", http.StatusInternalServerError)
		return
	}

	// Prevent directory traversal attacks
	if !strings.HasPrefix(absPath, absProjectDir) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Read the file
	content, err := os.ReadFile(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
		}
		return
	}

	// Determine content type based on file extension
	ext := filepath.Ext(filePath)
	contentType := "text/plain"
	switch ext {
	case ".md":
		contentType = "text/markdown"
	case ".json":
		contentType = "application/json"
	case ".html":
		contentType = "text/html"
	case ".css":
		contentType = "text/css"
	case ".js":
		contentType = "application/javascript"
	case ".go":
		contentType = "text/x-go"
	}

	w.Header().Set("Content-Type", contentType+"; charset=utf-8")
	w.Write(content)
}

// handleDevelopmentPlan returns the development plan for a project/task
func (s *Server) handleDevelopmentPlan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	projectID := r.URL.Query().Get("project")
	if projectID == "" {
		http.Error(w, "Missing project parameter", http.StatusBadRequest)
		return
	}

	// Build project directory path
	projectDir := filepath.Join(s.projectDir, projectID)

	// Load development plan
	plan, err := orchestrator.LoadDevelopmentPlan(projectDir)
	if err != nil {
		if os.IsNotExist(err) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"error": "No development plan found"}`))
			return
		}
		http.Error(w, "Failed to load development plan", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plan)
}

// handleSubTasks returns subtasks for a specific phase
func (s *Server) handleSubTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	projectID := r.URL.Query().Get("project")
	phaseIndexStr := r.URL.Query().Get("phase")

	if projectID == "" || phaseIndexStr == "" {
		http.Error(w, "Missing project or phase parameter", http.StatusBadRequest)
		return
	}

	phaseIndex, err := strconv.Atoi(phaseIndexStr)
	if err != nil {
		http.Error(w, "Invalid phase index", http.StatusBadRequest)
		return
	}

	projectDir := filepath.Join(s.projectDir, projectID)

	plan, err := orchestrator.LoadDevelopmentPlan(projectDir)
	if err != nil {
		http.Error(w, "Failed to load development plan", http.StatusInternalServerError)
		return
	}

	phase := plan.GetPhase(phaseIndex)
	if phase == nil {
		http.Error(w, "Phase not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(phase.SubTasks)
}

// handleQAReviews returns QA review history for a project
func (s *Server) handleQAReviews(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	projectID := r.URL.Query().Get("project")
	if projectID == "" {
		http.Error(w, "Missing project parameter", http.StatusBadRequest)
		return
	}

	projectDir := filepath.Join(s.projectDir, projectID)

	history, err := orchestrator.LoadQAReviewHistory(projectDir)
	if err != nil {
		http.Error(w, "Failed to load QA review history", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

// handlePhaseStatus returns current phase status for a project
func (s *Server) handlePhaseStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	projectID := r.URL.Query().Get("project")
	if projectID == "" {
		http.Error(w, "Missing project parameter", http.StatusBadRequest)
		return
	}

	projectDir := filepath.Join(s.projectDir, projectID)

	plan, err := orchestrator.LoadDevelopmentPlan(projectDir)
	if err != nil {
		// No development plan means we're in standard phases
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"has_plan":      false,
			"current_phase": s.orchestrator.GetCurrentPhase(),
		})
		return
	}

	currentPhase := plan.GetCurrentPhase()
	totalSubTasks := plan.GetTotalSubTasks()
	completedSubTasks := plan.GetCompletedSubTasks()

	status := map[string]interface{}{
		"has_plan":           true,
		"current_phase":      s.orchestrator.GetCurrentPhase(),
		"total_phases":       len(plan.Phases),
		"total_subtasks":     totalSubTasks,
		"completed_subtasks": completedSubTasks,
		"is_complete":        plan.IsComplete(),
	}

	if currentPhase != nil {
		status["current_dev_phase"] = map[string]interface{}{
			"index":       currentPhase.Index,
			"name":        currentPhase.Name,
			"status":      currentPhase.Status,
			"qa_status":   currentPhase.QAStatus,
			"iteration":   currentPhase.Iteration,
			"qa_feedback": currentPhase.QAFeedback,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

