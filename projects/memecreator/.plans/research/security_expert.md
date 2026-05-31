# Security Research: Browser-Based Meme Creator

**Role:** Security Expert
**Date:** 2026-03-16
**Status:** Research Phase

---

## 1. Threat Assessment

### Attack Surface Overview

This is a **client-side single-page web app** that handles:
- User image uploads (file input)
- Canvas-based image manipulation
- Text overlay rendering
- Image export/download

### Threat Actors & Motivations

| Actor | Motivation | Likely Attack |
|-------|-----------|---------------|
| Script kiddie | Deface or exploit for fun | XSS via crafted image metadata or text input |
| Malicious user | Exploit other users (if shared) | Stored XSS, malicious file upload |
| Automated bot | Resource abuse, SEO spam | Abuse any server endpoints, DoS |
| Curious user | Data exfiltration | Inspect client-side secrets, API keys |

### Key Assumption

Since this is a **browser-only tool with no backend**, the attack surface is significantly reduced. The primary risks are client-side vulnerabilities. If a backend is added later (for saving/sharing memes), the threat model expands dramatically.

---

## 2. Vulnerability Analysis

### 2.1 Image Upload Handling

**Risk: Malicious File Upload (High)**

- Users upload arbitrary files via `<input type="file">`.
- Even without a server, malicious files can exploit the browser.
- SVG files can contain embedded JavaScript that executes on render.
- Polyglot files (e.g., GIFAR — GIF+JAR) can bypass type checks.
- Excessively large images can cause browser DoS (memory exhaustion).

**Competitor Research:**
- **Canva**: Validates file type server-side + client-side, enforces size limits (25MB for images), strips EXIF/metadata.
- **Imgflip**: Restricts to JPEG/PNG/GIF only, max 10MB, server-side validation.
- **Kapwing**: Client-side type check + server processing pipeline that re-encodes all uploads.

**Best Practice:**
- Validate MIME type AND file extension AND magic bytes (file signature).
- Restrict to safe raster formats: JPEG, PNG, WebP, GIF.
- **Block SVG uploads** — SVGs can contain `<script>` tags, `onload` handlers, and external resource references.
- Enforce maximum file size (e.g., 10MB) client-side.
- Enforce maximum image dimensions (e.g., 4096x4096) to prevent canvas memory bombs.
- Use `createImageBitmap()` or `new Image()` to validate the file is actually a renderable image before placing on canvas.

### 2.2 Text Input & XSS

**Risk: Cross-Site Scripting via Text Overlays (Medium)**

- Users type text that gets rendered on a canvas.
- Canvas 2D API `fillText()` / `strokeText()` is inherently safe — it renders text as pixels, not DOM.
- **However**, if text is also rendered in DOM elements (e.g., for a text editing UI, preview, or tooltip), XSS is possible.
- If text content is stored in `innerHTML` or injected into the DOM unsanitized, script injection can occur.

**Competitor Research:**
- **Canva**: All text editing happens in a controlled contenteditable div with sanitization. Canvas rendering is separate.
- **Imgflip**: Text inputs use standard form elements; output is canvas-only.
- **Photopea**: Complex DOM + Canvas hybrid — heavy sanitization on DOM-rendered content.

**Best Practice:**
- Use `textContent` instead of `innerHTML` for any DOM text rendering.
- Canvas `fillText()` is safe by nature — prefer rendering text directly on canvas.
- If using contenteditable elements for rich text editing, sanitize with DOMPurify before any DOM insertion.
- Never use `eval()`, `new Function()`, or `document.write()` with user-supplied text.

### 2.3 Canvas Security

**Risk: Canvas Tainting & Data Exfiltration (Medium)**

- HTML5 Canvas has a security model around "tainted" canvases.
- If an image from a different origin is drawn on canvas without CORS, `toDataURL()` and `toBlob()` will throw a `SecurityError`.
- This could break the export feature if cross-origin images are loaded.

**Best Practice:**
- Only load images from the user's local filesystem (via `FileReader` + `URL.createObjectURL()`). This avoids CORS issues entirely.
- If loading remote images in the future, use `crossOrigin = "anonymous"` and ensure the server sends proper CORS headers.
- Never load images from untrusted URLs without CORS validation.

### 2.4 Client-Side Data Handling

**Risk: Data Privacy / EXIF Leakage (Low-Medium)**

- Uploaded photos may contain EXIF metadata: GPS coordinates, device info, timestamps, camera serial numbers.
- If the exported meme retains this metadata, users may unknowingly share private location data.

**Competitor Research:**
- **Canva**: Strips all EXIF metadata on export.
- **Imgflip**: Server-side re-encoding strips metadata naturally.
- **Remove.bg / TinyPNG**: Explicitly strip metadata during processing.

**Best Practice:**
- Canvas `toDataURL()` and `toBlob()` naturally strip EXIF metadata because they re-encode the pixel data. This is a security benefit of the canvas export approach.
- Document this behavior for users — it's a privacy feature.
- If using libraries like `canvas-to-blob` or direct canvas export, verify metadata is not preserved.

### 2.5 Dependency Supply Chain

**Risk: Malicious Dependencies (Medium)**

- If the project uses npm packages (e.g., for canvas manipulation, text rendering, UI framework), supply chain attacks are a concern.
- Compromised packages can exfiltrate data, inject crypto miners, or modify exports.

**Competitor Research / Industry Incidents:**
- **event-stream incident (2018)**: Popular npm package compromised to steal cryptocurrency.
- **ua-parser-js (2021)**: Supply chain attack injected crypto miners.
- **Colors.js / Faker.js (2022)**: Maintainer sabotage.

**Best Practice:**
- Minimize dependencies — a meme creator can be built with vanilla JS + Canvas API.
- If using a framework (React, Vue, etc.), pin exact versions in `package-lock.json`.
- Use `npm audit` regularly.
- Consider using Subresource Integrity (SRI) for any CDN-loaded scripts.
- Prefer well-maintained, widely-used libraries with active security response teams.

### 2.6 Export & Download Security

**Risk: Zip Slip / Download Injection (Low)**

- The app generates an image and triggers a download.
- Using `<a download="meme.png" href="blob:...">` is safe.
- Ensure the filename cannot be manipulated to include path traversal characters.

**Best Practice:**
- Hardcode or sanitize the download filename (e.g., `meme_[timestamp].png`).
- Use `URL.createObjectURL()` with a Blob for the download — this is safe and standard.
- Set correct MIME type on the Blob (`image/png` or `image/jpeg`).

---

## 3. OWASP Top 10 Relevance Check

| OWASP Category | Relevant? | Notes |
|---|---|---|
| A01: Broken Access Control | Low | No server/auth — N/A for pure client-side. Revisit if backend added. |
| A02: Cryptographic Failures | Low | No sensitive data storage. No encryption needed for client-only tool. |
| A03: Injection | **Medium** | XSS via text input if rendered in DOM. Canvas rendering is safe. |
| A04: Insecure Design | **Medium** | Must design file upload validation and canvas security correctly from the start. |
| A05: Security Misconfiguration | Low | Minimal config for static SPA. CSP headers matter if hosted. |
| A06: Vulnerable Components | **Medium** | Dependency supply chain risk. Minimize deps. |
| A07: Auth Failures | N/A | No authentication in scope. |
| A08: Data Integrity Failures | Low | No server-side processing. |
| A09: Logging & Monitoring | Low | Client-side only — limited logging needed. |
| A10: SSRF | N/A | No server to exploit. |

---

## 4. Recommendations Summary

### Critical Priority

| # | Issue | Severity | Likelihood | Recommendation |
|---|-------|----------|------------|----------------|
| 1 | SVG upload with embedded JS | Critical | Medium | Block SVG uploads entirely. Whitelist JPEG, PNG, WebP, GIF only. |
| 2 | Image file validation bypass | High | Medium | Validate file signature (magic bytes), not just extension/MIME. |

### High Priority

| # | Issue | Severity | Likelihood | Recommendation |
|---|-------|----------|------------|----------------|
| 3 | XSS via text input in DOM | High | Medium | Use `textContent` not `innerHTML`. Sanitize if using contenteditable. |
| 4 | Memory DoS via large images | High | Medium | Enforce max file size (10MB) and max dimensions (4096x4096). |
| 5 | Dependency supply chain | High | Low | Minimize deps. Pin versions. Run `npm audit`. |

### Medium Priority

| # | Issue | Severity | Likelihood | Recommendation |
|---|-------|----------|------------|----------------|
| 6 | Canvas tainting from cross-origin images | Medium | Low | Use only local file uploads via FileReader/createObjectURL. |
| 7 | EXIF metadata privacy | Medium | Medium | Canvas re-encoding strips EXIF — verify and document this. |
| 8 | CSP headers on hosting | Medium | Low | Deploy with strict Content-Security-Policy headers. |

### Low Priority

| # | Issue | Severity | Likelihood | Recommendation |
|---|-------|----------|------------|----------------|
| 9 | Download filename injection | Low | Low | Sanitize or hardcode export filenames. |
| 10 | Clipboard API abuse | Low | Low | If implementing paste-to-upload, validate pasted content type. |

---

## 5. Implementation Notes for Developers

### File Upload Validation Pattern

```
Recommended validation chain:
1. Check file.type against whitelist: ['image/jpeg', 'image/png', 'image/webp', 'image/gif']
2. Check file extension against whitelist: ['.jpg', '.jpeg', '.png', '.webp', '.gif']
3. Check file size: max 10MB (10 * 1024 * 1024 bytes)
4. Read first 4-8 bytes and validate magic number signatures:
   - JPEG: FF D8 FF
   - PNG: 89 50 4E 47
   - WebP: 52 49 46 46 ... 57 45 42 50
   - GIF: 47 49 46 38
5. Create an Image() object and validate it loads successfully
6. Check naturalWidth/naturalHeight <= 4096
```

### Secure Text Rendering Pattern

```
SAFE:
- canvas.getContext('2d').fillText(userText, x, y)  // Canvas pixel rendering
- element.textContent = userText                     // DOM text node (no parsing)

UNSAFE:
- element.innerHTML = userText      // Parses HTML — XSS risk
- document.write(userText)          // Never use
- eval(userText)                    // Never use
- template literals in DOM without escaping
```

### Content Security Policy (for hosting)

```
Recommended CSP header:
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

Note: `img-src blob: data:` is required for canvas export and local file display. `unsafe-inline` for styles may be needed for dynamic text styling — consider using nonces if possible.

### Secure Export Pattern

```
Recommended approach:
1. canvas.toBlob(callback, 'image/png') — generates clean PNG, strips metadata
2. URL.createObjectURL(blob) — creates safe local blob URL
3. Programmatic <a> click with download attribute and sanitized filename
4. URL.revokeObjectURL() after download to free memory
```

---

## 6. Security Testing Suggestions

### Manual Tests

1. **Upload SVG with script tag** — Verify it's rejected, not rendered.
2. **Upload .html file renamed to .png** — Verify magic byte check catches it.
3. **Upload 100MB image** — Verify file size limit triggers before processing.
4. **Upload 20000x20000px image** — Verify dimension limit prevents canvas memory bomb.
5. **Enter `<script>alert(1)</script>` as text overlay** — Verify it renders as literal text, not executed.
6. **Enter `<img src=x onerror=alert(1)>` as text** — Verify no DOM injection.
7. **Check exported image for EXIF data** — Use exiftool to verify metadata is stripped.
8. **Inspect download filename** — Verify no path traversal characters accepted.

### Automated Tests

- Unit test file validation function with malicious file fixtures.
- Unit test text rendering ensures no DOM injection.
- CSP violation reporting (if hosted) to catch policy bypasses.
- `npm audit` in CI pipeline.
- Lighthouse security audit on deployed page.

---

## 7. Architecture Security Considerations

### Client-Only Architecture (Recommended for v1)

```
[Browser] --> [Local Files Only]
    |
    v
[Canvas API] --> [Blob Export] --> [Download]
```

**Security advantage:** No server means no server-side vulnerabilities, no data storage, no authentication complexity, no API abuse. All processing is local — user data never leaves their browser.

### If Backend is Added Later (v2+)

If sharing/saving features are added, the threat model changes significantly:

- **New risks:** Stored XSS, unauthorized access, data breach, SSRF, API abuse, DoS.
- **Required additions:** Authentication, authorization, input validation server-side, rate limiting, image re-encoding on server, content moderation, abuse reporting.
- **Recommendation:** Conduct a fresh security review before adding any server-side features.

---

## 8. Delegation

```
DELEGATE:
- architect: Ensure client-only architecture is maintained for v1. Design file upload validation as a dedicated module. Plan CSP headers for deployment.
- senior_dev: Implement file validation chain (type + extension + magic bytes + dimensions). Use textContent exclusively for DOM text. Implement secure export with blob URLs and proper cleanup.
```

---

**PHASE_COMPLETE: research**
