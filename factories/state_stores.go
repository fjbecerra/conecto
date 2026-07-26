package factories

import (
	"conecto/core/statestores"
	"conecto/states"
	"database/sql"
	"fmt"
)

type StateStore struct{
	StateStoreConfig StateStoreConfig
	connections Connections
}

func NewStateStore(stateStoreConfig StateStoreConfig, connections Connections) *StateStore {
	return &StateStore {
		StateStoreConfig: stateStoreConfig,
		connections: connections,
	}
}

func (c *StateStore) Build() statestores.StateStore {
	if c.StateStoreConfig.Type == "" {
		c.StateStoreConfig.Type = MemoryStateStore
	}
	var stateStore statestores.StateStore
	switch c.StateStoreConfig.Type {
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
