'use client';

import * as React from 'react';
import { cva, type VariantProps } from 'class-variance-authority';
import { cn } from '@/lib/utils';

const buttonVariants = cva(
  'inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-semibold transition-all duration-fast focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/30 focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 active:scale-[0.98]',
  {
    variants: {
      variant: {
        default:
          'bg-primary text-white hover:bg-primary-hover active:bg-primary-pressed shadow-sm dark:bg-dark-primary dark:hover:bg-dark-primary-hover',
        secondary:
          'bg-surface border border-border text-text-primary hover:bg-surface-hover dark:bg-dark-surface dark:border-dark-border dark:text-dark-text-primary dark:hover:bg-dark-surface-hover',
        ghost:
          'text-primary hover:bg-primary/10 dark:text-dark-primary dark:hover:bg-dark-primary/10',
        destructive:
          'text-error hover:bg-error/10 dark:text-dark-error dark:hover:bg-dark-error/10',
        link: 'text-primary underline-offset-4 hover:underline dark:text-dark-primary',
      },
      size: {
        default: 'h-11 px-6 py-3',
        sm: 'h-9 px-4 py-2 text-sm',
        lg: 'h-12 px-8 py-4',
        icon: 'h-11 w-11',
      },
    },
    defaultVariants: {
      variant: 'default',
      size: 'default',
    },
  }
);

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant, size, ...props }, ref) => {
    return (
      <button
        className={cn(buttonVariants({ variant, size, className }))}
        ref={ref}
        {...props}
      />
    );
  }
);
Button.displayName = 'Button';

export { Button, buttonVariants };
