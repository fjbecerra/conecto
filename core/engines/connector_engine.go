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
	ctx context.Context,
	state core.State,
	out chan<- core.Batch,
) error {

	current := state.Cursor

	if err := e.Connector.Open(ctx, current); err != nil {
		return err
	}
	defer e.Connector.Close()
	var batch core.Batch

	for {
		 err := e.Retry.Do(ctx, func() error {
    		var err error
    		batch, err = e.Connector.FetchBatch(ctx, current)
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
		case <-ctx.Done():
			return ctx.Err()
		}
	
		// END OF STREAM 
		// if batch.Cursor == nil {
		// 	return nil
		// }

		current = batch.Cursor
	}
}

func (c *ConnectorEngine) Shutdown(ctx context.Context) error {
	return c.Connector.Close()
}