package pipelines

import (
	"conecto/core"
	"conecto/core/transforms"
	"conecto/core/sources"
	"conecto/core/sources/rest"
	"conecto/core/sinks"
	"context"
	"encoding/json"
	"net/http"
	"runtime"
)

type Pipeline[T any] struct {
    source sources.Source[T]
    sink   sinks.Sink[T]
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
	config := core.LoadConfig(configPath)
	client := rest.NewRestClient(http.DefaultClient)
	paginationProvider := rest.PaginationProvider{
		Client :client,
		BaseUrl: config.BaseUrl,
		DataPath: config.Data.Path,
		ResponseNextPath: config.Pagination.Response.Next.Path,
		RequestParam: config.Pagination.Request.Param,
	}

	connector := rest.Connector {
		Provider: &paginationProvider,
	}

	extractor := transforms.Extractor{
		Fields: transforms.Fields(config.Data.FieldsConfig.Fields),
	}


	source := &transforms.MapSource[json.RawMessage, core.Record]{
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






