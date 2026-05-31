/**
 * Toast Notification Component
 *
 * Displays brief notifications that auto-dismiss.
 * - Appears top-right
 * - Auto-dismisses after 3s (configurable)
 * - Supports success/error/warning/info variants
 * - Slides in from right with 300ms transition
 * - Stacks vertically
 * - Uses textContent for XSS prevention
 * - Uses role="status" with aria-live="polite" for screen readers
 */

export type ToastVariant = 'success' | 'error' | 'warning' | 'info';

export interface ToastOptions {
  message: string;
  variant?: ToastVariant;
  duration?: number;
}

const TOAST_ICONS: Record<ToastVariant, string> = {
  success: '\u2713',
  error: '\u2717',
  warning: '\u26A0',
  info: '\u2139',
};

const TOAST_COLORS: Record<ToastVariant, string> = {
  success: '#3DB87A',
  error: '#D94B4B',
  warning: '#DFA23E',
  info: '#5B8AD4',
};

function getContainer(): HTMLElement {
  // Use the existing #toast-container from index.html if available
  const existing = document.getElementById('toast-container');
  if (existing) return existing;

  // Fallback: create container (e.g. in tests)
  const container = document.createElement('div');
  container.id = 'toast-container';
  container.className = 'fixed top-6 right-4 z-50 flex flex-col gap-2 pointer-events-none';
  container.setAttribute('aria-live', 'polite');
  container.setAttribute('role', 'status');
  document.body.appendChild(container);
  return container;
}

export const Toast = {
  show(options: ToastOptions): HTMLElement {
    const variant = options.variant ?? 'info';
    const duration = options.duration ?? 3000;
    const container = getContainer();

    // Create toast element
    const toast = document.createElement('div');
    toast.className = `toast-item toast-${variant}`;

    // Icon
    const icon = document.createElement('span');
    icon.style.cssText = `color: ${TOAST_COLORS[variant]}; font-size: 16px; font-weight: 700; width: 20px; text-align: center; flex-shrink: 0;`;
    icon.textContent = TOAST_ICONS[variant];

    // Message (XSS safe: textContent only)
    const message = document.createElement('span');
    message.style.cssText = 'flex: 1; font-size: 14px; line-height: 1.5;';
    message.className = 'text-content';
    message.textContent = options.message;

    // Close button
    const closeBtn = document.createElement('button');
    closeBtn.style.cssText = 'flex-shrink: 0; background: none; border: none; cursor: pointer; padding: 2px; font-size: 14px;';
    closeBtn.className = 'text-content-secondary hover:text-content transition-colors';
    closeBtn.setAttribute('aria-label', 'Dismiss notification');
    closeBtn.textContent = '\u2715';
    closeBtn.addEventListener('click', () => dismissToast(toast));

    toast.appendChild(icon);
    toast.appendChild(message);
    toast.appendChild(closeBtn);

    container.appendChild(toast);

    // Trigger slide-in animation (next frame)
    requestAnimationFrame(() => {
      toast.classList.add('visible');
    });

    // Auto-dismiss
    if (duration > 0) {
      setTimeout(() => {
        dismissToast(toast);
      }, duration);
    }

    return toast;
  },

  success(message: string, duration?: number): HTMLElement {
    return Toast.show({ message, variant: 'success', duration });
  },

  error(message: string, duration?: number): HTMLElement {
    return Toast.show({ message, variant: 'error', duration });
  },

  warning(message: string, duration?: number): HTMLElement {
    return Toast.show({ message, variant: 'warning', duration });
  },

  info(message: string, duration?: number): HTMLElement {
    return Toast.show({ message, variant: 'info', duration });
  },

  clearAll(): void {
    const container = document.getElementById('toast-container');
    if (container) {
      while (container.firstChild) {
        container.removeChild(container.firstChild);
      }
    }
  },
};

function dismissToast(toast: HTMLElement): void {
  if (!toast.parentNode) return;

  toast.classList.remove('visible');
  toast.classList.add('exiting');

  // Remove after exit transition (150ms)
  setTimeout(() => {
    toast.remove();
  }, 150);
}
