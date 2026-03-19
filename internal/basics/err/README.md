# ⚠️ Error Handling in Go

Errors in Go are values, not exceptions. They are treated as first-class citizens and returned as regular return values. This approach encourages explicit error handling and makes the control flow highly predictable.

---

## 1. Core Concepts

| Concept | Description / Purpose |
| :--- | :--- |
| **`error` Interface** | The built-in interface `type error interface { Error() string }`. |
| **Sentinel Errors** | Predefined package-level error variables (e.g., `io.EOF`, `sql.ErrNoRows`). |
| **Custom Errors** | Structs that implement the `error` interface to attach extra context or fields. |
| **Wrapping (`%w`)** | Adding context to an error while preserving the original error type/value. |
| **`errors.Is`** | Checks if any error in the wrapped chain matches a specific sentinel error. |
| **`errors.As`** | Finds the first error in the chain that matches a specific type and assigns it. |
| **`errors.Join`** | (Go 1.20+) Combines multiple independent errors into a single error value. |
| **`panic` & `recover`** | Mechanisms for handling catastrophic, unrecoverable states. |

---

## 2. 🖼️ Visual Representation

Go errors form a "chain" when wrapped.

```text
  +--------------------------------+
  | fmt.Errorf("db error: %w", err)|   (Top-level / Current Error)
  +--------------------------------+
                  |
              (Unwrap)
                  |
                  v
  +--------------------------------+
  |   *ValidationError{...}        |   (Custom Error Type)
  +--------------------------------+
                  |
              (Unwrap)
                  |
                  v
  +--------------------------------+
  |       ErrNotFound              |   (Root / Sentinel Error)
  +--------------------------------+
```

When you use `errors.Is` or `errors.As`, Go traverses down this chain until it finds a match or reaches the root.

---

## 3. 📝 Implementation Examples

### Sentinel & Custom Errors

```go
// 1. Sentinel Error (usually defined at the package level)
var ErrNotFound = errors.New("item not found")

// 2. Custom Error Type
type ValidationError struct {
    Field  string
    Reason string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed on %s: %s", e.Field, e.Reason)
}
```

### Wrapping and Checking Errors

```go
func processUser() error {
    err := validateUser()
    if err != nil {
        // Wrap the error using %w
        return fmt.Errorf("processUser failed: %w", err)
    }
    return nil
}

func handle() {
    err := processUser()
    
    // Check for a Sentinel Error anywhere in the chain
    if errors.Is(err, ErrNotFound) {
        fmt.Println("Handle not found case")
    }

    // Extract a Custom Error type from anywhere in the chain
    var valErr *ValidationError
    if errors.As(err, &valErr) {
        fmt.Printf("Validation failed on field: %s\n", valErr.Field)
    }
}
```

---

## 4. 🚀 Common Patterns & Use Cases

- **Wrapping with Context**: Instead of just returning `err`, wrap it with `fmt.Errorf("failed to open config: %w", err)`. This creates a breadcrumb trail of where the error occurred.
- **Multiple Errors (`errors.Join`)**: When processing a batch of items or running multiple concurrent tasks, you can accumulate errors and return them all at once.
- **Panic and Recover**: Use `panic` ONLY for truly unrecoverable programming errors (e.g., out of bounds, nil pointer). Use `recover` in a `defer` block at boundary layers (like HTTP server middleware) to prevent an individual failing request from crashing the entire application.

---

## 5. ⚠️ Critical Pitfalls & Best Practices

> [!WARNING]
> Never ignore errors. If a function returns an `error`, handle it, return it, or explicitly discard it with `_` if you are absolutely sure it's safe (which is rare).

1. **Don't use `%v` for wrapping**: If you use `fmt.Errorf("... %v", err)`, the error is converted to a simple string. The chain is broken, and `errors.Is`/`errors.As` will fail. **Always use `%w`**.
2. **Sentinel Error Naming**: Prefix sentinel errors with `Err` (e.g., `ErrNotFound`, `ErrTimeout`).
3. **Custom Error Naming**: Suffix custom error structs with `Error` (e.g., `ValidationError`, `ParseError`).
4. **Don't Over-Panic**: Go is not Java. Don't use `panic` for normal control flow or validation failures. Stick to returning `error`.

---

## 🧪 Running the Examples

Explore the unit tests for runnable patterns covering Sentinel Errors, Custom Types, Wrapping, `errors.Join`, and `panic`/`recover`.

```bash
# Run tests with verbose output
go test -v ./internal/basics/err/...
```

---

## 📚 Further Reading

- [Official Go Blog: Error handling and Go](https://go.dev/blog/error-handling-and-go)
- [Official Go Package: errors](https://pkg.go.dev/errors)
