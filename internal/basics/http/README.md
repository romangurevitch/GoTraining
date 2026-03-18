# HTTP Client

Go's `net/http` package includes a production-ready HTTP client. Always configure timeouts.

## Basic Client

```go
client := &http.Client{Timeout: 10 * time.Second}
resp, err := client.Get("https://api.example.com/resource")
if err != nil { return err }
defer resp.Body.Close()
```

## Pitfalls

- Always call `resp.Body.Close()` — unclosed bodies leak connections
- `http.DefaultClient` has no timeout — it can block indefinitely
- Check `resp.StatusCode` — a non-nil err only means the request failed, not a 4xx/5xx
