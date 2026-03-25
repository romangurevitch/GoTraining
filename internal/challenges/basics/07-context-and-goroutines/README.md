# 🔍 Case 07: Context and Goroutines

## The Detective Brief
A request arrived, but the context lost its tracking ID. Meanwhile, an expensive task was run synchronously, blocking the main thread!

## The Crime Scene
```bash
cd internal/challenges/basics/07-context-and-goroutines
```
```bash
go test -v -race ./...
```

## Your Mission
1. Implement `WithRequestID` to store a `requestID` in the context.
2. Implement `GetRequestID` to retrieve the `requestID` from the context.
3. Implement `RunAsync` to execute the provided task in a separate goroutine, so it doesn't block the caller.
4. Run the tests — they should pass.

## Key Lesson
> **Python vs Go:** In Python, passing request context often requires thread-locals or framework-specific features. Go uses `context.Context` to explicitly carry deadlines, cancellation signals, and request-scoped values across API boundaries. For concurrency, Go provides the `go` keyword to spawn lightweight goroutines instead of managing threads directly.

<details>
<summary>Hints (click to reveal)</summary>

1. Use `context.WithValue(ctx, key, value)`. You should define an unexported type for the key to avoid collisions (e.g. `type contextKey string`).
2. Use `ctx.Value(key)` to retrieve the value. You'll need to assert the type: `val, ok := ctx.Value(key).(string)`.
3. For `RunAsync`, just use `go task()`.
</details>

## Your Next Step
Concurrency mastered! Now, let's get back to basics and build something the whole bank can use. It's time to learn about package visibility and exporting.
Head over to **[Case 08: The Account Greeter](../08-account-greeter/README.md)**.

