# Architecture

The service uses a module-oriented structure:

- `handler.go`: HTTP request/response logic.
- `service.go`: business rules.
- `repository.go`: PostgreSQL queries and transactions.
- `dto.go`: API request payloads.
- `model.go`: response/domain models.
- `routes.go`: route registration.

PostgreSQL is the source of truth for transactional data. MongoDB is reserved for logs/events/analytics-style data when needed.

Every tenant-owned business table carries `restaurant_id`. Branch-scoped features also carry `branch_id`.
