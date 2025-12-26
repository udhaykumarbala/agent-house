package message

import (
	"encoding/json"
	"testing"
)

func TestNewMessage(t *testing.T) {
	msg := NewMessage(TypeTask, "user", "ceo", "Build a todo app")

	if msg.Type != TypeTask {
		t.Errorf("expected type %s, got %s", TypeTask, msg.Type)
	}
	if msg.From != "user" {
		t.Errorf("expected from 'user', got '%s'", msg.From)
	}
	if msg.To != "ceo" {
		t.Errorf("expected to 'ceo', got '%s'", msg.To)
	}
	if msg.Content != "Build a todo app" {
		t.Errorf("unexpected content: %s", msg.Content)
	}
	if msg.ID == "" {
		t.Error("expected ID to be generated")
	}
	if msg.Timestamp.IsZero() {
		t.Error("expected timestamp to be set")
	}
}

func TestNewTaskMessage(t *testing.T) {
	msg := NewTaskMessage("project-123", "Create a landing page")

	if msg.Type != TypeTask {
		t.Errorf("expected type %s, got %s", TypeTask, msg.Type)
	}
	if msg.Metadata.ProjectID != "project-123" {
		t.Errorf("expected project ID 'project-123', got '%s'", msg.Metadata.ProjectID)
	}
	if msg.Metadata.TaskID == "" {
		t.Error("expected task ID to be generated")
	}
}

func TestNewDelegateMessage(t *testing.T) {
	msg := NewDelegateMessage("ceo", []string{"pm", "architect"}, "Need requirements and tech stack")

	if msg.Type != TypeDelegate {
		t.Errorf("expected type %s, got %s", TypeDelegate, msg.Type)
	}
	if msg.From != "ceo" {
		t.Errorf("expected from 'ceo', got '%s'", msg.From)
	}
	if len(msg.Metadata.DelegateTo) != 2 {
		t.Errorf("expected 2 delegates, got %d", len(msg.Metadata.DelegateTo))
	}
}

func TestNewFileCreateMessage(t *testing.T) {
	msg := NewFileCreateMessage("pm", "docs/requirements.md", "# Requirements\n\n- Feature 1")

	if msg.Type != TypeFileCreate {
		t.Errorf("expected type %s, got %s", TypeFileCreate, msg.Type)
	}
	if msg.Metadata.FilePath != "docs/requirements.md" {
		t.Errorf("expected path 'docs/requirements.md', got '%s'", msg.Metadata.FilePath)
	}
	if msg.Metadata.FileContent == "" {
		t.Error("expected file content to be set")
	}
}

func TestMessageJSON(t *testing.T) {
	original := NewTaskMessage("proj-1", "Build something")
	original.Metadata.Tags = []string{"urgent", "mvp"}

	// Serialize
	data, err := original.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	// Deserialize
	parsed, err := FromJSON(data)
	if err != nil {
		t.Fatalf("FromJSON failed: %v", err)
	}

	// Compare
	if parsed.ID != original.ID {
		t.Errorf("ID mismatch: %s != %s", parsed.ID, original.ID)
	}
	if parsed.Type != original.Type {
		t.Errorf("Type mismatch: %s != %s", parsed.Type, original.Type)
	}
	if len(parsed.Metadata.Tags) != 2 {
		t.Errorf("Tags not preserved: got %d", len(parsed.Metadata.Tags))
	}
}

func TestMessageString(t *testing.T) {
	msg := NewMessage(TypeResponse, "ceo", "pm", "Please define the requirements for this project")
	str := msg.String()

	if str != "ceo → pm: Please define the requirements for this project" {
		t.Errorf("unexpected string: %s", str)
	}
}

func TestMessageStringTruncation(t *testing.T) {
	longContent := "This is a very long message that should be truncated because it exceeds the maximum length"
	msg := NewMessage(TypeResponse, "a", "b", longContent)
	str := msg.String()

	if len(str) > 60 { // "a → b: " is 7 chars, plus 50 max content
		t.Errorf("string not truncated properly: %s", str)
	}
}

func TestMessageTypes(t *testing.T) {
	types := []MessageType{
		TypeTask,
		TypeResponse,
		TypeDelegate,
		TypeFileCreate,
		TypeFileModify,
		TypeComplete,
		TypeError,
	}

	for _, mt := range types {
		msg := NewMessage(mt, "a", "b", "content")
		data, _ := json.Marshal(msg)
		if string(data) == "" {
			t.Errorf("failed to marshal message type %s", mt)
		}
	}
}
