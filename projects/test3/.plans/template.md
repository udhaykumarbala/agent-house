# Template Selection

## Chosen Template: `nextjs-frontend`

## Task Analysis

A specialized todo app for designers is an **interactive frontend application**, not a static page. The CEO directive is clear: this app must itself demonstrate design excellence — it's a tool for a design-conscious audience who will judge it by its own aesthetics. This means:

- Rich, polished UI interactions (drag-and-drop, smooth animations, micro-interactions)
- Component-based architecture for reusability and consistency
- State management for todo items (add, complete, delete, categorize, filter)
- Design system approach (tokens, consistent visual language)
- Responsive across device sizes (designers work on multiple screens)

## Justification

- **Interactive SPA requirements**: Todo apps need client-side state (add/complete/delete/reorder tasks), filtering, and real-time UI updates — all hallmarks of a single-page application that `static-html`/`static-enhanced` can't handle cleanly without hacky vanilla JS
- **Design system architecture**: The `nextjs-frontend` template includes shadcn/ui + Tailwind, which provides a component library that can be customized to demonstrate design excellence. Building on proven design primitives rather than raw CSS gives us more time to focus on what makes this app special
- **Designer audience expectation**: Designers expect tools that feel premium — smooth animations, cohesive visual language, thoughtful micro-interactions. A React-based approach with Tailwind enables this level of polish
- **Complexity matched**: `fullstack` would add unnecessary infrastructure complexity (Go backend + Docker + PostgreSQL) when there's no API requirement, auth, or persistent server-side logic needed. `nextjs-frontend` captures all the frontend benefits without the overhead
- **Tooling alignment**: Next.js 14 + TypeScript + Tailwind + shadcn/ui is a proven modern stack that the team can execute quickly while delivering high-quality output

## Customizations Needed

- **Override default shadcn/ui theme** with a bespoke design system tuned for a designer's aesthetic (custom color palette, typography, spacing, motion)
- **Add Framer Motion** for orchestrated animations and transitions
- **Include drag-and-drop library** (e.g., @dnd-kit) for reordering tasks
- **Persist state** via localStorage for MVP, with easy upgrade path to a backend later

---
Status: AWAITING_REVIEW
