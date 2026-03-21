package session

import (
	"testing"
)

func TestSessionManagerGetOrCreate(t *testing.T) {
	sm := NewSessionManager()
	defer sm.Shutdown()

	config := SessionConfig{
		AgentRole: "ceo",
		AgentName: "CEO Agent",
		ProjectID: "test-proj",
		WorkDir:   "/tmp",
	}

	// First call creates a session
	sess1, err := sm.GetOrCreate(config)
	if err != nil {
		t.Fatalf("GetOrCreate failed: %v", err)
	}
	if sess1 == nil {
		t.Fatal("expected non-nil session")
	}
	if sess1.Config.AgentRole != "ceo" {
		t.Errorf("expected role ceo, got %s", sess1.Config.AgentRole)
	}

	// Second call returns the same session
	sess2, err := sm.GetOrCreate(config)
	if err != nil {
		t.Fatalf("second GetOrCreate failed: %v", err)
	}
	if sess1.ID != sess2.ID {
		t.Error("expected same session on second call")
	}
}

func TestSessionManagerListByProject(t *testing.T) {
	sm := NewSessionManager()
	defer sm.Shutdown()

	// Create sessions for two projects
	sm.GetOrCreate(SessionConfig{AgentRole: "ceo", AgentName: "CEO", ProjectID: "proj-a", WorkDir: "/tmp"})
	sm.GetOrCreate(SessionConfig{AgentRole: "pm", AgentName: "PM", ProjectID: "proj-a", WorkDir: "/tmp"})
	sm.GetOrCreate(SessionConfig{AgentRole: "ceo", AgentName: "CEO", ProjectID: "proj-b", WorkDir: "/tmp"})

	sessionsA := sm.ListByProject("proj-a")
	if len(sessionsA) != 2 {
		t.Errorf("expected 2 sessions for proj-a, got %d", len(sessionsA))
	}

	sessionsB := sm.ListByProject("proj-b")
	if len(sessionsB) != 1 {
		t.Errorf("expected 1 session for proj-b, got %d", len(sessionsB))
	}
}

func TestSessionManagerGet(t *testing.T) {
	sm := NewSessionManager()
	defer sm.Shutdown()

	sm.GetOrCreate(SessionConfig{AgentRole: "architect", AgentName: "Architect", ProjectID: "proj-1", WorkDir: "/tmp"})

	sess := sm.Get("proj-1", "architect")
	if sess == nil {
		t.Fatal("expected to find session")
	}
	if sess.Config.AgentRole != "architect" {
		t.Errorf("expected architect, got %s", sess.Config.AgentRole)
	}

	noSess := sm.Get("proj-1", "nonexistent")
	if noSess != nil {
		t.Error("expected nil for nonexistent role")
	}
}

func TestSessionManagerTerminateProject(t *testing.T) {
	sm := NewSessionManager()
	defer sm.Shutdown()

	sm.GetOrCreate(SessionConfig{AgentRole: "ceo", AgentName: "CEO", ProjectID: "proj-x", WorkDir: "/tmp"})
	sm.GetOrCreate(SessionConfig{AgentRole: "pm", AgentName: "PM", ProjectID: "proj-x", WorkDir: "/tmp"})
	sm.GetOrCreate(SessionConfig{AgentRole: "ceo", AgentName: "CEO", ProjectID: "proj-y", WorkDir: "/tmp"})

	sm.TerminateProject("proj-x")

	if len(sm.ListByProject("proj-x")) != 0 {
		t.Error("expected 0 sessions after termination")
	}
	if len(sm.ListByProject("proj-y")) != 1 {
		t.Error("proj-y session should be unaffected")
	}
}

func TestSessionManagerListAll(t *testing.T) {
	sm := NewSessionManager()
	defer sm.Shutdown()

	sm.GetOrCreate(SessionConfig{AgentRole: "ceo", AgentName: "CEO", ProjectID: "proj-1", WorkDir: "/tmp"})
	sm.GetOrCreate(SessionConfig{AgentRole: "pm", AgentName: "PM", ProjectID: "proj-1", WorkDir: "/tmp"})

	all := sm.ListAll()
	if len(all) != 2 {
		t.Errorf("expected 2 sessions, got %d", len(all))
	}
}

func TestSessionManagerProjectMetrics(t *testing.T) {
	sm := NewSessionManager()
	defer sm.Shutdown()

	sm.GetOrCreate(SessionConfig{AgentRole: "ceo", AgentName: "CEO", ProjectID: "proj-1", WorkDir: "/tmp"})
	sm.GetOrCreate(SessionConfig{AgentRole: "pm", AgentName: "PM", ProjectID: "proj-1", WorkDir: "/tmp"})

	metrics := sm.GetProjectMetrics("proj-1")
	if metrics["project_id"] != "proj-1" {
		t.Errorf("expected project_id proj-1, got %v", metrics["project_id"])
	}
	agentMetrics, ok := metrics["agent_metrics"].(map[string]interface{})
	if !ok {
		t.Fatal("expected agent_metrics map")
	}
	if len(agentMetrics) != 2 {
		t.Errorf("expected 2 agent metrics, got %d", len(agentMetrics))
	}
}
