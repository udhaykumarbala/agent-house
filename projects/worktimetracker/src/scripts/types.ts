export interface Employee {
  id: string;
  name: string;
  role: string;
  department: string;
  hourlyRate: number | null;
  status: 'active' | 'inactive';
  color: string;
  createdAt: string;
}

export interface Project {
  id: string;
  name: string;
  description: string;
  color: string;
  status: 'active' | 'completed' | 'archived';
  createdAt: string;
}

export interface Task {
  id: string;
  projectId: string;
  name: string;
  status: 'todo' | 'in_progress' | 'done';
}

export interface TimeEntry {
  id: string;
  employeeId: string;
  projectId: string;
  taskId: string | null;
  date: string;
  clockIn: string;
  clockOut: string | null;
  duration: number | null;
  notes: string;
  type: 'clock' | 'manual';
}

export interface ActiveTimer {
  employeeId: string;
  projectId: string;
  taskId: string | null;
  clockIn: string;
  notes: string;
}

export interface AppSettings {
  timeFormat: '12h' | '24h';
  weekStartsOn: 'monday' | 'sunday';
  defaultView: string;
}

export type ToastType = 'success' | 'error' | 'warning' | 'info';

export interface ViewModule {
  render: (container: HTMLElement) => void | Promise<void>;
  destroy?: () => void;
}
