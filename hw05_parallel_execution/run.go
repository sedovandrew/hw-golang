package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrNoWorkers           = errors.New("no workers")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		return ErrNoWorkers
	}
	taskCh := make(chan Task)
	var errorCount int32
	var wg sync.WaitGroup

	// Producer
	wg.Add(1)
	go func(errorCountPtr *int32) {
		defer wg.Done()
		defer close(taskCh)
		for _, task := range tasks {
			if atomic.LoadInt32(errorCountPtr) >= int32(m) {
				break
			}
			taskCh <- task
		}
	}(&errorCount)

	// Create workers
	wg.Add(n)
	for w := 0; w < n; w++ {
		// Consumer
		go func(errorCounterPtr *int32) {
			defer wg.Done()
			for task := range taskCh {
				if err := task(); err != nil {
					atomic.AddInt32(errorCounterPtr, 1)
				}
			}
		}(&errorCount)
	}

	wg.Wait()

	if m > 0 && errorCount >= int32(m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}
