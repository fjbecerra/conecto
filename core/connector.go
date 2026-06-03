package core

import (
	"conecto/core/statestores"
	"context"
)

type Connector interface {
	Close() error 
	FetchBatch(context context.Context, state statestores.Cursor, ID string) (Batch, error)
}