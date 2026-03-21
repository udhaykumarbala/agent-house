package session

import (
	"encoding/json"
	"log"
	"strings"
	"time"
)

// Translator converts Claude NDJSON events into AgentEvents.
type Translator struct {
	agentRole string
	agentName string
	projectID string
	taskID    string
}

// NewTranslator creates a Translator bound to a specific agent context.
func NewTranslator(role, name, projectID string) *Translator {
	return &Translator{
		agentRole: role,
		agentName: name,
		projectID: projectID,
	}
}

// SetTaskID updates the current task context for event tagging.
func (t *Translator) SetTaskID(taskID string) {
	t.taskID = taskID
}

// Reset clears state for a new turn.
func (t *Translator) Reset() {
	// No-op — content index tracking removed since Claude Code
	// stream-json emits only new blocks per line (not cumulative).
}

// Translate converts a raw ClaudeEvent into zero or more AgentEvents.
func (t *Translator) Translate(raw ClaudeEvent) []AgentEvent {
	switch raw.Type {
	case "system":
		return t.translateSystem(raw)
	case "result":
		return t.translateResult(raw)
	case "assistant":
		return t.translateAssistant(raw)
	case "user":
		return t.translateUser(raw)
	default:
		return nil
	}
}

func (t *Translator) newEvent(eventType string) AgentEvent {
	return AgentEvent{
		Type:      eventType,
		AgentRole: t.agentRole,
		AgentName: t.agentName,
		ProjectID: t.projectID,
		TaskID:    t.taskID,
		Timestamp: time.Now().UnixMilli(),
	}
}

func (t *Translator) translateSystem(raw ClaudeEvent) []AgentEvent {
	switch raw.Subtype {
	case "init":
		ev := t.newEvent("session_meta")
		ev.Model = raw.Model
		return []AgentEvent{ev}
	case "result":
		ev := t.newEvent("turn_complete")
		if raw.Result != nil {
			var resultStr string
			if err := json.Unmarshal(raw.Result, &resultStr); err == nil {
				ev.Content = resultStr
			}
		}
		if raw.IsError {
			ev.Type = "error"
		}
		t.Reset()
		return []AgentEvent{ev}
	default:
		return nil
	}
}

func (t *Translator) translateResult(raw ClaudeEvent) []AgentEvent {
	ev := t.newEvent("turn_complete")
	ev.CostUSD = raw.TotalCostUSD

	if raw.Result != nil {
		var resultStr string
		if err := json.Unmarshal(raw.Result, &resultStr); err == nil {
			ev.Content = resultStr
		}
	}
	if raw.Subtype == "error" || raw.IsError {
		ev.Type = "error"
	}

	// Extract usage from modelUsage
	if raw.ModelUsage != nil {
		var modelUsage map[string]struct {
			InputTokens              int `json:"inputTokens"`
			OutputTokens             int `json:"outputTokens"`
			CacheReadInputTokens     int `json:"cacheReadInputTokens"`
			CacheCreationInputTokens int `json:"cacheCreationInputTokens"`
			ContextWindow            int `json:"contextWindow"`
		}
		if err := json.Unmarshal(raw.ModelUsage, &modelUsage); err == nil {
			for _, u := range modelUsage {
				ev.InputTokens = u.InputTokens + u.CacheReadInputTokens + u.CacheCreationInputTokens
				ev.OutputTokens = u.OutputTokens
				break
			}
		}
	}

	t.Reset()
	return []AgentEvent{ev}
}

func (t *Translator) translateAssistant(raw ClaudeEvent) []AgentEvent {
	if raw.Message == nil {
		return nil
	}

	var msg ClaudeMessage
	if err := json.Unmarshal(raw.Message, &msg); err != nil {
		log.Printf("translator[%s]: failed to parse assistant message: %v", t.agentRole, err)
		return nil
	}

	var events []AgentEvent

	// Claude Code stream-json (without --include-partial-messages) emits each
	// assistant NDJSON line with ONLY the new content blocks for that step, not
	// a cumulative array. So we always process all blocks from index 0.
	for i := 0; i < len(msg.Content); i++ {
		block := msg.Content[i]
		switch block.Type {
		case "text":
			ev := t.newEvent("text_delta")
			ev.Content = block.Text
			events = append(events, ev)

		case "thinking":
			text := block.Thinking
			if text == "" {
				text = block.Text
			}
			if text != "" {
				ev := t.newEvent("thinking_delta")
				ev.Content = text
				events = append(events, ev)
			}

		case "tool_use":
			ev := t.newEvent("tool_use")
			ev.ToolUseID = block.ID
			ev.ToolName = block.Name
			if block.Input != nil {
				ev.Input = string(block.Input)
			}
			events = append(events, ev)

		case "tool_result":
			ev := t.newEvent("tool_result")
			ev.ToolUseID = block.ToolUseID
			ev.Output = extractContent(block.Content)
			ev.IsError = block.IsError
			events = append(events, ev)
		}
	}

	return events
}

func (t *Translator) translateUser(raw ClaudeEvent) []AgentEvent {
	if raw.Message == nil {
		return nil
	}

	var msg ClaudeMessage
	if err := json.Unmarshal(raw.Message, &msg); err != nil {
		log.Printf("translator[%s]: failed to parse user message: %v", t.agentRole, err)
		return nil
	}

	var events []AgentEvent
	for _, block := range msg.Content {
		if block.Type == "tool_result" {
			ev := t.newEvent("tool_result")
			ev.ToolUseID = block.ToolUseID
			ev.Output = extractContent(block.Content)
			ev.IsError = block.IsError
			events = append(events, ev)
		}
	}
	return events
}

// extractContent handles the polymorphic content field in tool_result blocks.
func extractContent(raw json.RawMessage) string {
	if raw == nil {
		return ""
	}

	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}

	var blocks []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}
	if err := json.Unmarshal(raw, &blocks); err == nil {
		var parts []string
		for _, b := range blocks {
			if b.Text != "" {
				parts = append(parts, b.Text)
			}
		}
		return strings.Join(parts, "\n")
	}

	return string(raw)
}
