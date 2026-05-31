# UX Research: Todo List for UX Designers

**Version:** 1.0
**Date:** 2026-03-21
**Status:** READY_FOR_DESIGN
**Researcher:** Security Expert Agent (Research Phase)

---

## Executive Summary

This research examines the specific needs, workflows, and pain points of UX designers when it comes to task management. It analyzes competitive tools, identifies gaps in generic todo applications, and establishes design principles for building a todo app that UX designers would actually choose over generic alternatives. The core insight: **UX designers need task management that mirrors the double diamond design process — not a generic inbox with checkboxes.**

---

## 1. UX Designer Workflows & Processes

### 1.1 The UX Design Process Stages

UX designers work through a non-linear, iterative process typically described by frameworks like the Double Diamond (Discover → Define → Develop → Deliver) or similar variants:

| Stage | Activities | Task Nature |
|-------|-----------|-------------|
| **Discover / Research** | User interviews, competitive analysis, diary studies, survey design | Research-heavy, async, reference-gathering |
| **Define / Synthesize** | Affinity mapping, persona creation, journey mapping, problem statements | Synthesis, grouping, prioritization |
| **Develop / Ideate** | Sketching, wireframing, concept development, design critiques | Iteration-heavy, visual, collaborative |
| **Deliver / Prototype** | High-fidelity mockups, prototyping, usability testing prep, developer handoff | Output-focused, deadline-driven |

### 1.2 What Makes UX Designer Work Different

Unlike software development tasks (which are relatively discrete and linear), UX design work has unique characteristics:

1. **Iteration is inherent**: A "wireframe v1" task doesn't end cleanly — it flows into v2, v3, critique feedback incorporation, and so on.
2. **Soft deadlines**: Research synthesis doesn't have a clear "done" state — it's complete when insights feel sufficiently developed.
3. **Context-heavy**: A task like "Test navigation flow" requires attaching screenshots, prototype links, participant notes, and findings.
4. **Collaborative by default**: Design critiques, stakeholder reviews, and developer syncs mean tasks frequently involve others.
5. **Mixed granularity**: A designer might have a 2-month project-level goal ("Redesign checkout flow") alongside micro-tasks ("Export assets for developer").

### 1.3 Designer Task Categories

Research consistently surfaces these recurring task categories for UX designers:

- **Research Tasks**: "Conduct 5 user interviews", "Compile affinity map", "Write research synthesis doc"
- **Design Iteration Tasks**: "Sketch 3 concepts for login flow", "Create high-fidelity mockups", "Incorporate critique feedback"
- **Usability Testing Tasks**: "Recruit 5 participants", "Run moderated test session", "Analyze findings", "Write test report"
- **Stakeholder/Process Tasks**: "Present to product team", "Handoff to engineering", "Update design system"
- **Personal/Admin Tasks**: "Review competitor apps", "Update portfolio", "Schedule critique"

### 1.4 Daily/Weekly Designer Workflow Patterns

Based on research from UX team workflows (NN/g, Nielsen Norman Group research on creative professionals):

**Daily Core Loop:**
1. Review current iteration status — "Where am I in the design cycle?"
2. Check for critique feedback or stakeholder input
3. Work on active design task (typically 1-2 major tasks per day)
4. Document progress (update Figma, write notes)
5. Plan next day's priority

**Weekly Core Loop:**
1. Week kickoff: What did I accomplish? What's the focus this week?
2. Mid-week check: Am I on track for the project milestone?
3. Design critique participation (typically 1-2x per week)
4. Stakeholder syncs (review, feedback, alignment)
5. Week retrospective: What slowed me down? What needs to change?

### 1.5 Pain Points with Generic Todo Apps

Through analysis of designer forums (Designer News, UX Collective, Reddit r/userexperience), several recurring complaints emerge about using generic tools like Todoist, Things 3, or Apple Reminders for design work:

| Pain Point | Description | Frequency |
|-----------|-------------|-----------|
| **No iteration support** | Tasks are binary (done/not done), but design work requires "in progress, needs review, approved" states | Very High |
| **No attachment context** | Can't attach Figma files, prototype links, or research docs to tasks | Very High |
| **No design-specific categories** | "Work" and "Personal" are meaningless — designers need "Research", "Design", "Testing", "Review" | High |
| **No estimation/tracking** | Designers estimate in hours or story points, but generic apps don't support this | Medium-High |
| **No recurring design rituals** | "Weekly design critique" or "Bi-weekly stakeholder review" have specific recurrence needs | Medium |
| **Flat hierarchy** | Projects are nested (epics → features → tasks), but flat task lists can't represent this | Medium |
| **No time blocking** | Designers block time for deep work (2hr design sprint), but calendar integration is separate | Medium |
| **Ugly interface** | Designers notice bad design — if the tool looks generic, they distrust its utility | High |

---

## 2. Competitive Analysis

### 2.1 Tools Designers Currently Use (and Their Gaps)

#### **Notion**
- **Strengths**: Flexible databases, linked pages, rich text, excellent for documentation, template ecosystem
- **Weaknesses for UX Designers**: Over-engineered for simple tasks, slow, database views are complex to set up per-project, not designed for quick task capture
- **Designer Use Case**: Often used as a project wiki + task tracker hybrid, but the task part is always secondary
- **Verdict**: Too heavy; great for documentation, poor for lightweight daily task management

#### **Linear**
- **Strengths**: Beautiful UI, keyboard-first, fast, project/issue hierarchy, status workflow (Backlog → Todo → In Progress → Done)
- **Weaknesses for UX Designers**: Built for engineering teams; no design-specific task types, no research/test tracking built in, requires project/team context
- **Designer Use Case**: Some design teams use Linear, but it's a stretch — Linear's opinionated workflow doesn't match design iteration needs
- **Verdict**: Best-in-class UI, but engineering-centric workflow model

#### **Figma (built-in project management)**
- **Strengths**: Already in the designer's workflow, kanban boards, file linking, prototype embedding
- **Weaknesses**: Not a task manager — no due dates, reminders, recurring tasks, or personal task capture; collaborative only
- **Designer Use Case**: Great for tracking design iteration status within a single project file, but can't replace a personal task inbox
- **Verdict**: Excellent for design-specific project tracking within files, useless as a personal todo tool

#### **Things 3**
- **Strengths**: Gorgeous UI (macOS/iOS native), Areas + Projects + Tasks hierarchy, tags, deadlines, recurring tasks
- **Weaknesses for UX Designers**: No collaboration features, no attachments, no design-specific categories, Apple ecosystem only
- **Designer Use Case**: Many designers use Things 3 as their personal task manager, but it's a workaround — it wasn't designed with design workflows in mind
- **Verdict**: Best personal task manager for individuals, aesthetically impressive, but no design-specific features

#### **Todoist**
- **Strengths**: Cross-platform, natural language input, labels, filters, collaboration (business plan), integrations
- **Weaknesses for UX Designers**: Generic project management feel, no iteration states beyond "completed," no design-specific templates
- **Designer Use Case**: Often used with a "Design Projects" label structure, but the tool fights the workflow
- **Verdict**: Functional but uninspiring for designers; useful for cross-functional collaboration

#### **OmniFocus**
- **Strengths**: Powerful GTD (Getting Things Done) system, contexts, perspectives, defer dates, automation (AppleScript, URL schemes)
- **Weaknesses for UX Designers**: Complex setup, steep learning curve, expensive, visually dated, no collaboration
- **Designer Use Case**: Power users who subscribe to GTD methodology love it; others find it overwhelming
- **Verdict**: Powerful but not designed for the design process

#### **Craft**
- **Strengths**: Beautiful, document-first, notes + tasks hybrid, Apple ecosystem, design-aware aesthetic
- **Weaknesses for UX Designers**: Task features are secondary to document/note features; less flexible than Notion
- **Designer Use Case**: Some designers use Craft as a personal knowledge base + task manager, but it's not purpose-built for either
- **Verdict**: Aesthetically compelling, but tasks are an afterthought

#### **Height**
- **Strengths**: AI-powered, design-friendly UI, flexible project views (kanban, list, calendar), built for small teams
- **Weaknesses for UX Designers**: Newer tool, less established, collaboration-first (not personal task focus)
- **Designer Use Case**: Emerging choice for design teams who want Linear-like UI with more flexibility
- **Verdict**: Worth watching; bridges gap between personal and team task management

### 2.2 Key Competitor Feature Matrix

| Feature | Notion | Linear | Things 3 | Todoist | OmniFocus | Craft | Height |
|---------|--------|--------|----------|---------|-----------|-------|--------|
| Design-specific categories | No | No | No | Labels | Contexts | No | No |
| Iteration states (not just done/not-done) | Partial | Yes | No | No | No | No | Yes |
| Attachment support (Figma, links) | Yes | Yes | No | Yes | No | Yes | Yes |
| Quick capture | Medium | Fast | Fast | Fast | Fast | Medium | Medium |
| Visual aesthetic | Good | Excellent | Excellent | OK | Dated | Excellent | Good |
| Personal (not team) focus | Yes | No | Yes | Yes | Yes | Yes | Partial |
| Free tier | Yes | Yes (3 projects) | No | Yes | No | Yes | Yes |
| Keyboard-first | No | Yes | Yes | Yes | Yes | No | Yes |

### 2.3 Gaps in Current Market

The research reveals a clear market gap:

1. **No tool is design-process-aware**: Every tool treats tasks as discrete, completable units. None acknowledge the iterative, non-linear nature of design work.
2. **No design-specific task taxonomy built in**: Designers waste time creating custom categories. Pre-built "Research / Design / Testing / Review" would save setup friction.
3. **No lightweight Figma/link attachment**: Tools either have no attachments (Things 3) or over-engineered document systems (Notion). A simple "attach Figma file" or "attach prototype link" would be perfect.
4. **No estimation support**: Design tasks need time estimates (in hours), but no mainstream task tool supports this elegantly.
5. **No design critique tracking**: "Get feedback on X", "Incorporate feedback from critique", "Present to stakeholder" — these are design-specific recurring task patterns that no tool supports natively.

---

## 3. UX Design Patterns for Design-Targeted Tools

### 3.1 Principles for Designing Tools for Designers

Designers are **critical, aesthetic users**. They notice bad typography, inconsistent spacing, and generic interfaces. A todo app for designers must:

1. **Demonstrate craft through its own interface** — If the tool looks like a generic Bootstrap app, designers will dismiss it immediately. The tool must prove it was designed by someone who cares about design.
2. **Respect cognitive load** — Designers are already managing complex visual thinking. The task tool should be visually calm, not competing for attention.
3. **Enable quick capture** — Designers capture tasks in passing (during a meeting, while prototyping). The tool must support sub-3-second task entry.
4. **Support iteration states, not binary completion** — A design doesn't go from "todo" to "done." It goes through "In Progress → Needs Review → Approved." The tool must reflect this.
5. **Allow rich context attachments** — Tasks need Figma links, prototype URLs, screenshots, and notes. Text-only tasks are insufficient.
6. **Enable flexible categorization** — Designer's mental models vary. Some think in projects, others in phases, others in weekly sprints. The tool shouldn't force a rigid hierarchy.

### 3.2 Recommended UI Patterns

#### **Kanban Board (Primary View)**
- **Why**: Mirrors the design process naturally — tasks flow from "Backlog" through "In Progress" through "Review" to "Done"
- **Columns**: Backlog | In Progress | In Review | Approved/Done (configurable)
- **Best Practice**: WIDE columns (designer aesthetic), card-based task display, drag-and-drop reordering within and across columns
- **Reference**: Linear's kanban is the gold standard — clean, fast, keyboard-navigable

#### **Quick Add Bar (Always Visible)**
- **Why**: Designers capture tasks mid-workflow
- **Behavior**: Floating bar at bottom or top of screen, Cmd/Ctrl+K shortcut to focus, natural language date parsing ("tomorrow", "next Monday")
- **Best Practice**: Auto-categorize based on keywords ("research:", "test:", "review:") for zero-friction categorization
- **Reference**: Things 3's quick entry is excellent; Linear's Cmd+K command palette

#### **Design-Specific Task Cards**
- **Why**: Standard todo cards are too generic
- **Elements per card**:
  - Task title (primary)
  - Category tag (Research / Design / Testing / Review / Admin — color-coded)
  - Iteration state indicator (dot: not started / in progress / needs review)
  - Due date (if set)
  - Attachment indicator (Figma icon, link icon)
  - Time estimate (e.g., "2h")
- **States**: Default, hover (subtle lift), dragging (elevated shadow + rotation), completed (strikethrough + muted)
- **Reference**: Linear's issue cards; Notion's database cards

#### **Tag/Label System**
- **Why**: Designers think in multiple dimensions simultaneously (project + phase + priority + type)
- **Behavior**: Multiple tags per task, color-coded, filterable by tag combination
- **Pre-built tags**: Research, Wireframes, High-Fidelity, Prototype, Usability Test, Stakeholder Review, Critique, Handoff
- **Reference**: Todoist's labels system; Linear's labels

#### **Filter/Sort Bar**
- **Why**: One view doesn't fit all moments
- **Filters**: By category, by tag, by iteration state, by due date range, by project
- **Sorts**: By priority, by due date, by created date, by category
- **Best Practice**: Active filters should be visible as removable chips — don't hide filters in menus

### 3.3 Color & Visual Direction

Designers tend to prefer:
- **Neutral, sophisticated base palette** — Light mode: off-white (#FAFAFA) with dark charcoal text (#1A1A1A). Not stark white, not gray.
- **High-contrast accents** — A single bold accent color (designer-curated, not default blue/purple). Could be warm amber, muted teal, or coral. Avoid the typical "SaaS purple" gradient.
- **Dark mode as first-class citizen** — Designers often work in dark mode (easier on eyes, matches Figma dark mode). Support both with equal care.
- **System font preference** — SF Pro on Apple, Inter or Geist on web. Avoid Roboto or default sans-serifs. Designers will notice.
- **Generous whitespace** — If there's one thing designers appreciate, it's intentional negative space. Don't fill every pixel.
- **Subtle, purposeful animations** — Card drag: smooth spring physics. Task complete: satisfying micro-celebration (subtle check animation, not confetti). Filter changes: gentle fade transitions.

**Recommended Palette Direction:**
- Background: #FAFAFA (light) / #0F0F0F (dark)
- Surface/Cards: #FFFFFF (light) / #1A1A1A (dark)
- Text Primary: #1A1A1A (light) / #FAFAFA (dark)
- Text Secondary: #6B6B6B (light) / #8A8A8A (dark)
- Accent: Warm amber (#F59E0B) or Muted teal (#14B8A6) — or configurable by user
- Category tags: Each category gets a distinct muted color (Research=blue, Design=purple, Testing=green, Review=amber, Admin=gray)

### 3.4 Typography

- **Primary Font**: Inter (Google Fonts) — clean, readable, modern, designer-friendly
- **Monospace (for time estimates)**: JetBrains Mono — technical but warm
- **Scale**: 14px base, 1.5 line-height for body, generous letter-spacing on headings
- **Weights**: Regular (400) for body, Medium (500) for task titles, Semibold (600) for section headers

### 3.5 Interaction Patterns

#### **Drag and Drop**
- **Best Practice (from NN/g drag-drop research)**:
  - Show drag preview immediately on mousedown (no delay)
  - Use spring physics for drop animation (tension, friction)
  - Show drop zones with subtle highlight (not aggressive outlines)
  - Allow cross-column drag with smooth insertion animation
  - Keyboard alternative: arrow keys + space to pick up/drop
- **Anti-Patterns**: Stiff animations, no keyboard alternative, small drop targets

#### **Keyboard Shortcuts**
Designers are keyboard-first power users. Essential shortcuts:
- `Cmd/Ctrl + K`: Focus quick add
- `Cmd/Ctrl + F`: Focus filter/search
- `J/K` or Arrow keys: Navigate tasks
- `Enter`: Open task detail
- `D`: Set due date on selected task
- `L`: Add label/tag
- `1-5`: Set priority
- `Cmd/Ctrl + Enter`: Complete task
- `Esc`: Close modals, cancel actions

#### **Task Detail Modal**
- Slide-in panel (not blocking modal) from the right
- Sections: Title (editable inline), Category selector, Tags, Due date, Time estimate, Description (rich text), Attachments (Figma URLs, prototype links), Iteration state, Subtasks
- **Key Pattern**: Click outside or press Esc to close — no "Save" button needed (auto-save on every change)

#### **Empty States**
- **Important**: Designers are sensitive to empty states — they're the first interaction with the tool
- Should include: Friendly illustration (simple line art, not cartoon), encouraging copy ("Your design backlog is clear — add your first task"), prominent "Add Task" button
- **Reference**: Craft's empty states are warm and helpful, not generic

---

## 4. Mental Model & Information Architecture

### 4.1 How UX Designers Think About Tasks

**Primary Mental Model**: Projects with phases, not flat lists

A designer's mental model typically looks like:

```
└── Checkout Redesign (Project)
    ├── Research Phase
    │   ├── User interviews (5 sessions)
    │   ├── Competitive analysis
    │   └── Synthesis & persona updates
    ├── Design Phase
    │   ├── Wireframe explorations
    │   ├── High-fidelity mockups
    │   └── Design system updates
    ├── Testing Phase
    │   ├── Usability test recruitment
    │   ├── Test sessions
    │   └── Findings report
    └── Handoff
        ├── Dev sync
        └── Asset export
```

**Secondary Mental Model**: Weekly focus areas

Designers also think in time-boxed chunks:
- "This week: Get wireframes approved"
- "Next sprint: Usability testing for checkout"

### 4.2 Recommended Information Architecture

```
┌─────────────────────────────────────────────────────────┐
│  HEADER                                                 │
│  [App Name/Logo]          [Search] [Filter] [Settings]  │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌─────────────────────────────────────────────────┐  │
│  │  QUICK ADD BAR (always visible)                  │  │
│  │  "Add a task... (Cmd+K)"                          │  │
│  └─────────────────────────────────────────────────┘  │
│                                                         │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ │
│  │ Backlog  │ │In Progress│ │In Review │ │ Approved │ │
│  │ (count)  │ │ (count)   │ │ (count)  │ │ (count)  │ │
│  ├──────────┤ ├──────────┤ ├──────────┤ ├──────────┤ │
│  │ [Task]   │ │ [Task]   │ │ [Task]   │ │ [Task]   │ │
│  │ [Task]   │ │ [Task]   │ │          │ │ [Task]   │ │
│  │ [Task]   │ │          │ │          │ │          │ │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘ │
│                                                         │
├─────────────────────────────────────────────────────────┤
│  FOOTER / STATUS BAR                                    │
│  [Task count: X total, Y active]    [Keyboard hints]   │
└─────────────────────────────────────────────────────────┘
```

**Navigation Structure:**
1. **Board View** (default): Kanban columns as described above
2. **List View**: Flat or grouped list with same columns — for power users who prefer density
3. **Focus View**: Single task at a time — for deep work sessions
4. **Calendar View**: Tasks plotted on calendar — for deadline management

**Settings/Preferences:**
- Theme (Light / Dark / System)
- Default category for new tasks
- Kanban column customization
- Keyboard shortcut remapping
- Data export (JSON)

---

## 5. Key UX Metrics to Track

Based on research from similar productivity app research and Nielsen Norman Group UX metrics:

| Metric | Target | Measurement |
|--------|--------|-------------|
| **Time to first task added** | < 10 seconds | From app open to first task in list |
| **Task completion rate** | > 70% of created tasks | Tasks marked complete / tasks created |
| **Filter usage rate** | > 50% of sessions | Sessions where user applies at least one filter |
| **Return rate** | > 60% daily active users return next day | DAU / MAU ratio |
| **Drag-and-drop usage** | > 40% of task state changes | Tasks moved via drag / total state changes |
| **Keyboard shortcut adoption** | > 30% of power users | Users who use Cmd+K within first week |
| **Category assignment rate** | > 80% of tasks get categorized | Tasks with category / total tasks |
| **Attachment rate** | > 25% of tasks have links | Tasks with attachments / total tasks |

---

## 6. Anti-Patterns to Avoid

### Critical Anti-Patterns:

1. **Don't use default blue/purple gradients** — Designers will immediately classify this as "generic SaaS tool"
2. **Don't force account creation before first task** — Zero-friction entry is essential; sign-up after first task capture
3. **Don't hide filters in menus** — Visible, chip-based filters that users can see and remove
4. **Don't use "Importance" or "Priority" labels alone** — Use color + position + icon, never color alone (accessibility)
5. **Don't require mouse-only interactions** — Every action must have a keyboard alternative
6. **Don't use empty modals for task editing** — Slide-in panels that don't block the board view
7. **Don't add social/collaboration features by default** — This is a personal task tool; collaboration can be phase 2
8. **Don't make the onboarding a 5-step wizard** — One screen of getting started, then direct access to the app
9. **Don't forget about dark mode** — Design it concurrently, not as an afterthought
10. **Don't animate unnecessarily** — Every animation must serve a purpose (feedback, orientation, or delight — not decoration)

### Important Anti-Patterns:

11. **Don't use generic icons** — Use a cohesive icon set (Phosphor, Lucide, or similar designer-friendly library)
12. **Don't display "Notifications" prominently** — This is a personal tool; no notification badges or unread states
13. **Don't make "Complete" a single click that disappears** — Show completed tasks with strikethrough for at least 24 hours before archiving
14. **Don't mix projects and tasks in the same view** — Use a clean separation; tasks belong to columns, not nested under visible projects
15. **Don't use Toast notifications** — For personal tools, toasts are annoying; update the UI directly

---

## 7. Accessibility Considerations

Based on WCAG 2.1 AA (minimum) and designer expectations:

1. **Color independence**: Never convey meaning through color alone. Always pair with text labels, icons, or position.
2. **Touch targets**: Minimum 44x44px on all interactive elements (Apple HIG).
3. **Contrast ratios**: 4.5:1 minimum for normal text, 3:1 for large text.
4. **Focus indicators**: Visible, styled focus rings (not browser default) for keyboard navigation.
5. **Screen reader**: All interactive elements must have proper ARIA labels; drag-and-drop must have keyboard alternative.
6. **Reduced motion**: Respect `prefers-reduced-motion` — disable spring animations, use opacity fades instead.
7. **Font resizing**: UI must not break at 200% browser zoom.

---

## 8. Research Sources & References

### Designer Workflow Research:
- Nielsen Norman Group: "Creative Professionals" UX research
- Nielsen Norman Group: "Drag-Drop UX" best practices
- UX Collective (uxdesign.cc) — Designer productivity tool discussions
- Designer News — Designer tool preferences and critiques
- Reddit r/userexperience — Designer workflow discussions

### Competitive Product Analysis:
- Notion (notion.so) — Feature analysis, personal use
- Linear (linear.app) — UI patterns, workflow model
- Things 3 (culturedcode.com) — macOS/iOS design patterns
- Todoist (todoist.com) — Label/filter patterns
- OmniFocus (omnigroup.com) — GTD methodology integration
- Craft (craft.do) — Document-task hybrid patterns
- Height (height.app) — Emerging design-friendly task tool

### Design System References:
- Inter font (rsms.me/inter)
- Phosphor Icons (phosphoricons.com)
- Apple Human Interface Guidelines (for macOS/iOS native feel)
- Material Design 3 (for web component patterns)
- Gestalt principles (proximity, similarity, closure)

---

## 9. Summary: Key Research Findings

### The Core Opportunity

UX designers are an underserved market for task management tools. They currently cobble together solutions (Often Things 3 + Notion + Figma boards) because no single tool is:
1. **Design-process-aware** — reflects the iterative, non-binary nature of design work
2. **Aesthetically excellent** — beautiful enough that designers trust it
3. **Design-context-aware** — pre-built categories, attachment types, and workflows that match how designers actually work
4. **Fast and keyboard-first** — respects the designer's time and workflow

### Must-Have Features (from research):
1. **Kanban board with 4+ configurable columns** (Backlog, In Progress, In Review, Done)
2. **Design-specific category tags** (Research, Design, Testing, Review, Admin)
3. **Quick add bar** with natural language date parsing (Cmd+K)
4. **Drag-and-drop** reordering with keyboard alternative
5. **Task detail panel** with attachment support (Figma links, prototype URLs)
6. **Filter chips** — visible, removable, combinable
7. **Dark mode** — first-class, not an afterthought
8. **Keyboard shortcuts** — full keyboard navigation

### Differentiating Features (competitive edge):
1. **Design critique workflow** — special task type for "Get feedback", "Incorporate feedback"
2. **Time estimation** — built-in hours/minutes estimate per task with total project time
3. **Subtask support** — tasks can have checklist items (e.g., "Conduct 5 user interviews" → 5 checkboxes)
4. **Calendar view** — tasks plotted on calendar for deadline visualization
5. **Local-first** — data stored locally (localStorage), no account required, privacy-respecting

### Recommended Tech Stack (from template selection):
- **Vite + Tailwind CSS + TypeScript** — Fast development, Tailwind for design system consistency, TypeScript for maintainability
- **LocalStorage for persistence** — No backend needed, data stays on device, exportable as JSON
- **Single-page application** — All interaction in one view (board), no page reloads

---

**Status: READY_FOR_DESIGN**

*This research document provides the foundation for design decisions. The next phase is to take these findings and create a specific, visual design proposal.*
