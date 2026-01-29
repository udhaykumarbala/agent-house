package agentask

import (
	"time"

	"pty-claude-test/internal/agent"
)

// AgentTaskType represents the type of task an agent is performing
type AgentTaskType string

const (
	TaskTypeResearch     AgentTaskType = "research"      // Research and discovery
	TaskTypePlanning     AgentTaskType = "planning"      // Planning and specification
	TaskTypeCoding       AgentTaskType = "coding"        // Implementation
	TaskTypeReview       AgentTaskType = "review"        // Code review
	TaskTypeDiscussion   AgentTaskType = "discussion"    // Discussion and decision
	TaskTypeTemplateSelection AgentTaskType = "template_selection" // Template selection
)

// AgentTaskStatus represents the current state of an agent task
type AgentTaskStatus string

const (
	StatusPending     AgentTaskStatus = "pending"      // Not yet started
	StatusInProgress  AgentTaskStatus = "in_progress"  // Currently executing
	StatusCompleted   AgentTaskStatus = "completed"    // Successfully completed
	StatusFailed      AgentTaskStatus = "failed"       // Failed with error
	StatusCancelled   AgentTaskStatus = "cancelled"    // Cancelled by user or system
)

// AgentState represents the current activity state of an agent
type AgentState string

const (
	StateIdle      AgentState = "idle"       // Not doing anything
	StateWorking   AgentState = "working"    // Currently executing a task
	StateWaiting   AgentState = "waiting"    // Waiting for dependencies
	StateCompleted AgentState = "completed"  // Finished all tasks
	StateError     AgentState = "error"      // Encountered an error
)

// OutputType represents the type of output an agent produces
type OutputType string

const (
	OutputResearchFindings OutputType = "research_findings" // Research document
	OutputSpecification    OutputType = "specification"     // Technical spec
	OutputCodeFile         OutputType = "code_file"         // Code implementation
	OutputReviewComments   OutputType = "review_comments"   // Review feedback
	OutputDecision         OutputType = "decision"          // Decision document
	OutputTemplate         OutputType = "template"          // Template selection
)

// AgentTask represents a task assigned to a specific agent
type AgentTask struct {
	ID           string          `json:"id"`
	AgentRole    agent.Role      `json:"agent_role"`
	TaskType     AgentTaskType   `json:"task_type"`
	Status       AgentTaskStatus `json:"status"`
	Title        string          `json:"title"`        // Brief task description
	Description  string          `json:"description"`  // Detailed task description
	Progress     int             `json:"progress"`     // 0-100 percentage
	Outputs      []AgentOutput   `json:"outputs"`      // Produced outputs
	Error        string          `json:"error,omitempty"` // Error message if failed
	ProjectID    string          `json:"project_id"`
	ParentTaskID string          `json:"parent_task_id"` // Parent orchestrator task ID
	AssignedAt   time.Time       `json:"assigned_at"`
	StartedAt    *time.Time      `json:"started_at,omitempty"`
	CompletedAt  *time.Time      `json:"completed_at,omitempty"`
	Duration     time.Duration   `json:"duration,omitempty"` // Execution duration
}

// AgentOutput represents a structured output from an agent
type AgentOutput struct {
	ID        string     `json:"id"`
	Type      OutputType `json:"type"`
	Title     string     `json:"title"`     // Output title
	FilePath  string     `json:"file_path"` // Path to the output file
	Format    string     `json:"format"`    // markdown, code, json, etc.
	Size      int64      `json:"size"`      // File size in bytes
	CreatedAt time.Time  `json:"created_at"`
}

// AgentStatus represents the current status of an agent
type AgentStatus struct {
	Role            agent.Role    `json:"role"`
	State           AgentState    `json:"state"`
	CurrentTaskID   string        `json:"current_task_id,omitempty"`
	CurrentTaskType AgentTaskType `json:"current_task_type,omitempty"`
	TasksCompleted  int           `json:"tasks_completed"`
	TasksFailed     int           `json:"tasks_failed"`
	LastActiveAt    *time.Time    `json:"last_active_at,omitempty"`
}

// NewAgentTask creates a new agent task
func NewAgentTask(role agent.Role, taskType AgentTaskType, title, description, projectID, parentTaskID string) *AgentTask {
	return &AgentTask{
		ID:           generateTaskID(),
		AgentRole:    role,
		TaskType:     taskType,
		Status:       StatusPending,
		Title:        title,
		Description:  description,
		Progress:     0,
		Outputs:      make([]AgentOutput, 0),
		ProjectID:    projectID,
		ParentTaskID: parentTaskID,
		AssignedAt:   time.Now(),
	}
}

// Start marks the task as started
func (at *AgentTask) Start() {
	now := time.Now()
	at.StartedAt = &now
	at.Status = StatusInProgress
}

// Complete marks the task as completed
func (at *AgentTask) Complete() {
	now := time.Now()
	at.CompletedAt = &now
	at.Status = StatusCompleted
	at.Progress = 100
	if at.StartedAt != nil {
		at.Duration = now.Sub(*at.StartedAt)
	}
}

// Fail marks the task as failed
func (at *AgentTask) Fail(err error) {
	now := time.Now()
	at.CompletedAt = &now
	at.Status = StatusFailed
	at.Error = err.Error()
	if at.StartedAt != nil {
		at.Duration = now.Sub(*at.StartedAt)
	}
}

// AddOutput adds an output to the task
func (at *AgentTask) AddOutput(output AgentOutput) {
	at.Outputs = append(at.Outputs, output)
}

// UpdateProgress updates the task progress
func (at *AgentTask) UpdateProgress(progress int) {
	if progress < 0 {
		progress = 0
	}
	if progress > 100 {
		progress = 100
	}
	at.Progress = progress
}

// generateTaskID generates a unique task ID
func generateTaskID() string {
	return "at_" + time.Now().Format("20060102150405") + "_" + randomString(6)
}

// randomString generates a random alphanumeric string
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
		time.Sleep(time.Nanosecond)
	}
	return string(b)
}
