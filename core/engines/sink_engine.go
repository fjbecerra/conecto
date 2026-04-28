package engines

import (
	"conecto/core"
	"conecto/core/sinks"
	"conecto/core/transformers"
	"context"
)

type SinkEngine struct {
	Sink sinks.Sink
	BatchSize int
}

func (e *SinkEngine) Run(
	ctx context.Context,
	
	in <-chan core.Event,
	transformer transformers.Transformer,
) error {

	batch := make([]core.Event, 0, e.BatchSize)

	for ev := range in {

		batch = append(batch, ev)

		if len(batch) >= e.BatchSize {

			if err := e.process(ctx, transformer, batch); err != nil {
				return err
			}

			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		return e.process(ctx, transformer, batch)
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