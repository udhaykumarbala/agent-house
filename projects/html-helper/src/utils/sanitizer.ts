import DOMPurify from 'dompurify';

/**
 * Sanitizes HTML to prevent XSS attacks
 * Uses DOMPurify with strict configuration
 */
export class HTMLSanitizer {
  private purify: typeof DOMPurify;

  constructor() {
    this.purify = DOMPurify;
    this.configurePurify();
  }

  private configurePurify(): void {
    // Strict configuration for learning environment
    this.purify.setConfig({
      ALLOWED_TAGS: [
        'html', 'head', 'body', 'title', 'meta',
        'h1', 'h2', 'h3', 'h4', 'h5', 'h6',
        'p', 'br', 'hr', 'strong', 'em', 'b', 'i', 'u',
        'ul', 'ol', 'li',
        'a', 'img',
        'div', 'span',
        'table', 'thead', 'tbody', 'tr', 'th', 'td',
        'form', 'input', 'button', 'textarea', 'select', 'option', 'label'
      ],
      ALLOWED_ATTR: [
        'href', 'src', 'alt', 'title', 'class', 'id', 'style',
        'type', 'name', 'value', 'placeholder', 'rows', 'cols',
        'width', 'height', 'border', 'cellpadding', 'cellspacing'
      ],
      ALLOW_DATA_ATTR: false,
      ALLOW_UNKNOWN_PROTOCOLS: false,
      SAFE_FOR_TEMPLATES: true,
    });
  }

  /**
   * Sanitize HTML string
   */
  sanitize(html: string): string {
    return this.purify.sanitize(html, {
      RETURN_DOM_FRAGMENT: false,
      RETURN_DOM: false,
    }) as string;
  }

  /**
   * Check if HTML is safe (no sanitization changes)
   */
  isSafe(html: string): boolean {
    const sanitized = this.sanitize(html);
    return html.trim() === sanitized.trim();
  }

  /**
   * Validate HTML structure (basic well-formed check)
   */
  isWellFormed(html: string): boolean {
    try {
      const parser = new DOMParser();
      const doc = parser.parseFromString(html, 'text/html');
      const errors = doc.querySelectorAll('parsererror');
      return errors.length === 0;
    } catch {
      return false;
    }
  }
}

export const sanitizer = new HTMLSanitizer();
