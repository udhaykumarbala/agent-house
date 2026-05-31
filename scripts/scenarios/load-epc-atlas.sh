#!/usr/bin/env bash
# Atlas Construction · Site 02 — an EPC demo scenario.
# Seeds the global inbox with five realistic emails (vendor impersonation
# alert, RFI from Atlas Consulting, vendor invoice, client progress request,
# civil-engineer application) and adds two scoped cron jobs to the running
# scheduler. Open /conductor afterwards — the briefing reflects it.
#
# Re-runnable: each email uses a deterministic id (email_atlas_<n>) so a
# second run overwrites instead of duplicating.
set -euo pipefail

HOST="${HOST:-http://localhost:8099}"
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
INBOX="$ROOT/projects/inbox"
PROJ="$ROOT/projects/atlas-site"

mkdir -p "$INBOX" "$PROJ/.tasks" "$PROJ/inbox" "$PROJ/drawings"

now() { date -u +"%Y-%m-%dT%H:%M:%SZ"; }
ago() { # ago <minutes>
  if date -v -0M >/dev/null 2>&1; then
    date -u -v "-${1}M" +"%Y-%m-%dT%H:%M:%SZ"
  else
    date -u -d "${1} minutes ago" +"%Y-%m-%dT%H:%M:%SZ"
  fi
}

write_email() {
  local id="$1" from="$2" name="$3" subj="$4" cat="$5" trust="$6" reason="$7" ts="$8" body="$9" vendor_id="${10:-}"
  local path="$INBOX/${id}.json"
  jq -n \
    --arg id "$id" --arg from "$from" --arg name "$name" --arg to "ops@atlas-construction.com" \
    --arg subj "$subj" --arg body "$body" --arg date "$ts" \
    --arg cat "$cat" --arg trust "$trust" --arg reason "$reason" --arg vid "$vendor_id" \
    '{id:$id, from:$from, from_name:$name, to:$to, subject:$subj, body:$body,
      date:$date, read:false, replied:false, category:$cat,
      vendor_id:(if $vid=="" then null else $vid end),
      trust_status:$trust, trust_reason:$reason, direction:"inbound"}' \
    > "$path"
  echo "  ✓ $(basename "$path") · $cat · $trust"
}

echo "→ seeding Atlas Construction · Site 02 inbox at $INBOX"

write_email "email_atlas_1_impersonation" \
  "ahmed.r@gmail.com" "Ahmed Rahman — XYZ Steel" \
  "URGENT — Updated Bank Account Details for Payment" \
  "vendor" "impersonation" \
  "IMPERSONATION RISK — claims XYZ Steel Corp but sender is gmail.com; real domain is @vendorxyz.com" \
  "$(ago 38)" \
  "Dear Accounts Team,

This is Ahmed Rahman from XYZ Steel Corp. Due to a change in our banking arrangements, please update our payment details immediately.

New Bank Details:
- Bank: First National Bank
- Account: 8834-2291-0057
- SWIFT: FNBKUS33

Please use these details for all future payments including the pending Invoice #ST-0847 (\$340,000). This change is effective immediately.

Urgent regards,
Ahmed Rahman
XYZ Steel Corp" \
  "vendor_xyz_imposter"

write_email "email_atlas_2_rfi" \
  "rfi@atlasconsulting.co" "Atlas Consulting · RFI Desk" \
  "RFI #18 — Column tolerances on grid C-7" \
  "client" "trusted" "Trusted consulting partner · on whitelist" \
  "$(ago 162)" \
  "Hi Site 02 team,

Per the drawings revision A-102 v3, we have an open question on column tolerance limits on grid C-7. Specifically: maximum allowable plumb deviation for cast-in-place columns is shown as ±10mm, but the structural spec calls out ±6mm.

Please confirm the governing tolerance before pour. Attaching the relevant detail sheets (4 PDFs). cc: pm@site02, qa@site02.

Best,
RFI Desk · Atlas Consulting"

write_email "email_atlas_3_invoice" \
  "ar@vendorxyz.com" "XYZ Steel · Accounts Receivable" \
  "Invoice #ST-0847 — 150 t structural steel · Project Alpha" \
  "vendor" "trusted" "Verified sender on @vendorxyz.com" \
  "$(ago 240)" \
  "Hello Procurement,

Invoice attached for Order #ST-2026-0847 — 150 metric tons of structural steel for Project Alpha (Solar Farm). PO reference 2294.

Amount: \$340,000.00
Net 30 days.
Bank: Standard Chartered, A/C 1100-2294-77 (on file).

Note: there is a 2-week delivery delay on the next batch — separate email follows.

Regards,
A/R · XYZ Steel" \
  "vendor_xyz"

write_email "email_atlas_4_client" \
  "sarah.jones@clientabc.com" "Sarah Jones — Client ABC" \
  "Request for Phase 2 Progress Update — Project Alpha" \
  "client" "trusted" "Long-standing client; board update cadence" \
  "$(ago 1440)" \
  "Hi team,

Could you send the Phase 2 progress update for **Project Alpha (Solar Farm)** by end of week? Our board reviews on April 1.

We need: % completion, current-phase status, budget utilisation, updated timeline, and risk register.

Thanks,
Sarah Jones
Director · Client ABC"

write_email "email_atlas_5_applicant" \
  "raj.kumar@email.com" "Raj Kumar" \
  "Application — Site Supervisor · Highway Bridge Project" \
  "applicant" "new_contact" "First contact from this sender" \
  "$(ago 5760)" \
  "Dear HR,

I am applying for the Site Supervisor position on the Highway Bridge Project (Project Gamma).

Brief summary:
- 12 years civil engineering experience
- 4 completed bridge projects (including Highway-9 viaduct, 2024)
- PE and PMP certified
- Open to relocate to Site 02

CV and references attached. Available for interview from May 22.

Sincerely,
Raj Kumar"

echo
echo "→ adding scoped cron jobs for atlas-site"
post_cron() {
  local body="$1"
  local code
  code=$(curl -s -o /tmp/cron.out -w "%{http_code}" -X POST "$HOST/api/cron/" \
    -H 'Content-Type: application/json' -d "$body")
  if [[ "$code" != "200" ]]; then
    echo "  ✗ POST failed ($code): $(cat /tmp/cron.out)"
    return 1
  fi
  echo "  ✓ $(jq -r '.job | "\(.id)  \(.schedule)  \(.task[:64])"' < /tmp/cron.out)"
}

post_cron '{"task":"Daily schedule slip check for Atlas Site 02 — scan milestones, flag late, propose mitigations","schedule":"24h","agent_role":"project_manager","project_id":"atlas-site"}' || true
post_cron '{"task":"Hourly inbox sweep for Atlas Site 02 — route vendor / applicant / RFI emails, flag impersonation risks","schedule":"1h","agent_role":"ceo","project_id":"atlas-site"}' || true

# ── HR applicants ────────────────────────────────────────────────────
echo
echo "→ seeding HR applicants for atlas-site (via /api/cap/hr/applicants)"
post_app() {
  curl -s -X POST "$HOST/api/cap/hr/applicants?scope=atlas-site" \
    -H 'Content-Type: application/json' -d "$1" \
    | jq -r '"  ✓ " + .name + " (" + .applied_for + ", " + (.experience_years|tostring) + "y) [" + .id + "]"'
}
post_app '{"id":"app_raj","name":"Raj Kumar","email":"raj.kumar@email.com","applied_for":"Site Supervisor for Highway Bridge Project","experience_years":12,"key_skills":["civil","bridge","highway","concrete","rebar"],"certifications":["PE","PMP"]}'
post_app '{"id":"app_anita","name":"Anita Verma","email":"anita.verma@email.com","applied_for":"Senior Site Engineer","experience_years":9,"key_skills":["civil","structural","autocad","rfi","qa"],"certifications":["PE"]}'
post_app '{"id":"app_keshav","name":"Keshav Iyer","email":"keshav@email.com","applied_for":"HSE Officer","experience_years":6,"key_skills":["hse","safety","compliance","incident response"],"certifications":["NEBOSH"]}'

# ── Procurement vendors ──────────────────────────────────────────────
echo
echo "→ seeding Procurement vendors for atlas-site (via /api/cap/procurement/vendors)"
post_vendor() {
  curl -s -X POST "$HOST/api/cap/procurement/vendors?scope=atlas-site" \
    -H 'Content-Type: application/json' -d "$1" \
    | jq -r '"  ✓ " + .name + " (@" + .domain + ", active=" + (.contract_active|tostring) + ") [" + .id + "]"'
}
post_vendor '{"id":"vendor_xyz","name":"XYZ Steel","domain":"vendorxyz.com","trusted_emails":["ar@vendorxyz.com","ahmed.rahman@vendorxyz.com"],"contact_person":"Ahmed Rahman","contract_active":true}'
post_vendor '{"id":"vendor_alpha","name":"AlphaConcrete","domain":"alphaconcrete.io","trusted_emails":["accounts@alphaconcrete.io"],"contact_person":"Priya Shah","contract_active":true}'
post_vendor '{"id":"vendor_beta","name":"BetaElectrics","domain":"betaelectrics.com","trusted_emails":["finance@betaelectrics.com"],"contact_person":"Liu Chen","contract_active":false}'

# ── Schedule milestones (one deliberately slipping) ─────────────────
echo
echo "→ seeding Schedule milestones for atlas-site (via /api/cap/schedule/milestones)"
post_ms() {
  curl -s -X POST "$HOST/api/cap/schedule/milestones?scope=atlas-site" \
    -H 'Content-Type: application/json' -d "$1" \
    | jq -r '"  ✓ " + .title + " (due " + .due_date + ", " + .status + ", " + (.pct_complete|tostring) + "%) [" + .id + "]"'
}
# Past-due deliberately for slip detection (date is in the past from today's POV)
post_ms '{"id":"ms_site_prep","title":"Site preparation","due_date":"2025-12-15","status":"completed","pct_complete":100,"owner":"site_engineer"}'
post_ms '{"id":"ms_foundation","title":"Foundation work","due_date":"2026-03-30","status":"in_progress","pct_complete":78,"owner":"site_engineer"}'
post_ms '{"id":"ms_slab_c7","title":"Pour C-7 slab","due_date":"2026-05-15","status":"in_progress","pct_complete":40,"owner":"site_engineer"}'
post_ms '{"id":"ms_steel_erect","title":"Structural steel erection","due_date":"2026-06-15","status":"pending","pct_complete":0,"owner":"site_engineer"}'
post_ms '{"id":"ms_grid_test","title":"Grid connection & testing","due_date":"2026-07-15","status":"pending","pct_complete":0,"owner":"qa_inspector"}'

echo
echo "→ done. Seeded:"
echo "   - 5 emails"
echo "   - 2 cron jobs"
echo "   - $(curl -s "$HOST/api/cap/hr/applicants?scope=atlas-site" | jq -r .count) applicants"
echo "   - $(curl -s "$HOST/api/cap/procurement/vendors?scope=atlas-site" | jq -r .count) vendors"
echo "   - $(curl -s "$HOST/api/cap/schedule/milestones?scope=atlas-site" | jq -r .count) milestones"
echo "   - $(curl -s "$HOST/api/cap/schedule/slips?scope=atlas-site" | jq -r .count) detected slips"
echo "   Open $HOST/conductor — the briefing reflects new fires + scenarios."
