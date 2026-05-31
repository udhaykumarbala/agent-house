import type { AgentVM } from "@/lib/types";
import { Avatar } from "./Avatar";
import { Gear, Plus } from "./icons";

export function Sidebar({
  packName,
  packMeta,
  groups,
  selectedRole,
  onSelect,
  version,
}: {
  packName: string;
  packMeta: string;
  groups: { label: string; agents: AgentVM[] }[];
  selectedRole?: string;
  onSelect?: (role: string) => void;
  version: string;
}) {
  return (
    <aside className="side">
      <div className="side-head">
        <div className="glyph">{packName.slice(0, 2).toUpperCase()}</div>
        <div className="info">
          <div className="t">{packName} pack</div>
          <div className="d">{packMeta}</div>
        </div>
        <button className="icon-btn" title="Pack settings">
          <Gear />
        </button>
      </div>

      {groups.map((g) => (
        <div className="side-grp" key={g.label}>
          <div className="side-label">
            <span>{g.label}</span>
            <span className="count">{g.agents.length}</span>
          </div>
          {g.agents.map((a, i) => (
            <div
              key={`${a.role}-${i}`}
              className={`side-item${
                selectedRole && a.role === selectedRole ? " selected" : ""
              }`}
              onClick={() => onSelect?.(a.role)}
            >
              <Avatar initials={a.initials} state={a.state} size={32} />
              <div className="meta">
                <div className="name">{a.name}</div>
                <div className="role">{a.role}</div>
              </div>
              {a.badgeCount != null && (
                <span className={`tiny-chip ${a.badgeTone ?? ""}`}>
                  <span className="pip" />
                  {a.badgeCount}
                </span>
              )}
            </div>
          ))}
          {g.label === "Custom" && (
            <div className="side-add">
              <Plus />
              Add custom agent
            </div>
          )}
        </div>
      ))}

      <div className="side-foot">
        <span>
          <span className="dot" />
          registry · live
        </span>
        <span>{version}</span>
      </div>
    </aside>
  );
}
