# Service

The service package contains Go Bank's business logic, sitting between the HTTP API and the data store.

## Responsibilities

- Enforce business rules (e.g., cannot debit below zero)
- Orchestrate multi-step operations
- Return domain errors that the API layer maps to HTTP responses

## Challenge

Implement the service layer after completing `domain/` and `store/`.

See `internal/challenges/bank/README.md` for quest details.
