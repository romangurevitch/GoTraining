# Go Detective: Basic "fixme" Challenges

Welcome to the **Go Detective** series! Your goal is to identify and fix 105 common Go pitfalls (The Great Go Gotchas).

These challenges are specifically designed for developers coming from dynamic languages like Python, JavaScript, or Ruby.

## How to Play

1.  **Navigate** to this directory in your terminal.
2.  **Run the tests**:
    ```bash
    go test -v ./...
    ```
3.  **Investigate**: Read the failing test output and the "Detective Brief" in the test file.
4.  **Repair**: Fix the code until all tests pass.

## Topics Covered

- Slices & Arrays (Headers, Memory, Capacity)
- Maps (Init, Missing Keys, Concurrency)
- Pointers & Memory (Dereferencing, Traps)
- Receivers & Methods (Value vs Pointer)
- Interfaces (The Nil trap, Type Assertions)
- Error Handling (Wrapping, Custom Types)
- Concurrency (Goroutines, Channels, Sync)
- Context & Testing

Good luck, Detective!
