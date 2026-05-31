# Security Research: Image Resizer Website

## Role: Security Expert
## Date: 2026-03-16

---

## 1. Project Context

- **Type**: Client-side only image resizer (static HTML/CSS/JS)
- **Stack**: HTML5 Canvas API, File API, vanilla JS, no backend
- **Hosting**: Static files (no server-side processing)
- **Data flow**: User selects image → browser processes → user downloads result
- **Key constraint**: All processing happens in the browser. No uploads to any server.

This is a low-risk application by design. No backend, no database, no user accounts, no server-side file handling. However, client-side applications still have real attack surfaces.

---

## 2. Threat Assessment

### 2.1 Threat Actors & Motivation

| Actor | Motivation | Likelihood |
|-------|-----------|------------|
| Script kiddie | Deface or inject ads via XSS if hosted on shared domain | Low |
| Malicious advertiser | Inject tracking/crypto-mining scripts via CDN compromise | Low-Medium |
| Supply chain attacker | Compromise any external dependency to reach users | Low (if zero deps) |
| Malicious image crafter | Exploit browser image parsing bugs via crafted files | Very Low |

### 2.2 Attack Surface Summary

Since this is a pure client-side app with no backend:
- **No server-side attack surface** (no SQL injection, no SSRF, no auth bypass)
- **No user data stored server-side** (no data breaches possible from backend)
- **Primary risk**: XSS if any user-controlled content is rendered as HTML
- **Secondary risk**: Supply chain attacks via external CDN dependencies
- **Tertiary risk**: Hosting misconfiguration (missing security headers)

---

## 3. Vulnerability Analysis

### 3.1 Client-Side Vulnerabilities

#### A. Cross-Site Scripting (XSS) — OWASP A03:2021

**Risk: Medium**

Even without a backend, XSS is possible if:
- File names from uploaded images are rendered as innerHTML without sanitization
- Image metadata (EXIF) is extracted and displayed without escaping
- URL parameters are read and injected into the DOM (DOM-based XSS)

**Examples from real image tools**:
- Squoosh (Google) sanitizes all file name displays
- TinyPNG renders file names as textContent, never innerHTML
- iLoveIMG strips all metadata display to prevent injection

**Mitigations**:
- ALWAYS use `textContent` or `innerText` instead of `innerHTML` for user-derived strings
- Never insert file names, EXIF data, or URL params via `innerHTML`
- If using template literals to build HTML, escape all user values
- Use Content Security Policy (CSP) headers to block inline script execution

#### B. Malicious File Handling

**Risk: Low**

Users upload image files that are processed by the browser's native image decoder (via `<canvas>`, `<img>`, `createImageBitmap`). Risks:
- **Crafted images** could theoretically exploit browser image parsing bugs (historically: libpng, libjpeg vulnerabilities). This is a browser-level concern, not an app-level one.
- **Non-image files** disguised as images could cause unexpected behavior if not validated.
- **Very large files** could cause memory exhaustion / tab crash (DoS on the user's own browser).

**Mitigations**:
- Validate file type via both `file.type` (MIME) and magic bytes (file signature)
- Set a reasonable max file size (e.g., 50MB) and warn users before processing
- Wrap canvas operations in try/catch to handle decode failures gracefully
- Use `createImageBitmap()` for safer, async image decoding where supported

#### C. Data Leakage via EXIF/Metadata

**Risk: Medium**

Images often contain sensitive EXIF metadata:
- GPS coordinates (location where photo was taken)
- Camera/device info
- Timestamps
- Thumbnail of original image (even if cropped)

When users resize images, they may expect metadata to be stripped. Canvas API naturally strips EXIF when drawing to canvas and exporting, which is a **security benefit**.

**Mitigations**:
- Document that resized images have EXIF stripped (this is a feature, not a bug)
- If ever adding EXIF preservation as a feature, make it opt-in with clear warnings
- Never display extracted EXIF data without sanitization

#### D. Blob URL / Object URL Leaks

**Risk: Low**

When creating preview URLs via `URL.createObjectURL()`, forgetting to revoke them causes memory leaks and keeps references to file data in memory longer than necessary.

**Mitigations**:
- Always call `URL.revokeObjectURL()` when the preview is no longer needed
- Revoke on new file upload and on page unload

### 3.2 Hosting & Delivery Vulnerabilities

#### E. Missing Security Headers — OWASP A05:2021

**Risk: Medium**

Static sites still need proper HTTP headers when deployed.

**Required headers**:
```
Content-Security-Policy: default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' blob: data:; object-src 'none'; base-uri 'self'
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
Referrer-Policy: no-referrer
Permissions-Policy: camera=(), microphone=(), geolocation=()
```

**Notes**:
- `img-src` needs `blob:` and `data:` for canvas-generated previews and downloads
- `style-src 'unsafe-inline'` may be needed for brutalist inline styles — ideally move styles to external file and use `'self'` only
- `script-src 'self'` blocks inline scripts — all JS should be in external files
- `object-src 'none'` prevents Flash/plugin-based attacks

#### F. Supply Chain / Dependency Risk — OWASP A06:2021

**Risk: Low (if zero external dependencies)**

The template choice (static-html, no build tools, no node_modules) is excellent from a security perspective.

**Risks arise if**:
- External CDN fonts are loaded (Google Fonts = tracking + SPOF)
- External JS libraries are included via CDN without SRI (Subresource Integrity)
- Third-party analytics or ad scripts are added

**Mitigations**:
- Self-host any fonts used for brutalist typography
- If any external script is absolutely needed, use SRI hashes: `<script src="..." integrity="sha384-..." crossorigin="anonymous">`
- Prefer zero external dependencies (the template already recommends this)

#### G. Subdomain / Hosting Takeover

**Risk: Low**

If deployed on platforms like GitHub Pages, Netlify, Vercel:
- Ensure custom domain DNS is correctly configured (prevent subdomain takeover)
- Enable HTTPS (most platforms do this by default)

---

## 4. OWASP Top 10 (2021) Relevance Check

| OWASP Category | Relevant? | Notes |
|---------------|-----------|-------|
| A01: Broken Access Control | No | No auth, no access control needed |
| A02: Cryptographic Failures | No | No sensitive data transmitted/stored |
| A03: Injection (XSS) | **Yes** | DOM-based XSS via file names, metadata |
| A04: Insecure Design | Low | Simple design, minimal attack surface |
| A05: Security Misconfiguration | **Yes** | Missing CSP/security headers on deploy |
| A06: Vulnerable Components | Low | Only if external deps are added |
| A07: Auth Failures | No | No authentication |
| A08: Data Integrity Failures | Low | SRI for any external resources |
| A09: Logging/Monitoring | No | Client-side only, no logs to monitor |
| A10: SSRF | No | No server-side requests |

---

## 5. Competitor Security Practices

### Squoosh (Google)
- Pure client-side processing (WASM + Canvas)
- Strict CSP headers
- No external tracking scripts
- All processing in Web Workers (isolates from main thread)
- Open source, regularly audited

### TinyPNG/TinyJPG
- Server-side processing (uploads to their API)
- API key authentication for programmatic access
- Rate limiting on free tier
- HTTPS enforced
- Privacy policy: files deleted after processing

### iLoveIMG
- Server-side processing
- File upload with size limits
- Auto-deletion of files after 2 hours
- GDPR compliant privacy policy

### Photopea
- Mostly client-side (Canvas/WebGL)
- Minimal external dependencies
- CSP headers configured
- Handles untrusted file formats defensively

### Key takeaway for our project
Client-side processing (like Squoosh) is the gold standard for privacy and security. Our approach aligns with this. The main things to get right are: CSP headers, no innerHTML with user data, and defensive file handling.

---

## 6. Recommendations Summary

| # | Issue | Severity | Likelihood | Recommendation | Priority |
|---|-------|----------|------------|----------------|----------|
| 1 | DOM XSS via file names/metadata | High | Medium | Use textContent only, never innerHTML for user data | **Critical** |
| 2 | Missing CSP headers | Medium | High | Configure strict CSP on deployment | **High** |
| 3 | No file type validation | Medium | Medium | Validate MIME type + magic bytes, enforce size limit | **High** |
| 4 | Blob URL memory leaks | Low | High | Revoke object URLs when no longer needed | **Medium** |
| 5 | External dependency risk | Medium | Low | Zero external CDN deps; self-host fonts; use SRI if needed | **Medium** |
| 6 | EXIF data leakage awareness | Low | Medium | Document that Canvas strips EXIF (a privacy win) | **Low** |
| 7 | Large file DoS (self-DoS) | Low | Medium | Max file size check before processing | **Low** |

---

## 7. Secure Implementation Guidelines for Developers

### Must-Do (Before First Deploy)

1. **Never use `innerHTML` with any user-derived value** — file names, dimensions, metadata. Always `textContent`.
2. **Validate uploaded files**:
   ```
   - Check file.type starts with 'image/'
   - Check file size < MAX_SIZE (suggest 50MB)
   - Wrap image loading in error handlers
   ```
3. **Revoke blob URLs** — call `URL.revokeObjectURL(url)` after use.
4. **All JS in external files** — enables strict CSP with no `'unsafe-inline'` for scripts.
5. **Deploy with security headers** — CSP, X-Content-Type-Options, X-Frame-Options.

### Should-Do (Best Practice)

6. **Use Web Workers for processing** — isolates heavy canvas work from main thread, slight security boundary.
7. **Self-host fonts** — avoid Google Fonts CDN for privacy and reliability.
8. **Add SRI hashes** if any external resource is ever included.
9. **Sanitize download file names** — when generating the output file name, strip any characters that could be problematic in file systems (e.g., `../`, null bytes, control characters).

### Nice-to-Have

10. **Feature-Policy/Permissions-Policy** — explicitly deny camera, mic, geolocation.
11. **Add a simple privacy statement** — "Your images never leave your browser. All processing happens locally."

---

## 8. Security Testing Suggestions

| Test | Method | What to Check |
|------|--------|---------------|
| XSS via file name | Upload file named `<img src=x onerror=alert(1)>.png` | File name should render as text, not HTML |
| XSS via URL params | Add `?name=<script>alert(1)</script>` to URL | No script execution |
| File type bypass | Upload a .html file renamed to .png | Should reject or fail gracefully |
| Large file handling | Upload a 200MB image | Should warn/reject before processing, not crash |
| Blob URL cleanup | Upload multiple images sequentially | Memory should not grow unboundedly |
| CSP validation | Use browser DevTools / securityheaders.com | All recommended headers present |
| Output file name | Upload file with `../../etc/passwd` as name | Output name should be sanitized |

---

## 9. Architecture Security Notes

The chosen architecture (static HTML, client-side only, zero dependencies) is **inherently secure by design**:

- **No network attack surface** — no API, no server, no database
- **No data leaves the browser** — privacy by architecture
- **No supply chain** — no npm, no build tools, no CDN dependencies to compromise
- **Minimal code** — less code = fewer bugs = fewer vulnerabilities

This is about as secure as a web application can get. The remaining risks are all client-side hygiene issues that are straightforward to address during implementation.

---

## 10. Delegation

```
DELEGATE:
- architect: Ensure CSP header strategy is included in deployment config; confirm Web Worker isolation approach if adopted
- senior_dev: Implement file validation (MIME + magic bytes + size limit), use textContent exclusively for user data rendering, implement blob URL lifecycle management
- junior_dev: No security-sensitive tasks should be delegated without review
```

---

**Overall Risk Rating: LOW** — Client-side architecture with zero dependencies is inherently low-risk. Main focus should be on XSS prevention via proper DOM API usage and security headers at deploy time.

---
Status: PHASE_COMPLETE: research
