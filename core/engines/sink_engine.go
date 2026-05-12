package engines

import (
	"conecto/core"
	"conecto/core/sinks/committers"
	"conecto/core/transformers"
	"math/rand/v2"
	"time"
)

type SinkEngine struct {
	
	BatchSize 	int
	MaxRetries 	int
	Backoff    	time.Duration
	MaxBackoff time.Duration
	Rand 	   *rand.Rand	
	Commiter 	committers.Committer
}

func (e *SinkEngine) Run(
	runtime Runtime,
	in <-chan core.Batch,
	transformer transformers.Transformer,
) error {
	for {

		select {
		case <-runtime.Context.Done():
			return runtime.Context.Err()

		case batch, ok := <-in:

			if !ok {
				return nil
			}

			// FLUSH
			if err := e.flush(
				runtime,
				batch,
				transformer,
			); err != nil {
				return err
			}
		}
	}
}

func (e *SinkEngine) flush(
	runtime Runtime,
	batch core.Batch,
	transformer transformers.Transformer,
) error {

	for i := 0; i <= e.MaxRetries; i++ {

		err := e.Commiter.CommitBatch(runtime.Context, runtime.PipelineId, batch ,transformer)
		if err == nil {
			return nil
		}

		if i == e.MaxRetries {
			return err
		}

		delay := backoffWithJitter(e.Backoff, i, e.MaxBackoff, e.Rand)

		select {
		case <-time.After(delay):
		case <-runtime.Context.Done():
			return runtime.Context.Err()
		}
	}

	return nil
}