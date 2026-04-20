package contextandgoroutines

import "context"

// requestIDKey is an unexported type to prevent collisions with other context keys.
type requestIDKey string

const key = requestIDKey("requestID")

// WithRequestID returns a new context with the given requestID.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, key, requestID)
}

// GetRequestID retrieves the requestID from the context.
// It returns an empty string if the ID is not found.
func GetRequestID(ctx context.Context) string {
	val, ok := ctx.Value(key).(string)
	if !ok {
		return ""
	}
	return val
}

// RunAsync executes the given task in a new goroutine.
func RunAsync(task func()) {
	go task()
}
