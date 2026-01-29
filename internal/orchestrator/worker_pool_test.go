package orchestrator

import (
	"context"
	"errors"
	"testing"
	"time"

	"pty-claude-test/internal/agent"
)

// mockAgentGetter returns a mock agent that simulates work
func mockAgentGetter(simulateError bool, workDuration time.Duration) AgentGetter {
	return func(role agent.Role) (*agent.Agent, error) {
		if simulateError {
			return nil, errors.New("mock agent creation error")
		}
		// Return a minimal agent - actual processing will be mocked
		return &agent.Agent{
			ID:   string(role),
			Role: role,
			Name: string(role),
		}, nil
	}
}

func TestWorkerPoolCreation(t *testing.T) {
	ctx := context.Background()
	pool := NewAgentWorkerPool(ctx, 3, mockAgentGetter(false, 100*time.Millisecond))

	if pool == nil {
		t.Fatal("Pool should not be nil")
	}

	// Test active tracking
	if pool.IsActive(agent.RolePM) {
		t.Error("No agents should be active initially")
	}

	activeAgents := pool.ActiveAgents()
	if len(activeAgents) != 0 {
		t.Errorf("Expected 0 active agents, got %d", len(activeAgents))
	}

	pool.Shutdown()
}

func TestWorkerPoolSubmit(t *testing.T) {
	ctx := context.Background()
	pool := NewAgentWorkerPool(ctx, 3, mockAgentGetter(false, 50*time.Millisecond))

	job := AgentJob{
		Role:      agent.RolePM,
		Prompt:    "test task",
		Context:   map[string]interface{}{"workDir": "/tmp"},
		TaskID:    "test-1",
		ProjectID: "test-project",
	}

	// Test submission
	err := pool.Submit(job)
	if err != nil {
		t.Fatalf("Failed to submit job: %v", err)
	}

	// Shutdown immediately
	pool.Shutdown()
}

func TestWorkerPoolCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	pool := NewAgentWorkerPool(ctx, 2, mockAgentGetter(false, 1*time.Second))

	// Submit a job
	job := AgentJob{
		Role:      agent.RolePM,
		Prompt:    "test task",
		Context:   map[string]interface{}{"workDir": "/tmp"},
		TaskID:    "test-cancel",
		ProjectID: "test-project",
	}

	err := pool.Submit(job)
	if err != nil {
		t.Fatalf("Failed to submit job: %v", err)
	}

	// Cancel immediately
	cancel()
	time.Sleep(100 * time.Millisecond)

	// Pool should stop accepting jobs
	err = pool.Submit(job)
	if err == nil {
		t.Error("Expected error when submitting after cancellation")
	}

	pool.Shutdown()
}

func TestWorkerPoolTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	pool := NewAgentWorkerPool(ctx, 2, mockAgentGetter(false, 5*time.Second))

	job := AgentJob{
		Role:      agent.RolePM,
		Prompt:    "test task",
		Context:   map[string]interface{}{"workDir": "/tmp"},
		TaskID:    "test-timeout",
		ProjectID: "test-project",
	}

	err := pool.Submit(job)
	if err != nil {
		t.Fatalf("Failed to submit job: %v", err)
	}

	go pool.Wait()

	// Should get result with timeout error
	select {
	case result := <-pool.Results():
		if result.Error == nil {
			t.Error("Expected timeout error")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for result")
	}
}

func TestWorkerPoolActiveTracking(t *testing.T) {
	ctx := context.Background()
	pool := NewAgentWorkerPool(ctx, 1, mockAgentGetter(false, 200*time.Millisecond))

	// Initially no active agents
	if pool.IsActive(agent.RolePM) {
		t.Error("Expected no active agents initially")
	}

	// Submit job
	job := AgentJob{
		Role:      agent.RolePM,
		Prompt:    "test task",
		Context:   map[string]interface{}{"workDir": "/tmp"},
		TaskID:    "test-active",
		ProjectID: "test-project",
	}

	pool.Submit(job)

	// Give it time to start processing
	time.Sleep(50 * time.Millisecond)

	// Should be active now
	if !pool.IsActive(agent.RolePM) {
		t.Error("Expected PM to be active")
	}

	// Wait for completion
	go pool.Wait()
	<-pool.Results()

	// Should not be active anymore
	time.Sleep(50 * time.Millisecond)
	if pool.IsActive(agent.RolePM) {
		t.Error("Expected PM to be inactive after completion")
	}
}

func TestWorkerPoolAgentCreationError(t *testing.T) {
	ctx := context.Background()
	pool := NewAgentWorkerPool(ctx, 2, mockAgentGetter(true, 0))

	job := AgentJob{
		Role:      agent.RolePM,
		Prompt:    "test task",
		Context:   map[string]interface{}{"workDir": "/tmp"},
		TaskID:    "test-error",
		ProjectID: "test-project",
	}

	err := pool.Submit(job)
	if err != nil {
		t.Fatalf("Failed to submit job: %v", err)
	}

	go pool.Wait()

	// Should get result with agent creation error
	select {
	case result := <-pool.Results():
		if result.Error == nil {
			t.Error("Expected agent creation error")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for result")
	}
}

func BenchmarkWorkerPoolParallel(b *testing.B) {
	roles := []agent.Role{agent.RolePM, agent.RoleUX, agent.RoleUI, agent.RoleSecurity, agent.RoleArchitect}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		pool := NewAgentWorkerPool(ctx, 5, mockAgentGetter(false, 10*time.Millisecond))

		for j, role := range roles {
			job := AgentJob{
				Role:      role,
				Prompt:    "benchmark task",
				Context:   map[string]interface{}{"workDir": "/tmp"},
				TaskID:    string(rune('a' + j)),
				ProjectID: "bench",
			}
			pool.Submit(job)
		}

		go pool.Wait()

		// Drain results
		for j := 0; j < len(roles); j++ {
			<-pool.Results()
		}
	}
}
