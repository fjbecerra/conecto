package connections

import (
	"context"
	"time"
)

type Store interface {

	Get(ctx context.Context,id string) (Connection,error)

	Save(ctx context.Context, connection Connection) error

	UpdateStatus(ctx context.Context,id string,status string) error


	ClaimDueConnections(ctx context.Context) ([]Connection,error)


	MarkRunning(ctx context.Context,id string) error


	MarkCompleted(ctx context.Context, id string, next time.Time) error

}