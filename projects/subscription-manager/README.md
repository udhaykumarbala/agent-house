# SubTrack - Subscription Manager

A privacy-first, mobile-first subscription management web application. Track all your subscriptions in one beautiful place without connecting your bank account.

## Features

- **Dashboard** - View total monthly/yearly spending at a glance
- **Subscription Management** - Add, edit, and delete subscriptions
- **Popular Templates** - Quick-add common services (Netflix, Spotify, etc.)
- **Upcoming Renewals** - See what's due in the next 30 days
- **Category Filtering** - Organize by Entertainment, Productivity, etc.
- **Dark Mode** - Full dark theme support
- **Mobile-First** - Optimized for phones, responsive on all devices

## Tech Stack

### Frontend
- Next.js 14 (App Router)
- TypeScript
- Tailwind CSS
- React Query (TanStack)
- Zustand (state management)
- Lucide React (icons)

### Backend
- Go 1.21+
- Chi Router
- PostgreSQL 16
- Argon2id password hashing

## Quick Start

### Prerequisites
- Node.js 20+
- Go 1.21+
- Docker & Docker Compose
- PostgreSQL (or use Docker)

### Development Setup

1. **Clone and navigate:**
   ```bash
   cd subscription-manager
   ```

2. **Start PostgreSQL:**
   ```bash
   docker-compose up -d postgres
   ```

3. **Backend setup:**
   ```bash
   cd backend
   cp .env.example .env
   go mod download
   make run
   ```

4. **Frontend setup (new terminal):**
   ```bash
   cd frontend
   cp .env.local.example .env.local
   npm install
   npm run dev
   ```

5. **Open the app:**
   - Frontend: http://localhost:3000
   - API: http://localhost:8080
   - API Test: Open `api_test.html` in browser

### Docker (Full Stack)

```bash
# Build and run all services
docker-compose up --build

# Or use make
make build
make up
```

## Project Structure

```
subscription-manager/
├── frontend/                 # Next.js frontend
│   ├── src/
│   │   ├── app/             # App Router pages
│   │   ├── components/      # React components
│   │   ├── hooks/           # Custom hooks
│   │   ├── lib/             # Utilities, API client
│   │   ├── stores/          # Zustand stores
│   │   └── types/           # TypeScript types
│   └── public/              # Static assets
├── backend/                  # Go API
│   ├── cmd/api/             # Main entry point
│   └── internal/
│       ├── config/          # Configuration
│       ├── database/        # DB connection & migrations
│       ├── handler/         # HTTP handlers
│       ├── model/           # Data models
│       ├── repository/      # Data access
│       └── service/         # Business logic
├── docker-compose.yml       # Docker services
├── swagger.yaml             # API specification
├── api_test.html            # API testing UI
└── uiflow.md               # UI flow documentation
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | Register new user |
| POST | `/api/v1/auth/login` | Login user |
| POST | `/api/v1/auth/logout` | Logout user |
| GET | `/api/v1/auth/me` | Get current user |
| GET | `/api/v1/subscriptions` | List subscriptions |
| POST | `/api/v1/subscriptions` | Create subscription |
| GET | `/api/v1/subscriptions/:id` | Get subscription |
| PUT | `/api/v1/subscriptions/:id` | Update subscription |
| DELETE | `/api/v1/subscriptions/:id` | Delete subscription |
| GET | `/api/v1/analytics/summary` | Get spending summary |
| GET | `/api/v1/analytics/upcoming` | Get upcoming renewals |

## Design System

### Colors
- **Primary (Coral)**: `#FF7A5C` - Warm, memorable brand color
- **Background**: `#FAFAFA` (light) / `#0F0F12` (dark)
- **Surface**: `#FFFFFF` (light) / `#1A1A1E` (dark)

### Typography
- Font: Geist / Inter
- Tabular numbers for prices

## Security

- Session-based authentication with HttpOnly cookies
- Argon2id password hashing (64MB memory, 3 iterations)
- Rate limiting on auth endpoints (10 req/min)
- CORS restricted to frontend origin
- Security headers (HSTS, CSP, X-Frame-Options)
- Parameterized SQL queries (no injection)

## License

MIT
