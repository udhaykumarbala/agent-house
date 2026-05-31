/**
 * Keyboard Shortcuts Tests
 */
import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import { initKeyboardShortcuts, destroyKeyboardShortcuts } from '../keyboard';
import { Modal } from '../../components/modal';

describe('Keyboard Shortcuts', () => {
  beforeEach(() => {
    document.body.innerHTML = `
      <div id="toast-container"></div>
      <div id="modal-container"></div>
      <div id="view-container"></div>
    `;
    Modal._reset();
    initKeyboardShortcuts();
  });

  afterEach(() => {
    destroyKeyboardShortcuts();
    Modal._reset();
    document.body.innerHTML = '';
    window.location.hash = '';
  });

  function pressKey(key: string, options: Partial<KeyboardEventInit> = {}): void {
    const event = new KeyboardEvent('keydown', { key, bubbles: true, ...options });
    document.dispatchEvent(event);
  }

  it('should not fire shortcuts when input is focused', () => {
    const input = document.createElement('input');
    document.body.appendChild(input);
    input.focus();

    // "e" should navigate to employees normally, but not when input is focused
    const hashBefore = window.location.hash;
    pressKey('e');
    expect(window.location.hash).toBe(hashBefore);
  });

  it('should not fire shortcuts when textarea is focused', () => {
    const textarea = document.createElement('textarea');
    document.body.appendChild(textarea);
    textarea.focus();

    const hashBefore = window.location.hash;
    pressKey('r');
    expect(window.location.hash).toBe(hashBefore);
  });

  it('should navigate to employees on "e" key', () => {
    pressKey('e');
    expect(window.location.hash).toBe('#/employees');
  });

  it('should navigate to projects on "p" key', () => {
    pressKey('p');
    expect(window.location.hash).toBe('#/projects');
  });

  it('should navigate to reports on "r" key', () => {
    pressKey('r');
    expect(window.location.hash).toBe('#/reports');
  });

  it('should navigate to dashboard on "d" key', () => {
    pressKey('d');
    expect(window.location.hash).toBe('#/dashboard');
  });

  it('should close modal on Escape key', () => {
    const content = document.createElement('div');
    content.textContent = 'Test';
    Modal.open({ title: 'Test', content, hideActions: true });
    expect(Modal.isOpen()).toBe(true);

    pressKey('Escape');
    expect(Modal.isOpen()).toBe(false);
  });

  it('should not fire shortcuts when modifier keys are held', () => {
    const hashBefore = window.location.hash;
    pressKey('e', { ctrlKey: true });
    expect(window.location.hash).toBe(hashBefore);
  });

  it('should not fire navigation shortcuts when modal is open', () => {
    const content = document.createElement('div');
    Modal.open({ title: 'Test', content, hideActions: true });

    const hashBefore = window.location.hash;
    pressKey('e');
    // Should not navigate since modal is open
    expect(window.location.hash).toBe(hashBefore);
  });

  it('should show shortcuts help on "?" key', () => {
    pressKey('?');
    expect(Modal.isOpen()).toBe(true);
    const title = document.querySelector('#modal-title');
    expect(title?.textContent).toBe('Keyboard Shortcuts');
  });

  it('handles uppercase key presses', () => {
    pressKey('E');
    expect(window.location.hash).toBe('#/employees');
  });

  it('destroyKeyboardShortcuts removes listener', () => {
    destroyKeyboardShortcuts();

    const hashBefore = window.location.hash;
    pressKey('e');
    expect(window.location.hash).toBe(hashBefore);
  });
});
