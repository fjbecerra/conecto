// scheduler.go
package sync

import (
	"conecto/auth/connections"
	"context"
	"time"

	"github.com/google/uuid"
)

type Scheduler struct {
	SyncService    *SyncService
}

func (s *Scheduler) Run(ctx context.Context) {

	ticker := time.NewTicker(time.Minute)

	defer ticker.Stop()

	for range ticker.C {
		s.SyncService.ScheduleDueSyncs(ctx)
	}
}

func (s *Scheduler) CreateJob(
	ctx context.Context,
	conn connections.Connection,
) error {

	job := SyncJob{
		ID: uuid.NewString(),

		ConnectionID: conn.ID,

		PipelineID: conn.Provider,

		Status: JobPending,

		Attempt: 0,

		MaxRetries: 3,
	}

	if err := s.SyncService.jobRepository.Create(
		ctx,
		job,
	); err != nil {
		return err
	}

	s.SyncService.queue.Publish(job)

	return nil
}