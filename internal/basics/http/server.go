package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

// loggingMiddleware wraps a handler and logs each incoming request method and path.
// This is the standard middleware pattern: accept an http.Handler, return an http.Handler.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// NewServer creates an HTTP server with named routes, a logging middleware,
// and production-ready timeouts.
//
// Key timeouts:
//   - ReadTimeout:  time to read the full request (headers + body)
//   - WriteTimeout: time to write the full response
//   - IdleTimeout:  keep-alive connection idle time
func NewServer(addr string) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

	mux.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		msg := r.URL.Query().Get("msg")
		if msg == "" {
			http.Error(w, "missing msg parameter", http.StatusBadRequest)
			return
		}
		fmt.Fprintln(w, msg)
	})

	return &http.Server{
		Addr:         addr,
		Handler:      loggingMiddleware(mux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}

// StartWithGracefulShutdown starts the server in a goroutine and blocks until
// ctx is cancelled, then shuts down gracefully with a 5-second drain window.
//
// Graceful shutdown lets in-flight requests finish before the process exits.
func StartWithGracefulShutdown(ctx context.Context, srv *http.Server) error {
	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	}
}
