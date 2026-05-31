# APPROVED PLAN: WorkTime Tracker

Reviewed by: CEO
Date: 2026-03-18
Status: APPROVED FOR DEVELOPMENT

## Summary

WorkTime Tracker is a zero-setup, local-first time tracking tool for small teams. A manager opens the app, adds employees, creates projects, and starts clocking time — all data stays in the browser via localStorage. It features real-time timers, manual time entry, project/task assignment, and filtered report generation with CSV export. The design uses a distinctive "Warm Precision" aesthetic (deep indigo + warm amber, dark-mode-first) to differentiate from every major competitor.

## Approved Specifications

- Product Spec: .plans/specs/product-spec.md
- UX Spec: .plans/specs/ux-spec.md
- UI Spec: .plans/specs/ui-spec.md
- Architecture Spec: .plans/specs/architecture-spec.md
- Security Spec: .plans/specs/security-spec.md
- Development Plan: .plans/development-plan.json

## Key Decisions

1. **UI Spec colors are canonical.** The UI spec's color palette (Primary #5746B2, Accent #CF8A2E, etc.) takes precedence over the architecture spec's section 8.1. Developers must use UI spec hex values.
2. **Dark mode is the default and only theme for MVP.** The entire UI spec is designed around the dark palette. A dark/light toggle is deferred to P2.
3. **Employee email is optional, not required.** Include it in the form, validate format if provided, but do not require it.
4. **Toast notifications render bottom-right**, auto-dismiss after 3s, per UX spec.
5. **Breakpoints follow UX spec**: Mobile <768px, Tablet 768-1023px, Desktop >=1024px.
6. **Settings is a modal/panel** from the sidebar gear icon, not a separate routed view. Five main routes only: dashboard, employees, projects, timesheet, reports.
7. **Date range presets must include**: Today, This Week, This Month, Last Week, Last Month, Custom — per product spec.
8. **Security utilities built first.** sanitize.js, validation.js, export.js, and store.js must exist before any view code is written.
9. **No framework.** Vanilla JS + ES modules. Zero runtime dependency for core functionality.
10. **Hash-based SPA routing** with only the active view in DOM.

## Implementation Order

1. **Phase 1: Foundation** — Vite scaffolding, design tokens, data layer (store + services), router, reusable components (modal, toast, table), utility modules (dates, sanitize, eventBus, export)
2. **Phase 2: Core Views** — Employee management, Project/Task management, Timer widget (persistent top bar), Dashboard
3. **Phase 3: Timesheet & Reports** — Time entry list with manual entry, date range picker, filtered reports with CSS charts, CSV export, print CSS
4. **Phase 4: Polish** — Settings panel (data backup/import/clear), keyboard shortcuts, responsive refinement, final transitions/animations

## Notes for Development Team

- **XSS prevention is non-negotiable.** Use `element.textContent` for ALL user data. Never `innerHTML` with unsanitized input. The `escapeHTML()` utility is a fallback, not the primary defense.
- **CSV export must sanitize for formula injection.** Prefix cells starting with `=`, `+`, `-`, `@`, `\t`, `\r` with a single quote. Wrap all values in double quotes. Include UTF-8 BOM.
- **Timer state must persist.** Save active timer to `wtt_active_timer` every 30 seconds. Recover on page reload. This is a P0 feature — users will lose trust if a tab close loses an hours-long timer.
- **All localStorage reads must be wrapped in try/catch.** Corrupted JSON must not crash the app.
- **Use `font-variant-numeric: tabular-nums`** on all timer displays and numeric columns so digits don't shift width as they update.
- **Empty states are first impressions.** Every view must have a guided empty state that directs the user to the next logical action.
- **Performance target:** Dashboard and reports must render in <500ms for up to 10,000 time entries. Read localStorage once on init, cache in memory.
- **The "General" project must auto-exist.** It's the fallback for unassigned time entries — create it on first app load if no projects exist.
- **Lucide icons** (outlined, 1.5px stroke) for all iconography. Consistent 20px size for navigation.

---
BEGIN DEVELOPMENT
