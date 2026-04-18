package pipelines

import (
	"conecto/core"
	"conecto/core/extractors"
	"conecto/core/sources"
	"conecto/core/sources/rest"
	"conecto/sinks"
	"context"
	"encoding/json"
	"net/http"
	"runtime"
)

type Pipeline[T any] struct {
    source sources.Source[T]
    sink   core.Sink[T]
}

type PipelineFactory func() Pipeline[core.Record]


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

func NewPipeline(configPath string) *Pipeline[core.Record]{
	client := rest.NewRestClient(http.DefaultClient)
	paginationProvider := rest.NewPaginationProvider(
		client,
		configPath)

	connector := rest.Connector {
		Provider: &paginationProvider,
	}

	extractor := extractors.NewJsonExtractor(configPath)


	source := &core.MapSource[json.RawMessage, core.Record]{
		Upstream: &connector,
		MapFn: func(raw json.RawMessage) core.Record {
			record,_ := extractor.Extract(raw)
			return record
		},
		Workers: runtime.NumCPU(),
	}

	sink := sinks.NewMemorySink[core.Record]()

	return &Pipeline[core.Record]{
		source: source,
		sink:   sink,
	}
}






