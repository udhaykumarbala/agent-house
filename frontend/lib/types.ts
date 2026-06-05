// Mirrors the Go backend JSON shapes (internal/message, internal/checkpoint,
// internal/web handleAgents) plus a UI view-model the reference is built on.

export type AgentState =
  | "working"
  | "awaiting"
  | "live"
  | "idle"
  | "offline"
  | "fired";

export interface ApiPermission {
  CanCreateFiles?: boolean;
  CanModifyFiles?: boolean;
  CanExecuteCode?: boolean;
  AllowedExtensions?: string[];
  AllowedPaths?: string[];
  DeniedPaths?: string[];
}

export interface ApiAgent {
  role: string;
  name: string;
  active: boolean;
  color?: string;
  permissions?: ApiPermission;
}

export interface ApiMessage {
  id: string;
  type: string;
  from: string;
  to: string;
  content: string;
  metadata?: {
    project_id?: string;
    task_id?: string;
    priority?: string;
    tags?: string[];
    [k: string]: unknown;
  };
  timestamp: string;
}

export interface ApiCheckpoint {
  id: string;
  project_id: string;
  task_id: string;
  type: string;
  status: string; // "pending" | "resolved"
  artifact_path?: string;
  artifact_summary?: string;
  created_at: string;
  resolved_at?: string | null;
  // Annotations from /api/checkpoints/all (per-project run mode + timeout).
  project?: string;
  run_mode?: string;
  decision_timeout_minutes?: number;
  decision?: {
    action: string;
    feedback?: string;
    decided_by: string;
    decided_at?: string;
  } | null;
}

// ── UI view-models (what the reference components render) ──

export interface AgentVM {
  initials: string;
  name: string;
  role: string;
  state: AgentState;
  group: string;
  badgeTone?: "wait" | "info" | "fire";
  badgeCount?: number;
}

export interface HitlVM {
  id: string;
  initials: string;
  state: AgentState;
  title: string;
  detail: string;
  tag: string;
  timer: string;
  roleTag: string;
  desc: string;
  payload: [string, string, "" | "num" | "str"][];
  trace: string;
  primary: string;
  live: boolean; // backed by a real checkpoint id
  /** Full-auto only: epoch ms when the CEO auto-decides if the human hasn't. */
  deadline?: number;
}

export interface FeedVM {
  id: string;
  tone: "" | "fire" | "live" | "info" | "wait";
  triggered?: "fire" | "live";
  ts: string;
  head: string[];
  body: string; // may contain <strong> / <code>
  card?: { title: string; detail: string };
}

export interface StatVM {
  tone: "" | "fire" | "wait";
  n: string;
  l: string;
}

// ── Scenarios (capability-composing flows) ─────────────────────────
// These mirror internal/scenario.* — kept minimal and shape-only so the
// UI can render *any* scenario the backend declares, without per-name
// branches anywhere.

export interface ScenarioMeta {
  name: string;
  description: string;
  example: Record<string, unknown>;
}

export interface ScenarioStep {
  agent: string;
  action: string;
  input?: Record<string, unknown>;
  output?: Record<string, unknown>;
  ok: boolean;
  note?: string;
}

export interface ScenarioSuggestion {
  title: string;
  detail: string;
  action: string;
  payload?: Record<string, unknown>;
}

export interface ScenarioResult {
  scenario: string;
  scope: string;
  steps: ScenarioStep[];
  summary: string;
  suggestions: ScenarioSuggestion[];
  ok: boolean;
}
