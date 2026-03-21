# Agent House v2 — Claude Code Session-Per-Agent Integration

**Date:** 2026-03-21
**Status:** Implementation Plan
**Goal:** Replace one-shot `claude --print` calls with persistent Claude Code streaming sessions per agent, giving full tool access, real-time visibility, and detailed activity reports.

---

## Table of Contents

1. [Architecture Delta](#1-architecture-delta)
2. [New Package: internal/session](#2-new-package-internalsession)
3. [Modify: internal/agent](#3-modify-internalagent)
4. [Modify: internal/orchestrator](#4-modify-internalorchestrator)
5. [Modify: internal/web (WebSocket + API)](#5-modify-internalweb)
6. [Dashboard UI: Per-Agent Live Panels](#6-dashboard-ui-per-agent-live-panels)
7. [Human-in-the-Loop via Tool Approval](#7-human-in-the-loop-via-tool-approval)
8. [Cost & Token Tracking](#8-cost--token-tracking)
9. [Error Recovery & Safety](#9-error-recovery--safety)
10. [API Changes](#10-api-changes)
11. [WebSocket Protocol Changes](#11-websocket-protocol-changes)
12. [Migration Path](#12-migration-path)
13. [Step-by-Step Build Order](#13-step-by-step-build-order)
14. [Testing Plan](#14-testing-plan)

---

## 1. Architecture Delta

### Before (Current)

```
User → Orchestrator → Agent.Execute("claude --print -p <prompt>")
                         ↓
                    One-shot CLI call
                    Returns full text
                    Parsed for DELEGATE:/REVIEW:/files
                         ↓
                    Orchestrator gets text → next agent
```

- Agent = system prompt + single CLI invocation → returns text blob
- No streaming, no tool visibility, no thinking visibility
- PTY-based I/O via creack/pty
- Response parsed with regex for delegation signals

### After (New)

```
User → Orchestrator → AgentSession (persistent Claude Code process)
                         ↓
                    stdin: {"type":"user","message":{"role":"user","content":"<task>"}}
                    stdout: NDJSON stream of events
                         ↓
                    text_delta, thinking_delta, tool_use, tool_result, turn_complete
                         ↓
                    Orchestrator reads events in real-time
                    WebSocket hub broadcasts per-agent events to dashboard
                    Turn completes → orchestrator proceeds
```

- Agent = persistent Claude Code session with streaming I/O
- Full tool access (Read, Edit, Write, Bash, Grep, Glob, etc.)
- Real-time visibility into every action
- Structured event protocol (not regex parsing)
- Session persists across multiple tasks (agent has memory)

### Key Architectural Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Workspace model | Shared directory, orchestrator-controlled turns | Simpler, agents see each other's work naturally |
| Session lifetime | Per-project (created on first use, alive until project ends) | Agent retains context across phases |
| Permission mode | Role-based (`plan` for non-devs, `acceptEdits` for devs) | Non-devs research/plan but don't modify code |
| System prompt delivery | `--append-system-prompt` flag per session | Clean separation, no CLAUDE.md conflicts |
| Inter-agent context | Orchestrator passes summary via stdin messages | Orchestrator remains the coordinator |
| Parallel execution | Multiple sessions run concurrently, each in own goroutine | Same as current worker pool model |

---

## 2. New Package: `internal/session`

Adapted from mobile-code's session management. This is the core new infrastructure.

### Files to Create

#### 2.1 `internal/session/types.go`

All data structures for Claude Code NDJSON protocol.

```go
package session

// ClaudeEvent — raw NDJSON line from Claude Code stdout
type ClaudeEvent struct {
    Type       string              `json:"type"`        // system, assistant, user, result
    Subtype    string              `json:"subtype"`     // init, result, success, error
    Message    json.RawMessage     `json:"message"`     // ClaudeMessage JSON
    Result     json.RawMessage     `json:"result"`      // result data
    SessionID  string              `json:"session_id"`
    Model      string              `json:"model"`
    IsError    bool                `json:"is_error"`
    CostUSD    float64             `json:"total_cost_usd"`
    ModelUsage json.RawMessage     `json:"modelUsage"`
}

// ClaudeMessage — message payload inside ClaudeEvent
type ClaudeMessage struct {
    Content    []ClaudeContent     `json:"content"`
    Role       string              `json:"role"`
    StopReason string              `json:"stop_reason"`
}

// ClaudeContent — polymorphic content block
type ClaudeContent struct {
    Type      string              `json:"type"`       // text, thinking, tool_use, tool_result
    Text      string              `json:"text"`
    Thinking  string              `json:"thinking"`
    ID        string              `json:"id"`          // tool_use ID
    Name      string              `json:"name"`        // tool name
    Input     json.RawMessage     `json:"input"`
    Content   json.RawMessage     `json:"content"`     // tool_result content
    ToolUseID string              `json:"tool_use_id"`
    IsError   bool                `json:"is_error"`
}

// AgentEvent — normalized event emitted to orchestrator and dashboard
type AgentEvent struct {
    Type         string  `json:"type"`          // text_delta, thinking_delta, tool_use, tool_result, turn_complete, error, session_meta
    AgentRole    string  `json:"agent_role"`    // ceo, pm, ux, ui, security, architect, senior_dev, junior_dev
    AgentName    string  `json:"agent_name"`    // display name
    ProjectID    string  `json:"project_id"`
    TaskID       string  `json:"task_id"`
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

// SessionConfig — configuration for spawning an agent session
type SessionConfig struct {
    AgentRole      string   // role identifier
    AgentName      string   // display name
    SystemPrompt   string   // role-specific system prompt
    WorkDir        string   // project working directory
    PermissionMode string   // default, plan, acceptEdits, bypassPermissions
    ProjectID      string   // project identifier
    AllowedTools   []string // optional tool whitelist
    Model          string   // optional model override
    MaxTurns       int      // safety limit per task (default: 50)
}

// SessionStats — accumulated metrics for a session
type SessionStats struct {
    TotalInputTokens  int     `json:"total_input_tokens"`
    TotalOutputTokens int     `json:"total_output_tokens"`
    TotalCostUSD      float64 `json:"total_cost_usd"`
    TotalTurns        int     `json:"total_turns"`
    TotalToolCalls    int     `json:"total_tool_calls"`
    ToolBreakdown     map[string]int `json:"tool_breakdown"` // tool_name → call count
}
```

#### 2.2 `internal/session/translator.go`

Ported from mobile-code's `translator.go`. Converts Claude NDJSON → AgentEvent.

```go
package session

type Translator struct {
    lastContentIndex int
    agentRole        string
    agentName        string
    projectID        string
    taskID           string
}

func NewTranslator(role, name, projectID, taskID string) *Translator

func (t *Translator) SetTaskID(taskID string)
func (t *Translator) Reset()
func (t *Translator) Translate(raw ClaudeEvent) []AgentEvent
// Internal methods:
// translateSystem, translateResult, translateAssistant, translateUser
// extractContent — handles polymorphic tool_result content
```

Key difference from mobile-code: Every emitted `AgentEvent` is tagged with `AgentRole`, `AgentName`, `ProjectID`, `TaskID`, and `Timestamp`. This allows the dashboard to route events to the correct agent panel.

#### 2.3 `internal/session/session.go`

The core session lifecycle manager. One per agent.

```go
package session

type AgentSession struct {
    ID              string              // UUID
    Config          SessionConfig       // role, prompt, workdir, etc.
    ClaudeSessionID string              // Claude's internal session ID
    State           string              // active, idle, terminated
    Stats           SessionStats        // accumulated metrics
    CreatedAt       time.Time

    cmd             *exec.Cmd
    stdin           io.WriteCloser
    stdout          io.ReadCloser
    stderr          io.ReadCloser
    translator      *Translator
    mu              sync.Mutex

    // Event fan-out
    listeners       []chan AgentEvent    // orchestrator + websocket hub subscribe here
    done            chan struct{}        // closed when process exits
}
```

**Public Methods:**

```go
// Create and start a new agent session
func NewAgentSession(config SessionConfig) (*AgentSession, error)

// Send a task to the agent and stream events until turn completes
// Returns the final text response and any error
func (s *AgentSession) SendTask(ctx context.Context, prompt string) (response string, events []AgentEvent, err error)

// Send a follow-up message (agent retains conversation context)
func (s *AgentSession) SendFollowUp(ctx context.Context, prompt string) (string, []AgentEvent, error)

// Approve a tool use request
func (s *AgentSession) ApproveTool(toolUseID string) error

// Deny a tool use request
func (s *AgentSession) DenyTool(toolUseID string) error

// Abort the current turn (sends SIGINT)
func (s *AgentSession) Abort() error

// Subscribe to real-time events (returns channel and unsubscribe func)
func (s *AgentSession) Subscribe() (<-chan AgentEvent, func())

// Get accumulated stats
func (s *AgentSession) GetStats() SessionStats

// Check if session process is still alive
func (s *AgentSession) IsAlive() bool

// Gracefully terminate the session
func (s *AgentSession) Terminate() error
```

**Internal Methods:**

```go
// Spawn claude process with correct flags
func (s *AgentSession) spawn() error

// Read NDJSON from stdout, translate, fan out
func (s *AgentSession) readLoop()

// Drain stderr to prevent pipe blocking
func (s *AgentSession) drainStderr()

// Wait for process exit
func (s *AgentSession) watchExit()

// Fan out event to all listeners
func (s *AgentSession) fanOut(ev AgentEvent)

// Write JSON message to stdin
func (s *AgentSession) writeStdin(msg interface{}) error
```

**Claude CLI Command Construction:**

```go
func buildCommand(config SessionConfig) *exec.Cmd {
    args := []string{
        "-p",                                    // interactive/persistent mode
        "--input-format", "stream-json",         // NDJSON input
        "--output-format", "stream-json",        // NDJSON output
        "--verbose",                             // detailed events
    }

    if config.PermissionMode != "" {
        args = append(args, "--permission-mode", config.PermissionMode)
    }

    if config.SystemPrompt != "" {
        args = append(args, "--append-system-prompt", config.SystemPrompt)
    }

    if config.Model != "" {
        args = append(args, "--model", config.Model)
    }

    if config.MaxTurns > 0 {
        args = append(args, "--max-turns", strconv.Itoa(config.MaxTurns))
    }

    cmd := exec.Command("claude", args...)
    cmd.Dir = config.WorkDir

    // Filter CLAUDECODE env to prevent nesting issues
    env := os.Environ()
    filtered := make([]string, 0, len(env))
    for _, e := range env {
        if !strings.HasPrefix(e, "CLAUDECODE=") {
            filtered = append(filtered, e)
        }
    }
    cmd.Env = filtered

    return cmd
}
```

**SendTask Flow:**

```go
func (s *AgentSession) SendTask(ctx context.Context, prompt string) (string, []AgentEvent, error) {
    // 1. If session not spawned yet, spawn it
    if s.cmd == nil {
        if err := s.spawn(); err != nil {
            return "", nil, fmt.Errorf("spawn failed: %w", err)
        }
    }

    // 2. Subscribe to events for this turn
    events := make([]AgentEvent, 0, 100)
    var response strings.Builder

    // 3. Write user message to stdin
    msg := map[string]interface{}{
        "type": "user",
        "message": map[string]interface{}{
            "role":    "user",
            "content": prompt,
        },
    }
    if err := s.writeStdin(msg); err != nil {
        return "", nil, err
    }

    // 4. Read events until turn_complete or context cancelled
    evCh, unsub := s.Subscribe()
    defer unsub()

    for {
        select {
        case ev := <-evCh:
            events = append(events, ev)
            if ev.Type == "text_delta" {
                response.WriteString(ev.Content)
            }
            if ev.Type == "turn_complete" {
                return response.String(), events, nil
            }
            if ev.Type == "error" {
                return response.String(), events, fmt.Errorf("agent error: %s", ev.Content)
            }
        case <-ctx.Done():
            s.Abort()
            return response.String(), events, ctx.Err()
        case <-s.done:
            return response.String(), events, fmt.Errorf("agent process exited unexpectedly")
        }
    }
}
```

#### 2.4 `internal/session/manager.go`

Manages all agent sessions across projects.

```go
package session

type SessionManager struct {
    mu       sync.RWMutex
    sessions map[string]*AgentSession // key: "{projectID}:{agentRole}"

    // Global event bus — all agent events flow here
    eventBus chan AgentEvent
}

func NewSessionManager() *SessionManager

// Get or create a session for an agent role in a project
func (sm *SessionManager) GetOrCreate(config SessionConfig) (*AgentSession, error)

// Get existing session (nil if not found)
func (sm *SessionManager) Get(projectID, agentRole string) *AgentSession

// List all sessions for a project
func (sm *SessionManager) ListByProject(projectID string) []*AgentSession

// List all sessions
func (sm *SessionManager) ListAll() []*AgentSession

// Terminate all sessions for a project
func (sm *SessionManager) TerminateProject(projectID string)

// Terminate all sessions (server shutdown)
func (sm *SessionManager) Shutdown()

// Subscribe to ALL agent events across all sessions (for WebSocket hub)
func (sm *SessionManager) SubscribeAll() <-chan AgentEvent
```

The `eventBus` channel is the central pipeline. When any agent session emits an event, it flows through `eventBus` to the WebSocket hub, which broadcasts it to all connected dashboard clients.

---

## 3. Modify: `internal/agent`

### 3.1 `internal/agent/agent.go` — Major Refactor

**Remove:**
- PTY-based execution (`creack/pty` dependency)
- `claude --print` one-shot invocation
- Response text parsing (regex for DELEGATE:, REVIEW:, FILE:, CREATE_FILE:)
- 10-minute timeout on single CLI call

**Replace with:**
- Session-based execution via `internal/session`
- Structured event parsing (AgentEvent instead of regex)
- Agent holds a reference to its session

**New Agent struct:**

```go
type Agent struct {
    Role            string              // ceo, pm, ux, ui, security, architect, senior_dev, junior_dev
    Name            string              // "CEO Agent", "PM Agent", etc.
    Color           string              // terminal color code
    Permissions     []Permission        // file access permissions (keep existing)
    SystemPrompt    string              // role-specific prompt
    PermissionMode  string              // claude permission mode
    Session         *session.AgentSession  // persistent session (nil until first use)
    SessionManager  *session.SessionManager // reference to create/get sessions
}
```

**New Execute method:**

```go
func (a *Agent) Execute(ctx context.Context, projectID, workDir, task string) (*AgentResponse, error) {
    // 1. Get or create session
    sess, err := a.SessionManager.GetOrCreate(session.SessionConfig{
        AgentRole:      a.Role,
        AgentName:      a.Name,
        SystemPrompt:   a.SystemPrompt,
        WorkDir:        workDir,
        PermissionMode: a.PermissionMode,
        ProjectID:      projectID,
        MaxTurns:       50,
    })
    if err != nil {
        return nil, fmt.Errorf("session create failed: %w", err)
    }
    a.Session = sess

    // 2. Send task and collect events
    response, events, err := sess.SendTask(ctx, task)
    if err != nil {
        return nil, err
    }

    // 3. Parse structured response from events
    return a.parseEvents(response, events), nil
}
```

**New AgentResponse struct (replaces string parsing):**

```go
type AgentResponse struct {
    // Text output from the agent
    Text            string

    // Structured signals extracted from events
    Delegations     []Delegation        // if agent delegated to another
    Reviews         []Review            // if agent reviewed another's work
    FilesCreated    []string            // files created (from tool_use events)
    FilesModified   []string            // files modified (from tool_use events)
    CommandsRun     []CommandRecord     // bash commands executed
    IsComplete      bool                // agent signaled completion

    // Metrics
    InputTokens     int
    OutputTokens    int
    CostUSD         float64
    ToolCalls       int
    Duration        time.Duration

    // Raw events for detailed reporting
    Events          []session.AgentEvent
}

type CommandRecord struct {
    Command string
    Output  string
    IsError bool
}
```

**Event Parsing (replaces regex):**

```go
func (a *Agent) parseEvents(text string, events []session.AgentEvent) *AgentResponse {
    resp := &AgentResponse{Text: text, Events: events}

    for _, ev := range events {
        switch ev.Type {
        case "tool_use":
            switch ev.ToolName {
            case "Write":
                resp.FilesCreated = append(resp.FilesCreated, extractFilePath(ev.Input))
            case "Edit":
                resp.FilesModified = append(resp.FilesModified, extractFilePath(ev.Input))
            case "Bash":
                resp.CommandsRun = append(resp.CommandsRun, CommandRecord{
                    Command: extractCommand(ev.Input),
                })
            }
            resp.ToolCalls++

        case "tool_result":
            // Match with corresponding tool_use to capture output
            updateCommandOutput(resp, ev)

        case "turn_complete":
            resp.InputTokens = ev.InputTokens
            resp.OutputTokens = ev.OutputTokens
            resp.CostUSD = ev.CostUSD
            resp.IsComplete = true
        }
    }

    // Parse text for delegation/review signals (keep lightweight parsing)
    resp.Delegations = parseDelegations(text)
    resp.Reviews = parseReviews(text)

    return resp
}
```

### 3.2 `internal/agent/roles.go` — New File

Define all 8 agent roles with their configurations.

```go
package agent

var AgentRoles = map[string]AgentConfig{
    "ceo": {
        Name:           "CEO Agent",
        Color:          "#FF6B6B",
        PermissionMode: "plan",        // CEO reads and decides, doesn't write code
        SystemPrompt:   ceoPrompt,     // loaded from prompts/ceo.md or inline
    },
    "pm": {
        Name:           "PM Agent",
        Color:          "#4ECDC4",
        PermissionMode: "plan",
        SystemPrompt:   pmPrompt,
    },
    "architect": {
        Name:           "Architect Agent",
        Color:          "#45B7D1",
        PermissionMode: "plan",        // reads codebase, writes specs not code
        SystemPrompt:   architectPrompt,
    },
    "senior_dev": {
        Name:           "Senior Developer",
        Color:          "#96CEB4",
        PermissionMode: "acceptEdits", // can write/edit code
        SystemPrompt:   seniorDevPrompt,
    },
    "junior_dev": {
        Name:           "Junior Developer",
        Color:          "#FFEAA7",
        PermissionMode: "acceptEdits",
        SystemPrompt:   juniorDevPrompt,
    },
    "security": {
        Name:           "Security Expert",
        Color:          "#DDA0DD",
        PermissionMode: "plan",        // audits, doesn't modify
        SystemPrompt:   securityPrompt,
    },
    "ux": {
        Name:           "UX Designer",
        Color:          "#F0E68C",
        PermissionMode: "plan",
        SystemPrompt:   uxPrompt,
    },
    "ui": {
        Name:           "UI Designer",
        Color:          "#87CEEB",
        PermissionMode: "plan",
        SystemPrompt:   uiPrompt,
    },
}
```

### 3.3 `internal/agent/permissions.go` — Simplify

With Claude Code's built-in permission system handling file access, the custom permission layer becomes a secondary safety net. Keep it but simplify — it now validates after-the-fact rather than blocking before execution.

---

## 4. Modify: `internal/orchestrator`

### 4.1 `orchestrator.go` — Session-Aware Orchestration

**Add to Orchestrator struct:**

```go
type Orchestrator struct {
    // ... existing fields ...
    sessionManager  *session.SessionManager
    agents          map[string]*agent.Agent  // role → agent instance
}
```

**Modify NewOrchestrator:**

```go
func NewOrchestrator(/* ... existing params ... */, sm *session.SessionManager) *Orchestrator {
    o := &Orchestrator{
        // ... existing init ...
        sessionManager: sm,
        agents:         initAgents(sm),
    }
    return o
}

func initAgents(sm *session.SessionManager) map[string]*agent.Agent {
    agents := make(map[string]*agent.Agent)
    for role, config := range agent.AgentRoles {
        agents[role] = &agent.Agent{
            Role:           role,
            Name:           config.Name,
            Color:          config.Color,
            SystemPrompt:   config.SystemPrompt,
            PermissionMode: config.PermissionMode,
            SessionManager: sm,
        }
    }
    return agents
}
```

**Modify phase execution to use sessions:**

```go
func (o *Orchestrator) executePhaseAgent(ctx context.Context, projectID, workDir, agentRole, task string) (*agent.AgentResponse, error) {
    ag, ok := o.agents[agentRole]
    if !ok {
        return nil, fmt.Errorf("unknown agent role: %s", agentRole)
    }

    // Execute returns structured response with full event history
    resp, err := ag.Execute(ctx, projectID, workDir, task)
    if err != nil {
        return nil, err
    }

    // Broadcast completion to message store
    o.messageStore.Add(message.Message{
        Type:      "agent_complete",
        AgentRole: agentRole,
        AgentName: ag.Name,
        Content:   resp.Text,
        ProjectID: projectID,
        Metadata: map[string]interface{}{
            "files_created":  resp.FilesCreated,
            "files_modified": resp.FilesModified,
            "tool_calls":     resp.ToolCalls,
            "input_tokens":   resp.InputTokens,
            "output_tokens":  resp.OutputTokens,
            "cost_usd":       resp.CostUSD,
        },
    })

    return resp, nil
}
```

### 4.2 `worker_pool.go` — Session-Aware Parallel Execution

The worker pool currently runs agents in parallel. With sessions, each parallel agent gets its own session, and events stream concurrently.

```go
type AgentTask struct {
    AgentRole string
    Prompt    string
    ProjectID string
    WorkDir   string
}

type AgentResult struct {
    AgentRole string
    Response  *agent.AgentResponse
    Error     error
}

func (o *Orchestrator) runParallel(ctx context.Context, tasks []AgentTask) []AgentResult {
    results := make([]AgentResult, len(tasks))
    var wg sync.WaitGroup

    for i, task := range tasks {
        wg.Add(1)
        go func(idx int, t AgentTask) {
            defer wg.Done()
            resp, err := o.executePhaseAgent(ctx, t.ProjectID, t.WorkDir, t.AgentRole, t.Prompt)
            results[idx] = AgentResult{
                AgentRole: t.AgentRole,
                Response:  resp,
                Error:     err,
            }
        }(i, task)
    }

    wg.Wait()
    return results
}
```

### 4.3 `phases.go` — Updated Phase Definitions

Each phase now specifies which agents participate and what permission mode they need:

```go
var PhaseAgents = map[string][]PhaseAgent{
    "triage": {
        {Role: "ceo", Mode: "plan"},
    },
    "research": {
        {Role: "pm", Mode: "plan"},
        {Role: "ux", Mode: "plan"},
        {Role: "ui", Mode: "plan"},
        {Role: "security", Mode: "plan"},
        {Role: "architect", Mode: "plan"},
    },
    "planning": {
        {Role: "pm", Mode: "plan"},
        {Role: "ux", Mode: "plan"},
        {Role: "ui", Mode: "plan"},
        {Role: "security", Mode: "plan"},
        {Role: "architect", Mode: "plan"},
    },
    "discussion": {
        {Role: "ceo", Mode: "plan"},
    },
    "development": {
        {Role: "senior_dev", Mode: "acceptEdits"},
        {Role: "junior_dev", Mode: "acceptEdits"},
    },
    "qa": {
        {Role: "ceo", Mode: "plan"},
        {Role: "security", Mode: "plan"},
    },
}
```

### 4.4 Context Passing Between Phases

After each phase completes, the orchestrator summarizes results and feeds them to the next phase's agents:

```go
func (o *Orchestrator) buildPhaseContext(projectID string, previousResults []AgentResult) string {
    var ctx strings.Builder
    ctx.WriteString("## Context from previous phase\n\n")

    for _, r := range previousResults {
        if r.Error != nil {
            continue
        }
        ctx.WriteString(fmt.Sprintf("### %s Output\n", r.Response.AgentName))
        ctx.WriteString(r.Response.Text)
        ctx.WriteString("\n\n")

        if len(r.Response.FilesCreated) > 0 {
            ctx.WriteString("Files created: " + strings.Join(r.Response.FilesCreated, ", ") + "\n")
        }
        if len(r.Response.FilesModified) > 0 {
            ctx.WriteString("Files modified: " + strings.Join(r.Response.FilesModified, ", ") + "\n")
        }
        ctx.WriteString("\n")
    }

    return ctx.String()
}
```

---

## 5. Modify: `internal/web`

### 5.1 `server.go` — New API Routes

Add these endpoints to the existing router:

```go
// Agent session management
mux.HandleFunc("GET /api/sessions", handleListAgentSessions(sessionManager))
mux.HandleFunc("GET /api/sessions/{projectID}", handleListProjectSessions(sessionManager))
mux.HandleFunc("GET /api/sessions/{projectID}/{agentRole}", handleGetAgentSession(sessionManager))
mux.HandleFunc("GET /api/sessions/{projectID}/{agentRole}/events", handleGetAgentEvents(sessionManager))
mux.HandleFunc("GET /api/sessions/{projectID}/{agentRole}/stats", handleGetAgentStats(sessionManager))
mux.HandleFunc("POST /api/sessions/{projectID}/{agentRole}/abort", handleAbortAgent(sessionManager))

// Tool approval (human-in-the-loop)
mux.HandleFunc("POST /api/sessions/{projectID}/{agentRole}/approve", handleApproveTool(sessionManager))
mux.HandleFunc("POST /api/sessions/{projectID}/{agentRole}/deny", handleDenyTool(sessionManager))

// Agent activity report
mux.HandleFunc("GET /api/reports/{projectID}", handleProjectReport(sessionManager))
mux.HandleFunc("GET /api/reports/{projectID}/{agentRole}", handleAgentReport(sessionManager))
```

### 5.2 `websocket.go` — Per-Agent Event Streaming

**Extend WebSocket protocol to stream agent events:**

The WebSocket hub subscribes to `SessionManager.SubscribeAll()` and broadcasts every `AgentEvent` to connected dashboard clients.

```go
func (h *Hub) startAgentEventStream(sm *session.SessionManager) {
    eventCh := sm.SubscribeAll()

    go func() {
        for ev := range eventCh {
            // Wrap AgentEvent in WebSocket message
            wsMsg := WSMessage{
                Type:    "agent_event",
                Payload: ev,
            }
            h.broadcast(wsMsg)
        }
    }()
}
```

**New WebSocket message types (server → client):**

```json
{
    "type": "agent_event",
    "payload": {
        "type": "text_delta",
        "agent_role": "senior_dev",
        "agent_name": "Senior Developer",
        "project_id": "proj-123",
        "task_id": "task-456",
        "content": "Let me implement the login endpoint...",
        "timestamp": 1711036800000
    }
}
```

```json
{
    "type": "agent_event",
    "payload": {
        "type": "tool_use",
        "agent_role": "senior_dev",
        "agent_name": "Senior Developer",
        "project_id": "proj-123",
        "tool_use_id": "toolu_abc",
        "tool_name": "Edit",
        "input": "{\"file_path\": \"src/auth.go\", ...}",
        "timestamp": 1711036801000
    }
}
```

```json
{
    "type": "agent_event",
    "payload": {
        "type": "turn_complete",
        "agent_role": "senior_dev",
        "agent_name": "Senior Developer",
        "project_id": "proj-123",
        "input_tokens": 15000,
        "output_tokens": 3200,
        "cost_usd": 0.0456,
        "timestamp": 1711036850000
    }
}
```

**New WebSocket message types (client → server):**

```json
{"type": "approve_agent_tool", "agent_role": "senior_dev", "project_id": "proj-123", "tool_use_id": "toolu_abc"}
{"type": "deny_agent_tool", "agent_role": "senior_dev", "project_id": "proj-123", "tool_use_id": "toolu_abc"}
{"type": "abort_agent", "agent_role": "senior_dev", "project_id": "proj-123"}
```

---

## 6. Dashboard UI: Per-Agent Live Panels

### 6.1 Mission Control Dashboard — Agent Activity View

Add a new view/section to the existing Mission Control dashboard:

**Layout:**

```
┌──────────────────────────────────────────────────────────────────┐
│  MISSION CONTROL                                    [Phase: Dev] │
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌───────────┐ │
│  │  CEO        │ │  PM         │ │  Architect  │ │  Security │ │
│  │  ● Idle     │ │  ● Idle     │ │  ● Idle     │ │  ● Idle   │ │
│  │             │ │             │ │             │ │           │ │
│  │  0 tokens   │ │  0 tokens   │ │  0 tokens   │ │  0 tokens │ │
│  │  $0.00      │ │  $0.00      │ │  $0.00      │ │  $0.00    │ │
│  └─────────────┘ └─────────────┘ └─────────────┘ └───────────┘ │
│                                                                  │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌───────────┐ │
│  │  UX         │ │  UI         │ │  Senior Dev │ │ Junior Dev│ │
│  │  ● Idle     │ │  ● Idle     │ │  ◉ ACTIVE   │ │  ● Idle   │ │
│  │             │ │             │ │             │ │           │ │
│  │  0 tokens   │ │  0 tokens   │ │  12.4K tok  │ │  0 tokens │ │
│  │  $0.00      │ │  $0.00      │ │  $0.034     │ │  $0.00    │ │
│  └─────────────┘ └─────────────┘ └─────────────┘ └───────────┘ │
│                                                                  │
├──────────────────────────────────────────────────────────────────┤
│  SENIOR DEVELOPER — Live Activity                    [Abort ■]  │
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  💭 Thinking: Let me analyze the existing auth middleware...     │
│                                                                  │
│  🔧 Read src/middleware/auth.go                                  │
│     → 45 lines read                                              │
│                                                                  │
│  🔧 Edit src/middleware/auth.go                                  │
│     → Added JWT validation function (lines 23-45)               │
│                                                                  │
│  🔧 Bash: go test ./src/middleware/...                           │
│     → PASS (0.34s)                                               │
│                                                                  │
│  📝 Now implementing the login handler...                        │
│                                                                  │
│  🔧 Write src/handlers/login.go                    [APPROVE ✓]  │
│     → Creating new file (78 lines)                  [DENY ✗]    │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
```

### 6.2 Agent Detail View

Clicking an agent card expands to show full activity log:

```
┌──────────────────────────────────────────────────────────────────┐
│  SENIOR DEVELOPER — Full Report                      [← Back]   │
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Session: abc-123  |  Model: claude-sonnet-4-6  |  Status: Active│
│  Tokens: 15.2K in / 3.4K out  |  Cost: $0.045  |  Tools: 12    │
│                                                                  │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │ TASK 1: Implement JWT auth middleware                      │  │
│  │ Status: Complete ✓  Duration: 2m 34s                       │  │
│  │                                                            │  │
│  │ Timeline:                                                  │  │
│  │ 10:23:01  💭 Analyzing existing auth code...               │  │
│  │ 10:23:03  📖 Read src/middleware/auth.go                   │  │
│  │ 10:23:05  📖 Read go.mod (checking jwt library)            │  │
│  │ 10:23:08  🔧 Bash: go get github.com/golang-jwt/jwt/v5    │  │
│  │ 10:23:15  ✏️  Edit src/middleware/auth.go (+42 lines)       │  │
│  │ 10:23:22  📝 Write src/middleware/auth_test.go (new)       │  │
│  │ 10:23:30  🔧 Bash: go test ./src/middleware/... → PASS     │  │
│  │ 10:23:35  ✅ Turn complete                                 │  │
│  │                                                            │  │
│  │ Files Created: src/middleware/auth_test.go                  │  │
│  │ Files Modified: src/middleware/auth.go, go.mod, go.sum     │  │
│  │ Commands: 2 (go get, go test)                              │  │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │ TASK 2: Implement login handler                            │  │
│  │ Status: In Progress ◉                                      │  │
│  │ ...                                                        │  │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
```

### 6.3 Implementation (Vanilla JS — matches existing pattern)

All dashboard UI in agent-house is vanilla HTML/CSS/JS embedded in Go template strings. The agent activity view follows the same pattern:

- New handler: `internal/web/agent_activity.go` → serves the agent activity page
- JavaScript connects to existing WebSocket hub
- Filters `agent_event` messages by `agent_role`
- Renders events as a timeline with icons per tool type
- Auto-scrolls to latest activity
- Tool approval buttons send WebSocket messages back

---

## 7. Human-in-the-Loop via Tool Approval

The existing checkpoint system pauses at phase gates. With agent sessions, we get **granular tool-level approval** for free.

### 7.1 Permission Mode Strategy

```
┌──────────────────────────────────────────────────────────┐
│ Phase          │ Agents              │ Permission Mode    │
├──────────────────────────────────────────────────────────┤
│ Triage         │ CEO                 │ plan (read-only)   │
│ Research       │ PM,UX,UI,Sec,Arch   │ plan (read-only)   │
│ Planning       │ PM,UX,UI,Sec,Arch   │ plan (read + spec) │
│ Discussion     │ CEO                 │ plan (read-only)   │
│ Development    │ Senior, Junior Dev  │ acceptEdits        │
│ QA             │ CEO, Security       │ plan (read-only)   │
└──────────────────────────────────────────────────────────┘
```

- **plan mode**: Agent can read files and propose edits but must get approval. Perfect for research/planning agents — they can browse the codebase but cannot change anything without human consent.
- **acceptEdits mode**: Agent can read and write files freely but must get approval for bash commands. Good for dev agents — they can code but dangerous commands need human approval.
- **User override**: Dashboard has a toggle per-agent to switch to `default` mode (ask for everything) or `bypassPermissions` (fully autonomous).

### 7.2 Tool Approval Flow

```
Agent Session → tool_use event → Dashboard shows approval card
                                      ↓
                               User clicks Approve/Deny
                                      ↓
                               WebSocket → Server → session.ApproveTool(id) / DenyTool(id)
                                      ↓
                               Written to Claude stdin → Agent continues
```

### 7.3 Auto-Approval with Timeout

For autonomous operation, add configurable auto-approval:

```go
type ApprovalPolicy struct {
    AutoApproveReads    bool          // auto-approve Read, Glob, Grep
    AutoApproveEdits    bool          // auto-approve Edit, Write
    AutoApproveBash     bool          // auto-approve Bash commands
    TimeoutAction       string        // "approve" or "deny" on timeout
    TimeoutDuration     time.Duration // default: 5 minutes
}
```

---

## 8. Cost & Token Tracking

### 8.1 Per-Agent Metrics

Each `AgentSession.Stats` accumulates:
- Total input/output tokens
- Total cost in USD
- Total turns (task invocations)
- Total tool calls with breakdown by tool name

### 8.2 Per-Project Aggregate

```go
type ProjectMetrics struct {
    ProjectID       string                    `json:"project_id"`
    TotalCostUSD    float64                   `json:"total_cost_usd"`
    TotalTokens     int                       `json:"total_tokens"`
    AgentMetrics    map[string]SessionStats   `json:"agent_metrics"` // role → stats
    PhaseMetrics    map[string]PhaseStats     `json:"phase_metrics"` // phase → stats
    StartedAt       time.Time                 `json:"started_at"`
    Duration        time.Duration             `json:"duration"`
}
```

### 8.3 API Endpoints

```
GET /api/reports/{projectID}              → ProjectMetrics
GET /api/reports/{projectID}/{agentRole}  → detailed agent report with event timeline
```

---

## 9. Error Recovery & Safety

### 9.1 Agent Process Crash

```go
func (s *AgentSession) SendTask(ctx context.Context, prompt string) (string, []AgentEvent, error) {
    // ... attempt execution ...

    // If process died, respawn and retry once
    if !s.IsAlive() {
        log.Printf("session %s: process died, respawning", s.ID)
        if err := s.spawn(); err != nil {
            return "", nil, fmt.Errorf("respawn failed: %w", err)
        }
        // Retry the task
        return s.sendTaskInternal(ctx, prompt)
    }
}
```

### 9.2 Runaway Agent Protection

```go
// MaxTurns safety limit (passed to claude --max-turns)
// Default: 50 turns per task
// If agent exceeds, Claude Code auto-terminates the turn

// Additionally, context-level timeout:
ctx, cancel := context.WithTimeout(parentCtx, 10*time.Minute)
defer cancel()
resp, events, err := agent.Execute(ctx, projectID, workDir, task)
```

### 9.3 Cost Ceiling

```go
type CostLimits struct {
    MaxPerAgent   float64 // max USD per agent per task (default: $1.00)
    MaxPerProject float64 // max USD per project (default: $10.00)
}

// Check before each agent task
func (o *Orchestrator) checkCostLimits(projectID, agentRole string) error {
    stats := o.sessionManager.Get(projectID, agentRole).GetStats()
    if stats.TotalCostUSD > o.costLimits.MaxPerAgent {
        return fmt.Errorf("agent %s exceeded cost limit ($%.2f > $%.2f)",
            agentRole, stats.TotalCostUSD, o.costLimits.MaxPerAgent)
    }
    // ... check project total ...
}
```

### 9.4 File Conflict Detection

When multiple dev agents run in the development phase:

```go
// Track files being modified by each agent
type FileTracker struct {
    mu    sync.Mutex
    locks map[string]string // file_path → agent_role currently modifying
}

func (ft *FileTracker) Acquire(path, agentRole string) error {
    ft.mu.Lock()
    defer ft.mu.Unlock()
    if owner, ok := ft.locks[path]; ok && owner != agentRole {
        return fmt.Errorf("file %s is being modified by %s", path, owner)
    }
    ft.locks[path] = agentRole
    return nil
}
```

---

## 10. API Changes

### New Endpoints

| Method | Path | Purpose |
|--------|------|---------|
| GET | `/api/sessions` | List all agent sessions |
| GET | `/api/sessions/{projectID}` | List sessions for a project |
| GET | `/api/sessions/{projectID}/{agentRole}` | Get specific agent session |
| GET | `/api/sessions/{projectID}/{agentRole}/events` | Get agent event history |
| GET | `/api/sessions/{projectID}/{agentRole}/stats` | Get agent metrics |
| POST | `/api/sessions/{projectID}/{agentRole}/abort` | Abort agent's current task |
| POST | `/api/sessions/{projectID}/{agentRole}/approve` | Approve tool request |
| POST | `/api/sessions/{projectID}/{agentRole}/deny` | Deny tool request |
| GET | `/api/reports/{projectID}` | Project-level report |
| GET | `/api/reports/{projectID}/{agentRole}` | Agent-level detailed report |

### Modified Endpoints

| Endpoint | Change |
|----------|--------|
| POST `/api/tasks` | Response now includes session IDs for spawned agents |
| GET `/api/tasks/{id}` | Response includes per-agent metrics and event counts |

---

## 11. WebSocket Protocol Changes

### New Server → Client Messages

```javascript
// Agent streaming events (real-time)
{ type: "agent_event", payload: { type: "text_delta", agent_role, content, ... } }
{ type: "agent_event", payload: { type: "thinking_delta", agent_role, content, ... } }
{ type: "agent_event", payload: { type: "tool_use", agent_role, tool_name, input, ... } }
{ type: "agent_event", payload: { type: "tool_result", agent_role, output, is_error, ... } }
{ type: "agent_event", payload: { type: "turn_complete", agent_role, cost_usd, tokens, ... } }
{ type: "agent_event", payload: { type: "session_meta", agent_role, model, ... } }
{ type: "agent_event", payload: { type: "error", agent_role, content, ... } }

// Agent state changes
{ type: "agent_state", agent_role: "senior_dev", state: "active", task: "Implementing login" }
{ type: "agent_state", agent_role: "senior_dev", state: "idle" }
{ type: "agent_state", agent_role: "senior_dev", state: "waiting_approval", tool_use_id: "..." }
```

### New Client → Server Messages

```javascript
// Tool approval
{ type: "approve_agent_tool", project_id, agent_role, tool_use_id }
{ type: "deny_agent_tool", project_id, agent_role, tool_use_id }

// Agent control
{ type: "abort_agent", project_id, agent_role }

// Request agent event history
{ type: "subscribe_agent", project_id, agent_role }  // start receiving events for this agent
{ type: "unsubscribe_agent", project_id, agent_role }
```

---

## 12. Migration Path

### Phase A: Add Session Package (No Breaking Changes)

1. Create `internal/session/` package
2. Port types, translator, session, manager from mobile-code
3. Add `go.mod` dependency: `nhooyr.io/websocket` (if not already present)
4. Write unit tests for translator and session lifecycle
5. **Nothing breaks** — existing agent code untouched

### Phase B: Dual-Mode Agent Execution

1. Add `ExecuteWithSession()` method alongside existing `Execute()`
2. Feature flag: `USE_SESSIONS=true` enables session-based execution
3. Both modes coexist — can test sessions without breaking the current flow
4. Gradually validate each agent works correctly with sessions

### Phase C: WebSocket Extension

1. Add `agent_event` message type to WebSocket hub
2. Session manager's event bus connects to hub
3. Dashboard gets new JavaScript handler for `agent_event`
4. Build the agent activity panel UI
5. **Existing WebSocket messages unchanged** — additive only

### Phase D: Full Cutover

1. Remove old `claude --print` execution path
2. Remove `creack/pty` dependency
3. Update all orchestrator phases to use session-based execution
4. Enable tool approval flow in dashboard
5. Add cost tracking and reporting endpoints

### Phase E: Polish

1. Agent detail view with full timeline
2. Cost dashboard
3. Auto-approval policies
4. File conflict detection
5. Session persistence (survive server restart)

---

## 13. Step-by-Step Build Order

### Step 1: `internal/session/types.go`
- [ ] Define ClaudeEvent, ClaudeMessage, ClaudeContent structs
- [ ] Define AgentEvent struct (with agent_role, project_id tagging)
- [ ] Define SessionConfig, SessionStats structs
- [ ] No dependencies on existing code

### Step 2: `internal/session/translator.go`
- [ ] Port Translator from mobile-code
- [ ] Modify to emit AgentEvent instead of RelayEvent
- [ ] Tag every event with agent role, project ID, timestamp
- [ ] Unit test: feed sample NDJSON → verify AgentEvent output

### Step 3: `internal/session/session.go`
- [ ] Implement AgentSession struct
- [ ] Implement spawn() — construct Claude CLI command with flags
- [ ] Implement readLoop() — NDJSON parsing with 1MB buffer
- [ ] Implement drainStderr(), watchExit()
- [ ] Implement fanOut() — multi-listener event distribution
- [ ] Implement writeStdin() — JSON + newline to stdin
- [ ] Implement SendTask() — write prompt, collect events until turn_complete
- [ ] Implement SendFollowUp() — same but session retains context
- [ ] Implement ApproveTool(), DenyTool(), Abort()
- [ ] Implement Subscribe() — return event channel + unsubscribe func
- [ ] Implement Terminate() — graceful shutdown (stdin close → wait → kill)
- [ ] Integration test: spawn real Claude session, send "echo hello", verify events

### Step 4: `internal/session/manager.go`
- [ ] Implement SessionManager with map[string]*AgentSession
- [ ] Implement GetOrCreate() — keyed by "{projectID}:{agentRole}"
- [ ] Implement SubscribeAll() — merged event channel from all sessions
- [ ] Implement TerminateProject(), Shutdown()
- [ ] Unit test: create multiple sessions, verify event routing

### Step 5: `internal/agent/roles.go`
- [ ] Define AgentConfig for all 8 roles
- [ ] Set permission modes per role
- [ ] Load system prompts (from existing prompts/ directory or inline)

### Step 6: `internal/agent/agent.go` — Refactor
- [ ] Add Session field to Agent struct
- [ ] Add SessionManager reference
- [ ] Implement new Execute() using session.SendTask()
- [ ] Implement parseEvents() — extract files, commands, delegations from AgentEvents
- [ ] Define AgentResponse struct with structured data
- [ ] Keep old Execute as ExecuteLegacy() temporarily
- [ ] Test: execute one agent with session, verify response

### Step 7: `internal/orchestrator/orchestrator.go` — Integration
- [ ] Add sessionManager field
- [ ] Modify NewOrchestrator to accept SessionManager
- [ ] Modify initAgents to create session-aware agents
- [ ] Modify executePhaseAgent to use new Agent.Execute()
- [ ] Add buildPhaseContext() for inter-agent context passing
- [ ] Test: run triage phase with CEO session

### Step 8: `internal/orchestrator/worker_pool.go` — Update
- [ ] Modify worker pool to use session-based agent execution
- [ ] Test: run research phase with 5 parallel agent sessions

### Step 9: `internal/web/websocket.go` — Extend
- [ ] Add agent_event broadcast from SessionManager.SubscribeAll()
- [ ] Add handlers for approve_agent_tool, deny_agent_tool, abort_agent
- [ ] Add subscribe_agent/unsubscribe_agent for filtered streaming
- [ ] Test: connect WebSocket, submit task, verify agent events stream

### Step 10: `internal/web/server.go` — New Routes
- [ ] Add GET /api/sessions endpoints
- [ ] Add POST abort/approve/deny endpoints
- [ ] Add GET /api/reports endpoints
- [ ] Update api_test.html with new endpoints
- [ ] Update swagger.yaml with new endpoints
- [ ] Update uiflow.md with new flows

### Step 11: `internal/web/agent_activity.go` — Dashboard UI
- [ ] Agent grid view (8 cards with status indicators)
- [ ] Live activity feed per agent (text, thinking, tools)
- [ ] Tool approval cards with Approve/Deny buttons
- [ ] Agent detail view with full event timeline
- [ ] Cost/token display per agent and total
- [ ] CSS: color-coded per agent role, animations for active state

### Step 12: Cleanup
- [ ] Remove creack/pty dependency from go.mod
- [ ] Remove old ExecuteLegacy() code path
- [ ] Remove PTY-related code
- [ ] Remove regex-based response parsing (if no longer needed)
- [ ] Update all tests

### Step 13: Documentation
- [ ] Update api_test.html with all new endpoints
- [ ] Update swagger.yaml with full API spec
- [ ] Update uiflow.md with agent activity flows
- [ ] Update base.md / README with new architecture

---

## 14. Testing Plan

### Unit Tests

| Test | Package | What It Verifies |
|------|---------|-----------------|
| TestTranslateSystemInit | session | system/init → session_meta AgentEvent |
| TestTranslateTextDelta | session | assistant text → text_delta AgentEvent |
| TestTranslateThinking | session | assistant thinking → thinking_delta AgentEvent |
| TestTranslateToolUse | session | assistant tool_use → tool_use AgentEvent with parsed input |
| TestTranslateToolResult | session | user tool_result → tool_result AgentEvent |
| TestTranslateTurnComplete | session | result/success → turn_complete with metrics |
| TestTranslateContentIndex | session | repeated content blocks not re-emitted |
| TestExtractContent | session | polymorphic tool_result content (string vs array) |
| TestSessionState | session | state transitions: active → idle → terminated |
| TestSessionFanOut | session | events reach all subscribers |
| TestManagerGetOrCreate | session | idempotent session creation by key |
| TestManagerTerminateProject | session | all project sessions terminated |
| TestAgentParseEvents | agent | AgentResponse correctly extracts files, commands, metrics |
| TestAgentExecute | agent | end-to-end with mock session |

### Integration Tests

| Test | What It Verifies |
|------|-----------------|
| TestRealSession | Spawn real Claude session, send "write hello.txt", verify file created |
| TestSessionReuse | Same session handles two consecutive tasks, retains context |
| TestParallelSessions | 3 sessions run simultaneously without interference |
| TestSessionCrashRecovery | Kill Claude process, verify respawn and retry |
| TestToolApproval | Send task that triggers tool_use, approve via API, verify completion |
| TestAbort | Send long task, abort mid-execution, verify SIGINT delivered |
| TestWebSocketStream | Connect WS, submit task, verify agent_event messages received |
| TestCostTracking | Run task, verify token/cost metrics in session stats |

### End-to-End Tests

| Test | What It Verifies |
|------|-----------------|
| TestFullPipeline | Submit task, run all 5 phases with agent sessions, verify deliverables |
| TestDashboardLive | Submit task, verify dashboard receives real-time agent events via WS |
| TestToolApprovalUI | Submit task, tool_use triggers, approve via WS, agent continues |

---

## Appendix: Dependencies

### New Dependencies

```
nhooyr.io/websocket v1.8.17    — may already be in go.mod, else add
```

### Removed Dependencies

```
github.com/creack/pty           — no longer needed (was for PTY-based claude --print)
```

### Unchanged Dependencies

```
github.com/google/uuid          — still used for session IDs
All existing dependencies       — no conflicts
```

---

## Appendix: Environment Variables

| Variable | Default | Purpose |
|----------|---------|---------|
| `AUTH_TOKEN` | (required) | Bearer token for API auth |
| `USE_SESSIONS` | `true` | Enable session-based agent execution (Phase B flag) |
| `MAX_AGENT_COST` | `1.00` | Max USD per agent per task |
| `MAX_PROJECT_COST` | `10.00` | Max USD per project |
| `AGENT_TIMEOUT` | `10m` | Max duration per agent task |
| `MAX_TURNS` | `50` | Max Claude turns per agent task |
| `AUTO_APPROVE_READS` | `true` | Auto-approve read-only tool calls |

---

## Appendix: Risk Assessment

| Risk | Impact | Mitigation |
|------|--------|------------|
| 8 concurrent Claude processes uses significant RAM | High | Lazy session creation — only spawn when agent is needed for current phase |
| Agent gets stuck in a loop | Medium | MaxTurns limit + context timeout + abort capability |
| Cost runaway | High | Per-agent and per-project cost ceilings, checked before each task |
| File conflicts between parallel dev agents | Medium | Sequential dev execution or file locking tracker |
| Claude CLI version incompatibility | Low | Pin CLI version, test with CI |
| Session state lost on server restart | Medium | Phase E: add session persistence (save/restore session IDs) |
| NDJSON parse errors | Low | Graceful error handling — log and skip malformed lines |
