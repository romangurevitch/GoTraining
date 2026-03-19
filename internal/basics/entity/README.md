# 🏗️ Entities & Structs in Go

In Go, an **Entity** is typically represented by a `struct` that holds data and a set of **methods** that define its behavior. Unlike traditional OOP languages, Go achieves encapsulation and polymorphism through a combination of structs, interfaces, and visibility rules.

---

## 1. Core Concepts

| Concept | Description / Purpose |
| :--- | :--- |
| **`struct`** | A collection of fields used to group related data. |
| **Methods** | Functions attached to a struct using a **receiver** (e.g., `func (u *User) Greet()`). |
| **Encapsulation** | Controlled via capitalization (e.g., `Name` is public, `id` is private). |
| **Factory (`New`)** | A convention for initializing structs with sensible defaults. |
| **Struct Tags** | Metadata attached to fields used by libraries (e.g., `json:"id"`). |
| **Interface** | Defines a contract of behaviors that an entity must satisfy. |

---

## 2. 🖼️ Visual Representation

### Exported vs. Unexported (Encapsulation)
Go uses capitalization to determine visibility. If a field or struct starts with an uppercase letter, it is **Exported** (public); otherwise, it is **Unexported** (private to the package).

```text
  +-----------------------------------+
  |           struct User             |
  +-----------------------------------+
  |  ID        (Exported/Public)      | <--- Visible to other packages
  |  Email     (Exported/Public)      |
  |                                   |
  |  role      (Unexported/Private)   | <--- Internal to THIS package only
  |  secretKey (Unexported/Private)   |
  +-----------------------------------+
```

### The Interface Contract
Instead of class inheritance, Go uses interfaces to describe what an entity can **do**.

```text
      [ Interface: User ]              [ Concrete Struct: user ]
      |  - GetName()      | <-------   |  - Name string        |
      |  - IsAdmin()      |            |  - role string        |
                                       |  - GetName() { ... }  |
                                       |  - IsAdmin() { ... }  |
```

---

## 3. 📝 Implementation Examples

### Defining an Entity with Encapsulation

```go
// User is the public interface
type User interface {
    GetName() string
    IsAdmin() bool
}

// user is the private concrete implementation
type user struct {
    Name string // Exported field
    role string // Unexported field
}

// New is the factory function
func New(name string) User {
    return &user{
        Name: name,
        role: "customer", // Default
    }
}

func (u *user) GetName() string { return u.Name }
func (u *user) IsAdmin() bool   { return u.role == "admin" }
```

---

## 4. 🚀 Common Patterns & Use Cases

- **Domain Modeling**: Representing core business objects like `Account`, `Order`, or `User`.
- **The "New" Pattern**: Providing a clean API for object creation while hiding complex setup logic.
- **JSON Mapping**: Using struct tags (`json:"user_id"`) to bridge the gap between Go field names and external API naming conventions.

---

## 5. ⚠️ Critical Pitfalls & Best Practices

> [!WARNING]
> While fields can be exported, it's often better to keep them unexported and provide getter/setter methods if you need to enforce invariants (e.g., validating an email format).

1. **Value vs. Pointer Receivers**: Use pointer receivers (`*User`) if you need to modify the struct or if the struct is large; use value receivers (`User`) for small, immutable types.
2. **Favor Composition**: Use struct embedding (see the `embed` module) instead of trying to simulate inheritance.
3. **Keep Structs Small**: If a struct starts growing too large, it might be violating the Single Responsibility Principle.

---

## 🧪 Running the Examples

Explore the unit tests to see how entities are created, how encapsulation is enforced, and how struct tags work.

```bash
# Run tests for entity patterns
go test -v ./internal/basics/entity/...
```

---

## 📚 Further Reading

- [A Tour of Go: Structs](https://go.dev/tour/moretypes/2)
- [Effective Go: Allocation with new](https://go.dev/doc/effective_go#allocation_new)
