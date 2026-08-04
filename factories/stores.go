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
	storeSourceName string

}

type Stores struct {
	credentialStore credentials.Store
	stateStore statestores.StateStore
	connectionStore connections.Store
	jobRepository sync.JobRepository
}

func NewConectoStore(storeConfig DBConfig, connections Connections) *conectoStore {
	return &conectoStore {		
		connections: connections,
		storeSourceName: storeConfig.Source,
	}
}


func (c *conectoStore) Build() Stores {
	connection:= c.connections.connections[c.storeSourceName].OpenConnection
	connectionType:= c.connections.connections[c.storeSourceName].Type
	switch connectionType {
		case MemorySource:
			memory := connection.NewMemory
			return Stores{
				credentialStore: credentials.NewMemoryCredentialStore(memory),
				stateStore: &states.MemoryStateStore{
					Store :  memory,				
				},
				connectionStore: connections.NewMemoryStore(memory),
				jobRepository: sync.NewMemoryJobRepository(memory),
			}
			
		case PostgresSource:
			db := connection.DB
			connectionStore := connections.NewPostgresStore(db)
			createConnectionsTable("connections", db)
			credentialStore:= credentials.NewPostgresCredentialStore(db)
			createCredentialTable("credentials_store", db)
			stateStore:= states.NewStateStore(db)
			createStateTable("streams_state", db)			
			jobRepository:= sync.NewPostgresJobRepository(db)
			createJobRepositoryTable("sync_jobs", db)
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
			id SERIAL PRIMARY KEY,
			connection_id  UUID NOT NULL,
			ciphertext     BYTEA NOT NULL,
			nonce          BYTEA NOT NULL,
			key_version    TEXT NOT NULL,
			expires_at     TIMESTAMPTZ,
			created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			CONSTRAINT credentials_store_connection_id_unique UNIQUE (connection_id),
			FOREIGN KEY (connection_id) REFERENCES connections(id) ON DELETE CASCADE
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
		name TEXT NOT NULL,
		connection_id UUID NOT NULL,
		cursor JSONB NOT NULL,
		status TEXT NOT NULL,
		updated_at TIMESTAMP DEFAULT NOW() NOT NULL,
		CONSTRAINT streams_state_connection_id_unique UNIQUE (connection_id),
		CONSTRAINT streams_state_name_unique UNIQUE (name),
		FOREIGN KEY (connection_id) REFERENCES connections(id) ON DELETE CASCADE
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
	CREATE TABLE IF NOT EXISTS %s (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    provider TEXT NOT NULL,
    external_id TEXT,
    metadata JSONB,
    status TEXT NOT NULL,
	sync_status TEXT NOT NULL,
	last_sync_at TIMESTAMP,
	next_sync_at TIMESTAMP,
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
	CREATE TABLE IF NOT EXISTS %s (
    id UUID PRIMARY KEY,
    connection_id UUID NOT NULL,
    provider TEXT NOT NULL,
    status TEXT NOT NULL,
    attempt INT NOT NULL,
    max_retries INT NOT NULL,
    next_retry_at TIMESTAMP,
    last_error TEXT,
    created_at TIMESTAMP NOT NULL,
    started_at TIMESTAMP,
    finished_at TIMESTAMP,
	FOREIGN KEY (connection_id) REFERENCES connections(id) ON DELETE CASCADE
);
	`
	_, err := db.Exec(fmt.Sprintf(query, tableName))
	if err != nil {
		panic(err)
	}

	fmt.Println("table created or already exists")
}
