// Package scenario composes capabilities into named multi-step flows that
// model how the system *should* decompose a user request across agents. A
// scenario emits one TypeSystem message per step into the message store +
// WS hub so Conductor surfaces the multi-agent work live, then returns a
// structured Result with proactive Suggestions.
package scenario

import (
	"fmt"
	"sync"

	"pty-claude-test/internal/message"
)

// Step is one decomposed action attributed to a specific agent.
type Step struct {
	Agent  string         `json:"agent"`
	Action string         `json:"action"`
	Input  map[string]any `json:"input,omitempty"`
	Output map[string]any `json:"output,omitempty"`
	OK     bool           `json:"ok"`
	Note   string         `json:"note,omitempty"`
}

// Suggestion is a proactive follow-up the scenario produced as a side-
// effect. The Action key is the *type* of follow-up (e.g.
// "schedule_interviews"), so the UI / Conductor can dispatch it.
type Suggestion struct {
	Title   string         `json:"title"`
	Detail  string         `json:"detail"`
	Action  string         `json:"action"`
	Payload map[string]any `json:"payload,omitempty"`
}

// Result is the structured outcome of a scenario run.
type Result struct {
	Scenario    string       `json:"scenario"`
	Scope       string       `json:"scope"`
	Steps       []Step       `json:"steps"`
	Summary     string       `json:"summary"`
	Suggestions []Suggestion `json:"suggestions"`
	OK          bool         `json:"ok"`
}

// Emitter is what the engine gives each Runner so it can broadcast steps
// into the existing message infrastructure without importing internal/web.
type Emitter interface {
	Emit(m *message.Message)
}

// Runner is one scenario. Example() returns a sample input that any UI can
// pre-fill — this is how the frontend stays fully dynamic: it never names
// a specific scenario, it just renders whatever shape the backend ships.
type Runner interface {
	Name() string
	Description() string
	Example() map[string]any
	Run(ctx Context) (Result, error)
}

// Context is what a runner receives. DataRoot is where every capability
// persists; Scope is the per-tenant directory (e.g. "atlas-site").
type Context struct {
	DataRoot string
	Scope    string
	Input    map[string]any
	Emitter  Emitter
}

// Engine is the registry of scenarios.
type Engine struct {
	mu      sync.RWMutex
	runners map[string]Runner
}

func New() *Engine { return &Engine{runners: map[string]Runner{}} }

func (e *Engine) Register(r Runner) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.runners[r.Name()] = r
}

func (e *Engine) List() []map[string]any {
	e.mu.RLock()
	defer e.mu.RUnlock()
	out := make([]map[string]any, 0, len(e.runners))
	for _, r := range e.runners {
		out = append(out, map[string]any{
			"name":        r.Name(),
			"description": r.Description(),
			"example":     r.Example(),
		})
	}
	return out
}

func (e *Engine) Get(name string) (Runner, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	r, ok := e.runners[name]
	return r, ok
}

// HelpEmitStep is a small helper for runners — builds a TypeSystem message
// describing one step and emits it.
func HelpEmitStep(em Emitter, scope, scenario string, step Step) {
	if em == nil {
		return
	}
	content := fmt.Sprintf("[scenario:%s] %s · %s — %s",
		scenario, step.Agent, step.Action, step.Note)
	m := message.NewMessage(message.TypeSystem, step.Agent, "conductor", content)
	m.Metadata.ProjectID = scope
	m.Metadata.Tags = []string{"scenario", scenario, step.Agent}
	em.Emit(m)
}
