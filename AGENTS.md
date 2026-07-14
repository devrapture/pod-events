# AGENTS.md

Guide for AI coding agents working on this repository.

## Project Overview

**PodEvents** is a podcast notification platform. Users subscribe to Spotify shows and receive instant notifications when new episodes drop, delivered via Slack, Discord, Telegram, or WhatsApp.

### Tech Stack

| Layer | Technology |
|---|---|
| Frontend | Next.js 15 (App Router), React 19, TypeScript 5.8, Tailwind CSS v4, Bun |
| Backend | Go 1.25, Gin (HTTP), GORM (ORM), PostgreSQL 17 |
| Auth | Better Auth (GitHub OAuth, frontend), JWT (backend), Spotify OAuth (backend) |
| Database | PostgreSQL 17 (Docker), Atlas (schema migrations from GORM models) |
| Tooling | Biome (lint/format), Husky + lint-staged, Air (Go hot-reload), Docker Compose |
| Validation | Zod (frontend env), go-playground/validator (backend) |

### Architecture

```
Frontend (Next.js 15) --> Backend (Gin + GORM) --> PostgreSQL 17
                                |
                         Notifications (Slack/Discord/Telegram/WhatsApp)
```

Dual OAuth: Spotify OAuth on the backend (podcast data), GitHub OAuth on the frontend via Better Auth (user identity).

## Repository Structure

```
pod-events/
├── apps/
│   ├── backend/                # Go API server
│   │   ├── cmd/api/main.go     # Entry point
│   │   ├── internal/           # Private application code
│   │   │   ├── config/         # Env-based config (godotenv)
│   │   │   ├── database/       # GORM PostgreSQL connection
│   │   │   ├── dto/            # Request/response DTOs per domain
│   │   │   ├── errors/         # Sentinel errors (apperrors package)
│   │   │   ├── handler/        # HTTP handlers (Gin)
│   │   │   ├── middleware/     # Auth (JWT), logging
│   │   │   ├── models/         # GORM models (embed Base: UUID, timestamps, soft delete)
│   │   │   ├── notifications/  # Notifier interface + implementations
│   │   │   ├── repositories/   # Data access layer
│   │   │   ├── routes/         # Router setup, CORS, Swagger
│   │   │   ├── services/       # Business logic layer
│   │   │   └── spotify/        # Spotify Web API client
│   │   ├── migrations/         # Atlas SQL migration files
│   │   ├── pkg/                # Shared utilities (no internal imports)
│   │   │   ├── crypto/         # AES-256-GCM encrypt/decrypt
│   │   │   ├── jwt/            # JWT generation/validation
│   │   │   ├── logger/         # Zap logger factory
│   │   │   └── response/       # Standardized API response wrapper
│   │   └── docs/               # Generated Swagger documentation
│   └── frontend/               # Next.js 15 SPA
│       ├── src/
│       │   ├── app/            # Next.js App Router pages
│       │   │   ├── (marketing)/    # Landing page (route group)
│       │   │   ├── dashboard/      # Dashboard pages (import, search, subscriptions)
│       │   │   └── api/auth/       # Better Auth API routes
│       │   ├── components/
│       │   │   ├── ui/         # Shared UI primitives (shadcn/ui pattern: button, card, badge, toast)
│       │   │   ├── {feature}/  # Feature-specific components (overview, search, import, subscriptions)
│       │   │   └── sections/   # Landing page sections
│       │   ├── hooks/
│       │   │   ├── keys/       # TanStack Query key factories
│       │   │   ├── queries/    # TanStack Query query hooks
│       │   │   └── mutations/  # TanStack Query mutation hooks
│       │   ├── lib/            # Utilities, Axios setup, constants, auth helpers
│       │   ├── server/         # Server-side code (Better Auth config)
│       │   ├── services/       # Centralized API client + TypeScript types
│       │   └── styles/         # globals.css (Tailwind v4)
│       └── biome.jsonc         # Biome linter/formatter config
├── .github/
│   ├── workflows/              # CI (ci.yml) and security (security.yml) workflows
│   └── dependabot.yml          # Weekly dependency updates
├── docker-compose.yml          # PostgreSQL 17 container
├── Makefile                    # Dev workflow orchestration
└── README.md                   # Project documentation
```

## Development Workflow

### Prerequisites

- Go 1.25+
- Bun (JavaScript runtime and package manager — **not** npm, pnpm, or Yarn)
- Docker (PostgreSQL 17)
- Atlas CLI (`brew install arigaio/tap/atlas`)
- Air (`go install github.com/air-verse/air@latest`)
- Swag (`go install github.com/swaggo/swag/cmd/swag@latest`)

### Setup

```bash
cp apps/backend/.env.example apps/backend/.env
cp apps/frontend/.env.example apps/frontend/.env
# Edit .env files with your secrets
make frontend-install        # bun install in apps/frontend
make db-up                   # Start Postgres, wait for health, create DB
```

### Run

```bash
make dev                     # Postgres + backend (Air) + frontend (bun dev)
# API: http://localhost:8080
# Frontend: http://localhost:3000
# Swagger: http://localhost:8080/swagger/index.html
```

### Database Migrations

```bash
make migrate-diff NAME=describe_change   # Generate migration from GORM models
make migrate-up                          # Apply pending migrations (local)
make migrate-down                        # Rollback last migration
make migrate-status                      # Show migration state
```

## Coding Standards

### Backend (Go)

- **Architecture**: Layered: `handler → service → repository → model`, with DTOs for request/response
- **File naming**: `snake_case.go` (e.g., `auth_handler.go`, `user_repository.go`)
- **Package naming**: Lowercase, single-word (`handler`, `services`, `repositories`, `models`, `dto`, `apperrors`)
- **Handlers**: Struct-based with constructor functions (`NewAuthHandler(...)`)
- **Services**: Interface-based, constructor-injected (`NewAuthService(...)`)
- **Repositories**: Constructor-injected (`NewUserRepository(db)`)
- **Models**: Embed `Base` struct (UUID PK, `CreatedAt`, `UpdatedAt`, soft delete via `gorm.DeletedAt`). `BeforeCreate` hook generates UUID
- **Responses**: Use `response.SuccessResponse()` and `response.ErrorResponse()` from `pkg/response/`
- **Errors**: Sentinel errors in `internal/errors/errors.go` (`apperrors` package)
- **Auth**: JWT Bearer tokens, claims set in context via `c.Set("userID", ...)`
- **Config**: Struct-based, loaded from env with `godotenv`. `mustGetEnv()` for required, `getEnv()` for optional
- **Encryption**: AES-256-GCM for Spotify token storage (base64-encoded keys)
- **API versioning**: All routes prefixed with `/api/v1`
- **Logging**: Zap structured logging, passed as dependency through constructors
- **Swagger**: Go doc comments with `@Summary`, `@Tags`, `@Security` annotations

### Frontend (TypeScript/React)

- **Components**: Functional, named exports (`export function Button(...)`), `"use client"` directive when needed
- **UI primitives**: shadcn/ui pattern with `class-variance-authority` (cva) in `src/components/ui/`
- **Styling**: Tailwind CSS v4 exclusively. Use `cn()` (clsx + tailwind-merge) for class composition
- **State management**: TanStack React Query for server state
- **Data fetching**: Centralized API client in `src/services/api-services.ts` (Axios). Two instances: plain and auth-intercepted
- **Types**: All API types in `src/services/types.ts`, matching backend DTOs
- **Query keys**: Hierarchical key factory pattern (`showKeys.all`, `showKeys.saved(params)`)
- **Path aliases**: `@/*` maps to `./src/*`
- **File naming**: `kebab-case.tsx` for components, `kebab-case.ts` for utilities
- **Environment validation**: `@t3-oss/env-nextjs` with Zod schemas in `src/env.js`

### Formatting

- **Biome** for frontend: tabs for indentation, organized imports, sorted Tailwind classes
- **Go standard formatting** (gofmt/goimports via editor) for backend

## Git Workflow

### Branches

- **`staging`** — Integration/development branch (CI runs on push/PR to this)
- **Feature branches**: Always create a branch before making changes
  - Naming: `feat/short-description`, `fix/short-description`, `chore/short-description`, `ci/short-description`, `docs/short-description`

### Commits

Follow **Conventional Commits**: `type(scope): message`

- Types: `feat`, `fix`, `refactor`, `chore`, `ci`, `docs`
- Optional scopes: `(backend)`, `(frontend)`, `(deps)`

### Hooks (Husky)

- **pre-commit**: Runs `biome check --write` on staged files via lint-staged
- **pre-push**: Runs `bun run build`

## Guidelines for AI Agents

1. **Read surrounding code** before making edits. Understand the existing patterns in the file and adjacent files.
2. **Reuse existing patterns** instead of introducing new ones. Check `src/components/ui/` for UI patterns, `hooks/` for query patterns, `internal/` for Go layer patterns.
3. **Avoid unnecessary refactors**. Do not restructure code unrelated to the task.
4. **Keep changes focused** on the requested task. Do not add unrelated improvements.
5. **Preserve backward compatibility** unless explicitly instructed otherwise.
6. **Update documentation** when behavior changes (README, Swagger annotations, etc.).
7. **Run validation** before considering work complete.
8. **Use Bun** as the only JavaScript package manager. Never use npm, pnpm, or Yarn.

## Validation Checklist

Run these before committing:

### Frontend

```bash
cd apps/frontend
bun run check          # Lint (Biome)
bun run typecheck      # TypeScript type checking
bun run build          # Production build
```

### Backend

```bash
cd apps/backend
go test ./... -count=1  # Tests
go vet ./...            # Static analysis
```

## Common Commands

### Makefile (run from repo root)

| Command | Description |
|---|---|
| `make dev` | Start everything (Postgres, backend, frontend) |
| `make db-up` | Start Postgres, wait for health, create `podevents_dev` |
| `make db-down` | Stop Postgres |
| `make frontend-install` | Install frontend dependencies with Bun |
| `make swagger-docs` | Generate Swagger docs from Go annotations |
| `make test` | Run all backend tests |
| `make test-verbose` | Run tests with verbose output |
| `make test-coverage` | Run tests with coverage report |
| `make migrate-diff NAME=<desc>` | Generate a new migration |
| `make migrate-up` | Apply pending migrations (local) |
| `make migrate-down` | Rollback last migration (local) |
| `make migrate-status` | Show migration state (local) |
| `make generate-encryption-key` | Generate AES-256 key for `TOKEN_ENCRYPTION_KEY` |

### Frontend Scripts (run from `apps/frontend/`)

| Command | Description |
|---|---|
| `bun run dev` | Start dev server (Turbopack) |
| `bun run build` | Production build |
| `bun run check` | Lint with Biome |
| `bun run check:write` | Lint + auto-fix |
| `bun run check:unsafe` | Lint + auto-fix (unsafe) |
| `bun run typecheck` | TypeScript type checking |
| `bun run start` | Start production server |

### Backend (run from `apps/backend/`)

| Command | Description |
|---|---|
| `go test ./... -count=1` | Run all tests |
| `go test ./... -v -count=1` | Verbose test output |
| `go test ./... -cover -count=1` | Tests with coverage |
| `go vet ./...` | Static analysis |
