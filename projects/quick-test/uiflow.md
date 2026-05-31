# Rock Paper Scissors — UI Flow

## Main Game Flow

```
┌──────────────────────────────────────────┐
│           INITIAL STATE                  │
│                                          │
│    Rock  Paper  Scissors                 │
│    ┌────────────────────────────────┐    │
│    │  ⚔️ Choose your weapon         │    │
│    └────────────────────────────────┘    │
│                                          │
│    Score:  You: 0  |  Draws: 0  |  CPU: 0│
└──────────────────┬───────────────────────┘
                   │
                   │ User clicks choice (or presses 1/2/3 or r/p/s)
                   ▼
┌──────────────────────────────────────────┐
│           SUSPENSE STATE (500ms)         │
│                                          │
│    Buttons disabled + dimmed             │
│    ┌────────────────────────────────┐    │
│    │  🎲 Choosing... (shake anim)   │    │
│    └────────────────────────────────┘    │
└──────────────────┬───────────────────────┘
                   │
                   │ CPU choice resolved
                   ▼
┌──────────────────────────────────────────┐
│           RESULT STATE                   │
│                                          │
│    ┌────────────────────────────────┐    │
│    │  🪨 vs ✂️                       │    │
│    │  🎉 You win!                   │    │
│    │  🪨 rock beats scissors ✂️      │    │
│    └────────────────────────────────┘    │
│                                          │
│    Score updated with pop animation      │
│    History pip added (green/red/purple)  │
│    Winner button gets glow border        │
└──────────────────┬───────────────────────┘
                   │
                   │ After 800ms
                   ▼
┌──────────────────────────────────────────┐
│           READY STATE                    │
│                                          │
│    Buttons re-enabled                    │
│    Glow borders removed                  │
│    Result stays visible until next play  │
└──────────────────────────────────────────┘
```

## Streak Flow

```
Win 3+ in a row ──▶  🔥 Win streak: 3  (badge, top-right)
                      Re-animates on each subsequent win

Lose 3+ in a row ──▶  💀 Lose streak: 3  (badge, top-right)

Streak broken ──────▶  Badge hides
```

## Reset Flow

```
Click "Reset" ──▶  Scores → 0
                   History pips → cleared
                   Result area → initial "Choose your weapon"
                   Streak badge → hidden
```

## Result Area Styling

| Outcome | Border Color | Background  | Text Color |
|---------|-------------|-------------|------------|
| Win     | #4ADE80/30  | #4ADE80/5   | #4ADE80    |
| Lose    | #F87171/30  | #F87171/5   | #F87171    |
| Draw    | #A78BFA/30  | #A78BFA/5   | #A78BFA    |

## Color Palette

| Token          | Hex       | Usage                    |
|---------------|-----------|--------------------------|
| primary       | #E8A838   | Hover borders, accents   |
| background    | #0F0F13   | Page background          |
| surface       | #1A1A24   | Cards, buttons           |
| surface-light | #24243A   | Hover states             |
| surface-border| #2E2E48   | Borders                  |
| accent-win    | #4ADE80   | Win indicators           |
| accent-lose   | #F87171   | Loss indicators          |
| accent-draw   | #A78BFA   | Draw indicators          |
| text-primary  | #F0F0F5   | Main text                |
| text-secondary| #8888A0   | Labels                   |
| text-muted    | #55556A   | Subtle text              |
