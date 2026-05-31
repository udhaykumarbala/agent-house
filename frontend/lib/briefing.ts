import type { ApiAgent, ApiMessage, ApiCheckpoint } from "./types";

/**
 * The proactive opener model. Compiled from live backend state (agents,
 * messages, checkpoints), this is what the Conductor *says first*, every
 * time, instead of an empty composer waiting for a prompt.
 */
export interface BriefingItem {
  tone: "wait" | "fire" | "info" | "live";
  title: string;          // short imperative
  detail: string;         // one-line context
  action?: { label: string; href?: string };
}

export interface Briefing {
  greeting: string;
  headline: string;
  items: BriefingItem[];
  emptyMessage?: string;
}

function timeOfDay(d = new Date()): string {
  const h = d.getHours();
  if (h < 5) return "Working late";
  if (h < 12) return "Good morning";
  if (h < 17) return "Good afternoon";
  return "Good evening";
}

const TRIGGER_TYPES = new Set([
  "trigger",
  "trigger_fired",
  "hook",
  "cron",
  "webhook",
]);

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
    .slice(0, 24);
  const triggerFires = recent.filter((m) => TRIGGER_TYPES.has(m.type));
  const working = agents.filter((a) => a.active);
  const idle = agents.length - working.length;

  const items: BriefingItem[] = [];

  // Pending decisions first — the most consequential proactive items.
  for (const cp of pending.slice(0, 3)) {
    items.push({
      tone: "wait",
      title: cp.artifact_summary || cp.type,
      detail: `${cp.type} · task#${(cp.task_id || "").slice(0, 8)}`,
      action: { label: "Review", href: "/mission" },
    });
  }
  // Trigger-driven work the system did on its own.
  if (triggerFires.length > 0) {
    const recentTypes = triggerFires
      .slice(0, 3)
      .map((m) => m.type)
      .join(" · ");
    items.push({
      tone: "fire",
      title: `${triggerFires.length} trigger ${
        triggerFires.length === 1 ? "fire" : "fires"
      } today`,
      detail: recentTypes,
      action: { label: "See triggers", href: "/triggers" },
    });
  }
  // Who's mid-flight right now.
  if (working.length > 0) {
    items.push({
      tone: "info",
      title: `${working.length} agent${working.length === 1 ? "" : "s"} working now`,
      detail: working
        .map((a) => a.name || a.role)
        .slice(0, 3)
        .join(" · "),
    });
  }

  return {
    greeting: `${timeOfDay()}, ${user}.`,
    headline: `${pending.length} pending · ${triggerFires.length} fires today · ${working.length} working · ${idle} idle`,
    items,
    emptyMessage:
      items.length === 0
        ? "Nothing requires you right now. I'll surface anything that does."
        : undefined,
  };
}
