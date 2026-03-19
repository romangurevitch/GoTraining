package context

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 1. Basic Cancellation
// Demonstrate how to manually trigger cancellation.
func TestContextWithCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel() // Manually trigger cancellation
	}()

	select {
	case <-ctx.Done():
		assert.Equal(t, context.Canceled, ctx.Err())
	case <-time.After(1 * time.Second):
		t.Fatal("Context was never cancelled")
	}
}

// 2. Timeout (The most common usage)
// Automatically cancels after a duration.
func TestContextWithTimeout(t *testing.T) {
	// A context that expires in 50ms
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel() // Always call cancel, even on timeout

	select {
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for context timeout!")
	case <-ctx.Done():
		// Context should be cancelled by the timer
		assert.Equal(t, context.DeadlineExceeded, ctx.Err())
	}
}

// 3. Deadline
// Automatically cancels at a specific point in time.
func TestContextWithDeadline(t *testing.T) {
	deadline := time.Now().Add(50 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	select {
	case <-ctx.Done():
		assert.Equal(t, context.DeadlineExceeded, ctx.Err())
	case <-time.After(1 * time.Second):
		t.Fatal("Deadline was never reached")
	}
}

// 4. Hierarchy / Propagation
// When a parent is cancelled, ALL children are cancelled.
func TestContextHierarchy(t *testing.T) {
	parentCtx, parentCancel := context.WithCancel(context.Background())
	childCtx, childCancel := context.WithCancel(parentCtx)
	defer childCancel() // Good practice, though parentCancel() would handle it too

	parentCancel() // Cancel the parent

	select {
	case <-childCtx.Done():
		// The child is automatically cancelled because the parent was
		assert.Equal(t, context.Canceled, childCtx.Err())
	case <-time.After(1 * time.Second):
		t.Fatal("Child context didn't inherit parent cancellation")
	}
}

// 5. Context Values (Request-Scoped Data)
// Use custom types for keys to avoid collisions.
type requestIDKey int

const ridKey requestIDKey = 0

func TestContextWithValue(t *testing.T) {
	ctx := context.WithValue(context.Background(), ridKey, "abc-123")

	// Retrieval
	val := ctx.Value(ridKey)
	assert.Equal(t, "abc-123", val)

	// Values also propagate down the tree
	childCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	childVal := childCtx.Value(ridKey)
	assert.Equal(t, "abc-123", childVal)

	// Missing values return nil
	assert.Nil(t, ctx.Value("non-existent"))
}

// 6. Practical Pattern: Work Function
// A common way to structure functions that respect context.
func doWork(ctx context.Context) (int, error) {
	// Simulate some async work
	// Buffered so the goroutine can complete its send even if the caller has returned due to ctx cancellation
	resCh := make(chan int, 1)
	go func() {
		time.Sleep(50 * time.Millisecond)
		resCh <- 42
	}()

	select {
	case <-ctx.Done():
		return 0, ctx.Err() // Return the reason for cancellation
	case res := <-resCh:
		return res, nil
	}
}

func TestDoWorkRespectsContext(t *testing.T) {
	t.Run("SuccessCase", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		val, err := doWork(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 42, val)
	})

	t.Run("TimeoutCase", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		val, err := doWork(ctx)
		assert.Error(t, err)
		assert.Equal(t, context.DeadlineExceeded, err)
		assert.Equal(t, 0, val)
	})
}
