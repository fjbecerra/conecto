package core

import (
	"conecto/core/statestores"
	"context"
)

type Connector interface {
	FetchBatch(context context.Context, state statestores.State, connection Connection,) (Batch, error)
}