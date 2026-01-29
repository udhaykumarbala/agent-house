'use client';

import { Moon, Sun, Monitor, LogOut, Trash2, Download } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { useAppStore } from '@/stores/app-store';
import { useLogout } from '@/hooks/use-auth';
import { cn } from '@/lib/utils';

export default function SettingsPage() {
  const theme = useAppStore((state) => state.theme);
  const setTheme = useAppStore((state) => state.setTheme);
  const user = useAppStore((state) => state.user);
  const logout = useLogout();

  const handleThemeChange = (newTheme: 'light' | 'dark' | 'system') => {
    setTheme(newTheme);
    const root = document.documentElement;

    if (newTheme === 'dark') {
      root.classList.add('dark');
    } else if (newTheme === 'light') {
      root.classList.remove('dark');
    } else {
      if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
        root.classList.add('dark');
      } else {
        root.classList.remove('dark');
      }
    }
  };

  const themeOptions = [
    { value: 'light', label: 'Light', icon: Sun },
    { value: 'dark', label: 'Dark', icon: Moon },
    { value: 'system', label: 'System', icon: Monitor },
  ] as const;

  return (
    <div className="px-4">
      <h1 className="mb-6 text-h2 text-text-primary dark:text-dark-text-primary">
        Settings
      </h1>

      {/* Appearance */}
      <section className="mb-8">
        <h2 className="mb-4 text-label uppercase tracking-wider text-text-secondary dark:text-dark-text-secondary">
          Appearance
        </h2>

        <div className="rounded-lg bg-surface p-4 shadow-sm dark:bg-dark-surface dark:border dark:border-dark-border/50">
          <p className="mb-3 text-body-sm text-text-primary dark:text-dark-text-primary">
            Theme
          </p>
          <div className="flex gap-2">
            {themeOptions.map((option) => {
              const Icon = option.icon;
              const isSelected = theme === option.value;
              return (
                <button
                  key={option.value}
                  onClick={() => handleThemeChange(option.value)}
                  className={cn(
                    'flex flex-1 flex-col items-center gap-2 rounded-lg border p-3 transition-colors',
                    isSelected
                      ? 'border-primary bg-primary/5 text-primary dark:border-dark-primary dark:bg-dark-primary/10 dark:text-dark-primary'
                      : 'border-border text-text-secondary hover:bg-surface-hover dark:border-dark-border dark:text-dark-text-secondary dark:hover:bg-dark-surface-hover'
                  )}
                >
                  <Icon className="h-5 w-5" />
                  <span className="text-sm font-medium">{option.label}</span>
                </button>
              );
            })}
          </div>
        </div>
      </section>

      {/* Data */}
      <section className="mb-8">
        <h2 className="mb-4 text-label uppercase tracking-wider text-text-secondary dark:text-dark-text-secondary">
          Data
        </h2>

        <div className="space-y-2">
          <button className="flex w-full items-center justify-between rounded-lg bg-surface p-4 text-left shadow-sm transition-colors hover:bg-surface-hover dark:bg-dark-surface dark:border dark:border-dark-border/50 dark:hover:bg-dark-surface-hover">
            <div className="flex items-center gap-3">
              <Download className="h-5 w-5 text-text-secondary dark:text-dark-text-secondary" />
              <span className="text-body text-text-primary dark:text-dark-text-primary">
                Export Subscriptions
              </span>
            </div>
            <span className="text-body-sm text-text-tertiary dark:text-dark-text-tertiary">
              Coming soon
            </span>
          </button>

          <button className="flex w-full items-center gap-3 rounded-lg bg-surface p-4 text-left shadow-sm transition-colors hover:bg-error-bg dark:bg-dark-surface dark:border dark:border-dark-border/50 dark:hover:bg-dark-error/10">
            <Trash2 className="h-5 w-5 text-error dark:text-dark-error" />
            <span className="text-body text-error dark:text-dark-error">
              Clear All Data
            </span>
          </button>
        </div>
      </section>

      {/* Account */}
      {user && (
        <section className="mb-8">
          <h2 className="mb-4 text-label uppercase tracking-wider text-text-secondary dark:text-dark-text-secondary">
            Account
          </h2>

          <div className="rounded-lg bg-surface p-4 shadow-sm dark:bg-dark-surface dark:border dark:border-dark-border/50">
            <p className="mb-1 text-body text-text-primary dark:text-dark-text-primary">
              {user.email}
            </p>
            <p className="mb-4 text-body-sm text-text-secondary dark:text-dark-text-secondary">
              Member since {new Date(user.createdAt).toLocaleDateString()}
            </p>
            <Button
              variant="secondary"
              onClick={() => logout.mutate()}
              disabled={logout.isPending}
            >
              <LogOut className="mr-2 h-4 w-4" />
              {logout.isPending ? 'Logging out...' : 'Log Out'}
            </Button>
          </div>
        </section>
      )}

      {/* About */}
      <section>
        <h2 className="mb-4 text-label uppercase tracking-wider text-text-secondary dark:text-dark-text-secondary">
          About
        </h2>

        <div className="rounded-lg bg-surface p-4 shadow-sm dark:bg-dark-surface dark:border dark:border-dark-border/50">
          <div className="mb-3 flex items-center justify-between">
            <span className="text-body text-text-primary dark:text-dark-text-primary">
              Version
            </span>
            <span className="text-body-sm text-text-secondary dark:text-dark-text-secondary">
              1.0.0
            </span>
          </div>

          <div className="border-t border-border pt-3 dark:border-dark-border">
            <a
              href="#"
              className="mb-2 block text-body text-text-primary hover:text-primary dark:text-dark-text-primary dark:hover:text-dark-primary"
            >
              Privacy Policy
            </a>
            <a
              href="#"
              className="block text-body text-text-primary hover:text-primary dark:text-dark-text-primary dark:hover:text-dark-primary"
            >
              Terms of Service
            </a>
          </div>
        </div>

        <p className="mt-6 text-center text-body-sm text-text-tertiary dark:text-dark-text-tertiary">
          Made with care by the SubTrack team
        </p>
      </section>
    </div>
  );
}
