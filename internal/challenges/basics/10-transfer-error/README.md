# 🔍 Case 10: The Transfer Error ⭐ BONUS

## The Detective Brief
A transfer failed. The log says "error". Which account? How much? What rule? The compliance team needs answers, not just a string that says "error".

## The Crime Scene
```bash
go test -v ./internal/challenges/basics/10-transfer-error/...
```

## Your Mission
1. Run the tests — they panic (implement me).
2. Implement `TransferError.Error()` — a descriptive string.
3. Implement `Transfer()` — validate amount and account IDs, return `*TransferError` for failures.
4. Run the tests again.

## Key Lesson
> **Python vs Go:** Python exceptions carry arbitrary attributes. Go errors are values — define a struct that implements the `error` interface to carry structured context. Use `errors.As(err, &te)` to extract the concrete type, not type assertions (`err.(*TransferError)`). `errors.As` works through wrapped errors; direct type assertions don't.

<details>
<summary>Hints (click to reveal)</summary>

1. `TransferError.Error()` — use `fmt.Sprintf("transfer from %s to %s of $%.2f failed: %s", e.FromID, e.ToID, e.Amount, e.Reason)`.
2. `Transfer()` — check `amount <= 0` first, then `fromID == toID`.
3. Return `&TransferError{FromID: fromID, ToID: toID, Amount: amount, Reason: "amount must be positive"}` for invalid amount.
4. Return `&TransferError{FromID: fromID, ToID: toID, Amount: amount, Reason: "cannot transfer to same account"}` for same-account.
5. Return `nil` on success.
</details>
