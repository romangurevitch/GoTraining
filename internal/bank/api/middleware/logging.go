package middleware

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

// LoggingMiddleware returns a Gin middleware that logs one structured line per request.
// Fields logged: time, method, path, status, latency, trace_id, span_id, request_id.
// trace_id/span_id are extracted from the active OTel span in the request context.
//
// Source: replaces hand-rolled JSONLogMiddleware from go-training-cba-solution/internal/server/rest/middleware/logger.go
// Upgrade: uses slog-gin instead of logrus — zero-boilerplate OTel trace correlation.
func LoggingMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return sloggin.NewWithConfig(logger, sloggin.Config{
		DefaultLevel:     slog.LevelInfo,
		ClientErrorLevel: slog.LevelWarn,
		ServerErrorLevel: slog.LevelError,
		WithTraceID:      true,
		WithSpanID:       true,
		WithRequestID:    true,
	})
}
