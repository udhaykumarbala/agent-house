package agent

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// AgentDefinition holds the complete definition loaded from the filesystem.
type AgentDefinition struct {
	ID       string          `json:"id"`
	Name     string          `json:"name"`
	Color    string          `json:"color"`
	Icon     string          `json:"icon"`
	Prompt   string          `json:"-"` // loaded from agent.md
	Settings AgentSettings   `json:"-"` // loaded from settings.json
}

// AgentSettings maps to settings.json
type AgentSettings struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Color     string          `json:"color"`
	Icon      string          `json:"icon"`
	Execution ExecutionConfig `json:"execution"`
	Phases    []string        `json:"phases"`
	Tools     ToolsConfig     `json:"tools"`
}

// ExecutionConfig defines how the agent runs
type ExecutionConfig struct {
	DefaultMode     string            `json:"default_mode"`      // "oneshot" or "session"
	PhaseModes      map[string]string `json:"phase_modes"`       // phase → mode override
	Model           string            `json:"model"`             // model override (empty = default)
	MaxTurns        int               `json:"max_turns"`
	Timeout         string            `json:"timeout"`
	ClaudePermission string           `json:"claude_permission"` // bypassPermissions, acceptEdits, etc.
}

// ToolsConfig defines what tools the agent can access
type ToolsConfig struct {
	Builtin    []string `json:"builtin"`
	MCPServers []string `json:"mcp_servers"`
}

// Registry discovers and manages agent definitions from the filesystem.
type Registry struct {
	mu          sync.RWMutex
	agents      map[Role]*AgentDefinition
	agentsDir   string // path to agents/ directory
}

// NewRegistry creates a registry and loads agents from the given directory.
func NewRegistry(agentsDir string) *Registry {
	r := &Registry{
		agents:    make(map[Role]*AgentDefinition),
		agentsDir: agentsDir,
	}
	r.Load()
	return r
}

// Load discovers and loads all agent definitions from the filesystem.
func (r *Registry) Load() {
	r.mu.Lock()
	defer r.mu.Unlock()

	entries, err := os.ReadDir(r.agentsDir)
	if err != nil {
		log.Printf("[REGISTRY] Failed to read agents directory %s: %v", r.agentsDir, err)
		return
	}

	loaded := 0
	for _, entry := range entries {
		if !entry.IsDir() || entry.Name()[0] == '_' || entry.Name()[0] == '.' {
			continue
		}

		role := Role(entry.Name())
		agentDir := filepath.Join(r.agentsDir, entry.Name())

		def, err := loadAgentDefinition(agentDir, role)
		if err != nil {
			log.Printf("[REGISTRY] Failed to load agent %s: %v", role, err)
			continue
		}

		r.agents[role] = def
		loaded++
	}

	log.Printf("[REGISTRY] Loaded %d agents from %s", loaded, r.agentsDir)
}

// Get returns the definition for a role (nil if not found).
func (r *Registry) Get(role Role) *AgentDefinition {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.agents[role]
}

// List returns all loaded agent definitions.
func (r *Registry) List() []*AgentDefinition {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]*AgentDefinition, 0, len(r.agents))
	for _, def := range r.agents {
		list = append(list, def)
	}
	return list
}

// Roles returns all registered role IDs.
func (r *Registry) Roles() []Role {
	r.mu.RLock()
	defer r.mu.RUnlock()

	roles := make([]Role, 0, len(r.agents))
	for role := range r.agents {
		roles = append(roles, role)
	}
	return roles
}

// CreateAgent creates an Agent instance from a registry definition.
func (r *Registry) CreateAgent(role Role) (*Agent, error) {
	def := r.Get(role)
	if def == nil {
		return nil, fmt.Errorf("agent role %q not found in registry", role)
	}

	agent := &Agent{
		ID:           string(role),
		Name:         def.Name,
		Role:         role,
		SystemPrompt: def.Prompt,
		Color:        def.Color,
	}

	// Apply execution mode from settings
	if def.Settings.Execution.DefaultMode != "" {
		SetMode(role, ExecutionMode(def.Settings.Execution.DefaultMode))
	}

	// Apply phase-specific modes
	for phase, mode := range def.Settings.Execution.PhaseModes {
		// Store phase modes for lookup during execution
		SetPhaseMode(role, phase, ExecutionMode(mode))
	}

	return agent, nil
}

// loadAgentDefinition loads agent.md and settings.json from a directory.
func loadAgentDefinition(dir string, role Role) (*AgentDefinition, error) {
	def := &AgentDefinition{
		ID: string(role),
	}

	// Load settings.json
	settingsPath := filepath.Join(dir, "settings.json")
	settingsData, err := os.ReadFile(settingsPath)
	if err != nil {
		return nil, fmt.Errorf("settings.json: %w", err)
	}

	if err := json.Unmarshal(settingsData, &def.Settings); err != nil {
		return nil, fmt.Errorf("parse settings.json: %w", err)
	}

	def.Name = def.Settings.Name
	def.Color = def.Settings.Color
	def.Icon = def.Settings.Icon

	// Load agent.md (system prompt)
	promptPath := filepath.Join(dir, "agent.md")
	promptData, err := os.ReadFile(promptPath)
	if err != nil {
		return nil, fmt.Errorf("agent.md: %w", err)
	}
	def.Prompt = string(promptData)

	return def, nil
}
