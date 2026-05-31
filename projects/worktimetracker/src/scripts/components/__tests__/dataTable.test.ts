import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { DataTable } from '../dataTable';

describe('DataTable', () => {
  let container: HTMLElement;

  const sampleData = [
    { name: 'Sarah Chen', role: 'Developer', hours: 8.5 },
    { name: 'Mike Ross', role: 'Designer', hours: 6.0 },
    { name: 'Lisa Park', role: 'Developer', hours: 7.25 },
    { name: 'John Lee', role: 'Manager', hours: 4.0 },
  ];

  const columns = [
    { key: 'name', label: 'Name' },
    { key: 'role', label: 'Role' },
    { key: 'hours', label: 'Hours' },
  ];

  beforeEach(() => {
    container = document.createElement('div');
    document.body.appendChild(container);
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
    document.body.innerHTML = '';
  });

  it('should render a table with data', () => {
    new DataTable(container, { columns, data: sampleData });

    const table = container.querySelector('.data-table');
    expect(table).toBeTruthy();

    const rows = container.querySelectorAll('tbody tr');
    expect(rows.length).toBe(4);
  });

  it('should render column headers', () => {
    new DataTable(container, { columns, data: sampleData });

    const headers = container.querySelectorAll('thead th');
    expect(headers.length).toBe(3);
    expect(headers[0].textContent).toContain('Name');
    expect(headers[1].textContent).toContain('Role');
    expect(headers[2].textContent).toContain('Hours');
  });

  it('should render cell values using textContent (XSS safe)', () => {
    const xssData = [
      { name: '<script>alert("xss")</script>', role: 'Hacker', hours: 0 },
    ];

    new DataTable(container, { columns, data: xssData });

    const cell = container.querySelector('tbody td');
    expect(cell?.textContent).toBe('<script>alert("xss")</script>');
    expect(cell?.innerHTML).not.toContain('<script>');
  });

  it('should show empty state when no data', () => {
    new DataTable(container, { columns, data: [] });

    const emptyCell = container.querySelector('.empty-state');
    expect(emptyCell).toBeTruthy();
    expect(emptyCell?.textContent).toBe('No data available');
  });

  it('should show custom empty message', () => {
    new DataTable(container, {
      columns,
      data: [],
      emptyMessage: 'No employees found',
    });

    const emptyCell = container.querySelector('.empty-state');
    expect(emptyCell?.textContent).toBe('No employees found');
  });

  it('should sort ascending on first column click', () => {
    new DataTable(container, { columns, data: sampleData });

    // Click the "Name" column header
    const nameHeader = container.querySelector('thead th') as HTMLElement;
    nameHeader.click();

    const cells = container.querySelectorAll('tbody tr td:first-child');
    const names = Array.from(cells).map((c) => c.textContent);
    expect(names).toEqual(['John Lee', 'Lisa Park', 'Mike Ross', 'Sarah Chen']);
  });

  it('should sort descending on second click', () => {
    new DataTable(container, { columns, data: sampleData });

    const nameHeader = container.querySelector('thead th') as HTMLElement;
    nameHeader.click(); // asc
    nameHeader.click(); // desc

    const cells = container.querySelectorAll('tbody tr td:first-child');
    const names = Array.from(cells).map((c) => c.textContent);
    expect(names).toEqual(['Sarah Chen', 'Mike Ross', 'Lisa Park', 'John Lee']);
  });

  it('should clear sort on third click', () => {
    new DataTable(container, { columns, data: sampleData });

    const nameHeader = container.querySelector('thead th') as HTMLElement;
    nameHeader.click(); // asc
    nameHeader.click(); // desc
    nameHeader.click(); // clear

    const cells = container.querySelectorAll('tbody tr td:first-child');
    const names = Array.from(cells).map((c) => c.textContent);
    // Back to original order
    expect(names).toEqual(['Sarah Chen', 'Mike Ross', 'Lisa Park', 'John Lee']);
  });

  it('should sort numeric columns correctly', () => {
    new DataTable(container, { columns, data: sampleData });

    // Click "Hours" header (third column)
    const headers = container.querySelectorAll('thead th');
    (headers[2] as HTMLElement).click();

    const cells = container.querySelectorAll('tbody tr td:nth-child(3)');
    const hours = Array.from(cells).map((c) => c.textContent);
    expect(hours).toEqual(['4', '6', '7.25', '8.5']);
  });

  it('should show sort indicator on active column', () => {
    new DataTable(container, { columns, data: sampleData });

    const nameHeader = container.querySelector('thead th') as HTMLElement;
    nameHeader.click();

    // Re-query after render (sort triggers full re-render)
    const updatedHeader = container.querySelector('thead th') as HTMLElement;
    const indicator = updatedHeader.querySelector('.sort-indicator.active');
    expect(indicator).toBeTruthy();
    expect(indicator?.textContent).toContain('\u25B2'); // up arrow
  });

  it('should not sort when sortable is false', () => {
    const unsortableColumns = [
      { key: 'name', label: 'Name', sortable: false },
      { key: 'role', label: 'Role' },
    ];

    new DataTable(container, { columns: unsortableColumns, data: sampleData });

    const nameHeader = container.querySelector('thead th') as HTMLElement;
    // No sort indicator should be present
    expect(nameHeader.querySelector('.sort-indicator')).toBeNull();
  });

  it('should filter data with search input', () => {
    new DataTable(container, {
      columns,
      data: sampleData,
      searchable: true,
    });

    const searchInput = container.querySelector('input') as HTMLInputElement;
    expect(searchInput).toBeTruthy();

    // Type search term
    searchInput.value = 'developer';
    searchInput.dispatchEvent(new Event('input'));

    // Wait for debounce (300ms)
    vi.advanceTimersByTime(300);

    const rows = container.querySelectorAll('tbody tr');
    expect(rows.length).toBe(2); // Sarah and Lisa
  });

  it('should show empty state when search returns no results', () => {
    new DataTable(container, {
      columns,
      data: sampleData,
      searchable: true,
    });

    const searchInput = container.querySelector('input') as HTMLInputElement;
    searchInput.value = 'zzzzz';
    searchInput.dispatchEvent(new Event('input'));

    vi.advanceTimersByTime(300);

    const emptyCell = container.querySelector('.empty-state');
    expect(emptyCell).toBeTruthy();
  });

  it('should use custom search placeholder', () => {
    new DataTable(container, {
      columns,
      data: sampleData,
      searchable: true,
      searchPlaceholder: 'Filter employees...',
    });

    const searchInput = container.querySelector('input') as HTMLInputElement;
    expect(searchInput.placeholder).toBe('Filter employees...');
  });

  it('should render row actions', () => {
    const editFn = vi.fn();
    new DataTable(container, {
      columns,
      data: sampleData,
      actions: [{ label: 'Edit', onClick: editFn }],
    });

    // Should have 4 action columns (one per row) + actions header
    const actionBtns = container.querySelectorAll('.action-btn');
    expect(actionBtns.length).toBe(4);
    expect(actionBtns[0].textContent).toBe('Edit');
  });

  it('should call action onClick with the correct row data', () => {
    const editFn = vi.fn();
    new DataTable(container, {
      columns,
      data: sampleData,
      actions: [{ label: 'Edit', onClick: editFn }],
    });

    const firstActionBtn = container.querySelector('.action-btn') as HTMLElement;
    firstActionBtn.click();

    expect(editFn).toHaveBeenCalledTimes(1);
    expect(editFn).toHaveBeenCalledWith(sampleData[0]);
  });

  it('should call onRowClick when a row is clicked', () => {
    const rowClick = vi.fn();
    new DataTable(container, {
      columns,
      data: sampleData,
      onRowClick: rowClick,
    });

    const firstRow = container.querySelector('tbody tr') as HTMLElement;
    firstRow.click();

    expect(rowClick).toHaveBeenCalledTimes(1);
    expect(rowClick).toHaveBeenCalledWith(sampleData[0]);
  });

  it('should NOT trigger onRowClick when action button is clicked', () => {
    const rowClick = vi.fn();
    const editFn = vi.fn();
    new DataTable(container, {
      columns,
      data: sampleData,
      onRowClick: rowClick,
      actions: [{ label: 'Edit', onClick: editFn }],
    });

    const firstActionBtn = container.querySelector('.action-btn') as HTMLElement;
    firstActionBtn.click();

    expect(editFn).toHaveBeenCalledTimes(1);
    expect(rowClick).not.toHaveBeenCalled();
  });

  it('should support custom render functions', () => {
    const customColumns = [
      {
        key: 'name',
        label: 'Name',
        render: (value: unknown) => {
          const el = document.createElement('strong');
          el.textContent = String(value);
          return el;
        },
      },
      { key: 'role', label: 'Role' },
    ];

    new DataTable(container, { columns: customColumns, data: sampleData });

    const firstCell = container.querySelector('tbody td');
    const strong = firstCell?.querySelector('strong');
    expect(strong).toBeTruthy();
    expect(strong?.textContent).toBe('Sarah Chen');
  });

  it('should update data via setData()', () => {
    const dt = new DataTable(container, { columns, data: sampleData });

    expect(container.querySelectorAll('tbody tr').length).toBe(4);

    dt.setData([{ name: 'New Person', role: 'Intern', hours: 1 }]);

    expect(container.querySelectorAll('tbody tr').length).toBe(1);
    expect(container.querySelector('tbody td')?.textContent).toBe('New Person');
  });

  it('should have proper th scope attributes', () => {
    new DataTable(container, { columns, data: sampleData });

    const headers = container.querySelectorAll('thead th');
    headers.forEach((th) => {
      expect(th.getAttribute('scope')).toBe('col');
    });
  });

  it('should debounce search input (300ms)', () => {
    new DataTable(container, {
      columns,
      data: sampleData,
      searchable: true,
    });

    const searchInput = container.querySelector('input') as HTMLInputElement;
    searchInput.value = 'dev';
    searchInput.dispatchEvent(new Event('input'));

    // Before debounce fires, still shows all data
    expect(container.querySelectorAll('tbody tr').length).toBe(4);

    vi.advanceTimersByTime(300);

    // After debounce, shows filtered data
    expect(container.querySelectorAll('tbody tr').length).toBe(2);
  });
});
