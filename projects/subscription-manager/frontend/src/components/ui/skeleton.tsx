'use client';

import { cn } from '@/lib/utils';

interface SkeletonProps {
  className?: string;
}

export function Skeleton({ className }: SkeletonProps) {
  return (
    <div
      className={cn(
        'animate-pulse-slow rounded-md bg-border dark:bg-dark-border',
        className
      )}
    />
  );
}

export function SubscriptionCardSkeleton() {
  return (
    <div className="flex items-center gap-3 rounded-lg bg-surface p-4 shadow-md dark:bg-dark-surface dark:shadow-none dark:border dark:border-dark-border/50">
      <Skeleton className="h-12 w-12 rounded-lg" />
      <div className="flex-1">
        <Skeleton className="mb-2 h-4 w-24" />
        <Skeleton className="h-3 w-32" />
      </div>
      <div className="text-right">
        <Skeleton className="mb-2 h-5 w-16" />
        <Skeleton className="h-3 w-20" />
      </div>
    </div>
  );
}

export function DashboardSkeleton() {
  return (
    <div className="space-y-6">
      {/* Spend summary skeleton */}
      <div className="text-center py-6">
        <Skeleton className="mx-auto mb-2 h-10 w-40" />
        <Skeleton className="mx-auto h-5 w-32" />
      </div>

      {/* Filter pills skeleton */}
      <div className="flex gap-2">
        <Skeleton className="h-9 w-20 rounded-full" />
        <Skeleton className="h-9 w-24 rounded-full" />
      </div>

      {/* Cards skeleton */}
      <div className="space-y-3">
        {[1, 2, 3].map((i) => (
          <SubscriptionCardSkeleton key={i} />
        ))}
      </div>
    </div>
  );
}
