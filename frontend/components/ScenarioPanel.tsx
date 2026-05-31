"use client";

import { useEffect, useMemo, useState } from "react";
import type { ScenarioMeta, ScenarioResult } from "@/lib/types";
import { fetchScenarios, runScenario } from "@/lib/api";
import { Bolt } from "./icons";

/**
 * Lists all scenarios the backend declares. Driven entirely by
 * /api/scenarios — no scenario name appears in this component anywhere.
 * Click a scenario → inline form pre-filled from its declared example →
 * Run → onResult dispatches back to the parent which renders a card in
 * the conversation.
 */
export function ScenarioPanel({
  scope,
  onScope,
  onResult,
}: {
  scope: string;
  onScope?: (s: string) => void;
  onResult: (r: ScenarioResult) => void;
}) {
  const [scenarios, setScenarios] = useState<ScenarioMeta[]>([]);
  const [openName, setOpenName] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);

  useEffect(() => {
    fetchScenarios().then(setScenarios);
  }, []);

  const opened = useMemo(
    () => scenarios.find((s) => s.name === openName),
    [scenarios, openName]
  );

  return (
    <div className="right-section">
      <h5>
        Scenarios{" "}
        <span className="ct">{scenarios.length}</span>
      </h5>
      <div className="sp-scope">
        <label>scope</label>
        <input
          value={scope}
          onChange={(e) => onScope?.(e.target.value)}
          placeholder="default"
        />
      </div>
      <div className="sp-list">
        {scenarios.length === 0 && (
          <div style={{ color: "var(--c-muted)", fontSize: 12 }}>
            No scenarios registered.
          </div>
        )}
        {scenarios.map((s) => (
          <div key={s.name} className="sp-item">
            <button
              className="sp-row"
              onClick={() => setOpenName(openName === s.name ? null : s.name)}
              aria-expanded={openName === s.name}
            >
              <Bolt />
              <div className="sp-meta">
                <div className="sp-name">{s.name}</div>
                <div className="sp-desc">{s.description}</div>
              </div>
            </button>
            {openName === s.name && opened && (
              <RunForm
                key={s.name}
                meta={opened}
                busy={busy}
                scope={scope}
                onRun={async (input) => {
                  setBusy(true);
                  const r = await runScenario(opened.name, input, scope);
                  setBusy(false);
                  if (r) {
                    onResult(r);
                    setOpenName(null);
                  }
                }}
                onCancel={() => setOpenName(null)}
              />
            )}
          </div>
        ))}
      </div>
    </div>
  );
}

function RunForm({
  meta,
  busy,
  scope,
  onRun,
  onCancel,
}: {
  meta: ScenarioMeta;
  busy: boolean;
  scope: string;
  onRun: (input: Record<string, unknown>) => void;
  onCancel: () => void;
}) {
  const [text, setText] = useState(() => JSON.stringify(meta.example, null, 2));
  const [err, setErr] = useState<string | null>(null);

  const handleRun = () => {
    try {
      const parsed = text.trim() ? JSON.parse(text) : {};
      setErr(null);
      onRun(parsed);
    } catch (e) {
      setErr((e as Error).message);
    }
  };

  return (
    <div className="sp-form">
      <div className="sp-form-meta">
        input · JSON · target scope = {scope || "default"}
      </div>
      <textarea
        value={text}
        onChange={(e) => setText(e.target.value)}
        rows={Math.min(8, Math.max(3, text.split("\n").length))}
        spellCheck={false}
      />
      {err && <div className="sp-err">{err}</div>}
      <div className="sp-form-actions">
        <button className="btn btn--ghost btn--sm" onClick={onCancel} disabled={busy}>
          Cancel
        </button>
        <button className="btn btn--primary btn--sm" onClick={handleRun} disabled={busy}>
          {busy ? "Running…" : "Run scenario"}
        </button>
      </div>
    </div>
  );
}
