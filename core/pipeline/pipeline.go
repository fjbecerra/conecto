package pipeline

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
	Run(context context.Context) error
	Shutdown(context context.Context) error
}

type Pipeline struct {
	Connection 	connections.Connection
	Engine 		*engines.Engine
	StateStore  statestores.StateStore
}

func (p *Pipeline) Run(context context.Context) error {
	
	defer p.Shutdown(context)

	// CONTEXT (shared cancellation)
	g, ctxWithCancel := errgroup.WithContext(context)

	// LOAD CHECKPOINT (ONLY ONCE)
	state, err := p.StateStore.Load(ctxWithCancel, p.Connection.ID)
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
			p.Connection,
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
				p.Connection.ID,
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
			p.StateStore.Save(ctxWithCancel, p.Connection.ID, statestores.State{
				Status: statestores.Stopped,
			})
			return err

		default:
			p.StateStore.Save(ctxWithCancel,p.Connection.ID, statestores.State{
				Status: statestores.Failed,
			})
			return err
		}

}

func (p *Pipeline) Shutdown(ctx context.Context) error {

    var errs []error

    if err := p.Engine.ConnectorRunnable.Shutdown(ctx); err != nil {
        errs = append(errs, err)
    }

    if err := p.Engine.SinkCommiter.Shutdown(ctx); err != nil {
        errs = append(errs, err)
    }

    return errors.Join(errs...)
}
