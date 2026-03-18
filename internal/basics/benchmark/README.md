# Benchmark

Go's built-in benchmarking lets you measure function performance using `go test -bench`.

## Running

```bash
go test -bench=. -benchmem ./internal/basics/benchmark/...
```

## Key Concepts

- `BenchmarkX(b *testing.B)` — benchmark function signature
- `b.N` — iterations Go determines automatically
- `-benchmem` — reports memory allocations per operation
- Use `b.ResetTimer()` after expensive setup
