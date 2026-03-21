package session

import (
	"encoding/json"
	"testing"
)

func TestTranslateSystemInit(t *testing.T) {
	tr := NewTranslator("ceo", "CEO Agent", "proj-1")

	raw := ClaudeEvent{
		Type:    "system",
		Subtype: "init",
		Model:   "claude-sonnet-4-6",
	}

	events := tr.Translate(raw)
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	ev := events[0]
	if ev.Type != "session_meta" {
		t.Errorf("expected session_meta, got %s", ev.Type)
	}
	if ev.Model != "claude-sonnet-4-6" {
		t.Errorf("expected model claude-sonnet-4-6, got %s", ev.Model)
	}
	if ev.AgentRole != "ceo" {
		t.Errorf("expected agent_role ceo, got %s", ev.AgentRole)
	}
}

func TestTranslateTextDelta(t *testing.T) {
	tr := NewTranslator("pm", "PM Agent", "proj-1")

	msgJSON, _ := json.Marshal(ClaudeMessage{
		Role: "assistant",
		Content: []ClaudeContent{
			{Type: "text", Text: "Hello, I'll analyze this task."},
		},
	})

	raw := ClaudeEvent{
		Type:    "assistant",
		Message: msgJSON,
	}

	events := tr.Translate(raw)
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	ev := events[0]
	if ev.Type != "text_delta" {
		t.Errorf("expected text_delta, got %s", ev.Type)
	}
	if ev.Content != "Hello, I'll analyze this task." {
		t.Errorf("unexpected content: %s", ev.Content)
	}
	if ev.AgentRole != "pm" {
		t.Errorf("expected pm, got %s", ev.AgentRole)
	}
}

func TestTranslateThinking(t *testing.T) {
	tr := NewTranslator("architect", "Architect", "proj-1")

	msgJSON, _ := json.Marshal(ClaudeMessage{
		Role: "assistant",
		Content: []ClaudeContent{
			{Type: "thinking", Thinking: "Let me think about the architecture..."},
		},
	})

	raw := ClaudeEvent{
		Type:    "assistant",
		Message: msgJSON,
	}

	events := tr.Translate(raw)
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Type != "thinking_delta" {
		t.Errorf("expected thinking_delta, got %s", events[0].Type)
	}
	if events[0].Content != "Let me think about the architecture..." {
		t.Errorf("unexpected content: %s", events[0].Content)
	}
}

func TestTranslateToolUse(t *testing.T) {
	tr := NewTranslator("senior_dev", "Senior Developer", "proj-1")

	inputJSON, _ := json.Marshal(map[string]string{
		"file_path": "src/main.go",
		"content":   "package main",
	})

	msgJSON, _ := json.Marshal(ClaudeMessage{
		Role: "assistant",
		Content: []ClaudeContent{
			{
				Type:  "tool_use",
				ID:    "toolu_abc123",
				Name:  "Write",
				Input: inputJSON,
			},
		},
	})

	raw := ClaudeEvent{
		Type:    "assistant",
		Message: msgJSON,
	}

	events := tr.Translate(raw)
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	ev := events[0]
	if ev.Type != "tool_use" {
		t.Errorf("expected tool_use, got %s", ev.Type)
	}
	if ev.ToolName != "Write" {
		t.Errorf("expected tool name Write, got %s", ev.ToolName)
	}
	if ev.ToolUseID != "toolu_abc123" {
		t.Errorf("expected tool_use_id toolu_abc123, got %s", ev.ToolUseID)
	}
}

func TestTranslateToolResult(t *testing.T) {
	tr := NewTranslator("senior_dev", "Senior Developer", "proj-1")

	contentJSON, _ := json.Marshal("File written successfully")

	msgJSON, _ := json.Marshal(ClaudeMessage{
		Role: "user",
		Content: []ClaudeContent{
			{
				Type:      "tool_result",
				ToolUseID: "toolu_abc123",
				Content:   contentJSON,
				IsError:   false,
			},
		},
	})

	raw := ClaudeEvent{
		Type:    "user",
		Message: msgJSON,
	}

	events := tr.Translate(raw)
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	ev := events[0]
	if ev.Type != "tool_result" {
		t.Errorf("expected tool_result, got %s", ev.Type)
	}
	if ev.Output != "File written successfully" {
		t.Errorf("unexpected output: %s", ev.Output)
	}
	if ev.IsError {
		t.Error("expected is_error false")
	}
}

func TestTranslateTurnComplete(t *testing.T) {
	tr := NewTranslator("ceo", "CEO", "proj-1")

	modelUsage, _ := json.Marshal(map[string]interface{}{
		"claude-sonnet-4-6": map[string]int{
			"inputTokens":              1500,
			"outputTokens":             420,
			"cacheReadInputTokens":     100,
			"cacheCreationInputTokens": 0,
			"contextWindow":            200000,
		},
	})

	raw := ClaudeEvent{
		Type:         "result",
		Subtype:      "success",
		TotalCostUSD: 0.0123,
		ModelUsage:   modelUsage,
	}

	events := tr.Translate(raw)
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	ev := events[0]
	if ev.Type != "turn_complete" {
		t.Errorf("expected turn_complete, got %s", ev.Type)
	}
	if ev.CostUSD != 0.0123 {
		t.Errorf("expected cost 0.0123, got %f", ev.CostUSD)
	}
	if ev.InputTokens != 1600 { // 1500 + 100
		t.Errorf("expected input tokens 1600, got %d", ev.InputTokens)
	}
	if ev.OutputTokens != 420 {
		t.Errorf("expected output tokens 420, got %d", ev.OutputTokens)
	}
}

func TestMultipleAssistantMessages(t *testing.T) {
	tr := NewTranslator("pm", "PM", "proj-1")

	// First assistant message with text
	msg1, _ := json.Marshal(ClaudeMessage{
		ID:      "msg_1",
		Content: []ClaudeContent{
			{Type: "text", Text: "Hello"},
		},
	})
	events1 := tr.Translate(ClaudeEvent{Type: "assistant", Message: msg1})
	if len(events1) != 1 {
		t.Fatalf("expected 1 event from first message, got %d", len(events1))
	}

	// Second assistant message (same ID, new content — Claude Code behavior)
	msg2, _ := json.Marshal(ClaudeMessage{
		ID:      "msg_1",
		Content: []ClaudeContent{
			{Type: "tool_use", ID: "toolu_1", Name: "Write", Input: json.RawMessage(`{}`)},
		},
	})
	events2 := tr.Translate(ClaudeEvent{Type: "assistant", Message: msg2})
	if len(events2) != 1 {
		t.Fatalf("expected 1 event from second message, got %d", len(events2))
	}
	if events2[0].Type != "tool_use" {
		t.Errorf("expected tool_use, got %s", events2[0].Type)
	}
}

func TestExtractContent(t *testing.T) {
	// String content
	strJSON, _ := json.Marshal("hello world")
	if got := extractContent(strJSON); got != "hello world" {
		t.Errorf("string: expected 'hello world', got '%s'", got)
	}

	// Array content
	arrJSON, _ := json.Marshal([]map[string]string{
		{"type": "text", "text": "line1"},
		{"type": "text", "text": "line2"},
	})
	if got := extractContent(arrJSON); got != "line1\nline2" {
		t.Errorf("array: expected 'line1\\nline2', got '%s'", got)
	}

	// Nil content
	if got := extractContent(nil); got != "" {
		t.Errorf("nil: expected empty, got '%s'", got)
	}
}

func TestTranslateReset(t *testing.T) {
	tr := NewTranslator("ceo", "CEO", "proj-1")

	// Send a message
	msg, _ := json.Marshal(ClaudeMessage{
		Content: []ClaudeContent{
			{Type: "text", Text: "First turn"},
		},
	})
	tr.Translate(ClaudeEvent{Type: "assistant", Message: msg})

	// Simulate turn complete (resets index)
	tr.Translate(ClaudeEvent{Type: "result", Subtype: "success"})

	// Send another message - should start fresh
	msg2, _ := json.Marshal(ClaudeMessage{
		Content: []ClaudeContent{
			{Type: "text", Text: "Second turn"},
		},
	})
	events := tr.Translate(ClaudeEvent{Type: "assistant", Message: msg2})
	if len(events) != 1 {
		t.Fatalf("expected 1 event after reset, got %d", len(events))
	}
	if events[0].Content != "Second turn" {
		t.Errorf("expected 'Second turn', got '%s'", events[0].Content)
	}
}
