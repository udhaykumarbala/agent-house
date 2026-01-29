'use client';

import { Plus } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { POPULAR_TEMPLATES } from '@/lib/constants';
import { useAppStore } from '@/stores/app-store';

export function EmptyState() {
  const setAddSheetOpen = useAppStore((state) => state.setAddSheetOpen);

  return (
    <div className="flex flex-col items-center justify-center py-12 text-center">
      {/* Illustration */}
      <div className="mb-6 flex h-24 w-24 items-center justify-center rounded-full bg-primary/10 dark:bg-dark-primary/10">
        <span className="text-5xl">📋</span>
      </div>

      <h2 className="mb-2 text-h2 text-text-primary dark:text-dark-text-primary">
        No subscriptions yet
      </h2>
      <p className="mb-6 max-w-xs text-text-secondary dark:text-dark-text-secondary">
        Add your first subscription to start tracking your spending.
      </p>

      <Button onClick={() => setAddSheetOpen(true)} className="mb-8">
        <Plus className="mr-2 h-5 w-5" />
        Add Subscription
      </Button>

      <p className="mb-3 text-sm text-text-secondary dark:text-dark-text-secondary">
        Or choose from popular:
      </p>

      <div className="flex flex-wrap justify-center gap-2">
        {POPULAR_TEMPLATES.slice(0, 4).map((template) => (
          <button
            key={template.name}
            onClick={() => setAddSheetOpen(true)}
            className="flex items-center gap-2 rounded-lg border border-border px-3 py-2 transition-colors hover:bg-surface-hover dark:border-dark-border dark:hover:bg-dark-surface-hover"
          >
            <span className="text-xl">{template.icon}</span>
            <span className="text-sm text-text-primary dark:text-dark-text-primary">
              {template.name}
            </span>
          </button>
        ))}
      </div>
    </div>
  );
}
