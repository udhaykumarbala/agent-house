package agent

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseFileOperations_MarkdownCodeBlock(t *testing.T) {
	response := `Here's the implementation:

` + "```docs/requirements.md" + `
# Requirements

## User Stories
1. As a user, I want to create todos
2. As a user, I want to mark todos complete
` + "```" + `

That should cover the basics.`

	ops := ParseFileOperations(response)

	if len(ops) != 1 {
		t.Fatalf("Expected 1 operation, got %d", len(ops))
	}

	if ops[0].Path != "docs/requirements.md" {
		t.Errorf("Expected path 'docs/requirements.md', got '%s'", ops[0].Path)
	}

	if !strings.Contains(ops[0].Content, "User Stories") {
		t.Error("Content should contain 'User Stories'")
	}

	if ops[0].Action != "create" {
		t.Errorf("Expected action 'create', got '%s'", ops[0].Action)
	}
}

func TestParseFileOperations_FilePrefix(t *testing.T) {
	response := `FILE: src/utils.go
` + "```go" + `
package utils

func Add(a, b int) int {
	return a + b
}
` + "```"

	ops := ParseFileOperations(response)

	if len(ops) != 1 {
		t.Fatalf("Expected 1 operation, got %d", len(ops))
	}

	if ops[0].Path != "src/utils.go" {
		t.Errorf("Expected path 'src/utils.go', got '%s'", ops[0].Path)
	}

	if !strings.Contains(ops[0].Content, "func Add") {
		t.Error("Content should contain 'func Add'")
	}
}

func TestParseFileOperations_CreateFileBlock(t *testing.T) {
	response := `CREATE_FILE: config/settings.yaml
database:
  host: localhost
  port: 5432
END_FILE`

	ops := ParseFileOperations(response)

	if len(ops) != 1 {
		t.Fatalf("Expected 1 operation, got %d", len(ops))
	}

	if ops[0].Path != "config/settings.yaml" {
		t.Errorf("Expected path 'config/settings.yaml', got '%s'", ops[0].Path)
	}

	if !strings.Contains(ops[0].Content, "localhost") {
		t.Error("Content should contain 'localhost'")
	}
}

func TestParseFileOperations_MultipleFiles(t *testing.T) {
	response := `FILE: src/main.go
` + "```go" + `
package main
func main() {}
` + "```" + `

FILE: src/utils.go
` + "```go" + `
package main
func helper() {}
` + "```"

	ops := ParseFileOperations(response)

	if len(ops) != 2 {
		t.Fatalf("Expected 2 operations, got %d", len(ops))
	}

	if ops[0].Path != "src/main.go" {
		t.Errorf("Expected first path 'src/main.go', got '%s'", ops[0].Path)
	}

	if ops[1].Path != "src/utils.go" {
		t.Errorf("Expected second path 'src/utils.go', got '%s'", ops[1].Path)
	}
}

func TestParseFileOperations_SkipsLanguageIdentifiers(t *testing.T) {
	response := `Here's some code:

` + "```go" + `
package main
func main() {}
` + "```" + `

And some JavaScript:

` + "```javascript" + `
console.log("hello");
` + "```"

	ops := ParseFileOperations(response)

	if len(ops) != 0 {
		t.Errorf("Expected 0 operations (language identifiers should be skipped), got %d", len(ops))
	}
}

func TestIsLanguageIdentifier(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"go", true},
		{"javascript", true},
		{"python", true},
		{"main.go", false},
		{"src/utils.js", false},
		{"docs/readme.md", false},
		{"json", true},
		{"config.json", false},
	}

	for _, tt := range tests {
		result := isLanguageIdentifier(tt.input)
		if result != tt.expected {
			t.Errorf("isLanguageIdentifier(%q) = %v, expected %v", tt.input, result, tt.expected)
		}
	}
}

func TestExecuteFileOperations(t *testing.T) {
	// Create a temp directory for testing
	tmpDir, err := os.MkdirTemp("", "agent-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a senior dev agent (can create any file)
	dev := &Agent{
		Role: RoleSeniorDev,
		Name: "Senior Dev",
	}

	ops := []FileOperation{
		{
			Path:    "src/main.go",
			Content: "package main\n\nfunc main() {}\n",
			Action:  "create",
		},
		{
			Path:    "docs/readme.md",
			Content: "# README\n",
			Action:  "create",
		},
	}

	results := dev.ExecuteFileOperations(ops, tmpDir)

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	// Both should succeed for senior dev
	for i, result := range results {
		if !result.Success {
			t.Errorf("Result %d should be successful: %v", i, result.Error)
		}
	}

	// Verify files exist
	mainPath := filepath.Join(tmpDir, "src/main.go")
	if _, err := os.Stat(mainPath); os.IsNotExist(err) {
		t.Error("src/main.go was not created")
	}

	readmePath := filepath.Join(tmpDir, "docs/readme.md")
	if _, err := os.Stat(readmePath); os.IsNotExist(err) {
		t.Error("docs/readme.md was not created")
	}

	// Verify content
	content, _ := os.ReadFile(mainPath)
	if !strings.Contains(string(content), "func main()") {
		t.Error("main.go content is incorrect")
	}
}

func TestExecuteFileOperations_PermissionDenied(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "agent-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a PM agent (can only create .md files in docs/)
	pm := &Agent{
		Role: RolePM,
		Name: "PM",
	}

	ops := []FileOperation{
		{
			Path:    "src/main.go", // PM cannot create Go files
			Content: "package main",
			Action:  "create",
		},
	}

	results := pm.ExecuteFileOperations(ops, tmpDir)

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0].Success {
		t.Error("PM should not be able to create Go files")
	}

	if results[0].Error == nil {
		t.Error("Expected permission error")
	}

	// Verify file was NOT created
	mainPath := filepath.Join(tmpDir, "src/main.go")
	if _, err := os.Stat(mainPath); !os.IsNotExist(err) {
		t.Error("File should not have been created")
	}
}

func TestExecuteFileOperations_JuniorDevRestrictions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "agent-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	junior := &Agent{
		Role: RoleJuniorDev,
		Name: "Junior Dev",
	}

	ops := []FileOperation{
		{
			Path:    "src/utils.go", // Junior can create in src/
			Content: "package main",
			Action:  "create",
		},
		{
			Path:    "config/settings.yaml", // Junior cannot create in config/
			Content: "key: value",
			Action:  "create",
		},
	}

	results := junior.ExecuteFileOperations(ops, tmpDir)

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	// First should succeed
	if !results[0].Success {
		t.Errorf("Junior should be able to create in src/: %v", results[0].Error)
	}

	// Second should fail
	if results[1].Success {
		t.Error("Junior should NOT be able to create in config/")
	}
}

func TestProcessWithFileOps(t *testing.T) {
	// This test would require mocking Claude, so we just verify the function exists
	// and returns proper types
	pm := &Agent{
		Role:         RolePM,
		Name:         "PM",
		SystemPrompt: "You are a PM",
	}

	// Verify the method exists on Agent
	_ = pm.ProcessWithFileOps
}
