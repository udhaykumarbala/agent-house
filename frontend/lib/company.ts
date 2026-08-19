// Two-door company model. Agent House runs two companies behind one app: a
// SOFTWARE studio (the agent-driven SDLC that builds apps) and the EPC
// construction company (Alredaa). The active company drives the workforce
// roster, the scope sent to the Brain, and the chrome labels.

export type Company = "software" | "epc";

export interface CompanyConfig {
  id: Company;
  label: string; // shown in the switcher + as the project name
  short: string; // pack pill label
  glyph: string; // pack pill glyph
  scope: string; // sent to /api/chat so the Brain routes per company
  roles: string[]; // workforce roster (filters /api/agents)
  /** Display names for snake_case role ids the API returns. */
  roleNames: Record<string, string>;
  tagline: string;
}

export const COMPANIES: Record<Company, CompanyConfig> = {
  software: {
    id: "software",
    label: "Software Studio",
    short: "SW",
    glyph: "SW",
    scope: "software",
    roles: [
      "ceo",
      "pm",
      "architect",
      "ux",
      "ui",
      "security",
      "senior_dev",
      "junior_dev",
    ],
    roleNames: {
      ceo: "Chief Executive Officer",
      pm: "Product Manager",
      architect: "Architect",
      ux: "UX Designer",
      ui: "UI Designer",
      security: "Security Expert",
      senior_dev: "Senior Developer",
      junior_dev: "Junior Developer",
    },
    tagline: "Ship software — PRD to production, agent-driven.",
  },
  epc: {
    id: "epc",
    label: "Alredaa",
    short: "EPC",
    glyph: "AL",
    scope: "atlas-site",
    roles: [
      "ceo",
      "procurement",
      "project_manager",
      "site_engineer",
      "hse",
      "qa_inspector",
      "hr",
    ],
    roleNames: {
      ceo: "Chief Executive Officer",
      procurement: "Procurement Lead",
      project_manager: "Project Manager",
      site_engineer: "Site Engineer",
      hse: "HSE Officer",
      qa_inspector: "QA Inspector",
      hr: "HR Manager",
    },
    tagline: "Run the site — inbox, vendors, schedule, HSE.",
  },
};

const KEY = "agenthouse.company";

/** Resolve the active company: ?company= URL param wins (and is persisted),
 * else localStorage, else the software studio (the current build focus). */
export function getCompany(): Company {
  if (typeof window === "undefined") return "software";
  const q = new URL(window.location.href).searchParams.get("company");
  if (q === "software" || q === "epc") {
    localStorage.setItem(KEY, q);
    return q;
  }
  return localStorage.getItem(KEY) === "epc" ? "epc" : "software";
}

export function setCompany(c: Company): void {
  if (typeof window !== "undefined") localStorage.setItem(KEY, c);
}

export function companyConfig(c: Company): CompanyConfig {
  return COMPANIES[c];
}
