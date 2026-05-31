# Security Specification — UnitShift Converter

**Role**: Security Expert
**Date**: 2026-03-17
**Phase**: Planning
**Stack**: Vite + TypeScript + Tailwind CSS (static client-side SPA, no backend)

---

## 1. Threat Assessment

### Overall Risk: LOW

This is a purely client-side math utility with no backend, no authentication, no database, and no sensitive data. The attack surface is minimal. However, defense-in-depth principles still apply — the goal is to ship a hardened static site that follows best practices even when the risk profile is low.

### Threat Model

| Threat Actor | Attack Vector | Impact | Likelihood | Priority |
|---|---|---|---|---|
| Supply chain attacker | Compromised npm dependency injects malicious code into build output | High — full control of build | Low | **High** |
| MITM attacker | Script injection if served over HTTP | High — full page compromise | Very Low (HTTPS enforced by modern hosts) | **Medium** |
| Self-XSS / social engineering | Tricking user into pasting malicious input into converter fields | Medium — JS execution in user context | Very Low | **Low** |
| Malicious iframe embedder | Clickjacking — embedding app in a deceptive frame | Low — no actions worth hijacking | Low | **Low** |
| CDN compromise | Malicious Google Fonts payload or CSS exfiltration | Low — CSS-only attack surface | Very Low | **Low** |

### What Is NOT a Threat

These common web vulnerabilities **do not apply** to this project:

- **SQL Injection** — No database
- **Authentication/Session attacks** — No auth system
- **CSRF** — No state-changing server requests
- **SSRF** — No server-side code
- **Privilege escalation** — No roles or permissions
- **Data breach** — No PII or sensitive data stored

---

## 2. OWASP Top 10 Applicability

| OWASP Category | Applies? | Specification |
|---|---|---|
| A01: Broken Access Control | No | No protected resources |
| A02: Cryptographic Failures | No | No secrets or encrypted data |
| A03: Injection | **Minimal** | DOM-based XSS possible if `innerHTML` used — see Section 3.1 |
| A04: Insecure Design | No | Simple math utility, inherently safe design |
| A05: Security Misconfiguration | **Yes** | CSP, security headers required — see Section 3.2 |
| A06: Vulnerable Components | **Yes** | npm supply chain — see Section 3.3 |
| A07: Auth Failures | No | No authentication |
| A08: Data Integrity Failures | **Yes** | Build pipeline integrity — see Section 3.4 |
| A09: Logging & Monitoring | No | No backend to monitor |
| A10: SSRF | No | No server-side requests |

---

## 3. Security Requirements

### 3.1 Input Handling — Priority: MEDIUM

All user input fields must treat input as untrusted, even though values are expected to be numeric.

**Requirements:**

| ID | Requirement | Rationale |
|---|---|---|
| INP-01 | All DOM rendering of user-provided values MUST use `textContent`, never `innerHTML` | Prevents DOM-based XSS |
| INP-02 | Input values MUST be validated as numeric before conversion processing | Rejects non-numeric payloads |
| INP-03 | Use `Number()` or `parseFloat()` for parsing — these naturally reject script content | Type coercion as defense layer |
| INP-04 | Input fields MUST have `inputmode="decimal"` attribute on mobile | Constrains keyboard to numeric, reduces attack surface |
| INP-05 | Input MUST be bounded — reject `Infinity`, `NaN`, and values exceeding `Number.MAX_SAFE_INTEGER` | Prevents edge-case rendering issues |
| INP-06 | Non-numeric characters MUST be stripped or ignored (allow: digits `0-9`, decimal `.`, negative `-`) | Defense in depth |

**Reference implementation:**

```typescript
function sanitizeNumericInput(value: string): number | null {
  const trimmed = value.trim();
  if (trimmed === '' || trimmed === '-') return null;

  // Strip anything that isn't a digit, decimal, or leading negative
  const cleaned = trimmed.replace(/[^0-9.\-]/g, '');
  const num = Number(cleaned);

  if (!Number.isFinite(num)) return null;
  if (Math.abs(num) > Number.MAX_SAFE_INTEGER) return null;

  return num;
}
```

**Anti-patterns to enforce:**

```typescript
// FORBIDDEN — XSS vector
resultElement.innerHTML = `${userInput} converts to ${result}`;

// REQUIRED — safe rendering
resultElement.textContent = convertedValue.toFixed(precision);
```

---

### 3.2 Content Security Policy & Headers — Priority: HIGH

This is the single highest-value security measure for the project. A strict CSP prevents script injection even if other defenses fail.

**Requirements:**

| ID | Requirement | Rationale |
|---|---|---|
| CSP-01 | `index.html` MUST include a CSP meta tag in `<head>` before any scripts | Baseline protection for all hosting environments |
| CSP-02 | `script-src` MUST be `'self'` only — no `'unsafe-inline'`, no `'unsafe-eval'` | Blocks injected scripts |
| CSP-03 | `style-src` MUST be `'self' https://fonts.googleapis.com` (add `'unsafe-inline'` only if Tailwind requires it) | Restricts stylesheets to known sources |
| CSP-04 | `font-src` MUST be `'self' https://fonts.gstatic.com` | Allows Google Fonts only |
| CSP-05 | `connect-src` MUST be `'none'` | No network requests needed — blocks data exfiltration |
| CSP-06 | `object-src` MUST be `'none'` | Blocks Flash/plugin-based attacks |
| CSP-07 | `base-uri` MUST be `'self'` | Prevents base tag injection |
| CSP-08 | `img-src` MUST be `'self' data:` | Allows inline SVGs via data URIs only |

**CSP meta tag (place in `<head>` of `index.html`):**

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

**Additional security headers** (via hosting platform config — `_headers` for Netlify, `vercel.json` for Vercel):

| Header | Value | Purpose |
|---|---|---|
| `X-Content-Type-Options` | `nosniff` | Prevents MIME-type sniffing |
| `X-Frame-Options` | `DENY` | Prevents clickjacking |
| `Referrer-Policy` | `strict-origin-when-cross-origin` | Limits referrer leakage |
| `Permissions-Policy` | `camera=(), microphone=(), geolocation=(), payment=()` | Disables unused browser APIs |

**Hosting headers file (`public/_headers` for Netlify):**

```
/*
  X-Content-Type-Options: nosniff
  X-Frame-Options: DENY
  Referrer-Policy: strict-origin-when-cross-origin
  Permissions-Policy: camera=(), microphone=(), geolocation=(), payment=()
```

> **Note**: The CSP is also set via meta tag (CSP-01) as a fallback for hosts that don't support custom headers (e.g., GitHub Pages).

---

### 3.3 Supply Chain & Dependency Security — Priority: HIGH

npm supply chain attacks are the most realistic threat vector for this project.

**Requirements:**

| ID | Requirement | Rationale |
|---|---|---|
| DEP-01 | Keep runtime dependencies at **zero** — all npm packages must be devDependencies only | Minimizes attack surface in production bundle |
| DEP-02 | A lockfile (`package-lock.json` or `pnpm-lock.yaml`) MUST be committed and used for installs | Ensures deterministic builds |
| DEP-03 | Run `npm audit` before each deployment — zero high/critical vulnerabilities allowed | Catches known vulnerabilities |
| DEP-04 | Pin exact dependency versions (no `^` or `~` ranges) OR rely strictly on the lockfile | Prevents unexpected version upgrades |
| DEP-05 | Periodically review dependency tree — `npm ls --all` should be small and expected | Awareness of transitive dependencies |

**Minimal dependency set** (per architecture research):

- **Dev only**: `vite`, `typescript`, `tailwindcss`, `postcss`, `autoprefixer`
- **Runtime**: Zero packages — all conversion logic is hand-written

This is ideal. The fewer dependencies, the smaller the supply chain attack surface.

---

### 3.4 Build & Deployment Integrity — Priority: MEDIUM

**Requirements:**

| ID | Requirement | Rationale |
|---|---|---|
| BLD-01 | Production builds MUST disable source maps (`build.sourcemap: false` in `vite.config.ts`) | Source maps expose original TypeScript, aiding attackers |
| BLD-02 | HTTPS MUST be enforced on the deployment platform | Prevents MITM script injection |
| BLD-03 | Build output (`dist/`) MUST NOT be committed to the repository | Prevents serving stale/tampered builds |
| BLD-04 | If using SRI (Subresource Integrity) for any CDN-loaded resources, hashes MUST be verified | Ensures CDN payloads haven't been tampered with |

**Vite config:**

```typescript
// vite.config.ts
export default defineConfig({
  build: {
    sourcemap: false,  // BLD-01: No source maps in production
  }
});
```

---

### 3.5 Client-Side Storage — Priority: LOW

The app may store user preferences (last-used category, theme preference) in `localStorage`.

**Requirements:**

| ID | Requirement | Rationale |
|---|---|---|
| STR-01 | Only non-sensitive preference data may be stored in `localStorage` | No PII or secrets in client storage |
| STR-02 | Stored values MUST be validated when read back (don't trust localStorage blindly) | Another script on same origin could tamper with values |
| STR-03 | Use `JSON.parse()` inside a try/catch when reading stored data | Graceful handling of corrupted storage |

**Example:**

```typescript
function loadPreference<T>(key: string, fallback: T): T {
  try {
    const raw = localStorage.getItem(key);
    if (raw === null) return fallback;
    const parsed = JSON.parse(raw);
    // Validate shape/type before returning
    return isValidPreference(parsed) ? parsed : fallback;
  } catch {
    return fallback;
  }
}
```

---

### 3.6 Third-Party Assets — Priority: LOW

The architecture specifies Google Fonts (Space Grotesk + Inter).

**Requirements:**

| ID | Requirement | Rationale |
|---|---|---|
| EXT-01 | Font loading MUST use `font-display: swap` to prevent FOIT | UX + prevents blank screen if CDN is slow/blocked |
| EXT-02 | Consider self-hosting fonts in `public/fonts/` to eliminate the external dependency | Removes third-party request, improves privacy, eliminates CDN compromise risk |
| EXT-03 | If loading from Google CDN, font URLs MUST be covered by the CSP `font-src` and `style-src` directives | Ensures CSP doesn't break font loading |

**Recommendation**: Self-hosting fonts is preferred for a premium app. It eliminates a third-party dependency, removes a privacy/tracking concern, and improves loading performance (no DNS lookup to `fonts.googleapis.com`). If self-hosted, simplify CSP to remove Google domains.

---

## 4. Security Testing Checklist

### Manual Tests

| # | Test Case | Expected Result | Priority |
|---|---|---|---|
| T-01 | Type `<script>alert(1)</script>` into input field | No script execution; input treated as NaN/invalid | High |
| T-02 | Type `"><img src=x onerror=alert(1)>` into input | Safe rendering, no image element created | High |
| T-03 | Paste `javascript:alert(1)` into input | Treated as non-numeric, no execution | High |
| T-04 | Enter `Infinity`, `NaN`, `-Infinity` | Graceful handling — show "Invalid" or empty result | Medium |
| T-05 | Enter extremely large number (`1e308`) | No crash, handles overflow gracefully | Medium |
| T-06 | Enter Unicode/emoji (`🔢 42`) | Non-numeric chars stripped, processes `42` | Low |
| T-07 | Verify CSP meta tag present in production `dist/index.html` | CSP tag exists with correct directives | High |
| T-08 | Verify no source maps in `dist/` folder | No `.map` files present | Medium |
| T-09 | Attempt to embed app in an iframe on another domain | Should be blocked by X-Frame-Options | Medium |
| T-10 | Check `localStorage` contents after use | Only non-sensitive preference keys present | Low |

### Automated / CI Checks

| # | Check | Tool | Priority |
|---|---|---|---|
| A-01 | `npm audit --audit-level=high` passes with 0 issues | npm | High |
| A-02 | No `innerHTML` usage in source code (grep check) | grep / lint rule | High |
| A-03 | Source maps disabled in production build | Check vite config / build output | Medium |
| A-04 | CSP meta tag present in built HTML | grep on `dist/index.html` | High |
| A-05 | Lockfile present and committed | git check | Medium |

---

## 5. Risk Summary Matrix

| Issue | Severity | Likelihood | Overall Priority | Mitigation |
|---|---|---|---|---|
| Missing CSP headers | Medium | High (no CSP by default) | **HIGH** | Add CSP meta tag (Section 3.2) |
| npm supply chain compromise | High | Low | **HIGH** | Zero runtime deps, lockfile, audit (Section 3.3) |
| DOM XSS via innerHTML | Medium | Low | **MEDIUM** | Use textContent only, validate input (Section 3.1) |
| Missing security response headers | Low | Medium | **MEDIUM** | Add headers via hosting config (Section 3.2) |
| Source maps in production | Low | Medium | **LOW** | Disable in Vite config (Section 3.4) |
| Third-party font CDN dependency | Low | Very Low | **LOW** | Consider self-hosting (Section 3.6) |
| localStorage tampering | Very Low | Very Low | **LOW** | Validate on read (Section 3.5) |

---

## 6. Implementation Priorities

### Must-Have (Before First Deploy)

1. **CSP meta tag** in `index.html` (CSP-01 through CSP-08)
2. **Zero use of `innerHTML`** for any user-provided data (INP-01)
3. **Numeric input validation** with `sanitizeNumericInput()` (INP-02, INP-03, INP-05)
4. **Source maps disabled** in production build (BLD-01)
5. **Lockfile committed** (DEP-02)
6. **`npm audit` clean** (DEP-03)

### Should-Have (Before Public Launch)

7. **Security response headers** via hosting config (X-Frame-Options, X-Content-Type-Options, etc.)
8. **`inputmode="decimal"`** on input fields (INP-04)
9. **localStorage validation** on read (STR-02, STR-03)

### Nice-to-Have (Enhancement)

10. **Self-host fonts** to eliminate Google CDN dependency (EXT-02)
11. **CI automated checks** for innerHTML grep, npm audit, CSP presence (A-01 through A-05)

---

## 7. Delegation

```
DELEGATE:
- architect: Integrate CSP meta tag into base HTML template; configure security headers for chosen hosting platform; ensure Vite config disables source maps in production
- senior_dev: Use textContent exclusively (never innerHTML) for all DOM rendering; implement sanitizeNumericInput() as the single entry point for all user input; add inputmode="decimal" to input elements; validate localStorage reads with try/catch and type checking
```

---

PHASE_COMPLETE: planning

