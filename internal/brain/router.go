package brain

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"pty-claude-test/internal/session"
)

// Router handles Brain decision-making via oneshot API calls.
type Router struct {
	apiClient    *session.APIClient
	projectsDir  string
	conversations *ConversationStore
	systemPrompt string
}

// NewRouter creates a Brain router.
func NewRouter(apiClient *session.APIClient, projectsDir string, conversations *ConversationStore) *Router {
	return &Router{
		apiClient:     apiClient,
		projectsDir:   projectsDir,
		conversations: conversations,
		systemPrompt:  brainSystemPrompt,
	}
}

// Route processes a user message and returns a BrainDecision.
// This is the core Brain loop: assemble context → call API → parse decision.
// Supports re-calls for server-mediated actions (read_file, etc.)
func (r *Router) Route(ctx context.Context, userID, message string) (*BrainDecision, error) {
	// Store user message
	r.conversations.Add(userID, "user", message)

	// Assemble workspace context
	workspaceCtx := AssembleContext(r.projectsDir)

	// Get recent conversation history
	history := r.conversations.GetRecent(userID, 10)

	// Build the full prompt
	prompt := r.buildPrompt(workspaceCtx, history, message)

	// Call Brain (oneshot)
	decision, err := r.callBrain(ctx, prompt)
	if err != nil {
		return nil, err
	}

	// Handle server-mediated actions (re-call loop, max 3)
	for i := 0; i < 3; i++ {
		if decision.Action != ActionReadFile {
			break
		}

		// Read the requested file
		filePath := decision.Params["path"]
		project := decision.Params["project"]
		if project != "" && !strings.Contains(filePath, "/") {
			filePath = project + "/" + filePath
		}

		fullPath := r.projectsDir + "/" + filePath
		content, err := readFileContent(fullPath)
		if err != nil {
			content = fmt.Sprintf("Error reading file: %v", err)
		}

		// Re-call Brain with file content
		followUp := fmt.Sprintf("Here's the content of %s:\n\n```\n%s\n```\n\nNow respond to the user's original request.", filePath, truncateFile(content, 8000))

		r.conversations.Add(userID, "assistant", fmt.Sprintf("[Reading %s...]", filePath))
		prompt = r.buildPrompt(workspaceCtx, r.conversations.GetRecent(userID, 12), followUp)
		decision, err = r.callBrain(ctx, prompt)
		if err != nil {
			return nil, err
		}
	}

	// If still asking for files after 3 reads, escalate
	if decision.Action == ActionReadFile {
		decision.Action = ActionEscalate
		decision.Response = "This requires deeper analysis. Let me look into it more thoroughly."
	}

	// Store assistant response
	if decision.Response != "" {
		r.conversations.Add(userID, "assistant", decision.Response)
	}

	return decision, nil
}

func (r *Router) callBrain(ctx context.Context, prompt string) (*BrainDecision, error) {
	start := time.Now()

	resp, err := r.apiClient.SendMessage(ctx, r.systemPrompt, prompt)
	if err != nil {
		return nil, fmt.Errorf("brain API call failed: %w", err)
	}

	log.Printf("[BRAIN] API call: %dms, %d in / %d out tokens",
		time.Since(start).Milliseconds(), resp.Usage.InputTokens, resp.Usage.OutputTokens)

	// Parse structured JSON from response
	decision, err := parseDecision(resp.Text)
	if err != nil {
		// If parsing fails, treat the whole response as a text reply
		log.Printf("[BRAIN] Failed to parse decision JSON, treating as text response: %v", err)
		return &BrainDecision{
			Action:   ActionRespond,
			Response: resp.Text,
		}, nil
	}

	log.Printf("[BRAIN] Decision: action=%s", decision.Action)
	return decision, nil
}

func (r *Router) buildPrompt(workspaceCtx string, history []ChatMessage, currentMessage string) string {
	var sb strings.Builder

	sb.WriteString(workspaceCtx)
	sb.WriteString("\n")

	if len(history) > 1 {
		sb.WriteString("## Recent Conversation\n\n")
		// Show all but the last message (which is the current one)
		for _, msg := range history[:len(history)-1] {
			if msg.Role == "user" {
				sb.WriteString("**User:** " + msg.Content + "\n\n")
			} else {
				sb.WriteString("**Assistant:** " + msg.Content + "\n\n")
			}
		}
	}

	sb.WriteString("## Current Message\n\n")
	sb.WriteString(currentMessage)

	return sb.String()
}

func parseDecision(text string) (*BrainDecision, error) {
	text = strings.TrimSpace(text)

	// Try to extract JSON from the response (may be wrapped in markdown code blocks)
	jsonStr := text
	if idx := strings.Index(text, "{"); idx >= 0 {
		// Find the matching closing brace
		depth := 0
		for i := idx; i < len(text); i++ {
			if text[i] == '{' {
				depth++
			} else if text[i] == '}' {
				depth--
				if depth == 0 {
					jsonStr = text[idx : i+1]
					break
				}
			}
		}
	}

	var decision BrainDecision
	if err := json.Unmarshal([]byte(jsonStr), &decision); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w (text: %.200s)", err, text)
	}

	if decision.Action == "" {
		return nil, fmt.Errorf("missing action field")
	}

	return &decision, nil
}

func readFileContent(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func truncateFile(content string, maxLen int) string {
	if len(content) <= maxLen {
		return content
	}
	return content[:maxLen] + "\n... (truncated)"
}

// brainSystemPrompt is the core system prompt for the Brain agent.
var brainSystemPrompt = `You are the Brain of Agent House — an AI workspace manager.

You manage projects built by teams of AI agents (CEO, PM, UX, UI, Security, Architect, Senior Dev, Junior Dev).

## How to Respond

ALWAYS respond with a single JSON object. No other text outside the JSON.

## Available Actions

{"action": "respond", "response": "Your answer here"}
  → Answer the user directly. Use for questions, status requests, greetings.

{"action": "create_project", "params": {"project_id": "my-project", "task": "Build a..."}, "response": "Creating project..."}
  → Start a new project. Choose a short, lowercase, hyphenated project_id.

{"action": "delegate", "params": {"project_id": "existing-project", "agent_role": "senior_dev", "task": "Add dark mode"}, "response": "Delegating to Senior Dev..."}
  → Send a task to a specific agent in an existing project.

{"action": "project_status", "params": {"project_id": "my-project"}, "response": "Let me check..."}
  → Get detailed status of a project (server will provide and re-call you).

{"action": "read_file", "params": {"project": "my-project", "path": "src/index.html"}, "response": "Reading file..."}
  → Read a file from a project. Server reads it and re-calls you with content.

{"action": "list_projects", "response": "Here are your projects:"}
  → List all projects with status.

{"action": "escalate", "params": {"project_id": "my-project", "task": "Analyze auth security"}, "response": "This needs deep analysis..."}
  → For complex multi-file analysis. Spawns a full Claude Code session.

## Decision Rules

1. Greeting or simple question → respond
2. "What projects exist?" / "What's running?" → list_projects
3. "How's project X?" / "Status of X" → project_status
4. "Build/Create/Make [something]" → create_project (pick a good project_id)
5. "Change/Add/Fix [something] in [project]" → delegate to appropriate agent
6. "Show me [file] from [project]" → read_file
7. Complex analysis across multiple files → escalate
8. If unsure, ask for clarification via respond

## Project ID Rules
- Lowercase, hyphens only: "landing-page", "todo-app", "api-server"
- Short and descriptive
- No spaces or special characters

## Agent Roles for Delegation
- senior_dev: Code changes, new features, bug fixes
- junior_dev: Simple code tasks, file creation
- architect: Architecture decisions, structure changes
- pm: Requirements, documentation, specs
- security: Security audits, vulnerability fixes
- ux: UX improvements, flow changes
- ui: Visual design, styling changes
- ceo: Strategic decisions, reviews`
