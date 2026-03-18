# Module 4: Temporal Orchestration

> Workflow and activity code for this module is under development.

## Why Temporal?

Raw goroutines and channels are great for short-lived concurrent work, but break down for long-running operations that must survive failures, retries, and restarts. Temporal provides **durable execution**: workflows that automatically retry, resume, and track state across failures.

## Key Concepts

| Concept | Description |
|---|---|
| **Workflow** | Deterministic, replayable business logic. No I/O directly — delegates to Activities. |
| **Activity** | A single, retriable step (e.g., debit an account, send an email). Can have timeouts and retries. |
| **Worker** | A process that polls Temporal for tasks and executes workflows/activities. |
| **Task Queue** | The named queue a worker listens on. Workflows are dispatched to workers via task queues. |

## TransferFunds Demo

The demo shows a **Saga pattern** for reliable fund transfers:

1. Debit source account (Activity)
2. Credit destination account (Activity)
3. On failure: compensate by reversing the debit (Activity)

## Self-Paced Resources

- [Temporal Go SDK documentation](https://docs.temporal.io/develop/go)
- [Temporal tutorials](https://learn.temporal.io/)
- [Saga pattern explained](https://microservices.io/patterns/data/saga.html)

## Running Temporal Locally

```bash
# Start Temporal dev server (requires Temporal CLI)
temporal server start-dev

# Start the worker (once implemented)
go run ./cmd/worker/...
```
