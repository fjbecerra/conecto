package factories

import (
	"conecto/auth/oauth"
	"conecto/auth/oauth/state"
	"conecto/connectors"
	"conecto/connectors/shopify"
	"conecto/core/pipelines"
	"conecto/core/retry"
	"conecto/sync"
	"context"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
)

type IRunner interface {
	Closed() error
	RunScheduler(ctx context.Context)
	RunWorker(ctx context.Context) 
	RunRouter() *chi.Mux
}

type Runner struct {
	scheduler *sync.Scheduler
	worker *sync.Worker
	router *chi.Mux
	connections Connections
}

func (c *Runner) Closed() error{
	err := c.connections.CloseAll()
	if err != nil {
		return err
	}
	return nil
}

func (c *Runner) RunScheduler(ctx context.Context){
	 c.scheduler.Run(ctx)
}

func (c *Runner) RunWorker(ctx context.Context){
	 c.worker.Run(ctx)
}

func (c *Runner) RunRouter() *chi.Mux{
	return c.router
}


type Conecto struct{
	
	conectoConfig ConectoConfig
}

func NewConecto(conectoConfig ConectoConfig) *Conecto{
	return &Conecto {
		conectoConfig: conectoConfig,
	}
}

func (c *Conecto) Build() IRunner{
	random:= &RandomImpl{}
	connections:= NewSource(c.conectoConfig.SourcesConfig).Build()
	stores:= NewConectoStore(c.conectoConfig.DBConfig, connections).Build()
	stateStore := stores.stateStore
	credentialService:=  NewCredentialService(stores.credentialStore).Build()
	pipelineRegitry := pipelines.NewRegistry()
	connectorRegistry := connectors.NewRegistry()

	for _, path := range c.conectoConfig.PipelineRegistryConfig {
		pipelineConfig, error := LoadConfig[PipelineConfig](path)
		var connector connectors.Connector
		switch pipelineConfig.ID {
			case "shopify" :
				connector = &shopify.Connector{
					ClientID:     pipelineConfig.AuthorizeConfig.Oauth.ClientId,
					ClientSecret: os.Getenv(pipelineConfig.AuthorizeConfig.Oauth.ClientSecret),
					Scopes:       pipelineConfig.AuthorizeConfig.Oauth.Scopes,
					AppUrl: 	  pipelineConfig.AuthorizeConfig.Oauth.AppUrl,	
					EndpointApiProvider: &shopify.ShopifyEndpointProvider{},
					ResponseProvider: &shopify.ShopifyResponseProvider{},
					HttpClient: &http.Client{},
				}		
				connectorRegistry.Register(connector)
		}
		if(error!=nil){
			panic("path not found")
		}
		
		pipeline:= NewPipeline(connector, connections, random, stateStore, credentialService ,pipelineConfig).Build()
		pipelineRegitry.Register(pipeline)
	}


	var buffer sync.Buffer
	if(c.conectoConfig.SyncConfig.Buffer.BufferType == QueueType){
		buffer = sync.NewQueue(c.conectoConfig.SyncConfig.Buffer.Size)
	}else{
		panic("unknown buffer type")
	}


	retryPolicy := retry.Policy{
		MaxRetries:     c.conectoConfig.SyncConfig.Retry.MaxRetries,
		InitialBackoff: time.Duration(c.conectoConfig.SyncConfig.Retry.BackoffMS),
		MaxBackoff:     time.Duration(c.conectoConfig.SyncConfig.Retry.MaxBackoff),
		Jitter:         true,
	}
	retryExecutor := retry.Executor{
		Policy: retryPolicy,
		Random: random,
	}

	syncService := sync.NewSyncService(
		buffer, 
		pipelineRegitry, 
		stores.connectionStore, 
		stores.jobRepository,
		retryExecutor,
	)
	duration,_:= time.ParseDuration(c.conectoConfig.SyncConfig.Scheduler.Duration)
	scheduler := sync.NewScheduler(duration,syncService )
	worker := sync.NewWorker(buffer, syncService)
	stateSigner := state.NewHMACStateSigner(
			[]byte(os.Getenv("AUTH_STATE_SIGNER_KEY")),
			10*time.Minute)
	
	oauthService := oauth.NewService(
		stores.connectionStore, 
		credentialService, 
		stateSigner,
		connectorRegistry, 
		syncService,
	)

	handler := oauth.NewHandler(oauthService)
	router :=oauth.NewRouter(handler)

	return &Runner{
		scheduler: &scheduler,
		worker: worker,
		router: router,
		connections: connections,
	}
}