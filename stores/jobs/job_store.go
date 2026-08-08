package jobs

import (
	"context"
	"time"
)

type JobStore interface {

	Create(ctx context.Context,job SyncJob) error

	MarkRunning(ctx context.Context,jobID string) error

	MarkCompleted(ctx context.Context,jobID string) error


	ScheduleRetry(
		ctx context.Context,
		jobID string,
		err error,
		nextRetry time.Time,
	) error


	MarkFailed(ctx context.Context, jobID string, err error) error
}