<img width="1600" height="496" alt="image" src="https://github.com/user-attachments/assets/45864420-f9aa-492c-8780-e2e4d715dfee" />

# Pod Events

Podcast notification platform — sign in with Spotify, subscribe to shows, and get notified when new episodes drop via **Slack**, **Discord**, or **Telegram**.

## Features

- Spotify OAuth sign-in and JWT session for the SPA
- Search Spotify catalog and import saved/followed shows
- Subscribe to one or many shows; manage subscriptions
- Notification channels: Slack webhooks, Discord webhooks, Telegram (bot link)
- Dashboard summary (setup checklist, stats, recent activity)
- External cron-triggered episode polling with delivery logging
- Marketing landing page and privacy policy

## Architecture

```
┌──────────────┐     ┌──────────────┐     ┌────────────┐
│  Frontend    │────▶│  Backend     │────▶│  Postgres  │
│  Next.js 16  │     │  Gin + GORM  │     │    17      │
│  Tailwind v4 │     │  Spotify API │     └────────────┘
└──────────────┘     └──────┬───────┘
                            │
                     ┌──────▼───────────┐
                     │  Notifications   │
                     │  Slack / Discord │
                     │  Telegram        │
                     └──────────────────┘
```

**Auth:** Spotify OAuth (backend) → temporary exchange code → JWT for the frontend  
**Episode checks:** External scheduler calls `POST /cron/check-episodes` with `X-Cron-Secret`

## Prerequisites

- [Go](https://go.dev/dl/) 1.25+
- [Bun](https://bun.sh/) — JavaScript runtime & package manager
- [Docker](https://docs.docker.com/get-docker/) — Postgres 17
- [Atlas CLI](https://atlasgo.io/getting-started) — database migrations
- [Air](https://github.com/air-verse/air) — hot-reload for Go

```bash
brew install arigaio/tap/atlas
go install github.com/air-verse/air@latest
go install github.com/swaggo/swag/cmd/swag@latest
```

## Quick start

```bash
# 1. Clone and enter the repo
git clone <repo-url>
cd pod-events

# 2. Configure environment
cp apps/backend/.env.example apps/backend/.env
cp apps/frontend/.env.example apps/frontend/.env
# Edit both .env files with your secrets (see Environment section below)

# 3. Install frontend dependencies
make frontend-install

# 4. Start Postgres
make db-up

# 5. Run database migrations
make migrate-up

# 6. Start development servers (backend + frontend)
make dev
```

| Service | URL |
|---|---|
| API | http://localhost:8080 |
| Frontend | http://localhost:3000 |
| Swagger UI | http://localhost:8080/swagger/index.html |

## Environment Configuration

### Backend (`apps/backend/.env`)

| Variable | Description |
|---|---|
| `DATABASE_URL` | Postgres connection string (local example: host port **5433**, DB `podevents`) |
| `SPOTIFY_CLIENT_ID` | Spotify OAuth client ID |
| `SPOTIFY_CLIENT_SECRET` | Spotify OAuth client secret |
| `SPOTIFY_REDIRECT_URL` | Must match Spotify dashboard redirect URI (e.g. `http://localhost:8080/api/v1/auth/spotify/callback`) |
| `JWT_SECRET` | Random secret (`openssl rand -base64 32`) |
| `JWT_EXPIRES_IN_HOURS` | JWT lifetime in hours (default `24`) |
| `TOKEN_ENCRYPTION_KEY` | AES-256 key (`make generate-encryption-key`) |
| `FRONTEND_URL` | Frontend origin for CORS and post-login redirects |
| `TELEGRAM_BOT_TOKEN` | Telegram bot token from BotFather |
| `TELEGRAM_WEBHOOK_SECRET` | Random secret for webhook verification |
| `TELEGRAM_WEBHOOK_URL` | Public HTTPS webhook URL (e.g. ngrok in dev) |
| `BOT_NAME` | Telegram bot username/display name |
| `CRON_SECRET` | Secret for episode-check endpoint (`X-Cron-Secret` header) |
| `TRUSTED_PROXIES` | Optional comma-separated IPs/CIDRs allowed to set `X-Forwarded-For` for rate limiting; leave empty when the API is reached directly |
| `APP_ENV` | `development` or `production` (default `development`) |
| `PORT` | HTTP port (default `8080`) |
| `SENTRY_DSN` | Sentry DSN; leave empty to disable |
| `SENTRY_TRACES_SAMPLE_RATE` | Trace sample rate `0.0`–`1.0` (default `0.1`) |
| `SENTRY_ENABLE_LOGS` | Forward Info+ Zap logs to Sentry (default `true`) |
| `SENTRY_RELEASE` | Optional release id (e.g. Git SHA) |

### Frontend (`apps/frontend/.env`)

| Variable | Description |
|---|---|
| `NEXT_PUBLIC_API_URL` | Backend API base URL **including** `/api/v1` (e.g. `http://localhost:8080/api/v1` or an ngrok URL with that suffix) |

## Development

```bash
make dev           # Start Postgres, backend (Air), and frontend (bun dev)
make db-up         # Start Postgres only
make db-down       # Stop Postgres
make db-logs       # Tail Postgres logs
```

## API Documentation

Swagger 2.0 is generated from Go annotations.

```bash
make swagger-docs
```

With the API running:

- **Swagger UI:** http://localhost:8080/swagger/index.html
- **Raw spec:** http://localhost:8080/swagger/doc.json

## Database Migrations

Migrations are managed with [Atlas](https://atlasgo.io/). GORM models live in `apps/backend/internal/models/`; Atlas generates SQL under `apps/backend/migrations/` via the loader in `internal/migrations/`.

```bash
make migrate-diff NAME=describe_change   # Generate migration
make migrate-up                          # Apply pending migrations
make migrate-down                        # Rollback last migration
make migrate-status                      # Show migration state

# Production (set DATABASE_URL and DEV_DATABASE_URL)
make migrate-prod-up
```

Docker Compose exposes Postgres on **host port 5433**. `make db-up` starts the container and ensures the `podevents_dev` database exists (used as Atlas `dev` database).

## Available Make Targets

| Target | Description |
|---|---|
| `dev` | Start Postgres, backend (Air), and frontend |
| `db-up` | Start Postgres with Docker |
| `db-down` | Stop Postgres |
| `db-logs` | Tail Postgres logs |
| `test` | Run backend tests |
| `test-verbose` | Run tests with verbose output |
| `test-coverage` | Run tests with coverage report |
| `frontend-install` | Install frontend dependencies |
| `generate-encryption-key` | Generate a base64 AES-256 key |
| `swagger-docs` | Generate Swagger API docs |
| `migrate-diff NAME=<desc>` | Generate a migration |
| `migrate-up` | Apply migrations (local) |
| `migrate-down` | Rollback last migration |
| `migrate-status` | Show migration state |
| `migrate-prod-*` | Production migration variants |

## Authentication Flow

```
User → GET /api/v1/auth/spotify/login → Spotify authorize
→ GET /api/v1/auth/spotify/callback (code + state)
→ backend stores encrypted Spotify tokens, creates short-lived exchange code
→ redirect to frontend /auth/callback?code={exchangeCode}
→ POST ${NEXT_PUBLIC_API_URL}/auth/exchange
→ JWT + user returned to SPA; subsequent API calls use Authorization: Bearer
```

## Episode Monitoring & Notifications

New-episode detection is triggered by an external scheduler (not an in-process cron):

```text
POST /cron/check-episodes
Header: X-Cron-Secret: <CRON_SECRET>
```

The endpoint is guarded by a constant-time compare of `X-Cron-Secret` to `CRON_SECRET` — it is **not** user JWT auth.

When triggered, the backend:

1. Loads tracked shows and fetches each show’s latest available episode from Spotify (skips null placeholders in Spotify’s episode list; handles rate limits).
2. Persists new episodes and updates the show’s `latest_episode_id` watermark.
3. Notifies each subscriber on their active channels: **Slack**, **Discord**, or **Telegram**.

Delivery is recorded in `notification_logs`. Duplicate notifications for the same episode/channel are avoided; failed deliveries can be retried on a later run.

**IP rate limiting:** The API applies a global per-IP limiter. Configure `TRUSTED_PROXIES` only if a reverse proxy sits in front of the Go process and sets `X-Forwarded-For`. If browsers call the API URL directly (typical when the SPA is on Vercel and `NEXT_PUBLIC_API_URL` points at the API host), leave `TRUSTED_PROXIES` empty.

## Project Structure

```
├── apps/
│   ├── backend/                  # Go API (Gin + GORM + Postgres)
│   │   ├── cmd/api/main.go       # Entry point
│   │   ├── docs/                 # Generated Swagger docs
│   │   ├── internal/
│   │   │   ├── config/           # Env-based configuration
│   │   │   ├── cron/             # Episode check + notify orchestration
│   │   │   ├── database/         # GORM connection
│   │   │   ├── dto/              # Request/response DTOs
│   │   │   ├── errors/           # Sentinel errors
│   │   │   ├── handler/          # HTTP handlers
│   │   │   ├── metrics/          # Application metrics (Sentry)
│   │   │   ├── middleware/       # Auth, rate limit, cron secret, logging
│   │   │   ├── migrations/       # Atlas GORM loader
│   │   │   ├── models/           # GORM models
│   │   │   ├── notifications/    # Slack / Discord / Telegram notifiers
│   │   │   ├── repositories/     # Data access
│   │   │   ├── routes/           # Router + CORS + Swagger
│   │   │   ├── services/         # Business logic
│   │   │   └── spotify/          # Spotify API client
│   │   ├── migrations/           # Versioned SQL migrations
│   │   └── pkg/                  # crypto, jwt, logger, response, utils
│   └── frontend/                 # Next.js 16 + Tailwind v4
│       └── src/
│           ├── app/              # Marketing, dashboard, auth callback
│           ├── components/       # UI + feature + landing sections
│           ├── hooks/            # TanStack Query keys / queries / mutations
│           ├── lib/              # Auth, Axios, helpers
│           ├── services/         # API client + types
│           └── styles/           # Global styles
├── docker-compose.yml            # PostgreSQL 17 (port 5433)
└── Makefile                      # Dev + migration workflow
```

## Testing

```bash
# Backend
make test
# or
cd apps/backend && go test ./... -count=1

# Frontend
cd apps/frontend
bun run check
bun run typecheck
```

CI (GitHub Actions) runs frontend lint/typecheck/build and backend tests on pushes and PRs to **`staging`**.
