import { describe, it, expect } from 'vitest';
import { formatDuration, calcDuration, toUTC, toLocal, todayDate } from '../dates';

describe('dates', () => {
  describe('formatDuration', () => {
    it('should format minutes into Xh Ym', () => {
      expect(formatDuration(90)).toBe('1h 30m');
      expect(formatDuration(60)).toBe('1h 0m');
      expect(formatDuration(45)).toBe('0h 45m');
      expect(formatDuration(0)).toBe('0h 0m');
    });

    it('should handle negative minutes', () => {
      expect(formatDuration(-10)).toBe('0h 0m');
    });
  });

  describe('calcDuration', () => {
    it('should calculate duration in minutes between two dates', () => {
      const start = '2026-03-18T09:00:00.000Z';
      const end = '2026-03-18T10:30:00.000Z';
      expect(calcDuration(start, end)).toBe(90);
    });

    it('should work with Date objects', () => {
      const start = new Date('2026-03-18T09:00:00.000Z');
      const end = new Date('2026-03-18T09:45:00.000Z');
      expect(calcDuration(start, end)).toBe(45);
    });
  });

  describe('toUTC', () => {
    it('should return ISO string', () => {
      const date = new Date('2026-03-18T09:00:00.000Z');
      expect(toUTC(date)).toBe('2026-03-18T09:00:00.000Z');
    });
  });

  describe('toLocal', () => {
    it('should parse ISO string to Date', () => {
      const result = toLocal('2026-03-18T09:00:00.000Z');
      expect(result instanceof Date).toBe(true);
      expect(result.toISOString()).toBe('2026-03-18T09:00:00.000Z');
    });
  });

  describe('todayDate', () => {
    it('should return YYYY-MM-DD format', () => {
      const result = todayDate();
      expect(result).toMatch(/^\d{4}-\d{2}-\d{2}$/);
    });
  });
});
