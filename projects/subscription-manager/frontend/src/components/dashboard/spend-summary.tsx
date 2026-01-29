'use client';

import * as React from 'react';
import { formatCurrency } from '@/lib/utils';
import { Skeleton } from '@/components/ui/skeleton';
import type { Summary } from '@/types/api';

interface SpendSummaryProps {
  summary: Summary | undefined;
  isLoading: boolean;
}

export function SpendSummary({ summary, isLoading }: SpendSummaryProps) {
  const [viewMode, setViewMode] = React.useState<'monthly' | 'yearly'>('monthly');

  if (isLoading) {
    return (
      <div className="py-6 text-center">
        <Skeleton className="mx-auto mb-2 h-12 w-48" />
        <Skeleton className="mx-auto h-5 w-32" />
      </div>
    );
  }

  const amount = viewMode === 'monthly' ? summary?.totalMonthly : summary?.totalYearly;
  const label = viewMode === 'monthly' ? '/month' : '/year';

  return (
    <div className="py-6 text-center">
      <div className="mb-2">
        <span className="text-display tabular-nums text-text-primary dark:text-dark-text-primary">
          {formatCurrency(amount || 0)}
        </span>
        <span className="text-h3 text-text-secondary dark:text-dark-text-secondary">
          {label}
        </span>
      </div>

      <div className="flex items-center justify-center gap-2">
        <button
          onClick={() => setViewMode('monthly')}
          className={`rounded-full px-3 py-1 text-sm font-medium transition-colors ${
            viewMode === 'monthly'
              ? 'bg-primary/10 text-primary dark:bg-dark-primary/10 dark:text-dark-primary'
              : 'text-text-secondary hover:text-text-primary dark:text-dark-text-secondary dark:hover:text-dark-text-primary'
          }`}
        >
          Monthly
        </button>
        <span className="text-text-tertiary dark:text-dark-text-tertiary">|</span>
        <button
          onClick={() => setViewMode('yearly')}
          className={`rounded-full px-3 py-1 text-sm font-medium transition-colors ${
            viewMode === 'yearly'
              ? 'bg-primary/10 text-primary dark:bg-dark-primary/10 dark:text-dark-primary'
              : 'text-text-secondary hover:text-text-primary dark:text-dark-text-secondary dark:hover:text-dark-text-primary'
          }`}
        >
          Yearly
        </button>
      </div>

      {summary && summary.dueSoon > 0 && (
        <p className="mt-3 text-sm text-warning dark:text-dark-warning">
          {summary.dueSoon} subscription{summary.dueSoon > 1 ? 's' : ''} due soon (
          {formatCurrency(summary.dueSoonAmount)})
        </p>
      )}
    </div>
  );
}
