// scheduler.go
package sync

import (
	"context"
	"log/slog"
	"time"

)

type Scheduler struct {
	syncService    *SyncService
	duration time.Duration
}

func NewScheduler(duration time.Duration, syncService *SyncService) Scheduler {
	return Scheduler{
		syncService: syncService,
		duration: duration,
	}
}

func (s *Scheduler) Run(ctx context.Context) {
	slog.Info(
		"scheduler started",
		"interval",
		s.duration,
	)

	ticker := time.NewTicker(s.duration)

	defer ticker.Stop()

	for {
		select {
		case <- ctx.Done():
			return
		case <- ticker.C:
			s.syncService.ScheduleDueSyncs(ctx)
		}		
	}
}
