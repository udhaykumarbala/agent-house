# Product Specification: WorkTime Tracker

Based on research in: `.plans/research/product-research.md`, `ux-research.md`, `ui-research.md`, `architecture-research.md`, `security.md`

---

## Target User

**"Small Team Manager Sam"** — Age 28-45, manages 3-20 people. Small business owner, freelance team lead, or shift supervisor. Not deeply technical but comfortable with web apps. Currently tracks hours in spreadsheets or informal messages. Frustrated by per-user SaaS pricing and complex onboarding in tools like Toggl/Clockify. Works primarily on desktop, occasionally checks reports on tablet/mobile.

## Unique Value Proposition

**"Free, instant, local-first time tracking for small teams — no accounts, no subscriptions, no data leaving your browser."**

We sit below Toggl/Clockify in complexity, above spreadsheets in capability. Zero setup friction: open the page, add employees, start tracking.

---

## Features (Prioritized)

### P0 — Must Have (MVP)

- [ ] **F1: Employee Management** — Add, edit, deactivate employees with name, role, and department. The manager is the sole user; employees don't log in. Color-coded per employee for visual identification.
- [ ] **F2: One-Click Clock In/Out** — Prominent clock-in/out button that is the most visible element on screen. Select employee → optionally select project → click. Timer runs visibly. Clock out with one click. Auto-calculates duration on clock-out.
- [ ] **F3: Manual Time Entry** — Enter start time, end time, employee, and project directly without using the timer. For logging hours after the fact (e.g., yesterday's forgotten entry).
- [ ] **F4: Project & Task Management** — Create projects with color tags. Optionally add tasks within projects. Assign time entries to a project (and optionally a task). Support active/completed/archived status.
- [ ] **F5: Dashboard (Landing Page)** — Today's overview: who's clocked in (with status indicators), total hours tracked today, active projects, running timers. Quick clock-in action directly from the dashboard.
- [ ] **F6: Report Generation** — Daily, weekly, and monthly summary reports. Filter by employee, project, and date range. Show totals, averages, and breakdowns. Date presets: Today, This Week, This Month, Custom.
- [ ] **F7: CSV Export** — One-click export of any report to CSV. Properly sanitized against formula injection. UTF-8 BOM for Excel compatibility.
- [ ] **F8: Data Persistence (localStorage)** — All data stored locally in the browser. Prefixed keys (`wtt_*`) to avoid collisions. In-memory caching for fast reads with durable writes.
- [ ] **F9: Data Export/Import (Backup)** — Full JSON export of all data for backup. JSON import to restore data or migrate between machines. Clear All Data option in settings.

### P1 — Should Have

- [ ] **F10: Timesheet Grid View** — Weekly grid (rows = employees or projects, columns = days) for reviewing and editing entries in bulk. Running totals per row and column. Inline editing.
- [ ] **F11: Visual Report Charts** — CSS-based bar charts showing hours by employee or project. Summary cards (total hours, avg per day, top project) at the top of the reports view.
- [ ] **F12: Print-Friendly Reports (PDF)** — Print-optimized CSS stylesheet (`@media print`). Manager clicks "Print" to get a clean PDF via the browser's native print dialog. Zero dependency.
- [ ] **F13: Employee Status Indicators** — Color-coded status: green (clocked in), gray (clocked out), yellow (on break). Visible in the employee list and dashboard.
- [ ] **F14: Input Validation & Security** — Max length on all inputs, character validation, XSS prevention via `textContent` rendering, overlapping time entry detection, clock-out > clock-in enforcement.

### P2 — Nice to Have (Post-MVP)

- [ ] **F15: Keyboard Shortcuts** — Power-user shortcuts: Space (clock in/out), N (new entry), E/P/R (navigate views).
- [ ] **F16: Dark Mode** — Deep indigo dark mode as primary, warm off-white light mode as secondary.
- [ ] **F17: Project Summary Report** — Dedicated report type showing total hours per project broken down by employee and task.
- [ ] **F18: Data Archival** — Archive old time entries (e.g., > 6 months) to reduce localStorage size. Archived data exportable but not loaded in memory.
- [ ] **F19: Timer Recovery** — Save active timer state to localStorage every 30 seconds. Recover running timers on page reload or accidental tab close.

---

## User Stories

1. **As a** team manager, **I want** to clock an employee in with one click **so that** I can start tracking their time without navigating away from the main screen.
2. **As a** team manager, **I want** to see who's currently clocked in on a dashboard **so that** I have an at-a-glance view of my team's status.
3. **As a** team manager, **I want** to generate a weekly report filtered by project **so that** I can review how many hours were spent on each project for client billing.
4. **As a** team manager, **I want** to export a monthly report as CSV **so that** I can send it to payroll or import it into a spreadsheet.
5. **As a** team manager, **I want** to add employees with a name, role, and department **so that** I can organize my team and filter reports.
6. **As a** team manager, **I want** to create projects with tasks **so that** time entries have context and reports can break down hours by project.
7. **As a** team manager, **I want** to manually enter time for yesterday **so that** I can record hours for entries that weren't tracked in real time.
8. **As a** team manager, **I want** to back up all my data as a JSON file **so that** I don't lose my tracking data if I clear my browser.
9. **As a** team manager, **I want** the timer to keep running even if I navigate between views **so that** I don't have to stay on one screen while an employee is working.
10. **As a** team manager, **I want** to see visual charts in my reports **so that** I can quickly understand time distribution without reading raw numbers.

---

## Acceptance Criteria

### F1: Employee Management
- [ ] Can add an employee with name (required), role, and department
- [ ] Can edit an existing employee's details
- [ ] Can deactivate an employee (they remain in historical data but can't be clocked in)
- [ ] Employee list shows all employees with their status and today's hours
- [ ] Each employee has an auto-assigned color for visual identification
- [ ] Empty state: "No employees yet — Add your first employee to get started"

### F2: One-Click Clock In/Out
- [ ] Clock-in button is the most prominent UI element on the dashboard
- [ ] Clocking in requires: selecting an employee (required) and project (defaults to "General" if none selected)
- [ ] Running timer displays elapsed time updating every second
- [ ] Clock-out is a single click — entry is auto-saved with calculated duration
- [ ] Active timer is visible from any view (persistent top bar)
- [ ] Cannot clock in the same employee twice simultaneously

### F3: Manual Time Entry
- [ ] Can create an entry with: employee, project, date, start time, end time
- [ ] Duration auto-calculates from start and end time
- [ ] Validates: end time > start time, max 24h per entry, no overlapping entries for same employee
- [ ] Can edit existing time entries inline
- [ ] Can delete a time entry with confirmation

### F4: Project & Task Management
- [ ] Can create a project with name (required), description, and color
- [ ] Can add tasks to a project
- [ ] Can mark projects as active, completed, or archived
- [ ] Archived projects don't appear in the clock-in project selector
- [ ] A default "General" project exists for unassigned time
- [ ] Empty state: "No projects yet — Create a project to categorize time entries"

### F5: Dashboard
- [ ] Shows total hours tracked today across all employees
- [ ] Shows list of currently clocked-in employees with running timers
- [ ] Shows count of active employees and active projects
- [ ] Quick clock-in action available directly on dashboard
- [ ] Displays recent time entries (last 5-10)

### F6: Report Generation
- [ ] Daily report: hours per employee for a specific date, broken down by project
- [ ] Weekly report: 7-day summary with hours per employee per day and totals
- [ ] Monthly report: totals per employee and per project for the selected month
- [ ] Filter by: date range, specific employee(s), specific project(s)
- [ ] Date presets: Today, This Week, This Month, Last Week, Last Month, Custom
- [ ] Shows both summary totals and detailed entry list

### F7: CSV Export
- [ ] Export button is prominently placed on the reports page (top-right)
- [ ] Exports currently visible/filtered report data
- [ ] CSV includes: employee name, project, task, date, clock-in, clock-out, duration, notes
- [ ] All cell values are wrapped in double quotes
- [ ] Formula-dangerous characters (`=`, `+`, `-`, `@`) are prefixed with a single quote
- [ ] File includes UTF-8 BOM for proper Excel encoding
- [ ] Download triggers immediately with a descriptive filename (e.g., `worktime-weekly-2026-03-16.csv`)

### F8: Data Persistence
- [ ] All data persists across page reloads and browser sessions
- [ ] localStorage keys use `wtt_` prefix
- [ ] Data reads are cached in memory for performance
- [ ] Graceful handling of corrupted JSON (log warning, don't crash)
- [ ] Graceful handling of storage quota exceeded (show user-friendly error)

### F9: Data Export/Import
- [ ] "Export All Data" button in settings exports complete JSON backup
- [ ] "Import Data" button in settings accepts a JSON file and restores data
- [ ] Import validates data schema before applying
- [ ] "Clear All Data" button with confirmation dialog
- [ ] Import warns if it will overwrite existing data

---

## Success Metrics

- **Time to First Clock-In**: A new user should go from opening the app to clocking in their first employee in under 2 minutes (add employee + clock in)
- **Clock-In Speed**: Clocking in a returning employee should take under 5 seconds (2 clicks max)
- **Report Generation Speed**: Generating and exporting a weekly report should take under 10 seconds
- **Data Reliability**: Zero data loss from normal browser usage (page reload, tab close with no active timer, navigation)
- **Rendering Performance**: Dashboard and reports render in under 500ms for datasets up to 10,000 time entries

---

## Out of Scope (for MVP)

- User authentication / multi-user login (manager is sole operator)
- Server-side storage or cloud sync
- Mobile native app (responsive web only)
- GPS or location-based tracking
- Screenshot or activity monitoring
- Payroll integration or invoicing
- Email notifications or reminders
- Real-time collaboration (multiple managers editing simultaneously)
- Hourly rate / cost calculations in reports
- API or third-party integrations
- Multi-language / i18n support

---

## Technical Constraints (from Architecture & Security Research)

- **Template**: `static-enhanced` (Vite + Tailwind) — no backend
- **Storage**: localStorage with JSON, `wtt_` prefixed keys
- **Routing**: Hash-based SPA (`#/dashboard`, `#/employees`, `#/projects`, `#/timesheet`, `#/reports`)
- **Security**: XSS prevention via `textContent` only, CSV formula injection sanitization, input validation on all fields, try/catch on all localStorage reads
- **Charts**: CSS-based (no Chart.js or heavy libraries)
- **Export**: CSV via vanilla JS Blob, PDF via `window.print()` with print CSS
- **Time**: UTC storage, local display via `Intl.DateTimeFormat`
- **Dependencies**: Minimal runtime dependencies; zero-dependency core

---
Status: READY_FOR_REVIEW
