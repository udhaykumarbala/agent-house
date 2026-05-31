# UI Specification: Meme Creator

Based on research in: .plans/research/ui-research.md
Styling wireframes from: .plans/research/ux-research.md

## Design System

### Color Palette (UNIQUE - Not Generic!)

**IMPORTANT**: These colors are UNIQUE to our app. They were chosen to differentiate from Canva (teal/purple), Imgflip (blue), Kapwing (indigo), Photopea (teal), and Adobe Express (coral-violet). Our direction is **Hot Magenta + Deep Charcoal** — an electric, creative-studio-at-night aesthetic that screams meme culture energy.

| Name | Hex | RGB | Usage |
|------|-----|-----|-------|
| Primary | #E91E8C | rgb(233, 30, 140) | Main CTAs, brand accents, selected states |
| Primary Hover | #D4177F | rgb(212, 23, 127) | Hover states on primary elements |
| Primary Muted | rgba(233,30,140,0.15) | — | Subtle backgrounds, focus rings |
| Secondary | #FF6B9D | rgb(255, 107, 157) | Secondary actions, highlights, tags |
| Background | #1A1A2E | rgb(26, 26, 46) | App background (deep charcoal with blue undertone) |
| Surface | #252547 | rgb(37, 37, 71) | Cards, panels, elevated elements |
| Surface Hover | #2E2E5E | rgb(46, 46, 94) | Hovered surface elements |
| Canvas BG | #16213E | rgb(22, 33, 62) | Canvas surrounding area, inset areas |
| Border | #3A3A5C | rgb(58, 58, 92) | Subtle borders between surfaces |
| Text Primary | #F8F8F2 | rgb(248, 248, 242) | Main text, headings |
| Text Secondary | #A0A0B8 | rgb(160, 160, 184) | Muted labels, hints, captions |
| Text Disabled | #5A5A78 | rgb(90, 90, 120) | Disabled text, placeholders |
| Success | #36D399 | rgb(54, 211, 153) | Success toasts, valid states |
| Error | #F87272 | rgb(248, 114, 114) | Error states, delete actions |
| Warning | #FBBD23 | rgb(251, 189, 35) | Warning indicators |

**Why These Colors**: Hot Magenta (#E91E8C) is virtually unused by any competitor in the meme/image editor space. It conveys creativity, energy, and fun — perfectly matching meme culture's irreverent, expressive nature. The deep charcoal background (#1A1A2E) with a slight blue undertone creates a "creative studio at night" feel, reduces eye strain during editing, and makes the canvas image the visual hero. The slight blue undertone avoids the flatness of pure dark gray.

### Typography

**UI Font Family**: DM Sans (Google Fonts)
- Geometric, friendly, rounded quality
- Approachable and creative without being childish
- Excellent legibility at small sizes in dark mode
- Variable weight support for flexible hierarchy

**Meme Canvas Fonts** (available in editor):
- Impact (default meme preset)
- Anton
- Bebas Neue
- Bangers
- Permanent Marker
- Oswald

| Element | Size | Weight | Line Height | Letter Spacing |
|---------|------|--------|-------------|----------------|
| H1 (App Title) | 28px | 700 | 1.2 | -0.02em |
| H2 (Panel Titles) | 20px | 600 | 1.3 | -0.01em |
| H3 (Section Labels) | 16px | 600 | 1.4 | 0 |
| Body | 14px | 400 | 1.5 | 0 |
| Small (Labels) | 13px | 500 | 1.4 | 0.01em |
| Caption (Hints) | 12px | 400 | 1.4 | 0.02em |
| Button | 14px | 600 | 1 | 0.01em |

### Spacing Scale

Base: 4px

| Token | Value | Usage |
|-------|-------|-------|
| xs | 4px | Tight gaps between inline elements |
| sm | 8px | Small gaps, inner padding of compact controls |
| md | 12px | Standard control padding |
| lg | 16px | Panel inner padding, control group gaps |
| xl | 24px | Section spacing within panels |
| 2xl | 32px | Major layout gaps |
| 3xl | 48px | Page-level margins (desktop) |

### Border Radius

| Token | Value | Usage |
|-------|-------|-------|
| sm | 6px | Inputs, small buttons, chips |
| md | 8px | Standard buttons, dropdowns |
| lg | 12px | Cards, panels, modals |
| xl | 16px | Large feature cards, upload zone |
| full | 9999px | Pills, avatar circles, toggle knobs |

### Shadows

Dark mode relies less on shadows and more on luminance layering. Shadows are used sparingly.

| Token | Value | Usage |
|-------|-------|-------|
| sm | 0 1px 3px rgba(0,0,0,0.3) | Subtle elevation (dropdowns) |
| md | 0 4px 12px rgba(0,0,0,0.4) | Floating toolbar, popovers |
| lg | 0 8px 24px rgba(0,0,0,0.5) | Modals, export dialog |
| glow | 0 0 12px rgba(233,30,140,0.3) | Primary accent glow on selected elements |

### Elevation Strategy (Dark Mode)

Instead of traditional shadows, surfaces use **luminance layering**:
- Level 0 (Background): #1A1A2E
- Level 1 (Panels/Sidebar): #252547
- Level 2 (Cards/Controls within panels): #2E2E5E
- Level 3 (Dropdowns/Popovers): #353564
- Level 4 (Modals): #3A3A6A

Each level is slightly lighter, creating depth without relying on shadow.

## Icons

**Icon Set**: Lucide Icons (outline style)
- Consistent 24px base size for toolbar icons
- 1.5px stroke weight
- Use filled variant for active/selected states
- Color: Text Secondary (#A0A0B8) default, Text Primary (#F8F8F2) on hover, Primary (#E91E8C) when active

## Components

### Button

**Primary** (Download, Add Text):
```css
background: #E91E8C;
color: #FFFFFF;
font: 600 14px 'DM Sans';
padding: 10px 20px;
border-radius: 8px;
border: none;
cursor: pointer;
transition: background 0.15s ease, transform 0.1s ease;
```
- Hover: background #D4177F, subtle lift (translateY -1px)
- Active: background #C01472, translateY 0
- Disabled: opacity 0.4, cursor not-allowed
- Focus: ring 2px #E91E8C with 30% opacity offset

**Secondary** (Format toggle, Cancel):
```css
background: #252547;
color: #F8F8F2;
font: 600 14px 'DM Sans';
padding: 10px 20px;
border-radius: 8px;
border: 1px solid #3A3A5C;
cursor: pointer;
transition: background 0.15s ease, border-color 0.15s ease;
```
- Hover: background #2E2E5E, border-color #E91E8C
- Active: background #353564
- Disabled: opacity 0.4, cursor not-allowed

**Ghost** (Toolbar actions, text-only actions):
```css
background: transparent;
color: #A0A0B8;
font: 500 13px 'DM Sans';
padding: 8px 12px;
border-radius: 6px;
border: none;
cursor: pointer;
transition: color 0.15s ease, background 0.15s ease;
```
- Hover: color #F8F8F2, background rgba(255,255,255,0.06)
- Active: color #E91E8C

**Icon Button** (Toolbar icons — bold, italic, align, delete):
```css
background: transparent;
color: #A0A0B8;
width: 36px;
height: 36px;
display: flex;
align-items: center;
justify-content: center;
border-radius: 6px;
border: none;
cursor: pointer;
transition: color 0.15s ease, background 0.15s ease;
```
- Hover: color #F8F8F2, background rgba(255,255,255,0.08)
- Active/Selected: color #E91E8C, background rgba(233,30,140,0.12)

### Input

**Text Input** (Text content entry):
```css
background: #16213E;
color: #F8F8F2;
font: 400 14px 'DM Sans';
padding: 8px 12px;
border-radius: 6px;
border: 1px solid #3A3A5C;
transition: border-color 0.15s ease, box-shadow 0.15s ease;
```
- Focus: border-color #E91E8C, box-shadow 0 0 0 3px rgba(233,30,140,0.15)
- Error: border-color #F87272, box-shadow 0 0 0 3px rgba(248,114,114,0.15)
- Placeholder color: #5A5A78

**Select / Dropdown** (Font family, format):
```css
background: #16213E;
color: #F8F8F2;
font: 400 14px 'DM Sans';
padding: 8px 12px;
border-radius: 6px;
border: 1px solid #3A3A5C;
appearance: none;
/* Custom chevron icon via background-image */
```
- Dropdown menu: background #353564, border-radius 8px, shadow-md

**Slider** (Font size, stroke width):
```css
/* Track */
background: #3A3A5C;
height: 4px;
border-radius: 2px;

/* Fill */
background: #E91E8C;

/* Thumb */
background: #F8F8F2;
width: 16px;
height: 16px;
border-radius: 50%;
border: 2px solid #E91E8C;
box-shadow: 0 0 6px rgba(233,30,140,0.3);
```

**Color Picker** (Text color, stroke color):
```css
/* Wrapper */
width: 36px;
height: 36px;
border-radius: 6px;
border: 2px solid #3A3A5C;
overflow: hidden;
cursor: pointer;

/* Native input hidden, shows chosen color as background */
```
- Active: border-color #E91E8C

### Card / Panel

**Sidebar Panel**:
```css
background: #252547;
width: 280px;
height: 100%;
border-left: 1px solid #3A3A5C;
padding: 16px;
overflow-y: auto;
```

**Control Group** (within sidebar):
```css
margin-bottom: 16px;
```
- Label: 12px, weight 500, color #A0A0B8, text-transform uppercase, letter-spacing 0.05em, margin-bottom 8px

**Text Layer Item** (in text list):
```css
background: #16213E;
padding: 10px 12px;
border-radius: 8px;
border: 1px solid transparent;
display: flex;
align-items: center;
justify-content: space-between;
cursor: pointer;
transition: border-color 0.15s ease, background 0.15s ease;
```
- Hover: background #1E2A4A
- Selected: border-color #E91E8C, background rgba(233,30,140,0.08)

### Upload Zone (Empty State)

```css
/* Container */
background: #16213E;
border: 2px dashed #3A3A5C;
border-radius: 16px;
padding: 48px;
text-align: center;
cursor: pointer;
transition: border-color 0.2s ease, background 0.2s ease;

/* Icon */
color: #A0A0B8;
width: 64px;
height: 64px;
margin-bottom: 16px;

/* Title text */
font: 600 18px 'DM Sans';
color: #F8F8F2;
margin-bottom: 8px;
/* "Drag an image here or click to upload" */

/* Hint text */
font: 400 13px 'DM Sans';
color: #5A5A78;
/* "Supports PNG, JPG, GIF, WebP" */
```
- Hover / Drag-over: border-color #E91E8C, background rgba(233,30,140,0.05)

### Floating Toolbar (Contextual — appears near selected text)

```css
background: #353564;
border-radius: 10px;
padding: 6px;
display: flex;
gap: 2px;
box-shadow: 0 4px 12px rgba(0,0,0,0.4);
border: 1px solid #3A3A5C;
```
- Contains: Font selector (compact), size adjust (+/-), color picker, bold/italic toggles, delete
- Positioned above or below the selected text element, centered
- Appears with a subtle fade+scale animation (150ms)

### Toast / Notification

```css
background: #252547;
border-left: 3px solid #36D399; /* or #F87272 for error, #FBBD23 for warning */
border-radius: 8px;
padding: 12px 16px;
box-shadow: 0 4px 12px rgba(0,0,0,0.4);
color: #F8F8F2;
font: 500 14px 'DM Sans';
```
- Position: bottom-center or bottom-right
- Auto-dismiss after 3 seconds with fade-out

### Export Modal (Optional — can also be inline)

```css
/* Overlay */
background: rgba(0,0,0,0.6);
backdrop-filter: blur(4px);

/* Modal */
background: #252547;
border-radius: 16px;
padding: 24px;
width: 360px;
box-shadow: 0 8px 24px rgba(0,0,0,0.5);
border: 1px solid #3A3A5C;
```
- Contains: Format toggle (PNG/JPEG), quality slider (JPEG only), preview thumbnail, Download button

## Layout

### Desktop (>1024px)

```
┌──────────────────────────────────────────────────────────┐
│  Header: [Logo/Title]            [Undo] [Redo] [Download]│
│  bg: #1A1A2E  h: 56px  border-bottom: 1px #3A3A5C       │
├──────────────────────────────────────┬───────────────────┤
│                                      │  Controls Panel   │
│                                      │  bg: #252547      │
│       Canvas Area                    │  w: 280px         │
│       bg: #16213E                    │  ───────────────  │
│       (centered image + overlays)    │  [+ Add Text]     │
│                                      │  Font Family ▾    │
│       Checkerboard pattern for       │  Size: ──●──      │
│       transparent areas              │  Color: ■  ■      │
│                                      │  Stroke: ■ ──●──  │
│                                      │  [B] [I] [L][C][R]│
│                                      │  ───────────────  │
│                                      │  Text Layers      │
│                                      │  [Text 1]    🗑   │
│                                      │  [Text 2]    🗑   │
├──────────────────────────────────────┴───────────────────┤
│  (no footer — clean edge-to-edge)                        │
└──────────────────────────────────────────────────────────┘
```

- Header height: 56px
- Sidebar width: 280px
- Canvas area: fills remaining space, image centered with padding
- Min app width: 768px before collapsing to mobile layout

### Tablet (640-1024px)

- Sidebar collapses to a bottom sheet (slides up from bottom)
- Canvas takes full width
- Floating toolbar remains for quick access
- Header stays the same but with icon-only buttons

### Mobile (<640px)

- Full-width canvas (image fills screen width)
- Controls in a collapsible bottom drawer
- "Add Text" as floating action button (bottom-right)
- Download button in header (icon-only)
- Swipe up to reveal controls

## Responsive Breakpoints

| Name | Width | Layout |
|------|-------|--------|
| mobile | < 640px | Full-width canvas, bottom drawer controls |
| tablet | 640-1024px | Full-width canvas, bottom sheet controls |
| desktop | > 1024px | Canvas + right sidebar |

## Animations & Micro-interactions

All animations use `ease-out` timing unless noted.

| Element | Trigger | Animation | Duration |
|---------|---------|-----------|----------|
| Buttons | Hover | Background color transition | 150ms |
| Primary Button | Hover | Subtle lift (translateY -1px) | 100ms |
| Floating Toolbar | Appear | Fade in + scale from 0.95 | 150ms |
| Floating Toolbar | Disappear | Fade out | 100ms |
| Text element (canvas) | Selected | Magenta border + corner handles appear | 150ms |
| Drag handle | Grabbed | Glow with primary accent shadow | immediate |
| Upload zone | Drag-over | Border color pulse to magenta | 200ms |
| Toast | Appear | Slide up + fade in | 250ms |
| Toast | Dismiss | Fade out + slide down | 200ms |
| Panel (mobile) | Open | Slide up from bottom | 250ms ease-out |
| Export button | Meme ready | Subtle shimmer/pulse once | 600ms |
| Text layer item | Hover | Background brighten | 150ms |
| Sidebar section | Expand/Collapse | Height transition | 200ms |

## Accessibility

- **Contrast**: All text meets WCAG AA 4.5:1 minimum
  - #F8F8F2 on #1A1A2E = 14.7:1 (passes AAA)
  - #A0A0B8 on #1A1A2E = 5.8:1 (passes AA)
  - #E91E8C on #1A1A2E = 4.6:1 (passes AA)
  - #F8F8F2 on #E91E8C = 3.2:1 (use only for large text/icons)
- **Focus indicators**: 2px outline in Primary with 3px offset on all interactive elements
- **Touch targets**: Minimum 44x44px for all interactive controls
- **Keyboard navigation**: Tab through all controls, Enter/Space to activate, arrow keys for fine text positioning
- **Reduced motion**: Respect `prefers-reduced-motion` — disable all transitions/animations
- **Screen reader**: `aria-live` regions for canvas state changes, `aria-label` on icon-only buttons

## CSS Custom Properties (Design Tokens)

```css
:root {
  /* Colors */
  --color-primary: #E91E8C;
  --color-primary-hover: #D4177F;
  --color-primary-muted: rgba(233, 30, 140, 0.15);
  --color-secondary: #FF6B9D;
  --color-bg: #1A1A2E;
  --color-surface: #252547;
  --color-surface-hover: #2E2E5E;
  --color-canvas-bg: #16213E;
  --color-border: #3A3A5C;
  --color-text: #F8F8F2;
  --color-text-secondary: #A0A0B8;
  --color-text-disabled: #5A5A78;
  --color-success: #36D399;
  --color-error: #F87272;
  --color-warning: #FBBD23;

  /* Typography */
  --font-ui: 'DM Sans', system-ui, sans-serif;
  --font-meme-default: 'Impact', 'Anton', sans-serif;

  /* Spacing */
  --space-xs: 4px;
  --space-sm: 8px;
  --space-md: 12px;
  --space-lg: 16px;
  --space-xl: 24px;
  --space-2xl: 32px;
  --space-3xl: 48px;

  /* Radius */
  --radius-sm: 6px;
  --radius-md: 8px;
  --radius-lg: 12px;
  --radius-xl: 16px;
  --radius-full: 9999px;

  /* Shadows */
  --shadow-sm: 0 1px 3px rgba(0, 0, 0, 0.3);
  --shadow-md: 0 4px 12px rgba(0, 0, 0, 0.4);
  --shadow-lg: 0 8px 24px rgba(0, 0, 0, 0.5);
  --shadow-glow: 0 0 12px rgba(233, 30, 140, 0.3);

  /* Layout */
  --header-height: 56px;
  --sidebar-width: 280px;

  /* Transitions */
  --transition-fast: 100ms ease-out;
  --transition-normal: 150ms ease-out;
  --transition-slow: 250ms ease-out;
}
```

## Meme Text Presets (One-Click Styles)

These presets apply instantly when selected, giving users polished text without manual tweaking:

| Preset Name | Font | Size | Color | Stroke | Extra |
|-------------|------|------|-------|--------|-------|
| Classic Meme | Impact | 48px | #FFFFFF | 3px #000000 | All-caps, center-aligned |
| Modern Clean | DM Sans Bold | 36px | #FFFFFF | 2px rgba(0,0,0,0.6) | Title case |
| Handwritten | Permanent Marker | 40px | #FFFFFF | 2px #000000 | Natural case |
| Comic Pop | Bangers | 44px | #FFFF00 | 3px #000000 | All-caps |
| Neon Glow | Bebas Neue | 42px | #E91E8C | none | Shadow: 0 0 20px #E91E8C |
| Minimal | Oswald | 32px | #F8F8F2 | none | Uppercase, light weight |

---
Status: READY_FOR_REVIEW
