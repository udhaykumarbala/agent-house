'use client';

import * as React from 'react';
import { ChevronDown } from 'lucide-react';
import { cn } from '@/lib/utils';

export interface SelectProps
  extends React.SelectHTMLAttributes<HTMLSelectElement> {
  label?: string;
  error?: string;
  options: { value: string; label: string }[];
}

const Select = React.forwardRef<HTMLSelectElement, SelectProps>(
  ({ className, label, error, options, id, ...props }, ref) => {
    const selectId = id || React.useId();

    return (
      <div className="w-full">
        {label && (
          <label
            htmlFor={selectId}
            className="mb-2 block text-sm font-medium text-text-primary dark:text-dark-text-primary"
          >
            {label}
          </label>
        )}
        <div className="relative">
          <select
            id={selectId}
            className={cn(
              'flex h-12 w-full appearance-none rounded-md border border-border bg-surface px-4 py-3 pr-10 text-base text-text-primary transition-colors',
              'focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/20',
              'disabled:cursor-not-allowed disabled:opacity-50',
              'dark:border-dark-border dark:bg-dark-surface dark:text-dark-text-primary dark:focus:border-dark-primary dark:focus:ring-dark-primary/20',
              error && 'border-error focus:border-error focus:ring-error/20 dark:border-dark-error',
              className
            )}
            ref={ref}
            {...props}
          >
            {options.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
          <ChevronDown className="pointer-events-none absolute right-3 top-1/2 h-5 w-5 -translate-y-1/2 text-text-secondary dark:text-dark-text-secondary" />
        </div>
        {error && (
          <p className="mt-2 text-sm text-error dark:text-dark-error">{error}</p>
        )}
      </div>
    );
  }
);
Select.displayName = 'Select';

export { Select };
