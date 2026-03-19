# 🧪 HTTP Testing in Go

`net/http/httptest` lets you test HTTP handlers and servers without running a real server or performing actual network calls.

---

## 1. Core Concepts

| Concept | Description / Purpose |
| :--- | :--- |
| **`httptest.NewRecorder()`** | Provides an `http.ResponseWriter` implementation to capture the response. |
| **`httptest.NewRequest()`** | Creates a mock `http.Request` for testing handlers. |
| **`httptest.NewServer()`** | Spawns a real (but ephemeral) HTTP server for testing clients. |
| **`w.Result()`** | Returns the final `http.Response` from a recorder. |

---

## 2. 🗺️ Visual Representation

```text
  +-----------------------+                     +-----------------------+
  |      Mock Request     |      Process        |     Response Recorder |
  | (httptest.NewRequest) |  -------------->    | (httptest.NewRecorder)|
  +-----------------------+                     +-----------------------+
              |                                             |
              v                                             v
       Invoke MyHandler(w, req)           ------>      Assert w.Result()
```

---

## 3. 💻 Implementation Examples

```go
func TestHandler(t *testing.T) {
    // 1. Initialisation
    req := httptest.NewRequest(http.MethodGet, "/accounts/1", nil)
    w := httptest.NewRecorder()
    
    // 2. Execution
    MyHandler(w, req)
    resp := w.Result()
    defer resp.Body.Close()

    // 3. Finalisation (Verification)
    assert.Equal(t, http.StatusOK, resp.StatusCode)
}
```

---

## 4. 📋 Common Patterns & Use Cases

- **Unit Testing Handlers**: Passing a `ResponseRecorder` directly to a handler function.
- **Integration Testing Clients**: Using `httptest.NewServer` to mock a remote API.
- **JSON Body Verification**: Mocking JSON requests and verifying JSON response payloads.

---

## 5. ⚠️ Critical Pitfalls & Best Practices

> [!WARNING]
> `httptest.NewRecorder()` buffers the full response in memory. Do not use it for testing large streaming responses.

1. **Check StatusCode**: Always verify the status code before attempting to read or decode the response body.
2. **Close the Server**: Always `defer ts.Close()` when using `httptest.NewServer` to clean up resources.
3. **Use Subtests**: Run multiple test cases with `t.Run()` to isolate failures in table-driven tests.

---

## 🏃 Running the Examples

Explore the unit tests for runnable patterns:
- `httptest_test.go`: Coverage of handler recording and server mocking.

```bash
# Run tests with verbose output
go test -v ./internal/basics/httptest/...
```

---

## 📚 Further Reading

- [Official Go Documentation: net/http/httptest](https://pkg.go.dev/net/http/httptest)
- [Testing HTTP Handlers](https://go.dev/doc/tutorial/add-a-test)
