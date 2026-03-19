package concurrency

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 1. Goroutines & WaitGroup
// WaitGroup is the canonical way to wait for multiple goroutines to complete.
func TestWaitGroup(t *testing.T) {
	var wg sync.WaitGroup
	var counter int32

	numGoroutines := 10
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			// atomic is another way to handle simple numeric synchronization
			atomic.AddInt32(&counter, 1)
		}()
	}

	wg.Wait() // Block until all wg.Done() calls are made
	assert.Equal(t, int32(numGoroutines), counter)
}

// 2. Mutex (Mutual Exclusion)
// Use Mutex to protect shared state from concurrent access (Race Conditions).
func TestMutex(t *testing.T) {
	var mu sync.Mutex
	var wg sync.WaitGroup
	counter := 0

	numGoroutines := 100
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()

			mu.Lock()
			// Critical section
			counter++
			mu.Unlock()
		}()
	}

	wg.Wait()
	assert.Equal(t, numGoroutines, counter)
}

// 3. Channels (Communication)
// Channels are used for "orchestration" and passing ownership of data.
func TestUnbufferedChannel(t *testing.T) {
	ch := make(chan string)

	go func() {
		// This will block until the receiver is ready
		ch <- "ping"
	}()

	msg := <-ch // This blocks until a value is sent
	assert.Equal(t, "ping", msg)
}

func TestBufferedChannel(t *testing.T) {
	// Buffered channels don't block until the buffer is full
	ch := make(chan int, 2)

	ch <- 1
	ch <- 2

	assert.Equal(t, 2, len(ch))

	assert.Equal(t, 1, <-ch)
	assert.Equal(t, 2, <-ch)
}

// 4. Select (Multiplexing)
// Select allows waiting on multiple channel operations.
func TestSelect(t *testing.T) {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(10 * time.Millisecond)
		ch1 <- "one"
	}()
	go func() {
		time.Sleep(20 * time.Millisecond)
		ch2 <- "two"
	}()

	timeout := time.After(1 * time.Second)
	results := []string{}
	for i := 0; i < 2; i++ {
		select {
		case res := <-ch1:
			results = append(results, res)
		case res := <-ch2:
			results = append(results, res)
		case <-timeout:
			t.Fatal("timeout")
		}
	}

	assert.Contains(t, results, "one")
	assert.Contains(t, results, "two")
}

// 5. Worker Pool Pattern
// A classic pattern to limit concurrency.
func TestWorkerPool(t *testing.T) {
	jobs := make(chan int, 100)
	results := make(chan int, 100)

	var wg sync.WaitGroup
	// Start 3 workers
	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go func(id int, jobs <-chan int, results chan<- int) {
			defer wg.Done()
			for j := range jobs {
				// Simulate work
				results <- j * 2
			}
		}(w, jobs, results)
	}

	// Send 5 jobs
	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs) // Closing 'jobs' tells workers to stop when done

	// Wait for all workers to finish, then close results
	wg.Wait()
	close(results)

	// Collect results
	count := 0
	for res := range results {
		assert.True(t, res%2 == 0)
		count++
	}
	assert.Equal(t, 5, count)
}

// 6. sync.Once
// Ensures a function is only executed once, regardless of how many goroutines call it.
func TestSyncOnce(t *testing.T) {
	var once sync.Once
	var counter atomic.Int32

	increment := func() {
		counter.Add(1)
	}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			once.Do(increment)
		}()
	}

	wg.Wait()
	assert.Equal(t, int32(1), counter.Load())
}

// 7. Context Cancellation
// Proper way to stop goroutines.
func TestContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan int)

	go func() {
		i := 0
		for {
			select {
			case <-ctx.Done():
				return // Stop when context is cancelled
			case ch <- i:
				i++
			}
		}
	}()

	// Read 3 values
	assert.Equal(t, 0, <-ch)
	assert.Equal(t, 1, <-ch)
	assert.Equal(t, 2, <-ch)

	cancel() // Signal the goroutine to stop

	// Verify the context was cancelled
	assert.Equal(t, context.Canceled, ctx.Err())
}

// 8. Closing Channels
// range over a channel continues until the channel is closed.
func TestRangeOverClosedChannel(t *testing.T) {
	ch := make(chan int, 5)
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch)

	sum := 0
	for v := range ch {
		sum += v
	}
	assert.Equal(t, 6, sum)

	// Receiving from a closed channel returns the zero value and false
	val, ok := <-ch
	assert.Equal(t, 0, val)
	assert.False(t, ok)
}

func TestPanicClosingClosedChannel(t *testing.T) {
	ch := make(chan int)
	close(ch)

	assert.Panics(t, func() {
		close(ch)
	}, "closing a closed channel should panic")
}

func TestPanicWritingToClosedChannel(t *testing.T) {
	ch := make(chan int)
	close(ch)

	assert.Panics(t, func() {
		ch <- 1
	}, "writing to a closed channel should panic")
}
