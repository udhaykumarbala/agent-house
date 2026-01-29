'use client';

import { useRouter } from 'next/navigation';
import { SubscriptionCard } from '@/components/subscription/subscription-card';
import { SubscriptionCardSkeleton } from '@/components/ui/skeleton';
import { formatCurrency } from '@/lib/utils';
import { useUpcoming } from '@/hooks/use-analytics';

export default function UpcomingPage() {
  const router = useRouter();
  const { data, isLoading } = useUpcoming(30);

  if (isLoading) {
    return (
      <div className="px-4">
        <h1 className="mb-6 text-h2 text-text-primary dark:text-dark-text-primary">
          Upcoming
        </h1>
        <div className="space-y-6">
          {[1, 2, 3].map((i) => (
            <div key={i}>
              <div className="mb-3 h-5 w-24 animate-pulse rounded bg-border dark:bg-dark-border" />
              <div className="space-y-3">
                <SubscriptionCardSkeleton />
                <SubscriptionCardSkeleton />
              </div>
            </div>
          ))}
        </div>
      </div>
    );
  }

  const groups = data?.groups || [];
  const totalUpcoming = data?.total || 0;

  return (
    <div className="px-4">
      <h1 className="mb-6 text-h2 text-text-primary dark:text-dark-text-primary">
        Upcoming
      </h1>

      {groups.length === 0 ? (
        <div className="py-12 text-center">
          <div className="mb-4 text-5xl">📅</div>
          <h2 className="mb-2 text-h3 text-text-primary dark:text-dark-text-primary">
            Nothing coming up
          </h2>
          <p className="text-text-secondary dark:text-dark-text-secondary">
            No renewals in the next 30 days
          </p>
        </div>
      ) : (
        <div className="space-y-6">
          {groups.map((group) => (
            <div key={group.label}>
              <div className="mb-3 flex items-center justify-between">
                <h2 className="text-body-sm font-medium text-text-secondary dark:text-dark-text-secondary">
                  {group.label}
                </h2>
                <span className="text-body-sm font-medium text-text-primary dark:text-dark-text-primary">
                  {formatCurrency(group.total)}
                </span>
              </div>
              <div className="space-y-3">
                {group.subscriptions.map((subscription) => (
                  <SubscriptionCard
                    key={subscription.id}
                    subscription={subscription}
                    onClick={() => router.push(`/subscriptions/${subscription.id}`)}
                  />
                ))}
              </div>
            </div>
          ))}

          <div className="border-t border-border pt-4 dark:border-dark-border">
            <div className="flex items-center justify-between">
              <span className="text-body font-medium text-text-secondary dark:text-dark-text-secondary">
                Total upcoming (30 days)
              </span>
              <span className="text-h3 font-semibold text-text-primary dark:text-dark-text-primary">
                {formatCurrency(totalUpcoming)}
              </span>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
