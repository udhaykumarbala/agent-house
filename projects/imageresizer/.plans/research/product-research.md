# Product Research: Image Resizer — Feature Expansion (Phase 2)

## Current State Assessment

The MVP is fully built with all P0 and P1 features shipped:
- Drag-and-drop upload (JPG, PNG, WebP, up to 50MB)
- Resize by pixel dimensions with aspect ratio lock
- Live preview with progress indicator
- One-click download with sanitized filenames
- Output format conversion (JPEG, PNG, WebP)
- Quality slider (1-100)
- File size comparison (original vs resized with %)
- Preset dimensions (Instagram Post, Instagram Story, HD 1080p, Twitter/X Header, OG Image)
- Scale buttons (25%, 50%, 75%)
- New Image / Reset
- Friendly brutalist UI, mobile responsive, privacy messaging
- Multi-step downscaling for quality
- EXIF stripping, CSP-safe (no innerHTML)

**What's missing**: The app handles single-image resize well but lacks the "next tier" features that would make users choose it over Squoosh or return regularly.

---

## Competitor Analysis (Feature Gaps)

| Competitor | Features We Don't Have | What They Do Poorly |
|------------|----------------------|---------------------|
| **Squoosh (squoosh.app)** | Side-by-side before/after slider, AVIF output, resize method options (Lanczos, Mitchell), offline PWA, CLI tool | No presets, overwhelming for casual users, no batch |
| **iLoveIMG (iloveimg.com)** | Batch resize, crop, rotate/flip, compress-only mode, watermark tool | Uploads to server, ad-heavy, slow, privacy concern |
| **TinyPNG (tinypng.com)** | Batch compress (up to 20), API for developers, smart compression with quality preservation | No resize at all — compress only. Server-side. Limited free usage (500/mo) |
| **Birme (birme.net)** | Batch resize, crop with focal point, border/padding, rename pattern, ZIP download | Dated UI, no format conversion, no quality control |
| **Photopea (photopea.com)** | Full image editing (crop, rotate, filters, layers, text), batch actions, PSD/AI support | Massive overkill for resizing, steep learning curve, ad-supported |

### Key Takeaways
1. **Crop + Rotate** is the most common "next feature" across competitors — users frequently need both in the same workflow
2. **Batch processing** is the #1 power-user request but adds significant complexity
3. **Image compression** as a standalone mode (keep dimensions, reduce file size) is a huge use case we're ignoring
4. **Before/after comparison** is Squoosh's signature UX — we could do a simpler version
5. **No competitor combines privacy + presets + crop in a delightful UI** — that's our lane

---

## Target User Feedback Patterns

Based on common user behavior with image resize tools:

### "Quick-Fix Quinn" (existing persona) — What they'd want next:
- **Crop before resize** — "I need to crop this product photo to square, then resize to 1080x1080"
- **Rotate/flip** — "The photo is sideways from my phone camera"
- **Compress-only mode** — "I don't need to change the size, just make the file smaller for email"
- **More presets** — LinkedIn, YouTube thumbnail, Etsy listing, Shopify product, email header

### New persona: "Batch Barbara"
- **Demographics**: E-commerce seller, photographer, marketing team member
- **Need**: "I have 15 product photos that all need to be 800x800 for Shopify"
- **Current Solution**: Uses Birme or a desktop app, wishes there was something nicer
- **Key Feature**: Upload multiple images, apply same resize settings, download as ZIP

### New persona: "Developer Dave"
- **Demographics**: Front-end dev, web performance optimizer
- **Need**: "I need to convert this PNG to WebP and compress it under 100KB for my website"
- **Current Solution**: Squoosh or CLI tools (sharp, imagemagick)
- **Key Feature**: Precise file size targeting, WebP optimization, keyboard shortcuts

---

## Feature Opportunity Research

### 1. Image Cropping (High Value)
**Why**: Cropping is the most natural companion to resizing. Users frequently need to change aspect ratio before hitting a target dimension.
- **Implementation**: Canvas-based crop with draggable handles
- **Complexity**: Medium — need a crop overlay UI, touch support, and integration with the resize pipeline
- **Competitors**: iLoveIMG has it, Squoosh doesn't, Birme has it with focal points
- **Our angle**: Simple, fast crop — not a full editor. Preset crop ratios (1:1, 4:3, 16:9, free)

### 2. Rotate & Flip (High Value, Low Effort)
**Why**: Phone photos are often rotated wrong. Flip is useful for mirroring. These are trivial to implement with Canvas.
- **Implementation**: Canvas `rotate()` and `scale(-1,1)` / `scale(1,-1)`
- **Complexity**: Low — 4 buttons (rotate CW, rotate CCW, flip horizontal, flip vertical)
- **Competitors**: Nearly all have this; we're notably missing it

### 3. Batch Resize (High Value, High Effort)
**Why**: Power users and e-commerce sellers need to resize multiple images with the same settings.
- **Implementation**: Multi-file upload, process sequentially via Canvas, ZIP download using JSZip
- **Complexity**: High — UI for multi-file management, progress tracking, ZIP generation
- **Dependency**: Would need JSZip library (lightweight, ~10KB gzipped) or use Compression Streams API
- **Our angle**: Batch with the same privacy promise — all client-side, no uploads

### 4. Compress-Only Mode (High Value, Low Effort)
**Why**: Many users want to reduce file size without changing dimensions. Currently they have to manually enter the same dimensions.
- **Implementation**: A "Compress Only" toggle or mode that skips resize, just re-encodes at the selected quality
- **Complexity**: Low — the Resizer engine already supports this (just pass original dimensions)
- **Competitors**: TinyPNG is a whole product just for this; iLoveIMG separates it as a different tool

### 5. Before/After Comparison (Medium Value, Medium Effort)
**Why**: Squoosh's signature feature. Helps users see quality impact before downloading.
- **Implementation**: Side-by-side or slider overlay comparing original vs resized
- **Complexity**: Medium — need a comparison UI component
- **Our angle**: Simple side-by-side with file size comparison (we already show sizes)

### 6. More Presets (Medium Value, Low Effort)
**Why**: Users are always looking up dimensions. More presets = less friction.
- **New presets to add**:
  - LinkedIn Post (1200x627)
  - YouTube Thumbnail (1280x720)
  - Pinterest Pin (1000x1500)
  - Facebook Cover (820x312)
  - Shopify Product (2048x2048)
  - Email Header (600x200)
  - Passport Photo (413x531)
- **Implementation**: Just more buttons in the existing preset grid
- **Complexity**: Very low — HTML + existing JS handles it
- **UX consideration**: Group presets by category (Social, Web, E-commerce, Print) to avoid overwhelming

### 7. Custom Preset Saving (Medium Value, Medium Effort)
**Why**: Repeat users often resize to the same custom dimensions. Let them save their own presets.
- **Implementation**: localStorage-based preset storage
- **Complexity**: Medium — need add/edit/delete UI for custom presets
- **Privacy-aligned**: No server needed, data stays in browser

### 8. Keyboard Shortcuts (Low-Medium Value, Low Effort)
**Why**: Power users expect keyboard efficiency. Tab navigation already works, but shortcuts would be nice.
- **Shortcuts**: `Ctrl/Cmd+S` to download, `Ctrl/Cmd+Z` to reset, `Enter` to trigger resize
- **Complexity**: Low — global keydown listener
- **Accessibility bonus**: Improves keyboard-only usability

### 9. PWA / Offline Support (Medium Value, Medium Effort)
**Why**: Makes the tool installable and available without internet. Fits the "bookmarkable utility" strategy.
- **Implementation**: Service worker + manifest.json
- **Complexity**: Medium — need to cache assets properly
- **Distribution win**: Users can add to home screen on mobile

### 10. Drag-to-Reorder / Multi-Output (Low Value, High Effort)
**Why**: Download the same image in multiple sizes at once (e.g., generate all social media sizes from one upload)
- **Complexity**: High — significant UI and processing work
- **Defer**: Better suited for a v3.0

---

## Prioritized Feature Recommendations

### Tier 1 — Build Next (High impact, achievable scope)

| # | Feature | Value | Effort | Why Now |
|---|---------|-------|--------|---------|
| 1 | **Rotate & Flip** | High | Low | Missing basic expected feature; 4 buttons, trivial Canvas work |
| 2 | **Image Cropping** | High | Medium | Most requested companion to resize; preset ratios (1:1, 4:3, 16:9, free) |
| 3 | **Compress-Only Mode** | High | Low | Huge use case we're ignoring; almost no new code needed |
| 4 | **More Presets (categorized)** | Medium | Very Low | Quick win; add LinkedIn, YouTube, Pinterest, Facebook, Shopify, Email |

### Tier 2 — Build After (Strong value, more effort)

| # | Feature | Value | Effort | Why |
|---|---------|-------|--------|-----|
| 5 | **Batch Resize** | High | High | Power-user magnet; needs multi-file UI + ZIP generation |
| 6 | **Before/After Comparison** | Medium | Medium | UX polish; builds confidence for quality-sensitive users |
| 7 | **PWA / Offline** | Medium | Medium | Distribution win; fits "bookmarkable utility" strategy |
| 8 | **Custom Preset Saving** | Medium | Medium | Retention feature for repeat users; localStorage-based |

### Tier 3 — Consider Later

| # | Feature | Value | Effort | Why Wait |
|---|---------|-------|--------|----------|
| 9 | **Keyboard Shortcuts** | Low-Med | Low | Nice-to-have, not blocking any workflow |
| 10 | **Multi-Output Export** | Low | High | Complex; batch covers most of this use case |
| 11 | **AVIF Support** | Low-Med | High | Requires WASM; browser support still growing |

---

## Market Positioning Update

**Updated One-Liner**: "Crop, resize, and compress images instantly in your browser — no uploads, no ads, no nonsense."

**Competitive Moat After Expansion**:
- Only privacy-first tool with crop + resize + compress in one flow
- Only tool with social media presets AND a beautiful UI
- Still zero dependencies (except JSZip for future batch feature)
- Still 100% client-side

**Price Point**: Remains **Free** — utility tool, no monetization needed

---

## Technical Considerations

| Feature | Technical Notes |
|---------|----------------|
| Crop | Canvas-based; need draggable crop overlay with touch support. Apply crop before resize in pipeline. |
| Rotate/Flip | Canvas `rotate()` + `translate()`. Update state dimensions on 90-degree rotations (swap W/H). |
| Compress-Only | No new engine code. Just a UI toggle that sets target dimensions = original dimensions. |
| More Presets | HTML-only change. Consider a category accordion or tab UI to manage the growing list. |
| Batch | Multi-file input, sequential Canvas processing, ZIP via JSZip or Compression Streams API. |
| PWA | manifest.json + service worker. Cache index.html, CSS, JS, fonts. |

---

## Risk Assessment

| Risk | Impact | Mitigation |
|------|--------|------------|
| Crop UI adds complexity, confuses simple users | Medium | Make crop optional — a button to enter "crop mode" before resize |
| Too many presets overwhelm the UI | Low | Group by category with collapsible sections |
| Batch processing crashes on mobile | Medium | Limit batch to 10 images on mobile, 20 on desktop; show memory warnings |
| Feature creep turns us into Photopea | High | Strict scope: crop, resize, compress, format convert. Nothing else. |
| JSZip adds external dependency | Low | Only load JSZip when batch mode is used (lazy load); it's 10KB gzipped |

---
Status: READY_FOR_PLANNING
