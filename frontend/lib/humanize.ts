// snake_case → "Snake Case" with sensible acronym preservation.
// Used everywhere a backend ID (role, action, tag, key) reaches a user.

const ACRONYMS = new Set([
  "rfi", "po", "pm", "hr", "qa", "qc", "hse", "kyc", "bec",
  "id", "cv", "url", "api", "ceo", "cfo", "coo", "ui", "ux",
  "mep", "epc", "sop", "kpi", "tco", "roi", "sla", "eta",
  "vp", "ip", "it", "ai", "ml", "db", "ar", "ap", "fk", "pk",
]);

const WORD_FIXUPS: Record<string, string> = {
  // domain-specific casing
  "applicants": "Applicants",
  "applicant": "Applicant",
  "vendors": "Vendors",
  "vendor": "Vendor",
  "milestone": "Milestone",
  "milestones": "Milestones",
  "inbox": "Inbox",
  "atlas": "Atlas",
};

function titleWord(w: string): string {
  if (!w) return w;
  const low = w.toLowerCase();
  if (ACRONYMS.has(low)) return low.toUpperCase();
  if (WORD_FIXUPS[low]) return WORD_FIXUPS[low];
  return low.charAt(0).toUpperCase() + low.slice(1);
}

/**
 * "snake_case_id" → "Snake Case ID"
 * "route_rfi"     → "Route RFI"
 * "match_to_po"   → "Match To PO"
 * "project_manager" → "Project Manager"
 * "" / null      → ""
 */
export function humanize(id?: string | null): string {
  if (!id) return "";
  return id
    .replace(/[._-]+/g, " ")
    .split(/\s+/)
    .filter(Boolean)
    .map(titleWord)
    .join(" ");
}

/**
 * "role:procurement · epc" — preserves structure but humanises each
 * underscored segment between the punctuation.
 */
export function humanizeMixed(s?: string | null): string {
  if (!s) return "";
  return s.replace(/[a-z][a-z0-9_-]*/gi, (m) => {
    // Only humanise tokens that contain at least one separator — leaves
    // existing words ("role", "epc", "active") alone unless they're IDs.
    if (m.includes("_") || m.includes("-")) return humanize(m);
    return m;
  });
}
