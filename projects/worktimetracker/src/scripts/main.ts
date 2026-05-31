// WorkTime Tracker - Main Entry Point

import { createIcons } from 'lucide';
import { store } from './data/store';
import { projectService } from './data/projectService';
import { timeEntryService } from './data/timeEntryService';
import { initTimerWidget } from './components/timerWidget';
import { initSettings } from './components/settings';
import { initKeyboardShortcuts } from './utils/keyboard';
import { initRouter } from './router';

function init(): void {
  // Initialize localStorage cache
  store.init();

  // Ensure the default "General" project exists
  projectService.ensureGeneralProject();

  // Initialize Lucide icons for the shell (sidebar, timer bar, mobile nav)
  createIcons();

  // Initialize the persistent timer bar widget (works across all views)
  initTimerWidget();

  // Wire up Settings gear icon (sidebar + mobile)
  initSettings();

  // Initialize global keyboard shortcuts
  initKeyboardShortcuts();

  // Initialize mobile hamburger menu
  initMobileMenu();

  // Persist active timer state every 30s for crash recovery
  setInterval(() => {
    timeEntryService.persistTimerState();
  }, 30000);

  // Start hash-based router (loads initial view)
  initRouter();
}

function initMobileMenu(): void {
  const hamburger = document.getElementById('hamburger-btn');
  const sidebar = document.getElementById('sidebar');
  const overlay = document.getElementById('sidebar-overlay');

  if (!hamburger || !sidebar || !overlay) return;

  function openSidebar(): void {
    sidebar!.classList.add('open');
    overlay!.classList.remove('hidden');
  }

  function closeSidebar(): void {
    sidebar!.classList.remove('open');
    overlay!.classList.add('hidden');
  }

  hamburger.addEventListener('click', () => {
    if (sidebar.classList.contains('open')) {
      closeSidebar();
    } else {
      openSidebar();
    }
  });

  overlay.addEventListener('click', closeSidebar);

  // Close sidebar when a nav item is clicked (mobile)
  sidebar.querySelectorAll('.nav-item[data-route]').forEach((item) => {
    item.addEventListener('click', closeSidebar);
  });

  // Close sidebar when settings is clicked
  const settingsBtn = document.getElementById('settings-btn');
  if (settingsBtn) {
    settingsBtn.addEventListener('click', closeSidebar);
  }
}

// Run on DOM ready
if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', init);
} else {
  init();
}
