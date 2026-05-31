# UX Research: Image Resizer (Minimal Friendly Brutalism)

## Pattern Analysis

Referenced apps and their successful patterns:

### Squoosh (squoosh.app) - Google Chrome Labs
- **Pattern Used**: Single-page tool with a massive drag-and-drop zone as the landing state. Before/after slider for comparing original vs. compressed output. Settings panel as a collapsible sidebar.
- **Why It Works**: Zero onboarding friction — the entire page IS the drop zone. Users see results immediately with a visual diff slider. No accounts, no sign-ups, no distractions.
- **What We Can Learn**: The "whole page is the action" pattern is ideal. Drop an image → instantly see controls. Minimal chrome, maximum workspace. The before/after slider gives users confidence in what they're getting.

### TinyPNG (tinypng.com)
- **Pattern Used**: Centered drop zone with playful panda branding. Batch upload with a simple list showing compression results per file. Download all as ZIP.
- **Why It Works**: Friendly personality (the panda) makes a utilitarian tool feel approachable. The progress list with file sizes and percentage savings gives clear, satisfying feedback. Batch support respects power users.
- **What We Can Learn**: A single mascot or visual motif adds personality without complexity. Showing "before size → after size (X% smaller)" is the key feedback metric users want. Batch download as ZIP is expected.

### iLoveIMG (iloveimg.com)
- **Pattern Used**: Tool selector grid on homepage → dedicated tool page per function. Each tool page has upload area + settings + download. Breadcrumb-like step flow: Upload → Adjust → Download.
- **Why It Works**: Clear separation of tools prevents feature overload. The 3-step linear flow (upload → configure → download) matches the user's mental model perfectly.
- **What We Can Learn**: A linear 3-step flow is the clearest mental model for image tools. However, their multi-page approach adds unnecessary navigation — we should keep it single-page.

### Photopea (photopea.com)
- **Pattern Used**: Full Photoshop-like interface in the browser. Menu bars, floating panels, layer system.
- **Why It Works**: For power users who need a full editor, it's remarkable. But for a simple resize task, it's massive overkill.
- **What We Can Learn**: Anti-example for our use case. A resizer should NOT look like an editor. Resist feature creep. One job, done well.

### Brutalist Web Design References (brutalistwebsites.com, Poolsuite, Craigslist)
- **Pattern Used**: Thick borders, monospace or bold sans-serif typography, high-contrast colors, raw/exposed UI elements, minimal decoration, visible grid structure, chunky interactive elements.
- **Why It Works**: Brutalism strips away visual noise. When done with "friendly" warmth (soft colors, rounded corners, playful copy), it creates a memorable, confident aesthetic that feels both modern and honest.
- **What We Can Learn**: Use thick 2-4px black borders, bold type hierarchy, a warm accent color palette (not sterile gray), chunky buttons with visible hover states, and generous whitespace. "Friendly brutalism" = brutalist structure + warm colors + playful microcopy.

## User Journey Analysis

### Entry Point
User has an image that's too large (for email, website, social media, or storage). They search "resize image online" or "image resizer." They want the fastest path from "I have a big image" to "I have a correctly sized image."

### Core Loop
1. Drop/select image(s)
2. Set target dimensions (or pick a preset)
3. Preview the result
4. Download

This loop repeats for each batch of images. Many users resize one image at a time. Power users resize batches.

### Success State
- The resized image is downloaded to their device
- They can see the new dimensions and file size
- The quality looks acceptable
- Total time: under 30 seconds for a single image

## Mental Model

Users think of this task as: **"Shrink my picture"** — like squeezing a photo to fit into a frame. They think in terms of:
- "Make it smaller" (reduce dimensions)
- "Make it fit [X]" (social media dimensions, email limits)
- "Keep it looking good" (quality preservation)

Common terminology:
- Width/Height (not "horizontal/vertical resolution")
- "Size" (ambiguously means both dimensions AND file size)
- "Quality" (vague — they mean "doesn't look blurry")
- Presets like "Instagram post," "Twitter header," "HD," "4K"

Users do NOT think in terms of: resampling algorithms, DPI, color profiles, aspect ratios (though they expect it to be maintained by default).

## Anti-Patterns to Avoid

- **Don't**: Force account creation or sign-up before using the tool — **Why**: Image resizing is a quick utility task. Any friction before the core action causes immediate bounce. TinyPNG and Squoosh prove that no-auth tools win.
- **Don't**: Show advanced options by default (format conversion, DPI, resampling method, metadata stripping) — **Why**: Overwhelms 90% of users who just want to change width/height. Hide behind an "Advanced" toggle for the 10% who need it.
- **Don't**: Use a multi-page/step wizard flow — **Why**: For a single-purpose tool, everything should happen on one screen. Page transitions feel slow and disorienting for what should be a 15-second task.
- **Don't**: Auto-download without user action — **Why**: Feels invasive. Users want to preview first, then explicitly choose to download. Unexpected downloads trigger browser warnings and erode trust.
- **Don't**: Neglect mobile users — **Why**: Many people need to resize images from their phone (for messaging apps, social media). Touch targets must be large, and the drop zone should gracefully fall back to a file picker on mobile.
- **Don't**: Use tiny, fiddly input fields for dimensions — **Why**: In brutalist design especially, inputs should be large, bold, and easy to tap/click. Small inputs contradict both usability and the aesthetic.
- **Don't**: Crop the image by default when aspect ratio doesn't match — **Why**: Users expect resize to mean scale, not crop. Cropping without consent feels like data loss.

## Design Language: Friendly Brutalism

### Visual Principles
1. **Thick borders** (2-4px solid black or dark color) on all interactive elements
2. **Bold, chunky typography** — large headings in a geometric sans-serif (e.g., Space Grotesk, DM Sans, or Inter Black)
3. **Warm color palette** — not gray/corporate. Think: soft yellow (#FFF3CD), warm coral (#FF6B6B), mint green (#A8E6CF), lavender (#D4A5FF) as accents on a cream/off-white (#FEFAE0) background
4. **Visible shadows** — hard drop shadows (offset, no blur) in black or dark accent color, giving a "sticker" feel
5. **Raw, honest UI** — no gradients, no glass-morphism, no subtle effects. What you see is what you get.
6. **Generous spacing** — brutalism breathes with whitespace. Don't cram elements together.
7. **Playful microcopy** — "Drop your image here (we won't judge the file size)" instead of "Upload file"

### Color Palette Recommendation
| Role | Color | Hex |
|------|-------|-----|
| Background | Cream | #FEFAE0 |
| Surface/Cards | White | #FFFFFF |
| Primary Action | Coral | #FF6B6B |
| Secondary | Mint | #A8E6CF |
| Accent | Lavender | #D4A5FF |
| Text | Near-black | #1A1A2E |
| Borders | Dark | #1A1A2E |
| Success | Green | #4CAF50 |

### Typography
- Headings: Bold/Black weight, 24-48px
- Body: Medium weight, 16-18px
- Inputs: Monospace for dimensions (feels technical/precise), 20px+
- All caps for labels, mixed case for body text

## Recommended Patterns for Our App

Based on research, we should use:

1. **Single-page app with full-viewport drop zone as landing state** — because Squoosh and TinyPNG prove this is the fastest path to engagement. The entire page invites action. No navigation, no decisions before the core task.

2. **3-zone layout after upload: Preview (left/center) + Controls (right sidebar or bottom panel) + Download bar (bottom)** — because this keeps the image visible while adjusting settings, matching how Squoosh handles it. On mobile, stack vertically: preview on top, controls below.

3. **Dimension inputs with live aspect-ratio lock toggle** — because users expect aspect ratio to be maintained by default (lock icon ON). Large, bold number inputs styled in brutalist fashion (thick borders, monospace font). Include common presets as chunky pill buttons (e.g., "1080x1080 Instagram", "1920x1080 HD").

4. **Instant client-side processing with progress feedback** — because all resizing should happen in-browser (Canvas API / OffscreenCanvas) for privacy and speed. Show a chunky progress bar with personality ("Crunching pixels...").

5. **Prominent download button with file size comparison** — because the download moment is the payoff. Show "Original: 4.2MB → Resized: 380KB (91% smaller)" in bold type. Big, satisfying download button with hard shadow.

6. **Friendly brutalist aesthetic throughout** — thick borders, hard shadows, warm colors, bold type, playful copy. This differentiates us from the sterile/corporate feel of most image tools and makes the experience memorable.

7. **"Reset / New Image" as a clear escape hatch** — after download, make it dead simple to start over. A big "Resize Another" button resets the workspace.

## Interaction Model Summary

```
[LANDING STATE]
┌────────────────────────────────────┐
│                                    │
│     Drop your image here           │
│     (or click to browse)           │
│                                    │
│         ┌──────────┐               │
│         │  📁 icon │               │
│         └──────────┘               │
│                                    │
│     Supports JPG, PNG, WebP        │
└────────────────────────────────────┘

         ↓ (image dropped)

[WORKSPACE STATE]
┌────────────────────────────────────┐
│  Image Resizer          [New] [⚙]  │
├──────────────────┬─────────────────┤
│                  │  Width  [1920]  │
│   Image Preview  │  Height [1080]  │
│                  │  🔗 Lock Ratio  │
│                  │                 │
│                  │  ── Presets ──  │
│                  │  [Instagram]    │
│                  │  [HD 1080p]     │
│                  │  [Twitter]      │
│                  │  [Custom]       │
│                  │                 │
│                  │  Format: [JPG▾] │
│                  │  Quality: [85]  │
├──────────────────┴─────────────────┤
│  4.2MB → 380KB (-91%)  [DOWNLOAD] │
└────────────────────────────────────┘
```

---
Status: READY_FOR_PLANNING
