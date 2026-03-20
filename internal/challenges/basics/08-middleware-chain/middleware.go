package middleware

import (
	"net/http"
)

// Middleware wraps an http.Handler and returns a new http.Handler.
type Middleware func(http.Handler) http.Handler

// Chain applies middlewares in order to an initial handler.
// Chain(h, A, B, C) returns A(B(C(h))).
func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	// TODO: implement the chaining logic.
	// Hint: iterate from the end of the slice backwards.
	panic("not implemented")
}

// LoggingMiddleware logs request details.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: log "method path"
		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware enforces a basic token check.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: check X-Auth-Token header
		// If valid ("valid-token"), call next.ServeHTTP
		// If invalid, return 401 Unauthorized
		next.ServeHTTP(w, r)
	})
}
