package hw06pipelineexecution

import "fmt"

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

// ExecutePipeline creates pipeline from stages.
func ExecutePipeline(in In, done In, stages ...Stage) Out {
	stageCh := in
	for _, stage := range stages {
		stageCh = runStage(done, stageCh, stage)
	}
	return stageCh
}

// runStage starts stage execution asynchronously.
// It returns channel to connect next stage.
func runStage(done In, in In, stage Stage) Out {
	outCh := make(Bi)
	go func() {
		defer func() {
			close(outCh)
			for range outCh {
			}
		}()
		for v := range stage(in) {
			select {
			case outCh <- v:
			case <-done:
				fmt.Println("Done!")
				return
			}
		}
	}()
	return outCh
}
