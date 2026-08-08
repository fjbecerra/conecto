package connections

import (
	"conecto/core"
	"context"
	"time"
)

type Store interface {

	Get(ctx context.Context,id string) (core.Connection,error)

	Save(ctx context.Context, connection core.Connection) error

	UpdateStatus(ctx context.Context,id string,status core.ConnectionStatus) error


	ClaimDueConnections(ctx context.Context) ([]core.Connection,error)


	MarkRunning(ctx context.Context,id string) error


	MarkCompleted(ctx context.Context, id string, next time.Time) error

}