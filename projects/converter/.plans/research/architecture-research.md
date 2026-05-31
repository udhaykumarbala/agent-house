# Architecture Research: Premium Unit Converter

## 1. Competitor & Reference Analysis

### Existing Unit Converters (What to Learn From)

**Google Unit Converter (built-in search)**
- Strengths: Instant results, clean two-panel layout (from → to), dropdown selectors
- Weaknesses: Generic Google styling, no personality, no micro-interactions
- Takeaway: The two-field mirrored layout is the proven UX pattern — users expect it

**Apple Calculator (unit conversion mode)**
- Strengths: Minimal, elegant typography, smooth animations, haptic-like transitions
- Weaknesses: Buried feature, limited discoverability
- Takeaway: Apple proves a converter can feel premium with just typography + whitespace + subtle motion

**Convertio / UnitConverters.net**
- Strengths: Comprehensive unit coverage
- Weaknesses: Ad-heavy, cluttered, dated design, poor mobile UX
- Takeaway: These are exactly what we're differentiating from — "generic converter tools"

**Linear App / Raycast / Arc Browser (design inspiration, not converters)**
- Strengths: Monochrome palettes with accent colors, fluid micro-interactions, premium feel from restraint
- Takeaway: Premium = restraint. Fewer colors, more whitespace, deliberate motion

### Key Design Patterns in Premium Utility Apps

1. **Single-purpose focus** — One thing on screen at a time, no clutter
2. **Real-time conversion** — Results update as you type (no "Convert" button)
3. **Swap animation** — Animated unit swap (from ↔ to) is the signature interaction
4. **Category tabs** — Length | Weight | Temperature as clean tab/pill selectors
5. **Subtle depth** — Glass morphism or soft shadows, not flat or skeuomorphic
6. **Monospace or geometric sans-serif** for numbers — Gives a technical, precise feel

---

## 2. Technical Architecture

### Template: `static-enhanced` (Vite + Tailwind CSS + TypeScript)

This is a purely client-side app. No backend, no database, no API calls. All conversion logic is mathematical.

### Core Architecture Pattern

**Single-Page App with Module Pattern**

```
src/
├── main.ts              # Entry point, initializes app
├── style.css            # Tailwind directives + custom CSS animations
├── converter/
│   ├── engine.ts        # Pure conversion functions (all math)
│   ├── units.ts         # Unit definitions, categories, metadata
│   └── types.ts         # TypeScript interfaces
├── ui/
│   ├── app.ts           # Main app controller, state management
│   ├── input.ts         # Input field handling, formatting
│   ├── selector.ts      # Unit dropdown/selector component
│   ├── tabs.ts          # Category tab switching
│   └── animations.ts    # Micro-interaction helpers
└── utils/
    └── format.ts        # Number formatting, precision handling
```

### Why This Structure
- **Separation of concerns**: Conversion logic is pure functions (testable, no DOM), UI is separate
- **No framework needed**: Vanilla TS with module pattern keeps bundle tiny (<20KB)
- **Tailwind handles styling**: No CSS architecture decisions needed
- **Type safety**: TypeScript prevents conversion formula bugs

---

## 3. Conversion Logic Design

### Categories & Units

**Length**
| Unit | Key | Base (meters) |
|------|-----|---------------|
| Millimeter | mm | 0.001 |
| Centimeter | cm | 0.01 |
| Meter | m | 1 |
| Kilometer | km | 1000 |
| Inch | in | 0.0254 |
| Foot | ft | 0.3048 |
| Yard | yd | 0.9144 |
| Mile | mi | 1609.344 |

**Weight**
| Unit | Key | Base (grams) |
|------|-----|--------------|
| Milligram | mg | 0.001 |
| Gram | g | 1 |
| Kilogram | kg | 1000 |
| Ounce | oz | 28.3495 |
| Pound | lb | 453.592 |
| Ton (metric) | t | 1000000 |

**Temperature** (special — not ratio-based)
| Unit | Key | Formula from Celsius |
|------|-----|---------------------|
| Celsius | °C | identity |
| Fahrenheit | °F | (C × 9/5) + 32 |
| Kelvin | K | C + 273.15 |

### Conversion Strategy
- **Length & Weight**: Convert to base unit first, then to target unit (`value * fromFactor / toFactor`)
- **Temperature**: Direct formula-based conversion (switch statement with 6 paths: C↔F, C↔K, F↔K)
- **Precision**: Display up to 6 significant digits, trim trailing zeros

### Data Model (TypeScript)

```typescript
interface UnitCategory {
  id: 'length' | 'weight' | 'temperature';
  label: string;
  icon: string; // SVG inline or emoji
  units: Unit[];
}

interface Unit {
  id: string;          // 'km', 'lb', 'celsius'
  label: string;       // 'Kilometer', 'Pound', 'Celsius'
  symbol: string;      // 'km', 'lb', '°C'
  toBase: number | null;  // Conversion factor (null for temperature)
}

interface ConversionState {
  category: UnitCategory['id'];
  fromUnit: string;
  toUnit: string;
  fromValue: string;   // String to preserve user input formatting
  toValue: string;     // Computed result
}
```

---

## 4. UI/UX Architecture

### Layout Structure

```
┌─────────────────────────────────┐
│         Category Tabs           │  ← Length | Weight | Temperature
│    [ Length ] [ Weight ] [ °C ] │
├─────────────────────────────────┤
│                                 │
│   ┌───────────────────────┐     │
│   │  FROM                 │     │  ← Unit selector + input
│   │  [Kilometer    ▼]    │     │
│   │  42.5                 │     │  ← Large, editable number
│   └───────────────────────┘     │
│                                 │
│          [ ⇅ swap ]             │  ← Animated swap button
│                                 │
│   ┌───────────────────────┐     │
│   │  TO                   │     │  ← Unit selector + result
│   │  [Mile         ▼]    │     │
│   │  26.4097              │     │  ← Computed, updates live
│   └───────────────────────┘     │
│                                 │
│   ┌───────────────────────┐     │
│   │  Quick reference:     │     │  ← Optional: formula display
│   │  1 km = 0.6214 mi    │     │
│   └───────────────────────┘     │
│                                 │
└─────────────────────────────────┘
```

### Micro-Interactions (Key Differentiators)

1. **Tab switch**: Active tab slides with a pill background that animates between positions (not instant highlight)
2. **Swap units**: Button rotates 180°, values cross-fade to their new positions
3. **Live conversion**: Result digit updates with a subtle scale/fade micro-animation
4. **Input focus**: Card subtly elevates (shadow increase) on focus
5. **Unit dropdown**: Slides down with staggered item entrance
6. **Number formatting**: Digits animate in/out as user types (counter-style)
7. **Category change**: Content cross-fades with slight vertical shift

### Animation Implementation

- **CSS transitions** for simple state changes (focus, hover, active)
- **CSS @keyframes** for the swap rotation and digit animations
- **Tailwind `transition-*` utilities** for most interactions
- **`requestAnimationFrame`** only if needed for digit counter animation
- No animation library — keep it lightweight

---

## 5. Typography & Visual Identity Research

### Font Pairing Strategy

**Numbers**: A geometric or monospace font for values — gives precision and technical feel
- **Candidates**: Inter (geometric sans), JetBrains Mono (monospace), Space Grotesk (geometric)
- **Recommendation**: **Space Grotesk** for numbers — geometric, modern, excellent digit legibility

**UI Text**: Clean sans-serif for labels, tabs, unit names
- **Recommendation**: **Inter** — proven, excellent readability, variable font for performance

Both are Google Fonts (free, reliable CDN).

### Color Palette Direction

Premium feel = restrained palette. Two strong directions:

**Option A: Dark Mode Primary (Recommended)**
- Background: Near-black (`#0a0a0b`) or dark slate (`#0f172a`)
- Surface cards: Slightly lighter with subtle border (`#1e293b`)
- Accent: Single vibrant color — electric blue (`#3b82f6`) or emerald (`#10b981`)
- Text: White/gray hierarchy
- Why: Dark themes feel more premium for utility tools (see: Raycast, Linear, Arc)

**Option B: Light Mode with Depth**
- Background: Warm off-white (`#fafaf9`)
- Surface cards: White with soft shadows
- Accent: Deep indigo or rich violet
- Why: Clean, Apple-like, accessible

**Recommendation**: Dark mode primary with light mode toggle. Dark feels more "next-gen."

### Glassmorphism / Depth Effects

Subtle glass effect on the converter cards:
- `backdrop-filter: blur(12px)`
- Semi-transparent background (`bg-white/5` or `bg-slate-800/60`)
- Thin border (`border border-white/10`)
- This gives depth without being heavy

---

## 6. Performance Considerations

- **Bundle target**: < 20KB gzipped (Vite + Tailwind purge makes this easy)
- **No runtime dependencies**: Zero npm packages beyond dev tooling
- **Instant load**: Should feel native-app fast
- **Input debounce**: Not needed — conversions are O(1) math, can run on every keystroke
- **Font loading**: Use `font-display: swap` to prevent FOIT

---

## 7. Accessibility

- Full keyboard navigation (Tab between fields, Enter to swap)
- ARIA labels on all interactive elements
- Sufficient contrast ratios (4.5:1 minimum for text)
- Screen reader announces conversion results
- Respects `prefers-reduced-motion` for all animations
- Respects `prefers-color-scheme` for default theme

---

## 8. Key Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Floating point precision errors | Wrong conversions displayed | Use `toPrecision()` with smart rounding, test edge cases |
| Animation jank on low-end devices | Poor perceived quality | CSS-only animations (GPU accelerated), `will-change` hints |
| Font flash (FOUT/FOIT) | Perceived loading delay | `font-display: swap` + preload critical fonts |
| Dropdown accessibility | Unusable for keyboard/screen reader | Custom select built with proper ARIA roles |
| Temperature formula errors | Wrong results | Dedicated test suite for all 6 conversion paths |

---

## 9. Recommendations Summary

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Template | `static-enhanced` | No backend needed, Tailwind + TS add real value |
| Architecture | Vanilla TS modules | No framework overhead, tiny bundle, fast |
| Conversion approach | Base-unit normalization + temperature formulas | Clean, extensible, testable |
| Primary theme | Dark mode | Premium feel, differentiates from generic tools |
| Typography | Space Grotesk (numbers) + Inter (UI) | Modern, precise, great digit readability |
| Animations | CSS transitions + keyframes only | GPU-accelerated, no JS animation lib |
| State management | Simple object + event listeners | Overkill to use a library for 3 fields |

---

Status: COMPLETE
