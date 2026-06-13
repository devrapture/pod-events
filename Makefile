.PHONY: dev db-up db-down db-logs generate-encryption-key
.PHONY: test test-verbose test-coverage
.PHONY: migrate-diff migrate-up migrate-down migrate-status
.PHONY: migrate-prod-up migrate-prod-down migrate-prod-status

MIGRATE_DIR := apps/backend

dev: db-up
	cd apps/backend && air

db-up:
	docker compose up -d
	@echo "Waiting for Postgres to be healthy..."
	@while ! docker compose exec -T postgres pg_isready -U postgres -d podevents > /dev/null 2>&1; do sleep 1; done
	@docker compose exec -T postgres psql -U postgres -tc \
	  "SELECT 1 FROM pg_database WHERE datname='podevents_dev'" | grep -q 1 \
	  || docker compose exec -T postgres psql -U postgres -c "CREATE DATABASE podevents_dev;"

db-down:
	docker compose down

db-logs:
	docker compose logs -f

# ── Token Encryption ──────────────────────────────────────────

generate-encryption-key:                                              ## Generate a base64-encoded 32-byte AES-256 key for TOKEN_ENCRYPTION_KEY
	@echo "TOKEN_ENCRYPTION_KEY=$$(openssl rand -base64 32)"

# ── Tests ──────────────────────────────────────────────────

test:                                                                   ## Run all tests
	cd $(MIGRATE_DIR) && go test ./... -count=1

test-verbose:                                                           ## Run all tests with verbose output
	cd $(MIGRATE_DIR) && go test ./... -v -count=1

test-coverage:                                                          ## Run all tests with coverage report
	cd $(MIGRATE_DIR) && go test ./... -cover -count=1

# ── Database Migrations (Atlas) ──────────────────────────

migrate-diff:                                                           ## Generate a new migration from GORM models
	@if [ -z "$(NAME)" ]; then echo "Usage: make migrate-diff NAME=<description>"; exit 1; fi
	cd $(MIGRATE_DIR) && atlas migrate diff $(NAME) --env local

migrate-up:                                                             ## Apply pending migrations (local)
	cd $(MIGRATE_DIR) && atlas migrate apply --env local

migrate-down:                                                           ## Rollback the last migration (local)
	cd $(MIGRATE_DIR) && atlas migrate down --env local 1

migrate-status:                                                         ## Show migration status (local)
	cd $(MIGRATE_DIR) && atlas migrate status --env local

migrate-prod-up:                                                        ## Apply pending migrations (production)
	cd $(MIGRATE_DIR) && atlas migrate apply --env production

migrate-prod-down:                                                      ## Rollback the last migration (production)
	cd $(MIGRATE_DIR) && atlas migrate down --env production 1

migrate-prod-status:                                                    ## Show migration status (production)
	cd $(MIGRATE_DIR) && atlas migrate status --env production