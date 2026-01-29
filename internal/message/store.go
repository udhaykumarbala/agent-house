package message

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// Store manages message persistence
type Store struct {
	messages  []*Message
	mu        sync.RWMutex
	listeners []func(*Message) // Callbacks for new messages
}

// NewStore creates a new message store
func NewStore() *Store {
	return &Store{
		messages:  make([]*Message, 0),
		listeners: make([]func(*Message), 0),
	}
}

// Add stores a new message
func (s *Store) Add(msg *Message) {
	s.mu.Lock()
	s.messages = append(s.messages, msg)
	s.mu.Unlock()

	// Notify listeners
	for _, listener := range s.listeners {
		go listener(msg)
	}
}

// OnMessage registers a callback for new messages
func (s *Store) OnMessage(callback func(*Message)) {
	s.mu.Lock()
	s.listeners = append(s.listeners, callback)
	s.mu.Unlock()
}

// GetAll returns all messages
func (s *Store) GetAll() []*Message {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*Message, len(s.messages))
	copy(result, s.messages)
	return result
}

// GetByAgent returns messages from or to a specific agent
func (s *Store) GetByAgent(agentID string) []*Message {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*Message
	for _, msg := range s.messages {
		if msg.From == agentID || msg.To == agentID {
			result = append(result, msg)
		}
	}
	return result
}

// GetByType returns messages of a specific type
func (s *Store) GetByType(msgType MessageType) []*Message {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*Message
	for _, msg := range s.messages {
		if msg.Type == msgType {
			result = append(result, msg)
		}
	}
	return result
}

// GetByProject returns messages for a specific project
func (s *Store) GetByProject(projectID string) []*Message {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*Message
	for _, msg := range s.messages {
		if msg.Metadata.ProjectID == projectID {
			result = append(result, msg)
		}
	}
	return result
}

// GetAfter returns messages after a specific time
func (s *Store) GetAfter(t time.Time) []*Message {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*Message
	for _, msg := range s.messages {
		if msg.Timestamp.After(t) {
			result = append(result, msg)
		}
	}
	return result
}

// GetByID returns a message by its ID
func (s *Store) GetByID(id string) *Message {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, msg := range s.messages {
		if msg.ID == id {
			return msg
		}
	}
	return nil
}

// Count returns the total number of messages
func (s *Store) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.messages)
}

// Clear removes all messages
func (s *Store) Clear() {
	s.mu.Lock()
	s.messages = make([]*Message, 0)
	s.mu.Unlock()
}

// GetConversation returns messages sorted by timestamp
func (s *Store) GetConversation() []*Message {
	messages := s.GetAll()
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].Timestamp.Before(messages[j].Timestamp)
	})
	return messages
}

// SaveToFile persists messages to a JSON file
func (s *Store) SaveToFile(path string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(s.messages, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// LoadFromFile loads messages from a JSON file
func (s *Store) LoadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No file yet, that's OK
		}
		return err
	}

	var messages []*Message
	if err := json.Unmarshal(data, &messages); err != nil {
		return err
	}

	s.mu.Lock()
	s.messages = messages
	s.mu.Unlock()

	return nil
}

// GetByTaskID returns messages for a specific task
func (s *Store) GetByTaskID(taskID string) []*Message {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*Message
	for _, msg := range s.messages {
		if msg.Metadata.TaskID == taskID {
			result = append(result, msg)
		}
	}

	// Sort by timestamp
	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.Before(result[j].Timestamp)
	})

	return result
}
