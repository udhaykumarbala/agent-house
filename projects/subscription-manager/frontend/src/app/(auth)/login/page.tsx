'use client';

import * as React from 'react';
import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { useLogin } from '@/hooks/use-auth';

export default function LoginPage() {
  const [email, setEmail] = React.useState('');
  const [password, setPassword] = React.useState('');
  const login = useLogin();

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    login.mutate({ email, password });
  };

  return (
    <>
      <div className="mb-8 text-center">
        <h1 className="mb-2 text-h1 font-bold text-primary dark:text-dark-primary">
          SubTrack
        </h1>
        <p className="text-body text-text-secondary dark:text-dark-text-secondary">
          Welcome back
        </p>
      </div>

      <form onSubmit={handleSubmit} className="space-y-4">
        <Input
          label="Email"
          type="email"
          placeholder="you@example.com"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
        />
        <Input
          label="Password"
          type="password"
          placeholder="Enter your password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />

        <Button type="submit" className="w-full" disabled={login.isPending}>
          {login.isPending ? 'Signing in...' : 'Sign In'}
        </Button>
      </form>

      <p className="mt-6 text-center text-body-sm text-text-secondary dark:text-dark-text-secondary">
        Don&apos;t have an account?{' '}
        <Link
          href="/register"
          className="text-primary hover:underline dark:text-dark-primary"
        >
          Sign up
        </Link>
      </p>
    </>
  );
}
