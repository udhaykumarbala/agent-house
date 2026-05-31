# WorkTime Tracker — UI Flow Documentation

## Routes

| Route | View | Description |
|-------|------|-------------|
| `#/dashboard` | Dashboard | Today's overview, summary cards, active timers, quick clock-in |
| `#/employees` | Employees | Employee list with CRUD management |
| `#/timesheet` | Timesheet | Time entry list and manual entry |
| `#/projects` | Projects | Project & task management |
| `#/reports` | Reports | Filtered reports with charts and export |

## Phase 2 Implemented Views

### Dashboard (`#/dashboard`)

```
┌─────────────────────────────────────────────────────────┐
│  Dashboard                               March 18, 2026 │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐  │
│  │ 24h 30m  │ │  3/8     │ │    4     │ │   12     │  │
│  │ TOTAL    │ │ CLOCKED  │ │ ACTIVE   │ │ ENTRIES  │  │
│  │ HOURS    │ │ IN       │ │ PROJECTS │ │ TODAY    │  │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘  │
│                                                         │
│  Active Timers                                          │
│  ┌─────────────────────────────────────────────────┐   │
│  │ 🟢 Sarah Chen  — Project Alpha    02:34:17 [Stop]│  │
│  │ 🟢 Lisa Park   — Project Beta     01:12:05 [Stop]│  │
│  └─────────────────────────────────────────────────┘   │
│                                                         │
│  Quick Clock In                                         │
│  ┌─────────────────────────────────────────────────┐   │
│  │ [Employee ▼]  [Project ▼]       [▶ Clock In]    │   │
│  └─────────────────────────────────────────────────┘   │
│                                                         │
│  Today's Activity                                       │
│  ┌─────────────────────────────────────────────────┐   │
│  │ 🟢 Sarah     Alpha    08:00 → ...    2h 34m    │   │
│  │ ⚪ Mike      Beta     08:15 → 10:30  2h 15m    │   │
│  │ 🟢 Lisa      Beta     09:00 → ...    1h 12m    │   │
│  │ ⚪ John      Gamma    07:00 → 11:00  4h 00m    │   │
│  └─────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
```

**Empty State** (no employees):
```
┌─────────────────────────────────────────────────┐
│  Welcome to WorkTime Tracker!                   │
│  Add your first employee to get started.        │
│              [+ Add Employee]                   │
└─────────────────────────────────────────────────┘
```

**Data Flow:**
- Summary cards compute real-time from `timeEntryService` and `employeeService`
- Active timers update every 1s via `setInterval`
- Quick Clock In triggers `timeEntryService.clockIn()` → emits `timer:started`
- Dashboard subscribes to `timer:changed`, `timeEntries:changed`, `employees:changed`, `projects:changed`

---

### Projects (`#/projects`)

```
┌─────────────────────────────────────────────────────────┐
│  Projects                              [+ New Project]  │
├─────────────────────────────────────────────────────────┤
│  [All] [Active] [Completed] [Archived]                  │
│                                                         │
│  ┌──────────────────┐ ┌──────────────────┐             │
│  │▎Project Alpha     │ │▎Project Beta      │            │
│  │ API Integration   │ │ Design system     │            │
│  │                   │ │                   │            │
│  │ 📋 3 tasks  ⏱ 124h│ │ 📋 2 tasks  ⏱ 56h │           │
│  │ [Active]     ⌄    │ │ [Active]     ⌄    │           │
│  └──────────────────┘ └──────────────────┘             │
│                                                         │
│  ── Project Alpha (expanded) ──                         │
│  ┌─────────────────────────────────────────────────┐   │
│  │ Tasks                            [+ Add Task]    │   │
│  │ ☐ Frontend redesign             [In Progress]    │   │
│  │ ☐ API integration               [Todo]           │   │
│  │ ☑ Database migration            [Done]      ✕    │   │
│  └─────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
```

**Empty State** (no user-created projects):
```
┌─────────────────────────────────────────────────┐
│  Create your first project                       │
│  Create a project to organize your team's work.  │
│              [+ New Project]                     │
└─────────────────────────────────────────────────┘
```

**Interactions:**
- **New Project** → Opens modal (name*, description, color picker, status)
- **Card click** → Expands/collapses task accordion
- **Task status badge click** → Cycles: Todo → In Progress → Done
- **Add Task** → Inline input within expanded project
- **Edit/Delete** → Hover actions on project card (pencil/trash icons)
- **Status filter tabs** → Filters All / Active / Completed / Archived

**Data Flow:**
- Projects rendered from `projectService.getAll()`
- Tasks from `projectService.getTasksByProject()`
- Hours from `timeEntryService.getAll().filter(e => e.projectId)`
- Subscribes to `projects:changed`, `tasks:changed`, `timeEntries:changed`

---

## Event Bus Events

| Event | Emitter | Consumers |
|-------|---------|-----------|
| `timer:started` | timeEntryService.clockIn | Dashboard, Timer bar |
| `timer:stopped` | timeEntryService.clockOut | Dashboard, Timer bar |
| `timer:changed` | timeEntryService | Dashboard |
| `timeEntries:changed` | timeEntryService | Dashboard, Projects |
| `employees:changed` | employeeService | Dashboard |
| `projects:changed` | projectService | Dashboard, Projects |
| `tasks:changed` | projectService | Projects |

## Phase 3 Components

### DateRangePicker (`dateRangePicker.ts`)

```
┌─────────────────────────────────────────────────────────┐
│  [Today] [This Week] [This Month] [Last Week]           │
│  [Last Month] [Custom]                                   │
│                                                          │
│  (When Custom selected:)                                 │
│  From [2026-03-01]  To [2026-03-18]                     │
│                                                          │
│  Selection: This Week                                    │
└─────────────────────────────────────────────────────────┘
```

**Props / Options:**
- `onChange(range: DateRange)` — callback receiving `{ start, end, preset }`
- `initialPreset` — default: `'thisWeek'`
- `initialStart` / `initialEnd` — for restoring custom range state

**Presets:**
| Key | Label | Range |
|-----|-------|-------|
| `today` | Today | Single day |
| `thisWeek` | This Week | Mon–Sun of current week |
| `thisMonth` | This Month | 1st to last day of current month |
| `lastWeek` | Last Week | Mon–Sun of previous week |
| `lastMonth` | Last Month | 1st to last day of previous month |
| `custom` | Custom | User picks start and end dates |

**Validation:**
- Custom range: end date must be >= start date
- Invalid range shows inline error, inputs get `.error` class, onChange not called

**Styling:**
- Preset buttons use `.date-preset` / `.date-preset.selected` classes
- Custom inputs use `.input` class with compact sizing
- Error message uses `.form-error` class
- Root has `role="group"` and `aria-label="Date range picker"`
- Buttons have `aria-pressed` for screen readers

---

## Component Reuse

| Component | Used By |
|-----------|---------|
| `Modal` | Projects (Add/Edit/Delete), Employees, Settings, Import, Clear Data |
| `Toast` | All views (success/error notifications) |
| `DataTable` | Employees, Timesheet, Reports |
| `DateRangePicker` | Timesheet (filter bar), Reports (filter bar) |
| `Settings` | Sidebar gear icon, Mobile settings button |
| Status dots (`.status-dot`) | Dashboard (active timers, activity) |
| Status badges (`.badge-*`) | Projects, Employees |
| Stat cards (`.stat-card`) | Dashboard, Reports |

---

## Phase 4 — Settings, Keyboard Shortcuts, Responsive Polish

### Settings Panel (Modal)

Accessible from gear icon in sidebar (desktop/tablet) or mobile top bar settings button.

```
┌─────────────────────────────────────────────┐
│  Settings                               ✕   │
├─────────────────────────────────────────────┤
│                                              │
│  General                                     │
│  ┌──────────────────────────────────────┐   │
│  │ Time Format     [12-hour] [24-hour]  │   │
│  │ Week Starts On  [Monday ▼]           │   │
│  │ Default View    [Dashboard ▼]        │   │
│  └──────────────────────────────────────┘   │
│                                              │
│  Data Management                             │
│  ┌──────────────────────────────────────┐   │
│  │ [Export All Data (JSON)]              │   │
│  │ [Import Data (JSON)]                  │   │
│  │ [Clear All Data]   ← red/destructive  │   │
│  └──────────────────────────────────────┘   │
│                                              │
│  Storage: 1.2 KB of ~5.00 MB used     0%    │
│  ████░░░░░░░░░░░░░░░░                       │
│                                              │
└─────────────────────────────────────────────┘
```

**Actions:**
- **Export All Data**: Downloads `worktime-backup-YYYY-MM-DD.json` with all wtt_ collections
- **Import Data**: Opens file picker → validates JSON schema → offers Merge or Replace
- **Clear All Data**: Opens confirmation modal → wipes all wtt_ keys → re-creates General project
- **Time Format**: Toggles 12h/24h, saved to localStorage
- **Week Starts On**: Monday/Sunday, affects weekly reports
- **Default View**: Sets which view loads on app open

**Import Dialog:**
```
┌─────────────────────────────────────────────┐
│  Import Data                            ✕   │
├─────────────────────────────────────────────┤
│  Found 42 items across 3 collections:        │
│  employees, projects, time_entries           │
│                                              │
│  Choose how to import:                       │
│  [Merge with existing]  [Replace all data]   │
└─────────────────────────────────────────────┘
```

### Keyboard Shortcuts (`keyboard.ts`)

Global keyboard shortcuts, disabled when an input/textarea/select is focused.

| Key | Action | Condition |
|-----|--------|-----------|
| `Space` | Toggle clock in/out | Only on Dashboard |
| `N` | Open new time entry | Navigates to Timesheet |
| `D` | Go to Dashboard | — |
| `E` | Go to Employees | — |
| `P` | Go to Projects | — |
| `R` | Go to Reports | — |
| `Esc` | Close open modal | — |
| `?` | Show shortcuts help | — |

**Help Modal** (press `?`):
```
┌─────────────────────────────────────────────┐
│  Keyboard Shortcuts                     ✕   │
├─────────────────────────────────────────────┤
│  [Space]  Clock in/out (on Dashboard)        │
│  [N]      New time entry                     │
│  [D]      Go to Dashboard                    │
│  [E]      Go to Employees                    │
│  [P]      Go to Projects                     │
│  [R]      Go to Reports                      │
│  [Esc]    Close modal                        │
│  [?]      Show this help                     │
└─────────────────────────────────────────────┘
```

### Responsive Layout

| Breakpoint | Layout |
|------------|--------|
| < 768px (Mobile) | Bottom tab bar, hamburger menu for sidebar drawer, full-screen modals, stacked stat cards, compact timer bar |
| 768–1023px (Tablet) | Icon-only sidebar (64px), 2-column stat cards, horizontal table scroll |
| >= 1024px (Desktop) | Full sidebar (240px) + content area, 4-column stat cards, full tables |

**Mobile Top Bar:**
```
┌──────────────────────────────────────────┐
│  [☰]     ⏱ WorkTime           [⚙]       │
└──────────────────────────────────────────┘
```
- Hamburger opens sidebar as a slide-out drawer (left)
- Settings gear opens settings modal
- Bottom tab bar for primary navigation

**Mobile Sidebar Drawer:**
- Slides in from left with backdrop overlay
- Shows full nav labels + settings
- Closes on: backdrop click, nav item click, settings click

### Event Bus (Phase 4 additions)

| Event | Emitter | Consumers |
|-------|---------|-----------|
| `settings:changed` | Settings panel | (future use for live preference updates) |
