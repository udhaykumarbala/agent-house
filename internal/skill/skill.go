package skill

// Skill represents a modular capability that can be loaded into an agent's prompt
type Skill struct {
	Name         string   // Unique skill identifier (e.g., "delegation", "code-review")
	Description  string   // What this skill does and when to use it
	Trigger      string   // Trigger condition (e.g., "phase:research AND role:pm")
	AllowedTools []string // Tools this skill allows (e.g., ["Read", "Write"])
	Content      string   // Markdown content with instructions
	Path         string   // File path where skill was loaded from
}

// Context provides the current execution context for skill matching
type Context struct {
	Role     string // Current agent role (e.g., "pm", "senior_dev")
	Phase    string // Current phase (e.g., "research", "development")
	Template string // Selected template (e.g., "static-enhanced", "go-api")
}

// SkillType categorizes skills by their scope
type SkillType string

const (
	SkillTypeShared   SkillType = "shared"   // Cross-agent skills
	SkillTypeAgent    SkillType = "agent"    // Agent-specific skills
	SkillTypeTemplate SkillType = "template" // Template-specific skills
)

// DefaultToolsByRole defines the default allowed tools for each agent role
var DefaultToolsByRole = map[string][]string{
	"ceo":        {"Read", "Write"},
	"pm":         {"Read", "Write"},
	"ux":         {"Read", "Write"},
	"ui":         {"Read", "Write"},
	"security":   {"Read", "Write", "Grep"},
	"architect":  {"Read", "Write", "Grep"},
	"senior_dev": {"Read", "Write", "Edit", "Bash"},
	"junior_dev": {"Read", "Write", "Edit", "Bash"},
}

// GetDefaultTools returns the default tools for a given role
func GetDefaultTools(role string) []string {
	if tools, ok := DefaultToolsByRole[role]; ok {
		return tools
	}
	// Default to read-only for unknown roles
	return []string{"Read"}
}

// MergeTools combines tool lists, removing duplicates
func MergeTools(toolSets ...[]string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, tools := range toolSets {
		for _, tool := range tools {
			if !seen[tool] {
				seen[tool] = true
				result = append(result, tool)
			}
		}
	}

	return result
}
