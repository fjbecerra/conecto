package credentials

import (
	"conecto/auth/connections"
	"context"
	"database/sql"
	"errors"
	"log"
)

type PostgresCredentialDB struct {
	db *sql.DB
}

func NewPostgresCredentialDB(dsn string) *PostgresCredentialDB {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected!")

	return &PostgresCredentialDB{
		db: db,
	}
}

func (p *PostgresCredentialDB) SaveCredential(context context.Context, connection connections.Connection, record EncryptedCredential) error {

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
		connection.ID,
		record.Ciphertext,
		record.Nonce,
		record.KeyVersion,
		record.ExpiresAt,
	)

	return err
}

func (p *PostgresCredentialDB) GetCredential(context context.Context, connection connections.Connection) (EncryptedCredential, error) {

	query := `
	SELECT
		ciphertext,
		nonce,
		key_version,
		expires_at
	FROM oauth_tokens
	WHERE pipeline_id = $1
	`

	var record EncryptedCredential

	err := p.db.QueryRowContext(
		context,
		query,
		connection.ID,
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

func (p *PostgresCredentialDB) Close() error {
	return p.db.Close()
}