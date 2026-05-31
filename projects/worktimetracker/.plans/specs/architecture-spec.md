# Architecture Specification: WorkTime Tracker

**Author:** Architect Agent
**Date:** 2026-03-18
**Template:** `static-enhanced` (Vite + Tailwind CSS)
**Status:** AWAITING_REVIEW

---

## 1. System Overview

WorkTime Tracker is a zero-setup, local-first time tracking tool for small teams. It runs entirely in the browser with no backend — data persists in localStorage. A manager opens the page, adds employees, and starts tracking time immediately.

**Core capabilities:**
- Employee management (CRUD, status tracking)
- Clock in/out with timer + manual time entry
- Project and task assignment
- Report generation (daily, weekly, monthly) with CSV export
- Data backup via JSON export/import

---

## 2. Tech Stack

| Layer | Technology | Justification |
|-------|-----------|---------------|
| Build tool | Vite | Fast HMR, ES module bundling, tree-shaking, minimal config |
| Styling | Tailwind CSS v3 | Utility-first, purges unused CSS, consistent design tokens |
| Language | Vanilla JavaScript (ES modules) | Zero framework overhead, ~30KB estimated bundle, fast load |
| Storage | localStorage (JSON) | No server needed, instant persistence, sufficient for single-user scope |
| Routing | Hash-based SPA | No server config, works on static hosting |
| Charts | CSS-based bars + HTML tables | No library dependency, sufficient for summary visualizations |
| Export | Vanilla JS CSV + `window.print()` PDF | Zero-dependency, browser-native |
| Icons | Lucide (CDN or bundled) | Outlined style, 1.5px stroke, consistent with "Warm Precision" aesthetic |
| Fonts | Google Fonts — Manrope (headings), Inter (body/UI) | Per UI research: warmth + legibility, tabular nums for time |

---

## 3. Data Architecture

### 3.1 Core Entities

**Employee**
```
{
  id: string (crypto.randomUUID()),
  name: string (required, max 100 chars),
  role: string (e.g., "Developer", "Designer"),
  department: string,
  hourlyRate: number | null (optional, for cost reports),
  status: "active" | "inactive",
  color: string (hex, for UI identification),
  createdAt: string (ISO 8601 UTC)
}
```

**Project**
```
{
  id: string,
  name: string (required, max 100 chars),
  description: string,
  color: string (hex, for visual tagging),
  status: "active" | "completed" | "archived",
  createdAt: string (ISO 8601 UTC)
}
```

**Task**
```
{
  id: string,
  projectId: string (ref → Project),
  name: string (required),
  status: "todo" | "in_progress" | "done"
}
```

**TimeEntry**
```
{
  id: string,
  employeeId: string (ref → Employee),
  projectId: string (ref → Project),
  taskId: string | null (ref → Task, optional),
  date: string (YYYY-MM-DD),
  clockIn: string (ISO 8601 UTC),
  clockOut: string | null (null = currently clocked in),
  duration: number | null (minutes, computed on clock-out),
  notes: string,
  type: "clock" | "manual"
}
```

### 3.2 Entity Relationships

```
Employee 1───∞ TimeEntry
Project  1───∞ TimeEntry
Project  1───∞ Task
Task     1───∞ TimeEntry (optional)
```

### 3.3 localStorage Strategy

- **Keys:** `wtt_employees`, `wtt_projects`, `wtt_tasks`, `wtt_time_entries`, `wtt_settings`
- **Format:** JSON arrays per key
- **Prefix:** `wtt_` to avoid collisions
- **Capacity:** ~5-10MB limit. A year of data for 50 employees at 2 entries/day ≈ 7MB. Warn at 80%.
- **Caching:** Read localStorage once on app init, cache in memory. Write on every mutation.
- **Timer recovery:** Save active timer state to `wtt_active_timer` every 30 seconds. Recover on reload.

---

## 4. Application Structure

```
src/
├── main.js                  # App entry, router init, global event bus
├── router.js                # Hash-based router (#/dashboard, #/employees, etc.)
├── styles/
│   └── main.css             # Tailwind imports + custom component styles + print CSS
├── data/
│   ├── store.js             # localStorage wrapper: read/write/cache/error handling
│   ├── employeeService.js   # Employee CRUD + validation
│   ├── projectService.js    # Project/Task CRUD + validation
│   └── timeEntryService.js  # Time entry CRUD + timer logic + overlap validation
├── views/
│   ├── dashboard.js         # Landing page: today's overview, active timers, quick clock-in
│   ├── employees.js         # Employee list + CRUD modal
│   ├── projects.js          # Project & task management
│   ├── timesheet.js         # Time entry list, manual entry, inline editing
│   └── reports.js           # Filtered reports, charts, export
├── components/
│   ├── modal.js             # Reusable modal for create/edit forms
│   ├── dataTable.js         # Sortable, filterable table
│   ├── timerWidget.js       # Persistent clock-in/out timer (top bar)
│   ├── dateRangePicker.js   # Date range selector with presets
│   ├── toast.js             # Success/error notification toasts
│   └── sidebar.js           # Left sidebar navigation
└── utils/
    ├── dates.js             # Date formatting, duration calc, locale display
    ├── export.js            # CSV export with formula injection protection
    ├── sanitize.js          # escapeHTML(), input validation helpers
    └── eventBus.js          # Pub/sub for cross-view communication
```

### 4.1 Architecture Patterns

- **View switching:** Hash-based router renders one view at a time. Only active view is in DOM.
- **State management:** Pub/sub event bus. Services emit events on data changes; views subscribe and re-render.
- **Data layer:** Service classes encapsulate localStorage access, validation, and business logic. Views never touch localStorage directly.
- **Component reuse:** Modal, DataTable, Toast, DateRangePicker are generic and parameterized.

---

## 5. View Specifications

### 5.1 Dashboard (`#/dashboard`) — Landing Page

**Purpose:** At-a-glance command center for today.

**Layout:**
- **Summary cards row:** Total hours today | Employees clocked in | Active projects | Entries today
- **Active timers section:** List of employees currently clocked in with running timers, project, and elapsed time
- **Quick clock-in:** Prominent button/dropdown to clock in any employee with project selection
- **Recent activity:** Last 10 time entries (mini-timeline)

**Empty state:** "Welcome to WorkTime Tracker! Add your first employee to get started." → CTA button to Employees view.

### 5.2 Employees (`#/employees`)

**Purpose:** Manage the team.

**Layout:**
- **Employee grid/list:** Name, role, department, status indicator (green dot = active, gray = inactive), today's hours
- **Add Employee button** (top right) → opens modal
- **Edit/deactivate** via row actions
- **Search/filter** by name, department, status

**Employee form fields:** Name (required), Role, Department, Hourly Rate (optional), Status, Color (auto-assigned or picker)

### 5.3 Projects (`#/projects`)

**Purpose:** Manage projects and their tasks.

**Layout:**
- **Project list** with color tags, status badges, task count, total hours logged
- **Add Project button** → modal
- **Expand project** → shows tasks underneath
- **Add Task** inline within expanded project
- **Filter** by status (active/completed/archived)

### 5.4 Time Tracker / Timesheet (`#/timesheet`)

**Purpose:** View, create, and edit time entries.

**Layout:**
- **View toggle:** List view (default) | Weekly grid view
- **List view:** Sortable table of entries — Date, Employee, Project, Task, Clock In, Clock Out, Duration, Notes, Actions (edit/delete)
- **Weekly grid view:** Rows = employees, Columns = Mon-Sun, Cells = total hours (click to see entries)
- **Add Entry** button → modal with: Employee (required), Project (required), Task (optional), Date, Start Time, End Time, Notes
- **Inline editing** for quick corrections
- **Filter bar:** Date range, Employee, Project

### 5.5 Reports (`#/reports`)

**Purpose:** Generate filtered summaries with export.

**Layout:**
- **Filter bar (top):** Report type (Daily/Weekly/Monthly/Project), Date range (presets + custom picker), Employee filter, Project filter
- **Summary cards:** Total hours, Average hours/day, Top employee, Top project
- **Chart area:** CSS-based horizontal bar chart (hours by employee or by project)
- **Data table:** Detailed breakdown below chart, sortable
- **Export buttons (top right):** "Export CSV" and "Print/PDF" — prominent, single-click

**Report types:**
1. **Daily:** Hours per employee for a specific date, by project
2. **Weekly:** 7-day grid, hours per employee per day, with totals
3. **Monthly:** Totals per employee/project for the month
4. **By Project:** Total hours per project, broken down by employee and task

---

## 6. Navigation & Layout

**Global layout:**
- **Left sidebar** (collapsible): Dashboard, Employees, Projects, Timesheet, Reports icons + labels
- **Top bar:** App title/logo left, active timer widget center/right, settings gear right
- **Content area:** Full remaining width, scrollable

**Active timer banner:** When any employee is clocked in, a persistent indicator in the top bar shows: employee name, project, elapsed time, Stop button. Visible from all views.

---

## 7. Security Requirements

Per security research, these are non-negotiable:

### 7.1 XSS Prevention (Critical)
- **All user data** rendered via `textContent`, never `innerHTML`
- Centralized `escapeHTML()` utility in `utils/sanitize.js` for any template literal rendering
- No `eval()`, `Function()`, or `document.write()` with user data

### 7.2 CSV Export Sanitization (Critical)
- All cell values wrapped in double quotes
- Values starting with `=`, `+`, `-`, `@`, `\t`, `\r` prefixed with single quote
- UTF-8 BOM (`\uFEFF`) for Excel compatibility
- Double quotes within values escaped as `""`

### 7.3 Input Validation (High)
- Employee names: required, max 100 chars
- Time entries: clockOut > clockIn, no overlapping entries per employee, max 24h per entry
- Project/task names: required, max 100 chars
- All dates: valid ISO format
- Hourly rate: non-negative number

### 7.4 localStorage Resilience (High)
- All reads wrapped in try/catch for corrupted JSON
- Schema validation on load — reject malformed entries
- QuotaExceededError handling with user-visible warning
- "Clear All Data" button in settings

---

## 8. Design System

Per UI research — "Warm Precision" style.

### 8.1 Color Palette

| Token | Value | Usage |
|-------|-------|-------|
| `--color-primary` | `#4338CA` (Indigo 700) | Primary buttons, active nav, links |
| `--color-primary-dark` | `#2E1065` (Indigo 950) | Sidebar background, headers |
| `--color-accent` | `#D97706` (Amber 600) | CTAs, active timer, highlights |
| `--color-accent-light` | `#FEF3C7` (Amber 100) | Accent backgrounds, badges |
| `--color-surface` | `#FFFBF5` | Page background (light mode) |
| `--color-card` | `#FFFFFF` | Card/panel backgrounds |
| `--color-text` | `#1E1B4B` (Indigo 950) | Primary text |
| `--color-text-secondary` | `#6B7280` (Gray 500) | Secondary text, labels |
| `--color-border` | `rgba(99, 102, 241, 0.1)` | Subtle borders |
| `--color-success` | `#059669` (Emerald 600) | Clocked in, success states |
| `--color-warning` | `#D97706` (Amber 600) | On-break, warnings |
| `--color-danger` | `#DC2626` (Red 600) | Errors, delete actions |

### 8.2 Typography

- **Headings:** Manrope (Google Fonts), weights 600-700
- **Body/UI:** Inter (Google Fonts), weights 400-500-600
- **Timer display:** Inter with `font-variant-numeric: tabular-nums`
- **Scale:** 12px (small labels), 14px (body), 16px (large body), 20px (h3), 24px (h2), 30px (h1)

### 8.3 Spacing & Shape

- **Border radius:** 8px default, 12px for cards/modals
- **Shadows:** Soft, warm-tinted (e.g., `0 4px 12px rgba(67, 56, 202, 0.08)`)
- **Animations:** 200ms ease transitions for state changes
- **Responsive breakpoints:** 640px (mobile), 768px (tablet), 1024px (desktop)

---

## 9. Key Technical Decisions

| Decision | Choice | Trade-off |
|----------|--------|-----------|
| No framework | Vanilla JS + ES modules | More manual DOM work, but zero runtime overhead, ~30KB bundle |
| localStorage over IndexedDB | Simpler API, synchronous | 5-10MB limit, no queries — acceptable for scope |
| CSS charts over Chart.js | No dependency | Less interactivity — horizontal bars and tables are sufficient |
| Hash routing over History API | No server config needed | URLs have `#` — acceptable for local-first tool |
| Print CSS over PDF library | Zero dependency | Less control over PDF layout — acceptable via print styling |
| Pub/sub over reactive framework | Lightweight, explicit | Manual subscription management — manageable for 5 views |

---

## 10. Performance Targets

- **First paint:** < 500ms (static assets, no API calls)
- **Bundle size:** < 50KB gzipped (JS + CSS)
- **localStorage read:** Once on init, cached in memory
- **Timer update:** `setInterval` at 1s for running display
- **Debounced search/filter:** 300ms
- **Render strategy:** Only active view in DOM

---

## 11. Data Backup & Migration

- **Export All Data:** JSON download of entire localStorage state (all collections)
- **Import Data:** Upload JSON file, validate schema, merge or replace
- **CSV export per report:** Individual report data
- **Storage warning:** Alert when localStorage usage exceeds 80% of estimated limit

---

## 12. Keyboard Shortcuts (Power User)

| Key | Action |
|-----|--------|
| `Space` | Quick clock in/out on dashboard |
| `N` | New time entry |
| `E` | Go to Employees |
| `P` | Go to Projects |
| `R` | Go to Reports |
| `Esc` | Close modal |

---

Status: AWAITING_REVIEW
