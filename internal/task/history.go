package task

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"pty-claude-test/internal/message"
)

// TaskStatus represents the status of a task
type TaskStatus string

const (
	StatusPending   TaskStatus = "pending"
	StatusRunning   TaskStatus = "running"
	StatusCompleted TaskStatus = "completed"
	StatusFailed    TaskStatus = "failed"
)

// TaskSummary is a lightweight representation for the task list
type TaskSummary struct {
	TaskID       string     `json:"taskId"`
	Summary      string     `json:"summary"`
	CreatedAt    time.Time  `json:"createdAt"`
	CompletedAt  *time.Time `json:"completedAt,omitempty"`
	Status       TaskStatus `json:"status"`
	Turns        int        `json:"turns"`
	FilesCreated int        `json:"filesCreated"`
}

// TaskMetadata contains full metadata for a task
type TaskMetadata struct {
	TaskID       string     `json:"taskId"`
	ParentTaskID string     `json:"parentTaskId,omitempty"` // Task this continues from
	ProjectID    string     `json:"projectId"`
	Task         string     `json:"task"`
	CreatedAt    time.Time  `json:"createdAt"`
	CompletedAt  *time.Time `json:"completedAt,omitempty"`
	Status       TaskStatus `json:"status"`
	Turns        int        `json:"turns"`
	Duration     string     `json:"duration"`
	FilesCreated []string   `json:"filesCreated"`
	Error        string     `json:"error,omitempty"`
}

// TaskHistory contains the list of all tasks for a project
type TaskHistory struct {
	ProjectID string        `json:"projectId"`
	Tasks     []TaskSummary `json:"tasks"`
}

// HistoryManager manages task history for projects
type HistoryManager struct {
	mu sync.RWMutex
}

// NewHistoryManager creates a new history manager
func NewHistoryManager() *HistoryManager {
	return &HistoryManager{}
}

// getTasksDir returns the .tasks directory path for a project
func (h *HistoryManager) getTasksDir(projectDir string) string {
	return filepath.Join(projectDir, ".tasks")
}

// getTaskDir returns the directory for a specific task
func (h *HistoryManager) getTaskDir(projectDir, taskID string) string {
	return filepath.Join(h.getTasksDir(projectDir), taskID)
}

// getHistoryPath returns the path to history.json
func (h *HistoryManager) getHistoryPath(projectDir string) string {
	return filepath.Join(h.getTasksDir(projectDir), "history.json")
}

// SaveTask saves a task's conversation and metadata
func (h *HistoryManager) SaveTask(projectDir string, meta *TaskMetadata, messages []*message.Message) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Create task directory
	taskDir := h.getTaskDir(projectDir, meta.TaskID)
	if err := os.MkdirAll(taskDir, 0755); err != nil {
		return fmt.Errorf("failed to create task directory: %w", err)
	}

	// Save metadata
	metaPath := filepath.Join(taskDir, "metadata.json")
	metaData, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}
	if err := os.WriteFile(metaPath, metaData, 0644); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	// Save conversation
	convPath := filepath.Join(taskDir, "conversation.json")
	convData, err := json.MarshalIndent(messages, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal conversation: %w", err)
	}
	if err := os.WriteFile(convPath, convData, 0644); err != nil {
		return fmt.Errorf("failed to write conversation: %w", err)
	}

	// Update history index
	return h.updateHistoryIndex(projectDir, meta)
}

// updateHistoryIndex updates the history.json index file
func (h *HistoryManager) updateHistoryIndex(projectDir string, meta *TaskMetadata) error {
	historyPath := h.getHistoryPath(projectDir)

	// Load existing history
	history := &TaskHistory{
		ProjectID: filepath.Base(projectDir),
		Tasks:     make([]TaskSummary, 0),
	}

	data, err := os.ReadFile(historyPath)
	if err == nil {
		json.Unmarshal(data, history)
	}

	// Create summary from metadata
	summary := TaskSummary{
		TaskID:       meta.TaskID,
		Summary:      truncateTask(meta.Task, 50),
		CreatedAt:    meta.CreatedAt,
		CompletedAt:  meta.CompletedAt,
		Status:       meta.Status,
		Turns:        meta.Turns,
		FilesCreated: len(meta.FilesCreated),
	}

	// Check if task already exists, update it
	found := false
	for i, t := range history.Tasks {
		if t.TaskID == meta.TaskID {
			history.Tasks[i] = summary
			found = true
			break
		}
	}

	// Add new task if not found
	if !found {
		history.Tasks = append(history.Tasks, summary)
	}

	// Sort by created time (newest first)
	sort.Slice(history.Tasks, func(i, j int) bool {
		return history.Tasks[i].CreatedAt.After(history.Tasks[j].CreatedAt)
	})

	// Save history
	historyData, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal history: %w", err)
	}

	return os.WriteFile(historyPath, historyData, 0644)
}

// GetTaskList returns the list of tasks for a project
func (h *HistoryManager) GetTaskList(projectDir string) (*TaskHistory, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	historyPath := h.getHistoryPath(projectDir)

	history := &TaskHistory{
		ProjectID: filepath.Base(projectDir),
		Tasks:     make([]TaskSummary, 0),
	}

	data, err := os.ReadFile(historyPath)
	if err != nil {
		if os.IsNotExist(err) {
			return history, nil // Empty history is OK
		}
		return nil, fmt.Errorf("failed to read history: %w", err)
	}

	if err := json.Unmarshal(data, history); err != nil {
		return nil, fmt.Errorf("failed to parse history: %w", err)
	}

	return history, nil
}

// GetTask returns the full task data including conversation
func (h *HistoryManager) GetTask(projectDir, taskID string) (*TaskMetadata, []*message.Message, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	taskDir := h.getTaskDir(projectDir, taskID)

	// Read metadata
	metaPath := filepath.Join(taskDir, "metadata.json")
	metaData, err := os.ReadFile(metaPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read metadata: %w", err)
	}

	var meta TaskMetadata
	if err := json.Unmarshal(metaData, &meta); err != nil {
		return nil, nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	// Read conversation
	convPath := filepath.Join(taskDir, "conversation.json")
	convData, err := os.ReadFile(convPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read conversation: %w", err)
	}

	var messages []*message.Message
	if err := json.Unmarshal(convData, &messages); err != nil {
		return nil, nil, fmt.Errorf("failed to parse conversation: %w", err)
	}

	return &meta, messages, nil
}

// GenerateTaskID creates a unique task ID
func GenerateTaskID() string {
	return fmt.Sprintf("task_%d", time.Now().UnixNano())
}

// truncateTask truncates a task description for display
func truncateTask(task string, maxLen int) string {
	if len(task) <= maxLen {
		return task
	}
	return task[:maxLen-3] + "..."
}
