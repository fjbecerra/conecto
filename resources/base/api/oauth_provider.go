package api

import (
	"conecto/core"
	"conecto/stores/credentials"
	"context"
)

type OAuthProvider interface {
	Name() string
	AuthorizeURL(ctx context.Context, connection core.Connection, state string) (string,error)
	Exchange(ctx context.Context, connection core.Connection, code string) (credentials.Credential,error)
}	