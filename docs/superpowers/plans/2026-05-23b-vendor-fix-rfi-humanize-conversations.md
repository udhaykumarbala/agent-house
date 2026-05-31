# Vendor fix + RFI scenario + humanize + conversation lifecycle

> Plan: ship 4 things this turn, each verified before moving to the next.

## Phase A — Quick fixes
1. **Procurement**: when sender is trusted but `contract_active=false`, override `recommendation` so payment isn't blindly approved.
2. **Humanize lib (frontend)**: `humanize("snake_case_id") → "Snake Case ID"` with acronyms (RFI, PO, HSE, QA, CV, MEP, BEC, KYC, CEO, PM, HR). Apply across ScenarioCard, ActionCard, workforce panel.

## Phase B — route_rfi scenario
3. `internal/scenario/route_rfi.go` — given RFI text, Conductor classifies discipline (structural/electrical/civil/MEP/safety), Site Engineer assigned by default, HR matches pipeline applicants by skill keywords, suggests `route_to_site_engineer` / `escalate_to_specialist` / `hire_for_gap`.
4. Test coverage in `scripts/test/test-scenarios.sh`.

## Phase C — Conversation lifecycle
5. **Backend** — extend `ConversationStore` with `ListConversations()` (read files in dir, build metadata from msgs), `Search(q)` (text search across all). Filename is `<conv_id>.json`.
6. **HTTP**:
   - `GET /api/conversations` — list with name/last_at/count
   - `GET /api/conversations/{id}` — full thread (alias to existing chat GET)
   - `POST /api/conversations` — create empty conversation, returns id (no chat yet)
   - `GET /api/conversations/search?q=...` — text hits across all
7. **Frontend**:
   - `?conv=<id>` URL param drives ConductorApp state (default to a stable `default`)
   - Auto-name from first user message (first 40 chars)
   - "New conversation" → `conv_<ms>` id, replaceState URL, clear msgs
   - "History" → command palette overlay: list + text search → click to switch

## Phase D — Verify
8. Unit tests for procurement fix + humanize.
9. e2e: route_rfi + conversation list/search added to test harness.
10. Smoke open `/conductor`, fire route_rfi from /lab, create new conversation, search history.
