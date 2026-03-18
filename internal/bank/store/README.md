# Store

The store package implements the repository pattern for Go Bank using [go-jet](https://github.com/go-jet/jet).

## Responsibilities

- Type-safe SQL queries (SELECT, INSERT, UPDATE)
- Repository interfaces consumed by the service layer
- Database connection management

## Dependencies

- PostgreSQL 15 (via docker-compose)
- `github.com/go-jet/jet/v2` for query building
- `github.com/lib/pq` as the Postgres driver

## Challenge

Implement the store after running the migrations in `migration/`.

See `internal/challenges/bank/README.md` for quest details.
