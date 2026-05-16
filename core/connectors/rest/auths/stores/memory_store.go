package stores

import (
	"conecto/core"
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

func (p *MemoryStoreToken) SaveToken(runtime core.Runtime, record TokenRecord) error {

	p.Store["pipeline_id"] = runtime.PipelineId
	p.Store["provider"] = runtime.Provider
	p.Store["ciphertext"] = record.Ciphertext
	p.Store["nonce"] = record.Nonce
	p.Store["key_version"] = record.KeyVersion
	p.Store["updated_at"] = record.ExpiresAt

	return nil
}

func (p *MemoryStoreToken) GetToken(runtime core.Runtime) (TokenRecord, error) {
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