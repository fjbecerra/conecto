package engines

import (
	"conecto/core"
	"conecto/core/checkpoint"
	"conecto/core/sinks"
	"conecto/core/transformers"
	"context"
	"math/rand/v2"
	"time"
)

type SinkEngine struct {
	Committer 	checkpoint.Committer
	Sink 		sinks.Sink
	BatchSize 	int
	MaxRetries 	int
	Backoff    	time.Duration
	MaxBackoff time.Duration
	Rand 	   *rand.Rand	
}

func (e *SinkEngine) Run(ctx context.Context, in <-chan core.Event, transformer transformers.Transformer) error {

	batch := make([]core.Event, 0, e.BatchSize)
	state := core.State{}

	for ev := range in {

		batch = append(batch, ev)

		if ev.Timestamp.After(state.Watermark) {
			state.Watermark = ev.Timestamp
		}
		state.Cursor = ev.Cursor

		if len(batch) >= e.BatchSize {
			if err := e.commitWithRetry(ctx, batch, state, transformer); err != nil {
				return err
			}
			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		return e.commitWithRetry(ctx, batch, state, transformer)
	}

	return nil

}

func (e *SinkEngine) commitWithRetry(
	ctx context.Context,
	batch []core.Event,
	state core.State,
	transformer transformers.Transformer,
) error {

	for i := 0; i <= e.MaxRetries; i++ {

		err := e.process(ctx, batch, state, transformer)
		if err == nil {
			return nil
		}

		if i == e.MaxRetries {
			return err
		}

		delay := backoffWithJitter(e.Backoff, i, e.MaxBackoff, e.Rand)

		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

func (e *SinkEngine) process(
	ctx context.Context,
	batch []core.Event,
	state core.State,
	transformer transformers.Transformer,
) error {

	out, err := transformer.Transform(ctx, batch)
	if err != nil {
		return err
	}

	return e.Committer.Commit(ctx, out, state)
}