package statestores

import (
	"conecto/core"
	"context"
)

type StateStore interface {
	Load(ctx context.Context, pipelineID string) (core.State, error)
	Save(ctx context.Context, pipelineID string, state core.State) error
}

