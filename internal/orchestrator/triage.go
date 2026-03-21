package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"pty-claude-test/internal/agent"
	"pty-claude-test/internal/message"
)

// TriageResult holds the CEO's decision on how to handle a task
type TriageResult struct {
	TaskType       string   // "new_project", "improvement", "bug_fix", "feature_add"
	SkipPhases     []Phase  // Phases to skip
	RequiredAgents []string // Which agents are needed (empty = all)
	Direction      string   // CEO's strategic direction/summary
	Reasoning      string   // Why these decisions were made
}

// executeTriage runs the CEO agent to analyze the task and project context,
// then returns a TriageResult that determines which phases to run.
func (o *Orchestrator) executeTriage(taskStr, projectID string, result *Result) (*TriageResult, error) {
	o.currentPhase = PhaseTriage

	// Notify about triage phase
	phaseMsg := message.NewMessage(message.TypeSystem, "orchestrator", "all", "CEO is triaging the task...")
	phaseMsg.Metadata.ProjectID = projectID
	phaseMsg.Metadata.TaskID = result.TaskID
	o.store.Add(phaseMsg)
	o.notify(phaseMsg)

	// Build project context for CEO
	triageContext := buildTriageContext(o.config.ProjectDir, taskStr, projectID)

	// Run CEO agent
	ceoAgent, err := o.getAgent(agent.RoleCEO)
	if err != nil {
		return nil, fmt.Errorf("failed to get CEO agent for triage: %w", err)
	}

	log.Printf("[TRIAGE] CEO analyzing task for project=%s", projectID)

	var responseContent string

	if ceoAgent.SessionManager != nil {
		// Session-based triage
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		sessResp, sessErr := ceoAgent.ExecuteWithSession(ctx, projectID, o.config.ProjectDir, triageContext)
		cancel()
		if sessErr != nil {
			log.Printf("[TRIAGE] CEO triage (session) failed: %v, falling back to full pipeline", sessErr)
			return defaultTriage(), nil
		}
		responseContent = sessResp.Text
	} else {
		// Legacy triage
		response, _, legacyErr := ceoAgent.ProcessWithFileOps(triageContext, o.config.ProjectDir)
		if legacyErr != nil {
			log.Printf("[TRIAGE] CEO triage failed: %v, falling back to full pipeline", legacyErr)
			return defaultTriage(), nil
		}
		responseContent = response.Content
	}

	// Record the CEO's response
	responseMsg := message.NewMessage(message.TypeResponse, string(agent.RoleCEO), "orchestrator", responseContent)
	responseMsg.Metadata.ProjectID = projectID
	responseMsg.Metadata.TaskID = result.TaskID
	o.store.Add(responseMsg)
	o.notify(responseMsg)
	result.Messages = append(result.Messages, responseMsg)

	// Parse the triage decision
	triage := parseTriageResponse(responseContent)

	log.Printf("[TRIAGE] Decision: type=%s skip=%v agents=%v", triage.TaskType, triage.SkipPhases, triage.RequiredAgents)

	return triage, nil
}

// buildTriageContext assembles project information for the CEO
func buildTriageContext(projectDir, taskStr, projectID string) string {
	var b strings.Builder

	b.WriteString("You are the CEO of an AI development team. A new task has been submitted.\n\n")
	b.WriteString("YOUR JOB: Analyze the task and the current project state, then decide HOW to execute it.\n\n")

	b.WriteString("## Task\n")
	b.WriteString(taskStr)
	b.WriteString("\n\n## Project: " + projectID + "\n\n")

	// Check if project has existing files
	files := listProjectFiles(projectDir)
	if len(files) > 0 {
		b.WriteString("### Existing Project Files\n")
		for _, f := range files {
			b.WriteString("- " + f + "\n")
		}
		b.WriteString("\n")
	} else {
		b.WriteString("### Project Status: EMPTY (no files yet)\n\n")
	}

	// Check if template was already selected
	templatePath := filepath.Join(projectDir, ".plans", "template.md")
	if data, err := os.ReadFile(templatePath); err == nil {
		b.WriteString("### Existing Template\n")
		content := string(data)
		if len(content) > 500 {
			content = content[:500] + "..."
		}
		b.WriteString(content + "\n\n")
	}

	// Check previous task history
	historyPath := filepath.Join(projectDir, ".tasks", "history.json")
	if data, err := os.ReadFile(historyPath); err == nil {
		var history struct {
			Tasks []struct {
				Summary string `json:"summary"`
				Status  string `json:"status"`
			} `json:"tasks"`
		}
		if json.Unmarshal(data, &history) == nil && len(history.Tasks) > 0 {
			b.WriteString("### Previous Tasks in This Project\n")
			for _, t := range history.Tasks {
				b.WriteString(fmt.Sprintf("- [%s] %s\n", t.Status, t.Summary))
			}
			b.WriteString("\n")
		}
	}

	// Check if development plan exists
	planPath := filepath.Join(projectDir, ".plans", "development-plan.json")
	if _, err := os.Stat(planPath); err == nil {
		b.WriteString("### Existing Development Plan: YES (from previous build)\n\n")
	}

	// Check if approved plan exists
	approvedPath := filepath.Join(projectDir, ".plans", "final", "approved-plan.md")
	if data, err := os.ReadFile(approvedPath); err == nil {
		content := string(data)
		if len(content) > 800 {
			content = content[:800] + "..."
		}
		b.WriteString("### Existing Approved Plan\n" + content + "\n\n")
	}

	b.WriteString(`## Your Decision

Based on the task and project state, respond with your triage decision using this EXACT format:

TRIAGE_DECISION:
TASK_TYPE: [new_project | improvement | bug_fix | feature_add | refactor]
SKIP_PHASES: [comma-separated list of phases to skip, or "none"]
REQUIRED_AGENTS: [comma-separated agent roles needed, or "all"]
DIRECTION: [1-3 sentence strategic direction for the team]
REASONING: [Why you chose this approach]
END_TRIAGE

Valid phase names to skip: template_selection, research, planning, discussion
Valid agent roles: ceo, pm, ux, ui, security, architect, senior-dev, junior-dev

### Guidelines:
- **New empty project**: Don't skip anything. Full pipeline.
- **Improvement to existing project**: Skip template_selection (already chosen). Consider skipping or shortening research if the domain hasn't changed.
- **Bug fix**: Skip template_selection, research, and planning. Go straight to discussion → development.
- **Feature addition**: Skip template_selection. Research may be needed if the feature is in a new domain.
- **Refactor**: Skip template_selection and research. Planning → development.

Be decisive. Your team is waiting.
`)

	return b.String()
}

// listProjectFiles returns a list of files in the project dir (excluding .plans/.tasks)
func listProjectFiles(projectDir string) []string {
	var files []string
	filepath.WalkDir(projectDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		name := d.Name()
		if strings.HasPrefix(name, ".") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if d.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(projectDir, path)
		files = append(files, rel)
		if len(files) > 30 {
			return filepath.SkipAll
		}
		return nil
	})
	return files
}

// parseTriageResponse extracts the triage decision from CEO's response
func parseTriageResponse(content string) *TriageResult {
	triage := &TriageResult{
		TaskType: "new_project",
	}

	// Find TRIAGE_DECISION block
	start := strings.Index(content, "TRIAGE_DECISION:")
	end := strings.Index(content, "END_TRIAGE")
	if start < 0 {
		log.Printf("[TRIAGE] No TRIAGE_DECISION marker found, using defaults")
		return defaultTriage()
	}

	block := content[start:]
	if end > start {
		block = content[start:end]
	}

	lines := strings.Split(block, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "TASK_TYPE:") {
			triage.TaskType = strings.TrimSpace(strings.TrimPrefix(line, "TASK_TYPE:"))
		}
		if strings.HasPrefix(line, "SKIP_PHASES:") {
			val := strings.TrimSpace(strings.TrimPrefix(line, "SKIP_PHASES:"))
			if val != "none" && val != "" {
				for _, p := range strings.Split(val, ",") {
					p = strings.TrimSpace(p)
					if p != "" {
						triage.SkipPhases = append(triage.SkipPhases, Phase(p))
					}
				}
			}
		}
		if strings.HasPrefix(line, "REQUIRED_AGENTS:") {
			val := strings.TrimSpace(strings.TrimPrefix(line, "REQUIRED_AGENTS:"))
			if val != "all" && val != "" {
				for _, a := range strings.Split(val, ",") {
					a = strings.TrimSpace(a)
					if a != "" {
						triage.RequiredAgents = append(triage.RequiredAgents, a)
					}
				}
			}
		}
		if strings.HasPrefix(line, "DIRECTION:") {
			triage.Direction = strings.TrimSpace(strings.TrimPrefix(line, "DIRECTION:"))
		}
		if strings.HasPrefix(line, "REASONING:") {
			triage.Reasoning = strings.TrimSpace(strings.TrimPrefix(line, "REASONING:"))
		}
	}

	return triage
}

// defaultTriage returns the default full-pipeline triage for new projects
func defaultTriage() *TriageResult {
	return &TriageResult{
		TaskType:  "new_project",
		Direction: "Full pipeline — new project build",
	}
}

// shouldSkipPhase checks if a phase should be skipped based on triage
func (t *TriageResult) shouldSkipPhase(phase Phase) bool {
	for _, skip := range t.SkipPhases {
		if skip == phase {
			return true
		}
	}
	return false
}
