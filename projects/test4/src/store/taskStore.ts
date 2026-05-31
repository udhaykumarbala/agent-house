/**
 * DesignFlow — Task Store
 * Source of truth for all task data. Persisted to localStorage.
 */
import { create } from 'zustand'
import { persist, createJSONStorage } from 'zustand/middleware'
import { nanoid } from 'nanoid'
import type { Task, IterationState, Attachment, Subtask } from '../types'

// ============================================================
// Debounce Helper
// ============================================================

/**
 * Debounced persistence using a _version counter.
 * Zustand's persist middleware uses Object.is() to compare state.
 * By incrementing _version on every write, we guarantee the state object
 * reference changes, which makes the persist middleware fire its save.
 */
let saveTimeout: ReturnType<typeof setTimeout> | null = null

function triggerSave(set: (partial: { _version: number } | ((state: TaskStoreInner) => Partial<TaskStoreInner>)) => void, getVersion: () => number) {
  if (saveTimeout) clearTimeout(saveTimeout)
  saveTimeout = setTimeout(() => {
    // Increment _version to force persist middleware to detect a state change
    set({ _version: getVersion() + 1 })
  }, 300)
}

interface TaskStoreInner {
  tasks: Task[]
  _version: number
  saveStatus: 'idle' | 'saving' | 'saved'
}

interface TaskStore extends TaskStoreInner {
  // CRUD
  addTask: (task: Omit<Task, 'id' | 'createdAt' | 'updatedAt' | 'completedAt'>) => Task
  updateTask: (id: string, updates: Partial<Task>) => void
  deleteTask: (id: string) => void

  // Drag & Drop
  moveTask: (taskId: string, targetState: IterationState) => void
  reorderTasks: (activeId: string, overId: string, newState: IterationState) => void

  // Completion
  toggleTaskComplete: (id: string) => void

  // Subtasks
  addSubtask: (taskId: string, title: string) => void
  toggleSubtask: (taskId: string, subtaskId: string) => void
  deleteSubtask: (taskId: string, subtaskId: string) => void

  // Attachments
  addAttachment: (taskId: string, attachment: Omit<Attachment, 'id'>) => void
  deleteAttachment: (taskId: string, attachmentId: string) => void
}

// ============================================================
// Implementation
// ============================================================

export const useTaskStore = create<TaskStore>()(
  persist<TaskStore>(
    (set, get) => ({
      tasks: [],
      _version: 0,
      saveStatus: 'idle' as const,

      addTask: (taskData) => {
        const now = new Date().toISOString()
        const tasksInColumn = get().tasks.filter(t => t.iterationState === taskData.iterationState)
        const maxPosition = tasksInColumn.length > 0
          ? Math.max(...tasksInColumn.map(t => t.position))
          : -1

        const task: Task = {
          ...taskData,
          id: nanoid(),
          createdAt: now,
          updatedAt: now,
          completedAt: null,
          position: maxPosition + 1,
        }

        set({
          tasks: [...get().tasks, task],
          saveStatus: 'saving' as const,
        })
        triggerSave(set, () => get()._version)
        setTimeout(() => set({ saveStatus: 'saved' as const }), 350)
        setTimeout(() => set({ saveStatus: 'idle' as const }), 1500)

        return task
      },

      updateTask: (id: string, updates: Partial<Task>) => {
        set(state => ({
          tasks: state.tasks.map(t =>
            t.id === id
              ? { ...t, ...updates, updatedAt: new Date().toISOString() }
              : t
          ),
          saveStatus: 'saving' as const,
        }))
        triggerSave(set, () => get()._version)
        setTimeout(() => set({ saveStatus: 'saved' as const }), 350)
        setTimeout(() => set({ saveStatus: 'idle' as const }), 1500)
      },

      deleteTask: (id: string) => {
        set(state => ({
          tasks: state.tasks.filter(t => t.id !== id),
          saveStatus: 'saving' as const,
        }))
        triggerSave(set, () => get()._version)
        setTimeout(() => set({ saveStatus: 'saved' as const }), 350)
        setTimeout(() => set({ saveStatus: 'idle' as const }), 1500)
      },

      moveTask: (taskId: string, targetState: IterationState) => {
        set(state => {
          const tasksInTarget = state.tasks.filter(t => t.iterationState === targetState)
          const maxPosition = tasksInTarget.length > 0
            ? Math.max(...tasksInTarget.map(t => t.position))
            : -1

          return {
            tasks: state.tasks.map(t =>
              t.id === taskId
                ? { ...t, iterationState: targetState, position: maxPosition + 1, updatedAt: new Date().toISOString() }
                : t
            ),
            saveStatus: 'saving' as const,
          }
        })
        triggerSave(set, () => get()._version)
        setTimeout(() => set({ saveStatus: 'saved' as const }), 350)
        setTimeout(() => set({ saveStatus: 'idle' as const }), 1500)
      },

      reorderTasks: (activeId: string, overId: string, newState: IterationState) => {
        set(state => {
          const tasks = [...state.tasks]
          const activeIndex = tasks.findIndex(t => t.id === activeId)
          const overIndex = tasks.findIndex(t => t.id === overId)

          if (activeIndex === -1 || overIndex === -1) return state

          const activeTask = { ...tasks[activeIndex] }
          tasks.splice(activeIndex, 1)

          const targetIndex = tasks.findIndex(t => t.id === overId)
          tasks.splice(targetIndex, 0, activeTask)

          // Re-index positions within the same state
          const filtered = tasks.filter(t => t.iterationState === newState)
          filtered.forEach((t, i) => {
            const idx = tasks.findIndex(x => x.id === t.id)
            tasks[idx] = { ...tasks[idx], position: i, updatedAt: new Date().toISOString() }
          })

          return { tasks, saveStatus: 'saving' as const }
        })
        triggerSave(set, () => get()._version)
        setTimeout(() => set({ saveStatus: 'saved' as const }), 350)
        setTimeout(() => set({ saveStatus: 'idle' as const }), 1500)
      },

      toggleTaskComplete: (id: string) => {
        const task = get().tasks.find(t => t.id === id)
        if (!task) return

        const isCompleting = !task.completedAt
        const completedAt = isCompleting ? new Date().toISOString() : null

        set(state => ({
          tasks: state.tasks.map(t =>
            t.id === id
              ? { ...t, completedAt, iterationState: isCompleting ? 'approved' : t.iterationState, updatedAt: new Date().toISOString() }
              : t
          ),
          saveStatus: 'saving' as const,
        }))
        triggerSave(set, () => get()._version)
        setTimeout(() => set({ saveStatus: 'saved' as const }), 350)
        setTimeout(() => set({ saveStatus: 'idle' as const }), 1500)
      },

      addSubtask: (taskId: string, title: string) => {
        const subtask: Subtask = { id: nanoid(), title, completed: false }
        set(state => ({
          tasks: state.tasks.map(t =>
            t.id === taskId
              ? { ...t, subtasks: [...t.subtasks, subtask], updatedAt: new Date().toISOString() }
              : t
          ),
          saveStatus: 'saving' as const,
        }))
        triggerSave(set, () => get()._version)
        setTimeout(() => set({ saveStatus: 'saved' as const }), 350)
        setTimeout(() => set({ saveStatus: 'idle' as const }), 1500)
      },

      toggleSubtask: (taskId: string, subtaskId: string) => {
        set(state => ({
          tasks: state.tasks.map(t =>
            t.id === taskId
              ? {
                  ...t,
                  subtasks: t.subtasks.map(s =>
                    s.id === subtaskId ? { ...s, completed: !s.completed } : s
                  ),
                  updatedAt: new Date().toISOString(),
                }
              : t
          ),
          saveStatus: 'saving' as const,
        }))
        triggerSave(set, () => get()._version)
        setTimeout(() => set({ saveStatus: 'saved' as const }), 350)
        setTimeout(() => set({ saveStatus: 'idle' as const }), 1500)
      },

      deleteSubtask: (taskId: string, subtaskId: string) => {
        set(state => ({
          tasks: state.tasks.map(t =>
            t.id === taskId
              ? {
                  ...t,
                  subtasks: t.subtasks.filter(s => s.id !== subtaskId),
                  updatedAt: new Date().toISOString(),
                }
              : t
          ),
          saveStatus: 'saving' as const,
        }))
        triggerSave(set, () => get()._version)
        setTimeout(() => set({ saveStatus: 'saved' as const }), 350)
        setTimeout(() => set({ saveStatus: 'idle' as const }), 1500)
      },

      addAttachment: (taskId: string, attachment: Omit<Attachment, 'id'>) => {
        const newAttachment: Attachment = { ...attachment, id: nanoid() }
        set(state => ({
          tasks: state.tasks.map(t =>
            t.id === taskId
              ? { ...t, attachments: [...t.attachments, newAttachment], updatedAt: new Date().toISOString() }
              : t
          ),
          saveStatus: 'saving' as const,
        }))
        triggerSave(set, () => get()._version)
        setTimeout(() => set({ saveStatus: 'saved' as const }), 350)
        setTimeout(() => set({ saveStatus: 'idle' as const }), 1500)
      },

      deleteAttachment: (taskId: string, attachmentId: string) => {
        set(state => ({
          tasks: state.tasks.map(t =>
            t.id === taskId
              ? {
                  ...t,
                  attachments: t.attachments.filter(a => a.id !== attachmentId),
                  updatedAt: new Date().toISOString(),
                }
              : t
          ),
          saveStatus: 'saving' as const,
        }))
        triggerSave(set, () => get()._version)
        setTimeout(() => set({ saveStatus: 'saved' as const }), 350)
        setTimeout(() => set({ saveStatus: 'idle' as const }), 1500)
      },
    }),
    {
      name: 'designflow_tasks',
      storage: createJSONStorage(() => localStorage),
    }
  )
)

// ============================================================
// Selectors (computed values)
// ============================================================

export const selectTasksByState = (state: IterationState) =>
  useTaskStore.getState().tasks
    .filter(t => t.iterationState === state)
    .sort((a, b) => a.position - b.position)

export const selectTaskById = (id: string) =>
  useTaskStore.getState().tasks.find(t => t.id === id)

export const selectTotalTasks = () => useTaskStore.getState().tasks.length

export const selectInProgressCount = () =>
  useTaskStore.getState().tasks.filter(t => t.iterationState === 'in_progress').length
