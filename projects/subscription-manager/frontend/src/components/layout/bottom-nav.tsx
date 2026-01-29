'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { Home, Calendar, Settings } from 'lucide-react';
import { cn } from '@/lib/utils';

const navItems = [
  { href: '/', label: 'Home', icon: Home },
  { href: '/upcoming', label: 'Upcoming', icon: Calendar },
  { href: '/settings', label: 'Settings', icon: Settings },
];

export function BottomNav() {
  const pathname = usePathname();

  return (
    <nav className="fixed bottom-0 left-0 right-0 z-40 border-t border-border bg-surface pb-safe dark:border-dark-border dark:bg-dark-surface">
      <div className="flex h-16 items-center justify-around">
        {navItems.map((item) => {
          const isActive = pathname === item.href;
          const Icon = item.icon;

          return (
            <Link
              key={item.href}
              href={item.href}
              className={cn(
                'flex flex-col items-center gap-1 px-6 py-2 transition-colors',
                isActive
                  ? 'text-primary dark:text-dark-primary'
                  : 'text-text-secondary hover:text-text-primary dark:text-dark-text-secondary dark:hover:text-dark-text-primary'
              )}
            >
              <Icon className="h-6 w-6" />
              <span className="text-caption font-medium">{item.label}</span>
              {isActive && (
                <div className="absolute -bottom-0 h-0.5 w-6 rounded-full bg-primary dark:bg-dark-primary" />
              )}
            </Link>
          );
        })}
      </div>
    </nav>
  );
}
