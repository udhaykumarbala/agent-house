# Architecture Specification: HTML Learning Platform

**Project**: HTML Helper - Interactive HTML Learning with Animations
**Template**: `static-enhanced` (Vite + Tailwind + TypeScript)
**Architect**: Software Architect
**Date**: 2026-01-29
**Status**: READY_FOR_DEVELOPMENT

---

## 1. Executive Summary

### Technical Vision
Build an animation-first, interactive HTML learning platform that makes coding feel like play. The architecture prioritizes **security**, **performance**, and **delightful user experience** through sandboxed code execution, instant visual feedback, and carefully orchestrated animations.

### Key Architectural Decisions
1. **Pure Frontend Architecture**: No backend required - all content and logic client-side
2. **Sandboxed Execution**: User HTML code runs in isolated iframes with strict CSP
3. **Animation-Driven Feedback**: Every user action triggers purposeful animations
4. **Progressive Web App**: Offline-capable, fast, app-like experience
5. **Security-First Design**: XSS prevention through multi-layer defense

---

## 2. System Architecture

### 2.1 High-Level Architecture

```
┌─────────────────────────────────────────────────────────┐
│                   User Browser                           │
├─────────────────────────────────────────────────────────┤
│                                                           │
│  ┌────────────────────────────────────────────────┐    │
│  │  Main Application (Secure Context - HTTPS)     │    │
│  │  ├─ Navigation & Progress Component            │    │
│  │  ├─ Lesson Viewer Component                    │    │
│  │  ├─ Code Editor Component                      │    │
│  │  └─ Achievement System                         │    │
│  └────────────────┬───────────────────────────────┘    │
│                   │                                       │
│                   │ (Sanitized HTML via DOMPurify)       │
│                   ▼                                       │
│  ┌────────────────────────────────────────────────┐    │
│  │  Sandboxed iframe (Isolated Execution)         │    │
│  │  - sandbox="allow-scripts"                     │    │
│  │  - No parent window access                     │    │
│  │  - Execution timeout: 5 seconds                │    │
│  │  - CSP: script-src 'self' 'unsafe-inline'      │    │
│  └────────────────────────────────────────────────┘    │
│                                                           │
│  ┌────────────────────────────────────────────────┐    │
│  │  localStorage (User Progress)                   │    │
│  │  - Completed lessons                            │    │
│  │  - Achievements                                 │    │
│  │  - Preferences                                  │    │
│  │  - HMAC integrity checks                        │    │
│  └────────────────────────────────────────────────┘    │
│                                                           │
└─────────────────────────────────────────────────────────┘
```

### 2.2 Component Architecture

```
src/
├── main.ts                    # Entry point, app initialization
├── App.ts                     # Root component orchestrator
│
├── components/
│   ├── Navigation.ts          # Top nav, progress bar
│   ├── LessonViewer.ts        # Lesson content display
│   ├── CodeEditor.ts          # Interactive code editor
│   ├── PreviewPane.ts         # Sandboxed live preview
│   ├── ChallengeCard.ts       # Interactive challenges
│   ├── AchievementPopup.ts    # Achievement notifications
│   └── AnimatedButton.ts      # Reusable animated components
│
├── core/
│   ├── LessonManager.ts       # Lesson loading, navigation
│   ├── ProgressTracker.ts     # Progress persistence
│   ├── Animator.ts            # Animation orchestration
│   ├── Validator.ts           # Code validation
│   ├── Sanitizer.ts           # HTML sanitization wrapper
│   └── EventBus.ts            # Component communication
│
├── security/
│   ├── SandboxManager.ts      # iframe sandbox control
│   ├── CSPPolicy.ts           # CSP configuration
│   ├── InputValidator.ts      # Input validation rules
│   └── IntegrityChecker.ts    # localStorage integrity
│
├── data/
│   ├── lessons.json           # Lesson content
│   ├── achievements.json      # Achievement definitions
│   └── challenges.json        # Challenge specifications
│
├── utils/
│   ├── storage.ts             # localStorage helpers
│   ├── dom.ts                 # DOM utilities
│   └── debounce.ts            # Performance helpers
│
└── styles/
    ├── main.css               # Global styles
    ├── animations.css         # Animation keyframes
    ├── components.css         # Component styles
    └── themes.css             # Color themes
```

---

## 3. Technology Stack

### 3.1 Core Technologies

| Technology | Version | Purpose | Justification |
|------------|---------|---------|---------------|
| **Vite** | 5.x | Build tool | Fast HMR, optimized production builds, plugin ecosystem |
| **TypeScript** | 5.x | Type safety | Complex animation logic needs type safety, better DX |
| **Tailwind CSS** | 3.x | Styling | Rapid UI development, consistent design system |
| **Anime.js** | 3.x | Animations | Lightweight (9KB), perfect for DOM/SVG animations |
| **CodeMirror 6** | 6.x | Code editor | Modern, performant, extensible, mobile-friendly |
| **DOMPurify** | 3.x | HTML sanitization | Industry standard XSS prevention, zero dependencies |

### 3.2 Additional Libraries

| Library | Purpose | Bundle Impact |
|---------|---------|---------------|
| **Prism.js** | Syntax highlighting for static examples | ~2KB (modular) |
| **Phosphor Icons** | Icon set | ~5KB (tree-shakeable) |
| **localForage** | Enhanced localStorage | ~8KB (fallbacks to IndexedDB) |

### 3.3 Development Tools

- **ESLint** + **Prettier**: Code quality and formatting
- **Vitest**: Unit testing (fast, Vite-native)
- **Playwright**: E2E testing (optional for critical flows)
- **npm audit** + **Snyk**: Dependency security scanning

### 3.4 Bundle Size Targets

- **Initial load**: < 150KB (gzipped)
- **Code editor chunk**: ~80KB (lazy loaded)
- **Animation library**: ~9KB
- **Lessons**: Loaded on-demand per lesson (~5KB each)

**Optimization Strategies**:
- Code splitting by route/lesson
- Lazy load CodeMirror and Anime.js
- Tailwind CSS purging in production
- Image optimization (WebP with PNG fallback)
- Service Worker caching for offline support

---

## 4. Data Model

### 4.1 Lesson Structure

**Reference**: `.plans/research/architecture-research.md:195-239`

```typescript
interface Lesson {
  id: string;                    // e.g., "lesson-1-first-tag"
  title: string;                 // "Your First HTML Tag"
  category: LessonCategory;      // "basics" | "text" | "links" | "media"
  order: number;                 // Sequence in category
  duration: string;              // Estimated time "5 min"

  content: {
    introduction: string;        // Brief intro paragraph
    theory: string;              // Concept explanation (markdown)
    examples: Example[];         // Code examples to demonstrate
    keyPoints: string[];         // Bullet points of key takeaways
  };

  challenge: Challenge;          // Interactive coding challenge

  quiz?: Quiz;                   // Optional quiz (post-MVP)

  animations: AnimationTriggers; // When to trigger animations
}

interface Example {
  code: string;                  // HTML code
  description: string;           // What this example shows
  highlight?: number[];          // Lines to highlight
}

interface Challenge {
  type: "code" | "multiple-choice" | "drag-drop";
  instructions: string;          // What the user should do
  starterCode?: string;          // Pre-filled code (if any)

  validation: {
    type: "contains" | "exact" | "regex" | "custom";
    expected: string | string[] | RegExp;
    customValidator?: string;    // Function name for complex validation
  };

  hints?: string[];              // Progressive hints if stuck
  solution?: string;             // For "show solution" button
}

interface AnimationTriggers {
  onStart?: string;              // Animation when lesson starts
  onComplete?: string;           // Animation when challenge passed
  onMilestone?: string;          // Animation for category completion
}
```

### 4.2 User Progress Model

**Reference**: `.plans/research/architecture-research.md:241-258`

```typescript
interface UserProgress {
  completedLessons: string[];    // Array of lesson IDs
  currentLesson: string;         // Current lesson ID
  score: number;                 // Total XP points
  achievements: string[];        // Unlocked achievement IDs
  streakDays: number;            // Consecutive days of learning
  lastVisit: string;             // ISO timestamp

  preferences: UserPreferences;

  version: number;               // Data schema version for migrations
  integrity: string;             // HMAC for tamper detection
}

interface UserPreferences {
  theme: "light" | "dark";
  fontSize: number;              // 12-24
  animationsEnabled: boolean;    // Respect prefers-reduced-motion
  soundEnabled: boolean;         // Success sound effects
}
```

### 4.3 Achievement System

**Reference**: `.plans/research/product-research.md:86-90`

```typescript
interface Achievement {
  id: string;                    // "first-tag"
  title: string;                 // "First Steps"
  description: string;           // "Created your first HTML tag"
  icon: string;                  // Icon name or emoji
  rarity: "common" | "rare" | "epic";
  xp: number;                    // XP reward
  unlockCondition: {
    type: "lesson-complete" | "streak" | "speed" | "perfect-score";
    value: any;                  // Condition-specific data
  };
}
```

---

## 5. Security Architecture

**Reference**: `.plans/research/security_research.md` (entire document)

### 5.1 Defense in Depth Strategy

#### Layer 1: Input Validation
```typescript
// src/security/InputValidator.ts
class InputValidator {
  validateHTML(code: string): ValidationResult {
    // Length check
    if (code.length > 50000) {
      return { valid: false, error: "Code too long (max 50KB)" };
    }

    // Dangerous pattern detection
    const dangerousPatterns = [
      /eval\s*\(/,
      /Function\s*\(/,
      /document\.write/,
      /<base\s/i,
      /javascript:/i
    ];

    for (const pattern of dangerousPatterns) {
      if (pattern.test(code)) {
        return { valid: false, error: "Unsafe code pattern detected" };
      }
    }

    return { valid: true };
  }
}
```

#### Layer 2: HTML Sanitization
```typescript
// src/security/Sanitizer.ts
import DOMPurify from 'dompurify';

class HTMLSanitizer {
  private config = {
    ALLOWED_TAGS: [
      'div', 'span', 'p', 'h1', 'h2', 'h3', 'h4', 'h5', 'h6',
      'strong', 'em', 'u', 's', 'br', 'hr',
      'ul', 'ol', 'li', 'table', 'tr', 'td', 'th',
      'img', 'a', 'button', 'input', 'label', 'form'
    ],
    ALLOWED_ATTR: ['class', 'id', 'style', 'href', 'src', 'alt', 'title'],
    FORBID_TAGS: ['script', 'iframe', 'object', 'embed', 'base'],
    FORBID_ATTR: ['onerror', 'onload', 'onclick', 'onmouseover', 'formaction']
  };

  sanitize(html: string): string {
    return DOMPurify.sanitize(html, this.config);
  }
}
```

#### Layer 3: Sandboxed Execution
```typescript
// src/security/SandboxManager.ts
class SandboxManager {
  private iframe: HTMLIFrameElement;
  private timeout: number = 5000; // 5 seconds

  createSandbox(container: HTMLElement): void {
    this.iframe = document.createElement('iframe');
    this.iframe.sandbox.add('allow-scripts');
    // NO allow-same-origin to prevent parent access
    this.iframe.style.cssText = 'border: none; width: 100%; height: 100%;';
    container.appendChild(this.iframe);
  }

  executeCode(code: string): Promise<void> {
    return new Promise((resolve, reject) => {
      const timeoutId = setTimeout(() => {
        this.killExecution();
        reject(new Error('Execution timeout - possible infinite loop'));
      }, this.timeout);

      this.iframe.srcdoc = code;

      this.iframe.onload = () => {
        clearTimeout(timeoutId);
        resolve();
      };
    });
  }

  killExecution(): void {
    this.iframe.src = 'about:blank';
  }
}
```

#### Layer 4: Content Security Policy
```typescript
// vite.config.ts - CSP headers configuration
export default defineConfig({
  server: {
    headers: {
      'Content-Security-Policy': [
        "default-src 'self'",
        "script-src 'self' 'unsafe-inline' 'unsafe-eval'", // unsafe-eval for CodeMirror
        "style-src 'self' 'unsafe-inline'",
        "img-src 'self' data: https:",
        "frame-src 'self' blob:",
        "connect-src 'self'",
        "object-src 'none'",
        "base-uri 'self'"
      ].join('; ')
    }
  }
});
```

#### Layer 5: Integrity Monitoring
```typescript
// src/security/IntegrityChecker.ts
class IntegrityChecker {
  private secret = 'app-secret-key-from-env'; // In real app, from env

  async generateHMAC(data: string): Promise<string> {
    const encoder = new TextEncoder();
    const key = await crypto.subtle.importKey(
      'raw',
      encoder.encode(this.secret),
      { name: 'HMAC', hash: 'SHA-256' },
      false,
      ['sign']
    );

    const signature = await crypto.subtle.sign(
      'HMAC',
      key,
      encoder.encode(data)
    );

    return Array.from(new Uint8Array(signature))
      .map(b => b.toString(16).padStart(2, '0'))
      .join('');
  }

  async verifyIntegrity(data: string, hash: string): Promise<boolean> {
    const computed = await this.generateHMAC(data);
    return computed === hash;
  }
}
```

### 5.2 Security Checklist (Phase 1 - CRITICAL)

**Reference**: `.plans/research/security_research.md:686-697`

- [x] Sandboxed iframe for code execution
- [x] DOMPurify HTML sanitization
- [x] Content Security Policy headers
- [x] HTTPS enforcement (HSTS)
- [x] Input validation (length, type, patterns)
- [ ] Rate limiting on code execution (50 executions/min)
- [ ] Execution timeout mechanism (5 seconds)
- [ ] Subresource Integrity (SRI) for CDN resources

---

## 6. Animation Architecture

**Reference**: `.plans/research/ui-research.md:130-152` and `.plans/research/ux-research.md:112-125`

### 6.1 Animation Categories

| Category | Purpose | Duration | Easing | Priority |
|----------|---------|----------|--------|----------|
| **Micro-interactions** | Button clicks, hovers | 200ms | ease-out | P0 |
| **Feedback** | Success, error, validation | 300-500ms | bounce/elastic | P0 |
| **Transitions** | Lesson navigation | 300ms | ease-in-out | P1 |
| **Celebrations** | Achievement unlocks | 1-2s | custom | P1 |
| **Ambient** | Idle animations, breathing | 2-3s loop | ease-in-out | P2 |

### 6.2 Animation Manager

```typescript
// src/core/Animator.ts
import anime from 'animejs';

class Animator {
  private animationsEnabled: boolean = true;
  private activeAnimations: anime.AnimeInstance[] = [];

  constructor() {
    // Respect user preferences
    const prefersReduced = window.matchMedia('(prefers-reduced-motion: reduce)');
    this.animationsEnabled = !prefersReduced.matches;
  }

  // Button click feedback
  buttonClick(element: HTMLElement): void {
    if (!this.animationsEnabled) return;

    anime({
      targets: element,
      scale: [1, 0.95, 1],
      duration: 200,
      easing: 'easeOutQuad'
    });
  }

  // Success celebration
  celebrateSuccess(container: HTMLElement): void {
    if (!this.animationsEnabled) {
      // Show static success state
      container.classList.add('success');
      return;
    }

    // Confetti effect
    this.createConfetti(container);

    // Success message bounce in
    anime({
      targets: container.querySelector('.success-message'),
      scale: [0, 1.2, 1],
      opacity: [0, 1],
      duration: 500,
      easing: 'easeOutElastic(1, .8)'
    });
  }

  // Page transition
  transitionTo(oldPage: HTMLElement, newPage: HTMLElement): Promise<void> {
    if (!this.animationsEnabled) {
      oldPage.style.display = 'none';
      newPage.style.display = 'block';
      return Promise.resolve();
    }

    return new Promise((resolve) => {
      const timeline = anime.timeline({
        complete: resolve
      });

      timeline
        .add({
          targets: oldPage,
          opacity: [1, 0],
          translateX: [0, -50],
          duration: 300,
          easing: 'easeInQuad'
        })
        .add({
          targets: newPage,
          opacity: [0, 1],
          translateX: [50, 0],
          duration: 300,
          easing: 'easeOutQuad'
        }, '-=100');
    });
  }

  // Typewriter effect for hints
  typewriterEffect(element: HTMLElement, text: string): void {
    if (!this.animationsEnabled) {
      element.textContent = text;
      return;
    }

    element.textContent = '';
    const chars = text.split('');

    anime({
      targets: chars,
      duration: 50 * chars.length,
      easing: 'linear',
      update: (anim) => {
        const index = Math.floor((anim.progress / 100) * chars.length);
        element.textContent = chars.slice(0, index).join('');
      }
    });
  }

  // Kill all animations (for cleanup)
  killAll(): void {
    this.activeAnimations.forEach(anim => anim.pause());
    this.activeAnimations = [];
  }
}
```

### 6.3 Animation Performance Budget

- **Max concurrent animations**: 20
- **Target frame rate**: 60 FPS
- **Animation budget per frame**: 16ms
- **Particle systems**: Max 100 particles simultaneously
- **GPU acceleration**: Use `transform` and `opacity` only (avoid layout thrashing)

### 6.4 Accessibility Considerations

**Reference**: `.plans/research/security_research.md:836-844`

```typescript
// Check user preference
const prefersReducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)');

if (prefersReducedMotion.matches) {
  // Disable animations
  // Apply final states immediately
  // Use simple fade transitions only
}

// Photosensitivity warning
if (flashCount > 3 && duration < 1000) {
  showWarning('This animation contains flashing. Disable animations in settings.');
}
```

---

## 7. Performance Architecture

### 7.1 Performance Targets

| Metric | Target | Measurement Tool |
|--------|--------|------------------|
| **First Contentful Paint (FCP)** | < 1.5s | Lighthouse |
| **Time to Interactive (TTI)** | < 3.5s | Lighthouse |
| **Bundle Size (initial)** | < 150KB gzipped | webpack-bundle-analyzer |
| **Animation Frame Rate** | 60 FPS | Chrome DevTools Performance |
| **Code Execution** | < 100ms per lesson | Performance API |

### 7.2 Optimization Strategies

#### Code Splitting
```typescript
// main.ts - Lazy load heavy components
const loadCodeEditor = () => import('./components/CodeEditor');
const loadAnimator = () => import('./core/Animator');

// Only load when needed
if (lessonRequiresEditor) {
  const { CodeEditor } = await loadCodeEditor();
  // ... initialize
}
```

#### Debouncing Preview Updates
```typescript
// src/components/PreviewPane.ts
import { debounce } from '../utils/debounce';

class PreviewPane {
  private updatePreview = debounce((code: string) => {
    const sanitized = this.sanitizer.sanitize(code);
    this.sandbox.executeCode(sanitized);
  }, 300); // Update max every 300ms
}
```

#### Virtual Scrolling for Lessons
```typescript
// Only render visible lessons in sidebar
// Use intersection observer for lazy loading
const observer = new IntersectionObserver((entries) => {
  entries.forEach(entry => {
    if (entry.isIntersecting) {
      loadLesson(entry.target.dataset.lessonId);
    }
  });
});
```

#### Image Optimization
- Use WebP with PNG fallback
- Lazy load images below the fold
- Use `loading="lazy"` attribute
- Serve responsive images with `srcset`

### 7.3 Caching Strategy

```typescript
// Service Worker (PWA)
const CACHE_VERSION = 'v1';
const STATIC_CACHE = 'static-v1';
const DYNAMIC_CACHE = 'dynamic-v1';

// Cache static assets
self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(STATIC_CACHE).then(cache => {
      return cache.addAll([
        '/',
        '/index.html',
        '/assets/main.js',
        '/assets/main.css'
      ]);
    })
  );
});

// Cache lessons on demand
self.addEventListener('fetch', (event) => {
  if (event.request.url.includes('/lessons/')) {
    event.respondWith(
      caches.match(event.request).then(response => {
        return response || fetch(event.request).then(fetchResponse => {
          return caches.open(DYNAMIC_CACHE).then(cache => {
            cache.put(event.request, fetchResponse.clone());
            return fetchResponse;
          });
        });
      })
    );
  }
});
```

---

## 8. State Management

### 8.1 Application State

```typescript
// src/core/AppState.ts
interface AppState {
  // Navigation
  currentLesson: string | null;
  currentCategory: string | null;

  // UI State
  sidebarOpen: boolean;
  editorFocused: boolean;
  previewVisible: boolean;

  // User State
  progress: UserProgress;

  // Transient State
  loading: boolean;
  error: Error | null;
}

class StateManager {
  private state: AppState;
  private listeners: Map<string, Function[]> = new Map();

  subscribe(key: keyof AppState, callback: Function): void {
    if (!this.listeners.has(key)) {
      this.listeners.set(key, []);
    }
    this.listeners.get(key)!.push(callback);
  }

  setState<K extends keyof AppState>(key: K, value: AppState[K]): void {
    this.state[key] = value;
    this.notify(key, value);

    // Persist progress to localStorage
    if (key === 'progress') {
      this.persistProgress(value as UserProgress);
    }
  }

  private notify(key: string, value: any): void {
    const callbacks = this.listeners.get(key) || [];
    callbacks.forEach(cb => cb(value));
  }

  private async persistProgress(progress: UserProgress): Promise<void> {
    const data = JSON.stringify(progress);
    const integrity = await this.integrityChecker.generateHMAC(data);

    localStorage.setItem('user_progress', data);
    localStorage.setItem('user_progress_hash', integrity);
  }
}
```

### 8.2 Event Bus for Component Communication

```typescript
// src/core/EventBus.ts
type EventCallback = (data: any) => void;

class EventBus {
  private events: Map<string, EventCallback[]> = new Map();

  on(event: string, callback: EventCallback): void {
    if (!this.events.has(event)) {
      this.events.set(event, []);
    }
    this.events.get(event)!.push(callback);
  }

  emit(event: string, data?: any): void {
    const callbacks = this.events.get(event) || [];
    callbacks.forEach(cb => cb(data));
  }

  off(event: string, callback: EventCallback): void {
    const callbacks = this.events.get(event) || [];
    const index = callbacks.indexOf(callback);
    if (index > -1) {
      callbacks.splice(index, 1);
    }
  }
}

// Usage
eventBus.on('lesson:complete', (lessonId) => {
  progressTracker.markComplete(lessonId);
  animator.celebrateSuccess(container);
  achievementManager.checkUnlocks();
});

eventBus.emit('lesson:complete', 'lesson-1');
```

---

## 9. Deployment Architecture

### 9.1 Hosting Strategy

**Recommended**: **Netlify** (free tier)

**Reference**: `.plans/research/architecture-research.md:391-406`

**Why Netlify**:
- Free SSL with Let's Encrypt
- Global CDN (edge caching)
- Deploy from Git (auto-deploy on push)
- Form handling (for feedback forms)
- Serverless functions (future feature expansion)
- Built-in asset optimization

**Alternatives**:
- **Vercel**: Similar features, great for Next.js but overkill here
- **GitHub Pages**: Free but limited features, no form handling
- **Cloudflare Pages**: Fast CDN, good alternative

### 9.2 Build Pipeline

```bash
# .github/workflows/deploy.yml
name: Deploy to Netlify

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: 18

      - name: Install dependencies
        run: npm ci

      - name: Run tests
        run: npm test

      - name: Run security audit
        run: npm audit --audit-level=moderate

      - name: Build
        run: npm run build
        env:
          NODE_ENV: production

      - name: Deploy to Netlify
        uses: netlify/actions/cli@master
        with:
          args: deploy --prod --dir=dist
        env:
          NETLIFY_AUTH_TOKEN: ${{ secrets.NETLIFY_AUTH_TOKEN }}
          NETLIFY_SITE_ID: ${{ secrets.NETLIFY_SITE_ID }}
```

### 9.3 Environment Configuration

```typescript
// vite.config.ts
import { defineConfig } from 'vite';

export default defineConfig({
  build: {
    target: 'es2020',
    outDir: 'dist',
    sourcemap: false, // Disable in production
    minify: 'terser',
    rollupOptions: {
      output: {
        manualChunks: {
          'vendor': ['animejs', 'dompurify'],
          'editor': ['codemirror']
        }
      }
    }
  },
  server: {
    port: 3000,
    headers: {
      // Security headers
      'Strict-Transport-Security': 'max-age=31536000; includeSubDomains',
      'X-Content-Type-Options': 'nosniff',
      'X-Frame-Options': 'DENY',
      'X-XSS-Protection': '1; mode=block'
    }
  }
});
```

### 9.4 DNS and SSL

- **Domain**: Register custom domain (optional)
- **SSL**: Automatic via Let's Encrypt
- **HSTS**: Enforce HTTPS everywhere
- **Preload**: Add to HSTS preload list (hstspreload.org)

---

## 10. Testing Strategy

### 10.1 Unit Tests (Vitest)

```typescript
// src/security/Sanitizer.test.ts
import { describe, it, expect } from 'vitest';
import { HTMLSanitizer } from './Sanitizer';

describe('HTMLSanitizer', () => {
  const sanitizer = new HTMLSanitizer();

  it('removes script tags', () => {
    const malicious = '<script>alert("XSS")</script><p>Safe</p>';
    const clean = sanitizer.sanitize(malicious);
    expect(clean).not.toContain('<script>');
    expect(clean).toContain('<p>Safe</p>');
  });

  it('removes event handlers', () => {
    const malicious = '<img src=x onerror="alert(1)">';
    const clean = sanitizer.sanitize(malicious);
    expect(clean).not.toContain('onerror');
  });

  it('allows safe HTML', () => {
    const safe = '<div class="test"><h1>Hello</h1></div>';
    const clean = sanitizer.sanitize(safe);
    expect(clean).toBe(safe);
  });
});
```

### 10.2 Integration Tests (Vitest)

```typescript
// src/components/CodeEditor.test.ts
describe('CodeEditor Integration', () => {
  it('validates and sanitizes user input before preview', async () => {
    const editor = new CodeEditor();
    const malicious = '<script>alert(1)</script><p>Test</p>';

    editor.setValue(malicious);
    await editor.updatePreview();

    const previewContent = editor.getPreviewContent();
    expect(previewContent).not.toContain('<script>');
    expect(previewContent).toContain('<p>Test</p>');
  });
});
```

### 10.3 E2E Tests (Playwright) - Optional

```typescript
// tests/e2e/lesson-flow.spec.ts
import { test, expect } from '@playwright/test';

test('user can complete first lesson', async ({ page }) => {
  await page.goto('/');

  // Start first lesson
  await page.click('text=Start Learning');

  // Read lesson content
  await expect(page.locator('h1')).toContainText('Your First HTML Tag');

  // Complete challenge
  await page.fill('.code-editor', '<h1>Hello World</h1>');

  // Check for success animation
  await expect(page.locator('.success-message')).toBeVisible();

  // Verify progress saved
  const progress = await page.evaluate(() => localStorage.getItem('user_progress'));
  expect(progress).toContain('lesson-1');
});
```

### 10.4 Security Testing

**Reference**: `.plans/research/security_research.md:599-684`

- **OWASP ZAP**: Automated vulnerability scanning
- **npm audit**: Weekly dependency checks
- **Snyk**: Continuous monitoring (GitHub integration)
- **Manual XSS testing**: Test payloads in sandbox environment

---

## 11. Monitoring and Analytics

### 11.1 Error Tracking

**Sentry** integration for client-side errors:

```typescript
// src/main.ts
import * as Sentry from '@sentry/browser';

Sentry.init({
  dsn: import.meta.env.VITE_SENTRY_DSN,
  environment: import.meta.env.MODE,
  beforeSend(event) {
    // Don't send user code content (privacy)
    if (event.request?.data?.code) {
      delete event.request.data.code;
    }
    return event;
  }
});
```

### 11.2 Performance Monitoring

```typescript
// Track key performance metrics
const lessonLoadTime = performance.measure('lesson-load', 'lesson-start', 'lesson-end');
const animationFPS = performance.getEntriesByType('measure').filter(e => e.name.includes('animation'));

// Send to analytics
analytics.track('lesson_performance', {
  lessonId: currentLesson,
  loadTime: lessonLoadTime.duration,
  fps: calculateAverageFPS(animationFPS)
});
```

### 11.3 User Analytics (Privacy-First)

**No personal data collection** - only anonymous usage metrics:

```typescript
interface AnalyticsEvent {
  event: string;
  properties: {
    lessonId?: string;
    category?: string;
    timestamp: string;
    // No user identifiers
  };
}

// Track lesson completion (anonymously)
analytics.track('lesson_complete', {
  lessonId: 'lesson-1',
  category: 'basics',
  timeSpent: '5min'
});
```

---

## 12. Internationalization (i18n) - Future Consideration

### 12.1 i18n Architecture (Post-MVP)

```typescript
// Future: Multi-language support
interface Translations {
  en: Record<string, string>;
  es: Record<string, string>;
  fr: Record<string, string>;
}

// src/core/i18n.ts
class I18n {
  private locale: string = 'en';
  private translations: Translations;

  t(key: string): string {
    return this.translations[this.locale][key] || key;
  }
}

// Usage
<h1>{i18n.t('lesson.title')}</h1>
```

### 12.2 Lesson Content Translation

- Store lessons in JSON per language: `lessons.en.json`, `lessons.es.json`
- Use crowdsourcing for translations (GitHub PRs)
- Implement language switcher in navbar

---

## 13. Trade-offs and Decisions

### 13.1 Optimizing For

✅ **Engagement**: Animation-first approach makes learning fun
✅ **Security**: Multi-layer XSS prevention, sandboxed execution
✅ **Performance**: < 150KB initial bundle, 60 FPS animations
✅ **Accessibility**: Reduced motion support, keyboard navigation
✅ **Simplicity**: No backend, pure frontend architecture

### 13.2 Sacrificing

❌ **User Accounts**: No cloud sync (localStorage only)
❌ **Social Features**: No comments, sharing, community (for now)
❌ **Advanced Features**: No AI assistance, no video tutorials
❌ **Real-time Collaboration**: No multiplayer or live help
❌ **Complex Backend**: No server-side code execution or API

### 13.3 Key Architectural Principles

1. **Security by Design**: Every feature considers XSS risks first
2. **Progressive Enhancement**: Works without JavaScript (basic HTML)
3. **Performance Budget**: Every feature must justify its bundle cost
4. **Accessibility First**: Animations respect user preferences
5. **Privacy Focused**: No tracking, no personal data collection

---

## 14. Implementation Phases

### Phase 1: Foundation (Week 1)
- [ ] Set up Vite + TypeScript + Tailwind boilerplate
- [ ] Implement security layer (Sanitizer, SandboxManager, Validator)
- [ ] Build core components (Navigation, LessonViewer)
- [ ] Create lesson data structure and first 3 lessons
- [ ] Implement localStorage progress tracking

### Phase 2: Editor & Preview (Week 2)
- [ ] Integrate CodeMirror 6
- [ ] Build sandboxed preview pane
- [ ] Implement code validation and sanitization flow
- [ ] Add syntax highlighting with Prism.js
- [ ] Create challenge validation system

### Phase 3: Animations (Week 3)
- [ ] Integrate Anime.js
- [ ] Build Animator class with core animations
- [ ] Implement success celebrations (confetti, bounces)
- [ ] Add micro-interactions (button clicks, hovers)
- [ ] Add page transitions between lessons
- [ ] Implement reduced motion support

### Phase 4: Gamification (Week 4)
- [ ] Build achievement system
- [ ] Implement XP and progress tracking
- [ ] Create achievement popup animations
- [ ] Add streak counter and milestone celebrations
- [ ] Design and implement badge/icon system

### Phase 5: Polish & Testing (Week 5)
- [ ] Write unit tests for security components
- [ ] Perform manual XSS testing
- [ ] Run OWASP ZAP security scan
- [ ] Performance optimization (code splitting, lazy loading)
- [ ] Accessibility audit (keyboard nav, screen readers)
- [ ] Cross-browser testing (Chrome, Firefox, Safari, Edge)

### Phase 6: Deployment (Week 6)
- [ ] Set up Netlify account and project
- [ ] Configure build pipeline and environment variables
- [ ] Set up custom domain and SSL
- [ ] Implement service worker for PWA
- [ ] Add error tracking with Sentry
- [ ] Deploy to production and monitor

---

## 15. Success Metrics

**Reference**: `.plans/research/product-research.md:98-105`

### 15.1 Technical Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| Page Load Time | < 2s | Lighthouse |
| Bundle Size | < 150KB | webpack-bundle-analyzer |
| Animation FPS | 60 FPS | Chrome DevTools |
| Security Score | A+ | securityheaders.com |
| Accessibility Score | > 95 | Lighthouse |

### 15.2 User Engagement Metrics

| Metric | Target | Tracking |
|--------|--------|----------|
| Time on site | > 5 min avg | Analytics |
| Lesson completion | > 60% | Progress data |
| Return visits | > 30% within 7 days | localStorage timestamps |
| Error rate | < 1% | Sentry |

---

## 16. Risks and Mitigation

**Reference**: `.plans/research/architecture-research.md:417-426`

| Risk | Impact | Mitigation |
|------|--------|------------|
| **Over-animation causes distraction** | HIGH | User preference toggle, reduced motion support |
| **XSS vulnerability in sandbox** | CRITICAL | Multi-layer defense (validation, sanitization, CSP, sandbox) |
| **Large bundle size** | MEDIUM | Code splitting, lazy loading, tree shaking |
| **Editor doesn't work on mobile** | HIGH | Use mobile-friendly CodeMirror 6, test early |
| **Browser compatibility issues** | MEDIUM | Target modern browsers (last 2 versions), graceful degradation |
| **Infinite loops crash browser** | HIGH | Execution timeout (5s), infinite loop detection |

---

## 17. Next Steps

### Immediate Actions

1. **Senior Developer**: Set up Vite + TypeScript + Tailwind boilerplate
2. **Senior Developer**: Implement security foundation (Sanitizer, SandboxManager)
3. **UI Designer**: Create design system and component library in Figma
4. **PM**: Create detailed lesson content outline (20-30 lessons)
5. **UX Designer**: Create user flow diagrams and wireframes

### Dependencies

- UI design must be complete before implementing components
- Lesson content structure needed before building LessonManager
- Security layer must be built before editor integration

### Delegation

**DELEGATE:**

- **senior_dev**: Implement security layer, set up project boilerplate, integrate CodeMirror
- **ui**: Create design system, component library, animation specifications
- **ux**: Create user flows, wireframes, interaction patterns
- **pm**: Write lesson content, define achievement criteria, prioritize features

---

## 18. References

**Research Documents**:
- `.plans/research/architecture-research.md` - Technical architecture patterns
- `.plans/research/security_research.md` - Comprehensive security analysis
- `.plans/research/ui-research.md` - Visual design and animation strategy
- `.plans/research/ux-research.md` - User experience patterns and flows
- `.plans/research/product-research.md` - Product strategy and competitive analysis

**External Resources**:
- [Vite Documentation](https://vitejs.dev/)
- [Anime.js Documentation](https://animejs.com/)
- [CodeMirror 6 Documentation](https://codemirror.net/)
- [DOMPurify Documentation](https://github.com/cure53/DOMPurify)
- [OWASP XSS Prevention](https://cheatsheetseries.owasp.org/cheatsheets/XSS_Prevention_Cheat_Sheet.html)

---

**Status**: PHASE_COMPLETE: planning
**Next Phase**: Development (Senior Dev implements architecture)
**Estimated Timeline**: 6 weeks to MVP

---

## Appendix A: Color Palette (Electric Playground)

**Reference**: `.plans/research/ui-research.md:76-84` and `.plans/research/ui-research.md:202-209`

```css
:root {
  /* Primary Colors */
  --color-primary: #FF2E63;        /* Hot Pink/Magenta */
  --color-secondary: #00D9FF;      /* Electric Cyan */

  /* Background Colors */
  --color-bg-dark: #1A1A2E;        /* Charcoal */
  --color-bg-light: #EEEEF7;       /* Light Gray */
  --color-surface-dark: #16213E;   /* Dark Surface */
  --color-surface-light: #FFFFFF;  /* White */

  /* Semantic Colors */
  --color-success: #00FF88;        /* Bright Green */
  --color-error: #FF3366;          /* Bright Red */
  --color-warning: #FFB627;        /* Amber */
  --color-info: #00D9FF;           /* Cyan */

  /* Typography */
  --font-display: 'Space Grotesk', sans-serif;
  --font-body: 'Inter', sans-serif;
  --font-code: 'JetBrains Mono', monospace;

  /* Spacing */
  --spacing-xs: 4px;
  --spacing-sm: 8px;
  --spacing-md: 16px;
  --spacing-lg: 24px;
  --spacing-xl: 32px;

  /* Radius */
  --radius-sm: 4px;
  --radius-md: 8px;
  --radius-lg: 12px;
  --radius-xl: 16px;

  /* Shadows */
  --shadow-glow-primary: 0 0 20px rgba(255, 46, 99, 0.3);
  --shadow-glow-secondary: 0 0 20px rgba(0, 217, 255, 0.3);
}
```

---

## Appendix B: File Structure Template

```
html-helper/
├── .github/
│   └── workflows/
│       └── deploy.yml
├── public/
│   ├── favicon.ico
│   ├── manifest.json
│   └── robots.txt
├── src/
│   ├── assets/
│   │   ├── images/
│   │   ├── icons/
│   │   └── fonts/
│   ├── components/
│   │   ├── Navigation.ts
│   │   ├── LessonViewer.ts
│   │   ├── CodeEditor.ts
│   │   ├── PreviewPane.ts
│   │   ├── ChallengeCard.ts
│   │   ├── AchievementPopup.ts
│   │   └── AnimatedButton.ts
│   ├── core/
│   │   ├── LessonManager.ts
│   │   ├── ProgressTracker.ts
│   │   ├── Animator.ts
│   │   ├── Validator.ts
│   │   ├── EventBus.ts
│   │   └── AppState.ts
│   ├── security/
│   │   ├── SandboxManager.ts
│   │   ├── Sanitizer.ts
│   │   ├── InputValidator.ts
│   │   ├── IntegrityChecker.ts
│   │   └── CSPPolicy.ts
│   ├── data/
│   │   ├── lessons.json
│   │   ├── achievements.json
│   │   └── challenges.json
│   ├── utils/
│   │   ├── storage.ts
│   │   ├── dom.ts
│   │   └── debounce.ts
│   ├── styles/
│   │   ├── main.css
│   │   ├── animations.css
│   │   ├── components.css
│   │   └── themes.css
│   ├── types/
│   │   ├── Lesson.ts
│   │   ├── Achievement.ts
│   │   └── UserProgress.ts
│   ├── App.ts
│   └── main.ts
├── tests/
│   ├── unit/
│   │   └── security/
│   └── e2e/
│       └── lesson-flow.spec.ts
├── .env.example
├── .eslintrc.json
├── .gitignore
├── .prettierrc
├── index.html
├── package.json
├── tailwind.config.js
├── tsconfig.json
├── vite.config.ts
└── README.md
```

---

**END OF ARCHITECTURE SPECIFICATION**

PHASE_COMPLETE: planning
