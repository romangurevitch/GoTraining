# 🔍 Case 05: The Fee Calculator

## The Detective Brief
Two account types need to report monthly fees. One compiles. The other won't budge. The interface is the same — why does Go refuse one of them?

## The Crime Scene
```bash
cd internal/challenges/basics/05-fee-calculator
```
```bash
go test -v ./...
```

## Your Mission
1. Run the tests — they panic (implement me).
2. Implement `MonthlyFee` for `SavingsAccount`.
3. Implement `MonthlyFee` for `PremiumAccount`.
4. Observe the compiler errors: why does one work with a value receiver and the other require a pointer?
5. Run the tests again once they compile — they should pass.

## Key Lesson
> **Python vs Go:** Python uses duck typing at runtime — if it has the method, it works. Go checks interface satisfaction at **compile time**. A crucial distinction in Go is the **Method Set**: if a method is declared on a pointer receiver (`*T`), then only the pointer type satisfies the interface. If you pass a value (`T`) where an interface is expected, it may fail to compile if the methods are bound to the pointer.

<details>
<summary>Hints (click to reveal)</summary>

1. Think about how to declare a method on `SavingsAccount` that operates on a copy.
2. Think about how to declare a method on `PremiumAccount` that operates on a pointer.
3. In the test, `SavingsAccount{}` (value) satisfies `FeeCalculator`. `&PremiumAccount{}` (pointer) satisfies it after your fix.
4. Try changing `*PremiumAccount` to `PremiumAccount` in the test — the compiler will reject it. That's the lesson.
5. `TotalFees` is already implemented — don't modify it.
</details>
