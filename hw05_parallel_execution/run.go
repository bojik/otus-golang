package hw05parallelexecution

import (
	"errors"
	"io/ioutil"
	"log"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

var l = log.Default()

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, maxWorkers, maxErrors int) error {
	l.SetOutput(ioutil.Discard) // comment for debug
	tch := make(chan Task)
	ech := make(chan error)
	stop := make(chan bool)
	wg := sync.WaitGroup{}
	defer func() {
		l.Println("waiting goroutines")
		wg.Wait()
		l.Println("closing channels")
		close(tch)
		close(ech)
		close(stop)
	}()
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go worker(i, &wg, tch, ech, stop)
	}

	ret := process(tasks, maxWorkers, maxErrors, tch, ech)

	for i := 0; i < maxWorkers; i++ {
		l.Println("Sending stop signal")
		stop <- true
	}

	return ret
}

func process(tasks []Task, mw, me int, tch chan Task, ech chan error) error {
	var (
		i, done, errs, prc int
		ret                error
	)
	count := len(tasks)

	for {
		if i < count && ret == nil {
			select {
			case tch <- tasks[i]:
				l.Printf("Writing task: %d\n", i)
				i++
				prc++
			default:
			}
		}

		select {
		case err := <-ech:
			prc--
			done++
			l.Printf("done: %d\n", done)
			if err != nil {
				errs++
				if done < me && mw+me < count {
					count = mw + me
				}
				if errs == me {
					ret = ErrErrorsLimitExceeded
					l.Println("stopping workers by errors")
				}
			}
		default:
		}
		if prc == 0 && (ret != nil || done == count) {
			break
		}
	}
	l.Println("main loop was ended")
	return ret
}

func worker(num int, wg *sync.WaitGroup, tasks <-chan Task, errs chan<- error, stop <-chan bool) {
	l.Printf("#%d: starting\n", num)
	for {
		select {
		case t := <-tasks:
			l.Printf("#%d: starting task\n", num)
			err := t()
			l.Printf("#%d: task done\n", num)
			l.Printf("#%d: writing in error channel\n", num)
			errs <- err
			l.Printf("#%d: done writing in error channel\n", num)
		case <-stop:
			l.Printf("#%d: receiving stop signal\n", num)
			wg.Done()
			l.Printf("#%d: stop worker\n", num)
			return
		}
	}
}
