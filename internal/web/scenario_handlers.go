package web

import (
	"encoding/json"
	"net/http"
	"strings"

	"pty-claude-test/internal/message"
	"pty-claude-test/internal/scenario"
)

// scenarioEmitter bridges the scenario engine (which knows nothing about
// web infra) to the message store + WS hub the rest of the system uses.
type scenarioEmitter struct {
	store *message.Store
	hub   *Hub
}

func (e *scenarioEmitter) Emit(m *message.Message) {
	if e == nil || m == nil {
		return
	}
	if e.store != nil {
		e.store.Add(m)
	}
	if e.hub != nil {
		e.hub.Broadcast(m)
	}
}

// registerScenarios installs the three demo runners. Called once at startup.
func (s *Server) registerScenarios() {
	if s.scenarioEngine == nil {
		s.scenarioEngine = scenario.New()
	}
	s.scenarioEngine.Register(scenario.ProcessApplicants{})
	s.scenarioEngine.Register(scenario.ValidateInvoice{})
	s.scenarioEngine.Register(scenario.ScheduleCheck{})
	s.scenarioEngine.Register(scenario.ProcessInbox{})
	s.scenarioEngine.Register(scenario.MorningBriefing{})
	s.scenarioEngine.Register(scenario.RouteRFI{})
}

func (s *Server) handleScenarioList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", 405)
		return
	}
	if s.scenarioEngine == nil {
		writeJSON(w, map[string]any{"scenarios": []any{}})
		return
	}
	writeJSON(w, map[string]any{"scenarios": s.scenarioEngine.List()})
}

// handleScenarioRun handles POST /api/scenario/{name}/run
func (s *Server) handleScenarioRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/api/scenario/")
	name := strings.TrimSuffix(path, "/run")
	if name == "" || strings.Contains(name, "/") {
		http.Error(w, "scenario name required: POST /api/scenario/{name}/run", 400)
		return
	}
	runner, ok := s.scenarioEngine.Get(name)
	if !ok {
		http.Error(w, "unknown scenario: "+name, 404)
		return
	}

	var body map[string]any
	if r.ContentLength > 0 {
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "bad json: "+err.Error(), 400)
			return
		}
	}
	if body == nil {
		body = map[string]any{}
	}

	ctx := scenario.Context{
		DataRoot: dataRoot,
		Scope:    scopeOf(r),
		Input:    body,
		Emitter:  &scenarioEmitter{store: s.store, hub: s.hub},
	}
	result, err := runner.Run(ctx)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	writeJSON(w, result)
}
