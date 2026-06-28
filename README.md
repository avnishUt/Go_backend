# Restaurant Inventory API

Backend service for a multi-tenant restaurant inventory, order, and payment management product.

## Stack

- Go 1.25+
- Gin
- PostgreSQL
- MongoDB
- Stateless REST APIs

## Current Part

Part 2 through Part 6 are complete:

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
- PostgreSQL connection package
- MongoDB connection package
- DB-aware health checks
- PostgreSQL migration folder
- Local Docker Compose for Postgres and MongoDB
- Makefile commands for common workflows
- Multi-tenant restaurant tables
- Restaurant branch tables
- Restaurant settings table
- Restaurant module using handler/service/repository structure
- Restaurant and branch APIs
- User registration and login
- JWT access tokens
- Role model and seeded roles
- Product categories and products
- Inventory items and stock adjustment

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
internal/modules         business modules
internal/platform/mongo  MongoDB client
internal/platform/postgres PostgreSQL client
internal/server          HTTP server lifecycle
migrations/postgres      PostgreSQL migrations
```

## Run

```bash
go mod tidy
go run ./cmd/api
```

For local databases:

```bash
docker compose up -d postgres mongo
```

If you started Postgres manually, this URL matches your container:

```bash
POSTGRES_URL='postgres://postgres:postgres@localhost:5432/restaurant?sslmode=disable'
```

Run PostgreSQL migrations:

```bash
make migrate-up
```

Then export the database URL before running the app.

Health checks:

```text
GET /
GET /health
GET /api/v1/health
```

Restaurant tenant APIs:

```text
POST  /api/v1/restaurants
GET   /api/v1/restaurants
GET   /api/v1/restaurants/:id
PATCH /api/v1/restaurants/:id
POST  /api/v1/restaurants/:id/branches
GET   /api/v1/restaurants/:id/branches
GET   /api/v1/restaurants/:id/settings
PATCH /api/v1/restaurants/:id/settings
```

Auth APIs:

```text
POST /api/v1/auth/register
POST /api/v1/auth/login
GET  /api/v1/auth/me
```

Product APIs:

```text
POST  /api/v1/product-categories
GET   /api/v1/product-categories?restaurant_id=:restaurant_id
POST  /api/v1/products
GET   /api/v1/products?restaurant_id=:restaurant_id
GET   /api/v1/products/:id
PATCH /api/v1/products/:id
```

Inventory APIs:

```text
POST  /api/v1/inventory/items
GET   /api/v1/inventory/items?restaurant_id=:restaurant_id
GET   /api/v1/inventory/items/:id
PATCH /api/v1/inventory/items/:id
POST  /api/v1/inventory/items/:id/adjust
```

Default port is `8080`. Override it with:

```bash
PORT=9090 go run ./cmd/api
```
