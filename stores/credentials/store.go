package credentials

import (
	"conecto/core"
	"context"
	"errors"
)

var ErrCredentialNotFound = errors.New("credential not found")


type Store interface {
	Save(context context.Context, connection core.Connection, record EncryptedCredential) error
	Get(context context.Context, connection core.Connection) (EncryptedCredential, error)
	Close()error
}