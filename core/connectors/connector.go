package connectors

import (
	"conecto/core"
	"context"
)

type Connector interface {
	Open(ctx context.Context, state core.Cursor) error 
	Close() error 
	FetchBatch(ctx context.Context, state core.Cursor) (core.Batch, error)
}