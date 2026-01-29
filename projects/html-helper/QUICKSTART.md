# HTML Helper - Quick Start Guide

## Get Started in 3 Steps

### 1. Install Dependencies
```bash
npm install
```

### 2. Start Development Server
```bash
npm run dev
```

The app will open at `http://localhost:3000`

### 3. Build for Production
```bash
npm run build
```

The optimized build will be in the `dist` folder.

---

## What You'll See

### Landing Page
- **Interactive Demo**: Type HTML immediately and see it render
- **Feature Highlights**: Animations, instant feedback, achievements
- **Start Button**: Click to begin your HTML learning journey

### Learning Interface
- **Left Panel**: Challenge description, code editor, hints
- **Right Panel**: Live preview of your HTML code
- **Top Nav**: Progress bar, XP counter, streak counter

### Lesson Flow
1. Read the challenge
2. Type HTML code in the editor
3. See live preview as you type
4. Click "Check Answer" when ready
5. Celebrate with confetti on success! 🎉
6. Move to next lesson

---

## Key Features

### 🎨 Animations
- Confetti on lesson completion
- Smooth transitions between screens
- Button interactions and hover effects
- Success/error feedback animations

### 🔒 Security
- All user HTML is sanitized with DOMPurify
- Code runs in sandboxed iframe
- 5-second execution timeout
- Rate limiting (50 executions/minute)

### 💾 Progress Tracking
- Progress saves automatically in browser
- No account needed
- Track XP, streak, and achievements
- Resume where you left off

### 🏆 Gamification
- 15 progressive lessons
- XP rewards for completion
- Achievement badges
- Streak counter for daily practice

---

## Development Commands

```bash
# Start dev server with hot reload
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview

# Type check
npm run lint

# Run tests
npm test
```

---

## Project Structure

```
src/
├── components/        # UI components (Landing, Navigation, Lesson, Editor)
├── data/             # Lessons and achievements
├── utils/            # Core utilities (security, storage, animations)
├── styles/           # Global CSS with Tailwind
└── main.ts           # App entry point
```

---

## Tech Stack

- **Vite** - Fast build tool
- **TypeScript** - Type safety
- **Tailwind CSS** - Utility-first styling
- **CodeMirror 6** - Code editor
- **Anime.js** - Animations
- **DOMPurify** - XSS protection
- **LocalForage** - Enhanced storage

---

## Browser Support

- Chrome/Edge: Latest 2 versions
- Firefox: Latest 2 versions
- Safari: Latest 2 versions
- Mobile Safari: iOS 14+

---

## Deployment

Build the project and deploy the `dist` folder to:

- **Netlify** (recommended) - Drag & drop deploy
- **Vercel** - Connect GitHub repo
- **GitHub Pages** - Static hosting
- **Cloudflare Pages** - Fast global CDN

---

## Need Help?

- Check `README.md` for detailed documentation
- Review `uiflow.md` for UI/UX flow
- Open an issue on GitHub

---

**Ready to learn HTML?** Run `npm run dev` and let's get started! 🚀
