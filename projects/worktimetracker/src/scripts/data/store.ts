import { eventBus } from '../utils/eventBus';

const STORAGE_PREFIX = 'wtt_';

class Store {
  private cache: Map<string, unknown> = new Map();
  private initialized = false;

  init(): void {
    if (this.initialized) return;

    for (let i = 0; i < localStorage.length; i++) {
      const key = localStorage.key(i);
      if (key && key.startsWith(STORAGE_PREFIX)) {
        try {
          const raw = localStorage.getItem(key);
          if (raw !== null) {
            this.cache.set(key, JSON.parse(raw));
          }
        } catch (e) {
          console.warn(`Corrupted data in localStorage key "${key}":`, e);
        }
      }
    }

    this.initialized = true;
  }

  get<T>(key: string): T[] {
    const fullKey = STORAGE_PREFIX + key;

    if (this.cache.has(fullKey)) {
      return this.cache.get(fullKey) as T[];
    }

    try {
      const raw = localStorage.getItem(fullKey);
      if (raw === null) return [];
      const parsed = JSON.parse(raw) as T[];
      this.cache.set(fullKey, parsed);
      return parsed;
    } catch (e) {
      console.warn(`Failed to read "${fullKey}":`, e);
      return [];
    }
  }

  set<T>(key: string, data: T): void {
    const fullKey = STORAGE_PREFIX + key;

    try {
      const json = JSON.stringify(data);
      localStorage.setItem(fullKey, json);
      this.cache.set(fullKey, data);
    } catch (e) {
      if (e instanceof DOMException && e.name === 'QuotaExceededError') {
        eventBus.emit('storage:quota-exceeded');
        console.error('localStorage quota exceeded');
      }
      throw e;
    }
  }

  remove(key: string): void {
    const fullKey = STORAGE_PREFIX + key;
    try {
      localStorage.removeItem(fullKey);
      this.cache.delete(fullKey);
    } catch (e) {
      console.warn(`Failed to remove "${fullKey}":`, e);
    }
  }

  clearAll(): void {
    const keysToRemove: string[] = [];
    for (let i = 0; i < localStorage.length; i++) {
      const key = localStorage.key(i);
      if (key && key.startsWith(STORAGE_PREFIX)) {
        keysToRemove.push(key);
      }
    }
    keysToRemove.forEach(key => localStorage.removeItem(key));
    this.cache.clear();
  }

  getStorageUsage(): { used: number; total: number; percentage: number } {
    let charCount = 0;
    for (let i = 0; i < localStorage.length; i++) {
      const key = localStorage.key(i);
      if (key && key.startsWith(STORAGE_PREFIX)) {
        const value = localStorage.getItem(key);
        if (value) {
          charCount += key.length + value.length;
        }
      }
    }
    // UTF-16: ~2 bytes per character
    const usedBytes = charCount * 2;
    const totalBytes = 5 * 1024 * 1024;
    return {
      used: usedBytes,
      total: totalBytes,
      percentage: Math.round((usedBytes / totalBytes) * 100),
    };
  }

  exportAll(): Record<string, unknown> {
    const data: Record<string, unknown> = {};
    for (let i = 0; i < localStorage.length; i++) {
      const key = localStorage.key(i);
      if (key && key.startsWith(STORAGE_PREFIX)) {
        try {
          const raw = localStorage.getItem(key);
          if (raw !== null) {
            data[key.slice(STORAGE_PREFIX.length)] = JSON.parse(raw);
          }
        } catch (e) {
          console.warn(`Skipping corrupted key "${key}" during export:`, e);
        }
      }
    }
    return data;
  }

  importAll(data: Record<string, unknown>): void {
    for (const [key, value] of Object.entries(data)) {
      this.set(key, value);
    }
  }
}

export const store = new Store();
