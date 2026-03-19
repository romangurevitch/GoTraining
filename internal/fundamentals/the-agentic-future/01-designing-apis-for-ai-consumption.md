# Designing APIs for AI Consumption

---

## Human-Driven vs Agent-Driven API Consumption

```mermaid
graph TD
    subgraph Human["👨‍💻 Human-Driven Consumption"]
        direction TB
        HD1["Reads documentation"]
        HD2["Understands intent from prose"]
        HD3["Makes one deliberate call"]
    end

    Human ~~~ Agent

    subgraph Agent["🤖 Agent-Driven Consumption"]
        direction TB
        AD1["Reads machine-readable schema"]
        AD2["Infers intent from metadata"]
        AD3["Calls autonomously — possibly many times"]
    end

    Human -->|"APIs must serve both"| Agent
```

> An API designed only for humans will fail agents. Agents cannot read prose documentation — they need **structured metadata** that encodes intent, risk, and constraints directly in the schema.

---

## Syntax-Centric vs Intent-Centric Design

```mermaid
graph TD
    subgraph Syntax["❌ SYNTAX-CENTRIC — What AI agents see today"]
        direction TB
        SA["POST /api/v1/transfers"]
        SB["body: from_account, to_account, amount"]
        SC["Agent must infer: What is this for? Is it safe? When?"]
    end

    Syntax ~~~ Intent

    subgraph Intent["✅ INTENT-CENTRIC — What AI agents need"]
        direction TB
        IA["POST /api/v1/transfers"]
        IB["x-intent: move-funds-between-accounts"]
        IC["x-risk-profile: high — irreversible"]
        ID["x-constraints: amount ≤ daily_limit, mfa_required: true"]
        IE["x-agent-guidance: ALWAYS confirm with user before calling"]
    end

    Syntax -->|"add semantic metadata"| Intent
```

> APIs must shift from documenting **what** they accept to declaring **why** they exist, **when** they are safe to call, and **what constraints** govern them.

---

## The Four Pillars of Agent-Ready API Design

```mermaid
graph TD
    subgraph Pillars["Agent-Ready API Design"]
        direction TB
        P1["📋 **Semantic Metadata**<br/>Intent, risk profile, constraints<br/>encoded in the schema"]
        P2["🔍 **Discoverability**<br/>Machine-readable capability manifests<br/>no documentation spelunking"]
        P3["🔒 **Safety Signals**<br/>Irreversibility flags, confirmation hints<br/>rate limits, daily caps"]
        P4["📦 **Predictable Structure**<br/>Consistent envelope shapes<br/>typed errors, versioned responses"]
    end
```

---

## Semantic Metadata on OpenAPI Specs

```mermaid
graph TD
    subgraph Before["Before — Machine-readable but not intent-aware"]
        direction TB
        B1["operationId: createTransfer"]
        B2["parameters: from, to, amount"]
        B3["Agent guesses purpose and risk from field names alone"]
    end

    Before ~~~ After

    subgraph After["After — Intent-aware, agent-safe"]
        direction TB
        A1["operationId: createTransfer"]
        A2["x-intent: move-funds-between-accounts"]
        A3["x-risk-profile: high — irreversible financial operation"]
        A4["x-constraints:<br/>  max_amount: daily_limit<br/>  mfa_required: true<br/>  idempotency_key: required"]
        A5["x-agent-guidance: always confirm with user before calling"]
    end

    Before -->|"enrich with domain knowledge"| After
```

> The AI does not need to guess. The metadata **is** the contract.

---

## Structuring Responses Agents Can Reason Over

```mermaid
graph TD
    subgraph Bad["❌ Unstructured — Hard to reason over"]
        direction TB
        BR1["200 OK"]
        BR2["{'msg': 'done', 'x': true, 'ref': 'abc'}"]
        BR3["Agent must parse ambiguous fields — error-prone"]
    end

    Bad ~~~ Good

    subgraph Good["✅ Structured — Agent-friendly envelope"]
        direction TB
        GR1["201 Created"]
        GR2["{'transfer_id': 'txn_abc',<br/>'status': 'completed',<br/>'amount': 250.00,<br/>'reversible': false,<br/>'idempotency_key': 'key_xyz'}"]
        GR3["Agent knows: what was created, its state, and if it can undo"]
    end
```

---

## Risk Tiering: How Agents Should Treat Operations

```mermaid
graph TD
    direction TB
    START(["Agent selects an operation"])

    R1{"Risk profile?"}

    LOW["🟢 LOW<br/>read-only, reversible<br/>Call freely"]
    MED["🟡 MEDIUM<br/>write, reversible<br/>Log intent — proceed"]
    HIGH["🔴 HIGH<br/>write, irreversible<br/>Pause — confirm with user"]

    R1 -->|"read-only"| LOW
    R1 -->|"state-changing, reversible"| MED
    R1 -->|"irreversible / financial"| HIGH
```

---

## Agent-Safe API Checklist

```mermaid
graph TD
    C1["✅ operationId is unique and descriptive"]
    C2["✅ x-intent declares the business purpose"]
    C3["✅ x-risk-profile is set: low / medium / high"]
    C4["✅ x-constraints lists limits and preconditions"]
    C5["✅ x-agent-guidance provides calling instructions"]
    C6["✅ Idempotency-Key supported on all POST/PATCH"]
    C7["✅ Error responses include machine-readable error codes"]

    C1 --> C2 --> C3 --> C4 --> C5 --> C6 --> C7
```