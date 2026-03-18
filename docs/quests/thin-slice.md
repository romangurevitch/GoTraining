# Thin-Slice Banking Challenge

Welcome to the Go Bank Thin-Slice challenge. In this workshop, you will build a functional (though simplified) banking system consisting of a REST API and a CLI tool.

## Objectives
- Master Go 1.26 features and best practices.
- Implement structured logging with `log/slog`.
- Build a thread-safe in-memory data store.
- Create a RESTful API with Gin.
- Develop a command-line interface with Cobra.

---

## Quest 1: The Domain Model
**File:** `internal/bank/domain/account.go` & `transaction.go`

Review the `Account` and `Transaction` structs. Ensure you understand the JSON tags and how they map to the API responses. 
- **Task:** Verify that the `CanPerformTransaction` method correctly handles different account statuses.

---

## Quest 2: Memory Store
**File:** `internal/bank/store/memory.go`

The `MemoryStore` is partially implemented. You need to ensure it is thread-safe and supports transaction history.
- **Task:** Implement the `SaveTransaction` and `ListTransactions` methods. 
- **Constraint:** Use `s.mu.Lock()` or `s.mu.RLock()` appropriately to prevent race conditions.

---

## Quest 3: Business Logic
**File:** `internal/bank/service/service.go`

The `BankService` handles the core logic. `Deposit` is provided as an example.
- **Task:** Implement the `Withdraw` method. 
- **Requirements:**
    1. Validate the amount is positive.
    2. Retrieve the account from the store.
    3. Check if the account can perform transactions (use the domain method).
    4. Ensure the account has sufficient funds.
    5. Update the balance and save a `Transaction` record of type `WITHDRAWAL`.

---

## Quest 4: The REST API
**File:** `internal/bank/api/handlers.go`

Expose your banking logic via HTTP.
- **Task:** Implement the `GetAccount` handler. 
- **Requirements:**
    1. Extract the `id` from the URL path.
    2. Call the `bankService.GetAccount` method.
    3. Return `200 OK` with the account JSON, or `404 Not Found` if the account doesn't exist.

---

## Quest 5: The CLI Tool
**Files:** `internal/bank/cli/account/balance.go`

Make the bank accessible from the terminal.
- **Task:** Implement the `balance` command in `internal/bank/cli/account/balance.go`.
- **Requirements:** It should call `bankClient.GetAccount` and print the balance clearly to the user.

---

## Quest 6: Testing
**File:** `internal/bank/service/service_test.go` (You need to create this)

Validation is key to production-ready software.
- **Task:** Write a table-driven test for the `Deposit` method in `BankService`.
- **Scenarios to cover:**
    - Success scenario.
    - Invalid (negative) amount.
    - Deposit to a locked account.
    - Deposit to a non-existent account.

---

## Run and Verify

### 1. Start the API Server
```bash
go run cmd/bank-api/main.go
```

### 2. Use the CLI
```bash
# Create an account
go run cmd/bank-cli/main.go account create "Jane Doe"

# Check balance (Implementation target for Quest 5)
go run cmd/bank-cli/main.go account balance ACC-XXX
```
