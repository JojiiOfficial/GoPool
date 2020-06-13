package gopool

import (
	"runtime"
	"sort"
	"sync"
	"testing"
)

type taskData struct {
	tasksDone []int
	mx        sync.RWMutex
}

// Test if all jobs will be executed
func TestRunningAll(t *testing.T) {
	var td taskData
	total := 10000000
	workers := runtime.NumCPU()

	pool := New(total, workers, func(wg *sync.WaitGroup, pos, total, workerID int) interface{} {
		//	time.Sleep(time.Duration(rand.Int63n(20)*10+1) * time.Millisecond)

		td.mx.Lock()
		defer td.mx.Unlock()

		td.tasksDone = append(td.tasksDone, pos)
		return nil
	})

	pool.Run().Wait()

	if len(td.tasksDone) != total {
		t.Errorf("Not all jobs were executed! Executed tasks: %d", len(td.tasksDone))
	}

	sort.Sort(sort.IntSlice(td.tasksDone))

	// Check if all tasks were executed
	for i := 0; i < total-1; i++ {
		if td.tasksDone[i] != td.tasksDone[i+1]-1 {
			t.Errorf("Missing %d", td.tasksDone[i+1])
		}
	}
}
