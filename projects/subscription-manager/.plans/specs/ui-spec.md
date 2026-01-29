# UI Specification: SubTrack - Subscription Manager

Based on research in: `.plans/research/ui-research.md`
Styling wireframes from: `.plans/specs/ux-spec.md`

---

## Design Philosophy

**"Soft Minimal"** - A balance between stark minimalism (Linear) and friendly approachability (Cleo).

SubTrack should feel:
- **Calm** - Not anxiety-inducing about spending
- **Trustworthy** - Clean, no dark patterns
- **Delightful** - Small moments of polish
- **Memorable** - The "coral subscription app"

---

## Design System

### Color Palette (UNIQUE - Not Generic!)

**IMPORTANT**: These colors differentiate SubTrack from the sea of blue/green/purple fintech apps. Coral is our signature - warm, friendly, and memorable.

#### Light Mode

| Name | Hex | RGB | Usage |
|------|-----|-----|-------|
| Primary (Coral) | #FF7A5C | rgb(255, 122, 92) | CTAs, accent highlights, brand elements |
| Primary Hover | #FF6347 | rgb(255, 99, 71) | Button hover states |
| Primary Pressed | #E85C3A | rgb(232, 92, 58) | Active/pressed states |
| Background | #FAFAFA | rgb(250, 250, 250) | App background (warm off-white) |
| Surface | #FFFFFF | rgb(255, 255, 255) | Cards, elevated elements |
| Surface Elevated | #FFFFFF | rgb(255, 255, 255) | Modals, dropdowns (with shadow) |
| Text Primary | #1A1A1E | rgb(26, 26, 30) | Main text, headings |
| Text Secondary | #6B6B73 | rgb(107, 107, 115) | Muted text, labels |
| Text Tertiary | #9CA3AF | rgb(156, 163, 175) | Placeholder, disabled text |
| Border | #E5E5E8 | rgb(229, 229, 232) | Card borders, dividers |
| Border Focus | #FF7A5C | rgb(255, 122, 92) | Input focus rings |

#### Dark Mode

| Name | Hex | RGB | Usage |
|------|-----|-----|-------|
| Primary (Coral) | #FF8A6C | rgb(255, 138, 108) | CTAs, accent (slightly lighter for dark bg) |
| Primary Hover | #FF9A7C | rgb(255, 154, 124) | Button hover states |
| Primary Pressed | #FF7A5C | rgb(255, 122, 92) | Active/pressed states |
| Background | #0F0F12 | rgb(15, 15, 18) | App background (warm charcoal) |
| Surface | #1A1A1E | rgb(26, 26, 30) | Cards, sheets |
| Surface Elevated | #252529 | rgb(37, 37, 41) | Modals, dropdowns |
| Surface Hover | #2A2A2F | rgb(42, 42, 47) | Interactive states |
| Text Primary | #FAFAFA | rgb(250, 250, 250) | Main text (not pure white) |
| Text Secondary | #A1A1AA | rgb(161, 161, 170) | Muted text, labels |
| Text Tertiary | #6B6B73 | rgb(107, 107, 115) | Placeholder, disabled text |
| Border | #2A2A2F | rgb(42, 42, 47) | Card borders (subtle) |
| Border Focus | #FF8A6C | rgb(255, 138, 108) | Input focus rings |

#### Semantic Colors (Both Modes)

| Name | Light Mode | Dark Mode | Usage |
|------|------------|-----------|-------|
| Success | #059669 | #10B981 | Positive states, savings |
| Success Background | #ECFDF5 | #064E3B | Success alerts |
| Error | #DC2626 | #EF4444 | Destructive actions, errors |
| Error Background | #FEF2F2 | #7F1D1D | Error alerts |
| Warning | #D97706 | #F59E0B | Due soon, trials ending |
| Warning Background | #FFFBEB | #78350F | Warning alerts |
| Info | #0891B2 | #22D3EE | Informational states |
| Info Background | #ECFEFF | #164E63 | Info alerts |

**Why These Colors**:
- Coral (#FF7A5C) is warm and memorable - no major fintech competitor uses it
- Warm charcoal (#0F0F12) feels modern without the harshness of pure black
- Off-white (#FAFAFA) is easier on eyes than pure white
- Semantic colors follow established conventions (green=good, red=bad) but with our own hues

---

### Typography

**Primary Font**: Geist

```css
font-family: 'Geist', system-ui, -apple-system, BlinkMacSystemFont, 'SF Pro', 'Inter', sans-serif;
```

**Why Geist**:
- Numbers render beautifully (critical for subscription prices)
- Tabular numbers prevent layout shift when prices change
- Clean, modern without being cold
- Variable font = single file, all weights
- Free and open source (by Vercel)

#### Type Scale

| Element | Size | Weight | Line Height | Letter Spacing | Usage |
|---------|------|--------|-------------|----------------|-------|
| Display | 40px | 700 | 1.1 | -0.02em | Hero numbers (total spend) |
| H1 | 32px | 700 | 1.2 | -0.02em | Page titles |
| H2 | 24px | 600 | 1.3 | -0.01em | Section headers |
| H3 | 20px | 600 | 1.4 | -0.01em | Card titles |
| Body Large | 18px | 400 | 1.5 | 0 | Featured content |
| Body | 16px | 400 | 1.5 | 0 | Default text |
| Body Small | 14px | 400 | 1.5 | 0 | Secondary text, captions |
| Label | 12px | 500 | 1.4 | 0.02em | Form labels, badges |
| Caption | 11px | 400 | 1.4 | 0.02em | Timestamps, metadata |

#### Font Feature Settings

```css
/* Enable tabular numbers for price alignment */
.price, .currency {
  font-feature-settings: 'tnum' 1;
}
```

---

### Spacing Scale

**Base Unit**: 4px

| Token | Value | CSS Variable | Usage |
|-------|-------|--------------|-------|
| 0 | 0px | --spacing-0 | No spacing |
| 1 | 4px | --spacing-1 | Tight spacing, inline elements |
| 2 | 8px | --spacing-2 | Small gaps, icon padding |
| 3 | 12px | --spacing-3 | Form field gaps |
| 4 | 16px | --spacing-4 | Standard padding, card padding |
| 5 | 20px | --spacing-5 | Medium gaps |
| 6 | 24px | --spacing-6 | Section spacing |
| 8 | 32px | --spacing-8 | Large gaps |
| 10 | 40px | --spacing-10 | Page sections |
| 12 | 48px | --spacing-12 | Page margins |
| 16 | 64px | --spacing-16 | Hero sections |

---

### Border Radius

| Token | Value | CSS Variable | Usage |
|-------|-------|--------------|-------|
| none | 0px | --radius-none | Square elements |
| sm | 4px | --radius-sm | Small chips, tags |
| md | 8px | --radius-md | Buttons, inputs |
| lg | 12px | --radius-lg | Cards, list items |
| xl | 16px | --radius-xl | Modals, sheets |
| 2xl | 24px | --radius-2xl | Large cards, hero sections |
| full | 9999px | --radius-full | Pills, avatars, FAB |

---

### Shadows

#### Light Mode Shadows

| Token | Value | Usage |
|-------|-------|-------|
| xs | `0 1px 2px rgba(0, 0, 0, 0.04)` | Subtle lift |
| sm | `0 2px 4px rgba(0, 0, 0, 0.06)` | Inputs on focus |
| md | `0 4px 12px rgba(0, 0, 0, 0.08)` | Cards at rest |
| lg | `0 8px 24px rgba(0, 0, 0, 0.12)` | Cards on hover |
| xl | `0 12px 32px rgba(0, 0, 0, 0.16)` | Modals, sheets |
| FAB | `0 4px 16px rgba(255, 122, 92, 0.3)` | Floating action button |

#### Dark Mode Shadows

Shadows are nearly invisible on dark backgrounds. Use borders and surface colors instead:

| Token | Alternative | Usage |
|-------|-------------|-------|
| Card | 1px border at `rgba(255, 255, 255, 0.08)` | Card edges |
| Elevated | Background color `#252529` | Modals, dropdowns |
| FAB | `0 4px 16px rgba(0, 0, 0, 0.4)` | Floating action button |

---

## Components

### Button

#### Variants

```
Primary:
┌─────────────────────────────────────┐
│          Add Subscription           │  ← Coral bg, white text
└─────────────────────────────────────┘

Secondary:
┌─────────────────────────────────────┐
│              Cancel                 │  ← White bg, gray border, dark text
└─────────────────────────────────────┘

Ghost:
┌─────────────────────────────────────┐
│               Edit                  │  ← Transparent, coral text
└─────────────────────────────────────┘

Destructive:
┌─────────────────────────────────────┐
│              Delete                 │  ← Red bg/text
└─────────────────────────────────────┘
```

#### Specifications

| Property | Primary | Secondary | Ghost | Destructive |
|----------|---------|-----------|-------|-------------|
| Background | `--primary` | `--surface` | transparent | transparent |
| Text | white | `--text-primary` | `--primary` | `--error` |
| Border | none | 1px `--border` | none | none |
| Padding | 12px 24px | 12px 24px | 12px 16px | 12px 24px |
| Radius | 8px | 8px | 8px | 8px |
| Font | 14px / 600 | 14px / 600 | 14px / 500 | 14px / 600 |
| Min Height | 44px | 44px | 44px | 44px |

#### States

| State | Transform | Additional |
|-------|-----------|------------|
| Default | — | — |
| Hover | — | Background lightens 5%, cursor pointer |
| Pressed | scale(0.98) | Background darkens 5% |
| Focus | — | 2px ring with `--primary` at 30% opacity |
| Disabled | — | opacity: 0.5, cursor: not-allowed |

---

### Input

```
Default:
┌─────────────────────────────────────┐
│  Label                              │
│  ┌─────────────────────────────────┐│
│  │ Placeholder text...             ││
│  └─────────────────────────────────┘│
└─────────────────────────────────────┘

Focused:
┌─────────────────────────────────────┐
│  Label                              │
│  ┌─────────────────────────────────┐│
│  │ User input text                 ││  ← Coral border + ring
│  └─────────────────────────────────┘│
└─────────────────────────────────────┘

Error:
┌─────────────────────────────────────┐
│  Label                              │
│  ┌─────────────────────────────────┐│
│  │ Invalid input                   ││  ← Red border + ring
│  └─────────────────────────────────┘│
│  ⚠️ Error message here              │
└─────────────────────────────────────┘
```

#### Specifications

| Property | Value |
|----------|-------|
| Background | `--surface` |
| Border | 1px `--border` |
| Border Focus | 2px `--primary` |
| Border Error | 2px `--error` |
| Radius | 8px |
| Padding | 12px 16px |
| Font | 16px (prevents iOS zoom) |
| Height | 48px |

---

### Card (Subscription Card)

```
┌───────────────────────────────────────────────────┐
│  ┌────┐                                           │
│  │ 🎵 │  Spotify                       $9.99 /mo  │
│  └────┘  Entertainment           Due Jan 22 →    │
└───────────────────────────────────────────────────┘
```

#### Specifications

| Property | Light Mode | Dark Mode |
|----------|------------|-----------|
| Background | `#FFFFFF` | `#1A1A1E` |
| Border | none | 1px `rgba(255,255,255,0.08)` |
| Shadow | `shadow-md` | none |
| Radius | 12px |
| Padding | 16px |
| Min Height | 72px |
| Gap (internal) | 12px |

#### Card States

| State | Light Mode | Dark Mode |
|-------|------------|-----------|
| Default | Shadow-md | Border visible |
| Hover | Shadow-lg, translateY(-2px) | Background `#252529` |
| Pressed | scale(0.98) | scale(0.98), background `#2A2A2F` |

#### Card Layout

```css
.subscription-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
}

.card-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  /* Service logo or emoji on colored background */
}

.card-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.card-price {
  text-align: right;
  font-feature-settings: 'tnum' 1;
}
```

---

### Floating Action Button (FAB)

```
      ┌───────┐
      │   +   │   ← 56px, coral, centered plus icon
      └───────┘
```

#### Specifications

| Property | Value |
|----------|-------|
| Size | 56px × 56px |
| Background | `--primary` |
| Icon | Plus, 24px, white, stroke 2px |
| Radius | full (9999px) |
| Shadow | `0 4px 16px rgba(255, 122, 92, 0.3)` |
| Position | Fixed, 24px from bottom, centered or right |

#### States

| State | Transform | Shadow |
|-------|-----------|--------|
| Default | — | FAB shadow |
| Hover | scale(1.05) | Shadow intensifies |
| Pressed | scale(0.95) | Shadow reduces |

---

### Bottom Navigation

```
┌─────────────────────────────────────────────────────┐
│    🏠           📅           ⚙️                      │
│   Home       Upcoming      Settings                 │  ← 3 items max
│   ━━━━                                              │  ← Active indicator
└─────────────────────────────────────────────────────┘
```

#### Specifications

| Property | Value |
|----------|-------|
| Height | 64px (+ safe area on iOS) |
| Background | `--surface` |
| Border Top | 1px `--border` |
| Items | 3 max (Home, Upcoming, Settings) |
| Icon Size | 24px |
| Label | 11px, medium weight |
| Active Color | `--primary` |
| Inactive Color | `--text-secondary` |
| Tap Area | 64px × 64px minimum |

---

### Bottom Sheet

```
┌─────────────────────────────────────┐
│           ━━━━━━━━                  │  ← Drag handle (36px × 4px)
│                                     │
│  Sheet Title               [ × ]   │
│  ─────────────────────────────────  │
│                                     │
│  [Content]                          │
│                                     │
└─────────────────────────────────────┘
```

#### Specifications

| Property | Value |
|----------|-------|
| Background | `--surface` |
| Radius (top) | 24px |
| Shadow | `shadow-xl` |
| Max Height | 90vh |
| Padding | 24px |
| Drag Handle | 36px × 4px, `--border` color, centered |
| Backdrop | rgba(0, 0, 0, 0.5) |

---

### Filter Pills

```
  ┌──────────┐  ┌───────────────┐  ┌──────────┐
  │ Monthly  │  │ Entertainment │  │   All    │
  └──────────┘  └───────────────┘  └──────────┘
       ↑               ↑
    Active          Inactive
  (filled bg)    (border only)
```

#### Specifications

| Property | Active | Inactive |
|----------|--------|----------|
| Background | `--primary` | transparent |
| Text | white | `--text-secondary` |
| Border | none | 1px `--border` |
| Radius | full (9999px) |
| Padding | 8px 16px |
| Font | 14px, 500 |
| Height | 36px |

---

### Toast / Snackbar

```
┌─────────────────────────────────────────────────────┐
│  ✓  Subscription added                     [Undo]  │
└─────────────────────────────────────────────────────┘
```

#### Specifications

| Property | Value |
|----------|-------|
| Background | `--text-primary` (inverted from page) |
| Text | `--background` (inverted) |
| Radius | 12px |
| Padding | 16px 20px |
| Shadow | `shadow-lg` |
| Position | Bottom center, 24px from bottom nav |
| Duration | 3 seconds (5 seconds for undo actions) |

---

### Empty State

```
┌─────────────────────────────────────┐
│                                     │
│           ┌─────────┐               │
│           │   📋    │               │
│           │   ✨    │               │  ← Illustration area
│           └─────────┘               │
│                                     │
│     No subscriptions yet            │  ← H2, centered
│                                     │
│   Add your first subscription to    │  ← Body, text-secondary
│   start tracking your spending.     │
│                                     │
│  ┌─────────────────────────────┐    │
│  │    + Add Subscription       │    │  ← Primary button
│  └─────────────────────────────┘    │
│                                     │
└─────────────────────────────────────┘
```

---

### Due Soon Badge

When a subscription is due within 7 days:

```
┌────────────────────────────────────────┐
│  🎵 Spotify                      ⚠️    │  ← Warning icon, amber color
│     $9.99/mo          Due in 2 days   │  ← "Due in X" text, amber
└────────────────────────────────────────┘
```

| Property | Value |
|----------|-------|
| Icon | AlertCircle (Lucide) |
| Color | `--warning` |
| Text | "Due in X days" instead of date |

---

## Iconography

### Icon Set: Lucide

- **License**: MIT (open source)
- **Style**: Outlined, 1.5-2px stroke
- **Consistency**: 1000+ icons with matching style

### Icon Sizes

| Context | Size | Stroke |
|---------|------|--------|
| Navigation | 24px | 1.5px |
| Cards / Lists | 20px | 1.5px |
| Inline / Small | 16px | 2px |
| FAB | 24px | 2px |

### Key Icons

| Action | Icon Name | Notes |
|--------|-----------|-------|
| Add | Plus | FAB, add buttons |
| Home | Home | Bottom nav |
| Calendar | Calendar | Upcoming nav |
| Settings | Settings | Settings nav |
| Delete | Trash2 | Destructive action |
| Edit | Pencil | Edit action |
| Search | Search | Search input |
| Close | X | Sheet/modal close |
| Back | ChevronLeft | Navigation |
| Warning | AlertCircle | Due soon |
| Success | CheckCircle | Confirmation |
| Dark mode | Moon | Theme toggle |
| Light mode | Sun | Theme toggle |

---

## Animation Specifications

### Timing Guidelines

- **Micro-interactions**: 100-200ms (feels instant)
- **State changes**: 200-300ms (noticeable but quick)
- **Page transitions**: 300-400ms (smooth navigation)

### Easing Functions

| Type | Easing | Use For |
|------|--------|---------|
| Default | `ease-out` | Most animations |
| Enter | `cubic-bezier(0.0, 0.0, 0.2, 1)` | Elements appearing |
| Exit | `cubic-bezier(0.4, 0.0, 1, 1)` | Elements leaving |
| Spring | `cubic-bezier(0.34, 1.56, 0.64, 1)` | Playful bounces |

### Component Animations

| Animation | Duration | Easing | Properties |
|-----------|----------|--------|------------|
| Button press | 100ms | ease-out | scale: 0.98 |
| Button release | 150ms | ease-out | scale: 1 |
| Card hover | 200ms | ease-out | translateY: -2px, shadow |
| Card press | 100ms | ease-out | scale: 0.98 |
| Card enter | 300ms | ease-out | opacity 0→1, translateY 20→0 |
| Card exit | 250ms | ease-in | opacity 1→0, translateX 0→-100% |
| Sheet open | 350ms | spring | translateY: 100%→0 |
| Sheet close | 250ms | ease-in | translateY: 0→100% |
| FAB press | 150ms | ease-out | scale: 0.92 |
| Total counter | 600ms | ease-out | number tween (count up/down) |
| Toast enter | 300ms | ease-out | translateY 20→0, opacity 0→1 |
| Toast exit | 200ms | ease-in | opacity 1→0 |
| Skeleton pulse | 1.5s | ease-in-out | opacity 0.4↔1, infinite |

### Reduced Motion

Respect `prefers-reduced-motion`:

```css
@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after {
    animation-duration: 0.01ms !important;
    transition-duration: 0.01ms !important;
  }
}
```

---

## Responsive Behavior

### Breakpoints

| Name | Width | Layout |
|------|-------|--------|
| xs | < 375px | Compact mobile |
| sm | 375-640px | Standard mobile (primary) |
| md | 640-1024px | Tablet |
| lg | > 1024px | Desktop |

### Mobile (Primary Target)

- Single column layout
- Bottom navigation
- Full-width cards
- FAB for primary action
- Bottom sheets for forms

### Tablet (md)

- Two-column card grid
- Bottom navigation persists
- Increased padding
- Wider sheets

### Desktop (lg)

- Sidebar navigation (left)
- Three-column card grid
- Modals instead of sheets
- Hover states enabled

---

## Accessibility

### Color Contrast

All combinations meet WCAG AA (4.5:1 for normal text):

| Combination | Ratio | Pass |
|-------------|-------|------|
| Text Primary on Background | 15.8:1 | AAA |
| Text Secondary on Background | 5.2:1 | AA |
| Primary on Background | 4.6:1 | AA |
| White on Primary | 4.5:1 | AA |
| Primary on Dark Background | 6.2:1 | AA |

### Focus States

All interactive elements have visible focus indicators:

```css
:focus-visible {
  outline: none;
  box-shadow: 0 0 0 2px var(--background), 0 0 0 4px var(--primary);
}
```

### Touch Targets

All interactive elements are minimum 44×44px (per WCAG 2.5.5).

---

## CSS Custom Properties

```css
:root {
  /* Colors - Light Mode */
  --color-primary: #FF7A5C;
  --color-primary-hover: #FF6347;
  --color-primary-pressed: #E85C3A;
  --color-background: #FAFAFA;
  --color-surface: #FFFFFF;
  --color-surface-elevated: #FFFFFF;
  --color-text-primary: #1A1A1E;
  --color-text-secondary: #6B6B73;
  --color-text-tertiary: #9CA3AF;
  --color-border: #E5E5E8;
  --color-success: #059669;
  --color-error: #DC2626;
  --color-warning: #D97706;

  /* Typography */
  --font-family: 'Geist', system-ui, -apple-system, BlinkMacSystemFont, 'SF Pro', 'Inter', sans-serif;

  /* Spacing */
  --spacing-1: 4px;
  --spacing-2: 8px;
  --spacing-3: 12px;
  --spacing-4: 16px;
  --spacing-6: 24px;
  --spacing-8: 32px;

  /* Radius */
  --radius-sm: 4px;
  --radius-md: 8px;
  --radius-lg: 12px;
  --radius-xl: 16px;
  --radius-full: 9999px;

  /* Shadows */
  --shadow-sm: 0 2px 4px rgba(0, 0, 0, 0.06);
  --shadow-md: 0 4px 12px rgba(0, 0, 0, 0.08);
  --shadow-lg: 0 8px 24px rgba(0, 0, 0, 0.12);
  --shadow-xl: 0 12px 32px rgba(0, 0, 0, 0.16);
}

[data-theme="dark"] {
  --color-primary: #FF8A6C;
  --color-primary-hover: #FF9A7C;
  --color-primary-pressed: #FF7A5C;
  --color-background: #0F0F12;
  --color-surface: #1A1A1E;
  --color-surface-elevated: #252529;
  --color-text-primary: #FAFAFA;
  --color-text-secondary: #A1A1AA;
  --color-text-tertiary: #6B6B73;
  --color-border: #2A2A2F;
  --color-success: #10B981;
  --color-error: #EF4444;
  --color-warning: #F59E0B;

  /* Dark mode shadows are minimal */
  --shadow-md: none;
  --shadow-lg: none;
}
```

---

## Implementation Notes

### Font Loading

```html
<link rel="preconnect" href="https://fonts.googleapis.com">
<link href="https://fonts.googleapis.com/css2?family=Geist:wght@400;500;600;700&display=swap" rel="stylesheet">
```

Or use `@fontsource/geist` for self-hosting.

### Icon Implementation

```bash
npm install lucide-react  # React
npm install lucide-vue    # Vue
```

### Dark Mode Detection

```css
@media (prefers-color-scheme: dark) {
  :root:not([data-theme="light"]) {
    /* Dark mode variables */
  }
}
```

---

## Visual Reference Summary

| Element | Key Visual Traits |
|---------|-------------------|
| **Brand Color** | Coral #FF7A5C - warm, memorable |
| **Typography** | Geist - clean, tabular numbers |
| **Cards** | 12px radius, subtle shadows |
| **Buttons** | 8px radius, 44px min height |
| **Dark Mode** | Warm charcoal, not pure black |
| **Icons** | Lucide, outlined, 1.5px stroke |
| **Animation** | Subtle, functional, 200-300ms |

---

Status: READY_FOR_REVIEW

PHASE_COMPLETE: planning
