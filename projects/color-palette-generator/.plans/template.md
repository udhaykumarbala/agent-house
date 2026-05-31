# Template Selection

## Chosen Template: `static-html`

## Task Analysis
Build a random color palette generator tool. This is a single-page, client-side interactive tool that generates random color palettes. All logic (color generation, format conversion, clipboard) runs in the browser with no server dependency. The UI is straightforward: color swatches, hex/RGB labels, and action buttons.

## Justification
- **No backend required** - Color generation is pure client-side math; no API, database, or server needed
- **Low complexity** - A single-page tool with a flat DOM structure; no routing, state management, or component hierarchy needed
- **No build tools needed** - Vanilla CSS handles the styling (gradients, transitions, grid layout); Tailwind/Vite would add unnecessary overhead
- **Fast iteration** - Pure HTML/CSS/JS means zero build step, instant reload, and easy debugging
- **Minimal deployment** - Can be served from any static host or opened directly as a file

## Customizations Needed
- None - using standard template. Vanilla JS is sufficient for color math, clipboard API, and DOM manipulation.

## PM Review

**TEMPLATE_APPROVED: static-html**

From a product perspective:
- Users expect instant interaction — no load time, no build overhead
- Market fit favors simplicity — competing tools (Coolors, Colormind) are lightweight browser tools
- No backend needed — all color math is client-side, no unnecessary complexity
- Broad accessibility — works on any device, any browser, even offline

---
Status: PM_APPROVED
