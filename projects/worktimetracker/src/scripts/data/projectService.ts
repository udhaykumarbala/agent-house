import { store } from './store';
import { eventBus } from '../utils/eventBus';
import { validateRequired, validateMaxLength, sanitizeString } from '../utils/sanitize';
import type { Project, Task } from '../types';

const PROJECTS_KEY = 'projects';
const TASKS_KEY = 'tasks';

const GENERAL_PROJECT_ID = 'general';

function generateId(): string {
  return crypto.randomUUID();
}

export interface CreateProjectInput {
  name: string;
  description?: string;
  color?: string;
  status?: Project['status'];
}

export interface UpdateProjectInput {
  name?: string;
  description?: string;
  color?: string;
  status?: Project['status'];
}

export interface CreateTaskInput {
  projectId: string;
  name: string;
}

class ProjectService {
  /** Ensure the "General" project exists on first load. */
  ensureGeneralProject(): void {
    const projects = this.getAll();
    if (!projects.some(p => p.id === GENERAL_PROJECT_ID)) {
      const general: Project = {
        id: GENERAL_PROJECT_ID,
        name: 'General',
        description: 'Default project for unassigned time entries',
        color: '#8C83A8',
        status: 'active',
        createdAt: new Date().toISOString(),
      };
      projects.push(general);
      store.set(PROJECTS_KEY, projects);
    }
  }

  getAll(): Project[] {
    return store.get<Project>(PROJECTS_KEY);
  }

  getById(id: string): Project | undefined {
    return this.getAll().find(p => p.id === id);
  }

  getActive(): Project[] {
    return this.getAll().filter(p => p.status === 'active');
  }

  /** Get projects visible in the clock-in selector (excludes archived). */
  getSelectable(): Project[] {
    return this.getAll().filter(p => p.status !== 'archived');
  }

  create(input: CreateProjectInput): Project {
    const reqErr = validateRequired(input.name, 'Project name');
    if (reqErr) throw new Error(reqErr);

    const maxErr = validateMaxLength(input.name, 100, 'Project name');
    if (maxErr) throw new Error(maxErr);

    const projects = this.getAll();

    const project: Project = {
      id: generateId(),
      name: sanitizeString(input.name),
      description: sanitizeString(input.description || ''),
      color: input.color || '#5746B2',
      status: input.status || 'active',
      createdAt: new Date().toISOString(),
    };

    projects.push(project);
    store.set(PROJECTS_KEY, projects);
    eventBus.emit('projects:changed', project);
    eventBus.emit('projects:created', project);

    return project;
  }

  update(id: string, input: UpdateProjectInput): Project {
    const projects = this.getAll();
    const index = projects.findIndex(p => p.id === id);
    if (index === -1) throw new Error('Project not found');

    if (input.name !== undefined) {
      const reqErr = validateRequired(input.name, 'Project name');
      if (reqErr) throw new Error(reqErr);

      const maxErr = validateMaxLength(input.name, 100, 'Project name');
      if (maxErr) throw new Error(maxErr);
    }

    // Validate status transitions
    if (input.status !== undefined) {
      const current = projects[index].status;
      const validTransitions: Record<string, string[]> = {
        active: ['completed', 'archived'],
        completed: ['active', 'archived'],
        archived: ['active'],
      };
      if (!validTransitions[current]?.includes(input.status)) {
        throw new Error(`Cannot transition from ${current} to ${input.status}`);
      }
    }

    const updated: Project = {
      ...projects[index],
      ...(input.name !== undefined && { name: sanitizeString(input.name) }),
      ...(input.description !== undefined && { description: sanitizeString(input.description) }),
      ...(input.color !== undefined && { color: input.color }),
      ...(input.status !== undefined && { status: input.status }),
    };

    projects[index] = updated;
    store.set(PROJECTS_KEY, projects);
    eventBus.emit('projects:changed', updated);
    eventBus.emit('projects:updated', updated);

    return updated;
  }

  delete(id: string): void {
    if (id === GENERAL_PROJECT_ID) {
      throw new Error('Cannot delete the General project');
    }

    const projects = this.getAll();
    const filtered = projects.filter(p => p.id !== id);
    if (filtered.length === projects.length) throw new Error('Project not found');

    store.set(PROJECTS_KEY, filtered);

    // Also delete tasks for this project
    const tasks = this.getAllTasks().filter(t => t.projectId !== id);
    store.set(TASKS_KEY, tasks);

    eventBus.emit('projects:changed');
    eventBus.emit('projects:deleted', id);
  }

  // ---- Tasks ----

  getAllTasks(): Task[] {
    return store.get<Task>(TASKS_KEY);
  }

  getTasksByProject(projectId: string): Task[] {
    return this.getAllTasks().filter(t => t.projectId === projectId);
  }

  getTaskById(id: string): Task | undefined {
    return this.getAllTasks().find(t => t.id === id);
  }

  createTask(input: CreateTaskInput): Task {
    const reqErr = validateRequired(input.name, 'Task name');
    if (reqErr) throw new Error(reqErr);

    const maxErr = validateMaxLength(input.name, 100, 'Task name');
    if (maxErr) throw new Error(maxErr);

    const tasks = this.getAllTasks();

    const task: Task = {
      id: generateId(),
      projectId: input.projectId,
      name: sanitizeString(input.name),
      status: 'todo',
    };

    tasks.push(task);
    store.set(TASKS_KEY, tasks);
    eventBus.emit('tasks:changed', task);

    return task;
  }

  updateTaskStatus(id: string, status: Task['status']): Task {
    const tasks = this.getAllTasks();
    const index = tasks.findIndex(t => t.id === id);
    if (index === -1) throw new Error('Task not found');

    tasks[index] = { ...tasks[index], status };
    store.set(TASKS_KEY, tasks);
    eventBus.emit('tasks:changed', tasks[index]);

    return tasks[index];
  }

  deleteTask(id: string): void {
    const tasks = this.getAllTasks();
    const filtered = tasks.filter(t => t.id !== id);
    if (filtered.length === tasks.length) throw new Error('Task not found');

    store.set(TASKS_KEY, filtered);
    eventBus.emit('tasks:changed');
  }
}

export const projectService = new ProjectService();
