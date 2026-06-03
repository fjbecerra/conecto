package factories

import (
	"conecto/connectors/rest"
	"conecto/connectors/rest/auths"
	"conecto/connectors/rest/auths/stores"
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
	tokenProvider := buildTokenProvider(c.config.RestConfig.AuthenticationConfig)
	store := buildStore(c.config.RestConfig.TokenStoreConfig, c.connections)
	v1, _ := base64.StdEncoding.DecodeString(
		os.Getenv("TOKEN_ENCRYPTION_KEY_V1"),
	)
	keys:=map[string][]byte{
		"v1": v1,
	}
	keyManager:=auths.NewStaticKeyManager(keys, "v1")
	tokenStore:= auths.NewADBTokenStore(store, keyManager)

	var httpClient rest.IClient
	switch c.config.Type {
		case Rest:
			httpClient= &rest.HttpClient{
				Client: http.DefaultClient,
			}	
		case MockedRest:
			mockedPaths := map[int]string{}
			for i, path := range c.config.MockedRestConfig.ResponsePaths {
				json,_ := os.ReadFile(path)
				mockedPaths[i] = string(json)
			}
			httpClient = &rest.MockHttpClient{
				Calls: mockedPaths,
			}
		default:
			panic("unknown source type: " + c.config.Type)
	}
	restClient := *rest.NewRestClient(httpClient, tokenProvider, tokenStore)
	paginationProvider := rest.PaginationProvider{
		RestClient : restClient,
		BaseUrl: c.config.RestConfig.BaseUrl,
		DataPath: c.config.RestConfig.DataConfig.Path,
		ResponseNextPath: c.config.RestConfig.PaginationConfig.Response.Next.Path,
		RequestParam: c.config.RestConfig.PaginationConfig.Request.Param,
	}
	connector := &rest.RESTConnector{
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
				ParamName: authenticationConfig.QueryToken.ParamName,
			}
		case "bearer":
			return &auths.BearerTokenProvider{}	
		
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
