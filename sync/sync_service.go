package sync

import (
	"conecto/core"
	"conecto/core/retry"
	"conecto/pipelines"
	"conecto/stores/connections"
	"conecto/stores/jobs"
	"context"
	"errors"
	"time"
	"github.com/google/uuid"
)

type SyncService struct {
	buffer            Buffer
	pipelineRegistry  pipelines.PipelineRegistry
	connectionStore  connections.Store
	jobStore    	jobs.JobStore
	retry 			retry.Executor
}

func NewSyncService(buffer Buffer, pipelineRegistry pipelines.PipelineRegistry, connectionStore connections.Store, jobStore jobs.JobStore,retry retry.Executor) *SyncService{
	return &SyncService{
		buffer: buffer,
		pipelineRegistry: pipelineRegistry,
		connectionStore: connectionStore,
		jobStore: jobStore,
		retry: retry,
	}
}

func (s *SyncService) ScheduleDueSyncs(ctx context.Context,) error {

	connections, err := s.connectionStore.ClaimDueConnections(ctx)

	if err != nil {
		return err
	}

	for _, conn := range connections {

		job, err := s.createJob(ctx, conn)

		if err != nil {
			continue
		}

		s.buffer.Publish(job)
	}

	return nil
}


//this should be a backfill of last 90 days
func (s *SyncService) ScheduleConnectionSync(ctx context.Context, conn core.Connection) error {

	job, err := s.createJob(ctx, conn)

	if err != nil {
		return err
	}

	s.buffer.Publish(job)

	return nil
}

func (s *SyncService) createJob(ctx context.Context, conn core.Connection,
) (jobs.SyncJob, error) {

	job := jobs.SyncJob{
		ID: uuid.NewString(),
		Connection: conn,
		Provider: conn.Provider,
		Status: jobs.JobPending,
		Attempt: 0,
		MaxRetries: 3,
	}

	return job,s.jobStore.Create(ctx,job)
}


func (s *SyncService) ExecuteJob(ctx context.Context, job jobs.SyncJob) error {
	

	err := s.retry.Do(ctx, func() error {
    		var err error
    		if err := s.jobStore.MarkRunning(ctx, job.ID); err != nil {
				return err
			}

			err = s.execute(ctx, job)

			if err == nil {

				if err := s.jobStore.MarkCompleted(ctx, job.ID); err != nil {
					return err
				}

				return s.connectionStore.MarkCompleted(
					ctx,
					job.Connection.ID,
					time.Now().Add(24*time.Hour),//todo. put this in the config
				)
			}
    		return err
	})  
	if err != nil {
		return s.jobStore.MarkFailed(
			ctx,
			job.ID,
			err,)
	}

	return err

}

func (e *SyncService) execute(ctx context.Context, job jobs.SyncJob) error {
	
	pipeline := e.pipelineRegistry.Get(job.Provider)

	var errs error
	for _, stream := range pipeline.Streams {
		errs = errors.Join(errs, stream.Run(ctx,job.Connection))
	}

	return errs
}