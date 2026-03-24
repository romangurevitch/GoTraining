# Design: Temporal Challenge Fixes

**Date:** 2026-03-24
**Branch:** remove-saga-terminology
**Scope:** Fixes to `internal/challenges/temporal/` (README.md, PRD.md, GRADING_PROMPT.md), top-level `README.md`, `SCHEDULE.md`, and `internal/temporal/README.md`

---

## Context

A review of the Durable Transfer Quest challenge (and associated modified files on the `remove-saga-terminology` branch) against the official Temporal Developer Skill identified 15 issues spanning broken links, wrong ports, structural dead-ends, and content inconsistencies. This document captures the agreed design for fixing all of them.

Mechanical issues (1, 6, 7, 12, 13) require no design decision — only the correct value applied. Issues requiring design decisions are captured in the Decisions Log.

---

## Section 1: Mechanical Fixes

### Files affected
- `internal/challenges/temporal/README.md`
- `README.md` (top-level)
- `internal/temporal/README.md`

### Changes

| File | Location | Change |
|---|---|---|
| `internal/challenges/temporal/README.md` | Quest 6 | `http://localhost:8080` → `http://localhost:8233` (correct Temporal Web UI port) |
| `internal/challenges/temporal/README.md` | Quest 3 | "Workflow Path" label → "Compensation Path" (replace entire bullet — see Section 3 Quest 3) |
| `internal/challenges/temporal/README.md` | Quest 4 | `workflow.Selector` → `workflow.AwaitWithTimeout` (see Section 3 Quest 4) |
| `internal/challenges/temporal/README.md` | Quest 1 | Remove specific agent filenames list; replace with generic instruction (see Section 3 Quest 1) |
| `internal/challenges/temporal/README.md` | Architecture diagram | Update `COMPENSATION: Refund Debit` → `Activity: RefundDebitActivity` in Mermaid sequence diagram |
| `README.md` | Table of Contents | Second `8.` → `9.` (Further Learning & Resources) |
| `internal/temporal/README.md` | Further Reading | Update URL: `https://microservices.io/patterns/data/saga.html` → `https://microservices.io/patterns/data/compensating-transaction.html` |

---

## Section 2: PRD Fixes (`PRD.md`)

### Idempotency approach (Sections 3.1, 3.2, 3.3, 3.4)

**Two levels of idempotency are required and must be clearly distinguished:**

1. **Workflow-level idempotency** (prevents duplicate workflow executions): The API must start the workflow with a deterministic `WorkflowID` (e.g. `"transfer-{transferID}"`). Temporal guarantees a second `StartWorkflow` call with the same `WorkflowID` returns the existing execution rather than starting a new one.

2. **Activity-level idempotency** (prevents duplicate DB writes within a single execution): Each activity (`DebitAccount`, `CreditAccount`, `RefundDebitActivity`) receives a `transferID` parameter. Activities must use this `transferID` as a natural idempotency key when writing to the database (e.g. check for an existing transaction with that `transferID` before inserting, or use an `ON CONFLICT DO NOTHING` / `INSERT IF NOT EXISTS` pattern). Temporal may retry an activity multiple times within the same workflow execution.

**Changes to PRD:**

- **Add to Section 3.3 (Technical Constraints):**
  > "**Idempotency:** Two levels are required. (1) The workflow must be started with a deterministic `WorkflowID` (e.g. `"transfer-{transferID}"`) so that retrying the API call cannot create a duplicate workflow execution. (2) Each activity must use `transferID` as a natural idempotency key when writing to the database, so that Temporal activity retries do not produce duplicate ledger entries."

- **Replace Section 3.4 (Data Integrity & Observability):** Remove the `deduplication_events` table and `UNIQUE` constraint recommendations. Replace with:
  > "**Idempotency Persistence:** Activities use the `transferID` parameter as a natural key when writing to the database. The workflow is started with a deterministic `WorkflowID` at the API layer. Together these prevent both duplicate workflow executions and duplicate activity side-effects under retry."

### Compensation activity (Section 3.1 step 5, Section 3.2, Acceptance Criteria item 5)

**Section 3.1 step 5 — replace the entire bullet with:**
> "**Compensating Transaction:** If `CreditActivity` fails (e.g., destination account closed), the workflow must call `RefundDebitActivity` to reverse the debit on the `FromAccount`."

**Section 3.2 — add new activity entry:**
> `RefundDebitActivity(ctx, accountID, amount, transferID)`: Reverses a debit on the `FromAccount`. Used as compensation if `CreditActivity` fails. Must be idempotent using `transferID` as a natural key.

**Acceptance Criteria item 5 — replace with:**
> "**Compensation Path:** If `CreditActivity` fails, `RefundDebitActivity` is automatically called to reverse the debit on the `FromAccount`."

### Rejection semantics (Section 3.1 step 2)

Replace the third sub-bullet of step 2:

**Before:**
> `- If a RejectSignal is received, the workflow should terminate.`

**After:**
> `- If a RejectSignal is received, the workflow should return a temporal.NewNonRetryableApplicationError("transfer rejected", "REJECTED", nil).`

---

## Section 3: Challenge README Quest Fixes (`internal/challenges/temporal/README.md`)

### Quest 1 — Agent Setup

**Step 2 — replace entire bullet with:**
> "**Create your agent instruction file:** Create a configuration file for your AI agent per your tool's conventions. Refer to your tool's documentation for the correct filename and location."

**Step 3 — keep as-is** (the content of the instruction is tool-agnostic already).

### Quest 2 — Spec & Design

Add a note at the start of the quest:
> "**Note:** `internal/bank/temporal/` does not exist yet. Create it as a new Go package as part of this quest. Your `spec.md` will live at `internal/bank/temporal/spec.md`."

### Quest 3 — TDD

Replace the "Workflow Path" bullet entirely with:
> "- **Compensation Path:** Simulate a failure in `CreditActivity` and verify `RefundDebitActivity` is called as compensation."

### Quest 4 — Implement the Workflow

**Replace this bullet:**
> "- Use a `workflow.Selector` for the approval gate and timeout."

**With:**
> "- Use `workflow.AwaitWithTimeout` for the approval gate — it blocks until the `ApprovalSignal` arrives or the 24-hour timeout expires, then returns `(conditionMet bool, err error)`."

**Add new bullet immediately after:**
> "- Use `workflow.NewDisconnectedContext(ctx)` to obtain a new context before calling `RefundDebitActivity`. Call this inside the compensation branch (when `CreditActivity` returns an error), *before* executing the compensation activity. This ensures the compensation runs even if the workflow's parent context has been externally cancelled."

### Quest 5 — Activities

**Replace the entire "Definition of Done" and idempotency instruction with:**

Task bullet update:
> "- Implement `DebitAccount`, `CreditAccount`, and `RefundDebitActivity` in `internal/bank/temporal/activities.go`."
> "- Each activity receives a `transferID` parameter. Use it as a natural idempotency key when writing to the database — check for an existing record with that `transferID` before inserting, or use an `ON CONFLICT DO NOTHING` strategy. This prevents duplicate ledger entries if Temporal retries the activity."
> "- Map domain errors (e.g., `domain.ErrInsufficientFunds`) to `temporal.NewNonRetryableApplicationError` so Temporal does not retry business logic failures."

**Definition of Done update:**
> "- Each activity handles a Temporal retry correctly: running the same activity twice with the same `transferID` produces exactly one ledger entry."

Note: the WorkflowID idempotency instruction (preventing duplicate *workflow executions*) belongs in Quest 6 when participants wire the API endpoint. Do not instruct it here.

### Quest 6 — End-to-End

- Fix Temporal Web UI link: `localhost:8080` → `localhost:8233`

- **Add to the API wiring task:**
  > "- Start the workflow with a deterministic `WorkflowID` — e.g. `"transfer-" + req.TransferID`. Temporal guarantees that a second `StartWorkflow` call with the same `WorkflowID` returns the existing execution rather than creating a duplicate."

- **Add new sub-task after the manual E2E test:**
  > "**Export & Replay Test (Bonus):** After a successful E2E run, export the workflow history:
  > ```bash
  > temporal workflow show --workflow-id <your-workflow-id> --output json > internal/bank/temporal/testdata/transfer_history.json
  > ```
  > Then write a `Test_ReplayHistory` test using `worker.NewWorkflowReplayer()` to verify your code handles existing histories correctly — the standard safety check before deploying changes to a live workflow."

### Quest 7 — Production Hardening

**Add new task:**
> "**Static Analysis:** Run `workflowcheck ./...` to catch non-deterministic code in your workflow (accidental use of `time.Now()`, native goroutines, or native `select`). Fix any violations — this is a mandatory gate in production CI pipelines."

**Update the Idempotency Check sub-task** (currently: "verify that the `TransferID` prevents a duplicate charge"). Replace with:
> "**Idempotency Check:** Manually kill your worker during an activity execution. Upon restart, verify that no duplicate ledger entries exist — the activity's `transferID` check should have prevented the duplicate write."

### Engineering Pro-Tips

- **Remove** the "On Idempotency" pro-tip (the `deduplication_events` table approach is no longer used).
- Keep the "On Observability" pro-tip unchanged.

---

## Section 4: Grading & Schedule Fixes

### `GRADING_PROMPT.md`

**Replay Test bullet in Section 2 (Testing Excellence) — delete it entirely.**
Replace the current mandatory replay test bullet with nothing. The replay test moves to the Elite Bonus section.

**Elite Bonus section — add:**
> "**Replay Test (+5 pts):** A `Test_ReplayHistory` test using `worker.NewWorkflowReplayer()` against an exported workflow history JSON file. Demonstrates production-grade safety discipline."

**Idempotency criterion in Section 3 (Pattern Consistency) — replace:**

Before:
> "Is the `TransferID` used surgically to ensure no transaction can be duplicated under Temporal retry conditions?"

After:
> "Is the workflow started with a deterministic `WorkflowID` (e.g. `"transfer-{transferID}"`) to prevent duplicate workflow executions? Are the activities (`DebitAccount`, `CreditAccount`, `RefundDebitActivity`) also idempotent with respect to Temporal activity retries — i.e., does each activity use `transferID` as a natural key to prevent duplicate DB writes?"

**Self-audit bias note — add at the top under the Persona section:**
> "**Note:** This audit is most valuable when run in a fresh agent session without the implementation context. If you built the code yourself, bias toward stricter scoring on pattern consistency and agentic noise."

### `SCHEDULE.md`

Replace the `15:15` row Description column content. Remove the "**Live Demo:**" prefix entirely. New full description:
> "**Hands-On Challenge:** Each participant builds a high-value transfer workflow with human-in-the-loop approval using their AI agent. Work through the quests in order — get as far as you can. Full quest guide remains available for self-study after the session."

The Topic column ("Agentic Go & Durable Workflows") remains unchanged.

---

## Decisions Log

| Issue | Decision | Rationale |
|---|---|---|
| Issues 1, 7, 12 | Mechanical fix — no design decision | Wrong port, wrong label, broken TOC numbering |
| Issues 6, 13 | Mechanical fix — no design decision | PRD self-contradiction resolved by compensation rename; link URL updated to correct destination |
| Issue 2: `internal/bank/temporal/` missing scaffold | Add prose note in Quest 2 directing participants to create the package | Clear enough for participants without pre-creating stub files |
| Issue 3: DB-level idempotency impossible with existing Repository | Two-level idempotency: WorkflowID for workflow-level; transferID natural key for activity-level | Teaches both Temporal patterns correctly; no Repository interface changes needed |
| Issue 4: `CreditActivity` reused as compensation | Introduce `RefundDebitActivity`; update PRD step 5, Section 3.2, Acceptance Criteria item 5, and architecture diagram | Matches Temporal skill saga pattern; removes forward/reverse naming confusion |
| Issue 5: "terminate" for rejection wrong Temporal semantics | Return `temporal.NewNonRetryableApplicationError("transfer rejected", "REJECTED", nil)` | Correct SDK usage; gives implementers a concrete, consistent error shape |
| Issue 8: Replay test mandatory but impossible at Quest 3 | Add as Quest 6 bonus sub-task; delete from mandatory in Section 2 grading rubric; add to Elite Bonus | Achievable sequencing; fair rubric for 45-min session |
| Issue 9: `workflowcheck` never mentioned | Add to Quest 7 Production Hardening | Production readiness check, not a development hint |
| Issue 10: `workflow.Selector` for approval gate | Replace with `workflow.AwaitWithTimeout` | Idiomatic modern API; session is time-pressured |
| Issue 11: Tool-specific agent filenames in Quest 1 | Drop filename list; refer participants to their tool's conventions | Agent-agnostic; zero maintenance burden |
| Issue 14: SCHEDULE.md implied trainer demo | Reframe as individual hands-on challenge; remove "**Live Demo:**" prefix | Accurate; participants work independently |
| Issue 15: Grading self-administered by same agent | Add bias note to GRADING_PROMPT.md Persona section | Acknowledges the limitation; encourages stricter self-scoring |
