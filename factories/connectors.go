package factories

import (
	"conecto/connectors/_http"
	"conecto/connectors/_http/auths"
	"conecto/connectors/_http/auths/stores"
	"conecto/connectors/_http/graphql"
	"conecto/connectors/_http/rest"
	"conecto/core/engines"
	"conecto/core/retry"
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Connector struct {
	config ConnectorConfig
	random retry.Random
	connections Connections

}

func NewConnector(config ConnectorConfig, random retry.Random, connections Connections) *Connector{
	return &Connector{
		config: config,
		random: random,
		connections: connections,
	}
}

func (c *Connector)Build() engines.ConnectorRunnable {
	
	store := buildStore(c.config.RestConfig.TokenStoreConfig, c.connections)
	v1, _ := base64.StdEncoding.DecodeString(
		os.Getenv("TOKEN_ENCRYPTION_KEY_V1"),
	)
	keys:=map[string][]byte{
		"v1": v1,
	}
	keyManager:=auths.NewStaticKeyManager(keys, "v1")
	tokenStore:= auths.NewADBTokenStore(store, keyManager)
	var tokenProvider auths.TokenProvider
	var httpClient _http.IClient
	var builder _http.RequestBuilder
	var dataExtractor _http.DataExtractor
	var cursorExtractor _http.CursorExtractor

	switch c.config.Type {
		case Rest:
			tokenProvider = buildTokenProvider(c.config.RestConfig.AuthenticationConfig)
			httpClient= &_http.HttpClient{
				Client: http.DefaultClient,
			}	
			
			builder = &rest.RestRequestBuilder{
				BaseURL: c.config.RestConfig.BaseUrl,
				CursorParam: c.config.RestConfig.PaginationConfig.Response.Next.Path,
				Method: "GET",
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			}
			
			dataExtractor = &rest.RestDataExtractor{
					Path: c.config.RestConfig.DataConfig.Path,
			}

			cursorExtractor = &_http.JSONCursorExtractor{
				Path:  c.config.RestConfig.PaginationConfig.Request.Param,
			}	
			
		case Graphql:
			tokenProvider = buildTokenProvider(c.config.GraphqlConfig.AuthenticationConfig)
			builder = &graphql.GraphQLRequestBuilder{
				Endpoint: c.config.GraphqlConfig.BaseUrl,
				Query: c.config.GraphqlConfig.Query,
					
			}

			dataExtractor = &graphql.GraphQLDataExtractor{
				Path: c.config.GraphqlConfig.DataConfig.Path,
			}

			cursorExtractor = &graphql.GraphQLCursorExtractor{
				HasMorePath: c.config.GraphqlConfig.PaginationConfig.HasMorePath,
				CursorPath: c.config.GraphqlConfig.PaginationConfig.CursorPath,
			}

		case MockedRest:
			mockedPaths := map[int]string{}
			for i, path := range c.config.MockedRestConfig.ResponsePaths {
				json,_ := os.ReadFile(path)
				mockedPaths[i] = string(json)
			}
			httpClient = &_http.MockHttpClient{
				Calls: mockedPaths,
			}
		default:
			panic("unknown source type: " + c.config.Type)
	}
	client := *_http.NewClient(httpClient, tokenProvider, tokenStore)
	
	paginationProvider := _http.PaginationProvider{
		Client : &client,
		Builder: builder,
		Data: dataExtractor,
		Cursor: cursorExtractor,
  }
	connector := &_http.HttpConnector{
		Provider: &paginationProvider,
	}

	retryPolicy:= retry.Policy{
		MaxRetries: c.config.Retry.MaxRetries,
		InitialBackoff: time.Duration(c.config.Retry.BackoffMS),
		MaxBackoff: time.Duration(c.config.Retry.MaxBackoff),
		Jitter: true,
	}
	retryExecutor := retry.Executor {
		Policy: retryPolicy,
		Random: c.random,
	}

	return &engines.ConnectorEngine{
		Connector: connector,
		Retry: retryExecutor,
	}

}

func buildTokenProvider(authenticationConfig AuthenticationConfig) auths.TokenProvider{	
	switch authenticationConfig.Type {
		case "query":
			return &auths.QueryTokenProvider{
				ParamName: authenticationConfig.ParamName,
			}
		case "bearer":
			return &auths.BearerTokenProvider{}	
		
		case "header":
			return &auths.HeaderTokenProvider{
				HeaderName: authenticationConfig.ParamName,
			}
		
		default:
			panic("not token provider found")
	}
}

func buildStore(tokenStoreConfig TokenStoreConfig, connections Connections) stores.Store{
	
	switch tokenStoreConfig.Type{
		case PostgresTokenStore:
			connection:= connections[tokenStoreConfig.Source].OpenDB()
			if(tokenStoreConfig.AutoCreate){
				createPostgresTokenStoreTable(tokenStoreConfig, connection)			}
			
			return stores.NewPostgresTokenDB(connection)
		case MemoryTokenStore:
			return stores.NewMemoryStoreToken(make(map[string]any))
		default:
			panic("not token store found")
	}
}

func createPostgresTokenStoreTable(tokenStoreConfig TokenStoreConfig, db *sql.DB){
	query := `
		CREATE TABLE IF NOT EXISTS %s (
			pipeline_id    TEXT NOT NULL,

			ciphertext     BYTEA NOT NULL,
			nonce          BYTEA NOT NULL,

			key_version    TEXT NOT NULL,

			expires_at     TIMESTAMPTZ,

			created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),

			PRIMARY KEY (pipeline_id)
	);
	`
	_, err := db.Exec(fmt.Sprintf(query, tokenStoreConfig.Name))
	if err != nil {
		panic(err)
	}

	fmt.Println("table created or already exists")
}
