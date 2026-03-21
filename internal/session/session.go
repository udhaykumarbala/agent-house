package session

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
)

// AgentSession manages a persistent Claude Code CLI process for one agent.
type AgentSession struct {
	ID              string
	Config          SessionConfig
	ClaudeSessionID string
	State           string
	Stats           SessionStats
	CreatedAt       time.Time

	cmd        *exec.Cmd
	stdin      io.WriteCloser
	stdout     io.ReadCloser
	stderr     io.ReadCloser
	translator *Translator

	mu        sync.Mutex
	listeners []chan AgentEvent
	done      chan struct{} // closed when process exits
	spawned   bool
}

// NewAgentSession creates a new agent session (does not spawn until first use).
func NewAgentSession(config SessionConfig) *AgentSession {
	id := uuid.New().String()
	return &AgentSession{
		ID:        id,
		Config:    config,
		State:     StateIdle,
		CreatedAt: time.Now(),
		Stats: SessionStats{
			ToolBreakdown: make(map[string]int),
		},
		translator: NewTranslator(config.AgentRole, config.AgentName, config.ProjectID),
		listeners:  make([]chan AgentEvent, 0),
		done:       make(chan struct{}),
	}
}

// spawn starts the Claude CLI process.
func (s *AgentSession) spawn() error {
	s.ClaudeSessionID = uuid.New().String()

	args := []string{
		"-p",
		"--input-format", "stream-json",
		"--output-format", "stream-json",
		"--verbose",
	}

	if s.Config.PermissionMode != "" {
		args = append(args, "--permission-mode", s.Config.PermissionMode)
	} else {
		args = append(args, "--permission-mode", "bypassPermissions")
	}

	if s.Config.SystemPrompt != "" {
		args = append(args, "--append-system-prompt", s.Config.SystemPrompt)
	}

	if s.Config.Model != "" {
		args = append(args, "--model", s.Config.Model)
	}

	maxTurns := s.Config.MaxTurns
	if maxTurns <= 0 {
		maxTurns = 50
	}
	args = append(args, "--max-turns", strconv.Itoa(maxTurns))

	args = append(args, "--session-id", s.ClaudeSessionID)

	cmd := exec.Command("claude", args...)
	cmd.Dir = s.Config.WorkDir

	// Filter CLAUDECODE env to prevent nesting issues
	cmd.Env = filterEnv(os.Environ(), "CLAUDECODE")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("stderr pipe: %w", err)
	}

	s.cmd = cmd
	s.stdin = stdin
	s.stdout = stdout
	s.stderr = stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start claude: %w", err)
	}

	s.spawned = true
	s.State = StateActive

	log.Printf("[SESSION:%s:%s] spawned claude (pid=%d, session=%s, dir=%s)",
		s.Config.AgentRole, s.ID[:8], cmd.Process.Pid, s.ClaudeSessionID[:8], s.Config.WorkDir)

	go s.readLoop()
	go s.drainStderr()
	go s.watchExit()

	return nil
}

// DebugNDJSON enables raw NDJSON line logging for all sessions.
// Set via environment variable AGENT_DEBUG_NDJSON=1
var DebugNDJSON = os.Getenv("AGENT_DEBUG_NDJSON") == "1"

// readLoop reads NDJSON lines from Claude's stdout and translates them.
func (s *AgentSession) readLoop() {
	scanner := bufio.NewScanner(s.stdout)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024) // 1MB buffer

	lineNum := 0
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		lineNum++

		if DebugNDJSON {
			preview := string(line)
			if len(preview) > 300 {
				preview = preview[:300] + "..."
			}
			log.Printf("[NDJSON:%s:#%d] %s", s.Config.AgentRole, lineNum, preview)
		}

		var raw ClaudeEvent
		if err := json.Unmarshal(line, &raw); err != nil {
			log.Printf("[SESSION:%s] invalid NDJSON line #%d: %v", s.Config.AgentRole, lineNum, err)
			continue
		}

		if DebugNDJSON {
			log.Printf("[NDJSON:%s:#%d] parsed: type=%s subtype=%s", s.Config.AgentRole, lineNum, raw.Type, raw.Subtype)
			if raw.Message != nil {
				var msg ClaudeMessage
				if json.Unmarshal(raw.Message, &msg) == nil {
					for i, c := range msg.Content {
						log.Printf("[NDJSON:%s:#%d]   content[%d]: type=%s name=%s id=%s", s.Config.AgentRole, lineNum, i, c.Type, c.Name, c.ID)
					}
				}
			}
		}

		// Capture session ID from init event
		if raw.Type == "system" && raw.Subtype == "init" && raw.SessionID != "" {
			s.mu.Lock()
			s.ClaudeSessionID = raw.SessionID
			s.mu.Unlock()
		}

		events := s.translator.Translate(raw)

		if DebugNDJSON && len(events) > 0 {
			for _, ev := range events {
				log.Printf("[NDJSON:%s:#%d] → emit: type=%s tool=%s content_len=%d",
					s.Config.AgentRole, lineNum, ev.Type, ev.ToolName, len(ev.Content)+len(ev.Output))
			}
		}

		for _, ev := range events {
			// Update stats
			s.mu.Lock()
			switch ev.Type {
			case "tool_use":
				s.Stats.TotalToolCalls++
				s.Stats.ToolBreakdown[ev.ToolName]++
			case "turn_complete":
				s.Stats.TotalTurns++
				s.Stats.TotalInputTokens += ev.InputTokens
				s.Stats.TotalOutputTokens += ev.OutputTokens
				s.Stats.TotalCostUSD += ev.CostUSD
			}
			s.mu.Unlock()

			s.fanOut(ev)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("[SESSION:%s] stdout scanner error: %v", s.Config.AgentRole, err)
	}
	if DebugNDJSON {
		log.Printf("[NDJSON:%s] readLoop ended after %d lines", s.Config.AgentRole, lineNum)
	}
}

// drainStderr reads and logs stderr to prevent pipe blocking.
func (s *AgentSession) drainStderr() {
	scanner := bufio.NewScanner(s.stderr)
	for scanner.Scan() {
		log.Printf("[SESSION:%s:stderr] %s", s.Config.AgentRole, scanner.Text())
	}
}

// watchExit waits for the process to exit and updates state.
func (s *AgentSession) watchExit() {
	err := s.cmd.Wait()
	close(s.done)

	s.mu.Lock()
	defer s.mu.Unlock()

	if err != nil {
		log.Printf("[SESSION:%s] claude exited with error: %v", s.Config.AgentRole, err)
	} else {
		log.Printf("[SESSION:%s] claude exited cleanly", s.Config.AgentRole)
	}

	if s.State != StateTerminated {
		s.State = StateIdle
	}
}

// fanOut sends an event to all subscribed listeners.
func (s *AgentSession) fanOut(ev AgentEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, ch := range s.listeners {
		select {
		case ch <- ev:
		default:
			// Slow listener, drop event
		}
	}
}

// Subscribe returns a channel that receives all events from this session.
// Call the returned function to unsubscribe.
func (s *AgentSession) Subscribe() (<-chan AgentEvent, func()) {
	ch := make(chan AgentEvent, 512)

	s.mu.Lock()
	s.listeners = append(s.listeners, ch)
	s.mu.Unlock()

	unsubscribe := func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		for i, l := range s.listeners {
			if l == ch {
				s.listeners = append(s.listeners[:i], s.listeners[i+1:]...)
				break
			}
		}
	}

	return ch, unsubscribe
}

// writeStdin marshals msg to JSON, appends newline, and writes to stdin.
func (s *AgentSession) writeStdin(msg interface{}) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	data = append(data, '\n')
	if _, err := s.stdin.Write(data); err != nil {
		return fmt.Errorf("write stdin: %w", err)
	}
	return nil
}

// SendTask sends a task to the agent and blocks until turn completes.
// Returns the text response, all events, and any error.
func (s *AgentSession) SendTask(ctx context.Context, prompt string) (string, []AgentEvent, error) {
	s.mu.Lock()
	needsSpawn := !s.spawned || s.State == StateIdle
	s.mu.Unlock()

	// Spawn or respawn if needed
	if needsSpawn {
		if s.spawned {
			// Process died, need to respawn
			s.done = make(chan struct{})
			s.translator.Reset()
		}
		if err := s.spawn(); err != nil {
			return "", nil, fmt.Errorf("spawn failed: %w", err)
		}
		// Give Claude a moment to initialize
		time.Sleep(500 * time.Millisecond)
	}

	// Update task context in translator
	s.translator.SetTaskID("")

	// Subscribe to events for this turn
	evCh, unsub := s.Subscribe()
	defer unsub()

	// Write user message to stdin
	msg := map[string]interface{}{
		"type": "user",
		"message": map[string]string{
			"role":    "user",
			"content": prompt,
		},
	}

	s.mu.Lock()
	if s.State != StateActive {
		s.mu.Unlock()
		return "", nil, fmt.Errorf("session not active (state=%s)", s.State)
	}
	s.mu.Unlock()

	if err := s.writeStdin(msg); err != nil {
		return "", nil, fmt.Errorf("send message: %w", err)
	}

	log.Printf("[SESSION:%s] sent task (%d bytes)", s.Config.AgentRole, len(prompt))

	// Collect events until turn_complete or context cancellation
	var events []AgentEvent
	var response strings.Builder

	for {
		select {
		case ev := <-evCh:
			events = append(events, ev)

			switch ev.Type {
			case "text_delta":
				response.WriteString(ev.Content)
			case "turn_complete":
				return response.String(), events, nil
			case "error":
				if ev.Content != "" {
					return response.String(), events, fmt.Errorf("agent error: %s", ev.Content)
				}
				// Some errors are non-fatal (like process exit after turn_complete)
				return response.String(), events, nil
			}

		case <-ctx.Done():
			s.Abort()
			return response.String(), events, ctx.Err()

		case <-s.done:
			// Process exited - return what we have
			if response.Len() > 0 {
				return response.String(), events, nil
			}
			return "", events, fmt.Errorf("agent process exited unexpectedly")
		}
	}
}

// SendFollowUp sends a follow-up message (session retains conversation context).
func (s *AgentSession) SendFollowUp(ctx context.Context, prompt string) (string, []AgentEvent, error) {
	return s.SendTask(ctx, prompt)
}

// ApproveTool approves a pending tool use request.
func (s *AgentSession) ApproveTool(toolUseID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.State != StateActive {
		return fmt.Errorf("session not active (state=%s)", s.State)
	}

	msg := map[string]interface{}{
		"type":        "approve",
		"tool_use_id": toolUseID,
	}
	return s.writeStdin(msg)
}

// DenyTool denies a pending tool use request.
func (s *AgentSession) DenyTool(toolUseID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.State != StateActive {
		return fmt.Errorf("session not active (state=%s)", s.State)
	}

	msg := map[string]interface{}{
		"type":        "deny",
		"tool_use_id": toolUseID,
	}
	return s.writeStdin(msg)
}

// Abort sends SIGINT to the Claude process for mid-turn abort.
func (s *AgentSession) Abort() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cmd == nil || s.cmd.Process == nil {
		return fmt.Errorf("no process")
	}

	if err := s.cmd.Process.Signal(syscall.SIGINT); err != nil {
		return fmt.Errorf("sigint: %w", err)
	}

	log.Printf("[SESSION:%s] sent SIGINT (abort)", s.Config.AgentRole)
	return nil
}

// IsAlive returns whether the Claude process is still running.
func (s *AgentSession) IsAlive() bool {
	select {
	case <-s.done:
		return false
	default:
		return s.spawned
	}
}

// GetStats returns a copy of the accumulated session stats.
func (s *AgentSession) GetStats() SessionStats {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Copy the map
	breakdown := make(map[string]int, len(s.Stats.ToolBreakdown))
	for k, v := range s.Stats.ToolBreakdown {
		breakdown[k] = v
	}

	return SessionStats{
		TotalInputTokens:  s.Stats.TotalInputTokens,
		TotalOutputTokens: s.Stats.TotalOutputTokens,
		TotalCostUSD:      s.Stats.TotalCostUSD,
		TotalTurns:        s.Stats.TotalTurns,
		TotalToolCalls:    s.Stats.TotalToolCalls,
		ToolBreakdown:     breakdown,
	}
}

// GetInfo returns session metadata.
func (s *AgentSession) GetInfo() SessionInfo {
	s.mu.Lock()
	defer s.mu.Unlock()

	return SessionInfo{
		ID:              s.ID,
		AgentRole:       s.Config.AgentRole,
		AgentName:       s.Config.AgentName,
		ProjectID:       s.Config.ProjectID,
		State:           s.State,
		ClaudeSessionID: s.ClaudeSessionID,
		Stats:           s.Stats,
		CreatedAt:       s.CreatedAt.Format(time.RFC3339),
	}
}

// Terminate gracefully terminates the Claude process.
func (s *AgentSession) Terminate() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.State == StateTerminated {
		return nil
	}

	s.State = StateTerminated

	// Close stdin first for clean exit
	if s.stdin != nil {
		s.stdin.Close()
	}

	// Give Claude a moment to exit, then kill
	select {
	case <-s.done:
		// Already exited
	case <-time.After(3 * time.Second):
		if s.cmd != nil && s.cmd.Process != nil {
			s.cmd.Process.Kill()
			log.Printf("[SESSION:%s] force killed claude", s.Config.AgentRole)
		}
	}

	log.Printf("[SESSION:%s] terminated", s.Config.AgentRole)
	return nil
}

// filterEnv returns a copy of env with the specified key removed.
func filterEnv(env []string, key string) []string {
	prefix := key + "="
	filtered := make([]string, 0, len(env))
	for _, e := range env {
		if !strings.HasPrefix(e, prefix) {
			filtered = append(filtered, e)
		}
	}
	return filtered
}
