package orchestrator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"pty-claude-test/internal/agent"
)

// devPlanMu serializes all load-modify-save operations on development plans
// to prevent concurrent tasks from overwriting each other's changes.
var devPlanMu sync.Mutex

// PhaseStatus represents the status of a development phase
type PhaseStatus string

const (
	PhaseStatusPending       PhaseStatus = "pending"
	PhaseStatusInProgress    PhaseStatus = "in_progress"
	PhaseStatusCompleted     PhaseStatus = "completed"
	PhaseStatusNeedsRevision PhaseStatus = "needs_revision"
)

// QAStatus represents the QA review status
type QAStatus string

const (
	QAStatusPending  QAStatus = "pending"
	QAStatusApproved QAStatus = "approved"
	QAStatusRejected QAStatus = "rejected"
)

// DevelopmentPlan represents the multi-phase development plan
type DevelopmentPlan struct {
	TaskID     string             `json:"task_id"`
	Phases     []DevelopmentPhase `json:"phases"`
	CreatedBy  agent.Role         `json:"created_by"`
	ApprovedBy agent.Role         `json:"approved_by"`
	CreatedAt  time.Time          `json:"created_at"`
	ApprovedAt *time.Time         `json:"approved_at,omitempty"`
}

// DevelopmentPhase represents one phase of development
type DevelopmentPhase struct {
	Index       int         `json:"index"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	SubTasks    []SubTask   `json:"subtasks"`
	Status      PhaseStatus `json:"status"`
	QAStatus    QAStatus    `json:"qa_status"`
	QAFeedback  string      `json:"qa_feedback,omitempty"`
	Iteration   int         `json:"iteration"`
	StartedAt   *time.Time  `json:"started_at,omitempty"`
	CompletedAt *time.Time  `json:"completed_at,omitempty"`
}

// LoadDevelopmentPlan loads the development plan from .plans/development-plan.json
func LoadDevelopmentPlan(projectDir string) (*DevelopmentPlan, error) {
	planPath := filepath.Join(projectDir, ".plans", "development-plan.json")

	data, err := os.ReadFile(planPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("development plan not found at %s", planPath)
		}
		return nil, fmt.Errorf("failed to read development plan: %w", err)
	}

	var plan DevelopmentPlan
	if err := json.Unmarshal(data, &plan); err != nil {
		return nil, fmt.Errorf("failed to parse development plan: %w", err)
	}

	// Initialize statuses if not set
	for i := range plan.Phases {
		if plan.Phases[i].Status == "" {
			plan.Phases[i].Status = PhaseStatusPending
		}
		if plan.Phases[i].QAStatus == "" {
			plan.Phases[i].QAStatus = QAStatusPending
		}
		if plan.Phases[i].Iteration == 0 {
			plan.Phases[i].Iteration = 1
		}
	}

	return &plan, nil
}

// SaveDevelopmentPlan saves the development plan to .plans/development-plan.json
func SaveDevelopmentPlan(projectDir string, plan *DevelopmentPlan) error {
	planPath := filepath.Join(projectDir, ".plans", "development-plan.json")

	// Ensure .plans directory exists
	if err := os.MkdirAll(filepath.Join(projectDir, ".plans"), 0755); err != nil {
		return fmt.Errorf("failed to create .plans directory: %w", err)
	}

	data, err := json.MarshalIndent(plan, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal development plan: %w", err)
	}

	if err := os.WriteFile(planPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write development plan: %w", err)
	}

	return nil
}

// GetPhase returns a pointer to the phase at the given index
func (p *DevelopmentPlan) GetPhase(index int) *DevelopmentPhase {
	if index < 1 || index > len(p.Phases) {
		return nil
	}
	return &p.Phases[index-1]
}

// UpdatePhaseStatus updates the status of a phase and saves the plan
func (p *DevelopmentPlan) UpdatePhaseStatus(projectDir string, phaseIndex int, status PhaseStatus) error {
	devPlanMu.Lock()
	defer devPlanMu.Unlock()

	phase := p.GetPhase(phaseIndex)
	if phase == nil {
		return fmt.Errorf("invalid phase index: %d", phaseIndex)
	}

	phase.Status = status

	now := time.Now()
	if status == PhaseStatusInProgress && phase.StartedAt == nil {
		phase.StartedAt = &now
	}
	if status == PhaseStatusCompleted && phase.CompletedAt == nil {
		phase.CompletedAt = &now
	}

	return SaveDevelopmentPlan(projectDir, p)
}

// UpdatePhaseQAStatus updates the QA status of a phase and saves the plan
func (p *DevelopmentPlan) UpdatePhaseQAStatus(projectDir string, phaseIndex int, qaStatus QAStatus, feedback string) error {
	devPlanMu.Lock()
	defer devPlanMu.Unlock()

	phase := p.GetPhase(phaseIndex)
	if phase == nil {
		return fmt.Errorf("invalid phase index: %d", phaseIndex)
	}

	phase.QAStatus = qaStatus
	phase.QAFeedback = feedback

	if qaStatus == QAStatusRejected {
		phase.Status = PhaseStatusNeedsRevision
	}

	return SaveDevelopmentPlan(projectDir, p)
}

// IncrementIteration increments the iteration counter for a phase
func (p *DevelopmentPlan) IncrementIteration(projectDir string, phaseIndex int) error {
	devPlanMu.Lock()
	defer devPlanMu.Unlock()

	phase := p.GetPhase(phaseIndex)
	if phase == nil {
		return fmt.Errorf("invalid phase index: %d", phaseIndex)
	}

	phase.Iteration++
	phase.QAStatus = QAStatusPending
	phase.Status = PhaseStatusInProgress

	return SaveDevelopmentPlan(projectDir, p)
}

// GetCurrentPhase returns the first phase that is not completed
func (p *DevelopmentPlan) GetCurrentPhase() *DevelopmentPhase {
	for i := range p.Phases {
		if p.Phases[i].Status != PhaseStatusCompleted {
			return &p.Phases[i]
		}
	}
	return nil
}

// IsComplete returns true if all phases are completed and approved
func (p *DevelopmentPlan) IsComplete() bool {
	for _, phase := range p.Phases {
		if phase.Status != PhaseStatusCompleted || phase.QAStatus != QAStatusApproved {
			return false
		}
	}
	return true
}

// GetTotalSubTasks returns the total number of subtasks across all phases
func (p *DevelopmentPlan) GetTotalSubTasks() int {
	total := 0
	for _, phase := range p.Phases {
		total += len(phase.SubTasks)
	}
	return total
}

// GetCompletedSubTasks returns the number of completed subtasks across all phases
func (p *DevelopmentPlan) GetCompletedSubTasks() int {
	completed := 0
	for _, phase := range p.Phases {
		for _, subtask := range phase.SubTasks {
			if subtask.Status == SubTaskStatusCompleted {
				completed++
			}
		}
	}
	return completed
}
