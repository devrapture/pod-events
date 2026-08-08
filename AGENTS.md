# AGENTS.md

Guide for AI coding agents working on this repository.

## Project Overview

**PodEvents** is a podcast notification platform. Users sign in with Spotify, subscribe to shows, and receive notifications when new episodes drop via **Slack**, **Discord**, or **Telegram**.

### Tech Stack

| Layer | Technology |
|---|---|
| Frontend | Next.js 16 (App Router), React 19, TypeScript ~5.9, Tailwind CSS v4, Bun |
| Backend | Go 1.25, Gin (HTTP), GORM (ORM), PostgreSQL 17 |
| Auth | Spotify OAuth (backend), JWT (backend), React auth context + localStorage token (frontend) |
| Database | PostgreSQL 17 (Docker Compose), Atlas (SQL migrations generated from GORM models) |
| Observability | Sentry (backend: gin + zap; frontend: `@sentry/nextjs`) |
| Tooling | Biome (frontend lint/format), Husky + lint-staged, Air (Go hot-reload), Docker Compose, Swag |
| Validation | Zod via `@t3-oss/env-nextjs` (frontend env), go-playground/validator (backend DTOs) |

### Architecture

```
Frontend (Next.js 16) --> Backend (Gin + GORM) --> PostgreSQL 17
                                |
                         Notifications (Slack / Discord / Telegram)
                                ^
                         External cron --> POST /cron/check-episodes
```

Auth flow: User logs in via Spotify OAuth on the backend, receives a temporary exchange code, then the frontend exchanges it for a JWT. The JWT is stored (localStorage) and sent as `Authorization: Bearer` on API requests.

Episode detection is **not** an in-process scheduler. An external caller hits `POST /cron/check-episodes` with header `X-Cron-Secret: <CRON_SECRET>`.

## Repository Structure

```
pod-events/
├── apps/
│   ├── backend/                # Go API server
│   │   ├── cmd/api/main.go     # Entry point (graceful shutdown, Sentry, wiring)
│   │   ├── internal/
│   │   │   ├── config/         # Env-based config (godotenv)
│   │   │   ├── cron/           # Episode check + notification orchestration
│   │   │   ├── database/       # GORM PostgreSQL connection
│   │   │   ├── dto/            # Request/response DTOs
│   │   │   ├── errors/         # Sentinel errors (package apperrors)
│   │   │   ├── handler/        # HTTP handlers (Gin)
│   │   │   ├── metrics/        # Sentry metrics helpers
│   │   │   ├── middleware/     # JWT auth, IP rate limit, cron secret, logging, metrics
│   │   │   ├── migrations/     # Atlas GORM schema loader (//go:build atlas)
│   │   │   ├── models/         # GORM models (embed Base: UUID, timestamps, soft delete)
│   │   │   ├── notifications/  # Notifier interface + slack/, discord/, telegram/
│   │   │   ├── repositories/   # Data access layer
│   │   │   ├── routes/         # Router, CORS, Swagger, middleware stack
│   │   │   ├── services/       # Business logic
│   │   │   └── spotify/        # Spotify Web API client
│   │   ├── migrations/         # Versioned Atlas SQL migrations
│   │   ├── docs/               # Generated Swagger docs
│   │   ├── atlas.hcl           # Atlas local/production envs
│   │   └── pkg/                # Shared utilities (no imports from internal/)
│   │       ├── crypto/         # AES-256-GCM encrypt/decrypt + fingerprint
│   │       ├── jwt/            # JWT generation/validation
│   │       ├── logger/         # Zap logger factory
│   │       ├── response/       # Standardized API response wrapper
│   │       └── utils/          # Shared helpers (e.g. truncate)
│   └── frontend/               # Next.js App Router app
│       ├── src/
│       │   ├── app/            # Routes: (marketing), dashboard/*, auth/callback
│       │   ├── components/     # ui/, channels/, overview/, search/, import/, subscriptions/, sections/
│       │   ├── hooks/          # keys/, queries/, mutations/ (TanStack Query + react-query-kit)
│       │   ├── lib/            # Auth, Axios, API helpers, constants
│       │   ├── services/       # api-services.ts + types.ts
│       │   ├── styles/         # globals.css (Tailwind v4)
│       │   ├── middleware.ts   # Dashboard auth cookie redirect
│       │   ├── env.js          # Zod env validation
│       │   └── instrumentation*.ts  # Sentry registration
│       ├── sentry.*.config.ts  # Sentry client/server/edge init
│       └── biome.jsonc
├── .github/workflows/          # ci.yml (staging), security.yml
├── docker-compose.yml          # Postgres 17 on host port 5433
├── Makefile
├── AGENTS.md
└── README.md
```

## Development Workflow

### Prerequisites

- Go 1.25+
- Bun (only JS package manager — **not** npm/pnpm/Yarn)
- Docker (PostgreSQL 17)
- Atlas CLI (`brew install arigaio/tap/atlas`)
- Air (`go install github.com/air-verse/air@latest`)
- Swag (`go install github.com/swaggo/swag/cmd/swag@latest`) — for regenerating Swagger

### Setup

```bash
cp apps/backend/.env.example apps/backend/.env
cp apps/frontend/.env.example apps/frontend/.env
# Edit .env files with real secrets
make frontend-install        # bun install in apps/frontend
make db-up                   # Start Postgres, wait for health, ensure podevents_dev exists
make migrate-up              # Apply Atlas migrations to local DB
```

Local Atlas (`atlas.hcl` env `local`):

- App DB URL: `postgres://postgres:password@127.0.0.1:5433/podevents?sslmode=disable`
- Dev DB (Atlas shadow): `.../podevents_dev` (created by `make db-up` if missing)

### Run

```bash
make dev                     # db-up + backend (Air) + frontend (bun dev)
# API:      http://localhost:8080
# Frontend: http://localhost:3000
# Swagger:  http://localhost:8080/swagger/index.html
```

### Database Migrations

```bash
make migrate-diff NAME=describe_change
make migrate-up
make migrate-down
make migrate-status
# Production (requires DATABASE_URL + DEV_DATABASE_URL in env):
make migrate-prod-up
```

## Coding Standards

### Backend (Go)

- **Layers**: `handler → service → repository → model`; DTOs for HTTP request/response shapes
- **Files**: `snake_case.go`; packages: lowercase single-word (`handler`, `services`, `repositories`, `models`, `dto`, `apperrors`)
- **Handlers**: Struct + `NewXHandler(...)`; bind JSON with Gin/validator tags; map sentinel errors to HTTP status
- **Services**: Interface + constructor injection
- **Repositories**: Constructor injection; use context-aware DB access patterns already in package
- **Models**: Embed `Base` (UUID PK, timestamps, `gorm.DeletedAt`); `BeforeCreate` sets UUID when nil
- **Responses**: `response.SuccessResponse` / `response.ErrorResponse` / `response.ValidationError` from `pkg/response`
- **Errors**: Sentinel vars in `internal/errors` (`apperrors` import alias)
- **Auth**: JWT Bearer; `c.Set("userID", ...)` in `AuthMiddleware`
- **Config**: `config.Load()`; `mustGetEnv` required, `getEnv` optional
- **Encryption**: AES-256-GCM (`pkg/crypto`) for Spotify tokens and webhook destinations
- **API prefix**: App routes under `/api/v1`. Cron is **`/cron/*`** (not under `/api/v1`)
- **Rate limiting**: Global IP rate limiter in `routes.Setup` (5 req/s, burst 10). Trust `X-Forwarded-For` only when peer is in `TRUSTED_PROXIES`
- **Logging**: Zap, constructor-injected; Sentry zap core when DSN configured
- **Swagger**: Annotation comments on handlers; regenerate with `make swagger-docs`
- **Subscribe flow**: Upsert shows → create subscriptions → **best-effort** seed of `latest_episode_id` (must not fail subscribe). Spotify episode lists may contain leading `null` items — client skips nulls (`limit=5`)

### Frontend (TypeScript/React)

- **Components**: Functional, named exports; `"use client"` when needed
- **UI**: shadcn-style primitives under `src/components/ui/` (`cva`, `cn()`)
- **Styling**: Tailwind CSS v4 only
- **Server state**: TanStack Query via `react-query-kit` helpers in `hooks/queries` and `hooks/mutations`
- **API**: `src/services/api-services.ts` (Axios). Plain `server` vs `serverWithInterceptors` (auth + ngrok headers)
- **Types**: `src/services/types.ts` aligned with backend DTOs
- **Query keys**: Factories in `hooks/keys/`
- **Paths**: `@/*` → `./src/*`
- **Files**: `kebab-case` for components/utils
- **Env**: `src/env.js` with Zod; required client var `NEXT_PUBLIC_API_URL` (must include `/api/v1`)
- **Notification UI channels**: Slack, Discord, Telegram only (`add-channel-form.tsx`)

### Formatting

- Frontend: Biome (tabs, organize imports, sort Tailwind classes)
- Backend: gofmt / goimports

## Environment Variables

### Backend (`apps/backend/.env`) — from `.env.example`

| Variable | Required | Notes |
|---|---|---|
| `DATABASE_URL` | yes | Local example uses port **5433** and DB `podevents` |
| `SPOTIFY_CLIENT_ID` / `SPOTIFY_CLIENT_SECRET` / `SPOTIFY_REDIRECT_URL` | yes | Redirect must match Spotify dashboard |
| `TOKEN_ENCRYPTION_KEY` | yes | Base64 32-byte AES key (`make generate-encryption-key`) |
| `FRONTEND_URL` | yes | CORS + OAuth redirect target |
| `JWT_SECRET` | yes | |
| `JWT_EXPIRES_IN_HOURS` | no | Default `24` |
| `TELEGRAM_BOT_TOKEN` / `TELEGRAM_WEBHOOK_SECRET` / `TELEGRAM_WEBHOOK_URL` / `BOT_NAME` | for Telegram | Webhook public HTTPS URL |
| `CRON_SECRET` | yes (for cron) | Compared via `X-Cron-Secret` |
| `APP_ENV` / `PORT` | no | Defaults `development` / `8080` |
| `TRUSTED_PROXIES` | no | Comma-separated IPs/CIDRs; empty = never trust XFF |
| `SENTRY_DSN` / `SENTRY_TRACES_SAMPLE_RATE` / `SENTRY_ENABLE_LOGS` / `SENTRY_RELEASE` | no | Empty DSN disables Sentry |

### Frontend (`apps/frontend/.env`)

| Variable | Required | Notes |
|---|---|---|
| `NEXT_PUBLIC_API_URL` | yes | Full API base including `/api/v1` |

`SKIP_ENV_VALIDATION` can skip Zod env checks (builds).

## Main API Surface (summary)

Protected routes use `AuthMiddleware` unless noted.

| Area | Methods / paths |
|---|---|
| Health | `GET /api/v1/health` (public) |
| Auth | `GET /api/v1/auth/spotify/login`, `.../callback`, `POST /api/v1/auth/exchange`; `GET /api/v1/auth/me` (auth) |
| Dashboard | `GET /api/v1/dashboard/summary` |
| Shows | `GET /api/v1/shows/saved`, `GET /api/v1/shows/search`, `POST /api/v1/shows/subscribe` |
| Subscriptions | `GET /api/v1/subscriptions`, `DELETE /api/v1/subscriptions/:id` |
| Channels | `GET/POST /api/v1/channels`, `POST .../toggle`, `DELETE .../:channelID` |
| Telegram | `POST /api/v1/telegram/generate-link`; webhook `POST /api/v1/webhooks/telegram` |
| Cron | `POST /cron/check-episodes` + `X-Cron-Secret` |

## Git Workflow

### Branches

- **`staging`** — integration branch; CI runs on push/PR to `staging`
- Feature branches: `feat/…`, `fix/…`, `chore/…`, `ci/…`, `docs/…`

### Commits

Conventional Commits: `type(scope): message`  
Types: `feat`, `fix`, `refactor`, `chore`, `ci`, `docs`  
Scopes often: `(backend)`, `(frontend)`, `(deps)`

### Hooks (Husky)

- **pre-commit**: lint-staged → `biome check --write` on staged frontend files
- **pre-push**: `bun run build` (frontend)

## Guidelines for AI Agents

1. Read surrounding code and match existing patterns before editing.
2. Reuse UI primitives, query/mutation hooks, and Go layer patterns — do not invent parallel abstractions.
3. Keep diffs focused; avoid unrelated refactors.
4. Prefer Bun for all frontend package operations.
5. Do not commit secrets or `.env` files.
6. Update README/Swagger when user-facing or API behavior changes.
7. Run validation before considering work done.
8. WhatsApp is **not** a supported delivery channel in the UI or notifier layer; do not document or implement it as available unless implementing the full path.

## Validation Checklist

### Frontend

```bash
cd apps/frontend
bun run check
bun run typecheck
bun run build
```

### Backend

```bash
cd apps/backend
go test ./... -count=1
go vet ./...
```

## Common Commands

### Makefile (repo root)

| Command | Description |
|---|---|
| `make dev` | Postgres + Air backend + `bun dev` frontend |
| `make db-up` / `db-down` / `db-logs` | Postgres lifecycle |
| `make frontend-install` | `bun install` |
| `make swagger-docs` | Regenerate Swagger |
| `make test` / `test-verbose` / `test-coverage` | Backend tests |
| `make migrate-diff NAME=<desc>` | Generate migration |
| `make migrate-up` / `migrate-down` / `migrate-status` | Local migrations |
| `make migrate-prod-*` | Production migrations |
| `make generate-encryption-key` | Print `TOKEN_ENCRYPTION_KEY=...` |

### Frontend (`apps/frontend/`)

| Command | Description |
|---|---|
| `bun run dev` | Dev server (Turbopack) |
| `bun run build` / `start` | Production build / serve |
| `bun run check` / `check:write` / `check:unsafe` | Biome |
| `bun run typecheck` | `tsc --noEmit` |

### Backend (`apps/backend/`)

| Command | Description |
|---|---|
| `go test ./... -count=1` | All tests |
| `go test ./... -v -count=1` | Verbose |
| `go test ./... -cover -count=1` | Coverage |
| `go vet ./...` | Static analysis |
