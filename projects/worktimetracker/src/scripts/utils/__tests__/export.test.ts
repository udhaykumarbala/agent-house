import { describe, it, expect } from 'vitest';
import { generateCSV } from '../export';

describe('export', () => {
  describe('generateCSV', () => {
    it('should start with UTF-8 BOM', () => {
      const csv = generateCSV(['Name'], [['John']]);
      expect(csv.charCodeAt(0)).toBe(0xfeff);
    });

    it('should wrap all values in double quotes', () => {
      const csv = generateCSV(['Name', 'Hours'], [['John', '8']]);
      const lines = csv.split('\n');
      expect(lines[0]).toContain('"Name"');
      expect(lines[1]).toContain('"John"');
      expect(lines[1]).toContain('"8"');
    });

    it('should escape double quotes within values', () => {
      const csv = generateCSV(['Note'], [['He said "hello"']]);
      expect(csv).toContain('He said ""hello""');
    });

    it('should prefix formula-dangerous characters with single quote', () => {
      const csv = generateCSV(['Formula'], [['=SUM(A1:A10)']]);
      expect(csv).toContain("\"'=SUM(A1:A10)\"");
    });

    it('should prefix + character', () => {
      const csv = generateCSV(['Val'], [['+100']]);
      expect(csv).toContain("\"'+100\"");
    });

    it('should prefix - character', () => {
      const csv = generateCSV(['Val'], [['-100']]);
      expect(csv).toContain("\"'-100\"");
    });

    it('should prefix @ character', () => {
      const csv = generateCSV(['Val'], [['@mention']]);
      expect(csv).toContain("\"'@mention\"");
    });

    it('should handle multiple rows', () => {
      const csv = generateCSV(
        ['Name', 'Hours'],
        [
          ['Alice', '8'],
          ['Bob', '6'],
        ]
      );
      const lines = csv.split('\n');
      expect(lines.length).toBe(3); // header + 2 rows
    });
  });
});
