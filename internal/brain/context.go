package brain

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// WorkspaceContext holds assembled workspace state for the Brain.
type WorkspaceContext struct {
	Projects       []ProjectSummary `json:"projects"`
	ActiveAgents   []string         `json:"active_agents,omitempty"`
	RunningTasks   []string         `json:"running_tasks,omitempty"`
	RecentActivity []string         `json:"recent_activity,omitempty"`
}

// ProjectSummary is a compact view of a project for the Brain's context.
type ProjectSummary struct {
	ID           string   `json:"id"`
	Status       string   `json:"status"`
	Description  string   `json:"description,omitempty"`
	CurrentPhase string   `json:"current_phase,omitempty"`
	Team         []string `json:"team,omitempty"`
	FileCount    int      `json:"file_count"`
	LastEvent    string   `json:"last_event,omitempty"`
}

// AssembleContext reads the workspace state and returns a formatted context string.
// This is pure Go — no LLM calls. Takes <100ms.
func AssembleContext(projectsDir string) string {
	ctx := WorkspaceContext{}

	// Scan projects
	entries, err := os.ReadDir(projectsDir)
	if err != nil {
		return "No projects found."
	}

	for _, entry := range entries {
		if !entry.IsDir() || entry.Name()[0] == '.' {
			continue
		}

		projDir := filepath.Join(projectsDir, entry.Name())
		summary := ProjectSummary{ID: entry.Name()}

		// Try to load project.json
		metaPath := filepath.Join(projDir, "project.json")
		if data, err := os.ReadFile(metaPath); err == nil {
			var meta struct {
				Status       string   `json:"status"`
				Description  string   `json:"description"`
				CurrentPhase string   `json:"current_phase"`
				Team         []string `json:"team"`
				TotalCost    float64  `json:"total_cost"`
				FilesCreated int      `json:"files_created"`
				History      []struct {
					Time  time.Time `json:"time"`
					Event string    `json:"event"`
				} `json:"history"`
			}
			if json.Unmarshal(data, &meta) == nil {
				summary.Status = meta.Status
				summary.Description = truncate(meta.Description, 100)
				summary.CurrentPhase = meta.CurrentPhase
				summary.Team = meta.Team
				summary.FileCount = meta.FilesCreated
				if len(meta.History) > 0 {
					last := meta.History[len(meta.History)-1]
					summary.LastEvent = last.Event
				}
			}
		} else {
			// No project.json — count files
			count := 0
			filepath.Walk(projDir, func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return nil
				}
				base := filepath.Base(filepath.Dir(path))
				if base == ".tasks" || base == ".plans" || base == "node_modules" {
					return filepath.SkipDir
				}
				count++
				return nil
			})
			summary.FileCount = count
			summary.Status = "unknown"
		}

		ctx.Projects = append(ctx.Projects, summary)
	}

	return formatContext(ctx)
}

func formatContext(ctx WorkspaceContext) string {
	var sb strings.Builder

	sb.WriteString("## Workspace State\n\n")

	if len(ctx.Projects) == 0 {
		sb.WriteString("No projects yet.\n")
	} else {
		sb.WriteString(fmt.Sprintf("**%d projects:**\n", len(ctx.Projects)))
		for _, p := range ctx.Projects {
			status := p.Status
			if status == "" {
				status = "unknown"
			}
			sb.WriteString(fmt.Sprintf("- **%s** [%s]", p.ID, status))
			if p.Description != "" {
				sb.WriteString(fmt.Sprintf(" — %s", p.Description))
			}
			if p.CurrentPhase != "" {
				sb.WriteString(fmt.Sprintf(" (phase: %s)", p.CurrentPhase))
			}
			sb.WriteString(fmt.Sprintf(" [%d files]", p.FileCount))
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
