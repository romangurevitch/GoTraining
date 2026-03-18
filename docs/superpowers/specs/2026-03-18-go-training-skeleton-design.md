# Design Spec: GoTraining Repository Skeleton

**Date:** 2026-03-18
**Status:** Approved
**Module:** `github.com/romangurevitch/go-training`

---

## 1. Overview

Create the skeleton framework for a generic Go training repository called **GoTraining**. The repo teaches Go to engineers transitioning from Python/Bash, using a banking domain ("Go Bank") as the applied use case. It is designed to be delivered to any audience — not CBA-specific.

This spec covers the initial skeleton: directory structure, placeholder packages, migrated demo content, supporting infrastructure, and tooling. It does not cover implementing the full challenge content (bank API/store/service) or Temporal — those are separate efforts.

---

## 2. Source Repositories

Two repos inform the structure and patterns:

- **go-training-cba** (`/Users/romang/workspace/CBA/go-training-cba`): 6-session course with build-tag-based progressive unlocking, Makefile-driven testing, and quest markdown files. Provides the domain model, Makefile patterns, docker-compose, and the `demo/` folder content to migrate.
- **ConcurrencyWorkshop** (`/Users/romang/workspace/CBA/ConcurrencyWorkshop`): Full-day concurrency workshop with `fixme/` (debug bugs) and `implme/` (implement stubs) challenge patterns, plus rich README documentation per topic.

---

## 3. Key Design Decisions

| Decision | Choice | Rationale |
|---|---|---|
| Learning structure | Module-based (4 modules) | Matches design doc; each module is a logical unit |
| Progressive unlocking | Packages, not build tags | Modules are self-contained and run independently |
| Basics pattern | ConcurrencyWorkshop style | README explaining concepts + practical examples |
| Basics content source | Migrate from go-training-cba/demo/ | Avoids duplication; same pattern, adapted module path |
| Challenges location | `internal/challenges/` | Cleanly separated from demos and reference code |
| Bank stubs | Empty packages + README only | Content to be developed separately |
| Temporal | Empty package + README placeholder | Someone else developing content |
| Fundamentals | README only | Someone else developing content |
| Infrastructure | docker-compose + Postgres + migration stubs | Full setup needed for bank challenge |
| Module name | `github.com/romangurevitch/go-training` | Author's GitHub username |

---

## 4. Repository Structure

```text
GoTraining/
├── cmd/
│   ├── hello/              # First Go app (build/test/docker lab)
│   │   └── main.go
│   ├── bank-api/           # Go Bank REST server entry point
│   │   └── main.go
│   ├── bank-cli/           # Go Bank CLI tool entry point
│   │   └── main.go
│   └── worker/             # Temporal worker (placeholder — compilable stub only, excluded from make build)
│       └── main.go
├── internal/
│   ├── fundamentals/       # Module 1: API design principles (README only)
│   │   └── README.md
│   ├── basics/             # Module 2: Go language building blocks
│   │   ├── pointers/
│   │   ├── casting/
│   │   ├── entity/
│   │   ├── layout/
│   │   ├── parameters/
│   │   ├── embed/
│   │   ├── receivers/
│   │   ├── init/
│   │   ├── err/
│   │   ├── interface/
│   │   ├── concurrency/
│   │   ├── context/
│   │   ├── http/
│   │   ├── testing/
│   │   ├── testify/
│   │   ├── benchmark/
│   │   ├── httptest/
│   │   ├── mocking/
│   │   ├── buildtags/
│   │   └── generics/
│   ├── bank/               # Module 3: Go Bank service (challenge — empty packages)
│   │   ├── domain/         # Account & Transaction models
│   │   ├── api/            # Gin routers, middleware, handlers
│   │   ├── store/          # go-jet Postgres repository
│   │   └── service/        # Business logic
│   ├── temporal/           # Module 4: Temporal (placeholder)
│   │   └── README.md
│   └── challenges/         # Student exercises
│       ├── README.md
│       ├── basics/         # Go basics challenges
│       │   ├── fixme/      # Find and fix the bugs (ConcurrencyWorkshop pattern)
│       │   └── implme/     # Implement the function
│       └── bank/           # Bank service challenge quests
│           └── README.md
├── pkg/                    # Shared public contracts
│   └── errors/             # Error types (doc.go with package declaration only; no implementation yet)
├── migration/              # SQL migrations (Postgres)
│   ├── 001_create_accounts.sql       # Accounts table stub
│   └── 002_create_transactions.sql   # Transactions table stub
├── docs/                   # Workshop guides and module documentation
│   ├── module1-fundamentals.md
│   ├── module2-basics.md
│   ├── module3-bank.md
│   ├── module4-temporal.md
│   └── setup.md
├── docker-compose.yaml     # Postgres local stack
├── Makefile                # Build, test by module, lint, db targets
├── go.mod                  # Module: github.com/romangurevitch/go-training
├── go.sum
├── .golangci.yaml          # Linter config (adapted from go-training-cba)
├── .gitignore
└── README.md               # Top-level: setup, navigation, module overview
```

---

## 5. Module Details

### Module 1: Fundamentals (`internal/fundamentals/`)

Contains a single `README.md` only. Content is being developed separately. The README describes what this module covers (API design, security, observability, lifecycle, agentic future) and links to external resources.

### Module 2: Basics (`internal/basics/`)

Each subdirectory follows the **ConcurrencyWorkshop pattern**:
- `README.md` — explains the concept, shows pitfalls, links resources
- `.go` implementation file(s) — working examples demonstrating the concept
- `_test.go` file(s) — table-driven tests with `testify/assert`

Content is **migrated from go-training-cba/demo/** with the module path updated from `gitlab.mantelgroup.com.au/training/golang-demo` to `github.com/romangurevitch/go-training`.

Topics with demo source material (from go-training-cba/demo/):
- benchmark, buildtags, casting, concurrency, context, embed, entity, err, generics, http, httptest, init, interface, layout, mocking, parameters, receivers, testing, testify

Topics without demo source material (create README + minimal working example following the same pattern):
- `pointers/` — referenced in design doc but no demo/ counterpart; create a `README.md` explaining pointer semantics and a `pointers.go` + `pointers_test.go` with a minimal working example (e.g., a function showing value vs pointer mutation)

### Module 3: Bank (`internal/bank/`)

Empty packages with a `doc.go` (package declaration + doc comment) and `README.md` in each subdirectory. No implementation yet — this is the student challenge.

Subpackages:
- `domain/` — Account and Transaction domain models
- `api/` — HTTP handlers, routers, middleware (Gin)
- `store/` — Database repository (go-jet + Postgres)
- `service/` — Business logic and service layer

### Module 4: Temporal (`internal/temporal/`)

Single `README.md` only. Describes what Temporal is, why it's used for long-running workflows, and references the TransferFunds saga pattern. Content (workflow + activity code) to be added separately.

### Challenges (`internal/challenges/`)

Follows the **ConcurrencyWorkshop challenge pattern**:

- `basics/fixme/` — Contains buggy code for students to diagnose and fix
- `basics/implme/` — Contains `panic("implement me!")` stubs for students to implement
- `bank/` — Quest descriptions (markdown) for the bank challenge. Students implement the bank API, store, and service layer using the domain models.

For the skeleton: create the directory structure with `README.md` placeholders. Actual challenge content (buggy code, stubs, quest descriptions) to be populated separately.

---

## 6. Infrastructure

### docker-compose.yaml

```yaml
services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: gobank
      POSTGRES_USER: gotrainer
      POSTGRES_PASSWORD: verysecret
    ports:
      - "5432:5432"
```

### migration/

Placeholder SQL migration files for Go Bank schema:
- `001_create_accounts.sql` — accounts table stub
- `002_create_transactions.sql` — transactions table stub

### Makefile Targets

| Target | Description |
|---|---|
| `make all` | Tidy, lint, build, test |
| `make build` | Build all cmd/ binaries |
| `make test` | Run all tests |
| `make test-basics` | Run basics module tests only |
| `make test-bank` | Run bank module tests only |
| `make test-challenges` | Run challenge tests only |
| `make lint` | Run golangci-lint |
| `make fmt` | Format all Go source |
| `make bench` | Run all benchmarks |
| `make db-up` | Start Postgres via docker-compose |
| `make db-down` | Stop Postgres |
| `make migrate` | Run SQL migrations (stub: prints instructions; tool TBD by bank challenge implementer) |
| `make tidy` | Run go mod tidy |

### pkg/errors/

Contains only a `doc.go` with the package declaration and a doc comment. No types or functions in the skeleton — content added when the bank challenge is implemented.

### cmd/worker/

Contains a `main.go` with a compilable stub (`func main() { /* TODO: Temporal worker — see internal/temporal/README.md */ }`). This binary is **excluded from the `make build` target** to keep the skeleton clean; build it explicitly with `go build ./cmd/worker/...` once Temporal content is added.

---

## 7. Dependencies (go.mod)

Core dependencies to include:

| Package | Version | Purpose |
|---|---|---|
| `github.com/gin-gonic/gin` | v1.12.0 | HTTP router (bank API) |
| `github.com/go-jet/jet/v2` | v2.14.1 | Type-safe SQL builder (bank store) |
| `github.com/lib/pq` | v1.10.9 | Postgres driver |
| `github.com/spf13/cobra` | v1.10.2 | CLI framework (bank-cli) |
| `github.com/spf13/viper` | v1.21.0 | Configuration management |
| `github.com/stretchr/testify` | v1.11.1 | Test assertions and suites |
| `go.uber.org/mock` | v0.6.0 | Interface mocking (successor to deprecated golang/mock) |
| `log/slog` | stdlib | Structured logging (Go 1.21+) |
| `golang.org/x/sync` | v0.10.0 | errgroup and sync utilities |

---

## 8. Go Version

Go 1.26.1 (latest stable; supports generics, slog, errgroup improvements).

---

## 9. Naming Conventions

- No references to "CBA" anywhere. Domain is **"Go Bank"**.
- Module path: `github.com/romangurevitch/go-training`
- Package names: short, lowercase, no underscores (Go standard)
- Binary names: `hello`, `bank-api`, `bank-cli`, `worker`

---

## 10. Out of Scope (This Skeleton)

- Full bank API implementation (handlers, routes, middleware)
- Bank store implementation (Postgres queries, go-jet models)
- Bank service implementation (business logic)
- Temporal workflow and activity code
- Module 1 fundamentals content
- Actual challenge content (buggy code for fixme, stubs for implme, quest descriptions)
- CI/CD pipeline (.gitlab-ci.yml or .github/workflows)
- gRPC (not in design doc for this training)
