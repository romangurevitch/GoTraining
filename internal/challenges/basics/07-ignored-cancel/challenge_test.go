package ignoredcancel

import (
	"context"
	"testing"
	"time"
)

func TestWatch_ReportsAccountID(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // always defer cancel — prevents goroutine leaks

	reported := make(chan string, 1)
	stopped := make(chan struct{})

	Watch(ctx, "acc-123", func(id string) {
		select {
		case reported <- id:
		default:
		}
	}, stopped)

	select {
	case got := <-reported:
		if got != "acc-123" {
			t.Fatalf("expected report for acc-123, got %q", got)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Watch never called report — did you start a goroutine?")
	}
}

func TestWatch_StopsOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	reported := make(chan struct{}, 1)
	stopped := make(chan struct{})

	Watch(ctx, "acc-999", func(string) {
		select {
		case reported <- struct{}{}:
		default:
		}
	}, stopped)

	// Wait for at least one report (confirms the goroutine started)
	select {
	case <-reported:
	case <-time.After(2 * time.Second):
		t.Fatal("Watch never called report — goroutine did not start")
	}

	// Cancel the context
	cancel()

	// Confirm the goroutine stopped (closed the stopped channel)
	select {
	case <-stopped:
		// success
	case <-time.After(2 * time.Second):
		t.Fatal("Watch goroutine did not stop after context cancellation — did you select on ctx.Done()?")
	}
}
