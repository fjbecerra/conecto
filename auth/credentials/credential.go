package credentials

import "time"

type Credential struct {
	Type string
	Data map[string]string
	Expiry *time.Time
}

type EncryptedCredential struct {
	Ciphertext []byte
	Nonce []byte
	KeyVersion string
	ExpiresAt *time.Time
}
