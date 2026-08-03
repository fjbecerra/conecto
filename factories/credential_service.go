package factories

import (
	"conecto/auth/credentials"
	"encoding/base64"
	"os"
)


type CredentialService struct{
	credentialStore credentials.Store
}

func NewCredentialService(credentialStore credentials.Store) *CredentialService {
	return &CredentialService {		
		credentialStore: credentialStore,
	}
}

func (c *CredentialService) Build() credentials.CredentialService {	
	v1, _ := base64.StdEncoding.DecodeString(os.Getenv("CREDENTIAL_ENCRYPTION_KEY_V1"))
	keys := map[string][]byte{
		"v1": v1,
	}
	keyManager := credentials.NewStaticKeyManager(keys, "v1")
	return credentials.NewAESGCMCredentialService(c.credentialStore, keyManager)
}

