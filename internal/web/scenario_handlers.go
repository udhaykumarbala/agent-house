package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
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
	s.scenarioEngine.Register(scenario.WorkforceSnapshot{})
	s.scenarioEngine.Register(scenario.StaffProject{})
	s.scenarioEngine.Register(scenario.EstimateProject{})

	// Wire the Brain's run_scenario action to the engine, so a typed prompt like
	// "process my inbox" fans work out to specialist agents (steps stream to the
	// workforce panel via the same emitter used by the HTTP scenario route).
	if s.brainHandler != nil {
		s.brainHandler.executor.OnRunScenario = s.runScenarioForBrain
	}
}

// runScenarioForBrain runs a capability scenario for the Brain's run_scenario
// action and returns a human summary + suggestion titles for the chat lane.
func (s *Server) runScenarioForBrain(name, scope string, input map[string]string) (string, []string, error) {
	if s.scenarioEngine == nil {
		return "", nil, fmt.Errorf("scenario engine not ready")
	}
	runner, ok := s.scenarioEngine.Get(name)
	if !ok {
		return "", nil, fmt.Errorf("unknown scenario %q", name)
	}
	if scope == "" {
		scope = "atlas-site"
	}
	in := make(map[string]any, len(input))
	for k, v := range input {
		in[k] = v
	}
	res, err := runner.Run(scenario.Context{
		DataRoot: dataRoot,
		Scope:    scope,
		Input:    in,
		Emitter:  &scenarioEmitter{store: s.store, hub: s.hub},
	})
	if err != nil {
		return "", nil, err
	}
	var suggs []string
	for _, sg := range res.Suggestions {
		if sg.Title != "" {
			suggs = append(suggs, sg.Title)
		}
	}
	return res.Summary, suggs, nil
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

// handleReseed restores demo state for a scope by re-running its (idempotent)
// seed script — deterministic ids overwrite, so it's safe to call repeatedly on
// stage. This is a demo/dev convenience: it restores the inbox (incl. the BEC
// email, which is deletable) and the in-memory cron fires (lost on restart).
// POST /api/scenario/reseed?scope=atlas-site
func (s *Server) handleReseed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}
	scope := scopeOf(r)
	// Only the atlas-site demo scope ships a seed script today.
	const script = "scripts/scenarios/load-epc-atlas.sh"
	if _, err := os.Stat(script); err != nil {
		writeJSON(w, map[string]any{"success": false, "error": "no seed script available for scope " + scope})
		return
	}
	cmd := exec.Command("bash", script)
	cmd.Env = append(os.Environ(), "HOST=http://"+r.Host) // seed against this same server
	out, err := cmd.CombinedOutput()
	if err != nil {
		writeJSON(w, map[string]any{"success": false, "scope": scope, "error": err.Error(), "output": string(out)})
		return
	}
	writeJSON(w, map[string]any{"success": true, "scope": scope, "output": string(out)})
}
