package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRunConcurrency(t *testing.T) {
	t.Run("concurrency", func(t *testing.T) {
		tasksCount := 50
		workerCount := 5
		maxErrorsCount := 1
		tasks := make([]Task, 0, tasksCount)

		var runningTasksCount int32

		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				atomic.AddInt32(&runningTasksCount, 1)
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runningTasksCount, -1)
				return nil
			})
		}

		conditionConcurrencyFunc := func() bool {
			return atomic.LoadInt32(&runningTasksCount) == int32(workerCount)
		}

		go Run(tasks, workerCount, maxErrorsCount)
		require.Eventually(t, conditionConcurrencyFunc, time.Second, time.Millisecond)
	})
}

func TestRunZeroErrors(t *testing.T) {
	t.Run("no_errors", func(t *testing.T) {
		tasksCount := 10
		workerCount := 2
		maxErrorsCount := 0
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		err := Run(tasks, workerCount, maxErrorsCount)
		require.NoError(t, err)
	})

	t.Run("all_errors", func(t *testing.T) {
		tasksCount := 10
		workerCount := 2
		maxErrorsCount := 0
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		err := Run(tasks, workerCount, maxErrorsCount)
		require.NoError(t, err)
	})
}

func TestRunZeroWorkers(t *testing.T) {
	t.Run("zero_workers", func(t *testing.T) {
		tasksCount := 10
		workerCount := 0
		maxErrorsCount := 3
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		err := Run(tasks, workerCount, maxErrorsCount)
		require.ErrorIs(t, err, ErrNoWorkers)
	})

	t.Run("negative_number_workers", func(t *testing.T) {
		tasksCount := 10
		workerCount := -1
		maxErrorsCount := 1
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		err := Run(tasks, workerCount, maxErrorsCount)
		require.ErrorIs(t, err, ErrNoWorkers)
	})
}

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("fewer tasks then workers", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 15
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("comparison of single-threaded and multi-threaded mode", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, 0, tasksCount)

		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				n := 1_409_305_684_859
				for i := 2; i < n; i++ {
					if n%i == 0 {
						return nil
					}
				}
				return nil
			})
		}
		workersCount := 10
		maxErrorsCount := 1

		oneWorker := make(chan struct{})
		manyWorkers := make(chan struct{})
		go func() {
			Run(tasks, 1, maxErrorsCount)
			oneWorker <- struct{}{}
		}()
		go func() {
			Run(tasks, workersCount, maxErrorsCount)
			manyWorkers <- struct{}{}
		}()

		select {
		case <-oneWorker:
			<-manyWorkers
			require.FailNow(t, "Single-threaded mode is faster than multi-threaded mode.")
		case <-manyWorkers:
			<-oneWorker
		}
	})
}
