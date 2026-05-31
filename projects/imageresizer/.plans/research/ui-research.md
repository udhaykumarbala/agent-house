# UI Research: Image Resizer (Friendly Brutalism)

## Design Direction: Minimal & Friendly Brutalism

Brutalism in web design embraces raw aesthetics — bold typography, thick borders, exposed structure, and unapologetic use of color. "Friendly brutalism" softens the edges: it keeps the boldness and honesty but adds warmth, playfulness, and approachability. This is NOT cold, harsh, or intimidating — it's confident and fun.

---

## Design Inspiration

### Gumroad
- **Visual Style**: Friendly brutalism pioneer — thick black borders, solid color fills, no gradients, bold sans-serif type, playful illustrations
- **What Makes It Great**: Feels hand-crafted and human despite being a tech product. The thick outlines and flat colors create a comic-book quality that feels approachable. High contrast black borders on colored surfaces give everything a sticker/cut-out feel.
- **What We Can Borrow**: Thick 2-3px black borders on cards and buttons, solid flat color backgrounds, the "sticker on a surface" elevation style using offset box-shadows instead of soft drop shadows.

### Poolsuite (formerly Poolside FM)
- **Visual Style**: Retro-meets-brutalism with bold color blocking, chunky UI elements, nostalgic type, and a playful irreverence. Uses cream/warm backgrounds with saturated accent colors.
- **What Makes It Great**: Proves brutalism can feel warm and inviting. The warm background palette prevents the stark borders from feeling cold. Every element feels intentional and confident.
- **What We Can Borrow**: Warm background tones (cream, butter, soft peach) instead of pure white or gray. The idea that brutalist structure + warm palette = friendly.

### Figma's Community/Config Branding
- **Visual Style**: Bold geometric shapes, thick strokes, vibrant but not neon colors, playful compositions. Uses a mix of solid fills and strong outlines.
- **What Makes It Great**: Accessible brutalism — it's bold and graphic without being hostile. The color palette is warm and energetic (corals, yellows, warm purples). Typography is confident but readable.
- **What We Can Borrow**: The warm, energetic color palette approach. Bold section headers. Using geometric shapes as decorative elements to break visual monotony in a utility tool.

### Notion (Calendar & Recent Redesign Elements)
- **Visual Style**: Clean minimal with brutalist touches — the calendar product uses thick borders, bold type, and flat colors. The rest of Notion keeps things minimal with just enough personality.
- **What Makes It Great**: Shows that utility tools can have personality. The calendar's thick-bordered date cells feel sturdy and tactile. Proves minimalism and brutalism can coexist.
- **What We Can Borrow**: The idea of a utility-first interface where brutalist styling enhances clarity rather than decoration. Minimal chrome, maximum function, but with bold visual anchors.

### TinyPNG
- **Visual Style**: Playful, minimal, single-purpose. Uses a panda mascot and a simple green/white palette. The interface is stripped to essentials — a drop zone and results.
- **What Makes It Great**: Proves an image tool can be dead simple. One screen, one action, instant feedback. The personality comes from the mascot and micro-copy, not UI complexity.
- **What We Can Borrow**: The radical simplicity of a single-purpose tool. Big drop zone, immediate feedback, minimal navigation. But we'll replace their soft/generic style with brutalist confidence.

---

## Competitor Color Analysis

| Competitor | Primary Color | Style | Notes |
|------------|--------------|-------|-------|
| TinyPNG | Green #6DB33F | Soft/friendly green | Nature-themed, overused in utility tools — AVOID |
| Squoosh (Google) | Blue #4285F4 | Google Material blue | Corporate tech blue — AVOID, too generic |
| iLoveIMG | Orange #FF6F3C | Warm utility orange | Better differentiation, but now associated with their brand |
| ResizeImage.net | Blue #2196F3 | Material Design blue | Generic — AVOID |
| Bulk Resize Photos | Teal #00BCD4 | Material teal | Slightly different but still common tech palette |
| Canva (resize tool) | Purple #7B2FBE + Teal | Gradient tech | Overused AI/creative purple — AVOID |
| Photopea | Gray #444 | Dark utilitarian | Functional but lifeless — AVOID the dullness |
| Pixlr | Blue #1B8EF2 | Yet another blue | Confirms blue is the most overused color in image tools |

**Key Insight**: Image tools are drowning in blues, greens, and generic material colors. There's a massive opportunity to stand out with warm, bold, brutalist-inspired colors.

---

## Our Differentiation Strategy

To stand out, we should:
- **Avoid**: Blue (every image tool uses it), generic green, material design purple, soft/rounded "friendly" aesthetics that look like every other SaaS
- **Avoid**: Gradients, glassmorphism, soft shadows — these are the opposite of brutalism
- **Consider**: A warm, bold palette rooted in friendly brutalism — think butter yellow, warm coral, or terracotta paired with thick black borders
- **Lean Into**: The brutalist identity as our brand differentiator. No other image resizer looks brutalist. This IS our moat.

---

## Color Psychology

**App Mood**: Confident, playful, efficient — like a friendly shop teacher who gets things done with a smile

The app should feel:
- **Trustworthy** — "I'll resize your image correctly"
- **Fast** — "This won't waste your time"
- **Fun** — "Utility doesn't have to be boring"
- **Bold** — "I know exactly what I am"

### Recommended Color Directions

1. **Butter & Charcoal** — Warm butter yellow (#F5E6A3) as the primary brand color with charcoal black (#1A1A1A) borders and text. Feels like a friendly workshop — warm, capable, not sterile. Coral (#E8625C) for accents and CTAs. This is the strongest direction for "friendly brutalism."

2. **Terracotta & Cream** — Rich terracotta/burnt orange (#C75C2E) as primary on a cream (#FFF8F0) background with thick black borders. Earthy, warm, handmade feel. Sage green (#7A9E7E) for success states. Feels artisanal and tactile.

3. **Coral & Off-White** — Vibrant coral (#E85D50) as primary on warm off-white (#FEFBF6) with black structural elements. Energetic and approachable. Slate (#3D4F5F) for text. Feels like a Scandinavian design studio poster.

**Strongest Recommendation**: Direction 1 (Butter & Charcoal) — it's the most unique in the image tool space, the most aligned with friendly brutalism, and the most memorable. Nobody else uses warm yellow as a primary in this category.

---

## Typography Research

Brutalism demands typography with presence — no thin, wispy fonts. We need bold, confident type that holds its own against thick borders and strong colors.

### Display/Heading Font Options

- **Space Grotesk** (Google Fonts): Geometric sans-serif with personality. The slightly quirky letterforms (especially the 'a' and 'g') add warmth to its bold weights. Perfect for friendly brutalism — it's bold without being aggressive. **TOP PICK.**
- **DM Sans**: Clean geometric with a friendly feel. Less quirky than Space Grotesk but very readable. Good fallback if Space Grotesk feels too distinctive.
- **Sora**: Modern geometric with a slightly rounded quality. Feels contemporary and friendly. Works well in bold weights for brutalist headings.
- **Archivo Black** (display only): Ultra-bold, condensed display face. Great for hero text and oversized labels. Pure brutalist energy but might need a softer body font to balance.

### Body Font Options

- **Inter**: The workhorse. Excellent readability at all sizes, great for UI labels, inputs, and body text. Pairs beautifully with bolder display fonts. **TOP PICK for body.**
- **DM Sans**: Can double as both display and body if we want a single-font system. Keeps things simple and consistent.
- **IBM Plex Sans**: Slightly more mechanical feel that suits brutalism's honest aesthetic. Very readable.

### Recommended Pairing

**Space Grotesk** (headings, bold UI elements) + **Inter** (body, labels, inputs)

This gives us personality in headlines and reliability in body text. Both are free on Google Fonts.

---

## Visual Style Direction

### Recommended Style: Friendly Brutalism ("Neo-Brut")

**Core Principles:**
- **Thick black borders** (2-3px) on interactive elements — buttons, cards, inputs
- **Offset box-shadows** instead of soft shadows — e.g., `4px 4px 0px #1A1A1A` creates a sticker/stamp effect
- **Flat, solid colors** — no gradients, no opacity tricks
- **Bold, oversized typography** for headings and key UI labels
- **Warm background** — cream or butter, NOT white or gray
- **Minimal decoration** — the structure IS the decoration
- **Generous whitespace** — brutalism needs room to breathe

### Corners
**Mix of sharp and slightly rounded**: Primary elements (buttons, cards) use small radius (4-6px) to keep the brutalist edge while preventing them from feeling hostile. No fully rounded pills — they're too soft for this aesthetic.

### Shadows
**Hard offset shadows only**: `3px 3px 0 #1A1A1A` or `4px 4px 0 #1A1A1A`. These create the signature brutalist "stacked paper" or "sticker" effect. NO soft/diffused shadows — that's the opposite of what we want.

### Borders
**Thick and black**: 2-3px solid #1A1A1A on cards, buttons, inputs. This is the single most important brutalist element. Everything should feel like it has a confident outline.

### Animations
**Minimal and snappy**:
- Button hover: translate(-2px, -2px) with shadow growing to 5px 5px — feels like the element is lifting off the page
- Button active/click: translate(2px, 2px) with shadow shrinking to 0 — feels like pressing a physical button
- Drag-and-drop zone: dashed border animation on hover
- No spring animations, no bouncy easing — keep it mechanical and satisfying

### Iconography
**Thick-stroke line icons** (2-2.5px stroke) to match the border weight. Options:
- Phosphor Icons (bold weight) — excellent brutalist-compatible icon set
- Lucide with increased stroke weight
- Custom minimal icons if needed

### Layout
- **Single column, centered** — the tool should feel focused
- **Large drop zone** as the hero element — at least 50% of viewport on desktop
- **Chunky controls** — oversized sliders, big number inputs, large buttons
- **Card-based sections** with thick borders for settings/options
- **Fixed/sticky action bar** at bottom for the resize/download CTA

---

## Micro-Interactions & Polish

- **Drag hover**: Drop zone border animates (marching ants or color pulse)
- **File loaded**: Satisfying snap/pop animation as the preview appears
- **Resize in progress**: Bold progress bar with thick border, no spinner
- **Download ready**: Button grows slightly or does a subtle "tada" translate
- **Error state**: Red thick border + shake animation (2-3 quick horizontal shakes)

---

## Accessibility Notes for Brutalism

Brutalism's high-contrast aesthetic actually HELPS accessibility:
- Thick borders provide clear visual boundaries
- Bold type is easier to read
- Solid colors maintain contrast ratios easily
- The warm background (butter/cream) reduces eye strain vs. pure white

Must ensure:
- All text meets WCAG 4.5:1 contrast on chosen backgrounds
- Interactive elements have clear focus states (double border or outline offset)
- Touch targets remain 44x44px minimum despite the bold styling

---

Status: READY_FOR_PLANNING
