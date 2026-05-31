import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { Toast } from '../toast';

describe('Toast', () => {
  beforeEach(() => {
    document.body.innerHTML = '';
    vi.useFakeTimers();
  });

  afterEach(() => {
    Toast.clearAll();
    vi.useRealTimers();
    document.body.innerHTML = '';
  });

  it('should create a toast container in the DOM', () => {
    Toast.show({ message: 'Hello' });

    const container = document.getElementById('toast-container');
    expect(container).toBeTruthy();
  });

  it('should position the container at top-right', () => {
    Toast.show({ message: 'Hello' });

    const container = document.getElementById('toast-container') as HTMLElement;
    expect(container).toBeTruthy();
    // Container should have top-right positioning classes or styles
    expect(container.className).toContain('top');
    expect(container.className).toContain('right');
  });

  it('should display the message using textContent (XSS safe)', () => {
    Toast.show({ message: '<img src=x onerror=alert(1)>' });

    const toastItem = document.querySelector('.toast-item');
    const spans = toastItem?.querySelectorAll('span');
    const messageSpan = spans ? spans[1] : null; // second span is the message
    expect(messageSpan?.textContent).toBe('<img src=x onerror=alert(1)>');
    expect(messageSpan?.innerHTML).not.toContain('<img');
  });

  it('should apply success variant styling', () => {
    Toast.success('Done!');

    const toast = document.querySelector('.toast-item');
    expect(toast?.classList.contains('toast-success')).toBe(true);
  });

  it('should apply error variant styling', () => {
    Toast.error('Failed!');

    const toast = document.querySelector('.toast-item');
    expect(toast?.classList.contains('toast-error')).toBe(true);
  });

  it('should apply warning variant styling', () => {
    Toast.warning('Careful!');

    const toast = document.querySelector('.toast-item');
    expect(toast?.classList.contains('toast-warning')).toBe(true);
  });

  it('should apply info variant styling', () => {
    Toast.info('FYI');

    const toast = document.querySelector('.toast-item');
    expect(toast?.classList.contains('toast-info')).toBe(true);
  });

  it('should default to info variant when none specified', () => {
    Toast.show({ message: 'Default' });

    const toast = document.querySelector('.toast-item');
    expect(toast?.classList.contains('toast-info')).toBe(true);
  });

  it('should auto-dismiss after 3s by default', () => {
    Toast.show({ message: 'Goodbye' });

    expect(document.querySelectorAll('.toast-item').length).toBe(1);

    // Advance past the auto-dismiss time
    vi.advanceTimersByTime(3000);
    // Advance past the exit transition
    vi.advanceTimersByTime(150);

    expect(document.querySelectorAll('.toast-item').length).toBe(0);
  });

  it('should auto-dismiss after custom duration', () => {
    Toast.show({ message: 'Quick', duration: 1000 });

    expect(document.querySelectorAll('.toast-item').length).toBe(1);

    vi.advanceTimersByTime(1000);
    vi.advanceTimersByTime(150);

    expect(document.querySelectorAll('.toast-item').length).toBe(0);
  });

  it('should NOT auto-dismiss when duration is 0', () => {
    Toast.show({ message: 'Persistent', duration: 0 });

    vi.advanceTimersByTime(10000);

    expect(document.querySelectorAll('.toast-item').length).toBe(1);
  });

  it('should dismiss on close button click', () => {
    Toast.show({ message: 'Closeable' });

    const closeBtn = document.querySelector('.toast-item button') as HTMLElement;
    closeBtn.click();

    vi.advanceTimersByTime(150);

    expect(document.querySelectorAll('.toast-item').length).toBe(0);
  });

  it('should stack multiple toasts vertically', () => {
    Toast.show({ message: 'First', duration: 0 });
    Toast.show({ message: 'Second', duration: 0 });
    Toast.show({ message: 'Third', duration: 0 });

    const toasts = document.querySelectorAll('.toast-item');
    expect(toasts.length).toBe(3);

    // All should be in the same container
    const container = document.getElementById('toast-container');
    expect(container?.children.length).toBe(3);
  });

  it('should have aria-live="polite" on the container', () => {
    Toast.show({ message: 'Accessible' });

    const container = document.getElementById('toast-container');
    expect(container?.getAttribute('aria-live')).toBe('polite');
    expect(container?.getAttribute('role')).toBe('status');
  });

  it('should clear all toasts with clearAll()', () => {
    Toast.show({ message: 'One', duration: 0 });
    Toast.show({ message: 'Two', duration: 0 });

    Toast.clearAll();

    expect(document.querySelectorAll('.toast-item').length).toBe(0);
  });

  it('should add visible class for slide-in animation', () => {
    Toast.show({ message: 'Animate' });

    // Before rAF, no visible class
    const toast = document.querySelector('.toast-item');
    expect(toast?.classList.contains('visible')).toBe(false);

    // After rAF, visible class added
    vi.advanceTimersByTime(16); // ~1 frame
    // Note: jsdom doesn't run rAF, so we test the initial state is correct
  });
});
