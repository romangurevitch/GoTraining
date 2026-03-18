# Pointer Parameter vs Value Parameter

## Running this demo

```bash
go test -bench '.Parameter.' -benchmem
```
```bash
# Sample output
goos: darwin
goarch: amd64
pkg: github.com/romangurevitch/go-training/internal/basics/parameters
cpu: Intel(R) Core(TM) i7-8850H CPU @ 2.60GHz
BenchmarkPointerParameter_1-12                  17398620                59.97 ns/op           40 B/op          2 allocs/op
BenchmarkValueParameter_1-12                    30026976                42.62 ns/op           32 B/op          1 allocs/op
BenchmarkPointerParameter_100-12                 1292474               938.7 ns/op          2072 B/op          9 allocs/op
BenchmarkValueParameter_100-12                    540777              2254 ns/op            8160 B/op          8 allocs/op
BenchmarkPointerParameter_10000-12                 12037             90051 ns/op          386330 B/op         21 allocs/op
BenchmarkValueParameter_10000-12                    3024            513409 ns/op         1761258 B/op         21 allocs/op
BenchmarkPointerParameter_1000000-12                  16          68624290 ns/op        45188403 B/op         41 allocs/op
BenchmarkValueParameter_1000000-12                     1        1021090515 ns/op        1616732128 B/op       51 allocs/op
BenchmarkPointerParameter_100000000-12                 1        6078855511 ns/op        4935060840 B/op       72 allocs/op
BenchmarkValueParameter_100000000-12                   1        16039609332 ns/op       18827329504 B/op      62 allocs/op
PASS
ok      github.com/romangurevitch/go-training/internal/basics/parameters       38.480s
```

## Pointer Parameter Use Cases

- When we need to modify the parameter passed into the function, we have to use pointer type.
- When the parameter is a huge in terms of memory usage, using pointer can save memory by not copying the parameter.
However, it doesn't guarantee better performance. Therefore, we don't choose pointer parameter because of performance blindly.

## Value Parameter Use Cases

- When using value parameter, we don't need to worry about the nil exception. We also don't need to worry about
the parameter being modified accidentally.
