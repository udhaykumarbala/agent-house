/**
 * DesignFlow — UI Store
 * Ephemeral UI state. NOT persisted — resets on page load.
 */
import { create } from 'zustand'
import type { ViewMode, ActiveFilters, CategoryType } from '../types'

interface UIStore {
  // View
  view: ViewMode

  // Filters
  activeFilters: ActiveFilters

  // Panels
  selectedTaskId: string | null
  detailPanelOpen: boolean
  settingsOpen: boolean

  // Search
  searchQuery: string

  // Command palette
  commandPaletteOpen: boolean

  // Actions — View
  setView: (view: ViewMode) => void

  // Actions — Filters
  toggleCategoryFilter: (category: CategoryType) => void
  clearFilters: () => void

  // Actions — Panels
  openTaskDetail: (taskId: string) => void
  closeTaskDetail: () => void
  openSettings: () => void
  closeSettings: () => void

  // Actions — Search
  setSearchQuery: (query: string) => void
  clearSearch: () => void

  // Actions — Command Palette
  openCommandPalette: () => void
  closeCommandPalette: () => void
}

export const useUIStore = create<UIStore>((set) => ({
  // Initial state
  view: 'board',
  activeFilters: { categories: [], tags: [] },
  selectedTaskId: null,
  detailPanelOpen: false,
  settingsOpen: false,
  searchQuery: '',
  commandPaletteOpen: false,

  // View
  setView: (view) => set({ view }),

  // Filters
  toggleCategoryFilter: (category) =>
    set(state => {
      const categories = state.activeFilters.categories.includes(category)
        ? state.activeFilters.categories.filter(c => c !== category)
        : [...state.activeFilters.categories, category]
      return { activeFilters: { ...state.activeFilters, categories } }
    }),

  clearFilters: () =>
    set({ activeFilters: { categories: [], tags: [] } }),

  // Panels
  openTaskDetail: (taskId) =>
    set({ selectedTaskId: taskId, detailPanelOpen: true }),

  closeTaskDetail: () =>
    set({ selectedTaskId: null, detailPanelOpen: false }),

  openSettings: () => set({ settingsOpen: true }),

  closeSettings: () => set({ settingsOpen: false }),

  // Search
  setSearchQuery: (query) => set({ searchQuery: query }),

  clearSearch: () => set({ searchQuery: '' }),

  // Command Palette
  openCommandPalette: () => set({ commandPaletteOpen: true }),

  closeCommandPalette: () => set({ commandPaletteOpen: false }),
}))
