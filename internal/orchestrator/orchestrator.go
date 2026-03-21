package orchestrator

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"pty-claude-test/internal/agent"
	"pty-claude-test/internal/checkpoint"
	"pty-claude-test/internal/kanban"
	"pty-claude-test/internal/message"
	"pty-claude-test/internal/session"
	"pty-claude-test/internal/task"
	"pty-claude-test/internal/template"
)

// Config holds orchestrator configuration
type Config struct {
	ProjectDir   string
	MaxDepth     int  // Maximum delegation depth
	MaxTurns     int  // Maximum total turns
	EnableFiles  bool // Whether to execute file operations
	Verbose      bool // Print detailed logs
	EnablePhases bool // Enable multi-phase workflow (research → plan → discuss → develop)
}

// DefaultConfig returns sensible defaults
func DefaultConfig() Config {
	return Config{
		ProjectDir:   "./projects/default",
		MaxDepth:     3,  // Reduced from 5 to prevent deep chains
		MaxTurns:     8,  // Reduced from 20 to limit runaway tasks
		EnableFiles:  true,
		Verbose:      true,
		EnablePhases: true, // Enable multi-phase by default
	}
}

// CheckpointNotifyFunc is called when a checkpoint is reached or resolved
type CheckpointNotifyFunc func(cpID, cpType, summary string, phaseIndex int, resolved bool, decidedBy string, action string)

// Orchestrator manages agent collaboration
type Orchestrator struct {
	config           Config
	store            *message.Store
	agents           map[agent.Role]*agent.Agent
	mu               sync.RWMutex
	turns            int
	listeners        []func(msg *message.Message)
	historyManager   *task.HistoryManager
	currentPhase     Phase
	phaseConfig      PhaseConfig
	selectedTemplate template.TemplateType // Template selected for current task
	fileManager      *ConcurrentFileManager // Thread-safe file operations
	resultMu         sync.Mutex // Protects Result struct during parallel execution

	// Session management for Claude Code sessions per agent
	sessionManager   *session.SessionManager

	// Checkpoint system
	checkpointCh     chan *checkpoint.Decision // Receives human/CEO decisions
	pendingCP        *checkpoint.Checkpoint    // Currently pending checkpoint
	cpMu             sync.Mutex               // Protects pendingCP
	checkpointNotify CheckpointNotifyFunc      // Callback for checkpoint events
}

// New creates a new orchestrator
func New(config Config, store *message.Store) *Orchestrator {
	return &Orchestrator{
		config:           config,
		store:            store,
		agents:           make(map[agent.Role]*agent.Agent),
		historyManager:   task.NewHistoryManager(),
		phaseConfig:      DefaultPhaseConfig(),
		selectedTemplate: template.DefaultTemplate,
		fileManager:      NewConcurrentFileManager(),
		checkpointCh:     make(chan *checkpoint.Decision, 1),
	}
}

// NewWithSessions creates a new orchestrator with session-based agent execution.
func NewWithSessions(config Config, store *message.Store, sm *session.SessionManager) *Orchestrator {
	o := New(config, store)
	o.sessionManager = sm
	return o
}

// GetSessionManager returns the session manager (may be nil in legacy mode).
func (o *Orchestrator) GetSessionManager() *session.SessionManager {
	return o.sessionManager
}

// SetCheckpointNotify sets the callback for checkpoint events
func (o *Orchestrator) SetCheckpointNotify(fn CheckpointNotifyFunc) {
	o.checkpointNotify = fn
}

// GetPendingCheckpoint returns the currently pending checkpoint, if any
func (o *Orchestrator) GetPendingCheckpoint() *checkpoint.Checkpoint {
	o.cpMu.Lock()
	defer o.cpMu.Unlock()
	return o.pendingCP
}

// ResolveCheckpoint unblocks the orchestrator by sending a decision
func (o *Orchestrator) ResolveCheckpoint(dec *checkpoint.Decision) error {
	o.cpMu.Lock()
	pending := o.pendingCP
	o.cpMu.Unlock()

	if pending == nil {
		return fmt.Errorf("no pending checkpoint")
	}

	// Resolve on disk
	_, err := checkpoint.ResolveCheckpoint(o.config.ProjectDir, pending.ID, *dec)
	if err != nil {
		return fmt.Errorf("failed to resolve checkpoint: %w", err)
	}

	// Record the decision
	checkpoint.AddDecision(o.config.ProjectDir, *dec)

	// Unblock the orchestrator
	select {
	case o.checkpointCh <- dec:
	default:
		// Channel already has a value (auto-delegate beat us)
	}

	return nil
}

// waitForCheckpoint blocks execution until a human or CEO agent resolves the checkpoint.
// Returns nil decision if the checkpoint type is disabled in settings.
func (o *Orchestrator) waitForCheckpoint(cpType string, phaseIndex int, artifactPath, artifactSummary, projectID string, result *Result) (*checkpoint.Decision, error) {
	// Check if this checkpoint type is enabled
	settings, err := checkpoint.LoadSettings(o.config.ProjectDir)
	if err != nil {
		settings = &checkpoint.WorkflowSettings{}
		*settings = checkpoint.DefaultSettings(projectID)
	}
	if !settings.IsCheckpointEnabled(cpType) {
		return nil, nil // Skip — checkpoint disabled
	}

	// Create checkpoint record
	cp := checkpoint.Checkpoint{
		ID:              checkpoint.GenerateCheckpointID(),
		ProjectID:       projectID,
		TaskID:          result.TaskID,
		Type:            cpType,
		PhaseIndex:      phaseIndex,
		Status:          "pending",
		ArtifactPath:    artifactPath,
		ArtifactSummary: artifactSummary,
		CreatedAt:       time.Now(),
	}

	if err := checkpoint.CreateCheckpoint(o.config.ProjectDir, cp); err != nil {
		log.Printf("[CHECKPOINT] Failed to create checkpoint: %v", err)
		return nil, nil // Don't block on storage failure
	}

	// Set as pending
	o.cpMu.Lock()
	o.pendingCP = &cp
	o.cpMu.Unlock()

	// Notify listeners
	if o.checkpointNotify != nil {
		o.checkpointNotify(cp.ID, cpType, artifactSummary, phaseIndex, false, "", "")
	}

	// Broadcast system message
	sysMsg := message.NewMessage(message.TypeSystem, "system", "", fmt.Sprintf("Checkpoint: Waiting for approval — %s", artifactSummary))
	sysMsg.Metadata.ProjectID = projectID
	sysMsg.Metadata.TaskID = result.TaskID
	o.store.Add(sysMsg)
	o.notify(sysMsg)

	log.Printf("[CHECKPOINT] Waiting for resolution: type=%s id=%s", cpType, cp.ID)

	// Start auto-delegation timer if configured
	var cancelTimer context.CancelFunc
	if settings.AutoDelegateMinutes > 0 {
		var ctx context.Context
		ctx, cancelTimer = context.WithCancel(context.Background())
		go o.autoDelegateTimer(ctx, &cp, settings.AutoDelegateMinutes, projectID, result)
	}

	// Block until decision arrives
	dec := <-o.checkpointCh

	// Cancel timer if it was running
	if cancelTimer != nil {
		cancelTimer()
	}

	// Clear pending
	o.cpMu.Lock()
	o.pendingCP = nil
	o.cpMu.Unlock()

	// Notify resolved
	if o.checkpointNotify != nil {
		o.checkpointNotify(cp.ID, cpType, artifactSummary, phaseIndex, true, dec.DecidedBy, string(dec.Action))
	}

	log.Printf("[CHECKPOINT] Resolved: type=%s action=%s by=%s", cpType, dec.Action, dec.DecidedBy)

	return dec, nil
}

// autoDelegateTimer waits for the configured duration then delegates to CEO agent
func (o *Orchestrator) autoDelegateTimer(ctx context.Context, cp *checkpoint.Checkpoint, minutes int, projectID string, result *Result) {
	select {
	case <-time.After(time.Duration(minutes) * time.Minute):
		// Timer expired — delegate to CEO
	case <-ctx.Done():
		// Human resolved before timer expired
		return
	}

	log.Printf("[CHECKPOINT] Auto-delegating to CEO after %d minutes: %s", minutes, cp.ID)

	// Read artifact content
	artifactContent := ""
	if cp.ArtifactPath != "" {
		if data, err := os.ReadFile(cp.ArtifactPath); err == nil {
			artifactContent = string(data)
			if len(artifactContent) > 8000 {
				artifactContent = artifactContent[:8000] + "\n...(truncated)"
			}
		}
	}

	// Build CEO review prompt
	prompt := fmt.Sprintf(`You are acting as an automated reviewer because the human has not responded within %d minutes.

CHECKPOINT TYPE: %s
ARTIFACT SUMMARY: %s

ARTIFACT TO REVIEW:
%s

Review the above artifact and respond with EXACTLY one of these on its own line:
- CHECKPOINT_APPROVED: {brief justification}
- CHECKPOINT_REJECTED: {specific feedback for improvement}
- CHECKPOINT_ESCALATED: {reason this needs mandatory human review}`, minutes, cp.Type, cp.ArtifactSummary, artifactContent)

	// Spawn CEO agent
	ceoAgent, err := o.getAgent(agent.RoleCEO)
	if err != nil {
		log.Printf("[CHECKPOINT] Failed to get CEO agent: %v", err)
		return
	}

	var response *agent.Response
	if ceoAgent.SessionManager != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		sessResp, sessErr := ceoAgent.ExecuteWithSession(ctx, projectID, o.config.ProjectDir, prompt)
		cancel()
		if sessErr != nil {
			log.Printf("[CHECKPOINT] CEO auto-review (session) failed: %v", sessErr)
			return
		}
		response = &agent.Response{Content: sessResp.Text}
	} else {
		var respErr error
		response, _, respErr = ceoAgent.ProcessWithFileOps(prompt, o.config.ProjectDir)
		if respErr != nil {
			log.Printf("[CHECKPOINT] CEO auto-review failed: %v", respErr)
			return
		}
	}

	// Parse CEO decision
	dec := &checkpoint.Decision{
		ID:             checkpoint.GenerateDecisionID(),
		ProjectID:      projectID,
		TaskID:         result.TaskID,
		CheckpointType: cp.Type,
		PhaseIndex:     cp.PhaseIndex,
		DecidedBy:      "ceo_auto",
		DecidedAt:      time.Now(),
	}

	content := response.Content
	if strings.Contains(content, "CHECKPOINT_APPROVED") {
		dec.Action = checkpoint.DecisionApproved
		// Extract justification
		if idx := strings.Index(content, "CHECKPOINT_APPROVED:"); idx >= 0 {
			dec.Feedback = strings.TrimSpace(content[idx+len("CHECKPOINT_APPROVED:"):])
			if nl := strings.Index(dec.Feedback, "\n"); nl >= 0 {
				dec.Feedback = dec.Feedback[:nl]
			}
		}
	} else if strings.Contains(content, "CHECKPOINT_REJECTED") {
		dec.Action = checkpoint.DecisionRejected
		if idx := strings.Index(content, "CHECKPOINT_REJECTED:"); idx >= 0 {
			dec.Feedback = strings.TrimSpace(content[idx+len("CHECKPOINT_REJECTED:"):])
		}
	} else if strings.Contains(content, "CHECKPOINT_ESCALATED") {
		// Escalated — don't resolve, notify user
		log.Printf("[CHECKPOINT] CEO escalated — waiting for human: %s", cp.ID)
		sysMsg := message.NewMessage(message.TypeSystem, "ceo", "", "CEO escalated checkpoint — requires human review")
		sysMsg.Metadata.ProjectID = projectID
		o.store.Add(sysMsg)
		o.notify(sysMsg)
		return // Don't send to channel — keep blocking
	} else {
		// Default to approved if response is unclear
		dec.Action = checkpoint.DecisionApproved
		dec.Feedback = "Auto-approved (CEO response did not contain explicit markers)"
	}

	// Resolve on disk
	checkpoint.ResolveCheckpoint(o.config.ProjectDir, cp.ID, *dec)
	checkpoint.AddDecision(o.config.ProjectDir, *dec)

	// Send to channel
	select {
	case o.checkpointCh <- dec:
	default:
	}
}

// OnMessage registers a callback for new messages
// syncKanbanFromPlan creates kanban tasks from the development plan's subtasks
func (o *Orchestrator) syncKanbanFromPlan(plan *DevelopmentPlan, projectID string) {
	var phases []kanban.DevPlanPhase
	for _, p := range plan.Phases {
		phase := kanban.DevPlanPhase{Index: p.Index, Name: p.Name}
		for _, st := range p.SubTasks {
			agents := make([]string, len(st.AssignedAgents))
			for j, a := range st.AssignedAgents { agents[j] = string(a) }
			phase.SubTasks = append(phase.SubTasks, kanban.DevPlanSubTask{
				ID:                 st.ID,
				Title:              st.Title,
				Description:        st.Description,
				AssignedAgents:     agents,
				CompletionCriteria: st.CompletionCriteria,
			})
		}
		phases = append(phases, phase)
	}

	created, err := kanban.SyncFromDevPlan(o.config.ProjectDir, projectID, phases)
	if err != nil {
		log.Printf("[KANBAN] Failed to sync: %v", err)
	} else if created > 0 {
		log.Printf("[KANBAN] Synced %d tasks from development plan", created)
	}
}

func (o *Orchestrator) OnMessage(fn func(msg *message.Message)) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.listeners = append(o.listeners, fn)
}

// notify sends message to all listeners
func (o *Orchestrator) notify(msg *message.Message) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	for _, fn := range o.listeners {
		fn(msg)
	}
}

// notifyLifecycle broadcasts a lifecycle event (phase/task start/complete) via the message system.
func (o *Orchestrator) notifyLifecycle(eventType, projectID, taskID string, data map[string]interface{}) {
	content := eventType
	if name, ok := data["phase_name"]; ok {
		content = fmt.Sprintf("%s: %s", eventType, name)
	}
	msg := message.NewMessage("lifecycle", "orchestrator", "", content)
	msg.Metadata.ProjectID = projectID
	msg.Metadata.TaskID = taskID
	msg.Metadata.Extra = data
	msg.Metadata.Extra["event_type"] = eventType
	o.store.Add(msg)
	o.notify(msg)
}

// getAgent returns or creates an agent for the given role.
// If a SessionManager is configured, agents will use session-based execution.
func (o *Orchestrator) getAgent(role agent.Role) (*agent.Agent, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if a, ok := o.agents[role]; ok {
		return a, nil
	}

	a, err := agent.NewAgent(role)
	if err != nil {
		return nil, err
	}

	// Inject session manager if available
	if o.sessionManager != nil {
		a.SessionManager = o.sessionManager
	}

	o.agents[role] = a
	return a, nil
}

// ProcessTask handles a user task through the agent chain
func (o *Orchestrator) ProcessTask(projectID, taskStr string) (*Result, error) {
	return o.ProcessTaskWithParent(projectID, taskStr, "")
}

// ProcessTaskWithParent handles a task that may be continuing from a previous task
func (o *Orchestrator) ProcessTaskWithParent(projectID, taskStr, parentTaskID string) (*Result, error) {
	return o.ProcessTaskWithID(projectID, taskStr, parentTaskID, task.GenerateTaskID())
}

// ProcessTaskWithID handles a task with a pre-generated task ID
func (o *Orchestrator) ProcessTaskWithID(projectID, taskStr, parentTaskID, taskID string) (*Result, error) {
	o.turns = 0

	result := &Result{
		ProjectID:    projectID,
		TaskID:       taskID,
		ParentTaskID: parentTaskID,
		Task:         taskStr,
		StartTime:    time.Now(),
		Messages:     make([]*message.Message, 0),
		Files:        make([]FileResult, 0),
		VisitedRoles: make(map[agent.Role]bool),
		CreatedFiles: make(map[string]string),
	}

	// Create initial task message with TaskID
	taskMsg := message.NewTaskMessage(projectID, taskStr)
	taskMsg.Metadata.TaskID = taskID
	o.store.Add(taskMsg)
	o.notify(taskMsg)
	result.Messages = append(result.Messages, taskMsg)

	log.Printf("[TASK] ========================================")
	log.Printf("[TASK] New task: id=%s project=%s phases=%v", taskID, projectID, o.config.EnablePhases)
	log.Printf("[TASK] Task: %s", taskStr)
	log.Printf("[TASK] Config: maxTurns=%d maxDepth=%d files=%v verbose=%v", o.config.MaxTurns, o.config.MaxDepth, o.config.EnableFiles, o.config.Verbose)
	log.Printf("[TASK] ========================================")

	if o.config.Verbose {
		fmt.Printf("\n🎯 New Task: %s\n", taskStr)
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	}

	// Save task immediately with "running" status
	o.saveRunningTask(result)

	var err error

	// Use phase-based execution if enabled
	if o.config.EnablePhases {
		err = o.processWithPhases(taskStr, projectID, result)
	} else {
		// Legacy: Start with CEO (single-phase mode)
		err = o.processAgent(agent.RoleCEO, taskStr, projectID, 0, result)
	}

	if err != nil {
		log.Printf("[TASK] FAILED: %v", err)
		result.Error = err
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)
	result.TotalTurns = o.turns

	// Broadcast task_completed
	o.notifyLifecycle("task_completed", projectID, result.TaskID, map[string]interface{}{
		"turns":        o.turns,
		"files":        len(result.Files),
		"duration_ms":  result.Duration.Milliseconds(),
	})

	if o.config.Verbose {
		fmt.Printf("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("✅ Task Complete - %d turns, %d files created\n", o.turns, len(result.Files))
		fmt.Printf("⏱️  Duration: %s\n", result.Duration.Round(time.Second))
	}

	// Save task to history
	o.saveTaskToHistory(result)

	return result, err
}

// processWithPhases runs the 5-phase workflow: Template Selection → Research → Planning → Discussion → Development
func (o *Orchestrator) processWithPhases(taskStr, projectID string, result *Result) error {
	// Create .plans directory structure
	plansDir := filepath.Join(o.config.ProjectDir, ".plans")
	for _, subdir := range []string{"research", "specs", "final"} {
		if err := os.MkdirAll(filepath.Join(plansDir, subdir), 0755); err != nil {
			return fmt.Errorf("failed to create plans directory: %w", err)
		}
	}

	// Reset selected template to default for new task
	o.selectedTemplate = template.DefaultTemplate

	// ── CEO TRIAGE: Decide which phases to run ──
	triage, err := o.executeTriage(taskStr, projectID, result)
	if err != nil {
		log.Printf("[TRIAGE] Failed: %v, using full pipeline", err)
		triage = defaultTriage()
	}

	if o.config.Verbose && triage.Direction != "" {
		fmt.Printf("\n🎯 CEO Direction: %s\n", triage.Direction)
		fmt.Printf("   Task type: %s | Skipping: %v\n", triage.TaskType, triage.SkipPhases)
	}

	// Broadcast triage result
	triageMsg := message.NewMessage(message.TypeSystem, "ceo", "all",
		fmt.Sprintf("Triage: %s — %s", triage.TaskType, triage.Direction))
	triageMsg.Metadata.ProjectID = projectID
	triageMsg.Metadata.TaskID = result.TaskID
	o.store.Add(triageMsg)
	o.notify(triageMsg)

	// If the CEO's direction provides context, prepend it to the task for downstream agents
	enrichedTask := taskStr
	if triage.Direction != "" {
		enrichedTask = fmt.Sprintf("CEO Direction: %s\n\nOriginal Task: %s", triage.Direction, taskStr)
	}

	// Build phase list based on triage (skip what CEO says to skip)
	standardPhases := []Phase{}
	for _, p := range []Phase{PhaseTemplateSelection, PhaseResearch, PhasePlanning, PhaseDiscussion} {
		if !triage.shouldSkipPhase(p) {
			standardPhases = append(standardPhases, p)
		} else if o.config.Verbose {
			fmt.Printf("   ⏭  Skipping %s (CEO decision)\n", GetPhaseName(p))
		}
	}

	for _, phase := range standardPhases {
		o.currentPhase = phase

		if o.config.Verbose {
			fmt.Printf("\n%s PHASE: %s\n", GetPhaseEmoji(phase), GetPhaseName(phase))
			fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		}

		// Notify about phase change
		phaseMsg := message.NewMessage(message.TypeSystem, "orchestrator", "all", fmt.Sprintf("Entering %s phase", GetPhaseName(phase)))
		phaseMsg.Metadata.ProjectID = projectID
		phaseMsg.Metadata.TaskID = result.TaskID
		o.store.Add(phaseMsg)
		o.notify(phaseMsg)
		result.Messages = append(result.Messages, phaseMsg)

		// Reset visited roles for each phase (agents can be called again in different phases)
		result.VisitedRoles = make(map[agent.Role]bool)

		// Execute the phase (use enriched task with CEO direction)
		phaseStart := time.Now()
		log.Printf("[PHASE] Starting phase: %s (project=%s)", GetPhaseName(phase), projectID)

		// Broadcast phase_started
		o.notifyLifecycle("phase_started", projectID, result.TaskID, map[string]interface{}{
			"phase": string(phase), "phase_name": GetPhaseName(phase),
		})

		if err := o.executePhase(phase, enrichedTask, projectID, result); err != nil {
			log.Printf("[PHASE] Phase %s FAILED after %v: %v", GetPhaseName(phase), time.Since(phaseStart), err)
			return fmt.Errorf("phase %s failed: %w", phase, err)
		}
		phaseDuration := time.Since(phaseStart)
		log.Printf("[PHASE] Phase %s completed in %v", GetPhaseName(phase), phaseDuration)

		// Broadcast phase_completed
		o.notifyLifecycle("phase_completed", projectID, result.TaskID, map[string]interface{}{
			"phase": string(phase), "phase_name": GetPhaseName(phase),
			"duration_ms": phaseDuration.Milliseconds(),
		})

		// ── Post-phase checkpoints ──
		switch phase {
		case PhaseTemplateSelection:
			templatePath := filepath.Join(o.config.ProjectDir, ".plans", "template.md")
			if dec, err := o.waitForCheckpoint(
				checkpoint.TypeTemplateApproval, 0, templatePath,
				fmt.Sprintf("Selected template: %s", o.selectedTemplate),
				projectID, result,
			); err != nil {
				return fmt.Errorf("template checkpoint failed: %w", err)
			} else if dec != nil && dec.Action == checkpoint.DecisionRejected {
				return fmt.Errorf("template selection rejected: %s", dec.Feedback)
			}
		case PhaseResearch:
			// Optional research review checkpoint
			if _, err := o.waitForCheckpoint(
				checkpoint.TypeResearchReview, 0,
				filepath.Join(o.config.ProjectDir, ".plans", "research"),
				"Research phase complete — review before planning?",
				projectID, result,
			); err != nil {
				return fmt.Errorf("research review checkpoint failed: %w", err)
			}
		case PhasePlanning:
			// Optional spec review checkpoint
			if _, err := o.waitForCheckpoint(
				checkpoint.TypeSpecReview, 0,
				filepath.Join(o.config.ProjectDir, ".plans", "specs"),
				"Specifications ready for review",
				projectID, result,
			); err != nil {
				return fmt.Errorf("spec review checkpoint failed: %w", err)
			}
		}
	}

	// ── CHECKPOINT: Plan Approval (after Discussion, before Development) ──
	planPath := filepath.Join(o.config.ProjectDir, ".plans", "final", "approved-plan.md")
	if dec, err := o.waitForCheckpoint(
		checkpoint.TypePlanApproval, 0, planPath,
		"Development plan ready for review",
		projectID, result,
	); err != nil {
		return fmt.Errorf("plan approval checkpoint failed: %w", err)
	} else if dec != nil && dec.Action == checkpoint.DecisionRejected {
		// Re-run discussion phase with feedback
		log.Printf("[CHECKPOINT] Plan rejected, re-running discussion with feedback: %s", dec.Feedback)
		result.VisitedRoles = make(map[agent.Role]bool)
		feedbackCtx := PhasePromptContext(PhaseDiscussion, taskStr) + "\n\n### Human Feedback\n" + dec.Feedback
		if err := o.processAgentInPhase(o.phaseConfig.DiscussionLeader, feedbackCtx, projectID, result); err != nil {
			log.Printf("[CHECKPOINT] Discussion re-run had issues: %v", err)
		}
	} else if dec != nil && dec.Action == checkpoint.DecisionOverridden && dec.OverrideData != "" {
		os.WriteFile(planPath, []byte(dec.OverrideData), 0644)
	}

	// After Discussion phase, load development plan
	plan, err := LoadDevelopmentPlan(o.config.ProjectDir)
	if err != nil {
		// No JSON plan found — auto-generate a simple one from the task
		log.Printf("[DEV-PLAN] No development-plan.json found, auto-generating from task")
		plan = autoGeneratePlan(enrichedTask, projectID, o.phaseConfig.DevelopmentAgents)
		if saveErr := SaveDevelopmentPlan(o.config.ProjectDir, plan); saveErr != nil {
			log.Printf("[DEV-PLAN] Failed to save auto-generated plan: %v", saveErr)
		}
	}

	// Auto-sync kanban tasks from the development plan
	o.syncKanbanFromPlan(plan, projectID)

	// Execute each development phase with QA loop
	for i := range plan.Phases {
		devPhase := &plan.Phases[i]
		maxIterations := 3

		// Update phase status to in progress
		if err := plan.UpdatePhaseStatus(o.config.ProjectDir, devPhase.Index, PhaseStatusInProgress); err != nil {
			if o.config.Verbose {
				fmt.Printf("⚠️  Failed to update phase status: %v\n", err)
			}
		}

		for iteration := 1; iteration <= maxIterations; iteration++ {
			devPhase.Iteration = iteration

			// Save iteration to plan
			if err := SaveDevelopmentPlan(o.config.ProjectDir, plan); err != nil {
				if o.config.Verbose {
					fmt.Printf("⚠️  Failed to save development plan: %v\n", err)
				}
			}

			// Execute development phase
			o.currentPhase = PhaseDevelopment
			if iteration == 1 {
				if err := o.executeDevelopmentPhase(devPhase, taskStr, projectID, result); err != nil {
					return fmt.Errorf("development phase %d failed: %w", devPhase.Index, err)
				}
			} else {
				o.currentPhase = PhaseDevIteration
				if err := o.executeDevelopmentIteration(devPhase, taskStr, projectID, result); err != nil {
					return fmt.Errorf("development iteration phase %d failed: %w", devPhase.Index, err)
				}
			}

			// Execute QA review
			o.currentPhase = PhaseQA
			if err := o.executeQAPhase(devPhase, taskStr, projectID, result); err != nil {
				return fmt.Errorf("QA phase %d failed: %w", devPhase.Index, err)
			}

			// Reload plan to get updated QA status
			plan, err = LoadDevelopmentPlan(o.config.ProjectDir)
			if err != nil {
				return fmt.Errorf("failed to reload development plan: %w", err)
			}
			devPhase = &plan.Phases[i]

			// Check QA result
			if devPhase.QAStatus == QAStatusApproved {
				// ── CHECKPOINT: Phase Gate (after QA approval, before next phase) ──
				if _, err := o.waitForCheckpoint(
					checkpoint.TypePhaseGate, devPhase.Index,
					filepath.Join(o.config.ProjectDir, ".plans", "development-plan.json"),
					fmt.Sprintf("Phase %d: %s passed QA — proceed?", devPhase.Index, devPhase.Name),
					projectID, result,
				); err != nil {
					return fmt.Errorf("phase gate checkpoint failed: %w", err)
				}

				// Mark phase as completed
				if err := plan.UpdatePhaseStatus(o.config.ProjectDir, devPhase.Index, PhaseStatusCompleted); err != nil {
					if o.config.Verbose {
						fmt.Printf("⚠️  Failed to update phase status: %v\n", err)
					}
				}
				if o.config.Verbose {
					fmt.Printf("✅ Phase %d completed successfully\n", devPhase.Index)
				}
				break // Move to next phase
			}

			// QA rejected - check iteration limit
			if iteration >= maxIterations {
				return fmt.Errorf("phase %d (%s) failed after %d iterations - QA feedback: %s",
					devPhase.Index, devPhase.Name, maxIterations, devPhase.QAFeedback)
			}

			// Increment iteration for next attempt
			if err := plan.IncrementIteration(o.config.ProjectDir, devPhase.Index); err != nil {
				if o.config.Verbose {
					fmt.Printf("⚠️  Failed to increment iteration: %v\n", err)
				}
			}

			// Reload plan after increment
			plan, err = LoadDevelopmentPlan(o.config.ProjectDir)
			if err != nil {
				return fmt.Errorf("failed to reload development plan: %w", err)
			}
			devPhase = &plan.Phases[i]

			if o.config.Verbose {
				fmt.Printf("\n🔄 Starting Iteration %d for Phase %d: %s\n", devPhase.Iteration, devPhase.Index, devPhase.Name)
			}
		}
	}

	// ── CHECKPOINT: Final Acceptance ──
	if _, err := o.waitForCheckpoint(
		checkpoint.TypeFinalAcceptance, 0,
		filepath.Join(o.config.ProjectDir, ".plans", "development-plan.json"),
		"All development phases completed — final acceptance",
		projectID, result,
	); err != nil {
		return fmt.Errorf("final acceptance checkpoint failed: %w", err)
	}

	// All phases completed successfully
	if o.config.Verbose {
		fmt.Printf("\n✅ All development phases completed and approved\n")
	}

	return nil
}

// executePhase runs a specific phase with appropriate agents
func (o *Orchestrator) executePhase(phase Phase, taskStr, projectID string, result *Result) error {
	// Build phase-specific context
	var phaseContext string
	if phase == PhaseDevelopment {
		// Inject template guidelines for development phase
		phaseContext = PhasePromptContextWithTemplate(phase, taskStr, o.selectedTemplate)
	} else {
		phaseContext = PhasePromptContext(phase, taskStr)
	}

	switch phase {
	case PhaseTemplateSelection:
		// Template selection phase: Architect proposes, CEO/PM review
		if err := o.executeTemplateSelection(taskStr, projectID, result); err != nil {
			return err
		}

	case PhaseResearch:
		// Research phase: PM, UX, UI, Security, Architect research in parallel
		if err := o.executeAgentsParallel(o.phaseConfig.ResearchAgents, phaseContext, projectID, result); err != nil {
			if o.config.Verbose {
				fmt.Printf("⚠️  Research phase had issues: %v\n", err)
			}
		}

	case PhasePlanning:
		// Planning phase: Same agents create specs based on research in parallel
		if err := o.executeAgentsParallel(o.phaseConfig.PlanningAgents, phaseContext, projectID, result); err != nil {
			if o.config.Verbose {
				fmt.Printf("⚠️  Planning phase had issues: %v\n", err)
			}
		}

	case PhaseDiscussion:
		// Discussion phase: CEO moderates and approves
		if err := o.processAgentInPhase(o.phaseConfig.DiscussionLeader, phaseContext, projectID, result); err != nil {
			return err
		}

	case PhaseDevelopment:
		// Development phase: Developers implement based on approved specs in parallel
		if err := o.executeAgentsParallel(o.phaseConfig.DevelopmentAgents, phaseContext, projectID, result); err != nil {
			if o.config.Verbose {
				fmt.Printf("⚠️  Development phase had issues: %v\n", err)
			}
		}
	}

	return nil
}

// executeTemplateSelection handles Phase 0: Template Selection
func (o *Orchestrator) executeTemplateSelection(taskStr, projectID string, result *Result) error {
	phaseContext := PhasePromptContext(PhaseTemplateSelection, taskStr)

	// Step 1: Architect proposes template
	if o.config.Verbose {
		fmt.Printf("   🏗️  Architect is selecting template...\n")
	}
	if err := o.processAgentInPhase(o.phaseConfig.TemplateProposer, phaseContext, projectID, result); err != nil {
		if o.config.Verbose {
			fmt.Printf("⚠️  Architect template selection had issues: %v\n", err)
		}
	}

	// Step 2: CEO and PM review
	for _, role := range o.phaseConfig.TemplateReviewers {
		if o.config.Verbose {
			fmt.Printf("   %s %s is reviewing template...\n", getAgentEmoji(role), role)
		}
		if err := o.processAgentInPhase(role, phaseContext, projectID, result); err != nil {
			if o.config.Verbose {
				fmt.Printf("⚠️  %s template review had issues: %v\n", role, err)
			}
		}
	}

	// Step 3: Parse selected template from .plans/template.md if it exists
	templateFile := filepath.Join(o.config.ProjectDir, ".plans", "template.md")
	if content, err := os.ReadFile(templateFile); err == nil {
		o.selectedTemplate = o.parseTemplateFromContent(string(content))
		if o.config.Verbose {
			fmt.Printf("   ✅ Template selected: %s\n", o.selectedTemplate)
		}
	} else {
		// Fall back to default
		o.selectedTemplate = template.DefaultTemplate
		if o.config.Verbose {
			fmt.Printf("   ℹ️  Using default template: %s\n", o.selectedTemplate)
		}
	}

	return nil
}

// parseTemplateFromContent extracts the template type from template.md content
func (o *Orchestrator) parseTemplateFromContent(content string) template.TemplateType {
	// Look for patterns like:
	// "Chosen Template: `static-enhanced`"
	// "TEMPLATE_APPROVED: static-enhanced"
	// "## Chosen Template: static-html"

	patterns := []string{
		`Chosen Template:\s*` + "`" + `?(static-html|static-enhanced|nextjs-frontend|go-api|fullstack)` + "`" + `?`,
		`TEMPLATE_APPROVED:\s*(static-html|static-enhanced|nextjs-frontend|go-api|fullstack)`,
		`template:\s*(static-html|static-enhanced|nextjs-frontend|go-api|fullstack)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(`(?i)` + pattern)
		matches := re.FindStringSubmatch(content)
		if len(matches) > 1 {
			return template.ParseTemplateType(strings.ToLower(matches[1]))
		}
	}

	// Default if not found
	return template.DefaultTemplate
}

// processAgentInPhase processes a single agent within a phase context
func (o *Orchestrator) processAgentInPhase(role agent.Role, phaseContext, projectID string, result *Result) error {
	// Skip if already processed in this phase
	if result.VisitedRoles[role] {
		return nil
	}

	// Mark as visited
	result.VisitedRoles[role] = true
	o.turns++

	// Get agent
	a, err := o.getAgent(role)
	if err != nil {
		return fmt.Errorf("failed to get agent %s: %w", role, err)
	}

	if o.config.Verbose {
		fmt.Printf("\n%s %s is working...\n", getAgentEmoji(role), a.Name)
	}

	// Process with phase context — use session if available, else legacy
	var response *agent.Response
	var fileResults []agent.FileResult

	log.Printf("[ORCHESTRATOR] Agent %s starting (turn %d, files=%v, session=%v)", role, o.turns, o.config.EnableFiles, a.SessionManager != nil)

	if a.SessionManager != nil {
		// Snapshot files before execution for diff
		filesBefore := snapshotFiles(o.config.ProjectDir)

		// Session-based execution — full Claude Code with tools
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		sessResp, sessErr := a.ExecuteWithSession(ctx, projectID, o.config.ProjectDir, phaseContext)
		cancel()
		if sessErr != nil {
			err = sessErr
		} else if sessResp != nil {
			response = &agent.Response{
				Content:    sessResp.Text,
				DelegateTo: sessResp.Delegations,
				ReviewTo:   sessResp.Reviews,
				IsComplete: sessResp.IsComplete,
				Raw:        sessResp.Text,
			}

			// Use tool events for file tracking if available
			response.Files = append(sessResp.FilesCreated, sessResp.FilesModified...)
			for _, f := range sessResp.FilesCreated {
				fileResults = append(fileResults, agent.FileResult{Path: f, Success: true, Action: "create"})
			}
			for _, f := range sessResp.FilesModified {
				fileResults = append(fileResults, agent.FileResult{Path: f, Success: true, Action: "modify"})
			}

			// Fallback: filesystem diff to catch files created by Claude directly
			filesAfter := snapshotFiles(o.config.ProjectDir)
			for path := range filesAfter {
				if _, existed := filesBefore[path]; !existed {
					// New file found on disk
					if !containsPath(response.Files, path) {
						response.Files = append(response.Files, path)
						fileResults = append(fileResults, agent.FileResult{Path: path, Success: true, Action: "create"})
					}
				}
			}
		}
	} else if o.config.EnableFiles {
		response, fileResults, err = a.ProcessWithFileOps(phaseContext, o.config.ProjectDir)
	} else {
		response, err = a.Process(phaseContext)
	}

	if err != nil {
		log.Printf("[ORCHESTRATOR] Agent %s FAILED at turn %d: %v", role, o.turns, err)
		return fmt.Errorf("agent %s failed: %w", role, err)
	}
	log.Printf("[ORCHESTRATOR] Agent %s completed (turn %d), response: %d bytes, files: %d", role, o.turns, len(response.Content), len(fileResults))

	// Record response
	responseMsg := message.NewMessage(message.TypeResponse, string(role), "orchestrator", response.Content)
	responseMsg.Metadata.ProjectID = projectID
	responseMsg.Metadata.TaskID = result.TaskID
	o.store.Add(responseMsg)
	o.notify(responseMsg)
	result.Messages = append(result.Messages, responseMsg)

	if o.config.Verbose {
		content := response.Content
		if len(content) > 300 {
			content = content[:300] + "..."
		}
		fmt.Printf("   └─ %s\n", content)
	}

	// Handle file operations (lock for parallel safety)
	result.filesMu.Lock()
	for _, fr := range fileResults {
		if prevAgent, exists := result.CreatedFiles[fr.Path]; exists {
			if o.config.Verbose {
				fmt.Printf("   ⚠️  Skipping %s - already created by %s\n", fr.Path, prevAgent)
			}
			continue
		}

		fileMsg := message.NewMessage(message.TypeFileCreate, string(role), "filesystem", fr.Path)
		fileMsg.Metadata.ProjectID = projectID
		fileMsg.Metadata.TaskID = result.TaskID
		fileMsg.Metadata.FilePath = fr.Path

		if fr.Success {
			result.CreatedFiles[fr.Path] = string(role)
			if o.config.Verbose {
				fmt.Printf("   📄 Created: %s\n", fr.Path)
			}
		} else {
			fileMsg.Type = message.TypeError
			if o.config.Verbose {
				fmt.Printf("   ❌ Failed: %s - %v\n", fr.Path, fr.Error)
			}
		}

		o.store.Add(fileMsg)
		o.notify(fileMsg)

		result.Files = append(result.Files, FileResult{
			Path:    fr.Path,
			Success: fr.Success,
			Error:   fr.Error,
			Agent:   string(role),
		})
	}
	result.filesMu.Unlock()

	return nil
}

// executeAgentsParallel executes multiple agents in parallel within a phase
func (o *Orchestrator) executeAgentsParallel(roles []agent.Role, phaseContext, projectID string, result *Result) error {
	// Filter out already visited roles
	rolesToExecute := make([]agent.Role, 0, len(roles))
	for _, role := range roles {
		if !result.VisitedRoles[role] {
			rolesToExecute = append(rolesToExecute, role)
		}
	}

	if len(rolesToExecute) == 0 {
		return nil // All agents already processed
	}

	// Mark all agents as visited
	o.resultMu.Lock()
	for _, role := range rolesToExecute {
		result.VisitedRoles[role] = true
		o.turns++
	}
	o.resultMu.Unlock()

	// Create executor config
	executorConfig := ExecutorConfig{
		EnableParallel: true,
		MaxWorkers:     0, // Unlimited
		AgentTimeout:   5 * time.Minute,
		PhaseTimeout:   15 * time.Minute,
		WorkDir:        o.config.ProjectDir,
		ProjectID:      projectID,
		TaskID:         result.TaskID,
	}

	// Create phase executor
	executor := NewPhaseExecutor(executorConfig, o.getAgent, o.fileManager, o.store)

	// Execute agents in parallel
	if o.config.Verbose {
		fmt.Printf("\n   🚀 Executing %d agents in parallel...\n", len(rolesToExecute))
	}

	phaseResult, err := executor.ExecuteAgentsParallel(rolesToExecute, phaseContext)
	if err != nil {
		return fmt.Errorf("parallel execution failed: %w", err)
	}

	// Aggregate results
	o.resultMu.Lock()
	result.Messages = append(result.Messages, phaseResult.Messages...)
	o.resultMu.Unlock()

	// Handle file operations from all agents
	for _, msg := range phaseResult.Messages {
		// Notify listeners
		o.notify(msg)

		if o.config.Verbose {
			agentName := msg.From
			if a, err := o.getAgent(agent.Role(msg.From)); err == nil {
				agentName = a.Name
			}
			content := msg.Content
			if len(content) > 200 {
				content = content[:200] + "..."
			}
			fmt.Printf("\n%s %s completed:\n   └─ %s\n", getAgentEmoji(agent.Role(msg.From)), agentName, content)
		}
	}

	// Report errors
	if phaseResult.HasErrors() {
		if o.config.Verbose {
			fmt.Printf("⚠️  %d agents had errors:\n", phaseResult.ErrorCount())
			for _, ae := range phaseResult.Errors.Errors() {
				fmt.Printf("   - %s: %v\n", ae.Role, ae.Error)
			}
		}
	}

	// Report timing
	if o.config.Verbose {
		speedup := float64(phaseResult.MaxAgentDuration) / float64(phaseResult.Duration)
		fmt.Printf("   ⏱️  Phase completed in %v (max agent: %v, speedup: %.1fx)\n",
			phaseResult.Duration, phaseResult.MaxAgentDuration, speedup)
	}

	return nil
}

// saveRunningTask saves the task with "running" status when it starts
func (o *Orchestrator) saveRunningTask(result *Result) {
	meta := &task.TaskMetadata{
		TaskID:       result.TaskID,
		ParentTaskID: result.ParentTaskID,
		ProjectID:    result.ProjectID,
		Task:         result.Task,
		CreatedAt:    result.StartTime,
		Status:       task.StatusRunning,
		Turns:        0,
		Duration:     "0s",
		FilesCreated: []string{},
	}

	// Save to history (creates entry in history.json)
	err := o.historyManager.SaveTask(o.config.ProjectDir, meta, result.Messages)
	if err != nil && o.config.Verbose {
		fmt.Printf("⚠️  Failed to save running task: %v\n", err)
	}

	// Notify about task start
	startMsg := message.NewMessage(message.TypeSystem, "orchestrator", "all", "Task started")
	startMsg.Metadata.ProjectID = result.ProjectID
	startMsg.Metadata.TaskID = result.TaskID
	startMsg.Metadata.Tags = []string{"task_started"}
	o.notify(startMsg)
}

// saveTaskToHistory saves the completed task to history
func (o *Orchestrator) saveTaskToHistory(result *Result) {
	// Collect created file paths from CreatedFiles map (tracks all successfully created files)
	filePaths := make([]string, 0, len(result.CreatedFiles))
	for path := range result.CreatedFiles {
		filePaths = append(filePaths, path)
	}

	// Determine status
	status := task.StatusCompleted
	errMsg := ""
	if result.Error != nil {
		status = task.StatusFailed
		errMsg = result.Error.Error()
	}

	// Create metadata
	meta := &task.TaskMetadata{
		TaskID:       result.TaskID,
		ParentTaskID: result.ParentTaskID,
		ProjectID:    result.ProjectID,
		Task:         result.Task,
		CreatedAt:    result.StartTime,
		CompletedAt:  &result.EndTime,
		Status:       status,
		Turns:        result.TotalTurns,
		Duration:     result.Duration.Round(time.Second).String(),
		FilesCreated: filePaths,
		Error:        errMsg,
	}

	// Get messages for this task
	messages := o.store.GetByTaskID(result.TaskID)

	// Save to history
	err := o.historyManager.SaveTask(o.config.ProjectDir, meta, messages)
	if err != nil && o.config.Verbose {
		fmt.Printf("⚠️  Failed to save task history: %v\n", err)
	} else if o.config.Verbose {
		fmt.Printf("💾 Task saved to history: %s\n", result.TaskID)
	}

	// Notify about task completion
	completeMsg := message.NewMessage(message.TypeComplete, "orchestrator", "all",
		fmt.Sprintf("Task completed: %d turns, %d files created", result.TotalTurns, len(filePaths)))
	completeMsg.Metadata.ProjectID = result.ProjectID
	completeMsg.Metadata.TaskID = result.TaskID
	completeMsg.Metadata.Tags = []string{"task_completed"}
	o.store.Add(completeMsg)
	o.notify(completeMsg)
}

// processAgent handles a single agent's turn
func (o *Orchestrator) processAgent(role agent.Role, task, projectID string, depth int, result *Result) error {
	// Check limits
	if depth > o.config.MaxDepth {
		if o.config.Verbose {
			fmt.Printf("⚠️  Max depth reached (%d), stopping delegation\n", o.config.MaxDepth)
		}
		return nil
	}

	if o.turns >= o.config.MaxTurns {
		if o.config.Verbose {
			fmt.Printf("⚠️  Max turns reached (%d), stopping\n", o.config.MaxTurns)
		}
		return nil
	}

	// Check if agent already visited (prevent circular delegation)
	if result.VisitedRoles[role] {
		if o.config.Verbose {
			fmt.Printf("⚠️  Skipping %s - already processed (preventing circular delegation)\n", role)
		}
		return nil
	}

	// Mark agent as visited
	result.VisitedRoles[role] = true
	o.turns++

	// Get agent
	a, err := o.getAgent(role)
	if err != nil {
		return fmt.Errorf("failed to get agent %s: %w", role, err)
	}

	if o.config.Verbose {
		fmt.Printf("\n%s %s is thinking...\n", getAgentEmoji(role), a.Name)
	}

	// Process task — use session if available, else legacy
	var response *agent.Response
	var fileResults []agent.FileResult

	if a.SessionManager != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		sessResp, sessErr := a.ExecuteWithSession(ctx, projectID, o.config.ProjectDir, task)
		cancel()
		if sessErr != nil {
			err = sessErr
		} else if sessResp != nil {
			response = &agent.Response{
				Content:    sessResp.Text,
				DelegateTo: sessResp.Delegations,
				ReviewTo:   sessResp.Reviews,
				IsComplete: sessResp.IsComplete,
				Raw:        sessResp.Text,
			}
			response.Files = append(sessResp.FilesCreated, sessResp.FilesModified...)
		}
	} else if o.config.EnableFiles {
		response, fileResults, err = a.ProcessWithFileOps(task, o.config.ProjectDir)
	} else {
		response, err = a.Process(task)
	}

	if err != nil {
		return fmt.Errorf("agent %s failed: %w", role, err)
	}

	// Record response
	responseMsg := message.NewMessage(message.TypeResponse, string(role), "orchestrator", response.Content)
	responseMsg.Metadata.ProjectID = projectID
	responseMsg.Metadata.TaskID = result.TaskID
	o.store.Add(responseMsg)
	o.notify(responseMsg)
	result.Messages = append(result.Messages, responseMsg)

	if o.config.Verbose {
		// Print truncated response
		content := response.Content
		if len(content) > 300 {
			content = content[:300] + "..."
		}
		fmt.Printf("   └─ %s\n", content)

		// Log if agent signaled completion
		if response.IsComplete {
			fmt.Printf("   ✅ %s signaled task complete\n", a.Name)
		}
	}

	// Handle file operations (with deduplication)
	for _, fr := range fileResults {
		// Check if file was already created by another agent
		if prevAgent, exists := result.CreatedFiles[fr.Path]; exists {
			if o.config.Verbose {
				fmt.Printf("   ⚠️  Skipping %s - already created by %s\n", fr.Path, prevAgent)
			}
			continue
		}

		fileMsg := message.NewMessage(message.TypeFileCreate, string(role), "filesystem", fr.Path)
		fileMsg.Metadata.ProjectID = projectID
		fileMsg.Metadata.TaskID = result.TaskID
		fileMsg.Metadata.FilePath = fr.Path

		if fr.Success {
			// Track the file creation
			result.CreatedFiles[fr.Path] = string(role)
			if o.config.Verbose {
				fmt.Printf("   📄 Created: %s\n", fr.Path)
			}
		} else {
			fileMsg.Type = message.TypeError
			if o.config.Verbose {
				fmt.Printf("   ❌ Failed: %s - %v\n", fr.Path, fr.Error)
			}
		}

		o.store.Add(fileMsg)
		o.notify(fileMsg)

		result.Files = append(result.Files, FileResult{
			Path:    fr.Path,
			Success: fr.Success,
			Error:   fr.Error,
			Agent:   string(role),
		})
	}

	// Handle delegations (filter by hierarchy)
	validDelegations := agent.FilterValidDelegations(role, response.DelegateTo)

	// Log filtered delegations
	if o.config.Verbose && len(response.DelegateTo) != len(validDelegations) {
		for _, d := range response.DelegateTo {
			if !agent.CanDelegateTo(role, d) {
				fmt.Printf("   ⚠️  Ignoring invalid delegation: %s cannot delegate to %s\n", role, d)
			}
		}
	}

	if len(validDelegations) > 0 {
		for _, delegateRole := range validDelegations {
			// Create delegation message
			delegateMsg := message.NewMessage(
				message.TypeDelegate,
				string(role),
				string(delegateRole),
				fmt.Sprintf("Delegating from %s", a.Name),
			)
			delegateMsg.Metadata.ProjectID = projectID
			delegateMsg.Metadata.TaskID = result.TaskID
			delegateMsg.Metadata.DelegateTo = []string{string(delegateRole)}
			o.store.Add(delegateMsg)
			o.notify(delegateMsg)
			result.Messages = append(result.Messages, delegateMsg)

			if o.config.Verbose {
				fmt.Printf("   → Delegating to %s\n", delegateRole)
			}

			// Build context for delegated agent
			context := buildDelegationContext(a, response, task)

			// Recursively process
			err := o.processAgent(delegateRole, context, projectID, depth+1, result)
			if err != nil {
				// Log error but continue with other delegations
				errorMsg := message.NewMessage(message.TypeError, string(delegateRole), "orchestrator", err.Error())
				errorMsg.Metadata.ProjectID = projectID
				errorMsg.Metadata.TaskID = result.TaskID
				o.store.Add(errorMsg)
				o.notify(errorMsg)
			}
		}
	}

	// Handle review escalations (filter by hierarchy)
	validReviews := agent.FilterValidReviews(role, response.ReviewTo)

	// Log filtered reviews
	if o.config.Verbose && len(response.ReviewTo) != len(validReviews) {
		for _, r := range response.ReviewTo {
			if !agent.CanReviewTo(role, r) {
				fmt.Printf("   ⚠️  Ignoring invalid review: %s cannot escalate to %s\n", role, r)
			}
		}
	}

	if len(validReviews) > 0 {
		for _, reviewRole := range validReviews {
			// Create review message
			reviewMsg := message.NewMessage(
				message.TypeDelegate, // Use delegate type, could add TypeReview later
				string(role),
				string(reviewRole),
				fmt.Sprintf("Review request from %s", a.Name),
			)
			reviewMsg.Metadata.ProjectID = projectID
			reviewMsg.Metadata.TaskID = result.TaskID
			o.store.Add(reviewMsg)
			o.notify(reviewMsg)
			result.Messages = append(result.Messages, reviewMsg)

			if o.config.Verbose {
				fmt.Printf("   ↑ Review request to %s\n", reviewRole)
			}

			// Build review context
			context := buildReviewContext(a, response, task)

			// Process review (doesn't count against visited since it's upward)
			err := o.processAgent(reviewRole, context, projectID, depth+1, result)
			if err != nil {
				errorMsg := message.NewMessage(message.TypeError, string(reviewRole), "orchestrator", err.Error())
				errorMsg.Metadata.ProjectID = projectID
				errorMsg.Metadata.TaskID = result.TaskID
				o.store.Add(errorMsg)
				o.notify(errorMsg)
			}
		}
	}

	return nil
}

// buildReviewContext creates context for a review escalation
func buildReviewContext(from *agent.Agent, response *agent.Response, originalTask string) string {
	return fmt.Sprintf(`%s is requesting your review/decision.

ORIGINAL TASK:
%s

THEIR WORK:
%s

Please review and provide your decision or approval.
If approved, you can signal COMPLETE.
If changes needed, provide feedback.`,
		from.Name,
		originalTask,
		response.Content,
	)
}

// buildDelegationContext creates context for a delegated agent
func buildDelegationContext(from *agent.Agent, response *agent.Response, originalTask string) string {
	return fmt.Sprintf(`You are receiving a delegated task from %s.

ORIGINAL TASK:
%s

MESSAGE FROM %s:
%s

Please complete your part of this task according to your role and expertise.
If you need to delegate subtasks to other agents, use the DELEGATE format.
If you create files, use the appropriate file creation format.`,
		from.Name,
		originalTask,
		from.Name,
		response.Content,
	)
}

// getAgentEmoji returns an emoji for the agent role
// snapshotFiles returns a set of file paths under dir (excluding hidden dirs like .git, .plans, .tasks).
func snapshotFiles(dir string) map[string]bool {
	files := make(map[string]bool)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			base := filepath.Base(path)
			if strings.HasPrefix(base, ".") && path != dir {
				return filepath.SkipDir
			}
			return nil
		}
		files[path] = true
		return nil
	})
	return files
}

// containsPath checks if a path is in a slice.
func containsPath(paths []string, path string) bool {
	for _, p := range paths {
		if p == path {
			return true
		}
	}
	return false
}

func getAgentEmoji(role agent.Role) string {
	emojis := map[agent.Role]string{
		agent.RoleCEO:       "👔",
		agent.RolePM:        "📋",
		agent.RoleUX:        "🎨",
		agent.RoleUI:        "🖼️",
		agent.RoleSecurity:  "🔒",
		agent.RoleArchitect: "🏗️",
		agent.RoleSeniorDev: "👨‍💻",
		agent.RoleJuniorDev: "👩‍💻",
	}
	if e, ok := emojis[role]; ok {
		return e
	}
	return "🤖"
}

// Result holds the outcome of a task
type Result struct {
	ProjectID    string
	TaskID       string
	ParentTaskID string // ID of the task this continues from (if any)
	Task         string
	StartTime    time.Time
	EndTime      time.Time
	Duration     time.Duration
	TotalTurns   int
	Messages     []*message.Message
	Files        []FileResult
	Error        error
	VisitedRoles map[agent.Role]bool // Track visited agents to prevent circular delegation
	CreatedFiles map[string]string   // Track created files (path -> agent) to prevent duplicates
	filesMu      sync.Mutex          // Protects CreatedFiles and Files during parallel execution
}

// FileResult tracks a file operation result
type FileResult struct {
	Path    string
	Success bool
	Error   error
	Agent   string
}

// GetStore returns the message store
func (o *Orchestrator) GetStore() *message.Store {
	return o.store
}

// GetAgents returns all active agents
func (o *Orchestrator) GetAgents() map[agent.Role]*agent.Agent {
	o.mu.RLock()
	defer o.mu.RUnlock()

	// Return a copy
	agents := make(map[agent.Role]*agent.Agent)
	for k, v := range o.agents {
		agents[k] = v
	}
	return agents
}

// GetCurrentPhase returns the current workflow phase as a string
func (o *Orchestrator) GetCurrentPhase() string {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return string(o.currentPhase)
}

// executeSingleDevelopmentPhase executes development without a plan (fallback)
func (o *Orchestrator) executeSingleDevelopmentPhase(taskStr, projectID string, result *Result) error {
	o.currentPhase = PhaseDevelopment

	if o.config.Verbose {
		fmt.Printf("\n%s PHASE: %s\n", GetPhaseEmoji(PhaseDevelopment), GetPhaseName(PhaseDevelopment))
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	}

	phaseMsg := message.NewMessage(message.TypeSystem, "orchestrator", "all", "Entering Development phase")
	phaseMsg.Metadata.ProjectID = projectID
	phaseMsg.Metadata.TaskID = result.TaskID
	o.store.Add(phaseMsg)
	o.notify(phaseMsg)
	result.Messages = append(result.Messages, phaseMsg)

	result.VisitedRoles = make(map[agent.Role]bool)

	phaseContext := PhasePromptContextWithTemplate(PhaseDevelopment, taskStr, o.selectedTemplate)
	if err := o.executeAgentsParallel(o.phaseConfig.DevelopmentAgents, phaseContext, projectID, result); err != nil {
		if o.config.Verbose {
			fmt.Printf("⚠️  Development phase had issues: %v\n", err)
		}
	}

	return nil
}

// executeDevelopmentPhase executes a specific development phase with subtasks
func (o *Orchestrator) executeDevelopmentPhase(phase *DevelopmentPhase, taskStr, projectID string, result *Result) error {
	if o.config.Verbose {
		fmt.Printf("\n%s DEVELOPMENT PHASE %d: %s\n", GetPhaseEmoji(PhaseDevelopment), phase.Index, phase.Name)
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("Description: %s\n", phase.Description)
		fmt.Printf("Subtasks: %d\n", len(phase.SubTasks))
	}

	phaseMsg := message.NewMessage(
		message.TypeSystem,
		"orchestrator",
		"all",
		fmt.Sprintf("Development Phase %d: %s", phase.Index, phase.Name),
	)
	phaseMsg.Metadata.ProjectID = projectID
	phaseMsg.Metadata.TaskID = result.TaskID
	o.store.Add(phaseMsg)
	o.notify(phaseMsg)
	result.Messages = append(result.Messages, phaseMsg)

	result.VisitedRoles = make(map[agent.Role]bool)

	// Build development context with phase information
	phaseContext := buildDevelopmentPhaseContext(taskStr, phase, o.selectedTemplate)

	// Execute developers in parallel for this phase
	if err := o.executeAgentsParallel(o.phaseConfig.DevelopmentAgents, phaseContext, projectID, result); err != nil {
		if o.config.Verbose {
			fmt.Printf("⚠️  Development phase %d had issues: %v\n", phase.Index, err)
		}
	}

	return nil
}

// executeDevelopmentIteration executes a development iteration with QA feedback
func (o *Orchestrator) executeDevelopmentIteration(phase *DevelopmentPhase, taskStr, projectID string, result *Result) error {
	if o.config.Verbose {
		fmt.Printf("\n%s DEVELOPMENT ITERATION %d - Phase %d: %s\n",
			GetPhaseEmoji(PhaseDevIteration), phase.Iteration, phase.Index, phase.Name)
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("Addressing QA feedback from previous iteration\n")
	}

	phaseMsg := message.NewMessage(
		message.TypeSystem,
		"orchestrator",
		"all",
		fmt.Sprintf("Development Iteration %d - Phase %d: %s", phase.Iteration, phase.Index, phase.Name),
	)
	phaseMsg.Metadata.ProjectID = projectID
	phaseMsg.Metadata.TaskID = result.TaskID
	o.store.Add(phaseMsg)
	o.notify(phaseMsg)
	result.Messages = append(result.Messages, phaseMsg)

	result.VisitedRoles = make(map[agent.Role]bool)

	// Build iteration context with QA feedback
	iterationContext := buildIterationContext(taskStr, phase)

	// Execute developers in parallel for this iteration
	if err := o.executeAgentsParallel(o.phaseConfig.DevelopmentAgents, iterationContext, projectID, result); err != nil {
		if o.config.Verbose {
			fmt.Printf("⚠️  Development iteration %d had issues: %v\n", phase.Iteration, err)
		}
	}

	return nil
}

// buildDevelopmentPhaseContext creates context for a development phase
func buildDevelopmentPhaseContext(originalTask string, phase *DevelopmentPhase, templateType template.TemplateType) string {
	var sb strings.Builder

	sb.WriteString("## PHASE: DEVELOPMENT\n\n")
	sb.WriteString(fmt.Sprintf("### Development Phase %d: %s\n", phase.Index, phase.Name))
	sb.WriteString(fmt.Sprintf("**Description:** %s\n\n", phase.Description))

	sb.WriteString("### Original Task\n")
	sb.WriteString(originalTask)
	sb.WriteString("\n\n")

	sb.WriteString("### IMPORTANT: Read These First\n")
	sb.WriteString("Before writing ANY code, read these files:\n")
	sb.WriteString("1. `.plans/specs/product-spec.md` - Features and priorities\n")
	sb.WriteString("2. `.plans/specs/ux-spec.md` - Wireframes and flows\n")
	sb.WriteString("3. `.plans/specs/ui-spec.md` - Colors, fonts, spacing\n")
	sb.WriteString("4. `.plans/final/approved-plan.md` - CEO's approved plan\n")
	sb.WriteString("5. `.plans/development-plan.json` - Phase breakdown and subtasks\n\n")

	sb.WriteString("### Subtasks for This Phase\n\n")
	for i, subtask := range phase.SubTasks {
		sb.WriteString(fmt.Sprintf("%d. **%s**\n", i+1, subtask.Title))
		if subtask.Description != "" {
			sb.WriteString(fmt.Sprintf("   %s\n", subtask.Description))
		}
		sb.WriteString("   \n   **Assigned to:** ")
		for j, role := range subtask.AssignedAgents {
			if j > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(string(role))
		}
		sb.WriteString("\n\n   **Completion Criteria:**\n")
		for _, criterion := range subtask.CompletionCriteria {
			sb.WriteString(fmt.Sprintf("   - %s\n", criterion))
		}
		if len(subtask.Dependencies) > 0 {
			sb.WriteString("   \n   **Dependencies:** ")
			for j, dep := range subtask.Dependencies {
				if j > 0 {
					sb.WriteString(", ")
				}
				sb.WriteString(dep)
			}
		}
		sb.WriteString("\n\n")
	}

	sb.WriteString("### Your Assignment\n")
	sb.WriteString("Implement the subtasks assigned to you following the approved specs EXACTLY.\n\n")

	sb.WriteString("DO NOT:\n")
	sb.WriteString("- Deviate from the UI spec colors/fonts\n")
	sb.WriteString("- Skip reading the spec files\n")
	sb.WriteString("- Use generic AI colors (blue #3B82F6, green #10B981)\n")
	sb.WriteString("- Implement features not in this phase\n\n")

	sb.WriteString("DO:\n")
	sb.WriteString("- Use the EXACT colors from ui-spec.md\n")
	sb.WriteString("- Follow the wireframes from ux-spec.md\n")
	sb.WriteString("- Meet ALL completion criteria for your subtasks\n")
	sb.WriteString("- Work on subtasks assigned to your role\n\n")

	// Add template guidelines
	guidelines := template.GetGuidelines(templateType)
	sb.WriteString("### Template Guidelines\n")
	sb.WriteString(guidelines)
	sb.WriteString("\n\n")

	sb.WriteString("When done, signal: COMPLETE: Development phase finished\n")

	return sb.String()
}
