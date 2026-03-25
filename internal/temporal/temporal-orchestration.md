# ⏱️ Temporal: Durable Workflow Orchestration

Raw goroutines and channels are powerful for short-lived concurrent work, but they break down when tasks need to survive **process crashes**, **network failures**, and **restarts**. Temporal provides **durable execution**: your code runs to completion even if the machine dies halfway through.

---

## 1. The Problem Temporal Solves

Imagine a multi-step order flow: charge the customer → reserve inventory → send a confirmation email. Without Temporal, a crash between any two steps leaves inconsistent state and no automatic recovery.

```mermaid
sequenceDiagram
    participant App
    participant DB
    participant Email
    App->>DB: Charge customer ✅
    App--xEmail: 💥 Process crash — email never sent
    Note over App,Email: State is now inconsistent
```

With Temporal, the Workflow is **durable**: it resumes exactly where it left off.

```mermaid
sequenceDiagram
    participant W as Worker
    participant T as Temporal Server
    participant A as Activity

    W->>A: Charge customer
    A->>T: ✅ Done
    T->>T: 📝 Record: ActivityCompleted

    W->>A: Reserve inventory
    A->>T: ✅ Done
    T->>T: 📝 Record: ActivityCompleted

    Note over W: 💥 Worker crashes
    Note over W: ↩️Worker restarts  

    W->>T: Replay: Polls Temporal
    T-->>W: Replay: 2 steps already completed
    Note over W,T: Skips completed steps,<br>resumes here ↓

    W->>A: Send confirmation email
    A->>T: ✅ Done
    T->>T: 📝 Record: ActivityCompleted
```

---

## 2. Core Concepts

| Concept | What it is | Analogy |
|---|---|---|
| **Workflow** | Deterministic, replayable orchestration logic. Never does I/O directly. | The recipe |
| **Activity** | A single retriable step: an API call, DB write, email send. | One step in the recipe |
| **Worker** | A process that polls Temporal and executes Workflows and Activities. | The chef |
| **Task Queue** | A named queue Workers listen on; routes work to the right Workers. | The order ticket rail |
| **Signal** | An external event sent into a running Workflow to influence its state. | A customer calling to change their order |

---

## 3. How It Fits Together

```mermaid
flowchart LR
    Client["Your Code\n(Workflow Client)"]
    TS["Temporal Server\n(Durable State + Task Queue)"]
    W["Worker\n(your process)"]
    A1["Activity: Charge"]
    A2["Activity: Ship"]

    Client --"StartWorkflow"--> TS
    W --"Poll"--> TS
    TS --"Dispatch task"--> W
    W --> A1
    W --> A2
    A1 & A2 --"Result"--> TS
```

The **Temporal Server** is just a durable queue and state store — it holds no business logic. All your business logic lives in your Worker process.

---

## 4. Scaling, Limits & Best Practices

Understanding where Temporal scales freely and where it imposes hard limits is essential to designing systems that don't hit walls in production.

```mermaid
graph TB
    Client(["🖥️ Client\nStartWorkflow / Signal"])

    subgraph TS["Temporal Server  —  durable state & task queues"]
        direction TB
        TQ["Task Queue"]
        EH["📜 Event History\n⚠️ Hard limit: ~50k events\n⚠️ Payload limit: 2MB per event\n⚠️ Signal limit: ~10k pending signals"]
    end

    subgraph WP["Worker Pool  —  scale horizontally, no limit"]
        direction TB
        W1["Worker 1"]
        W2["Worker 2"]
        W3["Worker N"]
    end

    subgraph WF["Workflow  —  orchestrator"]
        direction TB
        WFN["⚠️ Single goroutine — not concurrent\n⚠️ Bounded by event history\n⚠️ Must be deterministic\n✅ Durable across crashes\n✅ Retries, timeouts, signals built-in\n✅ Queryable state at any point"]
    end

    subgraph ACT["Activities  —  the black boxes"]
        direction LR
        A1["🟦 Activity\n(DB write)"]
        A2["🟦 Activity\n(HTTP call)"]
        A3["🟦 Activity\n(Email send)"]
        AN["🟦 Activity\n(anything)"]
    end

    subgraph ASCALE["Activity Scaling"]
        direction TB
        AS1["✅ Unlimited horizontal scale\n✅ Any language, any runtime\n✅ Arbitrary I/O, side-effects OK\n⚠️ Lose durability if called outside Temporal\n⚠️ Lose retries if called directly (not via workflow.ExecuteActivity)"]
    end

    Client -->|"StartWorkflow\nor Signal"| TS
    TQ --> WP
    WP --> WF
    WF -->|"workflow.ExecuteActivity"| TQ
    TQ --> ACT
    ACT -->|"result / error"| EH
    WF -.->|"bounded by"| EH
    ACT -.->|"scales independently"| ASCALE
```

### The Core Trade-off

| | Workflow | Activity |
|---|---|---|
| **Scale** | One execution per workflow ID | Unlimited — add workers freely |
| **Concurrency** | Single goroutine (use `workflow.Go` for fan-out) | Fully concurrent, no restrictions |
| **Durability** | Every step recorded in event history | Only recorded when called via Temporal |
| **Retries** | Orchestrates retries of activities | Retried automatically by Temporal |
| **State** | Queryable at any point in time | Stateless black box |
| **Hard limits** | ~50k events, 2MB payload, ~10k signals | None imposed by Temporal |
| **I/O** | ❌ Never — use Activities instead | ✅ All I/O lives here |
| **Determinism** | ✅ Required — same inputs → same path | ❌ Not required — side-effects are fine |

### Best Practices

**Keep Workflows thin:**
- Workflows are orchestrators, not executors. They decide *what* to do and *when*; Activities do the actual work.
- A workflow that makes a direct HTTP call or writes to a DB is a bug waiting to surface during replay.

**Push all I/O into Activities:**
- Network calls, database writes, file access, randomness, current time — all belong in Activities.
- Activities are free to use any library, spawn goroutines, or block indefinitely.

**Use `continue-as-new` before hitting history limits:**
- Long-running workflows (polling loops, perpetual workflows) will eventually breach the ~50k event limit.
- Call `workflow.NewContinueAsNewError` to start a fresh execution that inherits your current state, resetting the history counter.

**Keep payloads small:**
- Each activity result is stored in the event history. Storing large blobs (images, CSVs) will exhaust the 2MB per-event limit fast.
- Store large data externally (S3, GCS) and pass only references (URLs, IDs) through Temporal.

**Size your Task Queues by workload type:**
- Separate CPU-bound activities (e.g., image processing) from I/O-bound ones (e.g., DB writes) onto different task queues.
- This lets you scale worker pools independently without one workload starving the other.

---

## 5. The Golden Rule: Workflows Must Be Deterministic

Temporal reconstructs a Workflow's state by **replaying its event history** after a crash. This means Workflow code must always make the same decisions given the same history.
For most non-deterministic functions used typically, Temporal's SDK offers deterministic alternatives. e.g.,

| ❌ Instead of this | ✅ Do this instead |
|---|---|
| `time.Now()` | `workflow.Now(ctx)` |
| `rand.Int()` | `workflow.SideEffect(...)` |
| `http.Get(url)` | Call an Activity |
| `os.Getenv(...)` | Pass as Workflow input |
| `go func() { ... }` | `workflow.Go(ctx, ...)` |

Breaking this rule causes **non-determinism errors** — Temporal detects that the replayed decisions don't match history and panics the Workflow.

Any complex non-deterministic code (e.g., network calls, I/O, database operations) should sit in an activity.

> [!TIP]
> Temporal will not re-execute activities when replaying workflows. 
> However, if an activity does not return or produce an error (i.e., the worker crashes or some other error prevents the activity from being recorded in Temporal's event history),
> the activity may be re-executed. Because of this, Temporal recommends activities be ***idempotent***.
>
> This means that, executing the same activity with the same set of inputs multiple times should be the same as executing the activity once.
> Be mindful when designing activities that they can be safely executed multiple times without causing unexpected side-effects.

## Your Next Step

Now that you understand the theory of durable orchestration, let's see it in action.

Explore the live demo in: **[Order Processing Demo](order/README.md)**.

---
[← Back to Module 4 Overview](README.md)
