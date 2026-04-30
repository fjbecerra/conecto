package engines

import (
	"conecto/core"
	"conecto/core/connectors"
	"context"
	"math/rand/v2"
	"time"
)

type ConnectorEngine struct {
	Connector connectors.Connector
	MaxRetries int
	Backoff    time.Duration
	MaxBackoff time.Duration
	Rand 	   *rand.Rand
}

func (e *ConnectorEngine) Run(
	ctx context.Context,
	state core.Cursor,
	out chan<- core.Event,
) error {

	if err := e.Connector.Open(ctx, state); err != nil {
		return err
	}
	defer e.Connector.Close()

	current := state

	for {

		var batch core.Batch
		var err error

		for i := 0; i <= e.MaxRetries; i++ {

			batch, err = e.Connector.FetchBatch(ctx, current)
			if err == nil {
				break
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

		if len(batch.Events) == 0 {
			return nil
		}

		for _, ev := range batch.Events {
			ev.Cursor = batch.Cursor

			select {
			case out <- ev:
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		if batch.Cursor == nil{
			return nil
		}

		current = batch.Cursor
	}
}