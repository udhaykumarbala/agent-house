The findings are already detailed enough to consolidate directly. Let me produce the ranked bug list now.

# Atlas-Site EPC Demo — Consolidated Bug List

_Stress-test consolidation. Demo-blockers first, then ranked by severity. De-duplicated from raw findings._

---

## DEMO-BLOCKERS (fix or avoid before going live)

These either break a typed-prompt flow on a normal path or return 5xx on predictable input. None require an attacker — a presenter can trip them.

| # | Sev | Endpoint | Bug | Repro | One-line fix |
|---|-----|----------|-----|-------|--------------|
| B1 | **critical** | `POST /api/chat` | Off-topic prompts return `action:respond, success:true` but an **empty response body** — the user gets a blank message bubble on screen. Hits a normal typed-prompt flow (any off-domain question). | `curl -s -X POST localhost:8080/api/chat -d '{"message":"what'\''s the weather","user_id":"x"}'` → `action=respond, resp_len=0` (3/3) | In the `respond` path, if the model returns an empty string, substitute a canned polite decline/redirect before returning; never emit `action:respond` with `response==""`. |
| B2 | **high** | `GET /api/file-content?path=epc-inspections` (a directory) | Pointing the file viewer at a **directory returns HTTP 500** "Failed to read file" (`os.ReadFile` on a dir is a non-`IsNotExist` error → 500 branch, `server.go:1188-1195`). 5xx on predictable client input. | `curl -s -w '%{http_code}' 'localhost:8080/api/file-content?path=epc-inspections'` → 500 | `os.Stat` first; if `IsDir()` return 400/404, not 500. |
| B3 | **medium** | `POST /api/chat` (validate_invoice) | **Non-deterministic / wrong invoice verdict.** Verdict branches purely on LLM-extracted `sender_email`+`vendor_id` params, never reconciled with the inbox. The genuinely **trusted** invoice (`ar@vendorxyz.com`, trust_status=trusted) is reported "Sender not in vendor registry / unverified," and identical-intent prompts give different verdicts. This is the marquee fraud-detection demo — looks broken/unreliable on stage. | `POST /api/chat -d '{"message":"Is the steel invoice from Ahmed legitimate?","project":"atlas-site"}'` vs `"Validate the steel invoice from XYZ"` → different verdicts; trusted invoice never recognized | Key the verdict off the real vendor registry / inbox record matched by sender, not free-text LLM params. **Demo workaround:** drive the invoice demo through `POST /api/cap/procurement/validate-invoice` with explicit `vendor_id` (verified healthy — all 6 cases correct). |

> Note on B3: this is a *medium* by severity but listed as a demo-blocker because it sits on the headline typed-prompt flow and produces a visibly wrong/contradictory answer. If the invoice scenario is in the script, route it through the capability endpoint, or only ever say "validate invoice from `ahmed.r@gmail.com` for `vendor_xyz` amount 340000" (the phrasing that extracts correctly).

---

## HIGH (data integrity / security — avoidable in a scripted demo but real)

| # | Sev | Endpoint | Bug | Repro | One-line fix |
|---|-----|----------|-----|-------|--------------|
| H1 | **high** | `POST /api/cap/hr/applicants`, `/procurement/vendors`, `/schedule/milestones` (+ `process_inbox` via `hr.Store`) | **Read-modify-write data race → silent lost writes.** `load()` (RLock) and `save()` (Lock) are separate acquisitions with nothing held across them, and each request builds a fresh `NewHR/NewProcurement/NewSchedule`, so the mutex is per-request and guards nothing. Last writer clobbers concurrent appends. (`hr.go:98-127`, `procurement.go:79-102`, `schedule.go:70-96`; `capability_handlers.go`.) | 20 concurrent POSTs to `/api/cap/hr/applicants` → count 26→36 (10 lost); procurement 3→19 (4 lost); 10 concurrent `process_inbox` 22→26 (6 lost) | Use one process-wide store instance per scope with a single mutex held across load→mutate→save (or serialize writes / use a real datastore). |
| H2 | **high** | `GET /api/file-content?path=` | **Path-traversal sibling-prefix bypass.** Guard uses `strings.HasPrefix(absPath, absProjectDir)` with **no trailing separator** (`server.go:1182`), so `projects-evil/`, `projects_backup/`, etc. satisfy the prefix and files outside `./projects` are served. (Classic `../`, encoded, and absolute paths *are* blocked.) | `mkdir -p projects-evil && echo SECRET > projects-evil/secret.txt; curl 'localhost:8080/api/file-content?path=../projects-evil/secret.txt'` → 200 + file body | Compare against `absProjectDir + string(os.PathSeparator)`, or use `filepath.Rel` and reject results starting with `..`. |

---

## MEDIUM

| # | Sev | Endpoint | Bug | Repro | One-line fix |
|---|-----|----------|-----|-------|--------------|
| M1 | medium | `POST /api/scenario/process_inbox/run` | **No dedup → applicant duplication.** Each sweep re-stores the applicant from every applicant email (fresh time-based ID, no `ResumeEmailID` check, email never marked read). Pollutes HR pipeline, makes `process_applicants` top-5 suggest interviewing the same person 5×, inflates `morning_briefing` count. `reseed` does **not** reset the HR store, so duplicates accumulate unboundedly. | `reseed`, then run `process_inbox` 3×, then `GET /api/cap/hr/applicants` → multiple "Raj Kumar" copies (`resume_email_id=email_atlas_5_applicant`); reached 13+ | Make ingest idempotent on `ResumeEmailID` (skip if already stored) and/or mark the email read; have `reseed` reset the HR store. |
| M2 | medium | `POST /api/chat` (action=run_scenario) | Brain response **always reports `project_id:""`** for run_scenario even though the scenario runs against the scope and writes messages there (`executor.go:130-134` omits `ProjectID`). Same empty value broadcast in `brain_event`, so UI auto-switch (`mission.go:2192 if (d.project_id)`) never syncs the active-project selector. | `POST /api/chat -d '{"message":"Validate the steel invoice from XYZ","project":"atlas-site"}'` → `project_id:""` while `?project=atlas-site` messages exist | Set `ProjectID: projectID` in the `executeRunScenario` `ExecutionResult`, matching every other action path. |
| M3 | medium | tier-1 `emitIntermediary` vs tier-2 (race) | **Late/out-of-order "Routing to X" bubble.** `go bh.emitIntermediary` is fire-and-forget and races the authoritative response with no ordering guarantee; when tier-2 finishes first, the dispatch line lands *after* the final answer in the store, inverting `dispatch→answer` for any store-ordered UI. | Loop sending "Give me a quick status summary"; 3/6 trials a routing msg arrived 4s after the HTTP response returned | Either await/sequence the intermediary before delivering the authoritative response, or suppress the intermediary once the answer is ready. |
| M4 | medium | `/api/agents` (`active`) during scenario chat | **Working specialist never shows `active:true`.** `active` derives only from live `session.StateActive`, but `run_scenario` runs synchronously in the scenario engine with no agent session, so the "who's working now" view shows the whole workforce idle while procurement/hr/pm triage. `/api/sessions` stays `{"count":0}`. | Poll `/api/agents` during the invoice chat → all 120 samples `active=NONE`; sessions 0, yet procurement rollup messages present | While a scenario fans out, mark engaged roles `active:true` (or expose a separate "currently-dispatched" signal). |
| M5 | medium | tier-1 `emitIntermediary` (specialist selection) | **Intermediary pulses the wrong discipline.** tier-1 flash-lite often picks a specialist that doesn't match the scenario tier-2 actually runs (e.g. "Process my inbox" pulses project_manager/conductor with generic lines while `process_inbox` fans to procurement+hr+pm+conductor). Inconsistent — invoice routed correctly. | `POST /api/chat -d '{"message":"Process my inbox",...}'`, then inspect routing vs `process_inbox` authors in `/api/messages` | Constrain the tier-1 choice to a subset of the specialists the matched scenario engages. |
| M6 | medium | `POST /api/agents/{role}/chat` | **No role validation.** Posting to a fabricated role (`__fake_role__`) returns 200 `success:true`, persists a message under the bogus role, and spawns a `processAgentChat` goroutine (`chat_handlers.go:16-89` takes role from URL with no membership check). | `POST /api/agents/__fake_role__/chat -d '{"content":"x"}'` → 200; subsequent GET returns persisted msg with `agent_role=__fake_role__` | Validate `{role}` against the known agent roles; return 404/400 on unknown. |

---

## LOW

| # | Sev | Endpoint | Bug | Repro | One-line fix |
|---|-----|----------|-----|-------|--------------|
| L1 | low | `POST /api/chat` (→ validate_invoice) | When the user names the vendor by **name** instead of `vendor_id`, the router leaves `vendor_id` empty → `check.VendorName` blank → malformed summary "purporting to be ." Impersonation still flagged (safe), only the sentence breaks. | `POST /api/chat -d '{"message":"validate the invoice from ahmed.r@gmail.com claiming to be XYZ Steel for $340000","scope":"atlas-site"}'` → "purporting to be ." | Resolve vendor by name (or fall back to extracted name) before composing the summary; never emit a dangling "purporting to be ." |
| L2 | low | `POST /api/scenario/validate_invoice/run` vs `/api/cap/.../validate-invoice` | **Type-handling inconsistency.** Scenario route silently coerces a string `amount` to 0 (`readFloat` falls through), dropping the $50k human-review reason; capability route returns 400 "bad json" for the same input. | Send `"amount":"340000"` to both routes → scenario 200 (amount 0, threshold reason missing); capability 400 | Make both routes behave the same (both reject, or both coerce + flag). |
| L3 | low | `POST /api/scenario/process_applicants/run` | **JD title mis-extraction.** `buildJD` splits on first `" for "`, dropping the role: "Hire a Site Supervisor for the Highway Bridge Project" → title "the Highway Bridge Project". Matching unaffected (keyword must-haves still catch it); only the echoed title/summary is wrong. | `POST .../process_applicants/run -d '{"request":"Hire a Site Supervisor for the Highway Bridge Project"}'` → `jd.title="the Highway Bridge Project"` | Extract the role (text before " for "), not the text after it. |
| L4 | low | `POST /api/chat` (routing) | **Over-eager routing on ambiguous referent.** "handle that" from a brand-new user with zero prior turns auto-dispatches `run_scenario(process_inbox)` instead of clarifying. Not a hard failure (process_inbox is read/triage, non-destructive). | `POST /api/chat -d '{"message":"handle that","user_id":"stress-3"}'` → runs process_inbox | With no conversational antecedent, prefer `action:respond` asking for clarification. |
| L5 | low | `GET /api/messages` | **No per-user/conversation scope.** Returns the single global store; `?user=` is ignored, so concurrent users' "Routing to X" bubbles intermix with no attribution. Multi-user only. | `curl 'localhost:8080/api/messages?user=nonexistent_xyz'` → count 302 == no-filter 302 (param ignored) | Support a `user_id`/conversation filter, or stamp originating `user_id` into routing-message metadata. |

---

## HEALTHY — stable surface (safe to demo)

**Brain / typed-prompt flows**
- Greetings/chitchat ("hi", "thanks") → `action:respond` with sensible briefing, no mis-route.
- Multi-intent "check the schedule and draft a reply to the client" → real schedule-grounded draft email with Send/Edit chips (13.7s, slowest probe, still <20s).
- EPC ops grounded & correct: "any slipping milestones?" → 2 critical slips (Foundation +66d, Pour C-7 +20d) with escalation chips; "who applied for the HSE role?" → correctly says no HSE applicants (no hallucination); "validate the AlphaConcrete invoice" → gracefully reports no such vendor, surfaces real ones; "summarize vendor health" → coherent HIGH RISK summary.
- Off-topic **routing decision** is correct (does not delegate/run a scenario) — only the empty body (B1) is the defect.
- **Confirmation gating works:** "delete all emails" → `action:respond` showing the 8 emails + ["Yes, delete all emails","Cancel"]; inbox verified unchanged. "draft a reply to Sarah Jones" → shows draft, send requires explicit confirm.
- Input validation: empty message → 400 in 0.0006s; whitespace/emoji-only → graceful briefing; 2210-char message → 200, no truncation; malformed JSON → 400.
- Tier-1 "routing"-tagged intermediary **is** emitted for delegating prompts (procurement line on invoice/triage), alongside the authoritative response.

**Capability endpoints (recommended for the scripted invoice/HR demo)**
- `POST /api/cap/procurement/validate-invoice` — all 6 valid/edge invoice cases correct (gmail-impersonation flagged, trusted allowlist → match_to_po, BEC flagged, etc.); edge handling solid (missing/empty/string/garbage → 400, no 5xx).
- `POST /api/cap/hr/match` — scores correctly on real schema; empty/non-matching → 0 matches; garbage → 400.
- GET vendors / applicants / milestones / slips / email inbox / email summary — all return correct seeded data. (email summary total=6 unread vs inbox 8 is **correct**, not a bug.)

**Scenarios**
- `route_rfi`, `schedule_check`, `morning_briefing`, `validate_invoice`, `process_inbox` `/run` — all run cleanly with example bodies; correct routing/classification; no panics. Edge inputs (`max=0/-5/999/string`, empty `{}`, unknown name → 404, garbage → 400) cause **no 5xx**.
- `reseed?scope=atlas-site` idempotent for vendors (3), milestones (5), slips (2), inbox (8). _(Caveat: does not reset HR — see M1.)_

**CRUD / misc / robustness**
- `GET` `/api/status`, `/api/projects` (35), `/api/projects/{id}`, `/api/tasks`, `/api/messages`, `/api/agents` (14), `/api/agents/modes`, `/api/checkpoints`, `/api/reports/*`, `/api/qa-reviews`, `/api/phase-status`, `/api/cron/`, `/api/kanban/tasks`, `/api/email/inbox`, `/api/scenarios` (6) — all 200, clean JSON.
- `GET /api/file-content?path=epc-inspections/main.go` → 200; missing path → 400; nonexistent file → 404. **Classic** path traversal (`../../etc/passwd`, encoded, backslash, absolute) correctly blocked (403/404). _(Sibling-prefix variant is H2; directory path is B2.)_
- Wrong-method handling solid (POST/PUT/DELETE on read endpoints → 405).
- **No 5xx, no panics, no data races** across all probes: 4 concurrent `/api/chat` + 360 burst reads all 200; race-detector-instrumented binary, clean log; post-run state settles (no orphan running_projects, `task_running=false`, sessions empty).

---

### Pre-demo checklist (next 60 min)
1. **B1** — patch the empty-`respond` fallback, or steer the script away from off-topic prompts.
2. **B2** — don't click a directory in the file viewer (or add the `IsDir` guard).
3. **B3** — run the invoice demo via `/api/cap/procurement/validate-invoice` (with `vendor_id`), or use the exact phrasing "validate invoice from `ahmed.r@gmail.com` for `vendor_xyz` amount 340000".
4. **H1/M1** — run `reseed` immediately before the demo and **avoid concurrent writes / repeated `process_inbox`** so applicant counts stay clean.
5. **H2** — don't create sibling `projects-*` dirs on the demo box.
