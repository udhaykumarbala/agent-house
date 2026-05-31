# UX Research: Image Resizer — New Feature Additions (v2)

## Current State Audit

The app already implements all P0 and P1 features from the original product spec:
- Drag-and-drop upload (JPG, PNG, WebP, up to 50MB)
- Resize by pixel dimensions with aspect ratio lock
- Live preview with file size comparison
- Format conversion (JPG/PNG/WebP) with quality slider
- Presets (Instagram Post, Instagram Story, HD 1080p, Twitter/X Header, OG Image)
- Scale buttons (25%, 50%, 75%)
- One-click download with success feedback
- Friendly brutalist UI, privacy messaging, mobile responsive

This research focuses on **new features** that would elevate the tool beyond a simple resizer into a more complete image utility — while staying true to the "one-page, privacy-first, brutalist" identity.

---

## Pattern Analysis — Feature Expansion

### Squoosh (squoosh.app)
- **Pattern Used**: Before/after comparison slider overlaying original and processed image. Users drag a divider to see original vs. compressed side by side on the same canvas.
- **Why It Works**: Gives instant visual confidence that quality is acceptable. Eliminates guesswork around lossy compression. Users feel in control.
- **What We Can Learn**: A comparison view is the #1 missing feature for any image tool that changes quality. Even a simple side-by-side (not a slider) would add significant value.

### Canva (canva.com) — Image Resize Tool
- **Pattern Used**: Platform-specific preset categories (Social Media, Print, Web) with sub-presets showing platform icons. "Custom size" input. One-click resize to multiple sizes simultaneously.
- **Why It Works**: Users don't think in pixels — they think in platforms ("I need this for LinkedIn"). Categorized presets with recognizable icons reduce cognitive load.
- **What We Can Learn**: Our preset list is flat. Grouping presets by category (Social, Web, Print) and adding more platform sizes (LinkedIn, YouTube, Pinterest, Facebook) would serve more use cases. Platform icons/logos next to preset names improve scannability.

### Remove.bg / Cleanup.pictures
- **Pattern Used**: Background removal as a one-click action. The result is shown instantly on a checkerboard pattern (transparency indicator). Download as PNG with transparency.
- **Why It Works**: Background removal is one of the most common companion tasks to resizing — users often resize AND need a transparent background.
- **What We Can Learn**: We already show a checkerboard in our preview card. Adding background removal would be a high-value feature, BUT it requires ML models (WASM or server-side), which conflicts with our "no dependencies" philosophy. **Recommendation**: Skip this for now but keep the checkerboard pattern for PNG transparency.

### TinyPNG (tinypng.com)
- **Pattern Used**: Batch upload — drag multiple files at once. Each file shows as a row with its own progress bar, original size, compressed size, and individual download button. "Download All" button generates a ZIP.
- **Why It Works**: Power users (social media managers, web developers) rarely resize just one image. Batch processing is the #1 requested feature in image tools.
- **What We Can Learn**: Batch resize with consistent settings is a major upgrade. The UX pattern is clear: upload multiple → apply same dimensions/format/quality → download as ZIP or individually.

### Figma (figma.com) — Export Panel
- **Pattern Used**: Multiple export formats and scales from a single asset. Users can export at 1x, 2x, 3x simultaneously. Format selection per export variant. Suffix-based naming (_@2x, _large).
- **Why It Works**: Web developers and designers frequently need multiple sizes of the same image (thumbnail, medium, large; or 1x, 2x, 3x for retina).
- **What We Can Learn**: Multi-size export from a single source image would be unique among simple resizers. Pattern: "Export as 1x + 2x" or "Export thumbnail + full" with customizable suffixes.

### Birme (birme.net)
- **Pattern Used**: Visual crop with drag handles before resize. Users position a crop frame on the image, then it resizes the cropped area to target dimensions. Supports both "fit" and "fill" resize modes.
- **Why It Works**: Cropping before resizing is a natural workflow — users often need to reframe before hitting a specific aspect ratio.
- **What We Can Learn**: A simple crop tool (not a full editor) that lets users select a region before resizing fills a major workflow gap. The "fit vs. fill" distinction is also valuable.

### Photon (photon.sh) / CSS Filters in Browser
- **Pattern Used**: Quick image adjustments — brightness, contrast, saturation, sharpness — as slider controls alongside resize. These use CSS filters or Canvas API filters, requiring no external libraries.
- **Why It Works**: Users doing a quick resize often also want a quick tweak — bump contrast, sharpen a thumbnail, desaturate for a design mockup.
- **What We Can Learn**: Basic adjustments via Canvas API (brightness, contrast, saturation) are lightweight to implement and complement resizing. Sharpening is especially useful after downscaling (images often look soft after resize).

---

## User Journey Analysis — Extended Workflows

### Entry Point (unchanged)
User has an image that needs resizing. They open the app.

### Extended Core Loops (NEW)

**Loop A: Single Image, Quick Resize (existing)**
1. Drop image → set dimensions → download

**Loop B: Single Image, Optimize + Resize (NEW)**
1. Drop image → adjust brightness/contrast → crop to focus area → set dimensions → download
2. Motivation: "I need a hero image for my blog post — crop out the background clutter, make it brighter, resize to 1200x630"

**Loop C: Batch Resize (NEW)**
1. Drop multiple images → set dimensions/format once → download all as ZIP
2. Motivation: "I have 20 product photos that all need to be 800x800 for my Shopify store"

**Loop D: Multi-Format Export (NEW)**
1. Drop image → export at multiple sizes simultaneously (thumbnail, medium, large)
2. Motivation: "I need this image as a 150x150 avatar, 600x600 card, and 1200x1200 full size"

### Success States (NEW)
- **Batch**: "All 15 images resized and downloaded as a ZIP in under 10 seconds"
- **Crop + Resize**: "Cropped out the messy background and got exactly the dimensions I needed"
- **Quick Adjust**: "Sharpened the thumbnail so it doesn't look blurry after downscaling"

---

## Mental Model — Feature Expectations

Users who return to an image resizer expect it to grow into a **lightweight image utility**, not a full editor. The mental model is:

> "This tool handles the quick stuff I'd normally open Photoshop for but shouldn't have to."

**Acceptable feature additions** (still feels like a utility tool):
- Crop, rotate, flip
- Brightness/contrast/saturation adjustments
- Batch processing
- More presets
- Image comparison (before/after)
- Copy to clipboard
- Drag to reorder/manage multiple images

**Feature creep danger zone** (starts feeling like an editor):
- Text overlay
- Filters/effects (sepia, blur, vignette)
- Layers
- Drawing tools
- Templates

---

## Anti-Patterns to Avoid

- **Don't**: Add features that require server-side processing — **Why**: Privacy-first is the core differentiator. Background removal (ML models) and advanced compression (MozJPEG, AVIF) would require WASM at minimum, breaking the zero-dependency philosophy.
- **Don't**: Make batch mode the default flow — **Why**: Single-image resize is still the 80% use case. Batch should be an opt-in mode, not clutter the main UI. TinyPNG gets this right; iLoveIMG doesn't.
- **Don't**: Add crop with complex controls (free-form polygon, magnetic lasso) — **Why**: Simple rectangular crop with aspect ratio constraints is sufficient. Anything more belongs in an editor.
- **Don't**: Auto-apply adjustments (auto-brightness, auto-contrast) — **Why**: Users want control. Auto-adjustments feel unpredictable and can ruin images. Always manual, always preview-first.
- **Don't**: Hide new features behind tabs or navigation — **Why**: The single-page identity is core to the UX. New features should fold into the existing controls panel, not create new pages or modals.
- **Don't**: Add a watermark or branding to exported images — **Why**: Utility tools that stamp output lose trust immediately.

---

## Recommended New Features (Prioritized)

Based on research, here are the features ranked by user value, implementation feasibility, and alignment with the app's identity:

### Tier 1 — High Value, Moderate Effort

#### 1. Image Cropping (Before Resize)
- **Pattern**: Visual crop overlay on the preview image with drag handles. Aspect ratio constraint options (free, original, 1:1, 16:9, 4:3). Crop is applied before resize.
- **Why**: Cropping + resizing is the most common compound workflow. Without it, users need a separate tool for reframing.
- **Reference**: Birme, iLoveIMG crop tool
- **UX**: Add a "Crop" toggle/button in the controls panel. When active, overlay handles appear on the preview. Confirm with a checkmark button. Non-destructive until confirmed.

#### 2. Batch Resize
- **Pattern**: Allow dropping multiple images. Show a file list with thumbnails, names, sizes. Apply current resize settings to all. Download individually or as ZIP.
- **Why**: Power users need this. Social media managers, e-commerce sellers, web developers all resize in batches.
- **Reference**: TinyPNG's batch list, iLoveIMG
- **UX**: The drop zone should accept multiple files. Editor view shows a scrollable file list on the left. Settings apply globally. "Download All (ZIP)" button replaces single download.

#### 3. Expanded Presets with Categories
- **Pattern**: Group presets into collapsible categories: Social Media, Web/SEO, Print, Video Thumbnails. Add missing platforms.
- **New presets to add**:
  - **Social**: Facebook Cover (820x312), LinkedIn Post (1200x627), LinkedIn Banner (1584x396), Pinterest Pin (1000x1500), YouTube Thumbnail (1280x720), TikTok (1080x1920)
  - **Web**: Favicon (32x32, 192x192), Email Header (600x200), App Store Icon (1024x1024)
  - **Print**: A4 300dpi (2480x3508), Letter 300dpi (2550x3300), 4x6 Photo (1200x1800)
- **Why**: More presets = fewer Google searches for "what size is a LinkedIn post?"
- **Reference**: Canva's platform-specific preset categories

### Tier 2 — Medium Value, Low Effort

#### 4. Image Rotation & Flip
- **Pattern**: Rotate 90 CW / CCW buttons, Horizontal flip, Vertical flip. Applied via Canvas transform before resize.
- **Why**: Photos from phones are frequently sideways. Users shouldn't need another tool just to rotate.
- **Reference**: Nearly every image tool has this as a basic control
- **UX**: Small icon buttons in a row: [Rotate Left] [Rotate Right] [Flip H] [Flip V]. Place above or near the crop control.

#### 5. Copy to Clipboard
- **Pattern**: "Copy" button alongside "Download". Uses the Clipboard API (`navigator.clipboard.write()`) to copy the resized image as PNG.
- **Why**: Many users paste images directly into Slack, email composers, Google Docs, or design tools. Downloading then re-uploading is friction.
- **Reference**: Screenshot tools (Cleanshot, ShareX) always offer clipboard as output
- **UX**: Secondary button next to Download: "Copy to Clipboard". Show "Copied!" feedback for 2s.

#### 6. Before/After Comparison
- **Pattern**: Toggle or slider that shows original vs. resized image side by side or overlapping. Shows dimension and file size comparison.
- **Why**: Users want visual confirmation that quality hasn't degraded, especially at lower quality settings or significant size reductions.
- **Reference**: Squoosh's comparison slider
- **UX**: Simple toggle: "Compare" button that splits the preview into original (left) / resized (right) with labels.

#### 7. Sharpen After Resize
- **Pattern**: A subtle sharpening pass (unsharp mask via Canvas convolution) applied automatically or via toggle after downscaling. Small slider for intensity (0-100).
- **Why**: Downscaled images always lose sharpness. A post-resize sharpen is the single most impactful image adjustment.
- **Reference**: Photoshop's "Sharpen for Web" workflow, Squoosh's output processing
- **UX**: Toggle switch: "Sharpen" with a small intensity slider. Default: off. Placed near quality slider.

### Tier 3 — Nice-to-Have, Varies in Effort

#### 8. Basic Adjustments (Brightness, Contrast, Saturation)
- **Pattern**: Three sliders: Brightness (-100 to +100), Contrast (-100 to +100), Saturation (-100 to +100). Applied via Canvas API `ctx.filter` or manual pixel manipulation.
- **Why**: Quick tweaks eliminate the need to open another tool. Especially useful for thumbnails that need to "pop."
- **Reference**: Photon, iOS Photos quick-edit panel
- **UX**: Collapsible "Adjustments" section in controls panel, below quality slider.

#### 9. Drag-and-Drop Reorder for Batch
- **Pattern**: In batch mode, drag to reorder files in the list. Useful when downloading as ZIP with numbered filenames.
- **Why**: Minor but appreciated for batch workflows where order matters.

#### 10. PWA / Offline Support
- **Pattern**: Service worker for offline caching. "Install" prompt on supported browsers. App icon on home screen.
- **Why**: Makes the tool available without internet. Reinforces the "nothing leaves your browser" messaging.
- **Reference**: Squoosh is a PWA
- **UX**: Small "Install" button in header (shown only when supported). No other UI change.

#### 11. Keyboard Shortcuts
- **Pattern**: `Cmd/Ctrl+S` to download, `Cmd/Ctrl+Z` to undo crop, `Tab` through controls, `Escape` to reset/go back.
- **Why**: Power user efficiency. Small effort, big quality-of-life improvement.
- **UX**: No visible UI needed. Optional "?" shortcut help overlay.

#### 12. Dark Mode
- **Pattern**: Toggle between light (current cream) and dark (near-black background with inverted accents) themes. Respects `prefers-color-scheme` by default.
- **Why**: Many users work at night or prefer dark interfaces. The brutalist style adapts well to dark mode.
- **Reference**: Squoosh has dark mode, most modern tools do
- **UX**: Toggle button in header. CSS custom properties already exist in variables.css — retheme via a `.dark` class on body.

#### 13. EXIF / Metadata Viewer
- **Pattern**: Expandable panel showing image metadata: camera, date, GPS (with option to strip), color profile, DPI.
- **Why**: Users uploading photos sometimes want to check metadata. The app already strips EXIF on resize (Canvas API doesn't preserve it), but surfacing what was stripped adds transparency.
- **UX**: Small "Info" icon/button near the file info line. Opens a collapsible panel.

#### 14. Undo/Redo History
- **Pattern**: Track state changes (dimension edits, crop, adjustments) and allow stepping backward/forward.
- **Why**: Users make mistakes. Without undo, they must manually revert changes or reset entirely.
- **UX**: Undo/Redo buttons in header. `Cmd/Ctrl+Z` / `Cmd/Ctrl+Shift+Z`.

---

## Feature Interaction Map

How new features relate to each other and the existing flow:

```
[Upload] ──► [Editor View]
              │
              ├── [Crop Tool] ──► applies before resize
              ├── [Rotate/Flip] ──► applies before resize
              ├── [Dimensions / Presets / Scale] (existing)
              ├── [Format / Quality] (existing)
              ├── [Adjustments: Brightness/Contrast/Saturation] ──► applies during resize
              ├── [Sharpen] ──► applies after resize
              ├── [Compare Toggle] ──► shows original vs result
              │
              └── [Download Bar]
                   ├── [Download] (existing)
                   ├── [Copy to Clipboard] (new)
                   └── [Download All ZIP] (batch mode)
```

```
[Batch Mode]
Upload multiple ──► [File List + Thumbnails]
                     ├── Select one to preview/crop individually
                     ├── Global settings apply to all
                     └── [Download All ZIP] / [Download Individual]
```

---

## Controls Panel Layout Recommendation

To accommodate new features without cluttering, use collapsible sections:

```
┌─────────────────────────────┐
│ TRANSFORM                   │  (new section)
│ [Crop] [↺ Left] [↻ Right]  │
│ [↔ Flip H] [↕ Flip V]      │
├─────────────────────────────┤
│ DIMENSIONS (existing)       │
│ Width [____] × Height [____]│
│ 🔗 Ratio locked             │
├─────────────────────────────┤
│ PRESETS (expanded)          │
│ ▸ Social Media              │
│   [Instagram Post] [Story]  │
│   [Facebook Cover] [TikTok] │
│   [YouTube] [LinkedIn]      │
│ ▸ Web & SEO                 │
│   [OG Image] [Favicon]      │
│ ▸ Scale                     │
│   [25%] [50%] [75%]         │
├─────────────────────────────┤
│ FORMAT & QUALITY (existing) │
│ Format: [JPG ▾]             │
│ Quality: ━━━━━━━━● 85%      │
├─────────────────────────────┤
│ ▸ ADJUSTMENTS (collapsed)   │  (new section)
│   Brightness: ━━━━●━━ +10   │
│   Contrast:   ━━━━━●━ +20   │
│   Saturation: ━━●━━━━ -15   │
│   Sharpen:    ━━━●━━━ 30%   │
└─────────────────────────────┘
```

---

## Implementation Priority Matrix

| Feature | User Value | Effort | Risk | Priority |
|---------|-----------|--------|------|----------|
| Expanded Presets | High | Low | Low | P0 — Do first |
| Rotate & Flip | High | Low | Low | P0 — Do first |
| Copy to Clipboard | Medium | Low | Low | P0 — Do first |
| Dark Mode | Medium | Low | Low | P1 — Do soon |
| Image Cropping | High | Medium | Medium | P1 — Do soon |
| Before/After Compare | Medium | Low | Low | P1 — Do soon |
| Keyboard Shortcuts | Medium | Low | Low | P1 — Do soon |
| Sharpen After Resize | Medium | Medium | Low | P1 — Do soon |
| Batch Resize | High | High | Medium | P2 — Plan carefully |
| Basic Adjustments | Medium | Medium | Low | P2 — Plan carefully |
| PWA / Offline | Medium | Medium | Low | P2 — Plan carefully |
| EXIF Metadata Viewer | Low | Medium | Low | P3 — If time allows |
| Undo/Redo History | Medium | High | Medium | P3 — If time allows |

---

## Accessibility Considerations for New Features

- **Crop tool**: Must be operable via keyboard (arrow keys to move handles, Enter to confirm, Escape to cancel)
- **Rotate/Flip**: Icon buttons must have `aria-label` text ("Rotate 90 degrees clockwise")
- **Batch file list**: Use `role="list"` with keyboard navigation and screen reader announcements for add/remove
- **Collapsible sections**: Use `aria-expanded` and `aria-controls` on section headers
- **Dark mode**: Maintain 4.5:1 contrast ratios in dark theme. Test all states.
- **Keyboard shortcuts**: Show in a discoverable help panel. Don't override browser defaults (Cmd+P, Cmd+L, etc.)
- **Copy to clipboard**: Announce "Image copied to clipboard" via `aria-live` region

---

## Summary: Recommended Feature Roadmap

**Phase 1 (Quick Wins)**: Expanded presets, rotate/flip, copy to clipboard, dark mode
**Phase 2 (Core Additions)**: Crop tool, before/after comparison, keyboard shortcuts, sharpen
**Phase 3 (Power Features)**: Batch resize, basic adjustments, PWA
**Phase 4 (Polish)**: EXIF viewer, undo/redo, multi-size export

Each phase preserves the single-page, privacy-first, friendly brutalist identity while meaningfully expanding utility.

---
Status: READY_FOR_PLANNING
