import { handleInput } from './app';

export function initInputs(): void {
  const fromInput = document.getElementById('from-input') as HTMLInputElement;
  const toInput = document.getElementById('to-input') as HTMLInputElement;

  fromInput.addEventListener('input', () => {
    handleInput('from', fromInput.value);
  });

  toInput.addEventListener('input', () => {
    handleInput('to', toInput.value);
  });

  // Focus glow on cards
  fromInput.addEventListener('focus', () => {
    fromInput.closest('.card')?.classList.add('card-focused');
  });
  fromInput.addEventListener('blur', () => {
    fromInput.closest('.card')?.classList.remove('card-focused');
  });

  toInput.addEventListener('focus', () => {
    toInput.closest('.card')?.classList.add('card-focused');
  });
  toInput.addEventListener('blur', () => {
    toInput.closest('.card')?.classList.remove('card-focused');
  });
}
