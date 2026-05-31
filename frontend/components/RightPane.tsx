import type { HitlVM } from "@/lib/types";
import { Avatar } from "./Avatar";
import { Close, Check } from "./icons";

export function RightPane({
  active,
  rest,
  busy,
  onClose,
  onSelect,
  onDecision,
}: {
  active?: HitlVM;
  rest: HitlVM[];
  busy: boolean;
  onClose: () => void;
  onSelect: (id: string) => void;
  onDecision: (id: string, d: "approved" | "rejected") => void;
}) {
  return (
    <aside className="right" id="rightPane">
      <div className="right-head">
        <h3>
          <span className="cp-tag">
            <span className="dot" />
            checkpoint
          </span>
          {active ? `task#${active.id.slice(0, 8)}` : "none"}
        </h3>
        <button className="icon-btn" title="Close" onClick={onClose}>
          <Close />
        </button>
      </div>

      {active ? (
        <div className="cp">
          <div className="cp-who">
            <Avatar initials={active.initials} state={active.state} size={40} />
            <div>
              <div className="ag-name">{active.tag.replace(/_/g, " ")}</div>
              <div className="ag-role">{active.roleTag}</div>
            </div>
          </div>
          <h4 className="cp-title">{active.title}</h4>
          <p
            className="cp-desc"
            dangerouslySetInnerHTML={{ __html: active.desc }}
          />
          <div className="cp-payload">
            {active.payload.map(([k, v, t], i) => (
              <div key={i}>
                <span className="k">{k.padEnd(10, " ")}</span>{" "}
                <span className={t === "num" ? "v-num" : t === "str" ? "v-str" : ""}>
                  {v}
                </span>
              </div>
            ))}
          </div>
          <div className="cp-trace">
            <span className="pip" />
            <span>{active.trace}</span>
          </div>
          <div className="cp-countdown">
            <span>expires in</span>
            <span>
              <strong>{active.timer}</strong>
            </span>
          </div>
          <div className="cp-actions">
            <button
              className="btn btn--primary btn--block"
              disabled={busy || !active.live}
              onClick={() => onDecision(active.id, "approved")}
            >
              <Check />
              {active.primary}
            </button>
            <button className="btn btn--secondary btn--block" disabled={busy}>
              Ask agent for changes…
            </button>
            <button
              className="btn btn--ghost btn--block"
              disabled={busy || !active.live}
              onClick={() => onDecision(active.id, "rejected")}
            >
              Reject
            </button>
          </div>
        </div>
      ) : (
        <div className="right-empty">
          <div>
            <div className="glyph">
              <Check />
            </div>
            <div className="t">Nothing awaiting you</div>
            <div className="d">
              When an agent hits a gated action it appears here for your
              decision.
            </div>
          </div>
        </div>
      )}

      {rest.length > 0 && (
        <div className="right-rest">
          <h5>Also awaiting · {rest.length}</h5>
          {rest.map((r) => (
            <div className="rest-row" key={r.id} onClick={() => onSelect(r.id)}>
              <Avatar initials={r.initials} state={r.state} size={28} />
              <div className="t">{r.title}</div>
              <span className="timer">{r.timer}</span>
            </div>
          ))}
        </div>
      )}
    </aside>
  );
}
