"use client";

import { useCallback, useEffect, useState } from "react";
import { Chrome } from "./Chrome";
import { ScenarioPanel } from "./ScenarioPanel";
import { ScenarioCard } from "./ScenarioCard";
import { MockMailComposer } from "./MockMailComposer";
import { API_BASE } from "@/lib/api";
import { ATLAS_APPLICANTS, ATLAS_VENDORS, ATLAS_MILESTONES } from "@/lib/seed";
import type { ScenarioResult, ScenarioSuggestion } from "@/lib/types";

interface Counts {
  applicants: number;
  vendors: number;
  milestones: number;
  slips: number;
}

export function LabApp() {
  const [scope, setScope] = useState("atlas-site");
  const [counts, setCounts] = useState<Counts | null>(null);
  const [results, setResults] = useState<ScenarioResult[]>([]);
  const [seeding, setSeeding] = useState(false);
  const [flash, setFlash] = useState<string | null>(null);

  const refresh = useCallback(async () => {
    const enc = encodeURIComponent(scope);
    const grab = async (path: string): Promise<number> => {
      try {
        const r = await fetch(`${API_BASE}/api/cap/${path}?scope=${enc}`);
        if (!r.ok) return 0;
        const d = await r.json();
        return Number(d.count ?? 0);
      } catch {
        return 0;
      }
    };
    const [a, v, m, s] = await Promise.all([
      grab("hr/applicants"),
      grab("procurement/vendors"),
      grab("schedule/milestones"),
      grab("schedule/slips"),
    ]);
    setCounts({ applicants: a, vendors: v, milestones: m, slips: s });
  }, [scope]);

  useEffect(() => {
    refresh();
  }, [refresh]);

  const flashFor = (msg: string) => {
    setFlash(msg);
    setTimeout(() => setFlash(null), 3500);
  };

  const seedAtlas = async () => {
    if (seeding) return;
    setSeeding(true);
    const enc = encodeURIComponent(scope);
    const post = async (path: string, body: unknown) =>
      fetch(`${API_BASE}/api/cap/${path}?scope=${enc}`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(body),
      });
    try {
      await Promise.all([
        ...ATLAS_APPLICANTS.map((a) => post("hr/applicants", a)),
        ...ATLAS_VENDORS.map((v) => post("procurement/vendors", v)),
        ...ATLAS_MILESTONES.map((m) => post("schedule/milestones", m)),
      ]);
      flashFor(`Seeded Atlas demo data into "${scope}".`);
      refresh();
    } catch (e) {
      flashFor(`Seed failed: ${(e as Error).message}`);
    } finally {
      setSeeding(false);
    }
  };

  const onResult = (r: ScenarioResult) => setResults((p) => [r, ...p]);
  const onSuggestion = (sug: ScenarioSuggestion) => {
    // In the Lab, suggestions just copy their action+payload to clipboard for
    // wiring. The Conductor will wire dispatch when actions get real backends.
    const payload = JSON.stringify({ action: sug.action, payload: sug.payload }, null, 2);
    navigator.clipboard?.writeText(payload).then(
      () => flashFor(`Copied action "${sug.action}" to clipboard.`),
      () => flashFor(`Suggested action: ${sug.action}`)
    );
  };

  return (
    <div className="app lay-lab" data-screen-label="Lab">
      <Chrome
        pageTitle="Lab"
        active="lab"
        project="Agent House"
        pack="dev"
        packGlyph="LB"
        packCount="console"
      />

      <main className="main lab-main">
        <div className="lab-header">
          <h1>Lab</h1>
          <p>
            Development &amp; demo console. Trigger scenarios, seed backing
            data, inspect capability state. Nothing here is wired into the
            production Conductor view.
          </p>
          {flash && (
            <div className="lab-flash">
              <span className="dot" />
              {flash}
            </div>
          )}
        </div>

        <div className="lab-grid">
          {/* Demo data — left column */}
          <section className="lab-card">
            <h3>Demo data</h3>
            <div className="lab-scope">
              <label>scope</label>
              <input
                value={scope}
                onChange={(e) => setScope(e.target.value)}
                placeholder="atlas-site"
              />
            </div>
            <table className="lab-table">
              <tbody>
                <tr>
                  <td>applicants</td>
                  <td>{counts?.applicants ?? "—"}</td>
                </tr>
                <tr>
                  <td>vendors</td>
                  <td>{counts?.vendors ?? "—"}</td>
                </tr>
                <tr>
                  <td>milestones</td>
                  <td>{counts?.milestones ?? "—"}</td>
                </tr>
                <tr>
                  <td>detected slips</td>
                  <td>{counts?.slips ?? "—"}</td>
                </tr>
              </tbody>
            </table>
            <div className="lab-actions">
              <button
                className="btn btn--primary btn--sm"
                disabled={seeding}
                onClick={seedAtlas}
              >
                {seeding ? "Seeding…" : "Seed Atlas demo into this scope"}
              </button>
              <button className="btn btn--secondary btn--sm" onClick={refresh}>
                Refresh
              </button>
            </div>
            <div className="lab-hint">
              Hits <code>/api/cap/*</code> endpoints. Re-running is idempotent
              by id. Reset by removing <code>data/{scope}/</code> on disk.
            </div>
          </section>

          {/* Scenarios — right column */}
          <section className="lab-card lab-scenarios">
            <h3>Scenarios</h3>
            <p className="lab-helper">
              Each scenario composes one or more capabilities and returns a
              structured result with proactive suggestions. Adding a new
              scenario on the backend automatically appears here — this UI
              never names a specific scenario.
            </p>
            <ScenarioPanel
              scope={scope}
              onScope={setScope}
              onResult={onResult}
            />
          </section>
        </div>

        {/* Mock mail composer */}
        <MockMailComposer onSent={refresh} />

        {/* Results stack */}
        {results.length > 0 && (
          <div className="lab-results">
            <div className="lab-results-head">
              <h3>Recent runs</h3>
              <button
                className="btn btn--ghost btn--sm"
                onClick={() => setResults([])}
              >
                Clear
              </button>
            </div>
            {results.map((r, i) => (
              <ScenarioCard
                key={i}
                result={r}
                onAction={onSuggestion}
              />
            ))}
          </div>
        )}
      </main>
    </div>
  );
}
