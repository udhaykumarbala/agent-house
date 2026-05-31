# Architecture Specification: Meme Creator

## Overview

A browser-based, single-page meme creator tool built with Vite + Tailwind CSS + TypeScript. Users upload images, add freely-positioned draggable text overlays, customize text styling, and export the final composition as a downloadable PNG/JPEG. All processing is 100% client-side — no backend, no accounts, no watermarks.

---

## Template

**`static-enhanced`** — Vite + Tailwind CSS + TypeScript

Rationale: No backend needed. Single canvas-based page doesn't warrant a framework like Next.js. Tailwind accelerates UI development for the controls/toolbar. TypeScript provides safety for canvas coordinate math and event handling.

---

## Tech Stack

| Layer | Technology | Justification |
|-------|-----------|---------------|
| Build | Vite | Fast HMR, minimal config, native TS support |
| Styling | Tailwind CSS | Utility-first CSS for rapid UI development |
| Language | TypeScript | Type safety for canvas coordinates, overlays, events |
| Canvas | Raw HTML5 Canvas API | No Fabric.js — our scope (text + image compositing) is well within raw Canvas capability. Avoids ~300KB dependency |
| Fonts | Google Fonts (DM Sans for UI; Impact, Anton, Bebas Neue, Bangers, Permanent Marker, Oswald for meme text) | DM Sans: friendly geometric sans for UI. Meme fonts: proven popular choices |
| Icons | Lucide (outline style) | Clean, modern, consistent. Small bundle with tree-shaking |
| Export | Canvas `toBlob()` / `toDataURL()` | Native browser API, strips EXIF metadata automatically |

### No External Dependencies For Core Logic

The canvas rendering, drag handling, text overlay management, and export are all vanilla TypeScript. No canvas libraries (Fabric.js, Konva), no state management libraries (Redux, Zustand), no UI frameworks (React, Vue).

---

## Architecture

### Core Modules

```
src/
  main.ts           -> App initialization, module wiring, event setup
  canvas.ts         -> Canvas rendering engine (draw image, draw overlays, hit-testing)
  overlays.ts       -> TextOverlay CRUD, state management, selection
  drag.ts           -> Pointer event handling, coordinate transforms, drag logic
  export.ts         -> Export to PNG/JPEG at full resolution with scaling
  ui.ts             -> Controls panel bindings (font, color, size, stroke inputs)
  upload.ts         -> File upload validation (type, magic bytes, size, dimensions)
  history.ts        -> Undo/redo snapshot stack
  types.ts          -> TypeScript interfaces and type definitions
  style.css         -> Tailwind directives + custom canvas styles

index.html          -> Single page: canvas + controls layout
```

### Data Flow

```
User Action -> Event Handler -> State Update -> Canvas Re-render
                                    |
                                    v
                              History Snapshot
```

All state mutations flow through a central state object. Every mutation triggers a canvas re-render via `requestAnimationFrame`. History snapshots are captured after each discrete user action (add text, move text, change style, delete text).

---

## Data Model

### Core State

```typescript
interface AppState {
  image: HTMLImageElement | null
  imageFileName: string
  overlays: TextOverlay[]
  selectedId: string | null
  canvasWidth: number
  canvasHeight: number
  exportFormat: 'png' | 'jpeg'
  exportQuality: number  // 0.0 - 1.0, for JPEG
}
```

### TextOverlay

```typescript
interface TextOverlay {
  id: string              // unique identifier (crypto.randomUUID())
  text: string            // user-entered text content
  x: number              // canvas x coordinate (center of text)
  y: number              // canvas y coordinate (baseline of text)
  fontSize: number       // px (range: 16-200)
  fontFamily: string     // Google Font name
  fillColor: string      // hex color for text fill
  strokeColor: string    // hex color for text outline
  strokeWidth: number    // px (0 = no stroke)
  bold: boolean
  italic: boolean
  allCaps: boolean
  textAlign: 'left' | 'center' | 'right'
  shadowEnabled: boolean
  shadowColor: string
  shadowBlur: number
  shadowOffsetX: number
  shadowOffsetY: number
  opacity: number        // 0.0 - 1.0
}
```

### Default Text Overlay (Meme Preset)

```typescript
const MEME_PRESET: Partial<TextOverlay> = {
  text: 'YOUR TEXT HERE',
  fontSize: 48,
  fontFamily: 'Impact',
  fillColor: '#FFFFFF',
  strokeColor: '#000000',
  strokeWidth: 3,
  bold: false,
  italic: false,
  allCaps: true,
  textAlign: 'center',
  shadowEnabled: false,
  opacity: 1.0,
}
```

---

## UI Layout

### Desktop (Primary)

```
+-------------------------------------------------------+
|  [Logo] MemeForge              [Undo][Redo] [Download] |
+----------------------------+---------------------------+
|                            |   Controls Panel          |
|                            |   ----------------------  |
|                            |   [+ Add Text]            |
|       Canvas Area          |   [Meme Preset]           |
|    (Image + Text Overlays) |                           |
|                            |   Font: [dropdown]        |
|                            |   Size: [===slider===]    |
|                            |   Color: [picker] Stroke: |
|                            |   Bold [B] Italic [I]     |
|                            |   Align: [L] [C] [R]     |
|                            |   CAPS: [toggle]          |
|                            |   ----------------------  |
|                            |   Text Layers             |
|                            |   [Text 1]         [x]   |
|                            |   [Text 2]         [x]   |
+----------------------------+---------------------------+
```

- Canvas area: ~70% width on desktop
- Controls panel: ~30% width, right sidebar
- Controls are only active when a text overlay is selected
- Canvas sits on a subtle dark background to clearly delineate the editing zone

### Mobile (Secondary / Responsive)

- Canvas fills full width
- Controls panel moves below canvas
- Collapsible sections to save space
- Larger touch targets (min 44x44px)

---

## Visual Design System

Reference: UI Research decisions

### Colors

| Token | Value | Usage |
|-------|-------|-------|
| `--bg-primary` | `#1A1A2E` | Main background, deep charcoal |
| `--bg-surface` | `#16213E` | Panels, cards, sidebar |
| `--bg-surface-hover` | `#252547` | Hover state for surfaces |
| `--accent` | `#E91E8C` | Primary accent, hot magenta |
| `--accent-hover` | `#D4177F` | Accent hover state |
| `--text-primary` | `#F8F8F2` | Primary text, off-white |
| `--text-secondary` | `#A0A0B8` | Secondary text, muted lavender-gray |
| `--border` | `#2A2A4A` | Subtle borders between panels |
| `--canvas-bg` | `#0D0D1A` | Canvas surrounding area |
| `--success` | `#22C55E` | Success states (download complete) |
| `--error` | `#EF4444` | Error states (invalid file) |

### Typography

- **UI Font**: DM Sans (Google Fonts) — all interface text
- **Meme Fonts** (on-canvas): Impact (default), Anton, Bebas Neue, Bangers, Permanent Marker, Oswald

### Corners & Spacing

- Buttons: `rounded-lg` (8px)
- Panels/cards: `rounded-xl` (12px)
- Inputs: `rounded-md` (6px)
- Spacing scale: Tailwind default (4px base unit)

---

## Key Implementation Details

### Canvas Rendering Pipeline

1. Clear canvas
2. Draw background image (scaled to fit canvas display area)
3. For each overlay in order:
   a. Set font, fill, stroke, shadow properties on context
   b. If `allCaps`, transform text to uppercase
   c. Call `ctx.fillText()` then `ctx.strokeText()`
4. If overlay is selected, draw selection border (magenta dashed rect) around text bounding box

Rendering is triggered via `requestAnimationFrame` — never directly in event handlers.

### Drag System (Pointer Events)

- `pointerdown` on canvas: Hit-test all overlays via `ctx.measureText()` bounding boxes. Select topmost hit. Begin drag if hit.
- `pointermove`: If dragging, update overlay x/y, request re-render.
- `pointerup`: End drag state.
- CSS: `touch-action: none` on canvas element to prevent scroll interference.

### Hit Testing

For each text overlay, compute bounding box:
```
width = ctx.measureText(text).width
height = fontSize * 1.2  (approximate line height)
```
Check if pointer coordinates fall within the bounding box. Test overlays in reverse order (topmost first).

### Image Upload Flow

1. User drops image or clicks file input
2. Validate: file type whitelist, file extension, magic bytes, file size (max 10MB)
3. Validate: load as `Image`, check `naturalWidth`/`naturalHeight` <= 4096
4. Compute canvas display size: fit image within max display area (e.g., 800x600) maintaining aspect ratio
5. Store original `Image` object for full-resolution export
6. Render image on canvas, enable controls

### Export Flow

1. Create offscreen canvas at original image dimensions
2. Draw original image at full size
3. Scale all text overlay positions proportionally: `exportX = overlay.x * (originalWidth / displayWidth)`
4. Scale font sizes proportionally: `exportFontSize = overlay.fontSize * (originalWidth / displayWidth)`
5. Render all overlays on offscreen canvas
6. Call `canvas.toBlob(callback, mimeType, quality)`
7. Create `<a>` with `URL.createObjectURL(blob)`, set `download` attribute, trigger click
8. Revoke object URL after download

### Undo/Redo

- Snapshot-based: Deep clone `overlays` array on each action
- Max 30 snapshots in history stack
- `Ctrl+Z` / `Cmd+Z` for undo, `Ctrl+Shift+Z` / `Cmd+Shift+Z` for redo
- Actions that create snapshots: add text, delete text, move text (on pointerup), change any style property

---

## Security Requirements

Reference: Security Research

### File Upload Validation (Critical)

1. Whitelist MIME types: `image/jpeg`, `image/png`, `image/webp`, `image/gif`
2. Whitelist extensions: `.jpg`, `.jpeg`, `.png`, `.webp`, `.gif`
3. **Block SVG** — can contain embedded JavaScript
4. Validate magic bytes (file signatures): JPEG `FF D8 FF`, PNG `89 50 4E 47`, WebP `52 49 46 46...57 45 42 50`, GIF `47 49 46 38`
5. Max file size: 10MB
6. Max dimensions: 4096x4096
7. Validate image loads successfully via `new Image()`

### XSS Prevention

- Use `textContent` exclusively for DOM text — never `innerHTML` with user input
- Canvas `fillText()`/`strokeText()` is inherently safe (pixel rendering)
- No `eval()`, `new Function()`, or `document.write()` with user data

### Export Security

- Sanitize/hardcode download filename: `meme_[timestamp].png`
- Use `URL.createObjectURL()` with properly typed Blob
- `URL.revokeObjectURL()` after download to free memory

### CSP Headers (Deployment)

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

---

## Performance Considerations

- Use `requestAnimationFrame` for all canvas re-renders (batch with drag events)
- Debounce text input changes (150ms) before re-render
- Downscale images over 4096px for editing display; export at original size
- Keep history stack at max 30 entries to limit memory
- Load Google Fonts with `font-display: swap` to avoid blocking render

---

## Browser Support

Target: Modern evergreen browsers (Chrome, Firefox, Safari, Edge — latest 2 versions). No IE11 support. No polyfills required.

---

## Accessibility

- Keyboard navigation: Tab between controls, Enter to add text
- Arrow keys for fine-positioning selected text overlay (1px per press, 10px with Shift)
- Focus indicators on all interactive elements
- `aria-live` region for announcing selection changes and export status
- Sufficient color contrast in controls panel (WCAG AA on dark background)
- All controls labeled with visible text or `aria-label`

---

## Trade-offs

| Optimizing For | Sacrificing |
|---------------|-------------|
| Simplicity — no framework, no heavy libs | Scalability for complex future features |
| Speed — zero backend, instant load | Server-side features (sharing, templates gallery) |
| Privacy — all client-side | Analytics, usage tracking |
| Bundle size — raw Canvas API | Convenience of Fabric.js object model |
| TypeScript safety | Slightly more boilerplate than plain JS |

---

## Out of Scope (v1)

- Meme template gallery
- Image filters/effects
- Text rotation
- Multi-line text auto-wrapping
- Social sharing integration
- User accounts or cloud save
- Backend of any kind
- Mobile-optimized touch gestures (basic touch works via pointer events, but no pinch-to-zoom/rotate)

---

*Specification created: 2026-03-16*
*Based on research from: product, architecture, UX, UI, and security teams*
