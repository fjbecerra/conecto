package sync

import (
	"context"
	"errors"
	"time"
)

var ErrJobNotFound = errors.New("job not found")

type MemoryJobRepository struct {

	jobs map[string]any
}

func NewMemoryJobRepository(store map[string]any) *MemoryJobRepository {
	return &MemoryJobRepository{
		jobs: store,
	}
}

func (r *MemoryJobRepository) Create(
	ctx context.Context,
	job SyncJob,
) error {

	r.jobs[job.ID] = &job

	return nil
}

func (r *MemoryJobRepository) MarkRunning(
	ctx context.Context,
	jobID string,
) error {

	
	job, ok := r.jobs[jobID].(SyncJob)

	if !ok {
		return ErrJobNotFound
	}


	job.Status = JobRunning


	return nil
}

func (r *MemoryJobRepository) MarkCompleted(
	ctx context.Context,
	jobID string,
) error {


	job, ok := r.jobs[jobID].(SyncJob)

	if !ok {
		return ErrJobNotFound
	}


	job.Status = JobCompleted

	job.NextRetryAt = nil


	return nil
}

func (r *MemoryJobRepository) ScheduleRetry(
	ctx context.Context,
	jobID string,
	err error,
	nextRetry time.Time,
) error {


	job, ok := r.jobs[jobID].(SyncJob)

	if !ok {
		return ErrJobNotFound
	}


	job.Status = JobRetrying

	job.Attempt++

	job.NextRetryAt = &nextRetry

	job.LastError = err.Error()


	return nil
}

func (r *MemoryJobRepository) MarkFailed(
	ctx context.Context,
	jobID string,
	err error,
) error {


	job, ok := r.jobs[jobID].(SyncJob)

	if !ok {
		return ErrJobNotFound
	}


	job.Status = JobFailed

	job.LastError = err.Error()

	job.NextRetryAt = nil


	return nil
}

