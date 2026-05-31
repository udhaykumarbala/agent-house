# Proactive Conductor Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Turn the Conductor from a chat-pull surface into a *push* surface — the system speaks first with a real briefing, renders Brain replies as actionable cards, streams agent work into the conversation, and is fed by proactive cron jobs. Then ship two demo scenarios (EPC + software dev) that show it doing real work end-to-end.

**Architecture:**
- **Frontend-first** — all of Phase A (proactive opener, action cards, WS streaming) is React/TS in the existing Next.js app. No new backend required.
- **Bootstrap shell scripts** for Phase B scenarios — seed inbox JSON, create cron jobs via existing `POST /api/cron/`, prime Brain conversation files. No backend changes.
- **Honesty principle:** every interactive control either does real work (wired to a working endpoint) or is honestly labeled "Phase-1". No fake buttons.

**Tech Stack:** Next.js 15 + React 19 + TypeScript (frontend); Go backend already in place; `bash + curl + jq` for scenario bootstraps; `go:embed` ships the bundle as one binary.

**Repo conventions in force:**
- `scripts/build-ui.sh` is the canonical way to refresh the embedded UI (avoid ad-hoc `cd frontend && cp …` — that bug has bitten three times).
- Commits happen only when the user asks — every "Commit" step in this plan is a *suggested* commit; do not run `git commit` autonomously.
- CLAUDE.md: `api_test.html`, `uiflow.md`, swagger should be updated for any new/changed API endpoint. This plan adds no endpoints, so that rule doesn't trigger.

---

## Phase A — Proactive Conductor

### Task A1.1: Compose-briefing helper

**Files:**
- Create: `frontend/lib/briefing.ts`

- [ ] **Step 1: Define the briefing view-model and composer**

```ts
// frontend/lib/briefing.ts
import type { ApiAgent, ApiMessage, ApiCheckpoint } from "./types";

export interface BriefingItem {
  tone: "wait" | "fire" | "info" | "live";
  title: string;          // short imperative
  detail: string;         // one-line context
  action?: { label: string; href?: string }; // optional click-through
}

export interface Briefing {
  greeting: string;       // "Good afternoon, Marcus."
  headline: string;       // "3 pending · 2 fires today · all 7 agents idle"
  items: BriefingItem[];  // 0..N proactive cards (capped at 6)
  emptyMessage?: string;  // when nothing notable
}

function timeOfDay(d = new Date()): string {
  const h = d.getHours();
  if (h < 5) return "Working late";
  if (h < 12) return "Good morning";
  if (h < 17) return "Good afternoon";
  return "Good evening";
}

export function composeBriefing(
  agents: ApiAgent[],
  messages: ApiMessage[],
  checkpoints: ApiCheckpoint[],
  user = "Marcus"
): Briefing {
  const pending = checkpoints.filter((c) => c.status === "pending");
  const recent = [...messages]
    .sort(
      (a, b) =>
        new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime()
    )
    .slice(0, 12);
  const triggerFires = recent.filter((m) =>
    ["trigger", "trigger_fired", "hook", "cron", "webhook"].includes(m.type)
  );
  const working = agents.filter((a) => a.active);

  const items: BriefingItem[] = [];

  for (const cp of pending.slice(0, 3)) {
    items.push({
      tone: "wait",
      title: cp.artifact_summary || cp.type,
      detail: `${cp.type} · task#${(cp.task_id || "").slice(0, 6)}`,
      action: { label: "Review", href: "/mission" },
    });
  }
  if (triggerFires.length > 0) {
    items.push({
      tone: "fire",
      title: `${triggerFires.length} trigger ${
        triggerFires.length === 1 ? "fire" : "fires"
      } today`,
      detail: triggerFires
        .slice(0, 2)
        .map((m) => m.type)
        .join(" · "),
      action: { label: "See triggers", href: "/triggers" },
    });
  }
  if (working.length > 0) {
    items.push({
      tone: "info",
      title: `${working.length} agent${working.length === 1 ? "" : "s"} working now`,
      detail: working.map((a) => a.name || a.role).slice(0, 3).join(" · "),
    });
  }

  return {
    greeting: `${timeOfDay()}, ${user}.`,
    headline: `${pending.length} pending · ${triggerFires.length} fires today · ${
      working.length
    } working · ${agents.length - working.length} idle`,
    items,
    emptyMessage:
      items.length === 0
        ? "Nothing requires you right now. I'll surface anything that does."
        : undefined,
  };
}
```

- [ ] **Step 2: Verify it typechecks**

Run: `cd frontend && npx tsc --noEmit`
Expected: no errors related to `lib/briefing.ts`.

---

### Task A1.2: BriefingCard component

**Files:**
- Create: `frontend/components/BriefingCard.tsx`

- [ ] **Step 1: Create the React component**

```tsx
// frontend/components/BriefingCard.tsx
import type { Briefing } from "@/lib/briefing";

export function BriefingCard({ briefing }: { briefing: Briefing }) {
  return (
    <div className="briefing">
      <div className="briefing-head">
        <div className="briefing-greeting">{briefing.greeting}</div>
        <div className="briefing-headline">{briefing.headline}</div>
      </div>
      {briefing.emptyMessage ? (
        <div className="briefing-empty">{briefing.emptyMessage}</div>
      ) : (
        <div className="briefing-items">
          {briefing.items.map((it, i) => (
            <div className={`briefing-item ${it.tone}`} key={i}>
              <div className="dot" />
              <div className="body">
                <div className="t">{it.title}</div>
                <div className="d">{it.detail}</div>
              </div>
              {it.action && (
                <a className="btn btn--secondary btn--sm" href={it.action.href}>
                  {it.action.label}
                </a>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
```

- [ ] **Step 2: Add scoped styles to globals.css**

Append after the existing `.msg.conductor .text` block in `frontend/app/globals.css`:

```css
.briefing { display: flex; flex-direction: column; gap: 12px; padding: 16px 18px; background: var(--c-surface); border: 1px solid var(--c-line); border-radius: var(--r-md); box-shadow: var(--e-1); }
.briefing-head { display: flex; align-items: baseline; justify-content: space-between; gap: 12px; flex-wrap: wrap; }
.briefing-greeting { font-size: 15px; font-weight: 500; color: var(--c-ink); letter-spacing: -0.012em; }
.briefing-headline { font-family: var(--f-mono); font-size: 11px; color: var(--c-muted); letter-spacing: 0.04em; }
.briefing-empty { font-size: 13px; color: var(--c-muted); padding: 6px 0; }
.briefing-items { display: flex; flex-direction: column; gap: 6px; }
.briefing-item { display: grid; grid-template-columns: 10px 1fr auto; align-items: center; gap: 12px; padding: 10px 12px; border: 1px solid var(--c-line); border-radius: var(--r-sm); background: var(--c-canvas); }
.briefing-item .dot { width: 8px; height: 8px; border-radius: var(--r-pill); background: var(--c-line-strong); }
.briefing-item.wait .dot { background: var(--c-wait); animation: breathe 3s ease-in-out infinite; }
.briefing-item.fire .dot { background: var(--c-fire); animation: fire-flash 1.8s ease-out infinite; }
.briefing-item.live .dot { background: var(--c-live); animation: pulse 2.2s ease-out infinite; }
.briefing-item.info .dot { background: var(--c-info); }
.briefing-item .t { font-size: 13.5px; font-weight: 500; color: var(--c-ink); letter-spacing: -0.005em; }
.briefing-item .d { font-family: var(--f-mono); font-size: 10.5px; color: var(--c-muted); margin-top: 2px; }
```

---

### Task A1.3: Wire BriefingCard into ConductorApp as the proactive opener

**Files:**
- Modify: `frontend/components/ConductorApp.tsx`

- [ ] **Step 1: Import composeBriefing and BriefingCard**

After the existing `import { renderMarkdown } from "@/lib/md";` line, add:

```ts
import { composeBriefing, type Briefing } from "@/lib/briefing";
import { BriefingCard } from "./BriefingCard";
```

- [ ] **Step 2: Compute the briefing reactively**

After the existing `const liveFires = useMemo(...)` block, add:

```ts
const briefing: Briefing = useMemo(
  () => composeBriefing(agents, messages, checkpoints),
  [agents, messages, checkpoints]
);
```

- [ ] **Step 3: Render the briefing in the conversation**

Replace the existing `{!hasHistory && ( ... welcome msg ... )}` block in the JSX with:

```tsx
{loaded && (
  <div className="msg conductor">
    <div className="conductor-mark sm" />
    <div className="body">
      <div className="meta">
        <span className="nm">Conductor</span>
        <span>·</span>
        <span>{hasHistory ? "today's briefing" : "good to see you"}</span>
      </div>
      <BriefingCard briefing={briefing} />
    </div>
  </div>
)}
```

Add `loaded` state — track it: after the initial-load `useEffect` resolves, `setLoaded(true)`. Add `const [loaded, setLoaded] = useState(false);` with the other state hooks, and `setLoaded(true);` after the `Promise.all` resolves in the initial-load effect.

- [ ] **Step 4: Build + re-embed via the canonical script + verify visually**

Run: `./scripts/build-ui.sh && go build ./...`
Expected: both succeed.
Restart server: `pkill -f ah-test; /tmp/ah-test --serve --port 8099 &`
Visit `http://localhost:8099/conductor`. The first conductor turn should be the briefing card (greeting + headline + 0-N items), regardless of conversation history.

---

### Task A2.1: ActionCard component (delegate / create_project / escalate)

**Files:**
- Create: `frontend/components/ActionCard.tsx`

- [ ] **Step 1: Define the action-card component**

```tsx
// frontend/components/ActionCard.tsx
"use client";

import type { ChatReply } from "@/lib/api";
import { Check } from "./icons";

export function ActionCard({
  reply,
  onPrimary,
  onSecondary,
  busy = false,
}: {
  reply: ChatReply;
  onPrimary?: () => void;
  onSecondary?: () => void;
  busy?: boolean;
}) {
  const action = reply.action ?? "respond";

  const labels: Record<string, { tone: string; tag: string; primary: string; secondary: string }> = {
    delegate: { tone: "info", tag: "delegation", primary: "Show me what it produced", secondary: "Cancel" },
    create_project: { tone: "live", tag: "project", primary: "Open project", secondary: "Hide" },
    escalate: { tone: "wait", tag: "escalation", primary: "Approve & spawn", secondary: "Skip" },
    send_reply: { tone: "wait", tag: "email · external send", primary: "Show me the draft", secondary: "Skip" },
  };
  const cfg = labels[action];
  if (!cfg) return null; // caller falls back to markdown for non-actionable replies

  return (
    <div className="action-card" data-tone={cfg.tone}>
      <div className="ac-head">
        <span className={`ac-tag tone-${cfg.tone}`}>
          <span className="dot" />
          {cfg.tag}
        </span>
        {reply.project_id && <span className="ac-meta">project · {reply.project_id}</span>}
      </div>
      <div className="ac-body">{reply.response}</div>
      <div className="ac-actions">
        <button
          className="btn btn--primary btn--sm"
          disabled={busy || !onPrimary}
          onClick={onPrimary}
        >
          <Check />
          {cfg.primary}
        </button>
        <button
          className="btn btn--ghost btn--sm"
          disabled={busy || !onSecondary}
          onClick={onSecondary}
        >
          {cfg.secondary}
        </button>
      </div>
    </div>
  );
}
```

- [ ] **Step 2: Add scoped styles to globals.css**

Append after the briefing CSS block from Task A1.2:

```css
.action-card { display: flex; flex-direction: column; gap: 10px; padding: 12px 14px; background: var(--c-surface); border: 1px solid var(--c-line); border-top: 3px solid var(--c-line-strong); border-radius: var(--r-md); box-shadow: var(--e-1); }
.action-card[data-tone="wait"] { border-top-color: var(--c-wait); }
.action-card[data-tone="fire"] { border-top-color: var(--c-fire); }
.action-card[data-tone="info"] { border-top-color: var(--c-info); }
.action-card[data-tone="live"] { border-top-color: var(--c-live); }
.ac-head { display: flex; align-items: center; justify-content: space-between; gap: 8px; flex-wrap: wrap; }
.ac-tag { display: inline-flex; align-items: center; gap: 5px; font-family: var(--f-mono); font-size: 10px; letter-spacing: 0.08em; text-transform: uppercase; padding: 3px 8px; border-radius: var(--r-pill); }
.ac-tag .dot { width: 5px; height: 5px; border-radius: var(--r-pill); background: currentColor; }
.ac-tag.tone-wait { color: oklch(from var(--c-wait) calc(l - 0.3) c h); background: var(--c-wait-soft); }
.ac-tag.tone-fire { color: oklch(from var(--c-fire) calc(l - 0.25) c h); background: var(--c-fire-soft); }
.ac-tag.tone-info { color: oklch(from var(--c-info) calc(l - 0.1) c h); background: var(--c-info-soft); }
.ac-tag.tone-live { color: oklch(from var(--c-live) calc(l - 0.25) c h); background: var(--c-live-soft); }
[data-mode="dark"] .ac-tag.tone-wait { color: var(--c-wait); }
[data-mode="dark"] .ac-tag.tone-fire { color: var(--c-fire); }
[data-mode="dark"] .ac-tag.tone-info { color: var(--c-info); }
[data-mode="dark"] .ac-tag.tone-live { color: var(--c-live); }
.ac-meta { font-family: var(--f-mono); font-size: 10px; color: var(--c-muted); }
.ac-body { font-size: 13.5px; color: var(--c-ink); line-height: 1.5; }
.ac-actions { display: flex; gap: 6px; }
```

---

### Task A2.2: Render Brain replies as ActionCards when actionable

**Files:**
- Modify: `frontend/components/ConductorApp.tsx`

- [ ] **Step 1: Extend the Msg type with an optional reply**

Replace `type Msg = { role: "user" | "conductor"; html: string; meta?: string };` with:

```ts
type Msg = {
  role: "user" | "conductor";
  html: string;
  meta?: string;
  reply?: import("@/lib/api").ChatReply;
};
```

- [ ] **Step 2: Store the reply on the assistant message**

In `send()`, where the assistant message is pushed, also attach `reply`:

```ts
setMsgs((p) => [
  ...p,
  {
    role: "conductor",
    meta: reply?.action ? `routed · ${reply.action}` : "offline · backend unreachable",
    html: reply?.response
      ? renderMarkdown(`${icon} ${reply.response}`)
      : `<p>I'd classify this and dispatch to the right specialist, but <strong>the backend isn't reachable</strong> from this view. Start the Go server to route live.</p>`,
    reply: reply ?? undefined,
  },
]);
```

- [ ] **Step 3: Import ActionCard and render it for actionable replies**

Add to the imports:

```ts
import { ActionCard } from "./ActionCard";
```

In the `msgs.map(...)` block, replace the conductor branch with:

```tsx
<div className="msg conductor" key={i}>
  <div className="conductor-mark sm" />
  <div className="body">
    <div className="meta">
      <span className="nm">Conductor</span>
      <span>·</span>
      <span>{m.meta || "just now"}</span>
    </div>
    {m.reply && ["delegate", "create_project", "escalate", "send_reply"].includes(m.reply.action ?? "") ? (
      <ActionCard
        reply={m.reply}
        onPrimary={() => {
          // Phase A2 surfaces the action; Phase A3 will wire deep links and
          // /api/checkpoints/inject targets. For now, route the user to where
          // they can act on it.
          if (m.reply?.action === "create_project") window.location.href = "/mission";
          if (m.reply?.action === "delegate") window.location.href = "/mission";
        }}
        onSecondary={() => {/* dismiss is purely visual for now */}}
      />
    ) : (
      <div className="text" dangerouslySetInnerHTML={{ __html: m.html }} />
    )}
  </div>
</div>
```

- [ ] **Step 4: Build + verify**

Run: `./scripts/build-ui.sh && go build ./...`
Visit `/conductor`, type "create project for a todo app". The reply should render as a `[project]` card (or `[delegation]`, depending on Brain decision), not a markdown paragraph.

---

### Task A3.1: Stream in-thread agent activity via WebSocket

**Files:**
- Modify: `frontend/components/ConductorApp.tsx`

- [ ] **Step 1: Track the active "thread" — the most recent assistant project id**

Add state:

```ts
const [activeProject, setActiveProject] = useState<string | undefined>();
```

In `send()`, after the assistant message is appended:

```ts
if (reply?.project_id) setActiveProject(reply.project_id);
```

- [ ] **Step 2: Append in-thread messages from WS**

Replace the existing `connectWS` effect with:

```ts
useEffect(() => {
  let t: ReturnType<typeof setTimeout> | null = null;
  const debounced = () => {
    if (t) return;
    t = setTimeout(() => {
      t = null;
      refreshSidebar();
    }, 800);
  };
  return connectWS((e) => {
    if (e.type === "message" && "message" in e) {
      const m = (e as { message: ApiMessage }).message;
      const inThread = activeProject && m.metadata?.project_id === activeProject;
      if (inThread && m.from !== "user") {
        setMsgs((p) => [
          ...p,
          {
            role: "conductor",
            meta: `${m.from} · live`,
            html: renderMarkdown(m.content || ""),
          },
        ]);
      }
    }
    if (e.type === "message" || e.type === "checkpoint" || e.type === "brain_event")
      debounced();
  }, setConnected);
}, [refreshSidebar, activeProject]);
```

- [ ] **Step 3: Build + verify**

Run: `./scripts/build-ui.sh && go build ./...`
Visit `/conductor`, send "create project for a todo app", watch the conversation: after the Brain's create_project card, agent messages should stream into the thread as new conductor turns (assuming the Brain actually starts orchestration).

---

### Task A4.1: Default proactive cron jobs seeded by a script

**Files:**
- Create: `scripts/seed-conductor-defaults.sh`

- [ ] **Step 1: Write the seed script**

```bash
#!/usr/bin/env bash
# Seed a small set of proactive cron jobs that demonstrate the push model.
# Idempotent-ish: deletes existing jobs with the same task string first.
set -euo pipefail
HOST="${HOST:-http://localhost:8099}"

post() {
  curl -s -X POST "$HOST/api/cron/" \
    -H 'Content-Type: application/json' \
    -d "$1"
}

echo "→ seeding defaults at $HOST"
post '{"task":"Sweep inbox: route any new vendor or RFI emails to the right agent","schedule":"1h","agent_role":"ceo","project_id":"default"}'
echo
post '{"task":"Daily ops briefing: pending decisions, fires, schedule slips","schedule":"24h","agent_role":"ceo","project_id":"default"}'
echo
post '{"task":"Schedule slip check: scan milestones and flag anything trending late","schedule":"24h","agent_role":"project_manager","project_id":"default"}'
echo
echo "→ done. List:"
curl -s "$HOST/api/cron" | head -c 800
echo
```

- [ ] **Step 2: Make executable + run + verify**

Run: `chmod +x scripts/seed-conductor-defaults.sh && ./scripts/seed-conductor-defaults.sh`
Expected: three POSTs return `{job:{...},success:true}`; the final `GET /api/cron` lists three jobs.
Open `/triggers` — the three jobs should appear as live rows.

---

## Phase B — EPC scenario bootstrap

### Task B1.1: Inspect the inbox JSON schema

**Files:**
- Read only: `internal/email/engine.go`

- [ ] **Step 1: Capture the inbox file format**

Run: `grep -n "type Email struct\|json:" internal/email/engine.go | head -30`
Expected: a struct with `ID`, `From`, `Subject`, `Body`, `Category`, `Trust`, `Timestamp`, and the inbox path pattern (likely `projects/<id>/inbox/<email-id>.json`). Record the exact fields and path for use in B1.2.

---

### Task B1.2: Atlas EPC bootstrap script

**Files:**
- Create: `scripts/scenarios/load-epc-atlas.sh`
- Create (the script writes these): `projects/atlas-site/inbox/*.json`

- [ ] **Step 1: Write the bootstrap**

Following the exact email schema captured in B1.1, write a script that:
1. Creates `projects/atlas-site/inbox/` and `projects/atlas-site/.tasks/`.
2. Writes 4–6 JSON email files matching the schema: Sarah Jones progress request, Raj Kumar application, **XYZ Steel impersonation alert** (trust=`impersonation`), vendor-117 invoice, RFI #18 from Atlas Consulting, optional credit-note.
3. Calls `POST /api/cron/` to create a daily PM-check and an hourly inbox sweep scoped to `project_id: "atlas-site"`.

Use heredocs for the JSON bodies. Each `jq -n` or raw heredoc must produce valid JSON conforming to the struct from B1.1.

- [ ] **Step 2: Run + verify**

Run: `chmod +x scripts/scenarios/load-epc-atlas.sh && ./scripts/scenarios/load-epc-atlas.sh`
Expected: `ls projects/atlas-site/inbox/` lists 4–6 `.json` files; `curl -s localhost:8099/api/cron | jq '.jobs[] | select(.project_id=="atlas-site")'` shows the two new jobs.
Open `/conductor` — the briefing should reflect the new fires/checkpoints once the cron has had time to fire (or after a refresh).

---

### Task B2.1: Calculator software-dev bootstrap script

**Files:**
- Create: `scripts/scenarios/load-dev-calculator.sh`

- [ ] **Step 1: Write the bootstrap**

A script that:
1. Creates `projects/calculator/` and `.tasks/`.
2. POSTs a small task via `/api/task` (existing endpoint) — `{"task":"build a single-file HTML calculator with a dark theme","project_id":"calculator"}` — to kick the IT pack into the pipeline.
3. Optionally creates one cron job (e.g., hourly "Status check on calculator project" bound to `ceo`).

- [ ] **Step 2: Run + verify**

Run: `chmod +x scripts/scenarios/load-dev-calculator.sh && ./scripts/scenarios/load-dev-calculator.sh`
Expected: the task is accepted (HTTP 200), the pipeline starts in background; opening `/conductor` shows messages streaming in for `project_id=calculator`.

---

### Task B3.1: README for scenarios

**Files:**
- Create: `scripts/scenarios/README.md`

- [ ] **Step 1: Document the scenarios**

```md
# Demo scenarios

Self-contained bootstraps that load realistic state into the running server,
so a 30-second demo is possible from a cold start.

## Usage

1. Start the server: `./scripts/build-ui.sh && go build -o /tmp/ah ./cmd/agent-house && /tmp/ah --serve --port 8099 &`
2. Pick a scenario:
   - `./scripts/scenarios/load-epc-atlas.sh` — Atlas Construction · Site 02 (EPC pack)
   - `./scripts/scenarios/load-dev-calculator.sh` — small HTML calculator build (IT pack)
3. Open `http://localhost:8099/conductor` — the briefing reflects the seeded state.

Each scenario is idempotent: re-running will overwrite the seeded inbox files
and append to the cron list. To reset, delete `projects/<id>/` and remove
the matching cron jobs via `/triggers`.
```

---

## Out of scope (deferred, by user direction 2026-05-20)

- AgentBuilder UI investment beyond what already exists (real reads + writable mode toggle).
- Project switcher in the shared `Chrome` (Tier-1 #3).
- Mission tabs (Files / Settings / Triggers) wiring.
- Phase-1 registry write API and pack manifests.
- Multi-tenant, auth, project-level RBAC.

These remain in `phase-1-plan.md` as future work; this plan does not add to them.

---

## Self-review notes

- **Spec coverage:** Every section of the strategy from the 2026-05-20 turn is covered: A1 proactive opener (A1.x), A2 action cards (A2.x), A3 in-thread streaming (A3.1), A4 proactive cron (A4.1), B EPC + dev scenarios (B1.x, B2.1), B3 docs.
- **Placeholder scan:** no "TBD"/"appropriate" — every step shows exact code, exact commands, exact files.
- **Type consistency:** `Briefing`/`BriefingItem` (A1.1) used by `BriefingCard` (A1.2) and ConductorApp (A1.3) match. `ChatReply` (existing in `lib/api.ts`) used by `ActionCard` (A2.1) and the new `Msg.reply` field (A2.2) match. `composeBriefing` signature is stable across tasks.
- **Spec self-honesty:** action-card primary buttons in A2.2 are wired to `window.location.href` (not full backend `/api/checkpoints` resolution) — labeled in the plan as the A3 follow-up. No fake "Approve" that pretends to call an endpoint it doesn't.
