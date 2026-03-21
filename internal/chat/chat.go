package chat

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ChatMessage represents a message between user and agent
type ChatMessage struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	AgentRole string    `json:"agent_role"`
	Direction string    `json:"direction"` // "user_to_agent" or "agent_to_user"
	Content   string    `json:"content"`
	Context   string    `json:"context"` // "idle_chat", "mid_workflow", "checkpoint"
	Timestamp time.Time `json:"timestamp"`
}

func GenerateMessageID() string {
	return fmt.Sprintf("chat_%d", time.Now().UnixNano())
}

// ─── Storage ────────────────────────────────────────────

var mu sync.Mutex

func chatDir(projectDir string) string {
	return filepath.Join(projectDir, ".tasks", "chat")
}

func chatPath(projectDir, agentRole string) string {
	return filepath.Join(chatDir(projectDir), agentRole+".json")
}

func ensureDir(projectDir string) error {
	return os.MkdirAll(chatDir(projectDir), 0755)
}

// LoadMessages returns chat history for an agent
func LoadMessages(projectDir, agentRole string, limit int) ([]ChatMessage, error) {
	mu.Lock()
	defer mu.Unlock()

	data, err := os.ReadFile(chatPath(projectDir, agentRole))
	if err != nil {
		if os.IsNotExist(err) {
			return []ChatMessage{}, nil
		}
		return nil, err
	}

	var msgs []ChatMessage
	if err := json.Unmarshal(data, &msgs); err != nil {
		return nil, err
	}

	if limit > 0 && len(msgs) > limit {
		msgs = msgs[len(msgs)-limit:]
	}
	return msgs, nil
}

// AddMessage appends a message and persists
func AddMessage(projectDir string, msg ChatMessage) error {
	mu.Lock()
	defer mu.Unlock()

	if err := ensureDir(projectDir); err != nil {
		return err
	}

	path := chatPath(projectDir, msg.AgentRole)
	data, err := os.ReadFile(path)
	var msgs []ChatMessage
	if err == nil {
		json.Unmarshal(data, &msgs)
	}

	msgs = append(msgs, msg)

	// Keep last 200 messages per agent
	if len(msgs) > 200 {
		msgs = msgs[len(msgs)-200:]
	}

	out, err := json.MarshalIndent(msgs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, out, 0644)
}

// BuildConversationContext creates a prompt with chat history for the agent
func BuildConversationContext(projectDir, agentRole, userMessage string) (string, error) {
	msgs, err := LoadMessages(projectDir, agentRole, 10) // Last 10 messages for context
	if err != nil {
		msgs = []ChatMessage{}
	}

	prompt := "You are in a direct chat with the user. Answer their question helpfully and concisely.\n"
	prompt += "You can read project files to provide context but do NOT create or modify files unless explicitly asked.\n\n"

	if len(msgs) > 0 {
		prompt += "## Recent Chat History\n"
		for _, m := range msgs {
			if m.Direction == "user_to_agent" {
				prompt += "User: " + m.Content + "\n"
			} else {
				prompt += "You: " + m.Content + "\n"
			}
		}
		prompt += "\n"
	}

	prompt += "## User's Message\n" + userMessage + "\n\n"
	prompt += "Respond directly. Be concise but thorough."

	return prompt, nil
}
