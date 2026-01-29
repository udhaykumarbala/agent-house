'use client';

import { Header } from '@/components/layout/header';
import { BottomNav } from '@/components/layout/bottom-nav';
import { FAB } from '@/components/layout/fab';
import { Sheet } from '@/components/ui/sheet';
import { SubscriptionForm } from '@/components/subscription/subscription-form';
import { useAppStore } from '@/stores/app-store';
import { useCreateSubscription } from '@/hooks/use-subscriptions';
import type { SubscriptionInput } from '@/types/subscription';

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const isAddSheetOpen = useAppStore((state) => state.isAddSheetOpen);
  const setAddSheetOpen = useAppStore((state) => state.setAddSheetOpen);
  const createSubscription = useCreateSubscription();

  const handleSubmit = (data: SubscriptionInput) => {
    createSubscription.mutate(data, {
      onSuccess: () => {
        setAddSheetOpen(false);
      },
    });
  };

  return (
    <div className="min-h-screen bg-background dark:bg-dark-background">
      <Header />
      <main className="pb-24 pt-4">{children}</main>
      <FAB />
      <BottomNav />

      <Sheet
        open={isAddSheetOpen}
        onClose={() => setAddSheetOpen(false)}
        title="Add Subscription"
      >
        <SubscriptionForm
          onSubmit={handleSubmit}
          onCancel={() => setAddSheetOpen(false)}
          isLoading={createSubscription.isPending}
        />
      </Sheet>
    </div>
  );
}
