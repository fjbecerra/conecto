package auths

import (
	"conecto/connectors/_http/auths/stores"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
)


type AESGCMTokenStore struct {
	Store  stores.Store
	keyManager KeyManager
}

func NewADBTokenStore(store stores.Store, keyManager KeyManager) *AESGCMTokenStore {
	return &AESGCMTokenStore{
		Store:  store,
		keyManager: keyManager,
	}
}

func (s *AESGCMTokenStore) Get(context context.Context, ID string) (Token, error) {

	record, err := s.Store.GetToken(context, ID)

	if err != nil {
		return Token{}, err
	}

	key, err := s.keyManager.Get(record.KeyVersion)
	if err != nil {
		return Token{}, err
	}

	plain, err := decrypt(
		key,
		record.Ciphertext,
		record.Nonce,
	)

	if err != nil {
		return Token{}, err
	}

	var token Token

	err = json.Unmarshal(plain, &token)

	return token, err
}

func (s *AESGCMTokenStore) Save(context context.Context, ID string, token Token) error {

	raw, err := json.Marshal(token)
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

	record := stores.TokenRecord{
		Ciphertext: ciphertext,
		Nonce:      nonce,
		KeyVersion: version,
		ExpiresAt:  token.Expiry,
	}

	return s.Store.SaveToken(context,ID, record)
}

func (s *AESGCMTokenStore) Close() error {
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

