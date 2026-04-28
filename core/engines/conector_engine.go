package engines

import (
	"conecto/core"
	"conecto/core/connectors"
	"context"
	"time"
)

type ConnectorEngine struct {
	Connector connectors.Connector
	MaxRetries int
	Backoff    time.Duration
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
			time.Sleep(e.Backoff * time.Duration(1<<i))
		}

		if err != nil {
			return err
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