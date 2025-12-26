package orchestrator

import (
	"testing"

	"pty-claude-test/internal/agent"
	"pty-claude-test/internal/message"
)

func TestNewOrchestrator(t *testing.T) {
	store := message.NewStore()
	config := DefaultConfig()

	orch := New(config, store)

	if orch == nil {
		t.Fatal("Expected orchestrator to be created")
	}

	if orch.store != store {
		t.Error("Store not set correctly")
	}

	if orch.config.MaxDepth != 3 {
		t.Errorf("Expected MaxDepth 3, got %d", orch.config.MaxDepth)
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.MaxDepth != 3 {
		t.Errorf("Expected MaxDepth 3, got %d", config.MaxDepth)
	}

	if config.MaxTurns != 8 {
		t.Errorf("Expected MaxTurns 8, got %d", config.MaxTurns)
	}

	if !config.EnableFiles {
		t.Error("Expected EnableFiles to be true")
	}

	if !config.Verbose {
		t.Error("Expected Verbose to be true")
	}
}

func TestGetAgent(t *testing.T) {
	store := message.NewStore()
	config := DefaultConfig()
	config.Verbose = false

	orch := New(config, store)

	// Test agent caching by pre-populating
	testAgent := &agent.Agent{
		Name: "Test CEO",
		Role: agent.RoleCEO,
	}
	orch.agents[agent.RoleCEO] = testAgent

	// Get cached agent
	ceo, err := orch.getAgent(agent.RoleCEO)
	if err != nil {
		t.Fatalf("Failed to get cached CEO agent: %v", err)
	}

	if ceo != testAgent {
		t.Error("Expected same agent instance from cache")
	}

	if ceo.Role != agent.RoleCEO {
		t.Errorf("Expected role CEO, got %s", ceo.Role)
	}
}

func TestGetAgents(t *testing.T) {
	store := message.NewStore()
	config := DefaultConfig()
	config.Verbose = false

	orch := New(config, store)

	// Pre-populate agents
	orch.agents[agent.RoleCEO] = &agent.Agent{Name: "CEO", Role: agent.RoleCEO}
	orch.agents[agent.RolePM] = &agent.Agent{Name: "PM", Role: agent.RolePM}

	agents := orch.GetAgents()

	if len(agents) != 2 {
		t.Errorf("Expected 2 agents, got %d", len(agents))
	}

	if _, ok := agents[agent.RoleCEO]; !ok {
		t.Error("CEO agent missing from GetAgents")
	}

	if _, ok := agents[agent.RolePM]; !ok {
		t.Error("PM agent missing from GetAgents")
	}
}

func TestOnMessage(t *testing.T) {
	store := message.NewStore()
	config := DefaultConfig()
	config.Verbose = false

	orch := New(config, store)

	received := make([]*message.Message, 0)
	orch.OnMessage(func(msg *message.Message) {
		received = append(received, msg)
	})

	// Create and notify a message
	msg := message.NewMessage(message.TypeTask, "user", "ceo", "test task")
	orch.notify(msg)

	if len(received) != 1 {
		t.Errorf("Expected 1 message received, got %d", len(received))
	}

	if received[0].Content != "test task" {
		t.Errorf("Expected 'test task', got '%s'", received[0].Content)
	}
}

func TestGetAgentEmoji(t *testing.T) {
	tests := []struct {
		role     agent.Role
		expected string
	}{
		{agent.RoleCEO, "👔"},
		{agent.RolePM, "📋"},
		{agent.RoleSeniorDev, "👨‍💻"},
		{agent.RoleJuniorDev, "👩‍💻"},
		{agent.RoleSecurity, "🔒"},
		{agent.RoleArchitect, "🏗️"},
		{agent.RoleUX, "🎨"},
		{agent.RoleUI, "🖼️"},
	}

	for _, tt := range tests {
		emoji := getAgentEmoji(tt.role)
		if emoji != tt.expected {
			t.Errorf("getAgentEmoji(%s) = %s, expected %s", tt.role, emoji, tt.expected)
		}
	}
}

func TestBuildDelegationContext(t *testing.T) {
	a := &agent.Agent{
		Name: "CEO",
		Role: agent.RoleCEO,
	}

	response := &agent.Response{
		Content: "Please handle the requirements",
	}

	context := buildDelegationContext(a, response, "Build a todo app")

	if context == "" {
		t.Error("Expected non-empty context")
	}

	// Check context contains required elements
	if !containsAll(context, []string{"CEO", "Build a todo app", "Please handle the requirements"}) {
		t.Error("Context missing required elements")
	}
}

func containsAll(s string, subs []string) bool {
	for _, sub := range subs {
		found := false
		for i := 0; i <= len(s)-len(sub); i++ {
			if s[i:i+len(sub)] == sub {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func TestResult(t *testing.T) {
	result := &Result{
		ProjectID:  "test-project",
		TotalTurns: 5,
		Messages:   make([]*message.Message, 3),
		Files:      make([]FileResult, 2),
	}

	if result.ProjectID != "test-project" {
		t.Errorf("Expected project ID 'test-project', got '%s'", result.ProjectID)
	}

	if result.TotalTurns != 5 {
		t.Errorf("Expected 5 turns, got %d", result.TotalTurns)
	}

	if len(result.Messages) != 3 {
		t.Errorf("Expected 3 messages, got %d", len(result.Messages))
	}

	if len(result.Files) != 2 {
		t.Errorf("Expected 2 files, got %d", len(result.Files))
	}
}
