'use client';

import * as React from 'react';
import { X } from 'lucide-react';
import { cn } from '@/lib/utils';

interface SheetProps {
  open: boolean;
  onClose: () => void;
  title?: string;
  children: React.ReactNode;
}

export function Sheet({ open, onClose, title, children }: SheetProps) {
  React.useEffect(() => {
    if (open) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = '';
    }
    return () => {
      document.body.style.overflow = '';
    };
  }, [open]);

  if (!open) return null;

  return (
    <>
      {/* Backdrop */}
      <div
        className="fixed inset-0 z-50 bg-black/50 transition-opacity animate-fade-in"
        onClick={onClose}
      />

      {/* Sheet */}
      <div
        className={cn(
          'fixed inset-x-0 bottom-0 z-50 max-h-[90vh] overflow-auto rounded-t-2xl bg-surface shadow-xl animate-slide-up',
          'dark:bg-dark-surface'
        )}
      >
        {/* Drag handle */}
        <div className="flex justify-center py-3">
          <div className="h-1 w-9 rounded-full bg-border dark:bg-dark-border" />
        </div>

        {/* Header */}
        {title && (
          <div className="flex items-center justify-between border-b border-border px-6 pb-4 dark:border-dark-border">
            <h2 className="text-h3 text-text-primary dark:text-dark-text-primary">
              {title}
            </h2>
            <button
              onClick={onClose}
              className="rounded-full p-2 text-text-secondary transition-colors hover:bg-surface-hover hover:text-text-primary dark:text-dark-text-secondary dark:hover:bg-dark-surface-hover dark:hover:text-dark-text-primary"
            >
              <X className="h-5 w-5" />
            </button>
          </div>
        )}

        {/* Content */}
        <div className="px-6 py-4 pb-8">{children}</div>
      </div>
    </>
  );
}
