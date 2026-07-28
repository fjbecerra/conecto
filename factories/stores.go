package factories

import (
	"conecto/auth/connections"
	"conecto/auth/credentials"
	"conecto/core/statestores"
	"conecto/states"
	"conecto/sync"
	"database/sql"
	"fmt"
)

type conectoStore struct{
	connections Connections
	storeType StoreType
	storeSourceName string

}

type Stores struct {
	credentialStore credentials.Store
	stateStore statestores.StateStore
	connectionStore connections.Store
	jobRepository sync.JobRepository
}

func NewConectoStore(storeConfig StoreConfig, connections Connections) *conectoStore {
	return &conectoStore {		
		connections: connections,
		storeType: storeConfig.StoreType,
		storeSourceName: storeConfig.Source,
	}
}


func (c *conectoStore) Build() Stores {
	switch c.storeType {
		case MemoryStore:
			return Stores{
				credentialStore: credentials.NewMemoryCredentialStore(),
				stateStore: &states.MemoryStateStore{
					Store :  map[string]statestores.State{},				
				},
				connectionStore: connections.NewMemoryStore(),
				jobRepository: sync.NewMemoryJobRepository(),
			}
			
		case PostgresStore:
			connection:= c.connections[c.storeSourceName].OpenDB()	
			credentialStore:= credentials.NewPostgresCredentialStore(connection)
			createCredentialTable("crendentials_store", connection)
			stateStore:= states.NewStateStore(connection)
			createStateTable("streams_state", connection)
			connectionStore := connections.NewPostgresStore(connection)
			createConnectionsTable("connections", connection)
			jobRepository:= sync.NewPostgresJobRepository(connection)
			createJobRepositoryTable("sync_jobs", connection)
			return Stores{
				credentialStore: credentialStore,
				stateStore: stateStore,
				connectionStore: connectionStore,
				jobRepository: jobRepository,
			}
			
		default: panic("Unkown state store type")
	}

}


func createCredentialTable(tableName string, db *sql.DB) {
	query := `
		CREATE TABLE IF NOT EXISTS %s (
			connection_id    TEXT NOT NULL,
			ciphertext     BYTEA NOT NULL,
			nonce          BYTEA NOT NULL,
			key_version    TEXT NOT NULL,
			expires_at     TIMESTAMPTZ,
			created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			PRIMARY KEY (connection_id)
	);
	`
	_, err := db.Exec(fmt.Sprintf(query, tableName))
	if err != nil {
		panic(err)
	}

	fmt.Println("table created or already exists")
}

func createStateTable(tableName string, db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS %s (
		id SERIAL PRIMARY KEY,
		pipeline_id TEXT NOT NULL UNIQUE,
		cursor JSONB NOT NULL,
		status TEXT NOT NULL,
		updated_at TIMESTAMP DEFAULT NOW() NOT NULL
	)
	`
	_, err := db.Exec(fmt.Sprintf(query, tableName))
	if err != nil {
		panic(err)
	}

	fmt.Println("table created or already exists")
}

func createConnectionsTable(tableName string, db *sql.DB) {
	query := `
	CREATE TABLE %s (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    provider TEXT NOT NULL,
    external_id TEXT,
    metadata JSONB,
    status TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
	`
	_, err := db.Exec(fmt.Sprintf(query, tableName))
	if err != nil {
		panic(err)
	}

	fmt.Println("table created or already exists")
}

func createJobRepositoryTable(tableName string, db *sql.DB) {
	query := `
	CREATE TABLE %s (
    id UUID PRIMARY KEY,
    connection_id UUID NOT NULL,
    pipeline_id TEXT NOT NULL,
    status TEXT NOT NULL,
    attempt INT NOT NULL,
    max_retries INT NOT NULL,
    next_retry_at TIMESTAMP,
    last_error TEXT,
    created_at TIMESTAMP NOT NULL,
    started_at TIMESTAMP,
    finished_at TIMESTAMP
);
	`
	_, err := db.Exec(fmt.Sprintf(query, tableName))
	if err != nil {
		panic(err)
	}

	fmt.Println("table created or already exists")
}