package stores

import (
	"conecto/core/statestores"
	"conecto/shared/clients"
	"conecto/stores/connections"
	"conecto/stores/credentials"
	"conecto/stores/jobs"
	"conecto/stores/states"
	"database/sql"
	"fmt"
)

	
	


type PostgresStore struct{
	db *sql.DB
}

func NewPostgresStores(client *clients.PostgresClient)*PostgresStore {
	db:= client.Get()		
	createConnectionsTable("connections", db)
	createCredentialTable("credentials_store", db)
	createStateTable("streams_state", db)	
	createJobRepositoryTable("sync_jobs", db)
	return &PostgresStore{
		db: client.Get(),
	}
}

func (p* PostgresStore) CredentialService(credentialKey []byte) credentials.CredentialService {
	credentialStore:= credentials.NewPostgresCredentialStore(p.db)	
	return credentials.NewAESGCMCredentialService(credentialStore, credentialKey)
}

func (p* PostgresStore) StateStore() statestores.StateStore {
	return states.NewPostgresStateStore(p.db)			
}

func (p* PostgresStore) ConnectionStore() connections.Store {
	return connections.NewPostgresConnectionStore(p.db)	
}

func (p* PostgresStore) JobStore() jobs.JobStore {
	return jobs.NewPostgresJobStore(p.db)
}

func (p* PostgresStore) Close() error{
	return p.db.Close()
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
		CONSTRAINT streams_state_connection_id_unique UNIQUE (connection_id, name),
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
	resource_name TEXT NOT NULL,
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
