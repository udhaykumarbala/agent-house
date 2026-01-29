'use client';

import { useQuery } from '@tanstack/react-query';
import { api } from '@/lib/api';

export function useSummary() {
  return useQuery({
    queryKey: ['summary'],
    queryFn: () => api.getSummary(),
  });
}

export function useUpcoming(days: number = 30) {
  return useQuery({
    queryKey: ['upcoming', days],
    queryFn: () => api.getUpcoming(days),
  });
}
