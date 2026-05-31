#!/usr/bin/env bash
# Small software-development demo: a single-file HTML calculator built by
# the IT pack. Submits one task via /api/task — the orchestrator picks it
# up in the background and the agents stream work into /conductor in real
# time via the same WS path the proactive briefing uses.
set -euo pipefail

HOST="${HOST:-http://localhost:8099}"
PROJECT_ID="${PROJECT_ID:-calculator}"

echo "→ submitting calculator build to $HOST (project=$PROJECT_ID)"

response=$(curl -s -X POST "$HOST/api/task" \
  -H 'Content-Type: application/json' \
  -d "$(jq -n --arg p "$PROJECT_ID" '{
    task: "Build a single-file HTML calculator with a refined dark theme. Pure HTML/CSS/vanilla JS, no frameworks. Buttons for 0-9, decimal, +, -, ×, ÷, =, AC, ±, %. Keyboard support. Display large; history of last 5 operations under the display. One file at index.html, accessible, mobile-friendly.",
    project_id: $p
  }')")

echo "$response" | jq . 2>/dev/null || echo "$response"

echo
echo "→ adding an hourly status-check cron for $PROJECT_ID"
curl -s -X POST "$HOST/api/cron/" \
  -H 'Content-Type: application/json' \
  -d "$(jq -n --arg p "$PROJECT_ID" '{
    task: "Hourly status check on the calculator project — what shipped this hour, what is in review, what is blocked",
    schedule: "1h",
    agent_role: "ceo",
    project_id: $p
  }')" | jq '.job | {id, schedule, task:(.task[:60])}' 2>/dev/null || true

echo
echo "→ done. Open $HOST/conductor — agent messages from the IT pack will"
echo "   stream into the conversation as they work. Files land under"
echo "   projects/$PROJECT_ID/."
