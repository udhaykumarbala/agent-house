# UX Specification: Random Color Palette Generator — Brutalist Reskin

Based on research in: .plans/research/ux-research.md
Implementing features from: .plans/specs/product-spec.md
UI direction: **Brutalist** (replacing "Warm Dark Craft")
Security constraints from: .plans/research/security-expert.md

---

## Design System Change: Warm Dark Craft → Brutalism

This is a **skin-only change**. All functionality, user flows, keyboard shortcuts, accessibility behavior, and interaction logic remain identical to the existing spec. What changes:

| Aspect | Before (Warm Dark Craft) | After (Brutalist) |
|--------|--------------------------|-------------------|
| Typography | Geist Sans + JetBrains Mono | **Monospace everywhere** — a single monospace stack (`'Space Mono', 'JetBrains Mono', 'Courier New', monospace`) for all text. Bold weights. Uppercase headings. |
| Colors | Warm obsidian (#1A1614) + amber (#D4913B) accents | **Stark black (#000000) + pure white (#FFFFFF)**. No warm tints. Accent: raw yellow (#FFFF00) or none. |
| Borders | Subtle 1px rgba warm-tinted | **Thick 3-4px solid black or white**. Visible, structural, intentional. |
| Corners | Rounded (6-12px radius) | **0px — all sharp corners**. No border-radius anywhere. |
| Shadows | Warm-tinted glows, soft elevation | **None**. Zero box-shadow. Depth conveyed through borders and offset, not shadows. |
| Spacing | Generous, breathing, premium feel | **Tight and dense** where possible, **large and jarring** where intentional. Asymmetric. |
| Animations | Smooth 200-300ms ease transitions | **Instant or intentionally abrupt**. No easing curves. `step-end` timing or 0ms transitions. Colors snap, don't fade. |
| Icons | Lucide line icons, 1.5px stroke | **Text labels or raw unicode glyphs**. Minimal SVG. Prefer `[LOCK]`, `[COPY]`, `[GEN]` over icons. |
| Overall feel | "Creative studio at night" | **"Raw construction site"** — exposed structure, nothing hidden, function over form. |

---

## Screen List

1. **Main Generator** — Single-page app. Full-viewport palette swatches with brutalist controls. Same as before.
2. **Toast/Notification** — Copy confirmation. Raw text, no decoration.

No routing, no separate pages. This is the same single-page tool.

---

## Wireframes

### Screen 1: Main Generator (Desktop — Brutalist)

```
╔═══════════════════════════════════════════════════════════════════════════╗
║  PALETTE GENERATOR          [RANDOM ▾]    [ GENERATE ]  [SPACE]         ║
║                                                                         ║
╠══════════════╦══════════════╦══════════════╦══════════════╦══════════════╣
║              ║              ║              ║              ║              ║
║              ║              ║              ║              ║              ║
║              ║              ║              ║              ║              ║
║   COLOR 1    ║   COLOR 2    ║   COLOR 3    ║   COLOR 4    ║   COLOR 5   ║
║              ║              ║              ║              ║              ║
║              ║              ║              ║              ║              ║
║  [LOCK]      ║  [LOCK]      ║  [LOCK]      ║  [LOCK]      ║  [LOCK]     ║
║  #3A86FF     ║  #FF006E     ║  #FB5607     ║  #FFBE0B     ║  #8338EC    ║
║  [COPY]      ║  [COPY]      ║  [COPY]      ║  [COPY]      ║  [COPY]     ║
║              ║              ║              ║              ║              ║
╠══════════════╩══════════════╩══════════════╩══════════════╩══════════════╣
║  [HEX] [RGB] [HSL]                                      [COPY ALL]     ║
╚═══════════════════════════════════════════════════════════════════════════╝
```

**Key Brutalist Changes from Previous Wireframe**:
- **Double/thick borders** (`═══`) between all elements — structural, visible grid lines. 3-4px solid borders in implementation.
- **ALL CAPS text labels** for buttons and headings — brutalist convention.
- **Text-based buttons** instead of icon buttons — `[LOCK]`, `[COPY]`, `[GEN]` replace SVG icons. Raw, exposed affordances.
- **No hover-to-reveal** — all actions visible at all times on desktop too (not just mobile). Nothing is hidden. Brutalism exposes everything.
- **Sharp corners** everywhere — no border-radius. Every element is a rectangle.
- **Monospace font** for everything — the hex codes, the buttons, the title. One typeface, one voice.
- **No decorative elements** — no glow, no gradient, no subtle overlays. Just flat color, thick borders, raw text.

### Screen 1: Main Generator (Desktop — Hover/Active State)

```
╔═══════════════════════════════════════════════════════════════════════════╗
║  PALETTE GENERATOR          [RANDOM ▾]    [>GENERATE<]  [SPACE]         ║
╠══════════════╦══════════════╦██████████████╦══════════════╦══════════════╣
║              ║              ║▓▓▓▓▓▓▓▓▓▓▓▓▓▓║              ║              ║
║              ║              ║▓▓ HOVERED  ▓▓║              ║              ║
║              ║              ║▓▓           ▓▓║              ║              ║
║              ║              ║▓▓ [LOCK]   ▓▓║              ║              ║
║              ║              ║▓▓           ▓▓║              ║              ║
║              ║              ║▓▓ #FB5607   ▓▓║              ║              ║
║              ║              ║▓▓ [COPY]   ▓▓║              ║              ║
║              ║              ║▓▓▓▓▓▓▓▓▓▓▓▓▓▓║              ║              ║
╠══════════════╩══════════════╩══════════════╩══════════════╩══════════════╣
║  [HEX] [RGB] [HSL]                                      [COPY ALL]     ║
╚═══════════════════════════════════════════════════════════════════════════╝
```

**Hover behavior (Brutalist)**:
- No smooth brightness overlay. Instead: **inverted colors** — swatch background inverts, text inverts. Hard, instant switch (no transition).
- Or alternatively: a thick **dashed border** appears inside the hovered swatch as the only indicator.
- Actions remain visible at all times — hover only adds the visual indicator, it does not reveal hidden controls.

### Screen 1: Main Generator (Mobile — Portrait, Brutalist)

```
╔═══════════════════════════╗
║ PALETTE GENERATOR         ║
║ [RANDOM ▾]  [GENERATE]    ║
╠═══════════════════════════╣
║       COLOR 1             ║
║  [LOCK]  #3A86FF  [COPY]  ║
╠═══════════════════════════╣
║       COLOR 2             ║
║  [LOCK]  #FF006E  [COPY]  ║
╠═══════════════════════════╣
║       COLOR 3             ║
║  [LOCK]  #FB5607  [COPY]  ║
╠═══════════════════════════╣
║       COLOR 4             ║
║  [LOCK]  #FFBE0B  [COPY]  ║
╠═══════════════════════════╣
║       COLOR 5             ║
║  [LOCK]  #8338EC  [COPY]  ║
╠═══════════════════════════╣
║ [HEX][RGB][HSL] [COPY ALL]║
╚═══════════════════════════╝
```

**Mobile Brutalist Adaptations**:
- Same thick borders between every swatch — structural grid visible
- All controls always visible (same as desktop — no hidden states)
- Text labels instead of icons — readable at small sizes
- Full-width stacked bars, same as before functionally

### Toast Notification (Brutalist)

```
┌────────────────────────────────┐
│  COPIED: #3A86FF               │
└────────────────────────────────┘
      (black bg, white text, thick border, no animation — appears/disappears instantly)
```

**Brutalist toast behavior**:
- **No slide animation**. Appears instantly (`display: block`), disappears instantly after 1.5s.
- **Solid black background, white monospace text, 3px white border**.
- Positioned bottom-center, same as before.
- No rounded corners, no shadows, no warm tints.

---

## User Flows

**All user flows are IDENTICAL to the previous spec.** Brutalism changes the visual skin, not the behavior. Copied here for completeness with brutalist-specific visual notes:

### Primary Flow: Generate and Copy a Color

1. **User opens the app** → Palette auto-generates with 5 colors. Swatches fill the viewport. **Brutalist note**: colors snap in instantly (no staggered wave animation).
2. **User evaluates the palette** → Scans colors visually (<1 second).
3. **User presses spacebar** (or clicks GENERATE) → New palette generates. **Brutalist note**: colors change instantly, no transition. Hard cut.
4. **User repeats steps 2-3** → Rapid iteration. The lack of transition actually makes this faster.
5. **User finds a color they like** → Clicks the hex code on the swatch.
6. **Color copied to clipboard** → Toast appears instantly: "COPIED: #3A86FF!" Disappears after 1.5s.
7. **Task complete** → User pastes into their project.

### Secondary Flow: Refine with Lock

1. **User generates palettes** → Finds colors they like but not the full set.
2. **User clicks [LOCK] text button on liked colors** → Text changes to `[LOCKED]` with a visual indicator (e.g., cross-hatch pattern overlay on locked swatch at 5% opacity, or a thick dashed inner border).
3. **User presses spacebar** → Only unlocked colors regenerate instantly.
4. **User repeats** → Converges on a full palette.
5. **User copies individual colors** or uses **[COPY ALL]**.

### Tertiary Flow: Change Harmony Mode

1. User clicks harmony mode dropdown → Raw `<select>` element (no custom dropdown styling — brutalism embraces native browser controls).
2. Selects a mode.
3. Palette regenerates with the selected harmony rule.

### Tertiary Flow: Switch Color Format

1. User clicks [RGB] or [HSL] in the format toggle.
2. All swatch labels update instantly.
3. Copy actions now copy the active format.

### Tertiary Flow: Export Full Palette

1. User clicks [COPY ALL] → Copies all 5 hex codes as comma-separated list.
2. Toast confirms: "COPIED: ALL COLORS"

---

## Interaction Patterns

| Action | Element | Response | Brutalist Feedback |
|--------|---------|----------|-------------------|
| Press Spacebar | Global | Generate new palette (respects locks) | Colors snap instantly — no transition |
| Click | [GENERATE] button | Generate new palette (respects locks) | Button inverts colors momentarily (black→white, white→black) |
| Click | Hex code on swatch | Copy to clipboard | Toast: "COPIED: #FF5733" — instant appear/disappear |
| Click | [LOCK] text button | Toggle lock state | Text changes: `[LOCK]` ↔ `[LOCKED]`. Cross-hatch overlay appears/disappears on swatch |
| Click | Harmony dropdown | Open native select | Browser-native dropdown — no custom styling |
| Click | Format toggle [HEX]/[RGB]/[HSL] | Switch format | Active button gets inverted colors (white bg, black text). Instant swap. |
| Click | [COPY ALL] | Copy all colors | Toast: "COPIED: ALL COLORS" |
| Hover | Color swatch (desktop) | Visual indicator only | Thick dashed inner border or inverted text — instant, no fade |
| Keyboard L | While hovering/focused swatch | Toggle lock | Same as click [LOCK] |
| Keyboard C | While hovering/focused swatch | Copy color | Same as click hex code |

---

## States

### Initial Load State
**When**: App first opens
**Show**: Auto-generated 5-color palette. First-time visitors see a raw text hint at top: `PRESS SPACEBAR TO GENERATE. CLICK HEX TO COPY.` — no tooltip, no popover, just a line of text. Dismisses on first generation.
**Brutalist note**: No fancy onboarding. Just tell them.

### Active/Default State
**When**: User has generated at least one palette
**Show**: 5 full-viewport swatches with hex codes, [LOCK] and [COPY] text buttons visible at all times. No hidden states.

### Locked Colors State
**When**: One or more colors are locked
**Show**: Locked swatches display `[LOCKED]` instead of `[LOCK]`. Visual indicator: **cross-hatch pattern overlay** (`repeating-linear-gradient` at 5% opacity) OR thick **dashed inner border** (3px dashed, semi-transparent). Clear and unambiguous.
**Behavior**: Spacebar/Generate only randomizes unlocked colors.

### All Colors Locked State
**When**: All 5 locked
**Show**: GENERATE button text changes to `[ALL LOCKED]`. If user presses spacebar, inline text appears below toolbar: `UNLOCK A COLOR FIRST.` — raw text, no modal, no tooltip, disappears after 2s.

### Copy Success State
**When**: User copies a color or palette
**Show**: Toast at bottom-center: `COPIED: [VALUE]` — black background, white text, thick border. Instant appear, instant disappear after 1.5s.

### Copy Failure State
**When**: Clipboard API fails
**Show**: Toast changes to: `SELECT AND COPY: [VALUE]` — same brutalist styling. Value displayed as selectable text.

### Error State
**When**: None expected (all client-side)
**Fallback**: If localStorage unavailable, hint text shows every visit. Core functionality unaffected.

---

## Accessibility Considerations

Brutalism must NOT compromise accessibility. The high-contrast, text-heavy nature of brutalist design is actually **beneficial** for accessibility in many ways.

### Keyboard Navigation
- **Same tab order** as previous spec: Generate → Harmony dropdown → Swatch 1 (lock, copy) → Swatch 2 → ... → Swatch 5 → Format toggle → Copy All
- **Spacebar**: When focus is NOT on an interactive element, generates palette. When on a button, activates it.
- **L key**: Toggle lock on focused/hovered swatch
- **C key**: Copy focused/hovered swatch color
- **1-5 keys**: Focus swatch directly

### Screen Reader Support
- Same ARIA labels as previous spec
- Swatches: `role="listitem"` within `role="list"`
- Lock buttons: `aria-pressed="true/false"`, `aria-label="Lock color #FF5733"`
- Copy buttons: `aria-label="Copy color #FF5733 to clipboard"`
- Generate: `aria-label="Generate new color palette"`
- Toast: `role="status"`, `aria-live="polite"`

### Touch Targets
- All interactive elements: minimum 44x44px
- Text-based buttons ([LOCK], [COPY]) must have sufficient padding to meet touch target size
- Buttons spaced at least 8px apart

### Color & Contrast
- **App chrome**: Pure white (#FFFFFF) text on pure black (#000000) background = **21:1 contrast ratio** — exceeds WCAG AAA
- **Swatch text**: Same auto-contrast logic: `luminance = 0.299*R + 0.587*G + 0.114*B; textColor = luminance > 150 ? '#000000' : '#FFFFFF'`
- **Focus indicators**: 3px solid white outline (or 3px solid black on light elements). Thicker than before — more visible, more brutalist.
- Brutalist high-contrast aesthetic inherently meets and exceeds WCAG AA requirements

### Reduced Motion
- Brutalism already minimizes animation. But still respect `prefers-reduced-motion: reduce` — ensure any remaining transitions (if any) are set to 0ms.
- The default brutalist behavior (instant color changes) is already reduced-motion friendly.

---

## Responsive Breakpoints

| Breakpoint | Layout | Swatch Orientation | Brutalist Notes |
|------------|--------|-------------------|-----------------|
| >=1024px (Desktop) | 5 equal vertical columns, full viewport height | Vertical columns | Thick borders between columns. All controls visible. |
| 768-1023px (Tablet) | 5 equal vertical columns, reduced height | Vertical columns | Same as desktop. No condensation — brutalism doesn't hide things. |
| <768px (Mobile) | 5 stacked horizontal bars, equal height | Horizontal bars | Thick borders between rows. All controls visible. |

---

## Micro-Interactions Detail (Brutalist)

### Palette Generation
- **Duration**: 0ms — **instant**. No transition, no stagger. Colors snap to new values.
- **Why**: Brutalism rejects decorative animation. The instant change is more honest and actually improves the rapid-fire generation loop (no 400ms wave to sit through).

### Lock Toggle
- **Duration**: 0ms — instant state change
- **Visual**: Text swaps from `[LOCK]` to `[LOCKED]`. Cross-hatch overlay appears immediately. No bounce, no morph, no scale animation.
- **If keeping SVG icons**: Icon swap is instant. No animation.

### Copy Confirmation Toast
- **Enter**: Instant (`display: block`, `opacity: 1`). No slide, no fade.
- **Display**: 1.5 seconds
- **Exit**: Instant (`display: none`). No fade-out.
- **Stacking**: Replace previous toast instantly.

### Hover on Swatch (Desktop)
- **Duration**: 0ms — instant
- **Effect**: Thick dashed inner border (3px, semi-transparent white or black depending on swatch luminance). No brightness overlay. No opacity fade on action buttons (they're always visible).

### Button Press (Generate, Lock, Copy, Format)
- **Effect**: Color inversion — background and text swap colors for the duration of `:active`. Releases immediately on mouseup. No easing.

---

## Design Tokens (Brutalist)

These replace the "Warm Dark Craft" tokens. The UI Designer's spec will define exact CSS values.

| Token | Value | Change from Previous |
|-------|-------|---------------------|
| `--bg-app` | `#000000` | Was `#1A1614` (warm dark). Now pure black. |
| `--color-surface` | `#000000` | Was `#242019`. Same as background — no elevation layers. |
| `--color-surface-raised` | `#000000` | Was `#2E2A22`. Flat. Borders define structure, not surface color. |
| `--text-primary` | `#FFFFFF` | Was `#E8E0D6` (warm white). Now pure white. |
| `--text-secondary` | `#AAAAAA` | Was `#9C9083` (warm gray). Now neutral gray. |
| `--text-on-swatch` | Auto `#000000` or `#FFFFFF` | Same logic, same values — no warm tinting. |
| `--accent` | `#FFFFFF` | Was `#D4913B` (amber). Now white — or optionally raw yellow `#FFFF00` for key CTAs. |
| `--color-border` | `#FFFFFF` | Was `rgba(255,255,255,0.08)`. Now solid white, visible, structural. |
| `--border-width` | `3px` | Was `1px`. Thick, visible borders define the grid. |
| `--radius-*` | `0px` | Was `6-16px`. All zero. Sharp corners everywhere. |
| `--font-ui` | `'Space Mono', 'JetBrains Mono', 'Courier New', monospace` | Was `Geist Sans`. Monospace for everything. |
| `--font-mono` | Same as `--font-ui` | No distinction between UI font and code font. Everything is mono. |
| `--transition-fast` | `0ms` | Was `150ms ease`. Instant. |
| `--transition-normal` | `0ms` | Was `200ms ease-out`. Instant. |
| `--transition-color` | `0ms` | Was `250ms ease-out`. Instant. |
| `--shadow-*` | `none` | Was various warm glows/elevations. No shadows in brutalism. |
| `--space-swatch-gap` | `0px` | Same as before — swatches edge-to-edge, separated by thick borders. |
| `--touch-target-min` | `44px` | Same — accessibility unchanged. |

---

## Export Formats

Same as previous spec — no changes. The export dropdown may use native `<select>` or raw text buttons rather than a styled custom dropdown.

| Format | Output Example | Copy Text |
|--------|---------------|-----------|
| CSS Variables | `--color-1: #3A86FF;` (one per line) | Multi-line string |
| JSON | `["#3A86FF","#FF006E","#FB5607","#FFBE0B","#8338EC"]` | Single-line array |
| Array | `#3A86FF, #FF006E, #FB5607, #FFBE0B, #8338EC` | Comma-separated |
| URL | `?colors=3A86FF,FF006E,FB5607,FFBE0B,8338EC` | Shareable link |

**Security note** (unchanged): All exported values validated hex strings. URL sharing sanitized with `/^[0-9A-Fa-f]{3,6}$/`. Use `textContent` not `innerHTML`.

---

## Keyboard Shortcut Reference

Same as previous spec — no changes.

| Shortcut | Action | Scope |
|----------|--------|-------|
| `Space` | Generate new palette | Global (when no input focused) |
| `L` | Lock/unlock focused/hovered color | Swatch context |
| `C` | Copy focused/hovered color | Swatch context |
| `1-5` | Focus swatch by number | Global |
| `Tab` | Move focus to next element | Standard |
| `Escape` | Close any open dropdown | Global |

---

## Summary of What Changes vs. What Stays

### STAYS THE SAME (do not touch)
- All JavaScript logic (color generation, harmony algorithms, clipboard, keyboard shortcuts)
- HTML structure and semantic elements
- ARIA labels and screen reader behavior
- User flows and interaction patterns
- Feature set (generate, lock, copy, format toggle, harmony modes, export)
- Responsive breakpoints and layout grid logic
- Security constraints (textContent, CSP, hex validation)

### CHANGES (CSS/visual only)
- All colors → black/white/optional raw yellow
- All border-radius → 0px
- All borders → 3px solid, visible
- All shadows → removed
- All transitions → 0ms (instant)
- All fonts → monospace only
- All text → uppercase for labels/headings
- Icons → text labels preferred (`[LOCK]`, `[COPY]`, etc.)
- Hover states → instant, no fade (dashed border indicator)
- Toast → instant appear/disappear, no slide animation
- Button active states → color inversion, no lift/glow

---

Status: READY_FOR_REVIEW
