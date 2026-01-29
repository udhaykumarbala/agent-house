# APPROVED PLAN: HTML Helper - Interactive HTML Learning Platform

**Reviewed by**: CEO
**Date**: 2026-01-29
**Status**: APPROVED FOR DEVELOPMENT

---

## Executive Summary

HTML Helper is an **animation-driven, interactive learning platform** that teaches HTML through play, not study. Every user interaction triggers delightful animations and instant visual feedback, making coding feel like a game.

**Target Audience**: "Curious Casey" - teens and young adults (14-35) who want to learn HTML but find traditional tutorials boring.

**Unique Value Proposition**: Learn HTML by playing with it - every click, type, and hover brings code to life with delightful animations.

**Differentiation**: The most fun you'll have learning HTML. Unlike Codecademy, freeCodeCamp, and W3Schools which add gamification as an afterthought, we're building the entire experience around delightful animations and instant visual feedback.

---

## Approved Specifications

All specification documents have been reviewed and approved:

### 1. Product Specification
**File**: `.plans/specs/product-spec.md`
**Status**: ✅ Approved
**Key Elements**:
- 15 progressive lessons (from `<p>` tag to complete webpage)
- Animation-driven feedback system
- Progress tracking with localStorage
- Achievement system for milestones
- 100% free, no signup required
- Success metrics: >60% lesson completion, >30% return rate

### 2. UX Specification
**File**: `.plans/specs/ux-spec.md`
**Status**: ✅ Approved
**Key Elements**:
- Split-screen learning interface (editor + live preview)
- Interactive code editor with real-time updates
- Character-driven animations (mascot reactions)
- Progressive lesson map with visual progress
- Comprehensive accessibility (keyboard nav, screen readers, reduced motion)
- Error states with gentle guidance (not harsh errors)

### 3. UI Specification
**File**: `.plans/specs/ui-spec.md`
**Status**: ✅ Approved
**Key Elements**:
- "Electric Playground" theme: Magenta (#FF2E63) + Cyan (#00D9FF)
- Typography: Space Grotesk (display), Inter (body), JetBrains Mono (code)
- Complete animation library (confetti, bounces, pulses, transitions)
- Dark mode primary, light mode optional
- WCAG 2.1 AA compliant (contrast ratios verified)
- Smooth 60 FPS animations with performance budgets

### 4. Architecture Specification
**File**: `.plans/specs/architecture-spec.md`
**Status**: ✅ Approved
**Key Elements**:
- Pure frontend architecture (Vite + TypeScript + Tailwind)
- Security-first design: DOMPurify + sandboxed iframes + CSP
- Animation engine: Anime.js (9KB, lightweight)
- Code editor: CodeMirror 6 (modern, mobile-friendly)
- Progressive Web App with offline support
- Bundle size target: <150KB gzipped

### 5. Security Specification
**File**: `.plans/specs/security-spec.md`
**Status**: ✅ Approved
**Key Elements**:
- 7-layer defense-in-depth (HTTPS, CSP, validation, sanitization, sandboxing, timeouts, monitoring)
- XSS prevention: DOMPurify sanitization + sandboxed iframe execution
- Execution timeout: 5 seconds (prevents infinite loops)
- Rate limiting: 50 executions per minute
- GDPR compliant (localStorage only, no server tracking)

---

## Key Decisions

### 1. Technology Stack (FINAL)
- **Build Tool**: Vite 5.x (fast HMR, optimized builds)
- **Language**: TypeScript 5.x (type safety for complex interactions)
- **Styling**: Tailwind CSS 3.x (rapid development, consistent design system)
- **Animations**: Anime.js 3.x (9KB, perfect for DOM animations)
- **Code Editor**: CodeMirror 6 (modern, performant, mobile-friendly)
- **HTML Sanitization**: DOMPurify 3.x (industry standard XSS prevention)
- **Icons**: Phosphor Icons (5KB, tree-shakeable)
- **Storage**: localForage (8KB, enhanced localStorage with IndexedDB fallback)

### 2. Template Choice
**Chosen**: `static-enhanced` (Vite + Tailwind + TypeScript)
**Rationale**: Modern build tools for animations, no backend needed, right complexity level

### 3. Security Approach
**Multi-layer defense**:
1. Input validation (length, patterns)
2. HTML sanitization (DOMPurify)
3. Sandboxed iframe execution (no parent access)
4. Content Security Policy (strict CSP headers)
5. HTTPS enforcement (HSTS)
6. Rate limiting (prevent abuse)
7. Monitoring (CSP violations, errors)

### 4. MVP Scope
**15 Lessons** (NOT "10-15", exactly 15 for clear milestone):
1. Your First Tag: `<p>`
2. Headings: `<h1>` through `<h6>`
3. Bold and Italic: `<strong>` and `<em>`
4. Lists: `<ul>`, `<ol>`, `<li>`
5. Links: `<a href="">`
6. Images: `<img src="" alt="">`
7. Line Breaks: `<br>` and `<hr>`
8. Divs: `<div>` for grouping
9. Spans: `<span>` for inline styling
10. Basic Structure: `<html>`, `<head>`, `<body>`
11. Meta Tags: `<title>`, `<meta>`
12. Forms Part 1: `<form>`, `<input>`, `<button>`
13. Forms Part 2: `<textarea>`, `<select>`, `<label>`
14. Tables: `<table>`, `<tr>`, `<td>`, `<th>`
15. Final Challenge: Build a complete webpage

**Post-MVP Features** (defer to Phase 2):
- CSS lessons
- JavaScript basics
- User accounts and cloud sync
- Social sharing
- Community features
- Advanced achievements

### 5. Hosting & Deployment
**Provider**: Netlify (free tier)
**Rationale**: Free SSL, global CDN, deploy from Git, serverless functions for future expansion
**Domain**: TBD (consider: htmlhelper.dev, learnhtml.fun, htmlplay.com)
**CI/CD**: GitHub Actions for automated builds and security scans

---

## Implementation Order

### Phase 1: Foundation (Week 1) - Senior Dev
**Goal**: Security foundation + basic structure

- [ ] Set up Vite + TypeScript + Tailwind boilerplate
- [ ] Implement security layer:
  - [ ] DOMPurify HTML sanitization wrapper
  - [ ] Sandboxed iframe manager
  - [ ] Input validator (length, patterns)
  - [ ] CSP configuration
- [ ] Build core components:
  - [ ] Navigation with progress bar
  - [ ] Lesson viewer (markdown rendering)
- [ ] Create lesson data structure (JSON schema)
- [ ] Implement localStorage progress tracking with HMAC integrity

**Deliverable**: Secure foundation with 3 demo lessons

---

### Phase 2: Editor & Preview (Week 2) - Senior Dev
**Goal**: Interactive code editor with live preview

- [ ] Integrate CodeMirror 6
  - [ ] HTML syntax highlighting
  - [ ] Auto-close tags
  - [ ] Debounced updates (300ms)
- [ ] Build sandboxed preview pane
  - [ ] Iframe with `sandbox="allow-scripts"` only
  - [ ] Execution timeout (5 seconds)
  - [ ] Rate limiter (50 executions/min)
- [ ] Implement code validation
  - [ ] Challenge validation (contains, exact, regex)
  - [ ] Progressive hints system
- [ ] Add tag palette (visual buttons for mobile)

**Deliverable**: Working code editor with live preview, 5 lessons complete

---

### Phase 3: Animations (Week 3) - Senior Dev
**Goal**: Delightful animations on every interaction

- [ ] Integrate Anime.js
- [ ] Build Animator class with core animations:
  - [ ] Button clicks (scale + bounce)
  - [ ] Success celebrations (confetti + glow)
  - [ ] Error feedback (gentle shake)
  - [ ] Page transitions (slide + fade)
  - [ ] Typing feedback (cursor blink, syntax highlighting)
- [ ] Implement confetti system (canvas-based particles)
- [ ] Add character reactions (mascot: idle, success, thinking)
- [ ] Respect `prefers-reduced-motion` (accessibility)

**Deliverable**: All P0 animations working smoothly (60 FPS), 10 lessons complete

---

### Phase 4: Gamification (Week 4) - Senior Dev
**Goal**: Progress tracking and achievements

- [ ] Build achievement system
  - [ ] Define 10-15 achievements (first tag, speed demon, perfectionist, etc.)
  - [ ] Unlock conditions (lesson-complete, streak, speed)
  - [ ] Badge popup animations
- [ ] Implement XP system
  - [ ] +25 XP per lesson
  - [ ] Count-up animation on XP gain
- [ ] Create lesson map
  - [ ] Visual progression path
  - [ ] Locked/unlocked/current states
  - [ ] Section grouping (basics, text, interactive)
- [ ] Add streak counter (consecutive days)
- [ ] Implement celebration milestones (5 lessons, 10 lessons, all complete)

**Deliverable**: Full gamification system, all 15 lessons complete

---

### Phase 5: Polish & Testing (Week 5) - Senior Dev
**Goal**: Production-ready quality

- [ ] Write security unit tests
  - [ ] HTML sanitization (XSS payloads)
  - [ ] Execution timeout (infinite loops)
  - [ ] Rate limiting (abuse prevention)
  - [ ] Input validation
- [ ] Manual XSS testing (penetration test)
- [ ] Run OWASP ZAP security scan
- [ ] Performance optimization
  - [ ] Code splitting (editor chunk, lesson chunks)
  - [ ] Lazy load Anime.js and CodeMirror
  - [ ] Image optimization (WebP with PNG fallback)
  - [ ] Tailwind CSS purging
- [ ] Accessibility audit
  - [ ] Keyboard navigation (all interactive elements)
  - [ ] Screen reader support (ARIA labels)
  - [ ] Focus indicators (visible on Tab)
  - [ ] Touch targets (44x44px minimum on mobile)
- [ ] Cross-browser testing (Chrome, Firefox, Safari, Edge)
- [ ] Mobile responsive testing (real devices)

**Deliverable**: <150KB bundle, >95 Lighthouse score, passing security audit

---

### Phase 6: Deployment (Week 6) - Senior Dev + PM
**Goal**: Launch to production

- [ ] Set up Netlify account and project
- [ ] Configure security headers in `netlify.toml`:
  - [ ] HSTS (Strict-Transport-Security)
  - [ ] CSP (Content-Security-Policy)
  - [ ] X-Frame-Options: DENY
  - [ ] X-Content-Type-Options: nosniff
- [ ] Set up custom domain and SSL (auto via Let's Encrypt)
- [ ] Implement service worker (PWA offline support)
- [ ] Add error tracking with Sentry (or alternative)
- [ ] Configure analytics (Plausible or Fathom - privacy-focused)
- [ ] Create feedback mechanism (simple form)
- [ ] Deploy to production
- [ ] Monitor for first 48 hours (error rates, CSP violations)

**Deliverable**: Live production site, monitoring dashboard

---

## Outstanding Clarifications

These minor items need clarification before/during development:

### Content Creation
**Owner**: PM
**Task**: Write all 15 lesson content files
**Format**: JSON with fields: title, introduction, theory, examples, challenge, hints
**Timeline**: Week 1-2 (parallel with development)

### Mascot Character Design
**Owner**: UI Designer
**Task**: Design mascot character with 3 states
**States**: Idle (watching), Success (celebrating), Thinking (confused)
**Format**: SVG (lightweight, scalable)
**Timeline**: Week 2 (needed for animations in Week 3)

### Error Message Copy
**Owner**: UX Designer
**Task**: Write friendly error messages for all error states
**Examples**: Timeout, rate limit, invalid code, network error
**Tone**: Encouraging, never harsh (see UX spec tone guidelines)
**Timeline**: Week 2 (needed for error handling)

### Hosting Provider Confirmation
**Owner**: Architect
**Decision**: Confirm Netlify (recommended) or Vercel
**Considerations**: Both have free tiers, similar features
**Recommendation**: Netlify (slightly better for static sites)
**Timeline**: Week 1

### Error Monitoring Tool
**Owner**: Architect
**Decision**: Sentry (industry standard) or alternative (Rollbar, Bugsnag)
**Budget**: Free tier should suffice for MVP
**Timeline**: Week 6 (deployment phase)

---

## Notes for Development Team

### Security is Non-Negotiable
- **ALL user HTML MUST be sanitized with DOMPurify**
- **ALL code MUST execute in sandboxed iframe** (no `allow-same-origin`)
- **Test XSS payloads regularly** (see Security Spec section 5.2)
- **Never use `innerHTML` with unsanitized input**
- **Run `npm audit` before every deployment**

### Performance Matters
- **Target**: <150KB initial bundle (gzipped)
- **Measure**: Use webpack-bundle-analyzer
- **Optimize**: Code split by lesson, lazy load heavy libraries
- **Maintain**: 60 FPS animations (Chrome DevTools Performance tab)

### Accessibility is Critical
- **Keyboard navigation**: Every interactive element must be Tab-accessible
- **Reduced motion**: Respect `prefers-reduced-motion` media query
- **Screen readers**: Add ARIA labels to all custom components
- **Color contrast**: Verify all text meets WCAG AA (4.5:1 minimum)
- **Touch targets**: Minimum 44x44px on mobile

### Testing Strategy
- **Unit tests**: Security components (sanitizer, validator, sandbox manager)
- **Integration tests**: Code editor + preview flow
- **Manual testing**: XSS payloads, cross-browser, mobile devices
- **Automated scanning**: OWASP ZAP, npm audit, Snyk

### Quality Bar
- **Lighthouse score**: >95 (performance, accessibility, best practices, SEO)
- **Security headers**: A+ on securityheaders.com
- **Browser support**: Last 2 versions of Chrome, Firefox, Safari, Edge
- **Mobile responsive**: Works on iPhone SE (smallest common screen)
- **Animation smoothness**: 60 FPS or use reduced motion fallback

---

## Success Criteria

### Launch Criteria (Must-Have)
- [ ] All 15 lessons complete with validation
- [ ] Progress tracking functional across sessions
- [ ] All P0 animations implemented and smooth (60fps)
- [ ] Mobile responsive design tested on real devices
- [ ] Security audit passed (all critical/high vulnerabilities fixed)
- [ ] Accessibility audit passed (WCAG AA compliance)
- [ ] Performance budget met (<150KB initial load, <2s FCP)
- [ ] Browser compatibility tested (Chrome, Firefox, Safari)
- [ ] Error monitoring setup (Sentry or alternative)
- [ ] Feedback mechanism in place (simple form)

### Post-Launch Success Metrics (Week 1-4)

**Engagement Metrics** (Primary):
- Time on site: >5 minutes average session
- Lesson completion: >60% complete first 3 lessons
- Return visit rate: >30% return within 7 days
- Lesson progression: Average 5+ lessons per session

**Delight Metrics** (Secondary):
- Animation disable rate: <10% (validates animations enhance, not annoy)
- Sandbox usage: >40% of users who complete 3 lessons
- Share rate: >5% share achievements or refer friends

**Technical Metrics**:
- Load time: <2 seconds on 3G
- Error rate: <5% of sessions encounter JavaScript errors
- Security incidents: 0 (no XSS exploits, no CSP violations)

---

## Risk Management

### Identified Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| **Over-animation causes user fatigue** | HIGH | User preference toggle, `prefers-reduced-motion`, animation testing with real users |
| **XSS vulnerability in sandbox** | CRITICAL | Multi-layer defense (validation → sanitization → sandboxing → CSP → monitoring) |
| **Large bundle size slows load** | MEDIUM | Code splitting, lazy loading, Tailwind purging, bundle analysis |
| **Editor difficult on mobile** | HIGH | CodeMirror 6 is mobile-optimized, tag palette for touch input, early mobile testing |
| **6-week timeline too aggressive** | MEDIUM | Clear prioritization (P0/P1/P2), defer post-MVP features, daily progress tracking |
| **Infinite loops crash browser** | HIGH | 5-second execution timeout, iframe destruction, rate limiting |

---

## Future Vision (Post-MVP)

### Phase 2: Enhanced Learning (Month 2-3)
- CSS learning module (same playful approach)
- JavaScript basics (DOM manipulation)
- Sandbox/playground mode (free experimentation)
- Code validation with smart hints

### Phase 3: Social Features (Month 4-6)
- Share code snippets with unique URLs
- Gallery of student projects
- Peer feedback and encouragement
- Community leaderboard

### Phase 4: Advanced Topics (Month 6+)
- Accessibility best practices
- Semantic HTML
- SEO optimization
- Performance optimization

---

## Approval Signatures

- [x] **CEO**: Vision and scope aligned with market opportunity ✓
- [x] **Product Manager**: Features prioritized, success metrics defined ✓
- [x] **UX Designer**: User flows validated, accessibility requirements clear ✓
- [x] **UI Designer**: Visual design system complete, animation specs detailed ✓
- [x] **Architect**: Technical architecture sound, implementation plan realistic ✓
- [x] **Security Expert**: Security controls comprehensive, XSS risks mitigated ✓

---

## Development Authorization

**Status**: ✅ **APPROVED FOR DEVELOPMENT**

**Start Date**: 2026-01-29
**Target Launch**: 2026-03-11 (6 weeks)
**Sprint Structure**: Weekly sprints aligned with phases 1-6

**Lead Developer**: Senior Dev
**Supporting Roles**: PM (content), UI Designer (mascot), UX Designer (copy)

---

## Communication & Reporting

### Daily Standups (Async)
- What did you complete yesterday?
- What will you complete today?
- Any blockers?

### Weekly Reviews (End of Each Phase)
- Demo working features
- Review against phase deliverables
- Adjust timeline if needed
- Security check-in (any new vulnerabilities?)

### Launch Checklist Review (Week 5)
- Run through entire launch criteria checklist
- Identify any gaps or issues
- Plan remediation before Week 6 deployment

---

## Emergency Contacts

| Role | Responsibility | Contact |
|------|----------------|---------|
| CEO | Vision alignment, go/no-go decisions | [TBD] |
| Senior Developer | Implementation, emergency fixes | [TBD] |
| Product Manager | Scope decisions, user communication | [TBD] |
| Security Expert | Security review, incident response | [TBD] |
| Architect | Technical decisions, architecture changes | [TBD] |

---

## Appendix: Quick Reference

### Core Technologies
- Vite 5.x, TypeScript 5.x, Tailwind CSS 3.x
- Anime.js 3.x, CodeMirror 6, DOMPurify 3.x

### Color Palette
- Primary: #FF2E63 (Magenta)
- Secondary: #00D9FF (Cyan)
- Success: #00FF9D (Green)
- Background: #1A1A2E (Dark)

### Security Essentials
- DOMPurify sanitization
- Sandboxed iframe (allow-scripts only)
- CSP headers (default-src 'self')
- 5-second execution timeout
- 50 executions per minute rate limit

### Performance Targets
- <150KB initial bundle (gzipped)
- <2s First Contentful Paint
- 60 FPS animations
- >95 Lighthouse score

---

**END OF APPROVED PLAN**

**Status**: PLAN APPROVED
**Next Step**: BEGIN DEVELOPMENT (Phase 1: Foundation)

---

**Reviewed and approved by CEO on 2026-01-29**

This plan represents the consolidated vision of the entire team's research and specifications. All agents have aligned on scope, timeline, and quality standards. Development may now proceed with confidence.

---

**COMPLETE: Plan approved. Development may begin.**
