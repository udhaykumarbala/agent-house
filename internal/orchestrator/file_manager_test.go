package orchestrator

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"pty-claude-test/internal/agent"
)

func TestFileManagerBasicWrite(t *testing.T) {
	tmpDir := t.TempDir()
	fm := NewConcurrentFileManager()
	defer fm.Shutdown()

	testPath := filepath.Join(tmpDir, "test.txt")
	content := "test content"

	err := fm.WriteFile(testPath, content, agent.RolePM)
	if err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Verify file exists and has correct content
	data, err := os.ReadFile(testPath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if string(data) != content {
		t.Errorf("Expected content %q, got %q", content, string(data))
	}

	// Verify tracking
	exists, creator := fm.IsFileCreated(testPath)
	if !exists {
		t.Error("File should be tracked as created")
	}
	if creator != agent.RolePM {
		t.Errorf("Expected creator PM, got %s", creator)
	}
}

func TestFileManagerDuplicatePrevention(t *testing.T) {
	tmpDir := t.TempDir()
	fm := NewConcurrentFileManager()
	defer fm.Shutdown()

	testPath := filepath.Join(tmpDir, "duplicate.txt")

	// First write should succeed
	err := fm.WriteFile(testPath, "first", agent.RolePM)
	if err != nil {
		t.Fatalf("First write failed: %v", err)
	}

	// Second write by different agent should fail
	err = fm.WriteFile(testPath, "second", agent.RoleUX)
	if err == nil {
		t.Error("Expected error for duplicate file write")
	}

	// Verify file still has first content
	data, _ := os.ReadFile(testPath)
	if string(data) != "first" {
		t.Errorf("File content changed, expected 'first', got %q", string(data))
	}

	// Verify creator is still PM
	_, creator := fm.IsFileCreated(testPath)
	if creator != agent.RolePM {
		t.Errorf("Expected creator PM, got %s", creator)
	}
}

func TestFileManagerConcurrentWrites(t *testing.T) {
	tmpDir := t.TempDir()
	fm := NewConcurrentFileManager()
	defer fm.Shutdown()

	numAgents := 5
	roles := []agent.Role{agent.RolePM, agent.RoleUX, agent.RoleUI, agent.RoleSecurity, agent.RoleArchitect}

	// All agents try to write the same file
	testPath := filepath.Join(tmpDir, "concurrent.txt")

	var wg sync.WaitGroup
	successCount := 0
	errorCount := 0
	var mu sync.Mutex

	for i := 0; i < numAgents; i++ {
		wg.Add(1)
		go func(role agent.Role, content string) {
			defer wg.Done()
			err := fm.WriteFile(testPath, content, role)
			mu.Lock()
			if err == nil {
				successCount++
			} else {
				errorCount++
			}
			mu.Unlock()
		}(roles[i], string(roles[i]))
	}

	wg.Wait()

	// Only one should succeed
	if successCount != 1 {
		t.Errorf("Expected exactly 1 success, got %d", successCount)
	}

	// Others should fail
	if errorCount != numAgents-1 {
		t.Errorf("Expected %d errors, got %d", numAgents-1, errorCount)
	}

	// File should exist with one agent's content
	data, err := os.ReadFile(testPath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	// Content should be one of the agent roles
	foundValidContent := false
	for _, role := range roles {
		if string(data) == string(role) {
			foundValidContent = true
			break
		}
	}
	if !foundValidContent {
		t.Errorf("File content %q doesn't match any agent role", string(data))
	}
}

func TestFileManagerMultipleFiles(t *testing.T) {
	tmpDir := t.TempDir()
	fm := NewConcurrentFileManager()
	defer fm.Shutdown()

	roles := []agent.Role{agent.RolePM, agent.RoleUX, agent.RoleUI}
	var wg sync.WaitGroup

	// Each agent writes its own file
	for i, role := range roles {
		wg.Add(1)
		go func(role agent.Role, idx int) {
			defer wg.Done()
			path := filepath.Join(tmpDir, string(role)+".txt")
			err := fm.WriteFile(path, string(role)+" content", role)
			if err != nil {
				t.Errorf("Agent %s failed to write: %v", role, err)
			}
		}(role, i)
	}

	wg.Wait()

	// Verify all files created
	files := fm.GetCreatedFiles()
	if len(files) != len(roles) {
		t.Errorf("Expected %d files, got %d", len(roles), len(files))
	}

	// Verify each file
	for _, role := range roles {
		path := filepath.Join(tmpDir, string(role)+".txt")
		exists, creator := fm.IsFileCreated(path)
		if !exists {
			t.Errorf("File for %s not tracked", role)
		}
		if creator != role {
			t.Errorf("Expected creator %s, got %s", role, creator)
		}
	}
}

func TestFileManagerNestedDirectories(t *testing.T) {
	tmpDir := t.TempDir()
	fm := NewConcurrentFileManager()
	defer fm.Shutdown()

	nestedPath := filepath.Join(tmpDir, "a", "b", "c", "test.txt")

	err := fm.WriteFile(nestedPath, "nested content", agent.RolePM)
	if err != nil {
		t.Fatalf("Failed to write nested file: %v", err)
	}

	// Verify file exists
	data, err := os.ReadFile(nestedPath)
	if err != nil {
		t.Fatalf("Failed to read nested file: %v", err)
	}

	if string(data) != "nested content" {
		t.Errorf("Expected 'nested content', got %q", string(data))
	}
}

func TestFileManagerMarkFileCreated(t *testing.T) {
	fm := NewConcurrentFileManager()
	defer fm.Shutdown()

	testPath := "/tmp/manual.txt"

	// Manually mark as created
	fm.MarkFileCreated(testPath, agent.RoleArchitect)

	// Verify tracking
	exists, creator := fm.IsFileCreated(testPath)
	if !exists {
		t.Error("File should be tracked")
	}
	if creator != agent.RoleArchitect {
		t.Errorf("Expected creator Architect, got %s", creator)
	}

	// Subsequent write should fail
	err := fm.WriteFile(testPath, "content", agent.RolePM)
	if err == nil {
		t.Error("Expected error when writing manually marked file")
	}
}

func TestFileManagerClear(t *testing.T) {
	tmpDir := t.TempDir()
	fm := NewConcurrentFileManager()
	defer fm.Shutdown()

	// Write a file
	testPath := filepath.Join(tmpDir, "clear.txt")
	fm.WriteFile(testPath, "content", agent.RolePM)

	// Verify tracked
	exists, _ := fm.IsFileCreated(testPath)
	if !exists {
		t.Error("File should be tracked")
	}

	// Clear tracking
	fm.Clear()

	// Should not be tracked anymore
	exists, _ = fm.IsFileCreated(testPath)
	if exists {
		t.Error("File should not be tracked after clear")
	}

	// Should be able to write again
	err := fm.WriteFile(testPath, "new content", agent.RoleUX)
	if err != nil {
		t.Errorf("Should be able to write after clear: %v", err)
	}
}

func BenchmarkFileManagerConcurrentWrites(b *testing.B) {
	tmpDir := b.TempDir()
	fm := NewConcurrentFileManager()
	defer fm.Shutdown()

	roles := []agent.Role{agent.RolePM, agent.RoleUX, agent.RoleUI, agent.RoleSecurity, agent.RoleArchitect}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		for j, role := range roles {
			wg.Add(1)
			go func(role agent.Role, idx int) {
				defer wg.Done()
				path := filepath.Join(tmpDir, "bench", string(role), "file.txt")
				fm.WriteFile(path, "content", role)
			}(role, j)
		}
		wg.Wait()
		fm.Clear()
	}
}
