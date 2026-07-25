package credentials

import (
	"conecto/auth/connections"
	"context"
)

type CredentialService interface {
	Save(context context.Context, connection connections.Connection, credential Credential) error
	Get(context context.Context, connection connections.Connection) (Credential, error)
	Close() error
}


