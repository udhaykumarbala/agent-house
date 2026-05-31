# Architecture Specification: UnitShift — Premium Minimal Unit Converter

## Template: `static-enhanced` (Vite + Tailwind CSS + TypeScript)

---

## 1. Technical Overview

A purely client-side, single-page unit converter for length, weight, and temperature. No backend, no database, no API calls. All conversion logic is mathematical and runs in the browser. The emphasis is on premium visual identity, micro-interactions, and modern typography.

**Key Constraints**:
- Bundle target: < 20KB gzipped
- Zero runtime npm dependencies (dev tooling only)
- Instant load — must feel native-app fast
- Mobile-first, responsive to desktop

---

## 2. Tech Stack

| Layer | Technology | Justification |
|-------|-----------|---------------|
| Build Tool | **Vite** | Fast HMR, optimized production builds, tree-shaking |
| Styling | **Tailwind CSS** | Utility-first for rapid iteration, built-in responsive/dark mode, purge for tiny CSS |
| Language | **TypeScript** | Type safety for conversion formulas, prevents unit/direction bugs |
| Fonts | **Space Grotesk** (numbers) + **Inter** (UI) | Google Fonts, free; geometric precision for numbers, clean readability for labels |
| Animations | **CSS transitions + @keyframes** | GPU-accelerated, no JS animation library needed |
| State | **Plain object + event listeners** | App state is trivial (category, fromUnit, toUnit, value) — no library warranted |
| Hosting | **Static** (Vercel/Netlify/GitHub Pages) | No server needed, HTTPS by default |

---

## 3. Project Structure

```
converter/
├── index.html              # Entry HTML with CSP meta tag
├── package.json
├── tsconfig.json
├── vite.config.ts
├── tailwind.config.ts
├── postcss.config.js
├── public/
│   └── favicon.svg         # App icon
├── src/
│   ├── main.ts             # Entry point — initializes app, wires events
│   ├── style.css           # Tailwind directives + custom CSS animations/keyframes
│   ├── converter/
│   │   ├── engine.ts       # Pure conversion functions (all math, no DOM)
│   │   ├── units.ts        # Unit definitions, categories, metadata, default pairs
│   │   └── types.ts        # TypeScript interfaces (UnitCategory, Unit, ConversionState)
│   ├── ui/
│   │   ├── app.ts          # Main app controller — state management, orchestrates UI updates
│   │   ├── input.ts        # Input field handling — numeric validation, formatting, events
│   │   ├── selector.ts     # Unit pill/chip selector component
│   │   ├── tabs.ts         # Category tab switching with animated pill indicator
│   │   ├── swap.ts         # Swap button logic and animation trigger
│   │   └── animations.ts   # Micro-interaction helpers (number morph, focus glow, transitions)
│   └── utils/
│       └── format.ts       # Number formatting — significant digits, trailing zero trimming
└── .plans/                 # Planning documents (not deployed)
```

---

## 4. Data Model

### TypeScript Interfaces

```typescript
type CategoryId = 'length' | 'weight' | 'temperature';

interface Unit {
  id: string;            // 'km', 'lb', 'celsius'
  label: string;         // 'Kilometer', 'Pound', 'Celsius'
  symbol: string;        // 'km', 'lb', '°C'
  toBase: number | null; // Conversion factor to base unit (null for temperature)
}

interface UnitCategory {
  id: CategoryId;
  label: string;         // 'Length', 'Weight', 'Temperature'
  icon: string;          // Inline SVG path or Unicode symbol
  units: Unit[];
  defaultFrom: string;   // Default "from" unit id
  defaultTo: string;     // Default "to" unit id
}

interface ConversionState {
  category: CategoryId;
  fromUnit: string;      // Unit id
  toUnit: string;        // Unit id
  fromValue: string;     // String to preserve user input formatting
  toValue: string;       // Computed result string
  inputDirection: 'from' | 'to'; // Which field user is typing in
}
```

### Unit Definitions

**Length** (base unit: meter)

| Unit | Symbol | Factor to meters |
|------|--------|-----------------|
| Millimeter | mm | 0.001 |
| Centimeter | cm | 0.01 |
| Meter | m | 1 |
| Kilometer | km | 1000 |
| Inch | in | 0.0254 |
| Foot | ft | 0.3048 |
| Yard | yd | 0.9144 |
| Mile | mi | 1609.344 |

Default pair: **km ↔ mi**

**Weight** (base unit: gram)

| Unit | Symbol | Factor to grams |
|------|--------|----------------|
| Milligram | mg | 0.001 |
| Gram | g | 1 |
| Kilogram | kg | 1000 |
| Ounce | oz | 28.3495 |
| Pound | lb | 453.592 |
| Metric Ton | t | 1,000,000 |

Default pair: **kg ↔ lb**

**Temperature** (formula-based, not ratio-based)

| Unit | Symbol | To Celsius | From Celsius |
|------|--------|-----------|-------------|
| Celsius | °C | identity | identity |
| Fahrenheit | °F | (F − 32) × 5/9 | C × 9/5 + 32 |
| Kelvin | K | K − 273.15 | C + 273.15 |

Default pair: **°C ↔ °F**

### Conversion Strategy

- **Length & Weight**: Normalize to base unit, then convert to target: `result = value × fromFactor / toFactor`
- **Temperature**: Direct formula-based conversion using a lookup of 6 paths (C→F, F→C, C→K, K→C, F→K, K→F)
- **Precision**: Up to 6 significant digits via `toPrecision(6)`, trim trailing zeros
- **Bidirectional**: User can type in either the "from" or "to" field; the other updates live

---

## 5. UI Architecture

### Layout (Single Screen)

```
┌──────────────────────────────────────┐
│           UnitShift                  │  ← App title, minimal branding
├──────────────────────────────────────┤
│   [ Length ]  [ Weight ]  [ Temp ]   │  ← Category tabs with sliding pill
├──────────────────────────────────────┤
│                                      │
│   ┌──────────────────────────────┐   │
│   │  FROM                        │   │
│   │  ┌──────────────────────┐    │   │  ← Pill selectors for units
│   │  │ mm cm [m] km in ft yd mi│   │   │
│   │  └──────────────────────┘    │   │
│   │                              │   │
│   │       42.5                   │   │  ← Large editable number (Space Grotesk)
│   └──────────────────────────────┘   │
│                                      │
│              [ ⇅ ]                   │  ← Swap button (rotates 180° on click)
│                                      │
│   ┌──────────────────────────────┐   │
│   │  TO                          │   │
│   │  ┌──────────────────────┐    │   │
│   │  │ mm cm m km in ft yd [mi]│   │   │
│   │  └──────────────────────┘    │   │
│   │                              │   │
│   │       26.4097                │   │  ← Computed result (live update)
│   └──────────────────────────────┘   │
│                                      │
│   1 m = 3.2808 ft                    │  ← Quick reference line
│                                      │
└──────────────────────────────────────┘
```

### Unit Selectors

Horizontal scrollable pill/chip selectors — NOT dropdowns. Each unit is a tappable chip. Active chip gets the accent color fill. This pattern is:
- Thumb-friendly on mobile (large tap targets)
- Visually scannable (see all options at once)
- Modern feeling (vs. dated dropdown menus)

### Responsive Behavior

| Breakpoint | Behavior |
|-----------|----------|
| Mobile (< 640px) | Full-width, stacked layout, larger tap targets, pill selectors scroll horizontally |
| Tablet (640–1024px) | Centered card, ~480px max width |
| Desktop (> 1024px) | Centered card, 480px max width, hover states on interactive elements, keyboard input |

---

## 6. Visual Identity

### Color Palette: Warm Obsidian + Copper

Based on UI research recommendation #1 — "digital Braun" aesthetic.

**Dark Mode (Primary)**

| Token | Value | Usage |
|-------|-------|-------|
| `--bg-primary` | `#1A1714` | Page background (warm near-black) |
| `--bg-surface` | `#242019` | Card/surface background |
| `--bg-surface-hover` | `#2E2A22` | Hovered surface |
| `--border-subtle` | `rgba(255,255,255,0.06)` | Card borders |
| `--accent` | `#C17F59` | Copper accent — active states, swap button, active tab |
| `--accent-hover` | `#D4936D` | Accent hover state |
| `--accent-muted` | `rgba(193,127,89,0.15)` | Accent background tint |
| `--text-primary` | `#F5F0EB` | Primary text (warm white) |
| `--text-secondary` | `#9C9488` | Secondary text (labels, units) |
| `--text-tertiary` | `#6B6560` | Tertiary text (hints, inactive) |

**Light Mode**

| Token | Value | Usage |
|-------|-------|-------|
| `--bg-primary` | `#F7F4F0` | Page background (warm off-white) |
| `--bg-surface` | `#FFFFFF` | Card background |
| `--border-subtle` | `rgba(0,0,0,0.08)` | Card borders |
| `--accent` | `#A86B4A` | Copper accent (slightly deeper for contrast) |
| `--text-primary` | `#1A1714` | Primary text |
| `--text-secondary` | `#6B6560` | Secondary text |

### Typography Scale

| Element | Font | Size | Weight | Properties |
|---------|------|------|--------|------------|
| Primary value (editable) | Space Grotesk | 3rem (48px) | 500 | Tabular nums, letter-spacing: -0.02em |
| Result value | Space Grotesk | 2.5rem (40px) | 400 | Tabular nums |
| Category tabs | Inter | 0.9375rem (15px) | 500 | Uppercase, letter-spacing: 0.05em |
| Unit pills | Inter | 0.8125rem (13px) | 500 | — |
| Labels ("FROM", "TO") | Inter | 0.75rem (12px) | 600 | Uppercase, letter-spacing: 0.08em, text-secondary color |
| Quick reference | Inter | 0.8125rem (13px) | 400 | text-tertiary color |

### Spacing System

8px base grid. Key values: 4, 8, 12, 16, 24, 32, 48, 64.

### Border Radius

- Cards/containers: 16px
- Input areas: 12px
- Pills/chips: 24px (fully rounded)
- Swap button: 50% (circle)

### Surface Treatment

- Subtle noise texture on `--bg-primary` at 2–3% opacity for tactile quality
- Cards use `--bg-surface` with 1px `--border-subtle` border
- No heavy drop shadows — depth via color layering

---

## 7. Micro-Interactions Specification

### 7.1 Category Tab Switch
- Sliding pill background animates from current tab to new tab
- Duration: 300ms, ease `cubic-bezier(0.4, 0, 0.2, 1)`
- Content area cross-fades (opacity 1→0→1) with slight vertical shift (4px)

### 7.2 Swap Button
- Button icon rotates 180° on click
- Duration: 400ms, ease `cubic-bezier(0.34, 1.56, 0.64, 1)` (spring overshoot)
- From/To values cross-fade to their new positions
- Unit selections swap simultaneously

### 7.3 Live Conversion (Value Update)
- Result digits update with a subtle fade+scale micro-animation
- Each digit transitions independently (0.15s stagger) for a rolling effect
- Duration: 200ms per digit, ease-out

### 7.4 Input Focus
- Active card gets increased border opacity: `--border-subtle` → `rgba(193,127,89,0.3)`
- Soft copper glow: `box-shadow: 0 0 0 3px rgba(193,127,89,0.1)`
- Duration: 200ms, ease-in-out

### 7.5 Unit Pill Selection
- Selected pill fills with `--accent` background, text becomes white
- Transition: 200ms background-color + color
- Subtle scale press effect: 0.95→1.0 on click (150ms spring)

### 7.6 Category Content Change
- Old content fades out (opacity 1→0, translateY 0→-8px) — 150ms
- New content fades in (opacity 0→1, translateY 8px→0) — 150ms after delay

### 7.7 Reduced Motion
- All animations respect `prefers-reduced-motion: reduce`
- When reduced motion is preferred: instant state changes, no transitions
- Implementation: wrap all animation classes in `@media (prefers-reduced-motion: no-preference)`

---

## 8. Accessibility

| Requirement | Implementation |
|------------|----------------|
| Keyboard navigation | Tab between all interactive elements (tabs, pills, inputs, swap button) |
| Focus indicators | Visible copper-tinted focus ring on all focusable elements |
| ARIA labels | All inputs, buttons, and selectors have descriptive `aria-label` |
| Live regions | Result value wrapped in `aria-live="polite"` so screen readers announce updates |
| Color contrast | All text meets WCAG AA (4.5:1 ratio minimum) |
| Reduced motion | All animations disabled via `prefers-reduced-motion` media query |
| Color scheme | Respects `prefers-color-scheme` for default theme |
| Semantic HTML | Proper heading hierarchy, `<nav>` for tabs, `<input>` for values |

---

## 9. Security (per Security Research)

| Measure | Implementation | Priority |
|---------|---------------|----------|
| CSP meta tag | Add to `<head>` of `index.html` — `default-src 'self'; script-src 'self'; ...` | High |
| Safe DOM rendering | Use `textContent` exclusively, never `innerHTML` for user values | High |
| Input validation | Validate all input as numeric via `Number()` + `Number.isFinite()` | High |
| Dependency pinning | Exact versions in `package.json`, commit lockfile | High |
| Security headers | Configure on hosting platform (X-Frame-Options, X-Content-Type-Options, etc.) | Medium |
| No source maps in prod | `build.sourcemap: false` in `vite.config.ts` | Low |
| Self-host fonts | Optionally self-host Space Grotesk + Inter to eliminate CDN dependency | Low |

---

## 10. State Management

Simple reactive state with a plain object and manual DOM updates. No library needed.

```
State flow:
  User types in "from" field
    → sanitizeNumericInput(value)
    → convert(fromUnit, toUnit, numericValue)
    → format result with significant digits
    → update "to" field textContent
    → update quick reference line

  User clicks swap button
    → swap fromUnit ↔ toUnit in state
    → trigger swap animation
    → recalculate conversion
    → update DOM

  User switches category tab
    → update state.category
    → load default units for new category
    → trigger tab slide animation + content crossfade
    → reset values or preserve if same units exist
    → update DOM
```

---

## 11. Performance Budget

| Metric | Target |
|--------|--------|
| HTML + CSS + JS (gzipped) | < 20KB |
| First Contentful Paint | < 500ms |
| Time to Interactive | < 800ms |
| Conversion calculation | < 1ms (O(1) math) |
| Animation frame rate | 60fps (CSS-only, GPU accelerated) |
| Lighthouse Performance | > 95 |

---

## 12. Trade-offs

| Optimizing For | Sacrificing |
|---------------|-------------|
| **Tiny bundle size** | No framework conveniences (React, Vue) — manual DOM updates |
| **Premium animations** | Slightly more complex CSS — but no JS animation library |
| **Focused scope** (3 categories) | No currency, volume, data size, or other converters |
| **Dark mode primary** | Extra work to support light mode toggle |
| **Pill selectors over dropdowns** | More horizontal space needed — requires scroll on mobile for 8 units |
| **No runtime dependencies** | Must implement all helpers from scratch (minimal effort for this scope) |

---

## 13. File Responsibilities

| File | Responsibility | Key Exports |
|------|---------------|-------------|
| `main.ts` | Bootstrap app, load fonts, initialize state, wire event listeners | — |
| `converter/types.ts` | All TypeScript interfaces and type definitions | `CategoryId`, `Unit`, `UnitCategory`, `ConversionState` |
| `converter/units.ts` | Unit category definitions, metadata, default pairs | `categories`, `getCategory()`, `getUnit()` |
| `converter/engine.ts` | Pure conversion functions — no DOM, no side effects | `convert(from, to, value)`, `getQuickReference(from, to)` |
| `ui/app.ts` | Central app controller — holds state, orchestrates UI updates | `initApp()`, `getState()`, `setState()` |
| `ui/input.ts` | Input field setup, numeric validation, keypress handling | `initInputs()`, `sanitizeInput()` |
| `ui/selector.ts` | Renders pill/chip unit selectors, handles selection events | `renderSelector()`, `setActiveUnit()` |
| `ui/tabs.ts` | Category tab rendering, sliding pill animation | `initTabs()`, `switchCategory()` |
| `ui/swap.ts` | Swap button click handler, rotation animation trigger | `initSwap()` |
| `ui/animations.ts` | Reusable animation helpers — number morph, crossfade, spring | `animateNumberChange()`, `crossfade()` |
| `utils/format.ts` | Number formatting — precision, trailing zeros, display | `formatResult()`, `formatInput()` |
| `style.css` | Tailwind directives, custom properties, @keyframes, noise texture | — |

---

Status: READY_FOR_DEVELOPMENT
