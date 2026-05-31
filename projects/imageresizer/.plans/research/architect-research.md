# Architect Research: Image Resizer Website

## 1. Competitor & Reference Analysis

### Existing Image Resizer Tools

| Tool | Approach | Strengths | Weaknesses |
|------|----------|-----------|------------|
| **Squoosh (Google)** | Client-side, WASM codecs | Multiple formats, side-by-side preview, advanced compression | Complex UI, overwhelming options for casual users |
| **iLoveIMG** | Server-side upload | Batch processing, multiple tools | Requires upload to server, privacy concerns, ads |
| **TinyPNG** | Server-side compression | Excellent PNG/JPEG compression | Limited to compression only, no resize, server dependency |
| **Birme** | Client-side Canvas API | Bulk resize, simple UI | Dated design, limited format support |
| **Promo Image Resizer** | Client-side | Social media presets, clean UI | Preset-locked, limited custom sizing |
| **Simple Image Resizer** | Server-side | Dead simple, one-step | Slow, server uploads, minimal control |

### Key Takeaways from Competitors
- **Best-in-class UX**: Squoosh's drag-and-drop + instant preview is the gold standard
- **Privacy advantage**: Client-side processing is a selling point — no uploads to servers
- **Common pain point**: Too many options overwhelm users. The best tools are opinionated
- **Missing niche**: No image resizer combines utility with a distinctive, memorable visual identity

---

## 2. "Friendly Brutalism" Design Research

### What is Brutalist Web Design?
Brutalist web design draws from Brutalist architecture — raw, unpolished, exposing structure. On the web, this means:
- Thick, visible borders (often black)
- Raw, system or monospace typography
- High-contrast color palettes
- Visible grid/structure
- Minimal decoration — no gradients, subtle shadows, or rounded corners
- Content-first layout

### What Makes Brutalism "Friendly"?
Friendly brutalism softens the raw aesthetic while keeping its boldness:
- **Warm, playful color palette** instead of only black/white (pastels, warm yellows, soft pinks)
- **Rounded corners on elements** — a key softener
- **Chunky, offset box-shadows** (e.g., `4px 4px 0px black`) instead of harsh flat borders
- **Playful typography** — bold sans-serifs that feel approachable (e.g., Space Grotesk, DM Sans, Archivo Black)
- **Micro-interactions** — hover effects that feel tactile (buttons that "push in")
- **Generous whitespace** — breathing room prevents visual aggression
- **Emoji or simple illustrations** as accents

### Reference Examples of Friendly Brutalism
- **Gumroad** — Bold borders, bright colors, chunky buttons, but approachable
- **Poolsuite** — Retro brutalism with warm tones and playful energy
- **The Outline (former)** — Bold typography, thick borders, colorful accents
- **Figma's community pages** — Card-based with thick borders and offset shadows
- **Notion templates market** — Clean brutalism with soft colors

### Core Visual Pattern: "Neo-Brutalism"
The specific sub-genre that fits "minimal & friendly" is often called **neo-brutalism**:
- Solid background colors (not white — use cream, pale yellow, soft lavender)
- Black thick borders (2-3px) on interactive elements
- Hard box-shadows (offset, no blur): `box-shadow: 4px 4px 0px #000`
- Bold, chunky buttons that feel "pressable"
- Limited color palette: 1 background + 1-2 accent colors + black/white
- Hover states: shadow shrinks + element translates (feels like pressing a physical button)

---

## 3. Technical Research: Client-Side Image Resizing

### Canvas API Approach (Recommended)
The HTML5 Canvas API is the standard for client-side image manipulation:

```
Flow: File Input → FileReader → Image() → Canvas.drawImage() → canvas.toBlob() → Download
```

**Key Capabilities:**
- Resize to exact pixel dimensions
- Maintain or change aspect ratio
- Output as JPEG, PNG, or WebP
- Control JPEG/WebP quality (0.0 - 1.0)
- Fast — no network round-trip

**Limitations:**
- Large images (50MP+) can cause memory issues on mobile
- No access to advanced codecs (AVIF, JPEG XL) without WASM
- Color profile handling is browser-dependent
- Maximum canvas size varies by browser (~16k x 16k pixels typically)

### Image Quality Considerations
- **Downscaling**: Use `imageSmoothingEnabled = true` and `imageSmoothingQuality = 'high'` for best results
- **Multi-step downscaling**: For extreme size reductions (>50%), stepping down in halves produces better quality than a single resize
- **Upscaling**: Canvas upscaling is bilinear interpolation — acceptable for small upscales, blurry for large ones. Should warn users.

### File Handling
- **FileReader API**: Read the uploaded file as a data URL or ArrayBuffer
- **Drag and Drop API**: `dragover` + `drop` events on a drop zone element
- **File input**: Standard `<input type="file" accept="image/*">`
- **Blob/URL.createObjectURL**: For memory-efficient large file handling
- **Download**: Create an `<a>` element with `download` attribute, set `href` to blob URL

### Format Support
| Format | Browser Read | Canvas Export | Notes |
|--------|-------------|---------------|-------|
| JPEG | All browsers | All browsers | Quality parameter 0-1 |
| PNG | All browsers | All browsers | Always lossless, larger files |
| WebP | All modern | All modern | Good compression, wide support |
| GIF | All browsers | No (becomes static) | Cannot export animated GIFs |
| SVG | All browsers | Rasterizes | Loses vector quality |
| AVIF | Chrome/Firefox | No | Would need WASM |

### Performance Notes
- Images under 10MP: instant processing on modern devices
- Images 10-30MP: sub-second on desktop, 1-3s on mobile
- Images 30MP+: may need Web Worker to avoid UI blocking
- **Recommendation**: For MVP, handle up to ~30MP without Web Workers. Add a simple loading indicator.

---

## 4. Feature Prioritization

### Must-Have (Core)
1. **Image upload** — drag-and-drop + file picker
2. **Resize by dimensions** — width/height in pixels
3. **Aspect ratio lock** — toggle to maintain proportions
4. **Live preview** — show resized image before download
5. **Download** — save resized image to device
6. **Format selection** — JPEG, PNG, WebP output
7. **Quality slider** — for JPEG/WebP (not applicable to PNG)

### Should-Have (Phase 2)
1. **Resize by percentage** — scale by 25%, 50%, 75%
2. **Preset sizes** — common dimensions (social media, wallpaper sizes)
3. **File size display** — show original and estimated output size
4. **Before/after comparison** — original vs resized dimensions shown

### Nice-to-Have (Phase 3)
1. **Batch resize** — multiple images at once
2. **Crop before resize** — simple crop tool
3. **Keyboard shortcuts** — for power users

---

## 5. Architecture Recommendations

### Single-Page Structure
```
index.html          — Semantic HTML5 structure
css/style.css       — All styles (brutalist design system)
js/app.js           — Core application logic
js/resizer.js       — Image processing (Canvas API)
```

### Why This Structure
- **Separation of concerns** without over-engineering
- `resizer.js` isolates the image processing logic for testability
- Single CSS file is sufficient — brutalist design is intentionally simple
- No build step, no bundler — matches the "minimal" ethos

### State Management
Simple state object in `app.js`:
- `originalImage`: the loaded Image object
- `originalFile`: file metadata (name, size, type)
- `settings`: { width, height, lockAspectRatio, format, quality }
- No framework needed — vanilla JS with direct DOM updates

### UX Flow
```
1. Land on page → See upload zone (prominent, centered)
2. Drop/select image → Image loads, original dimensions shown
3. Adjust settings → Width/height inputs, format dropdown, quality slider
4. See preview → Canvas renders resized version
5. Download → One-click download with sensible filename
```

---

## 6. Color Palette Direction (for UI spec)

Recommended palette for "friendly brutalism":

| Role | Color | Hex | Reasoning |
|------|-------|-----|-----------|
| Background | Warm cream | `#FFF8E7` | Softer than white, feels warm |
| Primary accent | Coral/salmon | `#FF6B6B` | Energetic but friendly |
| Secondary accent | Soft blue | `#4ECDC4` | Complementary, calming |
| Borders & text | Near-black | `#2D2D2D` | Softer than pure black |
| Card/surface | White | `#FFFFFF` | Clean containers |
| Hover/active | Warm yellow | `#FFE66D` | Playful interaction feedback |

### Typography Direction
- **Headings**: Bold, chunky sans-serif (Space Grotesk, Archivo, or DM Sans)
- **Body/UI**: Clean sans-serif (Inter or system font stack)
- **Monospace accent**: For dimensions/numbers (JetBrains Mono or system monospace)
- Load via Google Fonts CDN (aligns with static-html, no build step)

---

## 7. Key Technical Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Processing location | Client-side only | Privacy, speed, no server cost |
| Image manipulation | Canvas API | Native, no dependencies, sufficient for resize |
| State management | Vanilla JS object | Overkill to add a framework for single-page tool |
| Styling approach | Custom CSS | Brutalist design is intentionally hand-crafted |
| File handling | FileReader + Blob URLs | Standard, well-supported, memory-efficient |
| Responsive approach | CSS Grid + media queries | Simple, no framework needed |
| External dependencies | Google Fonts only | Minimal footprint, CDN-hosted |

---

## 8. Risk Assessment

| Risk | Impact | Mitigation |
|------|--------|------------|
| Large images crash mobile browsers | Medium | Cap preview size, warn on files >20MB |
| Canvas color profile mismatch | Low | Document limitation, use sRGB |
| WebP not supported in old Safari | Low | Default to JPEG, detect WebP support |
| User expects server-side features (batch, AI upscale) | Low | Clear messaging about client-side tool |

---

## Summary

This is a well-scoped project for `static-html`. The Canvas API handles all image processing needs without any server or dependencies. The "friendly brutalism" aesthetic (neo-brutalism) is well-defined with thick borders, offset shadows, warm colors, and chunky interactive elements. The architecture is intentionally flat — 4 files, no build step, no framework. Focus should be on nail the UX flow (upload → adjust → preview → download) and making the brutalist design feel polished and tactile.
