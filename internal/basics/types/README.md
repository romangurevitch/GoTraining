# 📦 Basic Types in Go

Go is a statically typed language, which means variables have specific types that are known at compile time.

## 1. Variable Declaration

Go provides multiple ways to declare variables:

```go
// 1. Explicit type, zero value initialization
var name string

// 2. Explicit type, with initial value
var age int = 30

// 3. Type inference (Go figures out the type)
var isStudent = false

// 4. Short variable declaration (only inside functions)
country := "Australia"
```

## 2. Best Practices

### Slices
- Prefer `make([]T, 0, capacity)` when you know the capacity. This avoids reallocations.
- Don't use pointer to slices `*[]T`. Slices are already small descriptors (pointer, length, capacity).

### Maps
- Initialize maps with `make(map[K]V, capacity)` if you know the approximate size.
- Reading from a map returns the zero-value if the key is missing. Use the "comma-ok" idiom to check for existence: `val, ok := m[key]`.

## 🏃 Running the Examples

Explore the unit tests for runnable patterns:
```bash
cd internal/basics/types
go test -v ./...
```
