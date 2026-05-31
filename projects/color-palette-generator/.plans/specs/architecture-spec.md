# Architecture Specification: Brutalist Design System Reskin

**Created by**: Architect
**Date**: 2026-03-16
**Template**: `static-html` (existing — no change)
**Status**: READY_FOR_REVIEW

---

## 1. Scope & Constraints

**What changes**: Visual design system only — CSS custom properties, typography, borders, corners, shadows, animations, and any HTML class/attribute changes needed to support the new aesthetic.

**What does NOT change**: All JavaScript logic (colors.js, app.js, clipboard.js), all functionality (generate, lock, copy, harmony modes, format toggle), HTML semantic structure, file organization, security measures (CSP, textContent usage, hex validation).

This is a **CSS-first reskin** with minimal HTML touch-ups (font links, minor class adjustments).

---

## 2. Design Language: Brutalism

### Core Principles

1. **Raw structure** — Expose the grid, the borders, the underlying geometry. Nothing is hidden behind soft gradients or blur.
2. **Bold typography** — Large, heavy, monospace or heavy grotesque type. Uppercase where it creates impact.
3. **High contrast** — Near-black and near-white. Stark, not warm. No cozy amber glow.
4. **Thick borders** — 3-4px solid black/white borders that are proudly visible, not subtle dividers.
5. **Zero decoration** — No rounded corners (0px radius everywhere). No box shadows. No backdrop-filter blur. No gradients.
6. **Intentionally unpolished** — Hard edges, blunt transitions, raw layout. The "rough concrete" of web design.

### Anti-Patterns (What Brutalism Is NOT)

- NOT ugly or broken — it's intentionally stark and honest
- NOT inaccessible — contrast and readability are actually enhanced
- NOT chaotic — the grid is rigid and the hierarchy is clear
- NOT anti-user — functionality and clarity are paramount

---

## 3. Design Tokens: Current vs. Brutalist

### Color Palette

| Token | Current (Warm Dark Craft) | Brutalist | Rationale |
|-------|--------------------------|-----------|-----------|
| `--color-primary` | `#D4913B` (amber) | `#000000` (black) | Stark, no decorative accent |
| `--color-primary-hover` | `#E0A04E` | `#333333` | Slight lift, still monochrome |
| `--color-primary-active` | `#C0832F` | `#000000` | Same as default — blunt, no state theatrics |
| `--color-primary-muted` | `rgba(212,145,59,0.15)` | `rgba(0,0,0,0.1)` | Subtle background tint |
| `--color-secondary` | `#8C7A6B` (warm stone) | `#666666` (neutral gray) | Cold, utilitarian |
| `--color-background` | `#1A1614` (warm obsidian) | `#FFFFFF` (pure white) | Brutalism favors stark white canvas |
| `--color-surface` | `#242019` | `#F0F0F0` (light gray) | Barely off-white for layering |
| `--color-surface-raised` | `#2E2A22` | `#E0E0E0` | Visible differentiation |
| `--color-border` | `rgba(255,255,255,0.08)` | `#000000` | Thick, visible, proud borders |
| `--color-border-hover` | `rgba(255,255,255,0.14)` | `#000000` | No change on hover — border is permanent |
| `--color-text-primary` | `#E8E0D6` (warm white) | `#000000` (pure black) | Maximum contrast on white |
| `--color-text-secondary` | `#9C9083` | `#444444` | Still readable, less prominent |
| `--color-text-tertiary` | `#6B6058` | `#888888` | Muted but not invisible |
| `--color-success` | `#5B9A6B` | `#000000` | Monochrome — success communicated by text, not color |
| `--color-error` | `#C75C4A` | `#FF0000` | Brutalist red — raw, unapologetic |
| `--color-warning` | `#C9943A` | `#000000` | Same rationale as success |
| `--color-info` | `#5A8FAD` | `#0000FF` | Brutalist blue — raw, primary color |

### Typography

| Token | Current | Brutalist | Rationale |
|-------|---------|-----------|-----------|
| `--font-ui` | `'Inter', system-ui, sans-serif` | `'Space Mono', 'Courier New', monospace` | Monospace is the brutalist signature. Space Mono has personality while remaining legible. |
| `--font-mono` | `'JetBrains Mono', monospace` | `'Space Mono', 'Courier New', monospace` | Same font for both — brutalism collapses the distinction between "decorative" and "functional" type. |

**Font sizing changes**:
- Toolbar title: Uppercase, bold, letter-spaced (`text-transform: uppercase; letter-spacing: 0.15em; font-weight: 700`)
- Swatch color codes: Larger, bolder (`font-size: 20px; font-weight: 700; text-transform: uppercase`)
- Buttons: Uppercase, letter-spaced
- Kbd hints: Same monospace, more prominent border

### Spacing

| Token | Current | Brutalist | Rationale |
|-------|---------|-----------|-----------|
| All spacing tokens | Same | **No change** | Spacing is functional, not decorative |

### Border Radius

| Token | Current | Brutalist | Rationale |
|-------|---------|-----------|-----------|
| `--radius-sm` | `6px` | `0px` | Sharp corners everywhere |
| `--radius-md` | `8px` | `0px` | No rounding |
| `--radius-lg` | `12px` | `0px` | No rounding |
| `--radius-xl` | `16px` | `0px` | No rounding |
| `--radius-full` | `9999px` | `0px` | No pills, no circles |

### Borders

| Element | Current | Brutalist |
|---------|---------|-----------|
| Toolbar border | `1px solid rgba(255,255,255,0.08)` | `3px solid #000000` |
| Button border | `none` or `1px solid rgba(...)` | `3px solid #000000` |
| Format toggle border | `1px solid var(--color-border)` | `3px solid #000000` |
| Harmony select border | `1px solid rgba(255,255,255,0.1)` | `3px solid #000000` |
| Toast border | `1px solid rgba(255,255,255,0.1)` | `3px solid #000000` |

### Shadows

| Token | Current | Brutalist |
|-------|---------|-----------|
| `--shadow-glow-primary` | `0 0 20px rgba(212,145,59,0.15)` | `none` — or optional hard-offset: `4px 4px 0px #000000` |
| `--shadow-elevation-md` | `0 4px 12px rgba(0,0,0,0.3)` | `none` |
| `--shadow-elevation-lg` | `0 8px 24px rgba(0,0,0,0.4)` | `none` |

**Hard-offset shadow** (optional brutalist flourish): `4px 4px 0px #000000` — a flat, geometric shadow with no blur. Used sparingly on the generate button and toast for depth without softness.

### Transitions & Animations

| Element | Current | Brutalist |
|---------|---------|-----------|
| `--transition-fast` | `150ms ease` | `0ms` or `50ms linear` — instant, no easing |
| `--transition-normal` | `200ms ease` | `0ms` or `50ms linear` |
| `--transition-color` | `250ms ease-out` | `0ms` — colors snap instantly, no smooth transition |
| Stagger delay | 80ms between swatches | `0ms` — all swatches change simultaneously |
| Copy pulse animation | scale bounce | `none` — no animation, text just changes |
| Lock bounce animation | scale bounce | `none` — icon just swaps |
| Button hover lift | `translateY(-1px)` | `none` — hover changes background/invert only |

**Rationale**: Brutalism rejects the notion of "easing" the user into changes. Things happen. Immediately. No cushioning.

---

## 4. Component Redesign

### Generate Button

```
Current:  Rounded amber button with glow hover
Brutalist: Black rectangle, white uppercase text, thick border, inverts on hover
```

```css
/* Default */
background: #000000;
color: #FFFFFF;
font-family: 'Space Mono', monospace;
font-weight: 700;
text-transform: uppercase;
letter-spacing: 0.1em;
padding: 12px 24px;
border: 3px solid #000000;
border-radius: 0;

/* Hover */
background: #FFFFFF;
color: #000000;
/* No shadow, no lift, no glow — just inversion */

/* Active */
background: #333333;
color: #FFFFFF;

/* Focus-visible */
outline: 3px solid #000000;
outline-offset: 3px;
```

### Toolbar

```
Current:  Frosted glass effect (backdrop-filter blur), warm dark, subtle border
Brutalist: Flat white/light gray, thick black top border, no blur, no transparency
```

```css
background: #FFFFFF;
border-top: 3px solid #000000;
backdrop-filter: none;
```

### Swatch Color Codes

```
Current:  JetBrains Mono 18px, weight 500, subtle
Brutalist: Space Mono 20px, weight 700, uppercase, prominent
```

### Lock Button

```
Current:  Fades in on hover, amber when locked, subtle
Brutalist: Always a stark presence, black/white based on swatch luminance, no fade — just visible or hidden via display
```

### Toast

```
Current:  Dark surface with warm tint, rounded, shadow, slide-up animation
Brutalist: White background, thick black border, hard-offset shadow, appears instantly (no slide)
```

```css
background: #FFFFFF;
color: #000000;
border: 3px solid #000000;
border-radius: 0;
box-shadow: 4px 4px 0px #000000;
/* No slide animation — opacity snap */
```

### Format Toggle

```
Current:  Warm surface background, rounded segments, amber active state
Brutalist: White background, thick border, active segment is inverted (black bg, white text)
```

### Harmony Select

```
Current:  Dark surface, subtle border, warm tones
Brutalist: White background, thick black border, monospace font, no custom arrow (or stark SVG arrow)
```

### Keyboard Hint (kbd)

```
Current:  Subtle, nearly invisible chip
Brutalist: Thick-bordered box, prominent monospace text, visible as an intentional UI element
```

### Locked Swatch Overlay

```
Current:  Subtle diagonal stripes at 3% opacity
Brutalist: Bold diagonal stripes at 8-10% opacity, thicker lines (4px stripe, 12px gap)
```

---

## 5. Focus & Accessibility Updates

Focus states move from amber ring to **thick black outline**:

```css
*:focus-visible {
  outline: 3px solid #000000;
  outline-offset: 3px;
}
```

On dark swatches, the focus ring should invert to white. This can be handled by adding a contrasting box-shadow as a secondary focus indicator:

```css
.swatch:focus-visible {
  outline: 3px solid #FFFFFF;
  outline-offset: -3px;
  box-shadow: inset 0 0 0 3px #000000;
}
```

**Contrast ratios remain AA-compliant or better**: Black text on white (#000 on #FFF) is 21:1. All brutalist token combinations exceed WCAG AA requirements.

The locked lock button color changes from `#D4913B` (amber) to using the contrast color (black or white depending on swatch luminance) — maintaining visibility without decorative color.

---

## 6. Swatch Hover Overlay

```
Current:  5% white overlay on dark swatches, 5% black overlay on light swatches
Brutalist: No overlay — or a stark 10% overlay with no transition. Hover is less important in brutalism; the grid speaks for itself.
```

**Decision**: Remove hover overlay entirely. The cursor change to `pointer` is sufficient feedback. This reduces visual noise and aligns with the "raw" aesthetic.

---

## 7. Files Changed

| File | Type of Change | Scope |
|------|---------------|-------|
| `css/styles.css` | **Major rewrite** | All design tokens, component styles, animations, responsive tweaks |
| `index.html` | **Minor edit** | Google Fonts link → Space Mono, remove Inter/JetBrains Mono references |
| `js/app.js` | **Minor edit** | Change hardcoded `#D4913B` lock color to use contrast color from `getContrastColor()` |
| `js/colors.js` | **No change** | |
| `js/clipboard.js` | **No change** | |

---

## 8. HTML Changes

### Font Link Update

```html
<!-- Current -->
<link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500&display=swap" rel="stylesheet">

<!-- Brutalist -->
<link href="https://fonts.googleapis.com/css2?family=Space+Mono:wght@400;700&display=swap" rel="stylesheet">
```

### CSP Update

Update font-src if needed (still fonts.gstatic.com — no change needed).

---

## 9. JS Changes

### app.js — Lock Button Color

```javascript
// Current (line ~112)
lockBtn.style.color = '#D4913B';

// Brutalist
lockBtn.style.color = contrastColor;
// Lock icon should match the auto-contrast text color, not a decorative accent
```

This is the ONLY JavaScript change needed. All other behavior remains identical.

---

## 10. CSS Change Summary

### Design Token Overrides (`:root` block)

Replace all custom property values in `:root`. No structural CSS changes — just token values + a few component-level overrides for borders/backgrounds.

### Component-Level Changes

1. **Remove** all `border-radius` (set to 0)
2. **Remove** all `box-shadow` (replace with hard-offset where noted)
3. **Remove** `backdrop-filter` and `-webkit-backdrop-filter` from toolbar
4. **Remove** all easing transitions (replace with instant or 50ms linear)
5. **Remove** stagger animation delays
6. **Remove** copy-pulse and lock-bounce keyframe animations
7. **Remove** swatch hover overlay (`.swatch::before`)
8. **Add** thick borders (3px solid) to toolbar, buttons, inputs, toast
9. **Add** `text-transform: uppercase` and `letter-spacing` to key text elements
10. **Add** hard-offset `box-shadow: 4px 4px 0 #000` to generate button and toast
11. **Update** `.format-btn.active` from amber muted to inverted black/white
12. **Update** `.swatch.locked::after` stripe pattern to bolder (thicker lines, higher opacity)
13. **Update** focus-visible outlines from amber to black, thicker offset

---

## 11. Trade-offs

| Optimizing For | Sacrificing |
|---------------|------------|
| Visual distinctiveness | The warm, approachable feel of the current design |
| High contrast / accessibility | Nuanced color states (amber active, warm glows) |
| Instant feedback (no transitions) | The "satisfying" staggered reveal animation |
| Design consistency (everything monochrome) | Visual hierarchy through color accent |
| Brutalist authenticity | Soft/premium aesthetic |

---

## 12. Delegation

**DELEGATE:**
- **senior_dev**: Update `js/app.js` — replace hardcoded amber `#D4913B` lock color with dynamic contrast color. Verify all functionality still works after CSS changes. Review that no JS references to old design tokens remain.
- **junior_dev**: Primary owner of the reskin. Rewrite `css/styles.css` design tokens and component styles per this spec. Update `index.html` font link. Ensure responsive breakpoints still work. Test all states (hover, focus, locked, active, toast).

---

Status: READY_FOR_REVIEW
