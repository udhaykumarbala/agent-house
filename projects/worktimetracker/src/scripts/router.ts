/**
 * Hash-based router for SPA view switching.
 * - 5 routes: dashboard, employees, timesheet, projects, reports
 * - Default route: #/dashboard
 * - Unknown routes redirect to #/dashboard
 * - Active sidebar item highlighted on route change
 * - Only active view is in DOM (lazy rendering)
 */

import { createIcons } from 'lucide';
import type { ViewModule } from './types';

type ViewLoader = () => Promise<{ render: ViewModule['render']; destroy?: () => void }>;

interface RouteConfig {
  path: string;
  load: ViewLoader;
}

const ROUTE_MAP: RouteConfig[] = [
  { path: '/dashboard', load: () => import('./views/dashboard') },
  { path: '/employees', load: () => import('./views/employees') },
  { path: '/timesheet', load: () => import('./views/timesheet') },
  { path: '/projects', load: () => import('./views/projects') },
  { path: '/reports', load: () => import('./views/reports') },
];

const routes: Map<string, RouteConfig> = new Map();
let currentRoute: string | null = null;
let currentDestroy: (() => void) | undefined;
const DEFAULT_ROUTE = '/dashboard';

function getHashPath(): string {
  const hash = window.location.hash;
  if (!hash || hash === '#' || hash === '#/') {
    return DEFAULT_ROUTE;
  }
  return hash.substring(1);
}

function updateActiveNav(path: string): void {
  // Desktop sidebar
  document.querySelectorAll('#sidebar .nav-item[data-route]').forEach((el) => {
    const route = el.getAttribute('data-route');
    if (route && path === `/${route}`) {
      el.classList.add('active');
      el.setAttribute('aria-current', 'page');
    } else {
      el.classList.remove('active');
      el.removeAttribute('aria-current');
    }
  });

  // Mobile nav
  document
    .querySelectorAll('#mobile-nav .mobile-nav-item[data-route]')
    .forEach((el) => {
      const route = el.getAttribute('data-route');
      if (route && path === `/${route}`) {
        el.classList.add('active');
      } else {
        el.classList.remove('active');
      }
    });
}

async function navigateTo(path: string): Promise<void> {
  const route = routes.get(path);

  if (!route) {
    window.location.hash = `#${DEFAULT_ROUTE}`;
    return;
  }

  if (currentRoute === path) return;
  currentRoute = path;

  // Destroy previous view if it has cleanup
  if (currentDestroy) {
    currentDestroy();
    currentDestroy = undefined;
  }

  updateActiveNav(path);

  const container = document.getElementById('view-container');
  if (!container) return;

  // Fade out
  container.classList.remove('view-active');
  container.classList.add('view-enter');

  // Lazy-load and render view
  const viewModule = await route.load();
  container.innerHTML = '';
  await viewModule.render(container);
  currentDestroy = viewModule.destroy;

  // Re-initialize Lucide icons in the new view
  try {
    createIcons();
  } catch {
    // Icons not yet loaded, that's fine
  }

  // Fade in
  requestAnimationFrame(() => {
    container.classList.remove('view-enter');
    container.classList.add('view-active');
  });
}

function handleHashChange(): void {
  const path = getHashPath();
  navigateTo(path);
}

/** Initialize the hash-based router and navigate to the initial route. */
export function initRouter(): void {
  // Register all routes
  for (const route of ROUTE_MAP) {
    routes.set(route.path, route);
  }

  window.addEventListener('hashchange', handleHashChange);

  // Navigate to current hash or default
  if (
    !window.location.hash ||
    window.location.hash === '#' ||
    window.location.hash === '#/'
  ) {
    window.location.hash = `#${DEFAULT_ROUTE}`;
  } else {
    navigateTo(getHashPath());
  }
}

/** Programmatic navigation. */
export function navigate(path: string): void {
  window.location.hash = `#${path}`;
}

/** Get current route path. */
export function getCurrentRoute(): string | null {
  return currentRoute;
}
