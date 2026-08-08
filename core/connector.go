package core

import (
	"conecto/core/statestores"
	"context"
)

type Connector interface {
	FetchBatch(context context.Context, state statestores.Cursor, connection Connection) (Batch, error)
}