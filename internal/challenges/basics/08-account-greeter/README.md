# 🔍 Case 08: The Account Greeter

## The Detective Brief
Build the greeting service — the simplest package CBA ships. It must be correct, exported correctly, and testable. One exported function. One unexported helper. No excuses.

## The Crime Scene
```bash
go test -v ./internal/challenges/basics/08-account-greeter/...
```

## Your Mission
1. Run the tests — they panic (implement me).
2. Implement `formatName` (unexported helper) in `greeting.go`.
3. Implement `Greet` (exported function) using `formatName` and `fmt.Sprintf`.
4. Run the tests again.

**Note:** The editable file here is `greeting.go`, not `challenge.go`. This mirrors how real Go packages name their files after what they contain.

## Key Lesson
> **Python vs Go:** Python uses `_name` convention for private members — it's advisory, not enforced. In Go, a lowercase first letter (`formatName`) means package-private and the compiler enforces it. An uppercase first letter (`Greet`) means exported and callable from any package. There are no other access modifiers.

<details>
<summary>Hints (click to reveal)</summary>

1. `formatName(first, last string) string` — return `first + " " + last`.
2. `Greet(accountID, first, last string) string` — use `fmt.Sprintf("Hello, %s! Your account %s is ready.", formatName(first, last), accountID)`.
3. Try calling `formatName` from a different package — the compiler will reject it.
4. The test file is in the same package (`package accountgreeter`) so it can test unexported behaviour indirectly.
</details>
