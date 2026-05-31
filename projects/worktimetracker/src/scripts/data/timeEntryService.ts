import { store } from './store';
import { eventBus } from '../utils/eventBus';
import { validateRequired, validateTimeRange } from '../utils/sanitize';
import { calcDuration, todayDate } from '../utils/dates';
import type { TimeEntry, ActiveTimer } from '../types';

const ENTRIES_KEY = 'time_entries';
const ACTIVE_TIMERS_KEY = 'active_timers';

function generateId(): string {
  return crypto.randomUUID();
}

export interface CreateManualEntryInput {
  employeeId: string;
  projectId: string;
  taskId?: string | null;
  date: string;
  clockIn: string;
  clockOut: string;
  notes?: string;
}

export interface UpdateEntryInput {
  projectId?: string;
  taskId?: string | null;
  clockIn?: string;
  clockOut?: string;
  notes?: string;
}

class TimeEntryService {
  // ---- Time Entries ----

  getAll(): TimeEntry[] {
    return store.get<TimeEntry>(ENTRIES_KEY);
  }

  getById(id: string): TimeEntry | undefined {
    return this.getAll().find(e => e.id === id);
  }

  getByDate(date: string): TimeEntry[] {
    return this.getAll().filter(e => e.date === date);
  }

  getByEmployee(employeeId: string): TimeEntry[] {
    return this.getAll().filter(e => e.employeeId === employeeId);
  }

  getByDateRange(startDate: string, endDate: string): TimeEntry[] {
    return this.getAll().filter(e => e.date >= startDate && e.date <= endDate);
  }

  getTodayEntries(): TimeEntry[] {
    return this.getByDate(todayDate());
  }

  /** Calculate total hours for an employee today. */
  getTodayHoursForEmployee(employeeId: string): number {
    const today = todayDate();
    const entries = this.getAll().filter(
      e => e.employeeId === employeeId && e.date === today
    );

    let total = 0;
    for (const entry of entries) {
      if (entry.duration !== null) {
        total += entry.duration;
      } else if (entry.clockOut === null) {
        // Active timer: calculate from clockIn to now
        total += calcDuration(entry.clockIn, new Date().toISOString());
      }
    }
    return total;
  }

  /** Calculate total hours today across all employees. */
  getTodayTotalHours(): number {
    const today = todayDate();
    const entries = this.getAll().filter(e => e.date === today);

    let total = 0;
    for (const entry of entries) {
      if (entry.duration !== null) {
        total += entry.duration;
      } else if (entry.clockOut === null) {
        total += calcDuration(entry.clockIn, new Date().toISOString());
      }
    }
    return total;
  }

  createManualEntry(input: CreateManualEntryInput): TimeEntry {
    const empErr = validateRequired(input.employeeId, 'Employee');
    if (empErr) throw new Error(empErr);

    const projErr = validateRequired(input.projectId, 'Project');
    if (projErr) throw new Error(projErr);

    const timeErr = validateTimeRange(input.clockIn, input.clockOut);
    if (timeErr) throw new Error(timeErr);

    // Check overlap
    const overlapErr = this.checkOverlap(
      input.employeeId,
      input.clockIn,
      input.clockOut
    );
    if (overlapErr) throw new Error(overlapErr);

    const duration = calcDuration(input.clockIn, input.clockOut);

    const entry: TimeEntry = {
      id: generateId(),
      employeeId: input.employeeId,
      projectId: input.projectId,
      taskId: input.taskId || null,
      date: input.date || todayDate(),
      clockIn: input.clockIn,
      clockOut: input.clockOut,
      duration,
      notes: input.notes?.trim() || '',
      type: 'manual',
    };

    const entries = this.getAll();
    entries.push(entry);
    store.set(ENTRIES_KEY, entries);

    eventBus.emit('timeEntries:changed', entry);
    eventBus.emit('timeEntries:created', entry);

    return entry;
  }

  update(id: string, input: UpdateEntryInput): TimeEntry {
    const entries = this.getAll();
    const index = entries.findIndex(e => e.id === id);
    if (index === -1) throw new Error('Time entry not found');

    const existing = entries[index];
    const newClockIn = input.clockIn || existing.clockIn;
    const newClockOut = input.clockOut || existing.clockOut;

    if (newClockOut) {
      const timeErr = validateTimeRange(newClockIn, newClockOut);
      if (timeErr) throw new Error(timeErr);

      // Check overlap (excluding self)
      const overlapErr = this.checkOverlap(
        existing.employeeId,
        newClockIn,
        newClockOut,
        id
      );
      if (overlapErr) throw new Error(overlapErr);
    }

    const duration = newClockOut ? calcDuration(newClockIn, newClockOut) : null;

    const updated: TimeEntry = {
      ...existing,
      ...(input.projectId !== undefined && { projectId: input.projectId }),
      ...(input.taskId !== undefined && { taskId: input.taskId }),
      ...(input.clockIn !== undefined && { clockIn: input.clockIn }),
      ...(input.clockOut !== undefined && { clockOut: input.clockOut }),
      ...(input.notes !== undefined && { notes: input.notes.trim() }),
      duration,
    };

    entries[index] = updated;
    store.set(ENTRIES_KEY, entries);

    eventBus.emit('timeEntries:changed', updated);
    eventBus.emit('timeEntries:updated', updated);

    return updated;
  }

  delete(id: string): void {
    const entries = this.getAll();
    const filtered = entries.filter(e => e.id !== id);
    if (filtered.length === entries.length) throw new Error('Time entry not found');

    store.set(ENTRIES_KEY, filtered);
    eventBus.emit('timeEntries:changed');
    eventBus.emit('timeEntries:deleted', id);
  }

  // ---- Active Timers (Clock In/Out) ----

  getActiveTimers(): ActiveTimer[] {
    return store.get<ActiveTimer>(ACTIVE_TIMERS_KEY);
  }

  isEmployeeClockedIn(employeeId: string): boolean {
    return this.getActiveTimers().some(t => t.employeeId === employeeId);
  }

  clockIn(employeeId: string, projectId: string, taskId?: string | null, notes?: string): ActiveTimer {
    if (this.isEmployeeClockedIn(employeeId)) {
      throw new Error('Employee is already clocked in');
    }

    const timer: ActiveTimer = {
      employeeId,
      projectId,
      taskId: taskId || null,
      clockIn: new Date().toISOString(),
      notes: notes?.trim() || '',
    };

    const timers = this.getActiveTimers();
    timers.push(timer);
    store.set(ACTIVE_TIMERS_KEY, timers);

    eventBus.emit('timer:started', timer);
    eventBus.emit('timer:changed');

    return timer;
  }

  clockOut(employeeId: string): TimeEntry {
    const timers = this.getActiveTimers();
    const timerIndex = timers.findIndex(t => t.employeeId === employeeId);
    if (timerIndex === -1) {
      throw new Error('No active timer for this employee');
    }

    const timer = timers[timerIndex];
    const clockOut = new Date().toISOString();
    const duration = calcDuration(timer.clockIn, clockOut);

    // Create the time entry
    const entry: TimeEntry = {
      id: generateId(),
      employeeId: timer.employeeId,
      projectId: timer.projectId,
      taskId: timer.taskId,
      date: todayDate(),
      clockIn: timer.clockIn,
      clockOut,
      duration,
      notes: timer.notes,
      type: 'clock',
    };

    const entries = this.getAll();
    entries.push(entry);
    store.set(ENTRIES_KEY, entries);

    // Remove the active timer
    timers.splice(timerIndex, 1);
    store.set(ACTIVE_TIMERS_KEY, timers);

    eventBus.emit('timer:stopped', entry);
    eventBus.emit('timer:changed');
    eventBus.emit('timeEntries:changed', entry);
    eventBus.emit('timeEntries:created', entry);

    return entry;
  }

  /** Persist active timer state (called periodically for crash recovery). */
  persistTimerState(): void {
    // Active timers are already stored in localStorage via store.set,
    // so this is a no-op write to ensure the cache is flushed.
    const timers = this.getActiveTimers();
    store.set(ACTIVE_TIMERS_KEY, timers);
  }

  // ---- Overlap Detection ----

  private checkOverlap(
    employeeId: string,
    clockIn: string,
    clockOut: string,
    excludeId?: string
  ): string | null {
    const startMs = new Date(clockIn).getTime();
    const endMs = new Date(clockOut).getTime();

    const entries = this.getAll().filter(
      e => e.employeeId === employeeId && e.id !== excludeId && e.clockOut !== null
    );

    for (const entry of entries) {
      const eStart = new Date(entry.clockIn).getTime();
      const eEnd = new Date(entry.clockOut!).getTime();

      // Two ranges overlap if start < otherEnd AND end > otherStart
      if (startMs < eEnd && endMs > eStart) {
        return `Overlaps with existing entry (${new Date(entry.clockIn).toLocaleTimeString()} - ${new Date(entry.clockOut!).toLocaleTimeString()})`;
      }
    }

    return null;
  }
}

export const timeEntryService = new TimeEntryService();
