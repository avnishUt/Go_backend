# Deployment

## Build Image

```bash
docker build -t restaurant-inventory-api:local .
```

## Run Migrations

```bash
docker run --rm \
  -e POSTGRES_URL='postgres://postgres:postgres@host.docker.internal:5432/restaurant?sslmode=disable' \
  restaurant-inventory-api:local \
  /app/migrate -path /app/migrations/postgres
```

## Seed Admin

```bash
docker run --rm \
  -e POSTGRES_URL='postgres://postgres:postgres@host.docker.internal:5432/restaurant?sslmode=disable' \
  restaurant-inventory-api:local \
  /app/seed -admin-email admin@example.com -admin-password password123
```

## Run API

```bash
docker run --rm -p 8080:8080 \
  -e APP_ENV=production \
  -e JWT_SECRET='replace-with-strong-secret' \
  -e POSTGRES_URL='postgres://postgres:postgres@host.docker.internal:5432/restaurant?sslmode=disable' \
  restaurant-inventory-api:local
```
