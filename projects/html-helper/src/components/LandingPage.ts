import { animator } from '../utils/animator';

export class LandingPage {
  private container: HTMLElement;

  constructor(container: HTMLElement) {
    this.container = container;
  }

  render(): void {
    this.container.innerHTML = `
      <div class="min-h-screen flex flex-col items-center justify-center p-6 bg-gradient-to-br from-background-dark via-surface-dark to-background-dark">
        <!-- Hero Section -->
        <div class="max-w-4xl w-full text-center mb-12">
          <div class="text-6xl mb-6 animate-float">✨</div>
          <h1 class="text-5xl md:text-6xl font-display font-bold mb-6 text-gradient-sunset">
            HTML Helper
          </h1>
          <p class="text-xl md:text-2xl text-text-secondary-dark mb-8">
            Learn HTML by playing with it - every tag you type makes something magical happen
          </p>

          <div class="flex flex-wrap items-center justify-center gap-4 mb-12">
            <div class="badge badge-primary text-base px-4 py-2">
              ✓ No signup needed
            </div>
            <div class="badge badge-success text-base px-4 py-2">
              ✓ Free forever
            </div>
            <div class="badge badge-info text-base px-4 py-2">
              ✓ 15 fun lessons
            </div>
          </div>
        </div>

        <!-- Interactive Demo -->
        <div class="max-w-4xl w-full card mb-8">
          <h2 class="text-2xl font-display font-bold mb-4 text-gradient-primary">
            Try it now!
          </h2>
          <p class="text-text-secondary-dark mb-6">
            Type some HTML and watch it come alive 👇
          </p>

          <div class="grid md:grid-cols-2 gap-4">
            <!-- Demo Editor -->
            <div class="bg-background-dark rounded-lg p-4">
              <p class="text-sm text-text-secondary-dark mb-2">Type HTML:</p>
              <textarea
                id="demo-input"
                class="w-full h-32 bg-surface-dark border border-border-dark rounded-md p-3 font-code text-sm text-text-primary-dark resize-none focus:outline-none focus:border-primary"
                placeholder="<h1>Hello, HTML!</h1>"
              ></textarea>
            </div>

            <!-- Demo Preview -->
            <div class="bg-white rounded-lg p-4">
              <p class="text-sm text-gray-600 mb-2">Live Preview:</p>
              <div id="demo-preview" class="h-32 overflow-auto border border-gray-200 rounded-md p-3">
                <p class="text-gray-400 text-sm italic">Your preview appears here...</p>
              </div>
            </div>
          </div>

          <div class="mt-6 text-center">
            <button id="start-learning-btn" class="btn-primary text-lg px-8 py-4">
              Start Learning →
            </button>
          </div>
        </div>

        <!-- Features -->
        <div class="max-w-4xl w-full grid md:grid-cols-3 gap-6 mb-8">
          <div class="card text-center">
            <div class="text-4xl mb-3">🎨</div>
            <h3 class="text-lg font-display font-bold mb-2 text-text-primary-dark">
              Beautiful Animations
            </h3>
            <p class="text-text-secondary-dark text-sm">
              Every action triggers delightful animations that make learning fun
            </p>
          </div>

          <div class="card text-center">
            <div class="text-4xl mb-3">⚡</div>
            <h3 class="text-lg font-display font-bold mb-2 text-text-primary-dark">
              Instant Feedback
            </h3>
            <p class="text-text-secondary-dark text-sm">
              See your code come alive in real-time as you type
            </p>
          </div>

          <div class="card text-center">
            <div class="text-4xl mb-3">🏆</div>
            <h3 class="text-lg font-display font-bold mb-2 text-text-primary-dark">
              Achievements
            </h3>
            <p class="text-text-secondary-dark text-sm">
              Unlock badges and track your progress as you learn
            </p>
          </div>
        </div>

        <!-- Footer -->
        <div class="text-center text-text-tertiary-dark text-sm">
          <p>Made with 💖 for learners everywhere</p>
        </div>
      </div>
    `;

    this.attachEventListeners();
  }

  private attachEventListeners(): void {
    // Demo input
    const demoInput = document.getElementById('demo-input') as HTMLTextAreaElement;
    const demoPreview = document.getElementById('demo-preview');

    if (demoInput && demoPreview) {
      let timeoutId: number;

      demoInput.addEventListener('input', () => {
        clearTimeout(timeoutId);
        timeoutId = window.setTimeout(() => {
          const code = demoInput.value;
          if (code.trim()) {
            demoPreview.innerHTML = code;
            animator.fadeIn(demoPreview);
          } else {
            demoPreview.innerHTML = '<p class="text-gray-400 text-sm italic">Your preview appears here...</p>';
          }
        }, 300);
      });

      // Add example on focus
      demoInput.addEventListener('focus', () => {
        if (!demoInput.value) {
          demoInput.value = '<h1>Hello, HTML!</h1>';
          demoInput.dispatchEvent(new Event('input'));
        }
      });
    }

    // Start button
    const startBtn = document.getElementById('start-learning-btn');
    if (startBtn) {
      startBtn.addEventListener('click', () => {
        animator.buttonClick(startBtn);
        setTimeout(() => {
          const event = new CustomEvent('start-learning');
          window.dispatchEvent(event);
        }, 300);
      });
    }
  }
}
