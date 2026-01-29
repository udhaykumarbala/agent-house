package template

// TemplateType represents the type of project template
type TemplateType string

const (
	StaticHTML     TemplateType = "static-html"
	StaticEnhanced TemplateType = "static-enhanced"
	NextJSFrontend TemplateType = "nextjs-frontend"
	GoAPI          TemplateType = "go-api"
	Fullstack      TemplateType = "fullstack"
)

// DefaultTemplate is used when no template is specified
var DefaultTemplate = StaticEnhanced

// Template defines a project template with its configuration
type Template struct {
	Type          TemplateType
	Name          string
	Description   string
	Complexity    string
	TechStack     []string
	RequiredFiles []string
	Guidelines    string
}

// Registry holds all available templates
var Registry = map[TemplateType]Template{
	StaticHTML:     staticHTMLTemplate,
	StaticEnhanced: staticEnhancedTemplate,
	NextJSFrontend: nextjsFrontendTemplate,
	GoAPI:          goAPITemplate,
	Fullstack:      fullstackTemplate,
}

// GetTemplate returns a template by type
func GetTemplate(t TemplateType) Template {
	if tmpl, ok := Registry[t]; ok {
		return tmpl
	}
	return Registry[DefaultTemplate]
}

// GetGuidelines returns the guidelines for a template
func GetGuidelines(t TemplateType) string {
	return GetTemplate(t).Guidelines
}

// AllTemplates returns a list of all available template types
func AllTemplates() []TemplateType {
	return []TemplateType{
		StaticHTML,
		StaticEnhanced,
		NextJSFrontend,
		GoAPI,
		Fullstack,
	}
}

// ParseTemplateType converts a string to TemplateType
func ParseTemplateType(s string) TemplateType {
	switch s {
	case "static-html":
		return StaticHTML
	case "static-enhanced":
		return StaticEnhanced
	case "nextjs-frontend":
		return NextJSFrontend
	case "go-api":
		return GoAPI
	case "fullstack":
		return Fullstack
	default:
		return DefaultTemplate
	}
}

// Template Definitions

var staticHTMLTemplate = Template{
	Type:        StaticHTML,
	Name:        "Static HTML",
	Description: "Simple static HTML/CSS/JS site for landing pages, portfolios, and prototypes",
	Complexity:  "Simple",
	TechStack:   []string{"HTML5", "CSS3", "Vanilla JavaScript", "Lucide Icons (CDN)", "Google Fonts (CDN)"},
	RequiredFiles: []string{
		"index.html",
		"styles/variables.css",
		"styles/main.css",
		"scripts/app.js",
		"README.md",
	},
	Guidelines: staticHTMLGuidelines,
}

var staticEnhancedTemplate = Template{
	Type:        StaticEnhanced,
	Name:        "Static Enhanced",
	Description: "Modern static site with Vite, Tailwind CSS, and TypeScript",
	Complexity:  "Basic",
	TechStack:   []string{"Vite 5.x", "Tailwind CSS 3.x", "TypeScript 5.x", "Lucide Icons", "Fontsource"},
	RequiredFiles: []string{
		"package.json",
		"tailwind.config.js",
		"tsconfig.json",
		"vite.config.ts",
		"src/index.html",
		"src/styles/main.css",
		"src/scripts/main.ts",
		"README.md",
	},
	Guidelines: staticEnhancedGuidelines,
}

var nextjsFrontendTemplate = Template{
	Type:        NextJSFrontend,
	Name:        "Next.js Frontend",
	Description: "Full-featured Next.js 14 application with App Router and shadcn/ui",
	Complexity:  "Medium",
	TechStack:   []string{"Next.js 14", "TypeScript 5.x", "Tailwind CSS 3.x", "shadcn/ui", "Zustand", "TanStack Query", "React Hook Form", "Zod"},
	RequiredFiles: []string{
		"package.json",
		"next.config.js",
		"tailwind.config.ts",
		"tsconfig.json",
		"components.json",
		"src/app/layout.tsx",
		"src/app/page.tsx",
		"src/app/globals.css",
		"src/lib/utils.ts",
		".env.example",
		"README.md",
	},
	Guidelines: nextjsFrontendGuidelines,
}

var goAPITemplate = Template{
	Type:        GoAPI,
	Name:        "Go API",
	Description: "REST API with Go, Chi router, and PostgreSQL",
	Complexity:  "Medium",
	TechStack:   []string{"Go 1.21+", "Chi 5.x", "PostgreSQL + pgx", "golang-migrate", "Viper", "slog", "go-playground/validator", "testify"},
	RequiredFiles: []string{
		"go.mod",
		"cmd/api/main.go",
		"internal/config/config.go",
		"internal/handler/handler.go",
		"internal/service/service.go",
		"internal/repository/repository.go",
		"Makefile",
		"Dockerfile",
		"docker-compose.yml",
		".env.example",
		"migrations/000001_init.up.sql",
		"docs/swagger.yaml",
		"README.md",
	},
	Guidelines: goAPIGuidelines,
}

var fullstackTemplate = Template{
	Type:        Fullstack,
	Name:        "Fullstack",
	Description: "Complete application with Next.js frontend, Go backend, and Docker",
	Complexity:  "Complex",
	TechStack:   []string{"Next.js 14", "Go 1.21+", "Chi", "PostgreSQL 16", "Redis 7", "Docker Compose", "TypeScript", "Tailwind CSS"},
	RequiredFiles: []string{
		"docker-compose.yml",
		"docker-compose.prod.yml",
		"Makefile",
		"swagger.yaml",
		"README.md",
		"frontend/package.json",
		"frontend/Dockerfile",
		"frontend/next.config.js",
		"frontend/tailwind.config.ts",
		"frontend/src/lib/api-client.ts",
		"frontend/.env.example",
		"backend/go.mod",
		"backend/Dockerfile",
		"backend/Makefile",
		"backend/cmd/api/main.go",
		"backend/internal/handler/routes.go",
		"backend/migrations/000001_init.up.sql",
		"backend/.env.example",
	},
	Guidelines: fullstackGuidelines,
}

// Guidelines as constants

const staticHTMLGuidelines = `## Template: static-html

### File Structure
` + "```" + `
project/
├── index.html              # Main entry point
├── styles/
│   ├── main.css           # Primary styles
│   ├── variables.css      # CSS custom properties (design tokens)
│   └── responsive.css     # Media queries
├── scripts/
│   ├── app.js             # Main application logic
│   └── utils.js           # Helper functions
├── assets/
│   ├── images/
│   └── icons/
├── api_test.html
├── uiflow.md
└── README.md
` + "```" + `

### Rules
1. **File Structure**: Follow the exact structure above
2. **CSS Variables**: Define ALL colors in variables.css
3. **No Build Tools**: Do not introduce npm, webpack, or bundlers
4. **CDN Only**: Use CDN for external libraries (no npm install)
5. **Progressive Enhancement**: JS should enhance, not break without it
6. **Mobile First**: Write mobile styles first, then desktop overrides

### Color Implementation
Copy EXACT hex codes from .plans/specs/ui-spec.md:
` + "```css" + `
/* variables.css */
:root {
  --color-primary: #[FROM-UI-SPEC];
  --color-background: #[FROM-UI-SPEC];
  --color-text: #[FROM-UI-SPEC];
  --color-surface: #[FROM-UI-SPEC];
}
` + "```" + `

### Build & Deploy
- No build step required
- Development: npx live-server --port=3000
- Deploy: Copy files to any static host
`

const staticEnhancedGuidelines = `## Template: static-enhanced

### File Structure
` + "```" + `
project/
├── src/
│   ├── index.html
│   ├── styles/
│   │   ├── main.css
│   │   └── components.css
│   ├── scripts/
│   │   ├── main.ts
│   │   ├── utils.ts
│   │   └── components/
│   └── assets/
│       └── images/
├── public/
│   ├── favicon.ico
│   └── robots.txt
├── dist/
├── package.json
├── tsconfig.json
├── tailwind.config.js
├── vite.config.ts
├── api_test.html
├── uiflow.md
└── README.md
` + "```" + `

### Rules
1. **Use Tailwind Classes**: Prefer utility classes over custom CSS
2. **Extend Theme**: Add colors to tailwind.config.js, NOT inline hex
3. **TypeScript Required**: All .ts files, no .js
4. **Component Pattern**: Extract reusable components to scripts/components/
5. **Build Must Pass**: Run npm run build before COMPLETE

### package.json
` + "```json" + `
{
  "name": "project-name",
  "version": "1.0.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "tsc && vite build",
    "preview": "vite preview",
    "lint": "eslint src --ext ts",
    "test": "vitest"
  },
  "devDependencies": {
    "typescript": "^5.3.0",
    "vite": "^5.0.0",
    "tailwindcss": "^3.4.0",
    "postcss": "^8.4.0",
    "autoprefixer": "^10.4.0",
    "eslint": "^8.56.0",
    "vitest": "^1.0.0"
  }
}
` + "```" + `

### Color Implementation
Add to tailwind.config.js theme.extend.colors:
` + "```javascript" + `
// tailwind.config.js
export default {
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: '#[FROM-UI-SPEC]',
          hover: '#[FROM-UI-SPEC]',
        },
        background: '#[FROM-UI-SPEC]',
        surface: '#[FROM-UI-SPEC]',
      },
    },
  },
}
` + "```" + `

Use classes like bg-primary, NOT bg-[#hexcode]
`

const nextjsFrontendGuidelines = `## Template: nextjs-frontend

### File Structure
` + "```" + `
project/
├── src/
│   ├── app/
│   │   ├── layout.tsx
│   │   ├── page.tsx
│   │   ├── globals.css
│   │   ├── (auth)/
│   │   │   ├── login/
│   │   │   └── register/
│   │   └── dashboard/
│   ├── components/
│   │   ├── ui/
│   │   ├── layout/
│   │   └── features/
│   ├── lib/
│   │   ├── utils.ts
│   │   ├── api.ts
│   │   └── validations.ts
│   ├── hooks/
│   ├── stores/
│   └── types/
├── public/
├── tests/
├── package.json
├── next.config.js
├── tailwind.config.ts
├── tsconfig.json
├── components.json
├── .env.example
├── api_test.html
├── uiflow.md
└── README.md
` + "```" + `

### Rules
1. **App Router Only**: Use app/ directory, not pages/
2. **Server Components Default**: Mark client components with 'use client'
3. **shadcn/ui**: Use for base components, customize with theme
4. **Colocation**: Keep related files together
5. **Type Everything**: No any types allowed

### Component Pattern
` + "```tsx" + `
// src/components/features/feature-name.tsx
'use client' // Only if needs client-side JS

import { cn } from '@/lib/utils'

interface FeatureNameProps {
  // typed props
}

export function FeatureName({ ...props }: FeatureNameProps) {
  return (...)
}
` + "```" + `

### Color Implementation
Use CSS variables in globals.css + tailwind.config.ts:
` + "```css" + `
/* globals.css */
@layer base {
  :root {
    --primary: [H] [S]% [L]%;
    --background: [H] [S]% [L]%;
  }
}
` + "```" + `
`

const goAPIGuidelines = `## Template: go-api

### File Structure
` + "```" + `
project/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── handler/
│   │   ├── handler.go
│   │   ├── health.go
│   │   ├── user.go
│   │   └── middleware.go
│   ├── service/
│   │   ├── service.go
│   │   └── user.go
│   ├── repository/
│   │   ├── repository.go
│   │   ├── postgres.go
│   │   └── user.go
│   ├── model/
│   │   ├── user.go
│   │   └── errors.go
│   └── dto/
│       ├── request.go
│       └── response.go
├── pkg/
│   ├── validator/
│   └── httputil/
├── migrations/
├── scripts/
├── docs/
│   └── swagger.yaml
├── go.mod
├── go.sum
├── Makefile
├── Dockerfile
├── docker-compose.yml
├── .env.example
├── api_test.html
├── uiflow.md
└── README.md
` + "```" + `

### Rules
1. **Clean Architecture**: handler -> service -> repository
2. **Interface-Driven**: Define interfaces, implement concretely
3. **Error Handling**: Wrap errors with context, use custom error types
4. **No Global State**: Pass dependencies via constructors
5. **Context Propagation**: Pass context.Context through all layers

### Handler Pattern
` + "```go" + `
type UserHandler struct {
    userService service.UserService
    logger      *slog.Logger
}

func NewUserHandler(us service.UserService, l *slog.Logger) *UserHandler {
    return &UserHandler{userService: us, logger: l}
}

func (h *UserHandler) Routes() chi.Router {
    r := chi.NewRouter()
    r.Get("/", h.List)
    r.Post("/", h.Create)
    r.Get("/{id}", h.Get)
    return r
}
` + "```" + `

### Error Response Format
` + "```json" + `
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input",
    "details": [{"field": "email", "message": "must be valid email"}]
  }
}
` + "```" + `
`

const fullstackGuidelines = `## Template: fullstack

### File Structure
` + "```" + `
project/
├── frontend/
│   ├── src/
│   │   ├── app/
│   │   ├── components/
│   │   ├── lib/
│   │   │   └── api-client.ts
│   │   └── types/
│   ├── package.json
│   ├── next.config.js
│   ├── tailwind.config.ts
│   ├── Dockerfile
│   └── .env.local
├── backend/
│   ├── cmd/
│   │   └── api/
│   ├── internal/
│   │   ├── handler/
│   │   ├── service/
│   │   ├── repository/
│   │   └── model/
│   ├── migrations/
│   ├── go.mod
│   ├── Dockerfile
│   └── .env
├── shared/
│   └── api-types.ts
├── docker/
│   ├── nginx/
│   └── postgres/
├── scripts/
├── docker-compose.yml
├── docker-compose.prod.yml
├── Makefile
├── api_test.html
├── uiflow.md
├── swagger.yaml
└── README.md
` + "```" + `

### Rules
1. **API-First**: Define OpenAPI spec first, generate types for both ends
2. **Shared Types**: Keep frontend types in sync with backend via generation
3. **Docker Always**: All development happens in containers
4. **Environment Parity**: .env files match between dev/staging/prod
5. **Database Migrations**: Never modify production DB directly

### API Client Pattern
` + "```typescript" + `
// frontend/src/lib/api-client.ts
const API_URL = process.env.NEXT_PUBLIC_API_URL

export async function fetchAPI<T>(
  endpoint: string,
  options?: RequestInit
): Promise<T> {
  const res = await fetch(` + "`${API_URL}${endpoint}`" + `, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options?.headers,
    },
  })
  if (!res.ok) throw new APIError(await res.json())
  return res.json()
}
` + "```" + `

### CORS Configuration (Backend)
` + "```go" + `
func CORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", os.Getenv("FRONTEND_URL"))
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        next.ServeHTTP(w, r)
    })
}
` + "```" + `
`
