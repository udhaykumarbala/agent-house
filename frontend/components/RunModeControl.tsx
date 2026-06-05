"use client";

import { useEffect, useState } from "react";
import { getRunMode, setRunMode } from "@/lib/api";

const MODES: { id: string; label: string; blurb: string }[] = [
  { id: "manual", label: "Manual", blurb: "Gate every stage — you decide everything." },
  {
    id: "semi_auto",
    label: "Semi-auto",
    blurb: "Show choices at PRD, design, phase gates & final — wait for you.",
  },
  {
    id: "full_auto",
    label: "Full-auto",
    blurb: "Show the choice with a countdown; the CEO decides (with reasoning) if you don't answer in time.",
  },
  { id: "blitz", label: "Blitz", blurb: "No gates — auto-approve everything, fastest build." },
];

/** Sets the GLOBAL default build-autonomy mode that new builds inherit
 * (persisted via the workflow-settings API on the `_global` pseudo-project). */
export function RunModeControl() {
  const [open, setOpen] = useState(false);
  const [mode, setMode] = useState("semi_auto");
  const [mins, setMins] = useState(5);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    getRunMode().then((s) => {
      if (!s) return;
      setMode(s.run_mode || "semi_auto");
      if (s.decision_timeout_minutes) setMins(s.decision_timeout_minutes);
    });
  }, []);

  const apply = async (m: string, t: number) => {
    setMode(m);
    setMins(t);
    setSaving(true);
    await setRunMode(m, t);
    setSaving(false);
  };

  const cur = MODES.find((x) => x.id === mode) || MODES[1];
  const label = mode === "full_auto" ? `${cur.label} · ${mins}m` : cur.label;

  return (
    <div className="runmode">
      <button
        className="btn btn--secondary btn--sm"
        onClick={() => setOpen((o) => !o)}
        title="Build autonomy mode for new builds"
      >
        <span className={`runmode-dot rm-${mode}`} /> {label}
      </button>
      {open && (
        <div className="runmode-menu" onMouseLeave={() => setOpen(false)}>
          <div className="runmode-head">Build autonomy</div>
          {MODES.map((m) => (
            <button
              key={m.id}
              className={`runmode-opt${m.id === mode ? " active" : ""}`}
              onClick={() => apply(m.id, mins)}
            >
              <span className={`runmode-dot rm-${m.id}`} />
              <span className="runmode-meta">
                <span className="runmode-label">{m.label}</span>
                <span className="runmode-blurb">{m.blurb}</span>
              </span>
            </button>
          ))}
          {mode === "full_auto" && (
            <div className="runmode-timeout">
              <label>Decision window</label>
              <div className="runmode-timeout-row">
                <input
                  type="number"
                  min={1}
                  max={120}
                  value={mins}
                  onChange={(e) => setMins(Math.max(1, Number(e.target.value) || 1))}
                  onBlur={() => apply("full_auto", mins)}
                />
                <span>min — then the CEO decides on its own.</span>
              </div>
            </div>
          )}
          {saving && <div className="runmode-saving">Saving…</div>}
        </div>
      )}
    </div>
  );
}
