package stores

import (
	"context"
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

func (p *PostgresTokenDB) SaveToken(context context.Context, id string, record TokenRecord) error {

	query := `
	INSERT INTO oauth_tokens (
		pipeline_id,
		ciphertext,
		nonce,
		key_version,
		expires_at
	)
	VALUES ($1, $2, $3, $4, $5)

	ON CONFLICT (pipeline_id)

	DO UPDATE SET
		ciphertext  = EXCLUDED.ciphertext,
		nonce       = EXCLUDED.nonce,
		key_version = EXCLUDED.key_version,
		expires_at  = EXCLUDED.expires_at,
		updated_at  = NOW()
	`

	_, err := p.db.ExecContext(
		context,
		query,
		id,
		record.Ciphertext,
		record.Nonce,
		record.KeyVersion,
		record.ExpiresAt,
	)

	return err
}

func (p *PostgresTokenDB) GetToken(context context.Context, ID string) (TokenRecord, error) {

	query := `
	SELECT
		ciphertext,
		nonce,
		key_version,
		expires_at
	FROM oauth_tokens
	WHERE pipeline_id = $1
	`

	var record TokenRecord

	err := p.db.QueryRowContext(
		context,
		query,
		ID,
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