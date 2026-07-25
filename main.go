package conecto

import (
	"conecto/auth/connections"
	"conecto/auth/credentials"
	"conecto/auth/oauth"
	"conecto/auth/oauth/state"
	"conecto/connectors"
	"conecto/connectors/shopify"
	"conecto/core/retry"
	"conecto/factories"
	"conecto/pipelines"
	"conecto/sync"
	"context"
	"encoding/base64"
	"net/http"
	"time"
)

func main() {

	setupRouter()

	// // Database
	// db := postgres.Connect(...)

	// // Stores
	// connectionStore := connections.NewPostgresStore(db)
	// credentialStore := credentials.NewAESStore(...)
	// stateStore := oauth2.NewPostgresStateStore(db)

	// // Connectors
	// shopify := shopify.New(...)
	// meta := meta.New(...)
	// google := googleads.New(...)

	// connectors := map[string]connector.Connector{
	// 	shopify.Name(): shopify,
	// 	meta.Name(): meta,
	// 	google.Name(): google,
	// }

	// // OAuth service
	// oauthService := oauth2.NewService(
	// 	connectionStore,
	// 	credentialStore,
	// 	stateStore,
	// 	connectors,
	// )

	// // HTTP handlers
	// oauthHandler := oauth2.NewHandler(oauthService)

	// // Router
	// router := NewRouter(oauthHandler)

	// http.ListenAndServe(":8080", router)
}



// func setupPipelines() pipeline.Registry {

// 	registry := pipeline.NewRegistry()


// 	shopify := loadPipeline(
// 		"./pipelines/shopify.json",
// 	)


// 	registry.Register(shopify)


// 	return registry
// }

func setupSync() {
	queue := sync.NewQueue(100)

    worker := sync.NewWorker(
        queue,
    )

    scheduler := sync.NewScheduler(syncService)

    go worker.Run(context.Background())
    go scheduler.Run(context.Background())
}

func setupRouter() {

	connectionStore := connections.NewPostgresStore("postgres://user:password@localhost:5432/mydb?sslmode=disable")
	
	shopifyConnector := shopify.Connector{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			Scopes:       []string{"read_orders"},
			HttpClient: &http.Client{
    			Timeout: 30 * time.Second,
			},
	}		
	registry := connectors.NewRegistry(shopifyConnector)

	queue := sync.NewQueue(10)
	config:= factories.LoadConfigPipeline("./testdata/orders_pipeline_flattened_data_to_memory.json")
	pipeline:= factories.BuildPipeline(config)
	pipeRegistry:= pipelines.NewRegistry()
	pipeRegistry.Register(pipeline)
	jobRepository:= sync.NewPostgresJobRepository()
	retryPolicy:= retry.Policy{
		MaxRetries: 10,
		InitialBackoff: 10,
		MaxBackoff: time.Duration(10),
		Jitter: true,
	}
	retryExecutor := retry.Executor {
		Policy: retryPolicy,
		Random: x,
	}

	syncService := sync.NewSyncService(queue, pipeRegistry, connectionStore, retryExecutor)

	stateSigner := state.NewHMACStateSigner(
			[]byte("z1Wbj51mq1GmIpmwAfLv9X5oSekOYEsC/9YXhOCuKjU="),
			10*time.Minute)
	
	v1, _ := base64.StdEncoding.DecodeString(
		"z1Wbj51mq1GmIpmwAfLv9X5oSekOYEsC/9YXhOCuKjU=",
	)
	keys:=map[string][]byte{
		"v1": v1,
	}
	keyManager:=credentials.NewStaticKeyManager(keys, "v1")

	credentialStore := credentials.NewPostgresCredentialDB("postgres://user:password@localhost:5432/mydb?sslmode=disable")
	
	credentialService := credentials.NewAESGCMCredentialService(credentialStore, keyManager)

	service := oauth.NewService(
				connectionStore,
				credentialService,
				stateSigner,
				*registry, *syncService)


	handler := oauth.NewHandler(service)

	router:=oauth.NewRouter(handler)
	http.ListenAndServe(":8080", router)
}