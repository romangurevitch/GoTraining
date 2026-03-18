package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type contextKey string

const requestIDKey contextKey = "request_id"

const headerRequestID = "X-Request-Id"

// RequestIDMiddleware generates a UUID per request, injects it into context
// and the response header. slog-gin reads it automatically when WithRequestID: true.
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(headerRequestID)
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Writer.Header().Set(headerRequestID, requestID)
		c.Set(string(requestIDKey), requestID)
		c.Next()
	}
}

// RequestIDFromCtx extracts the request ID from the Gin context.
func RequestIDFromCtx(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}
