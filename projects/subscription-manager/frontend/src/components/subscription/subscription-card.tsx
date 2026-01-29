'use client';

import { AlertCircle, ChevronRight } from 'lucide-react';
import { cn } from '@/lib/utils';
import { formatCurrency } from '@/lib/utils';
import { CATEGORY_COLORS } from '@/lib/constants';
import type { Subscription } from '@/types/subscription';
import {
  formatBillingCycle,
  formatDate,
  getDaysUntilBilling,
  isDueSoon,
} from '@/types/subscription';

interface SubscriptionCardProps {
  subscription: Subscription;
  onClick?: () => void;
}

export function SubscriptionCard({ subscription, onClick }: SubscriptionCardProps) {
  const daysUntil = getDaysUntilBilling(subscription.nextBillingDate);
  const dueSoon = isDueSoon(subscription.nextBillingDate);

  const getDueDateText = () => {
    if (daysUntil === 0) return 'Due today';
    if (daysUntil === 1) return 'Due tomorrow';
    if (daysUntil < 0) return 'Overdue';
    if (dueSoon) return `Due in ${daysUntil}d`;
    return `Due ${formatDate(subscription.nextBillingDate)}`;
  };

  return (
    <button
      onClick={onClick}
      className={cn(
        'flex w-full items-center gap-3 rounded-lg bg-surface p-4 text-left shadow-md transition-all',
        'hover:shadow-lg hover:-translate-y-0.5 active:scale-[0.98]',
        'dark:bg-dark-surface dark:shadow-none dark:border dark:border-dark-border/50 dark:hover:bg-dark-surface-hover'
      )}
    >
      {/* Icon */}
      <div
        className="flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-lg text-2xl"
        style={{
          backgroundColor: subscription.color
            ? `${subscription.color}20`
            : `${CATEGORY_COLORS[subscription.category]}20`,
        }}
      >
        {subscription.icon || '📦'}
      </div>

      {/* Content */}
      <div className="min-w-0 flex-1">
        <div className="flex items-center gap-2">
          <h3 className="truncate font-semibold text-text-primary dark:text-dark-text-primary">
            {subscription.name}
          </h3>
          {dueSoon && (
            <AlertCircle className="h-4 w-4 flex-shrink-0 text-warning dark:text-dark-warning" />
          )}
        </div>
        <p className="text-sm text-text-secondary dark:text-dark-text-secondary">
          {subscription.category.charAt(0).toUpperCase() + subscription.category.slice(1)}
        </p>
      </div>

      {/* Price and date */}
      <div className="text-right">
        <p className="font-semibold tabular-nums text-text-primary dark:text-dark-text-primary">
          {formatCurrency(subscription.price, subscription.currency)}
          <span className="text-sm font-normal text-text-secondary dark:text-dark-text-secondary">
            {formatBillingCycle(subscription.billingCycle)}
          </span>
        </p>
        <p
          className={cn(
            'text-sm',
            dueSoon
              ? 'text-warning dark:text-dark-warning'
              : 'text-text-secondary dark:text-dark-text-secondary'
          )}
        >
          {getDueDateText()}
        </p>
      </div>

      {/* Chevron */}
      <ChevronRight className="h-5 w-5 flex-shrink-0 text-text-tertiary dark:text-dark-text-tertiary" />
    </button>
  );
}
