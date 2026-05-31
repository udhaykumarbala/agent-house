/**
 * CSV export with formula injection protection.
 * - All cell values wrapped in double quotes
 * - Formula-dangerous characters (=, +, -, @, \t, \r) prefixed with single quote
 * - UTF-8 BOM for Excel compatibility
 * - Double quotes within values escaped as ""
 */

const FORMULA_CHARS = ['=', '+', '-', '@', '\t', '\r'];

function sanitizeCSVValue(value: string): string {
  let escaped = value.replace(/"/g, '""');
  if (FORMULA_CHARS.some((ch) => escaped.startsWith(ch))) {
    escaped = "'" + escaped;
  }
  return `"${escaped}"`;
}

export function generateCSV(headers: string[], rows: string[][]): string {
  const BOM = '\uFEFF';
  const headerLine = headers.map((h) => sanitizeCSVValue(h)).join(',');
  const dataLines = rows.map((row) =>
    row.map((cell) => sanitizeCSVValue(cell)).join(',')
  );
  return BOM + [headerLine, ...dataLines].join('\n');
}

export function downloadCSV(csv: string, filename: string): void {
  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' });
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = filename;
  link.style.display = 'none';
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  URL.revokeObjectURL(url);
}
