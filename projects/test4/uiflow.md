# DesignFlow — UI Flow Documentation

**Version:** 1.0
**Phase:** 1 (Foundation — frontend scaffold only)

---

## Architecture Overview

DesignFlow is a single-page application (SPA) built with Vite + React + TypeScript + Tailwind CSS.

```
Browser
  │
  ├── index.html (entry, theme flash prevention script)
  │       │
  │       └── src/main.tsx (React root)
  │               │
  │               └── src/App.tsx
  │                       │
  │                       ├── Zustand Stores (client-side state)
  │                       │       ├── taskStore     → localStorage: 'designflow_tasks'
  │                       │       ├── uiStore       → ephemeral (not persisted)
  │                       │       └── settingsStore  → localStorage: 'designflow_settings'
  │                       │
  │                       └── Components (built in Phase 2+)
  │                               ├── layout/       (AppShell, Header, StatusBar, FilterBar)
  │                               ├── board/        (KanbanBoard, KanbanColumn, TaskCard)
  │                               ├── task/         (TaskDetailPanel)
  │                               ├── shared/       (QuickAddBar, CategoryChip, EmptyState)
  │                               ├── settings/     (SettingsPanel)
  │                               └── modals/       (ConfirmDialog)
  │
  └── localStorage
          ├── designflow_tasks     (tasks array + _version counter)
          └── designflow_settings  (theme, columnOrder, categories, tags)
```

---

## State Flow

### Task Creation Flow
```
User types in QuickAddBar
        │
        ▼
naturalLanguage.ts parses → category, dueDate, timeEstimate
        │
        ▼
taskStore.addTask() → creates Task with nanoid
        │
        ▼
  _version++ → triggers persist middleware
        │
        ▼
localStorage.setItem('designflow_tasks', JSON.stringify(...))
        │
        ▼
saveStatus: 'idle' → 'saving' → 'saved' (after 350ms)
```

### Theme Flow
```
On page load:
  1. Inline <script> reads localStorage → sets 'dark' class on <html>
  2. React hydrates
  3. settingsStore rehydrates → applyTheme() called
  4. If theme='system' → matchMedia listener active

On user toggle:
  settingsStore.updateTheme() → applyTheme()
        │
        ▼
  document.documentElement.classList.toggle('dark')
        │
        ▼
  CSS custom properties switch (via .dark {})
```

### Drag & Drop Flow
```
User drags TaskCard
        │
        ▼
@dnd-kit DndContext → useSortable on card
        │
        ▼
Drop over column → moveTask() or reorderTasks()
        │
        ▼
taskStore updates → _version++ → persist
        │
        ▼
KanbanBoard re-renders with new order
```

---

## Future API Integration

When backend is added (future phase):

```
Frontend (React)          Backend (Go API)
      │                          │
      │  REST /api/tasks  ──────► │
      │                     ◄─────┤
      │                         │
      │  Zustand stores cache   │
      │  localStorage for       │
      │  offline-first          │
```

**Migration path:**
- Replace `taskStore` persistence from localStorage to REST calls
- Add `useApiStore` for server-synced state
- Keep localStorage as offline cache (service worker)
- Update `src/api/` client functions

---

## File Structure

```
test4/
├── index.html              # Entry, Google Fonts, theme flash prevention
├── package.json            # Dependencies
├── tsconfig.json           # Strict TypeScript
├── tailwind.config.js      # Design tokens → CSS vars
├── postcss.config.js
├── vite.config.ts
├── api_test.html           # API endpoint documentation (this project is frontend-only Phase 1)
├── uiflow.md               # This file
├── README.md
│
├── public/
│   └── favicon.svg
│
└── src/
    ├── main.tsx
    ├── App.tsx
    ├── index.css           # Tailwind + CSS custom properties + component classes
    ├── types/
    │   └── index.ts         # Task, Category, Tag, Settings, etc.
    ├── data/
    │   └── defaults.ts      # Default categories, tags, columns
    ├── store/
    │   ├── index.ts         # Barrel export
    │   ├── taskStore.ts     # Tasks (persisted)
    │   ├── uiStore.ts       # UI state (ephemeral)
    │   └── settingsStore.ts # Settings (persisted)
    └── components/          # Built in Phase 2+
        ├── layout/
        ├── board/
        ├── task/
        ├── shared/
        ├── settings/
        └── modals/
```

---

## Design Tokens Reference

| Token | Light | Dark | Usage |
|-------|-------|------|-------|
| `--color-primary` | #3730A3 | #818CF8 | CTAs, links |
| `--color-accent` | #F97316 | #FB923C | Highlights |
| `--color-background` | #FAF8F5 | #0C0A09 | Page bg |
| `--color-surface` | #FFFFFF | #1C1917 | Cards |
| `--color-border` | #E7E5E4 | #292524 | Dividers |
| `--color-text-primary` | #1C1917 | #FAFAF9 | Headings |
| `--color-text-secondary` | #78716C | #A8A29E | Metadata |

Category colors: Research (#4D7C0F), Design (#3730A3), Testing (#D97706), Review (#BE185D), Discovery (#0D9488), Handoff (#475569), Admin (#6B7280)
