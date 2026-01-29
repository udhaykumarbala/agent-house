import { storage, UserProgress } from '../utils/storage';

export class Navigation {
  private container: HTMLElement;
  private progress: UserProgress | null = null;

  constructor(container: HTMLElement) {
    this.container = container;
  }

  async render(): Promise<void> {
    this.progress = await storage.getProgress();

    const totalLessons = 15;
    const completedCount = this.progress.completedLessons.length;
    const progressPercent = (completedCount / totalLessons) * 100;

    this.container.innerHTML = `
      <nav class="bg-surface-dark border-b border-border-dark">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div class="flex items-center justify-between h-16">
            <!-- Logo -->
            <div class="flex items-center">
              <h1 class="text-2xl font-display font-bold text-gradient-primary">
                HTML Helper
              </h1>
            </div>

            <!-- Stats -->
            <div class="hidden md:flex items-center space-x-6">
              <div class="flex items-center space-x-2">
                <span class="text-text-secondary-dark">Progress:</span>
                <span class="font-semibold text-text-primary-dark">${completedCount}/${totalLessons}</span>
              </div>
              <div class="flex items-center space-x-2">
                <span class="text-2xl">⭐</span>
                <span class="font-semibold text-text-primary-dark">${this.progress.xp} XP</span>
              </div>
              ${this.progress.streak > 0 ? `
                <div class="flex items-center space-x-2">
                  <span class="text-2xl">🔥</span>
                  <span class="font-semibold text-text-primary-dark">${this.progress.streak} day${this.progress.streak > 1 ? 's' : ''}</span>
                </div>
              ` : ''}
            </div>

            <!-- Menu Button -->
            <button id="menu-btn" class="btn-ghost p-2">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"/>
              </svg>
            </button>
          </div>

          <!-- Progress Bar -->
          <div class="pb-4">
            <div class="progress-bar">
              <div class="progress-fill" style="width: ${progressPercent}%"></div>
            </div>
          </div>
        </div>
      </nav>
    `;

    this.attachEventListeners();
  }

  private attachEventListeners(): void {
    const menuBtn = document.getElementById('menu-btn');
    if (menuBtn) {
      menuBtn.addEventListener('click', () => {
        // Open lesson map or menu
        const event = new CustomEvent('open-menu');
        window.dispatchEvent(event);
      });
    }
  }
}
