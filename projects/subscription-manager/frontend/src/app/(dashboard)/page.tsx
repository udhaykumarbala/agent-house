'use client';

import * as React from 'react';
import { SpendSummary } from '@/components/dashboard/spend-summary';
import { EmptyState } from '@/components/dashboard/empty-state';
import { SubscriptionList } from '@/components/subscription/subscription-list';
import { CategoryFilter } from '@/components/subscription/category-filter';
import { DashboardSkeleton } from '@/components/ui/skeleton';
import { useSubscriptions } from '@/hooks/use-subscriptions';
import { useSummary } from '@/hooks/use-analytics';
import type { Category } from '@/types/subscription';

export default function DashboardPage() {
  const [selectedCategory, setSelectedCategory] = React.useState<Category | null>(null);
  const { data: subscriptions, isLoading: subsLoading } = useSubscriptions();
  const { data: summary, isLoading: summaryLoading } = useSummary();

  const isLoading = subsLoading || summaryLoading;

  const filteredSubscriptions = React.useMemo(() => {
    if (!subscriptions) return [];
    if (!selectedCategory) return subscriptions;
    return subscriptions.filter((sub) => sub.category === selectedCategory);
  }, [subscriptions, selectedCategory]);

  if (isLoading) {
    return (
      <div className="px-4">
        <DashboardSkeleton />
      </div>
    );
  }

  const isEmpty = !subscriptions || subscriptions.length === 0;

  return (
    <div className="px-4">
      {isEmpty ? (
        <EmptyState />
      ) : (
        <>
          <SpendSummary summary={summary} isLoading={summaryLoading} />

          <div className="mb-4">
            <CategoryFilter
              selected={selectedCategory}
              onChange={setSelectedCategory}
            />
          </div>

          <SubscriptionList
            subscriptions={filteredSubscriptions}
            isLoading={subsLoading}
          />

          {filteredSubscriptions.length === 0 && selectedCategory && (
            <p className="py-8 text-center text-text-secondary dark:text-dark-text-secondary">
              No subscriptions in this category
            </p>
          )}
        </>
      )}
    </div>
  );
}
