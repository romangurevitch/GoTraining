# API Versioning

---

## Two Strategies for Managing Breaking Changes

```mermaid
graph TD

    subgraph Strategy2["Option 2: Extend & Contract"]
        direction LR
        EXT["📈 EXTEND<br/>Add new optional fields<br/>alongside old ones"]
        WAIT["⏳ WAIT<br/>Let clients migrate<br/>to the new fields"]
        CONTRACT["📉 CONTRACT<br/>Remove old fields<br/>once unused"]
        EXT --> WAIT --> CONTRACT
    end

    subgraph Strategy1["Option 1: Explicit Versioning"]
        direction LR
        V1["📌 /api/v1/accounts<br/>Old contract — still running"]
        V2["📌 /api/v2/accounts<br/>New contract — new features"]
        DEP["🪦 Deprecate v1<br/>Once all clients have migrated"]
        V1 --> DEP
        V2 -.->|"replaces"| V1
    end

```

> Versioning gives clients explicit contracts. Extend & Contract avoids version sprawl by evolving in place.

---

## URL Versioning: Where the Version Lives

```mermaid
graph TD
    subgraph URLPath["✅ URL Path Versioning"]
        direction TB
        U1["GET /api/v1/accounts/{id}"]
        U2["GET /api/v2/accounts/{id}"]
        U3["Visible in logs, browser, curl<br/>Easy to route at gateway level"]
    end

    URLPath ~~~ Header

    subgraph Header["⚠️ Header Versioning"]
        direction TB
        H1["GET /api/accounts/{id}"]
        H2["API-Version: 2"]
        H3["Cleaner URLs<br/>Harder to test, cache, and discover"]
    end

    Header ~~~ Query

    subgraph Query["⚠️ Query Param Versioning"]
        direction TB
        Q1["GET /api/accounts/{id}?version=2"]
        Q2["Simple to implement<br/>Pollutes every URL, easy to forget"]
    end
```

> URL path versioning is the most explicit and widely adopted. Clients know exactly which contract they are using.

---

## What Constitutes a Breaking Change?

```mermaid
graph TD
    subgraph Breaking["🔴 Breaking — REQUIRES new version"]
        direction TB
        B1["Remove a field from a response"]
        B2["Rename a field"]
        B3["Change a field type (string → int)"]
        B4["Remove or rename an endpoint"]
        B5["Change required fields in a request"]
        B6["Change error response shape"]
    end

    Breaking ~~~ NonBreaking

    subgraph NonBreaking["🟢 Non-Breaking — safe to ship"]
        direction TB
        N1["Add a new optional request field"]
        N2["Add a new response field"]
        N3["Add a new endpoint"]
        N4["Add a new optional query parameter"]
        N5["Add a new enum value (with caution)"]
    end
```

> When in doubt: if a client compiled against the old spec will break on the new response — it is a breaking change.

---

## Extend & Contract: The Two Phases in Action

```mermaid
sequenceDiagram
    autonumber
    participant OldClient as 📱 Old Client (v1 shape)
    participant NewClient as 🖥️ New Client (v2 shape)
    participant API as 🚪 API

    Note over OldClient,API: PHASE 1 — EXTEND (both shapes accepted)

    OldClient->>API: POST /accounts {"name":"Savings","balance":500}
    API->>API: adaptV1() → internal Domain
    API-->>OldClient: 201 {"id":"acc_001","name":"Savings","balance":500}

    NewClient->>API: POST /accounts {"name":"Savings","balance":500,"currency":"AUD"}
    API->>API: adaptV2() → internal Domain
    API-->>NewClient: 201 {"id":"acc_002","name":"Savings","balance":500,"currency":"AUD"}

    Note over OldClient,API: PHASE 2 — CONTRACT (old fields removed after migration)

    OldClient->>API: POST /accounts {"name":"Savings","balance":500}
    API-->>OldClient: 400 {"code":"FIELD_REQUIRED","message":"currency is now required"}
```

---

## Adapter Pattern: Many Wire Shapes, One Domain

```mermaid
graph TD
    subgraph Requests["Inbound Wire Formats"]
        direction TB
        V1REQ["📦 V1Request"]
        V2REQ["📦 V2Request"]
        EXTREQ["📦 PartnerRequest"]
    end

    subgraph Adapters["Domain Adapters"]
        direction TB
        ADAPT1["🔄 adaptV1()"]
        ADAPT2["🔄 adaptV2()"]
        ADAPT3["🔄 adaptExternal()"]
    end

    Requests --> Adapters

    subgraph DomainLayer["Core Domain"]
        direction TB
        DOMAIN["⚙️ Domain Engine"]
        BIZ["💼 Business Logic"]
        DOMAIN --> BIZ
    end

    Adapters --> DomainLayer
```

> Every API shape is an adapter. The domain model stays clean. Business logic never knows about wire formats.

---

## Version Lifecycle: From GA to Removal

```mermaid
graph LR
    BETA["🧪 Beta<br/>Not stable<br/>Opt-in only"]
    GA["✅ GA<br/>Stable contract<br/>SLA applies"]
    DEP["⚠️ Deprecated<br/>Sunset headers injected<br/>Grace period active"]
    SUNSET["🌅 Sunset<br/>410 Gone<br/>Link to migration guide"]

    BETA --> GA --> DEP --> SUNSET
```

> Every version has a lifecycle. Communicate it early, enforce it consistently.
