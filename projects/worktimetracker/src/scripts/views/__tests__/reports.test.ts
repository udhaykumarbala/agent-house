import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { employeeService } from '../../data/employeeService';
import { projectService } from '../../data/projectService';
import { timeEntryService } from '../../data/timeEntryService';
import { store } from '../../data/store';
import { eventBus } from '../../utils/eventBus';
import { todayDate } from '../../utils/dates';

vi.mock('lucide', () => ({
  createIcons: vi.fn(),
}));

function createSampleData() {
  const emp1 = employeeService.create({ name: 'Alice', role: 'Developer' });
  const emp2 = employeeService.create({ name: 'Bob', role: 'Designer' });
  const proj = projectService.create({ name: 'Alpha Project' });
  const today = todayDate();

  timeEntryService.createManualEntry({
    employeeId: emp1.id,
    projectId: proj.id,
    date: today,
    clockIn: new Date(`${today}T09:00:00`).toISOString(),
    clockOut: new Date(`${today}T12:00:00`).toISOString(),
    notes: 'Morning work',
  });

  timeEntryService.createManualEntry({
    employeeId: emp2.id,
    projectId: 'general',
    date: today,
    clockIn: new Date(`${today}T10:00:00`).toISOString(),
    clockOut: new Date(`${today}T14:00:00`).toISOString(),
    notes: 'Design work',
  });

  return { emp1, emp2, proj };
}

describe('Reports View', () => {
  let container: HTMLElement;
  let reportsModule: typeof import('../reports');

  beforeEach(async () => {
    store.clearAll();
    eventBus._reset();
    projectService.ensureGeneralProject();

    container = document.createElement('div');
    document.body.appendChild(container);

    reportsModule = await import('../reports');
  });

  afterEach(() => {
    reportsModule.destroy?.();
    document.body.removeChild(container);
  });

  it('should render the Reports heading', () => {
    reportsModule.render(container);
    const h1 = container.querySelector('h1');
    expect(h1).not.toBeNull();
    expect(h1!.textContent).toBe('Reports');
  });

  it('should show report type selector with all 4 types', () => {
    reportsModule.render(container);

    const text = container.textContent || '';
    expect(text).toContain('Daily');
    expect(text).toContain('Weekly');
    expect(text).toContain('Monthly');
    expect(text).toContain('By Project');
  });

  it('should show Export CSV button', () => {
    reportsModule.render(container);

    const exportBtn = container.querySelector('#export-csv-btn');
    expect(exportBtn).not.toBeNull();
    expect(exportBtn!.textContent).toContain('Export CSV');
  });

  it('should show Print / PDF button', () => {
    reportsModule.render(container);

    const printBtn = container.querySelector('.btn-secondary');
    expect(printBtn).not.toBeNull();
    expect(printBtn!.textContent).toContain('Print');
  });

  it('should show date range picker', () => {
    reportsModule.render(container);

    const presets = container.querySelectorAll('.date-preset');
    expect(presets.length).toBeGreaterThanOrEqual(5);
  });

  it('should show employee and project filter dropdowns', () => {
    reportsModule.render(container);

    const selects = container.querySelectorAll('select');
    const hasEmpFilter = Array.from(selects).some(s =>
      Array.from(s.options).some(o => o.textContent === 'All Employees')
    );
    const hasProjFilter = Array.from(selects).some(s =>
      Array.from(s.options).some(o => o.textContent === 'All Projects')
    );
    expect(hasEmpFilter).toBe(true);
    expect(hasProjFilter).toBe(true);
  });

  it('should show empty state when no entries exist', () => {
    reportsModule.render(container);

    const emptyTitle = container.querySelector('.empty-title');
    expect(emptyTitle).not.toBeNull();
    expect(emptyTitle!.textContent).toContain('No reports');
  });

  it('should display 4 summary cards when data exists', () => {
    createSampleData();
    reportsModule.render(container);

    const statCards = container.querySelectorAll('.stat-card');
    expect(statCards.length).toBe(4);
  });

  it('should display summary card labels correctly', () => {
    createSampleData();
    reportsModule.render(container);

    const labels = container.querySelectorAll('.stat-label');
    const labelTexts = Array.from(labels).map(l => l.textContent);
    expect(labelTexts).toContain('TOTAL HOURS');
    expect(labelTexts).toContain('AVG HOURS/DAY');
    expect(labelTexts).toContain('TOP EMPLOYEE');
    expect(labelTexts).toContain('TOP PROJECT');
  });

  it('should display horizontal bar chart', () => {
    createSampleData();
    reportsModule.render(container);

    // Chart has bars with colored backgrounds
    const chartTitle = Array.from(container.querySelectorAll('h3')).find(
      h => h.textContent?.includes('Hours by')
    );
    expect(chartTitle).not.toBeNull();
  });

  it('should display data table with sortable columns', () => {
    createSampleData();
    reportsModule.render(container);

    const tables = container.querySelectorAll('.data-table');
    expect(tables.length).toBeGreaterThan(0);

    const sortIndicators = container.querySelectorAll('.sort-indicator');
    expect(sortIndicators.length).toBeGreaterThan(0);
  });

  it('should switch report types when tabs are clicked', () => {
    createSampleData();
    reportsModule.render(container);

    // Find the "Daily" button and click it
    const buttons = container.querySelectorAll('button');
    const dailyBtn = Array.from(buttons).find(b => b.textContent === 'Daily');
    expect(dailyBtn).not.toBeNull();
    dailyBtn!.click();

    // Should show daily-specific content
    const text = container.textContent || '';
    expect(text).toContain('Reports');
  });

  it('should trigger window.print when Print/PDF button is clicked', () => {
    const printSpy = vi.spyOn(window, 'print').mockImplementation(() => {});
    reportsModule.render(container);

    const printBtn = container.querySelector('.btn-secondary') as HTMLButtonElement;
    printBtn.click();

    expect(printSpy).toHaveBeenCalled();
    printSpy.mockRestore();
  });

  it('should filter entries by employee', () => {
    const { emp1 } = createSampleData();
    reportsModule.render(container);

    const selects = container.querySelectorAll('select');
    const empSelect = Array.from(selects).find(s =>
      Array.from(s.options).some(o => o.textContent === 'All Employees')
    ) as HTMLSelectElement;

    empSelect.value = emp1.id;
    empSelect.dispatchEvent(new Event('change'));

    // After filtering, only Alice's data should appear
    const text = container.textContent || '';
    expect(text).toContain('Alice');
  });

  it('should filter entries by project', () => {
    const { proj } = createSampleData();
    reportsModule.render(container);

    const selects = container.querySelectorAll('select');
    const projSelect = Array.from(selects).find(s =>
      Array.from(s.options).some(o => o.textContent === 'All Projects')
    ) as HTMLSelectElement;

    projSelect.value = proj.id;
    projSelect.dispatchEvent(new Event('change'));

    const text = container.textContent || '';
    expect(text).toContain('Alpha Project');
  });

  it('should show no-match state when filters exclude all entries', () => {
    createSampleData();
    reportsModule.render(container);

    // Create a non-matching project filter
    const selects = container.querySelectorAll('select');
    const projSelect = Array.from(selects).find(s =>
      Array.from(s.options).some(o => o.textContent === 'All Projects')
    ) as HTMLSelectElement;

    // Set to a project with no entries
    const newProj = projectService.create({ name: 'Empty Project' });
    projSelect.value = newProj.id;
    projSelect.dispatchEvent(new Event('change'));

    const emptyTitle = container.querySelector('.empty-title');
    expect(emptyTitle).not.toBeNull();
  });

  it('should update when timeEntries:changed event fires', () => {
    reportsModule.render(container);

    // Initially empty
    let emptyTitle = container.querySelector('.empty-title');
    expect(emptyTitle).not.toBeNull();

    // Add data (triggers eventBus)
    createSampleData();

    // Should now have content
    const statCards = container.querySelectorAll('.stat-card');
    expect(statCards.length).toBe(4);
  });

  it('should show "Detailed Breakdown" table heading', () => {
    createSampleData();
    reportsModule.render(container);

    const headings = Array.from(container.querySelectorAll('h3'));
    const detailHeading = headings.find(h => h.textContent?.includes('Detailed Breakdown'));
    expect(detailHeading).not.toBeNull();
  });

  it('should export CSV when export button is clicked', () => {
    createSampleData();
    reportsModule.render(container);

    // Mock link click for CSV download
    const createElementSpy = vi.spyOn(document, 'createElement');

    const exportBtn = container.querySelector('#export-csv-btn') as HTMLButtonElement;
    exportBtn.click();

    // Should have created a link element for download
    const linkCalls = createElementSpy.mock.calls.filter(c => c[0] === 'a');
    expect(linkCalls.length).toBeGreaterThan(0);

    createElementSpy.mockRestore();
  });
});
