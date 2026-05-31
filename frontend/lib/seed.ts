// Seed payloads — used by /lab when a developer wants demo state in a scope.
// Pure data (no scenario / orchestration logic) — the same payloads our
// scripts/scenarios/load-epc-atlas.sh hits the API with.

export const ATLAS_APPLICANTS = [
  {
    id: "app_raj",
    name: "Raj Kumar",
    email: "raj.kumar@email.com",
    applied_for: "Site Supervisor for Highway Bridge Project",
    experience_years: 12,
    key_skills: ["civil", "bridge", "highway", "concrete", "rebar"],
    certifications: ["PE", "PMP"],
  },
  {
    id: "app_anita",
    name: "Anita Verma",
    email: "anita.verma@email.com",
    applied_for: "Senior Site Engineer",
    experience_years: 9,
    key_skills: ["civil", "structural", "autocad", "rfi", "qa"],
    certifications: ["PE"],
  },
  {
    id: "app_keshav",
    name: "Keshav Iyer",
    email: "keshav@email.com",
    applied_for: "HSE Officer",
    experience_years: 6,
    key_skills: ["hse", "safety", "compliance", "incident response"],
    certifications: ["NEBOSH"],
  },
];

export const ATLAS_VENDORS = [
  {
    id: "vendor_xyz",
    name: "XYZ Steel",
    domain: "vendorxyz.com",
    trusted_emails: ["ar@vendorxyz.com", "ahmed.rahman@vendorxyz.com"],
    contact_person: "Ahmed Rahman",
    contract_active: true,
  },
  {
    id: "vendor_alpha",
    name: "AlphaConcrete",
    domain: "alphaconcrete.io",
    trusted_emails: ["accounts@alphaconcrete.io"],
    contact_person: "Priya Shah",
    contract_active: true,
  },
  {
    id: "vendor_beta",
    name: "BetaElectrics",
    domain: "betaelectrics.com",
    trusted_emails: ["finance@betaelectrics.com"],
    contact_person: "Liu Chen",
    contract_active: false,
  },
];

export const ATLAS_MILESTONES = [
  {
    id: "ms_site_prep",
    title: "Site preparation",
    due_date: "2025-12-15",
    status: "completed",
    pct_complete: 100,
    owner: "site_engineer",
  },
  {
    id: "ms_foundation",
    title: "Foundation work",
    due_date: "2026-03-30",
    status: "in_progress",
    pct_complete: 78,
    owner: "site_engineer",
  },
  {
    id: "ms_slab_c7",
    title: "Pour C-7 slab",
    due_date: "2026-05-15",
    status: "in_progress",
    pct_complete: 40,
    owner: "site_engineer",
  },
  {
    id: "ms_steel_erect",
    title: "Structural steel erection",
    due_date: "2026-06-15",
    status: "pending",
    pct_complete: 0,
    owner: "site_engineer",
  },
  {
    id: "ms_grid_test",
    title: "Grid connection & testing",
    due_date: "2026-07-15",
    status: "pending",
    pct_complete: 0,
    owner: "qa_inspector",
  },
];
