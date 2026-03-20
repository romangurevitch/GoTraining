package pool

import (
	"context"
	"sync"
)

// Job represents a unit of work.
type Job struct {
	ID      int
	Payload []byte
}

// Result holds the outcome of processing a Job.
type Result struct {
	JobID int
	Err   error
}

// Pool distributes jobs to N goroutines.
type Pool struct {
	workers int
	jobs    chan Job
	results chan Result
	wg      sync.WaitGroup
}

// NewPool creates a new worker pool.
func NewPool(workers int, bufferSize int) *Pool {
	return &Pool{
		workers: workers,
		jobs:    make(chan Job, bufferSize),
		results: make(chan Result, bufferSize),
	}
}

// Start launches worker goroutines that process jobs until ctx is cancelled.
func (p *Pool) Start(ctx context.Context) {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-p.jobs:
					if !ok {
						return
					}
					// TODO: Process the job and send to results
					p.results <- Result{JobID: job.ID, Err: nil}
				}
			}
		}()
	}
}

// Submit sends a job to the pool.
func (p *Pool) Submit(job Job) {
	p.jobs <- job
}

// Results returns the channel to receive job outcomes.
func (p *Pool) Results() <-chan Result {
	return p.results
}

// Stop closes the jobs channel and waits for all workers to finish.
func (p *Pool) Stop() {
	close(p.jobs)
	p.wg.Wait()
	close(p.results)
}
