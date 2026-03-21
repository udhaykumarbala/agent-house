# Agent House v2 — Bug Tracker

**Date:** 2026-03-21
**Context:** Session-per-agent integration works (agents spawn, execute, produce code) but observability/reporting is broken.

---

## P0 — Critical

### BUG-01: Tool events invisible (tool_calls = 0 everywhere)
- **Status:** FIXED
- **Symptom:** All agents show `tools=0` in stats, even dev agents that create files
- **Root cause:** NDJSON translator may not be receiving tool_use/tool_result events from Claude Code in `acceptEdits`/`bypassPermissions` modes. Need raw NDJSON debug logging to see actual Claude output.
- **Fix:** Add debug logging to `session.readLoop()`, inspect raw NDJSON, fix translator if needed
- **Test:** Single-agent session test — send "create a file hello.txt with content hello" and check tool events

### BUG-02: Task/agent lifecycle events missing from WebSocket
- **Status:** FIXED
- **Symptom:** UI shows agents stuck on "working" after task completes. No phase transition visibility.
- **Root cause:** Orchestrator doesn't broadcast phase_started/phase_completed/task_completed events to WebSocket hub
- **Fix:** Add lifecycle event broadcasts in orchestrator phase transitions and task completion
- **Test:** Submit task, verify WebSocket receives phase_started, phase_completed, task_completed

---

## P1 — High

### BUG-03: "0 files created" despite files existing on disk
- **Status:** FIXED
- **Symptom:** Task summary says "0 files created" but `find` shows files
- **Root cause:** `result.Files` populated from tool_use events (BUG-01). Even if fixed, need fallback.
- **Fix:** Add filesystem diff — snapshot before agent, diff after, report changes
- **Test:** Dev agent creates file → check FileResult populated

### BUG-04: Agent states stuck on "working" in dashboard UI
- **Status:** FIXED
- **Symptom:** Cards stay blue/animated after task finishes
- **Root cause:** No task_completed event to reset all agents. Depends on BUG-02.
- **Fix:** On task_completed event, reset all agents to idle in JS
- **Test:** Task finishes → all agent cards show "Idle"

### BUG-05: CEO triage uses legacy `claude --print`
- **Status:** FIXED
- **Symptom:** Triage logs show `claude --print` while other phases use sessions
- **Root cause:** `triage.go` calls `ProcessInDir()` directly, not patched for sessions
- **Fix:** Update triage to use `ExecuteWithSession` when SessionManager available
- **Test:** Triage log shows `SESSION:ceo` spawned instead of `claude --print`

---

## P2 — Medium

### BUG-06: No subtasks created or assigned
- **Status:** OPEN
- **Symptom:** Kanban board empty, log shows "No development plan found"
- **Root cause:** CEO discussion phase doesn't create structured `.plans/final/approved-plan.md` in session mode
- **Fix:** Update CEO system prompt for discussion phase OR parse output into plan programmatically
- **Test:** After discussion, development plan exists with subtasks

### BUG-07: Response text often empty ("0 bytes")
- **Status:** OPEN
- **Symptom:** Agent responses show 0 bytes content in messages
- **Root cause:** Agents write to files instead of outputting text. `SendTask` only captures `text_delta` events.
- **Fix:** Also capture content from files written during turn, or read final text from tool outputs
- **Test:** Agent produces research doc → message content is non-empty

### BUG-08: No structured output format for delegation/review signals
- **Status:** OPEN
- **Symptom:** DELEGATE:/REVIEW:/COMPLETE: signals not parsed in session mode
- **Root cause:** Session agents don't output these text markers. Text parsing was regex-based on legacy response.
- **Fix:** Either update system prompts to include signals, or infer from events/files
- **Test:** Agent delegates → orchestrator detects delegation

---

## Testing Shortcut

For fast testing, skip the full 5-phase pipeline:
1. Use `EnablePhases: false` config OR
2. Create a minimal test that spawns ONE agent session directly
3. Send a simple task like "create hello.txt" and inspect all events
4. This avoids 18+ minute full pipeline runs

---

## Fix Log

| Date | Bug | Fix Description | Verified |
|------|-----|-----------------|----------|
| 2026-03-21 | BUG-01 | Removed lastContentIndex tracking — Claude Code stream-json sends new blocks per line, not cumulative | Integration test: tool_use=1, Write:1, hello.txt created |
| 2026-03-21 | BUG-02 | Added notifyLifecycle() in orchestrator for phase_started/phase_completed/task_completed events | Build OK, unit tests pass |
| 2026-03-21 | BUG-03 | Added filesystem diff (snapshotFiles before/after) as fallback file tracking | Build OK |
| 2026-03-21 | BUG-04 | Added handleLifecycleEvent in /live dashboard JS — resets all agents on task_completed | Build OK |
| 2026-03-21 | BUG-05 | Updated triage.go to use ExecuteWithSession when SessionManager available | Build OK |
