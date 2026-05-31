import type { CategoryId } from '../converter/types';
import { switchCategory } from './app';

export function initTabs(): void {
  const tablist = document.querySelector<HTMLElement>('[role="tablist"]');
  const tabs = document.querySelectorAll<HTMLButtonElement>('.tab');
  const pill = document.querySelector<HTMLElement>('.tab-pill');
  const panel = document.getElementById('converter-panel');
  let switchTimeout: number | null = null;

  function activateTab(target: HTMLButtonElement, focus = false): void {
    const categoryId = target.dataset.category as CategoryId;

    // Skip if already active
    if (target.getAttribute('aria-selected') === 'true') return;

    // Update ARIA and tabindex (roving tabindex pattern)
    tabs.forEach((tab) => {
      tab.setAttribute('aria-selected', 'false');
      tab.setAttribute('tabindex', '-1');
      tab.classList.remove('tab-active');
    });
    target.setAttribute('aria-selected', 'true');
    target.setAttribute('tabindex', '0');
    target.classList.add('tab-active');

    if (focus) target.focus();

    // Update tabpanel labelledby
    if (panel) {
      panel.setAttribute('aria-labelledby', target.id);
    }

    // Slide pill to active tab
    if (pill) {
      pill.style.left = `${target.offsetLeft}px`;
      pill.style.width = `${target.offsetWidth}px`;
    }

    // Cancel any pending switch
    if (switchTimeout !== null) {
      clearTimeout(switchTimeout);
      const converter = document.querySelector('.converter') as HTMLElement;
      if (converter) {
        converter.classList.remove('content-switching');
        converter.style.transition = '';
        converter.style.opacity = '';
        converter.style.transform = '';
      }
    }

    const converter = document.querySelector('.converter') as HTMLElement;
    const reduced = window.matchMedia('(prefers-reduced-motion: reduce)').matches;

    if (reduced || !converter) {
      switchCategory(categoryId);
      return;
    }

    // Phase 1: Fade out + shift down
    converter.classList.add('content-switching');

    switchTimeout = window.setTimeout(() => {
      switchTimeout = null;

      // Switch content while invisible
      switchCategory(categoryId);

      // Phase 2: Position above without transition
      converter.classList.remove('content-switching');
      converter.style.transition = 'none';
      converter.style.opacity = '0';
      converter.style.transform = 'translateY(-4px)';

      // Force reflow so browser registers the state
      void converter.offsetHeight;

      // Phase 3: Re-enable transitions and animate to final position
      converter.style.transition = '';
      converter.style.opacity = '';
      converter.style.transform = '';
    }, 250);
  }

  // Click handlers
  tabs.forEach((tab) => {
    tab.addEventListener('click', () => activateTab(tab));
  });

  // Keyboard navigation: Arrow keys, Home, End
  if (tablist) {
    tablist.addEventListener('keydown', (e: KeyboardEvent) => {
      const tabArray = Array.from(tabs);
      const currentIndex = tabArray.findIndex((t) => t === document.activeElement);
      if (currentIndex === -1) return;

      let nextIndex: number | null = null;

      switch (e.key) {
        case 'ArrowRight':
        case 'ArrowDown':
          nextIndex = (currentIndex + 1) % tabArray.length;
          break;
        case 'ArrowLeft':
        case 'ArrowUp':
          nextIndex = (currentIndex - 1 + tabArray.length) % tabArray.length;
          break;
        case 'Home':
          nextIndex = 0;
          break;
        case 'End':
          nextIndex = tabArray.length - 1;
          break;
        default:
          return;
      }

      e.preventDefault();
      activateTab(tabArray[nextIndex], true);
    });
  }

  // Initialize pill position on the first active tab
  const activeTab = document.querySelector<HTMLButtonElement>('.tab[aria-selected="true"]');
  if (activeTab && pill) {
    requestAnimationFrame(() => {
      pill.style.left = `${activeTab.offsetLeft}px`;
      pill.style.width = `${activeTab.offsetWidth}px`;
    });
  }
}
