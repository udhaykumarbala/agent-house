# Template Selection

## Chosen Template: `static-enhanced`

## Task Analysis
A browser-based meme creator tool requiring: image upload, HTML5 Canvas-based editor with draggable text overlays, text styling controls (font, size, color, stroke), and client-side image export (canvas.toBlob/toDataURL). This is a single-page application with zero backend requirements — all processing happens in the browser.

## Justification
- **No backend needed** — Image manipulation, text rendering, and export all use browser-native Canvas API. Eliminates `go-api` and `fullstack`.
- **Single page, not a complex SPA** — No routing, state management libraries, or SSR/SSG required. `nextjs-frontend` is overkill for one page with a canvas.
- **Tailwind CSS accelerates UI development** — The toolbar, controls panel, font pickers, color inputs, and layout chrome benefit significantly from utility-first CSS. This pushes us from `static-html` to `static-enhanced`.
- **Vite provides fast dev iteration** — Hot module replacement speeds up the tight feedback loop needed when building interactive canvas tooling.
- **Right complexity level** — The core logic (Canvas API, drag handling, text rendering) is vanilla JS. Wrapping it in Vite/Tailwind adds minimal overhead while improving DX and UI quality.

## Why Not Other Templates
- `static-html`: Viable but would require hand-writing all CSS for a polished UI. Tailwind is a significant productivity win for the controls/toolbar.
- `nextjs-frontend`: No routing, no SSR, no complex state — unnecessary framework weight for a single canvas page.
- `go-api` / `fullstack`: No server-side processing needed whatsoever.

## Customizations Needed
- Add HTML5 Canvas utility patterns (coordinate transforms, hit testing for drag)
- Include a file input / drag-and-drop zone for image upload
- Use TypeScript for safer canvas coordinate math and event handling
- Consider adding a simple undo/redo stack (array-based, no external lib)

---
Status: AWAITING_REVIEW
