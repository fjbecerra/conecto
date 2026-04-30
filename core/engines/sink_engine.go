package engines

import (
	"conecto/core"
	"conecto/core/sinks"
	"conecto/core/transformers"
	"context"
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
}

func (e *SinkEngine) Run(ctx context.Context, in <-chan core.Event, transformer transformers.Transformer) error {

	batch := make([]core.Event, 0, e.BatchSize)

	for ev := range in {

		batch = append(batch, ev)
		if len(batch) >= e.BatchSize {
			if err := e.processWithRetry(ctx, transformer, batch); err != nil {
				return err
			}
			batch = batch[:0]
		}
	}

	// flush remaining
	if len(batch) > 0 {
		if err := e.processWithRetry(ctx, transformer, batch); err != nil {
			return err
		}
	}

	return nil

}

func (e *SinkEngine) process(
	ctx context.Context,
	transformer transformers.Transformer,
	batch []core.Event,
) error {

	out, err := transformer.Transform(ctx, batch)
	if err != nil {
		return err
	}

	return e.Sink.WriteBatch(ctx, out)
}

func (e *SinkEngine) processWithRetry(ctx context.Context, transformer transformers.Transformer,batch []core.Event) error {

	for i := 0; i <= e.MaxRetries; i++ {

		err := e.process(ctx, transformer, batch)

		if err == nil {
			return nil
		}

		if i == e.MaxRetries {
			return err
		}

		delay := backoffWithJitter(e.Backoff, i, e.MaxBackoff, e.Rand)

		// backoff with cancel support
		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}