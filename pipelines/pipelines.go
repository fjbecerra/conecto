package pipelines

import (
	"conecto/core"
	"conecto/core/sinks"
	"conecto/core/sources"
	"conecto/factories"
	"context"
	"fmt"
)

type Pipeline[T any] struct {
    source sources.Source[T]
    sink   sinks.Sink[T]
}

func (p *Pipeline[T]) Run(ctx context.Context) error {

    out, errCh1 := p.source.Fetch(ctx)
    errCh2 := p.sink.Write(ctx, out)

    for errCh1 != nil || errCh2 != nil {
        select {
        case err, ok := <-errCh1:
            if !ok {
                fmt.Println("PIPELINE: source errCh closed")
                errCh1 = nil
                continue
            }
            if err != nil {
                fmt.Println("PIPELINE: source error:", err)
                return err
            }

        case err, ok := <-errCh2:
            if !ok {
                fmt.Println("PIPELINE: sink errCh closed")
                errCh2 = nil
                continue
            }
            if err != nil {
                fmt.Println("PIPELINE: sink error:", err)
                return err
            }
        }
    }

    return nil
}

func BuildPipeline(config core.ConfigPipeline) PipelineRunner{	
	source := factories.NewSource(config.SourceConfig).Build()
	current := factories.NewTransform(source, config.TransformsConfig, config.AdditionalConfig).Build()
	sink :=factories.NewSink(config.SinkConfig, config.AdditionalConfig).Build()

	return &Pipeline[core.Record]{
		source: current,
		sink:   sink,
	}
}
type PipelineRunner interface {
	Run(ctx context.Context) error	
}

type PipelineFactory func() PipelineRunner
