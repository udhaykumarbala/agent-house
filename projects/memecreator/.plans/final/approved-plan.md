# APPROVED PLAN: MemeForge — Browser-Based Meme Creator

Reviewed by: CEO
Date: 2026-03-16
Status: APPROVED FOR DEVELOPMENT

## Summary

A zero-friction, 100% client-side meme creator web app. Users upload images, add freely-positioned draggable text overlays, customize styling (font, color, size, stroke, shadow), and export polished memes as downloadable PNG/JPEG — all without sign-ups, watermarks, or server uploads. Built with Vite + Tailwind CSS + TypeScript using the raw HTML5 Canvas API.

Target user: "Quick-Meme Quinn" — 18-35, digitally native, wants a meme created in under 60 seconds.

## Approved Specifications

- Product Spec: `.plans/specs/product-spec.md`
- UX Spec: `.plans/specs/ux-spec.md`
- UI Spec: `.plans/specs/ui-spec.md`
- Architecture Spec: `.plans/specs/architecture-spec.md`
- Security Spec: `.plans/specs/security-spec.md`
- Development Plan: `.plans/development-plan.json`

## Key Decisions

### 1. Layout: 2-Column (Architecture Spec)
Use the **2-column layout** from the architecture spec: canvas area (left/center) + controls sidebar (right), with action buttons in the header bar. The UX spec's 3-column layout with left toolbar and floating toolbar are deferred to future enhancement. The right sidebar contains Add Text, Meme Style preset, styling controls, and text layers list.

### 2. Visual Design Authority: UI Spec
The **UI spec (`ui-spec.md`) is the single source of truth** for all visual design decisions: colors, typography, spacing, components, animations. Where the architecture spec's color summary differs from the UI spec, the UI spec wins. Key tokens:
- Primary: `#E91E8C` (Hot Magenta)
- Background: `#1A1A2E` (Deep Charcoal)
- Surface: `#252547`
- Canvas BG: `#16213E`
- Text Primary: `#F8F8F2`
- Border: `#3A3A5C`
- Success: `#36D399`
- Error: `#F87272`

### 3. Font List: 7 Fonts
Impact (default), Anton, Bebas Neue, Bangers, Permanent Marker, Oswald, Arial. These are the 6 core meme fonts from the UI spec plus Arial as a clean-text option from the product spec.

### 4. Font Size Range: 16-120px
Per product spec. The 200px value in the architecture spec is overridden.

### 5. Meme Presets: Classic Meme Only for MVP
Only the "Classic Meme" preset (Impact, white, black 3px stroke, all-caps, centered) is required for P0. The full preset gallery (Modern Clean, Handwritten, Comic Pop, Neon Glow, Minimal) from the UI spec is a P1 enhancement.

### 6. Toast Notifications: Bottom-Right, 3s
Per UI spec. Auto-dismiss after 3 seconds. Success uses green accent with icon + text. Errors use red accent with icon + text.

### 7. MVP Scope = Product Spec P0 Only
Features F1-F10 from the product spec are MVP. The UX spec's vision includes P1/P2 features (inline editing, snap guides, text resize handles, floating toolbar) — these are NOT in MVP scope but are documented for future phases.

### 8. Error Messages: User-Friendly but Specific
Follow UX spec's approach: show specific, helpful error messages ("This file type isn't supported. Please use PNG, JPG, WebP, or GIF.") rather than generic ones. These don't expose security internals — they help users fix the problem.

## Implementation Order

### Phase 1: Core Features (MVP)
1. Project scaffolding + HTML layout (Vite/Tailwind/TS, dark theme, Google Fonts)
2. Image upload with full security validation chain
3. Canvas rendering engine + text overlay CRUD
4. Drag-to-position via Pointer Events
5. Controls panel UI (font, size, color, stroke, bold/italic, alignment, caps)
6. Image export/download at full resolution

### Phase 2: Enhanced Features
7. Undo/redo system (Ctrl+Z/Ctrl+Shift+Z, max 30 snapshots)
8. Text layer list + management in sidebar
9. Meme preset button + JPEG export format option
10. Shadow controls + opacity slider
11. Keyboard shortcuts (arrow keys, Delete, Escape)

### Phase 3: Polish & Refinement
12. Empty state design, loading states, toast notifications
13. Micro-interactions + visual polish (hover effects, focus rings, glow)
14. Responsive layout (mobile stacking) + accessibility pass (WCAG AA)

## Notes for Development Team

- **Colors**: Always reference `ui-spec.md` for exact hex values, not the architecture spec summary.
- **Security**: Follow the validation chain in `security-spec.md` Section 2.1 exactly. SVG blocking and magic byte checking are non-negotiable.
- **XSS**: Use `textContent` exclusively for DOM text. Never use `innerHTML` with user input. This is a hard rule — zero exceptions.
- **Canvas rendering**: Use `requestAnimationFrame` for all re-renders. Never render directly in event handlers.
- **Export scaling**: Text positions and font sizes must scale proportionally to original image dimensions on export.
- **Memory**: Revoke blob URLs after download. Limit history to 30 snapshots. These prevent memory leaks.
- **Fonts**: Load Google Fonts with `font-display: swap` to avoid render blocking.
- **Touch support**: Use Pointer Events API (not separate mouse/touch handlers). Set `touch-action: none` on canvas.
- **Testing**: Implement the security test cases from `security-spec.md` Section 5 — especially SVG blocking, magic byte validation, and XSS prevention.

---
BEGIN DEVELOPMENT
