/**
 * Settings Panel Component
 *
 * Accessible from the gear icon in the sidebar. Opens as a modal panel.
 * - General preferences: time format (12h/24h), week starts on, default view
 * - Data Management: Export All Data (JSON), Import Data (JSON), Clear All Data
 * - Storage usage indicator bar
 * - Toast notifications for each action
 */

import { store } from '../data/store';
import { projectService } from '../data/projectService';
import { eventBus } from '../utils/eventBus';
import { Modal } from './modal';
import { Toast } from './toast';
import type { AppSettings } from '../types';

const SETTINGS_KEY = 'settings';

const VALID_COLLECTIONS = [
  'employees',
  'projects',
  'tasks',
  'time_entries',
  'active_timers',
  'settings',
];

function getSettings(): AppSettings {
  const all = store.get<AppSettings>(SETTINGS_KEY);
  if (all.length > 0) return all[0] as unknown as AppSettings;
  return { timeFormat: '24h', weekStartsOn: 'monday', defaultView: 'dashboard' };
}

function saveSettings(settings: AppSettings): void {
  store.set(SETTINGS_KEY, [settings]);
  eventBus.emit('settings:changed', settings);
}

function formatBytes(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(2)} MB`;
}

function buildStorageIndicator(): HTMLElement {
  const usage = store.getStorageUsage();
  const section = document.createElement('div');
  section.style.marginTop = '16px';

  const label = document.createElement('div');
  label.style.cssText = 'display: flex; justify-content: space-between; align-items: center; margin-bottom: 6px;';

  const text = document.createElement('span');
  text.style.cssText = 'font-size: 13px; color: #8C83A8;';
  text.textContent = `Storage: ${formatBytes(usage.used)} of ~${formatBytes(usage.total)} used`;

  const pct = document.createElement('span');
  pct.style.cssText = 'font-size: 13px; color: #8C83A8;';
  pct.textContent = `${usage.percentage}%`;

  label.appendChild(text);
  label.appendChild(pct);
  section.appendChild(label);

  const barBg = document.createElement('div');
  barBg.style.cssText = 'height: 8px; border-radius: 9999px; background: rgba(91, 70, 178, 0.15); overflow: hidden;';

  const barFill = document.createElement('div');
  const color = usage.percentage >= 80 ? '#D94B4B' : usage.percentage >= 60 ? '#DFA23E' : '#5746B2';
  barFill.style.cssText = `height: 100%; border-radius: 9999px; background: ${color}; transition: width 300ms ease; width: ${Math.min(usage.percentage, 100)}%;`;

  barBg.appendChild(barFill);
  section.appendChild(barBg);

  return section;
}

function buildGeneralSection(settings: AppSettings): HTMLElement {
  const section = document.createElement('div');

  const heading = document.createElement('h3');
  heading.style.cssText = 'font-size: 16px; font-weight: 600; color: #ECE9F5; margin-bottom: 16px; font-family: Manrope, sans-serif;';
  heading.textContent = 'General';
  section.appendChild(heading);

  const grid = document.createElement('div');
  grid.style.cssText = 'display: flex; flex-direction: column; gap: 16px;';

  // Time Format
  const timeRow = document.createElement('div');
  timeRow.style.cssText = 'display: flex; align-items: center; justify-content: space-between;';
  const timeLabel = document.createElement('span');
  timeLabel.style.cssText = 'font-size: 14px; color: #ECE9F5;';
  timeLabel.textContent = 'Time Format';

  const timeToggle = document.createElement('div');
  timeToggle.style.cssText = 'display: flex; gap: 4px; background: rgba(91, 70, 178, 0.10); border-radius: 8px; padding: 2px;';

  const btn12 = document.createElement('button');
  btn12.type = 'button';
  btn12.textContent = '12-hour';
  btn12.style.cssText = `padding: 6px 12px; border-radius: 6px; font-size: 13px; border: none; cursor: pointer; transition: all 150ms ease; ${
    settings.timeFormat === '12h'
      ? 'background: #5746B2; color: #ECE9F5;'
      : 'background: transparent; color: #8C83A8;'
  }`;

  const btn24 = document.createElement('button');
  btn24.type = 'button';
  btn24.textContent = '24-hour';
  btn24.style.cssText = `padding: 6px 12px; border-radius: 6px; font-size: 13px; border: none; cursor: pointer; transition: all 150ms ease; ${
    settings.timeFormat === '24h'
      ? 'background: #5746B2; color: #ECE9F5;'
      : 'background: transparent; color: #8C83A8;'
  }`;

  btn12.addEventListener('click', () => {
    settings.timeFormat = '12h';
    saveSettings(settings);
    btn12.style.background = '#5746B2';
    btn12.style.color = '#ECE9F5';
    btn24.style.background = 'transparent';
    btn24.style.color = '#8C83A8';
    Toast.success('Time format set to 12-hour');
  });

  btn24.addEventListener('click', () => {
    settings.timeFormat = '24h';
    saveSettings(settings);
    btn24.style.background = '#5746B2';
    btn24.style.color = '#ECE9F5';
    btn12.style.background = 'transparent';
    btn12.style.color = '#8C83A8';
    Toast.success('Time format set to 24-hour');
  });

  timeToggle.appendChild(btn12);
  timeToggle.appendChild(btn24);
  timeRow.appendChild(timeLabel);
  timeRow.appendChild(timeToggle);
  grid.appendChild(timeRow);

  // Week starts on
  const weekRow = document.createElement('div');
  weekRow.style.cssText = 'display: flex; align-items: center; justify-content: space-between;';
  const weekLabel = document.createElement('span');
  weekLabel.style.cssText = 'font-size: 14px; color: #ECE9F5;';
  weekLabel.textContent = 'Week Starts On';

  const weekSelect = document.createElement('select');
  weekSelect.className = 'select';
  weekSelect.style.cssText = 'width: auto; min-width: 130px;';

  const optMon = document.createElement('option');
  optMon.value = 'monday';
  optMon.textContent = 'Monday';
  optMon.selected = settings.weekStartsOn === 'monday';

  const optSun = document.createElement('option');
  optSun.value = 'sunday';
  optSun.textContent = 'Sunday';
  optSun.selected = settings.weekStartsOn === 'sunday';

  weekSelect.appendChild(optMon);
  weekSelect.appendChild(optSun);
  weekSelect.addEventListener('change', () => {
    settings.weekStartsOn = weekSelect.value as 'monday' | 'sunday';
    saveSettings(settings);
    Toast.success(`Week starts on ${weekSelect.value === 'monday' ? 'Monday' : 'Sunday'}`);
  });

  weekRow.appendChild(weekLabel);
  weekRow.appendChild(weekSelect);
  grid.appendChild(weekRow);

  // Default view
  const viewRow = document.createElement('div');
  viewRow.style.cssText = 'display: flex; align-items: center; justify-content: space-between;';
  const viewLabel = document.createElement('span');
  viewLabel.style.cssText = 'font-size: 14px; color: #ECE9F5;';
  viewLabel.textContent = 'Default View';

  const viewSelect = document.createElement('select');
  viewSelect.className = 'select';
  viewSelect.style.cssText = 'width: auto; min-width: 130px;';

  const views = [
    { value: 'dashboard', label: 'Dashboard' },
    { value: 'employees', label: 'Employees' },
    { value: 'timesheet', label: 'Time Tracker' },
    { value: 'projects', label: 'Projects' },
    { value: 'reports', label: 'Reports' },
  ];

  for (const v of views) {
    const opt = document.createElement('option');
    opt.value = v.value;
    opt.textContent = v.label;
    opt.selected = settings.defaultView === v.value;
    viewSelect.appendChild(opt);
  }

  viewSelect.addEventListener('change', () => {
    settings.defaultView = viewSelect.value;
    saveSettings(settings);
    Toast.success(`Default view set to ${viewSelect.options[viewSelect.selectedIndex].textContent}`);
  });

  viewRow.appendChild(viewLabel);
  viewRow.appendChild(viewSelect);
  grid.appendChild(viewRow);

  section.appendChild(grid);
  return section;
}

function buildDataSection(): HTMLElement {
  const section = document.createElement('div');
  section.style.marginTop = '24px';

  const heading = document.createElement('h3');
  heading.style.cssText = 'font-size: 16px; font-weight: 600; color: #ECE9F5; margin-bottom: 16px; font-family: Manrope, sans-serif;';
  heading.textContent = 'Data Management';
  section.appendChild(heading);

  const btnGroup = document.createElement('div');
  btnGroup.style.cssText = 'display: flex; flex-direction: column; gap: 10px;';

  // Export All Data
  const exportBtn = document.createElement('button');
  exportBtn.type = 'button';
  exportBtn.className = 'btn-secondary';
  exportBtn.style.justifyContent = 'flex-start';
  exportBtn.textContent = 'Export All Data (JSON)';
  exportBtn.addEventListener('click', handleExportData);
  btnGroup.appendChild(exportBtn);

  // Import Data
  const importBtn = document.createElement('button');
  importBtn.type = 'button';
  importBtn.className = 'btn-secondary';
  importBtn.style.justifyContent = 'flex-start';
  importBtn.textContent = 'Import Data (JSON)';
  importBtn.addEventListener('click', handleImportData);
  btnGroup.appendChild(importBtn);

  // Clear All Data
  const clearBtn = document.createElement('button');
  clearBtn.type = 'button';
  clearBtn.style.cssText = 'display: inline-flex; align-items: center; justify-content: flex-start; padding: 10px 16px; border-radius: 8px; font-size: 14px; font-weight: 500; border: 1px solid rgba(217, 75, 75, 0.30); color: #D94B4B; background: transparent; cursor: pointer; transition: all 200ms ease; min-height: 40px;';
  clearBtn.textContent = 'Clear All Data';
  clearBtn.addEventListener('mouseenter', () => {
    clearBtn.style.background = 'rgba(217, 75, 75, 0.10)';
    clearBtn.style.borderColor = 'rgba(217, 75, 75, 0.50)';
  });
  clearBtn.addEventListener('mouseleave', () => {
    clearBtn.style.background = 'transparent';
    clearBtn.style.borderColor = 'rgba(217, 75, 75, 0.30)';
  });
  clearBtn.addEventListener('click', handleClearData);
  btnGroup.appendChild(clearBtn);

  section.appendChild(btnGroup);
  section.appendChild(buildStorageIndicator());

  return section;
}

function handleExportData(): void {
  try {
    const data = store.exportAll();
    const json = JSON.stringify(data, null, 2);
    const blob = new Blob([json], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    const dateStr = new Date().toISOString().split('T')[0];
    link.download = `worktime-backup-${dateStr}.json`;
    link.style.display = 'none';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
    Toast.success('Data exported successfully');
  } catch {
    Toast.error('Failed to export data');
  }
}

function handleImportData(): void {
  // Close settings modal first to avoid stacking
  Modal.close();

  // Small delay to let modal close animation complete
  setTimeout(() => {
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = '.json,application/json';
    input.style.display = 'none';

    input.addEventListener('change', () => {
      const file = input.files?.[0];
      if (!file) return;

      const reader = new FileReader();
      reader.onload = () => {
        try {
          const data = JSON.parse(reader.result as string);
          if (typeof data !== 'object' || data === null || Array.isArray(data)) {
            Toast.error('Invalid backup file: expected a JSON object');
            return;
          }

          const keys = Object.keys(data);
          const unknownKeys = keys.filter(k => !VALID_COLLECTIONS.includes(k));
          if (unknownKeys.length > 0) {
            Toast.warning(`Skipping unknown keys: ${unknownKeys.join(', ')}`);
          }

          for (const key of keys) {
            if (VALID_COLLECTIONS.includes(key) && !Array.isArray(data[key])) {
              Toast.error(`Invalid data: "${key}" must be an array`);
              return;
            }
          }

          showImportDialog(data, keys.filter(k => VALID_COLLECTIONS.includes(k)));
        } catch {
          Toast.error('Invalid JSON file');
        }
      };
      reader.readAsText(file);
    });

    document.body.appendChild(input);
    input.click();
    document.body.removeChild(input);
  }, 200);
}

function showImportDialog(data: Record<string, unknown>, validKeys: string[]): void {
  let totalItems = 0;
  for (const key of validKeys) {
    const arr = data[key];
    if (Array.isArray(arr)) totalItems += arr.length;
  }

  const content = document.createElement('div');

  const info = document.createElement('p');
  info.style.cssText = 'font-size: 14px; color: #8C83A8; margin-bottom: 16px;';
  info.textContent = `Found ${totalItems} items across ${validKeys.length} collections: ${validKeys.join(', ')}`;
  content.appendChild(info);

  const warning = document.createElement('p');
  warning.style.cssText = 'font-size: 13px; color: #DFA23E; margin-bottom: 16px;';
  warning.textContent = 'Choose how to import:';
  content.appendChild(warning);

  const btnRow = document.createElement('div');
  btnRow.style.cssText = 'display: flex; gap: 8px; flex-wrap: wrap;';

  const mergeBtn = document.createElement('button');
  mergeBtn.type = 'button';
  mergeBtn.className = 'btn-primary';
  mergeBtn.textContent = 'Merge with existing';
  mergeBtn.addEventListener('click', () => {
    for (const key of validKeys) {
      const importArr = data[key] as Array<{ id?: string }>;
      if (!Array.isArray(importArr)) continue;
      const existing = store.get<{ id?: string }>(key);
      const existingIds = new Set(existing.map(item => item.id));
      const merged = [...existing];
      for (const item of importArr) {
        if (item.id && !existingIds.has(item.id)) {
          merged.push(item);
        }
      }
      store.set(key, merged);
    }
    Modal.close();
    emitAllChangedEvents();
    Toast.success(`Merged ${totalItems} items from backup`);
  });

  const replaceBtn = document.createElement('button');
  replaceBtn.type = 'button';
  replaceBtn.style.cssText = 'display: inline-flex; align-items: center; justify-content: center; padding: 10px 16px; border-radius: 8px; font-size: 14px; font-weight: 500; border: 1px solid rgba(217, 75, 75, 0.30); color: #D94B4B; background: transparent; cursor: pointer; transition: all 200ms ease; min-height: 40px;';
  replaceBtn.textContent = 'Replace all data';
  replaceBtn.addEventListener('click', () => {
    store.clearAll();
    const filtered: Record<string, unknown> = {};
    for (const key of validKeys) {
      filtered[key] = data[key];
    }
    store.importAll(filtered);
    projectService.ensureGeneralProject();
    Modal.close();
    emitAllChangedEvents();
    Toast.success(`Replaced all data with ${totalItems} imported items`);
  });

  btnRow.appendChild(mergeBtn);
  btnRow.appendChild(replaceBtn);
  content.appendChild(btnRow);

  Modal.open({
    title: 'Import Data',
    content,
    hideActions: true,
  });
}

function handleClearData(): void {
  // Close settings modal first
  Modal.close();

  setTimeout(() => {
    const content = document.createElement('div');

    const msg = document.createElement('p');
    msg.style.cssText = 'font-size: 14px; color: #ECE9F5; margin-bottom: 8px;';
    msg.textContent = 'This will delete all employees, projects, time entries, and settings.';
    content.appendChild(msg);

    const warn = document.createElement('p');
    warn.style.cssText = 'font-size: 14px; font-weight: 600; color: #D94B4B;';
    warn.textContent = 'This cannot be undone.';
    content.appendChild(warn);

    Modal.open({
      title: 'Clear All Data',
      content,
      submitLabel: 'Delete Everything',
      cancelLabel: 'Cancel',
      onSubmit: () => {
        store.clearAll();
        projectService.ensureGeneralProject();
        Modal.close();
        emitAllChangedEvents();
        Toast.success('All data cleared');
      },
    });
  }, 200);
}

function emitAllChangedEvents(): void {
  eventBus.emit('employees:changed');
  eventBus.emit('projects:changed');
  eventBus.emit('tasks:changed');
  eventBus.emit('timeEntries:changed');
  eventBus.emit('timer:changed');
  eventBus.emit('settings:changed');
}

/** Open the settings panel as a modal. */
export function openSettings(): void {
  const settings = getSettings();

  const content = document.createElement('div');
  content.appendChild(buildGeneralSection(settings));
  content.appendChild(buildDataSection());

  Modal.open({
    title: 'Settings',
    content,
    hideActions: true,
  });
}

/** Initialize settings button in the sidebar. */
export function initSettings(): void {
  const settingsBtn = document.getElementById('settings-btn');
  if (settingsBtn) {
    settingsBtn.addEventListener('click', openSettings);
  }

  const mobileSettingsBtn = document.getElementById('mobile-settings-btn');
  if (mobileSettingsBtn) {
    mobileSettingsBtn.addEventListener('click', openSettings);
  }
}
