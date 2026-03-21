package session

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// MCPConfig is the project-level MCP configuration.
// Stored at {projectDir}/.agent-house/mcp.json
type MCPConfig struct {
	MCPServers map[string]MCPServerDef `json:"mcpServers"`
}

// MCPServerDef defines an MCP server with role-based access.
type MCPServerDef struct {
	Command string            `json:"command"`           // e.g. "npx"
	Args    []string          `json:"args"`              // e.g. ["-y", "@modelcontextprotocol/server-postgres"]
	Env     map[string]string `json:"env,omitempty"`     // environment variables
	Roles   []string          `json:"roles"`             // which agent roles can use this: ["*"] = all
}

// LoadMCPConfig reads the project MCP config from .agent-house/mcp.json.
// Returns nil if the file doesn't exist (no MCP tools for this project).
func LoadMCPConfig(projectDir string) *MCPConfig {
	configPath := filepath.Join(projectDir, ".agent-house", "mcp.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil
	}

	var config MCPConfig
	if err := json.Unmarshal(data, &config); err != nil {
		log.Printf("[MCP] Failed to parse %s: %v", configPath, err)
		return nil
	}

	log.Printf("[MCP] Loaded config with %d servers from %s", len(config.MCPServers), configPath)
	return &config
}

// BuildMCPConfigForRole filters MCP servers to only those accessible by the given role,
// expands environment variables, writes a temp config file, and returns its path.
// Returns empty string if no servers match.
func BuildMCPConfigForRole(config *MCPConfig, agentRole string) string {
	if config == nil || len(config.MCPServers) == 0 {
		return ""
	}

	// Filter servers for this role
	filtered := make(map[string]interface{})
	for name, server := range config.MCPServers {
		if !roleAllowed(server.Roles, agentRole) {
			continue
		}

		// Expand env vars in server env
		expandedEnv := make(map[string]string)
		for k, v := range server.Env {
			expandedEnv[k] = os.ExpandEnv(v)
		}

		// Expand env vars in args
		expandedArgs := make([]string, len(server.Args))
		for i, arg := range server.Args {
			expandedArgs[i] = os.ExpandEnv(arg)
		}

		entry := map[string]interface{}{
			"command": server.Command,
			"args":    expandedArgs,
		}
		if len(expandedEnv) > 0 {
			entry["env"] = expandedEnv
		}
		filtered[name] = entry
	}

	if len(filtered) == 0 {
		return ""
	}

	// Write temp config file
	tempConfig := map[string]interface{}{
		"mcpServers": filtered,
	}
	data, err := json.MarshalIndent(tempConfig, "", "  ")
	if err != nil {
		log.Printf("[MCP] Failed to marshal config for role %s: %v", agentRole, err)
		return ""
	}

	tmpFile, err := os.CreateTemp("", fmt.Sprintf("mcp-%s-*.json", agentRole))
	if err != nil {
		log.Printf("[MCP] Failed to create temp file for role %s: %v", agentRole, err)
		return ""
	}
	defer tmpFile.Close()

	if _, err := tmpFile.Write(data); err != nil {
		log.Printf("[MCP] Failed to write temp config for role %s: %v", agentRole, err)
		os.Remove(tmpFile.Name())
		return ""
	}

	log.Printf("[MCP] Built config for %s: %d servers → %s", agentRole, len(filtered), tmpFile.Name())
	return tmpFile.Name()
}

func roleAllowed(roles []string, agentRole string) bool {
	for _, r := range roles {
		if r == "*" || r == agentRole {
			return true
		}
	}
	return false
}
