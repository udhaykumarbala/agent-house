export function formatResult(value: number): string {
  if (!Number.isFinite(value)) return '';
  if (value === 0) return '0';

  const abs = Math.abs(value);

  // Scientific notation for very large or very small numbers
  if (abs >= 1e15 || (abs > 0 && abs < 1e-6)) {
    return value.toExponential(4);
  }

  // Up to 6 significant digits, trimmed trailing zeros
  const formatted = Number(value.toPrecision(6));
  return String(formatted);
}

export function sanitizeNumericInput(raw: string): string | null {
  if (raw === '' || raw === '-') return raw;

  // Allow digits, one decimal point, leading minus
  let sanitized = '';
  let hasDecimal = false;
  let hasDigit = false;

  for (let i = 0; i < raw.length; i++) {
    const ch = raw[i];
    if (ch === '-' && i === 0) {
      sanitized += ch;
    } else if (ch === '.' && !hasDecimal) {
      hasDecimal = true;
      sanitized += ch;
    } else if (ch >= '0' && ch <= '9') {
      hasDigit = true;
      sanitized += ch;
    }
  }

  // If we stripped everything meaningful, return null
  if (!hasDigit && sanitized !== '-' && sanitized !== '') return null;

  return sanitized || null;
}
