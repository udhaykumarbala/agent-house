"use client";

import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { Chrome } from "./Chrome";
import { Plus, Check, Bolt } from "./icons";
import {
  sendChat,
  fetchChatHistory,
  fetchAgents,
  fetchMessages,
  fetchCheckpoints,
  connectWS,
} from "@/lib/api";
import { mapWorkforce, mapFires, mapCheckpoints } from "@/lib/mapper";
import { renderMarkdown, cleanReply } from "@/lib/md";
import { getCompany, companyConfig, type Company } from "@/lib/company";
import { RunModeControl } from "./RunModeControl";
import { humanize } from "@/lib/humanize";
import { composeBriefing, type Briefing } from "@/lib/briefing";
import { BriefingCard } from "./BriefingCard";
import { ActionCard } from "./ActionCard";
import { ConversationPalette } from "./ConversationPalette";
import { createConversation } from "@/lib/api";
import type { ApiAgent, ApiMessage, ApiCheckpoint } from "@/lib/types";
import type { ChatReply } from "@/lib/api";

type Msg = {
  role: "user" | "conductor";
  html: string;
  meta?: string;
  reply?: ChatReply;
};

const DEFAULT_SUGGEST = [
  "What needs my attention today?",
  "Pause all triggers",
  "Add a custom agent…",
  "Email me a daily summary",
  "Show the dependency graph",
];

// Mock SAMPLE_WF / SAMPLE_FIRES were removed in the production-shape pass —
// the Conductor now renders honest empty states when there's no live data,
// and demo data is seeded from /lab.

function esc(s: string) {
  return s.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
}

export function ConductorApp() {
  const [msgs, setMsgs] = useState<Msg[]>([]);
  const [hasHistory, setHasHistory] = useState(false);
  const [thinking, setThinking] = useState(false);
  // Live backend activity shown WHILE thinking — scenario steps + tier-1 routing
  // narration stream in here so the wait reads as "agents working", not dead dots.
  const [liveSteps, setLiveSteps] = useState<string[]>([]);
  const [text, setText] = useState("");
  const [agents, setAgents] = useState<ApiAgent[]>([]);
  const [messages, setMessages] = useState<ApiMessage[]>([]);
  const [checkpoints, setCheckpoints] = useState<ApiCheckpoint[]>([]);
  const [suggest, setSuggest] = useState<string[]>(DEFAULT_SUGGEST);
  const [connected, setConnected] = useState(false);
  const [loaded, setLoaded] = useState(false);
  const [activeProject, setActiveProject] = useState<string | undefined>();
  // The conversation ID is the unit of persistence. The URL is the source
  // of truth (so reloading lands the user back in the same conversation),
  // backed by localStorage if no ?conv= is present.
  // Every page load starts a FRESH conversation — we never restore old history
  // into the thread (past conversations stay reachable via History/⌘K). Start
  // empty, then mint the id after mount to avoid a static-export hydration
  // mismatch (see MissionApp #418).
  const [convId, setConvId] = useState<string>("");
  // Active company (software studio vs EPC). Resolved after mount to avoid a
  // static-export hydration mismatch; drives the roster, scope + chrome labels.
  const [company, setCompanyState] = useState<Company>("software");
  const cfg = companyConfig(company);
  useEffect(() => {
    if (typeof window === "undefined") return;
    const fromQuery = new URL(window.location.href).searchParams.get("conv");
    setConvId(fromQuery || `conv_${Date.now()}`);
    setCompanyState(getCompany());
  }, []);
  const [paletteOpen, setPaletteOpen] = useState(false);
  const convRef = useRef<HTMLDivElement>(null);
  const taRef = useRef<HTMLTextAreaElement>(null);
  const sendingRef = useRef(false); // synchronous send-lock (state `thinking` races on rapid Enter+click)

  // Persist convId across reloads — URL + localStorage.
  useEffect(() => {
    if (typeof window === "undefined") return;
    window.localStorage.setItem("ah:conv", convId);
    const url = new URL(window.location.href);
    if (url.searchParams.get("conv") !== convId) {
      url.searchParams.set("conv", convId);
      window.history.replaceState({}, "", url.toString());
    }
  }, [convId]);

  // ⌘K / Ctrl-K opens the conversation palette.
  useEffect(() => {
    if (typeof window === "undefined") return;
    const onKey = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key === "k") {
        e.preventDefault();
        setPaletteOpen((o) => !o);
      }
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, []);

  const refreshSidebar = useCallback(() => {
    fetchAgents().then(setAgents);
    fetchMessages().then(setMessages);
    fetchCheckpoints().then(setCheckpoints);
  }, []);

  // initial load
  // Refetch history + reset the visible thread whenever convId changes —
  // covers initial load, "New conversation" clicks, and palette switches.
  useEffect(() => {
    if (!convId) return; // wait until the fresh conversation id is minted on mount
    let on = true;
    setLoaded(false);
    setMsgs([]);
    setHasHistory(false);
    (async () => {
      const [hist, a, m, c] = await Promise.all([
        fetchChatHistory(convId),
        fetchAgents(),
        fetchMessages(),
        fetchCheckpoints(),
      ]);
      if (!on) return;
      setAgents(a);
      setMessages(m);
      setCheckpoints(c);
      if (hist.length > 0) {
        setHasHistory(true);
        setMsgs(
          hist.map((h) => ({
            role: h.role === "user" ? "user" : "conductor",
            html:
              h.role === "user"
                ? esc(h.content)
                : renderMarkdown(cleanReply(h.content)),
            meta: h.role === "assistant" ? "earlier" : undefined,
          }))
        );
      }
      setLoaded(true);
    })();
    return () => {
      on = false;
    };
  }, [convId]);

  // live stream → keep sidebar fresh AND append in-thread agent messages
  useEffect(() => {
    let t: ReturnType<typeof setTimeout> | null = null;
    const debounced = () => {
      if (t) return;
      t = setTimeout(() => {
        t = null;
        refreshSidebar();
      }, 800);
    };
    return connectWS((e) => {
      if (e.type === "message" && "message" in e) {
        const m = (e as { message: ApiMessage }).message;
        const inThread =
          activeProject && m.metadata?.project_id === activeProject;
        // Scenario steps and tier-1 "routing" intermediaries are asides — they
        // pulse the workforce sidebar (via message recency) but should NOT clutter
        // the main chat thread, which stays = user prompts + the Brain's answers.
        const tags = m.metadata?.tags || [];
        const isAside = tags.includes("scenario") || tags.includes("routing");
        // While a request is in flight, surface scenario/routing asides as live
        // "what the agents are doing now" steps under the thinking indicator —
        // real backend activity reads far better than three bouncing dots.
        if (isAside && m.from !== "user" && sendingRef.current) {
          const body = (m.content || "").trim();
          if (body) {
            const line = `${humanize(m.from || "agent")} · ${body}`;
            setLiveSteps((p) =>
              p[p.length - 1] === line ? p : [...p.slice(-4), line]
            );
          }
        }
        if (inThread && m.from !== "user" && !isAside) {
          setMsgs((p) => [
            ...p,
            {
              role: "conductor",
              meta: `${m.from || "agent"} · live`,
              html: renderMarkdown(cleanReply(m.content || "")),
            },
          ]);
        }
      }
      if (
        e.type === "message" ||
        e.type === "checkpoint" ||
        e.type === "brain_event"
      )
        debounced();
    }, setConnected);
  }, [refreshSidebar, activeProject]);

  useEffect(() => {
    convRef.current?.scrollTo(0, convRef.current.scrollHeight);
  }, [msgs, thinking]);

  useEffect(() => {
    const ta = taRef.current;
    if (!ta) return;
    ta.style.height = "auto";
    ta.style.height = Math.min(ta.scrollHeight, 200) + "px";
  }, [text]);

  const pendingRoles = useMemo(
    () => new Set(mapCheckpoints(checkpoints).map((h) => h.tag.toLowerCase())),
    [checkpoints]
  );
  // A periodic tick (10s) drives the "active · Ns ago" counter forward and
  // lets the workforce row fall back to idle once the 60s activity window
  // expires, even when no new WS frames arrive.
  const [tick, setTick] = useState(0);
  useEffect(() => {
    const i = setInterval(() => setTick((n) => n + 1), 4_000);
    return () => clearInterval(i);
  }, []);
  // The workforce panel shows only the active company's roster, with proper
  // display names for the snake_case role ids the API returns.
  const roleSet = useMemo(() => new Set(cfg.roles), [cfg]);
  const companyAgents = useMemo(
    () =>
      agents
        .filter((a) => roleSet.has(a.role))
        .map((a) => ({ ...a, name: a.name || cfg.roleNames[a.role] || a.role })),
    [agents, roleSet, cfg]
  );
  const wf = useMemo(
    () => mapWorkforce(companyAgents, pendingRoles, messages),
    // tick intentionally included so we re-derive every 10s
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [companyAgents, pendingRoles, messages, tick]
  );
  // Counts + briefing reflect just the company roster, matching the panel.
  const epcAgents = companyAgents;
  const onlineCount = wf.filter((w) => w.state !== "offline").length;
  const activeCount = epcAgents.filter((a) => a.active).length;

  const fires = useMemo(() => mapFires(messages), [messages]);

  const briefing: Briefing = useMemo(
    () => composeBriefing(epcAgents, messages, checkpoints),
    // epcAgents derives from agents
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [agents, messages, checkpoints]
  );

  async function send(message?: string) {
    const m = (message ?? text).trim();
    if (!m || thinking || sendingRef.current) return; // ref-lock blocks the double-fire
    sendingRef.current = true;
    setText("");
    setHasHistory(true);
    setMsgs((p) => [...p, { role: "user", html: esc(m) }]);
    setThinking(true);
    setLiveSteps([]);
    let reply;
    try {
      reply = await sendChat(m, convId, cfg.scope);
    } finally {
      setThinking(false);
      sendingRef.current = false;
    }
    setMsgs((p) => [
      ...p,
      {
        role: "conductor",
        meta: reply?.action ? `routed · ${humanize(reply.action)}` : "offline · backend unreachable",
        html: reply?.response
          ? renderMarkdown(cleanReply(reply.response))
          : `<p>I'd classify this and dispatch to the right specialist, but <strong>the backend isn't reachable</strong> from this view. Start the Go server to route live.</p>`,
        reply: reply ?? undefined,
      },
    ]);
    if (reply?.suggestions && reply.suggestions.length > 0)
      setSuggest(reply.suggestions);
    if (reply?.project_id) setActiveProject(reply.project_id);
    refreshSidebar();
  }

  return (
    <div className="app lay-conductor" data-screen-label="Conductor">
      <Chrome
        pageTitle="Conductor"
        active="conductor"
        project={cfg.label}
        company={company}
        pack={cfg.short}
        packGlyph={cfg.glyph}
        packCount={epcAgents.length > 0 ? `${epcAgents.length} agents` : "no roster"}
      />

      <main className="main">
        <div className="conv-head">
          <div className="conductor-mark" />
          <div className="info">
            <h1>
              Conductor{" "}
              <span className="live-chip">
                <span className="pip" />
                {connected
                  ? `live · ${onlineCount} online${activeCount ? ` · ${activeCount} active` : ""}`
                  : "offline · reconnecting"}
              </span>
            </h1>
            <div className="sub">
              {company === "software"
                ? "Software Studio · describe an app and I'll run the team through the full build."
                : "The brain of the house. Ask anything · I dispatch to the right specialists."}
            </div>
          </div>
          <div className="actions">
            {company === "software" && <RunModeControl />}
            <button
              className="btn btn--secondary btn--sm"
              onClick={async () => {
                const c = await createConversation();
                // Even if the backend isn't reachable, mint a client-side id so
                // the UX always opens a fresh thread.
                setConvId(c?.id ?? `conv_${Date.now()}`);
              }}
            >
              <Plus />
              New conversation
            </button>
            <button
              className="btn btn--ghost btn--sm"
              onClick={() => setPaletteOpen(true)}
              title="History (⌘K)"
            >
              History
            </button>
          </div>
        </div>

        <div className="conv" ref={convRef}>
          {loaded && (
            <div className="msg conductor">
              <div className="conductor-mark sm" />
              <div className="body">
                <div className="meta">
                  <span className="nm">Conductor</span>
                  <span>·</span>
                  <span>{hasHistory ? "today's briefing" : "new conversation"}</span>
                </div>
                {/* A brand-new conversation starts clean — no carried-over thread
                    and no global activity dump. The briefing only leads a thread
                    that already has history. */}
                {hasHistory ? (
                  <BriefingCard briefing={briefing} />
                ) : (
                  <div style={{ color: "var(--text-secondary)", lineHeight: 1.6, fontSize: "14px" }}>
                    Fresh start. Tell me what you need and I will route it to the
                    right specialist. Pick a starter below, or just type.
                  </div>
                )}
              </div>
            </div>
          )}

          {msgs.map((m, i) =>
            m.role === "user" ? (
              <div className="msg user" key={i}>
                <div className="bubble" dangerouslySetInnerHTML={{ __html: m.html }} />
              </div>
            ) : (
              <div className="msg conductor" key={i}>
                <div className="conductor-mark sm" />
                <div className="body">
                  <div className="meta">
                    <span className="nm">Conductor</span>
                    <span>·</span>
                    <span>{m.meta || "just now"}</span>
                  </div>
                  {m.reply &&
                  ["delegate", "create_project", "escalate"].includes(
                    m.reply.action ?? ""
                  ) ? (
                    <ActionCard
                      reply={m.reply}
                      onPrimary={() => {
                        if (m.reply?.action === "escalate") {
                          // Run the deep-dive in-thread: the Brain produces the
                          // full cascading-impact mitigation plan (it has the
                          // escalation context from this conversation).
                          send(
                            "Yes, proceed — give me the full mitigation plan now: how the C-7 slab delay cascades through the Project Alpha schedule, the recovery options, and your recommended actions."
                          );
                        } else {
                          // create_project / delegate already started server-side;
                          // jump to Mission to watch the team work it.
                          window.location.href = "/mission";
                        }
                      }}
                      onSecondary={() => {
                        /* Skip = visual dismiss; the proposal simply isn't acted on. */
                      }}
                    />
                  ) : (
                    <div
                      className="text"
                      dangerouslySetInnerHTML={{ __html: m.html }}
                    />
                  )}
                  {m.reply?.suggestions && m.reply.suggestions.length > 0 && (
                    <div className="msg-sugs">
                      <div className="msg-sugs-label">Suggested follow-ups</div>
                      <div className="msg-sugs-row">
                        {m.reply.suggestions.map((s) => (
                          <button
                            className="msg-sug"
                            key={s}
                            onClick={() => send(s)}
                            title="Send this back to the brain"
                          >
                            {s}
                          </button>
                        ))}
                      </div>
                    </div>
                  )}
                </div>
              </div>
            )
          )}

          {thinking && (
            <div className="thinking-row">
              <div className="conductor-mark sm thinking" />
              {liveSteps.length > 0 ? (
                <div className="live-steps">
                  {liveSteps.map((s, i) => {
                    const dot = s.indexOf(" · ");
                    const who = dot > 0 ? s.slice(0, dot) : "";
                    const act = dot > 0 ? s.slice(dot + 3) : s;
                    const last = i === liveSteps.length - 1;
                    return (
                      <div
                        key={i}
                        className={`live-step${last ? " live-step--active" : ""}`}
                      >
                        <span className="live-step__pip" />
                        {who && <span className="live-step__who">{who}</span>}
                        <span className="live-step__act">{act}</span>
                      </div>
                    );
                  })}
                </div>
              ) : (
                <div className="dots">
                  <span />
                  <span />
                  <span />
                </div>
              )}
            </div>
          )}
        </div>

        <div className="composer">
          <div className="composer-inner">
            <textarea
              ref={taRef}
              value={text}
              onChange={(e) => setText(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === "Enter" && !e.shiftKey) {
                  e.preventDefault();
                  send();
                }
              }}
              placeholder="Ask anything — I'll route it to the right agent…"
              rows={1}
            />
            <button className="composer-send" onClick={() => send()}>
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.2" strokeLinecap="round">
                <path d="M5 12h14M13 5l7 7-7 7" />
              </svg>
            </button>
          </div>
          <div className="composer-hint">
            <span>
              <span className="kbd">@</span> mention · <span className="kbd">/</span>{" "}
              commands · <span className="kbd">⇧⏎</span> newline
            </span>
            <span>Conductor · routes via the brain</span>
          </div>
        </div>
      </main>

      <aside className="right">
        <div className="right-section">
          <h5>
            Workforce <span className="ct">{onlineCount} online</span>
          </h5>
          {wf.length === 0 ? (
            <div className="cond-empty">
              <div className="t">No agents in the roster yet.</div>
              <div className="d">
                Register agents or open{" "}
                <a className="lk" href="/lab">Lab</a> to seed demo state.
              </div>
            </div>
          ) : (
            wf.map((w, i) => (
              <div className={`wf-row ${w.rowTone}`} key={w.name + i}>
                <div className={`av ${w.state}`}>
                  {w.ini}
                  <div className="st" />
                </div>
                <div className="meta">
                  <div className="nm">{w.name}</div>
                  <div className="rl">{w.sub}</div>
                </div>
              </div>
            ))
          )}
        </div>

        <div className="right-section">
          <h5>
            Recent fires <span className="ct">{fires.length}</span>
          </h5>
          {fires.length === 0 ? (
            <div className="cond-empty">
              <div className="t">No system activity yet.</div>
              <div className="d">
                Triggers and scenarios will surface here as they fire. Use{" "}
                <a className="lk" href="/lab">Lab</a> to run one now.
              </div>
            </div>
          ) : (
            fires.map((f, i) => (
              <div className={`spark-row ${f.tone}`} key={f.t + i}>
                <span className="pip" />
                <div className="body">
                  <div className="t">{f.t}</div>
                  <div className="d">{f.d}</div>
                </div>
                <span className="ts">{f.ts}</span>
              </div>
            ))
          )}
        </div>

        <div className="right-section">
          <h5>Ask me</h5>
          <div className="suggest">
            {suggest.map((q) => (
              <div className="qq" key={q} onClick={() => send(q)}>
                <Bolt />
                {q}
                <span className="arrow">↵</span>
              </div>
            ))}
          </div>
        </div>

        <div className="right-foot">
          <span>
            <span className="dot" />
            brain · routes via classify
          </span>
          <span>v0.5</span>
        </div>
      </aside>

      <ConversationPalette
        open={paletteOpen}
        currentId={convId}
        onSelect={(id) => {
          setPaletteOpen(false);
          setConvId(id);
        }}
        onClose={() => setPaletteOpen(false)}
      />
    </div>
  );
}
