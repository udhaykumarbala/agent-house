import { Lesson, validateLessonCode } from '../data/lessons';
import { CodeEditor } from './CodeEditor';
import { SandboxManager } from '../utils/sandbox';
import { animator } from '../utils/animator';
import { storage } from '../utils/storage';
import { checkAchievements, getAchievementById } from '../data/achievements';

export class LessonView {
  private container: HTMLElement;
  private lesson: Lesson;
  private codeEditor: CodeEditor | null = null;
  private sandbox: SandboxManager | null = null;
  private currentHintIndex: number = 0;

  constructor(container: HTMLElement, lesson: Lesson) {
    this.container = container;
    this.lesson = lesson;
  }

  async render(): Promise<void> {
    this.container.innerHTML = `
      <div class="flex flex-col h-full lg:flex-row">
        <!-- Left Panel: Challenge + Code -->
        <div class="w-full lg:w-1/2 flex flex-col border-b lg:border-b-0 lg:border-r border-border-dark">
          <!-- Challenge Section -->
          <div class="p-6 border-b border-border-dark bg-surface-elevated">
            <h2 class="text-2xl font-display font-bold mb-2 text-gradient-primary">
              ${this.lesson.title}
            </h2>
            <p class="text-text-secondary-dark mb-4">
              ${this.lesson.description}
            </p>

            <div class="bg-background-dark rounded-lg p-4 mb-4">
              <p class="text-sm text-text-secondary-dark mb-2">Theory:</p>
              <p class="text-text-primary-dark leading-relaxed">
                ${this.lesson.theory}
              </p>
            </div>

            <div class="bg-surface-dark border border-primary/20 rounded-lg p-4">
              <p class="text-sm text-secondary mb-2">✨ Your Challenge:</p>
              <p class="text-text-primary-dark font-medium">
                ${this.lesson.challenge}
              </p>
            </div>

            <!-- Hints -->
            <div class="mt-4">
              <button id="hint-btn" class="btn-ghost text-sm">
                💡 Show Hint
              </button>
              <div id="hint-content" class="mt-2 hidden bg-background-dark rounded-lg p-3 text-sm text-text-secondary-dark">
              </div>
            </div>
          </div>

          <!-- Code Editor -->
          <div class="flex-1 flex flex-col min-h-[300px]">
            <div class="px-6 py-3 bg-surface-elevated border-b border-border-dark">
              <p class="text-sm font-semibold text-text-secondary-dark">Your Code:</p>
            </div>
            <div id="code-editor" class="flex-1 overflow-auto"></div>
          </div>

          <!-- Actions -->
          <div class="p-4 border-t border-border-dark bg-surface-dark flex items-center justify-between">
            <button id="prev-btn" class="btn-ghost">
              ← Previous
            </button>
            <button id="check-btn" class="btn-primary">
              Check Answer ✓
            </button>
            <button id="next-btn" class="btn-ghost opacity-50 cursor-not-allowed" disabled>
              Next →
            </button>
          </div>
        </div>

        <!-- Right Panel: Live Preview -->
        <div class="w-full lg:w-1/2 flex flex-col bg-white">
          <div class="px-6 py-3 bg-surface-elevated border-b border-border-dark">
            <p class="text-sm font-semibold text-text-primary-dark">Live Preview:</p>
          </div>
          <div id="preview-container" class="flex-1 overflow-auto"></div>
        </div>
      </div>
    `;

    await this.initializeEditor();
    this.attachEventListeners();
  }

  private async initializeEditor(): Promise<void> {
    const editorContainer = document.getElementById('code-editor');
    const previewContainer = document.getElementById('preview-container');

    if (!editorContainer || !previewContainer) return;

    // Initialize code editor
    this.codeEditor = new CodeEditor(editorContainer);
    this.codeEditor.initialize(this.lesson.starterCode, (code) => {
      this.updatePreview(code);
    });

    // Initialize sandbox
    this.sandbox = new SandboxManager(previewContainer);

    // Initial preview
    if (this.lesson.starterCode) {
      await this.updatePreview(this.lesson.starterCode);
    }

    // Focus editor
    setTimeout(() => this.codeEditor?.focus(), 100);
  }

  private async updatePreview(code: string): Promise<void> {
    if (!this.sandbox) return;

    try {
      await this.sandbox.execute(code);
    } catch (error) {
      console.error('Preview error:', error);
    }
  }

  private attachEventListeners(): void {
    // Hint button
    const hintBtn = document.getElementById('hint-btn');
    const hintContent = document.getElementById('hint-content');

    if (hintBtn && hintContent) {
      hintBtn.addEventListener('click', () => {
        if (this.currentHintIndex < this.lesson.hints.length) {
          hintContent.classList.remove('hidden');
          hintContent.textContent = this.lesson.hints[this.currentHintIndex];
          animator.fadeIn(hintContent);
          this.currentHintIndex++;

          if (this.currentHintIndex >= this.lesson.hints.length) {
            hintBtn.textContent = '💡 No more hints';
            (hintBtn as HTMLButtonElement).disabled = true;
          }
        }
      });
    }

    // Check button
    const checkBtn = document.getElementById('check-btn');
    if (checkBtn) {
      checkBtn.addEventListener('click', () => this.checkAnswer());
    }

    // Navigation buttons
    const prevBtn = document.getElementById('prev-btn');
    const nextBtn = document.getElementById('next-btn');

    if (prevBtn) {
      prevBtn.addEventListener('click', () => {
        const event = new CustomEvent('navigate-lesson', { detail: { direction: 'prev' } });
        window.dispatchEvent(event);
      });
    }

    if (nextBtn) {
      nextBtn.addEventListener('click', () => {
        const event = new CustomEvent('navigate-lesson', { detail: { direction: 'next' } });
        window.dispatchEvent(event);
      });
    }
  }

  private async checkAnswer(): Promise<void> {
    const code = this.codeEditor?.getCode() || '';
    const isValid = validateLessonCode(this.lesson.id, code);

    const checkBtn = document.getElementById('check-btn');
    if (!checkBtn) return;

    if (isValid) {
      // Success!
      animator.successCelebration(checkBtn);

      // Show XP gain
      await this.showXPGain();

      // Save progress
      await storage.completeLesson(this.lesson.id, this.lesson.xpReward);

      // Check for achievements
      const progress = await storage.getProgress();
      const newAchievements = checkAchievements(progress);

      for (const achievementId of newAchievements) {
        await storage.unlockAchievement(achievementId);
        await this.showAchievement(achievementId);
      }

      // Enable next button
      const nextBtn = document.getElementById('next-btn');
      if (nextBtn) {
        nextBtn.classList.remove('opacity-50', 'cursor-not-allowed');
        nextBtn.removeAttribute('disabled');
        animator.pulseGlow(nextBtn);
      }
    } else {
      // Error
      animator.errorShake(checkBtn);
      this.showError(this.lesson.validation.errorMessage || 'Not quite right. Try again!');
    }
  }

  private async showXPGain(): Promise<void> {
    const modal = document.createElement('div');
    modal.className = 'modal-backdrop';
    modal.innerHTML = `
      <div class="modal-content text-center">
        <div class="text-6xl mb-4">🎉</div>
        <h2 class="text-3xl font-display font-bold mb-2 text-gradient-primary">
          Lesson Complete!
        </h2>
        <p class="text-text-secondary-dark mb-4">
          Great job! You've mastered ${this.lesson.title}
        </p>
        <div id="xp-counter" class="text-4xl font-bold text-success mb-6">
          +${this.lesson.xpReward} XP
        </div>
        <button id="continue-btn" class="btn-primary">
          Continue →
        </button>
      </div>
    `;

    document.body.appendChild(modal);
    animator.fadeIn(modal);

    const continueBtn = modal.querySelector('#continue-btn');
    if (continueBtn) {
      continueBtn.addEventListener('click', () => {
        document.body.removeChild(modal);
      });
    }

    // Animate XP counter
    const xpCounter = modal.querySelector('#xp-counter') as HTMLElement;
    if (xpCounter) {
      animator.countUp(xpCounter, 0, this.lesson.xpReward);
    }
  }

  private async showAchievement(achievementId: string): Promise<void> {
    const achievement = getAchievementById(achievementId);
    if (!achievement) return;

    const modal = document.createElement('div');
    modal.className = 'modal-backdrop';
    modal.innerHTML = `
      <div class="modal-content text-center">
        <div class="text-6xl mb-4">${achievement.icon}</div>
        <h2 class="text-3xl font-display font-bold mb-2 text-gradient-sunset">
          Achievement Unlocked!
        </h2>
        <h3 class="text-2xl font-bold text-text-primary-dark mb-2">
          ${achievement.title}
        </h3>
        <p class="text-text-secondary-dark mb-6">
          ${achievement.description}
        </p>
        <button id="close-achievement" class="btn-primary">
          Awesome! ✨
        </button>
      </div>
    `;

    document.body.appendChild(modal);
    animator.bounceIn(modal.querySelector('.modal-content') as HTMLElement);

    const closeBtn = modal.querySelector('#close-achievement');
    if (closeBtn) {
      closeBtn.addEventListener('click', () => {
        document.body.removeChild(modal);
      });
    }
  }

  private showError(message: string): void {
    const errorDiv = document.createElement('div');
    errorDiv.className = 'fixed top-20 right-4 bg-error text-white px-6 py-3 rounded-lg shadow-lg z-50';
    errorDiv.textContent = message;

    document.body.appendChild(errorDiv);
    animator.slideInLeft(errorDiv);

    setTimeout(() => {
      if (errorDiv.parentNode) {
        document.body.removeChild(errorDiv);
      }
    }, 3000);
  }

  destroy(): void {
    if (this.codeEditor) {
      this.codeEditor.destroy();
    }
    if (this.sandbox) {
      this.sandbox.destroy();
    }
  }
}
