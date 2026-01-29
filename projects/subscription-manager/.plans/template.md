# Template Selection

## Chosen Template: `fullstack`

## Task Analysis
A subscription management web app requires:
- User authentication (each user has their own subscriptions)
- Persistent data storage (subscriptions, payment dates, costs)
- CRUD operations via API
- Mobile-first responsive UI with real-time state
- Calculations (spending summaries, upcoming renewals)

This is a data-driven application, not a static site.

## Justification
- **Requires backend**: User data must persist across sessions and devices - needs database + API
- **Requires authentication**: Each user needs private access to their subscriptions
- **Complex state management**: Subscription lists, filters, sorting, spending calculations need proper frontend framework
- **Mobile-first SPA**: Next.js with responsive Tailwind provides excellent mobile UX
- **Scalability**: Fullstack template allows adding features like notifications, sharing, analytics

## Customizations Needed
- Simplified database schema (single table for subscriptions to start)
- Mobile-first Tailwind configuration
- PWA support for mobile app-like experience (optional enhancement)
- Remove unused fullstack boilerplate to keep it minimal

## Tech Stack Summary
| Layer | Technology | Reason |
|-------|------------|--------|
| Frontend | Next.js 14 + TypeScript | SSR, routing, excellent DX |
| Styling | Tailwind CSS | Mobile-first utilities, minimal custom CSS |
| UI Components | shadcn/ui | Clean, accessible, customizable |
| State | Zustand | Lightweight, simple for this scope |
| Backend | Go + Chi | Fast, simple REST API |
| Database | PostgreSQL | Reliable, supports JSON for flexibility |
| Container | Docker | Consistent dev/prod environments |

## Alternative Considered
`nextjs-frontend` was considered (using localStorage), but rejected because:
- Data would be lost on device/browser change
- No cross-device sync
- Can't support future features (notifications, sharing)

---
Status: AWAITING_REVIEW
