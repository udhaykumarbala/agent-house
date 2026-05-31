import type {
  ApiAgent,
  ApiMessage,
  ApiCheckpoint,
  AgentVM,
  AgentState,
  HitlVM,
  FeedVM,
  StatVM,
} from "./types";
import { humanize } from "./humanize";

function initials(name: string, role: string): string {
  const src = (name || role || "?").trim();
  const parts = src.split(/[\s_·]+/).filter(Boolean);
  if (parts.length >= 2) return (parts[0][0] + parts[1][0]).toUpperCase();
  return src.slice(0, 2).toUpperCase();
}

function esc(s: string): string {
  return s.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
}

function hhmm(iso: string): string {
  const d = new Date(iso);
  return isNaN(d.getTime())
    ? ""
    : d.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit", hour12: false });
}

const PLATFORM = new Set(["ceo"]);

/** Group live agents into the reference's sidebar buckets. */
export function mapAgents(
  agents: ApiAgent[],
  pendingRoles: Set<string>
): { label: string; agents: AgentVM[] }[] {
  const groups: Record<string, AgentVM[]> = {
    Working: [],
    "Awaiting you": [],
    "Live · idle": [],
    Platform: [],
    Disabled: [],
  };

  for (const a of agents) {
    const platform = PLATFORM.has(a.role);
    const awaiting = pendingRoles.has(a.role);
    const vm: AgentVM = {
      initials: initials(a.name, a.role),
      name: a.name || a.role,
      role: a.role,
      state: !a.active
        ? platform
          ? "live"
          : "idle"
        : awaiting
        ? "awaiting"
        : "working",
      group: "",
    };
    if (platform) groups.Platform.push({ ...vm, group: "Platform" });
    else if (awaiting)
      groups["Awaiting you"].push({
        ...vm,
        group: "Awaiting you",
        badgeTone: "wait",
        badgeCount: 1,
      });
    else if (a.active)
      groups.Working.push({ ...vm, group: "Working", badgeTone: "info" });
    else groups["Live · idle"].push({ ...vm, group: "Live · idle" });
  }

  return Object.entries(groups)
    .filter(([, v]) => v.length > 0)
    .map(([label, ag]) => ({ label, agents: ag }));
}

const PROFESSIONAL = (role: string) =>
  `role:${role.replace(/[^a-z_]/gi, "") || "agent"}`;

export function mapCheckpoints(cps: ApiCheckpoint[]): HitlVM[] {
  return cps
    .filter((c) => c.status === "pending")
    .map((c) => {
      const summary = c.artifact_summary || c.type;
      return {
        id: c.id,
        initials: (c.type || "CP").slice(0, 2).toUpperCase(),
        state: "awaiting",
        title: summary,
        detail: `${c.type} · task#${(c.task_id || "").slice(0, 6)}`,
        tag: c.type.toUpperCase(),
        timer: "pending",
        roleTag: PROFESSIONAL(c.type),
        desc: c.artifact_summary
          ? esc(c.artifact_summary)
          : `Checkpoint <code>${esc(c.type)}</code> awaiting your decision.`,
        payload: [
          ["type", c.type, "str"],
          ["task_id", c.task_id || "", "str"],
          ["artifact", c.artifact_path || "—", "str"],
          ["status", c.status, ""],
        ],
        trace: `checkpoint#${c.id} · project ${c.project_id} · all events audited`,
        primary: "Approve",
        live: true,
      };
    });
}

const TRIGGER_TYPES = new Set(["trigger", "trigger_fired", "hook", "cron", "webhook"]);

export function mapMessages(msgs: ApiMessage[]): FeedVM[] {
  return [...msgs]
    .sort(
      (a, b) =>
        new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime()
    )
    .slice(0, 40)
    .map((m) => {
      const triggered = TRIGGER_TYPES.has(m.type);
      const tone: FeedVM["tone"] = triggered
        ? "fire"
        : m.type === "checkpoint" || m.metadata?.priority === "high"
        ? "wait"
        : m.from === "user"
        ? ""
        : "info";
      const who = m.from === "user" ? "you" : m.from || "agent";
      return {
        id: m.id,
        tone,
        triggered: triggered ? "fire" : undefined,
        ts: hhmm(m.timestamp),
        head: triggered ? [`${m.type} fired`] : [who, m.type],
        body:
          `<strong>${esc(who)}</strong> ` +
          esc(m.content || "").slice(0, 320),
      };
    });
}

// ── Conductor sidebar ──

export interface WfRow {
  ini: string;
  name: string;
  sub: string;
  state: AgentState;
  rowTone: "" | "wait" | "fire";
}

/** ACTIVE_WINDOW_MS — an agent that emitted a message inside this window is
 * considered "recently working" even when /api/agents reports active=false.
 * This lets scenario step emits visibly pulse the workforce panel (each
 * scenario step is a single message, but it's the strongest live signal
 * we have until orchestrator activity reporting is real-time). */
const ACTIVE_WINDOW_MS = 60_000;

export function mapWorkforce(
  agents: ApiAgent[],
  pendingRoles: Set<string>,
  messages: ApiMessage[] = []
): WfRow[] {
  // Index recent-by-role: most recent timestamp per `from` (role).
  const now = Date.now();
  const recent: Record<string, number> = {};
  for (const m of messages) {
    if (!m.from || m.from === "user" || m.from === "system") continue;
    const t = new Date(m.timestamp).getTime();
    if (isNaN(t) || now - t > ACTIVE_WINDOW_MS) continue;
    if (!recent[m.from] || t > recent[m.from]) recent[m.from] = t;
  }

  return agents.map((a) => {
    const platform = PLATFORM.has(a.role);
    const awaiting = pendingRoles.has(a.role);
    const recentlyActive = !!recent[a.role];
    const live = a.active || recentlyActive;

    const state: AgentState = !live
      ? platform
        ? "live"
        : "idle"
      : awaiting
      ? "awaiting"
      : "working";

    let sub = "idle";
    if (platform) sub = "platform";
    else if (awaiting) sub = "await · 1";
    else if (a.active) sub = "working";
    else if (recentlyActive) {
      const secs = Math.max(1, Math.round((now - recent[a.role]) / 1000));
      sub = `active · ${secs}s ago`;
    }

    return {
      ini: initials(a.name, a.role),
      // If the API didn't give a display name, derive one from the role id
      // ("project_manager" → "Project Manager") so the UI never shows snake_case.
      name: a.name || humanize(a.role),
      sub,
      state,
      rowTone: awaiting ? "wait" : "",
    };
  });
}

export interface FireRow {
  tone: "" | "fire" | "live";
  t: string;
  d: string;
  ts: string;
}

const REL = (iso: string): string => {
  const d = new Date(iso).getTime();
  if (isNaN(d)) return "";
  const s = Math.max(0, (Date.now() - d) / 1000);
  if (s < 90) return `${Math.round(s)}s`;
  if (s < 5400) return `${Math.round(s / 60)}m`;
  if (s < 129600) return `${Math.round(s / 3600)}h`;
  return `${Math.round(s / 86400)}d`;
};

export function mapFires(msgs: ApiMessage[]): FireRow[] {
  return [...msgs]
    .sort(
      (a, b) =>
        new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime()
    )
    .slice(0, 6)
    .map((m) => {
      const triggered = TRIGGER_TYPES.has(m.type);
      return {
        tone: triggered ? "fire" : m.from !== "user" ? "live" : "",
        t: (m.content || m.type).slice(0, 42),
        d: `${m.from || "system"} → ${m.to || "all"}`,
        ts: REL(m.timestamp),
      };
    });
}

export function statsFrom(
  agents: ApiAgent[],
  hitl: HitlVM[],
  feed: FeedVM[]
): StatVM[] {
  const fires = feed.filter((f) => f.triggered).length;
  const working = agents.filter((a) => a.active).length;
  return [
    { tone: "wait", n: String(hitl.length), l: "await you" },
    { tone: "fire", n: String(fires), l: "trigger fires" },
    { tone: "", n: String(working), l: "working" },
    { tone: "", n: "0", l: "errors · 7d" },
  ];
}
