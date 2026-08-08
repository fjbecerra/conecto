package engines

import (
	"conecto/core"
	"conecto/core/retry"
	"conecto/core/statestores"
	"context"
)

type ConnectorRunnable interface {
	Run(context context.Context, state statestores.State, out chan<- core.Batch, connection core.Connection) error
}

type ConnectorEngine struct {
	Connector core.Connector
	Retry retry.Executor
}

func (e *ConnectorEngine) Run(
	context context.Context,
	state statestores.State,
	out chan<- core.Batch,
	connection core.Connection,
) error {

	current := state.Cursor

	var batch core.Batch
	for {
		 err := e.Retry.Do(context, func() error {
    		var err error
    		batch, err = e.Connector.FetchBatch(context, current, connection)
    		return err
		})
		if err != nil {
			return err
		}
		
		if len(batch.Events) == 0 {
			return nil
		}		
		
		//end of stream		
		isLast := batch.Cursor == nil

		select {
		case out <- core.Batch{
			Events: batch.Events,
			Cursor: batch.Cursor,
			IsLast: isLast,
		}:
		case <-context.Done():
			return context.Err()
		}

		if isLast {
			return nil
		}

		current = batch.Cursor
	}
}
