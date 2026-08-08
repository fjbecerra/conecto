package factories

import (
	"conecto/http_server/oauth"
	"conecto/http_server/oauth/state"
	"conecto/pipelines"
	"conecto/resources"
	"conecto/shared"
	"conecto/shared/config"
	"conecto/stores"
	"conecto/sync"
	"context"
	"encoding/base64"
	"time"

	"github.com/go-chi/chi/v5"
)


type Runner interface {
	Shutdown() error
	RunScheduler(ctx context.Context)
	RunWorker(ctx context.Context) 
	RunRouter() *chi.Mux
}

type ConectoRunner struct {
	scheduler *sync.Scheduler
	worker *sync.Worker
	routers *chi.Mux
	resourcesRegistry resources.ResourcesRegistry
}

func NewConectoRunner(cfg config.Conecto) *ConectoRunner{
	stores:= stores.NewStores(cfg.Store)
	retry:= shared.NewRetry()
	routers := chi.NewRouter()
	credentialKey,_ := base64.StdEncoding.DecodeString(cfg.Security.CredentialKey)
	
	credentialService := stores.CredentialService(credentialKey)
	resourceRegistry:= resources.NewResourceRegistry(credentialService, stores.StateStore(),retry)
	resourceRegistry.Register(cfg.Resources)
	pipelineRegistry:= pipelines.NewPipelineRegistry(resourceRegistry)
	pipelineRegistry.Register(cfg.Pipelines)

	//this may go eventually into a registry
	var buffer sync.Buffer
	if(cfg.Sync.Buffer.BufferType == config.QueueType){
		buffer = sync.NewQueue(cfg.Sync.Buffer.Size)
	}else{
		panic("unknown buffer type")
	}

	syncService := sync.NewSyncService(
		buffer, 
		pipelineRegistry, 
		stores.ConnectionStore(), 
		stores.JobStore(),
		retry.CreateRetryExecutor(&cfg.Sync.Retry),
	)
	duration,_:= time.ParseDuration(cfg.Sync.Scheduler.Duration)
	scheduler := sync.NewScheduler(duration,syncService )
	worker := sync.NewWorker(buffer, syncService)	
	stateSignerKey,_ := base64.StdEncoding.DecodeString(cfg.Security.StateSignerKey)
	stateSigner := state.NewHMACStateSigner(
			stateSignerKey,
			10*time.Minute)
	oauthService := oauth.NewService(
		stores.ConnectionStore(), 
		credentialService, 
		stateSigner,
		resourceRegistry, 
		syncService,
	)
	oauthHandler := oauth.NewHandler(oauthService)
	oauthHandler.RegisterRoutes(routers)
	

	return &ConectoRunner{
		scheduler: &scheduler,
		worker: worker,
		routers: routers,
		resourcesRegistry: *resourceRegistry,
	}
}

func (c *ConectoRunner) Shutdown() error{
	err := c.resourcesRegistry.CloseAll()
	if err != nil {
		return err
	}
	return nil
}

func (c *ConectoRunner) RunScheduler(ctx context.Context){
	 c.scheduler.Run(ctx)
}

func (c *ConectoRunner) RunWorker(ctx context.Context){
	 c.worker.Run(ctx)
}

func (c *ConectoRunner) RunRouter() *chi.Mux{
	return c.routers
}

