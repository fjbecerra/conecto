// worker.go
package sync

import "context"

type Worker struct {
	queue            *Queue
	syncService		SyncService
}

func (w *Worker) Run(ctx context.Context) {

	for job := range w.queue.Consume() {
		w.syncService.ExecuteJob(ctx, job)		
	}
}
