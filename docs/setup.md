# Setup Guide

## Prerequisites

- [Go 1.26.1+](https://go.dev/dl/)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/) (for Postgres)
- [golangci-lint](https://golangci-lint.run/usage/install/) — `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`

## Clone and Install

```bash
git clone https://github.com/romangurevitch/go-training.git
cd go-training
go mod tidy
```

## Verify Setup

```bash
make build   # compiles all binaries
make test    # runs all tests
make lint    # lints the codebase
```

## Start the Database

```bash
make db-up   # starts Postgres 15 via docker-compose
```

Postgres will be available at `localhost:5432` with:
- DB: `gobank`
- User: `gotrainer`
- Password: `verysecret`

## IDE Setup

**VS Code:** Install the [Go extension](https://marketplace.visualstudio.com/items?itemName=golang.go).

**GoLand / IntelliJ:** No plugins required — Go is natively supported.
