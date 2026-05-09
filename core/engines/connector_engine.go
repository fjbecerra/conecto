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
	state core.State,
	out chan<- core.Batch,
) error {

	current := state.Cursor

	if err := e.Connector.Open(ctx, current); err != nil {
		return err
	}
	defer e.Connector.Close()

	for {

		var batch core.Batch
		var err error

		//RETRY FETCH
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
		// EMIT BATCH (checkpoint = CURRENT)
		select {
		case out <- core.Batch{
			Events: batch.Events,
			Cursor: current,
		}:
		case <-ctx.Done():
			return ctx.Err()
		}

		
		// ADVANCE FETCH CURSOR
		if batch.Cursor == nil {
			return nil
		}

		current = batch.Cursor
	}
}