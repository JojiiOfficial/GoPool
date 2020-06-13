package gopool

import (
	"sync"
)

// worker represents a single worker
// created by pool
type worker struct {
	wg       *sync.WaitGroup
	jobs     <-chan int
	results  chan interface{}
	action   Action
	workerID int
	total    int
}

// Create a new worker
func newWorker(wg *sync.WaitGroup, jobs <-chan int, results chan interface{}, workerID, total int, action Action) *worker {
	return &worker{
		wg:       wg,
		jobs:     jobs,
		results:  results,
		action:   action,
		workerID: workerID,
		total:    total,
	}
}

// Run the worker
func (w *worker) run() {
	// Run all jobs until the
	// channel is closed
	for i := range w.jobs {
		// Execute action and write its result to the result channel
		res := w.action(w.wg, i, w.total, w.workerID)

		// Write result to channel if set
		if w.results != nil {
			w.results <- res
		}

		// Set task as done
		w.wg.Done()
	}
}
