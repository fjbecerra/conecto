package api

import (
	"conecto/core"
	"context"
	"net/http"
)

type RequestBuilder interface {
	Build(ctx context.Context, cursor *PageCursor, connection core.Connection) (*http.Request, error)
}
