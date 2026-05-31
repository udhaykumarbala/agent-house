/**
 * Date/time formatting and duration utilities.
 * Times stored in UTC (ISO 8601), displayed in local via Intl.DateTimeFormat.
 */

export function formatDate(date: Date | string): string {
  const d = typeof date === 'string' ? new Date(date) : date;
  return new Intl.DateTimeFormat('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  }).format(d);
}

export function formatTime(date: Date | string): string {
  const d = typeof date === 'string' ? new Date(date) : date;
  return new Intl.DateTimeFormat('en-US', {
    hour: '2-digit',
    minute: '2-digit',
    hour12: true,
  }).format(d);
}

export function formatDuration(minutes: number): string {
  if (minutes < 0) return '0h 0m';
  const h = Math.floor(minutes / 60);
  const m = Math.round(minutes % 60);
  return `${h}h ${m}m`;
}

export function calcDuration(start: Date | string, end: Date | string): number {
  const s = typeof start === 'string' ? new Date(start) : start;
  const e = typeof end === 'string' ? new Date(end) : end;
  return (e.getTime() - s.getTime()) / 60000;
}

export function toUTC(date: Date): string {
  return date.toISOString();
}

export function toLocal(isoString: string): Date {
  return new Date(isoString);
}

/** Returns today's date as YYYY-MM-DD string in local timezone. */
export function todayDate(): string {
  const d = new Date();
  const year = d.getFullYear();
  const month = String(d.getMonth() + 1).padStart(2, '0');
  const day = String(d.getDate()).padStart(2, '0');
  return `${year}-${month}-${day}`;
}
