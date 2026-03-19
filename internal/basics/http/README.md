# 🌐 HTTP Client & Server in Go

Go's `net/http` package includes production-ready client and server implementations. Both need explicit configuration — defaults are intentionally minimal.

---

## 1. Core Concepts

| Concept | Description / Purpose |
| :--- | :--- |
| **`http.Client`** | Sends HTTP requests and receives responses. Must be configured with timeouts. |
| **`http.Server`** | Listens for incoming HTTP requests and routes them to handlers. |
| **`http.ServeMux`** | An HTTP request multiplexer (router) that matches URL patterns to handlers. |
| **`http.Handler`** | An interface with a `ServeHTTP` method that processes a single request. |
| **Middleware** | Functions that wrap handlers to provide cross-cutting concerns (logging, auth). |

---

## 2. 🗺️ Visual Representation

```mermaid
flowchart LR
    C["HTTP Client\n(Timeout, Transport)"]
    S["HTTP Server\n(Mux, Middleware)"]
    C --"Request"--> S
    S --"Response"--> C
```

---

## 3. 💻 Implementation Examples

```go
// 1. Client initialization with timeout
client := &http.Client{
    Timeout: 10 * time.Second,
}

// 2. Server with production-ready timeouts
srv := &http.Server{
    Addr:         ":8080",
    Handler:      loggingMiddleware(mux),
    ReadTimeout:  5 * time.Second,
    WriteTimeout: 10 * time.Second,
    IdleTimeout:  120 * time.Second,
}

// 3. Standard Middleware pattern
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}
```

---

## 4. 📋 Common Patterns & Use Cases

- **Graceful Shutdown**: Using `srv.Shutdown(ctx)` to drain in-flight requests before the process exits.
- **Client Connection Pooling**: Reusing `http.Client` and `http.Transport` instances to maintain persistent connections.
- **Route Specific Middleware**: Applying logic only to certain paths within a `ServeMux`.

---

## 5. ⚠️ Critical Pitfalls & Best Practices

> [!WARNING]
> Never use `http.DefaultClient` or `http.ListenAndServe` in production. They lack timeouts and can lead to resource exhaustion or "hanging" processes.

1. **Always Close Response Body**: Use `defer resp.Body.Close()` immediately after checking the error to prevent connection leaks.
2. **Check StatusCode**: A nil error from `client.Do()` only means the request was sent and a response received; you must still check if it was a 4xx or 5xx.
3. **Avoid Default Mux**: `http.DefaultServeMux` is a global. Libraries can register routes on it without your knowledge. Always use `http.NewServeMux()`.

---

## 🏃 Running the Examples

Explore the unit tests for runnable patterns:
- `client_test.go`: Demonstrates client configuration and error handling.

```bash
# Run tests with verbose output
go test -v ./internal/basics/http/...
```

---

## 📚 Further Reading

- [Official Go Documentation: net/http](https://pkg.go.dev/net/http)
- [Effective Go: HTTP Server](https://go.dev/doc/effective_go#http_server)
