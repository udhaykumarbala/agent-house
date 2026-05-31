# Product Specification: Color Palette Generator — Brutalist Redesign

Based on research in: .plans/research/product-research.md, .plans/research/ux-research.md, .plans/research/ui-research.md, .plans/research/architecture-research.md, .plans/research/security-expert.md

## Context: Design System Pivot

**CEO Directive**: Transform the visual design system from "Warm Dark Craft" (amber accents, rounded corners, warm obsidian, soft shadows) to a **brutalist design language**. All functionality remains identical — this is a skin change, not a feature change.

**Why Brutalism**: The palette generator market is saturated with polished, safe, corporate-looking tools (Coolors' clean blue, Adobe's charcoal chrome, ColorHunt's Pinterest-white). Brutalism is a deliberate counter-positioning move: raw, honest, memorable. It signals that this tool is built for creators who value function over decoration — the same audience that gravitates toward tools like Craigslist, Hacker News, or early StackOverflow, which succeeded precisely because they refused to "look nice."

## Target User

**"Quick-Pick Quinn"** — A 22-38 year old freelance web/graphic designer or front-end developer who needs fresh color inspiration fast during design phases. High tech comfort, lives in browser dev tools, works on multiple small-to-medium projects. Uses the tool in short, bursty sessions (under 2 minutes): open, generate 3-10 palettes, copy the winner, close.

**Brutalism fit for this user**: Designers and developers who use palette generators are visually literate. They recognize and appreciate brutalist design as an intentional aesthetic choice, not a lack of polish. The raw, high-contrast interface will feel refreshing against the sea of smooth-cornered, gradient-laden design tools — and the monospace typography naturally fits a tool built around color codes.

## Unique Value Proposition

"Instant, beautiful color palettes — no sign-up, no ads, one click to copy."

Unchanged. Brutalism amplifies this message: the UI is stripped to its essentials, making the value proposition self-evident. No decoration competing for attention — just colors and actions.

## Brutalist Design Principles (Replacing "Warm Dark Craft")

These principles replace the previous five ("Colors are the hero," "Warm, not cold," etc.):

1. **Raw structure, exposed grid** — No decorative elements. The layout grid IS the design. Thick borders define space instead of shadows or gradients.
2. **Bold monospace typography** — All type in monospace. Headings are oversized and heavy. Hex codes feel native, not styled.
3. **High contrast, binary palette** — Black and white as the primary system colors. No warm neutrals, no subtle tints. The generated palette colors provide ALL the chromatic interest.
4. **Thick borders, no shadows** — 2-4px solid black borders replace all shadows, glows, and subtle border effects. Hard edges only.
5. **Intentionally unpolished** — No rounded corners (0px radius everywhere). No smooth transitions on non-essential elements. Hover states are immediate, not eased. The tool looks like it was built to work, not to impress.

## Features (Prioritized)

### P0 - Must Have (MVP)

All features are functionally identical to the original spec. Only the visual treatment changes.

- [ ] **Random Palette Generation**: Generate a 5-color palette on page load and on demand via a "Generate" button — HSL with constrained saturation (40-90%) and lightness (30-70%)
- [ ] **Spacebar Shortcut**: Pressing spacebar generates a new palette instantly
- [ ] **Full-Viewport Swatch Layout**: 5 equal-width vertical columns spanning full viewport height on desktop; stacks to horizontal bars on mobile. **Brutalist treatment**: thick black borders (3px solid) between swatches instead of seamless edges
- [ ] **Hex Code Display with Auto-Contrast Text**: Each swatch displays its hex code in uppercase monospace, centered. Auto light/dark text based on luminance. **Brutalist treatment**: oversized monospace type (24px+), no letter-spacing refinement
- [ ] **One-Click Copy to Clipboard**: Clicking a swatch's hex code copies it via Clipboard API with a blunt "COPIED" confirmation — no toast slide-in, just a raw text flash that disappears
- [ ] **Lock/Unlock Colors**: Click a lock icon on any swatch to preserve that color during regeneration. **Brutalist treatment**: lock state shown with a bold "LOCKED" text stamp or heavy border instead of a subtle icon state change
- [ ] **Responsive Design**: Desktop: horizontal columns. Mobile: vertical stacked bars. Touch targets minimum 44x44px. **Brutalist treatment**: same thick-border grid system at all breakpoints

### P1 - Should Have

- [ ] **Color Harmony Modes**: Random, Analogous, Complementary, Split-Complementary, Triadic, Monochromatic. **Brutalist treatment**: plain text list or raw button group, no dropdown — all options visible at once
- [ ] **Multiple Color Format Display**: Toggle between HEX (default), RGB, HSL. **Brutalist treatment**: simple text buttons labeled `HEX | RGB | HSL` with an underline or background-invert on active state
- [ ] **Keyboard Shortcuts for Power Users**: Spacebar = generate, L = lock, C = copy. **Brutalist treatment**: shortcuts displayed prominently in the UI as raw `[SPACE]` `[L]` `[C]` labels, not hidden in tooltips
- [ ] **Copy All Palette**: Button to copy entire palette as comma-separated hex codes

### P2 - Nice to Have

- [ ] **Export as CSS Variables**: Copy full palette as CSS custom properties
- [ ] **URL-Based Sharing**: Encode palette in URL hash (strict hex validation)
- [ ] **Palette History via localStorage**: Store recent palettes in browser
- [ ] **Shade/Tint Variations**: Lighter/darker variants of each color

## User Stories

1. As a freelance designer, I want to press spacebar to rapidly generate new palettes so that I can find color inspiration without any friction
2. As a front-end developer, I want to click a hex code and have it copied to my clipboard so that I can immediately paste it into my CSS
3. As a designer refining a palette, I want to lock colors I like and regenerate the rest so that I can converge on a perfect combination without starting from scratch
4. As a developer working across formats, I want to toggle between HEX, RGB, and HSL display so that I can copy the format my project needs
5. As a mobile user browsing for inspiration, I want the tool to work on my phone with touch-friendly targets so that I can explore palettes on any device
6. As a user who values aesthetics, I want to select a color harmony mode so that generated palettes are grounded in color theory
7. As a developer integrating colors, I want to export the palette as CSS variables so that I can drop them directly into my stylesheet
8. As a visually literate designer, I want the tool's interface to have a distinctive brutalist aesthetic so that the UI feels intentional and memorable rather than generic

## Acceptance Criteria

### Feature: Brutalist Design System (NEW — Core of This Spec)
- [ ] All border-radius values are 0px — no rounded corners anywhere
- [ ] All borders are solid black, minimum 2px width (3-4px for primary containers)
- [ ] No box-shadows, no glows, no gradients anywhere in the UI
- [ ] All typography uses a monospace font stack (system monospace fallbacks: `'Courier New', 'Consolas', 'Liberation Mono', monospace`)
- [ ] Primary system colors are #000000 (black) and #FFFFFF (white) — the generated palette provides all chromatic color
- [ ] Headings are bold/black weight, oversized relative to conventional sizing (e.g., page title at 48px+)
- [ ] Interactive elements use immediate state changes (no transition easing for hover/active states)
- [ ] The generate button uses a high-contrast inverted style: black background, white text (or vice versa), thick border
- [ ] The overall visual impression is raw, structural, and intentionally "unpolished"
- [ ] The design is visually distinctive from Coolors, Adobe Color, and all competitors analyzed in research

### Feature: Random Palette Generation
- [ ] Page loads with a randomly generated 5-color palette
- [ ] Clicking the "Generate" button produces a new palette
- [ ] All generated colors have constrained saturation (40-90%) and lightness (30-70%)
- [ ] Generation is instant — no perceptible delay
- [ ] Each of the 5 swatches fills an equal portion of the viewport

### Feature: Spacebar Shortcut
- [ ] Pressing spacebar anywhere on the page generates a new palette
- [ ] Spacebar does NOT trigger generation when focused on an input/text field
- [ ] Works across Chrome, Firefox, Safari, Edge

### Feature: One-Click Copy to Clipboard
- [ ] Clicking a hex code copies the value to the system clipboard
- [ ] A "COPIED" confirmation appears in-place (not a sliding toast) and disappears after ~1.5 seconds
- [ ] The copied value is the clean hex string from the JS data model (not DOM innerHTML)
- [ ] Works on desktop (click) and mobile (tap)

### Feature: Lock/Unlock Colors
- [ ] Each swatch has a clearly visible lock control
- [ ] Locked state is visually prominent — bold "LOCKED" label, thick border change, or strong visual indicator (not a subtle icon swap)
- [ ] Locked colors are preserved during regeneration; only unlocked colors change
- [ ] Lock state persists across multiple generations within the same session
- [ ] Locking all 5 colors and pressing generate results in no change (no error)

### Feature: Responsive Design
- [ ] Desktop (>=768px): 5 equal-width vertical columns spanning full viewport height
- [ ] Mobile (<768px): 5 stacked horizontal bars spanning full width
- [ ] All interactive elements have minimum 44x44px touch targets
- [ ] Thick border grid system is maintained at all breakpoints
- [ ] Hex codes remain readable on all screen sizes

### Feature: Color Harmony Modes (P1)
- [ ] Mode selector shows all options visibly (not hidden in a dropdown) — brutalist preference for visible, raw controls
- [ ] Options: Random, Analogous, Complementary, Split-Complementary, Triadic, Monochromatic
- [ ] Each mode uses correct color theory angles
- [ ] Active mode indicated by inverted colors (black bg, white text) or heavy underline

### Feature: Auto-Contrast Text
- [ ] Text on light swatches uses black (#000000)
- [ ] Text on dark swatches uses white (#FFFFFF)
- [ ] Contrast computed using luminance formula: `0.299*R + 0.587*G + 0.114*B`
- [ ] No warm-white or soft-black variants — pure black and white only (consistent with brutalist system)

## Brutalist Visual System Summary

| Property | Old ("Warm Dark Craft") | New (Brutalist) |
|----------|------------------------|-----------------|
| Background | #1A1614 (warm obsidian) | #FFFFFF (white) or #000000 (black) |
| Accent color | #D4913B (amber) | None — generated palette colors only |
| Text color | #E8E0D6 (warm white) | #000000 (pure black) on white bg |
| Border style | 1px rgba(255,255,255,0.08) | 3px solid #000000 |
| Border radius | 6-12px | 0px everywhere |
| Shadows | Warm-tinted glows, subtle | None |
| Font family | Geist Sans + Geist Mono | Monospace only (system monospace stack) |
| Font weight | Mixed 400-700 | Heavy — 700-900 for headings, 400-500 for body |
| Transitions | 150-300ms ease | Instant (0ms) for state changes; optional minimal transition for palette generation only |
| Button style | Amber fill, rounded, subtle glow | Black/white inverted, square, thick border |
| Toast style | Floating card with shadow | Raw inline text replacement |
| Icons | Lucide line icons, warm tints | Text labels or raw Unicode glyphs instead of icon library |
| Overall feel | "Creative studio at night" | "Printed function sheet" / "Swiss poster" |

## Success Metrics

- **Time to first palette**: < 1 second from page load (no framework overhead)
- **Copy action rate**: > 60% of sessions include at least one copy action
- **Session depth**: Average 5+ generations per session
- **Bounce rate**: < 30%
- **Mobile usability**: Responsive layout renders correctly 320px-2560px
- **Design distinctiveness**: Tool is visually identifiable within 2 seconds as NOT a Coolors/Adobe clone (qualitative — team review)

## Out of Scope (for MVP)

- User accounts, sign-up, or authentication
- Server-side backend or API
- Image-based palette extraction
- Community/social features
- Drag-to-reorder swatches
- Undo/redo palette history
- Accessibility contrast checker between palette color pairs
- Download as PNG/SVG/PDF
- Dark/light mode toggle (ship with single brutalist theme — white background variant)
- Any monetization, ads, or premium tier
- The "Warm Dark Craft" design system — this is fully replaced by brutalist

## Technical Constraints (from Architecture & Security Research)

- **Template**: `static-html` — vanilla HTML/CSS/JS, no framework, no build step
- **Internal color representation**: HSL (convert to HEX/RGB for display)
- **Zero external JS dependencies**
- **Security**: `textContent` not `innerHTML`; validate URL params with hex regex; CSP meta tag; copy from JS data model
- **Typography**: System monospace font stack — no external font loading needed (reduces page weight, aligns with brutalist raw ethos)
- **No icon library**: Use text labels (`LOCK`, `COPY`, `GENERATE`) or Unicode characters instead of Lucide/Phosphor — reduces dependencies and reinforces brutalist aesthetic

---
Status: READY_FOR_REVIEW
