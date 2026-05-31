# Architecture Research: Browser-Based Meme Creator

## 1. Competitive Analysis

### Existing Meme Creators

| Tool | Approach | Strengths | Weaknesses |
|------|----------|-----------|------------|
| **Imgflip** | Server-rendered templates | Massive template library, simple UX | Requires account for watermark-free, server-dependent |
| **Kapwing** | Canvas-based SPA | Professional editor, layers panel | Heavy framework, slow load, freemium gates |
| **Canva** | Canvas + WebGL | Full design suite, drag-and-drop | Overkill for memes, requires login |
| **Make a Meme** | Simple form-based | Zero learning curve | No free positioning, limited styling |
| **Photopea** | Full Photoshop clone | Extreme capability | Overwhelming for meme creation |

### Key Takeaways
- **Sweet spot**: Between "form with top/bottom text" (too simple) and "full design editor" (too complex)
- **Must-haves users expect**: Free text placement, font selection, text outline/stroke, instant download
- **Differentiator opportunity**: Zero login, no watermarks, no server dependency, instant load

---

## 2. Core Technology: HTML5 Canvas API

### Why Canvas Over DOM-Based Text

| Approach | Pros | Cons |
|----------|------|------|
| **Canvas API** | Pixel-perfect export, single surface, consistent rendering, stroke/shadow built-in | Custom hit-testing for drag, manual text input handling |
| **DOM overlays + html2canvas** | Native drag via CSS, easier text editing | Export fidelity issues, cross-origin problems, inconsistent rendering |
| **Fabric.js (Canvas lib)** | Built-in object model, drag/resize/rotate, serialization | 300KB+ dependency, learning curve, may be overkill |

**Recommendation: Raw Canvas API with a thin abstraction layer**

Rationale:
- The feature set (text + image compositing) is well within what raw Canvas handles
- Avoids large dependency (Fabric.js is ~300KB minified)
- Full control over rendering pipeline and export quality
- Canvas `toBlob()` / `toDataURL()` gives clean PNG/JPEG export
- Custom drag logic is ~50 lines of code for this use case

### Key Canvas APIs Needed

```
// Core rendering
ctx.drawImage(img, 0, 0)           // Background image
ctx.fillText(text, x, y)           // Text rendering
ctx.strokeText(text, x, y)         // Text outline
ctx.measureText(text)              // Hit-testing & sizing

// Export
canvas.toBlob(callback, 'image/png')   // Download as PNG
canvas.toDataURL('image/jpeg', 0.9)    // Download as JPEG

// Text styling
ctx.font = '48px Impact'
ctx.fillStyle = '#ffffff'
ctx.strokeStyle = '#000000'
ctx.lineWidth = 3
ctx.textAlign = 'center'
ctx.shadowColor / shadowBlur / shadowOffsetX / shadowOffsetY
```

---

## 3. Drag-and-Drop Text Implementation

### Approach: Coordinate-Based Hit Testing

Each text overlay is stored as an object:
```typescript
interface TextOverlay {
  id: string
  text: string
  x: number        // canvas coordinates
  y: number
  fontSize: number
  fontFamily: string
  fillColor: string
  strokeColor: string
  strokeWidth: number
  shadowEnabled: boolean
  isDragging: boolean
}
```

### Drag Flow
1. `mousedown` / `touchstart` → Check if click hits any text bounding box via `ctx.measureText()`
2. `mousemove` / `touchmove` → Update `x, y` of dragged text, re-render canvas
3. `mouseup` / `touchend` → Release drag state

### Touch Support
- Use pointer events (`pointerdown`, `pointermove`, `pointerup`) for unified mouse+touch
- Add `touch-action: none` CSS on canvas to prevent scroll interference
- Consider `pointer-events: none` on overlapping UI during drag

---

## 4. Image Upload & Handling

### Upload Methods
1. **File input** (`<input type="file" accept="image/*">`) — standard, reliable
2. **Drag-and-drop zone** — enhanced UX, uses `dragenter/dragover/drop` events
3. **Paste from clipboard** — power user feature, `paste` event + `clipboardData.items`

### Image Processing
- Load via `FileReader.readAsDataURL()` → `new Image()` → `ctx.drawImage()`
- **Canvas sizing strategy**: Fit image to a max display size (e.g., 800×600) while maintaining aspect ratio
- Store original dimensions for high-quality export
- Handle EXIF orientation (modern browsers handle this natively with `image-orientation: from-image`)

### Cross-Origin Considerations
- Not an issue for local file uploads (no CORS)
- If adding URL-based image loading later, use `img.crossOrigin = "anonymous"`

---

## 5. Text Styling Options

### Essential (MVP)
| Feature | Implementation |
|---------|---------------|
| Font family | Dropdown with web-safe fonts + Google Fonts subset |
| Font size | Slider or number input (16–120px) |
| Text color | `<input type="color">` for fill color |
| Outline/stroke | Color picker + width slider (the classic meme look) |
| Bold/Italic | Toggle buttons modifying `ctx.font` string |
| Text alignment | Left/Center/Right buttons |

### Nice-to-Have (Phase 2)
| Feature | Implementation |
|---------|---------------|
| Text shadow | Shadow color/blur/offset controls |
| All-caps toggle | `text.toUpperCase()` on render |
| Line spacing | Custom multi-line text rendering |
| Text background | Filled rect behind text bounding box |
| Opacity | `ctx.globalAlpha` per text element |
| Rotation | `ctx.rotate()` around text center point |

### Font Strategy
- Bundle 6-8 popular meme fonts via Google Fonts: Impact, Arial Black, Comic Sans MS, Roboto Bold, Oswald, Bebas Neue, Permanent Marker, Bangers
- Load via `@import` or `<link>` with `font-display: swap`
- Impact is the "meme default" — pre-select it

---

## 6. Export Functionality

### Export Flow
1. Re-render canvas at export resolution (original image size, not display size)
2. Scale all text overlays proportionally
3. Call `canvas.toBlob()` for the export format
4. Create `<a>` element with `URL.createObjectURL(blob)`, set `download` attribute, trigger click

### Format Options
| Format | Use Case | Method |
|--------|----------|--------|
| PNG | Default, lossless, supports transparency | `canvas.toBlob(cb, 'image/png')` |
| JPEG | Smaller file size, photos | `canvas.toBlob(cb, 'image/jpeg', 0.92)` |

### Quality Considerations
- Render at original image resolution, not CSS display size
- Use `window.devicePixelRatio` awareness for retina displays during editing
- Offer quality slider for JPEG export

---

## 7. State Management

### Simple Array-Based State
No external state library needed. A plain array of `TextOverlay` objects + a render loop:

```
state = {
  image: HTMLImageElement | null
  overlays: TextOverlay[]
  selectedId: string | null
  canvasWidth: number
  canvasHeight: number
}
```

### Undo/Redo
- Snapshot-based: Push full state to history array on each action
- Keep max 20-30 snapshots to limit memory
- `Ctrl+Z` / `Ctrl+Shift+Z` keyboard shortcuts

---

## 8. UI Layout Pattern

### Common Editor Layout (from competitor analysis)

```
┌─────────────────────────────────────────────┐
│  Header: Logo / Title / Export Button        │
├──────────────────────────┬──────────────────┤
│                          │  Controls Panel   │
│                          │  ─────────────── │
│      Canvas Area         │  + Add Text      │
│      (Image + Text)      │  Font Family     │
│                          │  Font Size       │
│                          │  Colors          │
│                          │  Stroke          │
│                          │  Alignment       │
│                          │  ─────────────── │
│                          │  Text List       │
│                          │  [Text 1]  🗑    │
│                          │  [Text 2]  🗑    │
├──────────────────────────┴──────────────────┤
│  Footer (optional): Format / Quality         │
└─────────────────────────────────────────────┘
```

### Mobile Layout
- Stack controls below canvas
- Collapsible controls panel
- Larger touch targets for drag handles

---

## 9. Performance Considerations

- **requestAnimationFrame** for smooth drag rendering (don't redraw on every mousemove directly)
- **Debounce** text input changes before re-render
- **Canvas size limits**: Max ~4096×4096 on most mobile browsers, ~16384×16384 on desktop
- **Memory**: Large images consume significant memory when loaded into canvas. Consider downscaling images over 4000px for editing, export at original size

---

## 10. Browser Compatibility

| Feature | Support |
|---------|---------|
| Canvas 2D | All modern browsers (IE11+) |
| Pointer Events | All modern browsers |
| File API / FileReader | All modern browsers |
| canvas.toBlob() | All modern (polyfill for old Edge) |
| input[type=color] | All modern browsers |
| Google Fonts | All browsers |

No polyfills needed for modern browser targeting.

---

## 11. Accessibility Considerations

- Keyboard support: Tab between text overlays, arrow keys for fine positioning
- Screen reader: Announce selected text, canvas state changes via `aria-live`
- High contrast controls panel
- Focus indicators on all interactive elements

---

## 12. Recommended Architecture Summary

```
Template: static-enhanced (Vite + Tailwind + TypeScript)

src/
  main.ts          → App initialization, event wiring
  canvas.ts        → Canvas rendering engine (draw image, draw texts, hit-test)
  overlays.ts      → TextOverlay CRUD, state management
  drag.ts          → Pointer event handling, coordinate transforms
  export.ts        → Export to PNG/JPEG with resolution scaling
  ui.ts            → Controls panel bindings (font, color, size inputs)
  history.ts       → Undo/redo stack
  types.ts         → TypeScript interfaces
  style.css        → Tailwind imports + custom canvas styles

index.html         → Single page with canvas + controls layout
```

### Key Decisions
1. **Raw Canvas API** over Fabric.js — smaller bundle, sufficient for scope
2. **Pointer Events** over separate mouse/touch — unified, modern
3. **Snapshot undo** over command pattern — simpler for this scale
4. **TypeScript** — canvas coordinate math benefits from type safety
5. **No framework** — vanilla TS + Tailwind, no React/Vue needed for single-page canvas app

---

## 13. Risk Assessment

| Risk | Mitigation |
|------|------------|
| Text rendering differs across browsers | Use common web-safe fonts + Google Fonts |
| Large images cause lag | Downscale for editing, export at full resolution |
| Mobile drag UX is poor | Pointer events + touch-action CSS + larger hit targets |
| Canvas export quality loss | Render at original image dimensions for export |
| Multi-line text complexity | Custom line-break rendering with `measureText()` width checks |

---

*Research completed: 2026-03-16*
*Next: Architecture specification and development plan in Phase 2*
