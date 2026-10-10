package web

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func seedDocsProject(t *testing.T, root, project string) {
	t.Helper()
	dir := filepath.Join(root, project, "discovery")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(filepath.Join(dir, "01-prd.md"), []byte("# FitFlow PRD\n\nbody"), 0o644)
	_ = os.WriteFile(filepath.Join(dir, "04-architecture.md"), []byte("# System Architecture\n\n```mermaid\nflowchart TD\nA-->B\n```"), 0o644)
	_ = os.WriteFile(filepath.Join(dir, "manifest.json"), []byte(`{"project_id":"`+project+`","status":"done","steps":[]}`), 0o644)
	// .plans docs from the build pipeline should be listed too.
	plans := filepath.Join(root, project, ".plans", "research")
	_ = os.MkdirAll(plans, 0o755)
	_ = os.WriteFile(filepath.Join(plans, "pm-research.md"), []byte("# PM Research"), 0o644)
}

func TestDocsListForProject(t *testing.T) {
	root := t.TempDir()
	seedDocsProject(t, root, "fitflow")
	s := &Server{projectDir: root}

	rec := httptest.NewRecorder()
	s.handleDocs(rec, httptest.NewRequest("GET", "/api/docs?project=fitflow", nil))
	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var body struct {
		Docs []struct {
			Path  string `json:"path"`
			Title string `json:"title"`
		} `json:"docs"`
		Status string `json:"status"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &body)
	if len(body.Docs) != 3 {
		t.Fatalf("expected 3 docs (2 discovery + 1 plans), got %d: %s", len(body.Docs), rec.Body.String())
	}
	var sawPRD, sawPlans bool
	for _, d := range body.Docs {
		if d.Path == "discovery/01-prd.md" && d.Title == "FitFlow PRD" {
			sawPRD = true
		}
		if strings.HasPrefix(d.Path, ".plans/") {
			sawPlans = true
		}
	}
	if !sawPRD || !sawPlans {
		t.Fatalf("doc listing incomplete: %s", rec.Body.String())
	}
	if body.Status != "done" {
		t.Fatalf("expected manifest status surfaced, got %q", body.Status)
	}
}

func TestDocsListProjects(t *testing.T) {
	root := t.TempDir()
	seedDocsProject(t, root, "fitflow")
	_ = os.MkdirAll(filepath.Join(root, "no-docs-project"), 0o755)
	s := &Server{projectDir: root}

	rec := httptest.NewRecorder()
	s.handleDocs(rec, httptest.NewRequest("GET", "/api/docs", nil))
	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var body struct {
		Projects []struct {
			ID   string `json:"id"`
			Docs int    `json:"docs"`
		} `json:"projects"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &body)
	if len(body.Projects) != 1 || body.Projects[0].ID != "fitflow" || body.Projects[0].Docs != 3 {
		t.Fatalf("unexpected projects listing: %s", rec.Body.String())
	}
}

func TestDocsRejectsTraversal(t *testing.T) {
	s := &Server{projectDir: t.TempDir()}
	for _, p := range []string{"../etc", "a/b", "..", "x%2f..%2fy"} {
		rec := httptest.NewRecorder()
		s.handleDocs(rec, httptest.NewRequest("GET", "/api/docs?project="+p, nil))
		if rec.Code != 400 {
			t.Fatalf("project %q: expected 400, got %d", p, rec.Code)
		}
	}
}

func TestDiscoveryRunRejectsWithoutRequirements(t *testing.T) {
	s := &Server{projectDir: t.TempDir()} // empty inbox, no request
	rec := httptest.NewRecorder()
	s.handleDiscoveryRun(rec, httptest.NewRequest("POST", "/api/discovery/run", strings.NewReader(`{}`)))
	if rec.Code != 404 {
		t.Fatalf("expected 404 when no requirements anywhere, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestDiscoveryRunIsPOSTOnly(t *testing.T) {
	s := &Server{projectDir: t.TempDir()}
	rec := httptest.NewRecorder()
	s.handleDiscoveryRun(rec, httptest.NewRequest("GET", "/api/discovery/run", nil))
	if rec.Code != 405 {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestDiscoveryStatus(t *testing.T) {
	root := t.TempDir()
	seedDocsProject(t, root, "fitflow")
	s := &Server{projectDir: root}

	rec := httptest.NewRecorder()
	s.handleDiscoveryStatus(rec, httptest.NewRequest("GET", "/api/discovery/status?project=fitflow", nil))
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), `"done"`) {
		t.Fatalf("expected manifest, got %d: %s", rec.Code, rec.Body.String())
	}
	rec = httptest.NewRecorder()
	s.handleDiscoveryStatus(rec, httptest.NewRequest("GET", "/api/discovery/status?project=missing", nil))
	if rec.Code != 404 {
		t.Fatalf("expected 404 for missing manifest, got %d", rec.Code)
	}
}
