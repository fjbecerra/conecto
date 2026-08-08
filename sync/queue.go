package sync

import "conecto/stores/jobs"


type Buffer interface {
	 Publish(job jobs.SyncJob)
	 Consume()<-chan jobs.SyncJob
}

type Queue struct {
	jobs chan jobs.SyncJob
}


func NewQueue(size int)*Queue{
	return &Queue{
		jobs: make(chan jobs.SyncJob,size),
	}
}


func (q *Queue) Publish(job jobs.SyncJob){

	q.jobs <- job
}


func (q *Queue) Consume()<-chan jobs.SyncJob {
	return q.jobs
}