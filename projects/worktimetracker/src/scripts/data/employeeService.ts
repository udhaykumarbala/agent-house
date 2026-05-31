import { store } from './store';
import { eventBus } from '../utils/eventBus';
import { validateRequired, validateMaxLength, sanitizeString } from '../utils/sanitize';
import type { Employee } from '../types';

const STORE_KEY = 'employees';

// Auto-assign colors from the chart palette
const EMPLOYEE_COLORS = [
  '#CF8A2E', '#5746B2', '#E07B5F', '#2E9E8F',
  '#C75A8A', '#5B8AD4', '#8BA055', '#9B6DB0',
];

function getNextColor(employees: Employee[]): string {
  const usedColors = new Set(employees.map(e => e.color));
  const available = EMPLOYEE_COLORS.find(c => !usedColors.has(c));
  return available || EMPLOYEE_COLORS[employees.length % EMPLOYEE_COLORS.length];
}

function generateId(): string {
  return crypto.randomUUID();
}

export interface CreateEmployeeInput {
  name: string;
  role?: string;
  department?: string;
  hourlyRate?: number | null;
}

export interface UpdateEmployeeInput {
  name?: string;
  role?: string;
  department?: string;
  hourlyRate?: number | null;
  status?: 'active' | 'inactive';
  color?: string;
}

class EmployeeService {
  getAll(): Employee[] {
    return store.get<Employee>(STORE_KEY);
  }

  getById(id: string): Employee | undefined {
    return this.getAll().find(e => e.id === id);
  }

  getActive(): Employee[] {
    return this.getAll().filter(e => e.status === 'active');
  }

  create(input: CreateEmployeeInput): Employee {
    const errors = this.validate(input.name);
    if (errors) throw new Error(errors);

    const employees = this.getAll();

    // Check uniqueness
    const nameNorm = input.name.trim().toLowerCase();
    if (employees.some(e => e.name.trim().toLowerCase() === nameNorm)) {
      throw new Error('An employee with this name already exists');
    }

    const employee: Employee = {
      id: generateId(),
      name: sanitizeString(input.name),
      role: sanitizeString(input.role || ''),
      department: sanitizeString(input.department || ''),
      hourlyRate: input.hourlyRate ?? null,
      status: 'active',
      color: getNextColor(employees),
      createdAt: new Date().toISOString(),
    };

    employees.push(employee);
    store.set(STORE_KEY, employees);
    eventBus.emit('employees:changed', employee);
    eventBus.emit('employees:created', employee);

    return employee;
  }

  update(id: string, input: UpdateEmployeeInput): Employee {
    const employees = this.getAll();
    const index = employees.findIndex(e => e.id === id);
    if (index === -1) throw new Error('Employee not found');

    if (input.name !== undefined) {
      const errors = this.validate(input.name, id);
      if (errors) throw new Error(errors);

      // Check uniqueness (excluding self)
      const nameNorm = input.name.trim().toLowerCase();
      if (employees.some(e => e.id !== id && e.name.trim().toLowerCase() === nameNorm)) {
        throw new Error('An employee with this name already exists');
      }
    }

    const updated: Employee = {
      ...employees[index],
      ...(input.name !== undefined && { name: sanitizeString(input.name) }),
      ...(input.role !== undefined && { role: sanitizeString(input.role) }),
      ...(input.department !== undefined && { department: sanitizeString(input.department) }),
      ...(input.hourlyRate !== undefined && { hourlyRate: input.hourlyRate }),
      ...(input.status !== undefined && { status: input.status }),
      ...(input.color !== undefined && { color: input.color }),
    };

    employees[index] = updated;
    store.set(STORE_KEY, employees);
    eventBus.emit('employees:changed', updated);
    eventBus.emit('employees:updated', updated);

    return updated;
  }

  delete(id: string): void {
    const employees = this.getAll();
    const filtered = employees.filter(e => e.id !== id);
    if (filtered.length === employees.length) throw new Error('Employee not found');

    store.set(STORE_KEY, filtered);
    eventBus.emit('employees:changed');
    eventBus.emit('employees:deleted', id);
  }

  private validate(name: string, excludeId?: string): string | null {
    void excludeId;
    const reqErr = validateRequired(name, 'Name');
    if (reqErr) return reqErr;

    const maxErr = validateMaxLength(name, 100, 'Name');
    if (maxErr) return maxErr;

    return null;
  }
}

export const employeeService = new EmployeeService();
