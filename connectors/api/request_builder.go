package api

import (
	"conecto/auth/connections"
	"context"
	"net/http"
)

type RequestBuilder interface {
	Build(ctx context.Context, cursor *PageCursor, connection connections.Connection) (*http.Request, error)
}
