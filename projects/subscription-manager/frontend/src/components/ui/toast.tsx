'use client';

import * as React from 'react';
import { CheckCircle, XCircle, Info, X } from 'lucide-react';
import { cn } from '@/lib/utils';
import { useAppStore } from '@/stores/app-store';

const icons = {
  success: CheckCircle,
  error: XCircle,
  info: Info,
};

export function ToastContainer() {
  const toasts = useAppStore((state) => state.toasts);
  const removeToast = useAppStore((state) => state.removeToast);

  if (toasts.length === 0) return null;

  return (
    <div className="fixed bottom-20 left-0 right-0 z-50 flex flex-col items-center gap-2 px-4 md:bottom-6">
      {toasts.map((toast) => {
        const Icon = icons[toast.type];
        return (
          <div
            key={toast.id}
            className={cn(
              'flex w-full max-w-sm items-center gap-3 rounded-lg px-4 py-3 shadow-lg animate-fade-in',
              'bg-text-primary text-background',
              'dark:bg-dark-text-primary dark:text-dark-background'
            )}
          >
            <Icon className={cn(
              'h-5 w-5 flex-shrink-0',
              toast.type === 'success' && 'text-success dark:text-dark-success',
              toast.type === 'error' && 'text-error dark:text-dark-error',
              toast.type === 'info' && 'text-info'
            )} />
            <p className="flex-1 text-sm">{toast.message}</p>
            {toast.action && (
              <button
                onClick={toast.action.onClick}
                className="text-sm font-medium text-primary hover:underline dark:text-dark-primary"
              >
                {toast.action.label}
              </button>
            )}
            <button
              onClick={() => removeToast(toast.id)}
              className="flex-shrink-0 rounded p-1 transition-colors hover:bg-white/10"
            >
              <X className="h-4 w-4" />
            </button>
          </div>
        );
      })}
    </div>
  );
}
