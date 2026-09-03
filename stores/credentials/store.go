package credentials

import (
	"context"
	"errors"
)

var ErrCredentialNotFound = errors.New("credential not found")


type Store interface {
	Save(context context.Context, record EncryptedCredential) error
	GetByConnectionId(context context.Context, connectionId string) (EncryptedCredential, error)
	Close()error
}