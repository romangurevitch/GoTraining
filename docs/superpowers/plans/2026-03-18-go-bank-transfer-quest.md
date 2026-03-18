# Go Bank Transfer Quest — Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the pre-scaffolded Go Bank Transfer Quest exercise — a fully working bank API (accounts handler as reference, transfers as participant skeleton) with OTel tracing, slog logging, JWT auth, and Jaeger via Docker Compose.

**Architecture:** Restructure `GoTraining/` from flat `api/handlers.go` + in-memory store to one-package-per-resource API layout (`api/account/`, `api/transfer/`), add middleware layer, replace store with Postgres repository using go-jet, add OTel + slog throughout. Account handler is the complete working reference; transfer handler is a skeleton with 5 guided TODOs. All pre-built code is derived directly from `go-training-cba-solution/`.

**Tech Stack:** Go 1.26.1, Gin v1.12.0, go-jet v2.14.1, log/slog (stdlib), slog-otel, slog-gin, OTel SDK, OTLP/HTTP → Jaeger, golang-jwt/jwt/v5, go-playground/validator/v10, mockery (go.uber.org/mock), httptest, Docker Compose.

---

## File Map

### Files to Delete / Replace
- `internal/bank/store/memory.go` → **deleted** (replaced by Postgres repository)
- `internal/bank/store/store.go` → **deleted** (replaced by `repository/repository.go`)
- `internal/bank/api/handlers.go` → **deleted** (split into `api/account/` and `api/transfer/`)
- `internal/bank/api/middleware.go` → **deleted** (replaced by `api/middleware/` package)
- `internal/bank/api/server.go` → **replaced** (new structure)
- `internal/bank/client/client.go` → **replaced by** `pkg/client/bank/client.go` (bonus)
- `cmd/bank-api/main.go` → **replaced** (OTel bootstrap + slog setup)
- `internal/pkg/logger/` → **deleted** (slog replaces it; note: may not exist in `GoTraining/`, check)

### Files to Create
```
GoTraining/
├── pkg/
│   ├── api/error/error.go              # APIError type + NewAPIError (from reference, slog)
│   └── http/http.go                    # DoRequest, GetURL, HeaderApplicationJSON (from reference, ctx fix)
├── internal/bank/
│   ├── repository/
│   │   ├── repository.go               # Repository interface (replaces store.Store)
│   │   └── postgres/
│   │       └── repository.go           # go-jet Postgres implementation
│   ├── service/
│   │   └── service.go                  # MODIFIED: add Service interface + Transfer method
│   ├── api/
│   │   ├── server.go                   # REPLACED: new structure with middleware + route groups
│   │   ├── middleware/
│   │   │   ├── logging.go              # slog-gin logging middleware
│   │   │   ├── tracing.go              # OTel span-per-request middleware
│   │   │   ├── requestid.go            # UUID request ID middleware
│   │   │   └── auth.go                 # JWT middleware + RequireScope + ClaimsFromCtx
│   │   ├── auth/
│   │   │   └── handler.go              # POST /v1/token — issues JWT (pre-built)
│   │   ├── account/
│   │   │   ├── handler.go              # REFERENCE — fully implemented
│   │   │   ├── handler_test.go         # REFERENCE tests — fully implemented
│   │   │   └── types.go                # AccountResponse, CreateAccountRequest
│   │   └── transfer/
│   │       ├── handler.go              # SKELETON — participants implement
│   │       ├── handler_test.go         # SKELETON — participants complete
│   │       └── types.go                # TransferRequest, TransferResponse
│   └── domain/
│       └── transfer.go                 # NEW: Transfer domain type (pre-built)
├── docs/openapi/
│   ├── accounts.yaml                   # Complete reference spec
│   └── transfers.yaml                  # Partially filled — participants complete
├── cmd/bank-api/main.go                # REPLACED: OTel bootstrap + slog setup
└── docker-compose.yaml                 # MODIFIED: add Jaeger service
```

---

## Task 1: Add Dependencies to go.mod

**Files:**
- Modify: `go.mod`

- [ ] **Step 1: Add new dependencies**

```bash
cd /Users/romang/workspace/CBA/GoTraining
go get github.com/golang-jwt/jwt/v5
go get github.com/remychantenay/slog-otel
go get github.com/samber/slog-gin
go get go.opentelemetry.io/otel
go get go.opentelemetry.io/otel/trace
go get go.opentelemetry.io/otel/sdk/trace
go get go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp
go get go.opentelemetry.io/contrib/bridges/otelslog
go get go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin
go get github.com/google/uuid
go get github.com/lib/pq      # Postgres driver (already in go.mod as indirect — make it direct)
go get go.uber.org/mock       # mockery mock generation + gomock test assertions
```

- [ ] **Step 2: Tidy dependencies**

```bash
go mod tidy
```

Expected: `go.mod` updated with new requires, no errors.

- [ ] **Step 3: Verify build still compiles before any changes**

```bash
go build ./...
```

Expected: clean build (existing code compiles before restructure).

- [ ] **Step 4: Commit**

```bash
git add go.mod go.sum
git commit -m "feat: add OTel, slog-gin, slog-otel, JWT dependencies"
```

---

## Task 2: Create `pkg/api/error/error.go`

**Files:**
- Create: `pkg/api/error/error.go`

Source: `go-training-cba-solution/pkg/api/error/error.go`
Change: `logger.WithContext(ctx).WithError(err)...` → `slog.ErrorContext(ctx, ...)`

- [ ] **Step 1: Create directory and file**

Create `pkg/api/error/error.go`:

```go
package error

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

// ErrInternalServerError is the sentinel for unexpected server errors.
var ErrInternalServerError = errors.New("internal server error")

// NewAPIError logs the error with context (trace_id/span_id injected automatically by OtelHandler)
// and returns a tuple suitable for c.JSON(apierror.NewAPIError(...)).
func NewAPIError(ctx context.Context, status int, msg string, err error) (int, *APIError) {
	slog.ErrorContext(ctx, msg,
		slog.Int("status", status),
		slog.Any("error", err),
	)
	return status, &APIError{Message: msg}
}

// NewUnauthorizedError returns a 401 without logging (avoids noise from probing).
func NewUnauthorizedError() (int, *APIError) {
	return http.StatusUnauthorized, &APIError{Message: "unauthorized"}
}

// APIError is the JSON error body returned by all API endpoints.
type APIError struct {
	Message string `json:"message"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("api error: %v", e.Message)
}
```

- [ ] **Step 2: Build to verify no syntax errors**

```bash
go build ./pkg/...
```

Expected: compiles cleanly.

- [ ] **Step 3: Commit**

```bash
git add pkg/api/error/error.go
git commit -m "feat: add pkg/api/error — APIError + NewAPIError from reference solution"
```

---

## Task 3: Create `pkg/http/http.go`

**Files:**
- Create: `pkg/http/http.go`

Source: `go-training-cba-solution/pkg/http/http.go`
Fix: `_ context.Context` → `ctx` actually passed to `r.WithContext(ctx)`. `...interface{}` → `...any`.

- [ ] **Step 1: Create file**

Create `pkg/http/http.go`:

```go
package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	apierror "github.com/romangurevitch/go-training/pkg/api/error"
)

// DoRequest executes the request with the given context, checks the response status,
// and returns the body bytes. Returns *apierror.APIError on non-expected status codes.
//
// Source: go-training-cba-solution/pkg/http/http.go
// Fix: context is now passed to r.WithContext(ctx) (was ignored in reference)
func DoRequest(ctx context.Context, client *http.Client, r *http.Request, expectedResponses ...int) ([]byte, error) {
	resp, err := client.Do(r.WithContext(ctx))
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

// GetURL builds a full URL by joining baseURL with path p and appending a formatted query string.
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

// HeaderApplicationJSON returns the Content-Type header key and value for JSON.
func HeaderApplicationJSON() (key, value string) {
	return "Content-Type", "application/json"
}
```

- [ ] **Step 2: Build to verify**

```bash
go build ./pkg/...
```

Expected: compiles cleanly.

- [ ] **Step 3: Commit**

```bash
git add pkg/http/http.go
git commit -m "feat: add pkg/http — DoRequest + GetURL from reference solution, fix ctx usage"
```

---

## Task 4: Verify / Update Domain Package

**Files:**
- Verify: `internal/bank/domain/account.go` (already exists — check sentinel errors and types)
- Verify: `internal/bank/domain/transaction.go` (already exists — check TypeDeposit/TypeWithdrawal)
- Create: `internal/bank/domain/transfer.go` (new — Transfer domain type)

The domain package already exists in the current `GoTraining/`. All subsequent tasks depend on it compiling correctly.

- [ ] **Step 1: Verify `domain/account.go` has required sentinel errors and types**

Check that the file contains:
- `ErrAccountNotFound`, `ErrInsufficientFunds`, `ErrAccountLocked`, `ErrInvalidAmount`, `ErrAccountAlreadyExists`
- `AccountStatus` with `StatusOpen`, `StatusLocked`, `StatusClosed`
- `Account` struct with `ID`, `Owner`, `Balance`, `Status`, `CreatedAt`, `UpdatedAt` fields
- `CanPerformTransaction() error` method on `*Account`

```bash
cat internal/bank/domain/account.go
```

If anything is missing, add it now. The file as read in codebase exploration has all of these — no changes expected.

- [ ] **Step 2: Verify `domain/transaction.go` has required types**

Check that the file contains:
- `TransactionType` with `TypeDeposit = "DEPOSIT"` and `TypeWithdrawal = "WITHDRAWAL"`
- `Transaction` struct with `ID`, `AccountID`, `Amount`, `Type`, `CreatedAt` fields

```bash
cat internal/bank/domain/transaction.go
```

No changes expected — file already correct.

- [ ] **Step 3: Create `internal/bank/domain/transfer.go`**

The `Transfer` service method uses `domain.Transaction` entries for debit/credit. No separate `Transfer` domain type is strictly required for the exercise, but add a thin type for the API layer to reference:

```go
package domain

// Transfer represents a completed fund movement between two accounts.
// Used for logging and response mapping — not persisted as its own record.
// The underlying transactions are recorded as TypeDeposit and TypeWithdrawal entries.
type Transfer struct {
	FromAccountID string
	ToAccountID   string
	Amount        float64
}
```

- [ ] **Step 4: Build domain package**

```bash
go build ./internal/bank/domain/...
```

Expected: compiles cleanly.

- [ ] **Step 5: Commit**

```bash
git add internal/bank/domain/
git commit -m "feat: verify domain package — add transfer.go stub"
```

---

## Task 5: Create Repository Layer (replaces `store/`)

**Files:**
- Create: `internal/bank/repository/repository.go`
- Create: `internal/bank/repository/postgres/repository.go`

- [ ] **Step 1: Create `repository/repository.go`**

```go
package repository

import (
	"context"

	"github.com/romangurevitch/go-training/internal/bank/domain"
)

// Repository defines the data access contract for the bank.
// Replaces store.Store — clearer intent, Postgres-only.
type Repository interface {
	GetAccount(ctx context.Context, id string) (*domain.Account, error)
	SaveAccount(ctx context.Context, account *domain.Account) error
	ListTransactions(ctx context.Context, accountID string) ([]domain.Transaction, error)
	SaveTransaction(ctx context.Context, transaction *domain.Transaction) error
}
```

- [ ] **Step 2: Create `repository/postgres/repository.go`**

```go
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/romangurevitch/go-training/internal/bank/domain"
	"github.com/romangurevitch/go-training/internal/bank/repository"
)

// Ensure PostgresRepository implements Repository at compile time.
var _ repository.Repository = (*PostgresRepository)(nil)

// PostgresRepository implements Repository using go-jet for type-safe SQL.
// go-jet generated models live in internal/bank/repository/postgres/gen/ (pre-generated).
type PostgresRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetAccount(ctx context.Context, id string) (*domain.Account, error) {
	// TODO: implement with go-jet — pre-built for participants.
	// Pattern: SELECT FROM accounts WHERE id = ?
	// Map to domain.Account, return domain.ErrAccountNotFound if sql.ErrNoRows.
	return nil, fmt.Errorf("not implemented")
}

func (r *PostgresRepository) SaveAccount(ctx context.Context, account *domain.Account) error {
	// TODO: implement with go-jet — pre-built for participants.
	return fmt.Errorf("not implemented")
}

func (r *PostgresRepository) ListTransactions(ctx context.Context, accountID string) ([]domain.Transaction, error) {
	// TODO: implement with go-jet — pre-built for participants.
	return nil, fmt.Errorf("not implemented")
}

func (r *PostgresRepository) SaveTransaction(ctx context.Context, t *domain.Transaction) error {
	// TODO: implement with go-jet — pre-built for participants.
	return fmt.Errorf("not implemented")
}

// isNotFound maps sql.ErrNoRows to domain.ErrAccountNotFound.
func isNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
```

Note: The go-jet generated models will be added in Step 3 below via the `jet` CLI.

- [ ] **Step 3: Run go-jet code generation to produce typed model structs**

Ensure the Postgres DB is running first (`docker compose up postgres -d`), then:

```bash
# Install jet CLI if not present
go install github.com/go-jet/jet/v2/cmd/jet@latest

# Generate typed Go models from the live schema into gen/
jet -dsn="postgres://gotrainer:verysecret@localhost:5432/gobank?sslmode=disable" \
    -schema=public \
    -path=./internal/bank/repository/postgres/gen
```

Expected: `internal/bank/repository/postgres/gen/` created with Go files for each table (e.g. `gen/public/model/accounts.go`, `gen/public/table/accounts.go`).

- [ ] **Step 4: Implement the Postgres repository using go-jet**

Replace the stub implementations in `repository/postgres/repository.go` with real go-jet queries. Source: `go-training-cba-solution/internal/repository/postgres/` for the query patterns.

```go
// GetAccount implementation pattern using go-jet:
import (
    "database/sql"
    "errors"
    . "github.com/romangurevitch/go-training/internal/bank/repository/postgres/gen/public/table"
    "github.com/romangurevitch/go-training/internal/bank/repository/postgres/gen/public/model"
    "github.com/go-jet/jet/v2/postgres"
)

func (r *PostgresRepository) GetAccount(ctx context.Context, id string) (*domain.Account, error) {
    stmt := Accounts.SELECT(Accounts.AllColumns).
        WHERE(Accounts.ID.EQ(postgres.String(id)))

    var dest model.Accounts
    if err := stmt.QueryContext(ctx, r.db, &dest); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, domain.ErrAccountNotFound
        }
        return nil, err
    }
    return toDomainAccount(&dest), nil
}

func toDomainAccount(m *model.Accounts) *domain.Account {
    return &domain.Account{
        ID:        m.ID,
        Owner:     m.Owner,
        Balance:   m.Balance,
        Status:    domain.AccountStatus(m.Status),
        CreatedAt: m.CreatedAt,
        UpdatedAt: m.UpdatedAt,
    }
}
```

Implement `SaveAccount`, `ListTransactions`, `SaveTransaction` with the same go-jet pattern. For `SaveAccount` use `Accounts.INSERT(...).ON_CONFLICT(Accounts.ID).DO_UPDATE(...)` or an upsert. For transactions use `Transactions.INSERT(...)`.

- [ ] **Step 5: Build to verify**

```bash
go build ./internal/bank/repository/...
```

Expected: compiles cleanly.

- [ ] **Step 6: Commit**

```bash
git add internal/bank/repository/
git commit -m "feat: add repository layer — Repository interface + go-jet Postgres implementation"
```

---

## Task 6: Update Service Layer — Add Interface + Transfer Method

**Files:**
- Modify: `internal/bank/service/service.go`

The current `service.go` uses `store.Store`. We:
1. Add a `Service` interface (enables mock injection in handler tests)
2. Add `Transfer` method to `BankService`
3. Change the `store` field to use `repository.Repository`
4. Replace `logger.L()` with `slog.InfoContext`

- [ ] **Step 1: Rewrite `service/service.go`**

```go
package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/romangurevitch/go-training/internal/bank/domain"
	"github.com/romangurevitch/go-training/internal/bank/repository"
)

// Service is the business logic interface. Enables mock injection in handler tests.
// Replaces direct BankService usage — mirrors the reference solution's service abstraction.
type Service interface {
	CreateAccount(ctx context.Context, owner string) (*domain.Account, error)
	GetAccount(ctx context.Context, id string) (*domain.Account, error)
	Deposit(ctx context.Context, accountID string, amount float64) error
	Withdraw(ctx context.Context, accountID string, amount float64) error
	Transfer(ctx context.Context, fromID, toID string, amount float64) error
}

// BankService implements Service backed by a Repository.
// Source: existing service.go, updated:
//   - store.Store → repository.Repository
//   - logger.L().Info → slog.InfoContext
//   - Add Transfer method
type BankService struct {
	repo repository.Repository
}

// Ensure BankService implements Service at compile time.
var _ Service = (*BankService)(nil)

func NewBankService(repo repository.Repository) *BankService {
	return &BankService{repo: repo}
}

func (s *BankService) CreateAccount(ctx context.Context, owner string) (*domain.Account, error) {
	acc := &domain.Account{
		ID:        fmt.Sprintf("ACC-%d", time.Now().UnixNano()),
		Owner:     owner,
		Balance:   0,
		Status:    domain.StatusOpen,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.SaveAccount(ctx, acc); err != nil {
		return nil, fmt.Errorf("failed to save account: %w", err)
	}

	slog.InfoContext(ctx, "account created", slog.String("id", acc.ID), slog.String("owner", acc.Owner))
	return acc, nil
}

func (s *BankService) GetAccount(ctx context.Context, id string) (*domain.Account, error) {
	return s.repo.GetAccount(ctx, id)
}

func (s *BankService) Deposit(ctx context.Context, accountID string, amount float64) error {
	if amount <= 0 {
		return domain.ErrInvalidAmount
	}

	acc, err := s.repo.GetAccount(ctx, accountID)
	if err != nil {
		return err
	}

	if err := acc.CanPerformTransaction(); err != nil {
		return err
	}

	acc.Balance += amount
	acc.UpdatedAt = time.Now()

	if err := s.repo.SaveAccount(ctx, acc); err != nil {
		return err
	}

	t := &domain.Transaction{
		ID:        fmt.Sprintf("TRX-%d", time.Now().UnixNano()),
		AccountID: accountID,
		Amount:    amount,
		Type:      domain.TypeDeposit,
		CreatedAt: time.Now(),
	}

	return s.repo.SaveTransaction(ctx, t)
}

func (s *BankService) Withdraw(ctx context.Context, accountID string, amount float64) error {
	if amount <= 0 {
		return domain.ErrInvalidAmount
	}

	acc, err := s.repo.GetAccount(ctx, accountID)
	if err != nil {
		return err
	}

	if err := acc.CanPerformTransaction(); err != nil {
		return err
	}

	if acc.Balance < amount {
		return domain.ErrInsufficientFunds
	}

	acc.Balance -= amount
	acc.UpdatedAt = time.Now()

	if err := s.repo.SaveAccount(ctx, acc); err != nil {
		return err
	}

	t := &domain.Transaction{
		ID:        fmt.Sprintf("TRX-%d", time.Now().UnixNano()),
		AccountID: accountID,
		Amount:    amount,
		Type:      domain.TypeWithdrawal,
		CreatedAt: time.Now(),
	}

	return s.repo.SaveTransaction(ctx, t)
}

// Transfer moves funds from one account to another atomically.
// Pre-built for participants — they call this from the transfer handler.
func (s *BankService) Transfer(ctx context.Context, fromID, toID string, amount float64) error {
	if amount <= 0 {
		return domain.ErrInvalidAmount
	}

	from, err := s.repo.GetAccount(ctx, fromID)
	if err != nil {
		return err
	}

	to, err := s.repo.GetAccount(ctx, toID)
	if err != nil {
		return err
	}

	if err := from.CanPerformTransaction(); err != nil {
		return err
	}

	if err := to.CanPerformTransaction(); err != nil {
		return err
	}

	if from.Balance < amount {
		return domain.ErrInsufficientFunds
	}

	from.Balance -= amount
	from.UpdatedAt = time.Now()

	to.Balance += amount
	to.UpdatedAt = time.Now()

	if err := s.repo.SaveAccount(ctx, from); err != nil {
		return fmt.Errorf("failed to debit source account: %w", err)
	}

	if err := s.repo.SaveAccount(ctx, to); err != nil {
		return fmt.Errorf("failed to credit destination account: %w", err)
	}

	// Record debit transaction on source account
	debit := &domain.Transaction{
		ID:        fmt.Sprintf("TRX-%d-D", time.Now().UnixNano()),
		AccountID: fromID,
		Amount:    amount,
		Type:      domain.TypeWithdrawal,
		CreatedAt: time.Now(),
	}
	if err := s.repo.SaveTransaction(ctx, debit); err != nil {
		return err
	}

	// Record credit transaction on destination account
	credit := &domain.Transaction{
		ID:        fmt.Sprintf("TRX-%d-C", time.Now().UnixNano()),
		AccountID: toID,
		Amount:    amount,
		Type:      domain.TypeDeposit,
		CreatedAt: time.Now(),
	}

	slog.InfoContext(ctx, "transfer completed",
		slog.String("from_account_id", fromID),
		slog.String("to_account_id", toID),
		slog.Float64("amount", amount),
	)

	return s.repo.SaveTransaction(ctx, credit)
}

```

- [ ] **Step 2: Generate Service mock using mockery**

```bash
cd internal/bank/service
# Add go:generate directive to service.go above the Service interface:
# //go:generate go run go.uber.org/mock/mockgen -source=service.go -destination=mocks/mock_service.go -package=mocks
go generate ./...
```

Or run mockgen directly:

```bash
go run go.uber.org/mock/mockgen \
  -source=internal/bank/service/service.go \
  -destination=internal/bank/service/mocks/mock_service.go \
  -package=mocks
```

Expected: `internal/bank/service/mocks/mock_service.go` created with `MockService` struct.

- [ ] **Step 3: Build to verify**

```bash
go build ./internal/bank/service/...
```

Expected: compiles cleanly.

- [ ] **Step 4: Commit**

```bash
git add internal/bank/service/
git commit -m "feat: add Service interface + Transfer method to BankService, switch to repository.Repository"
```

---

## Task 7: Create Middleware Package

**Files:**
- Create: `internal/bank/api/middleware/requestid.go`
- Create: `internal/bank/api/middleware/tracing.go`
- Create: `internal/bank/api/middleware/logging.go`
- Create: `internal/bank/api/middleware/auth.go`

### Step 6a: `requestid.go`

- [ ] **Step 1: Create `middleware/requestid.go`**

```go
package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type contextKey string

const requestIDKey contextKey = "request_id"

const headerRequestID = "X-Request-Id"

// RequestIDMiddleware generates a UUID per request, injects it into context
// and the response header. slog-gin reads it automatically when WithRequestID: true.
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(headerRequestID)
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Writer.Header().Set(headerRequestID, requestID)
		c.Set(string(requestIDKey), requestID)
		c.Next()
	}
}

// RequestIDFromCtx extracts the request ID from the Gin context.
func RequestIDFromCtx(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}
```

### Step 6b: `tracing.go`

- [ ] **Step 2: Create `middleware/tracing.go`**

```go
package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// TracingMiddleware starts an OTel span for every incoming HTTP request.
// Span name: "{method} {path}" — e.g. "POST /v1/transfers"
// Uses otelgin which handles W3C trace context propagation automatically.
func TracingMiddleware(serviceName string) gin.HandlerFunc {
	return otelgin.Middleware(serviceName)
}
```

### Step 6c: `logging.go`

- [ ] **Step 3: Create `middleware/logging.go`**

```go
package middleware

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

// LoggingMiddleware returns a Gin middleware that logs one structured line per request.
// Fields logged: time, method, path, status, latency, trace_id, span_id, request_id.
// trace_id/span_id are extracted from the active OTel span in the request context.
//
// Source: replaces hand-rolled JSONLogMiddleware from go-training-cba-solution/internal/server/rest/middleware/logger.go
// Upgrade: uses slog-gin instead of logrus — zero-boilerplate OTel trace correlation.
func LoggingMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return sloggin.NewWithConfig(logger, sloggin.Config{
		DefaultLevel:     slog.LevelInfo,
		ClientErrorLevel: slog.LevelWarn,
		ServerErrorLevel: slog.LevelError,
		WithTraceID:      true,
		WithSpanID:       true,
		WithRequestID:    true,
	})
}
```

### Step 6d: `auth.go`

- [ ] **Step 4: Create `middleware/auth.go`**

```go
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	apierror "github.com/romangurevitch/go-training/pkg/api/error"
)

type claimsKey struct{}

// Claims extends jwt.RegisteredClaims with a Scope field.
// JWT payload example: { "sub": "alice", "scope": "accounts:read transfers:write", "exp": ... }
//
// Source: replaces gin.BasicAuth from go-training-cba-solution/internal/server/rest/middleware/auth.go
// Upgrade: JWT Bearer with scope-based authorization + sub claim for ownership checks.
type Claims struct {
	Scope string `json:"scope"`
	jwt.RegisteredClaims
}

// JWTMiddleware validates the Bearer token from Authorization header.
// On success: injects *Claims into Gin context under claimsKey{}.
// On failure: returns 401 and aborts — downstream handlers do not run.
func JWTMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.JSON(apierror.NewUnauthorizedError())
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(apierror.NewUnauthorizedError())
			c.Abort()
			return
		}

		// Store claims on request context so ClaimsFromCtx(ctx context.Context) works
		// in handlers and tests that receive ctx from c.Request.Context().
		ctx := context.WithValue(c.Request.Context(), claimsKey{}, claims)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// RequireScope checks that the injected Claims contain the required scope.
// Returns 403 if scope is missing — used as per-route middleware after JWTMiddleware.
//
// Usage: accounts.GET("/:id", middleware.RequireScope("accounts:read"), handler.GetAccount)
func RequireScope(scope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := ClaimsFromCtx(c.Request.Context())
		if claims == nil || !strings.Contains(claims.Scope, scope) {
			c.JSON(http.StatusForbidden, &apierror.APIError{Message: "insufficient scope"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// ClaimsFromCtx extracts the JWT claims injected by JWTMiddleware from the request context.
// Returns nil if not present (should not happen in protected routes).
func ClaimsFromCtx(ctx context.Context) *Claims {
	if c, ok := ctx.Value(claimsKey{}).(*Claims); ok {
		return c
	}
	return nil
}
```

- [ ] **Step 5: Build middleware package**

```bash
go build ./internal/bank/api/middleware/...
```

Expected: compiles cleanly.

- [ ] **Step 6: Commit**

```bash
git add internal/bank/api/middleware/
git commit -m "feat: add middleware — RequestID, Tracing (otelgin), Logging (slog-gin), JWT auth"
```

---

## Task 8: Create Auth Handler (`POST /v1/token`)

**Files:**
- Create: `internal/bank/api/auth/handler.go`

- [ ] **Step 1: Create `api/auth/handler.go`**

```go
package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/romangurevitch/go-training/internal/bank/api/middleware"
)

// Handler issues JWT tokens for the /v1/token endpoint.
// Pre-built — participants use curl to get tokens for manual testing.
type Handler struct {
	secret string
}

func New(secret string) *Handler {
	return &Handler{secret: secret}
}

// IssueTokenRequest is the request body for POST /v1/token.
type IssueTokenRequest struct {
	Sub   string `json:"sub"   binding:"required"` // account owner — becomes JWT Subject
	Scope string `json:"scope" binding:"required"` // space-separated scopes
}

// IssueTokenResponse wraps the signed token.
type IssueTokenResponse struct {
	Token string `json:"token"`
}

// IssueToken issues a signed JWT for the given sub and scope.
// No authentication required — this is a training tool, not production auth.
//
// Example request:
//   curl -X POST localhost:8080/v1/token \
//     -H 'Content-Type: application/json' \
//     -d '{"sub":"alice","scope":"accounts:read transfers:write"}'
func (h *Handler) IssueToken(c *gin.Context) {
	var req IssueTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request: " + err.Error()})
		return
	}

	claims := middleware.Claims{
		Scope: req.Scope,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   req.Sub,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(h.secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to sign token"})
		return
	}

	c.JSON(http.StatusOK, IssueTokenResponse{Token: signed})
}
```

- [ ] **Step 2: Build**

```bash
go build ./internal/bank/api/auth/...
```

Expected: compiles cleanly.

- [ ] **Step 3: Commit**

```bash
git add internal/bank/api/auth/
git commit -m "feat: add auth handler — POST /v1/token issues JWT for training use"
```

---

## Task 9: Create Account Handler (Reference Implementation)

**Files:**
- Create: `internal/bank/api/account/types.go`
- Create: `internal/bank/api/account/handler.go`
- Create: `internal/bank/api/account/handler_test.go`

This is the **fully working reference** participants read before implementing the transfer handler.

### Step 9a: Types

- [ ] **Step 1: Create `api/account/types.go`**

```go
package account

import (
	"time"

	"github.com/romangurevitch/go-training/internal/bank/domain"
)

// CreateAccountRequest is the JSON body for POST /v1/accounts.
type CreateAccountRequest struct {
	Owner string `json:"owner" binding:"required"`
}

// AccountResponse is the JSON body returned by GET /v1/accounts/:id and POST /v1/accounts.
type AccountResponse struct {
	ID        string    `json:"id"`
	Owner     string    `json:"owner"`
	Balance   float64   `json:"balance"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toAccountResponse(a *domain.Account) AccountResponse {
	return AccountResponse{
		ID:        a.ID,
		Owner:     a.Owner,
		Balance:   a.Balance,
		Status:    string(a.Status),
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}
```

### Step 9b: Handler

- [ ] **Step 2: Create `api/account/handler.go`**

```go
package account

// REFERENCE IMPLEMENTATION — participants read this before implementing api/transfer/handler.go.
//
// Source: go-training-cba-solution/internal/server/rest/account/handler.go
// Changes vs reference:
//   - logrus logger.WithContext(ctx).WithField(...).Info() → slog.InfoContext(ctx, ...)
//   - OTel span added per handler (reference had no per-handler tracing)
//   - Direct store access → service.Service interface

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/romangurevitch/go-training/internal/bank/domain"
	"github.com/romangurevitch/go-training/internal/bank/service"
	apierror "github.com/romangurevitch/go-training/pkg/api/error"
)

// Handler handles account-related HTTP requests.
// Identical pattern to go-training-cba-solution accountServer — just wired to service.Service.
type Handler struct {
	svc service.Service
}

func New(svc service.Service) *Handler {
	return &Handler{svc: svc}
}

// GetAccount handles GET /v1/accounts/:id
//
// Pattern demonstrated (in order):
//   1. Extract ctx from request
//   2. Start OTel span — defer span.End()
//   3. Read URL param
//   4. slog.InfoContext — structured log with context (trace_id auto-injected)
//   5. Call service method
//   6. Map errors with errors.Is — return apierror.NewAPIError tuple to c.JSON
//   7. On success: set span attribute + return 200
func (h *Handler) GetAccount(c *gin.Context) {
	ctx := c.Request.Context()

	ctx, span := otel.Tracer("bank").Start(ctx, "account.get")
	defer span.End()

	id := c.Param("id")
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

// CreateAccount handles POST /v1/accounts
//
// Pattern demonstrated (in addition to GetAccount):
//   1. ShouldBindJSON — bind + validate request body
//   2. On bind error: 400 with apierror
//   3. Call service, map domain errors
//   4. On success: log created entity + return 201
func (h *Handler) CreateAccount(c *gin.Context) {
	ctx := c.Request.Context()

	ctx, span := otel.Tracer("bank").Start(ctx, "account.create")
	defer span.End()

	var req CreateAccountRequest
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
		span.SetAttributes(attribute.String("account.id", result.ID))
		slog.InfoContext(ctx, "account created", slog.String("account_id", result.ID))
		c.JSON(http.StatusCreated, toAccountResponse(result))
	}
}
```

### Step 9c: Handler Tests

- [ ] **Step 3: Create `api/account/handler_test.go`**

```go
package account_test

// REFERENCE TEST — participants replicate this pattern for api/transfer/handler_test.go.
//
// Source: go-training-cba-solution/internal/server/rest/account/handler_test.go
// Patterns demonstrated:
//   - Table-driven tests with t.Run
//   - type fields struct { svc func(t *testing.T) service.Service }
//   - mockery-generated MockService: m.EXPECT().GetAccount(...).Return(...).Times(1)
//   - gin.SetMode(gin.TestMode) + httptest.NewRecorder()
//   - JWT token in Authorization header (replaces Basic Auth from reference)

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/romangurevitch/go-training/internal/bank/api/account"
	"github.com/romangurevitch/go-training/internal/bank/api/middleware"
	"github.com/romangurevitch/go-training/internal/bank/domain"
	"github.com/romangurevitch/go-training/internal/bank/service/mocks"
)

const testSecret = "test-secret"

// testToken issues a signed JWT for use in test Authorization headers.
// Replaces the `const auth = "Basic Z29pczp0aGViZXN0"` from the reference test.
func testToken(t *testing.T, sub, scope string) string {
	t.Helper()
	claims := middleware.Claims{
		Scope: scope,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   sub,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(testSecret))
	require.NoError(t, err)
	return signed
}

// setupRouter builds a minimal Gin engine with the account routes wired up.
// Mirrors NewServer() but scoped to account routes only — keeps tests isolated.
func setupRouter(svc *mocks.MockService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := account.New(svc)

	v1 := r.Group("/v1/accounts")
	v1.Use(middleware.JWTMiddleware(testSecret))
	{
		v1.GET("/:id", middleware.RequireScope("accounts:read"), h.GetAccount)
		v1.POST("", middleware.RequireScope("accounts:write"), h.CreateAccount)
	}
	return r
}

var testAccount = &domain.Account{
	ID:      "ACC-001",
	Owner:   "alice",
	Balance: 100.0,
	Status:  domain.StatusOpen,
}

func TestGetAccount(t *testing.T) {
	type fields struct {
		svc func(t *testing.T) *mocks.MockService
	}
	tests := []struct {
		name     string
		fields   fields
		id       string
		scope    string
		wantCode int
		wantBody any
	}{
		{
			name: "success — returns 200 with account",
			fields: fields{
				svc: func(t *testing.T) *mocks.MockService {
					m := mocks.NewMockService(gomock.NewController(t))
					m.EXPECT().GetAccount(gomock.Any(), "ACC-001").Return(testAccount, nil).Times(1)
					return m
				},
			},
			id:       "ACC-001",
			scope:    "accounts:read",
			wantCode: http.StatusOK,
			wantBody: map[string]any{"id": "ACC-001", "owner": "alice"},
		},
		{
			name: "not found — returns 404",
			fields: fields{
				svc: func(t *testing.T) *mocks.MockService {
					m := mocks.NewMockService(gomock.NewController(t))
					m.EXPECT().GetAccount(gomock.Any(), "MISSING").Return(nil, domain.ErrAccountNotFound).Times(1)
					return m
				},
			},
			id:       "MISSING",
			scope:    "accounts:read",
			wantCode: http.StatusNotFound,
		},
		{
			name: "wrong scope — returns 403",
			fields: fields{
				svc: func(t *testing.T) *mocks.MockService {
					return mocks.NewMockService(gomock.NewController(t)) // no calls expected
				},
			},
			id:       "ACC-001",
			scope:    "transfers:write", // missing accounts:read
			wantCode: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter(tt.fields.svc(t))
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/v1/accounts/"+tt.id, nil)
			req.Header.Set("Authorization", "Bearer "+testToken(t, "alice", tt.scope))
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)
			if tt.wantBody != nil {
				var got map[string]any
				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
				for k, v := range tt.wantBody.(map[string]any) {
					assert.Equal(t, v, got[k])
				}
			}
		})
	}
}

func TestCreateAccount(t *testing.T) {
	type fields struct {
		svc func(t *testing.T) *mocks.MockService
	}
	tests := []struct {
		name     string
		fields   fields
		body     any
		scope    string
		wantCode int
	}{
		{
			name: "success — returns 201",
			fields: fields{
				svc: func(t *testing.T) *mocks.MockService {
					m := mocks.NewMockService(gomock.NewController(t))
					m.EXPECT().CreateAccount(gomock.Any(), "alice").Return(testAccount, nil).Times(1)
					return m
				},
			},
			body:     map[string]string{"owner": "alice"},
			scope:    "accounts:write",
			wantCode: http.StatusCreated,
		},
		{
			name: "missing owner — returns 400",
			fields: fields{
				svc: func(t *testing.T) *mocks.MockService {
					return mocks.NewMockService(gomock.NewController(t))
				},
			},
			body:     map[string]string{},
			scope:    "accounts:write",
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter(tt.fields.svc(t))
			bodyBytes, _ := json.Marshal(tt.body)
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/v1/accounts", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+testToken(t, "alice", tt.scope))
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)
		})
	}
}
```

- [ ] **Step 4: Run account handler tests**

```bash
go test ./internal/bank/api/account/... -v
```

Expected: all tests PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/bank/api/account/
git commit -m "feat: add account handler — reference implementation with OTel + slog + JWT tests"
```

---

## Task 10: Create Transfer Handler Skeleton (Participant Quest)

**Files:**
- Create: `internal/bank/api/transfer/types.go`
- Create: `internal/bank/api/transfer/handler.go`
- Create: `internal/bank/api/transfer/handler_test.go`

This is what participants receive. The handler has 5 guided TODOs; the test has 2 pre-written cases and 3 TODOs.

### Step 10a: Types

- [ ] **Step 1: Create `api/transfer/types.go`**

```go
package transfer

// CreateTransferRequest is the JSON body for POST /v1/transfers.
// TODO (Step 1 — OpenAPI): You defined these fields in transfers.yaml — use the same names here.
type CreateTransferRequest struct {
	FromAccountID string  `json:"from_account_id" binding:"required"`
	ToAccountID   string  `json:"to_account_id"   binding:"required"`
	Amount        float64 `json:"amount"           binding:"required,gt=0"`
}

// TransferResponse is the JSON body returned on successful transfer.
// TODO (Step 1 — OpenAPI): Matches the 200 response you defined in transfers.yaml.
type TransferResponse struct {
	Status string `json:"status"`
}
```

### Step 10b: Handler Skeleton

- [ ] **Step 2: Create `api/transfer/handler.go`**

```go
package transfer

// PARTICIPANT QUEST — implement this handler.
//
// Before starting: read api/account/handler.go in full.
// Every TODO below points to the exact line in account/handler.go that demonstrates the pattern.
//
// Time estimate: 60-80 minutes for Step 3 of the quest.

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/romangurevitch/go-training/internal/bank/domain"
	"github.com/romangurevitch/go-training/internal/bank/api/middleware"
	"github.com/romangurevitch/go-training/internal/bank/service"
	apierror "github.com/romangurevitch/go-training/pkg/api/error"

	// TODO 2: uncomment when implementing OTel span
	// "go.opentelemetry.io/otel"
	// "go.opentelemetry.io/otel/attribute"
	// "log/slog"
	// "net/http"
)

// Handler handles transfer-related HTTP requests.
// Same shape as api/account/handler.go — see Handler and New there.
type Handler struct {
	svc service.Service
}

func New(svc service.Service) *Handler {
	return &Handler{svc: svc}
}

// CreateTransfer handles POST /v1/transfers.
//
// Quest steps:
//   TODO 1: Parse request body — api/account/handler.go CreateAccount, line ~40
//   TODO 2: Start OTel span — api/account/handler.go GetAccount, line ~20 (NEW: not in reference solution)
//   TODO 3: Verify ownership (JWT sub must be account owner) — NEW: not in reference solution
//   TODO 4: Call service.Transfer + map errors — api/account/handler.go CreateAccount, lines ~50-60
//   TODO 5: Log success + return 200 — api/account/handler.go CreateAccount, line ~65
func (h *Handler) CreateTransfer(c *gin.Context) {
	ctx := c.Request.Context()

	// TODO 1: Parse and validate request body.
	//   Pattern: api/account/handler.go CreateAccount — ShouldBindJSON + 400 on error
	//
	//   var req CreateTransferRequest
	//   if err := c.ShouldBindJSON(&req); err != nil {
	//       c.JSON(apierror.NewAPIError(ctx, http.StatusBadRequest, "bad request", err))
	//       return
	//   }

	// TODO 2: Start an OTel span.
	//   Pattern: api/account/handler.go GetAccount — otel.Tracer("bank").Start + defer span.End()
	//   NEW: the reference solution had no per-handler tracing — this is the upgrade.
	//
	//   ctx, span := otel.Tracer("bank").Start(ctx, "transfer.create")
	//   defer span.End()
	//   span.SetAttributes(
	//       attribute.String("from_account_id", req.FromAccountID),
	//       attribute.String("to_account_id",   req.ToAccountID),
	//       attribute.Float64("amount",         req.Amount),
	//   )

	// TODO 3: Verify ownership — JWT sub claim must match the from_account owner.
	//   NEW: not in reference solution.
	//
	//   claims := middleware.ClaimsFromCtx(ctx)
	//   fromAccount, err := h.svc.GetAccount(ctx, req.FromAccountID)
	//   switch {
	//   case errors.Is(err, domain.ErrAccountNotFound):
	//       c.JSON(apierror.NewAPIError(ctx, http.StatusNotFound, "account not found", err))
	//       return
	//   case err != nil:
	//       c.JSON(apierror.NewAPIError(ctx, http.StatusInternalServerError, "could not get account", err))
	//       return
	//   }
	//   if fromAccount.Owner != claims.Subject {
	//       c.JSON(apierror.NewAPIError(ctx, http.StatusForbidden, "forbidden: not account owner", nil))
	//       return
	//   }

	// TODO 4: Call service and map errors with errors.Is.
	//   Pattern: api/account/handler.go CreateAccount — switch errors.Is pattern
	//
	//   err = h.svc.Transfer(ctx, req.FromAccountID, req.ToAccountID, req.Amount)
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
	//       // TODO 5 goes here
	//   }

	// TODO 5: Log success and return 200.
	//   Pattern: api/account/handler.go CreateAccount — slog.InfoContext + c.JSON
	//   Note: trace_id and span_id are injected automatically by OtelHandler — no extra code.
	//
	//   slog.InfoContext(ctx, "transfer completed",
	//       slog.String("from_account_id", req.FromAccountID),
	//       slog.String("to_account_id",   req.ToAccountID),
	//       slog.Float64("amount",         req.Amount),
	//   )
	//   c.JSON(http.StatusOK, TransferResponse{Status: "completed"})

	// REMOVE THIS LINE when you implement TODO 1:
	_, _ = errors.New(""), middleware.ClaimsFromCtx(ctx) // silence unused imports
	_ = domain.ErrAccountNotFound                         // silence unused import
	_ = apierror.ErrInternalServerError                   // silence unused import
}
```

### Step 10c: Test Skeleton

- [ ] **Step 3: Create `api/transfer/handler_test.go`**

```go
package transfer_test

// PARTICIPANT QUEST — complete the test cases marked TODO.
//
// Before starting: read api/account/handler_test.go in full.
// The setup helpers (testToken, setupRouter) follow the identical pattern.
//
// Time estimate: 20-30 minutes for Step 4 of the quest.

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/romangurevitch/go-training/internal/bank/api/middleware"
	"github.com/romangurevitch/go-training/internal/bank/api/transfer"
	"github.com/romangurevitch/go-training/internal/bank/domain"
	"github.com/romangurevitch/go-training/internal/bank/service/mocks"
)

const testSecret = "test-secret"

// testToken issues a signed JWT for test Authorization headers.
// Pattern: identical to api/account/handler_test.go — copy it here.
func testToken(t *testing.T, sub, scope string) string {
	t.Helper()
	claims := middleware.Claims{
		Scope: scope,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   sub,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(testSecret))
	require.NoError(t, err)
	return signed
}

// setupRouter builds a minimal Gin engine with the transfer route.
// Pattern: identical to api/account/handler_test.go — copy and adapt for transfers.
func setupRouter(svc *mocks.MockService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := transfer.New(svc)

	v1 := r.Group("/v1/transfers")
	v1.Use(middleware.JWTMiddleware(testSecret))
	{
		v1.POST("", middleware.RequireScope("transfers:write"), h.CreateTransfer)
	}
	return r
}

var aliceAccount = &domain.Account{ID: "ACC-001", Owner: "alice", Balance: 500.0, Status: domain.StatusOpen}
var bobAccount = &domain.Account{ID: "ACC-002", Owner: "bob", Balance: 0.0, Status: domain.StatusOpen}

func TestCreateTransfer(t *testing.T) {
	type fields struct {
		svc func(t *testing.T) *mocks.MockService
	}
	tests := []struct {
		name     string
		fields   fields
		body     any
		tokenSub string
		wantCode int
	}{
		// PRE-WRITTEN: Happy path — transfer succeeds
		{
			name: "success — 200",
			fields: fields{
				svc: func(t *testing.T) *mocks.MockService {
					m := mocks.NewMockService(gomock.NewController(t))
					m.EXPECT().GetAccount(gomock.Any(), "ACC-001").Return(aliceAccount, nil).Times(1)
					m.EXPECT().Transfer(gomock.Any(), "ACC-001", "ACC-002", 50.0).Return(nil).Times(1)
					return m
				},
			},
			body:     map[string]any{"from_account_id": "ACC-001", "to_account_id": "ACC-002", "amount": 50.0},
			tokenSub: "alice",
			wantCode: http.StatusOK,
		},
		// PRE-WRITTEN: Invalid body — 400
		{
			name: "missing amount — 400",
			fields: fields{
				svc: func(t *testing.T) *mocks.MockService {
					return mocks.NewMockService(gomock.NewController(t)) // no calls expected
				},
			},
			body:     map[string]string{"from_account_id": "ACC-001", "to_account_id": "ACC-002"}, // amount missing
			tokenSub: "alice",
			wantCode: http.StatusBadRequest,
		},

		// TODO: Wrong owner — sub is "alice" but from_account is owned by "bob"
		// Expected: 403
		// Mock setup: GetAccount("ACC-002") returns bobAccount (owner: "bob")
		// Token sub: "alice"
		// {
		//     name: "wrong owner — 403",
		//     ...
		// },

		// TODO: Insufficient funds — transfer amount exceeds from_account balance
		// Expected: 422
		// Mock setup: GetAccount("ACC-001") returns aliceAccount; Transfer returns domain.ErrInsufficientFunds
		// Token sub: "alice"
		// {
		//     name: "insufficient funds — 422",
		//     ...
		// },

		// TODO: Source account not found
		// Expected: 404
		// Mock setup: GetAccount("MISSING") returns nil, domain.ErrAccountNotFound
		// Token sub: "alice"
		// {
		//     name: "source account not found — 404",
		//     ...
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter(tt.fields.svc(t))
			bodyBytes, _ := json.Marshal(tt.body)
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/v1/transfers", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+testToken(t, tt.tokenSub, "transfers:write"))
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)
		})
	}
}
```

- [ ] **Step 4: Run skeleton tests (pre-written cases should pass once handler TODOs are implemented)**

```bash
go build ./internal/bank/api/transfer/...
```

Expected: compiles (handler skeleton compiles before TODOs filled in).

- [ ] **Step 5: Commit**

```bash
git add internal/bank/api/transfer/
git commit -m "feat: add transfer handler skeleton — 5 guided TODOs for participant quest"
```

---

## Task 11: Create API Server (`api/server.go`)

**Files:**
- Modify: `internal/bank/api/server.go`

Replaces the existing flat `NewServer(h *Handler)` with the new multi-middleware, multi-group structure.

- [ ] **Step 1: Replace `api/server.go`**

```go
package api

import (
	"github.com/gin-gonic/gin"
	"log/slog"

	apiauth "github.com/romangurevitch/go-training/internal/bank/api/auth"
	"github.com/romangurevitch/go-training/internal/bank/api/account"
	"github.com/romangurevitch/go-training/internal/bank/api/middleware"
	"github.com/romangurevitch/go-training/internal/bank/api/transfer"
	"github.com/romangurevitch/go-training/internal/bank/service"
)

// Config holds the API server configuration.
type Config struct {
	JWTSecret   string
	ServiceName string // used by OTel span names
}

// NewServer builds the Gin engine with all middleware and routes pre-wired.
//
// Source: go-training-cba-solution/internal/server/rest/account/server.go
// Changes vs reference:
//   - gin.BasicAuth → JWTMiddleware + RequireScope per route
//   - JSONLogMiddleware (logrus) → LoggingMiddleware (slog-gin)
//   - Added TracingMiddleware (otelgin) and RequestIDMiddleware (uuid)
//   - Added transfer group (pre-wired reference; TODO for participants in server.go)
func NewServer(svc service.Service, logger *slog.Logger, cfg Config) *gin.Engine {
	r := gin.New()

	// Middleware stack — applied to all routes in order:
	r.Use(middleware.RequestIDMiddleware())           // 1. Generate + inject X-Request-Id
	r.Use(middleware.TracingMiddleware(cfg.ServiceName)) // 2. Start OTel span (otelgin)
	r.Use(middleware.LoggingMiddleware(logger))       // 3. Structured request log (slog-gin)
	r.Use(gin.Recovery())                            // 4. Recover panics → 500

	// Auth handler — no JWT required (issues tokens)
	authHandler := apiauth.New(cfg.JWTSecret)
	r.POST("/v1/token", authHandler.IssueToken)

	// Accounts — REFERENCE: fully wired, participants read this
	accountHandler := account.New(svc)
	accounts := r.Group("/v1/accounts")
	accounts.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	{
		accounts.GET("/:id", middleware.RequireScope("accounts:read"), accountHandler.GetAccount)
		accounts.POST("", middleware.RequireScope("accounts:write"), accountHandler.CreateAccount)
	}

	// Transfers — TODO for participants (Step 2 of quest)
	// Pattern: identical to the accounts group above.
	//
	// transferHandler := transfer.New(svc)
	// transfers := r.Group("/v1/transfers")
	// transfers.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	// {
	//     transfers.POST("", middleware.RequireScope("transfers:write"), transferHandler.CreateTransfer)
	// }
	var _ = transfer.New(nil) // prevent unused import error while TODO is commented out

	return r
}
```

- [ ] **Step 2: Build**

```bash
go build ./internal/bank/api/...
```

Expected: compiles cleanly.

- [ ] **Step 3: Commit**

```bash
git add internal/bank/api/server.go
git commit -m "feat: update api server — JWT auth, OTel tracing, slog logging middleware wired"
```

---

## Task 12: Update `cmd/bank-api/main.go` — OTel Bootstrap + slog Setup

**Files:**
- Modify: `cmd/bank-api/main.go`

Replaces the existing simple main with OTel SDK bootstrap + slog multi-handler setup. All pre-built for participants.

- [ ] **Step 1: Replace `cmd/bank-api/main.go`**

```go
package main

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	slogotel "github.com/remychantenay/slog-otel"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	bankapi "github.com/romangurevitch/go-training/internal/bank/api"
	"github.com/romangurevitch/go-training/internal/bank/repository/postgres"
	"github.com/romangurevitch/go-training/internal/bank/service"
)

func main() {
	ctx := context.Background()

	// --- Logging Setup ---
	// Go 1.26: slog.NewMultiHandler fans out to JSON stdout and OTel log bridge.
	// slogotel.OtelHandler wraps both — injects trace_id/span_id into every log line.
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	otelBridge := otelslog.NewHandler("bank-api")
	otelEnricher := slogotel.OtelHandler{Next: slog.NewMultiHandler(jsonHandler, otelBridge)}
	logger := slog.New(otelEnricher)
	slog.SetDefault(logger)

	// --- OTel Tracing Setup ---
	tp, err := setupTracer(ctx)
	if err != nil {
		log.Fatalf("failed to setup tracer: %v", err)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tp.Shutdown(shutdownCtx); err != nil {
			slog.ErrorContext(ctx, "tracer shutdown failed", slog.Any("error", err))
		}
	}()

	// --- Database ---
	dsn := envOrDefault("DATABASE_URL", "postgres://gotrainer:verysecret@localhost:5432/gobank?sslmode=disable")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	// --- Wire Dependencies ---
	repo := postgres.New(db)
	svc := service.NewBankService(repo)

	jwtSecret := envOrDefault("JWT_SECRET", "super-secret-for-training-only")
	cfg := bankapi.Config{
		JWTSecret:   jwtSecret,
		ServiceName: "bank-api",
	}

	router := bankapi.NewServer(svc, logger, cfg)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// --- Start Server ---
	go func() {
		slog.InfoContext(ctx, "server starting", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.ErrorContext(ctx, "server error", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	// --- Graceful Shutdown ---
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.InfoContext(ctx, "shutting down server")

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("server forced to shutdown:", err)
	}
	slog.InfoContext(ctx, "server exited")
}

// setupTracer initializes the OTel SDK with OTLP/HTTP exporter pointing at Jaeger.
// Jaeger all-in-one accepts OTLP on port 4318.
func setupTracer(ctx context.Context) (*sdktrace.TracerProvider, error) {
	endpoint := envOrDefault("OTEL_EXPORTER_OTLP_ENDPOINT", "http://localhost:4318")

	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("bank-api"),
		semconv.ServiceVersion("0.1.0"),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	return tp, nil
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
```

- [ ] **Step 2: Build**

```bash
go build ./cmd/bank-api/...
```

Expected: compiles cleanly.

- [ ] **Step 3: Commit**

```bash
git add cmd/bank-api/main.go
git commit -m "feat: update main.go — OTel SDK bootstrap, slog multi-handler, Postgres wiring"
```

---

## Task 13: Update `docker-compose.yaml` — Add Jaeger

**Files:**
- Modify: `docker-compose.yaml`

- [ ] **Step 1: Replace `docker-compose.yaml`**

```yaml
services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: gobank
      POSTGRES_USER: gotrainer
      POSTGRES_PASSWORD: verysecret
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U gotrainer -d gobank"]
      interval: 5s
      timeout: 5s
      retries: 5

  jaeger:
    image: jaegertracing/all-in-one:latest
    environment:
      COLLECTOR_OTLP_ENABLED: "true"
    ports:
      - "16686:16686"  # Jaeger UI
      - "4317:4317"    # OTLP gRPC
      - "4318:4318"    # OTLP HTTP — bank-api exports here

  bank-api:
    build: .
    depends_on:
      postgres:
        condition: service_healthy
      jaeger:
        condition: service_started
    environment:
      DATABASE_URL: "postgres://gotrainer:verysecret@postgres:5432/gobank?sslmode=disable"
      OTEL_EXPORTER_OTLP_ENDPOINT: "http://jaeger:4318"
      JWT_SECRET: "super-secret-for-training-only"
    ports:
      - "8080:8080"
```

- [ ] **Step 2: Verify Jaeger starts**

```bash
docker compose up jaeger -d
# Wait a few seconds
curl -s http://localhost:16686/api/services | head -1
```

Expected: JSON response (Jaeger UI is up).

- [ ] **Step 3: Commit**

```bash
git add docker-compose.yaml
git commit -m "feat: add Jaeger to docker-compose — OTLP HTTP on 4318, UI on 16686"
```

---

## Task 14: Create OpenAPI Specs

**Files:**
- Create: `docs/openapi/accounts.yaml`
- Create: `docs/openapi/transfers.yaml`

### Step 14a: Accounts spec (complete reference)

- [ ] **Step 1: Create `docs/openapi/accounts.yaml`**

```yaml
openapi: 3.0.3
info:
  title: Bank Accounts API
  version: 1.0.0
  description: Reference API for the Go Bank training exercise. Participants read this before designing transfers.yaml.

paths:
  /v1/accounts/{id}:
    get:
      summary: Get account by ID
      operationId: getAccount
      security:
        - bearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
          example: ACC-001
      responses:
        "200":
          description: Account found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/AccountResponse'
        "401":
          description: Missing or invalid JWT token
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        "403":
          description: Token does not have accounts:read scope
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        "404":
          description: Account not found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        "500":
          description: Internal server error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

  /v1/accounts:
    post:
      summary: Create a new account
      operationId: createAccount
      security:
        - bearerAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreateAccountRequest'
      responses:
        "201":
          description: Account created
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/AccountResponse'
        "400":
          description: Missing or invalid request body
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        "401":
          description: Missing or invalid JWT token
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        "403":
          description: Token does not have accounts:write scope
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        "409":
          description: Account already exists
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        "500":
          description: Internal server error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
  schemas:
    CreateAccountRequest:
      type: object
      required:
        - owner
      properties:
        owner:
          type: string
          description: Account owner name
          example: alice

    AccountResponse:
      type: object
      properties:
        id:
          type: string
          example: ACC-001
        owner:
          type: string
          example: alice
        balance:
          type: number
          format: double
          example: 100.00
        status:
          type: string
          enum: [OPEN, LOCKED, CLOSED]
          example: OPEN
        created_at:
          type: string
          format: date-time
        updated_at:
          type: string
          format: date-time

    ErrorResponse:
      type: object
      properties:
        message:
          type: string
          example: account not found
```

### Step 14b: Transfers spec (partial — participants complete)

- [ ] **Step 2: Create `docs/openapi/transfers.yaml`**

```yaml
openapi: 3.0.3
info:
  title: Bank Transfer API
  version: 1.0.0
  description: |
    PARTICIPANT QUEST — Step 1 (20 minutes).

    Complete the TODOs below using docs/openapi/accounts.yaml as your reference.
    - The schema shapes are identical to accounts.yaml
    - The security and error response patterns are identical
    - Focus on: what fields does the request need? What do each status code mean?

paths:
  /v1/transfers:
    post:
      summary: Transfer funds between two accounts
      security:
        - bearerAuth: []
      # TODO: Add operationId (hint: look at accounts.yaml for the pattern)

      requestBody:
        required: true
        content:
          application/json:
            schema:
              # TODO: Define the TransferRequest schema inline or via $ref.
              #
              # Fields to include:
              #   from_account_id: string (required) — source account
              #   to_account_id:   string (required) — destination account
              #   amount:          number (required) — must be > 0
              #
              # Reference: CreateAccountRequest in accounts.yaml for schema structure

      responses:
        "200":
          description: Transfer completed successfully
          content:
            application/json:
              schema:
                # TODO: Define the success response schema.
                # Body: { "status": "completed" }
                # Hint: one field — status (string)

        "400":
          # TODO: Describe when this happens and define the response body.
          # Hint: malformed JSON or missing required fields
          # ErrorResponse shape is identical to accounts.yaml — use $ref

        "401":
          description: Missing or invalid JWT token
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
          # Note: 401 is pre-filled — the pattern is identical to accounts.yaml

        "403":
          # TODO: Two distinct 403 cases exist for this endpoint:
          #   Case 1: JWT token is missing the transfers:write scope (middleware rejects)
          #   Case 2: The authenticated user (JWT sub) is not the owner of from_account_id (handler rejects)
          # Both return 403 — describe this in the description field.

        "404":
          # TODO: Account not found.
          # Hint: either source or destination — the response body does not distinguish (intentional)

        "422":
          # TODO: Business rule violation — the account exists but the transfer cannot proceed.
          # Hint: two cases — what are they? (check domain/account.go for the sentinel errors)

        "500":
          # TODO: Unexpected internal error — same shape as accounts.yaml 500

components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
  schemas:
    ErrorResponse:
      type: object
      properties:
        message:
          type: string
          example: insufficient funds
    # TODO: Add TransferRequest and TransferResponse schemas here
    # Reference: accounts.yaml components/schemas for structure
```

- [ ] **Step 3: Commit**

```bash
git add docs/openapi/
git commit -m "feat: add OpenAPI specs — accounts.yaml (complete) + transfers.yaml (participant TODO)"
```

---

## Task 15: Build Bonus Client (Required for Cleanup)

**Why now:** Task 16 deletes `internal/bank/client/client.go`. The `cmd/bank-cli/main.go` imports it. We must replace it with `pkg/client/bank/` *before* deletion or Task 16's build verification step will fail.

- [ ] **Step 1: Go to Task 18 and complete it now**

Jump ahead to **Task 18 (Bonus — Transfer Client)** and complete all its steps (creates `pkg/client/bank/api/types.go` and `pkg/client/bank/client.go`). Return here when done.

- [ ] **Step 2: Verify the new client package builds**

```bash
go build ./pkg/client/bank/...
```

Expected: compiles cleanly.

- [ ] **Step 3: Return to Task 16**

---

## Task 16: Remove Old Files

**Files:**
- Delete: `internal/bank/store/` (replaced by `repository/`)
- Delete: `internal/bank/api/handlers.go` (replaced by `api/account/` and `api/transfer/`)
- Delete: `internal/bank/api/middleware.go` (replaced by `api/middleware/` package)
- Delete: `internal/bank/client/client.go` (replaced by `pkg/client/bank/` in bonus)

- [ ] **Step 1: Delete old store package**

```bash
rm -rf internal/bank/store/
```

- [ ] **Step 2: Delete old flat API files**

```bash
rm internal/bank/api/handlers.go
rm internal/bank/api/middleware.go
```

- [ ] **Step 3: Delete old thin client (replaced by pkg/client/bank/ in bonus)**

```bash
rm internal/bank/client/client.go
```

- [ ] **Step 4: Update `cmd/bank-cli/main.go` imports** (the old `client.NewBankClient` is gone)

The CLI's bank client will use the new `pkg/client/bank/` client after the bonus step. For now, update `cmd/bank-cli/main.go` to use a placeholder or the new client:

```go
package main

import (
	"net/http"
	"os"

	"github.com/romangurevitch/go-training/internal/bank/cli"
	bankclient "github.com/romangurevitch/go-training/pkg/client/bank"
	"github.com/spf13/cobra"
)

func main() {
	apiURL := os.Getenv("BANK_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}

	bankClient := bankclient.New(apiURL, &http.Client{})
	cobra.CheckErr(cli.New(bankClient).Execute())
}
```

- [ ] **Step 4b: Update `internal/bank/cli/cli.go` to use the new client interface**

The old `cli.New` signature was `func New(bankClient client.BankClient)` where `client.BankClient` is the old `internal/bank/client/` interface. This type no longer exists after deletion.

Update `internal/bank/cli/cli.go`:
```go
// Old:
import "github.com/romangurevitch/go-training/internal/bank/client"
func New(bankClient client.BankClient) *cobra.Command { ... }

// New:
import bankclient "github.com/romangurevitch/go-training/pkg/client/bank"
func New(bankClient bankclient.Client) *cobra.Command { ... }
```

Update `internal/bank/cli/account/` files in the same way — replace any `client.BankClient` parameter types with `bankclient.Client`. The method names on both interfaces are identical (`GetAccount`, `CreateAccount`) so no further changes are needed inside the command implementations.

- [ ] **Step 5: Build everything to verify no broken imports**

```bash
go build ./...
```

Expected: clean build. Fix any remaining import errors:
- `cmd/bank-api/main.go` still importing `store` or `logger` packages → already fixed in Task 12
- `internal/bank/cli/` importing old client → fixed in Step 4b above

- [ ] **Step 6: Run all tests**

```bash
go test ./...
```

Expected: all existing tests pass; newly written tests pass.

- [ ] **Step 7: Commit**

```bash
git add -A
git commit -m "refactor: remove old store/, handlers.go, middleware.go — replaced by repository/ and api/ packages"
```

---

## Task 17: Update Makefile

**Files:**
- Modify: `Makefile`

- [ ] **Step 1: Add Jaeger, OTel, and quest targets to Makefile**

Add these targets:

```makefile
# --- Observability ---
jaeger-up: ## Start Jaeger (traces at localhost:16686)
	docker compose up jaeger -d

jaeger-down: ## Stop Jaeger
	docker compose stop jaeger

# --- Quest Helpers ---
token: ## Get a JWT token for testing (usage: make token SUB=alice SCOPE="transfers:write")
	@curl -s -X POST http://localhost:8080/v1/token \
		-H 'Content-Type: application/json' \
		-d "{\"sub\":\"$(SUB)\",\"scope\":\"$(SCOPE)\"}" | jq .

transfer: ## Execute a transfer (usage: make transfer TOKEN=<jwt> FROM=ACC-001 TO=ACC-002 AMOUNT=50)
	@curl -s -X POST http://localhost:8080/v1/transfers \
		-H 'Content-Type: application/json' \
		-H "Authorization: Bearer $(TOKEN)" \
		-d "{\"from_account_id\":\"$(FROM)\",\"to_account_id\":\"$(TO)\",\"amount\":$(AMOUNT)}" | jq .

# --- Mocks ---
mocks: ## Regenerate all mockery mocks
	go generate ./...
```

- [ ] **Step 2: Build to verify Makefile syntax**

```bash
make help
```

Expected: all targets listed.

- [ ] **Step 3: Commit**

```bash
git add Makefile
git commit -m "feat: update Makefile — add jaeger, token, transfer, mocks targets"
```

---

## Task 18: Bonus — Transfer Client (`pkg/client/bank/`)

**Files:**
- Create: `pkg/client/bank/api/types.go`
- Create: `pkg/client/bank/client.go`

For participants who finish early. Ported directly from `go-training-cba-solution/pkg/client/rest/account/client.go`.
Build this task BEFORE Task 16 (Remove Old Files) — see Task 15 prerequisite note above.

### Step 18a: Types (pkg/client/bank/api/types.go)

- [ ] **Step 1: Create `pkg/client/bank/api/types.go`**

```go
package api

// AccountResponse mirrors internal/bank/api/account/types.go — kept in sync manually.
// In a real project this would be an OpenAPI-generated shared schema.
type AccountResponse struct {
	ID      string  `json:"id"`
	Owner   string  `json:"owner"`
	Balance float64 `json:"balance"`
	Status  string  `json:"status"`
}

// CreateAccountRequest is the request body for POST /v1/accounts.
type CreateAccountRequest struct {
	Owner string `json:"owner"`
}

// TransferRequest is the request body for POST /v1/transfers.
type TransferRequest struct {
	FromAccountID string  `json:"from_account_id" validate:"required"`
	ToAccountID   string  `json:"to_account_id"   validate:"required"`
	Amount        float64 `json:"amount"           validate:"required,gt=0"`
}

// TransferResponse is the success body for POST /v1/transfers.
type TransferResponse struct {
	Status string `json:"status"`
}

// TokenRequest is the request body for POST /v1/token.
type TokenRequest struct {
	Sub   string `json:"sub"`
	Scope string `json:"scope"`
}

// TokenResponse wraps the JWT issued by POST /v1/token.
type TokenResponse struct {
	Token string `json:"token"`
}
```

### Step 18b: Client

- [ ] **Step 2: Create `pkg/client/bank/client.go`**

```go
package bank

// Bank API client — ported from go-training-cba-solution/pkg/client/rest/account/client.go.
// Changes vs reference:
//   - Basic Auth → Bearer JWT (GetToken first, then set on subsequent calls)
//   - interface{} → any (Go 1.18+)
//   - TraceDuration → slog defer pattern

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/romangurevitch/go-training/pkg/client/bank/api"
	httppkg "github.com/romangurevitch/go-training/pkg/http"
)

// Client is the typed interface for the Bank API.
type Client interface {
	GetToken(ctx context.Context, sub, scope string) (string, error)
	GetAccount(ctx context.Context, id string) (*api.AccountResponse, error)
	Transfer(ctx context.Context, req *api.TransferRequest) (*api.TransferResponse, error)
}

type client struct {
	basePath   string
	httpClient *http.Client
	token      string // set by GetToken, sent as Bearer on all subsequent calls
	validate   *validator.Validate
}

func New(basePath string, httpClient *http.Client) Client {
	return &client{
		basePath:   basePath,
		httpClient: httpClient,
		validate:   validator.New(),
	}
}

// GetToken calls POST /v1/token and stores the JWT for future requests.
// Source: adapted from reference client's Basic Auth setup.
func (c *client) GetToken(ctx context.Context, sub, scope string) (string, error) {
	payload, err := json.Marshal(api.TokenRequest{Sub: sub, Scope: scope})
	if err != nil {
		return "", err
	}

	urlPath, err := httppkg.GetURL(c.basePath, "v1/token", "")
	if err != nil {
		return "", err
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, urlPath, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	k, v := httppkg.HeaderApplicationJSON()
	r.Header.Set(k, v)

	body, err := httppkg.DoRequest(ctx, c.httpClient, r, http.StatusOK)
	if err != nil {
		return "", err
	}

	var res api.TokenResponse
	if err = json.Unmarshal(body, &res); err != nil {
		return "", err
	}

	c.token = res.Token
	return res.Token, nil
}

// GetAccount calls GET /v1/accounts/:id.
// Source: go-training-cba-solution/pkg/client/rest/account/client.go GetAccount
// Changes: Basic Auth header → Bearer JWT
func (c *client) GetAccount(ctx context.Context, id string) (*api.AccountResponse, error) {
	// Reference used: defer logger.TraceDuration(ctx, time.Now(), "GetAccount")
	// slog equivalent:
	start := time.Now()
	defer func() {
		slog.DebugContext(ctx, "GetAccount", slog.Duration("duration", time.Since(start)))
	}()

	if err := c.validate.Var(id, "required"); err != nil {
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

	body, err := httppkg.DoRequest(ctx, c.httpClient, r, http.StatusOK)
	if err != nil {
		return nil, err
	}

	var res api.AccountResponse
	return &res, json.Unmarshal(body, &res)
}

// Transfer calls POST /v1/transfers.
//
// PARTICIPANT BONUS QUEST — implement this method.
// Read GetAccount above first — the pattern is identical.
//
// TODO 1: Validate req using c.validate.Struct(req) — returns error if required fields missing
// TODO 2: Marshal req to JSON — json.Marshal(req)
// TODO 3: Build URL — httppkg.GetURL(c.basePath, "v1/transfers", "")
// TODO 4: Build request — http.NewRequestWithContext + set Content-Type + Authorization: Bearer c.token
// TODO 5: httppkg.DoRequest(ctx, c.httpClient, r, http.StatusOK) — returns *APIError on failure
// TODO 6: json.Unmarshal body into *api.TransferResponse and return
func (c *client) Transfer(ctx context.Context, req *api.TransferRequest) (*api.TransferResponse, error) {
	start := time.Now()
	defer func() {
		slog.DebugContext(ctx, "Transfer", slog.Duration("duration", time.Since(start)))
	}()

	// TODO 1-6 here
	return nil, nil
}
```

- [ ] **Step 3: Build**

```bash
go build ./pkg/client/bank/...
```

Expected: compiles cleanly.

- [ ] **Step 4: Commit**

```bash
git add pkg/client/bank/
git commit -m "feat: add pkg/client/bank — typed bank client with GetToken, GetAccount; Transfer TODO for bonus"
```

---

## Task 19: End-to-End Smoke Test

Verify the entire scaffold works before handing to participants.

- [ ] **Step 1: Start infrastructure**

```bash
docker compose up postgres jaeger -d
# Run migrations (the Makefile target just prints instructions — run psql directly):
psql -h localhost -U gotrainer -d gobank -f migration/001_create_accounts.sql
psql -h localhost -U gotrainer -d gobank -f migration/002_create_transactions.sql
```

- [ ] **Step 2: Start the API**

```bash
make run-bank-api
```

Expected: `{"time":"...","level":"INFO","msg":"server starting","addr":":8080"}` in JSON.

- [ ] **Step 3: Get a token**

```bash
curl -s -X POST http://localhost:8080/v1/token \
  -H 'Content-Type: application/json' \
  -d '{"sub":"alice","scope":"accounts:read accounts:write transfers:write"}' | jq .
```

Expected: `{ "token": "eyJ..." }`

- [ ] **Step 4: Create an account**

```bash
export TOKEN="<paste token>"
curl -s -X POST http://localhost:8080/v1/accounts \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"owner":"alice"}' | jq .
```

Expected: 201 with account ID.

- [ ] **Step 5: Try the transfer endpoint (expect 404 — route not wired yet)**

```bash
curl -s -X POST http://localhost:8080/v1/transfers \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"from_account_id":"ACC-001","to_account_id":"ACC-002","amount":50}' | jq .
```

Expected: 404 (transfer route is commented out in `server.go` — this is correct pre-quest state).

- [ ] **Step 6: Check Jaeger**

Open `http://localhost:16686` → select service `bank-api` → find `POST /v1/accounts` trace.

Expected: trace with `http.request` and `account.create` child span, both visible.

- [ ] **Step 7: Run all tests**

```bash
go test ./... -count=1
```

Expected: all tests PASS.

- [ ] **Step 8: Final commit**

```bash
git add -A
git commit -m "feat: complete go bank transfer quest scaffold — ready for participants"
```

---

## Task 20: Fix Stale Reference in Spec (cleanup)

**Files:**
- Modify: `docs/superpowers/specs/2026-03-18-go-bank-transfer-quest-design.md`

The spec's Success Criteria section has a stale reference (`logger.WithContext(ctx)` — old logrus pattern). Fix it to match the slog pattern used throughout.

- [ ] **Step 1: Fix Success Criteria in spec**

Find this line:
```
- [ ] Log lines for success and failure use `logger.WithContext(ctx)` with structured fields
```

Replace with:
```
- [ ] Log lines for success and failure use `slog.InfoContext(ctx, ...)` / `slog.ErrorContext(ctx, ...)` with typed slog attributes
```

- [ ] **Step 2: Commit**

```bash
git add docs/superpowers/specs/2026-03-18-go-bank-transfer-quest-design.md
git commit -m "docs: fix stale logrus reference in quest success criteria → slog"
```
