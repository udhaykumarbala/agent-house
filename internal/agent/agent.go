package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/creack/pty"
	"pty-claude-test/internal/session"
)

// Role represents an agent's role in the team
type Role string

const (
	RoleCEO       Role = "ceo"
	RolePM        Role = "pm"
	RoleUX        Role = "ux"
	RoleUI        Role = "ui"
	RoleSecurity  Role = "security"
	RoleArchitect Role = "architect"
	RoleSeniorDev Role = "senior_dev"
	RoleJuniorDev Role = "junior_dev"
)

// PermissionModes maps agent roles to their Claude Code permission mode.
// PermissionModes maps agent roles to their Claude Code permission mode.
// acceptEdits: can read/write files freely, bash needs approval
// bypassPermissions: can do everything without approval
var PermissionModes = map[Role]string{
	RoleCEO:       "acceptEdits",       // Needs to write plans, reviews
	RolePM:        "acceptEdits",       // Needs to write research docs, specs
	RoleUX:        "acceptEdits",       // Needs to write UX specs
	RoleUI:        "acceptEdits",       // Needs to write UI specs
	RoleSecurity:  "acceptEdits",       // Needs to write security audits
	RoleArchitect: "acceptEdits",       // Needs to write template/architecture docs
	RoleSeniorDev: "bypassPermissions", // Full access for coding + testing
	RoleJuniorDev: "bypassPermissions", // Full access for coding
}

// Agent represents a specialized AI agent
type Agent struct {
	ID             string
	Name           string
	Role           Role
	SystemPrompt   string
	Color          string // For UI display
	SessionManager *session.SessionManager
}

// Response from an agent
type Response struct {
	Content    string
	DelegateTo []Role   // Downward delegation - pass work to
	ReviewTo   []Role   // Upward escalation - get decisions from
	Files      []string // Files to create
	Raw        string   // Original unprocessed output
	IsComplete bool     // Agent signals task is complete (no further delegation needed)
}

// SessionResponse is the structured response from session-based execution.
type SessionResponse struct {
	Text          string
	Events        []session.AgentEvent
	FilesCreated  []string
	FilesModified []string
	CommandsRun   []CommandRecord
	Delegations   []Role
	Reviews       []Role
	IsComplete    bool

	// Metrics
	InputTokens  int
	OutputTokens int
	CostUSD      float64
	ToolCalls    int
	Duration     time.Duration
}

// CommandRecord captures a bash command and its output.
type CommandRecord struct {
	Command string
	Output  string
	IsError bool
}

// ExecuteWithSession sends a task using the agent's configured execution mode.
// Oneshot mode: Direct API call (fast, no tools).
// Session mode: Claude Code session (persistent, full tool access).
func (a *Agent) ExecuteWithSession(ctx context.Context, projectID, workDir, taskPrompt string) (*SessionResponse, error) {
	mode := GetMode(a.Role)

	// Try oneshot mode first if configured
	if mode == ModeOneshot {
		resp, err := a.executeOneshot(ctx, projectID, taskPrompt)
		if err == nil {
			return resp, nil
		}
		// Fall back to session if oneshot fails (e.g., no API key)
		log.Printf("[AGENT:%s] Oneshot failed, falling back to session: %v", a.Role, err)
	}

	if a.SessionManager == nil {
		return nil, fmt.Errorf("session manager not set for agent %s", a.Role)
	}

	permMode := PermissionModes[a.Role]
	if permMode == "" {
		permMode = "bypassPermissions"
	}

	// Check for MCP config and build role-specific config
	mcpPath := ""
	if mcpConfig := session.LoadMCPConfig(workDir); mcpConfig != nil {
		mcpPath = session.BuildMCPConfigForRole(mcpConfig, string(a.Role))
	}

	sess, err := a.SessionManager.GetOrCreate(session.SessionConfig{
		AgentRole:      string(a.Role),
		AgentName:      a.Name,
		SystemPrompt:   a.SystemPrompt,
		WorkDir:        workDir,
		PermissionMode: permMode,
		ProjectID:      projectID,
		MaxTurns:       50,
		MCPConfigPath:  mcpPath,
	})
	if err != nil {
		return nil, fmt.Errorf("session create failed: %w", err)
	}

	startTime := time.Now()

	text, events, err := sess.SendTask(ctx, taskPrompt)
	if err != nil {
		return nil, fmt.Errorf("session task failed: %w", err)
	}

	resp := &SessionResponse{
		Text:     text,
		Events:   events,
		Duration: time.Since(startTime),
	}

	// Parse structured data from events
	for _, ev := range events {
		switch ev.Type {
		case "tool_use":
			resp.ToolCalls++
			switch ev.ToolName {
			case "Write":
				path := extractFilePath(ev.Input)
				if path != "" {
					resp.FilesCreated = append(resp.FilesCreated, path)
				}
			case "Edit":
				path := extractFilePath(ev.Input)
				if path != "" {
					resp.FilesModified = append(resp.FilesModified, path)
				}
			case "Bash":
				cmd := extractCommand(ev.Input)
				resp.CommandsRun = append(resp.CommandsRun, CommandRecord{Command: cmd})
			}
		case "tool_result":
			// Match with last command to capture output
			if len(resp.CommandsRun) > 0 {
				last := &resp.CommandsRun[len(resp.CommandsRun)-1]
				if last.Output == "" {
					last.Output = ev.Output
					last.IsError = ev.IsError
				}
			}
		case "turn_complete":
			resp.InputTokens = ev.InputTokens
			resp.OutputTokens = ev.OutputTokens
			resp.CostUSD = ev.CostUSD
		}
	}

	// Parse text for delegation/review/completion signals
	resp.Delegations = parseDelegations(text)
	resp.Reviews = parseReviews(text)
	resp.IsComplete = parseCompletionSignal(text)

	// If text parsing didn't find completion, infer from events:
	// - If the agent wrote files or ran commands, it completed its work
	// - If turn_complete arrived, the agent finished successfully
	if !resp.IsComplete {
		hasTurnComplete := false
		for _, ev := range events {
			if ev.Type == "turn_complete" {
				hasTurnComplete = true
				break
			}
		}
		if hasTurnComplete && (len(resp.FilesCreated) > 0 || len(resp.FilesModified) > 0 || resp.ToolCalls > 0) {
			resp.IsComplete = true
		}
	}

	return resp, nil
}

// apiClient is a singleton API client for oneshot calls.
var (
	apiClient     *session.APIClient
	apiClientOnce sync.Once
)

func getAPIClient() *session.APIClient {
	apiClientOnce.Do(func() {
		apiClient = session.NewAPIClient()
	})
	return apiClient
}

// executeOneshot sends a task via direct API call (no tools, fast).
func (a *Agent) executeOneshot(ctx context.Context, projectID, taskPrompt string) (*SessionResponse, error) {
	client := getAPIClient()
	if !client.IsAvailable() {
		return nil, fmt.Errorf("API client not available")
	}

	startTime := time.Now()

	apiResp, err := client.SendMessage(ctx, a.SystemPrompt, taskPrompt)
	if err != nil {
		return nil, fmt.Errorf("oneshot API call failed: %w", err)
	}

	duration := time.Since(startTime)

	log.Printf("[AGENT:%s] Oneshot completed in %dms — %d in / %d out tokens",
		a.Role, apiResp.DurationMS, apiResp.Usage.InputTokens, apiResp.Usage.OutputTokens)

	// Build SessionResponse compatible with the rest of the pipeline
	resp := &SessionResponse{
		Text:         apiResp.Text,
		IsComplete:   true,
		InputTokens:  apiResp.Usage.InputTokens,
		OutputTokens: apiResp.Usage.OutputTokens,
		Duration:     duration,
		Events: []session.AgentEvent{
			{
				Type:         "session_meta",
				AgentRole:    string(a.Role),
				AgentName:    a.Name,
				ProjectID:    projectID,
				Model:        apiResp.Model,
				Timestamp:    startTime.UnixMilli(),
			},
			{
				Type:         "text_delta",
				AgentRole:    string(a.Role),
				AgentName:    a.Name,
				ProjectID:    projectID,
				Content:      apiResp.Text,
				Timestamp:    time.Now().UnixMilli(),
			},
			{
				Type:         "turn_complete",
				AgentRole:    string(a.Role),
				AgentName:    a.Name,
				ProjectID:    projectID,
				InputTokens:  apiResp.Usage.InputTokens,
				OutputTokens: apiResp.Usage.OutputTokens,
				Timestamp:    time.Now().UnixMilli(),
			},
		},
	}

	// Parse text for signals
	resp.Delegations = parseDelegations(apiResp.Text)
	resp.Reviews = parseReviews(apiResp.Text)

	return resp, nil
}

// extractFilePath extracts file_path from a tool input JSON string.
func extractFilePath(inputJSON string) string {
	var input map[string]interface{}
	if err := json.Unmarshal([]byte(inputJSON), &input); err != nil {
		return ""
	}
	if path, ok := input["file_path"].(string); ok {
		return path
	}
	return ""
}

// extractCommand extracts command from a Bash tool input JSON string.
func extractCommand(inputJSON string) string {
	var input map[string]interface{}
	if err := json.Unmarshal([]byte(inputJSON), &input); err != nil {
		return ""
	}
	if cmd, ok := input["command"].(string); ok {
		return cmd
	}
	return ""
}

// NewAgent creates a new agent with the given role
func NewAgent(role Role) (*Agent, error) {
	agent := &Agent{
		ID:   string(role),
		Role: role,
	}

	// Set agent-specific properties
	switch role {
	case RoleCEO:
		agent.Name = "CEO"
		agent.Color = "#FF6B6B" // Red
	case RolePM:
		agent.Name = "Product Manager"
		agent.Color = "#4ECDC4" // Teal
	case RoleUX:
		agent.Name = "UX Designer"
		agent.Color = "#45B7D1" // Blue
	case RoleUI:
		agent.Name = "UI Designer"
		agent.Color = "#96CEB4" // Green
	case RoleSecurity:
		agent.Name = "Security Expert"
		agent.Color = "#FFEAA7" // Yellow
	case RoleArchitect:
		agent.Name = "Architect"
		agent.Color = "#DDA0DD" // Plum
	case RoleSeniorDev:
		agent.Name = "Senior Developer"
		agent.Color = "#98D8C8" // Mint
	case RoleJuniorDev:
		agent.Name = "Junior Developer"
		agent.Color = "#F7DC6F" // Gold
	default:
		return nil, fmt.Errorf("unknown role: %s", role)
	}

	// Load system prompt from file
	promptPath := fmt.Sprintf("prompts/%s.md", role)
	promptData, err := os.ReadFile(promptPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load prompt for %s: %w", role, err)
	}
	agent.SystemPrompt = string(promptData)

	return agent, nil
}

// Process sends a task to the agent and returns the response
func (a *Agent) Process(task string) (*Response, error) {
	return a.ProcessInDir(task, "")
}

// ProcessInDir sends a task to the agent with a specific working directory
func (a *Agent) ProcessInDir(task string, workDir string) (*Response, error) {
	log.Printf("[AGENT:%s] Processing task in dir=%s", a.Role, workDir)
	taskPreview := task
	if len(taskPreview) > 150 {
		taskPreview = taskPreview[:150] + "..."
	}
	log.Printf("[AGENT:%s] Task: %s", a.Role, taskPreview)

	// Build the full prompt with system context
	fullPrompt := fmt.Sprintf(`%s

---
TASK:
%s

---
Respond with your analysis and recommendations. If you need other team members involved, clearly state who and why.
If you create files, they will be created in the project directory.`, a.SystemPrompt, task)

	// Call Claude with --print mode
	output, err := runClaudePrint(fullPrompt, workDir)
	if err != nil {
		log.Printf("[AGENT:%s] FAILED: %v", a.Role, err)
		return nil, fmt.Errorf("claude error: %w", err)
	}
	log.Printf("[AGENT:%s] SUCCESS, response length: %d", a.Role, len(output))

	// Parse the response
	response := &Response{
		Raw:     output,
		Content: strings.TrimSpace(output),
	}

	// Parse delegation instructions
	response.DelegateTo = parseDelegations(output)

	// Parse review escalations
	response.ReviewTo = parseReviews(output)

	// Check for completion signal
	response.IsComplete = parseCompletionSignal(output)

	return response, nil
}

// parseCompletionSignal checks if the agent signaled completion
func parseCompletionSignal(output string) bool {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "COMPLETE:") || strings.HasPrefix(line, "TASK_COMPLETE:") {
			return true
		}
	}
	return false
}

// parseReviews extracts review escalation instructions from agent response
func parseReviews(output string) []Role {
	var reviews []Role

	// Look for REVIEW: block
	lines := strings.Split(output, "\n")
	inReviewBlock := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "REVIEW:") {
			inReviewBlock = true
			continue
		}

		if inReviewBlock {
			// End of block if we hit empty line or non-list item
			if line == "" || line == "```" {
				break
			}

			// Parse "- agent: reason" format
			if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
				parts := strings.SplitN(line[2:], ":", 2)
				if len(parts) >= 1 {
					agentID := strings.TrimSpace(parts[0])
					role := Role(agentID)
					if isValidRole(role) {
						reviews = append(reviews, role)
					}
				}
			}
		}
	}

	return reviews
}

// parseDelegations extracts delegate instructions from agent response
func parseDelegations(output string) []Role {
	var delegations []Role

	// Look for DELEGATE: block
	lines := strings.Split(output, "\n")
	inDelegateBlock := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "DELEGATE:") {
			inDelegateBlock = true
			continue
		}

		if inDelegateBlock {
			// End of block if we hit empty line or non-list item
			if line == "" || line == "```" {
				break
			}

			// Parse "- agent: reason" format
			if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
				parts := strings.SplitN(line[2:], ":", 2)
				if len(parts) >= 1 {
					agentID := strings.TrimSpace(parts[0])
					role := Role(agentID)
					if isValidRole(role) {
						delegations = append(delegations, role)
					}
				}
			}
		}
	}

	return delegations
}

// isValidRole checks if a role is valid
func isValidRole(r Role) bool {
	switch r {
	case RoleCEO, RolePM, RoleUX, RoleUI, RoleSecurity, RoleArchitect, RoleSeniorDev, RoleJuniorDev:
		return true
	default:
		return false
	}
}

// AgentTimeout is the maximum time a single Claude CLI invocation can run.
const AgentTimeout = 10 * time.Minute

// runClaudePrint calls claude with tool permissions
func runClaudePrint(prompt string, workDir string) (string, error) {
	args := []string{
		"--print",
		"--tools", "Write,Read,Edit,Bash",
		"--dangerously-skip-permissions",
		prompt,
	}

	cmd := exec.Command("claude", args...)

	// Set working directory for file operations
	if workDir != "" {
		cmd.Dir = workDir
	}

	// Log command details
	promptPreview := prompt
	if len(promptPreview) > 200 {
		promptPreview = promptPreview[:200] + "..."
	}
	log.Printf("[AGENT] Executing: claude --print --tools Write,Read,Edit,Bash --dangerously-skip-permissions")
	log.Printf("[AGENT] Working dir: %s", workDir)
	log.Printf("[AGENT] Prompt preview: %s", promptPreview)

	startTime := time.Now()

	// Use PTY for proper terminal emulation
	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Printf("[AGENT] ERROR PTY start failed: %v", err)
		return "", fmt.Errorf("pty start: %w", err)
	}
	defer ptmx.Close()

	// Capture output
	var output bytes.Buffer
	done := make(chan error, 1)

	go func() {
		_, err := io.Copy(&output, ptmx)
		done <- err
	}()

	// Wait for command with timeout
	waitCh := make(chan error, 1)
	go func() {
		waitCh <- cmd.Wait()
	}()

	var cmdErr error
	select {
	case cmdErr = <-waitCh:
		// Command finished normally
	case <-time.After(AgentTimeout):
		// Kill the process on timeout
		log.Printf("[AGENT] ERROR Timeout after %v, killing process", AgentTimeout)
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		<-waitCh // Wait for process to actually exit
		return output.String(), fmt.Errorf("agent timed out after %v", AgentTimeout)
	}

	<-done

	elapsed := time.Since(startTime)
	rawOutput := output.String()

	if cmdErr != nil {
		log.Printf("[AGENT] ERROR Command failed after %v: %v", elapsed, cmdErr)
		log.Printf("[AGENT] ERROR Raw output (%d bytes): %s", len(rawOutput), truncate(rawOutput, 500))
		// Also try to get exit code
		if exitErr, ok := cmdErr.(*exec.ExitError); ok {
			log.Printf("[AGENT] ERROR Exit code: %d, stderr: %s", exitErr.ExitCode(), string(exitErr.Stderr))
		}
		return rawOutput, fmt.Errorf("command error (after %v): %w", elapsed, cmdErr)
	}

	cleanOutput := stripANSI(rawOutput)
	log.Printf("[AGENT] OK Completed in %v, output: %d bytes", elapsed, len(cleanOutput))
	if len(cleanOutput) > 0 {
		log.Printf("[AGENT] Response preview: %s", truncate(cleanOutput, 300))
	}

	return cleanOutput, nil
}

// truncate returns the first n characters of s
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

// stripANSI removes ANSI escape codes from text
func stripANSI(input string) string {
	// Simple regex-free approach for common codes
	result := input

	// Remove common escape sequences
	for strings.Contains(result, "\x1b[") {
		start := strings.Index(result, "\x1b[")
		if start == -1 {
			break
		}
		// Find the end of the escape sequence
		end := start + 2
		for end < len(result) && !isANSITerminator(result[end]) {
			end++
		}
		if end < len(result) {
			end++ // Include the terminator
		}
		result = result[:start] + result[end:]
	}

	// Clean up carriage returns
	result = strings.ReplaceAll(result, "\r\n", "\n")
	result = strings.ReplaceAll(result, "\r", "\n")

	return strings.TrimSpace(result)
}

func isANSITerminator(b byte) bool {
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z')
}
