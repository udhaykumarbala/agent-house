# Agent House Mission Control -- Product Specification
## Human-in-the-Loop, Task Board, and Sub-task Decomposition

**Version:** 1.0
**Date:** 2026-03-16
**Author:** Product Management

---

## Table of Contents

1. [Executive Summary](#1-executive-summary)
2. [Human-in-the-Loop System](#2-human-in-the-loop-system)
3. [Task Board / Kanban](#3-task-board--kanban)
4. [Sub-task Decomposition & Tracking](#4-sub-task-decomposition--tracking)
5. [View Architecture](#5-view-architecture)
6. [Data Model](#6-data-model)
7. [API Endpoints](#7-api-endpoints)
8. [WebSocket Events](#8-websocket-events)
9. [Priority & Phasing](#9-priority--phasing)

---

## 1. Executive Summary

Agent House Mission Control currently runs as a fully autonomous 8-agent orchestration pipeline. The user submits a task, and the system executes Template Selection, Research, Planning, Discussion, Development, and QA phases without any human touchpoints.

This spec adds three capabilities that enterprise users require:

- **Human-in-the-Loop** -- The ability to pause, review, approve, reject, and provide feedback at critical decision points during the pipeline.
- **Task Board** -- A kanban-style board showing every work item across all agents with manual task management capabilities.
- **Sub-task Decomposition** -- Granular tracking of individual deliverables within each phase, with per-subtask review, dependencies, and artifact attachment.

All features are built with vanilla HTML/CSS/JS embedded in Go template strings, served through the existing `handleStatic` router, communicating via the existing WebSocket hub and new REST endpoints.

---

## 2. Human-in-the-Loop System

### 2.1 When the System Pauses for Human Input

The orchestrator currently moves through phases automatically. We introduce **Checkpoints** -- defined moments where execution halts and waits for human action. There are two types:

**Mandatory Checkpoints (always pause):**

| Checkpoint | Trigger | What the User Sees |
|---|---|---|
| Template Approval | After Architect proposes template, before Research begins | Template choice + rationale; Approve/Reject/Override buttons |
| Plan Approval | After Discussion phase, before Development begins | Full approved-plan.md rendered; Approve/Reject/Edit buttons |
| Phase Gate | After each Development Phase passes QA, before next phase starts | QA review summary; Approve to proceed / Request rework |
| Final Acceptance | After all development phases complete | Full deliverable summary; Accept / Request changes |

**Optional Checkpoints (user-configurable):**

| Checkpoint | Default | Description |
|---|---|---|
| Research Review | OFF | Pause after Research to review research docs before Planning |
| Spec Review | OFF | Pause after Planning to review individual specs before Discussion |
| Pre-QA Review | OFF | Pause before QA review to allow human QA first |

The user can toggle optional checkpoints via a "Workflow Settings" panel accessible from the header.

### 2.2 Checkpoint Interaction Model

When the orchestrator hits a checkpoint:

1. **Orchestrator pauses** -- the `processWithPhases` loop checks a `pendingCheckpoint` field before advancing to the next phase.
2. **WebSocket broadcasts** a `checkpoint_reached` event to all connected clients.
3. **The UI shows a Checkpoint Banner** -- a full-width overlay at the top of the canvas area with:
   - Title: "Approval Required: [checkpoint name]"
   - Summary of what was produced
   - Three action buttons: **Approve**, **Reject with Feedback**, **Override**
   - A collapsible section showing the full artifact (rendered markdown)
4. **User acts** -- their action is sent via REST API, the orchestrator unblocks and continues.

**Checkpoint Banner HTML Structure:**

```
.checkpoint-banner
  .checkpoint-badge       -- "WAITING FOR APPROVAL"
  .checkpoint-title       -- "Template Selection Complete"
  .checkpoint-summary     -- Short description of deliverable
  .checkpoint-content     -- Expandable: full rendered artifact
  .checkpoint-actions
    button.approve        -- "Approve & Continue"
    button.reject         -- "Reject with Feedback"
    button.override       -- "Override" (opens modal with custom input)
```

The banner appears overlaid on the canvas area (z-index 20), with a soft backdrop blur. It does not block the roster or log areas, so the user can still review agent activity and messages while deciding.

### 2.3 Reject with Feedback Flow

When the user clicks "Reject with Feedback":

1. A textarea modal appears with the prompt: "What needs to change?"
2. The user writes specific feedback.
3. The feedback is sent to the orchestrator, which:
   - Records the feedback as a `Decision` entity
   - Broadcasts a `checkpoint_feedback` WebSocket event
   - Re-enters the relevant phase with the feedback injected into agent prompts
   - The feedback text is prepended to the next phase's `phaseContext` as a `### Human Feedback` section

### 2.4 Override Flow

When the user clicks "Override":

1. A modal appears with the current artifact content in an editable textarea.
2. The user can modify the content directly (e.g., change the selected template, edit the plan).
3. On submit, the modified content replaces the artifact file on disk (e.g., `.plans/template.md`).
4. The orchestrator proceeds with the overridden content.

### 2.5 Chat with Individual Agents

Clicking an agent node on the canvas or in the roster opens the detail panel. We add a new **"Chat"** tab alongside the existing Activity, Tasks, and Metrics tabs.

**Chat Tab Design:**

```
.dp-chat
  .dp-chat-messages       -- Scrollable message history
    .chat-bubble.agent    -- Agent's messages (left-aligned, agent color)
    .chat-bubble.user     -- User's messages (right-aligned, accent color)
  .dp-chat-input-row
    textarea.chat-input   -- Multi-line input
    button.chat-send      -- "Send" button
```

**How it works:**

1. The user types a message in the chat input.
2. The message is sent via `POST /api/agents/{role}/chat` with the message content.
3. The orchestrator:
   - Creates a new `ChatMessage` entity
   - If the agent is currently idle, spawns the agent with the chat message as a prompt (prepended with conversation history)
   - If the agent is currently working, queues the message for delivery after the current task completes
4. The agent's response is captured and broadcast as a `chat_response` WebSocket event.
5. The chat tab updates in real-time.

**Constraints:**
- Chat messages are scoped to a project.
- Chat history is persisted in `.tasks/chat/{role}.json`.
- Chat does not interrupt an actively running phase. Messages are queued.
- When an agent is idle and receives a chat, the agent runs in a "consultation" mode -- it can read project files and provide analysis but does NOT create/modify files unless the user explicitly says "go ahead and implement this."

### 2.6 Feedback to Specific Agents Mid-Workflow

While agents are working during a phase, the user can send targeted feedback through the chat interface. This feedback is:

1. Stored as a `ChatMessage` with `context: "mid_workflow"`
2. Injected into the agent's next prompt cycle if it performs another iteration
3. Visible in the agent's Activity tab as a highlighted feedback entry

This is non-blocking -- it does not pause the workflow. It is advisory guidance.

---

## 3. Task Board / Kanban

### 3.1 Board Layout

The kanban board is a new top-level view, accessible via a tab in the header. The layout:

```
Header: [Logo] [Pipeline] [KPIs] [View: Topology | Kanban | Timeline] [Controls]

Kanban View:
+---------------+---------------+---------------+---------------+
|   BACKLOG     |  IN PROGRESS  |   IN REVIEW   |     DONE      |
+---------------+---------------+---------------+---------------+
| [Task Card]   | [Task Card]   | [Task Card]   | [Task Card]   |
| [Task Card]   | [Task Card]   |               | [Task Card]   |
| [Task Card]   |               |               | [Task Card]   |
|               |               |               | [Task Card]   |
| [+ Add Task]  |               |               |               |
+---------------+---------------+---------------+---------------+

Bottom: [Activity Log + Task Input]
```

The four columns:

| Column | Maps To | Description |
|---|---|---|
| **Backlog** | `pending`, `blocked` | Tasks waiting to be picked up, or blocked by dependencies |
| **In Progress** | `in_progress` | Tasks currently being executed by an agent |
| **In Review** | `in_review`, `needs_revision` | Tasks awaiting QA or human review |
| **Done** | `completed` | Successfully completed tasks |

Failed tasks appear in their last column with a red error indicator.

### 3.2 Task Card Design

Each task card shows:

```
.kanban-card
  .card-header
    .card-priority-dot     -- Color: red=critical, orange=high, blue=medium, gray=low
    .card-id               -- "ST-4.2" (phase 4, subtask 2)
    .card-agent-badge      -- Agent avatar (colored circle with initial)
  .card-title              -- "Implement user authentication API"
  .card-tags               -- Phase badge, type badge (research/planning/coding)
  .card-progress-bar       -- Thin bar showing 0-100% completion
  .card-footer
    .card-deps             -- Dependency count icon: "2 deps"
    .card-outputs          -- File count icon: "3 files"
    .card-time             -- Duration or time since creation
```

Card dimensions: ~240px wide, variable height. Cards use the agent's role color as a left border accent.

**Card States:**
- Default: `--bg-raised` background, `--border-subtle` border
- Hover: `--bg-hover` background, slight scale transform
- Dragging: Elevated shadow, 0.9 opacity, `--border-active` border
- Blocked: Dashed border, muted opacity, lock icon overlay
- Failed: Red left border, error icon

### 3.3 How Tasks Flow Between Columns

**Automatic flow (from orchestrator):**
- When the orchestrator creates a subtask -> Backlog
- When an agent starts working on it -> In Progress
- When an agent completes it and QA begins -> In Review
- When QA approves -> Done
- When QA rejects -> In Progress (with feedback badge)

**Manual flow (drag and drop):**
- User can drag cards between columns
- Dragging to "In Progress" triggers the orchestrator to assign the task to the card's agent
- Dragging to "Done" marks it complete (human override)
- Dragging to "Backlog" pauses/deprioritizes the task
- Dragging to "In Review" is only allowed for tasks that are "In Progress"

**Drag-and-drop implementation:** Pure JS using `dragstart`, `dragover`, `drop` events on cards and column containers. No external libraries.

### 3.4 Manual Task Creation

Clicking "+ Add Task" at the bottom of the Backlog column opens an inline card editor:

```
.kanban-new-card
  input.card-title-input     -- "Task title"
  textarea.card-desc-input   -- "Description (optional)"
  select.card-agent-select   -- Dropdown of 8 agents
  select.card-priority-select -- Critical / High / Medium / Low
  select.card-phase-select   -- Which phase this belongs to
  .card-actions
    button.card-save         -- "Create"
    button.card-cancel       -- "Cancel"
```

Manually created tasks get `source: "manual"` in their metadata, distinguishing them from auto-generated tasks.

### 3.5 Filtering, Sorting, Grouping

A toolbar above the columns provides:

```
.kanban-toolbar
  .filter-group
    select.filter-agent     -- "All Agents" / specific agent
    select.filter-phase     -- "All Phases" / specific phase
    select.filter-priority  -- "All" / Critical / High / Medium / Low
    input.filter-search     -- Text search across titles
  .group-sort
    select.group-by         -- "None" / "Agent" / "Phase" / "Priority"
    select.sort-by          -- "Created" / "Priority" / "Updated"
```

When grouped by Agent, each column is subdivided with agent-colored section headers. When grouped by Phase, sections show phase names.

### 3.6 Integration with Topology Canvas View

The topology canvas and kanban board are alternative views of the same data:

- Clicking an agent node on topology -> opens detail panel with that agent's tasks (same data as kanban filtered to that agent)
- The "View" switcher in the header toggles between: **Topology** (current canvas), **Kanban** (new board), **Timeline** (future, out of scope for now)
- Both views share the same bottom log area and header
- When switching views, the selected agent and filters are preserved

---

## 4. Sub-task Decomposition & Tracking

### 4.1 How Sub-tasks Are Created

**Auto-generated by agents (primary flow):**

The existing `DevelopmentPlan` already has `SubTask` structs within each `DevelopmentPhase`. Currently, the CEO agent creates these during the Discussion phase when it writes `development-plan.json`. We enhance this:

1. During Discussion phase, the CEO agent's prompt is augmented to produce more granular subtasks with explicit:
   - Unique IDs (e.g., `st-1-1` = phase 1, subtask 1)
   - Dependency declarations between subtasks
   - Estimated effort (small/medium/large)
   - Assigned agent(s)
   - Acceptance criteria (already exists as `completion_criteria`)
   - Expected outputs (file paths the subtask should produce)

2. After the Discussion phase writes the plan, the orchestrator parses it and creates corresponding `KanbanTask` records (see Data Model), one per subtask.

**Manually created by user:**

Via the kanban board "+ Add Task" button or via the sub-task panel on a parent task card. Manual subtasks:
- Can be attached to an existing phase or be phase-independent
- Are flagged `source: "manual"`
- Are included in the orchestrator's context for relevant agents

### 4.2 Sub-task Card Design (Expanded View)

Clicking a task card on the kanban board or in the detail panel opens an expanded view in the right panel:

```
.subtask-detail
  .subtask-header
    .subtask-id            -- "ST-2.3"
    .subtask-status-badge  -- "In Progress"
    .subtask-priority      -- Priority dot
  .subtask-title           -- Editable title
  .subtask-description     -- Editable description (markdown)

  .subtask-meta
    .meta-row: Agent       -- Agent avatar + name (click to reassign)
    .meta-row: Phase       -- Phase badge
    .meta-row: Created     -- Timestamp
    .meta-row: Updated     -- Timestamp
    .meta-row: Effort      -- Small / Medium / Large tag

  .subtask-criteria        -- Checklist of completion criteria
    .criterion-item        -- Checkbox + text (auto-checked from QA)

  .subtask-deps            -- "Dependencies" section
    .dep-card              -- Linked subtask mini-card (clickable)
    button.add-dep         -- "+ Add Dependency"

  .subtask-outputs         -- "Deliverables" section
    .output-card           -- File card with name, size, agent who created it
      button.view-file     -- Opens file content in a modal
    button.add-output      -- "+ Attach Output" (select from project files)

  .subtask-review          -- "Review" section
    .review-status         -- Pending / Approved / Rejected
    .review-feedback       -- QA feedback text
    button.approve-subtask -- "Approve" (human review)
    button.reject-subtask  -- "Request Changes" (opens feedback input)

  .subtask-history         -- Activity log for this subtask
    .history-entry         -- Timestamp + event description
```

### 4.3 Dependencies Between Sub-tasks

The existing `SubTask.Dependencies []string` field already supports this. We add visualization:

**On the kanban board:** Blocked subtasks show a lock icon and a tooltip listing what they depend on. Hovering a card highlights its dependencies (faded) and dependents (glowing).

**In the subtask detail panel:** A "Dependencies" section shows mini-cards of upstream tasks with their status. A "+ Add Dependency" button opens a dropdown of other subtasks in the same phase.

**Dependency enforcement:**
- The orchestrator already respects dependencies via `IsBlocked()` and `GetTopologicalOrder()`.
- On the kanban board, dragging a blocked task to "In Progress" shows a warning: "This task has unmet dependencies: [list]. Proceed anyway?"
- If the user overrides, the task is unblocked (`source: "manual_override"`).

### 4.4 Review/Approval Workflow Per Sub-task

Currently, QA review happens at the phase level. We extend to per-subtask:

1. When a subtask moves to "In Review" (auto or manual), a review badge appears on its card.
2. Opening the subtask detail shows the Review section with:
   - Each completion criterion as a checkbox
   - The QA agent's feedback (if any)
   - "Approve" and "Request Changes" buttons
3. Approving a subtask:
   - Sets `subtask.status = "completed"` and `subtask.review_status = "approved"`
   - Broadcasts a `subtask_approved` WebSocket event
   - The kanban card moves to "Done"
4. Rejecting a subtask:
   - Opens a feedback textarea
   - Sets `subtask.status = "needs_revision"` and `subtask.review_status = "rejected"`
   - The feedback is stored and injected into the next iteration context
   - The card moves back to "In Progress" with a revision badge

### 4.5 Progress Tracking and Rollup

**Subtask level:** 0-100% progress, updated by the orchestrator as agents work. Display as a thin progress bar on the kanban card.

**Phase level (rollup):** Phase progress = (completed subtasks / total subtasks) * 100. Displayed in the pipeline header and in the phase gate checkpoint.

**Project level (rollup):** Project progress = (completed phases / total phases) * 100. Displayed in the header KPI strip.

### 4.6 Deliverable Attachment

Each subtask can have attached outputs:

- **Auto-attached:** When an agent creates a file during a subtask's phase, the file is auto-linked to subtasks assigned to that agent in that phase.
- **Manual attachment:** User can click "+ Attach Output" and select from the project file browser.
- **Output cards** show: file name, file type icon, size, created-by agent, timestamp.
- Clicking an output card opens a file viewer modal (syntax-highlighted for code, rendered for markdown).

---

## 5. View Architecture

### 5.1 Updated Layout

The app shell grid evolves to support view switching:

```
.app (grid)
  Row 1: header (52px)
  Row 2: main-content (flex 1)
  Row 3: log-area (180px, collapsible)
```

The `main-content` row holds whichever view is active:

- **Topology View** (current): `roster | canvas | detail-panel`
- **Kanban View** (new): `kanban-filters | kanban-board | detail-panel`

The detail panel is shared between both views. When a task card is clicked in kanban view, the detail panel opens showing the subtask detail. When an agent card is clicked, it shows the agent detail.

### 5.2 Header Changes

The header gains a view switcher:

```
.view-switcher
  button.view-btn[data-view="topology"]  -- "Topology" (icon: nodes)
  button.view-btn[data-view="kanban"]    -- "Kanban" (icon: columns)
```

And a workflow settings button:

```
button.hdr-btn#settingsBtn  -- gear icon, opens checkpoint settings modal
```

### 5.3 Navigation Map

```
Header
  |-- View: Topology
  |     |-- Roster (left)
  |     |-- Canvas (center) -- click node -> Detail Panel opens
  |     |-- Detail Panel (right)
  |     |     |-- Tab: Activity (existing)
  |     |     |-- Tab: Tasks (existing, enhanced)
  |     |     |-- Tab: Chat (NEW)
  |     |     |-- Tab: Metrics (existing)
  |     |-- Checkpoint Banner (overlay on canvas)
  |
  |-- View: Kanban
  |     |-- Filter Toolbar (top)
  |     |-- Board Columns (center) -- click card -> Detail Panel opens with subtask detail
  |     |-- Detail Panel (right)
  |     |     |-- Subtask Detail view
  |     |     |-- OR Agent Detail view (if agent badge clicked)
  |
  |-- Log Area (bottom, shared across views)
  |     |-- Activity Log
  |     |-- Task Input
```

### 5.4 Where Chat/Feedback UI Lives

- **Agent chat:** In the detail panel, "Chat" tab. Available in both topology and kanban views when an agent is selected.
- **Checkpoint feedback:** In the checkpoint banner overlay. Available in topology view (primary) and as a notification banner in kanban view.
- **Subtask review feedback:** In the subtask detail panel, "Review" section. Available when a subtask card is selected in kanban view.

### 5.5 How the Detail Panel Evolves

Current detail panel tabs: Activity, Tasks, Metrics

New detail panel modes:

**Mode 1: Agent Detail** (when agent selected)
- Tab: Activity (existing -- shows recent messages to/from this agent)
- Tab: Tasks (enhanced -- shows kanban-style mini cards for this agent's tasks)
- Tab: Chat (NEW -- chat thread with this agent)
- Tab: Metrics (existing -- message counts, task counts)

**Mode 2: Subtask Detail** (when task card selected in kanban)
- Single view: subtask-detail (see section 4.2)
- Back button to return to agent detail or close panel

### 5.6 Served Route

The kanban view is NOT a separate page. It is part of the existing `/mission` SPA. View switching happens client-side by toggling visibility of `.topology-view` and `.kanban-view` containers. This avoids page reloads and maintains WebSocket connection.

---

## 6. Data Model

### 6.1 New Entities

#### Decision

Records human decisions at checkpoints.

```go
type Decision struct {
    ID             string         `json:"id"`              // "dec_20260316..."
    ProjectID      string         `json:"project_id"`
    TaskID         string         `json:"task_id"`         // Orchestrator task ID
    CheckpointType string         `json:"checkpoint_type"` // "template_approval", "plan_approval", "phase_gate", "final_acceptance"
    PhaseIndex     int            `json:"phase_index,omitempty"`
    Action         DecisionAction `json:"action"`          // "approved", "rejected", "overridden"
    Feedback       string         `json:"feedback,omitempty"`
    OverrideData   string         `json:"override_data,omitempty"` // Modified content if overridden
    DecidedBy      string         `json:"decided_by"`      // "user" (future: user ID)
    DecidedAt      time.Time      `json:"decided_at"`
}

type DecisionAction string
const (
    DecisionApproved  DecisionAction = "approved"
    DecisionRejected  DecisionAction = "rejected"
    DecisionOverridden DecisionAction = "overridden"
)
```

#### ChatMessage

Records chat between user and individual agents.

```go
type ChatMessage struct {
    ID        string    `json:"id"`         // "chat_20260316..."
    ProjectID string    `json:"project_id"`
    AgentRole string    `json:"agent_role"` // "architect", "pm", etc.
    Direction string    `json:"direction"`  // "user_to_agent" or "agent_to_user"
    Content   string    `json:"content"`
    Context   string    `json:"context"`    // "idle_chat", "mid_workflow", "checkpoint"
    Timestamp time.Time `json:"timestamp"`
}
```

#### KanbanTask

Extends the existing `SubTask` model with kanban-specific fields. Rather than replacing `SubTask`, we create a wrapper that references it.

```go
type KanbanTask struct {
    ID             string         `json:"id"`              // "kt_20260316..."
    ProjectID      string         `json:"project_id"`
    SubTaskRef     string         `json:"subtask_ref"`     // References SubTask.ID (if auto-generated)
    PhaseIndex     int            `json:"phase_index"`
    Title          string         `json:"title"`
    Description    string         `json:"description"`
    AssignedAgent  string         `json:"assigned_agent"`  // agent.Role
    Status         KanbanStatus   `json:"status"`
    Priority       TaskPriority   `json:"priority"`
    Source         string         `json:"source"`          // "auto" or "manual"
    ReviewStatus   ReviewStatus   `json:"review_status"`
    ReviewFeedback string         `json:"review_feedback,omitempty"`
    Dependencies   []string       `json:"dependencies,omitempty"` // Other KanbanTask IDs
    Outputs        []TaskOutput   `json:"outputs,omitempty"`
    Criteria       []Criterion    `json:"criteria,omitempty"`
    Effort         string         `json:"effort,omitempty"` // "small", "medium", "large"
    Progress       int            `json:"progress"`         // 0-100
    Position       int            `json:"position"`         // Sort order within column
    CreatedAt      time.Time      `json:"created_at"`
    UpdatedAt      time.Time      `json:"updated_at"`
    StartedAt      *time.Time     `json:"started_at,omitempty"`
    CompletedAt    *time.Time     `json:"completed_at,omitempty"`
}

type KanbanStatus string
const (
    KanbanBacklog    KanbanStatus = "backlog"
    KanbanInProgress KanbanStatus = "in_progress"
    KanbanInReview   KanbanStatus = "in_review"
    KanbanDone       KanbanStatus = "done"
    KanbanBlocked    KanbanStatus = "blocked"
)

type TaskPriority string
const (
    PriorityCritical TaskPriority = "critical"
    PriorityHigh     TaskPriority = "high"
    PriorityMedium   TaskPriority = "medium"
    PriorityLow      TaskPriority = "low"
)

type ReviewStatus string
const (
    ReviewPending  ReviewStatus = "pending"
    ReviewApproved ReviewStatus = "approved"
    ReviewRejected ReviewStatus = "rejected"
)

type TaskOutput struct {
    ID        string    `json:"id"`
    FilePath  string    `json:"file_path"`
    FileName  string    `json:"file_name"`
    FileType  string    `json:"file_type"` // "go", "html", "md", "json", etc.
    Size      int64     `json:"size"`
    CreatedBy string    `json:"created_by"` // agent role
    CreatedAt time.Time `json:"created_at"`
}

type Criterion struct {
    Text    string `json:"text"`
    Checked bool   `json:"checked"`
}
```

#### Checkpoint

Tracks pending and resolved checkpoints.

```go
type Checkpoint struct {
    ID             string          `json:"id"`              // "cp_20260316..."
    ProjectID      string          `json:"project_id"`
    TaskID         string          `json:"task_id"`
    Type           string          `json:"type"`            // Same as DecisionAction checkpoint_type
    PhaseIndex     int             `json:"phase_index,omitempty"`
    Status         string          `json:"status"`          // "pending", "resolved"
    ArtifactPath   string          `json:"artifact_path"`   // Path to the artifact file
    ArtifactSummary string         `json:"artifact_summary"` // Short text summary
    Decision       *Decision       `json:"decision,omitempty"` // Null while pending
    CreatedAt      time.Time       `json:"created_at"`
    ResolvedAt     *time.Time      `json:"resolved_at,omitempty"`
}
```

#### WorkflowSettings

User-configurable checkpoint toggles and auto-delegation.

```go
type WorkflowSettings struct {
    ProjectID                string `json:"project_id"`
    RequireTemplateApproval  bool   `json:"require_template_approval"`   // default: true
    RequirePlanApproval      bool   `json:"require_plan_approval"`       // default: true
    RequirePhaseGate         bool   `json:"require_phase_gate"`          // default: true
    RequireFinalAcceptance   bool   `json:"require_final_acceptance"`    // default: true
    RequireResearchReview    bool   `json:"require_research_review"`     // default: false
    RequireSpecReview        bool   `json:"require_spec_review"`         // default: false
    RequirePreQAReview       bool   `json:"require_pre_qa_review"`       // default: false

    // Auto-delegation: if a checkpoint is not resolved within this duration,
    // automatically delegate the decision to the CEO agent instead of blocking forever.
    // 0 = disabled (block indefinitely until human responds).
    AutoDelegateMinutes      int    `json:"auto_delegate_minutes"`       // default: 0 (off)
}
```

**Auto-delegation to CEO agent:**

When `AutoDelegateMinutes > 0` and a checkpoint has been pending for longer than the configured duration:

1. The orchestrator's checkpoint wait loop checks elapsed time on each tick (every 10s).
2. When timeout is reached:
   - A `checkpoint_auto_delegated` WebSocket event is broadcast.
   - The CEO agent is spawned with a prompt containing:
     - The full artifact being reviewed
     - The checkpoint type and context
     - Instruction: "The human reviewer has not responded within {X} minutes. Review this deliverable and decide: APPROVE, REJECT with specific feedback, or flag for mandatory human review."
   - The CEO agent's response is parsed for `CHECKPOINT_APPROVED` / `CHECKPOINT_REJECTED: {feedback}` / `CHECKPOINT_ESCALATED` markers.
   - If `CHECKPOINT_APPROVED`: checkpoint resolves as approved with `decided_by: "ceo_auto"`.
   - If `CHECKPOINT_REJECTED`: checkpoint resolves as rejected, feedback injected into next iteration, `decided_by: "ceo_auto"`.
   - If `CHECKPOINT_ESCALATED` (agent deems it too risky): checkpoint remains pending, a notification is sent to the user, and the timer resets.
3. The decision is recorded as a `Decision` entity with `decided_by: "ceo_auto"` so the audit trail clearly shows it was auto-delegated.
4. The UI shows a distinct badge on auto-delegated decisions: "Auto-approved by CEO" in amber, vs "Approved by user" in green.

**UI for auto-delegation settings:**

In the Workflow Settings modal, below the checkpoint toggles:

```
AUTO-DELEGATION
───────────────────────────────────────
When you're away, the CEO agent can review
checkpoints on your behalf.

  Auto-delegate after: [  5  ▾] minutes
                       [Off | 2 | 5 | 10 | 15 | 30]

  ● CEO will approve routine checkpoints
  ● CEO will escalate risky decisions back to you
  ● All auto-delegated decisions are marked in the audit log
```

**Checkpoint Banner update when timer is active:**

When auto-delegation is enabled, the checkpoint banner shows a countdown:

```
┌──────────────────────────────────────────────────────────┐
│  ⏳ WAITING FOR APPROVAL — Template Selection Complete    │
│                                                          │
│  Auto-delegating to CEO in 3:42                          │
│                                                          │
│  [Approve]  [Reject with Feedback]  [Override]           │
└──────────────────────────────────────────────────────────┘
```

The countdown updates every second. Clicking any button cancels the timer and resolves the checkpoint as a human decision.

**New WebSocket event:**

```json
{
  "type": "checkpoint",
  "event": {
    "event_type": "checkpoint_auto_delegated",
    "checkpoint_id": "cp_...",
    "decided_by": "ceo_auto",
    "action": "approved",
    "feedback": "Architecture plan is sound. Template selection aligns with project scope..."
  }
}
```

### 6.2 State Machines

**KanbanTask Status Transitions:**

```
                    +---> blocked ----+
                    |                 |
backlog ----> in_progress ----> in_review ----> done
  ^               |    ^            |
  |               |    |            |
  +--- (manual) --+    +-- (reject) +
```

Valid transitions:
- `backlog` -> `in_progress` (agent starts or manual drag)
- `backlog` -> `blocked` (dependency check)
- `blocked` -> `backlog` (dependency met or manual override)
- `in_progress` -> `in_review` (agent completes)
- `in_progress` -> `backlog` (manual deprioritize)
- `in_review` -> `done` (approved)
- `in_review` -> `in_progress` (rejected, needs revision)
- Any -> `done` (manual override by user)

**Checkpoint Status Transitions:**

```
pending ----> resolved
```

Simple: created as pending, resolved when user makes a decision.

### 6.3 Relationship Map

```
Orchestrator Task (1)
  |
  +---> DevelopmentPlan (1)
  |       |
  |       +---> DevelopmentPhase (N)
  |               |
  |               +---> SubTask (N)
  |                       |
  |                       +---> KanbanTask (1:1 mapping)
  |
  +---> Checkpoint (N) --- each resolved by ---> Decision (1)
  |
  +---> ChatMessage (N) --- per agent role
  |
  +---> KanbanTask (N) --- includes manual tasks without SubTask ref
```

### 6.4 Storage

All new entities are stored as JSON files within the project directory:

```
projects/{project_id}/
  .tasks/
    kanban.json           -- Array of KanbanTask
    checkpoints.json      -- Array of Checkpoint
    decisions.json        -- Array of Decision
    settings.json         -- WorkflowSettings
    chat/
      ceo.json            -- Array of ChatMessage
      pm.json
      architect.json
      ...
```

This is consistent with the existing pattern (`.plans/development-plan.json`, `.plans/qa-reviews.json`, `.tasks/history.json`).

---

## 7. API Endpoints

### 7.1 Checkpoint Endpoints

#### `GET /api/checkpoints`
Returns all checkpoints for a project.

Query params: `project` (required)

Response:
```json
{
  "checkpoints": [
    {
      "id": "cp_20260316...",
      "type": "template_approval",
      "status": "pending",
      "artifact_summary": "Selected template: static-enhanced",
      "created_at": "2026-03-16T10:00:00Z"
    }
  ],
  "count": 1,
  "pending_count": 1
}
```

#### `GET /api/checkpoints/{id}`
Returns a single checkpoint with full artifact content.

Response:
```json
{
  "checkpoint": { ... },
  "artifact_content": "# Template Selection\n\nChosen Template: `static-enhanced`\n..."
}
```

#### `POST /api/checkpoints/{id}/decide`
Resolves a checkpoint.

Request body:
```json
{
  "action": "approved",
  "feedback": "",
  "override_data": ""
}
```

Response:
```json
{
  "success": true,
  "decision_id": "dec_20260316...",
  "checkpoint_status": "resolved"
}
```

### 7.2 Chat Endpoints

#### `GET /api/agents/{role}/chat`
Returns chat history for an agent.

Query params: `project` (required), `limit` (optional, default 50)

Response:
```json
{
  "messages": [
    {
      "id": "chat_...",
      "direction": "user_to_agent",
      "content": "Can you explain the auth strategy?",
      "timestamp": "2026-03-16T10:05:00Z"
    },
    {
      "id": "chat_...",
      "direction": "agent_to_user",
      "content": "The authentication uses JWT tokens with...",
      "timestamp": "2026-03-16T10:05:30Z"
    }
  ],
  "count": 2
}
```

#### `POST /api/agents/{role}/chat`
Sends a message to an agent.

Request body:
```json
{
  "project_id": "subscription-manager",
  "content": "Can you explain the auth strategy?",
  "context": "idle_chat"
}
```

Response:
```json
{
  "success": true,
  "message_id": "chat_...",
  "status": "queued"
}
```

The agent's response arrives asynchronously via WebSocket (`chat_response` event).

### 7.3 Kanban Task Endpoints

#### `GET /api/kanban/tasks`
Returns all kanban tasks for a project.

Query params: `project` (required), `status` (optional), `agent` (optional), `phase` (optional), `priority` (optional)

Response:
```json
{
  "tasks": [ ... ],
  "count": 24,
  "columns": {
    "backlog": 8,
    "in_progress": 5,
    "in_review": 3,
    "done": 8
  }
}
```

#### `POST /api/kanban/tasks`
Creates a manual kanban task.

Request body:
```json
{
  "project_id": "subscription-manager",
  "title": "Add rate limiting to API",
  "description": "Implement rate limiting middleware...",
  "assigned_agent": "senior_dev",
  "priority": "high",
  "phase_index": 2,
  "criteria": [
    {"text": "Rate limiter middleware exists", "checked": false},
    {"text": "Configurable per-endpoint limits", "checked": false}
  ]
}
```

Response:
```json
{
  "success": true,
  "task": { ... }
}
```

#### `PATCH /api/kanban/tasks/{id}`
Updates a kanban task (status change, priority change, reassignment, position change).

Request body:
```json
{
  "status": "in_progress",
  "position": 2,
  "assigned_agent": "senior_dev"
}
```

Response:
```json
{
  "success": true,
  "task": { ... }
}
```

#### `POST /api/kanban/tasks/{id}/review`
Submit a review for a kanban task.

Request body:
```json
{
  "action": "approved",
  "feedback": "",
  "criteria_updates": [
    {"index": 0, "checked": true},
    {"index": 1, "checked": true}
  ]
}
```

Response:
```json
{
  "success": true,
  "task": { ... }
}
```

#### `POST /api/kanban/tasks/{id}/outputs`
Attach an output file to a task.

Request body:
```json
{
  "file_path": "index.html",
  "file_name": "index.html",
  "file_type": "html",
  "created_by": "senior_dev"
}
```

#### `POST /api/kanban/tasks/{id}/dependencies`
Add a dependency to a task.

Request body:
```json
{
  "depends_on": "kt_20260316..."
}
```

#### `DELETE /api/kanban/tasks/{id}/dependencies/{dep_id}`
Remove a dependency.

### 7.4 Workflow Settings Endpoints

#### `GET /api/settings/workflow`
Returns workflow settings for a project.

Query params: `project` (required)

Response:
```json
{
  "require_template_approval": true,
  "require_plan_approval": true,
  "require_phase_gate": true,
  "require_final_acceptance": true,
  "require_research_review": false,
  "require_spec_review": false,
  "require_pre_qa_review": false,
  "auto_delegate_minutes": 5
}
```

#### `PUT /api/settings/workflow`
Updates workflow settings.

Request body:
```json
{
  "project_id": "subscription-manager",
  "require_template_approval": true,
  "require_plan_approval": true,
  "require_phase_gate": true,
  "require_final_acceptance": true,
  "require_research_review": true,
  "require_spec_review": false,
  "require_pre_qa_review": false,
  "auto_delegate_minutes": 5
}
```

---

## 8. WebSocket Events

New event types broadcast through the existing Hub:

### 8.1 Checkpoint Events

```json
{
  "type": "checkpoint",
  "event": {
    "event_type": "checkpoint_reached",
    "checkpoint_id": "cp_...",
    "checkpoint_type": "template_approval",
    "artifact_summary": "Selected template: static-enhanced",
    "phase_index": 0
  }
}
```

```json
{
  "type": "checkpoint",
  "event": {
    "event_type": "checkpoint_resolved",
    "checkpoint_id": "cp_...",
    "action": "approved"
  }
}
```

### 8.2 Chat Events

```json
{
  "type": "chat",
  "event": {
    "event_type": "chat_response",
    "agent_role": "architect",
    "message": {
      "id": "chat_...",
      "content": "The authentication strategy uses...",
      "timestamp": "2026-03-16T10:05:30Z"
    }
  }
}
```

### 8.3 Kanban Events

```json
{
  "type": "kanban",
  "event": {
    "event_type": "task_moved",
    "task_id": "kt_...",
    "from_status": "in_progress",
    "to_status": "in_review"
  }
}
```

```json
{
  "type": "kanban",
  "event": {
    "event_type": "task_created",
    "task": { ... }
  }
}
```

```json
{
  "type": "kanban",
  "event": {
    "event_type": "task_reviewed",
    "task_id": "kt_...",
    "review_status": "approved"
  }
}
```

### 8.4 Implementation in Hub

The Hub struct gains new channels:

```go
type Hub struct {
    clients           map[*Client]bool
    broadcast         chan *message.Message
    agentTaskEvents   chan *AgentTaskEvent
    checkpointEvents  chan *CheckpointEvent    // NEW
    chatEvents        chan *ChatEvent           // NEW
    kanbanEvents      chan *KanbanEvent         // NEW
    register          chan *Client
    unregister        chan *Client
    mu                sync.RWMutex
}
```

Each new channel is handled in the `Run()` select loop with the same broadcast pattern as existing events.

---

## 9. Priority & Phasing

### Phase 1: Foundation (Build First)
**Goal:** Enable human control over the pipeline.
**Estimated effort:** 3-4 days

1. **Checkpoint data model + storage** -- `Decision`, `Checkpoint`, `WorkflowSettings` structs and JSON storage
2. **Checkpoint API endpoints** -- `GET/POST /api/checkpoints`, `POST /api/checkpoints/{id}/decide`, `GET/PUT /api/settings/workflow`
3. **Orchestrator integration** -- Modify `processWithPhases` to check for checkpoints and block on a Go channel until resolved
4. **Checkpoint Banner UI** -- HTML/CSS/JS for the approval banner overlay on the canvas
5. **Workflow Settings modal** -- Toggle checkpoints on/off
6. **WebSocket checkpoint events** -- Broadcast `checkpoint_reached` and `checkpoint_resolved`

**Why first:** This is the single most important feature for enterprise users. Without human-in-the-loop, the platform cannot be trusted with real projects. The checkpoint system is also the foundation for the review workflow in Phase 3.

### Phase 2: Kanban Board (Build Second)
**Goal:** Give users a task-centric view of all agent work.
**Estimated effort:** 4-5 days

1. **KanbanTask data model + storage** -- `KanbanTask` struct, JSON storage, sync with existing `SubTask` model
2. **Kanban API endpoints** -- Full CRUD for tasks, status changes, filtering
3. **Board UI** -- Four columns, task cards, drag-and-drop
4. **Filter toolbar** -- Agent, phase, priority, search
5. **Manual task creation** -- Inline card editor
6. **View switcher** -- Header toggle between Topology and Kanban
7. **WebSocket kanban events** -- Real-time card movement

**Why second:** The kanban board provides the visual framework that the subtask decomposition system (Phase 3) populates. Building the board first gives users immediate value and establishes the interaction patterns.

### Phase 3: Sub-task Decomposition + Reviews (Build Third)
**Goal:** Granular tracking with per-deliverable review.
**Estimated effort:** 3-4 days

1. **Enhanced subtask generation** -- Augment CEO Discussion prompt to produce richer subtasks with effort, outputs, and granular criteria
2. **Auto-sync subtasks to kanban** -- After Discussion phase, create `KanbanTask` for each `SubTask`
3. **Subtask detail panel** -- Full expanded view with criteria checklist, dependencies, outputs, review section
4. **Per-subtask review workflow** -- Approve/reject individual subtasks with feedback
5. **Output attachment** -- Auto-link files to subtasks, manual attachment
6. **Dependency visualization** -- Hover highlighting on kanban cards

**Why third:** This builds on both Phase 1 (review workflow) and Phase 2 (kanban cards). The subtask decomposition also requires the most changes to agent prompts, which benefits from having the checkpoint system already in place for testing.

### Phase 4: Agent Chat (Build Last)
**Goal:** Direct communication channel with individual agents.
**Estimated effort:** 2-3 days

1. **ChatMessage data model + storage** -- Chat history per agent per project
2. **Chat API endpoints** -- `GET/POST /api/agents/{role}/chat`
3. **Chat tab in detail panel** -- Message bubbles, input area
4. **Agent consultation mode** -- Spawning agents for ad-hoc questions
5. **Message queueing** -- Queue chat during active phases, deliver when agent is idle
6. **WebSocket chat events** -- Real-time message delivery

**Why last:** Chat is a nice-to-have compared to the core control loop (checkpoints), visibility (kanban), and tracking (subtasks). It also introduces the most complex concurrent behavior (queuing messages, consultation mode), so it benefits from the stability of the earlier phases.

### Summary Timeline

```
Week 1:  Phase 1 (Checkpoints + Human-in-the-Loop)
Week 2:  Phase 2 (Kanban Board)
Week 3:  Phase 3 (Sub-task Decomposition) + Phase 4 (Agent Chat)
```

Total estimated effort: 12-16 days for a single developer.

---

## Appendix A: File Changes Map

| File | Changes |
|---|---|
| `internal/web/server.go` | Add new route registrations for checkpoint, chat, kanban, settings endpoints |
| `internal/web/mission.go` | Add kanban board HTML/CSS/JS, checkpoint banner, chat tab, view switcher, subtask detail panel |
| `internal/web/websocket.go` | Add CheckpointEvent, ChatEvent, KanbanEvent channels and broadcast methods |
| `internal/orchestrator/orchestrator.go` | Add checkpoint blocking logic in `processWithPhases`, add `pendingCheckpoint` channel |
| `internal/orchestrator/dev_plan.go` | No changes (existing model is sufficient) |
| `internal/orchestrator/subtask.go` | Add `ReviewStatus`, `Effort`, `ExpectedOutputs` fields |
| NEW `internal/checkpoint/checkpoint.go` | Checkpoint, Decision, WorkflowSettings structs and storage |
| NEW `internal/chat/chat.go` | ChatMessage struct, storage, agent consultation runner |
| NEW `internal/kanban/kanban.go` | KanbanTask struct, storage, sync with SubTask |
| NEW `internal/kanban/manager.go` | KanbanManager: CRUD, filtering, status transitions, dependency validation |

## Appendix B: API Summary Table

| Method | Path | Description | Phase |
|---|---|---|---|
| GET | `/api/checkpoints` | List checkpoints | 1 |
| GET | `/api/checkpoints/{id}` | Get checkpoint detail | 1 |
| POST | `/api/checkpoints/{id}/decide` | Resolve checkpoint | 1 |
| GET | `/api/settings/workflow` | Get workflow settings | 1 |
| PUT | `/api/settings/workflow` | Update workflow settings | 1 |
| GET | `/api/kanban/tasks` | List kanban tasks | 2 |
| POST | `/api/kanban/tasks` | Create manual task | 2 |
| PATCH | `/api/kanban/tasks/{id}` | Update task | 2 |
| POST | `/api/kanban/tasks/{id}/review` | Review a task | 3 |
| POST | `/api/kanban/tasks/{id}/outputs` | Attach output | 3 |
| POST | `/api/kanban/tasks/{id}/dependencies` | Add dependency | 3 |
| DELETE | `/api/kanban/tasks/{id}/dependencies/{dep_id}` | Remove dependency | 3 |
| GET | `/api/agents/{role}/chat` | Get chat history | 4 |
| POST | `/api/agents/{role}/chat` | Send chat message | 4 |
