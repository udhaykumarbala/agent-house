# Atlas EPC — Conductor Demo Script

A live walkthrough of the **Atlas Construction · Site 02** AI project office. The
Director sits at the **Conductor** (`/conductor`) and runs the whole site by
typing plain English. The Conductor understands the request, **dispatches to the
right specialist agents** (they pulse live in the workforce panel), and comes
back with an answer **plus proactive next-step suggestions** as one-click chips.

> Model: **Gemini 3.5 Flash** (native API). Warm responses ~3–6s. The server
> warms the model on boot, so the first prompt is already fast.

---

## 0. Pre-demo setup (do this once, ~30s before you present)

```bash
# 1. Start the server (from the repo root, with `claude` on PATH)
/tmp/agent-house --serve --port=8080 --project=./projects

# 2. Load / reset the Atlas demo state (vendors, applicants, milestones,
#    inbox incl. the impersonation email, and the proactive cron fires)
curl -sX POST 'http://localhost:8080/api/scenario/reseed?scope=atlas-site'

# 3. Open the Conductor
open http://localhost:8080/conductor
```

**To reset mid-demo at any time** (restores the BEC email + clean state):
`curl -sX POST 'http://localhost:8080/api/scenario/reseed?scope=atlas-site'`

Quick health check (optional): the briefing should show fires and the inbox
should hold 1 impersonation:
```bash
curl -s 'http://localhost:8080/api/cap/email/summary?scope=atlas-site' | jq .impersonations   # → 1
curl -s 'http://localhost:8080/api/cron/' | jq '.jobs|length'                                   # → 2
```

---

## The script — type each prompt into the Conductor

Everything below is **literally typeable**. Each row: what you type, what the
Conductor does, and the moment to point at.

### 1. The morning brief  ☕
**Type:**
> `Good morning. What do I need to handle today?`

**What happens:** action `run_scenario` → **morning_briefing**. The PM, inbox,
HR, and procurement are all consulted; you get a one-page brief:
- 2 schedule slips (Foundation **62 days late**, Slab C-7) — 2 critical
- **1 impersonation flag** in the inbox
- 7 unread emails, HR pipeline, vendor health

**Proactive chips:** *Escalate 2 critical schedule slips · Investigate 1 impersonation flag · Process 7 unread emails · Review HR pipeline*

**Say:** "I didn't ask it to check anything specific — it read the whole site and told me what matters."

---

### 2. Fan out to the team  👥  (the multi-agent moment)
**Type:**
> `Process my inbox`

**What happens:** action `run_scenario` → **process_inbox**. Watch the
**workforce panel**: **procurement, project_manager, and hr** light up as the
Conductor routes each email to the right discipline. Result: *7 triaged · **1 impersonation flag** · 1 new applicant · 3 client items.*

**Proactive chips:** *Block 1 impersonation sender · Draft 3 client replies for review · Match 1 new applicant against open JDs*

**Say:** "One sentence, and four specialists just worked in parallel. That's the difference from a chatbot."

---

### 3. Catch the fraud  🚨  (the wow)
**Type:**
> `Is that XYZ Steel invoice from Ahmed for 340000 legit?`

**What happens:** action `run_scenario` → **validate_invoice**. The Conductor
contrasts the **legitimate** invoice from `ar@vendorxyz.com` with the
**suspicious** email from `ahmed.r@gmail.com` claiming to be Ahmed Rahman —
flags it as a **BEC / impersonation attempt** (public-mail domain, >$50k,
payment-detail change).

**Proactive chips:** *Block this sender · Notify finance & PM · Verify via known channel*

**Say:** "It caught a $340,000 fraud attempt that looks exactly like a real vendor email — before anyone paid it."

---

### 4. Act on it (with a safety gate)  ✋
**Type:**
> `Block that gmail sender and notify finance`

**What happens:** the Conductor shows **what it will do first** and asks you to
confirm — nothing destructive happens without your click.

**Proactive chips:** *Send alert and archive email · Edit draft · Cancel*

**Say:** "Nothing happens behind my back. I approve, then it executes."

---

### 5. Route a technical query  🛠️
**Type:**
> `RFI 18 is about column tolerances on grid C-7 — who handles it?`

**What happens:** action `run_scenario` → **route_rfi**. Classified as
**structural**, assigned to the **Site Engineer**, and HR surfaces a matching
specialist already in the pipeline — **Anita Verma (structural, 9 yrs)**.

**Proactive chips:** *Assign to Site Engineer · Hire a structural specialist · Draft an RFI response*

**Say:** "It routed the question AND found me a person who can answer it."

---

### 6. Draft a client update  ✍️
**Type:**
> `Draft a progress update reply for the client Sarah Jones`

**What happens:** a full, ready-to-send draft appears in the chat (Project Alpha
Phase 2 progress), grounded in the site's real schedule data.

**Proactive chips:** *Send this update · Edit draft · Check Project Alpha schedule*

**Then type:**
> `Send it`

**What happens:** action `send_reply` → the update is sent and the email is
marked replied.

**Say:** "Draft, review, send — the human stays in the loop the whole way."

---

## Optional closer — the bid pipeline (Scene 5)
**Type:**
> `Start evaluating the Gamma highway bridge tender`

**What happens:** spins up a new project workspace and puts an **EPC specialist
team** on it. (Open `/agents` to show the roster working in parallel.)

---

## One-liner for the room
> "Everything you saw was one person typing plain English. The Conductor read
> the inbox, dispatched four specialists, caught a fraud attempt, routed a
> technical query, and drafted a client reply — and asked my permission before
> doing anything irreversible. In production we swap the demo data for your
> SAP / Primavera / Exchange. Zero agent changes."

## If something drifts
- Reset state: `curl -sX POST 'http://localhost:8080/api/scenario/reseed?scope=atlas-site'`
- The agents you'll see pulse: **procurement, project_manager, site_engineer, hr** (EPC) and the **Conductor** itself.
- All prompts above are verified working on Gemini 3.5 Flash.
