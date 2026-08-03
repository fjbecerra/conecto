package pipelines

import (
	"conecto/auth/connections"
	"conecto/core"
	"conecto/core/engines"
	"conecto/core/statestores"
	"context"
	"errors"

	"golang.org/x/sync/errgroup"
)

type EngineRunner interface {
	Run(context context.Context, connection connections.Connection) error
	Shutdown(context context.Context) error
}


type Stream struct {
	Engine 		*engines.Engine
	StateStore  statestores.StateStore
}

func (p *Stream) Run(context context.Context,connection connections.Connection) error {
	
	// CONTEXT (shared cancellation)
	g, ctxWithCancel := errgroup.WithContext(context)

	// LOAD CHECKPOINT (ONLY ONCE)
	state, err := p.StateStore.Load(ctxWithCancel, connection.ID)
	if err != nil {
		return err
	}
	if state.Status == statestores.Completed {
    	return nil // NOTHING TO DO
	}

	// CHANNEL
	batches := make(chan core.Batch)

	// CONNECTOR
	g.Go(func() error {


		defer close(batches)

		return p.Engine.ConnectorRunnable.Run(
			ctxWithCancel,
			state,
			batches,
			connection,
		)
	})

	// processor
    g.Go(func() error {

        for batch := range batches {

            events, err :=
                p.Engine.Transformer.Transform(
                    ctxWithCancel,
                    batch.Events,
                )

            if err != nil {
                return err
            }

            batch.Events = events

            if err := p.Engine.SinkCommiter.Commit(
                ctxWithCancel,
				connection.ID,
                batch,
            ); err != nil {
                return err
            }
        }

        return nil
    })

	// WAIT
	err = g.Wait()
	switch {

		case err == nil:
			return nil

		case errors.Is(err, ctxWithCancel.Err()):
			p.StateStore.Save(ctxWithCancel, connection.ID, statestores.State{
				Status: statestores.Stopped,
			})
			return err

		default:
			p.StateStore.Save(ctxWithCancel, connection.ID, statestores.State{
				Status: statestores.Failed,
			})
			return err
		}

}
