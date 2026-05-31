# Architecture Research: Work Time Tracker

## 1. Competitor Analysis

### Existing Time Tracking Tools

| Tool | Key Features | UX Approach | What Works | What Doesn't |
|------|-------------|-------------|------------|--------------|
| **Toggl Track** | One-click timer, project tagging, reports dashboard | Minimal, timer-centric, single prominent "Start" button | Instant clock-in UX, clean report charts, weekly calendar heat map | Overwhelming settings for simple use cases |
| **Clockify** | Manual & timer entry, team management, exportable reports | Tabbed layout (Timer / Timesheet / Reports) | Timesheet grid view for bulk entry, CSV/PDF export, project color coding | Cluttered navigation, too many features visible at once |
| **Harvest** | Timer + manual, invoicing, expense tracking | Dashboard with running timer always visible | Report filtering by person/project/date, visual weekly summaries | Focused on billing — overkill for internal tracking |
| **Hubstaff** | Auto-tracking, screenshots, GPS | Activity-focused dashboard | Per-employee activity breakdown, payroll integration | Privacy-invasive, heavy for simple tracking |
| **Timely** | AI-powered auto-tracking, memory timeline | Timeline-based daily view | Beautiful daily timeline visualization, automatic categorization | Complex, AI dependency not suitable for manual tracker |

### Key Takeaways from Competitors
1. **Timer-first UX** — The most successful tools make clock-in/out the most prominent action (Toggl's big play button)
2. **Timesheet grid** — Clockify's weekly timesheet grid is the most efficient for bulk data entry/review
3. **Color-coded projects** — Every major tool uses color tags for projects; instant visual identification
4. **Report filtering** — Date range + employee + project as the three primary filter axes
5. **Export simplicity** — CSV is universally supported; PDF for formal reports

---

## 2. Data Architecture

### Core Entities

```
Employee
├── id: string (uuid)
├── name: string
├── email: string
├── role: string (e.g., "Developer", "Designer")
├── department: string
├── hourlyRate: number (optional, for cost reports)
├── status: "active" | "inactive"
├── createdAt: ISO timestamp
└── color: string (hex, for UI identification)

Project
├── id: string (uuid)
├── name: string
├── description: string
├── color: string (hex, for visual tagging)
├── status: "active" | "completed" | "archived"
├── createdAt: ISO timestamp
└── tasks: Task[]

Task
├── id: string (uuid)
├── projectId: string (ref → Project)
├── name: string
├── description: string
└── status: "todo" | "in_progress" | "done"

TimeEntry
├── id: string (uuid)
├── employeeId: string (ref → Employee)
├── projectId: string (ref → Project)
├── taskId: string (ref → Task, optional)
├── date: ISO date string (YYYY-MM-DD)
├── clockIn: ISO timestamp
├── clockOut: ISO timestamp | null (null = currently clocked in)
├── duration: number (minutes, computed on clock-out)
├── notes: string
└── type: "clock" | "manual" (clock = timer-based, manual = manually entered)
```

### Entity Relationships
```
Employee 1──────∞ TimeEntry
Project  1──────∞ TimeEntry
Project  1──────∞ Task
Task     1──────∞ TimeEntry (optional link)
```

### localStorage Strategy
- **Key-per-collection**: `wtt_employees`, `wtt_projects`, `wtt_tasks`, `wtt_time_entries`
- **JSON arrays** stored per key
- **Prefix `wtt_`** to avoid collisions with other apps
- **Settings key**: `wtt_settings` for app config (default views, export preferences)
- **Size consideration**: localStorage limit is ~5-10MB. A year of data for a 50-person team with 2 entries/day ≈ 36,500 entries × ~200 bytes = ~7MB. For larger datasets, implement a data archival/export mechanism.

---

## 3. View Architecture

### Recommended Views (Hash-based routing)

| Route | View | Purpose |
|-------|------|---------|
| `#/dashboard` | Dashboard | Active timers, today's summary, quick clock-in |
| `#/employees` | Employee Management | CRUD employees, status overview |
| `#/projects` | Project & Task Management | CRUD projects/tasks, assignment |
| `#/timesheet` | Timesheet | Weekly grid view, manual entry, edit entries |
| `#/reports` | Reports | Filtered reports with charts, export |

### Navigation Pattern
- **Sidebar navigation** (left, collapsible) — best for 4-6 top-level views
- **Active timer banner** — persistent top bar showing running timer across all views
- **Breadcrumb-less** — flat hierarchy, no deep nesting needed

---

## 4. Report Generation Patterns

### Report Types
1. **Daily Summary** — Hours per employee for a specific date, broken down by project
2. **Weekly Summary** — 7-day grid showing hours per employee per day, with totals
3. **Monthly Summary** — Calendar heat map + totals per employee/project
4. **Project Summary** — Total hours per project, broken down by employee and task

### Chart Approach (No heavy library)
- **Bar charts**: Use CSS-based horizontal bars (percentage widths) for simplicity — no Chart.js needed
- **Summary tables**: Styled HTML tables with totals, averages
- **Heat map**: CSS grid with background-color opacity based on hours

### Export Strategy
- **CSV**: Vanilla JS — construct comma-separated string, create Blob, trigger download via `<a>` element with `download` attribute
- **PDF**: Use browser's `window.print()` with a print-specific CSS stylesheet (`@media print`). This avoids adding a PDF library entirely.
- **JSON**: Direct export of filtered data for developer/power-user use

---

## 5. Technical Patterns & Best Practices

### State Management (Vanilla JS)
- **Pub/Sub event bus** — Lightweight custom event system for cross-view communication
- **Data service layer** — `EmployeeService`, `ProjectService`, `TimeEntryService` classes wrapping localStorage CRUD + validation
- **Reactive rendering** — Each view has a `render()` method called when relevant data changes

### Time Handling
- **Store in UTC ISO strings** — Avoid timezone bugs
- **Display in local time** — Use `Intl.DateTimeFormat` for locale-aware formatting
- **Duration calculation** — Compute on clock-out: `(clockOut - clockIn) / 60000` for minutes
- **Running timer display** — `setInterval` updating every second for active timer

### Input Validation
- **Employee**: Name required, email format validation, unique email
- **Time Entry**: clockOut > clockIn, no overlapping entries for same employee, max 24h per entry
- **Project/Task**: Name required, unique within scope

### Keyboard Shortcuts (Power-user feature)
- `Space` — Quick clock in/out on dashboard
- `N` — New entry
- `E` — Employees view
- `P` — Projects view
- `R` — Reports view

---

## 6. UI Component Patterns

### Reusable Components to Build
1. **Modal** — For create/edit forms (employee, project, task, time entry)
2. **Data Table** — Sortable, filterable table for lists (employees, time entries)
3. **Timer Widget** — Clock display with start/stop/project selector
4. **Date Range Picker** — For report filtering (today, this week, this month, custom)
5. **Toast Notifications** — Success/error feedback on CRUD operations
6. **Dropdown Select** — Employee/project selection with search

### CSS Architecture
- **Tailwind utility-first** with a small set of custom components via `@apply`
- **CSS custom properties** for theme colors (easy dark mode later)
- **Responsive breakpoints**: Mobile-first, 768px tablet, 1024px desktop
- **Print stylesheet**: Separate `@media print` rules for clean report printing

---

## 7. Performance Considerations

- **Lazy rendering** — Only render visible view; don't render all views on load
- **Debounced search/filter** — 300ms debounce on search inputs
- **Efficient localStorage reads** — Cache data in memory, only read localStorage on app init; write on every mutation
- **Minimal dependencies** — Zero runtime JS libraries for core features; Tailwind is build-time only

---

## 8. Risk & Mitigation

| Risk | Impact | Mitigation |
|------|--------|------------|
| localStorage size limit hit | Data loss, app failure | Implement data export/archive, warn at 80% capacity |
| Browser tab closed during timer | Lost time entry | Save timer state to localStorage every 30s; recover on reload |
| Overlapping time entries | Inaccurate reports | Validate on save: reject overlaps for same employee |
| Timezone confusion | Wrong hours in reports | Store UTC, display local, show timezone in report headers |
| Large dataset → slow rendering | UI lag | Virtual scrolling for 100+ entries, pagination for reports |

---

## 9. Architecture Recommendations

### Recommended Approach: Single-Page App with ES Modules

```
src/
├── main.js                  # App entry, router init, global event bus
├── router.js                # Hash-based router
├── styles/
│   └── main.css             # Tailwind imports + custom components
├── data/
│   ├── store.js             # localStorage wrapper with caching
│   ├── employeeService.js   # Employee CRUD + validation
│   ├── projectService.js    # Project/Task CRUD + validation
│   └── timeEntryService.js  # Time entry CRUD + timer logic
├── views/
│   ├── dashboard.js         # Dashboard view with active timer
│   ├── employees.js         # Employee management view
│   ├── projects.js          # Project & task management view
│   ├── timesheet.js         # Timesheet grid view
│   └── reports.js           # Reports with filters & export
├── components/
│   ├── modal.js             # Reusable modal
│   ├── dataTable.js         # Sortable/filterable table
│   ├── timerWidget.js       # Clock in/out timer
│   ├── dateRangePicker.js   # Date range selection
│   └── toast.js             # Notification toasts
└── utils/
    ├── dates.js             # Date formatting, duration calculation
    ├── export.js            # CSV/print export utilities
    ├── validation.js        # Input validation helpers
    └── uuid.js              # UUID generation
```

### Why This Structure
- **Separation of concerns**: Data layer is independent of views
- **Testable**: Services can be unit tested without DOM
- **Scalable**: New views/components slot in without refactoring
- **No framework overhead**: Vanilla JS modules keep bundle tiny (~30KB estimated)

---

## 10. Summary of Key Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Template | `static-enhanced` | No backend needed; Vite + Tailwind gives fast DX and tiny bundles |
| Storage | localStorage with JSON | Simple, no server, sufficient for single-user/small-team scope |
| Routing | Hash-based SPA | No server config needed, works with static hosting |
| Charts | CSS-based (no library) | Keeps bundle small, sufficient for bar charts and summaries |
| Export | CSV (vanilla) + Print CSS (PDF) | Zero-dependency, browser-native |
| Time format | UTC storage, local display | Avoids timezone bugs in calculations |
| State | In-memory cache + localStorage persist | Fast reads, durable writes |

---

**Research Status: COMPLETE**
**Next Step: Architecture specification and development plan (Phase 2)**
