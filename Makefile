APP_PACKAGE=./cmd/api
POSTGRES_URL?=postgres://postgres:postgres@localhost:5432/restaurant?sslmode=disable
MIGRATIONS_PATH=migrations/postgres

.PHONY: tidy test run db-up db-down migrate-up migrate-down

tidy:
	go mod tidy

test:
	go test ./...

run:
	go run $(APP_PACKAGE)

db-up:
	docker compose up -d postgres mongo

db-down:
	docker compose down

migrate-up:
	go run ./cmd/migrate -path $(MIGRATIONS_PATH) -database "$(POSTGRES_URL)"

migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(POSTGRES_URL)" down 1
