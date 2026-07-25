package credentials

import (
	"conecto/auth/connections"
	"context"
	"time"
)

type MemoryStoreCredential struct {
	Store map[string]any
}

func NewMemoryStoreCredential(store map[string]any) *MemoryStoreCredential {
	return &MemoryStoreCredential{
		Store: store,
	}
}

func (p *MemoryStoreCredential) SaveCredential(context context.Context, connection connections.Connection, record EncryptedCredential) error {

	p.Store["pipeline_id"] = connection.ID
	p.Store["ciphertext"] = record.Ciphertext
	p.Store["nonce"] = record.Nonce
	p.Store["key_version"] = record.KeyVersion
	p.Store["updated_at"] = record.ExpiresAt

	return nil
}

func (p *MemoryStoreCredential) GetCredential(context context.Context, connection connections.Connection) (EncryptedCredential, error) {
	return EncryptedCredential{
		Ciphertext: p.Store["ciphertext"].([]byte),
		Nonce: p.Store["nonce"].([]byte),
		KeyVersion:  p.Store["key_version"].(string),
		ExpiresAt:  p.Store["updated_at"].(*time.Time),
	},nil	
}

func (p *MemoryStoreCredential) Close() error{
	return nil
}