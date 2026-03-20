package httptest

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// myHandler is a minimal handler used by the recorder examples below.
func myHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("key", "val")
	_, _ = w.Write([]byte("hello"))
}

// --- Recorder: handler unit tests (no network) ---

func TestHandlerUsingRecorder(t *testing.T) {
	// 1. Initialisation
	req := httptest.NewRequest(http.MethodGet, "https://example.com/foo", strings.NewReader("input"))
	w := httptest.NewRecorder()

	// 2. Execution
	myHandler(w, req)
	res := w.Result()
	defer func() { require.NoError(t, res.Body.Close()) }()

	// 3. Finalisation (Verification)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "val", res.Header.Get("key"))
}

// --- Test Server: client integration tests (real HTTP) ---

// TestMockServerResponse shows the basic test server pattern.
// https://golang.org/src/net/http/httptest/example_test.go
func TestMockServerResponse(t *testing.T) {
	// 1. Initialisation
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("hello"))
		assert.NoError(t, err)
	}))
	defer ts.Close()

	// 2. Execution
	res, err := http.Get(ts.URL)
	assert.NoError(t, err)
	defer func() { require.NoError(t, res.Body.Close()) }()

	// 3. Finalisation (Verification)
	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, "hello", string(body))
}

// TestPostWithJSONBody tests a POST handler that decodes a JSON request body
// and echoes it back. Shows the encode/decode pattern used in most REST APIs.
func TestPostWithJSONBody(t *testing.T) {
	type payload struct {
		Message string `json:"message"`
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		var p payload
		err := json.NewDecoder(r.Body).Decode(&p)
		require.NoError(t, err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		require.NoError(t, json.NewEncoder(w).Encode(p))
	}))
	defer ts.Close()

	body, _ := json.Marshal(payload{Message: "hello"})
	resp, err := http.Post(ts.URL, "application/json", strings.NewReader(string(body)))
	require.NoError(t, err)
	defer func() { require.NoError(t, resp.Body.Close()) }()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var got payload
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&got))
	assert.Equal(t, "hello", got.Message)
}

// TestNonSuccessStatusCode verifies that non-2xx responses are not Go errors.
// You must inspect resp.StatusCode — a 400/500 does not set err != nil.
func TestNonSuccessStatusCode(t *testing.T) {
	for _, code := range []int{http.StatusBadRequest, http.StatusInternalServerError} {
		code := code
		t.Run(http.StatusText(code), func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "error", code)
			}))
			defer ts.Close()

			resp, err := http.Get(ts.URL)
			require.NoError(t, err, "non-2xx is NOT a Go error")
			defer func() { require.NoError(t, resp.Body.Close()) }()

			assert.Equal(t, code, resp.StatusCode)
		})
	}
}

// TestRequestHeadersArePassedThrough verifies custom headers set on a request
// arrive at the server unchanged.
func TestRequestHeadersArePassedThrough(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := r.Header.Get("X-Custom-Header")
		assert.Equal(t, "my-value", val)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	require.NoError(t, err)
	req.Header.Set("X-Custom-Header", "my-value")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer func() { require.NoError(t, resp.Body.Close()) }()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
