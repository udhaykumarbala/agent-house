# APPROVED PLAN: Image Resizer (Friendly Brutalist)

Reviewed by: CEO
Date: 2026-03-16
Status: APPROVED FOR DEVELOPMENT

## Summary

A single-page, 100% client-side image resizer with a friendly brutalist (neo-brut) design. Users drop an image, set dimensions, and download the resized file — all within the browser. No server uploads, no accounts, no ads. The warm "Butter & Charcoal" aesthetic differentiates it from the sea of generic utility tools.

Target user: Content creators, social media managers, bloggers who need quick image resizing without installing software or uploading to unknown servers.

## Approved Specifications

- Product Spec: .plans/specs/product-spec.md
- UX Spec: .plans/specs/ux-spec.md
- UI Spec: .plans/specs/ui-spec.md
- Architecture Spec: .plans/specs/architecture-spec.md
- Security Spec: .plans/specs/security-spec.md

## Key Decisions

1. **Color palette**: UI spec's "Butter & Charcoal" palette is canonical (#FEF3E2 background, #E8625C coral primary, #F5D547 butter yellow secondary, #7A9E7E sage accent, #1A1A1A charcoal borders). All other specs must align to these colors.

2. **Desktop layout**: Two-column layout (60% preview + 40% controls sidebar) as defined in UX/UI specs. NOT single-column stacked as architect spec proposed.

3. **Supported file types**: JPG, PNG, WebP only (per product spec). No GIF/BMP despite security spec including them.

4. **Responsive breakpoints**: Mobile < 640px, Tablet 640-1024px, Desktop > 1024px (per UI spec).

5. **Typography**: Space Grotesk for headings, Inter for body/buttons, JetBrains Mono for dimension inputs and file sizes. All via Google Fonts.

6. **Icons**: Use inline SVG for the ~10 needed icons rather than Phosphor Icons CDN. Keeps CSP strict and dependencies at zero JS libraries.

7. **Download bar**: `position: fixed; bottom: 0` (always visible, not sticky).

8. **Button font**: Inter (per UI spec), not Space Grotesk.

9. **Template**: `static-html` — vanilla HTML/CSS/JS, no build tools, no framework.

10. **Security**: All DOM rendering via textContent only, file validation with MIME + magic bytes, CSP headers configured at deployment, blob URL lifecycle management.

## Implementation Order

1. **HTML structure + CSS design system** — Landing state (drop zone) and workspace state skeleton. Implement the full brutalist design system (colors, typography, borders, shadows, buttons, inputs). Mobile responsive from the start.

2. **Core resize engine** — `resizer.js` module: Canvas API resize, format conversion (JPG/PNG/WebP), quality control, multi-step downscaling for quality, file size estimation.

3. **UI logic and interactivity** — `app.js`: Drag-and-drop upload, file validation, state management (upload/editor views), dimension inputs with aspect ratio lock, preset buttons, format/quality controls, live preview, download with sanitized filename.

4. **Polish and edge cases** — Error states, loading states, success feedback, accessibility (ARIA labels, keyboard navigation, focus management, prefers-reduced-motion), privacy messaging, blob URL cleanup.

## Notes for Development Team

- **UI spec is the design authority.** When in doubt about any visual decision (colors, spacing, shadows, typography), refer to ui-spec.md CSS custom properties section.
- **Security spec is non-negotiable.** Zero innerHTML with user data. All file validation must happen before processing. Sanitize download filenames.
- **No inline scripts or event handlers in HTML.** All JS in external files, all events via addEventListener. This supports strict CSP.
- **Inline SVG icons** instead of Phosphor CDN. Copy only the needed icon paths (~10 icons) directly.
- **Test on mobile early.** Many users will resize images from their phones. Touch targets >= 44x44px.
- **Privacy message must be visible** on the landing state — "Your images never leave your browser."
- **Active preset buttons use butter yellow** (#F5D547) per UI spec, not coral.
- **Quality slider hides for PNG** (lossless format, no quality setting).
- **EXIF stripping is a feature** — Canvas naturally strips metadata. Mention this as a privacy benefit.
- **WebP output detection**: Check browser support, fallback to JPEG if WebP encoding not available.

---
BEGIN DEVELOPMENT
