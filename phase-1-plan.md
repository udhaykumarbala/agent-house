# Agent House — Phase 1 Plan

> **Compose any AI workforce. Add agents on demand. They react to your business events.**

**Status:** Approved direction · **Scope:** Phase 1 · **Data layer:** Filesystem, single-org · **Date:** 2026-05-16

This document is layered: **Part I** is the product narrative, **Part II** catalogs every page and feature, **Part III** is the functional specification, **Part IV** is the engineering execution plan. Parts I–II are for product/stakeholders; Parts III–IV are for engineering.

---

## Part I — Product Narrative

### 1. What Phase 1 is

Agent House today is a fixed 8-agent software-development pipeline. Phase 1 turns it into a **composable, reactive AI workforce platform**:

1. **Composable** — teams are assembled on demand from *industry packs* or *custom agents*, not hardcoded. Ships with two default packs: **EPC** (Engineering, Procurement, Construction) and **IT** (software development).
2. **Reactive** — the headline AI advancement. Agents stop waiting for a human to press "go." Inbound email, webhooks, schedules, and dropped files **wake the right agent automatically**, gated by the existing human-in-the-loop checkpoints.

Phase 1 deliberately stays on the **filesystem, single-org** runtime. No PostgreSQL, no multi-tenant, no generic tool catalog. The bet: prove "compose + react" with the smallest delta to what already works, on real pilot projects.

### 2. Positioning

| | Today (v1) | After Phase 1 |
|---|---|---|
| Team | Fixed 8 software roles | Any roster, composed per project from packs or custom agents |
| Industries | Software only | EPC + IT default packs, extensible by the user |
| Work trigger | Human submits a task | Email / webhook / cron / file-drop wake agents; humans still gate decisions |
| Adding an agent | Code change + recompile | Drop a folder or fill a form in the UI — no recompile |

Not a chatbot, not a copilot — an AI workforce you assemble and that responds to your operation's events.

### 3. The two default packs

**EPC pack** (`industries/epc`) — already scaffolded in `agents/`:
`project_manager` (hub) · `site_engineer` · `procurement` · `hse` · `qa_inspector` · `hr`
Hierarchy already wired in `internal/agent/hierarchy.go` (Project Manager delegates to Site Engineer/Procurement/HSE/QA/HR; review escalates up to Project Manager → CEO).

**IT pack** (`industries/it`) — the existing tech team, repackaged unchanged:
`ceo` · `pm` · `ux` · `ui` · `security` · `architect` · `senior_dev` · `junior_dev`

The **Brain/Orchestrator is pack-agnostic** — it routes natural language and runs the phase pipeline regardless of which roster a project uses. A pack supplies a *roster + delegation/review graph + workflow phases + default triggers*; the engine stays the same.

> **`ceo` is a platform role, not an IT-pack member.** The existing code (and the EPC review hierarchy) escalates to `ceo` as the top-level reviewer. Phase 1 keeps this: `ceo` is a built-in **cross-pack escalation/approval role** owned by the platform, available to every pack's review chain, and is *not* counted as part of the IT roster for composition purposes. The registry marks it `pack: "_platform"`.

### 4. Personas

- **Marcus — EPC Project Director.** Coordinates design/procurement/construction/HSE/QA across time zones. Wants RFIs, invoices, and inspection schedules to *route themselves* to the right discipline and land on his desk as decisions, not raw email.
- **Alex — Engineering Manager (IT).** Wants a GitHub PR or an incident webhook to wake the security/dev agents without anyone filing a ticket.
- **Dana — AI Workforce Builder (new).** Ops/program lead who composes teams: picks a pack, adds a custom agent (e.g., "Contracts Analyst"), wires its triggers, and turns it loose. Never edits Go.

### 5. The AI advancement: event-driven agents

The system moves from **human-initiated → reactive**. Four trigger classes in Phase 1:

- **Email** — inbound message (via the existing Email Simulator / inbox) matched by sender/subject/category → spawns the bound agent. *EPC: invoice → Procurement; RFI → Site Engineer. IT: incident report → Security.*
- **Webhook** — external system POSTs (extends working `hooks.go`). *IT: GitHub PR opened → Security + Senior Dev review.*
- **Cron/Schedule** — recurring (extends working `cron.go`). *EPC: daily 07:00 → Project Manager schedule-slippage check; weekly → HSE compliance sweep.*
- **File-watch** — a file dropped into a watched dir (drawings, invoices, inspection sheets) → bound agent. *EPC: drawing revision → Site Engineer; invoice PDF → Procurement.*

Every trigger-spawned task flows through the **existing checkpoint/HITL path** — reactive does not mean unsupervised. Dangerous/external actions still pause for human approval.

### 6. Success metrics (Phase 1)

- A user can add a working custom agent through the UI in < 5 minutes, no recompile.
- A user can bind a trigger and see it fire → produce a checkpoint-gated task end-to-end.
- ≥ 1 EPC pilot project running with ≥ 2 live triggers (e.g., invoice file-watch + daily PM cron).
- Zero hardcoded role enums remain in `agent.go`/`hierarchy.go`/`modes.go`/`permissions.go` — all registry-backed.
- Trigger-spawned tasks are 100% audited and 100% checkpoint-gated for dangerous actions.

---

## Part II — Pages & Features

Legend: **(keep)** unchanged · **(enhance)** modified in Phase 1 · **(new)** built in Phase 1.

### Existing surfaces

| Page | Route | Phase 1 | Notes |
|---|---|---|---|
| Mission v2 | `/mission` | **enhance** | Primary UI. Sidebar becomes pack-aware; activity feed tags trigger-initiated work; pack/team indicator in header. |
| Project Detail | `/project/{id}` | keep | Files/artifacts viewer. |
| Live Activity | `/live` | enhance | Trigger fires shown as a distinct event class. |
| Email Simulator | `/email-sim` | **enhance** | Becomes the live driver for the **email trigger** — a received email can fire an agent. |
| Office / Dashboard / Mission v1 | `/office`, `/`, `/mission-v1` | keep | Legacy; not invested in. |

### New surfaces

#### 2.1 Agent Library & Builder — `/agents` **(new)**

**Purpose:** Browse, create, edit, clone, disable agents without touching code.

**Layout:** Left filter rail (group by pack: EPC / IT / Custom; filter by mode/permission) · center grid of agent cards · right slide-in editor panel.

**Agent card:** role id + display name, pack badge, mode chip (`oneshot`/`session`), permission summary (create/modify/execute, allowed paths), tool list, bound-trigger count, status (active/disabled). Actions: Edit · Clone · Disable.

**Create / Edit Agent form:**
- Identity: role id (kebab-case, unique), display name, pack assignment (EPC / IT / Custom / none).
- Prompt: full markdown editor (writes `agents/{id}/agent.md`).
- Execution: mode = `oneshot` | `session`.
- Permissions: `canCreateFiles`, `canModifyFiles`, `canExecuteCode`, allowed extensions, allowed paths, denied paths.
- Hierarchy: delegates-to (multi-select of existing agents), reviews-to (multi-select).
- Triggers: optionally attach default triggers (links to Triggers Console schema).
- Save → writes `agents/{id}/agent.md` + `agents/{id}/settings.json`, hot-reloads the registry, no restart.

**User actions:** create custom agent, edit any agent's prompt/permissions/hierarchy, clone a default agent as a starting point, disable an agent (kept on disk, excluded from rosters).

#### 2.2 Industry Packs — `/packs` **(new)**

**Purpose:** See and apply industry packs; compose a custom team per project.

**Layout:** Pack list (EPC, IT, + any custom) · selected-pack detail showing the roster, an auto-rendered delegation/review graph, the workflow phase sequence, and bundled default triggers.

**User actions:**
- "Apply pack to project" — sets a project's roster to the pack's.
- "Compose custom team" — pick agents across packs into a per-project roster; the delegation/review graph is derived from each agent's declared `delegatesTo`/`reviewsTo`.
- "Save as pack" — persist a composed roster as a new custom pack manifest.

#### 2.3 Triggers / Automations Console — `/triggers` **(new)**

**Purpose:** Define what events wake which agents, and watch them fire.

**Layout:** Trigger table (type, name, binding, status, last fired, fire count) · "New trigger" panel · per-trigger fire-history drawer.

**New-trigger panel fields:**
- Type: `email` | `webhook` | `cron` | `file_watch`.
- Match conditions (type-specific): email → sender/subject regex/category; webhook → path + optional HMAC secret; cron → schedule (`5m`/`1h`/`24h`/`@daily`); file_watch → directory + glob.
- Binding: target = a specific **agent** OR a **project/workflow** (Brain `create_project` vs `inject`).
- Safety: dedupe key, rate limit (max fires/window), max concurrent, enabled toggle.

**User actions:** create/edit/delete trigger, enable/disable (kill-switch), **Test-fire** (synthetic event to verify wiring), inspect fire history (payload hash, spawned task id, checkpoint outcome).

#### 2.4 Mission v2 enhancements **(enhance)**

- **Pack-aware dynamic sidebar** — the roster reflects the project's composed team (extends the recent "Dynamic sidebar" commit). Pack badge per agent.
- **Header team/pack indicator** — shows active pack ("EPC" / "IT" / "Custom (n agents)").
- **Event-tagged activity feed** — trigger-spawned work is visually distinct ("⚡ Triggered by: invoice file-watch") vs human-initiated.
- **Trigger inbox badge** — count of pending trigger-spawned checkpoints.

### Feature catalog (Phase 1)

1. On-demand custom agent creation (UI + folder).
2. Industry packs (EPC + IT defaults, custom packs).
3. Per-project team composition across packs.
4. Registry-backed roles (no recompile to add/change an agent).
5. Email trigger (inbound → agent).
6. Webhook trigger (generalized binding over existing `hooks.go`).
7. Cron trigger (generalized binding over existing `cron.go`).
8. File-watch trigger (fsnotify on project/inbox dirs).
9. Trigger → orchestrator dispatch with dedupe + rate-limit + kill-switch.
10. Checkpoint-gated reactive execution (reuses existing HITL).
11. Trigger fire audit log.
12. Pack-aware Mission v2 UI.

---

## Part III — Functional Specification

### 3.1 Pack manifest schema — `industries/{pack}/pack.json`

```json
{
  "id": "epc",
  "name": "Engineering, Procurement & Construction",
  "description": "Capital project delivery crew.",
  "agents": ["project_manager", "site_engineer", "procurement", "hse", "qa_inspector", "hr"],
  "hierarchy": {
    "delegation": { "project_manager": ["site_engineer", "procurement", "hse", "qa_inspector", "hr"], "site_engineer": ["qa_inspector"] },
    "review":     { "qa_inspector": ["site_engineer", "project_manager"], "hse": ["project_manager", "ceo"] }
  },
  "workflow": { "phases": ["triage", "research", "planning", "discussion", "execution", "qa"] },
  "defaultTriggers": [
    { "type": "cron", "schedule": "@daily", "bind": { "agent": "project_manager" }, "name": "Daily schedule check" }
  ]
}
```

Built-in packs (`epc`, `it`) ship in-repo. Custom packs are written by the UI to the same path.

### 3.2 Agent definition schema (extends existing files)

`agents/{id}/agent.md` — system prompt (unchanged convention).
`agents/{id}/settings.json` — extended:

```json
{
  "id": "procurement",
  "displayName": "Procurement Manager",
  "pack": "epc",
  "mode": "oneshot",
  "permissions": {
    "canCreateFiles": true, "canModifyFiles": true, "canExecuteCode": false,
    "allowedExtensions": [".md", ".txt", ".json", ".csv"],
    "allowedPaths": ["docs/", "procurement/", "vendors/", "rfq/"],
    "deniedPaths": [".git/", ".env"]
  },
  "delegatesTo": [],
  "reviewsTo": ["project_manager", "ceo"],
  "tools": [],
  "disabled": false
}
```

This is a superset of today's `RolePermissions`/`DefaultModes`/`PermissionModes`/hierarchy entries — those Go map literals become the *seed values* serialized into these files for the two built-in packs.

### 3.3 Runtime registry (the engineering spine)

**Problem:** roles are currently hardcoded package-level maps in `internal/agent/agent.go` (`PermissionModes`), `modes.go` (`DefaultModes`), `permissions.go` (`RolePermissions`), `hierarchy.go` (`DelegationHierarchy`, `ReviewHierarchy`). Adding an agent requires editing five maps and recompiling.

**Solution:** a `Registry` that loads `agents/*/settings.json` + `industries/*/pack.json` at startup and on demand, and exposes the same lookups via functions:

- `Registry.Mode(role)` replaces `DefaultModes[role]`
- `Registry.PermissionMode(role)` replaces `PermissionModes[role]`
- `Registry.Permission(role)` replaces `RolePermissions[role]`
- `Registry.DelegatesTo(role)` / `Registry.ReviewsTo(role)` replace the hierarchy maps

The existing map literals are kept only as the *generator* that writes the built-in pack JSON once (a `seed` step), then deleted as live lookups. Hot-reload: a `POST` that creates/edits an agent rewrites the folder and calls `Registry.Reload()`.

**Backward compatibility:** built-in packs must seed byte-identical behavior to today's maps — covered by a golden test (Epic E1 acceptance).

### 3.4 Trigger engine

```
Trigger sources ─┐
  email poll      │
  webhook (HTTP)  ├─▶ Event Bus (in-process pub/sub, glob topics) ─▶ Dispatcher ─▶ Orchestrator
  cron tick       │                                                     │             (InjectTask /
  file-watch      ┘                                                     │              Brain create_project)
                                                                  dedupe + rate-limit
                                                                  + kill-switch + audit
```

- **Event Bus:** in-process buffered channels, subscribe by topic glob (e.g., `email.*`, `file.invoices.*`). No external broker (filesystem/single-org constraint).
- **Sources:**
  - *Email* — polls the existing inbox / Email-Sim receive path; matches conditions; emits `email.received`.
  - *Webhook* — reuse `internal/web/hooks.go`; generalize so a hook binds to a trigger record instead of a hardcoded task template.
  - *Cron* — reuse `internal/web/cron.go`; generalize binding the same way.
  - *File-watch* — `fsnotify` on configured dirs (defaults: `projects/<id>/inbox/`, `projects/<id>/drawings/`); debounce; emit `file.<bucket>.created`.
- **Dispatcher:** resolves the trigger's binding → builds a task → calls `orchestrator.InjectTask` (existing) or Brain `create_project` (existing). Applies dedupe key, rate limit, max-concurrent. Records a fire entry.
- **Persistence (filesystem):** trigger definitions in `.triggers/triggers.json`; fire history in `.triggers/history.json`.
- **Safety:** every trigger-spawned task enters the **existing checkpoint pipeline** — dangerous/external actions pause for human approval exactly as today. Global kill-switch disables all triggers. Per-trigger enable toggle. Rate-limit + dedupe prevent storms/loops. Every fire is audited (trigger id, payload hash, spawned task id, checkpoint outcome) into the existing message store + `.triggers/history.json`.

### 3.5 Brain integration

Triggers may dispatch directly to an agent (fast path) or route the event text through the existing Brain router for classification (`delegate` / `create_project`). Phase 1 default: direct bind for `webhook`/`cron`/`file_watch`; Brain-routed for `email` (so subject/body classification reuses existing Brain logic).

---

## Part IV — Engineering Execution

### Epics

**E1 — Pack & Agent Registry** *(spine; blocks E2/E3)*
- New `internal/registry/` (loader, validation, hot-reload, `Reload()`).
- Refactor `internal/agent/{agent,modes,permissions,hierarchy}.go`: replace package-level map *lookups* with `Registry` calls; keep literals only as a one-time `seed` that writes `industries/{epc,it}/pack.json` + `agents/*/settings.json`.
- Create `industries/epc/pack.json`, `industries/it/pack.json`.
- **Acceptance:** golden test proving registry-backed lookups are byte-identical to the old maps for all 14 existing roles; a new agent folder is picked up via `Reload()` with no recompile; `go build ./...` green.

**E2 — Agent Builder + Packs API & UI** *(depends E1)*
- API: `GET/POST/PUT/DELETE /api/agents`, `GET/POST /api/packs`, `POST /api/projects/{id}/team` (compose).
- UI: `/agents` (Library & Builder), `/packs` (pack viewer/composer).
- Write-through to disk + `Registry.Reload()`.
- **Acceptance:** create a custom agent via UI → it appears in a project roster and can be delegated to, with no restart.

**E3 — Trigger Engine** *(depends E1; parallel with E2)*
- New `internal/trigger/` (event bus, sources, dispatcher, safety).
- File-watch via `fsnotify`; email source over existing inbox; generalize `hooks.go`/`cron.go` bindings.
- Dispatch through existing `orchestrator.InjectTask` / Brain `create_project`; dedupe + rate-limit + kill-switch + fire audit.
- **Acceptance:** each of the 4 trigger types fires end-to-end → spawns a checkpoint-gated task; disabling a trigger stops it; a forced loop is stopped by rate-limit/dedupe.

**E4 — Triggers Console UI** *(depends E3)*
- `/triggers`: table, new-trigger panel, test-fire, fire-history drawer.
- **Acceptance:** a trigger created, test-fired, and inspected entirely from the UI.

**E5 — Mission v2 enhancements** *(depends E1, E3)*
- Pack-aware sidebar (extends current dynamic sidebar), header pack indicator, event-tagged activity feed, trigger checkpoint badge.
- **Acceptance:** an EPC project shows the EPC roster; a trigger-spawned task is visually tagged in the feed.

**E6 — Safety & audit hardening** *(depends E3)*
- Checkpoint-gating integration tests for trigger-spawned dangerous actions; kill-switch; audit completeness.
- **Acceptance:** no trigger-spawned dangerous action bypasses a checkpoint; 100% of fires are in the audit log.

**E7 — Docs & test harness** *(continuous, per global workflow rule)*
- Update `api_test.html`, `uiflow.md`, and Swagger for every new/changed endpoint in E2/E3/E4.
- **Acceptance:** every Phase 1 endpoint present in all three.

### Sequencing & dependencies

```
E1 ──┬──▶ E2 ──▶ E4 ──┐
     └──▶ E3 ──┬──────┴──▶ E5 ──▶ E6
               └──▶ (E3 also feeds E5/E6)
E7 runs continuously alongside E2/E3/E4
```

E1 first (everything depends on the registry). E2 and E3 parallelizable after E1. E4/E5 after their feeders. E6 last. E7 tracks every endpoint change.

### Explicitly out of scope (YAGNI for Phase 1)

- PostgreSQL / multi-tenant / RBAC / SSO — stays filesystem, single-org.
- Generic tool catalog & MCP router (future.md Phase 1) — deferred to Phase 2; agents keep their current tool access.
- Agent **marketplace** with ratings/versioning/community/one-click install — Phase 1 is folder + form only.
- Autonomous agent-to-agent **Signal System** — deferred (highest risk; Phase 2+).
- Multi-model routing, mobile dashboard, on-prem packaging, compliance report templates.
- New autonomous tool execution beyond what agents already do — triggers only *initiate*; the existing pipeline executes under existing gates.

### Risks & mitigations

| Risk | Mitigation |
|---|---|
| Registry refactor changes behavior for the 14 existing roles | Golden test (E1) asserts byte-identical lookups before deleting map literals. |
| Trigger storms / feedback loops (agent action re-fires a trigger) | Dedupe key + rate-limit + max-concurrent + global kill-switch (E3/E6). |
| Reactive work bypasses human oversight | All trigger-spawned tasks routed through the *existing* checkpoint pipeline; E6 integration tests assert no bypass. |
| File-watch flakiness across OS | `fsnotify` + debounce; fall back to directory polling where unreliable. |
| Hot-reload races with a running pipeline | `Reload()` is copy-on-write; in-flight tasks keep their resolved roster snapshot. |

### Definition of done (Phase 1)

All E1–E7 acceptance criteria pass · `go build ./...` green · the five success metrics in Part I §6 met · `api_test.html`/`uiflow.md`/Swagger current · an EPC pilot project running with ≥ 2 live triggers under checkpoint gating.
