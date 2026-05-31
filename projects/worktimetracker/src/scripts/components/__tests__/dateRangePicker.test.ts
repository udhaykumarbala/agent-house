import { describe, it, expect, beforeEach, vi } from 'vitest';
import { createDateRangePicker, computePresetRange } from '../dateRangePicker';
import type { DateRange, PresetKey } from '../dateRangePicker';

describe('DateRangePicker', () => {
  let container: HTMLElement;

  beforeEach(() => {
    document.body.innerHTML = '';
    container = document.createElement('div');
    document.body.appendChild(container);
  });

  function mountPicker(
    onChangeFn?: (range: DateRange) => void,
    initialPreset: PresetKey = 'thisWeek',
    initialStart?: string,
    initialEnd?: string,
  ): { picker: HTMLElement; onChange: ReturnType<typeof vi.fn> } {
    const fn = onChangeFn ?? vi.fn();
    const picker = createDateRangePicker({
      onChange: fn,
      initialPreset,
      initialStart,
      initialEnd,
    });
    container.appendChild(picker);
    return { picker, onChange: fn as ReturnType<typeof vi.fn> };
  }

  // --- Preset buttons ---

  it('should render all 6 preset buttons', () => {
    const { picker } = mountPicker();
    const buttons = picker.querySelectorAll('.date-preset');
    expect(buttons.length).toBe(6);

    const labels = Array.from(buttons).map((b) => b.textContent);
    expect(labels).toEqual([
      'Today',
      'This Week',
      'This Month',
      'Last Week',
      'Last Month',
      'Custom',
    ]);
  });

  it('should highlight the initial preset as selected', () => {
    const { picker } = mountPicker(vi.fn(), 'thisMonth');
    const selected = picker.querySelector('.date-preset.selected');
    expect(selected?.textContent).toBe('This Month');
    expect(selected?.getAttribute('aria-pressed')).toBe('true');
  });

  it('should default to This Week preset', () => {
    const { picker } = mountPicker();
    const selected = picker.querySelector('.date-preset.selected');
    expect(selected?.textContent).toBe('This Week');
  });

  // --- Preset click emits range ---

  it('should emit date range when a preset is clicked', () => {
    const onChange = vi.fn();
    const { picker } = mountPicker(onChange);

    const todayBtn = picker.querySelector('[data-preset="today"]') as HTMLButtonElement;
    todayBtn.click();

    expect(onChange).toHaveBeenCalledTimes(1);
    const range = onChange.mock.calls[0][0] as DateRange;
    expect(range.preset).toBe('today');
    expect(range.start).toBe(range.end); // Today is a single day
  });

  it('should switch selected class when clicking a different preset', () => {
    const { picker } = mountPicker(vi.fn(), 'thisWeek');

    const thisMonthBtn = picker.querySelector('[data-preset="thisMonth"]') as HTMLButtonElement;
    thisMonthBtn.click();

    const selectedButtons = picker.querySelectorAll('.date-preset.selected');
    expect(selectedButtons.length).toBe(1);
    expect(selectedButtons[0].textContent).toBe('This Month');

    // Previous should no longer be selected
    const thisWeekBtn = picker.querySelector('[data-preset="thisWeek"]') as HTMLButtonElement;
    expect(thisWeekBtn.classList.contains('selected')).toBe(false);
    expect(thisWeekBtn.getAttribute('aria-pressed')).toBe('false');
  });

  // --- Custom range ---

  it('should hide custom date inputs by default (non-custom preset)', () => {
    const { picker } = mountPicker(vi.fn(), 'today');
    const customContainer = picker.querySelectorAll('input[type="date"]');
    // Inputs exist in the DOM but container is hidden
    const parent = customContainer[0]?.parentElement;
    expect(parent?.style.display).toBe('none');
  });

  it('should show custom date inputs when Custom is clicked', () => {
    const { picker } = mountPicker();

    const customBtn = picker.querySelector('[data-preset="custom"]') as HTMLButtonElement;
    customBtn.click();

    const inputs = picker.querySelectorAll('input[type="date"]');
    expect(inputs.length).toBe(2);
    const parent = inputs[0]?.parentElement;
    expect(parent?.style.display).toBe('flex');
  });

  it('should emit range when custom dates are changed', () => {
    const onChange = vi.fn();
    const { picker } = mountPicker(onChange, 'custom', '2026-03-01', '2026-03-15');

    const startInput = picker.querySelector('#drp-start') as HTMLInputElement;
    const endInput = picker.querySelector('#drp-end') as HTMLInputElement;

    startInput.value = '2026-03-05';
    endInput.value = '2026-03-20';
    endInput.dispatchEvent(new Event('change'));

    expect(onChange).toHaveBeenCalledTimes(1);
    const range = onChange.mock.calls[0][0] as DateRange;
    expect(range.start).toBe('2026-03-05');
    expect(range.end).toBe('2026-03-20');
    expect(range.preset).toBe('custom');
  });

  // --- Validation ---

  it('should show error when end date is before start date', () => {
    const onChange = vi.fn();
    const { picker } = mountPicker(onChange, 'custom', '2026-03-10', '2026-03-15');

    const startInput = picker.querySelector('#drp-start') as HTMLInputElement;
    const endInput = picker.querySelector('#drp-end') as HTMLInputElement;

    startInput.value = '2026-03-20';
    endInput.value = '2026-03-10';
    endInput.dispatchEvent(new Event('change'));

    // Should NOT emit onChange
    expect(onChange).not.toHaveBeenCalled();

    // Should show error message
    const errorMsg = picker.querySelector('.form-error') as HTMLElement;
    expect(errorMsg.style.display).toBe('block');

    // Inputs should have error class
    expect(startInput.classList.contains('error')).toBe(true);
    expect(endInput.classList.contains('error')).toBe(true);
  });

  it('should clear error when valid dates are entered after an error', () => {
    const onChange = vi.fn();
    const { picker } = mountPicker(onChange, 'custom', '2026-03-10', '2026-03-15');

    const startInput = picker.querySelector('#drp-start') as HTMLInputElement;
    const endInput = picker.querySelector('#drp-end') as HTMLInputElement;

    // Set invalid range
    startInput.value = '2026-03-20';
    endInput.value = '2026-03-10';
    endInput.dispatchEvent(new Event('change'));

    // Now fix it
    startInput.value = '2026-03-05';
    startInput.dispatchEvent(new Event('change'));

    const errorMsg = picker.querySelector('.form-error') as HTMLElement;
    expect(errorMsg.style.display).toBe('none');
    expect(startInput.classList.contains('error')).toBe(false);
    expect(endInput.classList.contains('error')).toBe(false);
    expect(onChange).toHaveBeenCalledTimes(1);
  });

  it('should accept same start and end date (single day)', () => {
    const onChange = vi.fn();
    const { picker } = mountPicker(onChange, 'custom', '2026-03-10', '2026-03-10');

    const startInput = picker.querySelector('#drp-start') as HTMLInputElement;
    const endInput = picker.querySelector('#drp-end') as HTMLInputElement;

    startInput.value = '2026-03-15';
    endInput.value = '2026-03-15';
    endInput.dispatchEvent(new Event('change'));

    expect(onChange).toHaveBeenCalledTimes(1);
    const range = onChange.mock.calls[0][0] as DateRange;
    expect(range.start).toBe('2026-03-15');
    expect(range.end).toBe('2026-03-15');
  });

  // --- Selection display ---

  it('should have an aria-live selection display', () => {
    const { picker } = mountPicker();
    const display = picker.querySelector('[aria-live="polite"]');
    expect(display).toBeTruthy();
  });

  // --- Custom inputs shown when initialPreset is custom ---

  it('should show custom inputs when initialPreset is custom', () => {
    const { picker } = mountPicker(vi.fn(), 'custom', '2026-03-01', '2026-03-15');
    const startInput = picker.querySelector('#drp-start') as HTMLInputElement;
    const endInput = picker.querySelector('#drp-end') as HTMLInputElement;

    expect(startInput.value).toBe('2026-03-01');
    expect(endInput.value).toBe('2026-03-15');
    expect(startInput.parentElement?.style.display).toBe('flex');
  });

  // --- Accessibility ---

  it('should have role="group" and aria-label on root', () => {
    const { picker } = mountPicker();
    expect(picker.getAttribute('role')).toBe('group');
    expect(picker.getAttribute('aria-label')).toBe('Date range picker');
  });

  it('should set aria-pressed on preset buttons', () => {
    const { picker } = mountPicker(vi.fn(), 'today');
    const todayBtn = picker.querySelector('[data-preset="today"]');
    expect(todayBtn?.getAttribute('aria-pressed')).toBe('true');

    const thisWeekBtn = picker.querySelector('[data-preset="thisWeek"]');
    expect(thisWeekBtn?.getAttribute('aria-pressed')).toBe('false');
  });
});

describe('computePresetRange', () => {
  it('should return same start and end for "today"', () => {
    const range = computePresetRange('today');
    expect(range.start).toBe(range.end);
    // Should be a valid YYYY-MM-DD
    expect(range.start).toMatch(/^\d{4}-\d{2}-\d{2}$/);
  });

  it('should return a 7-day range for "thisWeek"', () => {
    const range = computePresetRange('thisWeek');
    const start = new Date(range.start);
    const end = new Date(range.end);
    const diff = (end.getTime() - start.getTime()) / (1000 * 60 * 60 * 24);
    expect(diff).toBe(6); // Mon-Sun = 6 day difference
    // Start should be a Monday
    expect(start.getDay()).toBe(1); // 1 = Monday
  });

  it('should return full month for "thisMonth"', () => {
    const range = computePresetRange('thisMonth');
    const start = new Date(range.start);
    const end = new Date(range.end);
    expect(start.getDate()).toBe(1);
    // End should be last day of same month
    expect(start.getMonth()).toBe(end.getMonth());
  });

  it('should return previous 7-day range for "lastWeek"', () => {
    const range = computePresetRange('lastWeek');
    const start = new Date(range.start);
    const end = new Date(range.end);
    const diff = (end.getTime() - start.getTime()) / (1000 * 60 * 60 * 24);
    expect(diff).toBe(6);
    expect(start.getDay()).toBe(1); // Monday

    // Should be before this week's start
    const thisWeek = computePresetRange('thisWeek');
    expect(range.end < thisWeek.start).toBe(true);
  });

  it('should return full previous month for "lastMonth"', () => {
    const range = computePresetRange('lastMonth');
    const start = new Date(range.start);
    const end = new Date(range.end);
    expect(start.getDate()).toBe(1);
    expect(start.getMonth()).toBe(end.getMonth());

    // Should be before current month
    const now = new Date();
    const expectedMonth = now.getMonth() === 0 ? 11 : now.getMonth() - 1;
    expect(start.getMonth()).toBe(expectedMonth);
  });
});
