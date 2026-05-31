# Research: Designer Todo App

**Date:** 2026-03-21
**Phase:** Research
**Status:** Complete

---

## 1. Executive Summary

A todo app for designers must not only manage tasks but *demonstrate design excellence* — its audience will judge it by its own aesthetics. This research covers the competitive landscape, designer-specific workflow needs, and the design patterns that separate premium productivity tools from generic ones.

**Key Finding:** The gap between generic todo apps (Todoist, Any.do) and designer-premium tools (Linear, Things 3, Craft) is enormous. Generic apps treat tasks as text items; premium apps treat them as *experiences* with rich metadata, fluid animations, and thoughtful hierarchy.

---

## 2. Competitive Landscape

### 2.1 Generic Todo Apps (Baseline)

| App | Strengths | Weaknesses for Designers |
|-----|----------|--------------------------|
| **Todoist** | Cross-platform, natural language input, integrations | Generic visual design, no visual project context, boring animations |
| **Any.do** | Clean mobile-first design | Surface-level, no depth, basic checkbox interactions |
| **TickTick** | Dense feature set | Visual noise, overwhelming for creative work |
| **Asana** | Project management depth | Too corporate, visual clutter, slow, made for teams not solo creatives |

**Bottom Line:** These apps treat tasks as text strings. Designers need richer context.

### 2.2 Designer-Premium Tools (Benchmarks)

#### Linear (linear.app)
The gold standard. Every pixel considered.

- **Visual System:** Dark-first with warm tones (#0D0D0F background), hairline borders instead of shadows, text hierarchy via opacity (100%/60%/40%), semantic muted status colors
- **Interaction Model:** Command palette (Cmd+K) as primary interaction — keyboard-first, always available
- **Animation:** Spring-based micro-interactions on every state change. Checkbox completion: fill → checkmark stroke draw → strikethrough animation chain
- **Priority:** Left-side colored bars (3px), not icon badges
- **Status Colors:** All at reduced saturation — never neon. Red for blockers, amber for in-progress, green for done
- **Why Designers Love It:** It *feels* like a design tool — fast, keyboard-driven, visually refined, GitHub-inspired issue model
- **Reference:** [Linear Design Blog](https://linear.app/design)

#### Things 3 (culturedcode.com)
Apple Design Award Winner (2017). The hierarchical standard.

- **Architecture:** Areas → Projects → Tasks → Checklists. Tags cross-cut all levels
- **Signature Pattern:** "Someday/Maybe" as a first-class top-level navigation item — this is critical for creative work where many ideas exist but aren't committed
- **Visual:** Pastel tag colors at low saturation, circular checkboxes with spring animations, 16px indent per hierarchy level
- **Animation:** Core Animation-driven, "outstanding use of animation for smooth interactions" (Apple citation)
- **Natural Language:** "tomorrow," "next friday," "in 3 weeks" — date parsing without calendar pickers
- **Why Designers Love It:** Clean Apple-native aesthetic, powerful hierarchy without complexity, Apple Pencil support
- **Reference:** [Things 3 on App Store](https://apps.apple.com/app/things-3/id904244726)

#### Craft (craft.do)
Document-first, content-centric.

- **Philosophy:** Content is king. UI recedes to let content shine
- **Signature Patterns:** Block-based documents (every element is a first-class block), callout blocks with colored left borders, floating drag handles on hover
- **Visual:** Narrow sidebar (240px) with lighter background, content area "pops" forward
- **Animation:** Toggle blocks with smooth chevron rotation (90°), 250ms ease-out slide
- **Why Designers Love It:** It respects visual hierarchy, supports page-level backgrounds, feels like a design tool for thoughts
- **Reference:** [Craft Documentation](https://www.craft.do)

#### OmniFocus (omniplan.com)
Power-user depth. BFD (Big Fanatic Deal) audience.

- **Architecture:** Perspective-based views, forecast perspective, review perspective
- **Strength:** Contexts (location-based or tag-based) for GTD methodology
- **Visual:** Functional rather than flashy, Apple-native, sidebar-based
- **Why Designers Use It:** Serious depth for complex workflows, Apple-only premium feel

### 2.3 Niche Competitors

| App | Niche | Relevance |
|-----|-------|-----------|
| **Sorted 3** | Timeline view with card-based display | Interesting date picker design |
| **2Do** | Extremely customizable | Overwhelming for most users |
| **Ulysses** (2017 ADA winner) | Distraction-free writing | Focus mode inspiration |
| **Agenda** (2018 ADA winner) | Timeline approach, notes linked to dates | Date-task relationship inspiration |

---

## 3. Designer-Specific Workflow Pain Points

### 3.1 Why Designers Struggle with Generic Todo Apps

1. **Visual work needs visual context.** A task "Design hero section" has no meaning without: color palette reference, typography choice, component mockup, or design file link. Generic apps store only text.

2. **Creative work has non-linear timelines.** "Someday/Maybe" is not a buried feature — it's a core workflow. Designers have dozens of ideas at various stages of commitment. Things 3 pioneered this as a first-class view.

3. **Iteration ≠ completion.** A design task goes through: ideation → wireframe → mockup → review → revision → approval. Generic todo apps have binary complete/incomplete. Designers need stages.

4. **Context switching is expensive.** Designers need to see *why* a task exists — the brief, the reference images, the brand guidelines. Generic apps strip all context.

5. **Projects vs. tasks blur.** "Design landing page" is both a project (with many sub-tasks) and a task (one item on the designer's plate). Generic apps force a choice.

6. **Due dates are anti-creative.** Forced due dates add anxiety without value for creative exploration. Designers need "no date" as a first-class state.

7. **Tags and folders are insufficient.** Creative work cuts across many dimensions: client, project type, urgency, phase. Flat tagging doesn't capture this richness.

### 3.2 Designer Metadata That Generic Apps Don't Support

- **Color swatches** — attach a color palette to a task/project
- **Reference images** — attach screenshots, mood board images
- **Font specimens** — attach typography choices
- **File links** — Figma, Sketch, Adobe XD file references
- **Iteration count** — "v3" of a design
- **Critique status** — needs review, in review, revisions requested, approved
- **Time estimate** — not for billing, but for creative planning ("this needs 2 hours")
- **Mood/energy tag** — "deep work," "quick win," "collaborative"

### 3.3 The Design Tool Parallel

Designers use tools like Figma, Sketch, and Framer because they feel *made for creative work*. The gap between Figma and PowerPoint is not features — it's *feel*. The same gap exists between Linear and Todoist. Our app must close this gap for task management.

---

## 4. Design Excellence Patterns

### 4.1 Animation Philosophy

Every animation must communicate something (Apple Design Award standard):

| Animation | Meaning |
|-----------|---------|
| Checkbox spring bounce | "Task completed, satisfaction delivered" |
| Task slide-out + collapse | "Task removed from your plate" |
| Drag lift + shadow | "You're holding this task, place it" |
| Staggered list entrance | "Here's the full picture, items arrive in order" |
| Skeleton shimmer | "Loading, layout preserved" |

**Timing Reference:**
- Micro (< 150ms): checkbox, button press, hover
- Transitions (150-300ms): panel open, list add/remove
- Page changes (250-400ms): view switches, modals
- Stagger offset: 30-50ms between list items

**Curves:**
- Spring (`cubic-bezier(0.34, 1.56, 0.64, 1)`) for interactive elements
- Ease-out (`cubic-bezier(0, 0, 0.2, 1)`) for entrances
- Ease-in-out (`cubic-bezier(0.4, 0, 0.2, 1)`) for page transitions

### 4.2 Typography System

**Primary:** Inter — tested at scale, excellent for UI, tabular figures
**Monospace:** JetBrains Mono — for dates, shortcuts, metadata
**Base size:** 13px for task list items (smaller than typical web apps)

```
Scale: 11px (labels) / 12px (secondary) / 13px (body) / 14px (emphasis) /
       16px (headings) / 20px (page titles) / 24px (display)
All on 4px grid. Letter-spacing: -0.01em for larger text, 0.01em for small.
```

### 4.3 Color System

**60-30-10 Rule:**
- 60% neutral/background
- 30% secondary surfaces
- 10% accent/interactive

**Dark Theme Palette (Primary):**
```
Background:    #0F0F11 (deepest) → #1A1A1E → #242428 → #2E2E34 (elevated)
Surface:       #1E1E24
Border:        rgba(255,255,255,0.08)
Text:          #F4F4F5 (primary) / #A1A1AA (secondary) / #71717A (tertiary)
Accent:        #8B5CF6 (violet) — refined, not garish
```

**Semantic Colors (all at 50-60% saturation):**
- Success: #22C55E
- Warning: #F59E0B
- Error: #EF4444
- Info: #3B82F6

**Status Colors (muted):**
- Todo: Neutral gray
- In Progress: Amber at 30% saturation
- Review: Blue at 30% saturation
- Done: Green at 30% saturation
- Blocked: Red at 30% saturation

### 4.4 Layout Architecture

```
┌──────────────┬──────────────────────────────────┬───────────────┐
│   SIDEBAR     │         MAIN CONTENT             │  DETAIL PANEL │
│   (240px)     │         (flexible)               │   (360px)     │
│              │                                  │              │
│  - Inbox      │  ┌─────────────────────────────┐ │  (appears on  │
│  - Today      │  │ View Header + Filter Bar   │ │   task select)│
│  - Upcoming   │  ├─────────────────────────────┤ │              │
│  - Someday    │  │                             │ │  - Title      ││
│  - Projects   │  │     Task List              │ │  - Status     ││
│    └ proj1    │  │     (48px rows, scrollable) │ │  - Due date   ││
│    └ proj2    │  │                             │ │  - Tags       ││
│  - Tags       │  │                             │ │  - Notes      ││
│              │  │                             │ │  - References ││
│              │  └─────────────────────────────┘ │              │
│  [+ Project] │  [  Quick add task...         ] │              │
└──────────────┴──────────────────────────────────┴───────────────┘
```

**Density Tiers:**
- Sidebar: 8px vertical padding, 40px row height
- Task list: 12px vertical padding, 48px row height
- Detail panel: 20px padding, 60px section gaps
- Empty states: 40-60px margins, centered

### 4.5 Information Architecture

**Views (time-based):**
1. **Inbox** — uncategorized capture, the daily review entry point
2. **Today** — committed tasks for today, highlighted with accent at 5%
3. **Upcoming** — tasks with future dates, grouped by week
4. **Someday/Maybe** — low-pressure backlog, visually de-emphasized
5. **All Tasks** — complete list with search

**Organization (structure-based):**
- **Projects** — max 2 levels deep (project + sub-project)
- **Tags** — cross-cutting, color-coded, max 10 distinct colors
- **Priority** — 3 levels (high/medium/low) with left-border indicators

**Task Lifecycle:**
```
Inbox → Project/Today/Someday → In Progress → Review → Done
                                    ↓
                                Blocked
```

---

## 5. Key Differentiators for a Designer Todo App

Based on competitive analysis and designer pain points, the following features would differentiate this app:

### 5.1 Must-Have (MVP)

1. **Rich task metadata** — color swatch attachment, reference image links, file links (Figma URL)
2. **Someday/Maybe as first-class view** — not buried, visible in sidebar
3. **Fluid drag-and-drop reordering** — @dnd-kit + Framer Motion, scale/shadow lift effect
4. **Command palette** — Cmd+K for quick actions (add, search, navigate)
5. **Dark-first design** — designers predominantly use dark UIs
6. **Keyboard-first** — every action accessible without mouse
7. **Optimistic UI** — immediate feedback on all actions
8. **Beautiful empty states** — each view has a crafted empty state with icon and CTA

### 5.2 Should-Have (V2)

1. **Visual iteration stages** — "ideation → wireframe → mockup → review → revision → approved" pipeline
2. **Tag color system** — 8-10 pastel colors, visually distinct but harmonious
3. **Natural language date input** — "tomorrow," "next monday," "in 2 weeks"
4. **Focus mode** — hides sidebar/chrome, dims everything except current task
5. **Recurring tasks** — with next-occurrence preview
6. **Project progress indicator** — circular or linear progress on project headers

### 5.3 Nice-to-Have (Future)

1. **Color palette attachment** — attach a hex palette to a project
2. **Figma/Framer embed preview** — inline preview of linked design files
3. **Critique workflow** — mark tasks as "needs review," "changes requested," "approved"
4. **Time blocking** — estimate time per task, see workload visualization
5. **Mood/energy tagging** — "deep work" vs. "quick win" filter

---

## 6. Technology Stack Alignment

The existing template already selects: **Next.js 14 + TypeScript + Tailwind + shadcn/ui + Framer Motion + @dnd-kit**

This is the correct choice:
- **shadcn/ui** provides component primitives that can be fully restyled — start from proven components, customize to designer-premium
- **Tailwind** enables the custom color system, spacing grid, and dark mode
- **Framer Motion** handles spring animations, AnimatePresence, and layoutId for drag-reorder
- **@dnd-kit** provides accessible, framework-agnostic drag-and-drop
- **localStorage** for MVP persistence — designers expect instant, local-first tools

**Critical customizations needed:**
- Override shadcn/ui theme with bespoke dark palette (#0F0F11 base)
- Replace default Inter with a customized weight/styling
- Add JetBrains Mono for metadata
- Build custom checkbox component with Framer Motion animation chain
- Build custom command palette (Cmd+K)
- Implement `prefers-reduced-motion` support

---

## 7. References

- Linear Design Blog: https://linear.app/design
- Things 3 (App Store): https://apps.apple.com/app/things-3/id904244726
- Craft: https://www.craft.do
- Apple Human Interface Guidelines: https://developer.apple.com/design/human-interface-guidelines
- Framer Motion Documentation: https://www.framer.com/motion/
- dnd-kit: https://dndkit.com

---

## 8. Summary

The opportunity is clear: there is no designer-premium todo app that combines the visual excellence of Linear/Craft with the hierarchical power of Things 3, specifically tailored for creative work. Generic apps are too flat; existing premium apps are either too corporate (Linear for dev teams) or too platform-locked (Things 3 for Apple).

A designer todo app should:
1. **Feel like a design tool itself** — every pixel considered, every animation purposeful
2. **Support creative lifecycle** — Someday/Maybe, iteration stages, visual references
3. **Be keyboard-first and fast** — Cmd+K command palette, shortcuts everywhere
4. **Use restrained, premium aesthetics** — warm dark palette, spring animations, generous whitespace
5. **Implement the hybrid architecture** — Projects (hierarchy) + Tags (cross-cut) + Views (time-based)

**PHASE_COMPLETE: research**
