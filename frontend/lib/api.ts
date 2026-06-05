import type { ApiAgent, ApiMessage, ApiCheckpoint } from "./types";

/**
 * When the bundle is embedded in the Go binary it is served same-origin,
 * so the base is "". In `next dev` point at the Go server via
 * NEXT_PUBLIC_API_BASE (e.g. http://localhost:8080).
 */
export const API_BASE =
  process.env.NEXT_PUBLIC_API_BASE?.replace(/\/$/, "") ?? "";

export const DEFAULT_PROJECT =
  process.env.NEXT_PUBLIC_PROJECT ?? "default";

function wsUrl(): string {
  if (typeof window === "undefined") return "";
  const base = API_BASE || window.location.origin;
  return base.replace(/^http/, "ws").replace(/\/$/, "") + "/ws";
}

async function getJSON<T>(path: string): Promise<T | null> {
  try {
    const res = await fetch(`${API_BASE}${path}`, { cache: "no-store" });
    if (!res.ok) return null;
    return (await res.json()) as T;
  } catch {
    return null; // backend not reachable (dev without server) — caller falls back
  }
}

export async function fetchAgents(): Promise<ApiAgent[]> {
  const data = await getJSON<{ agents: ApiAgent[] }>("/api/agents");
  return data?.agents ?? [];
}

/**
 * Fetch recent messages. Pass `project = ""` to get everything across
 * scopes — Conductor uses this so workforce activity from a scenario
 * running under "atlas-site" still flows in.
 */
export async function fetchMessages(project = ""): Promise<ApiMessage[]> {
  // `/api/messages` returns {count, messages: ApiMessage[]|null}. Older
  // assumptions of a flat array silently produced an empty list — the bug
  // that made the Conductor briefing show "0 fires today" even when the
  // store had real messages.
  const qs = project ? `?project=${encodeURIComponent(project)}` : "";
  const data = await getJSON<
    | ApiMessage[]
    | { count?: number; messages?: ApiMessage[] | null }
  >(`/api/messages${qs}`);
  if (!data) return [];
  if (Array.isArray(data)) return data;
  return data.messages ?? [];
}

export async function fetchCheckpoints(
  project = DEFAULT_PROJECT
): Promise<ApiCheckpoint[]> {
  const data = await getJSON<{ checkpoints: ApiCheckpoint[] } | ApiCheckpoint[]>(
    `/api/checkpoints?project=${encodeURIComponent(project)}`
  );
  if (!data) return [];
  return Array.isArray(data) ? data : (data.checkpoints ?? []);
}

/** Aggregate checkpoints across every project (each annotated with its run mode
 * + decision timeout). Lets Mission follow whichever build is running. */
export async function fetchAllCheckpoints(): Promise<ApiCheckpoint[]> {
  const data = await getJSON<{ checkpoints: ApiCheckpoint[] }>(
    `/api/checkpoints/all`
  );
  return data?.checkpoints ?? [];
}

export async function resolveCheckpoint(
  id: string,
  decision: "approved" | "rejected",
  feedback = ""
): Promise<boolean> {
  try {
    const res = await fetch(`${API_BASE}/api/checkpoints/${id}`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ decision, feedback }),
    });
    return res.ok;
  } catch {
    return false;
  }
}

export interface ChatHistMsg {
  role: "user" | "assistant";
  content: string;
  timestamp: number;
}

// ── Named conversations ────────────────────────────────────────────

export interface ConversationMeta {
  id: string;
  name: string;
  message_count: number;
  last_at: number;
  first_at: number;
  preview: string;
}

export interface ConversationHit {
  conversation_id: string;
  name: string;
  role: "user" | "assistant";
  snippet: string;
  timestamp: number;
  match_count?: number; // total messages in this conv matching the query
}

export async function fetchConversations(): Promise<ConversationMeta[]> {
  const data = await getJSON<{ conversations: ConversationMeta[] }>(
    "/api/conversations"
  );
  return data?.conversations ?? [];
}

export async function createConversation(): Promise<ConversationMeta | null> {
  try {
    const res = await fetch(`${API_BASE}/api/conversations`, {
      method: "POST",
    });
    if (!res.ok) return null;
    return (await res.json()) as ConversationMeta;
  } catch {
    return null;
  }
}

export interface ConversationSearchResult {
  mode: "semantic" | "lexical";
  hits: ConversationHit[];
}

/**
 * Inject a synthetic email into the global inbox. The Lab page uses this
 * to make scenarios demoable without a real SMTP server. Posts to
 * /api/cap/email/inbox and writes a JSON file to projects/inbox/.
 */
export async function injectMockEmail(
  payload: Record<string, unknown>
): Promise<{ ok: boolean; id?: string; error?: string }> {
  try {
    const res = await fetch(`${API_BASE}/api/cap/email/inbox`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    if (!res.ok) {
      const text = await res.text();
      return { ok: false, error: text };
    }
    const saved = (await res.json()) as { id?: string };
    return { ok: true, id: saved.id };
  } catch (e) {
    return { ok: false, error: (e as Error).message };
  }
}

export async function searchConversations(
  q: string
): Promise<ConversationSearchResult> {
  if (!q.trim()) return { mode: "lexical", hits: [] };
  const data = await getJSON<{ mode?: string; hits?: ConversationHit[] }>(
    `/api/conversations/search?q=${encodeURIComponent(q)}`
  );
  return {
    mode: (data?.mode === "semantic" ? "semantic" : "lexical"),
    hits: data?.hits ?? [],
  };
}

/**
 * Sometimes the Brain's stored history contains a full envelope as the
 * assistant's content (`{"action":"respond","response":"…","suggestions":[…]}`).
 * That's a backend artefact — it should be plain text. We defensively
 * unwrap it client-side so the rendered conversation only shows the inner
 * reply. Idempotent on already-clean text.
 */
function unwrapEnvelope(s: string): string {
  if (!s) return s;
  const t = s.trim();
  if (!(t.startsWith("{") && t.includes('"response"'))) return s;
  try {
    const obj = JSON.parse(t);
    if (obj && typeof obj.response === "string") {
      // Recurse — sometimes there are two levels of nesting.
      return unwrapEnvelope(obj.response);
    }
  } catch {
    /* not JSON; render as-is */
  }
  return s;
}

/** Load Brain conversation history. `conv` is the conversation id (the
 * same string the backend stores as `user_id` internally). Each named
 * conversation is one file under `.brain/conversations/<conv>.json`. */
export async function fetchChatHistory(
  conv = "default"
): Promise<ChatHistMsg[]> {
  const data = await getJSON<{ messages: ChatHistMsg[] }>(
    `/api/chat?user=${encodeURIComponent(conv)}`
  );
  const msgs = data?.messages ?? [];
  return msgs.map((m) =>
    m.role === "assistant" ? { ...m, content: unwrapEnvelope(m.content) } : m
  );
}

export interface ChatReply {
  action?: string;
  response?: string;
  success?: boolean;
  project_id?: string;
  suggestions?: string[];
}

// ── Run-mode / workflow settings ──
export interface WorkflowSettings {
  project_id?: string;
  run_mode: string; // manual | semi_auto | full_auto | blitz
  decision_timeout_minutes: number;
  require_template_approval?: boolean;
  require_plan_approval?: boolean;
  require_phase_gate?: boolean;
  require_final_acceptance?: boolean;
}

/** The global default run mode new builds inherit. `_global` is a pseudo-project
 * that holds the default WorkflowSettings via the same checkpoint settings API. */
export async function getRunMode(
  project = "_global"
): Promise<WorkflowSettings | null> {
  return getJSON<WorkflowSettings>(
    `/api/settings/workflow?project=${encodeURIComponent(project)}`
  );
}

export async function setRunMode(
  run_mode: string,
  decision_timeout_minutes: number,
  project = "_global"
): Promise<boolean> {
  try {
    const res = await fetch(
      `${API_BASE}/api/settings/workflow?project=${encodeURIComponent(project)}`,
      {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ run_mode, decision_timeout_minutes }),
      }
    );
    return res.ok;
  } catch {
    return false;
  }
}

/** Send a message through the Brain router (POST /api/chat). `conv` is the
 * conversation id (kept as `user_id` in the wire payload for backend
 * compatibility — that field is what the store keys on). */
export async function sendChat(
  message: string,
  conv = "default",
  scope?: string
): Promise<ChatReply | null> {
  try {
    const res = await fetch(`${API_BASE}/api/chat`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ message, user_id: conv, scope }),
    });
    if (!res.ok) return null;
    const raw = (await res.json()) as ChatReply;
    // Unwrap the same double-envelope artefact that can show up in fresh
    // replies if the Brain's prompt history was tainted earlier.
    if (raw?.response) raw.response = unwrapEnvelope(raw.response);
    return raw;
  } catch {
    return null;
  }
}

// ── Cron triggers (real, working scheduler) ──

export interface CronJob {
  id: string;
  schedule: string;
  task: string;
  agent_role: string;
  project_id: string;
  enabled: boolean;
  last_run?: string;
  next_run?: string;
}

export async function fetchCron(): Promise<CronJob[]> {
  const data = await getJSON<{ jobs: CronJob[] }>("/api/cron");
  return data?.jobs ?? [];
}

export async function createCron(body: {
  task: string;
  schedule: string;
  agent_role: string;
  project_id: string;
}): Promise<boolean> {
  try {
    const res = await fetch(`${API_BASE}/api/cron/`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    });
    return res.ok;
  } catch {
    return false;
  }
}

export async function deleteCron(id: string): Promise<boolean> {
  try {
    const res = await fetch(
      `${API_BASE}/api/cron/?id=${encodeURIComponent(id)}`,
      { method: "DELETE" }
    );
    return res.ok;
  } catch {
    return false;
  }
}

// ── Agent execution modes (real, writable) ──

export async function fetchAgentModes(): Promise<{
  modes: Record<string, string>;
  defaults: Record<string, string>;
}> {
  const data = await getJSON<{
    modes: Record<string, string>;
    defaults: Record<string, string>;
  }>("/api/agents/modes");
  return { modes: data?.modes ?? {}, defaults: data?.defaults ?? {} };
}

export async function setAgentMode(
  role: string,
  mode: "oneshot" | "session"
): Promise<boolean> {
  try {
    const res = await fetch(`${API_BASE}/api/agents/modes`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ role, mode }),
    });
    return res.ok;
  } catch {
    return false;
  }
}

// ── Scenarios (capability-composing flows) ──

import type { ScenarioMeta, ScenarioResult } from "./types";

export async function fetchScenarios(): Promise<ScenarioMeta[]> {
  const data = await getJSON<{ scenarios: ScenarioMeta[] }>("/api/scenarios");
  return data?.scenarios ?? [];
}

export async function runScenario(
  name: string,
  input: Record<string, unknown>,
  scope = "default"
): Promise<ScenarioResult | null> {
  try {
    const res = await fetch(
      `${API_BASE}/api/scenario/${encodeURIComponent(name)}/run?scope=${encodeURIComponent(scope)}`,
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(input),
      }
    );
    if (!res.ok) return null;
    return (await res.json()) as ScenarioResult;
  } catch {
    return null;
  }
}

export type WSEnvelope =
  | { type: "message"; message: ApiMessage }
  | { type: "checkpoint"; event: unknown }
  | { type: "agent_task_event"; event: unknown }
  | { type: string; [k: string]: unknown };

/** Reconnecting WebSocket to the Go hub at /ws. */
export function connectWS(
  onEvent: (e: WSEnvelope) => void,
  onStatus: (connected: boolean) => void
): () => void {
  let ws: WebSocket | null = null;
  let closed = false;
  let retry: ReturnType<typeof setTimeout> | null = null;

  const open = () => {
    if (closed) return;
    try {
      ws = new WebSocket(wsUrl());
    } catch {
      schedule();
      return;
    }
    ws.onopen = () => onStatus(true);
    ws.onclose = () => {
      onStatus(false);
      schedule();
    };
    ws.onerror = () => ws?.close();
    ws.onmessage = (ev) => {
      try {
        onEvent(JSON.parse(ev.data) as WSEnvelope);
      } catch {
        /* ignore malformed frames */
      }
    };
  };

  const schedule = () => {
    if (closed || retry) return;
    retry = setTimeout(() => {
      retry = null;
      open();
    }, 2500);
  };

  open();

  return () => {
    closed = true;
    if (retry) clearTimeout(retry);
    ws?.close();
  };
}
