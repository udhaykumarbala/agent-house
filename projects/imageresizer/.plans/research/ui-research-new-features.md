# UI Research: New Features for Image Resizer

## Context

The current Image Resizer is a fully functional MVP with: drag-and-drop upload, pixel resize, aspect ratio lock, live preview, format conversion (JPG/PNG/WebP), quality slider, dimension presets, scale buttons, file size comparison, and download. All 100% client-side with a Friendly Brutalist (Butter & Charcoal) design system.

This research explores which additional features would deliver the most value, how competitors handle them, and how they should look within our existing design language.

---

## Feature Research: What Modern Image Tools Offer Beyond Basic Resize

### 1. Image Cropping

**Competitor Analysis:**
- **Squoosh**: No crop. Resize only.
- **iLoveIMG**: Separate crop tool (different page). Clunky multi-step.
- **Canva**: Excellent inline crop with handles and preset aspect ratios.
- **Photopea**: Full Photoshop-style crop with rotation. Overkill for our use case.
- **iOS Photos**: Simple drag-to-crop with aspect ratio presets. Gold standard for quick crop.

**Why It Matters**: Users often need to crop AND resize in one step. A social media manager uploading a landscape photo to Instagram (square) needs to crop to 1:1 first, then resize to 1080x1080. Without crop, they need a second tool. This is the #1 most requested "missing feature" for simple resize tools.

**Design Direction for Our App**:
- Keep it minimal: a crop overlay with draggable handles on the preview image
- Aspect ratio presets for crop (Free, 1:1, 4:3, 16:9, 9:16)
- Crop happens BEFORE resize in the pipeline
- Should feel like an optional step, not a mandatory workflow change
- Brutalist treatment: thick border handles, no soft overlay dimming — instead use a solid-color mask outside the crop area

**Inspiration**: iOS Photos crop interface — simple, intuitive, handles + rotation dial

---

### 2. Image Rotation & Flip

**Competitor Analysis:**
- **Squoosh**: No rotation.
- **iLoveIMG**: Separate rotation tool page.
- **Canva**: Inline rotation with degree handle + 90-degree snap buttons.
- **TinyPNG**: No rotation.
- **macOS Preview**: Quick 90-degree rotation button in toolbar.

**Why It Matters**: Photos from phones are often sideways or upside down (EXIF orientation gets stripped by canvas). Quick 90-degree rotation and horizontal/vertical flip are near-zero-complexity features that save users from opening another tool.

**Design Direction for Our App**:
- Four small icon buttons above or below the preview: Rotate Left (90 CCW), Rotate Right (90 CW), Flip Horizontal, Flip Vertical
- These are instant operations on the canvas — no UI complexity needed
- Brutalist treatment: small square icon buttons (44x44px) with thick 2px border and 2px offset shadow
- Group them in a "Transform" toolbar strip

**Inspiration**: macOS Preview's toolbar — just simple icon buttons for rotate/flip

---

### 3. Paste from Clipboard

**Competitor Analysis:**
- **Squoosh**: Supports paste! Ctrl/Cmd+V to paste an image from clipboard.
- **iLoveIMG**: No paste support.
- **TinyPNG**: No paste support.
- **Most image tools**: Missing this entirely.

**Why It Matters**: Screenshots are the #2 most common image source (after photo files). Users screenshot something, then need to resize it. Currently they'd need to: save screenshot as file, then drag to our tool. With paste support: Cmd+V and done. Massive friction reduction.

**Design Direction for Our App**:
- Add "or paste from clipboard (Ctrl+V)" text below the dropzone subtext
- Listen for `paste` event on `document`
- Extract image data from `clipboardData.items`
- Process it exactly like a dropped file
- No visible UI change needed — just the hint text and event listener
- On paste success, transition to editor view like normal

**Inspiration**: Squoosh's paste behavior — invisible until you need it, then delightful

---

### 4. Copy to Clipboard (Output)

**Competitor Analysis:**
- **Squoosh**: No copy-to-clipboard for output.
- **Figma**: Excellent — right-click > Copy as PNG. Instant.
- **macOS Screenshot**: Automatically to clipboard.

**Why It Matters**: Sometimes users don't want to download a file — they want to paste the resized image directly into Slack, email, or a CMS. Copy-to-clipboard skips the file-download-then-upload dance.

**Design Direction for Our App**:
- Add a secondary "Copy" button next to the Download button in the download bar
- Uses `navigator.clipboard.write()` with `ClipboardItem`
- Ghost button style (outlined, no fill) to not compete with Download
- Brief "Copied!" success state (butter yellow flash, 1.5s)
- Feature-detect clipboard API — hide button if not supported

**Inspiration**: Figma's "Copy as PNG" — one click, in your clipboard

---

### 5. Before/After Comparison Slider

**Competitor Analysis:**
- **Squoosh**: THE gold standard — draggable vertical divider between original and compressed
- **TinyPNG**: Shows before/after sizes but no visual comparison
- **Photopea**: No comparison mode

**Why It Matters**: When changing quality/format, users want to SEE the difference. "Is 60% quality too low? Does WebP look different from JPG?" A comparison slider answers this instantly and builds confidence in the output quality.

**Design Direction for Our App**:
- Draggable vertical divider on the preview image
- Left side: original image, Right side: resized/compressed output
- Labels: "Original" / "Resized" with dimensions below each
- Toggle button to enter/exit comparison mode (not always visible)
- Brutalist treatment: thick 3px divider line in charcoal, chunky drag handle in coral
- Only relevant when format/quality changes, not pure dimension changes

**Inspiration**: Squoosh's comparison slider — clean, intuitive, informative

---

### 6. Batch Resize (Multiple Images)

**Competitor Analysis:**
- **iLoveIMG**: Core strength — batch upload, apply same resize to all, download ZIP
- **Bulk Resize Photos**: Dedicated batch tool — simple and fast
- **Squoosh**: Single image only
- **Birme**: Batch resize with crop, downloads ZIP

**Why It Matters**: Social media managers and e-commerce sellers frequently need to resize 5-20 images to the same dimensions. Doing them one at a time is painful. This is the power-user feature that turns a "nice tool" into a "daily driver."

**Design Direction for Our App**:
- Toggle between "Single" and "Batch" mode on the upload screen
- Batch mode: multi-file drop zone, shows thumbnail grid of queued images
- Apply same settings (dimensions, format, quality) to all images
- Progress bar per image + overall progress
- Download all as ZIP (using JSZip library — ~100KB, client-side)
- Brutalist treatment: thumbnail grid with thick borders, checkmark overlay on completed items
- This is a significant feature — should be a clearly separate mode, not bolted onto current UI

**Complexity Note**: This is the most complex feature. Requires a thumbnail grid, progress tracking, ZIP generation. Should be Phase 2 or later.

**Inspiration**: Birme's batch interface — grid of thumbnails with shared controls

---

### 7. Dark Mode

**Competitor Analysis:**
- **Squoosh**: Auto dark mode (follows system preference)
- **Photopea**: Dark by default (editor convention)
- **TinyPNG**: Light only
- **iLoveIMG**: Light only

**Why It Matters**: Many developers and designers work in dark mode. A bright butter-cream background at 11pm is jarring. Dark mode is increasingly expected in web tools. It also provides a natural way to show the app has "depth" in its design system.

**Design Direction for Our App**:
- Charcoal background (#1A1A1A) with warm dark surface (#2A2A2A)
- Coral primary (#E8625C) stays — it pops beautifully on dark backgrounds
- Butter yellow (#F5D547) secondary becomes more vivid on dark
- Borders become lighter: #4A4A4A or keep #1A1A1A with lighter surface contrast
- Shadows: subtle warm glow instead of hard offset? Or keep offset but in a lighter color
- Toggle: sun/moon icon button in header — chunky brutalist toggle
- Respect `prefers-color-scheme: dark` media query for auto-detection

**Design Challenge**: The warm butter-cream IS our brand identity. Dark mode needs to feel equally "warm" — not cold/corporate dark. Consider warm charcoal (#2C2420) instead of pure #1A1A1A for backgrounds.

**Inspiration**: Poolsuite's dark mode — warm, not sterile

---

### 8. PWA / Installable & Offline Support

**Competitor Analysis:**
- **Squoosh**: Full PWA — installable, works offline. Best-in-class.
- **TinyPNG**: Not a PWA (server-dependent anyway)
- **iLoveIMG**: Not a PWA

**Why It Matters**: Our tool is 100% client-side — it SHOULD work offline. A service worker + manifest would let users "install" it as a desktop/mobile app and use it without internet. This reinforces our privacy message ("it literally works without a server") and improves load times for repeat users.

**Design Direction for Our App**:
- Add `manifest.json` with app name, icons, theme color (butter cream)
- Service worker caches all assets (HTML, CSS, JS, fonts)
- "Install App" prompt in header (shown once, dismissible)
- Brutalist treatment for install prompt: bold card with thick border, "Install for offline use" CTA
- App icon: our existing favicon at multiple resolutions

**Complexity Note**: Low — mostly configuration. No UI changes needed beyond an optional install prompt.

---

### 9. Keyboard Shortcuts

**Competitor Analysis:**
- **Photopea**: Full keyboard shortcut set (Photoshop-compatible)
- **Squoosh**: Minimal keyboard support
- **Figma**: Extensive shortcut system with visual cheatsheet

**Why It Matters**: Power users (designers, developers) expect keyboard shortcuts. Quick actions like Cmd+S to download, Cmd+Z to undo dimensions, Cmd+V to paste, Tab to navigate — these make the tool feel professional.

**Design Direction for Our App**:
- Essential shortcuts only (not a full shortcut system):
  - `Cmd/Ctrl + V`: Paste image from clipboard
  - `Cmd/Ctrl + S`: Download (prevent default browser save)
  - `Cmd/Ctrl + Shift + C`: Copy to clipboard
  - `Cmd/Ctrl + N`: New image (reset)
  - `L`: Toggle aspect lock
  - `?`: Show shortcut cheatsheet
- Shortcut hints shown as small badges on buttons (e.g., `Ctrl+S` badge on Download)
- Cheatsheet: overlay modal in brutalist style (thick border card, list of shortcuts)
- Brutalist treatment for kbd badges: monospace text in small bordered boxes like `[Ctrl+S]`

**Inspiration**: Figma's keyboard shortcut overlay — clean, comprehensive, dismissible

---

### 10. EXIF Data Viewer (Before Strip)

**Competitor Analysis:**
- **Jeffrey's EXIF Viewer**: Dedicated tool, shows all metadata
- **Squoosh**: No EXIF display
- **Most resize tools**: Silently strip EXIF without telling the user what was there

**Why It Matters**: Photographers and privacy-conscious users want to know what metadata their image contains before sharing. Our app already strips EXIF via canvas — we should SHOW what we're stripping. This reinforces our privacy-first positioning and adds a "wow, I didn't know my photo had GPS data" moment.

**Design Direction for Our App**:
- Small expandable section below the preview: "Metadata stripped" with expand arrow
- Shows: Camera model, date taken, GPS (if present, highlighted in warning amber), dimensions, color space
- GPS data gets a special warning: "Location data found and will be removed"
- Brutalist treatment: collapsible card with thick border, metadata in monospace (JetBrains Mono), GPS warning in amber
- This is read-only information — no editing

**Inspiration**: iOS photo "Info" panel — clean summary of metadata

---

### 11. Custom Presets (Save Your Own)

**Competitor Analysis:**
- **Bulk Resize Photos**: Lets you save custom preset profiles
- **Most tools**: Only offer fixed presets

**Why It Matters**: Users who resize for the same platform repeatedly (e.g., "my Shopify product images are always 800x800") want to save their custom dimensions as a reusable preset. This converts one-time users into power users.

**Design Direction for Our App**:
- "+" button at the end of the preset grid
- Opens a small form: Name, Width, Height
- Saves to `localStorage` (no server needed)
- Custom presets appear alongside built-in presets with a different visual treatment (e.g., outlined vs filled background)
- Delete button (X) on custom presets only
- Brutalist treatment: same preset pill style but with a dashed border to distinguish user-created presets
- Limit: 5-10 custom presets to prevent UI overflow

**Inspiration**: Browser bookmark bar — simple save/recall pattern

---

### 12. Drag-to-Resize Handles on Preview

**Competitor Analysis:**
- **Figma**: Resize handles on selected elements
- **Canva**: Corner drag handles for images
- **Most image resizers**: Input-only, no visual resize

**Why It Matters**: Some users think visually, not numerically. Being able to grab a corner of the preview and drag to resize is more intuitive than typing numbers. This makes the tool feel more interactive and "real."

**Design Direction for Our App**:
- Four corner handles on the preview image (visible in editor view)
- Dragging a handle updates width/height inputs in real-time
- Shift+drag to maintain aspect ratio (or auto-maintain if lock is on)
- Brutalist treatment: chunky 12x12px square handles with thick borders, coral fill
- Corner handles only (not edge handles) to keep it simple
- Cursor changes to resize cursor on hover

**Complexity Note**: Medium — requires mouse/touch event handling, coordinate math, and real-time input sync. Can be laggy on mobile.

---

## Feature Priority Matrix

| Feature | User Value | Complexity | Fits Brand | Recommended Priority |
|---------|-----------|------------|------------|---------------------|
| Image Cropping | Very High | Medium | Yes | P1 - Should Have |
| Rotate/Flip | High | Low | Yes | P1 - Should Have |
| Paste from Clipboard | High | Low | Yes | P1 - Should Have |
| Copy to Clipboard | Medium | Low | Yes | P1 - Should Have |
| Before/After Slider | Medium | Medium | Yes | P1 - Should Have |
| Keyboard Shortcuts | Medium | Low | Yes | P1 - Should Have |
| Dark Mode | Medium | Medium | Needs care | P2 - Nice to Have |
| EXIF Data Viewer | Medium | Low | Yes (privacy story) | P2 - Nice to Have |
| Custom Presets | Medium | Low | Yes | P2 - Nice to Have |
| PWA / Offline | Medium | Low | Yes (privacy story) | P2 - Nice to Have |
| Batch Resize | High | High | Yes | P2 - Nice to Have |
| Drag-to-Resize Handles | Low | Medium | Yes | P3 - Future |

---

## UI Design Considerations for New Features

### Where New Features Should Live

The current UI has clear zones. New features must slot in without overwhelming the simple two-panel layout:

```
HEADER:  [Title]  [Rotate/Flip buttons]  [New Image]
         ────────────────────────────────────────────
PREVIEW: [Image with optional crop overlay]
         [Before/After slider toggle]
         [Metadata expand]
         ────────────────────────────────────────────
SIDEBAR: [Dimensions]
         [Aspect Lock]
         [Presets + Custom Presets + "+" button]
         [Scale]
         [Format]
         [Quality]
         ────────────────────────────────────────────
BOTTOM:  [Size info]  [Copy btn]  [Download btn]
```

### Key Design Principle: Progressive Disclosure

Not all features should be visible at once. Use:
- **Always visible**: Rotate/Flip (small toolbar), Paste hint (upload view)
- **On toggle**: Crop mode, Before/After comparison
- **On expand**: EXIF metadata, Keyboard shortcuts cheatsheet
- **On demand**: Custom preset creation, Batch mode (separate entry point)

### Brutalist Treatment for New Components

**Crop Overlay:**
- Solid charcoal mask outside crop area (not semi-transparent)
- Thick 3px white dashed border on crop boundary
- Chunky corner handles: 14x14px squares, coral fill, 2px charcoal border
- Aspect ratio pills below preview: `[Free] [1:1] [4:3] [16:9]`

**Transform Toolbar:**
- Horizontal strip of 4 icon buttons: `[↺] [↻] [↔] [↕]`
- 2px border, 2px offset shadow, 44x44px each
- Sits between header and preview, or above preview

**Comparison Slider:**
- 3px vertical divider line in charcoal
- Chunky drag handle: 32x16px rectangle, coral fill, centered on divider
- Labels: "Original" (left, cream badge) / "Resized" (right, cream badge)

**Copy Button:**
- Ghost button style: outlined, no fill, thick 2px border
- Clipboard icon (Phosphor Bold)
- Success state: butter yellow fill flash + "Copied!" text swap

**Dark Mode Toggle:**
- Sun/Moon icon in header
- 44x44px icon button with border
- Transition: instant swap (no slow fade — brutalism is decisive)

---

## Typography and Microcopy for New Features

| Feature | Copy | Tone |
|---------|------|------|
| Paste hint | "or paste from clipboard (Ctrl+V)" | Helpful, understated |
| Crop mode entry | "Crop" toggle button | Direct |
| Crop aspect pills | "Free", "1:1", "4:3", "16:9" | Terse, functional |
| Rotate tooltip | "Rotate left", "Rotate right" | Terse |
| Flip tooltip | "Flip horizontal", "Flip vertical" | Terse |
| Copy button | "Copy" → "Copied!" | Action → confirmation |
| Comparison toggle | "Compare" | Direct |
| EXIF expand | "Metadata stripped (tap to view)" | Informative, builds trust |
| EXIF GPS warning | "Location data found — removed" | Amber warning, reassuring |
| Custom preset save | "Save as preset" | Action-oriented |
| Keyboard cheatsheet title | "Keyboard shortcuts" | Functional |
| Batch mode toggle | "Batch mode" | Direct |
| Install PWA | "Install for offline use" | Value-first |

---

## Color Additions Needed

The existing palette covers most states. New features may need:

| Name | Hex | Usage |
|------|-----|-------|
| Overlay Mask | #1A1A1A (90% opacity) | Crop overlay outside area |
| Dark BG (dark mode) | #2C2420 | Warm dark background |
| Dark Surface (dark mode) | #3D3530 | Warm dark cards/panels |
| Dark Border (dark mode) | #5A5450 | Borders in dark mode |
| Highlight | #FFF3B0 | Keyboard shortcut badge bg |

---

## Recommended Implementation Order

**Phase 1 (Quick Wins — Low Complexity, High Impact):**
1. Paste from Clipboard — near-zero UI change, big UX win
2. Rotate/Flip buttons — simple toolbar, instant value
3. Copy to Clipboard button — one button in download bar
4. Keyboard Shortcuts — event listeners + optional cheatsheet

**Phase 2 (Medium Features):**
5. Image Cropping — most requested, medium complexity
6. Before/After Comparison Slider — polishes the quality/format experience
7. EXIF Data Viewer — reinforces privacy story
8. Custom Presets — power-user retention

**Phase 3 (Larger Initiatives):**
9. Dark Mode — requires full color system extension
10. PWA / Offline — configuration work, low UI impact
11. Batch Resize — significant UI/architecture work
12. Drag-to-Resize Handles — nice-to-have polish

---

## Risk Assessment

| Feature | Risk | Mitigation |
|---------|------|------------|
| Cropping | Scope creep — could become a full editor | Keep it basic: no free-form crop rotation, just rectangular crop with preset ratios |
| Batch | Performance — processing 20 images client-side | Web Workers for parallel processing, progress feedback, limit to 20 images |
| Dark Mode | Loses brand identity | Use warm dark tones (#2C2420), not cold grays. Test coral/yellow on dark — they pop beautifully |
| Drag Handles | Mobile touch conflicts | Use touch-action CSS, consider disabling on mobile (inputs are fine) |
| EXIF Viewer | Privacy concern — showing GPS | Make it clear data is REMOVED. Show as "we found and stripped this" |
| Clipboard API | Browser support | Feature-detect, hide button if unsupported. Safari has limitations |

---

Status: READY_FOR_PLANNING
