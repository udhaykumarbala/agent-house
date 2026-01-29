# HTML Helper - Interactive HTML Learning Platform

Learn HTML by playing with it - every tag you type makes something magical happen! ✨

## Features

- 🎨 **Beautiful Animations** - Delightful animations on every interaction
- ⚡ **Instant Feedback** - Real-time code preview as you type
- 🏆 **Achievements** - Unlock badges as you progress
- 🔒 **Secure** - Multi-layer XSS protection with DOMPurify and sandboxed execution
- 📱 **Responsive** - Works on desktop, tablet, and mobile
- 💾 **Progress Tracking** - Your progress saves automatically (no account needed)
- ✨ **15 Lessons** - From basic tags to complete webpages

## Getting Started

### Prerequisites

- Node.js 18+ and npm

### Installation

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

## Tech Stack

- **Build Tool**: Vite 5.x
- **Language**: TypeScript 5.x
- **Styling**: Tailwind CSS 3.x
- **Animations**: Anime.js 3.x
- **Code Editor**: CodeMirror 6
- **Security**: DOMPurify 3.x
- **Storage**: LocalForage

## Project Structure

```
src/
├── components/         # UI components
│   ├── LandingPage.ts
│   ├── Navigation.ts
│   ├── CodeEditor.ts
│   └── LessonView.ts
├── data/              # Lesson content and achievements
│   ├── lessons.ts
│   └── achievements.ts
├── utils/             # Core utilities
│   ├── sanitizer.ts   # HTML sanitization
│   ├── sandbox.ts     # Sandboxed execution
│   ├── storage.ts     # Progress persistence
│   └── animator.ts    # Animation utilities
├── styles/
│   └── main.css       # Global styles
└── main.ts            # App entry point
```

## Security

This platform implements multiple layers of security to protect users:

1. **Input Validation** - Length and pattern checks
2. **HTML Sanitization** - DOMPurify removes malicious code
3. **Sandboxed Execution** - Iframe with strict sandbox attributes
4. **Content Security Policy** - Strict CSP headers
5. **Rate Limiting** - Prevents abuse (50 executions per minute)
6. **Execution Timeout** - 5-second limit prevents infinite loops

## Design System

### Colors

- **Primary**: #FF2E63 (Magenta)
- **Secondary**: #00D9FF (Cyan)
- **Success**: #00FF9D (Green)
- **Background**: #1A1A2E (Dark)
- **Surface**: #16213E (Elevated)

### Typography

- **Display**: Space Grotesk (headings)
- **Body**: Inter (text)
- **Code**: JetBrains Mono (code editor)

## Lessons

1. Your First Tag: Paragraphs
2. Headings: Making Titles
3. Bold and Italic: Emphasizing Text
4. Lists: Organizing Information
5. Links: Connecting Pages
6. Images: Adding Pictures
7. Line Breaks: Spacing Things Out
8. Divs: Grouping Content
9. Spans: Inline Styling
10. Basic Structure: HTML Document
11. Meta Tags: Page Information
12. Forms Part 1: Input and Buttons
13. Forms Part 2: Textarea and Select
14. Tables: Organizing Data
15. Final Challenge: Complete Webpage

## Achievements

- 🏆 Tag Master - Created your first HTML tag
- ⭐ Tag Closer - Closed 25 tags correctly
- 💪 Link Legend - Created your first working link
- 🚀 Speed Typer - Completed a lesson in under 60 seconds
- 🎨 Halfway Hero - Completed 7 lessons
- 🔥 Streak Keeper - Practiced 3 days in a row
- 💎 HTML Wizard - Completed all 15 lessons!

## Performance

- **Bundle Size**: < 150KB gzipped
- **First Contentful Paint**: < 1.5s on 3G
- **Lighthouse Score**: > 95
- **Animation FPS**: 60 FPS (with reduced motion fallback)

## Accessibility

- WCAG 2.1 Level AA compliant
- Keyboard navigation support
- Screen reader compatible
- `prefers-reduced-motion` respected
- Minimum 4.5:1 contrast ratios
- 44x44px touch targets on mobile

## Browser Support

- Chrome/Edge: Last 2 versions
- Firefox: Last 2 versions
- Safari: Last 2 versions
- Mobile Safari: iOS 14+

## License

MIT License - Feel free to use for learning and teaching!

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

If you find this helpful, please star the repository! ⭐

---

Made with 💖 for learners everywhere
