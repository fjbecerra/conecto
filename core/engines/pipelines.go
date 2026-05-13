package engines

import (
	"conecto/core"
	"conecto/core/statestores"
	"conecto/core/transformers"
	"context"
	"errors"

	"golang.org/x/sync/errgroup"
)

type Pipeline struct {	
	ConnectorEngine * ConnectorEngine
	Transformer 	transformers.Transformer
	StateStore 		statestores.StateStore
	CommitStrategy   CommitStrategy
}

func (p *Pipeline) Run(runtime core.Runtime) error {
	defer p.Shutdown(context.Background())

	// CONTEXT (shared cancellation)
	g, ctx := errgroup.WithContext(runtime.Context)

	// LOAD CHECKPOINT (ONLY ONCE)
	state, err := p.StateStore.Load(runtime)
	if err != nil {
		return err
	}
	if state.Status == core.Completed {
    	return nil // NOTHING TO DO
	}

	// CHANNEL
	batches := make(chan core.Batch)

	// CONNECTOR
	g.Go(func() error {

		defer close(batches)

		return p.ConnectorEngine.Run(
			ctx,
			state,
			batches,
		)
	})

	// processor
    g.Go(func() error {

        for batch := range batches {

            events, err :=
                p.Transformer.Transform(
                    ctx,
                    batch.Events,
                )

            if err != nil {
                return err
            }

            batch.Events = events

            if err := p.CommitStrategy.Commit(
                runtime,
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

		case errors.Is(err, context.Canceled):
			p.StateStore.Save(runtime, core.State{
				Status: core.Stopped,
			})
			return err

		default:
			p.StateStore.Save(runtime, core.State{
				Status: core.Failed,
			})
			return err
		}

}

func (p *Pipeline) Shutdown(ctx context.Context) error {

    var errs []error

    if err := p.ConnectorEngine.Shutdown(ctx); err != nil {
        errs = append(errs, err)
    }

    if err := p.CommitStrategy.Shutdown(ctx); err != nil {
        errs = append(errs, err)
    }

    return errors.Join(errs...)
}

type PipelineRunner interface {
	Run(runtime core.Runtime) error	
}

