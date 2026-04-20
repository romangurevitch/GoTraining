# Technical Specification: Durable Transfer Workflow

## 1. Overview
This document defines the implementation details for the `DurableTransferWorkflow`, which handles high-value bank transfers with human-in-the-loop approval and automated compensation for failures.

## 2. Workflow Definition

### 2.1. Interface
- **Workflow Name:** `DurableTransferWorkflow`
- **Workflow ID:** `transfer-{transferID}` (deterministic — prevents duplicate workflow executions on API retry)
- **Task Queue:** Configurable via `config.Values.TemporalTaskQueue` (default: `bank-transfer-queue`)

**Input:**
```go
type TransferRequest struct {
    TransferID    string // Natural idempotency key
    FromAccountID string
    ToAccountID   string
    Amount        int64  // In cents
    Reference     string
}
```

**Output:**
```go
type TransferResponse struct {
    TransferID string
    Status     string // COMPLETED, REJECTED, FAILED
}
```

### 2.2. Signals
- `ApprovalSignal` (`"approval-signal"`): Approves a pending high-value transfer.
- `RejectSignal` (`"reject-signal"`): Rejects a pending high-value transfer.

### 2.3. Threshold
The approval gate activates when `Amount > 1000` cents. This value is specified directly in the PRD as the literal threshold. All amounts in the system are denominated in cents.

## 3. Activities

### 3.1. `ValidateAccounts(ctx, req TransferRequest) error`
- Verifies both source and destination accounts exist via `repository.GetAccount`.
- Checks both accounts pass `domain.Account.CanPerformTransaction()` (rejects `StatusLocked` and `StatusClosed`).
- Non-retryable errors: `ErrAccountNotFound`, `ErrAccountLocked`, `ErrAccountClosed`.

### 3.2. `DebitAccount(ctx, accountID, amount, transferID) error`
- Withdraws funds from the source account.
- Idempotency: derives transaction ID as `{transferID}-debit`. Checks `ListTransactions` for existing record before writing.
- Re-checks `CanPerformTransaction()` to guard against state changes since validation (TOCTOU mitigation).
- Non-retryable: `ErrInsufficientFunds`, `ErrAccountLocked`, `ErrAccountClosed`.

### 3.3. `CreditAccount(ctx, accountID, amount, transferID) error`
- Deposits funds into the destination account.
- Idempotency: derives transaction ID as `{transferID}-credit`.
- Re-checks `CanPerformTransaction()` to guard against state changes since validation.
- Non-retryable: `ErrAccountLocked`, `ErrAccountClosed` (triggers compensation).

### 3.4. `RefundDebitActivity(ctx, accountID, amount, transferID) error`
- Reverses a previous debit if the credit fails (compensation).
- Idempotency: derives transaction ID as `{transferID}-refund`.
- Executed using `workflow.NewDisconnectedContext` so it completes even if the parent workflow is cancelled.

## 4. Workflow Logic Flow

```
1. ValidateAccounts ──fail──▶ RETURN error (non-retryable)
        │
        ▼ success
2. Amount > 1000?
   ├─ yes ──▶ AwaitWithTimeout(24h) for ApprovalSignal / RejectSignal
   │           ├─ RejectSignal  ──▶ RETURN "REJECTED" (non-retryable)
   │           ├─ Timeout       ──▶ RETURN "FAILED"   (non-retryable)
   │           └─ ApprovalSignal──▶ continue
   └─ no ───▶ continue
        │
        ▼
3. DebitAccount ──fail──▶ RETURN error
        │
        ▼ success
4. CreditAccount ──fail──▶ 5. RefundDebitActivity (disconnected ctx)
        │                       └──▶ RETURN original credit error
        ▼ success
   RETURN "COMPLETED"
```

## 5. Failure Mode Analysis

### 5.1. TOCTOU (Time-of-Check-to-Time-of-Use)
Account status may change between `ValidateAccounts` and subsequent activities (e.g., account locked by compliance after validation passes). **Mitigation:** Both `DebitAccount` and `CreditAccount` re-check `CanPerformTransaction()` before modifying balances. This is a deliberate design choice — validation is an early fail-fast, activities are the authoritative check.

### 5.2. Partial Database Failures
Each activity performs two database writes: `SaveAccount` (balance update) and `SaveTransaction` (ledger entry). These are not wrapped in a database transaction. **Risk:** If `SaveAccount` succeeds but `SaveTransaction` fails, the balance is updated without a ledger record. **Mitigation:** Temporal will retry the activity, and the idempotency check on `transactionID` prevents duplicate balance modifications on retry. The failed `SaveTransaction` will be retried in the next attempt.

### 5.3. Compensation Failure
If `RefundDebitActivity` itself fails, the refund error is logged but the workflow returns the original credit error. Temporal's retry policy will attempt the refund activity multiple times before giving up. If all retries exhaust, manual intervention is required — this is surfaced via workflow error state in the Temporal Web UI.

### 5.4. Duplicate Workflow Prevention
The workflow is started with a deterministic `WorkflowID` (`transfer-{transferID}`). Temporal guarantees that a second `StartWorkflow` call with the same ID returns the existing execution rather than creating a duplicate.

## 6. Error Classification

| Error | Retryable | Rationale |
|-------|-----------|-----------|
| `ErrAccountNotFound` | No | Business logic — account doesn't exist |
| `ErrInsufficientFunds` | No | Business logic — can't debit |
| `ErrAccountLocked` | No | Compliance hold — requires manual resolution |
| `ErrAccountClosed` | No | Terminal state — cannot be reopened |
| `ErrInvalidAmount` | No | Validation — amount must be positive |
| DB connection timeout | Yes | Transient infrastructure failure |
| Network error | Yes | Transient infrastructure failure |

## 7. Activity Options

```go
ActivityOptions{
    StartToCloseTimeout: 10 * time.Second,
    RetryPolicy: {
        InitialInterval:    1s,
        BackoffCoefficient: 2.0,
        MaximumInterval:    1m,
    },
}
```
