let activeAnimation: Animation | null = null;

export function animateResultChange(el: HTMLElement): void {
  if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) return;

  // Cancel any in-flight animation before starting a new one (handles rapid typing)
  if (activeAnimation) {
    activeAnimation.cancel();
    activeAnimation = null;
  }

  activeAnimation = el.animate(
    [
      { opacity: 0.4, transform: 'scale(0.97)' },
      { opacity: 1, transform: 'scale(1)' },
    ],
    {
      duration: 200,
      easing: 'ease-out',
      fill: 'none',
    }
  );

  activeAnimation.onfinish = () => {
    activeAnimation = null;
  };
  activeAnimation.oncancel = () => {
    activeAnimation = null;
  };
}
