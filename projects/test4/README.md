# DesignFlow — UX Designer Todo App

A kanban-style task management app built specifically for UX designers. Mirrors the double diamond design process — tasks flow through **Backlog → In Progress → In Review → Approved** — rather than the binary complete/incomplete model that generic todo apps impose.

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Build | Vite 5 |
| UI Framework | React 18 |
| Styling | Tailwind CSS 3 |
| State | Zustand |
| Drag & Drop | @dnd-kit |
| Icons | Lucide React |
| Dates | date-fns |
| IDs | nanoid |

## Getting Started

```bash
# Install dependencies
npm install

# Start dev server
npm run dev

# Production build
npm run build
```

Dev server runs at `http://localhost:5173`

## Design System

- **Fonts**: Cabinet Grotesk (headings) + Inter (body) + JetBrains Mono (estimates)
- **Theme**: Light (warm cream) / Dark (deep warm black) — toggled via `dark` class on `<html>`
- **Palette**: Indigo primary (#3730A3) + Coral accent (#F97316)
- **Categories**: Research, Design, Testing, Review, Discovery, Handoff, Admin

## Data Storage

All data is stored in **localStorage**:
- `designflow_tasks` — task data with auto-save (300ms debounce)
- `designflow_settings` — theme, categories, tags

Export to JSON available from the Settings panel.

## Features (MVP Scope)

- [ ] Kanban board with 4 columns
- [ ] Quick add bar with natural language parsing
- [ ] Drag-and-drop reordering
- [ ] Task cards with category, priority, due date, time estimate
- [ ] Task detail panel (slide-in)
- [ ] Subtasks within tasks
- [ ] URL attachments
- [ ] Category filter bar
- [ ] Light + Dark mode
- [ ] Search
- [ ] JSON export

## Architecture

See [uiflow.md](uiflow.md) for detailed UI flow documentation.

## API

See [api_test.html](api_test.html) for the planned REST API contract (Phase 2+ backend).
