# Challenges

This directory contains all student exercises for the Go Training workshop.

## Structure

```
challenges/
├── basics/
│   ├── 01-structs-and-pointers/
│   ├── 02-interfaces-and-receivers/
│   └── ...
└── bank/         # Go Bank service quests
```

## basics/

Short, focused exercises for learning Go fundamentals. Each directory contains a `README.md` and a `challenge.go` file.
Inspired by the [ConcurrencyWorkshop](https://github.com/) fixme pattern.

## bank/

The **[Go Bank Transfer Quest](bank/README.md)** is our main challenge! 

You'll implement the `POST /v1/transfers` API endpoint in a pre-scaffolded service, focusing on:
- Idiomatic HTTP handler patterns using Gin.
- OpenTelemetry tracing and structured logging with `slog`.
- JWT authentication and scope-based authorisation.
- Table-driven unit testing for handlers.

Everything below the API layer is pre-built so you can focus on building production-grade APIs.

Run tests with: `make test-bank`

