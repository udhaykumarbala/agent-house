package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"pty-claude-test/internal/agent"
	"pty-claude-test/internal/discovery"
	"pty-claude-test/internal/session"
)

// ── Client discovery (software studio) ───────────────────────────────
//
// POST /api/discovery/run        start the PRD/research/design pipeline
// GET  /api/discovery/status     poll its manifest
// GET  /api/docs[?project=]      list rich documents for the Documents page

// apiLLM adapts the oneshot API client to the discovery.LLM interface.
type apiLLM struct{ c *session.APIClient }

func (a apiLLM) Generate(ctx context.Context, system, user string) (string, error) {
	resp, err := a.c.SendMessage(ctx, system, user)
	if err != nil {
		return "", err
	}
	return resp.Text, nil
}

// discoveryRunner builds a runner wired to the server's registry prompts
// and the message store/hub so steps stream to the workforce panel.
func (s *Server) discoveryRunner() *discovery.Runner {
	return &discovery.Runner{
		ProjectsDir: s.projectDir,
		LLM:         apiLLM{c: session.NewAPIClient()},
		Prompts: func(role string) string {
			if s.agentRegistry == nil {
				return ""
			}
			if a, err := s.agentRegistry.CreateAgent(agent.Role(role)); err == nil {
				return a.SystemPrompt
			}
			return ""
		},
		Emitter: &scenarioEmitter{store: s.store, hub: s.hub},
	}
}

var projectIDRe = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)

func safeProjectID(p string) bool {
	return projectIDRe.MatchString(p) && !strings.Contains(p, "..")
}

func (s *Server) handleDiscoveryRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}
	var req struct {
		EmailID   string `json:"email_id"`
		Request   string `json:"request"`
		Client    string `json:"client"`
		Subject   string `json:"subject"`
		ProjectID string `json:"project_id"`
	}
	if r.ContentLength > 0 {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad json: "+err.Error(), 400)
			return
		}
	}
	if req.ProjectID != "" && !safeProjectID(req.ProjectID) {
		http.Error(w, "invalid project_id", 400)
		return
	}
	runner := s.discoveryRunner()
	resolved, err := runner.Resolve(discovery.RunInput{
		EmailID: req.EmailID, Request: req.Request, Client: req.Client,
		Subject: req.Subject, ProjectID: req.ProjectID,
	})
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}

	// One pipeline per project at a time.
	s.mu.Lock()
	if s.discoveryRuns == nil {
		s.discoveryRuns = map[string]bool{}
	}
	if s.discoveryRuns[resolved.ProjectID] {
		s.mu.Unlock()
		writeJSON(w, map[string]any{"success": false, "error": "discovery already running for " + resolved.ProjectID})
		return
	}
	s.discoveryRuns[resolved.ProjectID] = true
	s.mu.Unlock()

	go func() {
		defer func() {
			s.mu.Lock()
			delete(s.discoveryRuns, resolved.ProjectID)
			s.mu.Unlock()
		}()
		if _, err := runner.Run(context.Background(), resolved); err != nil {
			fmt.Printf("[DISCOVERY] %s failed: %v\n", resolved.ProjectID, err)
		}
	}()

	writeJSON(w, map[string]any{
		"success": true, "started": true,
		"project_id": resolved.ProjectID,
		"client":     resolved.Client,
	})
}

func (s *Server) handleDiscoveryStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", 405)
		return
	}
	project := r.URL.Query().Get("project")
	if !safeProjectID(project) {
		http.Error(w, "invalid project", 400)
		return
	}
	data, err := os.ReadFile(filepath.Join(s.projectDir, project, "discovery", "manifest.json"))
	if err != nil {
		http.Error(w, "no discovery run for "+project, 404)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}

// docInfo is one renderable document.
type docInfo struct {
	Path     string `json:"path"` // relative to the project dir
	Name     string `json:"name"`
	Title    string `json:"title"`
	Size     int64  `json:"size"`
	Modified int64  `json:"modified"`
}

// handleDocs serves the Documents page: without ?project= it lists the
// projects that have documents; with it, that project's docs + pipeline
// status.
func (s *Server) handleDocs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", 405)
		return
	}
	project := r.URL.Query().Get("project")
	if project == "" {
		s.listDocProjects(w)
		return
	}
	if !safeProjectID(project) {
		http.Error(w, "invalid project", 400)
		return
	}
	projDir := filepath.Join(s.projectDir, project)
	docs := collectDocs(projDir)
	status := ""
	if raw, err := os.ReadFile(filepath.Join(projDir, "discovery", "manifest.json")); err == nil {
		var m struct {
			Status string `json:"status"`
		}
		if json.Unmarshal(raw, &m) == nil {
			status = m.Status
		}
	}
	writeJSON(w, map[string]any{"project": project, "docs": docs, "status": status})
}

func (s *Server) listDocProjects(w http.ResponseWriter) {
	type projRow struct {
		ID       string `json:"id"`
		Docs     int    `json:"docs"`
		Status   string `json:"status,omitempty"`
		Modified int64  `json:"modified"`
	}
	rows := []projRow{}
	entries, _ := os.ReadDir(s.projectDir)
	for _, e := range entries {
		if !e.IsDir() || strings.HasPrefix(e.Name(), ".") || e.Name() == "inbox" || e.Name() == "outbox" {
			continue
		}
		docs := collectDocs(filepath.Join(s.projectDir, e.Name()))
		if len(docs) == 0 {
			continue
		}
		row := projRow{ID: e.Name(), Docs: len(docs)}
		for _, d := range docs {
			if d.Modified > row.Modified {
				row.Modified = d.Modified
			}
		}
		if raw, err := os.ReadFile(filepath.Join(s.projectDir, e.Name(), "discovery", "manifest.json")); err == nil {
			var m struct {
				Status string `json:"status"`
			}
			if json.Unmarshal(raw, &m) == nil {
				row.Status = m.Status
			}
		}
		rows = append(rows, row)
	}
	sort.Slice(rows, func(a, b int) bool { return rows[a].Modified > rows[b].Modified })
	writeJSON(w, map[string]any{"projects": rows})
}

// collectDocs gathers markdown artifacts from the places agents write:
// discovery/ (this pipeline) and .plans/ (the SDLC orchestrator).
func collectDocs(projDir string) []docInfo {
	out := []docInfo{}
	// discovery/*.md — flat, ordered by filename prefix.
	if entries, err := os.ReadDir(filepath.Join(projDir, "discovery")); err == nil {
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
				continue
			}
			if d, ok := docFrom(projDir, filepath.Join("discovery", e.Name())); ok {
				out = append(out, d)
			}
		}
	}
	// .plans/**/*.md — the orchestrator's research/specs/final docs.
	plansRoot := filepath.Join(projDir, ".plans")
	_ = filepath.Walk(plansRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}
		rel, rerr := filepath.Rel(projDir, path)
		if rerr != nil {
			return nil
		}
		if d, ok := docFrom(projDir, rel); ok {
			out = append(out, d)
		}
		return nil
	})
	return out
}

func docFrom(projDir, rel string) (docInfo, bool) {
	full := filepath.Join(projDir, rel)
	info, err := os.Stat(full)
	if err != nil {
		return docInfo{}, false
	}
	d := docInfo{
		Path: filepath.ToSlash(rel), Name: filepath.Base(rel),
		Size: info.Size(), Modified: info.ModTime().Unix(),
	}
	d.Title = strings.TrimSuffix(d.Name, ".md")
	// First "# " heading becomes the display title.
	if raw, err := os.ReadFile(full); err == nil {
		for _, line := range strings.SplitN(string(raw), "\n", 20) {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "# ") {
				d.Title = strings.TrimSpace(strings.TrimPrefix(line, "# "))
				break
			}
		}
	}
	return d, true
}
