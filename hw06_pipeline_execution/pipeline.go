package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

// stageManager handle the stop signal.
func stageManager(in In, done In, stage Stage) Out {
	preStageIn := make(Bi)
	go func() {
		defer close(preStageIn)
		for n := range in {
			select {
			case <-done:
				return
			default:
				preStageIn <- n
			}
		}
	}()
	out := stage(preStageIn)
	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(Out)

	// Connect stages to the pipeline.
	for _, stage := range stages {
		out = stageManager(in, done, stage)
		in = out
	}

	return out
}
