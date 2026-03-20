# 📦 Order Processing Demo

This demo runs two order processing workflows to show Temporal's core capabilities in action: **durable execution**, **signal handling**, and **child workflows**. All infrastructure runs locally via Docker.

---

## Prerequisites

- **Docker Desktop** (running)
- **Go** (1.21+)

---

## Terminal Setup

For the best experience, use **2 terminals**:

| Terminal | Purpose |
|---|---|
| **Terminal 1** | Runs the Temporal server + Worker process (keep open to watch logs) |
| **Terminal 2** | Sends commands — starts workflows, sends signals |

---

## Step 1: Start Services & Worker (Terminal 1)

Start the Temporal server and WireMock in the background:

```bash
make temporal-up
```

Then start the Worker. It stays active and prints logs as work arrives:

```bash
make worker-start
```

| Service | URL |
|---|---|
| Temporal Web UI | http://localhost:8233 |
| WireMock (Inventory API) | http://localhost:8081 |

---

## Step 2: Run the Workflows (Terminal 2)

### Workflow A: Automated

Drives the order through every stage automatically with no human input:

```
PLACED → PICKED → SHIPPED → COMPLETED
```

```bash
make workflow-auto
```

Watch Terminal 1 — you'll see the Worker log each Activity as it executes.

### Workflow B: Signal-Driven

This workflow **pauses at each stage** and waits for you to send a signal before advancing. It demonstrates how Temporal Workflows can react to external events.

**Set a unique Workflow ID:**

```bash
export ID=my-order-1
```

**Start the workflow** (returns immediately — the Workflow runs in the background waiting for signals):

```bash
make workflow-signal
```

---

## Step 3: Drive the Signal Workflow (Terminal 2)

Send signals to advance the paused Workflow through each stage:

```bash
make workflow-pick     # PLACED → PICKED
make workflow-ship     # PICKED → SHIPPED
make workflow-deliver  # SHIPPED → COMPLETED

# Or cancel before picking:
make workflow-cancel   # PLACED → CANCELLED
```

> **Tip**: Open the [Temporal Web UI](http://localhost:8233) while the workflow runs. You can inspect the full event history, see inputs and outputs at each step, and even send signals directly from the UI — no terminal needed.

---

## What to Look For

| Observation | What it shows |
|---|---|
| Worker logs appearing in Terminal 1 | Activities executing inside your Worker process |
| Workflow pausing between signals | Temporal durably suspending state with zero CPU usage |
| Restarting the Worker mid-workflow | The Workflow resumes from where it left off — durable execution |
| Event history in the Web UI | The full audit trail Temporal uses to replay Workflows |
