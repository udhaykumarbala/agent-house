package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"pty-claude-test/internal/agent"
	"pty-claude-test/internal/message"
	"pty-claude-test/internal/orchestrator"
)

// CronJob defines a scheduled task.
type CronJob struct {
	ID        string `json:"id"`
	Schedule  string `json:"schedule"`   // cron expression: "*/5 * * * *" or simple: "5m", "1h", "24h"
	Task      string `json:"task"`       // task description for the agent
	AgentRole string `json:"agent_role"` // target agent
	ProjectID string `json:"project_id"`
	Enabled   bool   `json:"enabled"`
	LastRun   string `json:"last_run,omitempty"`
	NextRun   string `json:"next_run,omitempty"`
}

// CronScheduler manages recurring agent tasks.
type CronScheduler struct {
	mu     sync.Mutex
	jobs   map[string]*cronEntry
	orch   *orchestrator.Orchestrator
	store  *message.Store // for the proactive trigger_fired marker
	hub    *Hub           // for live broadcast to /ws subscribers
	stopCh chan struct{}
}

type cronEntry struct {
	job      CronJob
	interval time.Duration
	nextRun  time.Time
	lastRun  time.Time
}

// NewCronScheduler creates a new cron scheduler. store + hub are optional
// (nil-safe in tick) so existing callers without them still work; the
// trigger_fired marker is the only thing that needs them.
func NewCronScheduler(orch *orchestrator.Orchestrator, store *message.Store, hub *Hub) *CronScheduler {
	cs := &CronScheduler{
		jobs:   make(map[string]*cronEntry),
		orch:   orch,
		store:  store,
		hub:    hub,
		stopCh: make(chan struct{}),
	}
	go cs.runLoop()
	return cs
}

func (cs *CronScheduler) runLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			cs.tick()
		case <-cs.stopCh:
			return
		}
	}
}

func (cs *CronScheduler) tick() {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	now := time.Now()
	for _, entry := range cs.jobs {
		if !entry.job.Enabled {
			continue
		}
		if now.After(entry.nextRun) {
			// Fire the job
			log.Printf("[CRON] Firing job %s: %s → %s", entry.job.ID, entry.job.AgentRole, truncateForLog(entry.job.Task, 60))

			injected := &orchestrator.InjectedTask{
				ID:        fmt.Sprintf("cron_%s_%d", entry.job.ID, now.UnixMilli()),
				Task:      entry.job.Task,
				AgentRole: agent.Role(entry.job.AgentRole),
				ProjectID: entry.job.ProjectID,
				Priority:  "normal",
				CreatedAt: now,
			}
			cs.orch.InjectTask(injected)

			// Proactive marker — record that the system fired this on its own.
			// The Conductor briefing scans messages for trigger_fired to surface
			// "N fires today" — without this, the cron does the work silently.
			if cs.store != nil {
				m := message.NewMessage(
					message.TypeTriggerFired,
					"cron",
					string(injected.AgentRole),
					fmt.Sprintf("⚡ cron fired → %s", truncateForLog(entry.job.Task, 120)),
				)
				m.Metadata.ProjectID = entry.job.ProjectID
				m.Metadata.TaskID = injected.ID
				m.Metadata.Tags = []string{"trigger", "cron", entry.job.ID}
				cs.store.Add(m)
				if cs.hub != nil {
					cs.hub.Broadcast(m)
				}
			}

			entry.lastRun = now
			entry.nextRun = now.Add(entry.interval)
		}
	}
}

// AddJob adds a new cron job.
func (cs *CronScheduler) AddJob(job CronJob) error {
	interval, err := parseDuration(job.Schedule)
	if err != nil {
		return fmt.Errorf("invalid schedule %q: %w", job.Schedule, err)
	}

	cs.mu.Lock()
	defer cs.mu.Unlock()

	cs.jobs[job.ID] = &cronEntry{
		job:      job,
		interval: interval,
		nextRun:  time.Now().Add(interval),
	}

	log.Printf("[CRON] Added job %s: every %v → %s", job.ID, interval, job.AgentRole)
	return nil
}

// RemoveJob removes a cron job.
func (cs *CronScheduler) RemoveJob(id string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	delete(cs.jobs, id)
}

// ListJobs returns all jobs.
func (cs *CronScheduler) ListJobs() []CronJob {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	var jobs []CronJob
	for _, entry := range cs.jobs {
		job := entry.job
		if !entry.lastRun.IsZero() {
			job.LastRun = entry.lastRun.Format(time.RFC3339)
		}
		job.NextRun = entry.nextRun.Format(time.RFC3339)
		jobs = append(jobs, job)
	}
	return jobs
}

// Stop shuts down the scheduler.
func (cs *CronScheduler) Stop() {
	close(cs.stopCh)
}

// parseDuration parses simple duration strings: "5m", "1h", "30s", "24h"
func parseDuration(s string) (time.Duration, error) {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, err
	}
	if d < 30*time.Second {
		return 0, fmt.Errorf("minimum interval is 30s")
	}
	return d, nil
}

func truncateForLog(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

// HTTP Handlers

// handleCronList handles GET /api/cron — list all cron jobs
func (s *Server) handleCronList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if s.cronScheduler == nil {
		writeJSON(w, map[string]interface{}{"jobs": []interface{}{}})
		return
	}
	writeJSON(w, map[string]interface{}{"jobs": s.cronScheduler.ListJobs()})
}

// handleCronCreate handles POST /api/cron — create a new cron job
func (s *Server) handleCronCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var job CronJob
	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	if job.Task == "" || job.Schedule == "" {
		http.Error(w, "task and schedule are required", http.StatusBadRequest)
		return
	}
	if job.ID == "" {
		job.ID = fmt.Sprintf("cron_%d", time.Now().UnixMilli())
	}
	if job.AgentRole == "" {
		job.AgentRole = "senior_dev"
	}
	if job.ProjectID == "" {
		job.ProjectID = "default"
	}
	job.Enabled = true

	if s.cronScheduler == nil {
		http.Error(w, "Cron scheduler not initialized", http.StatusServiceUnavailable)
		return
	}

	if err := s.cronScheduler.AddJob(job); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, map[string]interface{}{
		"success": true,
		"job":     job,
	})
}

// handleCronDelete handles DELETE /api/cron/{id}
func (s *Server) handleCronDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}
	if s.cronScheduler != nil {
		s.cronScheduler.RemoveJob(id)
	}
	writeJSON(w, map[string]string{"status": "deleted"})
}

// handleCronRouting routes cron API requests.
func (s *Server) handleCronRouting(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleCronList(w, r)
	case http.MethodPost:
		s.handleCronCreate(w, r)
	case http.MethodDelete:
		s.handleCronDelete(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// LoadCronJobsFromConfig loads cron jobs from .agent-house/cron.json
func LoadCronJobsFromConfig(projectDir string, scheduler *CronScheduler) {
	path := projectDir + "/.agent-house/cron.json"
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	var config struct {
		Jobs []CronJob `json:"jobs"`
	}
	if err := json.Unmarshal(data, &config); err != nil {
		log.Printf("[CRON] Failed to parse %s: %v", path, err)
		return
	}
	for _, job := range config.Jobs {
		if job.Enabled {
			scheduler.AddJob(job)
		}
	}
	log.Printf("[CRON] Loaded %d jobs from %s", len(config.Jobs), path)
}
