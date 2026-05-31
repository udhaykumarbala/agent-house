# UI Research: DesignFlow — A Todo App for UX Designers

## Design Inspiration

### 1. Things 3 (Cultured Code)
- **Visual Style**: Warm, editorial with a cream/off-white background and a single vivid accent color (teal). Sparse, confident layout with generous whitespace.
- **What Makes It Great**: The checklist never competes with content — it breathes. The sidebar is ultra-minimal; the main area is uncluttered. Tag colors are muted pastels. It feels like a premium paper planner digitized.
- **What We Can Borrow**: Warm neutral base palette, single bold accent for CTAs, pastel category chips, the concept of "Areas" (grouped projects) as a nav structure.

### 2. Craft
- **Visual Style**: Rich, document-based with subtle textures, warm tones, and a parchment-like feel. Cards use soft shadows and warm borders.
- **What Makes It Great**: It treats every document as beautiful by default. The block-based editor metaphor extends to the UI. Dark mode feels like a cozy workspace, not a terminal.
- **What We Can Borrow**: The warm texture/cream palette, soft border treatment, the idea that even a simple todo feels premium when it has visual polish.

### 3. Linear
- **Visual Style**: Dark, precise, with vibrant accent colors (purple, blue, red) against near-black backgrounds. Sharp corners, tight spacing, micro-animations everywhere.
- **What Makes It Great**: Keyboard-first navigation, real-time sync feel, status indicators that are instantly scannable. The dark theme is the brand — not a toggle.
- **What We Can Borrow**: The status badge system (priority pills, iteration markers), keyboard shortcuts UX, smooth state transitions, the dark theme as a first-class option.

### 4. Superhuman
- **Visual Style**: Almost entirely monochrome — white or dark background, black text, single red accent for the most important actions. Extreme minimalism.
- **What Makes It Great**: Zero visual noise. Every pixel either does work or breathes. The command bar (/) is the entire product metaphor.
- **What We Can Borrow**: The command palette (`/`) as an input paradigm, extreme restraint in visual decoration, monospaced details for metadata.

### 5. Figma
- **Visual Style**: Neutral gray workspace with a vibrant purple (#A259FF) as the singular brand accent. Toolbar-dominant UI, minimal chrome.
- **What Makes It Great**: The canvas is the star — UI stays out of the way. Contextual panels appear when needed. The purple accent is consistent across the entire product.
- **What We Can Borrow**: The purple accent strategy, the toolbar/nav pattern where the sidebar collapses into icons, contextual right-click menus.

### 6. Raycast
- **Visual Style**: Floating command bar with translucent/blur backgrounds. Dark mode native, rounded pill UI, category icons in colored circles.
- **What Makes It Great**: The extension ecosystem is visually seamless. Search is instant and beautiful. The popover feel makes everything feel lightweight.
- **What We Can Borrow**: Blur/translucent surface treatment, floating command palette aesthetic, icon-forward category navigation.

## Competitor Color Analysis

| Todo App | Primary Color | Style | Notes |
|----------|-------------|-------|-------|
| Todoist | Red #E4433E | Action-oriented | AVOID — aggressive red feels urgent/stressful |
| Things 3 | Teal #528FF5 | Calm, focused | INTERESTING — muted teal on warm white |
| Linear | Purple #8B5CF6 | Precision tech | USED — but works well with dark mode |
| Notion | Charcoal #37352F | Neutral, calm | INTERESTING — near-black on white, minimal color |
| Things 3 | Warm cream palette | Editorial warmth | BEST IN CLASS for UX designers |
| Apple Reminders | Blue #007AFF | Apple ecosystem | AVOID — generic Apple blue |
| Todoist | Red #D93025 | Urgency | AVOID — stress-inducing for a daily tool |
| Superhuman | Red #FF3B30 | Radical focus | INTERESTING — but only one red dot, not full palette |

**Key Insight**: Most todo apps use blue (trust/focus) or red (urgency). Things 3's warm cream + teal is the most distinctive and beloved by designers. Linear's dark + purple shows how a single vibrant accent can carry a brand.

## Our Differentiation Strategy

To stand out, we should:
- **Avoid**: Blue primary (Todoist, Apple, generic SaaS), Red accents (stress-inducing for a daily-use creative tool), Purple-only (Linear already owns this for dark mode)
- **Consider**:
  - **Warm cream + deep indigo** — Editorial warmth meets sophisticated depth (evokes design tools + print magazines)
  - **Slate + coral/salmon** — Warm neutral dark mode with an unexpected warm accent (feels like a late-night creative session)
  - **Cream + forest green** — Grounded, natural, calm — like a designer's sketchbook (risky if too muted)

## Color Psychology

**App Mood**: "Calm clarity for creative work" — UX designers need focus without rigidity, structure without sterility. The tool should feel like a well-organized design studio: everything in its place, beautiful to look at, never distracting.

### Recommended color directions:

1. **Warm Cream + Deep Indigo (RECOMMENDED)**
   - Background: Warm Cream #FAF8F5
   - Primary: Deep Indigo #3730A3 (rich, sophisticated, not the typical blue)
   - Accent: Coral #F97316 for highlights/actions
   - Surface: Pure White #FFFFFF with warm shadow
   - **Feeling**: Editorial magazine meets modern design tool — confident, calm, premium
   - **Why**: Indigo reads as "professional" but is warmer than navy. Coral adds energy without the stress of red. Cream grounds everything.

2. **Slate Dark + Peach**
   - Background: Deep Slate #0F172A
   - Primary: Peach #FB923C
   - Surface: Slate #1E293B
   - **Feeling**: Late-night creative studio, warm glow
   - **Why**: Dark mode native with warmth — most dark tools feel cold. Designers who work late will appreciate the warmth.

3. **Warm Stone + Violet**
   - Background: Stone #F5F5F4
   - Primary: Deep Violet #7C3AED
   - Accent: Amber #F59E0B
   - **Feeling**: Grounded and creative
   - **Why**: Violet as primary (different from Linear's purple) with amber warmth — feels creative without being chaotic.

## Typography Research

**Display/Heading Font Options**:
- **Instrument Serif**: Elegant editorial serif — perfect for a designer's tool. Warm, readable, distinctive. Used by Figma in marketing materials.
- **Fraunces**: Soft serif with personality — optical sizing, warm curves. Very design-community approved.
- **Cabinet Grotesk**: Modern geometric sans with warmth — used extensively in the design tool space (Framer, Loom).
- **Geist**: Extremely clean, modern, excellent for UI. But very tech-forward — might feel cold for a creative audience.
- **Manrope**: Rounded geometric sans — friendly but professional. Good middle ground.

**Body Font Options**:
- **Inter**: Industry standard, excellent legibility at all sizes. Designers may find it boring but it's never wrong.
- **Satoshi**: Modern geometric with warmth — distinctive without being quirky.
- **Plus Jakarta Sans**: Geometric but approachable — trending in the design community.
- **Geist**: If we go full dark mode, Geist pairs perfectly with it.

**Recommendation**:
- **Primary: Cabinet Grotesk** (headings) + **Inter** (body) — The pairing of a distinctive display font with the workhorse body font gives us personality without sacrificing readability. Both are variable fonts (modern).
- **Alternative**: **Instrument Serif** (headings) + **Inter** (body) — For a more editorial, warm feel. Better if we lean heavily into the cream/indigo editorial direction.
- **Dark mode option**: **Geist** (both) — Pairs perfectly with the slate dark theme.

**Why Cabinet Grotesk**: It's used by Figma-adjacent tools (Loom, Framer) and is specifically beloved in the design community. It signals "made by designers, for designers."

## Visual Style Direction

**Recommended style: "Editorial Design Studio"**
- Warm cream backgrounds with a rich indigo primary and coral accents
- Generous whitespace — let tasks breathe like editorial content
- Subtle warm shadows (cream-tinted, not gray)
- Rounded corners: medium (8-10px) — approachable, not childish
- Minimal borders, prefer shadow for depth
- Smooth micro-animations: 200-300ms ease-out transitions
- Icons: Outlined, 1.5px stroke weight, indigo-tinted

| Element | Recommendation | Rationale |
|---------|---------------|-----------|
| Corners | 8px standard, 12px for cards | Approachable, modern |
| Shadows | Warm-tinted, soft | Depth without coldness |
| Animations | Subtle spring/overshoot | Playful precision — designers notice |
| Borders | Minimal, warm gray | Only when needed for grouping |
| Icons | Outlined, consistent stroke | Clean, scalable, modern |
| Dark Mode | Slate + Peach | Warm dark option, not cold |

### Category-Specific Color Accents

For UX Designer-specific categories, we need distinct but harmonious colors:
| Category | Color | Hex | Why |
|----------|-------|-----|-----|
| Research | Sage | #4D7C0F | Investigative, grounded |
| Design | Indigo | #3730A3 | Creative depth |
| Testing | Amber | #D97706 | Warm attention |
| Review | Rose | #BE185D | Critical but warm |
| Discovery | Teal | #0D9488 | Exploratory calm |
| Handoff | Slate | #475569 | Technical/developer handoff |

## UI Component Inspirations

### Command Palette (from Superhuman/Raycast)
- Floating overlay with backdrop blur
- Search with fuzzy matching
- Keyboard-first (`Cmd+K` to open)
- Results grouped by category with icons

### Task Cards (from Things 3 + Craft)
- Clean white card with warm shadow
- Category chip (colored pill) at top
- Task title prominent
- Metadata row: due date, priority, iteration tag
- Subtle checkbox animation on completion

### Kanban/Iteration View (from Linear)
- Horizontal scrolling swimlanes
- Each lane = a design phase (Exploration → Refinement → Validation → Handoff)
- Compact task cards with status dots
- Smooth drag-and-drop with haptic feel

### Progress/Stats (from Craft + Notion)
- Circular progress rings for iteration completion
- Weekly/monthly streak counters
- Clean stat cards with large numbers + small labels

---

Status: READY_FOR_PLANNING
