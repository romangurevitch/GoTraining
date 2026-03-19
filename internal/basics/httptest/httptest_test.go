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

// https://golang.org/src/net/http/httptest/example_test.go
func TestMockServerResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("hello"))
		assert.NoError(t, err)
	}))
	defer ts.Close()

	res, err := http.Get(ts.URL) // run actual http get request
	assert.NoError(t, err)
	defer func() {
		err = res.Body.Close()
		assert.NoError(t, err)
	}()

	greeting, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, "hello", string(greeting))
}

func TestHandlerUsingRequestRecoder(t *testing.T) {
	target := "https://example.com/foo"
	input := "input string"
	req := httptest.NewRequest(http.MethodGet, target, strings.NewReader(input))
	w := httptest.NewRecorder()

	handler := func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)

		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, target, r.RequestURI)
		assert.Equal(t, input, string(body))

		// Write response
		w.Header().Add("key", "val")
		_, err = w.Write([]byte("hello"))
		assert.NoError(t, err)
	}
	handler(w, req)

	// Verify response
	res := w.Result()
	defer func() {
		err := res.Body.Close()
		assert.NoError(t, err)
	}()

	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "val", res.Header.Get("key"))
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
	defer resp.Body.Close()

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
			defer resp.Body.Close()

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
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
