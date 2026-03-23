# 🔍 Case 07: The Ignored Cancel

## The Detective Brief
The fraud monitor launched goroutines for every request. The requests completed. The goroutines are still running. Memory climbs. The context was cancelled. Nobody listened.

## The Crime Scene
```bash
cd internal/challenges/basics/07-ignored-cancel
```
```bash
go test -v -race -timeout 10s ./...
```

## Your Mission
1. Run the tests — the goroutine panics with "implement me".
2. Read `Watch` in `challenge.go`.
3. Implement the two TODO cases inside the `select`.
4. Run the tests again — they should pass.

## Key Lesson
> **Python vs Go:** Python's `asyncio` tasks can be cancelled automatically. Go goroutines run until they return — you must explicitly `select` on `ctx.Done()` to stop them. Always `defer cancel()` when you create a context: if you forget, goroutines can leak for the lifetime of the process.

<details>
<summary>Hints (click to reveal)</summary>

1. Replace `panic("implement me: return when context is cancelled")` with just `return`.
2. Replace `panic("implement me: call report on each iteration")` with `report(accountID)`.
3. `defer close(stopped)` is already there — it signals the test when the goroutine exits. Don't remove it.
4. Notice `defer cancel()` in the test. Always defer context cancellation — even if you expect it to be called anyway.
5. Run with `-race` to check for data races: `go test -race ./...`
</details>
