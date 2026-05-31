import { handleSwap } from './app';

let swapRotation = 0;
let swapLocked = false;

export function initSwap(): void {
  const btn = document.getElementById('swap-btn');
  if (!btn) return;

  const icon = btn.querySelector('svg') as SVGElement;
  const fromInput = document.getElementById('from-input') as HTMLInputElement;
  const toInput = document.getElementById('to-input') as HTMLInputElement;

  btn.addEventListener('click', () => {
    if (swapLocked) return;

    // Rotate icon (accumulates 180° each click)
    swapRotation += 180;
    if (icon) {
      icon.style.transform = `rotate(${swapRotation}deg)`;
    }

    const reduced = window.matchMedia('(prefers-reduced-motion: reduce)').matches;
    if (reduced) {
      handleSwap();
      return;
    }

    // Cross-fade values: fade out → swap → fade in
    swapLocked = true;
    fromInput.classList.add('value-swapping');
    toInput.classList.add('value-swapping');

    // Wait for fade-out (100ms), then swap values and fade back in
    setTimeout(() => {
      handleSwap();

      // Remove swapping class in next frame so the transition to opacity 1 plays
      requestAnimationFrame(() => {
        fromInput.classList.remove('value-swapping');
        toInput.classList.remove('value-swapping');

        // Unlock after fade-in completes
        setTimeout(() => {
          swapLocked = false;
        }, 100);
      });
    }, 100);
  });
}
