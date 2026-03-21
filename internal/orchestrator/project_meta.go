package orchestrator

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// ProjectMeta holds project metadata, auto-created and maintained by the orchestrator.
type ProjectMeta struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Status       string          `json:"status"`        // active, completed, archived
	Type         string          `json:"type"`           // from CEO triage: new_project, improvement, etc.
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	CreatedBy    string          `json:"created_by"`     // "human" or "brain"
	Description  string          `json:"description"`    // original task description
	Direction    string          `json:"direction"`      // CEO's strategic direction
	Team         []string        `json:"team"`           // agent roles involved
	Pipeline     []string        `json:"pipeline"`       // phases that ran
	CurrentPhase string          `json:"current_phase"`
	TotalCost    float64         `json:"total_cost"`
	TotalTurns   int             `json:"total_turns"`
	FilesCreated int             `json:"files_created"`
	History      []ProjectEvent  `json:"history"`
}

// ProjectEvent is a timestamped entry in the project history.
type ProjectEvent struct {
	Time    time.Time `json:"time"`
	Event   string    `json:"event"`
	Agent   string    `json:"agent,omitempty"`
	Details string    `json:"details,omitempty"`
}

// LoadProjectMeta loads project.json from the project directory.
func LoadProjectMeta(projectDir string) (*ProjectMeta, error) {
	path := filepath.Join(projectDir, "project.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var meta ProjectMeta
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}
	return &meta, nil
}

// SaveProjectMeta writes project.json to the project directory.
func SaveProjectMeta(projectDir string, meta *ProjectMeta) error {
	meta.UpdatedAt = time.Now()
	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}
	path := filepath.Join(projectDir, "project.json")
	return os.WriteFile(path, data, 0644)
}

// InitProjectMeta creates an initial project.json when a task is first submitted.
func InitProjectMeta(projectDir, projectID, taskDescription string) *ProjectMeta {
	now := time.Now()
	meta := &ProjectMeta{
		ID:          projectID,
		Name:        projectID,
		Status:      "active",
		CreatedAt:   now,
		UpdatedAt:   now,
		CreatedBy:   "human",
		Description: taskDescription,
		History: []ProjectEvent{
			{Time: now, Event: "project_created", Details: taskDescription},
		},
	}

	if err := os.MkdirAll(projectDir, 0755); err != nil {
		log.Printf("[PROJECT] Failed to create directory %s: %v", projectDir, err)
	}

	if err := SaveProjectMeta(projectDir, meta); err != nil {
		log.Printf("[PROJECT] Failed to save project.json: %v", err)
	} else {
		log.Printf("[PROJECT] Created project.json for %s", projectID)
	}

	return meta
}

// AddProjectEvent appends an event to the project history and saves.
func AddProjectEvent(projectDir, event, agentRole, details string) {
	meta, err := LoadProjectMeta(projectDir)
	if err != nil {
		return
	}
	meta.History = append(meta.History, ProjectEvent{
		Time:    time.Now(),
		Event:   event,
		Agent:   agentRole,
		Details: details,
	})
	// Keep history bounded
	if len(meta.History) > 200 {
		meta.History = meta.History[len(meta.History)-100:]
	}
	SaveProjectMeta(projectDir, meta)
}

// UpdateProjectPhase updates the current phase and saves.
func UpdateProjectPhase(projectDir, phase string) {
	meta, err := LoadProjectMeta(projectDir)
	if err != nil {
		return
	}
	meta.CurrentPhase = phase
	if !containsString(meta.Pipeline, phase) {
		meta.Pipeline = append(meta.Pipeline, phase)
	}
	SaveProjectMeta(projectDir, meta)
}

// UpdateProjectCompletion updates final stats and marks project as completed.
func UpdateProjectCompletion(projectDir string, totalCost float64, totalTurns, filesCreated int, team []string) {
	meta, err := LoadProjectMeta(projectDir)
	if err != nil {
		return
	}
	meta.Status = "completed"
	meta.TotalCost = totalCost
	meta.TotalTurns = totalTurns
	meta.FilesCreated = filesCreated
	meta.Team = team
	meta.History = append(meta.History, ProjectEvent{
		Time:    time.Now(),
		Event:   "project_completed",
		Details: fmt.Sprintf("%d turns, %d files, $%.4f", totalTurns, filesCreated, totalCost),
	})
	SaveProjectMeta(projectDir, meta)
}

func containsString(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
