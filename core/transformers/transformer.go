package transformers

import (
	"conecto/core"
	"context"
)

type Transformer interface {
	Transform(ctx context.Context, batch []core.Event) ([]core.Event, error)
}