# Agent Capabilities + Scenario Engine — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans. Steps use checkbox (`- [ ]`) syntax.

**Goal:** Make EPC scenarios *demoable end-to-end* — agents with declared capabilities, backing data stores (applicants/vendors/milestones), and multi-step scenarios that compose them. After execution, scenarios surface proactive follow-up suggestions. Every layer testable via API.

**Architecture:**
- New `internal/capability/` package — pure Go data + business logic per capability area (HR, Procurement, Schedule), each persisted as JSON in `data/<scope>/`.
- New `internal/scenario/` package — named scripted flows that compose capabilities, emit messages into the existing store/WS hub, return structured result with suggested next actions.
- New HTTP layer at `/api/cap/...` and `/api/scenario/{name}/run` so both capabilities and scenarios are externally testable.
- Improved Atlas EPC bootstrap seeds applicants/vendors/milestones via the new APIs (not raw JSON on disk) — proves the APIs work.
- Test harness `scripts/test/test-scenarios.sh` exercises every endpoint and asserts expected state. Pass/fail per check, summary at end.

**Tech Stack:** Go (backend), bash + curl + jq (test harness). Frontend untouched — capabilities and scenarios are backend concerns; Conductor already streams the resulting messages via existing /ws.

---

## Phase 1 — Capabilities

### Task 1.1: HR capability (applicants store + JD matching)

**Files:**
- Create: `internal/capability/hr.go`
- Create: `data/atlas-site/applicants.json` (seeded by tests; not committed)

`HR` exposes:
- `Store(a Applicant) error` — write/append
- `List() []Applicant`
- `Match(jd JD) []ApplicantMatch` — score = simple keyword overlap on `applied_for`, `key_skills`, `certifications`; returns sorted top N
- `UpdateStatus(id, status) error`
- `data/<scope>/applicants.json` is the persistence file; created on first write.

Type:
```go
type JD struct { Title string; MustHave []string; NiceToHave []string; MinYears int }
type ApplicantMatch struct { Applicant Applicant; Score int; Matched []string; Missing []string }
```

### Task 1.2: Procurement capability (vendor lookup + invoice validation)

**Files:**
- Create: `internal/capability/procurement.go`

`Procurement` exposes:
- `ListVendors() []Vendor`
- `StoreVendor(v Vendor)`
- `ValidateInvoice(senderEmail, vendorID, amount) InvoiceCheck` — returns `{trusted, impersonation_risk, reasons, recommendation}` by comparing sender domain to vendor's known domain.

Type:
```go
type InvoiceCheck struct {
  Trusted bool; ImpersonationRisk bool;
  Reasons []string; Recommendation string;
}
```

### Task 1.3: Schedule capability (milestones + slip detection)

**Files:**
- Create: `internal/capability/schedule.go`

`Schedule` exposes:
- `ListMilestones() []Milestone`
- `UpdateMilestone(m Milestone)`
- `Slips() []Slip` — compares due dates to current date, returns anything > 1 day late

Type:
```go
type Milestone struct { ID, Title string; DueDate string; Status string; PctComplete int }
type Slip struct { Milestone Milestone; DaysLate int; Severity string }
```

### Task 1.4: HTTP layer for capabilities

**Files:**
- Create: `internal/web/capability_handlers.go`
- Modify: `internal/web/server.go` (register routes)

Routes:
- `GET  /api/cap/hr/applicants?scope=<id>` — list
- `POST /api/cap/hr/applicants?scope=<id>` — store one
- `POST /api/cap/hr/match?scope=<id>` — body: `JD` → returns matches
- `PATCH /api/cap/hr/applicants/{id}?scope=<id>` — body: `{status}`
- `GET  /api/cap/procurement/vendors?scope=<id>`
- `POST /api/cap/procurement/vendors?scope=<id>`
- `POST /api/cap/procurement/validate-invoice?scope=<id>` — body: `{sender_email, vendor_id, amount}`
- `GET  /api/cap/schedule/milestones?scope=<id>`
- `POST /api/cap/schedule/milestones?scope=<id>`
- `GET  /api/cap/schedule/slips?scope=<id>`

---

## Phase 2 — Scenario engine

### Task 2.1: Scenario runner core

**Files:**
- Create: `internal/scenario/engine.go`

```go
type Step struct { Agent string; Action string; Input map[string]any; Output map[string]any; OK bool; Note string }
type Result struct { Scenario, Project string; Steps []Step; Summary string; Suggestions []Suggestion; OK bool }
type Suggestion struct { Title, Detail string; Action string; Payload map[string]any }
type Runner interface { Run(scope string, input map[string]any) (Result, error) }
```

Engine emits a `TypeSystem` message per step into the existing store + hub so the Conductor's WS stream shows the multi-agent flow in real time.

### Task 2.2: process_applicants scenario

**Files:**
- Create: `internal/scenario/process_applicants.go`

Flow:
1. **Conductor** decomposes input ("hire Site Supervisor for Highway Bridge") into a JD object.
2. **HR** runs `Match(jd)` → returns ranked applicants.
3. If matches found → emit "HR found N applicants; top: X" → suggest `schedule_interviews`.
4. If none → suggest `post_job_listing`.

Returns `Result{Steps: [conductor.classify, hr.match], Suggestions: [...]}`.

### Task 2.3: validate_invoice scenario

**Files:**
- Create: `internal/scenario/validate_invoice.go`

Flow:
1. **Conductor** picks an email from inbox with `category="vendor"`.
2. **Procurement** runs `ValidateInvoice(email.from, vendor_id, parsed_amount)`.
3. If impersonation → suggest `block_vendor`, `notify_finance`.
4. If trusted → suggest `match_to_po`, `schedule_payment`.

### Task 2.4: schedule_check scenario

**Files:**
- Create: `internal/scenario/schedule_check.go`

Flow:
1. **PM** runs `Slips()`.
2. For each slip → propose mitigation (extra crew / re-sequence) as Suggestion.
3. If 0 slips → suggest `notify_client_on_track`.

### Task 2.5: Scenario HTTP layer

**Files:**
- Create: `internal/web/scenario_handlers.go`
- Modify: `internal/web/server.go`

Routes:
- `GET  /api/scenarios` — list registered scenarios
- `POST /api/scenario/{name}/run?scope=<id>` — body: scenario input → returns `Result` JSON. Also emits messages.
- `GET  /api/capabilities` — declares which agent owns which capability (for UI / docs).

---

## Phase 3 — Improved Atlas EPC seed

### Task 3.1: Rewrite Atlas bootstrap to seed via API

**Files:**
- Modify: `scripts/scenarios/load-epc-atlas.sh`

In addition to the existing 5 inbox emails:
- POST 3 applicants (Raj Kumar + 2 others with varying skills) via `/api/cap/hr/applicants?scope=atlas-site`
- POST 3 vendors (XYZ Steel legit, AlphaConcrete, BetaElectrics) via `/api/cap/procurement/vendors?scope=atlas-site`
- POST 5 milestones (some on-track, one slipping) via `/api/cap/schedule/milestones?scope=atlas-site`

---

## Phase 4 — Test harness

### Task 4.1: End-to-end test script

**Files:**
- Create: `scripts/test/test-scenarios.sh`

Checks (each prints PASS/FAIL and exits non-zero on any failure):
1. Server reachable
2. Seed Atlas via API; verify `GET /api/cap/hr/applicants?scope=atlas-site` returns 3
3. Verify `GET /api/cap/procurement/vendors?scope=atlas-site` returns 3
4. Verify `GET /api/cap/schedule/milestones?scope=atlas-site` returns 5
5. Run `process_applicants` scenario; assert `Result.Steps` has 2+ steps and `Suggestions` is non-empty
6. Run `validate_invoice` scenario with impersonation email; assert suggestion includes `block_vendor`
7. Run `validate_invoice` with trusted email; assert suggestion includes `match_to_po`
8. Run `schedule_check`; assert it identifies the seeded slip
9. Verify `/api/messages` now contains `TypeSystem` messages from the scenario runs
10. Verify `/api/scenarios` lists all 3 by name
11. PATCH applicant status, verify it persists

Output: numbered pass/fail lines; final `==> PASS (N/N)` or `==> FAIL (X/N)` line; exit code reflects.

---

## Out of scope (this plan)

- Conductor UI integration with scenario cards (next plan after backend solid).
- Agent prompt updates to use these tools (separate plan; needs MCP or output-parsing).
- Persistence beyond JSON files (Postgres etc.).
- Authentication on the new endpoints.
