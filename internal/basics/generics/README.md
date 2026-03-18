# Generics

Go 1.18+ supports type parameters for reusable, type-safe utilities.

## Type Parameters

```go
func Map[T, U any](s []T, f func(T) U) []U {
    result := make([]U, len(s))
    for i, v := range s { result[i] = f(v) }
    return result
}
```

## Constraints

```go
type Number interface { int | int64 | float64 }
func Sum[T Number](values []T) T { ... }
```

## Pitfalls

- Only use generics when the same logic applies to truly different types
- You cannot use operators on unconstrained `any` type parameters
