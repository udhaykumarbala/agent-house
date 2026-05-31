# Security Research — Unit Converter

**Role**: Security Expert
**Date**: 2026-03-17
**Project**: Premium Minimal Unit Converter (Length, Weight, Temperature)
**Stack**: Vite + TypeScript + Tailwind CSS (static client-side app, no backend)

---

## 1. Threat Assessment

### Attack Surface Summary

This is a **low-risk application** by nature — a purely client-side math utility with no backend, no user accounts, no database, and no sensitive data. However, "low risk" is not "no risk." The following threat actors and motivations apply:

| Threat Actor | Motivation | Likelihood |
|---|---|---|
| Opportunistic attacker | Deface or inject malicious scripts via supply chain | Low |
| Malicious CDN/dependency | Compromise via poisoned npm package | Low |
| Man-in-the-middle | Inject scripts if served over HTTP | Low (if HTTPS enforced) |
| Self-XSS / social engineering | Trick user into pasting malicious input | Very Low |

### Key Insight
The primary security concerns for this project are **supply chain integrity**, **Content Security Policy**, and **secure deployment** — not traditional web app vulnerabilities like SQL injection or authentication bypass, which don't apply here.

---

## 2. Vulnerability Analysis

### 2.1 Input Handling
- **Risk**: Users type numeric values for conversion. If these values are rendered back into the DOM unsafely (e.g., via `innerHTML`), it could create a DOM-based XSS vector.
- **Likelihood**: Low — TypeScript and modern frameworks typically use safe DOM APIs.
- **Severity**: Medium (if exploited, could execute arbitrary JS in user's browser context).

### 2.2 Supply Chain (npm Dependencies)
- **Risk**: Vite, Tailwind, and their transitive dependencies introduce third-party code. A compromised dependency could inject malicious code into the build output.
- **Likelihood**: Low but non-trivial — supply chain attacks (event-stream, ua-parser-js, colors.js) are a growing trend.
- **Severity**: High (full control over build output).

### 2.3 Deployment & Transport Security
- **Risk**: If deployed over HTTP (no TLS), a MITM attacker can inject scripts into the page.
- **Likelihood**: Low if hosted on modern platforms (Vercel, Netlify, GitHub Pages all enforce HTTPS).
- **Severity**: High (full page compromise).

### 2.4 Client-Side Storage
- **Risk**: If the app stores user preferences (e.g., last-used units, theme) in `localStorage`, this data is accessible to any script running on the same origin.
- **Likelihood**: Low — no sensitive data in a converter app.
- **Severity**: Very Low (no PII or secrets stored).

### 2.5 Third-Party Fonts / Assets
- **Risk**: Loading Google Fonts or other external assets creates external dependencies and potential tracking/privacy concerns. If the font CDN is compromised, it could serve malicious CSS or use CSS-based data exfiltration.
- **Likelihood**: Very Low.
- **Severity**: Low.

---

## 3. OWASP Top 10 Relevance Check

| OWASP Category | Applicable? | Notes |
|---|---|---|
| A01: Broken Access Control | **No** | No auth, no backend, no protected resources |
| A02: Cryptographic Failures | **No** | No data encryption needed; no secrets stored |
| A03: Injection | **Minimal** | DOM-based XSS if input is rendered via `innerHTML` |
| A04: Insecure Design | **Low** | Simple math utility — design is inherently safe |
| A05: Security Misconfiguration | **Yes** | Missing CSP headers, permissive CORS on hosting |
| A06: Vulnerable Components | **Yes** | npm dependency supply chain risk |
| A07: Auth Failures | **No** | No authentication |
| A08: Data Integrity Failures | **Yes** | Build pipeline integrity, dependency verification |
| A09: Logging & Monitoring | **No** | No backend to monitor |
| A10: SSRF | **No** | No server-side requests |

**Applicable categories**: A03 (minimal), A05, A06, A08.

---

## 4. Recommendations

### 4.1 Input Sanitization — Priority: Medium

Even though this is a calculator, all user input should be treated as untrusted.

- **Use `textContent` not `innerHTML`** when rendering conversion results to the DOM. TypeScript with safe DOM APIs prevents XSS by default.
- **Validate input is numeric** before processing. Reject or strip non-numeric characters (allow digits, decimal point, negative sign only).
- **Use `parseFloat()` / `Number()`** for conversion, which naturally reject script content.

```typescript
// GOOD: Safe rendering
resultElement.textContent = convertedValue.toFixed(2);

// BAD: Potential XSS
resultElement.innerHTML = `${userInput} converts to ${result}`;
```

### 4.2 Content Security Policy (CSP) — Priority: High

Add a strict CSP header (or meta tag) to prevent inline script injection. This is the single highest-value security measure for a static site.

```html
<meta http-equiv="Content-Security-Policy"
  content="default-src 'self';
           script-src 'self';
           style-src 'self' 'unsafe-inline' https://fonts.googleapis.com;
           font-src 'self' https://fonts.gstatic.com;
           img-src 'self' data:;
           connect-src 'none';
           object-src 'none';
           base-uri 'self';">
```

**Notes**:
- `'unsafe-inline'` for styles is needed if Tailwind injects inline styles. Prefer removing it if possible by using external stylesheets only.
- `connect-src 'none'` blocks all network requests (fetch, XHR, WebSocket) — appropriate since the app is offline-capable.
- `object-src 'none'` blocks Flash/Java plugins.
- `base-uri 'self'` prevents base tag injection attacks.

### 4.3 Dependency Security — Priority: High

- **Pin exact dependency versions** in `package.json` (no `^` or `~` ranges) or use a lockfile (`package-lock.json` / `pnpm-lock.yaml`).
- **Keep dependencies minimal** — the template already calls for no additional libraries, which is ideal.
- **Run `npm audit`** before each deployment.
- **Consider Subresource Integrity (SRI)** if loading any scripts from CDNs (unlikely with Vite bundling, but good practice).

### 4.4 Secure Deployment — Priority: Medium

- **Enforce HTTPS** — all modern hosting platforms do this by default (Vercel, Netlify, GitHub Pages).
- **Add security headers** via hosting config or `_headers` file:
  ```
  X-Content-Type-Options: nosniff
  X-Frame-Options: DENY
  Referrer-Policy: strict-origin-when-cross-origin
  Permissions-Policy: camera=(), microphone=(), geolocation=()
  ```
- **X-Frame-Options: DENY** prevents clickjacking (embedding the converter in a malicious iframe).
- **Permissions-Policy** disables APIs the app doesn't need.

### 4.5 Self-Hosting Fonts — Priority: Low

- Consider self-hosting Google Fonts rather than loading from Google's CDN to eliminate the external dependency and improve privacy.
- Removes a third-party request and a tracking vector.
- Improves performance (no DNS lookup to fonts.googleapis.com).

### 4.6 Build Integrity — Priority: Low

- Verify the Vite build output doesn't include source maps in production (`build.sourcemap: false` in vite config).
- Source maps expose original TypeScript source code, which while not a direct vulnerability, aids attackers in understanding app logic.

---

## 5. Implementation Notes for Developers

### Safe Input Handling Pattern
```typescript
function sanitizeNumericInput(value: string): number | null {
  const trimmed = value.trim();
  if (trimmed === '' || trimmed === '-') return null;
  const num = Number(trimmed);
  return Number.isFinite(num) ? num : null;
}
```

### CSP Meta Tag
Place in `<head>` of `index.html` before any script tags. This is the simplest approach for a static site.

### Security Headers
For Netlify, create a `_headers` file:
```
/*
  X-Content-Type-Options: nosniff
  X-Frame-Options: DENY
  Referrer-Policy: strict-origin-when-cross-origin
  Permissions-Policy: camera=(), microphone=(), geolocation=()
  Content-Security-Policy: default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com; img-src 'self' data:; connect-src 'none'; object-src 'none'; base-uri 'self'
```

For Vercel, add to `vercel.json` headers config. For GitHub Pages, use the CSP meta tag approach (no server-side header control).

---

## 6. Testing Suggestions

| Test | Type | Priority |
|---|---|---|
| Paste `<script>alert(1)</script>` into input field — verify no execution | Manual | High |
| Paste `"><img src=x onerror=alert(1)>` into input — verify safe rendering | Manual | High |
| Verify CSP header/meta tag is present in production build | Automated | High |
| Run `npm audit` — verify no known vulnerabilities | CI/CD | High |
| Verify no source maps in production build | Automated | Medium |
| Test with extremely large numbers (1e308, Infinity, NaN) — verify graceful handling | Manual | Medium |
| Verify HTTPS redirect on deployment platform | Manual | Medium |
| Test with Unicode/emoji input — verify no crashes | Manual | Low |

---

## 7. Risk Summary

| Issue | Severity | Likelihood | Priority | Recommendation |
|---|---|---|---|---|
| Missing CSP headers | Medium | High (default has none) | **High** | Add CSP meta tag to index.html |
| npm supply chain compromise | High | Low | **High** | Pin deps, use lockfile, run audit |
| DOM-based XSS via innerHTML | Medium | Low | **Medium** | Use textContent, validate numeric input |
| Missing security headers | Low | Medium | **Medium** | Add X-Frame-Options, X-Content-Type-Options, etc. |
| Source maps in production | Low | Medium | **Low** | Disable in Vite config |
| Third-party font dependency | Low | Very Low | **Low** | Consider self-hosting fonts |

### Overall Risk Rating: **LOW**

This is an inherently low-risk application. The recommended mitigations are best practices that add defense-in-depth without adding complexity. The two high-priority items (CSP and dependency management) are simple to implement and provide significant protection.

---

## 8. Delegation

```
DELEGATE:
- architect: Ensure CSP meta tag and security headers are included in the base HTML template and deployment config
- senior_dev: Use textContent (not innerHTML) for all DOM rendering; validate all input as numeric; disable source maps in production Vite config
```

---

**Status**: PHASE_COMPLETE: research
