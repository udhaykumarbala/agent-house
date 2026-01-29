package orchestrator

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"pty-claude-test/internal/agent"
)

// FileWriteJob represents a file write operation
type FileWriteJob struct {
	Path      string
	Content   string
	AgentRole agent.Role
	Response  chan FileWriteResult
}

// FileWriteResult represents the outcome of a file write
type FileWriteResult struct {
	Path      string
	Success   bool
	Error     error
	CreatedBy agent.Role // The agent that successfully created the file
}

// ConcurrentFileManager handles thread-safe file operations with deduplication
type ConcurrentFileManager struct {
	mu            sync.RWMutex
	createdFiles  map[string]agent.Role // path -> agent who created it
	pendingWrites chan FileWriteJob
	done          chan struct{}
	wg            sync.WaitGroup
}

// NewConcurrentFileManager creates a new file manager
func NewConcurrentFileManager() *ConcurrentFileManager {
	fm := &ConcurrentFileManager{
		createdFiles:  make(map[string]agent.Role),
		pendingWrites: make(chan FileWriteJob, 100),
		done:          make(chan struct{}),
	}

	// Start the writer goroutine
	fm.wg.Add(1)
	go fm.writerLoop()

	return fm
}

// writerLoop processes file write operations serially
func (fm *ConcurrentFileManager) writerLoop() {
	defer fm.wg.Done()

	for {
		select {
		case <-fm.done:
			// Process remaining jobs
			for {
				select {
				case job := <-fm.pendingWrites:
					fm.processWrite(job)
				default:
					return
				}
			}
		case job := <-fm.pendingWrites:
			fm.processWrite(job)
		}
	}
}

// processWrite handles a single file write operation
func (fm *ConcurrentFileManager) processWrite(job FileWriteJob) {
	result := FileWriteResult{
		Path: job.Path,
	}

	// Check if file already exists (first-wins policy)
	fm.mu.RLock()
	existingAgent, exists := fm.createdFiles[job.Path]
	fm.mu.RUnlock()

	if exists {
		result.Success = false
		result.Error = fmt.Errorf("file %s already created by agent %s", job.Path, existingAgent)
		result.CreatedBy = existingAgent
		job.Response <- result
		return
	}

	// Create directory if needed
	dir := filepath.Dir(job.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		result.Success = false
		result.Error = fmt.Errorf("failed to create directory %s: %w", dir, err)
		job.Response <- result
		return
	}

	// Write the file
	if err := os.WriteFile(job.Path, []byte(job.Content), 0644); err != nil {
		result.Success = false
		result.Error = fmt.Errorf("failed to write file %s: %w", job.Path, err)
		job.Response <- result
		return
	}

	// Mark file as created
	fm.mu.Lock()
	fm.createdFiles[job.Path] = job.AgentRole
	fm.mu.Unlock()

	result.Success = true
	result.CreatedBy = job.AgentRole
	job.Response <- result
}

// WriteFile writes a file asynchronously and waits for the result
func (fm *ConcurrentFileManager) WriteFile(path, content string, agentRole agent.Role) error {
	responseChan := make(chan FileWriteResult, 1)

	job := FileWriteJob{
		Path:      path,
		Content:   content,
		AgentRole: agentRole,
		Response:  responseChan,
	}

	// Submit job
	select {
	case fm.pendingWrites <- job:
	case <-fm.done:
		return fmt.Errorf("file manager is shutting down")
	}

	// Wait for result
	result := <-responseChan
	if !result.Success {
		return result.Error
	}

	return nil
}

// IsFileCreated checks if a file has been created by any agent
func (fm *ConcurrentFileManager) IsFileCreated(path string) (bool, agent.Role) {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	agentRole, exists := fm.createdFiles[path]
	return exists, agentRole
}

// GetCreatedFiles returns a map of all created files and their creators
func (fm *ConcurrentFileManager) GetCreatedFiles() map[string]agent.Role {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	// Return a copy
	files := make(map[string]agent.Role, len(fm.createdFiles))
	for path, role := range fm.createdFiles {
		files[path] = role
	}
	return files
}

// MarkFileCreated manually marks a file as created (for existing files)
func (fm *ConcurrentFileManager) MarkFileCreated(path string, agentRole agent.Role) {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	fm.createdFiles[path] = agentRole
}

// Shutdown stops the file manager and waits for pending writes
func (fm *ConcurrentFileManager) Shutdown() {
	close(fm.done)
	fm.wg.Wait()
}

// Clear resets the created files tracking
func (fm *ConcurrentFileManager) Clear() {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	fm.createdFiles = make(map[string]agent.Role)
}
