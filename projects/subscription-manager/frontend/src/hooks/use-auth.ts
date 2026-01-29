'use client';

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useRouter } from 'next/navigation';
import { api } from '@/lib/api';
import { useAppStore } from '@/stores/app-store';
import type { LoginInput, RegisterInput } from '@/types/user';

export function useUser() {
  const setUser = useAppStore((state) => state.setUser);

  return useQuery({
    queryKey: ['user'],
    queryFn: async () => {
      const user = await api.me();
      setUser(user);
      return user;
    },
    retry: false,
  });
}

export function useLogin() {
  const queryClient = useQueryClient();
  const setUser = useAppStore((state) => state.setUser);
  const addToast = useAppStore((state) => state.addToast);
  const router = useRouter();

  return useMutation({
    mutationFn: (input: LoginInput) => api.login(input),
    onSuccess: (user) => {
      setUser(user);
      queryClient.setQueryData(['user'], user);
      queryClient.invalidateQueries({ queryKey: ['subscriptions'] });
      queryClient.invalidateQueries({ queryKey: ['summary'] });
      router.push('/');
    },
    onError: (error: Error) => {
      addToast({ message: error.message || 'Login failed', type: 'error' });
    },
  });
}

export function useRegister() {
  const queryClient = useQueryClient();
  const setUser = useAppStore((state) => state.setUser);
  const addToast = useAppStore((state) => state.addToast);
  const router = useRouter();

  return useMutation({
    mutationFn: (input: RegisterInput) => api.register(input),
    onSuccess: (user) => {
      setUser(user);
      queryClient.setQueryData(['user'], user);
      addToast({ message: 'Welcome to SubTrack!', type: 'success' });
      router.push('/');
    },
    onError: (error: Error) => {
      addToast({ message: error.message || 'Registration failed', type: 'error' });
    },
  });
}

export function useLogout() {
  const queryClient = useQueryClient();
  const setUser = useAppStore((state) => state.setUser);
  const router = useRouter();

  return useMutation({
    mutationFn: () => api.logout(),
    onSuccess: () => {
      setUser(null);
      queryClient.clear();
      router.push('/login');
    },
  });
}
