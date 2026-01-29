'use client';

import { Moon, Sun } from 'lucide-react';
import { useAppStore } from '@/stores/app-store';

export function Header() {
  const theme = useAppStore((state) => state.theme);
  const setTheme = useAppStore((state) => state.setTheme);

  const toggleTheme = () => {
    if (theme === 'light') {
      setTheme('dark');
      document.documentElement.classList.add('dark');
    } else if (theme === 'dark') {
      setTheme('system');
      if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
        document.documentElement.classList.add('dark');
      } else {
        document.documentElement.classList.remove('dark');
      }
    } else {
      setTheme('light');
      document.documentElement.classList.remove('dark');
    }
  };

  return (
    <header className="sticky top-0 z-40 border-b border-border bg-background/80 backdrop-blur-sm dark:border-dark-border dark:bg-dark-background/80">
      <div className="flex h-14 items-center justify-between px-4">
        <h1 className="text-h3 font-bold text-primary dark:text-dark-primary">
          SubTrack
        </h1>
        <button
          onClick={toggleTheme}
          className="rounded-full p-2 text-text-secondary transition-colors hover:bg-surface-hover hover:text-text-primary dark:text-dark-text-secondary dark:hover:bg-dark-surface-hover dark:hover:text-dark-text-primary"
          aria-label="Toggle theme"
        >
          {theme === 'dark' ? (
            <Sun className="h-5 w-5" />
          ) : (
            <Moon className="h-5 w-5" />
          )}
        </button>
      </div>
    </header>
  );
}
