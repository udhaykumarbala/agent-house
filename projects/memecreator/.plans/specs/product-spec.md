# Product Specification: Meme Creator

Based on research in: .plans/research/product-research.md, .plans/research/ux-research.md, .plans/research/ui-research.md, .plans/research/architecture-research.md, .plans/research/security_expert.md

## Target User

**"Quick-Meme Quinn"** — 18–35, digitally native, active on Reddit/Twitter/Discord/Instagram. Mix of students, office workers, and casual content creators. Primarily desktop browser, occasionally mobile. Wants to make a meme in under 60 seconds without sign-ups, watermarks, or bloated design tools. Privacy-conscious — doesn't want images uploaded to external servers.

## Unique Value Proposition

**"Drop an image, drag your text anywhere, export your meme — no signup, no watermark, no upload."**

A zero-friction, canvas-based meme editor that sits in the sweet spot between basic meme generators (Imgflip) and full design tools (Canva/Kapwing). 100% client-side — images never leave the browser.

## Features (Prioritized)

### P0 - Must Have (MVP)

- [ ] **F1: Image Upload** — Drag-and-drop or file picker to load an image onto the canvas. Supports JPEG, PNG, WebP, GIF. Max 10MB, max 4096x4096px. Validates file type via MIME, extension, and magic bytes. Blocks SVG uploads.
- [ ] **F2: Canvas Editor** — Central canvas displays the uploaded image at a responsive display size while preserving original resolution for export. Dark workspace background (#1A1A2E) keeps focus on the image.
- [ ] **F3: Add Text Overlay** — "Add Text" button creates a new text element on the canvas with default meme styling (Impact, white, black outline). Multiple text overlays supported.
- [ ] **F4: Drag-to-Position Text** — Click and drag text elements anywhere on the canvas using pointer events (unified mouse + touch). Visual feedback on selection (accent-colored border).
- [ ] **F5: Text Styling Controls** — Right sidebar panel (visible when text is selected) with: font family (Impact, Anton, Bebas Neue, Bangers, Permanent Marker, Oswald, Arial), font size (slider 16–120px), text color picker, outline/stroke color + width, bold/italic toggles, text alignment (left/center/right).
- [ ] **F6: Meme Preset Style** — One-click "Classic Meme" button that applies Impact font, white fill, black outline, all-caps. Instant classic meme look without manual styling.
- [ ] **F7: Export/Download** — Prominent "Download" button (always visible, top-right). Exports at original image resolution. Default PNG. Sanitized filename (`meme_[timestamp].png`). Uses canvas.toBlob() + blob URL + programmatic `<a>` click. Revokes blob URL after download.
- [ ] **F8: Delete Text** — Remove individual text overlays via a delete button in the text list or a keyboard shortcut (Delete/Backspace when text selected).
- [ ] **F9: Empty State & Upload UX** — Large, inviting drop zone on initial load with "Drag an image here or click to upload" messaging. Clear accepted format hints. No blank canvas confusion.
- [ ] **F10: Secure Text Rendering** — All DOM text uses `textContent`, never `innerHTML`. Canvas `fillText()`/`strokeText()` for rendering. No eval or document.write with user input.

### P1 - Should Have

- [ ] **F11: Undo/Redo** — Ctrl+Z / Ctrl+Shift+Z (Cmd on Mac). Snapshot-based history (max 20-30 states). Safety net that encourages experimentation.
- [ ] **F12: Inline Text Editing** — Double-click a text element on canvas to edit text in-place rather than only via a separate input field.
- [ ] **F13: JPEG Export Option** — Toggle between PNG and JPEG export. JPEG quality slider (0.7–1.0). Smaller file size option for photo-based memes.
- [ ] **F14: Text Style Presets** — Quick-apply styles beyond classic meme: "Modern Clean" (DM Sans, white, subtle shadow), "Handwritten" (Permanent Marker), "Neon Glow" (Bangers, colored shadow).
- [ ] **F15: Clipboard Paste Upload** — Ctrl+V / Cmd+V to paste an image from clipboard directly onto the canvas. Validates pasted content type.
- [ ] **F16: Text Shadow** — Shadow color, blur, and offset controls in the styling panel for text depth effects.
- [ ] **F17: Responsive Mobile Layout** — Controls panel stacks below canvas on small screens. Larger touch targets (min 44x44px). Collapsible controls.

### P2 - Nice to Have (Future)

- [ ] **F18: Text Rotation** — Rotate text elements via a drag handle or rotation control.
- [ ] **F19: All-Caps Toggle** — Button to toggle uppercase on text overlays.
- [ ] **F20: Text Background** — Optional filled rectangle behind text bounding box for readability.
- [ ] **F21: Opacity Control** — Per-text-element opacity slider.
- [ ] **F22: Snap-to-Center Guides** — Visual alignment guides when dragging text near canvas center or edges.
- [ ] **F23: Keyboard Accessibility** — Tab between text overlays, arrow keys for fine positioning, screen reader announcements via aria-live.

## User Stories

1. As Quick-Meme Quinn, I want to **drag an image onto the page and immediately see it on the canvas** so that I can start creating a meme in seconds without navigating menus or signing up.

2. As Quick-Meme Quinn, I want to **add text and drag it to any position on the image** so that I'm not limited to top/bottom placement like classic meme generators.

3. As Quick-Meme Quinn, I want to **apply the classic meme look (white Impact text with black outline) in one click** so that I can make a recognizable meme without manually configuring each style property.

4. As Quick-Meme Quinn, I want to **customize font, size, color, and outline of each text overlay independently** so that I can create memes with different visual styles and multiple text layers.

5. As Quick-Meme Quinn, I want to **click "Download" and instantly get a high-quality PNG** so that I can share the meme on social media right away without watermarks or quality loss.

6. As Quick-Meme Quinn, I want to **see all my changes in real-time on the canvas** so that I can iterate visually without a "generate" or "preview" step.

7. As Quick-Meme Quinn, I want to **know my images never leave my browser** so that I can use the tool with workplace screenshots or personal photos without privacy concerns.

8. As Quick-Meme Quinn, I want to **undo my changes with Ctrl+Z** so that I can experiment with text placement and styling without fear of making irreversible mistakes.

## Acceptance Criteria

### F1: Image Upload
- [ ] Drag-and-drop a JPEG/PNG/WebP/GIF file onto the drop zone loads it onto the canvas
- [ ] Clicking the drop zone opens a native file picker filtered to image types
- [ ] Files over 10MB show an error message and are rejected
- [ ] Images over 4096x4096 are rejected with a clear error
- [ ] SVG files are blocked with a message explaining accepted formats
- [ ] Invalid files (non-image, corrupted) show an error and don't break the app
- [ ] File validation checks MIME type, extension, and magic bytes

### F2: Canvas Editor
- [ ] Uploaded image displays at a responsive size that fits the viewport
- [ ] Image maintains aspect ratio — never stretched or distorted
- [ ] Canvas sits center-stage with dark workspace background
- [ ] Original image dimensions are preserved in memory for export

### F3: Add Text Overlay
- [ ] Clicking "Add Text" creates a new text element on the canvas
- [ ] Default text reads "Your text here" in Impact, white fill, black outline
- [ ] Multiple text elements can coexist on the canvas
- [ ] Each text element has a unique ID and independent styling

### F4: Drag-to-Position Text
- [ ] Clicking a text element selects it (visual accent border appears)
- [ ] Dragging a selected text element repositions it on the canvas in real-time
- [ ] Text elements can be placed anywhere within the canvas bounds
- [ ] Clicking empty canvas area deselects all text elements
- [ ] Works with both mouse and touch via pointer events

### F5: Text Styling Controls
- [ ] Styling panel appears/becomes active when a text element is selected
- [ ] Font family dropdown lists all available fonts with visual preview
- [ ] Font size slider adjusts text size from 16px to 120px in real-time
- [ ] Color picker changes text fill color immediately on canvas
- [ ] Outline color and width controls update stroke in real-time
- [ ] Bold and italic toggles work correctly
- [ ] Text alignment buttons (left/center/right) update text positioning

### F6: Meme Preset Style
- [ ] "Classic Meme" button applies: Impact font, white fill, black 3px outline, all-caps
- [ ] Preset applies to the currently selected text element
- [ ] Visual style updates immediately on canvas

### F7: Export/Download
- [ ] "Download" button is always visible in the top-right area
- [ ] Click triggers an immediate PNG download — no modal, no extra steps
- [ ] Exported image is at original upload resolution, not display size
- [ ] All text overlays are composited correctly at export resolution
- [ ] Text positioning and styling are accurate in the exported image
- [ ] Exported PNG does not contain EXIF metadata from the original image
- [ ] Download filename follows pattern: `meme_[timestamp].png`
- [ ] Blob URL is revoked after download to free memory

### F8: Delete Text
- [ ] Selected text element can be deleted via a delete/trash button
- [ ] Pressing Delete or Backspace key removes the selected text
- [ ] Deleting a text element updates the canvas immediately
- [ ] Text list in sidebar updates to reflect deletion

### F9: Empty State
- [ ] Initial page load shows a large, clear upload drop zone
- [ ] Drop zone includes an icon, "Drag an image here or click to upload" text, and format hints
- [ ] No confusing blank canvas — the upload prompt is the hero element
- [ ] After image upload, drop zone is replaced by the canvas editor

## Success Metrics

- **Time-to-first-meme**: User can go from page load to downloaded meme in under 60 seconds (measure via session timing)
- **Upload success rate**: >95% of upload attempts result in a successfully loaded canvas (track error rates)
- **Export completion rate**: >90% of sessions with an uploaded image result in at least one download
- **Return usage**: Users return to create another meme within 30 days (cookie/local storage tracking)
- **Page load time**: Initial load under 2 seconds on 3G connection (Lighthouse performance score >90)
- **Zero server dependency**: 100% of features work without any network requests after initial page load

## Out of Scope (for MVP)

- User accounts, login, or registration
- Server-side image processing or storage
- Meme template gallery / pre-built templates
- Image-to-image sharing or social sharing integration
- Video or GIF creation
- Multi-image collage or layered images
- AI-powered meme text suggestions
- Text rotation (deferred to P2)
- Mobile-native app (web-only)
- Analytics or tracking beyond basic metrics
- Monetization features (ads, premium tier, watermarks)
- Image cropping, filters, or non-text image manipulation
- Collaborative editing

## Technical Constraints (from Architecture & Security Research)

- **Architecture**: Static SPA — Vite + Tailwind + TypeScript, no framework (vanilla TS)
- **Rendering**: Raw Canvas API (no Fabric.js) — keeps bundle small, gives full control
- **Interaction**: Pointer Events API for unified mouse/touch handling
- **Fonts**: Google Fonts loaded via `<link>` with `font-display: swap`
- **Security**: CSP headers on deployment, `textContent` only for DOM text, file validation chain (MIME + extension + magic bytes + dimensions), SVG blocked
- **Privacy**: 100% client-side — no server uploads, canvas re-encoding strips EXIF metadata
- **Browser support**: Modern browsers (Chrome, Firefox, Safari, Edge). No IE11.

---
Status: READY_FOR_REVIEW
