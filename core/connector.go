package core

import (
	"conecto/auth/connections"
	"conecto/core/statestores"
	"context"
)

type Connector interface {
	Close() error 
	FetchBatch(context context.Context, state statestores.Cursor, connection connections.Connection) (Batch, error)
}