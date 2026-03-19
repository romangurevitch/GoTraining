package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHowToCreateAnHTTPClient_Config(t *testing.T) {
	client := HowToCreateAnHTTPClient()

	assert.NotNil(t, client, "client must not be nil")
	assert.NotNil(t, client.Transport, "transport must be explicitly configured")
	assert.Equal(t, 10*time.Second, client.Timeout, "timeout must be 10s — never use zero (blocks forever)")
}

func TestHowToCreateAnHTTPClient_MakesRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("hello"))
		assert.NoError(t, err)
	}))
	defer ts.Close()

	client := HowToCreateAnHTTPClient()
	resp, err := client.Get(ts.URL)
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "hello", string(body))
}

// TestHTTPClient_NonSuccessIsNotAnError shows a key Go HTTP pitfall:
// a 4xx or 5xx response does NOT return a non-nil error from client.Do.
// You must check resp.StatusCode yourself.
func TestHTTPClient_NonSuccessIsNotAnError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer ts.Close()

	client := HowToCreateAnHTTPClient()
	resp, err := client.Get(ts.URL)
	require.NoError(t, err, "a 404 is NOT a Go error — you must check StatusCode")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
