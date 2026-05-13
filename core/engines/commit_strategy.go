package engines

import (
	"conecto/core"
	"context"
)


type CommitStrategy interface{
	Commit(runtime core.Runtime, batch core.Batch) error
 	Shutdown(ctx context.Context) error
}

