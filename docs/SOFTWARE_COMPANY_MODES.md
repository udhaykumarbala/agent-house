# Software Company — Run Modes & Two-Door Switcher

Design for turning the existing build pipeline into a fully agent-driven
software solution company with selectable autonomy modes. Decisions confirmed
with the user (2026-06-05).

## Context (what already exists)

- A real 6-phase SDLC runs in `internal/orchestrator/`: **Template → Research →
  Planning → Discussion (PRD) → Development → QA** (3-iteration loops + phase
  gates). It produces real apps under `projects/<id>/` with full artifacts
  (`.plans/research`, `specs`, `approved-plan.md`, `development-plan.json`, QA
  reviews).
- Human-in-the-loop **checkpoints** (`internal/checkpoint/`): template_approval,
  plan_approval, phase_gate, final_acceptance (+ optional research/spec/preQA).
  Stored in `<proj>/.tasks/{checkpoints,decisions,settings}.json`.
- **Auto-delegate timer already exists**: `WorkflowSettings.AutoDelegateMinutes`
  + `orchestrator.autoDelegateTimer` → after N minutes with no human response,
  the CEO agent reviews and decides (APPROVED/REJECTED/ESCALATED) **with a
  written justification**, recorded as `DecidedBy: "ceo_auto"`. This is ~70% of
  the requested autonomous mode.

## Gaps this work closes

1. No clean **run-mode model** (just a raw `auto_delegate_minutes` number).
2. No **UI** for the countdown, the choice, or the auto-decision reasoning.
3. No **software-company door** — the Conductor is hardwired to EPC (atlas-site,
   EPC role filter). Software builds work via `create_project` but aren't
   presented as their own company.
4. Conversations are named from the raw first message, not build/modify intent.

## Run modes (preset over the existing engine)

| Mode | Gates enabled | Behavior at a gate |
|------|---------------|--------------------|
| **Manual** | all (incl. research, spec, preQA) | wait for human, forever |
| **Semi-auto** | template, plan(PRD), phase_gate, final | wait for human, forever |
| **Full-auto** | template, plan(PRD), phase_gate, final | show choice + **countdown X min**; on expiry CEO decides with reasoning |
| **Blitz** | none | auto-approve immediately (logged decision), never block |

- `RunMode` + `DecisionTimeoutMinutes` added to `WorkflowSettings`.
- `ApplyRunModePreset(mode)` sets the `Require*` flags + timeout; default mode =
  `semi_auto`. Individual `Require*` toggles still editable after.
- `waitForCheckpoint` branches on `RunMode`:
  - `blitz` → write an approved Decision (`decided_by:"auto_blitz"`) and return
    immediately, no block.
  - `full_auto` → start `autoDelegateTimer(DecisionTimeoutMinutes)` (CEO decides
    with reasoning on expiry).
  - `manual` / `semi_auto` → block on `checkpointCh` (no timer).
- `checkpointNotify` carries the countdown **deadline** (full-auto) and the
  resolved **decision + reasoning** so the UI can show both.

## Settings

- Extend `PUT /api/settings/workflow` to accept `run_mode` +
  `decision_timeout_minutes`; server applies the preset.
- UI: a mode selector (4 chips) + a timeout input (minutes), in the software
  conductor / mission settings.

## Two-door switcher + Software conductor

- Top-bar **company switcher**: `Software Co` ⟷ `EPC · Atlas`. Stored in a
  client company context + threaded to `/api/chat` as `company`/`scope`.
- Software door: workforce = `ceo, pm, architect, ux, ui, security, senior_dev,
  junior_dev`; brain routing is **company-aware** — software requests route to
  `create_project` (new app) / `delegate` (modify) into the build pipeline and
  do **not** fire EPC scenarios.
- Conversation naming: a software create → `Build: <app>`; a modify → `Modify:
  <app>` (derived from the routed decision).

## Build order

1. Run-mode model in `checkpoint` (#23)
2. Orchestrator honors RunMode (#24)
3. Settings API + UI (#25)
4. Checkpoint card: countdown + reasoning (#26)
5. Company switcher + software conductor (#27)
6. Conversation naming by intent (#28)
