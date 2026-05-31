# UX Research: Unit Converter (Next-Gen Minimal)

## Pattern Analysis

Referenced apps and their successful patterns:

### 1. Apple Calculator (iOS Built-in)
- **Pattern Used**: Single-screen interface with large tap targets, immediate visual feedback on input, no page navigation for core functionality.
- **Why It Works**: Zero learning curve. The user sees the entire tool at once. Large typography makes values scannable at a glance. The dark theme with accent colors creates a premium, focused feel.
- **What We Can Learn**: A converter should feel like an instrument — everything on one screen, no hunting through menus. Large, confident typography for the primary value is essential.

### 2. Elk - Travel Currency Converter
- **Pattern Used**: Swipe-based unit switching, single prominent value display, minimal chrome. Uses a "ruler" metaphor — slide to change the value, tap to switch currencies. Beautiful type hierarchy with the converted value shown large and the source value secondary.
- **Why It Works**: It removes all friction. No dropdowns, no keypads blocking the screen. The gesture-driven interaction feels native and fluid. The app won Best of App Store for its polish.
- **What We Can Learn**: Gesture-based input (drag/swipe) feels more premium than tapping a keypad. Showing one conversion pair at a time with large type keeps focus. The "less UI, more result" philosophy is key for a next-gen feel.

### 3. Numi (macOS Calculator/Converter)
- **Pattern Used**: Text-based natural language input. Users type "5 kg in pounds" and get instant inline results. No buttons, no dropdowns — just a clean text field and results.
- **Why It Works**: Feels like magic for power users. The minimal interface (just text) is the ultimate "less is more." Real-time conversion as you type creates a responsive, alive feel.
- **What We Can Learn**: Real-time conversion (no "convert" button) is expected in modern tools. Inline results feel faster than navigating to a results screen.

### 4. Google Unit Converter (Search Widget)
- **Pattern Used**: Two input fields side by side (or stacked), dropdown selectors for units, swap button between them. Instant conversion on input change.
- **Why It Works**: Universally understood layout. The swap arrow is iconic — users immediately know they can reverse the conversion. Both fields are editable, so users can type in either direction.
- **What We Can Learn**: The swap/reverse button is a must-have convention. Bidirectional input (type in either field) reduces friction. This is the baseline — our app needs to match this convenience then exceed it in polish.

### 5. Converto (iOS Converter)
- **Pattern Used**: Category tabs at top (Length, Weight, Temp), large result display, custom numeric keypad at bottom, unit selectors as horizontal scrollable pills. Smooth animations on unit change.
- **Why It Works**: Clear category separation without deep navigation. The pill-style unit selectors are thumb-friendly and visually scannable. Animations on value change give satisfying feedback.
- **What We Can Learn**: Category switching should be effortless (tabs or segmented control, not a menu). Pill/chip selectors for units feel modern and are easier to scan than dropdowns. Micro-animations on value transitions create perceived quality.

### 6. Amount (iOS Converter by Marco Masser)
- **Pattern Used**: Minimalist single-screen with bold monospaced numbers, subtle color coding per category, haptic feedback on interactions, smooth spring animations.
- **Why It Works**: The restraint in design makes it feel premium. Every interaction has tactile feedback. The monospaced font for numbers prevents layout shifts during typing.
- **What We Can Learn**: Monospaced or tabular-number fonts prevent jank when values change. Haptic feedback (on mobile) elevates perceived quality. Subtle color theming per category (e.g., blue for temperature, green for weight) aids recognition without cluttering.

## User Journey Analysis

### Entry Point
User needs to convert a value — triggered by cooking, travel, fitness, homework, or work. They open the app expecting to get an answer in under 5 seconds.

### Core Loop
1. Select category (length / weight / temperature) — or it defaults to last used
2. Enter a numeric value
3. See the converted result instantly
4. Optionally swap direction or change units
5. Copy result or enter a new value

### Success State
The user sees the converted value clearly displayed, confirms it's correct, and either copies it or leaves the app. Total time: 3-8 seconds. No confirmation screen needed — the result IS the success state.

## Mental Model

Users think of this task as: **"I have a number in one unit, I need it in another."** It's a simple translation — input goes in, answer comes out. They think of it like a calculator, not like a form to fill out.

Common terminology:
- "Convert X to Y"
- "How many [unit] in [value] [unit]?"
- Categories: length/distance, weight/mass, temperature
- Units: km, miles, meters, feet, inches, cm, kg, lbs, oz, grams, Celsius, Fahrenheit, Kelvin

## Anti-Patterns to Avoid

- **Don't**: Use dropdown menus for unit selection — **Why**: Dropdowns require precise taps, obscure options behind a click, and feel dated. Pill selectors or scrollable lists are faster and more visual.
- **Don't**: Require a "Convert" button to trigger conversion — **Why**: Users expect instant/real-time results. A manual trigger adds friction and feels like a 2010-era web form.
- **Don't**: Show ALL unit categories and ALL units on one screen — **Why**: Information overload kills the "minimal" goal. Progressive disclosure (show one category, let user switch) keeps focus.
- **Don't**: Use the system default keyboard for number input — **Why**: The standard keyboard wastes space with letters. A custom numeric pad (or no pad at all with gesture input) feels more intentional and premium.
- **Don't**: Navigate to a separate "results" page — **Why**: Conversion is instant mental math — the UI should mirror that speed. Same-screen results are expected.
- **Don't**: Use skeuomorphic design or heavy gradients — **Why**: Clashes with the "next-gen minimal" brief. Clean, flat design with subtle depth (soft shadows, glassmorphism hints) reads as modern.
- **Don't**: Forget empty/zero state design — **Why**: The app at launch with "0" everywhere looks broken. A welcoming prompt or placeholder text ("Enter a value...") makes the first impression feel alive.

## Recommended Patterns for Our App

Based on research, we should use:

1. **Single-screen layout with category tabs** — Segmented control or horizontal tabs at top for Length / Weight / Temperature. No page navigation for the core task. Inspired by Converto and Apple's design language.

2. **Large, bold result-first typography** — The converted value should be the visual hero of the screen (48-64px, bold weight). Source value is secondary. Monospaced or tabular-lining numerals to prevent layout jank. Inspired by Elk and Amount.

3. **Pill/chip unit selectors** — Horizontal scrollable chips for selecting source and target units. Visually scannable, thumb-friendly, no dropdowns. A prominent swap/reverse button between the two unit rows. Inspired by Converto and Google Converter.

4. **Real-time conversion with micro-animations** — Value updates as user types, with a subtle number-roll or fade transition on the result. No "convert" button. Spring-based animations for swap action. Inspired by Elk and Amount.

5. **Subtle category color theming** — Each category gets a muted accent color (e.g., blue for temperature, amber for length, green for weight) applied to active elements. Keeps the interface minimal while providing recognition cues.

6. **Custom numeric input area** — A clean, minimal numeric keypad at the bottom third of the screen (mobile) or inline editable fields (desktop). Large touch targets (min 48px), decimal point, backspace, and clear. No unnecessary keys.

7. **Thoughtful empty and edge states** — Launch state shows a friendly prompt, not zeros. Error states (e.g., nonsensical temperature below absolute zero) handled with gentle inline messages, not alerts.

8. **Dark mode as default with light mode option** — Dark backgrounds with light text feel premium for utility apps (ref: Apple Calculator, many fintech apps). Support both modes, but lead with dark to match the "premium utility" brief.

## Typography & Visual Direction

- **Primary font**: Inter or SF Pro — clean, geometric, excellent tabular figures
- **Number display**: Tabular lining figures, 48-64px for primary result, 24-32px for input value
- **Hierarchy**: 3 levels max — result value > input value > labels/units
- **Spacing**: Generous whitespace, 16px minimum padding, 8px grid system
- **Colors**: Near-black background (#0A0A0A or #121212), white text, category accent colors at 60% opacity for subtlety
- **Radius**: 12-16px for cards/containers, 24px for pills/chips — rounded but not bubbly
- **Shadows**: Minimal, only for elevation (floating keypad, modals)

## Responsive Considerations

- **Mobile (primary)**: Full-screen single view, keypad at bottom, tabs at top, result in center
- **Tablet**: Same layout, wider pill rows, potentially show multiple conversions
- **Desktop/Web**: Centered card layout (max-width 480px), keyboard input replaces keypad, hover states on interactive elements

---
Status: READY_FOR_PLANNING
