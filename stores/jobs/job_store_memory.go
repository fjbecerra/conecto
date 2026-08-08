package jobs

import (
	"context"
	"errors"
	"time"
)

var ErrJobNotFound = errors.New("job not found")

type MemoryJobStore struct {

	jobs map[string]any
}

func NewMemoryJobStore(store map[string]any) *MemoryJobStore {
	return &MemoryJobStore{
		jobs: store,
	}
}

func (r *MemoryJobStore) Create(
	ctx context.Context,
	job SyncJob,
) error {

	r.jobs[job.ID] = &job

	return nil
}

func (r *MemoryJobStore) MarkRunning(
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

func (r *MemoryJobStore) MarkCompleted(
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

func (r *MemoryJobStore) ScheduleRetry(
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

func (r *MemoryJobStore) MarkFailed(
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

