# Architecture Research: Random Color Palette Generator

## 1. Competitor Analysis

### Coolors (coolors.co)
- **How it works**: Spacebar to generate 5-color palettes; lock colors to keep them while regenerating others
- **Key features**: Lock/unlock colors, adjust individual colors (hue/saturation/brightness sliders), copy hex/RGB/HSL, export as PNG/SVG/PDF, share via URL, accessibility contrast checker
- **UX pattern**: Horizontal 5-swatch layout fills the viewport; each swatch shows hex code centered; hover reveals action icons (lock, drag, copy, adjust, delete)
- **Keyboard shortcut**: Spacebar = generate — extremely intuitive, low-friction interaction
- **Takeaway**: The spacebar-to-generate pattern is the gold standard; lock/unlock is the #1 power-user feature

### Colormind (colormind.io)
- **How it works**: AI/ML-based palette generation trained on photographs, movies, and art
- **Key features**: Generate palette, lock colors, integrations for Bootstrap/Material, website preview
- **UX pattern**: 5 horizontal swatches, minimal UI, "Generate" button
- **Takeaway**: ML-backed generation produces aesthetically pleasing results but requires a backend; for our client-side tool, algorithmic approaches (color harmony theory) are the way to go

### Paletton (paletton.com)
- **How it works**: Color wheel-based; pick a base color and a harmony type (complementary, triadic, tetradic, analogous)
- **Key features**: Color wheel visualization, harmony presets, fine-tuning with distance/angle sliders
- **Takeaway**: Color harmony theory produces reliably pleasing palettes; we should offer harmony modes, not just pure random

### Color Hunt (colorhunt.co)
- **How it works**: Curated 4-color palettes submitted by users; browseable/searchable
- **Key features**: Browse by tag (pastel, vintage, neon), like/save, copy palette
- **Takeaway**: 4-color palettes are the sweet spot for quick use; tags/moods are a nice-to-have

### Adobe Color (color.adobe.com)
- **How it works**: Full-featured color wheel with harmony rules, extract from images, explore trending palettes
- **Key features**: Harmony rules (analogous, monochromatic, triad, complementary, split-complementary, double-split, square), accessibility tools, save to Creative Cloud
- **Takeaway**: The harmony rule set from Adobe is comprehensive and well-understood

---

## 2. Color Generation Algorithms

### 2.1 Pure Random (HSL-based)
Generate random H (0-360), S (40-90%), L (30-70%) values. Constraining S and L avoids muddy or washed-out colors.

```
h = Math.random() * 360
s = 40 + Math.random() * 50  // 40-90%
l = 30 + Math.random() * 40  // 30-70%
```
**Pros**: Simple, fast, unpredictable
**Cons**: No guarantee of harmony; adjacent colors may clash

### 2.2 Color Harmony Rules
Start from a random base hue, then derive the rest using angular relationships on the color wheel:

| Harmony | Angles from base | Colors |
|---------|-----------------|--------|
| Analogous | 0°, ±30° | 3-5 |
| Complementary | 0°, 180° | 2 (+variations) |
| Split-complementary | 0°, 150°, 210° | 3 |
| Triadic | 0°, 120°, 240° | 3 |
| Tetradic (square) | 0°, 90°, 180°, 270° | 4 |
| Monochromatic | Same hue, varying S/L | 3-5 |

**Recommended**: Default to random harmony selection for "surprise me" mode; let user pick specific harmony types.

### 2.3 Palette Size
- Industry standard: **5 colors** (Coolors, Colormind)
- Minimum useful: 3 colors
- Maximum practical: 8 colors
- **Recommendation**: Default 5, allow 3-8 range

---

## 3. Color Format Conversions

Must support these formats (all can be computed client-side):

| Format | Example | Use Case |
|--------|---------|----------|
| HEX | `#3A86FF` | Web CSS, most common |
| RGB | `rgb(58, 134, 255)` | CSS, design tools |
| HSL | `hsl(220, 100%, 61%)` | Intuitive for adjustments |

### Conversion Functions Needed
- HSL → RGB → HEX (generation path)
- HEX → RGB (display)
- RGB → HSL (for adjustment sliders)

All conversions are well-documented math — no library needed. Standard algorithms from CSS Color Module Level 4 spec.

---

## 4. Essential Features (MVP)

Based on competitor analysis, these are table-stakes:

1. **Generate random palette** — button click + spacebar shortcut
2. **Display 5 color swatches** — full-height vertical bars or large horizontal blocks
3. **Show color codes** — HEX displayed on each swatch
4. **Copy to clipboard** — click hex code or copy button to copy
5. **Lock/unlock colors** — click lock icon to preserve a color during regeneration
6. **Responsive layout** — works on mobile (stack vertically) and desktop (horizontal)

## 5. Nice-to-Have Features (Post-MVP)

1. **Harmony mode selector** — analogous, complementary, triadic, etc.
2. **Export options** — download as PNG image, CSS variables, JSON
3. **Color format toggle** — switch between HEX, RGB, HSL display
4. **Shade/tint variations** — show lighter/darker variants of each color
5. **Accessibility contrast check** — WCAG AA/AAA contrast ratio between pairs
6. **URL sharing** — encode palette in URL hash for sharing
7. **Undo/redo** — palette history navigation
8. **Drag to reorder** — rearrange colors in the palette

---

## 6. UX Patterns & Best Practices

### Layout
- **Desktop**: 5 equal-width vertical columns spanning full viewport height (Coolors pattern)
- **Mobile**: Stack as horizontal bars, full-width, equal-height
- **Each swatch**: Color fills the entire area; text overlays (hex code + action icons) centered or bottom-aligned

### Contrast for Text on Swatches
Text displayed over color swatches must be readable. Use luminance calculation to decide white vs. black text:

```
luminance = 0.299*R + 0.587*G + 0.114*B
textColor = luminance > 150 ? '#000000' : '#FFFFFF'
```

### Interactions
- **Spacebar**: Generate new palette (global keyboard listener)
- **Click on hex code**: Copy to clipboard with brief "Copied!" feedback
- **Click lock icon**: Toggle lock state (locked colors persist on regeneration)
- **Hover on swatch**: Reveal action buttons (copy, lock, adjust)

### Visual Feedback
- Toast/snackbar on copy: "Copied #3A86FF!" — auto-dismiss after 1.5s
- Lock icon state change: Filled lock when locked, outline when unlocked
- Smooth color transitions on regeneration (CSS transition on background-color)

---

## 7. Technical Architecture

### Stack (per approved template)
- **HTML5**: Semantic structure, single page
- **CSS3**: Grid/Flexbox layout, transitions, responsive media queries
- **Vanilla JS**: Color math, DOM manipulation, Clipboard API, keyboard events

### File Structure
```
/
├── index.html          # Single HTML file with structure
├── css/
│   └── styles.css      # All styles
├── js/
│   ├── app.js          # Main app logic, event listeners
│   ├── colors.js       # Color generation & conversion functions
│   └── utils.js        # Clipboard, UI helpers
└── assets/
    └── favicon.svg     # Optional favicon
```

### Key Technical Decisions
1. **HSL as internal representation** — most intuitive for generation and manipulation; convert to HEX/RGB for display
2. **No external dependencies** — no libraries, no CDN (except possibly Google Fonts for typography)
3. **Clipboard API** — use `navigator.clipboard.writeText()` with fallback to `document.execCommand('copy')`
4. **CSS Grid for layout** — `grid-template-columns: repeat(5, 1fr)` for equal swatches
5. **CSS transitions** — `transition: background-color 0.3s ease` for smooth palette changes
6. **LocalStorage** — optional, for saving favorite palettes (nice-to-have)

### Browser Support
- All modern browsers (Chrome, Firefox, Safari, Edge)
- Clipboard API supported in all modern browsers
- CSS Grid supported everywhere relevant
- No polyfills needed

---

## 8. Accessibility Considerations

- **Keyboard navigation**: Tab through swatches, Enter to copy, L to lock
- **ARIA labels**: Descriptive labels on buttons ("Copy color #3A86FF", "Lock this color")
- **Focus indicators**: Visible focus ring on interactive elements
- **Screen reader**: Announce color values and actions
- **Color contrast**: Ensure UI chrome (buttons, text) meets WCAG AA against swatch backgrounds

---

## 9. Performance

- **Zero build step** — instant load
- **Minimal DOM** — ~20-30 elements total
- **No network requests** — fully offline-capable after initial load
- **Smooth animations** — CSS-only transitions, no JS animation loops
- **Estimated load time**: <100ms (no assets to fetch beyond HTML/CSS/JS)

---

## 10. Recommendations Summary

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Palette size | 5 colors (default) | Industry standard, enough variety |
| Generation method | HSL + harmony rules | Produces aesthetically pleasing results |
| Primary interaction | Spacebar to generate | Proven UX from Coolors |
| Color format | HEX (default), toggle to RGB/HSL | HEX is most used in web dev |
| Layout | Full-viewport vertical columns | Maximizes color visibility |
| Text contrast | Auto black/white based on luminance | Ensures readability |
| Dependencies | Zero | No build step, instant load, simple |

---

**Research complete.** Ready to proceed to architecture specification and development planning.
