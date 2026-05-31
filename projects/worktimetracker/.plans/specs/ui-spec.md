# UI Specification: WorkTime Tracker

Based on research in: .plans/research/ui-research.md
Styling wireframes from: .plans/research/ux-research.md

## Design Philosophy

**"Warm Precision"** — Clean, structured layouts with warm undertones. The app should feel like a well-crafted watch: precise, reliable, and quietly premium. Dark-mode-first with deep indigo foundations and warm amber accents that convey focused productivity.

## Design System

### Color Palette (UNIQUE — Deep Indigo + Warm Amber)

**IMPORTANT**: This palette was chosen to differentiate from the sea of blue (Clockify, Time Doctor), green (Hubstaff), pink (Toggl), and teal (Timely, Rize) time-tracking apps. No major competitor uses indigo + amber.

#### Dark Mode (Primary)

| Name | Hex | RGB | Usage |
|------|-----|-----|-------|
| Primary | #5746B2 | rgb(87, 70, 178) | Brand elements, active nav, selected states |
| Primary Hover | #483A96 | rgb(72, 58, 150) | Hover on primary elements |
| Primary Muted | rgba(87, 70, 178, 0.15) | — | Subtle highlights, active backgrounds |
| Accent | #CF8A2E | rgb(207, 138, 46) | CTAs (Clock In/Out), key actions, timer display |
| Accent Hover | #B87724 | rgb(184, 119, 36) | Hover on accent elements |
| Accent Muted | rgba(207, 138, 46, 0.12) | — | Accent background tints |
| Background | #0D0A1A | rgb(13, 10, 26) | App background, base layer |
| Surface | #17132B | rgb(23, 19, 43) | Cards, sidebar, panels |
| Surface Raised | #201B38 | rgb(32, 27, 56) | Dropdowns, modals, tooltips, hover rows |
| Border | rgba(91, 70, 178, 0.15) | — | Subtle borders, dividers |
| Border Strong | rgba(91, 70, 178, 0.30) | — | Input borders, focused dividers |
| Text Primary | #ECE9F5 | rgb(236, 233, 245) | Headings, body text, primary labels |
| Text Secondary | #8C83A8 | rgb(140, 131, 168) | Muted text, descriptions, timestamps |
| Text Disabled | #564F6D | rgb(86, 79, 109) | Disabled labels, placeholders |
| Success | #3DB87A | rgb(61, 184, 122) | Clocked-in status, success toasts |
| Error | #D94B4B | rgb(217, 75, 75) | Validation errors, destructive actions |
| Warning | #DFA23E | rgb(223, 162, 62) | On-break status, caution states |
| Info | #5B8AD4 | rgb(91, 138, 212) | Informational badges, links |

#### Light Mode (Secondary)

| Name | Hex | RGB | Usage |
|------|-----|-----|-------|
| Primary | #5746B2 | rgb(87, 70, 178) | Same brand color across modes |
| Primary Hover | #483A96 | rgb(72, 58, 150) | Hover states |
| Primary Muted | rgba(87, 70, 178, 0.08) | — | Subtle highlights |
| Accent | #B87724 | rgb(184, 119, 36) | CTAs — slightly deeper amber for contrast on light |
| Background | #FFFBF5 | rgb(255, 251, 245) | Warm off-white base |
| Surface | #FFFFFF | rgb(255, 255, 255) | Cards, panels |
| Surface Raised | #F7F4EE | rgb(247, 244, 238) | Hover rows, dropdowns |
| Border | rgba(26, 20, 51, 0.10) | — | Subtle borders |
| Border Strong | rgba(26, 20, 51, 0.20) | — | Input borders |
| Text Primary | #1A1433 | rgb(26, 20, 51) | Deep indigo-black text |
| Text Secondary | #5E5676 | rgb(94, 86, 118) | Muted text |
| Text Disabled | #A9A2BC | rgb(169, 162, 188) | Disabled/placeholder |

**Why These Colors**: Indigo conveys depth, intelligence, and focus — it's the color of deep work. Amber/copper conveys energy, warmth, and productivity — like the glow of a desk lamp during a focused session. Together they evoke "quiet craftsmanship," differentiating from the cold blues and clinical greens of every other time tracker. The warm off-white light mode (#FFFBF5) avoids the sterile #F3F4F6 gray that plagues generic dashboards.

### Typography

**Heading Font**: Manrope (Google Fonts) — Semi-rounded geometric sans-serif with warmth and personality
**Body/UI Font**: Inter (Google Fonts) — Unmatched legibility, excellent tabular numbers for time displays
**Timer Display**: Inter with `font-variant-numeric: tabular-nums` — Digits align perfectly as the timer counts

| Element | Font | Size | Weight | Line Height | Letter Spacing |
|---------|------|------|--------|-------------|----------------|
| H1 (Page Title) | Manrope | 32px | 700 | 1.2 | -0.02em |
| H2 (Section Title) | Manrope | 24px | 600 | 1.3 | -0.01em |
| H3 (Card Title) | Manrope | 20px | 600 | 1.4 | 0 |
| H4 (Subsection) | Manrope | 16px | 600 | 1.4 | 0 |
| Body | Inter | 16px | 400 | 1.5 | 0 |
| Body Small | Inter | 14px | 400 | 1.5 | 0 |
| Caption | Inter | 12px | 400 | 1.4 | 0.01em |
| Label | Inter | 14px | 500 | 1.0 | 0.01em |
| Timer Display | Inter | 48px | 700 | 1.0 | 0.02em |
| Timer Compact | Inter | 20px | 600 | 1.0 | 0.02em |
| Table Header | Inter | 12px | 600 | 1.0 | 0.05em |
| Table Cell | Inter | 14px | 400 | 1.5 | 0 |
| Button | Inter | 14px | 500 | 1.0 | 0.01em |

### Spacing Scale

Base unit: 4px

| Token | Value | Usage |
|-------|-------|-------|
| xs | 4px | Icon-to-label gap, tight element spacing |
| sm | 8px | Intra-component spacing, input padding-y |
| md | 16px | Standard card padding, component gaps |
| lg | 24px | Section spacing, card padding for large cards |
| xl | 32px | Between major sections |
| 2xl | 48px | Page margins, major separators |
| 3xl | 64px | Top-level page padding on desktop |

### Border Radius

| Token | Value | Usage |
|-------|-------|-------|
| sm | 4px | Small badges, tags, chips |
| md | 8px | Buttons, inputs, small cards |
| lg | 12px | Cards, panels, dropdowns |
| xl | 16px | Modals, large containers |
| full | 9999px | Status dots, avatars, pill badges |

### Shadows (Warm-Tinted)

Shadows use indigo tint rather than pure black, adding subtle warmth to elevation.

| Token | Value | Usage |
|-------|-------|-------|
| sm | 0 1px 2px rgba(13, 10, 26, 0.15) | Subtle lift — buttons, badges |
| md | 0 4px 8px rgba(13, 10, 26, 0.20) | Cards, sidebar on mobile |
| lg | 0 12px 24px rgba(13, 10, 26, 0.30) | Modals, dropdowns, popovers |
| glow-primary | 0 0 20px rgba(87, 70, 178, 0.25) | Focus rings, active primary element |
| glow-accent | 0 0 20px rgba(207, 138, 46, 0.25) | Active timer glow, accent emphasis |

### Iconography

- **Icon Set**: Lucide Icons (outlined, 1.5px stroke)
- **Sizes**: 16px (inline/small), 20px (standard UI), 24px (navigation/headings)
- **Color**: Inherits text color (Text Secondary by default, Text Primary on hover/active)
- **Key icons**: Clock (timer), Users (employees), FolderKanban (projects), BarChart3 (reports), Play/Square (start/stop), Download (export), Plus (add), Search, Settings

## Components

### Button

```css
/* Primary — used for main CTAs (not clock in/out) */
.btn-primary {
  background: #5746B2;
  color: #ECE9F5;
  font: 500 14px/1 'Inter', sans-serif;
  padding: 10px 16px;
  border-radius: 8px;
  transition: all 200ms ease;
}
.btn-primary:hover {
  background: #483A96;
  box-shadow: 0 0 20px rgba(87, 70, 178, 0.25);
}

/* Accent — used for Clock In/Out and key actions */
.btn-accent {
  background: #CF8A2E;
  color: #0D0A1A;
  font: 600 14px/1 'Inter', sans-serif;
  padding: 10px 20px;
  border-radius: 8px;
  transition: all 200ms ease;
}
.btn-accent:hover {
  background: #B87724;
  box-shadow: 0 0 20px rgba(207, 138, 46, 0.25);
}

/* Secondary — bordered, for less prominent actions */
.btn-secondary {
  background: transparent;
  color: #ECE9F5;
  border: 1px solid rgba(91, 70, 178, 0.30);
  padding: 10px 16px;
  border-radius: 8px;
}
.btn-secondary:hover {
  background: #201B38;
  border-color: rgba(91, 70, 178, 0.50);
}

/* Ghost — text-only for tertiary actions */
.btn-ghost {
  background: transparent;
  color: #8C83A8;
  padding: 10px 16px;
  border-radius: 8px;
}
.btn-ghost:hover {
  color: #ECE9F5;
  background: rgba(87, 70, 178, 0.10);
}

/* Disabled state — all variants */
.btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
  box-shadow: none;
}
```

**Sizes**:
| Size | Padding | Font | Min Height |
|------|---------|------|------------|
| sm | 6px 12px | 12px | 32px |
| md (default) | 10px 16px | 14px | 40px |
| lg | 12px 24px | 16px | 48px |

### Input

```css
.input {
  background: #17132B;
  color: #ECE9F5;
  border: 1px solid rgba(91, 70, 178, 0.15);
  border-radius: 8px;
  padding: 10px 12px;
  font: 400 14px/1.5 'Inter', sans-serif;
  transition: border-color 200ms ease, box-shadow 200ms ease;
}
.input::placeholder {
  color: #564F6D;
}
.input:focus {
  border-color: #5746B2;
  box-shadow: 0 0 0 3px rgba(87, 70, 178, 0.15);
  outline: none;
}
.input:error {
  border-color: #D94B4B;
  box-shadow: 0 0 0 3px rgba(217, 75, 75, 0.15);
}
```

### Select / Dropdown

```css
.select {
  /* Same base as .input */
  background: #17132B;
  border: 1px solid rgba(91, 70, 178, 0.15);
  border-radius: 8px;
  padding: 10px 12px;
}
.select-dropdown {
  background: #201B38;
  border: 1px solid rgba(91, 70, 178, 0.15);
  border-radius: 8px;
  box-shadow: 0 12px 24px rgba(13, 10, 26, 0.30);
}
.select-option:hover {
  background: rgba(87, 70, 178, 0.15);
}
.select-option.selected {
  background: rgba(87, 70, 178, 0.20);
  color: #CF8A2E;
}
```

### Card

```css
.card {
  background: #17132B;
  border: 1px solid rgba(91, 70, 178, 0.10);
  border-radius: 12px;
  padding: 16px;
}
.card-elevated {
  /* Same as card but with shadow */
  box-shadow: 0 4px 8px rgba(13, 10, 26, 0.20);
}
.card-interactive:hover {
  border-color: rgba(91, 70, 178, 0.25);
  background: #1A1630;
}
```

### Summary Stat Card (Dashboard)

```css
.stat-card {
  background: #17132B;
  border: 1px solid rgba(91, 70, 178, 0.10);
  border-radius: 12px;
  padding: 20px;
}
.stat-card .stat-value {
  font: 700 32px/1.2 'Manrope', sans-serif;
  color: #ECE9F5;
}
.stat-card .stat-label {
  font: 400 12px/1.4 'Inter', sans-serif;
  color: #8C83A8;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}
.stat-card .stat-icon {
  color: #CF8A2E;
  width: 24px;
  height: 24px;
}
```

### Table

```css
.table-header {
  font: 600 12px/1 'Inter', sans-serif;
  color: #8C83A8;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  padding: 12px 16px;
  border-bottom: 1px solid rgba(91, 70, 178, 0.15);
}
.table-row {
  padding: 12px 16px;
  border-bottom: 1px solid rgba(91, 70, 178, 0.08);
  transition: background 150ms ease;
}
.table-row:hover {
  background: #201B38;
}
.table-cell {
  font: 400 14px/1.5 'Inter', sans-serif;
  color: #ECE9F5;
  font-variant-numeric: tabular-nums; /* critical for time/number columns */
}
```

### Status Badges

```css
/* Clocked In */
.badge-active {
  background: rgba(61, 184, 122, 0.12);
  color: #3DB87A;
  font: 500 12px/1 'Inter', sans-serif;
  padding: 4px 10px;
  border-radius: 9999px;
}
/* Clocked Out */
.badge-inactive {
  background: rgba(140, 131, 168, 0.12);
  color: #8C83A8;
  padding: 4px 10px;
  border-radius: 9999px;
}
/* On Break */
.badge-break {
  background: rgba(223, 162, 62, 0.12);
  color: #DFA23E;
  padding: 4px 10px;
  border-radius: 9999px;
}
```

### Status Dots (Inline)

```css
.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 9999px;
}
.status-dot.active { background: #3DB87A; box-shadow: 0 0 6px rgba(61, 184, 122, 0.4); }
.status-dot.inactive { background: #564F6D; }
.status-dot.break { background: #DFA23E; box-shadow: 0 0 6px rgba(223, 162, 62, 0.3); }
```

### Timer Widget

```css
/* Persistent timer bar — sits at top of content area */
.timer-bar {
  background: #17132B;
  border-bottom: 1px solid rgba(91, 70, 178, 0.15);
  padding: 12px 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.timer-display {
  font: 700 48px/1 'Inter', sans-serif;
  font-variant-numeric: tabular-nums;
  color: #CF8A2E;
  letter-spacing: 0.02em;
}
.timer-display.compact {
  font-size: 20px;
  font-weight: 600;
}
.timer-project-label {
  font: 400 14px/1.5 'Inter', sans-serif;
  color: #8C83A8;
}
/* Clock In button — large, amber, glowing when active */
.clock-in-btn {
  background: #CF8A2E;
  color: #0D0A1A;
  font: 600 14px/1 'Inter', sans-serif;
  padding: 12px 28px;
  border-radius: 8px;
  min-height: 48px;
}
.clock-in-btn:hover {
  background: #B87724;
  box-shadow: 0 0 24px rgba(207, 138, 46, 0.35);
}
/* Clock Out — switch to destructive amber-red when timing */
.clock-out-btn {
  background: #D94B4B;
  color: #ECE9F5;
  padding: 12px 28px;
  border-radius: 8px;
  min-height: 48px;
}
```

### Sidebar Navigation

```css
.sidebar {
  width: 240px;
  background: #0D0A1A;
  border-right: 1px solid rgba(91, 70, 178, 0.10);
  padding: 16px 0;
}
.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 20px;
  font: 400 14px/1 'Inter', sans-serif;
  color: #8C83A8;
  border-radius: 0 8px 8px 0;
  margin-right: 8px;
  transition: all 150ms ease;
}
.nav-item:hover {
  color: #ECE9F5;
  background: rgba(87, 70, 178, 0.10);
}
.nav-item.active {
  color: #ECE9F5;
  background: rgba(87, 70, 178, 0.15);
  border-left: 3px solid #CF8A2E;
  font-weight: 500;
}
.nav-item .icon {
  width: 20px;
  height: 20px;
  stroke-width: 1.5;
}
```

### Modal / Dialog

```css
.modal-overlay {
  background: rgba(13, 10, 26, 0.70);
  backdrop-filter: blur(4px);
}
.modal {
  background: #201B38;
  border: 1px solid rgba(91, 70, 178, 0.15);
  border-radius: 16px;
  box-shadow: 0 12px 24px rgba(13, 10, 26, 0.30);
  padding: 24px;
  max-width: 480px;
  width: 100%;
}
.modal-title {
  font: 600 20px/1.4 'Manrope', sans-serif;
  color: #ECE9F5;
}
```

### Toast / Notification

```css
.toast {
  background: #201B38;
  border: 1px solid rgba(91, 70, 178, 0.15);
  border-radius: 8px;
  padding: 12px 16px;
  box-shadow: 0 4px 8px rgba(13, 10, 26, 0.25);
  display: flex;
  align-items: center;
  gap: 12px;
}
.toast-success { border-left: 3px solid #3DB87A; }
.toast-error { border-left: 3px solid #D94B4B; }
.toast-warning { border-left: 3px solid #DFA23E; }
.toast-info { border-left: 3px solid #5B8AD4; }
```

### Date Range Picker

```css
.date-picker {
  background: #201B38;
  border: 1px solid rgba(91, 70, 178, 0.15);
  border-radius: 12px;
  box-shadow: 0 12px 24px rgba(13, 10, 26, 0.30);
}
/* Preset chips: Today, This Week, This Month */
.date-preset {
  background: rgba(87, 70, 178, 0.10);
  color: #8C83A8;
  padding: 6px 14px;
  border-radius: 9999px;
  font: 400 12px/1 'Inter', sans-serif;
}
.date-preset.selected {
  background: #5746B2;
  color: #ECE9F5;
}
```

### Empty State

```css
.empty-state {
  text-align: center;
  padding: 48px 24px;
}
.empty-state .icon {
  color: #564F6D;
  width: 48px;
  height: 48px;
  margin-bottom: 16px;
}
.empty-state .title {
  font: 600 20px/1.4 'Manrope', sans-serif;
  color: #ECE9F5;
  margin-bottom: 8px;
}
.empty-state .description {
  font: 400 14px/1.5 'Inter', sans-serif;
  color: #8C83A8;
  margin-bottom: 24px;
}
/* CTA button follows .btn-accent style */
```

## Data Visualization Palette

For charts and reports, use these project/category colors (distinct, warm-palette-harmonious):

| Slot | Color | Hex | Usage Example |
|------|-------|-----|---------------|
| 1 | Amber | #CF8A2E | Primary project |
| 2 | Indigo | #5746B2 | Secondary project |
| 3 | Coral | #E07B5F | Third project |
| 4 | Teal | #2E9E8F | Fourth project |
| 5 | Rose | #C75A8A | Fifth project |
| 6 | Slate Blue | #5B8AD4 | Sixth project |
| 7 | Olive | #8BA055 | Seventh project |
| 8 | Mauve | #9B6DB0 | Eighth project |

**Chart Guidelines**:
- Bar/area charts: Use solid fills at 80% opacity with full-opacity borders
- Pie/donut charts: Use 1px gap between segments for definition
- Timeline blocks: Use project color at 50% opacity for background, full for border
- Grid lines: Use rgba(91, 70, 178, 0.08) — barely visible
- Axis labels: Text Secondary (#8C83A8), 12px Inter
- Durations in charts: Always "Xh Ym" format for scanability

## Responsive Breakpoints

| Name | Width | Layout | Notes |
|------|-------|--------|-------|
| mobile | < 640px | Single column, bottom nav | Sidebar collapses to hamburger. Timer bar always visible at top. Cards stack vertically. |
| tablet | 640–1024px | Sidebar collapsible, 2-col grid | Sidebar can overlay or collapse. Dashboard uses 2-column stat cards. Tables scroll horizontally. |
| desktop | > 1024px | Full sidebar + content area | Sidebar always visible (240px). Timer bar spans content area. 4-column stat card grid. |

### Mobile-Specific Adaptations
- **Bottom navigation bar** replaces sidebar (5 icons: Dashboard, Employees, Timer, Projects, Reports)
- **Timer**: Full-width bar at top with compact display, large touch target clock in/out button (min 48px height)
- **Tables**: Convert to stacked card layout on mobile (each row becomes a mini-card)
- **Modals**: Slide up from bottom as half-sheets on mobile
- **Touch targets**: Minimum 44x44px for all interactive elements

## Animations & Transitions

| Element | Property | Duration | Easing | Notes |
|---------|----------|----------|--------|-------|
| Button hover | background, box-shadow | 200ms | ease | Subtle color shift + glow |
| Nav item hover | background, color | 150ms | ease | Quick response |
| Page transitions | opacity | 200ms | ease-in-out | Fade between views |
| Modal enter | opacity + translateY | 250ms | ease-out | Fade in + slide up 16px |
| Modal exit | opacity + translateY | 150ms | ease-in | Faster exit |
| Toast enter | translateX | 300ms | ease-out | Slide in from right |
| Timer counting | — | — | — | Number transitions use `font-variant-numeric: tabular-nums` for stable width. No animation on individual digits — just smooth CSS counter updates. |
| Status dot pulse | box-shadow | 2s | ease-in-out infinite | Active status dots pulse gently |
| Sidebar collapse | width | 200ms | ease-in-out | Smooth collapse on tablet |

**Animation Principles**:
- Functional, not decorative — every animation serves a purpose (feedback, orientation, or state change)
- Fast: 150–300ms max. This is a productivity tool, not an art piece.
- No bouncy or springy effects. Ease curves only.
- Reduce motion: Respect `prefers-reduced-motion` — disable all animations except opacity fades.

## Layout Structure

```
┌─────────────────────────────────────────────────┐
│  Sidebar (240px)  │  Timer Bar (persistent)      │
│                   │──────────────────────────────│
│  [Logo]           │                              │
│  ─────            │  Page Content                │
│  Dashboard        │                              │
│  Employees        │  [Summary Cards Grid]        │
│  Time Tracker     │                              │
│  Projects         │  [Main Content Area]         │
│  Reports          │                              │
│                   │                              │
│  ─────            │                              │
│  Settings (gear)  │                              │
└─────────────────────────────────────────────────┘
```

### Sidebar
- Fixed left, full height
- App logo/name at top: "WorkTime" in Manrope 600, 18px
- Nav items with Lucide icons (20px) + labels
- Active item: left amber border (3px), indigo background tint
- Settings gear icon at bottom
- Collapse to icons-only at tablet breakpoint

### Timer Bar
- Sticky top of content area
- Shows: employee selector, project selector, timer display, clock in/out button
- When timer running: amber glow on display, red clock-out button
- When idle: neutral state, amber clock-in button
- Compact mode: just timer + button when scrolling down on mobile

### Content Area
- Max width: 1200px (centered with auto margins on large screens)
- Padding: 24px (desktop), 16px (tablet), 12px (mobile)
- Scroll: Content area scrolls, sidebar and timer bar remain fixed

## Accessibility

- **Contrast ratios**: All text meets WCAG 2.1 AA (4.5:1 for body text, 3:1 for large text)
  - Text Primary (#ECE9F5) on Background (#0D0A1A): **15.2:1** ✓
  - Text Secondary (#8C83A8) on Background (#0D0A1A): **5.8:1** ✓
  - Text Primary (#1A1433) on Light Background (#FFFBF5): **16.1:1** ✓
  - Accent (#CF8A2E) on Background (#0D0A1A): **5.6:1** ✓
- **Focus indicators**: 3px ring using Primary glow (visible, not just color-based)
- **Status**: Never rely on color alone — pair status dots with labels ("Clocked In", "On Break", "Offline")
- **Touch targets**: 44x44px minimum
- **Screen reader**: Semantic HTML (nav, main, section, table). ARIA labels for icon-only buttons.
- **Keyboard**: Full keyboard navigation. Tab order follows visual layout. Escape closes modals.
- **Reduced motion**: `prefers-reduced-motion: reduce` disables all transitions except opacity.

---
Status: READY_FOR_REVIEW
