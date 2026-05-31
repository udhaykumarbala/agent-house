# APPROVED PLAN: UnitShift — Premium Minimal Unit Converter

Reviewed by: CEO
Date: 2026-03-17
Status: APPROVED FOR DEVELOPMENT

## Summary

UnitShift is a premium, client-side unit converter for length, weight, and temperature. Built with Vite + Tailwind CSS + TypeScript (`static-enhanced` template). Zero runtime dependencies. The visual identity is "Digital Braun" — warm obsidian dark mode with copper accent, Space Grotesk typography for numbers, Inter for UI text. Single-screen app, no backend, no navigation. The converted number is the hero.

## Approved Specifications

- Product Spec: .plans/specs/product-spec.md
- UX Spec: .plans/specs/ux-spec.md
- UI Spec: .plans/specs/ui-spec.md
- Architecture Spec: .plans/specs/architecture-spec.md
- Security Spec: .plans/specs/security-spec.md

## Key Decisions (Conflict Resolutions)

1. **Typography hierarchy**: Result value = 56px/500 weight (hero), Input value = 40px/400 weight. UI spec is authoritative. Architecture spec had these inverted — developers must follow UI spec sizes.

2. **Dark mode background**: Use `#131110` (UI spec), not `#1A1714`. All color tokens from UI spec Section "CSS Custom Properties" are the single source of truth.

3. **Category accent colors**: Length = `#D4A053` (Warm Amber), Weight = `#C47A7A` (Dusty Rose), Temperature = `#7A9EB2` (Cool Slate). Per UI spec — these are tints applied at low opacity, not replacement for the copper brand color.

4. **Unit pill border radius**: 6px (`radius-sm`) per UI spec — refined chip look, not full pill shape.

5. **Font weights**: Load Space Grotesk 400 + 500 only. Result uses 500 (not 700 bold). This keeps bundle lean.

6. **Theme toggle**: Persist preference to localStorage. Falls back to `prefers-color-scheme` on first visit.

7. **Custom numeric keypad**: CUT FROM MVP. Use native `inputmode="decimal"` on mobile. Revisit post-MVP if needed.

8. **Bidirectional input**: Move to P1 scope (desktop only). Both FROM and TO fields editable on desktop. Mobile keeps single-direction input.

## MVP Feature Scope (P0 + P1)

### P0 — Must Ship
- F1: Real-time conversion engine (instant, no button)
- F2: Category tabs with animated sliding pill
- F3: Unit pill/chip selectors (horizontal, scrollable on mobile)
- F4: Swap animation (180° rotation, spring easing, value cross-fade)
- F5: Typography-first large number display (Space Grotesk hero numbers)
- F6: Dark mode default (warm obsidian #131110 + copper #C17F59)
- F7: Responsive layout (mobile-first, 480px max-width card on desktop)
- F8: Accessible & motion-safe (keyboard nav, ARIA, reduced-motion)

### P1 — Should Ship
- F9: Light mode toggle with persistence
- F10: Number morphing animation (digit fade/roll on value change)
- F11: Smart defaults & localStorage memory (last category + units)
- F12: Copy result (tap/click result → clipboard + checkmark confirmation)
- F13: Quick reference formula below result
- F14: Bidirectional input (desktop only — type in either field)

### Post-MVP (Not in scope)
- Custom numeric keypad (mobile)
- PWA / service worker
- Keyboard shortcuts
- Noise texture overlay

## Implementation Order

1. **Project scaffolding** — Vite + Tailwind + TypeScript setup, design tokens (CSS custom properties), font loading, CSP meta tag
2. **Conversion engine** — Pure math module (`engine.ts`, `units.ts`, `types.ts`), all 3 categories, input sanitization
3. **Core UI** — HTML structure, category tabs, FROM/TO cards, unit pill selectors, responsive layout
4. **Live conversion** — Wire input → engine → output, real-time updates, number formatting
5. **Swap animation** — Swap button with 180° rotation, value cross-fade, unit swap
6. **Micro-interactions** — Tab sliding pill, chip selection transitions, focus glow, card elevation
7. **Light mode** — Theme toggle, light color tokens, localStorage persistence
8. **Polish** — Copy result, quick reference formula, bidirectional input (desktop), number morphing, edge cases (absolute zero warning, overflow)
9. **Accessibility audit** — Keyboard nav, ARIA labels, screen reader testing, contrast verification
10. **Security hardening** — CSP verification, no innerHTML check, npm audit, source maps disabled

## Notes for Development Team

- **UI spec is the single source of truth for all visual decisions** (colors, spacing, typography, radii, shadows). When any other spec conflicts with UI spec on visual matters, UI spec wins.
- **Product spec is the source of truth for feature scope**. If UX spec shows a feature not in P0/P1, check product spec priority first.
- **Zero runtime dependencies** — all conversion math, formatting, and DOM manipulation is hand-written. No lodash, no animation libraries, no state management.
- **Security**: Use `textContent` exclusively (never `innerHTML`) for rendering user values. Include CSP meta tag. Validate all input as numeric. Disable source maps in prod.
- **Performance target**: <20KB gzipped total (excluding fonts), <500ms FCP, 95+ Lighthouse.
- **Animations**: CSS transitions + keyframes only. All animations must respect `prefers-reduced-motion: reduce`. Use `cubic-bezier(0.34, 1.56, 0.64, 1)` for spring effects (swap, tab slide).
- **The custom keypad is deliberately excluded** — do not build it. Native mobile keyboard with `inputmode="decimal"` is the input method.

---
BEGIN DEVELOPMENT
