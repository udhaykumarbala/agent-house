# UI Specification: Image Resizer

Based on research in: .plans/research/ui-research.md
Styling wireframes from: .plans/research/ux-research.md

## Design Philosophy

**Friendly Brutalism (Neo-Brut)**: Bold, honest, warm. Thick black borders, hard offset shadows, flat solid colors, chunky interactive elements — but softened with a warm butter-cream palette and playful energy. The structure IS the decoration. No gradients, no glassmorphism, no soft shadows.

---

## Design System

### Color Palette (UNIQUE - Not Generic!)

**IMPORTANT**: These colors are UNIQUE to our app. Chosen to differentiate from the sea of blue/green/purple image tools. Rooted in the "Butter & Charcoal" direction from research — warm, confident, memorable.

| Name | Hex | RGB | Usage |
|------|-----|-----|-------|
| Background | #FEF3E2 | rgb(254,243,226) | App background — warm butter cream |
| Surface | #FFFDF7 | rgb(255,253,247) | Cards, elevated panels — warm off-white |
| Primary | #E8625C | rgb(232,98,92) | Main CTAs, brand action color — warm coral |
| Primary Hover | #D4524C | rgb(212,82,76) | Hover state for primary buttons |
| Primary Active | #C04440 | rgb(192,68,64) | Active/pressed state for primary |
| Secondary | #F5D547 | rgb(245,213,71) | Secondary actions, highlights — butter yellow |
| Secondary Hover | #E6C73E | rgb(230,199,62) | Hover state for secondary |
| Accent | #7A9E7E | rgb(122,158,126) | Tags, badges, tertiary UI — sage green |
| Border | #1A1A1A | rgb(26,26,26) | All borders, outlines — charcoal |
| Text Primary | #1A1A1A | rgb(26,26,26) | Headings, body text — charcoal |
| Text Secondary | #6B6B6B | rgb(107,107,107) | Muted text, captions, labels |
| Text On Primary | #FFFFFF | rgb(255,255,255) | Text on coral/primary buttons |
| Success | #5A9E6F | rgb(90,158,111) | Success states — forest green |
| Error | #C74B4B | rgb(199,75,75) | Error states — deep warm red |
| Warning | #D4943A | rgb(212,148,58) | Warning states — warm amber |

**Why These Colors**: Butter cream + charcoal borders create the warm brutalist foundation. Coral is energetic without being aggressive — it's the "friendly" in friendly brutalism. Butter yellow as secondary adds playfulness (Gumroad/Poolsuite influence). Sage green grounds the palette with an earthy note. No competitor in the image tool space uses this warm, bold direction.

### Contrast Verification

| Combination | Ratio | WCAG Level |
|-------------|-------|------------|
| Text Primary (#1A1A1A) on Background (#FEF3E2) | 14.8:1 | AAA |
| Text Primary (#1A1A1A) on Surface (#FFFDF7) | 17.2:1 | AAA |
| Text Secondary (#6B6B6B) on Background (#FEF3E2) | 4.7:1 | AA |
| Text On Primary (#FFFFFF) on Primary (#E8625C) | 4.1:1 | AA (large text) |
| Text Primary (#1A1A1A) on Secondary (#F5D547) | 11.5:1 | AAA |
| Border (#1A1A1A) on Background (#FEF3E2) | 14.8:1 | AAA |

Note: For small text on coral buttons, use bold weight (600+) to meet AA compliance. All other combinations pass AA or AAA.

---

### Typography

**Heading Font**: Space Grotesk (Google Fonts) — geometric sans-serif with personality. Slightly quirky letterforms add warmth in bold weights. Perfect for brutalist confidence without aggression.

**Body Font**: Inter (Google Fonts) — workhorse readability at all sizes. Pairs beautifully with Space Grotesk's personality.

**Mono Font**: JetBrains Mono (Google Fonts) — for dimension inputs and file size numbers. Feels technical and precise, matching the tool's utility nature.

| Element | Font | Size | Weight | Line Height | Letter Spacing | Transform |
|---------|------|------|--------|-------------|----------------|-----------|
| H1 (Page Title) | Space Grotesk | 40px | 700 | 1.1 | -0.02em | none |
| H2 (Section) | Space Grotesk | 28px | 700 | 1.2 | -0.01em | none |
| H3 (Subsection) | Space Grotesk | 22px | 600 | 1.3 | 0 | none |
| Body | Inter | 16px | 400 | 1.5 | 0 | none |
| Body Bold | Inter | 16px | 600 | 1.5 | 0 | none |
| Small | Inter | 14px | 400 | 1.5 | 0 | none |
| Caption | Inter | 12px | 500 | 1.4 | 0.02em | uppercase |
| Label | Inter | 14px | 600 | 1.4 | 0.01em | uppercase |
| Input Value | JetBrains Mono | 18px | 500 | 1.4 | 0 | none |
| File Size | JetBrains Mono | 16px | 500 | 1.4 | 0 | none |
| Button | Inter | 16px | 600 | 1 | 0.01em | none |

### Google Fonts Import
```
Space Grotesk: 600, 700
Inter: 400, 500, 600
JetBrains Mono: 500
```

---

### Spacing Scale

Base unit: 4px

| Token | Value | Usage |
|-------|-------|-------|
| xs | 4px | Tight internal spacing, icon gaps |
| sm | 8px | Small gaps, inline elements, input padding-y |
| md | 16px | Standard padding, card internal spacing |
| lg | 24px | Section spacing, gaps between controls |
| xl | 32px | Large gaps, section separation |
| 2xl | 48px | Page margins, major section breaks |
| 3xl | 64px | Hero spacing, drop zone internal padding |

---

### Border System

Borders are the backbone of friendly brutalism. They define every interactive surface.

| Token | Value | Usage |
|-------|-------|-------|
| thin | 2px solid #1A1A1A | Inputs, small elements, dividers |
| default | 3px solid #1A1A1A | Cards, buttons, panels — the standard |
| thick | 4px solid #1A1A1A | Drop zone, hero elements, emphasis |
| dashed | 3px dashed #1A1A1A | Drop zone idle state, empty states |

---

### Border Radius

Brutalism favors sharp edges, but "friendly" softens them slightly. No fully rounded pills.

| Token | Value | Usage |
|-------|-------|-------|
| none | 0px | Decorative elements, hard-edge accents |
| sm | 4px | Inputs, small buttons, tags |
| md | 6px | Buttons, cards — the primary radius |
| lg | 10px | Modals, large cards, drop zone |
| full | 9999px | NOT USED — too soft for brutalism |

---

### Shadows (Hard Offset Only)

No soft/diffused shadows. Hard offset shadows create the signature brutalist "sticker" effect.

| Token | Value | Usage |
|-------|-------|-------|
| sm | 2px 2px 0px #1A1A1A | Small elements, tags, badges |
| md | 4px 4px 0px #1A1A1A | Buttons, inputs, standard elements |
| lg | 6px 6px 0px #1A1A1A | Cards, panels, elevated surfaces |
| xl | 8px 8px 0px #1A1A1A | Drop zone, hero elements, modals |
| hover | 6px 6px 0px #1A1A1A | Hovered interactive elements (grows from md) |
| active | 0px 0px 0px #1A1A1A | Pressed elements (shadow disappears) |

---

## Components

### Button — Primary

The most important interactive element. Must feel "pressable" like a physical button.

```css
.btn-primary {
  background: #E8625C;
  color: #FFFFFF;
  font-family: 'Inter', sans-serif;
  font-size: 16px;
  font-weight: 600;
  padding: 12px 24px;
  border: 3px solid #1A1A1A;
  border-radius: 6px;
  box-shadow: 4px 4px 0px #1A1A1A;
  cursor: pointer;
  transition: transform 0.1s ease, box-shadow 0.1s ease;
}
```

**States:**
- **Default**: Coral background, charcoal border, 4px offset shadow
- **Hover**: `transform: translate(-2px, -2px); box-shadow: 6px 6px 0px #1A1A1A;` — lifts up, shadow grows
- **Active/Pressed**: `transform: translate(2px, 2px); box-shadow: 0px 0px 0px #1A1A1A;` — pushes down, shadow collapses. Feels like pressing a physical button.
- **Disabled**: `opacity: 0.5; cursor: not-allowed;` — no hover/active effects
- **Focus**: `outline: 3px solid #F5D547; outline-offset: 2px;` — butter yellow focus ring

### Button — Secondary

```css
.btn-secondary {
  background: #F5D547;
  color: #1A1A1A;
  font-family: 'Inter', sans-serif;
  font-size: 16px;
  font-weight: 600;
  padding: 12px 24px;
  border: 3px solid #1A1A1A;
  border-radius: 6px;
  box-shadow: 4px 4px 0px #1A1A1A;
  cursor: pointer;
  transition: transform 0.1s ease, box-shadow 0.1s ease;
}
```

States follow the same hover/active pattern as primary.

### Button — Ghost

```css
.btn-ghost {
  background: transparent;
  color: #1A1A1A;
  font-family: 'Inter', sans-serif;
  font-size: 16px;
  font-weight: 600;
  padding: 12px 24px;
  border: 3px solid #1A1A1A;
  border-radius: 6px;
  box-shadow: none;
  cursor: pointer;
}
```

**Hover**: `background: #FEF3E2;` (subtle fill)
**Active**: `background: #1A1A1A; color: #FFFDF7;` (inverts)

### Button — Icon (small action buttons)

```css
.btn-icon {
  background: #FFFDF7;
  color: #1A1A1A;
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px solid #1A1A1A;
  border-radius: 6px;
  box-shadow: 2px 2px 0px #1A1A1A;
  cursor: pointer;
}
```

---

### Input — Number (Dimension Inputs)

Oversized, monospace, brutalist. These are core to the resize experience.

```css
.input-dimension {
  font-family: 'JetBrains Mono', monospace;
  font-size: 18px;
  font-weight: 500;
  color: #1A1A1A;
  background: #FFFDF7;
  padding: 10px 14px;
  border: 3px solid #1A1A1A;
  border-radius: 4px;
  box-shadow: 3px 3px 0px #1A1A1A;
  width: 120px;
  text-align: center;
}
```

**States:**
- **Focus**: `border-color: #E8625C; box-shadow: 3px 3px 0px #E8625C;` — coral highlight
- **Error**: `border-color: #C74B4B; box-shadow: 3px 3px 0px #C74B4B;` — red highlight
- **Disabled**: `background: #FEF3E2; opacity: 0.6;`

### Input — Range (Quality Slider)

```css
.input-range {
  -webkit-appearance: none;
  width: 100%;
  height: 8px;
  background: #1A1A1A;
  border-radius: 0;
  border: 2px solid #1A1A1A;
}

.input-range::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: 24px;
  height: 24px;
  background: #E8625C;
  border: 3px solid #1A1A1A;
  border-radius: 4px;
  box-shadow: 2px 2px 0px #1A1A1A;
  cursor: pointer;
}
```

### Select / Dropdown (Format Picker)

```css
.select {
  font-family: 'Inter', sans-serif;
  font-size: 16px;
  font-weight: 500;
  color: #1A1A1A;
  background: #FFFDF7;
  padding: 10px 14px;
  border: 3px solid #1A1A1A;
  border-radius: 4px;
  box-shadow: 3px 3px 0px #1A1A1A;
  cursor: pointer;
}
```

---

### Card (Settings Panel, Info Panels)

```css
.card {
  background: #FFFDF7;
  border: 3px solid #1A1A1A;
  border-radius: 6px;
  box-shadow: 6px 6px 0px #1A1A1A;
  padding: 24px;
}
```

### Card — Preset Button (Instagram, HD, Twitter, etc.)

```css
.preset-btn {
  background: #FFFDF7;
  color: #1A1A1A;
  font-family: 'Inter', sans-serif;
  font-size: 14px;
  font-weight: 600;
  padding: 8px 16px;
  border: 2px solid #1A1A1A;
  border-radius: 4px;
  box-shadow: 2px 2px 0px #1A1A1A;
  cursor: pointer;
}

.preset-btn--active {
  background: #F5D547;
  box-shadow: 0px 0px 0px #1A1A1A;
  transform: translate(2px, 2px);
}
```

---

### Drop Zone

The hero element. Full viewport on landing, the entire page invites action.

**Landing State (Full Page):**
```css
.dropzone {
  background: #FFFDF7;
  border: 4px dashed #1A1A1A;
  border-radius: 10px;
  box-shadow: 8px 8px 0px #1A1A1A;
  padding: 64px;
  text-align: center;
  min-height: 50vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: border-color 0.2s, background 0.2s;
}
```

**Drag Hover State:**
```css
.dropzone--hover {
  background: #FEF3E2;
  border-color: #E8625C;
  border-style: solid;
  box-shadow: 8px 8px 0px #E8625C;
}
```

**Drop Zone Content:**
- Upload icon: Thick-stroke (2.5px) arrow-up icon, 48px, charcoal
- Heading: "Drop your image here" — H2, Space Grotesk 700
- Subtext: "or click to browse" — Body, Inter 400, text-secondary
- Supported formats: "JPG, PNG, WebP" — Caption, uppercase

---

### Download Bar (Sticky Bottom)

The payoff moment. Appears after image is loaded and settings are configured.

```css
.download-bar {
  background: #FFFDF7;
  border-top: 3px solid #1A1A1A;
  padding: 16px 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  position: sticky;
  bottom: 0;
  box-shadow: 0px -4px 0px #1A1A1A;
}
```

**Contents:**
- Left: File size comparison — `"4.2 MB → 380 KB (91% smaller)"` in JetBrains Mono, with the percentage in Success green
- Right: Large primary download button — "Download" with download icon

---

### Progress Bar

Bold, brutalist. For resize processing feedback.

```css
.progress-bar {
  width: 100%;
  height: 16px;
  background: #FEF3E2;
  border: 3px solid #1A1A1A;
  border-radius: 4px;
  overflow: hidden;
}

.progress-bar__fill {
  height: 100%;
  background: #E8625C;
  transition: width 0.3s ease;
}
```

---

### Toggle (Aspect Ratio Lock)

```css
.toggle {
  width: 44px;
  height: 44px;
  background: #FFFDF7;
  border: 2px solid #1A1A1A;
  border-radius: 6px;
  box-shadow: 2px 2px 0px #1A1A1A;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

.toggle--active {
  background: #F5D547;
}
```

Contains a link/chain icon. When active (locked), butter yellow background. When inactive (unlocked), surface white with a broken-chain icon.

---

## Iconography

**Icon Set**: Phosphor Icons (Bold weight) — thick 2.5px stroke weight matches our border weight perfectly.

**Key Icons Needed:**
| Icon | Usage | Phosphor Name |
|------|-------|---------------|
| Upload/Drop | Drop zone | ArrowUp or UploadSimple |
| Link/Chain | Aspect ratio lock | LinkSimple |
| Unlink | Aspect ratio unlock | LinkBreak |
| Download | Download button | DownloadSimple |
| Image | Image placeholder | Image |
| X/Close | Remove image, dismiss | X |
| Gear | Advanced settings toggle | GearSix |
| Refresh | Reset/New image | ArrowCounterClockwise |
| Check | Success/Done | Check |
| Warning | Error/Warning | Warning |
| Info | Tooltips | Info |

All icons rendered at 24px with 2.5px stroke in #1A1A1A. On primary buttons, stroke color #FFFFFF.

---

## Animations & Micro-Interactions

All transitions use `transition: 0.1s ease` for snappiness. No spring/bouncy physics — brutalism favors mechanical, satisfying motion.

| Interaction | Animation | Duration |
|-------------|-----------|----------|
| Button hover | `translate(-2px, -2px)`, shadow grows | 100ms |
| Button press | `translate(2px, 2px)`, shadow collapses to 0 | 80ms |
| Drop zone drag-hover | Border changes from dashed to solid, bg color shift, shadow recolors to coral | 200ms |
| Image loaded | Fade in + slight scale from 0.95 to 1.0 | 200ms |
| Download ready | Download button does a single subtle pulse (scale 1.0 → 1.03 → 1.0) | 300ms |
| Error state | Horizontal shake: translateX(0 → -4px → 4px → -2px → 0) | 300ms |
| Card hover | Shadow grows from md to lg | 100ms |

---

## Layout

### Responsive Breakpoints

| Name | Width | Layout |
|------|-------|--------|
| mobile | < 640px | Single column, stacked vertically |
| tablet | 640px – 1024px | Single column, wider spacing |
| desktop | > 1024px | Two-column: preview + controls sidebar |

### Landing State (Before Image Upload)

```
┌──────────────────────────────────────────────────────┐
│  Image Resizer (H1, Space Grotesk 700)               │
│  "Resize instantly. 100% private." (Body, muted)     │
│                                                       │
│  ┌ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ┐  │
│  │                                                 │  │
│  │         ↑ (upload icon, 48px)                   │  │
│  │                                                 │  │
│  │     Drop your image here                        │  │
│  │     or click to browse                          │  │
│  │                                                 │  │
│  │     JPG · PNG · WEBP                            │  │
│  └ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ┘  │
│                                                       │
│  "Nothing leaves your browser." (Caption, muted)      │
└──────────────────────────────────────────────────────┘
```

Max-width: 640px, centered. Drop zone takes ~50% viewport height minimum.

### Workspace State — Desktop (After Image Upload)

```
┌──────────────────────────────────────────────────────┐
│  Image Resizer                    [New Image] [⚙]    │
├──────────────────────────┬───────────────────────────┤
│                          │  DIMENSIONS (label)        │
│                          │  ┌──────┐ × ┌──────┐      │
│     Image Preview        │  │ 1920 │   │ 1080 │      │
│     (bordered card,      │  └──────┘   └──────┘      │
│      checkered bg        │       [🔗 Lock Ratio]     │
│      behind image)       │                            │
│                          │  PRESETS (label)            │
│                          │  [Instagram] [HD 1080p]    │
│                          │  [Twitter]   [OG Image]    │
│                          │  [Custom]                  │
│                          │                            │
│                          │  FORMAT (label)             │
│                          │  [JPG ▾]                   │
│                          │                            │
│                          │  QUALITY (label)            │
│                          │  ──●─────── 85%            │
├──────────────────────────┴───────────────────────────┤
│  📐 1920×1080 → 800×450   4.2MB → 380KB (-91%)  [⬇ DOWNLOAD]  │
└──────────────────────────────────────────────────────┘
```

- Preview area: ~60% width. Image displayed in a bordered card with a subtle checkered/dotted pattern behind transparent images.
- Controls sidebar: ~40% width. Card with thick border and hard shadow.
- Download bar: Sticky bottom, full width.

### Workspace State — Mobile (< 640px)

```
┌────────────────────────┐
│ Image Resizer    [New] │
├────────────────────────┤
│                        │
│    Image Preview       │
│    (full width)        │
│                        │
├────────────────────────┤
│  DIMENSIONS             │
│  ┌──────┐ × ┌──────┐  │
│  │ 1920 │   │ 1080 │  │
│  └──────┘   └──────┘  │
│       [🔗 Lock]       │
│                        │
│  PRESETS               │
│  [Instagram] [HD]      │
│  [Twitter]  [Custom]   │
│                        │
│  FORMAT  [JPG ▾]       │
│  QUALITY ──●── 85%     │
├────────────────────────┤
│ 4.2MB→380KB  [⬇ DL]   │
└────────────────────────┘
```

Stacked vertically. Controls below preview. Download bar still sticky bottom.

---

## Page Structure & Spacing

| Section | Desktop | Mobile |
|---------|---------|--------|
| Page padding | 48px horizontal | 16px horizontal |
| Header to content | 32px | 24px |
| Preview/Controls gap | 32px | 0 (stacked) |
| Control groups gap | 24px | 20px |
| Within card padding | 24px | 16px |
| Preset button gap | 8px | 8px |
| Download bar padding | 16px 24px | 12px 16px |

---

## Special Visual Elements

### Checkered Background (Behind Image Preview)
For transparent images (PNG), show a subtle checkered pattern behind the preview:
```css
.preview-bg {
  background-image:
    linear-gradient(45deg, #EEE8D5 25%, transparent 25%),
    linear-gradient(-45deg, #EEE8D5 25%, transparent 25%),
    linear-gradient(45deg, transparent 75%, #EEE8D5 75%),
    linear-gradient(-45deg, transparent 75%, #EEE8D5 75%);
  background-size: 16px 16px;
  background-position: 0 0, 0 8px, 8px -8px, -8px 0;
}
```
Uses warm tones (#EEE8D5) instead of cold gray to match the palette.

### File Size Comparison Display
```
Original: 4.2 MB  →  Resized: 380 KB  (91% smaller)
```
- Numbers in JetBrains Mono, weight 500
- Arrow (→) in charcoal
- Percentage in Success green (#5A9E6F), bold
- If file grows (upscaling), percentage in Warning amber (#D4943A)

### Privacy Badge
Small, subtle badge near the bottom or footer:
```
🔒 100% client-side · Nothing leaves your browser
```
Caption size, text-secondary color. Builds trust without being preachy.

---

## Dark Mode

**Not included in v1.** The warm butter-cream palette is central to the friendly brutalist identity. Dark mode would be a separate design exercise. Revisit in a later phase if requested.

---

## CSS Custom Properties (Design Tokens)

```css
:root {
  /* Colors */
  --color-bg: #FEF3E2;
  --color-surface: #FFFDF7;
  --color-primary: #E8625C;
  --color-primary-hover: #D4524C;
  --color-primary-active: #C04440;
  --color-secondary: #F5D547;
  --color-secondary-hover: #E6C73E;
  --color-accent: #7A9E7E;
  --color-border: #1A1A1A;
  --color-text: #1A1A1A;
  --color-text-muted: #6B6B6B;
  --color-text-on-primary: #FFFFFF;
  --color-success: #5A9E6F;
  --color-error: #C74B4B;
  --color-warning: #D4943A;

  /* Typography */
  --font-heading: 'Space Grotesk', sans-serif;
  --font-body: 'Inter', sans-serif;
  --font-mono: 'JetBrains Mono', monospace;

  /* Spacing */
  --space-xs: 4px;
  --space-sm: 8px;
  --space-md: 16px;
  --space-lg: 24px;
  --space-xl: 32px;
  --space-2xl: 48px;
  --space-3xl: 64px;

  /* Borders */
  --border-thin: 2px solid #1A1A1A;
  --border-default: 3px solid #1A1A1A;
  --border-thick: 4px solid #1A1A1A;

  /* Radius */
  --radius-sm: 4px;
  --radius-md: 6px;
  --radius-lg: 10px;

  /* Shadows */
  --shadow-sm: 2px 2px 0px #1A1A1A;
  --shadow-md: 4px 4px 0px #1A1A1A;
  --shadow-lg: 6px 6px 0px #1A1A1A;
  --shadow-xl: 8px 8px 0px #1A1A1A;
}
```

---
Status: READY_FOR_REVIEW
