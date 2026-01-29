package orchestrator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"pty-claude-test/internal/agent"
)

// QAReview represents a QA review of a development phase
type QAReview struct {
	ID             string     `json:"id"`
	TaskID         string     `json:"task_id"`
	PhaseIndex     int        `json:"phase_index"`
	PhaseName      string     `json:"phase_name"`
	Reviewer       agent.Role `json:"reviewer"`
	Status         QAStatus   `json:"status"`
	Feedback       string     `json:"feedback"`
	FailedCriteria []string   `json:"failed_criteria,omitempty"`
	ReviewedAt     time.Time  `json:"reviewed_at"`
	Iteration      int        `json:"iteration"`
}

// QAReviewHistory contains all QA reviews for a task
type QAReviewHistory struct {
	TaskID  string     `json:"task_id"`
	Reviews []QAReview `json:"reviews"`
}

// LoadQAReviewHistory loads QA review history from .plans/qa-reviews.json
func LoadQAReviewHistory(projectDir string) (*QAReviewHistory, error) {
	reviewPath := filepath.Join(projectDir, ".plans", "qa-reviews.json")

	// If file doesn't exist, return empty history
	if _, err := os.Stat(reviewPath); os.IsNotExist(err) {
		return &QAReviewHistory{
			Reviews: make([]QAReview, 0),
		}, nil
	}

	data, err := os.ReadFile(reviewPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read QA review history: %w", err)
	}

	var history QAReviewHistory
	if err := json.Unmarshal(data, &history); err != nil {
		return nil, fmt.Errorf("failed to parse QA review history: %w", err)
	}

	return &history, nil
}

// SaveQAReviewHistory saves QA review history to .plans/qa-reviews.json
func SaveQAReviewHistory(projectDir string, history *QAReviewHistory) error {
	reviewPath := filepath.Join(projectDir, ".plans", "qa-reviews.json")

	// Ensure .plans directory exists
	if err := os.MkdirAll(filepath.Join(projectDir, ".plans"), 0755); err != nil {
		return fmt.Errorf("failed to create .plans directory: %w", err)
	}

	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal QA review history: %w", err)
	}

	if err := os.WriteFile(reviewPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write QA review history: %w", err)
	}

	return nil
}

// AddReview adds a new QA review to the history
func (h *QAReviewHistory) AddReview(review QAReview) {
	h.Reviews = append(h.Reviews, review)
}

// GetReviewsForPhase returns all reviews for a specific phase
func (h *QAReviewHistory) GetReviewsForPhase(phaseIndex int) []QAReview {
	var reviews []QAReview
	for _, review := range h.Reviews {
		if review.PhaseIndex == phaseIndex {
			reviews = append(reviews, review)
		}
	}
	return reviews
}

// GetLatestReview returns the most recent review
func (h *QAReviewHistory) GetLatestReview() *QAReview {
	if len(h.Reviews) == 0 {
		return nil
	}
	return &h.Reviews[len(h.Reviews)-1]
}

// GetLatestReviewForPhase returns the most recent review for a specific phase
func (h *QAReviewHistory) GetLatestReviewForPhase(phaseIndex int) *QAReview {
	var latest *QAReview
	for i := range h.Reviews {
		review := &h.Reviews[i]
		if review.PhaseIndex == phaseIndex {
			if latest == nil || review.ReviewedAt.After(latest.ReviewedAt) {
				latest = review
			}
		}
	}
	return latest
}

// CreateQAReview creates a new QA review record
func CreateQAReview(taskID string, phaseIndex int, phaseName string, reviewer agent.Role, status QAStatus, feedback string, failedCriteria []string, iteration int) QAReview {
	return QAReview{
		ID:             fmt.Sprintf("qa_%s_p%d_i%d_%d", taskID, phaseIndex, iteration, time.Now().Unix()),
		TaskID:         taskID,
		PhaseIndex:     phaseIndex,
		PhaseName:      phaseName,
		Reviewer:       reviewer,
		Status:         status,
		Feedback:       feedback,
		FailedCriteria: failedCriteria,
		ReviewedAt:     time.Now(),
		Iteration:      iteration,
	}
}

// RecordQAReview records a QA review for a phase
func RecordQAReview(projectDir string, review QAReview) error {
	history, err := LoadQAReviewHistory(projectDir)
	if err != nil {
		return err
	}

	history.AddReview(review)

	return SaveQAReviewHistory(projectDir, history)
}
