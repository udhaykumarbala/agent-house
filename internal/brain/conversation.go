package brain

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ConversationStore manages chat history per user.
type ConversationStore struct {
	mu       sync.Mutex
	dir      string
	messages map[string][]ChatMessage // userID → messages
}

// NewConversationStore creates a store backed by the filesystem.
func NewConversationStore(dir string) *ConversationStore {
	os.MkdirAll(dir, 0755)
	return &ConversationStore{
		dir:      dir,
		messages: make(map[string][]ChatMessage),
	}
}

// Add appends a message to a user's conversation.
func (cs *ConversationStore) Add(userID, role, content string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	msgs := cs.messages[userID]
	msgs = append(msgs, ChatMessage{
		Role:      role,
		Content:   content,
		Timestamp: time.Now().UnixMilli(),
	})

	// Keep bounded
	if len(msgs) > 50 {
		msgs = msgs[len(msgs)-30:]
	}

	cs.messages[userID] = msgs
	cs.save(userID)
}

// GetRecent returns the last N messages for a user.
func (cs *ConversationStore) GetRecent(userID string, n int) []ChatMessage {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	// Load from disk if not in memory
	if _, ok := cs.messages[userID]; !ok {
		cs.load(userID)
	}

	msgs := cs.messages[userID]
	if len(msgs) <= n {
		return msgs
	}
	return msgs[len(msgs)-n:]
}

// GetAll returns all messages for a user.
func (cs *ConversationStore) GetAll(userID string) []ChatMessage {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if _, ok := cs.messages[userID]; !ok {
		cs.load(userID)
	}
	return cs.messages[userID]
}

func (cs *ConversationStore) save(userID string) {
	data, err := json.MarshalIndent(cs.messages[userID], "", "  ")
	if err != nil {
		return
	}
	path := filepath.Join(cs.dir, userID+".json")
	os.WriteFile(path, data, 0644)
}

func (cs *ConversationStore) load(userID string) {
	path := filepath.Join(cs.dir, userID+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		cs.messages[userID] = []ChatMessage{}
		return
	}
	var msgs []ChatMessage
	if json.Unmarshal(data, &msgs) == nil {
		cs.messages[userID] = msgs
	}
}
