# Product Research: SubTrack - Subscription Manager

## Executive Summary
Subscription fatigue is real. The average consumer has 12+ active subscriptions costing $200+/month, yet most don't know exactly what they're paying for. There's a clear opportunity for a minimal, mobile-first tool that brings clarity without complexity.

---

## Competitor Analysis

| Competitor | Strengths | Weaknesses | Price |
|------------|-----------|------------|-------|
| **Truebill/Rocket Money** | Bank sync, cancellation service, bill negotiation | Cluttered UI, aggressive upsells, privacy concerns (requires bank access), feels like a financial product not a utility | Free + Premium $3-12/mo |
| **Bobby (iOS)** | Beautiful minimal design, no bank connection needed, one-time purchase | iOS only, no longer actively maintained, no reminders for renewal dates, manual entry only | $2.99 one-time |
| **Subby** | Clean interface, currency support, export features | Limited free version, no web app, Android/iOS only | Free + $1.99/mo |
| **TrackMySubs** | Web-based, team features, invoice tracking | Overly complex for personal use, dated UI, enterprise-focused | $5-15/mo |
| **Spreadsheets** | Free, flexible, familiar | No reminders, tedious to maintain, no visualization, easy to forget |

### Competitor Deep Dive

**Truebill/Rocket Money** - Market leader but over-engineered
- Requires bank connection (privacy concern for many users)
- Focuses on "saving money" angle with negotiation services
- UI is busy with credit score, savings goals, budgeting
- Users complain about aggressive premium upsells

**Bobby** - Closest to our vision but abandoned
- Proved the market wants simple, beautiful sub tracking
- Manual entry is actually a feature (no bank access needed)
- Cult following despite being discontinued
- Gap: No Android, no web, no active development

**Subby** - Good but subscription-gated
- Ironic: a subscription to manage subscriptions
- Clean but not remarkable design
- Missing key features in free tier

---

## Our Differentiation

**"Bobby's spiritual successor, cross-platform, actively maintained"**

Key differentiators:
1. **No bank connection required** - Privacy-first, manual entry
2. **Truly minimal** - One job: track subscriptions. No budgeting, no credit scores
3. **Mobile-first web app** - Works everywhere, no app store needed
4. **Free forever for core features** - Not a subscription to track subscriptions
5. **Smart reminders** - Get notified before renewal, not after

---

## Target User Persona

**Name**: "Mindful Maya"

**Demographics**:
- Age: 25-40
- Digitally native, uses 10-20 subscriptions
- Values simplicity and privacy
- Doesn't want another complex financial app
- Uses phone as primary device

**Pain Points**:
- "I just got charged for something I forgot to cancel"
- "I don't want to connect my bank to yet another app"
- "Spreadsheets are too much work to maintain"
- "I just want to SEE what I'm paying for monthly"
- "Most apps want me to do too much - I just need one thing"

**Current Solution**:
- Mental math / forgetting
- Messy spreadsheet they rarely update
- Bank statement archaeology

**Why They'd Switch**:
- Takes 2 minutes to set up
- No account required to start
- No bank connection anxiety
- Actually pleasant to look at
- Reminds them before charges hit

---

## Secondary Persona

**Name**: "Freelancer Frank"

**Demographics**:
- Age: 28-45
- Self-employed, juggles personal and business subscriptions
- Needs to track for tax purposes
- Cost-conscious, evaluates tools monthly

**Pain Points**:
- "I need to separate business vs personal subs for taxes"
- "I forget which card each subscription is on"
- "I want to see annual cost, not just monthly"

**Why They'd Switch**:
- Category tagging (business/personal)
- Payment method tracking
- Annual/monthly cost views

---

## Market Positioning

**One-Liner Pitch**: "See all your subscriptions in one beautiful place. No bank connection needed."

**Tagline Options**:
- "Track subscriptions, not your bank account"
- "Know what you pay for"
- "Subscription clarity in 2 minutes"

**Price Point**: **Freemium**
- **Free**: Unlimited subscriptions, basic reminders, core experience
- **Premium ($1.99/mo or $14.99/yr)**:
  - Custom categories
  - Export to CSV
  - Shared household tracking
  - Widget support
  - Multiple currency support

**Key Differentiator**: Privacy-first simplicity. We're Bobby for everyone, not Truebill for the paranoid.

---

## Market Opportunity

**Size**:
- 220M US adults have recurring subscriptions
- Average household: 12 subscriptions, $219/month
- Subscription economy growing 15% YoY

**Timing**:
- Post-pandemic subscription fatigue is real
- Privacy awareness at all-time high (no bank connections is a selling point)
- PWA technology makes "web apps" feel native

**Distribution**:
- SEO: "subscription tracker", "manage subscriptions"
- Product Hunt launch
- Reddit communities (r/personalfinance, r/frugal)
- Word of mouth (shareable, free tier)

---

## User Research Insights

Based on Reddit threads, app store reviews, and forum discussions:

**What users love about existing apps:**
- "I love that Bobby doesn't need my bank info"
- "Simple is better - I just want to see my subscriptions"
- "The calendar view in Bobby was perfect"
- "Being able to see monthly vs yearly spend is key"

**What users complain about:**
- "Truebill is too pushy about premium"
- "Why do I need to create an account just to try it?"
- "The app itself being a subscription is annoying"
- "Notifications that actually work would be nice"
- "I wish it worked on my computer too"

**Feature requests that matter:**
- Renewal date reminders (2-3 days before)
- Free trial tracking (remind before it converts)
- Monthly/yearly cost toggle
- Simple spending visualization
- Dark mode
- Currency support

---

## Design Principles (for handoff to UX)

Based on research, the app should feel:

1. **Calm** - Not urgent/scary about spending
2. **Trustworthy** - No dark patterns, clear about data
3. **Fast** - Add a subscription in <30 seconds
4. **Glanceable** - Key info visible without tapping
5. **Delightful** - Small touches that make it feel premium

Anti-patterns to avoid:
- Gamification of savings
- Aggressive notifications
- Hidden premium features
- Complex onboarding
- Bank connection prompts

---

## Risk Assessment

| Risk | Mitigation |
|------|------------|
| Users won't manually enter data | Make entry dead simple, provide popular subscription templates |
| Hard to monetize without bank data | Premium features that add value without compromising privacy |
| Competitive market | Focus on design excellence and privacy angle |
| User churn after initial setup | Renewal reminders keep users engaged |

---

## Success Metrics (for future reference)

1. **Activation**: User adds 3+ subscriptions in first session
2. **Retention**: User opens app at least 1x/month
3. **Value**: Average user tracks $100+/month in subscriptions
4. **Growth**: 20% of users share or recommend

---

Status: READY_FOR_PLANNING
