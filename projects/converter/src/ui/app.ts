import type { CategoryId, ConversionState } from '../converter/types';
import { getCategory } from '../converter/units';
import { convert, getQuickReference } from '../converter/engine';
import { formatResult, sanitizeNumericInput } from '../utils/format';
import { initTabs } from './tabs';
import { renderSelector } from './selector';
import { initInputs } from './input';
import { initSwap } from './swap';
import { initTheme } from './theme';
import { animateResultChange } from './animations';

const state: ConversionState = {
  category: 'length',
  fromUnit: 'km',
  toUnit: 'mi',
  fromValue: '',
  toValue: '',
  inputDirection: 'from',
};

let announceTimeout: number | null = null;

export function getState(): ConversionState {
  return state;
}

export function runConversion(): void {
  const sourceValue = state.inputDirection === 'from' ? state.fromValue : state.toValue;

  if (sourceValue === '' || sourceValue === '-') {
    if (state.inputDirection === 'from') {
      state.toValue = '';
    } else {
      state.fromValue = '';
    }
    updateResultDisplay();
    updateQuickRef();
    announceResult();
    return;
  }

  const num = parseFloat(sourceValue);
  if (!Number.isFinite(num)) {
    if (state.inputDirection === 'from') {
      state.toValue = '';
    } else {
      state.fromValue = '';
    }
    updateResultDisplay();
    updateQuickRef();
    announceResult();
    return;
  }

  if (state.inputDirection === 'from') {
    const result = convert(state.category, state.fromUnit, state.toUnit, num);
    state.toValue = formatResult(result);
  } else {
    const result = convert(state.category, state.toUnit, state.fromUnit, num);
    state.fromValue = formatResult(result);
  }

  updateResultDisplay();
  updateQuickRef();
  announceResult();
}

function updateResultDisplay(): void {
  const fromInput = document.getElementById('from-input') as HTMLInputElement;
  const toInput = document.getElementById('to-input') as HTMLInputElement;

  if (state.inputDirection === 'from') {
    const oldValue = toInput.value;
    toInput.value = state.toValue;
    if (state.toValue && oldValue !== state.toValue && !toInput.classList.contains('value-swapping')) {
      animateResultChange(toInput);
    }
  } else {
    const oldValue = fromInput.value;
    fromInput.value = state.fromValue;
    if (state.fromValue && oldValue !== state.fromValue && !fromInput.classList.contains('value-swapping')) {
      animateResultChange(fromInput);
    }
  }
}

function updateQuickRef(): void {
  const el = document.getElementById('quick-ref');
  if (el) {
    el.textContent = getQuickReference(state.category, state.fromUnit, state.toUnit);
  }
}

function announceResult(): void {
  const srEl = document.getElementById('sr-announcement');
  if (!srEl) return;

  // Debounce announcements to avoid spamming screen readers during typing
  if (announceTimeout !== null) {
    clearTimeout(announceTimeout);
  }

  announceTimeout = window.setTimeout(() => {
    announceTimeout = null;
    const cat = getCategory(state.category);
    const fromUnit = cat.units.find((u) => u.id === state.fromUnit);
    const toUnit = cat.units.find((u) => u.id === state.toUnit);

    if (state.toValue && state.fromValue) {
      const resultVal = state.inputDirection === 'from' ? state.toValue : state.fromValue;
      const resultUnit = state.inputDirection === 'from' ? toUnit : fromUnit;
      srEl.textContent = `${resultVal} ${resultUnit?.label ?? ''}`;
    } else {
      srEl.textContent = '';
    }
  }, 500);
}

export function updateSymbols(): void {
  const cat = getCategory(state.category);
  const fromUnit = cat.units.find((u) => u.id === state.fromUnit);
  const toUnit = cat.units.find((u) => u.id === state.toUnit);

  const fromSymbol = document.getElementById('from-symbol');
  const toSymbol = document.getElementById('to-symbol');
  if (fromSymbol) fromSymbol.textContent = fromUnit?.symbol ?? '';
  if (toSymbol) toSymbol.textContent = toUnit?.symbol ?? '';
}

export function switchCategory(categoryId: CategoryId): void {
  const cat = getCategory(categoryId);
  state.category = categoryId;
  state.fromUnit = cat.defaultFrom;
  state.toUnit = cat.defaultTo;
  state.inputDirection = 'from';

  // Re-render selectors
  renderSelector('from-pills', cat, state.fromUnit, (unitId) => {
    state.fromUnit = unitId;
    updateSymbols();
    runConversion();
  });
  renderSelector('to-pills', cat, state.toUnit, (unitId) => {
    state.toUnit = unitId;
    updateSymbols();
    runConversion();
  });

  updateSymbols();

  // Re-convert with current input
  const fromInput = document.getElementById('from-input') as HTMLInputElement;
  const toInput = document.getElementById('to-input') as HTMLInputElement;

  state.fromValue = fromInput.value;
  state.toValue = '';
  state.inputDirection = 'from';
  toInput.value = '';

  runConversion();
}

export function handleSwap(): void {
  // Swap units
  const tmpUnit = state.fromUnit;
  state.fromUnit = state.toUnit;
  state.toUnit = tmpUnit;

  // Swap values
  const tmpVal = state.fromValue;
  state.fromValue = state.toValue;
  state.toValue = tmpVal;

  const fromInput = document.getElementById('from-input') as HTMLInputElement;
  const toInput = document.getElementById('to-input') as HTMLInputElement;
  fromInput.value = state.fromValue;
  toInput.value = state.toValue;

  state.inputDirection = 'from';

  // Re-render selectors with swapped active states
  const cat = getCategory(state.category);
  renderSelector('from-pills', cat, state.fromUnit, (unitId) => {
    state.fromUnit = unitId;
    updateSymbols();
    runConversion();
  });
  renderSelector('to-pills', cat, state.toUnit, (unitId) => {
    state.toUnit = unitId;
    updateSymbols();
    runConversion();
  });

  updateSymbols();
  runConversion();
}

export function handleInput(direction: 'from' | 'to', rawValue: string): void {
  const sanitized = sanitizeNumericInput(rawValue);

  if (sanitized === null) {
    // Revert to last valid value
    const input = document.getElementById(
      direction === 'from' ? 'from-input' : 'to-input'
    ) as HTMLInputElement;
    input.value = direction === 'from' ? state.fromValue : state.toValue;
    return;
  }

  state.inputDirection = direction;
  if (direction === 'from') {
    state.fromValue = sanitized;
  } else {
    state.toValue = sanitized;
  }

  runConversion();
}

export function initApp(): void {
  const cat = getCategory(state.category);

  initTheme();
  initTabs();

  renderSelector('from-pills', cat, state.fromUnit, (unitId) => {
    state.fromUnit = unitId;
    updateSymbols();
    runConversion();
  });
  renderSelector('to-pills', cat, state.toUnit, (unitId) => {
    state.toUnit = unitId;
    updateSymbols();
    runConversion();
  });

  initInputs();
  initSwap();
  updateSymbols();
  updateQuickRef();
}
