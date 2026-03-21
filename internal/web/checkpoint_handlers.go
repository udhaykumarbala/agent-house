package web

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"pty-claude-test/internal/checkpoint"
)

// handleCheckpoints handles GET /api/checkpoints?project=xxx
func (s *Server) handleCheckpoints(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	projectID := r.URL.Query().Get("project")
	if projectID == "" {
		projectID = "default"
	}

	projectDir := filepath.Join(s.projectDir, projectID)
	cps, err := checkpoint.LoadCheckpoints(projectDir)
	if err != nil {
		writeJSON(w, map[string]interface{}{"error": err.Error(), "checkpoints": []checkpoint.Checkpoint{}})
		return
	}

	pendingCount := 0
	for _, cp := range cps {
		if cp.Status == "pending" {
			pendingCount++
		}
	}

	writeJSON(w, map[string]interface{}{
		"checkpoints":   cps,
		"count":         len(cps),
		"pending_count": pendingCount,
	})
}

// handleCheckpointByID handles:
//   GET  /api/checkpoints/{id}         — returns checkpoint + artifact content
//   POST /api/checkpoints/{id}/decide  — resolves checkpoint
func (s *Server) handleCheckpointByID(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/checkpoints/")
	parts := strings.Split(path, "/")
	cpID := parts[0]

	if cpID == "" {
		http.Error(w, "Checkpoint ID required", http.StatusBadRequest)
		return
	}

	projectID := r.URL.Query().Get("project")
	if projectID == "" {
		projectID = "default"
	}

	projectDir := filepath.Join(s.projectDir, projectID)

	// POST /api/checkpoints/{id}/decide
	if len(parts) > 1 && parts[1] == "decide" && r.Method == http.MethodPost {
		s.handleCheckpointDecide(w, r, cpID, projectID, projectDir)
		return
	}

	// GET /api/checkpoints/{id}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cp, err := checkpoint.GetCheckpoint(projectDir, cpID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Read artifact content
	artifactContent := ""
	if cp.ArtifactPath != "" {
		if data, err := os.ReadFile(cp.ArtifactPath); err == nil {
			artifactContent = string(data)
		}
	}

	writeJSON(w, map[string]interface{}{
		"checkpoint":       cp,
		"artifact_content": artifactContent,
	})
}

// handleCheckpointDecide resolves a checkpoint with a decision
func (s *Server) handleCheckpointDecide(w http.ResponseWriter, r *http.Request, cpID, projectID, projectDir string) {
	var req struct {
		Action       string `json:"action"`
		Feedback     string `json:"feedback"`
		OverrideData string `json:"override_data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Action == "" {
		http.Error(w, "Action is required (approved, rejected, overridden)", http.StatusBadRequest)
		return
	}

	dec := &checkpoint.Decision{
		ID:             checkpoint.GenerateDecisionID(),
		ProjectID:      projectID,
		TaskID:         "",
		CheckpointType: "",
		Action:         checkpoint.DecisionAction(req.Action),
		Feedback:       req.Feedback,
		OverrideData:   req.OverrideData,
		DecidedBy:      "user",
		DecidedAt:      time.Now(),
	}

	// Fill in task/type from the pending checkpoint
	pending := s.orchestrator.GetPendingCheckpoint()
	if pending != nil {
		dec.TaskID = pending.TaskID
		dec.CheckpointType = pending.Type
		dec.PhaseIndex = pending.PhaseIndex
	}

	if err := s.orchestrator.ResolveCheckpoint(dec); err != nil {
		writeJSON(w, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	writeJSON(w, map[string]interface{}{
		"success":           true,
		"decision_id":       dec.ID,
		"checkpoint_status": "resolved",
	})
}

// handleWorkflowSettings handles GET/PUT /api/settings/workflow
func (s *Server) handleWorkflowSettings(w http.ResponseWriter, r *http.Request) {
	projectID := r.URL.Query().Get("project")
	if projectID == "" {
		projectID = "default"
	}

	projectDir := filepath.Join(s.projectDir, projectID)

	switch r.Method {
	case http.MethodGet:
		settings, err := checkpoint.LoadSettings(projectDir)
		if err != nil {
			writeJSON(w, checkpoint.DefaultSettings(projectID))
			return
		}
		writeJSON(w, settings)

	case http.MethodPut:
		var settings checkpoint.WorkflowSettings
		if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if settings.ProjectID == "" {
			settings.ProjectID = projectID
		}
		if err := checkpoint.SaveSettings(projectDir, &settings); err != nil {
			http.Error(w, "Failed to save settings", http.StatusInternalServerError)
			return
		}
		writeJSON(w, map[string]interface{}{
			"success":  true,
			"settings": settings,
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
