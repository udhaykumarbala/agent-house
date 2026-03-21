package brain

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

// ExecutionResult holds the outcome of executing a Brain decision.
type ExecutionResult struct {
	Response   string `json:"response"`
	Action     string `json:"action"`
	ProjectID  string `json:"project_id,omitempty"`
	Success    bool   `json:"success"`
	NeedsReply bool   `json:"-"` // true if the Brain should be re-called with more context
	ExtraContext string `json:"-"` // additional context for re-call
}

// Executor handles Brain action execution.
type Executor struct {
	projectsDir string
	// Callbacks for actions that need orchestrator access
	OnCreateProject func(projectID, task string) error
	OnDelegate      func(projectID, agentRole, task string) error
	OnEscalate      func(projectID, task string) (string, error)
}

// NewExecutor creates an executor with the given project directory.
func NewExecutor(projectsDir string) *Executor {
	return &Executor{
		projectsDir: projectsDir,
	}
}

// Execute runs a Brain decision and returns the result.
func (e *Executor) Execute(decision *BrainDecision) *ExecutionResult {
	switch decision.Action {
	case ActionRespond:
		return &ExecutionResult{
			Response: decision.Response,
			Action:   ActionRespond,
			Success:  true,
		}

	case ActionCreateProject:
		return e.executeCreateProject(decision)

	case ActionDelegate:
		return e.executeDelegate(decision)

	case ActionListProjects:
		return e.executeListProjects(decision)

	case ActionProjectStatus:
		return e.executeProjectStatus(decision)

	case ActionEscalate:
		return e.executeEscalate(decision)

	case ActionSetReminder:
		return &ExecutionResult{
			Response: "Reminders are not yet implemented. I'll add this feature soon.",
			Action:   ActionRespond,
			Success:  true,
		}

	default:
		return &ExecutionResult{
			Response: decision.Response,
			Action:   decision.Action,
			Success:  true,
		}
	}
}

func (e *Executor) executeCreateProject(decision *BrainDecision) *ExecutionResult {
	projectID := decision.Params["project_id"]
	task := decision.Params["task"]

	if projectID == "" {
		return &ExecutionResult{
			Response: "I need a project name to create a project.",
			Action:   ActionRespond,
			Success:  false,
		}
	}

	if task == "" {
		task = decision.Response
	}

	if e.OnCreateProject != nil {
		if err := e.OnCreateProject(projectID, task); err != nil {
			return &ExecutionResult{
				Response:  fmt.Sprintf("Failed to create project: %v", err),
				Action:    ActionCreateProject,
				ProjectID: projectID,
				Success:   false,
			}
		}
	}

	return &ExecutionResult{
		Response:  decision.Response,
		Action:    ActionCreateProject,
		ProjectID: projectID,
		Success:   true,
	}
}

func (e *Executor) executeDelegate(decision *BrainDecision) *ExecutionResult {
	projectID := decision.Params["project_id"]
	agentRole := decision.Params["agent_role"]
	task := decision.Params["task"]

	if e.OnDelegate != nil {
		if err := e.OnDelegate(projectID, agentRole, task); err != nil {
			return &ExecutionResult{
				Response:  fmt.Sprintf("Failed to delegate: %v", err),
				Action:    ActionDelegate,
				ProjectID: projectID,
				Success:   false,
			}
		}
	}

	return &ExecutionResult{
		Response:  decision.Response,
		Action:    ActionDelegate,
		ProjectID: projectID,
		Success:   true,
	}
}

func (e *Executor) executeListProjects(decision *BrainDecision) *ExecutionResult {
	entries, err := os.ReadDir(e.projectsDir)
	if err != nil {
		return &ExecutionResult{
			Response: "No projects found.",
			Action:   ActionListProjects,
			Success:  true,
		}
	}

	var lines []string
	for _, entry := range entries {
		if !entry.IsDir() || entry.Name()[0] == '.' {
			continue
		}

		projDir := e.projectsDir + "/" + entry.Name()
		metaPath := projDir + "/project.json"
		status := "unknown"
		desc := ""

		if data, err := os.ReadFile(metaPath); err == nil {
			var meta struct {
				Status      string `json:"status"`
				Description string `json:"description"`
			}
			if json.Unmarshal(data, &meta) == nil {
				status = meta.Status
				desc = meta.Description
				if len(desc) > 60 {
					desc = desc[:60] + "..."
				}
			}
		}

		line := fmt.Sprintf("- **%s** [%s]", entry.Name(), status)
		if desc != "" {
			line += " — " + desc
		}
		lines = append(lines, line)
	}

	if len(lines) == 0 {
		return &ExecutionResult{
			Response: "No projects yet. Tell me what to build!",
			Action:   ActionListProjects,
			Success:  true,
		}
	}

	return &ExecutionResult{
		Response: fmt.Sprintf("Here are your projects:\n\n%s", strings.Join(lines, "\n")),
		Action:   ActionListProjects,
		Success:  true,
	}
}

func (e *Executor) executeProjectStatus(decision *BrainDecision) *ExecutionResult {
	projectID := decision.Params["project_id"]
	if projectID == "" {
		return &ExecutionResult{
			Response: "Which project? Tell me the project name.",
			Action:   ActionRespond,
			Success:  false,
		}
	}

	metaPath := e.projectsDir + "/" + projectID + "/project.json"
	data, err := os.ReadFile(metaPath)
	if err != nil {
		return &ExecutionResult{
			Response:  fmt.Sprintf("Project '%s' not found.", projectID),
			Action:    ActionProjectStatus,
			ProjectID: projectID,
			Success:   false,
		}
	}

	var meta struct {
		Status       string   `json:"status"`
		Description  string   `json:"description"`
		Direction    string   `json:"direction"`
		CurrentPhase string   `json:"current_phase"`
		Team         []string `json:"team"`
		TotalCost    float64  `json:"total_cost"`
		TotalTurns   int      `json:"total_turns"`
		FilesCreated int      `json:"files_created"`
		Pipeline     []string `json:"pipeline"`
	}
	json.Unmarshal(data, &meta)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("**%s** — %s\n\n", projectID, meta.Status))
	if meta.Description != "" {
		sb.WriteString(fmt.Sprintf("**Task:** %s\n", meta.Description))
	}
	if meta.CurrentPhase != "" {
		sb.WriteString(fmt.Sprintf("**Phase:** %s\n", meta.CurrentPhase))
	}
	if len(meta.Team) > 0 {
		sb.WriteString(fmt.Sprintf("**Team:** %s\n", strings.Join(meta.Team, ", ")))
	}
	if len(meta.Pipeline) > 0 {
		sb.WriteString(fmt.Sprintf("**Pipeline:** %s\n", strings.Join(meta.Pipeline, " → ")))
	}
	if meta.TotalTurns > 0 {
		sb.WriteString(fmt.Sprintf("**Stats:** %d turns, %d files, $%.4f\n", meta.TotalTurns, meta.FilesCreated, meta.TotalCost))
	}

	return &ExecutionResult{
		Response:  sb.String(),
		Action:    ActionProjectStatus,
		ProjectID: projectID,
		Success:   true,
	}
}

func (e *Executor) executeEscalate(decision *BrainDecision) *ExecutionResult {
	projectID := decision.Params["project_id"]
	task := decision.Params["task"]

	if e.OnEscalate != nil {
		result, err := e.OnEscalate(projectID, task)
		if err != nil {
			log.Printf("[BRAIN] Escalation failed: %v", err)
			return &ExecutionResult{
				Response: fmt.Sprintf("Deep analysis failed: %v", err),
				Action:   ActionEscalate,
				Success:  false,
			}
		}
		return &ExecutionResult{
			Response:  result,
			Action:    ActionEscalate,
			ProjectID: projectID,
			Success:   true,
		}
	}

	return &ExecutionResult{
		Response:  decision.Response,
		Action:    ActionEscalate,
		ProjectID: projectID,
		Success:   true,
	}
}
