package committers

import (
	"conecto/core"
	"conecto/core/transformers"
	"context"
)


type Committer interface{
	CommitBatch(ctx context.Context, pipelineID string, batch core.Batch, transformer transformers.Transformer) error
 	close() error
}

