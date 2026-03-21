package session

import (
	"encoding/json"
	"time"
)

// Session states
const (
	StateActive     = "active"
	StateIdle       = "idle"
	StateTerminated = "terminated"
)

// ClaudeEvent is the raw NDJSON envelope from the Claude CLI stdout.
type ClaudeEvent struct {
	Type          string          `json:"type"`                      // system, assistant, user, progress, result, rate_limit_event
	Subtype       string          `json:"subtype,omitempty"`         // init, result, success, error
	Message       json.RawMessage `json:"message,omitempty"`         // ClaudeMessage JSON
	Result        json.RawMessage `json:"result,omitempty"`          // result data
	SessionID     string          `json:"session_id,omitempty"`      // from system/init
	Model         string          `json:"model,omitempty"`           // from system/init
	DurationMS    int64           `json:"duration_ms,omitempty"`
	DurationAPIMS int64           `json:"duration_api_ms,omitempty"`
	CostUSD       float64         `json:"cost_usd,omitempty"`
	TotalCostUSD  float64         `json:"total_cost_usd,omitempty"`
	IsError       bool            `json:"is_error,omitempty"`
	NumTurns      int             `json:"num_turns,omitempty"`
	ModelUsage    json.RawMessage `json:"modelUsage,omitempty"`
	Usage         json.RawMessage `json:"usage,omitempty"`
}

// ClaudeMessage is the message payload within a ClaudeEvent.
type ClaudeMessage struct {
	ID         string          `json:"id,omitempty"`
	Content    []ClaudeContent `json:"content"`
	Role       string          `json:"role,omitempty"`
	StopReason string          `json:"stop_reason,omitempty"`
}

// ClaudeContent is a polymorphic content block.
type ClaudeContent struct {
	Type      string          `json:"type"`                    // text, thinking, tool_use, tool_result
	Text      string          `json:"text,omitempty"`          // for text blocks
	Thinking  string          `json:"thinking,omitempty"`      // for thinking blocks
	ID        string          `json:"id,omitempty"`            // for tool_use
	Name      string          `json:"name,omitempty"`          // tool name
	Input     json.RawMessage `json:"input,omitempty"`         // tool input
	Content   json.RawMessage `json:"content,omitempty"`       // for tool_result: string or [{type,text}] array
	ToolUseID string          `json:"tool_use_id,omitempty"`   // for tool_result
	IsError   bool            `json:"is_error,omitempty"`      // for tool_result
}

// AgentEvent is a normalized event emitted to the orchestrator and dashboard.
// Every event is tagged with agent identity and project context.
type AgentEvent struct {
	Type         string  `json:"type"`                       // text_delta, thinking_delta, tool_use, tool_result, turn_complete, error, session_meta
	AgentRole    string  `json:"agent_role"`
	AgentName    string  `json:"agent_name"`
	ProjectID    string  `json:"project_id,omitempty"`
	TaskID       string  `json:"task_id,omitempty"`
	Content      string  `json:"content,omitempty"`
	ToolUseID    string  `json:"tool_use_id,omitempty"`
	ToolName     string  `json:"tool_name,omitempty"`
	Input        string  `json:"input,omitempty"`
	Output       string  `json:"output,omitempty"`
	IsError      bool    `json:"is_error,omitempty"`
	Model        string  `json:"model,omitempty"`
	InputTokens  int     `json:"input_tokens,omitempty"`
	OutputTokens int     `json:"output_tokens,omitempty"`
	CostUSD      float64 `json:"cost_usd,omitempty"`
	Timestamp    int64   `json:"timestamp"`
}

// SessionConfig holds configuration for spawning an agent session.
type SessionConfig struct {
	AgentRole      string // role identifier (ceo, pm, etc.)
	AgentName      string // display name
	SystemPrompt   string // role-specific system prompt
	WorkDir        string // project working directory
	PermissionMode string // default, plan, acceptEdits, bypassPermissions
	ProjectID      string // project identifier
	Model          string // optional model override
	MaxTurns       int    // safety limit per task (default: 50)
	MCPConfigPath  string // path to MCP config JSON (optional — gives agent external tools)
}

// SessionStats holds accumulated metrics for a session.
type SessionStats struct {
	TotalInputTokens  int            `json:"total_input_tokens"`
	TotalOutputTokens int            `json:"total_output_tokens"`
	TotalCostUSD      float64        `json:"total_cost_usd"`
	TotalTurns        int            `json:"total_turns"`
	TotalToolCalls    int            `json:"total_tool_calls"`
	ToolBreakdown     map[string]int `json:"tool_breakdown"`
}

// SessionInfo is the JSON response for session listing.
type SessionInfo struct {
	ID              string       `json:"id"`
	AgentRole       string       `json:"agent_role"`
	AgentName       string       `json:"agent_name"`
	ProjectID       string       `json:"project_id"`
	State           string       `json:"state"`
	ClaudeSessionID string       `json:"claude_session_id,omitempty"`
	Model           string       `json:"model,omitempty"`
	Stats           SessionStats `json:"stats"`
	CreatedAt       string       `json:"created_at"`
}

// NewAgentEvent creates a new AgentEvent with timestamp and agent info.
func NewAgentEvent(eventType, agentRole, agentName, projectID, taskID string) AgentEvent {
	return AgentEvent{
		Type:      eventType,
		AgentRole: agentRole,
		AgentName: agentName,
		ProjectID: projectID,
		TaskID:    taskID,
		Timestamp: time.Now().UnixMilli(),
	}
}
