# Countdown Timer - UI Flow

## States

```
[IDLE] → Start click → [RUNNING] → Timer hits 0 → [FINISHED]
  ↑                       ↓    ↑                      ↓
  └── Reset click ←── [PAUSED] ←── Pause click    Reset click
                                                       ↓
                                                    [IDLE]
```

## Screen: Main Timer View

### Layout
- Full viewport, centered vertically and horizontally
- Dark background (#0D0D0F)
- SVG progress ring (320x320) behind timer digits
- Timer display: MM:SS in JetBrains Mono, light weight
- Status label below timer (Ready / Counting down / Paused / Time's up)
- Action buttons below status

### User Interactions

| Action | Current State | Result |
|--------|--------------|--------|
| Click "Start" | Idle | Timer begins counting down from 60s |
| Click "Pause" | Running | Timer pauses, button becomes "Resume" |
| Click "Resume" | Paused | Timer resumes from where it stopped |
| Click "Reset" | Running/Paused/Finished | Timer resets to 60s, state → Idle |
| Press Space | Any | Same as clicking Start/Pause/Resume |
| Press R | Any | Same as clicking Reset |

### Visual Feedback
- **Progress ring**: Amber gold (#E8C547) arc depletes as time passes
- **Glow effect**: Soft amber glow on timer digits
- **Danger state**: When < 10s remaining, digits and ring turn red (#E85454)
- **Finish animation**: Timer digits pulse slowly for 4 seconds
- **Button micro-interactions**: Scale down on active press (0.97x)
- **Status label**: Fade-in animation on state change

### Responsive Behavior
- Timer font scales with `clamp(5rem, 20vw, 12rem)`
- Separator scales with `clamp(4rem, 16vw, 10rem)`
- SVG ring is fixed 320x320 for consistency
- Buttons use consistent padding with rounded-full shape
