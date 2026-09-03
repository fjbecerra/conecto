package engines

import (
	"conecto/core"
	"conecto/core/statestores"
	"context"
	"errors"

	"golang.org/x/sync/errgroup"
)

type EngineRunner interface {
	Run(context context.Context, connection core.Connection) error
	Shutdown(context context.Context) error
}


type Stream struct {
	Name string
	Engine 		*Engine
	StateStore  statestores.StateStore
}

func (p *Stream) Run(context context.Context,connection core.Connection) error {
	// CONTEXT (shared cancellation)
	g, ctxWithCancel := errgroup.WithContext(context)

	// LOAD CHECKPOINT (ONLY ONCE)
	state, err := p.StateStore.Load(ctxWithCancel, connection.ID, p.Name)
	if err != nil {
		return err
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

            batchTransformed, err :=
                p.Engine.Transformer.Transform(
                    ctxWithCancel,
                    &batch,
                )

            if err != nil {
                return err
            }

            batch = *batchTransformed

            if err := p.Engine.SinkCommiter.Commit(
                ctxWithCancel,
				connection.ID,
				p.Name,
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
			p.StateStore.Save(ctxWithCancel, connection.ID, p.Name, statestores.State{
				Status: statestores.Stopped,
			})
			return err

		default:
			p.StateStore.Save(ctxWithCancel, connection.ID, p.Name, statestores.State{
				Status: statestores.Failed,
			})
			return err
		}

}
