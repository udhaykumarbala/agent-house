package agent

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/creack/pty"
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

// Agent represents a specialized AI agent
type Agent struct {
	ID           string
	Name         string
	Role         Role
	SystemPrompt string
	Color        string // For UI display
}

// Response from an agent
type Response struct {
	Content    string
	DelegateTo []Role   // Agents to involve next
	Files      []string // Files to create
	Raw        string   // Original unprocessed output
	IsComplete bool     // Agent signals task is complete (no further delegation needed)
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
		return nil, fmt.Errorf("claude error: %w", err)
	}

	// Parse the response
	response := &Response{
		Raw:     output,
		Content: strings.TrimSpace(output),
	}

	// Parse delegation instructions
	response.DelegateTo = parseDelegations(output)

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

// runClaudePrint calls claude with tool permissions
func runClaudePrint(prompt string, workDir string) (string, error) {
	// Use claude with allowed tools for file operations
	// --tools enables tools in print mode
	// --dangerously-skip-permissions bypasses permission prompts
	cmd := exec.Command("claude",
		"--print",
		"--tools", "Write,Read,Edit,Bash",
		"--dangerously-skip-permissions",
		prompt,
	)

	// Set working directory for file operations
	if workDir != "" {
		cmd.Dir = workDir
	}

	// Use PTY for proper terminal emulation
	ptmx, err := pty.Start(cmd)
	if err != nil {
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

	// Wait for command to complete
	cmdErr := cmd.Wait()
	<-done

	if cmdErr != nil {
		return output.String(), fmt.Errorf("command error: %w", cmdErr)
	}

	// Strip ANSI codes
	return stripANSI(output.String()), nil
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
