package credentials

import (
	"conecto/auth/connections"
	"context"
	"errors"
)

var ErrCredentialNotFound = errors.New("credential not found")


type Store interface {
	Save(context context.Context, connection connections.Connection, record EncryptedCredential) error
	Get(context context.Context, connection connections.Connection) (EncryptedCredential, error)
	Close()error
}