package capability

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

type Milestone struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	DueDate     string `json:"due_date"` // YYYY-MM-DD
	Status      string `json:"status"`   // pending, in_progress, completed, blocked
	PctComplete int    `json:"pct_complete"`
	Owner       string `json:"owner,omitempty"`
}

type Slip struct {
	Milestone Milestone `json:"milestone"`
	DaysLate  int       `json:"days_late"`
	Severity  string    `json:"severity"` // minor (1-2d), moderate (3-7d), critical (>7d)
}

type Schedule struct {
	mu   sync.RWMutex
	path string
}

func NewSchedule(dataRoot, scope string) (*Schedule, error) {
	dir := filepath.Join(dataRoot, scope)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return &Schedule{path: filepath.Join(dir, "milestones.json")}, nil
}

func (s *Schedule) load() ([]Milestone, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Milestone{}, nil
		}
		return nil, err
	}
	var out []Milestone
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("parse %s: %w", s.path, err)
	}
	return out, nil
}

func (s *Schedule) save(in []Milestone) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}

func (s *Schedule) List() ([]Milestone, error) { return s.load() }

func (s *Schedule) Store(m Milestone) (Milestone, error) {
	if m.ID == "" {
		m.ID = fmt.Sprintf("ms_%d", time.Now().UnixNano())
	}
	if m.Status == "" {
		m.Status = "pending"
	}
	cur, err := s.load()
	if err != nil {
		return m, err
	}
	replaced := false
	for i := range cur {
		if cur[i].ID == m.ID {
			cur[i] = m
			replaced = true
			break
		}
	}
	if !replaced {
		cur = append(cur, m)
	}
	if err := s.save(cur); err != nil {
		return m, err
	}
	return m, nil
}

// Slips returns milestones past their due date that are not yet completed.
// Severity buckets follow the docstring on the Slip type.
func (s *Schedule) Slips() ([]Slip, error) {
	ms, err := s.load()
	if err != nil {
		return nil, err
	}
	today := time.Now().UTC().Truncate(24 * time.Hour)
	out := []Slip{}
	for _, m := range ms {
		if m.Status == "completed" {
			continue
		}
		due, err := time.Parse("2006-01-02", m.DueDate)
		if err != nil {
			continue
		}
		days := int(today.Sub(due).Hours() / 24)
		if days < 1 {
			continue
		}
		sev := "minor"
		switch {
		case days > 7:
			sev = "critical"
		case days >= 3:
			sev = "moderate"
		}
		out = append(out, Slip{Milestone: m, DaysLate: days, Severity: sev})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].DaysLate > out[j].DaysLate })
	return out, nil
}
