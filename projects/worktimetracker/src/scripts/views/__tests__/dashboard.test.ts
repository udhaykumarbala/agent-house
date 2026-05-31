import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { employeeService } from '../../data/employeeService';
import { projectService } from '../../data/projectService';
import { timeEntryService } from '../../data/timeEntryService';
import { store } from '../../data/store';
import { eventBus } from '../../utils/eventBus';

// Mock lucide
vi.mock('lucide', () => ({
  createIcons: vi.fn(),
}));

describe('Dashboard View', () => {
  let container: HTMLElement;
  let dashboardModule: typeof import('../dashboard');

  beforeEach(async () => {
    store.clearAll();
    eventBus._reset();
    projectService.ensureGeneralProject();

    container = document.createElement('div');
    document.body.appendChild(container);

    dashboardModule = await import('../dashboard');
  });

  afterEach(() => {
    dashboardModule.destroy?.();
    document.body.removeChild(container);
  });

  it('should render the Dashboard heading', () => {
    dashboardModule.render(container);
    const h1 = container.querySelector('h1');
    expect(h1).not.toBeNull();
    expect(h1!.textContent).toBe('Dashboard');
  });

  it('should show empty state when no employees exist', () => {
    dashboardModule.render(container);
    const emptyTitle = container.querySelector('.empty-title');
    expect(emptyTitle).not.toBeNull();
    expect(emptyTitle!.textContent).toContain('Welcome');
  });

  it('should show summary cards when employees exist', () => {
    employeeService.create({ name: 'John Doe', role: 'Developer' });
    dashboardModule.render(container);

    const statCards = container.querySelectorAll('.stat-card');
    expect(statCards.length).toBe(4);
  });

  it('should display correct total hours today in summary card', () => {
    employeeService.create({ name: 'Jane Smith', role: 'Designer' });
    dashboardModule.render(container);

    const hoursEl = container.querySelector('[data-stat="total-hours"]');
    expect(hoursEl).not.toBeNull();
    expect(hoursEl!.textContent).toContain('h');
  });

  it('should display clocked in count', () => {
    employeeService.create({ name: 'Alice', role: 'Dev' });
    dashboardModule.render(container);

    const text = container.textContent || '';
    // Should show "0/1" for clocked in
    expect(text).toContain('0/1');
  });

  it('should show Quick Clock In section', () => {
    employeeService.create({ name: 'Bob', role: 'QA' });
    dashboardModule.render(container);

    const text = container.textContent || '';
    expect(text).toContain('Quick Clock In');
  });

  it('should have employee and project dropdowns in Quick Clock In', () => {
    employeeService.create({ name: 'Charlie', role: 'Dev' });
    dashboardModule.render(container);

    const selects = container.querySelectorAll('select');
    expect(selects.length).toBeGreaterThanOrEqual(2);
  });

  it('should show active timers section when employees are clocked in', () => {
    const emp = employeeService.create({ name: 'Dave', role: 'Dev' });
    timeEntryService.clockIn(emp.id, 'general');

    dashboardModule.render(container);

    const text = container.textContent || '';
    expect(text).toContain('Active Timers');
    expect(text).toContain('Dave');
  });

  it('should show clock out button for active timers', () => {
    const emp = employeeService.create({ name: 'Eve', role: 'Dev' });
    timeEntryService.clockIn(emp.id, 'general');

    dashboardModule.render(container);

    const stopBtns = container.querySelectorAll('.btn-clock-out');
    expect(stopBtns.length).toBeGreaterThan(0);
  });

  it('should not show activity section when no entries today', () => {
    employeeService.create({ name: 'Frank', role: 'Dev' });
    dashboardModule.render(container);

    const text = container.textContent || '';
    expect(text).toContain('No activity yet today');
  });

  it('should show Add Employee CTA in empty state', () => {
    dashboardModule.render(container);

    const ctaLink = container.querySelector('a[href="#/employees"]');
    expect(ctaLink).not.toBeNull();
    expect(ctaLink!.textContent).toContain('Add Employee');
  });

  it('should exclude clocked-in employees from Quick Clock In dropdown', () => {
    const emp1 = employeeService.create({ name: 'Gary', role: 'Dev' });
    const emp2 = employeeService.create({ name: 'Helen', role: 'Dev' });
    timeEntryService.clockIn(emp1.id, 'general');

    dashboardModule.render(container);

    const empSelect = container.querySelector('#quick-employee') as HTMLSelectElement;
    expect(empSelect).not.toBeNull();
    const options = Array.from(empSelect.options).map(o => o.value);
    // Gary should not be in the dropdown (already clocked in)
    expect(options).not.toContain(emp1.id);
    // Helen should be in the dropdown
    expect(options).toContain(emp2.id);
  });

  it('should update on eventBus events', () => {
    const emp = employeeService.create({ name: 'Ivy', role: 'Dev' });
    dashboardModule.render(container);

    // Clock in triggers event which re-renders
    timeEntryService.clockIn(emp.id, 'general');

    // The eventBus handler should trigger re-render
    const text = container.textContent || '';
    expect(text).toContain('Active Timers');
  });
});
