import type { AgentState } from "@/lib/types";

export function Avatar({
  initials,
  state,
  size = 32,
}: {
  initials: string;
  state: AgentState;
  size?: 28 | 32 | 40;
}) {
  return (
    <div className={`av av-${size} ${state}`}>
      {initials}
      <div className="st" />
    </div>
  );
}
