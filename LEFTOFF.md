# LEFTOFF — Agent House continuation notes

Handoff for resuming work. Last updated 2026-06-08.

- **Branch:** `epc-software-orchestration` (pushed to `origin`)
- **HEAD:** `db795a5` — "Software company: run modes, two-door switcher, autonomous build pipeline"
- **PR:** not opened yet — https://github.com/udhaykumarbala/agent-house/pull/new/epc-software-orchestration

---

## TL;DR — what works right now

Agent House runs **two companies behind one app**, switchable from the top bar:

1. **Software Studio** (`scope=software`) — an agent-driven SDLC that builds real apps
   (Template → Research → Planning → PRD → Development → QA), with selectable
   **run modes** governing checkpoint autonomy.
2. **Atlas EPC** (`scope=atlas-site`) — the construction company (inbox triage, BEC
   detection, vendors, schedule, HR) driven by deterministic scenarios.

**Proven this session:** a Blitz build ran end-to-end unattended (6 gates auto-approved)
and produced a real React+TS+Vite+Tailwind tip calculator — 24/24 unit tests pass,
`tsc` + `vite build` clean. App lives at `projects/tip-calculator/`.

---

## How to run

```bash
# Go toolchain for this repo:
export PATH="$HOME/.local/bin:$HOME/.gvm/gos/go1.24.11/bin:$PATH"

# 1) Build the embedded UI (Next static export → internal/web/webui/dist)
bash scripts/build-ui.sh

# 2) Build + run the server (UI is embedded in the binary)
go build -o /tmp/agent-house ./cmd/agent-house
/tmp/agent-house --serve --port=8080 --project=./projects

# 3) Open the app
open http://localhost:8080/conductor          # Software Studio is the default door
# switch company in the top-left, or force it:
#   http://localhost:8080/conductor?company=software
#   http://localhost:8080/conductor?company=epc
```

Secrets live in `.env` (gitignored): `GEMINI_KEY` (Brain via Gemini native),
`ANTHROPIC_API_KEY`/`ANTHROPIC_BASE_URL` (MiniMax — the per-agent `claude` CLI
sessions that actually write code). The Brain uses `ONESHOT_MODEL=gemini-3.5-flash`
+ `ONESHOT_ROUTER_MODEL=gemini-3.1-flash-lite` (tiered).

### Run modes (build autonomy)

Set the **global default** new builds inherit:
```bash
curl -X PUT 'localhost:8080/api/settings/workflow?project=_global' \
  -H 'Content-Type: application/json' \
  -d '{"run_mode":"full_auto","decision_timeout_minutes":5}'
# run_mode ∈ manual | semi_auto | full_auto | blitz
```
or use the **"Build autonomy" control** in the software Conductor header.

| Mode | Gates | On a gate |
|------|-------|-----------|
| manual | all (incl. research/spec) | wait for human, forever |
| semi_auto (default) | PRD + Design + phase + final | wait for human |
| full_auto | PRD + Design + phase + final | live countdown X min → CEO auto-decides w/ reasoning |
| blitz | none effectively | auto-approve instantly (logged) |

> ⚠️ **Global mode is currently `blitz`** (set for the demo). Reset to `semi_auto`
> for normal interactive use.

### Drive a build
Software Conductor → *"Build a &lt;thing&gt;"* → routes to `create_project`, spawns the
orchestrator on its own goroutine, runs the SDLC. Watch it on **Mission → Software
Studio** (live countdown for full-auto, `⚡/🤖 auto-approved …` in the feed).

### Reseed the EPC demo
```bash
curl -X POST 'localhost:8080/api/scenario/reseed?scope=atlas-site'
# (shells out to scripts/scenarios/load-epc-atlas.sh — resets inbox + stores)
```

---

## Architecture map (where things live)

**Run-mode engine**
- `internal/checkpoint/checkpoint.go` — `WorkflowSettings.RunMode` + `DecisionTimeoutMinutes`,
  `ApplyRunModePreset()`, `RunMode*` consts, back-compat in `LoadSettings`.
- `internal/orchestrator/orchestrator.go` — `waitForCheckpoint()` blitz short-circuit
  (~line 350) + `autoDelegateTimer()` (full-auto CEO auto-decide + feed broadcast ~line 528).
- `internal/web/checkpoint_handlers.go` — PUT `/api/settings/workflow` applies the preset;
  `handleAllCheckpoints` = GET `/api/checkpoints/all` (cross-project aggregation w/ run_mode+timeout).
- `internal/web/brain_handlers.go` `handleCreateProject` — new builds inherit `_global` mode.

**Two-door / company model**
- `frontend/lib/company.ts` — company config (roster, scope, labels, role display names),
  `getCompany()`/`setCompany()`.
- `frontend/components/Chrome.tsx` — top-bar company switcher dropdown.
- `frontend/components/ConductorApp.tsx` — reads company → roster filter, scope to
  `sendChat`, `RunModeControl` in header (software only).
- `frontend/components/MissionApp.tsx` — company-aware roster/labels, `fetchAllCheckpoints`,
  1s countdown tick, company-filtered feed.
- `frontend/components/RunModeControl.tsx` — the 4-mode selector + timeout input.
- `internal/brain/router.go` `Route()` — company-aware system-prompt banner (software vs EPC).

**Conversation naming** — `internal/brain/conversation.go` `intentTitle()` ("Build:"/"Modify:").

**Brain parser hardening** — `internal/brain/router.go` `parseDecision` Try-4 salvage +
`extractJSONStringField`; client unwrap in `frontend/lib/md.ts` `cleanReply`.

**EPC** — `internal/scenario/*` (process_inbox idempotent applicant store, validate_invoice,
route_rfi), `internal/capability/*` (procurement BEC + vendor domain, hr, schedule, inbox),
`scripts/scenarios/load-epc-atlas.sh` (seed, with a reset preamble).

**Design docs** — `docs/SOFTWARE_COMPANY_MODES.md`, `docs/SYSTEM_MAP.md`,
`docs/EPC_*` (demo plan, API stability, conversation audit).

---

## Loose ends / next steps (prioritized)

1. **🔴 Checkpoint resolve routing for human-approved software builds.**
   Each build runs on its own `newOrch` (in `handleCreateProject`) with its own
   `checkpointCh`. The HTTP decide handler (`handleCheckpointDecide`) calls
   `s.orchestrator.ResolveCheckpoint` — the **main** orchestrator's channel — so a
   human Approve/Reject on a software build's gate in **manual/semi-auto** likely
   never unblocks the build. Blitz + full-auto resolve *internally* (so they work).
   **Fix before shipping semi/manual for software:** route the decision to the
   build's orchestrator (registry of orchestrators by project, or share one
   checkpointCh, or disk-poll resolution).

2. **🟡 EPC `atlas-emails.json` manifest is missing.** `load-epc-atlas.sh` has a
   data-driven loop that reads `scripts/scenarios/atlas-emails.json` for 10 new EPC
   emails (lookalike-domain BEC, multi-vendor invoices, RFIs, etc.), but the file
   was never written (pivoted to software mid-task). The loop is a guarded no-op.
   The authored content is in the workflow output from that session
   (`enrich-atlas-epc-data` run) — reconstruct the JSON or re-run the workflow.
   The new vendors/applicants/milestones in the seed script DO apply.

3. **🟡 Reset global run mode to `semi_auto`** (currently `blitz`).

4. **🟢 Software Conductor "Ask me" suggestions** are still EPC-flavored defaults
   ("Pause all triggers", "Email me a daily summary"). Make them company-aware in
   `ConductorApp` / the suggestions source.

5. **🟢 Dev-iteration 10-min agent timeout** — each dev/QA iteration can hit a 10-min
   session cap (`Phase completed in 10m0.003s`), making big builds slow / forcing
   QA iterations. Consider raising `MaxTurns`/session timeout or splitting subtasks.

6. **🟢 Subtask-completion tracking quirk** — `development-plan.json` subtasks can stay
   `pending` even after a phase is QA-approved (phase advances, subtask flags don't).
   Cosmetic, but misleads the Mission/dev-plan view.

7. **🟢 `projects/` runtime is dirty** (conversation logs, inbox, demo apps incl.
   `tip-calculator`). Regenerable — intentionally left uncommitted. Decide whether to
   keep `tip-calculator` as an example or `git clean` it.

---

## Verify (quick re-checks)

```bash
# generated app still builds + tests pass
cd projects/tip-calculator && npm test -- --run && npm run build

# run-mode presets apply
curl -s 'localhost:8080/api/settings/workflow?project=_global'
curl -X PUT 'localhost:8080/api/settings/workflow?project=_global' -d '{"run_mode":"blitz"}'

# company-aware brain (software scope must NOT run EPC scenarios / fabricate inbox)
curl -s -X POST localhost:8080/api/chat -H 'Content-Type: application/json' \
  -d '{"message":"process my inbox","user_id":"t","scope":"software"}'

# brain salvage unit tests
go test ./internal/brain/ -run 'Salvage|ExtractJSONString' -v
```

## Gotchas
- Go PATH: `export PATH="$HOME/.local/bin:$HOME/.gvm/gos/go1.24.11/bin:$PATH"`.
- Commit messages: **no Claude co-author** (per user global instruction).
- Never commit `.env` (gitignored) or `node_modules/` (gitignored).
- The UI is embedded — after a frontend change you must `build-ui.sh` **and** rebuild
  the Go binary for it to take effect.
</content>
