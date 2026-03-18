# Go Bank Challenge

Implement the Go Bank service across four packages. Work in order — each layer depends on the one before it.

## Quest 1: Domain Models

**Package:** `internal/bank/domain/`

Define the core entities:

- `Account` struct with fields: `ID`, `FirstName`, `LastName`, `Email`, `Balance`, `Status` (`Active`/`Locked`), `CreatedAt`, `UpdatedAt`
- `Transaction` struct with fields: `ID`, `AccountID`, `Type` (`Credit`/`Debit`), `Amount`, `CreatedAt`
- Status and Type as typed enums (use `string` or `int` constants)

**Acceptance criteria:** Structs compile, JSON tags present, enums defined.

## Quest 2: Store / Repository

**Package:** `internal/bank/store/`

Implement a `Store` interface and a Postgres-backed implementation using go-jet:

- `FindAccountByID(ctx, id) (*domain.Account, error)`
- `InsertAccount(ctx, account) error`
- `UpdateAccount(ctx, account) error`
- `FindTransactionsByAccountID(ctx, accountID) ([]domain.Transaction, error)`
- `InsertTransaction(ctx, transaction) error`

Run migrations first: `make db-up && make migrate`

**Acceptance criteria:** `make test-bank` passes with a running Postgres.

## Quest 3: Service / Business Logic

**Package:** `internal/bank/service/`

Implement a `Service` interface:

- `GetAccount(ctx, id) (*domain.Account, error)`
- `Deposit(ctx, accountID, amount) error` — credit the account, insert transaction
- `Withdraw(ctx, accountID, amount) error` — debit the account, reject if insufficient funds, insert transaction

Use mocks (`go.uber.org/mock`) to test the service without a real database.

**Acceptance criteria:** `make test-bank` passes with mocked store.

## Quest 4: HTTP API

**Package:** `internal/bank/api/`

Wire up Gin routes:

- `GET /v1/accounts/:id` — return account JSON
- `POST /v1/accounts` — create account
- `POST /v1/accounts/:id/deposit` — deposit funds
- `POST /v1/accounts/:id/withdraw` — withdraw funds

Use `httptest` to test handlers without a running server.

**Acceptance criteria:** `make test-bank` passes, `make build` compiles the bank-api binary.
