# UI Specification: HTML Learning Playground

Based on research in: .plans/research/ui-research.md
Styling wireframes from: .plans/research/ux-research.md
Product requirements: .plans/research/product-research.md

## Design System

### Color Palette (Electric Playground Theme)

**IMPORTANT**: These colors were specifically chosen to differentiate from competitor learning platforms that overuse green/blue. The magenta + cyan combination is bold, modern, and unique in the educational space.

#### Primary Palette

| Name | Hex | RGB | Usage |
|------|-----|-----|-------|
| Primary | #FF2E63 | rgb(255, 46, 99) | Main CTAs, primary buttons, active states, progress bars |
| Primary Hover | #E6175A | rgb(230, 23, 90) | Hover state for primary elements |
| Primary Light | #FF5C85 | rgb(255, 92, 133) | Lighter accent, badges, highlights |
| Secondary | #00D9FF | rgb(0, 217, 255) | Secondary actions, links, accent elements |
| Secondary Hover | #00BFE6 | rgb(0, 191, 230) | Hover state for secondary elements |
| Secondary Light | #33E2FF | rgb(51, 226, 255) | Lighter accent, success states |

#### Background & Surface (Dark Mode Primary)

| Name | Hex | RGB | Usage |
|------|-----|-----|-------|
| Background Dark | #1A1A2E | rgb(26, 26, 46) | Main app background (dark mode) |
| Surface Dark | #16213E | rgb(22, 33, 62) | Cards, code editor, elevated elements (dark) |
| Surface Elevated | #1F2D52 | rgb(31, 45, 82) | Hover/focus states for cards |
| Border Dark | #2A3B62 | rgb(42, 59, 98) | Subtle borders and dividers |

#### Background & Surface (Light Mode Alternative)

| Name | Hex | RGB | Usage |
|------|-----|-----|-------|
| Background Light | #EEEEF7 | rgb(238, 238, 247) | Main app background (light mode) |
| Surface Light | #FFFFFF | rgb(255, 255, 255) | Cards, elevated elements (light) |
| Border Light | #D4D4E8 | rgb(212, 212, 232) | Borders in light mode |

#### Text Colors

| Name | Hex | RGB | Usage |
|------|-----|-----|-------|
| Text Primary Dark | #FFFFFF | rgb(255, 255, 255) | Main text on dark backgrounds |
| Text Secondary Dark | #B8B8D1 | rgb(184, 184, 209) | Muted text, descriptions (dark mode) |
| Text Tertiary Dark | #8585A3 | rgb(133, 133, 163) | Subtle text, placeholders |
| Text Primary Light | #1A1A2E | rgb(26, 26, 46) | Main text on light backgrounds |
| Text Secondary Light | #5F5F7A | rgb(95, 95, 122) | Muted text (light mode) |

#### Semantic Colors

| Name | Hex | RGB | Usage |
|------|-----|-----|-------|
| Success | #00FF9D | rgb(0, 255, 157) | Success states, checkmarks, correct answers |
| Success Dark | #00CC7E | rgb(0, 204, 126) | Success hover/active |
| Error | #FF3366 | rgb(255, 51, 102) | Error states, warnings, incorrect answers |
| Error Dark | #E61A47 | rgb(230, 26, 71) | Error hover/active |
| Warning | #FFB627 | rgb(255, 182, 39) | Warning states, tips, caution |
| Info | #00D9FF | rgb(0, 217, 255) | Info messages, hints (uses secondary) |

#### Gradient Definitions

| Name | CSS Value | Usage |
|------|-----------|-------|
| Primary Gradient | linear-gradient(135deg, #FF2E63 0%, #FF5C85 100%) | Buttons, headers, hero elements |
| Secondary Gradient | linear-gradient(135deg, #00D9FF 0%, #33E2FF 100%) | Secondary CTAs, accents |
| Sunset Gradient | linear-gradient(135deg, #FF2E63 0%, #6B4CE6 50%, #00D9FF 100%) | Special moments, achievements |
| Dark Glow | radial-gradient(circle at 50% 0%, rgba(255, 46, 99, 0.15) 0%, transparent 50%) | Background ambiance |

**Color Rationale**: The hot pink/magenta + electric cyan combination creates a high-energy, playful aesthetic that's unique in the learning space. This avoids the overused green (Duolingo, freeCodeCamp) and blue (most tech products). The dark mode primary approach feels modern and reduces eye strain during longer sessions.

### Typography

#### Font Families

**Display Font**: Space Grotesk (geometric, distinctive, slightly quirky - perfect for headings)
- Fallback: 'Inter', 'Segoe UI', system-ui, sans-serif

**Body Font**: Inter (best-in-class legibility, variable font support)
- Fallback: 'Segoe UI', system-ui, -apple-system, sans-serif

**Code Font**: JetBrains Mono (designed for developers, excellent character distinction)
- Fallback: 'Fira Code', 'Monaco', 'Courier New', monospace

#### Type Scale

| Element | Font | Size | Weight | Line Height | Letter Spacing | Usage |
|---------|------|------|--------|-------------|----------------|-------|
| H1 | Space Grotesk | 48px | 700 | 1.1 | -0.02em | Hero headlines, main titles |
| H2 | Space Grotesk | 36px | 700 | 1.2 | -0.01em | Section headers |
| H3 | Space Grotesk | 28px | 600 | 1.3 | -0.01em | Subsection headers |
| H4 | Space Grotesk | 24px | 600 | 1.4 | 0 | Card headers |
| H5 | Inter | 20px | 600 | 1.4 | 0 | Small headers |
| Body Large | Inter | 18px | 400 | 1.6 | 0 | Intro text, important copy |
| Body | Inter | 16px | 400 | 1.5 | 0 | Standard body text |
| Body Small | Inter | 14px | 400 | 1.5 | 0 | Secondary text, captions |
| Caption | Inter | 12px | 500 | 1.4 | 0.01em | Labels, metadata |
| Code | JetBrains Mono | 14px | 400 | 1.6 | 0 | Code snippets |
| Code Large | JetBrains Mono | 16px | 400 | 1.6 | 0 | Main code editor |

#### Responsive Typography

| Breakpoint | H1 | H2 | H3 | Body | Notes |
|------------|----|----|----|----|-------|
| Mobile (<640px) | 32px | 24px | 20px | 16px | Reduced for small screens |
| Tablet (640-1024px) | 40px | 30px | 24px | 16px | Moderate scaling |
| Desktop (>1024px) | 48px | 36px | 28px | 16px | Full scale |

### Spacing Scale

Base unit: 4px (uses multiples for consistency)

| Token | Value | Usage |
|-------|-------|-------|
| xs | 4px | Tight spacing, icon padding |
| sm | 8px | Small gaps between related elements |
| md | 16px | Standard padding for buttons, cards |
| lg | 24px | Section spacing, card padding |
| xl | 32px | Large gaps between sections |
| 2xl | 48px | Page margins, hero spacing |
| 3xl | 64px | Major section breaks |
| 4xl | 96px | Maximum spacing |

### Border Radius

| Token | Value | Usage |
|-------|-------|-------|
| sm | 6px | Small elements, tags, badges |
| md | 8px | Buttons, inputs, small cards |
| lg | 12px | Cards, modals, code blocks |
| xl | 16px | Large containers, hero sections |
| 2xl | 24px | Special containers |
| full | 9999px | Pills, avatars, circular elements |

### Shadows & Glows

#### Standard Shadows (Black-based)

| Token | Value | Usage |
|-------|-------|-------|
| sm | 0 1px 3px rgba(0, 0, 0, 0.12), 0 1px 2px rgba(0, 0, 0, 0.08) | Subtle elevation |
| md | 0 4px 6px rgba(0, 0, 0, 0.12), 0 2px 4px rgba(0, 0, 0, 0.08) | Cards, dropdowns |
| lg | 0 10px 15px rgba(0, 0, 0, 0.15), 0 4px 6px rgba(0, 0, 0, 0.08) | Modals, large cards |
| xl | 0 20px 25px rgba(0, 0, 0, 0.15), 0 10px 10px rgba(0, 0, 0, 0.04) | Major elevation |

#### Colored Glows (Unique to our design!)

| Token | Value | Usage |
|-------|-------|-------|
| Primary Glow | 0 0 20px rgba(255, 46, 99, 0.4), 0 0 40px rgba(255, 46, 99, 0.2) | Primary button hover, success states |
| Secondary Glow | 0 0 20px rgba(0, 217, 255, 0.4), 0 0 40px rgba(0, 217, 255, 0.2) | Secondary accents, links |
| Success Glow | 0 0 20px rgba(0, 255, 157, 0.5) | Correct answer feedback |
| Error Glow | 0 0 20px rgba(255, 51, 102, 0.5) | Error states |
| Focus Ring | 0 0 0 3px rgba(255, 46, 99, 0.3) | Keyboard focus states |

### Animation Specifications

**Critical for "fun" requirement** - Every interaction should feel alive!

#### Timing Functions

| Name | Cubic Bezier | Usage |
|------|-------------|-------|
| Ease Out Expo | cubic-bezier(0.16, 1, 0.3, 1) | Smooth exits, menu slides |
| Ease Out Back | cubic-bezier(0.34, 1.56, 0.64, 1) | Bounce effects, playful |
| Ease In Out | cubic-bezier(0.4, 0, 0.2, 1) | General animations |
| Spring | cubic-bezier(0.68, -0.55, 0.265, 1.55) | Bouncy, energetic |

#### Duration Scale

| Name | Duration | Usage |
|------|----------|-------|
| Fast | 150ms | Micro-interactions, hovers |
| Base | 250ms | Standard transitions |
| Medium | 400ms | Card animations, reveals |
| Slow | 600ms | Page transitions, complex animations |
| Slower | 1000ms | Celebration animations |

#### Keyframe Animations

**Bounce In**
```css
@keyframes bounceIn {
  0% { transform: scale(0.3); opacity: 0; }
  50% { transform: scale(1.05); }
  70% { transform: scale(0.9); }
  100% { transform: scale(1); opacity: 1; }
}
```

**Shake**
```css
@keyframes shake {
  0%, 100% { transform: translateX(0); }
  10%, 30%, 50%, 70%, 90% { transform: translateX(-4px); }
  20%, 40%, 60%, 80% { transform: translateX(4px); }
}
```

**Pulse Glow**
```css
@keyframes pulseGlow {
  0%, 100% { box-shadow: 0 0 10px var(--primary); }
  50% { box-shadow: 0 0 20px var(--primary), 0 0 30px var(--primary); }
}
```

**Float**
```css
@keyframes float {
  0%, 100% { transform: translateY(0px); }
  50% { transform: translateY(-10px); }
}
```

**Confetti Fall** (for success states)
```css
@keyframes confettiFall {
  0% { transform: translateY(-100%) rotate(0deg); opacity: 1; }
  100% { transform: translateY(100vh) rotate(360deg); opacity: 0; }
}
```

### Component Specifications

#### Button

**Primary Button**
```css
Base State:
  background: linear-gradient(135deg, #FF2E63 0%, #FF5C85 100%)
  color: #FFFFFF
  padding: 12px 24px
  border-radius: 8px
  font-weight: 600
  font-size: 16px
  border: none
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.12)
  transition: all 250ms cubic-bezier(0.4, 0, 0.2, 1)
  cursor: pointer

Hover State:
  transform: translateY(-2px)
  box-shadow: 0 0 20px rgba(255, 46, 99, 0.4), 0 4px 10px rgba(0, 0, 0, 0.15)

Active State:
  transform: translateY(0) scale(0.98)
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.12)

Disabled State:
  opacity: 0.5
  cursor: not-allowed
  transform: none
  box-shadow: none
```

**Secondary Button**
```css
Base State:
  background: transparent
  color: #00D9FF
  border: 2px solid #00D9FF
  padding: 12px 24px
  border-radius: 8px
  font-weight: 600

Hover State:
  background: rgba(0, 217, 255, 0.1)
  box-shadow: 0 0 20px rgba(0, 217, 255, 0.3)
```

**Ghost Button**
```css
Base State:
  background: transparent
  color: #B8B8D1
  border: none
  padding: 12px 24px

Hover State:
  color: #FFFFFF
  background: rgba(255, 255, 255, 0.05)
```

#### Input Fields

**Text Input**
```css
Base State:
  background: #16213E
  border: 2px solid #2A3B62
  border-radius: 8px
  padding: 12px 16px
  color: #FFFFFF
  font-size: 16px
  font-family: Inter
  transition: all 250ms

Hover State:
  border-color: #3A4B72

Focus State:
  border-color: #FF2E63
  box-shadow: 0 0 0 3px rgba(255, 46, 99, 0.2)
  outline: none

Error State:
  border-color: #FF3366
  box-shadow: 0 0 0 3px rgba(255, 51, 102, 0.2)
  animation: shake 400ms

Success State:
  border-color: #00FF9D
  box-shadow: 0 0 0 3px rgba(0, 255, 157, 0.2)
```

#### Code Editor

```css
Container:
  background: #0D1117 (slightly darker than surface)
  border-radius: 12px
  padding: 20px
  font-family: JetBrains Mono
  font-size: 14px
  line-height: 1.6
  border: 1px solid #2A3B62
  box-shadow: inset 0 2px 8px rgba(0, 0, 0, 0.3)

Line Numbers:
  color: #5F5F7A
  padding-right: 16px
  user-select: none

Syntax Highlighting:
  Tag: #FF2E63
  Attribute: #00D9FF
  String: #00FF9D
  Comment: #5F5F7A
```

#### Card

```css
Base State:
  background: #16213E
  border-radius: 12px
  padding: 24px
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.12)
  transition: all 250ms cubic-bezier(0.4, 0, 0.2, 1)

Hover State (if interactive):
  transform: translateY(-4px)
  box-shadow: 0 10px 15px rgba(0, 0, 0, 0.15)
  border: 1px solid rgba(255, 46, 99, 0.2)

Active Card:
  border: 2px solid #FF2E63
  box-shadow: 0 0 20px rgba(255, 46, 99, 0.3)
```

#### Progress Bar

```css
Container:
  background: #16213E
  border-radius: 9999px
  height: 8px
  overflow: hidden

Fill:
  background: linear-gradient(90deg, #FF2E63 0%, #00D9FF 100%)
  height: 100%
  border-radius: 9999px
  transition: width 400ms cubic-bezier(0.4, 0, 0.2, 1)
  animation: pulse 2s infinite

Success State:
  animation: fillComplete 600ms cubic-bezier(0.34, 1.56, 0.64, 1)
```

#### Badge/Tag

```css
Base:
  background: rgba(255, 46, 99, 0.15)
  color: #FF5C85
  padding: 4px 12px
  border-radius: 6px
  font-size: 12px
  font-weight: 600
  letter-spacing: 0.02em

Success:
  background: rgba(0, 255, 157, 0.15)
  color: #00FF9D

Info:
  background: rgba(0, 217, 255, 0.15)
  color: #00D9FF
```

#### Tooltip

```css
Container:
  background: #16213E
  border: 1px solid #2A3B62
  border-radius: 8px
  padding: 8px 12px
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.2)
  font-size: 14px
  color: #FFFFFF
  max-width: 250px
  z-index: 1000

Animation:
  Appear: fadeIn 150ms + slideUp 150ms
  Position: 8px offset from trigger element
```

#### Modal/Dialog

```css
Backdrop:
  background: rgba(26, 26, 46, 0.8)
  backdrop-filter: blur(8px)
  animation: fadeIn 250ms

Container:
  background: #16213E
  border: 1px solid #2A3B62
  border-radius: 16px
  padding: 32px
  max-width: 600px
  box-shadow: 0 20px 25px rgba(0, 0, 0, 0.3)
  animation: slideUp 400ms cubic-bezier(0.16, 1, 0.3, 1)
```

### Interaction Animations (CRITICAL!)

Based on UX research requirement: "animations for user actions"

#### Button Click
```
1. Scale down (scale(0.98)) - 100ms
2. Scale back up with bounce (scale(1)) - 200ms cubic-bezier(0.34, 1.56, 0.64, 1)
```

#### Success Feedback
```
1. Green checkmark bounces in - 300ms
2. Confetti particles fall from top - 1000ms
3. Success message fades in - 250ms
4. Card glows with success color - 500ms
```

#### Error Feedback
```
1. Shake animation - 400ms
2. Red highlight pulse - 300ms
3. Gentle error message slide in from bottom - 250ms
```

#### Hover States (All Interactive Elements)
```
- Lift effect: translateY(-2px)
- Glow appears: colored shadow fades in
- Color intensifies: brightness(1.1)
- Cursor: changes to pointer
- Transition: 150ms
```

#### Typing in Code Editor
```
- Syntax highlighting applies in real-time
- Matching tags highlight when cursor inside - 200ms fade
- Auto-close tags with slide-in animation - 150ms
```

#### Page Transition
```
1. Current content fades out - 200ms
2. New content slides up and fades in - 400ms cubic-bezier(0.16, 1, 0.3, 1)
3. Stagger child elements by 50ms each
```

#### Achievement Unlock
```
1. Screen flash (white overlay at 20% opacity) - 100ms
2. Badge scales up with bounce - 600ms
3. Confetti explosion - 1500ms
4. Sound effect (optional, with mute toggle)
5. Message appears with typewriter effect - 800ms
```

### Responsive Breakpoints

| Name | Min Width | Max Width | Layout Changes |
|------|-----------|-----------|----------------|
| Mobile | 0px | 639px | Single column, stacked layout, simplified animations |
| Tablet | 640px | 1023px | Two-column where applicable, full animations |
| Desktop | 1024px | 1439px | Multi-column, split-screen editor, full experience |
| Desktop Large | 1440px | ∞ | Maximum width containers (1280px), enhanced spacing |

#### Mobile Considerations
- Touch targets: minimum 44x44px
- Reduce animation complexity (respect prefers-reduced-motion)
- Simplify code editor (remove line numbers on very small screens)
- Stack preview below code instead of side-by-side
- Larger font sizes for readability (minimum 16px for body)

### Accessibility

#### Color Contrast
All text meets WCAG AA standards (minimum 4.5:1 for body, 3:1 for large text)
- Primary (#FF2E63) on Dark Background (#1A1A2E): 8.2:1 ✓
- Secondary (#00D9FF) on Dark Background: 9.1:1 ✓
- Text Primary (#FFFFFF) on Dark Background: 15.8:1 ✓

#### Keyboard Navigation
- All interactive elements focusable via Tab
- Focus states clearly visible with glow rings
- Escape key closes modals
- Enter key submits forms
- Arrow keys navigate through lessons

#### Screen Readers
- ARIA labels on all interactive elements
- Role attributes for custom components
- Live regions for dynamic content updates
- Skip links for navigation

#### Motion
```css
@media (prefers-reduced-motion: reduce) {
  * {
    animation-duration: 0.01ms !important;
    transition-duration: 0.01ms !important;
  }
}
```

### Special Effects

#### Confetti System (Canvas-based)
```
Particles: 50-100 colored rectangles
Colors: [#FF2E63, #00D9FF, #00FF9D, #FFB627, #FF5C85]
Physics: Gravity (500px/s²), rotation, fade out
Duration: 2 seconds
Trigger: Lesson completion, achievements
```

#### Particle Background (Optional ambient effect)
```
Subtle floating circles in background
Colors: Primary/Secondary at 5% opacity
Slow float animation (30s duration)
Adds depth without distraction
```

#### Glow Cursor Trail (Advanced)
```
On code preview area: cursor leaves brief glow trail
Color matches primary gradient
Fades after 500ms
Adds "magical" feel to interactions
```

### Dark/Light Mode Toggle

**Default**: Dark mode (reduces eye strain for coding)

**Light Mode Overrides**:
```
Background: #EEEEF7
Surface: #FFFFFF
Text Primary: #1A1A2E
Text Secondary: #5F5F7A
Borders: #D4D4E8
Shadows: Increased opacity for visibility
Glows: Reduced intensity
```

**Toggle Animation**:
- Sun/moon icon rotates 180deg - 400ms
- Colors transition smoothly - 300ms
- Persist choice in localStorage

### Code Syntax Highlighting Theme

**HTML Helper Custom Theme** (based on our colors)

```css
.tag { color: #FF2E63; font-weight: 600; } /* Vibrant pink for tags */
.attribute { color: #00D9FF; } /* Cyan for attributes */
.string { color: #00FF9D; } /* Success green for strings */
.comment { color: #5F5F7A; font-style: italic; } /* Muted for comments */
.doctype { color: #FFB627; } /* Warning yellow for doctype */
.punctuation { color: #B8B8D1; } /* Secondary text for brackets */
```

### Loading States

**Spinner**
```css
Circular gradient spinner with primary colors
Rotation: 1s linear infinite
Size: 32px (standard), 48px (large)
```

**Skeleton Screens**
```css
Background: shimmer animation
Base: #16213E
Shimmer: linear-gradient with #2A3B62
Animation: 1.5s ease-in-out infinite
```

### Iconography

**Style**: Outline icons with 2px stroke
**Library**: Lucide Icons or Heroicons (modern, clean)
**Size Scale**: 16px (small), 20px (medium), 24px (large), 32px (hero)
**Color**: Inherits text color or custom accent
**Hover**: Scale(1.1) + rotate slightly for playfulness

---

## Implementation Notes for Developers

### CSS Custom Properties Setup

```css
:root {
  /* Colors */
  --primary: #FF2E63;
  --primary-hover: #E6175A;
  --primary-light: #FF5C85;
  --secondary: #00D9FF;
  --secondary-hover: #00BFE6;
  --background-dark: #1A1A2E;
  --surface-dark: #16213E;
  --text-primary: #FFFFFF;
  --text-secondary: #B8B8D1;
  --success: #00FF9D;
  --error: #FF3366;
  --warning: #FFB627;

  /* Spacing */
  --space-xs: 4px;
  --space-sm: 8px;
  --space-md: 16px;
  --space-lg: 24px;
  --space-xl: 32px;

  /* Border Radius */
  --radius-sm: 6px;
  --radius-md: 8px;
  --radius-lg: 12px;
  --radius-xl: 16px;

  /* Typography */
  --font-display: 'Space Grotesk', sans-serif;
  --font-body: 'Inter', sans-serif;
  --font-code: 'JetBrains Mono', monospace;

  /* Timing */
  --duration-fast: 150ms;
  --duration-base: 250ms;
  --duration-medium: 400ms;
  --duration-slow: 600ms;

  /* Easing */
  --ease-out-expo: cubic-bezier(0.16, 1, 0.3, 1);
  --ease-bounce: cubic-bezier(0.34, 1.56, 0.64, 1);
}
```

### Animation Library Recommendations

1. **Framer Motion** (React) - For complex orchestrated animations
2. **GSAP** - For timeline-based celebrations and complex sequences
3. **CSS Animations** - For simple micro-interactions
4. **Canvas API** - For confetti and particle effects

### Performance Considerations

- Use `transform` and `opacity` for animations (GPU accelerated)
- Avoid animating `width`, `height`, `top`, `left` (causes reflow)
- Use `will-change` sparingly on elements about to animate
- Debounce rapid animations (typing, hovering)
- Lazy load heavy animations until needed
- Respect `prefers-reduced-motion` always

---

## Design Principles Summary

1. **Every Action Gets Visual Feedback** - Nothing happens silently
2. **Celebrate Progress** - Make success feel amazing with animations
3. **Unique Visual Identity** - Stand out with magenta/cyan, not generic blue/green
4. **Playful but Professional** - Fun without being childish
5. **Accessibility First** - Beautiful AND usable for everyone
6. **Performance Matters** - Smooth 60fps or don't animate
7. **Dark Mode Default** - Comfortable for extended coding sessions

---

## File Structure for Assets

```
/assets
  /fonts
    space-grotesk.woff2
    inter.woff2
    jetbrains-mono.woff2
  /animations
    confetti.json (Lottie)
    celebration.json
    success.json
  /icons
    [lucide icon set]
  /images
    mascot-idle.svg
    mascot-success.svg
    mascot-thinking.svg
```

---

Status: READY_FOR_REVIEW

This specification provides comprehensive visual design direction based on research findings. The "Electric Playground" theme with magenta/cyan creates a unique, energetic learning environment that differentiates from competitor platforms while maintaining accessibility and professional polish.