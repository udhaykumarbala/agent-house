package web

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"pty-claude-test/internal/agent"
	"pty-claude-test/internal/message"
	"pty-claude-test/internal/orchestrator"
)

// Server handles HTTP requests for the dashboard
type Server struct {
	store        *message.Store
	orchestrator *orchestrator.Orchestrator
	projectDir   string
	mu           sync.RWMutex
	taskRunning  bool
	hub          *Hub
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
		ProjectDir:  config.ProjectDir,
		MaxDepth:    5,
		MaxTurns:    20,
		EnableFiles: true,
		Verbose:     true,
	}

	hub := NewHub()

	server := &Server{
		store:        config.Store,
		orchestrator: orchestrator.New(orchConfig, config.Store),
		projectDir:   config.ProjectDir,
		hub:          hub,
	}

	// Register orchestrator message callback for WebSocket broadcast
	server.orchestrator.OnMessage(func(msg *message.Message) {
		hub.Broadcast(msg)
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
	mux.HandleFunc("/api/files", s.handleFiles)
	mux.HandleFunc("/api/status", s.handleStatus)

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
	Task      string `json:"task"`
	ProjectID string `json:"project_id"`
}

// handleTask handles task submission
func (s *Server) handleTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if task is already running
	s.mu.RLock()
	running := s.taskRunning
	s.mu.RUnlock()

	if running {
		writeJSON(w, map[string]interface{}{
			"error":   "Task already running",
			"success": false,
		})
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

	// Start task in background
	s.mu.Lock()
	s.taskRunning = true
	s.mu.Unlock()

	go func() {
		defer func() {
			s.mu.Lock()
			s.taskRunning = false
			s.mu.Unlock()
		}()

		projectDir := filepath.Join(s.projectDir, req.ProjectID)
		os.MkdirAll(projectDir, 0755)

		// Update orchestrator project dir
		config := orchestrator.Config{
			ProjectDir:  projectDir,
			MaxDepth:    5,
			MaxTurns:    20,
			EnableFiles: true,
			Verbose:     true,
		}
		s.orchestrator = orchestrator.New(config, s.store)

		// Re-register WebSocket callback on new orchestrator
		s.orchestrator.OnMessage(func(msg *message.Message) {
			s.hub.Broadcast(msg)
		})

		s.orchestrator.ProcessTask(req.ProjectID, req.Task)
	}()

	writeJSON(w, map[string]interface{}{
		"success":    true,
		"message":    "Task started",
		"project_id": req.ProjectID,
	})
}

// FileInfo represents a file in the project
type FileInfo struct {
	Path     string `json:"path"`
	Name     string `json:"name"`
	IsDir    bool   `json:"is_dir"`
	Size     int64  `json:"size"`
	Modified string `json:"modified"`
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
	running := s.taskRunning
	s.mu.RUnlock()

	wsClients := 0
	if s.hub != nil {
		wsClients = s.hub.ClientCount()
	}

	writeJSON(w, map[string]interface{}{
		"status":        "ok",
		"task_running":  running,
		"message_count": s.store.Count(),
		"project_dir":   s.projectDir,
		"ws_clients":    wsClients,
	})
}

// handleStatic serves static files
func (s *Server) handleStatic(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" || r.URL.Path == "/index.html" {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(dashboardHTML))
		return
	}

	http.NotFound(w, r)
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

