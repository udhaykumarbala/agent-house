package message

import (
	"encoding/json"
	"time"
)

// MessageType represents the type of message
type MessageType string

const (
	TypeTask       MessageType = "task"        // Initial task from user
	TypeResponse   MessageType = "response"    // Agent's response
	TypeDelegate   MessageType = "delegate"    // Delegation to another agent
	TypeFileCreate MessageType = "file_create" // Request to create a file
	TypeFileModify MessageType = "file_modify" // Request to modify a file
	TypeComplete   MessageType = "complete"    // Task completion signal
	TypeError      MessageType = "error"       // Error message
	TypeSystem     MessageType = "system"      // System/orchestrator message (phase changes, etc.)
)

// Message represents communication between agents or from user
type Message struct {
	ID        string      `json:"id"`
	Type      MessageType `json:"type"`
	From      string      `json:"from"`      // Agent ID or "user"
	To        string      `json:"to"`        // Agent ID or "all"
	Content   string      `json:"content"`   // Message content
	Metadata  Metadata    `json:"metadata"`  // Additional data
	Timestamp time.Time   `json:"timestamp"`
}

// Metadata holds additional message information
type Metadata struct {
	ProjectID   string                 `json:"project_id,omitempty"`
	TaskID      string                 `json:"task_id,omitempty"`
	DelegateTo  []string               `json:"delegate_to,omitempty"`  // For delegate messages
	FilePath    string                 `json:"file_path,omitempty"`    // For file operations
	FileContent string                 `json:"file_content,omitempty"` // For file creation
	Priority    string                 `json:"priority,omitempty"`     // high, medium, low
	Tags        []string               `json:"tags,omitempty"`         // Labels for filtering
	Extra       map[string]interface{} `json:"extra,omitempty"`        // For lifecycle events
}

// NewMessage creates a new message with generated ID and timestamp
func NewMessage(msgType MessageType, from, to, content string) *Message {
	return &Message{
		ID:        generateID(),
		Type:      msgType,
		From:      from,
		To:        to,
		Content:   content,
		Timestamp: time.Now(),
	}
}

// NewTaskMessage creates a task message from user
func NewTaskMessage(projectID, taskContent string) *Message {
	msg := NewMessage(TypeTask, "user", "ceo", taskContent)
	msg.Metadata.ProjectID = projectID
	msg.Metadata.TaskID = generateID()
	return msg
}

// NewDelegateMessage creates a delegation message
func NewDelegateMessage(from string, delegateTo []string, reason string) *Message {
	msg := NewMessage(TypeDelegate, from, delegateTo[0], reason)
	msg.Metadata.DelegateTo = delegateTo
	return msg
}

// NewFileCreateMessage creates a file creation request
func NewFileCreateMessage(from, filePath, content string) *Message {
	msg := NewMessage(TypeFileCreate, from, "system", "Create file: "+filePath)
	msg.Metadata.FilePath = filePath
	msg.Metadata.FileContent = content
	return msg
}

// ToJSON serializes the message to JSON
func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// FromJSON deserializes a message from JSON
func FromJSON(data []byte) (*Message, error) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// generateID creates a simple unique ID
func generateID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(6)
}

// randomString generates a random alphanumeric string
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
		time.Sleep(time.Nanosecond) // Ensure uniqueness
	}
	return string(b)
}

// String returns a human-readable representation
func (m *Message) String() string {
	return m.From + " → " + m.To + ": " + truncate(m.Content, 50)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
