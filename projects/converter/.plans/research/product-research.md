# Product Research: UnitShift — Premium Minimal Unit Converter

## Competitor Analysis

| Competitor | Strengths | Weaknesses | Price |
|------------|-----------|------------|-------|
| **Google Unit Converter** (built into search) | Instant access, zero friction, supports hundreds of units, trusted brand | No personality — feels like a utility widget, not a product; cluttered search page context; no micro-interactions or visual delight; can't bookmark or save conversions | Free |
| **ConvertPad (iOS/Android)** | Comprehensive categories, favorites, history, offline support | Dated UI with dense lists and small tap targets; ad-heavy free tier feels cheap; tries to do too much (currency, cooking, data) which dilutes focus | Free w/ ads / $2.99 Pro |
| **Unit Converter by Digit Grove** | Clean layout, fast, good category organization | Generic Material Design look; no visual identity; transitions feel stock; no dark mode; banner ads break the experience | Free w/ ads |
| **Numi (macOS)** | Beautiful natural-language input, premium feel, elegant typography | Desktop-only, not web; steep learning curve for simple conversions; $20 price point limits audience | $20 one-time |

## Our Differentiation

Most unit converters fall into two camps:
1. **Utilitarian tools** (Google, ConvertPad) — functional but visually forgettable
2. **Over-scoped apps** — try to convert everything (currency, cooking, data sizes) and become cluttered

**UnitShift occupies the gap**: a focused, beautifully crafted converter for the three most common conversion categories (length, weight, temperature) that *feels like a premium product* without the complexity. Think of it as the "Linear of unit converters" — opinionated design, restrained scope, delightful details.

**Key differentiators:**
- **Micro-interactions that reward use** — smooth unit-swap animations, live value morphing, subtle haptic-feel transitions
- **Modern typography-first design** — oversized numbers, refined type scale, generous whitespace
- **Focused scope** — three categories done perfectly vs. 30 done mediocrely
- **No ads, no clutter** — the product IS the experience
- **Instant web access** — no app install required, works on any device

## Target User Persona

**Name**: "Quick-Convert Quinn"
**Demographics**: 20–40 years old, design-conscious professional or student, uses both metric and imperial units regularly (lives in the US, travels internationally, or collaborates across regions)

**Pain Points**:
- Googles "kg to lbs" multiple times a week and gets a sterile widget buried in search results
- Existing converter apps feel outdated, ad-ridden, or overwhelmingly complex
- Wants something that feels as polished as the other tools they use daily (Notion, Linear, Arc)
- On mobile, needs large tap targets and fast results — not tiny dropdowns with 200 unit options

**Current Solution**: Google search ("X kg in lbs"), occasionally a free converter app they tolerate but don't enjoy

**Why They'd Switch**:
- Bookmarkable, instant-load web app that's faster than googling
- Visually delightful — something they'd actually want to keep open or share
- Focused on the three conversions they actually need, with zero noise
- Feels premium without costing anything

## Device & Context

- **Primary**: Mobile browser (quick conversions on the go — cooking, travel, gym)
- **Secondary**: Desktop browser (work context — international collaboration, shipping dimensions)
- **Context**: Usually mid-task (cooking a recipe, packing for a trip, reading a product spec) — needs answer in < 2 seconds

## Market Positioning

**One-Liner Pitch**: "The unit converter that feels like it was designed by the people who made your favorite apps."

**Price Point**: Free (open web tool) — the value is in the experience and brand, not monetization. Could serve as a portfolio/brand showcase piece or gateway to a broader utility suite later.

**Key Differentiator**: Design-forward craftsmanship applied to an everyday utility. No one has made a *beautiful* unit converter for the web — this is the opportunity.

## UX Patterns & Best Practices to Adopt

### Input Interaction
- **Live conversion** — result updates as user types, no "Convert" button needed
- **Swap animation** — tapping a swap icon flips the from/to units with a smooth rotation
- **Smart defaults** — pre-select the most common pair per category (km ↔ mi, kg ↔ lbs, °C ↔ °F)

### Visual Design Cues
- **Large numeric display** — inspired by calculator apps (Numi, iOS Calculator), numbers are the hero
- **Minimal chrome** — no visible borders on inputs, use spacing and typography to create hierarchy
- **Subtle color coding** — each category (length/weight/temp) gets a muted accent color for identity
- **Dark mode** — default dark with light option, matching the "premium tool" aesthetic

### Micro-Interactions
- **Number morphing** — digits animate when value changes (odometer-style or fade)
- **Category tabs** — smooth underline slider when switching between length/weight/temp
- **Focus glow** — soft glow on active input field
- **Unit selector** — elegant dropdown or segmented control, not a massive scrollable list

### Typography
- **Monospace or tabular figures for numbers** — keeps digits aligned as values change
- **Sans-serif for labels** — Inter, Geist, or similar modern geometric sans
- **Type scale**: Input numbers at 2.5–3rem, unit labels at 0.875rem, category headers at 1.125rem

## Risk & Scope Considerations

| Risk | Mitigation |
|------|------------|
| "Too simple" — users might not see value | Nail the polish; simplicity IS the value proposition |
| Limited discoverability (no app store) | SEO for "unit converter", shareable URL, PWA potential |
| Scope creep (add currency! add volume!) | Hard MVP boundary: length, weight, temperature ONLY |
| Accessibility gaps in animation-heavy UI | Respect `prefers-reduced-motion`, ensure WCAG AA contrast |

---
Status: READY_FOR_PLANNING
