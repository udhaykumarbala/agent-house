package orchestrator

import (
	"fmt"
	"sync"
	"time"

	"pty-claude-test/internal/agent"
	"pty-claude-test/internal/message"
)

// Config holds orchestrator configuration
type Config struct {
	ProjectDir   string
	MaxDepth     int  // Maximum delegation depth
	MaxTurns     int  // Maximum total turns
	EnableFiles  bool // Whether to execute file operations
	Verbose      bool // Print detailed logs
}

// DefaultConfig returns sensible defaults
func DefaultConfig() Config {
	return Config{
		ProjectDir:  "./projects/default",
		MaxDepth:    3,  // Reduced from 5 to prevent deep chains
		MaxTurns:    8,  // Reduced from 20 to limit runaway tasks
		EnableFiles: true,
		Verbose:     true,
	}
}

// Orchestrator manages agent collaboration
type Orchestrator struct {
	config    Config
	store     *message.Store
	agents    map[agent.Role]*agent.Agent
	mu        sync.RWMutex
	turns     int
	listeners []func(msg *message.Message)
}

// New creates a new orchestrator
func New(config Config, store *message.Store) *Orchestrator {
	return &Orchestrator{
		config: config,
		store:  store,
		agents: make(map[agent.Role]*agent.Agent),
	}
}

// OnMessage registers a callback for new messages
func (o *Orchestrator) OnMessage(fn func(msg *message.Message)) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.listeners = append(o.listeners, fn)
}

// notify sends message to all listeners
func (o *Orchestrator) notify(msg *message.Message) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	for _, fn := range o.listeners {
		fn(msg)
	}
}

// getAgent returns or creates an agent for the given role
func (o *Orchestrator) getAgent(role agent.Role) (*agent.Agent, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if a, ok := o.agents[role]; ok {
		return a, nil
	}

	a, err := agent.NewAgent(role)
	if err != nil {
		return nil, err
	}

	o.agents[role] = a
	return a, nil
}

// ProcessTask handles a user task through the agent chain
func (o *Orchestrator) ProcessTask(projectID, task string) (*Result, error) {
	o.turns = 0

	result := &Result{
		ProjectID: projectID,
		StartTime: time.Now(),
		Messages:  make([]*message.Message, 0),
		Files:     make([]FileResult, 0),
	}

	// Create initial task message
	taskMsg := message.NewTaskMessage(projectID, task)
	o.store.Add(taskMsg)
	o.notify(taskMsg)
	result.Messages = append(result.Messages, taskMsg)

	if o.config.Verbose {
		fmt.Printf("\n🎯 New Task: %s\n", task)
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	}

	// Start with CEO
	err := o.processAgent(agent.RoleCEO, task, projectID, 0, result)
	if err != nil {
		result.Error = err
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)
	result.TotalTurns = o.turns

	if o.config.Verbose {
		fmt.Printf("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("✅ Task Complete - %d turns, %d files created\n", o.turns, len(result.Files))
		fmt.Printf("⏱️  Duration: %s\n", result.Duration.Round(time.Second))
	}

	return result, err
}

// processAgent handles a single agent's turn
func (o *Orchestrator) processAgent(role agent.Role, task, projectID string, depth int, result *Result) error {
	// Check limits
	if depth > o.config.MaxDepth {
		if o.config.Verbose {
			fmt.Printf("⚠️  Max depth reached (%d), stopping delegation\n", o.config.MaxDepth)
		}
		return nil
	}

	if o.turns >= o.config.MaxTurns {
		if o.config.Verbose {
			fmt.Printf("⚠️  Max turns reached (%d), stopping\n", o.config.MaxTurns)
		}
		return nil
	}

	o.turns++

	// Get agent
	a, err := o.getAgent(role)
	if err != nil {
		return fmt.Errorf("failed to get agent %s: %w", role, err)
	}

	if o.config.Verbose {
		fmt.Printf("\n%s %s is thinking...\n", getAgentEmoji(role), a.Name)
	}

	// Process task
	var response *agent.Response
	var fileResults []agent.FileResult

	if o.config.EnableFiles {
		response, fileResults, err = a.ProcessWithFileOps(task, o.config.ProjectDir)
	} else {
		response, err = a.Process(task)
	}

	if err != nil {
		return fmt.Errorf("agent %s failed: %w", role, err)
	}

	// Record response
	responseMsg := message.NewMessage(message.TypeResponse, string(role), "orchestrator", response.Content)
	responseMsg.Metadata.ProjectID = projectID
	o.store.Add(responseMsg)
	o.notify(responseMsg)
	result.Messages = append(result.Messages, responseMsg)

	if o.config.Verbose {
		// Print truncated response
		content := response.Content
		if len(content) > 300 {
			content = content[:300] + "..."
		}
		fmt.Printf("   └─ %s\n", content)
	}

	// Handle file operations
	for _, fr := range fileResults {
		fileMsg := message.NewMessage(message.TypeFileCreate, string(role), "filesystem", fr.Path)
		fileMsg.Metadata.ProjectID = projectID
		fileMsg.Metadata.FilePath = fr.Path

		if fr.Success {
			if o.config.Verbose {
				fmt.Printf("   📄 Created: %s\n", fr.Path)
			}
		} else {
			fileMsg.Type = message.TypeError
			if o.config.Verbose {
				fmt.Printf("   ❌ Failed: %s - %v\n", fr.Path, fr.Error)
			}
		}

		o.store.Add(fileMsg)
		o.notify(fileMsg)

		result.Files = append(result.Files, FileResult{
			Path:    fr.Path,
			Success: fr.Success,
			Error:   fr.Error,
			Agent:   string(role),
		})
	}

	// Handle delegations
	if len(response.DelegateTo) > 0 {
		for _, delegateRole := range response.DelegateTo {
			// Create delegation message
			delegateMsg := message.NewMessage(
				message.TypeDelegate,
				string(role),
				string(delegateRole),
				fmt.Sprintf("Delegating from %s", a.Name),
			)
			delegateMsg.Metadata.ProjectID = projectID
			delegateMsg.Metadata.DelegateTo = []string{string(delegateRole)}
			o.store.Add(delegateMsg)
			o.notify(delegateMsg)
			result.Messages = append(result.Messages, delegateMsg)

			if o.config.Verbose {
				fmt.Printf("   → Delegating to %s\n", delegateRole)
			}

			// Build context for delegated agent
			context := buildDelegationContext(a, response, task)

			// Recursively process
			err := o.processAgent(delegateRole, context, projectID, depth+1, result)
			if err != nil {
				// Log error but continue with other delegations
				errorMsg := message.NewMessage(message.TypeError, string(delegateRole), "orchestrator", err.Error())
				errorMsg.Metadata.ProjectID = projectID
				o.store.Add(errorMsg)
				o.notify(errorMsg)
			}
		}
	}

	return nil
}

// buildDelegationContext creates context for a delegated agent
func buildDelegationContext(from *agent.Agent, response *agent.Response, originalTask string) string {
	return fmt.Sprintf(`You are receiving a delegated task from %s.

ORIGINAL TASK:
%s

MESSAGE FROM %s:
%s

Please complete your part of this task according to your role and expertise.
If you need to delegate subtasks to other agents, use the DELEGATE format.
If you create files, use the appropriate file creation format.`,
		from.Name,
		originalTask,
		from.Name,
		response.Content,
	)
}

// getAgentEmoji returns an emoji for the agent role
func getAgentEmoji(role agent.Role) string {
	emojis := map[agent.Role]string{
		agent.RoleCEO:       "👔",
		agent.RolePM:        "📋",
		agent.RoleUX:        "🎨",
		agent.RoleUI:        "🖼️",
		agent.RoleSecurity:  "🔒",
		agent.RoleArchitect: "🏗️",
		agent.RoleSeniorDev: "👨‍💻",
		agent.RoleJuniorDev: "👩‍💻",
	}
	if e, ok := emojis[role]; ok {
		return e
	}
	return "🤖"
}

// Result holds the outcome of a task
type Result struct {
	ProjectID  string
	StartTime  time.Time
	EndTime    time.Time
	Duration   time.Duration
	TotalTurns int
	Messages   []*message.Message
	Files      []FileResult
	Error      error
}

// FileResult tracks a file operation result
type FileResult struct {
	Path    string
	Success bool
	Error   error
	Agent   string
}

// GetStore returns the message store
func (o *Orchestrator) GetStore() *message.Store {
	return o.store
}

// GetAgents returns all active agents
func (o *Orchestrator) GetAgents() map[agent.Role]*agent.Agent {
	o.mu.RLock()
	defer o.mu.RUnlock()

	// Return a copy
	agents := make(map[agent.Role]*agent.Agent)
	for k, v := range o.agents {
		agents[k] = v
	}
	return agents
}
