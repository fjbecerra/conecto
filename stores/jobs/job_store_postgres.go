package jobs

import (
	"context"
	"database/sql"
	"time"
)

type PostgresJobStore struct {
	db *sql.DB
}

func NewPostgresJobStore(db *sql.DB) *PostgresJobStore {
	return &PostgresJobStore{
		db: db,
	}
}

func (r *PostgresJobStore) Create(ctx context.Context, job SyncJob) error {

	_, err := r.db.ExecContext(
		ctx,
		`
		INSERT INTO sync_jobs (
			id,
			connection_id,
			provider,
			status,
			attempt,
			max_retries,
			next_retry_at,
			last_error,
			created_at
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			NOW()
		)
		`,
		job.ID,
		job.Connection.ID,
		job.Provider,
		job.Status,
		job.Attempt,
		job.MaxRetries,
		job.NextRetryAt,
		job.LastError,
	)

	return err
}

func (r *PostgresJobStore) MarkRunning(ctx context.Context,jobID string) error {

	_, err := r.db.ExecContext(
		ctx,
		`
		UPDATE sync_jobs

		SET
			status = 'running',
			started_at = NOW()

		WHERE id = $1
		`,
		jobID,
	)

	return err
}

func (r *PostgresJobStore) MarkCompleted(ctx context.Context,jobID string) error {

	_, err := r.db.ExecContext(
		ctx,
		`
		UPDATE sync_jobs

		SET
			status = 'completed',
			finished_at = NOW()

		WHERE id = $1
		`,
		jobID,
	)

	return err
}

func (r *PostgresJobStore) ScheduleRetry(
	ctx context.Context,
	jobID string,
	err error,
	nextRetry time.Time,
) error {


	_, dbErr := r.db.ExecContext(
		ctx,
		`
		UPDATE sync_jobs

		SET
			status = 'retrying',
			attempt = attempt + 1,
			next_retry_at = $2,
			last_error = $3

		WHERE id = $1
		`,
		jobID,
		nextRetry,
		err.Error(),
	)


	return dbErr
}

func (r *PostgresJobStore) MarkFailed(ctx context.Context, jobID string, err error) error {


	_, dbErr := r.db.ExecContext(
		ctx,
		`
		UPDATE sync_jobs

		SET
			status = 'failed',
			finished_at = NOW(),
			last_error = $2

		WHERE id = $1
		`,
		jobID,
		err.Error(),
	)


	return dbErr
}



