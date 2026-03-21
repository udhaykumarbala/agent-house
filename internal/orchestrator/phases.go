package orchestrator

import (
	"pty-claude-test/internal/agent"
	"pty-claude-test/internal/template"
)

// Phase represents a stage in the multi-phase workflow
type Phase string

const (
	// PhaseTriage - CEO reviews the task and project context, decides which phases to run
	PhaseTriage Phase = "triage"
	// PhaseTemplateSelection - architect selects template, CEO/PM approve (NEW)
	PhaseTemplateSelection Phase = "template_selection"
	// PhaseResearch - agents research competitors, patterns, design inspiration
	PhaseResearch Phase = "research"
	// PhasePlanning - agents create specs based on research
	PhasePlanning Phase = "planning"
	// PhaseDiscussion - CEO moderates, agents iterate until approval
	PhaseDiscussion Phase = "discussion"
	// PhaseDevelopment - developers write code based on approved specs
	PhaseDevelopment Phase = "development"
	// PhaseQA - CEO reviews development phase for quality
	PhaseQA Phase = "qa"
	// PhaseDevIteration - development revision based on QA feedback
	PhaseDevIteration Phase = "dev_iteration"
)

// PhaseConfig defines which agents participate in each phase
type PhaseConfig struct {
	// Phase 0: Template Selection
	TemplateProposer  agent.Role   // Architect proposes
	TemplateReviewers []agent.Role // CEO and PM review
	// Phase 1-4: Existing phases
	ResearchAgents    []agent.Role
	PlanningAgents    []agent.Role
	DiscussionLeader  agent.Role
	DevelopmentAgents []agent.Role
}

// DefaultPhaseConfig returns the standard phase configuration
func DefaultPhaseConfig() PhaseConfig {
	return PhaseConfig{
		// Phase 0: Template Selection - architect proposes, CEO/PM review
		TemplateProposer:  agent.RoleArchitect,
		TemplateReviewers: []agent.Role{agent.RoleCEO, agent.RolePM},
		// Phase 1: Research - these agents create research .md files
		ResearchAgents: []agent.Role{
			agent.RolePM,        // Product research, competitors
			agent.RoleUX,        // UX pattern research
			agent.RoleUI,        // Design inspiration research
			agent.RoleSecurity,  // Security considerations (if applicable)
			agent.RoleArchitect, // Technical research (if applicable)
		},
		// Phase 2: Planning - same agents create spec .md files
		PlanningAgents: []agent.Role{
			agent.RolePM,
			agent.RoleUX,
			agent.RoleUI,
			agent.RoleSecurity,
			agent.RoleArchitect,
		},
		// Phase 3: Discussion - CEO moderates
		DiscussionLeader: agent.RoleCEO,
		// Phase 4: Development - developers implement
		DevelopmentAgents: []agent.Role{
			agent.RoleSeniorDev,
			agent.RoleJuniorDev,
		},
	}
}

// PhasePromptContext builds the context for an agent based on current phase
func PhasePromptContext(phase Phase, originalTask string) string {
	switch phase {
	case PhaseTemplateSelection:
		return buildTemplateSelectionContext(originalTask)
	case PhaseResearch:
		return buildResearchContext(originalTask)
	case PhasePlanning:
		return buildPlanningContext(originalTask)
	case PhaseDiscussion:
		return buildDiscussionContext(originalTask)
	case PhaseDevelopment:
		return buildSimpleDevelopmentContext(originalTask)
	default:
		return originalTask
	}
}

// PhasePromptContextWithTemplate builds the context with template guidelines
func PhasePromptContextWithTemplate(phase Phase, originalTask string, templateType template.TemplateType) string {
	baseContext := PhasePromptContext(phase, originalTask)

	if phase == PhaseDevelopment {
		// Inject template guidelines into development context
		guidelines := template.GetGuidelines(templateType)
		return baseContext + "\n\n### Template Guidelines\n" + guidelines
	}

	return baseContext
}

func buildTemplateSelectionContext(task string) string {
	return `## PHASE 0: TEMPLATE SELECTION

You are in the TEMPLATE SELECTION phase. Before any research begins, we must choose the right project template.

### Original Task
` + task + `

### Available Templates

| Template | Best For | Complexity |
|----------|----------|------------|
| static-html | Landing pages, portfolios, simple sites | Simple |
| static-enhanced | Marketing sites, blogs with build tools (DEFAULT) | Basic |
| nextjs-frontend | SPAs, SSR/SSG apps, dashboards | Medium |
| go-api | REST APIs, microservices, backends | Medium |
| fullstack | Complete apps with Next.js + Go + Docker | Complex |

### For Architect
Analyze the task and propose a template. Create .plans/template.md with:
1. Chosen template (default to static-enhanced if unclear)
2. Justification (3-5 bullet points)
3. Any customizations needed

### Selection Criteria
- Does it need a backend? -> go-api or fullstack
- Is it a complex frontend? -> nextjs-frontend
- Is it simple/static? -> static-html or static-enhanced
- Does it need build tools? -> static-enhanced over static-html

### For CEO/PM (Reviewers)
Review .plans/template.md and either:
- TEMPLATE_APPROVED: [template-name]
- TEMPLATE_CHANGE_REQUESTED: [reason and suggested template]

When template is approved, signal: PHASE_COMPLETE: template_selection
`
}

func buildResearchContext(task string) string {
	return `## PHASE 1: RESEARCH

You are in the RESEARCH phase. Your job is to research and document your findings.

### Original Task
` + task + `

### Your Assignment
Create a research document (.md file) in .plans/research/ folder.

DO NOT:
- Create any code files yet
- Skip research
- Make assumptions without research

DO:
- Research competitors, patterns, best practices
- Document your findings in a structured .md file
- Be specific with examples and references

When done, signal: PHASE_COMPLETE: research
`
}

func buildPlanningContext(task string) string {
	return `## PHASE 2: PLANNING

You are in the PLANNING phase. Research has been completed.

### Original Task
` + task + `

### Available Research
Check .plans/research/ folder for research documents created by other agents.

### Your Assignment
Based on the research, create your specification document in .plans/specs/ folder.

DO NOT:
- Create any code files yet
- Ignore the research documents
- Repeat research work

DO:
- Read the research files first
- Create a detailed spec .md file
- Reference decisions from research

When done, signal: PHASE_COMPLETE: planning
`
}

func buildDiscussionContext(task string) string {
	return `## PHASE 3: DISCUSSION

You are in the DISCUSSION phase. All specs have been created.

### Original Task
` + task + `

### Available Specs
Check .plans/specs/ folder for all specification documents.

### Your Assignment (CEO)
1. Review all spec documents
2. Identify conflicts or gaps
3. Ask agents to update their specs if needed
4. When satisfied, approve and create .plans/final/approved-plan.md

### For Other Agents
If CEO requests changes, update your spec file accordingly.

When CEO approves, signal: PHASE_COMPLETE: discussion
`
}

func buildSimpleDevelopmentContext(task string) string {
	return `## PHASE 4: DEVELOPMENT

You are in the DEVELOPMENT phase. The plan has been approved.

### Original Task
` + task + `

### IMPORTANT: Read These First
Before writing ANY code, read these files:
1. .plans/specs/product-spec.md - Features and priorities
2. .plans/specs/ux-spec.md - Wireframes and flows
3. .plans/specs/ui-spec.md - Colors, fonts, spacing
4. .plans/final/approved-plan.md - CEO's approved plan

### Your Assignment
Implement the application following the approved specs EXACTLY.

DO NOT:
- Deviate from the UI spec colors/fonts
- Skip reading the spec files
- Use generic AI colors (blue #3B82F6, green #10B981)

DO:
- Use the EXACT colors from ui-spec.md
- Follow the wireframes from ux-spec.md
- Implement features in the order from approved-plan.md

When done, signal: COMPLETE: Development finished following approved specs.
`
}

// GetPhaseName returns a human-readable name for the phase
func GetPhaseName(phase Phase) string {
	names := map[Phase]string{
		PhaseTriage:            "CEO Triage",
		PhaseTemplateSelection: "Template Selection",
		PhaseResearch:          "Research",
		PhasePlanning:          "Planning",
		PhaseDiscussion:        "Discussion",
		PhaseDevelopment:       "Development",
		PhaseQA:                "QA Review",
		PhaseDevIteration:      "Development Iteration",
	}
	if name, ok := names[phase]; ok {
		return name
	}
	return string(phase)
}

// GetPhaseEmoji returns an emoji for the phase
func GetPhaseEmoji(phase Phase) string {
	emojis := map[Phase]string{
		PhaseTriage:            "👔",
		PhaseTemplateSelection: "🎯",
		PhaseResearch:          "🔍",
		PhasePlanning:          "📋",
		PhaseDiscussion:        "💬",
		PhaseDevelopment:       "💻",
		PhaseQA:                "✅",
		PhaseDevIteration:      "🔄",
	}
	if emoji, ok := emojis[phase]; ok {
		return emoji
	}
	return "📌"
}

// AllPhases returns all phases in order
func AllPhases() []Phase {
	return []Phase{
		PhaseTemplateSelection,
		PhaseResearch,
		PhasePlanning,
		PhaseDiscussion,
		PhaseDevelopment,
	}
}
