package stores

import (
	"context"
	"time"
)

type MemoryStoreToken struct {
	Store map[string]any
}

func NewMemoryStoreToken(store map[string]any) *MemoryStoreToken {
	return &MemoryStoreToken{
		Store: store,
	}
}

func (p *MemoryStoreToken) SaveToken(context context.Context, ID string, record TokenRecord) error {

	p.Store["pipeline_id"] = ID
	p.Store["ciphertext"] = record.Ciphertext
	p.Store["nonce"] = record.Nonce
	p.Store["key_version"] = record.KeyVersion
	p.Store["updated_at"] = record.ExpiresAt

	return nil
}

func (p *MemoryStoreToken) GetToken(context context.Context, ID string) (TokenRecord, error) {
	return TokenRecord{
		Ciphertext: p.Store["ciphertext"].([]byte),
		Nonce: p.Store["nonce"].([]byte),
		KeyVersion:  p.Store["key_version"].(string),
		ExpiresAt:  p.Store["updated_at"].(time.Time),
	},nil	
}

func (p *MemoryStoreToken) Close() error{
	return nil
}