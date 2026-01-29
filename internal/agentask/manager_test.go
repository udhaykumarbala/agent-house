package agentask

import (
	"testing"

	"pty-claude-test/internal/agent"
)

func TestManagerCreateTask(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewManager(tmpDir)

	task, err := mgr.CreateTask(
		agent.RolePM,
		TaskTypeResearch,
		"Research user needs",
		"Analyze user requirements and pain points",
		"test-project",
		"parent-123",
	)

	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	if task.AgentRole != agent.RolePM {
		t.Errorf("Expected role PM, got %s", task.AgentRole)
	}

	if task.Status != StatusPending {
		t.Errorf("Expected status pending, got %s", task.Status)
	}

	if task.Progress != 0 {
		t.Errorf("Expected progress 0, got %d", task.Progress)
	}
}

func TestManagerGetTask(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewManager(tmpDir)

	created, _ := mgr.CreateTask(
		agent.RolePM,
		TaskTypeResearch,
		"Test task",
		"Description",
		"test-project",
		"parent-123",
	)

	retrieved, err := mgr.GetTask(created.ID)
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}

	if retrieved.ID != created.ID {
		t.Errorf("Task ID mismatch: %s != %s", retrieved.ID, created.ID)
	}

	if retrieved.Title != created.Title {
		t.Errorf("Task title mismatch: %s != %s", retrieved.Title, created.Title)
	}
}

func TestManagerStartTask(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewManager(tmpDir)

	task, _ := mgr.CreateTask(
		agent.RolePM,
		TaskTypeResearch,
		"Test task",
		"Description",
		"test-project",
		"parent-123",
	)

	// Start the task
	err := mgr.StartTask(task.ID)
	if err != nil {
		t.Fatalf("Failed to start task: %v", err)
	}

	// Verify status updated
	updated, _ := mgr.GetTask(task.ID)
	if updated.Status != StatusInProgress {
		t.Errorf("Expected status in_progress, got %s", updated.Status)
	}

	if updated.StartedAt == nil {
		t.Error("StartedAt should be set")
	}

	// Verify agent status updated
	status, _ := mgr.GetAgentStatus(agent.RolePM)
	if status.State != StateWorking {
		t.Errorf("Expected agent state working, got %s", status.State)
	}

	if status.CurrentTaskID != task.ID {
		t.Errorf("Expected current task %s, got %s", task.ID, status.CurrentTaskID)
	}
}

func TestManagerCompleteTask(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewManager(tmpDir)

	task, _ := mgr.CreateTask(
		agent.RolePM,
		TaskTypeResearch,
		"Test task",
		"Description",
		"test-project",
		"parent-123",
	)

	mgr.StartTask(task.ID)
	err := mgr.CompleteTask(task.ID)
	if err != nil {
		t.Fatalf("Failed to complete task: %v", err)
	}

	// Verify status updated
	updated, _ := mgr.GetTask(task.ID)
	if updated.Status != StatusCompleted {
		t.Errorf("Expected status completed, got %s", updated.Status)
	}

	if updated.Progress != 100 {
		t.Errorf("Expected progress 100, got %d", updated.Progress)
	}

	if updated.CompletedAt == nil {
		t.Error("CompletedAt should be set")
	}

	// Verify agent status updated
	status, _ := mgr.GetAgentStatus(agent.RolePM)
	if status.State != StateIdle {
		t.Errorf("Expected agent state idle, got %s", status.State)
	}

	if status.TasksCompleted != 1 {
		t.Errorf("Expected 1 completed task, got %d", status.TasksCompleted)
	}
}

func TestManagerFailTask(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewManager(tmpDir)

	task, _ := mgr.CreateTask(
		agent.RolePM,
		TaskTypeResearch,
		"Test task",
		"Description",
		"test-project",
		"parent-123",
	)

	mgr.StartTask(task.ID)
	err := mgr.FailTask(task.ID, &testError{"test error"})
	if err != nil {
		t.Fatalf("Failed to fail task: %v", err)
	}

	// Verify status updated
	updated, _ := mgr.GetTask(task.ID)
	if updated.Status != StatusFailed {
		t.Errorf("Expected status failed, got %s", updated.Status)
	}

	if updated.Error == "" {
		t.Error("Error message should be set")
	}

	// Verify agent status updated
	status, _ := mgr.GetAgentStatus(agent.RolePM)
	if status.State != StateError {
		t.Errorf("Expected agent state error, got %s", status.State)
	}

	if status.TasksFailed != 1 {
		t.Errorf("Expected 1 failed task, got %d", status.TasksFailed)
	}
}

func TestManagerUpdateProgress(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewManager(tmpDir)

	task, _ := mgr.CreateTask(
		agent.RolePM,
		TaskTypeResearch,
		"Test task",
		"Description",
		"test-project",
		"parent-123",
	)

	// Update progress
	err := mgr.UpdateProgress(task.ID, 50)
	if err != nil {
		t.Fatalf("Failed to update progress: %v", err)
	}

	// Verify progress updated
	updated, _ := mgr.GetTask(task.ID)
	if updated.Progress != 50 {
		t.Errorf("Expected progress 50, got %d", updated.Progress)
	}

	// Test bounds
	mgr.UpdateProgress(task.ID, 150) // Should cap at 100
	updated, _ = mgr.GetTask(task.ID)
	if updated.Progress != 100 {
		t.Errorf("Expected progress capped at 100, got %d", updated.Progress)
	}
}

func TestManagerAddOutput(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewManager(tmpDir)

	task, _ := mgr.CreateTask(
		agent.RolePM,
		TaskTypeResearch,
		"Test task",
		"Description",
		"test-project",
		"parent-123",
	)

	output := AgentOutput{
		ID:       "out-1",
		Type:     OutputResearchFindings,
		Title:    "Research Summary",
		FilePath: "/tmp/research.md",
		Format:   "markdown",
	}

	err := mgr.AddOutput(task.ID, output)
	if err != nil {
		t.Fatalf("Failed to add output: %v", err)
	}

	// Verify output added
	updated, _ := mgr.GetTask(task.ID)
	if len(updated.Outputs) != 1 {
		t.Errorf("Expected 1 output, got %d", len(updated.Outputs))
	}

	if updated.Outputs[0].Title != "Research Summary" {
		t.Errorf("Output title mismatch: %s", updated.Outputs[0].Title)
	}
}

func TestManagerGetAgentTasks(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewManager(tmpDir)

	// Create tasks for different agents
	mgr.CreateTask(agent.RolePM, TaskTypeResearch, "PM Task 1", "", "proj", "p1")
	mgr.CreateTask(agent.RolePM, TaskTypeResearch, "PM Task 2", "", "proj", "p1")
	mgr.CreateTask(agent.RoleUX, TaskTypeResearch, "UX Task 1", "", "proj", "p1")

	// Get PM tasks
	pmTasks := mgr.GetAgentTasks(agent.RolePM)
	if len(pmTasks) != 2 {
		t.Errorf("Expected 2 PM tasks, got %d", len(pmTasks))
	}

	// Get UX tasks
	uxTasks := mgr.GetAgentTasks(agent.RoleUX)
	if len(uxTasks) != 1 {
		t.Errorf("Expected 1 UX task, got %d", len(uxTasks))
	}
}

func TestManagerGetProjectTasks(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewManager(tmpDir)

	// Create tasks for different projects
	mgr.CreateTask(agent.RolePM, TaskTypeResearch, "Task 1", "", "project-a", "p1")
	mgr.CreateTask(agent.RoleUX, TaskTypeResearch, "Task 2", "", "project-a", "p1")
	mgr.CreateTask(agent.RolePM, TaskTypeResearch, "Task 3", "", "project-b", "p2")

	// Get project-a tasks
	projATasks := mgr.GetProjectTasks("project-a")
	if len(projATasks) != 2 {
		t.Errorf("Expected 2 tasks for project-a, got %d", len(projATasks))
	}

	// Get project-b tasks
	projBTasks := mgr.GetProjectTasks("project-b")
	if len(projBTasks) != 1 {
		t.Errorf("Expected 1 task for project-b, got %d", len(projBTasks))
	}
}

func TestManagerFilterTasks(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewManager(tmpDir)

	// Create various tasks
	task1, _ := mgr.CreateTask(agent.RolePM, TaskTypeResearch, "Task 1", "", "proj-a", "p1")
	task2, _ := mgr.CreateTask(agent.RoleUX, TaskTypeResearch, "Task 2", "", "proj-a", "p1")
	mgr.CreateTask(agent.RolePM, TaskTypePlanning, "Task 3", "", "proj-b", "p2")

	// Verify task2 created
	if task2.ID == "" {
		t.Fatal("task2 should have an ID")
	}

	// Start one task
	mgr.StartTask(task1.ID)
	mgr.CompleteTask(task1.ID)

	// Filter by project
	projAResults := mgr.FilterTasks("proj-a", "", "")
	if len(projAResults) != 2 {
		t.Errorf("Expected 2 tasks for proj-a, got %d", len(projAResults))
	}

	// Filter by role
	pmResults := mgr.FilterTasks("", agent.RolePM, "")
	if len(pmResults) != 2 {
		t.Errorf("Expected 2 PM tasks, got %d", len(pmResults))
	}

	// Filter by status
	completedResults := mgr.FilterTasks("", "", StatusCompleted)
	if len(completedResults) != 1 {
		t.Errorf("Expected 1 completed task, got %d", len(completedResults))
	}

	// Combined filters
	combinedResults := mgr.FilterTasks("proj-a", agent.RolePM, StatusCompleted)
	if len(combinedResults) != 1 {
		t.Errorf("Expected 1 task with combined filters, got %d", len(combinedResults))
	}
}

func TestManagerGetAllAgentStatuses(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewManager(tmpDir)

	// Create tasks for multiple agents
	task1, _ := mgr.CreateTask(agent.RolePM, TaskTypeResearch, "Task 1", "", "proj", "p1")
	task2, _ := mgr.CreateTask(agent.RoleUX, TaskTypeResearch, "Task 2", "", "proj", "p1")

	mgr.StartTask(task1.ID)
	mgr.CompleteTask(task1.ID)

	mgr.StartTask(task2.ID)

	// Verify task2 is working
	if task2.Status != StatusPending {
		// This check uses task2
	}

	// Get all statuses
	statuses := mgr.GetAllAgentStatuses()

	// Should have status for PM and UX
	if len(statuses) < 2 {
		t.Errorf("Expected at least 2 agent statuses, got %d", len(statuses))
	}

	// Verify PM status
	if pmStatus, ok := statuses[agent.RolePM]; ok {
		if pmStatus.TasksCompleted != 1 {
			t.Errorf("Expected PM to have 1 completed task, got %d", pmStatus.TasksCompleted)
		}
	} else {
		t.Error("PM status not found")
	}

	// Verify UX status
	if uxStatus, ok := statuses[agent.RoleUX]; ok {
		if uxStatus.State != StateWorking {
			t.Errorf("Expected UX state working, got %s", uxStatus.State)
		}
	} else {
		t.Error("UX status not found")
	}
}

// testError is a simple error type for testing
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
