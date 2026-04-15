package pipelines

import (
	"conecto/connectors"
	"conecto/core"
	"conecto/core/extractors"
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
    connector core.Connector[json.RawMessage]
    extractor core.Extractor[json.RawMessage]
    sink      core.Sink
}

func (p *JsonPipeline) Run(ctx context.Context) error {
    // 1. Connector
    out, errCh1 := p.connector.Run(ctx)

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

func NewRegistry() *Registry {
    return &Registry{
        Factories: map[string]PipelineFactory{
            "facebookAdInsight": func() Pipeline {
                return NewFacebookAdInsightPipeline()
            },
        },
    }
}

func NewFacebookAdInsightPipeline() *JsonPipeline {
    schema := core.LoadSchema("./schemas/facebook_ad_insight.json")

    return &JsonPipeline{
        connector: connectors.NewJsonRestConnector(
            http.DefaultClient, //todo set more advance options
            "url",
            schema),
        extractor: extractors.NewJsonExtractor(schema),
        sink:      sinks.NewMemorySink(),
    }
}



