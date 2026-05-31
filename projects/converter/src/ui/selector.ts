import type { UnitCategory } from '../converter/types';

export function renderSelector(
  containerId: string,
  category: UnitCategory,
  activeUnitId: string,
  onSelect: (unitId: string) => void
): void {
  const container = document.getElementById(containerId);
  if (!container) return;

  // Clear existing pills
  while (container.firstChild) {
    container.removeChild(container.firstChild);
  }

  category.units.forEach((unit) => {
    const pill = document.createElement('button');
    pill.type = 'button';
    pill.className = 'pill';
    pill.setAttribute('role', 'radio');
    pill.setAttribute('aria-checked', unit.id === activeUnitId ? 'true' : 'false');
    pill.setAttribute('aria-label', unit.label);
    pill.dataset.unitId = unit.id;
    pill.textContent = unit.symbol;

    // Roving tabindex: only active pill is tabbable
    pill.setAttribute('tabindex', unit.id === activeUnitId ? '0' : '-1');

    if (unit.id === activeUnitId) {
      pill.classList.add('pill-active');
    }

    pill.addEventListener('click', () => {
      selectPill(container, pill, onSelect);
    });

    container.appendChild(pill);
  });

  // Arrow key navigation within radiogroup
  container.addEventListener('keydown', (e: KeyboardEvent) => {
    const pills = Array.from(container.querySelectorAll<HTMLButtonElement>('.pill'));
    const currentIndex = pills.findIndex((p) => p === document.activeElement);
    if (currentIndex === -1) return;

    let nextIndex: number | null = null;

    switch (e.key) {
      case 'ArrowRight':
      case 'ArrowDown':
        nextIndex = (currentIndex + 1) % pills.length;
        break;
      case 'ArrowLeft':
      case 'ArrowUp':
        nextIndex = (currentIndex - 1 + pills.length) % pills.length;
        break;
      case 'Home':
        nextIndex = 0;
        break;
      case 'End':
        nextIndex = pills.length - 1;
        break;
      default:
        return;
    }

    e.preventDefault();
    const nextPill = pills[nextIndex];
    selectPill(container, nextPill, onSelect);
    nextPill.focus();
  });
}

function selectPill(
  container: HTMLElement,
  pill: HTMLButtonElement,
  onSelect: (unitId: string) => void
): void {
  // Deactivate all pills in this container
  container.querySelectorAll<HTMLButtonElement>('.pill').forEach((p) => {
    p.classList.remove('pill-active');
    p.setAttribute('aria-checked', 'false');
    p.setAttribute('tabindex', '-1');
  });

  // Activate clicked pill
  pill.classList.add('pill-active');
  pill.setAttribute('aria-checked', 'true');
  pill.setAttribute('tabindex', '0');

  // Trigger press animation
  pill.classList.remove('pill-press');
  void pill.offsetWidth;
  pill.classList.add('pill-press');
  pill.addEventListener('animationend', () => {
    pill.classList.remove('pill-press');
  }, { once: true });

  const unitId = pill.dataset.unitId;
  if (unitId) {
    onSelect(unitId);
  }
}
