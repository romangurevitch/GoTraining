# Module 4: Temporal Orchestration

## Why Temporal?

Raw goroutines and channels are great for short-lived concurrent work, but break down for long-running operations that must survive failures, retries, and restarts. Temporal provides **durable execution**: workflows that automatically retry, resume, and track state across failures.

## Prerequisites

- **Docker Desktop** (running)
- **Go** (1.21+)

## Key Concepts

| Concept | Description |
|---|---|
| **Workflow** | Deterministic, replayable business logic. No I/O directly — delegates to Activities. |
| **Activity** | A single, retriable step (e.g., debit an account, send an email). Can have timeouts and retries. |
| **Worker** | A process that polls Temporal for tasks and executes workflows/activities. |
| **Task Queue** | The named queue a worker listens on. Workflows are dispatched to workers via task queues. |

## Terminal Setup Strategy

For the best experience, we recommend using **2 terminals**:

1.  **Terminal 1 (Services & Worker)**: Starts the background infrastructure and then runs the Worker process. **Keep this terminal open** to watch the worker logs.
2.  **Terminal 2 (Client & Signals)**: Used for all interaction—starting workflows and sending signals.

---

## Order Processing Demo (Order Saga)

The demo features an **Order Processing Saga** demonstrating durable execution, signal handling, and child workflows. While the core banking service uses standard database transactions, this module explores how to orchestrate complex, multi-step business processes that require reliable state management across many services.

### 1. Start Services & Worker (Terminal 1)

First, start the Temporal server and WireMock (they run in the background):
```bash
make temporal-up
```

Then, start the Worker in the **same terminal**. It will stay active and show you logs as work is processed:
```bash
make worker-start
```

- **Temporal Web UI**: [http://localhost:8233](http://localhost:8233)
- **WireMock (Inventory API)**: [http://localhost:8081](http://localhost:8081)

### 2. Execute Workflows (Terminal 2)

Switch to your second terminal to interact with the system. We will explore two different ways of processing orders.

#### Step 1: Automated Workflow
Run a fully automated workflow. This drives the order through every stage automatically (`PLACED` → `PICKED` → `SHIPPED` → `COMPLETED`) using Activities.

```bash
make workflow-auto
```

#### Step 2: Signal-Driven Workflow
Next, run a workflow that requires external human/system interaction. This workflow will pause at each stage and wait for you to send a signal before proceeding.

**A. Set your Workflow ID**
Pick a unique name for your order:
```bash
export ID=my-order-1
```

**B. Start the workflow**
```bash
make workflow-signal
```
*Note: This command returns immediately so you can use this same terminal for the next steps.*

### 3. Interacting with your Workflow (Terminal 2)

While the signal-driven workflow is running, you must send signals to move it forward. These commands use the `ID` variable you set above:

```bash
# Pick the order (moves from PLACED to PICKED)
make workflow-pick

# Ship the order (moves to SHIPPED)
make workflow-ship

# Mark as delivered (moves to COMPLETED)
make workflow-deliver

# Or cancel the order (before picking)
make workflow-cancel
```

> **Pro Tip**: Open the [Temporal Web UI](http://localhost:8233) to see your workflow's progress visually. You can see the event history, input/output data, and even send signals directly from the UI.

## Self-Paced Resources

- [Temporal Go SDK documentation](https://docs.temporal.io/develop/go)
- [Temporal tutorials](https://learn.temporal.io/)
- [Saga pattern explained](https://microservices.io/patterns/data/saga.html)
