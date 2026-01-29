# Architecture Research: SubTrack - Subscription Manager

## Executive Summary

This document defines the technical architecture for a privacy-first, mobile-first subscription management web application. Based on the product requirements and research from PM, UX, UI, and Security teams, we'll design a fullstack system using Next.js 14 + Go API + PostgreSQL.

**Key Architectural Decisions:**
- Fullstack template with clear frontend/backend separation
- Session-based authentication for MVP simplicity
- Mobile-first PWA with offline support
- PostgreSQL for reliable, structured data storage
- Privacy-first: no bank connections, user-controlled data

---

## 1. Technical Requirements Analysis

### Functional Requirements

| Requirement | Technical Implication |
|-------------|----------------------|
| User authentication | Session management, secure password hashing |
| CRUD subscriptions | REST API endpoints, database schema |
| Total spend calculations | Server-side aggregation or client-side compute |
| Renewal reminders | Background jobs, push notifications |
| Category filtering | Database indexing, query optimization |
| Multi-currency support (future) | Currency field, conversion considerations |
| Data export | CSV/JSON generation endpoint |
| Dark mode | CSS variables, system preference detection |

### Non-Functional Requirements

| Requirement | Target | Implementation |
|-------------|--------|----------------|
| Performance | <100ms API response, <3s initial load | CDN, code splitting, optimized queries |
| Availability | 99.9% uptime | Health checks, graceful degradation |
| Scalability | 10K concurrent users | Stateless API, connection pooling |
| Security | OWASP Top 10 compliant | See security research document |
| Mobile-first | Works on 375px+ screens | Responsive design, touch-friendly |
| Offline support | View subscriptions offline | Service worker, IndexedDB caching |

---

## 2. System Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                           CLIENTS                                    │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│   ┌──────────────┐    ┌──────────────┐    ┌──────────────┐          │
│   │   Mobile     │    │   Desktop    │    │   Tablet     │          │
│   │   Browser    │    │   Browser    │    │   Browser    │          │
│   └──────┬───────┘    └──────┬───────┘    └──────┬───────┘          │
│          │                   │                   │                   │
│          └───────────────────┼───────────────────┘                   │
│                              │                                       │
│                    ┌─────────▼─────────┐                            │
│                    │   Next.js PWA     │                            │
│                    │   (Frontend)      │                            │
│                    │   - React 18      │                            │
│                    │   - TypeScript    │                            │
│                    │   - Tailwind      │                            │
│                    │   - Service Worker│                            │
│                    └─────────┬─────────┘                            │
│                              │                                       │
└──────────────────────────────┼───────────────────────────────────────┘
                               │ HTTPS
┌──────────────────────────────┼───────────────────────────────────────┐
│                              │                                       │
│                    ┌─────────▼─────────┐                            │
│                    │   API Gateway     │                            │
│                    │   (Chi Router)    │                            │
│                    │   - Rate Limiting │                            │
│                    │   - Auth Middleware│                           │
│                    │   - CORS          │                            │
│                    └─────────┬─────────┘                            │
│                              │                                       │
│              ┌───────────────┼───────────────┐                       │
│              │               │               │                       │
│    ┌─────────▼─────┐ ┌──────▼──────┐ ┌─────▼─────────┐             │
│    │    Auth       │ │ Subscription│ │  Notification │             │
│    │   Service     │ │   Service   │ │    Service    │             │
│    └───────┬───────┘ └──────┬──────┘ └───────┬───────┘             │
│            │                │                │                       │
│            └────────────────┼────────────────┘                       │
│                             │                                        │
│                    ┌────────▼────────┐                              │
│                    │   PostgreSQL    │                              │
│                    │   Database      │                              │
│                    └─────────────────┘                              │
│                                                                      │
│                           BACKEND                                    │
└──────────────────────────────────────────────────────────────────────┘
```

### Component Responsibilities

| Component | Responsibility |
|-----------|----------------|
| **Next.js Frontend** | UI rendering, client-side state, PWA shell, offline caching |
| **Go API** | Business logic, authentication, data validation, database operations |
| **PostgreSQL** | Data persistence, transactions, querying |
| **Service Worker** | Offline support, push notifications, background sync |

---

## 3. Tech Stack Decision Matrix

### Frontend

| Technology | Alternative | Decision | Rationale |
|------------|-------------|----------|-----------|
| **Next.js 14** | Vite + React, Remix | Next.js | SSR for SEO, App Router, excellent DX, Vercel ecosystem |
| **TypeScript** | JavaScript | TypeScript | Type safety, better IDE support, fewer runtime errors |
| **Tailwind CSS** | CSS Modules, styled-components | Tailwind | Mobile-first utilities, rapid prototyping, consistent design |
| **shadcn/ui** | Radix, Chakra, MUI | shadcn/ui | Accessible, customizable, no bundle bloat (copy-paste) |
| **Zustand** | Redux, Jotai, Context | Zustand | Minimal boilerplate, small bundle, simple mental model |
| **React Query** | SWR, RTK Query | React Query | Caching, refetching, offline support, devtools |

### Backend

| Technology | Alternative | Decision | Rationale |
|------------|-------------|----------|-----------|
| **Go** | Node.js, Python, Rust | Go | Fast, simple, excellent for APIs, low memory footprint |
| **Chi Router** | Gin, Echo, Fiber | Chi | Idiomatic Go, middleware support, net/http compatible |
| **sqlc** | GORM, sqlx, raw SQL | sqlc | Type-safe queries, no ORM magic, compile-time checking |
| **PostgreSQL** | MySQL, SQLite, MongoDB | PostgreSQL | ACID compliance, JSON support, mature ecosystem |
| **Argon2** | bcrypt, scrypt | Argon2 | Modern, memory-hard, OWASP recommended |

### Infrastructure

| Technology | Alternative | Decision | Rationale |
|------------|-------------|----------|-----------|
| **Docker** | Bare metal, Podman | Docker | Consistent environments, easy deployment |
| **Docker Compose** | K8s (overkill for MVP) | Compose | Simple local development, easy to understand |
| **Vercel** (frontend) | Netlify, Cloudflare | Vercel | Next.js creator, excellent DX, edge functions |
| **Railway/Render** (backend) | Fly.io, AWS | Railway | Simple Go deployment, PostgreSQL included |

---

## 4. Data Model Design

### Entity Relationship Diagram

```
┌─────────────────────────────────────────────────────────────────────┐
│                         DATABASE SCHEMA                              │
└─────────────────────────────────────────────────────────────────────┘

┌──────────────────────┐         ┌──────────────────────┐
│        users         │         │    subscriptions     │
├──────────────────────┤         ├──────────────────────┤
│ id (uuid) PK         │────────<│ id (uuid) PK         │
│ email (unique)       │         │ user_id (uuid) FK    │
│ password_hash        │         │ name                 │
│ created_at           │         │ description          │
│ updated_at           │         │ price                │
│ email_verified       │         │ currency             │
│ preferences (jsonb)  │         │ billing_cycle        │
└──────────────────────┘         │ billing_day          │
                                 │ next_billing_date    │
┌──────────────────────┐         │ category             │
│      sessions        │         │ icon                 │
├──────────────────────┤         │ color                │
│ id (uuid) PK         │         │ notes                │
│ user_id (uuid) FK    │────────<│ is_active            │
│ token_hash           │         │ reminder_days        │
│ expires_at           │         │ created_at           │
│ created_at           │         │ updated_at           │
│ ip_address           │         └──────────────────────┘
│ user_agent           │
└──────────────────────┘         ┌──────────────────────┐
                                 │  subscription_logs   │
                                 ├──────────────────────┤
                                 │ id (uuid) PK         │
                                 │ subscription_id FK   │
                                 │ action               │
                                 │ old_price            │
                                 │ new_price            │
                                 │ changed_at           │
                                 └──────────────────────┘
```

### Table Definitions

#### users
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    email_verified BOOLEAN DEFAULT FALSE,
    preferences JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
```

#### sessions
```sql
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(64) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    ip_address INET,
    user_agent TEXT
);

CREATE INDEX idx_sessions_token ON sessions(token_hash);
CREATE INDEX idx_sessions_user ON sessions(user_id);
CREATE INDEX idx_sessions_expiry ON sessions(expires_at);
```

#### subscriptions
```sql
CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),
    currency VARCHAR(3) DEFAULT 'USD',
    billing_cycle VARCHAR(20) NOT NULL CHECK (billing_cycle IN ('weekly', 'monthly', 'yearly', 'custom')),
    billing_day INTEGER CHECK (billing_day >= 1 AND billing_day <= 31),
    custom_days INTEGER, -- for custom billing cycles
    next_billing_date DATE NOT NULL,
    category VARCHAR(50) DEFAULT 'other',
    icon VARCHAR(50), -- lucide icon name
    color VARCHAR(7), -- hex color
    notes TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    reminder_days INTEGER DEFAULT 3,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_subscriptions_user ON subscriptions(user_id);
CREATE INDEX idx_subscriptions_next_billing ON subscriptions(next_billing_date);
CREATE INDEX idx_subscriptions_active ON subscriptions(user_id, is_active);
```

### Enums & Constants

```typescript
// Frontend types (TypeScript)
type BillingCycle = 'weekly' | 'monthly' | 'yearly' | 'custom';

type Category =
  | 'entertainment'   // Netflix, Spotify, gaming
  | 'productivity'    // Adobe, Notion, Figma
  | 'lifestyle'       // Gym, meal kits
  | 'utilities'       // Cloud storage, VPN
  | 'news'            // NYT, WSJ, newsletters
  | 'education'       // Courses, learning platforms
  | 'finance'         // Banking, investing apps
  | 'health'          // Health apps, wellness
  | 'other';          // Uncategorized

type Currency = 'USD' | 'EUR' | 'GBP' | 'CAD' | 'AUD' | 'INR';
```

---

## 5. API Design

### RESTful Endpoints

```
Base URL: /api/v1

Authentication:
  POST   /auth/register         - Create new account
  POST   /auth/login            - Login, get session
  POST   /auth/logout           - Invalidate session
  POST   /auth/forgot-password  - Request password reset
  POST   /auth/reset-password   - Reset password with token
  GET    /auth/me               - Get current user

Subscriptions:
  GET    /subscriptions         - List user's subscriptions
  POST   /subscriptions         - Create subscription
  GET    /subscriptions/:id     - Get subscription details
  PUT    /subscriptions/:id     - Update subscription
  DELETE /subscriptions/:id     - Delete subscription

Analytics:
  GET    /analytics/summary     - Total spend, by category, etc.
  GET    /analytics/upcoming    - Next 30 days renewals

User:
  PUT    /user/preferences      - Update preferences
  GET    /user/export           - Export all data (GDPR)
  DELETE /user/account          - Delete account
```

### Request/Response Examples

#### Create Subscription
```http
POST /api/v1/subscriptions
Content-Type: application/json
Authorization: Bearer {session_token}

{
  "name": "Netflix",
  "price": 15.99,
  "currency": "USD",
  "billing_cycle": "monthly",
  "billing_day": 15,
  "category": "entertainment",
  "icon": "tv",
  "color": "#E50914",
  "notes": "Premium plan, 4 screens",
  "reminder_days": 3
}
```

Response:
```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Netflix",
    "price": 15.99,
    "currency": "USD",
    "billing_cycle": "monthly",
    "next_billing_date": "2024-02-15",
    "category": "entertainment",
    "icon": "tv",
    "color": "#E50914",
    "is_active": true,
    "created_at": "2024-01-18T10:30:00Z"
  }
}
```

#### Get Analytics Summary
```http
GET /api/v1/analytics/summary
Authorization: Bearer {session_token}
```

Response:
```json
{
  "data": {
    "total_monthly": 156.47,
    "total_yearly": 1877.64,
    "active_count": 12,
    "by_category": {
      "entertainment": 45.97,
      "productivity": 62.50,
      "utilities": 48.00
    },
    "next_renewal": {
      "date": "2024-01-20",
      "subscription": "Spotify",
      "amount": 10.99
    }
  }
}
```

### Error Response Format

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid subscription data",
    "details": [
      { "field": "price", "message": "Price must be positive" }
    ]
  }
}
```

---

## 6. Project Structure

### Frontend (Next.js)

```
frontend/
├── app/                          # Next.js 14 App Router
│   ├── (auth)/                   # Auth route group
│   │   ├── login/
│   │   │   └── page.tsx
│   │   ├── register/
│   │   │   └── page.tsx
│   │   └── layout.tsx
│   ├── (dashboard)/              # Protected routes
│   │   ├── page.tsx              # Dashboard home
│   │   ├── subscriptions/
│   │   │   ├── page.tsx          # List view
│   │   │   ├── [id]/
│   │   │   │   └── page.tsx      # Detail/edit
│   │   │   └── new/
│   │   │       └── page.tsx      # Add new
│   │   ├── analytics/
│   │   │   └── page.tsx
│   │   └── settings/
│   │       └── page.tsx
│   ├── layout.tsx                # Root layout
│   ├── globals.css
│   └── manifest.json             # PWA manifest
├── components/
│   ├── ui/                       # shadcn/ui components
│   │   ├── button.tsx
│   │   ├── card.tsx
│   │   ├── input.tsx
│   │   └── ...
│   ├── subscription/
│   │   ├── subscription-card.tsx
│   │   ├── subscription-form.tsx
│   │   ├── subscription-list.tsx
│   │   └── category-filter.tsx
│   ├── dashboard/
│   │   ├── spend-summary.tsx
│   │   ├── upcoming-renewals.tsx
│   │   └── quick-actions.tsx
│   ├── layout/
│   │   ├── header.tsx
│   │   ├── bottom-nav.tsx
│   │   └── sidebar.tsx
│   └── common/
│       ├── loading.tsx
│       ├── empty-state.tsx
│       └── error-boundary.tsx
├── lib/
│   ├── api.ts                    # API client (fetch wrapper)
│   ├── auth.ts                   # Auth utilities
│   ├── utils.ts                  # General utilities
│   └── constants.ts              # App constants
├── hooks/
│   ├── use-subscriptions.ts      # React Query hooks
│   ├── use-analytics.ts
│   └── use-auth.ts
├── stores/
│   └── app-store.ts              # Zustand store
├── types/
│   ├── subscription.ts
│   ├── user.ts
│   └── api.ts
├── public/
│   ├── icons/                    # PWA icons
│   └── sw.js                     # Service worker
├── tailwind.config.ts
├── next.config.js
├── tsconfig.json
└── package.json
```

### Backend (Go)

```
backend/
├── cmd/
│   └── api/
│       └── main.go               # Entry point
├── internal/
│   ├── config/
│   │   └── config.go             # Environment config
│   ├── handler/
│   │   ├── auth.go               # Auth handlers
│   │   ├── subscription.go       # Subscription handlers
│   │   ├── analytics.go          # Analytics handlers
│   │   ├── user.go               # User handlers
│   │   └── middleware.go         # HTTP middleware
│   ├── service/
│   │   ├── auth.go               # Auth business logic
│   │   ├── subscription.go       # Subscription logic
│   │   └── analytics.go          # Analytics logic
│   ├── repository/
│   │   ├── user.go               # User DB operations
│   │   ├── session.go            # Session DB operations
│   │   └── subscription.go       # Subscription DB operations
│   ├── model/
│   │   ├── user.go               # User model
│   │   ├── session.go            # Session model
│   │   └── subscription.go       # Subscription model
│   ├── database/
│   │   ├── postgres.go           # DB connection
│   │   └── migrations/           # SQL migrations
│   │       ├── 001_init.up.sql
│   │       └── 001_init.down.sql
│   └── validator/
│       └── validator.go          # Input validation
├── pkg/
│   ├── password/
│   │   └── argon2.go             # Password hashing
│   └── response/
│       └── json.go               # JSON response helpers
├── sqlc/
│   ├── queries/
│   │   ├── users.sql
│   │   ├── sessions.sql
│   │   └── subscriptions.sql
│   ├── schema.sql
│   └── sqlc.yaml
├── Dockerfile
├── go.mod
├── go.sum
└── Makefile
```

### Root Project Structure

```
subscription-manager/
├── frontend/                     # Next.js app
├── backend/                      # Go API
├── docker-compose.yml            # Local development
├── docker-compose.prod.yml       # Production config
├── .env.example                  # Environment template
├── Makefile                      # Common commands
├── README.md
└── .plans/                       # Planning documents
    ├── research/
    └── specs/
```

---

## 7. Authentication Architecture

### Session-Based Flow (MVP)

```
┌─────────┐                         ┌─────────┐                    ┌────────────┐
│  Client │                         │   API   │                    │  Database  │
└────┬────┘                         └────┬────┘                    └─────┬──────┘
     │                                   │                               │
     │  POST /auth/login                 │                               │
     │  {email, password}                │                               │
     │ ─────────────────────────────────>│                               │
     │                                   │                               │
     │                                   │  Verify password              │
     │                                   │ ─────────────────────────────>│
     │                                   │                               │
     │                                   │  Create session               │
     │                                   │ ─────────────────────────────>│
     │                                   │                               │
     │  Set-Cookie: session=xxx          │                               │
     │  HttpOnly; Secure; SameSite=Strict│                               │
     │ <─────────────────────────────────│                               │
     │                                   │                               │
     │  GET /subscriptions               │                               │
     │  Cookie: session=xxx              │                               │
     │ ─────────────────────────────────>│                               │
     │                                   │                               │
     │                                   │  Validate session             │
     │                                   │ ─────────────────────────────>│
     │                                   │                               │
     │  {subscriptions: [...]}           │                               │
     │ <─────────────────────────────────│                               │
     │                                   │                               │
```

### Session Token Strategy

- **Token Generation**: 32 bytes of cryptographically secure random data, hex encoded
- **Storage**: Hash of token stored in database (not the token itself)
- **Cookie Settings**: `HttpOnly`, `Secure`, `SameSite=Strict`, `Path=/`
- **Expiry**: 7 days, sliding window (extends on activity)
- **Rotation**: New token on privilege change (password update, etc.)

---

## 8. Offline Support Strategy

### Service Worker Caching

```javascript
// Cache strategy by resource type
const CACHE_STRATEGIES = {
  // App shell - cache first
  '/': 'cache-first',
  '/dashboard': 'cache-first',
  '/_next/static/*': 'cache-first',

  // API calls - network first with cache fallback
  '/api/v1/subscriptions': 'network-first',
  '/api/v1/analytics/*': 'network-first',

  // Static assets - cache first
  '/icons/*': 'cache-first',
  '/fonts/*': 'cache-first',
};
```

### Offline Data Flow

```
┌──────────────────────────────────────────────────────────────────┐
│                        ONLINE MODE                                │
│                                                                   │
│   [User Action] → [API Call] → [Server] → [Update UI]            │
│                                   ↓                               │
│                            [Update Cache]                         │
└──────────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────────┐
│                       OFFLINE MODE                                │
│                                                                   │
│   [User Action] → [Read Cache] → [Update UI]                     │
│         ↓                                                        │
│   [Queue Mutation] → [Sync when online]                          │
└──────────────────────────────────────────────────────────────────┘
```

### IndexedDB Schema

```typescript
// Offline storage structure
interface OfflineStore {
  subscriptions: Subscription[];
  pendingMutations: Mutation[];
  lastSyncAt: Date;
  userPreferences: Preferences;
}
```

---

## 9. Performance Considerations

### Frontend Performance

| Metric | Target | Strategy |
|--------|--------|----------|
| LCP | <2.5s | SSR critical content, optimize images |
| FID | <100ms | Code split, lazy load non-critical |
| CLS | <0.1 | Reserve space for dynamic content |
| TTI | <3s | Minimize JS bundle, defer non-essential |

### Backend Performance

| Metric | Target | Strategy |
|--------|--------|----------|
| API Response | <100ms p95 | Connection pooling, query optimization |
| Database | <50ms p95 | Proper indexing, prepared statements |
| Cold Start | <500ms | Minimal dependencies, compiled binary |

### Optimization Strategies

1. **Database**
   - Connection pooling (25 connections per instance)
   - Query result caching for analytics
   - Proper indexing (covered in schema)

2. **API**
   - Response compression (gzip/brotli)
   - ETag caching for subscription lists
   - Pagination for large result sets

3. **Frontend**
   - Static generation where possible
   - Image optimization (next/image)
   - Component lazy loading
   - Prefetch on hover for navigation

---

## 10. Deployment Architecture

### MVP Deployment (Single Server)

```
┌────────────────────────────────────────────────────────────────┐
│                      Vercel (Frontend)                          │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │                  Next.js Application                      │  │
│  │  - Edge Functions for middleware                          │  │
│  │  - Static assets on CDN                                   │  │
│  │  - ISR for semi-static pages                              │  │
│  └──────────────────────────────────────────────────────────┘  │
└────────────────────────────────────────────────────────────────┘
                              │
                              │ HTTPS
                              ▼
┌────────────────────────────────────────────────────────────────┐
│                    Railway (Backend)                            │
│  ┌──────────────────────┐  ┌────────────────────────────────┐  │
│  │    Go API Service    │  │        PostgreSQL              │  │
│  │    (containerized)   │──│        (managed)               │  │
│  │                      │  │                                │  │
│  └──────────────────────┘  └────────────────────────────────┘  │
└────────────────────────────────────────────────────────────────┘
```

### Docker Compose (Local Development)

```yaml
version: '3.8'

services:
  frontend:
    build: ./frontend
    ports:
      - "3000:3000"
    environment:
      - NEXT_PUBLIC_API_URL=http://localhost:8080
    depends_on:
      - api

  api:
    build: ./backend
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://user:pass@db:5432/subtrack
      - SESSION_SECRET=${SESSION_SECRET}
    depends_on:
      - db

  db:
    image: postgres:16-alpine
    ports:
      - "5432:5432"
    environment:
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=pass
      - POSTGRES_DB=subtrack
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

---

## 11. Trade-off Analysis

### Decisions Made

| Decision | Trade-off | Rationale |
|----------|-----------|-----------|
| **Session auth over JWT** | No stateless scaling | Simpler, more secure for MVP, can add JWT later |
| **PostgreSQL over NoSQL** | Less flexible schema | Strong consistency, better for financial data |
| **Go over Node.js** | Smaller ecosystem | Better performance, lower memory, strong typing |
| **Zustand over Redux** | Less middleware ecosystem | Simpler, sufficient for app complexity |
| **shadcn/ui over MUI** | Build your own components | No vendor lock-in, smaller bundle, full control |

### What We're Optimizing For

1. **Developer Experience** - Fast iteration, clear patterns
2. **User Experience** - Snappy performance, offline support
3. **Security** - Privacy-first, minimal attack surface
4. **Simplicity** - Minimal dependencies, clear architecture

### What We're Sacrificing

1. **Horizontal scaling** - Session store limits scaling (acceptable for MVP)
2. **Real-time sync** - No WebSocket for live updates (can add later)
3. **Complex querying** - SQL limits vs graph databases (unnecessary complexity)

---

## 12. Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Scope creep | High | High | Strict MVP features, defer nice-to-haves |
| Performance issues | Medium | Medium | Performance budgets, early testing |
| Security vulnerabilities | Medium | Critical | Security review, automated scanning |
| Database migration issues | Low | High | Backwards-compatible migrations |
| Third-party service outages | Low | Medium | Graceful degradation, offline support |

---

## 13. Implementation Priorities

### Phase 1: Core MVP
1. User authentication (register, login, logout)
2. Subscription CRUD
3. Dashboard with total spend
4. Basic category filtering
5. Mobile-responsive UI

### Phase 2: Enhanced Experience
1. PWA with offline support
2. Push notification reminders
3. Analytics dashboard
4. Dark mode
5. Data export

### Phase 3: Growth Features
1. Multi-currency support
2. Recurring reminder emails
3. Subscription templates (popular services)
4. Shared household subscriptions

---

## 14. Recommendations for Team

**DELEGATE:**
- **senior_dev**: Implement backend API structure, database setup, authentication system, input validation middleware
- **ui**: Build frontend component library based on UI research, implement design system with Tailwind
- **security**: Review authentication flow, set up security headers, implement rate limiting

---

## 15. Questions for Team Review

1. Should we support social login (Google/Apple) in MVP or defer?
2. Is PWA offline support a must-have for MVP or Phase 2?
3. Should we implement email verification before launch?
4. Do we need admin dashboard for MVP?

---

## Summary

This architecture provides a solid foundation for a privacy-first subscription management app:

- **Clear separation** between frontend and backend
- **Type-safe** from database to UI with TypeScript and Go
- **Secure by design** with session-based auth and proper validation
- **Mobile-first** with PWA capabilities
- **Scalable** for future growth while keeping MVP simple

The fullstack template choice is validated - we need persistent storage, authentication, and real-time calculations that a static site can't provide.

---

Status: PHASE_COMPLETE: research
