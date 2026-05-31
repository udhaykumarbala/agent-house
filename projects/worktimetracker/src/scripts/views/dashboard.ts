/**
 * Dashboard View
 *
 * - 4 summary cards (total hours today, employees clocked in, active projects, entries today)
 * - Active timers section listing clocked-in employees with running time
 * - Quick clock-in: select employee + project, start timer
 * - Recent activity timeline (last 10 entries)
 * - Empty state for new users
 * - Updates on timer start/stop via eventBus
 */

import type { ViewModule } from '../types';
import { employeeService } from '../data/employeeService';
import { projectService } from '../data/projectService';
import { timeEntryService } from '../data/timeEntryService';
import { eventBus } from '../utils/eventBus';
import { formatDuration, formatTime, calcDuration } from '../utils/dates';
import { Toast } from '../components/toast';
import { createIcons } from 'lucide';

let containerRef: HTMLElement | null = null;
let timerInterval: ReturnType<typeof setInterval> | null = null;

function formatElapsed(clockInISO: string): string {
  const ms = Date.now() - new Date(clockInISO).getTime();
  const totalSec = Math.max(0, Math.floor(ms / 1000));
  const h = Math.floor(totalSec / 3600);
  const m = Math.floor((totalSec % 3600) / 60);
  const s = totalSec % 60;
  return `${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`;
}

function updateTimerDisplays(): void {
  if (!containerRef) return;
  const timerEls = containerRef.querySelectorAll('[data-timer-clockin]');
  timerEls.forEach((el) => {
    const clockIn = el.getAttribute('data-timer-clockin');
    if (clockIn) {
      el.textContent = formatElapsed(clockIn);
    }
  });

  // Also update today's hours card since active timers contribute
  const hoursEl = containerRef.querySelector('[data-stat="total-hours"]');
  if (hoursEl) {
    const totalMinutes = timeEntryService.getTodayTotalHours();
    hoursEl.textContent = formatDuration(totalMinutes);
  }
}

function handleClockIn(employeeId: string, projectId: string): void {
  try {
    timeEntryService.clockIn(employeeId, projectId);
    const emp = employeeService.getById(employeeId);
    const proj = projectService.getById(projectId);
    Toast.success(`${emp?.name || 'Employee'} clocked in on ${proj?.name || 'project'}`);
    renderView();
  } catch (err) {
    Toast.error((err as Error).message);
  }
}

function handleClockOut(employeeId: string): void {
  try {
    const entry = timeEntryService.clockOut(employeeId);
    const emp = employeeService.getById(entry.employeeId);
    const proj = projectService.getById(entry.projectId);
    const dur = entry.duration !== null ? formatDuration(entry.duration) : '';
    Toast.success(`${emp?.name || 'Employee'} clocked out — ${dur} on ${proj?.name || 'project'}`);
    renderView();
  } catch (err) {
    Toast.error((err as Error).message);
  }
}

function buildSummaryCards(): HTMLElement {
  const grid = document.createElement('div');
  grid.style.cssText = 'display: grid; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); gap: 16px; margin-bottom: 24px;';

  const totalMinutes = timeEntryService.getTodayTotalHours();
  const activeTimers = timeEntryService.getActiveTimers();
  const totalEmployees = employeeService.getActive().length;
  const activeProjectCount = projectService.getActive().length;
  const todayEntries = timeEntryService.getTodayEntries();

  const cards: { icon: string; value: string; label: string; statKey?: string }[] = [
    {
      icon: 'clock',
      value: formatDuration(totalMinutes),
      label: 'TOTAL HOURS TODAY',
      statKey: 'total-hours',
    },
    {
      icon: 'users',
      value: `${activeTimers.length}/${totalEmployees}`,
      label: 'CLOCKED IN',
    },
    {
      icon: 'folder-kanban',
      value: String(activeProjectCount),
      label: 'ACTIVE PROJECTS',
    },
    {
      icon: 'list-checks',
      value: String(todayEntries.length),
      label: 'ENTRIES TODAY',
    },
  ];

  for (const c of cards) {
    const card = document.createElement('div');
    card.className = 'stat-card';

    const iconRow = document.createElement('div');
    iconRow.style.cssText = 'display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px;';

    const iconWrap = document.createElement('div');
    iconWrap.className = 'stat-icon';
    iconWrap.innerHTML = `<i data-lucide="${c.icon}" class="w-6 h-6"></i>`;

    iconRow.appendChild(iconWrap);
    card.appendChild(iconRow);

    const val = document.createElement('div');
    val.className = 'stat-value';
    val.textContent = c.value;
    if (c.statKey) val.setAttribute('data-stat', c.statKey);

    const label = document.createElement('div');
    label.className = 'stat-label';
    label.textContent = c.label;

    card.appendChild(val);
    card.appendChild(label);
    grid.appendChild(card);
  }

  return grid;
}

function buildActiveTimers(): HTMLElement | null {
  const activeTimers = timeEntryService.getActiveTimers();
  if (activeTimers.length === 0) return null;

  const section = document.createElement('div');
  section.style.marginBottom = '24px';

  const heading = document.createElement('h2');
  heading.style.cssText = 'font-size: 18px; font-weight: 600; color: #ECE9F5; margin-bottom: 12px;';
  heading.textContent = 'Active Timers';
  section.appendChild(heading);

  const list = document.createElement('div');
  list.style.cssText = 'display: flex; flex-direction: column; gap: 8px;';

  for (const timer of activeTimers) {
    const emp = employeeService.getById(timer.employeeId);
    const proj = projectService.getById(timer.projectId);

    const row = document.createElement('div');
    row.className = 'card';
    row.style.cssText = 'display: flex; align-items: center; justify-content: space-between; padding: 12px 16px; gap: 12px;';

    const left = document.createElement('div');
    left.style.cssText = 'display: flex; align-items: center; gap: 10px; min-width: 0; flex: 1;';

    const dot = document.createElement('span');
    dot.className = 'status-dot active';
    dot.setAttribute('aria-label', 'Clocked in');

    const info = document.createElement('div');
    info.style.cssText = 'min-width: 0;';

    const name = document.createElement('div');
    name.style.cssText = 'font-size: 14px; font-weight: 500; color: #ECE9F5; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;';
    name.textContent = emp?.name || 'Unknown';

    const project = document.createElement('div');
    project.style.cssText = 'font-size: 12px; color: #8C83A8;';
    project.textContent = proj?.name || 'General';

    info.appendChild(name);
    info.appendChild(project);
    left.appendChild(dot);
    left.appendChild(info);

    const right = document.createElement('div');
    right.style.cssText = 'display: flex; align-items: center; gap: 12px; flex-shrink: 0;';

    const elapsed = document.createElement('span');
    elapsed.style.cssText = 'font-size: 16px; font-weight: 600; color: #CF8A2E; font-variant-numeric: tabular-nums; letter-spacing: 0.02em;';
    elapsed.textContent = formatElapsed(timer.clockIn);
    elapsed.setAttribute('data-timer-clockin', timer.clockIn);

    const stopBtn = document.createElement('button');
    stopBtn.type = 'button';
    stopBtn.className = 'btn-clock-out btn-sm';
    stopBtn.innerHTML = '<i data-lucide="square" class="w-3.5 h-3.5 mr-1 inline-block"></i>Stop';
    stopBtn.addEventListener('click', () => handleClockOut(timer.employeeId));

    right.appendChild(elapsed);
    right.appendChild(stopBtn);

    row.appendChild(left);
    row.appendChild(right);
    list.appendChild(row);
  }

  section.appendChild(list);
  return section;
}

function buildQuickClockIn(): HTMLElement {
  const section = document.createElement('div');
  section.className = 'card';
  section.style.cssText = 'margin-bottom: 24px; padding: 20px;';

  const heading = document.createElement('h2');
  heading.style.cssText = 'font-size: 18px; font-weight: 600; color: #ECE9F5; margin-bottom: 16px;';
  heading.textContent = 'Quick Clock In';
  section.appendChild(heading);

  const form = document.createElement('div');
  form.style.cssText = 'display: flex; gap: 12px; flex-wrap: wrap; align-items: flex-end;';

  // Employee select
  const empWrap = document.createElement('div');
  empWrap.style.cssText = 'flex: 1; min-width: 160px;';
  const empLabel = document.createElement('label');
  empLabel.className = 'form-label';
  empLabel.textContent = 'Employee';
  const empSelect = document.createElement('select');
  empSelect.className = 'select w-full';
  empSelect.id = 'quick-employee';

  const empDefault = document.createElement('option');
  empDefault.value = '';
  empDefault.textContent = 'Select employee...';
  empSelect.appendChild(empDefault);

  const activeEmployees = employeeService.getActive();
  const activeTimers = timeEntryService.getActiveTimers();
  const clockedInIds = new Set(activeTimers.map(t => t.employeeId));

  for (const emp of activeEmployees) {
    if (clockedInIds.has(emp.id)) continue; // Skip already clocked in
    const opt = document.createElement('option');
    opt.value = emp.id;
    opt.textContent = emp.name;
    empSelect.appendChild(opt);
  }

  empWrap.appendChild(empLabel);
  empWrap.appendChild(empSelect);

  // Project select
  const projWrap = document.createElement('div');
  projWrap.style.cssText = 'flex: 1; min-width: 160px;';
  const projLabel = document.createElement('label');
  projLabel.className = 'form-label';
  projLabel.textContent = 'Project';
  const projSelect = document.createElement('select');
  projSelect.className = 'select w-full';
  projSelect.id = 'quick-project';

  const selectableProjects = projectService.getSelectable();
  for (const proj of selectableProjects) {
    const opt = document.createElement('option');
    opt.value = proj.id;
    opt.textContent = proj.name;
    if (proj.id === 'general') opt.selected = true;
    projSelect.appendChild(opt);
  }

  projWrap.appendChild(projLabel);
  projWrap.appendChild(projSelect);

  // Clock In button
  const btnWrap = document.createElement('div');
  const clockInBtn = document.createElement('button');
  clockInBtn.type = 'button';
  clockInBtn.className = 'btn-accent btn-lg';
  clockInBtn.innerHTML = '<i data-lucide="play" class="w-4 h-4 mr-1.5 inline-block"></i>Clock In';
  clockInBtn.addEventListener('click', () => {
    const employeeId = empSelect.value;
    const projectId = projSelect.value;
    if (!employeeId) {
      Toast.warning('Please select an employee');
      return;
    }
    handleClockIn(employeeId, projectId || 'general');
  });

  // Disable if no employees available
  if (activeEmployees.length === 0 || activeEmployees.every(e => clockedInIds.has(e.id))) {
    clockInBtn.disabled = true;
    empSelect.disabled = true;
  }

  btnWrap.appendChild(clockInBtn);

  form.appendChild(empWrap);
  form.appendChild(projWrap);
  form.appendChild(btnWrap);
  section.appendChild(form);

  return section;
}

function buildRecentActivity(): HTMLElement | null {
  const todayEntries = timeEntryService.getTodayEntries();
  // Sort newest first, limit to 10
  const sorted = [...todayEntries]
    .sort((a, b) => new Date(b.clockIn).getTime() - new Date(a.clockIn).getTime())
    .slice(0, 10);

  if (sorted.length === 0) return null;

  const section = document.createElement('div');

  const heading = document.createElement('h2');
  heading.style.cssText = 'font-size: 18px; font-weight: 600; color: #ECE9F5; margin-bottom: 12px;';
  heading.textContent = "Today's Activity";
  section.appendChild(heading);

  const list = document.createElement('div');
  list.className = 'card';
  list.style.cssText = 'padding: 0; overflow: hidden;';

  for (let i = 0; i < sorted.length; i++) {
    const entry = sorted[i];
    const emp = employeeService.getById(entry.employeeId);
    const proj = projectService.getById(entry.projectId);
    const isActive = entry.clockOut === null;

    const row = document.createElement('div');
    row.style.cssText = `
      display: flex; align-items: center; gap: 12px; padding: 12px 16px;
      ${i < sorted.length - 1 ? 'border-bottom: 1px solid rgba(91, 70, 178, 0.08);' : ''}
      transition: background 150ms ease;
    `;
    row.addEventListener('mouseenter', () => { row.style.background = 'rgba(87, 70, 178, 0.06)'; });
    row.addEventListener('mouseleave', () => { row.style.background = 'transparent'; });

    // Status dot
    const dot = document.createElement('span');
    dot.className = isActive ? 'status-dot active' : 'status-dot inactive';
    dot.setAttribute('aria-label', isActive ? 'Active' : 'Completed');

    // Employee name
    const empName = document.createElement('span');
    empName.style.cssText = 'font-size: 14px; font-weight: 500; color: #ECE9F5; min-width: 100px;';
    empName.textContent = emp?.name || 'Unknown';

    // Project name
    const projName = document.createElement('span');
    projName.style.cssText = 'font-size: 13px; color: #8C83A8; flex: 1; min-width: 80px;';
    projName.textContent = proj?.name || 'General';

    // Time range
    const timeRange = document.createElement('span');
    timeRange.style.cssText = 'font-size: 13px; color: #8C83A8; font-variant-numeric: tabular-nums; white-space: nowrap;';
    const startTime = formatTime(entry.clockIn);
    if (isActive) {
      timeRange.textContent = `${startTime} → ...`;
    } else {
      timeRange.textContent = `${startTime} → ${formatTime(entry.clockOut!)}`;
    }

    // Duration
    const duration = document.createElement('span');
    duration.style.cssText = 'font-size: 13px; font-weight: 500; font-variant-numeric: tabular-nums; white-space: nowrap; min-width: 60px; text-align: right;';
    if (isActive) {
      const mins = calcDuration(entry.clockIn, new Date().toISOString());
      duration.style.color = '#CF8A2E';
      duration.textContent = formatDuration(mins);
      duration.setAttribute('data-timer-clockin', entry.clockIn);
    } else {
      duration.style.color = '#ECE9F5';
      duration.textContent = entry.duration !== null ? formatDuration(entry.duration) : '';
    }

    row.appendChild(dot);
    row.appendChild(empName);
    row.appendChild(projName);
    row.appendChild(timeRange);
    row.appendChild(duration);
    list.appendChild(row);
  }

  section.appendChild(list);
  return section;
}

function buildEmptyState(): HTMLElement {
  const empty = document.createElement('div');
  empty.className = 'empty-state';
  empty.innerHTML = `
    <div class="empty-icon">
      <i data-lucide="layout-dashboard" class="w-12 h-12 mx-auto text-content-disabled"></i>
    </div>
    <div class="empty-title">Welcome to WorkTime Tracker!</div>
    <p class="empty-description">Add your first employee to get started tracking time.</p>
  `;
  const ctaBtn = document.createElement('a');
  ctaBtn.href = '#/employees';
  ctaBtn.className = 'btn-accent';
  ctaBtn.innerHTML = '<i data-lucide="plus" class="w-4 h-4 mr-1.5 inline-block"></i>Add Employee';
  empty.appendChild(ctaBtn);
  return empty;
}

function renderView(): void {
  if (!containerRef) return;
  containerRef.innerHTML = '';

  const employees = employeeService.getAll();
  const hasEmployees = employees.length > 0;

  // Header
  const header = document.createElement('div');
  header.style.cssText = 'display: flex; align-items: center; justify-content: space-between; margin-bottom: 24px;';

  const titleWrap = document.createElement('div');
  const h1 = document.createElement('h1');
  h1.className = 'text-content font-heading';
  h1.textContent = 'Dashboard';

  const dateStr = new Intl.DateTimeFormat('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  }).format(new Date());

  const subtitle = document.createElement('p');
  subtitle.className = 'text-content-secondary text-sm';
  subtitle.style.marginTop = '4px';
  subtitle.textContent = dateStr;

  titleWrap.appendChild(h1);
  titleWrap.appendChild(subtitle);
  header.appendChild(titleWrap);
  containerRef.appendChild(header);

  if (!hasEmployees) {
    containerRef.appendChild(buildEmptyState());
    try { createIcons(); } catch { /* ok */ }
    return;
  }

  // Summary cards
  containerRef.appendChild(buildSummaryCards());

  // Active timers
  const activeTimersSection = buildActiveTimers();
  if (activeTimersSection) {
    containerRef.appendChild(activeTimersSection);
  }

  // Quick clock in
  containerRef.appendChild(buildQuickClockIn());

  // Recent activity
  const activity = buildRecentActivity();
  if (activity) {
    containerRef.appendChild(activity);
  } else {
    // No activity today message
    const noActivity = document.createElement('div');
    noActivity.style.cssText = 'text-align: center; padding: 24px; color: #8C83A8; font-size: 14px;';
    noActivity.textContent = 'No activity yet today. Clock someone in to start tracking.';
    containerRef.appendChild(noActivity);
  }

  // Re-init lucide icons
  try { createIcons(); } catch { /* ok */ }
}

function onDataChange(): void {
  renderView();
}

function startTimerUpdates(): void {
  if (timerInterval) clearInterval(timerInterval);
  timerInterval = setInterval(updateTimerDisplays, 1000);
}

function stopTimerUpdates(): void {
  if (timerInterval) {
    clearInterval(timerInterval);
    timerInterval = null;
  }
}

export const render: ViewModule['render'] = (container) => {
  containerRef = container;

  eventBus.on('timer:changed', onDataChange);
  eventBus.on('timeEntries:changed', onDataChange);
  eventBus.on('employees:changed', onDataChange);
  eventBus.on('projects:changed', onDataChange);

  renderView();
  startTimerUpdates();
};

export const destroy = (): void => {
  eventBus.off('timer:changed', onDataChange);
  eventBus.off('timeEntries:changed', onDataChange);
  eventBus.off('employees:changed', onDataChange);
  eventBus.off('projects:changed', onDataChange);
  stopTimerUpdates();
  containerRef = null;
};
