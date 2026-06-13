# Pod Events

## Prerequisites

- [Go](https://go.dev/dl/) 1.25+
- [Docker](https://docs.docker.com/get-docker/)
- [Atlas CLI](https://atlasgo.io/getting-started)

```bash
brew install arigaio/tap/atlas
```

## Quick start

```bash
make db-up              # Start Postgres
make dev                # Start the API with hot-reload (Air)
```

## Database migrations

This project uses [Atlas](https://atlasgo.io/) to manage database schema
migrations. Models are defined as GORM structs in
`apps/backend/internal/models/`, and Atlas generates SQL migration files
from them.

### Migration workflow

1. Edit or add GORM model structs in `apps/backend/internal/models/`.

2. Generate a migration:

   ```bash
   make migrate-diff NAME=describe_your_change
   ```

   This compares your GORM models against the local database and writes a
   versioned SQL file to `apps/backend/migrations/`.

3. Apply the migration locally:

   ```bash
   make migrate-up
   ```

4. Commit the generated SQL file alongside your code changes.

### Managing migrations

| Command | Description |
|---|---|
| `make migrate-diff NAME=<desc>` | Generate a new migration from GORM model changes |
| `make migrate-up` | Apply all pending migrations (local) |
| `make migrate-down` | Rollback the last migration (local) |
| `make migrate-status` | Show applied / pending migration state (local) |

### Production

For production, set these environment variables before running migration
commands:

- `DATABASE_URL` — target database connection string
- `DEV_DATABASE_URL` — a separate dev/staging database (Atlas uses this to
  plan migrations safely)

```bash
DATABASE_URL=postgres://user:pass@host:5432/db?sslmode=require \
DEV_DATABASE_URL=postgres://user:pass@dev-host:5432/db?sslmode=require \
make migrate-prod-status

DATABASE_URL=... DEV_DATABASE_URL=... make migrate-prod-up
DATABASE_URL=... DEV_DATABASE_URL=... make migrate-prod-down
```

| Command | Description |
|---|---|
| `make migrate-prod-up` | Apply pending migrations to production |
| `make migrate-prod-down` | Rollback the last migration on production |
| `make migrate-prod-status` | Show migration state on production |

### First-time setup

If the database is empty and you need to create the initial schema:

```bash
make db-up
make migrate-diff NAME=initial_schema
make migrate-up
```

### How it works

```
GORM models (internal/models/)
    |
    |   make migrate-diff (Atlas runs the GORM loader program)
    v
SQL migration file (apps/backend/migrations/)
    |
    |   make migrate-up (Atlas applies the file to the database)
    v
Database
```

- `apps/backend/atlas.hcl` defines two environments (`local` and
  `production`) and uses `data "external_schema"` with a `program` that
  runs `go run -tags=atlas ./internal/migrations/loader/` to load your
  GORM models.
- The loader is build-tagged with `//go:build atlas` so it only compiles
  during migration runs, not in the main binary.
- The `local` environment uses `podevents_dev` as its dev database (a
  clean database used by Atlas to compute schema diffs and plan
  rollbacks). This database is created automatically by `make db-up`.

## Project structure

```
├── apps/
│   ├── backend/            # Go API (Gin + GORM + Postgres)
│   │   ├── atlas.hcl       # Atlas migration config
│   │   ├── cmd/api/        # Entry point
│   │   ├── internal/
│   │   │   ├── config/     # Env-based configuration
│   │   │   ├── database/   # GORM connection setup
│   │   │   ├── handler/    # HTTP handlers
│   │   │   ├── middleware/ # Request logging
│   │   │   ├── migrations/ # Atlas GORM loader
│   │   │   ├── models/     # GORM model definitions
│   │   │   └── routes/     # Gin router
│   │   ├── migrations/     # Versioned SQL migration files
│   │   └── pkg/            # Shared utilities
│   └── frontend/           # (placeholder)
├── docker-compose.yml      # Postgres 17
└── Makefile                # Dev + migration workflow
```
