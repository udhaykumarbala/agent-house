# UI Research: Unit Converter (Next-Gen Minimal)

## Design Inspiration

### Elk (Unit Converter by Clean Shaven Apps)
- **Visual Style**: Minimal, single-screen interface with a muted warm color scheme, large type for values, and a subtle translucent material backdrop. Uses haptic-like micro-animations on unit swap.
- **What Makes It Great**: It strips the converter down to essentials — one input, one output, a unit picker. No clutter. The typography IS the UI. Numbers feel like they have weight and presence.
- **What We Can Borrow**: The "typography as interface" approach — oversized numbers, minimal chrome. The single-screen constraint that forces simplicity.

### Amounts (Currency Converter by Wunderbucket)
- **Visual Style**: Dark mode with a single accent color (warm amber/gold), monospaced number display, smooth slide-to-convert gesture. Very iOS-native feel with depth through layered surfaces.
- **What Makes It Great**: The input feels physical — numbers slide, snap, and settle. The dark background makes the gold accent feel premium. Zero unnecessary UI elements.
- **What We Can Borrow**: The layered dark surface approach — background, card, elevated card — creating depth without borders. The warm accent on dark canvas for a premium feel.

### Nuumi (Calculator by Stacking)
- **Visual Style**: Cream/off-white background with charcoal text, a single coral accent, generous whitespace, and a custom sans-serif that feels editorial. Animations are spring-based and feel organic.
- **What Makes It Great**: It feels like a well-designed print piece that happens to be interactive. The restraint is extreme — one accent color, one font weight hierarchy, nothing extra.
- **What We Can Borrow**: The editorial quality — treating a utility app like a design object. The spring-based micro-interactions that make inputs feel alive.

### Apple Calculator (iOS 18)
- **Visual Style**: Pure black background, orange accent, large SF Pro Display numbers, satisfying haptic feedback on key press. The grid layout is mathematically precise.
- **What Makes It Great**: Ultimate clarity. Every element earns its place. The contrast between black and orange is iconic and instantly recognizable.
- **What We Can Borrow**: The clarity-first approach — never sacrifice readability for aesthetics. The way the primary action (equals/convert) gets the only color.

### Linear (Project Management)
- **Visual Style**: Ultra-clean with a monochromatic palette punctuated by contextual color. Subtle gradients, refined micro-animations on state changes, keyboard-first interactions.
- **What Makes It Great**: The attention to detail in transitions — elements don't just appear, they arrive with purpose. The surface materials feel layered and tactile despite being flat.
- **What We Can Borrow**: The polish in transitions and state changes. When a user converts a value, the result should arrive with intention, not just pop in.

## Competitor Color Analysis

| Competitor | Primary Color | Style | Notes |
|------------|--------------|-------|-------|
| Google Unit Converter | Blue #4285F4 | Standard Google Material | AVOID - generic Google blue |
| Apple Calculator | Orange #FF9F0A | Bold accent on black | Iconic but owned by Apple |
| Elk Converter | Warm Gray #8B8178 | Muted earth tones | Nice restraint, but too muted |
| Unit Converter (popular Android) | Teal #009688 | Material Design teal | AVOID - dated Material look |
| Convertio | Blue #2979FF | Standard tech blue | AVOID - generic and forgettable |
| XE Currency | Blue/Navy #00457C | Financial blue | AVOID - corporate feel |
| Amounts | Amber/Gold #D4A853 | Warm premium accent | Good direction but currency-specific |
| CalcBot | Purple #7C4DFF | Playful tech purple | AVOID - overused AI/tech purple |

## Our Differentiation Strategy

To stand out, we should:
- **Avoid**: Blue (every converter uses it), teal/green (dated Material Design), purple (AI/tech cliche), plain white backgrounds (generic)
- **Avoid**: Heavy Material Design or Bootstrap-looking patterns — they scream "template"
- **Consider**: A warm, sophisticated palette that feels more like a luxury product than a utility app
- **Consider**: Deep charcoal or warm off-black as the base — dark mode done right, not just "invert colors"
- **Consider**: A single, unexpected accent color that becomes our signature — something in the terracotta/burnt sienna/warm copper range that no converter uses

## Color Psychology

**App Mood**: Precise yet warm. Trustworthy but not corporate. Premium but not pretentious. Like a well-made tool — a Leica camera or a Braun calculator.

Recommended color directions:

1. **Warm Obsidian + Copper**: Deep charcoal-black (#1A1714) as the canvas with burnished copper (#C17F59) as the sole accent. This feels like a premium instrument — think Dieter Rams meets Japanese craft. The warmth in the black (slight brown undertone) avoids the cold/tech feeling. Copper signals quality without being flashy.

2. **Ink + Terracotta**: Near-black with a blue-green undertone (#0F1419) paired with muted terracotta (#C4654A). This has an editorial quality — like a beautifully typeset reference book. The terracotta is unexpected for a converter and creates instant differentiation.

3. **Stone + Saffron**: Warm stone gray (#E8E2D9) as a light-mode base with deep saffron/turmeric (#C7882A) as accent, and dark walnut (#2C2419) for text. This feels like a Scandinavian design object — natural materials, intentional warmth, nothing synthetic. Unique in a space dominated by cold whites and blues.

## Typography Research

**Display Font Options** (for large numbers — the hero of the UI):
- **Geist Mono** (by Vercel): Clean, modern monospace with excellent number forms. The monospace ensures numbers align perfectly during conversion animations. Feels technical but refined.
- **Inter**: Excellent tabular number support, optimized for screens. The variable font allows precise weight tuning. The "cv01" stylistic set gives distinct number forms.
- **Space Grotesk**: Geometric sans with personality. Its numbers have a distinctive quality — slightly quirky but highly legible. Would give the app character.
- **JetBrains Mono**: Beautiful number forms with ligature support. More technical feeling, but the number design is exceptional.

**Body Font Options** (for labels, unit names):
- **Inter**: The workhorse — legible at all sizes, excellent language support, tabular figures for aligned data.
- **Manrope**: Geometric sans-serif with a warm, slightly rounded quality. Pairs well with monospace display fonts.
- **Satoshi**: Modern geometric with subtle humanist qualities. Feels fresh and contemporary without being trendy.

**Recommended Pairing**: **Space Grotesk** for display numbers (gives character and distinctiveness) + **Inter** for body/labels (reliable, clear, great at small sizes). The contrast between the slightly expressive display and the neutral body creates visual hierarchy naturally.

## Visual Style Direction

**Recommended style**: Modern minimal with warm material quality — "digital Braun"

- **Corners**: Rounded (12-16px for cards, 8px for inputs/buttons) — soft but not bubbly. Avoid fully rounded pill shapes except for small tags/badges.
- **Shadows**: Minimal — prefer layered surfaces with subtle background color shifts over drop shadows. If shadows are used, they should be warm-tinted (rgba with brown, not pure black).
- **Animations**:
  - Spring-based micro-interactions (not linear easing)
  - Number values should animate between states with a rolling/counting effect
  - Unit swap should feel physical — like flipping a card or rotating a dial
  - Subtle scale feedback on tap (0.97 → 1.0 spring)
  - Result should "settle" into place with slight overshoot
- **Layout**: Single-screen, centered content, generous vertical rhythm. The conversion should feel like a conversation: input → transform → output, flowing top to bottom.
- **Iconography**: Minimal line icons, 1.5px stroke, rounded caps. Only where truly needed (swap, settings). Prefer text labels over icons when space allows.
- **Surface Treatment**: Subtle noise texture on backgrounds (2-3% opacity) to add tactile quality. Slight gradient on primary surfaces to suggest depth without obvious gradients.

## Micro-Interaction Concepts

1. **Value Input**: Numbers type in with a subtle scale-up spring. Cursor blinks with a warm glow, not harsh.
2. **Live Conversion**: As user types, the output value counts/rolls to the new number in real-time. Slight delay (50ms) for perceived intelligence.
3. **Unit Swap**: The from/to units swap with a smooth rotation animation. The values cross-fade and recalculate.
4. **Category Switch** (Length/Weight/Temp): Horizontal slide with content parallax — the new category slides in while old slides out, numbers at a different speed than labels.
5. **Copy Result**: Subtle flash + checkmark morph on the value, brief haptic pulse.

---
Status: READY_FOR_PLANNING
