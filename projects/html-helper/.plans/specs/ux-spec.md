# UX Specification: Interactive HTML Learning Platform

Based on research in: .plans/research/ux-research.md, .plans/research/product-research.md, .plans/research/ui-research.md

## Executive Summary

Creating an HTML learning platform that feels like **playing a game, not taking a course**. Every user interaction triggers delightful animations and instant visual feedback. Users learn HTML by seeing their code come alive through animated characters and visual effects.

**Core Philosophy**: Show, don't tell. Do, don't read. Play, don't study.

---

## Screen Architecture

### Primary Screens
1. **Landing/Demo Screen** - Immediate interactive demo (no text walls)
2. **Main Learning Interface** - Split-screen code editor + live animated preview
3. **Lesson Map** - Visual progression path with unlocked/locked lessons
4. **Achievement Gallery** - Collection of earned badges and milestones
5. **Settings/Preferences** - Theme toggle, sound controls, keyboard shortcuts

### Screen Flow Priority
Landing → Main Learning Interface (80% of time spent here) → Lesson Map → Back to Learning

---

## Detailed Wireframes

### Screen 1: Landing/Demo Screen
```
┌─────────────────────────────────────────────────────────────┐
│  [Logo: HTML Helper]                    [Start Learning →]  │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│           ┌─────────────────────────────────┐                │
│           │  TRY IT NOW                     │                │
│           │  ┌────────────┐  ┌────────────┐ │                │
│           │  │ Type here: │  │ [Preview]  │ │                │
│           │  │            │  │            │ │                │
│           │  │ <h1>Hi!</  │  │    Hi!     │ │                │
│           │  │  [cursor]  │  │  [bounces] │ │                │
│           │  │            │  │            │ │                │
│           │  └────────────┘  └────────────┘ │                │
│           │  ← Type HTML    See it animate! │                │
│           └─────────────────────────────────┘                │
│                                                               │
│   "Learn HTML by playing with it - every tag you type        │
│    makes something magical happen"                           │
│                                                               │
│   [No signup needed • Works in browser • Free forever]       │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

**Key Elements**:
- **Interactive Demo Box**: Live mini-editor where visitors can type immediately
- **Real-time Animation**: As they type `<h1>`, text appears and animates in preview
- **Zero Friction**: No "Learn More" buttons, no feature lists - just interactive fun
- **Clear CTA**: "Start Learning" button glows/pulses subtly

**Animations**:
- Demo preview updates on every keystroke
- Text bounces/scales when valid HTML is typed
- Confetti burst when user completes a valid tag
- Mascot character peeks from corner, reacting to typing

---

### Screen 2: Main Learning Interface (Primary Screen)
```
┌──────────────────────────────────────────────────────────────────┐
│ [≡ Menu] Lesson 3: Headings          [Progress: 12/50] ⭐ 125 XP │
│ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │
├──────────────────────────┬───────────────────────────────────────┤
│                          │                                       │
│  CHALLENGE               │           PREVIEW                     │
│  ──────────              │                                       │
│  Make the robot say      │         ┌──────────┐                 │
│  "Hello World" in        │         │  🤖      │                 │
│  big bold text!          │         │          │                 │
│                          │         │  [speech │                 │
│  ┌────────────────────┐  │         │   bubble]│                 │
│  │ HINT (tap to see)  │  │         └──────────┘                 │
│  └────────────────────┘  │                                       │
│                          │                                       │
│  YOUR CODE:              │                                       │
│  ┌────────────────────┐  │    Updates live as you type!         │
│  │ <h1>Hello World    │  │                                       │
│  │  [cursor blinks]   │  │                                       │
│  │                    │  │                                       │
│  └────────────────────┘  │                                       │
│                          │                                       │
│  [Tag Palette]           │                                       │
│  <h1> <p> <button>       │                                       │
│  <img> <a> <div>         │                                       │
│                          │                                       │
│  [← Previous]  [Next →]  │         [🔊 Sound On]  [🌙 Theme]    │
└──────────────────────────┴───────────────────────────────────────┘
```

**Key Elements**:

**Left Panel (Challenge + Code)**:
- **Challenge Description**: Written conversationally ("Make the robot wave" not "Implement a heading element")
- **Collapsible Hint System**: Hints appear one at a time, progressive disclosure
- **Code Editor**: Syntax-highlighted, auto-closing tags, real-time validation
- **Tag Palette**: Visual buttons for common tags (especially useful on mobile/tablets)
- **Navigation**: Always accessible to move between lessons

**Right Panel (Live Preview)**:
- **Animated Character/Scene**: Robot, animals, or abstract shapes that respond to HTML
- **Real-time Updates**: Changes render immediately (no "Run" button needed)
- **Visual Feedback**: Elements fade/bounce/slide in as tags are typed
- **Success State**: When challenge complete, preview shows celebration animation

**Top Bar**:
- **Progress Indicator**: Shows lesson X of Y with visual progress bar
- **XP Counter**: Gamification element with animation on XP gain
- **Menu**: Access lesson map, achievements, settings

**Animations**:
- Code editor: Cursor blinks, tag suggestions appear with fade-in
- Preview: Elements scale in when HTML is valid, shake gently if invalid
- Character reactions: Watches while typing, celebrates on success, looks confused on errors
- Success moment: Confetti particles fall, screen flashes green briefly, robot dances, "Next" button pulses

---

### Screen 3: Lesson Map
```
┌─────────────────────────────────────────────────────────────┐
│  [← Back to Learning]         Your Learning Path            │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│   Section 1: HTML Basics  ✓ Complete                         │
│   ┌────┐   ┌────┐   ┌────┐   ┌────┐   ┌────┐               │
│   │ ✓  │───│ ✓  │───│ ✓  │───│ ✓  │───│ ✓  │               │
│   │ 1  │   │ 2  │   │ 3  │   │ 4  │   │ 5  │               │
│   └────┘   └────┘   └────┘   └────┘   └────┘               │
│    Tags     Head    Parag    Links    Images                │
│                                                               │
│   Section 2: Text Formatting  [In Progress]                  │
│   ┌────┐   ┌────┐   ┌────┐   ┌────┐   ┌────┐               │
│   │ ✓  │───│ 🔥 │───│ 🔒 │───│ 🔒 │───│ 🔒 │               │
│   │ 6  │   │ 7  │   │ 8  │   │ 9  │   │ 10 │               │
│   └────┘   └────┘   └────┘   └────┘   └────┘               │
│    Bold    Italic   Lists    Tables   Quote                 │
│                      ↑                                        │
│                   [Continue]                                 │
│                                                               │
│   Section 3: Interactive Elements  [Locked]                  │
│   ┌────┐   ┌────┐   ┌────┐   ┌────┐                        │
│   │ 🔒 │───│ 🔒 │───│ 🔒 │───│ 🔒 │                        │
│   │ 11 │   │ 12 │   │ 13 │   │ 14 │                        │
│   └────┘   └────┘   └────┘   └────┘                        │
│   Forms   Buttons  Inputs   Submit                          │
│                                                               │
│   [🎯 25 XP]  [⚡ 3-day streak]  [🏆 5 badges earned]        │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

**Key Elements**:
- **Visual Path**: Lessons connected by lines showing progression
- **Status Icons**:
  - ✓ = Completed (green)
  - 🔥 = Current lesson (pulsing animation)
  - 🔒 = Locked (grayscale)
- **Section Grouping**: Related lessons grouped under clear headers
- **Progress Stats**: XP, streaks, badges prominently displayed
- **Quick Resume**: "Continue" button on current lesson

**Interactions**:
- Hover over lesson: Shows title + XP reward in tooltip
- Click completed lesson: Can replay for practice
- Click locked lesson: Gentle shake + tooltip "Complete previous lessons first"
- Completed lessons have subtle check animation on load

**Animations**:
- Map loads with stagger animation (lessons appear one by one)
- Current lesson has subtle glow/pulse
- Progress lines draw in sequentially
- Stats counter animates up from 0 on page load

---

### Screen 4: Achievement Gallery
```
┌─────────────────────────────────────────────────────────────┐
│  [← Back]                 Your Achievements                  │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│   ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│   │   🏆     │  │   ⭐     │  │   💪     │  │   🚀     │   │
│   │          │  │          │  │          │  │          │   │
│   │ Tag      │  │ Tag      │  │ Link     │  │ Speed    │   │
│   │ Master   │  │ Closer   │  │ Legend   │  │ Typer    │   │
│   │          │  │          │  │          │  │          │   │
│   │ Created  │  │ Closed   │  │ Made 10  │  │ Finished │   │
│   │ your     │  │ 25 tags  │  │ working  │  │ lesson   │   │
│   │ first    │  │ correctly│  │ links    │  │ in 60s   │   │
│   │ HTML tag │  │          │  │          │  │          │   │
│   └──────────┘  └──────────┘  └──────────┘  └──────────┘   │
│                                                               │
│   ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│   │   🎨     │  │   📱     │  │   🔥     │  │   💎     │   │
│   │          │  │          │  │          │  │  (gray)  │   │
│   │ Style    │  │ Mobile   │  │ Streak   │  │ Master   │   │
│   │ Star     │  │ Master   │  │ Keeper   │  │ Builder  │   │
│   │          │  │          │  │          │  │          │   │
│   │ Added    │  │ Made a   │  │ 7 days   │  │ Complete │   │
│   │ your     │  │ page     │  │ in a row │  │ all 50   │   │
│   │ first    │  │ mobile   │  │          │  │ lessons  │   │
│   │ inline   │  │ friendly │  │          │  │ [LOCKED] │   │
│   │ style    │  │          │  │          │  │          │   │
│   └──────────┘  └──────────┘  └──────────┘  └──────────┘   │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

**Key Elements**:
- **Badge Cards**: Visual representation with icon, name, description
- **Earned vs Locked**: Colored badges earned, grayscale for locked ones
- **Progress Hints**: Locked badges show what's needed to unlock
- **Celebration on Entry**: Recently earned badges have glow effect

**Animations**:
- Badges load with stagger + scale animation
- Hover on earned badge: Slight lift + glow increase
- Click earned badge: Flips to show earned date + stats
- Locked badges: Gentle pulse to encourage progress

---

### Screen 5: Settings/Preferences
```
┌─────────────────────────────────────────────────────────────┐
│  [← Back]                    Settings                        │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│   Appearance                                                  │
│   ┌─────────────────────────────────────────────────────┐   │
│   │  Theme:  ( ) Light  (•) Dark  ( ) Auto              │   │
│   │                                                       │   │
│   │  Preview:  ┌─────────┐                              │   │
│   │           │ Sample   │  [Changes live]              │   │
│   │           │ preview  │                              │   │
│   │           └─────────┘                              │   │
│   └─────────────────────────────────────────────────────┘   │
│                                                               │
│   Sound & Animations                                          │
│   ┌─────────────────────────────────────────────────────┐   │
│   │  Sound Effects:      [ON]  OFF                       │   │
│   │  Animations:         [ON]  OFF                       │   │
│   │  Reduced Motion:     ON   [OFF]                      │   │
│   │                                                       │   │
│   │  ⚠️ Turning off animations reduces the fun!          │   │
│   └─────────────────────────────────────────────────────┘   │
│                                                               │
│   Code Editor                                                 │
│   ┌─────────────────────────────────────────────────────┐   │
│   │  Font Size:   [14px ▼]                              │   │
│   │  Auto-close tags:    [ON]  OFF                       │   │
│   │  Code hints:         [ON]  OFF                       │   │
│   └─────────────────────────────────────────────────────┘   │
│                                                               │
│   Data                                                        │
│   ┌─────────────────────────────────────────────────────┐   │
│   │  Progress saved locally (no account needed)          │   │
│   │                                                       │   │
│   │  [Export Progress]  [Reset All Progress]            │   │
│   └─────────────────────────────────────────────────────┘   │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

**Key Elements**:
- **Live Preview**: Theme changes show immediately
- **Accessibility First**: Reduced motion option prominent
- **Clear Warning**: Explain trade-offs (turning off animations)
- **Data Transparency**: Explain local storage approach

---

## User Flows

### Primary Flow: First-Time User Learning Journey

1. **Land on site** → sees interactive demo (no text, just typing box + preview)
2. **Types in demo** → sees immediate animation (text bounces, confetti on complete tag)
3. **Clicks "Start Learning"** → goes directly to Lesson 1 (no signup form!)
4. **Sees Challenge 1** → "Make your first HTML tag: Type `<h1>Hello</h1>`"
5. **Starts typing `<h1>`** → preview shows character preparing to speak
6. **Types `Hello`** → character says "Hello" in speech bubble
7. **Types `</h1>`** → character celebrates, confetti falls, XP counter animates "+5 XP"
8. **"Next" button pulses** → user clicks, slides to Lesson 2
9. **Completes 3 lessons** → achievement badge appears "Tag Master!" with celebration
10. **Clicks Menu** → sees lesson map with progress, feels motivated to continue

**Flow Duration**: Should feel fast - entire first lesson completable in 30-60 seconds.

---

### Secondary Flow: Returning User

1. **Opens site** → automatically goes to last lesson (saved in localStorage)
2. **Sees "Welcome back!" message** → with streak counter if applicable
3. **Continues from where left off** → no navigation needed
4. **Completes lesson** → celebration + auto-advances OR shows lesson map
5. **Checks achievements** → clicks menu → achievement gallery → sees new badge earned

---

### Tertiary Flow: Struggling User (Error State)

1. **User types incorrect HTML** → `<h1>Hello` (forgot closing tag)
2. **Preview shows partial state** → character looks confused (not angry/red X)
3. **After 3 seconds** → hint bubble appears: "Did you close the tag? Try `</h1>`"
4. **User hovers over hint icon** → tooltip expands with example
5. **User clicks hint** → correct code briefly highlights in editor
6. **User fixes code** → character immediately cheers, positive reinforcement

**Key Principle**: Never punish mistakes. Guide gently. Make errors feel like puzzles, not failures.

---

## Interaction Patterns

### Core Interactions Table

| Action | Element | Response | Animation |
|--------|---------|----------|-----------|
| **Type in editor** | Code input | Live preview updates | Preview elements fade in/scale |
| **Complete valid tag** | Code input | Success feedback | Confetti burst + character dance + XP animation |
| **Hover button** | Any button | Visual lift | Scale 1.05 + shadow increase + glow |
| **Click button** | Any button | Confirm action | Squish (scale 0.95) then bounce back |
| **Click "Next Lesson"** | Navigation | Load next challenge | Slide left transition + fade |
| **Open hint** | Hint button | Reveal help text | Slide down + fade in, gentle bounce |
| **Hover lesson (map)** | Lesson node | Show details | Tooltip fade in + lesson glow |
| **Click locked lesson** | Lesson node | Can't access yet | Gentle shake + "Complete previous" tooltip |
| **Earn achievement** | System trigger | Show badge | Modal fade in + badge scale/rotate + confetti |
| **Type incorrect HTML** | Code input | Gentle guidance | Character looks puzzled, preview dims slightly |
| **Idle 30 seconds** | System | Encourage action | Character waves, hint button pulses |
| **Complete section** | System milestone | Big celebration | Full-screen confetti + fanfare + badge reveal |

---

## Animation Specifications

### Micro-Animations

**Button Interactions**:
- Hover: `transform: translateY(-2px) scale(1.05); box-shadow: 0 4px 12px rgba(255,46,99,0.3);` Duration: 200ms ease-out
- Click: `transform: scale(0.95);` Duration: 100ms ease-in, then spring back

**Code Editor**:
- Cursor blink: 1s interval, fade opacity 0 to 1
- Tag completion: Scale from 0.9 to 1 with opacity 0 to 1, 150ms ease-out
- Syntax highlight: Color transition 200ms when tag becomes valid

**Preview Area**:
- Element appears: Scale from 0.8 to 1 + fade in, 300ms ease-out with slight bounce
- Element updates: Pulse effect (scale 1 → 1.1 → 1) 200ms
- Character reactions:
  - Typing: Character "watches" cursor (head turns slightly)
  - Success: Jump animation (translateY -20px → 0) with rotation
  - Error: Head shake (rotate -5deg → 5deg → 0) 3 times

**Success Celebrations**:
- **Lesson Complete**:
  - Confetti particles: 30-50 particles fall from top, random colors (magenta/cyan/yellow)
  - Screen flash: Green overlay 100ms fade in/out
  - XP counter: Count up animation 500ms with easing
  - Character: Dance animation (rotate + scale sequence)

- **Achievement Unlocked**:
  - Modal backdrop: Fade in 200ms
  - Badge: Scale from 0 with rotation (360deg) and bounce, 600ms
  - Particle burst: Radial explosion from badge center
  - Sound: Optional celebratory chime

**Loading States**:
- Spinner: Rotate 360deg continuously with gradient color shift
- Skeleton screens: Pulse opacity 0.6 to 1, 1.5s infinite

---

## State Management

### Empty State
**When**: User completes all available lessons
**Show**:
```
┌─────────────────────────────────┐
│     🎉 You did it all! 🎉      │
│                                 │
│  You've mastered HTML basics!  │
│                                 │
│  [Review Lessons]               │
│  [Share Achievement]            │
│  [What's Next? CSS →]           │
└─────────────────────────────────┘
```
**Animation**: Trophy scales up with glow effect

---

### Loading State
**When**: Initial app load or lesson transition
**Show**:
- Skeleton screens with pulsing placeholders
- Character animation loop (mascot doing something cute)
- Progress bar if loading takes >500ms

**Never show**: Generic spinners or "Loading..." text

---

### Error State
**When**: Invalid HTML typed OR system error
**Show**:
- For code errors: Character confused expression + gentle hint
- For system errors: Friendly message: "Oops! Something went wrong. [Refresh]"

**Never show**: Stack traces, technical errors, red X marks

---

### Success State
**When**: Challenge completed correctly
**Show**:
- Confetti animation
- Character celebration
- XP gain (+5 XP animates up)
- "Next Lesson" button pulses with glow
- Optional: Achievement badge if milestone reached

**Duration**: 2-3 seconds of celebration before allowing next action

---

## Accessibility Considerations

### Keyboard Navigation
- **Tab order**: Challenge → Code Editor → Hint → Tag Palette → Navigation
- **Shortcuts**:
  - `Ctrl/Cmd + Enter`: Submit/check code
  - `Ctrl/Cmd + /`: Toggle hint
  - `Esc`: Close modals
  - Arrow keys: Navigate lesson map

### Screen Reader Support
- All interactive elements have `aria-label`
- Code editor has `role="textbox"` with `aria-multiline="true"`
- Success states announce "Challenge complete! You earned 5 XP"
- Hints have `aria-expanded` state
- Progress bar has `aria-valuenow` and `aria-valuemax`

### Visual Accessibility
- **Color contrast**: All text meets WCAG AA (4.5:1 minimum)
- **Focus indicators**: Visible focus ring (2px magenta outline) on all interactive elements
- **Touch targets**: Minimum 44x44px on mobile
- **Text scaling**: Layout remains functional at 200% zoom

### Motion Sensitivity
- **Respect `prefers-reduced-motion`**:
  - Reduce confetti to simple fade
  - Replace bounces with fades
  - Keep functional animations only
- **Settings toggle**: Manual animation disable option
- **Never**: Flashing effects, rapid color changes, parallax scrolling

---

## Responsive Behavior

### Desktop (>1024px)
- Split-screen layout (50/50 code and preview)
- Tag palette visible
- Full keyboard shortcuts enabled
- Hover states prominent

### Tablet (768px - 1024px)
- Split-screen with 40/60 ratio (more preview space)
- Tag palette collapsible
- Touch-optimized buttons
- Simplified animations

### Mobile (<768px)
- Stacked layout: Challenge → Code → Preview (vertically)
- Tabs to switch between Code and Preview views
- Tag palette as bottom sheet
- Larger touch targets (48x48px minimum)
- Essential animations only (respects data usage)

---

## Copy & Tone Guidelines

**DO use**:
- "Make the robot wave" (action-oriented)
- "Nice work!" (encouraging)
- "Try adding..." (gentle suggestion)
- "You're on fire! 🔥" (celebratory)

**DON'T use**:
- "Implement the following code" (too formal)
- "Incorrect" (discouraging)
- "Error: Invalid syntax" (technical jargon)
- "You must complete..." (demanding)

**Character Voice**:
- Friendly peer, not teacher
- Encouraging, never condescending
- Playful without being childish
- Clear and concise

---

## Edge Cases & Special States

### First Visit
- Show 5-second animated intro explaining concept
- Highlight where to start typing
- Show one example animation to set expectations

### Return After Long Break
- "Welcome back! Let's continue where you left off"
- Show progress summary: "You've completed 12 lessons!"
- Offer quick refresher or continue directly

### Rapid Clicking/Typing
- Throttle animations to prevent overwhelming the user
- Queue celebrations if multiple achievements happen simultaneously
- Ensure animations complete before starting new ones

### Offline State
- Show friendly message: "You're offline! Progress saves locally."
- Disable lesson downloads but allow practice on completed lessons
- Sync progress when back online

### Browser Compatibility
- Graceful degradation: Works without animations if CSS/JS unsupported
- Polyfills for older browsers (if supporting IE11, though not recommended)
- Feature detection before using modern APIs

---

## Success Metrics (UX-focused)

**Primary Metrics**:
1. **Time to First Success**: How fast user completes first lesson (target: <60s)
2. **Engagement Duration**: Session length (target: >5 min average)
3. **Completion Rate**: % who finish first 3 lessons (target: >60%)
4. **Return Rate**: % who come back within 7 days (target: >30%)

**Delight Indicators**:
- Social shares of achievements
- "This is fun!" sentiment in feedback
- Low bounce rate on landing page (<40%)
- High interaction rate with animations

**Red Flags**:
- High drop-off on specific lesson (indicates difficulty spike)
- Settings toggle to disable animations (could mean too overwhelming)
- Long idle time (indicates confusion or boredom)

---

## Design Principles Summary

1. **Immediate Gratification**: Every action gets instant visual feedback
2. **Learn by Doing**: Zero reading before first interaction
3. **Fail Forward**: Mistakes are learning opportunities, not punishments
4. **Progressive Complexity**: Start simple, add features gradually
5. **Celebrate Everything**: Positive reinforcement drives motivation
6. **Respect Users**: No dark patterns, no forced signups, no paywalls
7. **Accessibility First**: Works for everyone, not just visual learners
8. **Performance Matters**: Smooth 60fps animations, fast load times

---

## Next Steps for Development

**Phase 1 - Core Experience**:
- Implement main learning interface with live preview
- Build 5 starter lessons (tags, headings, paragraphs, links, images)
- Add basic animations (typing feedback, success confetti)
- Create simple character reactions

**Phase 2 - Gamification**:
- Add progress tracking and XP system
- Create lesson map with visual progression
- Implement achievement badges
- Add celebration moments

**Phase 3 - Polish**:
- Add settings and preferences
- Implement hint system
- Create comprehensive animation library
- Optimize for mobile/tablet
- Add accessibility features

**Phase 4 - Content Expansion**:
- Add remaining lessons (forms, tables, semantic HTML)
- Create advanced challenges
- Add "free build" playground mode

---

## File References for Development

**Key Research Documents**:
- UX Patterns: `.plans/research/ux-research.md`
- User Personas: `.plans/research/product-research.md`
- Visual Style: `.plans/research/ui-research.md`

**Design Tokens** (from UI research):
- Primary Color: Magenta #FF2E63
- Secondary Color: Cyan #00D9FF
- Background Dark: #1A1A2E
- Typography: Space Grotesk (display), Inter (body), JetBrains Mono (code)

---

Status: READY_FOR_REVIEW
