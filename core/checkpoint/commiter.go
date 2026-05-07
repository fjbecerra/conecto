package checkpoint

import (
	"conecto/core"
	"context"
)

type Record map[string]interface{}


type Committer interface {
	Commit(ctx context.Context, batch []core.Event, state core.State) error
}

