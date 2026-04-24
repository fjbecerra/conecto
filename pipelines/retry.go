package pipelines

import (
	"context"
	"time"
)

type RetryRunner struct {
	inner      PipelineRunner
	maxRetries int
    backoff    time.Duration
}

func (r *RetryRunner) Run(ctx context.Context) error {
	var err error

	for i := 0; i <= r.maxRetries; i++ {

		err = r.inner.Run(ctx)
		if err == nil {
			return nil
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		time.Sleep(r.backoff * time.Duration(1<<i))
	}

	return err
}