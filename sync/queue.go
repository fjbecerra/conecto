package sync

type Queue struct {
	jobs chan SyncJob
}


func NewQueue(size int)*Queue{
	return &Queue{
		jobs: make(chan SyncJob,size),
	}
}


func (q *Queue) Publish(job SyncJob){

	q.jobs <- job
}


func (q *Queue) Consume()<-chan SyncJob {
	return q.jobs
}