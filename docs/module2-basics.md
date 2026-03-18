# Module 2: Go Language Fundamentals

**Duration:** Day 1, afternoon (13:30–16:00)
**Location:** `internal/basics/`

## What You'll Learn

Each topic in `internal/basics/` has a `README.md` explaining the concept and working code examples to run and explore.

| Topic | Package | Key Concepts |
|---|---|---|
| Pointers | `basics/pointers/` | Value vs pointer, mutation, pointer receivers |
| Type Assertions | `basics/casting/` | Interface → concrete type, type switch |
| Structs & Entities | `basics/entity/` | Struct design, JSON tags, interfaces |
| Package Layout | `basics/layout/` | `cmd/`, `internal/`, `pkg/`, naming |
| Parameters | `basics/parameters/` | Value vs pointer parameters, performance |
| Embedding | `basics/embed/` | Struct embedding, interface composition |
| Receivers | `basics/receivers/` | Method sets, pointer vs value receivers |
| init() | `basics/init/` | Init order, execution, pitfalls |
| Errors | `basics/err/` | Error types, wrapping, Is/As, sentinels |
| Interfaces | `basics/interface/` | Consumer-side vs producer-side definition |
| Concurrency | `basics/concurrency/` | Goroutines, channels, mutex, WaitGroup |
| Context | `basics/context/` | Cancel, timeout, values |
| HTTP | `basics/http/` | HTTP client, timeouts |
| Testing | `basics/testing/` | Table-driven tests, subtests |
| Testify | `basics/testify/` | assert, require, test suites |
| Benchmark | `basics/benchmark/` | go test -bench, benchmem |
| HTTP Testing | `basics/httptest/` | httptest.NewRecorder, NewServer |
| Mocking | `basics/mocking/` | go.uber.org/mock, EXPECT() |
| Build Tags | `basics/buildtags/` | //go:build constraints |
| Generics | `basics/generics/` | Type parameters, constraints |

## Running Examples

```bash
make test-basics               # run all basics tests
go test ./internal/basics/...  # same thing
go test -v ./internal/basics/pointers/...  # single topic, verbose
```
