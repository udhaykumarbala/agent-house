#!/usr/bin/env bash
# End-to-end test harness for the capability + scenario layers.
# Runs against a live server. Asserts shape + business outcomes for every
# capability endpoint and every scenario. Prints a pass/fail table; exits
# non-zero on any failure.
set -uo pipefail

HOST="${HOST:-http://localhost:8099}"
SCOPE="test-scope-$(date +%s)"

PASS=0
FAIL=0
declare -a FAILS

ok()  { PASS=$((PASS+1)); printf "  \033[32m✓\033[0m %s\n" "$1"; }
bad() { FAIL=$((FAIL+1)); FAILS+=("$1"); printf "  \033[31m✗\033[0m %s\n" "$1"; if [[ -n "${2:-}" ]]; then printf "      \033[2m%s\033[0m\n" "$2"; fi; }

# assert_eq <label> <expected> <actual>
assert_eq() {
  local label="$1" want="$2" got="$3"
  if [[ "$want" == "$got" ]]; then ok "$label"; else bad "$label" "expected=$want got=$got"; fi
}

# assert_contains <label> <haystack> <needle>
assert_contains() {
  local label="$1" hay="$2" needle="$3"
  if echo "$hay" | grep -qF "$needle"; then ok "$label"; else bad "$label" "needle '$needle' not found in: $(echo "$hay" | head -c 200)"; fi
}

# assert_status <label> <expected-code> <method> <path> [body]
assert_status() {
  local label="$1" want="$2" method="$3" path="$4" body="${5:-}"
  local code
  if [[ -n "$body" ]]; then
    code=$(curl -s -o /tmp/test.out -w "%{http_code}" -X "$method" "$HOST$path" \
      -H 'Content-Type: application/json' -d "$body")
  else
    code=$(curl -s -o /tmp/test.out -w "%{http_code}" -X "$method" "$HOST$path")
  fi
  if [[ "$code" == "$want" ]]; then ok "$label ($method $path → $code)"
  else bad "$label ($method $path)" "expected=$want got=$code body=$(head -c 200 /tmp/test.out)"; fi
}

section() { printf "\n\033[1m── %s ──\033[0m\n" "$1"; }

section "0. server reachability"
assert_status "server up"          200 GET /api/status
assert_status "capabilities index" 200 GET /api/capabilities
assert_status "scenarios index"    200 GET /api/scenarios

section "1. HR capability — store + list + match + update"
# Clean slate: seed three applicants into the test scope.
for body in \
  '{"id":"t_raj","name":"Raj Kumar","email":"raj@x.com","applied_for":"Site Supervisor Bridge","experience_years":12,"key_skills":["civil","bridge","highway"],"certifications":["PE","PMP"]}' \
  '{"id":"t_anita","name":"Anita Verma","email":"anita@x.com","applied_for":"Senior Site Engineer","experience_years":9,"key_skills":["civil","structural","rfi"],"certifications":["PE"]}' \
  '{"id":"t_hse","name":"Keshav Iyer","email":"k@x.com","applied_for":"HSE Officer","experience_years":6,"key_skills":["hse","safety"],"certifications":["NEBOSH"]}' ; do
  curl -s -X POST "$HOST/api/cap/hr/applicants?scope=$SCOPE" -H 'Content-Type: application/json' -d "$body" >/dev/null
done

LIST=$(curl -s "$HOST/api/cap/hr/applicants?scope=$SCOPE")
COUNT=$(echo "$LIST" | jq -r .count)
assert_eq "list returns 3 applicants" 3 "$COUNT"
assert_contains "list contains Raj Kumar" "$LIST" "Raj Kumar"

# Match against a JD that should find Raj first (civil + bridge + highway + PE)
MATCH=$(curl -s -X POST "$HOST/api/cap/hr/match?scope=$SCOPE" -H 'Content-Type: application/json' \
  -d '{"title":"Site Supervisor","must_have":["civil","bridge","highway"],"nice_to_have":["pmp"],"min_years":8}')
MCOUNT=$(echo "$MATCH" | jq -r .count)
TOP=$(echo "$MATCH" | jq -r '.matches[0].applicant.name')
TOPSCORE=$(echo "$MATCH" | jq -r '.matches[0].score')
assert_eq "match returns ≥1 result" "true" "$(echo "$MATCH" | jq -r '.count>0')"
assert_eq "top match is Raj Kumar" "Raj Kumar" "$TOP"
[[ "$TOPSCORE" -gt 5 ]] && ok "top match score > 5 (got $TOPSCORE)" || bad "top match score > 5" "got $TOPSCORE"

# Update status
UP=$(curl -s -X PATCH "$HOST/api/cap/hr/applicants/t_raj?scope=$SCOPE" \
  -H 'Content-Type: application/json' -d '{"status":"shortlisted"}')
assert_eq "update_status sets shortlisted" "shortlisted" "$(echo "$UP" | jq -r .status)"

# Verify persistence
PERSIST=$(curl -s "$HOST/api/cap/hr/applicants?scope=$SCOPE" | jq -r '.applicants[] | select(.id=="t_raj") | .status')
assert_eq "status persisted across requests" "shortlisted" "$PERSIST"

section "2. Procurement capability — vendors + invoice validation"
for body in \
  '{"id":"v_xyz","name":"XYZ Steel","domain":"vendorxyz.com","trusted_emails":["ar@vendorxyz.com"],"contact_person":"Ahmed Rahman","contract_active":true}' \
  '{"id":"v_alpha","name":"AlphaConcrete","domain":"alphaconcrete.io","trusted_emails":[],"contact_person":"Priya","contract_active":true}' ; do
  curl -s -X POST "$HOST/api/cap/procurement/vendors?scope=$SCOPE" -H 'Content-Type: application/json' -d "$body" >/dev/null
done
VC=$(curl -s "$HOST/api/cap/procurement/vendors?scope=$SCOPE" | jq -r .count)
assert_eq "2 vendors stored" 2 "$VC"

# Trusted sender → trusted result, recommendation match_to_po
TRUSTED=$(curl -s -X POST "$HOST/api/cap/procurement/validate-invoice?scope=$SCOPE" \
  -H 'Content-Type: application/json' \
  -d '{"sender_email":"ar@vendorxyz.com","vendor_id":"v_xyz","amount":24500}')
assert_eq "trusted invoice: trusted=true"  "true"        "$(echo "$TRUSTED" | jq -r .trusted)"
assert_eq "trusted invoice: impersonation=false" "false" "$(echo "$TRUSTED" | jq -r .impersonation_risk)"
assert_eq "trusted invoice: recommend match_to_po" "match_to_po" "$(echo "$TRUSTED" | jq -r .recommendation)"

# Impersonation: same vendor, gmail sender
SPOOF=$(curl -s -X POST "$HOST/api/cap/procurement/validate-invoice?scope=$SCOPE" \
  -H 'Content-Type: application/json' \
  -d '{"sender_email":"ahmed.r@gmail.com","vendor_id":"v_xyz","amount":340000}')
assert_eq "spoof: trusted=false"             "false" "$(echo "$SPOOF" | jq -r .trusted)"
assert_eq "spoof: impersonation_risk=true"   "true"  "$(echo "$SPOOF" | jq -r .impersonation_risk)"
assert_eq "spoof: recommend block_vendor_and_notify_finance" \
  "block_vendor_and_notify_finance" "$(echo "$SPOOF" | jq -r .recommendation)"
assert_contains "spoof: reasons mention public mail provider" "$SPOOF" "public mail provider"

section "3. Schedule capability — milestones + slips"
TODAY=$(date -u +%Y-%m-%d)
# One slipping (3 weeks ago, in_progress), one future (1 month out)
for body in \
  '{"id":"m_late","title":"Pour C-7 slab","due_date":"2026-04-01","status":"in_progress","pct_complete":40,"owner":"site_engineer"}' \
  "{\"id\":\"m_future\",\"title\":\"Grid connection\",\"due_date\":\"2026-08-15\",\"status\":\"pending\",\"pct_complete\":0,\"owner\":\"qa_inspector\"}" ; do
  curl -s -X POST "$HOST/api/cap/schedule/milestones?scope=$SCOPE" -H 'Content-Type: application/json' -d "$body" >/dev/null
done
MS=$(curl -s "$HOST/api/cap/schedule/milestones?scope=$SCOPE" | jq -r .count)
assert_eq "2 milestones stored" 2 "$MS"

SLIPS=$(curl -s "$HOST/api/cap/schedule/slips?scope=$SCOPE")
SC=$(echo "$SLIPS" | jq -r .count)
[[ "$SC" -ge 1 ]] && ok "≥1 slip detected (got $SC)" || bad "≥1 slip detected" "got $SC"
TOP_SLIP=$(echo "$SLIPS" | jq -r '.slips[0].milestone.title')
assert_eq "top slip is Pour C-7 slab" "Pour C-7 slab" "$TOP_SLIP"
SEV=$(echo "$SLIPS" | jq -r '.slips[0].severity')
assert_contains "top slip severity is critical/moderate" "critical moderate" "$SEV"

section "4. Scenario: process_applicants (Conductor → HR → suggestions)"
RES=$(curl -s -X POST "$HOST/api/scenario/process_applicants/run?scope=$SCOPE" \
  -H 'Content-Type: application/json' \
  -d '{"request":"Hire a Site Supervisor for the Highway Bridge Project"}')
assert_eq "scenario ok"                        "true" "$(echo "$RES" | jq -r .ok)"
assert_eq "scenario has 2 steps"               2      "$(echo "$RES" | jq -r '.steps|length')"
assert_eq "step 1 is conductor:decompose"      "conductor" "$(echo "$RES" | jq -r '.steps[0].agent')"
assert_eq "step 2 is hr:match_against_jd"      "hr" "$(echo "$RES" | jq -r '.steps[1].agent')"
SUG=$(echo "$RES" | jq -r '.suggestions|length')
[[ "$SUG" -ge 1 ]] && ok "scenario emitted ≥1 suggestion (got $SUG)" || bad "scenario emitted ≥1 suggestion" "got $SUG"
assert_contains "suggests schedule_interviews" "$(echo "$RES" | jq -r '.suggestions[].action')" "schedule_interviews"

section "5. Scenario: validate_invoice — impersonation path"
RES2=$(curl -s -X POST "$HOST/api/scenario/validate_invoice/run?scope=$SCOPE" \
  -H 'Content-Type: application/json' \
  -d '{"sender_email":"ahmed.r@gmail.com","vendor_id":"v_xyz","amount":340000}')
assert_eq "scenario2 ok"                       "true" "$(echo "$RES2" | jq -r .ok)"
assert_contains "suggests block_vendor"        "$(echo "$RES2" | jq -r '.suggestions[].action')" "block_vendor"
assert_contains "suggests notify_finance"      "$(echo "$RES2" | jq -r '.suggestions[].action')" "notify_finance"

section "6. Scenario: validate_invoice — trusted path"
RES3=$(curl -s -X POST "$HOST/api/scenario/validate_invoice/run?scope=$SCOPE" \
  -H 'Content-Type: application/json' \
  -d '{"sender_email":"ar@vendorxyz.com","vendor_id":"v_xyz","amount":24500}')
assert_contains "trusted: suggests match_to_po"      "$(echo "$RES3" | jq -r '.suggestions[].action')" "match_to_po"
assert_contains "trusted: suggests schedule_payment" "$(echo "$RES3" | jq -r '.suggestions[].action')" "schedule_payment"

section "7. Scenario: schedule_check (PM → slips → mitigations)"
RES4=$(curl -s -X POST "$HOST/api/scenario/schedule_check/run?scope=$SCOPE" \
  -H 'Content-Type: application/json' -d '{}')
assert_eq "scenario4 ok"                       "true" "$(echo "$RES4" | jq -r .ok)"
SLIP_COUNT=$(echo "$RES4" | jq -r '.steps[1].output.count')
[[ "$SLIP_COUNT" -ge 1 ]] && ok "schedule_check found ≥1 slip" || bad "schedule_check found ≥1 slip" "got $SLIP_COUNT"
SUGS=$(echo "$RES4" | jq -r '.suggestions[].action')
echo "$SUGS" | grep -qE "(escalate_slip|resequence_workstream|ask_for_catchup_plan)" \
  && ok "scenario suggests a slip mitigation" \
  || bad "scenario suggests a slip mitigation" "got actions: $SUGS"

section "8. Scenarios emit messages into the store"
MSG_COUNT=$(curl -s "$HOST/api/messages?project=$SCOPE" \
  | jq -r '[.messages[]? | select(.type=="system" and (.metadata.tags // [] | any(. == "scenario")))] | length')
[[ "$MSG_COUNT" -ge 6 ]] && ok "≥6 scenario messages in store (got $MSG_COUNT)" \
  || bad "≥6 scenario messages in store" "got $MSG_COUNT — scenarios may not be emitting through the hub"

section "9. Scenario index reports all 6"
INDEX=$(curl -s "$HOST/api/scenarios" | jq -r '.scenarios[].name' | sort | tr '\n' ',')
assert_eq "scenarios index lists all 6" \
  "morning_briefing,process_applicants,process_inbox,route_rfi,schedule_check,validate_invoice," \
  "$INDEX"

section "10. Inbox capability"
INBOX=$(curl -s "$HOST/api/cap/email/inbox")
ICOUNT=$(echo "$INBOX" | jq -r .count)
[[ "$ICOUNT" -gt 0 ]] && ok "inbox returns >0 emails (got $ICOUNT)" \
  || bad "inbox returns >0 emails" "got $ICOUNT — run scripts/scenarios/load-epc-atlas.sh first"
SUM=$(curl -s "$HOST/api/cap/email/summary")
assert_contains "summary has by_category" "$SUM" '"by_category"'
assert_contains "summary tracks impersonations" "$SUM" '"impersonations"'

section "11. Scenario: process_inbox (multi-agent triage)"
PI=$(curl -s -X POST "$HOST/api/scenario/process_inbox/run?scope=$SCOPE" \
  -H 'Content-Type: application/json' -d '{"max":8}')
assert_eq "process_inbox ok"      "true" "$(echo "$PI" | jq -r .ok)"
STEPCT=$(echo "$PI" | jq -r '.steps|length')
[[ "$STEPCT" -ge 2 ]] && ok "process_inbox emits ≥2 steps (got $STEPCT)" \
  || bad "process_inbox emits ≥2 steps" "got $STEPCT"

section "12. Scenario: morning_briefing (composite)"
MB=$(curl -s -X POST "$HOST/api/scenario/morning_briefing/run?scope=$SCOPE" \
  -H 'Content-Type: application/json' -d '{}')
assert_eq "morning_briefing ok"   "true" "$(echo "$MB" | jq -r .ok)"
MBSTEPS=$(echo "$MB" | jq -r '.steps|length')
[[ "$MBSTEPS" -ge 4 ]] && ok "morning_briefing has ≥4 agent reports (got $MBSTEPS)" \
  || bad "morning_briefing has ≥4 agent reports" "got $MBSTEPS"
assert_contains "morning_briefing engages PM" "$(echo "$MB" | jq -r '.steps[].agent')" "project_manager"
assert_contains "morning_briefing engages HR" "$(echo "$MB" | jq -r '.steps[].agent')" "hr"
assert_contains "morning_briefing engages Procurement" "$(echo "$MB" | jq -r '.steps[].agent')" "procurement"

section "13a. Inactive contract overrides recommendation (vendor_beta fix)"
# Seed an inactive-contract vendor + a trusted sender on its allowlist.
curl -s -X POST "$HOST/api/cap/procurement/vendors?scope=$SCOPE" -H 'Content-Type: application/json' \
  -d '{"id":"v_inactive","name":"BetaElectrics","domain":"betaelectrics.com","trusted_emails":["finance@betaelectrics.com"],"contract_active":false}' >/dev/null
INACTIVE=$(curl -s -X POST "$HOST/api/cap/procurement/validate-invoice?scope=$SCOPE" \
  -H 'Content-Type: application/json' \
  -d '{"sender_email":"finance@betaelectrics.com","vendor_id":"v_inactive","amount":15000}')
assert_eq "inactive contract: trusted still true"  "true" "$(echo "$INACTIVE" | jq -r .trusted)"
assert_eq "inactive contract: recommendation paused" \
  "pause_payment_until_contract_renewed" \
  "$(echo "$INACTIVE" | jq -r .recommendation)"
assert_contains "inactive contract: reasons mention contract" "$INACTIVE" "contract not active"

section "13b. Scenario: validate_invoice surfaces pause_payment for inactive contract"
INACTIVE_SCEN=$(curl -s -X POST "$HOST/api/scenario/validate_invoice/run?scope=$SCOPE" \
  -H 'Content-Type: application/json' \
  -d '{"sender_email":"finance@betaelectrics.com","vendor_id":"v_inactive","amount":15000}')
SUGS=$(echo "$INACTIVE_SCEN" | jq -r '.suggestions[].action')
assert_contains "inactive scenario: suggests pause_payment"      "$SUGS" "pause_payment"
assert_contains "inactive scenario: suggests reactivate_contract" "$SUGS" "reactivate_contract"

section "14. Scenario: route_rfi (Conductor → Site Engineer → HR pipeline)"
RFI=$(curl -s -X POST "$HOST/api/scenario/route_rfi/run?scope=$SCOPE" \
  -H 'Content-Type: application/json' \
  -d '{"rfi_id":"RFI-18","subject":"Column tolerances on grid C-7","body":"Per drawing A-102 v3, max plumb deviation on cast-in-place columns is ±10mm but spec calls ±6mm."}')
assert_eq "route_rfi ok"                          "true"        "$(echo "$RFI" | jq -r .ok)"
assert_eq "route_rfi step 1 is conductor"         "conductor"   "$(echo "$RFI" | jq -r '.steps[0].agent')"
assert_eq "route_rfi step 2 is site_engineer"     "site_engineer" "$(echo "$RFI" | jq -r '.steps[1].agent')"
assert_eq "route_rfi step 3 is hr"                "hr"          "$(echo "$RFI" | jq -r '.steps[2].agent')"
assert_eq "route_rfi classified as structural"    "structural"  "$(echo "$RFI" | jq -r '.steps[0].output.discipline')"
assert_contains "route_rfi suggests assign_to_site_engineer" \
  "$(echo "$RFI" | jq -r '.suggestions[].action')" "assign_to_site_engineer"

section "15. Conversation lifecycle (named threads + search)"
LIST=$(curl -s "$HOST/api/conversations")
assert_contains "GET /api/conversations returns object" "$LIST" '"conversations"'
COUNT=$(echo "$LIST" | jq -r '.conversations | length')
[[ "$COUNT" -ge 1 ]] && ok "lists ≥1 conversation (got $COUNT)" \
  || bad "lists ≥1 conversation" "got $COUNT — has /api/chat been used at least once?"

# auto-derived name is human, not snake_case
FIRST_NAME=$(echo "$LIST" | jq -r '.conversations[0].name')
if [[ "$FIRST_NAME" != *_* ]] || [[ ${#FIRST_NAME} -gt 0 && "${FIRST_NAME:0:1}" =~ [[:upper:]] ]]; then
  ok "auto-name is human-readable: \"$FIRST_NAME\""
else
  bad "auto-name is human-readable" "got: $FIRST_NAME"
fi

# create
CREATED=$(curl -s -X POST "$HOST/api/conversations")
NEW_ID=$(echo "$CREATED" | jq -r .id)
[[ "$NEW_ID" == conv_* ]] && ok "POST creates new id ($NEW_ID)" \
  || bad "POST creates new id" "got: $NEW_ID"

# get one
ONE=$(curl -s "$HOST/api/conversations/default")
assert_contains "GET /api/conversations/{id}" "$ONE" '"messages"'

# search across stored conversations — also asserts mode is reported
SEARCH=$(curl -s "$HOST/api/conversations/search?q=impersonation")
SHITS=$(echo "$SEARCH" | jq -r '.hits | length')
[[ "$SHITS" -gt 0 ]] && ok "search 'impersonation' returns hits (got $SHITS)" \
  || bad "search 'impersonation' returns hits" "got $SHITS"
SMODE=$(echo "$SEARCH" | jq -r .mode)
if [[ "$SMODE" == "semantic" || "$SMODE" == "lexical" ]]; then
  ok "search reports mode (got $SMODE)"
else
  bad "search reports mode" "got: $SMODE"
fi

section "16. Mock mail injection (Lab → inbox)"
BEFORE=$(curl -s "$HOST/api/cap/email/inbox" | jq -r .count)
MOCK_ID="mock_test_$(date +%s)"
RESP=$(curl -s -X POST "$HOST/api/cap/email/inbox" \
  -H 'Content-Type: application/json' \
  -d "{\"id\":\"$MOCK_ID\",\"from\":\"alpha@vendorxyz.com\",\"from_name\":\"Alpha Vendor\",\"subject\":\"Test invoice — mock\",\"body\":\"Test body\",\"category\":\"vendor\",\"trust_status\":\"trusted\",\"vendor_id\":\"vendor_xyz\"}")
RETURNED_ID=$(echo "$RESP" | jq -r .id)
assert_eq "POST returns same id"   "$MOCK_ID" "$RETURNED_ID"
AFTER=$(curl -s "$HOST/api/cap/email/inbox" | jq -r .count)
[[ "$AFTER" == "$((BEFORE+1))" ]] && ok "inbox count increased by 1 ($BEFORE → $AFTER)" \
  || bad "inbox count increased by 1" "before=$BEFORE after=$AFTER"
# Confirm the file is on disk and contains the expected fields.
FILE_PATH="projects/inbox/$MOCK_ID.json"
if [[ -f "$FILE_PATH" ]]; then
  ok "file written: $FILE_PATH"
  assert_contains "file has the subject" "$(cat "$FILE_PATH")" "Test invoice — mock"
else
  bad "file written" "$FILE_PATH not found"
fi
# POST without required field is rejected (400)
BADRESP=$(curl -s -o /tmp/test.out -w "%{http_code}" -X POST "$HOST/api/cap/email/inbox" \
  -H 'Content-Type: application/json' -d '{"subject":"no from"}')
assert_eq "POST without from is 400" "400" "$BADRESP"
# Cleanup so the file doesn't accumulate over many test runs.
rm -f "$FILE_PATH"

section "17. /api/agents includes the EPC roster"
ROLES=$(curl -s "$HOST/api/agents" | jq -r '.agents[].role' | sort | tr '\n' ',')
for r in hr procurement project_manager site_engineer hse qa_inspector; do
  if echo "$ROLES" | grep -q "$r,"; then ok "/api/agents includes $r"
  else bad "/api/agents includes $r" "ROLES=$ROLES"; fi
done

# ── summary ──────────────────────────────────────────────────────────
echo
echo "──────────────────────────────────────────────"
TOTAL=$((PASS+FAIL))
if [[ "$FAIL" == "0" ]]; then
  printf "\033[32m==> PASS\033[0m  %d/%d\n" "$PASS" "$TOTAL"
  exit 0
else
  printf "\033[31m==> FAIL\033[0m  %d/%d passed; %d failure(s):\n" "$PASS" "$TOTAL" "$FAIL"
  for f in "${FAILS[@]}"; do printf "    - %s\n" "$f"; done
  exit 1
fi
