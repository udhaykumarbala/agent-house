"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import {
  fetchAgents,
  fetchMessages,
  fetchAllCheckpoints,
  resolveCheckpoint,
  connectWS,
  DEFAULT_PROJECT,
} from "@/lib/api";
import { getCompany, companyConfig, type Company } from "@/lib/company";
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
  const [company, setCompanyState] = useState<Company>("software");
  const cfg = companyConfig(company);
  useEffect(() => {
    setMounted(true);
    setCompanyState(getCompany());
  }, []);

  // initial load
  useEffect(() => {
    let on = true;
    (async () => {
      const [a, m, c] = await Promise.all([
        fetchAgents(),
        fetchMessages(),
        fetchAllCheckpoints(),
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
      fetchAllCheckpoints().then(setCheckpoints);
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

  // Fast (1s) tick only while a full-auto gate is counting down, so its
  // "auto in M:SS" stays live without churning the page otherwise.
  const [tickMs, setTickMs] = useState(0);
  const hasCountdown = useMemo(
    () =>
      checkpoints.some(
        (c) => c.status === "pending" && c.run_mode === "full_auto"
      ),
    [checkpoints]
  );
  useEffect(() => {
    if (!hasCountdown) return;
    setTickMs(Date.now());
    const i = setInterval(() => setTickMs(Date.now()), 1000);
    return () => clearInterval(i);
  }, [hasCountdown]);

  // ── derive view-models; fall back to the reference sample per-section ──
  const liveHitl = useMemo(
    () => mapCheckpoints(checkpoints, tickMs),
    [checkpoints, tickMs]
  );
  const hitl: HitlVM[] = liveHitl.length > 0 ? liveHitl : SAMPLE_HITL;

  // Only the active company's roster, with proper display names.
  const roster = useMemo(
    () =>
      agents
        .filter((a) => cfg.roles.includes(a.role))
        .map((a) => ({ ...a, name: a.name || cfg.roleNames[a.role] || a.role })),
    [agents, cfg]
  );
  const pendingRoles = useMemo(
    () => new Set(liveHitl.map((h) => h.tag.toLowerCase())),
    [liveHitl]
  );
  const liveGroups = useMemo(
    () => mapAgents(roster, pendingRoles),
    [roster, pendingRoles]
  );
  const groups =
    liveGroups.reduce((n, g) => n + g.agents.length, 0) > 0
      ? liveGroups
      : SAMPLE_AGENT_GROUPS;

  // The message store is shared across companies; show only this company's
  // agents (+ system/conductor) so the EPC roster doesn't bleed into the
  // software feed and vice-versa.
  const companyMsgs = useMemo(() => {
    const rs = new Set<string>(cfg.roles);
    return messages.filter(
      (m) => !m.from || m.from === "system" || m.from === "conductor" || rs.has(m.from)
    );
  }, [messages, cfg]);
  const liveFeed = useMemo(() => mapMessages(companyMsgs), [companyMsgs]);
  const feed = liveFeed.length > 0 ? liveFeed : SAMPLE_FEED;

  const stats =
    liveHitl.length > 0 || liveFeed.length > 0
      ? statsFrom(roster, hitl, feed)
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
        project={cfg.label}
        site={`Mission`}
        pack={cfg.short}
        packCount={`${groups.reduce((n, g) => n + g.agents.length, 0)}`}
        pendingCount={hitl.length}
        user="Marcus L."
      />
      <Sidebar
        packName={cfg.short}
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
