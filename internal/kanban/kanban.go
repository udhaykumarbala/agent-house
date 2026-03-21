package kanban

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ─── Status & Priority ─────────────────────────────────

type Status string

const (
	StatusBacklog    Status = "backlog"
	StatusInProgress Status = "in_progress"
	StatusInReview   Status = "in_review"
	StatusDone       Status = "done"
	StatusBlocked    Status = "blocked"
)

type Priority string

const (
	PriorityCritical Priority = "critical"
	PriorityHigh     Priority = "high"
	PriorityMedium   Priority = "medium"
	PriorityLow      Priority = "low"
)

type ReviewStatus string

const (
	ReviewPending  ReviewStatus = "pending"
	ReviewApproved ReviewStatus = "approved"
	ReviewRejected ReviewStatus = "rejected"
)

// ─── Task ───────────────────────────────────────────────

type Task struct {
	ID             string       `json:"id"`
	ProjectID      string       `json:"project_id"`
	SubTaskRef     string       `json:"subtask_ref,omitempty"`
	PhaseIndex     int          `json:"phase_index"`
	PhaseName      string       `json:"phase_name,omitempty"`
	Title          string       `json:"title"`
	Description    string       `json:"description,omitempty"`
	AssignedAgent  string       `json:"assigned_agent"`
	Status         Status       `json:"status"`
	Priority       Priority     `json:"priority"`
	Source         string       `json:"source"` // "auto" or "manual"
	ReviewStatus   ReviewStatus `json:"review_status"`
	ReviewFeedback string       `json:"review_feedback,omitempty"`
	Dependencies   []string     `json:"dependencies,omitempty"`
	Outputs        []Output     `json:"outputs,omitempty"`
	Criteria       []Criterion  `json:"criteria,omitempty"`
	Effort         string       `json:"effort,omitempty"` // small, medium, large
	Progress       int          `json:"progress"`
	Position       int          `json:"position"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	StartedAt      *time.Time   `json:"started_at,omitempty"`
	CompletedAt    *time.Time   `json:"completed_at,omitempty"`
}

type Output struct {
	FilePath  string    `json:"file_path"`
	FileName  string    `json:"file_name"`
	FileType  string    `json:"file_type"`
	Size      int64     `json:"size"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

type Criterion struct {
	Text    string `json:"text"`
	Checked bool   `json:"checked"`
}

// ─── ID Generation ──────────────────────────────────────

func GenerateTaskID() string {
	return fmt.Sprintf("kt_%d", time.Now().UnixNano())
}

// ─── Storage ────────────────────────────────────────────

var mu sync.Mutex

func tasksPath(projectDir string) string {
	return filepath.Join(projectDir, ".tasks", "kanban.json")
}

func ensureDir(projectDir string) error {
	return os.MkdirAll(filepath.Join(projectDir, ".tasks"), 0755)
}

func LoadTasks(projectDir string) ([]Task, error) {
	mu.Lock()
	defer mu.Unlock()
	return loadTasksUnsafe(projectDir)
}

func loadTasksUnsafe(projectDir string) ([]Task, error) {
	data, err := os.ReadFile(tasksPath(projectDir))
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil
		}
		return nil, err
	}
	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func SaveTasks(projectDir string, tasks []Task) error {
	if err := ensureDir(projectDir); err != nil {
		return err
	}
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(tasksPath(projectDir), data, 0644)
}

func saveTasksUnsafe(projectDir string, tasks []Task) error {
	if err := ensureDir(projectDir); err != nil {
		return err
	}
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(tasksPath(projectDir), data, 0644)
}

// CreateTask adds a new task and persists.
func CreateTask(projectDir string, t Task) (*Task, error) {
	mu.Lock()
	defer mu.Unlock()

	tasks, _ := loadTasksUnsafe(projectDir)

	if t.ID == "" {
		t.ID = GenerateTaskID()
	}
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	if t.Status == "" {
		t.Status = StatusBacklog
	}
	if t.Priority == "" {
		t.Priority = PriorityMedium
	}
	if t.ReviewStatus == "" {
		t.ReviewStatus = ReviewPending
	}
	if t.Source == "" {
		t.Source = "manual"
	}

	// Set position to end of its column
	maxPos := 0
	for _, existing := range tasks {
		if existing.Status == t.Status && existing.Position >= maxPos {
			maxPos = existing.Position + 1
		}
	}
	t.Position = maxPos

	tasks = append(tasks, t)
	if err := saveTasksUnsafe(projectDir, tasks); err != nil {
		return nil, err
	}
	return &t, nil
}

// GetTask returns a single task by ID.
func GetTask(projectDir, taskID string) (*Task, error) {
	tasks, err := LoadTasks(projectDir)
	if err != nil {
		return nil, err
	}
	for i := range tasks {
		if tasks[i].ID == taskID {
			return &tasks[i], nil
		}
	}
	return nil, fmt.Errorf("task not found: %s", taskID)
}

// UpdateTask applies partial updates to a task.
func UpdateTask(projectDir, taskID string, updates map[string]interface{}) (*Task, error) {
	mu.Lock()
	defer mu.Unlock()

	tasks, err := loadTasksUnsafe(projectDir)
	if err != nil {
		return nil, err
	}

	var target *Task
	for i := range tasks {
		if tasks[i].ID == taskID {
			target = &tasks[i]
			break
		}
	}
	if target == nil {
		return nil, fmt.Errorf("task not found: %s", taskID)
	}

	now := time.Now()
	target.UpdatedAt = now

	if v, ok := updates["status"]; ok {
		newStatus := Status(v.(string))
		target.Status = newStatus
		if newStatus == StatusInProgress && target.StartedAt == nil {
			target.StartedAt = &now
		}
		if newStatus == StatusDone && target.CompletedAt == nil {
			target.CompletedAt = &now
			target.Progress = 100
		}
	}
	if v, ok := updates["priority"]; ok {
		target.Priority = Priority(v.(string))
	}
	if v, ok := updates["assigned_agent"]; ok {
		target.AssignedAgent = v.(string)
	}
	if v, ok := updates["position"]; ok {
		switch p := v.(type) {
		case float64:
			target.Position = int(p)
		case int:
			target.Position = p
		}
	}
	if v, ok := updates["progress"]; ok {
		switch p := v.(type) {
		case float64:
			target.Progress = int(p)
		case int:
			target.Progress = p
		}
	}
	if v, ok := updates["title"]; ok {
		target.Title = v.(string)
	}
	if v, ok := updates["review_status"]; ok {
		target.ReviewStatus = ReviewStatus(v.(string))
	}
	if v, ok := updates["review_feedback"]; ok {
		target.ReviewFeedback = v.(string)
	}

	if err := saveTasksUnsafe(projectDir, tasks); err != nil {
		return nil, err
	}
	return target, nil
}

// ReviewTask applies a review decision to a task.
func ReviewTask(projectDir, taskID, action, feedback string, criteriaUpdates []map[string]interface{}) (*Task, error) {
	mu.Lock()
	defer mu.Unlock()

	tasks, err := loadTasksUnsafe(projectDir)
	if err != nil {
		return nil, err
	}

	var target *Task
	for i := range tasks {
		if tasks[i].ID == taskID {
			target = &tasks[i]
			break
		}
	}
	if target == nil {
		return nil, fmt.Errorf("task not found: %s", taskID)
	}

	now := time.Now()
	target.UpdatedAt = now

	if action == "approved" {
		target.ReviewStatus = ReviewApproved
		target.Status = StatusDone
		target.CompletedAt = &now
		target.Progress = 100
	} else if action == "rejected" {
		target.ReviewStatus = ReviewRejected
		target.ReviewFeedback = feedback
		target.Status = StatusInProgress
	}

	// Update criteria checkboxes
	for _, cu := range criteriaUpdates {
		idx, ok1 := cu["index"]
		checked, ok2 := cu["checked"]
		if ok1 && ok2 {
			var i int
			switch v := idx.(type) {
			case float64:
				i = int(v)
			case int:
				i = v
			}
			if i >= 0 && i < len(target.Criteria) {
				target.Criteria[i].Checked = checked.(bool)
			}
		}
	}

	if err := saveTasksUnsafe(projectDir, tasks); err != nil {
		return nil, err
	}
	return target, nil
}

// FilterTasks returns tasks matching the given filters.
func FilterTasks(projectDir string, status Status, agentRole string, phaseIdx int, priority Priority) ([]Task, error) {
	tasks, err := LoadTasks(projectDir)
	if err != nil {
		return nil, err
	}

	var result []Task
	for _, t := range tasks {
		if status != "" && t.Status != status {
			continue
		}
		if agentRole != "" && t.AssignedAgent != agentRole {
			continue
		}
		if phaseIdx > 0 && t.PhaseIndex != phaseIdx {
			continue
		}
		if priority != "" && t.Priority != priority {
			continue
		}
		result = append(result, t)
	}
	return result, nil
}

// GetColumnCounts returns task counts per status column.
func GetColumnCounts(tasks []Task) map[string]int {
	counts := map[string]int{
		"backlog":     0,
		"in_progress": 0,
		"in_review":   0,
		"done":        0,
	}
	for _, t := range tasks {
		s := string(t.Status)
		if s == "blocked" {
			s = "backlog" // Blocked tasks count in backlog column
		}
		counts[s]++
	}
	return counts
}

// SyncFromDevPlan creates kanban tasks from a development plan's subtasks.
// Only creates tasks that don't already exist (based on subtask_ref).
func SyncFromDevPlan(projectDir, projectID string, phases []DevPlanPhase) (int, error) {
	mu.Lock()
	defer mu.Unlock()

	tasks, _ := loadTasksUnsafe(projectDir)

	// Index existing refs
	existingRefs := make(map[string]bool)
	for _, t := range tasks {
		if t.SubTaskRef != "" {
			existingRefs[t.SubTaskRef] = true
		}
	}

	created := 0
	for _, phase := range phases {
		for _, st := range phase.SubTasks {
			ref := fmt.Sprintf("p%d_st%s", phase.Index, st.ID)
			if existingRefs[ref] {
				continue
			}

			agent := ""
			if len(st.AssignedAgents) > 0 {
				agent = st.AssignedAgents[0]
			}

			var criteria []Criterion
			for _, c := range st.CompletionCriteria {
				criteria = append(criteria, Criterion{Text: c, Checked: false})
			}

			t := Task{
				ID:            GenerateTaskID(),
				ProjectID:     projectID,
				SubTaskRef:    ref,
				PhaseIndex:    phase.Index,
				PhaseName:     phase.Name,
				Title:         st.Title,
				Description:   st.Description,
				AssignedAgent: agent,
				Status:        StatusBacklog,
				Priority:      PriorityMedium,
				Source:         "auto",
				ReviewStatus:  ReviewPending,
				Criteria:      criteria,
				Effort:        "",
				Progress:      0,
				Position:      created,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}
			tasks = append(tasks, t)
			created++
		}
	}

	if created > 0 {
		if err := saveTasksUnsafe(projectDir, tasks); err != nil {
			return 0, err
		}
	}
	return created, nil
}

// DevPlanPhase is a minimal representation for sync purposes.
type DevPlanPhase struct {
	Index    int
	Name     string
	SubTasks []DevPlanSubTask
}

type DevPlanSubTask struct {
	ID                 string
	Title              string
	Description        string
	AssignedAgents     []string
	CompletionCriteria []string
}
