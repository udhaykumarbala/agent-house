/**
 * Global Keyboard Shortcuts
 *
 * Active when no input/textarea/select is focused:
 * - Space: Toggle clock in/out on dashboard
 * - N: Open new time entry (navigate to timesheet)
 * - E: Navigate to employees
 * - P: Navigate to projects
 * - R: Navigate to reports
 * - Esc: Close any open modal
 * - ?: Show keyboard shortcuts help
 *
 * Shortcuts are disabled when typing in input/textarea/select fields.
 */

import { Modal } from '../components/modal';
import { navigate, getCurrentRoute } from '../router';

/** Tags that indicate the user is typing. */
const INPUT_TAGS = new Set(['INPUT', 'TEXTAREA', 'SELECT']);

function isInputFocused(): boolean {
  const active = document.activeElement;
  if (!active) return false;
  if (INPUT_TAGS.has(active.tagName)) return true;
  if ((active as HTMLElement).isContentEditable) return true;
  return false;
}

function handleKeydown(e: KeyboardEvent): void {
  // Don't interfere with modifier keys (Ctrl, Alt, Meta)
  if (e.ctrlKey || e.altKey || e.metaKey) return;

  // Esc always works (modal close is handled by modal.ts, but we ensure it here too)
  if (e.key === 'Escape') {
    if (Modal.isOpen()) {
      Modal.close();
      e.preventDefault();
    }
    return;
  }

  // All other shortcuts require no input focus and no modal open
  if (isInputFocused()) return;
  if (Modal.isOpen()) return;

  switch (e.key) {
    case ' ': {
      // Space: toggle clock in/out on dashboard
      e.preventDefault();
      const route = getCurrentRoute();
      if (route === '/dashboard') {
        // Click the Clock In button if available
        const clockInBtn = document.querySelector<HTMLButtonElement>('.btn-accent.btn-lg');
        if (clockInBtn && !clockInBtn.disabled) {
          clockInBtn.click();
          return;
        }
        // Or try to clock out the first active timer via the stop button
        const stopBtn = document.querySelector<HTMLButtonElement>('#timer-stop-btn');
        if (stopBtn) {
          stopBtn.click();
        }
      }
      break;
    }
    case 'n':
    case 'N':
      e.preventDefault();
      // Navigate to timesheet which has the new entry modal
      if (getCurrentRoute() !== '/timesheet') {
        navigate('/timesheet');
      }
      // Trigger the add entry button after a short delay for view to render
      setTimeout(() => {
        const addBtn = document.querySelector<HTMLButtonElement>('[data-action="add-entry"]');
        if (addBtn) addBtn.click();
      }, 300);
      break;

    case 'e':
    case 'E':
      e.preventDefault();
      navigate('/employees');
      break;

    case 'p':
    case 'P':
      e.preventDefault();
      navigate('/projects');
      break;

    case 'r':
    case 'R':
      e.preventDefault();
      navigate('/reports');
      break;

    case 'd':
    case 'D':
      e.preventDefault();
      navigate('/dashboard');
      break;

    case '?':
      e.preventDefault();
      showShortcutsHelp();
      break;
  }
}

function showShortcutsHelp(): void {
  if (Modal.isOpen()) return;

  const content = document.createElement('div');
  content.style.cssText = 'display: flex; flex-direction: column; gap: 10px;';

  const shortcuts = [
    { key: 'Space', desc: 'Clock in/out (on Dashboard)' },
    { key: 'N', desc: 'New time entry' },
    { key: 'D', desc: 'Go to Dashboard' },
    { key: 'E', desc: 'Go to Employees' },
    { key: 'P', desc: 'Go to Projects' },
    { key: 'R', desc: 'Go to Reports' },
    { key: 'Esc', desc: 'Close modal' },
    { key: '?', desc: 'Show this help' },
  ];

  for (const s of shortcuts) {
    const row = document.createElement('div');
    row.style.cssText = 'display: flex; align-items: center; gap: 12px;';

    const keyBadge = document.createElement('kbd');
    keyBadge.style.cssText = 'display: inline-flex; align-items: center; justify-content: center; min-width: 32px; padding: 4px 8px; background: rgba(91, 70, 178, 0.15); border: 1px solid rgba(91, 70, 178, 0.30); border-radius: 6px; font-size: 13px; font-weight: 600; color: #CF8A2E; font-family: Inter, sans-serif;';
    keyBadge.textContent = s.key;

    const desc = document.createElement('span');
    desc.style.cssText = 'font-size: 14px; color: #ECE9F5;';
    desc.textContent = s.desc;

    row.appendChild(keyBadge);
    row.appendChild(desc);
    content.appendChild(row);
  }

  Modal.open({
    title: 'Keyboard Shortcuts',
    content,
    hideActions: true,
  });
}

/** Initialize global keyboard shortcut listener. */
export function initKeyboardShortcuts(): void {
  document.addEventListener('keydown', handleKeydown);
}

/** Remove global keyboard shortcut listener. */
export function destroyKeyboardShortcuts(): void {
  document.removeEventListener('keydown', handleKeydown);
}
