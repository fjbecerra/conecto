package sinks

import (
	"conecto/core"
	"context"
)

type Sink interface {
    Open(ctx context.Context) error
	Close() error
	WriteBatch(ctx context.Context, batch []core.Event) error
	Commit(ctx context.Context, cursor core.Cursor) error
}