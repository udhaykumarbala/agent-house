# UX Specification: WorkTime Tracker

Based on research in: .plans/research/ux-research.md
UI direction from: .plans/research/ui-research.md
Architecture from: .plans/research/architecture-research.md
Security considerations from: .plans/research/security.md
Product context from: .plans/research/product-research.md

---

## Design Principles

1. **Timer-first** — Clock in/out is the #1 action; it must be reachable from every screen in one click (Toggl pattern)
2. **Glanceable** — Dashboard answers "what's happening now?" in under 3 seconds (Hubstaff pattern)
3. **Minimal input** — Only Employee + Project are required to start tracking; everything else is optional (product research insight)
4. **Export in one click** — Reports and CSV export are never more than one button away (Harvest/Clockify pattern)
5. **Progressive disclosure** — Empty states guide new users; power features stay out of the way until needed
6. **Safe by default** — All user input rendered via `textContent`, never `innerHTML` (security research requirement)

---

## Screen List

1. **Dashboard** — Today's overview: active timers, summary cards, recent entries
2. **Employees** — Employee list with status indicators, CRUD management
3. **Time Tracker** — Primary clock in/out interface with timeline of today's entries
4. **Projects** — Project and task management with color-coded cards
5. **Reports** — Filtered reports with charts, tables, and export

---

## Navigation Structure

### Sidebar Navigation (Left, Fixed)

```
┌──────────┐
│  ⏱ WTT   │  ← App logo/name
├──────────┤
│ ◉ Dash   │  ← Dashboard (default landing)
│ 👤 Team  │  ← Employees
│ ⏰ Track  │  ← Time Tracker
│ 📁 Proj  │  ← Projects
│ 📊 Report│  ← Reports
├──────────┤
│ ⚙ Settings│  ← Settings (bottom)
└──────────┘
```

**Behavior**:
- Fixed left sidebar, 240px wide on desktop
- Collapses to icon-only (64px) on tablet (<1024px)
- Transforms to bottom tab bar on mobile (<768px)
- Active view highlighted with indigo accent + filled icon
- Icons use Lucide icon set (outlined, 1.5px stroke)

### Persistent Timer Bar (Top)

```
┌─────────────────────────────────────────────────────────────────┐
│  🟢 Sarah Chen — Project Alpha  │  02:34:17  │ [⏹ Clock Out]  │
└─────────────────────────────────────────────────────────────────┘
```

**Behavior**:
- Appears across all views when any employee has an active timer
- Shows: employee name, current project, elapsed time (live-updating every second)
- "Clock Out" button is always accessible
- If multiple timers active, shows a count badge; clicking expands to show all
- Collapses to a minimal bar on mobile (name + time + stop button)
- When no timer is active, this bar is hidden (no wasted vertical space)

---

## Wireframes

### Screen 1: Dashboard

```
┌──────────┬──────────────────────────────────────────────────────┐
│          │  [🟢 Sarah — Project Alpha  02:34:17  ⏹ Clock Out ] │
│  ◉ Dash  ├──────────────────────────────────────────────────────┤
│  👤 Team │                                                      │
│  ⏰ Track │  Dashboard                            March 18, 2026│
│  📁 Proj ├──────────────────────────────────────────────────────┤
│  📊 Report│                                                      │
│          │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌─────────┐ │
│          │  │ Total Hrs│ │ Clocked  │ │ Active   │ │ Entries │ │
│          │  │  24.5h   │ │ In: 3/8  │ │ Projects │ │ Today   │ │
│          │  │ today    │ │ employees│ │    4     │ │   12    │ │
│          │  └──────────┘ └──────────┘ └──────────┘ └─────────┘ │
│          │                                                      │
│          │  Quick Clock In                                      │
│          │  ┌───────────────────────────────────────────────┐   │
│          │  │ [Select Employee ▼]  [Select Project ▼]       │   │
│          │  │                            [▶ Clock In]       │   │
│          │  └───────────────────────────────────────────────┘   │
│          │                                                      │
│          │  Today's Activity                                    │
│          │  ┌───────────────────────────────────────────────┐   │
│          │  │ 🟢 Sarah Chen    Project Alpha   8:00→ ...    │   │
│          │  │ ⚪ Mike Ross     Project Beta    8:15→10:30   │   │
│          │  │ 🟢 Lisa Park    Project Alpha   9:00→ ...    │   │
│          │  │ ⚪ John Lee     Project Gamma   7:00→11:00   │   │
│          │  │ ⚪ Mike Ross     Project Alpha  10:45→12:15   │   │
│          │  └───────────────────────────────────────────────┘   │
│          │                                                      │
│  ⚙ Set  │                                                      │
└──────────┴──────────────────────────────────────────────────────┘
```

**Key Elements**:
- **Summary cards** (top row): Total hours today, employees clocked in (count/total), active projects, total entries today. Cards use warm amber accent for numbers.
- **Quick Clock In** widget: Two dropdowns (Employee, Project) + a large "Clock In" button. Pre-fills last used project per employee. This is the fastest path to the core action.
- **Today's Activity** feed: Chronological list of today's time entries. Green dot = currently active. Shows employee, project, start time → end time (or "..." if active). Inline click to view/edit entry.

---

### Screen 2: Employees

```
┌──────────┬──────────────────────────────────────────────────────┐
│          │  Employees                        [+ Add Employee]  │
│  Sidebar ├──────────────────────────────────────────────────────┤
│          │                                                      │
│          │  ┌─────────────────────────────────────────────────┐ │
│          │  │ Name          │ Role      │ Status  │ Today    │ │
│          │  ├───────────────┼───────────┼─────────┼──────────┤ │
│          │  │ 🟢 Sarah Chen │ Developer │ Active  │ 2h 34m   │ │
│          │  │ ⚪ Mike Ross  │ Designer  │ Offline │ 3h 45m   │ │
│          │  │ 🟢 Lisa Park │ Developer │ Active  │ 1h 12m   │ │
│          │  │ ⚪ John Lee  │ Manager   │ Offline │ 4h 00m   │ │
│          │  │ 🟡 Amy Tan   │ QA        │ Break   │ 2h 10m   │ │
│          │  └─────────────────────────────────────────────────┘ │
│          │                                                      │
│          │  Showing 5 employees (3 active, 1 inactive)         │
│          │                                                      │
└──────────┴──────────────────────────────────────────────────────┘
```

**Key Elements**:
- **Employee table**: Sortable columns — Name, Role/Department, Status (green=clocked in, gray=offline, yellow=break), Today's hours
- **Status indicators**: Color-coded dots (research: Hubstaff pattern). 🟢 Clocked In, ⚪ Clocked Out, 🟡 On Break
- **Add Employee** button: Top-right, opens a modal form
- **Row click**: Opens employee detail/edit modal
- **Inline actions**: Hover reveals edit/archive icons per row

### Add/Edit Employee Modal

```
┌─────────────────────────────────────┐
│  Add Employee                    ✕  │
├─────────────────────────────────────┤
│                                     │
│  Name *        [________________]   │
│  Email         [________________]   │
│  Role          [________________]   │
│  Department    [Select ▼        ]   │
│  Hourly Rate   [________________]   │
│                                     │
│            [Cancel]  [Save]         │
└─────────────────────────────────────┘
```

**Validation**:
- Name: required, max 100 characters
- Email: optional, format validated
- Role: optional, free text
- Department: optional, dropdown (user creates departments)
- Hourly Rate: optional, numeric, for cost reports
- All inputs sanitized on save (security requirement)

---

### Screen 3: Time Tracker

```
┌──────────┬──────────────────────────────────────────────────────┐
│          │  Time Tracker                                        │
│  Sidebar ├──────────────────────────────────────────────────────┤
│          │                                                      │
│          │  ┌─────────────────────────────────────────────────┐ │
│          │  │           ┌─────────────────────┐               │ │
│          │  │           │     02:34:17        │               │ │
│          │  │           │  [large timer display]              │ │
│          │  │           └─────────────────────┘               │ │
│          │  │  [Employee ▼]  [Project ▼]  [Task ▼ (optional)] │ │
│          │  │  Notes: [________________________________]      │ │
│          │  │                                                 │ │
│          │  │       [⏹ Clock Out]    or    [▶ Clock In]      │ │
│          │  └─────────────────────────────────────────────────┘ │
│          │                                                      │
│          │  ── Mode: [Timer] [Manual] ──                       │
│          │                                                      │
│          │  Today's Entries                     March 18, 2026  │
│          │  ┌─────────────────────────────────────────────────┐ │
│          │  │ Time        Employee     Project    Duration    │ │
│          │  ├─────────────────────────────────────────────────┤ │
│          │  │ 8:00-       Sarah Chen   Alpha      2h 34m ▶   │ │
│          │  │ 10:45-12:15 Mike Ross    Alpha      1h 30m ✏   │ │
│          │  │ 8:15-10:30  Mike Ross    Beta       2h 15m ✏   │ │
│          │  │ 9:00-       Lisa Park    Alpha      1h 12m ▶   │ │
│          │  │ 7:00-11:00  John Lee     Gamma      4h 00m ✏   │ │
│          │  └─────────────────────────────────────────────────┘ │
│          │                                                      │
└──────────┴──────────────────────────────────────────────────────┘
```

**Key Elements**:
- **Timer display**: Large, monospace, tabular-nums font. Shows HH:MM:SS updating live. Uses Inter with `font-variant-numeric: tabular-nums` so digits don't shift width.
- **Input row**: Employee dropdown, Project dropdown, optional Task dropdown, optional Notes field. Only Employee and Project required.
- **Action button**: Large, prominent. Shows "Clock In" (green/indigo) when idle, "Clock Out" (amber/red) when timer running. State change is obvious.
- **Mode toggle**: Switch between Timer mode (real-time clock in/out) and Manual mode (enter start/end times directly). Timer is default.
- **Entry list**: Today's entries in a table. Active entries show ▶ indicator. Completed entries show ✏ (edit) on hover. Click to edit inline.
- **Inline editing**: Click any cell in the entry list to edit directly — no modal needed for quick fixes (Clockify pattern).

### Manual Entry Mode

```
┌─────────────────────────────────────────────────┐
│  Manual Time Entry                              │
│                                                 │
│  Employee *  [Select ▼]    Project * [Select ▼] │
│  Date        [2026-03-18]                       │
│  Start Time  [08:00]       End Time  [12:00]    │
│  Task        [Select ▼ (optional)]              │
│  Notes       [________________________________] │
│                                                 │
│              [Cancel]  [Save Entry]             │
└─────────────────────────────────────────────────┘
```

**Behavior**:
- Duration auto-calculates from start/end times
- Validates: end > start, no overlap with existing entries for same employee, max 24h
- Date defaults to today, can be changed for backfilling

---

### Screen 4: Projects

```
┌──────────┬──────────────────────────────────────────────────────┐
│          │  Projects                          [+ New Project]  │
│  Sidebar ├──────────────────────────────────────────────────────┤
│          │                                                      │
│          │  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐ │
│          │  │ 🔵 Project   │ │ 🟠 Project   │ │ 🟣 Project   │ │
│          │  │    Alpha     │ │    Beta      │ │    Gamma     │ │
│          │  │              │ │              │ │              │ │
│          │  │ 3 tasks      │ │ 2 tasks      │ │ 5 tasks      │ │
│          │  │ 124h tracked │ │ 56h tracked  │ │ 89h tracked  │ │
│          │  │ Active       │ │ Active       │ │ Completed    │ │
│          │  └──────────────┘ └──────────────┘ └──────────────┘ │
│          │                                                      │
│          │  ── Project Alpha (expanded) ──                     │
│          │  ┌─────────────────────────────────────────────────┐ │
│          │  │ Tasks                           [+ Add Task]    │ │
│          │  │ ☐ Frontend redesign              In Progress    │ │
│          │  │ ☐ API integration                Todo           │ │
│          │  │ ☑ Database migration             Done           │ │
│          │  └─────────────────────────────────────────────────┘ │
│          │                                                      │
└──────────┴──────────────────────────────────────────────────────┘
```

**Key Elements**:
- **Project cards**: Color-coded left border (user picks color on creation), shows name, task count, total hours tracked, status badge
- **Card click**: Expands to show tasks list inline (accordion pattern — no page navigation needed)
- **Task list**: Simple checklist-style. Name + status (Todo / In Progress / Done). Click to edit inline.
- **Status filter**: Tabs at top — All / Active / Completed / Archived
- **Color picker**: On project create/edit, a row of 8-10 preset colors to choose from (no custom hex input needed)

### Add/Edit Project Modal

```
┌─────────────────────────────────────┐
│  New Project                     ✕  │
├─────────────────────────────────────┤
│                                     │
│  Project Name * [________________]  │
│  Description    [________________]  │
│  Color          ● ● ● ● ● ● ● ●   │
│  Status         [Active ▼       ]   │
│                                     │
│            [Cancel]  [Create]       │
└─────────────────────────────────────┘
```

---

### Screen 5: Reports

```
┌──────────┬──────────────────────────────────────────────────────┐
│          │  Reports                           [📥 Export CSV]  │
│  Sidebar ├──────────────────────────────────────────────────────┤
│          │                                                      │
│          │  ┌─────────────────────────────────────────────────┐ │
│          │  │ Period: [Today] [This Week] [This Month] [Custom]│
│          │  │ Employee: [All ▼]    Project: [All ▼]           │ │
│          │  │ Group by: [Employee] [Project] [Date]           │ │
│          │  └─────────────────────────────────────────────────┘ │
│          │                                                      │
│          │  Summary                                             │
│          │  ┌──────────┐ ┌──────────┐ ┌──────────┐            │
│          │  │ Total    │ │ Avg/Day  │ │ Top      │            │
│          │  │ 186.5h   │ │  6.2h    │ │ Project  │            │
│          │  │ this week│ │ per emp  │ │ Alpha    │            │
│          │  └──────────┘ └──────────┘ └──────────┘            │
│          │                                                      │
│          │  Hours by Day (Bar Chart)                            │
│          │  ┌─────────────────────────────────────────────────┐ │
│          │  │ Mon ████████████████████  32h                   │ │
│          │  │ Tue ██████████████████    28h                   │ │
│          │  │ Wed ████████████████████████  38h               │ │
│          │  │ Thu ██████████████████    27.5h                 │ │
│          │  │ Fri ████████████████████████████  42h           │ │
│          │  │ Sat ████  6h                                    │ │
│          │  │ Sun ██████  13h                                 │ │
│          │  └─────────────────────────────────────────────────┘ │
│          │                                                      │
│          │  Detailed Breakdown                                  │
│          │  ┌─────────────────────────────────────────────────┐ │
│          │  │ Employee     │ Project │ Hours │ % of Total     │ │
│          │  ├──────────────┼─────────┼───────┼────────────────┤ │
│          │  │ Sarah Chen   │ Alpha   │ 32.5h │ ████ 17.4%    │ │
│          │  │ Sarah Chen   │ Beta    │  8.0h │ █ 4.3%        │ │
│          │  │ Mike Ross    │ Alpha   │ 24.0h │ ███ 12.9%     │ │
│          │  │ Mike Ross    │ Gamma   │ 16.0h │ ██ 8.6%       │ │
│          │  │ Lisa Park    │ Alpha   │ 28.0h │ ███ 15.0%     │ │
│          │  │ ...          │ ...     │ ...   │ ...            │ │
│          │  ├──────────────┼─────────┼───────┼────────────────┤ │
│          │  │ TOTAL        │         │186.5h │ 100%           │ │
│          │  └─────────────────────────────────────────────────┘ │
│          │                                                      │
└──────────┴──────────────────────────────────────────────────────┘
```

**Key Elements**:
- **Filter bar** (top): Date range presets (Today, This Week, This Month, Custom date picker), Employee dropdown, Project dropdown, Group-by toggle. Filters apply instantly.
- **Summary cards**: Total hours, average per day per employee, top project. Quick KPIs.
- **Bar chart**: CSS-based horizontal bars (no Chart.js — architecture decision). Shows hours by day/employee/project depending on group-by selection. Uses warm amber gradient fill.
- **Detail table**: Sortable, shows employee, project, hours, and a mini percentage bar. Bottom row shows totals.
- **Export CSV** button: Prominent, top-right. Single click downloads the currently filtered data as CSV. CSV sanitized against formula injection (security requirement).
- **Print/PDF**: Browser print button in export dropdown — uses `@media print` stylesheet for clean formatting.

---

### Screen 6: Settings

```
┌──────────┬──────────────────────────────────────────────────────┐
│          │  Settings                                            │
│  Sidebar ├──────────────────────────────────────────────────────┤
│          │                                                      │
│          │  General                                             │
│          │  ┌─────────────────────────────────────────────────┐ │
│          │  │ Time Format     ( ) 12-hour  (●) 24-hour       │ │
│          │  │ Week Starts On  [Monday ▼]                      │ │
│          │  │ Default View    [Dashboard ▼]                   │ │
│          │  └─────────────────────────────────────────────────┘ │
│          │                                                      │
│          │  Data Management                                     │
│          │  ┌─────────────────────────────────────────────────┐ │
│          │  │ [📥 Export All Data (JSON)]                     │ │
│          │  │ [📤 Import Data (JSON)]                         │ │
│          │  │ [🗑 Clear All Data]  ← requires confirmation    │ │
│          │  └─────────────────────────────────────────────────┘ │
│          │                                                      │
│          │  Storage: 1.2 MB of ~5 MB used                      │
│          │  ████░░░░░░░░░░░░░░░░  24%                          │
│          │                                                      │
└──────────┴──────────────────────────────────────────────────────┘
```

**Key Elements**:
- **Time format toggle**: 12h vs 24h display preference
- **Week start**: Monday (default) or Sunday — affects weekly reports
- **Export/Import**: JSON backup for data portability (product research: local-first, data integrity)
- **Clear All Data**: Red button, requires confirmation modal ("This will delete all employees, projects, and time entries. This cannot be undone.")
- **Storage indicator**: Visual bar showing localStorage usage (security research: warn at 80% capacity)

---

## User Flows

### Primary Flow: Clock In an Employee

1. User opens app → lands on **Dashboard**
2. Sees Quick Clock In widget at top of content area
3. Selects employee from dropdown (recent employees shown first)
4. Selects project from dropdown (last-used project pre-selected)
5. Clicks **[▶ Clock In]** button
6. **Feedback**: Button changes to "Clock Out" state, timer bar appears at top, entry appears in Today's Activity with green dot, summary cards update
7. Timer runs until user clicks **[⏹ Clock Out]** (from timer bar, dashboard, or Time Tracker view)
8. **Feedback**: Toast notification "Sarah clocked out — 2h 34m on Project Alpha", entry updates with end time and duration

**Alternative path**: User navigates to Time Tracker view for the same flow with a larger timer display and additional options (task, notes).

### Secondary Flow: Manual Time Entry

1. User navigates to **Time Tracker** view
2. Clicks **[Manual]** mode toggle
3. Fills in: Employee, Project, Date, Start Time, End Time
4. Optionally adds Task and Notes
5. Duration auto-calculates and displays
6. Clicks **[Save Entry]**
7. **Validation**: If times overlap with existing entry for that employee, shows inline error: "Overlaps with existing entry (8:00-12:00 on Project Alpha)"
8. **Feedback**: Toast "Time entry saved", entry appears in list below

### Tertiary Flow: Generate a Report

1. User navigates to **Reports** view
2. Selects date range preset (e.g., "This Week") or picks custom dates
3. Optionally filters by Employee and/or Project
4. Charts and table update instantly (no "Generate" button needed — reactive filtering)
5. Reviews summary cards and bar chart
6. Scrolls to detailed breakdown table
7. Clicks **[📥 Export CSV]**
8. **Feedback**: CSV file downloads immediately. Toast: "Report exported"

### Onboarding Flow: First-Time User

1. User opens app for the first time → Dashboard shows **empty state**
2. Empty state message: "Welcome to WorkTime Tracker! Let's get started."
3. Three guided steps shown as cards:
   - Step 1: "Add your first employee" → [+ Add Employee] button → opens Employee modal
   - Step 2: "Create a project" → [+ Create Project] button → opens Project modal
   - Step 3: "Clock someone in!" → Quick Clock In widget (enabled once steps 1 & 2 done)
4. After first clock-in, dashboard transitions to normal view with activity feed

### Employee Management Flow

1. User navigates to **Employees** view
2. Clicks **[+ Add Employee]** → modal opens
3. Fills in Name (required), optional fields
4. Clicks **[Save]** → modal closes, employee appears in list with ⚪ status
5. To edit: clicks employee row → same modal opens in edit mode
6. To archive: in edit modal, changes status to "Inactive" → employee grayed out in list, excluded from dropdowns

### Project Management Flow

1. User navigates to **Projects** view
2. Clicks **[+ New Project]** → modal opens
3. Fills in Name (required), picks a color, optional description
4. Clicks **[Create]** → project card appears in grid
5. Clicks project card → expands to show tasks
6. Clicks **[+ Add Task]** → inline input appears below task list
7. Types task name, presses Enter → task saved
8. Clicks task status chip to cycle: Todo → In Progress → Done

---

## Interaction Patterns

| Action | Element | Response |
|--------|---------|----------|
| Click | Clock In button | Timer starts, button changes to Clock Out state, timer bar appears, toast confirmation |
| Click | Clock Out button | Timer stops, duration saved, toast shows summary, entry finalized in list |
| Click | Sidebar nav item | View transitions with subtle 200ms slide, active item highlighted |
| Click | Employee row | Edit modal opens with pre-filled data |
| Click | Time entry in list | Inline edit mode activates for that row |
| Click | Project card | Accordion expands to show tasks |
| Click | Report date preset | Charts and table update immediately (no submit button) |
| Click | Export CSV | File download starts, toast confirms |
| Hover | Table row | Row highlights, action icons (edit/delete) appear on right |
| Focus | Dropdown | Opens with search-as-you-type filtering |
| Keyboard | Space (on Dashboard) | Quick clock in/out for last-used employee+project |
| Keyboard | N | Opens new entry form |
| Keyboard | Escape | Closes any open modal or dropdown |

---

## States

### Empty States

| View | Empty Condition | Message | CTA |
|------|----------------|---------|-----|
| Dashboard (first use) | No employees exist | "Welcome! Add your first team member to get started." | [+ Add Employee] |
| Dashboard (no activity) | Employees exist but none clocked in today | "No activity yet today. Clock someone in to start tracking." | Quick Clock In widget (already visible) |
| Employees | No employees | "No team members yet. Add employees to start tracking their time." | [+ Add Employee] |
| Projects | No projects | "Create a project to organize your team's work." | [+ New Project] |
| Time Tracker | No entries today | "No time entries for today. Use the timer above or add a manual entry." | Timer widget (already visible) |
| Reports | No data for filter | "No time entries match your filters. Try a different date range or remove filters." | [Clear Filters] |

### Loading State

- **Initial app load**: Skeleton screens (gray placeholder blocks matching layout) for 0-200ms while localStorage is read and parsed
- **Timer**: Smooth number transition animation (digits roll/fade, not jump)
- **Report recalculation**: Instant for <1000 entries; for larger datasets, show a subtle progress indicator on the chart area

### Error States

| Error | Display | Recovery |
|-------|---------|----------|
| Overlapping time entry | Inline red text below time fields: "This overlaps with [entry details]" | User adjusts times |
| Required field empty | Red border + "Required" label below field | User fills in field |
| localStorage full | Banner at top of app: "Storage is full. Export and clear old data to continue." | [Export Data] [Go to Settings] |
| Corrupted data on load | Modal: "Some data could not be loaded. [View Details] [Restore Defaults]" | Export what's recoverable, or reset |
| Invalid time (end < start) | Inline red text: "End time must be after start time" | User corrects times |

### Success States

| Action | Feedback |
|--------|----------|
| Clock in | Timer bar animates in, green dot appears, toast: "Clocked in" |
| Clock out | Toast: "[Name] clocked out — [duration] on [Project]" |
| Save entry | Toast: "Time entry saved" with checkmark |
| Save employee/project | Modal closes, item appears in list with brief highlight animation |
| Export CSV | Toast: "Report exported as CSV" |
| Data import | Toast: "[N] employees, [N] projects, [N] entries imported" |

---

## Accessibility Considerations

### Keyboard Navigation
- Full tab navigation through all interactive elements
- Sidebar items focusable with arrow keys (up/down to navigate, Enter to select)
- Modal focus trap: Tab cycles within modal when open, Escape closes
- Skip-to-content link (hidden until focused) to bypass sidebar
- Timer Start/Stop accessible via Space bar on Dashboard

### Screen Reader Support
- All views have `<h1>` for the page title announced on navigation
- Timer status announced: "Timer running: 2 hours 34 minutes on Project Alpha for Sarah Chen"
- Status dots have aria-labels: "Status: Clocked In" / "Status: Offline"
- Table headers use `<th scope="col">` for proper association
- Toast notifications use `role="status"` with `aria-live="polite"`
- Form fields have associated `<label>` elements
- Error messages linked to inputs via `aria-describedby`

### Visual Accessibility
- Touch targets: Minimum 44x44px for all interactive elements
- Color contrast: 4.5:1 minimum for text, 3:1 for large text and UI components
- Status indicators use both color AND shape/icon (not color alone): 🟢 = filled circle + "Active" text, ⚪ = outline circle + "Offline" text
- Focus indicators: Visible 2px outline in indigo accent color on all focusable elements
- Reduced motion: Respect `prefers-reduced-motion` — disable timer animation, use instant transitions
- Font sizes: Base 16px, minimum 14px for secondary text. User can zoom to 200% without horizontal scroll

### Responsive Behavior

| Breakpoint | Layout Change |
|------------|--------------|
| ≥1024px (Desktop) | Full sidebar (240px) + content area |
| 768-1023px (Tablet) | Collapsed sidebar (64px icons only) + content area |
| <768px (Mobile) | No sidebar — bottom tab bar (5 icons). Timer bar becomes compact. Tables become card lists. Modals become full-screen sheets. |

---

## Duration Formatting

Consistent across all views:
- **Short format**: `2h 30m` (used in tables, cards, summaries)
- **Timer format**: `02:34:17` (used in active timer display only)
- **Never**: "2.5 hours" or "150 minutes" — these are not scannable (UX research anti-pattern)

## Date Formatting

- **Display**: Locale-aware via `Intl.DateTimeFormat` (e.g., "March 18, 2026" or "18 Mar 2026")
- **Input**: Native date/time pickers where available
- **Storage**: ISO 8601 UTC (architecture decision)

---

## Component Inventory

| Component | Used In | Notes |
|-----------|---------|-------|
| Sidebar Nav | All views | Fixed, collapsible, icon + label |
| Timer Bar | All views (when active) | Persistent top banner, live-updating |
| Summary Card | Dashboard, Reports | Stat number + label + optional icon |
| Data Table | Employees, Time Tracker, Reports | Sortable columns, hover actions, inline edit |
| Modal | Employee CRUD, Project CRUD, Settings confirmations | Centered overlay, focus trap, Escape to close |
| Dropdown Select | Time Tracker, Reports filters | Search-as-you-type, recent items first |
| Toast Notification | All views | Bottom-right, auto-dismiss 3s, stacks vertically |
| Date Range Picker | Reports | Presets + custom range |
| Project Card | Projects view | Color-coded, expandable accordion for tasks |
| Empty State | All views | Illustration + message + CTA button |
| Status Indicator | Employees, Dashboard | Colored dot + text label |
| Bar Chart (CSS) | Reports | Horizontal bars, percentage-width, no JS library |
| Quick Clock In | Dashboard | Compact employee + project + button row |

---
Status: READY_FOR_REVIEW
