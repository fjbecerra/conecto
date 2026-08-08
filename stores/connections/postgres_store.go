package connections

import (
	"conecto/core"
	"context"
	"database/sql"
	"encoding/json"
	"time"
)


type PostgresStore struct {
	db *sql.DB
}


func NewPostgresConnectionStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{
		db: db,
	}
}

func (s *PostgresStore) Get(ctx context.Context,id string) (core.Connection,error){

	query := `
	SELECT
		id,
		tenant_id,
		provider,
		metadata,
		status,
		resource_name

	FROM connections

	WHERE id = $1
	`


	var c core.Connection

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
			&metadata,
			&c.Status,
			&c.ResourceName,
		)


	if err != nil {
		return core.Connection{},err
	}


	if len(metadata) > 0 {

		err =
			json.Unmarshal(
				metadata,
				&c.Metadata,
			)

		if err != nil {
			return core.Connection{},err
		}
	}


	return c,nil
}

func (s *PostgresStore) Save(ctx context.Context,connection core.Connection) error {


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
		updated_at,
		resource_name
	)

	VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		NOW(),
		NOW(),
		$6,
	)

	ON CONFLICT (id)

	DO UPDATE SET

		tenant_id = EXCLUDED.tenant_id,

		provider = EXCLUDED.provider,

		external_id = EXCLUDED.external_id,

		metadata = EXCLUDED.metadata,

		updated_at = NOW(),

		resource_name=EXCLUDED.resource_name
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

			connection.ResourceName,
		)


	return err
}

func (s *PostgresStore) UpdateStatus(
	ctx context.Context,
	id string,
	status core.ConnectionStatus,
) error {

	_, err := s.db.ExecContext(
		ctx,
		`
		UPDATE connections
		SET status = $2
		WHERE id = $1
		`,
		id,
		status,
	)

	return err
}

func (s *PostgresStore) ClaimDueConnections(ctx context.Context) ([]core.Connection, error) {


	rows, err := s.db.QueryContext(
		ctx,
		`
		UPDATE connections

		SET sync_status = 'queued'

		WHERE id IN (

			SELECT id
			FROM connections

			WHERE status = 'connected'
			AND sync_status = 'idle'
			AND next_sync_at <= NOW()

			FOR UPDATE SKIP LOCKED
		)

		RETURNING
			id,
			tenant_id,
			provider,
			status,
			sync_status,
			next_sync_at,
			last_sync_at,
			resource_name
		`,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()


	var result []core.Connection


	for rows.Next() {

		var conn core.Connection


		err := rows.Scan(
			&conn.ID,
			&conn.TenantID,
			&conn.Provider,
			&conn.Status,
			&conn.SyncStatus,
			&conn.NextSyncAt,
			&conn.LastSyncAt,
			&conn.ResourceName,
		)

		if err != nil {
			return nil, err
		}


		result = append(
			result,
			conn,
		)
	}


	return result, rows.Err()
}

func (s *PostgresStore) MarkRunning(
	ctx context.Context,
	id string,
) error {


	_, err := s.db.ExecContext(
		ctx,
		`
		UPDATE connections

		SET sync_status = 'running'

		WHERE id = $1
		`,
		id,
	)


	return err
}

func (s *PostgresStore) MarkCompleted(ctx context.Context, id string, next time.Time) error {


	_, err := s.db.ExecContext(
		ctx,
		`
		UPDATE connections

		SET
			sync_status = 'idle',
			last_sync_at = NOW(),
			next_sync_at = $2

		WHERE id = $1
		`,
		id,
		next,
	)


	return err
}