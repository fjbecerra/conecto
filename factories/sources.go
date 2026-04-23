package factories

import (
	"conecto/core"
	"conecto/core/sources/rest"
	"conecto/testutils"
	"net/http"
	"os"
)

type Source struct {
	Config core.SourceConfig
}

func NewSource(config core.SourceConfig) *Source{
	return &Source{
		Config: config,
	}
}

func (source *Source)Build() any {
	switch source.Config.Type {
	case core.SourceRest:
		return buildRest(*source.Config.RestConfig)
	case core.SourceMockedRest:
		return buildMockedRest(*source.Config.MockedRestConfig)
	default:
		panic("unknown source type: " + source.Config.Type)
	}
}

func buildRest(config core.RestConfig) *rest.Connector {
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

func buildMockedRest(config core.MockedRestConfig) *rest.Connector{
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