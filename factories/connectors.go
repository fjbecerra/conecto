package factories

import (
	"conecto/core/connectors/rest"
	"conecto/core/connectors/rest/auths"
	"conecto/core/connectors/rest/auths/stores"
	"conecto/core/engines"
	"conecto/core/retry"
	"conecto/testutils"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Connector struct {
	Config ConnectorConfig
	RandFn func() float64
	DBConnection DBConnection

}

func NewConnector(config ConnectorConfig, randFn func() float64, dbConnection DBConnection) *Connector{
	return &Connector{
		Config: config,
		RandFn: randFn,
		DBConnection: dbConnection,
	}
}

func (c *Connector)Build() engines.ConnectorEngine {
	tokenProvider := buildTokenProvider(c.Config.RestConfig.AuthenticationConfig)
	store := buildStore(c.Config.RestConfig.TokenStoreConfig, c.DBConnection)
	v1, _ := base64.StdEncoding.DecodeString(
		os.Getenv("TOKEN_ENCRYPTION_KEY_V1"),
	)
	keys:=map[string][]byte{
		"v1": v1,
	}
	keyManager:=auths.NewStaticKeyManager(keys, "v1")
	tokenStore:= auths.NewADBTokenStore(store, keyManager)

	var httpClient rest.IClient
	switch c.Config.Type {
		case Rest:
			httpClient= &rest.HttpClient{
				Client: http.DefaultClient,
			}	
		case MockedRest:
			mockedPaths := map[int]string{}
			for i, path := range c.Config.MockedRestConfig.ResponsePaths {
				json,_ := os.ReadFile(path)
				mockedPaths[i] = string(json)
			}
			httpClient = &testutils.MockHttpClient{
				Calls: mockedPaths,
			}
		default:
			panic("unknown source type: " + c.Config.Type)
	}
	restClient := *rest.NewRestClient(httpClient, tokenProvider, tokenStore)
	paginationProvider := rest.PaginationProvider{
		RestClient : restClient,
		BaseUrl: c.Config.RestConfig.BaseUrl,
		DataPath: c.Config.RestConfig.DataConfig.Path,
		ResponseNextPath: c.Config.RestConfig.PaginationConfig.Response.Next.Path,
		RequestParam: c.Config.RestConfig.PaginationConfig.Request.Param,
	}
	connector := &rest.RESTConnector{
		Provider: &paginationProvider,
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

func buildStore(tokenStoreConfig TokenStoreConfig, dbConnection DBConnection) stores.Store{
	
	switch tokenStoreConfig.Type{
		case PostgresTokenStore:
			if(tokenStoreConfig.AutoCreate){
				createTokenStoreTable(tokenStoreConfig, dbConnection)
			}
			return stores.NewPostgresTokenDB(dbConnection.DB)
		case MemoryTokenStore:
			return stores.NewMemoryStoreToken(make(map[string]any))
		default:
			panic("not token store found")
	}
}

func createTokenStoreTable(tokenStoreConfig TokenStoreConfig, dbConnection DBConnection){
	query := `
		CREATE TABLE IF NOT EXISTS %s (
			pipeline_id    TEXT NOT NULL,
			provider       TEXT NOT NULL,

			ciphertext     BYTEA NOT NULL,
			nonce          BYTEA NOT NULL,

			key_version    TEXT NOT NULL,

			expires_at     TIMESTAMPTZ,

			created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),

			PRIMARY KEY (pipeline_id, provider)
	);
	`
	_, err := dbConnection.DB.Exec(fmt.Sprintf(query, tokenStoreConfig.Name))
	if err != nil {
		panic(err)
	}

	fmt.Println("table created or already exists")
}
