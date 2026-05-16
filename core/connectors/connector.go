package connectors

import (
	"conecto/core"
)

type Connector interface {
	Close() error 
	FetchBatch(runtime core.Runtime, state core.Cursor) (core.Batch, error)
}