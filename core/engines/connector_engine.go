package engines

import (
	"conecto/auth/connections"
	"conecto/core"
	"conecto/core/retry"
	"conecto/core/statestores"
	"context"
)

type ConnectorRunnable interface {
	Run(context context.Context, state statestores.State, out chan<- core.Batch, connection connections.Connection) error
	Shutdown(ctx context.Context) error
}

type ConnectorEngine struct {
	Connector core.Connector
	Retry retry.Executor
}

func (e *ConnectorEngine) Run(
	context context.Context,
	state statestores.State,
	out chan<- core.Batch,
	connection connections.Connection,
) error {

	current := state.Cursor

	defer e.Connector.Close()
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
		var isLast bool
		if batch.Cursor == nil {
			isLast = true
		}
		select {
		case out <- core.Batch{
			Events: batch.Events,
			Cursor: batch.Cursor,
			IsLast: isLast,
		}:
		case <-context.Done():
			return context.Err()
		}

		current = batch.Cursor
	}
}

func (c *ConnectorEngine) Shutdown(ctx context.Context) error {
	return c.Connector.Close()
}