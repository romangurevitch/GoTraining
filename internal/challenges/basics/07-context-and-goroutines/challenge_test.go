package contextandgoroutines

import (
	"context"
	"sync"
	"testing"
)

func TestContextValues(t *testing.T) {
	ctx := context.Background()

	// Initial context should not have the request ID
	id := GetRequestID(ctx)
	if id != "" {
		t.Errorf("expected empty string, got %q", id)
	}

	// Add request ID to context
	ctx = WithRequestID(ctx, "req-12345")

	// Retrieve request ID
	id = GetRequestID(ctx)
	if id != "req-12345" {
		t.Errorf("expected req-12345, got %q", id)
	}
}

func TestRunAsync(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	executed := false
	task := func() {
		executed = true
		wg.Done()
	}

	RunAsync(task)

	// Wait for the goroutine to finish
	wg.Wait()

	if !executed {
		t.Error("expected task to be executed asynchronously")
	}
}
