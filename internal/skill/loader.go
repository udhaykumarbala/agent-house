package skill

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

// Loader handles skill discovery, parsing, and caching
type Loader struct {
	projectDir string
	skillsDir  string
	cache      map[string]*Skill
	mu         sync.RWMutex
}

// NewLoader creates a new skill loader for the given project
func NewLoader(projectDir string) *Loader {
	return &Loader{
		projectDir: projectDir,
		skillsDir:  filepath.Join(projectDir, ".claude", "skills"),
		cache:      make(map[string]*Skill),
	}
}

// LoadAll loads all skills from the skills directory
func (l *Loader) LoadAll() ([]*Skill, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	var skills []*Skill

	// Check if skills directory exists
	if _, err := os.Stat(l.skillsDir); os.IsNotExist(err) {
		return skills, nil // No skills directory, return empty
	}

	// Walk through all subdirectories
	err := filepath.Walk(l.skillsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Look for SKILL.md files
		if info.Name() == "SKILL.md" {
			skill, err := l.parseSkillFile(path)
			if err != nil {
				// Log but continue
				fmt.Printf("Warning: Failed to parse skill at %s: %v\n", path, err)
				return nil
			}
			l.cache[skill.Name] = skill
			skills = append(skills, skill)
		}
		return nil
	})

	return skills, err
}

// LoadShared loads all shared skills (always loaded for all agents)
func (l *Loader) LoadShared() ([]*Skill, error) {
	sharedDir := filepath.Join(l.skillsDir, "shared")
	return l.loadFromDir(sharedDir)
}

// LoadForRole loads skills specific to an agent role
func (l *Loader) LoadForRole(role string, phase string) ([]*Skill, error) {
	agentDir := filepath.Join(l.skillsDir, "agents", role)
	skills, err := l.loadFromDir(agentDir)
	if err != nil {
		return nil, err
	}

	// Filter by trigger conditions
	var filtered []*Skill
	ctx := Context{Role: role, Phase: phase}
	for _, s := range skills {
		if l.matchesTrigger(s.Trigger, ctx) {
			filtered = append(filtered, s)
		}
	}

	return filtered, nil
}

// LoadForTemplate loads skills specific to a project template
func (l *Loader) LoadForTemplate(template string) ([]*Skill, error) {
	templateDir := filepath.Join(l.skillsDir, "templates", template)
	return l.loadFromDir(templateDir)
}

// LoadForContext loads all skills relevant to the given context
func (l *Loader) LoadForContext(ctx Context) ([]*Skill, error) {
	var allSkills []*Skill

	// 1. Load shared skills (always)
	shared, err := l.LoadShared()
	if err == nil {
		allSkills = append(allSkills, shared...)
	}

	// 2. Load agent-specific skills
	agentSkills, err := l.LoadForRole(ctx.Role, ctx.Phase)
	if err == nil {
		allSkills = append(allSkills, agentSkills...)
	}

	// 3. Load template skills (only in development phase)
	if ctx.Phase == "development" && ctx.Template != "" {
		templateSkills, err := l.LoadForTemplate(ctx.Template)
		if err == nil {
			allSkills = append(allSkills, templateSkills...)
		}
	}

	return allSkills, nil
}

// loadFromDir loads all skills from a directory
func (l *Loader) loadFromDir(dir string) ([]*Skill, error) {
	var skills []*Skill

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return skills, nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			skillPath := filepath.Join(dir, entry.Name(), "SKILL.md")
			if _, err := os.Stat(skillPath); err == nil {
				skill, err := l.parseSkillFile(skillPath)
				if err != nil {
					continue
				}
				skills = append(skills, skill)
			}
		}
	}

	return skills, nil
}

// parseSkillFile parses a SKILL.md file into a Skill struct
func (l *Loader) parseSkillFile(path string) (*Skill, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return ParseSkillContent(string(content), path)
}

// ParseSkillContent parses skill content with YAML frontmatter
func ParseSkillContent(content string, path string) (*Skill, error) {
	skill := &Skill{
		Path: path,
	}

	// Check for YAML frontmatter (between --- markers)
	frontmatterRegex := regexp.MustCompile(`(?s)^---\n(.+?)\n---\n(.*)$`)
	matches := frontmatterRegex.FindStringSubmatch(content)

	if len(matches) == 3 {
		// Parse frontmatter
		frontmatter := matches[1]
		skill.Content = strings.TrimSpace(matches[2])

		// Parse YAML-like frontmatter manually (simple key: value parsing)
		lines := strings.Split(frontmatter, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			parts := strings.SplitN(line, ":", 2)
			if len(parts) != 2 {
				continue
			}

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			switch key {
			case "name":
				skill.Name = value
			case "description":
				skill.Description = value
			case "trigger":
				skill.Trigger = value
			case "allowed-tools":
				// Parse comma-separated tools
				tools := strings.Split(value, ",")
				for _, tool := range tools {
					tool = strings.TrimSpace(tool)
					if tool != "" {
						skill.AllowedTools = append(skill.AllowedTools, tool)
					}
				}
			}
		}
	} else {
		// No frontmatter, treat entire content as the skill content
		skill.Content = content
		// Extract name from directory
		dir := filepath.Dir(path)
		skill.Name = filepath.Base(dir)
	}

	if skill.Name == "" {
		return nil, fmt.Errorf("skill has no name: %s", path)
	}

	return skill, nil
}

// matchesTrigger checks if the context matches the trigger condition
func (l *Loader) matchesTrigger(trigger string, ctx Context) bool {
	if trigger == "" {
		return true // No trigger means always match
	}

	return MatchTrigger(trigger, ctx)
}

// GetSkill returns a cached skill by name
func (l *Loader) GetSkill(name string) *Skill {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.cache[name]
}

// ClearCache clears the skill cache
func (l *Loader) ClearCache() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.cache = make(map[string]*Skill)
}

// ComposePrompt combines a base prompt with loaded skills
func ComposePrompt(basePrompt string, skills []*Skill) string {
	if len(skills) == 0 {
		return basePrompt
	}

	var builder strings.Builder
	builder.WriteString(basePrompt)
	builder.WriteString("\n\n---\n\n# Active Skills\n\n")

	for _, s := range skills {
		builder.WriteString("## ")
		builder.WriteString(s.Name)
		builder.WriteString("\n\n")
		if s.Description != "" {
			builder.WriteString("*")
			builder.WriteString(s.Description)
			builder.WriteString("*\n\n")
		}
		builder.WriteString(s.Content)
		builder.WriteString("\n\n---\n\n")
	}

	return builder.String()
}

// AggregateTools collects all allowed tools from skills and merges with defaults
func AggregateTools(role string, skills []*Skill) []string {
	// Start with default tools for the role
	tools := GetDefaultTools(role)

	// Collect tools from skills
	for _, s := range skills {
		if len(s.AllowedTools) > 0 {
			tools = MergeTools(tools, s.AllowedTools)
		}
	}

	return tools
}
