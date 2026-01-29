'use client';

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '@/lib/api';
import type { SubscriptionInput } from '@/types/subscription';
import { useAppStore } from '@/stores/app-store';

export function useSubscriptions() {
  return useQuery({
    queryKey: ['subscriptions'],
    queryFn: () => api.getSubscriptions(),
  });
}

export function useSubscription(id: string) {
  return useQuery({
    queryKey: ['subscription', id],
    queryFn: () => api.getSubscription(id),
    enabled: !!id,
  });
}

export function useCreateSubscription() {
  const queryClient = useQueryClient();
  const addToast = useAppStore((state) => state.addToast);

  return useMutation({
    mutationFn: (input: SubscriptionInput) => api.createSubscription(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['subscriptions'] });
      queryClient.invalidateQueries({ queryKey: ['summary'] });
      queryClient.invalidateQueries({ queryKey: ['upcoming'] });
      addToast({ message: 'Subscription added', type: 'success' });
    },
    onError: (error: Error) => {
      addToast({ message: error.message || 'Failed to add subscription', type: 'error' });
    },
  });
}

export function useUpdateSubscription() {
  const queryClient = useQueryClient();
  const addToast = useAppStore((state) => state.addToast);

  return useMutation({
    mutationFn: ({ id, input }: { id: string; input: SubscriptionInput }) =>
      api.updateSubscription(id, input),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['subscriptions'] });
      queryClient.invalidateQueries({ queryKey: ['subscription', variables.id] });
      queryClient.invalidateQueries({ queryKey: ['summary'] });
      queryClient.invalidateQueries({ queryKey: ['upcoming'] });
      addToast({ message: 'Subscription updated', type: 'success' });
    },
    onError: (error: Error) => {
      addToast({ message: error.message || 'Failed to update subscription', type: 'error' });
    },
  });
}

export function useDeleteSubscription() {
  const queryClient = useQueryClient();
  const addToast = useAppStore((state) => state.addToast);

  return useMutation({
    mutationFn: (id: string) => api.deleteSubscription(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['subscriptions'] });
      queryClient.invalidateQueries({ queryKey: ['summary'] });
      queryClient.invalidateQueries({ queryKey: ['upcoming'] });
      addToast({
        message: 'Subscription deleted',
        type: 'success',
      });
    },
    onError: (error: Error) => {
      addToast({ message: error.message || 'Failed to delete subscription', type: 'error' });
    },
  });
}
