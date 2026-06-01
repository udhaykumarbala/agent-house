package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"pty-claude-test/internal/agent"
	"pty-claude-test/internal/message"
)

func newTestServer() *Server {
	store := message.NewStore()
	tmpDir, _ := os.MkdirTemp("", "web-test-*")

	config := Config{
		Port:       8080,
		ProjectDir: tmpDir,
		Store:      store,
	}

	return NewServer(config)
}

func TestHandleStatus(t *testing.T) {
	server := newTestServer()

	req := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	w := httptest.NewRecorder()

	server.handleStatus(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok', got %v", response["status"])
	}

	if response["task_running"] != false {
		t.Errorf("Expected task_running false, got %v", response["task_running"])
	}
}

func TestHandleAgents(t *testing.T) {
	server := newTestServer()

	req := httptest.NewRequest(http.MethodGet, "/api/agents", nil)
	w := httptest.NewRecorder()

	server.handleAgents(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	agents, ok := response["agents"].([]interface{})
	if !ok {
		t.Fatal("Expected agents array")
	}

	// Roster is the IT pack (8) + the EPC pack (6) = 14 roles.
	if len(agents) != 14 {
		t.Errorf("Expected 14 agents, got %d", len(agents))
	}

	// With no live sessions, no agent should report active=true. This guards the
	// fix that derives "active" from live session state instead of permanent
	// orchestrator-map membership (which left agents stuck "working" forever).
	for _, a := range agents {
		m, _ := a.(map[string]interface{})
		if active, _ := m["active"].(bool); active {
			t.Errorf("agent %v should not be active with no live sessions", m["role"])
		}
	}
}

func TestHandleMessages(t *testing.T) {
	server := newTestServer()

	// Add some messages
	server.store.Add(message.NewMessage(message.TypeTask, "user", "ceo", "Test task"))
	server.store.Add(message.NewMessage(message.TypeResponse, "ceo", "user", "Response"))

	req := httptest.NewRequest(http.MethodGet, "/api/messages", nil)
	w := httptest.NewRecorder()

	server.handleMessages(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	count := response["count"].(float64)
	if count != 2 {
		t.Errorf("Expected 2 messages, got %v", count)
	}
}

func TestHandleMessagesWithFilter(t *testing.T) {
	server := newTestServer()

	// Add messages for different agents
	msg1 := message.NewMessage(message.TypeResponse, "ceo", "pm", "CEO message")
	msg1.Metadata.ProjectID = "project1"
	server.store.Add(msg1)

	msg2 := message.NewMessage(message.TypeResponse, "pm", "dev", "PM message")
	msg2.Metadata.ProjectID = "project2"
	server.store.Add(msg2)

	// Test agent filter
	req := httptest.NewRequest(http.MethodGet, "/api/messages?agent=ceo", nil)
	w := httptest.NewRecorder()

	server.handleMessages(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	count := response["count"].(float64)
	if count != 1 {
		t.Errorf("Expected 1 message for CEO, got %v", count)
	}

	// Test project filter
	req = httptest.NewRequest(http.MethodGet, "/api/messages?project=project1", nil)
	w = httptest.NewRecorder()

	server.handleMessages(w, req)

	json.Unmarshal(w.Body.Bytes(), &response)

	count = response["count"].(float64)
	if count != 1 {
		t.Errorf("Expected 1 message for project1, got %v", count)
	}
}

func TestHandleTask(t *testing.T) {
	server := newTestServer()

	// Test POST request
	body := bytes.NewBufferString(`{"task": "Build a todo app", "project_id": "test"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/task", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.handleTask(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["success"] != true {
		t.Errorf("Expected success true, got %v", response["success"])
	}

	if response["project_id"] != "test" {
		t.Errorf("Expected project_id 'test', got %v", response["project_id"])
	}
}

func TestHandleTaskMissingBody(t *testing.T) {
	server := newTestServer()

	body := bytes.NewBufferString(`{"project_id": "test"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/task", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.handleTask(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestHandleFiles(t *testing.T) {
	server := newTestServer()

	// Create some test files
	projectPath := filepath.Join(server.projectDir, "default")
	os.MkdirAll(filepath.Join(projectPath, "src"), 0755)
	os.WriteFile(filepath.Join(projectPath, "README.md"), []byte("# Test"), 0644)
	os.WriteFile(filepath.Join(projectPath, "src", "main.go"), []byte("package main"), 0644)

	req := httptest.NewRequest(http.MethodGet, "/api/files?project=default", nil)
	w := httptest.NewRecorder()

	server.handleFiles(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	files, ok := response["files"].([]interface{})
	if !ok {
		t.Fatal("Expected files array")
	}

	if len(files) < 2 {
		t.Errorf("Expected at least 2 files, got %d", len(files))
	}
}

func TestHandleFilesEmptyProject(t *testing.T) {
	server := newTestServer()

	req := httptest.NewRequest(http.MethodGet, "/api/files?project=nonexistent", nil)
	w := httptest.NewRecorder()

	server.handleFiles(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	files, _ := response["files"].([]interface{})
	if len(files) != 0 {
		t.Errorf("Expected 0 files, got %d", len(files))
	}
}

func TestHandleStatic(t *testing.T) {
	server := newTestServer()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	server.handleStatic(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "text/html" {
		t.Errorf("Expected Content-Type text/html, got %s", contentType)
	}

	body := w.Body.String()
	if len(body) < 100 {
		t.Error("Expected HTML content")
	}
}

func TestMethodNotAllowed(t *testing.T) {
	server := newTestServer()

	tests := []struct {
		path   string
		method string
	}{
		{"/api/messages", http.MethodPost},
		{"/api/agents", http.MethodPost},
		{"/api/task", http.MethodGet},
		{"/api/files", http.MethodPost},
		{"/api/status", http.MethodPost},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.path, nil)
		w := httptest.NewRecorder()

		switch tt.path {
		case "/api/messages":
			server.handleMessages(w, req)
		case "/api/agents":
			server.handleAgents(w, req)
		case "/api/task":
			server.handleTask(w, req)
		case "/api/files":
			server.handleFiles(w, req)
		case "/api/status":
			server.handleStatus(w, req)
		}

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("%s %s: Expected 405, got %d", tt.method, tt.path, w.Code)
		}
	}
}

func TestGetRoleName(t *testing.T) {
	tests := []struct {
		role     agent.Role
		expected string
	}{
		{agent.RoleCEO, "Chief Executive Officer"},
		{agent.RolePM, "Product Manager"},
		{agent.RoleSeniorDev, "Senior Developer"},
	}

	for _, tt := range tests {
		name := getRoleName(tt.role)
		if name != tt.expected {
			t.Errorf("getRoleName(%s) = %s, expected %s", tt.role, name, tt.expected)
		}
	}
}

func TestGetRoleColor(t *testing.T) {
	tests := []struct {
		role     agent.Role
		expected string
	}{
		{agent.RoleCEO, "#4A90A4"},
		{agent.RoleSeniorDev, "#2ECC71"},
	}

	for _, tt := range tests {
		color := getRoleColor(tt.role)
		if color != tt.expected {
			t.Errorf("getRoleColor(%s) = %s, expected %s", tt.role, color, tt.expected)
		}
	}
}

func TestCORSMiddleware(t *testing.T) {
	handler := corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Test OPTIONS request
	req := httptest.NewRequest(http.MethodOptions, "/api/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for OPTIONS, got %d", w.Code)
	}

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("Missing CORS header")
	}

	// Test regular request
	req = httptest.NewRequest(http.MethodGet, "/api/test", nil)
	w = httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("Missing CORS header on regular request")
	}
}
