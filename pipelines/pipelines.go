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

    out, srcErr := p.source.Fetch(ctx)
    sinkErr,done := p.sink.Write(ctx, out)

    for srcErr != nil || sinkErr != nil || done !=nil {
        select {
            case err, ok := <-srcErr:
                if !ok {
                    fmt.Println("PIPELINE: source errCh closed")
                    srcErr = nil
                    continue
                }
                if err != nil {
                    fmt.Println("PIPELINE: source error:", err)
                    return err
                }

            case err, ok := <-sinkErr:
                if !ok {
                    fmt.Println("PIPELINE: sink errCh closed")
                    sinkErr = nil
                    continue
                }
                if err != nil {
                    fmt.Println("PIPELINE: sink error:", err)
                    return err
                }
            case <- done:
                done = nil
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
