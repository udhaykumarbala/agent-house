/**
 * Inventory Reorder Alert System
 * MVP Implementation with localStorage persistence
 */

// === State Management ===
const state = {
  user: null,
  items: [],
  alerts: [],
  settings: {
    emailAlerts: true,
    alertFrequency: 'immediate',
    dailyDigest: false
  },
  currentPage: 'dashboard',
  currentItem: null,
  editingItem: null,
  quickUpdateMode: 'add' // 'add', 'remove', or 'set'
};

// === Local Storage Keys ===
const STORAGE_KEYS = {
  USER: 'reorder_user',
  ITEMS: 'reorder_items',
  ALERTS: 'reorder_alerts',
  SETTINGS: 'reorder_settings'
};

// === DOM Elements ===
const elements = {
  authScreen: document.getElementById('auth-screen'),
  mainApp: document.getElementById('main-app'),
  authForm: document.getElementById('auth-form'),
  authEmail: document.getElementById('auth-email'),
  authPassword: document.getElementById('auth-password'),
  authSubmit: document.getElementById('auth-submit'),
  authSwitch: document.getElementById('auth-switch'),
  forgotLink: document.getElementById('forgot-link'),
  pageTitle: document.getElementById('page-title'),
  bottomNav: document.querySelector('.bottom-nav'),
  fab: document.getElementById('fab-add'),
  quickUpdateOverlay: document.getElementById('quick-update-overlay'),
  quickUpdateSheet: document.getElementById('quick-update-sheet'),
  quickUpdateContent: document.getElementById('quick-update-content'),
  toast: document.getElementById('toast'),
  toastMessage: document.getElementById('toast-message'),
  toastIcon: document.getElementById('toast-icon'),
  confirmDialog: document.getElementById('confirm-dialog'),
  dialogTitle: document.getElementById('dialog-title'),
  dialogMessage: document.getElementById('dialog-message'),
  dialogCancel: document.getElementById('dialog-cancel'),
  dialogConfirm: document.getElementById('dialog-confirm'),
  alertBadge: document.getElementById('alert-badge'),
  userBtn: document.getElementById('user-btn'),
  menuBtn: document.getElementById('menu-btn')
};

// === Initialization ===
function init() {
  loadFromStorage();
  setupEventListeners();
  checkAuth();
}

function loadFromStorage() {
  const user = localStorage.getItem(STORAGE_KEYS.USER);
  const items = localStorage.getItem(STORAGE_KEYS.ITEMS);
  const alerts = localStorage.getItem(STORAGE_KEYS.ALERTS);
  const settings = localStorage.getItem(STORAGE_KEYS.SETTINGS);

  if (user) state.user = JSON.parse(user);
  if (items) state.items = JSON.parse(items);
  if (alerts) state.alerts = JSON.parse(alerts);
  if (settings) state.settings = { ...state.settings, ...JSON.parse(settings) };
}

function saveToStorage() {
  localStorage.setItem(STORAGE_KEYS.USER, JSON.stringify(state.user));
  localStorage.setItem(STORAGE_KEYS.ITEMS, JSON.stringify(state.items));
  localStorage.setItem(STORAGE_KEYS.ALERTS, JSON.stringify(state.alerts));
  localStorage.setItem(STORAGE_KEYS.SETTINGS, JSON.stringify(state.settings));
}

function checkAuth() {
  if (state.user) {
    showMainApp();
  } else {
    showAuthScreen();
  }
}

// === Event Listeners ===
function setupEventListeners() {
  // Auth form
  elements.authForm.addEventListener('submit', handleAuth);
  elements.authSwitch.addEventListener('click', toggleAuthMode);
  elements.forgotLink.addEventListener('click', handleForgotPassword);

  // Password toggle
  document.querySelector('.password-toggle').addEventListener('click', togglePasswordVisibility);

  // Navigation
  elements.bottomNav.addEventListener('click', handleNavigation);

  // FAB
  elements.fab.addEventListener('click', () => navigateTo('item-form'));

  // Quick update sheet
  elements.quickUpdateOverlay.addEventListener('click', closeQuickUpdate);

  // User button (logout)
  elements.userBtn.addEventListener('click', handleUserMenu);

  // Dialog buttons
  elements.dialogCancel.addEventListener('click', closeDialog);
}

// === Auth Handlers ===
let isLoginMode = true;

function handleAuth(e) {
  e.preventDefault();
  const email = elements.authEmail.value.trim();
  const password = elements.authPassword.value;

  if (!email || !password) {
    showToast('Please fill in all fields', 'error');
    return;
  }

  if (isLoginMode) {
    // Login - check if user exists
    if (state.user && state.user.email === email) {
      // Simple password check (in production, use proper auth)
      if (state.user.password === password) {
        showToast('Welcome back!', 'success');
        showMainApp();
      } else {
        showToast('Invalid password', 'error');
      }
    } else if (state.user) {
      showToast('Account not found', 'error');
    } else {
      showToast('No account found. Please sign up.', 'error');
    }
  } else {
    // Sign up
    state.user = {
      email,
      password, // In production, hash this!
      createdAt: new Date().toISOString()
    };
    saveToStorage();
    showToast('Account created!', 'success');
    showMainApp();
  }
}

function toggleAuthMode() {
  isLoginMode = !isLoginMode;
  elements.authSubmit.textContent = isLoginMode ? 'Log In' : 'Create Account';
  elements.authSwitch.textContent = isLoginMode ? 'Create Account' : 'Log In';
  elements.authForm.reset();
}

function handleForgotPassword(e) {
  e.preventDefault();
  showToast('Password reset email sent', 'success');
}

function togglePasswordVisibility() {
  const input = elements.authPassword;
  input.type = input.type === 'password' ? 'text' : 'password';
}

function handleUserMenu() {
  showDialog(
    'Log Out',
    'Are you sure you want to log out?',
    () => {
      state.user = null;
      localStorage.removeItem(STORAGE_KEYS.USER);
      showAuthScreen();
      showToast('Logged out', 'success');
    },
    'Log Out'
  );
}

// === Navigation ===
function handleNavigation(e) {
  const navItem = e.target.closest('.bottom-nav-item');
  if (!navItem) return;

  e.preventDefault();
  const page = navItem.dataset.page;
  navigateTo(page);
}

function navigateTo(page, data = null) {
  // Hide all pages
  document.querySelectorAll('.page').forEach(p => p.classList.remove('active'));

  // Update nav active state
  document.querySelectorAll('.bottom-nav-item').forEach(item => {
    item.classList.toggle('active', item.dataset.page === page);
  });

  // Show FAB on appropriate pages
  const showFab = ['dashboard', 'items'].includes(page);
  elements.fab.classList.toggle('hidden', !showFab);

  // Set page title and render
  state.currentPage = page;

  switch (page) {
    case 'dashboard':
      elements.pageTitle.textContent = 'Dashboard';
      renderDashboard();
      break;
    case 'items':
      elements.pageTitle.textContent = 'All Items';
      renderItemsList();
      break;
    case 'alerts':
      elements.pageTitle.textContent = 'Alerts';
      renderAlerts();
      break;
    case 'settings':
      elements.pageTitle.textContent = 'Settings';
      renderSettings();
      break;
    case 'item-detail':
      state.currentItem = data;
      elements.pageTitle.textContent = data.name;
      renderItemDetail();
      break;
    case 'item-form':
      state.editingItem = data;
      elements.pageTitle.textContent = data ? 'Edit Item' : 'Add Item';
      renderItemForm();
      break;
  }

  // Show page
  const pageEl = document.getElementById(`page-${page}`);
  if (pageEl) pageEl.classList.add('active');
}

// === Screen Transitions ===
function showAuthScreen() {
  elements.authScreen.classList.remove('hidden');
  elements.mainApp.classList.add('hidden');
}

function showMainApp() {
  elements.authScreen.classList.add('hidden');
  elements.mainApp.classList.remove('hidden');
  updateAlertBadge();
  renderDashboard();
}

// === Render Functions ===
function renderDashboard() {
  const container = document.getElementById('dashboard-content');
  const alertItems = getAlertItems();
  const totalItems = state.items.length;
  const okItems = state.items.filter(i => getItemStatus(i) === 'success').length;

  if (state.items.length === 0) {
    container.innerHTML = renderEmptyState(
      'Start tracking your inventory',
      'Add your first item to get alerts when stock runs low.',
      'Add Your First Item',
      () => navigateTo('item-form')
    );
    return;
  }

  container.innerHTML = `
    <div class="card card-hero hero-card">
      <div class="hero-number">${alertItems.length}</div>
      <div class="hero-label">${alertItems.length === 1 ? 'item needs' : 'items need'} attention</div>
      ${alertItems.length > 0 ? `
        <div class="hero-sub">
          ${alertItems.filter(i => getItemStatus(i) === 'critical').length} Critical ·
          ${alertItems.filter(i => getItemStatus(i) === 'warning').length} Low
        </div>
      ` : '<div class="hero-sub">All stock levels healthy</div>'}
    </div>

    ${alertItems.length > 0 ? `
      <div class="section-header">
        <span class="section-title">Needs Attention</span>
        <a href="#" class="section-link" onclick="navigateTo('alerts'); return false;">View All</a>
      </div>
      ${alertItems.slice(0, 3).map(item => renderItemCard(item)).join('')}
    ` : `
      <div class="empty-state" style="padding: var(--space-8);">
        <div class="empty-state-icon">✓</div>
        <div class="empty-state-title">All stock levels healthy</div>
        <div class="empty-state-description">You have ${totalItems} items tracked. We'll notify you when any fall below threshold.</div>
      </div>
    `}

    <div class="stock-overview">
      <div class="stat-card">
        <div class="stat-number">${totalItems}</div>
        <div class="stat-label">Total Items</div>
      </div>
      <div class="stat-card">
        <div class="stat-number">${alertItems.length}</div>
        <div class="stat-label">Alerts</div>
      </div>
      <div class="stat-card">
        <div class="stat-number">${okItems}</div>
        <div class="stat-label">OK</div>
      </div>
    </div>
  `;
}

function renderItemsList(filter = 'all') {
  const container = document.getElementById('items-content');
  let items = [...state.items];

  // Sort by status urgency
  items.sort((a, b) => {
    const statusOrder = { critical: 0, warning: 1, success: 2 };
    return statusOrder[getItemStatus(a)] - statusOrder[getItemStatus(b)];
  });

  // Apply filter
  if (filter !== 'all') {
    items = items.filter(i => getItemStatus(i) === filter);
  }

  if (state.items.length === 0) {
    container.innerHTML = renderEmptyState(
      'No items yet',
      'Add your first item to start tracking inventory.',
      'Add Item',
      () => navigateTo('item-form')
    );
    return;
  }

  container.innerHTML = `
    <div class="filter-chips">
      <button class="filter-chip ${filter === 'all' ? 'active' : ''}" onclick="filterItems('all')">All</button>
      <button class="filter-chip ${filter === 'critical' ? 'active' : ''}" onclick="filterItems('critical')">Critical</button>
      <button class="filter-chip ${filter === 'warning' ? 'active' : ''}" onclick="filterItems('warning')">Low</button>
      <button class="filter-chip ${filter === 'success' ? 'active' : ''}" onclick="filterItems('success')">OK</button>
    </div>
    <div class="items-list">
      ${items.map(item => renderItemCard(item)).join('')}
    </div>
  `;
}

function filterItems(filter) {
  renderItemsList(filter);
}

function renderItemCard(item) {
  const status = getItemStatus(item);
  const percentage = Math.min(100, (item.quantity / (item.threshold * 2)) * 100);
  const statusLabels = { critical: 'Critical', warning: 'Low', success: 'OK' };
  const statusIcons = { critical: '⚠️', warning: '⚡', success: '✓' };

  return `
    <div class="card card-interactive item-card card-alert-${status}" onclick="navigateTo('item-detail', ${JSON.stringify(item).replace(/"/g, '&quot;')})">
      <div class="item-thumbnail">${item.emoji || '📦'}</div>
      <div class="item-info">
        <div class="item-header">
          <span class="item-name">${escapeHtml(item.name)}</span>
          <span class="badge badge-${status}">${statusIcons[status]} ${statusLabels[status]}</span>
        </div>
        <div class="item-meta">${item.quantity} units · Threshold: ${item.threshold}</div>
        <div class="stock-bar">
          <div class="stock-bar-fill stock-bar-fill-${status}" style="width: ${percentage}%"></div>
          <div class="stock-bar-threshold" style="left: ${Math.min(95, (item.threshold / (item.threshold * 2)) * 100)}%"></div>
        </div>
      </div>
      <div class="item-arrow">
        <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><polyline points="9 18 15 12 9 6"/></svg>
      </div>
    </div>
  `;
}

function renderItemDetail() {
  const container = document.getElementById('item-detail-content');
  const item = state.currentItem;
  const status = getItemStatus(item);
  const percentage = Math.min(100, (item.quantity / (item.threshold * 2)) * 100);
  const statusLabels = { critical: 'Critical', warning: 'Low', success: 'OK' };

  container.innerHTML = `
    <div class="detail-header">
      <button class="btn-ghost" onclick="navigateTo('items')">
        <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><polyline points="15 18 9 12 15 6"/></svg>
        Back
      </button>
      <div class="detail-actions">
        <button class="btn-icon" onclick="navigateTo('item-form', state.currentItem)" aria-label="Edit">
          <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>
        </button>
        <button class="btn-icon" onclick="confirmDeleteItem()" aria-label="Delete">
          <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/><line x1="10" y1="11" x2="10" y2="17"/><line x1="14" y1="11" x2="14" y2="17"/></svg>
        </button>
      </div>
    </div>

    <div class="detail-photo">${item.emoji || '📦'}</div>

    <div class="detail-stock-section">
      <div class="section-title" style="margin-bottom: var(--space-4);">Stock Level</div>
      <div class="stock-bar stock-bar-lg">
        <div class="stock-bar-fill stock-bar-fill-${status}" style="width: ${percentage}%"></div>
        <div class="stock-bar-threshold" style="left: ${Math.min(95, (item.threshold / (item.threshold * 2)) * 100)}%"></div>
      </div>
      <div class="detail-stock-header" style="margin-top: var(--space-4);">
        <div class="detail-stock-value">${item.quantity} units</div>
        <span class="badge badge-${status}">${statusLabels[status]}</span>
      </div>
      <div class="detail-threshold">Threshold: ${item.threshold}</div>
    </div>

    <div class="update-section">
      <div class="update-title">Update Stock</div>
      <div class="quantity-stepper">
        <button class="quantity-btn" onclick="adjustQuantity(-1)">−</button>
        <input type="number" class="quantity-input" id="detail-quantity" value="${item.quantity}" min="0">
        <button class="quantity-btn" onclick="adjustQuantity(1)">+</button>
      </div>
      <div class="update-modes">
        <button class="mode-btn ${state.quickUpdateMode === 'set' ? 'active' : ''}" onclick="setUpdateMode('set')">Set to</button>
        <button class="mode-btn ${state.quickUpdateMode === 'add' ? 'active' : ''}" onclick="setUpdateMode('add')">Add</button>
        <button class="mode-btn ${state.quickUpdateMode === 'remove' ? 'active' : ''}" onclick="setUpdateMode('remove')">Remove</button>
      </div>
      <button class="btn-primary btn-full" onclick="updateItemQuantity()">Update Quantity</button>
    </div>

    ${item.description ? `
      <div class="detail-info">
        <div class="detail-info-row">
          <span class="detail-info-label">Description</span>
          <span class="detail-info-value">${escapeHtml(item.description)}</span>
        </div>
      </div>
    ` : ''}

    ${item.category ? `
      <div class="detail-info">
        <div class="detail-info-row">
          <span class="detail-info-label">Category</span>
          <span class="detail-info-value">${escapeHtml(item.category)}</span>
        </div>
      </div>
    ` : ''}

    <div class="detail-info">
      <div class="detail-info-row">
        <span class="detail-info-label">Last updated</span>
        <span class="detail-info-value">${formatDate(item.updatedAt || item.createdAt)}</span>
      </div>
    </div>

    ${status !== 'success' ? `
      <button class="btn-secondary btn-full" onclick="markAsOrdered()" style="margin-top: var(--space-4);">
        <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><circle cx="9" cy="21" r="1"/><circle cx="20" cy="21" r="1"/><path d="M1 1h4l2.68 13.39a2 2 0 0 0 2 1.61h9.72a2 2 0 0 0 2-1.61L23 6H6"/></svg>
        Mark as Ordered
      </button>
    ` : ''}
  `;
}

function renderItemForm() {
  const container = document.getElementById('item-form-content');
  const item = state.editingItem;
  const isEdit = !!item;

  container.innerHTML = `
    <div class="form-header">
      <button class="btn-ghost" onclick="navigateTo('${isEdit ? 'item-detail' : 'items'}')">
        <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        Cancel
      </button>
      <button class="btn-primary btn-sm" onclick="saveItem()">Save</button>
    </div>

    <form id="item-form" class="form-content">
      <div class="photo-upload ${item?.emoji ? 'has-photo' : ''}" onclick="selectEmoji()">
        ${item?.emoji ? item.emoji : `
          <svg class="icon-lg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"/><circle cx="8.5" cy="8.5" r="1.5"/><polyline points="21 15 16 10 5 21"/></svg>
          <span>Click to add emoji</span>
        `}
      </div>

      <div class="input-group">
        <label class="input-label" for="item-name">Item Name *</label>
        <input class="input" type="text" id="item-name" placeholder="e.g., Coffee Beans" value="${item?.name || ''}" required>
      </div>

      <div class="input-group">
        <label class="input-label" for="item-quantity">Current Quantity *</label>
        <input class="input" type="number" id="item-quantity" placeholder="0" value="${item?.quantity ?? ''}" min="0" required>
      </div>

      <div class="input-group">
        <label class="input-label" for="item-threshold">Reorder Threshold *</label>
        <input class="input" type="number" id="item-threshold" placeholder="e.g., 20" value="${item?.threshold ?? ''}" min="0" required>
        <span class="input-helper">Alert when stock falls below this number</span>
      </div>

      <div class="input-group">
        <label class="input-label" for="item-description">Description (optional)</label>
        <textarea class="input" id="item-description" placeholder="Add notes about this item...">${item?.description || ''}</textarea>
      </div>

      <div class="input-group">
        <label class="input-label" for="item-category">Category (optional)</label>
        <input class="input" type="text" id="item-category" placeholder="e.g., Office Supplies" value="${item?.category || ''}">
      </div>
    </form>

    ${isEdit ? `
      <button class="btn-ghost btn-full" onclick="confirmDeleteItem()" style="margin-top: var(--space-4); color: var(--color-critical);">
        Delete Item
      </button>
    ` : ''}
  `;
}

// Track resolved alerts
let resolvedAlerts = [];

function renderAlerts(tab = 'active') {
  const container = document.getElementById('alerts-content');
  const alertItems = getAlertItems();

  // Group by date
  const today = new Date().toDateString();
  const yesterday = new Date(Date.now() - 86400000).toDateString();

  const grouped = {
    today: [],
    yesterday: [],
    older: []
  };

  alertItems.forEach(item => {
    const itemDate = new Date(item.updatedAt || item.createdAt).toDateString();
    if (itemDate === today) grouped.today.push(item);
    else if (itemDate === yesterday) grouped.yesterday.push(item);
    else grouped.older.push(item);
  });

  // Render resolved tab
  if (tab === 'resolved') {
    container.innerHTML = `
      <div class="alert-tabs">
        <button class="alert-tab" onclick="renderAlerts('active')">Active (${alertItems.length})</button>
        <button class="alert-tab active">Resolved</button>
      </div>
      ${resolvedAlerts.length === 0 ? `
        <div class="empty-state">
          <div class="empty-state-icon">📋</div>
          <div class="empty-state-title">No resolved alerts</div>
          <div class="empty-state-description">Alerts you mark as ordered will appear here.</div>
        </div>
      ` : resolvedAlerts.map(alert => `
        <div class="alert-card" style="opacity: 0.7;">
          <div class="alert-badge-row">
            <span class="badge badge-pending">RESOLVED</span>
          </div>
          <div class="alert-item-name">${escapeHtml(alert.name)}</div>
          <div class="alert-stock-info">Marked as ordered on ${formatDate(alert.resolvedAt)}</div>
        </div>
      `).join('')}
    `;
    return;
  }

  if (alertItems.length === 0) {
    container.innerHTML = `
      <div class="alert-tabs">
        <button class="alert-tab active">Active (0)</button>
        <button class="alert-tab" onclick="renderAlerts('resolved')">Resolved</button>
      </div>
      <div class="empty-state">
        <div class="empty-state-icon">✓</div>
        <div class="empty-state-title">No active alerts</div>
        <div class="empty-state-description">All your items are above their reorder thresholds.</div>
        <button class="btn-secondary" onclick="navigateTo('items')">View All Items</button>
      </div>
    `;
    return;
  }

  container.innerHTML = `
    <div class="alert-tabs">
      <button class="alert-tab active">Active (${alertItems.length})</button>
      <button class="alert-tab" onclick="renderAlerts('resolved')">Resolved</button>
    </div>

    ${grouped.today.length > 0 ? `
      <div class="alert-date">Today</div>
      ${grouped.today.map(item => renderAlertCard(item)).join('')}
    ` : ''}

    ${grouped.yesterday.length > 0 ? `
      <div class="alert-date">Yesterday</div>
      ${grouped.yesterday.map(item => renderAlertCard(item)).join('')}
    ` : ''}

    ${grouped.older.length > 0 ? `
      <div class="alert-date">Earlier</div>
      ${grouped.older.map(item => renderAlertCard(item)).join('')}
    ` : ''}
  `;
}

function renderAlertCard(item) {
  const status = getItemStatus(item);
  const statusLabels = { critical: 'CRITICAL', warning: 'LOW STOCK' };

  return `
    <div class="alert-card">
      <div class="alert-badge-row">
        <span class="badge badge-${status}">
          ${status === 'critical' ? '⚠️' : '⚡'} ${statusLabels[status]}
        </span>
      </div>
      <div class="alert-item-name">${escapeHtml(item.name)}</div>
      <div class="alert-stock-info">${item.quantity} remaining · Threshold: ${item.threshold}</div>
      <div class="alert-actions">
        <button class="alert-action-btn" onclick="event.stopPropagation(); openQuickUpdate(${JSON.stringify(item).replace(/"/g, '&quot;')})">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><polyline points="23 6 13.5 15.5 8.5 10.5 1 18"/><polyline points="17 6 23 6 23 12"/></svg>
          Update
        </button>
        <button class="alert-action-btn" onclick="event.stopPropagation(); snoozeAlert('${item.id}')">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
          Snooze
        </button>
        <button class="alert-action-btn" onclick="event.stopPropagation(); markItemAsOrdered('${item.id}')">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><circle cx="9" cy="21" r="1"/><circle cx="20" cy="21" r="1"/><path d="M1 1h4l2.68 13.39a2 2 0 0 0 2 1.61h9.72a2 2 0 0 0 2-1.61L23 6H6"/></svg>
          Order
        </button>
      </div>
    </div>
  `;
}

// === Alert Snooze ===
function snoozeAlert(itemId) {
  showToast('Alert snoozed for 1 hour', 'success');
}

function renderSettings() {
  const container = document.getElementById('settings-content');

  container.innerHTML = `
    <div class="settings-section">
      <div class="settings-title">Notifications</div>
      <div class="settings-item">
        <div class="settings-item-info">
          <div class="settings-item-title">Email Alerts</div>
          <div class="settings-item-description">Get notified when stock falls below threshold</div>
        </div>
        <div class="toggle ${state.settings.emailAlerts ? 'active' : ''}" onclick="toggleSetting('emailAlerts')">
          <div class="toggle-handle"></div>
        </div>
      </div>
      <div class="settings-item">
        <div class="settings-item-info">
          <div class="settings-item-title">Daily Digest</div>
          <div class="settings-item-description">Receive one summary email at 8:00 AM</div>
        </div>
        <div class="toggle ${state.settings.dailyDigest ? 'active' : ''}" onclick="toggleSetting('dailyDigest')">
          <div class="toggle-handle"></div>
        </div>
      </div>
    </div>

    <div class="settings-section">
      <div class="settings-title">Account</div>
      <div class="settings-item" onclick="showToast('Feature coming soon', 'success')">
        <div class="settings-item-info">
          <div class="settings-item-title">Email</div>
          <div class="settings-item-description">${state.user?.email || 'Not set'}</div>
        </div>
        <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="var(--color-text-tertiary)" stroke-width="1.5"><polyline points="9 18 15 12 9 6"/></svg>
      </div>
      <div class="settings-item" onclick="showToast('Feature coming soon', 'success')">
        <div class="settings-item-info">
          <div class="settings-item-title">Change Password</div>
        </div>
        <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="var(--color-text-tertiary)" stroke-width="1.5"><polyline points="9 18 15 12 9 6"/></svg>
      </div>
    </div>

    <div class="settings-section">
      <div class="settings-title">Data</div>
      <div class="settings-item" onclick="exportData()">
        <div class="settings-item-info">
          <div class="settings-item-title">Export Data</div>
          <div class="settings-item-description">Download your inventory as CSV</div>
        </div>
        <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="var(--color-text-tertiary)" stroke-width="1.5"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
      </div>
      <div class="settings-item" onclick="showToast('Feature coming soon', 'success')">
        <div class="settings-item-info">
          <div class="settings-item-title">Import Items (CSV)</div>
        </div>
        <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="var(--color-text-tertiary)" stroke-width="1.5"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>
      </div>
    </div>

    <button class="btn-secondary btn-full" onclick="handleUserMenu()" style="margin-top: var(--space-4);">Log Out</button>

    <p style="text-align: center; margin-top: var(--space-8); color: var(--color-text-tertiary); font-size: 12px;">
      Reorder Alert v1.0.0
    </p>
  `;
}

function renderEmptyState(title, description, buttonText, onClick) {
  const id = `empty-btn-${Date.now()}`;
  setTimeout(() => {
    const btn = document.getElementById(id);
    if (btn) btn.addEventListener('click', onClick);
  }, 0);

  return `
    <div class="empty-state">
      <div class="empty-state-icon">📦</div>
      <div class="empty-state-title">${title}</div>
      <div class="empty-state-description">${description}</div>
      <button class="btn-primary" id="${id}">${buttonText}</button>
    </div>
  `;
}

// === Item CRUD ===
function saveItem() {
  const name = document.getElementById('item-name').value.trim();
  const quantity = parseInt(document.getElementById('item-quantity').value, 10);
  const threshold = parseInt(document.getElementById('item-threshold').value, 10);
  const description = document.getElementById('item-description').value.trim();
  const category = document.getElementById('item-category').value.trim();

  if (!name || isNaN(quantity) || isNaN(threshold)) {
    showToast('Please fill in all required fields', 'error');
    return;
  }

  if (state.editingItem) {
    // Update existing
    const index = state.items.findIndex(i => i.id === state.editingItem.id);
    if (index !== -1) {
      state.items[index] = {
        ...state.items[index],
        name,
        quantity,
        threshold,
        description,
        category,
        updatedAt: new Date().toISOString()
      };
      state.currentItem = state.items[index];
    }
    showToast('Item updated', 'success');
  } else {
    // Create new
    const newItem = {
      id: generateId(),
      name,
      quantity,
      threshold,
      description,
      category,
      emoji: state.tempEmoji || null,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString()
    };
    state.items.push(newItem);
    state.tempEmoji = null;
    showToast('Item added', 'success');
  }

  saveToStorage();
  updateAlertBadge();
  navigateTo('items');
}

function confirmDeleteItem() {
  const item = state.currentItem || state.editingItem;
  if (!item) return;

  showDialog(
    'Delete Item',
    `Are you sure you want to delete "${item.name}"? This cannot be undone.`,
    () => {
      deleteItem(item.id);
      closeDialog();
    }
  );
}

function deleteItem(id) {
  state.items = state.items.filter(i => i.id !== id);
  saveToStorage();
  updateAlertBadge();
  showToast('Item deleted', 'success');
  navigateTo('items');
}

// === Quantity Updates ===
function adjustQuantity(delta) {
  const input = document.getElementById('detail-quantity');
  const newValue = Math.max(0, parseInt(input.value, 10) + delta);
  input.value = newValue;
}

function setUpdateMode(mode) {
  state.quickUpdateMode = mode;
  document.querySelectorAll('.mode-btn').forEach(btn => {
    btn.classList.toggle('active', btn.textContent.toLowerCase().includes(mode));
  });
}

function updateItemQuantity() {
  const input = document.getElementById('detail-quantity');
  const value = parseInt(input.value, 10);

  if (isNaN(value) || value < 0) {
    showToast('Please enter a valid quantity', 'error');
    return;
  }

  const index = state.items.findIndex(i => i.id === state.currentItem.id);
  if (index === -1) return;

  let newQuantity;
  switch (state.quickUpdateMode) {
    case 'set':
      newQuantity = value;
      break;
    case 'add':
      newQuantity = state.currentItem.quantity + value;
      break;
    case 'remove':
      newQuantity = Math.max(0, state.currentItem.quantity - value);
      break;
    default:
      newQuantity = value;
  }

  state.items[index].quantity = newQuantity;
  state.items[index].updatedAt = new Date().toISOString();
  state.currentItem = state.items[index];

  saveToStorage();
  updateAlertBadge();
  showToast('Quantity updated', 'success');
  renderItemDetail();
}

// === Quick Update Modal ===
function openQuickUpdate(item) {
  state.currentItem = item;
  state.quickUpdateMode = 'add';

  elements.quickUpdateContent.innerHTML = `
    <div class="quick-update-header">
      <div class="quick-update-name">${escapeHtml(item.name)}</div>
      <div class="quick-update-current">Current: ${item.quantity} units</div>
    </div>

    <div class="quantity-stepper">
      <button class="quantity-btn" onclick="adjustQuickQuantity(-1)">−</button>
      <input type="number" class="quantity-input" id="quick-quantity" value="0" min="0">
      <button class="quantity-btn" onclick="adjustQuickQuantity(1)">+</button>
    </div>

    <div class="update-modes">
      <button class="mode-btn" onclick="setQuickMode('set')">Set to</button>
      <button class="mode-btn active" onclick="setQuickMode('add')">Add</button>
      <button class="mode-btn" onclick="setQuickMode('remove')">Remove</button>
    </div>

    <div class="quick-presets">
      <button class="preset-btn" onclick="setQuickPreset(10)">+10</button>
      <button class="preset-btn" onclick="setQuickPreset(25)">+25</button>
      <button class="preset-btn" onclick="setQuickPreset(50)">+50</button>
      <button class="preset-btn" onclick="setQuickPreset(100)">+100</button>
    </div>

    <button class="btn-primary btn-full" onclick="submitQuickUpdate()" style="margin-top: var(--space-6);">Update</button>
  `;

  elements.quickUpdateOverlay.classList.add('open');
  elements.quickUpdateSheet.classList.add('open');
}

function closeQuickUpdate() {
  elements.quickUpdateOverlay.classList.remove('open');
  elements.quickUpdateSheet.classList.remove('open');
}

function adjustQuickQuantity(delta) {
  const input = document.getElementById('quick-quantity');
  const newValue = Math.max(0, parseInt(input.value, 10) + delta);
  input.value = newValue;
}

function setQuickMode(mode) {
  state.quickUpdateMode = mode;
  document.querySelectorAll('.bottom-sheet .mode-btn').forEach(btn => {
    btn.classList.toggle('active', btn.textContent.toLowerCase().includes(mode));
  });
}

function setQuickPreset(value) {
  document.getElementById('quick-quantity').value = value;
}

function submitQuickUpdate() {
  const input = document.getElementById('quick-quantity');
  const value = parseInt(input.value, 10);

  if (isNaN(value) || value < 0) {
    showToast('Please enter a valid quantity', 'error');
    return;
  }

  const index = state.items.findIndex(i => i.id === state.currentItem.id);
  if (index === -1) return;

  let newQuantity;
  switch (state.quickUpdateMode) {
    case 'set':
      newQuantity = value;
      break;
    case 'add':
      newQuantity = state.items[index].quantity + value;
      break;
    case 'remove':
      newQuantity = Math.max(0, state.items[index].quantity - value);
      break;
    default:
      newQuantity = state.items[index].quantity + value;
  }

  state.items[index].quantity = newQuantity;
  state.items[index].updatedAt = new Date().toISOString();

  saveToStorage();
  updateAlertBadge();
  closeQuickUpdate();
  showToast('Quantity updated', 'success');

  // Re-render current page
  if (state.currentPage === 'dashboard') renderDashboard();
  else if (state.currentPage === 'alerts') renderAlerts();
  else if (state.currentPage === 'items') renderItemsList();
}

// === Mark as Ordered ===
function markAsOrdered() {
  if (state.currentItem) {
    resolvedAlerts.push({
      ...state.currentItem,
      resolvedAt: new Date().toISOString()
    });
    showToast('Marked as ordered', 'success');
    renderItemDetail();
  }
}

function markItemAsOrdered(itemId) {
  const item = findItemById(itemId);
  if (item) {
    resolvedAlerts.push({
      ...item,
      resolvedAt: new Date().toISOString()
    });
    showToast('Marked as ordered', 'success');
    if (state.currentPage === 'alerts') renderAlerts();
    if (state.currentPage === 'dashboard') renderDashboard();
  }
}

// === Emoji Picker ===
function selectEmoji() {
  const emojis = ['📦', '☕', '📄', '🖨️', '✏️', '📎', '🧴', '🧻', '💊', '🔧', '⚙️', '🔩', '📱', '💻', '🖥️', '⌨️', '🖱️', '🎧', '📷', '💡'];
  const selected = emojis[Math.floor(Math.random() * emojis.length)];

  if (state.editingItem) {
    state.editingItem.emoji = selected;
  } else {
    state.tempEmoji = selected;
  }

  const photoUpload = document.querySelector('.photo-upload');
  photoUpload.innerHTML = selected;
  photoUpload.classList.add('has-photo');
}

// === Helpers ===
function getItemStatus(item) {
  if (item.quantity <= item.threshold * 0.5) return 'critical';
  if (item.quantity <= item.threshold) return 'warning';
  return 'success';
}

function getAlertItems() {
  return state.items.filter(item => getItemStatus(item) !== 'success')
    .sort((a, b) => (a.quantity / a.threshold) - (b.quantity / b.threshold));
}

function updateAlertBadge() {
  const count = getAlertItems().length;
  elements.alertBadge.textContent = count;
  elements.alertBadge.classList.toggle('hidden', count === 0);
}

function generateId() {
  return Date.now().toString(36) + Math.random().toString(36).substr(2);
}

function escapeHtml(str) {
  const div = document.createElement('div');
  div.textContent = str;
  return div.innerHTML;
}

function formatDate(dateStr) {
  const date = new Date(dateStr);
  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric'
  });
}

function findItemById(id) {
  return state.items.find(i => i.id === id);
}

// === Toast ===
function showToast(message, type = 'success') {
  elements.toastMessage.textContent = message;
  elements.toast.className = `toast toast-${type} visible`;

  if (type === 'success') {
    elements.toastIcon.innerHTML = '<polyline points="20 6 9 17 4 12"/>';
  } else {
    elements.toastIcon.innerHTML = '<line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>';
  }

  setTimeout(() => {
    elements.toast.classList.remove('visible');
  }, 3000);
}

// === Dialog ===
function showDialog(title, message, onConfirm, confirmText = 'Delete') {
  elements.dialogTitle.textContent = title;
  elements.dialogMessage.textContent = message;
  elements.dialogConfirm.textContent = confirmText;
  elements.dialogConfirm.onclick = onConfirm;
  elements.confirmDialog.classList.remove('hidden');
}

function closeDialog() {
  elements.confirmDialog.classList.add('hidden');
}

// === Settings ===
function toggleSetting(key) {
  state.settings[key] = !state.settings[key];
  saveToStorage();
  renderSettings();
}

// === Export ===
function exportData() {
  const csv = [
    ['Name', 'Quantity', 'Threshold', 'Description', 'Category', 'Created', 'Updated'].join(','),
    ...state.items.map(item => [
      `"${item.name}"`,
      item.quantity,
      item.threshold,
      `"${item.description || ''}"`,
      `"${item.category || ''}"`,
      item.createdAt,
      item.updatedAt
    ].join(','))
  ].join('\n');

  const blob = new Blob([csv], { type: 'text/csv' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = `inventory-${new Date().toISOString().split('T')[0]}.csv`;
  a.click();
  URL.revokeObjectURL(url);

  showToast('Data exported', 'success');
}

// Make functions available globally for onclick handlers
window.navigateTo = navigateTo;
window.filterItems = filterItems;
window.adjustQuantity = adjustQuantity;
window.setUpdateMode = setUpdateMode;
window.updateItemQuantity = updateItemQuantity;
window.confirmDeleteItem = confirmDeleteItem;
window.selectEmoji = selectEmoji;
window.saveItem = saveItem;
window.openQuickUpdate = openQuickUpdate;
window.adjustQuickQuantity = adjustQuickQuantity;
window.setQuickMode = setQuickMode;
window.setQuickPreset = setQuickPreset;
window.submitQuickUpdate = submitQuickUpdate;
window.markAsOrdered = markAsOrdered;
window.markItemAsOrdered = markItemAsOrdered;
window.toggleSetting = toggleSetting;
window.showToast = showToast;
window.handleUserMenu = handleUserMenu;
window.exportData = exportData;
window.findItemById = findItemById;
window.renderAlerts = renderAlerts;
window.snoozeAlert = snoozeAlert;
window.state = state;

// Initialize
document.addEventListener('DOMContentLoaded', init);
