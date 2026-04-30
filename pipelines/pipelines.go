package pipelines

import (
	"conecto/core"
	"conecto/core/engines"
	"conecto/core/transformers"
	"conecto/factories"
	"context"
	"math/rand/v2"
	"time"
)

type settings struct {
	bufferSize 		int
}

type Pipeline struct {
	connectorEngine * engines.ConnectorEngine
	sinkEngine 		* engines.SinkEngine
	transformer 	transformers.Transformer
	settings 		settings
}

func (p *Pipeline) Run(ctx context.Context) error {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	events := make(chan core.Event, p.settings.bufferSize)
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
	seed := time.Now().UnixNano()

	r := rand.New(rand.NewPCG(
		uint64(seed),
		uint64(seed>>1),
	))
	
	connector := factories.NewConnector(config.ConnectorConfig, r).Build()
	transform := factories.NewTransform(config.TransformersConfig, config.AdditionalConfig).Build()
	sink :=factories.NewSink(config.SinkConfig, config.AdditionalConfig, r).Build()
	settings := settings {
		bufferSize: config.AdditionalConfig.BufferSize,
	}

	return &Pipeline{
		connectorEngine: &connector,
		sinkEngine:   &sink,
        transformer: transform,
        settings: settings,
	}
}
type PipelineRunner interface {
	Run(ctx context.Context) error	
}

type PipelineFactory func() PipelineRunner
