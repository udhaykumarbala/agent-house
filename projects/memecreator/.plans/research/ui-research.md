# UI Research: Meme Creator

## Design Inspiration

### Canva (canva.com)
- **Visual Style**: Clean white workspace with a left-side panel for tools/templates, vibrant accent colors (teal/purple gradient brand). The canvas sits center-stage on a neutral gray background, drawing full focus to the creation area.
- **What Makes It Great**: The toolbar is contextual — it morphs based on what element you select (text tools appear when text is selected, image tools when an image is selected). Drag handles are intuitive, and the property panel stays out of the way until needed.
- **What We Can Borrow**: The canvas-centric layout with floating toolbars. The idea of a clean workspace where the image dominates and tools orbit around it. Also, the snap-to-center guides when dragging elements.

### Imgflip (imgflip.com)
- **Visual Style**: Utilitarian, older-style web design. Dark header, white content area, very tool-focused with meme templates front and center.
- **What Makes It Great**: It's the go-to meme generator for a reason — zero friction. Pick template, add text, download. The text placement is dead simple with classic top/bottom meme format.
- **What We Can Borrow**: The speed-to-creation flow. Users want to go from idea to meme in under 30 seconds. We should match that speed but with a far more polished visual experience.

### Kapwing (kapwing.com)
- **Visual Style**: Modern, minimal with a dark sidebar and light canvas. Uses a soft indigo/violet brand color. The editor feels like a lightweight Photoshop — layers panel, timeline, property editors.
- **What Makes It Great**: Full creative freedom — you can place text anywhere, resize freely, add multiple layers. The export flow is clean with format/quality options. It bridges the gap between "simple meme maker" and "full editor."
- **What We Can Borrow**: The free-placement text system with drag handles and rotation. The property panel for text styling (font, color, stroke, shadow). The dark sidebar + light canvas contrast that separates tools from content.

### Photopea (photopea.com)
- **Visual Style**: Full Photoshop clone — dark UI, dense toolbars, professional feel. Uses a teal/cyan accent on a charcoal base.
- **What Makes It Great**: Incredibly powerful for a browser app. Shows that complex image editing is fully achievable in the browser with canvas APIs.
- **What We Can Borrow**: The technical proof that canvas-based text overlay with full styling control works well in browsers. Their approach to text rendering on canvas (stroke, shadow, alignment) is the gold standard.

### Adobe Express (express.adobe.com)
- **Visual Style**: Sleek, modern with gradient accents (coral-to-violet). Left panel navigation, center canvas, right properties panel. Feels premium without being intimidating.
- **What Makes It Great**: The text styling options are extensive but presented simply — font picker with visual previews, color picker with brand colors, effects like shadow/outline presented as one-click presets.
- **What We Can Borrow**: The preset text styles concept — "Meme Classic", "Bold Impact", "Subtle Caption" — so users get great-looking text instantly without tweaking every property. Also the clean export modal.

## Competitor Color Analysis

| Competitor | Primary Color | Style | Notes |
|------------|--------------|-------|-------|
| Canva | Teal #00C4CC + Purple #7D2AE8 | Gradient, vibrant | AVOID - strongly associated with Canva |
| Imgflip | Blue #0A74DA | Generic tech blue | AVOID - dated and generic |
| Kapwing | Indigo #6366F1 | Modern tech purple-blue | AVOID - common SaaS purple |
| Photopea | Teal #18A0A0 | Professional dark | AVOID - too close to Canva's teal |
| Adobe Express | Coral-Violet gradient | Premium, warm | AVOID - Adobe brand territory |
| Mematic (iOS) | Orange #FF6B35 | Playful, warm | Noted - warm tones work for creative tools |
| CapCut | Black + Cyan #00F0FF | High-contrast, Gen-Z | Noted - dark mode with neon is trending |

## Our Differentiation Strategy

To stand out, we should:
- **Avoid**: Teal, standard blues, indigo/violet, coral-violet gradients — these are already claimed by major competitors
- **Avoid**: Generic gray toolbars that look like every other SaaS
- **Consider**: A warm, energetic direction that says "fun and creative" — meme creation is playful, not corporate
- **Consider**: Dark mode as default (meme culture is internet culture, which skews dark-mode-native)
- **Consider**: A bold, unexpected primary color — something that pops on dark backgrounds and feels irreverent like meme culture itself

**Direction**: Dark workspace with a warm, punchy accent color. Think "creative studio at night" — the focus is on the image/canvas, while the UI wraps around it in a dark, non-distracting shell with energetic accent pops.

## Color Psychology

**App Mood**: Playful, creative, fast, slightly irreverent — but still polished and trustworthy

Meme creation is about humor, speed, and self-expression. The UI should feel:
- **Fun** — not corporate or sterile
- **Fast** — dark backgrounds reduce visual noise, keeping focus on the canvas
- **Expressive** — accent colors should feel energetic and bold
- **Modern** — Gen-Z and millennial users expect contemporary design

Recommended color directions:

1. **Hot Magenta + Charcoal** — `#E91E8C` on `#1A1A2E` — Electric, bold, unmistakably creative. Magenta is rarely used as a primary in this space, making it highly distinctive. Associated with creativity, energy, and fun. High contrast on dark backgrounds.

2. **Amber/Saffron + Deep Slate** — `#F59E0B` on `#0F172A` — Warm, inviting, energetic. Gold/amber conveys creativity and optimism. Uncommon in the meme/image editor space. Feels premium yet approachable.

3. **Electric Lime + Near-Black** — `#AAFF00` on `#121212` — High-energy, internet-native, bold. Lime/chartreuse screams "creative tool" and stands out against every competitor. Very Gen-Z. Risk: may feel too aggressive for some users.

**Recommended**: Direction 1 (Hot Magenta + Charcoal). It's distinctive, energetic, gender-neutral, and reads as "creative tool" immediately. Magenta has strong associations with imagination and artistic expression while being virtually unused by competitors.

## Typography Research

Meme text and UI text have different needs:

### Meme Text (On-Canvas)
The classic meme font is **Impact** (bold, condensed, white with black stroke). We must support this but also offer modern alternatives:
- **Impact**: The OG meme font — must be available as a preset
- **Anton**: Google Font alternative to Impact, slightly more refined
- **Bebas Neue**: Modern condensed sans-serif, popular in YouTube thumbnails
- **Permanent Marker**: Handwritten/marker style for casual memes
- **Bangers**: Comic/pop-art style, great for expressive memes
- **Oswald**: Clean condensed option for more "designed" memes

### UI Font (App Interface)

**Display Font Options**:
- **Inter**: The modern web standard — extremely legible, great for UI. Has variable weight support. Slightly safe but proven.
- **Geist Sans (by Vercel)**: Fresh, geometric, modern feel. Gaining popularity in dev/creative tools. Slightly more personality than Inter.
- **DM Sans**: Geometric sans with a friendly, rounded quality. Feels approachable and creative without being childish.

**Body Font Options**:
- **Inter**: Best readability for small UI text, extensive character support
- **DM Sans**: Pairs well as both display and body, reducing font load

**Recommended**: **DM Sans** for the UI — it has enough personality to feel creative and fun (matching our app's mood) while remaining highly legible. It's geometric and friendly, pairing well with the playful nature of meme creation. For the meme canvas, offer Impact as default with Anton, Bebas Neue, Bangers, and Permanent Marker as alternatives.

## Visual Style Direction

**Recommended style**: Dark-mode-first creative workspace with vibrant accent pops

- **Layout**: Canvas-centric — the meme image occupies the majority of screen real estate. Tools live in a compact sidebar (left) and a contextual property bar (right or floating).
- **Background**: Deep charcoal/navy (`#1A1A2E`) — not pure black, which feels too harsh. Slight blue undertone adds depth.
- **Surfaces**: Slightly lighter dark (`#16213E` or `#252547`) for panels/cards — creating subtle layering without hard borders.
- **Corners**: Rounded (8-12px for cards/panels, 6-8px for buttons/inputs) — rounded corners feel friendlier and more modern, matching the playful mood.
- **Shadows**: Minimal — in dark mode, shadows are less effective. Use subtle borders or luminance differences for elevation instead.
- **Animations**: Subtle and functional — smooth transitions on panel open/close, gentle scale on hover for buttons, spring animation on text element selection. Nothing gratuitous.
- **Icons**: Outline style (Lucide or Phosphor) — clean, modern, consistent. Filled icons for active/selected states.
- **Canvas Area**: Slightly lighter than surrounding UI (maybe a subtle checkerboard pattern for transparent areas, like Photoshop) to clearly delineate the editing zone.
- **Text on Dark**: White (`#F8F8F2`) for primary text, muted lavender-gray (`#A0A0B8`) for secondary text — avoids the harshness of pure white while maintaining readability.

### Micro-interactions
- Drag handles glow with the primary accent color when grabbed
- Text elements show a subtle magenta border when selected
- Export button has a satisfying pulse/shimmer when the meme is ready
- Tool icons subtly brighten on hover

### Key UI Patterns for Meme Creator
- **Floating toolbar**: Appears near selected text element with quick-access styling options (bold, color, size)
- **Drag-to-place**: Click canvas to place text, drag to reposition — no coordinate inputs
- **Style presets**: One-click text styles ("Classic Meme", "Modern Clean", "Handwritten", "Neon Glow")
- **Quick export**: Single button, sensible defaults (PNG, original resolution), with expandable options

---
Status: READY_FOR_PLANNING
