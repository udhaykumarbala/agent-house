# Architecture Research: Inventory Reorder Alert System

## Executive Summary

This document provides architectural research and technical design considerations for an alert-first inventory reorder system. Based on the product requirements (5-minute setup, multi-channel alerts, one-click reorder) and security constraints (multi-tenant isolation, SSRF protection, audit logging), this research evaluates architecture patterns, technology choices, and data models.

---

## 1. System Architecture Patterns Analysis

### 1.1 Monolith vs Microservices Evaluation

| Approach | Pros | Cons | Fit for MVP |
|----------|------|------|-------------|
| **Monolith** | Simple deployment, easy debugging, lower operational overhead, faster initial development | Harder to scale individual components, tech stack lock-in | **Recommended** |
| **Microservices** | Independent scaling, tech flexibility per service, team autonomy | Complex operations, distributed tracing needed, overkill for small team | Not recommended for MVP |
| **Modular Monolith** | Best of both - modular code with simple deployment, easy refactor to services later | Requires discipline to maintain boundaries | **Good alternative** |

**Recommendation**: Start with a **Modular Monolith** architecture. The codebase should be organized into clear modules (inventory, alerts, notifications, users) that can be extracted to services later if needed. This matches the product's "start simple" philosophy.

### 1.2 Competitor Architecture Patterns

| System | Architecture | Stack | Learnings |
|--------|--------------|-------|-----------|
| **Sortly** | Mobile-first SPA + API | React Native, Rails, PostgreSQL | Mobile-first pays off for warehouse workers |
| **inFlow** | Desktop client + cloud sync | .NET, SQL Server | Avoid desktop-first - limits adoption |
| **Zoho Inventory** | Multi-tenant SaaS | Java, MongoDB | Heavy framework slows iteration |
| **Stockpile** | Simple SPA | React, Node.js, PostgreSQL | Simplicity works for SMB market |
| **TradeGecko** | API-first platform | Ruby, PostgreSQL | API-first enables integrations |

**Key Insight**: Successful SMB inventory tools are API-first, use PostgreSQL for reliability, and prioritize mobile experience.

### 1.3 Event-Driven vs Request-Response

For an alert system, event-driven architecture is natural:

```
[Stock Update] → [Event Bus] → [Threshold Checker] → [Alert Generator] → [Notification Dispatcher]
                     ↓
              [Audit Logger]
```

**Hybrid Approach Recommended**:
- Request-response for CRUD operations (items, thresholds, users)
- Event-driven for alert processing and notifications
- Background job queue for email/SMS/webhook delivery

---

## 2. Technology Stack Analysis

### 2.1 Backend Framework Comparison

| Framework | Language | Pros | Cons | Recommendation |
|-----------|----------|------|------|----------------|
| **Express.js** | Node.js | Lightweight, huge ecosystem, async I/O | No structure, callback patterns | Possible |
| **Fastify** | Node.js | Fast, schema validation, modern | Smaller community than Express | Good option |
| **NestJS** | TypeScript | Enterprise structure, DI, modular | Learning curve, verbose | **Recommended** |
| **Go (Gin/Echo)** | Go | Fast, low memory, great for APIs | Smaller web ecosystem | Good option |
| **FastAPI** | Python | Async, auto-docs, type hints | Python GIL limits concurrency | Good option |
| **Rails** | Ruby | Rapid development, conventions | Performance concerns at scale | Possible |

**Recommendation**: **NestJS with TypeScript**
- TypeScript provides type safety for inventory quantities and business logic
- Modular architecture aligns with our modular monolith approach
- Built-in validation, guards, and interceptors for security
- Good WebSocket support for real-time dashboard updates
- Large ecosystem for notification integrations

### 2.2 Frontend Framework Comparison

| Framework | Pros | Cons | Recommendation |
|-----------|------|------|----------------|
| **React** | Largest ecosystem, mature, great tooling | JSX learning curve, no conventions | Good option |
| **Next.js** | SSR/SSG, API routes, great DX | Overkill for SPA, Vercel-centric | **Recommended** |
| **Vue 3** | Gentler learning curve, reactive, Composition API | Smaller ecosystem | Possible |
| **SvelteKit** | Fastest, smallest bundles, intuitive | Smaller ecosystem, newer | Future consideration |

**Recommendation**: **Next.js 14+ with App Router**
- Server components reduce client bundle size
- API routes can serve as BFF (Backend for Frontend)
- Excellent TypeScript integration
- Easy deployment to Vercel/AWS
- React ecosystem for component libraries

### 2.3 Mobile Strategy

| Approach | Pros | Cons | Recommendation |
|----------|------|------|----------------|
| **PWA** | Single codebase, no app store, auto-updates | Limited native features, iOS restrictions | **Recommended for MVP** |
| **React Native** | Native feel, code sharing | Two codebases to maintain, complex setup | Future v2 |
| **Flutter** | Fast UI, single codebase | Dart learning curve, larger bundles | Alternative |
| **Native iOS/Android** | Best performance, full features | Most expensive, longest development | Not recommended |

**Recommendation**: **PWA first** with Next.js
- Matches "5-minute setup" philosophy (no app store download)
- Push notifications work on Android, improving on iOS
- Barcode scanning via browser camera API
- Consider React Native for v2 if native features needed

### 2.4 Database Selection

| Database | Type | Pros | Cons | Fit |
|----------|------|------|------|-----|
| **PostgreSQL** | Relational | ACID, mature, JSON support, row-level security | Vertical scaling | **Recommended** |
| **MySQL** | Relational | Simple, widespread, good performance | Fewer advanced features | Alternative |
| **MongoDB** | Document | Flexible schema, horizontal scaling | Consistency challenges, injection risks | Not recommended |
| **SQLite** | Embedded | Zero config, great for edge | No concurrent writes, limited features | Dev/test only |

**Recommendation**: **PostgreSQL 15+**
- Row-Level Security (RLS) for multi-tenant isolation (security requirement)
- JSONB for flexible item attributes
- Full-text search for item catalog
- Excellent with Prisma ORM
- Proven reliability for inventory data

### 2.5 Message Queue / Background Jobs

| Technology | Pros | Cons | Recommendation |
|------------|------|------|----------------|
| **BullMQ** | Redis-based, TypeScript, feature-rich | Requires Redis | **Recommended** |
| **RabbitMQ** | Enterprise-grade, routing, acknowledgments | Operational complexity | Overkill for MVP |
| **AWS SQS** | Managed, scales infinitely | AWS lock-in, latency | Alternative |
| **Inngest** | Serverless functions, built-in retry | Newer, vendor dependency | Interesting option |

**Recommendation**: **BullMQ with Redis**
- TypeScript native, integrates well with NestJS
- Handles notification job queuing (retry, delay, rate limiting)
- Redis also serves as cache and session store
- Can migrate to managed Redis (Upstash, ElastiCache) easily

### 2.6 Notification Infrastructure

| Service | Type | Pricing | Reliability | Recommendation |
|---------|------|---------|-------------|----------------|
| **SendGrid** | Email | Free tier to $20/mo | High | **Recommended** |
| **Resend** | Email | Free tier, modern API | High | Alternative |
| **Twilio** | SMS | Pay per message | High | **Recommended** |
| **Slack API** | Chat | Free | High | **Recommended** |
| **Microsoft Graph** | Teams | Free with Azure | Medium | For Teams integration |

---

## 3. Data Model Design

### 3.1 Core Entities

```
┌─────────────────────────────────────────────────────────────────────┐
│                           DATA MODEL                                 │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌──────────────┐       ┌──────────────┐       ┌──────────────┐    │
│  │   Tenant     │───────│    User      │       │   Category   │    │
│  │              │   1:N │              │       │              │    │
│  │ - id (uuid)  │       │ - id (uuid)  │       │ - id (uuid)  │    │
│  │ - name       │       │ - tenant_id  │       │ - tenant_id  │    │
│  │ - plan       │       │ - email      │       │ - name       │    │
│  │ - settings   │       │ - role       │       │ - parent_id  │    │
│  └──────────────┘       │ - password   │       └──────────────┘    │
│         │               └──────────────┘              │             │
│         │                      │                      │             │
│         │ 1:N                  │                      │ 1:N         │
│         ▼                      ▼                      ▼             │
│  ┌──────────────┐       ┌──────────────┐       ┌──────────────┐    │
│  │    Item      │───────│  StockUpdate │       │    Alert     │    │
│  │              │   1:N │              │       │              │    │
│  │ - id (uuid)  │       │ - id (uuid)  │       │ - id (uuid)  │    │
│  │ - tenant_id  │       │ - item_id    │       │ - tenant_id  │    │
│  │ - name       │       │ - quantity   │       │ - item_id    │    │
│  │ - sku        │       │ - prev_qty   │       │ - type       │    │
│  │ - quantity   │       │ - reason     │       │ - status     │    │
│  │ - threshold  │       │ - user_id    │       │ - triggered  │    │
│  │ - max_stock  │       │ - timestamp  │       │ - resolved   │    │
│  │ - unit       │       └──────────────┘       │ - resolved_by│    │
│  │ - category_id│                              └──────────────┘    │
│  │ - attributes │                                     │             │
│  └──────────────┘                                     │ 1:N         │
│         │                                             ▼             │
│         │                                      ┌──────────────┐    │
│         │ 1:N                                  │ Notification │    │
│         ▼                                      │              │    │
│  ┌──────────────┐                              │ - id (uuid)  │    │
│  │   Supplier   │                              │ - alert_id   │    │
│  │              │                              │ - channel    │    │
│  │ - id (uuid)  │                              │ - recipient  │    │
│  │ - tenant_id  │                              │ - status     │    │
│  │ - name       │                              │ - sent_at    │    │
│  │ - email      │                              │ - error      │    │
│  │ - phone      │                              └──────────────┘    │
│  │ - webhook_url│                                                  │
│  └──────────────┘                                                  │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

### 3.2 Key Relationships

| Relationship | Cardinality | Notes |
|--------------|-------------|-------|
| Tenant → User | 1:N | Users belong to one tenant |
| Tenant → Item | 1:N | Items are tenant-scoped |
| Item → StockUpdate | 1:N | Audit trail of quantity changes |
| Item → Alert | 1:N | One item can have multiple historical alerts |
| Alert → Notification | 1:N | One alert can send multiple notifications |
| Category → Item | 1:N | Optional categorization |
| Item → Supplier | N:M | Items can have preferred suppliers |

### 3.3 Multi-Tenancy Pattern

**Recommended: Shared Database, Tenant Column**

```sql
-- Every query includes tenant filter
CREATE POLICY tenant_isolation ON items
  FOR ALL
  USING (tenant_id = current_setting('app.tenant_id')::uuid);

-- Application sets tenant context
SET app.tenant_id = 'uuid-of-tenant';
SELECT * FROM items;  -- Only sees tenant's items
```

**Why this pattern**:
- Simple implementation with Prisma/TypeORM
- PostgreSQL RLS provides defense-in-depth
- Single database reduces operational complexity
- Easy backup and restore per tenant (filtered export)

---

## 4. API Design

### 4.1 RESTful Resource Design

```
Base URL: /api/v1

Authentication:
  POST   /auth/login              # Login, returns JWT
  POST   /auth/register           # Register new tenant + admin user
  POST   /auth/refresh            # Refresh JWT token
  POST   /auth/forgot-password    # Initiate password reset
  POST   /auth/reset-password     # Complete password reset

Items:
  GET    /items                   # List items (paginated, filterable)
  POST   /items                   # Create item
  GET    /items/:id               # Get item detail
  PUT    /items/:id               # Update item
  DELETE /items/:id               # Delete item
  PATCH  /items/:id/quantity      # Quick quantity update
  POST   /items/bulk              # Bulk import

Alerts:
  GET    /alerts                  # List active alerts
  GET    /alerts/history          # Alert history
  PATCH  /alerts/:id/resolve      # Mark alert resolved
  PATCH  /alerts/:id/snooze       # Snooze alert

Categories:
  GET    /categories              # List categories
  POST   /categories              # Create category
  PUT    /categories/:id          # Update category
  DELETE /categories/:id          # Delete category

Settings:
  GET    /settings/notifications  # Get notification preferences
  PUT    /settings/notifications  # Update notification preferences
  GET    /settings/integrations   # List integrations
  POST   /settings/integrations   # Add integration (Slack, etc.)

Analytics:
  GET    /analytics/consumption   # Consumption trends
  GET    /analytics/alerts        # Alert statistics

Webhooks (for integrations):
  POST   /webhooks/incoming/:provider  # Receive external webhooks
```

### 4.2 Real-Time Updates

**WebSocket Events**:
```typescript
// Client subscribes to tenant channel
socket.on('alert:new', (alert) => { /* New alert triggered */ });
socket.on('alert:resolved', (alertId) => { /* Alert resolved */ });
socket.on('item:updated', (item) => { /* Stock level changed */ });
socket.on('item:critical', (item) => { /* Item hit critical level */ });
```

**Implementation**: Socket.IO with Redis adapter for horizontal scaling

### 4.3 API Versioning Strategy

- URL-based versioning: `/api/v1/`, `/api/v2/`
- Maintain backward compatibility within major versions
- Deprecation headers for sunset features
- OpenAPI 3.0 specification with SwaggerUI

---

## 5. Notification Architecture

### 5.1 Notification Flow

```
┌─────────────────────────────────────────────────────────────────────┐
│                     NOTIFICATION PIPELINE                            │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  [Stock Update]                                                      │
│       │                                                              │
│       ▼                                                              │
│  ┌─────────────┐    No     ┌─────────────┐                          │
│  │  Threshold  │──────────▶│   No-op     │                          │
│  │   Check     │           └─────────────┘                          │
│  └─────────────┘                                                     │
│       │ Yes                                                          │
│       ▼                                                              │
│  ┌─────────────┐                                                     │
│  │   Create    │                                                     │
│  │   Alert     │                                                     │
│  └─────────────┘                                                     │
│       │                                                              │
│       ▼                                                              │
│  ┌─────────────┐                                                     │
│  │   Check     │  Cooldown?   ┌─────────────┐                       │
│  │  Cooldown   │─────────────▶│ Skip notify │                       │
│  └─────────────┘              └─────────────┘                        │
│       │ No                                                           │
│       ▼                                                              │
│  ┌─────────────┐    ┌─────────────────────────────────────┐         │
│  │   Queue     │───▶│          NOTIFICATION JOBS          │         │
│  │   Jobs      │    │                                     │         │
│  └─────────────┘    │  ┌─────┐ ┌─────┐ ┌─────┐ ┌───────┐ │         │
│                     │  │Email│ │ SMS │ │Slack│ │Webhook│ │         │
│                     │  └─────┘ └─────┘ └─────┘ └───────┘ │         │
│                     └─────────────────────────────────────┘         │
│                                    │                                 │
│                                    ▼                                 │
│                     ┌─────────────────────────────────────┐         │
│                     │     DELIVERY with RETRY LOGIC       │         │
│                     │  - Exponential backoff              │         │
│                     │  - Dead letter queue                │         │
│                     │  - Delivery confirmation            │         │
│                     └─────────────────────────────────────┘         │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

### 5.2 Alert Coalescing

To prevent alert fatigue:

```typescript
// Alert rules
const alertRules = {
  cooldownPeriod: '4 hours',      // Don't re-alert same item within period
  digestMode: 'daily',             // Option: batch into daily digest
  criticalBypass: true,            // Critical items bypass cooldown
  maxAlertsPerDay: 50,             // Rate limit per tenant
  businessHours: {                 // Optional: only alert during work hours
    enabled: true,
    start: '08:00',
    end: '18:00',
    timezone: 'America/New_York'
  }
};
```

### 5.3 Webhook Security (SSRF Prevention)

Based on security research requirements:

```typescript
// Webhook URL validation
const validateWebhookUrl = async (url: string): Promise<boolean> => {
  // 1. Parse and validate URL format
  const parsed = new URL(url);

  // 2. Only allow HTTPS in production
  if (process.env.NODE_ENV === 'production' && parsed.protocol !== 'https:') {
    throw new Error('HTTPS required for webhooks');
  }

  // 3. Resolve DNS to get IP
  const addresses = await dns.resolve4(parsed.hostname);

  // 4. Check against blocked IP ranges
  const blockedRanges = [
    '10.0.0.0/8',      // Private
    '172.16.0.0/12',   // Private
    '192.168.0.0/16',  // Private
    '127.0.0.0/8',     // Loopback
    '169.254.0.0/16',  // Link-local (cloud metadata)
  ];

  for (const ip of addresses) {
    if (isInRange(ip, blockedRanges)) {
      throw new Error('Webhook URL resolves to private IP');
    }
  }

  return true;
};
```

---

## 6. Scalability Considerations

### 6.1 Horizontal Scaling Path

```
MVP (0-1000 tenants):
┌─────────┐    ┌─────────┐    ┌─────────┐
│ Next.js │───▶│ NestJS  │───▶│PostgreSQL│
│   PWA   │    │   API   │    │  + Redis │
└─────────┘    └─────────┘    └─────────┘

Growth (1000-10000 tenants):
┌─────────┐    ┌─────────┐    ┌─────────────┐
│   CDN   │───▶│Load     │───▶│ API Cluster │───▶│ RDS + ElastiCache
│(Vercel) │    │Balancer │    │ (3 nodes)   │    │ (managed)
└─────────┘    └─────────┘    └─────────────┘    └─────────────┘
                                    │
                                    ▼
                           ┌───────────────┐
                           │ Worker Cluster│
                           │ (notifications)│
                           └───────────────┘
```

### 6.2 Performance Targets

| Metric | Target | Approach |
|--------|--------|----------|
| API response time (p95) | < 200ms | Database indexing, caching |
| Dashboard load time | < 1s | Server components, edge caching |
| Notification delivery | < 30s | Dedicated worker queue |
| Concurrent users | 1000+ | Horizontal API scaling |
| Items per tenant | 100,000+ | Pagination, cursor-based |

### 6.3 Caching Strategy

```typescript
// Redis caching layers
const cacheConfig = {
  // Tenant settings - changes rarely
  tenantSettings: { ttl: '1 hour', key: 'tenant:{id}:settings' },

  // Item list - moderate changes
  itemList: { ttl: '5 minutes', key: 'tenant:{id}:items:page:{page}' },

  // Dashboard stats - computed periodically
  dashboardStats: { ttl: '1 minute', key: 'tenant:{id}:stats' },

  // Alert counts - real-time important
  alertCounts: { ttl: '30 seconds', key: 'tenant:{id}:alerts:count' },
};
```

---

## 7. Infrastructure & Deployment

### 7.1 Recommended Stack

| Layer | Technology | Rationale |
|-------|------------|-----------|
| **Hosting** | Vercel (frontend) + Railway/Fly.io (backend) | Simple deployment, good free tier |
| **Database** | Supabase PostgreSQL or Neon | Managed PostgreSQL, serverless options |
| **Redis** | Upstash | Serverless Redis, pay-per-use |
| **File Storage** | Cloudflare R2 | S3-compatible, no egress fees |
| **Email** | SendGrid | Reliable, good free tier |
| **Monitoring** | Sentry + Better Stack | Error tracking + uptime |

### 7.2 Alternative: AWS Stack

| Layer | Technology | When to Use |
|-------|------------|-------------|
| **Compute** | ECS Fargate | Scale requirements exceed Railway |
| **Database** | RDS PostgreSQL | Need enterprise SLA |
| **Cache** | ElastiCache Redis | High throughput needed |
| **Queue** | SQS + Lambda | Serverless job processing |
| **CDN** | CloudFront | Global distribution needed |

### 7.3 CI/CD Pipeline

```yaml
# GitHub Actions workflow
stages:
  1. Lint & Type Check
  2. Unit Tests
  3. Integration Tests (test DB)
  4. Security Scan (Snyk)
  5. Build Docker images
  6. Deploy to staging
  7. E2E Tests (Playwright)
  8. Deploy to production (manual approval)
```

---

## 8. Security Implementation Summary

Based on security research, key architectural security controls:

| Control | Implementation |
|---------|----------------|
| **Authentication** | JWT + refresh tokens, Argon2id passwords |
| **Authorization** | RBAC with tenant isolation middleware |
| **Multi-tenancy** | Tenant ID in JWT, PostgreSQL RLS |
| **Input Validation** | Zod schemas on all endpoints |
| **Rate Limiting** | Redis-based, per-tenant quotas |
| **SSRF Prevention** | URL validation, IP blocklist |
| **Audit Logging** | Structured logs to persistent storage |
| **Encryption** | TLS 1.3, AES-256-GCM for secrets at rest |

---

## 9. Trade-offs Summary

| Decision | Optimizing For | Sacrificing |
|----------|---------------|-------------|
| Modular monolith | Development speed, simplicity | Independent service scaling |
| PostgreSQL over MongoDB | Data integrity, security (RLS) | Schema flexibility |
| PWA over native | Faster launch, single codebase | Some native features |
| BullMQ over SQS | Development experience, cost | AWS ecosystem lock-in |
| TypeScript everywhere | Type safety, maintainability | Initial development time |

---

## 10. Recommended Tech Stack Summary

### Final Stack

| Layer | Choice |
|-------|--------|
| **Frontend** | Next.js 14 (App Router) + TypeScript |
| **Backend** | NestJS + TypeScript |
| **Database** | PostgreSQL 15 (Supabase/Neon) |
| **Cache/Queue** | Redis + BullMQ (Upstash) |
| **ORM** | Prisma |
| **Auth** | Custom JWT or Lucia Auth |
| **Email** | SendGrid |
| **SMS** | Twilio |
| **Real-time** | Socket.IO with Redis adapter |
| **Hosting** | Vercel (frontend) + Railway (backend) |
| **Monitoring** | Sentry + Better Stack |

### Development Tools

| Tool | Purpose |
|------|---------|
| **pnpm** | Package management (monorepo) |
| **Turborepo** | Monorepo build system |
| **Prisma** | Database ORM + migrations |
| **Zod** | Schema validation |
| **Swagger/OpenAPI** | API documentation |
| **Playwright** | E2E testing |
| **Vitest** | Unit testing |

---

## 11. Project Structure Recommendation

```
inventory-alerts/
├── apps/
│   ├── web/                    # Next.js frontend
│   │   ├── app/
│   │   │   ├── (auth)/         # Auth pages (login, register)
│   │   │   ├── (dashboard)/    # Dashboard pages
│   │   │   ├── api/            # API routes (BFF)
│   │   │   └── layout.tsx
│   │   ├── components/
│   │   └── lib/
│   │
│   └── api/                    # NestJS backend
│       ├── src/
│       │   ├── modules/
│       │   │   ├── auth/
│       │   │   ├── items/
│       │   │   ├── alerts/
│       │   │   ├── notifications/
│       │   │   └── users/
│       │   ├── common/         # Guards, filters, pipes
│       │   └── main.ts
│       └── prisma/
│           └── schema.prisma
│
├── packages/
│   ├── shared/                 # Shared types, utils
│   ├── ui/                     # Shared UI components
│   └── config/                 # Shared config (eslint, tsconfig)
│
├── docker-compose.yml
├── turbo.json
└── package.json
```

---

## Summary

The Inventory Reorder Alert System should be built as a **modular monolith** using **NestJS + Next.js + PostgreSQL + Redis**. This stack provides:

1. **Type safety** across the entire codebase
2. **Multi-tenant security** via PostgreSQL RLS
3. **Reliable notifications** via queued background jobs
4. **Real-time updates** for dashboard experience
5. **Simple deployment** with modern hosting platforms
6. **Clear scaling path** when growth demands it

Key architectural priorities align with product goals:
- **5-minute setup** → Simple PWA, no app store
- **Multi-channel alerts** → Queue-based notification pipeline
- **One-click reorder** → Fast API, optimistic UI
- **Security** → Tenant isolation, SSRF protection, audit logging

---

## DELEGATE:

- **senior_dev**: Implement the core backend with NestJS - auth module, items module, alerts module, notification queue with BullMQ
- **ui**: Design component library based on UI research - stock level bars, alert cards, dashboard layout
- **security**: Review final architecture for security gaps, validate multi-tenant isolation approach

---

**Status**: PHASE_COMPLETE: research
