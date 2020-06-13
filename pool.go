package gopool

import (
	"sync"
)

// Action for pool to do
type Action func(wg *sync.WaitGroup, pos, total, workerID int) interface{}

// ThreadPool of workers
type ThreadPool struct {
	wg *sync.WaitGroup // Waitgroup used for pool

	Total       int    // Total amount of tasks
	WorkerCount int    // Count of parallel running workers
	Action      Action // The action to run

	jobs    chan int
	results chan interface{}
	worker  []*worker
}

// New pool
func New(total int, workers int, action Action) *ThreadPool {
	// Prevent using more workers
	// than tasks
	if workers > total {
		workers = total
	}

	pool := &ThreadPool{
		wg:          &sync.WaitGroup{},
		Total:       total,
		Action:      action,
		WorkerCount: workers,
	}

	return pool
}

// init the ThreadPool
func (pool *ThreadPool) init() {
	pool.jobs = make(chan int, pool.Total)
	//pool.results = make(chan interface{}, pool.Total)
	pool.worker = make([]*worker, pool.WorkerCount)

	for i := range pool.worker {
		pool.worker[i] = newWorker(pool.wg, pool.jobs, pool.results, i, pool.Total, pool.Action)
	}
}

// return true if pool is initialized
func (pool ThreadPool) didInit() bool {
	return pool.jobs != nil && pool.worker != nil
}

// WithWG use a custom Waitgroup for pool
func (pool *ThreadPool) WithWG(wg *sync.WaitGroup) *ThreadPool {
	pool.wg = wg
	return pool
}

// WithResultChan use a custom channel to write results to
// Should be buffered using at least Total elements
func (pool *ThreadPool) WithResultChan(c chan interface{}) *ThreadPool {
	pool.results = c
	return pool
}

// Run the pool
func (pool *ThreadPool) Run() *ThreadPool {
	// Initialize pool
	if !pool.didInit() {
		pool.init()
	}

	// Start the workers
	for i := range pool.worker {
		go pool.worker[i].run()
	}

	// Fill jobs channel
	for i := 0; i < pool.Total; i++ {
		pool.wg.Add(1)
		pool.jobs <- i
	}

	close(pool.jobs)

	return pool
}

// Wait for pool to be
func (pool *ThreadPool) Wait() *ThreadPool {
	pool.wg.Wait()
	return pool
}
