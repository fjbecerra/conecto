package auths

import "time"

type Token struct {
	AccessToken  string
	RefreshToken string
	Expiry       time.Time
}

type EncryptedToken struct {
	Ciphertext []byte
	Nonce      []byte
}
