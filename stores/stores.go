package stores

import (
	"conecto/core/statestores"
	"conecto/shared/clients"
	"conecto/shared/config"
	"conecto/stores/connections"
	"conecto/stores/credentials"
	"conecto/stores/jobs"
)

type StoreType string
const (
	PostgresStoreType StoreType = "postgres"
	MemoryStoreType StoreType = "memory"
)

type Stores interface {
	CredentialService(key []byte) credentials.CredentialService
	StateStore() statestores.StateStore
	ConnectionStore() connections.Store
	JobStore() jobs.JobStore
	Close() error
}

type storyFactory func (config.Store)Stores

func NewStores(cfg config.Store) Stores {
	factories := map[StoreType]storyFactory{
		PostgresStoreType: func(cfg config.Store)Stores{
			postgrestConfig,_ := config.Unmarshal[clients.PostgresConfig](cfg.Client, config.FormatJSON)
			client:= clients.CreatePostgresClient(postgrestConfig)
			return NewPostgresStores(client)
		},
		MemoryStoreType: func(cfg config.Store) Stores {
			store:=make(map[string]any)
			return NewMemoryStore(&store)
		},
	}
	
	factory, ok := factories[StoreType(cfg.Type)]
	if !ok {
		panic("store not supported")
	}

    return factory(cfg)

}

