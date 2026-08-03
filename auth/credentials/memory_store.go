package credentials

import (
	"conecto/auth/connections"
	"context"
	"time"
)

type MemoryCredentialStore struct {
	store map[string]any
}

func NewMemoryCredentialStore(store map[string]any) *MemoryCredentialStore {
	return &MemoryCredentialStore{
		store: store,
	}
}

func (p *MemoryCredentialStore) Save(context context.Context, connection connections.Connection, record EncryptedCredential) error {

	p.store["pipeline_id"] = connection.ID
	p.store["ciphertext"] = record.Ciphertext
	p.store["nonce"] = record.Nonce
	p.store["key_version"] = record.KeyVersion
	p.store["updated_at"] = record.ExpiresAt

	return nil
}

func (p *MemoryCredentialStore) Get(context context.Context, connection connections.Connection) (EncryptedCredential, error) {
	return EncryptedCredential{
		Ciphertext: p.store["ciphertext"].([]byte),
		Nonce: p.store["nonce"].([]byte),
		KeyVersion:  p.store["key_version"].(string),
		ExpiresAt:  p.store["updated_at"].(*time.Time),
	},nil	
}

func (p *MemoryCredentialStore) Close() error{
	return nil
}