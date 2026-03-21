package web

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"pty-claude-test/internal/kanban"
)

// handleKanbanTasks handles GET /api/kanban/tasks and POST /api/kanban/tasks
func (s *Server) handleKanbanTasks(w http.ResponseWriter, r *http.Request) {
	projectID := r.URL.Query().Get("project")
	if projectID == "" {
		projectID = "default"
	}
	projectDir := filepath.Join(s.projectDir, projectID)

	switch r.Method {
	case http.MethodGet:
		// Filter params
		status := kanban.Status(r.URL.Query().Get("status"))
		agent := r.URL.Query().Get("agent")
		phaseStr := r.URL.Query().Get("phase")
		priority := kanban.Priority(r.URL.Query().Get("priority"))

		phaseIdx := 0
		if phaseStr != "" {
			phaseIdx, _ = strconv.Atoi(phaseStr)
		}

		tasks, err := kanban.FilterTasks(projectDir, status, agent, phaseIdx, priority)
		if err != nil {
			tasks = []kanban.Task{}
		}

		writeJSON(w, map[string]interface{}{
			"tasks":   tasks,
			"count":   len(tasks),
			"columns": kanban.GetColumnCounts(tasks),
		})

	case http.MethodPost:
		var t kanban.Task
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if t.Title == "" {
			http.Error(w, "Title is required", http.StatusBadRequest)
			return
		}
		t.ProjectID = projectID

		created, err := kanban.CreateTask(projectDir, t)
		if err != nil {
			http.Error(w, "Failed to create task: "+err.Error(), http.StatusInternalServerError)
			return
		}

		writeJSON(w, map[string]interface{}{
			"success": true,
			"task":    created,
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleKanbanTaskByID handles:
//
//	GET    /api/kanban/tasks/{id}
//	PATCH  /api/kanban/tasks/{id}
//	POST   /api/kanban/tasks/{id}/review
func (s *Server) handleKanbanTaskByID(w http.ResponseWriter, r *http.Request) {
	projectID := r.URL.Query().Get("project")
	if projectID == "" {
		projectID = "default"
	}
	projectDir := filepath.Join(s.projectDir, projectID)

	path := strings.TrimPrefix(r.URL.Path, "/api/kanban/tasks/")
	parts := strings.Split(path, "/")
	taskID := parts[0]

	if taskID == "" {
		http.Error(w, "Task ID required", http.StatusBadRequest)
		return
	}

	// POST /api/kanban/tasks/{id}/review
	if len(parts) > 1 && parts[1] == "review" && r.Method == http.MethodPost {
		var req struct {
			Action          string                   `json:"action"`
			Feedback        string                   `json:"feedback"`
			CriteriaUpdates []map[string]interface{} `json:"criteria_updates"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		task, err := kanban.ReviewTask(projectDir, taskID, req.Action, req.Feedback, req.CriteriaUpdates)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		writeJSON(w, map[string]interface{}{"success": true, "task": task})
		return
	}

	// GET /api/kanban/tasks/{id}
	if r.Method == http.MethodGet {
		task, err := kanban.GetTask(projectDir, taskID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		writeJSON(w, task)
		return
	}

	// PATCH /api/kanban/tasks/{id}
	if r.Method == http.MethodPatch {
		var updates map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		task, err := kanban.UpdateTask(projectDir, taskID, updates)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		writeJSON(w, map[string]interface{}{"success": true, "task": task})
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// handleKanbanSync syncs kanban tasks from the development plan
// POST /api/kanban/sync?project=xxx
func (s *Server) handleKanbanSync(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	projectID := r.URL.Query().Get("project")
	if projectID == "" {
		projectID = "default"
	}
	projectDir := filepath.Join(s.projectDir, projectID)

	// Load development plan
	plan, err := loadDevPlanForKanban(projectDir)
	if err != nil {
		writeJSON(w, map[string]interface{}{
			"success": false,
			"error":   "No development plan found",
		})
		return
	}

	created, err := kanban.SyncFromDevPlan(projectDir, projectID, plan)
	if err != nil {
		writeJSON(w, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	writeJSON(w, map[string]interface{}{
		"success": true,
		"created": created,
	})
}

// loadDevPlanForKanban reads the development plan and converts to kanban's format
func loadDevPlanForKanban(projectDir string) ([]kanban.DevPlanPhase, error) {
	plan, err := loadRawDevPlan(projectDir)
	if err != nil {
		return nil, err
	}

	var phases []kanban.DevPlanPhase
	for _, p := range plan.Phases {
		phase := kanban.DevPlanPhase{
			Index: p.Index,
			Name:  p.Name,
		}
		for _, st := range p.SubTasks {
			sub := kanban.DevPlanSubTask{
				ID:                 st.ID,
				Title:              st.Title,
				Description:        st.Description,
				AssignedAgents:     st.AssignedAgents,
				CompletionCriteria: st.CompletionCriteria,
			}
			phase.SubTasks = append(phase.SubTasks, sub)
		}
		phases = append(phases, phase)
	}
	return phases, nil
}

// rawDevPlan mirrors the orchestrator's DevelopmentPlan for JSON loading
type rawDevPlan struct {
	Phases []rawDevPhase `json:"phases"`
}

type rawDevPhase struct {
	Index    int             `json:"index"`
	Name     string          `json:"name"`
	SubTasks []rawDevSubTask `json:"sub_tasks"`
}

type rawDevSubTask struct {
	ID                 string   `json:"id"`
	Title              string   `json:"title"`
	Description        string   `json:"description"`
	AssignedAgents     []string `json:"assigned_agents"`
	CompletionCriteria []string `json:"completion_criteria"`
}

func loadRawDevPlan(projectDir string) (*rawDevPlan, error) {
	data, err := os.ReadFile(filepath.Join(projectDir, ".plans", "development-plan.json"))
	if err != nil {
		return nil, err
	}
	var plan rawDevPlan
	if err := json.Unmarshal(data, &plan); err != nil {
		return nil, err
	}
	return &plan, nil
}
