#!/usr/bin/env bash
# Seed a small set of proactive cron jobs so the system starts speaking
# first on its own — daily ops briefing, hourly inbox sweep, daily slip
# check. Idempotent enough: the cron scheduler accepts duplicates but
# you can DELETE them via the /triggers UI or `DELETE /api/cron/?id=...`.
set -euo pipefail
HOST="${HOST:-http://localhost:8099}"

post() {
  local body="$1"
  local code
  code=$(curl -s -o /tmp/cron.out -w "%{http_code}" -X POST "$HOST/api/cron/" \
    -H 'Content-Type: application/json' \
    -d "$body")
  if [[ "$code" != "200" ]]; then
    echo "  ✗ POST failed ($code): $(cat /tmp/cron.out)"
    return 1
  fi
  echo "  ✓ $(jq -r '.job | "\(.id)  \(.schedule)  \(.task)"' < /tmp/cron.out 2>/dev/null || cat /tmp/cron.out)"
}

echo "→ seeding proactive defaults at $HOST"
post '{"task":"Hourly inbox sweep — route any new vendor, applicant, or RFI email to the right agent and surface anything requiring a decision","schedule":"1h","agent_role":"ceo","project_id":"default"}'
post '{"task":"Daily ops briefing — summarise pending decisions, trigger fires, schedule slips, and anomalies; post a one-page brief","schedule":"24h","agent_role":"ceo","project_id":"default"}'
post '{"task":"Daily schedule slip check — scan active milestones, flag anything trending late, propose mitigations","schedule":"24h","agent_role":"project_manager","project_id":"default"}'

echo
echo "→ current cron jobs:"
curl -s "$HOST/api/cron" | jq '.jobs[] | {id, schedule, task: (.task[:60]), agent: .agent_role, project: .project_id, enabled}' 2>/dev/null \
  || curl -s "$HOST/api/cron"
