# Research: Designer-Specific Todo App

## 1. What Makes Designer Workflows Unique

### The Creative Brief to Completion Cycle

Designers work differently from developers or general knowledge workers. Their work follows a **non-linear, iterative cycle** that generic todo apps fail to capture:

- **Brief/Context** → **Ideation** → **Exploration** → **Refinement** → **Review** → **Revision** → **Handoff/Delivery**
- Tasks are rarely "done" in one shot — they cycle through feedback loops
- Work is often **parallel** (multiple projects at different stages) rather than sequential
- **Time perception** is different: designers think in "sessions" not hours (e.g., "I have a 2-hour focused design block")

### Task Types Designers Manage

Unlike generic todos ("Buy groceries"), designers track:

| Task Type | Examples |
|-----------|----------|
| **Creative tasks** | "Explore 3 logo directions", "Sketch homepage layout" |
| **Feedback tasks** | "Incorporate client comments on mockup v2", "Address review notes" |
| **Asset tasks** | "Export SVGs for dev handoff", "Generate retina assets" |
| **Research tasks** | "Find 10 references for landing page", "Audit competitor color palettes" |
| **Meeting prep** | "Prepare presentation for design review" |
| **Revision cycles** | "Round 2 revisions: typography + spacing", "Round 3: mobile responsiveness" |
| **Administrative** | "File project to Figma library", "Update status in Notion" |

### Designer-Specific Context

- **Projects, not just tasks**: Design work is organized by project (Brand Redesign, Landing Page, App Icon), with tasks living inside them
- **Visual references matter**: "Attach a screenshot" or "link to Figma frame" makes a task actionable in ways plain text can't
- **Due dates are fluid**: Deadlines exist but creative work doesn't fit rigid date slots — designers need **flexible scheduling** (e.g., "Today", "This Week", "Someday")
- **Energy/timing awareness**: Designers have peak creative hours (often mornings) and administrative hours (end of day). Tasks can be tagged by **energy level** (Deep Work, Quick Win, Admin)
- **Client/collaborator context**: Tasks often belong to a client or project stakeholder

---

## 2. Competitive Landscape

### Tools Designers Currently Use (and Their Gaps)

| Tool | What Designers Love | Gaps for Todo Management |
|------|---------------------|--------------------------|
| **Figma** | Native to design work, projects = files, prototypes | Not a task manager; no personal todos across projects |
| **Notion** | Flexible databases, linked pages, aesthetic | Over-flexible (analysis paralysis), too heavy for quick capture, defaults are ugly |
| **Linear** | Fast, keyboard-driven, beautiful UI, project + issue tracking | Too engineering-focused ("issue" terminology), no personal/timestamped todos, no creative project templates |
| **Things 3** | Beautiful macOS/iOS app, Areas + Projects + Todos hierarchy, elegant UX | Mac/iOS only, no collaboration, no visual attachments, no cross-project view |
| **Todoist** | Cross-platform, natural language input, labels/filters | Generic interface, not designed for project context, boring aesthetics |
| **Craft** | Beautiful documents, blocks, dark mode, visual hierarchy | Document-focused, not task-focused, steep learning curve |
| **Amie** | Gorgeous UI, calendar + todos in one, social feel | Still early, limited project hierarchy, no attachments |
| **Raycast** | Lightning fast, extensions, scriptable | Not a dedicated todo app, developer-focused |
| **Miro** | Visual whiteboarding, infinite canvas | Not a todo manager, but designers use sticky notes on Miro boards as informal task boards |

### Niche Tools Targeting Designers

- **GoVisually** — Design review and approval platform specifically for creative agencies. Revision rounds tracked, approval recorded with timestamps. Solves *one problem* (design review) better than any general tool.
- **Frame.io** — Gold standard for **visual feedback on video design work**. Timestamp-linked comments, side-by-side version comparison, stakeholder review links, approval workflows.
- **Height** — Emerging AI-native project management. Auto-tasking (AI generates task breakdowns from descriptions), contextual AI summaries, adaptive workflows that learn team conventions.
- **Red Pen** — Design feedback/annotation tool
- **Mockplus** — Prototyping tool with lightweight task features
- **Pixeltrue** — Design feedback + task tracking

**Gap in the market**: No tool combines **Things 3-level polish** with **Linear-level speed** and **project-based hierarchy** specifically for design professionals. Most designers cobble together 3-4 tools.

---

## 3. Pain Points with Generic Todo Apps

1. **Binary completion model**: Design work is non-linear. "I implemented the feedback" ≠ "client approved it." A task might be "done from my side" but pending external approval. Generic apps treat tasks as `done` or `not done` — but design exists on a spectrum.
2. **Feedback lives elsewhere**: Design feedback arrives via Slack, email, Figma comments, Loom, meetings, WhatsApp. Generic apps have zero awareness of this feedback. Designers must manually convert every message into a task — a massive overhead.
3. **No visual context**: When a designer sees "Redesign hero section" in a generic app, they must open Figma, find the file, navigate to the frame. The visual reference should live alongside the task.
4. **No project context**: A task "Fix the button" is meaningless without knowing which project/file it's for.
5. **Flat lists only**: Designers think in hierarchies (Project → Phase → Task), not flat lists.
6. **Boring defaults**: Generic todo apps look like productivity software from 2010. Designers notice bad typography, inconsistent spacing, and clashing colors immediately.
7. **Wrong terminology**: "Issues", "Tickets", "Bugs" — engineering language that breaks creative flow. Design tools should use design language.
8. **Scope creep is invisible**: Generic tools have no concept of how many times a task has been revised or how far it has drifted from the original brief.
9. **No sense of creative progress**: Generic apps show "3 of 12 tasks complete" — designers want to see they've progressed from "rough sketch" to "polished mockup" through the design phases.
10. **Static, lifeless UI**: No motion, no personality, no delight. The tool itself looks like it wasn't designed.
11. **No energy/time awareness**: Can't mark a task as "2-hour deep work block" vs "15-min quick edit."
12. **Collaborator-blind**: Most personal todo apps ignore "waiting on client" or "in design review" states — the most common states in creative work.
13. **The "in-between" work is invisible**: Research, incubation, exploration, informal conversations that shape direction — generic tools capture only tasks with clear beginnings and ends.

---

## 4. Unique Feature Opportunities

### 4.1 Design-Phase-Aware Tasks

Tasks have a **phase tag** that maps to the creative workflow:
- **Brief** → **Ideate** → **Draft** → **Refine** → **Review** → **Approved**

This gives designers a **sense of creative progress** rather than just "completed/not completed".

### 4.2 Visual Task Cards

Each task can have:
- A **color tag** (e.g., brand project = teal, marketing = coral)
- A **small preview** (Figma thumbnail, screenshot, or emoji canvas)
- **Priority indicator** that uses visual weight, not just a number

### 4.3 Figma-First Linking

- Paste a Figma URL → auto-fetches the frame/file thumbnail
- Tasks show inline previews of linked designs
- One-click "Open in Figma" action

### 4.4 Energy Level Tags

- **Deep Work** (2+ hours, focused, no interruptions)
- **Medium Lift** (45-90 min, some context switching)
- **Quick Win** (15-30 min, low cognitive load)
- **Admin** (meetings, emails, reviews)

### 4.5 Session-Based Time Blocking

Instead of due dates, designers think in **design sessions**:
- "This morning" / "This afternoon" / "Tomorrow" / "This week" / "Someday"
- Optional: actual time slots (9:00 AM - 11:00 AM)
- Visual calendar-adjacent view alongside task list

### 4.6 Creative Status States

Beyond "done/not done":
- **To Do** / **In Progress** / **In Review** / **Changes Requested** / **Approved** / **Done**
- These map to the feedback cycle designers actually live in

### 4.7 Project Dashboard

Quick overview of all active projects with:
- Progress bar (based on phase completion)
- Next action highlight
- Days since last update
- Overdue indicator

### 4.8 Focus Mode

One task visible at a time, full-screen, with:
- Figma preview (if linked)
- Timer option
- Minimal chrome
- Next task peek

---

## 5. UI/UX Patterns That Signal "Designed for Designers"

### 5.1 Visual Language

**Reference tools**: Linear, Things 3, Craft, Amie, Raycast

| Element | What Great Tools Do |
|---------|---------------------|
| **Typography** | Custom or carefully chosen system font — not Inter defaults. Linear uses system font but with perfect weight contrast. Things 3 uses its own "Now" font. |
| **Color** | Muted, sophisticated palette with one strong accent. Linear's purple, Amie's coral/indigo, Craft's warm neutrals. Avoid pure black/white backgrounds — use subtle tints. |
| **Spacing** | Generous whitespace. 8px grid system. Consistent internal padding (16px, 24px, 32px). Nothing feels cramped. |
| **Shadows** | Subtle, layered shadows for elevation. Things 3 uses soft shadows. Linear uses no shadows but uses background color shifts for depth. |
| **Borders** | 1px borders with low-contrast colors (not black). Rounded corners — usually 6-12px. |
| **Motion** | Spring-based animations (Framer Motion). 200-300ms transitions. Staggered list animations. Drag-and-drop with physics. Nothing feels instant or jarring. |

### 5.2 Interaction Patterns

**Quick Capture**: Cmd+N / Cmd+K to capture a task instantly, anywhere in the app. No navigating, no modals. Things 3 and Todoist pioneered this.

**Keyboard-first navigation**: Every action reachable from keyboard. J/K to navigate, X to select, E to edit. Linear is the gold standard here — designers who use Linear expect keyboard speed everywhere.

**Drag-and-drop with context**: Reorder tasks by dragging. Drop zones highlight with visual feedback. Smooth 60fps animations.

**Inline editing**: Click to edit, not "select → edit button → modal". Double-click to rename, click checkbox to complete.

**Swipe actions**: On mobile/touch, swipe left for delete, swipe right for complete. Things 3 and Linear both use this.

**Contextual menus**: Right-click for options, not cluttered toolbars.

### 5.3 Layout Patterns

**Sidebar + Content (Two-Pane)**: Classic pattern used by Linear, Notion. Sidebar for projects/navigation, main area for task list. Sidebar can collapse.

**Kanban View** (optional): For design teams managing multiple rounds of review. Columns = phases (Ideate → Draft → Review → Approved). Drag cards between columns.

**Calendar View** (optional): Session-based time blocks laid out on a weekly calendar. Integrates todos with time.

**Focus Mode**: Single-task full-screen view with minimal UI.

### 5.4 Color System

Great designer tools use **semantic color naming** with dark/light mode:

```
Background:      #0F0F0F (dark) / #FAFAFA (light)
Surface:        #1A1A1A (dark) / #FFFFFF (light)
Surface Raised: #262626 (dark) / #F5F5F5 (light)
Border:         #333333 (dark) / #E5E5E5 (light)
Text Primary:   #FFFFFF (dark) / #0A0A0A (light)
Text Secondary: #A0A0A0 (dark) / #6B6B6B (light)
Accent:         #7C5CFF (purple — Linear's approach) or
                #E07A5F (coral — warm alternative)
Success:        #4ADE80
Warning:        #FBBF24
Destructive:    #F87171
```

### 5.5 Motion Design Principles

1. **Spring physics**: Use `stiffness: 300, damping: 30` or similar spring configs. Never linear easing.
2. **Staggered reveals**: List items animate in with 30-50ms stagger between each item.
3. **Drag feedback**: Items scale up slightly (1.02) and get a shadow boost when dragged.
4. **Micro-confirmations**: Checkbox completion triggers a satisfying scale + color animation.
5. **Page transitions**: Subtle fade + slide (translateY 8px → 0) between views.
6. **Orchestrated loading**: Skeleton screens that fade in with stagger, not spinners.

### 5.6 Dark Mode

Designers overwhelmingly prefer dark mode for design work (reduces eye strain when matching colors). The app must:
- Have a first-class dark mode (not an afterthought)
- Use slightly warm dark tones (#0F0F0F over pure black #000000)
- Ensure text has sufficient contrast (WCAG AA minimum)
- Provide system-level auto-switching

---

## 6. Information Architecture

### Navigation Hierarchy

```
├── Inbox (quick capture, unsorted)
├── Today
│   ├── Morning Session
│   └── Afternoon Session
├── Upcoming
│   ├── This Week
│   └── Next Week
├── Someday (backburner / exploring)
├── Archive
└── Projects (collapsible list)
    ├── Brand Redesign
    │   ├── Phase: Concept / Design / Review / Handoff
    │   ├── Tasks (active)
    │   ├── Completed
    │   └── Resources
    ├── Landing Page
    └── App Icon
```

### Enhanced Task Data Model

```typescript
interface Task {
  id: string;
  title: string;
  description?: string;           // Rich text, can include Figma links, notes

  // Organization
  projectId?: string;
  colorTag?: string;             // Project color for visual identification

  // Design workflow state
  phase: 'brief' | 'ideate' | 'concept' | 'design' | 'review' | 'revision' | 'approved' | 'handoff' | 'delivered';
  status: 'todo' | 'in_progress' | 'in_review' | 'changes_requested' | 'approved' | 'done';
  iterationRound: number;         // Tracks how many revision rounds (Round 1, Round 2...)

  // Time & energy
  energy: 'deep' | 'medium' | 'quick' | 'admin';
  dueSession?: 'morning' | 'afternoon' | 'evening' | 'any';
  dueDate?: Date;

  // Visual context
  figmaUrl?: string;
  figmaThumbnail?: string;        // Auto-fetched thumbnail preview
  coverImage?: string;           // Optional cover/attachment

  // Priority
  priority: 1 | 2 | 3;

  // Timestamps
  createdAt: Date;
  updatedAt: Date;
  completedAt?: Date;
}

interface Project {
  id: string;
  name: string;
  color: string;               // Hex color for visual identification
  description?: string;
  coverImage?: string;         // Project cover / hero image
  currentPhase: ProjectPhase;  // Aggregate phase across all tasks
  tasks: Task[];
  createdAt: Date;
}
```

---

## 7. Technical Considerations

### Frontend Stack (nextjs-frontend)

- **Next.js 14** with App Router — for component architecture and routing
- **TypeScript** — for type safety across the data model
- **Tailwind CSS** — for styling (enables rapid iteration on design tokens)
- **Framer Motion** — for spring-based animations and orchestrated transitions
- **@dnd-kit** — for drag-and-drop reordering (accessible, performant)
- **Zustand** — lightweight state management (simpler than Redux, more structure than useState)
- **localStorage** — MVP persistence, upgrade path to backend later

### Design System

- Use **CSS custom properties** for design tokens (colors, spacing, typography)
- Build a **component library** on top of shadcn/ui primitives (override everything to match the design language)
- Establish **motion tokens** (spring configs, durations) as constants
- Never use default shadcn styling — every component must be customized

### Accessibility

- Full keyboard navigation (this is non-negotiable for power users)
- Focus indicators that match the design language (not browser defaults)
- ARIA labels on all interactive elements
- Motion respects `prefers-reduced-motion`
- Color contrast WCAG AA minimum

---

## 8. Strategic Positioning

### This App Is The Product

The CEO directive is critical: **the tool must demonstrate design excellence**. Designers will judge a "todo app for designers" the same way they judge any design tool — if the todo app looks generic, they've already lost credibility.

### Core Differentiators

1. **Phase-aware tasks** — Maps to the actual creative process, not a generic checkbox
2. **Figma-native linking** — First-class integration with the designer's primary tool
3. **Session-based scheduling** — Respects how designers actually plan their days
4. **Design-phase polish** — Every interaction feels intentional, smooth, and beautiful
5. **Energy tagging** — Helps designers match tasks to their creative energy

### Scope Control (MVP vs Full)

**MVP (Phase 1)**:
- Task CRUD with phase/status/energy
- Project organization with color tags
- Session-based time slots (Today/This Week/Someday)
- Dark/light mode with first-class dark mode
- Keyboard navigation (J/K, Cmd+K, Enter to edit)
- Beautiful, polished UI with spring animations
- Inline editing throughout
- Quick capture (Cmd+N anywhere)

**Post-MVP**:
- Figma URL link previews with auto-thumbnails
- Kanban view with drag between phases
- Focus mode with timer (Pomodoro-style)
- Calendar view (session-based time blocks)
- Iteration round tracking
- Cloud sync / backend
- Natural language input parsing
- Auto-generated design spec export on handoff

---

## 9. Key Interaction Patterns for Creative Professionals

### Command Palette (Cmd+K)

The universal command palette is the single most impactful interaction pattern for power users. Designers should be able to: create tasks, switch projects, search across all content, run actions, and navigate entirely from the command palette — without touching the mouse.

### Natural Language Input

"Design hero section by Friday" should:
- Parse "hero section" as the task title
- Parse "Friday" as the deadline
- Auto-assign to the current project
- Create a well-formed task in under 2 seconds

### Smart Defaults with Manual Override

When a designer creates a task, intelligent defaults should be offered:
- Suggested project (based on recent work)
- Suggested phase/status (based on project workflow)
- Suggested energy level (based on task type)
- Suggested tags (based on project type)

But every suggestion should be one-click dismissible or changeable.

### Kanban + List Toggle

Different tasks require different views:
- **Kanban**: For tracking status across phases (what's in review? what's approved?)
- **List**: For detailed work (what's on my plate today?)
- **Focus**: Single-task full-screen view with distraction-free canvas
- **Calendar**: For deadline management with session-based time blocks

### Visual Comparison (Before/After)

For revision-heavy workflows, side-by-side comparison is essential:
- Show "before" and "after" versions of a design task
- Highlight what changed between rounds
- Allow stakeholders to compare without opening Figma

### Iteration Counting

Track **Round 1, Round 2, Round 3** visually on task cards:
- Badge showing current round
- History of all rounds with timestamps
- Comparison between rounds (for post-MVP)

### Empty State as Teacher

When a designer opens a new project with no tasks, the empty state should:
- Show example tasks to demonstrate the workflow
- Offer a template to populate the project
- Explain what this project type typically needs
- Feel encouraging, not sterile

### Review Rituals Built In

Creative work benefits from regular review. The app should:
- Prompt daily: "What's on your mind for today?"
- Prompt per-project: "What decisions were made this week?"
- Make reviews feel natural, not like extra work

---

## 10. Key Principles for a Designer-Focused Task Tool

### Principle 1: Visual Context Is the Foundation
Every task should be able to show its current design state. The tool should integrate with Figma, not replicate it. Visual previews make tasks actionable in ways plain text cannot.

### Principle 2: Feedback Must Be First-Class
Design feedback is visual, multi-source, and multi-priority. The tool should make feedback aggregation easy — or at minimum, not require manual conversion of every Slack message/email into a task.

### Principle 3: Iteration Is Native
Design is inherently iterative. Round counting, version comparison, and rollback are not edge cases — they are the core workflow. Tasks must track "how many times have I revised this?"

### Principle 4: The Tool Must Be Beautiful
Designers will not trust their workflow to an ugly tool. The interface is a statement of values. If the todo app looks generic, designers have already lost trust — before they even use a single feature.

### Principle 5: Speed Is a Feature
Designers are power users who work fast. Sub-second response times, keyboard shortcuts everywhere, command palette access. Linear proved that a project management tool can feel as fast as a terminal.

### Principle 6: Multiple Contexts Simultaneously
Designers filter by client, by phase, by energy level, by deadline. Rigid views do not work. Custom filters, tags, and simultaneous views are essential. A task should belong to a project AND have an energy level AND have a deadline — not forced to choose one.

### Principle 7: Non-Linear Is the Default
Designers jump between projects, contexts, and creative modes constantly. Quick capture, global search, and one-click project switching enable this workflow. The app should never require drilling through menus.

### Principle 8: Integration Over Duplication
The tool should not try to replace Figma, Loom, or Notion. It should integrate with them — pulling design previews from Figma, video feedback from Loom, briefs from Notion. Meet designers where they already are.

### Principle 9: Respect the Invisible Work
Research, incubation, exploration, and informal conversation are all part of the creative process. The tool should create space for ambient work — "Someday" / "Backburner" / "Exploring" states — not force everything into executable tasks.

### Principle 10: Opinionated Defaults
The tool should work out of the box with designer-relevant phases, standard project templates, and sensible conventions. Designers can customize, but they should not *need* to configure everything from scratch. A well-designed default experience is itself a feature.

---

## 11. The App's Personality & Voice

A designer-focused todo app must have a distinct **point of view**. It should feel like it was made by a designer who was frustrated with existing tools — not a committee.

### Tone

- **Confident, not prescriptive**: It has opinions about how design work should be organized, but makes it easy to adapt.
- **Quiet, not empty**: The UI is restrained. Information density is high when needed, quiet when focus is required.
- **Precise**: Typography, spacing, and color are exact. Nothing feels accidental.
- **Warm**: Not cold and corporate. Not playful and bubbly. Sophisticated but human.

### Brand Positioning

The app occupies the space between:
- **Things 3** (personal, beautiful, macOS-native) and **Linear** (fast, keyboard-driven, opinionated)
- **Personal task manager** and **project management tool**
- **For developers** (Linear, Notion) and **for designers** (this app)

The tagline could be: *"The todo app designed for how designers actually work."*

### Visual Identity Direction

Drawing from the reference tools:
- **Linear's dark mode** — the reference standard for dark productivity apps
- **Craft's warmth** — slightly warm tones rather than cold grays
- **Amie's typography** — carefully chosen, distinctive font stack
- **Things 3's information hierarchy** — projects, areas, tasks with clear nesting

---

## 12. Key Research Sources & References

| Tool | Lesson for This App |
|------|---------------------|
| **Linear** | Speed as a feature. Keyboard-first navigation. Opinionated design. Dark mode default. |
| **Things 3** | Quiet productivity. Daily review ritual. Elegant hierarchy. Session-based organization. |
| **Craft** | Aesthetic investment creates trust. Information block system. Apple-native feel. |
| **Amie** | Calendar+todos integration. Stunning warm aesthetics. |
| **Figma** | How designers organize visual work. File/project model. Contextual comments. |
| **Raycast** | Command palette as the primary interface. Extension ecosystem. |
| **GoVisually** | Design review is a first-class problem. Revision rounds with approval tracking. |
| **Frame.io** | Visual feedback with timestamp-linked comments. Version comparison. |
| **Height** | AI-assisted task breakdown. Adaptive workflows. |
| **Notion** | Database relationships as power feature. The danger of over-flexibility. |
