# 🏗️ Entities & Structs in Go

In Go, an **Entity** is typically represented by a `struct` that holds data and a set of **methods** that define its behavior. It's the building block of your domain model.

---

## 1. What is a Struct?

A `struct` (short for "structure") is a collection of fields. It's Go's way of grouping related data together.

### 🖼️ Exported vs. Unexported Fields
Go uses capitalization to determine visibility (encapsulation).

```text
  +-----------------------------------+
  |           struct User             |
  +-----------------------------------+
  |  ID        (Exported/Public)      | <--- Visible to other packages
  |  Email     (Exported/Public)      |
  |                                   |
  |  password  (Unexported/Private)   | <--- Internal to THIS package only
  |  secretKey (Unexported/Private)   |
  +-----------------------------------+
```

---

## 2. Behavior (Methods)

An entity is more than just data; it has behavior. We attach methods to structs using **receivers**.

```go
type User struct {
    Name string
}

// Greet is a method attached to the User struct
func (u User) Greet() string {
    return "Hello, I am " + u.Name
}
```

---

## 3. Interfaces: The Contract

In Go, we often define an `interface` to describe what an entity can **do**, rather than what it **is**. This allows for "duck typing" and easier mocking in tests.

### 🖼️ The Interface Contract
```text
      [ Interface: Account ]           [ Concrete Struct: SavingsAccount ]
      |  - GetBalance()    | <-------  |  - balance float64              |
      |  - Withdraw(amt)   |           |  - GetBalance() { ... }         |
                                       |  - Withdraw(amt) { ... }        |
```

---

## 4. Struct Tags (Metadata)

Struct tags allow you to attach metadata to fields, which is used by libraries for JSON encoding, database mapping, or validation.

```go
type User struct {
    ID    int    `json:"id" db:"user_id"`
    Email string `json:"email" validate:"required"`
}
```

---

## 5. Construction (The `New` Pattern)

Since Go doesn't have constructors, we use "factory functions" (conventionally named `New...`) to initialize entities with sensible defaults.

```go
func NewUser(name string) *User {
    return &User{
        Name: name,
    }
}
```

---

## 🧪 Running the Examples

Explore `entity_test.go` to see how entities are created and used.

```bash
go test -v ./internal/basics/entity/...
```
