# APPROVED PLAN: Color Palette Generator — Brutalist Reskin

Reviewed by: CEO
Date: 2026-03-16
Status: APPROVED FOR DEVELOPMENT

## Summary
Transform the existing color palette generator's visual design system from "Warm Dark Craft" (warm obsidian + amber accents, rounded corners, soft shadows) to a brutalist design language. Bold monospace typography, raw/exposed UI elements, high contrast, thick borders, zero decoration, and instant state changes. All functionality remains identical — this is a CSS-first reskin with minimal HTML and JS touch-ups.

## Approved Specifications
- Product Spec: .plans/specs/product-spec.md
- UX Spec: .plans/specs/ux-spec.md
- UI Spec: .plans/specs/ui-spec.md
- Architecture Spec: .plans/specs/architecture-spec.md
- Security Spec: .plans/specs/security-spec.md

## Key Decisions (Conflict Resolutions)

1. **Dark Background (#0A0A0A)**: The UI and UX specs agree on dark brutalism. A dark frame makes generated palette colors pop harder. The architecture spec's white-background proposal is overruled. All borders, text, and chrome use white (#FFFFFF) on dark background.

2. **Electric Yellow Accent (#EBFF00)**: One harsh industrial accent color for the Generate button (primary CTA) and focus rings only. Not used for general decoration or secondary elements. This balances the PM's "no accent" preference with the UI designer's need for a clear call-to-action.

3. **Text Labels, No Icons**: `[LOCK]`, `[COPY]`, `[GENERATE]` — text-based buttons, no Lucide icon library. Reduces dependencies, reinforces the raw brutalist aesthetic. The UI spec's Lucide icons proposal is overruled in favor of the product and UX specs.

4. **Actions Always Visible**: All swatch controls (lock, copy, hex code) visible at all times on both desktop and mobile. No hover-to-reveal pattern. Brutalism exposes everything. The UI spec's "appears on hover" is overruled in favor of the UX spec.

5. **Space Mono from Google Fonts CDN**: The visual distinctiveness of Space Mono outweighs the minor CDN dependency concern. CSP must be updated to include `fonts.googleapis.com` (style-src) and `fonts.gstatic.com` (font-src). Fallback: `'Courier New', Courier, monospace`.

6. **Styled Select for Harmony Mode**: Custom-styled thick-bordered monospace select (per UI spec), not raw native browser `<select>`. Maintains visual consistency while staying minimal.

7. **No Shadows of Any Kind**: Zero `box-shadow`, zero hard-offset shadows, zero `filter: drop-shadow`. Borders are the only depth cue. The architecture spec's hard-offset shadow proposal is overruled.

8. **Lock Indicator = Thick Top Border**: 3px solid border in contrasting text color (black on light swatches, white on dark swatches). No cross-hatch patterns, no diagonal stripes, no overlays.

9. **Mobile Breakpoint at 768px**: `<768px` = stacked horizontal bars (mobile), `>=768px` = vertical columns (desktop). Matches existing implementation.

## Implementation Order

1. **Phase 1 — Design Token Overhaul**: Replace all CSS custom properties in `:root` — colors, typography, borders, radius, shadows, transitions. Update `index.html` font link to Space Mono. Update CSP meta tag.
2. **Phase 2 — Component Reskin**: Restyle toolbar, generate button, format toggle, harmony select, toast, keyboard hints, swatch action buttons (text labels replacing icons).
3. **Phase 3 — Swatch & Layout Polish**: Swatch dividers (1px solid #333333), lock indicator (thick top border), always-visible controls, responsive layout verification, focus states (3px yellow outline).

## Notes for Development Team

- **CSS-first change**: The reskin is primarily `css/styles.css`. Minimal `index.html` edits (font link, CSP). One `js/app.js` edit (replace hardcoded amber `#D4913B` lock color with dynamic contrast color).
- **Design tokens to override**: Background `#0A0A0A`, Surface `#141414`, Surface Raised `#1E1E1E`, Text `#FFFFFF`, Border `#FFFFFF`, Border Muted `#333333`, Accent `#EBFF00`. All radius = 0px. All shadows = none. All transitions = none (instant).
- **Typography**: Space Mono 700 for headings/buttons/codes (uppercase, letter-spaced), Space Mono 400 for body. Everything monospace.
- **Text labels replace icons**: `[LOCK]`/`[LOCKED]`, `[COPY]`, `[GENERATE]`, `[COPY ALL]`. Remove Lucide icon references if present.
- **Swatch controls always visible**: No opacity transitions on hover — controls are visible at all times on all breakpoints.
- **Instant state changes**: `transition: none` everywhere. Colors snap, buttons invert instantly, toast appears/disappears without animation.
- **Security (non-negotiable)**: Verify no `innerHTML` introduced. Update CSP font-src/style-src for Space Mono. All other security controls (R1-R5) unchanged.
- **Accessibility**: Focus rings = 3px solid #EBFF00, offset 3px. Touch targets = 44px minimum. High contrast black/white exceeds WCAG AAA.

---
BEGIN DEVELOPMENT
