/**
 * Projects & Task Management View
 *
 * - Project cards with color tags and status badges
 * - Add/Edit Project modal with color picker
 * - Expand/collapse project to show tasks (accordion)
 * - Inline task add within expanded project
 * - Edit/delete project and task actions
 * - Filter by project status
 * - Total hours logged per project
 * - Empty state CTA
 */

import type { ViewModule } from '../types';
import type { Project, Task } from '../types';
import { projectService } from '../data/projectService';
import { timeEntryService } from '../data/timeEntryService';
import { eventBus } from '../utils/eventBus';
import { formatDuration } from '../utils/dates';
import { Modal } from '../components/modal';
import { Toast } from '../components/toast';
import { createIcons } from 'lucide';

// Preset colors from the data visualization palette
const PROJECT_COLORS = [
  '#CF8A2E', '#5746B2', '#E07B5F', '#2E9E8F',
  '#C75A8A', '#5B8AD4', '#8BA055', '#9B6DB0',
];

type StatusFilter = 'all' | 'active' | 'completed' | 'archived';

let currentFilter: StatusFilter = 'all';
let expandedProjectId: string | null = null;
let containerRef: HTMLElement | null = null;

function getProjectTotalHours(projectId: string): number {
  const entries = timeEntryService.getAll().filter(e => e.projectId === projectId);
  let total = 0;
  for (const entry of entries) {
    if (entry.duration !== null) {
      total += entry.duration;
    }
  }
  return total;
}

function getFilteredProjects(): Project[] {
  const all = projectService.getAll();
  if (currentFilter === 'all') return all;
  return all.filter(p => p.status === currentFilter);
}

function getStatusBadgeClass(status: Project['status']): string {
  switch (status) {
    case 'active': return 'badge-active';
    case 'completed': return 'badge-inactive';
    case 'archived': return 'badge-inactive';
  }
}

function getStatusLabel(status: Project['status']): string {
  return status.charAt(0).toUpperCase() + status.slice(1);
}

function getTaskStatusLabel(status: Task['status']): string {
  switch (status) {
    case 'todo': return 'Todo';
    case 'in_progress': return 'In Progress';
    case 'done': return 'Done';
  }
}

function cycleTaskStatus(current: Task['status']): Task['status'] {
  switch (current) {
    case 'todo': return 'in_progress';
    case 'in_progress': return 'done';
    case 'done': return 'todo';
  }
}

function getTaskStatusBadgeStyle(status: Task['status']): string {
  switch (status) {
    case 'todo': return 'background: rgba(140, 131, 168, 0.12); color: #8C83A8;';
    case 'in_progress': return 'background: rgba(207, 138, 46, 0.12); color: #CF8A2E;';
    case 'done': return 'background: rgba(61, 184, 122, 0.12); color: #3DB87A;';
  }
}

function buildColorPicker(selectedColor: string): HTMLElement {
  const wrapper = document.createElement('div');
  wrapper.style.cssText = 'display: flex; gap: 8px; flex-wrap: wrap;';

  for (const color of PROJECT_COLORS) {
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
  return selected?.getAttribute('data-color') || PROJECT_COLORS[0];
}

function openProjectModal(project?: Project): void {
  const isEdit = !!project;
  const content = document.createElement('div');
  content.innerHTML = `
    <div style="display: flex; flex-direction: column; gap: 16px;">
      <div>
        <label class="form-label">Project Name *</label>
        <input id="project-name" type="text" class="input" maxlength="100" placeholder="e.g. Project Alpha" value="" />
        <div id="name-error" class="form-error" style="display: none;"></div>
      </div>
      <div>
        <label class="form-label">Description</label>
        <input id="project-desc" type="text" class="input" placeholder="Brief description (optional)" value="" />
      </div>
      <div>
        <label class="form-label">Color</label>
        <div id="color-picker-slot"></div>
      </div>
      <div>
        <label class="form-label">Status</label>
        <select id="project-status" class="select w-full">
          <option value="active">Active</option>
          <option value="completed">Completed</option>
          <option value="archived">Archived</option>
        </select>
      </div>
    </div>
  `;

  const nameInput = content.querySelector('#project-name') as HTMLInputElement;
  const descInput = content.querySelector('#project-desc') as HTMLInputElement;
  const statusSelect = content.querySelector('#project-status') as HTMLSelectElement;
  const nameError = content.querySelector('#name-error') as HTMLElement;
  const colorSlot = content.querySelector('#color-picker-slot') as HTMLElement;

  const defaultColor = project?.color || PROJECT_COLORS[0];
  const colorPicker = buildColorPicker(defaultColor);
  colorSlot.appendChild(colorPicker);

  if (isEdit) {
    nameInput.value = project!.name;
    descInput.value = project!.description;
    statusSelect.value = project!.status;

    // Disable invalid status transitions
    if (project!.id === 'general') {
      statusSelect.disabled = true;
    }
  }

  Modal.open({
    title: isEdit ? 'Edit Project' : 'New Project',
    content,
    submitLabel: isEdit ? 'Save' : 'Create',
    onSubmit: () => {
      const name = nameInput.value.trim();
      if (!name) {
        nameError.textContent = 'Project name is required';
        nameError.style.display = 'block';
        nameInput.classList.add('error');
        return;
      }

      const color = getSelectedColor(colorPicker);

      try {
        if (isEdit) {
          const updates: Record<string, unknown> = {
            name,
            description: descInput.value.trim(),
            color,
          };
          if (statusSelect.value !== project!.status) {
            updates.status = statusSelect.value;
          }
          projectService.update(project!.id, updates as Parameters<typeof projectService.update>[1]);
          Toast.success('Project updated');
        } else {
          projectService.create({
            name,
            description: descInput.value.trim(),
            color,
            status: statusSelect.value as Project['status'],
          });
          Toast.success('Project created');
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

function confirmDeleteProject(project: Project): void {
  const content = document.createElement('div');
  const msg = document.createElement('p');
  msg.className = 'text-content-secondary';
  msg.style.fontSize = '14px';
  msg.textContent = `Are you sure you want to delete "${project.name}"? This will also remove all tasks in this project. This cannot be undone.`;
  content.appendChild(msg);

  Modal.open({
    title: 'Delete Project',
    content,
    submitLabel: 'Delete',
    onSubmit: () => {
      try {
        projectService.delete(project.id);
        Toast.success(`"${project.name}" deleted`);
        if (expandedProjectId === project.id) {
          expandedProjectId = null;
        }
        Modal.close();
        renderView();
      } catch (err) {
        Toast.error((err as Error).message);
      }
    },
  });
}

function buildTaskSection(project: Project): HTMLElement {
  const section = document.createElement('div');
  section.style.cssText = 'padding: 16px; border-top: 1px solid rgba(91, 70, 178, 0.10);';

  const tasks = projectService.getTasksByProject(project.id);

  // Header
  const header = document.createElement('div');
  header.style.cssText = 'display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px;';

  const title = document.createElement('h4');
  title.textContent = 'Tasks';
  title.style.cssText = 'font-size: 14px; font-weight: 600; color: #8C83A8; text-transform: uppercase; letter-spacing: 0.05em;';

  const addBtn = document.createElement('button');
  addBtn.type = 'button';
  addBtn.className = 'btn-ghost btn-sm';
  addBtn.innerHTML = '<i data-lucide="plus" class="w-3.5 h-3.5 mr-1 inline-block"></i>Add Task';
  addBtn.addEventListener('click', () => showInlineTaskAdd(project.id, section, tasks.length));

  header.appendChild(title);
  header.appendChild(addBtn);
  section.appendChild(header);

  // Task list
  if (tasks.length === 0) {
    const empty = document.createElement('p');
    empty.style.cssText = 'color: #564F6D; font-size: 13px; text-align: center; padding: 12px 0;';
    empty.textContent = 'No tasks yet. Click "Add Task" to create one.';
    section.appendChild(empty);
  } else {
    const list = document.createElement('div');
    list.style.cssText = 'display: flex; flex-direction: column; gap: 4px;';

    for (const task of tasks) {
      const row = document.createElement('div');
      row.style.cssText = 'display: flex; align-items: center; gap: 10px; padding: 8px 10px; border-radius: 8px; transition: background 150ms ease;';
      row.addEventListener('mouseenter', () => { row.style.background = 'rgba(87, 70, 178, 0.08)'; });
      row.addEventListener('mouseleave', () => { row.style.background = 'transparent'; });

      // Checkbox icon
      const check = document.createElement('span');
      check.style.cssText = `
        width: 18px; height: 18px; border-radius: 4px; display: flex;
        align-items: center; justify-content: center; font-size: 12px;
        cursor: pointer; transition: all 150ms ease; flex-shrink: 0;
        ${task.status === 'done'
          ? 'background: rgba(61, 184, 122, 0.15); color: #3DB87A;'
          : 'border: 1.5px solid rgba(91, 70, 178, 0.30); color: transparent;'}
      `;
      check.textContent = task.status === 'done' ? '\u2713' : '';

      // Name
      const name = document.createElement('span');
      name.style.cssText = `flex: 1; font-size: 14px; color: ${task.status === 'done' ? '#564F6D' : '#ECE9F5'}; ${task.status === 'done' ? 'text-decoration: line-through;' : ''}`;
      name.textContent = task.name;

      // Status badge
      const badge = document.createElement('button');
      badge.type = 'button';
      badge.style.cssText = `
        font-size: 11px; font-weight: 500; padding: 3px 8px;
        border-radius: 9999px; border: none; cursor: pointer;
        transition: all 150ms ease; ${getTaskStatusBadgeStyle(task.status)}
      `;
      badge.textContent = getTaskStatusLabel(task.status);
      badge.setAttribute('aria-label', `Change status of ${task.name}`);
      badge.addEventListener('click', () => {
        const newStatus = cycleTaskStatus(task.status);
        projectService.updateTaskStatus(task.id, newStatus);
        renderView();
      });

      // Also toggle on checkbox click
      check.style.cursor = 'pointer';
      check.addEventListener('click', () => {
        const newStatus = cycleTaskStatus(task.status);
        projectService.updateTaskStatus(task.id, newStatus);
        renderView();
      });

      // Delete button
      const delBtn = document.createElement('button');
      delBtn.type = 'button';
      delBtn.style.cssText = 'background: none; border: none; color: #564F6D; cursor: pointer; font-size: 14px; padding: 2px 4px; opacity: 0; transition: opacity 150ms ease, color 150ms ease;';
      delBtn.textContent = '\u2715';
      delBtn.setAttribute('aria-label', `Delete task ${task.name}`);
      delBtn.addEventListener('click', () => {
        projectService.deleteTask(task.id);
        Toast.success('Task deleted');
        renderView();
      });

      row.addEventListener('mouseenter', () => { delBtn.style.opacity = '1'; });
      row.addEventListener('mouseleave', () => { delBtn.style.opacity = '0'; });

      row.appendChild(check);
      row.appendChild(name);
      row.appendChild(badge);
      row.appendChild(delBtn);
      list.appendChild(row);
    }

    section.appendChild(list);
  }

  return section;
}

function showInlineTaskAdd(projectId: string, section: HTMLElement, _taskCount: number): void {
  // Don't add another input if one exists
  if (section.querySelector('.inline-task-input')) return;

  const wrapper = document.createElement('div');
  wrapper.className = 'inline-task-input';
  wrapper.style.cssText = 'display: flex; gap: 8px; margin-top: 8px;';

  const input = document.createElement('input');
  input.type = 'text';
  input.className = 'input';
  input.placeholder = 'Task name...';
  input.maxLength = 100;
  input.style.flex = '1';

  const saveBtn = document.createElement('button');
  saveBtn.type = 'button';
  saveBtn.className = 'btn-primary btn-sm';
  saveBtn.textContent = 'Add';

  const cancelBtn = document.createElement('button');
  cancelBtn.type = 'button';
  cancelBtn.className = 'btn-ghost btn-sm';
  cancelBtn.textContent = 'Cancel';
  cancelBtn.addEventListener('click', () => wrapper.remove());

  const doAdd = () => {
    const name = input.value.trim();
    if (!name) {
      input.classList.add('error');
      return;
    }
    try {
      projectService.createTask({ projectId, name });
      Toast.success('Task added');
      renderView();
    } catch (err) {
      Toast.error((err as Error).message);
    }
  };

  saveBtn.addEventListener('click', doAdd);
  input.addEventListener('keydown', (e) => {
    if (e.key === 'Enter') doAdd();
    if (e.key === 'Escape') wrapper.remove();
  });

  wrapper.appendChild(input);
  wrapper.appendChild(saveBtn);
  wrapper.appendChild(cancelBtn);
  section.appendChild(wrapper);

  input.focus();
}

function buildProjectCard(project: Project): HTMLElement {
  const tasks = projectService.getTasksByProject(project.id);
  const totalMinutes = getProjectTotalHours(project.id);
  const isExpanded = expandedProjectId === project.id;

  const card = document.createElement('div');
  card.className = 'card card-interactive';
  card.style.cssText = `border-left: 4px solid ${project.color}; padding: 0; overflow: hidden;`;

  // Main card content
  const main = document.createElement('div');
  main.style.cssText = 'padding: 16px 16px 16px 12px; cursor: pointer;';
  main.addEventListener('click', () => {
    expandedProjectId = isExpanded ? null : project.id;
    renderView();
  });

  // Top row: name + actions
  const topRow = document.createElement('div');
  topRow.style.cssText = 'display: flex; align-items: flex-start; justify-content: space-between; gap: 8px;';

  const nameWrap = document.createElement('div');
  nameWrap.style.cssText = 'flex: 1; min-width: 0;';

  const nameEl = document.createElement('h3');
  nameEl.style.cssText = 'font-size: 16px; font-weight: 600; color: #ECE9F5; margin: 0 0 4px 0;';
  nameEl.textContent = project.name;

  const descEl = document.createElement('p');
  descEl.style.cssText = 'font-size: 13px; color: #8C83A8; margin: 0; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;';
  descEl.textContent = project.description || '';

  nameWrap.appendChild(nameEl);
  if (project.description) nameWrap.appendChild(descEl);

  // Action buttons
  const actions = document.createElement('div');
  actions.style.cssText = 'display: flex; gap: 4px; opacity: 0; transition: opacity 150ms ease;';
  actions.className = 'project-actions';

  const editBtn = document.createElement('button');
  editBtn.type = 'button';
  editBtn.className = 'action-btn';
  editBtn.innerHTML = '<i data-lucide="pencil" class="w-3.5 h-3.5"></i>';
  editBtn.setAttribute('aria-label', `Edit ${project.name}`);
  editBtn.addEventListener('click', (e) => {
    e.stopPropagation();
    openProjectModal(project);
  });

  if (project.id !== 'general') {
    const delBtn = document.createElement('button');
    delBtn.type = 'button';
    delBtn.className = 'action-btn';
    delBtn.innerHTML = '<i data-lucide="trash-2" class="w-3.5 h-3.5"></i>';
    delBtn.setAttribute('aria-label', `Delete ${project.name}`);
    delBtn.addEventListener('click', (e) => {
      e.stopPropagation();
      confirmDeleteProject(project);
    });
    actions.appendChild(delBtn);
  }

  actions.insertBefore(editBtn, actions.firstChild);

  topRow.appendChild(nameWrap);
  topRow.appendChild(actions);

  // Hover to show actions
  main.addEventListener('mouseenter', () => { actions.style.opacity = '1'; });
  main.addEventListener('mouseleave', () => { actions.style.opacity = '0'; });

  // Bottom row: stats
  const bottomRow = document.createElement('div');
  bottomRow.style.cssText = 'display: flex; align-items: center; gap: 16px; margin-top: 12px;';

  const taskCount = document.createElement('span');
  taskCount.style.cssText = 'font-size: 12px; color: #8C83A8; display: flex; align-items: center; gap: 4px;';
  taskCount.innerHTML = `<i data-lucide="list-checks" class="w-3.5 h-3.5"></i>`;
  const taskCountText = document.createElement('span');
  taskCountText.textContent = `${tasks.length} task${tasks.length !== 1 ? 's' : ''}`;
  taskCount.appendChild(taskCountText);

  const hours = document.createElement('span');
  hours.style.cssText = 'font-size: 12px; color: #8C83A8; display: flex; align-items: center; gap: 4px;';
  hours.innerHTML = `<i data-lucide="clock" class="w-3.5 h-3.5"></i>`;
  const hoursText = document.createElement('span');
  hoursText.textContent = formatDuration(totalMinutes);
  hours.appendChild(hoursText);

  const badge = document.createElement('span');
  badge.className = getStatusBadgeClass(project.status);
  badge.textContent = getStatusLabel(project.status);

  // Expand indicator
  const chevron = document.createElement('span');
  chevron.style.cssText = `
    margin-left: auto; color: #564F6D; font-size: 12px;
    transition: transform 200ms ease;
    ${isExpanded ? 'transform: rotate(180deg);' : ''}
  `;
  chevron.innerHTML = '<i data-lucide="chevron-down" class="w-4 h-4"></i>';

  bottomRow.appendChild(taskCount);
  bottomRow.appendChild(hours);
  bottomRow.appendChild(badge);
  bottomRow.appendChild(chevron);

  main.appendChild(topRow);
  main.appendChild(bottomRow);
  card.appendChild(main);

  // Expanded task section
  if (isExpanded) {
    const taskSection = buildTaskSection(project);
    card.appendChild(taskSection);
  }

  return card;
}

function renderView(): void {
  if (!containerRef) return;
  containerRef.innerHTML = '';

  const projects = getFilteredProjects();
  const allProjects = projectService.getAll();
  const hasProjects = allProjects.length > 1 || (allProjects.length === 1 && allProjects[0].id !== 'general');

  // Header
  const header = document.createElement('div');
  header.style.cssText = 'display: flex; align-items: center; justify-content: space-between; margin-bottom: 24px; flex-wrap: wrap; gap: 12px;';

  const titleWrap = document.createElement('div');
  const h1 = document.createElement('h1');
  h1.className = 'text-content font-heading';
  h1.textContent = 'Projects';
  const subtitle = document.createElement('p');
  subtitle.className = 'text-content-secondary text-sm';
  subtitle.style.marginTop = '4px';
  subtitle.textContent = 'Manage projects and tasks';
  titleWrap.appendChild(h1);
  titleWrap.appendChild(subtitle);

  const addBtn = document.createElement('button');
  addBtn.type = 'button';
  addBtn.className = 'btn-primary';
  addBtn.innerHTML = '<i data-lucide="plus" class="w-4 h-4 mr-1.5 inline-block"></i>New Project';
  addBtn.addEventListener('click', () => openProjectModal());

  header.appendChild(titleWrap);
  header.appendChild(addBtn);
  containerRef.appendChild(header);

  // Filter tabs
  if (hasProjects) {
    const filterBar = document.createElement('div');
    filterBar.style.cssText = 'display: flex; gap: 6px; margin-bottom: 20px; flex-wrap: wrap;';

    const filters: { label: string; value: StatusFilter }[] = [
      { label: 'All', value: 'all' },
      { label: 'Active', value: 'active' },
      { label: 'Completed', value: 'completed' },
      { label: 'Archived', value: 'archived' },
    ];

    for (const f of filters) {
      const btn = document.createElement('button');
      btn.type = 'button';
      btn.className = currentFilter === f.value ? 'date-preset selected' : 'date-preset';
      btn.textContent = f.label;
      btn.addEventListener('click', () => {
        currentFilter = f.value;
        renderView();
      });
      filterBar.appendChild(btn);
    }

    containerRef.appendChild(filterBar);
  }

  // Empty state
  if (!hasProjects) {
    const empty = document.createElement('div');
    empty.className = 'empty-state';
    empty.innerHTML = `
      <div class="empty-icon">
        <i data-lucide="folder-kanban" class="w-12 h-12 mx-auto text-content-disabled"></i>
      </div>
      <div class="empty-title">Create your first project</div>
      <p class="empty-description">Create a project to organize your team's work.</p>
    `;
    const ctaBtn = document.createElement('button');
    ctaBtn.type = 'button';
    ctaBtn.className = 'btn-accent';
    ctaBtn.innerHTML = '<i data-lucide="plus" class="w-4 h-4 mr-1.5 inline-block"></i>New Project';
    ctaBtn.addEventListener('click', () => openProjectModal());
    empty.appendChild(ctaBtn);
    containerRef.appendChild(empty);
    try { createIcons(); } catch { /* ok */ }
    return;
  }

  // No results for filter
  if (projects.length === 0) {
    const noResults = document.createElement('div');
    noResults.className = 'empty-state';
    noResults.innerHTML = `
      <div class="empty-icon">
        <i data-lucide="search" class="w-12 h-12 mx-auto text-content-disabled"></i>
      </div>
      <div class="empty-title">No ${currentFilter} projects</div>
      <p class="empty-description">No projects match the current filter.</p>
    `;
    containerRef.appendChild(noResults);
    try { createIcons(); } catch { /* ok */ }
    return;
  }

  // Project grid
  const grid = document.createElement('div');
  grid.style.cssText = 'display: grid; grid-template-columns: repeat(auto-fill, minmax(340px, 1fr)); gap: 16px;';

  for (const project of projects) {
    grid.appendChild(buildProjectCard(project));
  }

  containerRef.appendChild(grid);

  // Re-init lucide icons for dynamically created elements
  try { createIcons(); } catch { /* ok */ }
}

function onDataChange(): void {
  renderView();
}

export const render: ViewModule['render'] = (container) => {
  containerRef = container;
  currentFilter = 'all';

  eventBus.on('projects:changed', onDataChange);
  eventBus.on('tasks:changed', onDataChange);
  eventBus.on('timeEntries:changed', onDataChange);

  renderView();
};

export const destroy = (): void => {
  eventBus.off('projects:changed', onDataChange);
  eventBus.off('tasks:changed', onDataChange);
  eventBus.off('timeEntries:changed', onDataChange);
  containerRef = null;
  expandedProjectId = null;
};
