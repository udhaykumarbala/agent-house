# Security Research: New Feature Opportunities & Threat Analysis

## Role: Security Expert
## Date: 2026-03-16
## Scope: Research new features that can be added, evaluated through a security lens

---

## 1. Current State Assessment

### What's Already Built
The image resizer is a fully functional MVP with:
- Drag-and-drop upload (JPG, PNG, WebP)
- Resize by pixel dimensions with aspect ratio lock
- Preset sizes (Instagram, HD, Twitter/X, OG Image)
- Scale buttons (25%, 50%, 75%)
- Format conversion (JPEG, PNG, WebP)
- Quality slider for lossy formats
- Live preview with debounced updates
- File size comparison (original vs resized)
- One-click download with sanitized filename
- 100% client-side, privacy-first architecture

### Security Posture of Current Build
| Control | Status | Notes |
|---------|--------|-------|
| DOM XSS prevention (textContent only) | Implemented | `app.js:4` documents the rule; all rendering uses `textContent` |
| File type validation (MIME + size) | Implemented | `utils.js:123-143` — validates type allowlist + 50MB limit |
| File name sanitization | Implemented | `utils.js:26-33` — strips traversal, control chars, reserved chars |
| Blob URL lifecycle management | Implemented | `app.js:384-385, 474-476` — revokes on new preview and reset |
| No innerHTML usage | Verified | Zero `innerHTML` in the entire codebase |
| No inline scripts/handlers | Verified | All JS in external files, events via `addEventListener` |
| CSP-ready architecture | Verified | No inline JS, structure supports strict CSP |
| Magic bytes validation | **NOT implemented** | Only MIME type checked, not file header bytes |
| `beforeunload` blob cleanup | **NOT implemented** | Blob URLs not revoked on tab close |
| CSP `<meta>` tag in HTML | **NOT implemented** | No CSP meta tag in `index.html` |

### Security Gaps to Address Alongside New Features
1. **Magic bytes validation** — Current validation only checks `file.type` (MIME). A renamed `.html` file could pass validation. Add header byte verification.
2. **`beforeunload` cleanup** — Blob URLs survive tab close. Add `window.addEventListener('beforeunload', ...)` to revoke.
3. **CSP meta tag** — Add `<meta http-equiv="Content-Security-Policy">` to `index.html` for defense-in-depth even before server deployment.

---

## 2. New Feature Opportunities — Security-Prioritized Analysis

I've evaluated potential new features across two axes: **user value** and **security risk**. Features are grouped by risk tier.

### Tier 1: Low Risk, High Value (Recommend First)

These features stay within the existing security model (client-side only, no new attack surface).

#### F1: Image Cropping (Before Resize)
**What**: Allow users to select a rectangular crop area before resizing. Drag handles on the preview image to define the crop region.

**User Value**: High — users frequently need to crop before resizing (e.g., remove whitespace, center a subject). Currently requires a separate tool.

**Security Analysis**:
| Threat | Risk | Notes |
|--------|------|-------|
| XSS via crop coordinates | None | Coordinates are numeric, no user strings rendered |
| Canvas manipulation abuse | None | Same Canvas API already in use |
| Memory usage on large crops | Low | Canvas operations already handle this |

**Implementation Security Notes**:
- Crop coordinates must be validated as positive integers within original image bounds
- Clamp values: `x >= 0, y >= 0, w >= 1, h >= 1, x+w <= naturalWidth, y+h <= naturalHeight`
- No new user-string rendering needed — all values are computed
- Use `ctx.drawImage(img, sx, sy, sw, sh, 0, 0, dw, dh)` — same API, different args

**Verdict**: Safe to build. No new attack surface.

---

#### F2: Image Rotation (90° increments + flip)
**What**: Rotate image 90° CW/CCW and flip horizontal/vertical. Button controls, not free-form rotation.

**User Value**: Medium-High — common quick-fix need, especially for mobile photos with wrong orientation.

**Security Analysis**:
| Threat | Risk | Notes |
|--------|------|-------|
| Canvas dimension swap bugs | Low | Must correctly swap width/height on 90° rotation |
| Memory spike on rotation | Low | Creates new canvas of same pixel area |

**Implementation Security Notes**:
- Rotation is pure Canvas transformation (`ctx.rotate()`, `ctx.translate()`)
- No user input to sanitize — buttons trigger fixed transformations
- Validate that rotated dimensions don't exceed browser canvas limits

**Verdict**: Safe to build. Minimal complexity.

---

#### F3: Batch/Multi-Image Resize
**What**: Upload multiple images, apply same resize settings, download all as individual files or a ZIP archive.

**User Value**: Very High — the #1 feature gap vs competitors like iLoveIMG. Content creators resize dozens of images at once.

**Security Analysis**:
| Threat | Risk | Notes |
|--------|------|-------|
| Memory exhaustion (DoS) | **Medium** | Processing many large images simultaneously can crash the tab |
| ZIP bomb (output) | None | We create the ZIP, not consume it |
| File name collision in ZIP | Low | Multiple files with same name |
| Increased blob URL leaks | Medium | More URLs to track and revoke |

**Implementation Security Notes**:
- **Queue processing**: Process images sequentially, not in parallel, to limit memory pressure
- **Cap batch size**: Max 20 images per batch (configurable). Reject with friendly message.
- **Per-file validation**: Every file in the batch must pass the same validation (MIME, size, magic bytes)
- **Memory management**: Process one image at a time, revoke blob URLs after each is added to ZIP
- **ZIP library**: Use [JSZip](https://stuk.github.io/jszip/) (MIT, well-maintained, no eval/innerHTML). Include via local copy, NOT CDN. Add SRI hash if CDN is ever considered.
- **File name deduplication**: If multiple files produce the same output name, append counter (e.g., `photo-resized-1.jpg`, `photo-resized-2.jpg`)
- **Total size limit**: Enforce combined input size limit (e.g., 200MB total) to prevent memory crashes
- **Progress indication**: Show per-file progress to prevent user from thinking it's frozen

**Dependency Risk — JSZip**:
- JSZip is the standard for client-side ZIP creation
- 17K+ GitHub stars, actively maintained, no known CVEs
- Must be self-hosted (copy into `scripts/vendor/jszip.min.js`) not loaded from CDN
- Verify integrity hash on download
- Alternative: Use `CompressionStream` API (native, no dep) but browser support is limited

**Verdict**: Safe to build with proper memory management. First feature requiring an external dependency — handle with care.

---

#### F4: Custom Preset Management
**What**: Let users save their own preset sizes (e.g., "My Blog Header: 800x400") that persist across sessions.

**User Value**: Medium — power users who resize to the same dimensions repeatedly.

**Security Analysis**:
| Threat | Risk | Notes |
|--------|------|-------|
| XSS via preset names | **Medium** | User-defined strings stored and rendered |
| localStorage injection | Low | Malicious page on same origin could tamper with stored presets |
| Prototype pollution via JSON parse | Low | If preset data structure is manipulated |

**Implementation Security Notes**:
- **Render preset names via `textContent` only** — never innerHTML. This is the critical rule.
- **Validate preset data on read**: When loading from localStorage, validate structure:
  ```
  - name: string, max 50 chars, stripped of control characters
  - width: positive integer, 1-99999
  - height: positive integer, 1-99999
  ```
- **Use `JSON.parse()` with try/catch** — malformed data must not crash the app
- **Sanitize preset names on save**: Strip `<>`, control characters, limit length
- **Limit preset count**: Max 20 custom presets to prevent localStorage abuse
- **No eval() or Function()** on stored data — ever
- **Consider prefixing localStorage keys**: `imgresizer_presets` to reduce collision with other apps on same origin

**Verdict**: Safe to build with proper sanitization. Moderate care needed around user-defined strings.

---

#### F5: Drag-to-Reorder / Compare Preview
**What**: Side-by-side or slider comparison between original and resized image (like Squoosh's comparison slider).

**User Value**: Medium — helps users evaluate quality loss before downloading.

**Security Analysis**:
| Threat | Risk | Notes |
|--------|------|-------|
| Additional blob URL | Low | One more blob URL to manage |
| No user input rendering | None | Pure visual comparison |

**Implementation Security Notes**:
- Reuses existing preview image and original image — no new file handling
- Slider position is a numeric value — no sanitization needed
- Ensure original image blob URL is also tracked for revocation

**Verdict**: Safe to build. Trivial security surface.

---

#### F6: Copy to Clipboard
**What**: Button to copy the resized image directly to clipboard (for pasting into design tools, chat apps, etc.).

**User Value**: Medium — faster workflow than download-then-upload for many use cases.

**Security Analysis**:
| Threat | Risk | Notes |
|--------|------|-------|
| Clipboard API permissions | Low | Requires user gesture (click), browser handles permission |
| Clipboard data exposure | None | User explicitly triggers the copy |

**Implementation Security Notes**:
- Use `navigator.clipboard.write()` with `ClipboardItem` API
- Requires a user gesture (button click) — cannot be automated/abused
- Feature-detect before offering: `if (navigator.clipboard && navigator.clipboard.write)`
- Fallback: Don't show the button if API not available
- **CSP note**: No additional CSP changes needed

**Verdict**: Safe to build. Browser-managed permissions.

---

#### F7: Dark Mode / Theme Toggle
**What**: Toggle between light (current butter/cream) and dark mode. Persisted via localStorage.

**User Value**: Medium — expected feature for modern web tools. Reduces eye strain.

**Security Analysis**:
| Threat | Risk | Notes |
|--------|------|-------|
| localStorage preference tampering | Negligible | Theme is cosmetic only |
| CSS injection via theme value | None if using class toggle | Only toggle a body class |

**Implementation Security Notes**:
- Store preference as simple string in localStorage: `"light"` or `"dark"`
- On load: read value, validate it's one of exactly `"light"` or `"dark"`, apply class
- Toggle via `document.body.classList.toggle('theme-dark')`
- NO user-provided CSS values — all theme variations defined in `variables.css`

**Verdict**: Safe to build. Zero attack surface.

---

#### F8: EXIF Data Viewer (Read-Only)
**What**: Display extracted EXIF metadata (camera, ISO, aperture, GPS coordinates) from the original image before resize. Informational only.

**User Value**: Medium — photographers appreciate seeing metadata. Also educates users about what data their images carry (privacy awareness).

**Security Analysis**:
| Threat | Risk | Notes |
|--------|------|-------|
| **XSS via EXIF strings** | **HIGH** | EXIF fields like `ImageDescription`, `Artist`, `Copyright` contain arbitrary strings that could include `<script>` tags or HTML |
| EXIF parsing library vulnerabilities | Medium | Third-party library needed |
| Privacy concern — displaying GPS | Medium | Showing location data could surprise users |

**Implementation Security Notes**:
- **CRITICAL**: Every EXIF value rendered to the DOM MUST use `textContent`, never `innerHTML`
- **Library choice**: Use [exifr](https://github.com/nickytonline/exifr) (lightweight, modern) or [exif-js](https://github.com/nickytonline/exif-js). Self-host, not CDN.
- **Whitelist displayed fields**: Only show known safe fields (dimensions, camera model, ISO, aperture, focal length, date). Do NOT display arbitrary/custom EXIF tags.
- **GPS display**: If showing GPS, present as "Location data detected" warning, not as rendered coordinates by default. Give user option to view.
- **Sanitize all string values**: Even with textContent, strip control characters from displayed strings
- **Memory**: Parse EXIF lazily (on-demand when user clicks "View metadata"), not on every upload

**Verdict**: Buildable but requires careful string handling. Medium risk due to untrusted EXIF strings.

---

### Tier 2: Medium Risk, High Value (Recommend with Caution)

These features expand the attack surface in controlled ways.

#### F9: Watermark Overlay
**What**: Add text or image watermark to the resized image. Text watermark with customizable position, font size, opacity.

**User Value**: High — photographers and content creators frequently need watermarks. Currently requires Photoshop or separate tools.

**Security Analysis**:
| Threat | Risk | Notes |
|--------|------|-------|
| **XSS via watermark text** | **Medium** | User types arbitrary text that could be rendered unsafely |
| Canvas text rendering injection | None | `ctx.fillText()` is safe — it renders as pixels, not HTML |
| Watermark image as attack vector | Low | Second image upload needs same validation |

**Implementation Security Notes**:
- **Text watermark**: Canvas `ctx.fillText()` is inherently safe — text is rasterized to pixels. No XSS risk on the canvas itself.
- **Text preview in UI**: If showing watermark text in any DOM element, use `textContent` only
- **Watermark image**: Must pass same file validation (MIME, size, magic bytes) as primary image
- **Font selection**: Offer only predefined fonts — do NOT allow loading arbitrary font files
- **Opacity/position**: Numeric values, standard validation

**Verdict**: Safe to build. Canvas text rendering is inherently XSS-proof. Just guard DOM display of user text.

---

#### F10: Image Filters / Adjustments
**What**: Basic image adjustments — brightness, contrast, saturation, grayscale, sepia. Slider controls.

**User Value**: Medium — nice-to-have enhancement. Keeps users from needing a separate tool for basic adjustments.

**Security Analysis**:
| Threat | Risk | Notes |
|--------|------|-------|
| Canvas pixel manipulation | None | `getImageData()`/`putImageData()` operates on raw pixel arrays |
| Performance / DoS on large images | Medium | Pixel-level operations on 30MP images are expensive |
| Web Worker requirement | Low | May need to move to Worker to avoid blocking UI |

**Implementation Security Notes**:
- All operations are numeric transformations on pixel arrays — no string handling
- **Performance guard**: For images over 10MP, process in a Web Worker to avoid freezing the UI
- **Web Worker security**: Workers run in isolated context. Use `new Worker('scripts/filter-worker.js')` — must be same-origin
- **CSP impact**: `worker-src 'self'` must be added to CSP if using Web Workers
- **No `eval()` in workers** — pass filter parameters as structured data via `postMessage()`

**Verdict**: Safe to build. Computationally intensive but no security surface.

---

#### F11: PWA / Offline Support
**What**: Add a Service Worker to cache the app for offline use. Installable as a standalone app.

**User Value**: High — users can install it and use it without internet. Perfect for a privacy-first tool.

**Security Analysis**:
| Threat | Risk | Notes |
|--------|------|-------|
| **Service Worker cache poisoning** | **Medium** | If a compromised version is cached, it persists |
| Stale cache serving vulnerable code | Medium | Users may run old, unpatched versions |
| SW scope and registration | Low | Must be registered at correct scope |

**Implementation Security Notes**:
- **Cache versioning**: Use versioned cache names (`imgresizer-v1`, `imgresizer-v2`). On update, delete old caches in the `activate` event.
- **Cache-first with network fallback**: For assets. But include a version check mechanism.
- **SW scope**: Register at root (`/`) to cover all app paths
- **Update detection**: On `navigator.serviceWorker.controllerchange`, prompt user to refresh
- **CSP note**: Service Workers respect the page's CSP. No changes needed.
- **manifest.json**: Include with minimum permissions. No `background_sync` or `push` unless needed.
- **SW file location**: Must be at root or scope root — `/sw.js`

**Verdict**: Safe to build. Standard PWA patterns. Main concern is cache staleness.

---

#### F12: Undo/Redo History
**What**: Allow users to undo/redo resize, crop, rotation, and filter changes. Stack-based history.

**User Value**: Medium — quality-of-life feature. Reduces friction when experimenting with settings.

**Security Analysis**:
| Threat | Risk | Notes |
|--------|------|-------|
| Memory exhaustion from history stack | **Medium** | Storing multiple canvas states in memory |
| Blob URL accumulation | Medium | Each history state may have a blob URL |

**Implementation Security Notes**:
- **Limit history depth**: Max 10-20 states. Discard oldest when limit reached.
- **Store settings, not images**: Instead of caching canvas outputs, store the *settings* (width, height, crop, rotation, filters) and re-render from original. This is much more memory-efficient and avoids blob URL accumulation.
- **Original image is immutable**: Never modify the original image reference. All operations derive from it.

**Verdict**: Safe if settings-based (not image-based) history. Manageable memory.

---

### Tier 3: Higher Risk (Requires Careful Design)

#### F13: Share / Direct Link
**What**: Generate a shareable link that encodes resize settings (e.g., `?w=1080&h=1080&format=webp&q=85`). Not the image itself — just the tool with pre-filled settings.

**User Value**: Low-Medium — niche use case (teams sharing resize templates).

**Security Analysis**:
| Threat | Risk | Notes |
|--------|------|-------|
| **DOM-based XSS via URL parameters** | **HIGH** | Any URL parameter rendered in the DOM is an XSS vector |
| Open redirect | None | No navigation based on URL params |
| Parameter tampering | Low | Only affects user's own session |

**Implementation Security Notes**:
- **NEVER render URL parameters via innerHTML** — always `textContent`
- **Parse parameters defensively**: Use `URLSearchParams` with strict validation
  ```
  const params = new URLSearchParams(window.location.search);
  const width = parseInt(params.get('w'), 10);
  if (isNaN(width) || width < 1 || width > 99999) { ignore; }
  ```
- **Whitelist allowed parameters**: Only `w`, `h`, `format`, `q`. Ignore everything else.
- **Validate format**: Must be exactly `"jpeg"`, `"png"`, or `"webp"`
- **Do NOT reflect raw parameter values in the DOM** — always validate and cast to expected type first
- **No `eval()` or `Function()` on URL data** — ever
- **Test**: Navigate to `?w=<script>alert(1)</script>` — must not execute

**Verdict**: Buildable with strict input validation. HIGH risk if implemented carelessly. Every param must be validated before use.

---

#### F14: Paste from Clipboard
**What**: Support pasting an image directly from clipboard (Ctrl+V / Cmd+V). Common workflow when screenshotting.

**User Value**: High — very convenient for screenshots. Squoosh supports this.

**Security Analysis**:
| Threat | Risk | Notes |
|--------|------|-------|
| Clipboard may contain non-image data | Low | Must validate pasted data type |
| HTML content in clipboard | **Medium** | Clipboard may contain HTML blobs alongside images |
| Crafted image in clipboard | Low | Same risk as file upload |

**Implementation Security Notes**:
- **Listen to `paste` event on document**:
  ```
  document.addEventListener('paste', (e) => {
    const items = e.clipboardData.items;
    for (const item of items) {
      if (item.type.startsWith('image/')) {
        const file = item.getAsFile();
        // Same validation as file upload
      }
    }
  });
  ```
- **Only process `image/*` MIME types** — ignore `text/html`, `text/plain`, etc.
- **Run the same `validateFile()` on pasted image** as on dropped/selected files
- **Do NOT read `text/html` clipboard items** — potential XSS if rendered
- **CSP note**: No changes needed. Clipboard read is event-driven, not API-initiated.

**Verdict**: Safe to build with proper MIME filtering. Must ignore non-image clipboard content.

---

#### F15: AI-Powered Smart Crop / Background Removal
**What**: Use a WASM-based ML model (e.g., MediaPipe, ONNX Runtime) for intelligent cropping or background removal.

**User Value**: Very High — killer feature that would differentiate from all competitors.

**Security Analysis**:
| Threat | Risk | Notes |
|--------|------|-------|
| **WASM supply chain** | **HIGH** | Large binary dependency that's hard to audit |
| Model file integrity | High | ML models downloaded must be verified |
| Memory consumption | High | ML models + image processing = significant memory |
| Execution time DoS | Medium | ML inference on large images can freeze the tab |

**Implementation Security Notes**:
- **Self-host WASM files** — never load from third-party CDN
- **Verify integrity**: SRI hashes for all WASM binaries
- **CSP update**: Add `'wasm-unsafe-eval'` to `script-src` (required for WASM execution)
- **Web Worker isolation**: Run inference in a Worker to avoid UI blocking
- **Memory limits**: Set explicit limits on input image size for AI features (e.g., max 5MP for background removal)
- **Dependency audit**: WASM binaries must be from reputable sources (Google MediaPipe, ONNX Runtime). Pin versions.
- **User consent**: Clearly communicate that processing is still local — "AI runs in your browser, not on a server"

**Verdict**: High-value but high-complexity. Requires significant dependency management. Recommend as Phase 3+.

---

## 3. Feature Priority Matrix (Security Expert Recommendation)

| Priority | Feature | User Value | Security Risk | Effort | Recommendation |
|----------|---------|------------|---------------|--------|----------------|
| 1 | **F14: Paste from Clipboard** | High | Low | Small | Build now — easy win, minimal risk |
| 2 | **F6: Copy to Clipboard** | Medium | None | Small | Build now — pairs with paste |
| 3 | **F1: Image Cropping** | High | None | Medium | Build now — no new attack surface |
| 4 | **F2: Rotation/Flip** | Medium-High | None | Small | Build now — trivial security surface |
| 5 | **F7: Dark Mode** | Medium | None | Small | Build now — zero risk |
| 6 | **F3: Batch Resize** | Very High | Medium | Large | Build next — needs JSZip dependency review |
| 7 | **F9: Watermark** | High | Low | Medium | Build next — canvas-safe |
| 8 | **F5: Before/After Compare** | Medium | None | Medium | Build next — pure visual |
| 9 | **F11: PWA / Offline** | High | Low | Medium | Build next — standard patterns |
| 10 | **F4: Custom Presets** | Medium | Low | Small | Build when needed — localStorage concerns manageable |
| 11 | **F10: Image Filters** | Medium | None | Large | Build later — computationally intensive |
| 12 | **F12: Undo/Redo** | Medium | Low | Medium | Build later — memory management needed |
| 13 | **F8: EXIF Viewer** | Medium | Medium | Medium | Build carefully — untrusted string rendering |
| 14 | **F13: Share Link** | Low-Medium | High | Small | Build carefully — URL parameter XSS risk |
| 15 | **F15: AI Smart Crop** | Very High | High | Very Large | Build later (Phase 3+) — complex dependency chain |

---

## 4. Security Hardening Recommendations (Before New Features)

Before adding features, close the existing security gaps:

### 4.1 Add Magic Bytes Validation
Currently `utils.js:123-143` only checks `file.type`. Add file header verification:

```
JPEG: first 3 bytes = FF D8 FF
PNG:  first 8 bytes = 89 50 4E 47 0D 0A 1A 0A
WebP: bytes 0-3 = 52 49 46 46 AND bytes 8-11 = 57 45 42 50
```

Read file header with `FileReader.readAsArrayBuffer()` on the first 12 bytes, then verify.

### 4.2 Add `beforeunload` Cleanup
```js
window.addEventListener('beforeunload', function() {
  if (state.resizedBlobUrl) URL.revokeObjectURL(state.resizedBlobUrl);
});
```

### 4.3 Add CSP Meta Tag
Add to `<head>` in `index.html`:
```html
<meta http-equiv="Content-Security-Policy" content="default-src 'self'; script-src 'self'; style-src 'self' https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com; img-src 'self' blob: data:; connect-src 'none'; object-src 'none'; frame-src 'none'; base-uri 'self'; form-action 'none'">
```

### 4.4 Add Max Canvas Dimension Check
Before creating canvases for resize, check against browser limits:
```js
const MAX_CANVAS_DIMENSION = 16384; // Conservative cross-browser limit
if (width > MAX_CANVAS_DIMENSION || height > MAX_CANVAS_DIMENSION) {
  // Reject or cap dimensions with user warning
}
```

---

## 5. Dependency Risk Assessment for Potential New Dependencies

| Dependency | Feature | GitHub Stars | Last Updated | Known CVEs | Self-Host? | Risk |
|-----------|---------|-------------|-------------|------------|------------|------|
| JSZip | Batch resize (ZIP) | 9.8K+ | Active | None current | Yes (required) | Low |
| exifr | EXIF viewer | 900+ | Active | None current | Yes (required) | Low |
| exif-js | EXIF viewer (alt) | 5K+ | Stale (2019) | None known | Avoid — unmaintained | Medium |
| FileSaver.js | Better downloads | 21K+ | Active | None current | Optional — native approach preferred | Low |
| MediaPipe WASM | AI smart crop | Google-maintained | Active | N/A | Yes (required) | Medium |

**Policy**: Any new dependency must be:
1. Self-hosted (copied to `scripts/vendor/`)
2. Verified via SHA-256 hash
3. Licensed as MIT or equivalent
4. Actively maintained (commit in last 12 months)
5. Reviewed for `eval()`, `innerHTML`, `document.write()` usage before inclusion

---

## 6. Feature-Specific CSP Impact

| Feature | CSP Change Needed? | Details |
|---------|-------------------|---------|
| Crop, Rotate, Flip | No | Same Canvas API |
| Clipboard (Copy/Paste) | No | Browser-managed |
| Dark Mode | No | CSS class toggle |
| Batch Resize (JSZip) | No | JSZip uses Blob APIs, no eval |
| Watermark | No | Canvas text/image rendering |
| Custom Presets | No | localStorage is CSP-exempt |
| Image Filters | Maybe | Add `worker-src 'self'` if using Web Workers |
| PWA/Offline | Maybe | Add `worker-src 'self'` for Service Worker |
| EXIF Viewer | No | Library reads ArrayBuffer, outputs data objects |
| Share Link | No | URL parsing is native |
| AI Features | **Yes** | Add `'wasm-unsafe-eval'` to `script-src` + `worker-src 'self'` |

---

## 7. Competitor Feature Comparison (Security Perspective)

| Feature | Squoosh | iLoveIMG | TinyPNG | Photopea | Our App (Current) | Opportunity |
|---------|---------|----------|---------|----------|-------------------|-------------|
| Crop | Yes | Yes | No | Yes | **No** | High priority |
| Rotate/Flip | Yes | Yes | No | Yes | **No** | High priority |
| Batch | No | Yes (server) | Yes (server) | No | **No** | High value, ours would be client-side (privacy win) |
| Clipboard paste | Yes | No | No | Yes | **No** | Easy win |
| Clipboard copy | No | No | No | Yes | **No** | Easy win |
| Watermark | No | Yes (server) | No | Yes | **No** | Differentiator if client-side |
| Dark mode | No | No | No | Yes | **No** | Expected feature |
| Filters | Yes (advanced) | No | No | Yes (full) | **No** | Medium priority |
| Offline/PWA | Yes | No | No | No | **No** | Matches privacy narrative |
| EXIF viewer | No | No | No | Yes | **No** | Nice-to-have |
| Background removal | No | Yes (server) | No | Yes | **No** | Killer feature if client-side |
| Compare slider | Yes | No | No | No | **No** | Medium priority |

**Key insight**: The biggest competitive gaps are **crop, rotate, batch, and clipboard support**. All four are low-to-medium security risk and align with the client-side privacy model.

---

## 8. Recommended Implementation Phases

### Phase 2A: Quick Wins (1-2 days each)
1. Paste from Clipboard (F14)
2. Copy to Clipboard (F6)
3. Image Rotation/Flip (F2)
4. Dark Mode (F7)
5. Security hardening (4.1-4.4)

### Phase 2B: Core Enhancements (3-5 days each)
6. Image Cropping (F1)
7. Before/After Comparison Slider (F5)
8. Custom Presets with localStorage (F4)
9. PWA / Offline Support (F11)

### Phase 3: Power Features (1-2 weeks each)
10. Batch Resize with ZIP download (F3) — first external dependency
11. Watermark Overlay (F9)
12. Basic Image Filters (F10)
13. Undo/Redo History (F12)

### Phase 4: Advanced (Research Required)
14. AI Smart Crop / Background Removal (F15)
15. Share Link with URL Parameters (F13) — high XSS risk, needs careful design

---

## 9. Security Testing Plan for New Features

| Feature | Test Case | Expected Result |
|---------|-----------|-----------------|
| Crop | Set crop coordinates to negative values | Clamped to 0 |
| Crop | Set crop region larger than image | Clamped to image bounds |
| Rotate | Rotate image exceeding max canvas dimensions | Warning shown, operation blocked |
| Clipboard paste | Paste HTML content | Ignored — only image types processed |
| Clipboard paste | Paste crafted image with XSS filename | Filename rendered as text |
| Custom presets | Name a preset `<script>alert(1)</script>` | Rendered as plain text |
| Custom presets | Tamper localStorage with malformed JSON | App loads with default presets, no crash |
| Dark mode | Set localStorage theme to `<script>` | Value rejected, default theme applied |
| Batch resize | Upload 50 files simultaneously | Rejected with "max 20 files" message |
| Batch resize | Upload 20x 50MB files (1GB total) | Rejected with total size limit message |
| Watermark text | Enter `<img src=x onerror=alert(1)>` as watermark | Canvas renders it as literal text pixels |
| Share link | `?w=<script>alert(1)</script>` | Parameter ignored, no script execution |
| EXIF viewer | Image with `<script>` in EXIF description | Rendered as plain text via textContent |

---

## 10. Delegation

```
DELEGATE:
- architect: Review dependency policy for JSZip (batch) and exifr (EXIF). Design Web Worker architecture for filters. Plan PWA service worker caching strategy. Decide on CSP meta tag vs deployment headers.
- senior_dev: Implement magic bytes validation, beforeunload cleanup, CSP meta tag, canvas dimension limits. Build clipboard paste/copy, crop, rotate. Ensure all new features follow textContent-only DOM rule.
- junior_dev: Can implement dark mode toggle, rotation buttons (with review). Should NOT implement URL parameter parsing (share link) or EXIF rendering without senior review.
```

---

**Summary**: The image resizer has a strong security foundation. Most high-value new features (crop, rotate, clipboard, dark mode, batch) are low-risk and stay within the existing client-side security model. The main risks come from: (1) any feature that renders user-provided strings to the DOM, (2) new external dependencies, and (3) URL parameter parsing. Prioritize the quick wins first — they add significant user value with negligible security cost.

---
Status: PHASE_COMPLETE: research
