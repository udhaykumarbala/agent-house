'use client';

import { cn } from '@/lib/utils';
import { CATEGORIES, type Category } from '@/types/subscription';

interface CategoryFilterProps {
  selected: Category | null;
  onChange: (category: Category | null) => void;
}

export function CategoryFilter({ selected, onChange }: CategoryFilterProps) {
  return (
    <div className="flex gap-2 overflow-x-auto pb-2 scrollbar-hide">
      <button
        onClick={() => onChange(null)}
        className={cn(
          'flex-shrink-0 rounded-full px-4 py-2 text-sm font-medium transition-colors',
          selected === null
            ? 'bg-primary text-white dark:bg-dark-primary'
            : 'border border-border text-text-secondary hover:bg-surface-hover dark:border-dark-border dark:text-dark-text-secondary dark:hover:bg-dark-surface-hover'
        )}
      >
        All
      </button>
      {CATEGORIES.map((category) => (
        <button
          key={category.value}
          onClick={() => onChange(category.value)}
          className={cn(
            'flex-shrink-0 rounded-full px-4 py-2 text-sm font-medium transition-colors',
            selected === category.value
              ? 'bg-primary text-white dark:bg-dark-primary'
              : 'border border-border text-text-secondary hover:bg-surface-hover dark:border-dark-border dark:text-dark-text-secondary dark:hover:bg-dark-surface-hover'
          )}
        >
          {category.icon} {category.label}
        </button>
      ))}
    </div>
  );
}
