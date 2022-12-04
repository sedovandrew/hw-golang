package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	taskCh := make(chan Task)
	errCh := make(chan int)

	go func() {
		// Create workers
		var wg sync.WaitGroup
		wg.Add(n)
		for w := 0; w < n; w++ {
			go func() {
				defer wg.Done()
				for task := range taskCh {
					err := task()
					if err != nil {
						errCh <- 1
					}
				}
			}()
		}
		// Waiting for all workers to finish
		wg.Wait()
		close(errCh)
	}()

	errors := 0
	var returnError error
Loop:
	for i := 0; i < len(tasks); {
		select {
		case taskCh <- tasks[i]:
			i++
		case err := <-errCh:
			errors += err
			if errors >= m {
				returnError = ErrErrorsLimitExceeded
				break Loop
			}
		}
	}
	// Inform the workers to complete
	close(taskCh)
	// Waiting for all workers to finish
	for range errCh {
	}
	return returnError
}
