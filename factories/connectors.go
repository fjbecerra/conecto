package factories

import (
	"conecto/auth/credentials"
	"conecto/connectors/api"
	"conecto/connectors/api/graphql"
	"conecto/connectors/api/rest"
	"conecto/core/engines"
	"conecto/core/retry"
	"os"
	"time"
)

type Connector struct {
	config       ConnectorConfig
	streamConfig StreamConfig
	random       retry.Random
	credentialService credentials.CredentialService
	connectorType ConnectorType
	connections Connections

}

func NewConnector(
	config ConnectorConfig, 
	streamConfig StreamConfig, 
	random retry.Random, 
	credentialService credentials.CredentialService,
	connections Connections,
	) *Connector {
		return &Connector{
			config:       config,
			streamConfig: streamConfig,
			random:       random,
			credentialService: credentialService,
			connections: connections,
		}
}

func (c *Connector) Build() engines.ConnectorRunnable {

	var provider api.Provider
	var httpClient api.IClient
	var builder api.RequestBuilder
	var dataExtractor api.DataExtractor
	var cursorExtractor api.CursorExtractor
	var endpointProvider api.EndPointProvider

	switch c.config.ApiConfig.EndpointConfig.EndpointType{
		case DinamicEndpointType:
			endpointProvider = api.NewDinamicEndpointProvider(c.config.ApiConfig.EndpointConfig.Base, c.config.ApiConfig.EndpointConfig.MetadataKeys)
		case StaticEndpointType:
			endpointProvider = api.NewStaticEndpointProvider(c.config.ApiConfig.EndpointConfig.Base)
		default:
			panic("not endpoint type found")
	}

	switch c.config.Type{
		case Api: switch c.config.ApiConfig.Type {
			case Rest:
				provider = buildProvider(c.config.ApiConfig.RestConfig.AuthenticationConfig)
				builder = &rest.RestRequestBuilder{
					EndPointProvider: endpointProvider,
					CursorParam: c.config.ApiConfig.RestConfig.PaginationConfig.Response.Next.Path,
					Method:      "GET",
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				}

				dataExtractor = &rest.RestDataExtractor{
					Path: c.config.ApiConfig.RestConfig.DataConfig.Path,
				}

				cursorExtractor = &api.JSONCursorExtractor{
					Path: c.config.ApiConfig.RestConfig.PaginationConfig.Request.Param,
				}

			case Graphql:
				provider = buildProvider(c.config.ApiConfig.GraphqlConfig.AuthenticationConfig)
				builder = &graphql.GraphQLRequestBuilder{
					EndPointProvider:  endpointProvider,
					Query:    c.streamConfig.Query,
				}

				dataExtractor = &graphql.GraphQLDataExtractor{
					Path: c.config.ApiConfig.GraphqlConfig.DataConfig.Path,
				}

				cursorExtractor = &graphql.GraphQLCursorExtractor{
					HasMorePath: c.config.ApiConfig.GraphqlConfig.PaginationConfig.HasMorePath,
					CursorPath:  c.config.ApiConfig.GraphqlConfig.PaginationConfig.CursorPath,
				}
			default:
				panic("unknown source type: " + c.config.Type)
			}
	}
	

	if c.streamConfig.MockedRestConfig != nil {
		mockedPaths := map[int]string{}
		for i, path := range c.streamConfig.MockedRestConfig.ResponsePaths {
			json, _ := os.ReadFile(path)
			mockedPaths[i] = string(json)
		}
		httpClient = &api.MockHttpClient{
			Calls: mockedPaths,
		}
	} else {
		httpClient = &api.HttpClient{
			Client: c.connections.connections[c.config.ApiConfig.Source].OpenConnection.httpClient,
		}
	}

	client := *api.NewClient(httpClient, provider, c.credentialService)

	paginationProvider := api.PaginationProvider{
		Client:  &client,
		Builder: builder,
		Data:    dataExtractor,
		Cursor:  cursorExtractor,
	}
	connector := &api.HttpConnector{
		Provider: &paginationProvider,
	}

	retryPolicy := retry.Policy{
		MaxRetries:     c.config.Retry.MaxRetries,
		InitialBackoff: time.Duration(c.config.Retry.BackoffMS),
		MaxBackoff:     time.Duration(c.config.Retry.MaxBackoff),
		Jitter:         true,
	}
	retryExecutor := retry.Executor{
		Policy: retryPolicy,
		Random: c.random,
	}

	return &engines.ConnectorEngine{
		Connector: connector,
		Retry:     retryExecutor,
	}

}

func buildProvider(authenticationConfig AuthenticationConfig) api.Provider {
	switch authenticationConfig.Type {
	case Query:
		return &api.QueryProvider{
			Param: authenticationConfig.ParamName,
		}
	case Bearer:
		return &api.BearerProvider{
			Key: authenticationConfig.ParamName,
		}

	case Header:
		return &api.HeaderProvider{
			Name: authenticationConfig.ParamName,
		}

	default:
		panic("not token provider found")
	}
}
