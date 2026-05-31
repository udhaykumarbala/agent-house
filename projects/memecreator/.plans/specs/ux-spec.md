# UX Specification: Meme Creator

Based on research in: .plans/research/ux-research.md, .plans/research/ui-research.md
Implementing features from: .plans/research/product-research.md
Security considerations from: .plans/research/security_expert.md
Architecture guidance from: .plans/research/architecture-research.md

---

## Screen List

1. **Empty State / Upload Screen** — Landing view where user uploads an image to begin
2. **Canvas Editor** — Main workspace with image, text overlays, and controls
3. **Export Confirmation** — Inline toast/notification confirming successful download

> Note: This is a single-page app. "Screens" refer to states of the same page, not separate routes.

---

## Wireframes

### Screen 1: Empty State / Upload Screen

```
┌──────────────────────────────────────────────────────────┐
│  [🎨 MemeForge]                           [? Shortcuts]  │
├──────────────────────────────────────────────────────────┤
│                                                          │
│                                                          │
│          ┌────────────────────────────────┐               │
│          │                                │               │
│          │     ┌──────────────────┐       │               │
│          │     │   📷 Icon        │       │               │
│          │     └──────────────────┘       │               │
│          │                                │               │
│          │   Drop an image here           │               │
│          │   or click to upload           │               │
│          │                                │               │
│          │   PNG, JPG, WebP, GIF          │               │
│          │   Max 10MB                     │               │
│          │                                │               │
│          │   [ Browse Files ]             │               │
│          │                                │               │
│          └────────────────────────────────┘               │
│          - - - - - - - - - - - - - - - - -               │
│          Or paste an image (Ctrl+V)                       │
│                                                          │
│                                                          │
└──────────────────────────────────────────────────────────┘
```

**Key Elements**:
- **Drop Zone**: Large dashed-border area, the hero element. Responds to drag hover (border animates, background changes). Covers ~60% of viewport.
- **File Type Hint**: Shows accepted formats (PNG, JPG, WebP, GIF) and 10MB max — sets expectations and aligns with security validation.
- **Browse Button**: Secondary CTA for users who prefer click-to-browse over drag-and-drop.
- **Paste Hint**: Subtle text below drop zone for power users. Supports `Ctrl+V` / `Cmd+V` clipboard paste.
- **No sign-up, no account prompt**: Page is immediately usable. Zero friction per product research.

**Design Notes** (from UI research):
- Dark background (`#1A1A2E`), drop zone uses lighter surface (`#252547`) with dashed border in muted lavender-gray.
- On drag hover: border becomes solid magenta (`#E91E8C`), subtle glow effect.
- DM Sans font for all UI text.

---

### Screen 2: Canvas Editor (Main Workspace)

```
┌──────────────────────────────────────────────────────────────────────┐
│  [🎨 MemeForge]          [↩ Undo] [↪ Redo]    [⬇ Download ▾]      │
├──────────┬───────────────────────────────────────┬───────────────────┤
│ TOOLBAR  │                                       │  STYLE PANEL      │
│          │                                       │                   │
│ [🖼 New  │     ┌───────────────────────────┐     │  (appears when    │
│  Image]  │     │                           │     │   text selected)  │
│          │     │    Uploaded Image          │     │                   │
│ [T Add   │     │                           │     │  Font Family      │
│  Text]   │     │   ┌─────────────┐         │     │  [Impact     ▾]  │
│          │     │   │ DRAG ME!    │←floating │     │                   │
│ [⭐ Meme │     │   └─────────────┘ text    │     │  Font Size        │
│  Style]  │     │                           │     │  [──●────] 48px   │
│          │     │                           │     │                   │
│          │     │        ┌─────────┐        │     │  Style            │
│          │     │        │ LOL     │        │     │  [B] [I] [U] [Aa] │
│          │     │        └─────────┘        │     │                   │
│          │     │                           │     │  Fill Color       │
│          │     └───────────────────────────┘     │  [■ #FFFFFF]      │
│          │                                       │                   │
│          │     ───── Canvas Area ─────           │  Stroke           │
│          │     (neutral dark bg with             │  [■ #000000] 3px  │
│          │      subtle grid pattern)             │  [──●────]        │
│          │                                       │                   │
│          │                                       │  Shadow           │
│          │                                       │  [toggle off]     │
│          │                                       │                   │
│          │                                       │  Alignment        │
│          │                                       │  [◁] [≡] [▷]     │
│          │                                       │                   │
│          │                                       │  ─────────────    │
│          │                                       │  TEXT LAYERS      │
│          │                                       │  ┌─────────────┐  │
│          │                                       │  │ "DRAG ME!" 🗑│  │
│          │                                       │  │ "LOL"       🗑│  │
│          │                                       │  └─────────────┘  │
├──────────┴───────────────────────────────────────┴───────────────────┤
│  Status: 2 text layers  │  Image: 1200×800  │  Zoom: Fit           │
└──────────────────────────────────────────────────────────────────────┘
```

**Key Elements**:

**Header Bar**:
- **Logo**: Top-left, links back to empty state (confirms "start over" via dialog).
- **Undo/Redo**: Icon buttons with keyboard shortcut hints on hover (`Ctrl+Z` / `Ctrl+Shift+Z`).
- **Download Button**: Primary accent color (magenta `#E91E8C`), always visible, top-right. Dropdown arrow reveals format options (PNG default, JPEG with quality slider). This is the #1 CTA — per research, it must never be hidden.

**Left Toolbar** (compact, icon-based):
- **New Image**: Replace current image (confirms if edits exist).
- **Add Text**: Creates a new text overlay at canvas center with default "Your text" placeholder. Primary action after upload.
- **Meme Style Preset**: One-click applies Impact font, white fill, black 3px stroke, all-caps — the classic meme look. Key differentiator identified in UX research.

**Canvas Area** (center, dominant):
- Image displayed at fit-to-container size, maintaining aspect ratio.
- Dark neutral background (`#1A1A2E`) surrounds the image, clearly delineating the edit zone.
- Text overlays rendered as draggable elements with visible selection handles.
- Selected text shows magenta border + corner drag handles for resize.
- Snap guides appear when dragging near center (horizontal/vertical) — per Canva pattern.
- Double-click on text to enter inline edit mode (cursor appears in text, type directly).

**Right Style Panel** (contextual — only visible when a text element is selected):
- **Font Family Dropdown**: Pre-loaded meme fonts — Impact (default), Anton, Bebas Neue, Bangers, Permanent Marker, Oswald, Arial Black, Comic Sans MS. Visual font name preview in dropdown.
- **Font Size Slider**: Range 16–120px with numeric input for precision. Slider is the primary control (per research: users think "bigger/smaller," not point sizes).
- **Style Toggles**: Bold, Italic, Uppercase, Strikethrough. Uppercase is prominent — meme text is often all-caps.
- **Fill Color**: Color picker (`<input type="color">`) with hex value display. Default: `#FFFFFF`.
- **Stroke/Outline**: Color picker + width slider (0–8px). Default: `#000000` at 3px. This is essential for the classic meme look.
- **Shadow**: Toggle on/off, with color, blur, and offset controls when expanded. Off by default.
- **Text Alignment**: Left / Center / Right icon buttons. Default: Center.
- **Text Layers List**: Shows all text overlays by content preview. Click to select, trash icon to delete. Drag to reorder z-index.

**Floating Toolbar** (appears near selected text on canvas):
```
┌─────────────────────────────────────────────┐
│ [Impact ▾] [48 ▾] [■ Color] [B] [I] [🗑]  │
└─────────────────────────────────────────────┘
```
- Compact version of the most-used controls from the style panel.
- Positioned above the selected text element (or below if near top edge).
- Reduces mouse travel for quick edits — per Adobe Express and Google Docs patterns.
- Font family, font size, color, bold, italic, delete.

**Status Bar** (bottom, subtle):
- Layer count, image dimensions, zoom level.
- Non-interactive, purely informational.

**Design Notes**:
- Left toolbar: Dark surface (`#16213E`), icon-only with tooltips. ~60px wide.
- Style panel: Dark surface (`#16213E`), ~280px wide. Slides in/out with animation.
- Canvas background: Slightly different from UI background to create depth.
- All controls minimum 44x44px touch targets.

---

### Screen 3: Export Confirmation (Toast)

```
┌──────────────────────────────────────────┐
│  ✓  Image downloaded!                     │
│     meme_20260316_143022.png              │
│                                    [OK]   │
└──────────────────────────────────────────┘
```

**Key Elements**:
- **Toast notification**: Appears top-center, auto-dismisses after 4 seconds, or user clicks "OK".
- Shows filename for reference.
- Canvas remains intact for further edits or re-export.
- Success state uses a green accent, not the primary magenta.

---

## User Flows

### Primary Flow: Create and Download a Meme

```
1. User opens app
   → Sees empty state with prominent upload zone

2. User drops/selects an image
   → Image loads onto canvas, editor UI appears
   → Left toolbar + header bar become active
   → Hint tooltip: "Click 'Add Text' or double-click the image"

3. User clicks "Add Text" (or double-clicks canvas)
   → New text overlay appears at canvas center: "Your text"
   → Text is auto-selected, inline edit mode active (cursor blinking)
   → Style panel slides in on the right
   → Floating toolbar appears near the text

4. User types their meme text
   → Text updates in real-time on canvas
   → Press Enter or click outside to confirm

5. User drags text to desired position
   → Text follows pointer, snap guides appear at center lines
   → Release to place

6. User styles text (optional)
   → Adjusts font, size, color, stroke via floating toolbar or style panel
   → All changes render instantly on canvas (real-time preview)

7. User clicks "Meme Style" preset (optional)
   → Selected text becomes: Impact, white, black 3px stroke, uppercase
   → One-click shortcut for the classic meme look

8. User adds more text (repeat steps 3-7)

9. User clicks "Download"
   → Image exports as PNG at original resolution
   → Browser download dialog or direct save
   → Success toast appears

10. User can continue editing or start over
```

### Secondary Flow: Quick Meme (Speed Path)

```
1. User drops image → 2. Clicks "Add Text" → 3. Types text
→ 4. Clicks "Meme Style" → 5. Clicks "Download"

Total interactions: 5 clicks + typing
Target time: Under 30 seconds
```

### Secondary Flow: Replace Image

```
1. User clicks "New Image" in left toolbar
2. If edits exist → confirmation dialog: "Replace image? Text overlays will be kept."
3. User selects new image
4. Canvas updates with new image, text overlays remain in place
5. User adjusts text positions if needed
```

### Secondary Flow: Undo a Mistake

```
1. User makes an unwanted change (delete text, move, style change)
2. User presses Ctrl+Z (or clicks Undo button)
3. Previous state is restored instantly
4. Redo available via Ctrl+Shift+Z
```

### Secondary Flow: Multiple Text Layers Management

```
1. User has 3+ text overlays on canvas
2. User clicks a text layer in the Layers List (right panel)
3. That text is selected on canvas, style panel shows its properties
4. User can reorder layers by dragging in the list
5. User can delete a layer via the trash icon (no confirmation for single text delete)
```

---

## Interaction Patterns

| Action | Element | Response |
|--------|---------|----------|
| Drop file | Upload zone | Image loads, transitions to editor state |
| Click | "Browse Files" button | File picker opens, filtered to images |
| Ctrl+V | Anywhere (empty state) | Paste image from clipboard |
| Click | "Add Text" button | New text overlay at canvas center, enters edit mode |
| Double-click | Canvas (no text selected) | New text overlay at click position |
| Double-click | Existing text overlay | Enter inline text edit mode |
| Single-click | Text overlay | Select it (shows handles + floating toolbar + style panel) |
| Click | Empty canvas area | Deselect current text |
| Drag | Selected text overlay | Move to new position, snap guides at center |
| Drag | Corner handle of text | Resize text (font size scales proportionally) |
| Click | "Meme Style" button | Apply preset: Impact, white, black stroke, uppercase |
| Change | Any style control | Instant real-time update on canvas |
| Click | "Download" button | Export PNG at full resolution, trigger download |
| Click | Download dropdown arrow | Show format options (PNG / JPEG + quality) |
| Ctrl+Z | Anywhere | Undo last action |
| Ctrl+Shift+Z | Anywhere | Redo |
| Delete/Backspace | With text selected (not editing) | Delete text overlay |
| Escape | During text edit | Exit edit mode, keep changes |
| Click | Trash icon in layers list | Delete that text overlay |
| Drag | Layer in layers list | Reorder z-index |
| Click | "New Image" | Replace image (confirm if edits exist) |

---

## States

### Empty State (No Image)
**When**: App first loads, or after user clears the canvas.
**Show**: Full-screen upload drop zone with icon, instructional text, and browse button.
**Action**: "Drop an image here or click to upload" — this is the only possible action, keeping focus clear.
**Design**: Dark background, drop zone is a centered card with dashed border. Subtle animation (icon gently pulses) to draw attention.

### Image Loaded, No Text
**When**: Image uploaded but no text overlays added yet.
**Show**: Image on canvas, active toolbar with "Add Text" button highlighted/pulsing subtly.
**Action**: Tooltip hint appears once: "Click 'Add Text' to start, or double-click the image."
**Design**: Style panel is hidden (no text selected). Canvas is the hero.

### Text Selected
**When**: User clicks on a text overlay.
**Show**: Selection handles (magenta border + corner resize handles), floating toolbar above text, style panel slides in on right.
**Action**: All style controls are active and reflect the selected text's current properties.

### Text Editing (Inline)
**When**: User double-clicks a text overlay.
**Show**: Text becomes editable with a blinking cursor. Text has a subtle text-input background.
**Action**: Typing updates text in real-time. Enter confirms. Escape cancels (reverts to pre-edit text).

### Dragging Text
**When**: User is actively dragging a text overlay.
**Show**: Text follows pointer. Snap guides appear (dashed lines) when near horizontal/vertical center. Cursor changes to grab/move.
**Action**: Release to place. Floating toolbar hides during drag to reduce visual noise.

### Loading State
**When**: Image is being processed after upload (large files may take a moment).
**Show**: Skeleton loader over the canvas area. Upload zone shows a progress indicator.
**Duration**: Typically < 1 second. If > 2 seconds, show a "Processing image..." message.

### Error State: Invalid File
**When**: User uploads an unsupported file type or file exceeds limits.
**Show**: Inline error message within the upload zone (red accent, icon + text).
**Message Examples**:
- "This file type isn't supported. Please use PNG, JPG, WebP, or GIF."
- "This image is too large (max 10MB). Try a smaller file."
- "This image couldn't be loaded. It may be corrupted."
**Action**: Upload zone remains active — user can try again immediately. No page reload needed.

### Error State: Export Failure
**When**: Canvas export fails (rare, usually browser memory issue with very large images).
**Show**: Error toast at top-center: "Export failed. Try a smaller image or different format."
**Action**: Editor remains intact. User can try JPEG (smaller) or reduce quality.

### Success State (Export Complete)
**When**: Download triggered successfully.
**Show**: Success toast (green accent) with filename, auto-dismisses in 4 seconds.
**Action**: Canvas stays intact. User can re-export, continue editing, or start over.

---

## Accessibility Considerations

### Keyboard Navigation
- **Tab order**: Upload zone → Header buttons (Undo, Redo, Download) → Left toolbar (New Image, Add Text, Meme Style) → Canvas (cycles through text overlays) → Style panel controls.
- **Arrow keys**: When a text overlay is focused (not in edit mode), arrow keys nudge position by 1px. Shift+Arrow nudges by 10px. Provides precision without mouse.
- **Enter**: On a focused text overlay, enters inline edit mode.
- **Escape**: Exits edit mode or deselects current text.
- **Delete/Backspace**: Deletes focused text overlay (when not in edit mode).
- **Tab between text overlays**: When canvas is focused, Tab cycles through text overlays in z-order.

### Screen Reader Support
- **Upload zone**: `role="button"` with `aria-label="Upload an image. Drag and drop or click to browse. Accepted formats: PNG, JPG, WebP, GIF. Maximum size 10 megabytes."`
- **Canvas area**: `aria-label="Meme editor canvas"` with `aria-live="polite"` region announcing state changes (e.g., "Text added", "Text deleted", "Image uploaded").
- **Text overlays**: Each announced as "Text overlay: [content], position [x, y]" when focused.
- **Style controls**: Standard form labels. Slider values announced on change.
- **Download button**: `aria-label="Download meme as PNG image"`.
- **Toasts**: Use `role="status"` and `aria-live="polite"` for announcements.

### Visual Accessibility
- **Touch targets**: All interactive elements minimum 44x44px (per WCAG 2.5.5).
- **Color contrast**: All text meets WCAG AA (4.5:1 for normal text, 3:1 for large text). Primary text (`#F8F8F2`) on dark backgrounds (`#1A1A2E`) = 13.5:1 ratio. Secondary text (`#A0A0B8`) on dark = 5.2:1 ratio.
- **Focus indicators**: Visible focus ring (2px magenta outline with 2px offset) on all interactive elements. Never rely on color alone.
- **Motion**: Respect `prefers-reduced-motion` — disable pulse animations, snap guide transitions, and toolbar slide-in when enabled.
- **Icons**: All icon-only buttons have tooltips and `aria-label`. No icon used without accompanying text for screen readers.

### Color Blindness Considerations
- Error states use icon + text label, not just red color.
- Success states use icon + text label, not just green color.
- Selected state uses border + handles, not just color highlight.

---

## Responsive Behavior

### Desktop (≥1024px)
- Full three-column layout: left toolbar + canvas + right style panel.
- Style panel slides in/out. Canvas resizes to fill remaining space.

### Tablet (768px–1023px)
- Left toolbar collapses to a floating action bar at bottom-center.
- Style panel becomes a bottom sheet (slides up from bottom, covers ~40% of screen).
- Canvas takes full width.

### Mobile (< 768px)
- Upload zone fills viewport.
- Canvas fills width, controls in a bottom sheet.
- Floating toolbar is the primary style control (style panel available via expanding bottom sheet).
- Larger drag handles (56x56px) for touch accuracy.
- Pinch-to-zoom on canvas for precise placement.

> Note: Per product research, mobile is a P1 enhancement, not MVP. Desktop-first design. Mobile layout is defined here for future implementation.

---

## Design Tokens Reference

These align with UI research decisions for developer handoff:

| Token | Value | Usage |
|-------|-------|-------|
| `--bg-primary` | `#1A1A2E` | Page background, canvas surround |
| `--bg-surface` | `#16213E` | Panels, cards, toolbar bg |
| `--bg-surface-hover` | `#252547` | Hover states on surfaces |
| `--accent-primary` | `#E91E8C` | CTA buttons, selection highlights, focus rings |
| `--accent-primary-hover` | `#FF2DA5` | Hover state for accent |
| `--text-primary` | `#F8F8F2` | Primary UI text |
| `--text-secondary` | `#A0A0B8` | Secondary/hint text |
| `--success` | `#22C55E` | Success toasts, confirmations |
| `--error` | `#EF4444` | Error messages, invalid states |
| `--warning` | `#F59E0B` | Warning messages |
| `--border-default` | `#2D2D4A` | Subtle borders between panels |
| `--border-focus` | `#E91E8C` | Focus ring color |
| `--radius-sm` | `6px` | Buttons, inputs |
| `--radius-md` | `8px` | Cards, panels |
| `--radius-lg` | `12px` | Modals, large containers |
| `--font-ui` | `'DM Sans', sans-serif` | All UI text |
| `--font-meme-default` | `'Impact', sans-serif` | Default meme text |

---

## Meme Style Presets

Quick-apply text styles — one click to apply a complete look:

| Preset Name | Font | Size | Fill | Stroke | Stroke Width | Uppercase | Shadow |
|-------------|------|------|------|--------|-------------|-----------|--------|
| Classic Meme | Impact | 48px | `#FFFFFF` | `#000000` | 3px | Yes | No |
| Modern Clean | Bebas Neue | 42px | `#FFFFFF` | none | 0 | No | Yes (soft black) |
| Handwritten | Permanent Marker | 36px | `#FFFFFF` | `#000000` | 2px | No | No |
| Bold Impact | Anton | 56px | `#FFD700` | `#000000` | 4px | Yes | No |
| Neon Glow | Bangers | 44px | `#00FF88` | `#003322` | 2px | No | Yes (green glow) |
| Subtle Caption | DM Sans | 24px | `#FFFFFF` | none | 0 | No | Yes (soft black) |

The "Meme Style" button in the toolbar applies "Classic Meme" by default. Users can access other presets via a dropdown.

---

## Keyboard Shortcuts Reference

| Shortcut | Action |
|----------|--------|
| `Ctrl/Cmd + Z` | Undo |
| `Ctrl/Cmd + Shift + Z` | Redo |
| `Ctrl/Cmd + S` | Download (override browser save) |
| `Ctrl/Cmd + V` | Paste image (empty state) |
| `T` | Add new text (when not editing text) |
| `Delete` / `Backspace` | Delete selected text (when not editing) |
| `Escape` | Deselect / exit edit mode |
| `Enter` | Enter edit mode on selected text |
| `Arrow keys` | Nudge selected text 1px |
| `Shift + Arrow` | Nudge selected text 10px |

---

Status: READY_FOR_REVIEW
