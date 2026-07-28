package credentials

import (
	"conecto/auth/connections"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
)


type AESGCMCredentialService struct {
	Store  Store
	keyManager KeyManager
}

func NewAESGCMCredentialService(store Store, keyManager KeyManager) *AESGCMCredentialService {
	return &AESGCMCredentialService{
		Store:  store,
		keyManager: keyManager,
	}
}

func (s *AESGCMCredentialService) Get(context context.Context, connection connections.Connection) (Credential, error) {

	record, err := s.Store.Get(context, connection)

	if err != nil {
		return Credential{}, err
	}

	key, err := s.keyManager.Get(record.KeyVersion)
	if err != nil {
		return Credential{}, err
	}

	plain, err := decrypt(
		key,
		record.Ciphertext,
		record.Nonce,
	)

	if err != nil {
		return Credential{}, err
	}

	var token Credential

	err = json.Unmarshal(plain, &token)

	return token, err
}

func (s *AESGCMCredentialService) Save(context context.Context, connection connections.Connection, credential Credential) error {

	raw, err := json.Marshal(credential)
	if err != nil {
		return err
	}

	version, key, err := s.keyManager.Latest()
	if err != nil {
		return err
	}

	ciphertext, nonce, err := encrypt(key, raw)
	if err != nil {
		return err
	}

	record := EncryptedCredential{
		Ciphertext: ciphertext,
		Nonce:      nonce,
		KeyVersion: version,
		ExpiresAt:  credential.Expiry,
	}

	return s.Store.Save(context, connection, record)
}

func (s *AESGCMCredentialService) Close() error {
	return s.Store.Close()
}

func encrypt(key []byte, plain []byte) (ciphertext []byte, nonce []byte, err error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce = make([]byte, gcm.NonceSize())

	if _, err := rand.Read(nonce); err != nil {
		return nil, nil, err
	}

	ciphertext = gcm.Seal(nil, nonce, plain, nil)

	return ciphertext, nonce, nil
}

func decrypt(key []byte, ciphertext []byte, nonce []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return gcm.Open(nil, nonce, ciphertext, nil)
}

