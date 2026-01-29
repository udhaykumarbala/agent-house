'use client';

import { Plus } from 'lucide-react';
import { cn } from '@/lib/utils';
import { useAppStore } from '@/stores/app-store';

export function FAB() {
  const setAddSheetOpen = useAppStore((state) => state.setAddSheetOpen);

  return (
    <button
      onClick={() => setAddSheetOpen(true)}
      className={cn(
        'fixed bottom-24 right-4 z-30 flex h-14 w-14 items-center justify-center rounded-full bg-primary text-white shadow-fab transition-all',
        'hover:scale-105 hover:shadow-lg active:scale-95',
        'dark:bg-dark-primary'
      )}
      aria-label="Add subscription"
    >
      <Plus className="h-6 w-6" strokeWidth={2.5} />
    </button>
  );
}
