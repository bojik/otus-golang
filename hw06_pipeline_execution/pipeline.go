package hw06pipelineexecution

import "sync"

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

// ExecutePipeline creates pipeline from stages.
func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		in = runStage(done, in, stage)
	}
	return in
}

// runStage starts stage execution asynchronously.
// It returns channel to connect next stage.
func runStage(done In, in In, stage Stage) Out {
	stageCh := make(Bi)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for v := range in {
			select {
			case stageCh <- v:
			case <-done:
				return
			}
		}
	}()
	go func() {
		wg.Wait()
		close(stageCh)
	}()
	return stage(stageCh)
}
