package brain

import (
	"encoding/json"
	"strconv"
)

// BrainDecision is the structured output from the Brain router.
type BrainDecision struct {
	Action      string    `json:"action"`
	Params      paramMap  `json:"params,omitempty"`
	Response    string    `json:"response"`
	Suggestions []string  `json:"suggestions,omitempty"` // Quick action buttons for the user
}

// paramMap is a string-valued param map that tolerantly decodes scalar JSON
// values (numbers, bools) by stringifying them. Models routinely emit e.g.
// "amount": 340000 (a number); without this the whole decision fails to parse
// and silently degrades to a plain "respond".
type paramMap map[string]string

func (m *paramMap) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	out := make(map[string]string, len(raw))
	for k, v := range raw {
		switch x := v.(type) {
		case string:
			out[k] = x
		case float64:
			if x == float64(int64(x)) {
				out[k] = strconv.FormatInt(int64(x), 10)
			} else {
				out[k] = strconv.FormatFloat(x, 'f', -1, 64)
			}
		case bool:
			out[k] = strconv.FormatBool(x)
		case nil:
			out[k] = ""
		default:
			if bb, err := json.Marshal(x); err == nil {
				out[k] = string(bb)
			}
		}
	}
	*m = out
	return nil
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
	ActionDeleteEmail        = "delete_email"        // Delete/archive an email
	ActionSendReply          = "send_reply"          // Send email reply
	ActionShortlistApplicant = "shortlist_applicant" // Mark applicant as shortlisted
	ActionArchiveEmail       = "archive_email"       // Archive (not delete) an email
	ActionRunScenario        = "run_scenario"        // Run a multi-agent capability scenario (fans out to specialists)
)

// ChatMessage represents a message in the conversation.
type ChatMessage struct {
	Role      string `json:"role"`    // "user" or "assistant"
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}
