/**
 * DesignFlow — Default Data
 * Categories, tags, and column configurations.
 * Colors from: .plans/specs/ui-spec.md (Category Colors section)
 */

import type { Category, Tag, ColumnConfig, IterationState } from '../types';

// ============================================================
// Column Config (kanban stages)
// ============================================================

export const DEFAULT_COLUMNS: ColumnConfig[] = [
  { id: 'backlog',     label: 'Backlog' },
  { id: 'in_progress', label: 'In Progress' },
  { id: 'in_review',   label: 'In Review' },
  { id: 'approved',    label: 'Approved' },
];

export const DEFAULT_COLUMN_ORDER: IterationState[] = [
  'backlog',
  'in_progress',
  'in_review',
  'approved',
];

// ============================================================
// Categories (UX design process phases)
// ============================================================

export const DEFAULT_CATEGORIES: Category[] = [
  {
    id: 'research',
    name: 'Research',
    colorLight: '#4D7C0F',
    colorDark:  '#84CC16',
    keywords:   ['#research', '#r'],
  },
  {
    id: 'design',
    name: 'Design',
    colorLight: '#3730A3',
    colorDark:  '#818CF8',
    keywords:   ['#design', '#d'],
  },
  {
    id: 'testing',
    name: 'Testing',
    colorLight: '#D97706',
    colorDark:  '#FBBF24',
    keywords:   ['#testing', '#test', '#t'],
  },
  {
    id: 'review',
    name: 'Review',
    colorLight: '#BE185D',
    colorDark:  '#F472B6',
    keywords:   ['#review', '#rev', '#rv'],
  },
  {
    id: 'discovery',
    name: 'Discovery',
    colorLight: '#0D9488',
    colorDark:  '#2DD4BF',
    keywords:   ['#discovery'],
  },
  {
    id: 'handoff',
    name: 'Handoff',
    colorLight: '#475569',
    colorDark:  '#94A3B8',
    keywords:   ['#handoff', '#handover'],
  },
  {
    id: 'admin',
    name: 'Admin',
    colorLight: '#6B7280',
    colorDark:  '#9CA3AF',
    keywords:   ['#admin', '#a'],
  },
];

// ============================================================
// Default Tags
// ============================================================

export const DEFAULT_TAGS: Tag[] = [
  { id: 'tag-1', name: 'Urgent',    color: '#DC2626' },
  { id: 'tag-2', name: 'High Priority', color: '#D97706' },
  { id: 'tag-3', name: 'Blocked',   color: '#7C3AED' },
  { id: 'tag-4', name: 'In Sprint', color: '#0D9488' },
  { id: 'tag-5', name: 'Tech Debt', color: '#475569' },
  { id: 'tag-6', name: 'Design QA', color: '#BE185D' },
  { id: 'tag-7', name: 'Stakeholder', color: '#3730A3' },
];

// ============================================================
// Iteration State → Column Label helper
// ============================================================

export const ITERATION_LABELS: Record<IterationState, string> = {
  backlog:     'Backlog',
  in_progress: 'In Progress',
  in_review:   'In Review',
  approved:    'Approved',
};
