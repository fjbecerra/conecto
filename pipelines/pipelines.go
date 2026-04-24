package pipelines

import (
	"conecto/core"
	"conecto/core/sinks"
	"conecto/core/sources"
    "conecto/factories"
	"context"
)

type Pipeline[T any] struct {
    source sources.Source[T]
    sink   sinks.Sink[T]
}

type PipelineRunner interface {
	Run(ctx context.Context) error	
}

func (p *Pipeline[T]) Run(ctx context.Context) error {

    out, errCh1 := p.source.Fetch(ctx)
    errCh2 := p.sink.Write(ctx, out)

    for errCh1 != nil || errCh2 != nil {
        select {
        case err, ok := <-errCh1:
            if !ok {
                errCh1 = nil
                continue
            }
            if err != nil {
                return err
            }

        case err, ok := <-errCh2:
            if !ok {
                errCh2 = nil
                continue
            }
            if err != nil {
                return err
            }
        }
    }

    return nil
}

func BuildPipeline(configPath string) PipelineRunner{
	config := core.LoadConfigPipeline(configPath)
	source := factories.NewSource(config.SourceConfig).Build()
	current := factories.NewTransform(source, config.TransformsConfig, config.AdditionalConfigs).Build()
	sink :=factories.NewSink(config.SinkConfig, config.AdditionalConfigs).Build()

	return &Pipeline[core.Record]{
		source: current,
		sink:   sink,
	}
}

type PipelineFactory func() PipelineRunner
