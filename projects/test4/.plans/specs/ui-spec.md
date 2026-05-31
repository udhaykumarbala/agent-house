# UI Specification: DesignFlow

Based on research in: `.plans/research/ui-research.md` and `.plans/research/ux-research.md`

---

## Design System

### Color Palette

**IMPORTANT**: These colors are UNIQUE to DesignFlow. They were chosen to differentiate from generic todo apps (blue/red/green SaaS) and competitor tools (Linear purple, Things 3 teal). The palette evokes "Editorial Design Studio" — warm, confident, and premium.

#### Light Theme (Default)

| Name | Hex | RGB | Usage |
|------|-----|-----|-------|
| Primary | `#3730A3` | rgb(55,48,163) | Main CTAs, active states, links, kanban column headers |
| Primary Hover | `#312E81` | rgb(49,46,129) | Hover state for primary elements |
| Primary Light | `#EEF2FF` | rgb(238,242,255) | Selected backgrounds, active chips |
| Accent | `#F97316` | rgb(249,115,22) | Highlights, notifications, quick-add focus ring |
| Accent Light | `#FFF7ED` | rgb(255,247,237) | Accent chip backgrounds |
| Background | `#FAF8F5` | rgb(250,248,245) | App background — warm cream, not cold white |
| Surface | `#FFFFFF` | rgb(255,255,255) | Cards, modals, elevated panels |
| Surface Hover | `#FAFAF8` | rgb(250,250,248) | Card hover state |
| Border | `#E7E5E4` | rgb(231,229,228) | Dividers, input borders, card outlines |
| Border Focus | `#3730A3` | rgb(55,48,163) | Focus rings on inputs |
| Text Primary | `#1C1917` | rgb(28,25,23) | Main text — warm near-black |
| Text Secondary | `#78716C` | rgb(120,113,108) | Muted labels, metadata |
| Text Tertiary | `#A8A29E` | rgb(168,162,158) | Placeholder text, disabled states |
| Success | `#4D7C0F` | rgb(77,124,15) | Completed states, success feedback |
| Success Light | `#F7F8E8` | rgb(247,248,232) | Completed task background |
| Error | `#DC2626` | rgb(220,38,38) | Error states, destructive actions |
| Warning | `#D97706` | rgb(217,119,6) | Warning states, due-soon indicators |

#### Dark Theme

| Name | Hex | RGB | Usage |
|------|-----|-----|-------|
| Primary | `#818CF8` | rgb(129,140,248) | Main CTAs, active states (lighter indigo for dark bg) |
| Primary Hover | `#6366F1` | rgb(99,102,241) | Hover state |
| Primary Light | `#312E81` | rgb(49,46,129) | Selected backgrounds |
| Accent | `#FB923C` | rgb(251,146,60) | Highlights, peach warmth — softer than coral on dark |
| Accent Light | `#292524` | rgb(41,37,36) | Accent chip backgrounds |
| Background | `#0C0A09` | rgb(12,10,9) | App background — deep warm black |
| Surface | `#1C1917` | rgb(28,25,23) | Cards, panels |
| Surface Hover | `#292524` | rgb(41,37,36) | Card hover |
| Border | `#292524` | rgb(41,37,36) | Dividers, borders |
| Border Focus | `#818CF8` | rgb(129,140,248) | Focus rings |
| Text Primary | `#FAFAF9` | rgb(250,250,249) | Main text |
| Text Secondary | `#A8A29E` | rgb(168,162,158) | Muted labels |
| Text Tertiary | `#78716C` | rgb(120,113,108) | Placeholder text |
| Success | `#84CC16` | rgb(132,204,22) | Lighter lime for dark bg |
| Success Light | `#1C1917` | rgb(28,25,23) | Completed task bg |
| Error | `#EF4444` | rgb(239,68,68) | Slightly brighter for dark bg |
| Warning | `#FBBF24` | rgb(251,191,36) | Amber for dark bg |

#### Category Colors (Both Themes)

These colors represent UX design process phases. They use muted, harmonious tones — not harsh primaries.

| Category | Light Hex | Dark Hex | Usage |
|----------|-----------|----------|-------|
| Research | `#4D7C0F` | `#84CC16` | Sage green — investigative, grounded |
| Design | `#3730A3` | `#818CF8` | Indigo — creative depth, primary brand |
| Testing | `#D97706` | `#FBBF24` | Amber — warm attention, usability focus |
| Review | `#BE185D` | `#F472B6` | Rose — critical but warm, stakeholder review |
| Discovery | `#0D9488` | `#2DD4BF` | Teal — exploratory, early-phase calm |
| Handoff | `#475569` | `#94A3B8` | Slate — technical, developer-oriented |
| Admin | `#6B7280` | `#9CA3AF` | Gray — personal, miscellaneous |

**Why These Colors**: The indigo primary creates a sophisticated brand anchor while coral/peach accents add warmth without the stress of red. Category colors use muted, earthy tones (sage, teal, slate) rather than saturated primaries — this feels more like a design studio than a corporate dashboard.

---

### Typography

**Font Family**: Cabinet Grotesk (headings) + Inter (body) — loaded from Google Fonts / CDN

- Cabinet Grotesk signals "designer-made" — used by Framer, Loom, and beloved in the design community
- Inter is the workhorse — excellent legibility at all sizes, industry standard for UI

| Element | Font | Size | Weight | Line Height | Letter Spacing |
|---------|------|------|--------|-------------|----------------|
| App Logo | Cabinet Grotesk | 20px | 700 | 1.2 | -0.02em |
| H1 (Page Title) | Cabinet Grotesk | 28px | 700 | 1.2 | -0.02em |
| H2 (Section) | Cabinet Grotesk | 20px | 600 | 1.3 | -0.01em |
| H3 (Column Header) | Cabinet Grotesk | 14px | 600 | 1.4 | 0 |
| Task Title | Inter | 15px | 500 | 1.4 | 0 |
| Body Text | Inter | 14px | 400 | 1.5 | 0 |
| Small/Meta | Inter | 12px | 400 | 1.4 | 0 |
| Chip/Tag | Inter | 11px | 500 | 1.2 | 0.02em |
| Time Estimate | JetBrains Mono | 12px | 400 | 1.4 | 0 |

**Fallback Stack**:
```css
font-family: 'Cabinet Grotesk', 'Inter', system-ui, -apple-system, sans-serif;
font-family: 'JetBrains Mono', 'SF Mono', 'Consolas', monospace;
```

---

### Spacing Scale

Base: 4px

| Token | Value | Tailwind | Usage |
|-------|-------|----------|-------|
| 0 | 0px | 0 | Reset |
| 0.5 | 2px | 1 | Tight internal padding |
| 1 | 4px | 1 | Micro gaps |
| 2 | 8px | 2 | Small gaps, chip padding |
| 3 | 12px | 3 | Input padding, compact spacing |
| 4 | 16px | 4 | Standard padding |
| 5 | 20px | 5 | Card padding |
| 6 | 24px | 6 | Section spacing |
| 8 | 32px | 8 | Large gaps |
| 10 | 40px | 10 | Major section margins |
| 12 | 48px | 12 | Page margins (desktop) |
| 16 | 64px | 16 | Hero spacing |

---

### Border Radius

| Token | Value | Tailwind | Usage |
|-------|-------|----------|-------|
| none | 0px | 0 | Sharp edges (internal dividers) |
| sm | 4px | rounded | Inputs, checkboxes, small elements |
| md | 8px | rounded-md | Buttons, chips, cards |
| lg | 12px | rounded-lg | Modals, large panels |
| xl | 16px | rounded-xl | Feature cards, overlays |
| full | 9999px | rounded-full | Pills, avatars, circular buttons |

---

### Shadows

**Principle**: Warm-tinted shadows (warm gray, not cool gray). Shadows add depth without coldness.

#### Light Theme Shadows

| Token | Value | Usage |
|-------|-------|-------|
| sm | `0 1px 2px rgba(28,25,23,0.05)` | Subtle elevation, card default |
| md | `0 4px 12px rgba(28,25,23,0.08)` | Cards on hover, dropdowns |
| lg | `0 12px 32px rgba(28,25,23,0.12)` | Modals, drag preview |
| xl | `0 24px 48px rgba(28,25,23,0.16)` | Command palette |

#### Dark Theme Shadows

| Token | Value | Usage |
|-------|-------|-------|
| sm | `0 1px 2px rgba(0,0,0,0.3)` | Subtle elevation |
| md | `0 4px 12px rgba(0,0,0,0.4)` | Cards on hover, dropdowns |
| lg | `0 12px 32px rgba(0,0,0,0.5)` | Modals, drag preview |
| xl | `0 24px 48px rgba(0,0,0,0.6)` | Command palette |

---

### Motion & Animation

**Principle**: Subtle spring/overshoot — playful precision. Designers notice and appreciate good motion; bad motion breaks trust.

| Token | Value | Usage |
|-------|-------|-------|
| fast | 150ms | Checkbox toggles, immediate feedback |
| base | 200ms | Hover states, color transitions |
| slow | 300ms | Panel slides, modal enters |
| spring | cubic-bezier(0.34,1.56,0.64,1) | Card drag release, completion celebration |
| smooth | cubic-bezier(0.4,0,0.2,1) | Panel slides, page transitions |

**Reduced Motion**: All animations respect `prefers-reduced-motion`. When enabled, use opacity fades (200ms) instead of transforms/springs.

---

## Components

### 1. Quick Add Bar

```
┌─────────────────────────────────────────────────────────────┐
│  ⊕  Add a task...  (Cmd+K)                    [🔍 Filter]  │
└─────────────────────────────────────────────────────────────┘
```

**States**:
- Default: Surface bg, border, placeholder text in Text Tertiary
- Focused: Border Focus ring (indigo), accent glow, placeholder fades
- Active (typing): Text Primary, category chip auto-detection shown below
- Disabled: opacity-50, cursor-not-allowed

**Behavior**:
- `Cmd+K` / `Ctrl+K` focuses this bar from anywhere
- Natural language parsing: "Design hero section next Monday 2h" auto-assigns category, date, estimate
- Auto-suggestions appear below as chips (category, date, estimate)
- `Enter` creates task and refocuses bar; `Shift+Enter` for newline in title

### 2. Task Card

```
┌──────────────────────────────────────────────────────┐
│  [Category Chip]                    [Priority dot]   │
│  Task title goes here, can wrap to two lines max    │
│                                                      │
│  [📎] [link icon]          [2h]    [Mar 25]  ☰     │
└──────────────────────────────────────────────────────┘
```

**States**:
- Default: Surface bg, sm shadow, md border-radius
- Hover: md shadow, surface-hover bg, subtle translateY(-1px)
- Dragging: lg shadow, slight rotation (1-2deg), scale(1.02), opacity-90
- Completed: strikethrough title, success-light bg, muted text, checkmark overlay
- Selected: Primary-light bg, primary border (2px)

**Anatomy**:
- Category chip: top-left, 6px padding, rounded-full, category color
- Priority indicator: top-right, 6px colored dot (High=red, Medium=amber, Low=gray)
- Title: Task Title typography, max 2 lines with ellipsis
- Footer row: attachment icon (if has links), time estimate (JetBrains Mono), due date (Small), overflow menu (☰)

### 3. Kanban Column

```
┌────────────────────┐
│  ● In Progress (3) │  ← Column header with count
├────────────────────┤
│  [Task Card]        │
│  [Task Card]        │
│  [+ Add task]       │  ← Inline add at bottom
└────────────────────┘
```

**Header**: H3 typography, colored left border (4px, category color), task count in Text Secondary

**Drop zone**: When dragging over, column bg shifts to primary-light, dashed border appears

**Column widths**: min-width 280px, max-width 320px, equal flex distribution

### 4. Category Chip

```
[● Research]   [● Design]   [● Testing]   [● Review]   [● Discovery]   [● Handoff]
```

- Pill shape: rounded-full, 2px vertical padding, 8px horizontal padding
- Typography: Chip typography (11px, 500 weight)
- Colors: Category-specific (see palette table above), light-mode bg is color at 10% opacity
- Interactive chips (for filters): hover brightens, click toggles filter
- Active filter chip: filled bg (not just outline), checkmark icon

### 5. Filter Bar

```
┌─────────────────────────────────────────────────────────────┐
│  [● Research] [● Design] [● Testing]    Sort: [Due Date ▼] │
│                                        [👤 All Projects]    │
└─────────────────────────────────────────────────────────────┘
```

- Active filters shown as removable chips (× on hover)
- Sort dropdown: Due Date, Priority, Created, Alphabetical
- Project filter: dropdown with project names
- Filter bar collapses to single row on mobile

### 6. Command Palette (Cmd+K)

```
┌─────────────────────────────────────────────────────┐
│  🔍  Search tasks or type a command...               │
├─────────────────────────────────────────────────────┤
│  TASKS                                              │
│    🔍 Conduct 5 user interviews           Research  │
│    🔍 Create high-fidelity mockups       Design    │
│    🔍 Recruit usability test participants Testing  │
│                                                      │
│  COMMANDS                                           │
│    ⊕ Add new task                    Cmd+Enter     │
│    ☐ Show completed tasks                           │
│    🌙 Toggle dark mode                             │
│    ⚙  Open settings                               │
└─────────────────────────────────────────────────────┘
```

- Centered overlay, xl border-radius, lg shadow (xl shadow in dark)
- Backdrop: semi-transparent with blur(8px)
- Search input at top, results grouped below
- Keyboard navigation: ↑↓ to navigate, Enter to select, Esc to close
- 280px min-width, 480px max-width, auto height (max 400px)

### 7. Task Detail Panel

Slide-in panel from right (640px wide on desktop, full-screen on mobile):

```
┌──────────────────────────────────────────────┐
│  ← Back                         [Delete] [⋯] │
├──────────────────────────────────────────────┤
│  [Research ▼]                                │
│  Conduct 5 user interviews                    │
│  ─────────────────────────────────           │
│  Status: [● Not Started ▼]                   │
│  Priority: [High ▼]                          │
│  Due: [Mar 25, 2026]                        │
│  Estimate: [2h]                              │
│  Project: [Checkout Redesign ▼]              │
│  ─────────────────────────────────           │
│  Description                                 │
│  [Rich text area...]                         │
│  ─────────────────────────────────           │
│  Attachments                                 │
│  [🔗 https://figma.com/file/...]            │
│  [🔗 Prototype link]                        │
│  ─────────────────────────────────           │
│  Subtasks (3/5)                              │
│  [ ] Schedule interviews                     │
│  [ ] Prepare interview guide                 │
│  [✓] Send recruitment email                  │
└──────────────────────────────────────────────┘
```

- Auto-save on every change (no Save button)
- Sections collapse/expand with smooth animation
- Attachment: paste Figma URL → auto-detected and shown as linked card

### 8. Empty State

```
┌─────────────────────────────────────────────────────────────┐
│                                                             │
│                    [Line art illustration]                  │
│                                                             │
│               Your design backlog is clear.                  │
│          Add your first task to get started.                 │
│                                                             │
│               [+ Add your first task]                        │
│                                                             │
│               Tip: Press Cmd+K to add a task quickly         │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

- Simple line-art illustration (CSS-drawn or SVG — no external images)
- Centered, generous vertical padding (96px+)
- Warm, encouraging copy
- Primary CTA button
- Subtle keyboard shortcut hint

### 9. Buttons

```css
Primary:   bg-[primary] text-white rounded-md px-4 py-2
           hover:bg-[primary-hover] active:scale-[0.98]
Secondary: bg-transparent text-[primary] border border-[border] rounded-md px-4 py-2
           hover:bg-[primary-light]
Ghost:     bg-transparent text-[primary] px-4 py-2
           hover:bg-[primary-light] rounded-md
Destructive: bg-[error] text-white rounded-md px-4 py-2
           hover:opacity-90
Icon:      bg-transparent text-[text-secondary] p-2 rounded-full
           hover:bg-[surface-hover] hover:text-[text-primary]
```

**Sizes**: sm (py-1.5 px-3, text-sm), md (py-2 px-4, text-sm), lg (py-2.5 px-5, text-base)

**States**: Default → Hover → Active (scale-down) → Disabled (opacity-50, cursor-not-allowed) → Loading (spinner)

### 10. Input Fields

```css
Default:   bg-[surface] border border-[border] rounded-md px-3 py-2
           text-[text-primary] placeholder-[text-tertiary]
           text-sm leading-5
Focus:     border-[primary] ring-2 ring-[primary]/20
Error:     border-[error] ring-2 ring-[error]/20
Disabled:  bg-[surface-hover] text-[text-tertiary] cursor-not-allowed
```

**Textarea**: Same base, min-height 120px, resize-y

**Select/Dropdown**: Same as input, with chevron-down icon, dropdown uses lg border-radius

### 11. Progress Ring

```
      ┌───┐
    ╱   75%  ╲
    │        │
    │   ◯    │
    ╲        ╱
      └───┘
```

- SVG-based circular progress
- Track: border color at 20% opacity
- Fill: Primary color, rounded stroke-linecap
- Center: percentage text (H3 typography) or fraction ("3/5")
- Size variants: sm (32px), md (48px), lg (64px)

### 12. Badge / Pill

```css
Status badge: bg-[status-color]/10 text-[status-color]
               rounded-full px-2.5 py-1 text-xs font-medium
Priority:      6px circle, solid status color
Count:         bg-[primary] text-white text-xs font-semibold
               rounded-full min-w-[20px] h-5 px-1.5 flex items-center justify-center
```

---

## Responsive Breakpoints

| Name | Width | Layout |
|------|-------|--------|
| mobile | < 640px | Single column, horizontal scroll for kanban, full-width cards |
| tablet | 640-1024px | 2-column kanban, collapsible sidebar |
| desktop | > 1024px | Full 4-column kanban, visible header |

### Mobile Adaptations:
- Quick Add bar: full-width, 48px height, bottom-positioned option
- Task cards: full-width, reduced padding
- Task detail: full-screen panel (replaces content)
- Filter bar: horizontal scrollable chips
- Kanban: horizontal scroll, snap to column

### Tablet Adaptations:
- 2 kanban columns visible at once
- Quick add bar: inline with header
- Task detail: 480px side panel

### Desktop:
- 4 kanban columns
- Full command palette
- Persistent filter bar
- Optional sidebar (project list)

---

## Icon System

**Library**: Lucide React (outline icons, 1.5px stroke)

**Key icons used**:
- `Plus` — Add task, new item
- `Search` — Search, filter
- `Filter` — Filter bar
- `Check` / `CheckCircle2` — Complete task
- `Circle` — Unchecked state
- `GripVertical` — Drag handle
- `Calendar` — Due date
- `Clock` — Time estimate
- `Link` — Attachments
- `Tag` — Labels/tags
- `ChevronDown` — Dropdowns
- `MoreHorizontal` — Overflow menu
- `Sun` / `Moon` — Theme toggle
- `Settings` — Settings
- `Command` — Keyboard shortcut hint
- `Figma` — Figma attachment (custom SVG)
- `X` — Close, remove
- `ArrowLeft` — Back navigation
- `LayoutGrid` — Kanban view
- `List` — List view
- `CalendarDays` — Calendar view

**Icon sizing**: 16px for inline (within text), 20px for standalone actions, 24px for navigation items.

---

## Accessibility

| Requirement | Implementation |
|-------------|----------------|
| Color contrast | 4.5:1 minimum for text, 3:1 for UI components (WCAG AA) |
| Focus indicators | Custom focus ring: 2px primary outline, 2px offset, visible on all interactive elements |
| Keyboard navigation | Full tab order, J/K or arrow navigation in task list, Esc closes panels |
| Screen readers | ARIA labels on all icon buttons, proper heading hierarchy (h1→h2→h3), role="list" for task lists |
| Reduced motion | All animations disabled or simplified when `prefers-reduced-motion: reduce` |
| Touch targets | Minimum 44x44px on all interactive elements |
| Color independence | All color-coded states have text/icon alternatives |
| High contrast mode | Dark theme has higher contrast (text-primary on bg = 14.8:1) |

---

## Component CSS Reference (Tailwind)

```css
/* Button variants */
.btn-primary { @apply bg-indigo-700 text-white rounded-md px-4 py-2 hover:bg-indigo-800 active:scale-[0.98] transition-all duration-150; }
.btn-secondary { @apply bg-transparent text-indigo-700 border border-stone-200 rounded-md px-4 py-2 hover:bg-indigo-50 transition-all duration-150; }
.btn-ghost { @apply bg-transparent text-indigo-700 px-4 py-2 hover:bg-indigo-50 rounded-md transition-all duration-150; }
.btn-icon { @apply bg-transparent text-stone-500 p-2 rounded-full hover:bg-stone-100 hover:text-stone-700 transition-all duration-150; }

/* Card */
.card { @apply bg-white rounded-lg shadow-sm border border-stone-100 p-5 hover:shadow-md hover:-translate-y-px transition-all duration-200; }
.card-dragging { @apply shadow-lg rotate-1 scale-[1.02] opacity-90; }
.card-completed { @apply bg-green-50 border-green-100; }

/* Input */
.input { @apply bg-white border border-stone-200 rounded-md px-3 py-2 text-stone-800 placeholder-stone-400 focus:border-indigo-700 focus:ring-2 focus:ring-indigo-700/20 outline-none transition-all duration-150; }
.input-error { @apply border-red-600 ring-2 ring-red-600/20; }

/* Category chips */
.chip { @apply inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium transition-colors duration-150; }
.chip-research { @apply bg-amber-100 text-amber-800; }
.chip-design { @apply bg-indigo-100 text-indigo-800; }
.chip-testing { @apply bg-orange-100 text-orange-800; }
.chip-review { @apply bg-pink-100 text-pink-800; }

/* Shadows (CSS variables for theme switching) */
.shadow-warm { @apply shadow-[0_4px_12px_rgba(28,25,23,0.08)]; }
.shadow-warm-lg { @apply shadow-[0_12px_32px_rgba(28,25,23,0.12)]; }
```

---

Status: READY_FOR_REVIEW
