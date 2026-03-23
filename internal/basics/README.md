# Module 2: Go Language Fundamentals

This module covers the core building blocks of the Go programming language. Each directory contains code examples and a focused `README.md` to guide you through the concepts.

## Topics

<a name="the-mental-shift"></a>
### The Mental Shift
- **[Pointers](pointers/README.md)** — Memory addresses and indirect values
- **[Type Assertions & Casting](casting/README.md)** — Working with dynamic types and interfaces
- **[Parameters](parameters/README.md)** — Passing values vs. pointers

<a name="structs--layout"></a>
### Structs & Layout
- **[Entities](entity/README.md)** — Defining data structures
- **[Package Layout](layout/README.md)** — Organizing your code idiomatically
- **[Embedding](embed/README.md)** — Composition over inheritance

<a name="behaviours"></a>
### Behaviours
- **[Receivers](receivers/README.md)** — Adding methods to types
- **[init()](init/README.md)** — Package initialization
- **[Error Handling](err/README.md)** — Sentinel errors, wrapping, and the `error` interface
- **[Interfaces](interface/README.md)** — Implicit implementation and decoupling

<a name="concurrency--context"></a>
### Concurrency & Context
- **[Concurrency](concurrency/README.md)** — Goroutines and Channels
- **[Context](context/README.md)** — Cancellation, deadlines, and request-scoped values

<a name="testing--benchmarking"></a>
### Testing & Benchmarking
- **[Testing](testing/README.md)** — Unit testing and table-driven tests
- **[Testify](testify/README.md)** — Fluent assertions and requirements
- **[Benchmark](benchmark/README.md)** — Performance measurement
- **[HTTP Testing](httptest/README.md)** — Testing handlers without a network

<a name="advanced-features"></a>
### Advanced Features
- **[Generics](generics/README.md)** — Writing type-agnostic code
- **[Mocking](mocking/README.md)** — Using Mockery for dependency isolation
- **[Build Tags](buildtags/README.md)** — Conditional compilation
