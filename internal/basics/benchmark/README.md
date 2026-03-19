# 🚀 Benchmarking in Go

Benchmarking is the process of measuring the performance of your code. Go provides a powerful, built-in benchmarking tool as part of the `testing` package to help you establish baselines and prove optimizations.

---

## 1. Core Concepts

| Concept | Description / Purpose |
| :--- | :--- |
| **`b.N`** | The number of iterations. Go automatically increases this until results are statistically significant. |
| **`b.ResetTimer()`** | Resets the timer to exclude expensive setup logic from the benchmark results. |
| **`-benchmem`** | A flag to include memory allocation statistics (B/op and allocs/op). |
| **Compiler Optimization** | The risk where the compiler removes "unused" code, leading to misleadingly fast results. |

---

## 2. 🖼️ Visual Representation

Go's benchmark runner uses an iterative workflow to find the stable execution time.

```text
  +-------------------------------------------------------+
  |                   Benchmark Workflow                  |
  +-------------------------------------------------------+
  |  1. Start with b.N = 1                                |
  |  2. Run the loop b.N times                            |
  |  3. If too fast, increase b.N (e.g., 1, 2, 5, 10...)   |
  |  4. Repeat until time limit (default 1s) is reached   |
  |  5. Calculate averages: Total Time / b.N              |
  +-------------------------------------------------------+
            |
            v
  +-------------------------------------------------------+
  | Result: BenchmarkSomeFunc  1000000  1050 ns/op         |
  +-------------------------------------------------------+
```

---

## 3. 📝 Implementation Examples

### Anatomy of a Benchmark

```go
func BenchmarkMyFunction(b *testing.B) {
    // 1. Expensive Setup (optional)
    data := prepareHugeDataSet()
    
    // 2. Reset the timer so setup time isn't included
    b.ResetTimer()

    // 3. The Core Loop
    for i := 0; i < b.N; i++ {
        MyFunction(data)
    }
}
```

---

## 4. 🚀 Common Patterns & Use Cases

- **Performance Regression Testing**: Running benchmarks in CI to ensure new PRs don't slow down critical paths.
- **Implementation Comparison**: Comparing two algorithms (e.g., Iterative vs. Recursive) to see which scales better as input size grows.
- **Memory Profiling**: Identifying functions that allocate too much on the heap, leading to GC pressure.

---

## 5. ⚠️ Critical Pitfalls & Best Practices

> [!WARNING]
> If your function is too simple and its result isn't used, the Go compiler might "optimise" it away entirely. To prevent this, assign the result to a package-level variable.

1. **Always use `-benchmem`**: It's almost always relevant to see how much memory your code allocates.
2. **Avoid Side Effects**: Ensure your benchmark loop doesn't have side effects that affect subsequent iterations (e.g., appending to the same global slice).
3. **Reset the Timer**: Use `b.ResetTimer()` if you have more than a few microseconds of setup.

---

## 🧪 Running the Examples

Explore the unit tests to see a comparison between iterative and recursive factorial implementations.

```bash
# Run all benchmarks with memory statistics
go test -bench=. -benchmem ./internal/basics/benchmark/...

# Run for a specific time to get more stable results
go test -bench=. -benchtime=5s ./internal/basics/benchmark/...
```

---

## 📚 Further Reading

- [Go Documentation: Package testing](https://pkg.go.dev/testing#hdr-Benchmarks)
- [Go by example: Testing and Benchmarking](https://gobyexample.com/testing-and-benchmarking)
