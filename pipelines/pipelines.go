package pipelines

import (
	"conecto/core"
	"conecto/core/engines"
	"conecto/core/transformers"
	"conecto/factories"
	"context"
)

type Pipeline struct {
	connectorEngine * engines.ConnectorEngine
	sinkEngine * engines.SinkEngine
	transformer transformers.Transformer
	bufferSize int
}

func (p *Pipeline) Run(ctx context.Context) error {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	events := make(chan core.Event, p.bufferSize)
	errCh := make(chan error, 1)
	doneCh := make(chan struct{}) 


	// SOURCE
	go func() {
		defer close(events)

		if err := p.connectorEngine.Run(ctx, nil, events); err != nil {
			errCh <- err
			cancel()
		}
	}()

	// SINK
	go func() {
		defer close(doneCh)
		if err := p.sinkEngine.Run(ctx, events, p.transformer); err != nil {
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

func BuildPipeline(config core.ConfigPipeline) PipelineRunner{	
	connector := factories.NewConnector(config.ConnectorConfig).Build()
	transform := factories.NewTransform(config.TransformersConfig, config.AdditionalConfig).Build()
	sink :=factories.NewSink(config.SinkConfig, config.AdditionalConfig).Build()

	return &Pipeline{
		connectorEngine: &connector,
		sinkEngine:   &sink,
        transformer: transform,
        bufferSize: 10,
	}
}
type PipelineRunner interface {
	Run(ctx context.Context) error	
}

type PipelineFactory func() PipelineRunner
