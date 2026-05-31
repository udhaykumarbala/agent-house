# UI Research: Designer Todo App

## Design Inspiration

### Linear
- **Visual Style**: Dark-first, monochromatic with electric purple accent. Ultra-minimal, dense information display without clutter.
- **What Makes It Great**: Keyboard-driven, blazing fast, information-dense yet calm. The UI disappears and lets you work. Dark mode by default feels like a premium developer tool.
- **What We Can Borrow**: Dark-first aesthetic, single accent color philosophy, keyboard shortcuts, subtle depth through layered surfaces, monospace metadata.

### Craft
- **Visual Style**: Warm, inviting light mode with amber/orange accent. Rich card design with visual hierarchy that feels like a premium document editor.
- **What Makes It Great**: The "documents as building blocks" philosophy — every task can expand into a rich page. Beautiful type hierarchy. Feels handcrafted rather than templated.
- **What We Can Borrow**: Warm color temperature (cream/paper backgrounds), premium typography feel, expandable task cards, amber accent on light backgrounds.

### Things 3
- **Visual Style**: MacOS-native feel with generous whitespace, soft shadows, and a green accent. Clean, calm, desktop-app quality.
- **What Makes It Great**: Pioneered the "calm productivity" aesthetic. Multiple "Areas" for organizing life, elegant project hierarchies, beautiful animations on task completion.
- **What We Can Borrow**: Generous whitespace philosophy, soft rounded corners (8-12px), green accent for completion states, subtle spring animations, area/project hierarchy.

### Sunsama
- **Visual Style**: Warm terracotta/coral accents on dark backgrounds. Feels like a "daily planning ritual" tool — calm, intentional, meditative.
- **What Makes It Great**: Daily focus planning, warm color temperature in dark mode, "channel" concept for organizing work by mode (creative vs communication). Popular among creative professionals.
- **What We Can Borrow**: Warm coral accent on dark backgrounds, daily planning focus, "mode" concept for creative vs execution work, calm celebration moments.

### Figma (FigJam)
- **Visual Style**: Canvas-based, freeform, sticker-like elements, infinite canvas with gentle dot grid.
- **What Makes It Great**: Visual thinking is first-class. Brainstorming feels natural, not constrained by list views. Mood boards and tasks coexist.
- **What We Can Borrow**: Visual reference embedding, canvas-like organization, dot grid backgrounds, freeform card placement option.

## Competitor Color Analysis

| Competitor | Primary Color | Hex | Style | Notes |
|------------|-------------|-----|-------|-------|
| Linear | Electric Purple | #8B5CF6 | Dark tech-minimal | Too common among dev tools |
| Craft | Warm Amber | #E8853A | Warm premium docs | Good for light mode only |
| Things 3 | Soft Green | #3ECF8E | Calm completion | Excellent for success states |
| Notion | Warm Red | #E03E3E | Warm neutral | Too corporate/generic |
| Sunsama | Coral Terracotta | #E07A5F | Warm dark mode | Excellent differentiated tone |
| Height | Warm Coral | #FF6B6B | Modern AI-native | Fresh but slightly loud |
| Figma | Dark Purple | #0D99FF (blue) | Brand-driven | Not a productivity accent |
| Monday | Multicolor | #FF3D00 | Overwhelming | AVOID - visual chaos |

## Our Differentiation Strategy

To stand out, we should:
- **Avoid**: Generic blue/purple gradients (too AI-app), enterprise gray/white (too corporate), multicolor (too Monday.com)
- **Consider**: **Deep teal + warm gold** — a jewel-tone direction that feels premium, unique, and calming. Or **charcoal + soft sage** for a nature-inspired calm aesthetic.

**Unique Direction**: **Deep Ink + Amber Gold** — A rich, almost noir dark mode with warm amber/gold accents. Feels like a luxury creative studio tool. Light mode: **Warm Parchment + Deep Amber**. This differentiates from the sea of purple/blue AI tools and the generic gray productivity apps.

## Color Psychology

**App Mood**: "The Designer's Studio at Night" — calm, focused, premium, intentional.

Designers experience:
- Creative block and decision fatigue
- Context switching between multiple projects
- Review-dependent work blocked on others
- Deep work vs communication mode switching

Recommended color directions:
1. **Deep Ink + Amber Gold** (Recommended): Rich near-black (#0E0E10) with warm amber (#E8A838) accents. Feels like a luxury studio tool. Warm amber = creative energy, focus, premium quality.
2. **Charcoal + Sage**: Deep gray (#1C1C1E) with soft sage green (#8AAF8A) accents. Nature-inspired calm, reduces anxiety. Good for stressed designers.
3. **Midnight + Coral**: Deep navy (#0F172A) with warm coral (#FF7B6B) accents. More energetic than #1, good for fast-paced studios.

## Typography Research

**Display Font Options**:
- **Plus Jakarta Sans**: Warm geometric sans with personality. Has a "designed by humans" feel rather than algorithmic. Excellent for headings. 700 weight for impact.
- **Satoshi**: Geometric with subtle humanist touches. Modern, slightly playful without being casual.
- **Outfit**: Clean geometric, slightly rounded. Very versatile. Excellent readability.
- **General Sans**: Premium, by Commercial Type. Distinctive but readable.

**Body Font Options**:
- **Inter** (safe choice): Ubiquitous but excellent. Consistent, readable, variable font support.
- **Plus Jakarta Sans**: Can work for both display and body if weight is managed.
- **Geist**: Modern, technical feel. Good if going for "dev-designer" hybrid tool.

**Recommended**: **Plus Jakarta Sans** for both headings and body. It has enough personality to signal "designer tool" while remaining highly readable. Use weights 400 (body), 500 (medium emphasis), 700 (headings). Monospace: **JetBrains Mono** for metadata (dates, tags, IDs).

## Visual Style Direction

Recommended style: **"Noir Studio" — Dark-first with warm amber accents**
- Corners: 8px for cards/buttons (rounded but not bubbly), 12px for modals, 4px for inputs
- Shadows: Soft, layered — `0 1px 3px rgba(0,0,0,0.3), 0 4px 12px rgba(0,0,0,0.2)` for elevated elements
- Animations: Spring-based with `cubic-bezier(0.34, 1.56, 0.64, 1)` for bouncy completion moments. 200-300ms duration for micro-interactions.
- Background: Layered darks — base #0E0E10, elevated surfaces #161618, highest surfaces #1E1E21
- Accent: Amber #E8A838 with glow effect on focus states

**What NOT to do**:
- Avoid Tailwind blue-500 (#3B82F6) — the single most overused color in AI/generic apps
- Avoid flat gray backgrounds (#1F2937) — too GitHub Copilot
- Avoid purple gradients — too generic "AI assistant"
- Avoid green success states that feel like a healthcare app
- Avoid "enterprise" feeling with hard edges and corporate blue

## Component Aesthetic Targets

### Task Cards
- Subtle left border accent in amber for priority/highlight
- Soft inner shadow for depth
- Hover: slight elevation + amber glow on border
- Completion: Satisfying checkmark animation with spring bounce

### Buttons
- Primary: Amber fill, dark text, subtle glow on hover
- Secondary: Transparent with amber border, amber text
- Ghost: Text only with underline animation on hover

### Input Fields
- Dark surface (#1E1E21) with subtle border
- Focus: Amber border with soft amber glow (`box-shadow: 0 0 0 3px rgba(232,168,56,0.15)`)
- Monospace placeholder text in muted gray

### Navigation
- Minimal icon + text sidebar
- Active state: Amber left border + text highlight
- Subtle hover states without heavy backgrounds

---
Status: READY_FOR_PLANNING
