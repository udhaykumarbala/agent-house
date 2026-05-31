# Security Specification: Image Resizer

## Role: Security Expert
## Date: 2026-03-16
## Status: Planning Phase

---

## 1. Threat Model Summary

### Architecture Context
- **Type**: Static client-side SPA (HTML/CSS/JS)
- **Processing**: Browser-only via Canvas API — no server, no database, no API
- **Data flow**: `File Input → FileReader → Image → Canvas → Blob → Download`
- **External dependencies**: Google Fonts (CDN) — the only external resource

### Overall Risk: LOW

The pure client-side architecture eliminates entire categories of server-side vulnerabilities. No data leaves the browser. The remaining attack surface is limited to client-side DOM manipulation and hosting configuration.

### Threat Actors

| Actor | Target | Likelihood | Impact |
|-------|--------|------------|--------|
| XSS attacker | DOM injection via file names or URL params | Medium | High (session context) |
| Supply chain attacker | Compromised CDN resource (Google Fonts) | Low | Medium |
| Hosting misconfiguration | Missing security headers enabling framing/sniffing | Medium | Medium |
| Malicious file crafter | Browser image decoder exploits | Very Low | Low (browser-level) |

---

## 2. Security Requirements

### 2.1 DOM Security — CRITICAL

**Rule: Never use `innerHTML`, `outerHTML`, or `insertAdjacentHTML` with any user-derived value.**

All user-derived strings must be rendered via safe DOM APIs only.

| Data Source | Rendering Method | Forbidden |
|-------------|-----------------|-----------|
| File name (`file.name`) | `element.textContent` | `innerHTML` |
| Image dimensions | `element.textContent` | Template literal → `innerHTML` |
| File size | `element.textContent` | `document.write()` |
| Any URL parameter | `element.textContent` | `innerHTML`, `eval()` |
| Error messages | `element.textContent` | `innerHTML` |

**Implementation rules:**
```
REQUIRED:
  element.textContent = fileName;
  element.setAttribute('download', sanitizedName);

FORBIDDEN:
  element.innerHTML = `<span>${fileName}</span>`;
  container.innerHTML = `Size: ${width} x ${height}`;
  document.write(anything);
  eval(anything);
```

**Why this matters**: Even without a backend, DOM-based XSS is possible if file names like `<img src=x onerror=alert(1)>.png` are rendered as HTML. This is the #1 risk for this application.

### 2.2 File Input Validation — HIGH

All uploaded files must be validated before processing.

**Validation chain (in order):**

1. **MIME type check**: `file.type` must start with `image/`
2. **Extension check**: Must be one of `.jpg`, `.jpeg`, `.png`, `.webp`, `.gif`, `.bmp`, `.svg`
3. **File size check**: Must be ≤ 50MB. Warn at > 20MB.
4. **Magic bytes check**: Validate file signature matches claimed type

**Magic bytes reference:**

| Format | Magic Bytes (hex) |
|--------|------------------|
| JPEG | `FF D8 FF` |
| PNG | `89 50 4E 47 0D 0A 1A 0A` |
| GIF | `47 49 46 38` |
| WebP | `52 49 46 46` (RIFF) + `57 45 42 50` at offset 8 |
| BMP | `42 4D` |

**Implementation approach:**
```
function validateFile(file) {
  // 1. MIME type
  if (!file.type.startsWith('image/')) → reject with user message

  // 2. Extension
  const ext = file.name.split('.').pop().toLowerCase();
  if (!ALLOWED_EXTENSIONS.includes(ext)) → reject with user message

  // 3. Size
  if (file.size > MAX_FILE_SIZE) → reject with user message
  if (file.size > WARN_FILE_SIZE) → show warning, allow proceed

  // 4. Magic bytes (read first 12 bytes via FileReader/ArrayBuffer)
  const header = await readFileHeader(file, 12);
  if (!matchesMagicBytes(header, file.type)) → reject with user message
}
```

**Constants:**
- `MAX_FILE_SIZE`: 50 * 1024 * 1024 (50MB)
- `WARN_FILE_SIZE`: 20 * 1024 * 1024 (20MB)
- `ALLOWED_EXTENSIONS`: `['jpg', 'jpeg', 'png', 'webp', 'gif', 'bmp']`
- Note: SVG is intentionally excluded from initial support due to script injection risks in SVG files

**Error handling:**
- Wrap all Canvas/Image operations in try/catch
- Use `createImageBitmap()` where supported for safer async decoding
- Display user-friendly error on decode failure: "This file couldn't be processed. It may be corrupted or not a supported image format."

### 2.3 Download File Name Sanitization — HIGH

When generating the output file name for download, sanitize the original file name.

**Sanitization rules:**
```
function sanitizeFileName(name) {
  // 1. Remove path traversal sequences
  name = name.replace(/\.\.\//g, '').replace(/\.\.\\/g, '');

  // 2. Remove null bytes and control characters
  name = name.replace(/[\x00-\x1F\x7F]/g, '');

  // 3. Remove characters problematic in file systems
  name = name.replace(/[<>:"/\\|?*]/g, '');

  // 4. Trim whitespace and dots from edges
  name = name.replace(/^[\s.]+|[\s.]+$/g, '');

  // 5. Fallback if name is empty after sanitization
  if (!name || name.length === 0) {
    name = 'resized-image';
  }

  // 6. Truncate to reasonable length
  if (name.length > 200) {
    name = name.substring(0, 200);
  }

  return name;
}
```

### 2.4 Blob URL Lifecycle — MEDIUM

Object URLs created via `URL.createObjectURL()` must be revoked when no longer needed.

**Rules:**
- Store the current blob URL in a variable
- Before creating a new blob URL, revoke the previous one
- Revoke on "new image" / reset action
- Revoke on page unload (via `beforeunload` event)

```
let currentBlobUrl = null;

function setPreviewUrl(blob) {
  if (currentBlobUrl) {
    URL.revokeObjectURL(currentBlobUrl);
  }
  currentBlobUrl = URL.createObjectURL(blob);
  previewElement.src = currentBlobUrl;
}

window.addEventListener('beforeunload', () => {
  if (currentBlobUrl) URL.revokeObjectURL(currentBlobUrl);
});
```

### 2.5 SVG Exclusion — MEDIUM

SVG files must NOT be accepted as input for the initial release.

**Why**: SVG files can contain embedded `<script>` tags, event handlers (`onload`, `onerror`), and external resource references. Loading an SVG into an `<img>` tag is generally safe (scripts don't execute), but processing SVGs through other paths (e.g., `innerHTML`, `XMLSerializer`) can introduce XSS.

**Future consideration**: If SVG support is added later, SVGs must be sanitized via DOMParser with script/event handler stripping before any processing.

---

## 3. Content Security Policy (CSP)

### Recommended CSP Header

```
Content-Security-Policy:
  default-src 'self';
  script-src 'self';
  style-src 'self' https://fonts.googleapis.com;
  font-src 'self' https://fonts.gstatic.com;
  img-src 'self' blob: data:;
  connect-src 'none';
  object-src 'none';
  frame-src 'none';
  base-uri 'self';
  form-action 'none';
  frame-ancestors 'none';
```

### Directive Rationale

| Directive | Value | Why |
|-----------|-------|-----|
| `script-src 'self'` | Only local JS files | Blocks inline scripts and external script injection. ALL JavaScript must be in `.js` files, not inline `<script>` tags or `onclick` handlers. |
| `style-src 'self' https://fonts.googleapis.com` | Local CSS + Google Fonts CSS | Google Fonts serves CSS from `fonts.googleapis.com`. If we self-host fonts, tighten to `'self'` only. |
| `font-src 'self' https://fonts.gstatic.com` | Font file delivery | Google Fonts serves font files from `fonts.gstatic.com`. If self-hosted, tighten to `'self'`. |
| `img-src 'self' blob: data:` | Allow canvas-generated images | Canvas `toBlob()` and `toDataURL()` outputs need `blob:` and `data:` schemes. |
| `connect-src 'none'` | No network requests | This app makes zero fetch/XHR calls. Blocks any injected beacon/tracking. |
| `object-src 'none'` | No plugins | Blocks Flash, Java, and other plugin-based attack vectors. |
| `frame-src 'none'` | No iframes | App doesn't embed any third-party content. |
| `form-action 'none'` | No form submissions | App has no forms that submit to a server. Blocks CSRF-style form hijacking. |
| `frame-ancestors 'none'` | Cannot be framed | Prevents clickjacking. Equivalent to `X-Frame-Options: DENY`. |

### Important Implementation Note

Because `script-src 'self'` is specified:
- **No inline `<script>` blocks** in HTML — all JS must be in external `.js` files
- **No inline event handlers** — no `onclick="..."`, `onload="..."` etc. in HTML. Use `addEventListener()` in JS files.
- **No `eval()`** or `Function()` constructor
- **No `javascript:` URLs**

If the brutalist design requires inline styles, that's acceptable — `style-src` allows inline styles implicitly when using the `style` attribute. However, if `<style>` blocks in HTML are needed, add `'unsafe-inline'` to `style-src` (acceptable trade-off for a no-backend app). Alternatively, move all styles to the external CSS file (preferred).

---

## 4. Additional Security Headers

### Full Header Set for Deployment

```
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
Referrer-Policy: no-referrer
Permissions-Policy: camera=(), microphone=(), geolocation=(), payment=(), usb=()
Cross-Origin-Opener-Policy: same-origin
Cross-Origin-Embedder-Policy: require-corp
```

| Header | Purpose |
|--------|---------|
| `X-Content-Type-Options: nosniff` | Prevents MIME-type sniffing — browser respects declared content types |
| `X-Frame-Options: DENY` | Backup clickjacking protection (CSP `frame-ancestors` is primary) |
| `Referrer-Policy: no-referrer` | No referrer sent on any navigation — prevents URL leakage |
| `Permissions-Policy` | Explicitly denies access to hardware APIs the app doesn't use |
| `Cross-Origin-Opener-Policy` | Isolates the browsing context from cross-origin windows |
| `Cross-Origin-Embedder-Policy` | Ensures all subresources are same-origin or explicitly opted in |

### Deployment Configuration

These headers must be set at the hosting level. Examples for common platforms:

**Netlify** (`_headers` file):
```
/*
  Content-Security-Policy: default-src 'self'; script-src 'self'; style-src 'self' https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com; img-src 'self' blob: data:; connect-src 'none'; object-src 'none'; frame-src 'none'; base-uri 'self'; form-action 'none'; frame-ancestors 'none'
  X-Content-Type-Options: nosniff
  X-Frame-Options: DENY
  Referrer-Policy: no-referrer
  Permissions-Policy: camera=(), microphone=(), geolocation=(), payment=(), usb=()
```

**Vercel** (`vercel.json`):
```json
{
  "headers": [
    {
      "source": "/(.*)",
      "headers": [
        { "key": "Content-Security-Policy", "value": "default-src 'self'; script-src 'self'; style-src 'self' https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com; img-src 'self' blob: data:; connect-src 'none'; object-src 'none'; frame-src 'none'; base-uri 'self'; form-action 'none'; frame-ancestors 'none'" },
        { "key": "X-Content-Type-Options", "value": "nosniff" },
        { "key": "X-Frame-Options", "value": "DENY" },
        { "key": "Referrer-Policy", "value": "no-referrer" },
        { "key": "Permissions-Policy", "value": "camera=(), microphone=(), geolocation=(), payment=(), usb=()" }
      ]
    }
  ]
}
```

**GitHub Pages**: Cannot set custom headers natively. Use a `<meta>` tag fallback for CSP:
```html
<meta http-equiv="Content-Security-Policy" content="default-src 'self'; script-src 'self'; style-src 'self' https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com; img-src 'self' blob: data:; connect-src 'none'; object-src 'none'; frame-src 'none'; base-uri 'self'; form-action 'none'">
```
Note: `frame-ancestors` cannot be set via `<meta>` tag — only via HTTP header.

---

## 5. External Dependency Security

### Google Fonts Risk

Google Fonts is the only external dependency per the architecture research.

**Risks:**
- Google Fonts CDN serves CSS and font files — a compromised CDN could inject malicious CSS
- Google Fonts requests reveal user IP and browsing behavior to Google (privacy concern)
- CDN downtime = broken typography

**Recommended approach (in order of preference):**

1. **Self-host fonts** (BEST): Download Space Grotesk and Inter font files, include in project `/fonts/` directory. Eliminates CDN dependency entirely. Tighten CSP to `style-src 'self'; font-src 'self'`.

2. **Use Subresource Integrity (SRI)**: If using CDN, add `integrity` and `crossorigin` attributes. Note: Google Fonts dynamically generates CSS, making SRI hashes unstable. This approach is fragile.

3. **Use system font fallback**: `font-family: system-ui, -apple-system, 'Segoe UI', sans-serif`. Zero external requests. Loses design specificity but maximizes security and performance.

**Decision**: Defer to architect and UI team. Security recommends option 1 (self-host). If CDN is used, accept the risk as LOW given Google's infrastructure security.

### No Other Dependencies

Per the architect research, the app uses:
- No npm packages
- No build tools
- No JavaScript libraries
- No analytics scripts
- No third-party widgets

This is the ideal security posture. **Any future addition of external scripts must go through security review.**

---

## 6. Privacy Specification

### Privacy by Architecture

The client-side-only design provides inherent privacy guarantees:

- Images never leave the browser
- No telemetry, analytics, or tracking
- No cookies required
- No server to be breached

### Privacy Statement

Include a visible privacy statement on the page (aligns with UX research on trust-building):

> "Your images never leave your browser. All processing happens locally on your device. We don't upload, store, or see your files."

### EXIF Metadata Handling

The Canvas API naturally strips EXIF metadata when drawing to canvas and exporting. This is a privacy benefit — resized images won't contain GPS coordinates, device info, or timestamps from the original.

**Specification**: Do NOT add EXIF preservation. If it's ever requested as a feature, it must be opt-in with a clear warning about metadata exposure.

---

## 7. Security Testing Checklist

### Pre-Launch Tests

| # | Test | Method | Expected Result | Priority |
|---|------|--------|-----------------|----------|
| 1 | XSS via file name | Upload file named `<img src=x onerror=alert(1)>.png` | File name displays as plain text, no script execution | **Critical** |
| 2 | XSS via URL params | Navigate to `?q=<script>alert(1)</script>` | No script execution, params ignored | **Critical** |
| 3 | File type bypass | Rename `.html` file to `.png` and upload | Rejected by magic bytes validation | **High** |
| 4 | Oversized file | Upload 100MB+ file | Rejected with friendly error before processing | **High** |
| 5 | CSP validation | Check headers via browser DevTools → Network tab | All CSP directives present | **High** |
| 6 | CSP violation test | Inject `<script>alert(1)</script>` in DevTools | Blocked by CSP, violation logged in console | **High** |
| 7 | Security headers | Scan with securityheaders.com | A+ or A rating | **Medium** |
| 8 | Download name injection | Upload file named `../../etc/passwd.png` | Download name is sanitized (`etcpasswd.png` or similar) | **Medium** |
| 9 | Blob URL cleanup | Upload 10 images sequentially, check DevTools Memory | No unbounded memory growth from blob URLs | **Medium** |
| 10 | Null byte in name | Upload file with `image%00.png` name | Null byte stripped, processed normally or rejected | **Medium** |
| 11 | Zero-byte file | Upload empty file | Rejected gracefully with error message | **Low** |
| 12 | Corrupted image | Upload file with valid magic bytes but corrupted body | Canvas error caught, user sees friendly error | **Low** |

### Automated Checks

If CI/CD is ever added:
- Lint JS for `innerHTML` usage (ban via ESLint rule `no-inner-html` or custom rule)
- Scan for `eval()`, `Function()`, `document.write()` usage
- Validate CSP header syntax

---

## 8. Secure Coding Checklist for Developers

Before every PR/code change, verify:

- [ ] No `innerHTML` with user-derived values
- [ ] No inline `<script>` blocks or `onclick` handlers in HTML
- [ ] No `eval()`, `Function()`, or `document.write()`
- [ ] File validation runs before any processing
- [ ] Download file names are sanitized
- [ ] Blob URLs are revoked after use
- [ ] Error messages don't expose internal details (stack traces, file paths)
- [ ] No new external dependencies added without security review
- [ ] CSP meta tag / header is not weakened

---

## 9. Risk Register

| # | Risk | Severity | Likelihood | Current Status | Mitigation |
|---|------|----------|------------|----------------|------------|
| S1 | DOM XSS via file names | High | Medium | Mitigated by spec | Use `textContent` exclusively — §2.1 |
| S2 | Missing security headers | Medium | High (if forgotten) | Mitigated by spec | Deploy with headers — §3, §4 |
| S3 | Invalid file type processing | Medium | Medium | Mitigated by spec | Multi-layer validation — §2.2 |
| S4 | Download name injection | Medium | Low | Mitigated by spec | Name sanitization — §2.3 |
| S5 | Memory leak via blob URLs | Low | High | Mitigated by spec | Lifecycle management — §2.4 |
| S6 | SVG script injection | High | Low | Mitigated by spec | SVG excluded from input — §2.5 |
| S7 | Google Fonts CDN compromise | Medium | Very Low | Accepted risk | Self-host preferred — §5 |
| S8 | Browser image decoder exploit | High | Very Low | Accepted risk | Browser vendor responsibility |
| S9 | Clickjacking | Medium | Low | Mitigated by spec | `frame-ancestors 'none'` — §3 |

---

## 10. Delegation

```
DELEGATE:
- architect: Include CSP meta tag in HTML template if deploying to GitHub Pages; include _headers or equivalent for Netlify/Vercel deployment config; decide on self-hosting fonts vs CDN
- senior_dev: Implement file validation chain (§2.2), file name sanitization (§2.3), blob URL lifecycle (§2.4); ensure zero innerHTML usage with user data; all event handlers via addEventListener only
```

---

**Summary**: This is a low-risk application by design. The security requirements are straightforward — DOM safety, file validation, security headers, and dependency minimalism. No complex auth, encryption, or server-side hardening needed. The main discipline required is consistent use of safe DOM APIs and proper deployment headers.

---
Status: PHASE_COMPLETE: planning
