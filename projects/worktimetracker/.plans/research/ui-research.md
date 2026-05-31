# UI Research: WorkTime Tracker

## Design Inspiration

### Toggl Track
- **Visual Style**: Clean light mode with a bold pink/magenta (#E57CD8) primary accent on white. Minimal chrome, generous whitespace, soft rounded UI elements. Timer is front-and-center with large typography.
- **What Makes It Great**: The timer interface is dead-simple — one big start/stop button. Reports use colorful bar charts and pie charts with vibrant, distinct project colors. The sidebar is slim and unobtrusive.
- **What We Can Borrow**: The "timer-first" layout philosophy. The way reports use color blocks to make time distribution instantly scannable. Their excellent use of whitespace to prevent overwhelm in data-heavy screens.

### Clockify
- **Visual Style**: Blue primary (#03A9F4) with a clean, structured dashboard. Uses a traditional sidebar + content area layout. Heavy use of tables for time entries with inline editing.
- **What Makes It Great**: Comprehensive yet organized — lots of data without feeling cluttered. Calendar view for time entries is intuitive. The team dashboard gives a bird's-eye view of everyone's status.
- **What We Can Borrow**: The calendar/timeline visualization for daily entries. The team overview grid showing who's working on what. Inline editing of time entries (click to edit duration).

### Linear (Project Management)
- **Visual Style**: Dark-mode-first with ultra-clean lines, subtle borders, and a sophisticated muted color palette. Uses purple (#5E6AD2) as primary. Keyboard-first interaction model. Animations are buttery smooth but restrained.
- **What Makes It Great**: Feels like a premium developer tool. The density of information is high but never overwhelming thanks to careful typography hierarchy and spacing. Status indicators use small colored dots rather than heavy badges.
- **What We Can Borrow**: The information density approach — fitting lots of data without clutter. Their subtle status indicators. The smooth transitions between views. The command palette pattern for quick actions.

### Harvest
- **Visual Style**: Warm, approachable design with orange (#FA5D00) as primary on light backgrounds. Uses illustration and friendly copywriting. Rounded corners throughout, soft shadows.
- **What Makes It Great**: Feels less "corporate" than competitors. The weekly timesheet grid is a classic for a reason — it's fast to fill out. Invoice generation flows directly from tracked time.
- **What We Can Borrow**: The weekly grid timesheet view for bulk entry. The warm, human-friendly tone. The direct connection between time tracking and reporting/export.

### Rize (Time Analytics)
- **Visual Style**: Modern dark mode with teal/cyan accents on deep charcoal. Uses gradients sparingly but effectively for data viz. Focus timer has a meditative, calm aesthetic.
- **What Makes It Great**: The analytics dashboards are gorgeous — they use area charts and heatmaps to show productivity patterns over time. The daily breakdown is a stacked timeline showing each activity as a colored block.
- **What We Can Borrow**: The timeline/heatmap approach for visualizing work patterns. The calm, focused aesthetic that makes you want to stay productive. Their use of gradients in charts.

## Competitor Color Analysis

| Competitor | Primary Color | Style | Notes |
|------------|--------------|-------|-------|
| Toggl Track | Pink #E57CD8 | Playful, energetic | Unique choice but too playful for enterprise |
| Clockify | Blue #03A9F4 | Standard tech blue | AVOID - generic SaaS blue |
| Harvest | Orange #FA5D00 | Warm, approachable | Warm direction interesting but orange is overused in productivity |
| Hubstaff | Green #36B37E | Activity/nature | AVOID - overused in tracking apps |
| Time Doctor | Blue #2F80ED | Corporate tech | AVOID - indistinguishable from generic dashboards |
| Timely | Teal #00C2B2 | Modern, fresh | Interesting but specific shade is taken |
| Monday.com | Red/Purple gradient | Bold, energetic | Too busy for a focused tool |
| Rize | Teal #0EA5AB | Calm, analytical | Good direction but exact shade is taken |

## Our Differentiation Strategy

To stand out, we should:
- **Avoid**: Blues (#03A9F4, #2F80ED), generic greens (#36B37E), standard teals (#00C2B2) — these are the "time tracker default" colors
- **Avoid**: Bright oranges, hot pinks — too playful for a professional tool
- **Consider**: A warm, grounded palette that feels professional yet distinctive
- **Direction**: **Deep indigo paired with warm amber/copper accents** — this combines the seriousness of dark tones with the energy of warm metallic accents. No major competitor uses this combination. It evokes "focused craftsmanship" — like a well-made watch or a quality leather notebook.

## Color Psychology

**App Mood**: "Focused and reliable, with warmth" — Users need to feel that this tool respects their time (professional, trustworthy) while also feeling approachable enough for daily use (not cold or clinical).

Recommended color directions:

1. **Deep Indigo + Warm Amber** — A rich, dark indigo (#2E1065 to #4338CA range) as the brand anchor paired with warm amber/copper (#D97706 to #B45309) as the action color. Indigo conveys depth, intelligence, and focus. Amber conveys energy, warmth, and productivity. Together they feel like "a focused evening of deep work." This is our top recommendation.

2. **Slate + Terracotta** — A sophisticated cool gray foundation (#334155) with terracotta/burnt sienna accents (#C2410C). Feels architectural and grounded. Terracotta is rising in design trends (2025-2026) and is rare in SaaS. Risk: may feel too muted for some users.

3. **Charcoal + Electric Teal-Cyan** — Deep charcoal backgrounds (#1C1917) with a custom teal-cyan (#06B6D4 shifted toward #0891B2) for interactive elements. Modern and sharp. Risk: closer to competitors like Rize/Timely, though the dark foundation differentiates.

**Recommendation**: Direction 1 (Deep Indigo + Warm Amber). It's the most unique in the time-tracking space, carries the right emotional weight, and provides excellent contrast opportunities.

## Typography Research

**Display Font Options**:
- **Inter**: The modern web standard. Clean, highly legible at all sizes, excellent for data-heavy interfaces. Free, variable font with great number rendering (important for time displays). Used by Linear, Vercel, and other premium tools. **Top pick for body/UI text.**
- **Geist**: Vercel's typeface. Ultra-modern, slightly geometric, excellent for dashboards. Pairs well with data. Slightly more distinctive than Inter. Good for headings.
- **Manrope**: Semi-rounded geometric sans-serif. Warmer than Inter, with a friendly-professional balance. Great for headings to add character without sacrificing readability.

**Body Font Options**:
- **Inter**: Best-in-class legibility for interfaces. Tabular numbers feature is critical for time tracking (digits align in columns). Variable weight support allows fine-tuning hierarchy.
- **IBM Plex Sans**: Slightly more character than Inter, excellent for data. Good tabular figures. Feels "engineered" which suits a productivity tool.

**Recommended Pairing**:
- **Headings**: Manrope (adds warmth and personality)
- **Body/UI**: Inter (unmatched legibility, tabular numbers for time displays)
- **Monospace numbers**: Inter with `font-variant-numeric: tabular-nums` for timer displays

## Visual Style Direction

**Recommended style**: "Warm Precision" — Clean and structured layouts with warm undertones, subtle depth through layered surfaces, and purposeful micro-interactions.

- **Corners**: Rounded (8px default, 12-16px for cards/modals) — modern and approachable without being bubbly
- **Shadows**: Soft, warm-tinted shadows (using indigo/amber tints rather than pure black) — adds depth without feeling heavy
- **Borders**: Subtle, using semi-transparent borders (1px with ~10% opacity) for definition without harshness
- **Animations**: Subtle and functional — 200ms ease transitions for state changes, smooth number transitions for timer counting, gentle slide-ins for panels. No bouncy or playful animations — this is a professional tool.
- **Data Visualization**: Use the amber/warm palette for productive time, cooler tones for breaks/idle. Avoid red/green for time categories (accessibility). Use gradients sparingly in charts for polish.
- **Dark Mode**: Primary mode uses the deep indigo background tones. Light mode available as secondary option with warm off-white (#FFFBF5) background.
- **Iconography**: Outlined icons (not filled), 1.5px stroke weight, matching the clean-but-warm aesthetic. Consider Lucide or Phosphor icon sets.

## Key UI Patterns for Time Tracking

1. **Timer Widget**: Always visible, minimal footprint when running. Large readable time, single tap start/stop. Project/task selector inline.
2. **Weekly Timesheet Grid**: Rows = projects, columns = days. Quick-fill cells. Running totals per row and column.
3. **Daily Timeline**: Vertical timeline showing time blocks color-coded by project. Gaps are visible (encourages complete tracking).
4. **Reports Dashboard**: Summary cards at top (total hours, avg/day, top project), charts below. Export button prominent.
5. **Employee Overview**: Grid/list of team members with status indicators (active/idle/offline), today's hours, current task.

---
Status: READY_FOR_PLANNING
