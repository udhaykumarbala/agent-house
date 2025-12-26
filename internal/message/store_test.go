package message

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestStoreAdd(t *testing.T) {
	store := NewStore()

	msg := NewMessage(TypeTask, "user", "ceo", "Build something")
	store.Add(msg)

	if store.Count() != 1 {
		t.Errorf("expected count 1, got %d", store.Count())
	}
}

func TestStoreGetAll(t *testing.T) {
	store := NewStore()

	store.Add(NewMessage(TypeTask, "user", "ceo", "Task 1"))
	store.Add(NewMessage(TypeResponse, "ceo", "user", "Response 1"))

	all := store.GetAll()
	if len(all) != 2 {
		t.Errorf("expected 2 messages, got %d", len(all))
	}
}

func TestStoreGetByAgent(t *testing.T) {
	store := NewStore()

	store.Add(NewMessage(TypeTask, "user", "ceo", "Task"))
	store.Add(NewMessage(TypeResponse, "ceo", "pm", "Response"))
	store.Add(NewMessage(TypeDelegate, "pm", "ux", "Delegate"))

	ceoMessages := store.GetByAgent("ceo")
	if len(ceoMessages) != 2 {
		t.Errorf("expected 2 messages for CEO, got %d", len(ceoMessages))
	}

	pmMessages := store.GetByAgent("pm")
	if len(pmMessages) != 2 {
		t.Errorf("expected 2 messages for PM, got %d", len(pmMessages))
	}
}

func TestStoreGetByType(t *testing.T) {
	store := NewStore()

	store.Add(NewMessage(TypeTask, "user", "ceo", "Task 1"))
	store.Add(NewMessage(TypeTask, "user", "pm", "Task 2"))
	store.Add(NewMessage(TypeResponse, "ceo", "user", "Response"))

	tasks := store.GetByType(TypeTask)
	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}

	responses := store.GetByType(TypeResponse)
	if len(responses) != 1 {
		t.Errorf("expected 1 response, got %d", len(responses))
	}
}

func TestStoreGetByProject(t *testing.T) {
	store := NewStore()

	msg1 := NewTaskMessage("project-a", "Task A")
	msg2 := NewTaskMessage("project-b", "Task B")
	msg3 := NewTaskMessage("project-a", "Task A2")

	store.Add(msg1)
	store.Add(msg2)
	store.Add(msg3)

	projectA := store.GetByProject("project-a")
	if len(projectA) != 2 {
		t.Errorf("expected 2 messages for project-a, got %d", len(projectA))
	}
}

func TestStoreGetByID(t *testing.T) {
	store := NewStore()

	msg := NewMessage(TypeTask, "user", "ceo", "Find me")
	store.Add(msg)

	found := store.GetByID(msg.ID)
	if found == nil {
		t.Error("expected to find message by ID")
	}
	if found.Content != "Find me" {
		t.Errorf("wrong message found: %s", found.Content)
	}

	notFound := store.GetByID("nonexistent")
	if notFound != nil {
		t.Error("expected nil for nonexistent ID")
	}
}

func TestStoreGetAfter(t *testing.T) {
	store := NewStore()

	oldMsg := NewMessage(TypeTask, "user", "ceo", "Old")
	oldMsg.Timestamp = time.Now().Add(-1 * time.Hour)
	store.Add(oldMsg)

	time.Sleep(10 * time.Millisecond) // Ensure timestamp difference

	newMsg := NewMessage(TypeTask, "user", "pm", "New")
	store.Add(newMsg)

	cutoff := time.Now().Add(-30 * time.Minute)
	recent := store.GetAfter(cutoff)

	if len(recent) != 1 {
		t.Errorf("expected 1 recent message, got %d", len(recent))
	}
}

func TestStoreClear(t *testing.T) {
	store := NewStore()

	store.Add(NewMessage(TypeTask, "a", "b", "msg1"))
	store.Add(NewMessage(TypeTask, "a", "b", "msg2"))

	if store.Count() != 2 {
		t.Errorf("expected 2 messages before clear")
	}

	store.Clear()

	if store.Count() != 0 {
		t.Errorf("expected 0 messages after clear, got %d", store.Count())
	}
}

func TestStoreGetConversation(t *testing.T) {
	store := NewStore()

	msg1 := NewMessage(TypeTask, "user", "ceo", "First")
	msg1.Timestamp = time.Now().Add(-2 * time.Second)

	msg2 := NewMessage(TypeResponse, "ceo", "user", "Second")
	msg2.Timestamp = time.Now().Add(-1 * time.Second)

	msg3 := NewMessage(TypeDelegate, "ceo", "pm", "Third")
	msg3.Timestamp = time.Now()

	// Add out of order
	store.Add(msg3)
	store.Add(msg1)
	store.Add(msg2)

	conv := store.GetConversation()
	if len(conv) != 3 {
		t.Errorf("expected 3 messages, got %d", len(conv))
	}
	if conv[0].Content != "First" {
		t.Errorf("expected First first, got %s", conv[0].Content)
	}
	if conv[2].Content != "Third" {
		t.Errorf("expected Third last, got %s", conv[2].Content)
	}
}

func TestStoreSaveAndLoad(t *testing.T) {
	// Create temp file
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "messages.json")

	// Create and populate store
	store1 := NewStore()
	store1.Add(NewMessage(TypeTask, "user", "ceo", "Saved message"))
	store1.Add(NewMessage(TypeResponse, "ceo", "user", "Response"))

	// Save
	if err := store1.SaveToFile(path); err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("file was not created")
	}

	// Load into new store
	store2 := NewStore()
	if err := store2.LoadFromFile(path); err != nil {
		t.Fatalf("LoadFromFile failed: %v", err)
	}

	if store2.Count() != 2 {
		t.Errorf("expected 2 messages after load, got %d", store2.Count())
	}

	messages := store2.GetAll()
	if messages[0].Content != "Saved message" {
		t.Errorf("content not preserved: %s", messages[0].Content)
	}
}

func TestStoreOnMessage(t *testing.T) {
	store := NewStore()

	received := make(chan *Message, 1)
	store.OnMessage(func(msg *Message) {
		received <- msg
	})

	store.Add(NewMessage(TypeTask, "user", "ceo", "Notify me"))

	select {
	case msg := <-received:
		if msg.Content != "Notify me" {
			t.Errorf("wrong message received: %s", msg.Content)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("listener not called")
	}
}
