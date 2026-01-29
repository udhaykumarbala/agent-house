package orchestrator

import (
	"testing"
	"time"

	"pty-claude-test/internal/agent"
	"pty-claude-test/internal/message"
)

func TestPhaseExecutorParallelExecution(t *testing.T) {
	store := message.NewStore()

	config := ExecutorConfig{
		EnableParallel: true,
		MaxWorkers:     0,
		AgentTimeout:   5 * time.Minute,
		PhaseTimeout:   15 * time.Minute,
		WorkDir:        t.TempDir(),
		ProjectID:      "test-project",
		TaskID:         "test-task",
	}

	fm := NewConcurrentFileManager()
	defer fm.Shutdown()

	executor := NewPhaseExecutor(config, mockAgentGetter(false, 100*time.Millisecond), fm, store)

	roles := []agent.Role{agent.RolePM, agent.RoleUX, agent.RoleUI, agent.RoleSecurity, agent.RoleArchitect}

	start := time.Now()
	result, err := executor.ExecuteAgentsParallel(roles, "Test task prompt")
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("ExecuteAgentsParallel failed: %v", err)
	}

	// Should complete much faster than sequential (5 * 100ms = 500ms)
	// With parallel, should be around 100-200ms
	if duration > 1*time.Second {
		t.Errorf("Parallel execution took too long: %v", duration)
	}

	if !result.Success {
		t.Error("Expected successful execution")
	}

	if len(result.Messages) != len(roles) {
		t.Errorf("Expected %d messages, got %d", len(roles), len(result.Messages))
	}

	// Verify speedup
	if result.MaxAgentDuration == 0 {
		t.Error("MaxAgentDuration should be set")
	}

	speedup := float64(result.MaxAgentDuration) / float64(result.Duration)
	if speedup < 1.0 {
		t.Errorf("Expected speedup >= 1.0, got %.2f", speedup)
	}

	t.Logf("Parallel execution: %v (max agent: %v, speedup: %.1fx)",
		result.Duration, result.MaxAgentDuration, speedup)
}

func TestPhaseExecutorSequentialFallback(t *testing.T) {
	store := message.NewStore()

	config := ExecutorConfig{
		EnableParallel: false, // Disable parallel
		WorkDir:        t.TempDir(),
		ProjectID:      "test-project",
		TaskID:         "test-task",
	}

	fm := NewConcurrentFileManager()
	defer fm.Shutdown()

	executor := NewPhaseExecutor(config, mockAgentGetter(false, 50*time.Millisecond), fm, store)

	roles := []agent.Role{agent.RolePM, agent.RoleUX, agent.RoleUI}

	start := time.Now()
	result, err := executor.ExecuteAgentsParallel(roles, "Test task prompt")
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Sequential execution failed: %v", err)
	}

	// Sequential should take at least 3 * 50ms = 150ms
	if duration < 100*time.Millisecond {
		t.Errorf("Sequential execution too fast: %v (expected >= 100ms)", duration)
	}

	if !result.Success {
		t.Error("Expected successful execution")
	}

	t.Logf("Sequential execution: %v", duration)
}

func TestPhaseExecutorErrorHandling(t *testing.T) {
	store := message.NewStore()

	config := ExecutorConfig{
		EnableParallel: true,
		MaxWorkers:     0,
		WorkDir:        t.TempDir(),
		ProjectID:      "test-project",
		TaskID:         "test-task",
	}

	fm := NewConcurrentFileManager()
	defer fm.Shutdown()

	// Mock getter that returns errors
	executor := NewPhaseExecutor(config, mockAgentGetter(true, 0), fm, store)

	roles := []agent.Role{agent.RolePM, agent.RoleUX}

	result, err := executor.ExecuteAgentsParallel(roles, "Test task prompt")

	if err != nil {
		t.Fatalf("ExecuteAgentsParallel failed: %v", err)
	}

	// Should have errors
	if !result.HasErrors() {
		t.Error("Expected errors from failed agent creation")
	}

	if result.ErrorCount() != len(roles) {
		t.Errorf("Expected %d errors, got %d", len(roles), result.ErrorCount())
	}

	if result.Success {
		t.Error("Expected unsuccessful execution")
	}
}

func TestConcurrentFileOperations(t *testing.T) {
	tmpDir := t.TempDir()
	fm := NewConcurrentFileManager()
	defer fm.Shutdown()

	store := message.NewStore()

	config := ExecutorConfig{
		EnableParallel: true,
		MaxWorkers:     0,
		WorkDir:        tmpDir,
		ProjectID:      "test-project",
		TaskID:         "test-task",
	}

	executor := NewPhaseExecutor(config, mockAgentGetter(false, 50*time.Millisecond), fm, store)

	// All agents will try to create files
	roles := []agent.Role{agent.RolePM, agent.RoleUX, agent.RoleUI}

	result, err := executor.ExecuteAgentsParallel(roles, "Create files")

	if err != nil {
		t.Fatalf("Execution failed: %v", err)
	}

	// File manager should track all created files
	files := fm.GetCreatedFiles()
	if len(files) > 0 {
		t.Logf("Created %d files", len(files))
	}

	// All agents should complete
	if len(result.Messages) != len(roles) {
		t.Errorf("Expected %d messages, got %d", len(roles), len(result.Messages))
	}
}

func BenchmarkParallelExecution5Agents(b *testing.B) {
	store := message.NewStore()
	fm := NewConcurrentFileManager()
	defer fm.Shutdown()

	config := ExecutorConfig{
		EnableParallel: true,
		MaxWorkers:     0,
		WorkDir:        b.TempDir(),
		ProjectID:      "bench",
		TaskID:         "bench",
	}

	executor := NewPhaseExecutor(config, mockAgentGetter(false, 10*time.Millisecond), fm, store)
	roles := []agent.Role{agent.RolePM, agent.RoleUX, agent.RoleUI, agent.RoleSecurity, agent.RoleArchitect}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		executor.ExecuteAgentsParallel(roles, "Benchmark task")
	}
}

func BenchmarkSequentialExecution5Agents(b *testing.B) {
	store := message.NewStore()
	fm := NewConcurrentFileManager()
	defer fm.Shutdown()

	config := ExecutorConfig{
		EnableParallel: false, // Sequential
		WorkDir:        b.TempDir(),
		ProjectID:      "bench",
		TaskID:         "bench",
	}

	executor := NewPhaseExecutor(config, mockAgentGetter(false, 10*time.Millisecond), fm, store)
	roles := []agent.Role{agent.RolePM, agent.RoleUX, agent.RoleUI, agent.RoleSecurity, agent.RoleArchitect}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		executor.ExecuteAgentsParallel(roles, "Benchmark task")
	}
}

func BenchmarkParallelVsSequential(b *testing.B) {
	roles := []agent.Role{agent.RolePM, agent.RoleUX, agent.RoleUI, agent.RoleSecurity, agent.RoleArchitect}

	b.Run("Parallel", func(b *testing.B) {
		store := message.NewStore()
		fm := NewConcurrentFileManager()
		defer fm.Shutdown()

		config := ExecutorConfig{
			EnableParallel: true,
			MaxWorkers:     0,
			WorkDir:        b.TempDir(),
			ProjectID:      "bench",
			TaskID:         "bench",
		}

		executor := NewPhaseExecutor(config, mockAgentGetter(false, 10*time.Millisecond), fm, store)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			executor.ExecuteAgentsParallel(roles, "Benchmark task")
		}
	})

	b.Run("Sequential", func(b *testing.B) {
		store := message.NewStore()
		fm := NewConcurrentFileManager()
		defer fm.Shutdown()

		config := ExecutorConfig{
			EnableParallel: false,
			WorkDir:        b.TempDir(),
			ProjectID:      "bench",
			TaskID:         "bench",
		}

		executor := NewPhaseExecutor(config, mockAgentGetter(false, 10*time.Millisecond), fm, store)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			executor.ExecuteAgentsParallel(roles, "Benchmark task")
		}
	})
}
