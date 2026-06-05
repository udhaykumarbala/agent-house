I have enough verified evidence to synthesize the plan. The reviews' file:line references check out against the actual source. Producing the final build plan now.

# EPC Demo — Tech Lead Build Plan (ship a fully functional, dynamic demo TODAY)

The Brain is already the strongest asset and demo-ready. The work today is: (1) un-break the two centerpiece moments (BEC catch + proactive fires), (2) wire the chat lane into the multi-agent scenario engine so a typed prompt visibly fans work out, (3) make scenarios react to scope + live state, and (4) tune the Brain to delegate proactively and feel snappy. Everything below is ranked P0 (must-have for the demo) / P1 / P2.

---

## 1. Demo narrative (this becomes `scenario.md`)

Run on `scope=atlas-site`. **Pre-demo ritual (5 min, P0):** `HOST=http://localhost:8080 bash scripts/scenarios/load-epc-atlas.sh` then confirm `curl -s 'http://localhost:8080/api/cap/email/summary?scope=atlas-site' | jq .impersonations` returns `1` and `curl -s http://localhost:8080/api/cron/ | jq '.jobs|length'` returns `2`.

Persona: **Director of an EPC firm** sitting down at 8am. Open on `/conductor`.

| # | Director types | Expected response | Proactive "wow" beat |
|---|---|---|---|
| 0 | *(nothing — just opens /conductor)* | Morning briefing auto-loads (run_scenario `morning_briefing` on first load) | **A "fire" card appears unprompted ~20-30s in:** "New email flagged: possible impersonation of XYZ Steel — review?" The Director didn't ask. This is the single line that separates Agent House from ChatGPT. |
| 1 | `Good morning. What do I need to handle today?` | Prioritized brief: CRITICAL BEC alert, contract penalty math ($11,900–$23,800), client deliverable, RFI #18, applicant Raj Kumar. 4 action chips. | Chips are real entities, not generic. "This isn't a chatbot, it read my inbox." |
| 2 | `Process my inbox` | Brain returns `action: run_scenario` → process_inbox fans out: **procurement, hr, project_manager agents each pulse in the workforce panel**, per-email triage table. | The conductor visibly delegates to 4 disciplines. Multi-agent value made tangible from ONE sentence. |
| 3 | `That XYZ Steel invoice — is it legit?` | `run_scenario validate_invoice` → **IMPERSONATION RISK: true** with 3 layered reasons (domain mismatch `ahmed.r@gmail.com` vs `@vendorxyz.com`, public mail provider = classic BEC, >$50k threshold). Cross-refs the legit AR email on file. | The fraud catch. Recommends out-of-band verification. Director leans in. |
| 4 | `Block this sender and notify finance` | Confirmation gate → "Confirm?" → `Yes` → executes block_vendor + notify_finance. | One-click structured action card executes deterministically (no LLM round-trip). The checkpoint/trust beat. |
| 5 | `RFI #18 is about column tolerances — who handles it?` | `run_scenario route_rfi` → classifies **structural**, and HR step now **surfaces Anita (structural specialist, score 3)** instead of 0. Delegates to architect/site_engineer in the same turn. | Brain delegates, doesn't narrate. The RFI is routed AND a matching specialist is found from the live pipeline. |
| 6 | `Draft the progress update for Sarah Jones` | `action: respond` with full draft + chips ["Send this reply","Edit draft","Cancel"]. | Draft-then-confirm safety on display. |
| 7 | `Send it` | `action: send_reply`, success. Email marked replied. | Real mutation, human-in-the-loop honored. |
| 8 | `Start evaluating the Gamma tender` | `action: create_project` → spawns the **EPC specialist roster (Cost Controller + Procurement + Safety/HSE + Proposal Writer) in parallel** on `gamma-feasibility`. | Scene-5 wow: 3-4 EPC agents working simultaneously in /agents — looks like an EPC project office, not a software repo. |

**Recovery line if state drifts mid-demo:** "Let me reset the office" → hit the reseed button (P1 item §2.6) → re-run from step 1.

---

## 2. API bugs to fix (ranked)

### P0-1 — Proactive fires never register (centerpiece "dynamic" beat is dark)
- **File:** `scripts/scenarios/load-epc-atlas.sh:12` — `HOST="${HOST:-http://localhost:8099}"` defaults to the wrong port; live server is `8080`, so the two cron jobs were never registered.
- **Repro:** `curl -s http://localhost:8080/api/cron/` → `{"jobs":null}`.
- **Fix:** change default to `8080`, AND always run the seed with explicit `HOST=http://localhost:8080`. Verify `cron.go:97-105` emits `trigger_fired` markers (it does) so the conductor briefing surfaces "N fires".

### P0-2 — process_inbox cannot detect the BEC email it was built to catch
- **File:** `internal/scenario/process_inbox.go:84-101` (vendor branch) calls `proc.ValidateInvoice(e.From, e.VendorID, 0)`; the seeded attack carries a fake `vendor_id="vendor_xyz_imposter"`, so `procurement.go:130-135` returns early with `ImpersonationRisk=false`.
- **Repro:** `curl -s -X POST '.../api/scenario/process_inbox/run?scope=atlas-site' -d '{"max":8}' | jq '.steps[-1].output.impersonations'` → `0`, while `/api/cap/email/summary` shows `1`.
- **Fix (two parts):**
  1. In the vendor branch, **trust on-disk `e.TrustStatus` first**: if `e.TrustStatus=="impersonation"` set `Severity="critical"`, `Verdict` from `e.TrustReason`, `Suggestion="block_vendor_and_notify_finance"`, `impersonations++` — before relying on ValidateInvoice.
  2. In `procurement.go:130-135`, when `vendorID` is unknown, **fall back to resolving the vendor by claimed name/domain** and run the domain-mismatch check anyway — that fake-id-with-real-name pattern IS the BEC signal.

### P0-3 — Demo state is consumable; centerpiece BEC email can be permanently deleted
- **File:** no reset route (`server.go` ~236-249). The chat delete handler (`brain_handlers.go:418-424`) physically removed `email_atlas_1_impersonation.json` during testing.
- **Repro:** `ls projects/inbox/email_atlas_1_impersonation.json` → MISSING after a session.
- **Fix:** add `POST /api/scenario/reseed?scope=atlas-site` that re-runs the seed logic idempotently (deterministic ids overwrite). See §2.6 / build order.

### P1-4 — /api/chat silently drops `scope`; Brain is scope-blind
- **File:** `internal/web/brain_handlers.go:91-94` — request struct is only `{Message, UserID}`.
- **Repro:** `curl -X POST .../api/chat -d '{"message":"How is the Atlas site going?","scope":"atlas-site"}'` → "Project atlas-site not found" (scope misread as project id). Contrast `scenario_handlers.go:88` which honors `scopeOf(r)`.
- **Fix:** add `Scope string \`json:"scope"\`` to the struct, thread into `Route`, prepend a one-line scope banner to the assembled context.

### P1-5 — route_rfi HR specialist match always returns 0
- **File:** `internal/scenario/route_rfi.go:86-89` builds `MustHave = append([]string{discipline}, hits...)`. In `hr.go:197` each missing must-have costs −2; applicant skills never contain RFI jargon ("column","plumb"), so a real structural engineer scores negative → clamped to 0 (`hr.go:218-219`) → dropped by the `score>0` filter (`hr.go:163`).
- **Repro:** `/api/cap/hr/match` with `must_have:["structural","column","tolerance","plumb"]` → `count:0`; moving the RFI keywords to `nice_to_have` → Anita at score 3.
- **Fix (one line):** `MustHave: []string{discipline}, NiceToHave: hits`. Optionally relax the gate in `hr.go:163` to `score >= 0` for partial-match specialists.

### P1-6 — Gamma create spawns the generic software orchestrator, not an EPC team
- **File:** `executor.go:90` (`executeCreateProject`) → `OnCreateProject` callback wired in `server.go` runs the generic template (role `architect`, phase `Template Selection`); `project.json` `team:null`.
- **Repro:** create Gamma → `GET /api/sessions` shows `agent_role:architect`; `project.json` `current_phase:"Template Selection"`.
- **Fix:** detect EPC scope/intent in the create path and select an EPC team template (Cost Controller, Procurement, Safety/HSE, Proposal Writer) + EPC phase set. See §3.4.

### P2-7 — job_application emails silently archived
- **File:** `process_inbox.go:148` default branch — switch has no `job_application` case; seeded `email_1774348831940` (a real CV) → verdict "internal · no action".
- **Fix:** add a `case "job_application":` that mirrors the `applicant` case (HR `.Store` + "added to pipeline").

### P2-8 — Triage JSON keys are PascalCase, breaking the shape-driven UI
- **File:** `process_inbox.go:69-72` local `Triage` struct has no json tags → serializes `ID/From/Subject/...` while the rest of the API is snake_case.
- **Fix:** add `json:"id"`, `json:"from"`, etc. (7 tags).

### P2-9 — `max` param ignores string input
- **File:** `process_inbox.go:28` — `ctx.Input["max"].(float64)` with no string fallback.
- **Fix:** add a string-coercion fallback.

---

## 3. Make scenarios dynamic (adapt to live data + input, not hardcoded)

### 3.1 — Thread `scope` into the chat lane (P1)
`brain_handlers.go:91` add `Scope` to the request struct; pass into `Route(ctx, userID, message, scope)`; in `context.go AssembleContext` prepend: `"You are operating on EPC scope: <scope>."` and pass `scope` into every `run_scenario` Input so scenarios filter to the right site. Unblocks every scoped scenario from chat.

### 3.2 — Add a `run_scenario` action so chat dispatches the multi-agent engine (P0, highest leverage)
- `actions.go:24` add `ActionRunScenario = "run_scenario"`.
- `executor.go:43` Execute add `case ActionRunScenario:` → new callback `e.OnRunScenario(name, scope, input)` that calls `scenarioEngine.Get(name).Run(scenario.Context{DataRoot, Scope, Input, Emitter})`. The wiring already exists in `scenario_handlers.go:85-91` — lift it into a callback set on the executor in the brain handler construction, reusing the same `scenarioEmitter{store, hub}` so per-agent steps stream to the workforce panel.
- The scenario result's structured `suggestions` (title/detail/action/payload) flow back as the chat response's suggestion cards (see §4.3).

This is the change that turns "process my inbox" from a solo answer into a visible fan-out to procurement/hr/pm/architect.

### 3.3 — Data-driven classifiers (remove recompile-bound Go literals) (P1/P2)
- Move `disciplineKeywords` (`route_rfi.go:32-42`) and the EPC keyword bag (`process_applicants.go:159-166`) into per-scope `data/<scope>/taxonomy.json` (`{disciplines, skills, open_roles}`), loaded by the runner.
- Add an open-roles store so `process_applicants` matches new applicants against **real configured JDs** instead of a derived static bag. New discipline/skill/role → no recompile, and scenarios react to the actual configured org.

### 3.4 — EPC team template for create_project (P1)
Map EPC-scope `create_project` to an EPC roster + phase set (Feasibility → Cost Estimate → Bid Decision) and write `team:[cost_controller, procurement, safety_officer, proposal_writer]` into `project.json`, so /agents shows 3-4 EPC agents in parallel — the Scene-5 wow.

### 3.5 — Scope-aware, path-robust inbox (P1)
`inbox.go:46-55` `NewInbox(_ string)` ignores scope and hardcodes the relative `projects/inbox`; a CWD change silently empties it (swallows `os.IsNotExist`). Anchor the path to an absolute `dataRoot` and accept `scope` (back-compat read of `projects/inbox`). Fixes the inconsistent tenancy where `morning_briefing` on a fresh scope shows the wrong inbox.

### 3.6 — Reseed endpoint (P1) — see §2 P0-3
`POST /api/scenario/reseed?scope=atlas-site` wraps the seed logic behind a handler (deterministic ids overwrite). Makes the demo repeatable on stage and restores the BEC email. Add a Lab/conductor button.

### 3.7 — Payload-carrying suggestions (P2)
Make slip + morning_briefing suggestions carry the specific `milestone_id`/`email_id` in `Payload` so a typed prompt or clicked chip dispatches against real records.

---

## 4. Brain / AI tuning (proactive + dynamic feel)

### 4.1 — Force proactive delegation, kill narration (P0, system-prompt only)
In Decision Rules (`router.go:501-516`) add:
> When a task clearly belongs to a discipline (RFI→architect/site_engineer, hiring→hr, vendor/invoice→procurement, schedule→pm, inbox triage→run_scenario), **DELEGATE or RUN_SCENARIO in the SAME turn** — do not merely state which agent should handle it. **If you name an agent in your `response`, your `action` MUST be `delegate` or `run_scenario`, never `respond`.**

Keep descriptive prose in `response` for the human; force the action. Removes the extra click that today makes delegation a turn-2 event.

### 4.2 — Document run_scenario + mapping in the system prompt (P0)
After `router.go:498` add `run_scenario` with the 6 names and triggers:
- "process/triage/sweep inbox" → `process_inbox`
- "check schedule/milestones/slipping" → `schedule_check`
- "morning briefing / what needs attention today" → `morning_briefing`
- "incoming RFI / who handles" → `route_rfi`
- "verify/validate vendor invoice / is this legit" → `validate_invoice`
- "rank/screen applicants" → `process_applicants`

### 4.3 — Unify the two proactive vocabularies (P1)
Chat returns `suggestions []string` (`actions.go:8`); scenarios return rich `{title,detail,action,payload}`. Promote `BrainDecision.Suggestions` to a struct array `{label, action, params}` (keep accepting plain strings for back-compat), and render scenario action cards (`block_vendor`, `notify_finance`) as clickable chips in the same chat lane. A click executes the pre-resolved action **deterministically** — no LLM re-route. Today `mission_v2.go:385` posts chip text back as a fresh message (turn-2 round-trip + mis-route risk).

### 4.4 — Structured JSON output from Gemini (P1)
`apiclient.go:298-301` generationConfig sets only `maxOutputTokens`+`temperature`. Add `responseMimeType:"application/json"` (and ideally `responseSchema` mirroring `BrainDecision`), drop temperature to ~0.3 for the routing call. Eliminates the prose-wrapped-JSON failure mode that forces the 3-strategy parser + `inferActionFromText` (`router.go:130-136, 329-380`).

### 4.5 — Startup warmup (P1, trivial)
After the client is built in `server.go` (~:106), fire one throwaway `SendMessage` in a goroutine (system "reply OK", user "warmup"). Observed cold start 6-13s vs warm 3-4s; the first on-stage message must feel instant.

### 4.6 — Surface morning_briefing on first load (P2)
Once run_scenario lands, call `morning_briefing` on `/conductor` first load so the brief is already on screen when the Director sits down — reinforcing proactivity before they type anything.

---

## 5. Build order (fastest path to a solid demo first, polish second)

### PHASE A — Make the demo *exist and not break* (≈45 min, all P0)
1. **Seed-port fix + run seed** (`load-epc-atlas.sh:12` → 8080; run with explicit HOST). Restores BEC email + registers the 2 cron fires. *(5 min)*
2. **process_inbox BEC detection** — trust `e.TrustStatus` first + name/domain fallback in `ValidateInvoice` (§2 P0-2). *(20 min)*
3. **Reseed endpoint** `POST /api/scenario/reseed` (§3.6) — wrap seed behind a handler so state is recoverable on stage. *(20 min)*

*Checkpoint:* re-run the narrative steps 1, 3, 4 manually via curl — BEC catch must report impersonation, fires must show in `/api/cron/`.

### PHASE B — Make it *visibly multi-agent and proactive* (≈3-4 hrs, P0)
4. **run_scenario action** end-to-end: `actions.go` + `executor.go` case + `OnRunScenario` callback in brain handler reusing `scenarioEmitter` (§3.2). *(2h)*
5. **System-prompt tuning**: force-delegate rule + run_scenario mapping (§4.1, §4.2). *(30 min, prompt only)*
6. **Scripted proactive fire**: a short-interval cron (or "Simulate incoming email" button) that, ~20-30s in, emits a `trigger_fired` impersonation card to the conductor (narrative step 0). *(1h)*

*Checkpoint:* "process my inbox" fans out 4 agents in the workforce panel; a fire card appears unprompted. This is the demo's spine.

### PHASE C — EPC polish + reliability (≈2-3 hrs, P1)
7. **route_rfi match fix** (one line, §2 P1-5) — Anita surfaces. *(15 min)*
8. **scope threading** into /api/chat (§3.1). *(30 min)*
9. **EPC team template** for Gamma create_project (§3.4) — Scene-5 parallel agents. *(2h)*
10. **Warmup goroutine** (§4.5) + **structured JSON output** (§4.4). *(45 min)*
11. **Unify suggestion cards** as clickable structured chips (§4.3). *(if time)*

### PHASE D — Nice-to-have (P2, only if Phases A-C land)
12. job_application case + Triage json tags + `max` string coercion (§2 P2-7/8/9).
13. Data-driven taxonomy.json classifiers (§3.3).
14. Scope-aware inbox path (§3.5).
15. morning_briefing on first load (§4.6) + minimal /email-sim modal.
16. Shift seeded dates forward so "handle today" framing isn't past-due (low; framing only — slip math is correct).

---

**Critical path to a working demo:** Phase A + items 4, 5, 6 of Phase B. If only those land, every narrative step except 8 (EPC parallel team) works, and steps 0/2/3 (proactive fire, fan-out, BEC catch) — the three biggest wows — all fire. Phase C makes it feel like an EPC product rather than a software demo; treat item 9 (EPC team) as the highest-value P1 if time allows.
