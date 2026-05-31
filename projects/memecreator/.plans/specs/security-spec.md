# Security Specification: Browser-Based Meme Creator

**Role:** Security Expert
**Date:** 2026-03-16
**Phase:** Planning
**References:** `.plans/research/security_expert.md`, `.plans/research/architecture-research.md`

---

## 1. Threat Assessment Summary

### Architecture Security Posture

This is a **100% client-side single-page application** with no backend, no authentication, no data storage, and no server communication. This dramatically reduces the attack surface compared to a typical web app.

**Data flow:**
```
User's browser → Local file upload → Canvas rendering → Blob export → Local download
```

No user data ever leaves the browser. This is a core security feature and privacy advantage that must be preserved.

### Threat Actors

| Actor | Risk Level | Primary Attack Vector |
|-------|-----------|----------------------|
| Malicious file crafter | Medium | SVG with embedded JS, polyglot files, oversized images |
| XSS attacker | Medium | Script injection via text input if rendered in DOM |
| Supply chain attacker | Low-Medium | Compromised npm dependency |
| Curious user / researcher | Low | Client-side code inspection (acceptable — no secrets to protect) |

---

## 2. Security Requirements

### 2.1 File Upload Validation (CRITICAL)

Per research finding #1 and #2, file uploads are the highest-risk attack vector.

**Requirements:**

| ID | Requirement | Priority |
|----|-------------|----------|
| SEC-FU-01 | Whitelist allowed MIME types: `image/jpeg`, `image/png`, `image/webp`, `image/gif` | Critical |
| SEC-FU-02 | Whitelist allowed extensions: `.jpg`, `.jpeg`, `.png`, `.webp`, `.gif` | Critical |
| SEC-FU-03 | **Block SVG uploads entirely** — SVGs can contain `<script>`, `onload`, and external resource references | Critical |
| SEC-FU-04 | Validate file magic bytes (signatures) — MIME type and extension can be spoofed | Critical |
| SEC-FU-05 | Enforce max file size: **10 MB** (10 * 1024 * 1024 bytes) | High |
| SEC-FU-06 | Enforce max image dimensions: **4096 x 4096 pixels** (prevents canvas memory exhaustion) | High |
| SEC-FU-07 | Validate file is a renderable image using `new Image()` or `createImageBitmap()` before placing on canvas | High |
| SEC-FU-08 | Apply same validation to drag-and-drop uploads and clipboard paste uploads | High |

**Implementation — Magic Byte Signatures:**

| Format | Magic Bytes (hex) |
|--------|------------------|
| JPEG | `FF D8 FF` |
| PNG | `89 50 4E 47` |
| WebP | `52 49 46 46` at offset 0, `57 45 42 50` at offset 8 |
| GIF | `47 49 46 38` |

**Validation Chain (must execute in order):**
1. Check `file.size <= 10 * 1024 * 1024`
2. Check `file.type` against MIME whitelist
3. Check file extension against extension whitelist
4. Read first 12 bytes via `FileReader` / `ArrayBuffer` and validate magic bytes
5. Load into `new Image()` — confirm `onload` fires (not `onerror`)
6. Check `image.naturalWidth <= 4096 && image.naturalHeight <= 4096`
7. Only then draw to canvas

If any step fails, reject the file with a user-friendly error message (do not expose technical details about which validation failed — just "Unsupported file format" or "Image too large").

### 2.2 Text Input & XSS Prevention (HIGH)

Per research finding #3, text inputs rendered in the DOM are an XSS risk.

**Requirements:**

| ID | Requirement | Priority |
|----|-------------|----------|
| SEC-XSS-01 | All user-entered text rendered in DOM must use `textContent`, never `innerHTML` | Critical |
| SEC-XSS-02 | Canvas `fillText()` / `strokeText()` for rendering text on the canvas (inherently safe — renders pixels) | High |
| SEC-XSS-03 | If `contenteditable` is used for inline text editing, sanitize output before any DOM insertion | High |
| SEC-XSS-04 | Never use `eval()`, `new Function()`, `document.write()`, or template literal injection with user text | Critical |
| SEC-XSS-05 | Text input fields should have `maxlength` attribute (suggest 500 characters) to prevent abuse | Medium |
| SEC-XSS-06 | If text values are stored in data attributes, use `element.dataset.prop = value` (safe) not attribute string concatenation | High |

**Safe Patterns:**
```
SAFE:   element.textContent = userText
SAFE:   ctx.fillText(userText, x, y)
SAFE:   element.dataset.label = userText

UNSAFE: element.innerHTML = userText
UNSAFE: element.setAttribute('onclick', userText)
UNSAFE: eval(userText)
```

### 2.3 Canvas Security (MEDIUM)

**Requirements:**

| ID | Requirement | Priority |
|----|-------------|----------|
| SEC-CV-01 | Only load images from the user's local filesystem via `FileReader` + `URL.createObjectURL()` — no cross-origin image loading in v1 | High |
| SEC-CV-02 | If cross-origin images are ever supported, always set `img.crossOrigin = "anonymous"` and verify CORS headers | Medium |
| SEC-CV-03 | Call `URL.revokeObjectURL()` after image is loaded to free memory and prevent URL reuse | Medium |
| SEC-CV-04 | Canvas export via `toBlob()` / `toDataURL()` naturally strips EXIF metadata — verify this in testing | Medium |

### 2.4 Export & Download Security (MEDIUM)

**Requirements:**

| ID | Requirement | Priority |
|----|-------------|----------|
| SEC-EX-01 | Sanitize download filename — allow only alphanumeric, hyphens, underscores. Recommended pattern: `meme_[timestamp].png` | Medium |
| SEC-EX-02 | Reject any path traversal characters in filename (`..`, `/`, `\`) | Medium |
| SEC-EX-03 | Set correct MIME type on export Blob: `image/png` or `image/jpeg` | Medium |
| SEC-EX-04 | Use `URL.createObjectURL(blob)` for download link — this is safe and avoids data URI size limits | Medium |
| SEC-EX-05 | Call `URL.revokeObjectURL()` after download completes to free memory | Low |

**Secure Export Pattern:**
```
1. canvas.toBlob(callback, 'image/png')
2. const url = URL.createObjectURL(blob)
3. Create <a> with download="meme_1710590400.png" and href=url
4. Programmatic click
5. URL.revokeObjectURL(url)
```

### 2.5 Dependency Supply Chain (MEDIUM)

Per architecture research: the app uses Vite + Tailwind + TypeScript, no framework.

**Requirements:**

| ID | Requirement | Priority |
|----|-------------|----------|
| SEC-DEP-01 | Minimize dependencies — prefer vanilla JS/Canvas API over libraries where possible | High |
| SEC-DEP-02 | Pin exact dependency versions in `package-lock.json` | High |
| SEC-DEP-03 | Run `npm audit` before any deployment | High |
| SEC-DEP-04 | If loading Google Fonts via external `<link>`, use Subresource Integrity (SRI) where supported | Medium |
| SEC-DEP-05 | No CDN-loaded JavaScript — all JS should be bundled locally via Vite | High |
| SEC-DEP-06 | Review any new dependency for: maintenance status, known vulnerabilities, download count, scope of permissions | Medium |

### 2.6 Content Security Policy (MEDIUM)

When deployed/hosted, the app must serve proper security headers.

**Required CSP Header:**
```
Content-Security-Policy:
  default-src 'self';
  script-src 'self';
  style-src 'self' 'unsafe-inline' https://fonts.googleapis.com;
  font-src 'self' https://fonts.gstatic.com;
  img-src 'self' blob: data:;
  connect-src 'self';
  object-src 'none';
  base-uri 'self';
  frame-ancestors 'none';
```

**Notes:**
- `img-src blob: data:` — required for canvas export and local file display
- `style-src 'unsafe-inline'` — needed for dynamic text styling on canvas; consider nonce-based approach if possible
- `object-src 'none'` — blocks Flash, Java, and other plugin content
- `frame-ancestors 'none'` — prevents clickjacking (equivalent to `X-Frame-Options: DENY`)

**Additional Security Headers:**

| Header | Value | Purpose |
|--------|-------|---------|
| `X-Content-Type-Options` | `nosniff` | Prevents MIME-type sniffing |
| `X-Frame-Options` | `DENY` | Clickjacking protection (legacy fallback) |
| `Referrer-Policy` | `strict-origin-when-cross-origin` | Limits referrer leakage |
| `Permissions-Policy` | `camera=(), microphone=(), geolocation=()` | Disables unnecessary browser APIs |

---

## 3. OWASP Top 10 Applicability

| OWASP Category | Applies? | Implementation |
|---|---|---|
| A01: Broken Access Control | No | No auth/server — N/A for client-only |
| A02: Cryptographic Failures | No | No sensitive data stored |
| **A03: Injection** | **Yes** | SEC-XSS-01 through SEC-XSS-06 |
| **A04: Insecure Design** | **Yes** | SEC-FU-01 through SEC-FU-08 (file validation by design) |
| **A05: Security Misconfiguration** | **Yes** | CSP headers, security headers on deployment |
| **A06: Vulnerable Components** | **Yes** | SEC-DEP-01 through SEC-DEP-06 |
| A07: Auth Failures | No | No authentication |
| A08: Data Integrity Failures | No | No server-side processing |
| A09: Logging & Monitoring | No | Client-side only |
| A10: SSRF | No | No server |

---

## 4. Risk Register

### Critical

| ID | Issue | Severity | Likelihood | Mitigation |
|----|-------|----------|------------|------------|
| RISK-01 | SVG upload with embedded JavaScript | Critical | Medium | Block SVG uploads entirely (SEC-FU-03) |
| RISK-02 | File type validation bypass via spoofed extension/MIME | High | Medium | Magic byte validation (SEC-FU-04) |
| RISK-03 | XSS via `innerHTML` with user text | Critical | Medium | Enforce `textContent` only (SEC-XSS-01) |

### High

| ID | Issue | Severity | Likelihood | Mitigation |
|----|-------|----------|------------|------------|
| RISK-04 | Memory exhaustion via oversized image | High | Medium | File size + dimension limits (SEC-FU-05, SEC-FU-06) |
| RISK-05 | Compromised npm dependency | High | Low | Minimize deps, pin versions, audit (SEC-DEP-01–03) |

### Medium

| ID | Issue | Severity | Likelihood | Mitigation |
|----|-------|----------|------------|------------|
| RISK-06 | Canvas tainting from cross-origin images | Medium | Low | Local files only in v1 (SEC-CV-01) |
| RISK-07 | EXIF metadata privacy leakage | Medium | Medium | Canvas re-encoding strips EXIF — verify (SEC-CV-04) |
| RISK-08 | Missing CSP allows injected scripts | Medium | Low | Deploy with strict CSP (Section 2.6) |

### Low

| ID | Issue | Severity | Likelihood | Mitigation |
|----|-------|----------|------------|------------|
| RISK-09 | Download filename path traversal | Low | Low | Sanitize/hardcode filenames (SEC-EX-01, SEC-EX-02) |
| RISK-10 | Clipboard paste with unexpected content | Low | Low | Apply same file validation to paste (SEC-FU-08) |

---

## 5. Security Testing Plan

### Manual Test Cases

| Test | Expected Result | Validates |
|------|----------------|-----------|
| Upload `.svg` file with `<script>alert(1)</script>` | File rejected with error message | SEC-FU-03 |
| Upload `.html` file renamed to `.png` | Rejected by magic byte check | SEC-FU-04 |
| Upload 50 MB image | Rejected before processing | SEC-FU-05 |
| Upload 20000x20000px PNG | Rejected after dimension check | SEC-FU-06 |
| Type `<script>alert(1)</script>` as text overlay | Renders as literal text on canvas, no execution | SEC-XSS-01 |
| Type `<img src=x onerror=alert(1)>` as text | Renders as literal text, no DOM injection | SEC-XSS-01 |
| Type extremely long text (10000+ chars) | Truncated or rejected by maxlength | SEC-XSS-05 |
| Export image and inspect with `exiftool` | No EXIF/GPS metadata in exported file | SEC-CV-04 |
| Check download filename for path traversal | Filename is sanitized, no `../` in path | SEC-EX-01 |
| Drag-and-drop an SVG file onto upload zone | Rejected same as file input | SEC-FU-08 |
| Paste image from clipboard | Same validation applied as file upload | SEC-FU-08 |

### Automated Test Cases (Unit Tests)

```
describe('FileValidator', () => {
  it('rejects SVG files')
  it('rejects files exceeding 10MB')
  it('rejects files with spoofed MIME type (html renamed to png)')
  it('rejects images exceeding 4096x4096 dimensions')
  it('accepts valid JPEG with correct magic bytes')
  it('accepts valid PNG with correct magic bytes')
  it('accepts valid WebP with correct magic bytes')
  it('accepts valid GIF with correct magic bytes')
})

describe('TextSanitization', () => {
  it('renders script tags as literal text')
  it('renders HTML entities as literal text')
  it('handles empty string input')
  it('handles extremely long input')
  it('handles unicode/emoji input')
})

describe('ExportSecurity', () => {
  it('generates sanitized filename without path traversal')
  it('produces valid PNG blob with correct MIME type')
  it('exported image contains no EXIF metadata')
  it('revokes object URL after download')
})
```

### CI/CD Security Checks

- `npm audit` — fail build on critical/high vulnerabilities
- Lighthouse security audit on deployed page
- CSP violation reporting endpoint (if hosting supports it)

---

## 6. Architecture Security Constraints

### v1 Invariants (Must Not Change Without Security Review)

1. **No server-side processing** — all computation stays in the browser
2. **No user data transmission** — images never leave the client
3. **No external image loading** — only local file uploads
4. **No dynamic script loading** — all JS bundled at build time
5. **No user authentication** — no accounts, no sessions, no tokens

### If Backend Is Added (v2+) — Triggers Full Security Review

Adding any of these features would dramatically expand the threat model and require a new security assessment:
- Server-side image storage or processing
- User accounts or authentication
- Meme sharing or social features
- Template library loaded from API
- Analytics that transmit user data

---

## 7. Implementation Module: `fileValidator.ts`

The file validation logic should be implemented as a dedicated, testable module. Recommended structure:

```
src/
  security/
    fileValidator.ts    → All upload validation logic (SEC-FU-01 through SEC-FU-08)
    sanitize.ts         → Filename sanitization (SEC-EX-01, SEC-EX-02)
```

This keeps security logic isolated, testable, and easy to audit.

---

## 8. Delegation

```
DELEGATE:
- architect: Ensure file validation is a dedicated module in the architecture. Plan CSP headers and security headers in the deployment configuration. Maintain client-only invariant.
- senior_dev: Implement fileValidator.ts with the full validation chain (type + extension + magic bytes + Image load + dimensions). Use textContent exclusively for DOM text. Implement secure export with blob URLs, sanitized filenames, and URL.revokeObjectURL() cleanup. Add unit tests for all security test cases in Section 5.
```

---

PHASE_COMPLETE: planning
