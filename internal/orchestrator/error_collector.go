package orchestrator

import (
	"fmt"
	"strings"
	"sync"

	"pty-claude-test/internal/agent"
)

// AgentError represents an error from a specific agent
type AgentError struct {
	Role  agent.Role
	Error error
}

// ErrorCollector aggregates errors from parallel agent execution
type ErrorCollector struct {
	mu     sync.Mutex
	errors []AgentError
}

// NewErrorCollector creates a new error collector
func NewErrorCollector() *ErrorCollector {
	return &ErrorCollector{
		errors: make([]AgentError, 0),
	}
}

// Add adds an error for a specific agent
func (ec *ErrorCollector) Add(role agent.Role, err error) {
	if err == nil {
		return
	}

	ec.mu.Lock()
	defer ec.mu.Unlock()

	ec.errors = append(ec.errors, AgentError{
		Role:  role,
		Error: err,
	})
}

// HasErrors returns true if any errors were collected
func (ec *ErrorCollector) HasErrors() bool {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	return len(ec.errors) > 0
}

// Count returns the number of errors
func (ec *ErrorCollector) Count() int {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	return len(ec.errors)
}

// Errors returns a copy of all collected errors
func (ec *ErrorCollector) Errors() []AgentError {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	// Return a copy
	result := make([]AgentError, len(ec.errors))
	copy(result, ec.errors)
	return result
}

// Error implements the error interface, combining all errors
func (ec *ErrorCollector) Error() string {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	if len(ec.errors) == 0 {
		return "no errors"
	}

	if len(ec.errors) == 1 {
		return fmt.Sprintf("%s: %v", ec.errors[0].Role, ec.errors[0].Error)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d agent errors:\n", len(ec.errors)))
	for _, ae := range ec.errors {
		sb.WriteString(fmt.Sprintf("  - %s: %v\n", ae.Role, ae.Error))
	}
	return sb.String()
}

// Clear resets the error collector
func (ec *ErrorCollector) Clear() {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	ec.errors = make([]AgentError, 0)
}
