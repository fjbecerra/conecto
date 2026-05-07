package engines

import (
	"conecto/core"
	"conecto/core/checkpoint"
	"conecto/core/transformers"
	"context"
)

type Settings struct {
	BufferSize 		int
}

type Pipeline struct {
	ID 				string
	ConnectorEngine * ConnectorEngine
	SinkEngine 		* SinkEngine
	Transformer 	transformers.Transformer
	Settings 		Settings
	StateStore		checkpoint.StateStore	
}

func (p *Pipeline) Run(ctx context.Context) error {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	state, err := p.StateStore.Load(ctx, p.ID)
	if err != nil {
		return err
	}

	events := make(chan core.Event, p.Settings.BufferSize)
	errCh := make(chan error, 1)
	doneCh := make(chan struct{}) 


	// SOURCE
	go func() {
		defer close(events)

		if err := p.ConnectorEngine.Run(ctx, state, events); err != nil {
			errCh <- err
			cancel()
		}
	}()

	// SINK
	go func() {
		defer close(doneCh)
		if err := p.SinkEngine.Run(ctx, events, p.Transformer); err != nil {
			errCh <- err
			cancel()
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-doneCh: 
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

type PipelineRunner interface {
	Run(ctx context.Context) error	
}

