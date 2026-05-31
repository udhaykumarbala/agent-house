/**
 * Employee Management View
 *
 * - Employee list with color-coded status indicators (green/gray dots)
 * - Add Employee modal with validation (name required, max 100 chars)
 * - Edit employee via modal, pre-filled with current data
 * - Deactivate/activate toggle via row actions
 * - Search by name with 300ms debounce
 * - Filter by department and status dropdowns
 * - Empty state with working "Add your first employee" CTA
 * - Today's hours shown per employee
 * - EventBus listeners for employees:changed and timeEntries:changed
 */

import type { ViewModule, Employee } from '../types';
import { employeeService } from '../data/employeeService';
import { timeEntryService } from '../data/timeEntryService';
import { eventBus } from '../utils/eventBus';
import { formatDuration } from '../utils/dates';
import { Modal } from '../components/modal';
import { Toast } from '../components/toast';
import { createIcons } from 'lucide';

const EMPLOYEE_COLORS = [
  '#CF8A2E', '#5746B2', '#E07B5F', '#2E9E8F',
  '#C75A8A', '#5B8AD4', '#8BA055', '#9B6DB0',
];

let containerRef: HTMLElement | null = null;
let searchQuery = '';
let statusFilter: 'all' | 'active' | 'inactive' = 'all';
let departmentFilter = 'all';
let debounceTimer: ReturnType<typeof setTimeout> | null = null;

function getFilteredEmployees(): Employee[] {
  let employees = employeeService.getAll();

  if (statusFilter !== 'all') {
    employees = employees.filter(e => e.status === statusFilter);
  }

  if (departmentFilter !== 'all') {
    employees = employees.filter(e => e.department === departmentFilter);
  }

  if (searchQuery.trim()) {
    const q = searchQuery.trim().toLowerCase();
    employees = employees.filter(e => e.name.toLowerCase().includes(q));
  }

  return employees;
}

function getAllDepartments(): string[] {
  const all = employeeService.getAll();
  const depts = new Set<string>();
  for (const e of all) {
    if (e.department) depts.add(e.department);
  }
  return Array.from(depts).sort();
}

function buildColorPicker(selectedColor: string): HTMLElement {
  const wrapper = document.createElement('div');
  wrapper.style.cssText = 'display: flex; gap: 8px; flex-wrap: wrap;';

  for (const color of EMPLOYEE_COLORS) {
    const swatch = document.createElement('button');
    swatch.type = 'button';
    swatch.style.cssText = `
      width: 32px; height: 32px; border-radius: 9999px;
      background: ${color}; border: 2px solid transparent;
      cursor: pointer; transition: border-color 150ms ease, transform 150ms ease;
    `;

    if (color === selectedColor) {
      swatch.style.borderColor = '#ECE9F5';
      swatch.style.transform = 'scale(1.15)';
    }

    swatch.setAttribute('data-color', color);
    swatch.setAttribute('aria-label', `Select color ${color}`);

    swatch.addEventListener('click', () => {
      wrapper.querySelectorAll('button').forEach(btn => {
        btn.style.borderColor = 'transparent';
        btn.style.transform = 'scale(1)';
      });
      swatch.style.borderColor = '#ECE9F5';
      swatch.style.transform = 'scale(1.15)';
    });

    wrapper.appendChild(swatch);
  }

  return wrapper;
}

function getSelectedColor(colorPicker: HTMLElement): string {
  const selected = colorPicker.querySelector('button[style*="border-color: rgb(236, 233, 245)"]') as HTMLElement | null;
  return selected?.getAttribute('data-color') || EMPLOYEE_COLORS[0];
}

function openEmployeeModal(employee?: Employee): void {
  const isEdit = !!employee;
  const content = document.createElement('div');
  content.innerHTML = `
    <div style="display: flex; flex-direction: column; gap: 16px;">
      <div>
        <label class="form-label">Name *</label>
        <input id="emp-name" type="text" class="input" maxlength="100" placeholder="e.g. John Smith" value="" />
        <div id="name-error" class="form-error" style="display: none;"></div>
      </div>
      <div>
        <label class="form-label">Role</label>
        <input id="emp-role" type="text" class="input" placeholder="e.g. Developer (optional)" value="" />
      </div>
      <div>
        <label class="form-label">Department</label>
        <input id="emp-department" type="text" class="input" placeholder="e.g. Engineering (optional)" value="" />
      </div>
      <div>
        <label class="form-label">Hourly Rate</label>
        <input id="emp-rate" type="number" class="input" placeholder="e.g. 50 (optional)" min="0" step="0.01" value="" />
      </div>
      <div>
        <label class="form-label">Color</label>
        <div id="color-picker-slot"></div>
      </div>
    </div>
  `;

  const nameInput = content.querySelector('#emp-name') as HTMLInputElement;
  const roleInput = content.querySelector('#emp-role') as HTMLInputElement;
  const deptInput = content.querySelector('#emp-department') as HTMLInputElement;
  const rateInput = content.querySelector('#emp-rate') as HTMLInputElement;
  const nameError = content.querySelector('#name-error') as HTMLElement;
  const colorSlot = content.querySelector('#color-picker-slot') as HTMLElement;

  const defaultColor = employee?.color || EMPLOYEE_COLORS[0];
  const colorPicker = buildColorPicker(defaultColor);
  colorSlot.appendChild(colorPicker);

  if (isEdit) {
    nameInput.value = employee!.name;
    roleInput.value = employee!.role;
    deptInput.value = employee!.department;
    rateInput.value = employee!.hourlyRate !== null ? String(employee!.hourlyRate) : '';
  }

  Modal.open({
    title: isEdit ? 'Edit Employee' : 'Add Employee',
    content,
    submitLabel: isEdit ? 'Save' : 'Add',
    onSubmit: () => {
      const name = nameInput.value.trim();
      if (!name) {
        nameError.textContent = 'Name is required';
        nameError.style.display = 'block';
        nameInput.classList.add('error');
        return;
      }
      if (name.length > 100) {
        nameError.textContent = 'Name must be 100 characters or less';
        nameError.style.display = 'block';
        nameInput.classList.add('error');
        return;
      }

      const color = getSelectedColor(colorPicker);
      const rate = rateInput.value ? parseFloat(rateInput.value) : null;

      try {
        if (isEdit) {
          employeeService.update(employee!.id, {
            name,
            role: roleInput.value.trim(),
            department: deptInput.value.trim(),
            hourlyRate: rate,
            color,
          });
          Toast.success('Employee updated');
        } else {
          employeeService.create({
            name,
            role: roleInput.value.trim(),
            department: deptInput.value.trim(),
            hourlyRate: rate,
          });
          Toast.success('Employee added');
        }
        Modal.close();
        renderView();
      } catch (err) {
        nameError.textContent = (err as Error).message;
        nameError.style.display = 'block';
      }
    },
  });
}

function toggleEmployeeStatus(employee: Employee): void {
  const newStatus = employee.status === 'active' ? 'inactive' : 'active';
  try {
    employeeService.update(employee.id, { status: newStatus });
    Toast.success(`${employee.name} ${newStatus === 'active' ? 'activated' : 'deactivated'}`);
    renderView();
  } catch (err) {
    Toast.error((err as Error).message);
  }
}

function buildEmployeeRow(employee: Employee): HTMLElement {
  const todayHours = timeEntryService.getTodayHoursForEmployee(employee.id);
  const isClockedIn = timeEntryService.isEmployeeClockedIn(employee.id);

  const row = document.createElement('tr');
  row.style.cursor = 'pointer';
  row.addEventListener('click', (e) => {
    if ((e.target as HTMLElement).closest('.row-actions')) return;
    openEmployeeModal(employee);
  });

  // Status + Name cell
  const nameCell = document.createElement('td');

  const nameCellInner = document.createElement('div');
  nameCellInner.style.cssText = 'display: flex; align-items: center; gap: 10px;';

  const dot = document.createElement('span');
  if (employee.status === 'inactive') {
    dot.className = 'status-dot inactive';
    dot.setAttribute('aria-label', 'Status: Inactive');
  } else if (isClockedIn) {
    dot.className = 'status-dot active';
    dot.setAttribute('aria-label', 'Status: Clocked In');
  } else {
    dot.className = 'status-dot inactive';
    dot.setAttribute('aria-label', 'Status: Offline');
  }

  const colorBar = document.createElement('span');
  colorBar.style.cssText = `width: 3px; height: 24px; border-radius: 2px; background: ${employee.color}; flex-shrink: 0;`;

  const nameWrap = document.createElement('div');
  nameWrap.style.cssText = 'min-width: 0;';

  const nameEl = document.createElement('div');
  nameEl.style.cssText = 'font-weight: 500; color: #ECE9F5; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;';
  nameEl.textContent = employee.name;
  if (employee.status === 'inactive') nameEl.style.opacity = '0.5';

  const roleEl = document.createElement('div');
  roleEl.style.cssText = 'font-size: 12px; color: #8C83A8;';
  roleEl.textContent = employee.role || '';

  nameWrap.appendChild(nameEl);
  if (employee.role) nameWrap.appendChild(roleEl);

  nameCellInner.appendChild(dot);
  nameCellInner.appendChild(colorBar);
  nameCellInner.appendChild(nameWrap);
  nameCell.appendChild(nameCellInner);

  // Department cell
  const deptCell = document.createElement('td');
  deptCell.textContent = employee.department || '—';
  deptCell.style.color = employee.department ? '#ECE9F5' : '#564F6D';

  // Status cell
  const statusCell = document.createElement('td');
  const badge = document.createElement('span');
  if (employee.status === 'inactive') {
    badge.className = 'badge-inactive';
    badge.textContent = 'Inactive';
  } else if (isClockedIn) {
    badge.className = 'badge-active';
    badge.textContent = 'Clocked In';
  } else {
    badge.className = 'badge-inactive';
    badge.textContent = 'Offline';
  }
  statusCell.appendChild(badge);

  // Today's hours cell
  const hoursCell = document.createElement('td');
  hoursCell.style.cssText = 'font-variant-numeric: tabular-nums;';
  hoursCell.textContent = formatDuration(todayHours);

  // Rate cell
  const rateCell = document.createElement('td');
  rateCell.textContent = employee.hourlyRate !== null ? `$${employee.hourlyRate.toFixed(2)}/hr` : '—';
  rateCell.style.color = employee.hourlyRate !== null ? '#ECE9F5' : '#564F6D';

  // Actions cell
  const actionsCell = document.createElement('td');
  const actions = document.createElement('div');
  actions.className = 'row-actions';

  const editBtn = document.createElement('button');
  editBtn.type = 'button';
  editBtn.className = 'action-btn';
  editBtn.innerHTML = '<i data-lucide="pencil" class="w-3.5 h-3.5"></i>';
  editBtn.setAttribute('aria-label', `Edit ${employee.name}`);
  editBtn.addEventListener('click', (e) => {
    e.stopPropagation();
    openEmployeeModal(employee);
  });

  const toggleBtn = document.createElement('button');
  toggleBtn.type = 'button';
  toggleBtn.className = 'action-btn';
  toggleBtn.innerHTML = employee.status === 'active'
    ? '<i data-lucide="user-x" class="w-3.5 h-3.5"></i>'
    : '<i data-lucide="user-check" class="w-3.5 h-3.5"></i>';
  toggleBtn.setAttribute('aria-label', employee.status === 'active' ? `Deactivate ${employee.name}` : `Activate ${employee.name}`);
  toggleBtn.addEventListener('click', (e) => {
    e.stopPropagation();
    toggleEmployeeStatus(employee);
  });

  actions.appendChild(editBtn);
  actions.appendChild(toggleBtn);
  actionsCell.appendChild(actions);

  row.appendChild(nameCell);
  row.appendChild(deptCell);
  row.appendChild(statusCell);
  row.appendChild(hoursCell);
  row.appendChild(rateCell);
  row.appendChild(actionsCell);

  return row;
}

function renderView(): void {
  if (!containerRef) return;
  containerRef.innerHTML = '';

  const allEmployees = employeeService.getAll();
  const hasEmployees = allEmployees.length > 0;

  // Header
  const header = document.createElement('div');
  header.style.cssText = 'display: flex; align-items: center; justify-content: space-between; margin-bottom: 24px; flex-wrap: wrap; gap: 12px;';

  const titleWrap = document.createElement('div');
  const h1 = document.createElement('h1');
  h1.className = 'text-content font-heading';
  h1.textContent = 'Team';
  const subtitle = document.createElement('p');
  subtitle.className = 'text-content-secondary text-sm';
  subtitle.style.marginTop = '4px';
  subtitle.textContent = 'Manage your employees';
  titleWrap.appendChild(h1);
  titleWrap.appendChild(subtitle);

  const addBtn = document.createElement('button');
  addBtn.type = 'button';
  addBtn.className = 'btn-primary';
  addBtn.innerHTML = '<i data-lucide="plus" class="w-4 h-4 mr-1.5 inline-block"></i>Add Employee';
  addBtn.addEventListener('click', () => openEmployeeModal());

  header.appendChild(titleWrap);
  header.appendChild(addBtn);
  containerRef.appendChild(header);

  // Empty state
  if (!hasEmployees) {
    const empty = document.createElement('div');
    empty.className = 'empty-state';
    empty.innerHTML = `
      <div class="empty-icon">
        <i data-lucide="users" class="w-12 h-12 mx-auto text-content-disabled"></i>
      </div>
      <div class="empty-title">Add your first employee</div>
      <p class="empty-description">Add employees to start tracking their time.</p>
    `;
    const ctaBtn = document.createElement('button');
    ctaBtn.type = 'button';
    ctaBtn.className = 'btn-accent';
    ctaBtn.innerHTML = '<i data-lucide="plus" class="w-4 h-4 mr-1.5 inline-block"></i>Add Employee';
    ctaBtn.addEventListener('click', () => openEmployeeModal());
    empty.appendChild(ctaBtn);
    containerRef.appendChild(empty);
    try { createIcons(); } catch { /* ok */ }
    return;
  }

  // Search + Filters bar
  const filterBar = document.createElement('div');
  filterBar.style.cssText = 'display: flex; gap: 12px; margin-bottom: 20px; flex-wrap: wrap; align-items: center;';

  // Search input
  const searchWrap = document.createElement('div');
  searchWrap.style.cssText = 'flex: 1; min-width: 200px; position: relative;';
  const searchInput = document.createElement('input');
  searchInput.type = 'text';
  searchInput.className = 'input';
  searchInput.placeholder = 'Search by name...';
  searchInput.value = searchQuery;
  searchInput.style.paddingLeft = '36px';

  const searchIcon = document.createElement('span');
  searchIcon.style.cssText = 'position: absolute; left: 12px; top: 50%; transform: translateY(-50%); pointer-events: none; color: #564F6D;';
  searchIcon.innerHTML = '<i data-lucide="search" class="w-4 h-4"></i>';

  searchInput.addEventListener('input', () => {
    if (debounceTimer) clearTimeout(debounceTimer);
    debounceTimer = setTimeout(() => {
      searchQuery = searchInput.value;
      renderView();
    }, 300);
  });

  searchWrap.appendChild(searchIcon);
  searchWrap.appendChild(searchInput);

  // Department filter
  const deptSelect = document.createElement('select');
  deptSelect.className = 'select';
  deptSelect.style.minWidth = '140px';

  const deptAll = document.createElement('option');
  deptAll.value = 'all';
  deptAll.textContent = 'All Departments';
  deptSelect.appendChild(deptAll);

  const departments = getAllDepartments();
  for (const dept of departments) {
    const opt = document.createElement('option');
    opt.value = dept;
    opt.textContent = dept;
    if (dept === departmentFilter) opt.selected = true;
    deptSelect.appendChild(opt);
  }

  deptSelect.addEventListener('change', () => {
    departmentFilter = deptSelect.value;
    renderView();
  });

  // Status filter
  const statusSelect = document.createElement('select');
  statusSelect.className = 'select';
  statusSelect.style.minWidth = '120px';

  const statusOpts: { value: string; label: string }[] = [
    { value: 'all', label: 'All Status' },
    { value: 'active', label: 'Active' },
    { value: 'inactive', label: 'Inactive' },
  ];

  for (const opt of statusOpts) {
    const option = document.createElement('option');
    option.value = opt.value;
    option.textContent = opt.label;
    if (opt.value === statusFilter) option.selected = true;
    statusSelect.appendChild(option);
  }

  statusSelect.addEventListener('change', () => {
    statusFilter = statusSelect.value as typeof statusFilter;
    renderView();
  });

  filterBar.appendChild(searchWrap);
  filterBar.appendChild(deptSelect);
  filterBar.appendChild(statusSelect);
  containerRef.appendChild(filterBar);

  // Filtered employees
  const filtered = getFilteredEmployees();

  if (filtered.length === 0) {
    const noResults = document.createElement('div');
    noResults.className = 'empty-state';
    noResults.innerHTML = `
      <div class="empty-icon">
        <i data-lucide="search" class="w-12 h-12 mx-auto text-content-disabled"></i>
      </div>
      <div class="empty-title">No employees found</div>
      <p class="empty-description">Try adjusting your search or filter criteria.</p>
    `;
    containerRef.appendChild(noResults);
    try { createIcons(); } catch { /* ok */ }
    return;
  }

  // Employee table
  const tableWrapper = document.createElement('div');
  tableWrapper.className = 'card table-scroll';
  tableWrapper.style.padding = '0';

  const table = document.createElement('table');
  table.className = 'data-table';

  // Table head
  const thead = document.createElement('thead');
  const headerRow = document.createElement('tr');
  const columns = ['Employee', 'Department', 'Status', "Today's Hours", 'Rate', ''];
  for (const col of columns) {
    const th = document.createElement('th');
    th.textContent = col;
    if (col === '') th.style.width = '80px';
    headerRow.appendChild(th);
  }
  thead.appendChild(headerRow);
  table.appendChild(thead);

  // Table body
  const tbody = document.createElement('tbody');
  for (const employee of filtered) {
    tbody.appendChild(buildEmployeeRow(employee));
  }
  table.appendChild(tbody);

  tableWrapper.appendChild(table);
  containerRef.appendChild(tableWrapper);

  // Re-init lucide icons
  try { createIcons(); } catch { /* ok */ }
}

function onDataChange(): void {
  renderView();
}

export const render: ViewModule['render'] = (container) => {
  containerRef = container;
  searchQuery = '';
  statusFilter = 'all';
  departmentFilter = 'all';

  eventBus.on('employees:changed', onDataChange);
  eventBus.on('timeEntries:changed', onDataChange);
  eventBus.on('timer:changed', onDataChange);

  renderView();
};

export const destroy = (): void => {
  eventBus.off('employees:changed', onDataChange);
  eventBus.off('timeEntries:changed', onDataChange);
  eventBus.off('timer:changed', onDataChange);
  if (debounceTimer) clearTimeout(debounceTimer);
  containerRef = null;
};
