const STORAGE_KEY = 'unitshift-theme';

function getSystemPreference(): 'dark' | 'light' {
  return window.matchMedia('(prefers-color-scheme: light)').matches ? 'light' : 'dark';
}

function applyTheme(theme: 'dark' | 'light'): void {
  const root = document.documentElement;

  if (theme === 'light') {
    root.setAttribute('data-theme', 'light');
  } else {
    root.removeAttribute('data-theme');
  }

  updateToggleIcon(theme);
}

function updateToggleIcon(theme: 'dark' | 'light'): void {
  const sunIcon = document.getElementById('theme-icon-sun');
  const moonIcon = document.getElementById('theme-icon-moon');
  if (!sunIcon || !moonIcon) return;

  if (theme === 'light') {
    sunIcon.style.display = 'none';
    moonIcon.style.display = 'block';
  } else {
    sunIcon.style.display = 'block';
    moonIcon.style.display = 'none';
  }
}

function getCurrentTheme(): 'dark' | 'light' {
  return document.documentElement.hasAttribute('data-theme') ? 'light' : 'dark';
}

export function initTheme(): void {
  const stored = localStorage.getItem(STORAGE_KEY) as 'dark' | 'light' | null;
  const theme = stored ?? getSystemPreference();

  // Apply immediately (no transition on initial load)
  applyTheme(theme);

  // Enable transitions after first paint so initial load has no flash
  requestAnimationFrame(() => {
    requestAnimationFrame(() => {
      document.documentElement.classList.add('theme-transitions');
    });
  });

  const btn = document.getElementById('theme-toggle');
  if (!btn) return;

  btn.addEventListener('click', () => {
    const current = getCurrentTheme();
    const next = current === 'dark' ? 'light' : 'dark';
    localStorage.setItem(STORAGE_KEY, next);
    applyTheme(next);
    btn.setAttribute('aria-label', next === 'dark' ? 'Switch to light mode' : 'Switch to dark mode');
  });
}
