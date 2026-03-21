package checkpoint

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ─── Types ─────────────────────────────────────────────

type DecisionAction string

const (
	DecisionApproved   DecisionAction = "approved"
	DecisionRejected   DecisionAction = "rejected"
	DecisionOverridden DecisionAction = "overridden"
)

// Checkpoint types
const (
	TypeTemplateApproval = "template_approval"
	TypePlanApproval     = "plan_approval"
	TypePhaseGate        = "phase_gate"
	TypeFinalAcceptance  = "final_acceptance"
	TypeResearchReview   = "research_review"
	TypeSpecReview       = "spec_review"
	TypePreQAReview      = "pre_qa_review"
)

// Decision records a human or auto-delegated decision at a checkpoint.
type Decision struct {
	ID             string         `json:"id"`
	ProjectID      string         `json:"project_id"`
	TaskID         string         `json:"task_id"`
	CheckpointType string         `json:"checkpoint_type"`
	PhaseIndex     int            `json:"phase_index,omitempty"`
	Action         DecisionAction `json:"action"`
	Feedback       string         `json:"feedback,omitempty"`
	OverrideData   string         `json:"override_data,omitempty"`
	DecidedBy      string         `json:"decided_by"` // "user" or "ceo_auto"
	DecidedAt      time.Time      `json:"decided_at"`
}

// Checkpoint tracks a pending or resolved approval point.
type Checkpoint struct {
	ID              string     `json:"id"`
	ProjectID       string     `json:"project_id"`
	TaskID          string     `json:"task_id"`
	Type            string     `json:"type"`
	PhaseIndex      int        `json:"phase_index,omitempty"`
	Status          string     `json:"status"` // "pending", "resolved"
	ArtifactPath    string     `json:"artifact_path"`
	ArtifactSummary string     `json:"artifact_summary"`
	Decision        *Decision  `json:"decision,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	ResolvedAt      *time.Time `json:"resolved_at,omitempty"`
}

// WorkflowSettings controls which checkpoints are active and auto-delegation.
type WorkflowSettings struct {
	ProjectID               string `json:"project_id"`
	RequireTemplateApproval bool   `json:"require_template_approval"`
	RequirePlanApproval     bool   `json:"require_plan_approval"`
	RequirePhaseGate        bool   `json:"require_phase_gate"`
	RequireFinalAcceptance  bool   `json:"require_final_acceptance"`
	RequireResearchReview   bool   `json:"require_research_review"`
	RequireSpecReview       bool   `json:"require_spec_review"`
	RequirePreQAReview      bool   `json:"require_pre_qa_review"`
	AutoDelegateMinutes     int    `json:"auto_delegate_minutes"`
}

// ─── Defaults ──────────────────────────────────────────

func DefaultSettings(projectID string) WorkflowSettings {
	return WorkflowSettings{
		ProjectID:               projectID,
		RequireTemplateApproval: true,
		RequirePlanApproval:     true,
		RequirePhaseGate:        true,
		RequireFinalAcceptance:  true,
		RequireResearchReview:   false,
		RequireSpecReview:       false,
		RequirePreQAReview:      false,
		AutoDelegateMinutes:     0,
	}
}

// IsCheckpointEnabled returns whether a checkpoint type is enabled.
func (s *WorkflowSettings) IsCheckpointEnabled(cpType string) bool {
	switch cpType {
	case TypeTemplateApproval:
		return s.RequireTemplateApproval
	case TypePlanApproval:
		return s.RequirePlanApproval
	case TypePhaseGate:
		return s.RequirePhaseGate
	case TypeFinalAcceptance:
		return s.RequireFinalAcceptance
	case TypeResearchReview:
		return s.RequireResearchReview
	case TypeSpecReview:
		return s.RequireSpecReview
	case TypePreQAReview:
		return s.RequirePreQAReview
	default:
		return false
	}
}

// ─── ID Generation ─────────────────────────────────────

func GenerateCheckpointID() string {
	return fmt.Sprintf("cp_%d", time.Now().UnixNano())
}

func GenerateDecisionID() string {
	return fmt.Sprintf("dec_%d", time.Now().UnixNano())
}

// ─── File-level mutexes ────────────────────────────────

var (
	checkpointMu sync.Mutex
	settingsMu   sync.Mutex
	decisionMu   sync.Mutex
)

func tasksDir(projectDir string) string {
	return filepath.Join(projectDir, ".tasks")
}

func ensureTasksDir(projectDir string) error {
	return os.MkdirAll(tasksDir(projectDir), 0755)
}

// ─── Checkpoint Storage ────────────────────────────────

func checkpointsPath(projectDir string) string {
	return filepath.Join(tasksDir(projectDir), "checkpoints.json")
}

func LoadCheckpoints(projectDir string) ([]Checkpoint, error) {
	checkpointMu.Lock()
	defer checkpointMu.Unlock()

	data, err := os.ReadFile(checkpointsPath(projectDir))
	if err != nil {
		if os.IsNotExist(err) {
			return []Checkpoint{}, nil
		}
		return nil, err
	}

	var cps []Checkpoint
	if err := json.Unmarshal(data, &cps); err != nil {
		return nil, err
	}
	return cps, nil
}

func SaveCheckpoints(projectDir string, cps []Checkpoint) error {
	if err := ensureTasksDir(projectDir); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cps, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(checkpointsPath(projectDir), data, 0644)
}

func GetCheckpoint(projectDir, cpID string) (*Checkpoint, error) {
	cps, err := LoadCheckpoints(projectDir)
	if err != nil {
		return nil, err
	}
	for i := range cps {
		if cps[i].ID == cpID {
			return &cps[i], nil
		}
	}
	return nil, fmt.Errorf("checkpoint not found: %s", cpID)
}

func GetPendingCheckpoint(projectDir string) (*Checkpoint, error) {
	cps, err := LoadCheckpoints(projectDir)
	if err != nil {
		return nil, err
	}
	for i := range cps {
		if cps[i].Status == "pending" {
			return &cps[i], nil
		}
	}
	return nil, nil
}

func CreateCheckpoint(projectDir string, cp Checkpoint) error {
	checkpointMu.Lock()
	defer checkpointMu.Unlock()

	data, err := os.ReadFile(checkpointsPath(projectDir))
	var cps []Checkpoint
	if err == nil {
		json.Unmarshal(data, &cps)
	}
	cps = append(cps, cp)

	if err := ensureTasksDir(projectDir); err != nil {
		return err
	}
	out, err := json.MarshalIndent(cps, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(checkpointsPath(projectDir), out, 0644)
}

func ResolveCheckpoint(projectDir, cpID string, dec Decision) (*Checkpoint, error) {
	checkpointMu.Lock()
	defer checkpointMu.Unlock()

	data, err := os.ReadFile(checkpointsPath(projectDir))
	if err != nil {
		return nil, fmt.Errorf("failed to read checkpoints: %w", err)
	}

	var cps []Checkpoint
	if err := json.Unmarshal(data, &cps); err != nil {
		return nil, err
	}

	var resolved *Checkpoint
	now := time.Now()
	for i := range cps {
		if cps[i].ID == cpID {
			cps[i].Status = "resolved"
			cps[i].Decision = &dec
			cps[i].ResolvedAt = &now
			resolved = &cps[i]
			break
		}
	}

	if resolved == nil {
		return nil, fmt.Errorf("checkpoint not found: %s", cpID)
	}

	out, err := json.MarshalIndent(cps, "", "  ")
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(checkpointsPath(projectDir), out, 0644); err != nil {
		return nil, err
	}

	return resolved, nil
}

// ─── Decision Storage ──────────────────────────────────

func decisionsPath(projectDir string) string {
	return filepath.Join(tasksDir(projectDir), "decisions.json")
}

func AddDecision(projectDir string, dec Decision) error {
	decisionMu.Lock()
	defer decisionMu.Unlock()

	data, err := os.ReadFile(decisionsPath(projectDir))
	var decs []Decision
	if err == nil {
		json.Unmarshal(data, &decs)
	}
	decs = append(decs, dec)

	if err := ensureTasksDir(projectDir); err != nil {
		return err
	}
	out, err := json.MarshalIndent(decs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(decisionsPath(projectDir), out, 0644)
}

func LoadDecisions(projectDir string) ([]Decision, error) {
	decisionMu.Lock()
	defer decisionMu.Unlock()

	data, err := os.ReadFile(decisionsPath(projectDir))
	if err != nil {
		if os.IsNotExist(err) {
			return []Decision{}, nil
		}
		return nil, err
	}
	var decs []Decision
	if err := json.Unmarshal(data, &decs); err != nil {
		return nil, err
	}
	return decs, nil
}

// ─── Settings Storage ──────────────────────────────────

func settingsPath(projectDir string) string {
	return filepath.Join(tasksDir(projectDir), "settings.json")
}

func LoadSettings(projectDir string) (*WorkflowSettings, error) {
	settingsMu.Lock()
	defer settingsMu.Unlock()

	data, err := os.ReadFile(settingsPath(projectDir))
	if err != nil {
		if os.IsNotExist(err) {
			s := DefaultSettings("")
			return &s, nil
		}
		return nil, err
	}

	var s WorkflowSettings
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func SaveSettings(projectDir string, s *WorkflowSettings) error {
	settingsMu.Lock()
	defer settingsMu.Unlock()

	if err := ensureTasksDir(projectDir); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(settingsPath(projectDir), data, 0644)
}
