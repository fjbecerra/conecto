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
)

type Pipeline interface {
    Run(ctx context.Context) error
}

type PipelineFactory func() Pipeline

type Registry struct {
    Factories map[string]PipelineFactory
}

type JsonPipeline struct {
    connector sources.Source[json.RawMessage]
    extractor extractors.Extractor[json.RawMessage]
    sink      core.Sink
}

func (p *JsonPipeline) Run(ctx context.Context) error {
    // 1. Connector
    out, errCh1 := p.connector.Fetch(ctx)

    // 2. Extractor
    recordsCh, errCh2 := p.extractor.Extract(ctx, out)

    // 3. Sink
    errCh3 := p.sink.Write(ctx, recordsCh)

    // 4. Handle errors
    for errCh1 != nil || errCh2 != nil || errCh3 != nil {
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

        case err, ok := <-errCh3:
            if !ok {
                errCh3 = nil
                continue
            }
            if err != nil {
                return err
            }
        }
    }

    return nil
}

//i may put the config outside and discover them by looping 
func NewRegistry() *Registry {
    return &Registry{
        Factories: map[string]PipelineFactory{
            "facebookAdInsight": func() Pipeline {
                return NewPipeline("./schemas/facebook_ad_insight.json")
            },
        },
    }
}

func NewPipeline(configPath string) *JsonPipeline {
    //Intialize http
    //Initiazize restclient
    //initialize pagination provider
    //initialize connector
    client := rest.NewRestClient(http.DefaultClient)
	paginationProvider := rest.NewPaginationProvider(
		client,
		configPath, 
		"http://base-url")

	connector := Connector {
		Provider: &paginationProvider,
	}

    return &JsonPipeline{
        connector: connectors.NewJsonRestConnector(
            http.DefaultClient, //todo set more advance options
            "url",
            config),
        extractor: extractors.NewJsonExtractor(config),
        sink:      sinks.NewMemorySink(),
    }
}



