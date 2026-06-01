"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import {
  fetchAgents,
  fetchMessages,
  fetchCheckpoints,
  resolveCheckpoint,
  connectWS,
  DEFAULT_PROJECT,
} from "@/lib/api";
import {
  mapAgents,
  mapCheckpoints,
  mapMessages,
  statsFrom,
} from "@/lib/mapper";
import {
  SAMPLE_AGENT_GROUPS,
  SAMPLE_HITL,
  SAMPLE_FEED,
  SAMPLE_STATS,
} from "@/lib/sample";
import type {
  ApiAgent,
  ApiMessage,
  ApiCheckpoint,
  HitlVM,
} from "@/lib/types";
import { TopBar } from "./TopBar";
import { Sidebar } from "./Sidebar";
import { MissionMain } from "./MissionMain";
import { RightPane } from "./RightPane";

export function MissionApp() {
  const [agents, setAgents] = useState<ApiAgent[]>([]);
  const [messages, setMessages] = useState<ApiMessage[]>([]);
  const [checkpoints, setCheckpoints] = useState<ApiCheckpoint[]>([]);
  const [connected, setConnected] = useState(false);
  const [selected, setSelected] = useState<string | undefined>();
  const [busy, setBusy] = useState(false);
  const [now, setNow] = useState(() => new Date());
  // Guard time-dependent rendering so the server prerender and the client's
  // first (hydration) render are identical. The app is a static export, so the
  // server HTML freezes the build-time clock; rendering live time on first
  // paint causes React hydration error #418. Render the clock only after mount.
  const [mounted, setMounted] = useState(false);
  useEffect(() => setMounted(true), []);

  // initial load
  useEffect(() => {
    let on = true;
    (async () => {
      const [a, m, c] = await Promise.all([
        fetchAgents(),
        fetchMessages(),
        fetchCheckpoints(),
      ]);
      if (!on) return;
      setAgents(a);
      setMessages(m);
      setCheckpoints(c);
    })();
    const clock = setInterval(() => setNow(new Date()), 30_000);
    return () => {
      on = false;
      clearInterval(clock);
    };
  }, []);

  // live stream
  useEffect(() => {
    const refetch = () => {
      fetchMessages().then(setMessages);
      fetchCheckpoints().then(setCheckpoints);
      fetchAgents().then(setAgents);
    };
    return connectWS((e) => {
      if (e.type === "message" && "message" in e) {
        setMessages((prev) => [e.message as ApiMessage, ...prev].slice(0, 200));
      } else if (e.type === "checkpoint") {
        refetch();
      } else {
        refetch();
      }
    }, setConnected);
  }, []);

  // ── derive view-models; fall back to the reference sample per-section ──
  const liveHitl = useMemo(() => mapCheckpoints(checkpoints), [checkpoints]);
  const hitl: HitlVM[] = liveHitl.length > 0 ? liveHitl : SAMPLE_HITL;

  const pendingRoles = useMemo(
    () => new Set(liveHitl.map((h) => h.tag.toLowerCase())),
    [liveHitl]
  );
  const liveGroups = useMemo(
    () => mapAgents(agents, pendingRoles),
    [agents, pendingRoles]
  );
  const groups =
    liveGroups.reduce((n, g) => n + g.agents.length, 0) > 0
      ? liveGroups
      : SAMPLE_AGENT_GROUPS;

  const liveFeed = useMemo(() => mapMessages(messages), [messages]);
  const feed = liveFeed.length > 0 ? liveFeed : SAMPLE_FEED;

  const stats =
    liveHitl.length > 0 || liveFeed.length > 0
      ? statsFrom(agents, hitl, feed)
      : SAMPLE_STATS;

  const activeHitl = hitl.find((h) => h.id === selected) ?? hitl[0];
  const rest = hitl.filter((h) => h.id !== activeHitl?.id);

  const onDecision = useCallback(
    async (id: string, d: "approved" | "rejected") => {
      setBusy(true);
      const ok = await resolveCheckpoint(id, d);
      setBusy(false);
      if (ok) {
        setCheckpoints((prev) => prev.filter((c) => c.id !== id));
        setSelected(undefined);
      }
    },
    []
  );

  const dateLabel = mounted
    ? now
        .toLocaleString("en-US", {
          weekday: "short",
          day: "2-digit",
          month: "short",
          hour: "2-digit",
          minute: "2-digit",
          hour12: false,
          timeZone: "UTC",
        })
        .toUpperCase()
        .replace(",", " ·")
        .concat(" UTC")
    : "";

  return (
    <div className="app" data-screen-label="Mission v2 · workspace">
      <TopBar
        project="Atlas Construction"
        site={`Site 02 · Mission`}
        pack="EPC"
        packCount={`${groups.reduce((n, g) => n + g.agents.length, 0)}`}
        pendingCount={hitl.length}
        user="Marcus L."
      />
      <Sidebar
        packName="EPC"
        packMeta={`${groups.reduce(
          (n, g) => n + g.agents.length,
          0
        )} agents · live triggers · v1.0`}
        groups={groups}
        selectedRole={activeHitl?.tag.toLowerCase()}
        version={`${DEFAULT_PROJECT} · v0.4`}
      />
      <MissionMain
        dateLabel={dateLabel}
        stats={stats}
        hitl={hitl}
        feed={feed}
        selectedHitl={activeHitl?.id}
        onSelectHitl={setSelected}
        onReview={setSelected}
      />
      <RightPane
        active={activeHitl}
        rest={rest}
        busy={busy}
        onClose={() => setSelected("__none__")}
        onSelect={setSelected}
        onDecision={onDecision}
      />
      {!connected && (
        <div className="conn-banner">
          <span className="dot" />
          offline — showing reference data · reconnecting…
        </div>
      )}
    </div>
  );
}
