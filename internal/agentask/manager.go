package agentask

import (
	"fmt"
	"sync"

	"pty-claude-test/internal/agent"
)

// Manager manages agent tasks
type Manager struct {
	mu            sync.RWMutex
	tasks         map[string]*AgentTask           // taskID -> AgentTask
	agentTasks    map[agent.Role][]*AgentTask     // agentRole -> tasks
	projectTasks  map[string][]*AgentTask         // projectID -> tasks
	agentStatuses map[agent.Role]*AgentStatus     // agentRole -> status
	storage       *Storage
}

// NewManager creates a new agent task manager
func NewManager(baseDir string) *Manager {
	return &Manager{
		tasks:         make(map[string]*AgentTask),
		agentTasks:    make(map[agent.Role][]*AgentTask),
		projectTasks:  make(map[string][]*AgentTask),
		agentStatuses: make(map[agent.Role]*AgentStatus),
		storage:       NewStorage(baseDir),
	}
}

// CreateTask creates a new agent task
func (m *Manager) CreateTask(role agent.Role, taskType AgentTaskType, title, description, projectID, parentTaskID string) (*AgentTask, error) {
	task := NewAgentTask(role, taskType, title, description, projectID, parentTaskID)

	m.mu.Lock()
	defer m.mu.Unlock()

	// Store in maps
	m.tasks[task.ID] = task
	m.agentTasks[role] = append(m.agentTasks[role], task)
	m.projectTasks[projectID] = append(m.projectTasks[projectID], task)

	// Initialize agent status if needed
	if _, exists := m.agentStatuses[role]; !exists {
		m.agentStatuses[role] = &AgentStatus{
			Role:           role,
			State:          StateIdle,
			TasksCompleted: 0,
			TasksFailed:    0,
		}
	}

	// Persist to storage
	if err := m.storage.SaveTask(task); err != nil {
		return nil, fmt.Errorf("failed to save task: %w", err)
	}

	return task, nil
}

// GetTask retrieves a task by ID
func (m *Manager) GetTask(taskID string) (*AgentTask, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	task, exists := m.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("task not found: %s", taskID)
	}

	return task, nil
}

// GetAgentTasks retrieves all tasks for a specific agent
func (m *Manager) GetAgentTasks(role agent.Role) []*AgentTask {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tasks := m.agentTasks[role]
	result := make([]*AgentTask, len(tasks))
	copy(result, tasks)
	return result
}

// GetProjectTasks retrieves all tasks for a specific project
func (m *Manager) GetProjectTasks(projectID string) []*AgentTask {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tasks := m.projectTasks[projectID]
	result := make([]*AgentTask, len(tasks))
	copy(result, tasks)
	return result
}

// GetAllTasks retrieves all tasks
func (m *Manager) GetAllTasks() []*AgentTask {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*AgentTask, 0, len(m.tasks))
	for _, task := range m.tasks {
		result = append(result, task)
	}
	return result
}

// StartTask marks a task as started
func (m *Manager) StartTask(taskID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	task, exists := m.tasks[taskID]
	if !exists {
		return fmt.Errorf("task not found: %s", taskID)
	}

	task.Start()

	// Update agent status
	if status, exists := m.agentStatuses[task.AgentRole]; exists {
		status.State = StateWorking
		status.CurrentTaskID = taskID
		status.CurrentTaskType = task.TaskType
		now := task.StartedAt
		status.LastActiveAt = now
	}

	// Persist
	if err := m.storage.SaveTask(task); err != nil {
		return fmt.Errorf("failed to save task: %w", err)
	}

	return nil
}

// CompleteTask marks a task as completed
func (m *Manager) CompleteTask(taskID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	task, exists := m.tasks[taskID]
	if !exists {
		return fmt.Errorf("task not found: %s", taskID)
	}

	task.Complete()

	// Update agent status
	if status, exists := m.agentStatuses[task.AgentRole]; exists {
		status.State = StateIdle
		status.CurrentTaskID = ""
		status.TasksCompleted++
		now := task.CompletedAt
		status.LastActiveAt = now
	}

	// Persist
	if err := m.storage.SaveTask(task); err != nil {
		return fmt.Errorf("failed to save task: %w", err)
	}

	return nil
}

// FailTask marks a task as failed
func (m *Manager) FailTask(taskID string, err error) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	task, exists := m.tasks[taskID]
	if !exists {
		return fmt.Errorf("task not found: %s", taskID)
	}

	task.Fail(err)

	// Update agent status
	if status, exists := m.agentStatuses[task.AgentRole]; exists {
		status.State = StateError
		status.CurrentTaskID = ""
		status.TasksFailed++
		now := task.CompletedAt
		status.LastActiveAt = now
	}

	// Persist
	if saveErr := m.storage.SaveTask(task); saveErr != nil {
		return fmt.Errorf("failed to save task: %w", saveErr)
	}

	return nil
}

// UpdateProgress updates the progress of a task
func (m *Manager) UpdateProgress(taskID string, progress int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	task, exists := m.tasks[taskID]
	if !exists {
		return fmt.Errorf("task not found: %s", taskID)
	}

	task.UpdateProgress(progress)

	// Persist
	if err := m.storage.SaveTask(task); err != nil {
		return fmt.Errorf("failed to save task: %w", err)
	}

	return nil
}

// AddOutput adds an output to a task
func (m *Manager) AddOutput(taskID string, output AgentOutput) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	task, exists := m.tasks[taskID]
	if !exists {
		return fmt.Errorf("task not found: %s", taskID)
	}

	task.AddOutput(output)

	// Persist
	if err := m.storage.SaveTask(task); err != nil {
		return fmt.Errorf("failed to save task: %w", err)
	}

	return nil
}

// GetAgentStatus retrieves the status of an agent
func (m *Manager) GetAgentStatus(role agent.Role) (*AgentStatus, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	status, exists := m.agentStatuses[role]
	if !exists {
		return &AgentStatus{
			Role:           role,
			State:          StateIdle,
			TasksCompleted: 0,
			TasksFailed:    0,
		}, nil
	}

	// Return a copy
	statusCopy := *status
	return &statusCopy, nil
}

// GetAllAgentStatuses retrieves statuses for all agents
func (m *Manager) GetAllAgentStatuses() map[agent.Role]*AgentStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[agent.Role]*AgentStatus, len(m.agentStatuses))
	for role, status := range m.agentStatuses {
		statusCopy := *status
		result[role] = &statusCopy
	}
	return result
}

// FilterTasks filters tasks based on criteria
func (m *Manager) FilterTasks(projectID string, role agent.Role, status AgentTaskStatus) []*AgentTask {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*AgentTask

	for _, task := range m.tasks {
		// Apply filters
		if projectID != "" && task.ProjectID != projectID {
			continue
		}
		if role != "" && task.AgentRole != role {
			continue
		}
		if status != "" && task.Status != status {
			continue
		}

		result = append(result, task)
	}

	return result
}
