"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { Chrome } from "./Chrome";
import { Search, Plus, Bolt, Close } from "./icons";
import { fetchCron, createCron, deleteCron, type CronJob } from "@/lib/api";

type TrigType = "cron" | "email" | "webhook" | "filewatch";

const GLYPHS: Record<TrigType, React.ReactNode> = {
  email: (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round">
      <rect x="3" y="5" width="18" height="14" rx="2" /><path d="M3 7l9 7 9-7" />
    </svg>
  ),
  webhook: (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round">
      <circle cx="8" cy="16" r="3" /><circle cx="16" cy="8" r="3" /><path d="M10 14l4-4" />
    </svg>
  ),
  cron: (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round">
      <circle cx="12" cy="12" r="9" /><path d="M12 7v5l3 2" />
    </svg>
  ),
  filewatch: (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round">
      <path d="M5 4h14v16H5z" /><path d="M5 9h14M9 4v16" />
    </svg>
  ),
};

function rel(s?: string): string {
  if (!s) return "never";
  const d = new Date(s).getTime();
  if (isNaN(d)) return s;
  const sec = Math.max(0, (Date.now() - d) / 1000);
  if (sec < 90) return `${Math.round(sec)}s ago`;
  if (sec < 5400) return `${Math.round(sec / 60)}m ago`;
  if (sec < 129600) return `${Math.round(sec / 3600)}h ago`;
  return `${Math.round(sec / 86400)}d ago`;
}

export function TriggersApp() {
  const [jobs, setJobs] = useState<CronJob[]>([]);
  const [loaded, setLoaded] = useState(false);
  const [expanded, setExpanded] = useState<string | null>(null);
  const [armed, setArmed] = useState(false);
  const [busy, setBusy] = useState(false);
  const [flash, setFlash] = useState<string | null>(null);

  // new-trigger form
  const [fName, setFName] = useState("Daily schedule check");
  const [fSchedule, setFSchedule] = useState("24h");
  const [fRole, setFRole] = useState("project_manager");
  const [fProject, setFProject] = useState("default");

  const load = useCallback(() => {
    fetchCron().then((j) => {
      setJobs(j);
      setLoaded(true);
    });
  }, []);
  useEffect(() => {
    load();
    const t = setInterval(load, 15000);
    return () => clearInterval(t);
  }, [load]);

  const counts = useMemo(
    () => ({
      cron: jobs.length,
      email: 0,
      webhook: 0,
      filewatch: 0,
    }),
    [jobs]
  );

  async function save() {
    if (!fName.trim() || !fSchedule.trim()) return;
    setBusy(true);
    const ok = await createCron({
      task: fName.trim(),
      schedule: fSchedule.trim(),
      agent_role: fRole.trim() || "senior_dev",
      project_id: fProject.trim() || "default",
    });
    setBusy(false);
    if (ok) {
      setFlash("Trigger saved · scheduler will fire it on schedule");
      setTimeout(() => setFlash(null), 4000);
      load();
    } else {
      setFlash("Save failed · backend unreachable");
      setTimeout(() => setFlash(null), 4000);
    }
  }

  async function remove(id: string) {
    setBusy(true);
    const ok = await deleteCron(id);
    setBusy(false);
    if (ok) {
      setExpanded(null);
      load();
    }
  }

  return (
    <div className="app lay-triggers" data-screen-label="Triggers Console">
      <Chrome pageTitle="Triggers" active="triggers" project="Agent House" showCmdk />

      <main className="main">
        <div className="page-head">
          <div>
            <h1>
              Triggers{" "}
              <span className="count">
                {jobs.length} cron {jobs.length === 1 ? "job" : "jobs"} · live scheduler
              </span>
            </h1>
            <p className="desc">
              Recurring jobs wake an agent every interval — backed by the live
              cron scheduler (Go duration: 30s · 5m · 1h · 24h). Email,
              webhook, and file-watch types arrive with the Phase-1 trigger
              engine.
            </p>
          </div>
          <div className="actions">
            <button
              className={`kill-toggle${armed ? " armed" : ""}`}
              onClick={() => setArmed((a) => !a)}
              title="Global pause lands with the Phase-1 trigger engine — visual only for now"
            >
              <span className="switch" />
              <span>Kill-switch</span>
            </button>
            <button className="btn btn--secondary" onClick={load}>
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round">
                <path d="M21 12a9 9 0 1 1-3-6.7M21 3v6h-6" />
              </svg>
              Refresh
            </button>
          </div>
        </div>

        <div className="filters">
          <div className="f-group">
            <span className="f-label">type</span>
            <button className="f-chip active">
              All <span className="ct">{jobs.length}</span>
            </button>
            <button className="f-chip">
              <span style={{ width: 11, height: 11, display: "grid", placeItems: "center" }}>
                {GLYPHS.cron}
              </span>
              cron <span className="ct">{counts.cron}</span>
            </button>
            {(["email", "webhook", "filewatch"] as TrigType[]).map((t) => (
              <button
                className="f-chip"
                key={t}
                disabled
                title="Phase-1 trigger engine"
                style={{ opacity: 0.5, cursor: "not-allowed" }}
              >
                <span style={{ width: 11, height: 11, display: "grid", placeItems: "center" }}>
                  {GLYPHS[t]}
                </span>
                {t} <span className="ct">soon</span>
              </button>
            ))}
          </div>
          <div className="f-search">
            <Search />
            <input placeholder="Search by task, agent…" />
          </div>
        </div>

        <div className="table-wrap">
          <table className="tr-table">
            <thead>
              <tr>
                <th style={{ width: 260 }}>trigger</th>
                <th>bind</th>
                <th>last fired</th>
                <th>schedule</th>
                <th style={{ width: 90 }}>status</th>
                <th style={{ width: 40 }} />
              </tr>
            </thead>
            <tbody>
              {jobs.length === 0 && (
                <tr>
                  <td colSpan={6} style={{ padding: "32px", textAlign: "center", color: "var(--c-muted)" }}>
                    {loaded
                      ? "No cron triggers yet — create one in the panel →"
                      : "Loading from scheduler…"}
                  </td>
                </tr>
              )}
              {jobs.map((j) => (
                <Row
                  key={j.id}
                  j={j}
                  expanded={expanded === j.id}
                  onToggle={() => setExpanded((e) => (e === j.id ? null : j.id))}
                  onDelete={() => remove(j.id)}
                  busy={busy}
                />
              ))}
            </tbody>
          </table>
        </div>
      </main>

      <aside className="panel">
        <div className="panel-head">
          <div className="step">new · cron trigger</div>
          <div className="row1">
            <h3>New trigger</h3>
            <button className="icon-btn"><Close /></button>
          </div>
        </div>

        <div className="panel-section">
          <h5>Type</h5>
          <p className="helper">
            Cron is live. Email / webhook / file-watch land with the Phase-1
            trigger engine.
          </p>
          <div className="type-grid">
            <div className="type-card selected">
              <div className="glyph" style={{ width: 26, height: 26 }}>{GLYPHS.cron}</div>
              <div className="name">Cron</div>
              <div className="desc">every 30s · 5m · 1h · 24h</div>
            </div>
            {(["email", "webhook", "filewatch"] as TrigType[]).map((t) => (
              <div
                key={t}
                className="type-card"
                style={{ opacity: 0.45, cursor: "not-allowed" }}
                title="Phase-1 trigger engine"
              >
                <div className="glyph" style={{ width: 26, height: 26 }}>{GLYPHS[t]}</div>
                <div className="name">{t[0].toUpperCase() + t.slice(1)}</div>
                <div className="desc">planned · Phase-1</div>
              </div>
            ))}
          </div>
        </div>

        <div className="panel-section">
          <h5>Definition</h5>
          <div className="field">
            <div className="label">Task <span className="req">*</span></div>
            <div className="input">
              <input
                value={fName}
                onChange={(e) => setFName(e.target.value)}
                placeholder="What the agent should do each run"
              />
            </div>
          </div>
          <div className="row2">
            <div className="field">
              <div className="label">Schedule <span className="req">*</span></div>
              <div className="input mono">
                <input
                  value={fSchedule}
                  onChange={(e) => setFSchedule(e.target.value)}
                  placeholder="24h · 1h · 5m · 30s"
                />
              </div>
            </div>
            <div className="field">
              <div className="label">Agent role</div>
              <div className="input mono">
                <input value={fRole} onChange={(e) => setFRole(e.target.value)} />
              </div>
            </div>
          </div>
          <div className="field">
            <div className="label">Project</div>
            <div className="input mono">
              <input value={fProject} onChange={(e) => setFProject(e.target.value)} />
            </div>
          </div>
        </div>

        <div className="panel-section">
          <h5>Safety</h5>
          <div className="safety">
            <div className="safety-row">
              <div className="l">
                <div className="t">Enabled at save</div>
                <div className="d">scheduler starts firing immediately</div>
              </div>
              <div className="switch on" />
            </div>
            <div className="safety-row">
              <div className="l">
                <div className="t">Dedupe / rate-limit</div>
                <div className="d">Phase-1 trigger engine</div>
              </div>
              <div className="switch" />
            </div>
          </div>
        </div>

        {flash && (
          <div className="test-flash">
            <span className="dot" />
            {flash}
          </div>
        )}

        <div className="panel-foot">
          <button className="btn btn--primary" style={{ flex: 1 }} disabled={busy} onClick={save}>
            {busy ? "Saving…" : "Save trigger"}
          </button>
        </div>
      </aside>
    </div>
  );
}

function Row({
  j,
  expanded,
  onToggle,
  onDelete,
  busy,
}: {
  j: CronJob;
  expanded: boolean;
  onToggle: () => void;
  onDelete: () => void;
  busy: boolean;
}) {
  return (
    <>
      <tr className={`tr-row${expanded ? " expanded" : ""}`} onClick={onToggle}>
        <td>
          <div className="tr-type">
            <div className="tr-glyph cron">{GLYPHS.cron}</div>
            <div className="info">
              <div className="name">{j.task.slice(0, 48) || "(untitled)"}</div>
              <div className="detail">id {j.id}</div>
            </div>
          </div>
        </td>
        <td>
          <span className="bind">
            <span className="arrow">→</span>agent:{j.agent_role}
          </span>
        </td>
        <td><span className="ts">{rel(j.last_run)}</span></td>
        <td><span className="ts">{j.schedule}</span></td>
        <td>
          <span className={`chip chip--${j.enabled ? "live" : "idle"}`}>
            <span className="dot" />
            {j.enabled ? "live" : "off"}
          </span>
        </td>
        <td>
          <button
            className="icon-btn"
            title="Delete trigger"
            disabled={busy}
            onClick={(e) => {
              e.stopPropagation();
              onDelete();
            }}
          >
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round">
              <path d="M4 7h16M9 7V5h6v2M7 7l1 13h8l1-13" />
            </svg>
          </button>
        </td>
      </tr>
      {expanded && (
        <tr className="tr-drawer-row">
          <td colSpan={6}>
            <div className="tr-drawer">
              <div>
                <h5>Configuration</h5>
                <div className="drawer-cfg">
                  {[
                    ["id", j.id],
                    ["schedule", j.schedule],
                    ["task", j.task],
                    ["agent", j.agent_role],
                    ["project", j.project_id],
                    ["enabled", String(j.enabled)],
                    ["last run", j.last_run || "never"],
                    ["next run", j.next_run || "—"],
                  ].map(([k, v]) => (
                    <span key={k} style={{ display: "contents" }}>
                      <span className="k">{k}</span>
                      <span className="v">{v}</span>
                    </span>
                  ))}
                </div>
                <div className="drawer-actions">
                  <button
                    className="btn btn--ghost btn--sm"
                    disabled={busy}
                    onClick={onDelete}
                  >
                    Delete
                  </button>
                </div>
              </div>
              <div>
                <h5>Notes</h5>
                <div className="drawer-cfg">
                  <span className="k">backend</span>
                  <span className="v">cron scheduler · live</span>
                  <span className="k">dedupe</span>
                  <span className="v">Phase-1</span>
                  <span className="k">audit</span>
                  <span className="v">Phase-1</span>
                </div>
              </div>
            </div>
          </td>
        </tr>
      )}
    </>
  );
}
