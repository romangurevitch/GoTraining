# Go Bank Transfer Quest — Design Spec

**Date:** 2026-03-18
**Status:** Approved
**Audience:** Go training participants (coming from another language, ~2-3 hours available)

---

## Overview

A focused hands-on exercise where participants implement a `POST /v1/transfers` endpoint in a pre-scaffolded Go bank service. The goal is to teach idiomatic Go HTTP handler patterns, OpenTelemetry tracing, structured logging, JWT authentication/authorization, and handler testing — without getting distracted by DB or store concerns.

Everything below the service layer is pre-built. The **account handler is the fully working reference** — participants read it, understand every pattern, then replicate it for transfers.

All code in this spec is derived directly from `go-training-cba-solution/` and updated with Go 1.26.1 improvements and identified gaps (missing tracing, JWT auth, typed errors in the client, logrus → slog).

---

## Gaps Identified in Existing Challenge vs Reference Solution

| Area | Current Challenge (`GoTrainig/`) | Reference Solution | What We Add |
|------|----------------------------------|-------------------|-------------|
| Auth | None | Basic Auth (`gin.BasicAuth`) | JWT Bearer + scope-based authz (upgrade) |
| Tracing | None | `logger.TraceDuration` only (log-based) | OTel spans + Jaeger (upgrade) |
| Request ID | None | None | `RequestIDMiddleware` (new) |
| Client | Thin, no typed errors | Typed client + `DoRequest` + `APIError` | Port pattern to `pkg/client/bank/` |
| Handler structure | Flat `handlers.go` | One package per resource | Split to `api/account/`, `api/transfer/` |
| Error handling | `fmt.Errorf("status %d")` | `error.NewAPIError(ctx, status, msg, err)` | Adopt reference pattern |
| Store naming | `store/memory.go` | N/A | `repository/postgres/` (remove in-memory) |
| Logger | Basic logrus usage | `logger.WithContext(ctx)` package (logrus) | Replace with `log/slog` stdlib + `slog-otel` for trace correlation |

---

## Project Structure

The module root is `GoTrainig/`. Existing `cmd/` entry points are unchanged. The structure mirrors `go-training-cba-solution/` with one package per resource.

```
GoTrainig/                                  # module root (go.mod here)
├── cmd/                                    # unchanged — existing entry points
│   ├── bank-api/main.go
│   └── bank-cli/main.go
├── internal/bank/
│   ├── domain/
│   │   ├── account.go                      # Account, AccountStatus, sentinel errors
│   │   └── transaction.go                  # Transaction, TransactionType
│   ├── repository/
│   │   ├── repository.go                   # Repository interface (renamed from store/)
│   │   └── postgres/
│   │       └── repository.go               # go-jet Postgres implementation (pre-built)
│   ├── service/
│   │   └── service.go                      # Service interface + BankService (pre-built)
│   └── api/
│       ├── server.go                       # router wiring — participants wire transfer routes here
│       ├── middleware/
│       │   ├── logging.go                  # slog-gin middleware — replaces logrus JSONLogMiddleware
│       │   ├── tracing.go                  # OTel span per request (NEW — not in reference)
│       │   ├── requestid.go                # RequestID inject + propagate (NEW — not in reference)
│       │   └── auth.go                     # JWTMiddleware + RequireScope + ClaimsFromCtx (NEW — replaces BasicAuth)
│       ├── auth/
│       │   └── handler.go                  # POST /v1/token — issues JWT (pre-built)
│       ├── account/
│       │   └── handler.go                  # REFERENCE — fully implemented, do not modify
│       └── transfer/
│           └── handler.go                  # PARTICIPANTS IMPLEMENT
├── pkg/
│   ├── api/
│   │   └── error/
│   │       └── error.go                    # APIError type + NewAPIError — from reference solution
│   ├── http/
│   │   └── http.go                         # DoRequest, GetURL, HeaderApplicationJSON — from reference solution
│   └── client/
│       └── bank/
│           ├── client.go                   # Client interface + GetAccount/GetToken (pre-built reference)
│           └── api/
│               └── types.go                # Typed request/response structs
└── docs/
    └── openapi/
        ├── accounts.yaml                   # pre-written reference spec
        └── transfers.yaml                  # partially filled — participants complete Step 1
```

### Naming Changes from Current Challenge

| Old | New | Reason |
|-----|-----|--------|
| `store/` | `repository/` | Repository pattern — intent is clearer |
| `store/memory.go` | removed | Real Postgres only; in-memory was scaffolding |
| `api/handlers.go` | `api/account/handler.go`, `api/transfer/handler.go` | One package per resource — mirrors reference solution |
| `internal/bank/client/` | `pkg/client/bank/` | Reusable across projects, not bank-internal |
| `pkg/api/error/` | same | Adopted directly from reference solution |

---

## Dependencies

### Go Version
Go 1.26.1 — same as `go-training-cba-solution/go.mod`. Use range-over-func iterators where natural.

### Key Packages
| Package | Source | Purpose |
|---------|--------|---------|
| `github.com/gin-gonic/gin v1.10.1` | reference go.mod | HTTP router and middleware |
| `github.com/go-jet/jet/v2 v2.14.1` | reference go.mod | Type-safe SQL builder (pre-built) |
| `log/slog` | stdlib (Go 1.21+) | Structured JSON logging — **replaces logrus** |
| `slog.NewMultiHandler` | stdlib (Go 1.26+) | Fan out to JSON stdout + OTel bridge |
| `github.com/remychantenay/slog-otel` | NEW | Injects `trace_id`/`span_id` from OTel context into slog JSON output |
| `github.com/samber/slog-gin` | NEW | Gin request logging middleware using slog |
| `github.com/go-playground/validator/v10 v10.30.1` | reference go.mod | Struct/field validation |
| `go.opentelemetry.io/otel` | NEW | Tracing SDK |
| `go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp` | NEW | OTLP exporter → Jaeger |
| `go.opentelemetry.io/contrib/bridges/otelslog` | NEW | OTel Logs bridge (used with `NewMultiHandler`) |
| `github.com/golang-jwt/jwt/v5` | NEW | JWT parsing and validation (replaces BasicAuth) |
| `github.com/stretchr/testify v1.11.1` | reference go.mod | Assertions and mocking |
| `github.com/spf13/cobra v1.10.2` | reference go.mod | CLI framework |
| `github.com/spf13/viper v1.21.0` | reference go.mod | Config management |

### Docker Compose Services
```yaml
services:
  postgres:   # PostgreSQL 15 — same as reference solution
  jaeger:     # All-in-one Jaeger (UI: localhost:16686, OTLP: localhost:4318) — NEW
  bank-api:   # The Go service
```

---

## Pre-Built Scaffold (Participants Do Not Implement)

### `pkg/api/error/error.go` — Adopted from reference, updated to slog

The reference solution's `pkg/api/error/error.go` structure is kept — tuple return fits Gin's `c.JSON` directly. The internal logrus call is replaced with `slog.ErrorContext`:

```go
// Source: go-training-cba-solution/pkg/api/error/error.go
// Change: logger.WithContext(ctx).WithError(err)... → slog.ErrorContext(ctx, ...)

var ErrInternalServerError = errors.New("internal server error")

func NewAPIError(ctx context.Context, status int, msg string, err error) (int, *APIError) {
    slog.ErrorContext(ctx, msg,
        slog.Int("status", status),
        slog.Any("error", err),
    )
    return status, &APIError{Message: msg}
}

type APIError struct {
    Message string `json:"message"`
}

func (e APIError) Error() string {
    return fmt.Sprintf("api error: %v", e.Message)
}
```

Participants call `c.JSON(apierror.NewAPIError(ctx, http.StatusBadRequest, "bad request", err))` — unchanged from reference. The slog call inside automatically carries `trace_id`/`span_id` via the `OtelHandler` set up in `main.go`.

### `internal/bank/api/middleware/logging.go` — slog-gin (replaces logrus JSONLogMiddleware)

The reference solution used a hand-rolled `JSONLogMiddleware` with logrus. We replace it with `samber/slog-gin` which natively integrates OTel trace correlation — the reference's main logging gap:

```go
// Reference had: hand-rolled JSONLogMiddleware using logrus (logger.go)
// Upgrade: slog-gin with OTel trace_id/span_id injection via slog-otel handler chain

import (
    sloggin "github.com/samber/slog-gin"
)

func LoggingMiddleware(logger *slog.Logger) gin.HandlerFunc {
    return sloggin.NewWithConfig(logger, sloggin.Config{
        DefaultLevel:     slog.LevelInfo,
        ClientErrorLevel: slog.LevelWarn,
        ServerErrorLevel: slog.LevelError,
        WithTraceID:      true,  // reads from OTel context — no manual extraction needed
        WithSpanID:       true,
        WithRequestID:    true,  // from X-Request-Id header set by RequestIDMiddleware
        // Fields logged per request: time, method, path, status, latency, trace_id, span_id, request_id
    })
}
```

Logger is bootstrapped in `main.go` (pre-built) using Go 1.26's `slog.NewMultiHandler`:

```go
// cmd/bank-api/main.go (pre-built)
jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
otelBridge  := otelslog.NewHandler("bank-api")  // ships logs as OTel log signal

// slog-otel wraps both: injects trace_id/span_id into JSON output AND forwards to OTel
otelEnricher := slogotel.OtelHandler{Next: slog.NewMultiHandler(jsonHandler, otelBridge)}

logger := slog.New(otelEnricher)
slog.SetDefault(logger)
```

Result: every `slog.InfoContext(ctx, ...)` call anywhere in the codebase automatically includes `trace_id` and `span_id` from whatever OTel span is active in `ctx` — zero per-handler boilerplate.

### `internal/bank/api/middleware/auth.go` — NEW (replaces BasicAuth from reference)

The reference uses `gin.BasicAuth`. We upgrade to JWT + scope-based authorization:

```go
// Reference used: gin.BasicAuth(gin.Accounts{"gois": "thebest"})
// Upgrade: JWT Bearer token with custom scope claim

type Claims struct {
    Scope string `json:"scope"`
    jwt.RegisteredClaims          // provides Subject (sub), ExpiresAt, etc.
}

// JWTMiddleware validates the Bearer token and injects Claims into context.
// Used identically to middleware.Auth() in the reference — just registered on a route group.
func JWTMiddleware(secret string) gin.HandlerFunc

// RequireScope checks claims.Scope contains the required scope.
// Returns 403 if not — same pattern as if the BasicAuth check failed.
func RequireScope(scope string) gin.HandlerFunc

// ClaimsFromCtx extracts the injected Claims from context — used inside handlers.
func ClaimsFromCtx(ctx context.Context) *Claims
```

### `internal/bank/api/server.go` — Based on reference `account/server.go`

The reference server setup (`go-training-cba-solution/internal/server/rest/account/server.go`) is the pattern:

```go
// Reference: go-training-cba-solution/internal/server/rest/account/server.go
func NewServer(store store.Store) rest.Server {
    r := gin.New()
    r.Use(middleware.JSONLogMiddleware())
    r.Use(gin.Recovery())
    // ...routes...
}
```

Our server adds the three new middleware and a second route group. Accounts group is pre-wired as reference:

```go
func NewServer(svc service.Service, cfg *config.Config) *gin.Engine {
    r := gin.New()
    r.Use(middleware.RequestIDMiddleware())  // NEW
    r.Use(middleware.TracingMiddleware())    // NEW
    r.Use(middleware.JSONLogMiddleware())    // from reference + trace_id field
    r.Use(gin.Recovery())

    accountHandler := account.New(svc)
    r.POST("/v1/token", authHandler.IssueToken)

    // Pre-wired reference — participants read this to understand the pattern
    accounts := r.Group("/v1/accounts")
    accounts.Use(middleware.JWTMiddleware(cfg.JWTSecret))
    {
        accounts.GET("/:id", middleware.RequireScope("accounts:read"), accountHandler.GetAccount)
        accounts.POST("", middleware.RequireScope("accounts:write"), accountHandler.CreateAccount)
    }

    // TODO: Register POST /v1/transfers with JWTMiddleware and RequireScope("transfers:write")
    // Pattern: identical to the accounts group above

    return r
}
```

### `internal/bank/api/account/handler.go` — REFERENCE, adopted from reference solution

Derived from `go-training-cba-solution/internal/server/rest/account/handler.go`. Upgraded with OTel spans and slog (reference had neither):

```go
// Reference: go-training-cba-solution/internal/server/rest/account/handler.go
// Changes vs reference:
//   - logrus logger.WithContext(ctx).WithField(...).Info() → slog.InfoContext(ctx, ...)
//   - OTel span added per handler (reference had no tracing)
//   - service layer instead of direct store access

func (h *Handler) GetAccount(c *gin.Context) {
    ctx := c.Request.Context()

    ctx, span := otel.Tracer("bank").Start(ctx, "account.get")
    defer span.End()

    id := c.Param("id")
    // slog — replaces logger.WithContext(ctx).WithField("account_id", id).Info()
    slog.InfoContext(ctx, "get account", slog.String("account_id", id))

    result, err := h.svc.GetAccount(ctx, id)
    switch {
    case errors.Is(err, domain.ErrAccountNotFound):
        c.JSON(apierror.NewAPIError(ctx, http.StatusNotFound, "account not found", err))
    case err != nil:
        c.JSON(apierror.NewAPIError(ctx, http.StatusInternalServerError, "could not get account", err))
    default:
        span.SetAttributes(attribute.String("account.owner", result.Owner))
        c.JSON(http.StatusOK, toAccountResponse(result))
    }
}

func (h *Handler) CreateAccount(c *gin.Context) {
    ctx := c.Request.Context()

    ctx, span := otel.Tracer("bank").Start(ctx, "account.create")
    defer span.End()

    var req api.CreateAccountRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(apierror.NewAPIError(ctx, http.StatusBadRequest, "bad request", err))
        return
    }

    result, err := h.svc.CreateAccount(ctx, req.Owner)
    switch {
    case errors.Is(err, domain.ErrAccountAlreadyExists):
        c.JSON(apierror.NewAPIError(ctx, http.StatusConflict, "account already exists", err))
    case err != nil:
        c.JSON(apierror.NewAPIError(ctx, http.StatusInternalServerError, "could not create account", err))
    default:
        slog.InfoContext(ctx, "account created", slog.String("account_id", result.ID))
        c.JSON(http.StatusCreated, toAccountResponse(result))
    }
}
```

### `pkg/http/http.go` — Adopted directly from reference solution

Copied verbatim from `go-training-cba-solution/pkg/http/http.go`. One gap fixed: context is now actually used in `DoRequest` (it was `_ context.Context` in the reference):

```go
// Source: go-training-cba-solution/pkg/http/http.go
// Fix: context is passed to http.NewRequestWithContext instead of being ignored

func DoRequest(ctx context.Context, client *http.Client, r *http.Request, expectedResponses ...int) ([]byte, error) {
    resp, err := client.Do(r.WithContext(ctx))  // was: client.Do(r) with ctx ignored
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if err = checkResponse(resp, expectedResponses...); err != nil {
        return nil, err
    }
    return io.ReadAll(resp.Body)
}

func checkResponse(resp *http.Response, expectedResponses ...int) error {
    for _, expected := range expectedResponses {
        if expected == resp.StatusCode {
            return nil
        }
    }
    return getAPIError(resp)
}

func getAPIError(resp *http.Response) error {
    if resp.StatusCode == http.StatusUnauthorized {
        return apierror.APIError{Message: resp.Status}
    }
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    var apiErr *apierror.APIError
    if err = json.Unmarshal(body, &apiErr); err != nil {
        return apierror.ErrInternalServerError
    }
    return apiErr
}

func GetURL(baseURL, p, queryFormat string, args ...any) (string, error) {
    u, err := url.Parse(baseURL)
    if err != nil {
        return "", err
    }
    u.Path = path.Join(u.Path, p)
    if queryFormat != "" {
        u.RawQuery = fmt.Sprintf(queryFormat, args...)
    }
    return u.String(), nil
}

func HeaderApplicationJSON() (key, value string) {
    return "Content-Type", "application/json"
}
```

Note: `...interface{}` → `...any` (Go 1.18+ alias, cleaner).

### `internal/bank/api/auth/handler.go` — Token issuing endpoint (pre-built)

```
POST /v1/token
Body: { "sub": "alice", "scope": "accounts:read transfers:write" }
Returns: { "token": "<signed JWT>" }
```

Hardcoded secret from config. Participants use `curl` to get a token for manual testing.

---

## OpenAPI Spec Design (Step 1 — 20 min)

### Reference: `docs/openapi/accounts.yaml`

Complete spec for `GET /v1/accounts/{id}` and `POST /v1/accounts`. Participants read this first.

### Participant Task: `docs/openapi/transfers.yaml`

Partially filled. Participants complete the TODOs using `accounts.yaml` as reference:

```yaml
openapi: 3.0.3
info:
  title: Bank Transfer API
  version: 1.0.0

paths:
  /v1/transfers:
    post:
      summary: Transfer funds between two accounts
      security:
        - bearerAuth: []
      # TODO: Add operationId
      requestBody:
        required: true
        content:
          application/json:
            schema:
              # TODO: Define the request body schema
              #   Fields: from_account_id, to_account_id, amount
              #   Which fields are required? What types?
      responses:
        "200":
          # TODO: Define success response
          #   Hint: { "status": "completed" }
        "400":
          # TODO: When does this happen? Define error response body
          #   Hint: malformed or missing request fields — same shape as accounts 400
        "403":
          # TODO: Two distinct 403 cases — what are they?
          #   Hint 1: missing or wrong scope on the JWT token (middleware rejects)
          #   Hint 2: authenticated user is not the owner of from_account_id (handler rejects)
        "404":
          # TODO: Account not found
          #   Hint: either source or destination — response body does not distinguish (intentional simplification)
        "422":
          # TODO: Business rule violation — account exists but transfer cannot proceed
          #   Hint: insufficient funds, or account is locked
        "500":
          # TODO: Unexpected internal error — same shape as accounts 500

components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
  schemas:
    # TODO: Define TransferRequest and ErrorResponse schemas
    # Hint: ErrorResponse is identical to the one in accounts.yaml — reuse it
```

---

## Transfer Handler Quest (Steps 2-4)

### Step 2 — Wire Middleware (15 min)

In `api/server.go`, the accounts group is pre-wired as the reference pattern. Participants add the transfer group below it:

```go
// PRE-WIRED REFERENCE — read this to understand the pattern
accounts := r.Group("/v1/accounts")
accounts.Use(middleware.JWTMiddleware(cfg.JWTSecret))
{
    accounts.GET("/:id", middleware.RequireScope("accounts:read"), accountHandler.GetAccount)
    accounts.POST("", middleware.RequireScope("accounts:write"), accountHandler.CreateAccount)
}

// TODO: Add transfer group — same pattern as accounts above
// transfers := r.Group(...)
// transfers.Use(middleware.JWTMiddleware(...))
// {
//     transfers.POST(..., middleware.RequireScope(...), transferHandler.CreateTransfer)
// }
```

### Step 3 — Implement the Handler (60-80 min)

`Handler` struct follows the exact same shape as the reference account handler:

```go
// Same pattern as api/account/handler.go
type Handler struct {
    svc service.Service  // interface — enables mock injection in tests
}

func New(svc service.Service) *Handler {
    return &Handler{svc: svc}
}
```

Skeleton in `api/transfer/handler.go` — each TODO points directly to the reference line it mirrors:

```go
package transfer

// REFERENCE: api/account/handler.go — this handler follows the identical pattern.
// Read CreateAccount there first, then implement CreateTransfer here.

func (h *Handler) CreateTransfer(c *gin.Context) {
    ctx := c.Request.Context()

    // TODO 1: Parse and validate request body
    //   var req api.CreateTransferRequest
    //   if err := c.ShouldBindJSON(&req); err != nil { ... }
    //   Reference: account/handler.go CreateAccount — identical binding pattern

    // TODO 2: Start an OTel span (NEW — not in reference solution, this is the upgrade)
    //   ctx, span := otel.Tracer("bank").Start(ctx, "transfer.create")
    //   defer span.End()
    //   span.SetAttributes(
    //       attribute.String("from_account_id", req.FromAccountID),
    //       attribute.String("to_account_id", req.ToAccountID),
    //       attribute.Float64("amount", req.Amount),
    //   )

    // TODO 3: Verify ownership (JWT sub claim must match from_account owner)
    //   claims := middleware.ClaimsFromCtx(ctx)           // extracts sub from JWT
    //   fromAccount, err := h.svc.GetAccount(ctx, req.FromAccountID)
    //   Handle errors immediately — same switch pattern as reference:
    //     ErrAccountNotFound → 404, other → 500
    //   if fromAccount.Owner != claims.Subject {
    //       c.JSON(apierror.NewAPIError(ctx, http.StatusForbidden, "forbidden", nil))
    //       return
    //   }
    //   Note: this GetAccount is authorization-only. service.Transfer loads it again
    //   for business logic — separation of concerns is intentional.

    // TODO 4: Call service and map errors
    //   err = h.svc.Transfer(ctx, req.FromAccountID, req.ToAccountID, req.Amount)
    //   Reference: account/handler.go uses switch + errors.Is — use the same pattern:
    //   switch {
    //   case errors.Is(err, domain.ErrAccountNotFound):
    //       c.JSON(apierror.NewAPIError(ctx, http.StatusNotFound, "account not found", err))
    //   case errors.Is(err, domain.ErrInsufficientFunds):
    //       c.JSON(apierror.NewAPIError(ctx, http.StatusUnprocessableEntity, "insufficient funds", err))
    //   case errors.Is(err, domain.ErrAccountLocked):
    //       c.JSON(apierror.NewAPIError(ctx, http.StatusUnprocessableEntity, "account locked", err))
    //   case err != nil:
    //       c.JSON(apierror.NewAPIError(ctx, http.StatusInternalServerError, "transfer failed", err))
    //   default:
    //       ...
    //   }

    // TODO 5: Log and return
    //   Reference: account/handler.go uses slog.InfoContext(ctx, "...", slog.String(...))
    //   On success:
    //     slog.InfoContext(ctx, "transfer completed",
    //         slog.String("from_account_id", req.FromAccountID),
    //         slog.String("to_account_id",   req.ToAccountID),
    //         slog.Float64("amount",         req.Amount),
    //     )
    //   trace_id and span_id are injected automatically by the OtelHandler — no extra code needed
    //   c.JSON(http.StatusOK, api.TransferResponse{Status: "completed"})
}
```

### Step 4 — Write Handler Tests (20-30 min)

Based directly on `go-training-cba-solution/internal/server/rest/account/handler_test.go` — the exact same structure, just for the transfer handler.

The skeleton pre-provides (mirroring the reference test structure):
- `mockSvc` generated by mockery — same as `storemocks.NewStore(t)` in the reference
- `setupRouter(svc service.Service) *gin.Engine` — mirrors `NewServer(tt.fields.store(t)).(*accountServer)`
- `testToken(sub, scope string) string` — issues a signed JWT (replaces the `const auth = "Basic Z29pczp0aGViZXN0"` in the reference)
- Pre-written: **happy path (200)** and **invalid body (400)**

The reference test structure to replicate:

```go
// Reference: go-training-cba-solution/internal/server/rest/account/handler_test.go
func Test_accountServer_handleGetAccounts(t *testing.T) {
    type fields struct {
        store func(t *testing.T) store.Store
    }
    tests := []struct {
        name     string
        fields   fields
        args     args
        want     interface{}
        wantCode int
    }{
        {
            name: "success",
            fields: fields{
                store: func(t *testing.T) store.Store {
                    m := storemocks.NewStore(t)
                    m.EXPECT().FindAccountsByEmails(...).Return(...).Times(1)
                    return m
                },
            },
            want:     newMapper().toAccountResponse(&testmocks.Account1),
            wantCode: http.StatusOK,
        },
        // ...
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            gin.SetMode(gin.TestMode)
            server := NewServer(tt.fields.store(t)).(*accountServer)
            w := httptest.NewRecorder()
            req := httptest.NewRequest(http.MethodGet, tt.args.path, nil)
            req.Header.Set("Authorization", tt.args.auth)
            server.engine.ServeHTTP(w, req)
            httppkg.CheckAPIResponse(t, w.Body, w.Code, tt.wantCode, tt.want)
        })
    }
}
```

Participants add three cases to the pre-written transfer test:

| Case | Mock Setup | Expected |
|------|-----------|----------|
| Happy path | `GetAccount` → owner matches; `Transfer` → nil | 200 — pre-written |
| Invalid body | no mock needed | 400 — pre-written |
| Wrong owner | `GetAccount` → `Owner: "bob"`, token `sub: "alice"` | 403 — participant writes |
| Insufficient funds | `GetAccount` → owner matches; `Transfer` → `ErrInsufficientFunds` | 422 — participant writes |
| Account not found | `GetAccount` → `ErrAccountNotFound` | 404 — participant writes |

---

## Observability

### Tracing (NEW — not in reference solution)

The reference solution only had `logger.TraceDuration` (log-based timing). We add full OTel:

- OTel SDK bootstrapped in `cmd/bank-api/main.go` (pre-built)
- `TracingMiddleware` starts `"http.request"` span per request — wraps every handler
- Participants add `"transfer.create"` child span inside the handler (TODO 2)
- OTLP HTTP exporter → Jaeger at `localhost:4318`
- Traces visible at `localhost:16686` immediately after first request

Span hierarchy participants will see in Jaeger:
```
http.request
  └── transfer.create
        └── (service spans added by pre-built service layer)
```

### Logging — slog (replaces logrus from reference solution)

The reference used `logger.WithContext(ctx)` (a thin logrus wrapper). We replace it with stdlib `log/slog` — no wrapper package needed. The `OtelHandler` bootstrapped in `main.go` means every `slog.*Context(ctx, ...)` call automatically carries `trace_id` and `span_id`.

**Migration from reference pattern to slog:**

```go
// Reference (logrus):
logger.WithContext(ctx).WithField("account_id", id).Info()
logger.WithContext(ctx).WithError(err).WithField("status", 500).Error("storage error")

// slog equivalent (what participants write):
slog.InfoContext(ctx, "get account", slog.String("account_id", id))
slog.ErrorContext(ctx, "storage error", slog.Int("status", 500), slog.Any("error", err))
```

**`TraceDuration` equivalent** — the reference had `logger.TraceDuration(ctx, time.Now(), "GetAccount")`. With slog, defer-log duration inline:

```go
// Reference: defer logger.TraceDuration(ctx, time.Now(), "GetAccount")
// slog equivalent:
start := time.Now()
defer func() {
    slog.DebugContext(ctx, "GetAccount", slog.Duration("duration", time.Since(start)))
}()
```

Rules participants follow (demonstrated in account handler reference):
- Always `slog.InfoContext(ctx, ...)` — never bare `slog.Info(...)` (loses OTel correlation)
- Typed attrs: `slog.String(...)`, `slog.Float64(...)`, `slog.Any("error", err)` — not `"key", value` pairs
- Log errors once at the handler boundary using `slog.Any("error", err)` — don't re-log in service/repository
- One log line per request outcome (success or failure)

### Correlation (NEW)

`RequestIDMiddleware` generates a UUID per request and injects it into context and response header (`X-Request-Id`). The logging middleware reads it (`c.Writer.Header().Get("X-Request-Id")`). Every log line contains both `request_id` and `trace_id` — participants can copy either from a log line into Jaeger or grep.

---

## Success Criteria

A participant has completed the quest when:
- [ ] `transfers.yaml` spec is complete with correct request schema and all response codes
- [ ] `POST /v1/transfers` route is registered with JWT + scope middleware in `server.go`
- [ ] Handler parses request, verifies ownership (`sub` == `account.Owner`), calls service, maps all errors with `errors.Is()`
- [ ] OTel span `"transfer.create"` appears in Jaeger with `from_account_id`, `to_account_id`, `amount` attributes
- [ ] Log lines for success and failure use `logger.WithContext(ctx)` with structured fields
- [ ] All 5 table-driven test cases pass (200, 400, 403 wrong owner, 422 insufficient funds, 404 not found)

---

## What This Exercise Teaches

| Concept | Where | Reference |
|---------|-------|-----------|
| OpenAPI spec design | Step 1 — `transfers.yaml` | `docs/openapi/accounts.yaml` |
| Gin route groups + middleware wiring | Step 2 — `server.go` | `account/server.go` in reference |
| `ShouldBindJSON` + error mapping | Step 3 — request parsing | `handlePostAccounts` in reference |
| `errors.Is` switch pattern | Step 3 — error handling | `handlePostAccounts` in reference |
| JWT claims extraction + ownership check | Step 3 — `ClaimsFromCtx` | NEW — not in reference |
| OTel span creation + attributes | Step 3 — `otel.Tracer("bank").Start` | NEW — not in reference |
| `slog.InfoContext(ctx, ...)` structured logging | Step 3 — every log call | slog stdlib — replaces logrus from reference |
| Table-driven handler tests with `httptest` | Step 4 | `handler_test.go` in reference |
| Typed HTTP client with `DoRequest` | Bonus | `pkg/client/rest/account/client.go` in reference |

---

## Bonus Step — Transfer Client (`pkg/client/bank/`)

For participants who finish early. Ports the reference client pattern to `pkg/` for reuse.

### Pre-Built Reference

`pkg/client/bank/client.go` — based directly on `go-training-cba-solution/pkg/client/rest/account/client.go`, with Basic Auth replaced by JWT Bearer, and `...interface{}` replaced by `...any`:

```go
// Reference: go-training-cba-solution/pkg/client/rest/account/client.go
// Changes: Basic Auth → Bearer token, interface{} → any

type Client interface {
    GetToken(ctx context.Context, sub, scope string) (string, error)
    GetAccount(ctx context.Context, id string) (*api.AccountResponse, error)
    Transfer(ctx context.Context, req *api.TransferRequest) (*api.TransferResponse, error)
}

func New(basePath string, httpClient *http.Client) Client {
    return &client{basePath: basePath, HTTPClient: httpClient}
}

type client struct {
    basePath   string
    HTTPClient *http.Client
    token      string  // set by GetToken, sent as Bearer on subsequent calls
}

// GetAccount — pre-built reference (identical pattern to reference solution)
func (c *client) GetAccount(ctx context.Context, id string) (*api.AccountResponse, error) {
    // Reference used: defer logger.TraceDuration(ctx, time.Now(), "GetAccount")
    // slog equivalent — inline defer:
    start := time.Now()
    defer func() {
        slog.DebugContext(ctx, "GetAccount", slog.Duration("duration", time.Since(start)))
    }()

    if err := validator.Default().Var(id, "required"); err != nil {
        return nil, err
    }

    urlPath, err := httppkg.GetURL(c.basePath, "v1/accounts/"+id, "")
    if err != nil {
        return nil, err
    }

    r, err := http.NewRequestWithContext(ctx, http.MethodGet, urlPath, nil)
    if err != nil {
        return nil, err
    }
    r.Header.Set("Authorization", "Bearer "+c.token)

    body, err := httppkg.DoRequest(ctx, c.HTTPClient, r, http.StatusOK)
    if err != nil {
        return nil, err
    }

    var res api.AccountResponse
    return &res, json.Unmarshal(body, &res)
}
```

### Participant Task — Implement `Transfer`

`GetToken` and `GetAccount` are pre-built as reference. Participants implement `Transfer` — the pattern is identical:

```go
func (c *client) Transfer(ctx context.Context, req *api.TransferRequest) (*api.TransferResponse, error) {
    // TODO 1: Validate req — use validator.Default().Struct(req) like CreateAccount in reference

    // TODO 2: Marshal req to JSON — json.Marshal(req)

    // TODO 3: Build URL — httppkg.GetURL(c.basePath, "v1/transfers", "")

    // TODO 4: http.NewRequestWithContext(ctx, http.MethodPost, urlPath, bytes.NewBuffer(jsonPayload))
    //   r.Header.Add(httppkg.HeaderApplicationJSON())
    //   r.Header.Set("Authorization", "Bearer "+c.token)

    // TODO 5: httppkg.DoRequest(ctx, c.HTTPClient, r, http.StatusOK)
    //   Returns *APIError on non-200 — typed, not string

    // TODO 6: json.Unmarshal(body, &res) and return
}
```

### CLI Wiring

After implementing `Transfer`, participants wire `cmd/bank-cli`'s transfer command — skeleton at `internal/bank/cli/transfer/transfer.go`, using `internal/bank/cli/account/balance.go` as reference.

---

## Out of Scope

- Implementing repository or service layer (pre-built)
- Database migrations (pre-run via `make migrate`)
- OTel SDK bootstrapping (pre-built in `cmd/bank-api/main.go`)
- JWT secret management / production auth (hardcoded in config)
- Client retry logic / exponential backoff
