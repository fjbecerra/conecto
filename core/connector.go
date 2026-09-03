package core

import (
	"conecto/core/statestores"
	"context"
)

type Connector interface {
	FetchBatch(context context.Context, state statestores.Cursor, connection Connection, watermark *string) (Batch, error)
}