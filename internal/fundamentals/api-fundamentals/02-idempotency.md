# Idempotency

---

## What Is Idempotency?

```mermaid
graph TD
    DEF["**Idempotent operation:**<br/>Calling it once produces the same result<br/>as calling it N times"]

    DEF ~~~ Idempotent

    subgraph Idempotent["✅ Idempotent"]
        direction TB
        I1["GET  /accounts/acc_123<br/>→ same account every time"]
        I2["PUT  /accounts/acc_123 body={name: 'Savings'}<br/>→ same final state"]
        I3["DELETE /accounts/acc_123<br/>→ deleted whether called once or twice"]
    end

    Idempotent ~~~ NotIdempotent

    subgraph NotIdempotent["❌ NOT Idempotent by default"]
        direction TB
        N1["POST /payments body={amount: 250.00}<br/>→ creates a NEW payment each call"]
        N2["POST /notifications/send<br/>→ sends a NEW email each call"]
    end
```

> If a retry is safe and produces no extra side effects, the operation is idempotent.

---

## HTTP Methods and Idempotency

```mermaid
graph TD
    subgraph IdempotentMethods["✅ Idempotent Methods"]
        direction TB
        GET["**GET**<br/>📖 Read-only<br/>✅ Idempotent<br/>✅ Cacheable"]
        HEAD["**HEAD**<br/>📋 Metadata only<br/>✅ Idempotent"]
        PUT["**PUT**<br/>🔄 Full replace<br/>✅ Idempotent"]
        DELETE["**DELETE**<br/>🗑️ Remove<br/>✅ Idempotent"]
        OPTIONS["**OPTIONS**<br/>🔍 Preflight / CORS<br/>✅ Idempotent"]
    end

    IdempotentMethods ~~~ NonIdempotentMethods

    subgraph NonIdempotentMethods["❌ Non-Idempotent Methods"]
        direction TB
        PATCH["**PATCH**<br/>✏️ Partial update<br/>⚠️ Depends on implementation"]
        POST["**POST**<br/>➕ Create / action<br/>❌ NOT idempotent<br/>Use Idempotency-Key header"]
    end
```

---

## The Double-Charge Problem

```mermaid
sequenceDiagram
    autonumber
    participant App as 📱 Mobile App
    participant API as 🚪 API
    participant DB as 🐘 Database

    Note over App,DB: ❌ WITHOUT idempotency key

    App->>API: POST /payments {"amount": 250.00}
    API->>DB: INSERT INTO payments ...
    DB-->>API: payment_id = pay_001
    Note over App,API: ⏱️ Network timeout — app never receives response

    App->>API: POST /payments {"amount": 250.00}  ← retry
    API->>DB: INSERT INTO payments ...
    DB-->>API: payment_id = pay_002
    API-->>App: 201 Created {"id": "pay_002"}

    Note over App,DB: 💥 Customer charged TWICE
```

---

## The Fix: Idempotency-Key Header

```mermaid
sequenceDiagram
    autonumber
    participant App as 📱 Mobile App
    participant API as 🚪 API
    participant Cache as ⚡ Redis / DB
    participant DB as 🐘 Database

    Note over App,DB: ✅ WITH Idempotency-Key

    App->>API: POST /payments {"amount": 250.00}<br/>Idempotency-Key: key_abc123
    API->>Cache: Check key_abc123 → not found
    API->>DB: INSERT INTO payments ...
    DB-->>API: payment_id = pay_001
    API->>Cache: Store key_abc123 → {"id":"pay_001","status":201}
    API-->>App: 201 Created {"id": "pay_001"}

    Note over App,API: ⏱️ Network timeout — app retries

    App->>API: POST /payments {"amount": 250.00}<br/>Idempotency-Key: key_abc123  ← SAME key
    API->>Cache: Check key_abc123 → FOUND
    API-->>App: 201 Created {"id": "pay_001"}  ← cached response

    Note over App,DB: ✅ No duplicate. Same result. Safe to retry.
```

> The key must be generated **client-side** before the request is sent, and reused on every retry of the same logical operation.

---

## Idempotency Key Lifecycle

```mermaid
graph TD
    direction TB
    GEN["📱 Client generates key<br/>UUID or similar<br/>e.g. key_abc123"]

    FIRST["1️⃣ First request<br/>Key not in store → process normally<br/>Store key + response"]

    RETRY["🔄 Retry (timeout / 5xx)<br/>Key already in store → return cached response<br/>No side effects"]

    EXPIRE["⏱️ Key expires<br/>After TTL (e.g. 24 hours)<br/>Client must generate new key for new operations"]

    ERROR["❌ Key conflict<br/>Same key, different body → 422 Unprocessable<br/>Client bug: one key per logical operation"]

    GEN --> FIRST --> RETRY --> EXPIRE
    FIRST --> ERROR
```

---

## DELETE Idempotency: Already Gone Is Fine

```mermaid
sequenceDiagram
    autonumber
    participant Client as 📱 Client
    participant API as 🚪 API
    participant DB as 🐘 Database

    Client->>API: DELETE /accounts/acc_123
    API->>DB: DELETE WHERE id = acc_123
    DB-->>API: 1 row deleted
    API-->>Client: 204 No Content

    Note over Client,API: ⏱️ Client times out, retries

    Client->>API: DELETE /accounts/acc_123  ← retry
    API->>DB: DELETE WHERE id = acc_123
    DB-->>API: 0 rows deleted (already gone)
    API-->>Client: 204 No Content  ← same response

    Note over Client,API: ✅ Resource is gone either way. Safe.
```

> Return `204` (not `404`) on a repeat DELETE. The desired state — "this resource must not exist" — has been achieved.

---

## PATCH: Conditional Idempotency

```mermaid
graph TD

    subgraph Idem["✅ PATCH — Idempotent"]
        direction TB
        Q1["PATCH /accounts/acc_123<br/>body: {'name': 'New Savings Account'}"]
        Q2["Call 1: name updated"]
        Q3["Call 2: same name, same result"]
        Q1 --> Q2 --> Q3
    end

    Idem ~~~ NonIdem

    subgraph NonIdem["❌ PATCH — NOT Idempotent"]
        direction TB
        P1["PATCH /accounts/acc_123<br/>body: {'balance': {'increment': 100}}"]
        P2["Call 1: balance = 500 → 600"]
        P3["Call 2: balance = 600 → 700"]
        P1 --> P2 --> P3
    end
```

> A PATCH that **sets an absolute value** is idempotent. A PATCH that **increments** is not — use an Idempotency-Key.

---

## Idempotency at Financial Institutions: What Gets a Key

```mermaid
graph TD
    subgraph MUST_sub["🔴 MUST use Idempotency-Key"]
        direction TB
        M1["POST /payments"]
        M2["POST /transfers"]
        M3["POST /direct-debits"]
        M4["POST /notifications/send"]
    end

    MUST_sub ~~~ SHOULD_sub

    subgraph SHOULD_sub["🟡 SHOULD use Idempotency-Key"]
        direction TB
        S1["POST /accounts (onboarding)"]
        S2["POST /loans/apply"]
    end

    SHOULD_sub ~~~ NOTNEED_sub

    subgraph NOTNEED_sub["🟢 Does NOT need Idempotency-Key"]
        direction TB
        N1["GET (all read operations)"]
        N2["PUT /accounts/{id}"]
        N3["DELETE /accounts/{id}"]
    end
```

> Every operation with a **financial side effect** must be idempotent. Retries happen. Networks fail. Design for it.
