package factories

import (
	"conecto/auth/credentials"
	"conecto/connectors/api"
	"conecto/connectors/api/graphql"
	"conecto/connectors/api/rest"
	"conecto/core/engines"
	"conecto/core/retry"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Connector struct {
	config       ConnectorConfig
	streamConfig StreamConfig
	random       retry.Random
	credentialService credentials.CredentialService
}

func NewConnector(
	config ConnectorConfig, 
	streamConfig StreamConfig, 
	random retry.Random, 
	credentialService credentials.CredentialService,
	) *Connector {
		return &Connector{
			config:       config,
			streamConfig: streamConfig,
			random:       random,
			credentialService: credentialService,
		}
}

func (c *Connector) Build() engines.ConnectorRunnable {

	// v1, _ := base64.StdEncoding.DecodeString(
	// 	os.Getenv("TOKEN_ENCRYPTION_KEY_V1"),
	// )
	// keys := map[string][]byte{
	// 	"v1": v1,
	// }
	//keyManager := credentials.NewStaticKeyManager(keys, "v1")

	var provider api.Provider
	var httpClient api.IClient
	var builder api.RequestBuilder
	var dataExtractor api.DataExtractor
	var cursorExtractor api.CursorExtractor
	//var store credentials.Store

	switch c.config.Type{
		case Api: switch c.config.ApiConfig.Type {
			case Rest:
				//store = buildStore(c.config.ApiConfig.RestConfig.TokenStoreConfig, c.connections)
				provider = buildProvider(c.config.ApiConfig.RestConfig.AuthenticationConfig)
				builder = &rest.RestRequestBuilder{
					BaseURL:     c.config.ApiConfig.RestConfig.BaseUrl,
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
				//store = buildStore(c.config.ApiConfig.GraphqlConfig.TokenStoreConfig, c.connections)
				provider = buildProvider(c.config.ApiConfig.GraphqlConfig.AuthenticationConfig)
				builder = &graphql.GraphQLRequestBuilder{
					Endpoint: c.config.ApiConfig.GraphqlConfig.BaseUrl,
					Query:    c.config.ApiConfig.GraphqlConfig.Query,
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
			Client: http.DefaultClient,
		}
	}
	//credentialService := credentials.NewAESGCMCredentialService(store, keyManager)

	


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

func buildStore(tokenStoreConfig TokenStoreConfig, connections Connections) credentials.Store {

	switch tokenStoreConfig.Type {
	case PostgresTokenStore:
		connection := connections[tokenStoreConfig.Source].OpenDB()
		if tokenStoreConfig.AutoCreate {
			createPostgresTokenStoreTable(tokenStoreConfig, connection)
		}
		return credentials.NewPostgresCredentialDB(connection)
	case MemoryTokenStore:
		return credentials.NewMemoryStoreCredential(make(map[string]any))
	default:
		panic("not token store found")
	}
}

func createPostgresTokenStoreTable(tokenStoreConfig TokenStoreConfig, db *sql.DB) {
	query := `
		CREATE TABLE IF NOT EXISTS %s (
			connection_id    TEXT NOT NULL,

			ciphertext     BYTEA NOT NULL,
			nonce          BYTEA NOT NULL,

			key_version    TEXT NOT NULL,

			expires_at     TIMESTAMPTZ,

			created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),

			PRIMARY KEY (connection_id)
	);
	`
	_, err := db.Exec(fmt.Sprintf(query, tokenStoreConfig.Name))
	if err != nil {
		panic(err)
	}

	fmt.Println("table created or already exists")
}
