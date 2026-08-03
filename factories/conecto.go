package factories

import (
	"conecto/auth/oauth"
	"conecto/auth/oauth/state"
	"conecto/connectors"
	"conecto/connectors/shopify"
	"conecto/core/pipelines"
	"conecto/core/retry"
	"conecto/sync"
	"net/http"
	"os"
	"time"
	"github.com/go-chi/chi/v5"


)

type Runner struct {
	Scheduler *sync.Scheduler
	Worker *sync.Worker
	Router *chi.Mux
}

type Conecto struct{
	
	conectoConfig ConectoConfig
}

func NewConecto(conectoConfig ConectoConfig) *Conecto{
	return &Conecto {
		conectoConfig: conectoConfig,
	}
}

func (c *Conecto) Build() Runner{
	random:= &RandomImpl{}
	connections:= NewSource(c.conectoConfig.SourcesConfig).Build()
	stores:= NewConectoStore(c.conectoConfig.DBConfig, connections).Build()
	stateStore := stores.stateStore
	credentialService:=  NewCredentialService(stores.credentialStore).Build()
	pipelineRegitry := pipelines.NewRegistry()
	connectorRegistry := connectors.NewRegistry()

	for _, path := range c.conectoConfig.PipelineRegistryConfig {
		pipelineConfig, error := LoadConfig[PipelineConfig](path)

		switch pipelineConfig.ID {
			case "shopify" :
				shopifyConnector := shopify.Connector{
					ClientID:     pipelineConfig.AuthorizeConfig.Oauth.ClientId,
					ClientSecret: os.Getenv(pipelineConfig.AuthorizeConfig.Oauth.ClientSecret),
					Scopes:       pipelineConfig.AuthorizeConfig.Oauth.Scopes,
					AppUrl: 	  pipelineConfig.AuthorizeConfig.Oauth.AppUrl,	
					HttpClient: &http.Client{},
				}		
				connectorRegistry.Register(shopifyConnector)
		}
		if(error!=nil){
			panic("path not found")
		}
		
		pipeline:= NewPipeline(connections, random, stateStore, credentialService ,pipelineConfig).Build()
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

	return Runner{
		Scheduler: &scheduler,
		Worker: worker,
		Router: router,
	}
}