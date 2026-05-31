# Product Specification: Image Resizer (Friendly Brutalist)

Based on research in: .plans/research/product-research.md, .plans/research/ux-research.md, .plans/research/ui-research.md, .plans/research/architect-research.md, .plans/research/security.md

---

## Target User

**"Quick-Fix Quinn"** — Age 22–40, content creator / social media manager / blogger / small business owner. Not a designer, but regularly needs to resize images for posts, listings, emails, or uploads. Uses web tools daily, doesn't want to install desktop software for a 10-second task. Frustrated by ad-heavy, slow, untrustworthy resize tools that require server uploads.

## Unique Value Proposition

**"Resize images instantly in your browser — no uploads, no ads, no nonsense. Just bold vibes and fast results."**

The only image resizer that's both **completely private** (100% client-side) and **actually enjoyable to use** (friendly brutalist design with personality). Done in 3 clicks: drop, resize, download.

## Key Differentiators

1. **Privacy-first** — nothing leaves the browser. No server uploads, no tracking.
2. **Friendly brutalist design** — visually distinctive in a sea of generic utility tools. Bold, warm, memorable.
3. **Opinionated simplicity** — not Photoshop. One job, done fast and well.
4. **No ads, no accounts, no upsells** — pure utility with great vibes.

---

## Features (Prioritized)

### P0 — Must Have (MVP)

- [ ] **F1: Drag-and-Drop Image Upload** — Full-viewport drop zone as the landing state. Also supports click-to-browse via file picker. Accepts JPG, PNG, WebP. This IS the entry point — the entire page invites action.
- [ ] **F2: Resize by Pixel Dimensions** — Width and height inputs with direct pixel entry. Large, bold, monospace-styled inputs fitting the brutalist aesthetic. This is the core function of the app.
- [ ] **F3: Aspect Ratio Lock** — Toggle (on by default) that maintains proportions when changing width or height. Users expect this; not having it causes accidental distortion.
- [ ] **F4: Live Image Preview** — Instantly render the resized result in the workspace so users see what they're getting before download. Builds confidence, matches Squoosh's gold-standard UX.
- [ ] **F5: One-Click Download** — Prominent download button that saves the resized image to the user's device. Sensible auto-generated filename (e.g., `original-name-1080x720.jpg`).
- [ ] **F6: Output Format Selection** — Choose between JPEG, PNG, or WebP for the output. Dropdown or segmented toggle. Enables quick format conversion during resize (a gap most competitors miss).
- [ ] **F7: Quality Slider** — For JPEG and WebP output, a slider (1–100) to control compression. Not applicable when PNG is selected. Defaults to 85 (good balance of quality and size).
- [ ] **F8: File Size Display** — Show original file size and estimated output size with percentage change (e.g., "4.2MB → 380KB · 91% smaller"). This is the feedback metric users want most.
- [ ] **F9: Friendly Brutalist UI** — The visual identity IS the product differentiator. Butter/cream background, thick black borders, hard offset shadows, Space Grotesk + Inter typography, warm accent colors, playful microcopy. Not decorative — structural.
- [ ] **F10: Privacy Messaging** — Front-and-center statement: "Your images never leave your browser." No competitor highlights this prominently. It's a trust signal and differentiator.
- [ ] **F11: Mobile Responsive** — Fully usable on mobile. Drop zone falls back to file picker. Controls stack vertically. Touch targets ≥ 44x44px. Many users resize images from their phone.

### P1 — Should Have

- [ ] **F12: Preset Dimension Buttons** — Chunky pill buttons for common sizes: Instagram Post (1080×1080), Instagram Story (1080×1920), Twitter/X Header (1500×500), Facebook Cover (820×312), HD (1920×1080), OG Image (1200×630). Squoosh doesn't have this; iLoveIMG buries it. Major UX win.
- [ ] **F13: Resize by Percentage** — Scale by percentage (25%, 50%, 75%, custom %) as an alternative to pixel dimensions. Some users think in relative terms ("make it half the size").
- [ ] **F14: Before/After Comparison** — Show original vs. resized dimensions and file size side by side. Reinforces the value delivered.
- [ ] **F15: Reset / Resize Another** — Clear escape hatch after download. Big "Resize Another" button that resets the workspace to the drop zone. Keeps the user in the loop for multiple images.

### P2 — Nice to Have (Post-MVP)

- [ ] **F16: Batch Resize** — Upload and resize multiple images at once with the same settings. Download as ZIP. Power-user feature.
- [ ] **F17: Simple Crop** — Basic crop tool before resize. Keeps users from needing a separate tool.
- [ ] **F18: Keyboard Shortcuts** — Tab through controls, Enter to download. Power-user efficiency.
- [ ] **F19: PWA / Installable** — Make the app installable and usable offline. Supports the "bookmarkable utility" distribution model.

---

## User Stories

1. **As a** social media manager, **I want** to drop an image and instantly see resize controls **so that** I can resize assets without any setup or sign-up friction.
2. **As a** blogger, **I want** to resize an image by entering specific pixel dimensions **so that** my images fit my site's layout requirements exactly.
3. **As a** content creator, **I want** the aspect ratio to be locked by default **so that** my images don't get accidentally distorted.
4. **As a** small business owner, **I want** to see a live preview of the resized image **so that** I can verify quality before downloading.
5. **As a** privacy-conscious user, **I want** confirmation that my images never leave my browser **so that** I feel safe resizing photos of my clients' products.
6. **As a** social media manager, **I want** preset dimension buttons for Instagram, Twitter, etc. **so that** I don't have to look up the correct dimensions every time.
7. **As a** email marketer, **I want** to see the file size before and after resizing **so that** I know the image will fit within my email service's upload limit.
8. **As a** mobile user, **I want** the tool to work on my phone **so that** I can quickly resize a photo before posting to social media.
9. **As a** web developer, **I want** to convert PNG to WebP during resize **so that** I can optimize images for the web in one step.
10. **As a** returning user, **I want** the tool to look distinctive and memorable **so that** I can find it in my bookmarks and recognize it instantly.

---

## Acceptance Criteria

### F1: Drag-and-Drop Image Upload
- [ ] Full-viewport drop zone is the initial landing state
- [ ] Drag hover shows visual feedback (animated border, color change)
- [ ] Clicking the drop zone opens the native file picker
- [ ] Accepts JPEG, PNG, and WebP files
- [ ] Rejects non-image files with a friendly error message
- [ ] Validates file size ≤ 50MB; warns the user if exceeded
- [ ] Validates file type via MIME type check
- [ ] Transitions to workspace state after successful upload

### F2: Resize by Pixel Dimensions
- [ ] Width and height inputs are large, bold, and monospace-styled
- [ ] Inputs accept only positive integers
- [ ] Original dimensions are shown as reference
- [ ] Changing width/height triggers live preview update
- [ ] Inputs are pre-populated with original image dimensions on load

### F3: Aspect Ratio Lock
- [ ] Lock toggle is ON by default
- [ ] When locked, changing width auto-calculates height (and vice versa)
- [ ] Lock icon visually indicates locked/unlocked state
- [ ] Toggling off allows independent width/height entry

### F4: Live Image Preview
- [ ] Preview updates within 500ms of dimension change
- [ ] Preview shows the resized image at a viewport-appropriate display size
- [ ] Processing feedback shown for large images ("Crunching pixels...")

### F5: One-Click Download
- [ ] Download button is prominent with hard offset shadow (brutalist style)
- [ ] Downloaded file uses sensible name: `{original-name}-{width}x{height}.{format}`
- [ ] Download filename is sanitized (no path traversal, control characters)
- [ ] File downloads immediately on click (no redirect, no popup)

### F6: Output Format Selection
- [ ] JPEG, PNG, and WebP are available as output options
- [ ] Default output matches the input format
- [ ] Selecting PNG hides/disables the quality slider (PNG is lossless)
- [ ] Format change triggers preview and file size recalculation

### F7: Quality Slider
- [ ] Slider range: 1–100, default: 85
- [ ] Only visible/active for JPEG and WebP output
- [ ] Changing quality updates the estimated file size
- [ ] Slider is styled with thick borders matching brutalist design

### F8: File Size Display
- [ ] Shows original file size (e.g., "4.2 MB")
- [ ] Shows estimated output size after resize
- [ ] Shows percentage change (e.g., "91% smaller")
- [ ] Updates live as dimensions, format, or quality change

### F9: Friendly Brutalist UI
- [ ] Background is warm cream/butter, not white or gray
- [ ] All interactive elements have thick (2–3px) black borders
- [ ] Buttons use hard offset box-shadows (e.g., `4px 4px 0 #1A1A1A`)
- [ ] Hover: element lifts (translate -2px) with shadow growing
- [ ] Active/click: element presses (translate +2px) with shadow shrinking
- [ ] Typography uses Space Grotesk (headings) and Inter (body)
- [ ] Playful microcopy (e.g., "Drop your image here — we won't peek")
- [ ] No gradients, glassmorphism, or soft shadows

### F10: Privacy Messaging
- [ ] "Your images never leave your browser" is visible on the landing state
- [ ] Brief explanation of client-side processing is accessible (tooltip or footer)

### F11: Mobile Responsive
- [ ] Layout stacks vertically on screens < 768px
- [ ] Drop zone gracefully shows file picker on mobile
- [ ] All touch targets are ≥ 44x44px
- [ ] Inputs and buttons are large enough for thumb interaction
- [ ] No horizontal scroll on any viewport

---

## Security Requirements (from security research)

These are non-negotiable implementation constraints:

| # | Requirement | Priority |
|---|-------------|----------|
| S1 | Use `textContent` only — never `innerHTML` for file names, metadata, or user-derived strings | Critical |
| S2 | Validate file type (MIME check) and enforce max size (50MB) before processing | High |
| S3 | Configure Content Security Policy (CSP) headers on deployment | High |
| S4 | Revoke `URL.revokeObjectURL()` when previews are no longer needed | Medium |
| S5 | Sanitize download file names (strip `../`, null bytes, control characters) | Medium |
| S6 | All JS in external files (no inline scripts) to support strict CSP | Medium |
| S7 | Wrap canvas operations in try/catch for graceful error handling | Medium |

---

## Design Specifications (from UI/UX research)

### Color Palette

| Role | Color | Hex |
|------|-------|-----|
| Background | Warm cream | #FEFAE0 |
| Surface/Cards | White | #FFFFFF |
| Primary Action | Coral | #FF6B6B |
| Secondary | Mint | #A8E6CF |
| Accent | Lavender | #D4A5FF |
| Text / Borders | Near-black | #1A1A2E |
| Success | Green | #4CAF50 |
| Hover/Active | Warm yellow | #FFE66D |

### Typography
- **Headings**: Space Grotesk, Bold/Black, 24–48px
- **Body/Labels**: Inter, Medium, 16–18px
- **Dimension Inputs**: Monospace (JetBrains Mono or system), 20px+
- **Labels**: ALL CAPS, letter-spacing for emphasis

### Layout Pattern
- Single-page, no navigation
- Landing state: full-viewport drop zone (≥50% of viewport on desktop)
- Workspace state: Preview (left/center) + Controls (right sidebar) + Download bar (bottom)
- Mobile: vertical stack — preview on top, controls below, download sticky at bottom

### Interaction Wireframe
```
[LANDING STATE]
┌──────────────────────────────────────┐
│                                      │
│     Drop your image here             │
│     (or click to browse)             │
│                                      │
│     "Your images never leave         │
│      your browser"                   │
│                                      │
│     Supports JPG, PNG, WebP          │
└──────────────────────────────────────┘

         ↓ (image dropped)

[WORKSPACE STATE]
┌──────────────────────────────────────┐
│  Image Resizer              [New] [⚙]│
├────────────────────┬─────────────────┤
│                    │  Width  [1920]  │
│   Image Preview    │  Height [1080]  │
│                    │  🔗 Lock Ratio  │
│                    │                 │
│                    │  ── Presets ──  │
│                    │  [Instagram]    │
│                    │  [HD 1080p]     │
│                    │  [Twitter]      │
│                    │  [OG Image]     │
│                    │                 │
│                    │  Format: [JPG▾] │
│                    │  Quality: [85]  │
├────────────────────┴─────────────────┤
│  4.2MB → 380KB (91% smaller)        │
│                          [DOWNLOAD]  │
└──────────────────────────────────────┘
```

---

## Technical Constraints (from architect research)

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Processing | 100% client-side | Privacy, speed, no server cost |
| Image API | HTML5 Canvas API | Native, zero dependencies, sufficient for resize |
| Framework | Vanilla JS | Overkill to add a framework for a single-page tool |
| Styling | Custom CSS | Brutalist design is hand-crafted by nature |
| Build tools | None | No bundler, no npm — matches minimal ethos |
| External deps | Google Fonts only | Minimal footprint; consider self-hosting for privacy |
| File structure | index.html, css/style.css, js/app.js, js/resizer.js | Separation of concerns without over-engineering |

---

## Success Metrics

- **Time to first resize**: < 30 seconds from page load to downloaded file
- **Bounce rate**: < 40% (indicates the landing state is inviting, not confusing)
- **Return visits**: > 20% of users come back (indicates memorability and utility)
- **Core task completion**: > 90% of users who upload an image successfully download a resized version
- **Mobile usability**: Fully functional on iOS Safari and Android Chrome

---

## Out of Scope (for MVP)

- Batch processing (multiple images at once)
- Cropping or rotating tools
- AI upscaling or enhancement
- User accounts or saved preferences
- Server-side processing of any kind
- Image editing (filters, text overlay, etc.)
- AVIF or JPEG XL support (requires WASM, adds complexity)
- Analytics or tracking scripts
- Monetization features

---

## Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Large images crash mobile browsers | Medium | Validate file size ≤ 50MB, show warning for files > 20MB |
| Users expect server-side features (batch, AI) | Low | Clear messaging: "fast, private, single-image resizer" |
| Canvas color profile mismatch | Low | Document limitation, use sRGB |
| WebP not supported on older Safari | Low | Default to JPEG, detect WebP support |
| Design feels too niche/polarizing | Low | "Friendly" brutalism softens the edges; warm colors keep it approachable |

---

## Distribution Strategy

- **Primary**: Direct URL, bookmarkable
- **Secondary**: Shareable — the distinctive brutalist design is screenshot-worthy
- **Tertiary**: PWA potential (post-MVP) for installable offline use
- **SEO**: Target "resize image online," "image resizer no upload," "private image resizer"

---
Status: READY_FOR_REVIEW
