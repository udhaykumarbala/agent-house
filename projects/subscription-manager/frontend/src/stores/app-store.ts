import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { User } from '@/types/user';

interface Toast {
  id: string;
  message: string;
  type: 'success' | 'error' | 'info';
  action?: {
    label: string;
    onClick: () => void;
  };
}

interface AppState {
  // User
  user: User | null;
  setUser: (user: User | null) => void;

  // Theme
  theme: 'light' | 'dark' | 'system';
  setTheme: (theme: 'light' | 'dark' | 'system') => void;

  // Sheet state
  isAddSheetOpen: boolean;
  setAddSheetOpen: (open: boolean) => void;

  // Edit subscription
  editingSubscriptionId: string | null;
  setEditingSubscription: (id: string | null) => void;

  // Toast notifications
  toasts: Toast[];
  addToast: (toast: Omit<Toast, 'id'>) => void;
  removeToast: (id: string) => void;

  // Loading states
  isLoading: boolean;
  setLoading: (loading: boolean) => void;
}

export const useAppStore = create<AppState>()(
  persist(
    (set) => ({
      // User
      user: null,
      setUser: (user) => set({ user }),

      // Theme
      theme: 'system',
      setTheme: (theme) => set({ theme }),

      // Sheet state
      isAddSheetOpen: false,
      setAddSheetOpen: (open) => set({ isAddSheetOpen: open }),

      // Edit subscription
      editingSubscriptionId: null,
      setEditingSubscription: (id) => set({ editingSubscriptionId: id }),

      // Toast notifications
      toasts: [],
      addToast: (toast) => {
        const id = Math.random().toString(36).slice(2);
        set((state) => ({
          toasts: [...state.toasts, { ...toast, id }],
        }));
        setTimeout(() => {
          set((state) => ({
            toasts: state.toasts.filter((t) => t.id !== id),
          }));
        }, toast.action ? 5000 : 3000);
      },
      removeToast: (id) =>
        set((state) => ({
          toasts: state.toasts.filter((t) => t.id !== id),
        })),

      // Loading states
      isLoading: false,
      setLoading: (loading) => set({ isLoading: loading }),
    }),
    {
      name: 'subtrack-storage',
      partialize: (state) => ({
        theme: state.theme,
      }),
    }
  )
);
