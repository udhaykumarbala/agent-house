# Template Selection

## Chosen Template: `static-html`

## Task Analysis
This is a feature enhancement to an existing, fully built client-side image resizer. The project already uses `static-html` (pure HTML/CSS/JS with no build tools). All proposed new features — cropping, rotation, batch processing, filters, clipboard paste, dark mode — are achievable entirely client-side using the Canvas API, File API, and localStorage. No backend, no build tooling change needed.

## Justification
- **Existing project uses static-html** — Changing templates would require a full rewrite of working, QA-approved code. Continuity is the right call.
- **All new features are client-side** — Cropping, rotation, filters, and batch processing all use the Canvas API already in place. No server needed.
- **No new dependencies required** — Clipboard API, localStorage, and CSS custom properties for dark mode are all native browser APIs. Zero build complexity stays zero.
- **Brutalist design system is hand-crafted CSS** — The existing CSS variables and component styles extend naturally. Tailwind/Vite would add friction, not value.
- **Proven architecture** — The modular JS structure (utils.js, resizer.js, app.js) can be extended with new modules (crop.js, filters.js, batch.js) following the same pattern.

## Customizations Needed
- Add new JS modules: `crop.js`, `filters.js`, `batch.js`
- Extend `resizer.js` with rotation/flip transforms
- Add dark mode CSS via `prefers-color-scheme` + manual toggle
- Extend HTML with new UI sections (crop overlay, filter controls, batch list)

---
Status: AWAITING_REVIEW
