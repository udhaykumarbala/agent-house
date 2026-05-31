import type { Briefing } from "@/lib/briefing";

/**
 * The proactive opener — what the Conductor says first, unprompted, on
 * every load. Built from live state by `composeBriefing()`.
 */
export function BriefingCard({ briefing }: { briefing: Briefing }) {
  return (
    <div className="briefing">
      <div className="briefing-head">
        <div className="briefing-greeting">{briefing.greeting}</div>
        <div className="briefing-headline">{briefing.headline}</div>
      </div>
      {briefing.emptyMessage ? (
        <div className="briefing-empty">{briefing.emptyMessage}</div>
      ) : (
        <div className="briefing-items">
          {briefing.items.map((it, i) => (
            <div className={`briefing-item ${it.tone}`} key={i}>
              <div className="dot" />
              <div className="body">
                <div className="t">{it.title}</div>
                <div className="d">{it.detail}</div>
              </div>
              {it.action && (
                <a className="btn btn--secondary btn--sm" href={it.action.href}>
                  {it.action.label}
                </a>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
