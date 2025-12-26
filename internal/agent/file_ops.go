package agent

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// FileOperation represents a file creation/modification request
type FileOperation struct {
	Path    string
	Content string
	Action  string // "create" or "modify"
}

// FileResult represents the result of a file operation
type FileResult struct {
	Path    string
	Success bool
	Error   error
	Action  string
}

// ParseFileOperations extracts file creation requests from agent response
func ParseFileOperations(response string) []FileOperation {
	var ops []FileOperation

	// Pattern 1: Markdown code blocks with filename
	// ```filename.ext or ```path/to/file.ext
	// content
	// ```
	blockPattern := regexp.MustCompile("(?s)```([a-zA-Z0-9_/.-]+\\.[a-zA-Z0-9]+)\n(.*?)```")
	matches := blockPattern.FindAllStringSubmatch(response, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			path := strings.TrimSpace(match[1])
			content := match[2]

			// Skip if it looks like a language identifier (go, js, python, etc.)
			if isLanguageIdentifier(path) {
				continue
			}

			ops = append(ops, FileOperation{
				Path:    path,
				Content: content,
				Action:  "create",
			})
		}
	}

	// Pattern 2: FILE: path\n```\ncontent\n```
	filePattern := regexp.MustCompile("(?s)FILE:\\s*([^\n]+)\n```[a-z]*\n(.*?)```")
	matches = filePattern.FindAllStringSubmatch(response, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			path := strings.TrimSpace(match[1])
			content := match[2]

			ops = append(ops, FileOperation{
				Path:    path,
				Content: content,
				Action:  "create",
			})
		}
	}

	// Pattern 3: CREATE_FILE: path\ncontent...END_FILE
	createPattern := regexp.MustCompile("(?s)CREATE_FILE:\\s*([^\n]+)\n(.*?)END_FILE")
	matches = createPattern.FindAllStringSubmatch(response, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			path := strings.TrimSpace(match[1])
			content := strings.TrimSpace(match[2])

			ops = append(ops, FileOperation{
				Path:    path,
				Content: content,
				Action:  "create",
			})
		}
	}

	return ops
}

// isLanguageIdentifier checks if a string looks like a code block language
func isLanguageIdentifier(s string) bool {
	languages := []string{
		"go", "golang", "javascript", "js", "typescript", "ts",
		"python", "py", "java", "rust", "c", "cpp", "csharp", "cs",
		"ruby", "rb", "php", "swift", "kotlin", "scala",
		"html", "css", "scss", "sass", "less",
		"json", "yaml", "yml", "xml", "toml",
		"sql", "graphql", "markdown", "md",
		"bash", "sh", "shell", "zsh", "powershell", "ps1",
		"dockerfile", "makefile", "nginx", "apache",
		"plaintext", "text", "txt", "diff", "log",
	}

	s = strings.ToLower(s)
	for _, lang := range languages {
		if s == lang {
			return true
		}
	}
	return false
}

// ExecuteFileOperations performs file operations with permission checks
func (a *Agent) ExecuteFileOperations(ops []FileOperation, projectDir string) []FileResult {
	var results []FileResult

	for _, op := range ops {
		result := FileResult{
			Path:   op.Path,
			Action: op.Action,
		}

		// Check permissions
		fullPath := filepath.Join(projectDir, op.Path)

		var err error
		switch op.Action {
		case "create":
			err = a.CheckCreatePermission(op.Path)
		case "modify":
			if !a.CanModify(op.Path) {
				err = PermissionError{
					Agent:  a.Name,
					Action: "modify",
					Path:   op.Path,
					Reason: "modification not allowed",
				}
			}
		}

		if err != nil {
			result.Success = false
			result.Error = err
			results = append(results, result)
			continue
		}

		// Perform the operation
		err = writeFile(fullPath, op.Content)
		if err != nil {
			result.Success = false
			result.Error = err
		} else {
			result.Success = true
		}

		results = append(results, result)
	}

	return results
}

// writeFile creates a file with the given content
func writeFile(path, content string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write the file
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}

	return nil
}

// ProcessWithFileOps processes a task and executes any file operations
func (a *Agent) ProcessWithFileOps(task string, projectDir string) (*Response, []FileResult, error) {
	// Process the task with project directory so Claude can write files directly
	response, err := a.ProcessInDir(task, projectDir)
	if err != nil {
		return nil, nil, err
	}

	// Also parse FILE: blocks from response as fallback
	ops := ParseFileOperations(response.Content)

	// Execute any FILE: block operations (in case Claude used that format)
	var results []FileResult
	if len(ops) > 0 {
		results = a.ExecuteFileOperations(ops, projectDir)
	}

	// Add file paths to response
	for _, op := range ops {
		response.Files = append(response.Files, op.Path)
	}

	return response, results, nil
}
