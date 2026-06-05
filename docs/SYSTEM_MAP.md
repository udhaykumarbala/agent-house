All key facts confirmed: module is `pty-claude-test`, port default 8080, default checkpoints all block (template/plan/phase-gate/final), `.env` has the three Anthropic keys, `claude` CLI resolves to `~/.local/bin/claude`, and `scripts/build-ui.sh` exists. Now I'll write the runbook.

```markdown
# Agent House v2 — Engineering Runbook

> Authoritative operations guide for the Go multi-agent orchestration system. Module name is `pty-claude-test` (legacy); the repo is `agent-house-2`. Entry point: `cmd/agent-house/main.go`. Everything is flat-JSON-on-disk; there is no database.

---

## 1. System Overview

Agent House v2 turns one plain-English request into a working webapp by chaining four layers. The **Brain** (`internal/brain`) is the natural-language front door: it takes a chat message, assembles workspace context, and makes a single stateless LLM call that emits a JSON action — either answer directly, `delegate` a one-shot task to a specialist, or `create_project` (launch a full build). A build runs through the **Orchestrator** (`internal/orchestrator`), the pipeline engine that drives a task through CEO triage → template selection → research → planning → discussion → development → QA, pausing at human-approval **checkpoints** between phases. Each phase is executed by one or more **Agents** (`internal/agent`) — role definitions (CEO, PM, Architect, Senior/Junior Dev, plus an EPC construction team) loaded from `agents/<role>/{agent.md,settings.json}`. Each agent runs either as a fast text-only **oneshot** API call or as a persistent **Session** (`internal/session`) — a long-lived `claude -p --output-format stream-json` child process with full tool access (Read/Write/Edit/Bash) whose NDJSON output is normalized into events streamed live over `/ws`. Phases never pass objects to each other; they coordinate **entirely through files** written under `projects/<id>/.plans/`, and generated app code lands in `projects/<id>/`.

```
                            ┌─────────────────────────────────────────┐
   user (NL)                │              Agent House v2               │
      │                     │                                           │
      ▼   POST /api/chat     │   ┌──────────┐  create_project          │
 ┌──────────┐  ───────────▶ │   │  BRAIN   │ ───────────────┐          │
 │  HTTP    │               │   │ (1 LLM   │  delegate       │          │
 │  layer   │  POST /api/task│   │  call)   │ ──────┐        │          │
 │ (net/http│  ───────────────▶ └──────────┘       │        ▼          │
 │  ServeMux│               │                       │  ┌──────────────┐ │
 │  + CORS) │◀── /ws stream │                       │  │ ORCHESTRATOR │ │
 └──────────┘   (events)    │                  InjectTask  │ processWith- │
      ▲                     │                       │  │   Phases     │ │
      │ checkpoint decide   │                       ▼  │              │ │
      └─────────────────────┼──────────┐  ┌─────────────┐ triage→     │ │
                            │  CHECKPOINT│  │ injected/   │ template→   │ │
                            │  (blocks   │  │ immediate   │ research→   │ │
                            │  pipeline) │  └─────────────┘ planning→   │ │
                            │            │                 discussion→ │ │
                            │            │                 dev+QA loop  │ │
                            │            └────────┬───────────┬────────┘ │
                            │                     │           │          │
                            │           per phase ▼           ▼ writes    │
                            │        ┌──────────────┐   ┌──────────────┐ │
                            │        │   AGENTS     │   │  projects/   │ │
                            │        │ (role defs)  │   │   <id>/      │ │
                            │        └──────┬───────┘   │  .plans/*    │ │
                            │      oneshot  │ session    │  <app code>  │ │
                            │       (API)   ▼           │  .tasks/*    │ │
                            │        ┌──────────────┐   └──────────────┘ │
                            │        │  SESSIONS    │   files = the only │
                            │        │ claude -p    │   inter-phase data │
                            │        │ NDJSON→events│──────▶ /ws         │
                            │        └──────────────┘                    │
                            └─────────────────────────────────────────┘
```

---

## 2. Runtime Prerequisites

**The `claude` CLI must be on the server process's PATH.** Sessions spawn via `exec.Command("claude", ...)` with no fallback path (`internal/session/session.go:60`). Here it resolves to `/Users/udhaykumar/.local/bin/claude`. If you launch the server from a minimal-PATH environment (a LaunchAgent, non-interactive ssh), every session spawn fails — this is the #1 runtime failure mode.

**Build & run (preferred):**
```bash
cd /Users/udhaykumar/susanoox/agent-house-2
go build -o /tmp/agent-house ./cmd/agent-house
/tmp/agent-house --serve --port=8080 --project=./projects
```

**Or run directly:**
```bash
go run ./cmd/agent-house --serve --port=8080 --project=./projects
```

- `--serve` is **required** to start the HTTP server. Without it the binary runs CLI single-agent/orchestrate modes (legacy, no sessions) and `--task` is required.
- `--port` default **8080**. `--project` default **`./projects`** (base dir under which `projects/<project_id>/` are created).
- `loadEnv(".env")` runs at boot and reads `KEY=VALUE` lines but does **NOT** override env vars already set in the process environment.

**Environment variables** (`.env` in repo root):

| Var | Required? | Purpose |
|-----|-----------|---------|
| `ANTHROPIC_API_KEY` | **Mandatory for the Brain** (`/api/chat`) and for oneshot agents. Without it `/api/chat` fails with "API client not configured"; oneshot-default roles silently fall back to slower CLI sessions. | Auth for the oneshot Anthropic-compatible `/v1/messages` call. |
| `ANTHROPIC_BASE_URL` | Optional (default `api.anthropic.com`) | Lets the Brain + sessions point at a compatible API (MiniMax, etc.). Inherited by the spawned `claude` CLI. |
| `ONESHOT_MODEL` | Optional (default `claude-sonnet-4-6`) | Model for the **oneshot HTTP path only**. Session/CLI mode never passes `--model`; the CLI uses its own configured default. |
| `OPENAI_API_KEY` | Optional | Upgrades conversation search from lexical substring to semantic (cosine, `text-embedding-3-small`). Search-only — NOT used for routing. |
| `OPENAI_EMBEDDING_MODEL` | Optional | Override embedding model. |
| `RESEND_API_KEY` | Optional | Enables real outbound email. Without it `POST /api/email/send` returns `success:false` with a draft note (does not error, does not persist). |

The current `.env` contains only `ANTHROPIC_API_KEY`, `ANTHROPIC_BASE_URL`, `ONESHOT_MODEL`.

**Embedded Next.js UI** (`/mission`, `/conductor`, `/triggers`, `/agents`, `/lab`, `/_next/*`) requires a one-time build, else those routes return HTTP 500 "UI not built":
```bash
bash scripts/build-ui.sh
```
Legacy HTML routes (`/`, `/office`, `/mission-legacy`, `/mission-v1`, `/live`, `/email-sim`, `/project/*`) are always available without the build.

**Smoke test:**
```bash
curl -s localhost:8080/api/status
```

---

## 3. Complete HTTP API Reference

Single `http.ServeMux` + CORS wrapper (no router framework). Trailing-slash routes are **prefix** matches; handlers manually parse path segments. Method dispatch is inside each handler (wrong method → 405). **CORS caveat:** `Access-Control-Allow-Methods` advertises only `GET, POST, OPTIONS` even though several handlers accept `PATCH`/`PUT`/`DELETE` — cross-origin preflighted PATCH/PUT/DELETE from a browser may be blocked.

### Brain / Conversations
| Method | Path | Purpose | Key request fields |
|--------|------|---------|--------------------|
| POST | `/api/chat` | Route an NL message through Brain router+executor; returns `{response, action, project_id, success, suggestions[]}` | body `{message (req), user_id (default "default")}` |
| GET | `/api/chat` | Full conversation history | query `?user=<id>` |
| GET | `/api/conversations` | List conversations w/ metadata (palette) | — |
| POST | `/api/conversations` | Mint a fresh conversation id | empty body → `{id:"conv_<ms>"}` |
| GET | `/api/conversations/search` | Semantic (if `OPENAI_API_KEY`) else lexical search | query `?q=` (limit 30) |
| GET | `/api/conversations/{id}` | Full thread for a conversation id | path id |

### Mission / Pipeline
| Method | Path | Purpose | Key request fields |
|--------|------|---------|--------------------|
| POST | `/api/task` | **Create + run the full build pipeline** (background goroutine); returns `task_id` | `{task (req), project_id (default "default"), continue_from?}` |
| POST | `/api/inject` | Inject a task into a running pipeline (drained between phases) or run immediately if idle | `{task (req), agent_role? (empty=auto→senior_dev), project_id?, priority? high\|normal}` |
| GET | `/api/phase` | Current high-level phase + ordered phase list with status | — |
| GET | `/api/phase-status` | Detailed dev-phase status (has_plan, current_phase, subtask counts, qa_status) | query `?project=` (req) |
| GET | `/api/development-plan` | `.plans/development-plan.json` | query `?project=` (req) |
| GET | `/api/subtasks` | SubTasks for a dev-plan phase index | query `?project=` (req), `?phase=<int>` (req) |
| GET | `/api/qa-reviews` | QA review history | query `?project=` (req) |
| GET | `/api/status` | Server health: task_running, running_projects[], message_count, ws_clients, project_dir | — |
| GET | `/api/messages` | Message store; filter `?project=` / `?agent=` / `?after=<RFC3339>` | query |

### Checkpoints / Settings
| Method | Path | Purpose | Key request fields |
|--------|------|---------|--------------------|
| GET | `/api/checkpoints` | List checkpoints + `pending_count` (>0 = pipeline blocked) | query `?project=` (default "default") |
| GET | `/api/checkpoints/{id}` | One checkpoint + `artifact_content` (reads ArtifactPath) | path id; query `?project=` |
| POST | `/api/checkpoints/{id}/decide` | **Resolve a blocking checkpoint; unblocks the pipeline** | path id; query `?project=`; body `{action: approved\|rejected\|overridden (req), feedback?, override_data?}` |
| GET/PUT | `/api/settings/workflow` | Get/set WorkflowSettings (which gates block + auto_delegate_minutes) | query `?project=`; PUT body WorkflowSettings |

### Projects / Tasks / Files
| Method | Path | Purpose | Key request fields |
|--------|------|---------|--------------------|
| GET | `/api/projects` | List project dirs (id, name, task_count, last_updated) | — |
| GET | `/api/projects/{id}` | Full detail: meta, files[], tasks, session metrics, dev_plan | path id |
| GET | `/api/files` | Recursive non-hidden file list | query `?project=` |
| GET | `/api/file-content` | Raw file content (traversal-guarded) | query `?path=` (req) |
| GET | `/api/tasks` | Task history (`.tasks/history.json`) | query `?project=` |
| GET | `/api/task/{taskId}` | One task metadata + full conversation | path taskId; query `?project=` |
| GET | `/api/agent-tasks` | Fine-grained per-agent subtasks (READ-ONLY) | query `?project=`, `?role=`, `?status=` |
| GET | `/api/agent-tasks/{taskID}` | Single agent subtask | path taskID |

### Agents
| Method | Path | Purpose | Key request fields |
|--------|------|---------|--------------------|
| GET | `/api/agents` | All known roles (IT + EPC) with permissions, color, active flag | — |
| GET | `/api/agents/{role}/tasks` | Tasks for a role | path role |
| GET | `/api/agents/{role}/status` | Aggregate status for a role | path role |
| GET | `/api/agents/{role}/outputs` | All outputs across tasks for a role | path role |
| GET/POST | `/api/agents/{role}/chat` | GET last 50 idle-chat msgs; POST sends a message, spawns read-only response | path role; query `?project=`; POST `{content (req), context?}` (software roles only via NewAgent) |
| GET/POST | `/api/agents/modes` | GET per-role modes; POST set a role's mode | POST `{role, mode: oneshot\|session}` |

### Sessions / Reports
| Method | Path | Purpose | Key request fields |
|--------|------|---------|--------------------|
| GET | `/api/sessions` | List ALL Claude Code sessions across projects | — |
| GET | `/api/sessions/{projectID}` | Sessions for one project | path projectID |
| GET | `/api/sessions/{projectID}/{agentRole}` | One session's info/stats | path projectID, agentRole |
| POST | `/api/sessions/{projectID}/{agentRole}/{action}` | `abort` (SIGINT), `approve`, `deny` | path action; approve/deny body `{tool_use_id}` |
| GET | `/api/reports/{projectID}` | Aggregated per-project metrics (cost, tokens, per-agent) | path projectID |

### Kanban
| Method | Path | Purpose | Key request fields |
|--------|------|---------|--------------------|
| GET/POST | `/api/kanban/tasks` | GET list + column counts (filters status/agent/phase/priority); POST create | query `?project=`; POST kanban.Task (title req) |
| GET/PATCH | `/api/kanban/tasks/{id}` | GET one; PATCH partial update | path id; query `?project=`; PATCH `{field:value}` |
| POST | `/api/kanban/tasks/{id}/review` | QA review decision | `{action, feedback, criteria_updates:[{index,checked}]}` |
| POST | `/api/kanban/sync` | Generate/sync kanban from dev plan | query `?project=` |

### Capability (`/api/cap/*`) — per-`?scope=` JSON under top-level `data/`
| Method | Path | Purpose | Key request fields |
|--------|------|---------|--------------------|
| GET | `/api/capabilities` | Static descriptor of each agent's capabilities + endpoints | — |
| GET/POST | `/api/cap/hr/applicants` | List / store applicants | query `?scope=`; POST Applicant |
| PATCH/POST | `/api/cap/hr/applicants/{id}` | Update applicant status | path id; `{status}` |
| POST | `/api/cap/hr/match` | Match top-5 applicants vs JD | body JD `{title, must_have[], nice_to_have[], min_years}` |
| GET/POST | `/api/cap/procurement/vendors` | List / store vendors | POST Vendor (id req) |
| POST | `/api/cap/procurement/validate-invoice` | BEC/impersonation invoice check | `{sender_email (req), vendor_id, amount}` |
| GET/POST | `/api/cap/schedule/milestones` | List / store milestones | POST Milestone (title+due_date req) |
| GET | `/api/cap/schedule/slips` | Past-due milestones bucketed by severity | query `?scope=` |
| GET/POST | `/api/cap/email/inbox` | **GLOBAL** inbox (not scoped) list / inject | POST Email |
| GET | `/api/cap/email/summary` | Inbox triage stats (GLOBAL) | — |

### Scenario
| Method | Path | Purpose | Key request fields |
|--------|------|---------|--------------------|
| GET | `/api/scenarios` | List 6 registered runners | — |
| POST | `/api/scenario/{name}/run` | Run a capability-composing scenario; emits step messages to store+WS | path name; query `?scope=`; body = scenario Example map |

### Email (engine, rooted at ProjectDir)
| Method | Path | Purpose | Key request fields |
|--------|------|---------|--------------------|
| GET | `/api/email/inbox` (and `/{id}`) | List inbox; `/{id}` returns one (marks read) | optional path id |
| POST | `/api/email/receive` | Ingest inbound mail; trust + BEC scoring; returns trust+alert | `{from (req), from_name, to?, subject (req), body}` |
| POST | `/api/email/send` | Send via Resend (needs `RESEND_API_KEY`) | `{from?, to (req), subject (req), body}` |
| GET | `/api/email/applicants` | List parsed applicants | — |
| POST | `/api/email/trust` | Add address to vendor's trusted list | `{vendor_id, email}` |

### Triggers / Webhooks / Misc
| Method | Path | Purpose | Key request fields |
|--------|------|---------|--------------------|
| GET/POST/DELETE | `/api/cron` (and `/api/cron/`) | List / create / delete recurring agent triggers (min interval 30s, in-memory) | POST `{schedule e.g."5m" (req), task (req), agent_role? (default senior_dev), project_id?, id?}`; DELETE `?id=` |
| POST | `/api/hooks/{hook_id}` | Webhook → injected task via per-hook template; optional HMAC-SHA256 | path hook_id; arbitrary JSON; `?project`; header `X-Signature-256` |
| GET (Upgrade) | `/ws` | WebSocket: broadcasts `agent_event`, `message`, `checkpoint`, `brain_event`; inbound `approve_agent_tool`/`deny_agent_tool`/`abort_agent` | WS msg `{type, project_id, agent_role, tool_use_id}` |
| GET | `/` and static | Serves UI (legacy dashboard `/`, `/office`, `/email-sim`; Next.js `/mission` etc. need build) | — |

---

## 4. THE BUILD-A-WEBAPP RUNBOOK

This is the core operational flow. Use **`POST /api/task`** (not the Brain) for deterministic control — it drives `s.orchestrator`, which is the same instance the checkpoint-decide endpoint resolves against. (The Brain's `create_project` spins up a *separate* orchestrator whose checkpoints `/api/checkpoints/.../decide` may not resolve.)

> **Critical gotcha up front:** out of the box, 4 checkpoints BLOCK the pipeline — `template_approval` (after template selection), `plan_approval` (after discussion), `phase_gate` (after each dev phase passes QA), and `final_acceptance` (confirmed at `internal/checkpoint/checkpoint.go:77`, all `true`). A mission **stalls** at each one waiting for a human `decide` call. Choose attended (approve each) or unattended (disable gates first).

### Step 0 (optional) — Disable gates for a fully autonomous run
Do this *before* starting the task, so the run never blocks:
```bash
curl -sX PUT 'localhost:8080/api/settings/workflow?project=todo' \
  -H 'Content-Type: application/json' \
  -d '{"project_id":"todo","require_template_approval":false,"require_plan_approval":false,"require_phase_gate":false,"require_final_acceptance":false}'
```
(Alternative: set `auto_delegate_minutes` > 0 to let the CEO agent auto-approve after a timeout.)

### Step 1 — Create + start the pipeline
There is no separate "create project" endpoint — the project directory is created lazily by this call (`os.MkdirAll(./projects/todo)` at `server.go:426`).
```bash
curl -sX POST localhost:8080/api/task \
  -H 'Content-Type: application/json' \
  -d '{"task":"Build a todo webapp with add/complete/delete and localStorage","project_id":"todo"}'
```
**Expected response (immediate; pipeline runs in background):**
```json
{"success":true,"task_id":"task_1730000000000000000","project_id":"todo"}
```
If a task is already running for `todo` you get `{"success":false,"error":"Task already running for project todo"}` with HTTP 200 (one in-flight mission per project_id).

### Step 2 — Watch progress
Best: open the WebSocket and read lifecycle events (`phase_started`, `phase_completed`, `checkpoint_reached`, `task_completed`):
```bash
# any WS client, e.g.:
websocat ws://localhost:8080/ws
```
Or poll:
```bash
curl -s localhost:8080/api/status
curl -s 'localhost:8080/api/phase-status?project=todo'
curl -s 'localhost:8080/api/checkpoints?project=todo'
```

### Step 3 — Approve checkpoints (if gates are enabled)
A checkpoint is pending when `GET /api/checkpoints?project=todo` returns `pending_count > 0` (or you saw a `checkpoint_reached` WS frame with a `checkpoint_id`). Inspect the artifact before deciding:
```bash
curl -s 'localhost:8080/api/checkpoints/<cpID>?project=todo'   # includes artifact_content (template.md / approved-plan.md)
```
Approve to advance:
```bash
curl -sX POST 'localhost:8080/api/checkpoints/<cpID>/decide?project=todo' \
  -H 'Content-Type: application/json' \
  -d '{"action":"approved","feedback":""}'
```
`action` must be exactly `approved` | `rejected` | `overridden`. Reject semantics differ by gate: **template reject aborts the task**; **plan reject re-runs the Discussion phase** with your feedback; **`overridden` + `override_data`** overwrites `approved-plan.md` directly. You will hit, in order: `template_approval` → `plan_approval` → one `phase_gate` per development phase → `final_acceptance`. (The CEO triage can skip phases, so the exact set varies; `plan_approval` only fires if Discussion ran.)

### Step 4 — Know when it's done and where files land
**Done** = a `task_completed` WS lifecycle event (with files list), or `GET /api/tasks?project=todo` shows the task `completed`, or `GET /api/phase-status?project=todo` returns `is_complete:true`. A task that fails QA 3× on any dev phase ends `failed` ("phase N failed after 3 iterations").

**Files on disk** (`--project` base, default `./projects`):
```
projects/todo/                        ← generated webapp code (the deliverable)
projects/todo/.plans/
    template.md                       ← chosen template (gated by template_approval)
    research/*.md                     ← research-phase agent output
    specs/{product-spec,ux-spec,ui-spec}.md
    final/approved-plan.md            ← discussion output (gated by plan_approval)
    development-plan.json             ← dev phases/subtasks the dev loop executes
    qa-reviews.json                   ← QA verdict audit trail
projects/todo/.tasks/
    history.json                      ← task index (/api/tasks)
    <taskID>/{metadata.json,conversation.json}
    kanban.json, checkpoints.json, decisions.json, settings.json
    agent-tasks/<role>/*.json
projects/todo/project.json            ← phase/cost/turns/team meta
```
Inspect via API: `GET /api/projects/todo` (file list excludes `.plans`/`node_modules`), `GET /api/files?project=todo`, `GET /api/file-content?path=todo/index.html`.

### Iterating after a run
- **Inject a change** (drains between phases if running, else runs immediately):
  ```bash
  curl -sX POST localhost:8080/api/inject \
    -d '{"task":"Add dark mode toggle","agent_role":"senior_dev","project_id":"todo","priority":"high"}'
  ```
- **Continue from a finished task** (prepends prior context):
  ```bash
  curl -sX POST localhost:8080/api/task \
    -d '{"task":"Add auth","project_id":"todo","continue_from":"<prev_task_id>"}'
  ```

### Brain alternative (less deterministic)
```bash
curl -sX POST localhost:8080/api/chat -d '{"message":"Build me a todo app called todo","user_id":"default"}'
```
Functionally launches the same pipeline if the Brain classifies intent as `create_project`, but in its own orchestrator instance — prefer `/api/task` when you need checkpoint control.

---

## 5. Company Structure Setup (Software + EPC)

There are **three independent layers** that combine to define a "company," and there is **no single "activate company" API**.

### Layer A — Agent roster (`agents/<role>/`)
The live company = the union of role directories under the **top-level `agents/`** dir, auto-discovered at boot by `NewRegistry("agents")` (`server.go:96`). Each role needs **both** `agent.md` (entire file = system prompt) **and** `settings.json` (execution config); missing either → skipped at startup (log line, no crash). Dirs prefixed `_` or `.` are skipped. Adding a role dir requires a **server restart**.

- **Software team:** `ceo, pm, ux, ui, security, architect, senior_dev, junior_dev` (these 8 also have a hardcoded fallback in `agent.go` + `prompts/<role>.md`, so they're chat-able and pipeline-runnable).
- **EPC/construction team:** `hr, project_manager, procurement, site_engineer, hse, qa_inspector` (registry-only — `agent.NewAgent("site_engineer")` returns "unknown role"). **Caveat:** `isValidRole` (`agent.go:520`) recognizes only the 8 software roles, so EPC inter-agent `DELEGATE:`/`REVIEW:` text signals are silently dropped, and EPC roles are **not chat-able** via `/api/agents/{role}/chat`. EPC roles are pipeline-runnable but not first-class.
- `agents/demos/{epc,hospital,hr}/*` are **template** team configs (different role names like `project_director`) that are **NOT auto-loaded**. To use them you'd copy into a project's `.agent-house/agents/`. `team.json` is documented in `agents/demos/README.md` but **no Go loader reads it** — pipeline composition comes from CEO triage + `project.json`, not `team.json`.

**To run a "software & EPC" company today:** keep both role sets present under top-level `agents/` (they already coexist). Which roles actually run is decided dynamically by CEO triage + the phase config + the scope of work you submit — software build tasks invoke the IT pack; EPC scenarios/seed data exercise the construction roles.

### Layer B — Capabilities + seed data (`data/<scope>/`, LLM-free)
Capabilities are deterministic business logic over per-`?scope=` JSON files (default scope `"default"`). `scope` = tenant directory. Seed them via the `/api/cap/*` endpoints (HR applicants, procurement vendors, schedule milestones) or the inbox (global, `projects/inbox/`).

### Layer C — Scenarios (composed flows, LLM-free)
6 runners (`GET /api/scenarios`): `process_applicants`, `validate_invoice`, `schedule_check`, `process_inbox`, `morning_briefing`, `route_rfi`. `POST /api/scenario/{name}/run?scope=<tenant>` executes a flow over already-seeded data, emitting per-step agent-attributed messages to the store + `/ws`.

### Assembling a company end-to-end — use the seed scripts
"Setup" = running a re-runnable shell script that drives the same `/api/*` endpoints (the project's documented rule: talk to `/api`, never side-load JSON).

**EPC company:**
```bash
HOST=http://localhost:8080 bash scripts/scenarios/load-epc-atlas.sh
```
Seeds scope `atlas-site`: 5 inbox emails (impersonation, RFI, invoice, client progress, applicant), 2 cron jobs, 3 applicants, 3 vendors, 5 milestones → `data/atlas-site/{applicants,vendors,milestones}.json`. One vendor (`BetaElectrics`) is inactive and one milestone is past-due so the contract/slip logic fires.

**Software company:**
```bash
HOST=http://localhost:8080 bash scripts/scenarios/load-dev-calculator.sh
```
POSTs a build task to `/api/task` (IT-pack agents stream into `/conductor`) + 1 hourly status cron.

> **PORT MISMATCH WARNING:** the seed scripts default `HOST` to `http://localhost:8099`, but the server defaults to `--port 8080`. Always export `HOST=http://localhost:8080` (or run the server on 8099), or the curls fail to connect.

After seeding, exercise the company:
```bash
curl -sX POST 'localhost:8080/api/scenario/validate_invoice/run?scope=atlas-site' \
  -d '{"sender_email":"ahmed.r@gmail.com","vendor_id":"vendor_xyz","amount":340000}'
curl -s 'localhost:8080/api/cap/schedule/slips?scope=atlas-site'
```

---

## 6. How to Test the Generated Webapp

Generated code lands in `./projects/<project_id>/` (the project dir itself; `.plans/` and `.tasks/` are metadata, excluded from `GET /api/projects/{id}` file listings). The template (`internal/template/registry.go`) determines tech stack and required files — default is **`static-enhanced`**; others are `static-html`, `nextjs-frontend`, `go-api`, `fullstack`. Templates mandate `api_test.html` and `uiflow.md` artifacts.

**Locate & inspect:**
```bash
ls projects/todo
curl -s 'localhost:8080/api/projects/todo' | jq '.files'
curl -s 'localhost:8080/api/file-content?path=todo/index.html'
```

**Serve & test by template type:**
- **static-html / static-enhanced** — open `projects/todo/index.html` directly, or:
  ```bash
  python3 -m http.server 5500 --directory projects/todo
  # → http://localhost:5500
  ```
  If an `api_test.html` was generated, open it to exercise endpoints.
- **nextjs-frontend / fullstack** —
  ```bash
  npm install --prefix projects/todo
  npm run dev --prefix projects/todo
  ```
- **go-api** —
  ```bash
  go run .   # run from inside projects/todo (cd there in one compound command)
  ```

Cross-check the spec the build targeted in `projects/todo/.plans/specs/ui-spec.md` and `final/approved-plan.md`, and the QA verdicts in `.plans/qa-reviews.json`.

---

## 7. Known Gotchas & Likely Failure Modes (ranked by likelihood)

1. **`claude` CLI not on PATH (most likely).** Sessions `exec.Command("claude", ...)` with no fallback (`session.go:60`). A minimal-PATH launch (LaunchAgent, non-interactive ssh) makes every session spawn fail with "session create failed". Verify: `which claude` → `/Users/udhaykumar/.local/bin/claude`. Launch the server from a shell where that dir is on PATH.

2. **Pipeline hangs on a checkpoint.** By default `template_approval`, `plan_approval`, `phase_gate`, `final_acceptance` all BLOCK the pipeline goroutine indefinitely on `o.checkpointCh`. If nobody calls `/api/checkpoints/{id}/decide`, the build never finishes. Mitigate: `PUT /api/settings/workflow` to disable gates, set `auto_delegate_minutes>0`, or poll + decide each. Only ONE pending checkpoint per orchestrator at a time.

3. **Port mismatch on seed scripts.** Server defaults to 8080; `scripts/scenarios/*.sh` and the scripts README default `HOST` to 8099. Export `HOST=http://localhost:8080` first.

4. **`ANTHROPIC_API_KEY` missing.** `/api/chat` (Brain) fails outright ("API client not configured"). Oneshot-default agents silently degrade to slower CLI sessions instead of erroring — easy to miss.

5. **Brain `create_project` ≠ `/api/task` for checkpoint control.** The Brain launches its OWN orchestrator instance; `/api/checkpoints/.../decide` operates on `s.orchestrator` and may not resolve a Brain-launched run's checkpoint. Use `POST /api/task` when you need to drive gates.

6. **"Task already running" returns HTTP 200, not an error.** One in-flight mission per `project_id` (`s.projectTasks` guard). The `success:false` body is easy to overlook. Each `POST /api/task` also REBUILDS `s.orchestrator` (preserving only the SessionManager) and re-wires callbacks — mutating shared state without a lock, a potential race if two projects start near-simultaneously.

7. **Project id defaults to `"default"` everywhere `?project=` is omitted.** Easy to read/write the wrong project. There is no "create project" endpoint — dirs are created lazily by `POST /api/task`.

8. **Phases communicate only through `.plans/*` files; Research/Planning/Development errors are swallowed.** If an agent fails to write `specs/ui-spec.md`, the next phase proceeds with missing input and no error. File tracking is dual (parsed Write/Edit tool events + filesystem snapshot diff) and skips `node_modules`, `dist`, `.git`, hidden dirs — files written there aren't tracked.

9. **QA parsing is brittle.** `parseQADecision` keys on literal `QA_APPROVED:`/`QA_REJECTED:` substrings in the CEO's free text. No marker → treated as not-approved → another iteration; 3 failures → the whole task errors out. Triage and template selection silently degrade to defaults on any parse failure.

10. **Embedded Next.js routes 500 until built.** Run `scripts/build-ui.sh` for `/mission`, `/conductor`, `/triggers`, `/agents`, `/lab`. Legacy HTML routes work regardless.

11. **WebSocket/event delivery is best-effort and lossy.** Every channel hop uses non-blocking `select{ default: drop }` (buffers 512/1024/256). A slow dashboard silently drops events and slow WS clients are evicted — don't treat the stream as a complete audit log (use `.plans/`, `.tasks/`, `/api/messages` for that).

12. **`settings.json` fields are largely inert.** `tools.builtin`, `tools.mcp_servers`, `claude_permission`, `model`, `max_turns`, `timeout` are parsed but NOT applied to sessions. All roles spawn `bypassPermissions` (hardcoded `agent.go:48`), so `/approve` and `/deny` (HTTP and WS) are effectively no-ops — only `abort` (SIGINT) is a meaningful live control. Session `--max-turns` is hardcoded to 50; sessions never get `--model`.

13. **Cron is in-memory only.** Jobs from `POST /api/cron` are lost on restart (min interval 30s; `LoadCronJobsFromConfig` exists but isn't called at Start). `RESEND_API_KEY` unset → `/api/email/send` returns `success:false` (no actual draft persisted). The capability inbox (`/api/cap/email/*`) is a single GLOBAL mailbox at `projects/inbox/`, separate from and ignoring `?scope=`.

14. **Module name is `pty-claude-test`** (not `agent-house`) — relevant when grepping imports (`pty-claude-test/internal/...`).
```
