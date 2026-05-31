# Demo scenarios

Self-contained bootstraps that load realistic state into a running Agent House,
so a cold-start system goes from blank to demoable in 30 seconds.

Each scenario is a single shell script — it talks to the live server over
the same `/api/*` endpoints the UI uses. Nothing is hard-coded into the
binary. Re-run safely; emails use deterministic IDs and overwrite cleanly.

## Prerequisites

1. Build & run the server:
   ```
   ./scripts/build-ui.sh
   go build -o /tmp/ah ./cmd/agent-house
   /tmp/ah --serve --port 8099 &
   ```
2. (Optional) Seed the universal proactive cron defaults:
   ```
   ./scripts/seed-conductor-defaults.sh
   ```

## Scenarios

### EPC — `load-epc-atlas.sh`

Atlas Construction · Site 02. Drops 5 emails into the global inbox
(`projects/inbox/`):

- **Vendor impersonation alert** — `ahmed.r@gmail.com` claiming to be XYZ Steel, asking for bank-detail change (classic BEC pattern). Flagged `trust_status: "impersonation"`.
- **RFI #18** — Atlas Consulting asking about column tolerance on grid C-7.
- **Vendor invoice** — `ar@vendorxyz.com` Invoice #ST-0847 for 150 t structural steel, \$340,000 (the legitimate counterpart to the impersonation above).
- **Client progress request** — Sarah Jones at Client ABC needs the Phase 2 update for the board.
- **Applicant** — Raj Kumar applying for Site Supervisor on the Highway Bridge Project.

Also adds two scoped cron jobs:

- `24h` schedule-slip check (Project Manager)
- `1h` inbox sweep (CEO)

```bash
./scripts/scenarios/load-epc-atlas.sh
```

Open `/conductor` afterwards — the briefing will reflect the new fires and
agent activity as the cron runs.

### Software dev — `load-dev-calculator.sh`

Small, real build: a single-file HTML calculator with dark theme, full
keyboard support, history of last 5 ops. Submits the task via
`POST /api/task` so the orchestrator picks it up in the background; the IT
pack agents stream their work into `/conductor` over `/ws`.

Adds one `1h` status-check cron for the project.

```bash
./scripts/scenarios/load-dev-calculator.sh
```

Files land under `projects/calculator/`. Watch the conversation in
`/conductor` — agent messages stream in-thread as each phase produces
output.

## Reset

Each scenario seeds a project under `projects/<id>/` and (for EPC) writes
emails to `projects/inbox/`. To reset:

```bash
rm -rf projects/atlas-site projects/calculator
# remove the scenario's cron jobs from the /triggers page,
# or via curl -X DELETE "http://localhost:8099/api/cron/?id=<id>"
# remove seeded inbox emails if desired:
rm -f projects/inbox/email_atlas_*.json
```

## Adding a new scenario

Three rules:

1. **One script, no state outside `projects/` and the running server.** A
   scenario is *load-and-go*; it never edits source or config.
2. **Use deterministic IDs** for any file you write (e.g.
   `email_<scenario>_<n>.json`) so re-runs overwrite cleanly.
3. **Talk to `/api/*`, not the database directly.** If the endpoint doesn't
   exist, that's a backend gap to fix, not a reason to side-load JSON onto
   disk.
