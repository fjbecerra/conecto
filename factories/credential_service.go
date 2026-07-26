package factories

import (
	"conecto/auth/credentials"
	"conecto/core/statestores"
	"conecto/states"
	"database/sql"
	"fmt"
)

type CredentialService struct{
	CredentialConfig CredentialConfig
	connections Connections
}

func NewCredentialService(credentialConfig CredentialConfig, connections Connections) *CredentialService {
	return &CredentialService {
		CredentialConfig: credentialConfig,
		connections: connections,
	}
}

func (c *CredentialService) Build() credentials.CredentialService {
	if c.CredentialStoreConfig.Type == "" {
		c.CredentialStoreConfig.Type = MemoryStateStore
	}
	var stateStore statestores.StateStore
	switch c.CredentialStoreConfig.Type {
		case MemoryStateStore:
			stateStore = &states.MemoryStateStore{
				Store :  map[string]statestores.State{},				
			}
		case PostgresStateStore:
			connection:= c.connections[c.StateStoreConfig.Source].OpenDB()	
			stateStore = &states.PostgresStateStore{
				DB: connection,
			}
			if c.StateStoreConfig.AutoCreate {
				createPostgresTableStore(c.StateStoreConfig, connection)
			}

		default: panic("Unkown state store type")
	}
	return stateStore
}

func createPostgresTableStore(config StateStoreConfig, db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS %s (
		id SERIAL PRIMARY KEY,
		pipeline_id TEXT NOT NULL UNIQUE,
		cursor JSONB NOT NULL,
		status TEXT NOT NULL,
		updated_at TIMESTAMP DEFAULT NOW() NOT NULL
	)
	`
	_, err := db.Exec(fmt.Sprintf(query, config.Name))
	if err != nil {
		panic(err)
	}

	fmt.Println("table created or already exists")
}
