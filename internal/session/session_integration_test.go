package session

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestSessionRealClaude spawns a real Claude Code session, sends a simple file-creation task,
// and inspects every NDJSON event to diagnose tool visibility issues.
//
// Run with: go test ./internal/session/ -run TestSessionRealClaude -v -timeout 120s
func TestSessionRealClaude(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION") != "1" {
		t.Skip("Skipping integration test (set RUN_INTEGRATION=1 to run)")
	}

	// Create temp dir
	tmpDir := t.TempDir()

	config := SessionConfig{
		AgentRole:      "senior_dev",
		AgentName:      "Senior Developer",
		SystemPrompt:   "You are a developer. Create files as requested. Be concise.",
		WorkDir:        tmpDir,
		PermissionMode: "bypassPermissions",
		ProjectID:      "test-integration",
		MaxTurns:       10,
	}

	sess := NewAgentSession(config)
	defer sess.Terminate()

	// Subscribe to ALL events
	evCh, unsub := sess.Subscribe()
	defer unsub()

	// Collect events in background
	var allEvents []AgentEvent
	done := make(chan struct{})
	go func() {
		for ev := range evCh {
			allEvents = append(allEvents, ev)
			t.Logf("[EVENT] type=%-16s tool=%-8s content_len=%d", ev.Type, ev.ToolName, len(ev.Content)+len(ev.Output))
		}
		close(done)
	}()

	// Send task
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	text, events, err := sess.SendTask(ctx, "Create a file called hello.txt with the content 'Hello World'. Just create the file, nothing else.")
	if err != nil {
		t.Fatalf("SendTask failed: %v", err)
	}

	t.Logf("\n=== RESULTS ===")
	t.Logf("Response text length: %d", len(text))
	t.Logf("Events collected: %d", len(events))
	t.Logf("Stats: %+v", sess.GetStats())

	// Categorize events
	counts := map[string]int{}
	for _, ev := range events {
		counts[ev.Type]++
	}
	t.Logf("Event breakdown: %v", counts)

	// Check file was created
	helloPath := filepath.Join(tmpDir, "hello.txt")
	if _, err := os.Stat(helloPath); os.IsNotExist(err) {
		t.Errorf("hello.txt was NOT created at %s", helloPath)
	} else {
		content, _ := os.ReadFile(helloPath)
		t.Logf("hello.txt content: %q", string(content))
	}

	// Verify tool events
	if counts["tool_use"] == 0 {
		t.Errorf("PROBLEM: 0 tool_use events — tools are invisible!")
		t.Log("This means Claude Code stream-json doesn't emit tool_use content blocks in this permission mode")
	}
	if counts["turn_complete"] == 0 {
		t.Error("PROBLEM: no turn_complete event")
	}

	// Print all events for debugging
	t.Logf("\n=== ALL EVENTS ===")
	for i, ev := range events {
		evJSON, _ := json.MarshalIndent(ev, "", "  ")
		t.Logf("Event %d:\n%s", i, string(evJSON))
	}
}

// TestSessionNDJSONDebug adds raw NDJSON logging to see exactly what Claude sends.
// This is the key diagnostic for BUG-01.
//
// Run with: go test ./internal/session/ -run TestSessionNDJSONDebug -v -timeout 120s
func TestSessionNDJSONDebug(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION") != "1" {
		t.Skip("Skipping integration test (set RUN_INTEGRATION=1 to run)")
	}

	tmpDir := t.TempDir()

	// Build command manually to capture raw NDJSON
	config := SessionConfig{
		AgentRole:      "test_dev",
		AgentName:      "Test Dev",
		SystemPrompt:   "You are a developer. Be extremely concise. Just do what's asked.",
		WorkDir:        tmpDir,
		PermissionMode: "bypassPermissions",
		ProjectID:      "ndjson-debug",
		MaxTurns:       5,
	}

	sess := NewAgentSession(config)
	defer sess.Terminate()

	// Subscribe
	evCh, unsub := sess.Subscribe()
	defer unsub()

	var events []AgentEvent
	go func() {
		for ev := range evCh {
			events = append(events, ev)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	text, returnedEvents, err := sess.SendTask(ctx, "Write a file called test.txt containing 'hello'. That's all.")
	if err != nil {
		t.Fatalf("SendTask error: %v", err)
	}

	t.Logf("Response: %q", text)
	t.Logf("Returned events: %d", len(returnedEvents))

	// Print stats
	stats := sess.GetStats()
	t.Logf("Stats: turns=%d tools=%d cost=$%.4f", stats.TotalTurns, stats.TotalToolCalls, stats.TotalCostUSD)
	t.Logf("Tool breakdown: %v", stats.ToolBreakdown)

	// Check file
	testPath := filepath.Join(tmpDir, "test.txt")
	if data, err := os.ReadFile(testPath); err == nil {
		t.Logf("test.txt exists: %q", string(data))
	} else {
		t.Logf("test.txt NOT found: %v", err)
	}

	// Dump all events
	for i, ev := range returnedEvents {
		fmt.Printf("EVENT[%d] type=%s tool=%s content_len=%d output_len=%d\n",
			i, ev.Type, ev.ToolName, len(ev.Content), len(ev.Output))
	}
}
