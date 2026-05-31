# UX Research: Work Time Tracker

## Pattern Analysis

Referenced apps and their successful patterns:

### Toggl Track
- **Pattern Used**: One-click timer with prominent "Start" button at top; timeline view for daily entries; minimal-input time entry (description + project tag + start/stop)
- **Why It Works**: Reduces friction to near-zero for the core action (tracking time). Users don't have to navigate away from the main screen to start/stop a timer. The bold, colored Start button uses Fitts's Law — large target, always visible.
- **What We Can Learn**: The primary action (clock in/out) must be the most prominent element on the screen. A persistent timer bar or header-level CTA ensures users never hunt for it. Pre-filling project/task from recent entries reduces repetitive input.

### Clockify
- **Pattern Used**: Left sidebar navigation (Dashboard, Time Tracker, Calendar, Reports, Projects, Team); top bar with running timer; tabular time entry list with inline editing; report filters (date range, project, team member) with chart + table view
- **Why It Works**: The sidebar gives orientation — users always know where they are. Inline editing lets users fix mistakes without a modal. Reports combine a visual summary chart on top with a detailed table below, satisfying both glancers and detail seekers.
- **What We Can Learn**: Sidebar navigation scales well for our 4–5 views. Reports should lead with a visual summary (bar chart or donut) and let users drill into the table underneath. Inline editing of time entries prevents disruptive context switches.

### Harvest
- **Pattern Used**: Week-view grid (rows = projects, columns = days) for time entry; "New Entry" as a plus button inside the grid cell; running timer indicator in the top nav; report export to CSV/PDF from a dedicated Reports section with date-range picker and grouping options
- **Why It Works**: The grid matches how managers think — "which project, which day, how many hours?" It mirrors a spreadsheet mental model that's already familiar. Export options are clearly labeled and accessible from the report view, not buried in settings.
- **What We Can Learn**: For report generation, offer clear date-range presets (Today, This Week, This Month) alongside a custom picker. Group reports by employee, by project, or by date. CSV export should be a single button, not a multi-step flow.

### Hubstaff
- **Pattern Used**: Employee-centric dashboard showing each team member's status (active/idle/offline), hours today, and current task; activity timeline with screenshots; weekly summary cards per employee
- **Why It Works**: Managers get an at-a-glance view of who's working on what. The status indicators (green dot = active, gray = offline) are universally understood. Summary cards condense a week into a scannable block.
- **What We Can Learn**: The employee list view should show current status (clocked in/out), today's hours, and current project at a glance. Use color-coded status indicators. Summary cards per employee are effective for the weekly/monthly views.

### Jira Work Logs (Atlassian)
- **Pattern Used**: Time logging tied directly to tasks/issues; "Log Work" modal with time spent, date, and description; time tracking report aggregated by issue, sprint, or person
- **Why It Works**: Linking time to tasks gives context to every entry. The modal is simple — just 3 fields. Reports are flexible because they can pivot on different dimensions.
- **What We Can Learn**: Every time entry should be linkable to a project and task. This enables meaningful reports ("How much time did Project X consume this month?"). Keep the entry form to 4–5 fields max.

## User Journey Analysis

### Entry Point
A manager or team lead opens the app to either:
- **Check status**: "Who's clocked in right now?" → lands on Dashboard
- **Record time**: "I need to clock someone in or log hours" → navigates to Time Tracker
- **Pull a report**: "What did the team do last week?" → navigates to Reports

The Dashboard should serve as the default landing page, showing today's overview.

### Core Loop
1. **Clock in** an employee (or employee self-clocks) → select employee, confirm project/task, start timer
2. **Work happens** → timer runs, optionally switch tasks
3. **Clock out** → stop timer, entry saved automatically
4. **Review** → check hours at end of day/week
5. **Report** → generate and export summary

The core loop repeats daily. Clock in → work → clock out → review.

### Success State
- **Immediate**: Timer stops, entry appears in the timeline with correct hours, project, and task. A brief confirmation (checkmark animation or toast) tells the user "saved."
- **End of period**: A clean report showing total hours per employee per project, exportable in one click. The manager can see at a glance if hours look right.

## Mental Model

Users think of this task as: **a digital punch card + timesheet**.

They have a mental model of a grid: **employees on one axis, days on the other, with hours in the cells**. Reports are summaries that collapse this grid along different dimensions (by employee, by project, by date range).

Common terminology:
- "Clock in / Clock out" (not "start session" or "begin tracking")
- "Time entry" or "log" (a single block of recorded time)
- "Timesheet" (a collection of entries for a period)
- "Hours worked" (not "duration" or "time spent")
- "Project" → "Task" (hierarchical: projects contain tasks)
- "Report" / "Summary" (aggregated view)
- "Export" (get it out as CSV/PDF)

## Anti-Patterns to Avoid

- **Don't**: Require too many fields to clock in — **Why**: Every extra field adds friction to the most frequent action. Toggl's success is built on 1-click start. If clocking in requires employee + project + task + description + category, users will abandon the tool. Keep mandatory fields to 2 (employee, project) and make everything else optional or pre-filled.

- **Don't**: Use a calendar-only view for time entry — **Why**: Clockify and Harvest both offer a list/timeline view because calendar views are hard to use for sub-hour entries and don't show enough metadata. Calendar can be a secondary view but not the primary input method.

- **Don't**: Hide export behind multiple clicks or settings — **Why**: Report export is a critical end-of-workflow action. Harvest and Clockify put the export button directly on the reports page. If users have to go to Settings > Export > Configure > Download, they'll screenshot the screen instead.

- **Don't**: Force users to manually calculate hours — **Why**: Time entries should auto-calculate duration from clock-in/clock-out times. Manual hour entry should be an alternative, not the default. Calculation errors are the #1 frustration in paper timesheets.

- **Don't**: Show an empty dashboard with no guidance — **Why**: First-time users need to add employees and projects before they can track time. An empty state that just says "No data" is a dead end. Guide them: "Add your first employee to get started."

- **Don't**: Use ambiguous time formats — **Why**: Always show both 12h and 24h consistently (based on user locale/preference). Display durations as "2h 30m" not "2.5 hours" or "150 minutes" — the h/m format is instantly scannable.

## Recommended Patterns for Our App

Based on research, we should use:

1. **Persistent top-bar timer with prominent Clock In/Out button** — because the primary action must be instantly accessible from any view (learned from Toggl and Clockify). The button should change color/state to indicate active timing.

2. **Left sidebar navigation with 5 sections: Dashboard, Employees, Time Tracker, Projects, Reports** — because this is the established pattern in Clockify, Hubstaff, and Harvest. Users expect horizontal categorization via sidebar in productivity tools. Icons + labels for clarity.

3. **Dashboard as landing page with today's summary cards** — because managers' first question is "what's happening today?" Show: total hours today, employees clocked in, active projects, and a mini-timeline of recent entries (pattern from Hubstaff).

4. **Inline time entry with minimal required fields (Employee + Project)** — because reducing friction on the core action is the #1 UX priority (learned from Toggl). Task and description should be optional. Support both timer mode (clock in/out) and manual entry mode (enter start/end times directly).

5. **Reports with visual summary + filterable table + one-click export** — because this is the universal pattern across Clockify, Harvest, and Hubstaff. Lead with a bar chart (hours by day or by employee), follow with a sortable table, and place the Export CSV button prominently at the top-right of the report view. Include date-range presets (Today, This Week, This Month, Custom).

6. **Color-coded status indicators for employees** — because Hubstaff's green/gray dots are instantly understood. Show clocked-in (green), clocked-out (gray), on-break (yellow) at a glance in the employee list.

7. **Progressive empty states with CTAs** — because new users need onboarding guidance. Each empty view should explain what it does and offer a primary action: "No employees yet → Add Employee", "No time entries → Clock someone in", "No projects → Create a project."

---
Status: READY_FOR_PLANNING
