package resources

import (
	"conecto/core"
	"conecto/core/engines"
	"conecto/core/statestores"
	"conecto/resources/memory"
	"conecto/resources/postgres"
	"conecto/resources/shopify"
	"conecto/shared"
	"conecto/shared/clients"
	"conecto/shared/config"
	"conecto/stores/credentials"
	"errors"
	"fmt"
)

type ResourceType string
const(
	PostgresResourceType ResourceType = "postgres"
	MemoryResourceType 	ResourceType = "memory"
    ShopifyResourceType ResourceType = "shopify")

type ResourceName string


type ResourcesRegistry struct{
    registry map[ResourceName]Resource
    credentialService credentials.CredentialService
	StateStore statestores.StateStore
    retry shared.Retry    
}

func NewResourceRegistry(
    credentialService credentials.CredentialService, 
    stateStore statestores.StateStore, 
    retry shared.Retry) *ResourcesRegistry{
    return &ResourcesRegistry{
        registry: make(map[ResourceName]Resource),
        credentialService: credentialService,
        StateStore: stateStore,
        retry: retry,
    }
}

type ResourceFactory func(config.Resource,)Resource

type Resource interface {
	 Close() error
     Connector(cfg config.ConfigBytes) engines.ConnectorRunnable
     Sink(cfg config.ConfigBytes, fieldSpecs config.FieldsSpecs)  engines.SinkCommiter
     Transformers() []core.Transformer
}

func (rr *ResourcesRegistry) Register(resources []config.Resource) error {
    factories := map[ResourceType]ResourceFactory{
        PostgresResourceType: func(r config.Resource) Resource {
            clientConfig,_:= config.Unmarshal[clients.PostgresConfig](r.Client, config.FormatJSON)
            client:= clients.CreatePostgresClient(clientConfig)
            retry:= rr.retry.CreateRetryExecutor(r.Retry)
            var  postgresResourceConfig postgres.PostgresResourceConfig
            if(r.Config != nil){
                postgresResourceConfig,_ = config.Unmarshal[postgres.PostgresResourceConfig](*r.Config, config.FormatJSON)
            }            
            return postgres.NewPostgresResource(client, &retry, rr.StateStore,postgresResourceConfig)
        },

        ShopifyResourceType: func(r config.Resource) Resource {
            clientConfig, _:= config.Unmarshal[clients.HttpClientConfig](r.Client, config.FormatJSON)
            client:=clients.CreateHttpClint(&clientConfig)           
            retry:= rr.retry.CreateRetryExecutor(r.Retry)
            var shopifyConfig shopify.ShopifyResourceConfig
            if(r.Config != nil){
                shopifyConfig, _= config.Unmarshal[shopify.ShopifyResourceConfig](*r.Config, config.FormatJSON)
            }
            return shopify.NewShopifyResource(client, rr.credentialService, shopifyConfig, &retry)
        },

        MemoryResourceType: func(r config.Resource) Resource {
            retry:= rr.retry.CreateRetryExecutor(r.Retry)
            store:=[]map[string]any{}
            return memory.NewMemoryResource(&store, &retry, rr.StateStore)     
        },
    }
    
    for _, resource := range resources {
        factory, ok := factories[ResourceType(resource.Type)]
        if !ok {
            return fmt.Errorf("unsupported resource")
        }

        rr.registry[ResourceName(resource.Name)] = factory(resource)
    }

    return nil
}

func (rr *ResourcesRegistry) Get(name ResourceName) Resource{
    return rr.registry[name]
}

func (rr *ResourcesRegistry) CloseAll() error {
    var err error
    for _, resource := range rr.registry {
        err = errors.Join(err, resource.Close())
    }
    return err
}

