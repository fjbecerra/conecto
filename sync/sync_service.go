// sync_service.go
package sync

import (
	"conecto/auth/connections"
	"conecto/core/pipelines"
	"conecto/core/retry"
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type SyncService struct {
	buffer            Buffer
	registry         pipelines.Registry
	connectionStore  connections.Store
	jobRepository    JobRepository
	retry 			retry.Executor
}

func NewSyncService(buffer Buffer, registry pipelines.Registry, connectionStore connections.Store, jobRepository JobRepository,retry retry.Executor) *SyncService{
	return &SyncService{
		buffer: buffer,
		registry: registry,
		connectionStore: connectionStore,
		jobRepository: jobRepository,
		retry: retry,
	}
}

func (s *SyncService) ScheduleDueSyncs(ctx context.Context,) error {

	connections, err := s.connectionStore.ClaimDueConnections(ctx)

	if err != nil {
		return err
	}

	for _, conn := range connections {

		pipeline, err := s.registry.Get(conn.Provider)

		if err != nil {
			slog.Error("Pipeline due to sync not register", "pipeline", pipeline.ID)
			continue
		}

		job, err := s.createJob(ctx, conn)

		if err != nil {
			continue
		}

		s.buffer.Publish(job)
	}

	return nil
}


//this should be a backfill of last 90 days
func (s *SyncService) ScheduleConnectionSync(ctx context.Context, conn connections.Connection) error {

	job, err := s.createJob(ctx, conn)

	if err != nil {
		return err
	}

	s.buffer.Publish(job)

	return nil
}

func (s *SyncService) createJob(ctx context.Context, conn connections.Connection,
) (SyncJob, error) {

	job := SyncJob{
		ID: uuid.NewString(),
		ConnectionID: conn.ID,
		Provider: conn.Provider,
		Status: JobPending,
		Attempt: 0,
		MaxRetries: 3,
	}

	return job,s.jobRepository.Create(ctx,job)
}


func (s *SyncService) ExecuteJob(ctx context.Context, job SyncJob) error {
	

	err := s.retry.Do(ctx, func() error {
    		var err error
    		if err := s.jobRepository.MarkRunning(ctx, job.ID); err != nil {
				return err
			}

			err = s.execute(ctx, job)

			if err == nil {

				if err := s.jobRepository.MarkCompleted(ctx, job.ID); err != nil {
					return err
				}

				return s.connectionStore.MarkCompleted(
					ctx,
					job.ConnectionID,
					time.Now().Add(24*time.Hour),//todo. put this in the config
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

func (e *SyncService) execute(ctx context.Context, job SyncJob) error {

	conn, err := e.connectionStore.Get(ctx,job.ConnectionID)

	if err != nil {
		return err
	}

	pipeline, err := e.registry.Get(job.Provider)

	if err != nil {
		return err
	}

	for _, stream := range pipeline.Streams {


		err := stream.Run(ctx,conn)

		if err != nil {
			//cannot leave
			return err
		}
	}

	return nil
}