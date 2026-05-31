"use client";

import type { ScenarioResult, ScenarioStep, ScenarioSuggestion } from "@/lib/types";
import { humanize } from "@/lib/humanize";
import { Check, Bolt } from "./icons";

/**
 * Generic scenario result renderer. Anything the backend ships as a
 * `ScenarioResult` renders here — there are no per-scenario branches,
 * no hardcoded action labels, no name-keyed dispatch. Adding a new
 * scenario on the backend is enough; the UI picks it up.
 */
export function ScenarioCard({
  result,
  onAction,
}: {
  result: ScenarioResult;
  onAction?: (suggestion: ScenarioSuggestion) => void;
}) {
  return (
    <div className="scenario-card" data-ok={result.ok ? "1" : "0"}>
      <div className="sc-head">
        <span className="sc-tag">
          <span className="dot" />
          scenario · {humanize(result.scenario)}
        </span>
        <span className="sc-scope">scope · {result.scope}</span>
      </div>
      {result.summary && <div className="sc-summary">{result.summary}</div>}

      {result.steps.length > 0 && (
        <ol className="sc-steps">
          {result.steps.map((s, i) => (
            <li key={i} className={`sc-step${s.ok ? "" : " ko"}`}>
              <span className="sc-step-num">{i + 1}</span>
              <div className="sc-step-body">
                <div className="sc-step-head">
                  <span className="sc-agent">{humanize(s.agent)}</span>
                  <span className="sc-sep">·</span>
                  <span className="sc-action">{humanize(s.action)}</span>
                </div>
                {s.note && <div className="sc-step-note">{s.note}</div>}
                <StepOutput step={s} />
              </div>
            </li>
          ))}
        </ol>
      )}

      {result.suggestions.length > 0 && (
        <div className="sc-suggestions">
          <div className="sc-sug-label">Suggested next actions</div>
          {result.suggestions.map((sug, i) => (
            <SuggestionRow key={i} sug={sug} onAction={onAction} />
          ))}
        </div>
      )}
    </div>
  );
}

/**
 * StepOutput renders whatever shape the Step's output happens to be,
 * lightly formatted. No per-scenario knowledge: a couple of well-known
 * common keys ("matches", "slips", "check") get a slightly nicer summary,
 * everything else falls back to a generic key/value list.
 */
function StepOutput({ step }: { step: ScenarioStep }) {
  const out = step.output;
  if (!out) return null;

  // Generic, shape-driven enhancements — these all detect the shape, not
  // the scenario name. Any future scenario producing the same shape gets
  // the same treatment for free.
  if (
    "matches" in out &&
    Array.isArray((out as { matches: unknown }).matches) &&
    (out as { matches: unknown[] }).matches.length > 0
  ) {
    const matches = out.matches as Array<{
      applicant: { name: string; applied_for?: string; experience_years?: number };
      score: number;
      matched?: string[];
    }>;
    return (
      <div className="sc-mini">
        {matches.slice(0, 3).map((m, i) => (
          <div className="sc-mini-row" key={i}>
            <span className="rank">#{i + 1}</span>
            <span className="title">{m.applicant?.name ?? "unknown"}</span>
            <span className="meta">
              score {m.score}
              {m.applicant?.experience_years
                ? ` · ${m.applicant.experience_years}y`
                : ""}
            </span>
          </div>
        ))}
      </div>
    );
  }
  if (
    "slips" in out &&
    Array.isArray((out as { slips: unknown }).slips) &&
    (out as { slips: unknown[] }).slips.length > 0
  ) {
    const slips = out.slips as Array<{
      milestone: { title: string };
      days_late: number;
      severity: string;
    }>;
    return (
      <div className="sc-mini">
        {slips.slice(0, 3).map((s, i) => (
          <div className={`sc-mini-row sev-${s.severity}`} key={i}>
            <span className="rank">{s.severity[0].toUpperCase()}</span>
            <span className="title">{s.milestone.title}</span>
            <span className="meta">+{s.days_late}d · {s.severity}</span>
          </div>
        ))}
      </div>
    );
  }
  if ("check" in out) {
    const c = (out as { check: { recommendation: string; reasons?: string[]; trusted: boolean; impersonation_risk: boolean } }).check;
    return (
      <div className="sc-mini">
        <div className="sc-mini-row">
          <span className="rank">→</span>
          <span className="title">{c.recommendation}</span>
          <span className="meta">
            {c.impersonation_risk
              ? "impersonation"
              : c.trusted
              ? "trusted"
              : "unverified"}
          </span>
        </div>
        {c.reasons?.slice(0, 2).map((r, i) => (
          <div className="sc-mini-row dim" key={i}>
            <span className="rank">·</span>
            <span className="title">{r}</span>
          </div>
        ))}
      </div>
    );
  }
  return null;
}

function SuggestionRow({
  sug,
  onAction,
}: {
  sug: ScenarioSuggestion;
  onAction?: (s: ScenarioSuggestion) => void;
}) {
  return (
    <div className="sc-sug">
      <Bolt />
      <div className="sc-sug-body">
        <div className="sc-sug-title">{sug.title}</div>
        <div className="sc-sug-detail">{sug.detail}</div>
      </div>
      <button
        className="btn btn--secondary btn--sm"
        onClick={() => onAction?.(sug)}
        title={`action id: ${sug.action}`}
      >
        <Check />
        {humanize(sug.action)}
      </button>
    </div>
  );
}
