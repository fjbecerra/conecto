package pipelines

import (
	"conecto/core"
	"conecto/core/sinks"
	"conecto/core/sources"
	"conecto/core/sources/rest"
	"conecto/core/transforms"
	"conecto/testutils"
	"context"
	"encoding/json"
	"net/http"
	"os"
)

type Pipeline[T any] struct {
    source sources.Source[T]
    sink   sinks.Sink[T]
}

type PipelineRunner interface {
	Run(ctx context.Context) error	
	TestResult() []core.Record //only for tests porpuses
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

//only for tests purposes
func (p *Pipeline[T]) TestResult() []T {
	return p.sink.(*sinks.SinkMemory[T]).Data()	
}


func BuildSource(config core.SourceConfig) any {
	switch config.Type {
	case core.SourceRest:
		return buildRestSource(*config.RestConfig)
	case core.SourceMockedRest:
		return buildMockedRestSource(*config.MockedRestConfig)
	default:
		panic("unknown source type: " + config.Type)
	}
}

func buildRestSource(config core.RestConfig) *rest.Connector {
	var tokenProvider rest.TokenProvider
	switch config.Authentication.Type {
		case "query":
			tokenProvider = &rest.QueryTokenProvider{
				ParamName: config.Authentication.QueryToken.ParamName,
			}
		case "bearer":
			tokenProvider = &rest.BearerTokenProvider{}	
	}
	client := rest.NewRestClient(http.DefaultClient, tokenProvider)
	paginationProvider := rest.PaginationProvider{
		Client :client,
		BaseUrl: config.BaseUrl,
		DataPath: config.BaseRestConfig.Data.Path,
		ResponseNextPath: config.BaseRestConfig.Pagination.Response.Next.Path,
		RequestParam: config.BaseRestConfig.Pagination.Request.Param,
	}

	return &rest.Connector {
			Provider: &paginationProvider,
	}	
}

func buildMockedRestSource(config core.MockedRestConfig) *rest.Connector{
	mockedPaths := map[int]string{}
	for i, path := range config.ResponsePaths {
		json,_ := os.ReadFile(path)
		mockedPaths[i] = string(json)
	}
	
	paginationProvider := rest.PaginationProvider{
		Client : &testutils.MockClient{
			Calls: mockedPaths,
		},
		DataPath: config.BaseRestConfig.Data.Path,
		ResponseNextPath: config.BaseRestConfig.Pagination.Response.Next.Path,
		RequestParam: config.BaseRestConfig.Pagination.Request.Param,
	}

	return &rest.Connector {
			Provider: &paginationProvider,
	}
}

func BuildSink(config core.SinkConfig) sinks.Sink[core.Record] {
	switch config.Type {
		case core.SinkMemory:
			return BuildMemorySink()
		default:
			panic("unknown source type: " + config.Type)
	}
}

func BuildMemorySink() *sinks.SinkMemory[core.Record] {
	return sinks.NewMemorySink[core.Record]()
}

func BuildTransform (source any, config []core.TransformConfig) sources.Source[core.Record] {
	var current sources.Source[core.Record]

	for _, t := range config {
		switch src := source.(type) {
			case sources.Source[json.RawMessage]:
				var fn  core.Transform[json.RawMessage,core.Record]
				switch t.Type {
					case core.TransformExtractor :
						fn = BuildExtractor(*t.ExtractorConfig)
					default:
						panic("unsupported transform")
				}
				current = &core.MapSource[json.RawMessage, core.Record]{
					Upstream: src,
					MapFn: fn,
					Workers: t.Workers,
				}
			default:
				panic("unsupported source type")
		}
	}
	return current
}

func BuildExtractor(config core.ExtractorConfig) core.Transform[json.RawMessage, core.Record] {
	extractor :=transforms.Extractor{
					Fields : transforms.Fields(config.FieldsConfig),
				}
    return  func (in json.RawMessage) core.Record {
		record, _ := extractor.Extract(in)
		return record
	}  
}

func BuildPipeline(configPath string) PipelineRunner{
	config := core.LoadConfigPipeline(configPath)
	source := BuildSource(config.SourceConfig)
	current := BuildTransform(source, config.TransformsConfig)
	sink := BuildSink(config.SinkConfig)

	return &Pipeline[core.Record]{
		source: current,
		sink:   sink,
	}
}

type PipelineFactory func() PipelineRunner
