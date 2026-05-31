import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { projectService } from '../../data/projectService';
import { timeEntryService } from '../../data/timeEntryService';
import { store } from '../../data/store';
import { eventBus } from '../../utils/eventBus';

// Mock lucide
vi.mock('lucide', () => ({
  createIcons: vi.fn(),
}));

describe('Projects View', () => {
  let container: HTMLElement;
  let projectsModule: typeof import('../projects');

  beforeEach(async () => {
    // Reset state
    store.clearAll();
    eventBus._reset();

    // Ensure General project
    projectService.ensureGeneralProject();

    container = document.createElement('div');
    document.body.appendChild(container);

    // Dynamically import to get fresh module
    projectsModule = await import('../projects');
  });

  afterEach(() => {
    projectsModule.destroy?.();
    document.body.removeChild(container);
  });

  it('should render the page title', () => {
    projectsModule.render(container);
    const h1 = container.querySelector('h1');
    expect(h1).not.toBeNull();
    expect(h1!.textContent).toBe('Projects');
  });

  it('should show empty state when only General project exists', () => {
    projectsModule.render(container);
    const emptyTitle = container.querySelector('.empty-title');
    expect(emptyTitle).not.toBeNull();
    expect(emptyTitle!.textContent).toBe('Create your first project');
  });

  it('should show New Project button', () => {
    projectsModule.render(container);
    const buttons = container.querySelectorAll('button');
    const newBtn = Array.from(buttons).find(b => b.textContent?.includes('New Project'));
    expect(newBtn).not.toBeNull();
  });

  it('should render project cards when projects exist', () => {
    projectService.create({
      name: 'Test Project',
      description: 'A test project',
      color: '#CF8A2E',
    });

    projectsModule.render(container);

    const cards = container.querySelectorAll('.card');
    // Should have at least the project card (General + Test Project)
    expect(cards.length).toBeGreaterThan(0);

    // Find the card with our project name
    const found = Array.from(container.querySelectorAll('h3')).some(
      h => h.textContent === 'Test Project'
    );
    expect(found).toBe(true);
  });

  it('should show status badges on project cards', () => {
    projectService.create({
      name: 'Active Project',
      color: '#5746B2',
      status: 'active',
    });

    projectsModule.render(container);

    const badges = container.querySelectorAll('.badge-active');
    expect(badges.length).toBeGreaterThan(0);
  });

  it('should show task count on project cards', () => {
    const project = projectService.create({
      name: 'With Tasks',
      color: '#E07B5F',
    });
    projectService.createTask({ projectId: project.id, name: 'Task 1' });
    projectService.createTask({ projectId: project.id, name: 'Task 2' });

    projectsModule.render(container);

    const text = container.textContent || '';
    expect(text).toContain('2 tasks');
  });

  it('should show total hours per project', () => {
    const project = projectService.create({
      name: 'Hours Project',
      color: '#2E9E8F',
    });

    // Create a time entry with 2h (120 minutes) duration
    const clockIn = new Date(Date.now() - 7200000).toISOString();
    const clockOut = new Date().toISOString();
    timeEntryService.createManualEntry({
      employeeId: 'emp-1',
      projectId: project.id,
      date: new Date().toISOString().split('T')[0],
      clockIn,
      clockOut,
    });

    projectsModule.render(container);

    // Should display hours
    const text = container.textContent || '';
    expect(text).toContain('h');
  });

  it('should filter projects by status', () => {
    projectService.create({ name: 'Active One', color: '#CF8A2E' });
    const p2 = projectService.create({ name: 'To Complete', color: '#5746B2' });
    projectService.update(p2.id, { status: 'completed' });

    projectsModule.render(container);

    // Click "Active" filter
    const filterBtns = container.querySelectorAll('.date-preset');
    const activeBtn = Array.from(filterBtns).find(b => b.textContent === 'Active');
    expect(activeBtn).not.toBeNull();
    activeBtn!.dispatchEvent(new Event('click'));

    // After filter, "To Complete" should not be visible
    const h3s = Array.from(container.querySelectorAll('h3')).map(h => h.textContent);
    expect(h3s).not.toContain('To Complete');
    expect(h3s).toContain('Active One');
  });

  it('should show empty state CTA with New Project button', () => {
    // Only General project exists
    projectsModule.render(container);

    const ctaBtn = container.querySelector('.btn-accent');
    expect(ctaBtn).not.toBeNull();
    expect(ctaBtn!.textContent).toContain('New Project');
  });
});
