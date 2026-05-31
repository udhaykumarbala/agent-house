# UX Specification: Image Resizer (Minimal Friendly Brutalism)

Based on research in: .plans/research/ux-research.md, .plans/research/ui-research.md
Implementing features from: .plans/research/product-research.md
Security guidance from: .plans/research/security.md
Architecture from: .plans/research/architect-research.md

---

## Design Philosophy

**One sentence**: A single-purpose, client-side image resizer that feels like a friendly workshop — bold, honest, and fast.

**Principles (ranked)**:
1. **Zero friction** — Drop image → adjust → download. Under 30 seconds.
2. **Privacy by default** — Nothing leaves the browser. Say it loud.
3. **Friendly brutalism** — Thick borders, hard shadows, warm colors, playful copy.
4. **One job, done well** — Resist feature creep. This is not Photoshop.

---

## Screen List

1. **Landing / Upload** — Full-viewport drop zone. The entire page IS the action.
2. **Workspace** — Image preview + resize controls + download bar. The core tool.
3. **States** — Empty, loading, error, success (overlays/inline, not separate pages).

Everything is a **single page**. No routing, no navigation, no page transitions.

---

## Wireframes

### Screen 1: Landing / Upload State

```
┌──────────────────────────────────────────────────────┐
│                                                      │
│   IMAGE RESIZER                     privacy note ↗   │
│                                                      │
├──────────────────────────────────────────────────────┤
│ ┌──────────────────────────────────────────────────┐ │
│ │ ┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄ │ │
│ │ ┄                                              ┄ │ │
│ │ ┄                                              ┄ │ │
│ │ ┄              ┌────────────┐                  ┄ │ │
│ │ ┄              │   📁 icon  │                  ┄ │ │
│ │ ┄              └────────────┘                  ┄ │ │
│ │ ┄                                              ┄ │ │
│ │ ┄      Drop your image here                    ┄ │ │
│ │ ┄      (or click to browse)                    ┄ │ │
│ │ ┄                                              ┄ │ │
│ │ ┄      JPG, PNG, WebP — up to 50MB            ┄ │ │
│ │ ┄                                              ┄ │ │
│ │ ┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄ │ │
│ └──────────────────────────────────────────────────┘ │
│                                                      │
│   Your images never leave your browser.              │
│                                                      │
└──────────────────────────────────────────────────────┘
```

**Key Elements**:
- **Title "IMAGE RESIZER"**: Top-left, Space Grotesk Bold, all-caps, 36-48px. Establishes the tool identity immediately.
- **Privacy note link**: Top-right, small text. Opens a brief inline tooltip/popover: "All processing happens in your browser. Nothing is uploaded to any server."
- **Drop zone**: Dashed thick border (3px, #1A1A2E), cream fill (#FEFAE0). Takes up ~60% of viewport height. Entire zone is clickable and droppable.
- **File icon**: Simple thick-stroke icon (Phosphor Bold), centered.
- **Primary copy**: "Drop your image here" — Space Grotesk Bold, 24px.
- **Secondary copy**: "(or click to browse)" — Inter Medium, 16px, slightly muted.
- **Format hint**: "JPG, PNG, WebP — up to 50MB" — Inter, 14px, muted text.
- **Privacy statement**: Bottom of page, Inter, 14px. Always visible. Builds trust.

**Interaction**:
- Dragging a file over the zone: border changes from dashed to solid, background shifts to light coral tint, copy changes to "Let go to resize!"
- Clicking anywhere in the zone triggers native file picker (`accept="image/*"`)
- Invalid file type: zone border flashes red, shows "That's not an image file" briefly

---

### Screen 2: Workspace State (Desktop)

```
┌──────────────────────────────────────────────────────┐
│  IMAGE RESIZER                  [↻ New Image]  [⚙]  │
├────────────────────────────┬─────────────────────────┤
│                            │                         │
│                            │  DIMENSIONS             │
│                            │  ┌────────┐  ┌────────┐│
│                            │  │ Width  │  │ Height ││
│     ┌──────────────────┐   │  │  1920  │  │  1080  ││
│     │                  │   │  └────────┘  └────────┘│
│     │                  │   │       [🔗 Locked]       │
│     │  Image Preview   │   │                         │
│     │   (fitted to     │   │  PRESETS                │
│     │    available     │   │  ┌──────────────────┐  │
│     │    space)        │   │  │ Instagram 1080²  │  │
│     │                  │   │  │ HD 1920×1080     │  │
│     │                  │   │  │ Twitter 1200×675 │  │
│     └──────────────────┘   │  │ OG Image 1200×630│  │
│                            │  │ 50%  │  25%  │ …  │  │
│     original.jpg           │  └──────────────────┘  │
│     4032 × 3024 · 4.2 MB   │                         │
│                            │  FORMAT                 │
│                            │  [JPG ▾] Quality [85──]│
│                            │                         │
├────────────────────────────┴─────────────────────────┤
│  4.2 MB → ~380 KB (91% smaller)     [ ⬇ DOWNLOAD ] │
└──────────────────────────────────────────────────────┘
```

**Key Elements**:

**Header Bar**:
- Title "IMAGE RESIZER": Same as landing, top-left anchor.
- "New Image" button: Top-right, secondary style (outlined, thick border). Resets workspace to landing state.
- Settings gear icon: Top-right, opens advanced settings (hidden by default).

**Preview Area (left, ~60% width)**:
- Image displayed fit-to-container with checkerboard background for transparency (PNG/WebP).
- Below image: original filename (as `textContent`, never `innerHTML` — per security research), original dimensions, original file size. Inter, 14px, muted.
- Image has a thick 2px border and hard shadow (4px 4px 0 #1A1A2E).

**Controls Panel (right sidebar, ~40% width)**:
- **DIMENSIONS section**:
  - Two large input fields: Width and Height. Monospace font (JetBrains Mono or system mono), 24px, thick border, hard shadow. These are the star of the controls.
  - Aspect ratio lock toggle between them: chain-link icon. Defaults to LOCKED. When locked, changing width auto-calculates height and vice versa. Toggle changes icon (linked → unlinked) with a satisfying click feel.
  - Inputs show pixel values. Typing updates the other field in real-time when locked.

- **PRESETS section**:
  - Chunky pill/tag buttons with thick borders. One row or wrapped grid.
  - Presets: "Instagram 1080×1080", "HD 1920×1080", "Twitter 1200×675", "OG Image 1200×630", "50%", "25%"
  - Tapping a preset fills Width/Height inputs and highlights the active preset.
  - Active preset: filled coral background. Inactive: white/cream background.

- **FORMAT section**:
  - Format dropdown: styled select with thick border. Options: JPG (default), PNG, WebP.
  - Quality slider: only visible for JPG and WebP. Range 1-100, default 85. Shows numeric value. Thick track, chunky thumb.
  - PNG selected: quality slider hides (PNG is always lossless).

**Download Bar (bottom, full width, sticky)**:
- Fixed to viewport bottom. Thick top border. White/cream background.
- Left side: file size comparison — "4.2 MB → ~380 KB (91% smaller)" in bold. This is the payoff moment.
- Right side: Large "DOWNLOAD" button. Coral (#FF6B6B) fill, thick black border, hard shadow. This is the most prominent button on the entire page.
- The tilde (~) on estimated size communicates it's an approximation before actual export.

---

### Screen 2b: Workspace State (Mobile, < 768px)

```
┌──────────────────────────┐
│ IMAGE RESIZER   [↻] [⚙] │
├──────────────────────────┤
│ ┌──────────────────────┐ │
│ │                      │ │
│ │    Image Preview     │ │
│ │                      │ │
│ └──────────────────────┘ │
│ photo.jpg                │
│ 4032×3024 · 4.2 MB      │
├──────────────────────────┤
│ DIMENSIONS               │
│ ┌──────────┐┌──────────┐│
│ │ W: 1920  ││ H: 1080  ││
│ └──────────┘└──────────┘│
│      [🔗 Locked]         │
│                          │
│ PRESETS                  │
│ [Insta] [HD] [Twitter]  │
│ [OG] [50%] [25%]        │
│                          │
│ FORMAT                   │
│ [JPG ▾]  Quality [85──] │
├──────────────────────────┤
│ 4.2→380KB  [ ⬇ DOWNLOAD]│
└──────────────────────────┘
```

**Mobile Adaptations**:
- Single-column stack: preview → controls → download bar.
- Preview area: full-width, constrained height (~40vh).
- Dimension inputs: side-by-side, slightly smaller (20px mono).
- Presets: horizontal scroll or 2×3 grid of smaller pills.
- Download bar: still fixed to bottom. Compressed layout.
- "New Image" abbreviated to icon-only (↻).
- Touch targets: all interactive elements minimum 44×44px.

---

### Advanced Settings Panel (gear icon → slide-in or popover)

```
┌──────────────────────────┐
│ ADVANCED          [✕]    │
├──────────────────────────┤
│                          │
│ Resize Method            │
│ ○ Pixels (default)       │
│ ○ Percentage             │
│                          │
│ Smoothing                │
│ [High ▾]                 │
│                          │
└──────────────────────────┘
```

**Note**: This panel is hidden by default. 90% of users will never open it. It exists for the 10% who want fine control. Minimal options to avoid scope creep.

---

## User Flows

### Primary Flow: Resize a Single Image

```
1. User lands on page
   → Sees full-viewport drop zone with "Drop your image here"
   → Privacy note visible at bottom

2. User drops/selects an image
   → Drop zone transitions to workspace (smooth, fast — no page load)
   → Image preview appears with original dimensions and file size
   → Width/Height auto-populated with original dimensions
   → Aspect ratio lock defaults to ON
   → Format defaults to original format (or JPG if unsupported output)
   → Download bar appears at bottom with estimated output size

3. User adjusts dimensions
   → Types new width → height auto-calculates (lock on)
   → OR taps a preset → both fields update
   → Estimated file size updates in download bar

4. User (optionally) changes format/quality
   → Selects format from dropdown
   → Adjusts quality slider if JPG/WebP

5. User downloads
   → Taps "DOWNLOAD" button
   → Browser saves file as "[original-name]-resized.[ext]"
   → Download bar shows brief success state: "Downloaded!"

6. User resizes another (optional)
   → Taps "New Image" → workspace resets to landing/drop zone
   → OR drops another image directly onto the workspace
```

### Secondary Flow: Use a Preset Size

```
1. User uploads image (same as primary flow steps 1-2)

2. User taps "Instagram 1080×1080" preset
   → Width input → 1080
   → Height input → 1080
   → Aspect ratio lock automatically UNLOCKS (preset overrides ratio)
   → Preview updates to show the new dimensions
   → Note displayed: "Image will be stretched" if aspect ratio differs significantly
     OR: the image is scaled to fit within 1080×1080 maintaining ratio

3. User downloads (same as primary flow step 5)
```

### Edge Flow: Large File Warning

```
1. User drops a 45MB image
   → Warning toast: "This is a large file. Processing may take a moment."
   → Processing begins with progress indicator
   → Workspace loads normally once processed

2. User drops a 60MB image (over 50MB limit)
   → Drop zone shows error: "File too large (60MB). Max size is 50MB."
   → Drop zone remains active for retry
```

### Edge Flow: Invalid File

```
1. User drops a .pdf or .txt file
   → Drop zone border flashes red
   → Message: "That doesn't look like an image. Try JPG, PNG, or WebP."
   → Drop zone remains active for retry
```

---

## Interaction Patterns

| Action | Element | Response | Timing |
|--------|---------|----------|--------|
| Drag over | Drop zone | Border solid + coral tint + "Let go!" copy | Instant |
| Drag leave | Drop zone | Revert to dashed border + original copy | Instant |
| Drop valid image | Drop zone | Transition to workspace | < 500ms |
| Drop invalid file | Drop zone | Red flash + error message | Instant, message fades after 3s |
| Click | Drop zone | Opens native file picker | Native |
| Type in Width | Width input | Height auto-updates (if locked) | Debounced 300ms |
| Type in Height | Height input | Width auto-updates (if locked) | Debounced 300ms |
| Click | Aspect lock toggle | Toggles lock on/off, icon changes | Instant |
| Click | Preset pill | Fills dimensions, highlights pill | Instant |
| Change | Format dropdown | Updates format, shows/hides quality slider | Instant |
| Drag | Quality slider | Updates quality value + estimated size | Debounced 200ms |
| Click | Download button | Processes + triggers browser download | < 1s for images under 10MP |
| Click | New Image button | Resets to landing state | Instant (revoke old blob URLs) |
| Hover | Any button | translate(-2px, -2px), shadow grows | 100ms ease |
| Click/Active | Any button | translate(2px, 2px), shadow shrinks to 0 | 50ms |

---

## States

### Empty State (Landing)
- **When**: Page load, after "New Image" reset
- **Show**: Full-viewport drop zone with copy and file icon
- **CTA**: "Drop your image here (or click to browse)"

### Loading / Processing State
- **When**: Image is being read from disk, or resize is being processed for download
- **Show**: Bold progress bar inside the download bar area. Thick border, coral fill, animated. Microcopy: "Crunching pixels..." or "Almost there..."
- **Duration**: Typically < 1 second. Only visible for large images.
- **Behavior**: UI controls remain visible but download button is disabled and shows "Processing..."

### Error State: Invalid File
- **When**: Non-image file dropped, or image fails to decode
- **Show**: Drop zone border turns red (#FF6B6B), error message appears centered in drop zone
- **Copy**: "That doesn't look like an image. Try JPG, PNG, or WebP."
- **Recovery**: Drop zone remains active. Error message fades after 3 seconds.

### Error State: File Too Large
- **When**: File exceeds 50MB
- **Show**: Same as invalid file, but with size-specific message
- **Copy**: "That file is too heavy ([X] MB). Max size is 50MB."
- **Recovery**: Same — zone stays active.

### Error State: Processing Failure
- **When**: Canvas API fails (corrupt image, exceeded canvas size limit)
- **Show**: Inline error in workspace area where preview would be
- **Copy**: "Something went wrong processing that image. Try a different file?"
- **Recovery**: "New Image" button to reset. Drop zone also accepts a new file.

### Success State: Download Complete
- **When**: File download triggered successfully
- **Show**: Download button briefly changes to "Downloaded!" with green background (#4CAF50) for 2 seconds, then reverts
- **Follow-up**: "Resize Another" text link appears next to or below the download bar

---

## Component Inventory

### 1. Drop Zone
- Dashed thick border (3px), cream background
- States: default, drag-hover, error
- Full-viewport on landing, re-accessible from workspace via "New Image"

### 2. Image Preview
- Fit-to-container with aspect ratio preserved
- Thick 2px border + hard shadow
- Checkerboard background for transparent images
- Below: filename, dimensions, file size (text only, never innerHTML)

### 3. Dimension Inputs
- Two large monospace inputs side by side
- Thick border, hard shadow
- Numeric only, pixel values
- Connected by aspect ratio lock toggle

### 4. Aspect Ratio Lock Toggle
- Chain-link icon (linked/unlinked)
- Positioned between width and height inputs
- Default: locked
- Chunky, clearly tappable (44×44px minimum)

### 5. Preset Pills
- Horizontal row of chunky bordered pills
- Active: coral fill, white text
- Inactive: white fill, dark text, thick border
- Hard shadow on all

### 6. Format Dropdown
- Styled native select with thick border
- Options: JPG, PNG, WebP
- Matches overall brutalist aesthetic

### 7. Quality Slider
- Visible only for JPG/WebP
- Thick track, chunky square thumb
- Numeric display alongside (e.g., "85")
- Range: 1-100

### 8. Download Bar
- Fixed to bottom of viewport
- Thick top border
- Contains: file size comparison (left) + download button (right)
- Download button: coral fill, largest button on page, bold text

### 9. "New Image" Button
- Secondary style: outlined, no fill
- Resets entire workspace
- Desktop: text + icon. Mobile: icon only.

### 10. Privacy Note
- Small, unobtrusive, always visible
- Landing: bottom of page
- Workspace: in header or footer
- Copy: "Your images never leave your browser."

---

## Responsive Breakpoints

| Breakpoint | Layout | Key Changes |
|-----------|--------|-------------|
| ≥ 1024px (Desktop) | Two-column: preview (60%) + sidebar (40%) | Full layout as wireframed |
| 768–1023px (Tablet) | Two-column with narrower sidebar | Presets wrap to 2 rows, smaller inputs |
| < 768px (Mobile) | Single-column stack | Preview → controls → sticky download bar |

---

## Accessibility Considerations

### Keyboard Navigation
- **Tab order**: Drop zone → (after upload) Width → Height → Lock toggle → Presets → Format → Quality → Download → New Image
- **Enter/Space** on drop zone: opens file picker
- **Arrow keys** on quality slider: adjust value
- **Escape**: close any open popover (advanced settings, privacy note)

### Screen Reader
- Drop zone: `aria-label="Image upload area. Drop an image or press Enter to browse files."`
- Dimension inputs: `aria-label="Width in pixels"` / `aria-label="Height in pixels"`
- Lock toggle: `aria-pressed="true/false"`, `aria-label="Aspect ratio lock"`
- Presets: each pill is a button with descriptive label (e.g., "Set dimensions to 1080 by 1080 for Instagram")
- Download button: announces file size comparison on focus
- Status updates (processing, error, success): use `aria-live="polite"` region

### Visual Accessibility
- **Color contrast**: All text meets WCAG AA 4.5:1 minimum against their backgrounds
  - #1A1A2E on #FEFAE0 (cream) = 14.6:1 — passes
  - #FFFFFF on #FF6B6B (coral button) = 3.7:1 — use #1A1A2E text on coral buttons for AA compliance
- **Focus indicators**: Thick 3px outline offset (2px) in coral (#FF6B6B) on all focusable elements. Brutalist aesthetic naturally supports visible focus states.
- **Touch targets**: All interactive elements minimum 44×44px
- **Motion**: Respect `prefers-reduced-motion` — disable hover translations and progress animations

### Brutalism Accessibility Advantages (from UI research)
- Thick borders = clear visual boundaries
- Bold typography = easier to read
- Solid colors = easy to maintain contrast
- Warm background = less eye strain than pure white

---

## Microcopy Guide

| Context | Copy | Tone |
|---------|------|------|
| Landing heading | "Drop your image here" | Direct, inviting |
| Landing subtext | "(or click to browse)" | Helpful, understated |
| Format hint | "JPG, PNG, WebP — up to 50MB" | Factual, brief |
| Privacy | "Your images never leave your browser." | Reassuring, confident |
| Drag hover | "Let go to resize!" | Encouraging |
| Processing | "Crunching pixels..." | Playful, brief |
| Download ready | "⬇ DOWNLOAD" | Bold, action-oriented |
| Success | "Downloaded!" | Celebratory, minimal |
| Error: bad file | "That doesn't look like an image. Try JPG, PNG, or WebP." | Friendly, helpful |
| Error: too large | "That file is too heavy (X MB). Max size is 50MB." | Empathetic, clear |
| Error: processing | "Something went wrong. Try a different file?" | Honest, recovery-focused |
| Lock on | "Ratio locked" | Terse, status |
| Lock off | "Ratio unlocked" | Terse, status |
| Reset action | "Resize Another" | Forward-looking |

---

## Design Tokens Reference (from UI research)

| Token | Value | Usage |
|-------|-------|-------|
| Background | #FEFAE0 | Page background |
| Surface | #FFFFFF | Cards, panels |
| Primary (CTA) | #FF6B6B | Download button, active preset |
| Secondary | #A8E6CF | Success states |
| Accent | #D4A5FF | Hover highlights (optional) |
| Text | #1A1A2E | All text, borders |
| Border weight | 2-3px | All interactive elements |
| Shadow | 4px 4px 0 #1A1A2E | Cards, buttons, inputs |
| Border radius | 4-6px | Buttons, cards, inputs |
| Heading font | Space Grotesk Bold | Titles, labels |
| Body font | Inter | Body text, descriptions |
| Mono font | JetBrains Mono / system | Dimension inputs, file sizes |

---

## Security-Informed UX Decisions

Per security research findings:
1. **File names displayed via `textContent` only** — never `innerHTML`. Prevents XSS via crafted filenames.
2. **File validation before processing** — check MIME type + size. Show friendly error for invalid files.
3. **Blob URL cleanup on "New Image"** — revoke old object URLs to prevent memory leaks.
4. **Download filename sanitization** — output as `[sanitized-original-name]-resized.[ext]`. Strip path traversal characters.
5. **No inline scripts** — all JS in external files to support strict CSP.
6. **EXIF stripping is a feature** — Canvas naturally strips metadata. Mention this in privacy note as a benefit.

---

Status: READY_FOR_REVIEW
