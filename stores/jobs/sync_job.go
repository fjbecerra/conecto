package jobs

import (
	"conecto/core"
	"time"
)

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

	Connection core.Connection

	Provider core.Provider

	Status JobStatus

	Attempt int

	MaxRetries int

	NextRetryAt *time.Time

	LastError string
}