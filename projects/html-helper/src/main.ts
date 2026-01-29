import './styles/main.css';
import { LandingPage } from './components/LandingPage';
import { Navigation } from './components/Navigation';
import { LessonView } from './components/LessonView';
import { getLessonById, lessons } from './data/lessons';
import { storage } from './utils/storage';

class App {
  private appContainer: HTMLElement;
  private currentLessonId: number = 1;
  private currentLessonView: LessonView | null = null;

  constructor() {
    const container = document.getElementById('app');
    if (!container) {
      throw new Error('App container not found');
    }
    this.appContainer = container;
  }

  async init(): Promise<void> {
    // Update streak on app start
    await storage.updateStreak();

    // Check if user has progress
    const progress = await storage.getProgress();

    if (progress.completedLessons.length > 0) {
      // Returning user - go straight to learning
      this.currentLessonId = progress.currentLesson;
      this.showLearningView();
    } else {
      // New user - show landing page
      this.showLandingPage();
    }

    // Set up event listeners
    this.setupEventListeners();
  }

  private showLandingPage(): void {
    this.appContainer.innerHTML = '<div id="landing-container"></div>';

    const container = document.getElementById('landing-container');
    if (container) {
      const landingPage = new LandingPage(container);
      landingPage.render();
    }
  }

  private async showLearningView(): Promise<void> {
    // Clean up previous lesson view
    if (this.currentLessonView) {
      this.currentLessonView.destroy();
    }

    // Render structure
    this.appContainer.innerHTML = `
      <div class="h-full flex flex-col">
        <div id="nav-container"></div>
        <div id="lesson-container" class="flex-1 overflow-auto"></div>
      </div>
    `;

    // Render navigation
    const navContainer = document.getElementById('nav-container');
    if (navContainer) {
      const nav = new Navigation(navContainer);
      await nav.render();
    }

    // Render lesson
    await this.loadLesson(this.currentLessonId);
  }

  private async loadLesson(lessonId: number): Promise<void> {
    const lesson = getLessonById(lessonId);
    if (!lesson) {
      console.error('Lesson not found:', lessonId);
      return;
    }

    const lessonContainer = document.getElementById('lesson-container');
    if (!lessonContainer) return;

    this.currentLessonId = lessonId;
    this.currentLessonView = new LessonView(lessonContainer, lesson);
    await this.currentLessonView.render();
  }

  private setupEventListeners(): void {
    // Start learning from landing page
    window.addEventListener('start-learning', () => {
      this.showLearningView();
    });

    // Navigate between lessons
    window.addEventListener('navigate-lesson', async (event: Event) => {
      const customEvent = event as CustomEvent;
      const direction = customEvent.detail?.direction;

      if (direction === 'next' && this.currentLessonId < lessons.length) {
        await this.loadLesson(this.currentLessonId + 1);
      } else if (direction === 'prev' && this.currentLessonId > 1) {
        await this.loadLesson(this.currentLessonId - 1);
      }
    });

    // Open menu/lesson map
    window.addEventListener('open-menu', () => {
      this.showLessonMap();
    });
  }

  private showLessonMap(): void {
    // Create lesson map modal
    const modal = document.createElement('div');
    modal.className = 'modal-backdrop';
    modal.innerHTML = `
      <div class="modal-content max-w-4xl">
        <div class="flex items-center justify-between mb-6">
          <h2 class="text-3xl font-display font-bold text-gradient-primary">
            Your Learning Path
          </h2>
          <button id="close-map" class="text-text-secondary-dark hover:text-text-primary-dark text-2xl">
            ×
          </button>
        </div>

        <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-4 mb-6">
          ${this.renderLessonNodes()}
        </div>

        <div class="text-center">
          <button id="resume-learning" class="btn-primary">
            Continue Learning →
          </button>
        </div>
      </div>
    `;

    document.body.appendChild(modal);

    // Event listeners
    const closeBtn = modal.querySelector('#close-map');
    const resumeBtn = modal.querySelector('#resume-learning');

    if (closeBtn) {
      closeBtn.addEventListener('click', () => {
        document.body.removeChild(modal);
      });
    }

    if (resumeBtn) {
      resumeBtn.addEventListener('click', () => {
        document.body.removeChild(modal);
      });
    }

    // Lesson node clicks
    modal.querySelectorAll('.lesson-node').forEach((node) => {
      node.addEventListener('click', async () => {
        const lessonId = parseInt(node.getAttribute('data-lesson-id') || '1');
        const progress = await storage.getProgress();

        // Check if lesson is unlocked
        if (lessonId === 1 || progress.completedLessons.includes(lessonId - 1)) {
          await this.loadLesson(lessonId);
          document.body.removeChild(modal);
        }
      });
    });
  }

  private renderLessonNodes(): string {
    return lessons.map((lesson) => {
      const isCompleted = false; // Will be updated dynamically
      const isCurrent = lesson.id === this.currentLessonId;
      const isLocked = lesson.id > 1 && !isCompleted;

      return `
        <div
          class="lesson-node card-interactive text-center cursor-pointer ${isLocked ? 'opacity-50' : ''}"
          data-lesson-id="${lesson.id}"
        >
          <div class="text-3xl mb-2">
            ${isCompleted ? '✓' : isCurrent ? '🔥' : isLocked ? '🔒' : '○'}
          </div>
          <div class="text-sm font-semibold text-text-primary-dark">
            Lesson ${lesson.id}
          </div>
          <div class="text-xs text-text-secondary-dark mt-1">
            ${lesson.title.split(':')[0]}
          </div>
        </div>
      `;
    }).join('');
  }
}

// Initialize app
const app = new App();
app.init().catch((error) => {
  console.error('Failed to initialize app:', error);
  document.body.innerHTML = `
    <div class="min-h-screen flex items-center justify-center bg-background-dark text-text-primary-dark">
      <div class="text-center">
        <h1 class="text-4xl font-bold mb-4">Oops!</h1>
        <p>Something went wrong. Please refresh the page.</p>
      </div>
    </div>
  `;
});
