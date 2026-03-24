package contextandgoroutines

import "context"

// requestIDKey is an unexported type to prevent collisions with other context keys.
type requestIDKey string

const key = requestIDKey("requestID")

// WithRequestID returns a new context with the given requestID.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	// TODO: Use context.WithValue to store the requestID using the unexported key.
	panic("implement me: WithRequestID")
}

// GetRequestID retrieves the requestID from the context.
// It returns an empty string if the ID is not found.
func GetRequestID(ctx context.Context) string {
	// TODO: Retrieve the value from the context using ctx.Value.
	// Don't forget to use type assertion: val, ok := ctx.Value(key).(string)
	panic("implement me: GetRequestID")
}

// RunAsync executes the given task in a new goroutine.
func RunAsync(task func()) {
	// TODO: Run the task asynchronously using the 'go' keyword.
	panic("implement me: RunAsync")
}
