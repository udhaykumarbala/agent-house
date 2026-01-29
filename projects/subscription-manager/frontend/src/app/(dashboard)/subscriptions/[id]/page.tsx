'use client';

import * as React from 'react';
import { useRouter, useParams } from 'next/navigation';
import { ChevronLeft, Trash2, Pencil } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Sheet } from '@/components/ui/sheet';
import { SubscriptionForm } from '@/components/subscription/subscription-form';
import { Skeleton } from '@/components/ui/skeleton';
import { formatCurrency } from '@/lib/utils';
import { CATEGORY_COLORS } from '@/lib/constants';
import {
  useSubscription,
  useUpdateSubscription,
  useDeleteSubscription,
} from '@/hooks/use-subscriptions';
import {
  formatBillingCycle,
  formatFullDate,
  getDaysUntilBilling,
  CATEGORIES,
  type SubscriptionInput,
} from '@/types/subscription';

export default function SubscriptionDetailPage() {
  const router = useRouter();
  const params = useParams();
  const id = params.id as string;

  const [isEditOpen, setIsEditOpen] = React.useState(false);
  const [showDeleteConfirm, setShowDeleteConfirm] = React.useState(false);

  const { data: subscription, isLoading } = useSubscription(id);
  const updateSubscription = useUpdateSubscription();
  const deleteSubscription = useDeleteSubscription();

  const handleUpdate = (data: SubscriptionInput) => {
    updateSubscription.mutate(
      { id, input: data },
      {
        onSuccess: () => {
          setIsEditOpen(false);
        },
      }
    );
  };

  const handleDelete = () => {
    deleteSubscription.mutate(id, {
      onSuccess: () => {
        router.push('/');
      },
    });
  };

  if (isLoading) {
    return (
      <div className="px-4">
        <button
          onClick={() => router.back()}
          className="mb-4 flex items-center gap-1 text-text-secondary dark:text-dark-text-secondary"
        >
          <ChevronLeft className="h-5 w-5" />
          Back
        </button>
        <div className="flex flex-col items-center py-8">
          <Skeleton className="mb-4 h-20 w-20 rounded-xl" />
          <Skeleton className="mb-2 h-8 w-32" />
          <Skeleton className="h-6 w-24" />
        </div>
      </div>
    );
  }

  if (!subscription) {
    return (
      <div className="px-4 py-12 text-center">
        <p className="text-text-secondary dark:text-dark-text-secondary">
          Subscription not found
        </p>
        <Button variant="ghost" onClick={() => router.push('/')} className="mt-4">
          Go back home
        </Button>
      </div>
    );
  }

  const daysUntil = getDaysUntilBilling(subscription.nextBillingDate);
  const categoryInfo = CATEGORIES.find((c) => c.value === subscription.category);

  return (
    <div className="px-4">
      {/* Header */}
      <div className="mb-6 flex items-center justify-between">
        <button
          onClick={() => router.back()}
          className="flex items-center gap-1 text-text-secondary hover:text-text-primary dark:text-dark-text-secondary dark:hover:text-dark-text-primary"
        >
          <ChevronLeft className="h-5 w-5" />
          Back
        </button>
        <button
          onClick={() => setShowDeleteConfirm(true)}
          className="rounded-full p-2 text-error transition-colors hover:bg-error-bg dark:text-dark-error dark:hover:bg-dark-error/10"
        >
          <Trash2 className="h-5 w-5" />
        </button>
      </div>

      {/* Main info */}
      <div className="mb-8 flex flex-col items-center text-center">
        <div
          className="mb-4 flex h-20 w-20 items-center justify-center rounded-xl text-4xl"
          style={{
            backgroundColor: subscription.color
              ? `${subscription.color}20`
              : `${CATEGORY_COLORS[subscription.category]}20`,
          }}
        >
          {subscription.icon || '📦'}
        </div>
        <h1 className="mb-2 text-h1 text-text-primary dark:text-dark-text-primary">
          {subscription.name}
        </h1>
        <p className="text-h2 tabular-nums text-text-primary dark:text-dark-text-primary">
          {formatCurrency(subscription.price, subscription.currency)}
          <span className="text-body-lg text-text-secondary dark:text-dark-text-secondary">
            {formatBillingCycle(subscription.billingCycle)}
          </span>
        </p>
      </div>

      {/* Details */}
      <div className="space-y-4">
        <DetailRow
          label="Next billing"
          value={formatFullDate(subscription.nextBillingDate)}
          subValue={
            daysUntil === 0
              ? 'Due today'
              : daysUntil === 1
              ? 'Due tomorrow'
              : daysUntil < 0
              ? 'Overdue'
              : `In ${daysUntil} days`
          }
          onEdit={() => setIsEditOpen(true)}
        />

        <DetailRow
          label="Billing cycle"
          value={
            subscription.billingCycle.charAt(0).toUpperCase() +
            subscription.billingCycle.slice(1)
          }
          onEdit={() => setIsEditOpen(true)}
        />

        <DetailRow
          label="Category"
          value={`${categoryInfo?.icon || ''} ${categoryInfo?.label || subscription.category}`}
          onEdit={() => setIsEditOpen(true)}
        />

        {subscription.notes && (
          <div className="rounded-lg bg-surface p-4 shadow-sm dark:bg-dark-surface dark:border dark:border-dark-border/50">
            <p className="mb-1 text-label text-text-secondary dark:text-dark-text-secondary">
              Notes
            </p>
            <p className="text-body text-text-primary dark:text-dark-text-primary">
              {subscription.notes}
            </p>
          </div>
        )}
      </div>

      {/* Metadata */}
      <div className="mt-8 border-t border-border pt-4 dark:border-dark-border">
        <p className="text-body-sm text-text-tertiary dark:text-dark-text-tertiary">
          Added {new Date(subscription.createdAt).toLocaleDateString()}
        </p>
      </div>

      {/* Edit Sheet */}
      <Sheet
        open={isEditOpen}
        onClose={() => setIsEditOpen(false)}
        title="Edit Subscription"
      >
        <SubscriptionForm
          initialData={{
            name: subscription.name,
            price: subscription.price,
            billingCycle: subscription.billingCycle,
            nextBillingDate: subscription.nextBillingDate.split('T')[0],
            category: subscription.category,
            icon: subscription.icon || undefined,
            color: subscription.color || undefined,
            notes: subscription.notes || undefined,
          }}
          onSubmit={handleUpdate}
          onCancel={() => setIsEditOpen(false)}
          isLoading={updateSubscription.isPending}
          submitLabel="Save Changes"
        />
      </Sheet>

      {/* Delete confirmation */}
      {showDeleteConfirm && (
        <>
          <div
            className="fixed inset-0 z-50 bg-black/50"
            onClick={() => setShowDeleteConfirm(false)}
          />
          <div className="fixed inset-x-4 bottom-4 z-50 rounded-xl bg-surface p-6 shadow-xl dark:bg-dark-surface-elevated">
            <h3 className="mb-2 text-h3 text-text-primary dark:text-dark-text-primary">
              Delete {subscription.name}?
            </h3>
            <p className="mb-6 text-body text-text-secondary dark:text-dark-text-secondary">
              This action cannot be undone.
            </p>
            <div className="flex gap-3">
              <Button
                variant="secondary"
                className="flex-1"
                onClick={() => setShowDeleteConfirm(false)}
              >
                Cancel
              </Button>
              <Button
                variant="destructive"
                className="flex-1 bg-error text-white hover:bg-error/90 dark:bg-dark-error"
                onClick={handleDelete}
                disabled={deleteSubscription.isPending}
              >
                {deleteSubscription.isPending ? 'Deleting...' : 'Delete'}
              </Button>
            </div>
          </div>
        </>
      )}
    </div>
  );
}

function DetailRow({
  label,
  value,
  subValue,
  onEdit,
}: {
  label: string;
  value: string;
  subValue?: string;
  onEdit?: () => void;
}) {
  return (
    <div className="flex items-center justify-between rounded-lg bg-surface p-4 shadow-sm dark:bg-dark-surface dark:border dark:border-dark-border/50">
      <div>
        <p className="mb-1 text-label text-text-secondary dark:text-dark-text-secondary">
          {label}
        </p>
        <p className="text-body text-text-primary dark:text-dark-text-primary">
          {value}
        </p>
        {subValue && (
          <p className="text-body-sm text-text-tertiary dark:text-dark-text-tertiary">
            {subValue}
          </p>
        )}
      </div>
      {onEdit && (
        <button
          onClick={onEdit}
          className="rounded-full p-2 text-text-secondary transition-colors hover:bg-surface-hover hover:text-primary dark:text-dark-text-secondary dark:hover:bg-dark-surface-hover dark:hover:text-dark-primary"
        >
          <Pencil className="h-4 w-4" />
        </button>
      )}
    </div>
  );
}
