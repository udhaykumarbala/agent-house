package web

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"pty-claude-test/internal/checkpoint"
)

// handleAllCheckpoints aggregates checkpoints across every project dir, each
// annotated with its project's run mode + decision timeout (so the UI can show
// a live countdown for full-auto gates), sorted newest-first.
func (s *Server) handleAllCheckpoints(w http.ResponseWriter, r *http.Request) {
	type cpOut struct {
		checkpoint.Checkpoint
		Project                string `json:"project"`
		RunMode                string `json:"run_mode"`
		DecisionTimeoutMinutes int    `json:"decision_timeout_minutes"`
	}
	out := []cpOut{}
	pending := 0
	if entries, err := os.ReadDir(s.projectDir); err == nil {
		for _, e := range entries {
			if !e.IsDir() || strings.HasPrefix(e.Name(), ".") || strings.HasPrefix(e.Name(), "_") {
				continue
			}
			projDir := filepath.Join(s.projectDir, e.Name())
			cps, err := checkpoint.LoadCheckpoints(projDir)
			if err != nil || len(cps) == 0 {
				continue
			}
			rm, to := "", 0
			if st, err := checkpoint.LoadSettings(projDir); err == nil {
				rm, to = st.RunMode, st.DecisionTimeoutMinutes
			}
			for _, cp := range cps {
				out = append(out, cpOut{Checkpoint: cp, Project: e.Name(), RunMode: rm, DecisionTimeoutMinutes: to})
				if cp.Status == "pending" {
					pending++
				}
			}
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.After(out[j].CreatedAt) })
	writeJSON(w, map[string]interface{}{"checkpoints": out, "count": len(out), "pending_count": pending})
}

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

	// GET /api/checkpoints/all — aggregate checkpoints across every project,
	// each annotated with its project's run mode + decision timeout so the UI
	// can render a countdown. Lets Mission follow whichever build is running
	// (each build runs on its own orchestrator + project dir).
	if cpID == "all" && r.Method == http.MethodGet {
		s.handleAllCheckpoints(w, r)
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
		// When a known run mode is named, (re)apply its preset so the gate flags
		// + decision timeout match the mode (the request's decision_timeout_minutes
		// is honored for full_auto). Any other value ("custom"/"") persists the
		// raw Require* flags as sent, for advanced per-gate tweaking.
		switch settings.RunMode {
		case checkpoint.RunModeManual, checkpoint.RunModeSemiAuto,
			checkpoint.RunModeFullAuto, checkpoint.RunModeBlitz:
			settings.ApplyRunModePreset(settings.RunMode)
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
