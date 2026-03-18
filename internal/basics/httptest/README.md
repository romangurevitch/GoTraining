# HTTP Handler Testing

`net/http/httptest` lets you test HTTP handlers without running a real server.

## Test a Handler

```go
req := httptest.NewRequest(http.MethodGet, "/accounts/1", nil)
w := httptest.NewRecorder()
MyHandler(w, req)
resp := w.Result()
assert.Equal(t, http.StatusOK, resp.StatusCode)
```

## Pitfalls

- `httptest.NewRecorder()` buffers the full response — call `w.Result()` to get the response
