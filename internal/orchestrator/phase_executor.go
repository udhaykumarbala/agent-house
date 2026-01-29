package orchestrator

import (
	"context"
	"fmt"
	"sync"
	"time"

	"pty-claude-test/internal/agent"
	"pty-claude-test/internal/message"
)

// ExecutorConfig holds configuration for the phase executor
type ExecutorConfig struct {
	EnableParallel bool          // Enable parallel agent execution
	MaxWorkers     int           // Max concurrent agents per phase (0 = unlimited)
	AgentTimeout   time.Duration // Timeout per agent
	PhaseTimeout   time.Duration // Timeout per phase
	WorkDir        string        // Working directory for agents
	ProjectID      string        // Project ID
	TaskID         string        // Task ID
}

// DefaultExecutorConfig returns sensible defaults
func DefaultExecutorConfig() ExecutorConfig {
	return ExecutorConfig{
		EnableParallel: true,
		MaxWorkers:     0, // Unlimited
		AgentTimeout:   5 * time.Minute,
		PhaseTimeout:   15 * time.Minute,
	}
}

// PhaseExecutor manages parallel execution of agents within a phase
type PhaseExecutor struct {
	config       ExecutorConfig
	getAgent     AgentGetter
	fileManager  *ConcurrentFileManager
	messageStore *message.Store
	mu           sync.Mutex
}

// NewPhaseExecutor creates a new phase executor
func NewPhaseExecutor(config ExecutorConfig, getAgent AgentGetter, fileManager *ConcurrentFileManager, messageStore *message.Store) *PhaseExecutor {
	return &PhaseExecutor{
		config:       config,
		getAgent:     getAgent,
		fileManager:  fileManager,
		messageStore: messageStore,
	}
}

// ExecuteAgentsParallel executes multiple agents in parallel and returns aggregated results
func (pe *PhaseExecutor) ExecuteAgentsParallel(roles []agent.Role, taskPrompt string) (*PhaseResult, error) {
	result := &PhaseResult{
		Messages:     make([]*message.Message, 0),
		Files:        make([]string, 0),
		Errors:       NewErrorCollector(),
		StartTime:    time.Now(),
	}

	// If parallel execution is disabled, fall back to sequential
	if !pe.config.EnableParallel {
		return pe.executeSequential(roles, taskPrompt)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), pe.config.PhaseTimeout)
	defer cancel()

	// Determine number of workers
	numWorkers := len(roles)
	if pe.config.MaxWorkers > 0 && numWorkers > pe.config.MaxWorkers {
		numWorkers = pe.config.MaxWorkers
	}

	// Create worker pool
	pool := NewAgentWorkerPool(ctx, numWorkers, pe.getAgent)

	// Submit jobs for all agents
	for _, role := range roles {
		job := AgentJob{
			Role:      role,
			Prompt:    taskPrompt,
			Context: map[string]interface{}{
				"workDir": pe.config.WorkDir,
			},
			TaskID:    pe.config.TaskID,
			ProjectID: pe.config.ProjectID,
		}

		if err := pool.Submit(job); err != nil {
			result.Errors.Add(role, fmt.Errorf("failed to submit job: %w", err))
		}
	}

	// Collect results
	go pool.Wait()

	resultsReceived := 0
	expectedResults := len(roles)

	for agentResult := range pool.Results() {
		resultsReceived++

		// Handle errors
		if agentResult.Error != nil {
			result.Errors.Add(agentResult.Role, agentResult.Error)
		}

		// Aggregate messages
		pe.mu.Lock()
		result.Messages = append(result.Messages, agentResult.Messages...)
		pe.mu.Unlock()

		// Store messages in message store
		for _, msg := range agentResult.Messages {
			if pe.messageStore != nil {
				pe.messageStore.Add(msg)
			}
		}

		// Aggregate files
		pe.mu.Lock()
		result.Files = append(result.Files, agentResult.Files...)
		pe.mu.Unlock()

		// Update duration stats
		if agentResult.Duration > result.MaxAgentDuration {
			result.MaxAgentDuration = agentResult.Duration
		}

		// Check if we've received all results
		if resultsReceived >= expectedResults {
			break
		}
	}

	result.Duration = time.Since(result.StartTime)
	result.Success = !result.Errors.HasErrors()

	return result, nil
}

// executeSequential executes agents one by one (fallback)
func (pe *PhaseExecutor) executeSequential(roles []agent.Role, taskPrompt string) (*PhaseResult, error) {
	result := &PhaseResult{
		Messages:  make([]*message.Message, 0),
		Files:     make([]string, 0),
		Errors:    NewErrorCollector(),
		StartTime: time.Now(),
	}

	for _, role := range roles {
		// Get agent
		agentInstance, err := pe.getAgent(role)
		if err != nil {
			result.Errors.Add(role, fmt.Errorf("failed to get agent: %w", err))
			continue
		}

		// Execute agent
		startTime := time.Now()
		response, err := agentInstance.ProcessInDir(taskPrompt, pe.config.WorkDir)
		duration := time.Since(startTime)

		if err != nil {
			result.Errors.Add(role, err)
			continue
		}

		// Create message
		if response != nil {
			msg := message.NewMessage(
				message.TypeResponse,
				string(role),
				"orchestrator",
				response.Content,
			)
			msg.Metadata.ProjectID = pe.config.ProjectID
			msg.Metadata.TaskID = pe.config.TaskID

			result.Messages = append(result.Messages, msg)
			result.Files = append(result.Files, response.Files...)

			// Store message
			if pe.messageStore != nil {
				pe.messageStore.Add(msg)
			}
		}

		// Update duration stats
		if duration > result.MaxAgentDuration {
			result.MaxAgentDuration = duration
		}
	}

	result.Duration = time.Since(result.StartTime)
	result.Success = !result.Errors.HasErrors()

	return result, nil
}

// PhaseResult holds the aggregated results from a phase execution
type PhaseResult struct {
	Messages          []*message.Message
	Files             []string
	Errors            *ErrorCollector
	Duration          time.Duration
	MaxAgentDuration  time.Duration // Longest individual agent execution
	StartTime         time.Time
	Success           bool
}

// HasErrors returns true if any agent encountered errors
func (pr *PhaseResult) HasErrors() bool {
	return pr.Errors.HasErrors()
}

// ErrorCount returns the number of errors
func (pr *PhaseResult) ErrorCount() int {
	return pr.Errors.Count()
}
