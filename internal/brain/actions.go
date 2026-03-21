package brain

// BrainDecision is the structured output from the Brain router.
type BrainDecision struct {
	Action   string            `json:"action"`
	Params   map[string]string `json:"params,omitempty"`
	Response string            `json:"response"`
}

// Action types
const (
	ActionRespond       = "respond"        // Direct text response
	ActionCreateProject = "create_project" // Create project + start pipeline
	ActionDelegate      = "delegate"       // Inject task to specific agent
	ActionProjectStatus = "project_status" // Get detailed project status
	ActionListProjects  = "list_projects"  // List all projects
	ActionReadFile      = "read_file"      // Server reads file, re-calls Brain
	ActionEscalate      = "escalate"       // Spawn session for deep analysis
	ActionSetReminder   = "set_reminder"   // Create cron job
)

// ChatMessage represents a message in the conversation.
type ChatMessage struct {
	Role      string `json:"role"`    // "user" or "assistant"
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}
