package api

import (
	"conecto/core"
	"conecto/core/statestores"
	"context"
	"net/http"
)

type RequestBuilder interface {
	Build(ctx context.Context, cursor *PageCursor, connection core.Connection, syncState *statestores.SyncState) (*http.Request, error)
}
