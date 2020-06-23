# GoPool
A simple thread pool library for golang.

# Usage

```go
	// Create pool
	pool := gopool.New(MAX, THREADS, func(wg *sync.WaitGroup, pos, total, workerID int) interface{} {
	    // Do something
		return nil
	})

	// Use custom result channel
	pool.WithResultChan(resultChan)

	// Start pool and wait for it to complete
	pool.Run().Wait()

```