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
	Vendors        []VendorSummary  `json:"vendors,omitempty"`
	RecentEmails   []EmailSummary   `json:"recent_emails,omitempty"`
}

// VendorSummary is a compact vendor view with contract info.
type VendorSummary struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Domain        string `json:"domain"`
	Contact       string `json:"contact"`
	Active        bool   `json:"active"`
	ContractTerms string `json:"contract_terms,omitempty"` // serialized terms
}

// EmailSummary is an email for Brain context.
type EmailSummary struct {
	ID       string `json:"id"`
	From     string `json:"from"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
	Trust    string `json:"trust"`
	TrustReason string `json:"trust_reason"`
	Category string `json:"category"`
	Read     bool   `json:"read"`
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

	// Load vendors
	vendorsPath := filepath.Join(projectsDir, "vendors.json")
	if data, err := os.ReadFile(vendorsPath); err == nil {
		var vendors []struct {
			ID            string          `json:"id"`
			Name          string          `json:"name"`
			Domain        string          `json:"domain"`
			ContactPerson string          `json:"contact_person"`
			ContractActive bool           `json:"contract_active"`
			ContractTerms json.RawMessage `json:"contract_terms"`
		}
		if json.Unmarshal(data, &vendors) == nil {
			for _, v := range vendors {
				terms := ""
				if v.ContractTerms != nil {
					terms = string(v.ContractTerms)
				}
				ctx.Vendors = append(ctx.Vendors, VendorSummary{
					ID: v.ID, Name: v.Name, Domain: v.Domain,
					Contact: v.ContactPerson, Active: v.ContractActive,
					ContractTerms: terms,
				})
			}
		}
	}

	// Load recent emails from inbox (full content for Brain to reference)
	inboxDir := filepath.Join(projectsDir, "inbox")
	if entries, err := os.ReadDir(inboxDir); err == nil {
		for _, entry := range entries {
			if !strings.HasSuffix(entry.Name(), ".json") {
				continue
			}
			data, err := os.ReadFile(filepath.Join(inboxDir, entry.Name()))
			if err != nil {
				continue
			}
			var email struct {
				ID          string `json:"id"`
				From        string `json:"from"`
				FromName    string `json:"from_name"`
				Subject     string `json:"subject"`
				Body        string `json:"body"`
				TrustStatus string `json:"trust_status"`
				TrustReason string `json:"trust_reason"`
				Category    string `json:"category"`
				Read        bool   `json:"read"`
			}
			if json.Unmarshal(data, &email) == nil {
				body := email.Body
				if len(body) > 1000 {
					body = body[:1000] + "... (truncated)"
				}
				ctx.RecentEmails = append(ctx.RecentEmails, EmailSummary{
					ID:          email.ID,
					From:        email.FromName + " <" + email.From + ">",
					Subject:     email.Subject,
					Body:        body,
					Trust:       email.TrustStatus,
					TrustReason: email.TrustReason,
					Category:    email.Category,
					Read:        email.Read,
				})
			}
		}
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

	if len(ctx.Vendors) > 0 {
		sb.WriteString("\n**Vendors & Contracts:**\n")
		for _, v := range ctx.Vendors {
			active := "active"
			if !v.Active { active = "inactive" }
			sb.WriteString(fmt.Sprintf("- **%s** (%s) — contact: %s [%s]\n", v.Name, v.Domain, v.Contact, active))
			if v.ContractTerms != "" {
				sb.WriteString(fmt.Sprintf("  Contract: %s\n", v.ContractTerms))
			}
		}
	}

	if len(ctx.RecentEmails) > 0 {
		sb.WriteString(fmt.Sprintf("\n**Inbox (%d emails):**\n\n", len(ctx.RecentEmails)))
		for i, e := range ctx.RecentEmails {
			trust := ""
			if e.Trust == "impersonation" { trust = " 🚨IMPERSONATION — " + e.TrustReason }
			if e.Trust == "new_contact" { trust = " ⚠️NEW CONTACT — " + e.TrustReason }
			readStatus := "UNREAD"
			if e.Read { readStatus = "read" }
			sb.WriteString(fmt.Sprintf("**Email %d** [%s] [%s]%s\n", i+1, e.Category, readStatus, trust))
			sb.WriteString(fmt.Sprintf("From: %s\nSubject: %s\n", e.From, e.Subject))
			if e.Body != "" {
				sb.WriteString(fmt.Sprintf("Body:\n%s\n", e.Body))
			}
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
