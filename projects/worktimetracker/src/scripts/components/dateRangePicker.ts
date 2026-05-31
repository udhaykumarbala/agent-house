/**
 * DateRangePicker Component
 *
 * Reusable date range picker with preset buttons and custom range inputs.
 * Presets: Today, This Week, This Month, Last Week, Last Month, Custom.
 * Used by Timesheet and Reports views.
 *
 * - Validates end >= start for custom ranges
 * - Emits selected range via callback
 * - Displays current selection clearly (`.date-preset.selected`)
 * - Styled consistent with design system
 * - Uses textContent for user data (XSS prevention)
 */

export interface DateRange {
  start: string; // YYYY-MM-DD
  end: string;   // YYYY-MM-DD
  preset: PresetKey;
}

export type PresetKey = 'today' | 'thisWeek' | 'thisMonth' | 'lastWeek' | 'lastMonth' | 'custom';

export interface DateRangePickerOptions {
  onChange: (range: DateRange) => void;
  initialPreset?: PresetKey;
  initialStart?: string;
  initialEnd?: string;
}

interface PresetDef {
  key: PresetKey;
  label: string;
}

const PRESETS: PresetDef[] = [
  { key: 'today', label: 'Today' },
  { key: 'thisWeek', label: 'This Week' },
  { key: 'thisMonth', label: 'This Month' },
  { key: 'lastWeek', label: 'Last Week' },
  { key: 'lastMonth', label: 'Last Month' },
  { key: 'custom', label: 'Custom' },
];

/** Get Monday of the current week (ISO week: Mon–Sun). */
function getWeekStart(date: Date): Date {
  const d = new Date(date);
  const day = d.getDay();
  // Sunday = 0, Monday = 1, etc.
  const diff = day === 0 ? 6 : day - 1;
  d.setDate(d.getDate() - diff);
  return d;
}

function toYMD(date: Date): string {
  const y = date.getFullYear();
  const m = String(date.getMonth() + 1).padStart(2, '0');
  const d = String(date.getDate()).padStart(2, '0');
  return `${y}-${m}-${d}`;
}

/** Compute start/end dates for a given preset key. */
export function computePresetRange(key: PresetKey): { start: string; end: string } {
  const now = new Date();

  switch (key) {
    case 'today': {
      const today = toYMD(now);
      return { start: today, end: today };
    }
    case 'thisWeek': {
      const weekStart = getWeekStart(now);
      const weekEnd = new Date(weekStart);
      weekEnd.setDate(weekEnd.getDate() + 6);
      return { start: toYMD(weekStart), end: toYMD(weekEnd) };
    }
    case 'thisMonth': {
      const monthStart = new Date(now.getFullYear(), now.getMonth(), 1);
      const monthEnd = new Date(now.getFullYear(), now.getMonth() + 1, 0);
      return { start: toYMD(monthStart), end: toYMD(monthEnd) };
    }
    case 'lastWeek': {
      const thisWeekStart = getWeekStart(now);
      const lastWeekStart = new Date(thisWeekStart);
      lastWeekStart.setDate(lastWeekStart.getDate() - 7);
      const lastWeekEnd = new Date(lastWeekStart);
      lastWeekEnd.setDate(lastWeekEnd.getDate() + 6);
      return { start: toYMD(lastWeekStart), end: toYMD(lastWeekEnd) };
    }
    case 'lastMonth': {
      const lastMonthStart = new Date(now.getFullYear(), now.getMonth() - 1, 1);
      const lastMonthEnd = new Date(now.getFullYear(), now.getMonth(), 0);
      return { start: toYMD(lastMonthStart), end: toYMD(lastMonthEnd) };
    }
    case 'custom':
    default: {
      const today = toYMD(now);
      return { start: today, end: today };
    }
  }
}

export function createDateRangePicker(options: DateRangePickerOptions): HTMLElement {
  const { onChange, initialPreset = 'thisWeek', initialStart, initialEnd } = options;

  let activePreset: PresetKey = initialPreset;
  let currentRange: DateRange;

  // Compute the initial range
  if (initialPreset === 'custom' && initialStart && initialEnd) {
    currentRange = { start: initialStart, end: initialEnd, preset: 'custom' };
  } else {
    const computed = computePresetRange(initialPreset);
    currentRange = { ...computed, preset: initialPreset };
  }

  // Root container
  const root = document.createElement('div');
  root.className = 'flex flex-wrap items-center gap-2';
  root.setAttribute('role', 'group');
  root.setAttribute('aria-label', 'Date range picker');

  // Preset buttons container
  const presetsRow = document.createElement('div');
  presetsRow.className = 'flex flex-wrap gap-2';

  const presetButtons: Map<PresetKey, HTMLButtonElement> = new Map();

  PRESETS.forEach((preset) => {
    const btn = document.createElement('button');
    btn.type = 'button';
    btn.className = 'date-preset';
    btn.textContent = preset.label;
    btn.dataset.preset = preset.key;
    btn.setAttribute('aria-pressed', String(preset.key === activePreset));

    if (preset.key === activePreset) {
      btn.classList.add('selected');
    }

    btn.addEventListener('click', () => handlePresetClick(preset.key));
    presetButtons.set(preset.key, btn);
    presetsRow.appendChild(btn);
  });

  // Custom date inputs container (hidden by default)
  const customContainer = document.createElement('div');
  customContainer.className = 'flex items-center gap-2 mt-2 sm:mt-0';
  customContainer.style.display = activePreset === 'custom' ? 'flex' : 'none';

  const startLabel = document.createElement('label');
  startLabel.className = 'text-content-secondary text-xs font-medium';
  startLabel.textContent = 'From';
  startLabel.setAttribute('for', 'drp-start');

  const startInput = document.createElement('input');
  startInput.type = 'date';
  startInput.id = 'drp-start';
  startInput.className = 'input';
  startInput.style.cssText = 'width: 150px; padding: 6px 10px; font-size: 13px;';
  startInput.value = currentRange.start;

  const toLabel = document.createElement('label');
  toLabel.className = 'text-content-secondary text-xs font-medium';
  toLabel.textContent = 'To';
  toLabel.setAttribute('for', 'drp-end');

  const endInput = document.createElement('input');
  endInput.type = 'date';
  endInput.id = 'drp-end';
  endInput.className = 'input';
  endInput.style.cssText = 'width: 150px; padding: 6px 10px; font-size: 13px;';
  endInput.value = currentRange.end;

  // Validation error message
  const errorMsg = document.createElement('span');
  errorMsg.className = 'form-error';
  errorMsg.style.display = 'none';
  errorMsg.textContent = 'End date must be on or after start date';

  customContainer.appendChild(startLabel);
  customContainer.appendChild(startInput);
  customContainer.appendChild(toLabel);
  customContainer.appendChild(endInput);
  customContainer.appendChild(errorMsg);

  // Current selection display
  const selectionDisplay = document.createElement('span');
  selectionDisplay.className = 'text-content-secondary text-xs ml-2 hidden sm:inline';
  selectionDisplay.setAttribute('aria-live', 'polite');
  updateSelectionDisplay();

  root.appendChild(presetsRow);
  root.appendChild(customContainer);
  root.appendChild(selectionDisplay);

  // Event handlers for custom date inputs
  startInput.addEventListener('change', handleCustomDateChange);
  endInput.addEventListener('change', handleCustomDateChange);

  function handlePresetClick(key: PresetKey): void {
    activePreset = key;

    // Update button states
    presetButtons.forEach((btn, k) => {
      const isSelected = k === key;
      btn.classList.toggle('selected', isSelected);
      btn.setAttribute('aria-pressed', String(isSelected));
    });

    if (key === 'custom') {
      customContainer.style.display = 'flex';
      // Keep current custom dates or default to current range
      startInput.value = currentRange.start;
      endInput.value = currentRange.end;
    } else {
      customContainer.style.display = 'none';
      errorMsg.style.display = 'none';
      const computed = computePresetRange(key);
      currentRange = { ...computed, preset: key };
      updateSelectionDisplay();
      onChange(currentRange);
    }
  }

  function handleCustomDateChange(): void {
    const start = startInput.value;
    const end = endInput.value;

    if (!start || !end) return;

    // Validate end >= start
    if (end < start) {
      errorMsg.style.display = 'block';
      startInput.classList.add('error');
      endInput.classList.add('error');
      return;
    }

    errorMsg.style.display = 'none';
    startInput.classList.remove('error');
    endInput.classList.remove('error');

    currentRange = { start, end, preset: 'custom' };
    updateSelectionDisplay();
    onChange(currentRange);
  }

  function updateSelectionDisplay(): void {
    const presetDef = PRESETS.find((p) => p.key === currentRange.preset);
    if (currentRange.preset === 'custom') {
      selectionDisplay.textContent = `${currentRange.start} — ${currentRange.end}`;
    } else {
      selectionDisplay.textContent = presetDef?.label ?? '';
    }
  }

  return root;
}
