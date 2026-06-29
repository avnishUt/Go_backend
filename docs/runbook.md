# Runbook

## Local Startup

```bash
docker run -d \
  --name restaurant-postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=restaurant \
  -p 5432:5432 \
  postgres:16

make migrate-up
POSTGRES_URL='postgres://postgres:postgres@localhost:5432/restaurant?sslmode=disable' go run ./cmd/api
```

## Health

```bash
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/health
```

Postgres should report `healthy` when `POSTGRES_URL` is configured.

## Checks Before Handoff

```bash
make tidy
make test
make build
make docs-check
```

## Production Notes

- Set a strong `JWT_SECRET`.
- Set `APP_ENV=production`.
- Use restricted `CORS_ALLOWED_ORIGINS`.
- Put the API behind HTTPS.
- Use managed PostgreSQL backups.
- Keep `RATE_LIMIT_ENABLED=true`.
