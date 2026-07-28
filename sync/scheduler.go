// scheduler.go
package sync

import (
	"conecto/auth/connections"
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type Scheduler struct {
	syncService    *SyncService
	duration time.Duration
}

func NewScheduler(duration time.Duration, syncService *SyncService) Scheduler {
	return Scheduler{
		syncService: syncService,
		duration: duration,
	}
}

func (s *Scheduler) Run(ctx context.Context) {
	slog.Info(
		"scheduler started",
		"interval",
		s.duration,
	)

	ticker := time.NewTicker(s.duration)

	defer ticker.Stop()

	for range ticker.C {
		s.syncService.ScheduleDueSyncs(ctx)
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

	if err := s.syncService.jobRepository.Create(
		ctx,
		job,
	); err != nil {
		return err
	}

	s.syncService.buffer.Publish(job)

	return nil
}