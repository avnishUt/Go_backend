APP_PACKAGE=./cmd/api
POSTGRES_URL?=postgres://postgres:postgres@localhost:5432/restaurant?sslmode=disable
MIGRATIONS_PATH=migrations/postgres

.PHONY: tidy test build run seed worker db-up db-down migrate-up migrate-down docs-check docker-build

tidy:
	go mod tidy

test:
	go test ./...

build:
	go build ./...

run:
	go run $(APP_PACKAGE)

seed:
	go run ./cmd/seed -database "$(POSTGRES_URL)"

worker:
	go run ./cmd/worker -database "$(POSTGRES_URL)"

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
	test -f docs/deployment.md
	test -f docs/reports.md

docker-build:
	docker build -t restaurant-inventory-api:local .
