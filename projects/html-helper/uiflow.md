# HTML Helper - UI Flow

## User Journey

### First-Time User Flow

```
Landing Page
    ↓
Try Interactive Demo (type HTML, see live preview)
    ↓
Click "Start Learning" Button
    ↓
Lesson 1: Your First Tag
    ↓
Type code → See live preview → Click "Check Answer"
    ↓
Success! → Confetti animation → XP gained → Achievement unlocked
    ↓
Click "Next" → Lesson 2
    ↓
... Continue through 15 lessons ...
    ↓
Lesson 15: Complete → All Done! → HTML Wizard badge
```

### Returning User Flow

```
App Loads
    ↓
Check localStorage for progress
    ↓
Go directly to last lesson OR lesson map
    ↓
Continue learning from where left off
    ↓
Streak counter updates (if consecutive day)
```

## Screen Flows

### 1. Landing Page
- Hero section with animated emoji
- Interactive demo (try HTML immediately)
- Feature highlights (animations, feedback, achievements)
- "Start Learning" CTA button

### 2. Learning Interface (Main Screen)
**Left Panel: Challenge + Code**
- Lesson title and description
- Theory explanation
- Challenge instructions
- Hint button (progressive hints)
- Code editor (CodeMirror)
- Navigation buttons (Previous, Check, Next)

**Right Panel: Live Preview**
- Real-time HTML rendering
- Sandboxed iframe for security
- Updates on every keystroke (300ms debounce)

**Top Navigation**
- Logo
- Progress indicator (X/15 lessons)
- XP counter
- Streak counter
- Menu button (opens lesson map)

### 3. Lesson Map (Modal)
- Grid of all 15 lessons
- Visual status indicators:
  - ✓ Completed (green)
  - 🔥 Current (pulsing)
  - 🔒 Locked (grayscale)
  - ○ Available
- Click lesson to navigate
- "Continue Learning" button

### 4. Success Modal
- Appears on lesson completion
- Confetti animation
- "Lesson Complete!" message
- XP gain counter animation (+25 XP)
- "Continue" button

### 5. Achievement Modal
- Appears when achievement unlocked
- Badge icon with bounce animation
- Achievement title and description
- "Awesome!" button to dismiss

## Component Hierarchy

```
App
├── Landing Page
│   ├── Hero Section
│   ├── Interactive Demo
│   │   ├── Demo Input (textarea)
│   │   └── Demo Preview
│   └── Feature Cards
│
└── Learning View
    ├── Navigation
    │   ├── Logo
    │   ├── Progress Bar
    │   ├── Stats (XP, Streak)
    │   └── Menu Button
    │
    ├── Lesson View
    │   ├── Left Panel
    │   │   ├── Challenge Section
    │   │   │   ├── Title
    │   │   │   ├── Description
    │   │   │   ├── Theory
    │   │   │   ├── Challenge
    │   │   │   └── Hint Button
    │   │   ├── Code Editor
    │   │   └── Action Buttons
    │   │
    │   └── Right Panel
    │       └── Live Preview (Sandboxed Iframe)
    │
    └── Modals
        ├── Lesson Map
        ├── Success Modal
        └── Achievement Modal
```

## State Management

### User Progress State
```typescript
{
  completedLessons: number[],
  currentLesson: number,
  xp: number,
  streak: number,
  lastVisit: string,
  achievements: string[],
  preferences: {
    theme: 'dark' | 'light',
    soundEnabled: boolean,
    animationsEnabled: boolean
  }
}
```

### Lesson State
```typescript
{
  currentLessonId: number,
  code: string,
  hintsShown: number,
  isComplete: boolean
}
```

## Animations

### Micro-interactions
- Button click: scale down + bounce back
- Button hover: lift + glow
- Input focus: border color + ring
- Hint reveal: slide down + fade in

### Success Celebrations
- Confetti particles fall (canvas-based)
- Screen flash (green overlay)
- XP counter count-up animation
- Badge bounce in with rotation
- Check mark scale + bounce

### Error Feedback
- Shake animation (gentle)
- Red border pulse
- Error message slide in

### Page Transitions
- Fade out → Slide in new content
- Stagger child elements (50ms delay each)

## Responsive Breakpoints

- **Mobile** (< 640px): Stacked layout, simplified nav
- **Tablet** (640px - 1024px): Two-column where applicable
- **Desktop** (> 1024px): Full split-screen experience

## Accessibility Features

- Keyboard navigation (Tab, Enter, Escape, Arrows)
- Focus indicators (visible rings on all interactive elements)
- ARIA labels on all components
- Screen reader announcements for dynamic content
- `prefers-reduced-motion` support (instant state changes)
- Minimum 4.5:1 contrast ratios
- 44x44px touch targets on mobile

## Security Flow

```
User types HTML
    ↓
Debounce 300ms
    ↓
Input validation (length, patterns)
    ↓
DOMPurify sanitization
    ↓
Create sandboxed iframe
    ↓
Set execution timeout (5 seconds)
    ↓
Write sanitized HTML to iframe
    ↓
Monitor for completion/timeout
    ↓
Clear timeout on success
    ↓
Display rendered output
```

## Data Flow

```
User Action (type code)
    ↓
CodeEditor component
    ↓
onChange callback (debounced)
    ↓
SandboxManager.execute(code)
    ↓
Sanitizer.sanitize(code)
    ↓
Sandbox iframe renders
    ↓
User clicks "Check Answer"
    ↓
Validation logic
    ↓
If valid:
    - Animator.successCelebration()
    - Storage.completeLesson()
    - Check for achievements
    - Show XP gain modal
    - Enable next button
If invalid:
    - Animator.errorShake()
    - Show error message
    - Allow retry
```

## Navigation Flow

```
Landing → Learning (Lesson 1)
    ↓
Lesson 1 → Lesson 2 → ... → Lesson 15
    ↑                           ↓
    ← Previous    Next →
    ↓
Menu → Lesson Map
    ↓
Click any unlocked lesson → Navigate
```

---

This UI flow document describes the complete user journey through HTML Helper, from first visit to completing all lessons, with detailed component hierarchy and state management.
