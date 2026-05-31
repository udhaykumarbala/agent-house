# UI Flow — Image Resizer

## States

1. **Upload View** (Landing)
2. **Editor View** (Workspace)
3. **Download Complete** (Temporary feedback state)

---

## Flow Diagram

```
[Page Load]
     │
     ▼
┌─────────────┐
│ Upload View │ ◄──────────────────────────────────┐
│             │                                     │
│ Drop Zone:  │                                     │
│ - Drag file │                                     │
│ - Click to  │                                     │
│   browse    │                                     │
│             │                                     │
│ Format      │                                     │
│ Badges:     │                                     │
│ [JPG][PNG]  │                                     │
│ [WebP]      │                                     │
│             │                                     │
│ Privacy:    │                                     │
│ "Your images│                                     │
│  never leave│                                     │
│  your       │                                     │
│  browser."  │                                     │
└─────┬───────┘                                     │
      │                                             │
      │ (file selected / dropped)                   │
      │                                             │
      ▼                                             │
┌─────────────┐                                     │
│ Validate    │                                     │
│ - Type OK?  │──── No ──► Show error + shake ──────┘
│ - Size OK?  │           in Upload View
└─────┬───────┘
      │ Yes
      ▼
┌──────────────────────────────────────────────┐
│ Editor View (fade in transition)             │
│                                               │
│ ┌──────────────┐  ┌────────────────────────┐ │
│ │ Preview      │  │ Controls:              │ │
│ │ (live,       │  │ - Width / Height       │ │
│ │  hover lifts │  │ - Aspect Ratio Lock    │ │
│ │  shadow)     │  │                        │ │
│ └──────────────┘  │ Presets:               │ │
│                   │ [Instagram Post]       │ │
│ ┌──────────────┐  │ [Instagram Story]      │ │
│ │ File Info:   │  │ [HD 1080p]             │ │
│ │ name, dims,  │  │ [Twitter/X Header]     │ │
│ │ size         │  │ [OG Image]             │ │
│ └──────────────┘  │                        │ │
│                   │ Scale:                 │ │
│                   │ [25%] [50%] [75%]      │ │
│                   │                        │ │
│                   │ - Format (JPG/PNG/WebP)│ │
│                   │ - Quality Slider       │ │
│                   └────────────────────────┘ │
│                                               │
│ ┌──────────────────────────────────────────┐ │
│ │ Download Bar (fixed bottom)              │ │
│ │ original size → resized size (% change)  │ │
│ │                 [DOWNLOAD] (pulse anim)  │ │
│ └──────────────────────────────────────────┘ │
└──────────────────────────────────────────────┘
      │                         │
      │ [New Image]             │ [Download]
      │                         │
      ▼                         ▼
  Back to Upload View     Save file, show
  (fade transition)       "Downloaded!" for 2s
```

---

## Footer

Visible on Upload View, hidden on Editor View (download bar replaces it).

```
┌──────────────────────────────────────────────────────────┐
│ 🔒 100% client-side · Your images never leave your       │
│    browser · EXIF data stripped automatically              │
└──────────────────────────────────────────────────────────┘
```

---

## Interactions

| Action | Element | Result |
|--------|---------|--------|
| Drag file over drop zone | Drop zone | Border solid + coral tint, text: "Let go to resize!" |
| Drag leave | Drop zone | Revert to dashed border |
| Drop valid image | Drop zone | Transition to Editor View (fadeIn 200ms) |
| Drop invalid file | Drop zone | Red flash + shake animation + error message (3s) |
| Click drop zone | File input | Opens native file picker |
| Enter/Space on drop zone | File input | Opens native file picker |
| Type in Width | Width input | Height auto-updates if aspect locked, preview updates (300ms debounce), clears active preset |
| Type in Height | Height input | Width auto-updates if aspect locked, preview updates (300ms debounce), clears active preset |
| Click aspect lock | Lock button | Toggle locked/unlocked, icon changes |
| Click preset button | Preset btn | Sets width/height to preset values, unlocks aspect ratio, highlights button (yellow bg) |
| Click scale button | Scale btn | Calculates dimensions from original at %, re-locks aspect ratio, highlights button |
| Change format | Format select | Update format, show/hide quality slider, update preview |
| Drag quality slider | Quality slider | Update quality value display, update preview (300ms debounce) |
| Click Download | Download btn | Generate blob, trigger browser download, show "Downloaded!" for 2s |
| Click New Image | New Image btn | Reset state, clear presets, revoke blob URLs, show Upload View |
| Hover any button | Buttons | translate(-2px, -2px), shadow grows to 6px (100ms ease) |
| Click/Active any button | Buttons | translate(2px, 2px), shadow collapses to 0 (100ms ease) |
| Hover preview card | Preview card | Shadow grows from lg to xl (100ms) |
| Preview ready | Download btn | Pulse animation (scale 1→1.03→1, 300ms) |

---

## Micro-Interactions

| Interaction | Animation | Duration |
|-------------|-----------|----------|
| Button hover | `translate(-2px, -2px)`, shadow grows | 100ms |
| Button press | `translate(2px, 2px)`, shadow collapses to 0 | 100ms |
| Drop zone drag-hover | Border dashed→solid, bg shifts to cream, shadow recolors to coral | 200ms |
| View transition | fadeIn (opacity 0→1, scale 0.95→1) | 200ms |
| Image preview loaded | fadeIn animation | 200ms |
| Download ready | Download button pulse (scale 1→1.03→1) | 300ms |
| Error state | Horizontal shake (translateX: 0→-4→4→-2→0) | 300ms |
| Card hover | Shadow grows from lg to xl | 100ms |
| All animations | Respect `prefers-reduced-motion: reduce` | — |

---

## API Surface

### Utils (scripts/utils.js)

| Function | Input | Output |
|----------|-------|--------|
| `Utils.formatFileSize(bytes)` | `number` | `string` — e.g., "4.2 MB" |
| `Utils.sanitizeFilename(name)` | `string` | `string` — safe filename |
| `Utils.formatToExtension(format)` | `'jpeg'\|'png'\|'webp'` | `string` — e.g., "jpg" |
| `Utils.formatToMime(format)` | `'jpeg'\|'png'\|'webp'` | `string` — e.g., "image/jpeg" |
| `Utils.mimeToFormat(mime)` | `string` | `'jpeg'\|'png'\|'webp'` |
| `Utils.getBaseName(filename)` | `string` | `string` — without extension |
| `Utils.buildDownloadFilename(original, format)` | `string, string` | `string` — e.g., "photo-resized.jpg" |
| `Utils.calcSizeChange(original, resized)` | `number, number` | `{ percent, label, smaller }` |
| `Utils.debounce(fn, delay)` | `Function, number` | `Function` |
| `Utils.validateFile(file)` | `File` | `{ valid, error }` |

### Resizer (scripts/resizer.js)

| Function | Input | Output |
|----------|-------|--------|
| `Resizer.resize(image, options)` | `HTMLImageElement, { width, height, format, quality }` | `Promise<Blob>` |

**Options:**
- `width` — target width (pixels, min 1)
- `height` — target height (pixels, min 1)
- `format` — `'jpeg'`, `'png'`, or `'webp'`
- `quality` — `0.0` to `1.0` (ignored for PNG)

**Behavior:**
- Multi-step downscaling when reduction > 50%
- `imageSmoothingEnabled = true`, `imageSmoothingQuality = 'high'`
- Returns Blob with correct MIME type

### App (scripts/app.js) — Internal Functions

| Function | Description |
|----------|-------------|
| `showUploadView()` | Switch to upload state, show footer |
| `showEditorView()` | Switch to editor state, hide footer |
| `handleFile(file)` | Validate and load an image file |
| `toggleAspectLock()` | Toggle aspect ratio lock on/off |
| `clearActivePreset()` | Deselect any active preset/scale button |
| `onPresetClick(e)` | Handle preset button click — set dimensions, unlock ratio |
| `onScaleClick(e)` | Handle scale button click — calculate from original, lock ratio |
| `updatePreview()` | Resize and display preview (debounced 300ms) |
| `downloadImage()` | Generate and trigger file download |
| `resetApp()` | Reset all state, clear presets, return to upload view |

---

## Preset Sizes

| Name | Width | Height |
|------|-------|--------|
| Instagram Post | 1080 | 1080 |
| Instagram Story | 1080 | 1920 |
| HD 1080p | 1920 | 1080 |
| Twitter/X Header | 1500 | 500 |
| OG Image | 1200 | 630 |

## Scale Options

| Scale | Behavior |
|-------|----------|
| 25% | `originalWidth * 0.25` × `originalHeight * 0.25` |
| 50% | `originalWidth * 0.50` × `originalHeight * 0.50` |
| 75% | `originalWidth * 0.75` × `originalHeight * 0.75` |

---

## File Structure

```
project/
├── index.html              # Main entry point
├── styles/
│   ├── variables.css       # CSS custom properties (design tokens)
│   ├── main.css            # Primary styles & components
│   └── responsive.css      # Media queries
├── scripts/
│   ├── app.js              # Main application logic
│   ├── resizer.js          # Canvas API resize engine
│   └── utils.js            # Helper functions
├── assets/
│   ├── images/
│   └── icons/
├── api_test.html           # Test page for Utils & Resizer
└── uiflow.md               # This file
```
