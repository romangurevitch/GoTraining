# 🔍 Case 05: The Fee Calculator

## The Detective Brief
Two account types need to report monthly fees. One compiles. The other won't budge. The interface is the same — why does Go refuse one of them?

## The Crime Scene
```bash
go test -v ./internal/challenges/basics/05-fee-calculator/...
```

## Your Mission
1. Run the tests — they panic (implement me).
2. Implement `MonthlyFee` for `SavingsAccount` with a **value** receiver.
3. Implement `MonthlyFee` for `PremiumAccount` with a **pointer** receiver.
4. Run the tests — they should pass.
5. Experiment: try a value receiver for `PremiumAccount`. What happens?

## Key Lesson
> **Python vs Go:** Python uses duck typing at runtime — if it has the method, it works. Go checks interface satisfaction at **compile time**. A pointer receiver method (`func (p *PremiumAccount) MonthlyFee()`) is only in `*PremiumAccount`'s method set, not `PremiumAccount`'s. Pass `&PremiumAccount{}`, not `PremiumAccount{}`.

<details>
<summary>Hints (click to reveal)</summary>

1. `SavingsAccount.MonthlyFee()` should use `(s SavingsAccount)` — value receiver. Return 5.0.
2. `PremiumAccount.MonthlyFee()` should use `(p *PremiumAccount)` — pointer receiver. Return 25.0.
3. In the test, `SavingsAccount{}` (value) satisfies `FeeCalculator`. `&PremiumAccount{}` (pointer) satisfies it after your fix.
4. Try changing `*PremiumAccount` to `PremiumAccount` in the test — the compiler will reject it. That's the lesson.
5. `TotalFees` is already implemented — don't modify it.
</details>
