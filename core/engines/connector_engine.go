package engines

import (
	"conecto/core"
	"conecto/core/connectors"
	"conecto/core/retry"
	"context"
)



type ConnectorEngine struct {
	Connector connectors.Connector
	Retry retry.Executor
}

func (e *ConnectorEngine) Run(
	runtime core.Runtime,
	state core.State,
	out chan<- core.Batch,
) error {

	current := state.Cursor

	defer e.Connector.Close()
	var batch core.Batch

	for {
		 err := e.Retry.Do(runtime.Context, func() error {
    		var err error
    		batch, err = e.Connector.FetchBatch(runtime, current)
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
		case <-runtime.Context.Done():
			return runtime.Context.Err()
		}

		current = batch.Cursor
	}
}

func (c *ConnectorEngine) Shutdown(ctx context.Context) error {
	return c.Connector.Close()
}