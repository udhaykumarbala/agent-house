# Template Selection

## Chosen Template: `static-enhanced`

## Task Analysis
Build a premium, minimal unit converter (length, weight, temperature) with modern typography, clean micro-interactions, and a refined visual identity. This is a purely client-side utility — all conversion logic is mathematical with no backend or database needed. The emphasis is on polished UI/UX that feels like a premium app.

## Justification
- **No backend required** — All unit conversions (length, weight, temperature) are pure math operations that run entirely in the browser, ruling out `go-api` and `fullstack`
- **Not complex enough for Next.js** — A single-page converter tool doesn't need SSR/SSG, routing, or React's component model; `nextjs-frontend` would be over-engineered
- **Build tools earn their keep** — Tailwind CSS enables rapid iteration on the "premium" visual identity with utility classes, consistent spacing/typography scales, and easy dark mode; Vite provides fast HMR during development and optimized production builds
- **TypeScript adds safety for conversion logic** — Unit conversion formulas benefit from type checking to prevent bugs (e.g., mixing up conversion directions), which `static-html` wouldn't provide
- **Modern CSS animations via Tailwind** — The "clean micro-interactions" requirement is well-served by Tailwind's transition/animation utilities combined with custom CSS, without needing a JS animation library

## Customizations Needed
- Add CSS custom properties for the premium color palette and typography scale
- Include smooth transition utilities for micro-interactions (input focus, unit switching, value changes)
- Configure Tailwind with a custom theme for the refined visual identity (custom fonts, extended color palette)
- No additional libraries needed — keep it minimal and fast

---
Status: AWAITING_REVIEW
