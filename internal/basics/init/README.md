# 🚀 The `init()` Function in Go

`init()` runs automatically when a package is loaded, before `main()`. It is used for package-level initialization and setup.

---

## 1. Core Concepts

| Concept | Description / Purpose |
| :--- | :--- |
| **Package Initialisation** | The process where package-level variables are assigned their initial values. |
| **`init()`** | A special function that runs after variables are initialized but before `main()`. |
| **Execution Order** | Variables -> `init()` (in file order) -> `main()`. |

---

## 2. 🗺️ Visual Representation

```text
  +-----------------------+                     +-----------------------+
  |    Global Variables   |      Assigned       |    init() Functions   |
  |  (Memory Allocated)   |  -------------->    |  (Setup Logic Ran)    |
  +-----------------------+                     +-----------------------+
              |                                             |
              v                                             v
       Next Step: main()         ------>             (Ready for Use)
```

---

## 3. 💻 Implementation Examples

```go
package initializer

var localVar string

func init() {
    // 1. Initialisation
    localVar = "is that what you expect?"
}

func GetVar() string {
    // 2. Execution (Ready for use by callers)
    return localVar
}
```

---

## 4. 📋 Common Patterns & Use Cases

- **Package Level Setup**: Initializing configurations or global states (though use sparingly).
- **Self-registration**: Registering drivers (e.g., `sql` drivers or `image` formats) into a central registry.
- **Environment Variable Checks**: Validating that required environment variables are set during startup.

---

## 5. ⚠️ Critical Pitfalls & Best Practices

> [!WARNING]
> `init()` functions cannot be called explicitly and run without context. Avoid complex logic or side effects that make unit testing harder.

1. **Avoid DB Connections**: Do not open database connections in `init()` as it hides errors and prevents proper resource management.
2. **Deterministic Order**: Within a package, `init()` functions run in the order their source files are compiled (typically alphabetical). Relying on this is fragile.
3. **Prefer Explicit Setup**: For critical components, use a `New()` or `Setup()` function instead of relying on `init()`.

---

## 🏃 Running the Examples

Explore the unit tests for runnable patterns:
- `init_test.go`: Shows the order of initialization across different files.

```bash
# Run tests with verbose output
go test -v ./internal/basics/init/...
```

---

## 📚 Further Reading

- [Effective Go: The init function](https://go.dev/doc/effective_go#init)
- [Official Go Documentation: Program initialization](https://go.dev/ref/spec#Package_initialization)
