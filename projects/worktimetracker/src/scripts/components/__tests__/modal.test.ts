import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { Modal } from '../modal';

describe('Modal', () => {
  beforeEach(() => {
    vi.useFakeTimers();
    // Synchronous cleanup of any leftover modal state from previous tests
    Modal._reset();
    document.body.innerHTML = '';
  });

  afterEach(() => {
    vi.useRealTimers();
    document.body.innerHTML = '';
  });

  function createContent(text: string): HTMLElement {
    const div = document.createElement('div');
    div.textContent = text;
    return div;
  }

  it('should open a modal and add it to the DOM', () => {
    Modal.open({ title: 'Test Modal', content: createContent('Hello') });

    const overlay = document.querySelector('.modal-overlay');
    expect(overlay).toBeTruthy();
    expect(overlay?.querySelector('.modal-container')).toBeTruthy();
  });

  it('should display the title using textContent (XSS safe)', () => {
    Modal.open({ title: '<script>alert("xss")</script>', content: createContent('Content') });

    const title = document.querySelector('#modal-title');
    expect(title?.textContent).toBe('<script>alert("xss")</script>');
    // Ensure it was set via textContent, not innerHTML
    expect(title?.innerHTML).not.toContain('<script>');
  });

  it('should render the content slot', () => {
    const content = createContent('Form goes here');
    Modal.open({ title: 'Test', content });

    expect(document.body.textContent).toContain('Form goes here');
  });

  it('should show Cancel and Save buttons by default', () => {
    Modal.open({
      title: 'Test',
      content: createContent('Content'),
      onSubmit: () => {},
    });

    const buttons = document.querySelectorAll('.modal-container button');
    const buttonTexts = Array.from(buttons).map((b) => b.textContent?.trim());
    expect(buttonTexts).toContain('Cancel');
    expect(buttonTexts).toContain('Save');
  });

  it('should use custom button labels', () => {
    Modal.open({
      title: 'Test',
      content: createContent('Content'),
      submitLabel: 'Create',
      cancelLabel: 'Discard',
      onSubmit: () => {},
    });

    const buttons = document.querySelectorAll('.modal-container button');
    const buttonTexts = Array.from(buttons).map((b) => b.textContent?.trim());
    expect(buttonTexts).toContain('Create');
    expect(buttonTexts).toContain('Discard');
  });

  it('should hide action buttons when hideActions is true', () => {
    Modal.open({
      title: 'Test',
      content: createContent('Content'),
      hideActions: true,
    });

    const buttons = document.querySelectorAll('.modal-container .flex.justify-end button');
    expect(buttons.length).toBe(0);
  });

  it('should close on close button click', () => {
    Modal.open({ title: 'Test', content: createContent('Content') });
    expect(Modal.isOpen()).toBe(true);

    const closeBtn = document.querySelector('[aria-label="Close modal"]') as HTMLElement;
    closeBtn.click();

    vi.advanceTimersByTime(150);
    expect(Modal.isOpen()).toBe(false);
    expect(document.querySelector('.modal-overlay')).toBeNull();
  });

  it('should close on Escape key', () => {
    Modal.open({ title: 'Test', content: createContent('Content') });

    document.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' }));

    vi.advanceTimersByTime(150);
    expect(Modal.isOpen()).toBe(false);
  });

  it('should close on backdrop click', () => {
    const overlay = Modal.open({ title: 'Test', content: createContent('Content') });

    // Click on the overlay itself (not the modal container)
    overlay.dispatchEvent(new MouseEvent('click', { bubbles: true }));

    vi.advanceTimersByTime(150);
    expect(Modal.isOpen()).toBe(false);
  });

  it('should NOT close when clicking inside the modal container', () => {
    Modal.open({ title: 'Test', content: createContent('Content') });

    const modalContainer = document.querySelector('.modal-container') as HTMLElement;
    modalContainer.click();

    vi.advanceTimersByTime(150);
    expect(Modal.isOpen()).toBe(true);
  });

  it('should fire onClose callback when closed', () => {
    const onClose = vi.fn();
    Modal.open({ title: 'Test', content: createContent('Content'), onClose });

    Modal.close();

    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it('should fire onSubmit callback when Save is clicked', () => {
    const onSubmit = vi.fn();
    Modal.open({ title: 'Test', content: createContent('Content'), onSubmit });

    const saveBtn = Array.from(document.querySelectorAll('button')).find(
      (b) => b.textContent?.trim() === 'Save'
    );
    saveBtn?.click();

    expect(onSubmit).toHaveBeenCalledTimes(1);
  });

  it('should close previous modal when opening a new one', () => {
    Modal.open({ title: 'First', content: createContent('First') });
    Modal.open({ title: 'Second', content: createContent('Second') });

    vi.advanceTimersByTime(150);
    const overlays = document.querySelectorAll('.modal-overlay');
    // Only the second should remain (first removed after transition)
    expect(overlays.length).toBe(1);
    expect(document.querySelector('#modal-title')?.textContent).toBe('Second');
  });

  it('should have aria attributes for accessibility', () => {
    Modal.open({ title: 'Accessible Modal', content: createContent('Content') });

    const overlay = document.querySelector('.modal-overlay');
    expect(overlay?.getAttribute('role')).toBe('dialog');
    expect(overlay?.getAttribute('aria-modal')).toBe('true');
    expect(overlay?.getAttribute('aria-labelledby')).toBe('modal-title');
  });

  it('should report isOpen correctly', () => {
    expect(Modal.isOpen()).toBe(false);
    Modal.open({ title: 'Test', content: createContent('Content') });
    expect(Modal.isOpen()).toBe(true);
    Modal.close();
    vi.advanceTimersByTime(150);
    expect(Modal.isOpen()).toBe(false);
  });
});
