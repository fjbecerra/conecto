package credentials

import (
	"conecto/core"
	"context"
)

type CredentialRefresher interface {
    Refresh(ctx context.Context,connection core.Connection,credential Credential) (Credential, error)
}

type CredentialService interface {
	Save(context context.Context, connection core.Connection, credential Credential) error
	Get(context context.Context, connection core.Connection) (Credential, error)
	GetValid(ctx context.Context, connection core.Connection, refresher CredentialRefresher) (Credential, error)
	Close() error
}


