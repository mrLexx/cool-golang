package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func worker(in <-chan Task, e *int32) {
	for w := range in {
		if err := w(); err != nil {
			atomic.AddInt32(e, 1)
		}
	}
}

func Run(tasks []Task, n, m int) error {
	var errorsCount int32
	var wg sync.WaitGroup
	jobsCh := make(chan Task)

	switch {
	case m == 0:
		m = 1
	case m == -1:
		m = len(tasks) + 1
	}

	wg.Add(n)
	for range n {
		go func() {
			defer wg.Done()
			worker(jobsCh, &errorsCount)
		}()
	}
	for _, t := range tasks {
		if int(atomic.LoadInt32(&errorsCount)) >= m {
			break
		}
		jobsCh <- t
	}
	close(jobsCh)
	wg.Wait()

	if int(errorsCount) >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
