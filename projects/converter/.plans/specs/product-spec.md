# Product Specification: UnitShift — Premium Minimal Unit Converter

Based on research in:
- `.plans/research/product-research.md`
- `.plans/research/ux-research.md`
- `.plans/research/ui-research.md`
- `.plans/research/architecture-research.md`
- `.plans/research/security_expert.md`

---

## Target User

**"Quick-Convert Quinn"** — A design-conscious professional or student (20–40), who regularly converts between metric and imperial units for work, travel, cooking, or fitness. Currently uses Google search ("5 kg in lbs") and tolerates ad-ridden converter apps but wants something that feels as polished as the tools they use daily (Notion, Linear, Arc). Primarily on mobile, mid-task, needs an answer in under 5 seconds.

## Unique Value Proposition

**"The unit converter that feels like it was designed by the people who made your favorite apps."**

Free, focused, and beautifully crafted — three conversion categories done perfectly, with zero ads and zero clutter. The Linear of unit converters.

---

## Features (Prioritized)

### P0 — Must Have (MVP)

- [ ] **F1: Real-Time Conversion Engine** — Accurate, instant conversion for length, weight, and temperature as the user types. No "Convert" button. This is the core utility — without it, nothing else matters.
- [ ] **F2: Category Tabs** — Segmented control to switch between Length, Weight, and Temperature. Each category gets a subtle muted accent color for identity. Essential for organizing the three conversion domains.
- [ ] **F3: Unit Selector (Pill/Chip Style)** — Horizontally scrollable pill selectors for source and target units within each category. Replaces dropdown menus (an anti-pattern per UX research). Must be thumb-friendly with large tap targets (min 48px).
- [ ] **F4: Swap Animation** — A prominent swap button between from/to that rotates and cross-fades values. This is the signature micro-interaction that elevates the app from "tool" to "experience."
- [ ] **F5: Typography-First Large Number Display** — Oversized numeric display (48–64px) using Space Grotesk for values and Inter for labels. Numbers are the hero — monospaced/tabular figures to prevent layout jank during live conversion.
- [ ] **F6: Dark Mode Default** — Premium dark theme (#1A1714 warm obsidian base) with warm copper (#C17F59) accent. Dark mode is the primary brand identity per UI research. Must meet WCAG AA contrast (4.5:1).
- [ ] **F7: Responsive Layout** — Mobile-first single-screen design. Centered card (max-width 480px) on desktop. No page navigation for the core task. Works on any device and any modern browser.
- [ ] **F8: Accessible & Motion-Safe** — Full keyboard navigation, ARIA labels on all interactive elements, screen reader announces results, `prefers-reduced-motion` respected for all animations.

### P1 — Should Have

- [ ] **F9: Light Mode Toggle** — User can switch to a light theme (warm stone palette). Respects `prefers-color-scheme` for initial default.
- [ ] **F10: Number Morphing Animation** — Digits animate with a subtle roll/fade transition when the converted value updates. Adds the "alive" feeling identified in UX research (Elk, Amount patterns).
- [ ] **F11: Smart Defaults & Memory** — Pre-select the most common unit pair per category (km ↔ mi, kg ↔ lbs, °C ↔ °F). Remember last-used category and units via localStorage.
- [ ] **F12: Copy Result** — Tap/click the converted value to copy it. Subtle flash + checkmark morph confirmation. Users are often mid-task and need to paste the value elsewhere.
- [ ] **F13: Quick Reference Formula** — Display the conversion formula below the result (e.g., "1 km = 0.6214 mi"). Educates the user and builds trust in the result.

### P2 — Nice to Have (Post-MVP)

- [ ] **F14: Bidirectional Input** — Both from and to fields are editable. Typing in either triggers conversion in the opposite direction. Matches Google Converter convenience.
- [ ] **F15: PWA Support** — Service worker for offline use, installable on home screen. Makes it a true app replacement.
- [ ] **F16: Keyboard Shortcuts (Desktop)** — Tab to switch categories, Enter to swap, Cmd/Ctrl+C to copy result. Inspired by Linear's keyboard-first approach.
- [ ] **F17: Subtle Noise Texture** — 2–3% opacity background noise for tactile quality (per UI research). Adds the "digital Braun" materiality.

---

## User Stories

1. **As a traveler**, I want to type a distance in kilometers and instantly see it in miles, so that I can understand driving distances abroad without Googling.
2. **As a home cook**, I want to quickly convert ounces to grams while following a recipe, so that I don't have to leave my recipe app for long.
3. **As a fitness enthusiast**, I want to swap between kg and lbs with one tap, so that I can track my gym lifts in the units my program uses.
4. **As a student**, I want to convert Celsius to Fahrenheit and see the formula, so that I learn the relationship while getting the answer.
5. **As a mobile user**, I want the converter to load instantly and work with large tap targets, so that I can get an answer in under 5 seconds while on the go.
6. **As a design-conscious user**, I want the app to look and feel premium with smooth animations, so that I enjoy using it and want to share it.
7. **As a user with motion sensitivity**, I want animations to be disabled when I have reduced-motion enabled, so that I can use the tool comfortably.
8. **As a desktop user**, I want to copy the converted result with one click, so that I can paste it into a document or message without retyping.

---

## Acceptance Criteria

### F1: Real-Time Conversion Engine
- [ ] Conversion result updates within 50ms of each keystroke (no perceptible delay)
- [ ] Length and weight use base-unit normalization (value × fromFactor / toFactor)
- [ ] Temperature uses direct formula paths (C↔F, C↔K, F↔K — 6 paths)
- [ ] Displays up to 6 significant digits, trims trailing zeros
- [ ] Handles edge cases gracefully: empty input shows placeholder, NaN/Infinity shows "—", negative temperatures are valid, negative length/weight shows result but no error
- [ ] Input accepts digits, one decimal point, and negative sign only — strips all other characters

### F2: Category Tabs
- [ ] Three tabs visible: Length, Weight, Temperature
- [ ] Active tab has an animated sliding pill/underline background
- [ ] Each category has a distinct muted accent color (e.g., amber for length, green for weight, blue for temperature)
- [ ] Switching categories preserves any entered value (re-converts with new default units)
- [ ] Tab transition uses smooth horizontal slide with content cross-fade

### F3: Unit Selector (Pill/Chip Style)
- [ ] Units displayed as horizontal pill buttons, not dropdowns
- [ ] Active unit pill is visually highlighted with accent color
- [ ] Length units: mm, cm, m, km, in, ft, yd, mi
- [ ] Weight units: mg, g, kg, oz, lb, t
- [ ] Temperature units: °C, °F, K
- [ ] Pill row scrolls horizontally if units overflow on mobile
- [ ] Minimum touch target: 48×48px

### F4: Swap Animation
- [ ] Swap button centered between from/to sections
- [ ] On tap: button rotates 180°, from/to units exchange positions, values cross-fade and recalculate
- [ ] Animation uses spring-based easing (not linear)
- [ ] Works via keyboard (Enter or Space when focused)

### F5: Typography-First Large Number Display
- [ ] Primary result value: 48–64px, Space Grotesk, bold weight
- [ ] Input value: 24–32px
- [ ] Labels and unit names: Inter, 14px
- [ ] All numbers use tabular/monospaced figures (no layout shift on digit change)

### F6: Dark Mode Default
- [ ] Default theme is dark on first visit
- [ ] Background: warm near-black (#1A1714 or similar)
- [ ] Card surfaces: slightly lighter with subtle border (semi-transparent)
- [ ] Accent color: warm copper (#C17F59)
- [ ] All text meets WCAG AA contrast ratio (4.5:1 for normal text, 3:1 for large)

### F7: Responsive Layout
- [ ] Mobile (< 640px): full-width, keypad-friendly spacing, tabs at top
- [ ] Tablet (640–1024px): centered with wider pill rows
- [ ] Desktop (> 1024px): centered card, max-width 480px, keyboard input, hover states
- [ ] No horizontal scrolling on any viewport except unit pills row

### F8: Accessible & Motion-Safe
- [ ] All interactive elements reachable via Tab key
- [ ] Swap button and unit pills have visible focus indicators
- [ ] Screen reader announces: selected category, selected units, and conversion result on change
- [ ] `prefers-reduced-motion: reduce` disables all animations and transitions
- [ ] No reliance on color alone to convey information

---

## Success Metrics

- **Time to result**: User gets converted value in < 5 seconds from page load (measured via performance audit)
- **Bundle size**: < 20KB gzipped total (HTML + CSS + JS + fonts excluded)
- **Lighthouse score**: 95+ on Performance, 100 on Accessibility
- **First Contentful Paint**: < 1 second on 4G connection
- **Interaction latency**: Conversion updates in < 50ms (no perceptible input lag)

---

## Out of Scope (for MVP)

- Currency conversion (adds API dependency, changes architecture from static to connected)
- Volume, area, speed, data size, or any category beyond length/weight/temperature
- User accounts, login, or cloud sync
- Conversion history or saved favorites
- Natural language input ("5 kg in pounds")
- Custom themes or theme editor
- Backend or API of any kind
- App store distribution (web-only for MVP)
- Internationalization / multi-language support

---

## Security Requirements (from security research)

- Use `textContent` (never `innerHTML`) for all DOM rendering of user input
- Validate all input as numeric before processing
- Include CSP meta tag in `index.html`
- Pin exact dependency versions; run `npm audit` before deployment
- Disable source maps in production build
- Add security headers on deployment (X-Frame-Options, X-Content-Type-Options, Referrer-Policy, Permissions-Policy)

---

## Technical Constraints (from architecture research)

- **Template**: `static-enhanced` (Vite + Tailwind CSS + TypeScript)
- **No runtime dependencies**: Zero npm packages beyond dev tooling
- **Fonts**: Space Grotesk (numbers) + Inter (UI) — both Google Fonts
- **Animations**: CSS transitions + keyframes only, no JS animation library
- **State**: Simple object + event listeners, no state management library

---

Status: READY_FOR_REVIEW
