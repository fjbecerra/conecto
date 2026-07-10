package credentials

import (
	"context"
	"errors"
)

var ErrCredentialNotFound = errors.New("credential not found")


type Store interface {
	SaveCredential(context context.Context, ID string, record EncryptedCredential) error
	GetCredential(context context.Context, ID string) (EncryptedCredential, error)
	Close()error
}