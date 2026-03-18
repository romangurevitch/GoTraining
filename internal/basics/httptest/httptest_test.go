package httptest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
