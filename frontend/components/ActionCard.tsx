"use client";

import type { ChatReply } from "@/lib/api";
import { humanize } from "@/lib/humanize";
import { Check } from "./icons";

/**
 * Render a Brain reply as a typed action card when the action is one of the
 * actionable kinds. The plain-text `respond` action falls back to markdown
 * in the caller — this component returns null for unknown actions.
 *
 * Phase A2 surfaces the action and routes the user to where they can act on
 * it. Phase A3+ will wire `onPrimary` to real backend mutations
 * (/api/checkpoints, /api/inject) directly from the card.
 */
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

  const cfg: Record<
    string,
    { tone: "info" | "live" | "wait" | "fire"; tag: string; primary: string; secondary: string }
  > = {
    delegate: {
      tone: "info",
      tag: "delegation",
      primary: "Open mission",
      secondary: "Dismiss",
    },
    create_project: {
      tone: "live",
      tag: "project · started",
      primary: "Open project",
      secondary: "Hide",
    },
    escalate: {
      tone: "wait",
      tag: "escalation",
      primary: "Approve & spawn",
      secondary: "Skip",
    },
    send_reply: {
      tone: "wait",
      tag: "email · external send",
      primary: "Show me the draft",
      secondary: "Skip",
    },
  };

  const c = cfg[action];
  if (!c) return null; // caller renders markdown fallback

  return (
    <div className="action-card" data-tone={c.tone}>
      <div className="ac-head">
        <span className={`ac-tag tone-${c.tone}`}>
          <span className="dot" />
          {c.tag}
        </span>
        {reply.project_id && (
          <span className="ac-meta">project · {reply.project_id}</span>
        )}
      </div>
      <div className="ac-body">{reply.response}</div>
      <div className="ac-actions">
        <button
          className="btn btn--primary btn--sm"
          disabled={busy || !onPrimary}
          onClick={onPrimary}
        >
          <Check />
          {c.primary}
        </button>
        <button
          className="btn btn--ghost btn--sm"
          disabled={busy || !onSecondary}
          onClick={onSecondary}
        >
          {c.secondary}
        </button>
      </div>
    </div>
  );
}
