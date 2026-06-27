# Restaurant Inventory API

Backend service for a multi-tenant restaurant inventory, order, and payment management product.

## Stack

- Go
- Gin
- PostgreSQL
- MongoDB
- Stateless REST APIs

## Current Part

Part 1 foundation is complete:

- Go module
- Gin HTTP server
- Environment config loader with HTTP timeouts
- Health check routes
- Standard JSON response helper
- Application error model
- Request ID middleware
- Request logger middleware
- Panic recovery middleware
- CORS middleware
- Not found and method-not-allowed handlers
- Graceful shutdown
- App bootstrap package
- `.env.example`

## Planned Parts

1. Backend foundation and architecture
2. PostgreSQL, MongoDB, and migrations setup
3. Multi-tenant restaurant structure
4. Users, roles, auth, and permissions
5. Product and category management
6. Inventory items and stock model
7. Inventory transactions
8. Orders and order items
9. Payment and mock UPI flow
10. Production hardening, docs, and tests

## Project Structure

```text
cmd/api                  application entry point
internal/app             app bootstrap
internal/config          environment config
internal/errors          application error model
internal/http/health     health handlers
internal/http/middleware HTTP middleware
internal/http/response   standard API responses
internal/http/router     route registration
internal/server          HTTP server lifecycle
```

## Run

```bash
go mod tidy
go run ./cmd/api
```

Health checks:

```text
GET /
GET /health
GET /api/v1/health
```

Default port is `8080`. Override it with:

```bash
PORT=9090 go run ./cmd/api
```
