/**
 * Timesheet View
 *
 * - Time entry table sorted by date (newest first)
 * - Sortable by any column (date, employee, project, duration)
 * - Add manual entry modal: employee (required), project (required), task, date, start time, end time, notes
 * - Validation: end time > start time, overlap detection, max 24h
 * - Inline editing of duration and notes
 * - Delete entry with confirmation toast
 * - Filter by date range (DateRangePicker), employee dropdown, project dropdown
 * - Duration displayed as Xh Ym format
 * - EventBus listeners for timeEntries:changed
 */

import type { ViewModule, TimeEntry } from '../types';
import { timeEntryService } from '../data/timeEntryService';
import { employeeService } from '../data/employeeService';
import { projectService } from '../data/projectService';
import { eventBus } from '../utils/eventBus';
import { formatDuration, formatDate, formatTime, calcDuration, todayDate } from '../utils/dates';
import { createDateRangePicker, computePresetRange } from '../components/dateRangePicker';
import type { DateRange, PresetKey } from '../components/dateRangePicker';
import { Modal } from '../components/modal';
import { Toast } from '../components/toast';
import { createIcons } from 'lucide';

let containerRef: HTMLElement | null = null;
let currentDateRange: DateRange;
let employeeFilter = 'all';
let projectFilter = 'all';
let sortColumn = 'date';
let sortDirection: 'asc' | 'desc' = 'desc';

// Initialize default date range to this week
const defaultRange = computePresetRange('thisWeek');
currentDateRange = { ...defaultRange, preset: 'thisWeek' as PresetKey };

function getFilteredEntries(): TimeEntry[] {
  let entries = timeEntryService.getByDateRange(currentDateRange.start, currentDateRange.end);

  if (employeeFilter !== 'all') {
    entries = entries.filter(e => e.employeeId === employeeFilter);
  }

  if (projectFilter !== 'all') {
    entries = entries.filter(e => e.projectId === projectFilter);
  }

  // Sort
  entries.sort((a, b) => {
    const dir = sortDirection === 'asc' ? 1 : -1;

    switch (sortColumn) {
      case 'date': {
        const cmp = a.date.localeCompare(b.date) || a.clockIn.localeCompare(b.clockIn);
        return cmp * dir;
      }
      case 'employee': {
        const empA = employeeService.getById(a.employeeId)?.name || '';
        const empB = employeeService.getById(b.employeeId)?.name || '';
        return empA.localeCompare(empB) * dir;
      }
      case 'project': {
        const projA = projectService.getById(a.projectId)?.name || '';
        const projB = projectService.getById(b.projectId)?.name || '';
        return projA.localeCompare(projB) * dir;
      }
      case 'duration': {
        const durA = a.duration ?? 0;
        const durB = b.duration ?? 0;
        return (durA - durB) * dir;
      }
      default:
        return 0;
    }
  });

  return entries;
}

function handleSort(column: string): void {
  if (sortColumn === column) {
    sortDirection = sortDirection === 'asc' ? 'desc' : 'asc';
  } else {
    sortColumn = column;
    sortDirection = column === 'date' ? 'desc' : 'asc';
  }
  renderView();
}

function buildSortHeader(label: string, key: string): HTMLElement {
  const th = document.createElement('th');
  th.setAttribute('scope', 'col');
  th.style.cursor = 'pointer';

  const span = document.createElement('span');
  span.textContent = label;
  th.appendChild(span);

  const indicator = document.createElement('span');
  indicator.className = 'sort-indicator';
  if (sortColumn === key) {
    indicator.classList.add('active');
    indicator.textContent = sortDirection === 'asc' ? ' \u25B2' : ' \u25BC';
  } else {
    indicator.textContent = ' \u25B4';
  }
  th.appendChild(indicator);

  th.addEventListener('click', () => handleSort(key));
  return th;
}

function openAddEntryModal(): void {
  const content = document.createElement('div');
  content.style.cssText = 'display: flex; flex-direction: column; gap: 16px;';

  const employees = employeeService.getActive();
  const projects = projectService.getSelectable();

  // Employee select
  const empGroup = document.createElement('div');
  const empLabel = document.createElement('label');
  empLabel.className = 'form-label';
  empLabel.textContent = 'Employee *';
  empLabel.setAttribute('for', 'entry-employee');
  const empSelect = document.createElement('select');
  empSelect.id = 'entry-employee';
  empSelect.className = 'select w-full';
  const empDefault = document.createElement('option');
  empDefault.value = '';
  empDefault.textContent = 'Select employee...';
  empSelect.appendChild(empDefault);
  for (const emp of employees) {
    const opt = document.createElement('option');
    opt.value = emp.id;
    opt.textContent = emp.name;
    empSelect.appendChild(opt);
  }
  empGroup.appendChild(empLabel);
  empGroup.appendChild(empSelect);

  // Project select
  const projGroup = document.createElement('div');
  const projLabel = document.createElement('label');
  projLabel.className = 'form-label';
  projLabel.textContent = 'Project *';
  projLabel.setAttribute('for', 'entry-project');
  const projSelect = document.createElement('select');
  projSelect.id = 'entry-project';
  projSelect.className = 'select w-full';
  const projDefault = document.createElement('option');
  projDefault.value = '';
  projDefault.textContent = 'Select project...';
  projSelect.appendChild(projDefault);
  for (const proj of projects) {
    const opt = document.createElement('option');
    opt.value = proj.id;
    opt.textContent = proj.name;
    projSelect.appendChild(opt);
  }
  projGroup.appendChild(projLabel);
  projGroup.appendChild(projSelect);

  // Task select (optional, populated based on project)
  const taskGroup = document.createElement('div');
  const taskLabel = document.createElement('label');
  taskLabel.className = 'form-label';
  taskLabel.textContent = 'Task (optional)';
  taskLabel.setAttribute('for', 'entry-task');
  const taskSelect = document.createElement('select');
  taskSelect.id = 'entry-task';
  taskSelect.className = 'select w-full';
  const taskDefault = document.createElement('option');
  taskDefault.value = '';
  taskDefault.textContent = 'No task';
  taskSelect.appendChild(taskDefault);
  taskGroup.appendChild(taskLabel);
  taskGroup.appendChild(taskSelect);

  projSelect.addEventListener('change', () => {
    // Repopulate tasks
    taskSelect.innerHTML = '';
    const noTask = document.createElement('option');
    noTask.value = '';
    noTask.textContent = 'No task';
    taskSelect.appendChild(noTask);
    if (projSelect.value) {
      const tasks = projectService.getTasksByProject(projSelect.value);
      for (const task of tasks) {
        const opt = document.createElement('option');
        opt.value = task.id;
        opt.textContent = task.name;
        taskSelect.appendChild(opt);
      }
    }
  });

  // Date
  const dateGroup = document.createElement('div');
  const dateLabel = document.createElement('label');
  dateLabel.className = 'form-label';
  dateLabel.textContent = 'Date';
  dateLabel.setAttribute('for', 'entry-date');
  const dateInput = document.createElement('input');
  dateInput.type = 'date';
  dateInput.id = 'entry-date';
  dateInput.className = 'input';
  dateInput.value = todayDate();
  dateGroup.appendChild(dateLabel);
  dateGroup.appendChild(dateInput);

  // Time row
  const timeRow = document.createElement('div');
  timeRow.style.cssText = 'display: grid; grid-template-columns: 1fr 1fr; gap: 12px;';

  const startGroup = document.createElement('div');
  const startLabel = document.createElement('label');
  startLabel.className = 'form-label';
  startLabel.textContent = 'Start Time';
  startLabel.setAttribute('for', 'entry-start');
  const startInput = document.createElement('input');
  startInput.type = 'time';
  startInput.id = 'entry-start';
  startInput.className = 'input';
  startInput.value = '09:00';
  startGroup.appendChild(startLabel);
  startGroup.appendChild(startInput);

  const endGroup = document.createElement('div');
  const endLabel = document.createElement('label');
  endLabel.className = 'form-label';
  endLabel.textContent = 'End Time';
  endLabel.setAttribute('for', 'entry-end');
  const endInput = document.createElement('input');
  endInput.type = 'time';
  endInput.id = 'entry-end';
  endInput.className = 'input';
  endInput.value = '17:00';
  endGroup.appendChild(endLabel);
  endGroup.appendChild(endInput);

  timeRow.appendChild(startGroup);
  timeRow.appendChild(endGroup);

  // Notes
  const notesGroup = document.createElement('div');
  const notesLabel = document.createElement('label');
  notesLabel.className = 'form-label';
  notesLabel.textContent = 'Notes (optional)';
  notesLabel.setAttribute('for', 'entry-notes');
  const notesInput = document.createElement('input');
  notesInput.type = 'text';
  notesInput.id = 'entry-notes';
  notesInput.className = 'input';
  notesInput.placeholder = 'Add a note...';
  notesGroup.appendChild(notesLabel);
  notesGroup.appendChild(notesInput);

  // Error display
  const errorDiv = document.createElement('div');
  errorDiv.className = 'form-error';
  errorDiv.style.display = 'none';

  content.appendChild(empGroup);
  content.appendChild(projGroup);
  content.appendChild(taskGroup);
  content.appendChild(dateGroup);
  content.appendChild(timeRow);
  content.appendChild(notesGroup);
  content.appendChild(errorDiv);

  Modal.open({
    title: 'Add Time Entry',
    content,
    submitLabel: 'Save Entry',
    onSubmit: () => {
      errorDiv.style.display = 'none';

      const employeeId = empSelect.value;
      const projectId = projSelect.value;
      const taskId = taskSelect.value || null;
      const date = dateInput.value;
      const startTime = startInput.value;
      const endTime = endInput.value;
      const notes = notesInput.value;

      if (!employeeId) {
        errorDiv.textContent = 'Employee is required';
        errorDiv.style.display = 'block';
        return;
      }
      if (!projectId) {
        errorDiv.textContent = 'Project is required';
        errorDiv.style.display = 'block';
        return;
      }
      if (!date || !startTime || !endTime) {
        errorDiv.textContent = 'Date, start time, and end time are required';
        errorDiv.style.display = 'block';
        return;
      }

      // Build ISO datetime strings
      const clockIn = new Date(`${date}T${startTime}:00`).toISOString();
      const clockOut = new Date(`${date}T${endTime}:00`).toISOString();

      try {
        timeEntryService.createManualEntry({
          employeeId,
          projectId,
          taskId,
          date,
          clockIn,
          clockOut,
          notes,
        });
        Toast.success('Time entry saved');
        Modal.close();
        renderView();
      } catch (err) {
        errorDiv.textContent = (err as Error).message;
        errorDiv.style.display = 'block';
      }
    },
  });
}

function startInlineEdit(td: HTMLElement, entry: TimeEntry, field: 'notes' | 'duration'): void {
  if (td.querySelector('input')) return; // Already editing

  const currentValue = field === 'notes' ? entry.notes : formatDuration(entry.duration ?? 0);

  td.innerHTML = '';
  const input = document.createElement('input');
  input.type = 'text';
  input.className = 'input';
  input.style.cssText = 'padding: 4px 8px; font-size: 13px; width: 100%;';
  input.value = currentValue;

  function save(): void {
    const newValue = input.value.trim();
    try {
      if (field === 'notes') {
        timeEntryService.update(entry.id, { notes: newValue });
        Toast.success('Notes updated');
      } else {
        // Parse duration like "2h 30m" or just minutes
        const parsed = parseDurationInput(newValue);
        if (parsed === null) {
          Toast.error('Invalid duration format. Use "Xh Ym" (e.g., 2h 30m)');
          renderView();
          return;
        }
        // Calculate new clockOut from clockIn + duration
        const clockInMs = new Date(entry.clockIn).getTime();
        const newClockOut = new Date(clockInMs + parsed * 60000).toISOString();
        timeEntryService.update(entry.id, { clockOut: newClockOut });
        Toast.success('Duration updated');
      }
      renderView();
    } catch (err) {
      Toast.error((err as Error).message);
      renderView();
    }
  }

  input.addEventListener('blur', save);
  input.addEventListener('keydown', (e) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      input.blur();
    }
    if (e.key === 'Escape') {
      renderView();
    }
  });

  td.appendChild(input);
  input.focus();
  input.select();
}

function parseDurationInput(str: string): number | null {
  // Parse "Xh Ym", "Xh", "Ym", or just a number (minutes)
  const match = str.match(/^(?:(\d+)\s*h)?\s*(?:(\d+)\s*m)?$/i);
  if (match && (match[1] || match[2])) {
    const hours = parseInt(match[1] || '0', 10);
    const mins = parseInt(match[2] || '0', 10);
    return hours * 60 + mins;
  }
  // Try plain number as minutes
  const num = parseFloat(str);
  if (!isNaN(num) && num >= 0) return num;
  return null;
}

function deleteEntry(entry: TimeEntry): void {
  const empName = employeeService.getById(entry.employeeId)?.name || 'Unknown';
  const projName = projectService.getById(entry.projectId)?.name || 'Unknown';

  try {
    timeEntryService.delete(entry.id);
    Toast.success(`Deleted entry: ${empName} on ${projName}`);
    renderView();
  } catch (err) {
    Toast.error((err as Error).message);
  }
}

function buildEntryRow(entry: TimeEntry): HTMLElement {
  const tr = document.createElement('tr');

  const emp = employeeService.getById(entry.employeeId);
  const proj = projectService.getById(entry.projectId);
  const task = entry.taskId ? projectService.getTaskById(entry.taskId) : null;

  // Date cell
  const dateCell = document.createElement('td');
  dateCell.textContent = formatDate(entry.date + 'T00:00:00');
  tr.appendChild(dateCell);

  // Employee cell
  const empCell = document.createElement('td');
  if (emp) {
    const empInner = document.createElement('div');
    empInner.style.cssText = 'display: flex; align-items: center; gap: 8px;';
    const colorDot = document.createElement('span');
    colorDot.style.cssText = `width: 8px; height: 8px; border-radius: 9999px; background: ${emp.color}; flex-shrink: 0;`;
    const nameSpan = document.createElement('span');
    nameSpan.textContent = emp.name;
    empInner.appendChild(colorDot);
    empInner.appendChild(nameSpan);
    empCell.appendChild(empInner);
  } else {
    empCell.textContent = 'Unknown';
    empCell.style.color = '#564F6D';
  }
  tr.appendChild(empCell);

  // Project cell
  const projCell = document.createElement('td');
  if (proj) {
    const projInner = document.createElement('div');
    const projName = document.createElement('span');
    projName.textContent = proj.name;
    projInner.appendChild(projName);
    if (task) {
      const taskSpan = document.createElement('span');
      taskSpan.style.cssText = 'font-size: 12px; color: #8C83A8; margin-left: 6px;';
      taskSpan.textContent = `/ ${task.name}`;
      projInner.appendChild(taskSpan);
    }
    projCell.appendChild(projInner);
  } else {
    projCell.textContent = 'Unknown';
    projCell.style.color = '#564F6D';
  }
  tr.appendChild(projCell);

  // Time cell (start - end)
  const timeCell = document.createElement('td');
  const startStr = formatTime(entry.clockIn);
  const endStr = entry.clockOut ? formatTime(entry.clockOut) : '...';
  timeCell.textContent = `${startStr} – ${endStr}`;
  timeCell.style.fontVariantNumeric = 'tabular-nums';
  if (!entry.clockOut) {
    timeCell.style.color = '#3DB87A';
  }
  tr.appendChild(timeCell);

  // Duration cell (inline editable)
  const durCell = document.createElement('td');
  durCell.style.fontVariantNumeric = 'tabular-nums';
  if (entry.duration !== null) {
    durCell.textContent = formatDuration(entry.duration);
    durCell.style.cursor = 'pointer';
    durCell.title = 'Click to edit';
    durCell.addEventListener('click', (e) => {
      e.stopPropagation();
      startInlineEdit(durCell, entry, 'duration');
    });
  } else {
    // Active timer — calculate live
    const elapsed = calcDuration(entry.clockIn, new Date().toISOString());
    const span = document.createElement('span');
    span.style.color = '#3DB87A';
    span.textContent = formatDuration(elapsed) + ' ▶';
    durCell.appendChild(span);
  }
  tr.appendChild(durCell);

  // Notes cell (inline editable)
  const notesCell = document.createElement('td');
  notesCell.textContent = entry.notes || '—';
  notesCell.style.color = entry.notes ? '#ECE9F5' : '#564F6D';
  notesCell.style.cursor = 'pointer';
  notesCell.style.maxWidth = '200px';
  notesCell.style.overflow = 'hidden';
  notesCell.style.textOverflow = 'ellipsis';
  notesCell.style.whiteSpace = 'nowrap';
  notesCell.title = 'Click to edit';
  notesCell.addEventListener('click', (e) => {
    e.stopPropagation();
    startInlineEdit(notesCell, entry, 'notes');
  });
  tr.appendChild(notesCell);

  // Actions cell
  const actionsCell = document.createElement('td');
  const actionsDiv = document.createElement('div');
  actionsDiv.className = 'row-actions';

  const deleteBtn = document.createElement('button');
  deleteBtn.type = 'button';
  deleteBtn.className = 'action-btn';
  deleteBtn.style.color = '#D94B4B';
  deleteBtn.innerHTML = '<i data-lucide="trash-2" class="w-3.5 h-3.5"></i>';
  deleteBtn.setAttribute('aria-label', 'Delete entry');
  deleteBtn.addEventListener('click', (e) => {
    e.stopPropagation();
    deleteEntry(entry);
  });

  actionsDiv.appendChild(deleteBtn);
  actionsCell.appendChild(actionsDiv);
  tr.appendChild(actionsCell);

  return tr;
}

function renderView(): void {
  if (!containerRef) return;
  containerRef.innerHTML = '';

  const allEntries = timeEntryService.getAll();
  const employees = employeeService.getActive();
  const projects = projectService.getSelectable();

  // Header
  const header = document.createElement('div');
  header.style.cssText = 'display: flex; align-items: center; justify-content: space-between; margin-bottom: 24px; flex-wrap: wrap; gap: 12px;';

  const titleWrap = document.createElement('div');
  const h1 = document.createElement('h1');
  h1.className = 'text-content font-heading';
  h1.textContent = 'Timesheet';
  const subtitle = document.createElement('p');
  subtitle.className = 'text-content-secondary text-sm';
  subtitle.style.marginTop = '4px';
  subtitle.textContent = 'View and manage time entries';
  titleWrap.appendChild(h1);
  titleWrap.appendChild(subtitle);

  const addBtn = document.createElement('button');
  addBtn.type = 'button';
  addBtn.className = 'btn-primary';
  addBtn.setAttribute('data-action', 'add-entry');
  addBtn.innerHTML = '<i data-lucide="plus" class="w-4 h-4 mr-1.5 inline-block"></i>Add Entry';
  addBtn.addEventListener('click', () => openAddEntryModal());

  header.appendChild(titleWrap);
  header.appendChild(addBtn);
  containerRef.appendChild(header);

  // Filter bar
  const filterBar = document.createElement('div');
  filterBar.style.cssText = 'display: flex; gap: 12px; margin-bottom: 20px; flex-wrap: wrap; align-items: center;';

  // Date range picker
  const datePickerWrapper = document.createElement('div');
  datePickerWrapper.style.cssText = 'flex: 1; min-width: 300px;';
  const datePicker = createDateRangePicker({
    initialPreset: currentDateRange.preset,
    initialStart: currentDateRange.start,
    initialEnd: currentDateRange.end,
    onChange: (range: DateRange) => {
      currentDateRange = range;
      renderView();
    },
  });
  datePickerWrapper.appendChild(datePicker);
  filterBar.appendChild(datePickerWrapper);

  // Employee filter
  const empSelect = document.createElement('select');
  empSelect.className = 'select';
  empSelect.style.minWidth = '140px';
  const empAll = document.createElement('option');
  empAll.value = 'all';
  empAll.textContent = 'All Employees';
  empSelect.appendChild(empAll);
  for (const emp of employees) {
    const opt = document.createElement('option');
    opt.value = emp.id;
    opt.textContent = emp.name;
    if (emp.id === employeeFilter) opt.selected = true;
    empSelect.appendChild(opt);
  }
  empSelect.addEventListener('change', () => {
    employeeFilter = empSelect.value;
    renderView();
  });
  filterBar.appendChild(empSelect);

  // Project filter
  const projSelect = document.createElement('select');
  projSelect.className = 'select';
  projSelect.style.minWidth = '140px';
  const projAll = document.createElement('option');
  projAll.value = 'all';
  projAll.textContent = 'All Projects';
  projSelect.appendChild(projAll);
  for (const proj of projects) {
    const opt = document.createElement('option');
    opt.value = proj.id;
    opt.textContent = proj.name;
    if (proj.id === projectFilter) opt.selected = true;
    projSelect.appendChild(opt);
  }
  projSelect.addEventListener('change', () => {
    projectFilter = projSelect.value;
    renderView();
  });
  filterBar.appendChild(projSelect);

  containerRef.appendChild(filterBar);

  // Get filtered entries
  const filteredEntries = getFilteredEntries();

  // Empty state
  if (allEntries.length === 0) {
    const empty = document.createElement('div');
    empty.className = 'empty-state';
    empty.innerHTML = `
      <div class="empty-icon">
        <i data-lucide="clock" class="w-12 h-12 mx-auto text-content-disabled"></i>
      </div>
      <div class="empty-title">No time entries yet</div>
      <p class="empty-description">Use the timer or add a manual entry to get started.</p>
    `;
    const ctaBtn = document.createElement('button');
    ctaBtn.type = 'button';
    ctaBtn.className = 'btn-accent';
    ctaBtn.innerHTML = '<i data-lucide="plus" class="w-4 h-4 mr-1.5 inline-block"></i>Add Entry';
    ctaBtn.addEventListener('click', () => openAddEntryModal());
    empty.appendChild(ctaBtn);
    containerRef.appendChild(empty);
    try { createIcons(); } catch { /* ok */ }
    return;
  }

  if (filteredEntries.length === 0) {
    const noResults = document.createElement('div');
    noResults.className = 'empty-state';
    noResults.innerHTML = `
      <div class="empty-icon">
        <i data-lucide="search" class="w-12 h-12 mx-auto text-content-disabled"></i>
      </div>
      <div class="empty-title">No entries match your filters</div>
      <p class="empty-description">Try a different date range or remove filters.</p>
    `;
    containerRef.appendChild(noResults);
    try { createIcons(); } catch { /* ok */ }
    return;
  }

  // Table
  const tableWrapper = document.createElement('div');
  tableWrapper.className = 'card table-scroll';
  tableWrapper.style.padding = '0';

  const table = document.createElement('table');
  table.className = 'data-table';

  // Thead
  const thead = document.createElement('thead');
  const headerRow = document.createElement('tr');

  headerRow.appendChild(buildSortHeader('Date', 'date'));
  headerRow.appendChild(buildSortHeader('Employee', 'employee'));
  headerRow.appendChild(buildSortHeader('Project', 'project'));

  const timeTh = document.createElement('th');
  timeTh.setAttribute('scope', 'col');
  timeTh.textContent = 'Time';
  headerRow.appendChild(timeTh);

  headerRow.appendChild(buildSortHeader('Duration', 'duration'));

  const notesTh = document.createElement('th');
  notesTh.setAttribute('scope', 'col');
  notesTh.textContent = 'Notes';
  headerRow.appendChild(notesTh);

  const actionsTh = document.createElement('th');
  actionsTh.setAttribute('scope', 'col');
  actionsTh.textContent = '';
  actionsTh.style.width = '60px';
  headerRow.appendChild(actionsTh);

  thead.appendChild(headerRow);
  table.appendChild(thead);

  // Tbody
  const tbody = document.createElement('tbody');
  for (const entry of filteredEntries) {
    tbody.appendChild(buildEntryRow(entry));
  }
  table.appendChild(tbody);

  tableWrapper.appendChild(table);
  containerRef.appendChild(tableWrapper);

  // Summary
  const totalMinutes = filteredEntries.reduce((sum, e) => {
    if (e.duration !== null) return sum + e.duration;
    if (e.clockOut === null) return sum + calcDuration(e.clockIn, new Date().toISOString());
    return sum;
  }, 0);

  const summary = document.createElement('div');
  summary.style.cssText = 'margin-top: 16px; display: flex; gap: 16px; flex-wrap: wrap;';

  const entriesCount = document.createElement('span');
  entriesCount.className = 'text-content-secondary text-sm';
  entriesCount.textContent = `${filteredEntries.length} ${filteredEntries.length === 1 ? 'entry' : 'entries'}`;

  const totalHours = document.createElement('span');
  totalHours.className = 'text-content-secondary text-sm';
  totalHours.textContent = `Total: ${formatDuration(totalMinutes)}`;

  summary.appendChild(entriesCount);
  summary.appendChild(totalHours);
  containerRef.appendChild(summary);

  try { createIcons(); } catch { /* ok */ }
}

function onDataChange(): void {
  renderView();
}

export const render: ViewModule['render'] = (container) => {
  containerRef = container;
  employeeFilter = 'all';
  projectFilter = 'all';
  sortColumn = 'date';
  sortDirection = 'desc';

  const range = computePresetRange('thisWeek');
  currentDateRange = { ...range, preset: 'thisWeek' as PresetKey };

  eventBus.on('timeEntries:changed', onDataChange);
  eventBus.on('timer:changed', onDataChange);
  eventBus.on('employees:changed', onDataChange);
  eventBus.on('projects:changed', onDataChange);

  renderView();
};

export const destroy = (): void => {
  eventBus.off('timeEntries:changed', onDataChange);
  eventBus.off('timer:changed', onDataChange);
  eventBus.off('employees:changed', onDataChange);
  eventBus.off('projects:changed', onDataChange);
  containerRef = null;
};
