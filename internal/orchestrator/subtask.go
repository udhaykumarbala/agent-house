package orchestrator

import (
	"fmt"
	"time"

	"pty-claude-test/internal/agent"
)

// SubTaskStatus represents the status of a subtask
type SubTaskStatus string

const (
	SubTaskStatusPending       SubTaskStatus = "pending"
	SubTaskStatusInProgress    SubTaskStatus = "in_progress"
	SubTaskStatusCompleted     SubTaskStatus = "completed"
	SubTaskStatusNeedsRevision SubTaskStatus = "needs_revision"
)

// SubTask represents a concrete work item within a development phase
type SubTask struct {
	ID                 string          `json:"id"`
	TaskID             string          `json:"task_id"`
	PhaseIndex         int             `json:"phase_index"`
	Title              string          `json:"title"`
	Description        string          `json:"description"`
	AssignedAgents     []agent.Role    `json:"assigned_agents"`
	Dependencies       []string        `json:"dependencies,omitempty"`
	CompletionCriteria []string        `json:"completion_criteria"`
	Status             SubTaskStatus   `json:"status"`
	CreatedAt          time.Time       `json:"created_at"`
	StartedAt          *time.Time      `json:"started_at,omitempty"`
	CompletedAt        *time.Time      `json:"completed_at,omitempty"`
	Iteration          int             `json:"iteration"`
}

// UpdateStatus updates the subtask status and timestamps
func (st *SubTask) UpdateStatus(status SubTaskStatus) {
	st.Status = status

	now := time.Now()
	if status == SubTaskStatusInProgress && st.StartedAt == nil {
		st.StartedAt = &now
	}
	if status == SubTaskStatusCompleted && st.CompletedAt == nil {
		st.CompletedAt = &now
	}
}

// IsBlocked returns true if any dependencies are not completed
func (st *SubTask) IsBlocked(phase *DevelopmentPhase) bool {
	if len(st.Dependencies) == 0 {
		return false
	}

	// Check if all dependencies are completed
	for _, depID := range st.Dependencies {
		found := false
		for _, subtask := range phase.SubTasks {
			if subtask.ID == depID {
				found = true
				if subtask.Status != SubTaskStatusCompleted {
					return true
				}
				break
			}
		}
		if !found {
			// Dependency not found, consider blocked
			return true
		}
	}

	return false
}

// GetNextAvailableSubTask returns the next subtask that can be started
func GetNextAvailableSubTask(phase *DevelopmentPhase) *SubTask {
	for i := range phase.SubTasks {
		subtask := &phase.SubTasks[i]
		if subtask.Status == SubTaskStatusPending && !subtask.IsBlocked(phase) {
			return subtask
		}
	}
	return nil
}

// GetSubTasksByAgent returns all subtasks assigned to a specific agent
func GetSubTasksByAgent(phase *DevelopmentPhase, role agent.Role) []*SubTask {
	var tasks []*SubTask
	for i := range phase.SubTasks {
		subtask := &phase.SubTasks[i]
		for _, assignedRole := range subtask.AssignedAgents {
			if assignedRole == role {
				tasks = append(tasks, subtask)
				break
			}
		}
	}
	return tasks
}

// ValidateDependencies checks for circular dependencies in subtasks
func ValidateDependencies(phase *DevelopmentPhase) error {
	// Build dependency graph
	graph := make(map[string][]string)
	for _, subtask := range phase.SubTasks {
		graph[subtask.ID] = subtask.Dependencies
	}

	// Check for cycles using DFS
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	var hasCycle func(id string) bool
	hasCycle = func(id string) bool {
		visited[id] = true
		recStack[id] = true

		for _, dep := range graph[id] {
			if !visited[dep] {
				if hasCycle(dep) {
					return true
				}
			} else if recStack[dep] {
				return true
			}
		}

		recStack[id] = false
		return false
	}

	for id := range graph {
		if !visited[id] {
			if hasCycle(id) {
				return fmt.Errorf("circular dependency detected in subtasks")
			}
		}
	}

	return nil
}

// GetTopologicalOrder returns subtasks in dependency order
func GetTopologicalOrder(phase *DevelopmentPhase) ([]*SubTask, error) {
	// Validate no circular dependencies
	if err := ValidateDependencies(phase); err != nil {
		return nil, err
	}

	// Build in-degree map
	inDegree := make(map[string]int)
	graph := make(map[string][]string)

	for i := range phase.SubTasks {
		subtask := &phase.SubTasks[i]
		inDegree[subtask.ID] = len(subtask.Dependencies)
		for _, dep := range subtask.Dependencies {
			graph[dep] = append(graph[dep], subtask.ID)
		}
	}

	// Kahn's algorithm for topological sort
	var result []*SubTask
	queue := make([]*SubTask, 0)

	// Start with tasks that have no dependencies
	for i := range phase.SubTasks {
		subtask := &phase.SubTasks[i]
		if inDegree[subtask.ID] == 0 {
			queue = append(queue, subtask)
		}
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		result = append(result, current)

		for _, nextID := range graph[current.ID] {
			inDegree[nextID]--
			if inDegree[nextID] == 0 {
				// Find the subtask
				for i := range phase.SubTasks {
					if phase.SubTasks[i].ID == nextID {
						queue = append(queue, &phase.SubTasks[i])
						break
					}
				}
			}
		}
	}

	if len(result) != len(phase.SubTasks) {
		return nil, fmt.Errorf("unable to determine topological order")
	}

	return result, nil
}
