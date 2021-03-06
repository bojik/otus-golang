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

var (
	task Task = func() error {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
		return nil
	}

	errTask Task = func() error {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
		return errors.New("Hzhz")
	}
)

func TestSimple(t *testing.T) {
	t.Run("simple test", func(t *testing.T) {
		tasks := []Task{task, errTask, task}
		err := Run(tasks, 2, 2)
		require.Nil(t, err)
	})
	t.Run("workers > tasks", func(t *testing.T) {
		tasks := []Task{task, errTask, task}
		err := Run(tasks, 10, 2)
		require.Nil(t, err)
	})
}

func TestError(t *testing.T) {
	t.Run("error of limit exceeded is expected", func(t *testing.T) {
		tasks := []Task{task, errTask, task, errTask, task}
		err := Run(tasks, 2, 2)
		require.ErrorIs(t, err, ErrErrorsLimitExceeded)
	})
	t.Run("error of maxErrors less zero is expected", func(t *testing.T) {
		tasks := []Task{task, errTask, task, errTask, task}
		err := Run(tasks, 2, -2)
		require.ErrorIs(t, err, ErrMaxErrorsLessZero)
	})
	t.Run("error of maxWorkers less one is expected", func(t *testing.T) {
		tasks := []Task{task, errTask, task, errTask, task}
		err := Run(tasks, 0, 100)
		require.ErrorIs(t, err, ErrMaxWorkersLessOne)
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
				<-time.After(time.Millisecond * time.Duration(rand.Intn(100)))
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

	t.Run("if maxErrors equals zero, no errors expected", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				<-time.After(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 0
		err := Run(tasks, workersCount, maxErrorsCount)
		require.Nil(t, err)
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
				<-time.After(taskSleep)
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
}
