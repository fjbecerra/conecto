package credentials

import (
	"context"
	"database/sql"
	"errors"
)

type PostgresCredentialStore struct {
	db *sql.DB
}

func NewPostgresCredentialStore(db *sql.DB) *PostgresCredentialStore {
	return &PostgresCredentialStore{
		db: db,
	}
}

func (p *PostgresCredentialStore) Save(context context.Context, record EncryptedCredential) error {

	query := `
	INSERT INTO credentials_store (
		connection_id,
		ciphertext,
		nonce,
		key_version,
		expires_at
	)
	VALUES ($1, $2, $3, $4, $5)

	ON CONFLICT (connection_id)

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
		record.Connection.ID,
		record.Ciphertext,
		record.Nonce,
		record.KeyVersion,
		record.ExpiresAt,
	)

	return err
}

func (p *PostgresCredentialStore) GetByConnectionId(context context.Context, connectionId string) (EncryptedCredential, error) {

	query := `
	SELECT
		ciphertext,
		nonce,
		key_version,
		expires_at
	FROM credentials_store
	WHERE connection_id = $1
	`

	var record EncryptedCredential

	err := p.db.QueryRowContext(
		context,
		query,
		connectionId,
	).Scan(
		&record.Ciphertext,
		&record.Nonce,
		&record.KeyVersion,
		&record.ExpiresAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return EncryptedCredential{}, ErrCredentialNotFound
		}

		return EncryptedCredential{}, err
	}

	return record, nil
}

func (p *PostgresCredentialStore) Close() error {
	return p.db.Close()
}