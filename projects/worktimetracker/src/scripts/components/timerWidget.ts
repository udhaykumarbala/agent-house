/**
 * Timer Widget — Persistent timer bar component.
 *
 * Manages the top-bar timer display across all views:
 * - Shows/hides based on active timers
 * - Updates elapsed time every second
 * - Persists timer state to localStorage every 30s
 * - Supports multiple simultaneous timers with expandable list
 * - Recovers active timers on page reload
 */

import { timeEntryService } from '../data/timeEntryService';
import { employeeService } from '../data/employeeService';
import { projectService } from '../data/projectService';
import { eventBus } from '../utils/eventBus';
import { formatDuration } from '../utils/dates';
import { Toast } from './toast';
import type { ActiveTimer } from '../types';

let tickInterval: ReturnType<typeof setInterval> | null = null;
let persistInterval: ReturnType<typeof setInterval> | null = null;
let expanded = false;

/** Format elapsed time as HH:MM:SS from a clock-in ISO timestamp. */
function formatElapsed(clockIn: string): string {
  const diffMs = Date.now() - new Date(clockIn).getTime();
  const totalSec = Math.max(0, Math.floor(diffMs / 1000));
  const h = Math.floor(totalSec / 3600);
  const m = Math.floor((totalSec % 3600) / 60);
  const s = totalSec % 60;
  return `${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`;
}

/** Update the timer bar display. Called every second when timers are active. */
function updateDisplay(): void {
  const timerBar = document.getElementById('timer-bar');
  if (!timerBar) return;

  const timers = timeEntryService.getActiveTimers();

  if (timers.length === 0) {
    timerBar.classList.add('hidden');
    timerBar.classList.remove('timer-active');
    expanded = false;
    // Clean up dynamic elements
    timerBar.querySelector('.timer-count-badge')?.remove();
    timerBar.querySelector('.timer-extra-list')?.remove();
    return;
  }

  timerBar.classList.remove('hidden');
  timerBar.classList.add('timer-active');

  // Update primary timer display (first active timer)
  const primary = timers[0];
  const employee = employeeService.getById(primary.employeeId);
  const project = projectService.getById(primary.projectId);

  const employeeEl = document.getElementById('timer-employee');
  const projectEl = document.getElementById('timer-project');
  const displayEl = document.getElementById('timer-display');

  if (employeeEl) employeeEl.textContent = employee?.name ?? 'Unknown';
  if (projectEl) projectEl.textContent = project?.name ?? 'General';
  if (displayEl) displayEl.textContent = formatElapsed(primary.clockIn);

  // Handle multi-timer badge and expanded list
  updateMultiTimerUI(timers);
}

/** Render/update the multi-timer badge and expandable list. */
function updateMultiTimerUI(timers: ActiveTimer[]): void {
  const timerBar = document.getElementById('timer-bar');
  if (!timerBar) return;

  let badge = timerBar.querySelector<HTMLButtonElement>('.timer-count-badge');
  let extraList = timerBar.querySelector<HTMLElement>('.timer-extra-list');

  // Single timer: remove multi-timer UI
  if (timers.length <= 1) {
    badge?.remove();
    extraList?.remove();
    expanded = false;
    return;
  }

  // Create or update the "+N more" badge
  const leftSection = timerBar.querySelector('.flex.items-center.gap-3');
  if (!badge && leftSection) {
    badge = document.createElement('button');
    badge.className = 'timer-count-badge';
    badge.type = 'button';
    badge.addEventListener('click', () => {
      expanded = !expanded;
      updateDisplay();
    });
    leftSection.appendChild(badge);
  }
  if (badge) {
    badge.textContent = `+${timers.length - 1} more`;
  }

  // Show expanded timer list
  if (expanded) {
    if (!extraList) {
      extraList = document.createElement('div');
      extraList.className = 'timer-extra-list';
      timerBar.appendChild(extraList);
    }

    extraList.innerHTML = '';

    for (let i = 1; i < timers.length; i++) {
      const timer = timers[i];
      const emp = employeeService.getById(timer.employeeId);
      const proj = projectService.getById(timer.projectId);

      const row = document.createElement('div');
      row.className = 'flex items-center justify-between px-4 py-2 lg:px-6';
      row.style.borderTop = '1px solid rgba(91, 70, 178, 0.08)';

      // Left: status dot + name + project
      const left = document.createElement('div');
      left.className = 'flex items-center gap-3 min-w-0';

      const dot = document.createElement('span');
      dot.className = 'status-dot active shrink-0';

      const name = document.createElement('span');
      name.className = 'text-sm font-medium text-content truncate';
      name.textContent = emp?.name ?? 'Unknown';

      const dash = document.createElement('span');
      dash.className = 'text-content-secondary hidden sm:inline';
      dash.textContent = '\u2014';

      const projLabel = document.createElement('span');
      projLabel.className = 'text-sm text-content-secondary truncate hidden sm:inline';
      projLabel.textContent = proj?.name ?? 'General';

      left.appendChild(dot);
      left.appendChild(name);
      left.appendChild(dash);
      left.appendChild(projLabel);

      // Right: elapsed time + stop button
      const right = document.createElement('div');
      right.className = 'flex items-center gap-4 shrink-0';

      const timeEl = document.createElement('span');
      timeEl.className = 'text-timer-compact tabular-nums text-accent';
      timeEl.textContent = formatElapsed(timer.clockIn);

      const stopBtn = document.createElement('button');
      stopBtn.className = 'btn-clock-out text-sm';
      stopBtn.type = 'button';
      stopBtn.textContent = 'Stop';
      const eid = timer.employeeId;
      stopBtn.addEventListener('click', () => handleClockOut(eid));

      right.appendChild(timeEl);
      right.appendChild(stopBtn);

      row.appendChild(left);
      row.appendChild(right);
      extraList.appendChild(row);
    }
  } else {
    extraList?.remove();
  }
}

/** Clock out an employee, save the time entry, and show feedback toast. */
export function handleClockOut(employeeId: string): void {
  try {
    const employee = employeeService.getById(employeeId);
    const entry = timeEntryService.clockOut(employeeId);
    const dur = formatDuration(entry.duration ?? 0);
    const project = projectService.getById(entry.projectId);
    Toast.success(
      `${employee?.name ?? 'Employee'} clocked out \u2014 ${dur} on ${project?.name ?? 'General'}`
    );
  } catch (e) {
    Toast.error((e as Error).message);
  }
}

/** Clock in an employee and show feedback toast. */
export function handleClockIn(
  employeeId: string,
  projectId: string,
  taskId?: string | null,
  notes?: string
): void {
  try {
    timeEntryService.clockIn(employeeId, projectId, taskId, notes);
    const employee = employeeService.getById(employeeId);
    Toast.success(`${employee?.name ?? 'Employee'} clocked in`);
  } catch (e) {
    Toast.error((e as Error).message);
  }
}

function startTicking(): void {
  if (tickInterval) return;
  tickInterval = setInterval(updateDisplay, 1000);
}

function stopTicking(): void {
  if (tickInterval) {
    clearInterval(tickInterval);
    tickInterval = null;
  }
}

function startPersisting(): void {
  if (persistInterval) return;
  persistInterval = setInterval(() => {
    timeEntryService.persistTimerState();
  }, 30000);
}

function stopPersisting(): void {
  if (persistInterval) {
    clearInterval(persistInterval);
    persistInterval = null;
  }
}

function handleTimerChanged(): void {
  const timers = timeEntryService.getActiveTimers();
  updateDisplay();

  if (timers.length > 0) {
    startTicking();
    startPersisting();
  } else {
    stopTicking();
    stopPersisting();
  }
}

/** Initialize the timer widget. Call once after DOM ready and store initialized. */
export function initTimerWidget(): void {
  // Wire up the primary clock-out button in the static HTML timer bar
  const stopBtn = document.getElementById('timer-stop-btn');
  if (stopBtn) {
    stopBtn.addEventListener('click', () => {
      const timers = timeEntryService.getActiveTimers();
      if (timers.length > 0) {
        handleClockOut(timers[0].employeeId);
      }
    });
  }

  // Subscribe to timer lifecycle events
  eventBus.on('timer:changed', handleTimerChanged);

  // Recover active timers from localStorage on page load
  handleTimerChanged();
}

/** Tear down the timer widget (intervals + event listeners). */
export function destroyTimerWidget(): void {
  stopTicking();
  stopPersisting();
  eventBus.off('timer:changed', handleTimerChanged);
}
