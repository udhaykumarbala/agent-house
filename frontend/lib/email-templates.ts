// Pre-canned email templates the Lab uses to drop realistic inbound
// messages into the inbox. Each template targets one or more scenarios so
// the user can immediately see the downstream effect.

export interface EmailTemplate {
  id: string;
  label: string;
  description: string;
  /** Scenarios that consume this kind of email. Used to tell the user
   *  what to try after sending. */
  scenarios: string[];
  payload: {
    from: string;
    from_name?: string;
    to?: string;
    subject: string;
    body: string;
    category: "vendor" | "applicant" | "client" | "internal";
    trust_status?: "trusted" | "impersonation" | "new_contact";
    trust_reason?: string;
    vendor_id?: string;
  };
}

export const EMAIL_TEMPLATES: EmailTemplate[] = [
  {
    id: "client_estimation",
    label: "Client — new project estimation request",
    description:
      "A prospective client asks for a manpower + cost estimation for a new build. Drives the estimate_project story: AI drafts the crew, HR checks live HRMS availability and costs it.",
    scenarios: ["estimate_project"],
    payload: {
      from: "projects@newco-industrial.com",
      from_name: "NewCo Industrial · Projects",
      to: "estimation@alredaa.com",
      subject: "Request for Estimation — Warehouse Construction, Dammam",
      body: `Dear Alredaa team,

We are planning a new 4,500 sqm warehouse (2 bays, mezzanine office block) at Dammam 2nd Industrial City and would like your manpower and cost estimation.

Scope highlights:
- Structural concrete + steel erection
- MEP first fix and second fix
- Target duration: 6 months from mobilization

Please share your estimated crew composition, availability, and indicative monthly cost so we can shortlist contractors this month.

Best regards,
Projects Team
NewCo Industrial`,
      category: "client",
      trust_status: "new_contact",
      trust_reason: "First contact from @newco-industrial.com — no prior thread",
    },
  },
  {
    id: "vendor_impersonation",
    label: "Vendor — impersonation (classic BEC)",
    description:
      "Someone claiming to be a known vendor, asking to change bank details, sending from a public-mail domain.",
    scenarios: ["validate_invoice", "process_inbox"],
    payload: {
      from: "ahmed.r@gmail.com",
      from_name: "Ahmed Rahman — XYZ Steel",
      to: "ops@atlas-construction.com",
      subject: "URGENT — Updated Bank Account Details for Payment",
      body: `Dear Accounts Team,

This is Ahmed Rahman from XYZ Steel Corp. Due to a change in our banking arrangements, please update our payment details immediately.

New Bank Details:
- Bank: First National Bank
- Account: 8834-2291-0057
- SWIFT: FNBKUS33

Please use these details for all future payments, starting with the pending Invoice #ST-0847 ($340,000). This change is effective immediately.

Urgent regards,
Ahmed Rahman
XYZ Steel Corp`,
      category: "vendor",
      trust_status: "impersonation",
      trust_reason:
        "Claims XYZ Steel but sender domain is gmail.com — known vendor is @vendorxyz.com",
      vendor_id: "vendor_xyz_imposter",
    },
  },
  {
    id: "vendor_invoice",
    label: "Vendor — trusted invoice",
    description:
      "Legitimate invoice from a known vendor on its registered domain.",
    scenarios: ["validate_invoice", "process_inbox"],
    payload: {
      from: "ar@vendorxyz.com",
      from_name: "XYZ Steel · Accounts Receivable",
      to: "procurement@atlas-construction.com",
      subject: "Invoice #ST-0847 — 150 t structural steel · Project Alpha",
      body: `Hello Procurement,

Invoice attached for Order #ST-2026-0847 — 150 metric tons of structural steel for Project Alpha (Solar Farm). PO reference 2294.

Amount: $340,000.00
Net 30 days.
Bank: Standard Chartered, A/C 1100-2294-77 (on file).

Regards,
A/R · XYZ Steel`,
      category: "vendor",
      trust_status: "trusted",
      trust_reason: "Verified sender on registered domain @vendorxyz.com",
      vendor_id: "vendor_xyz",
    },
  },
  {
    id: "client_rfi",
    label: "Client — RFI (structural)",
    description:
      "Request For Information from the client about a structural detail; routes through route_rfi to discipline classification.",
    scenarios: ["route_rfi", "process_inbox"],
    payload: {
      from: "rfi@atlasconsulting.co",
      from_name: "Atlas Consulting · RFI Desk",
      to: "site02@atlas-construction.com",
      subject: "RFI #18 — Column tolerances on grid C-7",
      body: `Hi Site 02 team,

Per the drawings revision A-102 v3, we have an open question on column tolerance limits on grid C-7. Specifically: maximum allowable plumb deviation for cast-in-place columns is shown as ±10mm, but the structural spec calls out ±6mm.

Please confirm the governing tolerance before pour. Attaching the relevant detail sheets (4 PDFs). cc: pm@site02, qa@site02.

Best,
RFI Desk · Atlas Consulting`,
      category: "client",
      trust_status: "trusted",
      trust_reason: "Long-standing consulting partner",
    },
  },
  {
    id: "client_progress_request",
    label: "Client — progress update request",
    description:
      "Client asking for the Phase 2 board update. PM-flagged for client response.",
    scenarios: ["process_inbox"],
    payload: {
      from: "sarah.jones@clientabc.com",
      from_name: "Sarah Jones — Client ABC",
      to: "pm@atlas-construction.com",
      subject: "Request for Phase 2 Progress Update — Project Alpha",
      body: `Hi team,

Could you send the Phase 2 progress update for Project Alpha (Solar Farm) by end of week? Our board reviews on April 1.

We need: % completion, current-phase status, budget utilisation, updated timeline, and risk register.

Thanks,
Sarah Jones
Director · Client ABC`,
      category: "client",
      trust_status: "trusted",
      trust_reason: "Long-standing client; board update cadence",
    },
  },
  {
    id: "applicant_civil",
    label: "Applicant — Civil / bridge supervisor",
    description:
      "Inbound CV for a Site Supervisor role. HR will store; process_applicants can match against open JDs.",
    scenarios: ["process_applicants", "process_inbox", "route_rfi"],
    payload: {
      from: "raj.kumar@email.com",
      from_name: "Raj Kumar",
      to: "hr@atlas-construction.com",
      subject: "Application — Site Supervisor · Highway Bridge Project",
      body: `Dear HR,

I am applying for the Site Supervisor position on the Highway Bridge Project (Project Gamma).

Brief summary:
- 12 years civil engineering experience
- 4 completed bridge projects (including Highway-9 viaduct, 2024)
- PE and PMP certified
- Open to relocate to Site 02

CV and references attached. Available for interview from May 22.

Sincerely,
Raj Kumar`,
      category: "applicant",
      trust_status: "new_contact",
      trust_reason: "First contact from this sender",
    },
  },
  {
    id: "applicant_hse",
    label: "Applicant — HSE Officer",
    description: "HSE specialist applying. Good demo for skill-matching.",
    scenarios: ["process_applicants", "process_inbox"],
    payload: {
      from: "keshav@email.com",
      from_name: "Keshav Iyer",
      to: "hr@atlas-construction.com",
      subject: "Application — HSE Officer",
      body: `Hello,

I'm applying for the HSE Officer position. Six years on chemical and infrastructure sites, NEBOSH IGC certified, incident-response experience across three jurisdictions.

CV attached. Available immediately.

— Keshav Iyer`,
      category: "applicant",
      trust_status: "new_contact",
    },
  },
  {
    id: "vendor_delivery_delay",
    label: "Vendor — delivery delay notice",
    description:
      "Legitimate vendor announcing a 2-week steel delivery delay. Useful for schedule_check downstream.",
    scenarios: ["process_inbox", "schedule_check"],
    payload: {
      from: "ahmed.rahman@vendorxyz.com",
      from_name: "Ahmed Rahman — XYZ Steel",
      to: "procurement@atlas-construction.com",
      subject: "Steel Delivery Update — Project Alpha",
      body: `Hi team,

Heads up: we have a 2-week delay on the next steel batch for Project Alpha (Order #ST-2026-0847). Mill output is constrained; current ETA is April 19 (was April 5).

Apologies for the inconvenience. Happy to discuss alternative sourcing for the C-7 slab pour if it's critical.

Regards,
Ahmed`,
      category: "vendor",
      trust_status: "trusted",
      vendor_id: "vendor_xyz",
    },
  },
  {
    id: "internal_safety_incident",
    label: "Internal — HSE incident report",
    description: "Safety incident reported from site to ops.",
    scenarios: ["process_inbox"],
    payload: {
      from: "site-hse@atlas-construction.com",
      from_name: "Site 02 · HSE",
      to: "hse@atlas-construction.com",
      subject: "Incident report — minor near-miss on grid C-7",
      body: `Incident summary:

- Date/time: today, 09:42 local
- Location: grid C-7, formwork inspection area
- Type: near-miss — loose rebar tie spotted before pour
- People affected: none
- Action taken: corrected on the spot; toolbox talk scheduled

No injury, no work-stopper. Filing per protocol.`,
      category: "internal",
    },
  },
];
