# Agent House — Enterprise AI Workforce Platform

**Version:** 2.0
**Date:** 2026-03-22
**Status:** Strategic Plan
**Vision:** Any team of AI specialists, any industry, any workflow — with human oversight at every critical gate.

---

## Design Process: 5 Iterations of Self-Contradiction

> This plan was developed through iterative refinement, where each version was challenged and contradicted to arrive at something creative yet grounded.

**Iteration 1 — "Just add industry templates"**
→ Contradicted: Every AI company pitches templates. Agents that write markdown about hospital schedules aren't useful. They need to query ACTUAL systems.

**Iteration 2 — "It's all about tool integrations"**
→ Contradicted: If it's just API wrappers, why agents? The reasoning layer that handles "figure out why the Chennai project is over budget" across 5 systems — that's the value.

**Iteration 3 — "Build reasoning + data access for every industry"**
→ Contradicted: Enterprises won't trust AI with SAP/EMR access without SOC2, HIPAA, audit trails. Building compliance for every industry is multi-year, multi-million.

**Iteration 4 — "Start narrow, go deep in one vertical"**
→ Contradicted: Too narrow = miss the market window. OpenClaw went viral by being general. Need to show the vision while building depth.

**Iteration 5 — The synthesis:**
→ Build the PLATFORM that makes any industry deployable. Don't build industry solutions — build the tool that lets industry experts configure their own AI teams. Ship with software dev as the proven template. Let early adopters in EPC/healthcare/HR configure their own.

---

## Core Thesis

**Every industry runs on the same pattern:**

```
Project submitted → Specialists research → Plan created → Leader approves
→ Executors deliver → Quality reviewed → Delivered
```

| Industry | "Project" | Specialists | Quality Gate |
|----------|-----------|-------------|-------------|
| Software | Feature request | PM, Designer, Dev, Security | Code review, QA |
| EPC | Construction bid | Cost engineer, Procurement, Safety | Director review |
| Hospital | Monthly roster | Scheduler, Analyst, Compliance | Admin approval |
| Legal | Client matter | Associate, Paralegal, Researcher | Partner review |
| Marketing | Campaign | Strategist, Copywriter, Designer | Director review |
| HR | Hiring req | Recruiter, Screener, Coordinator | Manager approval |
| Accounting | Quarter close | Auditor, Analyst, Compliance | CFO sign-off |

Agent House abstracts this pattern into a configurable runtime.

---

## Architecture: The Four Layers

```
┌─────────────────────────────────────────────────────────┐
│                   INTERFACE LAYER                        │
│                                                          │
│  Chat (Web/Slack/Teams/WhatsApp)  ·  Dashboard  ·  API  │
│  Email trigger  ·  Voice  ·  Embedded widget             │
└────────────────────────┬────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────┐
│                    BRAIN LAYER                           │
│                                                          │
│  Understands intent → Checks workspace → Decides action  │
│                                                          │
│  Outcomes:                                               │
│  · Answer directly (simple question)                     │
│  · Route to existing project                             │
│  · Create new project + assemble team                    │
│  · Execute single-agent task (OpenClaw-style)            │
│  · Set up automation (cron/webhook)                      │
│  · Store in knowledge base                               │
│  · Escalate to human                                     │
│                                                          │
│  Context sources:                                        │
│  · /workspace/projects/ — all active projects            │
│  · /workspace/knowledge/ — persistent memory             │
│  · /workspace/inbox/ — pending items                     │
│  · User conversation history                             │
└────────────────────────┬────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────┐
│               WORKSPACE MANAGER                          │
│                                                          │
│  · Project lifecycle (create/archive/resume)             │
│  · Knowledge base (cross-project, persistent)            │
│  · File organization (auto-structured per project type)  │
│  · project.json — metadata, status, team, timeline       │
│  · Daily activity logs                                   │
│  · Cross-project insights ("we solved this before")      │
└────────────────────────┬────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────┐
│              AGENT RUNTIME  ← (built today)              │
│                                                          │
│  · Agent registry (file-based: agent.md + settings.json) │
│  · Hybrid execution (oneshot API + Claude Code sessions)  │
│  · Orchestrator (phased pipeline with QA gates)          │
│  · MCP tool access (DB, email, APIs, storage)            │
│  · Webhooks, cron, task injection                        │
│  · Per-agent cost tracking and observability             │
└─────────────────────────────────────────────────────────┘
```

---

## Industry Examples: How It Works

### Example 1: EPC Company — Construction Bid Management

**The Problem:**
An EPC company receives a tender for a water treatment plant. Today, 6 people spend 3 weeks assembling the bid — querying SAP for historical costs, checking vendor availability, calculating safety requirements, writing the technical proposal.

**Agent House Solution:**

```
agents/
  project-director/
    agent.md          "You evaluate project scope and make go/no-go decisions"
    settings.json     mode: oneshot, tools: [SAP-read, email]

  cost-controller/
    agent.md          "You estimate project costs using historical data"
    settings.json     mode: session, tools: [SAP, Primavera, spreadsheet]

  procurement-lead/
    agent.md          "You assess vendor availability and pricing"
    settings.json     mode: session, tools: [vendor-DB, email, HTTP]

  safety-officer/
    agent.md          "You review compliance requirements and flag risks"
    settings.json     mode: oneshot, tools: [compliance-DB]

  proposal-writer/
    agent.md          "You write technical proposals from team inputs"
    settings.json     mode: session, tools: [document-gen, template-engine]

  qa-reviewer/
    agent.md          "You verify all numbers, references, and compliance"
    settings.json     mode: session, tools: [SAP-read, compliance-DB]
```

**Pipeline:**
```
1. Project Director (triage)
   → "This is a mid-size water treatment project. Run full bid process."

2. Research Phase (parallel)
   → Cost Controller queries SAP for 5 similar past projects
   → Procurement Lead checks vendor rates and availability
   → Safety Officer pulls regulatory requirements

3. Planning Phase (parallel)
   → Cost Controller builds cost estimate spreadsheet
   → Procurement Lead creates vendor shortlist with quotes
   → Safety Officer writes compliance checklist

4. Checkpoint: Project Director reviews estimates
   → "Cost looks high on civil works. Ask cost controller to find alternatives."
   → Inject task → Cost Controller revises

5. Proposal Writing
   → Proposal Writer assembles technical + commercial proposal
   → Incorporates all team outputs

6. QA Review
   → QA Reviewer checks all numbers against SAP data
   → Verifies compliance references
   → QA_APPROVED or QA_REJECTED with specific issues

7. Final Checkpoint: Project Director signs off

Output: Complete bid document, cost breakdown, vendor plan, compliance matrix
```

**Tools needed (via MCP):**
- SAP MCP server (read project history, cost data)
- Primavera P6 MCP server (scheduling)
- Document generation (PDF/Excel from templates)
- Email (vendor communications)
- Compliance database

**Time saved:** 3 weeks → 2 days (with human review checkpoints)

---

### Example 2: Hospital — Monthly Roster & Report Generation

**The Problem:**
A hospital creates monthly duty rosters for 200+ doctors across 15 departments. Today, 3 admin staff spend 5 days negotiating schedules, checking labor compliance, and handling swap requests. Monthly reports (patient volume, resource utilization, quality metrics) take another 3 days.

**Agent House Solution:**

```
agents/
  scheduler/
    agent.md          "You create optimal duty rosters balancing preferences and coverage"
    settings.json     mode: session, tools: [EMR, calendar, spreadsheet]

  analyst/
    agent.md          "You analyze patient data and predict volume patterns"
    settings.json     mode: session, tools: [EMR, analytics-DB, charts]

  compliance-officer/
    agent.md          "You verify schedules meet labor laws and hospital policies"
    settings.json     mode: oneshot, tools: [compliance-rules]

  report-generator/
    agent.md          "You create formatted reports from data analysis"
    settings.json     mode: session, tools: [analytics-DB, document-gen, charts]

  notifier/
    agent.md          "You communicate schedules and updates to staff"
    settings.json     mode: oneshot, tools: [email, slack, SMS]
```

**Roster Workflow:**
```
1. Scheduler queries EMR for doctor availability, leave, preferences
2. Analyst pulls last 3 months patient volume by department
3. Analyst predicts next month's volume → recommends staffing levels
4. Scheduler generates optimized roster matching supply to demand
5. Compliance Officer checks: max hours, rest periods, certification requirements
6. Checkpoint: Admin reviews roster
7. Notifier sends schedule to all staff + calendar invites
```

**Report Workflow:**
```
1. Analyst queries EMR for monthly metrics (admissions, discharges, LOS)
2. Analyst queries resource utilization (bed occupancy, OT utilization)
3. Report Generator creates formatted monthly report (PDF + dashboard)
4. Checkpoint: Medical Director reviews
5. Notifier distributes to department heads
```

**Tools needed (via MCP):**
- EMR/HIS MCP server (patient data, doctor schedules)
- Analytics database (historical metrics)
- Calendar API (Google/Outlook)
- Document generation (PDF reports with charts)
- Email/SMS (staff notifications)

**Compliance note:** All EMR queries are read-only. No patient data is stored in Agent House workspace — only aggregated, de-identified metrics. Audit trail logs every query.

---

### Example 3: HR Department — Hiring Pipeline

```
agents/
  recruiter/
    agent.md          "You source candidates and manage the pipeline"
    settings.json     tools: [ATS, LinkedIn, email]

  screener/
    agent.md          "You evaluate resumes against job criteria"
    settings.json     tools: [ATS-read, document-parser]

  coordinator/
    agent.md          "You schedule interviews and manage logistics"
    settings.json     tools: [calendar, email, ATS]

  onboarding/
    agent.md          "You prepare welcome packages and IT setup"
    settings.json     tools: [HRIS, email, slack, IT-ticketing]
```

**Workflow:**
```
Trigger: New job requisition approved in ATS
→ Recruiter receives notification via webhook
→ Recruiter analyzes role requirements, writes job description
→ Checkpoint: Hiring manager reviews JD
→ Screener evaluates incoming applications (batch, daily)
→ Screener creates shortlist with scoring rationale
→ Checkpoint: Hiring manager selects interview candidates
→ Coordinator schedules interviews across availability
→ [After offer accepted]
→ Onboarding creates IT tickets, welcome email, Day 1 agenda
```

---

### Example 4: Accounting — Quarter Close

```
agents/
  auditor/          "Reconcile accounts, flag discrepancies"
  analyst/          "Generate financial statements and variance analysis"
  compliance/       "Check regulatory requirements (tax, GAAP/IFRS)"
  report-writer/    "Create board-ready financial reports"
```

---

## Platform Features Roadmap

### Phase 1: Foundation (Current — Q2 2026)
**Status: Built**

What exists today:
- [x] 8 specialized dev agents with Claude Code sessions
- [x] Phased pipeline: triage → research → plan → discuss → develop → QA
- [x] Hybrid execution (oneshot API + full sessions)
- [x] MCP tool passthrough for external integrations
- [x] Webhooks, cron scheduler, task injection
- [x] Live activity dashboard with per-agent observability
- [x] Human-in-the-loop checkpoints
- [x] Project detail view with file browser

### Phase 2: Agent Platform (Q3 2026)
**Goal: Any agent, any workflow — no Go code changes**

- [ ] File-based agent definitions (agent.md + settings.json)
- [ ] Agent registry — auto-discovers from /agents/ directory
- [ ] Custom pipeline definitions (team.json per project)
- [ ] Phase-specific mode switching declared in settings.json
- [ ] Per-project agent overrides (.agent-house/agents/)
- [ ] Hot-reload agent definitions on file change
- [ ] Agent creation UI in dashboard

### Phase 3: Brain & Workspace (Q4 2026)
**Goal: From task executor to intelligent assistant**

- [ ] Brain agent — conversation-first interface
- [ ] Auto-project creation (Brain decides when/how to create projects)
- [ ] project.json — rich metadata per project (status, team, timeline, context)
- [ ] Workspace manager — organized folder structure per project type
- [ ] Knowledge base — /workspace/knowledge/ for cross-project memory
- [ ] Conversation history with context carry-over
- [ ] Daily activity summaries (auto-generated)
- [ ] Cross-project search ("how did we handle auth in project X?")

### Phase 4: Enterprise Tool Packs (Q1 2027)
**Goal: Connect to the systems enterprises actually use**

MCP server integrations, grouped by industry need:

**Core Pack (all industries):**
- [ ] Email (SMTP send + IMAP read)
- [ ] Calendar (Google + Outlook)
- [ ] Slack / Microsoft Teams
- [ ] Document generation (PDF, Excel, PowerPoint)
- [ ] Cloud storage (S3, GCS, Google Drive)
- [ ] HTTP client (any REST API)

**Engineering Pack:**
- [ ] PostgreSQL / MySQL / MongoDB
- [ ] GitHub / GitLab
- [ ] Jira / Linear
- [ ] CI/CD (GitHub Actions, Jenkins)
- [ ] Docker / Kubernetes

**Business Pack:**
- [ ] Salesforce / HubSpot (CRM)
- [ ] QuickBooks / Xero (Accounting)
- [ ] SAP (ERP — read-only initially)
- [ ] Stripe (Payments)

**Healthcare Pack:**
- [ ] HL7 FHIR client (standard EMR interface)
- [ ] DICOM viewer (medical imaging metadata)
- [ ] Compliance rules engine

**Operations Pack:**
- [ ] Primavera P6 / MS Project (scheduling)
- [ ] Procurement databases
- [ ] Asset management

### Phase 5: Enterprise Readiness (Q2 2027)
**Goal: Production deployment for real companies**

- [ ] Multi-tenant workspaces (team/department isolation)
- [ ] SSO (SAML/OIDC) + RBAC (who can configure which agents)
- [ ] Audit trail — every agent action, tool call, decision logged
- [ ] Data residency controls (where does data live)
- [ ] Encryption at rest + in transit
- [ ] SOC 2 Type II compliance
- [ ] HIPAA-ready mode (no PII storage, read-only EMR, audit everything)
- [ ] On-premise / VPC deployment option
- [ ] Cost budgets and alerts (per project, per department)
- [ ] SLA monitoring (agent response times, error rates)

### Phase 6: Intelligence Layer (H2 2027)
**Goal: The assistant gets smarter over time**

- [ ] Cross-project learning — "We solved vendor negotiation like this before"
- [ ] Proactive alerts — "Project X hasn't had activity in 3 days"
- [ ] Performance analytics — which agents/workflows produce best results
- [ ] Automatic workflow optimization — adjust pipelines based on outcomes
- [ ] Knowledge graph — relationships between projects, people, decisions
- [ ] Agent marketplace — share/import agent definitions
- [ ] Industry benchmarking — "Your bid process takes 5 days, top firms do 2"

---

## Pricing Model (Projected)

| Tier | Target | Agents | Tool Packs | Price |
|------|--------|--------|------------|-------|
| **Starter** | Individual devs | 8 dev agents | Core | Free (BYOK) |
| **Team** | Small teams (5-20) | 20 custom agents | Core + 1 industry | $99/mo |
| **Business** | Departments (20-100) | Unlimited agents | All packs | $499/mo |
| **Enterprise** | Company-wide | Unlimited + custom | All + on-prem | Custom |

Revenue model: Platform fee + per-agent-execution cost (pass-through LLM costs + margin).

---

## Competitive Positioning

```
                    Simple tasks ←─────────→ Complex projects
                         │                        │
  Single agent ──── OpenClaw ─────────────────────│──
                         │                        │
                         │          Agent House ───│──
  Multi-agent  ──────────│───── CrewAI / AutoGen ──│──
                         │                        │
  Enterprise  ───────────│──────────── Salesforce AgentForce
                         │                        │
                    No oversight            Human-in-the-loop
```

**Agent House's unique quadrant:** Multi-agent + Human-in-the-loop + Enterprise-ready.

- **vs OpenClaw:** We're teams, not individuals. We have approval gates, not just chat.
- **vs CrewAI/AutoGen:** We have a production UI, not just a Python framework.
- **vs Salesforce AgentForce:** We're open, configurable, model-agnostic, not locked to one vendor.
- **vs Devin/OpenHands:** We're multi-industry, not dev-only.

---

## What Makes This Defensible

1. **Multi-agent orchestration with quality gates** — hardest part to replicate, and it's the thing enterprises need most (accountability, traceability, human control)

2. **Industry-specific agent configurations** — once an EPC company configures their procurement workflow with 6 agents, they won't switch. High switching cost.

3. **Knowledge base accumulation** — the workspace gets smarter over time. "We negotiated with Vendor X 6 months ago, here's what worked." This context is irreplaceable.

4. **Compliance-first architecture** — audit trails, RBAC, data residency. Enterprise buyers require this. Most AI tools bolt it on later.

5. **Model agnostic** — not tied to one LLM provider. Oneshot agents can use MiniMax, session agents use Claude Code, future agents could use local models. Enterprises want choice.

---

## Immediate Next Steps (This Week)

1. **Move agents to file-based definitions** — agent.md + settings.json (foundation for everything)
2. **Add project.json** — auto-created per project with metadata
3. **Build the Brain agent** — conversation interface that routes to orchestrator
4. **Perfect the dev workflow** — fix remaining edge cases, polish dashboard
5. **Write 3 industry demo configs** — EPC, Hospital, HR (just agent definitions, no custom tools yet)

---

## The One-Liner

**Agent House: Deploy an AI team for any business function. Configure the agents, connect the tools, define the workflow. Humans approve at every gate.**
