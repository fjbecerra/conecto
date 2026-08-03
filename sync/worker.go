// worker.go
package sync

import (
	"context"
	"log/slog"
)

type Worker struct {
	buffer          Buffer
	syncService		*SyncService
}

func NewWorker(buffer Buffer, syncServic *SyncService) *Worker {
	return &Worker{
		buffer: buffer,
		syncService: syncServic,
	}
}

func (w *Worker) Run(ctx context.Context) {
	jobs := w.buffer.Consume()
	slog.Info("worker started")

	for {
		select {
			case <-ctx.Done():
				return

			case job := <-jobs:
				w.syncService.ExecuteJob(ctx, job)
		}
	}
}
