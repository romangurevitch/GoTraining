# PRD: High-Value Durable Transfer Workflow

## 1. Overview
The current bank transfer implementation is a simple, synchronous database transaction. While effective for basic cases, it doesn't scale to complex, multi-stage processes that require:
- **Human-in-the-loop:** Approvals for large amounts.
- **External Integrations:** Notifying fraud detection or external ledgers.
- **Reliability:** Ensuring the transfer completes even if the system crashes between debit and credit.

This project introduces a **Temporal Workflow** to orchestrate a "High-Value Transfer Workflow".

## 2. Goals
- Implement a durable transfer process using Temporal.
- Demonstrate the **Compensation Pattern** for distributed transactions.
- Implement **Signal Handling** for manual approvals.
- Ensure **Idempotency** across retries.

## 3. Requirements

### 3.1. Workflow: `DurableTransferWorkflow`
- **Input:** `TransferRequest { FromAccountID, ToAccountID, Amount, Reference }`
- **Logic:**
    1. **Validation:** Check if the amount is positive and accounts exist (Activity).
    2. **Approval Gate:** If `Amount > 1000`, the workflow must wait for an `ApprovalSignal`.
        - If no signal is received within 24 hours, the workflow should time out and fail.
        - If a `RejectSignal` is received, the workflow should return a `temporal.NewNonRetryableApplicationError("transfer rejected", "REJECTED", nil)`.
    3. **Debit Stage:** Call `DebitActivity` to remove funds from the `FromAccount`.
    4. **Credit Stage:** Call `CreditActivity` to add funds to the `ToAccount`.
    5. **Compensating Transaction:** If `CreditActivity` fails (e.g., destination account closed), the workflow must call `RefundDebitActivity` to reverse the debit on the `FromAccount`.

### 3.2. Activities
- `ValidateAccounts(ctx, req)`: Checks account status and returns error if invalid.
- `DebitAccount(ctx, accountID, amount, transferID)`: Performs the debit. Must be idempotent using `transferID`.
- `CreditAccount(ctx, accountID, amount, transferID)`: Performs the credit. Must be idempotent using `transferID`.
- `RefundDebitActivity(ctx, accountID, amount, transferID)`: Reverses a debit on the `FromAccount`. Used as compensation if `CreditActivity` fails. Must be idempotent using `transferID` as a natural key.

### 3.3. Technical Constraints
- **Determinism:** Workflow logic must be strictly deterministic.
- **Timeouts:** Use appropriate Activity and Workflow timeouts.
- **Retries:** Configure Retry Policies to handle transient failures (e.g., DB locks) but fail on business errors (e.g., Insufficient Funds).
- **Idempotency:** Two levels are required. (1) The workflow must be started with a deterministic `WorkflowID` (e.g. `"transfer-{transferID}"`) so that retrying the API call cannot create a duplicate workflow execution. (2) Each activity must use `transferID` as a natural idempotency key when writing to the database, so that Temporal activity retries do not produce duplicate ledger entries.

### 3.4. Data Integrity & Observability
- **Idempotency Persistence:** Activities use the `transferID` parameter as a natural key when writing to the database. The workflow is started with a deterministic `WorkflowID` at the API layer. Together these prevent both duplicate workflow executions and duplicate activity side-effects under retry.
- **Trace Propagation:** Ensure that the `context.Context` is correctly passed through all activities to maintain OpenTelemetry trace continuity from the API to the DB.

## 4. Acceptance Criteria
1.  **Happy Path:** Transfer < $1,000 completes automatically.
2.  **Approval Path:** Transfer > $1,000 waits for signal, then completes upon approval.
3.  **Rejection Path:** Transfer > $1,000 fails immediately upon `RejectSignal`.
4.  **Timeout Path:** Transfer > $1,000 fails if no signal received in time.
5.  **Compensation Path:** If `CreditActivity` fails, `RefundDebitActivity` is automatically called to reverse the debit on the `FromAccount`.
6.  **Idempotency:** Restarting a worker or re-running the same Workflow ID does not result in double debits.
