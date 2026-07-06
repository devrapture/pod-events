# Pod Events

Podcast notification platform — subscribe to Spotify shows and get notifications via Slack, Discord, Telegram, or WhatsApp.

## Architecture

```
┌──────────────┐     ┌──────────────┐     ┌────────────┐
│  Frontend    │────▶│  Backend      │────▶│  Postgres  │
│  Next.js 15  │     │  Gin + GORM  │     │    17      │
│  Better Auth │     │  Spotify API │     └────────────┘
│  Tailwind v4 │     │  Telegram    │
└──────────────┘     └──────┬───────┘
                            │
                     ┌──────▼───────┐
                     │  Notifications │
                     │  Slack/Discord │
                     │  Telegram/WA   │
                     └──────────────┘
```

**Auth:** Spotify OAuth (backend) + GitHub OAuth via Better Auth (frontend)

## Prerequisites

- [Go](https://go.dev/dl/) 1.25+
- [Bun](https://bun.sh/) — JavaScript runtime & package manager
- [Docker](https://docs.docker.com/get-docker/) — Postgres 17
- [Atlas CLI](https://atlasgo.io/getting-started) — database migrations
- [Air](https://github.com/air-verse/air) — hot-reload for Go (installed via `go install`)

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

The API runs at `http://localhost:8080` and the frontend at `http://localhost:3000`.

## Environment Configuration

### Backend (`apps/backend/.env`)

| Variable | Description |
|---|---|
| `DATABASE_URL` | Postgres connection string |
| `SPOTIFY_CLIENT_ID` | Spotify OAuth client ID |
| `SPOTIFY_CLIENT_SECRET` | Spotify OAuth client secret |
| `SPOTIFY_REDIRECT_URL` | Must match Spotify dashboard redirect URI |
| `JWT_SECRET` | Random secret (`openssl rand -base64 32`) |
| `TOKEN_ENCRYPTION_KEY` | AES-256 key (`make generate-encryption-key`) |
| `FRONTEND_URL` | Frontend URL for CORS and redirects |
| `TELEGRAM_BOT_TOKEN` | Telegram bot token from BotFather |
| `TELEGRAM_WEBHOOK_SECRET` | Random secret for webhook auth |
| `TELEGRAM_WEBHOOK_URL` | Public HTTPS URL for Telegram (use ngrok) |

### Frontend (`apps/frontend/.env`)

| Variable | Description |
|---|---|
| `NEXT_PUBLIC_API_URL` | Backend API URL (e.g. ngrok tunnel in dev) |
| `BETTER_AUTH_SECRET` | Better Auth secret (`openssl rand -base64 32`) |
| `BETTER_AUTH_GITHUB_CLIENT_ID` | GitHub OAuth app client ID |
| `BETTER_AUTH_GITHUB_CLIENT_SECRET` | GitHub OAuth app secret |

## Development

```bash
make dev           # Start Postgres, backend (Air + hot-reload), and frontend
make db-up         # Start Postgres only
make db-down       # Stop Postgres
make db-logs       # Tail Postgres logs
```

## API Documentation

This project uses **Swagger 2.0** (OpenAPI) generated from Go annotations.

### Generate docs

```bash
make swagger-docs
```

### View docs

Start the API server (`make dev`) and visit:

- **Swagger UI:** http://localhost:8080/swagger/index.html
- **Raw spec:** http://localhost:8080/swagger/doc.json

## Database Migrations

Migrations are managed with [Atlas](https://atlasgo.io/). Models are defined as GORM structs in `apps/backend/internal/models/`, and Atlas generates SQL migration files from them.

```bash
make migrate-diff NAME=describe_change   # Generate migration
make migrate-up                          # Apply pending migrations
make migrate-down                        # Rollback last migration
make migrate-status                      # Show migration state

# Production
DATABASE_URL=... DEV_DATABASE_URL=... make migrate-prod-up
```

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

### Spotify OAuth (Backend)

```
User → /auth/spotify/login → redirect to Spotify → authorize
→ callback with code+state → exchange for Spotify tokens
→ encrypt tokens, store in DB → create JWT → redirect to frontend
→ frontend calls POST /auth/exchange → receives JWT
```

### GitHub OAuth (Frontend via Better Auth)

The frontend uses **Better Auth** for GitHub authentication, independent of the backend Spotify auth.

## Project Structure

```
├── apps/
│   ├── backend/                  # Go API (Gin + GORM + Postgres)
│   │   ├── cmd/api/main.go       # Entry point
│   │   ├── docs/                 # Generated Swagger docs
│   │   ├── internal/
│   │   │   ├── config/           # Env-based configuration
│   │   │   ├── database/         # GORM connection setup
│   │   │   ├── dto/              # Request/response DTOs
│   │   │   ├── errors/           # Sentinel errors
│   │   │   ├── handler/          # HTTP handlers
│   │   │   ├── middleware/       # Auth + request logging
│   │   │   ├── migrations/       # Atlas GORM loader
│   │   │   ├── models/           # GORM model definitions
│   │   │   ├── notifications/    # Notifier interface + Telegram
│   │   │   ├── repositories/     # Data access layer
│   │   │   ├── routes/           # Router + CORS
│   │   │   ├── services/         # Business logic layer
│   │   │   └── spotify/          # Spotify API client
│   │   ├── migrations/           # Versioned SQL migrations
│   │   └── pkg/                  # Shared utilities
│   │       ├── crypto/           # AES-256-GCM encrypt/decrypt
│   │       ├── jwt/              # JWT generation/validation
│   │       ├── logger/           # Zap logger factory
│   │       └── response/         # API response wrapper
│   └── frontend/                 # Next.js 15 + Better Auth + Tailwind v4
│       └── src/
│           ├── app/              # App router pages
│           ├── components/       # UI components
│           ├── lib/              # Auth, API client, env
│           ├── server/           # Server-side logic
│           └── styles/           # Global styles
├── docker-compose.yml            # Postgres 17
└── Makefile                      # Dev + migration workflow
```
