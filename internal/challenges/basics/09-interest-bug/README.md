# 🔍 Case 09: The Interest Bug ⭐ BONUS

## The Detective Brief
The interest calculator passed code review. The team wants to ship it. Your job: write the tests that prove it's broken before it hits production accounts.

## The Crime Scene
```bash
go test -v ./internal/challenges/basics/09-interest-bug/...
```

## Your Mission
1. Read `interest.go` — understand what `Calculate` should do.
2. Add at least **5 table cases** to `TestCalculate` in `interest_test.go`.
3. One of your cases must expose the hidden bug — watch it fail.
4. (Optional) Fix the bug in `interest.go` and watch your test go green.

**Note:** For this challenge you edit `interest_test.go`, not `challenge.go`.

## Key Lesson
> **Go testing:** Table-driven tests force you to enumerate edge cases you'd miss with ad-hoc tests. `t.Helper()` is already called in `checkResult` — notice how failure messages point to the table row, not the helper function.

<details>
<summary>Hints (click to reveal)</summary>

1. A normal case: `{name: "5% for 1 year", principal: 1000, rate: 0.05, years: 1, want: 1050.0}`
2. Zero years: any principal and rate → `want: principal` (no compounding)
3. Zero rate: any principal → `want: principal` (1+0 = 1, no growth)
4. Negative years: `wantErr: errors.New("years cannot be negative")` — but use `errors.Is` wisely
5. **The bug case**: `{name: "negative rate", principal: 1000, rate: -0.05, years: 1, want: 0, wantErr: ErrNegativeRate}` — add this and watch it fail!
</details>
