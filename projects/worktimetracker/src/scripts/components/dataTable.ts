/**
 * DataTable Component
 *
 * Sortable, filterable table built from arrays of objects.
 * - Column sorting (asc/desc toggle)
 * - Search filter across all string columns (300ms debounce)
 * - Row actions (edit, delete, etc.) appear on hover
 * - Uses textContent for all user data (XSS prevention)
 * - font-variant-numeric: tabular-nums on cells
 */

export interface Column {
  key: string;
  label: string;
  sortable?: boolean;
  render?: (value: unknown, row: Record<string, unknown>) => HTMLElement;
}

export interface RowAction {
  label: string;
  onClick: (row: Record<string, unknown>) => void;
}

export interface DataTableOptions {
  columns: Column[];
  data: Record<string, unknown>[];
  searchable?: boolean;
  searchPlaceholder?: string;
  onRowClick?: (row: Record<string, unknown>) => void;
  actions?: RowAction[];
  emptyMessage?: string;
}

type SortDirection = 'asc' | 'desc' | null;

export class DataTable {
  private container: HTMLElement;
  private options: DataTableOptions;
  private sortColumn: string | null = null;
  private sortDirection: SortDirection = null;
  private searchTerm = '';

  constructor(container: HTMLElement, options: DataTableOptions) {
    this.container = container;
    this.options = options;
    this.render();
  }

  setData(data: Record<string, unknown>[]): void {
    this.options.data = data;
    this.render();
  }

  refresh(): void {
    this.render();
  }

  private getFilteredData(): Record<string, unknown>[] {
    let data = [...this.options.data];

    // Apply search filter
    if (this.searchTerm) {
      const term = this.searchTerm.toLowerCase();
      data = data.filter((row) =>
        this.options.columns.some((col) => {
          const val = row[col.key];
          return val != null && String(val).toLowerCase().includes(term);
        })
      );
    }

    // Apply sort
    if (this.sortColumn && this.sortDirection) {
      const col = this.sortColumn;
      const dir = this.sortDirection === 'asc' ? 1 : -1;

      data.sort((a, b) => {
        const valA = a[col];
        const valB = b[col];

        if (valA == null && valB == null) return 0;
        if (valA == null) return 1;
        if (valB == null) return -1;

        if (typeof valA === 'number' && typeof valB === 'number') {
          return (valA - valB) * dir;
        }

        return String(valA).localeCompare(String(valB)) * dir;
      });
    }

    return data;
  }

  private handleSort(columnKey: string): void {
    if (this.sortColumn === columnKey) {
      if (this.sortDirection === 'asc') {
        this.sortDirection = 'desc';
      } else if (this.sortDirection === 'desc') {
        this.sortColumn = null;
        this.sortDirection = null;
      }
    } else {
      this.sortColumn = columnKey;
      this.sortDirection = 'asc';
    }

    this.render();
  }

  private handleSearch(term: string): void {
    this.searchTerm = term;
    this.render();
  }

  private render(): void {
    this.container.innerHTML = '';

    const wrapper = document.createElement('div');

    // Search input
    if (this.options.searchable) {
      const searchWrapper = document.createElement('div');
      searchWrapper.style.marginBottom = '16px';

      const searchInput = document.createElement('input');
      searchInput.type = 'text';
      searchInput.className = 'input';
      searchInput.placeholder = this.options.searchPlaceholder ?? 'Search...';
      searchInput.value = this.searchTerm;

      let debounceTimer: ReturnType<typeof setTimeout>;
      searchInput.addEventListener('input', () => {
        clearTimeout(debounceTimer);
        debounceTimer = setTimeout(() => {
          this.handleSearch(searchInput.value);
        }, 300);
      });

      searchWrapper.appendChild(searchInput);
      wrapper.appendChild(searchWrapper);
    }

    const filteredData = this.getFilteredData();

    // Table
    const table = document.createElement('table');
    table.className = 'data-table';

    // Thead
    const thead = document.createElement('thead');
    const headerRow = document.createElement('tr');

    for (const col of this.options.columns) {
      const th = document.createElement('th');
      th.setAttribute('scope', 'col');

      // Label (textContent for XSS safety)
      const labelSpan = document.createElement('span');
      labelSpan.textContent = col.label;
      th.appendChild(labelSpan);

      if (col.sortable !== false) {
        // Sort indicator
        const indicator = document.createElement('span');
        indicator.className = 'sort-indicator';

        if (this.sortColumn === col.key) {
          indicator.classList.add('active');
          indicator.textContent = this.sortDirection === 'asc' ? ' \u25B2' : ' \u25BC';
        } else {
          indicator.textContent = ' \u25B4';
        }

        th.appendChild(indicator);
        th.addEventListener('click', () => this.handleSort(col.key));
      }

      headerRow.appendChild(th);
    }

    // Actions column header
    if (this.options.actions && this.options.actions.length > 0) {
      const actionTh = document.createElement('th');
      actionTh.setAttribute('scope', 'col');
      actionTh.textContent = '';
      headerRow.appendChild(actionTh);
    }

    thead.appendChild(headerRow);
    table.appendChild(thead);

    // Tbody
    const tbody = document.createElement('tbody');

    if (filteredData.length === 0) {
      const emptyRow = document.createElement('tr');
      const emptyCell = document.createElement('td');
      const colSpan = this.options.columns.length + (this.options.actions ? 1 : 0);
      emptyCell.setAttribute('colspan', String(colSpan));
      emptyCell.className = 'empty-state';
      emptyCell.textContent = this.options.emptyMessage ?? 'No data available';
      emptyRow.appendChild(emptyCell);
      tbody.appendChild(emptyRow);
    } else {
      for (const row of filteredData) {
        const tr = document.createElement('tr');

        if (this.options.onRowClick) {
          tr.style.cursor = 'pointer';
          tr.addEventListener('click', (e) => {
            // Don't trigger row click if an action button was clicked
            if ((e.target as HTMLElement).closest('.row-actions')) return;
            this.options.onRowClick!(row);
          });
        }

        for (const col of this.options.columns) {
          const td = document.createElement('td');

          if (col.render) {
            // Custom render returns an HTMLElement (caller is responsible for safe creation)
            const rendered = col.render(row[col.key], row);
            td.appendChild(rendered);
          } else {
            // Default: textContent for XSS safety
            const val = row[col.key];
            td.textContent = val != null ? String(val) : '';
          }

          tr.appendChild(td);
        }

        // Row actions
        if (this.options.actions && this.options.actions.length > 0) {
          const actionTd = document.createElement('td');
          const actionsDiv = document.createElement('div');
          actionsDiv.className = 'row-actions';

          for (const action of this.options.actions) {
            const btn = document.createElement('button');
            btn.className = 'action-btn';
            btn.type = 'button';
            btn.textContent = action.label;
            btn.addEventListener('click', (e) => {
              e.stopPropagation();
              action.onClick(row);
            });
            actionsDiv.appendChild(btn);
          }

          actionTd.appendChild(actionsDiv);
          tr.appendChild(actionTd);
        }

        tbody.appendChild(tr);
      }
    }

    table.appendChild(tbody);
    wrapper.appendChild(table);
    this.container.appendChild(wrapper);
  }
}
