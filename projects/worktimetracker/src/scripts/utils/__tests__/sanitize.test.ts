import { describe, it, expect } from 'vitest';
import {
  escapeHTML,
  validateRequired,
  validateMaxLength,
  validateEmail,
  validateTimeRange,
  sanitizeString,
} from '../sanitize';

describe('sanitize', () => {
  describe('escapeHTML', () => {
    it('should escape HTML tags', () => {
      expect(escapeHTML('<script>alert("xss")</script>')).not.toContain('<script>');
    });

    it('should escape ampersands and angle brackets', () => {
      const result = escapeHTML('A & B <C>');
      expect(result).toContain('&amp;');
      expect(result).toContain('&lt;');
      expect(result).toContain('&gt;');
    });

    it('should return empty string for empty input', () => {
      expect(escapeHTML('')).toBe('');
    });
  });

  describe('validateRequired', () => {
    it('should return error for empty string', () => {
      expect(validateRequired('', 'Name')).toBe('Name is required');
    });

    it('should return error for whitespace-only string', () => {
      expect(validateRequired('   ', 'Name')).toBe('Name is required');
    });

    it('should return null for valid input', () => {
      expect(validateRequired('John', 'Name')).toBeNull();
    });
  });

  describe('validateMaxLength', () => {
    it('should return error when exceeding max', () => {
      expect(validateMaxLength('a'.repeat(101), 100, 'Name')).toBeTruthy();
    });

    it('should return null when within limit', () => {
      expect(validateMaxLength('short', 100, 'Name')).toBeNull();
    });
  });

  describe('validateEmail', () => {
    it('should return null for valid email', () => {
      expect(validateEmail('test@example.com')).toBeNull();
    });

    it('should return error for invalid email', () => {
      expect(validateEmail('not-an-email')).toBeTruthy();
    });

    it('should return null for empty string', () => {
      expect(validateEmail('')).toBeNull();
    });
  });

  describe('validateTimeRange', () => {
    it('should return error when clockOut <= clockIn', () => {
      expect(
        validateTimeRange('2026-03-18T10:00:00Z', '2026-03-18T09:00:00Z')
      ).toBeTruthy();
    });

    it('should return error when exceeding 24 hours', () => {
      expect(
        validateTimeRange('2026-03-18T09:00:00Z', '2026-03-20T09:00:00Z')
      ).toBeTruthy();
    });

    it('should return null for valid range', () => {
      expect(
        validateTimeRange('2026-03-18T09:00:00Z', '2026-03-18T17:00:00Z')
      ).toBeNull();
    });

    it('should return null for empty inputs', () => {
      expect(validateTimeRange('', '')).toBeNull();
    });
  });

  describe('sanitizeString', () => {
    it('should trim whitespace', () => {
      expect(sanitizeString('  hello  ')).toBe('hello');
    });
  });
});
