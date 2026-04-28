package factories

import (
	"conecto/core"
	"conecto/core/connectors/rest"
	"conecto/core/connectors"
	"conecto/core/engines"
	"conecto/testutils"
	"net/http"
	"os"
)

type Connector struct {
	Config core.ConnectorConfig
}

func NewConnector(config core.ConnectorConfig) *Connector{
	return &Connector{
		Config: config,
	}
}

func (c *Connector)Build() engines.ConnectorEngine {
	var connector connectors.Connector
	switch c.Config.Type {
	case core.SourceRest:
		connector = buildRest(*c.Config.RestConfig)
	case core.SourceMockedRest:
		connector = buildMockedRest(*c.Config.MockedRestConfig)
	default:
		panic("unknown source type: " + c.Config.Type)
	}
	return engines.ConnectorEngine{
		Connector: connector,
	}

}

func buildRest(config core.RestConfig) *rest.RESTConnector {
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

	return &rest.RESTConnector {
			Provider: &paginationProvider,
	}	
}

func buildMockedRest(config core.MockedRestConfig) *rest.RESTConnector{
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

	return &rest.RESTConnector {
			Provider: &paginationProvider,
	}
}