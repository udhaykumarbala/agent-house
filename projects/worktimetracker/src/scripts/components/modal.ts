/**
 * Modal Component
 *
 * Reusable modal dialog supporting form content slots.
 * - Opens/closes with 250ms/150ms transitions (UI spec)
 * - Closes on Esc key and backdrop click
 * - Focus trap for accessibility
 * - Uses textContent for user data (XSS prevention)
 */

export interface ModalOptions {
  title: string;
  content: HTMLElement;
  onClose?: () => void;
  onSubmit?: (e: Event) => void;
  submitLabel?: string;
  cancelLabel?: string;
  hideActions?: boolean;
}

let activeModal: HTMLElement | null = null;
let activeOnClose: (() => void) | null = null;
let previousFocus: HTMLElement | null = null;

function handleKeydown(e: KeyboardEvent): void {
  if (e.key === 'Escape') {
    Modal.close();
    return;
  }

  if (e.key === 'Tab' && activeModal) {
    trapFocus(e, activeModal);
  }
}

function handleBackdropClick(e: MouseEvent): void {
  const target = e.target as HTMLElement;
  if (target.classList.contains('modal-overlay')) {
    Modal.close();
  }
}

function trapFocus(e: KeyboardEvent, container: HTMLElement): void {
  const focusable = container.querySelectorAll<HTMLElement>(
    'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
  );

  if (focusable.length === 0) return;

  const first = focusable[0];
  const last = focusable[focusable.length - 1];

  if (e.shiftKey && document.activeElement === first) {
    e.preventDefault();
    last.focus();
  } else if (!e.shiftKey && document.activeElement === last) {
    e.preventDefault();
    first.focus();
  }
}

export const Modal = {
  open(options: ModalOptions): HTMLElement {
    // Close any existing modal first
    if (activeModal) {
      Modal.close();
    }

    previousFocus = document.activeElement as HTMLElement;
    activeOnClose = options.onClose ?? null;

    // Create overlay
    const overlay = document.createElement('div');
    overlay.className = 'modal-overlay';
    overlay.setAttribute('role', 'dialog');
    overlay.setAttribute('aria-modal', 'true');
    overlay.setAttribute('aria-labelledby', 'modal-title');

    // Create modal container
    const modal = document.createElement('div');
    modal.className = 'modal-container';

    // Header row: title + close button
    const header = document.createElement('div');
    header.style.cssText = 'display: flex; align-items: center; justify-content: space-between; margin-bottom: 24px;';

    const title = document.createElement('h2');
    title.id = 'modal-title';
    title.className = 'modal-title';
    title.textContent = options.title;

    const closeBtn = document.createElement('button');
    closeBtn.type = 'button';
    closeBtn.className = 'btn-ghost btn-sm';
    closeBtn.setAttribute('aria-label', 'Close modal');
    closeBtn.textContent = '\u2715';
    closeBtn.addEventListener('click', () => Modal.close());

    header.appendChild(title);
    header.appendChild(closeBtn);

    // Content slot
    const contentWrapper = document.createElement('div');
    contentWrapper.style.marginBottom = '24px';
    contentWrapper.appendChild(options.content);

    // Assemble modal
    modal.appendChild(header);
    modal.appendChild(contentWrapper);

    if (!options.hideActions) {
      const actions = document.createElement('div');
      actions.style.cssText = 'display: flex; justify-content: flex-end; gap: 8px;';

      const cancelBtn = document.createElement('button');
      cancelBtn.type = 'button';
      cancelBtn.className = 'btn-secondary';
      cancelBtn.textContent = options.cancelLabel ?? 'Cancel';
      cancelBtn.addEventListener('click', () => Modal.close());

      actions.appendChild(cancelBtn);

      if (options.onSubmit) {
        const submitBtn = document.createElement('button');
        submitBtn.type = 'button';
        submitBtn.className = 'btn-primary';
        submitBtn.textContent = options.submitLabel ?? 'Save';
        submitBtn.addEventListener('click', (e) => {
          options.onSubmit!(e);
        });
        actions.appendChild(submitBtn);
      }

      modal.appendChild(actions);
    }

    overlay.appendChild(modal);
    document.body.appendChild(overlay);
    activeModal = overlay;

    // Event listeners
    document.addEventListener('keydown', handleKeydown);
    overlay.addEventListener('click', handleBackdropClick);

    // Trigger open animation (next frame for CSS transition)
    requestAnimationFrame(() => {
      overlay.classList.add('active');
    });

    // Focus first focusable element in the modal
    requestAnimationFrame(() => {
      const firstFocusable = modal.querySelector<HTMLElement>(
        'input, select, textarea, button'
      );
      if (firstFocusable) {
        firstFocusable.focus();
      }
    });

    return overlay;
  },

  close(): void {
    if (!activeModal) return;

    const overlay = activeModal;
    const onClose = activeOnClose;

    // Clear state synchronously so isOpen() returns false immediately
    activeModal = null;
    activeOnClose = null;

    overlay.classList.remove('active');
    overlay.classList.add('closing');

    document.removeEventListener('keydown', handleKeydown);
    overlay.removeEventListener('click', handleBackdropClick);

    // Wait for exit transition (150ms per UI spec), then remove from DOM
    setTimeout(() => {
      overlay.remove();

      // Restore focus to previously focused element
      if (previousFocus) {
        previousFocus.focus();
        previousFocus = null;
      }
    }, 150);

    // Fire onClose callback
    if (onClose) {
      onClose();
    }
  },

  isOpen(): boolean {
    return activeModal !== null;
  },

  /** Synchronous cleanup — for testing only */
  _reset(): void {
    if (activeModal) {
      activeModal.remove();
    }
    activeModal = null;
    activeOnClose = null;
    previousFocus = null;
    document.removeEventListener('keydown', handleKeydown);
  },
};
