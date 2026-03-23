package ignoredcancel

import "context"

// Watch monitors accountID and calls report for each check cycle.
// It must stop when ctx is cancelled — no goroutine leaks allowed.
//
// stopped is closed by the goroutine when it exits. The test uses this
// to confirm the goroutine actually stopped (race-free, no time.Sleep).
//
// TODO: implement the goroutine inside Watch so that:
//  1. It calls report(accountID) on each iteration.
//  2. It selects on ctx.Done() and returns when the context is cancelled.
//  3. It closes stopped when it exits (use defer — it's already there).
func Watch(ctx context.Context, accountID string, report func(string), stopped chan struct{}) {
	go func() {
		defer close(stopped) // signals the test: this goroutine has exited
		for {
			select {
			case <-ctx.Done():
				// TODO: return here to stop the goroutine
				panic("implement me: return when context is cancelled")
			default:
				// TODO: call report(accountID) here
				panic("implement me: call report on each iteration")
			}
		}
	}()
}
