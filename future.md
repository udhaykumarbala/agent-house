# Agent House — Future Architecture: Run Any Company

**Version:** 1.0
**Date:** 2026-03-18
**Status:** Strategic Plan
**Goal:** Transform Agent House from a software dev tool into a universal AI company runtime

---

## Table of Contents

1. [Vision](#1-vision)
2. [Architecture Overview](#2-architecture-overview)
3. [Phase 1: Tool Access Layer](#3-phase-1-tool-access-layer)
4. [Phase 2: Hooks & Trigger System](#4-phase-2-hooks--trigger-system)
5. [Phase 3: Routine Automations](#5-phase-3-routine-automations)
6. [Phase 4: Agent Marketplace & Composability](#6-phase-4-agent-marketplace--composability)
7. [Phase 5: Multi-Tenant Company Runtime](#7-phase-5-multi-tenant-company-runtime)
8. [Data Architecture](#8-data-architecture)
9. [Security Model](#9-security-model)
10. [Implementation Roadmap](#10-implementation-roadmap)

---

## 1. Vision

Today Agent House deploys 8 AI agents that build software through structured phases. Tomorrow it should deploy **any team of agents that runs any business function** — sales, marketing, legal, finance, HR, operations, customer support — with real tools, real data, and real-world side effects.

The missing pieces are:

| Gap | What It Means | Why It Matters |
|-----|---------------|----------------|
| **Tool Access** | Agents can only read/write local files | Can't query databases, send emails, call APIs, or interact with business systems |
| **Hooks & Triggers** | Work only starts when a human submits a task | Can't react to emails, webhooks, Slack messages, schedule events, or system alerts |
| **Routine Automations** | No recurring work | Can't run daily reports, weekly syncs, monitoring, or scheduled pipelines |

Solving these three gaps turns Agent House into a **company operating system** — not just a project executor.

---

## 2. Architecture Overview

```
                          ┌──────────────────────────────┐
                          │     TRIGGER LAYER            │
                          │                              │
                          │  Email  Webhook  Cron  Slack │
                          │  GitHub  FileWatch  DB-Watch │
                          │  Calendar  SMS  Custom       │
                          └──────────┬───────────────────┘
                                     │
                                     ▼
┌────────────────────────────────────────────────────────────────┐
│                    EVENT BUS (internal)                         │
│                                                                │
│  event.Emit("email.received", payload)                        │
│  event.Emit("webhook.github.pr_opened", payload)              │
│  event.Emit("cron.daily_report", payload)                     │
│  event.Emit("agent.task_complete", payload)                   │
│                                                                │
│  ┌─────────────┐  ┌─────────────┐  ┌──────────────────────┐  │
│  │  Router /   │  │  Filter /   │  │  Priority Queue /    │  │
│  │  Matcher    │  │  Transform  │  │  Rate Limiter        │  │
│  └─────────────┘  └─────────────┘  └──────────────────────┘  │
└────────────────────────────┬───────────────────────────────────┘
                             │
                             ▼
┌────────────────────────────────────────────────────────────────┐
│                  ORCHESTRATOR (existing, enhanced)              │
│                                                                │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐     │
│  │  Triage  │  │ Research │  │ Planning │  │ Execute  │     │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘     │
│                                                                │
│  ┌──────────────────────────────────────────────────────┐     │
│  │            AGENT RUNTIME (per-agent sandbox)          │     │
│  │                                                       │     │
│  │   Agent ──→ Tool Router ──→ ┌── FileSystem (existing)│     │
│  │                             ├── Database              │     │
│  │                             ├── Email                 │     │
│  │                             ├── HTTP Client           │     │
│  │                             ├── Slack/Teams           │     │
│  │                             ├── CRM (Salesforce, HB)  │     │
│  │                             ├── Storage (S3, GCS)     │     │
│  │                             ├── Browser               │     │
│  │                             ├── Document Gen          │     │
│  │                             └── Custom MCP Tools      │     │
│  └──────────────────────────────────────────────────────┘     │
└────────────────────────────────────────────────────────────────┘
                             │
                             ▼
┌────────────────────────────────────────────────────────────────┐
│                  AUTOMATION ENGINE                              │
│                                                                │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────┐    │
│  │ Cron Scheduler│  │ Workflow     │  │ Monitoring &     │    │
│  │ (robfig/cron)│  │ Templates    │  │ Alerting         │    │
│  └──────────────┘  └──────────────┘  └──────────────────┘    │
│                                                                │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────┐    │
│  │ Pipeline     │  │ Approval     │  │ Audit Log        │    │
│  │ Engine       │  │ Chains       │  │                  │    │
│  └──────────────┘  └──────────────┘  └──────────────────┘    │
└────────────────────────────────────────────────────────────────┘
```

---

## 3. Phase 1: Tool Access Layer

### 3.1 Design Philosophy

Tools are **capabilities agents can invoke** during task execution. Today agents only have `Read`, `Write`, `Edit`, `Bash`, `Grep`. We need to expand this to every system a company uses — but with **permission boundaries**, **audit trails**, and **sandboxing**.

### 3.2 Tool Registry Architecture

```
internal/
  tools/
    registry.go          # Central tool registry, discovery, validation
    tool.go              # Tool interface & base types
    permission.go        # Per-role, per-tool permission matrix
    audit.go             # Every tool invocation logged
    sandbox.go           # Rate limits, quotas, dry-run mode

    database/
      postgres.go        # PostgreSQL: query, insert, update, migrate
      mysql.go           # MySQL adapter
      mongodb.go         # MongoDB adapter
      redis.go           # Redis: get/set/pub-sub
      sqlite.go          # SQLite for lightweight local DBs

    email/
      smtp.go            # Send emails (SMTP)
      imap.go            # Read/watch inbox (IMAP)
      template.go        # Email template engine
      parser.go          # Parse incoming email (headers, body, attachments)

    communication/
      slack.go           # Slack: send/read messages, channels, threads
      discord.go         # Discord bot integration
      teams.go           # Microsoft Teams
      twilio.go          # SMS/voice via Twilio

    http/
      client.go          # Generic HTTP client (GET/POST/PUT/DELETE)
      graphql.go         # GraphQL client
      oauth.go           # OAuth2 flow management
      webhook_sender.go  # Outbound webhook delivery

    storage/
      s3.go              # AWS S3: upload/download/list
      gcs.go             # Google Cloud Storage
      local.go           # Enhanced local filesystem
      gdrive.go          # Google Drive
      dropbox.go         # Dropbox

    crm/
      salesforce.go      # Salesforce: leads, contacts, opportunities
      hubspot.go         # HubSpot: deals, contacts, tickets

    erp/
      sap.go             # SAP integration
      netsuite.go        # NetSuite

    documents/
      pdf.go             # PDF generation (go-wkhtmltopdf or similar)
      spreadsheet.go     # Excel/CSV generation (excelize)
      presentation.go    # Slide deck generation
      markdown.go        # Markdown → HTML/PDF rendering

    browser/
      scraper.go         # Web scraping (chromedp/rod)
      form.go            # Form filling automation
      screenshot.go      # Page screenshots

    payments/
      stripe.go          # Stripe: invoices, subscriptions, charges

    calendar/
      google.go          # Google Calendar: create/read events
      outlook.go         # Outlook Calendar

    analytics/
      ga4.go             # Google Analytics 4
      mixpanel.go        # Mixpanel events/funnels
      posthog.go         # PostHog analytics
```

### 3.3 Tool Interface

```go
// internal/tools/tool.go

package tools

type Tool interface {
    // Identity
    Name() string                    // "database.postgres.query"
    Category() string                // "database"
    Description() string             // Human-readable description for agent prompts

    // Schema — what the agent sees
    InputSchema() JSONSchema         // Parameters the agent must provide
    OutputSchema() JSONSchema        // What the agent gets back

    // Execution
    Execute(ctx context.Context, input ToolInput) (ToolOutput, error)

    // Validation
    ValidateInput(input ToolInput) error

    // Permissions
    RequiredPermissions() []Permission
    DangerLevel() DangerLevel        // safe, moderate, dangerous, destructive
}

type DangerLevel int
const (
    DangerSafe        DangerLevel = iota  // Read-only, no side effects
    DangerModerate                         // Writes to internal systems
    DangerDangerous                        // Writes to external systems (emails, Slack)
    DangerDestructive                      // Irreversible (delete, payment)
)

type ToolInput struct {
    Parameters map[string]interface{}
    AgentRole  string
    ProjectID  string
    TaskID     string
}

type ToolOutput struct {
    Data      interface{}
    Metadata  map[string]string
    Duration  time.Duration
    BytesUsed int64
}
```

### 3.4 Permission Matrix

Permissions are defined per-role and per-tool, with escalation for dangerous operations:

```go
// internal/tools/permission.go

type PermissionRule struct {
    Role        string       // "ceo", "pm", "senior_dev", "*"
    ToolPattern string       // "database.*", "email.send", "*"
    Access      AccessLevel  // deny, read, write, admin
    Conditions  []Condition  // rate limits, time windows, approval required
}

type AccessLevel int
const (
    AccessDeny  AccessLevel = iota
    AccessRead                      // Can query/fetch but not modify
    AccessWrite                     // Can create/update
    AccessAdmin                     // Can delete, configure, grant
)
```

**Default Permission Matrix:**

| Tool Category | CEO | PM | UX/UI | Security | Architect | Sr Dev | Jr Dev |
|--------------|-----|-----|-------|----------|-----------|--------|--------|
| Database (read) | write | read | deny | read | read | write | read |
| Database (write) | write | deny | deny | deny | write | write | read |
| Email (send) | write | write | deny | deny | deny | deny | deny |
| Email (read) | read | read | deny | read | deny | deny | deny |
| Slack | write | write | deny | deny | deny | read | deny |
| HTTP Client | write | write | read | write | write | write | read |
| Storage | write | write | write | read | write | write | write |
| CRM | write | write | deny | deny | deny | deny | deny |
| Payments | admin | deny | deny | deny | deny | deny | deny |
| Browser | read | read | read | write | read | write | read |
| Documents | write | write | write | deny | write | write | write |
| Calendar | write | write | deny | deny | deny | deny | deny |

### 3.5 Tool Configuration (per-project)

```json
// project/.agent-house/tools.json
{
  "tools": {
    "database.postgres": {
      "enabled": true,
      "config": {
        "connection_string": "${POSTGRES_URL}",
        "max_connections": 5,
        "read_only_roles": ["pm", "security"],
        "allowed_tables": ["*"],
        "blocked_tables": ["users_credentials", "api_keys"],
        "max_rows_per_query": 10000
      }
    },
    "email.smtp": {
      "enabled": true,
      "config": {
        "host": "smtp.company.com",
        "port": 587,
        "from": "agents@company.com",
        "require_approval_for": ["external_domains"],
        "daily_send_limit": 100
      }
    },
    "slack": {
      "enabled": true,
      "config": {
        "bot_token": "${SLACK_BOT_TOKEN}",
        "allowed_channels": ["#agent-updates", "#project-feed"],
        "blocked_channels": ["#exec-private"],
        "mention_humans_require_approval": true
      }
    },
    "http": {
      "enabled": true,
      "config": {
        "allowed_domains": ["api.company.com", "api.github.com"],
        "blocked_domains": ["*.internal.company.com"],
        "timeout_seconds": 30,
        "max_response_size_mb": 10
      }
    }
  }
}
```

### 3.6 How Agents Use Tools

Tools are injected into agent prompts based on role permissions and project config:

```
You are the Product Manager agent. In addition to file operations, you have access to:

## Available Tools

### database.query
Query the project database (read-only for your role).
Parameters: { "sql": "SELECT ...", "params": [...] }

### email.send
Send an email on behalf of the team.
Parameters: { "to": "...", "subject": "...", "body": "...", "cc": [...] }

### slack.post
Post a message to an allowed Slack channel.
Parameters: { "channel": "#agent-updates", "text": "...", "thread_ts": "..." }

### http.get / http.post
Make HTTP requests to allowed domains.
Parameters: { "url": "...", "headers": {...}, "body": {...} }

To use a tool, output: TOOL_CALL: tool_name
{json parameters}
END_TOOL_CALL
```

The agent runtime parses `TOOL_CALL:` blocks (similar to how it parses `DELEGATE:` and `FILE:` today), routes them through the tool registry, and injects the result back into the agent's context.

### 3.7 Implementation: Tool Router in Agent Runtime

```go
// Enhanced internal/agent/agent.go

type AgentRuntime struct {
    agent       *Agent
    toolRouter  *tools.Router
    auditLog    *tools.AuditLog
    sandbox     *tools.Sandbox
}

func (r *AgentRuntime) ProcessResponse(response string) (*Response, error) {
    // Existing: parse DELEGATE:, REVIEW:, FILE:, CREATE_FILE:
    resp := ParseResponse(response)

    // NEW: parse TOOL_CALL: blocks
    toolCalls := ParseToolCalls(response)

    for _, call := range toolCalls {
        // Check permissions
        if err := r.toolRouter.CheckPermission(r.agent.Role, call.Tool); err != nil {
            resp.ToolResults = append(resp.ToolResults, ToolResult{
                Tool:  call.Tool,
                Error: fmt.Sprintf("Permission denied: %v", err),
            })
            continue
        }

        // Check sandbox limits (rate, quota)
        if err := r.sandbox.Check(call); err != nil {
            // Dangerous tools → request human approval via checkpoint
            if call.DangerLevel >= DangerDangerous {
                r.requestApproval(call)
                continue
            }
            resp.ToolResults = append(resp.ToolResults, ToolResult{
                Tool:  call.Tool,
                Error: fmt.Sprintf("Rate limited: %v", err),
            })
            continue
        }

        // Execute
        result, err := r.toolRouter.Execute(ctx, call)

        // Audit log
        r.auditLog.Record(AuditEntry{
            Agent:     r.agent.Role,
            Tool:      call.Tool,
            Input:     call.Input,
            Output:    result,
            Error:     err,
            Timestamp: time.Now(),
        })

        resp.ToolResults = append(resp.ToolResults, ToolResult{
            Tool:   call.Tool,
            Data:   result.Data,
            Error:  errStr(err),
        })
    }

    return resp, nil
}
```

### 3.8 MCP (Model Context Protocol) Integration

Instead of building every tool from scratch, leverage **MCP servers** as a universal tool adapter:

```go
// internal/tools/mcp/connector.go

type MCPConnector struct {
    servers map[string]*MCPServer  // "postgres", "slack", "email"
}

// Agent House can connect to any MCP server, instantly gaining
// access to hundreds of tools without custom Go code.
//
// Example: connect to the Postgres MCP server
// → agents get: query, insert, update, delete, schema_info tools
//
// Example: connect to the Slack MCP server
// → agents get: send_message, read_channel, create_channel tools
```

**Priority MCP servers to integrate:**

| MCP Server | Tools Gained | Priority |
|------------|-------------|----------|
| `@modelcontextprotocol/server-postgres` | SQL query, schema inspect | P0 |
| `@modelcontextprotocol/server-slack` | Send/read messages, channels | P0 |
| `@anthropic/mcp-email` | SMTP send, IMAP read | P0 |
| `@modelcontextprotocol/server-github` | Issues, PRs, repos | P0 |
| `@modelcontextprotocol/server-filesystem` | Enhanced file ops | P1 |
| `@modelcontextprotocol/server-google-drive` | Drive read/write | P1 |
| `@modelcontextprotocol/server-brave-search` | Web search | P1 |
| `@stripe/mcp-server` | Payments, invoices | P2 |
| Custom MCP servers | Any business-specific tool | P2 |

**MCP Config:**

```json
// project/.agent-house/mcp.json
{
  "mcpServers": {
    "postgres": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-postgres", "${POSTGRES_URL}"],
      "roles": ["ceo", "pm", "architect", "senior_dev"]
    },
    "slack": {
      "command": "npx",
      "args": ["-y", "@anthropic/mcp-slack"],
      "env": { "SLACK_BOT_TOKEN": "${SLACK_BOT_TOKEN}" },
      "roles": ["ceo", "pm"]
    },
    "github": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-github"],
      "env": { "GITHUB_TOKEN": "${GITHUB_TOKEN}" },
      "roles": ["*"]
    }
  }
}
```

---

## 4. Phase 2: Hooks & Trigger System

### 4.1 The Problem

Today, work only starts when a human types a task into the dashboard. A real company reacts to **events**: emails arrive, customers file tickets, PRs are opened, deadlines approach, metrics spike. Agents need to be **event-driven**.

### 4.2 Event Bus Architecture

```go
// internal/events/bus.go

package events

type EventBus struct {
    subscribers map[string][]Subscriber
    queue       chan Event
    mu          sync.RWMutex
}

type Event struct {
    ID        string                 `json:"id"`
    Type      string                 `json:"type"`      // "email.received", "webhook.github", "cron.fire"
    Source    string                 `json:"source"`    // "imap-watcher", "webhook-server", "cron-scheduler"
    Payload   map[string]interface{} `json:"payload"`
    Timestamp time.Time              `json:"timestamp"`
    ProjectID string                 `json:"project_id"`
    Priority  int                    `json:"priority"`  // 0=low, 10=critical
}

type Subscriber struct {
    Pattern string                   // "email.*", "webhook.github.pr_opened", "*"
    Handler func(Event) error
    Filter  func(Event) bool         // Additional filtering
}

func (bus *EventBus) Emit(event Event) {
    bus.queue <- event
}

func (bus *EventBus) Subscribe(pattern string, handler func(Event) error) {
    // Glob-style pattern matching on event types
}

func (bus *EventBus) processLoop() {
    for event := range bus.queue {
        // Match subscribers, fan-out, handle errors
        // Dead letter queue for failed deliveries
    }
}
```

### 4.3 Trigger Types

#### 4.3.1 Email Trigger

Watches an inbox (IMAP IDLE) and converts emails into agent tasks:

```go
// internal/triggers/email.go

type EmailTrigger struct {
    imapClient *imap.Client
    bus        *events.EventBus
    rules      []EmailRule
}

type EmailRule struct {
    Name       string   `json:"name"`
    Match      EmailMatch `json:"match"`
    Action     TriggerAction `json:"action"`
}

type EmailMatch struct {
    From      string `json:"from"`       // regex: ".*@client.com"
    Subject   string `json:"subject"`    // regex: ".*urgent.*"
    To        string `json:"to"`         // "support@company.com"
    HasAttachment bool `json:"has_attachment"`
    Labels    []string `json:"labels"`
}

type TriggerAction struct {
    Type      string `json:"type"`       // "create_task", "notify_agent", "run_workflow"
    ProjectID string `json:"project_id"`
    AgentRole string `json:"agent_role"` // Which agent handles it
    Template  string `json:"template"`   // Task template with {{.Subject}}, {{.Body}}, etc.
    Priority  int    `json:"priority"`
}
```

**Example: Customer Support Email → Agent Task**

```json
{
  "name": "support_email_to_task",
  "match": {
    "to": "support@company.com",
    "subject": ".*"
  },
  "action": {
    "type": "create_task",
    "project_id": "customer-support",
    "agent_role": "pm",
    "template": "New support ticket from {{.From}}:\n\nSubject: {{.Subject}}\n\n{{.Body}}\n\nPlease triage, draft a response, and escalate if needed.",
    "priority": 5
  }
}
```

#### 4.3.2 Webhook Trigger

Receives HTTP webhooks from external systems and routes them to agents:

```go
// internal/triggers/webhook.go

type WebhookServer struct {
    router *http.ServeMux
    bus    *events.EventBus
    hooks  map[string]WebhookConfig
}

// Registers: POST /hooks/{hook_id}
// Each webhook has a secret for HMAC verification

type WebhookConfig struct {
    ID          string `json:"id"`
    Secret      string `json:"secret"`       // HMAC-SHA256 verification
    Source      string `json:"source"`       // "github", "stripe", "custom"
    EventMap    map[string]TriggerAction `json:"event_map"`
}
```

**Example: GitHub PR Opened → Code Review Agent**

```json
{
  "id": "github-pr-review",
  "source": "github",
  "secret": "${GITHUB_WEBHOOK_SECRET}",
  "event_map": {
    "pull_request.opened": {
      "type": "create_task",
      "project_id": "{{.repository.name}}",
      "agent_role": "senior_dev",
      "template": "Review PR #{{.pull_request.number}}: {{.pull_request.title}}\n\nDescription: {{.pull_request.body}}\nBranch: {{.pull_request.head.ref}}\nFiles changed: {{.pull_request.changed_files}}\n\nPerform a thorough code review.",
      "priority": 7
    },
    "issues.opened": {
      "type": "create_task",
      "project_id": "{{.repository.name}}",
      "agent_role": "pm",
      "template": "New issue #{{.issue.number}}: {{.issue.title}}\n\n{{.issue.body}}\n\nTriage this issue.",
      "priority": 5
    }
  }
}
```

#### 4.3.3 Cron/Schedule Trigger

Time-based triggers for recurring work:

```go
// internal/triggers/cron.go

type CronScheduler struct {
    cron *cron.Cron  // robfig/cron
    bus  *events.EventBus
    jobs map[string]CronJob
}

type CronJob struct {
    ID        string        `json:"id"`
    Name      string        `json:"name"`
    Schedule  string        `json:"schedule"`   // "0 9 * * 1" (Mon 9am)
    Timezone  string        `json:"timezone"`   // "America/New_York"
    Action    TriggerAction `json:"action"`
    Enabled   bool          `json:"enabled"`
    LastRun   time.Time     `json:"last_run"`
    NextRun   time.Time     `json:"next_run"`
}
```

**Example Cron Jobs:**

```json
[
  {
    "id": "daily-standup-summary",
    "name": "Daily Standup Summary",
    "schedule": "0 9 * * 1-5",
    "timezone": "America/New_York",
    "action": {
      "type": "create_task",
      "project_id": "operations",
      "agent_role": "pm",
      "template": "Generate the daily standup summary:\n1. Review all kanban board changes in the last 24h\n2. Identify blockers\n3. Draft a summary and post to #standup on Slack"
    }
  },
  {
    "id": "weekly-security-scan",
    "name": "Weekly Security Scan",
    "schedule": "0 2 * * 0",
    "action": {
      "type": "create_task",
      "project_id": "security-ops",
      "agent_role": "security",
      "template": "Run the weekly security audit:\n1. Check all project dependencies for CVEs\n2. Review access logs for anomalies\n3. Generate compliance report\n4. Email report to security@company.com"
    }
  },
  {
    "id": "monthly-financial-report",
    "name": "Monthly Financial Report",
    "schedule": "0 8 1 * *",
    "action": {
      "type": "create_task",
      "project_id": "finance",
      "agent_role": "ceo",
      "template": "Generate monthly financial report:\n1. Query revenue from Stripe\n2. Pull expense data from ERP\n3. Calculate margins and KPIs\n4. Generate PDF report\n5. Email to exec team"
    }
  }
]
```

#### 4.3.4 File Watch Trigger

React when files appear in watched directories or cloud storage:

```go
// internal/triggers/filewatch.go

type FileWatcher struct {
    watchers map[string]*fsnotify.Watcher
    s3Polls  map[string]*S3Poller
    bus      *events.EventBus
}

type FileWatchRule struct {
    ID     string `json:"id"`
    Type   string `json:"type"`     // "local", "s3", "gcs"
    Path   string `json:"path"`     // "/data/incoming/" or "s3://bucket/prefix/"
    Pattern string `json:"pattern"` // "*.csv", "*.pdf"
    Action TriggerAction `json:"action"`
}
```

**Example: Invoice PDF uploaded → Finance Agent processes it**

```json
{
  "id": "invoice-processor",
  "type": "s3",
  "path": "s3://company-docs/invoices/incoming/",
  "pattern": "*.pdf",
  "action": {
    "type": "create_task",
    "project_id": "finance",
    "agent_role": "pm",
    "template": "New invoice received: {{.FileName}}\n\n1. Extract invoice details (vendor, amount, due date)\n2. Match against PO in database\n3. If match: approve and update ERP\n4. If no match: flag for human review\n5. Move to processed/ folder"
  }
}
```

#### 4.3.5 Database Watch Trigger

React to database changes (new rows, updated statuses):

```go
// internal/triggers/dbwatch.go

type DBWatcher struct {
    db   *sql.DB
    bus  *events.EventBus
    rules []DBWatchRule
}

type DBWatchRule struct {
    ID       string `json:"id"`
    Table    string `json:"table"`
    Event    string `json:"event"`     // "INSERT", "UPDATE", "DELETE"
    Condition string `json:"condition"` // "status = 'escalated'"
    PollInterval time.Duration `json:"poll_interval"`
    Action   TriggerAction `json:"action"`
}
```

**Example: Escalated support ticket → CEO agent reviews**

```json
{
  "id": "escalated-ticket-review",
  "table": "support_tickets",
  "event": "UPDATE",
  "condition": "status = 'escalated' AND assigned_to IS NULL",
  "poll_interval": "30s",
  "action": {
    "type": "create_task",
    "project_id": "customer-support",
    "agent_role": "ceo",
    "template": "Escalated ticket #{{.id}}: {{.subject}}\n\nCustomer: {{.customer_name}}\nPriority: {{.priority}}\nHistory: {{.conversation_summary}}\n\nReview and decide on resolution."
  }
}
```

#### 4.3.6 Slack/Chat Trigger

Agents respond to messages in communication channels:

```go
// internal/triggers/slack.go

type SlackTrigger struct {
    client *slack.Client
    bus    *events.EventBus
    rules  []SlackRule
}

type SlackRule struct {
    ID       string `json:"id"`
    Channel  string `json:"channel"`   // "#ask-agents"
    Mention  bool   `json:"mention"`   // Only when @agent-house mentioned
    Keywords []string `json:"keywords"` // Trigger on specific words
    Action   TriggerAction `json:"action"`
}
```

#### 4.3.7 Signal System (Agent-to-Agent Internal Triggers)

Agents can emit signals that trigger other agents — enabling reactive workflows without human intervention:

```go
// internal/triggers/signal.go

type Signal struct {
    Name    string                 `json:"name"`
    From    string                 `json:"from"`     // agent role that emitted
    Payload map[string]interface{} `json:"payload"`
}

// Examples:
// PM completes research → signal triggers Architect to start planning
// Security finds vulnerability → signal triggers Senior Dev for hotfix
// CEO approves budget → signal triggers PM to start procurement
```

**Signal Configuration:**

```json
{
  "signals": [
    {
      "name": "vulnerability_found",
      "from": "security",
      "triggers": {
        "type": "create_task",
        "agent_role": "senior_dev",
        "template": "URGENT: Security vulnerability found.\n\nSeverity: {{.severity}}\nComponent: {{.component}}\nDetails: {{.description}}\n\nPatch immediately.",
        "priority": 9
      }
    },
    {
      "name": "customer_churn_risk",
      "from": "pm",
      "triggers": {
        "type": "create_task",
        "agent_role": "ceo",
        "template": "Customer {{.customer_name}} showing churn signals:\n\n{{.signals}}\n\nReview account and decide on retention strategy.",
        "priority": 7
      }
    }
  ]
}
```

### 4.4 Trigger Configuration UI

Dashboard additions for managing triggers:

```
┌─────────────────────────────────────────────────────────────┐
│  TRIGGERS                                              [+New]│
├──────┬───────────────────────┬────────┬──────────┬──────────┤
│ Type │ Name                  │ Status │ Last Fire│ Actions  │
├──────┼───────────────────────┼────────┼──────────┼──────────┤
│ 📧   │ Support Email → PM    │ Active │ 2m ago   │ ⚙️ 🗑️   │
│ 🔗   │ GitHub PR → Review    │ Active │ 1h ago   │ ⚙️ 🗑️   │
│ ⏰   │ Daily Standup         │ Active │ 9:00 AM  │ ⚙️ 🗑️   │
│ ⏰   │ Weekly Security Scan  │ Active │ Sun 2AM  │ ⚙️ 🗑️   │
│ 📁   │ Invoice Processor     │ Paused │ 3d ago   │ ⚙️ 🗑️   │
│ 🔔   │ Escalated Tickets     │ Active │ 45m ago  │ ⚙️ 🗑️   │
│ 💬   │ Slack #ask-agents     │ Active │ 15m ago  │ ⚙️ 🗑️   │
└──────┴───────────────────────┴────────┴──────────┴──────────┘
│ Event Log                                                    │
│ 10:32 📧 support@company.com → Task #847 (PM)              │
│ 10:15 🔗 PR #234 opened → Task #846 (Senior Dev)           │
│ 09:00 ⏰ Daily Standup fired → Task #845 (PM)              │
│ 08:45 🔔 Ticket #1092 escalated → Task #844 (CEO)          │
└──────────────────────────────────────────────────────────────┘
```

### 4.5 Trigger Pipeline (Event → Task)

```
Event Source
    │
    ▼
┌──────────┐     ┌───────────┐     ┌────────────┐     ┌──────────┐
│ Receiver │ ──→ │ Validator │ ──→ │ Transformer│ ──→ │  Router  │
│ (ingest) │     │ (auth,    │     │ (template  │     │ (match   │
│          │     │  schema)  │     │  render)   │     │  rules)  │
└──────────┘     └───────────┘     └────────────┘     └──────────┘
                                                            │
                              ┌──────────────────────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │ Action Executor  │
                    │                  │
                    │ • create_task    │
                    │ • notify_agent   │
                    │ • run_workflow   │
                    │ • send_webhook   │
                    │ • update_kanban  │
                    │ • emit_signal    │
                    └──────────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │ Audit + Metrics  │
                    │ (every trigger   │
                    │  fire logged)    │
                    └──────────────────┘
```

---

## 5. Phase 3: Routine Automations

### 5.1 Automation Engine

The automation engine combines triggers + tools + workflows into repeatable, schedulable, composable automation pipelines.

```go
// internal/automation/engine.go

type AutomationEngine struct {
    scheduler *CronScheduler
    workflows map[string]*Workflow
    bus       *events.EventBus
    executor  *WorkflowExecutor
}

type Workflow struct {
    ID          string         `json:"id"`
    Name        string         `json:"name"`
    Description string         `json:"description"`
    Trigger     TriggerConfig  `json:"trigger"`      // What starts it
    Steps       []WorkflowStep `json:"steps"`         // What happens
    ErrorPolicy ErrorPolicy    `json:"error_policy"`  // On failure
    Timeout     time.Duration  `json:"timeout"`
    Enabled     bool           `json:"enabled"`
}

type WorkflowStep struct {
    ID        string            `json:"id"`
    Name      string            `json:"name"`
    Type      StepType          `json:"type"`
    Config    map[string]interface{} `json:"config"`
    DependsOn []string          `json:"depends_on"`  // DAG dependencies
    Condition string            `json:"condition"`    // Go template condition
    Timeout   time.Duration     `json:"timeout"`
    Retries   int               `json:"retries"`
}

type StepType string
const (
    StepAgent     StepType = "agent"       // Run an agent with a task
    StepTool      StepType = "tool"        // Execute a tool directly
    StepCondition StepType = "condition"   // If/else branching
    StepParallel  StepType = "parallel"    // Run steps concurrently
    StepApproval  StepType = "approval"    // Human approval gate
    StepDelay     StepType = "delay"       // Wait N minutes/hours
    StepSubflow   StepType = "subflow"     // Invoke another workflow
)
```

### 5.2 Workflow Templates

Pre-built automation templates for common company operations:

#### 5.2.1 Customer Onboarding Pipeline

```yaml
id: customer-onboarding
name: "New Customer Onboarding"
trigger:
  type: webhook
  source: stripe
  event: customer.subscription.created

steps:
  - id: extract-info
    name: "Extract Customer Info"
    type: tool
    config:
      tool: http.get
      url: "https://api.stripe.com/v1/customers/{{.customer_id}}"

  - id: create-workspace
    name: "Set Up Customer Workspace"
    type: agent
    config:
      role: pm
      task: |
        New customer signed up: {{.extract-info.name}} ({{.extract-info.email}})
        Plan: {{.plan_name}}

        1. Create project folder structure
        2. Set up initial kanban board with onboarding tasks
        3. Draft welcome email
    depends_on: [extract-info]

  - id: send-welcome
    name: "Send Welcome Email"
    type: tool
    config:
      tool: email.send
      to: "{{.extract-info.email}}"
      subject: "Welcome to {{.company_name}}!"
      body: "{{.create-workspace.welcome_email_draft}}"
    depends_on: [create-workspace]

  - id: notify-team
    name: "Notify Sales Team"
    type: tool
    config:
      tool: slack.post
      channel: "#new-customers"
      text: "New customer: {{.extract-info.name}} on {{.plan_name}} plan"
    depends_on: [extract-info]

  - id: schedule-kickoff
    name: "Schedule Kickoff Call"
    type: tool
    config:
      tool: calendar.create
      title: "Kickoff: {{.extract-info.company}}"
      duration: 30m
      attendees: ["{{.extract-info.email}}", "csm@company.com"]
    depends_on: [send-welcome]
```

#### 5.2.2 Daily Operations Pipeline

```yaml
id: daily-ops
name: "Daily Operations Cycle"
trigger:
  type: cron
  schedule: "0 8 * * 1-5"
  timezone: "America/New_York"

steps:
  - id: gather-metrics
    name: "Gather Business Metrics"
    type: parallel
    config:
      steps:
        - tool: database.query
          sql: "SELECT COUNT(*), SUM(revenue) FROM orders WHERE date = CURRENT_DATE - 1"
        - tool: http.get
          url: "https://api.stripe.com/v1/balance"
        - tool: database.query
          sql: "SELECT COUNT(*) as open, COUNT(*) FILTER (WHERE priority='high') as urgent FROM support_tickets WHERE status='open'"

  - id: analyze
    name: "CEO Analyzes Daily State"
    type: agent
    config:
      role: ceo
      task: |
        Daily operations review for {{.date}}:

        Orders yesterday: {{.gather-metrics.0.count}}, Revenue: ${{.gather-metrics.0.sum}}
        Current balance: ${{.gather-metrics.1.available}}
        Open tickets: {{.gather-metrics.2.open}} ({{.gather-metrics.2.urgent}} urgent)

        1. Identify any anomalies or concerns
        2. Decide priorities for today
        3. Draft the daily briefing
    depends_on: [gather-metrics]

  - id: post-briefing
    name: "Post Daily Briefing"
    type: tool
    config:
      tool: slack.post
      channel: "#daily-briefing"
      text: "{{.analyze.briefing}}"
    depends_on: [analyze]

  - id: check-urgent
    name: "Handle Urgents"
    type: condition
    config:
      condition: "{{.gather-metrics.2.urgent}} > 5"
      if_true:
        type: agent
        config:
          role: pm
          task: "URGENT: {{.gather-metrics.2.urgent}} high-priority tickets. Triage and assign."
      if_false:
        type: tool
        config:
          tool: slack.post
          channel: "#ops"
          text: "All clear — {{.gather-metrics.2.urgent}} urgent tickets (within threshold)"
    depends_on: [gather-metrics]
```

#### 5.2.3 Content Pipeline

```yaml
id: content-pipeline
name: "Weekly Content Production"
trigger:
  type: cron
  schedule: "0 10 * * 1"  # Monday 10am

steps:
  - id: research-trends
    name: "Research Industry Trends"
    type: agent
    config:
      role: pm
      task: |
        Research this week's industry trends:
        1. Search top 10 industry blogs for trending topics
        2. Check competitor social media for themes
        3. Review our analytics for top-performing content
        4. Propose 5 content ideas with audience + channel

  - id: approve-topics
    name: "Approve Content Topics"
    type: approval
    config:
      approver: "human"
      message: "Review proposed content topics for this week"
      timeout: 4h
      auto_approve_after: 4h  # If no human responds, auto-approve
    depends_on: [research-trends]

  - id: create-content
    name: "Create Content"
    type: parallel
    config:
      for_each: "{{.approve-topics.approved_topics}}"
      step:
        type: agent
        config:
          role: ux  # Content writer role
          task: "Write content piece: {{.item.title}}\n\nBrief: {{.item.brief}}\nChannel: {{.item.channel}}\nAudience: {{.item.audience}}"
    depends_on: [approve-topics]

  - id: review-content
    name: "Editorial Review"
    type: agent
    config:
      role: ceo
      task: "Review all content pieces for brand voice, accuracy, and quality. Approve or request revisions."
    depends_on: [create-content]
```

### 5.3 Monitoring & Alerting Automations

```go
// internal/automation/monitor.go

type Monitor struct {
    ID        string        `json:"id"`
    Name      string        `json:"name"`
    Check     MonitorCheck  `json:"check"`
    Interval  time.Duration `json:"interval"`
    Condition string        `json:"condition"`   // "value > threshold"
    Alert     AlertConfig   `json:"alert"`
}

type MonitorCheck struct {
    Type   string `json:"type"`    // "http", "database", "metric"
    Config map[string]interface{} `json:"config"`
}

type AlertConfig struct {
    Channels []string `json:"channels"`  // ["slack:#alerts", "email:ops@co.com"]
    Severity string   `json:"severity"`  // "info", "warning", "critical"
    AgentTask *TriggerAction `json:"agent_task"` // Optionally trigger agent
}
```

**Example Monitors:**

```json
[
  {
    "id": "api-health",
    "name": "API Health Check",
    "check": { "type": "http", "config": { "url": "https://api.company.com/health" } },
    "interval": "60s",
    "condition": "status_code != 200 || response_time_ms > 2000",
    "alert": {
      "channels": ["slack:#alerts"],
      "severity": "critical",
      "agent_task": {
        "type": "create_task",
        "agent_role": "architect",
        "template": "API health check failing. Status: {{.status_code}}, Response time: {{.response_time_ms}}ms. Investigate immediately."
      }
    }
  },
  {
    "id": "revenue-anomaly",
    "name": "Revenue Anomaly Detection",
    "check": {
      "type": "database",
      "config": { "sql": "SELECT SUM(amount) FROM orders WHERE date = CURRENT_DATE" }
    },
    "interval": "1h",
    "condition": "value < (7_day_avg * 0.7)",
    "alert": {
      "channels": ["slack:#revenue", "email:cfo@company.com"],
      "severity": "warning",
      "agent_task": {
        "type": "create_task",
        "agent_role": "ceo",
        "template": "Revenue anomaly detected. Today's revenue (${{.value}}) is {{.pct_below}}% below 7-day average. Investigate causes."
      }
    }
  }
]
```

### 5.4 Automation Dashboard

```
┌─────────────────────────────────────────────────────────────────┐
│  AUTOMATIONS                                     [+ New Workflow]│
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─── Active Workflows ─────────────────────────────────────┐   │
│  │                                                           │   │
│  │  ✅ Customer Onboarding    │ Webhook │ 3 runs today      │   │
│  │  ✅ Daily Operations       │ Cron    │ Next: Tomorrow 8am │   │
│  │  ✅ Content Pipeline       │ Cron    │ Next: Monday 10am  │   │
│  │  ⏸️ Invoice Processing     │ S3 Watch│ Paused             │   │
│  │  ✅ PR Review Auto-assign  │ Webhook │ 12 runs today      │   │
│  │                                                           │   │
│  └───────────────────────────────────────────────────────────┘   │
│                                                                  │
│  ┌─── Monitors ────────────────────────────────────────────┐    │
│  │                                                          │    │
│  │  🟢 API Health        │ 60s  │ OK (142ms)               │    │
│  │  🟢 Revenue Tracking  │ 1h   │ OK ($12,400 today)       │    │
│  │  🟡 Support Queue     │ 5m   │ WARN: 8 urgent tickets   │    │
│  │  🟢 DB Connections    │ 30s  │ OK (23/100 used)         │    │
│  │                                                          │    │
│  └──────────────────────────────────────────────────────────┘    │
│                                                                  │
│  ┌─── Recent Runs ─────────────────────────────────────────┐    │
│  │                                                          │    │
│  │  10:32  Customer Onboarding  ✅ Success  (2m 14s)       │    │
│  │  10:15  PR Review            ✅ Success  (45s)          │    │
│  │  09:00  Daily Operations     ✅ Success  (4m 30s)       │    │
│  │  08:45  Support Escalation   ⚠️ Partial  (1m 02s)      │    │
│  │  08:00  DB Backup            ✅ Success  (12m 05s)      │    │
│  │                                                          │    │
│  └──────────────────────────────────────────────────────────┘    │
└──────────────────────────────────────────────────────────────────┘
```

---

## 6. Phase 4: Agent Marketplace & Composability

### 6.1 Custom Agent Roles

Move beyond the fixed 8-agent team. Let users **define custom agents** with custom roles, tools, and prompts:

```json
// .agent-house/agents/customer-success-manager.json
{
  "role": "csm",
  "name": "Customer Success Manager",
  "description": "Manages customer relationships, monitors health scores, prevents churn",
  "system_prompt": "You are the Customer Success Manager...",
  "tools": ["database.query", "email.send", "slack.post", "crm.hubspot.*", "calendar.*"],
  "permissions": {
    "database": "read",
    "email": "write",
    "crm": "write"
  },
  "triggers": {
    "churn_risk": {
      "type": "database_watch",
      "condition": "health_score < 50",
      "task_template": "Customer {{.name}} health score dropped to {{.health_score}}. Investigate and intervene."
    }
  }
}
```

### 6.2 Agent Templates per Industry

```
internal/
  templates/
    industries/
      software-dev/          # Current 8 agents (CEO, PM, UX, UI, Security, Arch, SrDev, JrDev)
      marketing-agency/      # Creative Director, Strategist, Copywriter, Designer, Media Buyer, Analyst
      legal-firm/            # Managing Partner, Associate, Paralegal, Compliance, Researcher, Clerk
      ecommerce/             # CEO, Product Manager, Inventory, Marketing, Support, Analyst, Fulfillment
      real-estate/           # Broker, Agent, Analyst, Marketing, Transaction Coordinator, Inspector
      saas-company/          # CEO, CTO, PM, Engineering Lead, DevOps, Support, Sales, Marketing
      consulting-firm/       # Partner, Engagement Manager, Consultant, Analyst, Knowledge Manager
      healthcare/            # Director, Clinician, Scheduler, Billing, Compliance, Records
      construction/          # Project Manager, Estimator, Safety, QA, Scheduler, Procurement
      finance/               # CFO, Controller, Analyst, Compliance, Auditor, Tax
```

### 6.3 Agent Skill Marketplace

Community-contributed agent skills and workflow templates:

```
Agent Marketplace (future SaaS feature)
│
├── Agent Roles
│   ├── Sales Development Rep (SDR)
│   ├── Data Analyst
│   ├── Social Media Manager
│   ├── QA Tester
│   └── DevOps Engineer
│
├── Workflow Templates
│   ├── Lead Qualification Pipeline
│   ├── Content Calendar Automation
│   ├── Incident Response Runbook
│   ├── Employee Onboarding
│   └── Invoice Processing
│
├── Tool Integrations
│   ├── Salesforce Connector
│   ├── Shopify Connector
│   ├── QuickBooks Connector
│   ├── Zendesk Connector
│   └── Jira Connector
│
└── Industry Packs
    ├── E-commerce Operations Pack
    ├── SaaS Growth Pack
    ├── Legal Practice Pack
    └── Construction Project Pack
```

---

## 7. Phase 5: Multi-Tenant Company Runtime

### 7.1 Organization Model

```go
type Organization struct {
    ID          string
    Name        string
    Plan        string           // "starter", "pro", "enterprise"
    Teams       []Team
    Integrations []Integration
    Settings    OrgSettings
}

type Team struct {
    ID       string
    Name     string              // "Engineering", "Marketing", "Finance"
    Agents   []AgentConfig       // Custom agent roster
    Projects []Project
    Tools    []ToolConfig        // Team-specific tool access
    Triggers []TriggerConfig
    Automations []WorkflowConfig
}
```

### 7.2 Cross-Team Collaboration

Agents from different teams can collaborate:

```
Engineering Team                    Marketing Team
┌──────────────┐                   ┌──────────────┐
│ Senior Dev   │ ←── signal ─────→ │ Content Lead │
│ "Feature X   │   "feature_shipped" │ "Write blog  │
│  is shipped" │                   │  about X"    │
└──────────────┘                   └──────────────┘
       │                                  │
       │                                  │
       ▼                                  ▼
   Updates Jira                    Posts to social
   Closes PR                       Schedules campaign
   Notifies QA                     Updates website
```

### 7.3 Resource Quotas & Billing

```go
type ResourceQuota struct {
    AgentInvocationsPerMonth int      // Claude API calls
    ToolCallsPerMonth        int      // External tool executions
    StorageGB                float64  // File storage
    AutomationRuns           int      // Workflow executions
    ConcurrentAgents         int      // Max parallel agents
}

// Pricing tiers
var Plans = map[string]ResourceQuota{
    "starter":    {1000, 5000, 5, 100, 2},
    "pro":        {10000, 50000, 50, 1000, 8},
    "enterprise": {100000, 500000, 500, 10000, 32},
}
```

---

## 8. Data Architecture

### 8.1 Migration: File System → PostgreSQL

Current state: JSON files on disk. Target state: PostgreSQL with the file system as a cache layer.

```sql
-- Core tables
CREATE TABLE organizations (...);
CREATE TABLE teams (...);
CREATE TABLE projects (...);
CREATE TABLE agents (...);
CREATE TABLE tasks (...);

-- Tool & Trigger tables
CREATE TABLE tool_configs (
    id UUID PRIMARY KEY,
    team_id UUID REFERENCES teams(id),
    tool_name VARCHAR(100),        -- "database.postgres"
    config JSONB,                   -- Connection details (encrypted)
    permissions JSONB,              -- Role-based access
    enabled BOOLEAN DEFAULT true
);

CREATE TABLE triggers (
    id UUID PRIMARY KEY,
    team_id UUID REFERENCES teams(id),
    name VARCHAR(200),
    type VARCHAR(50),               -- "email", "webhook", "cron", "filewatch", "dbwatch", "slack"
    config JSONB,                   -- Type-specific configuration
    action JSONB,                   -- What to do when triggered
    enabled BOOLEAN DEFAULT true,
    last_fired_at TIMESTAMPTZ,
    fire_count BIGINT DEFAULT 0
);

CREATE TABLE automations (
    id UUID PRIMARY KEY,
    team_id UUID REFERENCES teams(id),
    name VARCHAR(200),
    workflow JSONB,                  -- Full workflow definition
    trigger_id UUID REFERENCES triggers(id),
    enabled BOOLEAN DEFAULT true,
    last_run_at TIMESTAMPTZ,
    run_count BIGINT DEFAULT 0
);

CREATE TABLE automation_runs (
    id UUID PRIMARY KEY,
    automation_id UUID REFERENCES automations(id),
    status VARCHAR(20),             -- "running", "success", "failed", "partial"
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    steps JSONB,                    -- Step-by-step execution log
    error TEXT
);

-- Audit trail
CREATE TABLE audit_log (
    id BIGSERIAL PRIMARY KEY,
    timestamp TIMESTAMPTZ DEFAULT NOW(),
    org_id UUID,
    team_id UUID,
    agent_role VARCHAR(50),
    tool_name VARCHAR(100),
    action VARCHAR(50),             -- "execute", "read", "write", "delete"
    input JSONB,
    output JSONB,
    error TEXT,
    duration_ms INTEGER,
    approved_by VARCHAR(100)        -- Human who approved (if applicable)
);

-- Secret management
CREATE TABLE secrets (
    id UUID PRIMARY KEY,
    team_id UUID REFERENCES teams(id),
    name VARCHAR(100),              -- "POSTGRES_URL", "SLACK_BOT_TOKEN"
    encrypted_value BYTEA,          -- AES-256-GCM encrypted
    created_at TIMESTAMPTZ,
    rotated_at TIMESTAMPTZ
);
```

### 8.2 Secret Management

Secrets (API keys, DB passwords, tokens) are never stored in plaintext:

```go
// internal/secrets/vault.go

type Vault struct {
    masterKey []byte              // Derived from env var or HSM
    store     SecretStore         // PostgreSQL or file-based
}

func (v *Vault) Set(name string, value string) error {
    encrypted := aesGCMEncrypt(v.masterKey, []byte(value))
    return v.store.Save(name, encrypted)
}

func (v *Vault) Get(name string) (string, error) {
    encrypted, err := v.store.Load(name)
    decrypted := aesGCMDecrypt(v.masterKey, encrypted)
    return string(decrypted), err
}

// Secrets are injected into tool configs at runtime:
// "${POSTGRES_URL}" → vault.Get("POSTGRES_URL") → actual connection string
```

---

## 9. Security Model

### 9.1 Defense in Depth

```
Layer 1: Authentication & Authorization
  ├── RBAC: Org Admin > Team Lead > Member > Viewer
  ├── API keys with scoped permissions
  └── SSO/SAML for enterprise

Layer 2: Agent Sandboxing
  ├── Each agent runs in isolated context
  ├── Tool access gated by role + project + team permissions
  ├── No agent can access another team's data
  └── Resource quotas enforced per-agent

Layer 3: Tool-Level Security
  ├── Allowlists: specific tables, channels, domains, recipients
  ├── Blocklists: sensitive tables, private channels, internal URLs
  ├── Rate limits per tool, per agent, per hour
  └── Dangerous operations require human approval (checkpoint)

Layer 4: Audit & Compliance
  ├── Every tool invocation logged with full input/output
  ├── Every trigger fire logged
  ├── Every agent decision logged
  ├── Immutable audit log (append-only)
  └── Export to SIEM (Splunk, Datadog, etc.)

Layer 5: Data Protection
  ├── Secrets encrypted at rest (AES-256-GCM)
  ├── TLS for all external connections
  ├── PII detection and masking in agent outputs
  └── Data retention policies per org
```

### 9.2 Human Approval Gates for Dangerous Operations

```go
// Tools with DangerLevel >= Dangerous automatically trigger checkpoints

// Example: Agent wants to send an external email
// 1. Agent outputs: TOOL_CALL: email.send { "to": "client@external.com", ... }
// 2. Tool router detects: email to external domain → DangerDangerous
// 3. Checkpoint created: "Agent PM wants to send email to client@external.com"
// 4. Human approves/rejects in dashboard
// 5. If approved: email sent, logged
// 6. If rejected: agent informed, adapts approach

// Auto-approval rules (configurable):
// - Internal emails: auto-approve
// - Slack to #agent-updates: auto-approve
// - Database reads: auto-approve
// - Database writes to non-sensitive tables: auto-approve
// - External emails: require approval
// - Payments: always require approval
// - Deletes: always require approval
```

---

## 10. Implementation Roadmap

### Phase 1: Tool Access Layer (4-6 weeks)

```
Week 1-2: Foundation
  ├── [ ] Tool interface and registry
  ├── [ ] Permission system (role × tool matrix)
  ├── [ ] Audit logging for all tool calls
  ├── [ ] TOOL_CALL: parsing in agent response handler
  ├── [ ] Tool result injection back into agent context
  └── [ ] Secret vault (env-var based, later upgrade to encrypted store)

Week 3-4: Core Tools
  ├── [ ] database/postgres.go — Query, insert, update (parameterized only)
  ├── [ ] email/smtp.go — Send emails with templates
  ├── [ ] email/imap.go — Read inbox, search, parse
  ├── [ ] http/client.go — Generic REST client with allowlists
  ├── [ ] communication/slack.go — Post messages, read channels
  └── [ ] storage/s3.go — Upload, download, list

Week 5-6: MCP + Dashboard
  ├── [ ] MCP connector — run MCP servers as tool providers
  ├── [ ] Tool configuration UI in dashboard
  ├── [ ] Tool test/dry-run mode
  ├── [ ] Per-project tool config (.agent-house/tools.json)
  └── [ ] Sandbox: rate limits, quotas, dry-run
```

### Phase 2: Hooks & Trigger System (4-6 weeks)

```
Week 7-8: Event Bus + Core Triggers
  ├── [ ] Event bus (in-process, channel-based)
  ├── [ ] Email trigger (IMAP IDLE watcher)
  ├── [ ] Webhook receiver (POST /hooks/{id} with HMAC auth)
  ├── [ ] Cron scheduler (robfig/cron)
  └── [ ] Event → Task pipeline (template rendering, routing)

Week 9-10: Advanced Triggers
  ├── [ ] File watch trigger (fsnotify + S3 polling)
  ├── [ ] Database watch trigger (polling-based)
  ├── [ ] Slack trigger (bot events API)
  ├── [ ] Signal system (agent-to-agent events)
  └── [ ] GitHub webhook handler (PR, issues, releases)

Week 11-12: Trigger Dashboard + Config
  ├── [ ] Trigger management UI
  ├── [ ] Trigger event log viewer
  ├── [ ] Trigger testing (simulate event)
  ├── [ ] Trigger enable/disable/pause
  └── [ ] Trigger metrics (fire count, success rate, latency)
```

### Phase 3: Routine Automations (4-6 weeks)

```
Week 13-14: Workflow Engine
  ├── [ ] Workflow definition schema (YAML/JSON)
  ├── [ ] Step types: agent, tool, condition, parallel, approval, delay
  ├── [ ] DAG executor (topological sort, parallel branches)
  ├── [ ] Step result passing (output of step A → input of step B)
  └── [ ] Error handling (retry, skip, abort, notify)

Week 15-16: Monitoring + Templates
  ├── [ ] Monitor system (HTTP, DB, metric checks)
  ├── [ ] Alert routing (Slack, email, agent task)
  ├── [ ] Pre-built workflow templates (5-10 common patterns)
  ├── [ ] Workflow builder UI (visual DAG editor)
  └── [ ] Automation run history and replay

Week 17-18: Polish + Integration
  ├── [ ] Automation dashboard (active, paused, history, monitors)
  ├── [ ] Workflow import/export
  ├── [ ] Workflow versioning
  ├── [ ] Cross-project workflows
  └── [ ] Documentation and examples
```

### Phase 4: Agent Marketplace (6-8 weeks)

```
Week 19-22: Custom Agents + Industry Templates
  ├── [ ] Custom agent role definition (JSON config)
  ├── [ ] Agent prompt builder UI
  ├── [ ] 5 industry agent packs (marketing, legal, ecommerce, consulting, finance)
  ├── [ ] Agent skill composition (mix and match capabilities)
  └── [ ] Agent testing sandbox

Week 23-26: Marketplace
  ├── [ ] Agent/workflow sharing format
  ├── [ ] Community marketplace UI
  ├── [ ] Rating and review system
  ├── [ ] Version management
  └── [ ] One-click install for agents, workflows, tool configs
```

### Phase 5: Multi-Tenant Company Runtime (8-12 weeks)

```
Week 27-30: PostgreSQL Migration
  ├── [ ] Schema design and migration scripts
  ├── [ ] Dual-write (file + DB) transition period
  ├── [ ] Read path migration (DB primary, file fallback)
  ├── [ ] File system deprecation
  └── [ ] Performance benchmarks

Week 31-34: Multi-Tenancy
  ├── [ ] Organization/team/member model
  ├── [ ] RBAC system
  ├── [ ] Team isolation (agents, tools, data)
  ├── [ ] Cross-team collaboration (signals, shared projects)
  └── [ ] SSO/SAML integration

Week 35-38: Billing + Scale
  ├── [ ] Resource quota enforcement
  ├── [ ] Usage metering
  ├── [ ] Stripe billing integration
  ├── [ ] Horizontal scaling (multiple orchestrator instances)
  └── [ ] Observability (metrics, tracing, structured logging)
```

---

## Summary: What Each Phase Unlocks

| Phase | Duration | Unlocks |
|-------|----------|---------|
| **1. Tool Access** | 6 weeks | Agents can query DBs, send emails, call APIs, post to Slack, manage files in cloud storage. They become **useful beyond code**. |
| **2. Hooks & Triggers** | 6 weeks | Work starts automatically from emails, webhooks, schedules, file uploads, DB changes, chat messages. The system becomes **reactive**. |
| **3. Automations** | 6 weeks | Multi-step workflows run on schedule or on trigger with branching, approvals, and monitoring. The system becomes **autonomous**. |
| **4. Marketplace** | 8 weeks | Custom agents for any role, industry templates, community sharing. The system becomes **universal**. |
| **5. Multi-Tenant** | 12 weeks | Organizations, teams, billing, scale. The system becomes a **product**. |

**Total: ~38 weeks (9 months) from dev tool to company OS.**

After all 5 phases, Agent House can:
- Run a **marketing agency**: Content agents react to briefs, create campaigns, post to social, track analytics
- Run a **law firm**: Legal agents monitor filings, research cases, draft documents, manage deadlines
- Run an **e-commerce operation**: Agents process orders, manage inventory, handle support, run promotions
- Run a **SaaS company**: Agents triage bugs, review PRs, onboard customers, generate reports, handle billing
- Run **any company**: Custom agents + custom tools + custom triggers + custom workflows = infinite flexibility
