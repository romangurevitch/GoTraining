# Generics in Go

Generics (introduced in Go 1.18) allow you to write code that works with multiple types while maintaining type safety. This eliminates the need for `interface{}`-based casting or repetitive code for different numeric types.

---

## 1. Core Concepts

| Concept | Description / Purpose |
| :--- | :--- |
| **Type Parameters** | The list of placeholder types in square brackets (e.g., `[T any]`). |
| **`any`** | An alias for the empty interface `interface{}`. Used when any type is acceptable. |
| **`comparable`** | A built-in constraint for types that support `==` and `!=`. Required for map keys or equality checks. |
| **Type Constraints** | Interfaces that define the allowed set of types (e.g., `Number interface { int \| float64 }`). |
| **`~` (Tilde)** | Used in constraints to include types whose **underlying type** matches (e.g., `~int` matches `type MyInt int`). |
| **`cmp.Ordered`** | Stdlib constraint (Go 1.21+) covering all ordered built-in types: integers, floats, and strings. Prefer this over hand-written unions for ordering operations. |
| **Type Inference** | The compiler deduces type arguments from function call arguments. `Min(3, 7)` and `Min[int](3, 7)` are equivalent. Explicit form is needed only when T appears solely in the return type. |
| **Generic Interfaces** | Interfaces can themselves have type parameters: `type Container[T any] interface { Add(T); Get() T }`. This describes *behaviour on T*, distinct from using an interface purely as a constraint. |

---

## 2. Visual Representation

### Generic Function Workflow
Generics act as a "blueprint". The compiler generates a specific version of the code for each type you use (monomorphization).

```text
  Blueprint: func Print[T any](v T) { ... }

  +---------------------------------------+
  |              Compiler                 |
  +---------------------------------------+
      |                  |                 |
      v                  v                 v
  [Print(42)]       [Print("hi")]     [Print(1.1)]
      |                  |                 |
      v                  v                 v
  func PrintInt      func PrintStr     func PrintFloat
```

---

## 3. Implementation Examples

### Generic Functions and Constraints

```go
// Number constraint including underlying types
type Number interface {
    ~int | ~int64 | ~float64
}

// Sum uses the Number constraint
func Sum[T Number](vals []T) T {
    var total T
    for _, v := range vals {
        total += v
    }
    return total
}
```

### Generic Structs

```go
// Stack is a generic container
type Stack[T any] struct {
    elements []T
}

func (s *Stack[T]) Push(v T) {
    s.elements = append(s.elements, v)
}

func (s *Stack[T]) Pop() (T, bool) {
    if len(s.elements) == 0 {
        var zero T   // zero value pattern - see Section 4b
        return zero, false
    }
    // ... pop logic
}
```

### cmp.Ordered: The Stdlib Constraint for Ordering

Instead of writing `~int | ~int64 | ~float64 | ~string | ...`, use `cmp.Ordered` when you need comparison operators (`<`, `>`) on built-in types:

```go
import "cmp"

func Min[T cmp.Ordered](a, b T) T {
    if a < b {
        return a
    }
    return b
}

// Works with int, float64, string, and any type with an ordered underlying built-in type
Min(3, 7)          // int
Min(1.1, 2.2)      // float64
Min("apple", "banana") // string
```

> **Rule of thumb:** Use `cmp.Ordered` when you need ordering on built-in types. Write your own union constraint only when you need to restrict to a *subset* (e.g., only integers, not floats).

### Type Inference: Explicit vs Inferred

The compiler can infer type arguments from function call arguments. Explicit notation is only needed when the type parameter does not appear in the parameter list:

```go
// Inferred - compiler deduces T=int from the arguments (3, 7)
result := Min(3, 7)

// Explicit - equivalent, but rarely necessary
result := Min[int](3, 7)

// Explicit REQUIRED - T only appears in the return type, not the parameters
zero := ZeroOf[int]()   // compiler cannot infer int with no arguments
```

### Generic Interfaces with Type Parameters

An interface can itself carry a type parameter. This is different from using an interface as a constraint:

```go
// Constraint: restricts what T can be (used in [T Number])
type Number interface { ~int | ~float64 }

// Parameterized interface: describes behaviour operating on T
// A type implements Container[int] by having Add(int) and Get() int.
type Container[T any] interface {
    Add(T)
    Get() T
}

type Box[T any] struct{ value T }

func (b *Box[T]) Add(v T) { b.value = v }
func (b *Box[T]) Get() T  { return b.value }

// Compile-time check that Box[int] satisfies Container[int]
var _ Container[int] = &Box[int]{}
```

### Multi-Constraint Intersection

A type parameter can be required to satisfy multiple constraints simultaneously using an inline interface that combines them:

```go
import ("cmp"; "fmt")

// T must be BOTH ordered (supports <) AND have a String() method
func MinWithLabel[T interface{ cmp.Ordered; fmt.Stringer }](a, b T) string {
    if a < b {
        return a.String()
    }
    return b.String()
}

// Temperature satisfies the intersection:
// - underlying type float64 puts it in cmp.Ordered
// - String() method satisfies fmt.Stringer
type Temperature float64

func (t Temperature) String() string {
    return fmt.Sprintf("%.1f°C", float64(t))
}

MinWithLabel(Temperature(0), Temperature(100)) // "0.0°C"
```

---

## 4. Common Patterns & Use Cases

- **Generic Collections**: Building type-safe Stacks, Queues, Sets.
- **Utility Functions**: `Map`, `Filter`, `Reduce` that work on any slice.
  - `Filter`: select a subset using a predicate — `Filter(nums, func(n int) bool { return n > 0 })`
  - `Reduce`: collapse a slice to one value; T and U can differ — `Reduce([]int{1,2,3}, "", joinFn)` produces a `string`
- **Repository Pattern**: Defining a generic `Repository[T]` for database operations.
- **Generic Set**: A `Set[T comparable]` backed by `map[T]struct{}` for deduplication.

### Zero Value Pattern

The zero value of a type parameter `T` is obtained with `var zero T`. This is used whenever a generic function must return "nothing" of type `T`:

```go
// Canonical pattern used in Stack.Pop and First:
var zero T
return zero, false

// Standalone demonstration - explicit type arg required (no parameter to infer from):
ZeroOf[int]()     // 0
ZeroOf[string]()  // ""
ZeroOf[bool]()    // false
ZeroOf[*int]()    // nil
```

---

## 5. Critical Pitfalls & Best Practices

> [!WARNING]
> Do not use generics for everything! If your logic only works for one or two types, or if you find yourself using `any` and then type-asserting inside the function, you probably don't need generics.

1. **`any` vs `comparable`**: Use `any` for general-purpose containers (like a Stack). Use `comparable` only when you need to use the `==` operator or map keys.
2. **The `~` Tilde**: Always use `~` in your constraints (e.g., `~int` instead of `int`) if you want to support custom types that are aliases of those built-ins.
3. **Methods on Generic Structs**: While a struct can have type parameters, Go does not currently support **generic methods** (methods that introduce *new* type parameters not present on the receiver). Use a standalone function instead.
4. **`cmp.Ordered` vs custom union**: Prefer `cmp.Ordered` when you need comparison/ordering on built-in types. Write your own union only when you need a deliberate subset (e.g., integers only, excluding floats).
5. **Type inference limits**: Inference works from *function arguments*. If the type parameter only appears in the return type (like `ZeroOf[T]() T`), the compiler cannot infer it and you must be explicit.

---

## Running the Examples

```bash
go test -v ./internal/basics/generics/...
```

---

## Further Reading

- [Official Go Blog: An Introduction to Generics](https://go.dev/blog/intro-generics)
- [Go by Example: Generics](https://gobyexample.com/generics)
- [cmp package documentation](https://pkg.go.dev/cmp)
