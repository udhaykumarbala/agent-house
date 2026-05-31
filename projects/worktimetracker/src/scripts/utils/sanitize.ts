/**
 * Input sanitization and validation utilities.
 * escapeHTML uses textContent/innerHTML DOM approach for safe escaping.
 */

export function escapeHTML(str: string): string {
  const div = document.createElement('div');
  div.textContent = str;
  return div.innerHTML;
}

export function validateRequired(value: string, fieldName: string): string | null {
  if (!value || value.trim().length === 0) {
    return `${fieldName} is required`;
  }
  return null;
}

export function validateMaxLength(
  value: string,
  max: number,
  fieldName: string
): string | null {
  if (value && value.length > max) {
    return `${fieldName} must be ${max} characters or less`;
  }
  return null;
}

export function validateEmail(value: string): string | null {
  if (!value) return null;
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if (!emailRegex.test(value)) {
    return 'Invalid email address';
  }
  return null;
}

export function validateTimeRange(
  clockIn: string,
  clockOut: string
): string | null {
  if (!clockIn || !clockOut) return null;
  const start = new Date(clockIn);
  const end = new Date(clockOut);
  if (end <= start) {
    return 'Clock out must be after clock in';
  }
  const diffHours = (end.getTime() - start.getTime()) / (1000 * 60 * 60);
  if (diffHours > 24) {
    return 'Time entry cannot exceed 24 hours';
  }
  return null;
}

/** Trim and sanitize a string for storage. */
export function sanitizeString(value: string): string {
  return value.trim();
}
