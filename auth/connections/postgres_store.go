package connections

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"
)


type PostgresStore struct {
	db *sql.DB
}


func NewPostgresStore(dsn string) *PostgresStore {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected!")

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

func (s *PostgresStore) UpdateStatus(
	ctx context.Context,
	id string,
	status string,
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

func (s *PostgresStore) ClaimDueConnections(ctx context.Context) ([]Connection, error) {


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
			last_sync_at
		`,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()


	var result []Connection


	for rows.Next() {

		var conn Connection


		err := rows.Scan(
			&conn.ID,
			&conn.TenantID,
			&conn.Provider,
			&conn.Status,
			&conn.SyncStatus,
			&conn.NextSyncAt,
			&conn.LastSyncAt,
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

func (s *PostgresStore) MarkCompleted(
	ctx context.Context,
	id string,
	next time.Time,
) error {


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