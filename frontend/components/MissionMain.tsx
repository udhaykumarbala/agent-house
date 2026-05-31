import type { HitlVM, FeedVM, StatVM } from "@/lib/types";
import { Avatar } from "./Avatar";
import { Bolt, FileGlyph } from "./icons";

const TABS = ["Overview", "Triggers", "Activity", "Files", "Settings"];

export function MissionMain({
  dateLabel,
  stats,
  hitl,
  feed,
  selectedHitl,
  onSelectHitl,
  onReview,
}: {
  dateLabel: string;
  stats: StatVM[];
  hitl: HitlVM[];
  feed: FeedVM[];
  selectedHitl?: string;
  onSelectHitl: (id: string) => void;
  onReview: (id: string) => void;
}) {
  return (
    <main className="main">
      <div className="tabs">
        {TABS.map((t, i) => (
          <button key={t} className={`tab${i === 0 ? " active" : ""}`}>
            {t}
            {t === "Overview" && hitl.length > 0 && (
              <span className="tab-badge fire">{hitl.length}</span>
            )}
          </button>
        ))}
      </div>

      <div className="main-body">
        <div className="today">
          <h1>
            Today <span className="date">{dateLabel}</span>
          </h1>
          <div className="stats">
            {stats.map((s, i) => (
              <div className={`stat ${s.tone}`} key={i}>
                <div className="n">{s.n}</div>
                <div className="l">{s.l}</div>
              </div>
            ))}
          </div>
        </div>

        <div>
          <div className="s-head">
            <h2>
              Awaiting your approval{" "}
              <span className="tiny-chip wait">
                <span className="pip" />
                {hitl.length} pending
              </span>
            </h2>
            <a href="#" className="s-link">
              view all →
            </a>
          </div>
          <div className="hitl-list">
            {hitl.map((h) => (
              <div
                key={h.id}
                className={`hitl${selectedHitl === h.id ? " selected" : ""}`}
                onClick={() => onSelectHitl(h.id)}
              >
                <Avatar initials={h.initials} state={h.state} size={32} />
                <div className="body">
                  <div className="t">{h.title}</div>
                  <div className="d">{h.detail}</div>
                </div>
                <span className="tag">{h.tag}</span>
                <span className="timer">
                  <span className="dot" />
                  {h.timer}
                </span>
                <button
                  className="btn-mini"
                  onClick={(e) => {
                    e.stopPropagation();
                    onReview(h.id);
                  }}
                >
                  Review
                </button>
              </div>
            ))}
          </div>
        </div>

        <div>
          <div className="s-head">
            <h2>
              Activity <span className="meta">live</span>
            </h2>
            <a href="#" className="s-link">
              today · all events →
            </a>
          </div>
          <div className="feed">
            {feed.map((f) => (
              <div className={`feed-item ${f.tone}`} key={f.id}>
                <div className="feed-dot" />
                <div className="feed-head">
                  {f.triggered && (
                    <span className={`feed-tag ${f.triggered === "live" ? "live" : ""}`}>
                      <Bolt />
                      triggered
                    </span>
                  )}
                  {f.head.map((h, i) => (
                    <span key={i} className={i === 0 && !f.triggered ? "ts" : ""}>
                      {h}
                      {i < f.head.length - 1 ? " ·" : ""}
                    </span>
                  ))}
                  {!f.triggered && <span className="ts">{f.ts}</span>}
                  {f.triggered && (
                    <>
                      <span className="ts">{f.ts}</span>
                      <span>·</span>
                    </>
                  )}
                </div>
                <p
                  className="feed-body"
                  dangerouslySetInnerHTML={{ __html: f.body }}
                />
                {f.card && (
                  <div className="feed-card">
                    <div className="icon">
                      <FileGlyph />
                    </div>
                    <div>
                      <div className="t">{f.card.title}</div>
                      <div className="d">{f.card.detail}</div>
                    </div>
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>
      </div>
    </main>
  );
}
