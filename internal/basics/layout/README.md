# 🏗️ Package Layout in Go

Idiomatic Go projects follow a consistent layout that separates concerns. Proper organization prevents "god packages" and ensures clear boundaries.

---

## 1. Core Concepts

| Concept | Description / Purpose |
| :--- | :--- |
| **`cmd/`** | Entry points for your application. One folder per binary. |
| **`internal/`** | Private application code that cannot be imported by other projects. |
| **`pkg/`** | Public library code that can be safely imported by anyone. |
| **Package Responsibility** | Naming packages by what they *do*, not what they *contain*. |

---

## 2. 🗺️ Visual Representation

```text
  +-----------------------+                     +-----------------------+
  |      cmd/myapp/       |      Calls          |      internal/        |
  |  (Main Entry Point)   |  -------------->    |   (Private Logic)     |
  +-----------------------+                     +-----------------------+
              |                                             |
              v                                             v
       (Executable Binary)       ------>             (Library Logic)
```

---

## 3. 💻 Implementation Examples

```go
// 1. Better Layout: naming by responsibility
import "myapp/pkg/files"   // usage: files.Open()
import "myapp/pkg/strings" // usage: strings.ToUpper()

// 2. Package definition
package files

import "os"

func Open(path string) (*os.File, error) {
    return os.Open(path)
}
```

---

## 4. 📋 Common Patterns & Use Cases

- **Decoupling with Internal**: Using the `internal/` directory to prevent other developers from depending on your internal API implementation details.
- **Micro-Packages**: Dividing logic into many small packages (like `net/http`, `encoding/json`) instead of one giant `common` package.
- **Clear Entry Points**: Keeping `main.go` files small, serving only to glue components together.

---

## 5. ⚠️ Critical Pitfalls & Best Practices

> [!WARNING]
> Avoid "god packages" like `util`, `common`, or `helpers`. They have no clear responsibility and inevitably become a dumping ground for unrelated code.

1. **Explicit Boundaries**: Use the `internal/` directory to enforce encapsulation at the Go toolchain level.
2. **Naming Matters**: Prefer `files.Open` over `fileutil.OpenFile`. Use names that make the call site readable.
3. **Avoid Circular Imports**: Proper package structure is essential to prevent cyclic dependencies, which are disallowed in Go.

---

## 🏃 Running the Examples

Explore the unit tests for runnable patterns:
- `layout_test.go`: Demonstrates how code in different packages interacts.

```bash
# Run tests with verbose output
go test -v ./internal/basics/layout/...
```

---

## 📚 Further Reading

- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Effective Go: Package Names](https://go.dev/doc/effective_go#package-names)
