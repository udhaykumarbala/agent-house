# Product Specification: SubTrack - Subscription Manager

Based on research in:
- `.plans/research/product-research.md`
- `.plans/research/ux-research.md`
- `.plans/research/ui-research.md`

---

## Target User

**Primary Persona: "Mindful Maya"**
- Age 25-40, digitally native
- Manages 10-20 subscriptions, paying $200+/month
- Values simplicity and privacy
- Uses phone as primary device
- Frustrated by surprise charges and complex financial apps

**Secondary Persona: "Freelancer Frank"**
- Self-employed, needs to track business vs personal subscriptions
- Requires export/categorization for tax purposes

---

## Unique Value Proposition

**"See all your subscriptions in one beautiful place. No bank connection needed."**

Core differentiators:
1. Privacy-first (no bank linking required)
2. Truly minimal (one job: track subscriptions)
3. Mobile-first web app (works everywhere)
4. Free forever for core features
5. Smart reminders before renewals

---

## Features (Prioritized)

### P0 - Must Have (MVP)

- [ ] **Dashboard with total spend** - Display monthly/yearly spend prominently at top of screen
  - *Why essential*: Primary insight users want; creates awareness and urgency

- [ ] **Add subscription manually** - Form to enter service name, price, billing cycle, next billing date
  - *Why essential*: Core functionality; must be completable in <30 seconds

- [ ] **Subscription list view** - Card-based list showing logo, name, price, next billing date
  - *Why essential*: Users need to scan and manage their subscriptions at a glance

- [ ] **Popular service templates** - Pre-filled templates for common subscriptions (Netflix, Spotify, etc.)
  - *Why essential*: Reduces friction; addresses risk of users not entering data manually

- [ ] **Delete subscription** - Swipe-to-delete or tap-to-delete option
  - *Why essential*: Users need to remove cancelled subscriptions

- [ ] **Edit subscription** - Modify any subscription details
  - *Why essential*: Prices change, billing dates shift

- [ ] **Local storage** - Data persists on device without account
  - *Why essential*: No-account-required experience; privacy-first approach

- [ ] **Dark mode** - Full dark theme support
  - *Why essential*: User expectation; finance apps often checked at night

- [ ] **Mobile-responsive design** - Works on all screen sizes, optimized for mobile
  - *Why essential*: Mobile-first requirement

### P1 - Should Have

- [ ] **Renewal reminders** - Push/email notifications X days before billing date
  - *Why it adds value*: Prevents surprise charges; key differentiator from spreadsheets

- [ ] **Category tagging** - Assign categories (Entertainment, Productivity, Business, etc.)
  - *Why it adds value*: Helps users understand spending breakdown; supports Freelancer Frank

- [ ] **Monthly/yearly toggle** - Switch between monthly and annualized cost views
  - *Why it adds value*: Many subscriptions bill yearly; users want accurate total picture

- [ ] **Upcoming renewals view** - Calendar or timeline of next 30 days of charges
  - *Why it adds value*: Proactive planning; matches mental model of "bills coming up"

- [ ] **Free trial tracking** - Special flag for trials with end date reminder
  - *Why it adds value*: High-demand feature; prevents unwanted conversions

- [ ] **Search and filter** - Find subscriptions by name, category, or price range
  - *Why it adds value*: Essential for users with many subscriptions

### P2 - Nice to Have (Post-MVP)

- [ ] **User accounts with cloud sync** - Optional sign-up to sync across devices
  - *Future consideration*: Enables multi-device use without compromising privacy-first approach

- [ ] **Spending trends/history** - Chart showing spend over time
  - *Future consideration*: Helps users see if they're reducing subscription spend

- [ ] **Export to CSV** - Download subscription data
  - *Future consideration*: Useful for Freelancer Frank persona (tax purposes)

- [ ] **Payment method tracking** - Note which card each subscription uses
  - *Future consideration*: Helps when cards expire or need consolidation

- [ ] **Shared household tracking** - Collaborate with family members
  - *Future consideration*: Premium feature for family plans

- [ ] **Currency support** - Multiple currencies with conversion
  - *Future consideration*: International user support

- [ ] **Widget support** - Home screen widget showing total spend or upcoming renewals
  - *Future consideration*: Increases engagement and daily utility

---

## User Stories

### P0 Stories (MVP)

1. As a **Mindful Maya**, I want to **see my total monthly subscription cost** so that **I know exactly what I'm paying for recurring services**.

2. As a **Mindful Maya**, I want to **add a new subscription in under 30 seconds** so that **I actually maintain the list instead of abandoning it**.

3. As a **Mindful Maya**, I want to **select from popular service templates** so that **I don't have to type out common subscriptions like Netflix or Spotify**.

4. As a **Mindful Maya**, I want to **see all my subscriptions in a simple list** so that **I can quickly scan what I'm paying for**.

5. As a **Mindful Maya**, I want to **delete subscriptions I've cancelled** so that **my list stays accurate**.

6. As a **Mindful Maya**, I want to **use the app without creating an account** so that **I can try it immediately without commitment or privacy concerns**.

7. As a **Mindful Maya**, I want to **use dark mode** so that **I can check my subscriptions at night without eye strain**.

### P1 Stories

8. As a **Mindful Maya**, I want to **get reminded before a subscription renews** so that **I can cancel before being charged for something I don't use**.

9. As a **Freelancer Frank**, I want to **categorize subscriptions as business or personal** so that **I can track deductible expenses for taxes**.

10. As a **Mindful Maya**, I want to **see my annual cost, not just monthly** so that **I understand the true yearly impact of my subscriptions**.

11. As a **Mindful Maya**, I want to **track free trials separately** so that **I remember to cancel before they convert to paid**.

12. As a **Mindful Maya**, I want to **see upcoming renewals for the next month** so that **I can prepare for charges hitting my account**.

---

## Acceptance Criteria

### Feature: Dashboard with Total Spend
- [ ] Total monthly spend displays prominently at top of screen
- [ ] Amount updates immediately when subscriptions are added/edited/deleted
- [ ] Toggle to switch between monthly and yearly view (P1)
- [ ] Currency symbol displays correctly based on user locale
- [ ] Works correctly with mix of monthly/yearly/weekly billing cycles

### Feature: Add Subscription
- [ ] Form includes: Service name, Price, Billing cycle (weekly/monthly/yearly), Next billing date
- [ ] Optional fields: Category, Notes, Logo/color
- [ ] Form validates required fields before saving
- [ ] Date picker is mobile-friendly (native picker on mobile)
- [ ] Price accepts decimal values and formats correctly
- [ ] Can be completed in <30 seconds (measured UX goal)
- [ ] Haptic feedback on successful save (mobile)

### Feature: Subscription List
- [ ] Each card shows: Service name, Price, Billing cycle, Next billing date
- [ ] Cards are tappable to view/edit details
- [ ] List scrolls smoothly with many items (50+ subscriptions)
- [ ] Empty state shows helpful prompt to add first subscription
- [ ] Visual indicator for subscriptions due within 7 days

### Feature: Popular Service Templates
- [ ] Minimum 20 popular services pre-configured (Netflix, Spotify, Adobe, etc.)
- [ ] Templates include: Name, typical price, billing cycle, logo/icon
- [ ] Search/filter within templates
- [ ] User can modify template values before saving
- [ ] Templates update without requiring app update (remote config)

### Feature: Delete Subscription
- [ ] Swipe-to-delete gesture on mobile
- [ ] Confirmation prompt before permanent deletion
- [ ] Undo option available for 5 seconds after deletion
- [ ] Total spend updates immediately after deletion
- [ ] Animation provides clear feedback (card slides away)

### Feature: Local Storage
- [ ] Data persists after closing/reopening browser
- [ ] Works offline after initial load (PWA)
- [ ] Data survives browser restart
- [ ] Clear data option available in settings
- [ ] Export option before clearing (P2)

### Feature: Dark Mode
- [ ] Respects system preference by default
- [ ] Manual toggle available in settings
- [ ] All text meets WCAG AA contrast requirements
- [ ] No jarring white flashes during theme switch
- [ ] Persists user preference

---

## Success Metrics

| Metric | Target | How We Measure |
|--------|--------|----------------|
| **Activation** | 60% of users add 3+ subscriptions in first session | Analytics: subscription count at end of first session |
| **Time to value** | <2 minutes from landing to first subscription added | Analytics: timestamp difference |
| **Add completion rate** | 80% of started "add subscription" flows complete | Analytics: funnel tracking |
| **Retention** | 40% of users return within 30 days | Analytics: unique return visits |
| **Average subscriptions tracked** | 8+ per active user | Analytics: average across users with 1+ subscription |
| **NPS Score** | 40+ | In-app survey after 7 days |

---

## Out of Scope (for MVP)

Explicitly NOT building in v1:

- Bank/card linking or transaction import
- Automatic subscription detection
- Bill negotiation or cancellation services
- Budgeting or financial planning features
- Credit score or financial health tracking
- Social features or subscription sharing
- Native iOS/Android apps (web-only MVP)
- Multi-currency support
- Team/business accounts
- API or integrations
- Gamification or rewards

---

## Technical Constraints (for handoff to Architect)

Based on product requirements, the technical solution should:

1. **Work offline** - PWA with local storage required
2. **No backend required for MVP** - All data stored locally
3. **Fast initial load** - Target <3s on 3G (mobile-first)
4. **Mobile-optimized** - Touch targets, swipe gestures, responsive
5. **Privacy-first** - No tracking beyond basic analytics
6. **Future-ready** - Architecture should allow optional cloud sync later (P2)

---

## Competitive Positioning Summary

| Competitor | Their Approach | Our Differentiation |
|------------|----------------|---------------------|
| Truebill/Rocket Money | Bank linking, cancellation services, complex | No bank required, minimal, focused |
| Bobby | Beautiful, manual, discontinued | Spiritual successor, cross-platform, maintained |
| Subby | Subscription-gated features | Free forever for core features |
| Spreadsheets | DIY, no reminders, tedious | Delightful UX, smart reminders |

---

## Launch Strategy (Product Perspective)

### Phase 1: Soft Launch
- Product Hunt submission
- Reddit communities (r/personalfinance, r/frugal, r/apps)
- Target: 1,000 users, gather feedback

### Phase 2: Iteration
- Address top 5 user requests
- Add P1 features based on demand
- Target: 10,000 users

### Phase 3: Growth
- SEO optimization for "subscription tracker"
- Freemium model introduction
- Target: 50,000 users

---

## Open Questions (for CEO/Team Discussion)

1. **Freemium timing**: Should we launch with free-only, or include premium tier from day 1?
2. **Reminder delivery**: Push notifications require service worker + permission. Alternative: email reminders?
3. **Template maintenance**: Who maintains the popular services database? How often updated?
4. **Analytics tool**: Which privacy-respecting analytics tool? (Plausible, Fathom, Simple Analytics)

---

Status: READY_FOR_REVIEW
