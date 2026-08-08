package stores

import (
	"conecto/core/statestores"
	"conecto/stores/connections"
	"conecto/stores/credentials"
	"conecto/stores/jobs"
	"conecto/stores/states"
)

type MemoryStore struct {
	mstore *map[string]any
}

func NewMemoryStore(mstore *map[string]any)*MemoryStore{	
	 return &MemoryStore{
        mstore: mstore,
    }
}

func (m* MemoryStore) CredentialService(credentialKey []byte )  credentials.CredentialService {
	credentialStore:= credentials.NewMemoryCredentialStore(*m.mstore)	
	return credentials.NewAESGCMCredentialService(credentialStore,credentialKey)
}

func (m* MemoryStore) StateStore() statestores.StateStore {
	return states.NewMemoryStateStore(*m.mstore)			
}

func (m* MemoryStore) ConnectionStore() connections.Store {
	return connections.NewMemoryConnectionStore(*m.mstore)	
}

func (m* MemoryStore) JobStore() jobs.JobStore {
	return jobs.NewMemoryJobStore(*m.mstore)
}

func (m* MemoryStore) Close() error{
	return nil
}
