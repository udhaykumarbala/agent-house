# Architecture Specification: Image Resizer

## Overview

A single-page, client-side image resizer with a friendly brutalist (neo-brutalist) aesthetic. All image processing happens in the browser via the HTML5 Canvas API. Zero backend, zero dependencies, zero build tools.

---

## Template

**`static-html`** — Pure HTML/CSS/JS, no build step, no node_modules. Matches the minimal ethos of both the product and the brutalist design.

---

## Tech Stack

| Layer | Technology | Justification |
|-------|-----------|---------------|
| Markup | HTML5 (semantic) | Native, no templating needed for single page |
| Styling | Custom CSS3 | Brutalist design demands hand-crafted styles; Tailwind would fight the aesthetic |
| Logic | Vanilla JavaScript (ES6+) | Canvas API + File API are native; no framework needed for this scope |
| Image Processing | HTML5 Canvas API | Client-side resize, format conversion, quality control — no server round-trip |
| File Handling | File API + Drag & Drop API + Blob URLs | Standard browser APIs for upload, preview, and download |
| Fonts | Google Fonts (Space Grotesk + Inter) | Loaded via `<link>`, no build step. Consider self-hosting for privacy (see Security) |
| Icons | Phosphor Icons (bold weight, CDN) | Thick-stroke icons match the brutalist border weight (2-3px) |

### External Dependencies (Minimal)

1. **Google Fonts** — Space Grotesk (headings) + Inter (body)
2. **Phosphor Icons** — Bold weight line icons via CDN
3. Nothing else. Zero npm packages, zero build tools.

---

## File Structure

```
imageresizer/
├── index.html              # Single-page HTML structure
├── css/
│   └── style.css           # All styles — brutalist design system
├── js/
│   ├── app.js              # UI logic, state management, DOM interactions
│   └── resizer.js          # Image processing module (Canvas API)
├── assets/
│   └── (favicon, og-image) # Static assets only
└── .plans/                 # Planning docs (not deployed)
```

### Why This Structure
- **Separation of concerns** without over-engineering
- `resizer.js` isolates Canvas API logic — testable, replaceable
- `app.js` handles all DOM/UI — single source of truth for state
- Single CSS file — brutalist design is intentionally simple, no need for partials
- No build step = open `index.html` and go

---

## Application Architecture

### State Model

A single state object in `app.js`:

```
state = {
  originalImage: Image | null,       // The loaded Image element
  originalFile: {                     // Source file metadata
    name: string,
    size: number,                     // bytes
    type: string,                     // MIME type
    width: number,
    height: number
  } | null,
  settings: {
    width: number,
    height: number,
    lockAspectRatio: boolean,         // default: true
    format: 'jpeg' | 'png' | 'webp', // default: 'jpeg'
    quality: number                   // 0.0-1.0, default: 0.85
  },
  ui: {
    view: 'upload' | 'editor',       // current screen state
    isDragging: boolean,              // drag-over visual state
    isProcessing: boolean             // show loading indicator
  }
}
```

No framework, no reactivity library. Direct DOM updates via helper functions when state changes.

### Data Flow

```
User Action → State Update → DOM Update → (if settings changed) → Canvas Re-render

Detailed flow:
1. File Input/Drop → FileReader.readAsDataURL() → new Image() → state.originalImage
2. Image.onload → read naturalWidth/naturalHeight → populate state.settings
3. User adjusts width/height → aspect ratio calc → update state.settings
4. Preview button or auto-preview → resizer.resize(image, settings) → canvas.toBlob()
5. Download click → create <a download> with blob URL → trigger click → revoke URL
```

### Module Responsibilities

**`app.js`** — UI Controller
- Initialize event listeners (drag/drop, file input, form controls)
- Manage application state
- Update DOM when state changes
- Handle user interactions (dimension inputs, format select, quality slider)
- Coordinate between UI and resizer module
- Sanitize all user-derived content before DOM insertion (textContent only, never innerHTML)

**`resizer.js`** — Image Processor
- `resize(image, options)` → returns Promise<Blob>
- `getPreviewDataUrl(image, options)` → returns data URL for preview
- Multi-step downscaling for quality (halving approach for >50% reductions)
- Format conversion (JPEG, PNG, WebP)
- Quality control for lossy formats
- File size estimation

---

## Design System: Friendly Brutalism

### Color Palette — "Butter & Charcoal"

Based on UI research direction #1 (strongest recommendation):

| Token | Hex | Usage |
|-------|-----|-------|
| `--bg` | `#FFF8E7` | Page background (warm cream) |
| `--surface` | `#FFFFFF` | Cards, panels, input backgrounds |
| `--primary` | `#FF6B6B` | Primary CTA (coral) — download button, key actions |
| `--secondary` | `#4ECDC4` | Secondary actions, success states (teal) |
| `--accent` | `#FFE66D` | Hover feedback, highlights (warm yellow) |
| `--text` | `#1A1A1A` | Body text, headings |
| `--border` | `#1A1A1A` | All thick borders |
| `--muted` | `#6B7280` | Secondary text, placeholders |
| `--error` | `#E53E3E` | Error states |

### Typography

| Role | Font | Weight | Size Range |
|------|------|--------|------------|
| Headings | Space Grotesk | 700 | 28-48px |
| Body / UI labels | Inter | 400-500 | 16-18px |
| Dimension inputs | System monospace | 500 | 20-24px |
| Button text | Space Grotesk | 600 | 16-18px |

### Component Patterns

**Borders**: 2-3px solid `#1A1A1A` on all interactive elements (buttons, cards, inputs)

**Shadows**: Hard offset, no blur — `4px 4px 0px #1A1A1A`

**Corners**: 4-6px border-radius (slightly rounded, not pills)

**Buttons**:
- Default: solid fill + thick border + offset shadow
- Hover: `translate(-2px, -2px)` + shadow grows to `6px 6px`
- Active: `translate(2px, 2px)` + shadow to `0px 0px` (pressed feel)

**Inputs**: Thick bordered, large (min 48px height), monospace for numbers

**Cards**: White surface, thick border, offset shadow, generous padding (24-32px)

### Layout

- Single column, max-width 900px, centered
- **Upload state**: Full-viewport drop zone, centered content
- **Editor state**: Stacked vertically — preview on top, controls below, download bar at bottom
- Mobile: Same stacked layout, controls become full-width
- Breakpoint: 768px for minor adjustments (the stacked layout works at all sizes)

---

## UX Flow

### State 1: Upload (Landing)
- Full-screen drop zone with dashed border
- "Drop your image here" heading + file picker button
- Supported formats listed (JPG, PNG, WebP)
- Privacy message: "Your images never leave your browser"

### State 2: Editor (After Upload)
- Image preview (scaled to fit, maintaining aspect ratio)
- Original dimensions & file size displayed
- **Controls panel**:
  - Width / Height inputs with aspect ratio lock toggle
  - Preset buttons (Instagram 1080x1080, HD 1920x1080, Twitter 1200x675, 50%, 25%)
  - Format dropdown (JPEG, PNG, WebP)
  - Quality slider (for JPEG/WebP, hidden for PNG)
- **Action bar**:
  - File size comparison: "4.2MB → 380KB (91% smaller)"
  - Download button (primary CTA, large, prominent)
  - "New Image" button to reset

### Transitions
- Upload → Editor: Fade/slide transition, drop zone shrinks away
- Editor → Upload: Reset state, show drop zone again

---

## Security Implementation

Based on security research findings:

| Concern | Implementation |
|---------|---------------|
| XSS Prevention | `textContent` only for all user-derived values (file names, dimensions). Never `innerHTML`. |
| File Validation | Check `file.type` starts with `image/`, enforce 50MB max size, wrap Image loading in error handler |
| Blob URL Lifecycle | `URL.revokeObjectURL()` on every new upload and before download URL creation |
| Download File Name | Sanitize output filename — strip `../`, null bytes, control chars. Pattern: `{original-name}-resized.{format}` |
| CSP Headers | Configure at deploy time: `script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' blob: data:` |
| External Resources | Use SRI hashes for any CDN-loaded resources (fonts, icons) |

---

## Performance Considerations

| Scenario | Approach |
|----------|----------|
| Images < 10MP | Direct Canvas resize, instant |
| Images 10-30MP | Canvas resize with loading indicator |
| Images > 30MP | Warn user, proceed with loading indicator |
| Files > 50MB | Reject with friendly error message |
| Multi-step downscaling | For reductions > 50%, halve dimensions in steps for better quality |
| Preview rendering | Use lower-quality preview, full quality on download |

---

## Browser Support

- Modern evergreen browsers: Chrome, Firefox, Safari, Edge (last 2 versions)
- Canvas API: universally supported
- WebP export: Chrome, Firefox, Edge (Safari 16+). Detect support, fallback to JPEG.
- Drag & Drop API: universally supported, with file input fallback for mobile

---

## Key Architectural Decisions

| Decision | Choice | Alternative Considered | Why |
|----------|--------|----------------------|-----|
| Client-side only | Yes | Server-side processing | Privacy, speed, zero hosting cost, no GDPR concerns |
| No framework | Vanilla JS | React, Vue, Svelte | Single page with ~10 interactive elements doesn't justify framework overhead |
| No build tools | None | Vite, Webpack | Matches minimal ethos, zero config, instant dev workflow |
| Single CSS file | `style.css` | CSS modules, Tailwind | Brutalist design is small and intentional; partials add complexity without benefit |
| Canvas API | Native | WASM (Squoosh-style) | Sufficient for resize + format conversion; WASM is overkill for this scope |
| Google Fonts CDN | Yes | Self-hosted fonts | Simpler for MVP; can self-host later for privacy hardening |

---

## What This Architecture Does NOT Include (By Design)

- No batch processing (single image focus for MVP)
- No crop tool (resize only)
- No AVIF/JPEG-XL support (would need WASM)
- No Web Workers (not needed for <30MP images on modern devices)
- No service worker / PWA (can add later)
- No analytics or tracking
- No backend of any kind

---

Status: READY_FOR_DEVELOPMENT_PLAN
