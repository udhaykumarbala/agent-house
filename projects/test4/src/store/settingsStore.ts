/**
 * DesignFlow — Settings Store
 * User preferences. Persisted to localStorage.
 */
import { create } from 'zustand'
import { persist, createJSONStorage } from 'zustand/middleware'
import type { Settings, ThemeMode, Category, Tag, IterationState } from '../types'
import { DEFAULT_CATEGORIES, DEFAULT_TAGS, DEFAULT_COLUMN_ORDER } from '../data/defaults'

interface SettingsStore extends Settings {
  updateTheme: (theme: ThemeMode) => void
  updateColumnOrder: (order: IterationState[]) => void
  updateCategories: (categories: Category[]) => void
  updateTags: (tags: Tag[]) => void
  addCategory: (category: Category) => void
  removeCategory: (id: string) => void
  addTag: (tag: Tag) => void
  removeTag: (id: string) => void
}

export const useSettingsStore = create<SettingsStore>()(
  persist(
    (set) => ({
      theme: 'system',
      columnOrder: DEFAULT_COLUMN_ORDER,
      customCategories: DEFAULT_CATEGORIES,
      customTags: DEFAULT_TAGS,

      updateTheme: (theme) => {
        set({ theme })
        applyTheme(theme)
      },

      updateColumnOrder: (order) => set({ columnOrder: order }),

      updateCategories: (categories) => set({ customCategories: categories }),

      updateTags: (tags) => set({ customTags: tags }),

      addCategory: (category) =>
        set(state => ({
          customCategories: [...state.customCategories, category],
        })),

      removeCategory: (id) =>
        set(state => ({
          customCategories: state.customCategories.filter(c => c.id !== id),
        })),

      addTag: (tag) =>
        set(state => ({
          customTags: [...state.customTags, tag],
        })),

      removeTag: (id) =>
        set(state => ({
          customTags: state.customTags.filter(t => t.id !== id),
        })),
    }),
    {
      name: 'designflow_settings',
      storage: createJSONStorage(() => localStorage),
      onRehydrateStorage: () => (state) => {
        // Apply saved theme on page load (before first paint via inline script in index.html)
        if (state) {
          applyTheme(state.theme)
        }
      },
    }
  )
)

// ============================================================
// Theme Application Helper
// ============================================================

export function applyTheme(theme: ThemeMode) {
  const root = document.documentElement

  switch (theme) {
    case 'dark':
      root.classList.add('dark')
      break
    case 'light':
      root.classList.remove('dark')
      break
    case 'system':
      if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
        root.classList.add('dark')
      } else {
        root.classList.remove('dark')
      }
      break
  }
}

// Listen for system theme changes
if (typeof window !== 'undefined') {
  window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
    const state = useSettingsStore.getState()
    if (state.theme === 'system') {
      applyTheme('system')
    }
  })
}
