'use client';

import { useRouter } from 'next/navigation';
import { SubscriptionCard } from './subscription-card';
import { SubscriptionCardSkeleton } from '@/components/ui/skeleton';
import type { Subscription } from '@/types/subscription';

interface SubscriptionListProps {
  subscriptions: Subscription[] | undefined;
  isLoading: boolean;
}

export function SubscriptionList({ subscriptions, isLoading }: SubscriptionListProps) {
  const router = useRouter();

  if (isLoading) {
    return (
      <div className="space-y-3">
        {[1, 2, 3, 4].map((i) => (
          <SubscriptionCardSkeleton key={i} />
        ))}
      </div>
    );
  }

  if (!subscriptions || subscriptions.length === 0) {
    return null;
  }

  return (
    <div className="space-y-3">
      {subscriptions.map((subscription) => (
        <SubscriptionCard
          key={subscription.id}
          subscription={subscription}
          onClick={() => router.push(`/subscriptions/${subscription.id}`)}
        />
      ))}
    </div>
  );
}
