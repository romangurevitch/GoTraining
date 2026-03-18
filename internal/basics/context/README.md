# Context

`context.Context` carries deadlines, cancellation signals, and request-scoped values.

## Creating Contexts

```go
ctx, cancel := context.WithCancel(context.Background())
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel() // always defer cancel
```

## Pitfalls

- Always use a typed, unexported key for context values
- Always `defer cancel()` — forgetting it leaks goroutines
- Pass context as the first parameter in every I/O function
- Never store a context in a struct
