package factories

import (
	"conecto/auth/credentials"
	"conecto/connectors/_http"
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
	
	
	v1, _ := base64.StdEncoding.DecodeString(
		os.Getenv("TOKEN_ENCRYPTION_KEY_V1"),
	)
	keys:=map[string][]byte{
		"v1": v1,
	}
	keyManager:=credentials.NewStaticKeyManager(keys, "v1")
	
	var provider _http.Provider
	var httpClient _http.IClient
	var builder _http.RequestBuilder
	var dataExtractor _http.DataExtractor
	var cursorExtractor _http.CursorExtractor
	var store credentials.Store

	switch c.config.Type {
		case Rest:
			store = buildStore(c.config.RestConfig.TokenStoreConfig, c.connections)
			provider = buildProvider(c.config.RestConfig.AuthenticationConfig)
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
			store = buildStore(c.config.GraphqlConfig.TokenStoreConfig, c.connections)
			provider = buildProvider(c.config.GraphqlConfig.AuthenticationConfig)
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
		default:
			panic("unknown source type: " + c.config.Type)
	}
	
	if c.config.MockedRestConfig != nil {
		mockedPaths := map[int]string{}
		for i, path := range c.config.MockedRestConfig.ResponsePaths {
			json,_ := os.ReadFile(path)
			mockedPaths[i] = string(json)
		}
		httpClient = &_http.MockHttpClient{
			Calls: mockedPaths,
		}
	}else{
		httpClient= &_http.HttpClient{
			Client: http.DefaultClient,
		}
	}
	credentialService:= credentials.NewAESGCMCredentialService(store, keyManager)
	client := *_http.NewClient(httpClient, provider, credentialService)
	
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

func buildProvider(authenticationConfig AuthenticationConfig) _http.Provider{	
	switch authenticationConfig.Type {
		case Query:
			return &_http.QueryProvider{
				Param: authenticationConfig.ParamName,
			}
		case Bearer:
			return &_http.BearerProvider{
				Key: authenticationConfig.ParamName,
			}
		
		case Header:
			return &_http.HeaderProvider{
				Name: authenticationConfig.ParamName,
			}
		
		default:
			panic("not token provider found")
	}
}

func buildStore(tokenStoreConfig TokenStoreConfig, connections Connections) credentials.Store{
	
	switch tokenStoreConfig.Type{
		case PostgresTokenStore:
			connection:= connections[tokenStoreConfig.Source].OpenDB()
			if(tokenStoreConfig.AutoCreate){
				createPostgresTokenStoreTable(tokenStoreConfig, connection)			}
			
			return credentials.NewPostgresCredentialDB(connection)
		case MemoryTokenStore:
			return credentials.NewMemoryStoreCredential(make(map[string]any))
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
