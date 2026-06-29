APP_PACKAGE=./cmd/api
POSTGRES_URL?=postgres://postgres:postgres@localhost:5432/restaurant?sslmode=disable
MIGRATIONS_PATH=migrations/postgres

.PHONY: tidy test build run db-up db-down migrate-up migrate-down docs-check

tidy:
	go mod tidy

test:
	go test ./...

build:
	go build ./...

run:
	go run $(APP_PACKAGE)

db-up:
	docker compose up -d postgres mongo

db-down:
	docker compose down

migrate-up:
	go run ./cmd/migrate -path $(MIGRATIONS_PATH) -database "$(POSTGRES_URL)"

migrate-down:
	@echo "down migrations are stored in $(MIGRATIONS_PATH), but cmd/migrate currently supports up-only safe migrations"

docs-check:
	test -f docs/openapi.yaml
	test -f docs/runbook.md
