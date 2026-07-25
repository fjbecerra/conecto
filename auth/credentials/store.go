package credentials

import (
	"conecto/auth/connections"
	"context"
	"errors"
)

var ErrCredentialNotFound = errors.New("credential not found")


type Store interface {
	SaveCredential(context context.Context, connection connections.Connection, record EncryptedCredential) error
	GetCredential(context context.Context, connection connections.Connection) (EncryptedCredential, error)
	Close()error
}