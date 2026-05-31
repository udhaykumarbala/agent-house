// Package capability defines per-agent backing data + business logic so
// agents have something to actually *do* beyond producing text. Each
// capability persists state under data/<scope>/<area>.json so multi-tenant
// scoping ("atlas-site" vs "default") is a single directory boundary.
package capability

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// Applicant mirrors internal/email.Applicant but the capability layer owns
// its own copy so it does not depend on the email engine. Importing
// internal/email here would invert the dependency direction (email already
// imports broader package state).
type Applicant struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Email           string   `json:"email"`
	AppliedFor      string   `json:"applied_for"`
	ExperienceYears int      `json:"experience_years,omitempty"`
	KeySkills       []string `json:"key_skills,omitempty"`
	Certifications  []string `json:"certifications,omitempty"`
	ResumeEmailID   string   `json:"resume_email_id,omitempty"`
	ReceivedDate    string   `json:"received_date,omitempty"`
	Status          string   `json:"status"` // new, reviewed, shortlisted, rejected, hired
}

// JD is what Conductor decomposes a hiring request into.
type JD struct {
	Title      string   `json:"title"`
	MustHave   []string `json:"must_have"`
	NiceToHave []string `json:"nice_to_have"`
	MinYears   int      `json:"min_years"`
}

// ApplicantMatch is an Applicant scored against a JD with a human-readable
// breakdown of what matched and what didn't.
type ApplicantMatch struct {
	Applicant Applicant `json:"applicant"`
	Score     int       `json:"score"`
	Matched   []string  `json:"matched"`
	Missing   []string  `json:"missing"`
}

// HR is the per-scope applicant store + matcher. Safe for concurrent use.
type HR struct {
	mu   sync.RWMutex
	path string
}

// NewHR opens (creates if missing) the per-scope HR store.
func NewHR(dataRoot, scope string) (*HR, error) {
	dir := filepath.Join(dataRoot, scope)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return &HR{path: filepath.Join(dir, "applicants.json")}, nil
}

func (h *HR) load() ([]Applicant, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	data, err := os.ReadFile(h.path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Applicant{}, nil
		}
		return nil, err
	}
	var out []Applicant
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("parse %s: %w", h.path, err)
	}
	return out, nil
}

func (h *HR) save(in []Applicant) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	data, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(h.path, data, 0o644)
}

// List returns all applicants in this scope.
func (h *HR) List() ([]Applicant, error) { return h.load() }

// Store adds or replaces (by ID) an applicant.
func (h *HR) Store(a Applicant) (Applicant, error) {
	if a.ID == "" {
		a.ID = fmt.Sprintf("app_%d", time.Now().UnixNano())
	}
	if a.Status == "" {
		a.Status = "new"
	}
	if a.ReceivedDate == "" {
		a.ReceivedDate = time.Now().UTC().Format("2006-01-02")
	}
	cur, err := h.load()
	if err != nil {
		return a, err
	}
	replaced := false
	for i := range cur {
		if cur[i].ID == a.ID {
			cur[i] = a
			replaced = true
			break
		}
	}
	if !replaced {
		cur = append(cur, a)
	}
	if err := h.save(cur); err != nil {
		return a, err
	}
	return a, nil
}

// UpdateStatus updates one applicant's status and returns the updated record.
func (h *HR) UpdateStatus(id, status string) (*Applicant, error) {
	cur, err := h.load()
	if err != nil {
		return nil, err
	}
	for i := range cur {
		if cur[i].ID == id {
			cur[i].Status = status
			if err := h.save(cur); err != nil {
				return nil, err
			}
			return &cur[i], nil
		}
	}
	return nil, fmt.Errorf("applicant %q not found", id)
}

// Match scores every applicant against the JD and returns the top N matches
// (sorted descending by score). Score is a simple keyword overlap on
// applied_for, key_skills and certifications, plus an experience bonus.
// Honest about what it is: a deterministic matcher for demoability, not ML.
func (h *HR) Match(jd JD, topN int) ([]ApplicantMatch, error) {
	apps, err := h.load()
	if err != nil {
		return nil, err
	}
	if topN <= 0 {
		topN = 5
	}
	out := make([]ApplicantMatch, 0, len(apps))
	for _, a := range apps {
		m := scoreApplicant(a, jd)
		if m.Score > 0 {
			out = append(out, m)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Score > out[j].Score })
	if len(out) > topN {
		out = out[:topN]
	}
	return out, nil
}

func scoreApplicant(a Applicant, jd JD) ApplicantMatch {
	haystack := strings.ToLower(
		a.AppliedFor + " | " +
			strings.Join(a.KeySkills, " ") + " | " +
			strings.Join(a.Certifications, " "),
	)

	matched := []string{}
	missing := []string{}
	score := 0
	// Title overlap — coarse but visible signal.
	if jd.Title != "" && strings.Contains(haystack, strings.ToLower(jd.Title)) {
		score += 5
		matched = append(matched, "title:"+jd.Title)
	}
	for _, kw := range jd.MustHave {
		if kw == "" {
			continue
		}
		if strings.Contains(haystack, strings.ToLower(kw)) {
			score += 3
			matched = append(matched, "must:"+kw)
		} else {
			missing = append(missing, "must:"+kw)
			score -= 2 // missing a must-have is a real cost
		}
	}
	for _, kw := range jd.NiceToHave {
		if kw == "" {
			continue
		}
		if strings.Contains(haystack, strings.ToLower(kw)) {
			score += 1
			matched = append(matched, "nice:"+kw)
		}
	}
	// Experience bonus / penalty.
	if jd.MinYears > 0 {
		if a.ExperienceYears >= jd.MinYears {
			score += 2
			matched = append(matched, fmt.Sprintf("years:%d≥%d", a.ExperienceYears, jd.MinYears))
		} else if a.ExperienceYears > 0 {
			missing = append(missing, fmt.Sprintf("years:%d<%d", a.ExperienceYears, jd.MinYears))
		}
	}
	if score < 0 {
		score = 0
	}
	return ApplicantMatch{Applicant: a, Score: score, Matched: matched, Missing: missing}
}
