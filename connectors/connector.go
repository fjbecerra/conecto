package connectors

import (
	"conecto/auth/connections"
	"conecto/auth/credentials"
	"context"
)


type Connector interface {
	Name() string
	AuthorizeURL(ctx context.Context, connection connections.Connection, state string) (string,error)
	Exchange(ctx context.Context, connection connections.Connection, code string) (credentials.Credential,error)
	Refresh(ctx context.Context, connection connections.Connection, credential credentials.Credential) (credentials.Credential,error)
}