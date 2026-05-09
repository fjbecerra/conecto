package engines

import (
	"conecto/core"
	"conecto/core/statestores"
	"conecto/core/transformers"
	"context"
	"golang.org/x/sync/errgroup"
)

type Runtime struct {
	PipelineId 		string
	Context 		context.Context
}

type Pipeline struct {	
	ConnectorEngine * ConnectorEngine
	SinkEngine 		* SinkEngine
	Transformer 	transformers.Transformer
	StateStore 		statestores.StateStore
}

func (p *Pipeline) Run(runtime Runtime) error {

	// CONTEXT (shared cancellation)
	g, ctx := errgroup.WithContext(runtime.Context)

	// LOAD CHECKPOINT (ONLY ONCE)
	state, err := p.StateStore.Load(ctx, runtime.PipelineId)
	if err != nil {
		return err
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

	// SINK (owns checkpointing)
	g.Go(func() error {

		return p.SinkEngine.Run(
			Runtime{
				Context:    ctx,
				PipelineId: runtime.PipelineId,
			},
			batches,
			p.Transformer,
		)
	})

	// WAIT
	return g.Wait()
}

type PipelineRunner interface {
	Run(runtime Runtime) error	
}

