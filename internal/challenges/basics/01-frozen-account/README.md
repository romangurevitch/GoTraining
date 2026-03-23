# 🔍 Case 01: The Frozen Account

## The Detective Brief
Alice deposits $100 into her account. The code compiles. The function runs without error. Her balance is still $0. The teller is baffled.

## The Crime Scene
```bash
go test -v ./internal/challenges/basics/01-frozen-account/...
```

## Your Mission
1. Run the tests — they fail.
2. Read `challenge.go` — the bug is in the method signature.
3. Fix the bug.
4. Run the tests again — they should pass.

## Key Lesson
> **Python vs Go:** In Python, `self` is always a reference to the actual object. In Go, `func (a Account) Deposit(...)` receives a *copy* of the struct — mutations to `a` vanish when the function returns. Use `func (a *Account) Deposit(...)` to mutate the original.

<details>
<summary>Hints (click to reveal)</summary>

1. Look at the receiver type in the `Deposit` method declaration.
2. In Go, `(a Account)` means the method gets a copy of the struct. `(a *Account)` means it gets a pointer to the original.
3. Change exactly one token in the method signature.
4. After the fix, `account.Deposit(100)` will update the balance of the original `account` variable.
</details>
