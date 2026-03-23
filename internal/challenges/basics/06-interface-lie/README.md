# 🔍 Case 06: The Interface Lie

## The Detective Brief
The AML checker says it found no error. But the error handler fires anyway. The checker was skipped. The variable is nil. Yet `err != nil` is true. This is Go's most famous gotcha.

## The Crime Scene
```bash
go test -v ./internal/challenges/basics/06-interface-lie/...
```

## Your Mission
1. Run the tests — they fail, even though the "error" prints as `<nil>`.
2. Read `runAMLCheck` in `challenge.go`.
3. Understand why returning a typed nil pointer as `error` isn't the same as returning `nil`.
4. Fix it with a one-word change.

## Key Lesson
> **Python vs Go:** A Go interface holds two things: `(type, value)`. When you return `(*AMLChecker)(nil)` as an `error`, the interface is `(type=*AMLChecker, value=nil)` — the type is set, so the interface is NOT nil. Always return bare `nil` from functions returning interface types, never a typed nil variable.

<details>
<summary>Hints (click to reveal)</summary>

1. `var checker *AMLChecker` declares a nil pointer of type `*AMLChecker`.
2. Storing that nil pointer in an `error` interface creates `(type=*AMLChecker, value=nil)`.
3. That interface != nil, because the type field is set — even though the value prints as `<nil>`.
4. Fix: in the `if skip` branch, write `return nil` — not `return checker`.
5. The `Explain()` function in `challenge.go` can help you see what the interface actually holds.
</details>
