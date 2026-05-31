# UI Specification: UnitShift — Premium Minimal Unit Converter

Based on research in: .plans/research/ui-research.md
Styling wireframes from: .plans/research/ux-research.md
Product context from: .plans/research/product-research.md

## Design Philosophy

"Digital Braun" — precision meets warmth. Every element earns its place. The app should feel like a well-crafted instrument: a Leica camera, a Braun ET66 calculator. The visual identity is built on restraint, warm materiality, and typographic confidence.

---

## Design System

### Color Palette (Warm Obsidian + Copper — UNIQUE)

**Rationale**: Every competitor uses blue, teal, or purple. We differentiate with a warm, material palette — deep charcoal with brown undertones paired with burnished copper as the sole accent. This feels premium and tactile, like oxidized metal on dark wood. The warmth in the black avoids the cold/tech cliche.

#### Dark Mode (Default)

| Name | Hex | RGB | Usage |
|------|-----|-----|-------|
| Background | `#131110` | rgb(19, 17, 16) | App canvas — near-black with warm brown undertone |
| Surface | `#1E1B19` | rgb(30, 27, 25) | Cards, elevated containers |
| Surface Elevated | `#2A2623` | rgb(42, 38, 35) | Hover states on cards, active surfaces |
| Border | `#3D3632` | rgb(61, 54, 50) | Subtle dividers, input borders |
| Border Focus | `#C17F59` | rgb(193, 127, 89) | Focused input borders (copper accent) |
| Primary (Copper) | `#C17F59` | rgb(193, 127, 89) | Main accent — active tabs, swap button, CTAs |
| Primary Hover | `#D4956E` | rgb(212, 149, 110) | Hover state — lighter copper |
| Primary Muted | `rgba(193, 127, 89, 0.15)` | — | Subtle copper tint for backgrounds, active tab pill |
| Text Primary | `#F2ECE6` | rgb(242, 236, 230) | Main text — warm off-white, not pure white |
| Text Secondary | `#9B918A` | rgb(155, 145, 138) | Labels, unit names, muted text |
| Text Tertiary | `#6B5F58` | rgb(107, 95, 88) | Placeholder text, disabled states |
| Success | `#7BAE7F` | rgb(123, 174, 127) | Muted sage green — copy confirmation |
| Error | `#C4654A` | rgb(196, 101, 74) | Warm terracotta red — validation errors |
| Warning | `#C7882A` | rgb(199, 136, 42) | Deep saffron — edge case warnings |

#### Light Mode (Alternative)

| Name | Hex | RGB | Usage |
|------|-----|-----|-------|
| Background | `#F5F0EB` | rgb(245, 240, 235) | Warm stone — not pure white |
| Surface | `#FFFFFF` | rgb(255, 255, 255) | Cards on stone background |
| Surface Elevated | `#FAF7F4` | rgb(250, 247, 244) | Hover states |
| Border | `#DDD5CC` | rgb(221, 213, 204) | Warm gray borders |
| Border Focus | `#A5683C` | rgb(165, 104, 60) | Deeper copper for light mode contrast |
| Primary (Copper) | `#A5683C` | rgb(165, 104, 60) | Darker copper for AA contrast on white |
| Primary Hover | `#8E5530` | rgb(142, 85, 48) | Darker on hover |
| Primary Muted | `rgba(165, 104, 60, 0.1)` | — | Copper tint backgrounds |
| Text Primary | `#2C2419` | rgb(44, 36, 25) | Dark walnut text |
| Text Secondary | `#7A6E63` | rgb(122, 110, 99) | Muted warm gray |
| Text Tertiary | `#A89E94` | rgb(168, 158, 148) | Placeholder text |
| Success | `#4A7A4E` | rgb(74, 122, 78) | Deeper sage for contrast |
| Error | `#A84832` | rgb(168, 72, 50) | Deeper terracotta |
| Warning | `#9A6A1F` | rgb(154, 106, 31) | Deeper saffron |

**Contrast ratios** (verified for WCAG AA):
- Text Primary on Background (dark): ~14.5:1 ✓
- Text Primary on Surface (dark): ~12.8:1 ✓
- Text Secondary on Background (dark): ~5.2:1 ✓
- Primary Copper on Background (dark): ~5.8:1 ✓
- Text Primary on Background (light): ~13.1:1 ✓
- Primary Copper on Surface (light): ~4.6:1 ✓

### Category Accent Colors

Each conversion category gets a subtle color identity applied to active elements (tab indicator, input focus ring). These are used at low opacity to tint the copper accent — not replace it.

| Category | Accent Tint | Hex | Usage |
|----------|------------|-----|-------|
| Length | Warm Amber | `#D4A053` | Tab active state, subtle input glow |
| Weight | Dusty Rose | `#C47A7A` | Tab active state, subtle input glow |
| Temperature | Cool Slate | `#7A9EB2` | Tab active state, subtle input glow |

These tints are applied only to:
- The active tab pill background at 15% opacity
- The input focus ring at 20% opacity
- The swap button on hover at 10% opacity

The primary copper remains the dominant brand color throughout.

---

### Typography

**Display Font**: **Space Grotesk** (Google Fonts)
- Used for: All numeric values (input + result)
- Why: Geometric with personality, distinctive number forms, highly legible digits. Gives the app character vs. generic Inter/Roboto numbers. The slightly quirky letterforms become part of the visual identity.
- Feature: `font-variant-numeric: tabular-nums` to prevent layout shift as digits change

**Body Font**: **Inter** (Google Fonts)
- Used for: Labels, unit names, category tabs, UI text
- Why: Proven readability at all sizes, excellent language support, clean pairing with Space Grotesk's expressiveness. The contrast between expressive numbers and neutral labels creates natural hierarchy.
- Feature: Variable font for precise weight control

| Element | Font | Size | Weight | Line Height | Letter Spacing | Usage |
|---------|------|------|--------|-------------|----------------|-------|
| Result Value | Space Grotesk | 56px / 3.5rem | 500 | 1.1 | -0.02em | Primary converted result — the hero |
| Input Value | Space Grotesk | 40px / 2.5rem | 400 | 1.15 | -0.01em | User-entered value |
| Unit Symbol | Inter | 18px / 1.125rem | 500 | 1.4 | 0.02em | "km", "°F", "lb" next to values |
| Category Tab | Inter | 15px / 0.9375rem | 500 | 1.4 | 0.01em | "Length", "Weight", "Temperature" |
| Unit Selector Label | Inter | 14px / 0.875rem | 400 | 1.5 | 0 | Unit name in pill selectors |
| Caption / Formula | Inter | 13px / 0.8125rem | 400 | 1.5 | 0.01em | "1 km = 0.6214 mi" reference line |
| Small / Meta | Inter | 12px / 0.75rem | 400 | 1.4 | 0.02em | Tooltips, accessibility labels |

**Number formatting**: Tabular lining figures (`font-variant-numeric: tabular-nums lining-nums`) on all numeric displays so digits hold fixed width and don't cause layout jank during live conversion.

---

### Spacing Scale

Base unit: 4px

| Token | Value | Usage |
|-------|-------|-------|
| `space-1` | 4px | Tight internal gaps (icon-to-text) |
| `space-2` | 8px | Small gaps (between pill chips, input padding-y) |
| `space-3` | 12px | Medium-small (card internal padding-x) |
| `space-4` | 16px | Standard padding (card content padding) |
| `space-5` | 20px | Section internal spacing |
| `space-6` | 24px | Between major sections (from-card to swap to to-card) |
| `space-8` | 32px | Large gaps (page top/bottom margin) |
| `space-10` | 40px | Card outer margins on mobile |
| `space-12` | 48px | Page-level horizontal margins on desktop |

---

### Border Radius

| Token | Value | Usage |
|-------|-------|-------|
| `radius-sm` | 6px | Small elements — pill chips, badges |
| `radius-md` | 10px | Inputs, buttons, small interactive elements |
| `radius-lg` | 16px | Converter cards, main containers |
| `radius-xl` | 20px | Outer app container (on desktop) |
| `radius-full` | 9999px | Category tab pill, swap button, avatars |

The rounding is soft but not bubbly — no fully-round cards. The pill shapes are reserved for small, clearly "pill-shaped" UI like tabs and chips.

---

### Shadows

Shadows are warm-tinted (using the brown-black of our palette, not pure black) to maintain the warm material quality.

| Token | Value | Usage |
|-------|-------|-------|
| `shadow-sm` | `0 1px 3px rgba(19, 17, 16, 0.12)` | Subtle lift — pill chips, buttons |
| `shadow-md` | `0 4px 12px rgba(19, 17, 16, 0.16)` | Cards at rest |
| `shadow-lg` | `0 8px 24px rgba(19, 17, 16, 0.24)` | Cards on hover/focus, elevated surfaces |
| `shadow-glow` | `0 0 20px rgba(193, 127, 89, 0.12)` | Copper glow on focused inputs |

In dark mode, shadows are less visible — rely more on surface color differentiation for elevation. In light mode, shadows are the primary elevation cue.

---

### Surface Treatment

- **Noise texture**: Subtle 2-3% opacity noise overlay on the Background layer to add tactile grain. This prevents the flat digital look and makes the dark mode feel more like a material surface.
- **Gradient on surfaces**: Cards get a barely-perceptible vertical gradient (top 1-2% lighter) to suggest real-world light direction. Not visible as a gradient — just adds life.
- **Glassmorphism**: Not used on primary elements (too trendy/dated). Reserved only for a subtle backdrop-blur on the category tabs bar if it overlaps content during scroll.

---

## Components

### Category Tabs (Segmented Control)

The primary navigation — always visible at the top of the converter.

```
┌──────────────────────────────────────────────┐
│  [ Length ]   [ Weight ]   [ Temperature ]   │
│     ●                                         │
│  ╭────────╮                                   │  ← Animated pill slides behind active tab
│  │ active │                                   │
│  ╰────────╯                                   │
└──────────────────────────────────────────────┘
```

**Styling**:
- Container: `bg-surface`, `rounded-full`, `p-1` (pill-shaped container)
- Inactive tab: `text-secondary`, `px-4 py-2`, `font-medium`
- Active tab: `text-primary`, animated pill background in `primary-muted` behind it
- Active pill: `bg-primary-muted`, `rounded-full`, slides with spring animation (200ms, ease-out with slight overshoot)
- Hover (inactive): `text-primary` transition 150ms

**Behavior**: Clicking a tab slides the pill background smoothly to the new position. The pill width adapts to the text width of the active tab.

---

### Conversion Card (From / To)

The core UI element — two stacked cards separated by a swap button.

```
┌──────────────────────────────────────┐
│  FROM                                │  ← Label in text-tertiary
│  ┌──────────────────────────┐        │
│  │ Kilometer  ·  Mile  ·  ...│       │  ← Pill selector row
│  └──────────────────────────┘        │
│                                      │
│  42.5                          km    │  ← Input value (large) + unit symbol
│                                      │
└──────────────────────────────────────┘

              [ ⇅ ]                      ← Swap button (floating between cards)

┌──────────────────────────────────────┐
│  TO                                  │
│  ┌──────────────────────────┐        │
│  │ Kilometer  ·  Mile  ·  ...│       │
│  └──────────────────────────┘        │
│                                      │
│  26.4097                       mi    │  ← Result value (hero size) + unit
│  1 km = 0.6214 mi                    │  ← Formula reference (caption)
│                                      │
└──────────────────────────────────────┘
```

**Card Styling**:
```css
/* Card Container */
background: var(--surface);           /* #1E1B19 dark / #FFFFFF light */
border: 1px solid var(--border);      /* #3D3632 dark / #DDD5CC light */
border-radius: var(--radius-lg);      /* 16px */
padding: var(--space-5) var(--space-4); /* 20px 16px */
box-shadow: var(--shadow-md);
transition: box-shadow 200ms ease, border-color 200ms ease;

/* Card Focus State (when input inside is focused) */
border-color: var(--border-focus);    /* Copper #C17F59 */
box-shadow: var(--shadow-lg), var(--shadow-glow);
```

**Input Value Styling**:
```css
/* Input number field */
font-family: 'Space Grotesk', monospace;
font-size: 2.5rem;                     /* 40px for input */
font-weight: 400;
color: var(--text-primary);
background: transparent;
border: none;
outline: none;
width: 100%;
font-variant-numeric: tabular-nums lining-nums;
caret-color: var(--primary);           /* Copper cursor */
```

**Result Value Styling**:
```css
/* Output number display */
font-family: 'Space Grotesk', monospace;
font-size: 3.5rem;                     /* 56px — hero size */
font-weight: 500;
color: var(--text-primary);
font-variant-numeric: tabular-nums lining-nums;
```

**Unit Symbol** (next to value): `text-secondary`, `font-size: 1.125rem`, `font-weight: 500`, positioned at the baseline-end of the number.

**Formula Reference**: `text-tertiary`, `font-size: 0.8125rem`, positioned below the result value with `space-2` gap.

---

### Unit Pill Selector

Horizontal row of selectable pills for choosing units within a category.

```
┌──────┐  ┌──────┐  ┌──────┐  ┌──────┐  ┌──────┐
│  mm  │  │  cm  │  │  m   │  │  km  │  │  mi  │  ...
└──────┘  └──────┘  └──────┘  └──────┘  └──────┘
                      ↑ active (copper border + copper bg tint)
```

**Styling**:
```css
/* Pill Default */
background: var(--surface-elevated);   /* #2A2623 */
border: 1px solid var(--border);       /* #3D3632 */
border-radius: var(--radius-sm);       /* 6px */
padding: var(--space-1) var(--space-3); /* 4px 12px */
font-family: 'Inter', sans-serif;
font-size: 0.875rem;
font-weight: 400;
color: var(--text-secondary);
cursor: pointer;
transition: all 150ms ease;

/* Pill Hover */
border-color: var(--text-tertiary);
color: var(--text-primary);

/* Pill Active (selected) */
background: var(--primary-muted);       /* rgba(193, 127, 89, 0.15) */
border-color: var(--primary);           /* #C17F59 */
color: var(--primary);                  /* #C17F59 */
font-weight: 500;
```

**Overflow behavior**: Horizontal scroll with `overflow-x: auto`, no visible scrollbar. On desktop, all pills fit in one row (max 8 units). On mobile, horizontal swipe for overflow.

---

### Swap Button

Floating circular button between the From and To cards.

**Styling**:
```css
/* Swap Button */
width: 44px;
height: 44px;
border-radius: var(--radius-full);
background: var(--surface-elevated);
border: 1px solid var(--border);
color: var(--text-secondary);
display: flex;
align-items: center;
justify-content: center;
cursor: pointer;
transition: all 200ms ease;
z-index: 10;  /* Floats above cards */
margin: calc(-1 * var(--space-3)) auto; /* Overlaps card edges */

/* Swap Hover */
background: var(--primary-muted);
border-color: var(--primary);
color: var(--primary);
transform: scale(1.05);
box-shadow: var(--shadow-glow);

/* Swap Active (pressing) */
transform: scale(0.95);
```

**Icon**: Vertical double arrow (↕) — 1.5px stroke, rounded caps, 20px size. SVG inline.

**Animation on click**: Rotates 180° with spring easing (`cubic-bezier(0.34, 1.56, 0.64, 1)`, 300ms). The values in both cards cross-fade simultaneously.

---

### Button (General)

```css
/* Primary Button (if ever needed — e.g., copy result) */
background: var(--primary);
color: #FFFFFF;
border: none;
border-radius: var(--radius-md);       /* 10px */
padding: var(--space-2) var(--space-4); /* 8px 16px */
font-family: 'Inter', sans-serif;
font-size: 0.875rem;
font-weight: 500;
cursor: pointer;
transition: background 150ms ease, transform 100ms ease;

/* Hover */
background: var(--primary-hover);       /* #D4956E */

/* Active */
transform: scale(0.97);

/* Disabled */
opacity: 0.4;
cursor: not-allowed;
```

```css
/* Ghost Button (copy, settings) */
background: transparent;
color: var(--text-secondary);
border: none;
border-radius: var(--radius-md);
padding: var(--space-2);
cursor: pointer;
transition: color 150ms ease, background 150ms ease;

/* Hover */
color: var(--primary);
background: var(--primary-muted);
```

---

### Copy Result Interaction

A small ghost button (clipboard icon) appears on hover/focus of the result card.

**States**:
1. **Default**: Hidden or very muted (text-tertiary), clipboard icon
2. **Hover**: Visible, copper color
3. **Clicked**: Icon morphs to checkmark, brief copper flash on the result value, text changes to "Copied" for 1.5s, then reverts
4. **Animation**: checkmark draws in with a stroke animation (200ms)

---

## Micro-Interactions & Animation Spec

All animations respect `prefers-reduced-motion: reduce` — when enabled, transitions become instant (0ms duration) with no transforms.

| Interaction | Animation | Duration | Easing | Details |
|-------------|-----------|----------|--------|---------|
| Tab switch | Pill background slides | 200ms | `cubic-bezier(0.34, 1.56, 0.64, 1)` | Spring overshoot — pill slides to active tab position |
| Category content change | Cross-fade + slight Y shift | 250ms | `ease-out` | Old content fades out + shifts down 4px, new fades in + shifts up 4px |
| Value typing | Digit scale-in | 100ms | `ease-out` | New digit enters at 0.9 scale → 1.0 |
| Live conversion result | Number roll/morph | 150ms | `ease-out` | Result digits fade-transition to new values, 50ms stagger per digit |
| Swap button click | 180° rotation | 300ms | `cubic-bezier(0.34, 1.56, 0.64, 1)` | Button rotates, values cross-fade |
| Value swap | Cross-fade values | 250ms | `ease-in-out` | From/To values simultaneously cross-fade to swapped positions |
| Card focus | Shadow elevation + border | 200ms | `ease` | Border color transitions to copper, shadow grows |
| Pill select | Background + border color | 150ms | `ease` | Smooth transition between default and active pill states |
| Copy result | Icon morph + flash | 200ms | `ease-out` | Clipboard → checkmark SVG morph, brief copper glow on value |
| Hover feedback (buttons) | Scale + color | 150ms | `ease` | Scale 1.05 on hover, 0.97 on press |

**No animation library** — all done with CSS transitions, CSS keyframes, and minimal JS for digit animation orchestration.

---

## Layout

### Overall Structure

```
┌─────────────────────────────────────────────────┐
│                  UnitShift                       │  ← App title (optional, very subtle)
│                                                  │
│  ┌─────────────────────────────────────────┐    │
│  │  [ Length ]  [ Weight ]  [ Temperature ] │    │  ← Category tabs
│  └─────────────────────────────────────────┘    │
│                                                  │
│  ┌─────────────────────────────────────────┐    │
│  │  FROM                                    │    │
│  │  [ mm ] [ cm ] [ m ] [ km ] [ mi ] ...  │    │  ← Unit pills
│  │                                          │    │
│  │  42.5                              km   │    │  ← Input + symbol
│  └─────────────────────────────────────────┘    │
│                                                  │
│                   [ ⇅ ]                          │  ← Swap
│                                                  │
│  ┌─────────────────────────────────────────┐    │
│  │  TO                                      │    │
│  │  [ mm ] [ cm ] [ m ] [ km ] [ mi ] ...  │    │  ← Unit pills
│  │                                          │    │
│  │  26.4097                            mi  │    │  ← Result (hero)
│  │  1 km = 0.6214 mi                       │    │  ← Reference formula
│  └─────────────────────────────────────────┘    │
│                                                  │
└─────────────────────────────────────────────────┘
```

### App Container

- Centered, single-column layout
- `max-width: 480px` on all screen sizes (utility app, not a dashboard)
- Vertically centered on desktop (using flexbox centering)
- Full-bleed on mobile with horizontal padding

### App Title

- "UnitShift" displayed very subtly at top center
- `text-tertiary` color, `font-size: 0.75rem`, `letter-spacing: 0.1em`, `text-transform: uppercase`
- Should feel like a watermark, not a headline — the numbers are the real UI

---

## Responsive Breakpoints

| Name | Width | Layout Notes |
|------|-------|-------------|
| Mobile | < 640px | Full width, `px-4` horizontal padding, stacked cards, category tabs may truncate "Temperature" → "Temp" |
| Tablet | 640–1024px | Centered `max-width: 480px`, comfortable spacing |
| Desktop | > 1024px | Centered `max-width: 480px`, vertically centered in viewport, subtle hover states |

### Mobile-Specific Adjustments
- Result value: `48px` (down from 56px) to prevent overflow
- Input value: `36px` (down from 40px)
- Unit pills: Horizontal scroll with momentum, slight fade on edges to hint scrollability
- Touch targets: Minimum 44x44px on all interactive elements
- Category tabs: "Temperature" may shorten to "Temp" below 360px width

### Desktop-Specific Enhancements
- Hover states on all interactive elements
- Keyboard navigation: Tab between input fields, Enter to swap, arrow keys for unit selection
- Cursor changes: `pointer` on interactive, `text` on input

---

## Accessibility

- **Contrast**: All text meets WCAG AA (4.5:1 minimum). Copper on dark surface ≥ 4.5:1. Verified above.
- **Touch targets**: 44x44px minimum on all tap/click areas
- **Focus indicators**: Visible copper outline ring (`outline: 2px solid var(--primary); outline-offset: 2px`) on keyboard focus
- **ARIA labels**: All inputs have descriptive labels ("Enter value in kilometers"), swap button has `aria-label="Swap conversion direction"`, result has `aria-live="polite"` for screen reader announcements
- **Reduced motion**: `@media (prefers-reduced-motion: reduce)` disables all animations and transitions
- **Color scheme**: Respects `prefers-color-scheme` for initial theme selection
- **Semantic HTML**: Use `<input>`, `<button>`, `<nav>` — no div-buttons

---

## Theme Toggle

A small icon button in the top-right corner to switch between dark and light mode.

- **Icon**: Sun (light mode active) / Moon (dark mode active), 1.5px stroke
- **Size**: 32x32px ghost button
- **Position**: Top-right, `space-4` from edges
- **Transition**: Colors cross-fade over 300ms when toggling. No flash.
- **Persistence**: Saved to `localStorage`, falls back to `prefers-color-scheme`

---

## CSS Custom Properties (Design Tokens)

```css
:root {
  /* Colors — Dark Mode Default */
  --bg: #131110;
  --surface: #1E1B19;
  --surface-elevated: #2A2623;
  --border: #3D3632;
  --border-focus: #C17F59;
  --primary: #C17F59;
  --primary-hover: #D4956E;
  --primary-muted: rgba(193, 127, 89, 0.15);
  --text-primary: #F2ECE6;
  --text-secondary: #9B918A;
  --text-tertiary: #6B5F58;
  --success: #7BAE7F;
  --error: #C4654A;
  --warning: #C7882A;

  /* Category Tints */
  --tint-length: #D4A053;
  --tint-weight: #C47A7A;
  --tint-temperature: #7A9EB2;

  /* Typography */
  --font-display: 'Space Grotesk', system-ui, monospace;
  --font-body: 'Inter', system-ui, sans-serif;

  /* Spacing */
  --space-1: 4px;
  --space-2: 8px;
  --space-3: 12px;
  --space-4: 16px;
  --space-5: 20px;
  --space-6: 24px;
  --space-8: 32px;
  --space-10: 40px;
  --space-12: 48px;

  /* Radius */
  --radius-sm: 6px;
  --radius-md: 10px;
  --radius-lg: 16px;
  --radius-xl: 20px;
  --radius-full: 9999px;

  /* Shadows */
  --shadow-sm: 0 1px 3px rgba(19, 17, 16, 0.12);
  --shadow-md: 0 4px 12px rgba(19, 17, 16, 0.16);
  --shadow-lg: 0 8px 24px rgba(19, 17, 16, 0.24);
  --shadow-glow: 0 0 20px rgba(193, 127, 89, 0.12);
}

/* Light Mode Override */
[data-theme="light"] {
  --bg: #F5F0EB;
  --surface: #FFFFFF;
  --surface-elevated: #FAF7F4;
  --border: #DDD5CC;
  --border-focus: #A5683C;
  --primary: #A5683C;
  --primary-hover: #8E5530;
  --primary-muted: rgba(165, 104, 60, 0.1);
  --text-primary: #2C2419;
  --text-secondary: #7A6E63;
  --text-tertiary: #A89E94;
  --success: #4A7A4E;
  --error: #A84832;
  --warning: #9A6A1F;
  --shadow-sm: 0 1px 3px rgba(44, 36, 25, 0.08);
  --shadow-md: 0 4px 12px rgba(44, 36, 25, 0.1);
  --shadow-lg: 0 8px 24px rgba(44, 36, 25, 0.14);
  --shadow-glow: 0 0 20px rgba(165, 104, 60, 0.08);
}
```

---

## Font Loading Strategy

```html
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500&family=Space+Grotesk:wght@400;500&display=swap" rel="stylesheet">
```

- `font-display: swap` to prevent invisible text during load
- Only load weights 400 + 500 (all we need)
- Preconnect for faster resolution

---

## Icon System

Minimal line icons, 1.5px stroke, rounded line caps, 20×20px viewbox. Only icons used:

| Icon | Usage | Description |
|------|-------|-------------|
| Swap (↕) | Swap button | Two vertical arrows, pointing opposite directions |
| Sun | Theme toggle (light active) | Circle with rays |
| Moon | Theme toggle (dark active) | Crescent moon |
| Copy | Copy result | Clipboard outline |
| Check | Copy confirmation | Checkmark, animates in with stroke-dashoffset |

All icons are inline SVG — no icon library. Keeps bundle minimal and allows animation control.

---

Status: READY_FOR_REVIEW
