# Go Training

A practical Go training workshop for engineers transitioning from Python and Bash to Go. Uses a banking domain ("Go Bank") to teach production-ready API and CLI patterns.

## Structure

| Module | Location | Topic |
|---|---|---|
| 1 | `internal/fundamentals/` | API design, security, observability |
| 2 | `internal/basics/` | Go language building blocks |
| 3 | `internal/bank/` | Go Bank service (the challenge) |
| 4 | `internal/temporal/` | Temporal orchestration (demo) |

Challenges live in `internal/challenges/`.

## Quick Start

```bash
# 1. Setup
go mod tidy
make build
make test

# 2. Start the database
make db-up

# 3. Explore a topic
go test -v ./internal/basics/pointers/...
```

See [docs/setup.md](docs/setup.md) for full setup instructions.

## Workshop Timeline

- **Day 1:** Module 1 (API fundamentals) + Module 2 (Go basics)
- **Day 2:** Module 3 (Go Bank challenge) + Module 4 (Temporal demo)

## Tech Stack

- **Go 1.26.1** — generics, slog, standard library
- **Gin** — HTTP router
- **go-jet** — type-safe SQL
- **Cobra** — CLI framework
- **Testify** — test assertions
- **go.uber.org/mock** — interface mocking
- **PostgreSQL 15** — persistence
