package agentask

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Storage handles persistence of agent tasks
type Storage struct {
	baseDir string
	mu      sync.RWMutex
}

// NewStorage creates a new storage instance
func NewStorage(baseDir string) *Storage {
	return &Storage{
		baseDir: baseDir,
	}
}

// SaveTask persists a task to disk
func (s *Storage) SaveTask(task *AgentTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create directory structure: projects/{projectID}/.tasks/agent-tasks/{agentRole}/
	tasksDir := filepath.Join(s.baseDir, task.ProjectID, ".tasks", "agent-tasks", string(task.AgentRole))
	if err := os.MkdirAll(tasksDir, 0755); err != nil {
		return fmt.Errorf("failed to create tasks directory: %w", err)
	}

	// Save task as JSON
	taskPath := filepath.Join(tasksDir, task.ID+".json")
	data, err := json.MarshalIndent(task, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	if err := os.WriteFile(taskPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write task file: %w", err)
	}

	// Update index file
	if err := s.updateIndex(task); err != nil {
		return fmt.Errorf("failed to update index: %w", err)
	}

	return nil
}

// LoadTask loads a task from disk
func (s *Storage) LoadTask(projectID, taskID string) (*AgentTask, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Search for the task file across all agent directories
	agentTasksDir := filepath.Join(s.baseDir, projectID, ".tasks", "agent-tasks")

	var task *AgentTask
	err := filepath.Walk(agentTasksDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Base(path) == taskID+".json" {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			task = &AgentTask{}
			return json.Unmarshal(data, task)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to load task: %w", err)
	}
	if task == nil {
		return nil, fmt.Errorf("task not found: %s", taskID)
	}

	return task, nil
}

// LoadProjectTasks loads all tasks for a project
func (s *Storage) LoadProjectTasks(projectID string) ([]*AgentTask, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	agentTasksDir := filepath.Join(s.baseDir, projectID, ".tasks", "agent-tasks")
	if _, err := os.Stat(agentTasksDir); os.IsNotExist(err) {
		return []*AgentTask{}, nil
	}

	var tasks []*AgentTask
	err := filepath.Walk(agentTasksDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".json" && filepath.Base(path) != "index.json" {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			task := &AgentTask{}
			if err := json.Unmarshal(data, task); err != nil {
				return err
			}
			tasks = append(tasks, task)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to load project tasks: %w", err)
	}

	return tasks, nil
}

// updateIndex updates the index file with all tasks
func (s *Storage) updateIndex(task *AgentTask) error {
	indexPath := filepath.Join(s.baseDir, task.ProjectID, ".tasks", "agent-tasks", "index.json")

	// Load existing index
	var index map[string]interface{}
	if data, err := os.ReadFile(indexPath); err == nil {
		json.Unmarshal(data, &index)
	}
	if index == nil {
		index = make(map[string]interface{})
	}

	// Add task summary to index
	taskSummary := map[string]interface{}{
		"id":          task.ID,
		"agent_role":  task.AgentRole,
		"task_type":   task.TaskType,
		"status":      task.Status,
		"title":       task.Title,
		"progress":    task.Progress,
		"assigned_at": task.AssignedAt,
		"started_at":  task.StartedAt,
		"completed_at": task.CompletedAt,
	}

	if index["tasks"] == nil {
		index["tasks"] = make(map[string]interface{})
	}
	index["tasks"].(map[string]interface{})[task.ID] = taskSummary

	// Write index
	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal index: %w", err)
	}

	if err := os.WriteFile(indexPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write index: %w", err)
	}

	return nil
}

// SaveOutput saves an output file and creates metadata
func (s *Storage) SaveOutput(task *AgentTask, output AgentOutput, content string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create outputs directory
	outputsDir := filepath.Join(s.baseDir, task.ProjectID, ".tasks", "agent-tasks", string(task.AgentRole), "outputs")
	if err := os.MkdirAll(outputsDir, 0755); err != nil {
		return fmt.Errorf("failed to create outputs directory: %w", err)
	}

	// Save output file
	outputPath := filepath.Join(outputsDir, filepath.Base(output.FilePath))
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	// Save output metadata
	metadataPath := outputPath + ".meta.json"
	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal output metadata: %w", err)
	}

	if err := os.WriteFile(metadataPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write output metadata: %w", err)
	}

	return nil
}
