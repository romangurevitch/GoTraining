# 🔍 Case 04: Silent Errors

## The Detective Brief
The payment processor silently swallowed a failed withdrawal. No exception. No crash. No trace. The balance never changed. Compliance is not happy.

## The Crime Scene
```bash
cd internal/challenges/basics/04-silent-errors
```
```bash
go test -v ./...
```

## Your Mission
1. Run the tests — they panic (implement me).
2. Implement `Withdraw` in `challenge.go`.
3. Return the right sentinel error for each failure case.
4. Run the tests again — they should pass.

## Key Lesson
> **Python vs Go:** In Python you `raise ValueError(...)`. In Go you `return 0, ErrNegativeAmount`. Errors are values — callers must check them. Use `errors.Is(err, ErrX)` not `err.Error() == "some string"` — the latter breaks as soon as errors get wrapped.

<details>
<summary>Hints (click to reveal)</summary>

1. Check `amount <= 0` first and return `(0, ErrNegativeAmount)`.
2. Check `balance < amount` next and return `(0, ErrInsufficientFunds)`.
3. On success, return `(balance - amount, nil)`.
4. The test uses `errors.Is(err, ErrInsufficientFunds)` — your returned error must be (or wrap) the sentinel.
5. Never ignore the second return value from a Go function that returns `error`.
</details>

## Your Next Step
Errors are now loud and clear. But why does the compiler hate our new account types? It's time to learn the subtle art of interface satisfaction.
Head over to **[Case 05: The Fee Calculator](../05-fee-calculator/README.md)**.

