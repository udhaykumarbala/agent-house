/**
 * Reports View
 *
 * - Report type selector: Daily, Weekly, Monthly, By Project
 * - Filters: DateRangePicker, employee dropdown, project dropdown
 * - Summary cards: total hours, avg hours/day, top employee, top project
 * - CSS horizontal bar chart showing hours distribution
 * - Data table with sortable columns
 * - CSV export with descriptive filenames
 * - Print/PDF button triggering window.print()
 * - All 4 report type calculations
 */

import type { ViewModule, TimeEntry } from '../types';
import { timeEntryService } from '../data/timeEntryService';
import { employeeService } from '../data/employeeService';
import { projectService } from '../data/projectService';
import { eventBus } from '../utils/eventBus';
import { calcDuration } from '../utils/dates';
import { createDateRangePicker, computePresetRange } from '../components/dateRangePicker';
import type { DateRange, PresetKey } from '../components/dateRangePicker';
import { generateCSV, downloadCSV } from '../utils/export';
import { Toast } from '../components/toast';
import { createIcons } from 'lucide';

type ReportType = 'daily' | 'weekly' | 'monthly' | 'project';

// Chart colors from ui-spec.md data visualization palette
const CHART_COLORS = [
  '#CF8A2E', '#5746B2', '#E07B5F', '#2E9E8F',
  '#C75A8A', '#5B8AD4', '#8BA055', '#9B6DB0',
];

let containerRef: HTMLElement | null = null;
let currentDateRange: DateRange;
let reportType: ReportType = 'weekly';
let employeeFilter = 'all';
let projectFilter = 'all';
let tableSortColumn = '';
let tableSortDirection: 'asc' | 'desc' = 'desc';

const defaultRange = computePresetRange('thisWeek');
currentDateRange = { ...defaultRange, preset: 'thisWeek' as PresetKey };

interface ReportRow {
  [key: string]: string | number;
}

function getFilteredEntries(): TimeEntry[] {
  let entries = timeEntryService.getByDateRange(currentDateRange.start, currentDateRange.end);

  if (employeeFilter !== 'all') {
    entries = entries.filter(e => e.employeeId === employeeFilter);
  }

  if (projectFilter !== 'all') {
    entries = entries.filter(e => e.projectId === projectFilter);
  }

  return entries;
}

function getEntryDuration(entry: TimeEntry): number {
  if (entry.duration !== null) return entry.duration;
  if (entry.clockOut === null) return calcDuration(entry.clockIn, new Date().toISOString());
  return 0;
}

function computeSummary(entries: TimeEntry[]): { totalHours: number; avgPerDay: number; topEmployee: string; topProject: string } {
  const totalMinutes = entries.reduce((sum, e) => sum + getEntryDuration(e), 0);
  const totalHours = totalMinutes / 60;

  // Count unique days
  const uniqueDays = new Set(entries.map(e => e.date));
  const dayCount = Math.max(uniqueDays.size, 1);
  const avgPerDay = totalHours / dayCount;

  // Top employee
  const empHours: Record<string, number> = {};
  for (const e of entries) {
    empHours[e.employeeId] = (empHours[e.employeeId] || 0) + getEntryDuration(e);
  }
  const topEmpId = Object.entries(empHours).sort((a, b) => b[1] - a[1])[0]?.[0];
  const topEmployee = topEmpId ? (employeeService.getById(topEmpId)?.name || 'Unknown') : '—';

  // Top project
  const projHours: Record<string, number> = {};
  for (const e of entries) {
    projHours[e.projectId] = (projHours[e.projectId] || 0) + getEntryDuration(e);
  }
  const topProjId = Object.entries(projHours).sort((a, b) => b[1] - a[1])[0]?.[0];
  const topProject = topProjId ? (projectService.getById(topProjId)?.name || 'Unknown') : '—';

  return { totalHours, avgPerDay, topEmployee, topProject };
}

function computeChartData(entries: TimeEntry[]): { label: string; value: number; color: string }[] {
  const groups: Record<string, number> = {};

  if (reportType === 'daily' || reportType === 'weekly') {
    // Group by employee
    for (const e of entries) {
      const name = employeeService.getById(e.employeeId)?.name || 'Unknown';
      groups[name] = (groups[name] || 0) + getEntryDuration(e) / 60;
    }
  } else if (reportType === 'monthly') {
    // Group by project
    for (const e of entries) {
      const name = projectService.getById(e.projectId)?.name || 'Unknown';
      groups[name] = (groups[name] || 0) + getEntryDuration(e) / 60;
    }
  } else {
    // By Project: group by project
    for (const e of entries) {
      const name = projectService.getById(e.projectId)?.name || 'Unknown';
      groups[name] = (groups[name] || 0) + getEntryDuration(e) / 60;
    }
  }

  return Object.entries(groups)
    .sort((a, b) => b[1] - a[1])
    .map(([label, value], i) => ({
      label,
      value: Math.round(value * 10) / 10,
      color: CHART_COLORS[i % CHART_COLORS.length],
    }));
}

function computeTableData(entries: TimeEntry[]): { columns: string[]; rows: ReportRow[] } {
  switch (reportType) {
    case 'daily':
      return computeDailyTable(entries);
    case 'weekly':
      return computeWeeklyTable(entries);
    case 'monthly':
      return computeMonthlyTable(entries);
    case 'project':
      return computeProjectTable(entries);
    default:
      return { columns: [], rows: [] };
  }
}

function computeDailyTable(entries: TimeEntry[]): { columns: string[]; rows: ReportRow[] } {
  // Rows: employee + project breakdown for the selected date(s)
  const grouped: Record<string, Record<string, number>> = {};

  for (const e of entries) {
    const empName = employeeService.getById(e.employeeId)?.name || 'Unknown';
    const projName = projectService.getById(e.projectId)?.name || 'Unknown';
    if (!grouped[empName]) grouped[empName] = {};
    grouped[empName][projName] = (grouped[empName][projName] || 0) + getEntryDuration(e) / 60;
  }

  const rows: ReportRow[] = [];
  for (const [emp, projects] of Object.entries(grouped)) {
    let empTotal = 0;
    for (const [proj, hours] of Object.entries(projects)) {
      empTotal += hours;
      rows.push({ employee: emp, project: proj, hours: Math.round(hours * 10) / 10 });
    }
    rows.push({ employee: emp, project: 'TOTAL', hours: Math.round(empTotal * 10) / 10 });
  }

  return { columns: ['Employee', 'Project', 'Hours'], rows };
}

function computeWeeklyTable(entries: TimeEntry[]): { columns: string[]; rows: ReportRow[] } {
  // Rows: employee with daily breakdown
  const employees = new Set<string>();
  const dailyHours: Record<string, Record<string, number>> = {};

  for (const e of entries) {
    const empName = employeeService.getById(e.employeeId)?.name || 'Unknown';
    employees.add(empName);
    if (!dailyHours[empName]) dailyHours[empName] = {};
    dailyHours[empName][e.date] = (dailyHours[empName][e.date] || 0) + getEntryDuration(e) / 60;
  }

  // Get all unique dates sorted
  const allDates = Array.from(new Set(entries.map(e => e.date))).sort();
  const dayLabels = allDates.map(d => {
    const date = new Date(d + 'T00:00:00');
    return date.toLocaleDateString('en-US', { weekday: 'short', month: 'short', day: 'numeric' });
  });

  const columns = ['Employee', ...dayLabels, 'Total'];
  const rows: ReportRow[] = [];

  for (const emp of Array.from(employees).sort()) {
    const row: ReportRow = { employee: emp };
    let total = 0;
    for (let i = 0; i < allDates.length; i++) {
      const hours = dailyHours[emp]?.[allDates[i]] || 0;
      total += hours;
      row[dayLabels[i]] = Math.round(hours * 10) / 10;
    }
    row['total'] = Math.round(total * 10) / 10;
    rows.push(row);
  }

  return { columns, rows };
}

function computeMonthlyTable(entries: TimeEntry[]): { columns: string[]; rows: ReportRow[] } {
  // Show both employee totals and project totals
  const empHours: Record<string, number> = {};
  const projHours: Record<string, number> = {};

  for (const e of entries) {
    const empName = employeeService.getById(e.employeeId)?.name || 'Unknown';
    const projName = projectService.getById(e.projectId)?.name || 'Unknown';
    empHours[empName] = (empHours[empName] || 0) + getEntryDuration(e) / 60;
    projHours[projName] = (projHours[projName] || 0) + getEntryDuration(e) / 60;
  }

  const rows: ReportRow[] = [];

  // Employee rows
  for (const [name, hours] of Object.entries(empHours).sort((a, b) => b[1] - a[1])) {
    rows.push({ category: 'Employee', name, hours: Math.round(hours * 10) / 10 });
  }

  // Project rows
  for (const [name, hours] of Object.entries(projHours).sort((a, b) => b[1] - a[1])) {
    rows.push({ category: 'Project', name, hours: Math.round(hours * 10) / 10 });
  }

  return { columns: ['Category', 'Name', 'Hours'], rows };
}

function computeProjectTable(entries: TimeEntry[]): { columns: string[]; rows: ReportRow[] } {
  // Group by project, then employee + task breakdown
  const grouped: Record<string, { byEmployee: Record<string, number>; byTask: Record<string, number>; total: number }> = {};

  for (const e of entries) {
    const projName = projectService.getById(e.projectId)?.name || 'Unknown';
    const empName = employeeService.getById(e.employeeId)?.name || 'Unknown';
    const taskName = e.taskId ? (projectService.getTaskById(e.taskId)?.name || 'No Task') : 'No Task';

    if (!grouped[projName]) grouped[projName] = { byEmployee: {}, byTask: {}, total: 0 };
    const dur = getEntryDuration(e) / 60;
    grouped[projName].byEmployee[empName] = (grouped[projName].byEmployee[empName] || 0) + dur;
    grouped[projName].byTask[taskName] = (grouped[projName].byTask[taskName] || 0) + dur;
    grouped[projName].total += dur;
  }

  const rows: ReportRow[] = [];
  const sorted = Object.entries(grouped).sort((a, b) => b[1].total - a[1].total);

  for (const [projName, data] of sorted) {
    // Project total row
    rows.push({ project: projName, detail: 'Total', hours: Math.round(data.total * 10) / 10 });

    // Employee breakdown
    for (const [emp, hours] of Object.entries(data.byEmployee).sort((a, b) => b[1] - a[1])) {
      rows.push({ project: projName, detail: `Employee: ${emp}`, hours: Math.round(hours * 10) / 10 });
    }

    // Task breakdown
    for (const [task, hours] of Object.entries(data.byTask).sort((a, b) => b[1] - a[1])) {
      rows.push({ project: projName, detail: `Task: ${task}`, hours: Math.round(hours * 10) / 10 });
    }
  }

  return { columns: ['Project', 'Detail', 'Hours'], rows };
}

function handleTableSort(column: string): void {
  if (tableSortColumn === column) {
    tableSortDirection = tableSortDirection === 'asc' ? 'desc' : 'asc';
  } else {
    tableSortColumn = column;
    tableSortDirection = 'desc';
  }
  renderView();
}

function sortTableRows(rows: ReportRow[], columns: string[]): ReportRow[] {
  if (!tableSortColumn) return rows;

  // Map display column name to row key
  const keyMap: Record<string, string> = {};
  for (const col of columns) {
    keyMap[col] = col.toLowerCase().replace(/[^a-z0-9]/g, '');
  }

  // Find the actual key in the row data
  const sortKey = Object.keys(rows[0] || {}).find(k => {
    const normalized = k.toLowerCase().replace(/[^a-z0-9]/g, '');
    return normalized === tableSortColumn.toLowerCase().replace(/[^a-z0-9]/g, '');
  });

  if (!sortKey) return rows;

  const dir = tableSortDirection === 'asc' ? 1 : -1;
  return [...rows].sort((a, b) => {
    const valA = a[sortKey];
    const valB = b[sortKey];
    if (typeof valA === 'number' && typeof valB === 'number') return (valA - valB) * dir;
    return String(valA).localeCompare(String(valB)) * dir;
  });
}

function exportReport(entries: TimeEntry[]): void {
  const { columns, rows } = computeTableData(entries);
  if (rows.length === 0) {
    Toast.warning('No data to export');
    return;
  }

  // Build CSV data
  const keys = Object.keys(rows[0]);
  const headers = columns.length > 0 ? columns : keys;
  const csvRows = rows.map(row => keys.map(k => String(row[k] ?? '')));

  const csv = generateCSV(headers, csvRows);

  // Build descriptive filename
  const typeLabel = reportType;
  const dateRange = `${currentDateRange.start}-to-${currentDateRange.end}`;
  const filename = `${typeLabel}-report-${dateRange}.csv`;

  downloadCSV(csv, filename);
  Toast.success('Report exported as CSV');
}

function buildReportTypeSelector(): HTMLElement {
  const container = document.createElement('div');
  container.style.cssText = 'display: flex; gap: 4px; background: rgba(87, 70, 178, 0.08); border-radius: 8px; padding: 4px;';

  const types: { key: ReportType; label: string }[] = [
    { key: 'daily', label: 'Daily' },
    { key: 'weekly', label: 'Weekly' },
    { key: 'monthly', label: 'Monthly' },
    { key: 'project', label: 'By Project' },
  ];

  for (const t of types) {
    const btn = document.createElement('button');
    btn.type = 'button';
    btn.textContent = t.label;
    btn.style.cssText = `
      padding: 8px 16px; border-radius: 6px; font-size: 13px; font-weight: 500;
      border: none; cursor: pointer; transition: all 150ms ease;
      font-family: 'Inter', sans-serif;
    `;

    if (t.key === reportType) {
      btn.style.background = '#5746B2';
      btn.style.color = '#ECE9F5';
    } else {
      btn.style.background = 'transparent';
      btn.style.color = '#8C83A8';
    }

    btn.addEventListener('mouseenter', () => {
      if (t.key !== reportType) btn.style.color = '#ECE9F5';
    });
    btn.addEventListener('mouseleave', () => {
      if (t.key !== reportType) btn.style.color = '#8C83A8';
    });

    btn.addEventListener('click', () => {
      reportType = t.key;
      tableSortColumn = '';
      // Adjust date range preset based on report type
      if (t.key === 'daily') {
        const range = computePresetRange('today');
        currentDateRange = { ...range, preset: 'today' };
      } else if (t.key === 'weekly') {
        const range = computePresetRange('thisWeek');
        currentDateRange = { ...range, preset: 'thisWeek' };
      } else if (t.key === 'monthly') {
        const range = computePresetRange('thisMonth');
        currentDateRange = { ...range, preset: 'thisMonth' };
      }
      renderView();
    });

    container.appendChild(btn);
  }

  return container;
}

function buildSummaryCards(summary: { totalHours: number; avgPerDay: number; topEmployee: string; topProject: string }): HTMLElement {
  const grid = document.createElement('div');
  grid.style.cssText = 'display: grid; grid-template-columns: repeat(auto-fit, minmax(180px, 1fr)); gap: 16px; margin-bottom: 24px;';

  const cards: { label: string; value: string; icon: string }[] = [
    { label: 'TOTAL HOURS', value: `${Math.round(summary.totalHours * 10) / 10}h`, icon: 'clock' },
    { label: 'AVG HOURS/DAY', value: `${Math.round(summary.avgPerDay * 10) / 10}h`, icon: 'trending-up' },
    { label: 'TOP EMPLOYEE', value: summary.topEmployee, icon: 'user' },
    { label: 'TOP PROJECT', value: summary.topProject, icon: 'folder' },
  ];

  for (const card of cards) {
    const el = document.createElement('div');
    el.className = 'stat-card';

    const iconRow = document.createElement('div');
    iconRow.style.cssText = 'display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px;';

    const iconEl = document.createElement('i');
    iconEl.setAttribute('data-lucide', card.icon);
    iconEl.className = 'stat-icon';
    iconRow.appendChild(iconEl);

    const value = document.createElement('div');
    value.className = 'stat-value';
    value.textContent = card.value;

    const label = document.createElement('div');
    label.className = 'stat-label';
    label.textContent = card.label;

    el.appendChild(iconRow);
    el.appendChild(value);
    el.appendChild(label);
    grid.appendChild(el);
  }

  return grid;
}

function buildBarChart(chartData: { label: string; value: number; color: string }[]): HTMLElement {
  const wrapper = document.createElement('div');
  wrapper.className = 'card';
  wrapper.style.marginBottom = '24px';

  const title = document.createElement('h3');
  title.style.cssText = 'margin-bottom: 16px; font-size: 16px; font-weight: 600;';
  title.textContent = reportType === 'daily' || reportType === 'weekly' ? 'Hours by Employee' : 'Hours by Project';
  wrapper.appendChild(title);

  if (chartData.length === 0) {
    const empty = document.createElement('p');
    empty.className = 'text-content-secondary';
    empty.style.cssText = 'font-size: 14px; text-align: center; padding: 24px 0;';
    empty.textContent = 'No data to display';
    wrapper.appendChild(empty);
    return wrapper;
  }

  const maxValue = Math.max(...chartData.map(d => d.value), 1);

  const chartContainer = document.createElement('div');
  chartContainer.style.cssText = 'display: flex; flex-direction: column; gap: 10px;';

  for (const item of chartData) {
    const row = document.createElement('div');
    row.style.cssText = 'display: flex; align-items: center; gap: 12px;';

    const labelEl = document.createElement('span');
    labelEl.style.cssText = 'width: 120px; font-size: 13px; color: #ECE9F5; text-overflow: ellipsis; overflow: hidden; white-space: nowrap; flex-shrink: 0; text-align: right;';
    labelEl.textContent = item.label;
    labelEl.title = item.label;

    const barWrapper = document.createElement('div');
    barWrapper.style.cssText = 'flex: 1; height: 28px; background: rgba(87, 70, 178, 0.08); border-radius: 4px; overflow: hidden;';

    const bar = document.createElement('div');
    const percentage = Math.max((item.value / maxValue) * 100, 2);
    bar.style.cssText = `width: ${percentage}%; height: 100%; background: ${item.color}; opacity: 0.85; border-radius: 4px; transition: width 300ms ease;`;
    barWrapper.appendChild(bar);

    const valueEl = document.createElement('span');
    valueEl.style.cssText = 'width: 60px; font-size: 13px; color: #8C83A8; font-variant-numeric: tabular-nums; flex-shrink: 0;';
    valueEl.textContent = `${item.value}h`;

    row.appendChild(labelEl);
    row.appendChild(barWrapper);
    row.appendChild(valueEl);
    chartContainer.appendChild(row);
  }

  wrapper.appendChild(chartContainer);
  return wrapper;
}

function buildDataTable(columns: string[], rows: ReportRow[]): HTMLElement {
  const wrapper = document.createElement('div');
  wrapper.className = 'card table-scroll';
  wrapper.style.padding = '0';

  if (rows.length === 0) {
    const empty = document.createElement('div');
    empty.className = 'empty-state';
    empty.style.padding = '24px';
    empty.innerHTML = '<p class="text-content-secondary">No data available for the selected filters.</p>';
    wrapper.appendChild(empty);
    return wrapper;
  }

  const table = document.createElement('table');
  table.className = 'data-table';

  // Thead
  const thead = document.createElement('thead');
  const headerRow = document.createElement('tr');

  for (const col of columns) {
    const th = document.createElement('th');
    th.setAttribute('scope', 'col');
    th.style.cursor = 'pointer';

    const span = document.createElement('span');
    span.textContent = col;
    th.appendChild(span);

    const colKey = col.toLowerCase().replace(/[^a-z0-9]/g, '');
    const indicator = document.createElement('span');
    indicator.className = 'sort-indicator';
    if (tableSortColumn.toLowerCase().replace(/[^a-z0-9]/g, '') === colKey) {
      indicator.classList.add('active');
      indicator.textContent = tableSortDirection === 'asc' ? ' \u25B2' : ' \u25BC';
    } else {
      indicator.textContent = ' \u25B4';
    }
    th.appendChild(indicator);
    th.addEventListener('click', () => handleTableSort(col));
    headerRow.appendChild(th);
  }

  thead.appendChild(headerRow);
  table.appendChild(thead);

  // Sort rows
  const sortedRows = sortTableRows(rows, columns);

  // Tbody
  const tbody = document.createElement('tbody');
  const keys = Object.keys(rows[0]);

  for (const row of sortedRows) {
    const tr = document.createElement('tr');

    // Highlight total rows
    const isTotalRow = Object.values(row).some(v => typeof v === 'string' && v === 'TOTAL');
    if (isTotalRow) {
      tr.style.fontWeight = '600';
      tr.style.background = 'rgba(87, 70, 178, 0.05)';
    }

    for (const key of keys) {
      const td = document.createElement('td');
      const val = row[key];

      if (typeof val === 'number') {
        td.textContent = `${val}h`;
        td.style.fontVariantNumeric = 'tabular-nums';
      } else {
        td.textContent = String(val);
      }

      tr.appendChild(td);
    }

    tbody.appendChild(tr);
  }

  table.appendChild(tbody);
  wrapper.appendChild(table);
  return wrapper;
}

function renderView(): void {
  if (!containerRef) return;
  containerRef.innerHTML = '';

  const entries = getFilteredEntries();
  const employees = employeeService.getActive();
  const projects = projectService.getSelectable();

  // Header with export buttons
  const header = document.createElement('div');
  header.style.cssText = 'display: flex; align-items: center; justify-content: space-between; margin-bottom: 24px; flex-wrap: wrap; gap: 12px;';

  const titleWrap = document.createElement('div');
  const h1 = document.createElement('h1');
  h1.className = 'text-content font-heading';
  h1.textContent = 'Reports';
  const subtitle = document.createElement('p');
  subtitle.className = 'text-content-secondary text-sm';
  subtitle.style.marginTop = '4px';
  subtitle.textContent = 'View and export time reports';
  titleWrap.appendChild(h1);
  titleWrap.appendChild(subtitle);

  const btnGroup = document.createElement('div');
  btnGroup.style.cssText = 'display: flex; gap: 8px;';

  const printBtn = document.createElement('button');
  printBtn.type = 'button';
  printBtn.className = 'btn-secondary';
  printBtn.innerHTML = '<i data-lucide="printer" class="w-4 h-4 mr-1.5 inline-block"></i>Print / PDF';
  printBtn.addEventListener('click', () => {
    window.print();
  });

  const exportBtn = document.createElement('button');
  exportBtn.type = 'button';
  exportBtn.className = 'btn-accent';
  exportBtn.id = 'export-csv-btn';
  exportBtn.innerHTML = '<i data-lucide="download" class="w-4 h-4 mr-1.5 inline-block"></i>Export CSV';
  exportBtn.addEventListener('click', () => {
    exportReport(entries);
  });

  btnGroup.appendChild(printBtn);
  btnGroup.appendChild(exportBtn);

  header.appendChild(titleWrap);
  header.appendChild(btnGroup);
  containerRef.appendChild(header);

  // Report type selector
  const typeSelector = buildReportTypeSelector();
  typeSelector.style.marginBottom = '20px';
  containerRef.appendChild(typeSelector);

  // Filter bar
  const filterBar = document.createElement('div');
  filterBar.style.cssText = 'display: flex; gap: 12px; margin-bottom: 24px; flex-wrap: wrap; align-items: center;';

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

  // Empty state if no entries at all
  if (timeEntryService.getAll().length === 0) {
    const empty = document.createElement('div');
    empty.className = 'empty-state';
    empty.innerHTML = `
      <div class="empty-icon">
        <i data-lucide="bar-chart-3" class="w-12 h-12 mx-auto text-content-disabled"></i>
      </div>
      <div class="empty-title">No reports yet</div>
      <p class="empty-description">Track some time to see reports and analytics here.</p>
    `;
    containerRef.appendChild(empty);
    try { createIcons(); } catch { /* ok */ }
    return;
  }

  // No matching entries for filters
  if (entries.length === 0) {
    const noResults = document.createElement('div');
    noResults.className = 'empty-state';
    noResults.innerHTML = `
      <div class="empty-icon">
        <i data-lucide="search" class="w-12 h-12 mx-auto text-content-disabled"></i>
      </div>
      <div class="empty-title">No time entries match your filters</div>
      <p class="empty-description">Try a different date range or remove filters.</p>
    `;
    containerRef.appendChild(noResults);
    try { createIcons(); } catch { /* ok */ }
    return;
  }

  // Summary cards
  const summary = computeSummary(entries);
  containerRef.appendChild(buildSummaryCards(summary));

  // Bar chart
  const chartData = computeChartData(entries);
  containerRef.appendChild(buildBarChart(chartData));

  // Data table
  const { columns, rows } = computeTableData(entries);
  const tableTitle = document.createElement('h3');
  tableTitle.style.cssText = 'margin-bottom: 12px; font-size: 16px; font-weight: 600;';
  tableTitle.textContent = 'Detailed Breakdown';
  containerRef.appendChild(tableTitle);
  containerRef.appendChild(buildDataTable(columns, rows));

  try { createIcons(); } catch { /* ok */ }
}

function onDataChange(): void {
  renderView();
}

export const render: ViewModule['render'] = (container) => {
  containerRef = container;
  employeeFilter = 'all';
  projectFilter = 'all';
  reportType = 'weekly';
  tableSortColumn = '';
  tableSortDirection = 'desc';

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
