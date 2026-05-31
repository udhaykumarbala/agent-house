# UX Research: Random Color Palette Generator

## Pattern Analysis

Referenced apps and their successful patterns:

### Coolors (coolors.co)
- **Pattern Used**: Full-screen horizontal color swatches with spacebar to generate new palettes. Each swatch occupies an equal vertical strip across the viewport. Hex codes centered on each swatch with lock, copy, drag-reorder, and adjust icons appearing on hover.
- **Why It Works**: The spacebar shortcut creates an addictive "slot machine" feel — users can rapidly iterate through palettes with zero friction. Full-screen swatches give an accurate sense of how colors feel together at scale, not just as tiny chips. The lock mechanism lets users pin colors they like and randomize the rest, enabling convergent exploration.
- **What We Can Learn**: Make generation instant and repeatable with a single keypress. Show colors large — small swatches don't convey palette harmony. Let users lock individual colors to refine incrementally rather than starting from scratch each time.

### Adobe Color (color.adobe.com)
- **Pattern Used**: Color wheel as the primary interaction surface with harmony rules (complementary, analogous, triadic, split-complementary, etc.). Users drag points on the wheel; dependent colors update in real-time. A strip of five color swatches below the wheel shows the active palette.
- **Why It Works**: The color wheel provides a mental model users already understand from art education. Harmony rules give structure to randomness — users learn *why* certain colors work together. The wheel + swatch strip layout separates exploration (wheel) from evaluation (swatches).
- **What We Can Learn**: Offer harmony-based generation modes so palettes aren't just random — they're aesthetically grounded. Even in a "random" tool, giving users optional structure (e.g., "generate analogous palette") dramatically improves output quality.

### Colormind (colormind.io)
- **Pattern Used**: Five horizontal swatches in a single row with a prominent "Generate" button. Each swatch shows hex value and can be locked. Minimal UI — almost no chrome, just colors and controls. Uses deep learning to generate palettes from real-world color data.
- **Why It Works**: The extreme simplicity lowers the barrier to entry — no learning curve at all. Five colors is the sweet spot: enough for a usable design palette (primary, secondary, accent, background, text) without overwhelming the user. The AI-driven approach produces palettes that feel "designed" rather than random.
- **What We Can Learn**: Default to 5 colors — it maps to real design needs. Keep the interface minimal; the colors *are* the interface. A single prominent action button is all that's needed.

### Palette Generator by Canva (canva.com/colors/color-palette-generator)
- **Pattern Used**: Upload an image to extract a palette, plus a curated gallery of trending palettes. Each palette card shows 4 colors with hex values and a copy button. Clean card-based grid layout.
- **Why It Works**: Image extraction solves the "blank canvas" problem — users who don't know color theory can start from inspiration they already have. The curated gallery provides social proof and starting points. Card layout makes palettes scannable and comparable.
- **What We Can Learn**: Consider providing palette "seeds" or starting points, not just pure randomness. Showing palettes as cards enables visual comparison when users want to browse.

### Realtime Colors (realtimecolors.com)
- **Pattern Used**: Live preview of a palette applied to a realistic website mockup. Users adjust colors and immediately see how they'd look in context (buttons, text, backgrounds, cards). Exports to CSS variables, Tailwind, and other formats.
- **Why It Works**: It answers the question users actually have: "Will these colors work in my project?" Seeing colors in context is fundamentally different from seeing them in isolation — a palette that looks good as swatches may fail when applied to real UI elements.
- **What We Can Learn**: Even a simple preview showing how colors map to roles (background, text, accent, etc.) adds enormous value. Export in useful formats (CSS, Tailwind, etc.) bridges the gap from inspiration to implementation.

## User Journey Analysis

### Entry Point
Users arrive with one of three mindsets:
1. **"I need colors for a project"** — Designer or developer starting a new project, needs a cohesive palette quickly
2. **"I'm exploring for inspiration"** — Browsing palettes without a specific goal, looking for something that sparks interest
3. **"I need to match/extend existing colors"** — Has one or two colors already, needs complementary colors to complete a palette

### Core Loop
1. **Generate** → View a new random palette
2. **Evaluate** → Quickly judge if the palette feels right (this takes < 1 second)
3. **Refine** → Lock good colors, regenerate the rest OR adjust individual colors
4. **Repeat** → Steps 1-3 happen rapidly, often 10-30 times per session

The core loop must be *fast*. Any friction in the generate→evaluate cycle breaks the flow state.

### Success State
- User copies hex/RGB values to clipboard (the ultimate intent signal)
- User exports the full palette in a usable format
- Visual confirmation: a brief toast/animation confirming the copy action

## Mental Model

Users think of this task as: **"shuffling a deck of cards until I get a good hand"** — rapid iteration with occasional lucky finds. They expect the same instant gratification as pulling a slot machine lever.

Common terminology:
- "Palette" (not "color scheme" or "color set")
- "Generate" or "randomize" (not "create" or "build")
- "Lock" a color (not "pin" or "save")
- "Copy" the hex code
- Hex codes (#FF5733) are the universal language; RGB secondary; HSL for power users

## Anti-Patterns to Avoid

- **Don't**: Require registration or sign-in to use basic features — **Why**: Color palette generators are quick-use tools. Forcing sign-up before the user sees value causes immediate bounce. Coolors gets this right: full functionality without an account.

- **Don't**: Show colors only as small chips or tiny squares — **Why**: Small color samples don't convey how colors feel at scale. A hex value next to a 20px square tells you nothing about whether that color works as a background. Users need to *feel* the color, which requires visual weight.

- **Don't**: Generate purely random RGB values with no harmony constraints — **Why**: Truly random RGB produces clashing, muddy, or washed-out palettes 90% of the time. Users will blame the tool, not the math. Even "random" should be weighted toward pleasing combinations (e.g., using HSL with constrained saturation/lightness ranges).

- **Don't**: Hide the copy-to-clipboard action behind menus or extra clicks — **Why**: Copying a hex code is the #1 action users take after finding a color they like. It should be a single click/tap directly on or next to the color value. Every extra click loses users.

- **Don't**: Auto-play animations or transitions that delay seeing the new palette — **Why**: Users generate palettes rapidly (spacebar-spacebar-spacebar). Any animation that takes > 200ms between generations breaks the flow. Transitions should be instant or near-instant.

- **Don't**: Use a dark UI that competes with the palette colors — **Why**: The palette is the content. A busy or strongly-colored UI frame makes it harder to evaluate the actual palette. Use neutral, minimal chrome so the generated colors are the hero.

## Recommended Patterns for Our App

Based on research, we should use:

1. **Full-width horizontal color swatches (Coolors-style layout)** — because large color areas are essential for evaluating palette harmony, and this layout is the established convention users expect from palette generators. Each color gets equal visual weight.

2. **Single-key generation (spacebar)** — because the core loop must be frictionless. A spacebar trigger is the established convention (Coolors trained this behavior). Also provide a visible "Generate" button for discoverability and mobile users.

3. **Per-color lock toggle** — because users rarely love all 5 colors at once. Locking lets them iteratively converge on a perfect palette by keeping winners and re-rolling the rest. This is the single most important refinement feature.

4. **One-click copy with visual feedback** — because copying hex values is the primary success action. Clicking a color's hex code should copy it to clipboard and show a brief "Copied!" toast or checkmark. Also provide a "Copy All" option for the full palette.

5. **Default 5-color palette with smart randomization** — because 5 maps to real design needs (primary, secondary, accent, background, text). Generation should use HSL-based logic with constrained saturation and lightness to avoid muddy/clashing results, rather than pure RGB randomness.

6. **Minimal, neutral UI frame** — because the colors are the product. Use light gray or near-white chrome, minimal borders, and let the swatches dominate the viewport. No competing brand colors or heavy navigation.

7. **Multiple color format display (hex + RGB + HSL)** — because different users need different formats. Show hex by default (most common) with a toggle or secondary display for RGB and HSL values.

8. **Keyboard shortcuts for power users** — Spacebar to generate, L to lock hovered color, C to copy hovered color. These map to the rapid-fire workflow power users expect.

---
Status: READY_FOR_PLANNING
