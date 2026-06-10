.PHONY: dev db-up db-down db-logs

dev: db-up
	cd apps/backend && air

db-up:
	docker compose up -d

db-down:
	docker compose down

db-logs:
	docker compose logs -f
