# HTTP Client & Server

Go's `net/http` package includes production-ready client and server implementations. Both need explicit configuration — defaults are intentionally minimal.

## Why this matters

HTTP is the backbone of most Go services. Getting the client and server configurations wrong leads to resource leaks, hanging requests, and subtle bugs that only appear under load.

---

## Client

### Key rules

- **Always set a timeout** — `http.DefaultClient` has no timeout and can block forever.
- **Always close the response body** — unclosed bodies leak connections from the pool.
- **Check `StatusCode`, not just `err`** — a non-nil `err` means the request could not be made (network failure, DNS error). A 4xx/5xx response still has `err == nil`.

```go
client := &http.Client{Timeout: 10 * time.Second}

resp, err := client.Get("https://api.example.com/resource")
if err != nil {
    return err // network failure
}
defer resp.Body.Close() // ALWAYS close — even on non-2xx

if resp.StatusCode != http.StatusOK {
    return fmt.Errorf("unexpected status %d", resp.StatusCode)
}

body, err := io.ReadAll(resp.Body)
```

### Common pitfalls

- Using `http.DefaultClient` without a timeout in production code.
- Forgetting `defer resp.Body.Close()` — connection pool exhaustion under load.
- Treating a 404 or 500 as a Go error — it isn't. Check `resp.StatusCode`.

---

## Server

### Key rules

- **Always set server timeouts** — an unconfigured server accepts slow-loris attacks.
- **Use `http.ServeMux`** for routing; use middleware for cross-cutting concerns (logging, auth).
- **Graceful shutdown** — use `srv.Shutdown(ctx)` to drain in-flight requests before exiting.

```go
mux := http.NewServeMux()
mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    fmt.Fprintln(w, "ok")
})

srv := &http.Server{
    Addr:         ":8080",
    Handler:      loggingMiddleware(mux),
    ReadTimeout:  5 * time.Second,
    WriteTimeout: 10 * time.Second,
    IdleTimeout:  120 * time.Second,
}
```

### Middleware pattern

```go
// A middleware accepts an http.Handler and returns one.
// This lets you chain: A(B(C(handler)))
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r) // always call next
    })
}
```

### Graceful shutdown

```go
// Block until OS signal, then drain in-flight requests.
<-ctx.Done()
shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
srv.Shutdown(shutdownCtx)
```

### Common pitfalls

- No timeouts on the server → denial-of-service via slow clients.
- Calling `log.Fatal` or `os.Exit` instead of graceful shutdown → in-flight requests are dropped.
- Writing to `w` after `http.Error` or `w.WriteHeader` → headers already sent, second write is silently ignored.

---

## When NOT to use the default mux

`http.DefaultServeMux` is a package-level global. Libraries that call `http.Handle(...)` can accidentally register routes into it. Always create an explicit `http.NewServeMux()`.
