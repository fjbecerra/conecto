package credentials

import (
	"conecto/core"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
)


type AESGCMCredentialService struct {
	Store  Store
	credentialKey []byte
}

func NewAESGCMCredentialService(store Store, credentialKey []byte) *AESGCMCredentialService {
	return &AESGCMCredentialService{
		Store:  store,
		credentialKey: credentialKey,
	}
}

func (s *AESGCMCredentialService) Get(context context.Context, connection core.Connection) (Credential, error) {

	record, err := s.Store.Get(context, connection)

	if err != nil {
		return Credential{}, err
	}

	plain, err := decrypt(
		s.credentialKey,
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

func (s *AESGCMCredentialService) Save(context context.Context, connection core.Connection, credential Credential) error {

	raw, err := json.Marshal(credential)
	if err != nil {
		return err
	}

	ciphertext, nonce, err := encrypt(s.credentialKey, raw)
	if err != nil {
		return err
	}

	record := EncryptedCredential{
		Ciphertext: ciphertext,
		Nonce:      nonce,
		ExpiresAt:  credential.Expiry,
	}

	return s.Store.Save(context, connection, record)
}

func (s *AESGCMCredentialService) Close() error {
	return s.Store.Close()
}

func (s *AESGCMCredentialService) GetValid(ctx context.Context,connection core.Connection,refresher CredentialRefresher,
) (Credential, error) {

    credential, err := s.Get(ctx, connection)
    if err != nil {
        return Credential{}, err
    }

    if !credential.IsExpired() {
        return credential, nil
    }

    if refresher == nil {
        return Credential{}, errors.New(
            "credential expired but no refresher is configured",
        )
    }

    credential, err = refresher.Refresh(
        ctx,
        connection,
        credential,
    )
    if err != nil {
        return Credential{}, err
    }

    if err := s.Save(ctx, connection, credential); err != nil {
        return Credential{}, err
    }

    return credential, nil
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

