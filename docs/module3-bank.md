# Module 3: Building the Go Bank Service

**Duration:** Day 2 (09:30–15:15)
**Location:** `internal/bank/`, `internal/challenges/bank/`

## Overview

Module 3 is a hands-on challenge: implement the Go Bank service from scratch. The codebase provides empty package stubs and a README quest guide in each package.

## Architecture

```
cmd/bank-api/     ← entry point (wire everything together)
internal/bank/
├── domain/       ← 1. Start here: Account and Transaction models
├── store/        ← 2. Postgres repository (go-jet + lib/pq)
├── service/      ← 3. Business logic
└── api/          ← 4. HTTP layer (Gin)
```

## Quests

See `internal/challenges/bank/README.md` for detailed acceptance criteria per quest.

## Prerequisites

- Postgres running: `make db-up`
- Migrations applied: `make migrate`

## Running Tests

```bash
make test-bank                 # all bank tests
go test ./internal/bank/...    # same
```
