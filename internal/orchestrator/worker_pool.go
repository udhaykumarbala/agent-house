package orchestrator

import (
	"context"
	"fmt"
	"sync"
	"time"

	"pty-claude-test/internal/agent"
	"pty-claude-test/internal/message"
)

// AgentJob represents a unit of work for an agent
type AgentJob struct {
	Role        agent.Role
	Prompt      string
	Context     map[string]interface{}
	TaskID      string
	ProjectID   string
}

// AgentResult represents the outcome of an agent's work
type AgentResult struct {
	Role            agent.Role
	Messages        []*message.Message
	Files           []string
	Error           error
	Duration        time.Duration
	TaskID          string
	SessionResponse *agent.SessionResponse // populated when using session mode
}

// AgentGetter is a function that returns an agent for a given role
type AgentGetter func(role agent.Role) (*agent.Agent, error)

// AgentWorkerPool manages concurrent agent execution
type AgentWorkerPool struct {
	jobs        chan AgentJob
	results     chan AgentResult
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	numWorkers  int
	getAgent    AgentGetter
	mu          sync.RWMutex
	active      map[agent.Role]bool // Track active agents
	closeOnce   sync.Once           // Ensures results channel is closed exactly once
}

// NewAgentWorkerPool creates a new worker pool for parallel agent execution
func NewAgentWorkerPool(ctx context.Context, numWorkers int, getAgent AgentGetter) *AgentWorkerPool {
	poolCtx, cancel := context.WithCancel(ctx)

	pool := &AgentWorkerPool{
		jobs:        make(chan AgentJob, numWorkers),
		results:     make(chan AgentResult, numWorkers),
		ctx:         poolCtx,
		cancel:      cancel,
		numWorkers:  numWorkers,
		getAgent:    getAgent,
		active:      make(map[agent.Role]bool),
	}

	// Start workers
	for i := 0; i < numWorkers; i++ {
		pool.wg.Add(1)
		go pool.worker(i)
	}

	return pool
}

// worker is the goroutine that processes jobs
func (p *AgentWorkerPool) worker(id int) {
	defer p.wg.Done()

	for {
		select {
		case <-p.ctx.Done():
			return
		case job, ok := <-p.jobs:
			if !ok {
				return // Channel closed
			}

			// Mark agent as active
			p.mu.Lock()
			p.active[job.Role] = true
			p.mu.Unlock()

			// Execute the job
			result := p.executeJob(job)

			// Mark agent as inactive
			p.mu.Lock()
			delete(p.active, job.Role)
			p.mu.Unlock()

			// Send result
			select {
			case <-p.ctx.Done():
				return
			case p.results <- result:
			}
		}
	}
}

// executeJob runs the agent and returns the result
func (p *AgentWorkerPool) executeJob(job AgentJob) AgentResult {
	startTime := time.Now()

	result := AgentResult{
		Role:   job.Role,
		TaskID: job.TaskID,
	}

	// Get or create agent
	agentInstance, err := p.getAgent(job.Role)
	if err != nil {
		result.Error = fmt.Errorf("failed to get agent %s: %w", job.Role, err)
		result.Duration = time.Since(startTime)
		return result
	}

	// Create context with timeout
	jobCtx, cancel := context.WithTimeout(p.ctx, 5*time.Minute)
	defer cancel()

	// Execute agent with context
	messages, files, err, sessResp := p.executeWithContext(jobCtx, agentInstance, job)

	result.Messages = messages
	result.Files = files
	result.Error = err
	result.Duration = time.Since(startTime)
	result.SessionResponse = sessResp

	return result
}

// executeWithContext runs the agent with context support.
// Uses session-based execution if the agent has a SessionManager, otherwise falls back to legacy.
func (p *AgentWorkerPool) executeWithContext(ctx context.Context, agentInstance *agent.Agent, job AgentJob) ([]*message.Message, []string, error, *agent.SessionResponse) {
	// Extract workDir from context
	workDir := ""
	if wd, ok := job.Context["workDir"]; ok {
		if wdStr, ok := wd.(string); ok {
			workDir = wdStr
		}
	}

	// Try session-based execution first
	if agentInstance.SessionManager != nil {
		sessResp, err := agentInstance.ExecuteWithSession(ctx, job.ProjectID, workDir, job.Prompt)

		var messages []*message.Message
		var files []string

		if sessResp != nil {
			msg := message.NewMessage(
				message.TypeResponse,
				string(job.Role),
				"orchestrator",
				sessResp.Text,
			)
			msg.Metadata.ProjectID = job.ProjectID
			msg.Metadata.TaskID = job.TaskID
			messages = append(messages, msg)

			files = append(files, sessResp.FilesCreated...)
			files = append(files, sessResp.FilesModified...)
		}

		return messages, files, err, sessResp
	}

	// Legacy execution path (claude --print)
	type agentOutput struct {
		messages []*message.Message
		files    []string
		err      error
	}
	resultChan := make(chan agentOutput, 1)

	go func() {
		response, err := agentInstance.ProcessInDir(job.Prompt, workDir)

		var messages []*message.Message
		var files []string

		if response != nil {
			msg := message.NewMessage(
				message.TypeResponse,
				string(job.Role),
				"orchestrator",
				response.Content,
			)
			msg.Metadata.ProjectID = job.ProjectID
			msg.Metadata.TaskID = job.TaskID
			messages = append(messages, msg)
			files = response.Files
		}

		resultChan <- agentOutput{messages, files, err}
	}()

	select {
	case <-ctx.Done():
		return nil, nil, fmt.Errorf("agent %s timed out or cancelled: %w", job.Role, ctx.Err()), nil
	case output := <-resultChan:
		return output.messages, output.files, output.err, nil
	}
}

// Submit adds a job to the worker pool
func (p *AgentWorkerPool) Submit(job AgentJob) error {
	select {
	case <-p.ctx.Done():
		return fmt.Errorf("worker pool is shutting down")
	case p.jobs <- job:
		return nil
	}
}

// Results returns the results channel
func (p *AgentWorkerPool) Results() <-chan AgentResult {
	return p.results
}

// Wait closes the jobs channel and waits for all workers to finish
func (p *AgentWorkerPool) Wait() {
	close(p.jobs)
	p.wg.Wait()
	p.closeResults()
}

// Shutdown cancels the context and waits for workers
func (p *AgentWorkerPool) Shutdown() {
	p.cancel()
	p.wg.Wait()
	p.closeResults()
}

// closeResults safely closes the results channel exactly once
func (p *AgentWorkerPool) closeResults() {
	p.closeOnce.Do(func() {
		close(p.results)
	})
}

// IsActive returns whether a specific agent is currently working
func (p *AgentWorkerPool) IsActive(role agent.Role) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.active[role]
}

// ActiveAgents returns a list of currently active agents
func (p *AgentWorkerPool) ActiveAgents() []agent.Role {
	p.mu.RLock()
	defer p.mu.RUnlock()

	roles := make([]agent.Role, 0, len(p.active))
	for role := range p.active {
		roles = append(roles, role)
	}
	return roles
}
