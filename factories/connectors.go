package factories

import (
	"conecto/core/connectors"
	"conecto/core/connectors/rest"
	"conecto/core/engines"
	"conecto/core/idempotency"
	"conecto/core/retry"
	"conecto/testutils"
	"net/http"
	"os"
	"time"
)

type Connector struct {
	Config ConnectorConfig
	RandFn func() float64
}

func NewConnector(config ConnectorConfig, randFn func() float64) *Connector{
	return &Connector{
		Config: config,
		RandFn: randFn,
	}
}

func (c *Connector)Build() engines.ConnectorEngine {
	var connector connectors.Connector
	switch c.Config.Type {
	case SourceRest:
		connector = buildRest(*c.Config.RestConfig)
	case SourceMockedRest:
		connector = buildMockedRest(*c.Config.MockedRestConfig)
	default:
		panic("unknown source type: " + c.Config.Type)
	}

	retryPolicy:= retry.Policy{
		MaxRetries: c.Config.Retry.MaxRetries,
		InitialBackoff: time.Duration(c.Config.Retry.BackoffMS),
		MaxBackoff: time.Duration(c.Config.Retry.MaxBackoff),
		Jitter: true,
	}
	retryExecutor := retry.Executor {
		Policy: retryPolicy,
		Rand: c.RandFn,
	}

	return engines.ConnectorEngine{
		Connector: connector,
		Retry: retryExecutor,
	}

}

func buildRest(config RestConfig) *rest.RESTConnector {
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
	generator:= idempotency.HashGenerator{}
	return &rest.RESTConnector {
			Provider: &paginationProvider,
			Generator: &generator,
	}	
}

func buildMockedRest(config MockedRestConfig) *rest.RESTConnector{
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
	generator:= idempotency.HashGenerator{}
	return &rest.RESTConnector {
			Provider: &paginationProvider,
			Generator: &generator,
	}
}