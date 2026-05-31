# UX Research: Meme Creator

## Pattern Analysis

Referenced apps and their successful patterns:

### Imgflip (imgflip.com)
- **Pattern Used**: Template gallery as entry point, then a simple two-panel editor — image preview on the left, text controls on the right. Pre-populated top/bottom text fields match the classic meme format. One-click download.
- **Why It Works**: Zero learning curve. Users see familiar meme templates, click one, type text, and download. The preview updates in real-time so there's no "generate" step. The layout follows a scan-left-read-right pattern that matches natural reading flow.
- **What We Can Learn**: Real-time preview is non-negotiable. Users expect to see changes instantly on the canvas. A template gallery dramatically lowers the barrier to entry — but we should also make "upload your own image" equally prominent for custom memes.

### Canva (canva.com)
- **Pattern Used**: Full canvas editor with a left sidebar for tools/assets, a central canvas with drag-and-drop elements, and a floating toolbar that appears contextually when an element is selected. Text is added by clicking a "Text" tool, then clicking on the canvas.
- **Why It Works**: The contextual toolbar keeps the UI clean — controls only appear when relevant. Drag-to-position is intuitive and gives precise control. The canvas metaphor (infinite workspace, zoom/pan) is familiar from design tools. Undo/redo and keyboard shortcuts serve power users without cluttering the UI.
- **What We Can Learn**: Contextual toolbars (appearing on text selection) keep the interface uncluttered. Drag-and-drop positioning is the expected interaction for placing text on an image. We should show formatting controls only when a text element is selected.

### Kapwing (kapwing.com)
- **Pattern Used**: Browser-based editor with a prominent canvas area, timeline at the bottom (for video, but the image editor uses a similar layered approach), and a right-side panel for element properties. Uses a "click to add text" model where text appears as a draggable box on the canvas.
- **Why It Works**: The property panel on the right provides detailed control (font, size, color, outline, shadow) without overwhelming the canvas. The layer system lets users manage multiple text overlays clearly. Export options are prominent in the top-right corner.
- **What We Can Learn**: A property/styling panel is essential for text customization. Showing a layer list helps when users have multiple text overlays. Export should be a primary action button, always visible.

### Adobe Express (express.adobe.com)
- **Pattern Used**: Template-first approach with a clean editor. Uses inline text editing directly on the canvas (double-click to edit). Floating formatting bar appears above selected text. Color picker, font selector, and effects are in a compact toolbar.
- **Why It Works**: Inline editing (typing directly on the canvas where text appears) feels natural — it's WYSIWYG in the truest sense. The floating bar doesn't obscure the canvas. Pre-built text styles ("Add a heading", "Add a subheading") give users starting points.
- **What We Can Learn**: Inline text editing on the canvas is more intuitive than editing in a separate panel. Pre-styled text options (e.g., "Impact white with black outline" — the classic meme style) save users time.

### Mematic (mobile app, reference for interaction patterns)
- **Pattern Used**: Minimal UI — image fills the screen, text is added by tapping, and a bottom sheet slides up for styling options. Pinch-to-resize and drag-to-move for text elements. Large, thumb-friendly controls.
- **Why It Works**: The image is always the hero — UI never competes with content. Touch-first interactions (drag, pinch) are translated well to mouse (drag, scroll-to-resize) in browser. The bottom sheet pattern keeps tools accessible without blocking the canvas.
- **What We Can Learn**: Keep the image/canvas as the dominant visual element. Ensure text manipulation (move, resize) works with both mouse and touch. A bottom panel or sidebar for controls keeps the canvas unobstructed.

## User Journey Analysis

### Entry Point
User lands on the single-page app. They see two clear paths:
1. **Upload an image** — prominent upload area (drag-and-drop + click-to-browse)
2. Optionally, a small set of popular blank templates (if we choose to include them)

The upload area should be the hero element on the landing/empty state. No sign-up walls, no distractions.

### Core Loop
1. Upload or select an image
2. Add text overlay (click "Add Text" button or double-click on canvas)
3. Position text by dragging it on the canvas
4. Style text (font, size, color, outline/stroke, shadow)
5. Repeat steps 2-4 for additional text elements
6. Export/download the final image

This is a **create → tweak → export** loop. Most users will go through it once per meme. Speed matters — the entire flow should take under 60 seconds.

### Success State
User clicks "Download" and gets a high-quality image file (PNG or JPEG) saved to their device. A brief success toast ("Image downloaded!") confirms the action. The canvas remains intact so they can make further edits or download again.

## Mental Model

Users think of this task as: **"Putting stickers/labels on a photo"** — they have an image and want to stamp text on top of it in specific spots. It's closer to a collage/sticker metaphor than a document editing metaphor.

Common terminology users use:
- "Add text" (not "insert text element")
- "Move it here" (drag and drop)
- "Make it bigger/smaller" (resize)
- "Change the font/color"
- "Download" or "Save" (not "export" — though "export" is acceptable as a secondary label)
- "Meme text" (white Impact font with black outline)
- "Caption"

## Anti-Patterns to Avoid

- **Don't**: Require account creation or login before using the tool — **Why**: Meme creation is an impulse activity. Any friction before the editor kills conversion. Imgflip and Kapwing both let users create without signing in.

- **Don't**: Use a multi-step wizard or modal flow for adding/styling text — **Why**: Modals break the direct manipulation mental model. Users want to see their changes on the canvas in real-time, not configure in a dialog and then "apply." Every successful competitor uses direct canvas manipulation.

- **Don't**: Auto-position text only at top/bottom — **Why**: While classic memes use top/bottom text, modern memes place text anywhere. Restricting placement feels limiting. Imgflip's biggest user complaint is the rigid top/bottom default.

- **Don't**: Hide the download/export button behind menus or secondary screens — **Why**: Download is the primary goal of every session. It should be a persistent, prominent button (top-right corner or floating action). Kapwing's export flow has too many steps and users complain about it.

- **Don't**: Use tiny touch targets or cramped controls — **Why**: Many users will access this on tablets or touch-enabled laptops. Controls smaller than 44x44px cause frustration and mis-taps.

- **Don't**: Show a blank canvas with no guidance — **Why**: An empty canvas with no affordances is intimidating. Always show either the uploaded image or a clear upload prompt. After image upload, a hint like "Click to add text" guides the first interaction.

- **Don't**: Require manual font-size entry via number input only — **Why**: Users think in visual terms ("bigger", "smaller"), not point sizes. Provide a slider or drag handles for intuitive resizing, with a number input as a secondary precision option.

## Recommended Patterns for Our App

Based on research, we should use:

1. **Single-page canvas editor with sidebar controls** — The canvas (image + overlays) is the central, dominant element. A collapsible sidebar or panel on the right holds text styling controls. This mirrors Canva/Kapwing's proven layout and keeps the focus on the image. The sidebar only shows controls when a text element is selected.

2. **Direct manipulation for text placement** — Text elements are added to the canvas and positioned via drag-and-drop. Resize via corner drag handles. This matches every successful competitor and the user's mental model of "placing stickers." No coordinate inputs, no alignment dialogs.

3. **Contextual floating toolbar for quick styling** — When a text element is selected, a compact floating toolbar appears near it (above or below) with the most-used controls: font family, font size, text color, and bold/italic toggles. Advanced options (outline, shadow, opacity) live in the sidebar. This reduces mouse travel and follows Adobe Express / Google Docs patterns.

4. **"Meme preset" one-click style** — A single button that applies the classic meme look (Impact font, white fill, black outline, centered, all-caps). This covers the most common use case in one click and signals to the user that we understand their intent.

5. **Prominent, always-visible download button** — Top-right corner, primary color, labeled "Download." One click starts the export. Optional format selection (PNG/JPEG) as a dropdown on the button, defaulting to PNG.

6. **Drag-and-drop image upload with clear empty state** — The initial view shows a large drop zone with an illustration/icon, "Drag an image here or click to upload" text, and accepted format hints. This is the universal pattern across all competitors.

7. **Undo/Redo support** — Ctrl+Z / Ctrl+Y (Cmd+Z / Cmd+Y on Mac) for undo/redo. This is a basic expectation in any editor and provides a safety net that encourages experimentation.

## Technology Considerations (for handoff to engineering)

- **Canvas rendering**: HTML5 Canvas API for the final composite (needed for pixel-perfect export). However, for the interactive editor, using DOM elements (divs) overlaid on the image provides easier drag/resize/text-editing interactions. On export, composite everything onto a canvas for download.
- **Export method**: `canvas.toBlob()` or `canvas.toDataURL()` to generate downloadable PNG/JPEG. Use a hidden canvas sized to the original image dimensions for full-resolution export.
- **Touch support**: Pointer events API (not just mouse events) to support touch, pen, and mouse uniformly.
- **No backend required**: Entire app runs client-side. No uploads to a server, no storage, no auth. This keeps it fast and private.

---
Status: READY_FOR_PLANNING
