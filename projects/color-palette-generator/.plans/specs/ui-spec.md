# UI Specification: Color Palette Generator — Brutalist Redesign

Based on research in: .plans/research/ui-research.md
Styling wireframes from: .plans/specs/ux-spec.md
Replaces: "Warm Dark Craft" design system

---

## Design Philosophy

**"Raw Machine"** — A brutalist interface that strips away all decoration and exposes the structure. The tool looks like it was built by an engineer who respects the user's time more than their comfort. No rounded corners. No gradients. No glow effects. No subtle anything. Thick borders, monospace everything, high contrast, and aggressive typography. The generated palette colors remain the hero — but now they sit inside a raw, industrial frame that makes them pop harder.

The five principles shift from "Warm Dark Craft" to:
1. **Colors are still the hero** — Swatches dominate, but the frame is raw and structural
2. **Exposed, not hidden** — Borders are thick and visible. Structure is the decoration
3. **Monospace everything** — One typeface family. No typographic hierarchy tricks
4. **No polish** — Zero rounded corners, zero shadows, zero gradients, zero blur
5. **Instant, not smooth** — No easing, no transitions. State changes are immediate

---

## Design System

### Color Palette (UNIQUE — Brutalist Contrast)

**IMPORTANT**: Brutalism demands high contrast and limited palette. We use raw black, raw white, and a single harsh accent. No warm undertones, no soft neutrals. This is the opposite of every polished design tool.

| Name | Hex | RGB | Usage |
|------|-----|-----|-------|
| Primary (Electric Yellow) | #EBFF00 | rgb(235, 255, 0) | Main CTAs, active states, focus rings |
| Primary Hover | #D4E600 | rgb(212, 230, 0) | Hover state on primary elements |
| Secondary (Raw White) | #FFFFFF | rgb(255, 255, 255) | Secondary actions, text on dark |
| Background | #0A0A0A | rgb(10, 10, 10) | App background, deepest layer |
| Surface | #141414 | rgb(20, 20, 20) | Toolbar, panels |
| Surface Raised | #1E1E1E | rgb(30, 30, 30) | Dropdowns, toasts |
| Border | #FFFFFF | rgb(255, 255, 255) | All visible borders — thick, no transparency |
| Border Muted | #333333 | rgb(51, 51, 51) | Swatch dividers (between color columns) |
| Text Primary | #FFFFFF | rgb(255, 255, 255) | All primary text |
| Text Secondary | #888888 | rgb(136, 136, 136) | Muted labels, hints |
| Text Tertiary | #555555 | rgb(85, 85, 85) | Disabled states |
| Success | #00FF00 | rgb(0, 255, 0) | Success — raw, pure green |
| Error | #FF0000 | rgb(255, 0, 0) | Error — raw, pure red |
| Warning | #FFAA00 | rgb(255, 170, 0) | Warning states |
| Info | #00AAFF | rgb(0, 170, 255) | Informational |

**Why These Colors**: Brutalism rejects nuance. Black (#0A0A0A) is near-pure black — no warm or cool undertone. White is #FFFFFF — no cream, no off-white. The accent is Electric Yellow (#EBFF00) — an industrial, high-visibility color borrowed from caution tape and construction signage. It's aggressively NOT the amber/coral/teal of polished design tools. Success/Error use pure saturated values because brutalism doesn't soften semantic colors.

### Typography

**Font Family**: Space Mono (primary — everything), system monospace (fallback)

Rationale: Brutalism demands monospace. Space Mono by Colophon Foundry (available on Google Fonts) has the right amount of industrial character — wider than Courier, more geometric, designed for screens. Using one font for EVERYTHING (headers, body, buttons, codes) is intentionally monotone and raw. No typographic hierarchy beyond size and weight.

**Font Loading**: Load from Google Fonts CDN with `font-display: swap`. Fallback stack: `'Courier New', Courier, monospace`.

| Element | Font | Size | Weight | Line Height | Letter Spacing | Transform |
|---------|------|------|--------|-------------|----------------|-----------|
| H1 (Page Title) | Space Mono | 28px | 700 | 1.1 | 0.08em | uppercase |
| H2 (Section) | Space Mono | 20px | 700 | 1.2 | 0.06em | uppercase |
| H3 (Label) | Space Mono | 16px | 700 | 1.3 | 0.04em | uppercase |
| Body | Space Mono | 14px | 400 | 1.5 | 0.02em | none |
| Small | Space Mono | 12px | 400 | 1.4 | 0.02em | none |
| Caption | Space Mono | 11px | 400 | 1.3 | 0.04em | uppercase |
| Color Code (Hex) | Space Mono | 18px | 700 | 1.2 | 0.06em | uppercase |
| Color Code Small | Space Mono | 14px | 700 | 1.2 | 0.04em | uppercase |
| Button Label | Space Mono | 13px | 700 | 1 | 0.06em | uppercase |
| Keyboard Shortcut | Space Mono | 11px | 700 | 1 | 0.04em | uppercase |

**Key difference from previous**: ALL text is monospace. ALL headings and buttons are uppercase. Letter-spacing is wider across the board to feel mechanical. Weight contrast (400 vs 700) is the only typographic hierarchy.

### Spacing Scale

Base unit: 4px. Brutalism uses generous spacing to let raw elements breathe.

| Token | Value | Usage |
|-------|-------|-------|
| xs | 4px | Tight internal gaps |
| sm | 8px | Icon-to-text, small padding |
| md | 16px | Standard padding |
| lg | 24px | Section spacing, toolbar padding |
| xl | 32px | Major gaps |
| 2xl | 48px | Page margins |
| 3xl | 64px | Desktop hero spacing |

### Border Radius

| Token | Value | Usage |
|-------|-------|-------|
| none | 0px | **Everything** |

Brutalism uses zero border radius on all elements. No exceptions. Buttons are rectangles. Inputs are rectangles. Cards are rectangles. Toasts are rectangles. This is non-negotiable.

### Borders

Borders are THE decorative element in brutalism. They replace shadows, glows, and subtle separators.

| Token | Value | Usage |
|-------|-------|-------|
| border-thick | 3px solid #FFFFFF | Primary structural borders (toolbar edge, cards) |
| border-medium | 2px solid #FFFFFF | Interactive element borders (buttons, inputs, selects) |
| border-thin | 1px solid #333333 | Subtle dividers (swatch separators) |
| border-accent | 3px solid #EBFF00 | Focus states, active/selected elements |

### Shadows

| Token | Value | Usage |
|-------|-------|-------|
| **none** | none | **Everything** |

Brutalism rejects shadows entirely. Elevation is communicated through borders and z-index stacking, not drop shadows or blurs. No `box-shadow`. No `filter: drop-shadow`. No `backdrop-filter: blur`.

### Iconography

- **Icon Set**: Lucide Icons (same set, but rendered at heavier stroke weight)
- **Stroke Width**: 2.5px (heavier than the previous 1.5px for a bolder feel)
- **Color**: #FFFFFF default, #EBFF00 on hover/active
- **Size**: 20px default, 24px prominent
- **Style**: Strictly line icons, no fills. The thick stroke makes them feel industrial.

---

## Components

### Generate Button (Primary CTA)

The most important button — bold, impossible to miss, industrial.

```css
/* Default */
background: #EBFF00;
color: #0A0A0A;
font: 700 13px 'Space Mono', monospace;
letter-spacing: 0.06em;
text-transform: uppercase;
padding: 14px 28px;
border: 3px solid #EBFF00;
border-radius: 0;
cursor: pointer;
transition: none;
min-height: 44px;

/* Hover */
background: #0A0A0A;
color: #EBFF00;
border-color: #EBFF00;
/* Hard color swap — no transition, no ease, instant */

/* Active */
background: #EBFF00;
color: #0A0A0A;
border-color: #FFFFFF;

/* Focus-visible */
outline: 3px solid #EBFF00;
outline-offset: 3px;

/* Disabled / Muted */
opacity: 0.3;
cursor: not-allowed;
```

Keyboard hint `<kbd>` appears next to button text with border treatment.

### Secondary Button (Copy All, etc.)

```css
/* Default */
background: transparent;
color: #FFFFFF;
font: 700 13px 'Space Mono', monospace;
letter-spacing: 0.06em;
text-transform: uppercase;
border: 2px solid #FFFFFF;
border-radius: 0;
padding: 10px 20px;
cursor: pointer;
transition: none;
min-height: 44px;

/* Hover */
background: #FFFFFF;
color: #0A0A0A;
/* Instant inversion — no transition */

/* Active */
background: #888888;
color: #0A0A0A;
```

### Ghost/Icon Button (Lock, Copy on swatch)

```css
/* Default */
background: transparent;
color: inherit;
border: 2px solid transparent;
border-radius: 0;
padding: 8px;
cursor: pointer;
transition: none;
min-width: 44px;
min-height: 44px;

/* Hover */
border-color: currentColor;
/* Just adds a visible border — no color change, no background */

/* Active (e.g., lock is engaged) */
color: #EBFF00;
border-color: #EBFF00;
```

### Color Swatch (Core Element)

Each swatch is a full-height column (desktop) or full-width bar (mobile). Brutalist treatment: separated by hard 1px dark borders, no overlap, no gap.

```css
/* Swatch Container */
position: relative;
display: flex;
flex-direction: column;
align-items: center;
justify-content: center;
min-height: 100%;
cursor: pointer;
border-right: 1px solid #333333; /* visible divider between swatches */
transition: none; /* NO smooth color transitions — instant change */

/* Last swatch has no right border */
&:last-child { border-right: none; }

/* Mobile: border-bottom instead of border-right */
@media (max-width: 767px) {
  border-right: none;
  border-bottom: 1px solid #333333;
  &:last-child { border-bottom: none; }
}

/* Hex Code Display */
font: 700 18px 'Space Mono', monospace;
letter-spacing: 0.06em;
text-transform: uppercase;
/* Auto light/dark text based on luminance — same logic as before */
```

**No hover overlay**. No brightness shift. No ::before pseudo-element overlay. When you hover a swatch, the action buttons appear — that's it. The swatch color stays pure and unmodified.

**No staggered animation**. When generating, all 5 swatches change color simultaneously and instantly. No wave effect, no delays, no transition duration.

**Lock indicator**: When locked, display a thick top border (3px solid) in the swatch's contrasting text color. No diagonal stripes, no overlay patterns.

### Swatch Action Bar

Appears on hover. Same icon buttons as before but with brutalist treatment.

| Action | Icon | Behavior |
|--------|------|----------|
| Copy hex | `copy` | Copy to clipboard, icon swaps to `check` for 1.5s |
| Lock/Unlock | `lock` / `unlock` | Toggle lock state |

### Toast / Notification

```css
/* Container */
position: fixed;
bottom: calc(56px + 16px); /* toolbar height + spacing */
left: 50%;
transform: translateX(-50%);
background: #1E1E1E;
border: 2px solid #FFFFFF;
border-radius: 0;
padding: 12px 20px;
font: 700 13px 'Space Mono', monospace;
letter-spacing: 0.04em;
text-transform: uppercase;
color: #FFFFFF;

/* No shadow. No blur. */
box-shadow: none;

/* Animation: instant appear, instant disappear after 1.5s */
/* No slide, no fade — just display:block / display:none or opacity toggle */
```

Content format: `COPIED #FF5733` — uppercase, monospace, direct.

### Toolbar / Controls Bar

```css
/* Container */
position: relative; /* NOT sticky with backdrop-blur — brutalism rejects blur */
display: flex;
align-items: center;
justify-content: space-between;
padding: 0 24px;
background: #141414; /* Solid, opaque — no transparency */
border-top: 3px solid #FFFFFF; /* THICK visible border — the key brutalist element */
min-height: 56px;
z-index: 10;
gap: 16px;
```

Layout:
- **Left**: App title "PALETTE" in Space Mono 700, 16px, uppercase, letter-spacing 0.08em
- **Center**: Generate button with `[SPACE]` kbd hint
- **Right**: Harmony mode select, format toggle, copy-all button

### Harmony Mode Selector

```css
/* Select element */
background: #141414;
border: 2px solid #FFFFFF;
border-radius: 0;
padding: 8px 36px 8px 12px;
color: #FFFFFF;
font: 700 13px 'Space Mono', monospace;
letter-spacing: 0.04em;
text-transform: uppercase;
cursor: pointer;
transition: none;
min-height: 44px;
appearance: none;
-webkit-appearance: none;
/* Custom dropdown arrow: white chevron SVG */
background-image: url("data:image/svg+xml,..."); /* white arrow */
background-repeat: no-repeat;
background-position: right 12px center;

/* Hover */
border-color: #EBFF00;
color: #EBFF00;

/* Focus-visible */
outline: 3px solid #EBFF00;
outline-offset: 2px;

/* Options */
option {
  background: #1E1E1E;
  color: #FFFFFF;
}
```

### Color Format Toggle

```css
/* Container */
display: inline-flex;
background: #141414;
border: 2px solid #FFFFFF;
border-radius: 0;
padding: 0;

/* Segment button */
padding: 8px 12px;
font: 700 12px 'Space Mono', monospace;
letter-spacing: 0.04em;
text-transform: uppercase;
color: #888888;
background: transparent;
border: none;
border-right: 1px solid #333333; /* internal divider */
cursor: pointer;
transition: none;
min-height: 44px;
min-width: 44px;

/* Last segment — no right border */
&:last-child { border-right: none; }

/* Hover */
color: #FFFFFF;

/* Active Segment */
background: #FFFFFF;
color: #0A0A0A;
/* Hard fill — no muted backgrounds, no accent tints */
```

### Keyboard Shortcut Hint (kbd)

```css
display: inline-flex;
align-items: center;
padding: 2px 8px;
font: 700 11px 'Space Mono', monospace;
letter-spacing: 0.04em;
text-transform: uppercase;
color: #555555;
background: transparent;
border: 2px solid #555555;
border-radius: 0;
```

---

## Layout

### Desktop (> 1024px)

```
┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃           │           │           │           │              ┃
┃           │           │           │           │              ┃
┃   COLOR   │   COLOR   │   COLOR   │   COLOR   │   COLOR      ┃
┃     1     │     2     │     3     │     4     │     5        ┃
┃           │           │           │           │              ┃
┃  #A45C2F  │  #D4913B  │  #E8C97A  │  #5B9A6B  │ #2E4A3F     ┃
┃  [🔓] [📋]│  [🔓] [📋]│  [🔓] [📋]│  [🔓] [📋]│ [🔓] [📋]   ┃
┃           │           │           │           │              ┃
┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫  ← 3px solid white
┃ PALETTE    [■ GENERATE  SPACE]    RANDOM▾  HEX|RGB|HSL  COPY┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
```

- Swatches: `display: grid; grid-template-columns: repeat(5, 1fr); height: calc(100vh - 56px);`
- Swatch dividers: 1px solid #333333 between columns
- Toolbar: 3px solid white top border, solid #141414 background
- No backdrop blur. No transparency. No rounded anything.

### Mobile (< 640px)

```
┏━━━━━━━━━━━━━━━━━━━━┓
┃ PALETTE    [■]  [≡] ┃
┣━━━━━━━━━━━━━━━━━━━━┫ ← 3px solid white
┃  COLOR 1   #A45C2F  ┃
┃──────────────────── ┃ ← 1px solid #333
┃  COLOR 2   #D4913B  ┃
┃──────────────────── ┃
┃  COLOR 3   #E8C97A  ┃
┃──────────────────── ┃
┃  COLOR 4   #5B9A6B  ┃
┃──────────────────── ┃
┃  COLOR 5   #2E4A3F  ┃
┗━━━━━━━━━━━━━━━━━━━━┛
```

- Swatches stack vertically as full-width horizontal bars
- Each swatch: `min-height: calc((100vh - toolbar-height) / 5)`
- Lock and copy always visible on mobile
- Toolbar at top on mobile with compact controls

### Responsive Breakpoints

| Name | Width | Layout Change |
|------|-------|---------------|
| mobile | < 640px | Vertical stacked swatches, compact toolbar at top |
| tablet | 640–1024px | Horizontal swatches, condensed toolbar |
| desktop | > 1024px | Full horizontal layout, spacious toolbar at bottom |

---

## Animations & Micro-Interactions

Brutalism rejects smooth transitions. State changes are **instant**.

### Palette Generation

- **Trigger**: Spacebar or Generate button
- **Animation**: **NONE**. All 5 swatches change color simultaneously and instantly. No stagger, no delay, no transition-duration. `transition: none` on background-color.
- **Why**: Brutalism values immediacy. The jarring instant-switch IS the feedback. Users know something happened because the colors snapped.

### Copy to Clipboard

1. Click hex code → icon swaps from `copy` to `check` instantly (no morph, no crossfade)
2. No pulse animation on hex text
3. Toast appears instantly (no slide-up, no fade-in) — `opacity: 0` to `opacity: 1` with `transition: none`
4. Toast disappears after 1.5s — instant removal

### Lock Toggle

1. Click lock → icon swaps between `lock` and `unlock` instantly (no bounce, no scale animation)
2. When locked: thick top border appears on swatch (3px solid in contrasting text color)
3. No diagonal stripe overlay. The border IS the indicator.

### Hover States

- Swatch hover: action buttons appear instantly (opacity 0 → 1, no transition)
- Button hover: colors invert instantly (no ease, no duration)
- Generate button hover: bg/color swap instantly
- **No** `transform: translateY` lifts. **No** scale changes. **No** glow effects.

### General Transition Default

```css
/* The brutalist default: no transitions */
transition: none;

/* Only exception: toast auto-dismiss uses a simple opacity toggle */
```

### Reduced Motion

Already handled — since we use no transitions by default, `prefers-reduced-motion` is satisfied automatically. Include the media query as a no-op for completeness:

```css
@media (prefers-reduced-motion: reduce) {
  /* Already brutalist — no transitions to reduce */
}
```

---

## Accessibility

### Contrast Ratios (WCAG AA Compliance)

| Element | Foreground | Background | Ratio | Pass? |
|---------|-----------|------------|-------|-------|
| Body text | #FFFFFF | #0A0A0A | 19.3:1 | AAA |
| Secondary text | #888888 | #0A0A0A | 5.9:1 | AA |
| Primary accent | #EBFF00 | #0A0A0A | 16.2:1 | AAA |
| Button text | #0A0A0A | #EBFF00 | 16.2:1 | AAA |
| Swatch auto-text | Auto white/dark | Swatch color | ≥ 4.5:1 | AA |
| Toolbar text | #FFFFFF | #141414 | 17.1:1 | AAA |

Brutalism's high-contrast nature makes WCAG compliance easier than softer designs.

### Keyboard Navigation

Same as UX spec — no changes:

| Key | Action |
|-----|--------|
| `Space` | Generate new palette |
| `Tab` | Move focus between swatches and toolbar controls |
| `Enter` | Activate focused element |
| `L` | Lock/unlock focused swatch |
| `C` | Copy focused swatch color |
| `1-5` | Focus swatch directly |
| `Escape` | Close open dropdown |

### Focus States

```css
/* Brutalist focus: thick yellow outline, no subtle rings */
*:focus-visible {
  outline: 3px solid #EBFF00;
  outline-offset: 3px;
}

*:focus:not(:focus-visible) {
  outline: none;
}
```

### ARIA & Screen Reader

No changes from UX spec. Same roles, labels, and live regions.

### Touch Targets

All interactive elements: minimum 44x44px. Same as before.

---

## Visual States Summary

### Swatch States

| State | Visual Treatment |
|-------|-----------------|
| Default | Full-color background, hex code centered, actions hidden, 1px #333 border between swatches |
| Hover | Actions appear instantly (no fade), NO brightness overlay |
| Locked | 3px thick top border in contrasting text color, lock icon visible |
| Locked + Hover | Same as locked, actions visible |
| Focused (keyboard) | 3px yellow outline, actions visible |

### Button States (all types)

| State | Visual Treatment |
|-------|-----------------|
| Default | Per component spec |
| Hover | Instant color inversion (bg ↔ text swap) |
| Active | Slightly dimmed version of hover |
| Focus-visible | 3px solid #EBFF00 outline, 3px offset |
| Disabled | 30% opacity, cursor: not-allowed |

---

## CSS Custom Properties (Design Tokens)

Complete mapping for implementation:

```css
:root {
  /* Colors */
  --color-primary: #EBFF00;
  --color-primary-hover: #D4E600;
  --color-secondary: #FFFFFF;
  --color-background: #0A0A0A;
  --color-surface: #141414;
  --color-surface-raised: #1E1E1E;
  --color-border: #FFFFFF;
  --color-border-muted: #333333;
  --color-text-primary: #FFFFFF;
  --color-text-secondary: #888888;
  --color-text-tertiary: #555555;
  --color-success: #00FF00;
  --color-error: #FF0000;
  --color-warning: #FFAA00;
  --color-info: #00AAFF;

  /* Typography */
  --font-mono: 'Space Mono', 'Courier New', Courier, monospace;
  /* No --font-ui. Mono is the only font. */

  /* Spacing */
  --space-xs: 4px;
  --space-sm: 8px;
  --space-md: 16px;
  --space-lg: 24px;
  --space-xl: 32px;
  --space-2xl: 48px;
  --space-3xl: 64px;

  /* Layout */
  --toolbar-height: 56px;
  --touch-target: 44px;

  /* Borders */
  --border-thick: 3px solid #FFFFFF;
  --border-medium: 2px solid #FFFFFF;
  --border-thin: 1px solid #333333;
  --border-accent: 3px solid #EBFF00;
  --radius: 0px; /* One token, one value, everywhere */

  /* No transitions */
  --transition: none;

  /* No shadows */
  /* (intentionally absent — no shadow tokens defined) */
}
```

---

## Key Differences from "Warm Dark Craft"

| Aspect | Warm Dark Craft (Previous) | Raw Machine (Brutalist) |
|--------|---------------------------|------------------------|
| Font | Geist Sans + Geist Mono | Space Mono only |
| Border radius | 6–16px | 0px everywhere |
| Borders | 1px rgba subtle | 2–3px solid white, visible |
| Shadows | Warm glows, elevation | None |
| Backdrop blur | Yes (toolbar) | No |
| Transitions | 150–300ms ease | None (instant) |
| Accent color | Amber #D4913B | Electric Yellow #EBFF00 |
| Background | Warm obsidian #1A1614 | Pure black #0A0A0A |
| Text color | Warm white #E8E0D6 | Pure white #FFFFFF |
| Text transform | Mixed case | Uppercase on headings/buttons |
| Letter spacing | Tight (-0.02 to 0.04em) | Wide (0.04–0.08em) |
| Hover pattern | Subtle glow + lift | Instant color inversion |
| Generation anim | Staggered wave, 80ms delay | Instant, simultaneous |
| Lock indicator | Diagonal stripes at 3% | Thick top border |
| Overall mood | Creative studio warmth | Industrial machine |

---

Status: READY_FOR_REVIEW
