/**
 * DesignFlow — TypeScript Type Definitions
 * Source: .plans/specs/architecture-spec.md Section 6.3
 */

// ============================================================
// Enums / Union Types
// ============================================================

/** Kanban column / workflow stage */
export type IterationState = 'backlog' | 'in_progress' | 'in_review' | 'approved';

/** 1–5 priority scale (1 = highest) */
export type Priority = 1 | 2 | 3 | 4 | 5;

/** UX design phase categories */
export type CategoryType =
  | 'research'
  | 'design'
  | 'testing'
  | 'review'
  | 'discovery'
  | 'handoff'
  | 'admin';

/** App view modes */
export type ViewMode = 'board' | 'list' | 'focus';

/** Theme preference */
export type ThemeMode = 'light' | 'dark' | 'system';

// ============================================================
// Interfaces
// ============================================================

export interface Attachment {
  id: string;
  type: 'url';
  url: string;
  label?: string;
  favicon?: string;
}

export interface Subtask {
  id: string;
  title: string;
  completed: boolean;
}

export interface Task {
  id: string;
  title: string;
  description: string;
  category: CategoryType;
  iterationState: IterationState;
  priority: Priority;
  dueDate: string | null;       // ISO date string e.g. "2026-03-25"
  timeEstimate: number | null;  // minutes
  tags: string[];
  attachments: Attachment[];
  subtasks: Subtask[];
  createdAt: string;             // ISO timestamp
  updatedAt: string;             // ISO timestamp
  completedAt: string | null;    // ISO timestamp
  position: number;              // sort order within column
}

export interface Category {
  id: string;
  name: string;
  colorLight: string;
  colorDark: string;
  icon?: string;                 // Lucide icon name
  keywords: string[];           // for natural language auto-categorization
}

export interface Tag {
  id: string;
  name: string;
  color: string;
}

export interface Settings {
  theme: ThemeMode;
  columnOrder: IterationState[];
  customCategories: Category[];
  customTags: Tag[];
}

export interface ActiveFilters {
  categories: CategoryType[];
  tags: string[];
}

// ============================================================
// Column / Status Config
// ============================================================

export interface ColumnConfig {
  id: IterationState;
  label: string;
}
