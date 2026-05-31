import type { AgentVM, HitlVM, FeedVM, StatVM } from "./types";

/**
 * The reference's EPC content ("Atlas Construction · Site 02"). Used as the
 * fallback view-model so the Mission page always renders the full reference
 * design even when the backend has no live data yet. Live data supersedes
 * any of these sections the moment the backend returns it.
 */

export const SAMPLE_AGENT_GROUPS: { label: string; agents: AgentVM[] }[] = [
  {
    label: "Working",
    agents: [
      { initials: "SE", name: "Site Engineer", role: "site_engineer", state: "working", group: "Working", badgeTone: "info", badgeCount: 3 },
      { initials: "QA", name: "QA Inspector", role: "qa_inspector", state: "working", group: "Working", badgeTone: "info", badgeCount: 1 },
    ],
  },
  {
    label: "Awaiting you",
    agents: [
      { initials: "PR", name: "Procurement", role: "procurement", state: "awaiting", group: "Awaiting you", badgeTone: "wait", badgeCount: 1 },
      { initials: "SE", name: "Site Engineer", role: "+ email", state: "awaiting", group: "Awaiting you", badgeTone: "wait", badgeCount: 1 },
    ],
  },
  {
    label: "Live · idle",
    agents: [
      { initials: "PM", name: "Project Manager", role: "project_manager", state: "live", group: "Live · idle" },
      { initials: "HS", name: "HSE Officer", role: "hse", state: "live", group: "Live · idle" },
    ],
  },
  {
    label: "Disabled",
    agents: [{ initials: "HR", name: "HR Lead", role: "hr · off", state: "offline", group: "Disabled" }],
  },
  {
    label: "Platform",
    agents: [{ initials: "CE", name: "Chief Executive", role: "ceo · escalation", state: "live", group: "Platform" }],
  },
  {
    label: "Custom",
    agents: [{ initials: "CA", name: "Contracts Analyst", role: "contracts", state: "idle", group: "Custom" }],
  },
];

export const SAMPLE_STATS: StatVM[] = [
  { tone: "wait", n: "3", l: "await you" },
  { tone: "fire", n: "12", l: "trigger fires" },
  { tone: "", n: "2", l: "working" },
  { tone: "", n: "0", l: "errors · 7d" },
];

export const SAMPLE_HITL: HitlVM[] = [
  {
    id: "po",
    initials: "PR",
    state: "awaiting",
    title: "Issue PO #2294 to Vendor 117 · $24,500",
    detail: "⚡ from invoice file-watch · 12 line items · spawned task#a47e0d",
    tag: "PROCUREMENT",
    timer: "23h 14m",
    roleTag: "role:procurement · epc",
    desc: "Triggered by <code>invoice drop</code> on <code>vendor-117-PO-2294.pdf</code>. Agent matched 12 line items to the Q4 concrete admixture RFQ and is asking to commit funds.",
    payload: [
      ["action", "issue_purchase_order", ""],
      ["vendor_id", "117", "num"],
      ["amount", "24500.00 USD", "num"],
      ["items", "12", "num"],
      ["item", '"SikaPlast 320"', "str"],
      ["delivery", '"14 d · Site 02"', "str"],
    ],
    trace: "trg_4ax9 · sha 8c1e…d2 · spawned task#a47e0d · all events audited",
    primary: "Approve & issue PO",
    live: false,
  },
  {
    id: "overwrite",
    initials: "QA",
    state: "awaiting",
    title: "Overwrite drawings/A-101 v2.pdf with v3",
    detail: "4-item punchlist resolved · originated from QA review",
    tag: "QA_INSPECTOR",
    timer: "11h 02m",
    roleTag: "role:qa_inspector · epc",
    desc: "QA produced a punchlist of 4 items against <code>v2</code>. The agent wants to commit <code>v3</code> as the new canonical revision — this rewrites a tracked file.",
    payload: [
      ["action", "overwrite_file", ""],
      ["path", '"drawings/A-101.pdf"', "str"],
      ["from_sha", '"a4f9…02"', "str"],
      ["to_sha", '"7e21…b4"', "str"],
      ["notes", "4", "num"],
    ],
    trace: "task#a47e0b · review chain: site_engineer → qa_inspector → you",
    primary: "Approve & commit v3",
    live: false,
  },
  {
    id: "email",
    initials: "SE",
    state: "awaiting",
    title: "Email RFI response to Atlas Consulting",
    detail: "external_send · 4 attachments · 2 cc's",
    tag: "SITE_ENGINEER",
    timer: "2h 48m",
    roleTag: "role:site_engineer · epc",
    desc: "External send. Draft attaches 4 PDFs and replies to RFI #18 from Atlas regarding column tolerances. Per policy, all external sends are gated.",
    payload: [
      ["action", "send_email", ""],
      ["to", '"atlas-rfi@atlas.co"', "str"],
      ["cc", '"pm@site02, qa@site02"', "str"],
      ["attachments", "4", "num"],
      ["subject", '"RE: RFI #18 · column tol."', "str"],
    ],
    trace: "task#a47e08 · dangerous_action: external_send",
    primary: "Approve & send",
    live: false,
  },
];

export const SAMPLE_FEED: FeedVM[] = [
  {
    id: "f1",
    tone: "fire",
    triggered: "fire",
    ts: "14:02",
    head: ["invoice file-watch fired"],
    body: "<strong>Procurement</strong> spawned <code>task#a47e0d</code> from <code>vendor-117-PO-2294.pdf</code>. Summarised 12 line items · requesting authorisation to issue the PO.",
    card: { title: "vendor-117-PO-2294.pdf", detail: "sha 8c1e…d2 · 244 kB" },
  },
  {
    id: "f2",
    tone: "",
    ts: "13:48",
    head: ["you", "approved checkpoint"],
    body: "<strong>QA Inspector</strong> resumed work on <code>Drawing A-102 v3</code> and produced a 4-item punchlist.",
  },
  {
    id: "f3",
    tone: "info",
    ts: "13:30",
    head: ["Site Engineer", "queued review"],
    body: "Routed <code>Drawing A-102 v3</code> to <strong>QA Inspector</strong> with two notes on rebar spacing on grid C-7.",
  },
  {
    id: "f4",
    tone: "fire",
    triggered: "fire",
    ts: "13:21",
    head: ["invoice file-watch fired"],
    body: "<strong>Procurement</strong> classified <code>vendor-04-credit-note.pdf</code> as a credit note. Auto-filed under <code>vendors/04/</code>.",
  },
  {
    id: "f5",
    tone: "live",
    triggered: "live",
    ts: "07:00",
    head: ["daily schedule check fired"],
    body: "<strong>Project Manager</strong> reported on schedule. 14 milestones · 1 trending late: <code>Pour C-7 slab</code> +3 days.",
  },
];
