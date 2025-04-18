package hw06pipelineexecution

type (
	In  = <-chan any
	Out = In
	Bi  = chan any
)

type Stage func(in In) (out Out)

func drain(stream In) {
	//nolint:revive
	for range stream {
	}
}

func pipe(done In) Stage {
	return func(in In) Out {
		out := make(Bi)
		stream := in

		go func() {
			defer func() {
				close(out)
			}()

			for {
				select {
				case <-done:
					go drain(stream)
					return
				case v, ok := <-stream:
					if !ok {
						return
					}
					select {
					case out <- v:
					case <-done:
						go drain(stream)
						return
					}
				}
			}
		}()
		return out
	}
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	in = pipe(done)(in)
	for _, stage := range stages {
		in = pipe(done)(stage(in))
	}
	return in
}
