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
	slog.Info("worker started")
	for job := range w.buffer.Consume() {
		w.syncService.ExecuteJob(ctx, job)		
	}
}
