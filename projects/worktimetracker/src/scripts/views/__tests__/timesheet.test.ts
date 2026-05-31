import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { employeeService } from '../../data/employeeService';
import { projectService } from '../../data/projectService';
import { timeEntryService } from '../../data/timeEntryService';
import { store } from '../../data/store';
import { eventBus } from '../../utils/eventBus';
import { Modal } from '../../components/modal';
import { todayDate } from '../../utils/dates';

vi.mock('lucide', () => ({
  createIcons: vi.fn(),
}));

describe('Timesheet View', () => {
  let container: HTMLElement;
  let timesheetModule: typeof import('../timesheet');

  beforeEach(async () => {
    store.clearAll();
    eventBus._reset();
    Modal._reset();
    projectService.ensureGeneralProject();

    container = document.createElement('div');
    document.body.appendChild(container);

    timesheetModule = await import('../timesheet');
  });

  afterEach(() => {
    timesheetModule.destroy?.();
    document.body.removeChild(container);
  });

  it('should render the Timesheet heading', () => {
    timesheetModule.render(container);
    const h1 = container.querySelector('h1');
    expect(h1).not.toBeNull();
    expect(h1!.textContent).toBe('Timesheet');
  });

  it('should show Add Entry button', () => {
    timesheetModule.render(container);
    const addBtn = container.querySelector('.btn-primary');
    expect(addBtn).not.toBeNull();
    expect(addBtn!.textContent).toContain('Add Entry');
  });

  it('should show date range picker with preset buttons', () => {
    timesheetModule.render(container);
    const presets = container.querySelectorAll('.date-preset');
    expect(presets.length).toBeGreaterThanOrEqual(5);
  });

  it('should show employee filter dropdown', () => {
    employeeService.create({ name: 'Alice', role: 'Dev' });
    timesheetModule.render(container);

    const selects = container.querySelectorAll('select');
    const empSelect = Array.from(selects).find(s =>
      Array.from(s.options).some(o => o.textContent === 'All Employees')
    );
    expect(empSelect).not.toBeNull();
  });

  it('should show project filter dropdown', () => {
    timesheetModule.render(container);

    const selects = container.querySelectorAll('select');
    const projSelect = Array.from(selects).find(s =>
      Array.from(s.options).some(o => o.textContent === 'All Projects')
    );
    expect(projSelect).not.toBeNull();
  });

  it('should display empty state when no entries exist', () => {
    timesheetModule.render(container);

    const emptyTitle = container.querySelector('.empty-title');
    expect(emptyTitle).not.toBeNull();
    expect(emptyTitle!.textContent).toContain('No time entries');
  });

  it('should display time entries in a data table', () => {
    const emp = employeeService.create({ name: 'Bob', role: 'Dev' });
    const today = todayDate();
    const clockIn = new Date(`${today}T09:00:00`).toISOString();
    const clockOut = new Date(`${today}T17:00:00`).toISOString();

    timeEntryService.createManualEntry({
      employeeId: emp.id,
      projectId: 'general',
      date: today,
      clockIn,
      clockOut,
      notes: 'Test entry',
    });

    timesheetModule.render(container);

    const table = container.querySelector('.data-table');
    expect(table).not.toBeNull();

    const rows = table!.querySelectorAll('tbody tr');
    expect(rows.length).toBeGreaterThan(0);
  });

  it('should display duration in Xh Ym format', () => {
    const emp = employeeService.create({ name: 'Charlie', role: 'Dev' });
    const today = todayDate();
    const clockIn = new Date(`${today}T09:00:00`).toISOString();
    const clockOut = new Date(`${today}T11:30:00`).toISOString();

    timeEntryService.createManualEntry({
      employeeId: emp.id,
      projectId: 'general',
      date: today,
      clockIn,
      clockOut,
    });

    timesheetModule.render(container);

    const text = container.textContent || '';
    expect(text).toMatch(/\d+h \d+m/);
  });

  it('should have sortable column headers', () => {
    const emp = employeeService.create({ name: 'Dave', role: 'Dev' });
    const today = todayDate();
    timeEntryService.createManualEntry({
      employeeId: emp.id,
      projectId: 'general',
      date: today,
      clockIn: new Date(`${today}T09:00:00`).toISOString(),
      clockOut: new Date(`${today}T17:00:00`).toISOString(),
    });

    timesheetModule.render(container);

    const sortIndicators = container.querySelectorAll('.sort-indicator');
    expect(sortIndicators.length).toBeGreaterThanOrEqual(4);
  });

  it('should show delete button in row actions', () => {
    const emp = employeeService.create({ name: 'Eve', role: 'Dev' });
    const today = todayDate();
    timeEntryService.createManualEntry({
      employeeId: emp.id,
      projectId: 'general',
      date: today,
      clockIn: new Date(`${today}T09:00:00`).toISOString(),
      clockOut: new Date(`${today}T10:00:00`).toISOString(),
    });

    timesheetModule.render(container);

    const actionBtns = container.querySelectorAll('.action-btn');
    expect(actionBtns.length).toBeGreaterThan(0);
  });

  it('should delete entry when delete button is clicked', () => {
    const emp = employeeService.create({ name: 'Frank', role: 'Dev' });
    const today = todayDate();
    timeEntryService.createManualEntry({
      employeeId: emp.id,
      projectId: 'general',
      date: today,
      clockIn: new Date(`${today}T09:00:00`).toISOString(),
      clockOut: new Date(`${today}T10:00:00`).toISOString(),
    });

    timesheetModule.render(container);

    expect(timeEntryService.getAll().length).toBe(1);

    const deleteBtn = container.querySelector('.action-btn') as HTMLButtonElement;
    deleteBtn.click();

    expect(timeEntryService.getAll().length).toBe(0);
  });

  it('should open Add Entry modal when button is clicked', () => {
    employeeService.create({ name: 'Grace', role: 'Dev' });
    timesheetModule.render(container);

    const addBtn = container.querySelector('.btn-primary') as HTMLButtonElement;
    addBtn.click();

    expect(Modal.isOpen()).toBe(true);
    const modalTitle = document.querySelector('.modal-title');
    expect(modalTitle?.textContent).toContain('Add Time Entry');
  });

  it('should filter entries by employee', () => {
    const emp1 = employeeService.create({ name: 'Hank', role: 'Dev' });
    const emp2 = employeeService.create({ name: 'Irene', role: 'QA' });
    const today = todayDate();

    timeEntryService.createManualEntry({
      employeeId: emp1.id,
      projectId: 'general',
      date: today,
      clockIn: new Date(`${today}T09:00:00`).toISOString(),
      clockOut: new Date(`${today}T10:00:00`).toISOString(),
    });
    timeEntryService.createManualEntry({
      employeeId: emp2.id,
      projectId: 'general',
      date: today,
      clockIn: new Date(`${today}T11:00:00`).toISOString(),
      clockOut: new Date(`${today}T12:00:00`).toISOString(),
    });

    timesheetModule.render(container);

    // Both should be visible initially
    let rows = container.querySelectorAll('tbody tr');
    expect(rows.length).toBe(2);

    // Filter to emp1
    const selects = container.querySelectorAll('select');
    const empSelect = Array.from(selects).find(s =>
      Array.from(s.options).some(o => o.textContent === 'All Employees')
    ) as HTMLSelectElement;

    empSelect.value = emp1.id;
    empSelect.dispatchEvent(new Event('change'));

    rows = container.querySelectorAll('tbody tr');
    expect(rows.length).toBe(1);
  });

  it('should update when timeEntries:changed event fires', () => {
    timesheetModule.render(container);

    const emp = employeeService.create({ name: 'Jack', role: 'Dev' });
    const today = todayDate();
    timeEntryService.createManualEntry({
      employeeId: emp.id,
      projectId: 'general',
      date: today,
      clockIn: new Date(`${today}T09:00:00`).toISOString(),
      clockOut: new Date(`${today}T10:00:00`).toISOString(),
    });

    // eventBus triggers re-render
    const table = container.querySelector('.data-table');
    expect(table).not.toBeNull();
  });

  it('should show entry count and total hours summary', () => {
    const emp = employeeService.create({ name: 'Kate', role: 'Dev' });
    const today = todayDate();

    timeEntryService.createManualEntry({
      employeeId: emp.id,
      projectId: 'general',
      date: today,
      clockIn: new Date(`${today}T09:00:00`).toISOString(),
      clockOut: new Date(`${today}T17:00:00`).toISOString(),
    });

    timesheetModule.render(container);

    const text = container.textContent || '';
    expect(text).toContain('1 entry');
    expect(text).toContain('Total:');
  });

  it('should support inline editing for notes cells', () => {
    const emp = employeeService.create({ name: 'Leo', role: 'Dev' });
    const today = todayDate();

    timeEntryService.createManualEntry({
      employeeId: emp.id,
      projectId: 'general',
      date: today,
      clockIn: new Date(`${today}T09:00:00`).toISOString(),
      clockOut: new Date(`${today}T10:00:00`).toISOString(),
      notes: 'Original note',
    });

    timesheetModule.render(container);

    // Find the notes cell (contains "Original note")
    const cells = container.querySelectorAll('td');
    const notesCell = Array.from(cells).find(td => td.textContent === 'Original note');
    expect(notesCell).not.toBeNull();

    // Click to edit
    notesCell!.click();

    // Should now contain an input
    const input = notesCell!.querySelector('input');
    expect(input).not.toBeNull();
    expect(input!.value).toBe('Original note');
  });
});
