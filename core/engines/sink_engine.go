package engines

import (
	"conecto/core"
	"conecto/core/extractors"
	"conecto/core/sinks"
	"conecto/core/statestores"
	"conecto/core/transformers"
	"math/rand/v2"
	"time"
)

type SinkEngine struct {
	Sink 		sinks.Sink
	BatchSize 	int
	MaxRetries 	int
	Backoff    	time.Duration
	MaxBackoff time.Duration
	Rand 	   *rand.Rand
	WatermarkExtractor extractors.WatermarkExtractor
	StateStore statestores.StateStore
}

func (e *SinkEngine) Run(
	runtime Runtime,
	in <-chan core.Batch,
	transformer transformers.Transformer,
) error {

	state := core.State{}

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
				batch.Events,
				transformer,
			); err != nil {
				return err
			}

			// UPDATE CHECKPOINT
			state.Cursor = batch.Cursor

			// PERSIST CHECKPOINT IMMEDIATELY
			if err := e.StateStore.Save(
				runtime.Context,
				runtime.PipelineId,
				state,
			); err != nil {
				return err
			}
		}
	}
}

func (e *SinkEngine) flush(
	runtime Runtime,
	batch []core.Event,
	transformer transformers.Transformer,
) error {

	for i := 0; i <= e.MaxRetries; i++ {

		err := e.process(runtime, batch, transformer)
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

func (e *SinkEngine) process(
	runtime Runtime,
	batch []core.Event,
	transformer transformers.Transformer,
) error {

	// transform
	transformed, err := transformer.Transform(runtime.Context, batch)
	if err != nil {
		return err
	}

	// write to sink 
	 return e.Sink.WriteBatch(runtime.Context, transformed);	

}