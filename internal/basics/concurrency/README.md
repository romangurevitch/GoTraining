# Concurrency

Go's concurrency is built on goroutines and channels.

## Goroutines

```go
go func() { /* runs concurrently */ }()
```

## Channels

```go
ch := make(chan int)      // unbuffered
ch := make(chan int, 10)  // buffered
ch <- 42                  // send
v := <-ch                 // receive
close(ch)                 // signal done
```

## Sync Primitives

- `sync.WaitGroup` — wait for goroutines to finish
- `sync.Mutex` / `sync.RWMutex` — protect shared state

## Pitfalls

- Always use `-race` flag when testing: `go test -race ./...`
- Never close a channel from the receiver side
