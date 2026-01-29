import { sanitizer } from './sanitizer';

/**
 * Sandboxed iframe manager for safe HTML execution
 * Prevents XSS and malicious code execution
 */
export class SandboxManager {
  private iframe: HTMLIFrameElement | null = null;
  private container: HTMLElement;
  private executionTimeout: number = 5000; // 5 seconds
  private executionCount: number = 0;
  private rateLimitWindow: number = 60000; // 1 minute
  private maxExecutionsPerWindow: number = 50;
  private timeoutHandle: number | null = null;

  constructor(container: HTMLElement) {
    this.container = container;
    this.resetRateLimitWindow();
  }

  /**
   * Initialize sandboxed iframe
   */
  private createIframe(): HTMLIFrameElement {
    const iframe = document.createElement('iframe');

    // Security sandbox attributes
    iframe.setAttribute('sandbox', 'allow-scripts');
    iframe.setAttribute('width', '100%');
    iframe.setAttribute('height', '100%');
    iframe.style.border = 'none';
    iframe.style.background = 'white';

    return iframe;
  }

  /**
   * Execute HTML in sandboxed environment
   */
  async execute(html: string): Promise<void> {
    // Rate limiting check
    if (!this.checkRateLimit()) {
      throw new Error('Rate limit exceeded. Please wait a moment before trying again.');
    }

    // Sanitize HTML
    const sanitizedHTML = sanitizer.sanitize(html);

    // Destroy existing iframe
    this.destroy();

    // Create new iframe
    this.iframe = this.createIframe();
    this.container.appendChild(this.iframe);

    // Set timeout for execution
    this.timeoutHandle = window.setTimeout(() => {
      this.destroy();
      throw new Error('Execution timeout. Your code took too long to run.');
    }, this.executionTimeout);

    // Write sanitized HTML to iframe
    try {
      const iframeDoc = this.iframe.contentDocument || this.iframe.contentWindow?.document;
      if (iframeDoc) {
        iframeDoc.open();
        iframeDoc.write(sanitizedHTML);
        iframeDoc.close();
      }

      // Clear timeout on success
      if (this.timeoutHandle) {
        clearTimeout(this.timeoutHandle);
        this.timeoutHandle = null;
      }
    } catch (error) {
      this.destroy();
      throw new Error('Failed to execute HTML. Please check your code.');
    }
  }

  /**
   * Check rate limit
   */
  private checkRateLimit(): boolean {
    this.executionCount++;
    return this.executionCount <= this.maxExecutionsPerWindow;
  }

  /**
   * Reset rate limit window
   */
  private resetRateLimitWindow(): void {
    setInterval(() => {
      this.executionCount = 0;
    }, this.rateLimitWindow);
  }

  /**
   * Destroy iframe and cleanup
   */
  destroy(): void {
    if (this.timeoutHandle) {
      clearTimeout(this.timeoutHandle);
      this.timeoutHandle = null;
    }

    if (this.iframe && this.iframe.parentNode) {
      this.iframe.parentNode.removeChild(this.iframe);
      this.iframe = null;
    }
  }

  /**
   * Get iframe element
   */
  getIframe(): HTMLIFrameElement | null {
    return this.iframe;
  }
}
