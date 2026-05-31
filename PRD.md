# Agent House — Product Requirements Document

**Version:** 2.0
**Date:** 2026-03-18
**Status:** Approved for Engineering
**Audience:** Engineering, Product, Sales, Investors

---

## Table of Contents

1. [Executive Summary](#1-executive-summary)
2. [Problem Statement & Market Opportunity](#2-problem-statement--market-opportunity)
3. [Product Vision & Strategy](#3-product-vision--strategy)
4. [User Personas](#4-user-personas)
5. [Use Cases](#5-use-cases)
6. [Feature Requirements](#6-feature-requirements)
7. [System Architecture](#7-system-architecture)
8. [Data Model](#8-data-model)
9. [API Design](#9-api-design)
10. [Security & Compliance](#10-security--compliance)
11. [Success Metrics & KPIs](#11-success-metrics--kpis)
12. [Competitive Analysis](#12-competitive-analysis)
13. [Go-to-Market Strategy](#13-go-to-market-strategy)
14. [Roadmap](#14-roadmap)
15. [Risks & Mitigations](#15-risks--mitigations)
16. [Appendix: Industry Agent Configurations](#16-appendix-industry-agent-configurations)

---

## 1. Executive Summary

Agent House is an AI multi-agent orchestration platform that deploys specialized AI agent teams to execute complex, multi-phase knowledge work across industries. The platform currently operates a proven 8-agent software development pipeline — CEO, PM, UX, UI, Security, Architect, Senior Dev, Junior Dev — that autonomously takes a task from brief to shipped code through structured Research → Planning → Discussion → Development → QA phases, with human-in-the-loop checkpoints at each gate.

This document defines the expansion of Agent House from a software development tool into a **universal AI workforce platform** applicable to any industry where teams of specialists collaborate on projects: manufacturing, engineering & construction, marketing agencies, legal firms, and beyond.

The core insight is that **every industry runs on the same underlying pattern**: a project is submitted, specialists research and plan it, a leader facilitates alignment, executors deliver the work, and quality gatekeepers verify it. Agent House abstracts this pattern into a configurable runtime and ships industry-specific team templates on top.

**Target Outcomes for This PRD:**
- Guide 6+ months of engineering work across 4 engineering squads
- Define the product sufficiently to close Series A investment
- Enable sales to qualify and close enterprise contracts in 5 verticals
- Provide compliance and security posture for enterprise procurement reviews

**Current State (v1.0):** Single-organization, single-team, software development vertical, file-system persistence, single-user.

**Target State (v2.0):** Multi-organization SaaS, multi-team per org, 5 industry verticals, PostgreSQL persistence, RBAC, external integrations (email, ERP, CRM, file storage), full audit trail.

---

## 2. Problem Statement & Market Opportunity

### 2.1 The Problem

Modern organizations face a compounding crisis: knowledge work is exploding in volume and complexity while skilled talent is expensive, scarce, and difficult to coordinate. A mid-size company executing a marketing campaign, a construction project, or a software release must orchestrate 6–12 specialists across weeks or months. The coordination overhead alone — meetings, handoffs, reviews, rework — consumes 30–40% of total project effort (McKinsey, 2025).

Existing solutions fail in different ways:

- **Point AI tools** (Copilot, Jasper, Midjourney) augment individual contributors but do not orchestrate teams. Each person still works in isolation.
- **Project Management Software** (Jira, Asana, Monday) provides coordination rails but requires humans to do all the work.
- **Chatbot platforms** (ChatGPT, Gemini) are general-purpose single agents. They have no memory, no specialization, no workflow structure, and no handoff mechanism.
- **Developer-facing agent frameworks** (LangGraph, CrewAI, AutoGen) require Python engineering expertise to configure and maintain. They have no UI, no enterprise RBAC, no audit trail, and no industry-specific knowledge.

**The gap:** There is no enterprise-grade, UI-first, industry-configured platform where a business user can deploy a coordinated AI team to execute a complex project end-to-end, with human oversight at critical gates.

### 2.2 Market Opportunity

| Segment | TAM 2025 | CAGR | Notes |
|---|---|---|---|
| AI Automation Platforms | $18.5B | 38% | Gartner, 2025 |
| Enterprise Workflow Software | $26.3B | 22% | IDC, 2025 |
| AI Professional Services | $8.1B | 44% | McKinsey, 2025 |
| **Addressable (AI Workforce)** | **$4.2B** | **52%** | Agent House's initial SAM |

The serviceable addressable market (SAM) for AI workforce orchestration — platforms that replace or augment entire specialized teams — is estimated at $4.2B in 2026 and growing at 52% CAGR. This segment barely existed in 2024; the rapid improvement of frontier language models makes it viable at scale for the first time.

**Beachhead industries by deal size:**

1. **Software Development Teams** — Immediate fit, existing product. Average deal: $2,400/yr per team.
2. **Marketing Agencies** — High task volume, measurable output (content, campaigns). Average deal: $18,000/yr per agency.
3. **Legal Firms** — High value per document, strong audit requirement. Average deal: $48,000/yr per firm.
4. **EPC (Engineering, Procurement, Construction)** — Largest project budgets, longest cycles. Average deal: $120,000/yr per project office.
5. **Manufacturing** — Operational continuity premium, regulatory compliance. Average deal: $84,000/yr per plant.

### 2.3 Why Now

Three converging forces make 2026 the right moment:

1. **Model capability**: GPT-4o, Claude 3.7, Gemini Ultra achieve expert-level reasoning at domain tasks when given structured context and role-specific system prompts.
2. **Enterprise readiness**: Legal, compliance, and procurement teams at Fortune 500 companies have approved AI tooling categories and are actively seeking vendors.
3. **Labor market pressure**: Cost of specialized knowledge workers has increased 34% since 2022 (BLS). AI workforce tools are now economically compelling even at premium SaaS pricing.

---

## 3. Product Vision & Strategy

### 3.1 Vision

> "Every team in every industry can deploy an expert AI workforce that plans, executes, and delivers real work — not just answers."

### 3.2 Strategic Pillars

**Pillar 1: Structured Phase Execution**
Unlike chatbots that operate in a single turn, Agent House decomposes every project into phases (research, planning, discussion, execution, QA), each with specific agents, artifacts, and human checkpoints. This mirrors how real expert teams work and produces verifiable, auditable output.

**Pillar 2: Industry Specialization**
Generic AI produces generic output. Agent House ships with deep industry templates — role-specific system prompts, domain knowledge bases, industry-standard artifacts (HAZOP reports, SOW documents, legal briefs, campaign briefs), and pre-configured toolchains. A manufacturing plant manager gets a team that speaks ISO 9001 and OSHA, not generic project management jargon.

**Pillar 3: Human-in-the-Loop by Design**
AI autonomy without oversight is a liability, not a feature, for enterprise customers. Agent House's checkpoint system makes human review a first-class primitive: every phase gate can require approval, every decision is logged with the deciding party's identity, and any output can be overridden. This is the trust layer that makes enterprise adoption possible.

**Pillar 4: Integration Depth**
An AI team that cannot touch your real systems is a toy. Agent House v2 agents send emails, query your ERP, write to your database, generate PDF reports, and post to your CRM. The platform acts as a secure integration bus between AI agents and enterprise systems.

**Pillar 5: Audit & Compliance**
Every action every agent takes — every file read, every API call, every decision made — is logged with timestamp, agent identity, user identity, and full payload. Exportable as structured JSON or PDF. Data retention policies configurable per organization. SOC 2 Type II roadmap.

### 3.3 Positioning

**Not a chatbot.** Not a copilot. Not a framework. Agent House is an **AI workforce operating system** — the platform that enterprise teams use to deploy, manage, and audit AI agent teams that do real work.

### 3.4 Business Model

| Tier | Price | Included |
|---|---|---|
| **Starter** | $99/mo per team | 1 team, 1 project/month, 5 industry templates, email support |
| **Professional** | $499/mo per team | Unlimited projects, all integrations, custom prompts, RBAC |
| **Enterprise** | Custom ($2K–$15K/mo) | Multi-team, SSO/SAML, custom templates, SLA, dedicated CSM, SOC 2 |
| **Platform** | Custom ($15K+/mo) | White-label, on-prem deployment, custom model routing, professional services |

---

## 4. User Personas

### 4.1 Alex — Engineering Manager (Software Development)

**Background:** Engineering manager at a 120-person SaaS company. Manages 3 development teams. Under pressure to ship faster with frozen headcount.
**Goal:** Use Agent House to handle feature development for lower-priority backlog items without burning senior engineer capacity.
**Key needs:** Code quality gate, GitHub integration, ability to review and override agent decisions, audit trail for engineering leadership.
**Pain point:** Current AI tools give individual engineers copilot-style help but don't coordinate across the full development lifecycle.
**Quote:** "I need a team that can take a ticket from Jira to merged PR, not just autocomplete."

### 4.2 Priya — Operations Director (Manufacturing)

**Background:** Operations director at an automotive parts manufacturer. Oversees production planning, quality, and compliance across 3 shifts.
**Goal:** Use Agent House to run daily production analysis, draft QC reports, identify maintenance risks, and prepare shift handover documentation.
**Key needs:** Integration with SCADA/MES data via API, compliance with ISO 9001 documentation requirements, role-based access so shift leads can't modify quality reports.
**Pain point:** Her team spends 40% of time on documentation and reporting rather than actual operational improvements.
**Quote:** "If your AI team can handle my daily report cycle, I can focus on the floor."

### 4.3 Marcus — Project Director (EPC)

**Background:** Project director at an engineering firm managing a $200M power plant construction project. Coordinates design, procurement, construction, HSE, and QA/QC teams.
**Goal:** Use Agent House to run weekly project reviews, generate progress reports, flag procurement risks, and draft HSE documentation.
**Key needs:** Document management integration (SharePoint), API connection to Oracle Primavera for schedule data, strict access control (subcontractors cannot see commercial terms), full audit trail for regulatory submissions.
**Pain point:** Coordinating 6 specialist teams across 3 time zones with 200+ documents in flight simultaneously.
**Quote:** "I lose a week every month just getting everyone aligned on the same version of the truth."

### 4.4 Sofia — Creative Director (Marketing Agency)

**Background:** Creative director at a boutique digital marketing agency serving 15 B2B clients.
**Goal:** Use Agent House to run campaign strategy, content creation, SEO analysis, and performance reporting for multiple clients simultaneously.
**Key needs:** Integration with Google Analytics, SEMrush, and HubSpot. Client-specific isolation so client A's data never appears in client B's workspace. White-label output.
**Pain point:** Her team is brilliant strategists but drowns in execution — writing blog posts, drafting emails, compiling reports.
**Quote:** "I want to think strategically for clients. The AI team can do the writing."

### 4.5 James — Managing Partner (Legal Firm)

**Background:** Managing partner at a 30-attorney commercial law firm. Practice areas: M&A, contracts, employment law.
**Goal:** Use Agent House to handle first-pass contract review, due diligence research, precedent research, and client intake documentation.
**Key needs:** Strict data isolation (privilege), integration with legal research databases (Westlaw API), full audit trail (attorney-client privilege records), human review required before any client-facing output.
**Pain point:** Junior associate time ($300/hr) is consumed by document review that AI could do faster and cheaper.
**Quote:** "We can't afford the liability of unchecked AI output. But with proper oversight, this changes our economics."

### 4.6 Dev — Platform Administrator (Enterprise IT)

**Background:** Senior IT administrator responsible for deploying and governing Agent House across a 5,000-person enterprise.
**Goal:** Configure SSO, manage user access across departments, set data retention policies, monitor system health, and run compliance reports.
**Key needs:** SAML/OIDC integration, per-organization data isolation, centralized audit log export, API rate limiting per team.
**Pain point:** Shadow AI adoption is happening department by department with no governance. Needs a centralized platform with proper controls.
**Quote:** "I need one platform I can govern, not 40 different AI subscriptions."

---

## 5. Use Cases

### 5.1 Software Development

#### Use Case SD-1: Feature Development from Jira Ticket

**Actor:** Engineering Manager (Alex)
**Trigger:** New Jira ticket labeled `agent-house` is created for "Add two-factor authentication to user login"
**Preconditions:** Project connected to GitHub repo, Jira integration configured, team template: Software Development

**Step-by-Step Agent Interaction:**

1. **Triage (CEO Agent):** Receives task brief. Reads existing codebase via file tools. Identifies the task scope: authentication system, session management, email delivery, UI changes. Decides to run all phases.

2. **Template Selection (Architect + CEO + PM):** Architect proposes `feature-addition` template. CEO and PM review. Checkpoint: human approval required. Alex approves via dashboard.

3. **Research Phase (PM + UX + Security + Architect):**
   - PM Agent researches 2FA UX patterns, writes `.research/pm-research.md` covering TOTP vs SMS vs magic link tradeoffs.
   - UX Agent researches best-practice authentication flows, writes `.research/ux-research.md` with screen flow diagrams.
   - Security Agent researches OWASP authentication guidelines, writes `.research/security-research.md` with threat model.
   - Architect Agent reviews existing auth code, writes `.research/architect-research.md` with integration points.

4. **Planning Phase (PM + UX + UI + Security + Architect):**
   - PM: User stories, acceptance criteria
   - UX: Wireframes described in markdown, user flow
   - UI: Component spec (form design, error states, success states)
   - Security: Security spec (rate limiting, brute force protection, token expiry)
   - Architect: Technical spec (TOTP library selection, database schema changes, API endpoint design)

5. **Checkpoint: Plan Approval.** CEO Agent synthesizes into `development-plan.json`. Alex reviews the plan in the Checkpoint panel. Adds feedback: "Use TOTP only, not SMS. Keep the backup codes flow." CEO agent revises the plan with feedback.

6. **Discussion Phase (CEO moderates):** CEO moderates a 3-round discussion where PM, Security, and Architect debate implementation approach. Output: agreed specification document.

7. **Development Phase (Senior Dev + Junior Dev):**
   - Senior Dev: Core TOTP implementation, database migrations, API endpoints
   - Junior Dev: UI components, email templates, unit tests
   - Both agents produce files into the project directory

8. **QA Phase (CEO reviews):** CEO reviews all produced code. Finds that rate limiting is missing from the API endpoint. Creates revision task for Senior Dev.

9. **Dev Iteration:** Senior Dev adds rate limiting middleware.

10. **Checkpoint: Final Acceptance.** Alex reviews the complete implementation. Approves. Agent House creates a GitHub PR with all changes, linking back to the Jira ticket.

**Artifacts Produced:** `.research/` folder (4 files), `.plans/` folder (spec documents), `.tasks/development-plan.json`, implementation files in project directory, GitHub PR.

**Human Touchpoints:** Template approval, plan approval with feedback, final acceptance.

**Time:** Autonomous execution ~45 minutes. Total elapsed with human review ~2 hours.

---

#### Use Case SD-2: Security Audit of Existing Codebase

**Actor:** Engineering Manager
**Trigger:** Quarterly security review requirement
**Agent Interactions:** Security agent performs SAST-style review of codebase. Architect validates findings. PM translates to business risk. CEO synthesizes into executive security report. Human checkpoint before report is finalized and emailed to CISO.

---

### 5.2 Manufacturing Plant

#### Use Case MFG-1: Daily Production Report Generation

**Actor:** Priya (Operations Director)
**Trigger:** Automated daily schedule trigger at 06:00
**Preconditions:** MES integration configured (REST API to Siemens SIMATIC IT), production targets loaded, shift schedule configured

**Step-by-Step Agent Interaction:**

1. **Triage (Plant Manager Agent):** Receives daily report task. Queries MES API via the external integration system to pull last 24h production data: units produced, downtime events, defect counts, material consumption.

2. **Research Phase:**
   - **Production Supervisor Agent:** Pulls shift-by-shift breakdown. Reads `.data/production_metrics.json` (written by integration layer). Identifies that Shift B produced 12% below target.
   - **QC Agent:** Pulls quality data. Identifies that defect rate on Line 3 is 2.4% vs 1.5% target. Reads `.data/qc_incidents.json`.
   - **Maintenance Engineer Agent:** Queries maintenance log API. Identifies that Line 3 had a 45-minute unplanned downtime at 02:30.
   - **Safety Officer Agent:** Checks incident log. No safety incidents. Reviews near-miss reports.

3. **Planning Phase (Analysis):**
   - Production Supervisor: Root cause analysis for Shift B shortfall
   - QC Agent: Line 3 defect pattern analysis — correlates with the 02:30 downtime
   - Maintenance Engineer: Maintenance action required for Line 3 bearing (flagged in predictive maintenance data)
   - Supply Chain Agent: Checks material inventory levels, flags that Component X has 3-day supply remaining

4. **Discussion Phase (Plant Manager moderates):** Agents align on report narrative. Key findings: Line 3 requires scheduled maintenance within 48h. Component X procurement needs emergency escalation.

5. **Development Phase (Report Generation):**
   - Plant Manager Agent: Writes daily production report in structured format
   - Process Engineer Agent: Generates trend analysis charts (encoded as base64 SVG data)
   - Supply Chain Agent: Drafts emergency procurement request for Component X

6. **Checkpoint: Manager Approval.** Priya reviews the draft report in the dashboard. Adds note: "Escalate the Component X issue to procurement manager directly." Approves.

7. **Delivery:** System sends final PDF report to distribution list via SendGrid integration. Emergency procurement request emailed to procurement manager. Slack notification sent to Priya's phone.

**Artifacts Produced:** `daily-production-report-{date}.pdf`, `procurement-request-{date}.json`, audit log entry for every API call and data access.

---

#### Use Case MFG-2: Incident Investigation Report

**Actor:** QC Manager (Operator role)
**Trigger:** Quality incident flagged — batch rejection rate exceeded 5% threshold
**Agent Interactions:** QC Agent leads investigation. Gathers production data, material traceability data, equipment logs via API. Maintenance Engineer reviews equipment history. Process Engineer analyzes process parameters. HSE (Safety) Agent reviews whether incident triggers regulatory reporting requirement. Human checkpoint before any regulatory notification is sent.

---

#### Use Case MFG-3: Preventive Maintenance Planning

**Actor:** Maintenance Supervisor
**Trigger:** Monthly maintenance planning cycle
**Agent Interactions:** Maintenance Engineer Agent retrieves equipment health data from CMMS API. Process Engineer analyzes impact of proposed maintenance windows on production schedule. Supply Chain Agent checks spare parts inventory. Shift Lead Agent coordinates schedule. Plant Manager Agent approves final schedule. Produces maintenance work orders and schedule in Excel format.

---

### 5.3 EPC (Engineering, Procurement, Construction)

#### Use Case EPC-1: Weekly Project Progress Report

**Actor:** Marcus (Project Director)
**Trigger:** Every Monday 08:00, automated
**Preconditions:** Oracle Primavera integration (schedule data), SAP integration (cost data), SharePoint integration (document management)

**Step-by-Step Agent Interaction:**

1. **Triage (Project Director Agent):** Receives weekly report task. Determines reporting period (previous week). Identifies which project sections need updating.

2. **Research Phase (Data Gathering):**
   - **Planning Engineer Agent:** Queries Primavera API. Extracts: schedule performance index (SPI), critical path activities, 3-week lookahead, milestone status. Finds that Piping Installation milestone is 5 days behind.
   - **Design Engineer Agent:** Reviews document register in SharePoint. Counts IFC (Issued for Construction) drawings vs outstanding. Reports 23 drawings overdue.
   - **Procurement Manager Agent:** Queries SAP Materials Management API. Identifies 3 purchase orders with delivery risk. Flags long-lead item (transformer) as 8 weeks delayed from supplier.
   - **Site Supervisor Agent:** Reads field progress reports uploaded to SharePoint. Extracts manpower, weather delays, RFI status.
   - **HSE Manager Agent:** Reviews incident log. 0 LTI (Lost Time Incidents). 2 near-misses requiring follow-up.
   - **QA/QC Inspector Agent:** Reviews inspection records. 14 NCRs (Non-Conformance Reports) open, 3 new this week.

3. **Planning Phase (Analysis):**
   - Contracts Manager Agent: Identifies that the transformer delay may trigger a variation order from the client. Drafts preliminary variation notice language.
   - Planning Engineer Agent: Revises schedule forecast. New completion date: 12 days later than baseline.
   - Procurement Manager Agent: Identifies alternative transformer supplier. Provides lead time comparison.

4. **Discussion Phase (Project Director moderates):** Project Director Agent moderates alignment on report narrative. Resolved: recommend issuing variation notice. Recommend accelerating civil works to recover 6 days.

5. **Report Generation:** Planning Engineer produces Primavera schedule update. Project Director produces executive summary. All agents populate the standard weekly report template. System uploads report to SharePoint with correct document numbering convention.

6. **Checkpoint: Director Approval.** Marcus reviews. Approves. System distributes report to client via email (with configured recipients and email template). Logs distribution in audit trail.

**Artifacts Produced:** Weekly report PDF (30+ pages), updated Primavera baseline file, variation notice draft, SharePoint upload confirmation.

---

#### Use Case EPC-2: Procurement Risk Assessment

**Actor:** Procurement Manager
**Trigger:** New critical equipment identified in project BOM
**Agent Interactions:** Procurement Manager Agent researches suppliers (via web search tool and internal vendor database API). Design Engineer validates technical specs. HSE Manager reviews import compliance. Contracts Manager reviews standard terms. Planning Engineer assesses lead time impact on schedule. Output: Vendor evaluation matrix, recommended vendor, and purchase order draft.

---

#### Use Case EPC-3: HSE Incident Investigation

**Actor:** HSE Manager
**Trigger:** Reportable incident occurs on site
**Agent Interactions:** HSE Manager Agent leads investigation workflow. Gathers incident data (photos via file upload, witness statements as text input). Site Supervisor Agent provides site conditions context. Contracts Manager Agent reviews contractual notification obligations. Human checkpoint: HSE Director must approve before regulatory notification is submitted. System generates completed incident report in regulatory required format and emails to relevant authorities.

---

### 5.4 Marketing Agency

#### Use Case MKT-1: B2B Content Marketing Campaign

**Actor:** Sofia (Creative Director)
**Trigger:** Client submits new campaign brief for Q2 SaaS product launch
**Preconditions:** Google Analytics integration, SEMrush integration, HubSpot integration, client workspace isolated

**Step-by-Step Agent Interaction:**

1. **Triage (Creative Director Agent):** Reads campaign brief. Identifies deliverables: 8 blog posts, 3 case studies, 12 LinkedIn posts, 1 email nurture sequence (6 emails), SEO keyword strategy, 1 landing page copy.

2. **Research Phase:**
   - **Strategist Agent:** Reads Google Analytics data (via API) for client's existing content performance. Identifies top-converting content types. Queries HubSpot for buyer persona data and deal stage data.
   - **SEO Specialist Agent:** Queries SEMrush API for keyword opportunities. Identifies 24 target keywords with volume and difficulty scores. Maps keywords to buyer journey stages.
   - **Account Manager Agent:** Reviews client's previous campaigns (files in client workspace). Identifies brand voice guidelines. Extracts tone of voice rules.
   - **Analytics Lead Agent:** Builds competitive content analysis using web search tool. Identifies 3 competitor blog strategies.

3. **Planning Phase:**
   - Strategist: Content strategy document (8 pages) — themes, messaging pillars, content mix rationale
   - SEO Specialist: Keyword map — each content piece mapped to primary and secondary keywords, title tags, meta descriptions
   - Copywriter Agent: Content calendar — 90-day publishing schedule with titles and outlines for each piece
   - Designer Agent: Visual brief — image style guide, social graphics specs

4. **Checkpoint: Strategy Approval.** Sofia reviews the strategy. Approves with one note: "The landing page copy should focus on ROI, not features." Creative Director agent updates the brief.

5. **Execution Phase (Content Production):**
   - Copywriter Agent: Produces all 8 blog posts (800–1,200 words each), 6 email sequences, 12 LinkedIn posts
   - SEO Specialist Agent: Optimizes each blog post — adds schema markup recommendations, internal linking suggestions
   - Designer Agent: Produces image briefs for each asset (actual image generation handled via DALL-E integration or delivered as text brief to human designer)
   - Account Manager Agent: Formats content for client delivery, adds brand compliance checklist

6. **QA Phase:** Creative Director Agent reviews all content against brand voice guidelines and campaign strategy. Flags 2 blog posts that need revision (too technical). Sends for Copywriter revision.

7. **Checkpoint: Final Creative Review.** Sofia reviews final content package. Approves. System packages all content into client deliverable folder in Google Drive. Sends client notification email via SendGrid.

**Artifacts Produced:** 8 blog posts (markdown + HTML), 6 email HTML templates, 12 LinkedIn post texts, SEO keyword map (Excel), content calendar (Google Sheets format), landing page copy document. Uploaded to Google Drive, link emailed to client.

---

#### Use Case MKT-2: Monthly Performance Report

**Actor:** Analytics Lead
**Trigger:** 1st of each month, automated
**Agent Interactions:** Analytics Lead Agent queries Google Analytics, HubSpot, and LinkedIn APIs. Generates traffic analysis, lead attribution, conversion analysis. Account Manager Agent writes executive summary. Creative Director reviews before delivery. Delivered as PDF to client and uploaded to client's Drive folder.

---

### 5.5 Legal Firm

#### Use Case LGL-1: Contract Review and Redline

**Actor:** James (Managing Partner)
**Trigger:** Client uploads commercial agreement for review (NDA, SaaS subscription agreement, or M&A purchase agreement)
**Preconditions:** Legal research integration (Westlaw API configured), client matter isolated, privilege rules active

**Step-by-Step Agent Interaction:**

1. **Triage (Managing Partner Agent):** Reviews contract type (SaaS subscription agreement, 42 pages). Identifies governing law (New York), industry (fintech SaaS), client position (customer). Determines review team: Senior Associate, Contract Specialist, Compliance Officer.

2. **Research Phase:**
   - **Research Analyst Agent:** Queries Westlaw API for relevant New York case law on SaaS limitation of liability clauses. Queries firm's internal precedent database (via file tools) for previous similar agreements.
   - **Compliance Officer Agent:** Reviews contract against firm's regulatory checklist for fintech agreements. Checks GDPR/CCPA data processing terms, SOC 2 compliance provisions.
   - **Senior Associate Agent:** Performs initial clause-by-clause review. Tags provisions as: Acceptable / Negotiable / Must Change / Missing.
   - **Contract Specialist Agent:** Reviews defined terms. Flags inconsistencies in capitalization, undefined terms, circular definitions.

3. **Planning Phase (Issues Identification):**
   - Senior Associate: Issues list — 14 issues ranked by risk (Critical / High / Medium / Low). Key critical issues: unlimited liability provision, unilateral price change right, IP ownership ambiguity.
   - Compliance Officer: Compliance issues — GDPR DPA missing, sub-processor list not included.
   - Research Analyst: Relevant precedents for negotiating limitation of liability (3 cases cited).
   - Paralegal Agent: Formats issues list into standard firm template.

4. **Checkpoint: Attorney Review Required.** James reviews the issues list. This is a mandatory human gate — no AI agent may proceed to redlining without attorney approval of the issues analysis. James approves issues list with note: "Focus redlines on the top 5 critical issues only. We can accept the medium ones."

5. **Execution Phase (Redlining):**
   - Contract Specialist Agent: Produces tracked-changes redline of the agreement addressing the 5 approved critical issues. Uses firm's standard fallback language from precedent library.
   - Research Analyst Agent: Drafts comment explanations for each redline (to accompany delivery to opposing counsel).
   - Paralegal Agent: Prepares cover letter to opposing counsel summarizing the redline position.
   - Compliance Officer Agent: Adds GDPR Data Processing Addendum (from firm's standard template library).

6. **QA Phase (Paralegal + Senior Associate review):** Senior Associate Agent reviews redline for legal accuracy. Paralegal verifies formatting and cite accuracy.

7. **Checkpoint: Final Attorney Sign-Off (Mandatory).** James reviews the complete redlined agreement. This checkpoint cannot be auto-delegated — attorney review is legally required. James approves. System generates final delivery package.

8. **Delivery:** System creates client portal folder (via SharePoint integration). Uploads: redlined agreement (DOCX), clean redline PDF, issues memorandum, GDPR DPA. Sends client notification email. Logs all actions in matter audit trail.

**Artifacts Produced:** Redlined agreement (DOCX), issues memorandum, GDPR DPA, cover letter, client delivery email. Full audit trail: who requested, which agent performed each action, which attorney approved, timestamp of every step.

**Critical Rule:** No client-facing output can be delivered without attorney (Manager role) checkpoint approval. This is a hard system constraint, not a user setting.

---

#### Use Case LGL-2: Due Diligence for M&A

**Actor:** Senior Associate
**Trigger:** Client engages firm for sell-side M&A due diligence
**Agent Interactions:** Research Analyst reviews uploaded document tranche (200+ documents). Senior Associate Agent categorizes and prioritizes issues. Paralegal drafts due diligence report sections. Litigation Support Agent reviews pending litigation and regulatory matters. Managing Partner approval required at each report section before inclusion. Final report delivered to client partner.

---

#### Use Case LGL-3: Legal Research Memorandum

**Actor:** Associate
**Trigger:** Partner requests research on a specific legal question
**Agent Interactions:** Research Analyst Agent queries Westlaw API. Synthesizes case law, regulations, and secondary sources. Drafts research memo in firm's standard format. Senior Associate reviews for accuracy. Managing Partner reviews before delivery to client. Full Westlaw query log retained in audit trail.

---

## 6. Feature Requirements

### Priority Legend
- **P0:** Must have for v2.0 launch (blocking)
- **P1:** High priority for v2.0 (important but not blocking)
- **P2:** Planned for v2.1 / Q3

---

### 6.1 Core Platform — Agent Orchestration

| ID | Feature | Priority | Notes |
|---|---|---|---|
| CORE-01 | Multi-phase workflow engine (Triage → Research → Planning → Discussion → Development → QA) | P0 | Exists in v1.0; needs generalization |
| CORE-02 | Configurable phase sequence per industry template | P0 | New: phases should be config-driven, not hardcoded |
| CORE-03 | Agent role system (8+ roles per team) | P0 | Exists; extend Role type to support industry roles |
| CORE-04 | Checkpoint system (pending → approved/rejected/overridden) | P0 | Exists; extend with RBAC rules on who can approve |
| CORE-05 | Auto-delegation timer (configurable per checkpoint type) | P1 | Exists; make per-org configurable |
| CORE-06 | Parallel agent execution within a phase | P0 | Exists (worker pool); generalize |
| CORE-07 | Agent-to-agent messaging with full history | P0 | Exists; needs persistence to PostgreSQL |
| CORE-08 | Human feedback injection into agent prompts | P0 | Exists in v1.0 |
| CORE-09 | Direct chat with individual agents (consultation mode) | P1 | Exists in v1.0 |
| CORE-10 | Agent output artifact management | P0 | Extend to support cloud storage |
| CORE-11 | Multi-project support per organization | P0 | New; currently single project |
| CORE-12 | Concurrent project execution (multiple projects running simultaneously) | P1 | New |
| CORE-13 | Task history and searchable archive | P1 | Exists (HistoryManager); needs DB backend |
| CORE-14 | Subtask decomposition and tracking | P0 | Exists in v1.0 PRODUCT_SPEC |
| CORE-15 | Kanban board view for task management | P0 | Exists in v1.0 PRODUCT_SPEC |

### 6.2 Multi-Industry Teams

| ID | Feature | Priority | Notes |
|---|---|---|---|
| TEAM-01 | Software Development team template (8 roles) | P0 | Exists; formalize as template |
| TEAM-02 | Manufacturing Plant team template (8 roles) | P0 | New |
| TEAM-03 | EPC team template (8 roles) | P0 | New |
| TEAM-04 | Marketing Agency team template (8 roles) | P0 | New |
| TEAM-05 | Legal Firm team template (8 roles) | P0 | New |
| TEAM-06 | Custom role creation (org-defined roles with custom prompts) | P1 | New |
| TEAM-07 | Role-specific system prompt management (editable per org) | P1 | New |
| TEAM-08 | Role-specific tool permissions (which tools each role can use) | P0 | Extend existing permissions system |
| TEAM-09 | Team template marketplace (community + official) | P2 | Future |
| TEAM-10 | Role capability matrix — per-role: tools, file access, API access, phase participation | P0 | New |

### 6.3 External Integrations

| ID | Feature | Priority | Notes |
|---|---|---|---|
| INT-01 | Email send (SMTP, SendGrid, Gmail API) | P0 | New |
| INT-02 | Email receive/read (IMAP, Gmail API) | P1 | New |
| INT-03 | PostgreSQL read/write | P0 | New |
| INT-04 | MySQL read/write | P1 | New |
| INT-05 | MongoDB read/write | P2 | New |
| INT-06 | REST API caller (configurable per integration) | P0 | New; secure credential vault required |
| INT-07 | AWS S3 file storage (read/write) | P0 | New |
| INT-08 | Google Drive integration (read/write) | P0 | New |
| INT-09 | SharePoint / OneDrive integration | P1 | New |
| INT-10 | PDF report generation (wkhtmltopdf or Chromium headless) | P0 | New |
| INT-11 | Excel/CSV report generation | P0 | New |
| INT-12 | Slack notification | P1 | New |
| INT-13 | Microsoft Teams notification | P1 | New |
| INT-14 | Jira integration (read tickets, create tickets, update status) | P1 | New |
| INT-15 | GitHub integration (create PR, post comment) | P1 | New |
| INT-16 | HubSpot CRM integration | P2 | New |
| INT-17 | SAP integration (read-only, via REST API) | P2 | New |
| INT-18 | Oracle Primavera integration | P2 | New |
| INT-19 | Westlaw API integration (legal research) | P2 | New |
| INT-20 | Integration health monitoring and error alerting | P1 | New |
| INT-21 | Credential vault (encrypted storage for API keys, passwords) | P0 | New; required for all integrations |
| INT-22 | Integration audit log (every API call logged) | P0 | New; required for compliance |

### 6.4 Role-Based Access Control

| ID | Feature | Priority | Notes |
|---|---|---|---|
| RBAC-01 | Organization entity — multi-tenant isolation | P0 | New |
| RBAC-02 | Admin role — full platform access, team management | P0 | New |
| RBAC-03 | Manager role — submit tasks, approve checkpoints, view all | P0 | New |
| RBAC-04 | Operator role — submit tasks for their area, view assigned agents | P0 | New |
| RBAC-05 | Viewer role — read-only access to dashboards and reports | P0 | New |
| RBAC-06 | Custom role creation with granular permissions | P1 | New |
| RBAC-07 | SSO/SAML 2.0 integration | P1 | New; required for enterprise |
| RBAC-08 | OIDC / OAuth 2.0 integration | P1 | New |
| RBAC-09 | Per-project access control (project-level visibility rules) | P1 | New |
| RBAC-10 | Checkpoint approval authority — only configured roles can approve specific checkpoint types | P0 | New |
| RBAC-11 | Resource-level permissions (which integrations each role can trigger) | P1 | New |
| RBAC-12 | Invitation and user management API | P0 | New |

### 6.5 Configurable Workflows

| ID | Feature | Priority | Notes |
|---|---|---|---|
| WF-01 | Workflow template YAML format — define phases, agents, checkpoints | P0 | New |
| WF-02 | Phase-level configuration — which agents participate, required outputs | P0 | New |
| WF-03 | Checkpoint configuration — type, required approver role, auto-delegate timeout | P0 | New |
| WF-04 | Custom phase names and descriptions | P1 | New |
| WF-05 | Conditional phases — skip phases based on task triage output | P1 | Exists partially in triage.go |
| WF-06 | Workflow template sharing — export/import YAML | P1 | New |
| WF-07 | Workflow template versioning | P2 | New |
| WF-08 | A/B workflow testing | P2 | New |
| WF-09 | Workflow execution history with per-phase timing metrics | P1 | New |

### 6.6 Audit Trail & Compliance

| ID | Feature | Priority | Notes |
|---|---|---|---|
| AUD-01 | Agent action log — every agent turn logged with prompt, response, tools used | P0 | New; store in PostgreSQL |
| AUD-02 | Decision log — every checkpoint decision: who decided, when, what action | P0 | Extends checkpoint.go |
| AUD-03 | Integration access log — every external API call, file read/write | P0 | New |
| AUD-04 | User action log — every user action (login, task submit, approval, override) | P0 | New |
| AUD-05 | Audit log search and filter UI | P1 | New |
| AUD-06 | Audit log export (JSON, CSV, PDF) | P0 | New |
| AUD-07 | Data retention policy — configurable per org (30d, 90d, 1yr, indefinite) | P1 | New |
| AUD-08 | Immutable audit log — logs cannot be modified or deleted within retention period | P0 | New; append-only table |
| AUD-09 | Compliance report templates (SOC 2, ISO 9001, legal privilege log) | P2 | New |
| AUD-10 | PII data handling — flag and redact PII in logs per org policy | P2 | New |

### 6.7 Dashboard & Monitoring

| ID | Feature | Priority | Notes |
|---|---|---|---|
| DASH-01 | Mission Control canvas — real-time agent topology with status | P0 | Exists in v1.0 |
| DASH-02 | Kanban task board | P0 | Exists in v1.0 PRODUCT_SPEC |
| DASH-03 | Timeline view of project phases | P1 | Exists in v1.0 PRODUCT_SPEC |
| DASH-04 | Agent activity feed (real-time WebSocket) | P0 | Exists in v1.0 |
| DASH-05 | KPI panel — tasks completed, time per phase, agent utilization | P0 | Extends existing KPI strip |
| DASH-06 | Multi-project dashboard — status of all active projects | P1 | New |
| DASH-07 | Organization analytics — usage, cost, ROI metrics | P1 | New |
| DASH-08 | Integration health panel | P1 | New |
| DASH-09 | System health monitoring (agent process status, queue depth) | P1 | New |
| DASH-10 | Mobile-responsive dashboard | P2 | New |

---

## 7. System Architecture

### 7.1 High-Level Component Diagram

```
┌─────────────────────────────────────────────────────────────────────┐
│                         CLIENT LAYER                                │
│   Browser (Vanilla JS/CSS/HTML)  ←→  WebSocket Hub  ←→  REST API  │
└───────────────────────────┬─────────────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────────────┐
│                       WEB SERVER (Go)                               │
│  server.go  │  dashboard.go  │  mission.go  │  websocket.go         │
│  checkpoint_handlers.go  │  kanban_handlers.go  │  chat_handlers.go │
│  rbac_middleware.go  │  audit_middleware.go                         │
└───────────────────────────┬─────────────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────────────┐
│                    ORCHESTRATION LAYER (Go)                         │
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │                    Orchestrator                              │   │
│  │  orchestrator.go │ phases.go │ phase_executor.go           │   │
│  │  triage.go │ dev_plan.go │ qa_phase.go │ subtask.go        │   │
│  └──────────────────────────┬──────────────────────────────────┘   │
│                             │                                       │
│  ┌──────────────────────────▼──────────────────────────────────┐   │
│  │                   Agent Pool                                 │   │
│  │  agent.go │ hierarchy.go │ permissions.go │ file_ops.go     │   │
│  │  [CEO] [PM] [UX] [UI] [Security] [Architect] [Dev×2]        │   │
│  │  [PlantMgr] [ProdSup] [QC] [Maint] [Safety] [Supply] ...   │   │
│  └──────────────────────────┬──────────────────────────────────┘   │
│                             │                                       │
│  ┌──────────────────────────▼──────────────────────────────────┐   │
│  │               Support Systems                                │   │
│  │  checkpoint/  │  kanban/  │  message/  │  task/             │   │
│  │  chat/  │  template/  │  skill/  │  agentask/               │   │
│  └──────────────────────────┬──────────────────────────────────┘   │
└───────────────────────────┬─┘─────────────────────────────────────-┘
                            │
┌───────────────────────────▼─────────────────────────────────────────┐
│                    INTEGRATION LAYER (New in v2)                    │
│                                                                     │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐  │
│  │  Email      │ │  Database   │ │  File Store │ │  REST API   │  │
│  │  (SMTP/     │ │  Connector  │ │  Connector  │ │  Connector  │  │
│  │  SendGrid/  │ │  (PG/MySQL/ │ │  (S3/Drive/ │ │  (Generic + │  │
│  │  Gmail API) │ │  Mongo)     │ │  SharePoint)│ │  per-vendor)│  │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘  │
│                                                                     │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐                  │
│  │  Report     │ │  Credential │ │  Integration│                  │
│  │  Generator  │ │  Vault      │ │  Audit Log  │                  │
│  │  (PDF/Excel)│ │  (AES-256)  │ │             │                  │
│  └─────────────┘ └─────────────┘ └─────────────┘                  │
└───────────────────────────┬─────────────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────────────┐
│                    DATA LAYER                                       │
│                                                                     │
│  PostgreSQL (primary)  │  Redis (session cache, message queue)      │
│  File System (agent working dirs, artifacts)                        │
└─────────────────────────────────────────────────────────────────────┘
```

### 7.2 Data Flow — Task Execution

```
User Submits Task
      │
      ▼
[REST POST /api/projects/{id}/tasks]
      │ RBAC check: Manager or Operator role required
      ▼
[Task stored in PostgreSQL, TaskID generated]
      │
      ▼
[Orchestrator.Run(task)] ← goroutine spawned
      │
      ├─ Phase: Triage
      │    └─ CEO Agent reads task, project context
      │    └─ Decides: which phases to run (may skip some)
      │    └─ Checkpoint emitted if triage requires approval
      │
      ├─ Phase: Template Selection (if enabled)
      │    └─ Architect Agent proposes template
      │    └─ CEO + PM review
      │    └─ Checkpoint: template_approval → pauses, waits for human
      │         └─ WebSocket: {type: "checkpoint_pending", ...}
      │         └─ Human approves via REST POST /api/checkpoints/{id}/decide
      │         └─ Orchestrator resumes
      │
      ├─ Phase: Research
      │    └─ Worker Pool: N agents run in parallel goroutines
      │    └─ Each agent: reads files, calls external APIs (via Integration Layer)
      │    └─ Each agent writes output artifact to project dir
      │    └─ Audit log: every agent action, every integration call
      │    └─ WebSocket: {type: "agent_update", role, status, message}
      │
      ├─ Phase: Planning → Discussion → Development → QA
      │    └─ (Same pattern, different agents and checkpoint types)
      │
      └─ Task Complete
           └─ Final artifacts stored
           └─ Integration Layer: send email / upload to cloud storage
           └─ WebSocket: {type: "task_complete", taskId, artifacts}
           └─ Audit log: task completion record
```

### 7.3 Multi-Tenancy Architecture

Agent House v2 uses a **shared database, schema-separated** multi-tenant model:

- Each `organization_id` is a UUID that is present as a partition key on every data table
- Row-Level Security (PostgreSQL RLS) enforces org isolation at the database level
- No cross-org data leakage is possible even in the application layer
- Agent working directories are isolated at the filesystem level: `/data/orgs/{org_id}/projects/{project_id}/`
- Credentials stored in the vault are encrypted per-org with an org-specific key derived from the master key

### 7.4 Agent Execution Model

Each agent runs as a managed subprocess (Claude CLI via PTY, as in v1.0). In v2, agents are additionally:

- **Bound to a project context:** Agent's working directory is scoped to the project
- **Given an integration tool set:** Based on role permissions, agent receives tool definitions for permitted integrations
- **Audited:** Every agent turn (prompt + response) is captured and written to the audit log
- **Resource-limited:** Max tokens per turn, max turns per phase, max wall-clock time per phase — all configurable per org

---

## 8. Data Model

### 8.1 PostgreSQL Schema

```sql
-- Organizations (multi-tenant root)
CREATE TABLE organizations (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL,
    slug            VARCHAR(100) NOT NULL UNIQUE,
    plan            VARCHAR(50) NOT NULL DEFAULT 'starter', -- starter|professional|enterprise
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    settings        JSONB NOT NULL DEFAULT '{}'
);

-- Users
CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id          UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    email           VARCHAR(255) NOT NULL,
    name            VARCHAR(255) NOT NULL,
    password_hash   VARCHAR(255),                -- null for SSO-only users
    role            VARCHAR(50) NOT NULL,        -- admin|manager|operator|viewer
    status          VARCHAR(50) NOT NULL DEFAULT 'active', -- active|suspended|pending
    sso_provider    VARCHAR(100),
    sso_subject     VARCHAR(255),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_login_at   TIMESTAMPTZ,
    UNIQUE(org_id, email)
);

CREATE INDEX idx_users_org ON users(org_id);
CREATE INDEX idx_users_email ON users(email);

-- Sessions
CREATE TABLE sessions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash      VARCHAR(255) NOT NULL UNIQUE,
    expires_at      TIMESTAMPTZ NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    user_agent      TEXT,
    ip_address      INET
);

-- Team Templates
CREATE TABLE team_templates (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id          UUID REFERENCES organizations(id) ON DELETE CASCADE, -- null = system template
    name            VARCHAR(255) NOT NULL,
    industry        VARCHAR(100) NOT NULL, -- software|manufacturing|epc|marketing|legal|custom
    description     TEXT,
    is_system       BOOLEAN NOT NULL DEFAULT FALSE,
    config          JSONB NOT NULL,  -- full WorkflowConfig as JSON
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    version         INTEGER NOT NULL DEFAULT 1
);

-- Projects
CREATE TABLE projects (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id          UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    template_id     UUID REFERENCES team_templates(id),
    status          VARCHAR(50) NOT NULL DEFAULT 'active', -- active|archived|deleted
    working_dir     VARCHAR(500) NOT NULL,
    settings        JSONB NOT NULL DEFAULT '{}',
    created_by      UUID NOT NULL REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_projects_org ON projects(org_id);

-- Project Members (per-project access control)
CREATE TABLE project_members (
    project_id      UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role            VARCHAR(50) NOT NULL, -- project-level override
    added_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY(project_id, user_id)
);

-- Tasks (top-level work items submitted to orchestrator)
CREATE TABLE tasks (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id          UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    project_id      UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    title           VARCHAR(500) NOT NULL,
    description     TEXT NOT NULL,
    status          VARCHAR(50) NOT NULL DEFAULT 'pending',
    -- pending|running|checkpoint_pending|completed|failed|cancelled
    priority        VARCHAR(20) NOT NULL DEFAULT 'medium', -- critical|high|medium|low
    source          VARCHAR(50) NOT NULL DEFAULT 'manual', -- manual|api|scheduled|integration
    submitted_by    UUID REFERENCES users(id),
    started_at      TIMESTAMPTZ,
    completed_at    TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata        JSONB NOT NULL DEFAULT '{}',
    result          JSONB                                  -- final output summary
);

CREATE INDEX idx_tasks_org ON tasks(org_id);
CREATE INDEX idx_tasks_project ON tasks(project_id);
CREATE INDEX idx_tasks_status ON tasks(status);

-- Subtasks (granular work items within a task)
CREATE TABLE subtasks (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id         UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    project_id      UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    org_id          UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    external_id     VARCHAR(50) NOT NULL,               -- e.g. "st-4-2"
    title           VARCHAR(500) NOT NULL,
    description     TEXT,
    phase           VARCHAR(100) NOT NULL,
    assigned_role   VARCHAR(100) NOT NULL,
    status          VARCHAR(50) NOT NULL DEFAULT 'pending',
    priority        VARCHAR(20) NOT NULL DEFAULT 'medium',
    estimated_effort VARCHAR(20),                       -- small|medium|large
    dependencies    JSONB NOT NULL DEFAULT '[]',        -- array of subtask external_ids
    outputs         JSONB NOT NULL DEFAULT '[]',        -- expected file paths
    acceptance_criteria TEXT,
    source          VARCHAR(50) NOT NULL DEFAULT 'auto', -- auto|manual
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_subtasks_task ON subtasks(task_id);

-- Checkpoints
CREATE TABLE checkpoints (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id          UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    project_id      UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    task_id         UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    type            VARCHAR(100) NOT NULL,
    -- template_approval|plan_approval|phase_gate|final_acceptance|research_review|spec_review|pre_qa_review
    phase_index     INTEGER,
    status          VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending|resolved
    artifact_path   VARCHAR(500),
    artifact_summary TEXT,
    required_role   VARCHAR(50) NOT NULL DEFAULT 'manager', -- which role can approve
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    resolved_at     TIMESTAMPTZ,
    auto_delegate_at TIMESTAMPTZ                        -- when auto-delegation fires
);

CREATE INDEX idx_checkpoints_task ON checkpoints(task_id);
CREATE INDEX idx_checkpoints_status ON checkpoints(status);

-- Checkpoint Decisions
CREATE TABLE checkpoint_decisions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    checkpoint_id   UUID NOT NULL REFERENCES checkpoints(id) ON DELETE CASCADE,
    action          VARCHAR(50) NOT NULL,               -- approved|rejected|overridden
    feedback        TEXT,
    override_data   TEXT,
    decided_by_user UUID REFERENCES users(id),          -- null if auto-delegated
    decided_by      VARCHAR(100) NOT NULL,              -- "user:{user_id}" or "ceo_auto"
    decided_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Agent Messages
CREATE TABLE agent_messages (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id          UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    project_id      UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    task_id         UUID REFERENCES tasks(id) ON DELETE CASCADE,
    role            VARCHAR(100) NOT NULL,
    direction       VARCHAR(20) NOT NULL,               -- outbound (to AI) | inbound (from AI)
    content         TEXT NOT NULL,
    phase           VARCHAR(100),
    turn_index      INTEGER,
    tokens_used     INTEGER,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_messages_task ON agent_messages(task_id);
CREATE INDEX idx_messages_project ON agent_messages(project_id);

-- Chat Messages (user ↔ individual agent consultation)
CREATE TABLE chat_messages (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id          UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    project_id      UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    agent_role      VARCHAR(100) NOT NULL,
    sender          VARCHAR(50) NOT NULL,               -- "user" | "agent"
    user_id         UUID REFERENCES users(id),
    content         TEXT NOT NULL,
    context         VARCHAR(50) NOT NULL DEFAULT 'consultation', -- consultation|mid_workflow
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Integrations (configured per organization)
CREATE TABLE integrations (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id          UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    type            VARCHAR(100) NOT NULL,
    -- email_smtp|email_sendgrid|email_gmail|db_postgres|db_mysql|
    -- storage_s3|storage_gdrive|storage_sharepoint|api_rest|
    -- report_pdf|slack|teams|jira|github
    name            VARCHAR(255) NOT NULL,              -- human name, e.g. "Production DB"
    status          VARCHAR(50) NOT NULL DEFAULT 'active', -- active|disabled|error
    config          JSONB NOT NULL DEFAULT '{}',        -- non-secret config (host, bucket name, etc.)
    -- Secrets stored in credential vault, referenced by vault_key
    vault_key       VARCHAR(255),
    last_tested_at  TIMESTAMPTZ,
    last_error      TEXT,
    created_by      UUID NOT NULL REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_integrations_org ON integrations(org_id);

-- Integration Audit Log (immutable)
CREATE TABLE integration_audit_log (
    id              BIGSERIAL PRIMARY KEY,
    org_id          UUID NOT NULL,
    integration_id  UUID NOT NULL,
    task_id         UUID,
    agent_role      VARCHAR(100),
    action          VARCHAR(100) NOT NULL,              -- email_sent|db_query|file_read|file_write|api_call
    target          TEXT,                              -- email recipient / table name / file path / URL
    status          VARCHAR(50) NOT NULL,               -- success|failure
    request_summary TEXT,                              -- truncated request (no secrets)
    response_summary TEXT,                             -- truncated response
    duration_ms     INTEGER,
    error           TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- This table is append-only. No UPDATE or DELETE permitted (enforced via triggers).
CREATE INDEX idx_int_audit_org ON integration_audit_log(org_id, created_at DESC);
CREATE INDEX idx_int_audit_task ON integration_audit_log(task_id);

-- Audit Log (immutable — all platform actions)
CREATE TABLE audit_log (
    id              BIGSERIAL PRIMARY KEY,
    org_id          UUID NOT NULL,
    actor_type      VARCHAR(50) NOT NULL,               -- user|agent|system|scheduler
    actor_id        VARCHAR(255) NOT NULL,              -- user UUID or agent role
    action          VARCHAR(255) NOT NULL,
    resource_type   VARCHAR(100),                      -- project|task|checkpoint|integration|user
    resource_id     VARCHAR(255),
    details         JSONB,
    ip_address      INET,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Append-only enforced by trigger
CREATE INDEX idx_audit_org ON audit_log(org_id, created_at DESC);
CREATE INDEX idx_audit_actor ON audit_log(actor_type, actor_id);
CREATE INDEX idx_audit_resource ON audit_log(resource_type, resource_id);

-- Workflow Settings (per project)
CREATE TABLE workflow_settings (
    project_id                  UUID PRIMARY KEY REFERENCES projects(id) ON DELETE CASCADE,
    require_template_approval   BOOLEAN NOT NULL DEFAULT TRUE,
    require_plan_approval       BOOLEAN NOT NULL DEFAULT TRUE,
    require_phase_gate          BOOLEAN NOT NULL DEFAULT FALSE,
    require_final_acceptance    BOOLEAN NOT NULL DEFAULT TRUE,
    require_research_review     BOOLEAN NOT NULL DEFAULT FALSE,
    require_spec_review         BOOLEAN NOT NULL DEFAULT FALSE,
    require_pre_qa_review       BOOLEAN NOT NULL DEFAULT FALSE,
    auto_delegate_minutes       INTEGER NOT NULL DEFAULT 0,    -- 0 = no auto-delegation
    updated_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Kanban Tasks (extends subtasks with board metadata)
CREATE TABLE kanban_tasks (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subtask_id      UUID REFERENCES subtasks(id) ON DELETE CASCADE,
    project_id      UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    org_id          UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    column          VARCHAR(50) NOT NULL DEFAULT 'backlog',
    -- backlog|in_progress|in_review|done
    position        INTEGER NOT NULL DEFAULT 0,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### 8.2 Go Struct Definitions (Key Entities)

```go
// Organization represents a tenant
type Organization struct {
    ID        string            `json:"id" db:"id"`
    Name      string            `json:"name" db:"name"`
    Slug      string            `json:"slug" db:"slug"`
    Plan      string            `json:"plan" db:"plan"`
    CreatedAt time.Time         `json:"created_at" db:"created_at"`
    UpdatedAt time.Time         `json:"updated_at" db:"updated_at"`
    Settings  map[string]any    `json:"settings" db:"settings"`
}

// User represents a platform user
type User struct {
    ID           string     `json:"id" db:"id"`
    OrgID        string     `json:"org_id" db:"org_id"`
    Email        string     `json:"email" db:"email"`
    Name         string     `json:"name" db:"name"`
    PasswordHash string     `json:"-" db:"password_hash"`
    Role         UserRole   `json:"role" db:"role"`
    Status       string     `json:"status" db:"status"`
    SSOProvider  string     `json:"sso_provider,omitempty" db:"sso_provider"`
    SSOSubject   string     `json:"sso_subject,omitempty" db:"sso_subject"`
    CreatedAt    time.Time  `json:"created_at" db:"created_at"`
    LastLoginAt  *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
}

type UserRole string

const (
    RoleAdmin    UserRole = "admin"
    RoleManager  UserRole = "manager"
    RoleOperator UserRole = "operator"
    RoleViewer   UserRole = "viewer"
)

// AgentRoleConfig defines a configurable agent role within a team template
type AgentRoleConfig struct {
    Role         string            `json:"role" yaml:"role"`
    DisplayName  string            `json:"display_name" yaml:"display_name"`
    SystemPrompt string            `json:"system_prompt" yaml:"system_prompt"`
    Tools        []string          `json:"tools" yaml:"tools"`
    // Allowed: email|db_read|db_write|storage_read|storage_write|api_call|web_search|file_read|file_write
    Phases       []string          `json:"phases" yaml:"phases"`
    MaxTurns     int               `json:"max_turns" yaml:"max_turns"`
    MaxTokens    int               `json:"max_tokens" yaml:"max_tokens"`
    Metadata     map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

// WorkflowConfig is the complete configuration for a team template
type WorkflowConfig struct {
    Name        string            `json:"name" yaml:"name"`
    Industry    string            `json:"industry" yaml:"industry"`
    Description string            `json:"description" yaml:"description"`
    Phases      []PhaseDefinition `json:"phases" yaml:"phases"`
    Agents      []AgentRoleConfig `json:"agents" yaml:"agents"`
    Checkpoints []CheckpointRule  `json:"checkpoints" yaml:"checkpoints"`
}

// PhaseDefinition describes a workflow phase
type PhaseDefinition struct {
    ID              string   `json:"id" yaml:"id"`
    Name            string   `json:"name" yaml:"name"`
    Description     string   `json:"description" yaml:"description"`
    ParticipantRoles []string `json:"participant_roles" yaml:"participant_roles"`
    LeaderRole      string   `json:"leader_role,omitempty" yaml:"leader_role,omitempty"`
    Parallel        bool     `json:"parallel" yaml:"parallel"`
    RequiredOutputs []string `json:"required_outputs" yaml:"required_outputs"`
    Conditional     string   `json:"conditional,omitempty" yaml:"conditional,omitempty"`
    // e.g. "triage.decision.skip_research == true"
}

// CheckpointRule defines when a checkpoint fires and who can approve it
type CheckpointRule struct {
    Type          string `json:"type" yaml:"type"`
    AfterPhase    string `json:"after_phase" yaml:"after_phase"`
    RequiredRole  string `json:"required_role" yaml:"required_role"` // manager|admin
    CanAutoDelegate bool `json:"can_auto_delegate" yaml:"can_auto_delegate"`
    MandatoryHuman  bool `json:"mandatory_human" yaml:"mandatory_human"`
    // true = no auto-delegation allowed (e.g., attorney sign-off)
}

// Integration represents a configured external integration
type Integration struct {
    ID          string            `json:"id" db:"id"`
    OrgID       string            `json:"org_id" db:"org_id"`
    Type        IntegrationType   `json:"type" db:"type"`
    Name        string            `json:"name" db:"name"`
    Status      string            `json:"status" db:"status"`
    Config      map[string]any    `json:"config" db:"config"`
    VaultKey    string            `json:"-" db:"vault_key"`
    LastTestedAt *time.Time       `json:"last_tested_at,omitempty" db:"last_tested_at"`
    LastError   string            `json:"last_error,omitempty" db:"last_error"`
    CreatedBy   string            `json:"created_by" db:"created_by"`
    CreatedAt   time.Time         `json:"created_at" db:"created_at"`
}

type IntegrationType string

const (
    IntegrationEmailSMTP      IntegrationType = "email_smtp"
    IntegrationEmailSendGrid  IntegrationType = "email_sendgrid"
    IntegrationEmailGmail     IntegrationType = "email_gmail"
    IntegrationDBPostgres     IntegrationType = "db_postgres"
    IntegrationDBMySQL        IntegrationType = "db_mysql"
    IntegrationStorageS3      IntegrationType = "storage_s3"
    IntegrationStorageGDrive  IntegrationType = "storage_gdrive"
    IntegrationStorageSP      IntegrationType = "storage_sharepoint"
    IntegrationAPIRest        IntegrationType = "api_rest"
    IntegrationReportPDF      IntegrationType = "report_pdf"
    IntegrationSlack          IntegrationType = "slack"
    IntegrationTeams          IntegrationType = "teams"
    IntegrationJira           IntegrationType = "jira"
    IntegrationGitHub         IntegrationType = "github"
)

// AuditLogEntry — immutable record of every platform action
type AuditLogEntry struct {
    ID           int64          `json:"id" db:"id"`
    OrgID        string         `json:"org_id" db:"org_id"`
    ActorType    string         `json:"actor_type" db:"actor_type"` // user|agent|system
    ActorID      string         `json:"actor_id" db:"actor_id"`
    Action       string         `json:"action" db:"action"`
    ResourceType string         `json:"resource_type,omitempty" db:"resource_type"`
    ResourceID   string         `json:"resource_id,omitempty" db:"resource_id"`
    Details      map[string]any `json:"details,omitempty" db:"details"`
    IPAddress    string         `json:"ip_address,omitempty" db:"ip_address"`
    CreatedAt    time.Time      `json:"created_at" db:"created_at"`
}
```

### 8.3 Workflow Template YAML Example (Manufacturing Plant)

```yaml
name: Manufacturing Plant Operations
industry: manufacturing
description: |
  AI agent team for manufacturing plant operations, production planning,
  quality control, and compliance reporting. Configured for ISO 9001 environments.

agents:
  - role: plant_manager
    display_name: Plant Manager
    system_prompt: |
      You are the Plant Manager AI agent for a manufacturing facility.
      Your responsibilities: oversee all plant operations, review production
      performance, make operational decisions, escalate critical issues.
      You have deep knowledge of manufacturing KPIs, OEE (Overall Equipment
      Effectiveness), lean manufacturing principles, and ISO 9001 requirements.
      Always reference specific metrics when making recommendations.
      Your output must be clear, actionable, and suitable for senior management.
    tools: [file_read, file_write, api_call, email, storage_read, storage_write, report_pdf]
    phases: [triage, discussion, qa, report_generation]
    max_turns: 12
    max_tokens: 8000

  - role: production_supervisor
    display_name: Production Supervisor
    system_prompt: |
      You are the Production Supervisor AI agent. Responsibilities:
      analyze shift production data, identify bottlenecks, track OEE per line,
      compare actual vs target output, document shift handover notes.
      You speak in terms of units/hour, downtime %, shift utilization, and cycle time.
      Always identify the root cause of deviations, not just the symptom.
    tools: [file_read, file_write, api_call, storage_read]
    phases: [research, planning]
    max_turns: 8
    max_tokens: 6000

  - role: qc_engineer
    display_name: Quality Control Engineer
    system_prompt: |
      You are the QC Engineer AI agent. Responsibilities: analyze defect data,
      perform SPC (Statistical Process Control) analysis, track NCRs
      (Non-Conformance Reports), prepare quality reports per ISO 9001.
      Use 8D methodology for root cause analysis. Reference control charts,
      Cp, Cpk values when available. Flag any out-of-control conditions immediately.
    tools: [file_read, file_write, api_call, storage_read, storage_write]
    phases: [research, planning, qa]
    max_turns: 10
    max_tokens: 7000

  - role: maintenance_engineer
    display_name: Maintenance Engineer
    system_prompt: |
      You are the Maintenance Engineer AI agent. Responsibilities: review
      equipment health data, identify maintenance requirements, prioritize
      work orders, analyze MTBF/MTTR trends, flag safety-critical equipment
      issues. Reference OEM maintenance manuals and CMMS data when available.
      Use RCM (Reliability Centered Maintenance) principles.
    tools: [file_read, file_write, api_call, storage_read]
    phases: [research, planning]
    max_turns: 8
    max_tokens: 6000

  - role: safety_officer
    display_name: HSE Safety Officer
    system_prompt: |
      You are the HSE Safety Officer AI agent. Responsibilities: review
      incident logs, check compliance with OSHA regulations, assess near-miss
      reports, recommend corrective actions, prepare regulatory submissions.
      Never downplay safety concerns. Flag any condition requiring immediate
      human escalation clearly with [IMMEDIATE ACTION REQUIRED] prefix.
      You are knowledgeable about OSHA 29 CFR 1910, ISO 45001.
    tools: [file_read, file_write, api_call, storage_read, storage_write, email]
    phases: [research, planning, qa]
    max_turns: 8
    max_tokens: 6000

  - role: supply_chain
    display_name: Supply Chain Manager
    system_prompt: |
      You are the Supply Chain Manager AI agent. Responsibilities: monitor
      material inventory levels, identify shortage risks, coordinate
      procurement needs, track supplier performance. Use EOQ, safety stock,
      and reorder point calculations. Flag critical shortages with
      [PROCUREMENT ALERT] prefix.
    tools: [file_read, file_write, api_call, storage_read, email]
    phases: [research, planning]
    max_turns: 6
    max_tokens: 5000

  - role: process_engineer
    display_name: Process Engineer
    system_prompt: |
      You are the Process Engineer AI agent. Responsibilities: analyze
      process parameters, identify process capability issues, recommend
      process improvements, perform FMEA analysis, document process changes.
      Reference SPC data, control limits, and process capability indices.
    tools: [file_read, file_write, api_call, storage_read, storage_write, report_pdf]
    phases: [planning, report_generation]
    max_turns: 8
    max_tokens: 6000

  - role: shift_lead
    display_name: Shift Lead
    system_prompt: |
      You are the Shift Lead AI agent. Responsibilities: compile shift
      handover documentation, track hourly production targets, document
      immediate issues requiring next-shift attention. Keep language
      clear and direct. Use bullet points for action items.
    tools: [file_read, file_write, storage_read]
    phases: [research]
    max_turns: 5
    max_tokens: 4000

phases:
  - id: triage
    name: Triage
    description: Plant Manager reviews task scope and determines required analysis
    participant_roles: [plant_manager]
    parallel: false
    required_outputs: [triage-decision.json]

  - id: research
    name: Data Gathering
    description: All specialist agents query data sources and document findings
    participant_roles: [production_supervisor, qc_engineer, maintenance_engineer, safety_officer, supply_chain, shift_lead]
    parallel: true
    required_outputs:
      - research/production-data.md
      - research/qc-data.md
      - research/maintenance-data.md
      - research/safety-data.md
      - research/supply-chain-data.md

  - id: planning
    name: Analysis & Recommendations
    description: Agents analyze findings and form recommendations
    participant_roles: [production_supervisor, qc_engineer, maintenance_engineer, safety_officer, supply_chain, process_engineer]
    parallel: true
    required_outputs:
      - plans/production-analysis.md
      - plans/qc-analysis.md
      - plans/maintenance-plan.md
      - plans/safety-recommendations.md
      - plans/procurement-actions.md

  - id: discussion
    name: Alignment
    description: Plant Manager moderates alignment on report narrative and priorities
    participant_roles: [plant_manager, production_supervisor, qc_engineer, maintenance_engineer, safety_officer, supply_chain]
    leader_role: plant_manager
    parallel: false
    required_outputs: [report-outline.md]

  - id: report_generation
    name: Report Generation
    description: Agents produce final report sections and deliverables
    participant_roles: [plant_manager, process_engineer]
    parallel: false
    required_outputs:
      - output/daily-production-report.md
      - output/daily-production-report.pdf

  - id: qa
    name: QA Review
    description: Plant Manager and QC Engineer review report for accuracy
    participant_roles: [plant_manager, qc_engineer]
    parallel: false
    required_outputs: [qa-review.md]

checkpoints:
  - type: plan_approval
    after_phase: planning
    required_role: manager
    can_auto_delegate: true
    mandatory_human: false

  - type: final_acceptance
    after_phase: qa
    required_role: manager
    can_auto_delegate: false
    mandatory_human: true
```

---

## 9. API Design

All endpoints are prefixed with `/api/v2`. Authentication via Bearer token (JWT). Organization context inferred from the authenticated user's `org_id`.

### 9.1 Authentication

```
POST   /api/v2/auth/login
       Body: { email, password }
       Response: { token, user, expires_at }

POST   /api/v2/auth/logout
POST   /api/v2/auth/refresh
GET    /api/v2/auth/me
POST   /api/v2/auth/sso/saml/callback
POST   /api/v2/auth/sso/oidc/callback
```

### 9.2 Organizations

```
GET    /api/v2/org                        # Get current org details
PUT    /api/v2/org                        # Update org settings (Admin)
GET    /api/v2/org/usage                  # Usage metrics and limits
```

### 9.3 Users & RBAC

```
GET    /api/v2/users                      # List org users (Admin, Manager)
POST   /api/v2/users/invite               # Invite user (Admin)
       Body: { email, name, role }
GET    /api/v2/users/:id                  # Get user detail
PUT    /api/v2/users/:id                  # Update user (Admin)
       Body: { name, role, status }
DELETE /api/v2/users/:id                  # Deactivate user (Admin)
POST   /api/v2/users/:id/reset-password   # Trigger password reset
```

### 9.4 Projects

```
GET    /api/v2/projects                   # List projects (filtered by user access)
POST   /api/v2/projects                   # Create project (Manager, Admin)
       Body: { name, description, template_id, settings }
GET    /api/v2/projects/:id               # Get project detail
PUT    /api/v2/projects/:id               # Update project (Manager, Admin)
DELETE /api/v2/projects/:id               # Archive project (Admin)
GET    /api/v2/projects/:id/members       # List project members
POST   /api/v2/projects/:id/members       # Add member to project
DELETE /api/v2/projects/:id/members/:uid  # Remove member
GET    /api/v2/projects/:id/settings/workflow  # Get workflow settings
PUT    /api/v2/projects/:id/settings/workflow  # Update workflow settings
```

### 9.5 Tasks

```
GET    /api/v2/projects/:id/tasks          # List tasks (with filters: status, priority, date)
POST   /api/v2/projects/:id/tasks          # Submit new task (Manager, Operator)
       Body: { title, description, priority, metadata }
GET    /api/v2/projects/:id/tasks/:tid     # Get task detail + current phase
PUT    /api/v2/projects/:id/tasks/:tid     # Update task (cancel, reprioritize)
DELETE /api/v2/projects/:id/tasks/:tid     # Cancel task
GET    /api/v2/projects/:id/tasks/:tid/subtasks    # List subtasks
GET    /api/v2/projects/:id/tasks/:tid/artifacts   # List artifacts produced
GET    /api/v2/projects/:id/tasks/:tid/timeline    # Phase-by-phase execution log
```

### 9.6 Checkpoints

```
GET    /api/v2/projects/:id/checkpoints          # List all checkpoints (with filter: status)
GET    /api/v2/projects/:id/checkpoints/pending  # List pending checkpoints (alert count)
GET    /api/v2/projects/:id/checkpoints/:cid     # Get checkpoint detail with artifact preview
POST   /api/v2/projects/:id/checkpoints/:cid/decide
       # Requires Manager role (or Admin)
       # Requires mandatory_human=false for non-attorney approval
       Body: { action: "approved"|"rejected"|"overridden", feedback, override_data }
```

### 9.7 Agent Chat

```
GET    /api/v2/projects/:id/agents/:role/chat    # Get chat history
POST   /api/v2/projects/:id/agents/:role/chat    # Send message to agent
       Body: { message, context: "consultation"|"mid_workflow" }
```

### 9.8 Kanban

```
GET    /api/v2/projects/:id/kanban               # Get full kanban board state
PUT    /api/v2/projects/:id/kanban/tasks/:kid/move
       Body: { column: "backlog"|"in_progress"|"in_review"|"done" }
POST   /api/v2/projects/:id/kanban/tasks         # Create manual kanban task
PUT    /api/v2/projects/:id/kanban/tasks/:kid    # Update task details
DELETE /api/v2/projects/:id/kanban/tasks/:kid    # Delete manual task
```

### 9.9 Integrations

```
GET    /api/v2/integrations                      # List org integrations (Admin)
POST   /api/v2/integrations                      # Create integration (Admin)
       Body: { type, name, config, credentials }
       # credentials are written to vault, never returned
GET    /api/v2/integrations/:iid                 # Get integration (config only, no credentials)
PUT    /api/v2/integrations/:iid                 # Update integration (Admin)
DELETE /api/v2/integrations/:iid                 # Delete integration (Admin)
POST   /api/v2/integrations/:iid/test            # Test connection (Admin)
       Response: { ok: bool, latency_ms, error }
GET    /api/v2/integrations/:iid/audit           # Integration audit log
```

### 9.10 Team Templates

```
GET    /api/v2/templates                         # List available templates (system + org)
GET    /api/v2/templates/:tid                    # Get template detail
POST   /api/v2/templates                         # Create custom template (Admin)
       Body: { name, industry, description, config: WorkflowConfig }
PUT    /api/v2/templates/:tid                    # Update template (Admin)
DELETE /api/v2/templates/:tid                    # Delete custom template (Admin)
POST   /api/v2/templates/:tid/export             # Export as YAML
POST   /api/v2/templates/import                  # Import from YAML
```

### 9.11 Audit Log

```
GET    /api/v2/audit                             # Query audit log (Admin)
       Query: org_id, actor_type, actor_id, action, resource_type, resource_id,
              from, to, limit, offset
POST   /api/v2/audit/export                      # Export audit log (Admin)
       Body: { from, to, format: "json"|"csv"|"pdf", filters }
       Response: 202 Accepted + { export_id }
GET    /api/v2/audit/export/:eid                 # Poll export status + download URL
```

### 9.12 Analytics & Monitoring

```
GET    /api/v2/analytics/overview                # Org-level usage summary
GET    /api/v2/analytics/projects/:id            # Project-level metrics
       Response: { tasks_completed, avg_phase_duration, checkpoint_approval_rate,
                   agent_utilization, cost_estimate }
GET    /api/v2/analytics/agents                  # Agent performance metrics
GET    /api/v2/health                            # System health check (public)
GET    /api/v2/status                            # Detailed system status (Admin)
```

### 9.13 WebSocket Events

All real-time events are delivered via WebSocket at `/ws?project_id={id}&token={jwt}`.

```json
// Agent status update
{ "type": "agent_update", "role": "qc_engineer", "status": "working",
  "phase": "research", "message": "Analyzing defect data from Line 3..." }

// Phase transition
{ "type": "phase_change", "from": "research", "to": "planning",
  "task_id": "uuid", "timestamp": "..." }

// Checkpoint pending — pause for human
{ "type": "checkpoint_pending", "checkpoint_id": "uuid",
  "checkpoint_type": "plan_approval", "summary": "...",
  "artifact_path": "...", "required_role": "manager" }

// Checkpoint resolved
{ "type": "checkpoint_resolved", "checkpoint_id": "uuid",
  "action": "approved", "decided_by": "user:uuid" }

// Kanban update
{ "type": "kanban_update", "task_id": "uuid", "column": "in_progress",
  "agent": "senior_dev" }

// Chat response
{ "type": "chat_response", "role": "architect", "message": "...",
  "project_id": "uuid" }

// Task complete
{ "type": "task_complete", "task_id": "uuid",
  "artifacts": ["output/report.pdf", "output/data.json"],
  "duration_seconds": 2847 }

// Integration event
{ "type": "integration_event", "integration_id": "uuid",
  "action": "email_sent", "status": "success" }

// Error
{ "type": "error", "code": "AGENT_TIMEOUT", "message": "...",
  "task_id": "uuid", "phase": "development" }
```

---

## 10. Security & Compliance

### 10.1 Authentication & Authorization

- **Authentication:** JWT tokens (RS256), 8-hour expiry with refresh tokens (30-day, rotated on use)
- **SSO:** SAML 2.0 and OIDC support for enterprise (Okta, Azure AD, Google Workspace)
- **MFA:** TOTP (Google Authenticator compatible) for Admin and Manager roles — enforced by org policy
- **Session management:** Sessions invalidated on password change, role change, or user deactivation
- **RBAC enforcement:** All API endpoints check role on every request via middleware. Database RLS provides defense-in-depth.
- **Principle of least privilege:** Agent roles are scoped to minimum required tools. An agent with `db_read` permission cannot write. An agent without `email` permission cannot send email.

### 10.2 Data Security

- **Encryption at rest:** PostgreSQL transparent data encryption (TDE) on RDS. File system encryption (AES-256) for agent working directories.
- **Encryption in transit:** TLS 1.3 minimum for all connections (web, API, database, integration endpoints).
- **Credential vault:** Integration secrets (API keys, passwords, OAuth tokens) stored in a dedicated vault (HashiCorp Vault or AWS Secrets Manager). Never stored in PostgreSQL. Never logged. Never returned in API responses. Referenced by `vault_key` only.
- **Data isolation:** PostgreSQL Row-Level Security ensures org_id scoping on every table. Verified in security test suite.
- **Agent prompt security:** Credentials are never injected into agent prompts. Agents call a local integration proxy that handles authentication on their behalf.

### 10.3 Audit & Immutability

- The `audit_log` and `integration_audit_log` tables are protected by a PostgreSQL trigger that prevents `UPDATE` and `DELETE` operations.
- Audit log entries are written synchronously before the corresponding action is considered complete.
- For regulated industries (legal, manufacturing), audit logs can be exported to customer-controlled S3 for immutable archival.
- Data retention: minimum 90 days (default), configurable up to 7 years. Retention policy enforced by scheduled deletion job that respects org policy.

### 10.4 Agent Security Boundaries

- Each agent subprocess runs in a restricted environment: no network access except via the integration proxy.
- The integration proxy validates every request against the agent's role permissions before executing.
- File system access is sandboxed to the project working directory. No agent can read files outside its project scope.
- Agents cannot modify audit logs, system configuration, or user credentials.
- All agent prompts are logged. Any prompt containing patterns matching known prompt injection signatures (e.g., "ignore previous instructions") is flagged and reviewed.

### 10.5 Compliance Roadmap

| Standard | Target Date | Status |
|---|---|---|
| SOC 2 Type I | Q3 2026 | Planned |
| SOC 2 Type II | Q1 2027 | Planned |
| ISO 27001 | Q2 2027 | Planned |
| GDPR Article 28 (Data Processing Agreement) | Q2 2026 | In progress |
| HIPAA (healthcare vertical) | Q4 2027 | Future |

### 10.6 Responsible AI Guardrails

- **Mandatory human checkpoints** for high-stakes actions: sending external communications, regulatory submissions, financial approvals. These are hard constraints, not user-configurable.
- **Output review before delivery:** No client-facing or regulatory output is delivered without passing through the configured approval checkpoint.
- **Hallucination mitigation:** Agents are required to cite sources (from project files or integration data) for factual claims. Research phase agents must produce cited research, not unsourced assertions.
- **Prompt injection detection:** Integration layer validates that data read from external sources (email body, database fields, API responses) does not contain prompt injection payloads before injecting into agent context.

---

## 11. Success Metrics & KPIs

### 11.1 Product Metrics

| Metric | v2.0 Target (6 months post-launch) | Measurement |
|---|---|---|
| Monthly Active Organizations | 150 | Database |
| Average tasks run per org per month | 45 | Analytics |
| Task completion rate (tasks that reach Done) | > 85% | Analytics |
| Checkpoint approval rate (approved on first review) | > 75% | Analytics |
| Average task duration (research → completion) | < 90 minutes | Analytics |
| Integration setup success rate | > 90% | Analytics |
| Agent error rate (unhandled exceptions per 1000 tasks) | < 5 | Monitoring |
| P99 API response time | < 500ms | APM |

### 11.2 Business Metrics

| Metric | 6-Month Target | 12-Month Target |
|---|---|---|
| Monthly Recurring Revenue (MRR) | $180K | $620K |
| Annual Recurring Revenue (ARR) | $2.2M | $7.4M |
| Number of enterprise contracts (>$2K/mo) | 8 | 28 |
| Average Contract Value | $22K ARR | $26K ARR |
| Net Revenue Retention | > 110% | > 115% |
| Customer Acquisition Cost (CAC) | < $4,000 | < $3,200 |
| CAC Payback Period | < 9 months | < 7 months |
| Churn Rate (monthly) | < 2.5% | < 1.8% |

### 11.3 Engineering Quality Metrics

| Metric | Target |
|---|---|
| Test coverage (unit + integration) | > 80% |
| API uptime | 99.9% |
| Mean Time to Resolve critical incidents | < 4 hours |
| Deploy frequency | Weekly |
| Security vulnerability resolution (critical) | < 24 hours |

### 11.4 User Experience Metrics

| Metric | Target |
|---|---|
| Time to first task completion (new user) | < 15 minutes |
| Checkpoint decision time (user response to pending) | Median < 30 minutes |
| NPS Score | > 45 |
| Support ticket volume per active org | < 2/month |

---

## 12. Competitive Analysis

### 12.1 Direct Competitors

| Platform | Strengths | Weaknesses | Agent House Advantage |
|---|---|---|---|
| **CrewAI** | Python-first, flexible role definition, growing community | No UI, requires engineering to configure, no RBAC, no audit trail, developer-only | UI-first, enterprise RBAC, industry templates, no-code configuration |
| **AutoGen (Microsoft)** | Strong research backing, multi-agent conversation patterns, Azure integration | Research-grade, not production-ready, no industry templates, complex setup | Production-ready, enterprise SLA, industry-specific knowledge baked in |
| **LangGraph (LangChain)** | Graph-based workflow composition, strong ecosystem, good debugging tools | Developer-only, no industry templates, no RBAC, no UI, requires LangSmith for monitoring | Turn-key deployment, business-user accessible, built-in monitoring |
| **Relevance AI** | No-code agent builder, good integrations, reasonable UI | Limited team coordination (single-agent focus), no checkpoint system, weak audit | True multi-agent coordination, checkpoint system, full audit trail |
| **Zapier Central** | Huge integration library, non-technical accessible | Simple automation, not AI reasoning, limited agent intelligence | Deep reasoning, multi-phase projects, specialized domain knowledge |

### 12.2 Adjacent Competitors (Custom Solutions)

Large enterprises (Goldman Sachs, JP Morgan, KPMG) are building internal AI agent systems. These require $2–5M in engineering investment and 18+ months of build time. Agent House provides 80% of the capability at 5% of the cost, deployed in days rather than years.

### 12.3 Agent House Differentiation Matrix

| Capability | CrewAI | AutoGen | LangGraph | Relevance AI | **Agent House** |
|---|---|---|---|---|---|
| Industry templates | — | — | — | Partial | **5 industries** |
| Business UI | — | — | — | Basic | **Full Mission Control** |
| Human checkpoint system | — | Partial | Partial | — | **Full (configurable)** |
| RBAC | — | — | — | Basic | **4 roles + custom** |
| Audit trail | — | — | Partial | — | **Full (immutable)** |
| External integrations | Via LangChain | Limited | Via LangChain | Good | **20+ integrations** |
| No engineering to configure | — | — | — | Partial | **Yes** |
| Kanban task board | — | — | — | — | **Yes** |
| Enterprise SSO | — | — | — | — | **SAML + OIDC** |
| On-prem deployment | — | Partial | — | — | **Planned Q4** |

---

## 13. Go-to-Market Strategy

### 13.1 Phase 1 — Software & Marketing Verticals (Q2 2026)

**Target Buyers:** Engineering managers at SaaS companies (50–500 employees). Creative directors at digital marketing agencies (10–50 employees).

**Sales Motion:** Product-led growth (PLG). 14-day free trial, full functionality, 1 project limit. Conversion to paid via in-product upgrade prompt when project limit reached.

**Channels:**
- Product Hunt launch (software dev template)
- Hacker News Show HN
- Developer-focused content marketing (blog, YouTube)
- LinkedIn ads targeting engineering managers and creative directors
- Partnership with Claude / Anthropic AI showcase

**Key Metrics:** 500 signups in first 30 days, 8% trial-to-paid conversion, $50K MRR by end of Q2.

### 13.2 Phase 2 — Manufacturing & Legal Verticals (Q3 2026)

**Target Buyers:** Operations directors at mid-market manufacturers ($50M–$500M revenue). Managing partners at commercial law firms (15–80 attorneys).

**Sales Motion:** Sales-assisted (SDR outbound + AE close). Average deal $18K–$48K ARR. 90-day sales cycle.

**Channels:**
- Industry conference presence (Fabtech for manufacturing, ILTA for legal)
- Vertical-specific case studies (1 per vertical, developed with pilot customers)
- Partnership with industry consultants and system integrators
- Webinar series: "AI Workforce for [Industry]"

**Key Metrics:** 5 enterprise pilots in each vertical, 3 closed contracts per vertical by end of Q3.

### 13.3 Phase 3 — EPC & Enterprise Expansion (Q4 2026)

**Target Buyers:** Project directors at top-50 EPC firms. Enterprise IT leaders at Fortune 500 (multi-vertical deployment).

**Sales Motion:** Enterprise (named account, executive sponsor, 6-month sales cycle). Average deal $100K–$500K ARR for EPC. Multi-department enterprise deployment $200K–$1M ARR.

**Channels:**
- Engineering press (ENR, Construction Executive)
- Direct executive outreach via network
- Professional services partner channel (Big 4 consulting firms)
- Request for Proposal (RFP) responses (EPC firms frequently run formal procurement)

**Key Metrics:** 3 EPC enterprise contracts by Q4 ($300K+ ARR combined). 2 Fortune 500 enterprise deals.

### 13.4 Pricing Strategy

- Starter tier ($99/mo) anchors the market and captures SMB inbound
- Professional tier ($499/mo) is the primary revenue driver for mid-market
- Enterprise tier pricing is based on agent teams deployed, not seats — aligns cost with value delivered
- Industry-specific add-ons (e.g., "Legal Compliance Pack" with Westlaw integration) add $200–$500/mo per team

---

## 14. Roadmap

### Q1 2026 (Current State + Polish)

**Engineering Squads:** Platform (2), Frontend (1)

Milestone deliverables:
- v1.0 software development pipeline (complete)
- Human-in-the-loop checkpoints (complete)
- Kanban board view (complete)
- Subtask decomposition (complete)
- Direct agent chat (complete)
- PostgreSQL migration (in progress)
- Multi-organization data model (in progress)

**Exit Criteria:** Stable v1.5 release. 5 paying software dev customers. PostgreSQL migration complete.

---

### Q2 2026 — Multi-Industry Foundation

**Engineering Squads:** Platform (3), Frontend (2), Integrations (2)

**P0 Deliverables:**

| Feature | Squad | Timeline |
|---|---|---|
| RBAC system (4 roles, all endpoints) | Platform | Week 1–3 |
| Multi-organization API (users, projects, settings) | Platform | Week 2–5 |
| Team template engine (YAML-configurable) | Platform | Week 3–6 |
| Marketing Agency template | Platform | Week 4–6 |
| Legal Firm template | Platform | Week 4–6 |
| Manufacturing template | Platform | Week 5–7 |
| EPC template | Platform | Week 5–7 |
| Credential vault (AES-256 encrypted) | Platform | Week 1–3 |
| Email integration (SendGrid + SMTP) | Integrations | Week 2–5 |
| REST API integration (generic) | Integrations | Week 2–5 |
| S3 file storage integration | Integrations | Week 4–6 |
| Google Drive integration | Integrations | Week 5–7 |
| PDF report generation | Integrations | Week 4–6 |
| Organization settings UI | Frontend | Week 3–5 |
| RBAC management UI (user list, invite, roles) | Frontend | Week 3–5 |
| Integration management UI | Frontend | Week 5–7 |
| Audit log viewer UI | Frontend | Week 6–8 |

**P1 Deliverables:**

| Feature | Squad | Timeline |
|---|---|---|
| SSO SAML 2.0 | Platform | Week 6–8 |
| OIDC integration | Platform | Week 7–8 |
| PostgreSQL database integration | Integrations | Week 4–7 |
| SharePoint integration | Integrations | Week 6–8 |
| Slack notifications | Integrations | Week 7–8 |
| Multi-project dashboard | Frontend | Week 7–8 |

**Exit Criteria:** 5 industry templates live. 20+ paying customers across 3 verticals. $180K MRR. SOC 2 Type I audit initiated.

---

### Q3 2026 — Enterprise Hardening

**Engineering Squads:** Platform (3), Frontend (2), Integrations (2), Security (1)

**P0 Deliverables:**

| Feature | Squad | Timeline |
|---|---|---|
| Audit log export (JSON, CSV, PDF) | Platform | Week 1–2 |
| Data retention policy enforcement | Platform | Week 2–3 |
| Immutable audit log (PG triggers) | Platform | Week 1–2 |
| Advanced RBAC (custom roles, resource-level permissions) | Platform | Week 3–6 |
| Integration audit log | Platform | Week 2–4 |
| Jira integration | Integrations | Week 2–5 |
| GitHub integration | Integrations | Week 2–5 |
| MySQL integration | Integrations | Week 3–5 |
| Microsoft Teams integration | Integrations | Week 4–6 |
| HubSpot CRM integration | Integrations | Week 5–7 |
| Agent performance analytics | Platform | Week 5–7 |
| Configurable workflow templates (UI) | Frontend | Week 3–7 |
| Compliance report templates | Platform | Week 6–8 |
| Penetration test + remediation | Security | Week 4–8 |

**Exit Criteria:** SOC 2 Type I report. 8 enterprise contracts. $400K MRR. Zero critical security findings.

---

### Q4 2026 — Scale & Ecosystem

**Engineering Squads:** Platform (3), Frontend (2), Integrations (2), Infrastructure (1)

**P0 Deliverables:**

| Feature | Squad | Timeline |
|---|---|---|
| On-premise deployment option (Docker Compose + Helm chart) | Infrastructure | Week 2–8 |
| SAP integration (read-only) | Integrations | Week 2–6 |
| Oracle Primavera integration | Integrations | Week 2–6 |
| Westlaw API integration | Integrations | Week 3–6 |
| Template marketplace (community sharing) | Platform | Week 4–7 |
| Workflow A/B testing | Platform | Week 5–8 |
| Mobile-responsive dashboard | Frontend | Week 2–6 |
| Multi-model routing (Claude + GPT-4 + Gemini) | Platform | Week 3–7 |
| Webhook system (outbound events) | Platform | Week 1–3 |
| Public API v2 documentation + SDK | Platform | Week 2–5 |

**Exit Criteria:** On-prem deployment available for EPC enterprise clients. $620K MRR. 28 enterprise contracts. SOC 2 Type II audit in progress.

---

## 15. Risks & Mitigations

### 15.1 Technical Risks

| Risk | Probability | Impact | Mitigation |
|---|---|---|---|
| LLM API rate limits under concurrent load | Medium | High | Multi-provider routing (Anthropic + OpenAI + Gemini). Request queuing with backpressure. Per-org rate limit controls. |
| Agent hallucination in regulated industries | High | Critical | Mandatory human checkpoints before regulated outputs. Citation requirements in agent prompts. Hallucination detection layer. |
| Data leakage between organizations | Low | Critical | PostgreSQL RLS + application-layer org_id enforcement. Dual-layer validation in test suite. Quarterly security audit. |
| External integration failures breaking workflow | Medium | High | Integration circuit breakers. Graceful degradation (log failure, notify human, pause workflow). Retry with exponential backoff. |
| Agent subprocess hangs / zombie processes | Medium | Medium | Hard timeout per phase (configurable). Process supervisor with heartbeat. Goroutine leak detection in test suite. |
| Prompt injection via external data sources | Medium | High | Sanitization layer in integration proxy. Pattern detection for injection payloads. Isolated context injection (data never mixed with system prompt). |

### 15.2 Business Risks

| Risk | Probability | Impact | Mitigation |
|---|---|---|---|
| Enterprise sales cycle longer than projected | High | Medium | PLG motion generates SMB revenue while enterprise pipeline matures. Reduce enterprise friction with pilot program. |
| Competitor (Microsoft Copilot Studio, Salesforce Agentforce) enters multi-agent enterprise market | Medium | High | Industry-specific depth (templates, knowledge, artifacts) is 12–18 months ahead of generalist platforms. Speed to market. |
| LLM cost increase by model providers | Medium | Medium | Multi-provider routing to cheapest capable model per task type. Token budget controls per project. |
| Customer AI skepticism in regulated industries (legal, manufacturing) | High | Medium | Mandatory human checkpoints as a trust feature (not a workaround). Case studies with measurable ROI. Compliance-forward positioning. |
| Key talent loss (founding engineers) | Low | High | Competitive equity packages. Documentation culture. No single points of knowledge failure. |

### 15.3 Regulatory Risks

| Risk | Probability | Impact | Mitigation |
|---|---|---|---|
| EU AI Act compliance requirements | High | Medium | Audit trail already planned. Human oversight design already embedded. Engage compliance counsel by Q3 2026. |
| GDPR data subject access requests on agent logs | Medium | Medium | Implement DSAR workflow. PII tagging in audit log pipeline. |
| Legal industry state bar restrictions on AI-generated legal work | Medium | High | Mandatory attorney review checkpoints (hard constraint). Clear disclaimer in all legal outputs: "Drafted by AI, reviewed and approved by [attorney name]." |

---

## 16. Appendix: Industry Agent Configurations

### A. Software Development Team

**Workflow:** Triage → Template Selection → Research → Planning → Discussion → Development → QA → Dev Iteration (optional)

| Role | System Prompt Summary | Tools | Phases |
|---|---|---|---|
| **CEO** | Visionary tech CEO. Moderates discussion, sets quality bar, makes final architectural calls. Approves or rejects all phase outputs. | file_read, file_write, storage_read | triage, template_selection, discussion, qa |
| **PM** | Senior product manager. Writes user stories, acceptance criteria, competitive analysis. Ensures business requirements are met. | file_read, file_write, web_search | research, planning |
| **UX** | Senior UX designer. Documents user journeys, information architecture, interaction patterns, usability heuristics. Output: UX spec markdown. | file_read, file_write, web_search | research, planning |
| **UI** | Senior UI designer. Documents visual design system, component specs, color, typography. Output: UI spec markdown. | file_read, file_write, web_search | research, planning |
| **Security** | Senior security engineer. OWASP-aligned. Produces threat model, security requirements, security test cases. | file_read, file_write, web_search | research, planning, qa |
| **Architect** | Principal software architect. Proposes templates, designs system architecture, selects technology, writes ADRs. | file_read, file_write, web_search | template_selection, research, planning |
| **Senior Dev** | Expert full-stack developer (10+ years). Implements core business logic, APIs, data layer. Writes production-quality code with tests. | file_read, file_write, storage_read | development, dev_iteration |
| **Junior Dev** | Mid-level developer (3–5 years). Implements UI components, writes unit tests, handles boilerplate. Works under senior dev specs. | file_read, file_write | development, dev_iteration |

**Checkpoint Rules:**
- `template_approval` after template_selection: Manager required, auto-delegate after 30 min
- `plan_approval` after planning: Manager required, auto-delegate after 60 min
- `final_acceptance` after QA: Manager required, auto-delegate disabled (human required)

---

### B. Manufacturing Plant Team

**Workflow:** Triage → Data Gathering → Analysis → Alignment → Report Generation → QA Review

| Role | System Prompt Summary | Tools | Phases |
|---|---|---|---|
| **Plant Manager** | Oversees all plant operations. Makes operational decisions, escalates critical issues. Speaks in OEE, throughput, cost per unit. | file_read, file_write, api_call, email, storage_write, report_pdf | triage, discussion, qa, report_generation |
| **Production Supervisor** | Analyzes shift production data, tracks OEE per line, identifies bottlenecks. Root cause analysis required for deviations. | file_read, file_write, api_call, storage_read | research, planning |
| **QC Engineer** | Analyzes defect data, SPC analysis, NCR tracking, ISO 9001 quality reports. Uses 8D methodology. | file_read, file_write, api_call, storage_write | research, planning, qa |
| **Maintenance Engineer** | Reviews equipment health data, plans maintenance, analyzes MTBF/MTTR. Flags safety-critical equipment. | file_read, file_write, api_call, storage_read | research, planning |
| **Safety Officer (HSE)** | Reviews incident logs, assesses OSHA compliance, near-miss analysis, regulatory submissions. [IMMEDIATE ACTION REQUIRED] prefix for critical findings. | file_read, file_write, api_call, storage_write, email | research, planning, qa |
| **Supply Chain** | Monitors material inventory, procurement risk, supplier performance. Flags shortages with [PROCUREMENT ALERT]. | file_read, file_write, api_call, email | research, planning |
| **Process Engineer** | Analyzes process parameters, SPC, FMEA, process capability. Recommends process improvements. | file_read, file_write, api_call, storage_write, report_pdf | planning, report_generation |
| **Shift Lead** | Compiles shift handover documentation, tracks hourly targets, documents immediate issues. Concise bullet-point output. | file_read, file_write, storage_read | research |

**Industry-Specific Artifacts:** Daily Production Report (PDF), NCR Register (Excel), Maintenance Work Orders (JSON→PDF), Shift Handover Notes (markdown), Procurement Alerts (email).

**Checkpoint Rules:**
- `plan_approval` after analysis: Manager required, auto-delegate after 2 hours
- `final_acceptance` after QA: Manager required (mandatory human — no auto-delegate)
- `safety_escalation` (new type): Admin required, mandatory human, auto-delegate never

---

### C. EPC (Engineering, Procurement, Construction) Team

**Workflow:** Triage → Data Gathering → Analysis & Recommendations → Alignment → Report Generation → Director Review

| Role | System Prompt Summary | Tools | Phases |
|---|---|---|---|
| **Project Director** | Senior PM with 20+ years EPC experience. Speaks in SPI, CPI, earned value, critical path. Makes project-level decisions. | file_read, file_write, api_call, email, storage_write, report_pdf | triage, discussion, qa |
| **Design Engineer** | Manages engineering document register, IFC status, technical queries (TQ/RFI). Knows ISO drawing standards. | file_read, file_write, api_call, storage_read, storage_write | research, planning |
| **Procurement Manager** | Tracks purchase orders, long-lead items, vendor performance, expediting. Speaks in lead time, delivery schedule, spend vs. budget. | file_read, file_write, api_call, storage_read, email | research, planning |
| **Site Supervisor** | Compiles field progress data, manpower reports, weather delays, daily diary. Pragmatic, fact-based output. | file_read, file_write, api_call, storage_read | research |
| **HSE Manager** | Reviews site safety statistics (LTI, TRI, near-miss), STOP observations, toolbox talks compliance. ISO 45001 aligned. Mandatory escalation for LTI. | file_read, file_write, api_call, storage_write, email | research, planning, qa |
| **QA/QC Inspector** | Tracks NCRs, inspection records, punch lists, ITPs (Inspection & Test Plans). References project quality plan. | file_read, file_write, api_call, storage_read, storage_write | research, planning, qa |
| **Contracts Manager** | Reviews contract obligations, identifies variation triggers, drafts notices, tracks claims. Commercial focus. | file_read, file_write, storage_read | planning |
| **Planning Engineer** | Queries Primavera/P6 API, analyzes schedule float, identifies critical path delays, produces 3-week lookahead. | file_read, file_write, api_call, storage_read, storage_write, report_pdf | research, planning, report_generation |

**Industry-Specific Artifacts:** Weekly Progress Report (PDF, 30–50 pages), Variation Notice (DOCX), Schedule Update (P6 export), NCR Register (Excel), HSE Statistics Report (PDF), Procurement Status Report (Excel).

**Checkpoint Rules:**
- `plan_approval` after analysis: Manager required, auto-delegate after 4 hours (long cycle time expected)
- `final_acceptance` after director review: Director (Admin) required, mandatory human
- `regulatory_submission` (new type): Admin required, mandatory human, logged with submitter name

---

### D. Marketing Agency Team

**Workflow:** Triage → Strategy Research → Content Planning → Strategy Alignment → Content Execution → Quality Review

| Role | System Prompt Summary | Tools | Phases |
|---|---|---|---|
| **Creative Director** | Oversees creative quality and brand consistency. Final approver on all client-facing content. Knows brand strategy, storytelling, integrated marketing. | file_read, file_write, storage_read, storage_write | triage, discussion, qa |
| **Strategist** | Develops campaign strategy, positioning, messaging architecture, buyer journey mapping. Data-driven — references analytics in all recommendations. | file_read, file_write, api_call, web_search, storage_read | research, planning |
| **Copywriter** | Produces high-quality marketing copy (blog, email, social, web, ads). Adapts to brand voice. Understands copywriting principles (AIDA, PAS, StoryBrand). | file_read, file_write, storage_read | planning, execution |
| **Designer** | Produces visual briefs, asset specs, mood boards. Cannot produce images but delivers production-ready creative briefs. | file_read, file_write, storage_read | planning, execution |
| **SEO Specialist** | Keyword research, content optimization, on-page SEO recommendations, technical SEO audit. References search volume, difficulty, intent. | file_read, file_write, api_call, web_search, storage_read | research, planning, qa |
| **Social Media Manager** | Develops platform-specific content strategy, content calendar, community engagement scripts. Platform expertise: LinkedIn, Instagram, X, TikTok. | file_read, file_write, storage_read | planning, execution |
| **Analytics Lead** | Pulls and interprets performance data (GA4, HubSpot, LinkedIn Analytics). Builds attribution models, reports on ROI. | file_read, file_write, api_call, storage_write, report_pdf | research, qa |
| **Account Manager** | Client relationship context, brief interpretation, deliverable packaging, client communication. Brand compliance checker. | file_read, file_write, email, storage_write | triage, execution |

**Industry-Specific Artifacts:** Campaign Strategy Deck (PDF/markdown), Content Calendar (Excel), Blog Posts (markdown + HTML), Email Templates (HTML), Social Media Pack (per-platform text files), Performance Report (PDF).

**Checkpoint Rules:**
- `strategy_approval` after planning: Manager required, auto-delegate after 24 hours (client review cycle)
- `final_acceptance` after QA: Manager required, mandatory human (no AI output sent to client without human review)

---

### E. Legal Firm Team

**Workflow:** Triage → Legal Research → Issues Analysis → Attorney Review (mandatory checkpoint) → Drafting → Quality Review → Final Sign-Off (mandatory)

| Role | System Prompt Summary | Tools | Phases |
|---|---|---|---|
| **Managing Partner** | Senior attorney with 20+ years practice. Makes all client-facing decisions. Approves all deliverables. Ultimate responsibility. | file_read, file_write, storage_read | triage, discussion, qa |
| **Senior Associate** | Experienced attorney (7–12 years). Leads substantive legal analysis, reviews agent work product for accuracy, identifies legal risk. | file_read, file_write, api_call, storage_read | research, planning, qa |
| **Paralegal** | Experienced paralegal. Formats documents per firm standards, cite-checks, prepares filing packages, manages document logistics. | file_read, file_write, storage_read, storage_write | execution, qa |
| **Research Analyst** | Specialist legal researcher. Queries Westlaw/Lexis APIs, synthesizes case law, regulatory guidance, secondary sources. Cites all sources. | file_read, file_write, api_call, web_search, storage_read | research, planning |
| **Compliance Officer** | Regulatory compliance specialist. Reviews matters for compliance with relevant regulations (GDPR, SEC, OSHA, industry-specific). | file_read, file_write, api_call, storage_read | research, planning |
| **Contract Specialist** | Specialist in contract drafting and review. Identifies defined term inconsistencies, missing provisions, risky clauses. | file_read, file_write, storage_read | planning, execution |
| **Litigation Support** | Organizes case materials, document review, deposition preparation, case timeline. eDiscovery processes. | file_read, file_write, storage_read, storage_write | research, execution |
| **Client Relations** | Manages client communication context, drafts client-friendly summaries, manages expectation documentation. | file_read, file_write, email, storage_read | triage, execution |

**Industry-Specific Artifacts:** Legal Research Memorandum (PDF), Redlined Agreement (DOCX), Issues List (PDF), Due Diligence Report (PDF), Client Update Letter (PDF), Matter Audit Trail (PDF export).

**Checkpoint Rules (Legal firms have the most restrictive checkpoint rules):**
- `attorney_review` after issues_analysis: Manager (Attorney) required, mandatory human, NO auto-delegate
- `final_sign_off` after quality review: Admin (Managing Partner) required, mandatory human, NO auto-delegate
- `client_delivery` (new type): Manager required, mandatory human — any content going to client must be attorney-approved
- `regulatory_filing` (new type): Admin required, mandatory human — any regulatory submission requires named attorney authorization

**Special Legal Rules:**
- All matter files are privilege-protected — no cross-matter data access
- Client-facing output always includes attorney name and bar number in metadata
- Full Westlaw query log retained as part of research record
- Any AI-generated content delivered to client is labeled "Prepared by [Firm Name] using AI assistance, reviewed and approved by [Attorney Name], [Bar #]"
- Audit log for legal matters retained minimum 7 years (configurable to jurisdiction requirement)

---

*End of Document*

**Document Control:**
- Version 2.0 — Initial PRD for multi-industry expansion
- Author: Product Management
- Review: Engineering, Legal, Sales, Executive
- Next Review: 2026-06-18 (90 days)
- Approval: Required from CTO and CPO before engineering sprint planning
