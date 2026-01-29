import anime from 'animejs';

/**
 * Animation utilities using Anime.js
 * Provides delightful animations for user interactions
 */
export class Animator {
  private prefersReducedMotion: boolean;

  constructor() {
    this.prefersReducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches;
  }

  /**
   * Check if animations should be disabled
   */
  private shouldAnimate(): boolean {
    return !this.prefersReducedMotion;
  }

  /**
   * Button click animation
   */
  buttonClick(element: HTMLElement): void {
    if (!this.shouldAnimate()) return;

    anime({
      targets: element,
      scale: [1, 0.95, 1],
      duration: 300,
      easing: 'spring(1, 80, 10, 0)',
    });
  }

  /**
   * Success celebration with confetti
   */
  successCelebration(element: HTMLElement): void {
    if (!this.shouldAnimate()) return;

    // Element bounce
    anime({
      targets: element,
      scale: [1, 1.1, 1],
      rotate: [0, 5, -5, 0],
      duration: 600,
      easing: 'easeInOutQuad',
    });

    // Trigger confetti
    this.createConfetti();
  }

  /**
   * Error shake animation
   */
  errorShake(element: HTMLElement): void {
    if (!this.shouldAnimate()) return;

    anime({
      targets: element,
      translateX: [0, -10, 10, -10, 10, 0],
      duration: 400,
      easing: 'easeInOutSine',
    });
  }

  /**
   * Bounce in animation
   */
  bounceIn(element: HTMLElement, delay: number = 0): void {
    if (!this.shouldAnimate()) {
      element.style.opacity = '1';
      return;
    }

    anime({
      targets: element,
      scale: [0.3, 1],
      opacity: [0, 1],
      duration: 600,
      delay,
      easing: 'easeOutElastic(1, .8)',
    });
  }

  /**
   * Fade in animation
   */
  fadeIn(element: HTMLElement, duration: number = 300): void {
    if (!this.shouldAnimate()) {
      element.style.opacity = '1';
      return;
    }

    anime({
      targets: element,
      opacity: [0, 1],
      duration,
      easing: 'easeOutQuad',
    });
  }

  /**
   * Slide in from left
   */
  slideInLeft(element: HTMLElement): void {
    if (!this.shouldAnimate()) {
      element.style.transform = 'translateX(0)';
      return;
    }

    anime({
      targets: element,
      translateX: [-100, 0],
      opacity: [0, 1],
      duration: 400,
      easing: 'easeOutCubic',
    });
  }

  /**
   * Pulse glow animation
   */
  pulseGlow(element: HTMLElement): void {
    if (!this.shouldAnimate()) return;

    anime({
      targets: element,
      scale: [1, 1.05, 1],
      duration: 2000,
      loop: true,
      easing: 'easeInOutSine',
    });
  }

  /**
   * Count up animation for numbers
   */
  countUp(element: HTMLElement, from: number, to: number, duration: number = 500): void {
    const obj = { value: from };

    anime({
      targets: obj,
      value: to,
      duration,
      easing: 'easeOutExpo',
      round: 1,
      update: () => {
        element.textContent = `+${Math.round(obj.value)} XP`;
      },
    });
  }

  /**
   * Create confetti particles
   */
  private createConfetti(): void {
    const canvas = document.getElementById('confetti-canvas') as HTMLCanvasElement;
    if (!canvas) {
      // Create canvas if doesn't exist
      const newCanvas = document.createElement('canvas');
      newCanvas.id = 'confetti-canvas';
      newCanvas.style.position = 'fixed';
      newCanvas.style.top = '0';
      newCanvas.style.left = '0';
      newCanvas.style.width = '100vw';
      newCanvas.style.height = '100vh';
      newCanvas.style.pointerEvents = 'none';
      newCanvas.style.zIndex = '9999';
      document.body.appendChild(newCanvas);
      this.renderConfetti(newCanvas);
    } else {
      this.renderConfetti(canvas);
    }
  }

  /**
   * Render confetti on canvas
   */
  private renderConfetti(canvas: HTMLCanvasElement): void {
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    canvas.width = window.innerWidth;
    canvas.height = window.innerHeight;

    const colors = ['#FF2E63', '#00D9FF', '#00FF9D', '#FFB627', '#FF5C85'];
    const particleCount = 50;
    const particles: Array<{
      x: number;
      y: number;
      size: number;
      color: string;
      velocityX: number;
      velocityY: number;
      rotation: number;
      rotationSpeed: number;
      opacity: number;
    }> = [];

    // Create particles
    for (let i = 0; i < particleCount; i++) {
      particles.push({
        x: Math.random() * canvas.width,
        y: -20,
        size: Math.random() * 10 + 5,
        color: colors[Math.floor(Math.random() * colors.length)],
        velocityX: (Math.random() - 0.5) * 4,
        velocityY: Math.random() * 3 + 2,
        rotation: Math.random() * 360,
        rotationSpeed: (Math.random() - 0.5) * 10,
        opacity: 1,
      });
    }

    // Animation loop
    const animate = () => {
      ctx.clearRect(0, 0, canvas.width, canvas.height);

      particles.forEach((particle, index) => {
        particle.x += particle.velocityX;
        particle.y += particle.velocityY;
        particle.rotation += particle.rotationSpeed;
        particle.opacity -= 0.01;

        if (particle.opacity <= 0 || particle.y > canvas.height) {
          particles.splice(index, 1);
        }

        ctx.save();
        ctx.translate(particle.x, particle.y);
        ctx.rotate((particle.rotation * Math.PI) / 180);
        ctx.globalAlpha = particle.opacity;
        ctx.fillStyle = particle.color;
        ctx.fillRect(-particle.size / 2, -particle.size / 2, particle.size, particle.size);
        ctx.restore();
      });

      if (particles.length > 0) {
        requestAnimationFrame(animate);
      } else {
        // Remove canvas when done
        if (canvas.parentNode) {
          canvas.parentNode.removeChild(canvas);
        }
      }
    };

    animate();
  }
}

export const animator = new Animator();
