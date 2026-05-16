package stores

import (
	"conecto/core"
	"database/sql"
	"errors"
)

type PostgresTokenDB struct {
	db *sql.DB
}

func NewPostgresTokenDB(db *sql.DB) *PostgresTokenDB {
	return &PostgresTokenDB{
		db: db,
	}
}

func (p *PostgresTokenDB) SaveToken(runtime core.Runtime, record TokenRecord) error {

	query := `
	INSERT INTO oauth_tokens (
		pipeline_id,
		provider,
		ciphertext,
		nonce,
		key_version,
		expires_at
	)
	VALUES ($1, $2, $3, $4, $5, $6)

	ON CONFLICT (pipeline_id, provider)

	DO UPDATE SET
		ciphertext  = EXCLUDED.ciphertext,
		nonce       = EXCLUDED.nonce,
		key_version = EXCLUDED.key_version,
		expires_at  = EXCLUDED.expires_at,
		updated_at  = NOW()
	`

	_, err := p.db.ExecContext(
		runtime.Context,
		query,
		runtime.PipelineId,
		runtime.Provider,
		record.Ciphertext,
		record.Nonce,
		record.KeyVersion,
		record.ExpiresAt,
	)

	return err
}

func (p *PostgresTokenDB) GetToken(runtime core.Runtime) (TokenRecord, error) {

	query := `
	SELECT
		ciphertext,
		nonce,
		key_version,
		expires_at
	FROM oauth_tokens
	WHERE pipeline_id = $1
	  AND provider = $2
	`

	var record TokenRecord

	err := p.db.QueryRowContext(
		runtime.Context,
		query,
		runtime.PipelineId,
		runtime.Provider,
	).Scan(
		&record.Ciphertext,
		&record.Nonce,
		&record.KeyVersion,
		&record.ExpiresAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return TokenRecord{}, ErrTokenNotFound
		}

		return TokenRecord{}, err
	}

	return record, nil
}

func (p *PostgresTokenDB) Close() error {
	return p.db.Close()
}