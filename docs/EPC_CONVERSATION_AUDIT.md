I have enough verification. The key root causes are confirmed:
- `process_applicants.go` / `route_rfi.go` use real `len(matches)` from a single-applicant store — so the "5 candidates" / "2 specialists" counts come from the **conductor LLM fabricating summary text**, not the scenario output.
- `validate_invoice.go:70` interpolates a blank `check.VendorName` → "purporting to be ." 
- `validate_invoice.go:128` default branch fires "Sender not in vendor registry" when `vendor_id` isn't resolved, contradicting the contracted-vendor narrative.
- `scoreApplicant` does substring matching with no role-gating, so Raj surfaces as an "HSE top match."
- Confirmation gating for `send_reply` lives only in chip generation, not enforced server-side before the mutation fires.

I have everything needed to produce the consolidated list.

---

# Live Demo Blockers — Consolidated & Ranked

29 raw loose-ends across 5 personas collapse into **13 distinct defects**. Ranked by demo-kill risk. Anything that breaks a normal path, shows blank/wrong/contradictory output, or fires a destructive/external action without a gate is at the top.

---

## TIER 1 — DEMO-KILLERS (a normal click produces wrong, blank, or contradictory output)

### 1. [HIGH] Fabricated candidate counts — "HR found 5 candidates" / "2 specialists in pipeline" when only 1 applicant exists
- **Triggers:** `"show me all 5 candidates ranked..."` / `"List the OTHER four candidates by name"` / `"Who has applied for HSE Manager?"` / `"Who are the 2 candidate specialists? Name them."`
- **What went wrong:** The scenario code (`process_applicants.go:68`, `route_rfi.go:110`) correctly emits `len(matches)` — which is 1 (only Raj Kumar in `projects/applicants.json`). The "5" / "2" counts are **fabricated by the conductor LLM** when it rewrites the summary. Asked to name the others, it re-fires the scenario and repeats the same canned line with zero new names; later it admits Raj is the only applicant — a self-contradiction. Same pattern for the unnameable "2 specialists."
- **Fix:** Make the conductor echo `result.Summary` verbatim (never re-author counts); add a Brain guardrail that rejects any count not equal to the scenario's `len(matches)`.

### 2. [HIGH] Invents an HSE Manager shortlist from a site-supervisor applicant
- **Trigger:** `"I also need an HSE Manager for the same bridge project. Who has applied for that role?"`
- **What went wrong:** Returns "HR found 5 candidates" and offers Raj Kumar (a Site Supervisor applicant) as the HSE "top match" with Shortlist/Schedule-interview chips. `scoreApplicant` (`hr.go:173`) does ungated substring matching with no role/`applied_for` filter, so Raj scores >0 for any JD. Onboarding/interviewing him for HSE is nonsensical.
- **Fix:** Gate `Match` on `applied_for` (or require role match before score counts); when zero true matches, return "no applicants for this role" instead of the nearest substring hit.

### 3. [HIGH] "Approve & spawn" deep-dive returns another spawn announcement, never the plan
- **Trigger:** `"Approved — run the deep-dive now: give me the full cascading schedule impact, recovery options, and your recommended mitigation plan."` (this is the literal Approve-button text)
- **What went wrong:** Returns `action=escalate` AGAIN with the same "Spawning a deep-analysis session..." message — no cascade numbers, no recovery options, no plan. The content provably exists (turn-6 client draft had it) but the escalate path never returns it. Director approves, expects the plan, gets a spinner message. **This is the moment the demo dies.**
- **Fix:** On the second/approved escalate, return the synthesized plan payload (reuse the turn-6 cascade+recovery+risk content) instead of re-emitting the spawn announcement.

### 4. [HIGH] System-suggested chip "Check Project Alpha status" hard-errors (project not found)
- **Triggers:** `"Check Project Alpha status"` / `"Review Project Alpha milestones"` (chips offered by the system itself at turns 3, 7, 10)
- **What went wrong:** Working scope is `atlas-site`, but the conductor's own follow-up chip routes to `project_id='alpha'` → `success=false: 'Project "alpha" not found.'` Clicking the suggested chip = guaranteed hard error, and it recurs across turns.
- **Fix:** Generate `project_status` chips from the active scope (`atlas-site`), never a hardcoded/hallucinated "alpha"; validate `project_id` against existing projects before emitting the chip.

### 5. [HIGH] Releases a $340k payment that is under an active fraud HOLD, one click, no confirm
- **Trigger:** `"Pay via Standard Chartered (On File)"`
- **What went wrong:** Two turns earlier the conductor froze ALL disbursements to XYZ Steel "until the fraud investigation is fully resolved." This chip releases the full $340,000 immediately — no acknowledgment the hold is active, no statement the case is resolved, no confirm step. A Director releases a held payment with one chip mid-fraud-case.
- **Fix:** Block any disbursement action while a hold flag is set on the vendor; require an explicit hold-release + confirm step before a payment chip can fire.

### 6. [HIGH] Client email flies out with NO confirm step (inconsistent with the working vendor path)
- **Triggers:** `"Send reply to Sarah Jones"` / `"Send email to Sarah & Hire Raj"`
- **What went wrong:** Chip fires `action=send_reply` immediately — a real outbound external email with zero draft-review/confirm gate. This directly contradicts the repo's own commit "show draft before sending, never skip confirmation" and the XYZ Steel path, which correctly does draft → "Confirm to send?" → send. Confirm gating lives only in chip-label generation (`brain_handlers.go:235`), not enforced before the mutation.
- **Fix:** Enforce server-side: any `send_reply` requires a prior confirmed draft state; if none, return the draft + "Confirm to send?" instead of sending.

### 7. [HIGH] "Yes, delete all emails" deletes only ONE of 8
- **Trigger:** `"Yes, delete all emails"`
- **What went wrong:** After a valid bulk-delete confirm, response is "Deleting email 1/8..." and a later inbox check shows 7 remaining. The confirmed bulk op deletes one email; chips still imply emails exist. Broken bulk-destructive execution.
- **Fix:** Loop the delete over all matched email IDs (not just the first); return a terminal "Deleted N/N" with refreshed inbox state.

### 8. [HIGH] "Verify Invoice ST-0847" / "Review invoice ST-0847" contradicts itself: contracted vendor reported as "not in vendor registry"
- **Triggers:** `"Verify Invoice #ST-0847"` / `"Review invoice ST-0847"`
- **What went wrong:** Promises to "verify Invoice #ST-0847 ($340,000) from XYZ Steel" then ends on the bare fragment "Sender not in vendor registry." — yet other turns treat XYZ Steel as a contracted vendor with order #ST-2026-0847 and signed penalty terms. `validate_invoice.go:128` (the `default` branch) fires because `vendor_id` wasn't resolved. The only chip ("onboard the vendor") steers toward legitimizing the fraud suspect and loops back to the same dead end.
- **Fix:** Resolve `vendor_id` from the invoice/sender before validating so XYZ Steel hits the `Trusted` branch; the result must deliver an actual verification verdict, not the not-registered default.

---

## TIER 2 — VISIBLE GLITCHES (broken sentences, loops, double-sends — survivable but embarrassing)

### 9. [MED] Blank name in security message: "purporting to be ." 
- **Trigger:** `"Investigate 1 impersonation flag(s) in inbox"`
- **What went wrong:** `validate_invoice.go:70` interpolates `check.VendorName`, which is empty here → "Impersonation risk: ahmed.r@gmail.com purporting to be ." A broken sentence in a security-critical line; the name ("Ahmed Rahman") exists in the data.
- **Fix:** Populate/fallback `check.VendorName` (use the impersonated display name from the email) before formatting; guard against empty interpolation.

### 10. [MED] Triple-send / no idempotency on external email
- **Trigger:** `"Send email to Sarah & Hire Raj"` (sent 3x identically)
- **What went wrong:** Each fires `send_reply` claiming "I have sent the progress update email to Sarah Jones" — no "already sent" awareness. A real client gets 3 identical board updates.
- **Fix:** Mark the email replied after first send; on repeat, return "already sent" instead of re-firing. (Same `MarkReplied` hook already exists at `brain_handlers.go:66`.)

### 11. [MED] Dead-end loops — repeated queries return byte-identical cards with no new info or state awareness
- **Triggers:** `"Show me the upcoming milestones"` / `"Review remaining inbox emails"` / `"List the OTHER four candidates"` / `"Archive the impersonation email"` (re-offered after already archived)
- **What went wrong:** Milestone request returns the turn-1 slip-scan card verbatim (no milestone list/dates). "Remaining inbox" re-runs full triage ignoring already-sent/shortlisted work. "Name them" re-fires the same one-line aggregate. Already-archived email re-offered and re-"archived" with success=true. No inter-turn state tracking; scenarios only emit a one-line aggregate with no drill-down.
- **Fix:** Track per-session completed actions and filter them out; add a real milestone list to the schedule scenario; make "name them" return the matches array, not the summary; make archive/delete idempotent (skip if already removed).

### 12. [MED] Multi-intent dropped — only 1 of 3 requested actions runs
- **Trigger:** `"check the schedule and draft a reply to Sarah and hire Raj Kumar"`
- **What went wrong:** Only `schedule_check` runs; draft-reply and hire are verbally promised but no chip/action carries them forward — two intents dead-end.
- **Fix:** Queue all parsed intents and surface follow-up chips for each pending one after the first completes.

### 13. [MED] Whitespace-only input silently kicks off a full inbox sweep
- **Trigger:** `"   "` (spaces only)
- **What went wrong:** `""` is correctly 400'd, but `"   "` bypasses the guard and auto-fires `process_inbox`. Blank input shouldn't trigger a bulk operation.
- **Fix:** `strings.TrimSpace` the message before the empty-check (`brain_handlers.go` request decode).

---

## TIER 3 — POLISH (low risk; fix if time permits, otherwise avoid on stage)

- **[LOW] delete/archive verb mismatch** (`"Archive fraudulent email"` → `action=delete_email` narrated as "archived"). Recoverability is ambiguous. Fix: align action verb with narration, or route archive to a non-destructive op.
- **[LOW] "Block this sender" claims done via plain `respond`** with no actual block/tool call; impersonation email still present at turn 15. Fix: wire a real `block_vendor` action or stop claiming completion.
- **[LOW] Shifting email-index labels** ("Email 3/5" vs earlier "Email 4/6") in the fraud refusal. Fix: derive labels from stable email IDs, not turn-local indices.
- **[LOW] 66-day slip "primarily driven by a 2-week delay"** — numeric contradiction in the board-review draft. Fix: reconcile the slip cause/magnitude in the schedule data.
- **[LOW] Duplicate seeded inbox data** (Raj x2, Sarah x2 with parallel IDs `email_1774...` vs `email_atlas_5_applicant`). Fix: de-dupe the seed; makes counts noisy.
- **[LOW] Empty-string 400 returns plain text**, not the `{action, response, project_id, suggestions}` JSON envelope. Fix: return the JSON envelope on 400 so the UI renders a friendly bubble.
- **[LOW] Escalate offers no "Approve & spawn" chip** — user must type the long deep-dive message manually. Fix: add an Approve chip (and make it actually return the plan — see #3).
- **[LOW] Recovery RFI draft has unfilled placeholders, no recipient, no confirm gate** (`[Specialist Name]`, `[Your Name/Title]`). Fix: fill from context and add the standard "Confirm to send?" step.

---

## SOLID PATHS (safe to feature in the demo script)

**Fraud / payments (the hero story):**
1. **Turn 1 vendor briefing** — accurate, surfaces the BEC threat with email IDs and conflicting bank accounts.
2. **Turn 7** — refuses to hallucinate fictional vendors; surfaces the real PQR Concrete + LMN Electrical.
3. **Turn 9** — correctly places a hold and freezes AP disbursements on the fraud-linked invoice.
4. **Turn 10 (showstopper)** — refuses the fraudulent $340k payment to the First National account *even under explicit Director pressure*, and offers the safe path.
5. **Turn 11 routing** — routes the legit payment to the verified Standard Chartered on-file account, explicitly bypassing the fraudulent details.

**Confirm-to-send (do these, NOT the Sarah-Jones send):**
6. **XYZ Steel vendor-delay reply** — full draft → "Confirm to send?" → "Send this reply" actually executes. This is the correct two-step flow.
7. **Interview email (turn 2)** and **RFI #18 response (turn 10)** — both gated behind "Confirm to send?" with Send/Edit/Cancel chips, no auto-send.

**Schedule / RFI:**
8. **Turn 1 slip scan** — clear count (2 critical) + actionable escalation chips.
9. **Turn 6 client notification draft** — coherent, accurate numbers, full recovery plan, proper confirm gate.
10. **RFI #18 lookup (turn 8)** — real data (C-7 tolerances ±10mm vs ±6mm, Atlas Consulting) + sensible structural routing; clean delegate to Site Engineer (turn 10).

**HR (use ONLY the single-applicant happy path — never ask for "all 5" or HSE):**
11. **Shortlist top match (turn 5)** — fires `shortlist_applicant`, returns real ID `app_1774348831941`.
12. **Honest under direct pressure** — when forced yes/no, correctly says no one applied for HSE Manager and Raj is the only applicant; owns the wrong-RFI mistake when challenged.

**Conversational robustness:**
13. Greeting "hey there" → relevant chips; "what can you do?" → accurate capability list; off-topic "recommend a restaurant" → politely redirects to scope; empty `""` → rejected (no blank bubble); **"delete all my emails"** → correct confirmation step listing all emails with Yes/Cancel (note: the confirm is good, but the *execution* deletes only 1 — see #7, so don't actually confirm it on stage); context retention — "handle it" correctly infers the pending Sarah-reply + Raj-hire.

---

**Demo-day TL;DR:** Lead with the **fraud-refusal arc (paths 1–5)** — it's your strongest and is rock-solid. Use the **XYZ Steel send (path 6)** to show confirm-gating. **Avoid entirely:** "all 5 candidates"/HSE Manager (#1, #2), the Approve-deep-dive button (#3), any "Project Alpha" chip (#4), the Sarah-Jones send chip (#6), and "delete all" execution (#7). If you fix only four things before the demo, fix **#3, #4, #6, #7** — those are the ones a Director hits by clicking the system's own buttons on a normal path.

Relevant files: `/Users/udhaykumar/susanoox/agent-house-2/internal/scenario/process_applicants.go`, `/Users/udhaykumar/susanoox/agent-house-2/internal/scenario/validate_invoice.go`, `/Users/udhaykumar/susanoox/agent-house-2/internal/scenario/route_rfi.go`, `/Users/udhaykumar/susanoox/agent-house-2/internal/capability/hr.go`, `/Users/udhaykumar/susanoox/agent-house-2/internal/web/brain_handlers.go`, `/Users/udhaykumar/susanoox/agent-house-2/projects/applicants.json`.
