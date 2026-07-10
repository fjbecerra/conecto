package connections

import (
	"context"
	"database/sql"
	"encoding/json"
)


type PostgresStore struct {
	db *sql.DB
}


func NewPostgresStore(db *sql.DB) *PostgresStore {

	return &PostgresStore{
		db: db,
	}
}

func (s *PostgresStore) Get(ctx context.Context,id string) (Connection,error){

	query := `
	SELECT
		id,
		tenant_id,
		provider,
		external_id,
		metadata,
		status

	FROM connections

	WHERE id = $1
	`


	var c Connection

	var metadata []byte


	err :=
		s.db.QueryRowContext(
			ctx,
			query,
			id,
		).
		Scan(
			&c.ID,
			&c.TenantID,
			&c.Provider,
			&c.ExternalID,
			&metadata,
			&c.Status,
		)


	if err != nil {
		return Connection{},err
	}


	if len(metadata) > 0 {

		err =
			json.Unmarshal(
				metadata,
				&c.Metadata,
			)

		if err != nil {
			return Connection{},err
		}
	}


	return c,nil
}

func (s *PostgresStore) Save(ctx context.Context,connection Connection) error {


	metadata, err :=
		json.Marshal(
			connection.Metadata,
		)

	if err != nil {
		return err
	}


	query := `
	INSERT INTO connections (
		id,
		tenant_id,
		provider,
		external_id,
		metadata,
		created_at,
		updated_at
	)

	VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		NOW(),
		NOW()
	)

	ON CONFLICT (id)

	DO UPDATE SET

		tenant_id = EXCLUDED.tenant_id,

		provider = EXCLUDED.provider,

		external_id = EXCLUDED.external_id,

		metadata = EXCLUDED.metadata,

		updated_at = NOW()
	`


	_, err =
		s.db.ExecContext(
			ctx,
			query,

			connection.ID,

			connection.TenantID,

			connection.Provider,

			connection.ExternalID,

			metadata,
		)


	return err
}