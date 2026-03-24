# Temporal Challenge Fixes Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fix all 15 issues in the Durable Transfer Quest challenge files so the challenge is correct, consistent with the official Temporal Developer Skill, and ready to present.

**Architecture:** Pure documentation edits across 6 files — no code changes. Each task is isolated to one file. Changes are grouped by file to minimise context switching and make each commit reviewable on its own.

**Tech Stack:** Markdown, Mermaid diagrams, Git

**Spec:** `docs/superpowers/specs/2026-03-24-temporal-challenge-fixes-design.md`

> **Note for subagent:** All file edits must use content-based search (unique surrounding text), NOT line numbers. Line numbers are provided only as a rough orientation hint — they drift as earlier steps in the same file add lines.

---

## Files Modified

| File | Changes |
|---|---|
| `internal/challenges/temporal/README.md` | Quest 1 agent setup, Quest 2 scaffold note, Quest 3 label+activity name, Quest 4 AwaitWithTimeout+DisconnectedContext, Quest 5 activities+idempotency, Quest 6 port+WorkflowID+replay bonus, Quest 7 workflowcheck+idempotency check, remove Pro-Tip, fix architecture diagram |
| `internal/challenges/temporal/PRD.md` | Rejection semantics, compensation activity rename, RefundDebitActivity in 3.2, two-level idempotency in 3.3+3.4, Acceptance Criteria item 5 |
| `internal/challenges/temporal/GRADING_PROMPT.md` | Delete mandatory replay test bullet, add replay to Elite Bonus, update idempotency criterion, add self-audit bias note |
| `README.md` | Fix duplicate TOC numbering |
| `SCHEDULE.md` | Replace "Live Demo" framing with individual hands-on |
| `internal/temporal/README.md` | Update saga URL to compensating-transaction URL |

---

## Task 1: Fix `internal/temporal/README.md` and `README.md` (mechanical fixes)

**Files:**
- Modify: `internal/temporal/README.md`
- Modify: `README.md`

- [ ] **Step 1: Fix the saga URL in `internal/temporal/README.md`**

  Find the unique string:
  ```
  [Compensation pattern explained](https://microservices.io/patterns/data/saga.html)
  ```
  Replace with:
  ```
  [Compensation pattern explained](https://microservices.io/patterns/data/compensating-transaction.html)
  ```

- [ ] **Step 2: Fix the duplicate TOC numbering in `README.md`**

  Find the unique string (the second occurrence of `8.` in the TOC — the Further Learning line):
  ```
  8. [Further Learning & Resources](#further-learning--resources)
  ```
  Replace with:
  ```
  9. [Further Learning & Resources](#further-learning--resources)
  ```
  Note: the `8. [Challenges]` line immediately above must remain as `8.` — only the Further Learning line changes.

- [ ] **Step 3: Commit**

  ```bash
  git add internal/temporal/README.md README.md
  git commit -m "fix: correct compensation pattern URL and TOC numbering"
  ```

---

## Task 2: Fix `SCHEDULE.md`

**Files:**
- Modify: `SCHEDULE.md`

- [ ] **Step 1: Replace the 15:15 session description**

  Find the unique string:
  ```
  | 15:15 | **Agentic Go & Durable Workflows** | **Live Demo:** Leveraging AI agents and modern tooling to build a high-value transfer workflow with human-in-the-loop approval. | [Challenges](internal/challenges/temporal/README.md) |
  ```
  Replace with:
  ```
  | 15:15 | **Agentic Go & Durable Workflows** | **Hands-On Challenge:** Each participant builds a high-value transfer workflow with human-in-the-loop approval using their AI agent. Work through the quests in order — get as far as you can. Full quest guide remains available for self-study after the session. | [Challenges](internal/challenges/temporal/README.md) |
  ```

- [ ] **Step 2: Commit**

  ```bash
  git add SCHEDULE.md
  git commit -m "fix: reframe 15:15 session as individual hands-on challenge"
  ```

---

## Task 3: Fix `internal/challenges/temporal/PRD.md`

**Files:**
- Modify: `internal/challenges/temporal/PRD.md`

- [ ] **Step 1: Fix rejection semantics**

  Find the unique string (exact match including leading spaces):
  ```
        - If a `RejectSignal` is received, the workflow should terminate.
  ```
  Replace with:
  ```
        - If a `RejectSignal` is received, the workflow should return a `temporal.NewNonRetryableApplicationError("transfer rejected", "REJECTED", nil)`.
  ```

- [ ] **Step 2: Fix the compensation activity in step 5**

  Find the unique string:
  ```
      5. **Compensating Transaction:** If `CreditActivity` fails (e.g., destination account closed), the workflow must call `CreditActivity` (as a refund) back to the `FromAccount`.
  ```
  Replace with:
  ```
      5. **Compensating Transaction:** If `CreditActivity` fails (e.g., destination account closed), the workflow must call `RefundDebitActivity` to reverse the debit on the `FromAccount`.
  ```

- [ ] **Step 3: Add `RefundDebitActivity` to Section 3.2**

  Find the unique string:
  ```
  - `CreditAccount(ctx, accountID, amount, transferID)`: Performs the credit. Must be idempotent using `transferID`.
  ```
  Replace with:
  ```
  - `CreditAccount(ctx, accountID, amount, transferID)`: Performs the credit. Must be idempotent using `transferID`.
  - `RefundDebitActivity(ctx, accountID, amount, transferID)`: Reverses a debit on the `FromAccount`. Used as compensation if `CreditActivity` fails. Must be idempotent using `transferID` as a natural key.
  ```

- [ ] **Step 4: Add two-level idempotency constraint to Section 3.3**

  Find the unique string:
  ```
  - **Retries:** Configure Retry Policies to handle transient failures (e.g., DB locks) but fail on business errors (e.g., Insufficient Funds).
  ```
  Replace with:
  ```
  - **Retries:** Configure Retry Policies to handle transient failures (e.g., DB locks) but fail on business errors (e.g., Insufficient Funds).
  - **Idempotency:** Two levels are required. (1) The workflow must be started with a deterministic `WorkflowID` (e.g. `"transfer-{transferID}"`) so that retrying the API call cannot create a duplicate workflow execution. (2) Each activity must use `transferID` as a natural idempotency key when writing to the database, so that Temporal activity retries do not produce duplicate ledger entries.
  ```

- [ ] **Step 5: Replace Section 3.4 idempotency line**

  Find the unique string:
  ```
  - **Idempotency Persistence:** Activities must persist the `TransferID` to the database to ensure a "Restart" or "Retry" does not result in duplicate ledger entries.
  ```
  Replace with:
  ```
  - **Idempotency Persistence:** Activities use the `transferID` parameter as a natural key when writing to the database. The workflow is started with a deterministic `WorkflowID` at the API layer. Together these prevent both duplicate workflow executions and duplicate activity side-effects under retry.
  ```

- [ ] **Step 6: Fix Acceptance Criteria item 5**

  Find the unique string:
  ```
  5.  **Compensation Path:** If `Credit` fails, the `Debit` is automatically reversed via a refund.
  ```
  Replace with:
  ```
  5.  **Compensation Path:** If `CreditActivity` fails, `RefundDebitActivity` is automatically called to reverse the debit on the `FromAccount`.
  ```

- [ ] **Step 7: Commit**

  ```bash
  git add internal/challenges/temporal/PRD.md
  git commit -m "fix: introduce RefundDebitActivity, fix rejection semantics, clarify two-level idempotency"
  ```

---

## Task 4: Fix `internal/challenges/temporal/GRADING_PROMPT.md`

**Files:**
- Modify: `internal/challenges/temporal/GRADING_PROMPT.md`

- [ ] **Step 1: Add self-audit bias note under Persona section**

  Find the unique string:
  ```
  **Objective:** Perform a rigorous, holistic audit of the "Durable Transfer Workflow." This is a competition: evaluate the participant's ability to maintain architectural integrity within the **specific context** of this Go Bank codebase.
  ```
  Replace with:
  ```
  **Objective:** Perform a rigorous, holistic audit of the "Durable Transfer Workflow." This is a competition: evaluate the participant's ability to maintain architectural integrity within the **specific context** of this Go Bank codebase.

  > **Note:** This audit is most valuable when run in a fresh agent session without the implementation context. If you built the code yourself, bias toward stricter scoring on pattern consistency and agentic noise.
  ```

- [ ] **Step 2: Delete the mandatory Replay Test bullet from Section 2**

  Find the unique string:
  ```
  *   **The Replay Test:** **MANDATORY.** A failure to include a dedicated history replay test (ensuring current code can handle existing histories) is a disqualifying omission for a top-tier score.
  ```
  Replace with nothing (delete the entire line).

- [ ] **Step 3: Update the idempotency criterion in Section 3**

  Find the unique string:
  ```
  *   **Idempotency:** Is the `TransferID` used surgically to ensure no transaction can be duplicated under Temporal retry conditions?
  ```
  Replace with:
  ```
  *   **Idempotency:** Is the workflow started with a deterministic `WorkflowID` (e.g. `"transfer-{transferID}"`) to prevent duplicate workflow executions? Are the activities (`DebitAccount`, `CreditAccount`, `RefundDebitActivity`) also idempotent with respect to Temporal activity retries — i.e., does each activity use `transferID` as a natural key to prevent duplicate DB writes?
  ```

- [ ] **Step 4: Add Replay Test to Elite Bonus section**

  Find the unique string:
  ```
  Award points only for **proactive engineering**—decisions that were not in the PRD but improve the system (e.g., custom OTel spans, Advanced status queries, or sophisticated error-recovery strategies).
  ```
  Replace with:
  ```
  Award points only for **proactive engineering**—decisions that were not in the PRD but improve the system (e.g., custom OTel spans, Advanced status queries, or sophisticated error-recovery strategies).
  - **Replay Test (+5 pts):** A `Test_ReplayHistory` test using `worker.NewWorkflowReplayer()` against an exported workflow history JSON file. Demonstrates production-grade safety discipline.
  ```

- [ ] **Step 5: Commit**

  ```bash
  git add internal/challenges/temporal/GRADING_PROMPT.md
  git commit -m "fix: demote replay test to bonus, update idempotency criterion, add self-audit bias note"
  ```

---

## Task 5: Fix `internal/challenges/temporal/README.md` — Architecture diagram and Quest 1

**Files:**
- Modify: `internal/challenges/temporal/README.md`

- [ ] **Step 1: Fix the architecture diagram**

  Find the unique string (note: 8 spaces of indentation before `W->>A`):
  ```
          W->>A: COMPENSATION: Refund Debit
  ```
  Replace with:
  ```
          W->>A: Activity: RefundDebitActivity
  ```

- [ ] **Step 2: Fix Quest 1 Step 2 — remove tool-specific filenames**

  Find the unique string (note trailing space after `project.`):
  ```
  2.  **Create your Rules:** Create a file (e.g., `.cursorrules`, `.aierules`, `.github/copilot-instructions.md`, or `GEMINI.md`) in the root of the project.
  ```
  Replace with:
  ```
  2.  **Create your agent instruction file:** Create a configuration file for your AI agent per your tool's conventions. Refer to your tool's documentation for the correct filename and location.
  ```

- [ ] **Step 3: Commit**

  ```bash
  git add internal/challenges/temporal/README.md
  git commit -m "fix: update architecture diagram activity name and make Quest 1 agent-agnostic"
  ```

---

## Task 6: Fix `internal/challenges/temporal/README.md` — Quest 2 and Quest 3

**Files:**
- Modify: `internal/challenges/temporal/README.md`

- [ ] **Step 1: Add scaffold note at the start of Quest 2**

  Find the unique string (note trailing space):
  ```
  Now that your agent is "primed," work together to agree on the contract.
  ```
  Replace with:
  ```
  Now that your agent is "primed," work together to agree on the contract.

  > **Note:** `internal/bank/temporal/` does not exist yet. Create it as a new Go package as part of this quest. Your `spec.md` will live at `internal/bank/temporal/spec.md`.
  ```

- [ ] **Step 2: Fix Quest 3 — replace "Workflow Path" bullet**

  Find the unique string:
  ```
      - **Workflow Path:** Simulate a failure in `CreditActivity` and verify `DebitActivity` is compensated (called as a refund).
  ```
  Replace with:
  ```
      - **Compensation Path:** Simulate a failure in `CreditActivity` and verify `RefundDebitActivity` is called as compensation.
  ```

- [ ] **Step 3: Commit**

  ```bash
  git add internal/challenges/temporal/README.md
  git commit -m "fix: add scaffold note for Quest 2, rename Workflow Path to Compensation Path in Quest 3"
  ```

---

## Task 7: Fix `internal/challenges/temporal/README.md` — Quest 4

**Files:**
- Modify: `internal/challenges/temporal/README.md`

- [ ] **Step 1: Replace `workflow.Selector` with `workflow.AwaitWithTimeout`**

  Find the unique string:
  ```
  - Use a `workflow.Selector` for the approval gate and timeout.
  ```
  Replace with:
  ```
  - Use `workflow.AwaitWithTimeout` for the approval gate — it blocks until the `ApprovalSignal` arrives or the 24-hour timeout expires, then returns `(conditionMet bool, err error)`.
  - Use `workflow.NewDisconnectedContext(ctx)` inside the compensation branch (when `CreditActivity` returns an error), before calling `RefundDebitActivity`, to obtain a context not tied to the parent workflow's cancellation. This ensures compensation runs even if the workflow is externally cancelled.
  ```

- [ ] **Step 2: Commit**

  ```bash
  git add internal/challenges/temporal/README.md
  git commit -m "fix: replace workflow.Selector with AwaitWithTimeout, add NewDisconnectedContext guidance"
  ```

---

## Task 8: Fix `internal/challenges/temporal/README.md` — Quest 5

**Files:**
- Modify: `internal/challenges/temporal/README.md`

- [ ] **Step 1: Update the Task section — add RefundDebitActivity and fix idempotency instruction**

  Find the unique string:
  ```
  - Implement `DebitAccount` and `CreditAccount` in `internal/bank/temporal/activities.go`.
  - Ensure they are **idempotent** by using the `TransferID` as a unique key for the transaction.
  - Map domain errors (e.g., `InsufficientFunds`) to `temporal.NewNonRetryableApplicationError`.
  ```
  Replace with:
  ```
  - Implement `DebitAccount`, `CreditAccount`, and `RefundDebitActivity` in `internal/bank/temporal/activities.go`.
  - Each activity receives a `transferID` parameter. Use it as a natural idempotency key when writing to the database — check for an existing record with that `transferID` before inserting, or use an `ON CONFLICT DO NOTHING` strategy. This prevents duplicate ledger entries if Temporal retries the activity.
  - Map domain errors (e.g., `domain.ErrInsufficientFunds`) to `temporal.NewNonRetryableApplicationError` so Temporal does not retry business logic failures.
  ```

- [ ] **Step 2: Update the Definition of Done**

  Find the unique string:
  ```
  - Run the activities against a mock repository or the real database.
  ```
  Replace with:
  ```
  - Each activity handles a Temporal retry correctly: running the same activity twice with the same `transferID` produces exactly one ledger entry.
  ```

- [ ] **Step 3: Commit**

  ```bash
  git add internal/challenges/temporal/README.md
  git commit -m "fix: add RefundDebitActivity to Quest 5, clarify activity-level idempotency"
  ```

---

## Task 9: Fix `internal/challenges/temporal/README.md` — Quest 6

**Files:**
- Modify: `internal/challenges/temporal/README.md`

- [ ] **Step 1: Fix the Temporal Web UI port**

  Find the unique string:
  ```
  - You can observe the workflow history in the [Temporal Web UI](http://localhost:8080).
  ```
  Replace with:
  ```
  - You can observe the workflow history in the [Temporal Web UI](http://localhost:8233).
  ```

- [ ] **Step 2: Add WorkflowID instruction to the Quest 6 Task section**

  Find the unique string:
  ```
  - Wire up a new API endpoint `POST /v1/durable-transfers` that starts the Temporal workflow.
  ```
  Replace with:
  ```
  - Wire up a new API endpoint `POST /v1/durable-transfers` that starts the Temporal workflow.
  - Start the workflow with a deterministic `WorkflowID` — e.g. `"transfer-" + req.TransferID`. Temporal guarantees that a second `StartWorkflow` call with the same `WorkflowID` returns the existing execution rather than creating a duplicate.
  ```

- [ ] **Step 3: Add the Export & Replay Test bonus sub-task after the Definition of Done**

  Find the unique string:
  ```
  - The funds are correctly moved only after approval.
  ```
  Replace with:
  ````
  - The funds are correctly moved only after approval.

  #### Export & Replay Test (Bonus)
  After a successful E2E run, export the workflow history:
  ```bash
  temporal workflow show --workflow-id <your-workflow-id> --output json > internal/bank/temporal/testdata/transfer_history.json
  ```
  Then write a `Test_ReplayHistory` test using `worker.NewWorkflowReplayer()` to verify your code handles existing histories correctly — the standard safety check before deploying changes to a live workflow.
  ````

- [ ] **Step 4: Commit**

  ```bash
  git add internal/challenges/temporal/README.md
  git commit -m "fix: correct Temporal UI port, add WorkflowID instruction, add replay test bonus in Quest 6"
  ```

---

## Task 10: Fix `internal/challenges/temporal/README.md` — Quest 7 and Pro-Tips

**Files:**
- Modify: `internal/challenges/temporal/README.md`

- [ ] **Step 1: Add `workflowcheck` task to Quest 7**

  Find the unique string:
  ```
  - **Worker Registration:** Register your new Workflow and Activities in `cmd/temporal/worker/main.go`.
  ```
  Replace with:
  ```
  - **Worker Registration:** Register your new Workflow and Activities in `cmd/temporal/worker/main.go`.
  - **Static Analysis:** Run `workflowcheck ./...` to catch non-deterministic code in your workflow (accidental use of `time.Now()`, native goroutines, or native `select`). Fix any violations — this is a mandatory gate in production CI pipelines.
  ```

- [ ] **Step 2: Update the Idempotency Check sub-task in Quest 7**

  Find the unique string:
  ```
  - **Idempotency Check:** Manually "kill" your worker during an activity execution and verify that upon restart, the `TransferID` prevents a duplicate charge.
  ```
  Replace with:
  ```
  - **Idempotency Check:** Manually kill your worker during an activity execution. Upon restart, verify that no duplicate ledger entries exist — the activity's `transferID` check should have prevented the duplicate write.
  ```

- [ ] **Step 3: Remove the "On Idempotency" Pro-Tip section**

  Find the unique string:
  ```
  ### On Idempotency
  To get 100/100, your activities shouldn't just "be idempotent"; they should prove it. Consider adding a migration (`/migration`) to create a `deduplication_events` table or use a `UNIQUE` constraint on the `TransferID` in your `transactions` table.

  ```
  Replace with nothing (delete entirely, including the trailing blank line).

- [ ] **Step 4: Commit**

  ```bash
  git add internal/challenges/temporal/README.md
  git commit -m "fix: add workflowcheck to Quest 7, update idempotency check wording, remove outdated pro-tip"
  ```

---

## Task 11: Verify all changes

- [ ] **Step 1: Confirm only the expected files were changed**

  ```bash
  git log main..HEAD --oneline | wc -l
  ```
  Expected output: `10` (one commit per task).

  ```bash
  git diff main --name-only
  ```
  Expected output — exactly these 6 files, no others:
  ```
  README.md
  SCHEDULE.md
  internal/challenges/temporal/GRADING_PROMPT.md
  internal/challenges/temporal/PRD.md
  internal/challenges/temporal/README.md
  internal/temporal/README.md
  ```

- [ ] **Step 2: Verify no stale references remain**

  ```bash
  grep -rn "localhost:8080" internal/challenges/temporal/
  grep -rn "workflow\.Selector" internal/challenges/temporal/
  grep -rn "deduplication_events\|UNIQUE constraint" internal/challenges/temporal/
  grep -rn "the workflow should terminate" internal/challenges/temporal/
  grep -rn "COMPENSATION: Refund Debit" internal/challenges/temporal/
  grep -rn "\.cursorrules\|\.aierules\|copilot-instructions" internal/challenges/temporal/
  grep -rn "Workflow Path" internal/challenges/temporal/
  grep -rn "saga\.html" internal/temporal/
  grep -n "Live Demo" SCHEDULE.md
  grep -n "MANDATORY" internal/challenges/temporal/GRADING_PROMPT.md
  ```
  Expected: all greps return empty output (no matches).

- [ ] **Step 3: Verify key new content is present**

  ```bash
  grep -n "RefundDebitActivity" internal/challenges/temporal/README.md internal/challenges/temporal/PRD.md
  grep -n "AwaitWithTimeout" internal/challenges/temporal/README.md
  grep -n "NewDisconnectedContext" internal/challenges/temporal/README.md
  grep -n "workflowcheck" internal/challenges/temporal/README.md
  grep -n "transfer-{transferID}\|transfer-.*TransferID" internal/challenges/temporal/PRD.md
  grep -n "Hands-On Challenge" SCHEDULE.md
  grep -n "compensating-transaction" internal/temporal/README.md
  grep -n "fresh agent session" internal/challenges/temporal/GRADING_PROMPT.md
  grep -n "NewWorkflowReplayer" internal/challenges/temporal/GRADING_PROMPT.md
  grep -n "Compensation Path" internal/challenges/temporal/README.md
  ```
  Expected: each grep returns at least one match.

- [ ] **Step 4: Confirm working tree is clean**

  ```bash
  git status
  ```
  Expected: `nothing to commit, working tree clean`.
