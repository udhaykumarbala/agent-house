# EPC Demo — "A Day in the Life of a Project Director"

## Design Iterations

**Iteration 1 — Full EPC OS with SAP, Primavera, OCR**
→ Contradicted: Multi-month effort. Fake integrations look fake. A demo that's clearly mocked is worse than no demo.

**Iteration 2 — Simple report generation from dummy DB**
→ Contradicted: Any chatbot can generate a report. Doesn't show why Agent House is different from ChatGPT.

**Iteration 3 — Show workflow, not data**
→ Contradicted: Still a chatbot demo. The differentiator is multi-agent orchestration, not single-agent chat.

**Iteration 4 — Multi-agent bid generation**
→ Contradicted: One workflow isn't enough. EPC directors manage 5+ projects daily. Need to show ongoing management, not just project creation.

**Iteration 5 — The synthesis: A Day in the Life**
→ Pre-loaded scenario with 3 active projects, emails, invoices, vendors. Shows quick tasks (email reply, invoice check) AND full pipelines (new bid). Shows the Brain as a natural interface to a team of specialists.

---

## Demo Concept

**Audience:** EPC company executives (COO, Head of Projects, IT Director)
**Duration:** 20 minutes
**Tagline:** "Your AI project office. Ask anything, delegate anything, review everything."

### Pre-Loaded State

The demo environment starts with realistic state — not an empty workspace.

**3 Active Projects:**
- **Alpha (Solar Farm)** — In execution phase, 60% complete, on budget
- **Beta (Water Treatment)** — In planning phase, bid submitted, awaiting client
- **Gamma (Highway Bridge)** — New tender received yesterday, needs evaluation

**Data Files (agents read these like they'd read a database):**
```
demo/epc-workspace/
  data/
    projects.json         ← 3 projects with status, budget, timeline
    vendors.json          ← 15 vendors with contact, rating, contracts
    invoices.json         ← 8 invoices (2 pending approval)
    costs.json            ← Historical cost data for past projects
    employees.json        ← 30 employees with roles, availability
    schedules.json        ← Project milestones and deadlines

  inbox/
    email_001.json        ← Vendor XYZ: "Steel delivery delayed 2 weeks"
    email_002.json        ← Client ABC: "When can we expect Phase 2 update?"
    email_003.json        ← Internal: "Site engineer requesting additional resources"
    email_004.json        ← Vendor PQR: "Invoice #1042 for concrete supply"
    email_005.json        ← HR: "3 candidates shortlisted for site supervisor"

  outbox/                 ← Drafts go here for review
  reports/                ← Generated reports land here

  project.json            ← Agent House project metadata
```

---

## Demo Flow (20 Minutes)

### Scene 1: Morning Briefing (2 min)
**Shows: Brain as natural interface**

Director opens Mission Control. Types:

> "Good morning. What do I need to handle today?"

**Brain responds (~3 seconds):**
> Good morning. Here's your daily briefing:
>
> **3 active projects:**
> - **Alpha (Solar Farm)** — 60% complete, on budget. ⚠️ Vendor delay reported.
> - **Beta (Water Treatment)** — Awaiting client response on bid.
> - **Gamma (Highway Bridge)** — New tender received. Needs evaluation.
>
> **Action items:**
> - 2 invoices pending approval ($45K + $28K)
> - 1 vendor delay needs response (steel delivery, 2-week slip)
> - 1 client email awaiting reply (Phase 2 update request)
> - HR has 3 candidates for site supervisor role
>
> What would you like to tackle first?

*Brain assembled this by reading all JSON files in <100ms, then one API call.*

---

### Scene 2: Handle Vendor Delay (4 min)
**Shows: Multi-agent collaboration, email management**

> "Tell me about the vendor delay on Alpha"

**Brain → delegates to Procurement Lead (session mode):**
- Reads `inbox/email_001.json` (vendor's delay notification)
- Reads `vendors.json` (contract terms for Vendor XYZ)
- Reads `schedules.json` (Alpha project timeline impact)
- Returns analysis:

> **Vendor XYZ — Steel Delivery Delay**
>
> **Situation:** 2-week delay on structural steel (150 tons). Original delivery: April 5.
> **Contract terms:** Max allowable delay: 7 days. Penalty: 0.5% per day beyond.
> **Impact:** Critical path item. Delays foundation work by ~10 days.
> **Recommendation:** Invoke penalty clause. Request expedited partial delivery.

Director: "Draft a reply to the vendor. Firm but professional."

**Brain → delegates to Procurement Lead:**
- Writes draft to `outbox/reply_vendor_xyz.md`
- Shows draft in chat:

> Subject: Re: Delivery Schedule Update — Project Alpha Steel Supply
>
> Dear Mr. Rahman,
>
> Thank you for the update. A 2-week delay on the structural steel delivery
> significantly impacts our critical path on Project Alpha.
>
> Per our contract (Clause 12.4), the maximum allowable delay is 7 calendar
> days. Any delay beyond this triggers a penalty of 0.5% of the order value
> per additional day.
>
> To mitigate the impact, we propose:
> 1. Expedited partial delivery of 80 tons by April 8
> 2. Remaining 70 tons by April 12
>
> Please confirm feasibility within 24 hours.

Director: "Looks good. Send it."

*In the live dashboard, audience sees the Procurement agent's tool calls — reading files, writing the draft.*

---

### Scene 3: Invoice Processing (3 min)
**Shows: Data analysis, human-in-the-loop approval**

> "Process the pending invoices"

**Brain → delegates to Cost Controller (session mode):**
- Reads `invoices.json` (2 pending invoices)
- Reads `costs.json` (budget for each project)
- Reads `vendors.json` (contract amounts)
- Cross-references and returns:

> **Invoice Review:**
>
> | Invoice | Vendor | Amount | Contract | Status |
> |---------|--------|--------|----------|--------|
> | #1042 | PQR Concrete | $45,200 | $38,000 | ⚠️ $7,200 OVER (18.9%) |
> | #1043 | LMN Electrical | $28,500 | $30,000 | ✅ Within budget |
>
> **Recommendation:**
> - #1042: REJECT — exceeds contract by $7,200. Request itemized breakdown.
> - #1043: APPROVE — standard payment terms apply.

**CHECKPOINT appears in UI:** Director clicks Approve for #1043, Reject for #1042.

Cost Controller writes rejection email draft to `outbox/reject_invoice_1042.md`.

---

### Scene 4: Weekly Report (3 min)
**Shows: Report generation across multiple projects**

> "Generate the weekly report for all projects"

**Brain → delegates to Report Generator (session mode):**
- Reads ALL project data files
- Reads cost data, schedule data, vendor status
- Generates `reports/weekly-report-2026-W12.md`:

> # Weekly Project Report — Week 12, 2026
>
> ## Executive Summary
> 3 active projects. Total portfolio value: $12.4M. 1 risk item (Alpha steel delay).
>
> ## Project Alpha — Solar Farm
> - **Progress:** 60% → 62% (+2% this week)
> - **Budget:** $4.2M spent of $6.8M (61.8%) — ON TRACK
> - **Schedule:** 2-week delay on steel. Mitigation in progress.
> - **Risk:** HIGH — vendor delay impacts critical path
>
> ## Project Beta — Water Treatment
> - **Progress:** Bid submitted, awaiting client
> - **Budget:** $15K pre-bid costs of $120K allocation
> - **Schedule:** Client response expected by March 28
> - **Risk:** LOW
>
> ## Project Gamma — Highway Bridge
> - **Progress:** Tender received, evaluation not started
> - **Next step:** Cost estimation and feasibility review
>
> ## Action Items
> 1. Resolve Alpha steel delay (Procurement)
> 2. Follow up on Beta client response (PM)
> 3. Begin Gamma tender evaluation (Full team)

Director: "Send this to the management team"

Brain drafts email to `outbox/weekly-report-email.md` with report as attachment reference.

---

### Scene 5: New Bid Evaluation (5 min)
**Shows: Full multi-agent pipeline, live dashboard**

> "Start evaluating the Gamma tender for the highway bridge"

**Brain → creates project → full pipeline on /live dashboard:**

The audience watches the live dashboard as agents work:

1. **Cost Controller** reads historical data for similar bridge projects
2. **Procurement Lead** checks vendor availability for concrete, steel, heavy equipment
3. **Safety Officer** reviews highway construction regulations
4. **All three work in parallel** — tool calls streaming on live dashboard

**CHECKPOINT:** Director reviews preliminary findings.

5. **Proposal Writer** assembles initial feasibility summary
6. **Brain:** "Preliminary evaluation complete. Estimated cost: $3.2M. 2 vendors available. Compliance: 3 items to address. Should I proceed to full bid preparation?"

*This is the wow moment — the audience sees 4 agents working simultaneously, each with their own tool calls visible, producing a coordinated result in minutes.*

---

### Scene 6: HR Quick Task (2 min)
**Shows: Quick delegation, non-engineering workflow**

> "Review the site supervisor candidates from HR"

**Brain → delegates to HR Screener (oneshot):**
- Reads `employees.json` for job requirements
- Reads `inbox/email_005.json` for candidate details
- Returns:

> **Site Supervisor Candidates:**
>
> | Candidate | Experience | Certifications | Rating |
> |-----------|-----------|----------------|--------|
> | Raj Kumar | 12 years, 4 bridge projects | PE, PMP, OSHA-30 | ⭐⭐⭐⭐⭐ |
> | Sarah Chen | 8 years, 2 highway projects | PE, LEED | ⭐⭐⭐⭐ |
> | Mike Torres | 15 years, 6 projects | PMP, OSHA-30 | ⭐⭐⭐⭐ |
>
> **Recommendation:** Raj Kumar — best fit for bridge project scope. Schedule interview?

---

---

## Email System (Approach B — Built-in)

### Architecture

```
┌─────────────────────────────────────────────────────────┐
│                  EMAIL SIMULATOR UI                       │
│            /email-sim (separate page)                     │
│                                                          │
│  [Vendor Delay] [Invoice] [Job Application] [Client]     │
│  [Impersonation Attempt] [Internal Request]               │
│                                                          │
│  From: ___________  Subject: ___________                 │
│  Body: ________________________________                  │
│                          [Send to Agent House]            │
└──────────────────────────┬──────────────────────────────┘
                           │ POST /api/email/receive
┌──────────────────────────▼──────────────────────────────┐
│                  EMAIL ENGINE                             │
│                                                          │
│  1. Store in inbox.json                                  │
│  2. Classify: vendor | job_application | client | unknown│
│  3. TRUST CHECK (vendor emails):                         │
│     · Is sender in vendor.trusted_emails?                │
│     · Domain match vendor.domain?                        │
│     · Historical communication exists?                    │
│     → Trusted: process normally                          │
│     → Untrusted: ALERT user + show reasoning             │
│  4. JOB APPLICATION: parse & store in applicants.json    │
│  5. Route to Brain → delegate to appropriate agent       │
└──────────────────────────┬──────────────────────────────┘
                           │
┌──────────────────────────▼──────────────────────────────┐
│                  AGENT PROCESSING                        │
│                                                          │
│  Procurement → vendor emails, quotes, delays             │
│  Cost Controller → invoices                              │
│  HR Screener → job applications                          │
│  Report Generator → client update requests               │
│  Brain → everything else (routes intelligently)          │
└─────────────────────────────────────────────────────────┘
```

### Email Data Model

```json
// inbox/{id}.json
{
  "id": "email_001",
  "from": "ahmed.rahman@vendorxyz.com",
  "from_name": "Ahmed Rahman",
  "to": "projects@epc-demo.com",
  "subject": "Steel Delivery Update — Project Alpha",
  "body": "Dear Sir, We regret to inform you that...",
  "date": "2026-03-22T08:30:00Z",
  "read": false,
  "category": "vendor",
  "vendor_id": "vendor_xyz",
  "trust_status": "trusted",
  "trust_reason": "Sender matches vendor trusted email list"
}
```

### Vendor Trust Verification

```json
// vendors.json (each vendor has trusted_emails)
{
  "id": "vendor_xyz",
  "name": "XYZ Steel Corp",
  "domain": "vendorxyz.com",
  "trusted_emails": [
    "ahmed.rahman@vendorxyz.com",
    "sales@vendorxyz.com",
    "accounts@vendorxyz.com"
  ],
  "contact_person": "Ahmed Rahman",
  "contract_active": true
}
```

**Trust Check Logic:**

| Check | Result | Action |
|-------|--------|--------|
| Sender in `trusted_emails` list | ✅ Trusted | Process normally |
| Sender domain matches `vendor.domain` but email not in list | ⚠️ New contact | Alert: "New email from XYZ's domain. Add to trusted?" |
| Sender domain doesn't match ANY vendor | 🔴 Unknown | Alert: "Unknown sender. Not associated with any vendor." |
| Sender claims to be vendor but domain is different | 🚨 Impersonation | Alert: "IMPERSONATION RISK — claims to be XYZ but sending from @gmail.com" |
| Email mentions invoice/payment but sender untrusted | 🚨 High Risk | Alert: "PAYMENT REQUEST from untrusted source. Possible BEC fraud." |

**Alert Format in Brain:**
```
🚨 TRUST ALERT — Incoming Email

From: ahmed.r@gmail.com
Claims: Vendor XYZ Steel Corp
Subject: Urgent — Updated Payment Details

⚠️ RISK FACTORS:
1. Sender domain (gmail.com) does not match vendor domain (vendorxyz.com)
2. Known contact is ahmed.rahman@vendorxyz.com — this is ahmed.r@gmail.com
3. Email requests payment detail changes — common BEC tactic

RECOMMENDATION: DO NOT process. Verify with vendor via known phone number.

Actions:
  [Add to Trusted List]  [Reject & Flag]  [Verify Manually]
```

### Job Application Storage

When an email is classified as a job application:

```json
// applicants.json (auto-built from emails)
[
  {
    "id": "app_001",
    "name": "Raj Kumar",
    "email": "raj.kumar@email.com",
    "applied_for": "Site Supervisor",
    "experience_years": 12,
    "key_skills": ["bridge construction", "team management", "OSHA-30"],
    "certifications": ["PE", "PMP"],
    "resume_email_id": "email_005",
    "received_date": "2026-03-20",
    "status": "new"
  }
]
```

**When HR has a requirement:**
> "We need a site supervisor for the Gamma bridge project"

Brain → searches `applicants.json` → returns matches ranked by relevance.

### Email Simulator UI (`/email-sim`)

A standalone page for the demo. Pre-built scenarios with one-click send:

```
┌─────────────────────────────────────────────────────────┐
│  📧 Email Simulator                       [Agent House]  │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  SCENARIOS:                                              │
│                                                          │
│  [📦 Vendor Delay]  [💰 Invoice]  [📋 Job Application]  │
│  [👤 Client Query]  [🚨 Impersonation]  [📝 Internal]   │
│                                                          │
│  ─────────────────────────────────────────────────────  │
│                                                          │
│  From:    [ahmed.rahman@vendorxyz.com          ]         │
│  Name:    [Ahmed Rahman                        ]         │
│  To:      [projects@epc-demo.com               ]         │
│  Subject: [Steel Delivery Update               ]         │
│                                                          │
│  Body:                                                   │
│  ┌───────────────────────────────────────────────────┐  │
│  │ Dear Sir,                                         │  │
│  │                                                    │  │
│  │ We regret to inform you that the structural steel │  │
│  │ delivery for Project Alpha will be delayed by     │  │
│  │ approximately 2 weeks due to...                   │  │
│  └───────────────────────────────────────────────────┘  │
│                                                          │
│               [Send to Agent House →]                    │
│                                                          │
│  ─────────────────────────────────────────────────────  │
│  SENT LOG:                                               │
│  ✅ 10:30 — Vendor delay email sent                     │
│  ✅ 10:31 — Agent House processing...                   │
│  ✅ 10:32 — Procurement agent drafted reply              │
└─────────────────────────────────────────────────────────┘
```

**Scenarios:**

| Button | From | Subject | Category | Trust |
|--------|------|---------|----------|-------|
| 📦 Vendor Delay | ahmed.rahman@vendorxyz.com | Steel Delivery Delayed | vendor | ✅ Trusted |
| 💰 Invoice | accounts@pqrconcrete.com | Invoice #1042 — Concrete Supply | vendor/invoice | ✅ Trusted |
| 📋 Job Application | raj.kumar@email.com | Application — Site Supervisor | job_application | N/A |
| 👤 Client Query | sarah.jones@clientabc.com | Phase 2 Update Request | client | ✅ Known |
| 🚨 Impersonation | ahmed.r@gmail.com | URGENT — Updated Bank Details | vendor/payment | 🚨 ALERT |
| 📝 Internal | site.engineer@company.com | Additional Resources Needed | internal | ✅ Internal |

**Impersonation scenario is the wow moment:**
- User clicks "Impersonation" button
- Email arrives claiming to be Vendor XYZ but from a Gmail address
- Brain immediately alerts: domain mismatch, BEC risk, payment request from untrusted source
- Director sees the alert, clicks "Reject & Flag"
- Audience: "This catches fraud before it costs us $50K"

### Outbound Email (Resend Integration)

When Director approves a draft reply:

```
POST https://api.resend.com/emails
{
  "from": "projects@epc-demo.com",
  "to": "ahmed.rahman@vendorxyz.com",
  "subject": "Re: Steel Delivery Update — Project Alpha",
  "html": "<p>Dear Mr. Rahman,...</p>"
}
```

**Env var:** `RESEND_API_KEY` — set once, works immediately.

For the demo, we can actually send real emails to a test address to prove it works live.

---

## What We Need to Build

### 1. Email Engine — `/api/email/*` (3 hours)
- `POST /api/email/receive` — Receive email (from simulator or webhook)
- `GET /api/email/inbox` — List inbox with trust status
- `GET /api/email/inbox/{id}` — Read single email
- `POST /api/email/send` — Send via Resend API (needs `RESEND_API_KEY`)
- Trust verification logic (domain match, trusted list, BEC detection)
- Auto-classify: vendor, job_application, client, internal, unknown
- Job application parser → stores in `applicants.json`
- Brain wiring: "check email" → reads inbox, "send that" → sends via Resend

### 2. Email Simulator UI — `/email-sim` (2 hours)
- Standalone page with 6 scenario buttons
- Pre-fills from/to/subject/body per scenario
- User can edit before sending
- Sends to `POST /api/email/receive`
- Sent log shows what happened after each email
- Impersonation scenario demonstrates BEC detection

### 3. Sample Data Files (2 hours)
- `projects.json` — 3 projects with detailed status
- `vendors.json` — 15 vendors with contracts + trusted_emails
- `invoices.json` — 8 invoices (some with discrepancies)
- `costs.json` — Historical cost data
- `employees.json` — Team roster
- `applicants.json` — Pre-seeded with 3 past applicants
- `schedules.json` — Milestones and deadlines

### 4. EPC Agent Tuning (1 hour)
- Tune agent prompts for the demo scenario
- Add trust-checking instructions to Procurement agent
- Add resume parsing instructions to HR Screener
- Wire Brain to check vendor trust on incoming emails

### 5. Resend Integration (30 min)
- Simple HTTP POST to Resend API
- `RESEND_API_KEY` env var
- Template for outbound emails
- Draft review before send (Director approves in UI)

### 6. Demo Script (30 min)
- Exact commands/clicks for each scene
- Expected outputs and talking points
- Impersonation scenario script (the wow moment)
- Fallback plan if something goes wrong

---

## What We DON'T Need

| Skip This | Why |
|-----------|-----|
| Real database (PostgreSQL, SAP) | JSON files are read by agents the same way. Audience understands. |
| Real email (SMTP/IMAP) | File-based inbox/outbox demonstrates the workflow. "In production, this connects to Exchange/Gmail." |
| OCR for invoices | Invoice data is pre-structured in JSON. "In production, OCR feeds this." |
| Real scheduling tool | JSON schedule data. "In production, this connects to Primavera/MS Project." |
| Custom MCP servers | Agents use Read/Write tools on JSON files. Same pattern, simpler setup. |

**Key line for the demo:** "Everything you're seeing uses the same tool interface — Read, Write, Search. In production, we swap these JSON files for MCP connections to your SAP, Primavera, and Exchange servers. Zero agent code changes."

---

## Success Criteria

The demo is successful if the audience:

1. **Gets the Brain concept** — "I just talk to it and it handles everything"
2. **Sees multi-agent value** — "5 specialists working in parallel is faster than 1 person"
3. **Trusts the checkpoint system** — "Nothing happens without my approval"
4. **Understands extensibility** — "We can add our own agents and connect our systems"
5. **Asks about deployment** — "When can we try this?"

---

## Timeline

| Day | Task | Hours |
|-----|------|-------|
| Day 1 | Email engine (receive, inbox, send, trust check, classify) | 3h |
| Day 1 | Email simulator UI with 6 scenario buttons | 2h |
| Day 2 | Sample data files (projects, vendors, invoices, applicants) | 2h |
| Day 2 | EPC agent prompt tuning + Brain wiring for email | 2h |
| Day 2 | Resend integration for real outbound email | 0.5h |
| Day 3 | End-to-end test of all 6 scenes + impersonation scenario | 2h |
| Day 3 | Fix issues, polish Brain responses | 1.5h |
| Day 3 | Write demo script with talking points | 1h |
| Day 3 | Rehearsal run | 1h |
| **Total** | | **15h** |
