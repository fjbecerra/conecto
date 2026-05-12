package factories

import (
	"conecto/core"
	"conecto/core/sinks/statestores"
	"fmt"
)

type StateStore struct{
	StateStoreConfig StateStoreConfig
	DBConnection DBConnection
}

func NewStateStore(stateStoreConfig StateStoreConfig, dbConnection DBConnection) *StateStore {
	return &StateStore {
		StateStoreConfig: stateStoreConfig,
		DBConnection: dbConnection,
	}
}

func (c *StateStore) Build() statestores.StateStore {
	if c.StateStoreConfig.Type == "" {
		c.StateStoreConfig.Type = MemoryStateStore
	}
	var stateStore statestores.StateStore
	switch c.StateStoreConfig.Type {
		case MemoryStateStore:
			stateStore = &statestores.MemoryStateStore{
				Store :  map[string]core.State{},				
			}
		case PostgresStateStore:
			stateStore = &statestores.PostgresStateStore{
				DB: c.DBConnection.DB,
			}
			if c.StateStoreConfig.AutoCreate {
				c.createPostgresTableStore()
			}

		default: panic("Unkown state store type")
	}
	return stateStore
}

func (c *StateStore) createPostgresTableStore() {
	query := `
	CREATE TABLE IF NOT EXISTS %s (
		id SERIAL PRIMARY KEY,
		pipeline_id TEXT NOT NULL,
		cursor JSONB NOT NULL,
		created_at TIMESTAMP DEFAULT NOW() NOT NULL
	)
	`
	_, err := c.DBConnection.DB.Exec(fmt.Sprintf(query, c.StateStoreConfig.Name))
	if err != nil {
		panic(err)
	}

	fmt.Println("table created or already exists")
}
