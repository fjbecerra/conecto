// sync_service.go
package sync

import (
	"conecto/auth/connections"
	"conecto/core/retry"
	"conecto/core/pipelines"
	"context"
	"time"

	"github.com/google/uuid"
)

type SyncService struct {
	buffer            Buffer
	registry         pipelines.Registry
	connectionStore  connections.Store
	jobRepository    JobRepository
	executor         *Executor
	retry 			retry.Executor
}

func NewSyncService(buffer Buffer, registry pipelines.Registry, connectionStore connections.Store, jobRepository JobRepository,executor *Executor,retry retry.Executor) *SyncService{
	return &SyncService{
		buffer: buffer,
		registry: registry,
		connectionStore: connectionStore,
		jobRepository: jobRepository,
		executor: executor,
		retry: retry,
	}
}

func (s *SyncService) ScheduleDueSyncs(
	ctx context.Context,
) error {

	connections, err := s.connectionStore.ClaimDueConnections(ctx)

	if err != nil {
		return err
	}

	for _, conn := range connections {

		_, err := s.registry.Get(
			conn.Provider,
		)

		if err != nil {
			continue
		}

		job := SyncJob{
			ID: uuid.NewString(),

			ConnectionID: conn.ID,

			PipelineID: conn.Provider,

			Status: JobPending,

			Attempt: 0,

			MaxRetries: 3,
		}

		err = s.jobRepository.Create(
			ctx,
			job,
		)

		if err != nil {
			continue
		}

		s.buffer.Publish(job)
	}

	return nil
}

func (s *SyncService) ScheduleConnectionSync(
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


	if err := s.jobRepository.Create(
		ctx,
		job,
	); err != nil {
		return err
	}


	s.buffer.Publish(job)

	return nil
}

func (s *SyncService) ExecuteJob(ctx context.Context, job SyncJob) error {
	

	err := s.retry.Do(ctx, func() error {
    		var err error
    		if err := s.jobRepository.MarkRunning(ctx, job.ID); err != nil {
				return err
			}

			err = s.executor.Execute(ctx, job)

			if err == nil {

				if err := s.jobRepository.MarkCompleted(ctx, job.ID); err != nil {
					return err
				}

				return s.connectionStore.MarkCompleted(
					ctx,
					job.ConnectionID,
					time.Now().Add(24*time.Hour),
				)
			}
    		return err
	})
	if err != nil {
		return s.jobRepository.MarkFailed(
			ctx,
			job.ID,
			err,)
	}

	return err

}