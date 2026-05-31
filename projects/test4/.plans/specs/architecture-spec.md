# Architecture Specification: DesignFlow — UX Designer Todo App

**Version:** 1.0
**Date:** 2026-03-21
**Status:** READY_FOR_DEVELOPMENT
**Architect:** Software Architect Agent

---

## 1. Concept & Vision

**DesignFlow** is a todo application built *by* designers *for* designers. It mirrors the double diamond design process — tasks flow through Backlog → In Progress → In Review → Approved — rather than the binary complete/incomplete model that generic todo apps impose. The tool is a designer's trusted daily companion: calm, beautiful, and deeply aware of how creative work actually gets done.

**Personality**: "A well-organized design studio. Everything has its place, the aesthetic is warm and editorial, and the interface stays out of the way while you work."

---

## 2. Design Language

### 2.1 Color Palette

**Light Mode (Primary):**

| Token | Hex | Usage |
|-------|-----|-------|
| `bg-primary` | `#FAF8F5` | Main background — warm cream |
| `bg-surface` | `#FFFFFF` | Cards, panels — pure white |
| `bg-surface-elevated` | `#FFFFFF` | Modals, dropdowns |
| `text-primary` | `#1C1917` | Headings, task titles |
| `text-secondary` | `#78716C` | Metadata, labels |
| `text-muted` | `#A8A29E` | Placeholders, hints |
| `border` | `#E7E5E4` | Card borders, dividers |
| `accent-primary` | `#3730A3` | Indigo — CTAs, active states, links |
| `accent-warm` | `#F97316` | Coral — highlights, urgent items |
| `shadow` | `rgba(28,25,23,0.06)` | Warm-tinted soft shadows |

**Dark Mode:**

| Token | Hex | Usage |
|-------|-----|-------|
| `bg-primary` | `#0C0A09` | Deep warm black |
| `bg-surface` | `#1C1917` | Cards, panels |
| `bg-surface-elevated` | `#292524` | Modals, dropdowns |
| `text-primary` | `#FAF8F5` | Headings, task titles |
| `text-secondary` | `#A8A29E` | Metadata, labels |
| `text-muted` | `#78716C` | Placeholders, hints |
| `border` | `#292524` | Card borders, dividers |
| `accent-primary` | `#818CF8` | Indigo-light — CTAs, links |
| `accent-warm` | `#FB923C` | Peach — highlights |
| `shadow` | `rgba(0,0,0,0.3)` | Elevated surfaces |

**Category Colors (both modes):**

| Category | Light Hex | Dark Hex | Usage |
|----------|-----------|----------|-------|
| Research | `#166534` | `#86EFAC` | Sage green |
| Design | `#3730A3` / `#818CF8` | `#818CF8` | Indigo |
| Testing | `#B45309` | `#FCD34D` | Amber |
| Review | `#9D174D` | `#F9A8D4` | Rose |
| Discovery | `#0F766E` | `#5EEAD4` | Teal |
| Handoff | `#475569` | `#94A3B8` | Slate |
| Admin | `#6B7280` | `#9CA3AF` | Gray |

### 2.2 Typography

**Font Stack:**
- **Display/Headings**: `Cabinet Grotesk` (Google Fonts, variable) — warm geometric sans with personality
- **Body/UI**: `Inter` (Google Fonts, variable) — proven legibility for dense UI
- **Monospace** (time estimates, dates): `JetBrains Mono` (Google Fonts)

**Type Scale:**

| Token | Size | Weight | Line Height | Usage |
|-------|------|--------|-------------|-------|
| `text-xs` | 12px | 400 | 1.4 | Badges, timestamps |
| `text-sm` | 14px | 400 | 1.5 | Body, task metadata |
| `text-base` | 16px | 500 | 1.5 | Task titles (card) |
| `text-lg` | 18px | 600 | 1.3 | Column headers |
| `text-xl` | 24px | 700 | 1.2 | Page title |
| `text-2xl` | 32px | 700 | 1.1 | Empty state headings |

### 2.3 Spacing & Layout

- **Base unit**: 4px
- **Card padding**: 16px
- **Column gap**: 16px
- **Section spacing**: 24px
- **Page margins**: 32px (desktop), 16px (mobile)
- **Border radius**: `8px` (buttons, inputs), `12px` (cards), `16px` (panels)

### 2.4 Motion Philosophy

Every animation serves one of three purposes:
1. **Feedback** — confirms user action (button press, drag release)
2. **Orientation** — helps user understand spatial changes (panel slide-in)
3. **Delight** — subtle polish that rewards interaction (task completion)

| Animation | Duration | Easing | Purpose |
|-----------|----------|--------|---------|
| Card hover lift | 150ms | ease-out | Feedback |
| Drag spring | 200ms | spring(1, 80, 10) | Orientation |
| Panel slide-in | 250ms | ease-out | Orientation |
| Task complete | 300ms | ease-out | Delight |
| Filter fade | 200ms | ease-in-out | Orientation |
| Reduced motion | 100ms | linear | Accessibility |

**Reduced Motion**: All animations collapse to `opacity` fades (100ms) when `prefers-reduced-motion` is set.

### 2.5 Visual Assets

- **Icons**: Lucide React — consistent 1.5px stroke, outlined style
- **No images**: Pure CSS/SVG aesthetic — no external image dependencies
- **Decorative**: Subtle warm gradient on header, dot-grid pattern in empty column states

---

## 3. Layout & Structure

### 3.1 Page Architecture

```
┌─────────────────────────────────────────────────────────┐
│  HEADER (fixed, 56px)                                   │
│  [Logo] DesignFlow    [View Toggle] [Search] [⚙️] [🌙]   │
├─────────────────────────────────────────────────────────┤
│  QUICK ADD BAR (sticky, 64px)                           │
│  [🔍 +] Add a task... (Cmd+K)           [Category ▼]    │
├─────────────────────────────────────────────────────────┤
│  FILTER BAR (collapsible, ~48px)                         │
│  [All] [Research] [Design] [Testing] [Review]    [+tag] │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  KANBAN BOARD (scrollable, flex row)                     │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐   │
│  │ BACKLOG  │ │IN PROGRESS│ │IN REVIEW │ │ APPROVED │   │
│  │    12    │ │    3     │ │    2     │ │    8     │   │
│  ├──────────┤ ├──────────┤ ├──────────┤ ├──────────┤   │
│  │ [Card]   │ │ [Card]   │ │ [Card]   │ │ [Card]   │   │
│  │ [Card]   │ │ [Card]   │ │ [Card]   │ │ [Card]   │   │
│  │ [Card]   │ │          │ │          │ │ [Card]   │   │
│  │ ...      │ │          │ │          │ │ ...      │   │
│  │ + Add    │ │          │ │          │ │          │   │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘   │
│                                                          │
├─────────────────────────────────────────────────────────┤
│  STATUS BAR (32px)                                       │
│  25 tasks · 3 in progress · Last saved just now          │
└─────────────────────────────────────────────────────────┘
```

### 3.2 Views

1. **Kanban Board** (default) — 4 horizontal columns with drag-and-drop
2. **List View** — Dense flat list grouped by iteration state, sortable
3. **Focus View** — Single-column, one task expanded at a time

### 3.3 Responsive Strategy

- **Desktop (>1024px)**: Full kanban, 4 visible columns, side panel for task detail
- **Tablet (768-1024px)**: 2-3 visible columns, horizontal scroll, modal for task detail
- **Mobile (<768px)**: Single column list view, bottom sheet for task detail

---

## 4. Features & Interactions

### 4.1 Quick Add

**Trigger**: Click the quick add bar OR press `Cmd/Ctrl+K` from anywhere.

**Behavior**:
- Natural language parsing: "Design new CTA button #Design @tomorrow 2h"
- Auto-categorization from `#` prefix (maps to category)
- Auto-assigns iteration state from column context (if adding from specific column)
- Auto-focuses category dropdown if no `#` detected
- `Enter` to submit, `Esc` to cancel
- If the input is empty on submit, do nothing (no error flash — just ignore)

**Auto-categorization keywords**:
- `#research`, `#r` → Research
- `#design`, `#d` → Design
- `#test`, `#t` → Testing
- `#review`, `#rev`, `#rv` → Review
- `#discovery` → Discovery
- `#handoff`, `#handover` → Handoff
- `#admin`, `#a` → Admin

**Time estimate parsing**: `2h`, `30m`, `2.5h` appended anywhere in input

**Due date parsing**: `tomorrow`, `next monday`, `mar 25`, `+3d` (3 days from now)

### 4.2 Kanban Board

**Columns**: Backlog | In Progress | In Review | Approved
- Each column shows a task count badge
- Empty columns show a subtle dot-grid placeholder
- Horizontal scroll on tablet/mobile

**Drag & Drop**:
- Mouse: click-and-hold (150ms delay) to pick up, release to drop
- Visual feedback: card lifts with shadow, drop zone highlights with 2px indigo border
- Cross-column: smooth insertion animation
- Keyboard: `Space` to pick up selected card, `J/K` or arrows to move, `Space` to drop, `Esc` to cancel
- Touch: long-press (300ms) to activate drag mode

**Add Task Inline**:
- Each column has a "+ Add task" button at the bottom
- Clicking opens an inline input within the column
- `Enter` creates the task in that column with current filter's category
- `Esc` cancels

### 4.3 Task Cards

**Card anatomy**:
```
┌──────────────────────────────────────────┐
│ [Category Chip]              [Priority] │
│ Task title goes here and can wrap to     │
│ multiple lines if needed...              │
│                                          │
│ [📎] [2h] [Due: Mar 25]                 │
│ [subtask progress]                       │
└──────────────────────────────────────────┘
```

**States**:
- **Default**: White card, warm shadow, 12px radius
- **Hover**: Subtle lift (translateY -2px), slightly stronger shadow
- **Dragging**: Elevated shadow, 3deg rotation, 50% opacity at origin
- **Selected**: 2px indigo border
- **Completed**: Strikethrough title, muted colors, moved to "Completed" section (collapsible)

**Interactions**:
- **Click**: Opens task detail panel (slide-in from right)
- **Double-click title**: Inline edit
- **Checkbox click**: Mark complete with micro-animation, then fade to completed section after 2 seconds
- **Right-click**: Context menu (Edit, Duplicate, Move to column, Delete)
- **Hover actions** (appear on card hover): Quick-complete checkbox, open detail

### 4.4 Task Detail Panel

**Layout**: Slide-in panel from right, 400px wide (desktop), full-width bottom sheet (mobile)
**Close**: Click outside, press `Esc`, or click X button
**Save**: Auto-save on every field change — no save button needed

**Sections**:
1. **Title** — Large, editable inline (click to edit, blur to save)
2. **Category** — Dropdown selector with colored chips
3. **Iteration State** — Column selector (Backlog → In Progress → In Review → Approved)
4. **Priority** — 1-5 star rating, visual indicators
5. **Due Date** — Date picker with natural language support
6. **Time Estimate** — Hours/minutes input
7. **Tags** — Multi-select from pre-built list + create custom
8. **Attachments** — URL input (Figma links, prototype URLs) — each shows as a clickable chip with favicon
9. **Description** — Rich text area (bold, italic, links only — no full WYSIWYG)
10. **Subtasks** — Checklist items within the task
11. **Metadata** — Created date, last modified, time tracked (computed from subtask completion)

### 4.5 Filter Bar

- Category chips are toggle buttons — click to filter by that category
- Active filters show with filled background and checkmark
- Multiple category filters combine with OR logic
- "All" deselects all category filters
- Tags dropdown: multi-select tag filter
- Filter state persists in URL hash (sharable filtered views)
- Active filter chips are removable (X button)

### 4.6 Search

- `Cmd/Ctrl+F` or click search icon to focus
- Searches: task title, description, tags, attachment URLs
- Results highlight matching text
- Press `Enter` to jump to first result in kanban view
- `Esc` to clear and close

### 4.7 Settings

Accessible via gear icon in header. Settings panel (not modal, full-screen takeover or slide-in):
- **Appearance**: Theme toggle (Light / Dark / System)
- **Columns**: Rename columns, reorder columns, show/hide columns
- **Categories**: Manage category list (add, edit color, delete)
- **Tags**: Manage tag list (add, edit color, delete)
- **Data**: Export to JSON, Import from JSON, Clear all data (with confirmation)
- **Keyboard shortcuts**: Reference page showing all shortcuts

### 4.8 Data Persistence

- **localStorage** with key `designflow_tasks`
- **Auto-save**: Every state change triggers debounced save (300ms)
- **Export**: Download as `designflow-backup-YYYY-MM-DD.json`
- **Import**: File picker, validates JSON structure, merges or replaces
- **No account required** — zero-friction first use

### 4.9 Keyboard Shortcuts

| Shortcut | Action |
|----------|--------|
| `Cmd/Ctrl + K` | Focus quick add bar |
| `Cmd/Ctrl + F` | Focus search |
| `Cmd/Ctrl + /` | Show keyboard shortcuts overlay |
| `J` / `K` | Navigate down/up through tasks |
| `H` / `L` | Navigate left/right through columns |
| `Enter` | Open selected task detail |
| `Space` | Toggle task complete |
| `N` | Add new task in current column |
| `D` | Set due date on selected task |
| `E` | Edit title inline |
| `1-5` | Set priority |
| `Esc` | Close panel / deselect |
| `Cmd/Ctrl + Shift + D` | Toggle dark mode |

---

## 5. Component Inventory

### 5.1 AppShell
- Container for header, content, status bar
- Handles theme context (light/dark/system)
- Manages global keyboard shortcut listener

### 5.2 Header
- Fixed top, 56px height
- Logo + app name (left)
- View toggle buttons: Board | List | Focus (center-right)
- Search icon, Settings icon, Theme toggle (right)
- Subtle bottom border, no heavy shadow

### 5.3 QuickAddBar
- Sticky below header
- Left icon: `+` or search icon
- Input: placeholder "Add a task... (Cmd+K)"
- Right: Category quick-select dropdown
- Full width with max-width container
- Warm cream background with subtle border

### 5.4 FilterBar
- Below QuickAddBar
- Horizontal scrollable row of category chips
- "+ Tag" button opens tag dropdown
- Active filters as removable chips with X button
- Collapse/expand toggle for mobile

### 5.5 KanbanBoard
- Flex row of KanbanColumn components
- Horizontal scroll with snap points
- Empty state per column

### 5.6 KanbanColumn
- Header: column name + task count badge
- Droppable zone (highlights on drag-over)
- Scrollable task list
- "+ Add task" button at bottom
- States: default, drag-over (indigo border + light indigo bg)

### 5.7 TaskCard
- Compact card within column
- States: default, hover (lift), dragging (elevated + rotated), selected (indigo border), completed (strikethrough + muted)
- Click: opens detail panel
- Checkbox: completes task
- Priority indicator: 1-5 dots (filled for active)
- Category chip, time estimate, due date, attachment icon
- Subtask progress bar (thin, colored)

### 5.8 TaskDetailPanel
- Slide-in from right, 400px
- Overlay behind (semi-transparent backdrop, click to close)
- Sections: title, metadata grid, description, attachments, subtasks
- Auto-save on all field changes
- Close button (X) top-right

### 5.9 CategoryChip
- Small pill with category color + category name
- Sizes: sm (for cards), md (for filters)
- States: default, selected (filled bg)

### 5.10 TagChip
- Small rounded tag, gray by default or custom color
- Removable (X button) in filter context
- Selectable in task detail context

### 5.11 EmptyState
- Centered in column or full page
- Simple line-art SVG illustration
- Heading + subtext + CTA button
- Context-aware copy

### 5.12 SettingsPanel
- Slide-in from right or full-screen on mobile
- Grouped sections with dividers
- Toggle switches, dropdowns, buttons
- Data management at bottom

### 5.13 KeyboardShortcutsOverlay
- Modal overlay showing all shortcuts
- Grouped by category
- `Esc` to close

---

## 6. Technical Approach

### 6.1 Tech Stack

| Layer | Technology | Rationale |
|-------|-----------|-----------|
| Build | Vite 5 | Fast HMR, modern ESM, TypeScript support |
| UI Framework | React 18 | Component model matches design system, broad ecosystem |
| Styling | Tailwind CSS 3 | Design token system, responsive utilities, dark mode |
| State | Zustand | Minimal boilerplate, localStorage sync built-in |
| Drag & Drop | @dnd-kit | Best React DnD library — accessible, performant, touch-friendly |
| Icons | Lucide React | Consistent outlined style, tree-shakeable |
| Date handling | date-fns | Lightweight, tree-shakeable date utilities |
| ID generation | nanoid | Tiny, URL-safe unique IDs |

### 6.2 Project Structure

```
src/
├── main.tsx                 # Entry point
├── App.tsx                  # Root component, theme provider
├── index.css                # Tailwind directives + CSS variables
├── components/
│   ├── layout/
│   │   ├── AppShell.tsx
│   │   ├── Header.tsx
│   │   ├── StatusBar.tsx
│   │   └── FilterBar.tsx
│   ├── board/
│   │   ├── KanbanBoard.tsx
│   │   ├── KanbanColumn.tsx
│   │   ├── TaskCard.tsx
│   │   └── InlineAddTask.tsx
│   ├── task/
│   │   ├── TaskDetailPanel.tsx
│   │   ├── TaskTitle.tsx
│   │   ├── TaskMetadata.tsx
│   │   ├── TaskAttachments.tsx
│   │   └── TaskSubtasks.tsx
│   ├── shared/
│   │   ├── CategoryChip.tsx
│   │   ├── TagChip.tsx
│   │   ├── PriorityIndicator.tsx
│   │   ├── EmptyState.tsx
│   │   └── QuickAddBar.tsx
│   ├── settings/
│   │   ├── SettingsPanel.tsx
│   │   └── KeyboardShortcutsOverlay.tsx
│   └── modals/
│       └── ConfirmDialog.tsx
├── store/
│   ├── taskStore.ts         # Zustand store for tasks
│   ├── uiStore.ts           # Zustand store for UI state (view, filters, panels)
│   └── settingsStore.ts    # Zustand store for user settings
├── hooks/
│   ├── useKeyboardShortcuts.ts
│   ├── useDragAndDrop.ts
│   ├── useLocalStorage.ts
│   └── useNaturalLanguage.ts
├── utils/
│   ├── naturalLanguage.ts   # Parse dates, categories, estimates from text
│   ├── taskHelpers.ts      # Sort, filter, group tasks
│   └── exportImport.ts     # JSON export/import logic
├── types/
│   └── index.ts             # Task, Category, Tag, Settings types
└── data/
    └── defaults.ts          # Default categories, tags, column names
```

### 6.3 Data Model

```typescript
// types/index.ts

type IterationState = 'backlog' | 'in_progress' | 'in_review' | 'approved';
type Priority = 1 | 2 | 3 | 4 | 5;

interface Attachment {
  id: string;
  type: 'url'; // Future: 'file' | 'image'
  url: string;
  label?: string; // Extracted from URL or user-defined
  favicon?: string; // Google favicon service
}

interface Subtask {
  id: string;
  title: string;
  completed: boolean;
}

interface Task {
  id: string;
  title: string;
  description: string;
  category: CategoryType;
  iterationState: IterationState;
  priority: Priority;
  dueDate: string | null; // ISO date string
  timeEstimate: number | null; // minutes
  tags: string[];
  attachments: Attachment[];
  subtasks: Subtask[];
  createdAt: string; // ISO timestamp
  updatedAt: string; // ISO timestamp
  completedAt: string | null; // ISO timestamp
  position: number; // Sort order within column
}

interface Category {
  id: string;
  name: string;
  colorLight: string;
  colorDark: string;
  icon?: string; // Lucide icon name
  keywords: string[]; // For auto-categorization
}

interface Tag {
  id: string;
  name: string;
  color: string;
}

interface Settings {
  theme: 'light' | 'dark' | 'system';
  columnOrder: IterationState[];
  customCategories: Category[];
  customTags: Tag[];
}

type CategoryType = 'research' | 'design' | 'testing' | 'review' | 'discovery' | 'handoff' | 'admin';
```

### 6.4 State Architecture

**Three Zustand stores**:

1. **taskStore** — Source of truth for all tasks
   - `tasks: Task[]`
   - `addTask`, `updateTask`, `deleteTask`, `moveTask`, `reorderTasks`
   - Synced to localStorage via `persist` middleware
   - Computed: tasks grouped by iterationState, filtered by active filters

2. **uiStore** — Ephemeral UI state
   - `view: 'board' | 'list' | 'focus'`
   - `activeFilters: { categories: CategoryType[], tags: string[] }`
   - `selectedTaskId: string | null`
   - `detailPanelOpen: boolean`
   - `settingsOpen: boolean`
   - `searchQuery: string`
   - NOT persisted — resets on page load

3. **settingsStore** — User preferences
   - `settings: Settings`
   - `updateTheme`, `updateColumns`, `updateCategories`, `updateTags`
   - Synced to localStorage via `persist` middleware

### 6.5 Key Implementation Notes

1. **Drag & Drop**: Use `@dnd-kit/core` with `@dnd-kit/sortable`. The `DndContext` wraps the board, each column is a `SortableContext`, and each card is a `useSortable` item. Custom collision detection for cross-column drops.

2. **Natural Language Parsing**: Regex-based parser that extracts `#category`, `@dueDate`, and `Xh/Xm` time estimates from the quick-add input string. date-fns for date normalization.

3. **Dark Mode**: Tailwind's `dark:` variant with `class` strategy. Theme context from settingsStore sets `document.documentElement.classList`. CSS custom properties for colors.

4. **Accessibility**: All interactive elements have `aria-label`, `role`, and `tabIndex`. Drag-and-drop has full keyboard alternative. Focus trap in modal/panel contexts. ARIA live regions for dynamic content updates.

5. **Performance**: React `useMemo` for filtered/sorted task lists. Virtual scrolling not needed for MVP (assume <500 tasks). Zustand selectors to prevent unnecessary re-renders.

6. **Auto-save**: Debounced 300ms write to localStorage. Show "Saving..." → "Saved" indicator in status bar.

---

## 7. Edge Cases & Error Handling

| Scenario | Behavior |
|----------|----------|
| localStorage full | Show toast: "Storage full. Export your tasks to free up space." with export button |
| Import invalid JSON | Show inline error below file picker: "Invalid file format. Please select a DesignFlow backup file." |
| Import with existing tasks | Prompt: "Merge with existing tasks?" or "Replace all tasks?" |
| Delete task with subtasks | No confirmation needed for MVP. Subtasks are deleted with parent. |
| Clear all data | Confirmation dialog with explicit "Type DELETE to confirm" |
| Duplicate task titles | Allowed — tasks are identified by ID, not title |
| Empty quick-add submit | Ignored silently |
| Very long task titles | Truncate with ellipsis on card (3 lines max), full text in detail panel |
| No tasks in column | Show warm empty state with inline "+ Add" CTA |
| URL attachment without protocol | Auto-prepend `https://` |
| Keyboard shortcut conflict | User's system shortcuts take precedence; no override |

---

## 8. Scope for MVP (Phase 1)

The MVP focuses on the core kanban experience — everything needed for a designer to replace their current workaround tools.

**In MVP scope**:
- Kanban board with 4 columns (Backlog, In Progress, In Review, Approved)
- Quick add bar with natural language parsing (#category, @date, Xh time)
- Drag-and-drop reordering within and across columns
- Task cards with category chip, priority, due date, time estimate
- Task detail panel (slide-in) with all editable fields
- Subtasks within tasks
- URL attachments (Figma links, prototype URLs)
- Filter bar by category
- Light + Dark mode (toggle)
- Search (task title only)
- localStorage persistence
- Export to JSON

**Out of MVP scope (Phase 2+)**:
- List view, Focus view
- Tag management and tag filtering
- Column customization (rename, reorder, show/hide)
- Calendar view
- Import from JSON
- Keyboard shortcuts overlay
- Settings panel (beyond theme toggle)
- Rich text description editor
- Time tracking
- Mobile-specific optimizations
- Touch drag-and-drop polish

---

**Status: READY_FOR_DEVELOPMENT**

*This specification is the source of truth. All implementation decisions must trace back to this document. When in doubt, refer to research in `.plans/research/`. When a new requirement emerges, update this document first.*
