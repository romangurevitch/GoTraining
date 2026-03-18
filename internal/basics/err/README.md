# Error Handling

Go errors are values — not exceptions. Functions signal failure by returning an `error` as the last return value.

## Creating Errors

```go
errors.New("something went wrong")
fmt.Errorf("loading config: %w", err) // %w wraps the error
```

## Wrapping and Unwrapping

```go
var ErrNotFound = errors.New("not found")
if errors.Is(err, ErrNotFound) { ... }

var target *MyError
if errors.As(err, &target) { ... }
```

## Pitfalls

- Never discard errors with `_`
- Use `%w` (not `%v`) to preserve the error chain for `errors.Is`/`errors.As`
