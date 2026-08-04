// sync_job.go
package sync

import "time"

type JobStatus string

const (
	JobPending   JobStatus = "pending"
	JobRunning   JobStatus = "running"
	JobRetrying  JobStatus = "retrying"
	JobCompleted JobStatus = "completed"
	JobFailed    JobStatus = "failed"
)

type SyncJob struct {
	ID string

	ConnectionID string

	Provider string

	Status JobStatus

	Attempt int

	MaxRetries int

	NextRetryAt *time.Time

	LastError string
}