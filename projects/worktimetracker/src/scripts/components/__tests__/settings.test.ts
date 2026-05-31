/**
 * Settings Panel Tests
 */
import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import { openSettings, initSettings } from '../settings';
import { store } from '../../data/store';
import { projectService } from '../../data/projectService';
import { Modal } from '../modal';
import { Toast } from '../toast';
import { eventBus } from '../../utils/eventBus';

describe('Settings Panel', () => {
  beforeEach(() => {
    // Set up DOM
    document.body.innerHTML = `
      <div id="toast-container"></div>
      <div id="modal-container"></div>
      <button id="settings-btn"></button>
      <button id="mobile-settings-btn"></button>
    `;
    Modal._reset();
    Toast.clearAll();
    eventBus._reset();
    store.init();
    projectService.ensureGeneralProject();
  });

  afterEach(() => {
    Modal._reset();
    Toast.clearAll();
    eventBus._reset();
    store.clearAll();
    document.body.innerHTML = '';
  });

  describe('openSettings', () => {
    it('should open a modal with title "Settings"', () => {
      openSettings();
      expect(Modal.isOpen()).toBe(true);
      const title = document.querySelector('#modal-title');
      expect(title?.textContent).toBe('Settings');
    });

    it('should contain General section with display preferences', () => {
      openSettings();
      const modal = document.querySelector('.modal-container');
      expect(modal?.textContent).toContain('General');
      expect(modal?.textContent).toContain('Time Format');
      expect(modal?.textContent).toContain('Week Starts On');
      expect(modal?.textContent).toContain('Default View');
    });

    it('should contain Data Management section', () => {
      openSettings();
      const modal = document.querySelector('.modal-container');
      expect(modal?.textContent).toContain('Data Management');
      expect(modal?.textContent).toContain('Export All Data');
      expect(modal?.textContent).toContain('Import Data');
      expect(modal?.textContent).toContain('Clear All Data');
    });

    it('should show storage usage indicator', () => {
      openSettings();
      const modal = document.querySelector('.modal-container');
      expect(modal?.textContent).toContain('Storage');
      expect(modal?.textContent).toContain('used');
    });
  });

  describe('initSettings', () => {
    it('should wire up settings button click', () => {
      initSettings();
      const btn = document.getElementById('settings-btn')!;
      btn.click();
      expect(Modal.isOpen()).toBe(true);
    });

    it('should wire up mobile settings button click', () => {
      initSettings();
      const btn = document.getElementById('mobile-settings-btn')!;
      btn.click();
      expect(Modal.isOpen()).toBe(true);
    });
  });

  describe('Store integration', () => {
    it('store.exportAll returns data with all wtt_ keys', () => {
      store.set('test_collection', [{ id: '1', name: 'test' }]);
      const exported = store.exportAll();
      expect(exported).toHaveProperty('test_collection');
      expect(exported.test_collection).toEqual([{ id: '1', name: 'test' }]);
    });

    it('store.clearAll removes all wtt_ keys', () => {
      store.set('test_collection', [{ id: '1' }]);
      expect(store.get('test_collection')).toHaveLength(1);
      store.clearAll();
      expect(store.get('test_collection')).toHaveLength(0);
    });

    it('store.getStorageUsage returns valid usage data', () => {
      const usage = store.getStorageUsage();
      expect(usage).toHaveProperty('used');
      expect(usage).toHaveProperty('total');
      expect(usage).toHaveProperty('percentage');
      expect(usage.total).toBe(5 * 1024 * 1024);
      expect(usage.percentage).toBeGreaterThanOrEqual(0);
      expect(usage.percentage).toBeLessThanOrEqual(100);
    });

    it('store.importAll merges data correctly', () => {
      store.set('employees', [{ id: '1', name: 'Alice' }]);
      store.importAll({ employees: [{ id: '1', name: 'Alice' }, { id: '2', name: 'Bob' }] });
      // importAll replaces the key entirely
      const employees = store.get('employees');
      expect(employees).toHaveLength(2);
    });
  });
});
