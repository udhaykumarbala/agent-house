"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { Chrome } from "./Chrome";
import { Search, Plus, Close } from "./icons";
import { fetchAgents, fetchAgentModes, setAgentMode } from "@/lib/api";
import type { ApiAgent } from "@/lib/types";

type Pack = "epc" | "it" | "custom" | "platform";
interface BAgent {
  initials: string;
  name: string;
  role: string;
  pack: Pack;
  mode: string;
  perms: [boolean, boolean, boolean];
  paths: string[];
  denied: string[];
  state: string;
  active: boolean;
}

const PACK_OF = (role: string): Pack =>
  role === "ceo"
    ? "platform"
    : ["procurement", "project_manager", "site_engineer", "qa_inspector", "hse", "hr"].includes(role)
    ? "epc"
    : ["pm", "ux", "ui", "security", "architect", "senior_dev", "junior_dev"].includes(role)
    ? "it"
    : "custom";

function fromApi(a: ApiAgent, mode: string): BAgent {
  const parts = (a.name || a.role).split(/[\s_]+/).filter(Boolean);
  const initials = (
    parts.length >= 2 ? parts[0][0] + parts[1][0] : (a.name || a.role).slice(0, 2)
  ).toUpperCase();
  const p = a.permissions ?? {};
  return {
    initials,
    name: a.name || a.role,
    role: a.role,
    pack: PACK_OF(a.role),
    mode,
    perms: [!!p.CanCreateFiles, !!p.CanModifyFiles, !!p.CanExecuteCode],
    paths: p.AllowedPaths ?? [],
    denied: p.DeniedPaths ?? [],
    state: a.active ? "working" : "live",
    active: a.active,
  };
}

const TABS = ["Identity", "Prompt", "Permissions", "Hierarchy"];

export function AgentBuilderApp() {
  const [raw, setRaw] = useState<ApiAgent[]>([]);
  const [modes, setModes] = useState<Record<string, string>>({});
  const [defaults, setDefaults] = useState<Record<string, string>>({});
  const [selected, setSelected] = useState(0);
  const [tab, setTab] = useState(0);
  const [packFilter, setPackFilter] = useState<"all" | Pack>("all");
  const [loaded, setLoaded] = useState(false);
  const [savingMode, setSavingMode] = useState(false);

  const load = useCallback(async () => {
    const [a, m] = await Promise.all([fetchAgents(), fetchAgentModes()]);
    setRaw(a);
    setModes(m.modes);
    setDefaults(m.defaults);
    setLoaded(true);
  }, []);
  useEffect(() => {
    load();
  }, [load]);

  const agents = useMemo(
    () =>
      raw.map((a) =>
        fromApi(a, modes[a.role] ?? defaults[a.role] ?? "session")
      ),
    [raw, modes, defaults]
  );

  const shown = useMemo(
    () => agents.filter((a) => packFilter === "all" || a.pack === packFilter),
    [agents, packFilter]
  );
  const sel = shown[selected] ?? shown[0];

  const counts = useMemo(() => {
    const c = { epc: 0, it: 0, custom: 0, platform: 0 } as Record<Pack, number>;
    agents.forEach((a) => (c[a.pack] += 1));
    return c;
  }, [agents]);

  async function changeMode(role: string, mode: "oneshot" | "session") {
    if (!role) return;
    setSavingMode(true);
    const ok = await setAgentMode(role, mode);
    setSavingMode(false);
    if (ok) setModes((m) => ({ ...m, [role]: mode }));
  }

  return (
    <div className="app lay-builder" data-screen-label="Agent Builder">
      <Chrome
        pageTitle="Agents"
        active="agents"
        project="Agent House"
        pack={loaded ? "registry" : "EPC"}
        packGlyph={loaded ? "AG" : "EP"}
        packCount={loaded ? `${agents.length} agents` : "6 + 1"}
        showCmdk
      />

      <main className="main">
        <div className="page-head">
          <div>
            <h1>
              Agents{" "}
              <span className="count">
                {agents.length} in registry · {counts.epc} epc · {counts.it} it ·{" "}
                {counts.custom} custom · {counts.platform} platform
              </span>
            </h1>
            <p className="desc">
              Live from the agent registry. Mode (oneshot/session) saves
              instantly. Prompt, permissions, and new-agent creation land with
              the Phase-1 registry write API.
            </p>
          </div>
          <div className="actions">
            <button className="btn btn--secondary" onClick={load}>
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round">
                <path d="M21 12a9 9 0 1 1-3-6.7M21 3v6h-6" />
              </svg>
              Refresh
            </button>
            <button
              className="btn btn--primary"
              disabled
              title="New-agent creation lands with the Phase-1 registry write API"
              style={{ opacity: 0.55, cursor: "not-allowed" }}
            >
              <Plus />
              New agent
            </button>
          </div>
        </div>

        <div className="filters">
          <div className="f-group">
            <span className="f-label">pack</span>
            <button
              className={`f-chip${packFilter === "all" ? " active" : ""}`}
              onClick={() => {
                setPackFilter("all");
                setSelected(0);
              }}
            >
              All <span className="ct">{agents.length}</span>
            </button>
            {(
              [
                ["epc", "EPC", "oklch(56% 0.16 50)"],
                ["it", "IT", "oklch(54% 0.14 240)"],
                ["custom", "Custom", "oklch(40% 0.02 240)"],
                ["platform", "Platform", "var(--c-ink)"],
              ] as [Pack, string, string][]
            ).map(([k, label, color]) => (
              <button
                key={k}
                className={`f-chip${packFilter === k ? " active" : ""}`}
                onClick={() => {
                  setPackFilter(k);
                  setSelected(0);
                }}
              >
                <span className="pip" style={{ background: color }} />
                {label} <span className="ct">{counts[k]}</span>
              </button>
            ))}
          </div>
          <div className="f-search">
            <Search />
            <input placeholder="Search agents by role…" />
          </div>
        </div>

        <div className="grid-wrap">
          <div className="ag-grid">
            {!loaded && (
              <div style={{ color: "var(--c-muted)", padding: "24px" }}>
                Loading registry…
              </div>
            )}
            {shown.map((a, i) => (
              <div
                key={a.role}
                className={`agent-card${i === selected ? " selected" : ""}`}
                onClick={() => setSelected(i)}
              >
                <div className="top">
                  <div className={`av av-40 ${a.state}`}>
                    {a.initials}
                    <div className="st" />
                  </div>
                  <div className="meta">
                    <div className="name">{a.name}</div>
                    <div className="role">role:{a.role}</div>
                  </div>
                  <span className={`pack-badge ${a.pack}`}>
                    <span className="pip" />
                    {a.pack}
                  </span>
                </div>
                <div className="perm-row">
                  <span className="k">mode</span>
                  <span className="v">{a.mode}</span>
                  <span style={{ marginLeft: "auto" }} className="perm-bullets">
                    {(["C", "M", "X"] as const).map((c, j) => (
                      <span key={c} className={`pbull ${a.perms[j] ? "on" : "off"}`}>
                        {c}
                      </span>
                    ))}
                  </span>
                </div>
                <div className="foot">
                  <span className="trigs">
                    {a.paths.length ? `${a.paths.length} allowed paths` : "registry"}
                  </span>
                  <span className={`chip chip--${a.active ? "info" : "live"}`}>
                    <span className="dot" />
                    {a.active ? "working" : "live"}
                  </span>
                </div>
              </div>
            ))}
          </div>
        </div>
      </main>

      <aside className="panel">
        <div className="panel-head">
          <div className="step">edit · mode saves live</div>
          <div className="row1">
            <div className="who">
              <div className={`av av-40 ${sel?.state ?? "idle"}`}>
                {sel?.initials}
                <div className="st" />
              </div>
              <div>
                <div className="name">{sel?.name}</div>
                <div className="role">
                  role:{sel?.role} · {sel?.pack}
                </div>
              </div>
            </div>
            <button className="icon-btn"><Close /></button>
          </div>
        </div>

        <div className="ptabs">
          {TABS.map((t, i) => (
            <button
              key={t}
              className={`ptab${i === tab ? " active" : ""}`}
              onClick={() => setTab(i)}
            >
              {t}
            </button>
          ))}
        </div>

        {tab === 0 && (
          <div className="panel-section">
            <h5>Identity</h5>
            <div className="row2">
              <div className="field">
                <div className="label">Role id</div>
                <div className="input mono">
                  <span className="prefix">role:</span>
                  <input value={sel?.role ?? ""} readOnly />
                </div>
                <div className="help">from registry · read-only until Phase-1</div>
              </div>
              <div className="field">
                <div className="label">Display name</div>
                <div className="input">
                  <input value={sel?.name ?? ""} readOnly />
                </div>
              </div>
            </div>
            <div className="row2">
              <div className="field">
                <div className="label">Pack</div>
                <div className="input">
                  <input value={sel?.pack ?? ""} readOnly />
                </div>
              </div>
              <div className="field">
                <div className="label">
                  Mode {savingMode && <span className="help">saving…</span>}
                </div>
                <div className="seg" style={{ height: 32 }}>
                  {(["oneshot", "session"] as const).map((m) => (
                    <button
                      key={m}
                      className={sel?.mode === m ? "active" : ""}
                      disabled={savingMode}
                      onClick={() => sel && changeMode(sel.role, m)}
                    >
                      {m}
                    </button>
                  ))}
                </div>
                <div className="help">
                  {defaults[sel?.role ?? ""]
                    ? `default: ${defaults[sel!.role]}`
                    : "writes POST /api/agents/modes"}
                </div>
              </div>
            </div>
          </div>
        )}

        {tab === 1 && (
          <div className="panel-section">
            <h5>System prompt</h5>
            <p className="helper">
              Lives at <code>agents/{sel?.role}/agent.md</code>. Editing +
              hot-reload lands with the Phase-1 registry write API — read-only
              here.
            </p>
            <textarea
              className="input"
              readOnly
              value={`# ${sel?.name}\n\nPrompt content is served from the registry. The Phase-1 E1/E2 epics add the write path so this becomes editable with hot-reload.`}
            />
          </div>
        )}

        {tab === 2 && (
          <div className="panel-section">
            <h5>Permissions · live from registry</h5>
            <div className="perm-grid">
              {[
                ["C", "Create files", sel?.perms[0]],
                ["M", "Modify files", sel?.perms[1]],
                ["X", "Execute code", sel?.perms[2]],
              ].map(([b, t, on]) => (
                <div className="perm-toggle" key={b as string}>
                  <div className="l">
                    <span className={`pbull ${on ? "on" : "off"}`}>{b}</span>
                    <div className="info">
                      <div className="t">{t}</div>
                      <div className="d">{on ? "granted" : "denied"} · registry</div>
                    </div>
                  </div>
                  <div className={`switch${on ? " on" : ""}`} />
                </div>
              ))}
            </div>
            <div className="field" style={{ marginTop: 6 }}>
              <div className="label">Allowed paths</div>
              <div className="chip-input">
                {(sel?.paths.length ? sel.paths : ["—"]).map((p) => (
                  <span className="chip-tag" key={p}>
                    {p}
                  </span>
                ))}
              </div>
            </div>
            <div className="field">
              <div className="label">Denied paths</div>
              <div className="chip-input">
                {(sel?.denied.length ? sel.denied : ["—"]).map((p) => (
                  <span className="chip-tag" key={p}>
                    {p}
                  </span>
                ))}
              </div>
            </div>
          </div>
        )}

        {tab === 3 && (
          <div className="panel-section">
            <h5>Hierarchy</h5>
            <p className="helper">
              Delegation / review graph is registry-defined. Editing lands with
              the Phase-1 pack manifest work.
            </p>
            <div className="field">
              <div className="label">Pack</div>
              <div className="chip-input" style={{ minHeight: 32 }}>
                <span className="chip-tag">{sel?.pack}</span>
              </div>
            </div>
          </div>
        )}

        <div className="panel-foot">
          <div className="left">
            <span className="dot" />
            mode saves live · prompt/permissions = Phase-1
          </div>
          <button
            className="btn btn--primary"
            disabled
            title="Full save (prompt, permissions, hierarchy) lands with the Phase-1 registry write API"
            style={{ opacity: 0.55, cursor: "not-allowed" }}
          >
            Save
          </button>
        </div>
      </aside>
    </div>
  );
}
