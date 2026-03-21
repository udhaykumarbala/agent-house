package agent

import (
	"log"
	"sync"
)

// ExecutionMode determines how an agent processes tasks.
type ExecutionMode string

const (
	ModeOneshot ExecutionMode = "oneshot" // Direct API call — fast, no tools, text in/text out
	ModeSession ExecutionMode = "session" // Claude Code session — persistent, full tool access
)

// DefaultModes defines the default execution mode per agent role.
// Oneshot for text-only agents, session for tool-using agents.
var DefaultModes = map[Role]ExecutionMode{
	RoleCEO:       ModeOneshot, // Triage, review — just text analysis
	RolePM:        ModeOneshot, // Research docs, specs — text output
	RoleUX:        ModeOneshot, // UX research/specs — text output
	RoleUI:        ModeSession, // May need MCP tools (design APIs)
	RoleSecurity:  ModeSession, // Needs to read code for audits
	RoleArchitect: ModeSession, // Needs to browse codebase
	RoleSeniorDev: ModeSession, // Full tool access for coding
	RoleJuniorDev: ModeSession, // Full tool access for coding
}

// ModeOverrides stores runtime mode overrides set by the user.
var (
	modeOverrides   = make(map[Role]ExecutionMode)
	modeOverridesMu sync.RWMutex
)

// GetMode returns the execution mode for a role, checking overrides first.
func GetMode(role Role) ExecutionMode {
	modeOverridesMu.RLock()
	if mode, ok := modeOverrides[role]; ok {
		modeOverridesMu.RUnlock()
		return mode
	}
	modeOverridesMu.RUnlock()

	if mode, ok := DefaultModes[role]; ok {
		return mode
	}
	return ModeSession // Default to session if unknown role
}

// SetMode overrides the execution mode for a role at runtime.
func SetMode(role Role, mode ExecutionMode) {
	if mode != ModeOneshot && mode != ModeSession {
		return
	}
	modeOverridesMu.Lock()
	modeOverrides[role] = mode
	modeOverridesMu.Unlock()
	log.Printf("[MODES] %s → %s", role, mode)
}

// ResetMode removes the override for a role, reverting to default.
func ResetMode(role Role) {
	modeOverridesMu.Lock()
	delete(modeOverrides, role)
	modeOverridesMu.Unlock()
}

// GetAllModes returns the current mode for all roles.
func GetAllModes() map[string]string {
	modes := make(map[string]string)
	for role := range DefaultModes {
		modes[string(role)] = string(GetMode(role))
	}
	return modes
}
